// Copyright 2025 Esslingen University of Applied Sciences
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.
//
// Author: Dennis Grewe
// Version: 1.0.0
// Date: 2025-12-03

package http

import (
	"encoding/json"
	"errors"
	"net/http"

	"tic-tac-go/internal/models"
	"tic-tac-go/internal/service"
)

// healthResponse is the JSON structure returned by the /health endpoint.
type healthResponse struct {
	Status string `json:"status"`
}

// request / response DTOs
type createPlayerRequest struct {
	Name string `json:"name"`
}

type createPlayerResponse struct {
	PlayerID string `json:"playerId"`
	Name     string `json:"name"`
}

type createGameRequest struct {
	Mode string `json:"mode"`
}

type createGameResponse struct {
	GameID      string     `json:"gameId"`
	Mode        string     `json:"mode"`
	Board       [][]string `json:"board"`
	CurrentTurn string     `json:"currentTurn"`
	Status      string     `json:"status"`
	Winner      string     `json:"winner"`
}

// CreatePlayerHandler returns an http.HandlerFunc bound to a PlayerService.
func CreatePlayerHandler(playerSvc service.PlayerService) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req createPlayerRequest
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			http.Error(w, http.StatusText(http.StatusBadRequest), http.StatusBadRequest)
			return
		}
		if req.Name == "" {
			http.Error(w, "name is required", http.StatusBadRequest)
			return
		}

		player, err := playerSvc.CreatePlayer(r.Context(), req.Name)
		if err != nil {
			http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusCreated)

		resp := createPlayerResponse{
			PlayerID: player.ID,
			Name:     player.Name,
		}
		_ = json.NewEncoder(w).Encode(resp)
	}
}

func CreateGameHandler(gameSvc service.GameService) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		playerID := PlayerIDFromContext(r.Context())
		if playerID == "" {
			http.Error(w, "missing X-Player-Id header", http.StatusBadRequest)
			return
		}

		var req createGameRequest
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			http.Error(w, http.StatusText(http.StatusBadRequest), http.StatusBadRequest)
			return
		}

		mode := models.GameMode(req.Mode)
		gameState, err := gameSvc.CreateGame(r.Context(), playerID, mode)
		if err != nil {
			if errors.Is(err, service.ErrInvalidGameMode) {
				http.Error(w, "invalid mode", http.StatusBadRequest)
				return
			}
			http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
			return
		}

		// Convert board [3][3]Symbol to [][]string for JSON response.
		board := make([][]string, 3)
		for i := 0; i < 3; i++ {
			board[i] = make([]string, 3)
			for j := 0; j < 3; j++ {
				board[i][j] = string(gameState.Board[i][j])
			}
		}

		resp := createGameResponse{
			GameID:      gameState.ID,
			Mode:        string(gameState.Mode),
			Board:       board,
			CurrentTurn: string(gameState.CurrentTurn),
			Status:      string(gameState.Status),
			Winner:      gameState.Winner,
		}

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusCreated)
		_ = json.NewEncoder(w).Encode(resp)
	}
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
