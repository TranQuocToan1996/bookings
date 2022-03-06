package dbrepo

import (
	"context"
	"errors"
	"log"
	"time"

	"github.com/TranQuocToan1996/bookings/internal/models"
	"golang.org/x/crypto/bcrypt"
)

// implement for DatabaseRepo interface
func (p *postgresDBRepo) AllUsers() bool {
	return true
}

// InsertReservation inserts a reservation into database (C in CRUD)
func (p *postgresDBRepo) InsertReservation(res *models.Reservation) (int, error) {

	// Context, if for any reasons that Insert not complete within 3 seconds, cancel connection
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	// Insert post data into database and returning reservation id
	query := `insert into reservations
	(first_name, last_name, email, phone, start_date, end_date, room_id, created_at, updated_at) 
	values  ($1, $2, $3, $4, $5, $6, $7, $8, $9) returning id`
	var newID int
	err := p.DB.QueryRowContext(ctx, query,
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
func (p *postgresDBRepo) InsertRoomRestriction(r *models.RoomRestriction) error {
	// Context, if for any reasons that Insert not complete within 3 seconds, cancel connection
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	query := `insert into	room_restriction 
	(start_date, end_date, room_id, reservation_id , created_at, updated_at, restriction_id)
	values  ($1, $2, $3, $4, $5, $6, $7)`
	_, err := p.DB.ExecContext(ctx, query,
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
func (p *postgresDBRepo) SearchAvailabilityByRoomID(start, end time.Time, roomID int) (bool, error) {

	// Context, if for any reasons that Insert not complete within 3 seconds, cancel connection
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	var numRows int
	query := `select 
						count(id) 
					from 
						room_restriction
					where 
						$1 < end_date and $2 > start_date
						and room_id = $3;`
	err := p.DB.QueryRowContext(ctx, query,
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
func (p *postgresDBRepo) SearchAvailabilityForAllRooms(start, end time.Time) ([]models.Room, error) {
	// Context, if for any reasons that Insert not complete within 3 seconds, cancel connection
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	var rooms []models.Room
	query := `select 
						r.id, r.room_name 
					from
						rooms r
					where r.id not in (
							select room_id from room_restriction rr
							where $1 < rr.end_date and $2 > rr.start_date
					)`
	rows, err := p.DB.QueryContext(ctx, query,
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
func (p *postgresDBRepo) GetRoomByID(id int) (models.Room, error) {
	// Context, if for any reasons that Insert not complete within 3 seconds, cancel connection
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	var room models.Room
	query := `select id, room_name, created_at, updated_at from rooms	where id = $1`
	row := p.DB.QueryRowContext(ctx, query,
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

// GetUserByID return the user information by ID
func (p *postgresDBRepo) GetUserByID(id int) (models.User, error) {
	// Context, if for any reasons that Insert not complete within 3 seconds, cancel connection
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	query := `select id, first_name, last_name, email, password, access_level, created_at, updated_at
				from users where id=$1`
	var u models.User
	err := p.DB.QueryRowContext(ctx, query, id).Scan(
		&u.ID,
		&u.FirstName,
		&u.LastName,
		&u.Email,
		&u.Password,
		&u.AccessLevel,
		&u.CreateAt,
		&u.UpdateAt,
	)
	if err != nil {
		return u, err
	}

	return u, nil
}

// UpdateUser updates a user in the database
func (p *postgresDBRepo) UpdateUser(u models.User) error {
	// Context, if for any reasons that Insert not complete within 3 seconds, cancel connection
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	query := `update users set first_name=$1, last_name=$2, email=$3, access_level=$4, updated_at=$5`
	_, err := p.DB.ExecContext(ctx, query,
		u.FirstName,
		u.LastName,
		u.Email,
		u.AccessLevel,
		time.Now(),
	)
	if err != nil {
		return err
	}

	return nil
}

// Authenticate authenticates the user and send back user_id, hashPassword, and an error if any
func (p *postgresDBRepo) Authenticate(email, testPassword string) (int, string, error) {
	// Context, if for any reasons that Insert not complete within 3 seconds, cancel connection
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	var id int
	var hashedPassword string // This one will store in the database instead of plain text password

	row := p.DB.QueryRowContext(ctx, "select id, password from users where email=$1", email)
	err := row.Scan(&id, &hashedPassword)
	if err != nil {
		return id, "", err
	}

	// Compare between hashedPassword and password user typing in the form
	err = bcrypt.CompareHashAndPassword([]byte(hashedPassword), []byte(testPassword))
	if err == bcrypt.ErrMismatchedHashAndPassword {
		return 0, "", errors.New("incorrect password")
	} else if err != nil {
		return 0, "", err
	}

	return id, hashedPassword, nil
}

// AllReservations returns a slice of all reservations
func (p *postgresDBRepo) AllReservations() ([]models.Reservation, error) {
	// Context, if for any reasons that Insert not complete within 3 seconds, cancel connection
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	var reservations []models.Reservation
	query := `
			select r.id, r.first_name, r.last_name, r.email, r.phone, 
			r.start_date, r.end_date, r.room_id, r.created_at, r.updated_at, r.processed,
			rm.id, rm.room_name

			from reservations r
			left join rooms rm on (r.room_id = rm.id)
			order by r.start_date asc
	`

	rows, err := p.DB.QueryContext(ctx, query)
	if err != nil {
		return reservations, err
	}
	defer rows.Close()

	for rows.Next() {
		var item models.Reservation
		err := rows.Scan(
			&item.ID,
			&item.FirstName,
			&item.LastName,
			&item.Email,
			&item.Phone,
			&item.StartDate,
			&item.EndDate,
			&item.RoomID,
			&item.CreateAt,
			&item.UpdateAt,
			&item.Processed,

			&item.Room.ID,
			&item.Room.RoomName,
		)
		if err != nil {
			return reservations, err
		}
		reservations = append(reservations, item)
	}

	if err = rows.Err(); err != nil {
		return reservations, err
	}

	return reservations, nil
}

// AllNewReservations returns a slice of all new reservations
func (p *postgresDBRepo) AllNewReservations() ([]models.Reservation, error) {
	// Context, if for any reasons that Insert not complete within 3 seconds, cancel connection
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	var reservations []models.Reservation
	query := `
			select r.id, r.first_name, r.last_name, r.email, r.phone, 
			r.start_date, r.end_date, r.room_id, r.created_at, r.updated_at,
			rm.id, rm.room_name

			from reservations r
			left join rooms rm on (r.room_id = rm.id)
			where r.processed = 0
			order by r.start_date asc
	`

	rows, err := p.DB.QueryContext(ctx, query)
	if err != nil {
		return reservations, err
	}
	defer rows.Close()

	for rows.Next() {
		var item models.Reservation
		err := rows.Scan(
			&item.ID,
			&item.FirstName,
			&item.LastName,
			&item.Email,
			&item.Phone,
			&item.StartDate,
			&item.EndDate,
			&item.RoomID,
			&item.CreateAt,
			&item.UpdateAt,

			&item.Room.ID,
			&item.Room.RoomName,
		)
		if err != nil {
			return reservations, err
		}
		reservations = append(reservations, item)
	}

	if err = rows.Err(); err != nil {
		return reservations, err
	}

	return reservations, nil
}

func (p *postgresDBRepo) GetReservationByID(id int) (models.Reservation, error) {
	// Context, if for any reasons that Insert not complete within 3 seconds, cancel connection
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	var res models.Reservation
	query := `
			select r.id, r.first_name, r.last_name, r.email, r.phone, 
			r.start_date, r.end_date, r.room_id, r.created_at, r.updated_at, r.processed,
			rm.id, rm.room_name
			from reservations r
			left join rooms rm on (r.room_id = rm.id)
			where r.id = $1
	`

	row := p.DB.QueryRowContext(ctx, query, id)
	err := row.Scan(
		&res.ID,
		&res.FirstName,
		&res.LastName,
		&res.Email,
		&res.Phone,
		&res.StartDate,
		&res.EndDate,
		&res.RoomID,
		&res.CreateAt,
		&res.UpdateAt,
		&res.Processed,

		&res.Room.ID,
		&res.Room.RoomName,
	)

	if err != nil {
		return res, err
	}

	return res, nil
}

// UpdateReservation updates the reservation info in the database
func (p *postgresDBRepo) UpdateReservation(r models.Reservation) error {
	// Context, if for any reasons that Insert not complete within 3 seconds, cancel connection
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	query := `update reservations set first_name=$1, last_name=$2, email=$3, phone=$4, updated_at=$5
				where id = $6`
	_, err := p.DB.ExecContext(ctx, query,
		r.FirstName,
		r.LastName,
		r.Email,
		r.Phone,
		time.Now(),
		r.ID,
	)
	if err != nil {
		return err
	}

	return nil
}

// DeleteReservation deletes a reservation by ID
func (p *postgresDBRepo) DeleteReservation(id int) error {
	// Context, if for any reasons that Insert not complete within 3 seconds, cancel connection
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	query := `delete from reservations where id = $1`
	_, err := p.DB.ExecContext(ctx, query, id)
	if err != nil {
		return err
	}

	return nil
}

// UpdateProcessedForReservation updates processed-index in the database by id
func (p *postgresDBRepo) UpdateProcessedForReservation(id, processed int) error {
	// Context, if for any reasons that Insert not complete within 3 seconds, cancel connection
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	query := `update reservations set processed = $1 where id = $2`
	_, err := p.DB.ExecContext(ctx, query, processed, id)
	if err != nil {
		return err
	}

	return nil
}

// AllRooms returns all room from the database
func (p *postgresDBRepo) AllRooms() ([]models.Room, error) {
	// Context, if for any reasons that Insert not complete within 3 seconds, cancel connection
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	var rooms []models.Room
	query := `select id, room_name, created_at, updated_at from rooms order by room_name`

	rows, err := p.DB.QueryContext(ctx, query)
	if err != nil {
		return rooms, err
	}
	defer rows.Close()

	for rows.Next() {
		var rm models.Room
		err := rows.Scan(
			&rm.ID,
			&rm.RoomName,
			&rm.CreateAt,
			&rm.UpdateAt,
		)
		if err != nil {
			return rooms, err
		}
		rooms = append(rooms, rm)
	}

	if err = rows.Err(); err != nil {
		return rooms, err
	}

	return rooms, nil
}

func (p *postgresDBRepo) GetRestrictionsForRoomByDate(roomID int, start, end time.Time) ([]models.RoomRestriction, error) {
	// Context, if for any reasons that Insert not complete within 3 seconds, cancel connection
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	var restriction []models.RoomRestriction
	// coalesce: if reservation_id is null using 0 instead
	query := `select id, coalesce(reservation_id, 0), restriction_id, room_id, start_date, end_date
	from room_restriction where $1 < end_date and $2 > start_date and room_id = $3
	`

	rows, err := p.DB.QueryContext(ctx, query, start, end, roomID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		var r models.RoomRestriction
		err := rows.Scan(
			&r.ID,
			&r.ReservationID,
			&r.RestrictionID,
			&r.RoomID,
			&r.StartDate,
			&r.EndDate,
		)
		if err != nil {
			return nil, err
		}
		restriction = append(restriction, r)
	}
	if err = rows.Err(); err != nil {
		return nil, err
	}

	return restriction, nil
}

// InsertBlockForRoom inserts a room restriction
func (p *postgresDBRepo) InsertBlockForRoom(id int, startDate time.Time) error {
	// Context, if for any reasons that Insert not complete within 3 seconds, cancel connection
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	query := `insert into room_restriction (start_date, end_date, room_id, restriction_id, created_at, updated_at)
				values ($1, $2, $3, $4, $5, $6)
	`

	_, err := p.DB.ExecContext(ctx, query, startDate, startDate.AddDate(0, 0, 1), id, 2, time.Now(), time.Now())
	if err != nil {
		log.Println(err)
		return err
	}

	return nil
}

// DeleteBlockByID deletes a room restriction
func (p *postgresDBRepo) DeleteBlockByID(id int) error {
	// Context, if for any reasons that Insert not complete within 3 seconds, cancel connection
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	query := `delete from room_restriction where id = $1`

	_, err := p.DB.ExecContext(ctx, query, id)
	if err != nil {
		log.Println(err)
		return err
	}

	return nil
}
