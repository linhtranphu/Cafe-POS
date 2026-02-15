package batch

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

// BatchUsageLog records when a batch is used in an order
type BatchUsageLog struct {
	ID           primitive.ObjectID `bson:"_id,omitempty" json:"id"`
	BatchRecordID primitive.ObjectID `bson:"batch_record_id" json:"batch_record_id"`
	BatchName    string             `bson:"batch_name" json:"batch_name"`
	OrderID      primitive.ObjectID `bson:"order_id" json:"order_id"`
	MenuItemID   primitive.ObjectID `bson:"menu_item_id" json:"menu_item_id"`
	MenuItemName string             `bson:"menu_item_name" json:"menu_item_name"`
	QuantityUsed float64            `bson:"quantity_used" json:"quantity_used"`
	Unit         string             `bson:"unit" json:"unit"`
	CostPerUnit  float64            `bson:"cost_per_unit" json:"cost_per_unit"`
	TotalCost    float64            `bson:"total_cost" json:"total_cost"`
	UsedAt       time.Time          `bson:"used_at" json:"used_at"`
}

// CreateBatchUsageLogRequest represents the request to create a new batch usage log
type CreateBatchUsageLogRequest struct {
	BatchRecordID primitive.ObjectID `json:"batch_record_id" binding:"required"`
	BatchName     string             `json:"batch_name" binding:"required"`
	OrderID       primitive.ObjectID `json:"order_id" binding:"required"`
	MenuItemID    primitive.ObjectID `json:"menu_item_id" binding:"required"`
	MenuItemName  string             `json:"menu_item_name" binding:"required"`
	QuantityUsed  float64            `json:"quantity_used" binding:"required,gt=0"`
	Unit          string             `json:"unit" binding:"required"`
	CostPerUnit   float64            `json:"cost_per_unit" binding:"required,gte=0"`
}

// BatchUsageLogFilter represents filter options for querying batch usage logs
type BatchUsageLogFilter struct {
	BatchRecordID *primitive.ObjectID
	OrderID       *primitive.ObjectID
	MenuItemID    *primitive.ObjectID
	FromDate      *time.Time
	ToDate        *time.Time
	Page          int
	Limit         int
}

// CalculateTotalCost calculates the total cost of the batch usage
func (bul *BatchUsageLog) CalculateTotalCost() float64 {
	return bul.QuantityUsed * bul.CostPerUnit
}

// NewBatchUsageLog creates a new BatchUsageLog from a request
func NewBatchUsageLog(req CreateBatchUsageLogRequest) *BatchUsageLog {
	now := time.Now()
	log := &BatchUsageLog{
		BatchRecordID: req.BatchRecordID,
		BatchName:     req.BatchName,
		OrderID:       req.OrderID,
		MenuItemID:    req.MenuItemID,
		MenuItemName:  req.MenuItemName,
		QuantityUsed:  req.QuantityUsed,
		Unit:          req.Unit,
		CostPerUnit:   req.CostPerUnit,
		UsedAt:        now,
	}
	log.TotalCost = log.CalculateTotalCost()
	return log
}
