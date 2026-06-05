package server

import (
	"net/http"
	"time"

	"github.com/matipongpa/url-shortener/handler"
	"github.com/matipongpa/url-shortener/middleware"
)

func New(h *handler.Handler, limiter *middleware.Limiter) *http.Server {
	mux := http.NewServeMux()
	mux.HandleFunc("GET /health", h.Health)
	mux.HandleFunc("POST /shorten", h.Shorten)
	mux.HandleFunc("GET /{code}", h.Redirect)

	return &http.Server{
		Addr:              ":8080",
		Handler:           limiter.Middleware(mux),
		ReadHeaderTimeout: 2 * time.Second,
		ReadTimeout:       5 * time.Second,
		WriteTimeout:      10 * time.Second,
		IdleTimeout:       120 * time.Second,
	}
}
