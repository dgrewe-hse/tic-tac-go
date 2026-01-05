## Tic-Tac-Go Game Server 

<p align="center">
  <a href="LICENSE"><img alt="Educational" src="https://img.shields.io/badge/status-educational-blue"></a>
  <a href="docs/"><img alt="Docs" src="https://img.shields.io/badge/docs-available-brightgreen"></a>
  <a href="LICENSE"><img alt="License: Apache-2.0" src="https://img.shields.io/badge/license-Apache--2.0-blue"></a>
  <img alt="Made at HSE Esslingen" src="https://img.shields.io/badge/made%20at-HSE%20Esslingen-0a7ea4">
</p>

> IMPORTANT: This repository is for educational purposes only. It may contain unfinished, faulty, or even non-executable code used solely for teaching during the Internet Technologies course. Do not use this repository for any production system.

---

This repository contains a minimal Go backend for a Tic-Tac-Toe game server, designed to expose REST and WebSocket APIs for frontends such as Angular or Vue. This project is guided by the design document in `docs/spec.md`.

### Building and Running the server

#### Prerequisites

- Go 1.22 or higher installed
- Git (to clone the repository)

#### Setup

1. **Install dependencies:**
   ```bash
   go mod download
   ```

2. **Build the server (optional):**
   ```bash
   go build -o tic-tac-go-server ./cmd/server
   ```
   This creates an executable binary `tic-tac-go-server` that can be run directly.

3. **Run the server:**

   **Option A: Run directly with `go run` (recommended for development):**
   ```bash
   go run ./cmd/server
   ```

   **Option B: Run the compiled binary:**
   ```bash
   ./tic-tac-go-server
   ```

By default the server listens on port `8080`. You can override this via:

```bash
TICTACGO_PORT=9090 go run ./cmd/server
```

or for the compiled binary:

```bash
TICTACGO_PORT=9090 ./tic-tac-go-server
```

Once running, you can verify the basic health endpoint:

```bash
curl http://localhost:8080/health
```

You should receive:

```json
{"status":"ok"}
```

### HTTP API overview (for frontend developers)

The backend currently exposes the following REST endpoints (all paths are prefixed with `http://localhost:8080` by default):

- `POST /players`
  - Request body: `{"name": "Alice"}`
  - Response: `{"playerId": "...","name":"Alice"}`
  - Used to obtain a `playerId` that is then sent in the `X-Player-Id` header for all game-related calls.

- `POST /games`
  - Headers: `X-Player-Id: <playerId>`
  - Request body: `{"mode": "PVP"}` or `{"mode": "PVC"}`
  - Response: game state:
    - `gameId`, `mode`, `board` (`3x3` array of `"X" | "O" | ""`), `currentTurn`, `status`, `winner`.

- `GET /games`
  - Query parameters (optional):
    - `mode` = `PVP` or `PVC`
    - `status` = `WAITING_FOR_PLAYER` | `IN_PROGRESS` | `FINISHED`
    - `limit`, `offset` (pagination)
  - Response: `{ "games": [ { "gameId", "mode", "status", "createdAt", "createdBy": { "playerId", "name" } } ] }`
  - Typical frontend usage: list open PVP games with `GET /games?mode=PVP&status=WAITING_FOR_PLAYER`.

- `GET /games/{gameId}`
  - Response: same shape as `POST /games` response for that specific game.

- `POST /games/{gameId}/join`
  - Headers: `X-Player-Id: <playerId>`
  - Response: full game state after the player joined (PVP only; second player becomes `"O"`).

- `POST /games/{gameId}/moves`
  - Headers: `X-Player-Id: <playerId>`
  - Request body: `{"row": 0, "col": 2}`
  - Response: updated game state after the move (and, in PVC mode, after the AI response move if applicable).

For concrete example calls and typical flows (create player → create game → list games → join → make moves), see the shell scripts documented in `scripts/README.md`.

### WebSocket API overview (for frontend developers)

The backend also supports **WebSocket connections** for real-time game state updates. This is an **alternative to polling** the REST API.

#### Endpoint

- **`GET /ws/games/{gameId}`** (WebSocket upgrade)
  - URL: `ws://localhost:8080/ws/games/{gameId}` (or `wss://` for HTTPS)
  - **Behavior:**
    - Upgrades HTTP connection to WebSocket
    - Validates that the game exists (returns 404 if not found)
    - Immediately sends current game state to the newly connected client
    - Registers connection for that `gameId` to receive future broadcasts
    - Automatically broadcasts state updates when:
      - A player joins the game
      - A move is made (including AI moves in PVC mode)
      - Game status changes (win/draw)

