package services

import (
	"context"
	"errors"
	"fmt"
	"time"

	"cafe-pos/backend/domain/cashier"
	"cafe-pos/backend/domain/order"
	"cafe-pos/backend/infrastructure/mongodb"
	"go.mongodb.org/mongo-driver/bson/primitive"
	mongoDriver "go.mongodb.org/mongo-driver/mongo"
)

type PaymentOversightService struct {
	orderRepo       *mongodb.OrderRepository
	discrepancyRepo *mongodb.PaymentDiscrepancyRepository
	auditRepo       *mongodb.PaymentAuditRepository
	journalService  *JournalService
	shiftRepo       *mongodb.ShiftRepository
	mongoClient     *mongoDriver.Client
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

func (s *PaymentOversightService) SetJournalService(js *JournalService) {
	s.journalService = js
}

func (s *PaymentOversightService) SetShiftRepository(repo *mongodb.ShiftRepository) {
	s.shiftRepo = repo
}

func (s *PaymentOversightService) SetMongoClient(client *mongoDriver.Client) {
	s.mongoClient = client
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

// ChangePaymentMethod sửa phương thức thanh toán của đơn đã thanh toán.
// Cập nhật kế toán (journal), bộ đếm ca (shift counters), và ghi audit.
func (s *PaymentOversightService) ChangePaymentMethod(orderID string, newMethod order.PaymentMethod, cashierID string) error {
	orderObjID, err := primitive.ObjectIDFromHex(orderID)
	if err != nil {
		return errors.New("invalid order ID")
	}

	ord, err := s.orderRepo.FindByID(context.Background(), orderObjID)
	if err != nil {
		return errors.New("order not found")
	}

	if ord.Status == order.StatusLocked {
		return errors.New("cannot change payment method of locked order")
	}
	if ord.Status != order.StatusPaid {
		return errors.New("can only change payment method of paid orders")
	}
	if ord.PaymentMethod == newMethod {
		return fmt.Errorf("order already uses payment method %s", newMethod)
	}

	oldMethod := ord.PaymentMethod
	amount := ord.AmountPaid

	// Build old/new cash+transfer amounts for journal
	var oldCash, oldTransfer, newCash, newTransfer float64
	if oldMethod == order.PaymentCash {
		oldCash = amount
	} else {
		oldTransfer = amount
	}
	if newMethod == order.PaymentCash {
		newCash = amount
	} else {
		newTransfer = amount
	}

	// Wrap order update + journal into one transaction
	doChange := func(ctx context.Context) error {
		loaded, err := s.orderRepo.FindByID(ctx, orderObjID)
		if err != nil {
			return err
		}
		loaded.PaymentMethod = newMethod
		loaded.UpdatedAt = time.Now()
		if err := s.orderRepo.Update(ctx, orderObjID, loaded); err != nil {
			return err
		}

		if s.journalService != nil {
			cashierObjID, _ := primitive.ObjectIDFromHex(cashierID)
			if _, err := s.journalService.RecordPaymentMethodCorrectionInSession(
				ctx,
				oldCash, oldTransfer,
				newCash, newTransfer,
				loaded.ID,
				loaded.OrderNumber,
				loaded.WaiterID,
				loaded.WaiterName,
			); err != nil {
				// Log but don't block — cashierObjID used to silence unused warning
				_ = cashierObjID
				return fmt.Errorf("kế toán không ghi nhận được sửa HTTT: %w", err)
			}
		}
		return nil
	}

	if s.mongoClient != nil {
		session, err := s.mongoClient.StartSession()
		if err != nil {
			return fmt.Errorf("cannot start transaction: %w", err)
		}
		defer session.EndSession(context.Background())
		err = mongoDriver.WithSession(context.Background(), session, func(sc mongoDriver.SessionContext) error {
			if err := session.StartTransaction(); err != nil {
				return err
			}
			if err := doChange(sc); err != nil {
				session.AbortTransaction(sc)
				return err
			}
			return session.CommitTransaction(sc)
		})
		if err != nil {
			return err
		}
	} else {
		if err := doChange(context.Background()); err != nil {
			return err
		}
	}

	// Update shift revenue counters (non-fatal, outside transaction)
	if s.shiftRepo != nil && !ord.ShiftID.IsZero() {
		shift, err := s.shiftRepo.FindByID(context.Background(), ord.ShiftID)
		if err == nil && shift != nil {
			// Reverse old method
			if oldMethod == order.PaymentCash {
				shift.RemainingCash -= amount
				shift.CurrentCash -= amount
			} else {
				shift.TransferRevenue -= amount
				shift.RemainingTransfer -= amount
			}
			// Add new method
			if newMethod == order.PaymentCash {
				shift.RemainingCash += amount
				shift.CurrentCash += amount
			} else {
				shift.TransferRevenue += amount
				shift.RemainingTransfer += amount
			}
			_ = s.shiftRepo.Update(context.Background(), ord.ShiftID, shift)
		}
	}

	// Create audit record
	cashierObjID, _ := primitive.ObjectIDFromHex(cashierID)
	audit := &cashier.PaymentAudit{
		OrderID:          orderID,
		Action:           cashier.AuditActionPaymentMethodChange,
		CashierID:        cashierObjID.Hex(),
		Reason:           fmt.Sprintf("Sửa HTTT: %s → %s", oldMethod, newMethod),
		OldStatus:        string(order.StatusPaid),
		NewStatus:        string(order.StatusPaid),
		OldPaymentMethod: string(oldMethod),
		NewPaymentMethod: string(newMethod),
		Amount:           amount,
		AuditedAt:        time.Now(),
		CreatedAt:        time.Now(),
	}
	return s.auditRepo.Create(audit)
}