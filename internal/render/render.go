package render

import (
	"bytes"
	"log"
	"net/http"
	"path/filepath"
	"text/template"

	"github.com/gdalelio/bookings/internal/config"
	"github.com/gdalelio/bookings/internal/models"
	"github.com/justinas/nosurf"
)

var app *config.AppConfig

// NewTemplates sets the config for the template package
func NewTemplates(a *config.AppConfig) {
	app = a
}

// AddDefaultData adds any data for all pages
func AddDefaultData(templateData *models.TemplateData, r *http.Request) *models.TemplateData {
	templateData.CSRFToken = nosurf.Token(r)
	return templateData
}

// RenderTemplate renders templates using html/template
func RenderTemplate(w http.ResponseWriter, r *http.Request, tmpl string, templateData *models.TemplateData) {
	var templateCache map[string]*template.Template

	if app.UseCache {
		//get the templage cache from the app config
		templateCache = app.TemplateCache
	} else {
		templateCache, _ = CreateTemplateCache()
	}

	theTemplate, ok := templateCache[tmpl] //uses string passed in
	if !ok {
		log.Fatal("could not get template from template cahce")
	}
	//creating a buffer to hold template and see if I can execute from the buffer - for finer grain error checking
	buf := new(bytes.Buffer)

	templateData = AddDefaultData(templateData, r)

	_ = theTemplate.Execute(buf, templateData) //lets me know what the issue is if I cannot execute the template

	//render the template
	_, err := buf.WriteTo(w)
	if err != nil {
		log.Println(err)
	}
}

func CreateTemplateCache() (map[string]*template.Template, error) {
	// myCache := make(map[string]*template.Template)
	// creates an empty map of pointers to the templates that are in the cache
	myCache := map[string]*template.Template{}

	//need to start with the *.page.tmpl first and then the page you want to retrieve
	//get all of the files named *.tmpl of the ./templates folder
	pages, err := filepath.Glob("./templates/*.page.tmpl")
	if err != nil {
		return myCache, err
	}
	//range through all of the files ending *.page.tmpl
	for _, page := range pages {
		//want just the name and not the path
		name := filepath.Base(page)
		templateSet, err := template.New(name).ParseFiles(page)
		if err != nil {
			return myCache, err
		}

		matches, err := filepath.Glob("./templates/*.layout.tmpl")
		if err != nil {
			return myCache, err
		}

		if len(matches) > 0 {
			templateSet, err = templateSet.ParseGlob("./templates/*.layout.tmpl")
			if err != nil {
				return myCache, err
			}

		}
		//adds template to the map
		myCache[name] = templateSet

	}

	return myCache, err
}
