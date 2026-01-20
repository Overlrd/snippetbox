package main

import (
	"github.com/Overlrd/snippetbox/ui"
	"github.com/justinas/alice"
	"net/http"
)

func (app *application) routes() http.Handler {
	// initialize new servermux
	mux := http.NewServeMux()

	// Use http.FileServerFS() function to create a HTTP handler which
	// serves the embedded files in ui.Files. The static file are located
	// in the "static" folder of th ui.Files embedded filesystem, so there
	// is no more need to strip the prefix from the request URL
	mux.Handle("GET /static/", http.FileServerFS(ui.Files))

	// Unprotected application routes using the "dynamic" middleware chain.
	dynamic := alice.New(app.sessionManager.LoadAndSave, noSurf, app.authenticate)

	// Update these routes to use the nes dynamic middleware chain
	mux.Handle("GET /{$}", dynamic.ThenFunc(app.home))
	mux.Handle("GET /snippet/view/{id}", dynamic.ThenFunc(app.snippetView))
	mux.Handle("GET /user/signup", dynamic.ThenFunc(app.getUserSignup))
	mux.Handle("POST /user/signup", dynamic.ThenFunc(app.postUserSignup))
	mux.Handle("GET /user/login", dynamic.ThenFunc(app.getUserLogin))
	mux.Handle("POST /user/login", dynamic.ThenFunc(app.postUserLogin))

	// Protected (authenticated-only) application routes, using a new "protected"
	// middleware chain which includes the requireAuthentication middleware.
	protected := dynamic.Append(app.requireAuthentication)

	mux.Handle("GET /snippet/create", protected.ThenFunc(app.getSnippetCreate))
	mux.Handle("POST /snippet/create", protected.ThenFunc(app.postSnippetCreate))
	mux.Handle("POST /user/logout", protected.ThenFunc(app.postUserLogout))

	// Create a middleware chain containing our 'standard' middleware
	// which will be used for every request our application receives.
	standard := alice.New(app.recoverPanic, app.logRequest, commonHeader)

	return standard.Then(mux)
}
