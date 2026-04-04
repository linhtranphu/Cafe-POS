package services

import (
	"context"
	"errors"
	"time"

	"cafe-pos/backend/domain/cashier"
	"cafe-pos/backend/domain/order"
	"cafe-pos/backend/infrastructure/mongodb"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type PaymentOversightService struct {
	orderRepo       *mongodb.OrderRepository
	discrepancyRepo *mongodb.PaymentDiscrepancyRepository
	auditRepo       *mongodb.PaymentAuditRepository
}

func NewPaymentOversightService(
	orderRepo *mongodb.OrderRepository,
	discrepancyRepo *mongodb.PaymentDiscrepancyRepository,
	auditRepo *mongodb.PaymentAuditRepository,
) *PaymentOversightService {
	return &PaymentOversightService{
		orderRepo:       orderRepo,
		discrepancyRepo: discrepancyRepo,
		auditRepo:       auditRepo,
	}
}

type PaymentSummary struct {
	OrderID       string    `json:"order_id"`
	OrderNumber   string    `json:"order_number"`
	Amount        float64   `json:"amount"`
	PaymentMethod string    `json:"payment_method"`
	Status        string    `json:"status"`
	PaidAt        time.Time `json:"paid_at"`
}

// FR-CASH-04: Giám sát thanh toán
// In distribution model: Cashier monitors payments from waiter shifts
// Orders belong to waiter shifts, cashier receives cash via handovers
func (s *PaymentOversightService) GetPaymentsByShift(shiftID string) ([]*PaymentSummary, error) {
	shiftObjID, err := primitive.ObjectIDFromHex(shiftID)
	if err != nil {
		return nil, errors.New("invalid shift ID")
	}

	// Find orders by waiter shift ID
	orders, err := s.orderRepo.FindByShiftID(context.Background(), shiftObjID)
	if err != nil {
		return nil, err
	}

	var payments []*PaymentSummary
	for _, ord := range orders {
		// Include all paid orders
		if ord.Status == order.StatusPaid || ord.Status == order.StatusInProgress || 
		   ord.Status == order.StatusServed || ord.Status == order.StatusQueued {
			payment := &PaymentSummary{
				OrderID:       ord.ID.Hex(),
				OrderNumber:   ord.OrderNumber,
				Amount:        ord.Total,
				PaymentMethod: string(ord.PaymentMethod),
				Status:        string(ord.Status),
				PaidAt:        *ord.PaidAt,
			}
			payments = append(payments, payment)
		}
	}

	return payments, nil
}

// GetTodayPayments - For cashier to view all payments today
// In distribution model, cashier doesn't have direct orders
// Cashier monitors all payments and receives cash via handovers
func (s *PaymentOversightService) GetTodayPayments() ([]*PaymentSummary, error) {
	// Get all orders
	orders, err := s.orderRepo.FindAll(context.Background())
	if err != nil {
		return nil, err
	}

	// Filter today's paid orders
	today := time.Now().Truncate(24 * time.Hour)
	var payments []*PaymentSummary
	
	for _, ord := range orders {
		// Only include paid orders
		if ord.Status == order.StatusPaid || ord.Status == order.StatusInProgress || 
		   ord.Status == order.StatusServed || ord.Status == order.StatusQueued {
			// Check if paid today
			if ord.PaidAt != nil && ord.PaidAt.Truncate(24*time.Hour).Equal(today) {
				payment := &PaymentSummary{
					OrderID:       ord.ID.Hex(),
					OrderNumber:   ord.OrderNumber,
					Amount:        ord.Total,
					PaymentMethod: string(ord.PaymentMethod),
					Status:        string(ord.Status),
					PaidAt:        *ord.PaidAt,
				}
				payments = append(payments, payment)
			}
		}
	}

	return payments, nil
}

// OrderSummary represents an order for cashier overview
type OrderSummary struct {
	OrderID       string    `json:"order_id"`
	OrderNumber   string    `json:"order_number"`
	CustomerName  string    `json:"customer_name"`
	WaiterName    string    `json:"waiter_name"`
	Total         float64   `json:"total"`
	Status        string    `json:"status"`
	PaymentMethod string    `json:"payment_method"`
	BillPrinted   bool      `json:"bill_printed"`
	CreatedAt     time.Time `json:"created_at"`
}

// GetAllOrdersByShift returns all orders for a shift (including unpaid)
func (s *PaymentOversightService) GetAllOrdersByShift(shiftID string) ([]*OrderSummary, error) {
	shiftObjID, err := primitive.ObjectIDFromHex(shiftID)
	if err != nil {
		return nil, errors.New("invalid shift ID")
	}

	orders, err := s.orderRepo.FindByShiftID(context.Background(), shiftObjID)
	if err != nil {
		return nil, err
	}

	var result []*OrderSummary
	for _, ord := range orders {
		result = append(result, &OrderSummary{
			OrderID:       ord.ID.Hex(),
			OrderNumber:   ord.OrderNumber,
			CustomerName:  ord.CustomerName,
			WaiterName:    ord.WaiterName,
			Total:         ord.Total,
			Status:        string(ord.Status),
			PaymentMethod: string(ord.PaymentMethod),
			BillPrinted:   ord.BillPrinted,
			CreatedAt:     ord.CreatedAt,
		})
	}

	return result, nil
}

// FR-CASH-05: Xử lý sai lệch thanh toán
func (s *PaymentOversightService) ReportDiscrepancy(orderID, reason string, amount float64, cashierID string) error {
	// Validate order exists
	orderObjID, err := primitive.ObjectIDFromHex(orderID)
	if err != nil {
		return errors.New("invalid order ID")
	}

	_, err = s.orderRepo.FindByID(context.Background(), orderObjID)
	if err != nil {
		return errors.New("order not found")
	}

	discrepancy := &cashier.PaymentDiscrepancy{
		OrderID:    orderID,
		CashierID:  cashierID,
		Reason:     reason,
		Amount:     amount,
		Status:     cashier.DiscrepancyStatusPending,
		ReportedAt: time.Now(),
	}

	return s.discrepancyRepo.Create(discrepancy)
}

// FR-CASH-08: Hủy/điều chỉnh thanh toán
func (s *PaymentOversightService) OverridePayment(orderID, reason string, cashierID string) error {
	orderObjID, err := primitive.ObjectIDFromHex(orderID)
	if err != nil {
		return errors.New("invalid order ID")
	}

	ord, err := s.orderRepo.FindByID(context.Background(), orderObjID)
	if err != nil {
		return errors.New("order not found")
	}

	if ord.Status == order.StatusLocked {
		return errors.New("cannot override locked order")
	}

	oldStatus := string(ord.Status)

	// Create audit record
	audit := cashier.NewPaymentAudit(
		orderID,
		cashier.AuditActionOverride,
		cashierID,
		reason,
		oldStatus,
		string(order.StatusCreated), // Reset to CREATED for override
		ord.Total,
	)

	err = s.auditRepo.Create(audit)
	if err != nil {
		return err
	}

	// Update order status
	ord.Status = order.StatusCreated
	ord.PaymentMethod = ""
	ord.PaidAt = nil
	ord.AmountPaid = 0
	ord.UpdatedAt = time.Now()

	return s.orderRepo.Update(context.Background(), ord.ID, ord)
}

// FR-CASH-09: Khóa order
func (s *PaymentOversightService) LockOrder(orderID string, cashierID string) error {
	orderObjID, err := primitive.ObjectIDFromHex(orderID)
	if err != nil {
		return errors.New("invalid order ID")
	}

	ord, err := s.orderRepo.FindByID(context.Background(), orderObjID)
	if err != nil {
		return errors.New("order not found")
	}

	if !ord.CanTransitionTo(order.StatusLocked) {
		return errors.New("cannot lock order in current status")
	}

	oldStatus := string(ord.Status)

	// Create audit record
	audit := cashier.NewPaymentAudit(
		orderID,
		cashier.AuditActionLock,
		cashierID,
		"Order locked by cashier",
		oldStatus,
		string(order.StatusLocked),
		ord.Total,
	)

	err = s.auditRepo.Create(audit)
	if err != nil {
		return err
	}

	// Update order status
	ord.Status = order.StatusLocked
	ord.UpdatedAt = time.Now()
	now := time.Now()
	ord.LockedAt = &now

	return s.orderRepo.Update(context.Background(), ord.ID, ord)
}

func (s *PaymentOversightService) GetPendingDiscrepancies() ([]*cashier.PaymentDiscrepancy, error) {
	return s.discrepancyRepo.FindPendingDiscrepancies()
}

func (s *PaymentOversightService) ResolveDiscrepancy(discrepancyID string) error {
	return s.discrepancyRepo.UpdateStatus(discrepancyID, cashier.DiscrepancyStatusResolved)
}

func (s *PaymentOversightService) GetAuditsByOrder(orderID string) ([]*cashier.PaymentAudit, error) {
	return s.auditRepo.FindByOrderID(orderID)
}

func (s *PaymentOversightService) GetAuditsByCashier(cashierID string) ([]*cashier.PaymentAudit, error) {
	return s.auditRepo.FindByCashierID(cashierID)
}