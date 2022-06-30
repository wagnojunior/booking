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
	{"post-search-availability", "/search-availability", "POST", []postData{
		{key: "start", value: "2022-01-01"},
		{key: "end", value: "2022-01-02"},
	}, http.StatusOK},
	{"post-search-availability-json", "/search-availability-json", "POST", []postData{
		{key: "start", value: "2022-01-01"},
		{key: "end", value: "2022-01-02"},
	}, http.StatusOK},
	{"make-reservation", "/make-reservation", "POST", []postData{
		{key: "first_name", value: "John"},
		{key: "last_name", value: "Smith"},
		{key: "email", value: "me@here.com"},
		{key: "phone", value: "555-555-5555"},
	}, http.StatusOK},
}

func TestHandler(t *testing.T) {
	routes := getRoutes()

	// test server
	ts := httptest.NewTLSServer(routes)
	defer ts.Close()

	// Loop thruogh all test cases defined above
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
		} else { // POST method
			// creates a variable that has the same format expected by the server
			values := url.Values{}
			// populate with the test entries
			for _, x := range e.params {
				values.Add(x.key, x.value)
			}

			resp, err := ts.Client().PostForm(ts.URL+e.url, values)
			if err != nil {
				t.Log(err)
				t.Fatal(err)
			}

			if resp.StatusCode != e.expectedCode {
				t.Errorf("For %s, expected %d but got %d", e.name, e.expectedCode, resp.StatusCode)
			}
		}
	}
}
