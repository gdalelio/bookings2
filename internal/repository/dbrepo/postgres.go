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
func (model *postgresDBRepo) InsertReservation(reservation models.Reservation) error {

	//create a context for use with  execContext for executing the insert statment
	//uses time out to allow it to die after period of timme.
	//context.Background is available across the application
	contxt, cancel := context.WithTimeout(context.Background(), 3*time.Second)

	defer cancel()

	//building SQL statment for database - with arguments instead of real values - avoids injected sql
	stmt := `insert into reservations (first_name, last_name, email, phone, start_date, 
			end_date,  created_at, updated_at,room_id)
			values ($1,$2,$3,$4,$5,$6,$7,$8, $9)`
	//log.Printf("\n %s", stmt)  <-------for debugging the insert statement
	//returns an result and error - only care about the case of an error; uses context
	_, err := model.DB.ExecContext(contxt, stmt,
		reservation.FirstName,
		reservation.LastName,
		reservation.Email,
		reservation.Phone,
		reservation.StartDate,
		reservation.EndDate,
		time.Now(),
		time.Now(),
		reservation.RoomID,
	)
	if err != nil {
		log.Println("returning error after trying to insert")
		return err
	}
	return nil
}
