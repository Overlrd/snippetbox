package main

import (
	"github.com/justinas/alice"
	"net/http"
)

func (app *application) routes() http.Handler {
	// initialize new servermux
	mux := http.NewServeMux()

	// Create a file server which serves files out of the "./ui/static" directory.
	// Note that the path given to the http.Dir functin is relative to the project directory root.
	// We're using a custom FileSystem that checks if the requested path is a directory
	// If it is a directory we then try to Open() any index.html file in it. If no index.html
	// file exists, then this will return a os.ErrNotExist error (which in turn we return and
	// it will be transformed into a 404 Not Found response by http.Fileserver).
	fileserver := http.FileServer(neuteredFileSystem{http.Dir("./ui/static")})
	mux.Handle("./ui/static", http.NotFoundHandler())

	// Use the mux.Handler() function to register the file server as the handler for
	// all URL paths that start with "/static/". For matching paths, we strip the
	// "/static/" prefix before the request reaches the file server
	mux.Handle("GET /static/", http.StripPrefix("/static", fileserver))

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
