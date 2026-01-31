package order

import (
	"time"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type OrderStatus string

const (
	StatusCreated    OrderStatus = "CREATED"     // Order chưa thanh toán
	StatusPaid       OrderStatus = "PAID"        // Đã thu tiền, chưa giao cho pha chế
	StatusQueued     OrderStatus = "QUEUED"      // Đã gửi cho barista, chờ nhận
	StatusInProgress OrderStatus = "IN_PROGRESS" // Barista đã nhận và đang pha
	StatusReady      OrderStatus = "READY"       // Pha xong, chờ giao
	StatusServed     OrderStatus = "SERVED"      // Đã giao cho khách
	StatusCancelled  OrderStatus = "CANCELLED"   // Đã hủy
	StatusRefunded   OrderStatus = "REFUNDED"    // Đã hoàn tiền
	StatusLocked     OrderStatus = "LOCKED"      // Đã chốt ca
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
	BaristaID       primitive.ObjectID `bson:"barista_id,omitempty" json:"barista_id,omitempty"`
	BaristaName     string             `bson:"barista_name,omitempty" json:"barista_name,omitempty"`
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
	QueuedAt        *time.Time         `bson:"queued_at,omitempty" json:"queued_at,omitempty"`
	AcceptedAt      *time.Time         `bson:"accepted_at,omitempty" json:"accepted_at,omitempty"`
	ReadyAt         *time.Time         `bson:"ready_at,omitempty" json:"ready_at,omitempty"`
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
	// BR-01: State machine transitions
	transitions := map[OrderStatus][]OrderStatus{
		StatusCreated:    {StatusPaid, StatusCancelled},
		StatusPaid:       {StatusPaid, StatusQueued, StatusCancelled}, // Can stay PAID for edits/additional payments
		StatusQueued:     {StatusInProgress, StatusCancelled},         // Only barista can move to IN_PROGRESS
		StatusInProgress: {StatusReady},                               // Only barista can mark as READY
		StatusReady:      {StatusServed},                              // Only waiter can deliver
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
	// BR-08: Payment adjustments are allowed only before QUEUED
	return o.Status == StatusCreated || o.Status == StatusPaid
}

func (o *Order) IsLocked() bool {
	return o.Status == StatusLocked
}

func (o *Order) IsFullyPaid() bool {
	return o.AmountDue <= 0
}

func (o *Order) CanModify() bool {
	// BR-07: Once order enters IN_PROGRESS, no modification or refund is allowed
	return o.Status == StatusCreated || o.Status == StatusPaid || o.Status == StatusQueued
}

func (o *Order) CanRefund() bool {
	// BR-08: Refunds only allowed before QUEUED
	return o.Status == StatusPaid && o.AmountPaid > 0
}
