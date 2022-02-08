package main

import (
	"fmt"
	"net/http"
	"testing"
)

func TestNoSurf(t *testing.T) {

	var myH myHandler

	handlers := NoSurf(&myH)

	switch v := handlers.(type) {
	case http.Handler:
		// Do nothing
	default:
		t.Error(fmt.Sprintf("Type is not http.Handler, but is %T", v))
	}
}

func TestSessionLoad(t *testing.T) {

	var myH myHandler

	handlers := SessionLoad(&myH)

	switch v := handlers.(type) {
	case http.Handler:
		// Do nothing
	default:
		t.Error(fmt.Sprintf("Type is not http.Handler, but is %T", v))
	}
}
