package cashier

import (
	"time"
)

type PaymentAudit struct {
	ID        string    `json:"id" bson:"_id,omitempty"`
	OrderID   string    `json:"order_id" bson:"order_id"`
	Action    string    `json:"action" bson:"action"` // CANCEL, REFUND, OVERRIDE, LOCK
	CashierID string    `json:"cashier_id" bson:"cashier_id"`
	Reason    string    `json:"reason" bson:"reason"`
	OldStatus string    `json:"old_status" bson:"old_status"`
	NewStatus string    `json:"new_status" bson:"new_status"`
	Amount    float64   `json:"amount" bson:"amount"`
	AuditedAt time.Time `json:"audited_at" bson:"audited_at"`
	CreatedAt time.Time `json:"created_at" bson:"created_at"`
}

// Constants
const (
	AuditActionCancel   = "CANCEL"
	AuditActionRefund   = "REFUND"
	AuditActionOverride = "OVERRIDE"
	AuditActionLock     = "LOCK"
)

// Methods
func NewPaymentAudit(orderID, action, cashierID, reason, oldStatus, newStatus string, amount float64) *PaymentAudit {
	now := time.Now()
	return &PaymentAudit{
		OrderID:   orderID,
		Action:    action,
		CashierID: cashierID,
		Reason:    reason,
		OldStatus: oldStatus,
		NewStatus: newStatus,
		Amount:    amount,
		AuditedAt: now,
		CreatedAt: now,
	}
}