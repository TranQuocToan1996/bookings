package handlers

import (
	"encoding/json"
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/TranQuocToan1996/bookings/internal/config"
	"github.com/TranQuocToan1996/bookings/internal/driver"
	"github.com/TranQuocToan1996/bookings/internal/forms"
	"github.com/TranQuocToan1996/bookings/internal/models"
	"github.com/TranQuocToan1996/bookings/internal/render"
	"github.com/TranQuocToan1996/bookings/internal/repository"
	"github.com/TranQuocToan1996/bookings/internal/repository/dbrepo"
)

// Const variable layout for format time.Time
const layout string = "2006-01-02"

// Repo the respository used by the handler
var Repo *Repository

// Repository is the repository type
type Repository struct {
	App *config.AppConfig
	DB  repository.DatabaseRepo
}

// NewRepo creates a new Repository
func NewRepo(a *config.AppConfig, db *driver.DB) *Repository {
	return &Repository{
		App: a,
		DB:  dbrepo.NewPostgresRepo(db.SQL, a),
	}
}

// NewHandlers sets the repository for the handler
func NewHandlers(r *Repository) {
	Repo = r
}

// Contact renders the contact page
func (m *Repository) Contact(w http.ResponseWriter, r *http.Request) {
	render.Template(w, r, "contact.page.html", &models.TemplateData{})
}

// Another page handle request and send back home page response
func (m *Repository) Home(w http.ResponseWriter, r *http.Request) {
	/* 	Taking user IP address and save into session
	   	remoteIP := r.RemoteAddr
	   	m.App.Session.Put(r.Context(), "remote_ip", remoteIP) */
	render.Template(w, r, "home.page.html", &models.TemplateData{})
}

// About is the about page handler
func (m *Repository) About(w http.ResponseWriter, r *http.Request) {

	// Send data to template
	render.Template(w, r, "about.page.html", &models.TemplateData{})

}

// Reservation renders make-reservation page
func (m *Repository) Reservation(w http.ResponseWriter, r *http.Request) {

	// Get the reservation info from session (Choose-room) and render make-reservation page
	res, ok := m.App.Session.Get(r.Context(), "reservation").(models.Reservation)
	if !ok {
		m.App.Session.Put(r.Context(), "error", "can't get reservation information from session!")
		http.Redirect(w, r, "/", http.StatusTemporaryRedirect)
		return
	}

	// Get room struct by ID and saving into
	room, err := m.DB.GetRoomByID(res.RoomID)
	if err != nil {
		m.App.Session.Put(r.Context(), "error", "can't find rooms!")
		http.Redirect(w, r, "/", http.StatusTemporaryRedirect)
		return
	}
	res.Room.RoomName = room.RoomName

	// Update reservation into session (startDate, endDate, roomName, roomID) and this data will take in PostReservation
	m.App.Session.Put(r.Context(), "reservation", res)

	data := make(map[string]interface{})
	data["reservation"] = res

	// Start date and end date are type time.Time, convert to string to using in template
	stringMap := map[string]string{
		"start_date": res.StartDate.Format(layout),
		"end_date":   res.EndDate.Format(layout),
	}

	render.Template(w, r, "make-reservation.page.html", &models.TemplateData{
		Form:      forms.New(nil),
		Data:      data,
		StringMap: stringMap,
	})
}

