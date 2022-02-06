package render

import (
	"bytes"
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

func AddDefaultData(td *models.TemplateData, r *http.Request) *models.TemplateData {
	td.CSRFToken = nosurf.Token(r)
	return td
}

// NewTemplates sets the config for the template package
func NewTemplates(a *config.AppConfig) {
	app = a
}

func RenderTemplate(w http.ResponseWriter, r *http.Request, html string, td *models.TemplateData) {

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
		log.Fatal("Could not get template from cache")
	}

	// Turn template cache (In memory) into some bytes
	buf := new(bytes.Buffer)

	td = AddDefaultData(td, r)

	_ = t.Execute(buf, td)

	_, err := buf.WriteTo(w)
	if err != nil {
		fmt.Println("Error writing template to browser", err)
	}

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
	pages, err := filepath.Glob("./templates/*.page.html")
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
		matches, err := filepath.Glob("./templates/*.layout.html")
		if err != nil {

			return myCache, err
		}
		if len(matches) > 0 {
			ts, err = ts.ParseGlob("./templates/*.layout.html")
			if err != nil {

				return myCache, err
			}
		}

		// add template set to the Cache map
		myCache[name] = ts
	}

	return myCache, nil
}
