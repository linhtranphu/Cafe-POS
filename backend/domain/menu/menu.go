package menu

import (
	"time"
	"cafe-pos/backend/domain/ingredient"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

// CostStatus represents the status of cost calculation for a menu item
type CostStatus string

const (
	CostStatusFinal      CostStatus = "FINAL"      // Cost has been calculated and finalized
	CostStatusEstimated  CostStatus = "ESTIMATED"  // Cost is estimated (shift not closed)
	CostStatusIncomplete CostStatus = "INCOMPLETE" // Missing ingredient cost data
)

type Ingredient struct {
	Name     string               `bson:"name" json:"name"`
	Quantity float64              `bson:"quantity" json:"quantity"`
	Unit     ingredient.UnitType  `bson:"unit" json:"unit"`
}

type MenuItem struct {
	ID          primitive.ObjectID `bson:"_id,omitempty" json:"id"`
	Name        string             `bson:"name" json:"name"`
	Price       float64            `bson:"price" json:"price"`
	Category    string             `bson:"category" json:"category"`
	Description string             `bson:"description" json:"description"`
	Ingredients []Ingredient       `bson:"ingredients" json:"ingredients"`
	Available   bool               `bson:"available" json:"available"`
	
	// Cost tracking fields
	CurrentCost          float64    `bson:"current_cost" json:"current_cost"`
	CostLastCalculatedAt time.Time  `bson:"cost_last_calculated_at" json:"cost_last_calculated_at"`
	CostStatus           CostStatus `bson:"cost_status" json:"cost_status"`
	
	CreatedAt   time.Time `bson:"created_at" json:"created_at"`
	UpdatedAt   time.Time `bson:"updated_at" json:"updated_at"`
}

type CreateMenuItemRequest struct {
	Name        string       `json:"name" binding:"required"`
	Price       float64      `json:"price" binding:"required,min=0"`
	Category    string       `json:"category" binding:"required"`
	Description string       `json:"description"`
	Ingredients []Ingredient `json:"ingredients"`
}

type UpdateMenuItemRequest struct {
	Name        string       `json:"name"`
	Price       float64      `json:"price" binding:"min=0"`
	Category    string       `json:"category"`
	Description string       `json:"description"`
	Ingredients []Ingredient `json:"ingredients"`
	Available   *bool        `json:"available"`
}