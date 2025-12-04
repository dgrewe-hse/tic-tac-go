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


