package main

import (
	"net/http"

	"github.com/wagnojunior/booking/internal/config"
	"github.com/wagnojunior/booking/internal/handlers"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
)

// Instead of routing every page of the web application in the <main.go> files, it is a good practice to do it in separate file.
// For reference, the routing was done in the following manner:
// http.HandleFunc("/", handlers.Repo.Home)
// http.HandleFunc("/about", handlers.Repo.About)
func routes(app *config.AppConfig) http.Handler {
	// A http handlers is often times called a mux or a multiplexor

	// Create a new mux
	mux := chi.NewRouter()

	// Middleware allows you to process a web request as it comes and perform an action
	mux.Use(middleware.Recoverer) // Gracefully absorb panics and prints the stack trace
	mux.Use(NoSurf)               // Our own middleware which was created in <middleware.go> using the third-party package <nosurf>
	mux.Use(SessionLoad)

	// Get the http requests
	mux.Get("/", handlers.Repo.Home)
	mux.Get("/about", handlers.Repo.About)
	mux.Get("/panda-suite", handlers.Repo.PandaSuite)
	mux.Get("/bamboo-dorm", handlers.Repo.BambooDorm)
	mux.Get("/search-availability", handlers.Repo.SearchAvailability)
	mux.Get("/contact", handlers.Repo.Contact)
	mux.Get("/make-reservation", handlers.Repo.MakeReservation)
	mux.Get("/reservation-summary", handlers.Repo.ReservationSummary)

	// Post the http requests
	mux.Post("/search-availability", handlers.Repo.PostSearchAvailability) // Catch requests that POST to this url and send it to the specified handler
	mux.Post("/search-availability-json", handlers.Repo.AvailabilityJSON)
	mux.Post("/make-reservation", handlers.Repo.PostMakeReservation)

	// Creates a file server from which static files are retrieved
	fileServer := http.FileServer(http.Dir("./static/"))
	mux.Handle("/static/*", http.StripPrefix("/static", fileServer))

	return mux
}
