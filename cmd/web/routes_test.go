package main

import (
	"testing"

	"github.com/go-chi/chi/v5"
	"github.com/wagnojunior/booking/internal/config"
)

func TestRoutes(t *testing.T) {
	var app config.AppConfig

	mux := routes(&app)

	switch v := mux.(type) {
	case *chi.Mux:
		// do nothing, test passed
	default:
		t.Errorf("Type %t is not <*chi.Mux>", v)
	}
}
