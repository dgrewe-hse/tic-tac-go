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
# Usage: ./scripts/00_full_test_script.sh
#
# This script demonstrates a complete game flow:
# 1. Creates two players (Alice and Bob)
# 2. Alice creates a PVP game
# 3. Bob joins the game
# 4. Both players make moves in turn
# 5. Shows the final game state

set -euo pipefail

# Default to localhost:8080 if API_BASE is not set
API_BASE="${API_BASE:-http://localhost:8080}"

# Colors for output
GREEN='\033[0;32m'
BLUE='\033[0;34m'
YELLOW='\033[1;33m'
NC='\033[0m' # No Color

echo -e "${BLUE}=== Tic-Tac-Go Full Test Script ===${NC}\n"

# Step 1: Create Alice
echo -e "${YELLOW}Step 1: Creating player Alice...${NC}"
PLAYER_ID_ALICE=$(./01_create-player.sh "Alice" | jq -r '.playerId')
echo -e "${GREEN}✓ Alice created with ID: ${PLAYER_ID_ALICE}${NC}\n"

# Step 2: Create Bob
echo -e "${YELLOW}Step 2: Creating player Bob...${NC}"
PLAYER_ID_BOB=$(./01_create-player.sh "Bob" | jq -r '.playerId')
echo -e "${GREEN}✓ Bob created with ID: ${PLAYER_ID_BOB}${NC}\n"

# Step 3: Alice creates a PVP game
echo -e "${YELLOW}Step 3: Alice creates a PVP game...${NC}"
GAME_ID=$(PLAYER_ID="$PLAYER_ID_ALICE" MODE=PVP ./02_create-game.sh | jq -r '.gameId')
echo -e "${GREEN}✓ Game created with ID: ${GAME_ID}${NC}\n"

# Step 4: Show initial game state
echo -e "${YELLOW}Step 4: Initial game state:${NC}"
GAME_ID="$GAME_ID" ./04_get-game.sh | jq '{gameId, mode, status, currentTurn, board}'
echo ""

# Step 5: Bob joins the game
echo -e "${YELLOW}Step 5: Bob joins the game...${NC}"
PLAYER_ID="$PLAYER_ID_BOB" GAME_ID="$GAME_ID" ./05_join-game.sh | jq '{gameId, status, currentTurn, board}' > /dev/null
echo -e "${GREEN}✓ Bob joined the game${NC}\n"

# Step 6: Alice makes first move (top-left: 0,0)
echo -e "${YELLOW}Step 6: Alice makes move at (0,0)...${NC}"
PLAYER_ID="$PLAYER_ID_ALICE" GAME_ID="$GAME_ID" ROW=0 COL=0 ./06_make-move.sh | jq '{gameId, status, currentTurn, board, winner}' > /dev/null
echo -e "${GREEN}✓ Move made${NC}\n"

# Step 7: Bob makes move (top-middle: 0,1)
echo -e "${YELLOW}Step 7: Bob makes move at (0,1)...${NC}"
PLAYER_ID="$PLAYER_ID_BOB" GAME_ID="$GAME_ID" ROW=0 COL=1 ./06_make-move.sh | jq '{gameId, status, currentTurn, board, winner}' > /dev/null
echo -e "${GREEN}✓ Move made${NC}\n"

# Step 8: Alice makes move (center: 1,1)
echo -e "${YELLOW}Step 8: Alice makes move at (1,1)...${NC}"
PLAYER_ID="$PLAYER_ID_ALICE" GAME_ID="$GAME_ID" ROW=1 COL=1 ./06_make-move.sh | jq '{gameId, status, currentTurn, board, winner}' > /dev/null
echo -e "${GREEN}✓ Move made${NC}\n"

# Step 9: Bob makes move (bottom-right: 2,2)
echo -e "${YELLOW}Step 9: Bob makes move at (2,2)...${NC}"
PLAYER_ID="$PLAYER_ID_BOB" GAME_ID="$GAME_ID" ROW=2 COL=2 ./06_make-move.sh | jq '{gameId, status, currentTurn, board, winner}' > /dev/null
echo -e "${GREEN}✓ Move made${NC}\n"

# Step 10: Alice makes move (top-right: 0,2) - this should win for Alice (X)
echo -e "${YELLOW}Step 10: Alice makes move at (0,2)...${NC}"
FINAL_STATE=$(PLAYER_ID="$PLAYER_ID_ALICE" GAME_ID="$GAME_ID" ROW=0 COL=2 ./06_make-move.sh)
echo -e "${GREEN}✓ Move made${NC}\n"

# Step 11: Show final game state
echo -e "${YELLOW}Step 11: Final game state:${NC}"
echo "$FINAL_STATE" | jq '{
  gameId,
  mode,
  status,
  currentTurn,
  winner,
  board: .board | map(map(if . == "" then " " else . end))
}'

echo ""
echo -e "${BLUE}=== Test Complete ===${NC}"
echo -e "Game ID: ${GREEN}${GAME_ID}${NC}"
echo -e "Alice ID: ${GREEN}${PLAYER_ID_ALICE}${NC}"
echo -e "Bob ID: ${GREEN}${PLAYER_ID_BOB}${NC}"

