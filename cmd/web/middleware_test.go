package main

import (
	"fmt"
	"net/http"
	"testing"
)

func TestNoSurf(t *testing.T) {
	var myH myHandler //defined in setup_test.go
	//requires pointer to the handler - hence &
	h := NoSurf(&myH)

	//check for the type the handler (h)  returns
	switch v := h.(type) {
	case http.Handler:
		//do nothing test passed
	default:
		t.Errorf(fmt.Sprintf("type is not http.Handler, but is %T", v))
	}
}

func TestSessionLoad(t *testing.T) {
	var myH myHandler //defined in setup_test.go
	//requires pointer to the handler - hence &
	h := SessionLoad(&myH)

	//check for the type the handler (h)  returns
	switch v := h.(type) {
	case http.Handler:
		//do nothing
	default:
		t.Errorf(fmt.Sprintf("type is not http.Handler, but is %T", v))
	}
}
