package websocket

import (
	"encoding/json"
	"log"
	"net/http"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/gorilla/websocket"
)

var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
	CheckOrigin: func(r *http.Request) bool {
		// Allow all origins for development
		// In production, you should check the origin
		return true
	},
}

// SocketIOMessage represents a Socket.IO protocol message
type SocketIOMessage struct {
	Type int         `json:"type"` // 0=CONNECT, 1=DISCONNECT, 2=EVENT, 3=ACK, 4=ERROR
	NSP  string      `json:"nsp"`  // Namespace
	Data interface{} `json:"data"` // Event data
}

// HandleSocketIO handles Socket.IO compatible WebSocket connections
// This handles both HTTP polling (initial handshake) and WebSocket upgrade
func HandleSocketIO(hub *Hub) gin.HandlerFunc {
	return func(c *gin.Context) {
		// Check if this is a WebSocket upgrade request
		if c.GetHeader("Upgrade") == "websocket" {
			handleWebSocketUpgrade(c, hub)
			return
		}

		// Otherwise, handle as HTTP polling (Socket.IO v4 handshake)
		handleHTTPPolling(c, hub)
	}
}

// handleHTTPPolling handles Socket.IO HTTP polling requests (handshake)
func handleHTTPPolling(c *gin.Context, hub *Hub) {
	// Get EIO (Engine.IO) parameter
	eio := c.Query("EIO")
	transport := c.Query("transport")

	// Socket.IO v4 uses EIO=4
	if eio != "4" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Unsupported Engine.IO version"})
		return
	}

	// For polling transport, send handshake response
	if transport == "polling" {
		// Generate session ID
		sid := uuid.New().String()

		// Send handshake packet (Engine.IO OPEN packet)
		handshake := map[string]interface{}{
			"sid":          sid,
			"upgrades":     []string{"websocket"},
			"pingInterval": 25000,
			"pingTimeout":  60000,
			"maxPayload":   1000000,
		}

		handshakeData, _ := json.Marshal(handshake)
		
		// Engine.IO v4 packet format: "0" + JSON
		response := "0" + string(handshakeData)

		c.Header("Content-Type", "text/plain; charset=UTF-8")
		c.String(http.StatusOK, response)
		return
	}

	c.JSON(http.StatusBadRequest, gin.H{"error": "Unsupported transport"})
}

// handleWebSocketUpgrade handles WebSocket upgrade after HTTP polling handshake
func handleWebSocketUpgrade(c *gin.Context, hub *Hub) {
	// Get token from query params (Socket.IO sends it here)
	token := c.Query("token")
	if token == "" {
		// Try to get from auth header
		authHeader := c.GetHeader("Authorization")
		if strings.HasPrefix(authHeader, "Bearer ") {
			token = strings.TrimPrefix(authHeader, "Bearer ")
		}
	}

	// For now, we'll accept connections without strict auth validation
	// In production, you should validate the JWT token here
	userID := "anonymous"
	if token != "" {
		// TODO: Validate JWT and extract user ID
		// For now, just use a placeholder
		userID = "user-" + token[:min(8, len(token))]
	}

	// Upgrade HTTP connection to WebSocket
	conn, err := upgrader.Upgrade(c.Writer, c.Request, nil)
	if err != nil {
		log.Printf("[WebSocket] Upgrade error: %v", err)
		return
	}

	// Create client
	client := &Client{
		ID:     uuid.New().String(),
		UserID: userID,
		Send:   make(chan []byte, 256),
		hub:    hub,
	}

	// Register client
	hub.register <- client

	// Start goroutines for reading and writing
	go client.writePump(conn)
	go client.readPump(conn)

	log.Printf("[WebSocket] New connection established: %s", client.ID)
}

func min(a, b int) int {
	if a < b {
		return a
	}
	return b
}

// readPump pumps messages from the WebSocket connection to the hub
func (c *Client) readPump(conn *websocket.Conn) {
	defer func() {
		c.hub.unregister <- c
		conn.Close()
	}()

	conn.SetReadDeadline(time.Now().Add(60 * time.Second))
	conn.SetPongHandler(func(string) error {
		conn.SetReadDeadline(time.Now().Add(60 * time.Second))
		return nil
	})

	for {
		_, message, err := conn.ReadMessage()
		if err != nil {
			if websocket.IsUnexpectedCloseError(err, websocket.CloseGoingAway, websocket.CloseAbnormalClosure) {
				log.Printf("[WebSocket] Read error: %v", err)
			}
			break
		}

		// Log received message for debugging
		log.Printf("[WebSocket] Received message from %s: %s", c.ID, string(message))

		// Handle Socket.IO protocol messages
		// For now, we just acknowledge receipt
		// In a full implementation, you would parse and handle different message types
	}
}

// writePump pumps messages from the hub to the WebSocket connection
func (c *Client) writePump(conn *websocket.Conn) {
	ticker := time.NewTicker(54 * time.Second)
	defer func() {
		ticker.Stop()
		conn.Close()
	}()

	// Send Socket.IO v4 CONNECT packet immediately after WebSocket upgrade
	// Format: "40" for default namespace, or "40{json}" with connection data
	if err := conn.WriteMessage(websocket.TextMessage, []byte("40")); err != nil {
		log.Printf("[WebSocket] Error sending Socket.IO CONNECT: %v", err)
		return
	}
	log.Printf("[WebSocket] Sent Socket.IO CONNECT packet to client %s", c.ID)

	for {
		select {
		case message, ok := <-c.Send:
			conn.SetWriteDeadline(time.Now().Add(10 * time.Second))
			if !ok {
				// Hub closed the channel
				conn.WriteMessage(websocket.CloseMessage, []byte{})
				return
			}

			// Parse the event
			var event Event
			if err := json.Unmarshal(message, &event); err != nil {
				log.Printf("[WebSocket] Error unmarshaling event: %v", err)
				continue
			}

			// Format as Socket.IO v4 event message
			// Socket.IO protocol: "42" + JSON array [eventName, data]
			socketIOEvent := []interface{}{event.Type, event.Data}
			eventData, err := json.Marshal(socketIOEvent)
			if err != nil {
				log.Printf("[WebSocket] Error marshaling Socket.IO event: %v", err)
				continue
			}

			// Send with Socket.IO protocol prefix
			socketIOMessage := append([]byte("42"), eventData...)
			if err := conn.WriteMessage(websocket.TextMessage, socketIOMessage); err != nil {
				log.Printf("[WebSocket] Write error: %v", err)
				return
			}
			log.Printf("[WebSocket] Sent event '%s' to client %s", event.Type, c.ID)

		case <-ticker.C:
			conn.SetWriteDeadline(time.Now().Add(10 * time.Second))
			// Send ping (Engine.IO v4 protocol: "2" for PING)
			if err := conn.WriteMessage(websocket.TextMessage, []byte("2")); err != nil {
				return
			}
		}
	}
}
