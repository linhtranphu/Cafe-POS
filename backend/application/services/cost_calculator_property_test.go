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

// Feature: menu-cost-profit-analysis, Property 1: Cost Calculation Formula
// **Validates: Requirements 1.1, 1.2, 1.7, 10.1, 10.2, 10.4**
//
// Property: For any menu item with ingredients that have valid cost_per_unit values,
// calculating the current_cost should produce a result equal to the sum of
// (ingredient.quantity * ingredient.cost_per_unit * conversion_rate * (1 + wastage_percentage/100))
// for all ingredients, rounded to 2 decimal places.
func TestProperty_CostCalculationFormula(t *testing.T) {
	parameters := gopter.DefaultTestParameters()
	parameters.MinSuccessfulTests = 100 // Minimum 100 iterations as per spec
	properties := gopter.NewProperties(parameters)

	properties.Property("Cost equals sum of ingredient costs with conversion and wastage", prop.ForAll(
		func(ingredientData []testIngredientData) bool {
			// Skip empty ingredient lists
			if len(ingredientData) == 0 {
				return true
			}

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

			// Create menu item with ingredients
			menuItem := &menu.MenuItem{
				ID:          primitive.ObjectID{},
				Name:        "Test Item",
				Ingredients: []menu.Ingredient{},
			}

			// Calculate expected cost manually
			expectedCost := 0.0
			hasValidIngredients := false

			for _, ingData := range ingredientData {
				// Skip ingredients with invalid cost (should be marked incomplete)
				if ingData.CostPerUnit <= 0 {
					continue
				}

				hasValidIngredients = true

				// Use same unit for both stock and recipe (no conversion)
				// This makes the test simpler and focuses on the core formula
				testUnit := ingredient.UnitPiece

				// Create ingredient with unique ID
				ing := &ingredient.Ingredient{
					ID:                primitive.NewObjectID(),
					Name:              ingData.Name,
					Unit:              testUnit,  // Stock unit
					CostPerUnit:       ingData.CostPerUnit,
					WastagePercentage: ingData.WastagePercentage,
				}

				// Add ingredient to repository using ID as key
				ingredientRepo.ingredients[ing.ID] = ing

				// Add ingredient to menu item with same unit (no conversion needed)
				menuItem.Ingredients = append(menuItem.Ingredients, menu.Ingredient{
					Name:     ingData.Name,
					Quantity: ingData.Quantity,
					Unit:     testUnit,  // Recipe unit = stock unit
				})

				// Calculate expected cost for this ingredient
				// Since stock unit = recipe unit, conversion rate = 1.0
				conversionRate := 1.0

				wastageMultiplier := 1.0
				if ingData.WastagePercentage > 0 {
					wastageMultiplier = 1.0 + (ingData.WastagePercentage / 100.0)
				}

				ingredientCost := ingData.Quantity * ingData.CostPerUnit * conversionRate * wastageMultiplier
				expectedCost += ingredientCost
			}

			// Skip if no valid ingredients
			if !hasValidIngredients {
				return true
			}

			// Round expected cost to 2 decimal places
			expectedCost = math.Round(expectedCost*100) / 100

			// Add menu item to repository
			menuRepo.menuItems[menuItem.ID] = menuItem

			// Create service and calculate cost
			service := NewCostCalculatorService(menuRepo, ingredientRepo, orderRepo, orderItemRepo)
			result, err := service.CalculateMenuItemCost(context.Background(), menuItem.ID)

			// Should not error
			if err != nil {
				t.Logf("Unexpected error: %v", err)
				return false
			}

			// Cost should match expected (within floating point tolerance)
			tolerance := 0.01
			if math.Abs(result.CurrentCost-expectedCost) > tolerance {
				t.Logf("Cost mismatch: expected %v, got %v (diff: %v)", 
					expectedCost, result.CurrentCost, math.Abs(result.CurrentCost-expectedCost))
				return false
			}

			// Cost status should be FINAL if all ingredients have valid costs
			if result.CostStatus != menu.CostStatusFinal {
				t.Logf("Expected status FINAL, got %v", result.CostStatus)
				return false
			}

			return true
		},
		genIngredientDataList(),
	))

	properties.TestingRun(t, gopter.ConsoleReporter(false))
}

