package services

import (
	"context"
	"errors"
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"

	"cafe-pos/backend/domain/cashier"
	"cafe-pos/backend/domain/handover"
	"cafe-pos/backend/domain/order"
	"cafe-pos/backend/infrastructure/mongodb"
)

// CashHandoverService handles cash handover business logic
type CashHandoverService struct {
	handoverRepo         *mongodb.CashHandoverRepository
	discrepancyRepo      *mongodb.CashDiscrepancyRepository
	shiftRepo            *mongodb.ShiftRepository
	cashierShiftRepo     *mongodb.CashierShiftRepository
	orderRepo            *mongodb.OrderRepository
	discrepancyThreshold float64
}

// NewCashHandoverService creates a new cash handover service
func NewCashHandoverService(
	handoverRepo *mongodb.CashHandoverRepository,
	discrepancyRepo *mongodb.CashDiscrepancyRepository,
	shiftRepo *mongodb.ShiftRepository,
	cashierShiftRepo *mongodb.CashierShiftRepository,
	orderRepo *mongodb.OrderRepository,
) *CashHandoverService {
	return &CashHandoverService{
		handoverRepo:         handoverRepo,
		discrepancyRepo:      discrepancyRepo,
		shiftRepo:            shiftRepo,
		cashierShiftRepo:     cashierShiftRepo,
		orderRepo:            orderRepo,
		discrepancyThreshold: 100000, // 100k VND
	}
}

// CreateHandover creates a new handover request from waiter
func (s *CashHandoverService) CreateHandover(
	ctx context.Context,
	waiterShiftID primitive.ObjectID,
	req *handover.CreateHandoverRequest,
	waiterID, waiterName string,
) (*handover.CashHandover, error) {
	// 1. Validate waiter shift exists and is open
	waiterShift, err := s.shiftRepo.FindByID(ctx, waiterShiftID)
	if err != nil {
		return nil, errors.New("waiter shift not found")
	}
	if waiterShift.Status != order.ShiftOpen {
		return nil, errors.New("waiter shift is not open")
	}

	// 2. Check if waiter owns the shift
	waiterOID, err := primitive.ObjectIDFromHex(waiterID)
	if err != nil {
		return nil, errors.New("invalid waiter ID")
	}
	if waiterShift.UserID != waiterOID {
		return nil, errors.New("unauthorized: not your shift")
	}

	// 3. Check if there's already a pending handover
	existingPending, err := s.handoverRepo.FindPendingByWaiterShift(ctx, waiterShiftID)
	if err != nil {
		return nil, err
	}
	if existingPending != nil {
		return nil, errors.New("there is already a pending handover for this shift")
	}

	// 4. Validate declared amount
	if req.DeclaredAmount > waiterShift.RemainingCash {
		return nil, errors.New("declared amount exceeds remaining cash")
	}

	// 5. Find active cashier shift
	cashierShifts, err := s.cashierShiftRepo.FindByStatus(ctx, cashier.CashierShiftOpen)
	if err != nil || len(cashierShifts) == 0 {
		return nil, errors.New("no active cashier shift found")
	}
	cashierShift := cashierShifts[0] // Get first open shift

	// 6. Create handover record
	h := &handover.CashHandover{
		WaiterShiftID:  waiterShiftID,
		CashierShiftID: cashierShift.ID,
		WaiterID:       waiterOID,
		WaiterName:     waiterName,
		CashierID:      cashierShift.CashierID,
		CashierName:    cashierShift.CashierName,
		DeclaredAmount: req.DeclaredAmount,
		ActualAmount:   0, // Will be set by cashier
		Discrepancy:    0, // Will be calculated
		HandoverType:   req.HandoverType,
		Status:         handover.StatusPending,
		WaiterNote:     req.WaiterNote,
		HandoverAt:     time.Now(),
		CreatedAt:      time.Now(),
		UpdatedAt:      time.Now(),
	}

	if err := s.handoverRepo.Create(ctx, h); err != nil {
		return nil, err
	}

	return h, nil
}

