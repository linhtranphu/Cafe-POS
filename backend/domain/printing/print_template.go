package printing

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

// TemplateType represents the type of print template
type TemplateType string

const (
	TemplateTypeBill     TemplateType = "BILL"
	TemplateTypeLabel    TemplateType = "LABEL"
	TemplateTypeTempBill TemplateType = "TEMP_BILL" // Bill tạm
)

// PrintTemplate represents a print template in the system
type PrintTemplate struct {
	ID        primitive.ObjectID `bson:"_id,omitempty" json:"id"`
	Type      TemplateType       `bson:"type" json:"type"`
	Name      string             `bson:"name" json:"name"`
	Content   string             `bson:"content" json:"content"` // Template string
	IsDefault bool               `bson:"is_default" json:"is_default"`
	CreatedAt time.Time          `bson:"created_at" json:"created_at"`
	UpdatedAt time.Time          `bson:"updated_at" json:"updated_at"`
}

// TemplateData represents the data passed to templates for rendering
type TemplateData struct {
	ShopName          string
	ShopAddress       string
	ShopPhone         string
	LogoURL           string
	CustomMessage     string
	ShowLogo          bool
	ShowAddress       bool
	ShowPhone         bool
	ShowCustomMessage bool
	Order             interface{} // Will be *order.Order, using interface{} to avoid circular dependency
	PrintTime         time.Time
	// For labels
	ItemIndex  int
	TotalItems int
}

// CreatePrintTemplateRequest represents the request to create a new print template
type CreatePrintTemplateRequest struct {
	Type      TemplateType `json:"type" binding:"required,oneof=BILL LABEL TEMP_BILL"`
	Name      string       `json:"name" binding:"required"`
	Content   string       `json:"content" binding:"required"`
	IsDefault bool         `json:"is_default"`
}

// UpdatePrintTemplateRequest represents the request to update a print template
type UpdatePrintTemplateRequest struct {
	Type      TemplateType `json:"type" binding:"omitempty,oneof=BILL LABEL TEMP_BILL"`
	Name      string       `json:"name"`
	Content   string       `json:"content"`
	IsDefault *bool        `json:"is_default"`
}

// PrintTemplateFilter represents filter options for querying print templates
type PrintTemplateFilter struct {
	Type  TemplateType
	Page  int
	Limit int
}
