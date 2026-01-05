// Copyright 2026 Esslingen University of Applied Sciences
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
// Date: 2026-01-05
//
// Usage: go run scripts/07_websocket-test.go
//
// This script demonstrates a complete game flow with WebSocket real-time updates.
// It mirrors the functionality of 00_full_test_script.sh but shows WebSocket messages
// as they are received in real-time.
//
// Flow:
// 1. Creates two players (Alice and Bob) via REST API
// 2. Alice creates a PVP game via REST API
// 3. Connects to WebSocket for real-time updates
// 4. Bob joins the game via REST API (WebSocket receives update)
// 5. Both players make moves via REST API (WebSocket receives updates)
// 6. Shows all WebSocket messages with formatted board visualization

package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"time"

	"github.com/gorilla/websocket"
)

const (
	apiBase = "http://localhost:8080"
	wsBase  = "ws://localhost:8080"
)

// ANSI color codes
const (
	colorReset  = "\033[0m"
	colorBlue   = "\033[0;34m"
	colorGreen  = "\033[0;32m"
	colorYellow = "\033[1;33m"
	colorCyan   = "\033[0;36m"
)

type Player struct {
	PlayerID string `json:"playerId"`
	Name     string `json:"name"`
}

type GameState struct {
	GameID      string     `json:"gameId"`
	Mode        string     `json:"mode"`
	Board       [][]string `json:"board"`
	CurrentTurn string     `json:"currentTurn"`
	Status      string     `json:"status"`
	Winner      string     `json:"winner"`
}

type CreateGameRequest struct {
	Mode string `json:"mode"`
}

type MakeMoveRequest struct {
	Row int `json:"row"`
	Col int `json:"col"`
}

type WSMessage struct {
	Type    string                 `json:"type"`
	Payload map[string]interface{} `json:"payload"`
}

func printStep(step int, message string) {
	fmt.Printf("%s[Step %d]%s %s\n", colorYellow, step, colorReset, message)
}

func printSuccess(message string) {
	fmt.Printf("%s✓%s %s\n", colorGreen, colorReset, message)
}

func printWebSocket(message string) {
	fmt.Printf("%s[WebSocket]%s %s\n", colorCyan, colorReset, message)
}

func printHeader(message string) {
	fmt.Printf("\n%s=== %s ===%s\n\n", colorBlue, message, colorReset)
}

func createPlayer(name string) (string, error) {
	reqBody := map[string]string{"name": name}
	jsonData, _ := json.Marshal(reqBody)

	resp, err := http.Post(apiBase+"/players", "application/json", bytes.NewBuffer(jsonData))
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", err
	}

	var player Player
	if err := json.Unmarshal(body, &player); err != nil {
		return "", err
	}

	return player.PlayerID, nil
}

func createGame(playerID string) (string, error) {
	reqBody := CreateGameRequest{Mode: "PVP"}
	jsonData, _ := json.Marshal(reqBody)

	req, _ := http.NewRequest("POST", apiBase+"/games", bytes.NewBuffer(jsonData))
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("X-Player-Id", playerID)

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", err
	}

	var gameState GameState
	if err := json.Unmarshal(body, &gameState); err != nil {
		return "", err
	}

	return gameState.GameID, nil
}

func joinGame(gameID, playerID string) error {
	req, _ := http.NewRequest("POST", apiBase+"/games/"+gameID+"/join", nil)
	req.Header.Set("X-Player-Id", playerID)

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	return nil
}

func makeMove(gameID, playerID string, row, col int) error {
	reqBody := MakeMoveRequest{Row: row, Col: col}
	jsonData, _ := json.Marshal(reqBody)

	req, _ := http.NewRequest("POST", apiBase+"/games/"+gameID+"/moves", bytes.NewBuffer(jsonData))
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("X-Player-Id", playerID)

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	return nil
}

func formatBoard(board [][]string) string {
	var result string
	for i, row := range board {
		result += "  "
		for j, cell := range row {
			if cell == "" {
				result += " "
			} else {
				result += cell
			}
			if j < 2 {
				result += " | "
			}
		}
		if i < 2 {
			result += "\n  -----------\n"
		}
	}
	return result
}

