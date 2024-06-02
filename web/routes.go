package web

import (
	"embed"
	"net/http"
)

//go:embed static
var staticFs embed.FS

func (h *Handler) Router() http.Handler {
	router := http.NewServeMux()
	router.HandleFunc("GET /", h.Index)
	router.HandleFunc("POST /", h.Paste)
	router.HandleFunc("GET /p/{id}/{key}", h.View)
	router.HandleFunc("POST /p/{id}/{key}", h.Decrypt)
	router.Handle("GET /static/", http.FileServer(http.FS(staticFs)))
	return router
}
