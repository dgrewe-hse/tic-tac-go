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

package game

import (
	"testing"

	"tic-tac-go/internal/models"
)

func TestCheckWinner_RowWin(t *testing.T) {
	board := NewBoard()
	board[0][0] = models.SymbolX
	board[0][1] = models.SymbolX
	board[0][2] = models.SymbolX

	winner, isDraw := CheckWinner(board)
	if winner != models.SymbolX || isDraw {
		t.Fatalf("expected winner X and not draw, got winner=%q isDraw=%v", winner, isDraw)
	}
}

func TestCheckWinner_Draw(t *testing.T) {
	// Full board with no winner:
	// X O X
	// X O O
	// O X X
	board := NewBoard()
	board[0][0] = models.SymbolX
	board[0][1] = models.SymbolO
	board[0][2] = models.SymbolX

	board[1][0] = models.SymbolX
	board[1][1] = models.SymbolO
	board[1][2] = models.SymbolO

	board[2][0] = models.SymbolO
	board[2][1] = models.SymbolX
	board[2][2] = models.SymbolX

	winner, isDraw := CheckWinner(board)
	if winner != models.SymbolEmpty || !isDraw {
		t.Fatalf("expected draw with no winner, got winner=%q isDraw=%v", winner, isDraw)
	}
}
