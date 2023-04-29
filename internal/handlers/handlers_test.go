package handlers

import (
	"net/http"
	"net/http/httptest"
	"net/url"
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
	{"get search-availability", "/search-availability", "GET", []postData{}, http.StatusOK},
	{"contact", "/contact", "GET", []postData{}, http.StatusOK},
	{"make reservation get", "/make-reservation", "GET", []postData{}, http.StatusOK},
	{"search-availability-json post", "/search-availability-json", "POST", []postData{
		{key: "start", value: "2020-01-01"},
		{key: "end", value: "2020-01-02"},
	}, http.StatusOK},
	{"make reservation post", "/make-reservation", "POSt", []postData{
		{key: "first_name", value: "John"},
		{key: "last_name", value: "Smith"},
		{key: "email", value: "noone@nowhere.com"},
		{key: "phone", value: "555-555-5555"},
	}, http.StatusOK},
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
			values := url.Values{}
			for _, x := range e.parms {
				values.Add(x.key, x.value)
			}

			response, err := testServer.Client().PostForm(testServer.URL+e.url, values)
			if err != nil {
				t.Log(err)
				t.Fatal(err)
			}
			if response.StatusCode != e.expectedStatusCode {
				t.Errorf("for %s expected %d but got %d", e.name, e.expectedStatusCode, response.StatusCode)
			}
		}
	}
}
