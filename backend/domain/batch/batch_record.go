package batch

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

// BatchRecord represents a specific instance of a batch that has been prepared
type BatchRecord struct {
	ID                primitive.ObjectID `bson:"_id,omitempty" json:"id"`
	BatchDefinitionID primitive.ObjectID `bson:"batch_definition_id" json:"batch_definition_id"`
	BatchName         string             `bson:"batch_name" json:"batch_name"`
	QuantityProduced  float64            `bson:"quantity_produced" json:"quantity_produced"`
	QuantityRemaining float64            `bson:"quantity_remaining" json:"quantity_remaining"`
	Unit              string             `bson:"unit" json:"unit"`
	CostPerUnit       float64            `bson:"cost_per_unit" json:"cost_per_unit"`
	TotalCost         float64            `bson:"total_cost" json:"total_cost"`
	PreparedBy        string             `bson:"prepared_by" json:"prepared_by"`
	PreparedByName    string             `bson:"prepared_by_name" json:"prepared_by_name"` // Username of the person who prepared
	PreparedAt        time.Time          `bson:"prepared_at" json:"prepared_at"`
	ExpiresAt         time.Time          `bson:"expires_at" json:"expires_at"`
	Status            string             `bson:"status" json:"status"` // "available", "expired", "depleted"
	IngredientsUsed   []IngredientUsage  `bson:"ingredients_used" json:"ingredients_used"`
	CreatedAt         time.Time          `bson:"created_at" json:"created_at"`
	UpdatedAt         time.Time          `bson:"updated_at" json:"updated_at"`
}

// IngredientUsage is a value object that records the ingredients used to prepare a batch
type IngredientUsage struct {
	IngredientID   primitive.ObjectID `bson:"ingredient_id" json:"ingredient_id"`
	IngredientName string             `bson:"ingredient_name" json:"ingredient_name"`
	Quantity       float64            `bson:"quantity" json:"quantity"`
	Unit           string             `bson:"unit" json:"unit"`
	CostPerUnit    float64            `bson:"cost_per_unit" json:"cost_per_unit"`
	TotalCost      float64            `bson:"total_cost" json:"total_cost"`
}

// Batch status constants
const (
	BatchStatusAvailable = "available"
	BatchStatusExpired   = "expired"
	BatchStatusDepleted  = "depleted"
)

// CreateBatchRecordRequest represents the request to create a new batch record
type CreateBatchRecordRequest struct {
	BatchDefinitionID primitive.ObjectID `json:"batch_definition_id" binding:"required"`
	QuantityProduced  float64            `json:"quantity_produced" binding:"required,gt=0"`
	PreparedBy        string             `json:"prepared_by" binding:"required"`
}

// UpdateBatchQuantityRequest represents the request to update batch quantity
type UpdateBatchQuantityRequest struct {
	QuantityRemaining float64 `json:"quantity_remaining" binding:"required,gte=0"`
}

// BatchRecordFilter represents filter options for querying batch records
type BatchRecordFilter struct {
	BatchDefinitionID *primitive.ObjectID
	Status            string
	PreparedBy        string
	FromDate          *time.Time
	ToDate            *time.Time
	Page              int
	Limit             int
}

// IsExpired checks if the batch record has expired
func (br *BatchRecord) IsExpired() bool {
	return time.Now().After(br.ExpiresAt)
}

// IsDepleted checks if the batch record has been fully used
func (br *BatchRecord) IsDepleted() bool {
	return br.QuantityRemaining <= 0
}

// IsAvailable checks if the batch record is available for use
func (br *BatchRecord) IsAvailable() bool {
	return br.Status == BatchStatusAvailable && !br.IsExpired() && !br.IsDepleted()
}

// CalculateExpiryTime calculates the expiry time based on prepared time and shelf life
func (br *BatchRecord) CalculateExpiryTime(shelfLifeHours int) time.Time {
	return br.PreparedAt.Add(time.Duration(shelfLifeHours) * time.Hour)
}

// UpdateStatus updates the batch status based on current state
func (br *BatchRecord) UpdateStatus() {
	if br.IsDepleted() {
		br.Status = BatchStatusDepleted
	} else if br.IsExpired() {
		br.Status = BatchStatusExpired
	} else {
		br.Status = BatchStatusAvailable
	}
}

// DeductQuantity deducts the specified quantity from the batch
func (br *BatchRecord) DeductQuantity(quantity float64) error {
	if quantity < 0 {
		return ErrInvalidQuantity
	}
	if quantity > br.QuantityRemaining {
		return ErrInsufficientQuantity
	}
	br.QuantityRemaining -= quantity
	br.UpdatedAt = time.Now()
	br.UpdateStatus()
	return nil
}

// CalculateTotalCost calculates the total cost from ingredients used
func (br *BatchRecord) CalculateTotalCost() float64 {
	total := 0.0
	for _, ingredient := range br.IngredientsUsed {
		total += ingredient.TotalCost
	}
	return total
}

// CalculateCostPerUnit calculates the cost per unit of the batch
func (br *BatchRecord) CalculateCostPerUnit() float64 {
	if br.QuantityProduced <= 0 {
		return 0
	}
	return br.TotalCost / br.QuantityProduced
}

// Custom errors
var (
	ErrInvalidQuantity      = &BatchError{Code: "INVALID_QUANTITY", Message: "Quantity must be greater than zero"}
	ErrInsufficientQuantity = &BatchError{Code: "INSUFFICIENT_QUANTITY", Message: "Insufficient quantity available"}
)

// BatchError represents a batch-related error
type BatchError struct {
	Code    string
	Message string
}

func (e *BatchError) Error() string {
	return e.Message
}
