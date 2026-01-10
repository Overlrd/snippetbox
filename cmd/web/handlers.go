package main

import (
	"errors"
	"fmt"
	"net/http"
	"strconv"

	"github.com/Overlrd/snippetbox/internal/models"
	"github.com/Overlrd/snippetbox/internal/validator"
	"github.com/go-playground/form/v4"
)

// Update our snippetCreateForm struct to include tags which tell the
// decocer how to map HTML form values into the different struct fields.
// The struct tag 'form:"-"' telles the decoder to completely ignore a
// field during decoding
type snippetCreateForm struct {
	Title               string `form:"title"`
	Content             string `form:"content"`
	Expires             int    `form:"expires"`
	validator.Validator `form:"-"`
}

// home handler function
func (app *application) home(w http.ResponseWriter, r *http.Request) {
	// Restrict the home handler to the "/" url pattern
	// also consider adding the restriction when registering the handler
	// mux.HandleFunc("/{$}", home)
	// https://gopherbuilders.com/articles/avoiding-catchall-root-route-golang-servemux

	// if r.URL.Path != "/" {
	// 	http.NotFound(w, r)
	// 	return
	// }

	snippets, err := app.snippets.Latest()
	if err != nil {
		app.serverError(w, r, err)
		return
	}

	// Call the newTemplateData() helper to get a templateData struct
	// containing the 'default' data and add the snippets slice to it
	data := app.newTemplateData(r)
	data.Snippets = snippets

	// Use the new render helper
	app.render(w, r, http.StatusOK, "home.tmpl", data)
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

	data := app.newTemplateData(r)
	data.Snippet = snippet

	// Use the new render helper
	app.render(w, r, http.StatusOK, "view.tmpl", data)
}

// getSnippetCreate: Display a form for creating a new snippet
func (app application) getSnippetCreate(w http.ResponseWriter, r *http.Request) {
	data := app.newTemplateData(r)

	// Initialize a new createSnippetForm instance and pass it to the template
	data.Form = snippetCreateForm{
		Expires: 365,
	}

	app.render(w, r, http.StatusOK, "create.tmpl", data)
}

// postSnippetCreate: Save a new snippet
func (app *application) postSnippetCreate(w http.ResponseWriter, r *http.Request) {

	// Declare a new empty instance of the snippetCreateForm struct.
	var form snippetCreateForm

	err := app.decodePostForm(r, &form)
	if err != nil {
		app.clientError(w, http.StatusBadRequest)
		return
	}

	// Call the Decode() method on the form decoder, passing in the current
	// request and *a pointerù to our snippetCreateForm struct. This will
	// fill our struct with the relevant values from the HTML form.
	// If there's a problem return a 400 Bad Request response to the client.
	err = app.formDecoder.Decode(&form, r.PostForm)
	if err != nil {
		app.clientError(w, http.StatusBadRequest)
	}

	// Use the embedded Validator struct's CheckField() method to execute our
	// validation checks
	form.CheckField(validator.NotBlank(form.Title), "title", "This field cannot be blank")
	form.CheckField(validator.MaxChars(form.Title, 100), "title", "This field cannot be more than 100 characters long")
	form.CheckField(validator.NotBlank(form.Content), "content", "This field cannot be blank")
	form.CheckField(validator.PermittedValue(form.Expires, 1, 7, 365), "expires", "This field must equal 1, 7 or 365")

	// Use the valid() method to see if any checks failed.

	// If there are any validation errors, then re-display the create.tmpl template,
	// passing in the postSnippetCreate instance as dynamic data in the form
	// field. Note that we use the HTTP status code 422 Unprocessable Entity
	// when sending the response to indicate that there was a validation error.
	if !form.Valid() {
		data := app.newTemplateData(r)
		data.Form = form
		app.render(w, r, http.StatusUnprocessableEntity, "create.tmpl", data)
		return
	}

	id, err := app.snippets.Insert(form.Title, form.Content, form.Expires)
	if err != nil {
		app.serverError(w, r, err)
		return
	}

	http.Redirect(w, r, fmt.Sprintf("/snippet/view/%d", id), http.StatusSeeOther)
}
