package services

import (
	"context"
	"errors"
	"fmt"
	"log"
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"

	"cafe-pos/backend/domain"
	"cafe-pos/backend/domain/cashier"
	"cafe-pos/backend/domain/order"
	"cafe-pos/backend/infrastructure/mongodb"
)

// CashierShiftService handles business logic for cashier shifts.
// This is separate from the regular shift service to maintain clear separation
// between cashier shifts and waiter/barista shifts.
type CashierShiftService struct {
	cashierShiftRepo    *mongodb.CashierShiftRepository
	fundHandoverRepo    *mongodb.FundHandoverRepository
	waiterShiftRepo     ShiftRepository // To check if waiter shifts are closed
	stateMachineManager *domain.StateMachineManager
	fundService         *FundService
	mongoClient         *mongo.Client
}

// NewCashierShiftService creates a new service for managing cashier shifts.
func NewCashierShiftService(
	cashierShiftRepo *mongodb.CashierShiftRepository,
	fundHandoverRepo *mongodb.FundHandoverRepository,
	waiterShiftRepo ShiftRepository,
	stateMachineManager *domain.StateMachineManager,
	fundService *FundService,
	mongoClient *mongo.Client,
) *CashierShiftService {
	return &CashierShiftService{
		cashierShiftRepo:    cashierShiftRepo,
		fundHandoverRepo:    fundHandoverRepo,
		waiterShiftRepo:     waiterShiftRepo,
		stateMachineManager: stateMachineManager,
		fundService:         fundService,
		mongoClient:         mongoClient,
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

	// If starting float > 0, create fund withdrawal transaction
	// Use MongoDB transaction for atomicity
	if startingFloat > 0 {
		session, err := s.mongoClient.StartSession()
		if err != nil {
			return nil, fmt.Errorf("failed to start session: %w", err)
		}
		defer session.EndSession(ctx)

		err = mongo.WithSession(ctx, session, func(sc mongo.SessionContext) error {
			if err := session.StartTransaction(); err != nil {
				return err
			}

			// Create cashier shift first to get the ID
			if err := s.cashierShiftRepo.Create(sc, shift); err != nil {
				session.AbortTransaction(sc)
				return fmt.Errorf("failed to create cashier shift: %w", err)
			}

			// Create fund withdrawal for starting float
			reason := fmt.Sprintf("Tiền đầu ca cho thu ngân %s", cashierName)
			_, _, err := s.fundService.CreateWithdrawal(
				sc,
				startingFloat, // cash amount
				0,             // transfer amount
				reason,
				cashierID,
				cashierName,
				"cashier",
			)
			if err != nil {
				session.AbortTransaction(sc)
				return fmt.Errorf("failed to create fund withdrawal for starting float: %w", err)
			}

			return session.CommitTransaction(sc)
		})

		if err != nil {
			return nil, err
		}
	} else {
		// No starting float, just create the shift
		if err := s.cashierShiftRepo.Create(ctx, shift); err != nil {
			return nil, err
		}
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


// ============================================================================
// TRANSACTION-WRAPPED CLOSURE OPERATIONS
// ============================================================================

// InitiateClosureWithTransaction wraps InitiateClosure in a MongoDB transaction
func (s *CashierShiftService) InitiateClosureWithTransaction(
	ctx context.Context,
	shiftID primitive.ObjectID,
	userID, deviceID string,
) (*cashier.CashierShift, error) {
	fmt.Printf("🔵 [INITIATE CLOSURE TX] Starting for shift %s\n", shiftID.Hex())

	// Start MongoDB transaction
	session, err := s.mongoClient.StartSession()
	if err != nil {
		return nil, fmt.Errorf("failed to start transaction session: %w", err)
	}
	defer session.EndSession(ctx)

	var updatedShift *cashier.CashierShift

	// Execute transaction
	_, err = session.WithTransaction(ctx, func(sessCtx mongo.SessionContext) (interface{}, error) {
		// 1. Get the shift
		shift, err := s.cashierShiftRepo.FindByID(sessCtx, shiftID)
		if err != nil {
			return nil, err
		}

		// 2. Validate state transition
		if err := s.stateMachineManager.ValidateCashierShiftTransition(shift, cashier.EventInitiateClosure); err != nil {
			return nil, err
		}

		// 3. Initiate closure
		if err := shift.InitiateClosure(userID, deviceID, time.Now()); err != nil {
			return nil, err
		}

		// 4. Save the shift
		if err := s.cashierShiftRepo.Save(sessCtx, shift); err != nil {
			return nil, err
		}

		updatedShift = shift
		return nil, nil
	})

	if err != nil {
		fmt.Printf("❌ [INITIATE CLOSURE TX] Failed: %v\n", err)
		return nil, err
	}

	fmt.Printf("✅ [INITIATE CLOSURE TX] Success\n")
	return updatedShift, nil
}

// RecordActualCashWithTransaction wraps RecordActualCash in a MongoDB transaction
func (s *CashierShiftService) RecordActualCashWithTransaction(
	ctx context.Context,
	shiftID primitive.ObjectID,
	actualCash float64,
	userID, deviceID string,
) (*cashier.CashierShift, *cashier.Variance, error) {
	fmt.Printf("🔵 [RECORD CASH TX] Starting for shift %s, actual: %.0f\n", shiftID.Hex(), actualCash)

	// Start MongoDB transaction
	session, err := s.mongoClient.StartSession()
	if err != nil {
		return nil, nil, fmt.Errorf("failed to start transaction session: %w", err)
	}
	defer session.EndSession(ctx)

	var updatedShift *cashier.CashierShift
	var variance *cashier.Variance

	// Execute transaction
	_, err = session.WithTransaction(ctx, func(sessCtx mongo.SessionContext) (interface{}, error) {
		// 1. Get the shift
		shift, err := s.cashierShiftRepo.FindByID(sessCtx, shiftID)
		if err != nil {
			return nil, err
		}

		// 2. Validate step
		if err := s.stateMachineManager.ValidateCashierShiftStep(shift, "record_actual_cash"); err != nil {
			return nil, err
		}

		// 3. Record actual cash
		v, err := shift.RecordActualCash(actualCash, userID, deviceID, time.Now())
		if err != nil {
			return nil, err
		}
		variance = v

		// 4. Save the shift
		if err := s.cashierShiftRepo.Save(sessCtx, shift); err != nil {
			return nil, err
		}

		updatedShift = shift
		return nil, nil
	})

	if err != nil {
		fmt.Printf("❌ [RECORD CASH TX] Failed: %v\n", err)
		return nil, nil, err
	}

	fmt.Printf("✅ [RECORD CASH TX] Success, variance: %.0f\n", variance.Amount)
	return updatedShift, variance, nil
}

// DocumentVarianceWithTransaction wraps DocumentVariance in a MongoDB transaction
func (s *CashierShiftService) DocumentVarianceWithTransaction(
	ctx context.Context,
	shiftID primitive.ObjectID,
	reason cashier.VarianceReason,
	notes string,
	userID, deviceID string,
) (*cashier.CashierShift, error) {
	fmt.Printf("🔵 [DOCUMENT VARIANCE TX] Starting for shift %s\n", shiftID.Hex())

	// Start MongoDB transaction
	session, err := s.mongoClient.StartSession()
	if err != nil {
		return nil, fmt.Errorf("failed to start transaction session: %w", err)
	}
	defer session.EndSession(ctx)

	var updatedShift *cashier.CashierShift

	// Execute transaction
	_, err = session.WithTransaction(ctx, func(sessCtx mongo.SessionContext) (interface{}, error) {
		// 1. Get the shift
		shift, err := s.cashierShiftRepo.FindByID(sessCtx, shiftID)
		if err != nil {
			return nil, err
		}

		// 2. Validate step
		if err := s.stateMachineManager.ValidateCashierShiftStep(shift, "document_variance"); err != nil {
			return nil, err
		}

		// 3. Document variance
		if err := shift.DocumentVariance(reason, notes, userID, deviceID, time.Now()); err != nil {
			return nil, err
		}

		// 4. Save the shift
		if err := s.cashierShiftRepo.Save(sessCtx, shift); err != nil {
			return nil, err
		}

		updatedShift = shift
		return nil, nil
	})

	if err != nil {
		fmt.Printf("❌ [DOCUMENT VARIANCE TX] Failed: %v\n", err)
		return nil, err
	}

	fmt.Printf("✅ [DOCUMENT VARIANCE TX] Success\n")
	return updatedShift, nil
}

// CloseShiftWithTransaction wraps CloseShift in a MongoDB transaction
func (s *CashierShiftService) CloseShiftWithTransaction(
	ctx context.Context,
	shiftID primitive.ObjectID,
	userID, deviceID string,
) (*cashier.CashierShift, error) {
	fmt.Printf("🔵 [CLOSE SHIFT TX] Starting for shift %s\n", shiftID.Hex())

	// Start MongoDB transaction
	session, err := s.mongoClient.StartSession()
	if err != nil {
		return nil, fmt.Errorf("failed to start transaction session: %w", err)
	}
	defer session.EndSession(ctx)

	var updatedShift *cashier.CashierShift

	// Execute transaction
	_, err = session.WithTransaction(ctx, func(sessCtx mongo.SessionContext) (interface{}, error) {
		// 1. Get the shift
		shift, err := s.cashierShiftRepo.FindByID(sessCtx, shiftID)
		if err != nil {
			return nil, err
		}

		// 2. Validate state transition
		if err := s.stateMachineManager.ValidateCashierShiftTransition(shift, cashier.EventCloseShift); err != nil {
			return nil, err
		}

		// 3. Check if all waiter shifts are closed
		allClosed, openShifts, err := s.CheckWaiterShiftsClosed(sessCtx)
		if err != nil {
			return nil, err
		}
		if !allClosed {
			return nil, fmt.Errorf("cannot close cashier shift: %d waiter shift(s) still open", len(openShifts))
		}

		// 4. Close the shift
		if err := shift.Close(userID, deviceID, time.Now()); err != nil {
			return nil, err
		}

		// 5. Save the shift
		if err := s.cashierShiftRepo.Save(sessCtx, shift); err != nil {
			return nil, err
		}

		updatedShift = shift
		return nil, nil
	})

	if err != nil {
		fmt.Printf("❌ [CLOSE SHIFT TX] Failed: %v\n", err)
		return nil, err
	}

	fmt.Printf("✅ [CLOSE SHIFT TX] Success\n")
	return updatedShift, nil
}

// CancelClosureWithTransaction wraps CancelClosure in a MongoDB transaction
func (s *CashierShiftService) CancelClosureWithTransaction(
	ctx context.Context,
	shiftID primitive.ObjectID,
	userID, deviceID string,
) (*cashier.CashierShift, error) {
	fmt.Printf("🔵 [CANCEL CLOSURE TX] Starting for shift %s\n", shiftID.Hex())

	// Start MongoDB transaction
	session, err := s.mongoClient.StartSession()
	if err != nil {
		return nil, fmt.Errorf("failed to start transaction session: %w", err)
	}
	defer session.EndSession(ctx)

	var updatedShift *cashier.CashierShift

	// Execute transaction
	_, err = session.WithTransaction(ctx, func(sessCtx mongo.SessionContext) (interface{}, error) {
		// 1. Get the shift
		shift, err := s.cashierShiftRepo.FindByID(sessCtx, shiftID)
		if err != nil {
			return nil, err
		}

		// 2. Validate state transition
		if err := s.stateMachineManager.ValidateCashierShiftTransition(shift, cashier.EventCancelClosure); err != nil {
			return nil, err
		}

		// 3. Check if can cancel
		if !s.stateMachineManager.CanCancelCashierShiftClosure(shift) {
			return nil, errors.New("cannot cancel closure: critical steps have been completed")
		}

		// 4. Cancel closure
		if err := shift.CancelClosure(userID, deviceID, time.Now()); err != nil {
			return nil, err
		}

		// 5. Save the shift
		if err := s.cashierShiftRepo.Save(sessCtx, shift); err != nil {
			return nil, err
		}

		updatedShift = shift
		return nil, nil
	})

	if err != nil {
		fmt.Printf("❌ [CANCEL CLOSURE TX] Failed: %v\n", err)
		return nil, err
	}

	fmt.Printf("✅ [CANCEL CLOSURE TX] Success\n")
	return updatedShift, nil
}


// ============================================================================
// COMPLETE CLOSURE - FRONTEND-DRIVEN APPROACH
// ============================================================================

// CompleteClosure completes the entire closure workflow in one transaction.
// This is a frontend-driven approach where all data is collected on the client
// and submitted once. If user clicks "Back", no API call is made and no data is saved.
func (s *CashierShiftService) CompleteClosure(
	ctx context.Context,
	shiftID primitive.ObjectID,
	actualCash float64,
	varianceReason *cashier.VarianceReason,
	varianceNotes *string,
	userID, deviceID string,
) (*cashier.CashierShift, error) {
	fmt.Printf("🔵 [COMPLETE CLOSURE] Starting for shift %s\n", shiftID.Hex())

	// Start MongoDB transaction
	session, err := s.mongoClient.StartSession()
	if err != nil {
		return nil, fmt.Errorf("failed to start transaction session: %w", err)
	}
	defer session.EndSession(ctx)

	var updatedShift *cashier.CashierShift

	// Execute entire workflow in ONE transaction
	_, err = session.WithTransaction(ctx, func(sessCtx mongo.SessionContext) (interface{}, error) {
		now := time.Now()

		// 1. Get the shift
		shift, err := s.cashierShiftRepo.FindByID(sessCtx, shiftID)
		if err != nil {
			return nil, err
		}

		// 2. Validate shift is in correct state
		if shift.Status != cashier.CashierShiftOpen {
			return nil, fmt.Errorf("shift must be in OPEN status, current: %s", shift.Status)
		}

		// 3. Check if all waiter shifts are closed
		allClosed, openShifts, err := s.CheckWaiterShiftsClosed(sessCtx)
		if err != nil {
			return nil, err
		}
		if !allClosed {
			return nil, fmt.Errorf("cannot close cashier shift: %d waiter shift(s) still open", len(openShifts))
		}

		// 4. Initiate closure
		if err := shift.InitiateClosure(userID, deviceID, now); err != nil {
			return nil, fmt.Errorf("failed to initiate closure: %w", err)
		}

		// 5. Record actual cash
		variance, err := shift.RecordActualCash(actualCash, userID, deviceID, now)
		if err != nil {
			return nil, fmt.Errorf("failed to record actual cash: %w", err)
		}

		// 6. Document variance if needed
		if variance.RequiresDocumentation() {
			if varianceReason == nil || varianceNotes == nil {
				return nil, errors.New("variance requires documentation: reason and notes are required")
			}
			if err := shift.DocumentVariance(*varianceReason, *varianceNotes, userID, deviceID, now); err != nil {
				return nil, fmt.Errorf("failed to document variance: %w", err)
			}
		}

		// 7. Close the shift
		if err := shift.Close(userID, deviceID, now); err != nil {
			return nil, fmt.Errorf("failed to close shift: %w", err)
		}

		// 8. Save the shift
		if err := s.cashierShiftRepo.Save(sessCtx, shift); err != nil {
			return nil, fmt.Errorf("failed to save shift: %w", err)
		}

		updatedShift = shift
		return nil, nil
	})

	if err != nil {
		fmt.Printf("❌ [COMPLETE CLOSURE] Failed: %v\n", err)
		return nil, err
	}

	fmt.Printf("✅ [COMPLETE CLOSURE] Success\n")
	return updatedShift, nil
}

// ============================================================================
// FUND HANDOVER - NEW FUNCTIONALITY
// ============================================================================

// ManagedFundsSummary represents the summary of funds managed by a cashier
type ManagedFundsSummary struct {
	CashierShiftID    primitive.ObjectID `json:"cashier_shift_id"`
	StartingFloat     float64            `json:"starting_float"`
	ReceivedCash      float64            `json:"received_cash"`
	ReceivedTransfer  float64            `json:"received_transfer"`
	TotalManagedFunds float64            `json:"total_managed_funds"`
	ExpectedCash      float64            `json:"expected_cash"`
	HandoverCount     int                `json:"handover_count"`
}

// GetManagedFunds returns the current managed funds for a cashier shift.
// This shows the cashier how much money they are responsible for.
//
// Parameters:
//   - ctx: Context for the operation
//   - shiftID: The ID of the cashier shift
//
// Returns:
//   - *ManagedFundsSummary: Summary of managed funds
//   - error: If query fails
func (s *CashierShiftService) GetManagedFunds(
	ctx context.Context,
	shiftID primitive.ObjectID,
) (*ManagedFundsSummary, error) {
	shift, err := s.cashierShiftRepo.FindByID(ctx, shiftID)
	if err != nil {
		return nil, err
	}

	return &ManagedFundsSummary{
		CashierShiftID:    shift.ID,
		StartingFloat:     shift.StartingFloat,
		ReceivedCash:      shift.ReceivedCash,
		ReceivedTransfer:  shift.ReceivedTransfer,
		TotalManagedFunds: shift.ReceivedCash + shift.ReceivedTransfer,
		ExpectedCash:      shift.StartingFloat + shift.ReceivedCash,
		HandoverCount:     shift.HandoverCount,
	}, nil
}

// CloseShiftWithFundHandover closes a cashier shift and creates a fund handover record.
// This is the complete closure flow that includes handing over funds back to the fund.
//
// The operation is atomic - either everything succeeds or everything rolls back.
//
// Parameters:
//   - ctx: Context for the operation
//   - shiftID: The ID of the cashier shift to close
//   - actualCash: The actual cash amount counted by the cashier
//   - varianceReason: The reason for variance (if any)
//   - varianceNotes: Detailed notes about the variance (if any)
//   - receiverID: Optional ID of the person receiving the funds (for future use)
//   - userID: The ID of the user performing the operation
//   - deviceID: The ID of the device used
//
// Returns:
//   - *cashier.CashierShift: The closed shift
//   - *cashier.FundHandover: The fund handover record
//   - error: If validation or operation fails
func (s *CashierShiftService) CloseShiftWithFundHandover(
	ctx context.Context,
	shiftID primitive.ObjectID,
	actualCash float64,
	varianceReason *cashier.VarianceReason,
	varianceNotes *string,
	receiverID *primitive.ObjectID,
	userID, deviceID string,
) (*cashier.CashierShift, *cashier.FundHandover, error) {
	fmt.Printf("🔵 [CLOSE WITH FUND HANDOVER] Starting for shift %s\n", shiftID.Hex())

	// Start MongoDB transaction
	session, err := s.mongoClient.StartSession()
	if err != nil {
		return nil, nil, fmt.Errorf("failed to start transaction session: %w", err)
	}
	defer session.EndSession(ctx)

	var updatedShift *cashier.CashierShift
	var fundHandover *cashier.FundHandover

	// Execute entire workflow in ONE transaction
	_, err = session.WithTransaction(ctx, func(sessCtx mongo.SessionContext) (interface{}, error) {
		now := time.Now()

		// 1. Get the shift
		shift, err := s.cashierShiftRepo.FindByID(sessCtx, shiftID)
		if err != nil {
			return nil, err
		}

		// 2. Validate shift is in correct state
		if shift.Status != cashier.CashierShiftOpen {
			return nil, fmt.Errorf("shift must be in OPEN status, current: %s", shift.Status)
		}

		// 3. Check if all waiter shifts are closed
		allClosed, openShifts, err := s.CheckWaiterShiftsClosed(sessCtx)
		if err != nil {
			return nil, err
		}
		if !allClosed {
			return nil, fmt.Errorf("cannot close cashier shift: %d waiter shift(s) still open", len(openShifts))
		}

		// 4. Calculate expected cash
		expectedCash := shift.StartingFloat + shift.ReceivedCash

		// 5. Create fund handover record
		handover := cashier.NewFundHandover(
			shift.ID,
			shift.CashierID,
			shift.CashierName,
			actualCash,
			shift.ReceivedTransfer,
			expectedCash,
		)

		// Debug logging
		log.Printf("🔍 [VARIANCE DEBUG] actualCash=%.2f, expectedCash=%.2f, variance=%.10f, hasVariance=%v",
			actualCash, expectedCash, handover.VarianceAmount, handover.HasVariance())

		// 6. Document variance if exists
		if handover.HasVariance() {
			log.Printf("⚠️ [VARIANCE] Variance detected, requiring documentation")
			if varianceReason == nil || varianceNotes == nil {
				return nil, errors.New("variance requires documentation: reason and notes are required")
			}
			if err := handover.DocumentVariance(*varianceReason, *varianceNotes); err != nil {
				return nil, fmt.Errorf("failed to document variance: %w", err)
			}
		} else {
			log.Printf("✅ [NO VARIANCE] No variance detected, skipping documentation")
		}

		// 7. Set receiver if provided (for future use)
		if receiverID != nil {
			// TODO: Get receiver name from user service
			handover.SetReceiver(*receiverID, "Manager Name")
		}

		// 8. Validate fund handover
		if err := handover.Validate(); err != nil {
			return nil, fmt.Errorf("fund handover validation failed: %w", err)
		}

		// 9. Save fund handover record
		if err := s.fundHandoverRepo.Create(sessCtx, handover); err != nil {
			return nil, fmt.Errorf("failed to create fund handover: %w", err)
		}

		// 10. Initiate closure
		if err := shift.InitiateClosure(userID, deviceID, now); err != nil {
			return nil, fmt.Errorf("failed to initiate closure: %w", err)
		}

		// 11. Update system cash to match expected cash before recording actual cash
		// This ensures variance calculation is based on the correct expected amount
		shift.UpdateSystemCash(expectedCash)
		log.Printf("📊 [SYSTEM CASH UPDATE] Updated SystemCash to %.2f (expectedCash)", expectedCash)

		// 12. Record actual cash in shift (for variance tracking)
		variance, err := shift.RecordActualCash(actualCash, userID, deviceID, now)
		if err != nil {
			return nil, fmt.Errorf("failed to record actual cash: %w", err)
		}
		log.Printf("📊 [VARIANCE IN SHIFT] amount=%.10f, requiresDoc=%v", variance.Amount, variance.RequiresDocumentation())

		// 13. Document variance in shift if needed (already handled in FundHandover)
		// Variance documentation is handled in the FundHandover record (step 6)
		// No need to duplicate the check here
		if handover.HasVariance() && handover.VarianceReason != nil {
			// Sync variance documentation from FundHandover to CashierShift
			if err := shift.DocumentVariance(*handover.VarianceReason, handover.VarianceNotes, userID, deviceID, now); err != nil {
				return nil, fmt.Errorf("failed to sync variance documentation to shift: %w", err)
			}
		}

		// 13. Close the shift
		if err := shift.Close(userID, deviceID, now); err != nil {
			return nil, fmt.Errorf("failed to close shift: %w", err)
		}

		// 14. Save the shift
		if err := s.cashierShiftRepo.Save(sessCtx, shift); err != nil {
			return nil, fmt.Errorf("failed to save shift: %w", err)
		}

		updatedShift = shift
		fundHandover = handover
		return nil, nil
	})

	if err != nil {
		fmt.Printf("❌ [CLOSE WITH FUND HANDOVER] Failed: %v\n", err)
		return nil, nil, err
	}

	fmt.Printf("✅ [CLOSE WITH FUND HANDOVER] Success - Handover ID: %s\n", fundHandover.ID.Hex())
	return updatedShift, fundHandover, nil
}

// GetFundHandoverByShift retrieves the fund handover record for a cashier shift
func (s *CashierShiftService) GetFundHandoverByShift(
	ctx context.Context,
	shiftID primitive.ObjectID,
) (*cashier.FundHandover, error) {
	return s.fundHandoverRepo.FindByCashierShift(ctx, shiftID)
}

// GetFundHandoverHistory retrieves fund handover history for a cashier
func (s *CashierShiftService) GetFundHandoverHistory(
	ctx context.Context,
	cashierID primitive.ObjectID,
	page, pageSize int,
) ([]*cashier.FundHandover, int64, error) {
	return s.fundHandoverRepo.FindByCashier(ctx, cashierID, page, pageSize)
}

// GetFundHandoversByDateRange retrieves fund handovers within a date range
func (s *CashierShiftService) GetFundHandoversByDateRange(
	ctx context.Context,
	startDate, endDate time.Time,
	page, pageSize int,
) ([]*cashier.FundHandover, int64, error) {
	return s.fundHandoverRepo.FindByDateRange(ctx, startDate, endDate, page, pageSize)
}
