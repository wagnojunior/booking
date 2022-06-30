package render

import (
	"encoding/gob"
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
