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
# Usage: ./scripts/03_list-games.sh [MODE=PVP] [STATUS=WAITING_FOR_PLAYER]
# NOTE: You need to create a player first (here: Alice), then use the player ID to create a game, afterwards you can list the games.
# PLAYER_ID=$(./scripts/01_create-player.sh "Alice" | jq -r '.playerId')
# PLAYER_ID="$PLAYER_ID" MODE=PVP ./scripts/02_create-game.sh
# MODE=PVP STATUS=WAITING_FOR_PLAYER ./scripts/03_list-games.sh

set -euo pipefail

# Extract ENV variables from the command line
# Default to localhost:8080 if API_BASE is not set
API_BASE="${API_BASE:-http://localhost:8080}"
# Extract MODE from the command line, else PVP
MODE="${MODE:-PVP}"
# Extract STATUS from the command line, else WAITING_FOR_PLAYER
STATUS="${STATUS:-WAITING_FOR_PLAYER}"

# List all games and print the response
curl -sS "${API_BASE}/games?mode=${MODE}&status=${STATUS}" | jq .
