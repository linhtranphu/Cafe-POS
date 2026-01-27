package cashier

import (
	"time"
)

type CashReconciliation struct {
	ID               string    `json:"id" bson:"_id,omitempty"`
	ShiftID          string    `json:"shift_id" bson:"shift_id"`
	CashierID        string    `json:"cashier_id" bson:"cashier_id"`
	ExpectedCash     float64   `json:"expected_cash" bson:"expected_cash"`
	ActualCash       float64   `json:"actual_cash" bson:"actual_cash"`
	Difference       float64   `json:"difference" bson:"difference"`
	Status           string    `json:"status" bson:"status"` // MATCH, OVER, SHORT
	Notes            string    `json:"notes" bson:"notes"`
	ReconciliationAt time.Time `json:"reconciliation_at" bson:"reconciliation_at"`
	CreatedAt        time.Time `json:"created_at" bson:"created_at"`
	UpdatedAt        time.Time `json:"updated_at" bson:"updated_at"`
}

type PaymentDiscrepancy struct {
	ID         string     `json:"id" bson:"_id,omitempty"`
	OrderID    string     `json:"order_id" bson:"order_id"`
	CashierID  string     `json:"cashier_id" bson:"cashier_id"`
	Reason     string     `json:"reason" bson:"reason"`
	Amount     float64    `json:"amount" bson:"amount"`
	Status     string     `json:"status" bson:"status"` // PENDING, RESOLVED
	ReportedAt time.Time  `json:"reported_at" bson:"reported_at"`
	ResolvedAt *time.Time `json:"resolved_at,omitempty" bson:"resolved_at,omitempty"`
	CreatedAt  time.Time  `json:"created_at" bson:"created_at"`
	UpdatedAt  time.Time  `json:"updated_at" bson:"updated_at"`
}

// Constants
const (
	ReconciliationStatusMatch = "MATCH"
	ReconciliationStatusOver  = "OVER"
	ReconciliationStatusShort = "SHORT"

	DiscrepancyStatusPending  = "PENDING"
	DiscrepancyStatusResolved = "RESOLVED"
)

// Methods
func (cr *CashReconciliation) CalculateDifference() {
	cr.Difference = cr.ActualCash - cr.ExpectedCash
	if cr.Difference == 0 {
		cr.Status = ReconciliationStatusMatch
	} else if cr.Difference > 0 {
		cr.Status = ReconciliationStatusOver
	} else {
		cr.Status = ReconciliationStatusShort
	}
}

func (pd *PaymentDiscrepancy) Resolve() {
	pd.Status = DiscrepancyStatusResolved
	now := time.Now()
	pd.ResolvedAt = &now
	pd.UpdatedAt = now
}