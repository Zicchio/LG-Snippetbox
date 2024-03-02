package main

import (
	"net/http"

	"github.com/bmizerany/pat"
	"github.com/justinas/alice"
)

func (app *application) routes() http.Handler {
	middlewareChain := alice.New(app.recoverPanic, app.logRequest, secureHeaders)
	dynamicMiddleware := alice.New(app.session.Enable, noSurf, app.authenticate) // used for routes that requires session (handles via cookies)

	// Routing with base htpp package as of Go 1.21
	// mux := http.NewServeMux()
	// mux.HandleFunc("/", app.home)
	// mux.HandleFunc("/snippet", app.showSnippet)
	// mux.HandleFunc("/snippet/create", app.createSnippet)

	// Routing with pat.
	// NOTE: pattern are matched in the order they are registered
	mux := pat.New()
	mux.Get("/", dynamicMiddleware.ThenFunc(app.home)) // without alice: mux.Get("/", app.session.Enable(http.HandlerFunc(app.home)))

	mux.Get("/snippet/create", dynamicMiddleware.Append(app.requireAuthentication).ThenFunc(app.createSnippetForm))
	mux.Post("/snippet/create", dynamicMiddleware.Append(app.requireAuthentication).ThenFunc(app.createSnippet))

	mux.Get("/snippet/:id", dynamicMiddleware.ThenFunc(app.showSnippet)) // it is important that this is below "GET /snippet/create"

	mux.Get("/user/signup", dynamicMiddleware.ThenFunc(app.signupUserForm))
	mux.Post("/user/signup", dynamicMiddleware.ThenFunc(app.signupUser))
	mux.Get("/user/login", dynamicMiddleware.ThenFunc(app.loginUserForm))
	mux.Post("/user/login", dynamicMiddleware.ThenFunc(app.loginUser))
	mux.Post("/user/logout", dynamicMiddleware.Append(app.requireAuthentication).ThenFunc(app.logoutUser))

	fileServer := http.FileServer(http.Dir("./ui/static/"))
	mux.Get("/static/", http.StripPrefix("/static", fileServer))
	// return app.recoverPanic(app.logRequest(secureHeaders(mux))) // without alice package
	return middlewareChain.Then(mux)

}
