package handler

import (
	"encoding/json"
	"errors"
	"log"
	"net/http"

	"github.com/matipongpa/url-shortener/shortener"
	"github.com/matipongpa/url-shortener/storage"
)

type Payload struct {
	URL string `json:"url"`
}

type Handler struct {
	short   *shortener.Shortener
	baseURL string
}

func New(short *shortener.Shortener, base string) *Handler {
	return &Handler{
		short,
		base,
	}
}

func (h *Handler) Shorten(w http.ResponseWriter, r *http.Request) {
	r.Body = http.MaxBytesReader(w, r.Body, 1<<20)

	var payload Payload

	dec := json.NewDecoder(r.Body)
	dec.DisallowUnknownFields()
	if err := dec.Decode(&payload); err != nil {
		log.Printf("shorten decode: %v", err)
		ErrorResponse(w, http.StatusBadRequest, "invalid request body")
		return
	}

	if err := ValidateURL(payload.URL); err != nil {
		ErrorResponse(w, http.StatusBadRequest, "Url is invalid")
		return
	}

	code, err := h.short.Shorten(r.Context(), payload.URL)
	if err != nil {
		ErrorResponse(w, http.StatusInternalServerError, "Shorten Url Error")
		return
	}

	Response(w, http.StatusCreated, map[string]string{
		"short_url": h.baseURL + "/" + code,
		"code":      code,
	})
}

func (h *Handler) Redirect(w http.ResponseWriter, r *http.Request) {
	code := r.PathValue("code")
	if code == "" {
		ErrorResponse(w, http.StatusBadRequest, "Path is not valid")
		return
	}
	url, err := h.short.Resolve(r.Context(), code)
	if errors.Is(err, storage.ErrNotFound) {
		ErrorResponse(w, http.StatusNotFound, "code not found")
		return
	}
	if err != nil {
		ErrorResponse(w, http.StatusInternalServerError, "Redirect Error")
		return
	}
	http.Redirect(w, r, url, http.StatusFound)
}