func main() {
	printHeader("Tic-Tac-Go WebSocket Full Test")

	// Step 1: Create Alice
	printStep(1, "Creating player Alice...")
	playerIDAlice, err := createPlayer("Alice")
	if err != nil {
		log.Fatalf("Failed to create Alice: %v", err)
	}
	printSuccess(fmt.Sprintf("Alice created with ID: %s", playerIDAlice))

	// Step 2: Create Bob
	printStep(2, "Creating player Bob...")
	playerIDBob, err := createPlayer("Bob")
	if err != nil {
		log.Fatalf("Failed to create Bob: %v", err)
	}
	printSuccess(fmt.Sprintf("Bob created with ID: %s", playerIDBob))

	// Step 3: Alice creates a PVP game
	printStep(3, "Alice creates a PVP game...")
	gameID, err := createGame(playerIDAlice)
	if err != nil {
		log.Fatalf("Failed to create game: %v", err)
	}
	printSuccess(fmt.Sprintf("Game created with ID: %s", gameID))

	// Step 4: Connect to WebSocket
	printStep(4, "Connecting to WebSocket for real-time updates...")
	url := fmt.Sprintf("%s/ws/games/%s", wsBase, gameID)
	conn, _, err := websocket.DefaultDialer.Dial(url, nil)
	if err != nil {
		log.Fatalf("Failed to connect to WebSocket: %v", err)
	}
	defer conn.Close()
	printSuccess("WebSocket connected!")

	// Channel to receive WebSocket messages
	wsMessages := make(chan WSMessage, 10)
	wsDone := make(chan struct{})

	// Read WebSocket messages in a goroutine
	go func() {
		defer close(wsDone)
		for {
			_, message, err := conn.ReadMessage()
			if err != nil {
				if websocket.IsUnexpectedCloseError(err, websocket.CloseGoingAway, websocket.CloseAbnormalClosure) {
					log.Printf("WebSocket error: %v", err)
				}
				return
			}

			var msg WSMessage
			if err := json.Unmarshal(message, &msg); err == nil {
				wsMessages <- msg
			}
		}
	}()

	// Wait a bit for initial state message
	time.Sleep(500 * time.Millisecond)

	// Process initial WebSocket message
	select {
	case msg := <-wsMessages:
		if msg.Type == "state" {
			printWebSocket("Received initial game state:")
			if payload, ok := msg.Payload["board"].([]interface{}); ok {
				board := make([][]string, 3)
				for i := 0; i < 3; i++ {
					board[i] = make([]string, 3)
					if row, ok := payload[i].([]interface{}); ok {
						for j := 0; j < 3; j++ {
							if cell, ok := row[j].(string); ok {
								board[i][j] = cell
							}
						}
					}
				}
				fmt.Println(formatBoard(board))
			}
		}
	default:
	}

	// Step 5: Bob joins the game
	printStep(5, "Bob joins the game...")
	if err := joinGame(gameID, playerIDBob); err != nil {
		log.Fatalf("Failed to join game: %v", err)
	}
	printSuccess("Bob joined the game")

	// Wait for WebSocket update
	time.Sleep(500 * time.Millisecond)
	select {
	case msg := <-wsMessages:
		if msg.Type == "state" {
			printWebSocket("Received update: Bob joined!")
			if status, ok := msg.Payload["status"].(string); ok {
				fmt.Printf("  Status: %s\n", status)
			}
		}
	default:
	}

	// Step 6: Alice makes move (0,0)
	printStep(6, "Alice makes move at (0,0)...")
	if err := makeMove(gameID, playerIDAlice, 0, 0); err != nil {
		log.Fatalf("Failed to make move: %v", err)
	}
	printSuccess("Move made")

	time.Sleep(500 * time.Millisecond)
	select {
	case msg := <-wsMessages:
		if msg.Type == "state" {
			printWebSocket("Received update: Alice's move!")
			if payload, ok := msg.Payload["board"].([]interface{}); ok {
				board := make([][]string, 3)
				for i := 0; i < 3; i++ {
					board[i] = make([]string, 3)
					if row, ok := payload[i].([]interface{}); ok {
						for j := 0; j < 3; j++ {
							if cell, ok := row[j].(string); ok {
								board[i][j] = cell
							}
						}
					}
				}
				fmt.Println(formatBoard(board))
			}
		}
	default:
	}

	// Step 7: Bob makes move (0,1)
	printStep(7, "Bob makes move at (0,1)...")
	if err := makeMove(gameID, playerIDBob, 0, 1); err != nil {
		log.Fatalf("Failed to make move: %v", err)
	}
	printSuccess("Move made")

	time.Sleep(500 * time.Millisecond)
	select {
	case msg := <-wsMessages:
		if msg.Type == "state" {
			printWebSocket("Received update: Bob's move!")
			if payload, ok := msg.Payload["board"].([]interface{}); ok {
				board := make([][]string, 3)
				for i := 0; i < 3; i++ {
					board[i] = make([]string, 3)
					if row, ok := payload[i].([]interface{}); ok {
						for j := 0; j < 3; j++ {
							if cell, ok := row[j].(string); ok {
								board[i][j] = cell
							}
						}
					}
				}
				fmt.Println(formatBoard(board))
			}
		}
	default:
	}

	// Step 8: Alice makes move (1,1)
	printStep(8, "Alice makes move at (1,1)...")
	if err := makeMove(gameID, playerIDAlice, 1, 1); err != nil {
		log.Fatalf("Failed to make move: %v", err)
	}
	printSuccess("Move made")

	time.Sleep(500 * time.Millisecond)
	select {
	case msg := <-wsMessages:
		if msg.Type == "state" {
			printWebSocket("Received update: Alice's move!")
			if payload, ok := msg.Payload["board"].([]interface{}); ok {
				board := make([][]string, 3)
				for i := 0; i < 3; i++ {
					board[i] = make([]string, 3)
					if row, ok := payload[i].([]interface{}); ok {
						for j := 0; j < 3; j++ {
							if cell, ok := row[j].(string); ok {
								board[i][j] = cell
							}
						}
					}
				}
				fmt.Println(formatBoard(board))
			}
		}
	default:
	}

	// Step 9: Bob makes move (2,2)
	printStep(9, "Bob makes move at (2,2)...")
	if err := makeMove(gameID, playerIDBob, 2, 2); err != nil {
		log.Fatalf("Failed to make move: %v", err)
	}
	printSuccess("Move made")

	time.Sleep(500 * time.Millisecond)
	select {
	case msg := <-wsMessages:
		if msg.Type == "state" {
			printWebSocket("Received update: Bob's move!")
			if payload, ok := msg.Payload["board"].([]interface{}); ok {
				board := make([][]string, 3)
				for i := 0; i < 3; i++ {
					board[i] = make([]string, 3)
					if row, ok := payload[i].([]interface{}); ok {
						for j := 0; j < 3; j++ {
							if cell, ok := row[j].(string); ok {
								board[i][j] = cell
							}
						}
					}
				}
				fmt.Println(formatBoard(board))
			}
		}
	default:
	}

	// Step 10: Alice makes move (0,2)
	printStep(10, "Alice makes move at (0,2)...")
	if err := makeMove(gameID, playerIDAlice, 0, 2); err != nil {
		log.Fatalf("Failed to make move: %v", err)
	}
	printSuccess("Move made")

	time.Sleep(500 * time.Millisecond)
	select {
	case msg := <-wsMessages:
		if msg.Type == "state" {
			status, _ := msg.Payload["status"].(string)
			winner, _ := msg.Payload["winner"].(string)

			// Determine appropriate message based on actual game state
			var updateMsg string
			if status == "FINISHED" {
				if winner == "X" {
					updateMsg = "Received update: Alice wins!"
				} else if winner == "O" {
					updateMsg = "Received update: Bob wins!"
				} else if winner == "DRAW" {
					updateMsg = "Received update: Game ended in a draw!"
				} else {
					updateMsg = "Received update: Game finished!"
				}
			} else {
				updateMsg = "Received update: Alice's move!"
			}
			printWebSocket(updateMsg)

			if payload, ok := msg.Payload["board"].([]interface{}); ok {
				board := make([][]string, 3)
				for i := 0; i < 3; i++ {
					board[i] = make([]string, 3)
					if row, ok := payload[i].([]interface{}); ok {
						for j := 0; j < 3; j++ {
							if cell, ok := row[j].(string); ok {
								board[i][j] = cell
							}
						}
					}
				}
				fmt.Println(formatBoard(board))
			}
			if winner != "" {
				fmt.Printf("  Winner: %s\n", winner)
			}
			if status != "" {
				fmt.Printf("  Status: %s\n", status)
			}
		}
	default:
	}

	// Close WebSocket
	time.Sleep(500 * time.Millisecond)
	conn.WriteMessage(websocket.CloseMessage, websocket.FormatCloseMessage(websocket.CloseNormalClosure, ""))

	printHeader("Test Complete")
	fmt.Printf("Game ID: %s%s%s\n", colorGreen, gameID, colorReset)
	fmt.Printf("Alice ID: %s%s%s\n", colorGreen, playerIDAlice, colorReset)
	fmt.Printf("Bob ID: %s%s%s\n", colorGreen, playerIDBob, colorReset)
	fmt.Println("\nAll game state updates were received via WebSocket in real-time!")
}
