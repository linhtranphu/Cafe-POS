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

type CashierReportService struct {
	orderRepo          *mongodb.OrderRepository
	reconciliationRepo *mongodb.CashReconciliationRepository
	shiftRepo          *mongodb.ShiftRepository
	auditRepo          *mongodb.PaymentAuditRepository
}

func NewCashierReportService(
	orderRepo *mongodb.OrderRepository,
	reconciliationRepo *mongodb.CashReconciliationRepository,
	shiftRepo *mongodb.ShiftRepository,
	auditRepo *mongodb.PaymentAuditRepository,
) *CashierReportService {
	return &CashierReportService{
		orderRepo:          orderRepo,
		reconciliationRepo: reconciliationRepo,
		shiftRepo:          shiftRepo,
		auditRepo:          auditRepo,
	}
}

type ShiftReport struct {
	Shift           *order.Shift                `json:"shift"`
	TotalOrders     int                         `json:"total_orders"`
	TotalRevenue    float64                     `json:"total_revenue"`
	CashRevenue     float64                     `json:"cash_revenue"`
	TransferRevenue float64                     `json:"transfer_revenue"`
	QRRevenue       float64                     `json:"qr_revenue"`
	Reconciliation  *cashier.CashReconciliation `json:"reconciliation,omitempty"`
	Audits          []*cashier.PaymentAudit     `json:"audits"`
	GeneratedAt     time.Time                   `json:"generated_at"`
}

// FR-CASH-10: Báo cáo ca
func (s *CashierReportService) GenerateShiftReport(shiftID string) (*ShiftReport, error) {
	shiftObjID, err := primitive.ObjectIDFromHex(shiftID)
	if err != nil {
		return nil, errors.New("invalid shift ID")
	}

	shift, err := s.shiftRepo.FindByID(context.Background(), shiftObjID)
	if err != nil {
		return nil, errors.New("shift not found")
	}

	orders, err := s.orderRepo.FindByShiftID(context.Background(), shiftObjID)
	if err != nil {
		return nil, err
	}

	report := &ShiftReport{
		Shift:       shift,
		TotalOrders: len(orders),
		GeneratedAt: time.Now(),
	}

	// Calculate revenue by payment method
	for _, ord := range orders {
		if ord.Status == order.StatusPaid || ord.Status == order.StatusInProgress || ord.Status == order.StatusServed {
			report.TotalRevenue += ord.Total
			switch ord.PaymentMethod {
			case order.PaymentCash:
				report.CashRevenue += ord.Total
			case order.PaymentTransfer:
				report.TransferRevenue += ord.Total
			case order.PaymentQR:
				report.QRRevenue += ord.Total
			}
		}
	}

	// Get reconciliation if exists
	reconciliation, err := s.reconciliationRepo.FindByShiftID(shiftID)
	if err == nil {
		report.Reconciliation = reconciliation
	}

	// Get audit records for this shift's orders
	var allAudits []*cashier.PaymentAudit
	for _, ord := range orders {
		audits, err := s.auditRepo.FindByOrderID(ord.ID.Hex())
		if err == nil {
			allAudits = append(allAudits, audits...)
		}
	}
	report.Audits = allAudits

	return report, nil
}

func (s *CashierReportService) GetDailyReport(date time.Time) (*ShiftReport, error) {
	startOfDay := time.Date(date.Year(), date.Month(), date.Day(), 0, 0, 0, 0, date.Location())
	endOfDay := startOfDay.Add(24 * time.Hour)

	shifts, err := s.shiftRepo.FindByDateRange(context.Background(), startOfDay, endOfDay)
	if err != nil {
		return nil, err
	}

	report := &ShiftReport{
		GeneratedAt: time.Now(),
	}

	// Aggregate data from all shifts
	for _, shift := range shifts {
		orders, err := s.orderRepo.FindByShiftID(context.Background(), shift.ID)
		if err != nil {
			continue
		}

		report.TotalOrders += len(orders)
		for _, ord := range orders {
			if ord.Status == order.StatusPaid || ord.Status == order.StatusInProgress || ord.Status == order.StatusServed {
				report.TotalRevenue += ord.Total
				switch ord.PaymentMethod {
				case order.PaymentCash:
					report.CashRevenue += ord.Total
				case order.PaymentTransfer:
					report.TransferRevenue += ord.Total
				case order.PaymentQR:
					report.QRRevenue += ord.Total
				}
			}
		}
	}

	return report, nil
}