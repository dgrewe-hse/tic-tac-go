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

	"tic-tac-go/internal/models"
	"tic-tac-go/internal/store"

	"github.com/google/uuid"
)

// playerService is a concrete implementation of PlayerService.
type playerService struct {
	playerStore store.PlayerStore
}

// NewPlayerService constructs a PlayerService using the given PlayerStore.
func NewPlayerService(playerStore store.PlayerStore) PlayerService {
	return &playerService{
		playerStore: playerStore,
	}
}

func (s *playerService) CreatePlayer(ctx context.Context, name string) (*models.Player, error) {
	player := &models.Player{
		ID:   uuid.NewString(),
		Name: name,
	}
	if err := s.playerStore.Create(player); err != nil {
		return nil, err
	}
	return player, nil
}

func (s *playerService) GetPlayer(ctx context.Context, id string) (*models.Player, error) {
	return s.playerStore.Get(id)
}
