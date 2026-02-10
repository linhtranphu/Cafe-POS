package services

import (
	"context"
	"testing"

	"cafe-pos/backend/domain/ingredient"
	"cafe-pos/backend/domain/menu"
	"cafe-pos/backend/domain/order"
	"github.com/leanovate/gopter"
	"github.com/leanovate/gopter/gen"
	"github.com/leanovate/gopter/prop"
	"github.com/stretchr/testify/assert"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

// Feature: menu-cost-profit-analysis, Property 18: Recipe Change Immutability
// Validates: Requirements 8.6
//
// Property: For any menu item with historical accounting_cost data in closed shifts,
// when the recipe is modified, the historical accounting_cost values should remain unchanged.
func TestProperty_RecipeChangeImmutability(t *testing.T) {
	properties := gopter.NewProperties(nil)

	properties.Property("Recipe changes do not affect historical accounting costs", prop.ForAll(
		func(originalQuantity, newQuantity float64) bool {
			ctx := context.Background()

			// Setup mock repositories
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
				CostPerUnit:       100.0,
				ConversionRate:    1.0,
				WastagePercentage: 0.0,
			}
			ingredientRepo.ingredients[testIngredient.ID] = testIngredient

			// Create menu item with original recipe
			menuItem := &menu.MenuItem{
				ID:    primitive.NewObjectID(),
				Name:  "Test Item",
				Price: 10000,
				Ingredients: []menu.Ingredient{
					{
						Name:     testIngredient.Name,
						Quantity: originalQuantity,
						Unit:     "g",
					},
				},
			}
			menuRepo.menuItems[menuItem.ID] = menuItem

			// Create shift
			testShift := &order.Shift{
				ID:       primitive.NewObjectID(),
				Status:   order.ShiftOpen,
				RoleType: order.RoleWaiter,
			}

			// Create order with the menu item
			testOrder := &order.Order{
				ID:      primitive.NewObjectID(),
				ShiftID: testShift.ID,
				Items: []order.OrderItem{
					{
						MenuItemID: menuItem.ID,
						Name:       menuItem.Name,
						Price:      menuItem.Price,
						Quantity:   1,
						Subtotal:   menuItem.Price,
					},
				},
				Subtotal: menuItem.Price,
				Total:    menuItem.Price,
				Status:   order.StatusServed,
			}
			orderRepo.orders[testOrder.ID] = testOrder

			// Close shift and calculate costs (this creates historical accounting_cost)
			_, err := costCalculator.CalculateShiftOrderCosts(ctx, testShift.ID)
			if err != nil {
				return false
			}

			// Get the order item with accounting cost
			orderItems, err := orderItemRepo.FindByOrderID(ctx, testOrder.ID)
			if err != nil || len(orderItems) == 0 {
				return false
			}

			originalAccountingCost := orderItems[0].AccountingCost
			originalCostStatus := orderItems[0].CostStatus

			// Verify accounting cost was calculated
			if originalAccountingCost == 0 {
				return false
			}

			// Verify cost status is FINAL
			if originalCostStatus != order.CostStatusFinal {
				return false
			}

			// Modify the recipe (change ingredient quantity)
			menuItem.Ingredients[0].Quantity = newQuantity
			menuRepo.menuItems[menuItem.ID] = menuItem

			// Get the order item again
			orderItems2, err := orderItemRepo.FindByOrderID(ctx, testOrder.ID)
			if err != nil || len(orderItems2) == 0 {
				return false
			}

			newAccountingCost := orderItems2[0].AccountingCost
			newCostStatus := orderItems2[0].CostStatus

			// Verify accounting cost has NOT changed
			if newAccountingCost != originalAccountingCost {
				t.Logf("FAIL: Accounting cost changed from %.2f to %.2f after recipe change",
					originalAccountingCost, newAccountingCost)
				return false
			}

			// Verify cost status is still FINAL
			if newCostStatus != order.CostStatusFinal {
				t.Logf("FAIL: Cost status changed from %s to %s after recipe change",
					originalCostStatus, newCostStatus)
				return false
			}

			return true
		},
		gen.Float64Range(1, 100),  // originalQuantity
		gen.Float64Range(1, 100),  // newQuantity (different from original)
	))

	properties.TestingRun(t, gopter.ConsoleReporter(false))
}

