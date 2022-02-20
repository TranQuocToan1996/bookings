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

type postData struct {
	key   string
	value string
}

// Test data is the slice of struct
var theTestsGET = []struct {
	testName         string
	url              string
	requestMethod    string
	expectStatusCode int
}{
	// GET method
	{"home", "/", "GET", http.StatusOK},
	{"about", "/about", "GET", http.StatusOK},
	{"generalsQuater", "/generals-quarters", "GET", http.StatusOK},
	{"majorsSuite", "/majors-suite", "GET", http.StatusOK},
	{"seachAvailability", "/search-availability", "GET", http.StatusOK},
	{"contact", "/contact", "GET", http.StatusOK},
	// {"makeReservation", "/make-reservation", "GET", []postData{}, http.StatusOK},

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
	if responseRecorder.Code != http.StatusTemporaryRedirect {
		t.Errorf("Reservation handler returned wrong code: Got %d, wanted %d", responseRecorder.Code, http.StatusTemporaryRedirect)
	}

	// Test with non-exist room
	req, _ = http.NewRequest("GET", "/make-reservation", nil)
	ctx = getCtx(req)
	req = req.WithContext(ctx)
	responseRecorder = httptest.NewRecorder()
	reservation.RoomID = 999999 // Test case for roomID that out of range
	session.Put(ctx, "reservation", reservation)
	handler.ServeHTTP(responseRecorder, req)
	if responseRecorder.Code != http.StatusTemporaryRedirect {
		t.Errorf("Reservation handler returned wrong code: Got %d, wanted %d", responseRecorder.Code, http.StatusTemporaryRedirect)
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
	if responseRecorder.Code != http.StatusTemporaryRedirect {
		t.Errorf("Reservation handler returned wrong code: Got %d, wanted %d", responseRecorder.Code, http.StatusTemporaryRedirect)
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
	if responseRecorder.Code != http.StatusTemporaryRedirect {
		t.Errorf("Reservation handler returned wrong code: Got %d, wanted %d", responseRecorder.Code, http.StatusTemporaryRedirect)
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
	rr := httptest.NewRecorder()
	session.Put(ctx, "reservation", reservation)
	handler = http.HandlerFunc(Repo.PostReservation)
	handler.ServeHTTP(rr, req)

	if rr.Code != http.StatusTemporaryRedirect {
		t.Errorf("Reservation handler returned wrong code: Got %d, wanted %d", responseRecorder.Code, http.StatusTemporaryRedirect)
	}

	/* Case 5: InsertReservation error*/
	// set up
	var reservation = models.Reservation{
		RoomID: 2,
	}

	// Hard code error Insert when roomID == 2 in session
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
	if responseRecorder.Code != http.StatusTemporaryRedirect {
		t.Errorf("Reservation handler returned wrong code: Got %d, wanted %d", responseRecorder.Code, http.StatusTemporaryRedirect)
	}

	/* Case 6: InsertRoomRestriction error - add later*/

	/* 	// setup

	   	// Hard code error Insert when restrictionID == 2 in session
	   	postedData = url.Values{}
	   	postedData.Add("start_date", "2050-01-01")
	   	postedData.Add("end_date", "2050-01-02")
	   	postedData.Add("first_name", "John")
	   	postedData.Add("last_name", "Smith")
	   	postedData.Add("email", "example@example.com")
	   	postedData.Add("phone", "0999999999")
	   	postedData.Add("room_id", "1000")

	   	req, _ = http.NewRequest("POST", "/make-reservation", strings.NewReader(postedData.Encode()))
	   	ctx = getCtx(req)
	   	req = req.WithContext(ctx)
	   	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	   	responseRecorder = httptest.NewRecorder()
	   	session.Put(ctx, "reservation", reservation) // Add reservation to session
	   	handler = http.HandlerFunc(Repo.PostReservation)
	   	handler.ServeHTTP(responseRecorder, req)
	   	if responseRecorder.Code != http.StatusTemporaryRedirect {
	   		t.Errorf("Reservation handler returned wrong code: Got %d, wanted %d", responseRecorder.Code, http.StatusTemporaryRedirect)
	   	} */

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
		t.Errorf("ChooseRoom handler returned wrong response code: got %d, wanted %d", responseRecorder.Code, http.StatusTemporaryRedirect)
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

	if responseRecorder.Code != http.StatusTemporaryRedirect {
		t.Errorf("ChooseRoom handler returned wrong response code: got %d, wanted %d", responseRecorder.Code, http.StatusTemporaryRedirect)
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

	if responseRecorder.Code != http.StatusTemporaryRedirect {
		t.Errorf("ChooseRoom handler returned wrong response code: got %d, wanted %d", responseRecorder.Code, http.StatusTemporaryRedirect)
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

	if rr.Code != http.StatusTemporaryRedirect {
		t.Errorf("BookRoom handler returned wrong response code: got %d, wanted %d", rr.Code, http.StatusTemporaryRedirect)
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

	if rr.Code != http.StatusTemporaryRedirect {
		t.Errorf("BookRoom handler returned wrong response code: got %d, wanted %d", rr.Code, http.StatusTemporaryRedirect)
	}

	// Case 2 wrong: endDate
	req, _ = http.NewRequest("GET", "/book-room?s=2040-01-01&e=invalidEndDate&id=4", nil)
	ctx = getCtx(req)
	req = req.WithContext(ctx)

	rr = httptest.NewRecorder()

	handler = http.HandlerFunc(Repo.BookRoom)

	handler.ServeHTTP(rr, req)

	if rr.Code != http.StatusTemporaryRedirect {
		t.Errorf("BookRoom handler returned wrong response code: got %d, wanted %d", rr.Code, http.StatusTemporaryRedirect)
	}
	// Case 3 wrong: roomid
	req, _ = http.NewRequest("GET", "/book-room?s=2040-01-01&e=2040-01-02&id=invalidRoomID", nil)
	ctx = getCtx(req)
	req = req.WithContext(ctx)

	rr = httptest.NewRecorder()

	handler = http.HandlerFunc(Repo.BookRoom)

	handler.ServeHTTP(rr, req)

	if rr.Code != http.StatusTemporaryRedirect {
		t.Errorf("BookRoom handler returned wrong response code: got %d, wanted %d", rr.Code, http.StatusTemporaryRedirect)
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

	if rr.Code != http.StatusTemporaryRedirect {
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
	if rr.Code != http.StatusTemporaryRedirect {
		t.Errorf("Post availability when no rooms available gave wrong status code: got %d, wanted %d", rr.Code, http.StatusTemporaryRedirect)
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

	// since we have rooms available, we expect to get status http.StatusTemporaryRedirect
	if rr.Code != http.StatusTemporaryRedirect {
		t.Errorf("Post availability with empty request body (nil) gave wrong status code: got %d, wanted %d", rr.Code, http.StatusTemporaryRedirect)
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

	// since we have rooms available, we expect to get status http.StatusTemporaryRedirect
	if rr.Code != http.StatusTemporaryRedirect {
		t.Errorf("Post availability with invalid start date gave wrong status code: got %d, wanted %d", rr.Code, http.StatusTemporaryRedirect)
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

	// since we have rooms available, we expect to get status http.StatusTemporaryRedirect
	if rr.Code != http.StatusTemporaryRedirect {
		t.Errorf("Post availability with invalid end date gave wrong status code: got %d, wanted %d", rr.Code, http.StatusTemporaryRedirect)
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

	// since we have rooms available, we expect to get status http.StatusTemporaryRedirect
	if rr.Code != http.StatusTemporaryRedirect {
		t.Errorf("Post availability when database query fails gave wrong status code: got %d, wanted %d", rr.Code, http.StatusTemporaryRedirect)
	}
}
