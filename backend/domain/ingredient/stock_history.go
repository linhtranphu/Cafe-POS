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

// IngredientStockSummary is the aggregated stock movement for one ingredient in a date range
type IngredientStockSummary struct {
	IngredientID primitive.ObjectID `bson:"_id" json:"ingredient_id"`
	Consumed     float64            `bson:"consumed" json:"consumed"`       // deducted by orders
	Restocked    float64            `bson:"restocked" json:"restocked"`     // added by purchases
	Wasted       float64            `bson:"wasted" json:"wasted"`           // removed as waste
	Adjusted     float64            `bson:"adjusted" json:"adjusted"`       // net manual adjustments
	OrderCount   int                `bson:"order_count" json:"order_count"` // number of order transactions
}

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
	CostPerUnit  float64            `bson:"cost_per_unit" json:"cost_per_unit"`     // Price at time of transaction
	TotalCost    float64            `bson:"total_cost" json:"total_cost"`           // Total cost for this transaction
	CreatedAt    time.Time          `bson:"created_at" json:"created_at"`
}