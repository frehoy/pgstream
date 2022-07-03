package main

import (
	"encoding/json"
	"net/http"
)

type Work struct {
	Foo      string `json:"foo"`
	Category string `json:"category"`
}

type job_payload struct {
	Work Work `json:"work"`
}

func submit_job(client *http.Client, settings appSettings) error {
	data := job_payload{
		Work: Work{
			Foo:      "From job client",
			Category: randomCategory(),
		},
	}
	payload, err := json.Marshal(data)
	if err != nil {
		return err
	}
	write_to_api(client, settings.url, settings.token, payload, "/jobs")
	return nil
}
