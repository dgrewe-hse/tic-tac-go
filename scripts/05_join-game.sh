#!/bin/bash
# Copyright 2025 Esslingen University of Applied Sciences
#
# Licensed under the Apache License, Version 2.0 (the "License");
# you may not use this file except in compliance with the License.
# You may obtain a copy of the License at
#
#     http://www.apache.org/licenses/LICENSE-2.0
#
# Unless required by applicable law or agreed to in writing, software
# distributed under the License is distributed on an "AS IS" BASIS,
# WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
# See the License for the specific language governing permissions and
# limitations under the License.
#
# Author: Dennis Grewe
# Version: 1.0.0
# Date: 2025-12-03
#
# Usage: PLAYER_ID=<player-id> GAME_ID=<game-id> ./scripts/05_join-game.sh
# NOTE: You need to create a player first (here: Alice), then use the player ID to create a game, afterwards you can get the game using its ID. Next, you need to create a second player (here: Bob), then use the player ID to join the game.
# PLAYER_ID_ALICE=$(./scripts/01_create-player.sh "Alice" | jq -r '.playerId')
# GAME_ID=$(PLAYER_ID="$PLAYER_ID_ALICE" MODE=PVP ./scripts/02_create-game.sh | jq -r '.gameId')
# GAME_ID="$GAME_ID" ./scripts/04_get-game.sh
# PLAYER_ID_BOB=$(./scripts/01_create-player.sh "Bob" | jq -r '.playerId')
# PLAYER_ID="$PLAYER_ID_BOB" GAME_ID="$GAME_ID" ./scripts/05_join-game.sh

set -euo pipefail

# Extract ENV variables from the command line
# Default to localhost:8080 if API_BASE is not set
API_BASE="${API_BASE:-http://localhost:8080}"
# Extract PLAYER_ID from the command line (mandatory)
PLAYER_ID="${PLAYER_ID:-}"
# Extract GAME_ID from the command line (mandatory)
GAME_ID="${GAME_ID:-}"

# Check if PLAYER_ID and GAME_ID are set
if [[ -z "${PLAYER_ID}" || -z "${GAME_ID}" ]]; then
  echo "Usage: PLAYER_ID=<player-id> GAME_ID=<game-id> $0"
  exit 1
fi

# Join the game and print the response
curl -sS -X POST "${API_BASE}/games/${GAME_ID}/join" \
  -H "X-Player-Id: ${PLAYER_ID}" \
  -H 'Content-Type: application/json' \
  | jq .