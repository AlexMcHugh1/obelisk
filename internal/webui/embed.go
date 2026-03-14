package webui

import (
	"embed"
	"io/fs"
	"net/http"
	"path"
)

// @AlexMcHugh1: this is to embed the static files into the binary,
// so we don't have to worry about shipping the html files.
// The handler serves index.html for any unknown route to support client-side
// routing.

//go:embed static/*
var staticFiles embed.FS

func Handler() http.Handler {
	sub, err := fs.Sub(staticFiles, "static")
	if err != nil {
		panic(err)
	}

	fileServer := http.FileServer(http.FS(sub))

	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		cleanPath := path.Clean(r.URL.Path)
		if cleanPath == "/" {
			http.ServeFileFS(w, r, sub, "index.html")
			return
		}

		_, err := fs.Stat(sub, cleanPath[1:])
		if err != nil {
			http.ServeFileFS(w, r, sub, "index.html")
			return
		}

		fileServer.ServeHTTP(w, r)
	})
}
