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

import "tic-tac-go/internal/models"

// CheckWinner checks the board and returns:
//   - winner: "X" or "O" if someone has three in a row,
//     models.SymbolEmpty ("") otherwise.
//   - isDraw: true if the board is full and there is no winner.
func CheckWinner(board models.Board) (winner models.Symbol, isDraw bool) {
	// 1. Check rows and columns
	for i := 0; i < 3; i++ {
		// Row i
		if board[i][0] != models.SymbolEmpty &&
			board[i][0] == board[i][1] &&
			board[i][1] == board[i][2] {
			return board[i][0], false
		}
		// Column i
		if board[0][i] != models.SymbolEmpty &&
			board[0][i] == board[1][i] &&
			board[1][i] == board[2][i] {
			return board[0][i], false
		}
	}

	// 2. Diagonals
	if board[0][0] != models.SymbolEmpty &&
		board[0][0] == board[1][1] &&
		board[1][1] == board[2][2] {
		return board[0][0], false
	}
	if board[0][2] != models.SymbolEmpty &&
		board[0][2] == board[1][1] &&
		board[1][1] == board[2][0] {
		return board[0][2], false
	}

	// 3. Draw?
	if IsFull(board) {
		return models.SymbolEmpty, true
	}

	// 4. Game still in progress
	return models.SymbolEmpty, false
}
