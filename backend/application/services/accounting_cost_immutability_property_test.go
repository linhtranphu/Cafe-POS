package services

import (
	"context"
	"testing"
	"time"

	"cafe-pos/backend/domain/ingredient"
	"cafe-pos/backend/domain/menu"
	"cafe-pos/backend/domain/order"
	"github.com/leanovate/gopter"
	"github.com/leanovate/gopter/gen"
	"github.com/leanovate/gopter/prop"
	"github.com/stretchr/testify/assert"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

// Feature: menu-cost-profit-analysis, Property 6: Accounting Cost Immutability
// Validates: Requirements 5.8, 9.6
//
// Property: For any order item with cost_status = "FINAL" (from a closed shift),
// when ingredient cost_per_unit values change, the accounting_cost should remain unchanged.
func TestProperty_AccountingCostImmutability(t *testing.T) {
	parameters := gopter.DefaultTestParameters()
	parameters.MinSuccessfulTests = 100 // Minimum 100 iterations as per spec
	properties := gopter.NewProperties(parameters)

	properties.Property("Accounting cost remains unchanged when ingredient cost changes", prop.ForAll(
		func(initialCost, newCost float64) bool {
			// Skip invalid cases
			if initialCost <= 0 || newCost <= 0 {
				return true
			}

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

			// Create test ingredient
			testIngredient := &ingredient.Ingredient{
				ID:                primitive.NewObjectID(),
				Name:              "Test Ingredient",
				Category:          "Test",
				Unit:              ingredient.UnitGram,
				Quantity:          1000,
				MinStock:          100,
				CostPerUnit:       initialCost,
				ConversionRate:    1.0,
				WastagePercentage: 0.0,
				CreatedAt:         time.Now(),
				UpdatedAt:         time.Now(),
			}
			ingredientRepo.ingredients[testIngredient.ID] = testIngredient

			// Create test menu item
			testMenuItem := &menu.MenuItem{
				ID:          primitive.NewObjectID(),
				Name:        "Test Menu Item",
				Price:       50000,
				Category:    "Test",
				Description: "Test item",
				Ingredients: []menu.Ingredient{
					{
						Name:     testIngredient.Name,
						Quantity: 100,
						Unit:     testIngredient.Unit,
					},
				},
				Available: true,
				CreatedAt: time.Now(),
				UpdatedAt: time.Now(),
			}
			menuRepo.menuItems[testMenuItem.ID] = testMenuItem

			// Create test shift (closed)
			testShift := &order.Shift{
				ID:        primitive.NewObjectID(),
				UserID:    primitive.NewObjectID(),
				UserName:  "test_user",
				RoleType:  order.RoleWaiter,
				Status:    order.ShiftClosed,
				StartedAt: time.Now().Add(-2 * time.Hour),
				CreatedAt: time.Now().Add(-2 * time.Hour),
				UpdatedAt: time.Now().Add(-1 * time.Hour),
			}
			endedAt := time.Now().Add(-1 * time.Hour)
			testShift.EndedAt = &endedAt

			// Create test order
			testOrder := &order.Order{
				ID:       primitive.NewObjectID(),
				ShiftID:  testShift.ID,
				WaiterID: testShift.UserID,
				Items: []order.OrderItem{
					{
						MenuItemID: testMenuItem.ID,
						Name:       testMenuItem.Name,
						Price:      testMenuItem.Price,
						Quantity:   1,
						Subtotal:   testMenuItem.Price,
					},
				},
				Subtotal:  testMenuItem.Price,
				Total:     testMenuItem.Price,
				Status:    order.StatusServed,
				CreatedAt: time.Now().Add(-90 * time.Minute),
				UpdatedAt: time.Now().Add(-90 * time.Minute),
			}
			orderRepo.orders[testOrder.ID] = testOrder

			// Calculate shift order costs (simulating shift closure)
			_, err := costCalculator.CalculateShiftOrderCosts(ctx, testShift.ID)
			if err != nil {
				t.Logf("Failed to calculate shift order costs: %v", err)
				return false
			}

			// Get the order item with accounting cost
			orderItems, err := orderItemRepo.FindByOrderID(ctx, testOrder.ID)
			if err != nil || len(orderItems) == 0 {
				t.Logf("Failed to get order items: %v", err)
				return false
			}

			originalAccountingCost := orderItems[0].AccountingCost
			originalCostStatus := orderItems[0].CostStatus

			// Verify cost status is FINAL
			if originalCostStatus != "FINAL" {
				t.Logf("Expected cost status FINAL, got %s", originalCostStatus)
				return false
			}

			// Verify accounting cost was calculated
			if originalAccountingCost <= 0 {
				t.Logf("Expected positive accounting cost, got %f", originalAccountingCost)
				return false
			}

			// Update ingredient cost
			testIngredient.CostPerUnit = newCost
			testIngredient.UpdatedAt = time.Now()
			ingredientRepo.ingredients[testIngredient.ID] = testIngredient

			// Recalculate menu item cost (this should happen in background)
			_, err = costCalculator.CalculateMenuItemCost(ctx, testMenuItem.ID)
			if err != nil {
				t.Logf("Failed to recalculate menu item cost: %v", err)
				return false
			}

			// Get the order item again
			updatedOrderItems, err := orderItemRepo.FindByOrderID(ctx, testOrder.ID)
			if err != nil || len(updatedOrderItems) == 0 {
				t.Logf("Failed to get updated order items: %v", err)
				return false
			}

			updatedAccountingCost := updatedOrderItems[0].AccountingCost
			updatedCostStatus := updatedOrderItems[0].CostStatus

			// Verify accounting cost remains unchanged
			if updatedAccountingCost != originalAccountingCost {
				t.Logf("Accounting cost changed from %f to %f", originalAccountingCost, updatedAccountingCost)
				return false
			}

			// Verify cost status remains FINAL
			if updatedCostStatus != originalCostStatus {
				t.Logf("Cost status changed from %s to %s", originalCostStatus, updatedCostStatus)
				return false
			}

			return true
		},
		gen.Float64Range(100, 10000),  // initialCost
		gen.Float64Range(100, 10000),  // newCost
	))

	properties.TestingRun(t, gopter.ConsoleReporter(false))
}

// TestAccountingCostImmutability_UnitTest is a concrete unit test that verifies
// accounting cost immutability with specific values
func TestAccountingCostImmutability_UnitTest(t *testing.T) {
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

	// Create test ingredient with initial cost
	testIngredient := &ingredient.Ingredient{
		ID:                primitive.NewObjectID(),
		Name:              "Coffee Beans",
		Category:          "Beverage",
		Unit:              ingredient.UnitGram,
		Quantity:          5000,
		MinStock:          500,
		CostPerUnit:       200, // 200 VND per gram
		ConversionRate:    1.0,
		WastagePercentage: 0.0,
		CreatedAt:         time.Now(),
		UpdatedAt:         time.Now(),
	}
	ingredientRepo.ingredients[testIngredient.ID] = testIngredient

	// Create test menu item
	testMenuItem := &menu.MenuItem{
		ID:          primitive.NewObjectID(),
		Name:        "Espresso",
		Price:       45000,
		Category:    "Coffee",
		Description: "Strong coffee",
		Ingredients: []menu.Ingredient{
			{
				Name:     testIngredient.Name,
				Quantity: 30, // 30 grams
				Unit:     testIngredient.Unit,
			},
		},
		Available: true,
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}
	menuRepo.menuItems[testMenuItem.ID] = testMenuItem

	// Create test shift (closed)
	testShift := &order.Shift{
		ID:        primitive.NewObjectID(),
		UserID:    primitive.NewObjectID(),
		UserName:  "waiter1",
		RoleType:  order.RoleWaiter,
		Status:    order.ShiftClosed,
		StartedAt: time.Now().Add(-3 * time.Hour),
		CreatedAt: time.Now().Add(-3 * time.Hour),
		UpdatedAt: time.Now().Add(-1 * time.Hour),
	}
	endedAt := time.Now().Add(-1 * time.Hour)
	testShift.EndedAt = &endedAt

	// Create test order
	testOrder := &order.Order{
		ID:       primitive.NewObjectID(),
		ShiftID:  testShift.ID,
		WaiterID: testShift.UserID,
		Items: []order.OrderItem{
			{
				MenuItemID: testMenuItem.ID,
				Name:       testMenuItem.Name,
				Price:      testMenuItem.Price,
				Quantity:   2,
				Subtotal:   testMenuItem.Price * 2,
			},
		},
		Subtotal:  testMenuItem.Price * 2,
		Total:     testMenuItem.Price * 2,
		Status:    order.StatusServed,
		CreatedAt: time.Now().Add(-2 * time.Hour),
		UpdatedAt: time.Now().Add(-2 * time.Hour),
	}
	orderRepo.orders[testOrder.ID] = testOrder

	// Calculate shift order costs (simulating shift closure)
	summary, err := costCalculator.CalculateShiftOrderCosts(ctx, testShift.ID)
	assert.NoError(t, err)
	assert.NotNil(t, summary)
	assert.Equal(t, 1, summary.TotalOrders)
	assert.Equal(t, 1, summary.TotalItems) // 1 order item (even though quantity = 2)

	// Get the order items with accounting cost
	orderItems, err := orderItemRepo.FindByOrderID(ctx, testOrder.ID)
	assert.NoError(t, err)
	assert.Len(t, orderItems, 1)

	originalAccountingCost := orderItems[0].AccountingCost
	originalCostStatus := orderItems[0].CostStatus

	// Verify cost status is FINAL
	assert.Equal(t, order.CostStatusFinal, originalCostStatus)

	// Verify accounting cost was calculated correctly
	// Expected: 30 grams * 200 VND/gram * 2 quantity = 12000 VND total
	expectedCost := 12000.0
	assert.InDelta(t, expectedCost, originalAccountingCost, 0.01)

	// Update ingredient cost to a much higher value
	testIngredient.CostPerUnit = 500 // Increase from 200 to 500
	testIngredient.UpdatedAt = time.Now()
	ingredientRepo.ingredients[testIngredient.ID] = testIngredient

	// Recalculate menu item cost (this should update current_cost but not accounting_cost)
	menuCost, err := costCalculator.CalculateMenuItemCost(ctx, testMenuItem.ID)
	assert.NoError(t, err)
	assert.NotNil(t, menuCost)

	// Verify menu item current_cost was updated
	// Expected new cost: 30 grams * 500 VND/gram = 15000 VND per item
	expectedNewCost := 15000.0
	assert.InDelta(t, expectedNewCost, menuCost.CurrentCost, 0.01)

	// Get the order items again
	updatedOrderItems, err := orderItemRepo.FindByOrderID(ctx, testOrder.ID)
	assert.NoError(t, err)
	assert.Len(t, updatedOrderItems, 1)

	updatedAccountingCost := updatedOrderItems[0].AccountingCost
	updatedCostStatus := updatedOrderItems[0].CostStatus

	// CRITICAL: Verify accounting cost remains unchanged
	assert.Equal(t, originalAccountingCost, updatedAccountingCost,
		"Accounting cost should remain unchanged after ingredient cost update")

	// Verify cost status remains FINAL
	assert.Equal(t, originalCostStatus, updatedCostStatus,
		"Cost status should remain FINAL after ingredient cost update")
}
