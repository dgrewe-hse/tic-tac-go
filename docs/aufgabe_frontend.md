# Aufgabenstellung: Tic-Tac-Toe Frontend-Implementierung

## Übersicht

In dieser Aufgabe implementieren Sie ein **Frontend für ein Tic-Tac-Toe-Spiel** mit **Vue.js** oder **React.js**. Der Backend-Server ist bereits implementiert und stellt REST- und WebSocket-APIs zur Verfügung. Ihre Aufgabe ist es, eine benutzerfreundliche Web-Anwendung zu entwickeln, die mit diesem Server kommuniziert.

### Wichtiger Hinweis: Lokale Entwicklungsumgebung

**Alle Komponenten laufen lokal auf Ihrem Rechner:**

- ✅ **Kein gehosteter Server erforderlich** – Sie müssen sich nicht mit einem externen Server verbinden
- ✅ **Backend-Server** läuft lokal auf `http://localhost:8080`
- ✅ **Frontend-Anwendung** läuft lokal (z.B. auf `http://localhost:5173` für Vite/Vue oder `http://localhost:3000` für React)
- ✅ **Vollständige Kontrolle** über die Entwicklungsumgebung

### Quick-Start Checkliste

**Voraussetzungen prüfen:**
- [ ] Go 1.22+ installiert (`go version`)
- [ ] Node.js und npm installiert (`node --version`, `npm --version`)
- [ ] Backend-Server-Repository verfügbar/kopiert

**Setup-Schritte:**

1. [ ] **Backend-Server klonen/erhalten** (wird vom Kurs bereitgestellt)
2. [ ] **Backend-Server starten** (siehe Abschnitt "Server-Beschreibung" unten)
   ```bash
   cd <server-verzeichnis>
   go run ./cmd/server
   ```
3. [ ] **Server-Verfügbarkeit testen:**
   ```bash
   curl http://localhost:8080/health
   # Erwartet: {"status":"ok"}
   ```
4. [ ] **Frontend-Projekt erstellen** (Vue.js oder React.js)
5. [ ] **Frontend mit lokalem Server verbinden** (API-Base-URL: `http://localhost:8080`)

**Wichtig:** 
- Der Backend-Server **muss laufen**, bevor Sie das Frontend starten
- Das Frontend kommuniziert **ausschließlich** mit dem lokalen Server auf Port 8080
- Beide laufen in **separaten Terminal-Fenstern** parallel

---

## Lernziele

Nach Abschluss dieser Aufgabe sollten Sie:

1. **REST-APIs verstehen und nutzen** können
   - HTTP-Requests (GET, POST) mit Headers und Body-Daten senden
   - JSON-Daten verarbeiten und anzeigen
   - Fehlerbehandlung bei API-Aufrufen implementieren

2. **WebSocket-Verbindungen** für Echtzeit-Updates implementieren können
   - WebSocket-Verbindungen aufbauen und verwalten
   - Echtzeit-Nachrichten empfangen und verarbeiten
   - Verbindungsfehler behandeln und Reconnect-Logik implementieren

3. **State Management** in modernen Frontend-Frameworks beherrschen
   - Komponenten-basierte Architektur nutzen
   - Lokalen State und Props verwenden
   - Event-Handling implementieren

4. **Benutzerfreundliche UI/UX** gestalten
   - Intuitive Spieloberfläche entwickeln
   - Feedback für Benutzeraktionen bereitstellen
   - Loading-States und Fehlermeldungen anzeigen

5. **Moderne Frontend-Entwicklungspraktiken** anwenden
   - Komponenten strukturiert aufbauen
   - Code wiederverwendbar gestalten
   - Responsive Design berücksichtigen

---

## Server-Beschreibung

### Server-Repository

Der Backend-Server wird als **separates Repository** bereitgestellt. Sie müssen dieses Repository klonen oder herunterladen, um den Server lokal ausführen zu können.

**Voraussetzungen:**
- Go 1.22 oder höher installiert
- Git (zum Klonen des Repositories)

