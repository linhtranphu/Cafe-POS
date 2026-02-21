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
	Delete(ctx context.Context, id primitive.ObjectID) error
	FindByShiftID(ctx context.Context, shiftID primitive.ObjectID) ([]*order.Order, error)
	FindByWaiterID(ctx context.Context, waiterID primitive.ObjectID) ([]*order.Order, error)
	FindByStatus(ctx context.Context, status order.OrderStatus) ([]*order.Order, error)
	FindByOrderNumber(ctx context.Context, orderNumber string) (*order.Order, error)
	FindAll(ctx context.Context) ([]*order.Order, error)
}

type OrderService struct {
	orderRepo           OrderRepository
	shiftRepo           ShiftRepository
	menuRepo            MenuRepository
	stateMachineManager *domain.StateMachineManager
	batchUsageService   *BatchUsageService
}

func NewOrderService(
	orderRepo OrderRepository,
	shiftRepo ShiftRepository,
	menuRepo MenuRepository,
	stateMachineManager *domain.StateMachineManager,
	batchUsageService *BatchUsageService,
) *OrderService {
	return &OrderService{
		orderRepo:           orderRepo,
		shiftRepo:           shiftRepo,
		menuRepo:            menuRepo,
		stateMachineManager: stateMachineManager,
		batchUsageService:   batchUsageService,
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
		Items:        []order.OrderItem{}, // Initialize empty, will populate below
		Status:       order.StatusCreated,
		Note:         req.Note,
		AmountPaid:   0,
	}

	// Validate and populate order items with prices from menu
	for _, item := range req.Items {
		// Fetch menu item to get price and validate variant
		menuItem, err := s.menuRepo.FindByID(ctx, item.MenuItemID)
		if err != nil {
			return nil, fmt.Errorf("menu item not found: %w", err)
		}

		// Check availability
		if !menuItem.Available {
			return nil, fmt.Errorf("menu item %s is not available", menuItem.Name)
		}

		// Determine price and variant name based on menu item type
		var price float64
		var variantName string
		var variantID string

		if menuItem.HasVariants {
			// Multi-size item - variant_id required
			if item.VariantID == "" {
				return nil, fmt.Errorf("variant_id required for multi-size item: %s", menuItem.Name)
			}

			variant := menuItem.GetVariantByID(item.VariantID)
			if variant == nil {
				return nil, fmt.Errorf("invalid variant_id '%s' for item: %s", item.VariantID, menuItem.Name)
			}

			if !variant.Available {
				return nil, fmt.Errorf("variant %s is not available", variant.Name)
			}

			price = variant.Price
			variantName = variant.Name
			variantID = variant.ID
		} else {
			// Single-size item (backward compatible)
			// No variant_id needed, use item.price directly
			price = menuItem.Price
			variantName = "" // No variant for single-size
			variantID = ""
		}

		// Create order item with validated data
		orderItem := order.OrderItem{
			MenuItemID:  item.MenuItemID,
			VariantID:   variantID,
			Name:        menuItem.Name,
			VariantName: variantName,
			Price:       price,
			Quantity:    item.Quantity,
			Note:        item.Note,
		}

		o.Items = append(o.Items, orderItem)
	}

	o.CalculateTotal()
	
	// Create order first
	if err := s.orderRepo.Create(ctx, o); err != nil {
		return nil, err
	}

	// Deduct batch ingredients after order is created
	// This happens after order creation so we have an order ID for logging
	if s.batchUsageService != nil {
		batchCost, err := s.deductBatchIngredients(ctx, o)
		if err != nil {
			// Rollback: restore batch quantities and delete the order
			_ = s.batchUsageService.RollbackBatchUsage(ctx, o.ID)
			_ = s.orderRepo.Delete(ctx, o.ID)
			return nil, fmt.Errorf("failed to deduct batch ingredients: %w", err)
		}
		
		// Store batch cost information in order note for tracking
		// In a production system, you might want to add a dedicated field for this
		if batchCost > 0 {
			if o.Note != "" {
				o.Note += fmt.Sprintf(" [Batch Cost: %.2f VND]", batchCost)
			} else {
				o.Note = fmt.Sprintf("[Batch Cost: %.2f VND]", batchCost)
			}
			// Update order with batch cost info
			_ = s.orderRepo.Update(ctx, o.ID, o)
		}
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

	// Update shift cash if payment is cash and order has shift_id
	if req.PaymentMethod == order.PaymentCash && !o.ShiftID.IsZero() {
		fmt.Printf("DEBUG: Updating shift cash - ShiftID: %s, Amount: %.2f\n", o.ShiftID.Hex(), req.Amount)
		shift, err := s.shiftRepo.FindByID(ctx, o.ShiftID)
		if err == nil && shift != nil {
			fmt.Printf("DEBUG: Found shift - Current RemainingCash: %.2f\n", shift.RemainingCash)
			// Add cash to shift
			shift.RemainingCash += req.Amount
			shift.CurrentCash += req.Amount
			shift.TotalRevenue += req.Amount
			fmt.Printf("DEBUG: New RemainingCash: %.2f\n", shift.RemainingCash)
			
			// Update shift
			if err := s.shiftRepo.Update(ctx, o.ShiftID, shift); err != nil {
				// Log error but don't fail the payment
				fmt.Printf("ERROR: Failed to update shift cash: %v\n", err)
			} else {
				fmt.Printf("DEBUG: Shift cash updated successfully\n")
			}
		} else {
			fmt.Printf("DEBUG: Shift not found or error: %v\n", err)
		}
	} else {
		fmt.Printf("DEBUG: Not updating shift - PaymentMethod: %s, ShiftID.IsZero: %v\n", req.PaymentMethod, o.ShiftID.IsZero())
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

	// Validate and populate new items with prices from menu
	var validatedItems []order.OrderItem
	for _, item := range req.Items {
		// Fetch menu item to get price and validate variant
		menuItem, err := s.menuRepo.FindByID(ctx, item.MenuItemID)
		if err != nil {
			return nil, fmt.Errorf("menu item not found: %w", err)
		}

		// Check availability
		if !menuItem.Available {
			return nil, fmt.Errorf("menu item %s is not available", menuItem.Name)
		}

		// Determine price and variant name based on menu item type
		var price float64
		var variantName string
		var variantID string

		if menuItem.HasVariants {
			// Multi-size item - variant_id required
			if item.VariantID == "" {
				return nil, fmt.Errorf("variant_id required for multi-size item: %s", menuItem.Name)
			}

			variant := menuItem.GetVariantByID(item.VariantID)
			if variant == nil {
				return nil, fmt.Errorf("invalid variant_id '%s' for item: %s", item.VariantID, menuItem.Name)
			}

			if !variant.Available {
				return nil, fmt.Errorf("variant %s is not available", variant.Name)
			}

			price = variant.Price
			variantName = variant.Name
			variantID = variant.ID
		} else {
			// Single-size item (backward compatible)
			price = menuItem.Price
			variantName = ""
			variantID = ""
		}

		// Create order item with validated data
		orderItem := order.OrderItem{
			MenuItemID:  item.MenuItemID,
			VariantID:   variantID,
			Name:        menuItem.Name,
			VariantName: variantName,
			Price:       price,
			Quantity:    item.Quantity,
			Note:        item.Note,
		}

		validatedItems = append(validatedItems, orderItem)
	}

	// Update order details
	o.Items = validatedItems
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

// deductBatchIngredients deducts batch ingredients for all items in an order
// Returns the total batch cost used
func (s *OrderService) deductBatchIngredients(ctx context.Context, o *order.Order) (float64, error) {
	totalBatchCost := 0.0
	
	for _, item := range o.Items {
		// Get menu item to check ingredients
		menuItem, err := s.menuRepo.FindByID(ctx, item.MenuItemID)
		if err != nil {
			return 0, fmt.Errorf("failed to fetch menu item %s: %w", item.MenuItemID.Hex(), err)
		}

		// Get ingredients for this item (considering variants)
		ingredients := menuItem.GetIngredients(item.VariantID)

		// Process each ingredient
		for _, ing := range ingredients {
			// Only process batch ingredients
			if !ing.IsBatchIngredient() {
				continue
			}

			// Skip if no batch ID
			if ing.BatchID == nil {
				continue
			}

			// Calculate total quantity needed (ingredient quantity * order item quantity)
			quantityNeeded := ing.Quantity * float64(item.Quantity)

			// Deduct batch using BatchUsageService
			req := UseBatchRequest{
				BatchDefinitionID: *ing.BatchID,
				QuantityNeeded:    quantityNeeded,
				Unit:              string(ing.Unit), // Convert UnitType to string
				OrderID:           o.ID,
				MenuItemID:        item.MenuItemID,
				MenuItemName:      item.Name,
			}

			result, err := s.batchUsageService.UseBatch(ctx, req)
			if err != nil {
				return 0, fmt.Errorf("failed to use batch for ingredient %s: %w", ing.Name, err)
			}

			if !result.Success {
				return 0, fmt.Errorf("insufficient batch %s: %s", ing.Name, result.Message)
			}
			
			// Accumulate batch cost
			totalBatchCost += result.TotalCost
		}
	}

	return totalBatchCost, nil
}
