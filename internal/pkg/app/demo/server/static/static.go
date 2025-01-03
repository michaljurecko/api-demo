package static

import (
	"embed"
	"net/http"
)

//go:embed index.html
var embedFs embed.FS

func FS() http.Handler {
	return http.FileServer(http.FS(embedFs))
}
