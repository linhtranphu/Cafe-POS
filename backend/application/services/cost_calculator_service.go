package services

import (
	"context"
	"fmt"
	"math"
	"time"

	"cafe-pos/backend/domain/ingredient"
	"cafe-pos/backend/domain/menu"
	"cafe-pos/backend/domain/order"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

// OrderItemRepository interface for order items with cost tracking
type OrderItemRepository interface {
	CreateMany(ctx context.Context, items []*order.OrderItemWithCost) error
	FindByOrderIDs(ctx context.Context, orderIDs []primitive.ObjectID) ([]*order.OrderItemWithCost, error)
	FindByDateRange(ctx context.Context, startDate, endDate time.Time) ([]*order.OrderItemWithCost, error)
}

// CostCalculatorService handles cost calculation for menu items
type CostCalculatorService struct {
	menuRepo       MenuRepository
	ingredientRepo IngredientRepository
	orderRepo      OrderRepository
	orderItemRepo  OrderItemRepository
	
	// Reference to cost recalculation service for queuing background jobs
	recalcService *CostRecalculationService
	
	// Monitoring service for metrics and alerts
	monitoringService *MonitoringService
}

// NewCostCalculatorService creates a new cost calculator service
func NewCostCalculatorService(menuRepo MenuRepository, ingredientRepo IngredientRepository, orderRepo OrderRepository, orderItemRepo OrderItemRepository) *CostCalculatorService {
	return &CostCalculatorService{
		menuRepo:       menuRepo,
		ingredientRepo: ingredientRepo,
		orderRepo:      orderRepo,
		orderItemRepo:  orderItemRepo,
	}
}

// SetCostRecalculationService sets the cost recalculation service for queuing background jobs
// This is called after service initialization to avoid circular dependencies
func (s *CostCalculatorService) SetCostRecalculationService(recalcService *CostRecalculationService) {
	s.recalcService = recalcService
}

// SetMonitoringService sets the monitoring service for metrics collection
func (s *CostCalculatorService) SetMonitoringService(monitoringService *MonitoringService) {
	s.monitoringService = monitoringService
}

// MenuItemCostResult represents the result of a cost calculation
type MenuItemCostResult struct {
	MenuItemID           primitive.ObjectID
	CurrentCost          float64
	CostStatus           menu.CostStatus
	CostLastCalculatedAt time.Time
	MissingIngredients   []string // Names of ingredients with missing cost_per_unit
}

// CostCalculationSummary represents summary statistics for batch cost calculation
type CostCalculationSummary struct {
	TotalItems      int
	SuccessfulItems int
	IncompleteItems int
	FailedItems     int
	AverageCost     float64
}

// CalculateMenuItemCost calculates the current cost for a single menu item
// Formula: sum(quantity * cost_per_unit * conversion_rate * (1 + wastage_percentage/100))
// Returns INCOMPLETE status if any ingredient has missing cost_per_unit
func (s *CostCalculatorService) CalculateMenuItemCost(ctx context.Context, menuItemID primitive.ObjectID) (*MenuItemCostResult, error) {
	startTime := time.Now()
	
	// Fetch the menu item
	menuItem, err := s.menuRepo.FindByID(ctx, menuItemID)
	if err != nil {
		// Record failure metric
		if s.monitoringService != nil {
			s.monitoringService.RecordMetric(Metric{
				Type:      MetricTypeCostCalculation,
				Status:    MetricStatusFailure,
				Timestamp: time.Now(),
				Duration:  time.Since(startTime),
				Message:   fmt.Sprintf("Failed to fetch menu item: %v", err),
				Metadata: map[string]interface{}{
					"menu_item_id": menuItemID.Hex(),
					"error":        err.Error(),
				},
			})
		}
		return nil, fmt.Errorf("failed to fetch menu item: %w", err)
	}

	// If menu item has no ingredients, return zero cost with FINAL status
	if len(menuItem.Ingredients) == 0 {
		result := &MenuItemCostResult{
			MenuItemID:           menuItemID,
			CurrentCost:          0.0,
			CostStatus:           menu.CostStatusFinal,
			CostLastCalculatedAt: time.Now(),
			MissingIngredients:   []string{},
		}
		
		// Record success metric
		if s.monitoringService != nil {
			s.monitoringService.RecordMetric(Metric{
				Type:      MetricTypeCostCalculation,
				Status:    MetricStatusSuccess,
				Timestamp: time.Now(),
				Duration:  time.Since(startTime),
				Message:   "Cost calculation completed (no ingredients)",
				Metadata: map[string]interface{}{
					"menu_item_id": menuItemID.Hex(),
					"cost":         0.0,
				},
			})
		}
		
		return result, nil
	}

	// Fetch all ingredients to build a lookup map
	allIngredients, err := s.ingredientRepo.FindAll(ctx)
	if err != nil {
		// Record failure metric
		if s.monitoringService != nil {
			s.monitoringService.RecordMetric(Metric{
				Type:      MetricTypeCostCalculation,
				Status:    MetricStatusFailure,
				Timestamp: time.Now(),
				Duration:  time.Since(startTime),
				Message:   fmt.Sprintf("Failed to fetch ingredients: %v", err),
				Metadata: map[string]interface{}{
					"menu_item_id": menuItemID.Hex(),
					"error":        err.Error(),
				},
			})
		}
		return nil, fmt.Errorf("failed to fetch ingredients: %w", err)
	}

	// Build ingredient lookup map by name
	ingredientMap := make(map[string]*ingredient.Ingredient)
	for _, ing := range allIngredients {
		ingredientMap[ing.Name] = ing
	}

	// Calculate total cost
	var totalCost float64
	var missingIngredients []string
	hasIncompleteCost := false

	for _, menuIngredient := range menuItem.Ingredients {
		// Find the ingredient in our map
		ing, exists := ingredientMap[menuIngredient.Name]
		if !exists {
			// Ingredient not found in database
			missingIngredients = append(missingIngredients, menuIngredient.Name)
			hasIncompleteCost = true
			continue
		}

		// Check if cost_per_unit is missing (0 or not set)
		if ing.CostPerUnit <= 0 {
			missingIngredients = append(missingIngredients, menuIngredient.Name)
			hasIncompleteCost = true
			continue
		}

		// Calculate conversion rate dynamically based on stock unit and recipe unit
		// stockUnit = ing.Unit (e.g., "L")
		// recipeUnit = menuIngredient.Unit (e.g., "ml")
		conversionRate := ingredient.GetConversionRate(ing.Unit, menuIngredient.Unit)

		// Get wastage percentage (default 0.0 if not set)
		wastagePercentage := ing.WastagePercentage
		if wastagePercentage < 0 {
			wastagePercentage = 0.0
		}

		// Calculate cost for this ingredient
		// Formula: quantity * cost_per_unit * conversion_rate * (1 + wastage_percentage/100)
		ingredientCost := menuIngredient.Quantity * ing.CostPerUnit * conversionRate * (1 + wastagePercentage/100)
		totalCost += ingredientCost
	}

	// Round to 2 decimal places
	totalCost = math.Round(totalCost*100) / 100

	// Determine cost status
	costStatus := menu.CostStatusFinal
	metricStatus := MetricStatusSuccess
	if hasIncompleteCost {
		costStatus = menu.CostStatusIncomplete
		metricStatus = MetricStatusWarning
	}
	
	result := &MenuItemCostResult{
		MenuItemID:           menuItemID,
		CurrentCost:          totalCost,
		CostStatus:           costStatus,
		CostLastCalculatedAt: time.Now(),
		MissingIngredients:   missingIngredients,
	}
	
	// Record metric
	if s.monitoringService != nil {
		s.monitoringService.RecordMetric(Metric{
			Type:      MetricTypeCostCalculation,
			Status:    metricStatus,
			Timestamp: time.Now(),
			Duration:  time.Since(startTime),
			Message:   "Cost calculation completed",
			Metadata: map[string]interface{}{
				"menu_item_id":        menuItemID.Hex(),
				"cost":                totalCost,
				"cost_status":         string(costStatus),
				"missing_ingredients": len(missingIngredients),
			},
		})
	}

	return result, nil
}

// CalculateAllMenuItemCosts calculates current cost for all menu items
// Batch fetches all menu items and ingredients, calculates cost for each item,
// updates the menu_items collection with current_cost, and returns summary statistics
func (s *CostCalculatorService) CalculateAllMenuItemCosts(ctx context.Context) (*CostCalculationSummary, error) {
	// Batch fetch all menu items
	menuItems, err := s.menuRepo.FindAll(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch menu items: %w", err)
	}

	// Batch fetch all ingredients once for efficiency
	allIngredients, err := s.ingredientRepo.FindAll(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch ingredients: %w", err)
	}

	// Build ingredient lookup map by name
	ingredientMap := make(map[string]*ingredient.Ingredient)
	for _, ing := range allIngredients {
		ingredientMap[ing.Name] = ing
	}

	// Initialize summary statistics
	summary := &CostCalculationSummary{
		TotalItems: len(menuItems),
	}

	var totalCostSum float64

	// Calculate cost for each menu item
	for _, menuItem := range menuItems {
		// Calculate cost using the same logic as CalculateMenuItemCost
		costResult := s.calculateCostForMenuItem(menuItem, ingredientMap)

		// Update the menu item with calculated cost
		menuItem.CurrentCost = costResult.CurrentCost
		menuItem.CostStatus = costResult.CostStatus
		menuItem.CostLastCalculatedAt = costResult.CostLastCalculatedAt

		// Update in database
		err := s.menuRepo.Update(ctx, menuItem.ID, menuItem)
		if err != nil {
			// Log error but continue with other items
			summary.FailedItems++
			continue
		}

		// Update summary statistics
		if costResult.CostStatus == menu.CostStatusIncomplete {
			summary.IncompleteItems++
		} else {
			summary.SuccessfulItems++
		}

		totalCostSum += costResult.CurrentCost
	}

	// Calculate average cost (excluding failed items)
	processedItems := summary.SuccessfulItems + summary.IncompleteItems
	if processedItems > 0 {
		summary.AverageCost = math.Round((totalCostSum/float64(processedItems))*100) / 100
	}

	return summary, nil
}

// calculateCostForMenuItem is a helper method that calculates cost for a menu item
// given a pre-built ingredient map (for efficiency in batch operations)
func (s *CostCalculatorService) calculateCostForMenuItem(menuItem *menu.MenuItem, ingredientMap map[string]*ingredient.Ingredient) *MenuItemCostResult {
	result := &MenuItemCostResult{
		MenuItemID:           menuItem.ID,
		CurrentCost:          0.0,
		CostStatus:           menu.CostStatusFinal,
		CostLastCalculatedAt: time.Now(),
		MissingIngredients:   []string{},
	}

	// If menu item has no ingredients, return zero cost with FINAL status
	if len(menuItem.Ingredients) == 0 {
		return result
	}

	var totalCost float64
	hasIncompleteCost := false

	for _, menuIngredient := range menuItem.Ingredients {
		// Find the ingredient in our map
		ing, exists := ingredientMap[menuIngredient.Name]
		if !exists {
			// Ingredient not found in database
			result.MissingIngredients = append(result.MissingIngredients, menuIngredient.Name)
			hasIncompleteCost = true
			continue
		}

		// Check if cost_per_unit is missing (0 or not set)
		if ing.CostPerUnit <= 0 {
			result.MissingIngredients = append(result.MissingIngredients, menuIngredient.Name)
			hasIncompleteCost = true
			continue
		}

		// Calculate conversion rate dynamically based on stock unit and recipe unit
		conversionRate := ingredient.GetConversionRate(ing.Unit, menuIngredient.Unit)

		// Get wastage percentage (default 0.0 if not set)
		wastagePercentage := ing.WastagePercentage
		if wastagePercentage < 0 {
			wastagePercentage = 0.0
		}

		// Calculate cost for this ingredient
		// Formula: quantity * cost_per_unit * conversion_rate * (1 + wastage_percentage/100)
		ingredientCost := menuIngredient.Quantity * ing.CostPerUnit * conversionRate * (1 + wastagePercentage/100)
		totalCost += ingredientCost
	}

	// Round to 2 decimal places
	result.CurrentCost = math.Round(totalCost*100) / 100

	// Determine cost status
	if hasIncompleteCost {
		result.CostStatus = menu.CostStatusIncomplete
	}

	return result
}

// ShiftCostCalculationResult represents the result of shift order cost calculation
type ShiftCostCalculationResult struct {
	TotalOrders         int
	TotalItems          int
	ItemsWithFinalCost  int
	ItemsWithIncompleteCost int
	TotalAccountingCost float64
}

// CalculateShiftOrderCosts calculates accounting cost for all orders in a shift
// This is called during shift closure to snapshot the cost at the time of closing
// Requirements: 5.2, 5.3, 5.4
func (s *CostCalculatorService) CalculateShiftOrderCosts(ctx context.Context, shiftID primitive.ObjectID) (*ShiftCostCalculationResult, error) {
	startTime := time.Now()
	
	// Fetch all orders in the shift
	orders, err := s.orderRepo.FindByShiftID(ctx, shiftID)
	if err != nil {
		// Record failure metric
		if s.monitoringService != nil {
			s.monitoringService.RecordMetric(Metric{
				Type:      MetricTypeShiftClosure,
				Status:    MetricStatusFailure,
				Timestamp: time.Now(),
				Duration:  time.Since(startTime),
				Message:   fmt.Sprintf("Failed to fetch orders for shift: %v", err),
				Metadata: map[string]interface{}{
					"shift_id": shiftID.Hex(),
					"error":    err.Error(),
				},
			})
		}
		return nil, fmt.Errorf("failed to fetch orders for shift: %w", err)
	}

	// If no orders, return early
	if len(orders) == 0 {
		result := &ShiftCostCalculationResult{
			TotalOrders: 0,
			TotalItems:  0,
		}
		
		// Record success metric
		if s.monitoringService != nil {
			s.monitoringService.RecordMetric(Metric{
				Type:      MetricTypeShiftClosure,
				Status:    MetricStatusSuccess,
				Timestamp: time.Now(),
				Duration:  time.Since(startTime),
				Message:   "Shift closure cost calculation completed (no orders)",
				Metadata: map[string]interface{}{
					"shift_id": shiftID.Hex(),
				},
			})
		}
		
		return result, nil
	}

	// Fetch all ingredients once for efficiency
	allIngredients, err := s.ingredientRepo.FindAll(ctx)
	if err != nil {
		// Record failure metric
		if s.monitoringService != nil {
			s.monitoringService.RecordMetric(Metric{
				Type:      MetricTypeShiftClosure,
				Status:    MetricStatusFailure,
				Timestamp: time.Now(),
				Duration:  time.Since(startTime),
				Message:   fmt.Sprintf("Failed to fetch ingredients: %v", err),
				Metadata: map[string]interface{}{
					"shift_id": shiftID.Hex(),
					"error":    err.Error(),
				},
			})
		}
		return nil, fmt.Errorf("failed to fetch ingredients: %w", err)
	}

	// Build ingredient lookup map by name
	ingredientMap := make(map[string]*ingredient.Ingredient)
	for _, ing := range allIngredients {
		ingredientMap[ing.Name] = ing
	}

	// Fetch all menu items once for efficiency
	allMenuItems, err := s.menuRepo.FindAll(ctx)
	if err != nil {
		// Record failure metric
		if s.monitoringService != nil {
			s.monitoringService.RecordMetric(Metric{
				Type:      MetricTypeShiftClosure,
				Status:    MetricStatusFailure,
				Timestamp: time.Now(),
				Duration:  time.Since(startTime),
				Message:   fmt.Sprintf("Failed to fetch menu items: %v", err),
				Metadata: map[string]interface{}{
					"shift_id": shiftID.Hex(),
					"error":    err.Error(),
				},
			})
		}
		return nil, fmt.Errorf("failed to fetch menu items: %w", err)
	}

	// Build menu item lookup map by ID
	menuItemMap := make(map[primitive.ObjectID]*menu.MenuItem)
	for _, item := range allMenuItems {
		menuItemMap[item.ID] = item
	}

	// Initialize result
	result := &ShiftCostCalculationResult{
		TotalOrders: len(orders),
	}

	// Collect all order items with calculated costs
	var orderItemsWithCost []*order.OrderItemWithCost

	// Process each order
	for _, ord := range orders {
		// Process each item in the order
		for _, item := range ord.Items {
			result.TotalItems++

			// Find the menu item
			menuItem, exists := menuItemMap[item.MenuItemID]
			if !exists {
				// Menu item not found - mark as incomplete
				orderItemWithCost := &order.OrderItemWithCost{
					OrderID:          ord.ID,
					MenuItemID:       item.MenuItemID,
					Name:             item.Name,
					Price:            item.Price,
					Quantity:         item.Quantity,
					Note:             item.Note,
					Subtotal:         item.Subtotal,
					AccountingCost:   0.0,
					CostCalculatedAt: time.Now(),
					CostStatus:       order.CostStatusIncomplete,
				}
				orderItemsWithCost = append(orderItemsWithCost, orderItemWithCost)
				result.ItemsWithIncompleteCost++
				continue
			}

			// Calculate cost for this menu item using current ingredient costs
			costResult := s.calculateCostForMenuItem(menuItem, ingredientMap)

			// Calculate total accounting cost for this order item (cost per unit * quantity)
			accountingCostPerItem := costResult.CurrentCost
			totalAccountingCost := accountingCostPerItem * float64(item.Quantity)
			totalAccountingCost = math.Round(totalAccountingCost*100) / 100

			// Convert menu.CostStatus to order.CostStatus
			var orderCostStatus order.CostStatus
			switch costResult.CostStatus {
			case menu.CostStatusFinal:
				orderCostStatus = order.CostStatusFinal
			case menu.CostStatusEstimated:
				orderCostStatus = order.CostStatusEstimated
			case menu.CostStatusIncomplete:
				orderCostStatus = order.CostStatusIncomplete
			default:
				orderCostStatus = order.CostStatusIncomplete
			}

			// Create OrderItemWithCost
			orderItemWithCost := &order.OrderItemWithCost{
				OrderID:          ord.ID,
				MenuItemID:       item.MenuItemID,
				Name:             item.Name,
				Price:            item.Price,
				Quantity:         item.Quantity,
				Note:             item.Note,
				Subtotal:         item.Subtotal,
				AccountingCost:   totalAccountingCost,
				CostCalculatedAt: time.Now(),
				CostStatus:       orderCostStatus,
			}

			orderItemsWithCost = append(orderItemsWithCost, orderItemWithCost)

			// Update result statistics
			if orderCostStatus == order.CostStatusFinal {
				result.ItemsWithFinalCost++
				result.TotalAccountingCost += totalAccountingCost
			} else {
				result.ItemsWithIncompleteCost++
			}
		}
	}

	// Save all order items with costs to the order_items collection
	if len(orderItemsWithCost) > 0 {
		err = s.orderItemRepo.CreateMany(ctx, orderItemsWithCost)
		if err != nil {
			// Record failure metric
			if s.monitoringService != nil {
				s.monitoringService.RecordMetric(Metric{
					Type:      MetricTypeShiftClosure,
					Status:    MetricStatusFailure,
					Timestamp: time.Now(),
					Duration:  time.Since(startTime),
					Message:   fmt.Sprintf("Failed to save order items with costs: %v", err),
					Metadata: map[string]interface{}{
						"shift_id": shiftID.Hex(),
						"error":    err.Error(),
					},
				})
			}
			return nil, fmt.Errorf("failed to save order items with costs: %w", err)
		}
	}

	// Round total accounting cost
	result.TotalAccountingCost = math.Round(result.TotalAccountingCost*100) / 100
	
	// Record success metric
	metricStatus := MetricStatusSuccess
	if result.ItemsWithIncompleteCost > 0 {
		metricStatus = MetricStatusWarning
	}
	
	if s.monitoringService != nil {
		s.monitoringService.RecordMetric(Metric{
			Type:      MetricTypeShiftClosure,
			Status:    metricStatus,
			Timestamp: time.Now(),
			Duration:  time.Since(startTime),
			Message:   "Shift closure cost calculation completed",
			Metadata: map[string]interface{}{
				"shift_id":                shiftID.Hex(),
				"total_orders":            result.TotalOrders,
				"total_items":             result.TotalItems,
				"items_with_final_cost":   result.ItemsWithFinalCost,
				"items_with_incomplete_cost": result.ItemsWithIncompleteCost,
				"total_accounting_cost":   result.TotalAccountingCost,
			},
		})
	}

	return result, nil
}

// QueueCostRecalculation queues background jobs to recalculate costs for all menu items
// that use a specific ingredient. This is called when an ingredient's cost_per_unit changes.
// Requirements: 1.3, 9.1
func (s *CostCalculatorService) QueueCostRecalculation(ctx context.Context, ingredientID primitive.ObjectID) error {
	// Fetch the ingredient to get its name
	ingredient, err := s.ingredientRepo.FindByID(ctx, ingredientID)
	if err != nil {
		return fmt.Errorf("failed to fetch ingredient: %w", err)
	}

	// Find all menu items that use this ingredient
	menuItems, err := s.menuRepo.FindByIngredientName(ctx, ingredient.Name)
	if err != nil {
		return fmt.Errorf("failed to find menu items using ingredient: %w", err)
	}

	// If no recalculation service is configured, skip queuing
	if s.recalcService == nil {
		return nil
	}

	// Queue a background job for each menu item using the recalculation service
	for _, menuItem := range menuItems {
		err := s.recalcService.QueueRecalculation(menuItem.ID)
		if err != nil {
			// Log warning but continue with other items
			// In production, this should use proper logging
			fmt.Printf("Warning: failed to queue recalculation for menu item %s: %v\n", menuItem.ID.Hex(), err)
		}
	}

	return nil
}

// IngredientCostDetail represents cost detail for a single ingredient in a menu item
type IngredientCostDetail struct {
	Name              string
	Quantity          float64
	Unit              string
	CostPerUnit       float64
	ConversionRate    float64
	WastagePercentage float64
	TotalCost         float64
}

// MenuItemCostDetail represents detailed cost breakdown for a menu item
type MenuItemCostDetail struct {
	MenuItemID    primitive.ObjectID
	MenuItemName  string
	Ingredients   []IngredientCostDetail
	TotalCost     float64
	CostStatus    menu.CostStatus
	MissingCosts  []string // Names of ingredients with missing cost_per_unit
}

// CalculateMenuItemCostDetail calculates cost detail with ingredient breakdown for a menu item
// Requirements: 8.1, 8.2, 8.3
func (s *CostCalculatorService) CalculateMenuItemCostDetail(ctx context.Context, menuItemID primitive.ObjectID) (*MenuItemCostDetail, error) {
	// Fetch the menu item
	menuItem, err := s.menuRepo.FindByID(ctx, menuItemID)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch menu item: %w", err)
	}

	// Initialize result
	result := &MenuItemCostDetail{
		MenuItemID:   menuItemID,
		MenuItemName: menuItem.Name,
		Ingredients:  []IngredientCostDetail{},
		TotalCost:    0.0,
		CostStatus:   menu.CostStatusFinal,
		MissingCosts: []string{},
	}

	// If menu item has no ingredients, return empty breakdown
	if len(menuItem.Ingredients) == 0 {
		return result, nil
	}

	// Fetch all ingredients to build a lookup map
	allIngredients, err := s.ingredientRepo.FindAll(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch ingredients: %w", err)
	}

	// Build ingredient lookup map by name
	ingredientMap := make(map[string]*ingredient.Ingredient)
	for _, ing := range allIngredients {
		ingredientMap[ing.Name] = ing
	}

	// Calculate cost for each ingredient
	var totalCost float64
	hasIncompleteCost := false

	for _, menuIngredient := range menuItem.Ingredients {
		// Find the ingredient in our map
		ing, exists := ingredientMap[menuIngredient.Name]
		
		detail := IngredientCostDetail{
			Name:     menuIngredient.Name,
			Quantity: menuIngredient.Quantity,
			Unit:     string(menuIngredient.Unit),
		}

		if !exists {
			// Ingredient not found in database
			result.MissingCosts = append(result.MissingCosts, menuIngredient.Name)
			hasIncompleteCost = true
			detail.CostPerUnit = 0
			detail.ConversionRate = 1.0
			detail.WastagePercentage = 0.0
			detail.TotalCost = 0
			result.Ingredients = append(result.Ingredients, detail)
			continue
		}

		// Check if cost_per_unit is missing (0 or not set)
		if ing.CostPerUnit <= 0 {
			result.MissingCosts = append(result.MissingCosts, menuIngredient.Name)
			hasIncompleteCost = true
			detail.CostPerUnit = 0
			// Calculate conversion rate even if cost is missing (for display purposes)
			detail.ConversionRate = ingredient.GetConversionRate(ing.Unit, menuIngredient.Unit)
			detail.WastagePercentage = ing.WastagePercentage
			if detail.WastagePercentage < 0 {
				detail.WastagePercentage = 0.0
			}
			detail.TotalCost = 0
			result.Ingredients = append(result.Ingredients, detail)
			continue
		}

		// Calculate conversion rate dynamically based on stock unit and recipe unit
		conversionRate := ingredient.GetConversionRate(ing.Unit, menuIngredient.Unit)

		// Get wastage percentage (default 0.0 if not set)
		wastagePercentage := ing.WastagePercentage
		if wastagePercentage < 0 {
			wastagePercentage = 0.0
		}

		// Calculate cost for this ingredient
		// Formula: quantity * cost_per_unit * conversion_rate * (1 + wastage_percentage/100)
		ingredientCost := menuIngredient.Quantity * ing.CostPerUnit * conversionRate * (1 + wastagePercentage/100)
		ingredientCost = math.Round(ingredientCost*100) / 100

		detail.CostPerUnit = ing.CostPerUnit
		detail.ConversionRate = conversionRate
		detail.WastagePercentage = wastagePercentage
		detail.TotalCost = ingredientCost

		totalCost += ingredientCost
		result.Ingredients = append(result.Ingredients, detail)
	}

	// Round total cost to 2 decimal places
	result.TotalCost = math.Round(totalCost*100) / 100

	// Determine cost status
	if hasIncompleteCost {
		result.CostStatus = menu.CostStatusIncomplete
	}

	return result, nil
}
