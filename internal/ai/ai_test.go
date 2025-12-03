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
	"testing"

	"tic-tac-go/internal/game"
	"tic-tac-go/internal/models"
)

func TestChooseMove_TakesWinningMove(t *testing.T) {
	board := game.NewBoard()
	// AI is X and already has two in a row on the top row: X X _
	board[0][0] = models.SymbolX
	board[0][1] = models.SymbolX

	row, col := ChooseMove(board, models.SymbolX, models.SymbolO)

	if row != 0 || col != 2 {
		t.Fatalf("expected AI to win at (0,2), got (%d,%d)", row, col)
	}
}

func TestChooseMove_BlocksOpponentWinningMove(t *testing.T) {
	board := game.NewBoard()
	// Opponent is X, AI is O, opponent has two in a row on the first column:
	// X
	// X
	// _
	board[0][0] = models.SymbolX
	board[1][0] = models.SymbolX

	row, col := ChooseMove(board, models.SymbolO, models.SymbolX)

	if row != 2 || col != 0 {
		t.Fatalf("expected AI to block at (2,0), got (%d,%d)", row, col)
	}
}

func TestChooseMove_TakesCenterIfFree(t *testing.T) {
	board := game.NewBoard()

	row, col := ChooseMove(board, models.SymbolX, models.SymbolO)

	if row != 1 || col != 1 {
		t.Fatalf("expected AI to take center (1,1), got (%d,%d)", row, col)
	}
}
