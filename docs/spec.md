## Tic-Tac-Go Server – Architecture & Implementation Plan

### 1. Goals & Scope

- **Purpose**: Provide a Go backend for a Tic-Tac-Toe game that can be consumed by web frontends (e.g. Angular, Vue).
- **Capabilities**:
  - **REST APIs** for all core operations.
  - **WebSocket** support for real-time game updates.
  - Support **Player vs Player (PVP)** and **Player vs Computer (PVC)**.
  - Allow players to **discover open PVP games** that are waiting for an opponent.
- **Constraints**:
  - Implementation time: a few hours → focus on simplicity and robustness.
  - Use **in-memory storage** (no external DB).
  - Prioritize **clarity and testability** over micro-optimizations.

---

### 2. High-Level Architecture

- **HTTP Layer**:
  - Handles routing, JSON encoding/decoding, and basic validation.
  - Exposes both **REST** endpoints and a **WebSocket** endpoint.
- **Service Layer**:
  - Implements application use-cases:
    - Creating/joining games.
    - Making moves.
    - Listing open games.
  - Orchestrates between game engine, AI, store, and WebSocket hub.
- **Game Engine Layer**:
  - Pure game rules: board representation, move validation, win/draw detection.
- **AI Layer**:
  - Simple heuristic opponent for PVC mode.
- **Store Layer**:
  - In-memory storage for players and games using Go maps and mutexes.
- **WebSocket Hub**:
  - Manages WebSocket connections per game.
  - Broadcasts game state updates to connected clients.

Project layout (proposed):

- `cmd/server/main.go`
- `internal/`
  - `http/`
    - `router.go`
    - `handlers.go`
    - `middleware.go`
  - `ws/`
    - `hub.go`
    - `connection.go`
  - `game/`
    - `board.go`
    - `rules.go`
  - `ai/`
    - `ai.go`
  - `service/`
    - `game_service.go`
    - `player_service.go`
  - `store/`
    - `memory_game_store.go`
    - `memory_player_store.go`
  - `models/`
    - `types.go`

---

### 3. Domain Model

#### 3.1 Core Types

In `internal/models/types.go`:

- **Player**
  - `ID string`
  - `Name string`

- **GameMode**
  - Enum-like string: `"PVP"`, `"PVC"`.

- **GameStatus**
  - `"WAITING_FOR_PLAYER"`
  - `"IN_PROGRESS"`
  - `"FINISHED"`

- **Symbol**
  - `"X"` or `"O"`.

- **Board**
  - `type Board [3][3]Symbol` (or `string` if simpler).

- **GameState**
  - `ID string`
  - `Mode GameMode`
  - `Board Board`
  - `PlayerXID string`
  - `PlayerOID string` (`"AI"` for computer in PVC)
  - `CurrentTurn Symbol`
  - `Status GameStatus`
  - `Winner string` (`"", "X", "O", "DRAW"`)
  - `CreatedAt time.Time`
  - `UpdatedAt time.Time`

- **GameSummary** (for listing/open games):
  - `ID string`
  - `Mode GameMode`
  - `Status GameStatus`
  - `CreatedAt time.Time`
  - `CreatedByPlayerID string`
  - `CreatedByPlayerName string`

---

### 4. API Design – REST Endpoints

#### 4.1 Identification / Lightweight Auth

- Clients first create a player:
  - `POST /players` with `{"name": "Alice"}`.
  - Response includes `playerId`.
- All subsequent calls include:
  - Header: `X-Player-Id: <playerId>`.
- No passwords, sessions, or JWTs (to keep implementation small).

#### 4.2 Player Endpoints

- **Create player**
  - **POST** `/players`
  - Request:
    ```json
    { "name": "Alice" }
    ```
  - Response:
    ```json
    { "playerId": "uuid", "name": "Alice" }
    ```

#### 4.3 Game Endpoints

- **Create game**
  - **POST** `/games`
  - Headers: `X-Player-Id`
  - Request:
    ```json
    { "mode": "PVP" } // or "PVC"
    ```
  - Behavior:
    - `PVP`:
      - Creator becomes `PlayerXID`.
      - `PlayerOID` empty initially.
      - `Status = "WAITING_FOR_PLAYER"`.
    - `PVC`:
      - Creator is `PlayerXID`, AI is `PlayerOID = "AI"`.
      - `Status = "IN_PROGRESS"`.
  - Response (common game representation):
    ```json
    {
      "gameId": "uuid",
      "mode": "PVP",
      "board": [["", "", ""], ["", "", ""], ["", "", ""]],
      "currentTurn": "X",
      "status": "WAITING_FOR_PLAYER",
      "winner": ""
    }
    ```

