package ingredient

import (
	"time"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type TransactionType string

const (
	TransactionAdjustment TransactionType = "adjustment"
	TransactionOrder      TransactionType = "order"
	TransactionPurchase   TransactionType = "purchase"
	TransactionWaste      TransactionType = "waste"
)

type StockHistory struct {
	ID           primitive.ObjectID `bson:"_id,omitempty" json:"id"`
	IngredientID primitive.ObjectID `bson:"ingredient_id" json:"ingredient_id"`
	Type         TransactionType    `bson:"type" json:"type"`
	Quantity     float64            `bson:"quantity" json:"quantity"`
	BeforeQty    float64            `bson:"before_qty" json:"before_qty"`
	AfterQty     float64            `bson:"after_qty" json:"after_qty"`
	Reason       string             `bson:"reason" json:"reason"`
	OrderID      *primitive.ObjectID `bson:"order_id,omitempty" json:"order_id,omitempty"`
	UserID       primitive.ObjectID `bson:"user_id" json:"user_id"`
	Username     string             `bson:"username" json:"username"`
	CreatedAt    time.Time          `bson:"created_at" json:"created_at"`
}