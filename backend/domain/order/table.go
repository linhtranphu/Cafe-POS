package order

import (
	"time"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type TableStatus string

const (
	TableEmpty    TableStatus = "EMPTY"
	TableOccupied TableStatus = "OCCUPIED"
	TableReserved TableStatus = "RESERVED"
)

type Table struct {
	ID        primitive.ObjectID `bson:"_id,omitempty" json:"id"`
	Name      string             `bson:"name" json:"name"`
	Capacity  int                `bson:"capacity" json:"capacity"`
	Status    TableStatus        `bson:"status" json:"status"`
	Area      string             `bson:"area,omitempty" json:"area,omitempty"`
	CreatedAt time.Time          `bson:"created_at" json:"created_at"`
	UpdatedAt time.Time          `bson:"updated_at" json:"updated_at"`
}

type CreateTableRequest struct {
	Name     string `json:"name" binding:"required"`
	Capacity int    `json:"capacity" binding:"required,min=1"`
	Area     string `json:"area"`
}

type UpdateTableRequest struct {
	Name     string      `json:"name"`
	Capacity int         `json:"capacity" binding:"omitempty,min=1"`
	Status   TableStatus `json:"status"`
	Area     string      `json:"area"`
}
