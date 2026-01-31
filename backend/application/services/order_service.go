package services

import (
	"context"
	"errors"
	"fmt"
	"time"
	"cafe-pos/backend/domain"
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
	FindByOrderNumber(ctx context.Context, orderNumber string) (*order.Order, error)
	FindAll(ctx context.Context) ([]*order.Order, error)
}

type OrderService struct {
	orderRepo           OrderRepository
	shiftRepo           ShiftRepository
	stateMachineManager *domain.StateMachineManager
}

func NewOrderService(
	orderRepo OrderRepository,
	shiftRepo ShiftRepository,
	stateMachineManager *domain.StateMachineManager,
) *OrderService {
	return &OrderService{
		orderRepo:           orderRepo,
		shiftRepo:           shiftRepo,
		stateMachineManager: stateMachineManager,
	}
}

func (s *OrderService) CreateOrder(ctx context.Context, req *order.CreateOrderRequest, waiterID, waiterName string) (*order.Order, error) {
	shiftID, _ := primitive.ObjectIDFromHex(req.ShiftID)
	shift, err := s.shiftRepo.FindByID(ctx, shiftID)
	if err != nil || shift.Status != order.ShiftOpen {
		return nil, errors.New("no open shift found")
	}

	// Generate order number (format: YYYYMMDD-HHMMSS-XXX)
	now := time.Now()
	orderNumber := fmt.Sprintf("%s-%03d", now.Format("20060102-150405"), now.Nanosecond()/1000000%1000)

	waiterOID, _ := primitive.ObjectIDFromHex(waiterID)
	o := &order.Order{
		OrderNumber:  orderNumber,
		CustomerName: req.CustomerName,
		WaiterID:     waiterOID,
		WaiterName:   waiterName,
		ShiftID:      shiftID,
		Items:        req.Items,
		Status:       order.StatusCreated,
		Note:         req.Note,
		AmountPaid:   0,
	}

	o.CalculateTotal()
	
	if err := s.orderRepo.Create(ctx, o); err != nil {
		return nil, err
	}

	return o, nil
}

func (s *OrderService) CollectPayment(ctx context.Context, id primitive.ObjectID, req *order.PaymentRequest) (*order.Order, error) {
	o, err := s.orderRepo.FindByID(ctx, id)
	if err != nil {
		return nil, err
	}

	// Validate state transition using state machine
	if err := s.stateMachineManager.ValidateOrderTransition(o, order.EventPayOrder); err != nil {
		return nil, fmt.Errorf("payment validation failed: %w", err)
	}

	collectorID, _ := primitive.ObjectIDFromHex(req.CollectorID)
	now := time.Now()
	
	// Add to amount paid
	o.AmountPaid += req.Amount
	o.PaymentMethod = req.PaymentMethod
	o.CollectorID = collectorID
	o.CollectorName = req.CollectorName
	
	// Recalculate amounts
	o.CalculateTotal()
	
	// If fully paid, mark as PAID
	if o.IsFullyPaid() {
		o.Status = order.StatusPaid
		o.PaidAt = &now
	}

	if err := s.orderRepo.Update(ctx, id, o); err != nil {
		return nil, err
	}
	return o, nil
}

func (s *OrderService) EditOrder(ctx context.Context, id primitive.ObjectID, req *order.EditOrderRequest) (*order.EditOrderResponse, error) {
	o, err := s.orderRepo.FindByID(ctx, id)
	if err != nil {
		return nil, err
	}

	// Validate using state machine
	if !s.stateMachineManager.CanModifyOrder(o) {
		return nil, fmt.Errorf("cannot modify order in state %s", o.Status)
	}

	// Store old total for refund calculation
	oldTotal := o.Total
	oldAmountPaid := o.AmountPaid

	// Update order details
	o.Items = req.Items
	o.Discount = req.Discount
	o.Note = req.Note
	
	// Recalculate totals
	o.CalculateTotal()

	response := &order.EditOrderResponse{
		Order: o,
	}

	// Handle refund if new total is less than amount paid
	if o.Total < oldAmountPaid {
		excessAmount := oldAmountPaid - o.Total
		o.RefundAmount += excessAmount
		o.AmountPaid = o.Total // Adjust amount paid to match new total
		o.RefundReason = fmt.Sprintf("Auto refund due to order edit. Old total: %.2f, New total: %.2f", oldTotal, o.Total)
		
		// Recalculate after refund adjustment
		o.CalculateTotal()
		
		// Add refund info to response
		response.RefundAmount = excessAmount
		response.RefundReason = o.RefundReason
		response.Message = fmt.Sprintf("Order updated. Refund amount: %.2f VND", excessAmount)
	} else if o.Total > oldAmountPaid {
		// Need additional payment
		response.Message = fmt.Sprintf("Order updated. Additional payment needed: %.2f VND", o.AmountDue)
	} else {
		response.Message = "Order updated successfully"
	}

	if err := s.orderRepo.Update(ctx, id, o); err != nil {
		return nil, err
	}
	
	return response, nil
}

func (s *OrderService) RefundPartial(ctx context.Context, id primitive.ObjectID, req *order.RefundRequest) (*order.Order, error) {
	o, err := s.orderRepo.FindByID(ctx, id)
	if err != nil {
		return nil, err
	}

	// Validate state transition using state machine
	if err := s.stateMachineManager.ValidateOrderTransition(o, order.EventRefundOrder); err != nil {
		return nil, fmt.Errorf("refund validation failed: %w", err)
	}

	if req.Amount > o.AmountPaid {
		return nil, errors.New("refund amount cannot exceed amount paid")
	}

	// Reduce amount paid
	o.AmountPaid -= req.Amount
	o.RefundAmount += req.Amount
	o.RefundReason = req.Reason
	
	// Recalculate amounts
	o.CalculateTotal()

	if err := s.orderRepo.Update(ctx, id, o); err != nil {
		return nil, err
	}
	return o, nil
}

