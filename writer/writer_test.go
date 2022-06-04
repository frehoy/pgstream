package main

import (
	"os"
	"testing"
)

func TestN_threads(t *testing.T) {
	os.Setenv("N_THREADS", "-1")
	_, err := n_threads()
	if err == nil {
		t.Fatalf("n_threads() should error on negative thread count")
	}

	os.Setenv("N_THREADS", "0")
	_, err = n_threads()
	if err == nil {
		t.Fatalf("n_threads() should error on 0 thread count")
	}

	os.Setenv("N_THREADS", "65536")
	_, err = n_threads()
	if err == nil {
		// 65535 is max of uint16
		t.Fatalf("n_threads() should error on thread count larger than 65535")
	}
}
