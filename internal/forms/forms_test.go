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
	postedValues := url.Values{}
	//create form for the request
	form := New(postedValues)

	//creating a field on the form that does not exist and passing the request
	has := form.Has("whatever")

	//should return true - but the form has no fields
	if has {
		t.Error("form shows has field when it does not")
	}

	//creating new empty form  use for next test
	postedValues = url.Values{}

	//adding field with key = a and value of a
	postedValues.Add("a", "a")

	//reinitailizing the form passing the posted data
	form = New(postedValues)

	//Check form with the request to
	has = form.Has("a")

	//check to see if has is false
	if !has {
		t.Error("shows form does not have filed when it should")
	}

}

func TestForm_MinLength(t *testing.T) {
	//create a request as HAS needs a request when it is called
	postedValues := url.Values{}
	//create form for the request
	form := New(postedValues)

	form.MinLength("x", 10)

	if form.Valid() {
		t.Error("Form shows min length for non-existent field")
	}
	//adding field to form called some_field
	postedValues = url.Values{}

	//add some_field to the empty form with data
	postedValues.Add("some_field", "some value")

	//posting data for some_field into the form for testing
	form = New(postedValues)

	//calling the method MinLength with the field, min length to check and request
	form.MinLength("some_field", 100)

	//if form is valid comes back true - the form shows that test fails; min length is larger than data in field
	if form.Valid() {
		t.Error("Shows minlength of 100 met when data is shorter in length - fails check")
	}
	//re-initialize the form with empty values
	postedValues = url.Values{}

	//pass another field to the form
	postedValues.Add("another_field", "another value")

	//post values
	form = New(postedValues)

	//call MinLength to test
	form.MinLength("another_field", 5)

	if !form.Valid() {
		t.Error("Shows minlength of 5 was not met")
	}
}

func TestForm_IsEmail(t *testing.T) {

	//****  testing a form without an email field  ****//
	//create a request as IsMail needs a request when it is called
	postedValues := url.Values{}

	//create form for the request
	form := New(postedValues)

	// check to see that form is valid
	form.IsEmail("email")
	if form.Valid() {
		t.Error("form shows valid for non-existend field")
	}
	//****  testing a valid email address  ****//
	postedValues = url.Values{}
	postedValues.Add("email", "some@email.com")

	form = New(postedValues)
	form.IsEmail("email")
	if !form.Valid() {
		t.Error("got an invalid eamil when we should not have")
	}
	//****  testing an invalid email address  ****//
	postedValues = url.Values{}
	postedValues.Add("email", "x")

	form = New(postedValues)
	form.IsEmail("email")
	if form.Valid() {
		t.Error("got valid for invalid email address")
	}

}
func TestForm_IsPhone(t *testing.T) {
	//****  testing a form without an phone field  ****//
	//create a request as IsPhone needs a request when it is called
	postedValues := url.Values{}

	//create form for the request
	form := New(postedValues)

	// check to see that form is valid
	form.IsPhone("phone")
	if form.Valid() {
		t.Error("form shows valid for non-existent field")
	}

	//****  testing a valid phone number  ****
	postedValues = url.Values{}
	postedValues.Add("phone", "555-555555")

	form = New(postedValues)
	form.IsPhone("phone")
	if !form.Valid() {
		t.Error("got an invalid phone number type")
	}
	//****  testing an invalid phone number  ****
	postedValues = url.Values{}
	postedValues.Add("phone", "(703)-555-1212")

	form = New(postedValues)
	form.IsPhone("phone")
	if !form.Valid() {
		t.Error("got valid phone number type")
	}

}
