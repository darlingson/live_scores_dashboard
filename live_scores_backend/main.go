package main

import (
	"log"
	"net/http"

	"github.com/gorilla/websocket"
)

var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
	// Allow all connections by bypassing the origin check
	CheckOrigin: func(r *http.Request) bool {
		return true
	},
}

func main() {
	log.Println("Starting server on port 8080...")
	http.HandleFunc("/ws", handleWebSocket)
	log.Fatal(http.ListenAndServe(":8080", nil))
}

func handleWebSocket(w http.ResponseWriter, r *http.Request) {
	log.Println("Incoming WebSocket connection...")
	c, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Println("Error upgrading connection:", err)
		return
	}
	defer c.Close()
	log.Println("WebSocket connection established.")

	// Start reading messages
	for {
		messageType, message, err := c.ReadMessage()
		if err != nil {
			log.Println("Error reading message:", err)
			break
		}
		log.Printf("Received message: %s", message)

		// Echo the message back
		if err := c.WriteMessage(messageType, message); err != nil {
			log.Println("Error writing message:", err)
			break
		}
	}
}
