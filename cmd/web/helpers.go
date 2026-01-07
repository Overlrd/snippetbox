package main

import (
	"fmt"
	"net/http"
	"runtime/debug"
)

// The serverError helper writes a log entry at Error level (including the request
// method and URL as attributes), then sends a generic 500 Internal Server Error
// response to the user
func (app *application) serverError(w http.ResponseWriter, r *http.Request, err error) {
	var (
		method = r.Method
		url = r.URL.RequestURI()
		trace = string(debug.Stack())
	)

	app.logger.Error(err.Error(), "method", method, "url", url, "trace", trace)
	http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
}

// The clientError helper sends a specific status code and corresponding description
// to the user
func (app *application) clientError(w http.ResponseWriter, status int) {
	http.Error(w, http.StatusText(status), status)
}

func (app *application) render (w http.ResponseWriter, r *http.Request, status int, page string, data templateData) {
	// Retrive the appropriate template set from the cache besed on the page
	// name. If no entry exists in the cache with the provided name
	// then create a new error and call the serverError() helper method
	ts, ok := app.templateCache[page]
	if !ok {
		err := fmt.Errorf("the template %s does not exist", page)
		app.serverError(w, r, err)
		return
	}

	w.WriteHeader(status)

	// Execute the template set and write the response body
	err := ts.ExecuteTemplate(w, "base", data)
	if err != nil {
		app.serverError(w, r, err)
	}
}
