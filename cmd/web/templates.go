package main

import (
	"html/template"
	"io/fs"
	"path/filepath"
	"time"

	"snippetbox.gregor-pifko/internal/models"
	"snippetbox.gregor-pifko/ui"
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

	pages, err := fs.Glob(ui.Files, "html/pages/*.tmpl")
	if err != nil {
		return nil, err
	}

	for _, page := range pages {
		name := filepath.Base(page)

		// Slice of filepath patterns for the templates we want to parse
		patterns := []string{
			"html/base.tmpl",
			"html/partials/*.tmpl",
			page,
		}
		// Custom defined functions must be added beforehand and then templates can be parsed.
		// Parse the template files from the ui.Files embedded filesystem.
		ts, err := template.New(name).Funcs(functions).ParseFS(ui.Files, patterns...)
		if err != nil {
			return nil, err
		}

		cache[name] = ts
	}

	return cache, nil
}
