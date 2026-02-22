package websocket

import (
	"cafe-pos/backend/domain/printing"
	"context"
	"log"
	"time"
	
	socketio "github.com/googollee/go-socket.io"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

// PrinterConfigRepository interface for getting printer details
type PrinterConfigRepository interface {
	FindByID(ctx context.Context, id primitive.ObjectID) (*printing.PrinterConfig, error)
}

// Broadcaster provides methods to broadcast events via WebSocket
type Broadcaster struct {
	hub          *Hub
	printerRepo  PrinterConfigRepository
	socketIOServer *socketio.Server
}

// NewBroadcaster creates a new broadcaster
func NewBroadcaster(hub *Hub) *Broadcaster {
	return &Broadcaster{
		hub:            hub,
		printerRepo:    nil,
		socketIOServer: nil,
	}
}

// SetPrinterRepository sets the printer config repository
func (b *Broadcaster) SetPrinterRepository(repo PrinterConfigRepository) {
	b.printerRepo = repo
}

// SetSocketIOServer sets the Socket.IO server
func (b *Broadcaster) SetSocketIOServer(server *socketio.Server) {
	b.socketIOServer = server
}

// BroadcastPrintJobCreated broadcasts when a new print job is created
func (b *Broadcaster) BroadcastPrintJobCreated(job *printing.PrintJob) {
	// Get printer details
	var printerIP string
	var printerPort int
	if b.printerRepo != nil {
		ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
		defer cancel()
		
		printer, err := b.printerRepo.FindByID(ctx, job.PrinterID)
		if err == nil && printer != nil {
			printerIP = printer.IPAddress
			printerPort = printer.Port
		} else {
			log.Printf("[WebSocket] Failed to get printer details: %v", err)
		}
	}

	jobData := map[string]interface{}{
		"job": map[string]interface{}{
			"id":           job.ID.Hex(),
			"type":         string(job.Type),
			"order_id":     job.OrderID.Hex(),
			"order_number": job.OrderNumber,
			"printer_id":   job.PrinterID.Hex(),
			"printer_ip":   printerIP,
			"printer_port": printerPort,
			"content":      job.Content,
			"status":       string(job.Status),
			"created_at":   job.CreatedAt,
		},
	}

	// Broadcast using Socket.IO server
	if b.socketIOServer != nil {
		BroadcastPrintJobCreatedToSocketIO(b.socketIOServer, jobData)
	}

	// Also broadcast using old hub for backward compatibility
	if b.hub != nil {
		if err := b.hub.BroadcastEvent("print-job-created", jobData); err != nil {
			log.Printf("[WebSocket] Failed to broadcast print-job-created via hub: %v", err)
		}
	}
}

// BroadcastPrintJobStatusChanged broadcasts when a print job status changes
func (b *Broadcaster) BroadcastPrintJobStatusChanged(jobID string, status printing.PrintJobStatus, errorMsg string) {
	// Broadcast using Socket.IO server
	if b.socketIOServer != nil {
		BroadcastPrintJobStatusChangedToSocketIO(b.socketIOServer, jobID, string(status), errorMsg)
	}

	// Also broadcast using old hub for backward compatibility
	if b.hub != nil {
		data := map[string]interface{}{
			"job_id":    jobID,
			"status":    string(status),
			"error_msg": errorMsg,
		}

		if err := b.hub.BroadcastEvent("print-job-status-changed", data); err != nil {
			log.Printf("[WebSocket] Failed to broadcast print-job-status-changed via hub: %v", err)
		}
	}
}

// BroadcastPrintJobFailed broadcasts when a print job fails
func (b *Broadcaster) BroadcastPrintJobFailed(jobID string, errorMsg string) {
	// Broadcast using Socket.IO server
	if b.socketIOServer != nil {
		data := map[string]interface{}{
			"job_id":    jobID,
			"error_msg": errorMsg,
		}
		b.socketIOServer.BroadcastToNamespace("/", "print-job-failed", data)
	}

	// Also broadcast using old hub for backward compatibility
	if b.hub != nil {
		data := map[string]interface{}{
			"job_id":    jobID,
			"error_msg": errorMsg,
		}

		if err := b.hub.BroadcastEvent("print-job-failed", data); err != nil {
			log.Printf("[WebSocket] Failed to broadcast print-job-failed via hub: %v", err)
		}
	}
}

// BroadcastPrinterOffline broadcasts when a printer goes offline
func (b *Broadcaster) BroadcastPrinterOffline(printerID string, printerName string) {
	if b.hub == nil {
		log.Printf("[WebSocket] Hub not initialized, skipping broadcast")
		return
	}

	data := map[string]interface{}{
		"printer_id":   printerID,
		"printer_name": printerName,
	}

	if err := b.hub.BroadcastEvent("printer-offline", data); err != nil {
		log.Printf("[WebSocket] Failed to broadcast printer-offline: %v", err)
	}
}

// BroadcastPrinterOnline broadcasts when a printer comes online
func (b *Broadcaster) BroadcastPrinterOnline(printerID string, printerName string) {
	if b.hub == nil {
		log.Printf("[WebSocket] Hub not initialized, skipping broadcast")
		return
	}

	data := map[string]interface{}{
		"printer_id":   printerID,
		"printer_name": printerName,
	}

	if err := b.hub.BroadcastEvent("printer-online", data); err != nil {
		log.Printf("[WebSocket] Failed to broadcast printer-online: %v", err)
	}
}

// BroadcastPrinterError broadcasts when a printer encounters an error
func (b *Broadcaster) BroadcastPrinterError(printerID string, printerName string, errorMsg string) {
	if b.hub == nil {
		log.Printf("[WebSocket] Hub not initialized, skipping broadcast")
		return
	}

	data := map[string]interface{}{
		"printer_id":   printerID,
		"printer_name": printerName,
		"error_msg":    errorMsg,
	}

	if err := b.hub.BroadcastEvent("printer-error", data); err != nil {
		log.Printf("[WebSocket] Failed to broadcast printer-error: %v", err)
	}
}
