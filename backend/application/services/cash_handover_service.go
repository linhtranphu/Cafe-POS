package services

import (
	"context"
	"errors"
	"fmt"
	"math"
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"

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
	mongoClient          *mongo.Client
	discrepancyThreshold float64
}

// NewCashHandoverService creates a new cash handover service
func NewCashHandoverService(
	handoverRepo *mongodb.CashHandoverRepository,
	discrepancyRepo *mongodb.CashDiscrepancyRepository,
	shiftRepo *mongodb.ShiftRepository,
	cashierShiftRepo *mongodb.CashierShiftRepository,
	orderRepo *mongodb.OrderRepository,
	mongoClient *mongo.Client,
) *CashHandoverService {
	return &CashHandoverService{
		handoverRepo:         handoverRepo,
		discrepancyRepo:      discrepancyRepo,
		shiftRepo:            shiftRepo,
		cashierShiftRepo:     cashierShiftRepo,
		orderRepo:            orderRepo,
		mongoClient:          mongoClient,
		discrepancyThreshold: 100000, // 100k VND
	}
}

// ============================================================================
// CREATE HANDOVER
// ============================================================================

// CreateHandover creates a new handover request from waiter
// Supports both cash and transfer amounts
func (s *CashHandoverService) CreateHandover(
	ctx context.Context,
	waiterShiftID primitive.ObjectID,
	cashAmount, transferAmount float64,
	handoverType handover.HandoverType,
	waiterNote string,
	endCash, endTransfer float64,
	waiterID, waiterName string,
) (*handover.CashHandover, error) {
	fmt.Printf("🔵 [CREATE HANDOVER] Starting - Cash: %.0f, Transfer: %.0f\n", cashAmount, transferAmount)

	// 1. Validate at least one amount is provided
	if cashAmount == 0 && transferAmount == 0 {
		return nil, errors.New("at least one amount (cash or transfer) must be greater than 0")
	}

	// 2. Validate waiter shift exists and is open
	waiterShift, err := s.shiftRepo.FindByID(ctx, waiterShiftID)
	if err != nil {
		return nil, errors.New("waiter shift not found")
	}
	if waiterShift.Status != order.ShiftOpen {
		return nil, errors.New("waiter shift is not open")
	}

	fmt.Printf("🔵 [CREATE HANDOVER] Shift before - RemainingCash: %.0f, RemainingTransfer: %.0f\n", 
		waiterShift.RemainingCash, waiterShift.RemainingTransfer)

	// 3. Check if waiter owns the shift
	waiterOID, err := primitive.ObjectIDFromHex(waiterID)
	if err != nil {
		return nil, errors.New("invalid waiter ID")
	}
	if waiterShift.UserID != waiterOID {
		return nil, errors.New("unauthorized: not your shift")
	}

	// 4. Check if there's already a pending handover
	existingPending, err := s.handoverRepo.FindPendingByWaiterShift(ctx, waiterShiftID)
	if err != nil {
		return nil, err
	}
	if existingPending != nil {
		return nil, errors.New("there is already a pending handover for this shift")
	}

	// 5. Validate amounts against shift balances
	if cashAmount > waiterShift.RemainingCash {
		return nil, fmt.Errorf("cash amount (%.0f) exceeds remaining cash (%.0f)", cashAmount, waiterShift.RemainingCash)
	}
	if transferAmount > waiterShift.RemainingTransfer {
		return nil, fmt.Errorf("transfer amount (%.0f) exceeds remaining transfer (%.0f)", transferAmount, waiterShift.RemainingTransfer)
	}

	// 6. Find active cashier shift
	cashierShifts, err := s.cashierShiftRepo.FindByStatus(ctx, cashier.CashierShiftOpen)
	if err != nil || len(cashierShifts) == 0 {
		return nil, errors.New("no active cashier shift found")
	}
	cashierShift := cashierShifts[0]

	// 7. Create handover record
	h := &handover.CashHandover{
		WaiterShiftID:          waiterShiftID,
		CashierShiftID:         cashierShift.ID,
		WaiterID:               waiterOID,
		WaiterName:             waiterName,
		CashierID:              cashierShift.CashierID,
		CashierName:            cashierShift.CashierName,
		CashDeclaredAmount:     cashAmount,
		TransferDeclaredAmount: transferAmount,
		CashActualAmount:       0, // Will be set by cashier
		TransferActualAmount:   0, // Will be set by cashier
		CashDiscrepancy:        0, // Will be calculated
		TransferDiscrepancy:    0, // Will be calculated
		// Deprecated fields for backward compatibility
		DeclaredAmount: cashAmount + transferAmount,
		ActualAmount:   0,
		Discrepancy:    0,
		HandoverType:   handoverType,
		Status:         handover.StatusPending,
		WaiterNote:     waiterNote,
		EndCash:        endCash,
		EndTransfer:    endTransfer,
		HandoverAt:     time.Now(),
		CreatedAt:      time.Now(),
		UpdatedAt:      time.Now(),
	}

	if err := s.handoverRepo.Create(ctx, h); err != nil {
		return nil, err
	}

	// 8. CRITICAL: Reduce remaining amounts immediately
	// This is the ONLY place where remaining amounts are reduced
	waiterShift.RemainingCash -= cashAmount
	waiterShift.RemainingTransfer -= transferAmount
	waiterShift.UpdatedAt = time.Now()

	fmt.Printf("🔵 [CREATE HANDOVER] Shift after - RemainingCash: %.0f, RemainingTransfer: %.0f\n", 
		waiterShift.RemainingCash, waiterShift.RemainingTransfer)

	if err := s.shiftRepo.Update(ctx, waiterShiftID, waiterShift); err != nil {
		return nil, err
	}

	fmt.Printf("✅ [CREATE HANDOVER] Success - Handover ID: %s\n", h.ID.Hex())
	return h, nil
}

