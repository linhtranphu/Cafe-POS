package order

import (
	"time"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

// CostStatus represents the status of cost calculation for an order item
type CostStatus string

const (
	CostStatusFinal      CostStatus = "FINAL"      // Cost has been calculated and finalized (shift closed)
	CostStatusEstimated  CostStatus = "ESTIMATED"  // Cost is estimated (shift not closed)
	CostStatusIncomplete CostStatus = "INCOMPLETE" // Missing ingredient cost data
)

// OrderItemWithCost represents an order item with accounting cost tracking
// This is stored in a separate collection for efficient querying and aggregation
type OrderItemWithCost struct {
	ID               primitive.ObjectID `bson:"_id,omitempty" json:"id"`
	OrderID          primitive.ObjectID `bson:"order_id" json:"order_id"`
	MenuItemID       primitive.ObjectID `bson:"menu_item_id" json:"menu_item_id"`
	Name             string             `bson:"name" json:"name"`
	Price            float64            `bson:"price" json:"price"`
	Quantity         int                `bson:"quantity" json:"quantity"`
	Note             string             `bson:"note,omitempty" json:"note,omitempty"`
	Subtotal         float64            `bson:"subtotal" json:"subtotal"`
	
	// Accounting cost fields
	AccountingCost   float64            `bson:"accounting_cost" json:"accounting_cost"`
	CostCalculatedAt time.Time          `bson:"cost_calculated_at" json:"cost_calculated_at"`
	CostStatus       CostStatus         `bson:"cost_status" json:"cost_status"`
	
	CreatedAt        time.Time          `bson:"created_at" json:"created_at"`
}
