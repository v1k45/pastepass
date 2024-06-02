package web

import (
	"context"
	"fmt"
	"log/slog"
	"net/http"

	"github.com/v1k45/pastepass/db"
	"github.com/v1k45/pastepass/views"
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
		slog.Error("validation_error", "error", "paste content is required")
		errorResponse(w, http.StatusBadRequest, "Invalid Data", "Paste content is required.")
		return
	}

	expiresAt, err := getExpiresAt(r.FormValue("expiration"))
	if err != nil {
		slog.Error("validation_error", "error", err)
		errorResponse(w, http.StatusBadRequest, "Invalid Data", "Invalid expiration time.")
		return
	}

	paste, err := h.DB.NewPaste(pastedText, expiresAt)
	if err != nil {
		slog.Error("cannot_create_paste", "error", err)
		errorResponse(w, http.StatusInternalServerError, "Internal Server Error", "Failed to create paste, please try again later.")
		return
	}

	var scheme string
	if r.Header.Get("X-Forwarded-Proto") != "" {
		scheme = r.Header.Get("X-Forwarded-Proto")
	} else {
		scheme = "http"
	}

	url := fmt.Sprintf("%s://%s/p/%s/%s", scheme, r.Host, paste.ID, paste.Key)

	component := views.PasteSuccess(url)
	component.Render(context.Background(), w)
}

func (h *Handler) View(w http.ResponseWriter, r *http.Request) {
	id := r.PathValue("id")
	if _, err := h.DB.Get(id); err != nil {
		slog.Error("cannot_view_paste", "error", err, "id", id)
		errorResponse(w, http.StatusNotFound, "Not Found", "The paste you are looking for is either expired or does not exist.")
		return
	}

	component := views.View()
	component.Render(context.Background(), w)
}

func (h *Handler) Decrypt(w http.ResponseWriter, r *http.Request) {
	id, key := r.PathValue("id"), r.PathValue("key")
	decryptedText, err := h.DB.Decrypt(id, key)
	if err != nil {
		slog.Error("cannot_decrypt_paste", "error", err, "id", id)
		errorResponse(
			w, http.StatusInternalServerError,
			"Internal Server Error", "The paste you are looking for is either expired, corrputed or does not exist.")
		return
	}

	component := views.Decrypt(decryptedText)
	component.Render(context.Background(), w)
}