// CreateHandoverAndEndShift creates handover and prepares shift for closure
func (s *CashHandoverService) CreateHandoverAndEndShift(
	ctx context.Context,
	waiterShiftID primitive.ObjectID,
	req *handover.CreateHandoverAndEndShiftRequest,
	waiterID, waiterName string,
) (*handover.CashHandover, error) {
	// 1. Validate waiter shift exists and is open
	waiterShift, err := s.shiftRepo.FindByID(ctx, waiterShiftID)
	if err != nil {
		return nil, errors.New("waiter shift not found")
	}
	if waiterShift.Status != order.ShiftOpen {
		return nil, errors.New("waiter shift is not open")
	}

	// 2. Check if waiter owns the shift
	waiterOID, err := primitive.ObjectIDFromHex(waiterID)
	if err != nil {
		return nil, errors.New("invalid waiter ID")
	}
	if waiterShift.UserID != waiterOID {
		return nil, errors.New("unauthorized: not your shift")
	}

	// 3. Check if there's already a pending handover
	existingPending, err := s.handoverRepo.FindPendingByWaiterShift(ctx, waiterShiftID)
	if err != nil {
		return nil, err
	}
	if existingPending != nil {
		return nil, errors.New("there is already a pending handover for this shift")
	}

	// 4. Find active cashier shift
	cashierShifts, err := s.cashierShiftRepo.FindByStatus(ctx, cashier.CashierShiftOpen)
	if err != nil || len(cashierShifts) == 0 {
		return nil, errors.New("no active cashier shift found")
	}
	cashierShift := cashierShifts[0] // Get first open shift

	// 5. Amount must equal remaining cash for END_SHIFT
	handoverAmount := waiterShift.RemainingCash
	if handoverAmount <= 0 {
		return nil, errors.New("no remaining cash to handover")
	}

	// 6. Create handover record
	h := &handover.CashHandover{
		WaiterShiftID:  waiterShiftID,
		CashierShiftID: cashierShift.ID,
		WaiterID:       waiterOID,
		WaiterName:     waiterName,
		CashierID:      cashierShift.CashierID,
		CashierName:    cashierShift.CashierName,
		DeclaredAmount: handoverAmount,
		ActualAmount:   0, // Will be set by cashier
		Discrepancy:    0, // Will be calculated
		HandoverType:   handover.TypeEndShift,
		Status:         handover.StatusPending,
		WaiterNote:     req.WaiterNote,
		EndCash:        req.EndCash, // Store end cash for later use
		HandoverAt:     time.Now(),
		CreatedAt:      time.Now(),
		UpdatedAt:      time.Now(),
	}

	if err := s.handoverRepo.Create(ctx, h); err != nil {
		return nil, err
	}

	return h, nil
}

// ConfirmHandoverWithReconciliation confirms handover with actual amount and handles discrepancy
func (s *CashHandoverService) ConfirmHandoverWithReconciliation(
	ctx context.Context,
	handoverID primitive.ObjectID,
	req *handover.ConfirmHandoverRequest,
	cashierID string,
) error {
	// 1. Get handover record
	h, err := s.handoverRepo.FindByID(ctx, handoverID)
	if err != nil {
		return err
	}

	// 2. Validate cashier authorization
	cashierOID, err := primitive.ObjectIDFromHex(cashierID)
	if err != nil {
		return errors.New("invalid cashier ID")
	}
	if h.CashierID != cashierOID {
		return errors.New("unauthorized: not assigned to you")
	}

	// 3. Validate status
	if h.Status != handover.StatusPending {
		return errors.New("handover is not pending")
	}

	// 4. Validate actual_amount for CONFIRMED status
	if req.Status == handover.StatusConfirmed && req.ActualAmount == 0 {
		return errors.New("actual_amount is required when confirming handover")
	}

	// 5. Calculate discrepancy (only for CONFIRMED)
	var discrepancy float64
	if req.Status == handover.StatusConfirmed {
		discrepancy = req.ActualAmount - h.DeclaredAmount
	}

	// 6. Update handover with reconciliation data
	now := time.Now()
	h.Status = req.Status
	h.CashierNote = req.CashierNote
	h.ConfirmedAt = &now
	h.UpdatedAt = now

	if req.Status == handover.StatusConfirmed {
		h.ActualAmount = req.ActualAmount
		h.Discrepancy = discrepancy
		h.ReconciledAt = &now
	}

	// 7. Handle discrepancy if exists (only for CONFIRMED)
	if req.Status == handover.StatusConfirmed && h.HasDiscrepancy() {
		h.DiscrepancyReason = req.DiscrepancyReason
		h.DiscrepancyResponsibility = req.DiscrepancyResponsibility

		// Check if requires manager approval
		if h.RequiresManagerApproval(s.discrepancyThreshold) {
			h.RequiresApproval = true
			h.Status = handover.StatusDiscrepancy
		}

		// Create discrepancy record
		if err := s.createDiscrepancyRecord(ctx, h); err != nil {
			return err
		}
	}

	// 8. Update handover record
	if err := s.handoverRepo.Update(ctx, handoverID, h); err != nil {
		return err
	}

	// 9. If confirmed (and not requiring approval), update cash amounts
	if req.Status == handover.StatusConfirmed && !h.RequiresApproval {
		if err := s.updateCashAmounts(ctx, h); err != nil {
			return err
		}
	}

	return nil
}

