package menu

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

// MenuCategory represents a menu category
type MenuCategory struct {
	ID        primitive.ObjectID `bson:"_id,omitempty" json:"id"`
	Name      string             `bson:"name" json:"name"`
	CreatedAt time.Time          `bson:"created_at" json:"created_at"`
	UpdatedAt time.Time          `bson:"updated_at" json:"updated_at"`
}

// CreateMenuCategoryRequest represents the request to create a menu category
type CreateMenuCategoryRequest struct {
	Name string `json:"name" binding:"required"`
}

// UpdateMenuCategoryRequest represents the request to update a menu category
type UpdateMenuCategoryRequest struct {
	Name string `json:"name" binding:"required"`
}
