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
	ID          primitive.ObjectID `bson:"_id,omitempty" json:"id"`
	Name        string             `bson:"name" json:"name"`
	Category    string             `bson:"category" json:"category"`
	Unit        UnitType           `bson:"unit" json:"unit"`
	Quantity    float64            `bson:"quantity" json:"quantity"`
	MinStock    float64            `bson:"min_stock" json:"min_stock"`
	CostPerUnit float64            `bson:"cost_per_unit" json:"cost_per_unit"`
	Supplier    string             `bson:"supplier" json:"supplier"`
	CreatedAt   time.Time          `bson:"created_at" json:"created_at"`
	UpdatedAt   time.Time          `bson:"updated_at" json:"updated_at"`
}

type CreateIngredientRequest struct {
	Name        string   `json:"name" binding:"required"`
	Category    string   `json:"category" binding:"required"`
	Unit        UnitType `json:"unit" binding:"required"`
	Quantity    float64  `json:"quantity" binding:"required,min=0"`
	MinStock    float64  `json:"min_stock" binding:"min=0"`
	CostPerUnit float64  `json:"cost_per_unit" binding:"min=0"`
	Supplier    string   `json:"supplier"`
}

type UpdateIngredientRequest struct {
	Name        string   `json:"name"`
	Category    string   `json:"category"`
	Unit        UnitType `json:"unit"`
	Quantity    *float64 `json:"quantity" binding:"omitempty,min=0"`
	MinStock    *float64 `json:"min_stock" binding:"omitempty,min=0"`
	CostPerUnit *float64 `json:"cost_per_unit" binding:"omitempty,min=0"`
	Supplier    string   `json:"supplier"`
}

type StockAdjustmentRequest struct {
	Quantity float64 `json:"quantity" binding:"required"`
	Reason   string  `json:"reason" binding:"required"`
	UserID   string  `json:"user_id"`
	Username string  `json:"username"`
}