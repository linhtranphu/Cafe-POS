package services

import (
	"context"
	"errors"

	"cafe-pos/backend/domain"
	"cafe-pos/backend/domain/cashier"
	"cafe-pos/backend/domain/order"
	"cafe-pos/backend/infrastructure/mongodb"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

// CashierShiftService handles business logic for cashier shifts.
// This is separate from the regular shift service to maintain clear separation
// between cashier shifts and waiter/barista shifts.
type CashierShiftService struct {
	cashierShiftRepo    *mongodb.CashierShiftRepository
	waiterShiftRepo     ShiftRepository // To check if waiter shifts are closed
	stateMachineManager *domain.StateMachineManager
}

// NewCashierShiftService creates a new service for managing cashier shifts.
func NewCashierShiftService(
	cashierShiftRepo *mongodb.CashierShiftRepository,
	waiterShiftRepo ShiftRepository,
	stateMachineManager *domain.StateMachineManager,
) *CashierShiftService {
	return &CashierShiftService{
		cashierShiftRepo:    cashierShiftRepo,
		waiterShiftRepo:     waiterShiftRepo,
		stateMachineManager: stateMachineManager,
	}
}

// StartCashierShift creates a new cashier shift.
// It validates that the cashier doesn't already have an open shift.
//
// Parameters:
//   - ctx: Context for the operation
//   - cashierID: The ID of the cashier starting the shift
//   - cashierName: The name of the cashier for display
//   - startingFloat: The initial cash amount in the drawer
//
// Returns:
//   - *cashier.CashierShift: The created shift
//   - error: If validation fails or creation fails
func (s *CashierShiftService) StartCashierShift(
	ctx context.Context,
	cashierID primitive.ObjectID,
	cashierName string,
	startingFloat float64,
) (*cashier.CashierShift, error) {
	// Check if cashier already has an open shift
	existingShift, err := s.cashierShiftRepo.FindOpenByCashier(ctx, cashierID)
	if err != nil {
		return nil, err
	}
	if existingShift != nil {
		return nil, errors.New("cashier already has an open shift")
	}

	// Validate starting float is non-negative
	if startingFloat < 0 {
		return nil, errors.New("starting float must be non-negative")
	}

	// Create new cashier shift
	shift := cashier.NewCashierShift(cashierID, cashierName, startingFloat)

	// Save to repository
	if err := s.cashierShiftRepo.Create(ctx, shift); err != nil {
		return nil, err
	}

	return shift, nil
}

// GetCurrentCashierShift retrieves the current open cashier shift for a cashier.
//
// Parameters:
//   - ctx: Context for the operation
//   - cashierID: The ID of the cashier
//
// Returns:
//   - *cashier.CashierShift: The current open shift, or nil if none exists
//   - error: If query fails
func (s *CashierShiftService) GetCurrentCashierShift(
	ctx context.Context,
	cashierID primitive.ObjectID,
) (*cashier.CashierShift, error) {
	return s.cashierShiftRepo.FindOpenByCashier(ctx, cashierID)
}

// GetCashierShift retrieves a specific cashier shift by ID.
func (s *CashierShiftService) GetCashierShift(
	ctx context.Context,
	shiftID primitive.ObjectID,
) (*cashier.CashierShift, error) {
	return s.cashierShiftRepo.FindByID(ctx, shiftID)
}

// GetCashierShiftsByUser retrieves all cashier shifts for a specific cashier.
func (s *CashierShiftService) GetCashierShiftsByUser(
	ctx context.Context,
	cashierID primitive.ObjectID,
) ([]*cashier.CashierShift, error) {
	return s.cashierShiftRepo.FindByCashierID(ctx, cashierID)
}

// GetAllCashierShifts retrieves all cashier shifts.
func (s *CashierShiftService) GetAllCashierShifts(ctx context.Context) ([]*cashier.CashierShift, error) {
	return s.cashierShiftRepo.FindAll(ctx)
}

// CheckWaiterShiftsClosed checks if all waiter shifts are closed.
// This is used during cashier shift closure to ensure all waiter shifts are closed first.
//
// Returns:
//   - bool: true if all waiter shifts are closed, false otherwise
//   - []*order.Shift: List of open waiter shifts (if any)
//   - error: If query fails
func (s *CashierShiftService) CheckWaiterShiftsClosed(ctx context.Context) (bool, []*order.Shift, error) {
	// Get all open shifts
	openShifts, err := s.waiterShiftRepo.FindOpenShifts(ctx)
	if err != nil {
		return false, nil, err
	}

	// Filter for waiter shifts only (exclude barista shifts if needed)
	var openWaiterShifts []*order.Shift
	for _, shift := range openShifts {
		if shift.RoleType == order.RoleWaiter {
			openWaiterShifts = append(openWaiterShifts, shift)
		}
	}

	return len(openWaiterShifts) == 0, openWaiterShifts, nil
}

// SaveCashierShift saves an updated cashier shift.
// This is used during the shift closure workflow to persist state changes.
func (s *CashierShiftService) SaveCashierShift(
	ctx context.Context,
	shift *cashier.CashierShift,
) error {
	return s.cashierShiftRepo.Save(ctx, shift)
}

// CanCloseCashierShift checks if a cashier shift can be closed.
// It verifies that all waiter shifts are closed before allowing cashier shift closure.
//
// Parameters:
//   - ctx: Context for the operation
//   - shiftID: The ID of the cashier shift to close
//
// Returns:
//   - bool: true if the shift can be closed, false otherwise
//   - error: If validation fails
func (s *CashierShiftService) CanCloseCashierShift(
	ctx context.Context,
	shiftID primitive.ObjectID,
) (bool, error) {
	// Check if all waiter shifts are closed
	allClosed, _, err := s.CheckWaiterShiftsClosed(ctx)
	if err != nil {
		return false, err
	}

	if !allClosed {
		return false, errors.New("cannot close cashier shift: " + 
			"there are still open waiter shifts")
	}

	// Get the cashier shift
	shift, err := s.cashierShiftRepo.FindByID(ctx, shiftID)
	if err != nil {
		return false, err
	}

	// Validate using state machine
	if err := s.stateMachineManager.ValidateCashierShiftTransition(shift, cashier.EventCloseShift); err != nil {
		return false, err
	}

	return true, nil
}
