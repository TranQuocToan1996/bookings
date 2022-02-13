package render

import (
	"bytes"
	"errors"
	"fmt"
	"html/template"
	"log"
	"net/http"
	"path/filepath"

	"github.com/TranQuocToan1996/bookings/internal/config"
	"github.com/TranQuocToan1996/bookings/internal/models"

	"github.com/justinas/nosurf"
)

// Create func and pass to template
var functions = template.FuncMap{}

var app *config.AppConfig

var pathToTemplates = "./templates"

func AddDefaultData(td *models.TemplateData, r *http.Request) *models.TemplateData {
	// taking message from session to user, after that delete that message from session
	td.Flash = app.Session.PopString(r.Context(), "flash")
	td.Error = app.Session.PopString(r.Context(), "error")
	td.Warning = app.Session.PopString(r.Context(), "warning")

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

/* Old code - Load the template file from specified directory

t, _ := template.ParseFiles("./templates/" + html)
// Write to writter
err = t.Execute(w, nil)
if err != nil {
	fmt.Println("Error parsing template: ", err)
	return
}

Old code end*/

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
		fmt.Println("Getting name of page:", page)

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