#### Message Protocol

**Server → Client messages:**

1. **Game state update:**
   ```json
   {
     "type": "state",
     "payload": {
       "gameId": "uuid",
       "board": [["X", "", ""], ["", "O", ""], ["", "", ""]],
       "currentTurn": "O",
       "status": "IN_PROGRESS",
       "winner": ""
     }
   }
   ```

2. **Error message:**
   ```json
   {
     "type": "error",
     "payload": {
       "message": "Game not found"
     }
   }
   ```

**Client → Server:**

- Currently, **all actions are sent via REST API only**:
  - `POST /players` - Create player
  - `POST /games` - Create game
  - `POST /games/{gameId}/join` - Join game
  - `POST /games/{gameId}/moves` - Make move
- The WebSocket connection is **read-only** for receiving real-time updates
- **Recommended pattern:** Use REST for actions, WebSocket for receiving updates

#### Frontend Integration Example

**JavaScript/TypeScript:**
```javascript
const gameId = "your-game-id";
const ws = new WebSocket(`ws://localhost:8080/ws/games/${gameId}`);

ws.onopen = () => {
  console.log("WebSocket connected");
};

ws.onmessage = (event) => {
  const message = JSON.parse(event.data);
  if (message.type === "state") {
    // Update UI with new game state
    updateGameBoard(message.payload.board);
    updateCurrentTurn(message.payload.currentTurn);
    updateStatus(message.payload.status);
    if (message.payload.winner) {
      showGameOver(message.payload.winner);
    }
  } else if (message.type === "error") {
    console.error("WebSocket error:", message.payload.message);
  }
};

ws.onerror = (error) => {
  console.error("WebSocket error:", error);
};

ws.onclose = () => {
  console.log("WebSocket disconnected");
  // Optionally reconnect or fall back to polling
};
```

**Recommended pattern:**
- Use WebSocket for **real-time updates** (opponent moves, game state changes)
- Use REST API for **actions** (creating games, joining, making moves)
- Fall back to polling `GET /games/{gameId}` if WebSocket connection fails

#### Testing WebSocket

A Go-based test client is available (no external dependencies):
```bash
# Terminal 1: Connect via WebSocket
GAME_ID="..." go run scripts/07_websocket-test.go "$GAME_ID"

# Terminal 2: Make moves via REST
PLAYER_ID="..." GAME_ID="..." ROW=0 COL=0 ./scripts/06_make-move.sh

# Watch Terminal 1 for real-time updates
```

**Complete example flow:**
```bash
# 1. Create players and game via REST
PLAYER_ID_ALICE=$(./scripts/01_create-player.sh "Alice" | jq -r '.playerId')
GAME_ID=$(PLAYER_ID="$PLAYER_ID_ALICE" MODE=PVP ./scripts/02_create-game.sh | jq -r '.gameId')

# 2. Connect to WebSocket in one terminal
go run scripts/07_websocket-test.go "$GAME_ID"

# 3. In another terminal: join game and make moves via REST
PLAYER_ID_BOB=$(./scripts/01_create-player.sh "Bob" | jq -r '.playerId')
PLAYER_ID="$PLAYER_ID_BOB" GAME_ID="$GAME_ID" ./scripts/05_join-game.sh
PLAYER_ID="$PLAYER_ID_ALICE" GAME_ID="$GAME_ID" ROW=0 COL=0 ./scripts/06_make-move.sh

# 4. Watch the WebSocket terminal for real-time updates
```

### Docker Deployment

The server can be deployed using Docker for production environments.

**Quick start with Docker Compose:**
```bash
docker-compose up -d
```

**Manual Docker build and run:**
```bash
# Build the image
docker build -t tic-tac-go-server:latest .

# Run the container
docker run -d --name tic-tac-go-server -p 8080:8080 tic-tac-go-server:latest
```

For detailed deployment instructions, see [DEPLOYMENT.md](DEPLOYMENT.md).

### Development

If you prefer to develop or run the project inside a preconfigured container, this repository includes a Dev Container setup under the `.devcontainer` directory.

- Using VS Code with Dev Containers extension installed:
  - Open the project folder in your editor.
  - When prompted, choose "Reopen in Container" (or use the Command Palette: `Dev Containers: Reopen in Container`).
  - The container image defined in `.devcontainer/Dockerfile` will be built, dependencies will be prepared, and you will be attached to a shell with Go tools preinstalled.
  - From inside the container, you can run the server the same way:

```bash
go run ./cmd/server
```

## License
This project is licensed under the Apache License 2.0. See the [LICENSE](LICENSE) file for details.

