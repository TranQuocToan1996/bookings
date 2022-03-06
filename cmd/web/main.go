package main

import (
	"encoding/gob"
	"flag"
	"fmt"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/TranQuocToan1996/bookings/internal/config"
	"github.com/TranQuocToan1996/bookings/internal/driver"
	"github.com/TranQuocToan1996/bookings/internal/handlers"
	"github.com/TranQuocToan1996/bookings/internal/helpers"
	"github.com/TranQuocToan1996/bookings/internal/models"
	"github.com/TranQuocToan1996/bookings/internal/render"
	"github.com/alexedwards/scs/v2"
)

const portNumber = ":8080"

var app config.AppConfig
var session *scs.SessionManager
var infoLog *log.Logger
var errorLog *log.Logger

// Main application func
func main() {

	db, err := run()
	if err != nil {
		log.Fatal(err)
	}
	defer db.SQL.Close()
	defer close(app.MailChan)

	// Listen continuous for an email
	infoLog.Println("Starting mail listener!")
	listenForMail()

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

/*  Old code: Alternative in routes.go
	// Send a request -> process request -> send back a response
   	http.HandleFunc("/", handlers.Repo.Home)
   	http.HandleFunc("/about", handlers.Repo.About)
	http.ListenAndServe(portNumber, nil)
*/

func run() (*driver.DB, error) {
	// https://stackoverflow.com/questions/47071276/decode-gob-output-without-knowing-concrete-types
	// Tell application about things (Premitive types) we need store in session
	// In other words, we register Revervation to save in session
	gob.Register(models.Reservation{})
	gob.Register(models.User{})
	gob.Register(models.Room{})
	gob.Register(models.Restriction{})
	gob.Register(map[string]int{})

	// read flags from terminal
	inProduction := flag.Bool("production", true, "Application is in production")
	useCache := flag.Bool("cache", true, "Save loading template in cache for reduce loading times")
	dbName := flag.String("dbname", "", "Database name")
	dbHost := flag.String("dbhost", "localhost", "Database host")
	dbUser := flag.String("dbuser", "", "Database user")
	dbPass := flag.String("dbpass", "", "Database password")
	dbPort := flag.String("dbport", "5432", "Database port")
	dbSSL := flag.String("dbssl", "disable", "Database SSL setting (disable, prefer, require)")

	// Parse the flags
	flag.Parse()
	if *dbName == "" || *dbUser == "" {
		log.Println("Missing require flags")
		os.Exit(1)
	}

	// Production
	app.InProduction = *inProduction
	app.UseCache = *useCache

	// Create mail channel
	mailChan := make(chan models.MailData) // We not going to close channel here it will close after run() finish
	app.MailChan = mailChan

	// Declare logs for appconfig
	infoLog = log.New(os.Stdout, "INFO\t", log.Ldate|log.Ltime)
	app.InfoLog = infoLog
	errorLog = log.New(os.Stdout, "ERROR\t", log.Ldate|log.Ltime|log.Lshortfile)
	app.ErrorLog = errorLog
``
	session = scs.New()
	session.Lifetime = 24 * time.Hour
	// Keep session even after close window/browser
	session.Cookie.Persist = true
	// allows you to declare if your cookie should be restricted to a first-party or same-site context
	session.Cookie.SameSite = http.SameSiteLaxMode
	// HTTPS
	session.Cookie.Secure = app.InProduction
	app.Session = session

	// Connect to database
	connectionString := fmt.Sprintf("host=%s port=%s dbname=%s user=%s password=%s sslmode=%s", *dbHost, *dbPort, *dbName, *dbUser, *dbPass, *dbSSL)
	db, err := driver.ConnectSQL(connectionString)
	if err != nil {
		log.Fatal("Can't connect to database! Exiting...")
	}
	log.Println("Connected to database")

	// Create template cache (map data structure of Golang)
	tc, err := render.CreateTemplateCache()
	if err != nil {
		// If we can't get template cache, we can't show any pages
		log.Fatal("Can't create template cache: ", err)
		return nil, err
	}
	app.TemplateCache = tc

	repo := handlers.NewRepo(&app, db)
	// Pass new repo to handler
	handlers.NewHandlers(repo)
	render.NewRenderer(&app)
	helpers.NewHelpers(&app)

	return db, nil
}
