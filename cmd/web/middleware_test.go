package main

import (
	"net/http"
	"testing"
)

// In order to test <NoSurf()> we need a <http.Handler> as an argument
func TestNoSurf(t *testing.T) {
	// Creates a variable of type <http.Handler>
	var myH myHandler

	h := NoSurf(&myH)

	switch v := h.(type) {
	case http.Handler:
		// do nothing
	default:
		t.Errorf("Type %t is not <http.Handler>", v)
	}
}

// In order to test <SessionLoad()> we need a <http.Handler> as an argument
func TestSessionLoad(t *testing.T) {
	// Creates a variable of type <http.Handler>
	var myH myHandler

	h := SessionLoad(&myH)

	switch v := h.(type) {
	case http.Handler:
		// do nothing
	default:
		t.Errorf("Type %t is not <http.Handler>", v)
	}
}
