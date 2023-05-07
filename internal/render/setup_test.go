package render

import (
	"encoding/gob"
	"log"
	"net/http"
	"os"
	"testing"
	"time"

	"github.com/alexedwards/scs/v2"
	"github.com/gdalelio/bookings/internal/config"
	"github.com/gdalelio/bookings/internal/models"
)

// need for a session object for testing
var session *scs.SessionManager

// need to be able to create a copy of the app object in render.go - and will
// assign a pointer to it to alloow for testing
var testApp config.AppConfig

// TestMain is a tesing object for testing main
func TestMain(m *testing.M) {

	//what am I going to put in the session
	gob.Register(models.Reservation{})

	//change this to true when in production
	//using testApp as we don't have the app object as we are testing this
	testApp.InProduction = false

	//set up information log
	infoLog := log.New(os.Stdout, "INFO\t", log.Ldate|log.Ltime)
	testApp.InfoLog = infoLog
	//setting up error log
	errorLog := log.New(os.Stdout, "Error\t", log.Ldate|log.Ltime|log.Lshortfile)
	testApp.ErrorLog = errorLog

	//setting up session parameters
	session = scs.New()
	session.Lifetime = 24 * time.Hour
	session.Cookie.Persist = true
	session.Cookie.SameSite = http.SameSiteLaxMode

	//makes sure it is encrypted. Will need to change later as we have a localhost that doesnt have  https
	session.Cookie.Secure = testApp.InProduction

	testApp.Session = session

	app = &testApp

	//gets run before our test then exits
	os.Exit((m.Run()))
}
