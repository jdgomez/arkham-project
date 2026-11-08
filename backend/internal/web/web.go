// Package web serves the front end that `make build` embeds in the binary.
package web

import (
	"embed"
	"io/fs"
	"net/http"
)

//go:embed all:dist
var dist embed.FS

// Handler serves the built front end.
func Handler() http.Handler {
	site, err := fs.Sub(dist, "dist")
	if err != nil {
		panic(err)
	}
	return http.FileServerFS(site)
}
