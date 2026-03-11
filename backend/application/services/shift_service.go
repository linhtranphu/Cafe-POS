package services

import (
	"context"
	"errors"
	"log"
	"time"
	"cafe-pos/backend/domain"
	"cafe-pos/backend/domain/order"
	"cafe-pos/backend/infrastructure/mongodb"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type ShiftRepository interface {
	Create(ctx context.Context, s *order.Shift) error
	FindByID(ctx context.Context, id primitive.ObjectID) (*order.Shift, error)
	Update(ctx context.Context, id primitive.ObjectID, s *order.Shift) error
	FindOpenShiftByWaiter(ctx context.Context, waiterID primitive.ObjectID) (*order.Shift, error)
	FindOpenShiftByUser(ctx context.Context, userID primitive.ObjectID, roleType order.RoleType) (*order.Shift, error)
	FindOpenShifts(ctx context.Context) ([]*order.Shift, error)
	FindByWaiterID(ctx context.Context, waiterID primitive.ObjectID) ([]*order.Shift, error)
	FindByUserID(ctx context.Context, userID primitive.ObjectID, roleType order.RoleType) ([]*order.Shift, error)
	FindByDateRange(ctx context.Context, startDate, endDate time.Time) ([]*order.Shift, error)
	FindByRoleType(ctx context.Context, roleType order.RoleType) ([]*order.Shift, error)
	FindAll(ctx context.Context) ([]*order.Shift, error)
}

type ShiftService struct {
	shiftRepo           ShiftRepository
	orderRepo           OrderRepository
	stateMachineManager *domain.StateMachineManager
	journalService      *JournalService
	cashierShiftRepo    *mongodb.CashierShiftRepository
}

func NewShiftService(
	shiftRepo ShiftRepository,
	orderRepo OrderRepository,
	stateMachineManager *domain.StateMachineManager,
	journalService *JournalService,
	cashierShiftRepo *mongodb.CashierShiftRepository,
) *ShiftService {
	return &ShiftService{
		shiftRepo:           shiftRepo,
		orderRepo:           orderRepo,
		stateMachineManager: stateMachineManager,
		journalService:      journalService,
		cashierShiftRepo:    cashierShiftRepo,
	}
}

func (s *ShiftService) StartShift(ctx context.Context, req *order.StartShiftRequest, userID, userName string, roleType order.RoleType) (*order.Shift, error) {
	// Reject cashier role - cashier shifts are handled separately
	if roleType == "cashier" {
		return nil, errors.New("cashier shifts must be created using the cashier shift service")
	}
	
	// Validate role type is waiter or barista
	if !roleType.IsValid() {
		return nil, errors.New("invalid role type: must be waiter or barista")
	}
	
	userOID, _ := primitive.ObjectIDFromHex(userID)
	
	// Check if user already has an open shift for this role
	existingShift, _ := s.shiftRepo.FindOpenShiftByUser(ctx, userOID, roleType)
	
	// Validate using state machine
	if err := s.stateMachineManager.ValidateWaiterShiftStart(existingShift); err != nil {
		return nil, err
	}

	shift := &order.Shift{
		Type:              req.Type,
		Status:            order.ShiftOpen,
		RoleType:          roleType,
		UserID:            userOID,
		UserName:          userName,
		StartCash:         req.StartCash,
		CurrentCash:       req.StartCash,       // Initialize with start cash
		RemainingCash:     req.StartCash,       // Initialize with start cash
		TransferRevenue:   0,                   // Initialize transfer revenue
		RemainingTransfer: 0,                   // Initialize remaining transfer
		HandedOverCash:    0,                   // Initialize handed over cash
		HandedOverTransfer: 0,                  // Initialize handed over transfer
		StartedAt:         time.Now(),
	}

	if err := s.shiftRepo.Create(ctx, shift); err != nil {
		return nil, err
	}

	// Waiter phải có cashier ca mở trước khi mở ca
	if roleType == order.RoleWaiter && s.cashierShiftRepo != nil {
		openCashierShifts, err := s.cashierShiftRepo.FindOpen(ctx)
		if err != nil {
			return nil, errors.New("không thể kiểm tra ca thu ngân")
		}
		if len(openCashierShifts) == 0 {
			return nil, errors.New("thu ngân phải mở ca trước khi phục vụ mở ca")
		}

		cashierShift := openCashierShifts[0]

		// Ghi nhận tiền đầu ca: double-entry journal cash_drawer → waiter_float
		if req.StartCash > 0 && s.journalService != nil {
			if _, err := s.journalService.RecordWaiterShiftStart(ctx, req.StartCash, shift.ID, userOID, userName); err != nil {
				log.Printf("⚠️ [WAITER START FLOAT] Failed to record journal entry: %v (non-fatal)", err)
			} else {
				cashierShift.AddDistributedCash(req.StartCash)
				if saveErr := s.cashierShiftRepo.Save(ctx, cashierShift); saveErr != nil {
					log.Printf("⚠️ [WAITER START FLOAT] Failed to update cashier distributed cash: %v (non-fatal)", saveErr)
				} else {
					log.Printf("✅ [WAITER START FLOAT] Journal entry recorded: cash_drawer → waiter_float %.0f for %s", req.StartCash, userName)
				}
			}
		}
	}

	return shift, nil
}

func (s *ShiftService) EndShift(ctx context.Context, shiftID primitive.ObjectID, req *order.EndShiftRequest) (*order.Shift, error) {
	shift, err := s.shiftRepo.FindByID(ctx, shiftID)
	if err != nil {
		return nil, err
	}

	// Validate state transition using state machine
	if err := s.stateMachineManager.ValidateWaiterShiftTransition(shift, order.EventEndShift); err != nil {
		return nil, err
	}

	orders, err := s.orderRepo.FindByShiftID(ctx, shiftID)
	if err != nil {
		return nil, err
	}

	totalRevenue := 0.0
	for _, o := range orders {
		if o.Status == order.StatusPaid || o.Status == order.StatusInProgress || o.Status == order.StatusServed {
			totalRevenue += o.Total
		}
	}

	now := time.Now()
	shift.Status = order.ShiftClosed
	shift.EndCash = req.EndCash
	shift.TotalRevenue = totalRevenue
	shift.TotalOrders = len(orders)
	shift.EndedAt = &now

	if err := s.shiftRepo.Update(ctx, shiftID, shift); err != nil {
		return nil, err
	}
	return shift, nil
}

func (s *ShiftService) GetCurrentShift(ctx context.Context, userID primitive.ObjectID, roleType order.RoleType) (*order.Shift, error) {
	return s.shiftRepo.FindOpenShiftByUser(ctx, userID, roleType)
}

func (s *ShiftService) GetOpenShifts(ctx context.Context) ([]*order.Shift, error) {
	return s.shiftRepo.FindOpenShifts(ctx)
}

func (s *ShiftService) GetShiftsByUser(ctx context.Context, userID primitive.ObjectID, roleType order.RoleType) ([]*order.Shift, error) {
	return s.shiftRepo.FindByUserID(ctx, userID, roleType)
}

func (s *ShiftService) GetShiftsByWaiter(ctx context.Context, waiterID primitive.ObjectID) ([]*order.Shift, error) {
	return s.shiftRepo.FindByWaiterID(ctx, waiterID)
}

func (s *ShiftService) GetShiftsByRole(ctx context.Context, roleType order.RoleType) ([]*order.Shift, error) {
	return s.shiftRepo.FindByRoleType(ctx, roleType)
}

func (s *ShiftService) GetAllShifts(ctx context.Context) ([]*order.Shift, error) {
	return s.shiftRepo.FindAll(ctx)
}

func (s *ShiftService) GetShift(ctx context.Context, id primitive.ObjectID) (*order.Shift, error) {
	return s.shiftRepo.FindByID(ctx, id)
}

func (s *ShiftService) GetPendingOrdersByShift(ctx context.Context, shiftID primitive.ObjectID) ([]*order.Order, error) {
	orders, err := s.orderRepo.FindByShiftID(ctx, shiftID)
	if err != nil {
		return nil, err
	}

	pendingStatuses := map[order.OrderStatus]bool{
		order.StatusCreated:    true,
		order.StatusPaid:       true,
		order.StatusQueued:     true,
		order.StatusInProgress: true,
		order.StatusReady:      true,
	}

	var pending []*order.Order
	for _, o := range orders {
		if pendingStatuses[o.Status] {
			pending = append(pending, o)
		}
	}
	return pending, nil
}

func (s *ShiftService) CloseShiftAndLockOrders(ctx context.Context, shiftID primitive.ObjectID, req *order.EndShiftRequest) (*order.Shift, error) {
	// Get shift first to validate
	shift, err := s.shiftRepo.FindByID(ctx, shiftID)
	if err != nil {
		return nil, err
	}

	// Validate state transition using state machine
	if err := s.stateMachineManager.ValidateWaiterShiftTransition(shift, order.EventEndShift); err != nil {
		return nil, err
	}

	// End the shift
	shift, err = s.EndShift(ctx, shiftID, req)
	if err != nil {
		return nil, err
	}

	now := time.Now()
	orders, _ := s.orderRepo.FindByShiftID(ctx, shiftID)
	for _, o := range orders {
		// First pass: cancel any pending orders
		switch o.Status {
		case order.StatusCreated, order.StatusPaid, order.StatusQueued, order.StatusInProgress, order.StatusReady:
			o.Status = order.StatusCancelled
			o.CancelReason = "Ca làm việc đã kết thúc"
			o.UpdatedAt = now
			s.orderRepo.Update(ctx, o.ID, o)
		}
		// Second pass: lock all completed/cancelled orders
		if o.Status == order.StatusServed || o.Status == order.StatusCancelled {
			o.Status = order.StatusLocked
			o.LockedAt = &now
			o.UpdatedAt = now
			s.orderRepo.Update(ctx, o.ID, o)
		}
	}

	return shift, nil
}

// CalculateTransferRevenue calculates transfer revenue from orders and updates shift
func (s *ShiftService) CalculateTransferRevenue(ctx context.Context, shiftID primitive.ObjectID) error {
	shift, err := s.shiftRepo.FindByID(ctx, shiftID)
	if err != nil {
		return err
	}

	orders, err := s.orderRepo.FindByShiftID(ctx, shiftID)
	if err != nil {
		return err
	}

	// Calculate cash and transfer revenue separately
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

	// Update shift with transfer revenue
	shift.TransferRevenue = transferRevenue
	shift.RemainingTransfer = transferRevenue - shift.HandedOverTransfer
	shift.CurrentCash = shift.StartCash + cashRevenue
	shift.RemainingCash = shift.CurrentCash - shift.HandedOverCash
	shift.TotalRevenue = cashRevenue + transferRevenue
	shift.UpdatedAt = time.Now()

	return s.shiftRepo.Update(ctx, shiftID, shift)
}
