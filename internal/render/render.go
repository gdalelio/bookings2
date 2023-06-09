package render

import (
	"bytes"
	"fmt"
	"log"
	"net/http"
	"path/filepath"
	"text/template"

	"github.com/gdalelio/bookings/internal/config"
	"github.com/gdalelio/bookings/internal/models"
	"github.com/justinas/nosurf"
)

// Was in Instructors code but not used up to lessopn 74
var functions = template.FuncMap{}

var app *config.AppConfig
var pathToTemplates = "./templates"

// NewRenderer sets the config for the template package
func NewRenderer(a *config.AppConfig) {
	app = a
}

// AddDefaultData adds any data for all pages
func AddDefaultData(templateData *models.TemplateData, r *http.Request) *models.TemplateData {
	templateData.Flash = app.Session.PopString(r.Context(), "flash")
	templateData.Error = app.Session.PopString(r.Context(), "error")
	templateData.Warning = app.Session.PopString(r.Context(), "warning")
	templateData.CSRFToken = nosurf.Token(r)
	return templateData
}

//Template renders templates using html/template
func Template(w http.ResponseWriter, r *http.Request, tmpl string, templateData *models.TemplateData) {
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
	//need to remove hard coding to directory
	pages, err := filepath.Glob(fmt.Sprintf("%s/*.page.tmpl", pathToTemplates))
	if err != nil {
		return myCache, err
	}
	//range through all of the files ending *.page.tmpl
	for _, page := range pages {
		//want just the name and not the path
		name := filepath.Base(page)
		templateSet, err := template.New(name).Funcs(functions).ParseFiles(page)
		if err != nil {
			return myCache, err
		}

		matches, err := filepath.Glob(fmt.Sprintf("%s/*.layout.tmpl", pathToTemplates))
		if err != nil {
			return myCache, err
		}

		if len(matches) > 0 {
			templateSet, err = templateSet.ParseGlob(fmt.Sprintf("%s/*.layout.tmpl", pathToTemplates))
			if err != nil {
				return myCache, err
			}

		}
		//adds template to the map
		myCache[name] = templateSet

	}

	return myCache, err
}
