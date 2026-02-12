package services

import (
	"context"
	"testing"
	"time"

	"cafe-pos/backend/domain/expense"
	"cafe-pos/backend/domain/ingredient"
	"cafe-pos/backend/domain/menu"
	"cafe-pos/backend/domain/order"
	"github.com/stretchr/testify/assert"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

// Test incomplete ingredient data handling
// Requirements: 1.5, 1.6
func TestCalculateMenuItemCost_IncompleteIngredientData(t *testing.T) {
	tests := []struct {
		name           string
		menuItem       *menu.MenuItem
		ingredients    []*ingredient.Ingredient
		expectedStatus menu.CostStatus
		expectedCost   float64
		description    string
	}{
		{
			name: "Missing cost_per_unit",
			menuItem: &menu.MenuItem{
				ID:   primitive.NewObjectID(),
				Name: "Latte",
				Ingredients: []menu.Ingredient{
					{Name: "Espresso", Quantity: 30, Unit: ingredient.UnitMilliliter},
					{Name: "Milk", Quantity: 150, Unit: ingredient.UnitMilliliter},
				},
			},
			ingredients: []*ingredient.Ingredient{
				{
					ID:          primitive.NewObjectID(),
					Name:        "Espresso",
					CostPerUnit: 0, // Missing cost
				},
				{
					ID:          primitive.NewObjectID(),
					Name:        "Milk",
					CostPerUnit: 50,
				},
			},
			expectedStatus: menu.CostStatusIncomplete,
			expectedCost:   7500, // Only milk cost: 150 * 50
			description:    "Should mark as INCOMPLETE when ingredient has zero cost_per_unit",
		},
		{
			name: "Ingredient not found in database",
			menuItem: &menu.MenuItem{
				ID:   primitive.NewObjectID(),
				Name: "Special Coffee",
				Ingredients: []menu.Ingredient{
					{Name: "RareBean", Quantity: 20, Unit: ingredient.UnitGram},
				},
			},
			ingredients:    []*ingredient.Ingredient{}, // Empty - ingredient not found
			expectedStatus: menu.CostStatusIncomplete,
			expectedCost:   0,
			description:    "Should mark as INCOMPLETE when ingredient not found in database",
		},
		{
			name: "All ingredients missing cost",
			menuItem: &menu.MenuItem{
				ID:   primitive.NewObjectID(),
				Name: "Mystery Drink",
				Ingredients: []menu.Ingredient{
					{Name: "Ingredient1", Quantity: 10, Unit: ingredient.UnitMilliliter},
					{Name: "Ingredient2", Quantity: 20, Unit: ingredient.UnitMilliliter},
				},
			},
			ingredients: []*ingredient.Ingredient{
				{
					ID:          primitive.NewObjectID(),
					Name:        "Ingredient1",
					CostPerUnit: 0,
				},
				{
					ID:          primitive.NewObjectID(),
					Name:        "Ingredient2",
					CostPerUnit: 0,
				},
			},
			expectedStatus: menu.CostStatusIncomplete,
			expectedCost:   0,
			description:    "Should mark as INCOMPLETE and return zero cost when all ingredients missing cost",
		},
		{
			name: "Partial incomplete data",
			menuItem: &menu.MenuItem{
				ID:   primitive.NewObjectID(),
				Name: "Mixed Drink",
				Ingredients: []menu.Ingredient{
					{Name: "ValidIngredient", Quantity: 10, Unit: ingredient.UnitMilliliter},
					{Name: "InvalidIngredient", Quantity: 20, Unit: ingredient.UnitMilliliter},
				},
			},
			ingredients: []*ingredient.Ingredient{
				{
					ID:          primitive.NewObjectID(),
					Name:        "ValidIngredient",
					CostPerUnit: 100,
				},
				{
					ID:          primitive.NewObjectID(),
					Name:        "InvalidIngredient",
					CostPerUnit: 0,
				},
			},
			expectedStatus: menu.CostStatusIncomplete,
			expectedCost:   1000, // Only valid ingredient: 10 * 100
			description:    "Should calculate partial cost but mark as INCOMPLETE",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Setup
			menuRepo := &mockMenuRepository{
				menuItems: map[primitive.ObjectID]*menu.MenuItem{
					tt.menuItem.ID: tt.menuItem,
				},
			}
			ingredientRepo := &mockIngredientRepository{
				ingredients: make(map[primitive.ObjectID]*ingredient.Ingredient),
			}
			for _, ing := range tt.ingredients {
				ingredientRepo.ingredients[ing.ID] = ing
			}

			service := NewCostCalculatorService(menuRepo, ingredientRepo, nil, nil)

			// Execute
			result, err := service.CalculateMenuItemCost(context.Background(), tt.menuItem.ID)

			// Verify
			assert.NoError(t, err, tt.description)
			assert.Equal(t, tt.expectedStatus, result.CostStatus, tt.description)
			assert.Equal(t, tt.expectedCost, result.CurrentCost, tt.description)
			if tt.expectedStatus == menu.CostStatusIncomplete {
				assert.NotEmpty(t, result.MissingIngredients, "Should list missing ingredients")
			}
		})
	}
}

