package main

import (
	"net/http"

	"github.com/justinas/nosurf"
)

// In case a middleware does not come out of the box from the router, it is necessary to build our own middleware

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
