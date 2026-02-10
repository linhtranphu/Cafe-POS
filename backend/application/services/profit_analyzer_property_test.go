package services

import (
	"context"
	"fmt"
	"math"
	"testing"
	"time"

	"cafe-pos/backend/domain/expense"
	"cafe-pos/backend/domain/menu"
	"cafe-pos/backend/domain/order"
	"cafe-pos/backend/domain/settings"
	"github.com/leanovate/gopter"
	"github.com/leanovate/gopter/gen"
	"github.com/leanovate/gopter/prop"
	"go.mongodb.org/mongo-driver/bson/primitive"
)


// Feature: menu-cost-profit-analysis, Property 2: Profit Calculations
// **Validates: Requirements 2.1, 2.5, 2.6**
//
// Property: For any menu item with valid cost and price values, the profit_margin should equal
// ((price - cost) / price) * 100 rounded to 2 decimal places, and the absolute_profit should
// equal (price - cost).
func TestProperty_ProfitCalculations(t *testing.T) {
	parameters := gopter.DefaultTestParameters()
	parameters.MinSuccessfulTests = 100 // Minimum 100 iterations as per spec
	properties := gopter.NewProperties(parameters)

	properties.Property("Profit margin and absolute profit are correctly calculated", prop.ForAll(
		func(price, cost float64) bool {
			// Skip invalid cases where price <= 0 (promotional items handled separately)
			if price <= 0 {
				return true
			}

			// Setup test repositories
			menuRepo := &mockMenuRepo{
				items: []*menu.MenuItem{
					{
						ID:          primitive.NewObjectID(),
						Name:        "Test Item",
						Category:    "Test",
						Price:       price,
						CurrentCost: cost,
						CostStatus:  menu.CostStatusFinal,
						CostLastCalculatedAt: time.Now(),
					},
				},
			}

			// Create service
			service := NewProfitAnalyzerService(menuRepo, nil, nil, nil)

			// Calculate profit metrics
			result, err := service.CalculateMenuItemProfit(context.Background(), menuRepo.items[0].ID)

			// Should not error
			if err != nil {
				t.Logf("Unexpected error: %v", err)
				return false
			}

			// Calculate expected profit margin: ((price - cost) / price) * 100
			expectedProfitMargin := ((price - cost) / price) * 100
			// Round to 2 decimal places
			expectedProfitMargin = math.Round(expectedProfitMargin*100) / 100

			// Calculate expected absolute profit: price - cost
			expectedAbsoluteProfit := price - cost

			// Verify profit margin (with floating point tolerance)
			tolerance := 0.01
			if math.Abs(result.ProfitMargin-expectedProfitMargin) > tolerance {
				t.Logf("Profit margin mismatch: expected %v, got %v (price=%v, cost=%v)",
					expectedProfitMargin, result.ProfitMargin, price, cost)
				return false
			}

			// Verify absolute profit (with floating point tolerance)
			if math.Abs(result.AbsoluteProfit-expectedAbsoluteProfit) > tolerance {
				t.Logf("Absolute profit mismatch: expected %v, got %v (price=%v, cost=%v)",
					expectedAbsoluteProfit, result.AbsoluteProfit, price, cost)
				return false
			}

			return true
		},
		gen.Float64Range(1.0, 100000.0),  // price (always > 0)
		gen.Float64Range(0.0, 100000.0),  // cost (can be 0 or higher than price)
	))

	properties.TestingRun(t, gopter.ConsoleReporter(false))
}

// Feature: menu-cost-profit-analysis, Property 2: Profit Calculations (Rounding)
// **Validates: Requirements 2.5**
//
// Property: For any calculated profit margin, the result should be rounded to exactly 2 decimal places.
func TestProperty_ProfitCalculations_Rounding(t *testing.T) {
	parameters := gopter.DefaultTestParameters()
	parameters.MinSuccessfulTests = 100
	properties := gopter.NewProperties(parameters)

	properties.Property("Profit margin is always rounded to 2 decimal places", prop.ForAll(
		func(price, cost float64) bool {
			// Skip invalid cases where price <= 0
			if price <= 0 {
				return true
			}

			// Setup test repositories
			menuRepo := &mockMenuRepo{
				items: []*menu.MenuItem{
					{
						ID:          primitive.NewObjectID(),
						Name:        "Test Item",
						Category:    "Test",
						Price:       price,
						CurrentCost: cost,
						CostStatus:  menu.CostStatusFinal,
						CostLastCalculatedAt: time.Now(),
					},
				},
			}

			// Create service
			service := NewProfitAnalyzerService(menuRepo, nil, nil, nil)

			// Calculate profit metrics
			result, err := service.CalculateMenuItemProfit(context.Background(), menuRepo.items[0].ID)

			if err != nil {
				t.Logf("Unexpected error: %v", err)
				return false
			}

			// Check that profit margin has at most 2 decimal places
			// Multiply by 100, round, divide by 100 should give same result
			rounded := math.Round(result.ProfitMargin*100) / 100
			if math.Abs(result.ProfitMargin-rounded) > 0.0001 {
				t.Logf("Profit margin not properly rounded to 2 decimals: %v (rounded: %v)",
					result.ProfitMargin, rounded)
				return false
			}

			return true
		},
		gen.Float64Range(1.0, 100000.0),
		gen.Float64Range(0.0, 100000.0),
	))

	properties.TestingRun(t, gopter.ConsoleReporter(false))
}

// Feature: menu-cost-profit-analysis, Property 2: Profit Calculations (Negative Profit)
// **Validates: Requirements 2.3, 2.6, 2.7**
//
// Property: For any menu item where cost exceeds price, the profit_margin should be negative
// and the absolute_profit should be negative.
func TestProperty_ProfitCalculations_NegativeProfit(t *testing.T) {
	parameters := gopter.DefaultTestParameters()
	parameters.MinSuccessfulTests = 100
	properties := gopter.NewProperties(parameters)

	properties.Property("Items with cost > price have negative profit", prop.ForAll(
		func(price, costExcess float64) bool {
			// Ensure cost > price by adding excess to price
			cost := price + costExcess

			// Skip invalid cases where price <= 0
			if price <= 0 {
				return true
			}

			// Setup test repositories
			menuRepo := &mockMenuRepo{
				items: []*menu.MenuItem{
					{
						ID:          primitive.NewObjectID(),
						Name:        "Loss Item",
						Category:    "Test",
						Price:       price,
						CurrentCost: cost,
						CostStatus:  menu.CostStatusFinal,
						CostLastCalculatedAt: time.Now(),
					},
				},
			}

			// Create service
			service := NewProfitAnalyzerService(menuRepo, nil, nil, nil)

			// Calculate profit metrics
			result, err := service.CalculateMenuItemProfit(context.Background(), menuRepo.items[0].ID)

			if err != nil {
				t.Logf("Unexpected error: %v", err)
				return false
			}

			// Profit margin should be negative
			if result.ProfitMargin >= 0 {
				t.Logf("Expected negative profit margin for cost > price, got %v (price=%v, cost=%v)",
					result.ProfitMargin, price, cost)
				return false
			}

			// Absolute profit should be negative
			if result.AbsoluteProfit >= 0 {
				t.Logf("Expected negative absolute profit for cost > price, got %v (price=%v, cost=%v)",
					result.AbsoluteProfit, price, cost)
				return false
			}

			// Absolute profit should equal price - cost
			expectedAbsoluteProfit := price - cost
			tolerance := 0.01
			if math.Abs(result.AbsoluteProfit-expectedAbsoluteProfit) > tolerance {
				t.Logf("Absolute profit mismatch: expected %v, got %v",
					expectedAbsoluteProfit, result.AbsoluteProfit)
				return false
			}

			return true
		},
		gen.Float64Range(1.0, 100000.0),   // price
		gen.Float64Range(0.1, 50000.0),    // costExcess (amount by which cost exceeds price)
	))

	properties.TestingRun(t, gopter.ConsoleReporter(false))
}

// Feature: menu-cost-profit-analysis, Property 2: Profit Calculations (Zero Price Edge Case)
// **Validates: Requirements 2.4**
//
// Property: For any menu item with price = 0 or negative (promotional/gifted items),
// the profit_margin should be 0 and absolute_profit should be -cost.
func TestProperty_ProfitCalculations_ZeroPrice(t *testing.T) {
	parameters := gopter.DefaultTestParameters()
	parameters.MinSuccessfulTests = 100
	properties := gopter.NewProperties(parameters)

	properties.Property("Items with price <= 0 have profit_margin = 0 and absolute_profit = -cost", prop.ForAll(
		func(cost float64) bool {
			// Test with price = 0 (promotional item)
			price := 0.0

			// Setup test repositories
			menuRepo := &mockMenuRepo{
				items: []*menu.MenuItem{
					{
						ID:          primitive.NewObjectID(),
						Name:        "Free Item",
						Category:    "Promo",
						Price:       price,
						CurrentCost: cost,
						CostStatus:  menu.CostStatusFinal,
						CostLastCalculatedAt: time.Now(),
					},
				},
			}

			// Create service
			service := NewProfitAnalyzerService(menuRepo, nil, nil, nil)

			// Calculate profit metrics
			result, err := service.CalculateMenuItemProfit(context.Background(), menuRepo.items[0].ID)

			if err != nil {
				t.Logf("Unexpected error: %v", err)
				return false
			}

			// Profit margin should be 0 for promotional items
			if result.ProfitMargin != 0 {
				t.Logf("Expected profit margin 0 for price=0, got %v", result.ProfitMargin)
				return false
			}

			// Absolute profit should be -cost
			expectedAbsoluteProfit := -cost
			tolerance := 0.01
			if math.Abs(result.AbsoluteProfit-expectedAbsoluteProfit) > tolerance {
				t.Logf("Expected absolute profit %v for price=0, got %v",
					expectedAbsoluteProfit, result.AbsoluteProfit)
				return false
			}

			return true
		},
		gen.Float64Range(0.0, 100000.0), // cost
	))

	properties.TestingRun(t, gopter.ConsoleReporter(false))
}

