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
func (m *testDBRepo) AllUsers() bool {
	return true
}

// InsertReservation inserts a reservation into database
func (m *testDBRepo) InsertReservation(res *models.Reservation) (int, error) {
	// if room id 2, then fail; otherwise, pass
	if res.RoomID == 2 {
		return 0, errors.New("some err")
	}
	return 1, nil
}

// InsertRoomRestriction inserts Room restriction data into database
func (m *testDBRepo) InsertRoomRestriction(r *models.RoomRestriction) error {
	if r.RestrictionID == 2 {
		return errors.New("some err")
	}
	return nil
}

// SearchAvailabilityByDate checks availability of a specific room
func (m *testDBRepo) SearchAvailabilityByRoomID(start, end time.Time, roomID int) (bool, error) {

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
		return false, errors.New("Out of date range")
	}

	return true, nil
}

// SearchAvailabilityForAllRooms returns a slice of available room(s) if any for given date range
func (m *testDBRepo) SearchAvailabilityForAllRooms(start, end time.Time) ([]models.Room, error) {

	var rooms []models.Room
	startDate, _ := time.Parse("2006-01-02", "2050-01-01")
	if start == startDate {
		return nil, errors.New("room not available!")
	}

	// Logic to cancel booking After finalDate"
	finalDateParse, err := time.Parse(layout, finalDate)
	if err != nil {
		log.Println(err)
	}
	if start.After(finalDateParse) {
		return nil, errors.New("Out of date range")
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
func (m *testDBRepo) GetRoomByID(id int) (models.Room, error) {

	var room models.Room

	if id > 2 {
		return room, errors.New("Some error!")
	}

	return room, nil
}
