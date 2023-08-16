package dbrepo

import (
	"context"
	"log"
	"time"

	"github.com/gdalelio/bookings/internal/models"
)

// AllUsers function with a model returned as a postgres repo
func (model *postgresDBRepo) AllUsers() bool {

	return true
}

// InsertReservation inserts a reservation into the database
func (m *postgresDBRepo) InsertReservation(reservation models.Reservation) (int, error) {

	//create a context for use with  execContext for executing the insert statment
	//uses time out to allow it to die after period of timme.
	//context.Background is available across the application
	contxt, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	var newID int //used to obtain the rowID of the inserted row

	//building SQL statment for database - with arguments instead of real values - avoids injected sql
	stmt := `insert into reservations (first_name, last_name, email, phone, start_date, 
			end_date,room_id, created_at, updated_at)
			values ($1,$2,$3,$4,$5,$6,$7,$8, $9) returning id`

	//log.Printf("\n %s", stmt)  <-------for debugging the insert statement
	//returns an result and error - only care about the case of an error; uses context
	//need the id of inserted id to use in the room_restrictions_id in newID

	err := m.DB.QueryRowContext(contxt, stmt,
		reservation.FirstName,
		reservation.LastName,
		reservation.Email,
		reservation.Phone,
		reservation.StartDate,
		reservation.EndDate,
		reservation.RoomID,
		time.Now(),
		time.Now(),
	).Scan(&newID)

	if err != nil {
		log.Println("returning error after trying to insert reservation!")
		return 0, err
	}
	return newID, nil
}

// InsertRoomRestriction inserts a room restriction into the database
// TODO need to roll back reservation insert if the room restriction isn't succesful - we have the reservation #
func (m *postgresDBRepo) InsertRoomRestriction(restriction models.RoomRestriction) error {
	contxt, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	stmt := `insert into room_restrictions (start_date,	end_date, room_id, reservation_id, 
			 created_at, updated_at,restriction_id)
			 values
			 ($1, $2, $3, $4, $5, $6, $7)`

	_, err := m.DB.ExecContext(contxt, stmt,
		restriction.StartDate,
		restriction.EndDate,
		restriction.RoomID,
		restriction.ReservationID,
		time.Now(),
		time.Now(),
		restriction.RestrictionID,
	)
	if err != nil {
		log.Println("returning error after trying to insert roomRestriction row!")
		return err

	}

	return nil

}

// SearchAvailabilityByDatesByRoomID returns true if availability exists for room id, and false if no availability exists
func (m *postgresDBRepo) SearchAvailabilityByDatesByRoomID(startDT, endDT time.Time, roomID int) (bool, error) {
	contxt, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	//number of rows returned from query
	var numRows int

	//room id = $1, start date = $2, and end date is $3
	query := `
		select 
	   		count (id) 
		from 
			room_restrictions 
		where 
			room_id = $1
		 	and $2 < end_date and $3 > start_date;`

	row := m.DB.QueryRowContext(contxt, query, startDT, endDT)
	err := row.Scan(&numRows)

	if err != nil {
		log.Println("SearchAvailablityByDates failed")
		return false, err
	}
	if numRows == 0 {
		return true, nil
	}
	return false, nil
}

// SearchRoomAvailabilityForAllRooms returns a slice of available roooms if any, for a range of dates
func (m *postgresDBRepo) SearchRoomAvailabilityForAllRooms(startDT, endDT time.Time) ([]models.Room, error) {
	contxt, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	var rooms []models.Room

	// checking for availabilty of any rooms that have available dates; start_date = $1 and end_date =$2
	query := `
			Select
				r.id, r.room_name
			from
				rooms r
			where r.id not in 
			(select rr.room_id from room_restrictions rr where $1 < rr.end_date and $2 > rr.start_date);
			`
	rows, err := m.DB.QueryContext(contxt, query, startDT, endDT)
	if err != nil {
		return rooms, err
	}

	for rows.Next() {
		var room models.Room
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

// GetRoomByID returns the name of the room for the room id
func (m *postgresDBRepo) GetRoomByID(id int) (models.Room, error) {
	contxt, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	var room models.Room

	log.Printf("room id = %d", id)
	query := `
		select id, room_name, created_at, updated_at from rooms where id = $1
	`
	row := m.DB.QueryRowContext(contxt, query, id)
	err := row.Scan(
		&room.ID,
		&room.RoomName,
		&room.CreatedAt,
		&room.UpdatedAt,
	)
	if err != nil {
		return room, err
	}
	return room, nil
}
