package main

import (
	"log"
	"net/http"
	"os"

	"github.com/matipongpa/url-shortener/handler"
	"github.com/matipongpa/url-shortener/middleware"
	"github.com/matipongpa/url-shortener/server"
	"github.com/matipongpa/url-shortener/shortener"
	"github.com/matipongpa/url-shortener/storage"
)

func main() {
	baseURL := os.Getenv("BASE_URL")
	if baseURL == "" {
		baseURL = "http://localhost:8080"
	}
	store := storage.NewMemory()
	short := shortener.New(store)
	h := handler.New(short, baseURL)
	limiter := middleware.NewRateLimiter(10, 2)
	srv := server.New(h, limiter)

	log.Println("server starting on http://localhost:8080")
	if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
		log.Fatalf("server error: %v", err)
	}
}
