package services

import (
	"context"
	"testing"

	"cafe-pos/backend/domain/ingredient"
	"cafe-pos/backend/domain/menu"
	"cafe-pos/backend/domain/settings"
	"github.com/leanovate/gopter"
	"github.com/leanovate/gopter/gen"
	"github.com/leanovate/gopter/prop"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

// Feature: menu-cost-profit-analysis, Property 19: Summary Statistics Calculation
// Validates: Requirements 7.4
//
// Property: For any set of menu items, the summary statistics should correctly calculate:
// - total_items = count of all items
// - loss_count = count of items with cost > price
// - low_margin_count = count of items with profit_margin < threshold
// - average_profit_margin = average of all profit_margins (excluding items with price <= 0)
func TestProperty_SummaryStatisticsCalculation(t *testing.T) {
	properties := gopter.NewProperties(nil)

	properties.Property("Summary statistics are calculated correctly", prop.ForAll(
		func(prices []float64, costs []float64, threshold float64) bool {
			// Skip if arrays are empty or different lengths
			if len(prices) == 0 || len(costs) == 0 || len(prices) != len(costs) {
				return true
			}

			ctx := context.Background()

			// Setup mock repositories
			menuRepo := &mockMenuRepository{
				menuItems: make(map[primitive.ObjectID]*menu.MenuItem),
			}
			ingredientRepo := &mockIngredientRepository{
				ingredients: make(map[primitive.ObjectID]*ingredient.Ingredient),
			}
			settingsRepo := &mockSettingsRepo{
				settings: &settings.ShopSettings{
					ID:                 primitive.NewObjectID(),
					LowMarginThreshold: threshold,
				},
			}

			// Create service
			costCalculator := NewCostCalculatorService(menuRepo, ingredientRepo, nil, nil)
			profitAnalyzer := NewProfitAnalyzerService(menuRepo, nil, settingsRepo, nil)

			// Create a single ingredient for all menu items
			testIngredient := &ingredient.Ingredient{
				ID:                primitive.NewObjectID(),
				Name:              "Test Ingredient",
				CostPerUnit:       1.0, // We'll adjust quantities to match desired costs
				ConversionRate:    1.0,
				WastagePercentage: 0.0,
			}
			ingredientRepo.ingredients[testIngredient.ID] = testIngredient

			// Calculate expected statistics manually
			expectedTotalItems := len(prices)
			expectedLossCount := 0
			expectedLowMarginCount := 0
			totalProfitMargin := 0.0
			profitMarginCount := 0

			// Create menu items with specified prices and costs
			for i := 0; i < len(prices); i++ {
				price := prices[i]
				cost := costs[i]

				// Create menu item
				menuItem := &menu.MenuItem{
					ID:    primitive.NewObjectID(),
					Name:  "Test Item",
					Price: price,
					Ingredients: []menu.Ingredient{
						{
							Name:     testIngredient.Name,
							Quantity: cost, // Quantity = cost since cost_per_unit = 1.0
							Unit:     "g",
						},
					},
				}
				menuRepo.menuItems[menuItem.ID] = menuItem

				// Calculate cost
				costResult, err := costCalculator.CalculateMenuItemCost(ctx, menuItem.ID)
				if err != nil {
					return false
				}

				// Update menu item with calculated cost
				menuItem.CurrentCost = costResult.CurrentCost
				menuItem.CostStatus = costResult.CostStatus
				menuItem.CostLastCalculatedAt = costResult.CostLastCalculatedAt
				menuRepo.menuItems[menuItem.ID] = menuItem

				// Calculate expected statistics
				if price <= 0 {
					// Promotional items (price <= 0) are not counted as loss or low margin
					// They are excluded from profit margin calculations
					// But they ARE counted in total items
				} else if cost > price {
					expectedLossCount++
					// Loss items ARE included in average profit margin calculation
					profitMargin := ((price - cost) / price) * 100
					totalProfitMargin += profitMargin
					profitMarginCount++
				} else {
					profitMargin := ((price - cost) / price) * 100
					if profitMargin < threshold {
						expectedLowMarginCount++
					}
					totalProfitMargin += profitMargin
					profitMarginCount++
				}
			}

			expectedAverageProfitMargin := 0.0
			if profitMarginCount > 0 {
				expectedAverageProfitMargin = totalProfitMargin / float64(profitMarginCount)
			}

			// Get actual statistics from service
			filter := ProfitFilter{}
			response, err := profitAnalyzer.GetAllMenuItemProfits(ctx, filter)
			if err != nil {
				return false
			}

			actualTotalItems := response.Summary.TotalItems
			actualLossCount := response.Summary.LossCount
			actualLowMarginCount := response.Summary.LowMarginCount
			actualAverageProfitMargin := response.Summary.AverageProfitMargin

			// Verify statistics match
			if actualTotalItems != expectedTotalItems {
				t.Logf("FAIL: Total items mismatch: expected %d, got %d",
					expectedTotalItems, actualTotalItems)
				return false
			}

			if actualLossCount != expectedLossCount {
				t.Logf("FAIL: Loss count mismatch: expected %d, got %d",
					expectedLossCount, actualLossCount)
				return false
			}

			if actualLowMarginCount != expectedLowMarginCount {
				t.Logf("FAIL: Low margin count mismatch: expected %d, got %d",
					expectedLowMarginCount, actualLowMarginCount)
				return false
			}

			// Allow small tolerance for floating point comparison
			if profitMarginCount > 0 {
				diff := actualAverageProfitMargin - expectedAverageProfitMargin
				if diff < -0.1 || diff > 0.1 {
					t.Logf("FAIL: Average profit margin mismatch: expected %.2f, got %.2f",
						expectedAverageProfitMargin, actualAverageProfitMargin)
					return false
				}
			}

			return true
		},
		gen.SliceOfN(5, gen.Float64Range(0, 10000)),   // prices
		gen.SliceOfN(5, gen.Float64Range(0, 10000)),   // costs
		gen.Float64Range(10, 30),                       // threshold
	))

	properties.TestingRun(t, gopter.ConsoleReporter(false))
}

// Feature: menu-cost-profit-analysis, Property 19: Summary Statistics Calculation (Edge Cases)
// Test with specific edge cases
func TestProperty_SummaryStatistics_EdgeCases(t *testing.T) {
	testCases := []struct {
		name      string
		items     []struct{ price, cost float64 }
		threshold float64
		expected  struct {
			totalItems          int
			lossCount           int
			lowMarginCount      int
			averageProfitMargin float64
		}
	}{
		{
			name: "All profitable items",
			items: []struct{ price, cost float64 }{
				{10000, 3000}, // 70% margin
				{20000, 5000}, // 75% margin
				{15000, 3000}, // 80% margin
			},
			threshold: 20.0,
			expected: struct {
				totalItems          int
				lossCount           int
				lowMarginCount      int
				averageProfitMargin float64
			}{
				totalItems:          3,
				lossCount:           0,
				lowMarginCount:      0,
				averageProfitMargin: 75.0,
			},
		},
		{
			name: "All loss items",
			items: []struct{ price, cost float64 }{
				{10000, 15000}, // -50% margin
				{20000, 25000}, // -25% margin
			},
			threshold: 20.0,
			expected: struct {
				totalItems          int
				lossCount           int
				lowMarginCount      int
				averageProfitMargin float64
			}{
				totalItems:          2,
				lossCount:           2,
				lowMarginCount:      0,
				averageProfitMargin: -37.5,
			},
		},
		{
			name: "Mixed items",
			items: []struct{ price, cost float64 }{
				{10000, 3000},  // 70% margin - profitable
				{15000, 13500}, // 10% margin - low margin
				{12000, 16000}, // -33.33% margin - loss
				{6000, 6000},   // 0% margin - low margin
			},
			threshold: 20.0,
			expected: struct {
				totalItems          int
				lossCount           int
				lowMarginCount      int
				averageProfitMargin float64
			}{
				totalItems:          4,
				lossCount:           1,
				lowMarginCount:      2,
				averageProfitMargin: 11.67, // (70 + 10 + (-33.33) + 0) / 4
			},
		},
		{
			name: "Items with zero price (promotional)",
			items: []struct{ price, cost float64 }{
				{10000, 3000}, // 70% margin
				{0, 5000},     // Promotional item - excluded from average
				{20000, 5000}, // 75% margin
			},
			threshold: 20.0,
			expected: struct {
				totalItems          int
				lossCount           int
				lowMarginCount      int
				averageProfitMargin float64
			}{
				totalItems:          3,
				lossCount:           0,
				lowMarginCount:      0,
				averageProfitMargin: 72.5, // (70 + 75) / 2, excluding zero price item
			},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			ctx := context.Background()

			// Setup mock repositories
			menuRepo := &mockMenuRepository{
				menuItems: make(map[primitive.ObjectID]*menu.MenuItem),
			}
			ingredientRepo := &mockIngredientRepository{
				ingredients: make(map[primitive.ObjectID]*ingredient.Ingredient),
			}
			settingsRepo := &mockSettingsRepo{
				settings: &settings.ShopSettings{
					ID:                 primitive.NewObjectID(),
					LowMarginThreshold: tc.threshold,
				},
			}

			// Create services
			costCalculator := NewCostCalculatorService(menuRepo, ingredientRepo, nil, nil)
			profitAnalyzer := NewProfitAnalyzerService(menuRepo, nil, settingsRepo, nil)

			// Create test ingredient
			testIngredient := &ingredient.Ingredient{
				ID:                primitive.NewObjectID(),
				Name:              "Test Ingredient",
				CostPerUnit:       1.0,
				ConversionRate:    1.0,
				WastagePercentage: 0.0,
			}
			ingredientRepo.ingredients[testIngredient.ID] = testIngredient

			// Create menu items
			for i, item := range tc.items {
				menuItem := &menu.MenuItem{
					ID:    primitive.NewObjectID(),
					Name:  "Test Item " + string(rune(i)),
					Price: item.price,
					Ingredients: []menu.Ingredient{
						{
							Name:     testIngredient.Name,
							Quantity: item.cost,
							Unit:     "g",
						},
					},
				}
				menuRepo.menuItems[menuItem.ID] = menuItem

				// Calculate cost
				costResult, err := costCalculator.CalculateMenuItemCost(ctx, menuItem.ID)
				if err != nil {
					t.Fatalf("Failed to calculate cost: %v", err)
				}

				// Update menu item
				menuItem.CurrentCost = costResult.CurrentCost
				menuItem.CostStatus = costResult.CostStatus
				menuItem.CostLastCalculatedAt = costResult.CostLastCalculatedAt
				menuRepo.menuItems[menuItem.ID] = menuItem
			}

			// Get statistics
			filter := ProfitFilter{}
			response, err := profitAnalyzer.GetAllMenuItemProfits(ctx, filter)
			if err != nil {
				t.Fatalf("Failed to get menu item profits: %v", err)
			}

			// Verify statistics
			if response.Summary.TotalItems != tc.expected.totalItems {
				t.Errorf("Total items: expected %d, got %d",
					tc.expected.totalItems, response.Summary.TotalItems)
			}

			if response.Summary.LossCount != tc.expected.lossCount {
				t.Errorf("Loss count: expected %d, got %d",
					tc.expected.lossCount, response.Summary.LossCount)
			}

			if response.Summary.LowMarginCount != tc.expected.lowMarginCount {
				t.Errorf("Low margin count: expected %d, got %d",
					tc.expected.lowMarginCount, response.Summary.LowMarginCount)
			}

			// Allow small tolerance for average profit margin
			diff := response.Summary.AverageProfitMargin - tc.expected.averageProfitMargin
			if diff < -1.0 || diff > 1.0 {
				t.Errorf("Average profit margin: expected %.2f, got %.2f",
					tc.expected.averageProfitMargin, response.Summary.AverageProfitMargin)
			}

			t.Logf("✓ %s: Total=%d, Loss=%d, LowMargin=%d, AvgMargin=%.2f%%",
				tc.name,
				response.Summary.TotalItems,
				response.Summary.LossCount,
				response.Summary.LowMarginCount,
				response.Summary.AverageProfitMargin)
		})
	}

	t.Log("✅ Summary statistics edge cases test passed!")
}
