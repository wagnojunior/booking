// This file <setup_test.go> runs before any other test.
// It must have a function <TestMain (m *testing.M)>
package main

import (
	"net/http"
	"os"
	"testing"
)

func TestMain(m *testing.M) {
	// Place to setup the test environment

	// <m.Run()> runs all other tests, then exit
	os.Exit(m.Run())
}

// This is the handler used to test <NoSurf()>.
// It must implement the methods in <http.Handler> interface
// in order to be of type <http.Handler>
type myHandler struct{}

func (mh *myHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {

}
