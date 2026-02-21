package websocket

import (
	"encoding/json"
	"log"
	"sync"
)

// Event represents a WebSocket event
type Event struct {
	Type string      `json:"type"`
	Data interface{} `json:"data"`
}

// Client represents a connected WebSocket client
type Client struct {
	ID     string
	UserID string
	Send   chan []byte
	hub    *Hub
}

// Hub maintains the set of active clients and broadcasts messages
type Hub struct {
	// Registered clients
	clients map[*Client]bool

	// Inbound messages from clients
	broadcast chan []byte

	// Register requests from clients
	register chan *Client

	// Unregister requests from clients
	unregister chan *Client

	// Mutex for thread-safe operations
	mu sync.RWMutex
}

// NewHub creates a new Hub instance
func NewHub() *Hub {
	return &Hub{
		clients:    make(map[*Client]bool),
		broadcast:  make(chan []byte, 256),
		register:   make(chan *Client),
		unregister: make(chan *Client),
	}
}

// Run starts the hub's main loop
func (h *Hub) Run() {
	for {
		select {
		case client := <-h.register:
			h.mu.Lock()
			h.clients[client] = true
			h.mu.Unlock()
			log.Printf("[WebSocket] Client registered: %s (User: %s)", client.ID, client.UserID)
			log.Printf("[WebSocket] Total clients: %d", len(h.clients))

		case client := <-h.unregister:
			h.mu.Lock()
			if _, ok := h.clients[client]; ok {
				delete(h.clients, client)
				close(client.Send)
				log.Printf("[WebSocket] Client unregistered: %s", client.ID)
				log.Printf("[WebSocket] Total clients: %d", len(h.clients))
			}
			h.mu.Unlock()

		case message := <-h.broadcast:
			h.mu.RLock()
			for client := range h.clients {
				select {
				case client.Send <- message:
				default:
					// Client's send channel is full, close it
					close(client.Send)
					delete(h.clients, client)
				}
			}
			h.mu.RUnlock()
		}
	}
}

// BroadcastEvent broadcasts an event to all connected clients
func (h *Hub) BroadcastEvent(eventType string, data interface{}) error {
	event := Event{
		Type: eventType,
		Data: data,
	}

	message, err := json.Marshal(event)
	if err != nil {
		log.Printf("[WebSocket] Error marshaling event: %v", err)
		return err
	}

	h.broadcast <- message
	log.Printf("[WebSocket] Broadcasted event: %s to %d clients", eventType, len(h.clients))
	return nil
}

// BroadcastToUser broadcasts an event to a specific user
func (h *Hub) BroadcastToUser(userID string, eventType string, data interface{}) error {
	event := Event{
		Type: eventType,
		Data: data,
	}

	message, err := json.Marshal(event)
	if err != nil {
		log.Printf("[WebSocket] Error marshaling event: %v", err)
		return err
	}

	h.mu.RLock()
	defer h.mu.RUnlock()

	count := 0
	for client := range h.clients {
		if client.UserID == userID {
			select {
			case client.Send <- message:
				count++
			default:
				// Client's send channel is full, skip
			}
		}
	}

	log.Printf("[WebSocket] Broadcasted event: %s to user %s (%d clients)", eventType, userID, count)
	return nil
}

// GetClientCount returns the number of connected clients
func (h *Hub) GetClientCount() int {
	h.mu.RLock()
	defer h.mu.RUnlock()
	return len(h.clients)
}
