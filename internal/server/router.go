package server

import (
	"net/http"

	"github.com/go-chi/chi/v5"
)

// NewRouter builds the application's HTTP router.
func NewRouter() http.Handler {
	r := chi.NewRouter()

	r.Get("/health", healthHandler)

	return r
}