// Feature: menu-cost-profit-analysis, Property 2: Profit Calculations (Break-Even)
// **Validates: Requirements 2.2, 2.8**
//
// Property: For any menu item where cost equals price, the profit_margin should be 0
// and the absolute_profit should be 0.
func TestProperty_ProfitCalculations_BreakEven(t *testing.T) {
	parameters := gopter.DefaultTestParameters()
	parameters.MinSuccessfulTests = 100
	properties := gopter.NewProperties(parameters)

	properties.Property("Items with cost = price have zero profit", prop.ForAll(
		func(price float64) bool {
			// Skip invalid cases where price <= 0
			if price <= 0 {
				return true
			}

			// Set cost equal to price (break-even)
			cost := price

			// Setup test repositories
			menuRepo := &mockMenuRepo{
				items: []*menu.MenuItem{
					{
						ID:          primitive.NewObjectID(),
						Name:        "Break Even Item",
						Category:    "Test",
						Price:       price,
						CurrentCost: cost,
						CostStatus:  menu.CostStatusFinal,
						CostLastCalculatedAt: time.Now(),
					},
				},
			}

			// Create service
			service := NewProfitAnalyzerService(menuRepo, nil, nil, nil)

			// Calculate profit metrics
			result, err := service.CalculateMenuItemProfit(context.Background(), menuRepo.items[0].ID)

			if err != nil {
				t.Logf("Unexpected error: %v", err)
				return false
			}

			// Profit margin should be 0
			tolerance := 0.01
			if math.Abs(result.ProfitMargin) > tolerance {
				t.Logf("Expected profit margin 0 for break-even, got %v (price=%v, cost=%v)",
					result.ProfitMargin, price, cost)
				return false
			}

			// Absolute profit should be 0
			if math.Abs(result.AbsoluteProfit) > tolerance {
				t.Logf("Expected absolute profit 0 for break-even, got %v (price=%v, cost=%v)",
					result.AbsoluteProfit, price, cost)
				return false
			}

			return true
		},
		gen.Float64Range(1.0, 100000.0), // price (and cost)
	))

	properties.TestingRun(t, gopter.ConsoleReporter(false))
}

// Feature: menu-cost-profit-analysis, Property 10: Loss Detection
// **Validates: Requirements 3.1**
//
// Property: For any menu item where cost exceeds price, the warning_status should be marked as "loss".
func TestProperty_LossDetection(t *testing.T) {
	parameters := gopter.DefaultTestParameters()
	parameters.MinSuccessfulTests = 100 // Minimum 100 iterations as per spec
	properties := gopter.NewProperties(parameters)

	properties.Property("Items with cost > price are marked as loss", prop.ForAll(
		func(price, costExcess float64) bool {
			// Ensure cost > price by adding excess to price
			cost := price + costExcess

			// Skip invalid cases where price <= 0
			if price <= 0 {
				return true
			}

			// Setup test repositories
			menuRepo := &mockMenuRepo{
				items: []*menu.MenuItem{
					{
						ID:          primitive.NewObjectID(),
						Name:        "Loss Item",
						Category:    "Test",
						Price:       price,
						CurrentCost: cost,
						CostStatus:  menu.CostStatusFinal,
						CostLastCalculatedAt: time.Now(),
					},
				},
			}

			// Create service with mock settings repo
			settingsRepo := &mockSettingsRepo{
				settings: &settings.ShopSettings{
					LowMarginThreshold: 20.0, // Default threshold
				},
			}
			service := NewProfitAnalyzerService(menuRepo, nil, settingsRepo, nil)

			// Calculate profit metrics with threshold
			result := service.calculateProfitMetrics(menuRepo.items[0], 20.0)

			// Verify warning status is "loss"
			if result.WarningStatus != WarningLoss {
				t.Logf("Expected warning status 'loss' for cost > price, got %v (price=%v, cost=%v)",
					result.WarningStatus, price, cost)
				return false
			}

			return true
		},
		gen.Float64Range(1.0, 100000.0),   // price
		gen.Float64Range(0.1, 50000.0),    // costExcess (amount by which cost exceeds price)
	))

	properties.TestingRun(t, gopter.ConsoleReporter(false))
}

// Feature: menu-cost-profit-analysis, Property 11: Low Margin Detection
// **Validates: Requirements 3.2**
//
// Property: For any menu item where profit_margin is below the configured low_margin_threshold
// and cost does not exceed price, the warning_status should be marked as "low_margin".
func TestProperty_LowMarginDetection(t *testing.T) {
	parameters := gopter.DefaultTestParameters()
	parameters.MinSuccessfulTests = 100
	properties := gopter.NewProperties(parameters)

	properties.Property("Items with profit margin below threshold are marked as low_margin", prop.ForAll(
		func(price, threshold float64) bool {
			// Skip invalid cases
			if price <= 0 || threshold <= 0 || threshold >= 100 {
				return true
			}

			// Calculate cost that gives profit margin just below threshold
			// profit_margin = ((price - cost) / price) * 100
			// threshold = ((price - cost) / price) * 100
			// cost = price * (1 - threshold/100)
			// To be below threshold, add a small amount to cost
			cost := price * (1 - (threshold-1)/100) // 1% below threshold

			// Ensure cost doesn't exceed price (that would be a loss, not low margin)
			if cost >= price {
				return true
			}

			// Setup test repositories
			menuRepo := &mockMenuRepo{
				items: []*menu.MenuItem{
					{
						ID:          primitive.NewObjectID(),
						Name:        "Low Margin Item",
						Category:    "Test",
						Price:       price,
						CurrentCost: cost,
						CostStatus:  menu.CostStatusFinal,
						CostLastCalculatedAt: time.Now(),
					},
				},
			}

			// Create service
			settingsRepo := &mockSettingsRepo{
				settings: &settings.ShopSettings{
					LowMarginThreshold: threshold,
				},
			}
			service := NewProfitAnalyzerService(menuRepo, nil, settingsRepo, nil)

			// Calculate profit metrics with threshold
			result := service.calculateProfitMetrics(menuRepo.items[0], threshold)

			// Calculate actual profit margin
			actualMargin := ((price - cost) / price) * 100

			// If margin is below threshold and cost < price, should be low_margin
			if actualMargin < threshold && cost < price {
				if result.WarningStatus != WarningLowMargin {
					t.Logf("Expected warning status 'low_margin' for margin=%v < threshold=%v, got %v (price=%v, cost=%v)",
						actualMargin, threshold, result.WarningStatus, price, cost)
					return false
				}
			}

			return true
		},
		gen.Float64Range(100.0, 100000.0),  // price
		gen.Float64Range(10.0, 50.0),       // threshold (10-50%)
	))

	properties.TestingRun(t, gopter.ConsoleReporter(false))
}

// Feature: menu-cost-profit-analysis, Property 12: Warning Status Transitions
// **Validates: Requirements 3.6**
//
// Property: For any menu item, when cost or price changes such that the profit_margin crosses
// the low_margin_threshold or the cost-price relationship changes, the warning_status should
// update immediately to reflect the new state.
func TestProperty_WarningStatusTransitions(t *testing.T) {
	parameters := gopter.DefaultTestParameters()
	parameters.MinSuccessfulTests = 100
	properties := gopter.NewProperties(parameters)

	properties.Property("Warning status updates when cost/price changes", prop.ForAll(
		func(initialPrice, initialCost, newPrice, newCost, threshold float64) bool {
			// Skip invalid cases
			if initialPrice <= 0 || newPrice <= 0 || threshold <= 0 || threshold >= 100 {
				return true
			}

			// Setup test repositories with initial values
			menuItem := &menu.MenuItem{
				ID:          primitive.NewObjectID(),
				Name:        "Transition Item",
				Category:    "Test",
				Price:       initialPrice,
				CurrentCost: initialCost,
				CostStatus:  menu.CostStatusFinal,
				CostLastCalculatedAt: time.Now(),
			}

			menuRepo := &mockMenuRepo{
				items: []*menu.MenuItem{menuItem},
			}

			settingsRepo := &mockSettingsRepo{
				settings: &settings.ShopSettings{
					LowMarginThreshold: threshold,
				},
			}
			service := NewProfitAnalyzerService(menuRepo, nil, settingsRepo, nil)

			// Calculate initial warning status
			initialResult := service.calculateProfitMetrics(menuItem, threshold)
			initialStatus := initialResult.WarningStatus

			// Update price and cost
			menuItem.Price = newPrice
			menuItem.CurrentCost = newCost

			// Calculate new warning status
			newResult := service.calculateProfitMetrics(menuItem, threshold)
			newStatus := newResult.WarningStatus

			// Verify the new status is correct based on new values
			expectedStatus := service.determineWarningStatus(newPrice, newCost, newResult.ProfitMargin, threshold)

			if newStatus != expectedStatus {
				t.Logf("Warning status did not update correctly: expected %v, got %v (price: %v->%v, cost: %v->%v, threshold=%v)",
					expectedStatus, newStatus, initialPrice, newPrice, initialCost, newCost, threshold)
				return false
			}

			// Verify status transitions are logical:
			// 1. If cost > price, must be loss
			if newCost > newPrice && newStatus != WarningLoss {
				t.Logf("Expected loss status when cost > price, got %v (price=%v, cost=%v)",
					newStatus, newPrice, newCost)
				return false
			}

			// 2. If cost < price and margin < threshold, must be low_margin
			newMargin := ((newPrice - newCost) / newPrice) * 100
			if newCost < newPrice && newMargin < threshold && newStatus != WarningLowMargin {
				t.Logf("Expected low_margin status when margin=%v < threshold=%v, got %v (price=%v, cost=%v)",
					newMargin, threshold, newStatus, newPrice, newCost)
				return false
			}

			// 3. If cost < price and margin >= threshold, must be none
			if newCost < newPrice && newMargin >= threshold && newStatus != WarningNone {
				t.Logf("Expected none status when margin=%v >= threshold=%v, got %v (price=%v, cost=%v)",
					newMargin, threshold, newStatus, newPrice, newCost)
				return false
			}

			// Log transition for debugging (only if status changed)
			if initialStatus != newStatus {
				t.Logf("Status transition: %v -> %v (price: %v->%v, cost: %v->%v, margin: %v->%v, threshold=%v)",
					initialStatus, newStatus, initialPrice, newPrice, initialCost, newCost,
					initialResult.ProfitMargin, newResult.ProfitMargin, threshold)
			}

			return true
		},
		gen.Float64Range(100.0, 100000.0),  // initialPrice
		gen.Float64Range(0.0, 100000.0),    // initialCost
		gen.Float64Range(100.0, 100000.0),  // newPrice
		gen.Float64Range(0.0, 100000.0),    // newCost
		gen.Float64Range(10.0, 50.0),       // threshold (10-50%)
	))

	properties.TestingRun(t, gopter.ConsoleReporter(false))
}

