package main

import (
	"net/http"

	"github.com/justinas/alice"
)

func (app *application) routes() http.Handler {
	middlewareChain := alice.New(app.recoverPanic, app.logRequest, secureHeaders)
	mux := http.NewServeMux()
	mux.HandleFunc("/", app.home)
	mux.HandleFunc("/snippet", app.showSnippet)
	mux.HandleFunc("/snippet/create", app.createSnippet)

	fileServer := http.FileServer(http.Dir("./ui/static/"))
	mux.Handle("/static/", http.StripPrefix("/static", fileServer))
	// return app.recoverPanic(app.logRequest(secureHeaders(mux))) // without alice package
	return middlewareChain.Then(mux)

}
