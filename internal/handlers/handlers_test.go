package handlers

import (
	"net/http"
	"net/http/httptest"
	"net/url"
	"testing"
)

type postData struct {
	key   string
	value string
}

// Test data is the slice of struct
var theTests = []struct {
	testName         string
	url              string
	requestMethod    string
	parameters       []postData
	expectStatusCode int
}{
	// GET method
	{"home", "/", "GET", []postData{}, http.StatusOK},
	{"about", "/about", "GET", []postData{}, http.StatusOK},
	{"generalsQuater", "/generals-quarters", "GET", []postData{}, http.StatusOK},
	{"majorsSuite", "/majors-suite", "GET", []postData{}, http.StatusOK},
	{"seachAvailability", "/search-availability", "GET", []postData{}, http.StatusOK},
	{"contact", "/contact", "GET", []postData{}, http.StatusOK},
	{"makeReservation", "/make-reservation", "GET", []postData{}, http.StatusOK},

	// POST method
	{"postSearchAvailability", "/search-availability", "POST", []postData{
		{key: "start", value: "2020-01-01"},
		{key: "end", value: "2020-01-01"},
	}, http.StatusOK},
	{"postSearchAvailability-json", "/search-availability-json", "POST", []postData{
		{key: "start", value: "2020-01-01"},
		{key: "end", value: "2020-01-01"},
	}, http.StatusOK},
	{"makeReservation", "/make-reservation", "POST", []postData{
		{key: "first_name", value: "Toan"},
		{key: "last_name", value: "Tran"},
		{key: "email", value: "tranquoctoan@example.com"},
		{key: "phone", value: "0989xxxxxx"},
	}, http.StatusOK},
}

func TestHanlers(t *testing.T) {
	routes := getRoutes()

	// Create a server for testing, close it when TestHanlers finish
	testServer := httptest.NewTLSServer(routes)
	defer testServer.Close()

	for _, e := range theTests {
		if e.requestMethod == "GET" {
			// Create cliend brower and send get method
			resp, err := testServer.Client().Get(testServer.URL + e.url)
			if err != nil {
				t.Log(err)
				t.Fatal(err)
			}
			if resp.StatusCode != e.expectStatusCode {
				t.Errorf("For %s expecting %d but got %d", e.testName, e.expectStatusCode, resp.StatusCode)
			}

			// Test for Post
		} else if e.requestMethod == "POST" {
			// Hold information as POST request
			value := url.Values{}
			for _, x := range e.parameters {
				value.Add(x.key, x.value)
			}
			// Create cliend brower and send post method with data
			resp, err := testServer.Client().PostForm(testServer.URL+e.url, value)
			if err != nil {
				t.Log(err)
				t.Fatal(err)
			}
			if resp.StatusCode != e.expectStatusCode {
				t.Errorf("For %s expecting %d but got %d", e.testName, e.expectStatusCode, resp.StatusCode)
			}

		}
	}
}
