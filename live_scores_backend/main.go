package main

import (
    "encoding/json"
    "fmt"
    "log"
    "net/http"
    "sync"

    "github.com/gorilla/websocket"
)

var upgrader = websocket.Upgrader{
    CheckOrigin: func(r *http.Request) bool {
        return true // Allow all connections (for simplicity)
    },
}
type MatchEvent struct {
    Type       string `json:"type"`
    Scorer     string `json:"scorer"`
    Time       string `json:"time"`
    Score      string `json:"score"`
    Message    string `json:"message"`
}

var clients = make(map[*websocket.Conn]bool)
var mutex = &sync.Mutex{}

func broadcast(event MatchEvent) {
    message, err := json.Marshal(event)
    if err != nil {
        log.Printf("Error marshalling event: %v", err)
        return
    }

    mutex.Lock()
    defer mutex.Unlock()

    for client := range clients {
        err := client.WriteMessage(websocket.TextMessage, message)
        if err != nil {
            log.Printf("Error broadcasting to client: %v", err)
            client.Close()
            delete(clients, client)
        }
    }
}

func handleWebSocket(w http.ResponseWriter, r *http.Request) {
    conn, err := upgrader.Upgrade(w, r, nil)
    if err != nil {
        log.Printf("WebSocket upgrade error: %v", err)
        return
    }
    defer conn.Close()

    mutex.Lock()
    clients[conn] = true
    mutex.Unlock()

    for {
        _, _, err := conn.ReadMessage()
        if err != nil {
            log.Printf("Error reading WebSocket message: %v", err)
            break
        }
    }
}

func handleMatchEvent(w http.ResponseWriter, r *http.Request) {
    var event MatchEvent
    err := json.NewDecoder(r.Body).Decode(&event)
    if err != nil {
        http.Error(w, "Invalid request body", http.StatusBadRequest)
        return
    }

    broadcast(event)

    w.WriteHeader(http.StatusOK)
    fmt.Fprintf(w, "Event received: %+v", event)
}

func main() {
    http.HandleFunc("/ws", handleWebSocket)
    http.HandleFunc("/event", handleMatchEvent)
    port := ":8080"
    fmt.Printf("Server started on port %s\n", port)
    log.Fatal(http.ListenAndServe(port, nil))
}
