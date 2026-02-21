package printing

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

// NotificationType represents the type of notification
type NotificationType string

const (
	NotificationTypePrintFailure    NotificationType = "PRINT_FAILURE"
	NotificationTypePrinterOffline  NotificationType = "PRINTER_OFFLINE"
	NotificationTypeHardwareError   NotificationType = "HARDWARE_ERROR"
	NotificationTypePaperOut        NotificationType = "PAPER_OUT"
	NotificationTypePaperJam        NotificationType = "PAPER_JAM"
	NotificationTypeCoverOpen       NotificationType = "COVER_OPEN"
)

// NotificationSeverity represents the severity level
type NotificationSeverity string

const (
	NotificationSeverityInfo    NotificationSeverity = "INFO"
	NotificationSeverityWarning NotificationSeverity = "WARNING"
	NotificationSeverityError   NotificationSeverity = "ERROR"
)

// PrintNotification represents a print-related notification
type PrintNotification struct {
	ID          primitive.ObjectID   `bson:"_id,omitempty" json:"id"`
	Type        NotificationType     `bson:"type" json:"type"`
	Severity    NotificationSeverity `bson:"severity" json:"severity"`
	Title       string               `bson:"title" json:"title"`
	Message     string               `bson:"message" json:"message"`
	JobID       primitive.ObjectID   `bson:"job_id,omitempty" json:"job_id,omitempty"`
	OrderID     primitive.ObjectID   `bson:"order_id,omitempty" json:"order_id,omitempty"`
	OrderNumber string               `bson:"order_number,omitempty" json:"order_number,omitempty"`
	PrinterID   primitive.ObjectID   `bson:"printer_id,omitempty" json:"printer_id,omitempty"`
	PrinterName string               `bson:"printer_name,omitempty" json:"printer_name,omitempty"`
	ErrorMsg    string               `bson:"error_msg,omitempty" json:"error_msg,omitempty"`
	Read        bool                 `bson:"read" json:"read"`
	CreatedAt   time.Time            `bson:"created_at" json:"created_at"`
}

// PrintNotificationRepository defines the interface for notification storage
type PrintNotificationRepository interface {
	Create(notification *PrintNotification) error
	FindUnread() ([]*PrintNotification, error)
	FindByJobID(jobID primitive.ObjectID) ([]*PrintNotification, error)
	MarkAsRead(id primitive.ObjectID) error
	MarkAllAsRead() error
	DeleteOld(olderThan time.Time) error
}
