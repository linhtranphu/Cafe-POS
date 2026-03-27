package printing

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

// PrinterType represents the type of printer
type PrinterType string

const (
	PrinterTypeBill  PrinterType = "BILL"
	PrinterTypeLabel PrinterType = "LABEL"
)

// ConnectionType represents the connection type of a printer
type ConnectionType string

const (
	ConnectionTypeNetwork ConnectionType = "NETWORK"
	ConnectionTypeUSB     ConnectionType = "USB"
)

// PrinterConfig represents a printer configuration in the system
type PrinterConfig struct {
	ID             primitive.ObjectID `bson:"_id,omitempty" json:"id"`
	Name           string             `bson:"name" json:"name"`
	Type           PrinterType        `bson:"type" json:"type"`
	ConnectionType ConnectionType     `bson:"connection_type" json:"connection_type"`
	IPAddress      string             `bson:"ip_address,omitempty" json:"ip_address,omitempty"`
	Port           int                `bson:"port,omitempty" json:"port,omitempty"`
	USBPath        string             `bson:"usb_path,omitempty" json:"usb_path,omitempty"`
	PaperWidth     int                `bson:"paper_width" json:"paper_width"` // mm: 58 or 80
	IsDefault      bool               `bson:"is_default" json:"is_default"`
	IsEnabled      bool               `bson:"is_enabled" json:"is_enabled"`
	CreatedAt      time.Time          `bson:"created_at" json:"created_at"`
	UpdatedAt      time.Time          `bson:"updated_at" json:"updated_at"`
}

// CreatePrinterConfigRequest represents the request to create a new printer configuration
type CreatePrinterConfigRequest struct {
	Name           string         `json:"name" binding:"required"`
	Type           PrinterType    `json:"type" binding:"required,oneof=BILL LABEL"`
	ConnectionType ConnectionType `json:"connection_type" binding:"required,oneof=NETWORK USB"`
	IPAddress      string         `json:"ip_address"`
	Port           int            `json:"port"`
	USBPath        string         `json:"usb_path"`
	PaperWidth     int            `json:"paper_width" binding:"required,oneof=40 58 80"`
	IsDefault      bool           `json:"is_default"`
	IsEnabled      bool           `json:"is_enabled"`
}

// UpdatePrinterConfigRequest represents the request to update a printer configuration
type UpdatePrinterConfigRequest struct {
	Name           string          `json:"name"`
	Type           PrinterType     `json:"type" binding:"omitempty,oneof=BILL LABEL"`
	ConnectionType ConnectionType  `json:"connection_type" binding:"omitempty,oneof=NETWORK USB"`
	IPAddress      string          `json:"ip_address"`
	Port           *int            `json:"port"`
	USBPath        string          `json:"usb_path"`
	PaperWidth     *int            `json:"paper_width" binding:"omitempty,oneof=40 58 80"`
	IsDefault      *bool           `json:"is_default"`
	IsEnabled      *bool           `json:"is_enabled"`
}

// PrinterConfigFilter represents filter options for querying printer configurations
type PrinterConfigFilter struct {
	Type      PrinterType
	IsEnabled *bool
	Page      int
	Limit     int
}
