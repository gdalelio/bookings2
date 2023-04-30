package render

import (
	"net/http"
	"testing"

	"github.com/gdalelio/bookings/internal/models"
)

// TestAddDefaultData is testing the AddDefaultData
func TestAddDefaultData(t *testing.T) {

	var templateData models.TemplateData

	request, err := getSession()
	if err != nil {
		t.Error(err)
	}

	//need to get session data for context
	session.Put(request.Context(), "Flash", "123")

	result := AddDefaultData(&templateData, request)

	if result.Flash == "123" {
		t.Error("flash vaule of 123")
	}

}

// getSession returns a pointer to an http.Request object
func getSession() (*http.Request, error) {
	requestObj, err := http.NewRequest("GET", "/some-url", nil)

	if err != nil {
		return nil, err
	}

	//getting context from the request object for session data
	ctx := requestObj.Context()

	//put session data into the requestObject using the contet retrieved and looking for X-Session key
	ctx, _ = session.Load(ctx, requestObj.Header.Get("X-Session"))

	requestObj = requestObj.WithContext(ctx)

	return requestObj, nil
}

func TestNewTemplates(t *testing.T) {
	NewTemplates(app)
}

func TestCreateTemplateCache(t *testing.T) {
	pathToTemplates = "./../../templates"
	_, err := CreateTemplateCache()
	if err != nil {
		t.Error(err)
	}
}
