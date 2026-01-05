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
	"net/http"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"

	"tic-tac-go/internal/service"
	"tic-tac-go/internal/store"
	"tic-tac-go/internal/ws"
)

// NewRouter constructs the root HTTP router for the Tic-Tac-Go server.
// For now it only exposes a simple health endpoint; additional routes
// for game and player APIs will be added later.
func NewRouter() http.Handler {
	r := chi.NewRouter()

	// In-memory stores for players and games.
	playerStore := store.NewMemoryPlayerStore()
	gameStore := store.NewMemoryGameStore()

	// WebSocket hub for real-time game updates
	hub := ws.NewHub()
	go hub.Run()

	// Services using the stores.
	playerSvc := service.NewPlayerService(playerStore)
	// GameService with WebSocket broadcaster
	gameSvc := service.NewGameServiceWithBroadcaster(gameStore, playerStore, hub)

	// Basic middlewares for logging, recovery and timeouts.
	r.Use(middleware.RequestID)
	r.Use(middleware.RealIP)
	r.Use(middleware.Logger)
	r.Use(middleware.Recoverer)
	r.Use(middleware.Timeout(60 * time.Second))
	r.Use(WithPlayerID)

	// Health check endpoint to verify the server is up.
	r.Get("/health", healthHandler)
	// get a list of created games
	r.Get("/games", ListGamesHandler(gameSvc))
	// get existing game by id
	r.Get("/games/{gameId}", GetGameHandler(gameSvc))

	// Player endpoints.
	r.Post("/players", CreatePlayerHandler(playerSvc))
	// Game endpoints.
	r.Post("/games", CreateGameHandler(gameSvc))
	// join existing game by id
	r.Post("/games/{gameId}/join", JoinGameHandler(gameSvc))
	// make move within existing game
	r.Post("/games/{gameId}/moves", MakeMoveHandler(gameSvc))

	// WebSocket endpoint for real-time game updates
	r.Get("/ws/games/{gameId}", WebSocketHandler(hub, gameSvc))

	return r
}