// Test zero price edge case
// Requirements: 2.9
func TestEdgeCase_ZeroPrice(t *testing.T) {
	tests := []struct {
		name        string
		price       float64
		cost        float64
		description string
	}{
		{
			name:        "Zero price (promotional item)",
			price:       0,
			cost:        10000,
			description: "Should handle zero price gracefully",
		},
		{
			name:        "Negative price (gifted item)",
			price:       -1000,
			cost:        5000,
			description: "Should handle negative price gracefully",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Setup
			menuItem := &menu.MenuItem{
				ID:          primitive.NewObjectID(),
				Name:        "Test Item",
				Price:       tt.price,
				CurrentCost: tt.cost,
				CostStatus:  menu.CostStatusFinal,
			}

			menuRepo := &mockMenuRepository{
				menuItems: map[primitive.ObjectID]*menu.MenuItem{
					menuItem.ID: menuItem,
				},
			}

			service := NewProfitAnalyzerService(menuRepo, nil, nil, nil)

			// Execute
			profit, err := service.CalculateMenuItemProfit(context.Background(), menuItem.ID)

			// Verify
			assert.NoError(t, err, tt.description)
			// For zero or negative price, profit margin should be marked as N/A or handled specially
			// The service should not crash or return invalid values
			assert.NotNil(t, profit, tt.description)
		})
	}
}

// Test negative profit scenarios
// Requirements: 2.9
func TestEdgeCase_NegativeProfit(t *testing.T) {
	tests := []struct {
		name               string
		price              float64
		cost               float64
		expectedMargin     float64
		expectedAbsolute   float64
		expectedWarning    WarningStatus
		description        string
	}{
		{
			name:             "Cost exceeds price (loss)",
			price:            30000,
			cost:             50000,
			expectedMargin:   -66.67, // ((30000 - 50000) / 30000) * 100
			expectedAbsolute: -20000,
			expectedWarning:  WarningLoss,
			description:      "Should calculate negative margin when cost > price",
		},
		{
			name:             "Cost equals price (break-even)",
			price:            40000,
			cost:             40000,
			expectedMargin:   0,
			expectedAbsolute: 0,
			expectedWarning:  WarningNone,
			description:      "Should return zero margin when cost = price",
		},
		{
			name:             "Very small profit",
			price:            40000,
			cost:             39900,
			expectedMargin:   0.25, // ((40000 - 39900) / 40000) * 100
			expectedAbsolute: 100,
			expectedWarning:  WarningNone, // No threshold passed to CalculateMenuItemProfit
			description:      "Should handle very small profit margins",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Setup
			menuItem := &menu.MenuItem{
				ID:          primitive.NewObjectID(),
				Name:        "Test Item",
				Price:       tt.price,
				CurrentCost: tt.cost,
				CostStatus:  menu.CostStatusFinal,
			}

			menuRepo := &mockMenuRepository{
				menuItems: map[primitive.ObjectID]*menu.MenuItem{
					menuItem.ID: menuItem,
				},
			}

			service := NewProfitAnalyzerService(menuRepo, nil, nil, nil)

			// Execute
			profit, err := service.CalculateMenuItemProfit(context.Background(), menuItem.ID)

			// Verify
			assert.NoError(t, err, tt.description)
			assert.InDelta(t, tt.expectedMargin, profit.ProfitMargin, 0.01, tt.description)
			assert.Equal(t, tt.expectedAbsolute, profit.AbsoluteProfit, tt.description)

			// The warning status is already set in the profit object
			assert.Equal(t, tt.expectedWarning, profit.WarningStatus, tt.description)
		})
	}
}