// createDiscrepancyRecord creates a discrepancy record
func (s *CashHandoverService) createDiscrepancyRecord(ctx context.Context, h *handover.CashHandover) error {
	d := &handover.CashDiscrepancy{
		HandoverID:              h.ID,
		WaiterShiftID:           h.WaiterShiftID,
		CashierShiftID:          h.CashierShiftID,
		DeclaredAmount:          h.DeclaredAmount,
		ActualAmount:            h.ActualAmount,
		DiscrepancyAmount:       h.Discrepancy,
		DiscrepancyType:         h.GetDiscrepancyType(),
		DetailedReason:          h.DiscrepancyReason,
		Responsibility:          h.DiscrepancyResponsibility,
		ResolutionStatus:        "PENDING",
		RequiresManagerApproval: h.RequiresApproval,
		CreatedAt:               time.Now(),
		UpdatedAt:               time.Now(),
	}

	return s.discrepancyRepo.Create(ctx, d)
}

// updateCashAmounts updates cash amounts for both shifts
func (s *CashHandoverService) updateCashAmounts(ctx context.Context, h *handover.CashHandover) error {
	now := time.Now()

	// Update waiter shift - use actual amount received
	waiterShift, err := s.shiftRepo.FindByID(ctx, h.WaiterShiftID)
	if err != nil {
		return err
	}

	waiterShift.HandedOverCash += h.ActualAmount
	waiterShift.RemainingCash -= h.DeclaredAmount // Reduce by declared amount
	waiterShift.TotalDiscrepancy += h.Discrepancy
	waiterShift.HandoverCount++
	waiterShift.UpdatedAt = now

	// Handle END_SHIFT type
	if h.HandoverType == handover.TypeEndShift {
		// Calculate total revenue and orders
		orders, err := s.orderRepo.FindByShiftID(ctx, h.WaiterShiftID)
		if err != nil {
			return err
		}

		totalRevenue := 0.0
		for _, o := range orders {
			if o.Status == order.StatusPaid || o.Status == order.StatusInProgress || o.Status == order.StatusServed {
				totalRevenue += o.Total
			}
		}

		// End the shift
		waiterShift.Status = order.ShiftClosed
		waiterShift.EndCash = h.EndCash
		waiterShift.TotalRevenue = totalRevenue
		waiterShift.TotalOrders = len(orders)
		waiterShift.EndedAt = &now

		// Lock completed orders
		for _, o := range orders {
			if o.Status == order.StatusServed || o.Status == order.StatusCancelled {
				o.Status = order.StatusLocked
				o.LockedAt = &now
				if err := s.orderRepo.Update(ctx, o.ID, o); err != nil {
					return err
				}
			}
		}
	}

	if err := s.shiftRepo.Update(ctx, h.WaiterShiftID, waiterShift); err != nil {
		return err
	}

	// Update cashier shift
	cashierShift, err := s.cashierShiftRepo.FindByID(ctx, h.CashierShiftID)
	if err != nil {
		return err
	}

	cashierShift.ReceivedCash += h.ActualAmount
	cashierShift.TotalDiscrepancy += h.Discrepancy
	cashierShift.HandoverCount++
	if h.HasDiscrepancy() {
		cashierShift.DiscrepancyCount++
	}
	cashierShift.UpdatedAt = now

	if err := s.cashierShiftRepo.Save(ctx, cashierShift); err != nil {
		return err
	}

	return nil
}

// ApproveDiscrepancy approves or rejects a discrepancy (manager only)
func (s *CashHandoverService) ApproveDiscrepancy(
	ctx context.Context,
	handoverID primitive.ObjectID,
	managerID string,
	approved bool,
	note string,
) error {
	h, err := s.handoverRepo.FindByID(ctx, handoverID)
	if err != nil {
		return err
	}

	if !h.RequiresApproval {
		return errors.New("handover does not require approval")
	}

	now := time.Now()
	managerOID, err := primitive.ObjectIDFromHex(managerID)
	if err != nil {
		return errors.New("invalid manager ID")
	}

	h.ApprovedBy = managerOID
	h.ApprovedAt = &now
	h.UpdatedAt = now

	if approved {
		h.Status = handover.StatusConfirmed
		// Update cash amounts after approval
		if err := s.updateCashAmounts(ctx, h); err != nil {
			return err
		}
	} else {
		h.Status = handover.StatusRejected
		if h.CashierNote != "" {
			h.CashierNote += " | Manager rejected: " + note
		} else {
			h.CashierNote = "Manager rejected: " + note
		}
	}

	// Update discrepancy record
	d, err := s.discrepancyRepo.FindByHandoverID(ctx, handoverID)
	if err == nil && d != nil {
		d.ManagerApproved = approved
		d.ApprovedBy = managerOID
		d.ApprovedAt = &now
		d.ManagerNote = note
		d.ResolutionStatus = "RESOLVED"
		d.UpdatedAt = now
		if err := s.discrepancyRepo.Update(ctx, d.ID, d); err != nil {
			return err
		}
	}

	return s.handoverRepo.Update(ctx, handoverID, h)
}

