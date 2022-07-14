package main

import (
	"fmt"
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
	DoJob() (bool, error)
}

type Work interface {
	DoWork() (bool, error)
}

type JobPayload struct {
	Id      int
	Payload Work
}

func (image_job_payload ImageJobPayload) DoWork() (done bool, err error) {
	fmt.Printf("Doing  work on image image_job_payload %+v\n", image_job_payload)
	return true, nil
}

func (data_job_payload DataJobPayload) DoWork() (done bool, err error) {
	fmt.Printf("Doing  work on data data_job_payload %+v\n", data_job_payload)
	return true, nil
}

func (jp JobPayload) DoJob() (done bool, err error) {
	fmt.Printf("Doing job on j %+v\n", jp)
	jp.Payload.DoWork()
	return true, nil
}
