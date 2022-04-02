package handlers

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"net/url"
	"reflect"
	"strings"
	"testing"
	"time"

	"github.com/TranQuocToan1996/bookings/internal/driver"
	"github.com/TranQuocToan1996/bookings/internal/models"
)

// Reservation data for some tests require reservation in session
var reservation = models.Reservation{
	RoomID: 1,
	Room: models.Room{
		ID:       1,
		RoomName: "General's Quarters",
	},
}

// Test data is the slice of struct
var theTestsGET = []struct {
	testName         string
	url              string
	requestMethod    string
	expectStatusCode int
}{
	{"home", "/", "GET", http.StatusOK},
	{"about", "/about", "GET", http.StatusOK},
	{"gq", "/generals-quarters", "GET", http.StatusOK},
	{"ms", "/majors-suite", "GET", http.StatusOK},
	{"sa", "/search-availability", "GET", http.StatusOK},
	{"contact", "/contact", "GET", http.StatusOK},
	{"non-existent", "/green/eggs/and/ham", "GET", http.StatusNotFound},
	{"login", "/user/login", "GET", http.StatusOK},
	{"logout", "/user/logout", "GET", http.StatusOK},
	{"dashboard", "/admin/dashboard", "GET", http.StatusOK},
	{"new res", "/admin/reservations-new", "GET", http.StatusOK},
	{"all res", "/admin/reservations-all", "GET", http.StatusOK},
	{"show res", "/admin/reservations/new/1/show", "GET", http.StatusOK},
	{"show res cal", "/admin/reservations-calendar", "GET", http.StatusOK},
	{"show res cal with params", "/admin/reservations-calendar?y=2020&m=1", "GET", http.StatusOK},
}

func TestHanlers(t *testing.T) {
	routes := getRoutes()

	// Create a server for testing, close it when TestHanlers finish
	testServer := httptest.NewTLSServer(routes)
	defer testServer.Close()

	for _, e := range theTestsGET {
		// Create cliend brower and send get method
		resp, err := testServer.Client().Get(testServer.URL + e.url)
		if err != nil {
			t.Log(err)
			t.Fatal(err)
		}
		if resp.StatusCode != e.expectStatusCode {
			t.Errorf("For %s expecting %d but got %d", e.testName, e.expectStatusCode, resp.StatusCode)
		}

	}
}

func TestRepository_Reservation(t *testing.T) {

	// In this test function, we can't get reservation info from session so that we have to manually create a variable instead
	var reservation = models.Reservation{
		RoomID: 1,
		Room: models.Room{
			ID:       1,
			RoomName: "General's Quarters",
		},
	}

	// Create Request, the aim to push reservation into session of this request
	req, _ := http.NewRequest("GET", "/make-reservation", nil)
	// Get context for request
	ctx := getCtx(req)
	// Add infomation about X-session header (session) from ctx into request
	req = req.WithContext(ctx)

	// responseRecorder is an implementation of http.ResponseWriter that
	// records its mutations for later inspection in tests.
	// Simulation response cycle: Hit page -> pass request -> get response write -> response writer write to browser
	responseRecorder := httptest.NewRecorder()

	// Put into session
	session.Put(ctx, "reservation", reservation)

	// To directly execute (Don't need to hit the page in browser)
	handler := http.HandlerFunc(Repo.Reservation) // Turn Reservation() into HandlerFunc
	handler.ServeHTTP(responseRecorder, req)      // Execute Reservation()

	if responseRecorder.Code != http.StatusOK {
		t.Errorf("Reservation handler returned wrong code: Got %d, wanted %d", responseRecorder.Code, http.StatusOK)
	}

	// Test case where reservation is not in session (reset everything)
	req, _ = http.NewRequest("GET", "/make-reservation", nil)
	ctx = getCtx(req)
	req = req.WithContext(ctx)
	/* 	Notice that we remove the part add reservation information into session
	session.Put(ctx, "reservation", reservation)  */
	responseRecorder = httptest.NewRecorder()
	handler.ServeHTTP(responseRecorder, req)
	if responseRecorder.Code != http.StatusSeeOther {
		t.Errorf("Reservation handler returned wrong code: Got %d, wanted %d", responseRecorder.Code, http.StatusSeeOther)
	}

	// Test with non-exist room
	req, _ = http.NewRequest("GET", "/make-reservation", nil)
	ctx = getCtx(req)
	req = req.WithContext(ctx)
	responseRecorder = httptest.NewRecorder()
	reservation.RoomID = 999999 // Test case for roomID that out of range
	session.Put(ctx, "reservation", reservation)
	handler.ServeHTTP(responseRecorder, req)
	if responseRecorder.Code != http.StatusSeeOther {
		t.Errorf("Reservation handler returned wrong code: Got %d, wanted %d", responseRecorder.Code, http.StatusSeeOther)
	}
}