// ============================================================================
// CONFIRM HANDOVER
// ============================================================================

// ConfirmHandover confirms or rejects a handover with actual amounts
// Uses MongoDB transaction to ensure atomicity - if rejected, rollback remaining amounts
func (s *CashHandoverService) ConfirmHandover(
	ctx context.Context,
	handoverID primitive.ObjectID,
	actualCashAmount, actualTransferAmount float64,
	status handover.HandoverStatus,
	cashierNote, discrepancyReason string,
	discrepancyResponsibility handover.ResponsibilityType,
	cashierID string,
) error {
	fmt.Printf("🟢 [CONFIRM HANDOVER] Starting - ActualCash: %.0f, ActualTransfer: %.0f, Status: %s\n", 
		actualCashAmount, actualTransferAmount, status)

	// 1. Get handover record
	h, err := s.handoverRepo.FindByID(ctx, handoverID)
	if err != nil {
		return err
	}

	fmt.Printf("🟢 [CONFIRM HANDOVER] Handover - DeclaredCash: %.0f, DeclaredTransfer: %.0f\n", 
		h.CashDeclaredAmount, h.TransferDeclaredAmount)

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

	// 4. Validate actual amounts match declared amounts
	if status == handover.StatusConfirmed {
		if h.CashDeclaredAmount > 0 && actualCashAmount == 0 {
			return errors.New("actual_cash_amount is required when cash was declared")
		}
		if h.TransferDeclaredAmount > 0 && actualTransferAmount == 0 {
			return errors.New("actual_transfer_amount is required when transfer was declared")
		}
	}

	// 5. Calculate discrepancies
	cashDiscrepancy := actualCashAmount - h.CashDeclaredAmount
	transferDiscrepancy := actualTransferAmount - h.TransferDeclaredAmount
	totalDiscrepancy := cashDiscrepancy + transferDiscrepancy

	fmt.Printf("🟢 [CONFIRM HANDOVER] Discrepancy - Cash: %.0f, Transfer: %.0f, Total: %.0f\n", 
		cashDiscrepancy, transferDiscrepancy, totalDiscrepancy)

	// 6. START MONGODB TRANSACTION
	// This ensures atomicity - either all updates succeed or all rollback
	session, err := s.mongoClient.StartSession()
	if err != nil {
		return fmt.Errorf("failed to start transaction session: %w", err)
	}
	defer session.EndSession(ctx)

	// Execute transaction
	_, err = session.WithTransaction(ctx, func(sessCtx mongo.SessionContext) (interface{}, error) {
		// 6a. Update handover record
		now := time.Now()
		h.Status = status
		h.CashierNote = cashierNote
		h.ConfirmedAt = &now
		h.UpdatedAt = now

		if status == handover.StatusConfirmed {
			h.CashActualAmount = actualCashAmount
			h.TransferActualAmount = actualTransferAmount
			h.CashDiscrepancy = cashDiscrepancy
			h.TransferDiscrepancy = transferDiscrepancy
			// Update deprecated fields
			h.ActualAmount = actualCashAmount + actualTransferAmount
			h.Discrepancy = totalDiscrepancy
			h.ReconciledAt = &now
		}

		// 6b. Handle discrepancy if exists
		if status == handover.StatusConfirmed && totalDiscrepancy != 0 {
			h.DiscrepancyReason = discrepancyReason
			h.DiscrepancyResponsibility = discrepancyResponsibility

			// Check if requires manager approval
			totalAbsDiscrepancy := math.Abs(totalDiscrepancy)
			if totalAbsDiscrepancy > s.discrepancyThreshold {
				h.RequiresApproval = true
				h.Status = handover.StatusDiscrepancy
				fmt.Printf("⚠️  [CONFIRM HANDOVER] Large discrepancy - requires manager approval\n")
			}

			// Create discrepancy record
			if err := s.createDiscrepancyRecord(sessCtx, h); err != nil {
				return nil, err
			}
		}

		// 6c. Update handover record
		if err := s.handoverRepo.Update(sessCtx, handoverID, h); err != nil {
			return nil, err
		}

		// 6d. Handle REJECTION - Rollback remaining amounts
		if status == handover.StatusRejected {
			fmt.Printf("🔴 [CONFIRM HANDOVER] REJECTED - Rolling back remaining amounts\n")
			
			// Get waiter shift
			waiterShift, err := s.shiftRepo.FindByID(sessCtx, h.WaiterShiftID)
			if err != nil {
				return nil, err
			}

			fmt.Printf("🔴 [ROLLBACK] Before - RemainingCash: %.0f, RemainingTransfer: %.0f\n", 
				waiterShift.RemainingCash, waiterShift.RemainingTransfer)

			// RESTORE the amounts that were deducted in CreateHandover
			waiterShift.RemainingCash += h.CashDeclaredAmount
			waiterShift.RemainingTransfer += h.TransferDeclaredAmount
			waiterShift.UpdatedAt = now

			fmt.Printf("🔴 [ROLLBACK] After - RemainingCash: %.0f, RemainingTransfer: %.0f\n", 
				waiterShift.RemainingCash, waiterShift.RemainingTransfer)

			// Update waiter shift
			if err := s.shiftRepo.Update(sessCtx, h.WaiterShiftID, waiterShift); err != nil {
				return nil, err
			}

			fmt.Printf("✅ [ROLLBACK] Success - Amounts restored\n")
		}

		// 6e. If confirmed and not requiring approval, update balances
		if status == handover.StatusConfirmed && !h.RequiresApproval {
			if err := s.updateBalances(sessCtx, h); err != nil {
				return nil, err
			}
		}

		return nil, nil
	})

	if err != nil {
		fmt.Printf("❌ [CONFIRM HANDOVER] Transaction failed: %v\n", err)
		return err
	}

	fmt.Printf("✅ [CONFIRM HANDOVER] Success\n")
	return nil
}

