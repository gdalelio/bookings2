package models

import "time"

// User is the user model
type User struct {
	ID          int
	FirstName   string
	LastName    string
	Email       string
	Password    string
	AccessLevel int
	CreatedAt   time.Time
	UpdatedAt   time.Time
}

// Room is the rooms model
type Room struct {
	ID        int
	RoomName  string
	CreatedAt time.Time
	UpdatedAt time.Time
}

// Restriction is the restrictions model
type Restriction struct {
	ID               int
	RestrictionsName string
	CreatedAt        time.Time
	UpdatedAt        time.Time
}

// Reservation is the reservations model
type Reservation struct {
	FirstName string
	LastName  string
	Email     string
	Phone     string
	StartDate time.Time
	EndDate   time.Time
	RoomID    int
	CreatedAt time.Time
	UpdatedAt time.Time
	Room      Room
}

// RoomRestriction is the room restrictions model
type RoomRestriction struct {
	ID             int
	StartDate      time.Time
	EndDate        time.Time
	RoomID         int
	ReservationID int
	RestrictionID  int
	CreatedAt      time.Time
	UpdatedAt      time.Time
	Room           Room
	Reservation    Reservation
	Restriction    Restriction
}
