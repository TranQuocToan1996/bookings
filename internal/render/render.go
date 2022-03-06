package render

import (
	"bytes"
	"errors"
	"fmt"
	"html/template"
	"log"
	"net/http"
	"path/filepath"
	"time"

	"github.com/TranQuocToan1996/bookings/internal/config"
	"github.com/TranQuocToan1996/bookings/internal/models"

	"github.com/justinas/nosurf"
)

// Create func and pass to template for golang template
var functions = template.FuncMap{
	"humanDate":  HumanDate,
	"formatDate": FormatDate,
	"iterate":    Iterate,
	"add":        Add,
}

var app *config.AppConfig

var pathToTemplates = "./templates"

// Add returns a slice of int starting at 1 to count
func Add(a, b int) int {
	return a + b
}

// Iterate returns a slice of int starting at 1 to count
func Iterate(count int) []int {
	var i int
	var items []int
	for i = 0; i < count; i++ {
		items = append(items, i)
	}
	return items
}

// HumanDate formats time.Time into yyyy-mm-dd
func HumanDate(t time.Time) string {
	return t.Format("2006-01-02")
}

// FormatDate formats time.Time into yyyy-mm-dd
func FormatDate(t time.Time, format string) string {
	return t.Format(format)
}

// AddDefaultData adds data for all templates
func AddDefaultData(td *models.TemplateData, r *http.Request) *models.TemplateData {
	// taking message from session to user, after that delete that message from session
	td.Flash = app.Session.PopString(r.Context(), "flash")
	td.Error = app.Session.PopString(r.Context(), "error")
	td.Warning = app.Session.PopString(r.Context(), "warning")

	// If user already login return IsAuthenticate=1 when redering templates
	if app.Session.Exists(r.Context(), "user_id") {
		td.IsAuthenticate = 1
	}

	td.CSRFToken = nosurf.Token(r)
	return td
}

// NewRenderer sets the config for the template package
func NewRenderer(a *config.AppConfig) {
	app = a
}

// Template renders templates using html/template
func Template(w http.ResponseWriter, r *http.Request, html string, td *models.TemplateData) error {

	// Sometime in development, Rebuild the template on every requests
	var tc map[string]*template.Template
	if app.UseCache {
		// Get the template cache from config.go
		tc = app.TemplateCache
	} else {
		// Rebuild the template from templates directory
		tc, _ = CreateTemplateCache()
	}

	t, ok := tc[html]
	if !ok {
		log.Println("Could not get template from cache")
		return errors.New("can't get template from cache")
	}

	// Turn template cache (In memory) into some bytes
	buf := new(bytes.Buffer)

	td = AddDefaultData(td, r)

	_ = t.Execute(buf, td)

	_, err := buf.WriteTo(w)
	if err != nil {
		fmt.Println("Error writing template to browser", err)
		return err
	}

	return nil

}

// CreateTemplateCache return a map of template caches
func CreateTemplateCache() (map[string]*template.Template, error) {
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
		// fmt.Println("Getting name of page:", page)

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
