package services

import (
	"context"
	"errors"
	"time"
	"cafe-pos/backend/domain/order"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type ShiftRepository interface {
	Create(ctx context.Context, s *order.Shift) error
	FindByID(ctx context.Context, id primitive.ObjectID) (*order.Shift, error)
	Update(ctx context.Context, id primitive.ObjectID, s *order.Shift) error
	FindOpenShiftByWaiter(ctx context.Context, waiterID primitive.ObjectID) (*order.Shift, error)
	FindOpenShifts(ctx context.Context) ([]*order.Shift, error)
	FindByWaiterID(ctx context.Context, waiterID primitive.ObjectID) ([]*order.Shift, error)
	FindByDateRange(ctx context.Context, startDate, endDate time.Time) ([]*order.Shift, error)
	FindAll(ctx context.Context) ([]*order.Shift, error)
}

type ShiftService struct {
	shiftRepo ShiftRepository
	orderRepo OrderRepository
}

func NewShiftService(shiftRepo ShiftRepository, orderRepo OrderRepository) *ShiftService {
	return &ShiftService{
		shiftRepo: shiftRepo,
		orderRepo: orderRepo,
	}
}

func (s *ShiftService) StartShift(ctx context.Context, req *order.StartShiftRequest, waiterID, waiterName string) (*order.Shift, error) {
	waiterOID, _ := primitive.ObjectIDFromHex(waiterID)
	
	existingShift, _ := s.shiftRepo.FindOpenShiftByWaiter(ctx, waiterOID)
	if existingShift != nil {
		return nil, errors.New("waiter already has an open shift")
	}

	shift := &order.Shift{
		Type:       req.Type,
		Status:     order.ShiftOpen,
		WaiterID:   waiterOID,
		WaiterName: waiterName,
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

	if shift.Status != order.ShiftOpen {
		return nil, errors.New("shift is not open")
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

func (s *ShiftService) GetCurrentShift(ctx context.Context, waiterID primitive.ObjectID) (*order.Shift, error) {
	return s.shiftRepo.FindOpenShiftByWaiter(ctx, waiterID)
}

func (s *ShiftService) GetOpenShifts(ctx context.Context) ([]*order.Shift, error) {
	return s.shiftRepo.FindOpenShifts(ctx)
}

func (s *ShiftService) GetShiftsByWaiter(ctx context.Context, waiterID primitive.ObjectID) ([]*order.Shift, error) {
	return s.shiftRepo.FindByWaiterID(ctx, waiterID)
}

func (s *ShiftService) GetAllShifts(ctx context.Context) ([]*order.Shift, error) {
	return s.shiftRepo.FindAll(ctx)
}

func (s *ShiftService) GetShift(ctx context.Context, id primitive.ObjectID) (*order.Shift, error) {
	return s.shiftRepo.FindByID(ctx, id)
}

func (s *ShiftService) CloseShiftAndLockOrders(ctx context.Context, shiftID primitive.ObjectID, req *order.EndShiftRequest) (*order.Shift, error) {
	shift, err := s.EndShift(ctx, shiftID, req)
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
