package web

import (
	"context"
	"errors"
	"net/http"
	"time"

	"github.com/v1k45/paste/views"
)

var (
	expirationTimes = map[string]time.Duration{
		"1h": time.Hour,
		"1d": 24 * time.Hour,
		"1w": 7 * 24 * time.Hour,
		"2w": 2 * 7 * 24 * time.Hour,
		"4w": 4 * 7 * 24 * time.Hour,
	}
)

func getExpiresAt(expiresAt string) (time.Time, error) {
	expiresDuration, found := expirationTimes[expiresAt]
	if !found {
		return time.Time{}, errors.New("invalid expiration time")
	}

	return time.Now().Add(expiresDuration), nil
}

func errorResponse(w http.ResponseWriter, status int, title, message string) {
	w.WriteHeader(status)
	component := views.Error(title, message)
	component.Render(context.Background(), w)
}