func TestRepository_PostReservation(t *testing.T) {

	/* Case 1: Request with no reservation in the session */
	req, _ := http.NewRequest("POST", "/make-reservation", nil)
	ctx := getCtx(req)
	req = req.WithContext(ctx)
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	responseRecorder := httptest.NewRecorder()
	/* Skip this code so no reservation in session
	session.Put(ctx, "reservation", reservation) */
	handler := http.HandlerFunc(Repo.PostReservation)
	handler.ServeHTTP(responseRecorder, req)
	if responseRecorder.Code != http.StatusSeeOther {
		t.Errorf("Reservation handler returned wrong code: Got %d, wanted %d", responseRecorder.Code, http.StatusSeeOther)
	}

	/* Case 2: Request with no POST body*/
	req, _ = http.NewRequest("POST", "/make-reservation", nil)
	ctx = getCtx(req)
	req = req.WithContext(ctx)
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	responseRecorder = httptest.NewRecorder()
	session.Put(ctx, "reservation", reservation) // Add reservation to session
	handler = http.HandlerFunc(Repo.PostReservation)
	handler.ServeHTTP(responseRecorder, req)
	if responseRecorder.Code != http.StatusSeeOther {
		t.Errorf("Reservation handler returned wrong code: Got %d, wanted %d", responseRecorder.Code, http.StatusSeeOther)
	}

	/* Case 3: Request with POST body, reservation in session*/
	// Add POST body: Primary way
	postedData := url.Values{}
	postedData.Add("start_date", "2050-01-01")
	postedData.Add("end_date", "2050-01-02")
	postedData.Add("first_name", "John")
	postedData.Add("last_name", "Smith")
	postedData.Add("email", "example@example.com")
	postedData.Add("phone", "0999999999")

	req, _ = http.NewRequest("POST", "/make-reservation", strings.NewReader(postedData.Encode()))
	ctx = getCtx(req)
	req = req.WithContext(ctx)
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	responseRecorder = httptest.NewRecorder()
	session.Put(ctx, "reservation", reservation) // Add reservation to session
	handler = http.HandlerFunc(Repo.PostReservation)
	handler.ServeHTTP(responseRecorder, req)
	if responseRecorder.Code != http.StatusSeeOther {
		t.Errorf("Reservation handler returned wrong code: Got %d, wanted %d", responseRecorder.Code, http.StatusSeeOther)
	}

	/* Case 4: Form.Valid() == false*/

	// Add POST body: Secondary way
	// generation request body with some form params
	// start_date=2050-01-01&end_date=2050-01-02&first_name=John&last_name=Smith&email=john@smith.com&phone=123456789&room_id=1
	reqBody := "start_date=2050-01-01"
	reqBody = fmt.Sprintf("%s&%s", reqBody, "end_date=2050-01-02")
	reqBody = fmt.Sprintf("%s&%s", reqBody, "first_name=J")
	reqBody = fmt.Sprintf("%s&%s", reqBody, "last_name=Smith")
	reqBody = fmt.Sprintf("%s&%s", reqBody, "email=Invalid") // Invalid data
	reqBody = fmt.Sprintf("%s&%s", reqBody, "phone=Invalid")

	req, _ = http.NewRequest("POST", "/make-reservation", strings.NewReader(reqBody))
	ctx = getCtx(req)
	req = req.WithContext(ctx)
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	responseRecorder = httptest.NewRecorder()
	session.Put(ctx, "reservation", reservation)
	handler = http.HandlerFunc(Repo.PostReservation)
	handler.ServeHTTP(responseRecorder, req)

	// If form doesn't valid, we render the make-reservation.page.html with some data, so that the status must be 200(http.StatusOK)
	if responseRecorder.Code != http.StatusOK {
		t.Errorf("Reservation handler returned wrong code: Got %d, wanted %d", responseRecorder.Code, http.StatusSeeOther)
	}

	/* Case 5: InsertReservation error*/
	// set up
	var reservation = models.Reservation{
		RoomID: 2,
	}

	postedData = url.Values{}
	postedData.Add("start_date", "2050-01-01")
	postedData.Add("end_date", "2050-01-02")
	postedData.Add("first_name", "John")
	postedData.Add("last_name", "Smith")
	postedData.Add("email", "example@example.com")
	postedData.Add("phone", "0999999999")

	req, _ = http.NewRequest("POST", "/make-reservation", strings.NewReader(postedData.Encode()))
	ctx = getCtx(req)
	req = req.WithContext(ctx)
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	responseRecorder = httptest.NewRecorder()
	session.Put(ctx, "reservation", reservation) // Add reservation to session
	handler = http.HandlerFunc(Repo.PostReservation)
	handler.ServeHTTP(responseRecorder, req)
	if responseRecorder.Code != http.StatusSeeOther {
		t.Errorf("Reservation handler returned wrong code: Got %d, wanted %d", responseRecorder.Code, http.StatusSeeOther)
	}

}

