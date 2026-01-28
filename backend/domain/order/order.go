package order

import (
	"time"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type OrderStatus string

const (
	StatusCreated    OrderStatus = "CREATED"
	StatusPaid       OrderStatus = "PAID"
	StatusInProgress OrderStatus = "IN_PROGRESS"
	StatusServed     OrderStatus = "SERVED"
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
	OrderNumber     string             `bson:"order_number" json:"order_number"`
	CustomerName    string             `bson:"customer_name,omitempty" json:"customer_name,omitempty"`
	WaiterID        primitive.ObjectID `bson:"waiter_id" json:"waiter_id"`
	WaiterName      string             `bson:"waiter_name" json:"waiter_name"`
	ShiftID         primitive.ObjectID `bson:"shift_id" json:"shift_id"`
	Items           []OrderItem        `bson:"items" json:"items"`
	Subtotal        float64            `bson:"subtotal" json:"subtotal"`
	Discount        float64            `bson:"discount" json:"discount"`
	Total           float64            `bson:"total" json:"total"`
	AmountPaid      float64            `bson:"amount_paid" json:"amount_paid"`
	AmountDue       float64            `bson:"amount_due" json:"amount_due"`
	Status          OrderStatus        `bson:"status" json:"status"`
	PaymentMethod   PaymentMethod      `bson:"payment_method,omitempty" json:"payment_method,omitempty"`
	CollectorID     primitive.ObjectID `bson:"collector_id,omitempty" json:"collector_id,omitempty"`
	CollectorName   string             `bson:"collector_name,omitempty" json:"collector_name,omitempty"`
	Note            string             `bson:"note,omitempty" json:"note,omitempty"`
	CancelReason    string             `bson:"cancel_reason,omitempty" json:"cancel_reason,omitempty"`
	RefundAmount    float64            `bson:"refund_amount,omitempty" json:"refund_amount,omitempty"`
	RefundReason    string             `bson:"refund_reason,omitempty" json:"refund_reason,omitempty"`
	CreatedAt       time.Time          `bson:"created_at" json:"created_at"`
	UpdatedAt       time.Time          `bson:"updated_at" json:"updated_at"`
	PaidAt          *time.Time         `bson:"paid_at,omitempty" json:"paid_at,omitempty"`
	SentToBarAt     *time.Time         `bson:"sent_to_bar_at,omitempty" json:"sent_to_bar_at,omitempty"`
	ServedAt        *time.Time         `bson:"served_at,omitempty" json:"served_at,omitempty"`
	LockedAt        *time.Time         `bson:"locked_at,omitempty" json:"locked_at,omitempty"`
}

type CreateOrderRequest struct {
	CustomerName string      `json:"customer_name"`
	Items        []OrderItem `json:"items" binding:"required,min=1"`
	Note         string      `json:"note"`
	WaiterID     string      `json:"waiter_id"`
	ShiftID      string      `json:"shift_id"`
}

type PaymentRequest struct {
	PaymentMethod PaymentMethod `json:"payment_method" binding:"required"`
	Amount        float64       `json:"amount" binding:"required,gt=0"`
	CollectorID   string        `json:"collector_id"`
	CollectorName string        `json:"collector_name"`
}

type EditOrderRequest struct {
	Items    []OrderItem `json:"items" binding:"required,min=1"`
	Discount float64     `json:"discount"`
	Note     string      `json:"note"`
}

type EditOrderResponse struct {
	Order        *Order  `json:"order"`
	RefundAmount float64 `json:"refund_amount,omitempty"`
	RefundReason string  `json:"refund_reason,omitempty"`
	Message      string  `json:"message,omitempty"`
}

type RefundRequest struct {
	Amount float64 `json:"amount" binding:"required,gt=0"`
	Reason string  `json:"reason" binding:"required"`
}

type CancelOrderRequest struct {
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
	o.AmountDue = o.Total - o.AmountPaid
	if o.AmountDue < 0 {
		o.AmountDue = 0
	}
}

func (o *Order) CanTransitionTo(newStatus OrderStatus) bool {
	transitions := map[OrderStatus][]OrderStatus{
		StatusCreated:    {StatusPaid, StatusCancelled},
		StatusPaid:       {StatusPaid, StatusInProgress, StatusCancelled}, // Can stay PAID for edits/additional payments
		StatusInProgress: {StatusServed},
		StatusServed:     {StatusLocked},
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
	return o.Status == StatusCreated || o.Status == StatusPaid
}

func (o *Order) IsLocked() bool {
	return o.Status == StatusLocked
}

func (o *Order) IsFullyPaid() bool {
	return o.AmountDue <= 0
}
