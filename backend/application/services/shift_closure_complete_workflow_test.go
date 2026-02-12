package services

import (
	"context"
	"testing"
	"time"

	"cafe-pos/backend/domain/ingredient"
	"cafe-pos/backend/domain/menu"
	"cafe-pos/backend/domain/order"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

// TestCompleteShiftClosureWorkflow tests the complete shift closure workflow
// Requirements: 5.1, 5.2, 5.8, 9.6
// - Create shift with orders
// - Close shift
// - Verify accounting_cost is calculated
// - Update ingredient costs
// - Verify accounting_cost remains unchanged
func TestCompleteShiftClosureWorkflow(t *testing.T) {
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

	// Initialize services
	costCalculator := NewCostCalculatorService(menuRepo, ingredientRepo, orderRepo, orderItemRepo)
	
	// Step 1: Create test ingredients with initial costs
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
		CostPerUnit:       10, // 10 VND per gram
		ConversionRate:    1.0,
		WastagePercentage: 0.0,
		CreatedAt:         time.Now(),
		UpdatedAt:         time.Now(),
	}
	ingredientRepo.ingredients[sugarIngredient.ID] = sugarIngredient

	// Step 2: Create test menu items
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
				Unit:     espressoIngredient.Unit,
			},
			{
				Name:     milkIngredient.Name,
				Quantity: 150, // 150 ml
				Unit:     milkIngredient.Unit,
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
		Price:       50000,
		Category:    "Coffee",
		Description: "Smooth latte",
		Ingredients: []menu.Ingredient{
			{
				Name:     espressoIngredient.Name,
				Quantity: 30, // 30 grams
				Unit:     espressoIngredient.Unit,
			},
			{
				Name:     milkIngredient.Name,
				Quantity: 200, // 200 ml
				Unit:     milkIngredient.Unit,
			},
			{
				Name:     sugarIngredient.Name,
				Quantity: 10, // 10 grams
				Unit:     sugarIngredient.Unit,
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
		},
		Subtotal:  cappuccinoItem.Price * 2,
		Total:     cappuccinoItem.Price * 2,
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
				MenuItemID: latteItem.ID,
				Name:       latteItem.Name,
				Price:      latteItem.Price,
				Quantity:   1,
				Subtotal:   latteItem.Price,
			},
			{
				MenuItemID: cappuccinoItem.ID,
				Name:       cappuccinoItem.Name,
				Price:      cappuccinoItem.Price,
				Quantity:   1,
				Subtotal:   cappuccinoItem.Price,
			},
		},
		Subtotal:  latteItem.Price + cappuccinoItem.Price,
		Total:     latteItem.Price + cappuccinoItem.Price,
		Status:    order.StatusServed,
		CreatedAt: time.Now().Add(-1 * time.Hour),
		UpdatedAt: time.Now().Add(-1 * time.Hour),
	}
	orderRepo.orders[order2.ID] = order2

	// Step 5: Close shift and calculate costs
	t.Log("Closing shift and calculating costs...")
	summary, err := costCalculator.CalculateShiftOrderCosts(ctx, testShift.ID)
	require.NoError(t, err, "Should calculate shift order costs successfully")
	require.NotNil(t, summary, "Cost calculation summary should not be nil")

	// Verify summary statistics
	assert.Equal(t, 2, summary.TotalOrders, "Should have processed 2 orders")
	assert.Equal(t, 3, summary.TotalItems, "Should have processed 3 order items")
	assert.Equal(t, 3, summary.ItemsWithFinalCost, "All items should have FINAL cost status")
	assert.Equal(t, 0, summary.ItemsWithIncompleteCost, "No items should have incomplete cost")

	// Step 6: Verify accounting_cost is calculated for all order items
	t.Log("Verifying accounting costs were calculated...")
	
	// Get all order items
	allOrderItems := orderItemRepo.orderItems
	require.Len(t, allOrderItems, 3, "Should have 3 order items total")

	// Store original accounting costs for later verification
	originalCosts := make(map[primitive.ObjectID]float64)
	
	for _, item := range allOrderItems {
		t.Logf("Order item: %s, Cost: %.2f, Status: %s", item.Name, item.AccountingCost, item.CostStatus)
		
		// Verify accounting_cost is set
		assert.NotZero(t, item.AccountingCost, "Accounting cost should be calculated for %s", item.Name)
		
		// Verify cost_status is FINAL
		assert.Equal(t, order.CostStatusFinal, item.CostStatus, "Cost status should be FINAL for %s", item.Name)
		
		// Verify cost_calculated_at is set
		assert.False(t, item.CostCalculatedAt.IsZero(), "Cost calculated at should be set for %s", item.Name)
		
		// Store original cost
		originalCosts[item.ID] = item.AccountingCost
	}

	// Calculate expected costs manually
	// Cappuccino: 30*200*1.0*1.05 + 150*50*1.0*1.10 = 6300 + 8250 = 14550 per item
	// Latte: 30*200*1.0*1.05 + 200*50*1.0*1.10 + 10*10*1.0*1.0 = 6300 + 11000 + 100 = 17400 per item
	
	cappuccinoExpectedCostPerItem := 14550.0
	latteExpectedCostPerItem := 17400.0

	// Verify costs match expected values (with small tolerance for rounding)
	for _, item := range allOrderItems {
		if item.Name == "Cappuccino" {
			expectedTotal := cappuccinoExpectedCostPerItem * float64(item.Quantity)
			assert.InDelta(t, expectedTotal, item.AccountingCost, 1.0, 
				"Cappuccino cost should be approximately %.2f", expectedTotal)
		} else if item.Name == "Latte" {
			expectedTotal := latteExpectedCostPerItem * float64(item.Quantity)
			assert.InDelta(t, expectedTotal, item.AccountingCost, 1.0,
				"Latte cost should be approximately %.2f", expectedTotal)
		}
	}

	// Step 7: Update ingredient costs (simulate price changes)
	t.Log("Updating ingredient costs...")
	espressoIngredient.CostPerUnit = 300.0 // Increase from 200 to 300
	milkIngredient.CostPerUnit = 80.0      // Increase from 50 to 80
	ingredientRepo.ingredients[espressoIngredient.ID] = espressoIngredient
	ingredientRepo.ingredients[milkIngredient.ID] = milkIngredient

	// Step 8: Verify accounting_cost remains unchanged (immutability)
	t.Log("Verifying accounting costs remain unchanged after ingredient cost update...")
	
	// Get order items again
	allOrderItems2 := orderItemRepo.orderItems

	for _, item := range allOrderItems2 {
		originalCost, exists := originalCosts[item.ID]
		require.True(t, exists, "Should find original cost for item %s", item.Name)
		
		// Verify accounting cost has NOT changed
		assert.Equal(t, originalCost, item.AccountingCost,
			"Accounting cost for %s should remain unchanged at %.2f after ingredient cost update",
			item.Name, originalCost)
		
		// Verify cost status is still FINAL
		assert.Equal(t, order.CostStatusFinal, item.CostStatus,
			"Cost status for %s should remain FINAL", item.Name)
		
		t.Logf("✓ %s accounting cost unchanged: %.2f", item.Name, item.AccountingCost)
	}

	// Step 9: Verify menu item current_cost DOES change (for future orders)
	t.Log("Verifying menu item current_cost is updated for future orders...")
	
	// Recalculate current costs for menu items
	cappuccinoCost, err := costCalculator.CalculateMenuItemCost(ctx, cappuccinoItem.ID)
	require.NoError(t, err)
	
	latteCost, err := costCalculator.CalculateMenuItemCost(ctx, latteItem.ID)
	require.NoError(t, err)

	// New expected costs with updated ingredient prices:
	// Cappuccino: 30*300*1.0*1.05 + 150*80*1.0*1.10 = 9450 + 13200 = 22650
	// Latte: 30*300*1.0*1.05 + 200*80*1.0*1.10 + 10*10*1.0*1.0 = 9450 + 17600 + 100 = 27150

	assert.InDelta(t, 22650.0, cappuccinoCost.CurrentCost, 1.0,
		"Cappuccino current_cost should be updated to approximately 22650")
	t.Logf("✓ Cappuccino current_cost updated: %.2f", cappuccinoCost.CurrentCost)

	assert.InDelta(t, 27150.0, latteCost.CurrentCost, 1.0,
		"Latte current_cost should be updated to approximately 27150")
	t.Logf("✓ Latte current_cost updated: %.2f", latteCost.CurrentCost)

	t.Log("✅ Complete shift closure workflow test passed!")
	t.Log("   - Accounting costs calculated correctly at shift closure")
	t.Log("   - Accounting costs remain immutable after ingredient cost changes")
	t.Log("   - Menu item current_cost updated for future orders")
}
