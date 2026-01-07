package main

import "net/http"

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

	mux.HandleFunc("GET /{$}", app.home)
	mux.HandleFunc("GET /snippet/view/{id}", app.snippetView)
	mux.HandleFunc("GET /snippet/create", app.getSnippetCreate)
	mux.HandleFunc("POST /snippet/create", app.postSnippetCreate)

	return commonHeader(mux)
}