// Feature: menu-cost-profit-analysis, Property 14: Category Filtering
// **Validates: Requirements 4.3**
//
// Property: For any category filter, the API should return only menu items that belong to that category.
func TestProperty_CategoryFiltering(t *testing.T) {
	parameters := gopter.DefaultTestParameters()
	parameters.MinSuccessfulTests = 100 // Minimum 100 iterations as per spec
	properties := gopter.NewProperties(parameters)

	properties.Property("Category filter returns only items from that category", prop.ForAll(
		func(targetCategory string, numItems int) bool {
			// Skip invalid cases
			if targetCategory == "" || numItems < 1 || numItems > 50 {
				return true
			}

			// Generate menu items with various categories
			categories := []string{targetCategory, "Other1", "Other2", "Other3"}
			menuItems := make([]*menu.MenuItem, 0, numItems)

			for i := 0; i < numItems; i++ {
				// Randomly assign category (ensure at least one item has target category)
				category := categories[i%len(categories)]
				if i == 0 {
					category = targetCategory // Ensure at least one item with target category
				}

				menuItems = append(menuItems, &menu.MenuItem{
					ID:          primitive.NewObjectID(),
					Name:        fmt.Sprintf("Item %d", i),
					Category:    category,
					Price:       float64(10000 + i*1000),
					CurrentCost: float64(5000 + i*500),
					CostStatus:  menu.CostStatusFinal,
					CostLastCalculatedAt: time.Now(),
				})
			}

			// Setup mock repository that filters by category
			menuRepo := &mockMenuRepo{
				items: menuItems,
			}

			settingsRepo := &mockSettingsRepo{
				settings: &settings.ShopSettings{
					LowMarginThreshold: 20.0,
				},
			}

			service := NewProfitAnalyzerService(menuRepo, nil, settingsRepo, nil)

			// Apply category filter
			filter := ProfitFilter{
				Category: targetCategory,
			}

			response, err := service.GetAllMenuItemProfits(context.Background(), filter)

			if err != nil {
				t.Logf("Unexpected error: %v", err)
				return false
			}

			// Verify all returned items belong to target category
			for _, item := range response.Items {
				if item.Category != targetCategory {
					t.Logf("Category filter failed: expected category %v, got %v for item %v",
						targetCategory, item.Category, item.Name)
					return false
				}
			}

			// Verify we got at least one item (since we ensured one exists)
			if len(response.Items) == 0 {
				t.Logf("Category filter returned no items, but at least one should exist for category %v",
					targetCategory)
				return false
			}

			// Verify count matches expected
			expectedCount := 0
			for _, item := range menuItems {
				if item.Category == targetCategory {
					expectedCount++
				}
			}

			if len(response.Items) != expectedCount {
				t.Logf("Category filter count mismatch: expected %d items, got %d for category %v",
					expectedCount, len(response.Items), targetCategory)
				return false
			}

			return true
		},
		gen.OneConstOf("Coffee", "Tea", "Smoothie", "Dessert", "Snack"), // Target category
		gen.IntRange(5, 20), // Number of items to generate
	))

	properties.TestingRun(t, gopter.ConsoleReporter(false))
}

// Feature: menu-cost-profit-analysis, Property 14: Category Filtering (Empty Category)
// **Validates: Requirements 4.3**
//
// Property: When filtering by a category that has no items, the API should return an empty list.
func TestProperty_CategoryFiltering_EmptyCategory(t *testing.T) {
	parameters := gopter.DefaultTestParameters()
	parameters.MinSuccessfulTests = 100
	properties := gopter.NewProperties(parameters)

	properties.Property("Category filter returns empty list for non-existent category", prop.ForAll(
		func(numItems int) bool {
			// Skip invalid cases
			if numItems < 1 || numItems > 50 {
				return true
			}

			// Generate menu items with specific categories (not including "NonExistent")
			categories := []string{"Coffee", "Tea", "Smoothie"}
			menuItems := make([]*menu.MenuItem, 0, numItems)

			for i := 0; i < numItems; i++ {
				category := categories[i%len(categories)]

				menuItems = append(menuItems, &menu.MenuItem{
					ID:          primitive.NewObjectID(),
					Name:        fmt.Sprintf("Item %d", i),
					Category:    category,
					Price:       float64(10000 + i*1000),
					CurrentCost: float64(5000 + i*500),
					CostStatus:  menu.CostStatusFinal,
					CostLastCalculatedAt: time.Now(),
				})
			}

			// Setup mock repository
			menuRepo := &mockMenuRepo{
				items: menuItems,
			}

			settingsRepo := &mockSettingsRepo{
				settings: &settings.ShopSettings{
					LowMarginThreshold: 20.0,
				},
			}

			service := NewProfitAnalyzerService(menuRepo, nil, settingsRepo, nil)

			// Filter by non-existent category
			filter := ProfitFilter{
				Category: "NonExistentCategory",
			}

			response, err := service.GetAllMenuItemProfits(context.Background(), filter)

			if err != nil {
				t.Logf("Unexpected error: %v", err)
				return false
			}

			// Verify empty result
			if len(response.Items) != 0 {
				t.Logf("Expected empty result for non-existent category, got %d items",
					len(response.Items))
				return false
			}

			return true
		},
		gen.IntRange(5, 20), // Number of items to generate
	))

	properties.TestingRun(t, gopter.ConsoleReporter(false))
}

// Feature: menu-cost-profit-analysis, Property 15: Profit Margin Sorting
// **Validates: Requirements 4.4**
//
// Property: For any sort order (ascending or descending), the API should return menu items
// sorted by profit_margin in the specified order.
func TestProperty_ProfitMarginSorting(t *testing.T) {
	parameters := gopter.DefaultTestParameters()
	parameters.MinSuccessfulTests = 100 // Minimum 100 iterations as per spec
	properties := gopter.NewProperties(parameters)

	properties.Property("Items are sorted by profit margin in specified order", prop.ForAll(
		func(numItems int, ascending bool) bool {
			// Skip invalid cases
			if numItems < 2 || numItems > 50 {
				return true
			}

			// Generate menu items with random profit margins
			menuItems := make([]*menu.MenuItem, 0, numItems)

			for i := 0; i < numItems; i++ {
				// Generate random price and cost to create varied profit margins
				price := float64(10000 + i*1000)
				cost := float64(1000 + i*800) // Varied cost to create different margins

				menuItems = append(menuItems, &menu.MenuItem{
					ID:          primitive.NewObjectID(),
					Name:        fmt.Sprintf("Item %d", i),
					Category:    "Test",
					Price:       price,
					CurrentCost: cost,
					CostStatus:  menu.CostStatusFinal,
					CostLastCalculatedAt: time.Now(),
				})
			}

			// Setup mock repository
			menuRepo := &mockMenuRepo{
				items: menuItems,
			}

			settingsRepo := &mockSettingsRepo{
				settings: &settings.ShopSettings{
					LowMarginThreshold: 20.0,
				},
			}

			service := NewProfitAnalyzerService(menuRepo, nil, settingsRepo, nil)

			// Apply sort filter
			sortOrder := "desc"
			if ascending {
				sortOrder = "asc"
			}

			filter := ProfitFilter{
				SortBy:    "profit_margin",
				SortOrder: sortOrder,
			}

			response, err := service.GetAllMenuItemProfits(context.Background(), filter)

			if err != nil {
				t.Logf("Unexpected error: %v", err)
				return false
			}

			// Verify items are sorted correctly
			for i := 0; i < len(response.Items)-1; i++ {
				currentMargin := response.Items[i].ProfitMargin
				nextMargin := response.Items[i+1].ProfitMargin

				if ascending {
					// Ascending order: current <= next
					if currentMargin > nextMargin {
						t.Logf("Ascending sort failed: item[%d].ProfitMargin=%v > item[%d].ProfitMargin=%v",
							i, currentMargin, i+1, nextMargin)
						return false
					}
				} else {
					// Descending order: current >= next
					if currentMargin < nextMargin {
						t.Logf("Descending sort failed: item[%d].ProfitMargin=%v < item[%d].ProfitMargin=%v",
							i, currentMargin, i+1, nextMargin)
						return false
					}
				}
			}

			return true
		},
		gen.IntRange(5, 20),  // Number of items to generate
		gen.Bool(),           // Ascending or descending
	))

	properties.TestingRun(t, gopter.ConsoleReporter(false))
}

