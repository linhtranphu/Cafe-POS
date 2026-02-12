package services

import (
	"context"
	"math"
	"testing"

	"cafe-pos/backend/domain/ingredient"
	"cafe-pos/backend/domain/menu"
	"cafe-pos/backend/domain/order"
	"github.com/leanovate/gopter"
	"github.com/leanovate/gopter/gen"
	"github.com/leanovate/gopter/prop"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

// Feature: menu-cost-profit-analysis, Property 5: Shift Closure Cost Calculation
// **Validates: Requirements 5.2, 5.3**
//
// Property: For any shift, when the shift is closed, the system should calculate
// accounting_cost for all order items in that shift using the current cost_per_unit
// values at the time of closure, using the same calculation method as current menu
// item cost, and mark the cost_status as "FINAL".
func TestProperty_ShiftClosureCostCalculation(t *testing.T) {
	parameters := gopter.DefaultTestParameters()
	parameters.MinSuccessfulTests = 100 // Minimum 100 iterations as per spec
	properties := gopter.NewProperties(parameters)

	properties.Property("Shift closure calculates accounting cost using same formula as current cost", prop.ForAll(
		func(shiftData testShiftData) bool {
			// Skip if no orders
			if len(shiftData.Orders) == 0 {
				return true
			}

			ctx := context.Background()

			// Setup test repositories
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

			// Create ingredients
			ingredientMap := make(map[string]*ingredient.Ingredient)
			for _, ingData := range shiftData.Ingredients {
				ing := &ingredient.Ingredient{
					ID:                primitive.NewObjectID(),
					Name:              ingData.Name,
					CostPerUnit:       ingData.CostPerUnit,
					ConversionRate:    ingData.ConversionRate,
					WastagePercentage: ingData.WastagePercentage,
				}
				ingredientRepo.ingredients[ing.ID] = ing
				ingredientMap[ing.Name] = ing
			}

			// Create menu items
			menuItemMap := make(map[string]*menu.MenuItem)
			for _, menuData := range shiftData.MenuItems {
				menuItem := &menu.MenuItem{
					ID:          primitive.NewObjectID(),
					Name:        menuData.Name,
					Price:       menuData.Price,
					Ingredients: []menu.Ingredient{},
				}

				// Add ingredients to menu item
				for _, ingRef := range menuData.Ingredients {
					menuItem.Ingredients = append(menuItem.Ingredients, menu.Ingredient{
						Name:     ingRef.Name,
						Quantity: ingRef.Quantity,
						Unit:     "unit",
					})
				}

				menuRepo.menuItems[menuItem.ID] = menuItem
				menuItemMap[menuItem.Name] = menuItem
			}

			// Create orders for the shift
			shiftID := primitive.NewObjectID()
			for _, orderData := range shiftData.Orders {
				ord := &order.Order{
					ID:      primitive.NewObjectID(),
					ShiftID: shiftID,
					Items:   []order.OrderItem{},
				}

				// Add items to order
				for _, itemData := range orderData.Items {
					menuItem, exists := menuItemMap[itemData.MenuItemName]
					if !exists {
						continue
					}

					ord.Items = append(ord.Items, order.OrderItem{
						MenuItemID: menuItem.ID,
						Name:       itemData.MenuItemName,
						Price:      menuItem.Price,
						Quantity:   itemData.Quantity,
						Subtotal:   menuItem.Price * float64(itemData.Quantity),
					})
				}

				if len(ord.Items) > 0 {
					orderRepo.orders[ord.ID] = ord
				}
			}

			// Skip if no valid orders were created
			if len(orderRepo.orders) == 0 {
				return true
			}

			// Create service
			service := NewCostCalculatorService(menuRepo, ingredientRepo, orderRepo, orderItemRepo)

			// Calculate shift order costs (simulating shift closure)
			result, err := service.CalculateShiftOrderCosts(ctx, shiftID)
			if err != nil {
				t.Logf("CalculateShiftOrderCosts failed: %v", err)
				return false
			}

			// Verify that order items were created
			if len(orderItemRepo.orderItems) == 0 {
				t.Logf("No order items created")
				return false
			}

			// Verify each order item
			for _, orderItem := range orderItemRepo.orderItems {
				// Find the corresponding menu item
				menuItem, exists := menuRepo.menuItems[orderItem.MenuItemID]
				if !exists {
					t.Logf("Menu item not found for order item")
					return false
				}

				// Calculate expected cost using the same formula as CalculateMenuItemCost
				expectedCostPerItem := 0.0
				hasIncompleteCost := false

				for _, menuIng := range menuItem.Ingredients {
					ing, exists := ingredientMap[menuIng.Name]
					if !exists || ing.CostPerUnit <= 0 {
						hasIncompleteCost = true
						continue
					}

					conversionRate := ing.ConversionRate
					if conversionRate <= 0 {
						conversionRate = 1.0
					}

					wastagePercentage := ing.WastagePercentage
					if wastagePercentage < 0 {
						wastagePercentage = 0.0
					}

					ingredientCost := menuIng.Quantity * ing.CostPerUnit * conversionRate * (1 + wastagePercentage/100)
					expectedCostPerItem += ingredientCost
				}

				// Round to 2 decimal places
				expectedCostPerItem = math.Round(expectedCostPerItem*100) / 100

				// Calculate total expected cost for the order item (cost per item * quantity)
				expectedTotalCost := expectedCostPerItem * float64(orderItem.Quantity)
				expectedTotalCost = math.Round(expectedTotalCost*100) / 100

				// Verify accounting cost matches expected cost
				tolerance := 0.01
				if math.Abs(orderItem.AccountingCost-expectedTotalCost) > tolerance {
					t.Logf("Accounting cost mismatch for %s: expected %.2f, got %.2f (diff: %.4f)",
						orderItem.Name, expectedTotalCost, orderItem.AccountingCost,
						math.Abs(orderItem.AccountingCost-expectedTotalCost))
					return false
				}

				// Verify cost status
				if hasIncompleteCost {
					if orderItem.CostStatus != order.CostStatusIncomplete {
						t.Logf("Expected INCOMPLETE status for %s with missing costs, got %s",
							orderItem.Name, orderItem.CostStatus)
						return false
					}
				} else {
					if orderItem.CostStatus != order.CostStatusFinal {
						t.Logf("Expected FINAL status for %s, got %s",
							orderItem.Name, orderItem.CostStatus)
						return false
					}
				}

				// Verify cost_calculated_at is set
				if orderItem.CostCalculatedAt.IsZero() {
					t.Logf("CostCalculatedAt not set for order item")
					return false
				}
			}

			// Verify summary statistics
			if result.TotalItems != len(orderItemRepo.orderItems) {
				t.Logf("Total items mismatch: expected %d, got %d",
					len(orderItemRepo.orderItems), result.TotalItems)
				return false
			}

			return true
		},
		genShiftData(),
	))

	properties.TestingRun(t, gopter.ConsoleReporter(false))
}

// Feature: menu-cost-profit-analysis, Property 5: Shift Closure Cost Calculation (Multiple Orders)
// **Validates: Requirements 5.2, 5.3**
//
// Property: For any shift with multiple orders, all order items should have their
// accounting_cost calculated and stored with FINAL status.
func TestProperty_ShiftClosureCostCalculation_MultipleOrders(t *testing.T) {
	parameters := gopter.DefaultTestParameters()
	parameters.MinSuccessfulTests = 100
	properties := gopter.NewProperties(parameters)

	properties.Property("All order items in shift get accounting cost with FINAL status", prop.ForAll(
		func(numOrders int, itemsPerOrder int, ingredientData []testIngredientData) bool {
			// Skip invalid cases
			if numOrders <= 0 || itemsPerOrder <= 0 || len(ingredientData) == 0 {
				return true
			}

			ctx := context.Background()

			// Setup repositories
			menuRepo := &mockMenuRepository{menuItems: make(map[primitive.ObjectID]*menu.MenuItem)}
			ingredientRepo := &mockIngredientRepository{ingredients: make(map[primitive.ObjectID]*ingredient.Ingredient)}
			orderRepo := &mockOrderRepository{orders: make(map[primitive.ObjectID]*order.Order)}
			orderItemRepo := &mockOrderItemRepository{orderItems: make([]*order.OrderItemWithCost, 0)}

			// Create ingredients (all with valid costs)
			for _, ingData := range ingredientData {
				if ingData.CostPerUnit <= 0 {
					continue
				}

				ing := &ingredient.Ingredient{
					ID:                primitive.NewObjectID(),
					Name:              ingData.Name,
					CostPerUnit:       ingData.CostPerUnit,
					ConversionRate:    ingData.ConversionRate,
					WastagePercentage: ingData.WastagePercentage,
				}
				ingredientRepo.ingredients[ing.ID] = ing
			}

			// Skip if no valid ingredients
			if len(ingredientRepo.ingredients) == 0 {
				return true
			}

			// Create a menu item using first ingredient
			firstIng := ingredientData[0]
			menuItem := &menu.MenuItem{
				ID:    primitive.NewObjectID(),
				Name:  "Test Item",
				Price: 50000,
				Ingredients: []menu.Ingredient{
					{Name: firstIng.Name, Quantity: 10, Unit: ingredient.UnitPiece},
				},
			}
			menuRepo.menuItems[menuItem.ID] = menuItem

			// Create orders
			shiftID := primitive.NewObjectID()
			expectedTotalItems := 0

			for i := 0; i < numOrders; i++ {
				ord := &order.Order{
					ID:      primitive.NewObjectID(),
					ShiftID: shiftID,
					Items:   []order.OrderItem{},
				}

				for j := 0; j < itemsPerOrder; j++ {
					ord.Items = append(ord.Items, order.OrderItem{
						MenuItemID: menuItem.ID,
						Name:       menuItem.Name,
						Price:      menuItem.Price,
						Quantity:   1,
						Subtotal:   menuItem.Price,
					})
					expectedTotalItems++
				}

				orderRepo.orders[ord.ID] = ord
			}

			// Create service and calculate shift costs
			service := NewCostCalculatorService(menuRepo, ingredientRepo, orderRepo, orderItemRepo)
			result, err := service.CalculateShiftOrderCosts(ctx, shiftID)
			if err != nil {
				t.Logf("CalculateShiftOrderCosts failed: %v", err)
				return false
			}

			// Verify all items were processed
			if result.TotalItems != expectedTotalItems {
				t.Logf("Expected %d total items, got %d", expectedTotalItems, result.TotalItems)
				return false
			}

			// Verify all order items were created
			if len(orderItemRepo.orderItems) != expectedTotalItems {
				t.Logf("Expected %d order items created, got %d", expectedTotalItems, len(orderItemRepo.orderItems))
				return false
			}

			// Verify all have FINAL status (since all ingredients have valid costs)
			for _, item := range orderItemRepo.orderItems {
				if item.CostStatus != order.CostStatusFinal {
					t.Logf("Expected FINAL status, got %s", item.CostStatus)
					return false
				}
				if item.AccountingCost <= 0 {
					t.Logf("Expected positive accounting cost, got %.2f", item.AccountingCost)
					return false
				}
			}

			return true
		},
		gen.IntRange(1, 5),  // 1-5 orders
		gen.IntRange(1, 3),  // 1-3 items per order
		genIngredientDataList(),
	))

	properties.TestingRun(t, gopter.ConsoleReporter(false))
}

// Feature: menu-cost-profit-analysis, Property 5: Shift Closure Cost Calculation (Rounding)
// **Validates: Requirements 5.2, 5.3**
//
// Property: Accounting costs calculated during shift closure should be rounded to 2 decimal places.
func TestProperty_ShiftClosureCostCalculation_Rounding(t *testing.T) {
	parameters := gopter.DefaultTestParameters()
	parameters.MinSuccessfulTests = 100
	properties := gopter.NewProperties(parameters)

	properties.Property("Accounting costs are rounded to 2 decimal places", prop.ForAll(
		func(costPerUnit float64, quantity float64, orderQuantity int) bool {
			// Skip invalid cases
			if costPerUnit <= 0 || quantity <= 0 || orderQuantity <= 0 {
				return true
			}

			ctx := context.Background()

			// Setup repositories
			menuRepo := &mockMenuRepository{menuItems: make(map[primitive.ObjectID]*menu.MenuItem)}
			ingredientRepo := &mockIngredientRepository{ingredients: make(map[primitive.ObjectID]*ingredient.Ingredient)}
			orderRepo := &mockOrderRepository{orders: make(map[primitive.ObjectID]*order.Order)}
			orderItemRepo := &mockOrderItemRepository{orderItems: make([]*order.OrderItemWithCost, 0)}

			// Create ingredient with specific cost that may produce many decimal places
			ing := &ingredient.Ingredient{
				ID:                primitive.NewObjectID(),
				Name:              "TestIngredient",
				CostPerUnit:       costPerUnit,
				ConversionRate:    1.0,
				WastagePercentage: 0.0,
			}
			ingredientRepo.ingredients[ing.ID] = ing

			// Create menu item
			menuItem := &menu.MenuItem{
				ID:    primitive.NewObjectID(),
				Name:  "Test Item",
				Price: 50000,
				Ingredients: []menu.Ingredient{
					{Name: "TestIngredient", Quantity: quantity, Unit: ingredient.UnitPiece},
				},
			}
			menuRepo.menuItems[menuItem.ID] = menuItem

			// Create order
			shiftID := primitive.NewObjectID()
			ord := &order.Order{
				ID:      primitive.NewObjectID(),
				ShiftID: shiftID,
				Items: []order.OrderItem{
					{
						MenuItemID: menuItem.ID,
						Name:       menuItem.Name,
						Price:      menuItem.Price,
						Quantity:   orderQuantity,
						Subtotal:   menuItem.Price * float64(orderQuantity),
					},
				},
			}
			orderRepo.orders[ord.ID] = ord

			// Calculate shift costs
			service := NewCostCalculatorService(menuRepo, ingredientRepo, orderRepo, orderItemRepo)
			_, err := service.CalculateShiftOrderCosts(ctx, shiftID)
			if err != nil {
				t.Logf("CalculateShiftOrderCosts failed: %v", err)
				return false
			}

			// Verify order item was created
			if len(orderItemRepo.orderItems) != 1 {
				t.Logf("Expected 1 order item, got %d", len(orderItemRepo.orderItems))
				return false
			}

			orderItem := orderItemRepo.orderItems[0]

			// Verify rounding to 2 decimal places
			rounded := math.Round(orderItem.AccountingCost*100) / 100
			if math.Abs(orderItem.AccountingCost-rounded) > 0.0001 {
				t.Logf("Accounting cost not properly rounded to 2 decimals: %.10f (rounded: %.2f)",
					orderItem.AccountingCost, rounded)
				return false
			}

			return true
		},
		gen.Float64Range(1.0, 10000.0),   // Cost per unit
		gen.Float64Range(0.1, 1000.0),    // Quantity in recipe
		gen.IntRange(1, 10),               // Order quantity
	))

	properties.TestingRun(t, gopter.ConsoleReporter(false))
}

// Test data structures for shift closure property testing

type testShiftData struct {
	Ingredients []testIngredientData
	MenuItems   []testMenuItemData
	Orders      []testOrderData
}

type testMenuItemData struct {
	Name        string
	Price       float64
	Ingredients []testIngredientReference
}

type testIngredientReference struct {
	Name     string
	Quantity float64
}

type testOrderData struct {
	Items []testOrderItemData
}

type testOrderItemData struct {
	MenuItemName string
	Quantity     int
}

// Generator for shift data
func genShiftData() gopter.Gen {
	return gopter.CombineGens(
		genIngredientDataList(),
		genMenuItemDataList(),
		genOrderDataList(),
	).Map(func(values []interface{}) testShiftData {
		return testShiftData{
			Ingredients: values[0].([]testIngredientData),
			MenuItems:   values[1].([]testMenuItemData),
			Orders:      values[2].([]testOrderData),
		}
	})
}

// Generator for menu item data list
func genMenuItemDataList() gopter.Gen {
	return gen.SliceOf(genMenuItemData()).
		Map(func(slice []testMenuItemData) []testMenuItemData {
			if len(slice) == 0 {
				return []testMenuItemData{
					{
						Name:  "DefaultItem",
						Price: 50000,
						Ingredients: []testIngredientReference{
							{Name: "DefaultIngredient", Quantity: 10},
						},
					},
				}
			}
			if len(slice) > 5 {
				return slice[:5]
			}
			return slice
		})
}

// Generator for single menu item data
func genMenuItemData() gopter.Gen {
	return gopter.CombineGens(
		gen.Identifier(),
		gen.Float64Range(10000, 100000),
		genIngredientReferenceList(),
	).Map(func(values []interface{}) testMenuItemData {
		return testMenuItemData{
			Name:        values[0].(string),
			Price:       values[1].(float64),
			Ingredients: values[2].([]testIngredientReference),
		}
	})
}

// Generator for ingredient reference list
func genIngredientReferenceList() gopter.Gen {
	return gen.SliceOf(genIngredientReference()).
		Map(func(slice []testIngredientReference) []testIngredientReference {
			if len(slice) == 0 {
				return []testIngredientReference{
					{Name: "DefaultIngredient", Quantity: 10},
				}
			}
			if len(slice) > 5 {
				return slice[:5]
			}
			return slice
		})
}

// Generator for ingredient reference
func genIngredientReference() gopter.Gen {
	return gopter.CombineGens(
		gen.Identifier(),
		gen.Float64Range(0.1, 1000.0),
	).Map(func(values []interface{}) testIngredientReference {
		return testIngredientReference{
			Name:     values[0].(string),
			Quantity: values[1].(float64),
		}
	})
}

// Generator for order data list
func genOrderDataList() gopter.Gen {
	return gen.SliceOf(genOrderData()).
		Map(func(slice []testOrderData) []testOrderData {
			if len(slice) == 0 {
				return []testOrderData{
					{
						Items: []testOrderItemData{
							{MenuItemName: "DefaultItem", Quantity: 1},
						},
					},
				}
			}
			if len(slice) > 5 {
				return slice[:5]
			}
			return slice
		})
}

// Generator for single order data
func genOrderData() gopter.Gen {
	return gen.SliceOf(genOrderItemData()).
		Map(func(slice []testOrderItemData) testOrderData {
			if len(slice) == 0 {
				return testOrderData{
					Items: []testOrderItemData{
						{MenuItemName: "DefaultItem", Quantity: 1},
					},
				}
			}
			if len(slice) > 5 {
				slice = slice[:5]
			}
			return testOrderData{Items: slice}
		})
}

// Generator for order item data
func genOrderItemData() gopter.Gen {
	return gopter.CombineGens(
		gen.Identifier(),
		gen.IntRange(1, 10),
	).Map(func(values []interface{}) testOrderItemData {
		return testOrderItemData{
			MenuItemName: values[0].(string),
			Quantity:     values[1].(int),
		}
	})
}
