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

package ai

import (
	"math/rand"
	"time"

	"tic-tac-go/internal/game"
	"tic-tac-go/internal/models"
)

// ChooseMove picks the next move for the AI player based on a simple heuristic:
// 1) win if possible, 2) block opponent, 3) take center, 4) pick a random free cell.
func ChooseMove(board models.Board, aiSymbol, opponentSymbol models.Symbol) (row, col int) {
	// 1. Try to win.
	for _, move := range game.AvailableMoves(board) {
		r, c := move[0], move[1]
		b, _ := game.ApplyMove(board, r, c, aiSymbol)
		winner, _ := game.CheckWinner(b)
		if winner == aiSymbol {
			return r, c
		}
	}

	// 2. Block opponent's winning move.
	for _, move := range game.AvailableMoves(board) {
		r, c := move[0], move[1]
		b, _ := game.ApplyMove(board, r, c, opponentSymbol)
		winner, _ := game.CheckWinner(b)
		if winner == opponentSymbol {
			return r, c
		}
	}

	// 3. Take center if free.
	if game.IsValidMove(board, 1, 1) {
		return 1, 1
	}

	// 4. Random available move.
	moves := game.AvailableMoves(board)
	if len(moves) == 0 {
		return -1, -1 // should not happen for a valid in-progress game
	}

	// Use a local rand source so tests are deterministic if needed.
	rng := rand.New(rand.NewSource(time.Now().UnixNano()))
	choice := moves[rng.Intn(len(moves))]
	return choice[0], choice[1]
}
