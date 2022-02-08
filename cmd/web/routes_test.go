package main

import (
	"fmt"
	"testing"

	"github.com/TranQuocToan1996/bookings/internal/config"
	"github.com/go-chi/chi"
)

func TestRoutes(t *testing.T) {
	var app config.AppConfig

	mux := routes(&app)
	switch v := mux.(type) {
	case *chi.Mux:
		// do nothing: Test pass
	default:
		t.Error(fmt.Sprintf("Type is not *chi.Mux, type is %T", v))
	}
}
