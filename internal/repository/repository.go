package repository

import "github.com/gdalelio/bookings/internal/models"

//DatabaseRepo sets up an interface for operations on ther database
type DatabaseRepo interface {

	AllUsers() bool

	InsertReservation(res models.Reservation) error
}

