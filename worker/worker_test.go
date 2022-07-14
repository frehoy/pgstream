package main

import (
	"testing"
)

func TestWorker(t *testing.T) {
	image_job_payload := ImageJobPayload{
		Source_url:        "http://example.com/image.jpeg",
		Target_resolution: ImageResolution{x: 1024, y: 768},
	}

	data_job_payload := DataJobPayload{
		From:   "Start",
		To:     "End",
		Window: "Hourly",
	}

	job_payload_with_image := JobPayload{
		Id:      1,
		Payload: image_job_payload,
	}

	job_payload_with_data := JobPayload{
		Id:      1,
		Payload: data_job_payload,
	}

	err := job_payload_with_image.Payload.DoWork()
	if err != nil {
		t.Fatalf("%v != %v", err, nil)
	}

	err = job_payload_with_image.DoJob()
	if err != nil {
		t.Fatalf("%v != %v", err, nil)
	}

	err = job_payload_with_data.DoJob()
	if err != nil {
		t.Fatalf("%v != %v", err, nil)
	}
}

func submit_jobs(n_jobs int, incoming_jobs chan JobPayload) {
	for i := 0; i < n_jobs; i++ {
		incoming_jobs <- JobPayload{Id: i, Payload: DataJobPayload{}, Status: "in_progress"}
	}
}

// DoJobs takes jobs through a channel, runs them concurrently and
// puts them back on another channel
func TestDoJobs(t *testing.T) {

	n_jobs := 1000
	concurrency := 2000
	// Use tiny channels (size 1) to demonstrate we aren't deadlocking
	incoming_jobs := make(chan JobPayload, 1)
	processed_jobs := make(chan JobPayload, 1)
	defer close(incoming_jobs)
	defer close(processed_jobs)

	go DoJobsConcurrently(incoming_jobs, processed_jobs, concurrency)
	go submit_jobs(n_jobs, incoming_jobs)

	var finished_jobs []JobPayload
	for i := 0; i < n_jobs; i++ {
		j := <-processed_jobs
		finished_jobs = append(finished_jobs, j)
	}
	for _, j := range finished_jobs {
		if j.Status != "finished" {
			t.Fatalf("Job %v is not finished", j)
		}
	}
}
