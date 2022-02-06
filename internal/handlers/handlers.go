package handlers

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"

	"github.com/TranQuocToan1996/bookings/internal/config"
	"github.com/TranQuocToan1996/bookings/internal/forms"
	"github.com/TranQuocToan1996/bookings/internal/models"
	"github.com/TranQuocToan1996/bookings/internal/render"
)

// Repo the respository used by the handler
var Repo *Repository

// Repository is the repository type
type Repository struct {
	App *config.AppConfig
}

// NewRepo creates a new Repository
func NewRepo(a *config.AppConfig) *Repository {
	return &Repository{
		App: a,
	}
}

// NewHandlers sets the repository for the handler
func NewHandlers(r *Repository) {
	Repo = r
}

func (m *Repository) Contact(w http.ResponseWriter, r *http.Request) {
	render.RenderTemplate(w, r, "contact.page.html", &models.TemplateData{})
}

// Another page handle request and send back home page response
func (m *Repository) Home(w http.ResponseWriter, r *http.Request) {
	remoteIP := r.RemoteAddr
	m.App.Session.Put(r.Context(), "remote_ip", remoteIP)
	render.RenderTemplate(w, r, "home.page.html", &models.TemplateData{})
}

// About is the about page handler
func (m *Repository) About(w http.ResponseWriter, r *http.Request) {

	stringMap := make(map[string]string)
	stringMap["test"] = "Testing"

	remoteIP := m.App.Session.GetString(r.Context(), "remote_ip")
	stringMap["remote_ip"] = remoteIP

	// Send data to template
	render.RenderTemplate(w, r, "about.page.html", &models.TemplateData{
		StringMap: stringMap,
	})

}

// Renders the page
func (m *Repository) Reservation(w http.ResponseWriter, r *http.Request) {
	var emptyReservation models.Revervation
	data := make(map[string]interface{})
	data["reservation"] = emptyReservation

	render.RenderTemplate(w, r, "make-reservation.page.html", &models.TemplateData{
		Form: forms.New(nil),
		Data: data,
	})
}
func (m *Repository) PostReservation(w http.ResponseWriter, r *http.Request) {
	err := r.ParseForm()
	if err != nil {
		log.Println(err)
		return
	}

	reservation := models.Revervation{
		FirstName: r.Form.Get("first_name"),
		LastName:  r.Form.Get("last_name"),
		Email:     r.Form.Get("email"),
		Phone:     r.Form.Get("phone"),
	}

	// Create new form variable
	// r.PostForm contain all url and data.
	// PostForm contains the parsed form data from PATCH, POST or PUT body parameters.
	// This field is only available after ParseForm is called.
	form := forms.New(r.PostForm)

	// Check empty input in 4 fields
	form.Required("first_name", "last_name", "phone", "email")
	// Check length for first_name box
	form.MinLength("first_name", 3, r)
	form.MinLength("last_name", 3, r)
	//Check valid pattern
	form.IsPhoneNumber("phone")
	form.IsEmail("email")

	if !form.Valid() {
		data := make(map[string]interface{})
		data["reservation"] = reservation

		render.RenderTemplate(w, r, "make-reservation.page.html", &models.TemplateData{
			Form: form,
			Data: data,
		})
		return
	}

	// Write Reservation info into session, we will add logic to added this info into reservation-summary.page.html
	m.App.Session.Put(r.Context(), "reservation", reservation)

	// Redirect to /reservation-summary after send post request to prevent accident send post twice
	// http.StatusSeeOther - https://developer.mozilla.org/en-US/docs/Web/HTTP/Status/303
	http.Redirect(w, r, "/reservation-summary", http.StatusSeeOther)
}

func (m *Repository) Generals(w http.ResponseWriter, r *http.Request) {
	render.RenderTemplate(w, r, "generals.page.html", &models.TemplateData{})
}

func (m *Repository) Majors(w http.ResponseWriter, r *http.Request) {
	render.RenderTemplate(w, r, "majors.page.html", &models.TemplateData{})
}

func (m *Repository) Availability(w http.ResponseWriter, r *http.Request) {
	render.RenderTemplate(w, r, "search-availability.page.html", &models.TemplateData{})
}
func (m *Repository) PostAvailability(w http.ResponseWriter, r *http.Request) {
	// Render page for the after sending post request
	// Taking data from post request
	start := r.Form.Get("start")
	end := r.Form.Get("end")

	// Dont want render page again, just send a HTTP reply to brower
	w.Write([]byte(fmt.Sprintf("Start date is %s and End date is %s", start, end)))
}

func (m *Repository) ReservationSummary(w http.ResponseWriter, r *http.Request) {
	// Taking reservation info from session
	// Because the session don't know what type should return, so we use type assertion
	reservation, ok := m.App.Session.Get(r.Context(), "reservation").(models.Revervation)
	if !ok {
		m.App.Session.Put(r.Context(), "error", "Can't get reservation from session!")
		http.Redirect(w, r, "/", http.StatusTemporaryRedirect)
		return
	}

	// Remove reservation from session
	m.App.Session.Remove(r.Context(), "reservation")

	data := make(map[string]interface{})
	data["reservation"] = reservation
	render.RenderTemplate(w, r, "reservation-summary.page.html", &models.TemplateData{
		Data: data,
	})
}

type jsonResponse struct {
	// The member name must be captalize because JSON (un)marshaller uses reflection, it cannot read or write unexported fields
	// `` what the field will be recognized in json/xml
	OK      bool   `json:"ok"`
	Message string `json:"message"`
}

// Handle request from availability-page and send back JSON response
func (m *Repository) AvailabilityJSON(w http.ResponseWriter, r *http.Request) {
	resp := jsonResponse{
		OK:      true,
		Message: "Available!",
	}

	json, err := json.MarshalIndent(resp, "", "    ")
	if err != nil {
		log.Println(err)
	}
	// Send json back to client browser
	// log.Println(string(json))
	w.Header().Set("Content-Type", "application/json")
	w.Write(json)
}
