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
	"go.mongodb.org/mongo-driver/bson/primitive"
)

// Feature: menu-cost-profit-analysis, Property 3: Background Job Queuing on Ingredient Update
// **Validates: Requirements 1.3, 9.1**
//
// Property: For any ingredient, when its cost_per_unit is updated, the system should queue
// a background job to recalculate current_cost for all menu items that use that ingredient.
func TestProperty_BackgroundJobQueueing(t *testing.T) {
	parameters := gopter.DefaultTestParameters()
	parameters.MinSuccessfulTests = 20 // Reduced for faster execution
	properties := gopter.NewProperties(parameters)

	properties.Property("All menu items using an ingredient are queued when ingredient cost changes", prop.ForAll(
		func(queueTestData testQueueData) bool {
			// Skip if no menu items
			if len(queueTestData.MenuItems) == 0 {
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

			// Create the target ingredient (the one whose cost will be updated)
			targetIngredient := &ingredient.Ingredient{
				ID:                primitive.NewObjectID(),
				Name:              queueTestData.TargetIngredientName,
				CostPerUnit:       queueTestData.InitialCost,
				ConversionRate:    1.0,
				WastagePercentage: 0.0,
			}
			ingredientRepo.ingredients[targetIngredient.ID] = targetIngredient

			// Create other ingredients
			for _, otherIngName := range queueTestData.OtherIngredients {
				otherIng := &ingredient.Ingredient{
					ID:                primitive.NewObjectID(),
					Name:              otherIngName,
					CostPerUnit:       100.0,
					ConversionRate:    1.0,
					WastagePercentage: 0.0,
				}
				ingredientRepo.ingredients[otherIng.ID] = otherIng
			}

			// Create menu items
			// Track which menu items should be queued (those using target ingredient)
			expectedQueuedItems := make(map[primitive.ObjectID]bool)
			notExpectedQueuedItems := make(map[primitive.ObjectID]bool)

			for _, menuData := range queueTestData.MenuItems {
				menuItem := &menu.MenuItem{
					ID:          primitive.NewObjectID(),
					Name:        menuData.Name,
					Price:       menuData.Price,
					Ingredients: []menu.Ingredient{},
				}

				usesTargetIngredient := false
				for _, ingName := range menuData.IngredientNames {
					menuItem.Ingredients = append(menuItem.Ingredients, menu.Ingredient{
						Name:     ingName,
						Quantity: 10.0,
						Unit:     "unit",
					})

					if ingName == queueTestData.TargetIngredientName {
						usesTargetIngredient = true
					}
				}

				menuRepo.menuItems[menuItem.ID] = menuItem

				if usesTargetIngredient {
					expectedQueuedItems[menuItem.ID] = true
				} else {
					notExpectedQueuedItems[menuItem.ID] = true
				}
			}

			// Create service
			service := NewCostCalculatorService(menuRepo, ingredientRepo, orderRepo, orderItemRepo)

			// Create recalculation service and wire it up
			recalcService := NewCostRecalculationService(service, menuRepo, 1, 100)
			service.SetCostRecalculationService(recalcService)
			recalcService.Start()
			defer recalcService.Stop()

			// Queue cost recalculation for the target ingredient
			err := service.QueueCostRecalculation(ctx, targetIngredient.ID)
			if err != nil {
				t.Logf("QueueCostRecalculation failed: %v", err)
				return false
			}

			// Wait for recalculations to complete
			maxWaitTime := 3 * time.Second
			startTime := time.Now()
			for {
				status, _ := recalcService.GetRecalculationStatus(ctx)
				if status.QueuedItems == 0 {
					break
				}
				if time.Since(startTime) > maxWaitTime {
					t.Logf("Timeout waiting for recalculations to complete")
					return false
				}
				time.Sleep(50 * time.Millisecond)
			}

			// Give workers time to finish processing
			time.Sleep(100 * time.Millisecond)

			// Verify: All menu items using the target ingredient were recalculated
			// by checking that their CurrentCost was updated
			for expectedID := range expectedQueuedItems {
				menuItem, err := menuRepo.FindByID(ctx, expectedID)
				if err != nil {
					t.Logf("Failed to find menu item %s: %v", expectedID.Hex(), err)
					return false
				}
				
				// Menu items that use the target ingredient should have their cost calculated
				if menuItem.CurrentCost == 0 && len(menuItem.Ingredients) > 0 {
					t.Logf("Expected menu item %s to have cost calculated", expectedID.Hex())
					return false
				}
			}

			// Verify: Menu items NOT using the target ingredient were NOT recalculated
			// (their cost should remain 0)
			for notExpectedID := range notExpectedQueuedItems {
				menuItem, err := menuRepo.FindByID(ctx, notExpectedID)
				if err != nil {
					t.Logf("Failed to find menu item %s: %v", notExpectedID.Hex(), err)
					return false
				}
				
				// Menu items that don't use the target ingredient should not be recalculated
				if menuItem.CurrentCost != 0 {
					t.Logf("Menu item %s was recalculated but shouldn't have been", notExpectedID.Hex())
					return false
				}
			}

			return true
		},
		genQueueTestData(),
	))

	properties.TestingRun(t, gopter.ConsoleReporter(false))
}

// Feature: menu-cost-profit-analysis, Property 3: Background Job Queuing (Multiple Ingredients)
// **Validates: Requirements 1.3, 9.1**
//
// Property: When multiple menu items use the same ingredient, all of them should be queued
// when that ingredient's cost changes.
func TestProperty_BackgroundJobQueueing_MultipleMenuItems(t *testing.T) {
	parameters := gopter.DefaultTestParameters()
	parameters.MinSuccessfulTests = 20
	properties := gopter.NewProperties(parameters)

	properties.Property("All menu items using ingredient are queued regardless of count", prop.ForAll(
		func(numMenuItems int, ingredientName string) bool {
			// Skip invalid cases
			if numMenuItems <= 0 || numMenuItems > 50 {
				return true
			}

			ctx := context.Background()

			// Setup repositories
			menuRepo := &mockMenuRepository{menuItems: make(map[primitive.ObjectID]*menu.MenuItem)}
			ingredientRepo := &mockIngredientRepository{ingredients: make(map[primitive.ObjectID]*ingredient.Ingredient)}
			orderRepo := &mockOrderRepository{orders: make(map[primitive.ObjectID]*order.Order)}
			orderItemRepo := &mockOrderItemRepository{orderItems: make([]*order.OrderItemWithCost, 0)}

			// Create target ingredient
			targetIng := &ingredient.Ingredient{
				ID:                primitive.NewObjectID(),
				Name:              ingredientName,
				CostPerUnit:       100.0,
				ConversionRate:    1.0,
				WastagePercentage: 0.0,
			}
			ingredientRepo.ingredients[targetIng.ID] = targetIng

			// Create menu items that all use the target ingredient
			for i := 0; i < numMenuItems; i++ {
				menuItemName, _ := gen.Identifier().Sample()
				menuItem := &menu.MenuItem{
					ID:    primitive.NewObjectID(),
					Name:  menuItemName.(string),
					Price: 50000,
					Ingredients: []menu.Ingredient{
						{Name: ingredientName, Quantity: 10, Unit: "unit"},
					},
				}
				menuRepo.menuItems[menuItem.ID] = menuItem
			}

			// Create service and queue recalculation
			service := NewCostCalculatorService(menuRepo, ingredientRepo, orderRepo, orderItemRepo)
			
			// Create recalculation service and wire it up
			recalcService := NewCostRecalculationService(service, menuRepo, 1, 100)
			service.SetCostRecalculationService(recalcService)
			recalcService.Start()
			defer recalcService.Stop()
			
			err := service.QueueCostRecalculation(ctx, targetIng.ID)
			if err != nil {
				t.Logf("QueueCostRecalculation failed: %v", err)
				return false
			}

			// Wait for recalculations to complete
			maxWaitTime := 3 * time.Second
			startTime := time.Now()
			for {
				status, _ := recalcService.GetRecalculationStatus(ctx)
				if status.QueuedItems == 0 {
					break
				}
				if time.Since(startTime) > maxWaitTime {
					t.Logf("Timeout waiting for recalculations to complete")
					return false
				}
				time.Sleep(50 * time.Millisecond)
			}

			// Give workers time to finish
			time.Sleep(100 * time.Millisecond)

			// Verify all menu items were recalculated by checking their costs were updated
			for _, menuItem := range menuRepo.menuItems {
				if menuItem.CurrentCost == 0 && len(menuItem.Ingredients) > 0 {
					t.Logf("Expected menu item %s to have cost calculated", menuItem.ID.Hex())
					return false
				}
			}

			return true
		},
		gen.IntRange(1, 50),
		gen.Identifier(),
	))

	properties.TestingRun(t, gopter.ConsoleReporter(false))
}

// Feature: menu-cost-profit-analysis, Property 3: Background Job Queuing (No Menu Items)
// **Validates: Requirements 1.3, 9.1**
//
// Property: When an ingredient is not used by any menu items, no jobs should be queued.
func TestProperty_BackgroundJobQueueing_NoMenuItems(t *testing.T) {
	parameters := gopter.DefaultTestParameters()
	parameters.MinSuccessfulTests = 20
	properties := gopter.NewProperties(parameters)

	properties.Property("No jobs queued when ingredient not used by any menu items", prop.ForAll(
		func(targetIngName string, otherIngNames []string, menuItemData []testQueueMenuItem) bool {
			// Skip if target ingredient name conflicts with other ingredients
			for _, otherName := range otherIngNames {
				if otherName == targetIngName {
					return true
				}
			}

			ctx := context.Background()

			// Setup repositories
			menuRepo := &mockMenuRepository{menuItems: make(map[primitive.ObjectID]*menu.MenuItem)}
			ingredientRepo := &mockIngredientRepository{ingredients: make(map[primitive.ObjectID]*ingredient.Ingredient)}
			orderRepo := &mockOrderRepository{orders: make(map[primitive.ObjectID]*order.Order)}
			orderItemRepo := &mockOrderItemRepository{orderItems: make([]*order.OrderItemWithCost, 0)}

			// Create target ingredient (not used by any menu items)
			targetIng := &ingredient.Ingredient{
				ID:                primitive.NewObjectID(),
				Name:              targetIngName,
				CostPerUnit:       100.0,
				ConversionRate:    1.0,
				WastagePercentage: 0.0,
			}
			ingredientRepo.ingredients[targetIng.ID] = targetIng

			// Create other ingredients
			for _, otherName := range otherIngNames {
				otherIng := &ingredient.Ingredient{
					ID:                primitive.NewObjectID(),
					Name:              otherName,
					CostPerUnit:       100.0,
					ConversionRate:    1.0,
					WastagePercentage: 0.0,
				}
				ingredientRepo.ingredients[otherIng.ID] = otherIng
			}

			// Create menu items that DON'T use the target ingredient
			for _, menuData := range menuItemData {
				menuItem := &menu.MenuItem{
					ID:          primitive.NewObjectID(),
					Name:        menuData.Name,
					Price:       menuData.Price,
					Ingredients: []menu.Ingredient{},
				}

				// Only use other ingredients, not the target
				for _, ingName := range menuData.IngredientNames {
					// Skip if this is the target ingredient
					if ingName == targetIngName {
						continue
					}

					menuItem.Ingredients = append(menuItem.Ingredients, menu.Ingredient{
						Name:     ingName,
						Quantity: 10,
						Unit:     "unit",
					})
				}

				// Only add menu item if it has ingredients
				if len(menuItem.Ingredients) > 0 {
					menuRepo.menuItems[menuItem.ID] = menuItem
				}
			}

			// Create service and queue recalculation
			service := NewCostCalculatorService(menuRepo, ingredientRepo, orderRepo, orderItemRepo)
			
			// Create recalculation service and wire it up
			recalcService := NewCostRecalculationService(service, menuRepo, 1, 100)
			service.SetCostRecalculationService(recalcService)
			recalcService.Start()
			defer recalcService.Stop()
			
			err := service.QueueCostRecalculation(ctx, targetIng.ID)
			if err != nil {
				t.Logf("QueueCostRecalculation failed: %v", err)
				return false
			}

			// Wait for any potential recalculations
			time.Sleep(300 * time.Millisecond)

			// Verify no menu items were recalculated (all should have cost = 0)
			// since none of them use the target ingredient
			for _, menuItem := range menuRepo.menuItems {
				if menuItem.CurrentCost != 0 {
					t.Logf("Menu item %s was recalculated but shouldn't have been", menuItem.ID.Hex())
					return false
				}
			}

			return true
		},
		gen.Identifier(),
		gen.SliceOf(gen.Identifier()).Map(func(slice []string) []string {
			if len(slice) > 5 {
				return slice[:5]
			}
			return slice
		}),
		genQueueMenuItemList(),
	))

	properties.TestingRun(t, gopter.ConsoleReporter(false))
}

// Feature: menu-cost-profit-analysis, Property 3: Background Job Queuing (Selective Queuing)
// **Validates: Requirements 1.3, 9.1**
//
// Property: Only menu items that use the specific ingredient should be queued,
// not all menu items in the system.
func TestProperty_BackgroundJobQueueing_SelectiveQueuing(t *testing.T) {
	parameters := gopter.DefaultTestParameters()
	parameters.MinSuccessfulTests = 20
	properties := gopter.NewProperties(parameters)

	properties.Property("Only menu items using specific ingredient are queued", prop.ForAll(
		func(targetIngName string, otherIngName string, numWithTarget int, numWithoutTarget int) bool {
			// Skip invalid cases
			if numWithTarget <= 0 || numWithoutTarget <= 0 || numWithTarget > 20 || numWithoutTarget > 20 {
				return true
			}
			if targetIngName == otherIngName {
				return true
			}

			ctx := context.Background()

			// Setup repositories
			menuRepo := &mockMenuRepository{menuItems: make(map[primitive.ObjectID]*menu.MenuItem)}
			ingredientRepo := &mockIngredientRepository{ingredients: make(map[primitive.ObjectID]*ingredient.Ingredient)}
			orderRepo := &mockOrderRepository{orders: make(map[primitive.ObjectID]*order.Order)}
			orderItemRepo := &mockOrderItemRepository{orderItems: make([]*order.OrderItemWithCost, 0)}

			// Create ingredients
			targetIng := &ingredient.Ingredient{
				ID:                primitive.NewObjectID(),
				Name:              targetIngName,
				CostPerUnit:       100.0,
				ConversionRate:    1.0,
				WastagePercentage: 0.0,
			}
			ingredientRepo.ingredients[targetIng.ID] = targetIng

			otherIng := &ingredient.Ingredient{
				ID:                primitive.NewObjectID(),
				Name:              otherIngName,
				CostPerUnit:       100.0,
				ConversionRate:    1.0,
				WastagePercentage: 0.0,
			}
			ingredientRepo.ingredients[otherIng.ID] = otherIng

			// Create menu items WITH target ingredient
			for i := 0; i < numWithTarget; i++ {
				menuItemName, _ := gen.Identifier().Sample()
				menuItem := &menu.MenuItem{
					ID:    primitive.NewObjectID(),
					Name:  menuItemName.(string) + "_with_target",
					Price: 50000,
					Ingredients: []menu.Ingredient{
						{Name: targetIngName, Quantity: 10, Unit: "unit"},
					},
				}
				menuRepo.menuItems[menuItem.ID] = menuItem
			}

			// Create menu items WITHOUT target ingredient
			for i := 0; i < numWithoutTarget; i++ {
				menuItemName, _ := gen.Identifier().Sample()
				menuItem := &menu.MenuItem{
					ID:    primitive.NewObjectID(),
					Name:  menuItemName.(string) + "_without_target",
					Price: 50000,
					Ingredients: []menu.Ingredient{
						{Name: otherIngName, Quantity: 10, Unit: "unit"},
					},
				}
				menuRepo.menuItems[menuItem.ID] = menuItem
			}

			// Create service and queue recalculation
			service := NewCostCalculatorService(menuRepo, ingredientRepo, orderRepo, orderItemRepo)
			
			// Create recalculation service and wire it up
			recalcService := NewCostRecalculationService(service, menuRepo, 1, 100)
			service.SetCostRecalculationService(recalcService)
			recalcService.Start()
			defer recalcService.Stop()
			
			err := service.QueueCostRecalculation(ctx, targetIng.ID)
			if err != nil {
				t.Logf("QueueCostRecalculation failed: %v", err)
				return false
			}

			// Wait for recalculations to complete
			maxWaitTime := 3 * time.Second
			startTime := time.Now()
			for {
				status, _ := recalcService.GetRecalculationStatus(ctx)
				if status.QueuedItems == 0 {
					break
				}
				if time.Since(startTime) > maxWaitTime {
					t.Logf("Timeout waiting for recalculations to complete")
					return false
				}
				time.Sleep(50 * time.Millisecond)
			}

			// Give workers time to finish
			time.Sleep(100 * time.Millisecond)

			// Verify: Only menu items with target ingredient were recalculated
			recalculatedCount := 0
			for _, menuItem := range menuRepo.menuItems {
				if menuItem.CurrentCost > 0 {
					recalculatedCount++
					
					// Verify this menu item uses the target ingredient
					usesTarget := false
					for _, ing := range menuItem.Ingredients {
						if ing.Name == targetIngName {
							usesTarget = true
							break
						}
					}
					
					if !usesTarget {
						t.Logf("Menu item %s was recalculated but doesn't use target ingredient", menuItem.ID.Hex())
						return false
					}
				}
			}

			// Verify the count matches expected
			if recalculatedCount != numWithTarget {
				t.Logf("Expected %d items recalculated, got %d", numWithTarget, recalculatedCount)
				return false
			}

			return true
		},
		gen.Identifier(),
		gen.Identifier(),
		gen.IntRange(1, 20),
		gen.IntRange(1, 20),
	))

	properties.TestingRun(t, gopter.ConsoleReporter(false))
}

// Test data structures for queue property testing

type testQueueData struct {
	TargetIngredientName string
	InitialCost          float64
	UpdatedCost          float64
	OtherIngredients     []string
	MenuItems            []testQueueMenuItem
}

type testQueueMenuItem struct {
	Name            string
	Price           float64
	IngredientNames []string
}

// Generator for queue test data
func genQueueTestData() gopter.Gen {
	return gopter.CombineGens(
		gen.Identifier(),                // Target ingredient name
		gen.Float64Range(100.0, 1000.0), // Initial cost
		gen.Float64Range(100.0, 1000.0), // Updated cost
		genIngredientNameList(),         // Other ingredients
		genQueueMenuItemList(),          // Menu items
	).Map(func(values []interface{}) testQueueData {
		return testQueueData{
			TargetIngredientName: values[0].(string),
			InitialCost:          values[1].(float64),
			UpdatedCost:          values[2].(float64),
			OtherIngredients:     values[3].([]string),
			MenuItems:            values[4].([]testQueueMenuItem),
		}
	})
}

// Generator for ingredient name list
func genIngredientNameList() gopter.Gen {
	return gen.SliceOf(gen.Identifier()).
		Map(func(slice []string) []string {
			if len(slice) > 5 {
				return slice[:5]
			}
			return slice
		})
}

// Generator for queue menu item list
func genQueueMenuItemList() gopter.Gen {
	return gen.SliceOf(genQueueMenuItem()).
		Map(func(slice []testQueueMenuItem) []testQueueMenuItem {
			if len(slice) == 0 {
				return []testQueueMenuItem{
					{
						Name:            "DefaultItem",
						Price:           50000,
						IngredientNames: []string{"DefaultIngredient"},
					},
				}
			}
			if len(slice) > 10 {
				return slice[:10]
			}
			return slice
		})
}

// Generator for single queue menu item
func genQueueMenuItem() gopter.Gen {
	return gopter.CombineGens(
		gen.Identifier(),
		gen.Float64Range(10000, 100000),
		genIngredientNameList(),
	).Map(func(values []interface{}) testQueueMenuItem {
		ingredientNames := values[2].([]string)
		// Ensure at least one ingredient
		if len(ingredientNames) == 0 {
			ingredientNames = []string{"DefaultIngredient"}
		}
		return testQueueMenuItem{
			Name:            values[0].(string),
			Price:           values[1].(float64),
			IngredientNames: ingredientNames,
		}
	})
}
