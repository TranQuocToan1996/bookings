package handlers

import (
	"context"
	"encoding/gob"
	"fmt"
	"html/template"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"testing"
	"time"

	"github.com/TranQuocToan1996/bookings/internal/config"
	"github.com/TranQuocToan1996/bookings/internal/models"
	"github.com/TranQuocToan1996/bookings/internal/render"
	"github.com/TranQuocToan1996/bookings/internal/repository/dbrepo"
	"github.com/alexedwards/scs/v2"
	"github.com/go-chi/chi"
	"github.com/go-chi/chi/middleware"
	"github.com/justinas/nosurf"
)

var app config.AppConfig
var session *scs.SessionManager
var pathToTemplates = "./../../templates"
var functions = template.FuncMap{}

// NewRepo creates a new Repository
func NewTestRepo(a *config.AppConfig) *Repository {
	return &Repository{
		App: a,
		DB:  dbrepo.NewTestingRepo(a),
	}
}

func TestMain(m *testing.M) {
	// https://stackoverflow.com/questions/47071276/decode-gob-output-without-knowing-concrete-types
	// Tell application about things (Premitive types) we need store in session
	gob.Register(models.Reservation{})

	// Production
	app.InProduction = false

	infoLog := log.New(os.Stdout, "INFO\t", log.Ldate|log.Ltime)
	app.InfoLog = infoLog
	errorLog := log.New(os.Stdout, "ERROR\t", log.Ldate|log.Ltime|log.Lshortfile)
	app.ErrorLog = errorLog

	session = scs.New()
	session.Lifetime = 24 * time.Hour
	// Keep session even after close window/browser
	session.Cookie.Persist = true
	// allows you to declare if your cookie should be restricted to a first-party or same-site context
	session.Cookie.SameSite = http.SameSiteLaxMode
	// HTTPS
	session.Cookie.Secure = app.InProduction

	app.Session = session

	// Channel listen for email (Use for testing only)
	mailChan := make(chan models.MailData)
	app.MailChan = mailChan
	defer close(mailChan)
	listenForMail()

	// Create template cache (map data structure of Golang)
	tc, err := CreateTestTemplateCache()
	if err != nil {
		// If we can't get template cache, we can't show any pages
		log.Fatal("Can't create template cache: ", err)
	}
	app.TemplateCache = tc
	app.UseCache = true

	repo := NewTestRepo(&app)
	// Pass new repo to handler
	NewHandlers(repo)

	render.NewRenderer(&app)

	// start to running test, after that exit program
	os.Exit(m.Run())
}

func getRoutes() http.Handler {

	// Using chi
	mux := chi.NewRouter()
	// Recoverer middleware recover panic, print log panic and return page with code 500
	mux.Use(middleware.Recoverer)
	// nosurf is middleware that prevent Cross-Site Request Forgery (CSRF) attacks from all POST request (It should accept post request with CSRF-token)
	// mux.Use(NoSurf) -> nosurf are already tested in middleware_test.go

	// This my own middleware add to Handlers some features to test system
	mux.Use(WriteToConsole)
	// SessionLoad loads and saves the session on every request
	mux.Use(SessionLoad)

	// Handlers get request
	mux.Get("/", Repo.Home)
	mux.Get("/about", Repo.About)
	mux.Get("/generals-quarters", Repo.Generals)
	mux.Get("/majors-suite", Repo.Majors)
	mux.Get("/search-availability", Repo.Availability)
	mux.Get("/contact", Repo.Contact)
	mux.Get("/make-reservation", Repo.Reservation)
	mux.Get("/reservation-summary", Repo.ReservationSummary)

	// Handlers Post request
	mux.Post("/search-availability", Repo.PostAvailability)
	mux.Post("/search-availability-json", Repo.AvailabilityJSON)
	mux.Post("/make-reservation", Repo.PostReservation)

	// FileServer is the place to get static files
	fileServer := http.FileServer(http.Dir("./static/"))
	mux.Handle("/static/*", http.StripPrefix("/static", fileServer))

	return mux
}

// NoSurf adds CSRF protestion for all POST requests
func NoSurf(next http.Handler) http.Handler {

	csrfHandler := nosurf.New(next)
	// Using cookies to make sure the csrfToken available on a per page basic
	csrfHandler.SetBaseCookie(http.Cookie{
		HttpOnly: true,             // Decrease chance  risk of client side script accessing the protected cookie (stop Js in client side)
		Path:     "/",              // Entire site
		Secure:   app.InProduction, // HTTPs
		SameSite: http.SameSiteLaxMode,
	})

	return csrfHandler
}

// LoadAndSave provides middleware which automatically loads and saves session
// data for the current request, and communicates the session token to and from
// the client in a cookie.
func SessionLoad(next http.Handler) http.Handler {
	return session.LoadAndSave(next)
}

// WriteToConsole log some text to terminal when client load a page
func WriteToConsole(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// fmt.Println("Hit the page: " + r.Host + r.URL.Path)
		next.ServeHTTP(w, r)
	})
}

// CreateTestTemplateCache return a map of template caches
func CreateTestTemplateCache() (map[string]*template.Template, error) {
	// Create a map for store all templates
	myCache := map[string]*template.Template{}

	// Go to find all html template name
	pages, err := filepath.Glob(fmt.Sprintf("%s/*.page.html", pathToTemplates))
	if err != nil {
		return myCache, err
	}

	// return full path of template
	for _, page := range pages {
		// fmt.Println(filepath.Base("/foo/bar/baz.js")) return baz.js
		name := filepath.Base(page)

		// template set will has some functions
		// Must call .Funcs before Parse template
		ts, err := template.New(name).Funcs(functions).ParseFiles(page)
		if err != nil {

			return myCache, err
		}

		// Add all layout to []string
		matches, err := filepath.Glob(fmt.Sprintf("%s/*.layout.html", pathToTemplates))
		if err != nil {

			return myCache, err
		}
		if len(matches) > 0 {
			ts, err = ts.ParseGlob(fmt.Sprintf("%s/*.layout.html", pathToTemplates))
			if err != nil {

				return myCache, err
			}
		}

		// add template set to the Cache map
		myCache[name] = ts
	}

	return myCache, nil
}

// Get context include session data
func getCtx(r *http.Request) context.Context {
	// ctx is context contains session data
	ctx, err := session.Load(r.Context(), r.Header.Get("X-Session"))
	if err != nil {
		log.Println(err)
	}
	return ctx
}

func listenForMail() {
	go func() {
		for {
			// Get email but do nothing to it (To make app.MailChan waiting for receiver forever)
			_ = <-app.MailChan
		}
	}()
}