func (m *Repository) PostReservation(w http.ResponseWriter, r *http.Request) {

	// Get roomID, roomName, startDate, endDate from session Reservation(Get method)
	reservation, ok := m.App.Session.Get(r.Context(), "reservation").(models.Reservation)
	if !ok {
		m.App.Session.Put(r.Context(), "error", "can't get reservation information from session!")
		http.Redirect(w, r, "/", http.StatusTemporaryRedirect)
		return
	}

	/* 	Old code, getting data from the post form -> user had to type more than now
	   	//Get date and format to Vietnam date
	   	layout := "2006-01-02" // Mon Jan 2 15:04:05 -0700 MST 2006 (01/02 03:04:05PM '06 -0700) --> yyyy-mm-dd
	   	startDate, err := time.Parse(layout, r.Form.Get("start_date"))
	   	if err != nil {
	   		helpers.ServerError(w, err)
	   		return
	   	}
	   	endDate, err := time.Parse(layout, r.Form.Get("end_date"))
	   	if err != nil {
	   		helpers.ServerError(w, err)
	   		return
	   	}
	   	roomID, err := strconv.Atoi(r.Form.Get("room_id"))
	   	if err != nil {
	   		helpers.ServerError(w, err)
	   		return
	   	}
	   	reservation := models.Reservation{
	   		FirstName: r.Form.Get("first_name"),
	   		LastName:  r.Form.Get("last_name"),
	   		Email:     r.Form.Get("email"),
	   		Phone:     r.Form.Get("phone"),
	   		StartDate: startDate,
	   		EndDate:   endDate,
	   		RoomID:    roomID,
	   	} */

	// Update info from post form
	err := r.ParseForm()
	if err != nil {
		m.App.Session.Put(r.Context(), "error", "Can't parse form!")
		http.Redirect(w, r, "/", http.StatusTemporaryRedirect)
		return
	}
	reservation.FirstName = r.Form.Get("first_name")
	reservation.LastName = r.Form.Get("last_name")
	reservation.Phone = r.Form.Get("phone")
	reservation.Email = r.Form.Get("email")

	form := forms.New(r.PostForm)

	// Check input from post request
	form.Required("first_name", "last_name", "phone", "email")
	form.MinLength("first_name", 2)
	form.MinLength("last_name", 2)
	form.IsPhoneNumber("phone")
	form.IsEmail("email")

	// if not valid, take data back to form, highlight the part where error exist
	if !form.Valid() {
		data := make(map[string]interface{})
		data["reservation"] = reservation
		// Replies error string and HTTP code
		http.Error(w, "my own error message", http.StatusTemporaryRedirect)
		render.Template(w, r, "make-reservation.page.html", &models.TemplateData{
			Form: form,
			Data: data,
		})
		return
	}

	// after form validation, push data into database and get returned id
	newReservationID, err := m.DB.InsertReservation(&reservation)
	if err != nil {
		m.App.Session.Put(r.Context(), "error", "can't insert reservation into the database!")
		http.Redirect(w, r, "/", http.StatusTemporaryRedirect)
		return
	}

	restriction := models.RoomRestriction{
		StartDate:     reservation.StartDate,
		EndDate:       reservation.EndDate,
		RoomID:        reservation.RoomID,
		ReservationID: newReservationID,
		RestrictionID: 1,
	}

	err = m.DB.InsertRoomRestriction(&restriction)
	if err != nil {
		m.App.Session.Put(r.Context(), "error", "can't insert room restriction!")
		http.Redirect(w, r, "/", http.StatusTemporaryRedirect)
		return
	}

	// Update reservation into session
	// Write Reservation info into session, we will add logic to added this info into reservation-summary.page.html
	m.App.Session.Put(r.Context(), "reservation", reservation)

	// Redirect to /reservation-summary after send post request to prevent accident send post twice
	// http.StatusSeeOther - https://developer.mozilla.org/en-US/docs/Web/HTTP/Status/303
	http.Redirect(w, r, "/reservation-summary", http.StatusSeeOther)
}

// ReservationSummary displays reservation summary page
func (m *Repository) ReservationSummary(w http.ResponseWriter, r *http.Request) {
	// Taking reservation info from session
	// Because the session don't know what type should return, so we use type assertion
	reservation, ok := m.App.Session.Get(r.Context(), "reservation").(models.Reservation)
	if !ok {
		m.App.ErrorLog.Println("Can't get reservation from session!")
		m.App.Session.Put(r.Context(), "error", "Can't get reservation from session!")
		http.Redirect(w, r, "/", http.StatusTemporaryRedirect)
		return
	}

	// Remove reservation from session
	m.App.Session.Remove(r.Context(), "reservation")

	// Prepare data for Template rendering
	data := make(map[string]interface{})
	data["reservation"] = reservation
	stringMap := make(map[string]string)
	stringMap["start_date"] = reservation.StartDate.Format(layout)
	stringMap["end_date"] = reservation.EndDate.Format(layout)

	render.Template(w, r, "reservation-summary.page.html", &models.TemplateData{
		Data:      data,
		StringMap: stringMap,
	})
}

