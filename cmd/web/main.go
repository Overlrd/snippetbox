package main

import (
	"log"
	"net/http"
)

func main() {
	// initialize new servermux
	mux := http.NewServeMux()
	mux.HandleFunc("GET /{$}", home)
	mux.HandleFunc("GET /snippet/view/{id}", snippetView)
	mux.HandleFunc("GET /snippet/create", getSnippetCreate)
	mux.HandleFunc("POST /snippet/create", postSnippetCreate)

	// Start a new server with http.ListenAndServe
	log.Println("Starting server on: 4000")
	err := http.ListenAndServe(":4000", mux)
	log.Fatal(err)
}
