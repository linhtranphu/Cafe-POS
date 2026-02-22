package websocket

import (
	"log"
	"net/http"

	socketio "github.com/googollee/go-socket.io"
	"github.com/googollee/go-socket.io/engineio"
	"github.com/googollee/go-socket.io/engineio/transport"
	"github.com/googollee/go-socket.io/engineio/transport/polling"
	"github.com/googollee/go-socket.io/engineio/transport/websocket"
)

// NewSocketIOServer creates a new Socket.IO server with standard library
func NewSocketIOServer() (*socketio.Server, error) {
	// Allow all origins for development
	// In production, configure this properly
	allowOriginFunc := func(r *http.Request) bool {
		return true
	}

	server := socketio.NewServer(&engineio.Options{
		Transports: []transport.Transport{
			&polling.Transport{
				CheckOrigin: allowOriginFunc,
			},
			&websocket.Transport{
				CheckOrigin: allowOriginFunc,
			},
		},
	})

	// Handle connection
	server.OnConnect("/", func(s socketio.Conn) error {
		s.SetContext("")
		log.Printf("[Socket.IO] Client connected: %s", s.ID())
		return nil
	})

	// Handle disconnection
	server.OnDisconnect("/", func(s socketio.Conn, reason string) {
		log.Printf("[Socket.IO] Client disconnected: %s, reason: %s", s.ID(), reason)
	})

	// Handle errors
	server.OnError("/", func(s socketio.Conn, e error) {
		log.Printf("[Socket.IO] Error from client %s: %v", s.ID(), e)
	})

	return server, nil
}

// BroadcastPrintJobCreated broadcasts print job created event to all clients
func BroadcastPrintJobCreatedToSocketIO(server *socketio.Server, jobData map[string]interface{}) {
	if server == nil {
		log.Printf("[Socket.IO] Server not initialized, skipping broadcast")
		return
	}

	// Broadcast to all clients in default namespace
	server.BroadcastToNamespace("/", "print-job-created", jobData)
	log.Printf("[Socket.IO] Broadcasted print-job-created event")
}

// BroadcastPrintJobStatusChanged broadcasts print job status change
func BroadcastPrintJobStatusChangedToSocketIO(server *socketio.Server, jobID string, status string, errorMsg string) {
	if server == nil {
		log.Printf("[Socket.IO] Server not initialized, skipping broadcast")
		return
	}

	data := map[string]interface{}{
		"job_id":    jobID,
		"status":    status,
		"error_msg": errorMsg,
	}

	server.BroadcastToNamespace("/", "print-job-status-changed", data)
	log.Printf("[Socket.IO] Broadcasted print-job-status-changed event for job %s", jobID)
}