**Server-Repository klonen:**
```bash
git clone <server-repository-url>
cd tic-tac-go  # oder wie das Repository heißt
```

### Server starten

Der Backend-Server kann wie folgt gestartet werden:

```bash
# Im Server-Projektverzeichnis
go run ./cmd/server
```

Der Server läuft standardmäßig auf **Port 8080** (`http://localhost:8080`).

**Wichtig:** Lassen Sie den Server während der gesamten Frontend-Entwicklung laufen. Sie können den Server in einem separaten Terminal-Fenster starten.

Sie können den Port über eine Umgebungsvariable ändern:

```bash
TICTACGO_PORT=9090 go run ./cmd/server
```

**Hinweis:** Wenn Sie den Port ändern, müssen Sie auch die API-Base-URL in Ihrem Frontend entsprechend anpassen.

### Server-Verfügbarkeit prüfen

Nach dem Start können Sie testen, ob der Server läuft:

```bash
curl http://localhost:8080/health
```

Erwartete Antwort:
```json
{"status":"ok"}
```

Falls Sie diese Antwort erhalten, ist der Server erfolgreich gestartet und bereit für Frontend-Verbindungen.

### Server-Features

Der Server unterstützt:

- **REST API** für alle Spielaktionen (Spieler erstellen, Spiele erstellen/beitreten, Züge machen)
- **WebSocket API** für Echtzeit-Spielzustands-Updates
- **Zwei Spielmodi**:
  - **PVP (Player vs Player)**: Zwei menschliche Spieler spielen gegeneinander
  - **PVC (Player vs Computer)**: Ein Spieler spielt gegen eine KI
- **Spiel-Entdeckung**: Offene PVP-Spiele können aufgelistet werden, um beizutreten

---

## REST API Dokumentation

Alle Endpunkte sind unter `http://localhost:8080` erreichbar (oder dem konfigurierten Port).

### 1. Spieler erstellen

**Endpoint:** `POST /players`

**Request Body:**
```json
{
  "name": "Alice"
}
```

**Response:**
```json
{
  "playerId": "9772a11d-27ae-4952-a6a6-dad7b2802e5e",
  "name": "Alice"
}
```

**Verwendung:**
- Erstellt einen neuen Spieler
- Der `playerId` muss für alle weiteren API-Aufrufe im Header `X-Player-Id` mitgesendet werden
- **Wichtig:** Speichern Sie die `playerId` im `localStorage`, damit der Spieler bei erneutem Besuch wiedererkannt wird

**Beispiel (JavaScript):**
```javascript
const response = await fetch('http://localhost:8080/players', {
  method: 'POST',
  headers: { 'Content-Type': 'application/json' },
  body: JSON.stringify({ name: 'Alice' })
});
const player = await response.json();
localStorage.setItem('playerId', player.playerId);
```

---

### 2. Spiel erstellen

**Endpoint:** `POST /games`

**Headers:**
- `X-Player-Id: <playerId>` (erforderlich)
- `Content-Type: application/json`

**Request Body:**
```json
{
  "mode": "PVP"
}
```
oder
```json
{
  "mode": "PVC"
}
```

**Response:**
```json
{
  "gameId": "e2c58d69-5ee4-4779-b43f-7825494723c4",
  "mode": "PVP",
  "board": [["", "", ""], ["", "", ""], ["", "", ""]],
  "currentTurn": "X",
  "status": "WAITING_FOR_PLAYER",
  "winner": ""
}
```

**Verwendung:**
- Erstellt ein neues Spiel
- Im **PVP-Modus** ist der Status zunächst `WAITING_FOR_PLAYER` (wartet auf zweiten Spieler)
- Im **PVC-Modus** beginnt das Spiel sofort mit Status `IN_PROGRESS`
- Der Ersteller ist immer Spieler **X**

**Status-Werte:**
- `WAITING_FOR_PLAYER`: Spiel wartet auf zweiten Spieler (nur PVP)
- `IN_PROGRESS`: Spiel läuft
- `FINISHED`: Spiel beendet (Gewinner oder Unentschieden)

---

### 3. Spiele auflisten