// Feature: menu-cost-profit-analysis, Property 15: Profit Margin Sorting (Absolute Profit)
// **Validates: Requirements 4.4**
//
// Property: For any sort order, the API should correctly sort menu items by absolute_profit.
func TestProperty_AbsoluteProfitSorting(t *testing.T) {
	parameters := gopter.DefaultTestParameters()
	parameters.MinSuccessfulTests = 100
	properties := gopter.NewProperties(parameters)

	properties.Property("Items are sorted by absolute profit in specified order", prop.ForAll(
		func(numItems int, ascending bool) bool {
			// Skip invalid cases
			if numItems < 2 || numItems > 50 {
				return true
			}

			// Generate menu items with random absolute profits
			menuItems := make([]*menu.MenuItem, 0, numItems)

			for i := 0; i < numItems; i++ {
				// Generate varied price and cost combinations
				price := float64(5000 + i*2000)
				cost := float64(1000 + i*500)

				menuItems = append(menuItems, &menu.MenuItem{
					ID:          primitive.NewObjectID(),
					Name:        fmt.Sprintf("Item %d", i),
					Category:    "Test",
					Price:       price,
					CurrentCost: cost,
					CostStatus:  menu.CostStatusFinal,
					CostLastCalculatedAt: time.Now(),
				})
			}

			// Setup mock repository
			menuRepo := &mockMenuRepo{
				items: menuItems,
			}

			settingsRepo := &mockSettingsRepo{
				settings: &settings.ShopSettings{
					LowMarginThreshold: 20.0,
				},
			}

			service := NewProfitAnalyzerService(menuRepo, nil, settingsRepo, nil)

			// Apply sort filter
			sortOrder := "desc"
			if ascending {
				sortOrder = "asc"
			}

			filter := ProfitFilter{
				SortBy:    "absolute_profit",
				SortOrder: sortOrder,
			}

			response, err := service.GetAllMenuItemProfits(context.Background(), filter)

			if err != nil {
				t.Logf("Unexpected error: %v", err)
				return false
			}

			// Verify items are sorted correctly
			for i := 0; i < len(response.Items)-1; i++ {
				currentProfit := response.Items[i].AbsoluteProfit
				nextProfit := response.Items[i+1].AbsoluteProfit

				if ascending {
					// Ascending order: current <= next
					if currentProfit > nextProfit {
						t.Logf("Ascending sort by absolute profit failed: item[%d]=%v > item[%d]=%v",
							i, currentProfit, i+1, nextProfit)
						return false
					}
				} else {
					// Descending order: current >= next
					if currentProfit < nextProfit {
						t.Logf("Descending sort by absolute profit failed: item[%d]=%v < item[%d]=%v",
							i, currentProfit, i+1, nextProfit)
						return false
					}
				}
			}

			return true
		},
		gen.IntRange(5, 20),  // Number of items to generate
		gen.Bool(),           // Ascending or descending
	))

	properties.TestingRun(t, gopter.ConsoleReporter(false))
}

// Feature: menu-cost-profit-analysis, Property 15: Profit Margin Sorting (Name)
// **Validates: Requirements 4.4**
//
// Property: For any sort order, the API should correctly sort menu items by name.
func TestProperty_NameSorting(t *testing.T) {
	parameters := gopter.DefaultTestParameters()
	parameters.MinSuccessfulTests = 100
	properties := gopter.NewProperties(parameters)

	properties.Property("Items are sorted by name in specified order", prop.ForAll(
		func(numItems int, ascending bool) bool {
			// Skip invalid cases
			if numItems < 2 || numItems > 50 {
				return true
			}

			// Generate menu items with alphabetically varied names
			names := []string{"Americano", "Cappuccino", "Espresso", "Latte", "Mocha"}
			menuItems := make([]*menu.MenuItem, 0, numItems)

			for i := 0; i < numItems; i++ {
				name := names[i%len(names)] + fmt.Sprintf(" %d", i)

				menuItems = append(menuItems, &menu.MenuItem{
					ID:          primitive.NewObjectID(),
					Name:        name,
					Category:    "Coffee",
					Price:       float64(10000 + i*1000),
					CurrentCost: float64(5000 + i*500),
					CostStatus:  menu.CostStatusFinal,
					CostLastCalculatedAt: time.Now(),
				})
			}

			// Setup mock repository
			menuRepo := &mockMenuRepo{
				items: menuItems,
			}

			settingsRepo := &mockSettingsRepo{
				settings: &settings.ShopSettings{
					LowMarginThreshold: 20.0,
				},
			}

			service := NewProfitAnalyzerService(menuRepo, nil, settingsRepo, nil)

			// Apply sort filter
			sortOrder := "desc"
			if ascending {
				sortOrder = "asc"
			}

			filter := ProfitFilter{
				SortBy:    "name",
				SortOrder: sortOrder,
			}

			response, err := service.GetAllMenuItemProfits(context.Background(), filter)

			if err != nil {
				t.Logf("Unexpected error: %v", err)
				return false
			}

			// Verify items are sorted correctly
			for i := 0; i < len(response.Items)-1; i++ {
				currentName := response.Items[i].Name
				nextName := response.Items[i+1].Name

				if ascending {
					// Ascending order: current <= next (alphabetically)
					if currentName > nextName {
						t.Logf("Ascending sort by name failed: item[%d]=%v > item[%d]=%v",
							i, currentName, i+1, nextName)
						return false
					}
				} else {
					// Descending order: current >= next (alphabetically)
					if currentName < nextName {
						t.Logf("Descending sort by name failed: item[%d]=%v < item[%d]=%v",
							i, currentName, i+1, nextName)
						return false
					}
				}
			}

			return true
		},
		gen.IntRange(5, 20),  // Number of items to generate
		gen.Bool(),           // Ascending or descending
	))

	properties.TestingRun(t, gopter.ConsoleReporter(false))
}

// Feature: menu-cost-profit-analysis, Property 14 & 15: Combined Filtering and Sorting
// **Validates: Requirements 4.3, 4.4**
//
// Property: When both category filter and sort are applied, the API should return only items
// from that category, sorted in the specified order.
func TestProperty_CombinedFilteringAndSorting(t *testing.T) {
	parameters := gopter.DefaultTestParameters()
	parameters.MinSuccessfulTests = 100
	properties := gopter.NewProperties(parameters)

	properties.Property("Category filter and sort work together correctly", prop.ForAll(
		func(targetCategory string, numItems int, ascending bool) bool {
			// Skip invalid cases
			if targetCategory == "" || numItems < 2 || numItems > 50 {
				return true
			}

			// Generate menu items with various categories
			categories := []string{targetCategory, "Other1", "Other2"}
			menuItems := make([]*menu.MenuItem, 0, numItems)

			for i := 0; i < numItems; i++ {
				// Ensure at least 2 items have target category for meaningful sort test
				category := categories[i%len(categories)]
				if i < 2 {
					category = targetCategory
				}

				price := float64(10000 + i*1000)
				cost := float64(5000 + i*300)

				menuItems = append(menuItems, &menu.MenuItem{
					ID:          primitive.NewObjectID(),
					Name:        fmt.Sprintf("Item %d", i),
					Category:    category,
					Price:       price,
					CurrentCost: cost,
					CostStatus:  menu.CostStatusFinal,
					CostLastCalculatedAt: time.Now(),
				})
			}

			// Setup mock repository
			menuRepo := &mockMenuRepo{
				items: menuItems,
			}

			settingsRepo := &mockSettingsRepo{
				settings: &settings.ShopSettings{
					LowMarginThreshold: 20.0,
				},
			}

			service := NewProfitAnalyzerService(menuRepo, nil, settingsRepo, nil)

			// Apply both category filter and sort
			sortOrder := "desc"
			if ascending {
				sortOrder = "asc"
			}

			filter := ProfitFilter{
				Category:  targetCategory,
				SortBy:    "profit_margin",
				SortOrder: sortOrder,
			}

			response, err := service.GetAllMenuItemProfits(context.Background(), filter)

			if err != nil {
				t.Logf("Unexpected error: %v", err)
				return false
			}

			// Verify all items belong to target category
			for _, item := range response.Items {
				if item.Category != targetCategory {
					t.Logf("Category filter failed in combined test: expected %v, got %v",
						targetCategory, item.Category)
					return false
				}
			}

			// Verify items are sorted correctly
			for i := 0; i < len(response.Items)-1; i++ {
				currentMargin := response.Items[i].ProfitMargin
				nextMargin := response.Items[i+1].ProfitMargin

				if ascending {
					if currentMargin > nextMargin {
						t.Logf("Sort failed in combined test: item[%d]=%v > item[%d]=%v",
							i, currentMargin, i+1, nextMargin)
						return false
					}
				} else {
					if currentMargin < nextMargin {
						t.Logf("Sort failed in combined test: item[%d]=%v < item[%d]=%v",
							i, currentMargin, i+1, nextMargin)
						return false
					}
				}
			}

			return true
		},
		gen.OneConstOf("Coffee", "Tea", "Smoothie"), // Target category
		gen.IntRange(5, 20),  // Number of items to generate
		gen.Bool(),           // Ascending or descending
	))

	properties.TestingRun(t, gopter.ConsoleReporter(false))
}

