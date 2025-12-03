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

package store

import (
	"testing"

	"tic-tac-go/internal/models"
)

func TestMemoryPlayerStore_CreateAndGet(t *testing.T) {
	s := NewMemoryPlayerStore()

	player := &models.Player{
		ID:   "player-1",
		Name: "Alice",
	}

	if err := s.Create(player); err != nil {
		t.Fatalf("Create() error = %v, want nil", err)
	}

	got, err := s.Get("player-1")
	if err != nil {
		t.Fatalf("Get() error = %v, want nil", err)
	}
	if got.ID != player.ID || got.Name != player.Name {
		t.Fatalf("Get() = %+v, want %+v", got, player)
	}
}

func TestMemoryPlayerStore_Get_NotFound(t *testing.T) {
	s := NewMemoryPlayerStore()

	_, err := s.Get("does-not-exist")
	if err != ErrPlayerNotFound {
		t.Fatalf("expected ErrPlayerNotFound, got %v", err)
	}
}

func TestMemoryGameStore_CreateGetUpdate(t *testing.T) {
	s := NewMemoryGameStore()

	game := &models.GameState{
		ID:    "game-1",
		Mode:  models.GameModePVP,
		Board: models.Board{},
	}

	if err := s.Create(game); err != nil {
		t.Fatalf("Create() error = %v, want nil", err)
	}

	got, err := s.Get("game-1")
	if err != nil {
		t.Fatalf("Get() error = %v, want nil", err)
	}
	if got.ID != game.ID {
		t.Fatalf("Get().ID = %q, want %q", got.ID, game.ID)
	}

	// Update
	game.Status = models.GameStatusInProgress
	if err := s.Update(game); err != nil {
		t.Fatalf("Update() error = %v, want nil", err)
	}

	got, err = s.Get("game-1")
	if err != nil {
		t.Fatalf("Get() after Update error = %v, want nil", err)
	}
	if got.Status != models.GameStatusInProgress {
		t.Fatalf("expected status IN_PROGRESS, got %q", got.Status)
	}
}

func TestMemoryGameStore_Get_NotFound(t *testing.T) {
	s := NewMemoryGameStore()

	_, err := s.Get("missing")
	if err != ErrGameNotFound {
		t.Fatalf("expected ErrGameNotFound, got %v", err)
	}
}

func TestMemoryGameStore_List_WithFiltersAndPaging(t *testing.T) {
	s := NewMemoryGameStore()

	game1 := &models.GameState{ID: "g1", Mode: models.GameModePVP, Status: models.GameStatusWaitingForPlayer}
	game2 := &models.GameState{ID: "g2", Mode: models.GameModePVP, Status: models.GameStatusInProgress}
	game3 := &models.GameState{ID: "g3", Mode: models.GameModePVC, Status: models.GameStatusWaitingForPlayer}

	_ = s.Create(game1)
	_ = s.Create(game2)
	_ = s.Create(game3)

	mode := models.GameModePVP
	status := models.GameStatusWaitingForPlayer

	// Filter by mode + status
	games, err := s.List(GameFilter{
		Mode:   &mode,
		Status: &status,
	})
	if err != nil {
		t.Fatalf("List() error = %v, want nil", err)
	}
	if len(games) != 1 || games[0].ID != "g1" {
		t.Fatalf("expected [g1], got %#v", games)
	}

	// Paging: limit 1, offset 1 (no filters)
	all, err := s.List(GameFilter{Limit: 1, Offset: 1})
	if err != nil {
		t.Fatalf("List() with paging error = %v", err)
	}
	if len(all) != 1 {
		t.Fatalf("expected 1 game with limit=1 offset=1, got %d", len(all))
	}
}
