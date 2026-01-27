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
	TableName     string    `json:"table_name"`
	Amount        float64   `json:"amount"`
	PaymentMethod string    `json:"payment_method"`
	Status        string    `json:"status"`
	PaidAt        time.Time `json:"paid_at"`
}

// FR-CASH-04: Giám sát thanh toán
func (s *PaymentOversightService) GetPaymentsByShift(shiftID string) ([]*PaymentSummary, error) {
	shiftObjID, err := primitive.ObjectIDFromHex(shiftID)
	if err != nil {
		return nil, errors.New("invalid shift ID")
	}

	orders, err := s.orderRepo.FindByShiftID(context.Background(), shiftObjID)
	if err != nil {
		return nil, err
	}

	var payments []*PaymentSummary
	for _, ord := range orders {
		if ord.Status == order.StatusPaid || ord.Status == order.StatusInProgress || ord.Status == order.StatusServed {
			payment := &PaymentSummary{
				OrderID:       ord.ID.Hex(),
				TableName:     ord.TableName,
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
		string(order.StatusUnpaid), // Reset to unpaid for override
		ord.Total,
	)

	err = s.auditRepo.Create(audit)
	if err != nil {
		return err
	}

	// Update order status
	ord.Status = order.StatusUnpaid
	ord.PaymentMethod = ""
	ord.PaidAt = nil
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