**Endpoint:** `GET /games`

**Query-Parameter (optional):**
- `mode`: `PVP` oder `PVC`
- `status`: `WAITING_FOR_PLAYER`, `IN_PROGRESS`, oder `FINISHED`
- `limit`: Anzahl der Ergebnisse (Pagination)
- `offset`: Offset für Pagination

**Beispiel-Request:**
```
GET /games?mode=PVP&status=WAITING_FOR_PLAYER
```

**Response:**
```json
{
  "games": [
    {
      "gameId": "e2c58d69-5ee4-4779-b43f-7825494723c4",
      "mode": "PVP",
      "status": "WAITING_FOR_PLAYER",
      "createdAt": "2025-01-05T10:30:00Z",
      "createdBy": {
        "playerId": "9772a11d-27ae-4952-a6a6-dad7b2802e5e",
        "name": "Alice"
      }
    }
  ]
}
```

**Verwendung:**
- Zeigt verfügbare Spiele an
- Typischerweise für PVP-Spiele, die auf einen zweiten Spieler warten
- Ermöglicht es Spielern, einem bestehenden Spiel beizutreten

---

### 4. Spiel abrufen

**Endpoint:** `GET /games/{gameId}`

**Response:**
```json
{
  "gameId": "e2c58d69-5ee4-4779-b43f-7825494723c4",
  "mode": "PVP",
  "board": [["X", "O", ""], ["", "X", ""], ["", "", "O"]],
  "currentTurn": "X",
  "status": "IN_PROGRESS",
  "winner": ""
}
```

**Verwendung:**
- Ruft den aktuellen Zustand eines Spiels ab
- Kann für Polling verwendet werden (falls WebSocket nicht verfügbar ist)

---

### 5. Spiel beitreten

**Endpoint:** `POST /games/{gameId}/join`

**Headers:**
- `X-Player-Id: <playerId>` (erforderlich)

**Response:**
```json
{
  "gameId": "e2c58d69-5ee4-4779-b43f-7825494723c4",
  "mode": "PVP",
  "board": [["", "", ""], ["", "", ""], ["", "", ""]],
  "currentTurn": "X",
  "status": "IN_PROGRESS",
  "winner": ""
}
```

**Verwendung:**
- Nur für **PVP-Spiele** mit Status `WAITING_FOR_PLAYER`
- Der beitretende Spieler wird automatisch **Spieler O**
- Nach dem Beitreten ändert sich der Status zu `IN_PROGRESS`

**Fehler:**
- `400 Bad Request`: Spiel ist nicht im richtigen Status oder bereits voll
- `404 Not Found`: Spiel existiert nicht

---

### 6. Zug machen

**Endpoint:** `POST /games/{gameId}/moves`

**Headers:**
- `X-Player-Id: <playerId>` (erforderlich)
- `Content-Type: application/json`

**Request Body:**
```json
{
  "row": 0,
  "col": 2
}
```

**Response:**
```json
{
  "gameId": "e2c58d69-5ee4-4779-b43f-7825494723c4",
  "mode": "PVP",
  "board": [["X", "O", "X"], ["", "X", ""], ["", "", "O"]],
  "currentTurn": "O",
  "status": "IN_PROGRESS",
  "winner": ""
}
```

**Verwendung:**
- Macht einen Zug im Spiel
- `row` und `col` sind Indizes von 0-2
- Im **PVC-Modus** macht die KI automatisch ihren Zug nach dem Spielerzug
- Die Response enthält den aktualisierten Spielzustand (inklusive KI-Zug bei PVC)

**Fehler:**
- `400 Bad Request`: Ungültiger Zug (Zelle bereits belegt, außerhalb des Spielfelds, nicht am Zug)
- `403 Forbidden`: Spieler ist nicht Teilnehmer des Spiels
- `404 Not Found`: Spiel existiert nicht

**Gewinner/Unentschieden:**
Wenn das Spiel endet, enthält die Response:
```json
{
  "status": "FINISHED",
  "winner": "X"  // oder "O" oder "DRAW"
}
```

---

## WebSocket API Dokumentation

