package handover

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

// Enums
type HandoverStatus string
type HandoverType string
type DiscrepancyType string
type ResponsibilityType string

const (
	StatusPending     HandoverStatus = "PENDING"
	StatusConfirmed   HandoverStatus = "CONFIRMED"
	StatusRejected    HandoverStatus = "REJECTED"
	StatusDiscrepancy HandoverStatus = "DISCREPANCY" // Có chênh lệch cần xử lý

	TypePartial  HandoverType = "PARTIAL"
	TypeFull     HandoverType = "FULL"
	TypeEndShift HandoverType = "END_SHIFT"

	DiscrepancyShortage DiscrepancyType = "SHORTAGE" // Thiếu tiền
	DiscrepancyOverage  DiscrepancyType = "OVERAGE"  // Thừa tiền

	ResponsibilityWaiter   ResponsibilityType = "WAITER"
	ResponsibilityCashier  ResponsibilityType = "CASHIER"
	ResponsibilitySystem   ResponsibilityType = "SYSTEM"
	ResponsibilityCustomer ResponsibilityType = "CUSTOMER"
	ResponsibilityUnknown  ResponsibilityType = "UNKNOWN"
)

// CashHandover represents a cash handover from waiter to cashier
type CashHandover struct {
	ID             primitive.ObjectID `bson:"_id,omitempty" json:"id"`
	WaiterShiftID  primitive.ObjectID `bson:"waiter_shift_id" json:"waiter_shift_id"`
	CashierShiftID primitive.ObjectID `bson:"cashier_shift_id" json:"cashier_shift_id"`
	WaiterID       primitive.ObjectID `bson:"waiter_id" json:"waiter_id"`
	WaiterName     string             `bson:"waiter_name" json:"waiter_name"`
	CashierID      primitive.ObjectID `bson:"cashier_id" json:"cashier_id"`
	CashierName    string             `bson:"cashier_name" json:"cashier_name"`

	// Amounts
	DeclaredAmount float64 `bson:"declared_amount" json:"declared_amount"` // Waiter khai báo
	ActualAmount   float64 `bson:"actual_amount" json:"actual_amount"`     // Cashier thực nhận
	Discrepancy    float64 `bson:"discrepancy" json:"discrepancy"`         // Chênh lệch

	HandoverType HandoverType   `bson:"handover_type" json:"handover_type"`
	Status       HandoverStatus `bson:"status" json:"status"`

	// Notes and reasons
	WaiterNote                string             `bson:"waiter_note,omitempty" json:"waiter_note,omitempty"`
	CashierNote               string             `bson:"cashier_note,omitempty" json:"cashier_note,omitempty"`
	DiscrepancyReason         string             `bson:"discrepancy_reason,omitempty" json:"discrepancy_reason,omitempty"`
	DiscrepancyResponsibility ResponsibilityType `bson:"discrepancy_responsibility,omitempty" json:"discrepancy_responsibility,omitempty"`

	// Timestamps
	HandoverAt   time.Time  `bson:"handover_at" json:"handover_at"`
	ConfirmedAt  *time.Time `bson:"confirmed_at,omitempty" json:"confirmed_at,omitempty"`
	ReconciledAt *time.Time `bson:"reconciled_at,omitempty" json:"reconciled_at,omitempty"`

	// Metadata
	EndCash          float64            `bson:"end_cash,omitempty" json:"end_cash,omitempty"`
	RequiresApproval bool               `bson:"requires_approval" json:"requires_approval"`
	ApprovedBy       primitive.ObjectID `bson:"approved_by,omitempty" json:"approved_by,omitempty"`
	ApprovedAt       *time.Time         `bson:"approved_at,omitempty" json:"approved_at,omitempty"`

	CreatedAt time.Time `bson:"created_at" json:"created_at"`
	UpdatedAt time.Time `bson:"updated_at" json:"updated_at"`
}

// HasDiscrepancy checks if handover has discrepancy
func (h *CashHandover) HasDiscrepancy() bool {
	return h.Discrepancy != 0
}

// GetDiscrepancyType returns the type of discrepancy
func (h *CashHandover) GetDiscrepancyType() DiscrepancyType {
	if h.Discrepancy < 0 {
		return DiscrepancyShortage
	} else if h.Discrepancy > 0 {
		return DiscrepancyOverage
	}
	return ""
}

// RequiresManagerApproval checks if requires manager approval (large discrepancy)
func (h *CashHandover) RequiresManagerApproval(threshold float64) bool {
	return h.HasDiscrepancy() && (h.Discrepancy > threshold || h.Discrepancy < -threshold)
}

// Request structures
type CreateHandoverRequest struct {
	DeclaredAmount float64      `json:"declared_amount" binding:"required,gt=0"`
	HandoverType   HandoverType `json:"handover_type" binding:"required"`
	WaiterNote     string       `json:"waiter_note"`
}

type CreateHandoverAndEndShiftRequest struct {
	DeclaredAmount float64 `json:"declared_amount" binding:"required,gt=0"`
	WaiterNote     string  `json:"waiter_note"`
	EndCash        float64 `json:"end_cash" binding:"min=0"`
}

type ConfirmHandoverRequest struct {
	ActualAmount              float64            `json:"actual_amount" binding:"gte=0"`
	Status                    HandoverStatus     `json:"status" binding:"required"`
	CashierNote               string             `json:"cashier_note"`
	DiscrepancyReason         string             `json:"discrepancy_reason"`
	DiscrepancyResponsibility ResponsibilityType `json:"discrepancy_responsibility"`
}
