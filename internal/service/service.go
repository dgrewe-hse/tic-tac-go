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

package service

import (
	"context"
	"errors"

	"tic-tac-go/internal/models"
	"tic-tac-go/internal/store"
)

// some service layer error definitions
var (
	ErrInvalidGameMode  = errors.New("invalid game mode")
	ErrInvalidGameState = errors.New("invalid game state")
	ErrNotParticipant   = errors.New("player is not a participant in this game")
	ErrNotPlayersTurn   = errors.New("it is not this player's turn")
	ErrInvalidMove      = errors.New("invalid move")
)

// GameService defines the high-level use-cases for managing games
type GameService interface {
	CreateGame(ctx context.Context, creatorPlayerID string, mode models.GameMode) (*models.GameState, error)
	JoinGame(ctx context.Context, gameID, playerID string) (*models.GameState, error)
	GetGame(ctx context.Context, gameID string) (*models.GameState, error)
	MakeMove(ctx context.Context, gameID, playerID string, row, col int) (*models.GameState, error)
	ListGames(ctx context.Context, filter store.GameFilter) ([]*models.GameSummary, error)
}

// PlayerService defines use-cases for managing players
type PlayerService interface {
	CreatePlayer(ctx context.Context, name string) (*models.Player, error)
	GetPlayer(ctx context.Context, id string) (*models.Player, error)
}

// GameStateBroadcaster defines an interface for broadcasting game state updates.
// This allows the service layer to notify WebSocket clients without directly depending on the WebSocket implementation.
type GameStateBroadcaster interface {
	BroadcastGameState(gameID string, state *models.GameState)
}
