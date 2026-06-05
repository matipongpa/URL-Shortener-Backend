package handler

import (
	"encoding/json"
	"errors"
	"log"
	"net/http"
	"net/url"
)

func Response(w http.ResponseWriter, status int, data any) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	if err := json.NewEncoder(w).Encode(data); err != nil {
		log.Printf("response encode: %v", err)
	}

}

func ErrorResponse(w http.ResponseWriter, status int, msg string) {
	Response(w, status, map[string]string{"error": msg})
}

func ValidateURL(rawURL string) error {
	if len(rawURL) > 2048 {
		return errors.New("Url length over limit 2048")
	}
	u, err := url.ParseRequestURI(rawURL)
	if err != nil {
		return errors.New("Url is not valid")
	}
	if u.Scheme != "http" && u.Scheme != "https" {
		return errors.New("scheme must be http or https")
	}
	if u.Host == "" {
		return errors.New("url must have a host")
	}
	return nil
}
