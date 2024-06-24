package web

import (
	"net/http"
	"net/http/httptest"
	"net/url"
	"testing"
	"time"

	"github.com/PuerkitoBio/goquery"
	"github.com/stretchr/testify/assert"
	"github.com/v1k45/pastepass/config"
	"github.com/v1k45/pastepass/db"
)

func TestHandlerIndex(t *testing.T) {
	h := NewHandler(nil)

	w := httptest.NewRecorder()
	r := httptest.NewRequest("GET", "/", nil)

	h.Index(w, r)
	assert.Equal(t, http.StatusOK, w.Code)

	doc, err := goquery.NewDocumentFromReader(w.Body)
	assert.NoError(t, err)

	title := doc.Find("title").Text()
	assert.Contains(t, title, config.AppName)

	// has a form
	form := doc.Find("form")
	assert.Equal(t, 1, form.Length())

	// has a textarea
	textarea := doc.Find("textarea")
	assert.Equal(t, 1, textarea.Length())

	// has a button
	button := doc.Find("button")
	assert.Equal(t, 1, button.Length())
}

func TestHandlerPaste(t *testing.T) {
	h := NewHandler(nil)

	t.Run("empty text", func(t *testing.T) {
		w := httptest.NewRecorder()
		r := httptest.NewRequest("POST", "/", nil)

		h.Paste(w, r)
		assert.Equal(t, http.StatusBadRequest, w.Code)

		doc, err := goquery.NewDocumentFromReader(w.Body)
		assert.NoError(t, err)

		errorMessage := doc.Find("hgroup > small").Text()
		assert.Equal(t, "Paste content is required.", errorMessage)

	})

	t.Run("invalid expiration", func(t *testing.T) {
		w := httptest.NewRecorder()
		r := httptest.NewRequest("POST", "/", nil)
		r.Form = url.Values{"text": {"hello"}, "expiration": {"invalid"}}

		h.Paste(w, r)
		assert.Equal(t, http.StatusBadRequest, w.Code)

		doc, err := goquery.NewDocumentFromReader(w.Body)
		assert.NoError(t, err)

		errorMessage := doc.Find("hgroup > small").Text()
		assert.Equal(t, "Invalid expiration time.", errorMessage)
	})

	db, err := db.NewTestDB()
	assert.NoError(t, err)
	defer db.Reset()

	t.Run("success", func(t *testing.T) {
		h = NewHandler(db)

		w := httptest.NewRecorder()
		r := httptest.NewRequest("POST", "/", nil)
		r.Form = url.Values{"text": {"hello"}, "expiration": {"1h"}}

		h.Paste(w, r)
		assert.Equal(t, http.StatusOK, w.Code)

		doc, err := goquery.NewDocumentFromReader(w.Body)
		assert.NoError(t, err)

		pasteUrl := doc.Find("pre").Text()
		assert.Contains(t, pasteUrl, "/p/")
	})

	t.Run("url scheme", func(t *testing.T) {
		h = NewHandler(db)

		w := httptest.NewRecorder()
		r := httptest.NewRequest("POST", "/", nil)
		r.Header.Set("X-Forwarded-Proto", "https")
		r.Form = url.Values{"text": {"hello"}, "expiration": {"1h"}}

		h.Paste(w, r)
		assert.Equal(t, http.StatusOK, w.Code)

		doc, err := goquery.NewDocumentFromReader(w.Body)
		assert.NoError(t, err)

		pasteUrl := doc.Find("pre").Text()
		assert.Contains(t, pasteUrl, "https")
	})

}

func TestHandlerView(t *testing.T) {
	db, err := db.NewTestDB()
	assert.NoError(t, err)
	defer db.Reset()

	h := NewHandler(db)

	paste, err := db.NewPaste("test paste", time.Now().Add(time.Hour))
	assert.NoError(t, err)

	w := httptest.NewRecorder()
	r := httptest.NewRequest("GET", "/p/"+paste.ID+"/"+paste.Key, nil)
	r.SetPathValue("id", paste.ID)
	r.SetPathValue("key", paste.Key)

	h.View(w, r)
	assert.Equal(t, http.StatusOK, w.Code)

	doc, err := goquery.NewDocumentFromReader(w.Body)
	assert.NoError(t, err)

	// document does not contain the content of the paste
	// it is only displayed upon decryption
	assert.NotContains(t, doc.Text(), "test paste")

	t.Run("expired paste", func(t *testing.T) {
		paste, err := db.NewPaste("test paste", time.Now().Add(-time.Hour))
		assert.NoError(t, err)

		w := httptest.NewRecorder()
		r := httptest.NewRequest("GET", "/p/"+paste.ID+"/"+paste.Key, nil)
		r.SetPathValue("id", paste.ID)
		r.SetPathValue("key", paste.Key)

		h.View(w, r)
		assert.Equal(t, http.StatusNotFound, w.Code)
	})

	t.Run("nonexistent paste", func(t *testing.T) {
		w := httptest.NewRecorder()
		r := httptest.NewRequest("GET", "/p/nonexistent/key", nil)
		r.SetPathValue("id", "nonexistent")
		r.SetPathValue("key", "key")

		h.View(w, r)
		assert.Equal(t, http.StatusNotFound, w.Code)
	})
}

func TestHandlerDecrypt(t *testing.T) {
	db, err := db.NewTestDB()
	assert.NoError(t, err)
	defer db.Reset()

	h := NewHandler(db)

	paste, err := db.NewPaste("test paste", time.Now().Add(time.Hour))
	assert.NoError(t, err)

	w := httptest.NewRecorder()
	r := httptest.NewRequest("POST", "/p/"+paste.ID+"/"+paste.Key, nil)
	r.SetPathValue("id", paste.ID)
	r.SetPathValue("key", paste.Key)

	h.Decrypt(w, r)
	assert.Equal(t, http.StatusOK, w.Code)

	doc, err := goquery.NewDocumentFromReader(w.Body)
	assert.NoError(t, err)

	// document contains the content of the paste
	pre := doc.Find("pre").Text()
	assert.Contains(t, pre, "test paste")

	// paste is deleted after decryption
	_, err = db.Get(paste.ID)
	assert.Error(t, err)

	// try to decrypt again
	w = httptest.NewRecorder()
	h.Decrypt(w, r)
	assert.Equal(t, http.StatusInternalServerError, w.Code)
}

func TestHandlerStatic(t *testing.T) {
	h := NewHandler(nil)

	t.Run("css", func(t *testing.T) {
		w := httptest.NewRecorder()
		r := httptest.NewRequest("GET", "/static/pico.min.css", nil)

		h.Router().ServeHTTP(w, r)
		assert.Equal(t, http.StatusOK, w.Code)
	})

	t.Run("js", func(t *testing.T) {
		w := httptest.NewRecorder()
		r := httptest.NewRequest("GET", "/static/pastepass.js", nil)

		h.Router().ServeHTTP(w, r)
		assert.Equal(t, http.StatusOK, w.Code)
	})

	t.Run("404", func(t *testing.T) {
		w := httptest.NewRecorder()
		r := httptest.NewRequest("GET", "/static/notfound", nil)

		h.Router().ServeHTTP(w, r)
		assert.Equal(t, http.StatusNotFound, w.Code)
	})
}
