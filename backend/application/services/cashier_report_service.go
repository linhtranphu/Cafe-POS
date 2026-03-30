package services

import (
	"context"
	"errors"
	"time"

	"cafe-pos/backend/domain/cashier"
	"cafe-pos/backend/domain/handover"
	"cafe-pos/backend/domain/order"
	"cafe-pos/backend/infrastructure/mongodb"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type CashierReportService struct {
	orderRepo          *mongodb.OrderRepository
	reconciliationRepo *mongodb.CashReconciliationRepository
	shiftRepo          *mongodb.ShiftRepository
	auditRepo          *mongodb.PaymentAuditRepository
	handoverRepo       *mongodb.CashHandoverRepository
}

func NewCashierReportService(
	orderRepo *mongodb.OrderRepository,
	reconciliationRepo *mongodb.CashReconciliationRepository,
	shiftRepo *mongodb.ShiftRepository,
	auditRepo *mongodb.PaymentAuditRepository,
	handoverRepo *mongodb.CashHandoverRepository,
) *CashierReportService {
	return &CashierReportService{
		orderRepo:          orderRepo,
		reconciliationRepo: reconciliationRepo,
		shiftRepo:          shiftRepo,
		auditRepo:          auditRepo,
		handoverRepo:       handoverRepo,
	}
}

type ItemSalesBreakdown struct {
	Name        string  `json:"name"`
	VariantName string  `json:"variant_name,omitempty"`
	Quantity    int     `json:"quantity"`
	Revenue     float64 `json:"revenue"`
}

type HandoverEntry struct {
	WaiterName             string  `json:"waiter_name"`
	Status                 string  `json:"status"`
	CashDeclared           float64 `json:"cash_declared"`
	CashActual             float64 `json:"cash_actual"`
	CashDiscrepancy        float64 `json:"cash_discrepancy"`
	TransferDeclared       float64 `json:"transfer_declared"`
	TransferActual         float64 `json:"transfer_actual"`
	TransferDiscrepancy    float64 `json:"transfer_discrepancy"`
}

type HandoverSummary struct {
	Entries          []HandoverEntry `json:"entries"`
	TotalDeclared    float64         `json:"total_declared"`
	TotalActual      float64         `json:"total_actual"`
	TotalShortage    float64         `json:"total_shortage"`
	TotalOverage     float64         `json:"total_overage"`
}

type ShiftReport struct {
	Shift           *order.Shift                `json:"shift"`
	TotalOrders     int                         `json:"total_orders"`
	TotalRevenue    float64                     `json:"total_revenue"`
	CashRevenue     float64                     `json:"cash_revenue"`
	TransferRevenue float64                     `json:"transfer_revenue"`
	QRRevenue       float64                     `json:"qr_revenue"`
	TotalItemsSold  int                         `json:"total_items_sold"`
	ItemsSold       []ItemSalesBreakdown        `json:"items_sold"`
	Handover        *HandoverSummary            `json:"handover,omitempty"`
	Reconciliation  *cashier.CashReconciliation `json:"reconciliation,omitempty"`
	Audits          []*cashier.PaymentAudit     `json:"audits"`
	GeneratedAt     time.Time                   `json:"generated_at"`
}

// aggregateItems counts quantity and revenue per item+variant across a list of orders
func aggregateItems(orders []*order.Order) ([]ItemSalesBreakdown, int) {
	type key struct{ name, variant string }
	itemMap := make(map[key]*ItemSalesBreakdown)

	for _, ord := range orders {
		if ord.PaymentMethod == "" {
			continue // chưa thanh toán — không tính vào danh sách món đã bán
		}
		if ord.Status != order.StatusPaid && ord.Status != order.StatusQueued &&
			ord.Status != order.StatusInProgress && ord.Status != order.StatusReady &&
			ord.Status != order.StatusServed && ord.Status != order.StatusLocked {
			continue
		}
		for _, item := range ord.Items {
			k := key{item.Name, item.VariantName}
			if _, ok := itemMap[k]; !ok {
				itemMap[k] = &ItemSalesBreakdown{Name: item.Name, VariantName: item.VariantName}
			}
			itemMap[k].Quantity += item.Quantity
			itemMap[k].Revenue += item.Subtotal
		}
	}

	result := make([]ItemSalesBreakdown, 0, len(itemMap))
	total := 0
	for _, v := range itemMap {
		result = append(result, *v)
		total += v.Quantity
	}
	// Sort by quantity descending
	for i := 0; i < len(result)-1; i++ {
		for j := i + 1; j < len(result); j++ {
			if result[j].Quantity > result[i].Quantity {
				result[i], result[j] = result[j], result[i]
			}
		}
	}
	return result, total
}

// aggregateHandovers builds HandoverSummary from a list of handover records
func aggregateHandovers(hs []*handover.CashHandover) *HandoverSummary {
	if len(hs) == 0 {
		return nil
	}
	summary := &HandoverSummary{}
	for _, h := range hs {
		entry := HandoverEntry{
			WaiterName:          h.WaiterName,
			Status:              string(h.Status),
			CashDeclared:        h.CashDeclaredAmount,
			CashActual:          h.CashActualAmount,
			CashDiscrepancy:     h.CashDiscrepancy,
			TransferDeclared:    h.TransferDeclaredAmount,
			TransferActual:      h.TransferActualAmount,
			TransferDiscrepancy: h.TransferDiscrepancy,
		}
		summary.Entries = append(summary.Entries, entry)
		summary.TotalDeclared += h.CashDeclaredAmount + h.TransferDeclaredAmount
		summary.TotalActual += h.CashActualAmount + h.TransferActualAmount
		disc := h.CashDiscrepancy + h.TransferDiscrepancy
		if disc < 0 {
			summary.TotalShortage += -disc
		} else if disc > 0 {
			summary.TotalOverage += disc
		}
	}
	return summary
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
		if ord.PaymentMethod == "" {
			continue // chưa thanh toán (cancelled trước khi pay)
		}
		if ord.Status == order.StatusPaid || ord.Status == order.StatusQueued ||
			ord.Status == order.StatusInProgress || ord.Status == order.StatusReady ||
			ord.Status == order.StatusServed || ord.Status == order.StatusLocked {
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

	// Item breakdown
	report.ItemsSold, report.TotalItemsSold = aggregateItems(orders)

	// Handover summary
	if s.handoverRepo != nil {
		handovers, err := s.handoverRepo.FindByWaiterShift(context.Background(), shiftObjID)
		if err == nil {
			report.Handover = aggregateHandovers(handovers)
		}
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
	var allOrders []*order.Order
	for _, shift := range shifts {
		orders, err := s.orderRepo.FindByShiftID(context.Background(), shift.ID)
		if err != nil {
			continue
		}

		report.TotalOrders += len(orders)
		for _, ord := range orders {
			if ord.PaymentMethod == "" {
				continue // chưa thanh toán (cancelled trước khi pay)
			}
			if ord.Status == order.StatusPaid || ord.Status == order.StatusQueued ||
				ord.Status == order.StatusInProgress || ord.Status == order.StatusReady ||
				ord.Status == order.StatusServed || ord.Status == order.StatusLocked {
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
		allOrders = append(allOrders, orders...)
	}

	// Item breakdown
	report.ItemsSold, report.TotalItemsSold = aggregateItems(allOrders)

	// Handover summary across all shifts
	if s.handoverRepo != nil {
		handovers, err := s.handoverRepo.FindByDateRange(context.Background(), startOfDay, endOfDay)
		if err == nil {
			report.Handover = aggregateHandovers(handovers)
		}
	}

	return report, nil
}