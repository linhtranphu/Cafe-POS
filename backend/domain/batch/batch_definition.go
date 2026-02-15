package batch

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

// BatchDefinition defines how a batch is prepared, including the recipe and conversion rates
type BatchDefinition struct {
	ID                 primitive.ObjectID `bson:"_id,omitempty" json:"id"`
	Name               string             `bson:"name" json:"name"`
	Unit               string             `bson:"unit" json:"unit"`
	ShelfLifeHours     int                `bson:"shelf_life_hours" json:"shelf_life_hours"`
	ConversionRates    []ConversionRate   `bson:"conversion_rates" json:"conversion_rates"`
	LowStockThreshold  float64            `bson:"low_stock_threshold" json:"low_stock_threshold"`
	ExpiryWarningHours int                `bson:"expiry_warning_hours" json:"expiry_warning_hours"`
	CreatedAt          time.Time          `bson:"created_at" json:"created_at"`
	UpdatedAt          time.Time          `bson:"updated_at" json:"updated_at"`
}

// ConversionRate is a value object that defines the conversion from source ingredient to batch
type ConversionRate struct {
	SourceIngredientID   primitive.ObjectID `bson:"source_ingredient_id" json:"source_ingredient_id"`
	SourceIngredientName string             `bson:"source_ingredient_name" json:"source_ingredient_name"`
	SourceQuantity       float64            `bson:"source_quantity" json:"source_quantity"`
	SourceUnit           string             `bson:"source_unit" json:"source_unit"`
	BatchQuantity        float64            `bson:"batch_quantity" json:"batch_quantity"`
	WastageRate          float64            `bson:"wastage_rate" json:"wastage_rate"` // 0.0 to 1.0 (e.g., 0.1 = 10% wastage)
}

// CreateBatchDefinitionRequest represents the request to create a new batch definition
type CreateBatchDefinitionRequest struct {
	Name               string             `json:"name" binding:"required"`
	Unit               string             `json:"unit" binding:"required"`
	ShelfLifeHours     int                `json:"shelf_life_hours" binding:"required,min=1"`
	ConversionRates    []ConversionRate   `json:"conversion_rates" binding:"required,min=1,dive"`
	LowStockThreshold  float64            `json:"low_stock_threshold" binding:"min=0"`
	ExpiryWarningHours int                `json:"expiry_warning_hours" binding:"min=0"`
}

// UpdateBatchDefinitionRequest represents the request to update a batch definition
type UpdateBatchDefinitionRequest struct {
	Name               string           `json:"name"`
	Unit               string           `json:"unit"`
	ShelfLifeHours     *int             `json:"shelf_life_hours" binding:"omitempty,min=1"`
	ConversionRates    []ConversionRate `json:"conversion_rates" binding:"omitempty,min=1,dive"`
	LowStockThreshold  *float64         `json:"low_stock_threshold" binding:"omitempty,min=0"`
	ExpiryWarningHours *int             `json:"expiry_warning_hours" binding:"omitempty,min=0"`
}

// BatchDefinitionFilter represents filter options for querying batch definitions
type BatchDefinitionFilter struct {
	Search string
	Page   int
	Limit  int
}
