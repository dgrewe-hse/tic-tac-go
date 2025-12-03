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

package models

import "time"

// Board is a fixed 3x3 tic-tac-toe game board
type Board [3][3]Symbol

// Player represents a player in the game
type Player struct {
	ID   string `json:"id"`
	Name string `json:"name"`
}

// GameMode describes whether a game is player-vs-player or player-vs-computer
type GameMode string

const (
	GameModePVP GameMode = "PVP"
	GameModePVC GameMode = "PVC"
)

// GameStatus represents the lifecycle state of a game
type GameStatus string

const (
	GameStatusWaitingForPlayer GameStatus = "WAITING_FOR_PLAYER"
	GameStatusInProgress       GameStatus = "IN_PROGRESS"
	GameStatusFinished         GameStatus = "FINISHED"
)

// Symbol represents the mark used by players on the board
type Symbol string

const (
	SymbolEmpty Symbol = ""
	SymbolX     Symbol = "X"
	SymbolO     Symbol = "O"
)

// GameState holds the full state of a single tic-tac-toe game
type GameState struct {
	ID          string     `json:"id"`
	Mode        GameMode   `json:"mode"`
	Board       Board      `json:"board"`
	PlayerXID   string     `json:"playerXId"`
	PlayerOID   string     `json:"playerOId"`
	CurrentTurn Symbol     `json:"currentTurn"`
	Status      GameStatus `json:"status"`
	// Winner can be "X", "O", "DRAW", or "" (no winner yet)
	Winner    string    `json:"winner"`
	CreatedAt time.Time `json:"createdAt"`
	UpdatedAt time.Time `json:"updatedAt"`
}

// GameSummary is a lightweight representation used when listing games (e.g., in the lobby)
type GameSummary struct {
	ID                  string     `json:"id"`
	Mode                GameMode   `json:"mode"`
	Status              GameStatus `json:"status"`
	CreatedAt           time.Time  `json:"createdAt"`
	CreatedByPlayerID   string     `json:"createdByPlayerId"`
	CreatedByPlayerName string     `json:"createdByPlayerName"`
}