### Endpoint

**URL:** `ws://localhost:8080/ws/games/{gameId}`

(oder `wss://` für HTTPS in Produktion)

### Verhalten

- Die WebSocket-Verbindung wird über einen HTTP GET-Request mit Upgrade-Header etabliert
- Der Server sendet **sofort** den aktuellen Spielzustand nach dem Verbindungsaufbau
- Alle weiteren Spielzustandsänderungen werden automatisch an alle verbundenen Clients gesendet
- **Wichtig:** Alle Aktionen (Spiel erstellen, beitreten, Züge machen) erfolgen weiterhin über die REST API
- Die WebSocket-Verbindung ist **nur zum Empfangen** von Updates gedacht

### Nachrichten-Protokoll

#### Server → Client

**1. Spielzustands-Update:**
```json
{
  "type": "state",
  "payload": {
    "gameId": "e2c58d69-5ee4-4779-b43f-7825494723c4",
    "board": [["X", "O", ""], ["", "X", ""], ["", "", "O"]],
    "currentTurn": "O",
    "status": "IN_PROGRESS",
    "winner": ""
  }
}
```

**2. Fehlermeldung:**
```json
{
  "type": "error",
  "payload": {
    "message": "Game not found"
  }
}
```

#### Client → Server

- **Keine Nachrichten vom Client erforderlich**
- Alle Aktionen erfolgen über REST API

### Frontend-Integration (JavaScript/TypeScript)

**Beispiel-Implementierung:**

```javascript
// WebSocket-Verbindung aufbauen
const gameId = "e2c58d69-5ee4-4779-b43f-7825494723c4";
const ws = new WebSocket(`ws://localhost:8080/ws/games/${gameId}`);

// Verbindung geöffnet
ws.onopen = () => {
  console.log("WebSocket verbunden");
};

// Nachricht empfangen
ws.onmessage = (event) => {
  const message = JSON.parse(event.data);
  
  if (message.type === "state") {
    // UI mit neuem Spielzustand aktualisieren
    updateGameBoard(message.payload.board);
    updateCurrentTurn(message.payload.currentTurn);
    updateStatus(message.payload.status);
    
    if (message.payload.winner) {
      showGameOver(message.payload.winner);
    }
  } else if (message.type === "error") {
    console.error("WebSocket Fehler:", message.payload.message);
    showError(message.payload.message);
  }
};

// Fehlerbehandlung
ws.onerror = (error) => {
  console.error("WebSocket Verbindungsfehler:", error);
  // Fallback: Polling mit REST API
  startPolling();
};

// Verbindung geschlossen
ws.onclose = () => {
  console.log("WebSocket Verbindung geschlossen");
  // Optional: Reconnect-Logik oder Fallback zu Polling
  attemptReconnect();
};
```

**Empfohlenes Muster:**
- **REST API** für alle Aktionen (Spiel erstellen, beitreten, Züge machen)
- **WebSocket** für Echtzeit-Updates (Gegnerzüge, Spielzustandsänderungen)
- **Fallback zu Polling** (`GET /games/{gameId}`) falls WebSocket-Verbindung fehlschlägt

---

## Beispiel-Skripte

Im Repository finden Sie unter `scripts/` verschiedene Shell-Skripte, die die API-Nutzung demonstrieren:

### Vollständiger Test-Ablauf

```bash
# 1. Spieler erstellen
PLAYER_ID_ALICE=$(./scripts/01_create-player.sh "Alice" | jq -r '.playerId')

# 2. Spiel erstellen
GAME_ID=$(PLAYER_ID="$PLAYER_ID_ALICE" MODE=PVP ./scripts/02_create-game.sh | jq -r '.gameId')

# 3. Verfügbare Spiele auflisten
MODE=PVP STATUS=WAITING_FOR_PLAYER ./scripts/03_list-games.sh

# 4. Zweiten Spieler erstellen und beitreten
PLAYER_ID_BOB=$(./scripts/01_create-player.sh "Bob" | jq -r '.playerId')
PLAYER_ID="$PLAYER_ID_BOB" GAME_ID="$GAME_ID" ./scripts/05_join-game.sh