// SendToBar - Waiter sends order to barista queue
func (s *OrderService) SendToBar(ctx context.Context, id primitive.ObjectID) (*order.Order, error) {
	o, err := s.orderRepo.FindByID(ctx, id)
	if err != nil {
		return nil, err
	}

	// Validate state transition using state machine
	if err := s.stateMachineManager.ValidateOrderTransition(o, order.EventSendToBar); err != nil {
		return nil, fmt.Errorf("send to bar validation failed: %w", err)
	}

	now := time.Now()
	o.Status = order.StatusQueued
	o.QueuedAt = &now

	if err := s.orderRepo.Update(ctx, id, o); err != nil {
		return nil, err
	}
	return o, nil
}

// AcceptOrder - BR-06: Only Barista can move order to IN_PROGRESS
// BR-13: Barista must have an open shift to accept orders
func (s *OrderService) AcceptOrder(ctx context.Context, id primitive.ObjectID, baristaID, baristaName string) (*order.Order, error) {
	o, err := s.orderRepo.FindByID(ctx, id)
	if err != nil {
		return nil, err
	}

	// Validate state transition using state machine
	if err := s.stateMachineManager.ValidateOrderTransition(o, order.EventStartPreparing); err != nil {
		return nil, fmt.Errorf("accept order validation failed: %w", err)
	}

	// BR-13: Check if barista has an open shift
	baristaOID, _ := primitive.ObjectIDFromHex(baristaID)
	shift, err := s.shiftRepo.FindOpenShiftByUser(ctx, baristaOID, order.RoleBarista)
	if err != nil || shift == nil {
		return nil, errors.New("barista must open a shift before accepting orders")
	}

	now := time.Now()
	o.Status = order.StatusInProgress
	o.BaristaID = baristaOID
	o.BaristaName = baristaName
	o.AcceptedAt = &now

	if err := s.orderRepo.Update(ctx, id, o); err != nil {
		return nil, err
	}
	return o, nil
}

// FinishPreparing - BR-09: Barista marks drink as READY
func (s *OrderService) FinishPreparing(ctx context.Context, id primitive.ObjectID) (*order.Order, error) {
	o, err := s.orderRepo.FindByID(ctx, id)
	if err != nil {
		return nil, err
	}

	// Validate state transition using state machine
	if err := s.stateMachineManager.ValidateOrderTransition(o, order.EventMarkReady); err != nil {
		return nil, fmt.Errorf("finish preparing validation failed: %w", err)
	}

	now := time.Now()
	o.Status = order.StatusReady
	o.ReadyAt = &now

	if err := s.orderRepo.Update(ctx, id, o); err != nil {
		return nil, err
	}
	return o, nil
}

// ServeOrder - Waiter delivers drink to customer
func (s *OrderService) ServeOrder(ctx context.Context, id primitive.ObjectID) (*order.Order, error) {
	o, err := s.orderRepo.FindByID(ctx, id)
	if err != nil {
		return nil, err
	}

	// Validate state transition using state machine
	if err := s.stateMachineManager.ValidateOrderTransition(o, order.EventServeOrder); err != nil {
		return nil, fmt.Errorf("serve order validation failed: %w", err)
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

	// Validate state transition using state machine
	if err := s.stateMachineManager.ValidateOrderTransition(o, order.EventCancelOrder); err != nil {
		return nil, fmt.Errorf("cancel order validation failed: %w", err)
	}

	o.Status = order.StatusCancelled
	o.CancelReason = req.Reason

	if err := s.orderRepo.Update(ctx, id, o); err != nil {
		return nil, err
	}

	return o, nil
}

// GetQueuedOrders - Get orders waiting for barista
func (s *OrderService) GetQueuedOrders(ctx context.Context) ([]*order.Order, error) {
	return s.orderRepo.FindByStatus(ctx, order.StatusQueued)
}

// GetBaristaOrders - Get orders assigned to a barista
func (s *OrderService) GetBaristaOrders(ctx context.Context, baristaID primitive.ObjectID) ([]*order.Order, error) {
	// Get IN_PROGRESS, READY, and SERVED orders for barista
	inProgress, err := s.orderRepo.FindByStatus(ctx, order.StatusInProgress)
	if err != nil {
		return nil, err
	}
	ready, err := s.orderRepo.FindByStatus(ctx, order.StatusReady)
	if err != nil {
		return nil, err
	}
	served, err := s.orderRepo.FindByStatus(ctx, order.StatusServed)
	if err != nil {
		return nil, err
	}
	
	// Combine and filter by barista
	var result []*order.Order
	for _, o := range inProgress {
		if o.BaristaID == baristaID {
			result = append(result, o)
		}
	}
	for _, o := range ready {
		if o.BaristaID == baristaID {
			result = append(result, o)
		}
	}
	for _, o := range served {
		if o.BaristaID == baristaID {
			result = append(result, o)
		}
	}
	
	return result, nil
}

func (s *OrderService) LockOrder(ctx context.Context, id primitive.ObjectID) (*order.Order, error) {
	o, err := s.orderRepo.FindByID(ctx, id)
	if err != nil {
		return nil, err
	}

	// Validate state transition using state machine
	if err := s.stateMachineManager.ValidateOrderTransition(o, order.EventLockOrder); err != nil {
		return nil, fmt.Errorf("lock order validation failed: %w", err)
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
