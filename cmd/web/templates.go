package main

import (
	"html/template"
	"path/filepath"
	"time"

	"snippetbox.gregor-pifko/internal/models"
)

type templateData struct {
	CurrentYear         int
	Author              string
	Snippet             *models.Snippet
	Snippets            []*models.Snippet
	Form                any
	Flash               string
	IsAuthenticatedUser bool
	CSRFToken           string
}

func humanDate(t time.Time) string {
	return t.Format("02 Jan 2006 at 15:04")
}

// Custom defined template functions map
var functions = template.FuncMap{
	"humanDate": humanDate,
}

func newTemplateCache() (map[string]*template.Template, error) {
	cache := map[string]*template.Template{}

	pages, err := filepath.Glob("./ui/html/pages/*.tmpl")
	if err != nil {
		return nil, err
	}

	for _, page := range pages {
		name := filepath.Base(page)

		// Custom defined functions must be added beforehand and then base template can be parsed
		ts, err := template.New(name).Funcs(functions).ParseFiles("./ui/html/base.tmpl")
		if err != nil {
			return nil, err
		}

		// Add partials
		ts, err = ts.ParseGlob("./ui/html/partials/*.tmpl")
		if err != nil {
			return nil, err
		}

		// Add page
		ts, err = ts.ParseFiles(page)
		if err != nil {
			return nil, err
		}

		cache[name] = ts
	}

	return cache, nil
}