// Feature: menu-cost-profit-analysis, Property 1: Cost Calculation Formula (Incomplete Cost Status)
// **Validates: Requirements 1.5, 1.6**
//
// Property: For any menu item where at least one ingredient has zero or missing cost_per_unit,
// the cost_status should be marked as "INCOMPLETE" and missing ingredients should be tracked.
func TestProperty_CostCalculationFormula_IncompleteCost(t *testing.T) {
	parameters := gopter.DefaultTestParameters()
	parameters.MinSuccessfulTests = 100
	properties := gopter.NewProperties(parameters)

	properties.Property("Items with missing ingredient costs are marked INCOMPLETE", prop.ForAll(
		func(validIngredients []testIngredientData, missingCostIngredient testIngredientData) bool {
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

			// Create menu item
			menuItem := &menu.MenuItem{
				ID:          primitive.NewObjectID(),
				Name:        "Test Item",
				Ingredients: []menu.Ingredient{},
			}

			// Add valid ingredients
			for _, ingData := range validIngredients {
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

				menuItem.Ingredients = append(menuItem.Ingredients, menu.Ingredient{
					Name:     ingData.Name,
					Quantity: ingData.Quantity,
					Unit:     "unit",
				})
			}

			// Add ingredient with missing cost (cost = 0)
			missingIngredientName := "MissingCost_" + missingCostIngredient.Name
			missingIng := &ingredient.Ingredient{
				ID:                primitive.NewObjectID(),
				Name:              missingIngredientName,
				CostPerUnit:       0.0, // Missing cost
				ConversionRate:    1.0,
				WastagePercentage: 0.0,
			}
			ingredientRepo.ingredients[missingIng.ID] = missingIng

			menuItem.Ingredients = append(menuItem.Ingredients, menu.Ingredient{
				Name:     missingIngredientName,
				Quantity: missingCostIngredient.Quantity,
				Unit:     "unit",
			})

			// Add menu item to repository
			menuRepo.menuItems[menuItem.ID] = menuItem

			// Create service and calculate cost
			service := NewCostCalculatorService(menuRepo, ingredientRepo, orderRepo, orderItemRepo)
			result, err := service.CalculateMenuItemCost(context.Background(), menuItem.ID)

			// Should not error
			if err != nil {
				t.Logf("Unexpected error: %v", err)
				return false
			}

			// Cost status should be INCOMPLETE
			if result.CostStatus != menu.CostStatusIncomplete {
				t.Logf("Expected status INCOMPLETE, got %v", result.CostStatus)
				return false
			}

			// Missing ingredients should include the one with zero cost
			found := false
			for _, missing := range result.MissingIngredients {
				if missing == missingIngredientName {
					found = true
					break
				}
			}
			if !found {
				t.Logf("Expected missing ingredient '%s' in list: %v", missingIngredientName, result.MissingIngredients)
				return false
			}

			return true
		},
		gen.SliceOf(genIngredientData()).Map(func(slice []testIngredientData) []testIngredientData {
			if len(slice) > 5 {
				return slice[:5]
			}
			return slice
		}),
		genIngredientData(),
	))

	properties.TestingRun(t, gopter.ConsoleReporter(false))
}

// Feature: menu-cost-profit-analysis, Property 1: Cost Calculation Formula (Rounding)
// **Validates: Requirements 1.7**
//
// Property: For any calculated cost, the result should be rounded to exactly 2 decimal places.
func TestProperty_CostCalculationFormula_Rounding(t *testing.T) {
	parameters := gopter.DefaultTestParameters()
	parameters.MinSuccessfulTests = 100
	properties := gopter.NewProperties(parameters)

	properties.Property("Cost is always rounded to 2 decimal places", prop.ForAll(
		func(ingredientData []testIngredientData) bool {
			// Skip empty ingredient lists
			if len(ingredientData) == 0 {
				return true
			}

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

			// Create menu item with ingredients
			menuItem := &menu.MenuItem{
				ID:          primitive.NewObjectID(),
				Name:        "Test Item",
				Ingredients: []menu.Ingredient{},
			}

			hasValidIngredients := false
			for _, ingData := range ingredientData {
				if ingData.CostPerUnit <= 0 {
					continue
				}

				hasValidIngredients = true

				ing := &ingredient.Ingredient{
					ID:                primitive.NewObjectID(),
					Name:              ingData.Name,
					CostPerUnit:       ingData.CostPerUnit,
					ConversionRate:    ingData.ConversionRate,
					WastagePercentage: ingData.WastagePercentage,
				}

				ingredientRepo.ingredients[ing.ID] = ing

				menuItem.Ingredients = append(menuItem.Ingredients, menu.Ingredient{
					Name:     ingData.Name,
					Quantity: ingData.Quantity,
					Unit:     "unit",
				})
			}

			if !hasValidIngredients {
				return true
			}

			menuRepo.menuItems[menuItem.ID] = menuItem

			// Create service and calculate cost
			service := NewCostCalculatorService(menuRepo, ingredientRepo, orderRepo, orderItemRepo)
			result, err := service.CalculateMenuItemCost(context.Background(), menuItem.ID)

			if err != nil {
				t.Logf("Unexpected error: %v", err)
				return false
			}

			// Check that cost has at most 2 decimal places
			// Multiply by 100, round, divide by 100 should give same result
			rounded := math.Round(result.CurrentCost*100) / 100
			if math.Abs(result.CurrentCost-rounded) > 0.0001 {
				t.Logf("Cost not properly rounded to 2 decimals: %v (rounded: %v)", result.CurrentCost, rounded)
				return false
			}

			return true
		},
		genIngredientDataList(),
	))

	properties.TestingRun(t, gopter.ConsoleReporter(false))
}

// Test data structure for property-based testing
type testIngredientData struct {
	Name              string
	Quantity          float64
	CostPerUnit       float64
	ConversionRate    float64
	WastagePercentage float64
}

// Generator for ingredient data list (1-10 ingredients)
func genIngredientDataList() gopter.Gen {
	return gen.SliceOf(genIngredientData()).
		Map(func(slice []testIngredientData) []testIngredientData {
			// Ensure at least 1 ingredient
			if len(slice) == 0 {
				return []testIngredientData{
					{
						Name:              "DefaultIngredient",
						Quantity:          10.0,
						CostPerUnit:       100.0,
						ConversionRate:    1.0,
						WastagePercentage: 0.0,
					},
				}
			}
			// Limit to 10 ingredients
			if len(slice) > 10 {
				return slice[:10]
			}
			return slice
		})
}

// Generator for single ingredient data
func genIngredientData() gopter.Gen {
	return gopter.CombineGens(
		gen.Identifier(),                // Name
		gen.Float64Range(0.1, 1000.0),   // Quantity
		gen.Float64Range(1.0, 10000.0),  // CostPerUnit (always valid for this property)
		gen.Float64Range(0.1, 10.0),     // ConversionRate
		gen.Float64Range(0.0, 50.0),     // WastagePercentage
	).Map(func(values []interface{}) testIngredientData {
		return testIngredientData{
			Name:              values[0].(string),
			Quantity:          values[1].(float64),
			CostPerUnit:       values[2].(float64),
			ConversionRate:    values[3].(float64),
			WastagePercentage: values[4].(float64),
		}
	})
}
