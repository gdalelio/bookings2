package repository

import (
	"time"

	"github.com/gdalelio/bookings/internal/models"
)

// DatabaseRepo sets up an interface for operations on the database
type DatabaseRepo interface {
	AllUsers() bool

	InsertReservation(reservation models.Reservation) (int, error)

	//TODO - need to test to see if the room restriction fails to roll back the reservation
	//otherwise there is an extra row in reservations that don't match the room restrictions
	InsertRoomRestriction(restriction models.RoomRestriction) error

	SearchAvailabilityByRoomID(startDT, endDT time.Time, roomID int) (bool, error)
}
