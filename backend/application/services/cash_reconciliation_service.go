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

type CashReconciliationService struct {
	reconciliationRepo *mongodb.CashReconciliationRepository
	shiftRepo          *mongodb.ShiftRepository
	orderRepo          *mongodb.OrderRepository
}

func NewCashReconciliationService(
	reconciliationRepo *mongodb.CashReconciliationRepository,
	shiftRepo *mongodb.ShiftRepository,
	orderRepo *mongodb.OrderRepository,
) *CashReconciliationService {
	return &CashReconciliationService{
		reconciliationRepo: reconciliationRepo,
		shiftRepo:          shiftRepo,
		orderRepo:          orderRepo,
	}
}

type ShiftStatusResponse struct {
	Shift           *order.Shift                   `json:"shift"`
	TotalOrders     int                            `json:"total_orders"`
	TotalRevenue    float64                        `json:"total_revenue"`
	CashRevenue     float64                        `json:"cash_revenue"`
	TransferRevenue float64                        `json:"transfer_revenue"`
	QRRevenue       float64                        `json:"qr_revenue"`
	Reconciliation  *cashier.CashReconciliation    `json:"reconciliation,omitempty"`
	Closure         *cashier.ShiftClosure          `json:"closure,omitempty"`
}

// FR-CASH-06: Đối soát tiền mặt
func (s *CashReconciliationService) ReconcileCash(shiftID string, actualCash float64, notes string, cashierID string) (*cashier.CashReconciliation, error) {
	// Validate shift exists and is closed
	shiftObjID, err := primitive.ObjectIDFromHex(shiftID)
	if err != nil {
		return nil, errors.New("invalid shift ID")
	}

	shift, err := s.shiftRepo.FindByID(context.Background(), shiftObjID)
	if err != nil {
		return nil, errors.New("shift not found")
	}

	if shift.Status != order.ShiftClosed {
		return nil, errors.New("can only reconcile closed shifts")
	}

	// Calculate expected cash from orders
	orders, err := s.orderRepo.FindByShiftID(context.Background(), shiftObjID)
	if err != nil {
		return nil, err
	}

	expectedCash := 0.0
	for _, ord := range orders {
		if ord.PaymentMethod == order.PaymentCash && (ord.Status == order.StatusPaid || ord.Status == order.StatusInProgress || ord.Status == order.StatusServed) {
			expectedCash += ord.Total
		}
	}

	// Create reconciliation
	reconciliation := &cashier.CashReconciliation{
		ShiftID:          shiftID,
		CashierID:        cashierID,
		ExpectedCash:     expectedCash,
		ActualCash:       actualCash,
		Notes:            notes,
		ReconciliationAt: time.Now(),
	}

	reconciliation.CalculateDifference()

	err = s.reconciliationRepo.Create(reconciliation)
	if err != nil {
		return nil, err
	}

	return reconciliation, nil
}

// FR-CASH-02: Theo dõi trạng thái ca
func (s *CashReconciliationService) GetShiftStatus(shiftID string) (*ShiftStatusResponse, error) {
	shiftObjID, err := primitive.ObjectIDFromHex(shiftID)
	if err != nil {
		return nil, errors.New("invalid shift ID")
	}

	shift, err := s.shiftRepo.FindByID(context.Background(), shiftObjID)
	if err != nil {
		return nil, err
	}

	orders, err := s.orderRepo.FindByShiftID(context.Background(), shiftObjID)
	if err != nil {
		return nil, err
	}

	response := &ShiftStatusResponse{
		Shift:       shift,
		TotalOrders: len(orders),
	}

	// Calculate revenue by payment method
	for _, ord := range orders {
		if ord.Status == order.StatusPaid || ord.Status == order.StatusInProgress || ord.Status == order.StatusServed {
			response.TotalRevenue += ord.Total
			switch ord.PaymentMethod {
			case order.PaymentCash:
				response.CashRevenue += ord.Total
			case order.PaymentTransfer:
				response.TransferRevenue += ord.Total
			case order.PaymentQR:
				response.QRRevenue += ord.Total
			}
		}
	}

	// Get reconciliation if exists
	reconciliation, err := s.reconciliationRepo.FindByShiftID(shiftID)
	if err == nil {
		response.Reconciliation = reconciliation
	}

	return response, nil
}

func (s *CashReconciliationService) GetReconciliationsByDateRange(start, end time.Time) ([]*cashier.CashReconciliation, error) {
	return s.reconciliationRepo.FindByDateRange(start, end)
}

func (s *CashReconciliationService) GetReconciliationsByCashier(cashierID string) ([]*cashier.CashReconciliation, error) {
	return s.reconciliationRepo.FindByCashierID(cashierID)
}