package dbrepo

import (
	"errors"
	"log"
	"time"

	"github.com/TranQuocToan1996/bookings/internal/models"
)

// Format time.Time
const layout string = "2006-01-02"
const finalDate string = "2099-12-31"

// implement for DatabaseRepo interface
func (t *testDBRepo) AllUsers() bool {
	return true
}

// InsertReservation inserts a reservation into database
func (t *testDBRepo) InsertReservation(res *models.Reservation) (int, error) {
	// if room id 2, then fail; otherwise, pass
	if res.RoomID == 2 {
		return 0, errors.New("some err")
	}
	return 1, nil
}

// InsertRoomRestriction inserts Room restriction data into database
func (t *testDBRepo) InsertRoomRestriction(r *models.RoomRestriction) error {
	if r.RoomID == 1000 {
		return errors.New("some err")
	}
	return nil
}

// SearchAvailabilityByDate checks availability of a specific room
func (t *testDBRepo) SearchAvailabilityByRoomID(start, end time.Time, roomID int) (bool, error) {

	startDate, err := time.Parse("2006-01-02", "2050-01-01")
	if err != nil {
		log.Println(err)
	}
	if start == startDate {
		return false, errors.New("some err")
	}

	// Logic to cancel booking After finalDate"
	finalDateParse, err := time.Parse(layout, finalDate)
	if err != nil {
		log.Println(err)
	}
	if start.After(finalDateParse) {
		return false, errors.New("out of date range")
	}

	return true, nil
}

// SearchAvailabilityForAllRooms returns a slice of available room(s) if any for given date range
func (t *testDBRepo) SearchAvailabilityForAllRooms(start, end time.Time) ([]models.Room, error) {

	var rooms []models.Room
	startDate, _ := time.Parse("2006-01-02", "2050-01-01")
	if start == startDate {
		return nil, errors.New("room not available")
	}

	// Logic to cancel booking After finalDate"
	finalDateParse, err := time.Parse(layout, finalDate)
	if err != nil {
		log.Println(err)
	}
	if start.After(finalDateParse) {
		return nil, errors.New("out of date range")
	}

	// otherwise, put an entry into the slice, indicating that some room is
	// available for search dates
	room := models.Room{
		ID: 1,
	}
	rooms = append(rooms, room)
	return rooms, nil
}

// GetRoomByID gets a room struct by id
func (t *testDBRepo) GetRoomByID(id int) (models.Room, error) {

	var room models.Room

	if id > 2 {
		return room, errors.New("some error")
	}

	return room, nil
}

func (t *testDBRepo) GetUserByID(id int) (models.User, error) {
	var u models.User
	return u, nil
}

func (t *testDBRepo) UpdateUser(u models.User) error {
	return nil
}

func (t *testDBRepo) Authenticate(email, testPassword string) (int, string, error) {
	if email == "validEmail@here.com" {

		return 1, "", nil
	}
	return 0, "", errors.New("error Authenticate in testing mode, successful testing")
}

// AllReservations returns a slice of all reservations (R in CRUD)
func (t *testDBRepo) AllReservations() ([]models.Reservation, error) {

	var reservations []models.Reservation

	return reservations, nil
}

// AllNewReservations returns a slice of all new reservations
func (t *testDBRepo) AllNewReservations() ([]models.Reservation, error) {

	var reservations []models.Reservation

	return reservations, nil
}

func (t *testDBRepo) GetReservationByID(id int) (models.Reservation, error) {

	var res models.Reservation

	return res, nil
}

// UpdateReservation updates the reservation info in the database
func (t *testDBRepo) UpdateReservation(r models.Reservation) error {

	return nil
}

// DeleteReservation deletes a reservation by ID
func (t *testDBRepo) DeleteReservation(id int) error {

	return nil
}

// UpdateProcessedForReservation updates processed-index in the database by id
func (t *testDBRepo) UpdateProcessedForReservation(id, processed int) error {

	return nil
}

func (t *testDBRepo) AllRooms() ([]models.Room, error) {
	var rooms []models.Room
	return rooms, nil
}

func (t *testDBRepo) GetRestrictionsForRoomByDate(roomID int, start, end time.Time) ([]models.RoomRestriction, error) {

	var restriction []models.RoomRestriction

	return restriction, nil
}

// InsertBlockForRoom inserts a room restriction
func (t *testDBRepo) InsertBlockForRoom(id int, startDate time.Time) error {

	return nil
}

// DeleteBlockByID deletes a room restriction
func (t *testDBRepo) DeleteBlockByID(id int) error {

	return nil
}
