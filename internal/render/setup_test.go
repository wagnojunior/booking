package render

import (
	"encoding/gob"
	"log"
	"net/http"
	"os"
	"testing"
	"time"

	"github.com/alexedwards/scs/v2"
	"github.com/wagnojunior/booking/internal/config"
	"github.com/wagnojunior/booking/internal/models"
)

var session *scs.SessionManager
var testApp config.AppConfig

// This function is called before any of the tests are run, executes the body, and
// then call the tests themselves
func TestMain(m *testing.M) {
	// Things that will be put in sessions
	gob.Register(models.Reservation{})

	// Set the in development mode
	testApp.InProduction = false

	// Creates the infoLog. Write to the standard output (terminal),
	// prefixed by the tag INFO, and flagged by the date and time
	infoLog := log.New(os.Stdout, "INFO\t", log.Ldate|log.Ltime)
	testApp.InfoLog = infoLog

	errorLog := log.New(os.Stdout, "ERROR\t", log.Ldate|log.Ltime|log.Lshortfile)
	testApp.ErrorLog = errorLog

	// Set the configuration of sessions
	session = scs.New()                            // Creates a new session
	session.Lifetime = 24 * time.Hour              // Defines for how long the session will persist
	session.Cookie.Persist = true                  // Cookies will persist after the browser is closed by the end-user
	session.Cookie.SameSite = http.SameSiteLaxMode //
	session.Cookie.Secure = false

	// Set the <Session> field in the <AppConfig>, thus exposing this variable to all packages that import <config.go>
	testApp.Session = session

	// app is defined in rander.go
	app = &testApp

	os.Exit(m.Run())
}

type myWrite struct{}

func (tw *myWrite) Header() http.Header {
	var h http.Header
	return h
}

func (tw *myWrite) WriteHeader(i int) {

}

func (tw *myWrite) Write(b []byte) (int, error) {
	length := len(b)
	return length, nil
}