# 5. Züge machen
PLAYER_ID="$PLAYER_ID_ALICE" GAME_ID="$GAME_ID" ROW=0 COL=0 ./scripts/06_make-move.sh
PLAYER_ID="$PLAYER_ID_BOB" GAME_ID="$GAME_ID" ROW=0 COL=1 ./scripts/06_make-move.sh
```

### WebSocket-Test

```bash
# Vollständiger Test mit WebSocket-Updates
go run scripts/07_websocket-test.go
```

Diese Skripte dienen als **Referenz** für die API-Nutzung und können beim Debugging helfen.

---

## Anforderungen

### Funktionale Anforderungen

1. **Spieler-Registrierung**
   - Spieler kann seinen Namen eingeben
   - `playerId` wird im `localStorage` gespeichert
   - Bei erneutem Besuch wird der Spieler automatisch erkannt

2. **Spielmodus-Auswahl**
   - Spieler kann zwischen **PVP** und **PVC** wählen
   - Bei PVP: Spieler kann neues Spiel erstellen ODER einem bestehenden Spiel beitreten

3. **Spieloberfläche**
   - 3x3 Spielfeld visualisieren
   - Aktuellen Spielzustand anzeigen (Board, aktueller Zug, Status)
   - Klare visuelle Unterscheidung zwischen X und O
   - Leere Felder als klickbar markieren

4. **Spielablauf**
   - Spieler kann nur Züge machen, wenn er am Zug ist
   - Ungültige Züge werden verhindert/angezeigt
   - Nach jedem Zug wird der Spielzustand aktualisiert
   - Bei PVC: KI-Zug wird automatisch angezeigt (nach Spielerzug)

5. **Echtzeit-Updates**
   - WebSocket-Verbindung für Echtzeit-Updates
   - Gegnerzüge werden sofort angezeigt (ohne Seitenaktualisierung)
   - Fallback zu Polling, falls WebSocket nicht verfügbar

6. **Spielende**
   - Gewinner wird angezeigt
   - Unentschieden wird angezeigt
   - Möglichkeit, ein neues Spiel zu starten

7. **Fehlerbehandlung**
   - Netzwerkfehler werden angezeigt
   - Ungültige Aktionen werden verhindert/erklärt
   - Benutzerfreundliche Fehlermeldungen

### Nicht-funktionale Anforderungen

1. **Benutzerfreundlichkeit**
   - Intuitive Bedienung
   - Klare visuelle Feedback
   - Responsive Design (funktioniert auf Desktop und Tablet)

2. **Code-Qualität**
   - Strukturierter, lesbarer Code
   - Wiederverwendbare Komponenten
   - Kommentare wo nötig

3. **Performance**
   - Schnelle Reaktionszeiten
   - Keine unnötigen API-Aufrufe

---

## Implementierungshinweise

### Projekt-Setup

#### Vue.js

```bash
# Neues Vue-Projekt erstellen
npm create vue@latest tic-tac-toe-frontend
cd tic-tac-toe-frontend
npm install

# Optional: Axios für HTTP-Requests
npm install axios
```

**Empfohlene Projektstruktur:**
```
src/
  components/
    GameBoard.vue      # Spielfeld-Komponente
    GameStatus.vue    # Status-Anzeige
    PlayerForm.vue    # Spieler-Registrierung
    GameList.vue      # Liste verfügbarer Spiele
  services/
    api.js            # REST API Client
    websocket.js      # WebSocket Client
  stores/
    game.js           # State Management (optional: Pinia)
  App.vue
  main.js
```

#### React.js

```bash
# Neues React-Projekt erstellen
npx create-react-app tic-tac-toe-frontend
cd tic-tac-toe-frontend

# Optional: Axios für HTTP-Requests
npm install axios
```

**Empfohlene Projektstruktur:**
```
src/
  components/
    GameBoard.jsx     # Spielfeld-Komponente
    GameStatus.jsx    # Status-Anzeige
    PlayerForm.jsx    # Spieler-Registrierung
    GameList.jsx      # Liste verfügbarer Spiele
  services/
    api.js            # REST API Client
    websocket.js      # WebSocket Client
  hooks/
    useGame.js        # Custom Hook für Spiel-Logik
  App.js
  index.js