- **Join game (PVP)**
  - **POST** `/games/{gameId}/join`
  - Headers: `X-Player-Id`
  - Behavior:
    - Game must be `PVP` and `Status = "WAITING_FOR_PLAYER"`.
    - Assign joining player as `PlayerOID`.
    - Change status to `"IN_PROGRESS"`.
  - Response: full game representation.

- **Get game state**
  - **GET** `/games/{gameId}`
  - Headers: `X-Player-Id` (optional, but useful for authorization).
  - Response: full game representation.

- **Make move**
  - **POST** `/games/{gameId}/moves`
  - Headers: `X-Player-Id`
  - Request:
    ```json
    { "row": 0, "col": 2 }
    ```
  - Behavior:
    - Validate:
      - Game exists.
      - Game status is `"IN_PROGRESS"`.
      - Caller is one of the players.
      - It is the caller’s turn.
      - Move is within bounds and cell is empty.
    - Apply move via game engine.
    - Detect win/draw and update status and winner.
    - If PVC and the game is still in progress after player move:
      - Compute AI move.
      - Apply AI move and re-check status/winner.
    - Persist updated game state.
    - Broadcast new state to WebSocket clients for that game.
  - Response: updated game representation.

#### 4.4 List Open PVP Games

- **List games (with filters)**
  - **GET** `/games`
  - Query params:
    - `mode` (optional): `"PVP"` or `"PVC"`.
    - `status` (optional): `"WAITING_FOR_PLAYER"`, `"IN_PROGRESS"`, `"FINISHED"`.
    - `limit` (optional, default e.g. 20).
    - `offset` (optional, default 0).
  - Typical usage to find open PVP games:
    - `GET /games?mode=PVP&status=WAITING_FOR_PLAYER`
  - Response (list of `GameSummary` objects, not full games):
    ```json
    {
      "games": [
        {
          "gameId": "uuid-1",
          "mode": "PVP",
          "status": "WAITING_FOR_PLAYER",
          "createdAt": "2025-12-02T10:00:00Z",
          "createdBy": {
            "playerId": "player-123",
            "name": "Alice"
          }
        },
        {
          "gameId": "uuid-2",
          "mode": "PVP",
          "status": "WAITING_FOR_PLAYER",
          "createdAt": "2025-12-02T10:05:00Z",
          "createdBy": {
            "playerId": "player-456",
            "name": "Bob"
          }
        }
      ]
    }
    ```
  - Frontends:
    - Poll this endpoint or call on demand when user opens a “Join game” screen.
    - Once user selects a game, call `POST /games/{gameId}/join`.

---

### 5. WebSocket Design

#### 5.1 Endpoint

- **GET** `/ws/games/{gameId}`
  - Query or header for identification:
    - `X-Player-Id` header (preferred).
    - Or `?playerId=...` as a fallback.

#### 5.2 Connection Behavior

- On connect:
  - Validate `gameId` exists.
  - Optionally validate `playerId` is a participant or allow spectators.
  - Register connection with the hub for that `gameId`.
  - Immediately send current game state to the newly connected client.

#### 5.3 Message Protocol

- **Server → Client** messages:
  - Game state:
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
  - Error:
    ```json
    {
      "type": "error",
      "payload": { "message": "It's not your turn." }
    }
    ```

- **Client → Server**:
  - For the initial version, moves are sent via REST only:
    - Clients call `POST /games/{gameId}/moves`.
    - Server broadcasts updated state via WebSocket.
  - (Optional future extension) Support `type: "move"` messages to submit moves via WS.

---

### 6. Game Engine (Pure Logic)

Package `internal/game`:

- **Types**:
  - `type Symbol string`
  - `type Board [3][3]Symbol`

- **Functions**:
  - `func NewBoard() Board`
  - `func IsValidMove(board Board, row, col int) bool`
  - `func ApplyMove(board Board, row, col int, symbol Symbol) (Board, error)`
  - `func CheckWinner(board Board) (winner Symbol, isDraw bool)`
  - `func AvailableMoves(board Board) [][2]int`
  - `func IsFull(board Board) bool`

Responsibilities:

- Encapsulate all game rules and board manipulations.
- Be deterministic and side-effect free to simplify testing.

---

### 7. AI Design

Package `internal/ai`:

- **Goal**: Simple but non-trivial AI for PVC mode.
- **Public API**:
  - `func ChooseMove(board game.Board, aiSymbol, opponentSymbol game.Symbol) (row, col int)`

- **Heuristic**:
  1. If AI has a winning move this turn, take it.
  2. Else if opponent has a winning move next turn, block it.
  3. Else if center is free, take center.
  4. Else pick a random available corner or side.

This keeps implementation simple but provides reasonable challenge.

