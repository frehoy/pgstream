package main

import (
	"fmt"
	"os"
	"testing"
)

func TestConfig_from_env_happy(t *testing.T) {
	// Params
	wanted_url := "https://foo.bar/"
	wanted_token := "verysecret"

	// Setup
	os.Setenv("base_url", wanted_url)
	os.Setenv("secret_token", wanted_token)
	defer os.Unsetenv("base_url")
	defer os.Unsetenv("secret_token")

	config, err := config_from_env()
	if err != nil {
		t.Fatalf("Got error from config_from_env()")
	}
	got_url := config.base_url.String()
	if got_url != wanted_url {
		t.Fatalf(fmt.Sprintf("got_url: %s, wanted: %s", got_url, wanted_url))

	}

	got_token := config.secret_token
	if got_token != wanted_token {
		t.Fatalf(fmt.Sprintf("got_token: %s, wanted: %s", got_token, wanted_token))
	}
}

func TestConfig_from_env_required_variables(t *testing.T) {

	defer os.Unsetenv("base_url")
	defer os.Unsetenv("secret_token")

	// Set no env vars, should error
	_, err := config_from_env()
	if err == nil {
		t.Fatalf("config_from_env() should error when missing all env vars")
	}

	// Set one of two required env vars, should error
	os.Setenv("secret_token", "verysercret")
	_, err = config_from_env()
	if err == nil {
		t.Fatalf("config_from_env() should error when missing one env var")
	}

	// Set a broken url that can't be parsed, should error
	os.Setenv("base_url", ":://::http://foo@@@")
	_, err = config_from_env()
	if err == nil {
		t.Fatalf("config_from_env() should error when base_url not parseable")
	}

	// Set a good url that can be parsed, should not error
	os.Setenv("base_url", "http://example.com")
	_, err = config_from_env()
	if err != nil {
		t.Fatalf("config_from_env() should not error when base_url is parseable")
	}
}