var testData_AvailabilityJSON = []struct {
	testName  string
	startDate string
	endDate   string
	roomID    string
}{
	// 2050-01-01 is already booked, 2060-01-01 is available
	{"roomNotAvailable", "2050-01-01", "2060-01-01", "1"},
	{"roomAvailable", "2060-01-01", "2060-01-01", "1"},
	{"noRequestBody", "2060-01-01", "2060-01-01", "1"},
	{"databaseErrors", "2050-01-01", "2050-01-01", "1"},
	{"failConvertStartDate", "Invalid", "2060-01-02", "1"},
	{"failConvertEndDate", "2060-01-02", "Invalid", "1"},
	{"failConvertRoomID", "2060-01-01", "2060-01-01", "Invalid"},
}

func TestRepository_AvailabilityJSON(t *testing.T) {
	for _, testData := range testData_AvailabilityJSON {
		var req *http.Request
		if testData.testName != "noRequestBody" {
			// create our request body
			postedData := url.Values{}
			postedData.Add("start", testData.startDate)
			postedData.Add("end", testData.endDate)
			postedData.Add("room_id", testData.roomID)

			// create our request with attach body
			req, _ = http.NewRequest("POST", "/search-availability-json", strings.NewReader(postedData.Encode()))

		} else {
			req, _ = http.NewRequest("POST", "/search-availability-json", nil)
		}
		// get the context with session
		ctx := getCtx(req)
		req = req.WithContext(ctx)

		// set the request header
		req.Header.Set("Content-Type", "application/x-www-form-urlencoded")

		// create our response recorder, which satisfies the requirements for http.ResponseWriter
		rr := httptest.NewRecorder()

		// make our handler a http.HandlerFunc
		handler := http.HandlerFunc(Repo.AvailabilityJSON)

		// make the request to our handler
		handler.ServeHTTP(rr, req)

		// this time we want to parse JSON and get the expected response
		var j jsonResponse
		err := json.Unmarshal([]byte(rr.Body.String()), &j)
		if err != nil {
			t.Error("failed to parse json!")
		}
		switch testData.testName {
		case "roomNotAvailable":
			if j.OK {
				t.Error("Got availability when none was expected in AvailabilityJSON")
			}
		case "roomAvailable":
			if !j.OK {
				t.Error("Got no availability when some was expected in AvailabilityJSON")
			}
		case "noRequestBody":
			if j.OK || j.Message != "Internal server error!" {
				t.Error("Got availability when request body of POST method was empty")
			}
		case "databaseErrors":
			if j.OK || j.Message != "Connecting to database error!" {
				t.Error("Got availability when simulating database error")
			}
		case "failConvertStartDate":
			if j.OK || j.Message != "Can't get start date from POST request!" {
				t.Error("Got availability when simulating database error")
			}
		case "failConvertEndDate":
			if j.OK || j.Message != "Can't get end date from POST request" {
				t.Error("Got availability when simulating database error")
			}
		case "failConvertRoomID":
			if j.OK || j.Message != "Can't get room ID from POST request" {
				t.Error("Got availability when simulating database error")
			}
		}
	}
}