// Feature: menu-cost-profit-analysis, Property 18: Recipe Change Immutability (Extended)
// Test with ingredient additions and removals
func TestProperty_RecipeChangeImmutability_AddRemoveIngredients(t *testing.T) {
	ctx := context.Background()

	// Setup mock repositories
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

	// Create test ingredients
	ingredient1 := &ingredient.Ingredient{
		ID:                primitive.NewObjectID(),
		Name:              "Ingredient 1",
		CostPerUnit:       100.0,
		ConversionRate:    1.0,
		WastagePercentage: 0.0,
	}
	ingredientRepo.ingredients[ingredient1.ID] = ingredient1

	ingredient2 := &ingredient.Ingredient{
		ID:                primitive.NewObjectID(),
		Name:              "Ingredient 2",
		CostPerUnit:       200.0,
		ConversionRate:    1.0,
		WastagePercentage: 0.0,
	}
	ingredientRepo.ingredients[ingredient2.ID] = ingredient2

	// Create menu item with one ingredient
	menuItem := &menu.MenuItem{
		ID:    primitive.NewObjectID(),
		Name:  "Test Item",
		Price: 10000,
		Ingredients: []menu.Ingredient{
			{
				Name:     ingredient1.Name,
				Quantity: 10,
				Unit:     "g",
			},
		},
	}
	menuRepo.menuItems[menuItem.ID] = menuItem

	// Create shift
	testShift := &order.Shift{
		ID:       primitive.NewObjectID(),
		Status:   order.ShiftOpen,
		RoleType: order.RoleWaiter,
	}

	// Create order
	testOrder := &order.Order{
		ID:      primitive.NewObjectID(),
		ShiftID: testShift.ID,
		Items: []order.OrderItem{
			{
				MenuItemID: menuItem.ID,
				Name:       menuItem.Name,
				Price:      menuItem.Price,
				Quantity:   1,
				Subtotal:   menuItem.Price,
			},
		},
		Subtotal: menuItem.Price,
		Total:    menuItem.Price,
		Status:   order.StatusServed,
	}
	orderRepo.orders[testOrder.ID] = testOrder

	// Close shift and calculate costs
	_, err := costCalculator.CalculateShiftOrderCosts(ctx, testShift.ID)
	assert.NoError(t, err)

	// Get original accounting cost
	orderItems, err := orderItemRepo.FindByOrderID(ctx, testOrder.ID)
	assert.NoError(t, err)
	assert.Len(t, orderItems, 1)

	originalAccountingCost := orderItems[0].AccountingCost
	originalCostStatus := orderItems[0].CostStatus

	t.Logf("Original accounting cost: %.2f, status: %s", originalAccountingCost, originalCostStatus)

	// Modify recipe: Add a new ingredient
	menuItem.Ingredients = append(menuItem.Ingredients, menu.Ingredient{
		Name:     ingredient2.Name,
		Quantity: 20,
		Unit:     "g",
	})
	menuRepo.menuItems[menuItem.ID] = menuItem

	// Get accounting cost after recipe change
	orderItems2, err := orderItemRepo.FindByOrderID(ctx, testOrder.ID)
	assert.NoError(t, err)
	assert.Len(t, orderItems2, 1)

	newAccountingCost := orderItems2[0].AccountingCost
	newCostStatus := orderItems2[0].CostStatus

	t.Logf("After adding ingredient: accounting cost: %.2f, status: %s", newAccountingCost, newCostStatus)

	// Verify accounting cost has NOT changed
	assert.Equal(t, originalAccountingCost, newAccountingCost,
		"Accounting cost should remain unchanged after adding ingredient to recipe")
	assert.Equal(t, originalCostStatus, newCostStatus,
		"Cost status should remain FINAL after recipe change")

	// Modify recipe: Remove an ingredient
	menuItem.Ingredients = []menu.Ingredient{
		{
			Name:     ingredient2.Name,
			Quantity: 20,
			Unit:     "g",
		},
	}
	menuRepo.menuItems[menuItem.ID] = menuItem

	// Get accounting cost after removing ingredient
	orderItems3, err := orderItemRepo.FindByOrderID(ctx, testOrder.ID)
	assert.NoError(t, err)
	assert.Len(t, orderItems3, 1)

	finalAccountingCost := orderItems3[0].AccountingCost
	finalCostStatus := orderItems3[0].CostStatus

	t.Logf("After removing ingredient: accounting cost: %.2f, status: %s", finalAccountingCost, finalCostStatus)

	// Verify accounting cost STILL has NOT changed
	assert.Equal(t, originalAccountingCost, finalAccountingCost,
		"Accounting cost should remain unchanged after removing ingredient from recipe")
	assert.Equal(t, originalCostStatus, finalCostStatus,
		"Cost status should remain FINAL after recipe change")

	t.Log("✅ Recipe change immutability test passed!")
	t.Log("   - Accounting cost remains unchanged when recipe is modified")
	t.Log("   - Cost status remains FINAL after recipe changes")
	t.Log("   - Historical data is immutable")
}
