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

type HandoverData struct {
	FromCashierID string `json:"from_cashier_id"`
	ToCashierID   string `json:"to_cashier_id"`
	Notes         string `json:"notes"`
}

// FR-CASH-11: Bàn giao ca
func (s *CashierReportService) HandoverShift(data *HandoverData) error {
	// Find open shifts for from_cashier
	openShifts, err := s.shiftRepo.FindOpenShifts(context.Background())
	if err != nil {
		return err
	}

	var fromCashierShifts []*order.Shift
	for _, shift := range openShifts {
		if shift.WaiterID.Hex() == data.FromCashierID {
			fromCashierShifts = append(fromCashierShifts, shift)
		}
	}

	if len(fromCashierShifts) == 0 {
		return errors.New("no open shifts found for cashier")
	}

	// Transfer shifts to new cashier
	toCashierObjID, err := primitive.ObjectIDFromHex(data.ToCashierID)
	if err != nil {
		return errors.New("invalid to_cashier_id")
	}

	for _, shift := range fromCashierShifts {
		shift.WaiterID = toCashierObjID
		shift.UpdatedAt = time.Now()
		
		err = s.shiftRepo.Update(context.Background(), shift.ID, shift)
		if err != nil {
			return err
		}

		// Create audit record for handover
		audit := cashier.NewPaymentAudit(
			"", // No specific order
			"HANDOVER",
			data.FromCashierID,
			"Shift handover to "+data.ToCashierID+": "+data.Notes,
			"OPEN",
			"TRANSFERRED",
			0,
		)
		s.auditRepo.Create(audit)
	}

	return nil
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