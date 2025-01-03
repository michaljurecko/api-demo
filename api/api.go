package api

import (
	"embed"
	"io/fs"
	"net/http"
)

//go:embed gen/openapi.yaml
var embedFs embed.FS

func GenFS() http.Handler {
	dir, err := fs.Sub(embedFs, "gen")
	if err != nil {
		panic(err)
	}
	return http.FileServer(http.FS(dir))
}
