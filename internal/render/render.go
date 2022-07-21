package render

import (
	"bytes"
	"errors"
	"fmt"
	"html/template"
	"net/http"
	"path/filepath"

	"github.com/justinas/nosurf"
	"github.com/wagnojunior/booking/internal/config"
	"github.com/wagnojunior/booking/internal/models"
)

// Set the path to the templates
var pathToTemplates = "./templates"

// Map of functions that can be used in a template, usually functions that are not built into the language
var functions = template.FuncMap{}

// Local variable of typo <*AppConfig>
var app *config.AppConfig

// NewTemplates sets the config for the template package
func NewTemplates(a *config.AppConfig) {
	app = a
}

// AddDefaultData sets the data that will be available to all templates, that is the default data.
func AddDefaultData(td *models.TemplateData, r *http.Request) *models.TemplateData {
	td.Flash = app.Session.PopString(r.Context(), "flash")
	td.Error = app.Session.PopString(r.Context(), "error")
	td.Warning = app.Session.PopString(r.Context(), "warning")
	td.CSRFToken = nosurf.Token(r)
	return td
}

// RenderTemplate renders a template
func RenderTemplate(w http.ResponseWriter, r *http.Request, tmpl string, td *models.TemplateData) error {
	var tc map[string]*template.Template

	if app.UseCache {
		// get the template cache from the app config
		tc = app.TemplateCache
	} else {
		tc, _ = CreateTemplateCache()
	}

	t, ok := tc[tmpl]
	if !ok {
		// log.Fatal("Could not get template from template cache")
		return errors.New("Can't get template from cache")
	}

	buf := new(bytes.Buffer)

	// Sets the default data <td>
	td = AddDefaultData(td, r)

	// Data is passed to the template when it is executed; in this case <td> is passed to the template
	_ = t.Execute(buf, td)

	_, err := buf.WriteTo(w)
	if err != nil {
		fmt.Println("error writing template to browser", err)
		return err
	}

	return nil

}

// CreateTemplateCache creates a template cache as a map
func CreateTemplateCache() (map[string]*template.Template, error) {

	// <myCache> maps a string to a pointer to <template.Template>
	myCache := map[string]*template.Template{}

	// Gets the file path of all files in the folder <templates> that end with <.page.tmpl>
	pages, err := filepath.Glob(fmt.Sprintf("%s/*.page.tmpl", pathToTemplates))
	if err != nil {
		return myCache, err
	}

	// Loop through all the pages
	for _, page := range pages {
		// Gets the file path base
		name := filepath.Base(page)

		// <New> allocates a new HTML template with the given <name>
		// <Funcs> adds the elements of the argument map to the template's function map
		// <ParseFiles> parses the named files and associates the resulting templates with t
		ts, err := template.New(name).Funcs(functions).ParseFiles(page)
		if err != nil {
			return myCache, err
		}

		// <Glob> returns the names of all files matching pattern or nil if there is no matching file
		matches, err := filepath.Glob(fmt.Sprintf("%s/*.layout.tmpl", pathToTemplates))
		if err != nil {
			return myCache, err
		}

		// If the length of <matches> is grater than zero, there are might be layouts associated with templates
		// In fact, this is the case for <about.page.tmpl> and <home.page.tmpl> which reference <base.layout.tmpl>
		if len(matches) > 0 {
			// ParseGlob parses the template definitions in the files identified by the pattern and associates the
			//resulting templates with t.
			ts, err = ts.ParseGlob(fmt.Sprintf("%s/*.layout.tmpl", pathToTemplates))
			if err != nil {
				return myCache, err
			}
		}

		// Add the parsed template to the template map <myCache>
		myCache[name] = ts
	}

	return myCache, nil
}
