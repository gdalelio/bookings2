package main

import (
	"encoding/gob"
	"fmt"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/alexedwards/scs/v2"
	"github.com/gdalelio/bookings/internal/config"
	"github.com/gdalelio/bookings/internal/handlers"
	"github.com/gdalelio/bookings/internal/helpers"
	"github.com/gdalelio/bookings/internal/models"
	"github.com/gdalelio/bookings/internal/render"
)

// web based application for "Hello World!"
const portNumber = ":8080"

var app config.AppConfig
var session *scs.SessionManager
var infoLog *log.Logger
var errorLog *log.Logger

// main is the main application function
func main() {
	err := run()
	if err != run() {
		log.Fatal(err)
	}

	fmt.Printf("Starting application on port: %s", portNumber)
	//	_ = http.ListenAndServe(portNumber, nil) //starts up the webserver to listing on port 8080

	srv := &http.Server{
		Addr:    portNumber,
		Handler: routes(&app),
	}

	err = srv.ListenAndServe()
	if err != nil {
		log.Fatal(err)
	}

}

// Run does the majority of the work
func run() error {
	//what am I going to put in the session
	gob.Register(models.Reservation{})
	//change this to true when in production
	app.InProduction = false
	

	//set up information log
	infoLog = log.New(os.Stdout, "INFO\t", log.Ldate|log.Ltime)
	app.InfoLog = infoLog
	//setting up error log
	errorLog = log.New(os.Stdout, "Error\t", log.Ldate|log.Ltime|log.Lshortfile)
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
	templateCache, err := render.CreateTemplateCache()

	if err != nil {
		log.Fatal("cannot create template cache")
		return err
	}

	app.TemplateCache = templateCache
	app.UseCache = false

	//set up handlers
	repo := handlers.NewRepo(&app)
	handlers.NewHandlers(repo)

	render.NewTemplates(&app) //refernce to the app config
	helpers.NewHelpers(&app)  //refernce to the app config

	return nil
}