func TestRepository_ChooseRoom(t *testing.T) {
	/*****************************************
	// first case -- reservation in session
	*****************************************/

	req, _ := http.NewRequest("GET", "/choose-room/1", nil)
	ctx := getCtx(req)
	req = req.WithContext(ctx)
	// set the RequestURI on the request so that we can grab the ID
	// from the URL
	req.RequestURI = "/choose-room/1"

	responseRecorder := httptest.NewRecorder()
	session.Put(ctx, "reservation", reservation)

	handler := http.HandlerFunc(Repo.ChooseRoom)

	handler.ServeHTTP(responseRecorder, req)

	if responseRecorder.Code != http.StatusSeeOther {
		t.Errorf("ChooseRoom handler returned wrong response code: got %d, wanted %d", responseRecorder.Code, http.StatusSeeOther)
	}

	///*****************************************
	//// second case -- reservation not in session
	//*****************************************/
	req, _ = http.NewRequest("GET", "/choose-room/1", nil)
	ctx = getCtx(req)
	req = req.WithContext(ctx)
	req.RequestURI = "/choose-room/1"

	responseRecorder = httptest.NewRecorder()

	handler = http.HandlerFunc(Repo.ChooseRoom)

	handler.ServeHTTP(responseRecorder, req)

	if responseRecorder.Code != http.StatusSeeOther {
		t.Errorf("ChooseRoom handler returned wrong response code: got %d, wanted %d", responseRecorder.Code, http.StatusSeeOther)
	}

	///*****************************************
	//// third case -- missing url parameter, or malformed parameter
	//*****************************************/
	req, _ = http.NewRequest("GET", "/choose-room/fish", nil)
	ctx = getCtx(req)
	req = req.WithContext(ctx)
	req.RequestURI = "/choose-room/fish"

	responseRecorder = httptest.NewRecorder()

	handler = http.HandlerFunc(Repo.ChooseRoom)

	handler.ServeHTTP(responseRecorder, req)

	if responseRecorder.Code != http.StatusSeeOther {
		t.Errorf("ChooseRoom handler returned wrong response code: got %d, wanted %d", responseRecorder.Code, http.StatusSeeOther)
	}
}

func TestRepository_BookRoom(t *testing.T) {
	/*****************************************
	// first case -- database works
	*****************************************/

	req, _ := http.NewRequest("GET", "/book-room?s=2050-01-01&e=2050-01-02&id=1", nil)
	ctx := getCtx(req)
	req = req.WithContext(ctx)

	rr := httptest.NewRecorder()
	session.Put(ctx, "reservation", reservation)

	handler := http.HandlerFunc(Repo.BookRoom)

	handler.ServeHTTP(rr, req)

	if rr.Code != http.StatusSeeOther {
		t.Errorf("BookRoom handler returned wrong response code: got %d, wanted %d", rr.Code, http.StatusSeeOther)
	}

	/*****************************************
	// second case -- database failed
	*****************************************/
	req, _ = http.NewRequest("GET", "/book-room?s=2040-01-01&e=2040-01-02&id=4", nil)
	ctx = getCtx(req)
	req = req.WithContext(ctx)

	rr = httptest.NewRecorder()

	handler = http.HandlerFunc(Repo.BookRoom)

	handler.ServeHTTP(rr, req)

	if rr.Code != http.StatusSeeOther {
		t.Errorf("BookRoom handler returned wrong response code: got %d, wanted %d", rr.Code, http.StatusSeeOther)
	}

	/*****************************************
	// third case -- wrong URL
	*****************************************/

	// Case 1 wrong: startDate
	req, _ = http.NewRequest("GET", "/book-room?s=invalidStartDate&e=2040-01-02&id=4", nil)
	ctx = getCtx(req)
	req = req.WithContext(ctx)

	rr = httptest.NewRecorder()

	handler = http.HandlerFunc(Repo.BookRoom)

	handler.ServeHTTP(rr, req)

	if rr.Code != http.StatusSeeOther {
		t.Errorf("BookRoom handler returned wrong response code: got %d, wanted %d", rr.Code, http.StatusSeeOther)
	}

	// Case 2 wrong: endDate
	req, _ = http.NewRequest("GET", "/book-room?s=2040-01-01&e=invalidEndDate&id=4", nil)
	ctx = getCtx(req)
	req = req.WithContext(ctx)

	rr = httptest.NewRecorder()

	handler = http.HandlerFunc(Repo.BookRoom)

	handler.ServeHTTP(rr, req)

	if rr.Code != http.StatusSeeOther {
		t.Errorf("BookRoom handler returned wrong response code: got %d, wanted %d", rr.Code, http.StatusSeeOther)
	}
	// Case 3 wrong: roomid
	req, _ = http.NewRequest("GET", "/book-room?s=2040-01-01&e=2040-01-02&id=invalidRoomID", nil)
	ctx = getCtx(req)
	req = req.WithContext(ctx)

	rr = httptest.NewRecorder()

	handler = http.HandlerFunc(Repo.BookRoom)

	handler.ServeHTTP(rr, req)

	if rr.Code != http.StatusSeeOther {
		t.Errorf("BookRoom handler returned wrong response code: got %d, wanted %d", rr.Code, http.StatusSeeOther)
	}

}