// ============================================================================
// UPDATE BALANCES
// ============================================================================

// updateBalances updates handed_over amounts for both shifts
// CRITICAL: This function does NOT reduce remaining amounts
// Remaining amounts were already reduced in CreateHandover
func (s *CashHandoverService) updateBalances(ctx context.Context, h *handover.CashHandover) error {
	fmt.Printf("🟡 [UPDATE BALANCES] Starting\n")
	now := time.Now()

	// Update waiter shift
	waiterShift, err := s.shiftRepo.FindByID(ctx, h.WaiterShiftID)
	if err != nil {
		return err
	}

	fmt.Printf("🟡 [UPDATE BALANCES] Waiter shift before - HandedOverCash: %.0f, HandedOverTransfer: %.0f, RemainingCash: %.0f, RemainingTransfer: %.0f\n", 
		waiterShift.HandedOverCash, waiterShift.HandedOverTransfer, waiterShift.RemainingCash, waiterShift.RemainingTransfer)

	// ONLY update handed_over amounts - DO NOT touch remaining amounts
	waiterShift.HandedOverCash += h.CashActualAmount
	waiterShift.HandedOverTransfer += h.TransferActualAmount
	waiterShift.TotalDiscrepancy += h.TotalDiscrepancy()
	waiterShift.HandoverCount++
	waiterShift.UpdatedAt = now

	fmt.Printf("🟡 [UPDATE BALANCES] Waiter shift after - HandedOverCash: %.0f, HandedOverTransfer: %.0f, RemainingCash: %.0f, RemainingTransfer: %.0f\n", 
		waiterShift.HandedOverCash, waiterShift.HandedOverTransfer, waiterShift.RemainingCash, waiterShift.RemainingTransfer)

	// Handle END_SHIFT type
	if h.HandoverType == handover.TypeEndShift {
		fmt.Printf("🟡 [UPDATE BALANCES] Ending shift\n")
		orders, err := s.orderRepo.FindByShiftID(ctx, h.WaiterShiftID)
		if err != nil {
			return err
		}

		cashRevenue := 0.0
		transferRevenue := 0.0
		for _, o := range orders {
			if o.Status == order.StatusPaid || o.Status == order.StatusInProgress || o.Status == order.StatusServed {
				if o.PaymentMethod == order.PaymentCash {
					cashRevenue += o.Total
				} else if o.PaymentMethod == order.PaymentTransfer || o.PaymentMethod == order.PaymentQR {
					transferRevenue += o.Total
				}
			}
		}

		waiterShift.Status = order.ShiftClosed
		waiterShift.EndCash = h.EndCash
		waiterShift.EndTransfer = h.EndTransfer
		waiterShift.TotalRevenue = cashRevenue + transferRevenue
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

	fmt.Printf("🟡 [UPDATE BALANCES] Cashier shift before - ReceivedCash: %.0f, ReceivedTransfer: %.0f\n", 
		cashierShift.ReceivedCash, cashierShift.ReceivedTransfer)

	cashierShift.ReceivedCash += h.CashActualAmount
	cashierShift.ReceivedTransfer += h.TransferActualAmount
	cashierShift.TotalDiscrepancy += h.TotalDiscrepancy()
	cashierShift.HandoverCount++
	if h.TotalDiscrepancy() != 0 {
		cashierShift.DiscrepancyCount++
	}
	cashierShift.UpdatedAt = now

	fmt.Printf("🟡 [UPDATE BALANCES] Cashier shift after - ReceivedCash: %.0f, ReceivedTransfer: %.0f\n", 
		cashierShift.ReceivedCash, cashierShift.ReceivedTransfer)

	if err := s.cashierShiftRepo.Save(ctx, cashierShift); err != nil {
		return err
	}

	fmt.Printf("✅ [UPDATE BALANCES] Success\n")
	return nil
}

// ============================================================================
// HELPER FUNCTIONS
// ============================================================================

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

// ============================================================================
// QUERY FUNCTIONS
// ============================================================================

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

	// IMPORTANT: Restore remaining amounts when canceling
	waiterShift, err := s.shiftRepo.FindByID(ctx, h.WaiterShiftID)
	if err != nil {
		return err
	}

	waiterShift.RemainingCash += h.CashDeclaredAmount
	waiterShift.RemainingTransfer += h.TransferDeclaredAmount
	waiterShift.UpdatedAt = time.Now()

	if err := s.shiftRepo.Update(ctx, h.WaiterShiftID, waiterShift); err != nil {
		return err
	}

	return s.handoverRepo.Delete(ctx, handoverID)
}
