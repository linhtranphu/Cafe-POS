package menu

import (
	"time"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type Ingredient struct {
	Name     string  `bson:"name" json:"name"`
	Quantity float64 `bson:"quantity" json:"quantity"`
	Unit     string  `bson:"unit" json:"unit"`
}

type MenuItem struct {
	ID          primitive.ObjectID `bson:"_id,omitempty" json:"id"`
	Name        string             `bson:"name" json:"name"`
	Price       float64            `bson:"price" json:"price"`
	Category    string             `bson:"category" json:"category"`
	Description string             `bson:"description" json:"description"`
	Ingredients []Ingredient       `bson:"ingredients" json:"ingredients"`
	Available   bool               `bson:"available" json:"available"`
	CreatedAt   time.Time          `bson:"created_at" json:"created_at"`
	UpdatedAt   time.Time          `bson:"updated_at" json:"updated_at"`
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