// Test empty date ranges
// Requirements: 2.9
func TestEdgeCase_EmptyDateRange(t *testing.T) {
	tests := []struct {
		name        string
		dateRange   DateRange
		orderItems  []*order.OrderItemWithCost
		description string
	}{
		{
			name: "No orders in date range",
			dateRange: DateRange{
				Start: time.Date(2024, 1, 1, 0, 0, 0, 0, time.UTC),
				End:   time.Date(2024, 1, 31, 23, 59, 59, 0, time.UTC),
			},
			orderItems:  []*order.OrderItemWithCost{}, // Empty
			description: "Should return empty results for date range with no orders",
		},
		{
			name: "Orders outside date range",
			dateRange: DateRange{
				Start: time.Date(2024, 2, 1, 0, 0, 0, 0, time.UTC),
				End:   time.Date(2024, 2, 28, 23, 59, 59, 0, time.UTC),
			},
			orderItems: []*order.OrderItemWithCost{
				{
					MenuItemID:       primitive.NewObjectID(),
					Name:             "Coffee",
					Price:            45000,
					Quantity:         1,
					AccountingCost:   15000,
					CostStatus:       order.CostStatusFinal,
					CostCalculatedAt: time.Date(2024, 1, 15, 10, 0, 0, 0, time.UTC), // January
				},
			},
			description: "Should return empty results when orders are outside date range",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Setup
			menuRepo := &mockMenuRepository{
				menuItems: make(map[primitive.ObjectID]*menu.MenuItem),
			}
			orderItemRepo := &mockOrderItemRepository{
				orderItems: tt.orderItems,
			}

			service := NewProfitAnalyzerService(menuRepo, orderItemRepo, nil, nil)

			// Execute
			profits, err := service.GetCategoryProfits(context.Background(), tt.dateRange)

			// Verify
			assert.NoError(t, err, tt.description)
			// Should return empty array or array with zero values
			if len(profits) > 0 {
				for _, profit := range profits {
					assert.Equal(t, 0.0, profit.TotalRevenue, "Revenue should be zero for empty date range")
					assert.Equal(t, 0.0, profit.TotalCost, "Cost should be zero for empty date range")
					assert.Equal(t, 0, profit.OrderCount, "Order count should be zero for empty date range")
				}
			}
		})
	}
}

