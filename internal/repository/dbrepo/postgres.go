package dbrepo

import (
	"context"
	"time"

	"github.com/gdalelio/bookings/internal/models"
)

// AllUsers function with a mdoel returned as a postgres repo
func (model *postgresDBRepo) AllUsers() bool {

	return true
}

// InsertReservation inserts a reservation into the database
func (model *postgresDBRepo) InsertReservation(reservation models.Reservation) error {

	//create a context for use with  execContext for executing the insert statment
	//uses time out to allow it to die after period of timme.  
	//context.Background is available across the application
	contxt, cancel := context.WithTimeout(context.Background(), 3*time.Second)

	defer cancel()

	//building SQL statment for database - with arguments instead of real values - avoids injected sql
	stmt := `insert into reservations (first_name, last_name, email, phone, start_date, 
			end_date, room_id, created_at, updated_at)
			values ($1,$2,$3,$4,$5,$6,$7,$8, $9)`

	//returns an result and error - only care about the case of an error; uses context 
	_, err := model.DB.ExecContext(contxt, stmt,
		reservation.FirstName,
		reservation.LastName,
		reservation.Phone,
		reservation.StartDate,
		reservation.EndDate,
		reservation.RoomID,
		//use the time now for the created_at and updated_at as they only need time
		time.Now(),
		time.Now(),
	)
	if err != nil {
		return err
	}
	return nil
}
