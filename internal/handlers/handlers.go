package handlers

import (
	"encoding/json"
	"log"
	"net/http"
	"strconv"
	"time"

	"github.com/gdalelio/bookings/internal/config"
	"github.com/gdalelio/bookings/internal/driver"
	"github.com/gdalelio/bookings/internal/forms"
	"github.com/gdalelio/bookings/internal/helpers"
	"github.com/gdalelio/bookings/internal/models"
	"github.com/gdalelio/bookings/internal/render"
	"github.com/gdalelio/bookings/internal/repository"
	"github.com/gdalelio/bookings/internal/repository/dbrepo"
	"github.com/go-chi/chi/v5"
)

// Repo the repository used by the handlers
var Repo *Repository

// Repository Type
type Repository struct {
	App *config.AppConfig
	DB  repository.DatabaseRepo
}

// NewRepo creates a new repository with the app config and dbase connection pool
// returning a repo back
func NewRepo(a *config.AppConfig, dbase *driver.DB) *Repository {
	//creates a new repoository
	return &Repository{
		App: a,
		DB:  dbrepo.NewPostgresRepo(dbase.SQL, a),
	}
}

// NewHandlers sets the repository for the handlers
func NewHandlers(r *Repository) {
	Repo = r
}

// Home is the home page handler
func (m *Repository) Home(w http.ResponseWriter, r *http.Request) {

	render.Template(w, r, "home.page.tmpl", &models.TemplateData{})
}

// About is the about page handler
func (m *Repository) About(w http.ResponseWriter, r *http.Request) {
	//send the data to the template
	render.Template(w, r, "about.page.tmpl", &models.TemplateData{})
}

// Generals is the room page handler
func (m *Repository) Generals(w http.ResponseWriter, r *http.Request) {

	render.Template(w, r, "generals.page.tmpl", &models.TemplateData{})
}

// Majors is the room page handler
func (m *Repository) Majors(w http.ResponseWriter, r *http.Request) {

	render.Template(w, r, "majors.page.tmpl", &models.TemplateData{})
}

// CheckAvailability is the check availability page handler
func (m *Repository) Availability(w http.ResponseWriter, r *http.Request) {

	render.Template(w, r, "search-availability.page.tmpl", &models.TemplateData{})
}

// PostAvailability is the check availability page handler
func (m *Repository) PostAvailability(w http.ResponseWriter, r *http.Request) {
	//grabs strings for start and end
	start := r.Form.Get("start")
	end := r.Form.Get("end")

	//set the data format for parsing the date string into time.Time format
	layout := "2006-01-02"

	startDTParsed, err := time.Parse(layout, start)
	if err != nil {
		helpers.ServerError(w, err)
		return
	}

	endDTParsed, err := time.Parse(layout, end)
	if err != nil {
		helpers.ServerError(w, err)
		return
	}

	rooms, err := m.DB.SearchRoomAvailabilityForAllRooms(startDTParsed, endDTParsed)
	if err != nil {
		helpers.ServerError(w, err)
		return
	}

	//no rooms are available and redirecting to the search-avaialbility screen
	if len(rooms) == 0 {
		m.App.Session.Put(r.Context(), "error", "No avaialabilty")
		http.Redirect(w, r, "/search-availability", http.StatusSeeOther)
		return

	}
	//creting a map of strings that contains the room names
	data := make(map[string]interface{})
	//storing the rooms list into the data element to render to the page choose-room.page.tmpl
	data["rooms"] = rooms

	//create the start of a resrvation with start and end dates they searched
	reservation := models.Reservation{
		StartDate: startDTParsed,
		EndDate:   endDTParsed,
	}
	//store in the session for use with reservation
	m.App.Session.Put(r.Context(), "reservation", reservation)

	render.Template(w, r, "choose-room.page.tmpl", &models.TemplateData{
		Data: data,
	})

}

// jsonResponse keep as close to code that uses this type
type jsonResponse struct {
	OK      bool   `json:"ok"`
	Message string `message:"Available"`
}

// PostAvailabilityJSON handles request for availability and send JSON response
func (m *Repository) AvailabilityJSON(w http.ResponseWriter, r *http.Request) {
	resp := jsonResponse{
		OK:      false,
		Message: "Available!",
	}

	out, err := json.MarshalIndent(resp, "", "     ")
	if err != nil {
		helpers.ServerError(w, err)
		return
	}
	log.Println(string(out))
	w.Header().Set("Content-Type", "application/json")
	w.Write(out)

}

// Reservation is the check availability page handler
func (m *Repository) Reservation(w http.ResponseWriter, r *http.Request) {
	//create var to hold the empty reservation type
	var emptyReservation models.Reservation

	//create a map to store the entries on the form
	data := make(map[string]interface{})

	//store the data into the emptyReservation
	data["reservation"] = emptyReservation

	render.Template(w, r, "make-reservation.page.tmpl", &models.TemplateData{
		Form: forms.New(nil),
		Data: data,
	})

}

