package presentation

import (
	"fmt"
	"net/http"
	"sync"
	"github.com/gorilla/websocket"
)

// We configure the upgrader to accept connections from your Vercel domain
var upgrader = websocket.Upgrader{
	CheckOrigin: func(r *http.Request) bool {
		return true 
	},
}

// Hub maintains the list of active clients and broadcasts messages
type Hub struct {
	clients map[*websocket.Conn]bool
	mu      sync.Mutex // Mutex prevents crashes if two people connect at the exact same millisecond
}

func NewHub() *Hub {
	return &Hub{
		clients: make(map[*websocket.Conn]bool),
	}
}

// HandleConnections upgrades the HTTP request to a WebSocket
func (h *Hub) HandleConnections(w http.ResponseWriter, r *http.Request) {
	ws, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		fmt.Println("WebSocket Upgrade Error:", err)
		return
	}
	
	// Add the new browser tab to our list of clients
	h.mu.Lock()
	h.clients[ws] = true
	h.mu.Unlock()

	fmt.Println("🔌 New WebSocket client connected!")

	// Keep the connection alive until the user closes their browser
	for {
		if _, _, err := ws.NextReader(); err != nil {
			h.mu.Lock()
			delete(h.clients, ws)
			h.mu.Unlock()
			ws.Close()
			break
		}
	}
}

// Broadcast sends a "REFRESH" signal to every connected browser
func (h *Hub) Broadcast() {
	h.mu.Lock()
	defer h.mu.Unlock()
	
	for client := range h.clients {
		err := client.WriteMessage(websocket.TextMessage, []byte("REFRESH"))
		if err != nil {
			client.Close()
			delete(h.clients, client)
		}
	}
}