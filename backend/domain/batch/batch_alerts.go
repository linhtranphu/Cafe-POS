package batch

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

// BatchAlerts aggregates all types of batch alerts
type BatchAlerts struct {
	LowStock    []*LowStockAlert  `json:"low_stock"`
	Expiring    []*ExpiringAlert  `json:"expiring"`
	Expired     []*ExpiredAlert   `json:"expired"`
	LastChecked time.Time         `json:"last_checked"`
}

// LowStockAlert represents an alert when batch stock is below threshold
type LowStockAlert struct {
	BatchDefinitionID primitive.ObjectID `json:"batch_definition_id"`
	BatchName         string             `json:"batch_name"`
	CurrentStock      float64            `json:"current_stock"`
	Threshold         float64            `json:"threshold"`
	Unit              string             `json:"unit"`
}

// ExpiringAlert represents an alert when a batch is about to expire
type ExpiringAlert struct {
	BatchRecordID     primitive.ObjectID `json:"batch_record_id"`
	BatchName         string             `json:"batch_name"`
	QuantityRemaining float64            `json:"quantity_remaining"`
	Unit              string             `json:"unit"`
	ExpiresAt         time.Time          `json:"expires_at"`
	HoursUntilExpiry  int                `json:"hours_until_expiry"`
}

// ExpiredAlert represents an alert when a batch has expired
type ExpiredAlert struct {
	BatchRecordID  primitive.ObjectID `json:"batch_record_id"`
	BatchName      string             `json:"batch_name"`
	QuantityWasted float64            `json:"quantity_wasted"`
	Unit           string             `json:"unit"`
	CostWasted     float64            `json:"cost_wasted"`
	ExpiredAt      time.Time          `json:"expired_at"`
}

// NewBatchAlerts creates a new BatchAlerts aggregate
func NewBatchAlerts() *BatchAlerts {
	return &BatchAlerts{
		LowStock:    make([]*LowStockAlert, 0),
		Expiring:    make([]*ExpiringAlert, 0),
		Expired:     make([]*ExpiredAlert, 0),
		LastChecked: time.Now(),
	}
}

// AddLowStockAlert adds a low stock alert to the aggregate
func (ba *BatchAlerts) AddLowStockAlert(alert *LowStockAlert) {
	ba.LowStock = append(ba.LowStock, alert)
}

// AddExpiringAlert adds an expiring alert to the aggregate
func (ba *BatchAlerts) AddExpiringAlert(alert *ExpiringAlert) {
	ba.Expiring = append(ba.Expiring, alert)
}

// AddExpiredAlert adds an expired alert to the aggregate
func (ba *BatchAlerts) AddExpiredAlert(alert *ExpiredAlert) {
	ba.Expired = append(ba.Expired, alert)
}

// HasAlerts returns true if there are any alerts
func (ba *BatchAlerts) HasAlerts() bool {
	return len(ba.LowStock) > 0 || len(ba.Expiring) > 0 || len(ba.Expired) > 0
}

// TotalAlertCount returns the total number of alerts
func (ba *BatchAlerts) TotalAlertCount() int {
	return len(ba.LowStock) + len(ba.Expiring) + len(ba.Expired)
}

// NewLowStockAlert creates a new LowStockAlert
func NewLowStockAlert(batchDefID primitive.ObjectID, batchName string, currentStock, threshold float64, unit string) *LowStockAlert {
	return &LowStockAlert{
		BatchDefinitionID: batchDefID,
		BatchName:         batchName,
		CurrentStock:      currentStock,
		Threshold:         threshold,
		Unit:              unit,
	}
}

// IsLowStock checks if the current stock is below or equal to threshold
func (lsa *LowStockAlert) IsLowStock() bool {
	return lsa.CurrentStock <= lsa.Threshold
}

// NewExpiringAlert creates a new ExpiringAlert
func NewExpiringAlert(batchRecordID primitive.ObjectID, batchName string, quantityRemaining float64, unit string, expiresAt time.Time) *ExpiringAlert {
	hoursUntilExpiry := int(time.Until(expiresAt).Hours())
	return &ExpiringAlert{
		BatchRecordID:     batchRecordID,
		BatchName:         batchName,
		QuantityRemaining: quantityRemaining,
		Unit:              unit,
		ExpiresAt:         expiresAt,
		HoursUntilExpiry:  hoursUntilExpiry,
	}
}

// IsExpiringSoon checks if the batch is expiring within the warning hours
func (ea *ExpiringAlert) IsExpiringSoon(warningHours int) bool {
	return ea.HoursUntilExpiry <= warningHours && ea.HoursUntilExpiry > 0
}

// NewExpiredAlert creates a new ExpiredAlert
func NewExpiredAlert(batchRecordID primitive.ObjectID, batchName string, quantityWasted float64, unit string, costWasted float64, expiredAt time.Time) *ExpiredAlert {
	return &ExpiredAlert{
		BatchRecordID:  batchRecordID,
		BatchName:      batchName,
		QuantityWasted: quantityWasted,
		Unit:           unit,
		CostWasted:     costWasted,
		ExpiredAt:      expiredAt,
	}
}

// CalculateCostWasted calculates the cost wasted from an expired batch
func (ea *ExpiredAlert) CalculateCostWasted(costPerUnit float64) float64 {
	return ea.QuantityWasted * costPerUnit
}
