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
	"testing"

	"tic-tac-go/internal/models"
	"tic-tac-go/internal/store"
)

func TestPlayerService_CreateAndGet(t *testing.T) {
	ctx := context.Background()
	playerStore := store.NewMemoryPlayerStore()
	svc := NewPlayerService(playerStore)

	created, err := svc.CreatePlayer(ctx, "Alice")
	if err != nil {
		t.Fatalf("CreatePlayer error = %v, want nil", err)
	}
	if created.ID == "" {
		t.Fatalf("expected non-empty ID")
	}
	if created.Name != "Alice" {
		t.Fatalf("expected name Alice, got %q", created.Name)
	}

	got, err := svc.GetPlayer(ctx, created.ID)
	if err != nil {
		t.Fatalf("GetPlayer error = %v, want nil", err)
	}
	if got.ID != created.ID || got.Name != created.Name {
		t.Fatalf("GetPlayer = %+v, want %+v", got, created)
	}
}

func TestGameService_CreateGame_PVPAndPVC(t *testing.T) {
	ctx := context.Background()
	gameStore := store.NewMemoryGameStore()
	playerStore := store.NewMemoryPlayerStore()

	player := &models.Player{ID: "p1", Name: "Alice"}
	_ = playerStore.Create(player)

	svc := NewGameService(gameStore, playerStore)

	// PVP
	gamePVP, err := svc.CreateGame(ctx, "p1", models.GameModePVP)
	if err != nil {
		t.Fatalf("CreateGame PVP error = %v", err)
	}
	if gamePVP.Mode != models.GameModePVP {
		t.Fatalf("expected mode PVP, got %q", gamePVP.Mode)
	}
	if gamePVP.Status != models.GameStatusWaitingForPlayer {
		t.Fatalf("expected status WAITING_FOR_PLAYER, got %q", gamePVP.Status)
	}

	// PVC
	gamePVC, err := svc.CreateGame(ctx, "p1", models.GameModePVC)
	if err != nil {
		t.Fatalf("CreateGame PVC error = %v", err)
	}
	if gamePVC.Mode != models.GameModePVC {
		t.Fatalf("expected mode PVC, got %q", gamePVC.Mode)
	}
	if gamePVC.Status != models.GameStatusInProgress {
		t.Fatalf("expected status IN_PROGRESS, got %q", gamePVC.Status)
	}
	if gamePVC.PlayerOID != "AI" {
		t.Fatalf("expected PlayerOID AI, got %q", gamePVC.PlayerOID)
	}
}

func TestGameService_JoinGame_PVP(t *testing.T) {
	ctx := context.Background()
	gameStore := store.NewMemoryGameStore()
	playerStore := store.NewMemoryPlayerStore()

	creator := &models.Player{ID: "creator", Name: "Alice"}
	joiner := &models.Player{ID: "joiner", Name: "Bob"}
	_ = playerStore.Create(creator)
	_ = playerStore.Create(joiner)

	svc := NewGameService(gameStore, playerStore)

	gameState, err := svc.CreateGame(ctx, "creator", models.GameModePVP)
	if err != nil {
		t.Fatalf("CreateGame error = %v", err)
	}

	joined, err := svc.JoinGame(ctx, gameState.ID, "joiner")
	if err != nil {
		t.Fatalf("JoinGame error = %v", err)
	}

	if joined.PlayerOID != "joiner" {
		t.Fatalf("expected PlayerOID joiner, got %q", joined.PlayerOID)
	}
	if joined.Status != models.GameStatusInProgress {
		t.Fatalf("expected status IN_PROGRESS, got %q", joined.Status)
	}
}

func TestGameService_MakeMove_PVP_FirstMove(t *testing.T) {
	ctx := context.Background()
	gameStore := store.NewMemoryGameStore()
	playerStore := store.NewMemoryPlayerStore()

	creator := &models.Player{ID: "pX", Name: "Alice"}
	joiner := &models.Player{ID: "pO", Name: "Bob"}
	_ = playerStore.Create(creator)
	_ = playerStore.Create(joiner)

	svc := NewGameService(gameStore, playerStore)

	gameState, err := svc.CreateGame(ctx, "pX", models.GameModePVP)
	if err != nil {
		t.Fatalf("CreateGame error = %v", err)
	}

	_, err = svc.JoinGame(ctx, gameState.ID, "pO")
	if err != nil {
		t.Fatalf("JoinGame error = %v", err)
	}

	updated, err := svc.MakeMove(ctx, gameState.ID, "pX", 0, 0)
	if err != nil {
		t.Fatalf("MakeMove error = %v", err)
	}
	if updated.Board[0][0] != models.SymbolX {
		t.Fatalf("expected X at (0,0), got %q", updated.Board[0][0])
	}
}

func TestGameService_ListGames(t *testing.T) {
	ctx := context.Background()
	gameStore := store.NewMemoryGameStore()
	playerStore := store.NewMemoryPlayerStore()

	player := &models.Player{ID: "p1", Name: "Alice"}
	_ = playerStore.Create(player)

	svc := NewGameService(gameStore, playerStore)

	// Two PVP games, one waiting, one in progress.
	g1, _ := svc.CreateGame(ctx, "p1", models.GameModePVP)
	_, _ = svc.CreateGame(ctx, "p1", models.GameModePVP)

	// Make g1 IN_PROGRESS by joining.
	other := &models.Player{ID: "p2", Name: "Bob"}
	_ = playerStore.Create(other)
	_, _ = svc.JoinGame(ctx, g1.ID, "p2")

	mode := models.GameModePVP
	status := models.GameStatusWaitingForPlayer

	summaries, err := svc.ListGames(ctx, store.GameFilter{
		Mode:   &mode,
		Status: &status,
	})
	if err != nil {
		t.Fatalf("ListGames error = %v", err)
	}
	if len(summaries) != 1 {
		t.Fatalf("expected 1 waiting PVP game, got %d", len(summaries))
	}
	if summaries[0].CreatedByPlayerID != "p1" {
		t.Fatalf("expected CreatedByPlayerID p1, got %q", summaries[0].CreatedByPlayerID)
	}
}
