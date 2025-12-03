package http

import (
	"net/http"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
)

// NewRouter constructs the root HTTP router for the Tic-Tac-Go server.
// For now it only exposes a simple health endpoint; additional routes
// for game and player APIs will be added later.
func NewRouter() http.Handler {
	r := chi.NewRouter()

	// Basic middlewares for logging, recovery and timeouts.
	r.Use(middleware.RequestID)
	r.Use(middleware.RealIP)
	r.Use(middleware.Logger)
	r.Use(middleware.Recoverer)
	r.Use(middleware.Timeout(60 * time.Second))

	// Health check endpoint to verify the server is up.
	r.Get("/health", healthHandler)

	return r
}
