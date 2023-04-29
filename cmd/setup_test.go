package main

import (
	"net/http"
	"os"
	"testing"
)

// allows for setting up testing parameters, variables, etc. before the testing
func TestMain(m *testing.M) {

	os.Exit(m.Run())
}

type myHandler struct{}

func (mh *myHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
}
