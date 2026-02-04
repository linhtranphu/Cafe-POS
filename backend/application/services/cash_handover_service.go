package services

import (
	"context"
	"errors"
	"fmt"
	"time"

	"cafe-pos/backend/domain/handover"
	"cafe-pos/backend/infrastructure/mongodb"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

// Configuration constants
const (
	DefaultDiscrepancyThreshold = 50000.0 // VND - Large discrepancies requiring manager approval
	MaxHandoverAmount          = 10000000.0 // VND - Maximum amount per handover
)

type CashHandoverService struct {
	handoverRepo    *mongodb.CashHandoverRepository
	discrepancyRepo *mongodb.CashDiscrepancyRepository
	shiftRepo       *mongodb.ShiftRepository
	cashierShiftRepo *mongodb.CashierShiftRepository
	discrepancyThreshold float64
}

func NewCashHandoverService(
	handoverRepo *mongodb.CashHandoverRepository,
	discrepancyRepo *mongodb.CashDiscrepancyRepository,
	shiftRepo *mongodb.ShiftRepository,
	cashierShiftRepo *mongodb.CashierShiftRepository,
) *CashHandoverService {
	return &CashHandoverService{
		handoverRepo:         handoverRepo,
		discrepancyRepo:      discrepancyRepo,
		shiftRepo:            shiftRepo,
		cashierShiftRepo:     cashierShiftRepo,
		discrepancyThreshold: DefaultDiscrepancyThreshold,
	}
}

// SetDiscrepancyThreshold sets the threshold for requiring manager approval
func (s *CashHandoverService) SetDiscrepancyThreshold(threshold float64) {
	s.discrepancyThreshold = threshold
}

// CreateHandover creates a new cash handover request
func (s *CashHandoverService) CreateHandover(
	ctx context.Context,
	waiterShiftID primitive.ObjectID,
	req *handover.CreateHandoverRequest,
	waiterID primitive.ObjectID,
	waiterName string,
) (*handover.CashHandover, error) {
	// Validate request
	if req.RequestedAmount <= 0 {
		return nil, errors.New("requested amount must be greater than 0")
	}

	if req.RequestedAmount > MaxHandoverAmount {
		return nil, fmt.Errorf("requested amount exceeds maximum limit of %.0f", MaxHandoverAmount)
	}

	// Get waiter shift to validate ownership and cash availability
	shift, err := s.shiftRepo.FindByID(ctx, waiterShiftID)
	if err != nil {
		return nil, fmt.Errorf("failed to find shift: %w", err)
	}

	if shift == nil {
		return nil, errors.New("shift not found")
	}

	// Validate shift ownership
	if shift.UserID != waiterID {
		return nil, errors.New("unauthorized: shift does not belong to user")
	}

	// Check if shift can perform handover
	if err := shift.CanHandover(req.RequestedAmount); err != nil {
		return nil, fmt.Errorf("handover validation failed: %w", err)
	}

	// Check for existing pending handover
	existingHandover, err := s.handoverRepo.FindPendingByWaiterShift(ctx, waiterShiftID)
	if err != nil {
		return nil, fmt.Errorf("failed to check existing handovers: %w", err)
	}

	if existingHandover != nil {
		return nil, errors.New("there is already a pending handover for this shift")
	}

	// Create handover
	handover, err := handover.NewCashHandover(
		req.Type,
		waiterShiftID,
		waiterID,
		waiterName,
		req.RequestedAmount,
		req.WaiterNotes,
	)
	if err != nil {
		return nil, fmt.Errorf("failed to create handover: %w", err)
	}

	// Save handover
	if err := s.handoverRepo.Create(ctx, handover); err != nil {
		return nil, fmt.Errorf("failed to save handover: %w", err)
	}

	return handover, nil
}

// CreateHandoverAndEndShift creates a handover and ends the shift
func (s *CashHandoverService) CreateHandoverAndEndShift(
	ctx context.Context,
	waiterShiftID primitive.ObjectID,
	req *handover.CreateHandoverRequest,
	waiterID primitive.ObjectID,
	waiterName string,
) (*handover.CashHandover, error) {
	// Ensure this is an END_SHIFT type handover
	req.Type = handover.HandoverTypeEndShift

	// Create the handover
	handover, err := s.CreateHandover(ctx, waiterShiftID, req, waiterID, waiterName)
	if err != nil {
		return nil, err
	}

	// Note: The actual shift ending will be handled when the handover is confirmed
	// This prevents the shift from being ended if the handover is rejected

	return handover, nil
}

// ConfirmHandoverWithReconciliation confirms a handover with cash reconciliation
func (s *CashHandoverService) ConfirmHandoverWithReconciliation(
	ctx context.Context,
	handoverID primitive.ObjectID,
	req *handover.ConfirmHandoverRequest,
	cashierShiftID, cashierID primitive.ObjectID,
	cashierName string,
) error {
	// Get handover
	handover, err := s.handoverRepo.FindByID(ctx, handoverID)
	if err != nil {
		return fmt.Errorf("failed to find handover: %w", err)
	}

	if handover == nil {
		return errors.New("handover not found")
	}

	// Validate cashier shift exists and is open
	cashierShift, err := s.cashierShiftRepo.FindByID(ctx, cashierShiftID)
	if err != nil {
		return fmt.Errorf("failed to find cashier shift: %w", err)
	}

	if cashierShift == nil {
		return errors.New("cashier shift not found")
	}

	if cashierShift.Status != "OPEN" {
		return errors.New("cashier shift is not open")
	}

	// Confirm handover
	if err := handover.ConfirmHandover(
		cashierShiftID,
		cashierID,
		cashierName,
		req.ActualAmount,
		req.CashierNotes,
		s.discrepancyThreshold,
	); err != nil {
		return fmt.Errorf("failed to confirm handover: %w", err)
	}

	// Update handover in database
	if err := s.handoverRepo.Update(ctx, handoverID, handover); err != nil {
		return fmt.Errorf("failed to update handover: %w", err)
	}

	// Update cash amounts in shifts
	if err := s.updateCashAmounts(ctx, handover); err != nil {
		return fmt.Errorf("failed to update cash amounts: %w", err)
	}

	// Create discrepancy record if needed
	if handover.HasDiscrepancy() {
		if err := s.createDiscrepancyRecord(ctx, handover); err != nil {
			return fmt.Errorf("failed to create discrepancy record: %w", err)
		}
	}

	// If this was an END_SHIFT handover and it's confirmed, end the waiter shift
	if handover.Type == handover.HandoverTypeEndShift && handover.Status == handover.HandoverStatusConfirmed {
		if err := s.endWaiterShift(ctx, handover.WaiterShiftID); err != nil {
			return fmt.Errorf("failed to end waiter shift: %w", err)
		}
	}

	return nil
}

// RejectHandover rejects a handover request
func (s *CashHandoverService) RejectHandover(
	ctx context.Context,
	handoverID primitive.ObjectID,
	reason string,
	cashierShiftID, cashierID primitive.ObjectID,
	cashierName string,
) error {
	// Get handover
	handover, err := s.handoverRepo.FindByID(ctx, handoverID)
	if err != nil {
		return fmt.Errorf("failed to find handover: %w", err)
	}

	if handover == nil {
		return errors.New("handover not found")
	}

	// Reject handover
	if err := handover.RejectHandover(cashierShiftID, cashierID, cashierName, reason); err != nil {
		return fmt.Errorf("failed to reject handover: %w", err)
	}

	// Update handover in database
	if err := s.handoverRepo.Update(ctx, handoverID, handover); err != nil {
		return fmt.Errorf("failed to update handover: %w", err)
	}

	return nil
}

// ApproveDiscrepancy handles manager approval/rejection of discrepancies
func (s *CashHandoverService) ApproveDiscrepancy(
	ctx context.Context,
	handoverID primitive.ObjectID,
	managerID primitive.ObjectID,
	approved bool,
	managerNotes string,
) error {
	// Get handover
	handover, err := s.handoverRepo.FindByID(ctx, handoverID)
	if err != nil {
		return fmt.Errorf("failed to find handover: %w", err)
	}

	if handover == nil {
		return errors.New("handover not found")
	}

	// Approve/reject discrepancy
	if err := handover.ApproveDiscrepancy(managerID, approved, managerNotes); err != nil {
		return fmt.Errorf("failed to approve discrepancy: %w", err)
	}

	// Update handover in database
	if err := s.handoverRepo.Update(ctx, handoverID, handover); err != nil {
		return fmt.Errorf("failed to update handover: %w", err)
	}

	// Update discrepancy record if exists
	discrepancy, err := s.discrepancyRepo.FindByHandoverID(ctx, handoverID)
	if err != nil {
		return fmt.Errorf("failed to find discrepancy: %w", err)
	}

	if discrepancy != nil {
		if approved {
			if err := discrepancy.Resolve("Manager approved discrepancy"); err != nil {
				return fmt.Errorf("failed to resolve discrepancy: %w", err)
			}
		}

		if err := s.discrepancyRepo.Update(ctx, discrepancy.ID, discrepancy); err != nil {
			return fmt.Errorf("failed to update discrepancy: %w", err)
		}
	}

	return nil
}

// GetDiscrepancyStats gets discrepancy statistics for a date range
func (s *CashHandoverService) GetDiscrepancyStats(ctx context.Context, startDate, endDate time.Time) (*handover.DiscrepancyStats, error) {
	return s.discrepancyRepo.GetDiscrepancyStats(ctx, startDate, endDate)
}

// CancelHandover cancels a pending handover
func (s *CashHandoverService) CancelHandover(ctx context.Context, handoverID, waiterID primitive.ObjectID) error {
	// Get handover
	handover, err := s.handoverRepo.FindByID(ctx, handoverID)
	if err != nil {
		return fmt.Errorf("failed to find handover: %w", err)
	}

	if handover == nil {
		return errors.New("handover not found")
	}

	// Validate ownership
	if handover.WaiterID != waiterID {
		return errors.New("unauthorized: handover does not belong to user")
	}

	// Check if can cancel
	if !handover.CanCancel() {
		return errors.New("cannot cancel handover: handover is not pending")
	}

	// Delete handover
	if err := s.handoverRepo.Delete(ctx, handoverID); err != nil {
		return fmt.Errorf("failed to delete handover: %w", err)
	}

	return nil
}

// GetPendingHandovers gets pending handovers for a cashier
func (s *CashHandoverService) GetPendingHandovers(ctx context.Context, cashierID primitive.ObjectID) ([]*handover.CashHandover, error) {
	return s.handoverRepo.FindPendingByCashier(ctx, cashierID)
}

// GetTodayHandovers gets today's handovers
func (s *CashHandoverService) GetTodayHandovers(ctx context.Context) ([]*handover.CashHandover, error) {
	return s.handoverRepo.FindTodayHandovers(ctx)
}

// GetHandoverHistory gets handover history for a shift
func (s *CashHandoverService) GetHandoverHistory(ctx context.Context, shiftID primitive.ObjectID) ([]*handover.CashHandover, error) {
	return s.handoverRepo.FindByWaiterShift(ctx, shiftID)
}

// GetPendingApprovals gets handovers requiring manager approval
func (s *CashHandoverService) GetPendingApprovals(ctx context.Context) ([]*handover.CashHandover, error) {
	return s.handoverRepo.FindRequiringApproval(ctx)
}

// Private helper methods

// createDiscrepancyRecord creates a discrepancy record for tracking
func (s *CashHandoverService) createDiscrepancyRecord(ctx context.Context, handover *handover.CashHandover) error {
	if !handover.HasDiscrepancy() {
		return nil
	}

	// Determine responsibility (simplified logic - can be enhanced)
	responsibility := handover.ResponsibilityUnknown
	if handover.DiscrepancyReason != nil {
		// This would be set during reconciliation
		if handover.Responsibility != nil {
			responsibility = *handover.Responsibility
		}
	}

	reason := "Cash discrepancy during handover"
	if handover.DiscrepancyReason != nil {
		reason = *handover.DiscrepancyReason
	}

	discrepancy, err := handover.NewCashDiscrepancy(
		handover.ID,
		handover.RequestedAmount,
		*handover.ActualAmount,
		responsibility,
		reason,
		handover.WaiterID,
		handover.WaiterName,
		*handover.CashierID,
		*handover.CashierName,
	)
	if err != nil {
		return err
	}

	return s.discrepancyRepo.Create(ctx, discrepancy)
}

// updateCashAmounts updates cash amounts in both waiter and cashier shifts
func (s *CashHandoverService) updateCashAmounts(ctx context.Context, handover *handover.CashHandover) error {
	// Update waiter shift
	waiterShift, err := s.shiftRepo.FindByID(ctx, handover.WaiterShiftID)
	if err != nil {
		return fmt.Errorf("failed to find waiter shift: %w", err)
	}

	if waiterShift != nil {
		actualAmount := handover.RequestedAmount
		if handover.ActualAmount != nil {
			actualAmount = *handover.ActualAmount
		}

		discrepancyAmount := handover.GetDiscrepancyAmount()
		waiterShift.UpdateCashAfterHandover(actualAmount, discrepancyAmount)

		if err := s.shiftRepo.Update(ctx, handover.WaiterShiftID, waiterShift); err != nil {
			return fmt.Errorf("failed to update waiter shift: %w", err)
		}
	}

	// Update cashier shift
	if handover.CashierShiftID != nil {
		cashierShift, err := s.cashierShiftRepo.FindByID(ctx, *handover.CashierShiftID)
		if err != nil {
			return fmt.Errorf("failed to find cashier shift: %w", err)
		}

		if cashierShift != nil {
			actualAmount := handover.RequestedAmount
			if handover.ActualAmount != nil {
				actualAmount = *handover.ActualAmount
			}

			discrepancyAmount := handover.GetDiscrepancyAmount()
			cashierShift.UpdateCashAfterHandover(actualAmount, discrepancyAmount, handover.HasDiscrepancy())

			if err := s.cashierShiftRepo.Update(ctx, *handover.CashierShiftID, cashierShift); err != nil {
				return fmt.Errorf("failed to update cashier shift: %w", err)
			}
		}
	}

	return nil
}

// endWaiterShift ends a waiter shift after successful END_SHIFT handover
func (s *CashHandoverService) endWaiterShift(ctx context.Context, shiftID primitive.ObjectID) error {
	shift, err := s.shiftRepo.FindByID(ctx, shiftID)
	if err != nil {
		return fmt.Errorf("failed to find shift: %w", err)
	}

	if shift == nil {
		return errors.New("shift not found")
	}

	// End the shift
	now := time.Now()
	shift.Status = "CLOSED"
	shift.EndedAt = &now
	shift.EndCash = shift.RemainingCash // Should be 0 after full handover
	shift.UpdatedAt = now

	return s.shiftRepo.Update(ctx, shiftID, shift)
}