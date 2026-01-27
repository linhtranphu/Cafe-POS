package order

import (
	"time"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type OrderStatus string

const (
	StatusCreated    OrderStatus = "CREATED"
	StatusUnpaid     OrderStatus = "UNPAID"
	StatusPaid       OrderStatus = "PAID"
	StatusInProgress OrderStatus = "IN_PROGRESS"
	StatusServed     OrderStatus = "SERVED"
	StatusRefunded   OrderStatus = "REFUNDED"
	StatusCancelled  OrderStatus = "CANCELLED"
	StatusLocked     OrderStatus = "LOCKED"
)

type PaymentMethod string

const (
	PaymentCash     PaymentMethod = "CASH"
	PaymentTransfer PaymentMethod = "TRANSFER"
	PaymentQR       PaymentMethod = "QR"
)

type OrderItem struct {
	MenuItemID  primitive.ObjectID `bson:"menu_item_id" json:"menu_item_id"`
	Name        string             `bson:"name" json:"name"`
	Price       float64            `bson:"price" json:"price"`
	Quantity    int                `bson:"quantity" json:"quantity"`
	Note        string             `bson:"note,omitempty" json:"note,omitempty"`
	Subtotal    float64            `bson:"subtotal" json:"subtotal"`
}

type Order struct {
	ID              primitive.ObjectID `bson:"_id,omitempty" json:"id"`
	TableID         primitive.ObjectID `bson:"table_id" json:"table_id"`
	TableName       string             `bson:"table_name" json:"table_name"`
	WaiterID        primitive.ObjectID `bson:"waiter_id" json:"waiter_id"`
	WaiterName      string             `bson:"waiter_name" json:"waiter_name"`
	ShiftID         primitive.ObjectID `bson:"shift_id" json:"shift_id"`
	Items           []OrderItem        `bson:"items" json:"items"`
	Subtotal        float64            `bson:"subtotal" json:"subtotal"`
	Discount        float64            `bson:"discount" json:"discount"`
	Total           float64            `bson:"total" json:"total"`
	Status          OrderStatus        `bson:"status" json:"status"`
	PaymentMethod   PaymentMethod      `bson:"payment_method,omitempty" json:"payment_method,omitempty"`
	CollectorID     primitive.ObjectID `bson:"collector_id,omitempty" json:"collector_id,omitempty"`
	CollectorName   string             `bson:"collector_name,omitempty" json:"collector_name,omitempty"`
	Note            string             `bson:"note,omitempty" json:"note,omitempty"`
	CancelReason    string             `bson:"cancel_reason,omitempty" json:"cancel_reason,omitempty"`
	RefundReason    string             `bson:"refund_reason,omitempty" json:"refund_reason,omitempty"`
	CreatedAt       time.Time          `bson:"created_at" json:"created_at"`
	UpdatedAt       time.Time          `bson:"updated_at" json:"updated_at"`
	PaidAt          *time.Time         `bson:"paid_at,omitempty" json:"paid_at,omitempty"`
	SentToKitchenAt *time.Time         `bson:"sent_to_kitchen_at,omitempty" json:"sent_to_kitchen_at,omitempty"`
	ServedAt        *time.Time         `bson:"served_at,omitempty" json:"served_at,omitempty"`
	LockedAt        *time.Time         `bson:"locked_at,omitempty" json:"locked_at,omitempty"`
}

type CreateOrderRequest struct {
	TableID  string      `json:"table_id" binding:"required"`
	Items    []OrderItem `json:"items" binding:"required,min=1"`
	Note     string      `json:"note"`
	WaiterID string      `json:"waiter_id"`
	ShiftID  string      `json:"shift_id"`
}

type PaymentRequest struct {
	PaymentMethod PaymentMethod `json:"payment_method" binding:"required"`
	CollectorID   string        `json:"collector_id"`
	CollectorName string        `json:"collector_name"`
}

type CancelOrderRequest struct {
	Reason string `json:"reason" binding:"required"`
}

type RefundOrderRequest struct {
	Reason string `json:"reason" binding:"required"`
}

func (o *Order) CalculateTotal() {
	o.Subtotal = 0
	for i := range o.Items {
		o.Items[i].Subtotal = o.Items[i].Price * float64(o.Items[i].Quantity)
		o.Subtotal += o.Items[i].Subtotal
	}
	o.Total = o.Subtotal - o.Discount
	if o.Total < 0 {
		o.Total = 0
	}
}

func (o *Order) CanTransitionTo(newStatus OrderStatus) bool {
	transitions := map[OrderStatus][]OrderStatus{
		StatusCreated:    {StatusUnpaid, StatusCancelled},
		StatusUnpaid:     {StatusPaid, StatusCancelled},
		StatusPaid:       {StatusInProgress, StatusRefunded},
		StatusInProgress: {StatusServed, StatusRefunded},
		StatusServed:     {StatusLocked},
		StatusRefunded:   {StatusLocked},
		StatusCancelled:  {StatusLocked},
		StatusLocked:     {},
	}
	
	allowedTransitions, exists := transitions[o.Status]
	if !exists {
		return false
	}
	
	for _, allowed := range allowedTransitions {
		if allowed == newStatus {
			return true
		}
	}
	return false
}

func (o *Order) IsEditable() bool {
	return o.Status == StatusCreated || o.Status == StatusUnpaid
}

func (o *Order) IsLocked() bool {
	return o.Status == StatusLocked
}