// Feature: menu-cost-profit-analysis, Property 7: Category Profit Aggregation
// **Validates: Requirements 6.1, 6.2, 6.3**
//
// Property: For any category and date range, the category profit should be calculated as:
// total_revenue = sum of all order item revenues, total_cost = sum of all order item accounting_costs
// (not current_costs), total_profit = total_revenue - total_cost, and
// average_profit_margin = (total_profit / total_revenue) * 100.
func TestProperty_CategoryProfitAggregation(t *testing.T) {
	parameters := gopter.DefaultTestParameters()
	parameters.MinSuccessfulTests = 100 // Minimum 100 iterations as per spec
	properties := gopter.NewProperties(parameters)

	properties.Property("Category profit aggregation uses accounting_cost and calculates correctly", prop.ForAll(
		func(numOrderItems int, categoryIndex int) bool {
			// Skip invalid cases
			if numOrderItems < 1 || numOrderItems > 50 {
				return true
			}

			// Define test categories
			categories := []string{"Coffee", "Tea", "Smoothie", "Dessert"}
			targetCategory := categories[categoryIndex%len(categories)]

			// Generate order items with various categories
			orderItems := make([]*order.OrderItemWithCost, 0, numOrderItems)
			menuItems := make([]*menu.MenuItem, 0, numOrderItems)

			// Track expected values for verification
			expectedRevenue := make(map[string]float64)
			expectedCost := make(map[string]float64)
			expectedItemCount := make(map[string]int)
			orderIDsByCategory := make(map[string]map[primitive.ObjectID]bool)

			// Generate test data
			for i := 0; i < numOrderItems; i++ {
				// Assign category (ensure target category has at least one item)
				category := categories[i%len(categories)]
				if i == 0 {
					category = targetCategory
				}

				// Generate random price, quantity, and accounting cost
				price := float64(5000 + i*1000)
				quantity := 1 + (i % 5) // 1-5 items
				accountingCost := float64(1000+i*500) * float64(quantity) // Total cost for all items

				menuItemID := primitive.NewObjectID()
				orderID := primitive.NewObjectID()

				// Create menu item
				menuItem := &menu.MenuItem{
					ID:          menuItemID,
					Name:        fmt.Sprintf("Item %d", i),
					Category:    category,
					Price:       price,
					CurrentCost: float64(1000 + i*500), // Per-item cost (not used in aggregation)
					CostStatus:  menu.CostStatusFinal,
					CostLastCalculatedAt: time.Now(),
				}
				menuItems = append(menuItems, menuItem)

				// Create order item with accounting cost
				orderItem := &order.OrderItemWithCost{
					ID:               primitive.NewObjectID(),
					OrderID:          orderID,
					MenuItemID:       menuItemID,
					Name:             menuItem.Name,
					Price:            price,
					Quantity:         quantity,
					Subtotal:         price * float64(quantity),
					AccountingCost:   accountingCost, // Total cost for all items in this order item
					CostCalculatedAt: time.Now(),
					CostStatus:       order.CostStatusFinal,
					CreatedAt:        time.Now(),
				}
				orderItems = append(orderItems, orderItem)

				// Track expected values
				revenue := price * float64(quantity)
				expectedRevenue[category] += revenue
				expectedCost[category] += accountingCost
				expectedItemCount[category] += quantity

				// Track unique order IDs per category
				if orderIDsByCategory[category] == nil {
					orderIDsByCategory[category] = make(map[primitive.ObjectID]bool)
				}
				orderIDsByCategory[category][orderID] = true
			}

			// Setup mock repositories
			menuRepo := &mockMenuRepo{
				items: menuItems,
			}

			orderItemRepo := &mockOrderItemRepo{
				items: orderItems,
			}

			settingsRepo := &mockSettingsRepo{
				settings: &settings.ShopSettings{
					LowMarginThreshold: 20.0,
				},
			}

			service := NewProfitAnalyzerService(menuRepo, orderItemRepo, settingsRepo, nil)

			// Call GetCategoryProfits
			dateRange := DateRange{
				Start: time.Now().Add(-24 * time.Hour),
				End:   time.Now().Add(24 * time.Hour),
			}

			results, err := service.GetCategoryProfits(context.Background(), dateRange)

			if err != nil {
				t.Logf("Unexpected error: %v", err)
				return false
			}

			// Verify results for each category
			resultMap := make(map[string]CategoryProfit)
			for _, result := range results {
				resultMap[result.Category] = result
			}

			// Check each category that should have results
			for category, expRevenue := range expectedRevenue {
				result, exists := resultMap[category]
				if !exists {
					t.Logf("Category %v missing from results", category)
					return false
				}

				// Verify total_revenue = sum of all order item revenues
				tolerance := 0.01
				if math.Abs(result.TotalRevenue-expRevenue) > tolerance {
					t.Logf("Category %v: TotalRevenue mismatch: expected %v, got %v",
						category, expRevenue, result.TotalRevenue)
					return false
				}

				// Verify total_cost = sum of all order item accounting_costs (not current_costs)
				expCost := expectedCost[category]
				if math.Abs(result.TotalCost-expCost) > tolerance {
					t.Logf("Category %v: TotalCost mismatch: expected %v, got %v",
						category, expCost, result.TotalCost)
					return false
				}

				// Verify total_profit = total_revenue - total_cost
				expProfit := expRevenue - expCost
				if math.Abs(result.TotalProfit-expProfit) > tolerance {
					t.Logf("Category %v: TotalProfit mismatch: expected %v, got %v",
						category, expProfit, result.TotalProfit)
					return false
				}

				// Verify average_profit_margin = (total_profit / total_revenue) * 100
				expMargin := (expProfit / expRevenue) * 100
				expMargin = math.Round(expMargin*100) / 100 // Round to 2 decimal places
				if math.Abs(result.AverageProfitMargin-expMargin) > tolerance {
					t.Logf("Category %v: AverageProfitMargin mismatch: expected %v, got %v",
						category, expMargin, result.AverageProfitMargin)
					return false
				}

				// Verify item count
				expItemCount := expectedItemCount[category]
				if result.ItemCount != expItemCount {
					t.Logf("Category %v: ItemCount mismatch: expected %v, got %v",
						category, expItemCount, result.ItemCount)
					return false
				}

				// Verify order count (unique orders)
				expOrderCount := len(orderIDsByCategory[category])
				if result.OrderCount != expOrderCount {
					t.Logf("Category %v: OrderCount mismatch: expected %v, got %v",
						category, expOrderCount, result.OrderCount)
					return false
				}
			}

			return true
		},
		gen.IntRange(5, 30),  // Number of order items to generate
		gen.IntRange(0, 100), // Category index (will be modulo'd)
	))

	properties.TestingRun(t, gopter.ConsoleReporter(false))
}

// Feature: menu-cost-profit-analysis, Property 7: Category Profit Aggregation (Accounting Cost)
// **Validates: Requirements 6.2**
//
// Property: Category profit calculation must use accounting_cost from order items, not current_cost
// from menu items. This ensures historical accuracy.
func TestProperty_CategoryProfitAggregation_UsesAccountingCost(t *testing.T) {
	parameters := gopter.DefaultTestParameters()
	parameters.MinSuccessfulTests = 100
	properties := gopter.NewProperties(parameters)

	properties.Property("Category profit uses accounting_cost not current_cost", prop.ForAll(
		func(numItems int) bool {
			// Skip invalid cases
			if numItems < 1 || numItems > 30 {
				return true
			}

			category := "Coffee"
			orderItems := make([]*order.OrderItemWithCost, 0, numItems)
			menuItems := make([]*menu.MenuItem, 0, numItems)

			var expectedTotalCost float64

			// Generate test data where accounting_cost differs from current_cost
			for i := 0; i < numItems; i++ {
				price := float64(10000 + i*1000)
				quantity := 1 + (i % 3)
				
				// accounting_cost is the historical cost at shift closure
				accountingCost := float64(3000+i*200) * float64(quantity)
				
				// current_cost is the current cost (different from accounting_cost)
				currentCost := float64(5000 + i*300) // Per-item cost (should NOT be used)

				menuItemID := primitive.NewObjectID()
				orderID := primitive.NewObjectID()

				menuItem := &menu.MenuItem{
					ID:          menuItemID,
					Name:        fmt.Sprintf("Item %d", i),
					Category:    category,
					Price:       price,
					CurrentCost: currentCost, // This should NOT be used in aggregation
					CostStatus:  menu.CostStatusFinal,
					CostLastCalculatedAt: time.Now(),
				}
				menuItems = append(menuItems, menuItem)

				orderItem := &order.OrderItemWithCost{
					ID:               primitive.NewObjectID(),
					OrderID:          orderID,
					MenuItemID:       menuItemID,
					Name:             menuItem.Name,
					Price:            price,
					Quantity:         quantity,
					Subtotal:         price * float64(quantity),
					AccountingCost:   accountingCost, // This SHOULD be used
					CostCalculatedAt: time.Now(),
					CostStatus:       order.CostStatusFinal,
					CreatedAt:        time.Now(),
				}
				orderItems = append(orderItems, orderItem)

				// Track expected cost (should use accounting_cost)
				expectedTotalCost += accountingCost
			}

			// Setup mock repositories
			menuRepo := &mockMenuRepo{
				items: menuItems,
			}

			orderItemRepo := &mockOrderItemRepo{
				items: orderItems,
			}

			settingsRepo := &mockSettingsRepo{
				settings: &settings.ShopSettings{
					LowMarginThreshold: 20.0,
				},
			}

			service := NewProfitAnalyzerService(menuRepo, orderItemRepo, settingsRepo, nil)

			// Call GetCategoryProfits
			dateRange := DateRange{
				Start: time.Now().Add(-24 * time.Hour),
				End:   time.Now().Add(24 * time.Hour),
			}

			results, err := service.GetCategoryProfits(context.Background(), dateRange)

			if err != nil {
				t.Logf("Unexpected error: %v", err)
				return false
			}

			// Find the Coffee category result
			var coffeeResult *CategoryProfit
			for _, result := range results {
				if result.Category == category {
					r := result
					coffeeResult = &r
					break
				}
			}

			if coffeeResult == nil {
				t.Logf("Category %v not found in results", category)
				return false
			}

			// Verify that total_cost uses accounting_cost, not current_cost
			tolerance := 0.01
			if math.Abs(coffeeResult.TotalCost-expectedTotalCost) > tolerance {
				t.Logf("TotalCost should use accounting_cost: expected %v, got %v",
					expectedTotalCost, coffeeResult.TotalCost)
				return false
			}

			// Calculate what the cost would be if current_cost was incorrectly used
			var incorrectCost float64
			for _, item := range menuItems {
				// Find corresponding order item
				for _, orderItem := range orderItems {
					if orderItem.MenuItemID == item.ID {
						incorrectCost += item.CurrentCost * float64(orderItem.Quantity)
					}
				}
			}

			// Verify that the result does NOT match the incorrect calculation
			if math.Abs(coffeeResult.TotalCost-incorrectCost) < tolerance {
				// This would mean we're using current_cost instead of accounting_cost
				// Only fail if they're actually different
				if math.Abs(expectedTotalCost-incorrectCost) > tolerance {
					t.Logf("TotalCost appears to use current_cost instead of accounting_cost: got %v, should be %v (not %v)",
						coffeeResult.TotalCost, expectedTotalCost, incorrectCost)
					return false
				}
			}

			return true
		},
		gen.IntRange(5, 20), // Number of items to generate
	))

	properties.TestingRun(t, gopter.ConsoleReporter(false))
}

