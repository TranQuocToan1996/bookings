package render

import (
	"encoding/gob"
	"log"
	"net/http"
	"os"
	"testing"
	"time"

	"github.com/TranQuocToan1996/bookings/internal/config"
	"github.com/TranQuocToan1996/bookings/internal/models"
	"github.com/alexedwards/scs/v2"
)

var session *scs.SessionManager

// config data for testing render.go
var testApp config.AppConfig

func TestMain(m *testing.M) {

	// Set up session for testing
	gob.Register(models.Revervation{})
	testApp.InProduction = false
	infoLog := log.New(os.Stdout, "INFO\t", log.Ldate|log.Ltime)
	testApp.InfoLog = infoLog
	errorLog := log.New(os.Stdout, "ERROR\t", log.Ldate|log.Ltime|log.Lshortfile)
	testApp.ErrorLog = errorLog
	session = scs.New()
	session.Lifetime = 24 * time.Hour
	session.Cookie.Persist = true
	session.Cookie.SameSite = http.SameSiteLaxMode
	session.Cookie.Secure = false
	testApp.Session = session
	app = &testApp

	// Run all tests before exit
	os.Exit(m.Run())
}

type myWriter struct{}

// http.Header is key-value pairs in an HTTP header
func (w *myWriter) Header() http.Header {
	var h http.Header
	return h
}

func (w *myWriter) WriteHeader(i int) {}

func (w *myWriter) Write(b []byte) (int, error) {
	length := len(b)
	return length, nil
}
