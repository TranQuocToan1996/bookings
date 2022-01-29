package main

import (
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/TranQuocToan1996/bookings/pkg/config"
	"github.com/TranQuocToan1996/bookings/pkg/handlers"
	"github.com/TranQuocToan1996/bookings/pkg/render"
	"github.com/alexedwards/scs/v2"
)

const portNumber = ":8080"

var app config.AppConfig
var session *scs.SessionManager

// Main application func
func main() {

	// Production
	app.InProduction = false

	session = scs.New()
	session.Lifetime = 24 * time.Hour
	// Keep session even after close window/browser
	session.Cookie.Persist = true
	// allows you to declare if your cookie should be restricted to a first-party or same-site context
	session.Cookie.SameSite = http.SameSiteLaxMode
	// HTTPS
	session.Cookie.Secure = app.InProduction

	app.Session = session

	// Create template cache (map data structure of Golang)
	tc, err := render.CreateTemplateCache()
	if err != nil {
		// If we can't get template cache, we can't show any pages
		log.Fatal("Can't create template cache: ", err)
	}
	app.TemplateCache = tc
	app.UseCache = false

	repo := handlers.NewRepo(&app)
	// Pass new repo to handler
	handlers.NewHandlers(repo)

	render.NewTemplates(&app)

	/*  Alternative in routes.go
		// Send a request -> process request -> send back a response
	   	http.HandleFunc("/", handlers.Repo.Home)
	   	http.HandleFunc("/about", handlers.Repo.About)
		http.ListenAndServe(portNumber, nil)
	*/

	fmt.Println("Starting application on port:", portNumber)
	// Start the server
	srv := &http.Server{
		Addr:    portNumber,
		Handler: routes(&app),
	}
	// Start server. After server close or shutdown, return err
	err = srv.ListenAndServe()
	log.Fatal(err)
}
