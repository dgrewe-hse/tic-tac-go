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
# Usage: PLAYER_ID=<player-id> MODE=[PVP|PVC] ./scripts/02_create-game.sh
# NOTE: You need to create a player first (here: Alice), then use the player ID to create a game.
# PLAYER_ID=$(./scripts/01_create-player.sh "Alice" | jq -r '.playerId')
# PLAYER_ID="$PLAYER_ID" MODE=PVP ./scripts/02_create-game.sh

set -euo pipefail

# Extract ENV variables from the command line
# Default to localhost:8080 if API_BASE is not set
API_BASE="${API_BASE:-http://localhost:8080}"
# Extract PLAYER_ID from the command line (mandatory)
PLAYER_ID="${PLAYER_ID:-}"
# Extract MODE from the command line, else PVP
MODE="${MODE:-PVP}"


# Check if PLAYER_ID is set
if [[ -z "${PLAYER_ID}" ]]; then
  echo "Usage: PLAYER_ID=<id> MODE=[PVP|PVC] $0"
  exit 1
fi

# Create a new game and print the response
curl -sS -X POST "${API_BASE}/games" \
  -H "X-Player-Id: ${PLAYER_ID}" \
  -H 'Content-Type: application/json' \
  -d "{\"mode\":\"${MODE}\"}" \
  | jq .