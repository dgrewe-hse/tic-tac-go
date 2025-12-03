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
	"sync"
	"tic-tac-go/internal/models"
)

// MemoryPlayerStore is an in-memory implementation of PlayerStore.
// It is safe for concurrent use using mutex sync mechanism.
type MemoryPlayerStore struct {
	mu      sync.RWMutex // Read/Write Mutex
	players map[string]*models.Player
}

// NewMemoryPlayerStore constructs a new empty MemoryPlayerStore
func NewMemoryPlayerStore() *MemoryPlayerStore {
	return &MemoryPlayerStore{
		players: make(map[string]*models.Player),
	}
}

// Create a new player and store it
func (s *MemoryPlayerStore) Create(player *models.Player) error {
	s.mu.Lock() // Lock access to player store
	defer s.mu.Unlock()

	s.players[player.ID] = player
	return nil
}

// Lookup of players within the list using its id
func (s *MemoryPlayerStore) Get(id string) (*models.Player, error) {
	s.mu.RLock() // read lock to the player store
	defer s.mu.RUnlock()

	player, ok := s.players[id]
	if !ok {
		return nil, ErrPlayerNotFound
	}

	return player, nil
}

// MemoryGameStore is an in-memory implementation of GameStore.
// It is safe for concurrent use.
type MemoryGameStore struct {
	mu    sync.RWMutex
	games map[string]*models.GameState
}

// NewMemoryGameStore constructs a new empty MemoryGameStore.
func NewMemoryGameStore() *MemoryGameStore {
	return &MemoryGameStore{
		games: make(map[string]*models.GameState),
	}
}

func (s *MemoryGameStore) Create(game *models.GameState) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	s.games[game.ID] = game
	return nil
}

func (s *MemoryGameStore) Update(game *models.GameState) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	// Only update if it already exists.
	if _, ok := s.games[game.ID]; !ok {
		return ErrGameNotFound
	}
	s.games[game.ID] = game
	return nil
}

func (s *MemoryGameStore) Get(id string) (*models.GameState, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()

	game, ok := s.games[id]
	if !ok {
		return nil, ErrGameNotFound
	}
	return game, nil
}

// List of games within the GameStore
func (s *MemoryGameStore) List(filter GameFilter) ([]*models.GameState, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()

	var result []*models.GameState

	for _, g := range s.games {
		if filter.Mode != nil && g.Mode != *filter.Mode {
			continue
		}
		if filter.Status != nil && g.Status != *filter.Status {
			continue
		}
		result = append(result, g)
	}

	// Apply offset and limit.
	start := filter.Offset
	if start < 0 {
		start = 0
	}
	if start > len(result) {
		return []*models.GameState{}, nil
	}

	end := len(result)
	if filter.Limit > 0 && start+filter.Limit < end {
		end = start + filter.Limit
	}

	return result[start:end], nil
}
