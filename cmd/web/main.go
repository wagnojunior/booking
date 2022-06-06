package main

import (
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/wagnojunior/booking/pkg/config"
	"github.com/wagnojunior/booking/pkg/handlers"
	"github.com/wagnojunior/booking/pkg/render"

	"github.com/alexedwards/scs/v2"
)

// Constants
const portNumber = ":8080" // port number

// Package-level variables
var app config.AppConfig
var session *scs.SessionManager

func main() {

	// Set the
	app.InProduction = false

	// Set the configuration of sessions
	session = scs.New()
	session.Lifetime = 24 * time.Hour
	session.Cookie.Persist = true
	session.Cookie.SameSite = http.SameSiteLaxMode
	session.Cookie.Secure = app.InProduction

	app.Session = session

	tc, err := render.CreateTemplateCache()
	if err != nil {
		log.Fatal("cannot create template cache")
	}

	app.TemplateCache = tc
	app.UseCache = false

	repo := handlers.NewRepo(&app)
	handlers.NewHandlers(repo)

	render.NewTemplates(&app)

	fmt.Println(fmt.Sprintf("Staring application on port %s", portNumber))

	srv := &http.Server{
		Addr:    portNumber,
		Handler: routes(&app),
	}

	err = srv.ListenAndServe()
	log.Fatal(err)
}