```

### API-Client-Implementierung

**Wichtig:** Die API-Base-URL zeigt auf den **lokalen Server**. Stellen Sie sicher, dass der Backend-Server läuft, bevor Sie das Frontend starten.

**Beispiel (Vue.js / JavaScript):**

```javascript
// services/api.js
// API-Base-URL für lokalen Server
const API_BASE = 'http://localhost:8080';

// Optional: Für Entwicklung vs. Produktion
// const API_BASE = import.meta.env.VITE_API_BASE || 'http://localhost:8080';

export const api = {
  async createPlayer(name) {
    const response = await fetch(`${API_BASE}/players`, {
      method: 'POST',
      headers: { 'Content-Type': 'application/json' },
      body: JSON.stringify({ name })
    });
    if (!response.ok) throw new Error('Failed to create player');
    return response.json();
  },

  async createGame(playerId, mode) {
    const response = await fetch(`${API_BASE}/games`, {
      method: 'POST',
      headers: {
        'Content-Type': 'application/json',
        'X-Player-Id': playerId
      },
      body: JSON.stringify({ mode })
    });
    if (!response.ok) throw new Error('Failed to create game');
    return response.json();
  },

  async listGames(mode, status) {
    const params = new URLSearchParams({ mode, status });
    const response = await fetch(`${API_BASE}/games?${params}`);
    if (!response.ok) throw new Error('Failed to list games');
    return response.json();
  },

  async joinGame(playerId, gameId) {
    const response = await fetch(`${API_BASE}/games/${gameId}/join`, {
      method: 'POST',
      headers: { 'X-Player-Id': playerId }
    });
    if (!response.ok) throw new Error('Failed to join game');
    return response.json();
  },

  async makeMove(playerId, gameId, row, col) {
    const response = await fetch(`${API_BASE}/games/${gameId}/moves`, {
      method: 'POST',
      headers: {
        'Content-Type': 'application/json',
        'X-Player-Id': playerId
      },
      body: JSON.stringify({ row, col })
    });
    if (!response.ok) throw new Error('Failed to make move');
    return response.json();
  },

  async getGame(gameId) {
    const response = await fetch(`${API_BASE}/games/${gameId}`);
    if (!response.ok) throw new Error('Failed to get game');
    return response.json();
  }
};
```

### WebSocket-Client-Implementierung

**Wichtig:** Die WebSocket-URL zeigt auf den **lokalen Server**. Der Server muss laufen, damit die WebSocket-Verbindung funktioniert.

**Beispiel:**

```javascript
// services/websocket.js
export class GameWebSocket {
  constructor(gameId, onMessage) {
    this.gameId = gameId;
    this.onMessage = onMessage;
    this.ws = null;
    this.reconnectAttempts = 0;
    this.maxReconnectAttempts = 5;
    // WebSocket-Base-URL für lokalen Server
    this.wsBase = 'ws://localhost:8080';
  }

  connect() {
    const url = `${this.wsBase}/ws/games/${this.gameId}`;
    this.ws = new WebSocket(url);

    this.ws.onopen = () => {
      console.log('WebSocket verbunden');
      this.reconnectAttempts = 0;
    };

    this.ws.onmessage = (event) => {
      const message = JSON.parse(event.data);
      this.onMessage(message);
    };

    this.ws.onerror = (error) => {
      console.error('WebSocket Fehler:', error);
    };

    this.ws.onclose = () => {
      console.log('WebSocket geschlossen');
      this.attemptReconnect();
    };
  }

  attemptReconnect() {
    if (this.reconnectAttempts < this.maxReconnectAttempts) {
      this.reconnectAttempts++;
      setTimeout(() => {
        console.log(`Reconnect-Versuch ${this.reconnectAttempts}...`);
        this.connect();
      }, 1000 * this.reconnectAttempts);
    } else {
      console.error('Maximale Reconnect-Versuche erreicht');
      // Fallback zu Polling
    }
  }