---

### 8. Store Layer

Package `internal/store`:

- **Interfaces**:
  - `GameStore`:
    ```go
    type GameStore interface {
        Create(game *models.GameState) error
        Update(game *models.GameState) error
        Get(id string) (*models.GameState, error)
        List(filter GameFilter) ([]*models.GameState, error)
    }
    ```
  - `PlayerStore`:
    ```go
    type PlayerStore interface {
        Create(player *models.Player) error
        Get(id string) (*models.Player, error)
    }
    ```
  - `GameFilter`:
    ```go
    type GameFilter struct {
        Mode   *models.GameMode
        Status *models.GameStatus
        Limit  int
        Offset int
    }
    ```

- **In-memory Implementation**:
  - `MemoryGameStore` with:
    - `games map[string]*models.GameState`
    - `mu sync.RWMutex`
  - `MemoryPlayerStore` with:
    - `players map[string]*models.Player`
    - `mu sync.RWMutex`
  - `List` operation filters in-memory by mode/status and applies limit/offset.

---

### 9. Service Layer

Package `internal/service`:

#### 9.1 GameService

Interface:

```go
type GameService interface {
    CreateGame(ctx context.Context, creatorPlayerID string, mode models.GameMode) (*models.GameState, error)
    JoinGame(ctx context.Context, gameID, playerID string) (*models.GameState, error)
    GetGame(ctx context.Context, gameID string) (*models.GameState, error)
    MakeMove(ctx context.Context, gameID, playerID string, row, col int) (*models.GameState, error)
    ListGames(ctx context.Context, filter store.GameFilter) ([]*models.GameSummary, error)
}
```

Implementation dependencies:

- `GameStore`
- `PlayerStore`
- `AI` module (`ai.ChooseMove`)
- `WebSocketHub` (via an interface like `Broadcaster` with `BroadcastGameState`).

Key responsibilities:

- Enforce business rules on top of game engine:
  - Turn order, participant checks.
  - Game status transitions.
  - Winner/draw state updates.
- Trigger AI moves in PVC mode.
- Notify WebSocket hub of state changes.
- For `ListGames`, call `GameStore.List`, map `GameState` to `GameSummary`, join creator’s name from `PlayerStore`.

##### 9.1.1 `MakeMove` Logic (Pseudocode)

```text
function MakeMove(gameID, playerID, row, col):
  game = gameStore.Get(gameID)
  if game is nil:
      return error NotFound

  if game.Status != IN_PROGRESS:
      return error InvalidGameState

  // Determine symbol for this player
  if playerID == game.PlayerXID:
      symbol = "X"
  else if playerID == game.PlayerOID:
      symbol = "O"
  else:
      return error NotParticipant

  if game.CurrentTurn != symbol:
      return error NotPlayersTurn

  if !IsValidMove(game.Board, row, col):
      return error InvalidMove

  // Apply player move
  game.Board = ApplyMove(game.Board, row, col, symbol)
  winner, isDraw = CheckWinner(game.Board)

  if winner != "":
      game.Status = FINISHED
      game.Winner = winner
  else if isDraw:
      game.Status = FINISHED
      game.Winner = "DRAW"
  else:
      // Switch turn
      game.CurrentTurn = opposite(symbol)

      // AI turn in PVC
      if game.Mode == PVC and game.Status == IN_PROGRESS and game.CurrentTurn == aiSymbol:
          aiRow, aiCol = ChooseMove(game.Board, aiSymbol, playerSymbol)
          game.Board = ApplyMove(game.Board, aiRow, aiCol, aiSymbol)
          winner, isDraw = CheckWinner(game.Board)
          if winner != "":
              game.Status = FINISHED
              game.Winner = winner
          else if isDraw:
              game.Status = FINISHED
              game.Winner = "DRAW"
          else:
              game.CurrentTurn = playerSymbol

  game.UpdatedAt = now()
  gameStore.Update(game)
  hub.BroadcastGameState(gameID, game)
  return game
```

#### 9.2 PlayerService

Simple wrapper over `PlayerStore`:

- `CreatePlayer(ctx, name string) (*Player, error)`
- `GetPlayer(ctx, id string) (*Player, error)`

---

### 10. WebSocket Hub

Package `internal/ws`:

- **Connection struct**:
  - Wraps `*websocket.Conn`.
  - Has `send chan []byte` for outgoing messages.

- **Hub interface**:
  ```go
  type Hub interface {
      Register(gameID string, conn *Connection)
      Unregister(gameID string, conn *Connection)
      BroadcastGameState(gameID string, state *models.GameState)
  }
  ```

