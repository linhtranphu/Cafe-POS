package handover

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

// CashDiscrepancy represents a discrepancy in cash handover
type CashDiscrepancy struct {
	ID             primitive.ObjectID `bson:"_id,omitempty" json:"id"`
	HandoverID     primitive.ObjectID `bson:"handover_id" json:"handover_id"`
	WaiterShiftID  primitive.ObjectID `bson:"waiter_shift_id" json:"waiter_shift_id"`
	CashierShiftID primitive.ObjectID `bson:"cashier_shift_id" json:"cashier_shift_id"`

	// Discrepancy details
	DeclaredAmount    float64         `bson:"declared_amount" json:"declared_amount"`
	ActualAmount      float64         `bson:"actual_amount" json:"actual_amount"`
	DiscrepancyAmount float64         `bson:"discrepancy_amount" json:"discrepancy_amount"`
	DiscrepancyType   DiscrepancyType `bson:"discrepancy_type" json:"discrepancy_type"`

	// Analysis
	ReasonCategory string             `bson:"reason_category" json:"reason_category"`
	DetailedReason string             `bson:"detailed_reason" json:"detailed_reason"`
	Responsibility ResponsibilityType `bson:"responsibility" json:"responsibility"`

	// Resolution
	ResolutionStatus string             `bson:"resolution_status" json:"resolution_status"`
	ResolutionAction string             `bson:"resolution_action" json:"resolution_action"`
	ResolvedBy       primitive.ObjectID `bson:"resolved_by,omitempty" json:"resolved_by,omitempty"`
	ResolvedAt       *time.Time         `bson:"resolved_at,omitempty" json:"resolved_at,omitempty"`

	// Manager approval
	RequiresManagerApproval bool               `bson:"requires_manager_approval" json:"requires_manager_approval"`
	ManagerApproved         bool               `bson:"manager_approved" json:"manager_approved"`
	ApprovedBy              primitive.ObjectID `bson:"approved_by,omitempty" json:"approved_by,omitempty"`
	ApprovedAt              *time.Time         `bson:"approved_at,omitempty" json:"approved_at,omitempty"`
	ManagerNote             string             `bson:"manager_note,omitempty" json:"manager_note,omitempty"`

	CreatedAt time.Time `bson:"created_at" json:"created_at"`
	UpdatedAt time.Time `bson:"updated_at" json:"updated_at"`
}