func (m *Repository) Generals(w http.ResponseWriter, r *http.Request) {
	render.Template(w, r, "generals.page.html", &models.TemplateData{})
}

func (m *Repository) Majors(w http.ResponseWriter, r *http.Request) {
	render.Template(w, r, "majors.page.html", &models.TemplateData{})
}

func (m *Repository) Availability(w http.ResponseWriter, r *http.Request) {
	render.Template(w, r, "search-availability.page.html", &models.TemplateData{})
}

// PostAvailability Renders page for the after sending post request
func (m *Repository) PostAvailability(w http.ResponseWriter, r *http.Request) {
	err := r.ParseForm()
	if err != nil {
		m.App.Session.Put(r.Context(), "error", "can't parse form!")
		http.Redirect(w, r, "/", http.StatusTemporaryRedirect)
		return
	}

	// Taking data from post request as type string and convert into time.Time
	startDate, err := time.Parse(layout, r.Form.Get("start"))
	if err != nil {
		m.App.Session.Put(r.Context(), "error", "can't parse start date!")
		http.Redirect(w, r, "/", http.StatusTemporaryRedirect)
		return
	}
	endDate, err := time.Parse(layout, r.Form.Get("end"))
	if err != nil {
		m.App.Session.Put(r.Context(), "error", "can't parse end date!")
		http.Redirect(w, r, "/", http.StatusTemporaryRedirect)
		return
	}

	rooms, err := m.DB.SearchAvailabilityForAllRooms(startDate, endDate)
	if err != nil {
		m.App.Session.Put(r.Context(), "error", "can't get availability for rooms")
		http.Redirect(w, r, "/", http.StatusTemporaryRedirect)
		return
	}

	// If no room available, redirect and popup notie "no available room"
	if len(rooms) == 0 {
		m.App.Session.Put(r.Context(), "error", "No availability room!")
		http.Redirect(w, r, "/search-availability", http.StatusSeeOther)
		return
	}

	// {{$rooms := index .Data "rooms"}}
	data := make(map[string]interface{})
	data["rooms"] = rooms

	// Store date into session
	res := models.Reservation{
		StartDate: startDate,
		EndDate:   endDate,
	}
	m.App.Session.Put(r.Context(), "reservation", res)

	render.Template(w, r, "choose-room.page.html", &models.TemplateData{
		Data: data,
	})

	/* 	Old code
	//Don't want render page again, just send a HTTP reply to brower
	//w.Write([]byte(fmt.Sprintf("Start date is %s and End date is %s", startDate, endDate))) */
}

type jsonResponse struct {
	// The member name must be captalize because JSON (un)marshaller uses reflection, it cannot read or write unexported fields
	// `` what the field will be recognized in json/xml
	OK        bool   `json:"ok"`
	Message   string `json:"message"`
	RoomID    string `json:"room_id"`
	StartDate string `json:"start_date"`
	EndDate   string `json:"end_date"`
}

// Handle request from availability-page and send back JSON response
func (m *Repository) AvailabilityJSON(w http.ResponseWriter, r *http.Request) {

	// need to parse request body
	err := r.ParseForm()
	if err != nil {
		resp := jsonResponse{
			OK:      false,
			Message: "Internal server error!",
		}
		out, _ := json.MarshalIndent(resp, "", "     ")
		w.Header().Set("Content-Type", "application/json")
		w.Write(out)
		return
	}

	// Search availale for room by id
	roomID, err := strconv.Atoi(r.Form.Get("room_id"))
	if err != nil {
		resp := jsonResponse{
			OK:      false,
			Message: "Can't get room ID from POST request",
		}
		out, _ := json.MarshalIndent(resp, "", "    ")
		w.Header().Set("Content-type", "application/json")
		w.Write(out)
		return
	}
	startDate, err := time.Parse(layout, r.Form.Get("start"))
	if err != nil {
		resp := jsonResponse{
			OK:      false,
			Message: "Can't get start date from POST request!",
		}
		out, _ := json.MarshalIndent(resp, "", "    ")
		w.Header().Set("Content-type", "application/json")
		w.Write(out)
		return
	}
	endDate, err := time.Parse(layout, r.Form.Get("end"))
	if err != nil {
		resp := jsonResponse{
			OK:      false,
			Message: "Can't get end date from POST request",
		}
		out, _ := json.MarshalIndent(resp, "", "    ")
		w.Header().Set("Content-type", "application/json")
		w.Write(out)
		return
	}
	available, err := m.DB.SearchAvailabilityByRoomID(startDate, endDate, roomID)
	if err != nil {
		resp := jsonResponse{
			OK:      false,
			Message: "Connecting to database error!",
		}
		out, _ := json.MarshalIndent(resp, "", "    ")
		w.Header().Set("Content-type", "application/json")
		w.Write(out)
		return
	}

	// Create response json to client
	resp := jsonResponse{
		OK:        available,
		Message:   "",
		StartDate: r.Form.Get("start"),
		EndDate:   r.Form.Get("end"),
		RoomID:    strconv.Itoa(roomID),
	}
	json, _ := json.MarshalIndent(resp, "", "    ") // No error here because the resp json is constructed manually

	// Send json back to client browser (general.page.html check fetch function)
	// log.Println(string(json))
	w.Header().Set("Content-Type", "application/json")
	w.Write(json)
}