// CancelHandover cancels a pending handover
func (s *CashHandoverService) CancelHandover(ctx context.Context, handoverID primitive.ObjectID, userID string) error {
	h, err := s.handoverRepo.FindByID(ctx, handoverID)
	if err != nil {
		return err
	}

	if h.Status != handover.StatusPending {
		return errors.New("can only cancel pending handovers")
	}

	// Verify user is the waiter who created it
	userOID, err := primitive.ObjectIDFromHex(userID)
	if err != nil {
		return errors.New("invalid user ID")
	}
	if h.WaiterID != userOID {
		return errors.New("unauthorized: not your handover")
	}

	return s.handoverRepo.Delete(ctx, handoverID)
}

// GetPendingHandover gets pending handover for a shift
func (s *CashHandoverService) GetPendingHandover(ctx context.Context, shiftID primitive.ObjectID) (*handover.CashHandover, error) {
	return s.handoverRepo.FindPendingByWaiterShift(ctx, shiftID)
}

// GetHandoverByID gets handover by ID
func (s *CashHandoverService) GetHandoverByID(ctx context.Context, handoverID primitive.ObjectID) (*handover.CashHandover, error) {
	return s.handoverRepo.FindByID(ctx, handoverID)
}

// GetHandoverHistory gets handover history for a shift
func (s *CashHandoverService) GetHandoverHistory(ctx context.Context, shiftID primitive.ObjectID) ([]*handover.CashHandover, error) {
	return s.handoverRepo.FindByWaiterShift(ctx, shiftID)
}

// GetPendingByCashier gets pending handovers for a cashier
func (s *CashHandoverService) GetPendingByCashier(ctx context.Context, cashierID primitive.ObjectID) ([]*handover.CashHandover, error) {
	return s.handoverRepo.FindPendingByCashier(ctx, cashierID)
}

// GetAllPending gets all pending handovers
func (s *CashHandoverService) GetAllPending(ctx context.Context) ([]*handover.CashHandover, error) {
	return s.handoverRepo.FindAllPending(ctx)
}

// GetTodayByCashier gets today's handovers for a cashier
func (s *CashHandoverService) GetTodayByCashier(ctx context.Context, cashierID primitive.ObjectID) ([]*handover.CashHandover, error) {
	return s.handoverRepo.FindTodayByCashier(ctx, cashierID)
}

// GetRequiringApproval gets handovers requiring manager approval
func (s *CashHandoverService) GetRequiringApproval(ctx context.Context) ([]*handover.CashHandover, error) {
	return s.handoverRepo.FindRequiringApproval(ctx)
}

// DiscrepancyStats represents discrepancy statistics
type DiscrepancyStats struct {
	TotalHandovers   int     `json:"total_handovers"`
	TotalDiscrepancy float64 `json:"total_discrepancy"`
	ShortageCount    int     `json:"shortage_count"`
	OverageCount     int     `json:"overage_count"`
	ShortageAmount   float64 `json:"shortage_amount"`
	OverageAmount    float64 `json:"overage_amount"`
	RequiredApproval int     `json:"required_approval"`
}

// GetDiscrepancyStats gets discrepancy statistics for a date range
func (s *CashHandoverService) GetDiscrepancyStats(ctx context.Context, startDate, endDate time.Time) (*DiscrepancyStats, error) {
	handovers, err := s.handoverRepo.FindByDateRange(ctx, startDate, endDate)
	if err != nil {
		return nil, err
	}

	stats := &DiscrepancyStats{
		TotalHandovers:   len(handovers),
		TotalDiscrepancy: 0,
		ShortageCount:    0,
		OverageCount:     0,
		ShortageAmount:   0,
		OverageAmount:    0,
		RequiredApproval: 0,
	}

	for _, h := range handovers {
		if h.HasDiscrepancy() {
			stats.TotalDiscrepancy += h.Discrepancy
			if h.Discrepancy < 0 {
				stats.ShortageCount++
				stats.ShortageAmount += -h.Discrepancy
			} else {
				stats.OverageCount++
				stats.OverageAmount += h.Discrepancy
			}
			if h.RequiresApproval {
				stats.RequiredApproval++
			}
		}
	}

	return stats, nil
}
