package handover

import (
	"errors"
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

// HandoverStatus represents the current status of a cash handover
type HandoverStatus string

const (
	HandoverStatusPending     HandoverStatus = "PENDING"
	HandoverStatusConfirmed   HandoverStatus = "CONFIRMED"
	HandoverStatusRejected    HandoverStatus = "REJECTED"
	HandoverStatusDiscrepancy HandoverStatus = "DISCREPANCY"
)

// HandoverType represents the type of handover
type HandoverType string

const (
	HandoverTypePartial  HandoverType = "PARTIAL"
	HandoverTypeFull     HandoverType = "FULL"
	HandoverTypeEndShift HandoverType = "END_SHIFT"
)

// ResponsibilityType represents who is responsible for a discrepancy
type ResponsibilityType string

const (
	ResponsibilityWaiter  ResponsibilityType = "WAITER"
	ResponsibilityCashier ResponsibilityType = "CASHIER"
	ResponsibilitySystem  ResponsibilityType = "SYSTEM"
	ResponsibilityUnknown ResponsibilityType = "UNKNOWN"
)

// CashHandover represents a cash handover request from waiter to cashier
type CashHandover struct {
	ID primitive.ObjectID `bson:"_id,omitempty" json:"id"`

	// Handover details
	Type   HandoverType   `bson:"type" json:"type"`
	Status HandoverStatus `bson:"status" json:"status"`

	// Waiter information
	WaiterShiftID primitive.ObjectID `bson:"waiter_shift_id" json:"waiter_shift_id"`
	WaiterID      primitive.ObjectID `bson:"waiter_id" json:"waiter_id"`
	WaiterName    string             `bson:"waiter_name" json:"waiter_name"`

	// Cashier information (assigned when confirmed)
	CashierShiftID *primitive.ObjectID `bson:"cashier_shift_id,omitempty" json:"cashier_shift_id,omitempty"`
	CashierID      *primitive.ObjectID `bson:"cashier_id,omitempty" json:"cashier_id,omitempty"`
	CashierName    *string             `bson:"cashier_name,omitempty" json:"cashier_name,omitempty"`

	// Cash amounts
	RequestedAmount float64  `bson:"requested_amount" json:"requested_amount"`
	ActualAmount    *float64 `bson:"actual_amount,omitempty" json:"actual_amount,omitempty"`

	// Notes and reasons
	WaiterNotes  string  `bson:"waiter_notes,omitempty" json:"waiter_notes,omitempty"`
	CashierNotes *string `bson:"cashier_notes,omitempty" json:"cashier_notes,omitempty"`

	// Discrepancy information
	DiscrepancyAmount *float64            `bson:"discrepancy_amount,omitempty" json:"discrepancy_amount,omitempty"`
	DiscrepancyReason *string             `bson:"discrepancy_reason,omitempty" json:"discrepancy_reason,omitempty"`
	Responsibility    *ResponsibilityType `bson:"responsibility,omitempty" json:"responsibility,omitempty"`

	// Manager approval (for large discrepancies)
	RequiresManagerApproval bool                `bson:"requires_manager_approval" json:"requires_manager_approval"`
	ManagerApproved         *bool               `bson:"manager_approved,omitempty" json:"manager_approved,omitempty"`
	ManagerID               *primitive.ObjectID `bson:"manager_id,omitempty" json:"manager_id,omitempty"`
	ManagerNotes            *string             `bson:"manager_notes,omitempty" json:"manager_notes,omitempty"`

	// Timestamps
	RequestedAt  time.Time  `bson:"requested_at" json:"requested_at"`
	ConfirmedAt  *time.Time `bson:"confirmed_at,omitempty" json:"confirmed_at,omitempty"`
	RejectedAt   *time.Time `bson:"rejected_at,omitempty" json:"rejected_at,omitempty"`
	ManageredAt  *time.Time `bson:"managered_at,omitempty" json:"managered_at,omitempty"`
	CreatedAt    time.Time  `bson:"created_at" json:"created_at"`
	UpdatedAt    time.Time  `bson:"updated_at" json:"updated_at"`
}

// NewCashHandover creates a new cash handover request
func NewCashHandover(
	handoverType HandoverType,
	waiterShiftID, waiterID primitive.ObjectID,
	waiterName string,
	requestedAmount float64,
	waiterNotes string,
) (*CashHandover, error) {
	if requestedAmount <= 0 {
		return nil, errors.New("requested amount must be greater than 0")
	}

	if waiterName == "" {
		return nil, errors.New("waiter name is required")
	}

	now := time.Now()
	return &CashHandover{
		Type:                    handoverType,
		Status:                  HandoverStatusPending,
		WaiterShiftID:           waiterShiftID,
		WaiterID:                waiterID,
		WaiterName:              waiterName,
		RequestedAmount:         requestedAmount,
		WaiterNotes:             waiterNotes,
		RequiresManagerApproval: false,
		RequestedAt:             now,
		CreatedAt:               now,
		UpdatedAt:               now,
	}, nil
}

// ConfirmHandover confirms the handover with actual amount and cashier details
func (ch *CashHandover) ConfirmHandover(
	cashierShiftID, cashierID primitive.ObjectID,
	cashierName string,
	actualAmount float64,
	cashierNotes string,
	discrepancyThreshold float64,
) error {
	if ch.Status != HandoverStatusPending {
		return errors.New("can only confirm pending handovers")
	}

	if actualAmount < 0 {
		return errors.New("actual amount cannot be negative")
	}

	now := time.Now()
	ch.CashierShiftID = &cashierShiftID
	ch.CashierID = &cashierID
	ch.CashierName = &cashierName
	ch.ActualAmount = &actualAmount
	ch.CashierNotes = &cashierNotes
	ch.ConfirmedAt = &now
	ch.UpdatedAt = now

	// Calculate discrepancy
	discrepancy := actualAmount - ch.RequestedAmount
	if discrepancy != 0 {
		ch.DiscrepancyAmount = &discrepancy
		ch.Status = HandoverStatusDiscrepancy

		// Check if manager approval is required for large discrepancies
		if abs(discrepancy) >= discrepancyThreshold {
			ch.RequiresManagerApproval = true
		}
	} else {
		ch.Status = HandoverStatusConfirmed
	}

	return nil
}

// RejectHandover rejects the handover with reason
func (ch *CashHandover) RejectHandover(
	cashierShiftID, cashierID primitive.ObjectID,
	cashierName string,
	reason string,
) error {
	if ch.Status != HandoverStatusPending {
		return errors.New("can only reject pending handovers")
	}

	if reason == "" {
		return errors.New("rejection reason is required")
	}

	now := time.Now()
	ch.CashierShiftID = &cashierShiftID
	ch.CashierID = &cashierID
	ch.CashierName = &cashierName
	ch.CashierNotes = &reason
	ch.Status = HandoverStatusRejected
	ch.RejectedAt = &now
	ch.UpdatedAt = now

	return nil
}

// SetDiscrepancyDetails sets the discrepancy reason and responsibility
func (ch *CashHandover) SetDiscrepancyDetails(reason string, responsibility ResponsibilityType) error {
	if ch.Status != HandoverStatusDiscrepancy {
		return errors.New("can only set discrepancy details for handovers with discrepancies")
	}

	if reason == "" {
		return errors.New("discrepancy reason is required")
	}

	ch.DiscrepancyReason = &reason
	ch.Responsibility = &responsibility
	ch.UpdatedAt = time.Now()

	return nil
}

// ApproveDiscrepancy handles manager approval/rejection of discrepancies
func (ch *CashHandover) ApproveDiscrepancy(
	managerID primitive.ObjectID,
	approved bool,
	managerNotes string,
) error {
	if ch.Status != HandoverStatusDiscrepancy {
		return errors.New("can only approve discrepancies")
	}

	if !ch.RequiresManagerApproval {
		return errors.New("this handover does not require manager approval")
	}

	now := time.Now()
	ch.ManagerID = &managerID
	ch.ManagerApproved = &approved
	ch.ManagerNotes = &managerNotes
	ch.ManageredAt = &now
	ch.UpdatedAt = now

	// If approved, mark as confirmed
	if approved {
		ch.Status = HandoverStatusConfirmed
	}

	return nil
}

// HasDiscrepancy returns true if there is a discrepancy
func (ch *CashHandover) HasDiscrepancy() bool {
	return ch.DiscrepancyAmount != nil && *ch.DiscrepancyAmount != 0
}

// GetDiscrepancyAmount returns the discrepancy amount (0 if no discrepancy)
func (ch *CashHandover) GetDiscrepancyAmount() float64 {
	if ch.DiscrepancyAmount == nil {
		return 0
	}
	return *ch.DiscrepancyAmount
}

// IsShortage returns true if the discrepancy is a shortage (negative)
func (ch *CashHandover) IsShortage() bool {
	return ch.HasDiscrepancy() && *ch.DiscrepancyAmount < 0
}

// IsOverage returns true if the discrepancy is an overage (positive)
func (ch *CashHandover) IsOverage() bool {
	return ch.HasDiscrepancy() && *ch.DiscrepancyAmount > 0
}

// CanCancel returns true if the handover can be cancelled
func (ch *CashHandover) CanCancel() bool {
	return ch.Status == HandoverStatusPending
}

// Helper function to get absolute value
func abs(x float64) float64 {
	if x < 0 {
		return -x
	}
	return x
}

// Request and Response structs for API

// CreateHandoverRequest represents a request to create a handover
type CreateHandoverRequest struct {
	Type            HandoverType `json:"type" binding:"required"`
	RequestedAmount float64      `json:"requested_amount" binding:"required,gt=0"`
	WaiterNotes     string       `json:"waiter_notes,omitempty"`
}

// ConfirmHandoverRequest represents a request to confirm a handover
type ConfirmHandoverRequest struct {
	ActualAmount  float64 `json:"actual_amount" binding:"required,gte=0"`
	CashierNotes  string  `json:"cashier_notes,omitempty"`
	QuickConfirm  bool    `json:"quick_confirm,omitempty"` // For exact amount confirmations
}

// ReconcileHandoverRequest represents a request to reconcile a handover with discrepancy
type ReconcileHandoverRequest struct {
	ActualAmount      float64            `json:"actual_amount" binding:"required,gte=0"`
	DiscrepancyReason string             `json:"discrepancy_reason" binding:"required"`
	Responsibility    ResponsibilityType `json:"responsibility" binding:"required"`
	CashierNotes      string             `json:"cashier_notes,omitempty"`
}

// ApproveDiscrepancyRequest represents a manager's approval/rejection of a discrepancy
type ApproveDiscrepancyRequest struct {
	Approved     bool   `json:"approved" binding:"required"`
	ManagerNotes string `json:"manager_notes" binding:"required"`
}

// HandoverResponse represents the response for handover operations
type HandoverResponse struct {
	*CashHandover
	DiscrepancyText string `json:"discrepancy_text,omitempty"`
}

// NewHandoverResponse creates a response with additional computed fields
func NewHandoverResponse(handover *CashHandover) *HandoverResponse {
	response := &HandoverResponse{
		CashHandover: handover,
	}

	// Add discrepancy text
	if handover.HasDiscrepancy() {
		discrepancy := handover.GetDiscrepancyAmount()
		if discrepancy > 0 {
			response.DiscrepancyText = "Overage"
		} else {
			response.DiscrepancyText = "Shortage"
		}
	}

	return response
}