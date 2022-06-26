package main

import (
	"encoding/gob"
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/wagnojunior/booking/internal/config"
	"github.com/wagnojunior/booking/internal/handlers"
	"github.com/wagnojunior/booking/internal/models"
	"github.com/wagnojunior/booking/internal/render"

	"github.com/alexedwards/scs/v2"
)

// Constants
const portNumber = ":8080" // port number

// Package-level variables
var app config.AppConfig
var session *scs.SessionManager // Creates a variable <sessions>

func main() {
	// Things that will be put in sessions
	gob.Register(models.Reservation{})

	// Set the in development mode
	app.InProduction = false

	// Set the configuration of sessions
	session = scs.New()                            // Creates a new session
	session.Lifetime = 24 * time.Hour              // Defines for how long the session will persist
	session.Cookie.Persist = true                  // Cookies will persist after the browser is closed by the end-user
	session.Cookie.SameSite = http.SameSiteLaxMode //
	session.Cookie.Secure = app.InProduction

	// Set the <Session> field in the <AppConfig>, thus exposing this variable to all packages that import <config.go>
	app.Session = session

	// Creates the template cache
	tc, err := render.CreateTemplateCache()
	if err != nil {
		log.Fatal("cannot create template cache")
	}

	// Sets the <TemplateCache> field in the <AppConfig>
	app.TemplateCache = tc
	app.UseCache = false

	// Sets the variable <repo>, which points to the <AppConfig> app
	repo := handlers.NewRepo(&app)

	// Sends the local variable <repo> to <handlers.go> to initialize the variable <Repo> there
	handlers.NewHandlers(repo)

	// Initialized the variable <app> of type <*AppConfig> in <render.go>
	render.NewTemplates(&app)

	fmt.Println(fmt.Sprintf("Staring application on port %s", portNumber))

	// Initializes a server
	srv := &http.Server{
		Addr:    portNumber,
		Handler: routes(&app), // Instead of writing the handlers one by one for every webpage, pass the routes function
	}

	// Start the server
	err = srv.ListenAndServe()
	log.Fatal(err)
}
