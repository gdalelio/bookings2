package handlers

import (
	"net/http"
	"net/http/httptest"
	"testing"
)

type postData struct {
	key   string
	value string
}

var theTests = []struct {
	//mame of individual test
	name string
	url  string
	//get or post
	method string
	//values for testing
	parms []postData
	//http codes expected 400, 404, 500, etc.
	expectedStatusCode int
}{
	{"home", "/", "GET", []postData{}, http.StatusOK},
	{"about", "/about", "GET", []postData{}, http.StatusOK},
	{"generals quarters", "/generals-quarters", "GET", []postData{}, http.StatusOK},
	{"majors suite", "/majors-suite", "GET", []postData{}, http.StatusOK},
	{"sa", "/search-availability", "GET", []postData{}, http.StatusOK},
	{"contact", "/contact", "GET", []postData{}, http.StatusOK},
	{"make reservation", "/make-reservation", "GET", []postData{}, http.StatusOK},
}

func TestHandlers(t *testing.T) {
	routes := getRoutes()
	//set up test server to send the routes to
	testServer := httptest.NewTLSServer(routes)

	//close the test server when done
	defer testServer.Close()

	for _, e := range theTests {
		//testing two types - gets and posts

		if e.method == "GET" {
			//appending the Test Server URL to ours
			response, err := testServer.Client().Get(testServer.URL + e.url)
			//test for error
			if err != nil {
				//log the error
				t.Log(err)
				//fail the test due to error
				t.Fatal(err)
			}

			if response.StatusCode != e.expectedStatusCode {
				t.Errorf("for %s, expected %d but got %d", e.name, e.expectedStatusCode, response.StatusCode)
			}
		} else {

		}
	}
}
