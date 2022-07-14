package main

import (
	"math/rand"
	"time"
)

type ImageResolution struct {
	x int
	y int
}

type ImageJobPayload struct {
	Source_url        string
	Target_resolution ImageResolution
}

type DataJobPayload struct {
	From   string
	To     string
	Window string
}

type Job interface {
	DoJob() error
}

type Work interface {
	DoWork() error
}

type JobPayload struct {
	Id         int
	Payload    Work
	Status     string
	Try_number int
}

func (image_job_payload ImageJobPayload) DoWork() (err error) {
	return nil
}

func (data_job_payload DataJobPayload) DoWork() (err error) {
	return nil
}

// DoJob takes a pointer to JobPayload so it can mutate it, to set Status
func (jp *JobPayload) DoJob() (err error) {
	time.Sleep(time.Duration(rand.Intn(10)) * time.Millisecond)
	err = jp.Payload.DoWork()
	if err != nil {
		jp.Status = "failed"
	} else {
		jp.Status = "finished"
	}
	return err
}

func RunJobsFromChannel(incoming_jobs <-chan JobPayload, processed_jobs chan<- JobPayload) {
	for {
		job, chan_open := <-incoming_jobs
		if !chan_open {
			break
		}
		err := job.DoJob()
		if err != nil {

		}
		processed_jobs <- job
	}
}

func DoJobsConcurrently(incoming_jobs <-chan JobPayload, processed_jobs chan<- JobPayload, concurrency int) {
	for i := 0; i < concurrency; i++ {
		go RunJobsFromChannel(incoming_jobs, processed_jobs)
	}
}
