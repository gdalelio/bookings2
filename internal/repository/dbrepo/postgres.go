package dbrepo

import (
	"context"
	"log"
	"time"

	"github.com/gdalelio/bookings/internal/models"
)

// AllUsers function with a mdoel returned as a postgres repo
func (model *postgresDBRepo) AllUsers() bool {

	return true
}

// InsertReservation inserts a reservation into the database
func (model *postgresDBRepo) InsertReservation(reservation models.Reservation) (int, error) {

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

	err := model.DB.QueryRowContext(contxt, stmt,
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
func (model *postgresDBRepo) InsertRoomRestriction(restriction models.RoomRestriction) error {
	contxt, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	stmt := `insert into room_restrictions (start_date,	end_date, room_id, reservation_id, 
			 created_at, updated_at,restriction_id)
			 values
			 ($1, $2, $3, $4, $5, $6, $7)`

	_, err := model.DB.ExecContext(contxt, stmt,
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

//SearchAvailabilityByDates returns true if availability exists for room id, and false if no availability exists
func (model *postgresDBRepo) SearchAvailabilityByDates(startDT, endDT time.Time, roomID int) (bool, error) {
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

	row := model.DB.QueryRowContext(contxt, query, startDT, endDT)
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


//Rooms Available for the dates provided in search availability