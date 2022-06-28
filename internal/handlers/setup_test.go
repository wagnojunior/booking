package handlers

import (
	"encoding/gob"
	"fmt"
	"html/template"
	"log"
	"net/http"
	"path/filepath"
	"time"

	"github.com/alexedwards/scs/v2"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/justinas/nosurf"
	"github.com/wagnojunior/booking/internal/config"
	"github.com/wagnojunior/booking/internal/models"
	"github.com/wagnojunior/booking/internal/render"
)

var app config.AppConfig
var session *scs.SessionManager
var pathToTemplates = "./../../templates"
var functions = template.FuncMap{}

func getRoutes() http.Handler {
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
	tc, err := CreateTestTemplateCache()
	if err != nil {
		log.Fatal("cannot create template cache")
	}

	// Sets the <TemplateCache> field in the <AppConfig>
	app.TemplateCache = tc
	app.UseCache = true

	// Sets the variable <repo>, which points to the <AppConfig> app
	repo := NewRepo(&app)

	// Sends the local variable <repo> to <handlers.go> to initialize the variable <Repo> there
	NewHandlers(repo)

	// Initialized the variable <app> of type <*AppConfig> in <render.go>
	render.NewTemplates(&app)

	// Create a new mux
	mux := chi.NewRouter()

	// Middleware allows you to process a web request as it comes and perform an action
	mux.Use(middleware.Recoverer) // Gracefully absorb panics and prints the stack trace
	mux.Use(NoSurf)               // Our own middleware which was created in <middleware.go> using the third-party package <nosurf>
	mux.Use(SessionLoad)

	// Get the http requests
	mux.Get("/", Repo.Home)
	mux.Get("/about", Repo.About)
	mux.Get("/panda-suite", Repo.PandaSuite)
	mux.Get("/bamboo-dorm", Repo.BambooDorm)
	mux.Get("/search-availability", Repo.SearchAvailability)
	mux.Get("/contact", Repo.Contact)
	mux.Get("/make-reservation", Repo.MakeReservation)
	mux.Get("/reservation-summary", Repo.ReservationSummary)

	// Post the http requests
	mux.Post("/search-availability", Repo.PostSearchAvailability) // Catch requests that POST to this url and send it to the specified handler
	mux.Post("/search-availability-json", Repo.AvailabilityJSON)
	mux.Post("/make-reservation", Repo.PostMakeReservation)

	// Creates a file server from which static files are retrieved
	fileServer := http.FileServer(http.Dir("./static/"))
	mux.Handle("/static/*", http.StripPrefix("/static", fileServer))

	return mux
}

// NoSurf adds CSRF protection to all POST requests
func NoSurf(next http.Handler) http.Handler {
	csrfHandler := nosurf.New(next)

	// Set the base cookie
	csrfHandler.SetBaseCookie(http.Cookie{
		HttpOnly: true,
		Path:     "/", // Applies to the entire website
		Secure:   app.InProduction,
		SameSite: http.SameSiteLaxMode,
	})

	return csrfHandler
}

// SessionLoad loads and saves the session on every request.
// Sessions are persistent data about the user between page requests
func SessionLoad(next http.Handler) http.Handler {
	return session.LoadAndSave(next)
}

// CreateTemplateCache creates a template cache as a map
func CreateTestTemplateCache() (map[string]*template.Template, error) {

	// <myCache> maps a string to a pointer to <template.Template>
	myCache := map[string]*template.Template{}

	// Gets the file path of all files in the folder <templates> that end with <.page.tmpl>
	pages, err := filepath.Glob(fmt.Sprintf("%s/*.page.tmpl", pathToTemplates))
	if err != nil {
		return myCache, err
	}

	// Loop through all the pages
	for _, page := range pages {
		// Gets the file path base
		name := filepath.Base(page)

		// <New> allocates a new HTML template with the given <name>
		// <Funcs> adds the elements of the argument map to the template's function map
		// <ParseFiles> parses the named files and associates the resulting templates with t
		ts, err := template.New(name).Funcs(functions).ParseFiles(page)
		if err != nil {
			return myCache, err
		}

		// <Glob> returns the names of all files matching pattern or nil if there is no matching file
		matches, err := filepath.Glob(fmt.Sprintf("%s/*.layout.tmpl", pathToTemplates))
		if err != nil {
			return myCache, err
		}

		// If the length of <matches> is grater than zero, there are might be layouts associated with templates
		// In fact, this is the case for <about.page.tmpl> and <home.page.tmpl> which reference <base.layout.tmpl>
		if len(matches) > 0 {
			// ParseGlob parses the template definitions in the files identified by the pattern and associates the
			//resulting templates with t.
			ts, err = ts.ParseGlob(fmt.Sprintf("%s/*.layout.tmpl", pathToTemplates))
			if err != nil {
				return myCache, err
			}
		}

		// Add the parsed template to the template map <myCache>
		myCache[name] = ts
	}

	return myCache, nil
}
