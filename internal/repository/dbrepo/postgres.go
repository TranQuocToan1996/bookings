package dbrepo

import (
	"context"
	"time"

	"github.com/TranQuocToan1996/bookings/internal/models"
)

// implement for DatabaseRepo interface
func (m *postgresDBRepo) AllUsers() bool {
	return true
}

// InsertReservation inserts a reservation into database
func (m *postgresDBRepo) InsertReservation(res *models.Reservation) (int, error) {

	// Context, if for any reasons that Insert not complete within 3 seconds, cancel connection
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	// Insert post data into database and returning reservation id
	queryString := `insert into reservations
	(first_name, last_name, email, phone, start_date, end_date, room_id, created_at, updated_at) 
	values  ($1, $2, $3, $4, $5, $6, $7, $8, $9) returning id`
	var newID int
	err := m.DB.QueryRowContext(ctx, queryString,
		res.FirstName,
		res.LastName,
		res.Email,
		res.Phone,
		res.StartDate,
		res.EndDate,
		res.RoomID,
		time.Now(),
		time.Now(),
	).Scan(&newID)
	if err != nil {
		return 0, err
	}

	return newID, nil
}

// InsertRoomRestriction inserts Room restriction data into database
func (m *postgresDBRepo) InsertRoomRestriction(r *models.RoomRestriction) error {
	// Context, if for any reasons that Insert not complete within 3 seconds, cancel connection
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	queryString := `insert into	room_restriction 
	(start_date, end_date, room_id, reservation_id , created_at, updated_at, restriction_id)
	values  ($1, $2, $3, $4, $5, $6, $7)`
	_, err := m.DB.ExecContext(ctx, queryString,
		r.StartDate,
		r.EndDate,
		r.RoomID,
		r.ReservationID,
		time.Now(),
		time.Now(),
		r.RestrictionID,
	)
	if err != nil {
		return err
	}

	return nil
}

// SearchAvailabilityByDate checks availability of a specific room
func (m *postgresDBRepo) SearchAvailabilityByRoomID(start, end time.Time, roomID int) (bool, error) {

	// Context, if for any reasons that Insert not complete within 3 seconds, cancel connection
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	var numRows int
	queryString := `select 
						count(id) 
					from 
						room_restriction
					where 
						$1 < end_date and $2 > start_date
						and room_id = $3;`
	err := m.DB.QueryRowContext(ctx, queryString,
		start,
		end,
		roomID,
	).Scan(&numRows)
	if err != nil {
		return false, err
	}

	return numRows == 0, nil
}

// SearchAvailabilityForAllRooms returns a slice of available room(s) if any for given date range
func (m *postgresDBRepo) SearchAvailabilityForAllRooms(start, end time.Time) ([]models.Room, error) {
	// Context, if for any reasons that Insert not complete within 3 seconds, cancel connection
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	var rooms []models.Room
	queryString := `select 
						r.id, r.room_name 
					from
						rooms r
					where r.id not in (
							select room_id from room_restriction rr
							where $1 < rr.end_date and $2 > rr.start_date
					)`
	rows, err := m.DB.QueryContext(ctx, queryString,
		start,
		end,
	)
	if err != nil {
		return rooms, err
	}

	for rows.Next() {
		room := models.Room{}
		err := rows.Scan(
			&room.ID,
			&room.RoomName,
		)
		if err != nil {
			return rooms, err
		}

		rooms = append(rooms, room)
	}

	if err = rows.Err(); err != nil {
		return rooms, err
	}

	return rooms, nil
}

// GetRoomByID gets a room struct by id
func (m *postgresDBRepo) GetRoomByID(id int) (models.Room, error) {
	// Context, if for any reasons that Insert not complete within 3 seconds, cancel connection
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	var room models.Room
	queryString := `select id, room_name, created_at, updated_at from rooms	where id = $1`
	row := m.DB.QueryRowContext(ctx, queryString,
		id,
	)
	err := row.Scan(
		&room.ID,
		&room.RoomName,
		&room.CreateAt,
		&room.UpdateAt,
	)
	if err != nil {
		return room, err
	}

	return room, nil
}
