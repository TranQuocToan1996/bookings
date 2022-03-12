package handlers

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/TranQuocToan1996/bookings/internal/config"
	"github.com/TranQuocToan1996/bookings/internal/driver"
	"github.com/TranQuocToan1996/bookings/internal/forms"
	"github.com/TranQuocToan1996/bookings/internal/helpers"
	"github.com/TranQuocToan1996/bookings/internal/models"
	"github.com/TranQuocToan1996/bookings/internal/render"
	"github.com/TranQuocToan1996/bookings/internal/repository"
	"github.com/TranQuocToan1996/bookings/internal/repository/dbrepo"
	"github.com/go-chi/chi"
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
	reservation, ok := m.App.Session.Get(r.Context(), "reservation").(models.Reservation)
	if !ok {
		m.App.Session.Put(r.Context(), "error", "can't get reservation information from session!")
		http.Redirect(w, r, "/", http.StatusSeeOther)
		return
	}

	// Get room struct by ID and saving into
	room, err := m.DB.GetRoomByID(reservation.RoomID)
	if err != nil {
		m.App.Session.Put(r.Context(), "error", "can't find rooms!")
		http.Redirect(w, r, "/", http.StatusSeeOther)
		return
	}
	reservation.Room.RoomName = room.RoomName

	// Update reservation into session (startDate, endDate, roomName, roomID) and this data will take in PostReservation
	m.App.Session.Put(r.Context(), "reservation", reservation)

	data := make(map[string]interface{})
	data["reservation"] = reservation

	// Start date and end date are type time.Time, convert to string to using in template
	stringMap := map[string]string{
		"start_date": reservation.StartDate.Format(layout),
		"end_date":   reservation.EndDate.Format(layout),
	}

	render.Template(w, r, "make-reservation.page.html", &models.TemplateData{
		Form:      forms.New(nil),
		Data:      data,
		StringMap: stringMap,
	})
}