// Feature: menu-cost-profit-analysis, Property 7: Category Profit Aggregation (Skip Incomplete)
// **Validates: Requirements 6.2**
//
// Property: Category profit calculation should skip order items with cost_status = INCOMPLETE.
func TestProperty_CategoryProfitAggregation_SkipsIncomplete(t *testing.T) {
	parameters := gopter.DefaultTestParameters()
	parameters.MinSuccessfulTests = 100
	properties := gopter.NewProperties(parameters)

	properties.Property("Category profit skips items with INCOMPLETE cost status", prop.ForAll(
		func(numCompleteItems, numIncompleteItems int) bool {
			// Skip invalid cases
			if numCompleteItems < 1 || numCompleteItems > 20 || numIncompleteItems < 1 || numIncompleteItems > 20 {
				return true
			}

			category := "Coffee"
			orderItems := make([]*order.OrderItemWithCost, 0)
			menuItems := make([]*menu.MenuItem, 0)

			var expectedTotalRevenue float64
			var expectedTotalCost float64
			var expectedItemCount int

			// Generate complete items (should be included)
			for i := 0; i < numCompleteItems; i++ {
				price := float64(10000 + i*1000)
				quantity := 1 + (i % 3)
				accountingCost := float64(3000+i*200) * float64(quantity)

				menuItemID := primitive.NewObjectID()
				orderID := primitive.NewObjectID()

				menuItem := &menu.MenuItem{
					ID:          menuItemID,
					Name:        fmt.Sprintf("Complete Item %d", i),
					Category:    category,
					Price:       price,
					CurrentCost: float64(3000 + i*200),
					CostStatus:  menu.CostStatusFinal,
					CostLastCalculatedAt: time.Now(),
				}
				menuItems = append(menuItems, menuItem)

				orderItem := &order.OrderItemWithCost{
					ID:               primitive.NewObjectID(),
					OrderID:          orderID,
					MenuItemID:       menuItemID,
					Name:             menuItem.Name,
					Price:            price,
					Quantity:         quantity,
					Subtotal:         price * float64(quantity),
					AccountingCost:   accountingCost,
					CostCalculatedAt: time.Now(),
					CostStatus:       order.CostStatusFinal, // FINAL status
					CreatedAt:        time.Now(),
				}
				orderItems = append(orderItems, orderItem)

				// Track expected values (only complete items)
				expectedTotalRevenue += price * float64(quantity)
				expectedTotalCost += accountingCost
				expectedItemCount += quantity
			}

			// Generate incomplete items (should be skipped)
			for i := 0; i < numIncompleteItems; i++ {
				price := float64(15000 + i*1000)
				quantity := 1 + (i % 3)
				accountingCost := 0.0 // Incomplete items have no cost

				menuItemID := primitive.NewObjectID()
				orderID := primitive.NewObjectID()

				menuItem := &menu.MenuItem{
					ID:          menuItemID,
					Name:        fmt.Sprintf("Incomplete Item %d", i),
					Category:    category,
					Price:       price,
					CurrentCost: 0,
					CostStatus:  menu.CostStatusIncomplete,
					CostLastCalculatedAt: time.Now(),
				}
				menuItems = append(menuItems, menuItem)

				orderItem := &order.OrderItemWithCost{
					ID:               primitive.NewObjectID(),
					OrderID:          orderID,
					MenuItemID:       menuItemID,
					Name:             menuItem.Name,
					Price:            price,
					Quantity:         quantity,
					Subtotal:         price * float64(quantity),
					AccountingCost:   accountingCost,
					CostCalculatedAt: time.Now(),
					CostStatus:       order.CostStatusIncomplete, // INCOMPLETE status
					CreatedAt:        time.Now(),
				}
				orderItems = append(orderItems, orderItem)

				// Do NOT track incomplete items in expected values
			}

			// Setup mock repositories
			menuRepo := &mockMenuRepo{
				items: menuItems,
			}

			orderItemRepo := &mockOrderItemRepo{
				items: orderItems,
			}

			settingsRepo := &mockSettingsRepo{
				settings: &settings.ShopSettings{
					LowMarginThreshold: 20.0,
				},
			}

			service := NewProfitAnalyzerService(menuRepo, orderItemRepo, settingsRepo, nil)

			// Call GetCategoryProfits
			dateRange := DateRange{
				Start: time.Now().Add(-24 * time.Hour),
				End:   time.Now().Add(24 * time.Hour),
			}

			results, err := service.GetCategoryProfits(context.Background(), dateRange)

			if err != nil {
				t.Logf("Unexpected error: %v", err)
				return false
			}

			// Find the Coffee category result
			var coffeeResult *CategoryProfit
			for _, result := range results {
				if result.Category == category {
					r := result
					coffeeResult = &r
					break
				}
			}

			if coffeeResult == nil {
				t.Logf("Category %v not found in results", category)
				return false
			}

			// Verify that incomplete items were skipped
			tolerance := 0.01
			if math.Abs(coffeeResult.TotalRevenue-expectedTotalRevenue) > tolerance {
				t.Logf("TotalRevenue should exclude incomplete items: expected %v, got %v",
					expectedTotalRevenue, coffeeResult.TotalRevenue)
				return false
			}

			if math.Abs(coffeeResult.TotalCost-expectedTotalCost) > tolerance {
				t.Logf("TotalCost should exclude incomplete items: expected %v, got %v",
					expectedTotalCost, coffeeResult.TotalCost)
				return false
			}

			if coffeeResult.ItemCount != expectedItemCount {
				t.Logf("ItemCount should exclude incomplete items: expected %v, got %v",
					expectedItemCount, coffeeResult.ItemCount)
				return false
			}

			return true
		},
		gen.IntRange(3, 10), // Number of complete items
		gen.IntRange(1, 5),  // Number of incomplete items
	))

	properties.TestingRun(t, gopter.ConsoleReporter(false))
}

// Feature: menu-cost-profit-analysis, Property 8: Operating Profit Calculations
// **Validates: Requirements 6.5.1, 6.5.3, 6.5.4**
//
// Property: For any date range with orders and operating expenses, the operating_profit should equal
// gross_profit - total_expenses, where gross_profit = total_revenue - total_cogs, and
// operating_profit_margin = (operating_profit / total_revenue) * 100.
func TestProperty_OperatingProfitCalculations(t *testing.T) {
	parameters := gopter.DefaultTestParameters()
	parameters.MinSuccessfulTests = 100 // Minimum 100 iterations as per spec
	properties := gopter.NewProperties(parameters)

	properties.Property("Operating profit equals gross profit minus total expenses", prop.ForAll(
		func(numOrderItems int, hasExpenses bool) bool {
			// Skip invalid cases
			if numOrderItems < 1 || numOrderItems > 30 {
				return true
			}

			// Generate order items
			orderItems := make([]*order.OrderItemWithCost, 0, numOrderItems)
			menuItems := make([]*menu.MenuItem, 0, numOrderItems)

			var expectedTotalRevenue float64
			var expectedTotalCOGS float64

			for i := 0; i < numOrderItems; i++ {
				price := float64(5000 + i*1000)
				quantity := 1 + (i % 5)
				accountingCost := float64(2000+i*300) * float64(quantity)

				menuItemID := primitive.NewObjectID()
				orderID := primitive.NewObjectID()

				menuItem := &menu.MenuItem{
					ID:          menuItemID,
					Name:        fmt.Sprintf("Item %d", i),
					Category:    "Coffee",
					Price:       price,
					CurrentCost: float64(2000 + i*300),
					CostStatus:  menu.CostStatusFinal,
					CostLastCalculatedAt: time.Now(),
				}
				menuItems = append(menuItems, menuItem)

				orderItem := &order.OrderItemWithCost{
					ID:               primitive.NewObjectID(),
					OrderID:          orderID,
					MenuItemID:       menuItemID,
					Name:             menuItem.Name,
					Price:            price,
					Quantity:         quantity,
					Subtotal:         price * float64(quantity),
					AccountingCost:   accountingCost,
					CostCalculatedAt: time.Now(),
					CostStatus:       order.CostStatusFinal,
					CreatedAt:        time.Now(),
				}
				orderItems = append(orderItems, orderItem)

				// Track expected values
				expectedTotalRevenue += price * float64(quantity)
				expectedTotalCOGS += accountingCost
			}

			// Round expected values to 2 decimal places
			expectedTotalRevenue = math.Round(expectedTotalRevenue*100) / 100
			expectedTotalCOGS = math.Round(expectedTotalCOGS*100) / 100

			// Calculate expected gross profit
			expectedGrossProfit := expectedTotalRevenue - expectedTotalCOGS
			expectedGrossProfit = math.Round(expectedGrossProfit*100) / 100

			// Calculate expected gross profit margin
			var expectedGrossProfitMargin float64
			if expectedTotalRevenue > 0 {
				expectedGrossProfitMargin = (expectedGrossProfit / expectedTotalRevenue) * 100
				expectedGrossProfitMargin = math.Round(expectedGrossProfitMargin*100) / 100
			}

			// Generate operating expenses if needed
			var expenses []*expense.OperatingExpense
			var expectedTotalExpenses float64

			if hasExpenses {
				// Create expense for the date range
				staffSalary := float64(2000000 + numOrderItems*10000)
				rent := float64(1000000 + numOrderItems*5000)
				utilities := float64(500000 + numOrderItems*2000)
				marketingCosts := float64(300000 + numOrderItems*1000)
				otherExpenses := float64(200000 + numOrderItems*500)

				exp := &expense.OperatingExpense{
					ID:             primitive.NewObjectID(),
					PeriodStart:    time.Now().Add(-24 * time.Hour),
					PeriodEnd:      time.Now().Add(24 * time.Hour),
					StaffSalary:    staffSalary,
					Rent:           rent,
					Utilities:      utilities,
					MarketingCosts: marketingCosts,
					OtherExpenses:  otherExpenses,
					TotalExpenses:  staffSalary + rent + utilities + marketingCosts + otherExpenses,
					CreatedAt:      time.Now(),
					UpdatedAt:      time.Now(),
				}
				expenses = append(expenses, exp)

				// Calculate expected total expenses (rounded)
				expectedTotalExpenses = math.Round(exp.TotalExpenses*100) / 100
			}

			// Calculate expected operating profit
			expectedOperatingProfit := expectedGrossProfit - expectedTotalExpenses
			expectedOperatingProfit = math.Round(expectedOperatingProfit*100) / 100

			// Calculate expected operating profit margin
			var expectedOperatingProfitMargin float64
			if expectedTotalRevenue > 0 {
				expectedOperatingProfitMargin = (expectedOperatingProfit / expectedTotalRevenue) * 100
				expectedOperatingProfitMargin = math.Round(expectedOperatingProfitMargin*100) / 100
			}

			// Setup mock repositories
			menuRepo := &mockMenuRepo{
				items: menuItems,
			}

			orderItemRepo := &mockOrderItemRepo{
				items: orderItems,
			}

			operatingExpenseRepo := &mockOperatingExpenseRepo{
				expenses: expenses,
			}

			settingsRepo := &mockSettingsRepo{
				settings: &settings.ShopSettings{
					LowMarginThreshold: 20.0,
				},
			}

			service := NewProfitAnalyzerService(menuRepo, orderItemRepo, settingsRepo, operatingExpenseRepo)

			// Call GetOperatingProfit
			dateRange := DateRange{
				Start: time.Now().Add(-24 * time.Hour),
				End:   time.Now().Add(24 * time.Hour),
			}

			report, err := service.GetOperatingProfit(context.Background(), dateRange)

			if err != nil {
				t.Logf("Unexpected error: %v", err)
				return false
			}

			// Verify total revenue
			tolerance := 0.01
			if math.Abs(report.TotalRevenue-expectedTotalRevenue) > tolerance {
				t.Logf("TotalRevenue mismatch: expected %v, got %v",
					expectedTotalRevenue, report.TotalRevenue)
				return false
			}

			// Verify total COGS
			if math.Abs(report.TotalCOGS-expectedTotalCOGS) > tolerance {
				t.Logf("TotalCOGS mismatch: expected %v, got %v",
					expectedTotalCOGS, report.TotalCOGS)
				return false
			}

			// Verify gross profit = total_revenue - total_cogs
			if math.Abs(report.GrossProfit-expectedGrossProfit) > tolerance {
				t.Logf("GrossProfit mismatch: expected %v, got %v (revenue=%v, cogs=%v)",
					expectedGrossProfit, report.GrossProfit, expectedTotalRevenue, expectedTotalCOGS)
				return false
			}

			// Verify gross profit margin
			if math.Abs(report.GrossProfitMargin-expectedGrossProfitMargin) > tolerance {
				t.Logf("GrossProfitMargin mismatch: expected %v, got %v",
					expectedGrossProfitMargin, report.GrossProfitMargin)
				return false
			}

			// Verify total expenses
			if math.Abs(report.TotalExpenses-expectedTotalExpenses) > tolerance {
				t.Logf("TotalExpenses mismatch: expected %v, got %v",
					expectedTotalExpenses, report.TotalExpenses)
				return false
			}

			// Verify operating profit = gross_profit - total_expenses
			if math.Abs(report.OperatingProfit-expectedOperatingProfit) > tolerance {
				t.Logf("OperatingProfit mismatch: expected %v, got %v (gross_profit=%v, expenses=%v)",
					expectedOperatingProfit, report.OperatingProfit, expectedGrossProfit, expectedTotalExpenses)
				return false
			}

			// Verify operating profit margin = (operating_profit / total_revenue) * 100
			if math.Abs(report.OperatingProfitMargin-expectedOperatingProfitMargin) > tolerance {
				t.Logf("OperatingProfitMargin mismatch: expected %v, got %v",
					expectedOperatingProfitMargin, report.OperatingProfitMargin)
				return false
			}

			return true
		},
		gen.IntRange(5, 20),  // Number of order items
		gen.Bool(),           // Whether to include expenses
	))

	properties.TestingRun(t, gopter.ConsoleReporter(false))
}