// ChooseRoom displays list of available rooms
func (m *Repository) ChooseRoom(w http.ResponseWriter, r *http.Request) {
	// used to have next 6 lines
	//roomID, err := strconv.Atoi(chi.URLParam(r, "id"))
	//if err != nil {
	//	log.Println(err)
	//	m.App.Session.Put(r.Context(), "error", "missing url parameter")
	//	http.Redirect(w, r, "/", http.StatusTemporaryRedirect)
	//	return
	//}

	// changed to this, so we can test it more easily
	// split the URL up by /, and grab the 3rd element (id)
	exploded := strings.Split(r.RequestURI, "/") /* http://localhost:8080/choose-room/{id} */
	roomID, err := strconv.Atoi(exploded[2])
	if err != nil {
		m.App.Session.Put(r.Context(), "error", "missing url parameter")
		http.Redirect(w, r, "/", http.StatusTemporaryRedirect)
		return
	}

	res, ok := m.App.Session.Get(r.Context(), "reservation").(models.Reservation)
	if !ok {
		m.App.Session.Put(r.Context(), "error", "Can't get reservation from session")
		http.Redirect(w, r, "/", http.StatusTemporaryRedirect)
		return
	}

	res.RoomID = roomID

	m.App.Session.Put(r.Context(), "reservation", res)

	http.Redirect(w, r, "/make-reservation", http.StatusSeeOther)

}

// BookRoom takes url params, build a sessional varialbe, and redirect user to res page
func (m *Repository) BookRoom(w http.ResponseWriter, r *http.Request) {
	// get data from GET request: id, s, e
	// /book-room?s=2050-01-01&e=2050-01-02&id=1
	roomID, err := strconv.Atoi(r.URL.Query().Get("id"))
	if err != nil {
		m.App.Session.Put(r.Context(), "error", "can't convert roomID from GET method into int!")
		http.Redirect(w, r, "/", http.StatusTemporaryRedirect)
		return
	}

	startDate, err := time.Parse(layout, r.URL.Query().Get("s"))
	if err != nil {
		m.App.Session.Put(r.Context(), "error", "can't parse start date!")
		http.Redirect(w, r, "/", http.StatusTemporaryRedirect)
		return
	}
	endDate, err := time.Parse(layout, r.URL.Query().Get("e"))
	if err != nil {
		m.App.Session.Put(r.Context(), "error", "can't parse end date!")
		http.Redirect(w, r, "/", http.StatusTemporaryRedirect)
		return
	}
	room, err := m.DB.GetRoomByID(roomID)
	if err != nil {
		m.App.Session.Put(r.Context(), "error", "Error when querying room id from database!")
		http.Redirect(w, r, "/", http.StatusTemporaryRedirect)
		return
	}

	var res models.Reservation
	res.RoomID = roomID
	res.StartDate = startDate
	res.EndDate = endDate
	res.Room.RoomName = room.RoomName

	// Put res into session
	m.App.Session.Put(r.Context(), "reservation", res)

	http.Redirect(w, r, "/make-reservation", http.StatusSeeOther)
}