func TestNewRepo(t *testing.T) {
	var db driver.DB
	testRepo := NewRepo(&app, &db)

	if reflect.TypeOf(testRepo).String() != "*handlers.Repository" {
		t.Errorf("Did not get correct type from NewRepo: got %s, wanted *Repository", reflect.TypeOf(testRepo).String())
	}
}

func TestRepository_ReservationSummary(t *testing.T) {
	/*****************************************
	// first case -- reservation in session
	*****************************************/

	req, _ := http.NewRequest("GET", "/reservation-summary", nil)
	ctx := getCtx(req)
	req = req.WithContext(ctx)

	session.Put(ctx, "reservation", reservation)
	rr := httptest.NewRecorder()

	handler := http.HandlerFunc(Repo.ReservationSummary)

	handler.ServeHTTP(rr, req)

	if rr.Code != http.StatusOK {
		t.Errorf("ReservationSummary handler returned wrong response code: got %d, wanted %d", rr.Code, http.StatusOK)
	}

	/*****************************************
	// second case -- reservation not in session
	*****************************************/
	req, _ = http.NewRequest("GET", "/reservation-summary", nil)
	ctx = getCtx(req)
	req = req.WithContext(ctx)

	rr = httptest.NewRecorder()

	handler = http.HandlerFunc(Repo.ReservationSummary)

	handler.ServeHTTP(rr, req)

	if rr.Code != http.StatusSeeOther {
		t.Errorf("ReservationSummary handler returned wrong response code: got %d, wanted %d", rr.Code, http.StatusOK)
	}
}

