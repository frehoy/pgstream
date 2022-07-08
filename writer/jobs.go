package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	// "time"
)

type Work struct {
	Foo      string `json:"foo"`
	Category string `json:"category"`
}

type create_job_payload struct {
	Work Work `json:"work"`
}

type get_job_payload struct {
	Worker_id int `json:"worker_id"`
}

type mark_job_done_payload struct {
	Status string `json:"status"`
}

type Job struct {
	Job_id int  `json:"job_id"`
	Work   Work `json:"work"`
}

func submit_job(client *http.Client, settings appSettings) error {
	data := create_job_payload{
		Work: Work{
			Foo:      "From job client",
			Category: randomCategory(),
		},
	}
	payload, err := json.Marshal(data)
	if err != nil {
		return err
	}
	call_api(client, settings.url, settings.token, payload, "/jobs", "POST")
	return nil
}

func do_work(client *http.Client, settings appSettings) (err error) {
	job, err := get_job(client, settings)
	if err != nil {
		return err
	}
	run_job(job)
	err = mark_job_done(client, settings, job)
	if err != nil {
		return err
	}
	return nil
}

// TODO: Try again if we don't get a job.
func get_job(client *http.Client, settings appSettings) (job Job, err error) {

	payload, err := json.Marshal(get_job_payload{Worker_id: 5})
	if err != nil {
		return job, err
	}
	resp, err := call_api(client, settings.url, settings.token, payload, "/rpc/start_job", "POST")
	if err != nil {
		return job, err
	}

	err = json.Unmarshal(resp, &job)
	if err != nil {
		return job, err
	}
	if job.Job_id == 0 {
		job, err = get_job(client, settings)
		if err != nil {
			panic(err)
		}
	}

	return job, nil
}

func run_job(job Job) {
	// sleep_time := seededRand.Intn(10)
	// time.Sleep(time.Duration(sleep_time) * time.Millisecond)
}

func mark_job_done(client *http.Client, settings appSettings, job Job) (err error) {
	payload, err := json.Marshal(mark_job_done_payload{Status: "finished"})
	if err != nil {
		return err
	}
	u := *settings.url
	u.Path = "jobs"
	params := url.Values{}
	params.Add("id", fmt.Sprintf("eq.%d", job.Job_id))
	u.RawQuery = params.Encode()
	req, err := http.NewRequest("PATCH", u.String(), bytes.NewBuffer(payload))
	if err != nil {
		return err
	}
	req.Header.Add("Content-Type", "application/json")
	req.Header.Add("Authorization", fmt.Sprintf("Bearer %s", settings.token))

	resp, err := client.Do(req)
	if err != nil {
		return err
	}

	_, err = io.ReadAll(resp.Body)
	if err != nil {
		return err
	}

	err = resp.Body.Close()
	if err != nil {
		return err
	}

	return nil

}
