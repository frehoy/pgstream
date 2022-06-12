package main

import (
	"errors"
	"fmt"
	"net/url"
	"os"
	"strconv"
)

const default_parallelism = 4
const duration = 60 * 60

func getEnvVar(varname string) (value string, err error) {
	value, is_set := os.LookupEnv(varname)
	if !is_set {
		return "", errors.New(fmt.Sprintf("%s env var not set.", varname))
	}
	return value, nil
}

func url_from_env() (*url.URL, error) {
	env_url, err := getEnvVar("WRITE_ENDPOINT")
	if err != nil {
		return nil, err
	}

	parsed_url, err := url.Parse(env_url)
	if err != nil {
		return nil, err
	}

	return parsed_url, nil
}

func n_threads() (uint16, error) {
	value, is_set := os.LookupEnv("N_THREADS")
	if !is_set {
		return default_parallelism, nil
	}

	value_int, err := strconv.ParseUint(value, 10, 16)
	if err != nil {
		return 0, err
	}

	if value_int == 0 {
		return 0, errors.New("Invalid number of threads: 0")
	}

	return uint16(value_int), nil
}

func get_settings() (appSettings, error) {
	token, err := getEnvVar("TOKEN")
	if err != nil {
		return appSettings{}, err
	}

	url_from_env, err := url_from_env()
	if err != nil {
		return appSettings{}, err
	}

	n_threads, err := n_threads()
	if err != nil {
		return appSettings{}, err
	}

	return appSettings{
		url:       url_from_env,
		token:     token,
		n_threads: n_threads,
	}, nil
}
