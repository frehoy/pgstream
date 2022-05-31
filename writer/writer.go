package main

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"math/rand"
	"net/http"
	"net/url"
	"os"
	"time"
)

const charset = "abcdefghijklmnopqrstuvwxyz" +
	"ABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"

const token = "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJyb2xlIjoiYXBpX3VzZXIifQ.W9Fo49rxMbSVnhdK1lzjMwCgf_1MZCPy9GNbt9j10ds"
const parallelism = 16
const duration = 60

var seededRand *rand.Rand = rand.New(
	rand.NewSource(time.Now().UnixNano()))

func StringWithCharset(length int, charset string) string {
	b := make([]byte, length)
	for i := range b {
		b[i] = charset[seededRand.Intn(len(charset))]
	}
	return string(b)
}

func String(length int) string {
	return StringWithCharset(length, charset)
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

type appSettings struct {
	url   *url.URL
	token string
}

func get_url_from_env() (*url.URL, error) {
	env_url, is_set := os.LookupEnv("WRITE_ENDPOINT")
	if !is_set {
		return nil, errors.New("WRITE_ENDPOINT env var not set.")
	}

	parsed_url, err := url.Parse(env_url)
	if err != nil {
		return nil, err
	}

	return parsed_url, nil
}

func get_settings() appSettings {

	url_from_env, err := get_url_from_env()
	if err != nil {
		panic(err)
	}

	return appSettings{
		url:   url_from_env,
		token: token,
	}
}

type message struct {
	Message_type string `json:"message_type"`
	Search_terms string `json:"search_terms"`
}

type payload_struct struct {
	Message message `json:"message"`
}

func make_payload() ([]byte, error) {

	random_word := String(10)
	d := payload_struct{
		Message: message{
			Message_type: "search",
			Search_terms: random_word,
		},
	}
	json_payload, err := json.Marshal(d)
	if err != nil {
		return nil, err
	}

	payload := json_payload
	return payload, nil
}

func write_message_to_api(client *http.Client, settings appSettings) error {
	data, err := make_payload()
	if err != nil {
		return err
	}

	err = write_to_api(client, settings.url, token, data)
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
	settings := get_settings()
	for i := 0; i < parallelism; i++ {
		go spamit(settings)
	}
	time.Sleep(time.Second * duration)
}
