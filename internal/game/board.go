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
	"fmt"

	"tic-tac-go/internal/models"
)

// OppositeSymbol function returns the other player's symbol (so X <--> O)
func OppositeSymbol(s models.Symbol) models.Symbol {
	switch s {
	case models.SymbolX:
		return models.SymbolO
	case models.SymbolO:
		return models.SymbolX
	default:
		return models.SymbolEmpty
	}
}

// NewBoard creates a new empty 3x3 game board
func NewBoard() models.Board {
	var board models.Board
	for row := 0; row < 3; row++ {
		for col := 0; col < 3; col++ {
			board[row][col] = models.SymbolEmpty
		}
	}
	return board
}

// IsValidMove reports whether a move by a player is within the boiunds and on an empty cell of the board
func IsValidMove(board models.Board, row, col int) bool {
	if row < 0 || row >= 3 || col < 0 || col >= 3 {
		return false
	}
	return board[row][col] == models.SymbolEmpty // check if selected cell is empty
}

// ApplyMove returns a new board with the given symvol placed at (row, col)
// If the move is invalid, it returns an error.
func ApplyMove(board models.Board, row, col int, symbol models.Symbol) (models.Board, error) {
	if !IsValidMove(board, row, col) {
		return board, fmt.Errorf("invalid move at row=%d col=%d", row, col)
	}

	newBoard := board
	newBoard[row][col] = symbol
	return newBoard, nil
}

// IsFull reports whether the board as no empty cells left
func IsFull(board models.Board) bool {
	for row := 0; row < 3; row++ {
		for col := 0; col < 3; col++ {
			if board[row][col] == models.SymbolEmpty {
				return false
			}
		}
	}
	return true
}

// AvailableMoves return a slice of all empty positions as [row, col] pairs
func AvailableMoves(board models.Board) [][2]int {
	var moves [][2]int
	for row := 0; row < 3; row++ {
		for col := 0; col < 3; col++ {
			if board[row][col] == models.SymbolEmpty {
				moves = append(moves, [2]int{row, col})
			}
		}
	}
	return moves
}
