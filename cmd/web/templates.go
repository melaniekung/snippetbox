package main

import (
	"path/filepath"
	"text/template"
	"time"

	"github.com/melaniekung/snippetbox/internal/models"
)

// define struct type to hold dynamic data for templates
type templateData struct {
	CurrentYear int
	Snippet models.Snippet
	Snippets []models.Snippet
	Form any
}

func humanDate(t time.Time) string {
	return t.Format("02 Jan 2006 at 15:04")
}

var functions = template.FuncMap{
	"humanDate": humanDate,
}

func newTemplateCache() (map[string]*template.Template, error) {
	// initialize new map to act as cache
	cache := map[string]*template.Template{}

	// get slice of all filepaths matching pattern
	pages, err := filepath.Glob("./assets/html/pages/*.html")
	if err != nil {
		return nil, err
	}

	// loop through page filepaths
	for _, page := range pages {
		name := filepath.Base(page)

		// parse base template file into template set (ts)
		ts, err := template.New(name).Funcs(functions).ParseFiles("./assets/html/base.html")
		if err != nil {
			return nil, err
		}

		// add partials
		ts, err = ts.ParseGlob("./assets/html/partials/*.html")
		if err != nil {
			return nil, err
		}

		// add page template
		ts, err = ts.ParseFiles(page)
		if err != nil {
			return nil, err
		}

		// add template set to map
		cache[name] = ts
	}

	return cache, nil
}

