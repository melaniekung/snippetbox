package main

import (
	"io/fs"
	"path/filepath"
	"text/template"
	"time"

	"github.com/melaniekung/snippetbox/assets"
	"github.com/melaniekung/snippetbox/internal/models"
)

// define struct type to hold dynamic data for templates
type templateData struct {
	CurrentYear     int
	Snippet         models.Snippet
	Snippets        []models.Snippet
	Form            any
	Flash           string
	IsAuthenticated bool
	CSRFToken       string
}

func humanDate(t time.Time) string {
	if t.IsZero() {
		return ""
	}

	return t.UTC().Format("02 Jan 2006 at 15:04")
}

var functions = template.FuncMap{
	"humanDate": humanDate,
}

func newTemplateCache() (map[string]*template.Template, error) {
	// initialize new map to act as cache
	cache := map[string]*template.Template{}

	// get slice of all filepaths matching pattern
	pages, err := fs.Glob(assets.Files, "html/pages/*html")
	if err != nil {
		return nil, err
	}

	// loop through page filepaths
	for _, page := range pages {
		name := filepath.Base(page)

		// set filepath patterns for templates to parse
		patterns := []string{
			"html/base.html",
			"html/partials/*.html",
			page,
		}

		// parse template files from embedded filesystem
		ts, err := template.New(name).Funcs(functions).ParseFS(assets.Files, patterns...)
		if err != nil {
			return nil, err
		}

		// add template set to map
		cache[name] = ts
	}

	return cache, nil
}