func (m *Repository) PostReservation(w http.ResponseWriter, r *http.Request) {
	// Update info from post form
	err := r.ParseForm()
	if err != nil {
		m.App.Session.Put(r.Context(), "error", "Can't parse form!")
		http.Redirect(w, r, "/", http.StatusSeeOther)
		return
	}

	// 15.5  Creating and sending mail notifications
	// Get roomID, roomName, startDate, endDate from session Reservation(Get method)
	reservation, ok := m.App.Session.Get(r.Context(), "reservation").(models.Reservation)
	if !ok {
		m.App.Session.Put(r.Context(), "error", "can't get reservation information from session!")
		http.Redirect(w, r, "/", http.StatusSeeOther)
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

		// add these lines to fix bad data error
		stringMap := make(map[string]string)
		stringMap["start_date"] = reservation.StartDate.Format(layout)
		stringMap["end_date"] = reservation.EndDate.Format(layout)

		// Replies error string and HTTP code
		// http.Error(w, "my own error message", http.StatusSeeOther)
		render.Template(w, r, "make-reservation.page.html", &models.TemplateData{
			Form:      form,
			Data:      data,
			StringMap: stringMap,
		})
		return
	}

	// after form validation, push data into database and get returned id
	newReservationID, err := m.DB.InsertReservation(&reservation)
	if err != nil {
		m.App.Session.Put(r.Context(), "error", "can't insert reservation into the database!")
		http.Redirect(w, r, "/", http.StatusSeeOther)
		return
	}

	restriction := models.RoomRestriction{
		StartDate:     reservation.StartDate,
		EndDate:       reservation.EndDate,
		RoomID:        reservation.RoomID,
		ReservationID: newReservationID,
		RestrictionID: 1, // This id is for reservation
	}

	err = m.DB.InsertRoomRestriction(&restriction)
	if err != nil {
		m.App.Session.Put(r.Context(), "error", "can't insert room restriction!")
		http.Redirect(w, r, "/", http.StatusSeeOther)
		return
	}

	// Send notifications - first to guest who wants book room
	htmlMessageGuest := fmt.Sprintf(`
		<strong>Reservation confirmation</strong><br>
		Dear %s:, <br>
		This is confirmed your reservation from %s to %s.
	`, reservation.FirstName,
		reservation.StartDate.Format(layout),
		reservation.EndDate.Format(layout))

	msg := models.MailData{
		To:       reservation.Email,
		From:     "me@here.com",
		Subject:  "Reservation confirmation",
		Content:  htmlMessageGuest,
		Template: "basic.html",
	}
	m.App.MailChan <- msg

	// Send notifications - first to Owner rooms
	htmlMessageOwner := fmt.Sprintf(`
		<strong>Reservation alert</strong> <br>
		Dear %s, <br>
		This is alerted your room from %s to %s.
	`, reservation.FirstName,
		reservation.StartDate.Format(layout),
		reservation.EndDate.Format(layout))

	msg = models.MailData{
		To:       reservation.Email,
		From:     "me@here.com",
		Subject:  "Reservation alert",
		Content:  htmlMessageOwner,
		Template: "basic.html",
	}
	m.App.MailChan <- msg

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
		m.App.ErrorLog.Println("Can't get reservation from session! Skip this log if in testing")
		m.App.Session.Put(r.Context(), "error", "Can't get reservation from session!")
		http.Redirect(w, r, "/", http.StatusSeeOther)
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

// Handler for logout user
func (m *Repository) Logout(w http.ResponseWriter, r *http.Request) {
	// Destroy all datas in the session
	err := m.App.Session.Destroy(r.Context())
	if err != nil {
		log.Println(err)
	}

	// Renew session
	err = m.App.Session.RenewToken(r.Context())
	if err != nil {
		log.Println(err)
	}

	http.Redirect(w, r, "/user/login", http.StatusSeeOther)
}

func (m *Repository) AdminDashboard(w http.ResponseWriter, r *http.Request) {
	render.Template(w, r, "admin-dashboard.page.html", &models.TemplateData{})
}

// AdminNewReservations shows all new reservations in admin dashboard
func (m *Repository) AdminNewReservations(w http.ResponseWriter, r *http.Request) {
	reservations, err := m.DB.AllNewReservations()
	if err != nil {
		helpers.ServerError(w, err)
		return
	}

	data := make(map[string]interface{})
	data["reservations"] = reservations

	render.Template(w, r, "admin-new-reservations.page.html", &models.TemplateData{
		Data: data,
	})
}

// AdminAllReservations shows all reservations in admin dashboard
func (m *Repository) AdminAllReservations(w http.ResponseWriter, r *http.Request) {
	reservations, err := m.DB.AllReservations()
	if err != nil {
		helpers.ServerError(w, err)
		return
	}

	data := make(map[string]interface{})
	data["reservations"] = reservations

	render.Template(w, r, "admin-all-reservations.page.html", &models.TemplateData{
		Data: data,
	})
}

// AdminShowReservations shows reservation in the admin page
func (m *Repository) AdminShowReservations(w http.ResponseWriter, r *http.Request) {

	// Get the URL params from request
	exploded := strings.Split(r.RequestURI, "/") // r.RequestURI = "/admin/reservations/{src}/{id}"
	src := exploded[3]
	id, err := strconv.Atoi(exploded[4])
	if err != nil {
		helpers.ServerError(w, err)
		return
	}

	// Get the reservation from the database
	res, err := m.DB.GetReservationByID(id)
	if err != nil {
		helpers.ServerError(w, err)
		return
	}

	// Stores and sends data to templates
	stringMap := make(map[string]string)
	year := r.URL.Query().Get("y")
	month := r.URL.Query().Get("m")
	stringMap["month"] = month
	stringMap["year"] = year
	stringMap["src"] = src
	data := make(map[string]interface{})
	data["reservation"] = res

	render.Template(w, r, "admin-reservations-show.page.html", &models.TemplateData{
		StringMap: stringMap,
		Data:      data,
		Form:      forms.New(nil),
	})
}

// ShowLogin shows the login sreen
func (m *Repository) ShowLogin(w http.ResponseWriter, r *http.Request) {
	render.Template(w, r, "login.page.html", &models.TemplateData{
		Form: forms.New(nil),
	})
}

// PostShowLogin handles the POST request sending login info
func (m *Repository) PostShowLogin(w http.ResponseWriter, r *http.Request) {
	// Call RenewToken method whenever we have operation that change privilege (IE login, logout)
	// This one Prevent session fixation attack
	err := m.App.Session.RenewToken(r.Context())
	if err != nil {
		log.Println(err)
	}

	err = r.ParseForm()
	if err != nil {
		log.Println(err)
	}
	email := r.Form.Get("email")
	password := r.Form.Get("password")

	// Check if received email and Plain text password from POST request
	form := forms.New(r.PostForm)
	form.Required("email", "password") // Check for blank
	form.IsEmail("email")
	if !form.Valid() {
		// Take user back to the page
		render.Template(w, r, "login.page.html", &models.TemplateData{
			Form: form,
		})
		return
	}

	id, _, err := m.DB.Authenticate(email, password)
	if err != nil {
		log.Println(err)

		// Put error into session, and redirect back to login page
		m.App.Session.Put(r.Context(), "error", "Invalid login credentials")
		http.Redirect(w, r, "/user/login", http.StatusSeeOther)
		return
	}

	// Store user_id in the session so that remembers the user login
	// AddDefaultData() will check login status by using this id and send that info to Template()
	m.App.Session.Put(r.Context(), "user_id", id)

	// Inform to user login success and redirect into home page
	m.App.Session.Put(r.Context(), "flash", "Logged in successfully!")
	http.Redirect(w, r, "/", http.StatusSeeOther)
}

// PostAvailability Renders page for the after sending post request
func (m *Repository) PostAvailability(w http.ResponseWriter, r *http.Request) {
	err := r.ParseForm()
	if err != nil {
		m.App.Session.Put(r.Context(), "error", "can't parse form!")
		http.Redirect(w, r, "/", http.StatusSeeOther)
		return
	}

	// Taking data from post request as type string and convert into time.Time
	startDate, err := time.Parse(layout, r.Form.Get("start"))
	if err != nil {
		m.App.Session.Put(r.Context(), "error", "can't parse start date!")
		http.Redirect(w, r, "/", http.StatusSeeOther)
		return
	}
	endDate, err := time.Parse(layout, r.Form.Get("end"))
	if err != nil {
		m.App.Session.Put(r.Context(), "error", "can't parse end date!")
		http.Redirect(w, r, "/", http.StatusSeeOther)
		return
	}

	rooms, err := m.DB.SearchAvailabilityForAllRooms(startDate, endDate)
	if err != nil {
		m.App.Session.Put(r.Context(), "error", "can't get availability for rooms")
		http.Redirect(w, r, "/", http.StatusSeeOther)
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
	//	http.Redirect(w, r, "/", http.StatusSeeOther)
	//	return
	//}

	// changed to this, so we can test it more easily
	// split the URL up by /, and grab the 3rd element (id)
	exploded := strings.Split(r.RequestURI, "/") /* http://localhost:8080/choose-room/{id} */
	roomID, err := strconv.Atoi(exploded[2])
	if err != nil {
		m.App.Session.Put(r.Context(), "error", "missing url parameter")
		http.Redirect(w, r, "/", http.StatusSeeOther)
		return
	}

	res, ok := m.App.Session.Get(r.Context(), "reservation").(models.Reservation)
	if !ok {
		m.App.Session.Put(r.Context(), "error", "Can't get reservation from session")
		http.Redirect(w, r, "/", http.StatusSeeOther)
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
		http.Redirect(w, r, "/", http.StatusSeeOther)
		return
	}

	startDate, err := time.Parse(layout, r.URL.Query().Get("s"))
	if err != nil {
		m.App.Session.Put(r.Context(), "error", "can't parse start date!")
		http.Redirect(w, r, "/", http.StatusSeeOther)
		return
	}
	endDate, err := time.Parse(layout, r.URL.Query().Get("e"))
	if err != nil {
		m.App.Session.Put(r.Context(), "error", "can't parse end date!")
		http.Redirect(w, r, "/", http.StatusSeeOther)
		return
	}
	room, err := m.DB.GetRoomByID(roomID)
	if err != nil {
		m.App.Session.Put(r.Context(), "error", "Error when querying room id from database!")
		http.Redirect(w, r, "/", http.StatusSeeOther)
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

// AdminPostShowReservations updates the reservation from the POST form into the database
func (m *Repository) AdminPostShowReservations(w http.ResponseWriter, r *http.Request) {
	// Update info from post form
	err := r.ParseForm()
	if err != nil {
		// this is in admin page so we can render the error for the administration
		helpers.ServerError(w, err)
		return
	}

	// Get the URL params from request
	exploded := strings.Split(r.RequestURI, "/") // r.RequestURI = "/admin/reservations/{src}/{id}"
	src := exploded[3]
	id, err := strconv.Atoi(exploded[4])
	if err != nil {
		helpers.ServerError(w, err)
		return
	}

	// Stores and sends data to templates
	stringMap := make(map[string]string)
	stringMap["src"] = src

	// Get the reservation from the database
	res, err := m.DB.GetReservationByID(id)
	if err != nil {
		helpers.ServerError(w, err)
		return
	}

	// Update the reservation with the data from the form
	res.FirstName = r.Form.Get("first_name")
	res.LastName = r.Form.Get("last_name")
	res.Email = r.Form.Get("email")
	res.Phone = r.Form.Get("phone")
	err = m.DB.UpdateReservation(res)
	if err != nil {
		helpers.ServerError(w, err)
		return
	}

	// Get month and year from post form (input tag)
	month := r.Form.Get("month")
	year := r.Form.Get("year")

	// Inform "Changes saved" to user and redirect to suitable page
	m.App.Session.Put(r.Context(), "flash", "Changes saved")
	if year == "" {
		http.Redirect(w, r, fmt.Sprintf("/admin/reservations-%s", src), http.StatusSeeOther)
	} else {
		http.Redirect(w, r, fmt.Sprintf("/admin/reservations-calendar?y=%s&m=%s", year, month), http.StatusSeeOther)
	}
}

// AdminProcessReservation marks a reservation as processed status
func (m *Repository) AdminProcessReservation(w http.ResponseWriter, r *http.Request) {
	// get URL params from "/admin/process-reservation/cal/1/do"
	id, _ := strconv.Atoi(chi.URLParam(r, "id"))

	src := chi.URLParam(r, "src")

	err := m.DB.UpdateProcessedForReservation(id, 1) // 1 is already processed status for a reservation
	if err != nil {
		helpers.ServerError(w, err)
		return
	}

	year := r.URL.Query().Get("y")
	month := r.URL.Query().Get("m")

	// Inform "Reservation marked as processed" to user and redirect to source page
	m.App.Session.Put(r.Context(), "flash", "Reservation marked as processed")

	if year == "" {
		http.Redirect(w, r, fmt.Sprintf("/admin/reservations-%s", src), http.StatusSeeOther)
	} else {
		http.Redirect(w, r, fmt.Sprintf("/admin/reservations-calendar?y=%s&m=%s", year, month), http.StatusSeeOther)
	}
}

// AdminDeleteReservation deletes a reservation from database
func (m *Repository) AdminDeleteReservation(w http.ResponseWriter, r *http.Request) {
	// get URL params from "/admin/reservations/{src}/{id}""
	id, _ := strconv.Atoi(chi.URLParam(r, "id"))

	src := chi.URLParam(r, "src")

	err := m.DB.DeleteReservation(id)
	if err != nil {
		helpers.ServerError(w, err)
		return
	}

	year := r.URL.Query().Get("y")
	month := r.URL.Query().Get("m")

	// Inform "Reservation marked as processed" to user and redirect to source page
	m.App.Session.Put(r.Context(), "flash", "Reservation delete")

	if year == "" {
		http.Redirect(w, r, fmt.Sprintf("/admin/reservations-%s", src), http.StatusSeeOther)
	} else {
		http.Redirect(w, r, fmt.Sprintf("/admin/reservations-calendar?y=%s&m=%s", year, month), http.StatusSeeOther)
	}
}

// AdminPostReservationsCalendar handles post of reservation calendar
func (m *Repository) AdminPostReservationsCalendar(w http.ResponseWriter, r *http.Request) {
	err := r.ParseForm()
	if err != nil {
		helpers.ServerError(w, err)
		return
	}
	form := forms.New(r.PostForm)

	// Get hidden filed y and m in calendar page
	year, _ := strconv.Atoi(r.Form.Get("y"))
	month, _ := strconv.Atoi(r.Form.Get("m"))

	rooms, err := m.DB.AllRooms()
	if err != nil {
		helpers.ServerError(w, err)
		return
	}

	// Process unblocking in calendar page
	for _, room := range rooms {
		// Get the blockMap from the session (From AdminReservationsCalendar handler)
		/* Loop through the blockMap, if we have an entry in the map that does not exist in our posted data, and if the restriction id > 0, then it is a block we need to remove */
		currentMap := m.App.Session.Get(r.Context(), fmt.Sprintf("block_map_%d", room.ID)).(map[string]int)
		for name := range currentMap {
			// check whether name in th blockMap
			if val, ok := currentMap[name]; ok {
				// Pay attention only to value > 0 and not in the form post (Uncheck)
				if val > 0 {
					if !form.Has(fmt.Sprintf("remove_block_%d_%s", room.ID, name)) {
						// delete restriction by id
						err := m.DB.DeleteBlockByID(val)
						if err != nil {
							log.Println(err)
						}
						log.Println("delete block successful the:", name)
					}
				}
			}
		}

		// Process new blocking
		for name := range r.PostForm {
			if strings.HasPrefix(name, "add_block") {
				exploded := strings.Split(name, "_")
				roomID, _ := strconv.Atoi(exploded[2])
				startDate, _ := time.Parse("2006-01-2", exploded[3])
				// Inserts a new block
				err := m.DB.InsertBlockForRoom(roomID, startDate)
				if err != nil {
					log.Println(err)
				}
				log.Println("insert block complete the startdate:", startDate)

			}
		}
	}

	m.App.Session.Put(r.Context(), "flash", "Changes saved")
	http.Redirect(w, r, fmt.Sprintf("/admin/reservations-calendar?y=%d&m=%d", year, month), http.StatusSeeOther)
}

// AdminReservationsCalendar displays the reservation calendar
func (m *Repository) AdminReservationsCalendar(w http.ResponseWriter, r *http.Request) {
	// Passing current time into template
	now := time.Now()

	// get the data from URL query
	if r.URL.Query().Get("y") != "" && r.URL.Query().Get("m") != "" {
		year, err := strconv.Atoi(r.URL.Query().Get("y"))
		if err != nil {
			helpers.ServerError(w, err)
			return
		}
		month, err := strconv.Atoi(r.URL.Query().Get("m"))
		if err != nil {
			helpers.ServerError(w, err)
			return
		}
		now = time.Date(year, time.Month(month), 1, 0, 0, 0, 0, time.UTC)
	}
	// Pass data to the templates
	data := make(map[string]interface{})
	data["now"] = now

	// time.Time format
	nextMonth_timeTime := now.AddDate(0, 1, 0)
	lastMonth_timeTime := now.AddDate(0, -1, 0)

	// Format to string the month with the "01" month layout
	nextMonth := nextMonth_timeTime.Format("01")
	nextMonthYear := nextMonth_timeTime.Format("2006")
	lastMonth := lastMonth_timeTime.Format("01")
	lastMonthYear := lastMonth_timeTime.Format("2006")

	// Pass data to the templates
	stringMap := make(map[string]string)
	stringMap["last_month"] = lastMonth
	stringMap["last_month_year"] = lastMonthYear
	stringMap["this_month"] = now.Format("01")
	stringMap["this_month_year"] = now.Format("2006")
	stringMap["next_month"] = nextMonth
	stringMap["next_month_year"] = nextMonthYear

	// Get the first and last day of the month and passing into templates
	currentYear, currentMonth, _ := now.Date()
	currentLocation := now.Location()
	firstOfMonth := time.Date(currentYear, currentMonth, 1, 0, 0, 0, 0, currentLocation)
	lastOfMonth := firstOfMonth.AddDate(0, 1, -1)
	intMap := make(map[string]int)
	intMap["days_in_month"] = lastOfMonth.Day()

	// Get rooms from database
	rooms, err := m.DB.AllRooms()
	if err != nil {
		helpers.ServerError(w, err)
		return
	}
	data["rooms"] = rooms

	for _, room := range rooms {
		// create maps
		reservationMap := make(map[string]int)
		blockMap := make(map[string]int)

		for d := firstOfMonth; !d.After(lastOfMonth); d = d.AddDate(0, 0, 1) {
			reservationMap[d.Format("2006-01-2")] = 0
			blockMap[d.Format("2006-01-2")] = 0
		}

		// get all the restrictions for the current room
		restrictions, err := m.DB.GetRestrictionsForRoomByDate(room.ID, firstOfMonth, lastOfMonth)
		if err != nil {
			helpers.ServerError(w, err)
			return
		}

		for _, restriction := range restrictions {
			if restriction.ReservationID > 0 {
				// it's a reservation
				for d := restriction.StartDate; !d.After(restriction.EndDate); d = d.AddDate(0, 0, 1) {
					reservationMap[d.Format("2006-01-2")] = restriction.ReservationID
				}
			} else {
				// it's a block
				blockMap[restriction.StartDate.Format("2006-01-2")] = restriction.ID
			}
		}
		data[fmt.Sprintf("reservation_map_%d", room.ID)] = reservationMap
		data[fmt.Sprintf("block_map_%d", room.ID)] = blockMap

		m.App.Session.Put(r.Context(), fmt.Sprintf("block_map_%d", room.ID), blockMap)
	}

	render.Template(w, r, "admin-reservations-calendar.page.html", &models.TemplateData{
		StringMap: stringMap,
		Data:      data,
		IntMap:    intMap,
	})
}
