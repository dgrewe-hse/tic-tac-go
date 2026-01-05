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
	"log"
	"net/http"
	"strconv"
	"time"

	"tic-tac-go/internal/models"
	"tic-tac-go/internal/service"
	"tic-tac-go/internal/store"
	"tic-tac-go/internal/ws"

	"github.com/go-chi/chi/v5"
	"github.com/gorilla/websocket"
)

// ------------------------------------------
// DATA TRANSFER OBJECTS / REQUEST / RESPONSE

// healthResponse is the JSON structure returned by the /health endpoint.
type healthResponse struct {
	Status string `json:"status"`
}

// CREATE request / response DTOs
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

// GAME SUMMARY DTO
type gameSummaryDTO struct {
	GameID    string `json:"gameId"`
	Mode      string `json:"mode"`
	Status    string `json:"status"`
	CreatedAt string `json:"createdAt"`
	CreatedBy struct {
		PlayerID string `json:"playerId"`
		Name     string `json:"name"`
	} `json:"createdBy"`
}

type listGamesResponse struct {
	Games []gameSummaryDTO `json:"games"`
}

type makeMoveRequest struct {
	Row int `json:"row"`
	Col int `json:"col"`
}

// ----------------------

// ----------------------
// HANDLER IMPLEMENTATIONS

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

// ----------------------

// -----------------------------
// GAME HANDLERS

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

// ListGamesHandler provides a handle to receive the list of crated games
// to list it within a dashboard and to provide means to join specific games
func ListGamesHandler(gameSvc service.GameService) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		q := r.URL.Query()
		var filter store.GameFilter

		if modeStr := q.Get("mode"); modeStr != "" {
			m := models.GameMode(modeStr)
			filter.Mode = &m
		}
		if statusStr := q.Get("status"); statusStr != "" {
			s := models.GameStatus(statusStr)
			filter.Status = &s
		}
		if limitStr := q.Get("limit"); limitStr != "" {
			if v, err := strconv.Atoi(limitStr); err == nil {
				filter.Limit = v
			}
		}
		if offsetStr := q.Get("offset"); offsetStr != "" {
			if v, err := strconv.Atoi(offsetStr); err == nil {
				filter.Offset = v
			}
		}

		summaries, err := gameSvc.ListGames(r.Context(), filter)
		if err != nil {
			http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
			return
		}

		resp := listGamesResponse{
			Games: make([]gameSummaryDTO, 0, len(summaries)),
		}

		for _, g := range summaries {
			dto := gameSummaryDTO{
				GameID:    g.ID,
				Mode:      string(g.Mode),
				Status:    string(g.Status),
				CreatedAt: g.CreatedAt.Format(time.RFC3339),
			}
			dto.CreatedBy.PlayerID = g.CreatedByPlayerID
			dto.CreatedBy.Name = g.CreatedByPlayerName
			resp.Games = append(resp.Games, dto)
		}

		w.Header().Set("Content-Type", "application/json")
		_ = json.NewEncoder(w).Encode(resp)
	}
}

func JoinGameHandler(gameSvc service.GameService) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		playerID := PlayerIDFromContext(r.Context())
		if playerID == "" {
			http.Error(w, "missing X-Player-Id header", http.StatusBadRequest)
			return
		}

		gameID := chi.URLParam(r, "gameId")
		if gameID == "" {
			http.Error(w, "missing gameId", http.StatusBadRequest)
			return
		}

		gameState, err := gameSvc.JoinGame(r.Context(), gameID, playerID)
		if err != nil {
			switch {
			case errors.Is(err, service.ErrInvalidGameState):
				http.Error(w, "invalid game state", http.StatusBadRequest)
				return
			case errors.Is(err, store.ErrGameNotFound):
				http.Error(w, http.StatusText(http.StatusNotFound), http.StatusNotFound)
				return
			default:
				http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
				return
			}
		}

		// build response using existing helper-style logic
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
		_ = json.NewEncoder(w).Encode(resp)
	}
}

func GetGameHandler(gameSvc service.GameService) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		gameID := chi.URLParam(r, "gameId")
		if gameID == "" {
			http.Error(w, "missing gameId", http.StatusBadRequest)
			return
		}

		gameState, err := gameSvc.GetGame(r.Context(), gameID)
		if err != nil {
			if errors.Is(err, store.ErrGameNotFound) {
				http.Error(w, http.StatusText(http.StatusNotFound), http.StatusNotFound)
				return
			}
			http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
			return
		}

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
		_ = json.NewEncoder(w).Encode(resp)
	}
}

func MakeMoveHandler(gameSvc service.GameService) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		playerID := PlayerIDFromContext(r.Context())
		if playerID == "" {
			http.Error(w, "missing X-Player-Id header", http.StatusBadRequest)
			return
		}

		gameID := chi.URLParam(r, "gameId")
		if gameID == "" {
			http.Error(w, "missing gameId", http.StatusBadRequest)
			return
		}

		var req makeMoveRequest
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			http.Error(w, http.StatusText(http.StatusBadRequest), http.StatusBadRequest)
			return
		}

		gameState, err := gameSvc.MakeMove(r.Context(), gameID, playerID, req.Row, req.Col)
		if err != nil {
			switch {
			case errors.Is(err, store.ErrGameNotFound):
				http.Error(w, http.StatusText(http.StatusNotFound), http.StatusNotFound)
				return
			case errors.Is(err, service.ErrNotParticipant):
				http.Error(w, http.StatusText(http.StatusForbidden), http.StatusForbidden)
				return
			case errors.Is(err, service.ErrNotPlayersTurn),
				errors.Is(err, service.ErrInvalidMove),
				errors.Is(err, service.ErrInvalidGameState):
				http.Error(w, http.StatusText(http.StatusBadRequest), http.StatusBadRequest)
				return
			default:
				http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
				return
			}
		}

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
		_ = json.NewEncoder(w).Encode(resp)
	}
}

// WebSocketHandler handles WebSocket connections for game updates.
func WebSocketHandler(hub *ws.Hub, gameSvc service.GameService) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		gameID := chi.URLParam(r, "gameId")
		if gameID == "" {
			http.Error(w, "missing gameId", http.StatusBadRequest)
			return
		}

		// Verify game exists
		gameState, err := gameSvc.GetGame(r.Context(), gameID)
		if err != nil {
			if errors.Is(err, store.ErrGameNotFound) {
				http.Error(w, http.StatusText(http.StatusNotFound), http.StatusNotFound)
				return
			}
			http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
			return
		}

		// Upgrade HTTP connection to WebSocket
		upgrader := websocket.Upgrader{
			ReadBufferSize:  1024,
			WriteBufferSize: 1024,
			CheckOrigin: func(r *http.Request) bool {
				// Allow all origins for development
				return true
			},
		}

		conn, err := upgrader.Upgrade(w, r, nil)
		if err != nil {
			log.Printf("websocket upgrade error: %v", err)
			return
		}

		// Create connection wrapper
		wsConn := ws.NewConnection(hub, conn)
		hub.Register(gameID, wsConn)

		// Send initial game state immediately
		hub.BroadcastGameState(gameID, gameState)

		// Start connection pumps
		go wsConn.WritePump()
		go wsConn.ReadPump()
	}
}
