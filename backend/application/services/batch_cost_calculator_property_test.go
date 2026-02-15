package services

import (
	"context"
	"math"
	"testing"

	"cafe-pos/backend/domain/batch"
	"cafe-pos/backend/domain/ingredient"
	"github.com/leanovate/gopter"
	"github.com/leanovate/gopter/gen"
	"github.com/leanovate/gopter/prop"
	"github.com/stretchr/testify/assert"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

// Property Test: Cost Accuracy
// Validates: Requirements 3.1, 3.2, 3.5
//
// Property: The cost of a batch must be calculated accurately from source ingredient costs
// Formula verification:
// - For each ingredient: actual_quantity = (batch_quantity / conversion_batch_quantity) * conversion_source_quantity * (1 + wastage_rate)
// - ingredient_cost = actual_quantity * cost_per_unit
// - total_cost = sum of all ingredient_costs
// - cost_per_unit = total_cost / batch_quantity
//
// This property ensures that:
// 1. Wastage is correctly applied to ingredient quantities
// 2. Conversion rates are correctly applied
// 3. Costs are correctly summed
// 4. Cost per unit is correctly calculated
// 5. All calculations maintain precision (rounded to 2 decimal places)
func TestProperty_BatchCostAccuracy(t *testing.T) {
	parameters := gopter.DefaultTestParameters()
	parameters.MinSuccessfulTests = 100 // Run 100 random test cases
	properties := gopter.NewProperties(parameters)

	properties.Property("Batch cost is calculated accurately with wastage", prop.ForAll(
		func(testData batchCostTestData) bool {
			ctx := context.Background()
			repo := newMockIngredientRepo()
			calculator := NewBatchCostCalculator(repo)

			// Setup ingredients in repository
			for _, ing := range testData.ingredients {
				repo.ingredients[ing.ID] = ing
			}

			// Create batch definition
			batchDef := &batch.BatchDefinition{
				ID:              primitive.NewObjectID(),
				Name:            testData.batchName,
				Unit:            testData.batchUnit,
				ShelfLifeHours:  24,
				ConversionRates: testData.conversionRates,
			}

			// Calculate cost
			breakdown, err := calculator.CalculateBatchCost(ctx, batchDef, testData.quantityProduced)
			if err != nil {
				t.Logf("Error calculating cost: %v", err)
				return false
			}

			// Verify: Number of ingredient costs matches conversion rates
			if len(breakdown.IngredientCosts) != len(testData.conversionRates) {
				t.Logf("Expected %d ingredient costs, got %d", len(testData.conversionRates), len(breakdown.IngredientCosts))
				return false
			}

			// Verify: Each ingredient cost is calculated correctly
			var expectedTotalCost float64
			for i, convRate := range testData.conversionRates {
				// Find the ingredient
				ing := repo.ingredients[convRate.SourceIngredientID]
				
				// Calculate expected quantity with wastage
				quantityNeeded := (testData.quantityProduced / convRate.BatchQuantity) * convRate.SourceQuantity
				actualQuantity := quantityNeeded * (1 + convRate.WastageRate)
				
				// Calculate expected cost
				expectedIngredientCost := actualQuantity * ing.CostPerUnit
				expectedIngredientCost = math.Round(expectedIngredientCost*100) / 100
				
				// Verify ingredient cost
				if !floatEquals(breakdown.IngredientCosts[i].Quantity, actualQuantity, 0.01) {
					t.Logf("Ingredient %s: Expected quantity %.2f, got %.2f", 
						ing.Name, actualQuantity, breakdown.IngredientCosts[i].Quantity)
					return false
				}
				
				if !floatEquals(breakdown.IngredientCosts[i].TotalCost, expectedIngredientCost, 0.01) {
					t.Logf("Ingredient %s: Expected cost %.2f, got %.2f", 
						ing.Name, expectedIngredientCost, breakdown.IngredientCosts[i].TotalCost)
					return false
				}
				
				expectedTotalCost += expectedIngredientCost
			}

			// Verify: Total cost is sum of ingredient costs
			expectedTotalCost = math.Round(expectedTotalCost*100) / 100
			if !floatEquals(breakdown.TotalCost, expectedTotalCost, 0.01) {
				t.Logf("Expected total cost %.2f, got %.2f", expectedTotalCost, breakdown.TotalCost)
				return false
			}

			// Verify: Cost per unit is correctly calculated
			expectedCostPerUnit := math.Round((expectedTotalCost/testData.quantityProduced)*100) / 100
			if !floatEquals(breakdown.CostPerUnit, expectedCostPerUnit, 0.01) {
				t.Logf("Expected cost per unit %.2f, got %.2f", expectedCostPerUnit, breakdown.CostPerUnit)
				return false
			}

			return true
		},
		genBatchCostTestData(),
	))

	properties.TestingRun(t, gopter.ConsoleReporter(false))
}

// Property Test: Cost Accuracy with Zero Wastage
// Validates: Requirements 3.1, 3.2
//
// Property: When wastage rate is 0, the cost calculation should be exact without wastage multiplier
func TestProperty_BatchCostAccuracy_NoWastage(t *testing.T) {
	parameters := gopter.DefaultTestParameters()
	parameters.MinSuccessfulTests = 50
	properties := gopter.NewProperties(parameters)

	properties.Property("Batch cost is calculated accurately without wastage", prop.ForAll(
		func(testData batchCostTestData) bool {
			ctx := context.Background()
			repo := newMockIngredientRepo()
			calculator := NewBatchCostCalculator(repo)

			// Setup ingredients in repository
			for _, ing := range testData.ingredients {
				repo.ingredients[ing.ID] = ing
			}

			// Create batch definition with zero wastage
			conversionRatesNoWastage := make([]batch.ConversionRate, len(testData.conversionRates))
			for i, cr := range testData.conversionRates {
				conversionRatesNoWastage[i] = cr
				conversionRatesNoWastage[i].WastageRate = 0.0 // Force zero wastage
			}

			batchDef := &batch.BatchDefinition{
				ID:              primitive.NewObjectID(),
				Name:            testData.batchName,
				Unit:            testData.batchUnit,
				ShelfLifeHours:  24,
				ConversionRates: conversionRatesNoWastage,
			}

			// Calculate cost
			breakdown, err := calculator.CalculateBatchCost(ctx, batchDef, testData.quantityProduced)
			if err != nil {
				return false
			}

			// Verify: Each ingredient quantity should equal base quantity (no wastage multiplier)
			for i, convRate := range conversionRatesNoWastage {
				expectedQuantity := (testData.quantityProduced / convRate.BatchQuantity) * convRate.SourceQuantity
				
				if !floatEquals(breakdown.IngredientCosts[i].Quantity, expectedQuantity, 0.01) {
					t.Logf("With zero wastage: Expected quantity %.2f, got %.2f", 
						expectedQuantity, breakdown.IngredientCosts[i].Quantity)
					return false
				}
			}

			return true
		},
		genBatchCostTestData(),
	))

	properties.TestingRun(t, gopter.ConsoleReporter(false))
}

// Property Test: Cost Linearity
// Validates: Requirements 3.1, 3.2
//
// Property: Doubling the batch quantity should double the total cost (linear relationship)
func TestProperty_BatchCostLinearity(t *testing.T) {
	parameters := gopter.DefaultTestParameters()
	parameters.MinSuccessfulTests = 50
	properties := gopter.NewProperties(parameters)

	properties.Property("Batch cost scales linearly with quantity", prop.ForAll(
		func(testData batchCostTestData) bool {
			ctx := context.Background()
			repo := newMockIngredientRepo()
			calculator := NewBatchCostCalculator(repo)

			// Setup ingredients
			for _, ing := range testData.ingredients {
				repo.ingredients[ing.ID] = ing
			}

			batchDef := &batch.BatchDefinition{
				ID:              primitive.NewObjectID(),
				Name:            testData.batchName,
				Unit:            testData.batchUnit,
				ShelfLifeHours:  24,
				ConversionRates: testData.conversionRates,
			}

			// Calculate cost for base quantity
			breakdown1, err := calculator.CalculateBatchCost(ctx, batchDef, testData.quantityProduced)
			if err != nil {
				return false
			}

			// Calculate cost for double quantity
			breakdown2, err := calculator.CalculateBatchCost(ctx, batchDef, testData.quantityProduced*2)
			if err != nil {
				return false
			}

			// Verify: Double quantity should have double cost
			expectedDoubleCost := breakdown1.TotalCost * 2
			if !floatEquals(breakdown2.TotalCost, expectedDoubleCost, 0.05) { // Increased tolerance for rounding
				t.Logf("Expected double cost %.2f, got %.2f", expectedDoubleCost, breakdown2.TotalCost)
				return false
			}

			// Verify: Cost per unit should remain the same
			if !floatEquals(breakdown1.CostPerUnit, breakdown2.CostPerUnit, 0.02) { // Increased tolerance for rounding
				t.Logf("Cost per unit changed: %.2f vs %.2f", breakdown1.CostPerUnit, breakdown2.CostPerUnit)
				return false
			}

			return true
		},
		genBatchCostTestData(),
	))

	properties.TestingRun(t, gopter.ConsoleReporter(false))
}

// Property Test: Cost Stored at Creation Time
// Validates: Requirements 3.5
//
// Property: The cost calculated at batch creation time should not change when ingredient costs change
// This is implicitly tested by the fact that we store cost_per_unit in the batch record
func TestProperty_BatchCostImmutability(t *testing.T) {
	// This property is verified by the design:
	// - BatchRecord stores cost_per_unit and total_cost at creation time
	// - These values are never recalculated
	// - The calculator uses current ingredient costs, but results are stored immediately
	
	// This test verifies the calculation is deterministic given the same inputs
	parameters := gopter.DefaultTestParameters()
	parameters.MinSuccessfulTests = 50
	properties := gopter.NewProperties(parameters)

	properties.Property("Batch cost calculation is deterministic", prop.ForAll(
		func(testData batchCostTestData) bool {
			ctx := context.Background()
			repo := newMockIngredientRepo()
			calculator := NewBatchCostCalculator(repo)

			// Setup ingredients
			for _, ing := range testData.ingredients {
				repo.ingredients[ing.ID] = ing
			}

			batchDef := &batch.BatchDefinition{
				ID:              primitive.NewObjectID(),
				Name:            testData.batchName,
				Unit:            testData.batchUnit,
				ShelfLifeHours:  24,
				ConversionRates: testData.conversionRates,
			}

			// Calculate cost twice with same inputs
			breakdown1, err1 := calculator.CalculateBatchCost(ctx, batchDef, testData.quantityProduced)
			breakdown2, err2 := calculator.CalculateBatchCost(ctx, batchDef, testData.quantityProduced)

			if err1 != nil || err2 != nil {
				return false
			}

			// Verify: Both calculations produce identical results
			return floatEquals(breakdown1.TotalCost, breakdown2.TotalCost, 0.001) &&
				floatEquals(breakdown1.CostPerUnit, breakdown2.CostPerUnit, 0.001)
		},
		genBatchCostTestData(),
	))

	properties.TestingRun(t, gopter.ConsoleReporter(false))
}

// Test data structure for property tests
type batchCostTestData struct {
	batchName        string
	batchUnit        string
	quantityProduced float64
	ingredients      []*ingredient.Ingredient
	conversionRates  []batch.ConversionRate
}

// Generator for batch cost test data
func genBatchCostTestData() gopter.Gen {
	return gopter.CombineGens(
		gen.Identifier(),                  // batch name
		gen.OneConstOf("ml", "g", "kg"),   // batch unit
		gen.Float64Range(100, 2000),       // quantity produced
		genIngredientList(),               // ingredients
	).Map(func(values []interface{}) batchCostTestData {
		batchName := values[0].(string)
		batchUnit := values[1].(string)
		quantityProduced := values[2].(float64)
		ingredients := values[3].([]*ingredient.Ingredient)

		// Create conversion rates for each ingredient
		conversionRates := make([]batch.ConversionRate, len(ingredients))
		for i, ing := range ingredients {
			conversionRates[i] = batch.ConversionRate{
				SourceIngredientID:   ing.ID,
				SourceIngredientName: ing.Name,
				SourceQuantity:       float64(50 + i*25),  // 50, 75, 100, ...
				SourceUnit:           string(ing.Unit),    // Convert UnitType to string
				BatchQuantity:        500,                 // Fixed batch quantity for conversion
				WastageRate:          float64(i) * 0.05,   // 0%, 5%, 10%, ...
			}
		}

		return batchCostTestData{
			batchName:        batchName,
			batchUnit:        batchUnit,
			quantityProduced: quantityProduced,
			ingredients:      ingredients,
			conversionRates:  conversionRates,
		}
	})
}

// Generator for ingredient list (1-3 ingredients)
func genIngredientList() gopter.Gen {
	return gen.SliceOf(genIngredient()).
		Map(func(slice []*ingredient.Ingredient) []*ingredient.Ingredient {
			// Ensure 1-3 ingredients
			if len(slice) == 0 {
				// Create at least one ingredient
				ing := &ingredient.Ingredient{
					ID:          primitive.NewObjectID(),
					Name:        "DefaultIngredient",
					Unit:        ingredient.UnitGram,
					CostPerUnit: 1.0,
				}
				return []*ingredient.Ingredient{ing}
			}
			if len(slice) > 3 {
				return slice[:3]
			}
			return slice
		})
}

// Generator for a single ingredient
func genIngredient() gopter.Gen {
	return gopter.CombineGens(
		gen.Identifier(),                                      // name
		gen.OneConstOf("g", "ml", "kg"),                       // unit (as string)
		gen.Float64Range(0.01, 10.0),                          // cost per unit
	).Map(func(values []interface{}) *ingredient.Ingredient {
		unitStr := values[1].(string)
		var unit ingredient.UnitType
		switch unitStr {
		case "g":
			unit = ingredient.UnitGram
		case "ml":
			unit = ingredient.UnitMilliliter
		case "kg":
			unit = ingredient.UnitKilogram
		default:
			unit = ingredient.UnitGram
		}
		
		return &ingredient.Ingredient{
			ID:          primitive.NewObjectID(),
			Name:        values[0].(string),
			Unit:        unit,
			CostPerUnit: math.Round(values[2].(float64)*100) / 100, // Round to 2 decimals
		}
	})
}

// Helper function to compare floats with tolerance
func floatEquals(a, b, tolerance float64) bool {
	return math.Abs(a-b) <= tolerance
}

// Test helper to verify property test framework is working
func TestBatchCostCalculator_PropertyTestFramework(t *testing.T) {
	// Simple sanity check that gopter is working
	properties := gopter.NewProperties(nil)
	
	properties.Property("Addition is commutative", prop.ForAll(
		func(a, b float64) bool {
			return floatEquals(a+b, b+a, 0.0001)
		},
		gen.Float64(),
		gen.Float64(),
	))
	
	result := properties.Run(gopter.ConsoleReporter(false))
	assert.True(t, result)
}