  disconnect() {
    if (this.ws) {
      this.ws.close();
      this.ws = null;
    }
  }
}
```

### Komponenten-Beispiele

**GameBoard-Komponente (Vue.js):**

```vue
<template>
  <div class="game-board">
    <div
      v-for="(row, rowIndex) in board"
      :key="rowIndex"
      class="board-row"
    >
      <button
        v-for="(cell, colIndex) in row"
        :key="colIndex"
        class="board-cell"
        :disabled="!canMakeMove(rowIndex, colIndex)"
        @click="handleCellClick(rowIndex, colIndex)"
      >
        {{ cell || ' ' }}
      </button>
    </div>
  </div>
</template>

<script>
export default {
  props: {
    board: {
      type: Array,
      required: true
    },
    currentTurn: String,
    playerSymbol: String,
    status: String
  },
  methods: {
    canMakeMove(row, col) {
      return (
        this.status === 'IN_PROGRESS' &&
        this.currentTurn === this.playerSymbol &&
        this.board[row][col] === ''
      );
    },
    handleCellClick(row, col) {
      if (this.canMakeMove(row, col)) {
        this.$emit('move', row, col);
      }
    }
  }
};
</script>
```

---

## Testanweisungen

### Manuelle Tests

1. **Spieler-Registrierung**
   - Neuen Spieler erstellen
   - Seite neu laden → Spieler sollte wiedererkannt werden

2. **PVP-Spiel**
   - Spieler A erstellt PVP-Spiel
   - Spieler B sieht Spiel in Liste
   - Spieler B tritt bei
   - Beide Spieler können abwechselnd Züge machen
   - WebSocket-Updates funktionieren in Echtzeit

3. **PVC-Spiel**
   - Spieler erstellt PVC-Spiel
   - Spieler macht Zug
   - KI macht automatisch Zug
   - Spieler kann weiter machen

4. **Fehlerbehandlung**
   - Ungültige Züge werden verhindert
   - Netzwerkfehler werden angezeigt
   - WebSocket-Verbindungsfehler werden behandelt

### Automatisierte Tests (Optional)

- Unit-Tests für Komponenten
- Integration-Tests für API-Aufrufe
- E2E-Tests mit Playwright oder Cypress

---

## Bewertungskriterien

### Funktionale Korrektheit (40%)
- Alle Anforderungen sind implementiert
- Spielablauf funktioniert korrekt (PVP und PVC)
- WebSocket-Updates funktionieren
- Fehlerbehandlung ist vorhanden

### Code-Qualität (30%)
- Strukturierter, lesbarer Code
- Wiederverwendbare Komponenten
- Angemessene Kommentierung
- Keine offensichtlichen Bugs

### Benutzerfreundlichkeit (20%)
- Intuitive Bedienung
- Klare visuelle Feedback
- Responsive Design
- Gute Fehlermeldungen

### Dokumentation (10%)
- README mit Setup-Anweisungen
- Code-Kommentare wo nötig
- API-Integration dokumentiert

---

## Abgabe

### Abgabeformat

1. **Git-Repository** (z.B. GitHub, GitLab)
   - Vollständiger Quellcode des Frontends
   - README.md mit Setup-Anweisungen
   - `.gitignore` korrekt konfiguriert
   - **Hinweis:** Nur das Frontend-Repository abgeben, nicht das Backend-Server-Repository

2. **Kurze Dokumentation** (README.md)
   - Projekt-Setup (inkl. Backend-Server-Start)
   - Verwendete Technologien
   - Bekannte Probleme/Limitationen
   - **Wichtig:** Dokumentieren Sie, dass der Backend-Server lokal gestartet werden muss

### Setup-Anweisungen für Abgabe

Ihr README.md sollte folgende Informationen enthalten:

```markdown
## Setup

### Voraussetzungen
- Node.js und npm installiert
- Go installiert (für Backend-Server)
- Backend-Server-Repository verfügbar

