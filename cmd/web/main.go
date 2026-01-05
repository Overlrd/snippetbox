package main

import (
	"flag"
	"log/slog"
	"net/http"
	"os"
	"path/filepath"
)

type application struct {
	logger *slog.Logger
}

func main() {
	// Define command line flags
	addr := flag.String("addr", ":4000", "HTTP network address")
	flag.Parse()

	// Initialize a new logger
	logger := slog.New(slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{
		Level: slog.LevelDebug,
		AddSource: true,
	}))

	// Initialize a new instance of the application struct, containing the
	// dependencies
	app := &application{
		logger:logger,
	}

	logger.Info("starting server", "addr", *addr)

	// Start a new server with http.ListenAndServe
	err := http.ListenAndServe(*addr, app.routes())
	logger.Error(err.Error())
	os.Exit(1)
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
