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

func TestIsValidMove(t *testing.T) {
	board := NewBoard()

	tests := []struct {
		name string
		row  int
		col  int
		want bool
	}{
		{"inside empty cell", 0, 0, true},
		{"negative row", -1, 0, false},
		{"row too large", 3, 0, false},
		{"negative col", 0, -1, false},
		{"col too large", 0, 3, false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := IsValidMove(board, tt.row, tt.col)
			if got != tt.want {
				t.Fatalf("IsValidMove(%d,%d) = %v, want %v", tt.row, tt.col, got, tt.want)
			}
		})
	}
}

func TestApplyMove_ValidAndInvalid(t *testing.T) {
	board := NewBoard()

	// valid move
	newBoard, err := ApplyMove(board, 1, 1, models.SymbolX)
	if err != nil {
		t.Fatalf("expected no error for valid move, got %v", err)
	}
	if newBoard[1][1] != models.SymbolX {
		t.Fatalf("expected symbol X at (1,1), got %q", newBoard[1][1])
	}

	// invalid move (same place again)
	_, err = ApplyMove(newBoard, 1, 1, models.SymbolO)
	if err == nil {
		t.Fatalf("expected error for invalid move, got nil")
	}
}
