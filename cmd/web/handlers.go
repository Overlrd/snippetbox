package main

import (
	"fmt"
	"net/http"
	"strconv"
)

// home handler function
func home(w http.ResponseWriter, r *http.Request) {
	// Restrict the home handler to the "/" url pattern
	// also consider adding the restriction when registering the handler
	// mux.HandleFunc("/{$}", home)
	// https://gopherbuilders.com/articles/avoiding-catchall-root-route-golang-servemux

	// if r.URL.Path != "/" {
	// 	http.NotFound(w, r)
	// 	return
	// }

	w.Write([]byte("hello, world"))
}

// snippetView: Display a specific snippet
func snippetView(w http.ResponseWriter, r *http.Request) {
	// Extract the value of the "id" wildcard from the request
	// and sanitize it
	id, err := strconv.Atoi(r.PathValue("id"))
	if err != nil || id < 1 {
		http.NotFound(w, r)
		return
	}

	fmt.Fprintf(w, "Display a specific snippet with ID %d...", id)
}

// getSnippetCreate: Display a form for creating a new snippet
func getSnippetCreate(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("Creates a snippet"))
}

// postSnippetCreate: Save a new snippet
func postSnippetCreate(w http.ResponseWriter, r *http.Request) {
	// Use w.WriteHeader() method to send a 201 status code.
	// Any changes made to the header map after calling w.WriteHeader()
	// or w.Write() will have no effect on the headers that the user receives.
	w.WriteHeader(http.StatusCreated)

	w.Write([]byte("Save a new snippet..."))
}
