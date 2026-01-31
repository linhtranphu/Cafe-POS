package order

import (
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
