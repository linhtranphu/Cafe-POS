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

type Shift struct {
	ID            primitive.ObjectID `bson:"_id,omitempty" json:"id"`
	Type          ShiftType          `bson:"type" json:"type"`
	Status        ShiftStatus        `bson:"status" json:"status"`
	WaiterID      primitive.ObjectID `bson:"waiter_id" json:"waiter_id"`
	WaiterName    string             `bson:"waiter_name" json:"waiter_name"`
	CashierID     primitive.ObjectID `bson:"cashier_id,omitempty" json:"cashier_id,omitempty"`
	CashierName   string             `bson:"cashier_name,omitempty" json:"cashier_name,omitempty"`
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
	WaiterID  string    `json:"waiter_id"`
}

type EndShiftRequest struct {
	EndCash float64 `json:"end_cash" binding:"min=0"`
}
