package services

import (
	"context"
	"errors"
	"time"
	"cafe-pos/backend/domain"
	"cafe-pos/backend/domain/order"
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
}

func NewShiftService(
	shiftRepo ShiftRepository,
	orderRepo OrderRepository,
	stateMachineManager *domain.StateMachineManager,
) *ShiftService {
	return &ShiftService{
		shiftRepo:           shiftRepo,
		orderRepo:           orderRepo,
		stateMachineManager: stateMachineManager,
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
		Type:       req.Type,
		Status:     order.ShiftOpen,
		RoleType:   roleType,
		UserID:     userOID,
		UserName:   userName,
		StartCash:  req.StartCash,
		StartedAt:  time.Now(),
	}

	if err := s.shiftRepo.Create(ctx, shift); err != nil {
		return nil, err
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

	orders, _ := s.orderRepo.FindByShiftID(ctx, shiftID)
	for _, o := range orders {
		// Lock orders that are completed (served or cancelled)
		if o.Status == order.StatusServed || o.Status == order.StatusCancelled {
			now := time.Now()
			o.Status = order.StatusLocked
			o.LockedAt = &now
			s.orderRepo.Update(ctx, o.ID, o)
		}
	}

	return shift, nil
}