func TestRepository_PostAvailability(t *testing.T) {
	/*****************************************
	// first case -- rooms are not available
	*****************************************/
	// create our request body
	// 2060-01-01 is available date
	postedData := url.Values{}
	postedData.Add("start", "2050-01-01")
	postedData.Add("end", "2050-01-01")

	// create our request
	req, _ := http.NewRequest("POST", "/search-availability", strings.NewReader(postedData.Encode()))

	// get the context with session
	ctx := getCtx(req)
	req = req.WithContext(ctx)

	// set the request header
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")

	// create our response recorder, which satisfies the requirements
	// for http.ResponseWriter
	rr := httptest.NewRecorder()

	// make our handler a http.HandlerFunc
	handler := http.HandlerFunc(Repo.PostAvailability)

	// make the request to our handler
	handler.ServeHTTP(rr, req)

	// since we have no rooms available, we expect to get status http.StatusSeeOther
	if rr.Code != http.StatusSeeOther {
		t.Errorf("Post availability when no rooms available gave wrong status code: got %d, wanted %d", rr.Code, http.StatusSeeOther)
	}

	/*****************************************
	// second case -- rooms are available
	*****************************************/
	// 2050-01-01 is not room available
	postedData = url.Values{}
	postedData.Add("start", "2060-01-01")
	postedData.Add("end", "2060-01-01")

	// create our request
	req, _ = http.NewRequest("POST", "/search-availability", strings.NewReader(postedData.Encode()))

	// get the context with session
	ctx = getCtx(req)
	req = req.WithContext(ctx)

	// set the request header
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")

	// create our response recorder, which satisfies the requirements
	// for http.ResponseWriter
	rr = httptest.NewRecorder()

	// make our handler a http.HandlerFunc
	handler = http.HandlerFunc(Repo.PostAvailability)

	// make the request to our handler
	handler.ServeHTTP(rr, req)

	// since we have rooms available, we expect to get status http.StatusOK
	if rr.Code != http.StatusOK {
		t.Errorf("Post availability when rooms are available gave wrong status code: got %d, wanted %d", rr.Code, http.StatusOK)
	}

	/*****************************************
	// third case -- empty post body
	*****************************************/
	// create our request with a nil body, so parsing form fails
	req, _ = http.NewRequest("POST", "/search-availability", nil)

	// get the context with session
	ctx = getCtx(req)
	req = req.WithContext(ctx)

	// set the request header
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")

	// create our response recorder, which satisfies the requirements
	// for http.ResponseWriter
	rr = httptest.NewRecorder()

	// make our handler a http.HandlerFunc
	handler = http.HandlerFunc(Repo.PostAvailability)

	// make the request to our handler
	handler.ServeHTTP(rr, req)

	// since we have rooms available, we expect to get status http.StatusSeeOther
	if rr.Code != http.StatusSeeOther {
		t.Errorf("Post availability with empty request body (nil) gave wrong status code: got %d, wanted %d", rr.Code, http.StatusSeeOther)
	}

	/*****************************************
	// fourth case -- start date in wrong format
	*****************************************/
	// this time, we specify a start date in the wrong format
	postedData = url.Values{}
	postedData.Add("start", "invalid")
	postedData.Add("end", "2040-01-02")
	req, _ = http.NewRequest("POST", "/search-availability", strings.NewReader(postedData.Encode()))

	// get the context with session
	ctx = getCtx(req)
	req = req.WithContext(ctx)

	// set the request header
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")

	// create our response recorder, which satisfies the requirements
	// for http.ResponseWriter
	rr = httptest.NewRecorder()

	// make our handler a http.HandlerFunc
	handler = http.HandlerFunc(Repo.PostAvailability)

	// make the request to our handler
	handler.ServeHTTP(rr, req)

	// since we have rooms available, we expect to get status http.StatusSeeOther
	if rr.Code != http.StatusSeeOther {
		t.Errorf("Post availability with invalid start date gave wrong status code: got %d, wanted %d", rr.Code, http.StatusSeeOther)
	}

	/*****************************************
	// fifth case -- end date in wrong format
	*****************************************/
	// this time, we specify a start date in the wrong format
	postedData = url.Values{}
	postedData.Add("start", "2040-01-01")
	postedData.Add("end", "invalid")

	req, _ = http.NewRequest("POST", "/search-availability", strings.NewReader(postedData.Encode()))

	// get the context with session
	ctx = getCtx(req)
	req = req.WithContext(ctx)

	// set the request header
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")

	// create our response recorder, which satisfies the requirements
	// for http.ResponseWriter
	rr = httptest.NewRecorder()

	// make our handler a http.HandlerFunc
	handler = http.HandlerFunc(Repo.PostAvailability)

	// make the request to our handler
	handler.ServeHTTP(rr, req)

	// since we have rooms available, we expect to get status http.StatusSeeOther
	if rr.Code != http.StatusSeeOther {
		t.Errorf("Post availability with invalid end date gave wrong status code: got %d, wanted %d", rr.Code, http.StatusSeeOther)
	}

	/*****************************************
	// sixth case -- database query fails
	*****************************************/
	// this time, we specify a start date of out of range (2099-12-31), which will cause
	// our testdb repo to return an error
	postedData = url.Values{}
	postedData.Add("start", "2100-01-01")
	postedData.Add("end", "2100-01-02")
	req, _ = http.NewRequest("POST", "/search-availability", strings.NewReader(postedData.Encode()))

	// get the context with session
	ctx = getCtx(req)
	req = req.WithContext(ctx)

	// set the request header
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")

	// create our response recorder, which satisfies the requirements
	// for http.ResponseWriter
	rr = httptest.NewRecorder()

	// make our handler a http.HandlerFunc
	handler = http.HandlerFunc(Repo.PostAvailability)

	// make the request to our handler
	handler.ServeHTTP(rr, req)

	// since we have rooms available, we expect to get status http.StatusSeeOther
	if rr.Code != http.StatusSeeOther {
		t.Errorf("Post availability when database query fails gave wrong status code: got %d, wanted %d", rr.Code, http.StatusSeeOther)
	}
}

var loginTests = []struct {
	name               string
	email              string
	expectedStatusCode int
	expectedHTML       string
	expectedLocation   string
}{
	// Hard code valid email in Authenticate(test-repo.go)
	{
		"valid-credentials",
		"validEmail@here.com",
		http.StatusSeeOther,
		"",
		"/",
	},
	{
		"invalid-credentials",
		"invalidEmail@here.com",
		http.StatusSeeOther,
		"",
		"/user/login",
	},
	{
		"invalid-data",
		"j",
		http.StatusOK,
		`action="/user/login"`,
		"",
	},
}

