package forms

import (
	"net/http"
	"net/http/httptest"
	"net/url"
	"testing"
)

func TestForm_Valid(t *testing.T) {
	request := httptest.NewRequest("POST", "/whaterver", nil)

	form := New(request.PostForm)

	isValid := form.Valid()

	if !isValid {
		t.Error("got invalid when should have been valid")
	}
}

func TestForm_Required(t *testing.T) {
	request := httptest.NewRequest("POST", "/whatever", nil)
	form := New(request.PostForm)

	form.Required("a", "b", "c")
	if form.Valid() {
		t.Error("Form shows valid when required fields are missing")
	}
	postedData := url.Values{}
	postedData.Add("a", "a")
	postedData.Add("b", "b")
	postedData.Add("c", "c")

	request, _ = http.NewRequest("POST", "/whatever", nil)

	request.PostForm = postedData
	form = New(request.PostForm)
	form.Required("a", "b", "c")

	if !form.Valid() {
		t.Error("Shows doe not have required fileds when it does")
	}

}

func TestForm_Has(t *testing.T) {

	//create a request as HAS needs a request when it is called
	request := httptest.NewRequest("POST", "/whatever", nil)
	//create form for the request
	form := New(request.PostForm)

	//creating a filed on the form that does not exist and passing the request
	has := form.Has("whatever", request)

	//should return true - but the form has no fields
	if has {
		t.Error("form shows has field when it does not")
	}

	//creating new empty form  use for next test
	postedData := url.Values{}

	//adding field with key = a and value of a
	postedData.Add("a", "a")

	//reinitailizing the form passing the posted data
	form = New(postedData)

	//Check form with the request to
	has = form.Has("a", request)

	//check to see if has is false
	if !has {
		t.Error("shows form does not have filed when it should")
	}

}
