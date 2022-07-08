package main

import (
	"bytes"
	"fmt"
	"io"
	"log"
	"net/http"
	"net/url"
	"time"
)

type appSettings struct {
	url       *url.URL
	token     string
	n_threads uint16
}

func call_api(client *http.Client, url *url.URL, token string, payload []byte, path string, method string) (response []byte, err error) {
	// Copy the url before setting path so we don't get hairy with threads
	this_url := *url
	this_url.Path = path
	req, err := http.NewRequest("POST", this_url.String(), bytes.NewBuffer(payload))
	if err != nil {
		return nil, err
	}

	req.Header.Add("Content-Type", "application/json")
	req.Header.Add("Authorization", fmt.Sprintf("Bearer %s", token))

	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}

	data, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	err = resp.Body.Close()
	if err != nil {
		return nil, err
	}

	return data, nil
}

func write_message_to_api(client *http.Client, settings appSettings) error {
	payload, err := make_payload()
	if err != nil {
		return err
	}

	_, err = call_api(client, settings.url, settings.token, payload, "/events", "POST")
	if err != nil {
		return err
	}
	return nil
}

func spam_events(settings appSettings, sent_messages chan int) {
	client := &http.Client{}
	for {
		err := write_message_to_api(client, settings)
		if err != nil {
			panic(err)
		}
		sent_messages <- 1
	}
}

func spam_jobs(settings appSettings, submitted_jobs chan int) {
	client := &http.Client{}
	for {
		err := submit_job(client, settings)
		if err != nil {
			panic(err)
		}
		submitted_jobs <- 1
	}
}

func show_count_per_second(message string, count_channel chan int) {
	var total int
	tick := time.Tick(1 * time.Second)

	for {
		select {
		case <-tick:
			log.Println(fmt.Sprintf("%s: %d", message, total))
			total = 0
		default:
			total += <-count_channel
		}
	}
}

func do_jobs(settings appSettings, done_jobs chan int) (err error) {
	client := &http.Client{}
	for {
		err = do_work(client, settings)
		if err != nil {
			panic(err)
		}
		done_jobs <- 1
	}
}

func main() {
	settings := get_settings_or_fail()
	sent_messages := make(chan int)
	submitted_jobs := make(chan int)
	finished_jobs := make(chan int)

	for i := uint16(0); i < 1; i++ {
		go spam_events(settings, sent_messages)
	}
	for i := uint16(0); i < 1; i++ {
		go spam_jobs(settings, submitted_jobs)
	}
	for i := uint16(0); i < 16; i++ {
		go do_jobs(settings, finished_jobs)
	}

	go show_count_per_second("sent_messages", sent_messages)
	go show_count_per_second("submitted_jobs", submitted_jobs)
	go show_count_per_second("finished_jobs", finished_jobs)
	time.Sleep(time.Second * duration)
}