func TestLogin(t *testing.T) {
	for _, loginTest := range loginTests {
		// Create post data
		postData := url.Values{}
		postData.Add("email", loginTest.email)
		postData.Add("password", "password")

		// Create request
		req, _ := http.NewRequest("POST", "/user/login", strings.NewReader(postData.Encode()))
		ctx := getCtx(req)
		req = req.WithContext(ctx)

		// Set header
		req.Header.Set("Content-Type", "application/x-www-form-urlencoded")

		responseRecorder := httptest.NewRecorder()

		// Call the handlers
		handler := http.HandlerFunc(Repo.PostShowLogin)
		handler.ServeHTTP(responseRecorder, req)

		if responseRecorder.Code != loginTest.expectedStatusCode {
			t.Errorf("failed %s: expected code %d, but got %d", loginTest.name, loginTest.expectedStatusCode, responseRecorder.Code)
		}

		// Check the URL response from the handler
		if loginTest.expectedLocation != "" {
			// Get url from test
			actualLoc, _ := responseRecorder.Result().Location()
			if actualLoc.String() != loginTest.expectedLocation {
				t.Errorf("failed %s: expected location %s, but got %s", loginTest.name, loginTest.expectedLocation, actualLoc.String())
			}
		}

		// checking expected for value in HTML
		if loginTest.expectedHTML != "" {
			// Read response body into string
			html := responseRecorder.Body.String()

			if !strings.Contains(html, loginTest.expectedHTML) {
				t.Errorf("failed %s: expected location %s, but got %s", loginTest.name, loginTest.expectedHTML, html)
			}
		}
	}
}

var adminProcessReservationTests = []struct {
	name                 string
	queryParams          string
	expectedResponseCode int
	expectedLocation     string
}{
	{
		name:                 "process-reservation",
		queryParams:          "",
		expectedResponseCode: http.StatusSeeOther,
		expectedLocation:     "",
	},
	{
		name:                 "process-reservation-back-to-cal",
		queryParams:          "?y=2021&m=12",
		expectedResponseCode: http.StatusSeeOther,
		expectedLocation:     "",
	},
}

func TestAdminProcessReservation(t *testing.T) {
	for _, e := range adminProcessReservationTests {
		req, _ := http.NewRequest("GET", fmt.Sprintf("/admin/process-reservation/cal/1/do%s", e.queryParams), nil)
		ctx := getCtx(req)
		req = req.WithContext(ctx)

		rr := httptest.NewRecorder()

		handler := http.HandlerFunc(Repo.AdminProcessReservation)
		handler.ServeHTTP(rr, req)

		if rr.Code != http.StatusSeeOther {
			t.Errorf("failed %s: expected code %d, but got %d", e.name, e.expectedResponseCode, rr.Code)
		}
	}
}

var adminDeleteReservationTests = []struct {
	name                 string
	queryParams          string
	expectedResponseCode int
	expectedLocation     string
}{
	{
		name:                 "delete-reservation",
		queryParams:          "",
		expectedResponseCode: http.StatusSeeOther,
		expectedLocation:     "",
	},
	{
		name:                 "delete-reservation-back-to-cal",
		queryParams:          "?y=2021&m=12",
		expectedResponseCode: http.StatusSeeOther,
		expectedLocation:     "",
	},
}

func TestAdminDeleteReservation(t *testing.T) {
	for _, e := range adminDeleteReservationTests {
		req, _ := http.NewRequest("GET", fmt.Sprintf("/admin/process-reservation/cal/1/do%s", e.queryParams), nil)
		ctx := getCtx(req)
		req = req.WithContext(ctx)

		rr := httptest.NewRecorder()

		handler := http.HandlerFunc(Repo.AdminDeleteReservation)
		handler.ServeHTTP(rr, req)

		if rr.Code != http.StatusSeeOther {
			t.Errorf("failed %s: expected code %d, but got %d", e.name, e.expectedResponseCode, rr.Code)
		}
	}
}

var adminPostReservationCalendarTests = []struct {
	name                 string
	postedData           url.Values
	expectedResponseCode int
	expectedLocation     string
	expectedHTML         string
	blocks               int
	reservations         int
}{
	{
		name: "cal",
		postedData: url.Values{
			"year":  {time.Now().Format("2006")},
			"month": {time.Now().Format("01")},
			fmt.Sprintf("add_block_1_%s", time.Now().AddDate(0, 0, 2).Format("2006-01-2")): {"1"},
		},
		expectedResponseCode: http.StatusSeeOther,
	},
	{
		name:                 "cal-blocks",
		postedData:           url.Values{},
		expectedResponseCode: http.StatusSeeOther,
		blocks:               1,
	},
	{
		name:                 "cal-res",
		postedData:           url.Values{},
		expectedResponseCode: http.StatusSeeOther,
		reservations:         1,
	},
}

