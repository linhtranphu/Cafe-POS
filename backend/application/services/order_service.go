package services

import (
	"context"
	"errors"
	"time"
	"cafe-pos/backend/domain/order"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type OrderRepository interface {
	Create(ctx context.Context, o *order.Order) error
	FindByID(ctx context.Context, id primitive.ObjectID) (*order.Order, error)
	Update(ctx context.Context, id primitive.ObjectID, o *order.Order) error
	FindByShiftID(ctx context.Context, shiftID primitive.ObjectID) ([]*order.Order, error)
	FindByWaiterID(ctx context.Context, waiterID primitive.ObjectID) ([]*order.Order, error)
	FindByStatus(ctx context.Context, status order.OrderStatus) ([]*order.Order, error)
	FindByTableID(ctx context.Context, tableID primitive.ObjectID) ([]*order.Order, error)
	FindAll(ctx context.Context) ([]*order.Order, error)
}

type OrderService struct {
	orderRepo OrderRepository
	shiftRepo ShiftRepository
	tableRepo TableRepository
}

func NewOrderService(orderRepo OrderRepository, shiftRepo ShiftRepository, tableRepo TableRepository) *OrderService {
	return &OrderService{
		orderRepo: orderRepo,
		shiftRepo: shiftRepo,
		tableRepo: tableRepo,
	}
}

func (s *OrderService) CreateOrder(ctx context.Context, req *order.CreateOrderRequest, waiterID, waiterName string) (*order.Order, error) {
	shiftID, _ := primitive.ObjectIDFromHex(req.ShiftID)
	shift, err := s.shiftRepo.FindByID(ctx, shiftID)
	if err != nil || shift.Status != order.ShiftOpen {
		return nil, errors.New("no open shift found")
	}

	tableID, _ := primitive.ObjectIDFromHex(req.TableID)
	table, err := s.tableRepo.FindByID(ctx, tableID)
	if err != nil {
		return nil, errors.New("table not found")
	}

	waiterOID, _ := primitive.ObjectIDFromHex(waiterID)
	o := &order.Order{
		TableID:    tableID,
		TableName:  table.Name,
		WaiterID:   waiterOID,
		WaiterName: waiterName,
		ShiftID:    shiftID,
		Items:      req.Items,
		Status:     order.StatusCreated,
		Note:       req.Note,
	}

	o.CalculateTotal()
	
	if err := s.orderRepo.Create(ctx, o); err != nil {
		return nil, err
	}

	s.tableRepo.UpdateStatus(ctx, tableID, order.TableOccupied)
	return o, nil
}

func (s *OrderService) ConfirmOrder(ctx context.Context, id primitive.ObjectID, discount float64) (*order.Order, error) {
	o, err := s.orderRepo.FindByID(ctx, id)
	if err != nil {
		return nil, err
	}

	if !o.CanTransitionTo(order.StatusUnpaid) {
		return nil, errors.New("cannot confirm order in current state")
	}

	o.Discount = discount
	o.CalculateTotal()
	o.Status = order.StatusUnpaid

	if err := s.orderRepo.Update(ctx, id, o); err != nil {
		return nil, err
	}
	return o, nil
}

func (s *OrderService) PayOrder(ctx context.Context, id primitive.ObjectID, req *order.PaymentRequest) (*order.Order, error) {
	o, err := s.orderRepo.FindByID(ctx, id)
	if err != nil {
		return nil, err
	}

	if !o.CanTransitionTo(order.StatusPaid) {
		return nil, errors.New("cannot pay order in current state")
	}

	collectorID, _ := primitive.ObjectIDFromHex(req.CollectorID)
	now := time.Now()
	
	o.Status = order.StatusPaid
	o.PaymentMethod = req.PaymentMethod
	o.CollectorID = collectorID
	o.CollectorName = req.CollectorName
	o.PaidAt = &now

	if err := s.orderRepo.Update(ctx, id, o); err != nil {
		return nil, err
	}
	return o, nil
}

func (s *OrderService) SendToKitchen(ctx context.Context, id primitive.ObjectID) (*order.Order, error) {
	o, err := s.orderRepo.FindByID(ctx, id)
	if err != nil {
		return nil, err
	}

	if o.Status != order.StatusPaid {
		return nil, errors.New("order must be paid before sending to kitchen")
	}

	if !o.CanTransitionTo(order.StatusInProgress) {
		return nil, errors.New("cannot send order to kitchen in current state")
	}

	now := time.Now()
	o.Status = order.StatusInProgress
	o.SentToKitchenAt = &now

	if err := s.orderRepo.Update(ctx, id, o); err != nil {
		return nil, err
	}
	return o, nil
}

func (s *OrderService) ServeOrder(ctx context.Context, id primitive.ObjectID) (*order.Order, error) {
	o, err := s.orderRepo.FindByID(ctx, id)
	if err != nil {
		return nil, err
	}

	if !o.CanTransitionTo(order.StatusServed) {
		return nil, errors.New("cannot serve order in current state")
	}

	now := time.Now()
	o.Status = order.StatusServed
	o.ServedAt = &now

	if err := s.orderRepo.Update(ctx, id, o); err != nil {
		return nil, err
	}
	return o, nil
}

func (s *OrderService) CancelOrder(ctx context.Context, id primitive.ObjectID, req *order.CancelOrderRequest) (*order.Order, error) {
	o, err := s.orderRepo.FindByID(ctx, id)
	if err != nil {
		return nil, err
	}

	if !o.CanTransitionTo(order.StatusCancelled) {
		return nil, errors.New("cannot cancel order in current state")
	}

	o.Status = order.StatusCancelled
	o.CancelReason = req.Reason

	if err := s.orderRepo.Update(ctx, id, o); err != nil {
		return nil, err
	}

	s.tableRepo.UpdateStatus(ctx, o.TableID, order.TableEmpty)
	return o, nil
}

func (s *OrderService) RefundOrder(ctx context.Context, id primitive.ObjectID, req *order.RefundOrderRequest) (*order.Order, error) {
	o, err := s.orderRepo.FindByID(ctx, id)
	if err != nil {
		return nil, err
	}

	if !o.CanTransitionTo(order.StatusRefunded) {
		return nil, errors.New("cannot refund order in current state")
	}

	o.Status = order.StatusRefunded
	o.RefundReason = req.Reason

	if err := s.orderRepo.Update(ctx, id, o); err != nil {
		return nil, err
	}

	s.tableRepo.UpdateStatus(ctx, o.TableID, order.TableEmpty)
	return o, nil
}

func (s *OrderService) LockOrder(ctx context.Context, id primitive.ObjectID) (*order.Order, error) {
	o, err := s.orderRepo.FindByID(ctx, id)
	if err != nil {
		return nil, err
	}

	if !o.CanTransitionTo(order.StatusLocked) {
		return nil, errors.New("cannot lock order in current state")
	}

	now := time.Now()
	o.Status = order.StatusLocked
	o.LockedAt = &now

	if err := s.orderRepo.Update(ctx, id, o); err != nil {
		return nil, err
	}
	return o, nil
}

func (s *OrderService) GetOrdersByWaiter(ctx context.Context, waiterID primitive.ObjectID) ([]*order.Order, error) {
	return s.orderRepo.FindByWaiterID(ctx, waiterID)
}

func (s *OrderService) GetOrdersByShift(ctx context.Context, shiftID primitive.ObjectID) ([]*order.Order, error) {
	return s.orderRepo.FindByShiftID(ctx, shiftID)
}

func (s *OrderService) GetAllOrders(ctx context.Context) ([]*order.Order, error) {
	return s.orderRepo.FindAll(ctx)
}

func (s *OrderService) GetOrder(ctx context.Context, id primitive.ObjectID) (*order.Order, error) {
	return s.orderRepo.FindByID(ctx, id)
}
