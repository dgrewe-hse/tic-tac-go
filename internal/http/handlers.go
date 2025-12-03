package http

import (
	"encoding/json"
	"net/http"
)

// healthResponse is the JSON structure returned by the /health endpoint.
type healthResponse struct {
	Status string `json:"status"`
}

// healthHandler serves a minimal health check response so that clients
// and deployment environments can verify the server is running.
func healthHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	resp := healthResponse{Status: "ok"}

	if err := json.NewEncoder(w).Encode(resp); err != nil {
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return
	}
}


