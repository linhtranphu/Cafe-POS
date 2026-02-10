package services

import (
	"context"
	"testing"

	"cafe-pos/backend/domain/ingredient"
	"cafe-pos/backend/domain/menu"
	"cafe-pos/backend/domain/settings"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

// TestWarningDetectionWorkflow tests the complete warning detection workflow
// Requirements: 3.1, 3.2, 3.3, 3.6
// - Create menu items with various cost/price ratios
// - View menu cost list with warnings
// - Adjust low margin threshold
// - Verify warning status updates
func TestWarningDetectionWorkflow(t *testing.T) {
	ctx := context.Background()

	// Create mock repositories
	menuRepo := &mockMenuRepository{
		menuItems: make(map[primitive.ObjectID]*menu.MenuItem),
	}
	ingredientRepo := &mockIngredientRepository{
		ingredients: make(map[primitive.ObjectID]*ingredient.Ingredient),
	}
	settingsRepo := &mockSettingsRepo{
		settings: &settings.ShopSettings{
			ID:                 primitive.NewObjectID(),
			LowMarginThreshold: 20.0, // Default threshold
		},
	}

	// Initialize services
	costCalculator := NewCostCalculatorService(menuRepo, ingredientRepo, nil, nil)
	profitAnalyzer := NewProfitAnalyzerService(menuRepo, nil, settingsRepo, nil)

	// Step 1: Create test ingredients
	espressoIngredient := &ingredient.Ingredient{
		ID:                primitive.NewObjectID(),
		Name:              "Espresso",
		CostPerUnit:       200,
		ConversionRate:    1.0,
		WastagePercentage: 0.0,
	}
	ingredientRepo.ingredients[espressoIngredient.ID] = espressoIngredient

	milkIngredient := &ingredient.Ingredient{
		ID:                primitive.NewObjectID(),
		Name:              "Milk",
		CostPerUnit:       50,
		ConversionRate:    1.0,
		WastagePercentage: 0.0,
	}
	ingredientRepo.ingredients[milkIngredient.ID] = milkIngredient

	// Step 2: Create menu items with various cost/price ratios
	t.Log("Creating menu items with various cost/price ratios...")

	// Item 1: Profitable item (high margin > 20%)
	// Cost: 30*200 = 6000, Price: 45000, Margin: 86.67%
	profitableItem := &menu.MenuItem{
		ID:       primitive.NewObjectID(),
		Name:     "Profitable Coffee",
		Price:    45000,
		Category: "Coffee",
		Ingredients: []menu.Ingredient{
			{Name: espressoIngredient.Name, Quantity: 30, Unit: "g"},
		},
		Available: true,
	}
	menuRepo.menuItems[profitableItem.ID] = profitableItem

	// Item 2: Low margin item (margin < 20%)
	// Cost: 30*200 + 150*50 = 6000 + 7500 = 13500, Price: 15000, Margin: 10%
	lowMarginItem := &menu.MenuItem{
		ID:       primitive.NewObjectID(),
		Name:     "Low Margin Coffee",
		Price:    15000,
		Category: "Coffee",
		Ingredients: []menu.Ingredient{
			{Name: espressoIngredient.Name, Quantity: 30, Unit: "g"},
			{Name: milkIngredient.Name, Quantity: 150, Unit: "ml"},
		},
		Available: true,
	}
	menuRepo.menuItems[lowMarginItem.ID] = lowMarginItem

	// Item 3: Loss item (cost > price)
	// Cost: 30*200 + 200*50 = 6000 + 10000 = 16000, Price: 12000, Margin: -33.33%
	lossItem := &menu.MenuItem{
		ID:       primitive.NewObjectID(),
		Name:     "Loss Coffee",
		Price:    12000,
		Category: "Coffee",
		Ingredients: []menu.Ingredient{
			{Name: espressoIngredient.Name, Quantity: 30, Unit: "g"},
			{Name: milkIngredient.Name, Quantity: 200, Unit: "ml"},
		},
		Available: true,
	}
	menuRepo.menuItems[lossItem.ID] = lossItem

	// Item 4: Break-even item (margin ≈ 0%)
	// Cost: 30*200 = 6000, Price: 6000, Margin: 0%
	breakEvenItem := &menu.MenuItem{
		ID:       primitive.NewObjectID(),
		Name:     "Break Even Coffee",
		Price:    6000,
		Category: "Coffee",
		Ingredients: []menu.Ingredient{
			{Name: espressoIngredient.Name, Quantity: 30, Unit: "g"},
		},
		Available: true,
	}
	menuRepo.menuItems[breakEvenItem.ID] = breakEvenItem

	// Step 3: Calculate costs for all menu items
	t.Log("Calculating costs for all menu items...")
	for id, item := range menuRepo.menuItems {
		costResult, err := costCalculator.CalculateMenuItemCost(ctx, id)
		require.NoError(t, err)
		
		// Update the menu item with the calculated cost
		item.CurrentCost = costResult.CurrentCost
		item.CostStatus = costResult.CostStatus
		item.CostLastCalculatedAt = costResult.CostLastCalculatedAt
		menuRepo.menuItems[id] = item
	}

	// Step 4: View menu cost list with warnings (default threshold = 20%)
	t.Log("Viewing menu cost list with warnings (threshold = 20%)...")
	
	filter := ProfitFilter{
		SortBy:    "profit_margin",
		SortOrder: "desc",
	}
	
	menuCosts, err := profitAnalyzer.GetAllMenuItemProfits(ctx, filter)
	require.NoError(t, err)
	require.Len(t, menuCosts.Items, 4, "Should have 4 menu items")

	// Verify warning statuses with default threshold (20%)
	var lowMarginCount, lossCount, noneCount int
	for _, item := range menuCosts.Items {
		t.Logf("Item: %s, Price: %.2f, Cost: %.2f, Margin: %.2f%%, Warning: %s",
			item.Name, item.Price, item.CurrentCost, item.ProfitMargin, item.WarningStatus)
		
		switch item.WarningStatus {
		case "loss":
			lossCount++
			assert.True(t, item.CurrentCost > item.Price,
				"Loss item %s should have cost > price", item.Name)
		case "low_margin":
			lowMarginCount++
			assert.True(t, item.ProfitMargin < 20.0 && item.CurrentCost <= item.Price,
				"Low margin item %s should have margin < 20%% and cost <= price", item.Name)
		case "none":
			noneCount++
			assert.True(t, item.ProfitMargin >= 20.0 && item.CurrentCost <= item.Price,
				"Normal item %s should have margin >= 20%% and cost <= price", item.Name)
		}
	}

	assert.Equal(t, 1, lossCount, "Should have 1 loss item")
	assert.Equal(t, 2, lowMarginCount, "Should have 2 low margin items (low margin + break-even)")
	assert.Equal(t, 1, noneCount, "Should have 1 normal item")

	t.Logf("✓ Warning detection with threshold 20%%: %d loss, %d low margin, %d normal",
		lossCount, lowMarginCount, noneCount)

	// Step 5: Get warnings summary
	warnings, err := profitAnalyzer.DetectWarningStatus(ctx, 20.0)
	require.NoError(t, err)
	
	assert.Equal(t, 1, warnings.LossCount, "Should have 1 loss item")
	assert.Equal(t, 2, warnings.LowMarginCount, "Should have 2 low margin items")
	assert.Len(t, warnings.LossItems, 1, "Should have 1 loss item in list")
	assert.Len(t, warnings.LowMarginItems, 2, "Should have 2 low margin items in list")

	// Step 6: Adjust low margin threshold to 15%
	t.Log("Adjusting low margin threshold to 15%...")
	settingsRepo.settings.LowMarginThreshold = 15.0

	// Step 7: View menu cost list with new threshold
	t.Log("Viewing menu cost list with warnings (threshold = 15%)...")
	
	menuCosts2, err := profitAnalyzer.GetAllMenuItemProfits(ctx, filter)
	require.NoError(t, err)
	require.Len(t, menuCosts2.Items, 4, "Should have 4 menu items")

	// Verify warning statuses with new threshold (15%)
	var lowMarginCount2, lossCount2, noneCount2 int
	for _, item := range menuCosts2.Items {
		t.Logf("Item: %s, Price: %.2f, Cost: %.2f, Margin: %.2f%%, Warning: %s",
			item.Name, item.Price, item.CurrentCost, item.ProfitMargin, item.WarningStatus)
		
		switch item.WarningStatus {
		case "loss":
			lossCount2++
		case "low_margin":
			lowMarginCount2++
			assert.True(t, item.ProfitMargin < 15.0 && item.CurrentCost <= item.Price,
				"Low margin item %s should have margin < 15%% with new threshold", item.Name)
		case "none":
			noneCount2++
			assert.True(t, item.ProfitMargin >= 15.0 && item.CurrentCost <= item.Price,
				"Normal item %s should have margin >= 15%% with new threshold", item.Name)
		}
	}

	assert.Equal(t, 1, lossCount2, "Should still have 1 loss item")
	assert.Equal(t, 2, lowMarginCount2, "Should still have 2 low margin items (low margin + break-even)")
	assert.Equal(t, 1, noneCount2, "Should still have 1 normal item")

	t.Logf("✓ Warning detection with threshold 15%%: %d loss, %d low margin, %d normal",
		lossCount2, lowMarginCount2, noneCount2)

	// Step 8: Adjust threshold to 30% (higher threshold)
	t.Log("Adjusting low margin threshold to 30%...")
	settingsRepo.settings.LowMarginThreshold = 30.0

	// Step 9: View menu cost list with higher threshold
	t.Log("Viewing menu cost list with warnings (threshold = 30%)...")
	
	menuCosts3, err := profitAnalyzer.GetAllMenuItemProfits(ctx, filter)
	require.NoError(t, err)

	// Verify warning statuses with higher threshold (30%)
	var lowMarginCount3, lossCount3, noneCount3 int
	for _, item := range menuCosts3.Items {
		t.Logf("Item: %s, Price: %.2f, Cost: %.2f, Margin: %.2f%%, Warning: %s",
			item.Name, item.Price, item.CurrentCost, item.ProfitMargin, item.WarningStatus)
		
		switch item.WarningStatus {
		case "loss":
			lossCount3++
		case "low_margin":
			lowMarginCount3++
			assert.True(t, item.ProfitMargin < 30.0 && item.CurrentCost <= item.Price,
				"Low margin item %s should have margin < 30%% with higher threshold", item.Name)
		case "none":
			noneCount3++
			assert.True(t, item.ProfitMargin >= 30.0 && item.CurrentCost <= item.Price,
				"Normal item %s should have margin >= 30%% with higher threshold", item.Name)
		}
	}

	assert.Equal(t, 1, lossCount3, "Should still have 1 loss item")
	assert.Equal(t, 2, lowMarginCount3, "Should have 2 low margin items with higher threshold")
	assert.Equal(t, 1, noneCount3, "Should have 1 normal item (only the highly profitable one)")

	t.Logf("✓ Warning detection with threshold 30%%: %d loss, %d low margin, %d normal",
		lossCount3, lowMarginCount3, noneCount3)

	// Step 10: Verify warning status transitions
	t.Log("Verifying warning status transitions...")
	
	// Get warnings with different thresholds
	warnings20, err := profitAnalyzer.DetectWarningStatus(ctx, 20.0)
	require.NoError(t, err)
	
	warnings30, err := profitAnalyzer.DetectWarningStatus(ctx, 30.0)
	require.NoError(t, err)
	
	// With threshold 20%, we should have 2 low margin items
	assert.Equal(t, 2, warnings20.LowMarginCount, "Should have 2 low margin items with 20%% threshold")
	
	// With threshold 30%, we should have 2 low margin items (same items, but different threshold)
	assert.Equal(t, 2, warnings30.LowMarginCount, "Should have 2 low margin items with 30%% threshold")
	
	// Loss count should be the same regardless of threshold
	assert.Equal(t, warnings20.LossCount, warnings30.LossCount, "Loss count should not change with threshold")

	t.Log("✅ Warning detection workflow test passed!")
	t.Log("   - Menu items created with various cost/price ratios")
	t.Log("   - Warning statuses detected correctly with default threshold")
	t.Log("   - Low margin threshold adjusted successfully")
	t.Log("   - Warning statuses updated correctly with new thresholds")
	t.Log("   - Warning status transitions verified")
}
