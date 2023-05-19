package handlers

import (
	"encoding/gob"
	"fmt"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"time"

	"text/template"

	"github.com/alexedwards/scs/v2"
	"github.com/gdalelio/bookings/internal/config"
	"github.com/gdalelio/bookings/internal/models"
	"github.com/gdalelio/bookings/internal/render"
	"github.com/go-chi/chi"
	"github.com/go-chi/chi/middleware"
	"github.com/justinas/nosurf"
)

var app config.AppConfig
var session *scs.SessionManager
var pathToTemplates = "./../../templates"
var functions = template.FuncMap{}

func getRoutes() http.Handler {
	//creating var to allow for setting up testing objects

	//what am I going to put in the session
	gob.Register(models.Reservation{})
	//change this to true when in production
	app.InProduction = false

	//set up information log
	infoLog := log.New(os.Stdout, "INFO\t", log.Ldate|log.Ltime)
	app.InfoLog = infoLog
	//setting up error log
	errorLog := log.New(os.Stdout, "Error\t", log.Ldate|log.Ltime|log.Lshortfile)
	app.ErrorLog = errorLog

	//setting up session parameters
	session = scs.New()
	session.Lifetime = 24 * time.Hour
	session.Cookie.Persist = true
	session.Cookie.SameSite = http.SameSiteLaxMode

	//makes sure it is encrypted. Will need to change later as we have a localhost that doesnt have  https
	session.Cookie.Secure = app.InProduction

	app.Session = session

	//create the template cache
	templateCache, err := CreateTestTemplateCache()

	if err != nil {
		log.Fatal("cannot create template cache - insider setup_test for handler")
	}

	app.TemplateCache = templateCache
	app.UseCache = true

	//set up handlers
	//TODO need to add databae repo to the function call
	repo := NewRepo(&app)
	NewHandlers(repo)

	//refernce to the app config
	render.NewRenderer(&app)

	mux := chi.NewRouter()

	mux.Use(middleware.Recoverer)
	//mux.Use(NoSurf)  removed this to allow for testing and not checking for the cross site token
	mux.Use(SessionLoad)

	mux.Get("/", Repo.Home)
	mux.Get("/about", Repo.About)
	mux.Get("/generals-quarters", Repo.Generals)
	mux.Get("/majors-suite", Repo.Majors)

	mux.Get("/search-availability", Repo.Availability)
	mux.Post("/search-availability", Repo.PostAvailability)
	mux.Post("/search-availability-json", Repo.AvailabilityJSON)

	mux.Get("/contact", Repo.Contact)

	//mux.Get("/reservation", Repo.Reservation)
	mux.Get("/make-reservation", Repo.Reservation)
	mux.Post("/make-reservation", Repo.PostReservation)
	mux.Get("/reservation-summary", Repo.ReservationSummary)

	fileServer := http.FileServer(http.Dir("./static/"))
	mux.Handle("/static/*", http.StripPrefix("/static", fileServer))

	return mux
}

// NoSurf prevents Cross-Site Request Forgery attacks
func NoSurf(next http.Handler) http.Handler {
	csrfHandler := nosurf.New(next)
	csrfHandler.SetBaseCookie(http.Cookie{
		HttpOnly: true,
		Path:     "/",
		Secure:   app.InProduction,
		SameSite: http.SameSiteLaxMode,
	})
	return csrfHandler
}

// SessionLoad loads the session information and saves the session on every request
func SessionLoad(next http.Handler) http.Handler {
	return session.LoadAndSave(next)
}
func CreateTestTemplateCache() (map[string]*template.Template, error) {
	// myCache := make(map[string]*template.Template)
	// creates an empty map of pointers to the templates that are in the cache
	myCache := map[string]*template.Template{}

	//need to start with the *.page.tmpl first and then the page you want to retrieve
	//get all of the files named *.tmpl of the ./templates folder
	//need to remove hard coding to directory
	pages, err := filepath.Glob(fmt.Sprintf("%s/*.page.tmpl", pathToTemplates))
	if err != nil {
		return myCache, err
	}
	//range through all of the files ending *.page.tmpl
	for _, page := range pages {
		//want just the name and not the path
		name := filepath.Base(page)
		templateSet, err := template.New(name).Funcs(functions).ParseFiles(page)
		if err != nil {
			return myCache, err
		}

		matches, err := filepath.Glob(fmt.Sprintf("%s/*.layout.tmpl", pathToTemplates))
		if err != nil {
			return myCache, err
		}

		if len(matches) > 0 {
			templateSet, err = templateSet.ParseGlob(fmt.Sprintf("%s/*.layout.tmpl", pathToTemplates))
			if err != nil {
				return myCache, err
			}

		}
		//adds template to the map
		myCache[name] = templateSet

	}

	return myCache, err
}
