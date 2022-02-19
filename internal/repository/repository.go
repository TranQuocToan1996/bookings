package repository

import (
	"time"

	"github.com/TranQuocToan1996/bookings/internal/models"
)

// Contains method to contact with table in database
type DatabaseRepo interface {
	AllUsers() bool

	InsertReservation(res *models.Reservation) (int, error)

	InsertRoomRestriction(r *models.RoomRestriction) error

	SearchAvailabilityByRoomID(start, end time.Time, roomID int) (bool, error)

	SearchAvailabilityForAllRooms(start, end time.Time) ([]models.Room, error)

	GetRoomByID(id int) (models.Room, error)
}
