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
	"tic-tac-go/internal/ai"
	"tic-tac-go/internal/game"
	"tic-tac-go/internal/models"
	"tic-tac-go/internal/store"
	"time"

	"github.com/google/uuid"
)

// gameService is a concrete implementation of GameService.
type gameService struct {
	gameStore   store.GameStore
	playerStore store.PlayerStore
}

// NewGameService constructs a GameService with the given dependencies.
func NewGameService(gameStore store.GameStore, playerStore store.PlayerStore) GameService {
	return &gameService{
		gameStore:   gameStore,
		playerStore: playerStore,
	}
}

// CreateGame creates a new game in either PVP or PVC mode.
func (s *gameService) CreateGame(ctx context.Context, creatorPlayerID string, mode models.GameMode) (*models.GameState, error) {
	// Ensure creator exists.
	if _, err := s.playerStore.Get(creatorPlayerID); err != nil {
		return nil, err
	}

	if mode != models.GameModePVP && mode != models.GameModePVC {
		return nil, ErrInvalidGameMode
	}

	now := time.Now().UTC()

	gameState := &models.GameState{
		ID:        uuid.NewString(),
		Mode:      mode,
		Board:     game.NewBoard(),
		PlayerXID: creatorPlayerID,
		Status:    models.GameStatusInProgress,
		Winner:    "",
		CreatedAt: now,
		UpdatedAt: now,
	}

	// For PVP, wait for second player.
	if mode == models.GameModePVP {
		gameState.Status = models.GameStatusWaitingForPlayer
	}

	// For PVC, PlayerO is the AI.
	if mode == models.GameModePVC {
		gameState.PlayerOID = "AI"
		gameState.CurrentTurn = models.SymbolX
	} else {
		gameState.CurrentTurn = models.SymbolX
	}

	if err := s.gameStore.Create(gameState); err != nil {
		return nil, err
	}

	return gameState, nil
}

// GetGame returns the current state of a game by ID.
func (s *gameService) GetGame(ctx context.Context, gameID string) (*models.GameState, error) {
	return s.gameStore.Get(gameID)
}

func (s *gameService) JoinGame(ctx context.Context, gameID, playerID string) (*models.GameState, error) {
	// Ensure player exists.
	if _, err := s.playerStore.Get(playerID); err != nil {
		return nil, err
	}

	gameState, err := s.gameStore.Get(gameID)
	if err != nil {
		return nil, err
	}

	// Game must be PVP and waiting.
	if gameState.Mode != models.GameModePVP || gameState.Status != models.GameStatusWaitingForPlayer {
		return nil, ErrInvalidGameState
	}

	// Joining player becomes O.
	gameState.PlayerOID = playerID
	gameState.Status = models.GameStatusInProgress
	gameState.CurrentTurn = models.SymbolX
	gameState.UpdatedAt = time.Now().UTC()

	if err := s.gameStore.Update(gameState); err != nil {
		return nil, err
	}

	return gameState, nil
}

func (s *gameService) MakeMove(ctx context.Context, gameID, playerID string, row, col int) (*models.GameState, error) {
	// Load game.
	gameState, err := s.gameStore.Get(gameID)
	if err != nil {
		return nil, err
	}

	// Game must be in progress.
	if gameState.Status != models.GameStatusInProgress {
		return nil, ErrInvalidGameState
	}

	// Determine symbol for this player and ensure they are a participant.
	var symbol models.Symbol
	var opponentSymbol models.Symbol

	if playerID == gameState.PlayerXID {
		symbol = models.SymbolX
		opponentSymbol = models.SymbolO
	} else if playerID == gameState.PlayerOID {
		symbol = models.SymbolO
		opponentSymbol = models.SymbolX
	} else {
		return nil, ErrNotParticipant
	}

	// Enforce turn order.
	if gameState.CurrentTurn != symbol {
		return nil, ErrNotPlayersTurn
	}

	// Validate move.
	if !game.IsValidMove(gameState.Board, row, col) {
		return nil, ErrInvalidMove
	}

	// Apply player's move.
	newBoard, err := game.ApplyMove(gameState.Board, row, col, symbol)
	if err != nil {
		return nil, ErrInvalidMove
	}
	gameState.Board = newBoard

	// Check winner / draw after player's move.
	winner, isDraw := game.CheckWinner(gameState.Board)
	if winner != models.SymbolEmpty {
		gameState.Status = models.GameStatusFinished
		gameState.Winner = string(winner)
	} else if isDraw {
		gameState.Status = models.GameStatusFinished
		gameState.Winner = "DRAW"
	} else {
		// Switch turn.
		gameState.CurrentTurn = game.OppositeSymbol(symbol)
	}

	// If PVC and still in progress and it's AI's turn, let AI move.
	if gameState.Mode == models.GameModePVC &&
		gameState.Status == models.GameStatusInProgress &&
		gameState.CurrentTurn == opponentSymbol {

		aiRow, aiCol := ai.ChooseMove(gameState.Board, opponentSymbol, symbol)

		aiBoard, err := game.ApplyMove(gameState.Board, aiRow, aiCol, opponentSymbol)
		if err == nil {
			gameState.Board = aiBoard

			winner, isDraw = game.CheckWinner(gameState.Board)
			if winner != models.SymbolEmpty {
				gameState.Status = models.GameStatusFinished
				gameState.Winner = string(winner)
			} else if isDraw {
				gameState.Status = models.GameStatusFinished
				gameState.Winner = "DRAW"
			} else {
				// Back to human.
				gameState.CurrentTurn = symbol
			}
		}
	}

	gameState.UpdatedAt = time.Now().UTC()

	if err := s.gameStore.Update(gameState); err != nil {
		return nil, err
	}

	return gameState, nil
}

func (s *gameService) ListGames(ctx context.Context, filter store.GameFilter) ([]*models.GameSummary, error) {
	games, err := s.gameStore.List(filter)
	if err != nil {
		return nil, err
	}

	var summaries []*models.GameSummary

	for _, g := range games {
		summary := &models.GameSummary{
			ID:        g.ID,
			Mode:      g.Mode,
			Status:    g.Status,
			CreatedAt: g.CreatedAt,
		}

		// Enrich with creator info if available.
		if g.PlayerXID != "" {
			player, err := s.playerStore.Get(g.PlayerXID)
			if err == nil {
				summary.CreatedByPlayerID = player.ID
				summary.CreatedByPlayerName = player.Name
			}
		}

		summaries = append(summaries, summary)
	}

	return summaries, nil
}
