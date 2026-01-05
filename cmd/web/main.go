package main

import (
	"log"
	"net/http"
	"path/filepath"
)

func main() {
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

	mux.HandleFunc("GET /{$}", home)
	mux.HandleFunc("GET /snippet/view/{id}", snippetView)
	mux.HandleFunc("GET /snippet/create", getSnippetCreate)
	mux.HandleFunc("POST /snippet/create", postSnippetCreate)

	// Start a new server with http.ListenAndServe
	log.Println("Starting server on: 4000")
	err := http.ListenAndServe(":4000", mux)
	log.Fatal(err)
}

type neuteredFileSystem struct {
	fs http.FileSystem
}

func (nfs neuteredFileSystem) Open(path string) (http.File, error)  {
	f, err := nfs.fs.Open(path)	
	if err != nil {
		return nil, err
	}

	s, err := f.Stat()
	if err != nil {
		return nil, err
	}

	if s.IsDir() {
		index := filepath.Join(path, "index.html")
		if _, err := nfs.fs.Open(index); err != nil {
			closeErr := f.Close()
			if closeErr != nil {
				return nil, closeErr
			}

			return nil, err
		}
	}

	return f, nil
}




