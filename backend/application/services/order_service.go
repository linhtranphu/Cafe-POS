package services

import (
	"context"
	"errors"
	"fmt"
	"log"
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
	FindByIDs(ctx context.Context, ids []primitive.ObjectID) ([]*order.Order, error)
	FindAll(ctx context.Context) ([]*order.Order, error)
}

type OrderService struct {
	orderRepo           OrderRepository
	shiftRepo           ShiftRepository
	menuRepo            MenuRepository
	stateMachineManager *domain.StateMachineManager
	batchUsageService   *BatchUsageService
	ingredientService   *IngredientService
	printService        PrintService              // Optional: for auto-printing
	settingsRepo        ShopSettingsRepository    // Optional: for auto-print setting
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
		printService:        nil, // Will be set via SetPrintService
	}
}

// SetIngredientService sets the ingredient service for raw stock deduction
func (s *OrderService) SetIngredientService(ingredientService *IngredientService) {
	s.ingredientService = ingredientService
}

// SetPrintService sets the print service for auto-printing (optional)
func (s *OrderService) SetPrintService(printService PrintService) {
	s.printService = printService
}

// SetSettingsRepository sets the settings repository for auto-print setting (optional)
func (s *OrderService) SetSettingsRepository(settingsRepo ShopSettingsRepository) {
	s.settingsRepo = settingsRepo
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

		// Deduct inventory when order is PAID
		if err := s.deductIngredients(ctx, o); err != nil {
			log.Printf("WARNING: inventory deduction failed for paid order %s: %v", o.OrderNumber, err)
			// Do not rollback payment — log only
		}
	}

	// Update shift revenue based on payment method
	if !o.ShiftID.IsZero() {
		log.Printf("💰 [PAYMENT] Received - ShiftID: %s, Method: %s, Amount: %.0f VND", 
			o.ShiftID.Hex(), req.PaymentMethod, req.Amount)
		
		shift, err := s.shiftRepo.FindByID(ctx, o.ShiftID)
		if err != nil {
			log.Printf("❌ [PAYMENT] Failed to find shift: %v", err)
		} else if shift == nil {
			log.Printf("❌ [PAYMENT] Shift is nil")
		} else {
			log.Printf("✅ [PAYMENT] Found shift - ID: %s", shift.ID.Hex())
			log.Printf("📊 [PAYMENT] BEFORE UPDATE:")
			log.Printf("   - CurrentCash: %.0f VND", shift.CurrentCash)
			log.Printf("   - RemainingCash: %.0f VND", shift.RemainingCash)
			log.Printf("   - TransferRevenue: %.0f VND", shift.TransferRevenue)
			log.Printf("   - RemainingTransfer: %.0f VND", shift.RemainingTransfer)
			log.Printf("   - TotalRevenue: %.0f VND", shift.TotalRevenue)
			
			// Update total revenue for all payment methods
			shift.TotalRevenue += req.Amount
			
			// Update specific payment method fields
			if req.PaymentMethod == order.PaymentCash {
				log.Printf("💵 [PAYMENT] Processing CASH payment")
				shift.RemainingCash += req.Amount
				shift.CurrentCash += req.Amount
			} else if req.PaymentMethod == order.PaymentTransfer || req.PaymentMethod == order.PaymentQR {
				log.Printf("🏦 [PAYMENT] Processing TRANSFER/QR payment")
				shift.TransferRevenue += req.Amount
				shift.RemainingTransfer += req.Amount
			} else {
				log.Printf("⚠️  [PAYMENT] Unknown payment method: %s", req.PaymentMethod)
			}
			
			log.Printf("📊 [PAYMENT] AFTER UPDATE (in memory):")
			log.Printf("   - CurrentCash: %.0f VND", shift.CurrentCash)
			log.Printf("   - RemainingCash: %.0f VND", shift.RemainingCash)
			log.Printf("   - TransferRevenue: %.0f VND", shift.TransferRevenue)
			log.Printf("   - RemainingTransfer: %.0f VND", shift.RemainingTransfer)
			log.Printf("   - TotalRevenue: %.0f VND", shift.TotalRevenue)
			
			// Update shift
			if err := s.shiftRepo.Update(ctx, o.ShiftID, shift); err != nil {
				log.Printf("❌ [PAYMENT] Failed to update shift in DB: %v", err)
			} else {
				log.Printf("✅ [PAYMENT] Shift updated successfully in DB")
				
				// Verify the update by reading back
				verifyShift, verifyErr := s.shiftRepo.FindByID(ctx, o.ShiftID)
				if verifyErr == nil && verifyShift != nil {
					log.Printf("🔍 [PAYMENT] VERIFY (from DB):")
					log.Printf("   - CurrentCash: %.0f VND", verifyShift.CurrentCash)
					log.Printf("   - RemainingCash: %.0f VND", verifyShift.RemainingCash)
					log.Printf("   - TransferRevenue: %.0f VND", verifyShift.TransferRevenue)
					log.Printf("   - RemainingTransfer: %.0f VND", verifyShift.RemainingTransfer)
					log.Printf("   - TotalRevenue: %.0f VND", verifyShift.TotalRevenue)
				} else {
					log.Printf("❌ [PAYMENT] Failed to verify shift update: %v", verifyErr)
				}
			}
		}
	} else {
		log.Printf("⚠️  [PAYMENT] Not updating shift - ShiftID is zero")
	}

	if err := s.orderRepo.Update(ctx, id, o); err != nil {
		return nil, err
	}

	// Emit OrderCreated event for printing if order is now PAID
	// This happens after order is committed to ensure order persistence before print jobs
	if o.Status == order.StatusPaid && s.printService != nil {
		// Check auto-print setting before creating print jobs
		autoPrintEnabled := true // Default to true if settings not available
		if s.settingsRepo != nil {
			settings, err := s.settingsRepo.FindFirst(ctx)
			if err == nil && settings != nil {
				autoPrintEnabled = settings.AutoPrintEnabled
			}
		}

		if autoPrintEnabled {
			// Call print service asynchronously to not block order creation
			// Errors in print job creation should not affect order creation
			go func() {
				// Use background context to avoid cancellation
				printCtx := context.Background()
				if err := s.printService.CreatePrintJobsForOrder(printCtx, o); err != nil {
					// Log error but don't fail the order creation
					fmt.Printf("ERROR: Failed to create print jobs for order %s: %v\n", o.OrderNumber, err)
				} else {
					fmt.Printf("INFO: Print jobs created for order %s\n", o.OrderNumber)
				}
			}()
		} else {
			fmt.Printf("INFO: Auto-print disabled, skipping print jobs for order %s\n", o.OrderNumber)
		}
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

// deductIngredients deducts inventory for all ingredients with DeductInventory=true when order is PAID.
// Batch ingredients are deducted via BatchUsageService; raw ingredients via IngredientService.DeductStock.
func (s *OrderService) deductIngredients(ctx context.Context, o *order.Order) error {
	for _, item := range o.Items {
		menuItem, err := s.menuRepo.FindByID(ctx, item.MenuItemID)
		if err != nil {
			log.Printf("WARNING: failed to fetch menu item %s for deduction: %v", item.MenuItemID.Hex(), err)
			continue
		}

		ingredients := menuItem.GetIngredients(item.VariantID)

		for _, ing := range ingredients {
			if !ing.DeductInventory {
				continue
			}

			qty := ing.Quantity * float64(item.Quantity)

			switch {
			case ing.IsBatchIngredient() && ing.BatchID != nil && s.batchUsageService != nil:
				req := UseBatchRequest{
					BatchDefinitionID: *ing.BatchID,
					QuantityNeeded:    qty,
					Unit:              string(ing.Unit),
					OrderID:           o.ID,
					MenuItemID:        item.MenuItemID,
					MenuItemName:      item.Name,
				}
				result, err := s.batchUsageService.UseBatch(ctx, req)
				if err != nil {
					log.Printf("WARNING: failed to deduct batch ingredient %s for order %s: %v", ing.Name, o.OrderNumber, err)
				} else if !result.Success {
					log.Printf("WARNING: insufficient batch %s for order %s: %s", ing.Name, o.OrderNumber, result.Message)
				}

			case ing.IsRawIngredient() && ing.IngredientID != nil && s.ingredientService != nil:
				if err := s.ingredientService.DeductStock(ctx, *ing.IngredientID, qty, ing.Unit, o.ID, o.OrderNumber); err != nil {
					log.Printf("WARNING: failed to deduct raw ingredient %s for order %s: %v", ing.Name, o.OrderNumber, err)
				}
			}
		}
	}

	return nil
}

// PrintTemporaryBill - Mark order as bill printed (for temporary bill printing)
func (s *OrderService) PrintTemporaryBill(ctx context.Context, id primitive.ObjectID) (*order.Order, error) {
	o, err := s.orderRepo.FindByID(ctx, id)
	if err != nil {
		return nil, fmt.Errorf("order not found: %w", err)
	}

	log.Printf("[ORDER] PrintTemporaryBill called for order %s (ID: %s)", o.OrderNumber, o.ID.Hex())

	// Only allow printing bill for CREATED status orders
	if o.Status != order.StatusCreated {
		return nil, fmt.Errorf("can only print temporary bill for unpaid orders")
	}

	// Mark as bill printed
	o.BillPrinted = true
	o.UpdatedAt = time.Now()

	if err := s.orderRepo.Update(ctx, o.ID, o); err != nil {
		return nil, fmt.Errorf("failed to update order: %w", err)
	}

	log.Printf("[ORDER] Order %s marked as bill_printed=true", o.OrderNumber)

	// Trigger print job if print service is available
	if s.printService != nil {
		log.Printf("[ORDER] Calling CreateTempBillJob for order %s", o.OrderNumber)
		if err := s.printService.CreateTempBillJob(ctx, o); err != nil {
			// Log error but don't fail the operation
			log.Printf("[ORDER ERROR] Failed to print temp bill for order %s: %v", o.OrderNumber, err)
			fmt.Printf("Warning: Failed to print temp bill for order %s: %v\n", o.OrderNumber, err)
		} else {
			log.Printf("[ORDER] CreateTempBillJob completed successfully for order %s", o.OrderNumber)
		}
	} else {
		log.Printf("[ORDER WARNING] Print service is nil, cannot create temp bill job")
	}

	return o, nil
}

// MergeOrders - Merge multiple unpaid orders into one
func (s *OrderService) MergeOrders(ctx context.Context, req *order.MergeOrdersRequest, userID, userName string) (*order.MergeOrdersResponse, error) {
	// Validate minimum orders
	if len(req.OrderIDs) < 2 {
		return nil, errors.New("phải chọn ít nhất 2 orders để gộp")
	}

	// Convert string IDs to ObjectIDs
	orderIDs := make([]primitive.ObjectID, 0, len(req.OrderIDs))
	for _, idStr := range req.OrderIDs {
		id, err := primitive.ObjectIDFromHex(idStr)
		if err != nil {
			return nil, fmt.Errorf("invalid order ID: %s", idStr)
		}
		orderIDs = append(orderIDs, id)
	}

	// Fetch all orders
	orders, err := s.orderRepo.FindByIDs(ctx, orderIDs)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch orders: %w", err)
	}

	// Validate all orders exist
	if len(orders) != len(orderIDs) {
		return nil, errors.New("một số orders không tồn tại")
	}

	// Validate all orders are mergeable
	var shiftID primitive.ObjectID
	hasPaidOrder := false
	totalAmountPaid := 0.0
	allItems := []order.OrderItem{}
	totalDiscount := 0.0
	cancelledOrderNumbers := []string{}

	for i, o := range orders {
		// Check if order is mergeable (CREATED or PAID only)
		if !o.IsMergeable() {
			return nil, fmt.Errorf("order %s không thể gộp (status: %s)", o.OrderNumber, o.Status)
		}

		// First order sets the shift ID
		if i == 0 {
			shiftID = o.ShiftID
		} else {
			// All orders must be in the same shift
			if o.ShiftID != shiftID {
				return nil, errors.New("các orders phải thuộc cùng một ca làm việc")
			}
		}

		// Check if any order is PAID
		if o.Status == order.StatusPaid {
			hasPaidOrder = true
		}

		// Accumulate amount paid
		totalAmountPaid += o.AmountPaid

		// Collect all items
		allItems = append(allItems, o.Items...)

		// Accumulate discounts
		totalDiscount += o.Discount

		// Track order numbers for cancellation
		cancelledOrderNumbers = append(cancelledOrderNumbers, o.OrderNumber)
	}

	// Validate shift is still open
	shift, err := s.shiftRepo.FindByID(ctx, shiftID)
	if err != nil || shift.Status != order.ShiftOpen {
		return nil, errors.New("ca làm việc đã đóng, không thể gộp orders")
	}

	// Generate new order number
	now := time.Now()
	newOrderNumber := fmt.Sprintf("%s-%03d", now.Format("20060102-150405"), now.Nanosecond()/1000000%1000)

	// Determine customer name (use provided or from first order)
	customerName := req.CustomerName
	if customerName == "" {
		customerName = orders[0].CustomerName
	}

	// Determine waiter (use first order's waiter)
	waiterID := orders[0].WaiterID
	waiterName := orders[0].WaiterName

	// Determine status and payment info
	var newStatus order.OrderStatus
	var paymentMethod order.PaymentMethod
	var collectorID primitive.ObjectID
	var collectorName string
	var paidAt *time.Time

	if hasPaidOrder {
		// If any order is PAID, merged order is PAID
		newStatus = order.StatusPaid
		// Use payment info from first PAID order
		for _, o := range orders {
			if o.Status == order.StatusPaid {
				paymentMethod = o.PaymentMethod
				collectorID = o.CollectorID
				collectorName = o.CollectorName
				paidAt = o.PaidAt
				break
			}
		}
	} else {
		// All orders are CREATED
		newStatus = order.StatusCreated
	}

	// Create merged order
	mergedOrder := &order.Order{
		OrderNumber:   newOrderNumber,
		CustomerName:  customerName,
		WaiterID:      waiterID,
		WaiterName:    waiterName,
		ShiftID:       shiftID,
		Items:         allItems,
		Discount:      totalDiscount,
		Status:        newStatus,
		AmountPaid:    totalAmountPaid,
		Note:          req.Note,
		PaymentMethod: paymentMethod,
		CollectorID:   collectorID,
		CollectorName: collectorName,
		PaidAt:        paidAt,
		BillPrinted:   false, // New order hasn't been printed yet
	}

	// Calculate totals
	mergedOrder.CalculateTotal()

	// Create the merged order
	if err := s.orderRepo.Create(ctx, mergedOrder); err != nil {
		return nil, fmt.Errorf("failed to create merged order: %w", err)
	}

	log.Printf("[ORDER] Created merged order %s from %d orders", mergedOrder.OrderNumber, len(orders))

	// Cancel all original orders
	cancelReason := fmt.Sprintf("Đã gộp vào order #%s", mergedOrder.OrderNumber)
	for _, o := range orders {
		o.Status = order.StatusCancelled
		o.CancelReason = cancelReason
		o.UpdatedAt = time.Now()

		if err := s.orderRepo.Update(ctx, o.ID, o); err != nil {
			log.Printf("[ORDER ERROR] Failed to cancel order %s: %v", o.OrderNumber, err)
			// Continue with other orders even if one fails
		} else {
			log.Printf("[ORDER] Cancelled order %s (merged into %s)", o.OrderNumber, mergedOrder.OrderNumber)
		}
	}

	return &order.MergeOrdersResponse{
		MergedOrder:     mergedOrder,
		CancelledOrders: cancelledOrderNumbers,
		Message:         fmt.Sprintf("Đã gộp %d orders thành order #%s", len(orders), mergedOrder.OrderNumber),
	}, nil
}
