package main

import (
	"bytes"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"time"
)

type appSettings struct {
	url       *url.URL
	token     string
	n_threads uint16
}

func write_to_api(client *http.Client, url *url.URL, token string, data []byte) error {

	req, err := http.NewRequest("POST", url.String(), bytes.NewBuffer(data))
	if err != nil {
		return err
	}

	req.Header.Add("Content-Type", "application/json")
	req.Header.Add("Authorization", fmt.Sprintf("Bearer %s", token))

	resp, err := client.Do(req)
	if err != nil {
		return err
	}

	_, err = io.ReadAll(resp.Body)
	if err != nil {
		return err
	}

	return nil
}

func write_message_to_api(client *http.Client, settings appSettings) error {
	data, err := make_payload()
	if err != nil {
		return err
	}

	err = write_to_api(client, settings.url, settings.token, data)
	if err != nil {
		return err
	}
	return nil
}

func spamit(settings appSettings) {
	client := &http.Client{}
	for {
		err := write_message_to_api(client, settings)
		if err != nil {
			panic(err)
		}
	}
}

func main() {
	settings, err := get_settings()
	if err != nil {
		panic(err)
	}

	for i := uint16(0); i < settings.n_threads; i++ {
		go spamit(settings)
	}
	time.Sleep(time.Second * duration)
}