// Test shift closure with incomplete ingredient data
// Requirements: 1.5, 1.6
func TestCalculateShiftOrderCosts_IncompleteData(t *testing.T) {
	// Setup
	shiftID := primitive.NewObjectID()
	menuItemID := primitive.NewObjectID()

	menuRepo := &mockMenuRepository{
		menuItems: map[primitive.ObjectID]*menu.MenuItem{
			menuItemID: {
				ID:   menuItemID,
				Name: "Incomplete Item",
				Ingredients: []menu.Ingredient{
					{Name: "MissingIngredient", Quantity: 10, Unit: ingredient.UnitMilliliter},
				},
			},
		},
	}

	ingredientRepo := &mockIngredientRepository{
		ingredients: make(map[primitive.ObjectID]*ingredient.Ingredient),
		// No ingredients - simulating missing data
	}

	orderRepo := &mockOrderRepository{
		orders: map[primitive.ObjectID]*order.Order{
			primitive.NewObjectID(): {
				ID:      primitive.NewObjectID(),
				ShiftID: shiftID,
				Items: []order.OrderItem{
					{
						MenuItemID: menuItemID,
						Name:       "Incomplete Item",
						Price:      50000,
						Quantity:   2,
						Subtotal:   100000,
					},
				},
			},
		},
	}

	orderItemRepo := &mockOrderItemRepository{
		orderItems: make([]*order.OrderItemWithCost, 0),
	}

	service := NewCostCalculatorService(menuRepo, ingredientRepo, orderRepo, orderItemRepo)

	// Execute
	result, err := service.CalculateShiftOrderCosts(context.Background(), shiftID)

	// Verify
	assert.NoError(t, err, "Should not error even with incomplete data")
	assert.Equal(t, 1, result.TotalOrders, "Should process the order")
	assert.Equal(t, 1, result.TotalItems, "Should process the order item")
	assert.Equal(t, 1, result.ItemsWithIncompleteCost, "Should mark item as incomplete")
	assert.Equal(t, 0, result.ItemsWithFinalCost, "Should have no items with final cost")

	// Verify order item was created with INCOMPLETE status
	assert.Len(t, orderItemRepo.orderItems, 1, "Should create order item")
	assert.Equal(t, order.CostStatusIncomplete, orderItemRepo.orderItems[0].CostStatus, "Should mark as INCOMPLETE")
}

// Test date range validation
func TestOperatingExpenseService_InvalidDateRange(t *testing.T) {
	tests := []struct {
		name        string
		startDate   string
		endDate     string
		shouldError bool
		description string
	}{
		{
			name:        "Start date after end date",
			startDate:   "2024-02-01",
			endDate:     "2024-01-01",
			shouldError: true,
			description: "Should error when start date is after end date",
		},
		{
			name:        "Invalid date format",
			startDate:   "2024-13-01", // Invalid month
			endDate:     "2024-12-31",
			shouldError: true,
			description: "Should error with invalid date format",
		},
		{
			name:        "Valid date range",
			startDate:   "2024-01-01",
			endDate:     "2024-01-31",
			shouldError: false,
			description: "Should succeed with valid date range",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Setup
			expenseRepo := &mockOperatingExpenseRepository{
				expenses: make(map[primitive.ObjectID]*expense.OperatingExpense),
			}
			service := NewOperatingExpenseService(expenseRepo)

			req := &expense.OperatingExpenseRequest{
				PeriodStart:    tt.startDate,
				PeriodEnd:      tt.endDate,
				StaffSalary:    1000000,
				Rent:           500000,
				Utilities:      200000,
				MarketingCosts: 100000,
				OtherExpenses:  100000,
			}

			// Execute
			_, err := service.UpsertOperatingExpense(context.Background(), req)

			// Verify
			if tt.shouldError {
				assert.Error(t, err, tt.description)
			} else {
				assert.NoError(t, err, tt.description)
			}
		})
	}
}

// Test negative expense amounts validation
func TestOperatingExpenseService_NegativeAmounts(t *testing.T) {
	// Setup
	expenseRepo := &mockOperatingExpenseRepository{
		expenses: make(map[primitive.ObjectID]*expense.OperatingExpense),
	}
	service := NewOperatingExpenseService(expenseRepo)

	req := &expense.OperatingExpenseRequest{
		PeriodStart:    "2024-01-01",
		PeriodEnd:      "2024-01-31",
		StaffSalary:    -1000000, // Negative amount
		Rent:           500000,
		Utilities:      200000,
		MarketingCosts: 100000,
		OtherExpenses:  100000,
	}

	// Execute
	_, err := service.UpsertOperatingExpense(context.Background(), req)

	// Verify
	assert.Error(t, err, "Should error with negative expense amount")
}
