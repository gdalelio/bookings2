package main

import (
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/alexedwards/scs/v2"
	"github.com/gdalelio/bookings/pkg/config"
	"github.com/gdalelio/bookings/pkg/handlers"
	"github.com/gdalelio/bookings/pkg/render"
)

// web based application for "Hello World!"
const portNumber = ":8080"

var app config.AppConfig
var session *scs.SessionManager

// main is the main application function
func main() {

	//change this to true when in production
	app.InProduction = false

	//setting up session parameters
	session = scs.New()
	session.Lifetime = 24 * time.Hour
	session.Cookie.Persist = true
	session.Cookie.SameSite = http.SameSiteLaxMode
	session.Cookie.Secure = app.InProduction //makes sure it is encrypted. Will need to change later as we have a localhost that doesnt have  https

	app.Session = session

	//create the template cache
	templateCache, err := render.CreateTemplateCache()

	if err != nil {
		log.Fatal("cannot create template cache")
	}

	app.TemplateCache = templateCache
	app.UseCache = false

	//set up handlers
	repo := handlers.NewRepo(&app)
	handlers.NewHandlers(repo)

	render.NewTemplates(&app) //refernce to the app config

	fmt.Printf(fmt.Sprintf("Starting application on port: %s", portNumber))
	//	_ = http.ListenAndServe(portNumber, nil) //starts up the webserver to listing on port 8080

	srv := &http.Server{
		Addr:    portNumber,
		Handler: routes(&app),
	}

	err = srv.ListenAndServe()
	log.Fatal(err)

}
