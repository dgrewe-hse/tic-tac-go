// Copyright 2026 Esslingen University of Applied Sciences
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
// Date: 2026-01-05

package ws

import (
	"encoding/json"
	"sync"

	"tic-tac-go/internal/models"
)

// Hub manages WebSocket connections grouped by game ID.
// It is safe for concurrent use.
type Hub struct {
	// clients maps gameID -> set of connections for that game
	clients map[string]map[*Connection]struct{}
	mu      sync.RWMutex

	// register channel for new connections
	register chan *Connection

	// unregister channel for disconnected clients
	unregister chan *Connection
}

// NewHub creates a new WebSocket hub.
func NewHub() *Hub {
	return &Hub{
		clients:    make(map[string]map[*Connection]struct{}),
		register:   make(chan *Connection),
		unregister: make(chan *Connection),
	}
}

// Run starts the hub's event loop, processing register/unregister events.
func (h *Hub) Run() {
	for {
		select {
		case conn := <-h.register:
			h.mu.Lock()
			if h.clients[conn.gameID] == nil {
				h.clients[conn.gameID] = make(map[*Connection]struct{})
			}
			h.clients[conn.gameID][conn] = struct{}{}
			h.mu.Unlock()

		case conn := <-h.unregister:
			h.mu.Lock()
			if clients, ok := h.clients[conn.gameID]; ok {
				delete(clients, conn)
				if len(clients) == 0 {
					delete(h.clients, conn.gameID)
				}
			}
			close(conn.send)
			h.mu.Unlock()
		}
	}
}

// Register adds a connection to the hub for a specific game.
func (h *Hub) Register(gameID string, conn *Connection) {
	conn.gameID = gameID
	h.register <- conn
}

// Unregister removes a connection from the hub.
func (h *Hub) Unregister(conn *Connection) {
	h.unregister <- conn
}

// BroadcastGameState sends the game state to all connections registered for the given gameID.
func (h *Hub) BroadcastGameState(gameID string, state *models.GameState) {
	h.mu.RLock()
	clients, ok := h.clients[gameID]
	if !ok {
		h.mu.RUnlock()
		return
	}

	// Convert board to [][]string for JSON
	board := make([][]string, 3)
	for i := 0; i < 3; i++ {
		board[i] = make([]string, 3)
		for j := 0; j < 3; j++ {
			board[i][j] = string(state.Board[i][j])
		}
	}

	message := map[string]interface{}{
		"type": "state",
		"payload": map[string]interface{}{
			"gameId":      state.ID,
			"board":       board,
			"currentTurn": string(state.CurrentTurn),
			"status":      string(state.Status),
			"winner":      state.Winner,
		},
	}

	msgBytes, err := json.Marshal(message)
	if err != nil {
		h.mu.RUnlock()
		return
	}

	// Copy clients slice to avoid holding lock while sending
	conns := make([]*Connection, 0, len(clients))
	for conn := range clients {
		conns = append(conns, conn)
	}
	h.mu.RUnlock()

	// Send to all connections (non-blocking)
	for _, conn := range conns {
		select {
		case conn.send <- msgBytes:
		default:
			// If send buffer is full, skip this connection
			// (it will be cleaned up on next unregister)
		}
	}
}

// BroadcastError sends an error message to a specific connection.
func (h *Hub) BroadcastError(conn *Connection, message string) {
	msg := map[string]interface{}{
		"type": "error",
		"payload": map[string]interface{}{
			"message": message,
		},
	}

	msgBytes, err := json.Marshal(msg)
	if err != nil {
		return
	}

	select {
	case conn.send <- msgBytes:
	default:
		// Send buffer full, skip
	}
}
