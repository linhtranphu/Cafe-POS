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
	PaperWidth         int                `bson:"paper_width" json:"paper_width"`                               // 58 or 80 (mm)
	LabelSize          string             `bson:"label_size" json:"label_size"`                                 // "40x30", "50x30", "60x40" (mm)
	ShowLogo           bool               `bson:"show_logo" json:"show_logo"`                                   // Show logo on bill
	ShowAddress        bool               `bson:"show_address" json:"show_address"`                             // Show address on bill
	ShowPhone          bool               `bson:"show_phone" json:"show_phone"`                                 // Show phone on bill
	ShowCustomMessage  bool               `bson:"show_custom_message" json:"show_custom_message"`               // Show custom message on bill
	LowMarginThreshold float64            `bson:"low_margin_threshold" json:"low_margin_threshold"`             // Default: 20.0
	AutoPrintEnabled   bool               `bson:"auto_print_enabled" json:"auto_print_enabled"`                 // Default: true
	CreatedAt          time.Time          `bson:"created_at" json:"created_at"`
	UpdatedAt          time.Time          `bson:"updated_at" json:"updated_at"`
}

// NewShopSettings creates a new ShopSettings with default values
func NewShopSettings(shopName string) *ShopSettings {
	now := time.Now()
	return &ShopSettings{
		ShopName:           shopName,
		PaperWidth:         80,      // Default: 80mm
		LabelSize:          "60x40", // Default: 60x40mm
		ShowLogo:           true,
		ShowAddress:        true,
		ShowPhone:          true,
		ShowCustomMessage:  true,
		LowMarginThreshold: 20.0, // Default threshold
		AutoPrintEnabled:   true, // Default: auto-print enabled
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
func (s *ShopSettings) UpdatePrintSettings(address, phone, logoURL, customMessage string, paperWidth int, labelSize string) error {
	if paperWidth != 58 && paperWidth != 80 {
		return ErrInvalidPaperWidth
	}
	
	validLabelSizes := map[string]bool{
		"40x30": true,
		"50x30": true,
		"60x40": true,
	}
	if !validLabelSizes[labelSize] {
		return ErrInvalidLabelSize
	}
	
	s.ShopAddress = address
	s.ShopPhone = phone
	s.LogoURL = logoURL
	s.CustomMessage = customMessage
	s.PaperWidth = paperWidth
	s.LabelSize = labelSize
	s.UpdatedAt = time.Now()
	return nil
}

// SetFieldVisibility updates field visibility settings
func (s *ShopSettings) SetFieldVisibility(showLogo, showAddress, showPhone, showCustomMessage bool) {
	s.ShowLogo = showLogo
	s.ShowAddress = showAddress
	s.ShowPhone = showPhone
	s.ShowCustomMessage = showCustomMessage
	s.UpdatedAt = time.Now()
}