### Backend-Server starten
1. Backend-Server-Repository klonen
2. `cd` in das Server-Verzeichnis
3. `go run ./cmd/server` ausführen
4. Server läuft auf http://localhost:8080

### Frontend starten
1. `npm install`
2. `npm run dev` (oder entsprechendes Start-Kommando)
3. Frontend läuft auf http://localhost:5173 (oder Port Ihrer Wahl)
```

### Abgabefrist

Siehe Kurs-Informationen.

---

## Häufige Fragen (FAQ)

### Muss ich einen gehosteten Server verwenden?

**Nein!** Alles läuft lokal auf Ihrem Rechner:
- Backend-Server: `http://localhost:8080`
- Frontend: `http://localhost:5173` (oder ähnlich)
- Keine externe Verbindung erforderlich

### Wie speichere ich die playerId?

Verwenden Sie `localStorage`:

```javascript
// Speichern
localStorage.setItem('playerId', playerId);

// Abrufen
const playerId = localStorage.getItem('playerId');
```

### Was passiert, wenn der Server nicht erreichbar ist?

Zeigen Sie eine benutzerfreundliche Fehlermeldung an und ermöglichen Sie einen Retry.

**Häufige Ursachen:**
- Server wurde nicht gestartet → Starten Sie den Server mit `go run ./cmd/server`
- Server läuft auf anderem Port → Überprüfen Sie die Port-Konfiguration
- Port 8080 ist bereits belegt → Ändern Sie den Port mit `TICTACGO_PORT=9090`
- CORS-Fehler → Stellen Sie sicher, dass der Server CORS für `localhost` erlaubt (sollte bereits konfiguriert sein)

### Wie teste ich, ob der Server läuft?

```bash
# In einem Terminal
curl http://localhost:8080/health

# Erwartete Antwort:
# {"status":"ok"}
```

Falls Sie eine Fehlermeldung erhalten, ist der Server nicht erreichbar. Überprüfen Sie:
1. Wurde der Server gestartet?
2. Läuft der Server auf Port 8080?
3. Gibt es Firewall-Probleme?

### Soll ich Polling oder WebSocket verwenden?

**Empfehlung:** WebSocket für Echtzeit-Updates, mit Fallback zu Polling falls WebSocket fehlschlägt.

### Wie handle ich WebSocket-Verbindungsfehler?

Implementieren Sie eine Reconnect-Logik mit exponentieller Backoff-Strategie. Falls alle Reconnect-Versuche fehlschlagen, wechseln Sie zu Polling.

### Kann ich State Management (Vuex/Pinia, Redux) verwenden?

Ja, ist optional aber empfohlen für größere Anwendungen. Für diese Aufgabe reicht lokaler State in Komponenten aus.

### Kann ich den Server und das Frontend auf verschiedenen Ports laufen lassen?

Ja, das ist sogar der Standard:
- Backend-Server: Port 8080
- Frontend (Vite/Vue): Port 5173
- Frontend (React): Port 3000

Das Frontend kommuniziert mit dem Backend über HTTP/WebSocket, daher funktioniert dies problemlos.

### Was mache ich, wenn Port 8080 bereits belegt ist?

Sie haben zwei Optionen:

1. **Backend-Port ändern:**
   ```bash
   TICTACGO_PORT=9090 go run ./cmd/server
   ```
   Dann im Frontend die API-Base-URL anpassen:
   ```javascript
   const API_BASE = 'http://localhost:9090';
   ```

2. **Anderen Prozess auf Port 8080 beenden** (falls möglich)

---

## Zusätzliche Ressourcen

- **Vue.js Dokumentation:** https://vuejs.org/
- **React.js Dokumentation:** https://react.dev/
- **WebSocket API:** https://developer.mozilla.org/en-US/docs/Web/API/WebSocket
- **Fetch API:** https://developer.mozilla.org/en-US/docs/Web/API/Fetch_API

---

## Support

Bei Fragen wenden Sie sich an:
- Kurs-Forum
- Sprechstunden
- E-Mail an Kursleiter

**Viel Erfolg bei der Implementierung!**

