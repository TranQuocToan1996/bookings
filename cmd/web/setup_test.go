package main

import (
	"net/http"
	"os"
	"testing"
)

// M is a type passed to a TestMain function to run the actual tests.

func TestMain(m *testing.M) {

	/* The above code of TestMain will run first, after that m.Run run all test (_test.go) and passing exit code to
	os.Exit Finally */
	os.Exit(m.Run())
}

type myHandler struct{}

func (mh *myHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {

}
