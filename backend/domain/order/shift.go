package order

import (
	"errors"
	"time"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type ShiftStatus string

const (
	ShiftOpen   ShiftStatus = "OPEN"
	ShiftClosed ShiftStatus = "CLOSED"
)

type ShiftType string

const (
	ShiftMorning   ShiftType = "MORNING"
	ShiftAfternoon ShiftType = "AFTERNOON"
	ShiftEvening   ShiftType = "EVENING"
)

type RoleType string

const (
	RoleWaiter  RoleType = "waiter"
	RoleBarista RoleType = "barista"
)

// ParseRoleType converts a string to RoleType
// Note: Cashier shifts are now handled separately in the cashier domain
func ParseRoleType(role string) RoleType {
	switch role {
	case "waiter":
		return RoleWaiter
	case "barista":
		return RoleBarista
	default:
		return RoleWaiter // default fallback
	}
}

// IsValid checks if the RoleType is valid
func (r RoleType) IsValid() bool {
	switch r {
	case RoleWaiter, RoleBarista:
		return true
	default:
		return false
	}
}

// String returns the string representation
func (r RoleType) String() string {
	return string(r)
}

// Shift represents a work period for a waiter or barista.
// Note: Cashier shifts are handled separately in the cashier domain with CashierShift.
type Shift struct {
	ID            primitive.ObjectID `bson:"_id,omitempty" json:"id"`
	Type          ShiftType          `bson:"type" json:"type"`
	Status        ShiftStatus        `bson:"status" json:"status"`
	RoleType      RoleType           `bson:"role_type" json:"role_type"`
	UserID        primitive.ObjectID `bson:"user_id" json:"user_id"`
	UserName      string             `bson:"user_name" json:"user_name"`
	
	StartCash     float64            `bson:"start_cash" json:"start_cash"`
	EndCash       float64            `bson:"end_cash" json:"end_cash"`
	TotalRevenue  float64            `bson:"total_revenue" json:"total_revenue"`
	TotalOrders   int                `bson:"total_orders" json:"total_orders"`
	
	// Cash handover tracking fields
	CurrentCash      float64 `bson:"current_cash" json:"current_cash"`           // Current cash amount in possession
	HandedOverCash   float64 `bson:"handed_over_cash" json:"handed_over_cash"`   // Total amount handed over to cashiers
	RemainingCash    float64 `bson:"remaining_cash" json:"remaining_cash"`       // Cash remaining after handovers
	TotalDiscrepancy float64 `bson:"total_discrepancy" json:"total_discrepancy"` // Total discrepancy from all handovers
	HandoverCount    int     `bson:"handover_count" json:"handover_count"`       // Number of handovers made
	
	StartedAt     time.Time          `bson:"started_at" json:"started_at"`
	EndedAt       *time.Time         `bson:"ended_at,omitempty" json:"ended_at,omitempty"`
	CreatedAt     time.Time          `bson:"created_at" json:"created_at"`
	UpdatedAt     time.Time          `bson:"updated_at" json:"updated_at"`
}

type StartShiftRequest struct {
	Type      ShiftType `json:"type" binding:"required"`
	StartCash float64   `json:"start_cash" binding:"min=0"`
	UserID    string    `json:"user_id"`
	RoleType  RoleType  `json:"role_type"`
}

type EndShiftRequest struct {
	EndCash float64 `json:"end_cash" binding:"min=0"`
}

// UpdateCashAfterHandover updates the cash amounts after a handover
func (s *Shift) UpdateCashAfterHandover(handedOverAmount, discrepancyAmount float64) {
	s.HandedOverCash += handedOverAmount
	s.TotalDiscrepancy += discrepancyAmount
	s.HandoverCount++
	s.RemainingCash = s.CurrentCash - s.HandedOverCash
	s.UpdatedAt = time.Now()
}

// CanHandover checks if the shift can perform a handover of the specified amount
func (s *Shift) CanHandover(amount float64) error {
	if s.Status != ShiftOpen {
		return errors.New("cannot handover cash from closed shift")
	}
	
	if amount <= 0 {
		return errors.New("handover amount must be greater than 0")
	}
	
	if amount > s.RemainingCash {
		return errors.New("handover amount exceeds remaining cash")
	}
	
	return nil
}

// GetAvailableCash returns the amount of cash available for handover
func (s *Shift) GetAvailableCash() float64 {
	return s.RemainingCash
}
