package forms

import (
	"fmt"
	"log"
	"net/url"
	"regexp"
	"strings"

	"github.com/asaskevich/govalidator"
)

// Form creates a custom form struct, embeds a url.Values object
type Form struct {
	url.Values
	Errors errors
}

// Valid returns true if there are no errors, otherwise false
func (f *Form) Valid() bool {
	return len(f.Errors) == 0
}

// New initializes a form struct
func New(data url.Values) *Form {
	return &Form{
		data,
		errors(map[string][]string{}),
	}

}

// Required checks for required fields
func (f *Form) Required(fields ...string) {
	for _, field := range fields {
		value := f.Get(field)
		//check if the value field is empty and if it is add the error message for the field
		if strings.TrimSpace(value) == "" {
			f.Errors.Add(field, "This field cannot be blank")
		}
	}
}

// Has checks if the form field is in post and not empty
func (f *Form) Has(field string) bool {
	x := f.Get(field)
	if x == "" {
		return false
	}
	return true
}

// MinLength checks for minimum length for field
func (f *Form) MinLength(field string, length int) bool {
	x := f.Get(field)

	if len(x) < length {
		f.Errors.Add(field, fmt.Sprintf("This field must be at least %d charachters long", length))
		return false
	}
	return true
}

// IsEmail checks for valid eamil address
func (f *Form) IsEmail(field string) {
	if !govalidator.IsEmail(f.Get(field)) {
		f.Errors.Add(field, "Invalid email address")
	}
}

func (f *Form) IsPhone(field string) {

	//set up the regex pattern to be matched against
	re := regexp.MustCompile(`^(?:(?:\(?(?:00|\+)([1-4]\d\d|[1-9]\d?)\)?)?[\-\.\ \\\/]?)?((?:\(?\d{1,}\)?[\-\.\ \\\/]?){0,})(?:[\-\.\ \\\/]?(?:#|ext\.?|extension|x)[\-\.\ \\\/]?(\d+))?$`)
	str := f.Get(field)
	strLen := len(str)
	minLength := 10

	//check to see if the phone number is the minimum length
	if strLen < minLength {
		log.Println("****************************************************************")
		log.Printf("number doesn't meet min length for phone number %s. Was %d and should be %d", str, strLen, minLength)
		log.Println("****************************************************************")
	} else {
		//check to see if the string matches the phone format
		valid := re.MatchString(str)

		if !valid {
			log.Printf("Phone number is invalid format: %s", str)
		} else {
			log.Printf("Phone number is valid format: %s", str)
		}
	}

}
