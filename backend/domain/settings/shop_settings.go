package settings

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

// ShopSettings represents the shop configuration and settings
type ShopSettings struct {
	ID                 primitive.ObjectID `bson:"_id,omitempty" json:"id"`
	ShopName           string             `bson:"shop_name" json:"shop_name"`
	ShopAddress        string             `bson:"shop_address,omitempty" json:"shop_address,omitempty"`
	ShopPhone          string             `bson:"shop_phone,omitempty" json:"shop_phone,omitempty"`
	LogoURL            string             `bson:"logo_url,omitempty" json:"logo_url,omitempty"`
	CustomMessage      string             `bson:"custom_message,omitempty" json:"custom_message,omitempty"`
	PrintBridgeURL     string             `bson:"print_bridge_url,omitempty" json:"print_bridge_url,omitempty"` // Local print bridge URL
	ShowLogo           bool               `bson:"show_logo" json:"show_logo"`                                   // Show logo on bill
	ShowAddress        bool               `bson:"show_address" json:"show_address"`                             // Show address on bill
	ShowPhone          bool               `bson:"show_phone" json:"show_phone"`                                 // Show phone on bill
	ShowCustomMessage  bool               `bson:"show_custom_message" json:"show_custom_message"`               // Show custom message on bill
	LowMarginThreshold float64            `bson:"low_margin_threshold" json:"low_margin_threshold"`             // Default: 20.0
	AutoPrintEnabled   bool               `bson:"auto_print_enabled" json:"auto_print_enabled"`                 // Default: true
	PaperWidth         int                `bson:"paper_width" json:"paper_width"`                               // Bill paper width in mm (40, 58, 80)
	
	// Label Printer Settings (replaces temp bill printer)
	LabelPrinterEnabled bool   `bson:"label_printer_enabled" json:"label_printer_enabled"`
	LabelPrinterIP      string `bson:"label_printer_ip" json:"label_printer_ip"`
	LabelPrinterPort    int    `bson:"label_printer_port" json:"label_printer_port"`
	LabelWidth          int    `bson:"label_width" json:"label_width"`   // mm
	LabelHeight         int    `bson:"label_height" json:"label_height"` // mm
	
	CreatedAt          time.Time          `bson:"created_at" json:"created_at"`
	UpdatedAt          time.Time          `bson:"updated_at" json:"updated_at"`
}

// NewShopSettings creates a new ShopSettings with default values
func NewShopSettings(shopName string) *ShopSettings {
	now := time.Now()
	return &ShopSettings{
		ShopName:           shopName,
		ShowLogo:           true,
		ShowAddress:        true,
		ShowPhone:          true,
		ShowCustomMessage:  true,
		LowMarginThreshold: 20.0, // Default threshold
		AutoPrintEnabled:   true, // Default: auto-print enabled
		PaperWidth:         58,   // Default: 58mm thermal paper
		CreatedAt:          now,
		UpdatedAt:          now,
	}
}

// UpdateLowMarginThreshold updates the low margin threshold
func (s *ShopSettings) UpdateLowMarginThreshold(threshold float64) error {
	if threshold < 0 {
		return ErrInvalidThreshold
	}
	s.LowMarginThreshold = threshold
	s.UpdatedAt = time.Now()
	return nil
}

// SetAutoPrintEnabled updates the auto-print enabled setting
func (s *ShopSettings) SetAutoPrintEnabled(enabled bool) {
	s.AutoPrintEnabled = enabled
	s.UpdatedAt = time.Now()
}

// UpdatePrintSettings updates print-related settings
func (s *ShopSettings) UpdatePrintSettings(address, phone, logoURL, customMessage string) {
	s.ShopAddress = address
	s.ShopPhone = phone
	s.LogoURL = logoURL
	s.CustomMessage = customMessage
	s.UpdatedAt = time.Now()
}

// SetFieldVisibility updates field visibility settings
func (s *ShopSettings) SetFieldVisibility(showLogo, showAddress, showPhone, showCustomMessage bool) {
	s.ShowLogo = showLogo
	s.ShowAddress = showAddress
	s.ShowPhone = showPhone
	s.ShowCustomMessage = showCustomMessage
	s.UpdatedAt = time.Now()
}
