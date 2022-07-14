package main

import (
	"testing"
)

func TestWorker(t *testing.T) {
	image_job_payload := ImageJobPayload{
		Source_url:        "http://example.com/image.jpeg",
		Target_resolution: ImageResolution{x: 1024, y: 768},
	}

	data_job_payload := DataJobPayload {
		From: "Start",
		To: "End",
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

	var expected bool
	var response bool
	expected = true

	response, _ = job_payload_with_image.Payload.DoWork()
	if response != expected {
		t.Fatalf("%v != %v", response, expected)
	}

	response, _ = job_payload_with_image.DoJob()
	if response != expected {
		t.Fatalf("%v != %v", response, expected)
	}

	response, _ = job_payload_with_data.DoJob()
	if response != expected {
		t.Fatalf("%v != %v", response, expected)
	}
}