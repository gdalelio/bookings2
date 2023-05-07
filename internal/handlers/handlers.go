package handlers

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"

	"github.com/gdalelio/bookings/internal/config"
	"github.com/gdalelio/bookings/internal/forms"
	"github.com/gdalelio/bookings/internal/helpers"
	"github.com/gdalelio/bookings/internal/models"
	"github.com/gdalelio/bookings/internal/render"
)

// Repo the repository used by the handlers
var Repo *Repository

// Repository Type
type Repository struct {
	App *config.AppConfig
}

// NewRepo creates a new repository
func NewRepo(a *config.AppConfig) *Repository {
	return &Repository{
		App: a,
	}
}

// NewHandlers sets the repository for the handlers
func NewHandlers(r *Repository) {
	Repo = r
}

// Home is the home page handler
func (m *Repository) Home(w http.ResponseWriter, r *http.Request) {
	render.RenderTemplate(w, r, "home.page.tmpl", &models.TemplateData{})
}

// About is the about page handler
func (m *Repository) About(w http.ResponseWriter, r *http.Request) {
	//send the data to the template
	render.RenderTemplate(w, r, "about.page.tmpl", &models.TemplateData{})
}

// Generals is the room page handler
func (m *Repository) Generals(w http.ResponseWriter, r *http.Request) {

	render.RenderTemplate(w, r, "generals.page.tmpl", &models.TemplateData{})
}

// Majors is the room page handler
func (m *Repository) Majors(w http.ResponseWriter, r *http.Request) {

	render.RenderTemplate(w, r, "majors.page.tmpl", &models.TemplateData{})
}

// CheckAvailability is the check availability page handler
func (m *Repository) Availability(w http.ResponseWriter, r *http.Request) {

	render.RenderTemplate(w, r, "search-availability.page.tmpl", &models.TemplateData{})
}

// PostAvailability is the check availability page handler
func (m *Repository) PostAvailability(w http.ResponseWriter, r *http.Request) {
	start := r.Form.Get("start")
	end := r.Form.Get("end")

	w.Write([]byte(fmt.Sprintf("Start Date is %s and end date is %s", start, end)))
}

// Keep as close to code that uses this type
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

	render.RenderTemplate(w, r, "make-reservation.page.tmpl", &models.TemplateData{
		Form: forms.New(nil),
		Data: data,
	})

}

// Post Reservation handles the posting of a reservation form
func (m *Repository) PostReservation(w http.ResponseWriter, r *http.Request) {
	err := r.ParseForm()
	if err != nil {
		helpers.ServerError(w, err)
		return
	}

	reservation := models.Reservation{
		FirstName: r.Form.Get("first_name"),
		LastName:  r.Form.Get("last_name"),
		Phone:     r.Form.Get("phone"),
		Email:     r.Form.Get("email"),
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

		render.RenderTemplate(w, r, "make-reservation.page.tmpl", &models.TemplateData{
			Form: form,
			Data: data,
		})

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
	render.RenderTemplate(w, r, "contact.page.tmpl", &models.TemplateData{})
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

	render.RenderTemplate(w, r, "reservation-summary.page.tmpl", &models.TemplateData{
		Data: data,
	})
}