func TestAdminPostReservationCalendar(t *testing.T) {
	for _, e := range adminPostReservationCalendarTests {
		var req *http.Request
		if e.postedData != nil {
			req, _ = http.NewRequest("POST", "/admin/reservations-calendar", strings.NewReader(e.postedData.Encode()))
		} else {
			req, _ = http.NewRequest("POST", "/admin/reservations-calendar", nil)
		}
		ctx := getCtx(req)
		req = req.WithContext(ctx)

		now := time.Now()
		bm := make(map[string]int)
		rm := make(map[string]int)

		currentYear, currentMonth, _ := now.Date()
		currentLocation := now.Location()

		firstOfMonth := time.Date(currentYear, currentMonth, 1, 0, 0, 0, 0, currentLocation)
		lastOfMonth := firstOfMonth.AddDate(0, 1, -1)

		for d := firstOfMonth; d.After(lastOfMonth) == false; d = d.AddDate(0, 0, 1) {
			rm[d.Format("2006-01-2")] = 0
			bm[d.Format("2006-01-2")] = 0
		}

		if e.blocks > 0 {
			bm[firstOfMonth.Format("2006-01-2")] = e.blocks
		}

		if e.reservations > 0 {
			rm[lastOfMonth.Format("2006-01-2")] = e.reservations
		}

		session.Put(ctx, "block_map_1", bm)
		session.Put(ctx, "reservation_map_1", rm)

		// set the header
		req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
		rr := httptest.NewRecorder()

		// call the handler
		handler := http.HandlerFunc(Repo.AdminPostReservationsCalendar)
		handler.ServeHTTP(rr, req)

		if rr.Code != e.expectedResponseCode {
			t.Errorf("failed %s: expected code %d, but got %d", e.name, e.expectedResponseCode, rr.Code)
		}

	}
}

var adminPostShowReservationTests = []struct {
	name                 string
	url                  string
	postedData           url.Values
	expectedResponseCode int
	expectedLocation     string
	expectedHTML         string
}{
	{
		name: "valid-data-from-new",
		url:  "/admin/reservations/new/1/show",
		postedData: url.Values{
			"first_name": {"John"},
			"last_name":  {"Smith"},
			"email":      {"john@smith.com"},
			"phone":      {"555-555-5555"},
		},
		expectedResponseCode: http.StatusSeeOther,
		expectedLocation:     "/admin/reservations-new",
		expectedHTML:         "",
	},
	{
		name: "valid-data-from-all",
		url:  "/admin/reservations/all/1/show",
		postedData: url.Values{
			"first_name": {"John"},
			"last_name":  {"Smith"},
			"email":      {"john@smith.com"},
			"phone":      {"555-555-5555"},
		},
		expectedResponseCode: http.StatusSeeOther,
		expectedLocation:     "/admin/reservations-all",
		expectedHTML:         "",
	},
	{
		name: "valid-data-from-cal",
		url:  "/admin/reservations/cal/1/show",
		postedData: url.Values{
			"first_name": {"John"},
			"last_name":  {"Smith"},
			"email":      {"john@smith.com"},
			"phone":      {"555-555-5555"},
			"year":       {"2022"},
			"month":      {"01"},
		},
		expectedResponseCode: http.StatusSeeOther,
		expectedLocation:     "/admin/reservations-calendar?y=2022&m=01",
		expectedHTML:         "",
	},
}

// TestAdminPostShowReservations tests the AdminPostReservation handler
func TestAdminPostShowReservations(t *testing.T) {
	for _, e := range adminPostShowReservationTests {
		var req *http.Request
		if e.postedData != nil {
			req, _ = http.NewRequest("POST", "/user/login", strings.NewReader(e.postedData.Encode()))
		} else {
			req, _ = http.NewRequest("POST", "/user/login", nil)
		}
		ctx := getCtx(req)
		req = req.WithContext(ctx)
		req.RequestURI = e.url

		// set the header
		req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
		rr := httptest.NewRecorder()

		// call the handler
		handler := http.HandlerFunc(Repo.AdminPostShowReservations)
		handler.ServeHTTP(rr, req)

		if rr.Code != e.expectedResponseCode {
			t.Errorf("failed %s: expected code %d, but got %d", e.name, e.expectedResponseCode, rr.Code)
		}

		if e.expectedLocation != "" {
			// get the URL from test
			actualLoc, _ := rr.Result().Location()
			if actualLoc.String() != e.expectedLocation {
				t.Errorf("failed %s: expected location %s, but got location %s", e.name, e.expectedLocation, actualLoc.String())
			}
		}

		// checking for expected values in HTML
		if e.expectedHTML != "" {
			// read the response body into a string
			html := rr.Body.String()
			if !strings.Contains(html, e.expectedHTML) {
				t.Errorf("failed %s: expected to find %s but did not", e.name, e.expectedHTML)
			}
		}
	}
}
