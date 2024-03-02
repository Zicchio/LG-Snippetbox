package main

import (
	"context"
	"errors"
	"fmt"
	"net/http"

	"github.com/Zicchio/LG-Snippetbox/pkg/models"
	"github.com/justinas/nosurf"
)

// Funciton secureHeaders set some header which provides resilience
// against cross site scripting attack. Since it is usefull against every
// request, it should be used before servemux
func secureHeaders(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("X-XSS-Protection", "1; mode=block")
		w.Header().Set("X-Frame-Options", "deny")

		next.ServeHTTP(w, r)
	})
}

func (app *application) logRequest(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		app.infoLog.Printf("%s - %s %s %s", r.RemoteAddr, r.Proto, r.Method, r.URL.RequestURI())
		next.ServeHTTP(w, r)
	})
}

func (app *application) recoverPanic(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// panic recovery is defined as a deferred function. Deferred functions are
		// always run in this case as Go unwinds the stack
		defer func() {
			if err := recover(); err != nil {
				// in case of error, close the connection and log the stack trace
				w.Header().Set("Connection", "close")
				app.serverError(w, fmt.Errorf("%s", err))
			}
		}()
		next.ServeHTTP(w, r)
	})
}

func noSurf(next http.Handler) http.Handler {
	csrfHandler := nosurf.New(next)
	csrfHandler.SetBaseCookie(http.Cookie{
		HttpOnly: true,
		Path:     "/",
		Secure:   true,
	})
	return csrfHandler
}

func (app *application) requireAuthentication(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if !app.isAuthenticated(r) {
			http.Redirect(w, r, "/user/login", http.StatusSeeOther)
			return
		}
		// Setting "Cache-Control: no-store" so that pages that require
		// authentication are not stored in the cache (either browser or
		// other caches)
		w.Header().Add("Cache-Control", "no-store")
		next.ServeHTTP(w, r)
	})
}

// authenticate will add an user to the context if there is a currently
// active user in the session, it exisst and it is active. Otherwise, move
// to the next element in the chain
func (app *application) authenticate(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Check if a authenticatedUserID value exists in the session. If this *isn't
		// present* then call the next handler in the chain as normal.
		exists := app.session.Exists(r, "authenticatedUserID") // NOTE: this was the previous definition of "app.isAuthneticated(r)""
		if !exists {                                           // if not authenticated, go on
			next.ServeHTTP(w, r)
			return
		}

		// If authenticated, check if it exists and is active.
		// If no, delete authentication and go on
		// If yes, add to context
		user, err := app.users.Get(app.session.GetInt(r, "authenticatedUserID"))
		if errors.Is(err, models.ErrNoRecord) || !user.Active {
			app.session.Remove(r, "authenticatedUserID")
			next.ServeHTTP(w, r)
			return
		} else if err != nil {
			app.serverError(w, err)
			return
		}
		ctx := context.WithValue(r.Context(), contextKeyIsAuthenticated, true)
		next.ServeHTTP(w, r.WithContext(ctx))

	})
}
