package main

import (
	"html/template"
	"path/filepath"
	"time"

	"github.com/Zicchio/LG-Snippetbox/pkg/forms"
	"github.com/Zicchio/LG-Snippetbox/pkg/models"
)

// Type templateData is used to wrap the content of all template data that
// are required in one or more templates.
// The author of the book chooses to merge all templates into one struct,
// which might or might not be a scalable approach.
// Personal NOTE: I don't love the fact that the same template wrapper
// is used in multiple places, but I guess it is valid anyway
type templateData struct {
	CSRFToken       string
	CurrentYear     int
	Flash           string
	Form            *forms.Form
	IsAuthenticated bool
	Snippet         *models.Snippet
	Snippets        []*models.Snippet
}

// Custom format the date. The custom format is similar to RFC822
func humanDate(t time.Time) string {
	if t.IsZero() {
		return ""
	}
	return t.UTC().Format("02 Jan 2006 at 15:04") // Memo: format must contain a "rearrangment" of the value Jan 2 15:04:05 2006 MST
}

// Memo: custom template functions must have at most 1 return value, or at most (value, err)
// Reference: https://pkg.go.dev/text/template#FuncMap
// NOTE: the function map must be "registered", which is done in the function below
var functions = template.FuncMap{
	"humanDate": humanDate,
}

func newTemplateCache(dir string) (map[string]*template.Template, error) {
	cache := make(map[string]*template.Template)
	// get all files matching the template suffix corresponding to full pages in the input directory
	pages, err := filepath.Glob(filepath.Join(dir, "*.page.gohtml"))
	if err != nil {
		return nil, err
	}

	// for each page, add the layout and the partial
	for _, page := range pages {
		fname := filepath.Base(page)
		// fmt.Printf("caching %s (file %s)\n", page, fname)
		ts, err := template.New(fname).Funcs(functions).ParseFiles(page)
		if err != nil {
			return nil, err
		}
		ts, err = ts.ParseGlob(filepath.Join(dir, "*.layout.gohtml"))
		if err != nil {
			return nil, err
		}
		ts, err = ts.ParseGlob(filepath.Join(dir, "*.partial.gohtml"))
		if err != nil {
			return nil, err
		}
		cache[fname] = ts
	}
	return cache, nil
}
