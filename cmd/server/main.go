package main

import (
	"log"
	"net/http"
	"os"
	"time"

	httpserver "tic-tac-go/internal/http"
)

// main is the entrypoint for the Tic-Tac-Go server application.
// It constructs the HTTP router and starts listening for incoming requests.
func main() {
	// Allow overriding the port via environment variable for convenience.
	port := os.Getenv("TICTACGO_PORT")
	if port == "" {
		port = "8080"
	}

	router := httpserver.NewRouter()

	server := &http.Server{
		Addr:              ":" + port,
		Handler:           router,
		ReadHeaderTimeout: 5 * time.Second,
	}

	log.Printf("Tic-Tac-Go server listening on :%s\n", port)

	if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
		log.Fatalf("server failed: %v", err)
	}
}
