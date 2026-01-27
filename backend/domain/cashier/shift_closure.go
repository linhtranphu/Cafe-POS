package cashier

import (
	"time"
)

type ShiftClosure struct {
	ID              string    `json:"id" bson:"_id,omitempty"`
	ShiftID         string    `json:"shift_id" bson:"shift_id"`
	CashierID       string    `json:"cashier_id" bson:"cashier_id"`
	TotalOrders     int       `json:"total_orders" bson:"total_orders"`
	TotalRevenue    float64   `json:"total_revenue" bson:"total_revenue"`
	CashRevenue     float64   `json:"cash_revenue" bson:"cash_revenue"`
	TransferRevenue float64   `json:"transfer_revenue" bson:"transfer_revenue"`
	QRRevenue       float64   `json:"qr_revenue" bson:"qr_revenue"`
	Discrepancies   []string  `json:"discrepancies" bson:"discrepancies"`
	Status          string    `json:"status" bson:"status"` // PENDING, COMPLETED
	ClosedAt        time.Time `json:"closed_at" bson:"closed_at"`
	CreatedAt       time.Time `json:"created_at" bson:"created_at"`
	UpdatedAt       time.Time `json:"updated_at" bson:"updated_at"`
}

// Constants
const (
	ClosureStatusPending   = "PENDING"
	ClosureStatusCompleted = "COMPLETED"
)

// Methods
func (sc *ShiftClosure) AddDiscrepancy(discrepancy string) {
	sc.Discrepancies = append(sc.Discrepancies, discrepancy)
}

func (sc *ShiftClosure) HasDiscrepancies() bool {
	return len(sc.Discrepancies) > 0
}

func (sc *ShiftClosure) Complete() {
	sc.Status = ClosureStatusCompleted
	sc.UpdatedAt = time.Now()
}