- **Implementation**:
  - Maintains:
    - `clients map[string]map[*Connection]struct{}` (keyed by `gameID`).
    - Mutex to protect `clients`.
  - `BroadcastGameState`:
    - Serializes state into JSON `{"type": "state", "payload": ...}`.
    - Sends to all connections registered for that `gameID`.

- **HTTP WS Handler**:
  - Upgrades HTTP request to WebSocket.
  - Extracts `gameId` from path and `playerId` from header/query.
  - Validates `gameId` exists.
  - Registers connection in hub.
  - Starts read/write goroutines; on close, unregisters.

---

### 11. HTTP Layer Implementation

Package `internal/http`:

- **Router** (using `chi` or similar):
  - `POST /players`
  - `POST /games`
  - `GET /games` (list with filters; used to list open PVP games)
  - `POST /games/{gameId}/join`
  - `GET /games/{gameId}`
  - `POST /games/{gameId}/moves`
  - `GET /ws/games/{gameId}` (WebSocket upgrade; handler can live in `ws` package but mounted here).

- **Middleware**:
  - Extract `X-Player-Id` into context for handlers that require a logged-in player.

- **Handlers**:
  - Decode JSON request body.
  - Validate input (basic checks).
  - Call the appropriate service method.
  - Map domain errors to HTTP status codes:
    - 400: invalid arguments/moves.
    - 403: player not participant.
    - 404: game or player not found.
    - 500: unexpected internal error.

---

### 12. Testing Strategy

#### 12.1 High Priority

- **Game engine tests** (`internal/game`):
  - `CheckWinner`:
    - All winning combinations (rows, columns, diagonals) for X and O.
    - Draw detection on a full board with no winner.
  - `IsValidMove`:
    - Valid positions.
    - Out-of-bounds indexes.
    - Already occupied cells.
  - `ApplyMove`:
    - Correctly updates board state.

- **AI tests** (`internal/ai`):
  - Picks winning move if available.
  - Blocks imminent opponent win.
  - Always returns a valid free cell.

- **Service tests** (`internal/service`):
  - `CreateGame`:
    - Correct setup for PVP and PVC.
  - `JoinGame`:
    - Transitions from `WAITING_FOR_PLAYER` to `IN_PROGRESS`.
    - Error cases (wrong mode/status).
  - `MakeMove`:
    - Turn enforcement.
    - Game completion (win/draw).
    - AI response in PVC.
  - `ListGames`:
    - Filters by mode and status.
    - Supports pagination fields (limit/offset) in-memory.

- **HTTP handler tests** (`internal/http`):
  - Use `httptest` to:
    - Create player.
    - Create game.
    - Join game.
    - Make valid/invalid moves.
    - List open PVP games with `GET /games?mode=PVP&status=WAITING_FOR_PLAYER`.

#### 12.2 Nice-to-have

- WebSocket hub tests using mocked connections.
- Run `go test -race ./...` locally to catch race conditions.

---

### 13. Implementation Plan (Step-by-Step)

1. **Bootstrap project**
   - `go mod init github.com/yourname/tic-tac-go`
   - Create directory structure under `cmd/` and `internal/`.

2. **Models & Game Engine**
   - Implement types in `internal/models/types.go`.
   - Implement `internal/game` with board and rule functions.
   - Add unit tests for `game`.

3. **AI**
   - Implement `internal/ai.ChooseMove` with heuristic logic.
   - Add unit tests for AI behavior.

4. **Stores**
   - Implement `MemoryGameStore` and `MemoryPlayerStore`.
   - Implement `List` with filtering on mode/status for open games.

5. **Service Layer**
   - Implement `GameService` and `PlayerService`.
   - Integrate game engine, AI, and stores.
   - Add service tests for create/join/move/list flows.

6. **HTTP Layer**
   - Implement router, handlers, and middleware.
   - Add HTTP tests for key endpoints, especially `GET /games` for open PVP games.
   - Implement `cmd/server/main.go` to wire everything and start an HTTP server on a configurable port.

7. **WebSocket Hub**
   - Implement hub and connection management in `internal/ws`.
   - Implement `/ws/games/{gameId}` handler.
   - Integrate hub with `GameService` to broadcast state on each update.

8. **Manual Testing & Polish**
   - Run `go test ./...`.
   - Start server and manually test via `curl` or Postman:
     - Create players and games.
     - Join PVP games.
     - List open PVP games (`GET /games?mode=PVP&status=WAITING_FOR_PLAYER`).
     - Make moves in PVP and PVC modes.
   - Connect via WebSocket client to verify real-time updates.

This specification provides the detailed architecture and implementation steps necessary to implement the Tic-Tac-Go backend within a few focused hours, while keeping the system clean, testable, and ready for frontends such as Angular or Vue.

