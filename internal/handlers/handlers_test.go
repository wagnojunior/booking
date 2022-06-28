package handlers

import (
	"net/http"
	"net/http/httptest"
	"testing"
)

type postData struct {
	key   string
	value string
}

var theTest = []struct {
	name         string // name of the test
	url          string
	method       string // get os post
	params       []postData
	expectedCode int // 200 for OK, 404 page not found, 300 redirect
}{
	{"home", "/", "GET", []postData{}, http.StatusOK},                                   // first entry of the test
	{"about", "/about", "GET", []postData{}, http.StatusOK},                             // first entry of the test
	{"panda-suite", "/panda-suite", "GET", []postData{}, http.StatusOK},                 // first entry of the test
	{"bamboo-dorm", "/bamboo-dorm", "GET", []postData{}, http.StatusOK},                 // first entry of the test
	{"search-availability", "/search-availability", "GET", []postData{}, http.StatusOK}, // first entry of the test
	{"contact", "/contact", "GET", []postData{}, http.StatusOK},                         // first entry of the test
	{"make-reservation", "/make-reservation", "GET", []postData{}, http.StatusOK},       // first entry of the test
}

func TestHandler(t *testing.T) {
	routes := getRoutes()

	// test server
	ts := httptest.NewTLSServer(routes)
	defer ts.Close()

	for _, e := range theTest {
		if e.method == "GET" {
			resp, err := ts.Client().Get(ts.URL + e.url)
			if err != nil {
				t.Log(err)
				t.Fatal(err)
			}

			if resp.StatusCode != e.expectedCode {
				t.Errorf("For %s, expected %d but got %d", e.name, e.expectedCode, resp.StatusCode)
			}
		} else {

		}
	}
}
