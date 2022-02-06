package main

import (
	"net/http"

	"github.com/TranQuocToan1996/bookings/internal/config"
	"github.com/TranQuocToan1996/bookings/internal/handlers"
	"github.com/go-chi/chi"
	"github.com/go-chi/chi/middleware"
)

func routes(app *config.AppConfig) http.Handler {
	/*     // Create multiplexed Using pat Package
	mux := pat.New()
	mux.Get("/", http.HandlerFunc(handlers.Repo.Home))
	mux.Get("/about", http.HandlerFunc(handlers.Repo.About))
	return mux */

	// Using chi
	mux := chi.NewRouter()
	// Recoverer middleware recover panic, print log panic and return page with code 500
	mux.Use(middleware.Recoverer)
	// nosurf is middleware that prevent Cross-Site Request Forgery (CSRF) attacks from all POST request (It should accept post request with CSRF-token)
	mux.Use(NoSurf)

	// This my own middleware add to Handlers some features to test system
	mux.Use(WriteToConsole)
	// SessionLoad loads and saves the session on every request
	mux.Use(SessionLoad)

	// Handlers get request
	mux.Get("/", handlers.Repo.Home)
	mux.Get("/about", handlers.Repo.About)
	mux.Get("/generals-quarters", handlers.Repo.Generals)
	mux.Get("/majors-suite", handlers.Repo.Majors)
	mux.Get("/search-availability", handlers.Repo.Availability)
	mux.Get("/contact", handlers.Repo.Contact)
	mux.Get("/make-reservation", handlers.Repo.Reservation)
	mux.Get("/reservation-summary", handlers.Repo.ReservationSummary)

	// Handlers Post request
	mux.Post("/search-availability", handlers.Repo.PostAvailability)
	mux.Post("/search-availability-json", handlers.Repo.AvailabilityJSON)
	mux.Post("/make-reservation", handlers.Repo.PostReservation)

	// FileServer is the place to get static files
	fileServer := http.FileServer(http.Dir("./static/"))
	mux.Handle("/static/*", http.StripPrefix("/static", fileServer))

	return mux
}
