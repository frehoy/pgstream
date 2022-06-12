package main

import (
	"encoding/json"
	"math/rand"
	"time"
)

const charset = "abcdefghijklmnopqrstuvwxyz" +
	"ABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"

var seededRand *rand.Rand = rand.New(
	rand.NewSource(time.Now().UnixNano()))

func randomCategory() string {
	categories := []string{
		"Page",
		"Image",
		"Map",
		"Video",
		"News",
	}
	category := categories[rand.Intn(len(categories))]
	return category
}

func stringWithCharset(length int, charset string) string {
	b := make([]byte, length)
	for i := range b {
		b[i] = charset[seededRand.Intn(len(charset))]
	}
	return string(b)
}

func randomString(length int) string {
	return stringWithCharset(length, charset)
}

type message struct {
	Message_type string `json:"message_type"`
	Search       search `json:"search"`
}

type payload_struct struct {
	Message message `json:"message"`
}

type search struct {
	Query    string `json:"query"`
	Category string `json:"category"`
}

func make_payload() ([]byte, error) {
	data := payload_struct{
		Message: message{
			Message_type: "search",
			Search: search{
				Query:    randomString(3),
				Category: randomCategory(),
			},
		},
	}
	json_payload, err := json.Marshal(data)
	if err != nil {
		return nil, err
	}

	payload := json_payload
	return payload, nil
}
