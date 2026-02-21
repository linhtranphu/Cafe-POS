package websocket

import (
	"cafe-pos/backend/domain/printing"
	"log"
)

// Broadcaster provides methods to broadcast events via WebSocket
type Broadcaster struct {
	hub *Hub
}

// NewBroadcaster creates a new broadcaster
func NewBroadcaster(hub *Hub) *Broadcaster {
	return &Broadcaster{hub: hub}
}

// BroadcastPrintJobCreated broadcasts when a new print job is created
func (b *Broadcaster) BroadcastPrintJobCreated(job *printing.PrintJob) {
	if b.hub == nil {
		log.Printf("[WebSocket] Hub not initialized, skipping broadcast")
		return
	}

	data := map[string]interface{}{
		"job": map[string]interface{}{
			"id":           job.ID.Hex(),
			"type":         string(job.Type),
			"order_id":     job.OrderID.Hex(),
			"order_number": job.OrderNumber,
			"printer_id":   job.PrinterID.Hex(),
			"content":      job.Content,
			"status":       string(job.Status),
			"created_at":   job.CreatedAt,
		},
	}

	if err := b.hub.BroadcastEvent("print-job-created", data); err != nil {
		log.Printf("[WebSocket] Failed to broadcast print-job-created: %v", err)
	}
}

// BroadcastPrintJobStatusChanged broadcasts when a print job status changes
func (b *Broadcaster) BroadcastPrintJobStatusChanged(jobID string, status printing.PrintJobStatus, errorMsg string) {
	if b.hub == nil {
		log.Printf("[WebSocket] Hub not initialized, skipping broadcast")
		return
	}

	data := map[string]interface{}{
		"job_id":    jobID,
		"status":    string(status),
		"error_msg": errorMsg,
	}

	if err := b.hub.BroadcastEvent("print-job-status-changed", data); err != nil {
		log.Printf("[WebSocket] Failed to broadcast print-job-status-changed: %v", err)
	}
}

// BroadcastPrintJobFailed broadcasts when a print job fails
func (b *Broadcaster) BroadcastPrintJobFailed(jobID string, errorMsg string) {
	if b.hub == nil {
		log.Printf("[WebSocket] Hub not initialized, skipping broadcast")
		return
	}

	data := map[string]interface{}{
		"job_id":    jobID,
		"error_msg": errorMsg,
	}

	if err := b.hub.BroadcastEvent("print-job-failed", data); err != nil {
		log.Printf("[WebSocket] Failed to broadcast print-job-failed: %v", err)
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
