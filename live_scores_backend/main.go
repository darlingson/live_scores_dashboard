package main

import (
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/gorilla/websocket" // Using Gorilla WebSocket for easier handling
)

// Game represents the structure of a football game
type Game struct {
	ID         string    `json:"id"`
	HomeTeam   string    `json:"homeTeam"`
	AwayTeam   string    `json:"awayTeam"`
	HomeScore  int       `json:"homeScore"`
	AwayScore  int       `json:"awayScore"`
	Scorers    []Scorer  `json:"scorers"`
	Status     string    `json:"status"` // e.g., "pending", "active", "finished"
	LastUpdate time.Time `json:"lastUpdate"`
}

// Scorer represents a goal scorer
type Scorer struct {
	PlayerName string `json:"playerName"`
	Team       string `json:"team"`
	Minute     int    `json:"minute"`
}

// Message represents the WebSocket message format
type Message struct {
	Type string      `json:"type"` // e.g., "gameUpdate", "initialGames"
	Data interface{} `json:"data"`
}

// We'll use a map to store currently active games, keyed by game ID
var activeGames = make(map[string]*Game)

// Configure the upgrader to allow cross-origin requests (for frontend running on different port)
var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
	CheckOrigin: func(r *http.Request) bool {
		return true // Allow all origins for development
	},
}

// broadcastMessage sends a message to all connected clients
func broadcastMessage(msg Message) {
	// In a real application, you'd manage a pool of connections and iterate them
	// For this example, we'll just log that a message would be sent
	// (we'll connect only one client for simplicity, or handle multiple if needed later)
	// For now, let's just assume ws is the connection from handleConnections
	// In a full app, you'd have a map of connections to iterate and send to.
	// For this simplified example, this function would need to know about all active websockets.
	// We'll simulate this with a global slice of connections for this example.
	for client := range clients {
		err := client.WriteJSON(msg)
		if err != nil {
			log.Printf("Error broadcasting to client: %v", err)
			client.Close()
			delete(clients, client)
		}
	}
}

// clients holds all active WebSocket connections (simplified for this example)
var clients = make(map[*websocket.Conn]bool) // Concurrent map for storing clients
var register = make(chan *websocket.Conn)
var unregister = make(chan *websocket.Conn)
var broadcast = make(chan Message)

func handleMessages() {
	for {
		select {
		case conn := <-register:
			clients[conn] = true
			log.Println("Client registered. Total clients:", len(clients))
		case conn := <-unregister:
			delete(clients, conn)
			log.Println("Client unregistered. Total clients:", len(clients))
		case message := <-broadcast:
			for client := range clients {
				err := client.WriteJSON(message)
				if err != nil {
					log.Printf("error: %v", err)
					client.Close()
					delete(clients, client)
				}
			}
		}
	}
}

func removeClient(clientWs *websocket.Conn) {
	unregister <- clientWs
}

// simulateGameUpdates simulates real-time game updates
func simulateGameUpdates() {
	// Initialize some dummy games
	activeGames["game1"] = &Game{
		ID:         "game1",
		HomeTeam:   "Nyasa Big Bullets",
		AwayTeam:   "Mighty Wanderers",
		HomeScore:  0,
		AwayScore:  0,
		Scorers:    []Scorer{},
		Status:     "active",
		LastUpdate: time.Now(),
	}
	activeGames["game2"] = &Game{
		ID:         "game2",
		HomeTeam:   "Silver Strikers",
		AwayTeam:   "Civil Service United",
		HomeScore:  0,
		AwayScore:  0,
		Scorers:    []Scorer{},
		Status:     "pending",
		LastUpdate: time.Now(),
	}
	activeGames["game3"] = &Game{
		ID:         "game3",
		HomeTeam:   "Blue Eagles",
		AwayTeam:   "Ekwendeni Hammers",
		HomeScore:  1,
		AwayScore:  0,
		Scorers:    []Scorer{{"Chikondi Banda", "Blue Eagles", 25}},
		Status:     "active",
		LastUpdate: time.Now(),
	}

	ticker := time.NewTicker(5 * time.Second) // Update every 5 seconds
	defer ticker.Stop()

	for range ticker.C {
		// Simulate a score change for game1
		game1 := activeGames["game1"]
		if game1.Status == "active" {
			if game1.HomeScore < 3 { // Cap scores for demo
				game1.HomeScore++
				scorer := Scorer{PlayerName: fmt.Sprintf("Player %d", game1.HomeScore), Team: game1.HomeTeam, Minute: int(time.Now().Unix()%90) + 1}
				game1.Scorers = append(game1.Scorers, scorer)
				game1.LastUpdate = time.Now()
				log.Printf("Simulating goal for %s: %s %d-%d %s", game1.ID, game1.HomeTeam, game1.HomeScore, game1.AwayScore, game1.AwayTeam)
				broadcast <- Message{Type: "gameUpdate", Data: game1}
			} else {
				// Finish game1 after max score
				game1.Status = "finished"
				game1.LastUpdate = time.Now()
				log.Printf("Game %s has finished!", game1.ID)
				broadcast <- Message{Type: "gameUpdate", Data: game1}
			}
		}

		// Simulate game2 starting
		game2 := activeGames["game2"]
		if game2.Status == "pending" {
			game2.Status = "active"
			game2.LastUpdate = time.Now()
			log.Printf("Game %s has started!", game2.ID)
			broadcast <- Message{Type: "gameUpdate", Data: game2}
		}
	}
}

func main() {
	go handleMessages()      // Start the message handling goroutine
	go simulateGameUpdates() // Start simulating game updates in a goroutine

	http.HandleFunc("/ws", func(w http.ResponseWriter, r *http.Request) {
		ws, err := upgrader.Upgrade(w, r, nil)
		if err != nil {
			log.Println("Upgrade error:", err)
			return
		}

		register <- ws // Register new client
		defer func() {
			unregister <- ws
			ws.Close()
		}()

		log.Println("New WebSocket client connected!")

		// Send initial game data
		gamesList := []*Game{}
		for _, game := range activeGames {
			gamesList = append(gamesList, game)
		}
		initialMsg := Message{
			Type: "initialGames",
			Data: gamesList,
		}
		if err := ws.WriteJSON(initialMsg); err != nil {
			log.Printf("Error sending initial games: %v", err)
			return
		}

		for {
			_, _, err := ws.ReadMessage()
			if err != nil {
				log.Printf("Error reading message: %v", err)
				break
			}
		}
	})

	log.Println("Go WebSocket server started on :8080")
	log.Fatal(http.ListenAndServe(":8080", nil))
}
