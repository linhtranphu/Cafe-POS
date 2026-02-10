package services

import (
	"context"
	"testing"
	"time"

	"cafe-pos/backend/domain/ingredient"
	"cafe-pos/backend/domain/menu"
	"cafe-pos/backend/domain/order"
	"github.com/stretchr/testify/assert"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

// TestShiftClosureWorkflow_Integration tests the complete shift closure workflow
// Requirements: 5.1, 5.2, 5.3
//
// This integration test verifies:
// 1. Create shift with orders
// 2. Close shift
// 3. Verify all order items have accounting_cost
// 4. Verify cost_status = FINAL
func TestShiftClosureWorkflow_Integration(t *testing.T) {
	ctx := context.Background()

	// Create mock repositories
	menuRepo := &mockMenuRepository{
		menuItems: make(map[primitive.ObjectID]*menu.MenuItem),
	}
	ingredientRepo := &mockIngredientRepository{
		ingredients: make(map[primitive.ObjectID]*ingredient.Ingredient),
	}
	orderRepo := &mockOrderRepository{
		orders: make(map[primitive.ObjectID]*order.Order),
	}
	orderItemRepo := &mockOrderItemRepository{
		orderItems: make([]*order.OrderItemWithCost, 0),
	}

	// Create service
	costCalculator := NewCostCalculatorService(menuRepo, ingredientRepo, orderRepo, orderItemRepo)

	// Step 1: Setup test data - Create ingredients
	espressoIngredient := &ingredient.Ingredient{
		ID:                primitive.NewObjectID(),
		Name:              "Espresso",
		Category:          "Beverage",
		Unit:              ingredient.UnitGram,
		Quantity:          5000,
		MinStock:          500,
		CostPerUnit:       200, // 200 VND per gram
		ConversionRate:    1.0,
		WastagePercentage: 5.0, // 5% wastage
		CreatedAt:         time.Now(),
		UpdatedAt:         time.Now(),
	}
	ingredientRepo.ingredients[espressoIngredient.ID] = espressoIngredient

	milkIngredient := &ingredient.Ingredient{
		ID:                primitive.NewObjectID(),
		Name:              "Milk",
		Category:          "Beverage",
		Unit:              ingredient.UnitMilliliter,
		Quantity:          10000,
		MinStock:          1000,
		CostPerUnit:       50, // 50 VND per ml
		ConversionRate:    1.0,
		WastagePercentage: 10.0, // 10% wastage
		CreatedAt:         time.Now(),
		UpdatedAt:         time.Now(),
	}
	ingredientRepo.ingredients[milkIngredient.ID] = milkIngredient

	sugarIngredient := &ingredient.Ingredient{
		ID:                primitive.NewObjectID(),
		Name:              "Sugar",
		Category:          "Sweetener",
		Unit:              ingredient.UnitGram,
		Quantity:          2000,
		MinStock:          200,
		CostPerUnit:       30, // 30 VND per gram
		ConversionRate:    1.0,
		WastagePercentage: 0.0,
		CreatedAt:         time.Now(),
		UpdatedAt:         time.Now(),
	}
	ingredientRepo.ingredients[sugarIngredient.ID] = sugarIngredient

	// Step 2: Create menu items
	cappuccinoItem := &menu.MenuItem{
		ID:          primitive.NewObjectID(),
		Name:        "Cappuccino",
		Price:       45000,
		Category:    "Coffee",
		Description: "Classic cappuccino",
		Ingredients: []menu.Ingredient{
			{
				Name:     espressoIngredient.Name,
				Quantity: 30, // 30 grams
				Unit:     string(espressoIngredient.Unit),
			},
			{
				Name:     milkIngredient.Name,
				Quantity: 150, // 150 ml
				Unit:     string(milkIngredient.Unit),
			},
		},
		Available: true,
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}
	menuRepo.menuItems[cappuccinoItem.ID] = cappuccinoItem

	latteItem := &menu.MenuItem{
		ID:          primitive.NewObjectID(),
		Name:        "Latte",
		Price:       42000,
		Category:    "Coffee",
		Description: "Smooth latte",
		Ingredients: []menu.Ingredient{
			{
				Name:     espressoIngredient.Name,
				Quantity: 30, // 30 grams
				Unit:     string(espressoIngredient.Unit),
			},
			{
				Name:     milkIngredient.Name,
				Quantity: 200, // 200 ml
				Unit:     string(milkIngredient.Unit),
			},
			{
				Name:     sugarIngredient.Name,
				Quantity: 10, // 10 grams
				Unit:     string(sugarIngredient.Unit),
			},
		},
		Available: true,
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}
	menuRepo.menuItems[latteItem.ID] = latteItem

	// Step 3: Create shift
	testShift := &order.Shift{
		ID:        primitive.NewObjectID(),
		UserID:    primitive.NewObjectID(),
		UserName:  "waiter1",
		RoleType:  order.RoleWaiter,
		Status:    order.ShiftOpen,
		StartedAt: time.Now().Add(-3 * time.Hour),
		CreatedAt: time.Now().Add(-3 * time.Hour),
		UpdatedAt: time.Now().Add(-3 * time.Hour),
	}

	// Step 4: Create orders in the shift
	order1 := &order.Order{
		ID:       primitive.NewObjectID(),
		ShiftID:  testShift.ID,
		WaiterID: testShift.UserID,
		Items: []order.OrderItem{
			{
				MenuItemID: cappuccinoItem.ID,
				Name:       cappuccinoItem.Name,
				Price:      cappuccinoItem.Price,
				Quantity:   2,
				Subtotal:   cappuccinoItem.Price * 2,
			},
			{
				MenuItemID: latteItem.ID,
				Name:       latteItem.Name,
				Price:      latteItem.Price,
				Quantity:   1,
				Subtotal:   latteItem.Price,
			},
		},
		Subtotal:  cappuccinoItem.Price*2 + latteItem.Price,
		Total:     cappuccinoItem.Price*2 + latteItem.Price,
		Status:    order.StatusServed,
		CreatedAt: time.Now().Add(-2 * time.Hour),
		UpdatedAt: time.Now().Add(-2 * time.Hour),
	}
	orderRepo.orders[order1.ID] = order1

	order2 := &order.Order{
		ID:       primitive.NewObjectID(),
		ShiftID:  testShift.ID,
		WaiterID: testShift.UserID,
		Items: []order.OrderItem{
			{
				MenuItemID: cappuccinoItem.ID,
				Name:       cappuccinoItem.Name,
				Price:      cappuccinoItem.Price,
				Quantity:   3,
				Subtotal:   cappuccinoItem.Price * 3,
			},
		},
		Subtotal:  cappuccinoItem.Price * 3,
		Total:     cappuccinoItem.Price * 3,
		Status:    order.StatusServed,
		CreatedAt: time.Now().Add(-1 * time.Hour),
		UpdatedAt: time.Now().Add(-1 * time.Hour),
	}
	orderRepo.orders[order2.ID] = order2

	// Step 5: Close shift and calculate costs
	summary, err := costCalculator.CalculateShiftOrderCosts(ctx, testShift.ID)
	assert.NoError(t, err, "Failed to calculate shift order costs")
	assert.NotNil(t, summary, "Cost calculation summary should not be nil")

	// Verify summary statistics
	assert.Equal(t, 2, summary.TotalOrders, "Should have processed 2 orders")
	assert.Equal(t, 3, summary.TotalItems, "Should have processed 3 order items (2 from order1, 1 from order2)")
	assert.Equal(t, 3, summary.ItemsWithFinalCost, "All items should have FINAL cost status")
	assert.Equal(t, 0, summary.ItemsWithIncompleteCost, "No items should have incomplete cost")
	assert.Greater(t, summary.TotalAccountingCost, 0.0, "Total accounting cost should be positive")

	// Step 6: Verify all order items have accounting_cost
	order1Items, err := orderItemRepo.FindByOrderID(ctx, order1.ID)
	assert.NoError(t, err, "Failed to get order1 items")
	assert.Len(t, order1Items, 2, "Order1 should have 2 order items")

	// Verify order1 item 1 (Cappuccino x2)
	cappuccinoOrderItem := order1Items[0]
	assert.Equal(t, cappuccinoItem.ID, cappuccinoOrderItem.MenuItemID)
	assert.Equal(t, order.CostStatusFinal, cappuccinoOrderItem.CostStatus, "Cost status should be FINAL")
	assert.Greater(t, cappuccinoOrderItem.AccountingCost, 0.0, "Accounting cost should be positive")
	assert.NotZero(t, cappuccinoOrderItem.CostCalculatedAt, "Cost calculated timestamp should be set")

	// Calculate expected cost for Cappuccino
	// Espresso: 30g * 200 VND/g * 1.05 (5% wastage) = 6300 VND
	// Milk: 150ml * 50 VND/ml * 1.10 (10% wastage) = 8250 VND
	// Total per item: 14550 VND
	// For quantity 2: 29100 VND
	expectedCappuccinoCost := 29100.0
	assert.InDelta(t, expectedCappuccinoCost, cappuccinoOrderItem.AccountingCost, 0.01,
		"Cappuccino accounting cost should match expected value")

	// Verify order1 item 2 (Latte x1)
	latteOrderItem := order1Items[1]
	assert.Equal(t, latteItem.ID, latteOrderItem.MenuItemID)
	assert.Equal(t, order.CostStatusFinal, latteOrderItem.CostStatus, "Cost status should be FINAL")
	assert.Greater(t, latteOrderItem.AccountingCost, 0.0, "Accounting cost should be positive")

	// Calculate expected cost for Latte
	// Espresso: 30g * 200 VND/g * 1.05 = 6300 VND
	// Milk: 200ml * 50 VND/ml * 1.10 = 11000 VND
	// Sugar: 10g * 30 VND/g * 1.0 = 300 VND
	// Total: 17600 VND
	expectedLatteCost := 17600.0
	assert.InDelta(t, expectedLatteCost, latteOrderItem.AccountingCost, 0.01,
		"Latte accounting cost should match expected value")

	// Verify order2 items (Cappuccino x3)
	order2Items, err := orderItemRepo.FindByOrderID(ctx, order2.ID)
	assert.NoError(t, err, "Failed to get order2 items")
	assert.Len(t, order2Items, 1, "Order2 should have 1 order item")

	cappuccinoOrderItem2 := order2Items[0]
	assert.Equal(t, cappuccinoItem.ID, cappuccinoOrderItem2.MenuItemID)
	assert.Equal(t, order.CostStatusFinal, cappuccinoOrderItem2.CostStatus, "Cost status should be FINAL")

	// For quantity 3: 14550 * 3 = 43650 VND
	expectedCappuccinoCost2 := 43650.0
	assert.InDelta(t, expectedCappuccinoCost2, cappuccinoOrderItem2.AccountingCost, 0.01,
		"Cappuccino accounting cost for order2 should match expected value")

	// Step 7: Verify total accounting cost
	expectedTotalCost := expectedCappuccinoCost + expectedLatteCost + expectedCappuccinoCost2
	assert.InDelta(t, expectedTotalCost, summary.TotalAccountingCost, 0.01,
		"Total accounting cost should match sum of all order items")
}

// TestShiftClosureWorkflow_WithIncompleteData tests shift closure with missing ingredient costs
// Requirements: 5.1, 5.2, 5.3
func TestShiftClosureWorkflow_WithIncompleteData(t *testing.T) {
	ctx := context.Background()

	// Create mock repositories
	menuRepo := &mockMenuRepository{
		menuItems: make(map[primitive.ObjectID]*menu.MenuItem),
	}
	ingredientRepo := &mockIngredientRepository{
		ingredients: make(map[primitive.ObjectID]*ingredient.Ingredient),
	}
	orderRepo := &mockOrderRepository{
		orders: make(map[primitive.ObjectID]*order.Order),
	}
	orderItemRepo := &mockOrderItemRepository{
		orderItems: make([]*order.OrderItemWithCost, 0),
	}

	// Create service
	costCalculator := NewCostCalculatorService(menuRepo, ingredientRepo, orderRepo, orderItemRepo)

	// Create ingredient with valid cost
	espressoIngredient := &ingredient.Ingredient{
		ID:                primitive.NewObjectID(),
		Name:              "Espresso",
		CostPerUnit:       200,
		ConversionRate:    1.0,
		WastagePercentage: 0.0,
	}
	ingredientRepo.ingredients[espressoIngredient.ID] = espressoIngredient

	// Create ingredient with missing cost
	chocolateIngredient := &ingredient.Ingredient{
		ID:                primitive.NewObjectID(),
		Name:              "Chocolate",
		CostPerUnit:       0, // Missing cost
		ConversionRate:    1.0,
		WastagePercentage: 0.0,
	}
	ingredientRepo.ingredients[chocolateIngredient.ID] = chocolateIngredient

	// Create menu item with incomplete ingredient data
	mochaItem := &menu.MenuItem{
		ID:    primitive.NewObjectID(),
		Name:  "Mocha",
		Price: 50000,
		Ingredients: []menu.Ingredient{
			{Name: espressoIngredient.Name, Quantity: 30, Unit: "g"},
			{Name: chocolateIngredient.Name, Quantity: 20, Unit: "g"}, // Missing cost
		},
	}
	menuRepo.menuItems[mochaItem.ID] = mochaItem

	// Create shift and order
	testShift := &order.Shift{
		ID:       primitive.NewObjectID(),
		Status:   order.ShiftOpen,
		RoleType: order.RoleWaiter,
	}

	testOrder := &order.Order{
		ID:      primitive.NewObjectID(),
		ShiftID: testShift.ID,
		Items: []order.OrderItem{
			{
				MenuItemID: mochaItem.ID,
				Name:       mochaItem.Name,
				Price:      mochaItem.Price,
				Quantity:   1,
				Subtotal:   mochaItem.Price,
			},
		},
		Subtotal: mochaItem.Price,
		Total:    mochaItem.Price,
		Status:   order.StatusServed,
	}
	orderRepo.orders[testOrder.ID] = testOrder

	// Close shift and calculate costs
	summary, err := costCalculator.CalculateShiftOrderCosts(ctx, testShift.ID)
	assert.NoError(t, err)
	assert.NotNil(t, summary)

	// Verify summary shows incomplete items
	assert.Equal(t, 1, summary.TotalOrders)
	assert.Equal(t, 1, summary.TotalItems)
	assert.Equal(t, 0, summary.ItemsWithFinalCost, "Should have 0 items with FINAL cost")
	assert.Equal(t, 1, summary.ItemsWithIncompleteCost, "Should have 1 item with INCOMPLETE cost")

	// Verify order item has INCOMPLETE status
	orderItems, err := orderItemRepo.FindByOrderID(ctx, testOrder.ID)
	assert.NoError(t, err)
	assert.Len(t, orderItems, 1)

	assert.Equal(t, order.CostStatusIncomplete, orderItems[0].CostStatus,
		"Cost status should be INCOMPLETE when ingredient cost is missing")
	assert.Greater(t, orderItems[0].AccountingCost, 0.0,
		"Accounting cost should include cost from available ingredients")
}