// Feature: menu-cost-profit-analysis, Property 8: Operating Profit Calculations (No Expenses)
// **Validates: Requirements 6.5.1, 6.5.7**
//
// Property: When no operating expenses are entered for a period, operating_profit should equal
// gross_profit and the system should display a note "Chưa nhập chi phí vận hành".
func TestProperty_OperatingProfitCalculations_NoExpenses(t *testing.T) {
	parameters := gopter.DefaultTestParameters()
	parameters.MinSuccessfulTests = 100
	properties := gopter.NewProperties(parameters)

	properties.Property("Operating profit equals gross profit when no expenses", prop.ForAll(
		func(numOrderItems int) bool {
			// Skip invalid cases
			if numOrderItems < 1 || numOrderItems > 30 {
				return true
			}

			// Generate order items
			orderItems := make([]*order.OrderItemWithCost, 0, numOrderItems)
			menuItems := make([]*menu.MenuItem, 0, numOrderItems)

			var expectedTotalRevenue float64
			var expectedTotalCOGS float64

			for i := 0; i < numOrderItems; i++ {
				price := float64(5000 + i*1000)
				quantity := 1 + (i % 5)
				accountingCost := float64(2000+i*300) * float64(quantity)

				menuItemID := primitive.NewObjectID()
				orderID := primitive.NewObjectID()

				menuItem := &menu.MenuItem{
					ID:          menuItemID,
					Name:        fmt.Sprintf("Item %d", i),
					Category:    "Coffee",
					Price:       price,
					CurrentCost: float64(2000 + i*300),
					CostStatus:  menu.CostStatusFinal,
					CostLastCalculatedAt: time.Now(),
				}
				menuItems = append(menuItems, menuItem)

				orderItem := &order.OrderItemWithCost{
					ID:               primitive.NewObjectID(),
					OrderID:          orderID,
					MenuItemID:       menuItemID,
					Name:             menuItem.Name,
					Price:            price,
					Quantity:         quantity,
					Subtotal:         price * float64(quantity),
					AccountingCost:   accountingCost,
					CostCalculatedAt: time.Now(),
					CostStatus:       order.CostStatusFinal,
					CreatedAt:        time.Now(),
				}
				orderItems = append(orderItems, orderItem)

				// Track expected values
				expectedTotalRevenue += price * float64(quantity)
				expectedTotalCOGS += accountingCost
			}

			// Round expected values
			expectedTotalRevenue = math.Round(expectedTotalRevenue*100) / 100
			expectedTotalCOGS = math.Round(expectedTotalCOGS*100) / 100

			// Calculate expected gross profit
			expectedGrossProfit := expectedTotalRevenue - expectedTotalCOGS
			expectedGrossProfit = math.Round(expectedGrossProfit*100) / 100

			// Setup mock repositories with NO expenses
			menuRepo := &mockMenuRepo{
				items: menuItems,
			}

			orderItemRepo := &mockOrderItemRepo{
				items: orderItems,
			}

			operatingExpenseRepo := &mockOperatingExpenseRepo{
				expenses: []*expense.OperatingExpense{}, // No expenses
			}

			settingsRepo := &mockSettingsRepo{
				settings: &settings.ShopSettings{
					LowMarginThreshold: 20.0,
				},
			}

			service := NewProfitAnalyzerService(menuRepo, orderItemRepo, settingsRepo, operatingExpenseRepo)

			// Call GetOperatingProfit
			dateRange := DateRange{
				Start: time.Now().Add(-24 * time.Hour),
				End:   time.Now().Add(24 * time.Hour),
			}

			report, err := service.GetOperatingProfit(context.Background(), dateRange)

			if err != nil {
				t.Logf("Unexpected error: %v", err)
				return false
			}

			// Verify operating profit equals gross profit when no expenses
			tolerance := 0.01
			if math.Abs(report.OperatingProfit-expectedGrossProfit) > tolerance {
				t.Logf("OperatingProfit should equal GrossProfit when no expenses: expected %v, got %v",
					expectedGrossProfit, report.OperatingProfit)
				return false
			}

			// Verify operating profit margin equals gross profit margin when no expenses
			if math.Abs(report.OperatingProfitMargin-report.GrossProfitMargin) > tolerance {
				t.Logf("OperatingProfitMargin should equal GrossProfitMargin when no expenses: expected %v, got %v",
					report.GrossProfitMargin, report.OperatingProfitMargin)
				return false
			}

			// Verify total expenses is zero
			if report.TotalExpenses != 0 {
				t.Logf("TotalExpenses should be 0 when no expenses: got %v", report.TotalExpenses)
				return false
			}

			// Verify allocation note
			if report.AllocationNote != "Chưa nhập chi phí vận hành" {
				t.Logf("Expected allocation note 'Chưa nhập chi phí vận hành', got '%v'",
					report.AllocationNote)
				return false
			}

			return true
		},
		gen.IntRange(5, 20), // Number of order items
	))

	properties.TestingRun(t, gopter.ConsoleReporter(false))
}

