package models

import "time"

// User struct hold information for users model
type User struct {
	ID          int
	FirstName   string
	LastName    string
	Email       string
	Password    string
	AccessLevel int
	CreateAt    time.Time
	UpdateAt    time.Time
}

// Room is the rooms model
type Room struct {
	ID       int
	RoomName string
	CreateAt time.Time
	UpdateAt time.Time
}

// Restriction is the restrictions model
type Restriction struct {
	ID              int
	RestrictionName string
	CreateAt        time.Time
	UpdateAt        time.Time
}

// Revervation is the Revervations model
type Reservation struct {
	ID        int
	FirstName string
	LastName  string
	Email     string
	Phone     string
	StartDate time.Time
	EndDate   time.Time
	RoomID    int
	CreateAt  time.Time
	UpdateAt  time.Time
	Room      Room
}

// RoomRestriction is the RoomRestriction model
type RoomRestriction struct {
	ID            int
	StartDate     time.Time
	EndDate       time.Time
	RoomID        int
	ReservationID int
	RestrictionID int
	CreateAt      time.Time
	UpdateAt      time.Time
	Room          Room
	Reservation   Reservation
	Restriction   Restriction
}
