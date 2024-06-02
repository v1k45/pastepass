package web

import "net/http"

func (h *Handler) Router() http.Handler {
	router := http.NewServeMux()
	router.HandleFunc("GET /", h.Index)
	router.HandleFunc("POST /", h.Paste)
	router.HandleFunc("GET /p/{id}/{key}", h.View)
	router.HandleFunc("POST /p/{id}/{key}", h.Decrypt)
	return router
}
