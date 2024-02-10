package main

import (
	"errors"
	"fmt"
	"net/http"
	"strconv"

	"github.com/Zicchio/LG-Snippetbox/pkg/models"
)

func (app *application) home(w http.ResponseWriter, r *http.Request) {
	// NOTE: the 3 lines below are not required with pat, as pat library will check / exactly
	// if r.URL.Path != "/" {
	// 	app.notFound(w)
	// 	return
	// }

	s, err := app.snippets.Latest()
	if err != nil {
		app.serverError(w, err)
		return
	}
	// for _, snip := range s {
	// 	fmt.Fprintf(w, "%v\n", snip)
	// }

	tmplData := &templateData{Snippets: s}

	app.render(w, r, "home.page.go.tmpl", tmplData)
}

func (app *application) showSnippet(w http.ResponseWriter, r *http.Request) {
	// NOTE: the 4 lines below are required if "id" is query parameter instead of path variable
	// id, err := strconv.Atoi(r.URL.Query().Get("id"))
	// if err != nil || id < 1 {
	// 	app.notFound(w)
	// 	return
	// }

	id, err := strconv.Atoi(r.URL.Query().Get(":id"))
	if err != nil || id < 1 {
		app.notFound(w)
		return
	}

	s, err := app.snippets.Get(id)
	if err != nil {
		if errors.Is(err, models.ErrNoRecord) {
			app.notFound(w)
			return
		}
		app.serverError(w, err)
		return
	}
	tmplData := &templateData{Snippet: s}
	app.render(w, r, "show.page.go.tmpl", tmplData)
}

func (app *application) createSnippet(w http.ResponseWriter, r *http.Request) {
	// NOTE: the 5 lines below are commented as they were used to check that the method was a post, required as of routing from Go1.21 http package
	// if r.Method != http.MethodPost {
	// 	w.Header().Set("Allow", http.MethodPost)
	// 	app.clientError(w, http.StatusMethodNotAllowed)
	// 	return
	// }

	// Create some variables holding dummy data. We'll remove these later on
	// during the build.
	title := "O snail"
	content := "O snail\nClimb Mount Fuji,\nBut slowly, slowly!\n\n– Kobayashi Issa"
	expires := "7"

	id, err := app.snippets.Insert(title, content, expires)
	if err != nil {
		app.serverError(w, err)
		return
	}
	// redirect the user to the newly created snippet
	http.Redirect(w, r, fmt.Sprintf("/snippet?id=%d", id), http.StatusSeeOther)
}

func (app *application) createSnippetForm(w http.ResponseWriter, r *http.Request) {
	// TODO
	w.Write([]byte("This page will visualize a form creating a new snippet"))
}