// Feature: menu-cost-profit-analysis, Property 8: Operating Profit Calculations (Expense Breakdown)
// **Validates: Requirements 6.5.3, 6.5.5**
//
// Property: The operating profit report should show breakdown of all expense categories with
// individual amounts, and total_expenses should equal the sum of all expense categories.
func TestProperty_OperatingProfitCalculations_ExpenseBreakdown(t *testing.T) {
	parameters := gopter.DefaultTestParameters()
	parameters.MinSuccessfulTests = 100
	properties := gopter.NewProperties(parameters)

	properties.Property("Total expenses equals sum of all expense categories", prop.ForAll(
		func(staffSalary, rent, utilities, marketingCosts, otherExpenses float64) bool {
			// Skip invalid cases (negative expenses)
			if staffSalary < 0 || rent < 0 || utilities < 0 || marketingCosts < 0 || otherExpenses < 0 {
				return true
			}

			// Generate minimal order items
			menuItemID := primitive.NewObjectID()
			orderID := primitive.NewObjectID()

			menuItem := &menu.MenuItem{
				ID:          menuItemID,
				Name:        "Test Item",
				Category:    "Coffee",
				Price:       10000,
				CurrentCost: 5000,
				CostStatus:  menu.CostStatusFinal,
				CostLastCalculatedAt: time.Now(),
			}

			orderItem := &order.OrderItemWithCost{
				ID:               primitive.NewObjectID(),
				OrderID:          orderID,
				MenuItemID:       menuItemID,
				Name:             "Test Item",
				Price:            10000,
				Quantity:         1,
				Subtotal:         10000,
				AccountingCost:   5000,
				CostCalculatedAt: time.Now(),
				CostStatus:       order.CostStatusFinal,
				CreatedAt:        time.Now(),
			}

			// Create expense with specified breakdown
			exp := &expense.OperatingExpense{
				ID:             primitive.NewObjectID(),
				PeriodStart:    time.Now().Add(-24 * time.Hour),
				PeriodEnd:      time.Now().Add(24 * time.Hour),
				StaffSalary:    staffSalary,
				Rent:           rent,
				Utilities:      utilities,
				MarketingCosts: marketingCosts,
				OtherExpenses:  otherExpenses,
				TotalExpenses:  staffSalary + rent + utilities + marketingCosts + otherExpenses,
				CreatedAt:      time.Now(),
				UpdatedAt:      time.Now(),
			}

			// Setup mock repositories
			menuRepo := &mockMenuRepo{
				items: []*menu.MenuItem{menuItem},
			}

			orderItemRepo := &mockOrderItemRepo{
				items: []*order.OrderItemWithCost{orderItem},
			}

			operatingExpenseRepo := &mockOperatingExpenseRepo{
				expenses: []*expense.OperatingExpense{exp},
			}

			settingsRepo := &mockSettingsRepo{
				settings: &settings.ShopSettings{
					LowMarginThreshold: 20.0,
				},
			}

			service := NewProfitAnalyzerService(menuRepo, orderItemRepo, settingsRepo, operatingExpenseRepo)

			// Call GetOperatingProfit
			dateRange := DateRange{
				Start: time.Now().Add(-24 * time.Hour),
				End:   time.Now().Add(24 * time.Hour),
			}

			report, err := service.GetOperatingProfit(context.Background(), dateRange)

			if err != nil {
				t.Logf("Unexpected error: %v", err)
				return false
			}

			// Verify individual expense categories are present
			tolerance := 0.01
			if math.Abs(report.StaffSalary-math.Round(staffSalary*100)/100) > tolerance {
				t.Logf("StaffSalary mismatch: expected %v, got %v",
					math.Round(staffSalary*100)/100, report.StaffSalary)
				return false
			}

			if math.Abs(report.Rent-math.Round(rent*100)/100) > tolerance {
				t.Logf("Rent mismatch: expected %v, got %v",
					math.Round(rent*100)/100, report.Rent)
				return false
			}

			if math.Abs(report.Utilities-math.Round(utilities*100)/100) > tolerance {
				t.Logf("Utilities mismatch: expected %v, got %v",
					math.Round(utilities*100)/100, report.Utilities)
				return false
			}

			if math.Abs(report.MarketingCosts-math.Round(marketingCosts*100)/100) > tolerance {
				t.Logf("MarketingCosts mismatch: expected %v, got %v",
					math.Round(marketingCosts*100)/100, report.MarketingCosts)
				return false
			}

			if math.Abs(report.OtherExpenses-math.Round(otherExpenses*100)/100) > tolerance {
				t.Logf("OtherExpenses mismatch: expected %v, got %v",
					math.Round(otherExpenses*100)/100, report.OtherExpenses)
				return false
			}

			// Verify total expenses equals sum of all categories
			calculatedTotal := report.StaffSalary + report.Rent + report.Utilities + report.MarketingCosts + report.OtherExpenses
			calculatedTotal = math.Round(calculatedTotal*100) / 100

			if math.Abs(report.TotalExpenses-calculatedTotal) > tolerance {
				t.Logf("TotalExpenses should equal sum of categories: expected %v, got %v",
					calculatedTotal, report.TotalExpenses)
				return false
			}

			// Note: We don't verify exact match with expectedTotalExpenses because rounding order matters
			// (rounding individual components then summing vs summing then rounding can differ by 1-2 cents)
			// The important property is that TotalExpenses = sum of all category fields in the report

			return true
		},
		gen.Float64Range(0, 5000000),  // staffSalary
		gen.Float64Range(0, 3000000),  // rent
		gen.Float64Range(0, 1000000),  // utilities
		gen.Float64Range(0, 1000000),  // marketingCosts
		gen.Float64Range(0, 1000000),  // otherExpenses
	))

	properties.TestingRun(t, gopter.ConsoleReporter(false))
}

// Feature: menu-cost-profit-analysis, Property 8: Operating Profit Calculations (Negative Profit)
// **Validates: Requirements 6.5.4**
//
// Property: When total expenses exceed gross profit, operating_profit should be negative and
// operating_profit_margin should be negative.
func TestProperty_OperatingProfitCalculations_NegativeProfit(t *testing.T) {
	parameters := gopter.DefaultTestParameters()
	parameters.MinSuccessfulTests = 100
	properties := gopter.NewProperties(parameters)

	properties.Property("Operating profit can be negative when expenses exceed gross profit", prop.ForAll(
		func(numOrderItems int, expenseMultiplier float64) bool {
			// Skip invalid cases
			if numOrderItems < 1 || numOrderItems > 20 || expenseMultiplier < 1.5 || expenseMultiplier > 10 {
				return true
			}

			// Generate order items
			orderItems := make([]*order.OrderItemWithCost, 0, numOrderItems)
			menuItems := make([]*menu.MenuItem, 0, numOrderItems)

			var expectedTotalRevenue float64
			var expectedTotalCOGS float64

			for i := 0; i < numOrderItems; i++ {
				price := float64(5000 + i*1000)
				quantity := 1 + (i % 3)
				accountingCost := float64(2000+i*300) * float64(quantity)

				menuItemID := primitive.NewObjectID()
				orderID := primitive.NewObjectID()

				menuItem := &menu.MenuItem{
					ID:          menuItemID,
					Name:        fmt.Sprintf("Item %d", i),
					Category:    "Coffee",
					Price:       price,
					CurrentCost: float64(2000 + i*300),
					CostStatus:  menu.CostStatusFinal,
					CostLastCalculatedAt: time.Now(),
				}
				menuItems = append(menuItems, menuItem)

				orderItem := &order.OrderItemWithCost{
					ID:               primitive.NewObjectID(),
					OrderID:          orderID,
					MenuItemID:       menuItemID,
					Name:             menuItem.Name,
					Price:            price,
					Quantity:         quantity,
					Subtotal:         price * float64(quantity),
					AccountingCost:   accountingCost,
					CostCalculatedAt: time.Now(),
					CostStatus:       order.CostStatusFinal,
					CreatedAt:        time.Now(),
				}
				orderItems = append(orderItems, orderItem)

				expectedTotalRevenue += price * float64(quantity)
				expectedTotalCOGS += accountingCost
			}

			// Calculate gross profit
			expectedGrossProfit := expectedTotalRevenue - expectedTotalCOGS

			// Create expenses that exceed gross profit
			totalExpenses := expectedGrossProfit * expenseMultiplier

			exp := &expense.OperatingExpense{
				ID:             primitive.NewObjectID(),
				PeriodStart:    time.Now().Add(-24 * time.Hour),
				PeriodEnd:      time.Now().Add(24 * time.Hour),
				StaffSalary:    totalExpenses * 0.5,
				Rent:           totalExpenses * 0.3,
				Utilities:      totalExpenses * 0.1,
				MarketingCosts: totalExpenses * 0.05,
				OtherExpenses:  totalExpenses * 0.05,
				TotalExpenses:  totalExpenses,
				CreatedAt:      time.Now(),
				UpdatedAt:      time.Now(),
			}

			// Setup mock repositories
			menuRepo := &mockMenuRepo{
				items: menuItems,
			}

			orderItemRepo := &mockOrderItemRepo{
				items: orderItems,
			}

			operatingExpenseRepo := &mockOperatingExpenseRepo{
				expenses: []*expense.OperatingExpense{exp},
			}

			settingsRepo := &mockSettingsRepo{
				settings: &settings.ShopSettings{
					LowMarginThreshold: 20.0,
				},
			}

			service := NewProfitAnalyzerService(menuRepo, orderItemRepo, settingsRepo, operatingExpenseRepo)

			// Call GetOperatingProfit
			dateRange := DateRange{
				Start: time.Now().Add(-24 * time.Hour),
				End:   time.Now().Add(24 * time.Hour),
			}

			report, err := service.GetOperatingProfit(context.Background(), dateRange)

			if err != nil {
				t.Logf("Unexpected error: %v", err)
				return false
			}

			// Verify operating profit is negative
			if report.OperatingProfit >= 0 {
				t.Logf("Expected negative operating profit when expenses exceed gross profit, got %v (gross_profit=%v, expenses=%v)",
					report.OperatingProfit, report.GrossProfit, report.TotalExpenses)
				return false
			}

			// Verify operating profit margin is negative
			if report.OperatingProfitMargin >= 0 {
				t.Logf("Expected negative operating profit margin when expenses exceed gross profit, got %v",
					report.OperatingProfitMargin)
				return false
			}

			// Verify the calculation is correct
			expectedOperatingProfit := report.GrossProfit - report.TotalExpenses
			expectedOperatingProfit = math.Round(expectedOperatingProfit*100) / 100

			tolerance := 0.01
			if math.Abs(report.OperatingProfit-expectedOperatingProfit) > tolerance {
				t.Logf("OperatingProfit calculation incorrect: expected %v, got %v",
					expectedOperatingProfit, report.OperatingProfit)
				return false
			}

			return true
		},
		gen.IntRange(3, 15),       // Number of order items
		gen.Float64Range(1.5, 5.0), // Expense multiplier (expenses = gross_profit * multiplier)
	))

	properties.TestingRun(t, gopter.ConsoleReporter(false))
}