// PostReservation handles the posting of a reservation form
func (m *Repository) PostReservation(w http.ResponseWriter, r *http.Request) {
	err := r.ParseForm()
	if err != nil {
		helpers.ServerError(w, err)
		return
	}

	//dates from the form are strings and will need to be translated to date
	startDT := r.Form.Get("start_date")
	endDT := r.Form.Get("end_date")

	//golang data format - 2020-01-01 -- 01//02 03:04:05PM '06  -0700
	layout := "2006-01-02"

	log.Printf("\n starttDt before parsing: %s", startDT)
	log.Printf("\n endtDt before parsing: %s", endDT)

	startDTParsed, err := time.Parse(layout, startDT)
	if err != nil {
		helpers.ServerError(w, err)
		return
	}

	endDTParsed, err := time.Parse(layout, endDT)
	if err != nil { //check for eror
		helpers.ServerError(w, err)
		return
	}

	//convert the room id from as string to an int for reservation model
	roomID, err := strconv.Atoi(r.Form.Get("room_id"))
	if err != nil { //check for eror
		helpers.ServerError(w, err)
		return
	}

	reservation := models.Reservation{
		FirstName: r.Form.Get("first_name"),
		LastName:  r.Form.Get("last_name"),
		Phone:     r.Form.Get("phone"),
		Email:     r.Form.Get("email"),
		StartDate: startDTParsed,
		EndDate:   endDTParsed,
		RoomID:    roomID,
	}
	form := forms.New(r.PostForm)

	form.Required("first_name", "last_name", "email")
	form.MinLength("first_name", 3)
	form.IsEmail("email")
	form.IsPhone("phone")

	if !form.Valid() {
		//validated the form and creating map of data from the form
		data := make(map[string]interface{})
		//putting the reservation information into the data slice for use later
		data["reservation"] = reservation

		render.Template(w, r, "make-reservation.page.tmpl", &models.TemplateData{
			Form: form,
			Data: data,
		})
		return
	}

	//pass in the reservation to the model
	newReservationID, err := m.DB.InsertReservation(reservation)
	if err != nil {
		helpers.ServerError(w, err)
		return
	}

	restriction := models.RoomRestriction{
		StartDate:     startDTParsed,
		EndDate:       endDTParsed,
		RoomID:        roomID,
		ReservationID: newReservationID,
		RestrictionID: 1,
	}

	err = m.DB.InsertRoomRestriction(restriction)
	if err != nil {
		helpers.ServerError(w, err)
		return
	}

	//reservation summary presented
	m.App.Session.Put(r.Context(), "reservation", reservation)

	//do a redirect to prevent from double submission of the form  with Status code 300
	http.Redirect(w, r, "/reservation-summary", http.StatusSeeOther)
}

// Contact page handler
func (m *Repository) Contact(w http.ResponseWriter, r *http.Request) {
	//pulling IP and capturing it in the session model
	render.Template(w, r, "contact.page.tmpl", &models.TemplateData{})
}

// ReservationSummary
func (m *Repository) ReservationSummary(w http.ResponseWriter, r *http.Request) {
	//pulling reservation out of the session and casting it to models.Reservation
	reservation, ok := m.App.Session.Get(r.Context(), "reservation").(models.Reservation)
	if !ok {
		m.App.ErrorLog.Println("can't get error from session")
		m.App.Session.Put(r.Context(), "error", "can't get reservation from session")
		//redirect the user to the home page temporarily
		http.Redirect(w, r, "/", http.StatusTemporaryRedirect)
		return
	}
	//removing data from the reservation
	m.App.Session.Remove(r.Context(), "reservation")

	//map for reservation data
	data := make(map[string]interface{})
	//looking up the reservation using the "reservation" keu
	data["reservation"] = reservation

	render.Template(w, r, "reservation-summary.page.tmpl", &models.TemplateData{
		Data: data,
	})
}

// ChooseRoom displays the list of available
func (m *Repository) ChooseRoom(w http.ResponseWriter, r *http.Request) {
	//get the id from the room chosen on the choose-room page id 1 -> General's Room, id 2 -> Major's Room
	roomID, err := strconv.Atoi(chi.URLParam(r, "id"))
	if err != nil {
		helpers.ServerError(w, err)
		return
	}
	//get reservation variable in the session and stick the room id into the varaiable and put into reservation
	reservation, ok := m.App.Session.Get(r.Context(), "reservation").(models.Reservation)
	if !ok {
		helpers.ServerError(w, err)
		return
	}
	//putting the roomID into the reservation variable
	reservation.RoomID = roomID
	m.App.Session.Put(r.Context(), "reservation", reservation)

	http.Redirect(w, r, "/make-reservation", http.StatusSeeOther)
}
