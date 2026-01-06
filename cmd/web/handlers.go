package main

import (
	"errors"
	"fmt"
	// "html/template"
	"net/http"
	"strconv"

	"github.com/Overlrd/snippetbox/internal/models"
)

// home handler function
func (app *application)home(w http.ResponseWriter, r *http.Request) {
	// Restrict the home handler to the "/" url pattern
	// also consider adding the restriction when registering the handler
	// mux.HandleFunc("/{$}", home)
	// https://gopherbuilders.com/articles/avoiding-catchall-root-route-golang-servemux

	// if r.URL.Path != "/" {
	// 	http.NotFound(w, r)
	// 	return
	// }
	w.Header().Add("Server", "Go")

	snippets, err := app.snippets.Latest()
	if err != nil {
		app.serverError(w, r, err)
		return
	}

	for _, snippet := range snippets {
		fmt.Fprintf(w, "%+v\n", snippet)
	}
	
	// // Initialize a slice containing the paths to the template files 
	// // The file containing our base template must be the *first* file in the slice
	// files := []string{
	// 	"./ui/html/base.tmpl",
	// 	"./ui/html/partials/nav.tmpl",
	// 	"./ui/html/pages/home.tmpl",
	// }
	//
	// // Use template.ParseFiles() function to read the template files into a
	// // template set and handle error
	// ts, err := template.ParseFiles(files...)
	// if err != nil {
	// 	app.serverError(w, r, err) // Use the serverError helper
	// 	return
	// }
	//
	// // The we use the Execute() method on the template set to write the
	// // template content as the response body. The last parameter to Execute()
	// // represents any dynamic data that we want to pass in
	// // err = ts.Execute(w, nil)
	//
	// // Use the ExecuteTemplate() method to write the content of the "base"
	// // template as the response body
	// err = ts.ExecuteTemplate(w, "base", nil)
	// if err != nil {
	// 	app.serverError(w, r, err)
	// }
}

// snippetView: Display a specific snippet
func (app *application) snippetView(w http.ResponseWriter, r *http.Request) {
	// Extract the value of the "id" wildcard from the request
	// and sanitize it
	id, err := strconv.Atoi(r.PathValue("id"))
	if err != nil || id < 1 {
		http.NotFound(w, r)
		return
	}

	// Use the SnippetModel.Get() method to retrieve the data for a
	// specific record based on its ID. If no matching record is found,
	// return a 404 Not Found response
	snippet, err := app.snippets.Get(id)
	if err != nil {
		if errors.Is(err, models.ErrNoRecord) {
			http.NotFound(w, r)
		} else {
			app.serverError(w, r, err)
		}
		return
	}
	// Write the snippet data as plain-text HTTP response
	fmt.Fprintf(w, "%+v", snippet)

}

// getSnippetCreate: Display a form for creating a new snippet
func (app application) getSnippetCreate(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("Creates a snippet"))
}

// postSnippetCreate: Save a new snippet
func (app *application) postSnippetCreate(w http.ResponseWriter, r *http.Request) {
	// Use w.WriteHeader() method to send a 201 status code.
	// Any changes made to the header map after calling w.WriteHeader()
	// or w.Write() will have no effect on the headers that the user receives.
	// w.WriteHeader(http.StatusCreated)
	//
	// w.Write([]byte("Save a new snippet..."))
	title := "O snall"
	content := "O snail\nClimb Mount Fuji,\nBut slowly, slowly!\n\n– Kobayashi Issa" 
	expires := 7

	// Pass the data to the SnippetModel.Insert() method
	id, err := app.snippets.Insert(title, content, expires)
	if err != nil {
		app.serverError(w, r, err)
		return
	}

	// Redirect the user to the relevant page for the snippet
	http.Redirect(w, r, fmt.Sprintf("/snippet/view/%d", id), http.StatusSeeOther)
} 




