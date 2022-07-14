package main

import (
	"errors"
	"net/url"
	"os"
)

type configuration struct {
	base_url     *url.URL
	secret_token string
}

func config_from_env() (configuration, error) {

	// base_url
	base_url_str, is_set := os.LookupEnv("base_url")
	if !is_set {
		return configuration{}, errors.New("base_url env var not set")
	}
	base_url, err := url.Parse(base_url_str)
	if err != nil {
		return configuration{}, err
	}

	// secret_token
	secret_token, is_set := os.LookupEnv("secret_token")
	if !is_set {
		return configuration{}, errors.New("secret_token env var not set")
	}

	return configuration{
		base_url:     base_url,
		secret_token: secret_token,
	}, nil

}
