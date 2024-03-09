package main

import (
	"errors"
	"fmt"
	"net/http"
	"strconv"

	"github.com/Zicchio/LG-Snippetbox/pkg/forms"
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

	app.render(w, r, "home.page.gohtml", tmplData)
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
	app.render(w, r, "show.page.gohtml", tmplData)
}

func (app *application) createSnippet(w http.ResponseWriter, r *http.Request) {
	// NOTE: the 5 lines below are commented as they were used to check that the method was a post, required as of routing from Go1.21 http package
	// if r.Method != http.MethodPost {
	// 	w.Header().Set("Allow", http.MethodPost)
	// 	app.clientError(w, http.StatusMethodNotAllowed)
	// 	return
	// }

	err := r.ParseForm()
	if err != nil {
		app.clientError(w, http.StatusBadRequest)
		return
	}

	form := forms.New(r.PostForm)
	form.Required("title", "content", "expires")
	form.MaxLenght("title", 100)
	form.AdmittedValues("expires", "365", "7", "1")

	if !form.Valid() {
		app.render(w, r, "create.page.gohtml", &templateData{
			Form: form,
		})
		return
	}

	id, err := app.snippets.Insert(form.Get("title"), form.Get("content"), form.Get("expires"))
	if err != nil {
		app.serverError(w, err)
		return
	}

	app.session.Put(r, "flash", "Snippet successfully created!")
	http.Redirect(w, r, fmt.Sprintf("/snippet/%d", id), http.StatusSeeOther)
}

func (app *application) createSnippetForm(w http.ResponseWriter, r *http.Request) {
	app.render(w, r, "create.page.gohtml", &templateData{
		Form: forms.New(nil),
	})
}

func (app *application) signupUserForm(w http.ResponseWriter, r *http.Request) {
	app.render(w, r, "signup.page.gohtml", &templateData{
		Form: forms.New(nil),
	})
}

func (app *application) signupUser(w http.ResponseWriter, r *http.Request) {
	err := r.ParseForm()
	if err != nil {
		app.clientError(w, http.StatusBadRequest)
		return
	}

	form := forms.New(r.PostForm)
	form.Required("name")
	form.MaxLenght("name", 255)
	form.Required("email")
	form.MaxLenght("email", 255)
	form.MatchesPattern("email", forms.EmailRX)
	form.Required("password")
	form.MinLength("password", 10)
	form.MaxLenght("password", 255)

	if !form.Valid() {
		app.render(w, r, "signup.page.gohtml", &templateData{
			Form: form,
		})
		return
	}

	err = app.users.Insert(form.Get("name"), form.Get("email"), form.Get("password"))
	if err != nil {
		if errors.Is(err, models.ErrDuplicateEmail) {
			form.Errors.Add("email", "Email address already in use")
			app.render(w, r, "signup.page.gohtml", &templateData{Form: form})
		} else {
			app.serverError(w, err)
		}
		return
	}
	app.session.Put(r, "flash", "Signup was successfull. Please log in.")
	http.Redirect(w, r, "/user/login", http.StatusSeeOther)
}

func (app *application) loginUserForm(w http.ResponseWriter, r *http.Request) {
	app.render(w, r, "login.page.gohtml", &templateData{Form: forms.New(nil)})
}

func (app *application) loginUser(w http.ResponseWriter, r *http.Request) {
	err := r.ParseForm()
	if err != nil {
		app.clientError(w, http.StatusBadRequest)
		return
	}

	form := forms.New(r.PostForm)
	id, err := app.users.Authenticate(form.Get("email"), form.Get("password"))
	if err != nil {
		if errors.Is(err, models.ErrInvalidCredentials) {
			form.Errors.Add("generic", "Email or password was incorrect")
			app.render(w, r, "login.page.gohtml", &templateData{Form: form})
		} else {
			app.serverError(w, err)
		}
		return
	}
	app.session.Put(r, "authenticatedUserID", id)
	app.session.Put(r, "flash", "Login was succesfull.")
	http.Redirect(w, r, "/snippet/create", http.StatusSeeOther)
}

func (app *application) logoutUser(w http.ResponseWriter, r *http.Request) {
	// NOTE: this code should not be available to unauthenticated users
	app.session.Remove(r, "authenticatedUserID")
	app.session.Put(r, "flash", "You have been succesfully logged out")
	http.Redirect(w, r, "/", http.StatusSeeOther)
}
