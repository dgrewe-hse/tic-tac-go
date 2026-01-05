## Tic-Tac-Go helper scripts

This directory contains small shell scripts that demonstrate how to interact with the Tic-Tac-Go HTTP API.
They are intended as runnable examples and quick references for frontend developers.

All scripts assume:

- The server is running at `http://localhost:8080` (override with `API_BASE` env var).
- `curl` and `jq` are installed.

Make all scripts executable once:

```bash
chmod +x scripts/*.sh
```

---

### 00_full_test_script.sh

- **Purpose**: Complete end-to-end test demonstrating a full game flow from start to finish.
- **Usage**:

```bash
./scripts/00_full_test_script.sh
```

**What it does:**
1. Creates two players: Alice and Bob
2. Alice creates a PVP game
3. Bob joins the game
4. Both players make moves in turn (Alice → Bob → Alice → Bob → Alice)
5. Shows the final game state

**Output:**
- Colored console output showing each step
- Final game state with board visualization
- Summary with all IDs (game, players)

**Example output:**
```
=== Tic-Tac-Go Full Test Script ===

Step 1: Creating player Alice...
✓ Alice created with ID: abc-123...

Step 2: Creating player Bob...
✓ Bob created with ID: def-456...

Step 3: Alice creates a PVP game...
✓ Game created with ID: game-789...

... (continues with moves) ...

=== Test Complete ===
Game ID: game-789...
Alice ID: abc-123...
Bob ID: def-456...
```

**Use cases:**
- Quick verification that the server is working correctly
- Demonstration of the complete API flow
- Integration testing
- Reference for frontend developers

---

### 01_create-player.sh

- **Purpose**: Create a new player and return its `playerId`.
- **Usage**:

```bash
scripts/01_create-player.sh "Alice"
```

- **Output**: JSON:

```json
{ "playerId": "...","name":"Alice" }
```

You typically capture `playerId` for use in later calls:

```bash
PLAYER_ID=$(scripts/01_create-player.sh "Alice" | jq -r '.playerId')
```

### 02_create-game.sh

- **Purpose**: Create a new game in `PVP` or `PVC` mode as a specific player.
- **Environment variables**:
  - `PLAYER_ID` (required): ID of the creator (must exist).
  - `MODE` (optional): `PVP` (default) or `PVC`.
- **Usage**:

```bash
PLAYER_ID="$PLAYER_ID" MODE=PVP scripts/02_create-game.sh
```

- **Output**: JSON game state including `gameId`:

```bash
GAME_ID=$(PLAYER_ID="$PLAYER_ID" MODE=PVP scripts/02_create-game.sh | jq -r '.gameId')
```

### 03_list-games.sh

- **Purpose**: List games, typically "open" PVP games that are waiting for a second player.
- **Environment variables**:
  - `MODE` (optional): game mode filter, defaults to `PVP`.
  - `STATUS` (optional): status filter, defaults to `WAITING_FOR_PLAYER`.
- **Usage**:

```bash
MODE=PVP STATUS=WAITING_FOR_PLAYER scripts/03_list-games.sh
```

- **Output**: JSON:

```json
{
  "games": [
    {
      "gameId": "...",
      "mode": "PVP",
      "status": "WAITING_FOR_PLAYER",
      "createdAt": "...",
      "createdBy": { "playerId": "...", "name": "Alice" }
    }
  ]
}
```

### 04_get-game.sh

- **Purpose**: Fetch the full state of a specific game.
- **Environment variables**:
  - `GAME_ID` (required): ID of the game to fetch.
- **Usage**:

```bash
GAME_ID="$GAME_ID" scripts/04_get-game.sh
```

### 05_join-game.sh

- **Purpose**: Join an existing PVP game as the second player.
- **Environment variables**:
  - `PLAYER_ID` (required): ID of the joining player.
  - `GAME_ID` (required): ID of the game to join.
- **Usage**:

```bash
PLAYER_ID_BOB=$(scripts/01_create-player.sh "Bob" | jq -r '.playerId')
PLAYER_ID="$PLAYER_ID_BOB" GAME_ID="$GAME_ID" scripts/05_join-game.sh
```

### 06_make-move.sh

- **Purpose**: Make a move in an existing game as one of the players.
- **Environment variables**:
  - `PLAYER_ID` (required): ID of the player making the move.
  - `GAME_ID` (required): ID of the game.
  - `ROW` (required): row index `0..2`.
  - `COL` (required): column index `0..2`.
- **Usage** (example flow):

```bash
# Create players
PLAYER_ID_ALICE=$(scripts/01_create-player.sh "Alice" | jq -r '.playerId')
PLAYER_ID_BOB=$(scripts/01_create-player.sh "Bob" | jq -r '.playerId')

# Create game as Alice
GAME_ID=$(PLAYER_ID="$PLAYER_ID_ALICE" MODE=PVP scripts/02_create-game.sh | jq -r '.gameId')

# Bob joins the game
PLAYER_ID="$PLAYER_ID_BOB" GAME_ID="$GAME_ID" scripts/05_join-game.sh

# Alice makes a move at (0,0)
PLAYER_ID="$PLAYER_ID_ALICE" GAME_ID="$GAME_ID" ROW=0 COL=0 scripts/06_make-move.sh
```

---

### 07_websocket-test.go

- **Purpose**: Complete end-to-end test demonstrating WebSocket real-time updates during a full game flow.
- **Requirements**: Go (already available in devcontainer)
- **Usage**:

```bash
go run scripts/07_websocket-test.go
```

**What it does:**
1. Creates two players (Alice and Bob) via REST API
2. Alice creates a PVP game via REST API
3. Connects to WebSocket for real-time game state updates
4. Bob joins the game via REST API (WebSocket receives update)
5. Both players make moves via REST API (WebSocket receives updates in real-time)
6. Shows all WebSocket messages with formatted board visualization

**Output:**
- Colored console output showing each step
- Real-time WebSocket messages as game state changes
- Formatted board visualization for each update
- Final summary with all IDs

**Example output:**
```
=== Tic-Tac-Go WebSocket Full Test ===

[Step 1] Creating player Alice...
✓ Alice created with ID: abc-123...

[Step 2] Creating player Bob...
✓ Bob created with ID: def-456...

[Step 3] Alice creates a PVP game...
✓ Game created with ID: game-789...

[Step 4] Connecting to WebSocket for real-time updates...
✓ WebSocket connected!

[WebSocket] Received initial game state:
    |   |  
  -----------
    |   |  
  -----------
    |   |  

[Step 5] Bob joins the game...
✓ Bob joined the game
[WebSocket] Received update: Bob joined!
  Status: IN_PROGRESS

[Step 6] Alice makes move at (0,0)...
✓ Move made
[WebSocket] Received update: Alice's move!
  X |   |  
  -----------
    |   |  
  -----------
    |   |  

... (continues with moves) ...

=== Test Complete ===
Game ID: game-789...
Alice ID: abc-123...
Bob ID: def-456...

All game state updates were received via WebSocket in real-time!
```

**Use cases:**
- Demonstrates WebSocket real-time update functionality
- Shows the recommended pattern: REST for actions, WebSocket for updates
- Integration testing for WebSocket implementation
- Reference for frontend developers on WebSocket usage

**Advantages:**
- No external dependencies (uses gorilla/websocket already in the project)
- Works immediately in the devcontainer
- Complete end-to-end test in a single command
- Shows both REST API usage and WebSocket message reception

---

## Notes for Frontend Developers

- **Player ID persistence:** Frontends should store the `playerId` (e.g., in localStorage) and send it in the `X-Player-Id` header for all game-related requests.
- **Real-time updates:** Use WebSocket (`GET /ws/games/{gameId}`) for real-time game state updates. See the main `README.md` for WebSocket API documentation.
- **Error handling:** Scripts use `set -euo pipefail` for strict error handling. In your frontend, handle HTTP error status codes (400, 403, 404, 500) appropriately.


