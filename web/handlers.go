package web

import (
	"context"
	"fmt"
	"net/http"

	"github.com/v1k45/paste/db"
	"github.com/v1k45/paste/views"
)

type Handler struct {
	DB *db.DB
}

func NewHandler(db *db.DB) *Handler {
	return &Handler{DB: db}
}

func (h *Handler) Index(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusOK)
	component := views.Index()
	component.Render(context.Background(), w)
}

func (h *Handler) Paste(w http.ResponseWriter, r *http.Request) {
	pastedText := r.FormValue("text")
	if pastedText == "" {
		errorResponse(w, http.StatusBadRequest, "Invalid Data", "Paste content is required.")
		return
	}

	expiresAt, err := getExpiresAt(r.FormValue("expiration"))
	if err != nil {
		errorResponse(w, http.StatusBadRequest, "Invalid Data", "Invalid expiration time.")
		return
	}

	paste, err := h.DB.NewPaste(pastedText, expiresAt)
	if err != nil {
		errorResponse(w, http.StatusInternalServerError, "Internal Server Error", "Failed to create paste, please try again later.")
		return
	}

	var scheme string
	if r.TLS == nil {
		scheme = "http"
	} else {
		scheme = "https"
	}
	url := fmt.Sprintf("%s://%s/p/%s/%s", scheme, r.Host, paste.ID, paste.Key)

	component := views.PasteSuccess(url)
	component.Render(context.Background(), w)
}

func (h *Handler) View(w http.ResponseWriter, r *http.Request) {
	if _, err := h.DB.Get(r.PathValue("id")); err != nil {
		errorResponse(w, http.StatusNotFound, "Not Found", "The paste you are looking for is either expired or does not exist.")
		return
	}

	component := views.View()
	component.Render(context.Background(), w)
}

func (h *Handler) Decrypt(w http.ResponseWriter, r *http.Request) {
	decryptedText, err := h.DB.Decrypt(r.PathValue("id"), r.PathValue("key"))
	if err != nil {
		errorResponse(
			w, http.StatusInternalServerError,
			"Internal Server Error", "The paste you are looking for is either expired, corrputed or does not exist.")
		return
	}

	component := views.Decrypt(decryptedText)
	component.Render(context.Background(), w)
}
