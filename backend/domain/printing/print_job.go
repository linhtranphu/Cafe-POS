package printing

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

// PrintJobType represents the type of print job
type PrintJobType string

const (
	PrintJobTypeBill  PrintJobType = "BILL"
	PrintJobTypeLabel PrintJobType = "LABEL"
)

// PrintJobStatus represents the status of a print job
type PrintJobStatus string

const (
	PrintJobStatusPending   PrintJobStatus = "PENDING"
	PrintJobStatusPrinting  PrintJobStatus = "PRINTING"
	PrintJobStatusCompleted PrintJobStatus = "COMPLETED"
	PrintJobStatusFailed    PrintJobStatus = "FAILED"
)

// PrintJob represents a print job in the system
type PrintJob struct {
	ID          primitive.ObjectID `bson:"_id,omitempty" json:"id"`
	Type        PrintJobType       `bson:"type" json:"type"`
	OrderID     primitive.ObjectID `bson:"order_id" json:"order_id"`
	OrderNumber string             `bson:"order_number" json:"order_number"`
	PrinterID   primitive.ObjectID `bson:"printer_id" json:"printer_id"`
	Content     string             `bson:"content" json:"content"` // Rendered content (text or base64 encoded binary)
	ContentType string             `bson:"content_type,omitempty" json:"content_type,omitempty"` // "text" or "binary" (default: "text")
	Status      PrintJobStatus     `bson:"status" json:"status"`
	RetryCount  int                `bson:"retry_count" json:"retry_count"`
	MaxRetries  int                `bson:"max_retries" json:"max_retries"`
	ErrorMsg    string             `bson:"error_msg,omitempty" json:"error_msg,omitempty"`
	CreatedAt   time.Time          `bson:"created_at" json:"created_at"`
	UpdatedAt   time.Time          `bson:"updated_at" json:"updated_at"`
	PrintedAt   *time.Time         `bson:"printed_at,omitempty" json:"printed_at,omitempty"`
}

// CreatePrintJobRequest represents the request to create a new print job
type CreatePrintJobRequest struct {
	Type        PrintJobType       `json:"type" binding:"required,oneof=BILL LABEL"`
	OrderID     primitive.ObjectID `json:"order_id" binding:"required"`
	OrderNumber string             `json:"order_number" binding:"required"`
	PrinterID   primitive.ObjectID `json:"printer_id" binding:"required"`
	Content     string             `json:"content" binding:"required"`
}

// PrintJobFilter represents filter options for querying print jobs
type PrintJobFilter struct {
	OrderID primitive.ObjectID
	Status  PrintJobStatus
	Type    PrintJobType
	Page    int
	Limit   int
}
