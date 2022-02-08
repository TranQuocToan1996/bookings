package render

import (
	"net/http"
	"testing"

	"github.com/TranQuocToan1996/bookings/internal/models"
)

func TestAddDefaultData(t *testing.T) {
	var name string = "TestAddDefaultData"
	var td models.TemplateData

	r, err := getSession()
	if err != nil {
		t.Errorf("%s: %s", name, err)
	}

	session.Put(r.Context(), "flash", "123")
	result := AddDefaultData(&td, r)
	if result == nil {
		t.Errorf("%s: Can't add default data for template", name)
	}
	if result.Flash != "123" {
		t.Errorf("%s: flash value of 123 not found in session", name)
	}

}

func TestNewTemplate(t *testing.T) {
	NewTemplates(app)

}

func TestRenderTemplate(t *testing.T) {
	testName := "TestRenderTemplate"
	pathToTemplates = "./../../templates"

	// Create cache for templates
	tc, err := CreateTemplateCache()
	if err != nil {
		t.Errorf("%s: %s", testName, err)
	}
	app.TemplateCache = tc

	// Create request
	r, err := getSession()
	if err != nil {
		t.Errorf("%s: %s", testName, err)
	}

	// Create repsonse writer
	var ww myWriter

	// Check exist page
	err = RenderTemplate(&ww, r, "home.page.html", &models.TemplateData{})
	if err != nil {
		t.Errorf("%s: %s", testName, "Can't writing template to browser")
	}

	// Check non-exist page
	err = RenderTemplate(&ww, r, "non-exist.page.html", &models.TemplateData{})
	if err == nil {
		t.Errorf("%s: %s", testName, "Render a template that not exist")
	}

}

func TestCreateTemplateCache(t *testing.T) {
	testName := "TestCreateTemplateCache"
	pathToTemplates = "./../../templates"

	_, err := CreateTemplateCache()
	if err != nil {
		t.Errorf("%s: %s", testName, err)
	}
}

func getSession() (*http.Request, error) {
	r, err := http.NewRequest("GET", "/some-url", nil)
	if err != nil {
		return nil, err
	}

	ctx := r.Context()
	// Put
	ctx, _ = session.Load(ctx, r.Header.Get("X-Session"))
	// Put context back to request
	r = r.WithContext(ctx)
	return r, nil
}
