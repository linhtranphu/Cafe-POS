package ingredient

import (
	"time"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type UnitType string

const (
	// Mass units (kg base)
	UnitKilogram UnitType = "kg"
	UnitGram     UnitType = "g"
	
	// Volume units (L base)
	UnitLiter      UnitType = "L"
	UnitMilliliter UnitType = "ml"
	
	// Count units
	UnitPiece UnitType = "piece"
	UnitBox   UnitType = "box"
	UnitPack  UnitType = "pack"
)

type Ingredient struct {
	ID                primitive.ObjectID `bson:"_id,omitempty" json:"id"`
	Name              string             `bson:"name" json:"name"`
	Category          string             `bson:"category" json:"category"`
	Unit              UnitType           `bson:"unit" json:"unit"`
	Quantity          float64            `bson:"quantity" json:"quantity"`
	MinStock          float64            `bson:"min_stock" json:"min_stock"`
	CostPerUnit       float64            `bson:"cost_per_unit" json:"cost_per_unit"`
	Supplier          string             `bson:"supplier" json:"supplier"`
	ConversionRate    float64            `bson:"conversion_rate" json:"conversion_rate"`       // Default: 1.0 - for unit conversion
	WastagePercentage float64            `bson:"wastage_percentage" json:"wastage_percentage"` // Default: 0.0 - wastage factor
	CreatedAt         time.Time          `bson:"created_at" json:"created_at"`
	UpdatedAt         time.Time          `bson:"updated_at" json:"updated_at"`
}

type CreateIngredientRequest struct {
	Name              string    `json:"name" binding:"required"`
	Category          string    `json:"category" binding:"required"`
	Unit              UnitType  `json:"unit" binding:"required"`
	Quantity          float64   `json:"quantity" binding:"min=0"`
	MinStock          float64   `json:"min_stock" binding:"min=0"`
	CostPerUnit       float64   `json:"cost_per_unit" binding:"min=0"`
	Supplier          string    `json:"supplier"`
	ConversionRate    *float64  `json:"conversion_rate" binding:"omitempty,min=0"`       // Optional: default 1.0
	WastagePercentage *float64  `json:"wastage_percentage" binding:"omitempty,min=0"`    // Optional: default 0.0
	CreatedDate       *string   `json:"created_date"` // Optional: custom creation date (ISO 8601 format)
}

type UpdateIngredientRequest struct {
	Name              string   `json:"name"`
	Category          string   `json:"category"`
	Unit              UnitType `json:"unit"`
	Quantity          *float64 `json:"quantity" binding:"omitempty,min=0"`
	MinStock          *float64 `json:"min_stock" binding:"omitempty,min=0"`
	CostPerUnit       *float64 `json:"cost_per_unit" binding:"omitempty,min=0"`
	Supplier          string   `json:"supplier"`
	ConversionRate    *float64 `json:"conversion_rate" binding:"omitempty,min=0"`
	WastagePercentage *float64 `json:"wastage_percentage" binding:"omitempty,min=0"`
}

// Stock Operation Types - Must match frontend constants
type StockOperationType string

const (
	StockOperationIn     StockOperationType = "in"     // Stock IN (purchase)
	StockOperationOut    StockOperationType = "out"    // Stock OUT (usage/waste)
	StockOperationAdjust StockOperationType = "adjust" // Stock ADJUST (inventory correction)
)

// StockInRequest - For purchasing/receiving stock
type StockInRequest struct {
	Quantity    float64 `json:"quantity" binding:"required,gt=0"`
	CostPerUnit float64 `json:"cost_per_unit" binding:"min=0"` // Optional: if 0, use current price
	Reason      string  `json:"reason"`
	UserID      string  `json:"user_id"`
	Username    string  `json:"username"`
}

// StockOutRequest - For using/removing stock
type StockOutRequest struct {
	Quantity float64 `json:"quantity" binding:"required,gt=0"`
	Reason   string  `json:"reason" binding:"required"`
	UserID   string  `json:"user_id"`
	Username string  `json:"username"`
}

// StockAdjustRequest - For inventory corrections (set to specific quantity)
type StockAdjustRequest struct {
	NewQuantity float64 `json:"new_quantity" binding:"required,min=0"` // Target quantity
	CostPerUnit float64 `json:"cost_per_unit" binding:"min=0"`         // Optional: if increase due to purchase
	Reason      string  `json:"reason" binding:"required"`
	UserID      string  `json:"user_id"`
	Username    string  `json:"username"`
}

// Legacy request for backward compatibility
type StockAdjustmentRequest struct {
	Quantity    float64 `json:"quantity" binding:"required"`
	Reason      string  `json:"reason" binding:"required"`
	CostPerUnit float64 `json:"cost_per_unit"` // Optional: price for this purchase (if different from current)
	UserID      string  `json:"user_id"`
	Username    string  `json:"username"`
}

// IngredientCategory represents a category for ingredients
type IngredientCategory struct {
	ID        primitive.ObjectID `bson:"_id,omitempty" json:"id"`
	Name      string             `bson:"name" json:"name"`
	CreatedAt time.Time          `bson:"created_at" json:"created_at"`
	UpdatedAt time.Time          `bson:"updated_at" json:"updated_at"`
}

type CreateCategoryRequest struct {
	Name string `json:"name" binding:"required"`
}
