package handover

import (
	"errors"
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

// DiscrepancyStatus represents the status of a discrepancy resolution
type DiscrepancyStatus string

const (
	DiscrepancyStatusPending  DiscrepancyStatus = "PENDING"
	DiscrepancyStatusResolved DiscrepancyStatus = "RESOLVED"
	DiscrepancyStatusEscalated DiscrepancyStatus = "ESCALATED"
)

// DiscrepancyType represents the type of discrepancy
type DiscrepancyType string

const (
	DiscrepancyTypeShortage DiscrepancyType = "SHORTAGE"
	DiscrepancyTypeOverage  DiscrepancyType = "OVERAGE"
)

// CashDiscrepancy represents a detailed record of cash discrepancies for tracking and analysis
type CashDiscrepancy struct {
	ID primitive.ObjectID `bson:"_id,omitempty" json:"id"`

	// Related handover
	HandoverID primitive.ObjectID `bson:"handover_id" json:"handover_id"`

	// Discrepancy details
	Type   DiscrepancyType   `bson:"type" json:"type"`
	Status DiscrepancyStatus `bson:"status" json:"status"`
	Amount float64           `bson:"amount" json:"amount"` // Positive for overage, negative for shortage

	// Amounts for reference
	ExpectedAmount float64 `bson:"expected_amount" json:"expected_amount"`
	ActualAmount   float64 `bson:"actual_amount" json:"actual_amount"`

	// Responsibility and reason
	Responsibility ResponsibilityType `bson:"responsibility" json:"responsibility"`
	Reason         string             `bson:"reason" json:"reason"`
	Notes          string             `bson:"notes,omitempty" json:"notes,omitempty"`

	// People involved
	WaiterID    primitive.ObjectID `bson:"waiter_id" json:"waiter_id"`
	WaiterName  string             `bson:"waiter_name" json:"waiter_name"`
	CashierID   primitive.ObjectID `bson:"cashier_id" json:"cashier_id"`
	CashierName string             `bson:"cashier_name" json:"cashier_name"`

	// Manager resolution (if escalated)
	ManagerID       *primitive.ObjectID `bson:"manager_id,omitempty" json:"manager_id,omitempty"`
	ManagerName     *string             `bson:"manager_name,omitempty" json:"manager_name,omitempty"`
	ManagerApproved *bool               `bson:"manager_approved,omitempty" json:"manager_approved,omitempty"`
	ManagerNotes    *string             `bson:"manager_notes,omitempty" json:"manager_notes,omitempty"`

	// Resolution details
	ResolutionAction *string    `bson:"resolution_action,omitempty" json:"resolution_action,omitempty"`
	ResolvedAt       *time.Time `bson:"resolved_at,omitempty" json:"resolved_at,omitempty"`

	// Timestamps
	OccurredAt time.Time `bson:"occurred_at" json:"occurred_at"`
	CreatedAt  time.Time `bson:"created_at" json:"created_at"`
	UpdatedAt  time.Time `bson:"updated_at" json:"updated_at"`
}

// NewCashDiscrepancy creates a new cash discrepancy record
func NewCashDiscrepancy(
	handoverID primitive.ObjectID,
	expectedAmount, actualAmount float64,
	responsibility ResponsibilityType,
	reason string,
	waiterID primitive.ObjectID, waiterName string,
	cashierID primitive.ObjectID, cashierName string,
) (*CashDiscrepancy, error) {
	if reason == "" {
		return nil, errors.New("discrepancy reason is required")
	}

	if waiterName == "" {
		return nil, errors.New("waiter name is required")
	}

	if cashierName == "" {
		return nil, errors.New("cashier name is required")
	}

	discrepancyAmount := actualAmount - expectedAmount
	if discrepancyAmount == 0 {
		return nil, errors.New("cannot create discrepancy record for zero discrepancy")
	}

	var discrepancyType DiscrepancyType
	if discrepancyAmount > 0 {
		discrepancyType = DiscrepancyTypeOverage
	} else {
		discrepancyType = DiscrepancyTypeShortage
	}

	now := time.Now()
	return &CashDiscrepancy{
		HandoverID:     handoverID,
		Type:           discrepancyType,
		Status:         DiscrepancyStatusPending,
		Amount:         discrepancyAmount,
		ExpectedAmount: expectedAmount,
		ActualAmount:   actualAmount,
		Responsibility: responsibility,
		Reason:         reason,
		WaiterID:       waiterID,
		WaiterName:     waiterName,
		CashierID:      cashierID,
		CashierName:    cashierName,
		OccurredAt:     now,
		CreatedAt:      now,
		UpdatedAt:      now,
	}, nil
}

// Resolve marks the discrepancy as resolved with an action
func (cd *CashDiscrepancy) Resolve(action string) error {
	if cd.Status != DiscrepancyStatusPending {
		return errors.New("can only resolve pending discrepancies")
	}

	if action == "" {
		return errors.New("resolution action is required")
	}

	now := time.Now()
	cd.Status = DiscrepancyStatusResolved
	cd.ResolutionAction = &action
	cd.ResolvedAt = &now
	cd.UpdatedAt = now

	return nil
}

// Escalate escalates the discrepancy to manager
func (cd *CashDiscrepancy) Escalate() error {
	if cd.Status != DiscrepancyStatusPending {
		return errors.New("can only escalate pending discrepancies")
	}

	cd.Status = DiscrepancyStatusEscalated
	cd.UpdatedAt = time.Now()

	return nil
}

// SetManagerResolution sets the manager's resolution for escalated discrepancies
func (cd *CashDiscrepancy) SetManagerResolution(
	managerID primitive.ObjectID,
	managerName string,
	approved bool,
	managerNotes string,
) error {
	if cd.Status != DiscrepancyStatusEscalated {
		return errors.New("can only set manager resolution for escalated discrepancies")
	}

	if managerName == "" {
		return errors.New("manager name is required")
	}

	cd.ManagerID = &managerID
	cd.ManagerName = &managerName
	cd.ManagerApproved = &approved
	cd.ManagerNotes = &managerNotes
	cd.UpdatedAt = time.Now()

	// If approved, mark as resolved
	if approved {
		action := "Manager approved discrepancy"
		cd.Resolve(action)
	}

	return nil
}

// IsShortage returns true if this is a shortage discrepancy
func (cd *CashDiscrepancy) IsShortage() bool {
	return cd.Type == DiscrepancyTypeShortage
}

// IsOverage returns true if this is an overage discrepancy
func (cd *CashDiscrepancy) IsOverage() bool {
	return cd.Type == DiscrepancyTypeOverage
}

// GetAbsoluteAmount returns the absolute value of the discrepancy amount
func (cd *CashDiscrepancy) GetAbsoluteAmount() float64 {
	if cd.Amount < 0 {
		return -cd.Amount
	}
	return cd.Amount
}

// RequiresManagerApproval returns true if the discrepancy requires manager approval
func (cd *CashDiscrepancy) RequiresManagerApproval(threshold float64) bool {
	return cd.GetAbsoluteAmount() >= threshold
}

// DiscrepancyStats represents statistics about discrepancies
type DiscrepancyStats struct {
	TotalDiscrepancies int     `json:"total_discrepancies"`
	TotalShortages     int     `json:"total_shortages"`
	TotalOverages      int     `json:"total_overages"`
	TotalShortageAmount float64 `json:"total_shortage_amount"`
	TotalOverageAmount  float64 `json:"total_overage_amount"`
	NetDiscrepancy     float64 `json:"net_discrepancy"`
	PendingCount       int     `json:"pending_count"`
	ResolvedCount      int     `json:"resolved_count"`
	EscalatedCount     int     `json:"escalated_count"`
}

// NewDiscrepancyStats creates a new stats object
func NewDiscrepancyStats() *DiscrepancyStats {
	return &DiscrepancyStats{}
}

// AddDiscrepancy adds a discrepancy to the stats
func (ds *DiscrepancyStats) AddDiscrepancy(discrepancy *CashDiscrepancy) {
	ds.TotalDiscrepancies++

	if discrepancy.IsShortage() {
		ds.TotalShortages++
		ds.TotalShortageAmount += discrepancy.GetAbsoluteAmount()
	} else {
		ds.TotalOverages++
		ds.TotalOverageAmount += discrepancy.Amount
	}

	ds.NetDiscrepancy += discrepancy.Amount

	switch discrepancy.Status {
	case DiscrepancyStatusPending:
		ds.PendingCount++
	case DiscrepancyStatusResolved:
		ds.ResolvedCount++
	case DiscrepancyStatusEscalated:
		ds.EscalatedCount++
	}
}