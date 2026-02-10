package services

import (
	"context"
	"testing"
	"time"

	"cafe-pos/backend/domain/expense"
	"cafe-pos/backend/domain/ingredient"
	"cafe-pos/backend/domain/menu"
	"cafe-pos/backend/domain/order"
	"cafe-pos/backend/domain/settings"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

// TestCompleteProfitAnalysisWorkflow tests the complete profit analysis workflow
// Requirements: 6.1, 6.5.1
// - Create orders with various menu items
// - Close shift
// - View category profit report
// - Add operating expenses
// - View operating profit report
func TestCompleteProfitAnalysisWorkflow(t *testing.T) {
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
	operatingExpenseRepo := &mockOperatingExpenseRepository{
		expenses: make(map[primitive.ObjectID]*expense.OperatingExpense),
	}
	settingsRepo := &mockSettingsRepo{
		settings: &settings.ShopSettings{
			ID:                 primitive.NewObjectID(),
			LowMarginThreshold: 20.0,
		},
	}

	// Initialize services
	costCalculator := NewCostCalculatorService(menuRepo, ingredientRepo, orderRepo, orderItemRepo)
	profitAnalyzer := NewProfitAnalyzerService(menuRepo, orderItemRepo, settingsRepo, operatingExpenseRepo)

	// Step 1: Create test ingredients
	espressoIngredient := &ingredient.Ingredient{
		ID:                primitive.NewObjectID(),
		Name:              "Espresso",
		CostPerUnit:       200,
		ConversionRate:    1.0,
		WastagePercentage: 5.0,
	}
	ingredientRepo.ingredients[espressoIngredient.ID] = espressoIngredient

	milkIngredient := &ingredient.Ingredient{
		ID:                primitive.NewObjectID(),
		Name:              "Milk",
		CostPerUnit:       50,
		ConversionRate:    1.0,
		WastagePercentage: 10.0,
	}
	ingredientRepo.ingredients[milkIngredient.ID] = milkIngredient

	teaLeavesIngredient := &ingredient.Ingredient{
		ID:                primitive.NewObjectID(),
		Name:              "Tea Leaves",
		CostPerUnit:       100,
		ConversionRate:    1.0,
		WastagePercentage: 0.0,
	}
	ingredientRepo.ingredients[teaLeavesIngredient.ID] = teaLeavesIngredient

	// Step 2: Create test menu items in different categories
	cappuccinoItem := &menu.MenuItem{
		ID:       primitive.NewObjectID(),
		Name:     "Cappuccino",
		Price:    45000,
		Category: "Coffee",
		Ingredients: []menu.Ingredient{
			{Name: espressoIngredient.Name, Quantity: 30, Unit: "g"},
			{Name: milkIngredient.Name, Quantity: 150, Unit: "ml"},
		},
		Available: true,
	}
	menuRepo.menuItems[cappuccinoItem.ID] = cappuccinoItem

	latteItem := &menu.MenuItem{
		ID:       primitive.NewObjectID(),
		Name:     "Latte",
		Price:    50000,
		Category: "Coffee",
		Ingredients: []menu.Ingredient{
			{Name: espressoIngredient.Name, Quantity: 30, Unit: "g"},
			{Name: milkIngredient.Name, Quantity: 200, Unit: "ml"},
		},
		Available: true,
	}
	menuRepo.menuItems[latteItem.ID] = latteItem

	greenTeaItem := &menu.MenuItem{
		ID:       primitive.NewObjectID(),
		Name:     "Green Tea",
		Price:    30000,
		Category: "Tea",
		Ingredients: []menu.Ingredient{
			{Name: teaLeavesIngredient.Name, Quantity: 10, Unit: "g"},
		},
		Available: true,
	}
	menuRepo.menuItems[greenTeaItem.ID] = greenTeaItem

	milkTeaItem := &menu.MenuItem{
		ID:       primitive.NewObjectID(),
		Name:     "Milk Tea",
		Price:    35000,
		Category: "Tea",
		Ingredients: []menu.Ingredient{
			{Name: teaLeavesIngredient.Name, Quantity: 10, Unit: "g"},
			{Name: milkIngredient.Name, Quantity: 100, Unit: "ml"},
		},
		Available: true,
	}
	menuRepo.menuItems[milkTeaItem.ID] = milkTeaItem

	// Step 3: Create shift
	testShift := &order.Shift{
		ID:        primitive.NewObjectID(),
		UserID:    primitive.NewObjectID(),
		UserName:  "waiter1",
		RoleType:  order.RoleWaiter,
		Status:    order.ShiftOpen,
		StartedAt: time.Now().Add(-4 * time.Hour),
		CreatedAt: time.Now().Add(-4 * time.Hour),
		UpdatedAt: time.Now().Add(-4 * time.Hour),
	}

	// Step 4: Create orders with various menu items
	t.Log("Creating orders with various menu items...")
	
	// Order 1: Coffee items
	order1 := &order.Order{
		ID:       primitive.NewObjectID(),
		ShiftID:  testShift.ID,
		WaiterID: testShift.UserID,
		Items: []order.OrderItem{
			{
				MenuItemID: cappuccinoItem.ID,
				Name:       cappuccinoItem.Name,
				Price:      cappuccinoItem.Price,
				Quantity:   3,
				Subtotal:   cappuccinoItem.Price * 3,
			},
			{
				MenuItemID: latteItem.ID,
				Name:       latteItem.Name,
				Price:      latteItem.Price,
				Quantity:   2,
				Subtotal:   latteItem.Price * 2,
			},
		},
		Subtotal:  cappuccinoItem.Price*3 + latteItem.Price*2,
		Total:     cappuccinoItem.Price*3 + latteItem.Price*2,
		Status:    order.StatusServed,
		CreatedAt: time.Now().Add(-3 * time.Hour),
		UpdatedAt: time.Now().Add(-3 * time.Hour),
	}
	orderRepo.orders[order1.ID] = order1

	// Order 2: Tea items
	order2 := &order.Order{
		ID:       primitive.NewObjectID(),
		ShiftID:  testShift.ID,
		WaiterID: testShift.UserID,
		Items: []order.OrderItem{
			{
				MenuItemID: greenTeaItem.ID,
				Name:       greenTeaItem.Name,
				Price:      greenTeaItem.Price,
				Quantity:   4,
				Subtotal:   greenTeaItem.Price * 4,
			},
			{
				MenuItemID: milkTeaItem.ID,
				Name:       milkTeaItem.Name,
				Price:      milkTeaItem.Price,
				Quantity:   3,
				Subtotal:   milkTeaItem.Price * 3,
			},
		},
		Subtotal:  greenTeaItem.Price*4 + milkTeaItem.Price*3,
		Total:     greenTeaItem.Price*4 + milkTeaItem.Price*3,
		Status:    order.StatusServed,
		CreatedAt: time.Now().Add(-2 * time.Hour),
		UpdatedAt: time.Now().Add(-2 * time.Hour),
	}
	orderRepo.orders[order2.ID] = order2

	// Order 3: Mixed items
	order3 := &order.Order{
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
			{
				MenuItemID: greenTeaItem.ID,
				Name:       greenTeaItem.Name,
				Price:      greenTeaItem.Price,
				Quantity:   1,
				Subtotal:   greenTeaItem.Price,
			},
		},
		Subtotal:  cappuccinoItem.Price*2 + greenTeaItem.Price,
		Total:     cappuccinoItem.Price*2 + greenTeaItem.Price,
		Status:    order.StatusServed,
		CreatedAt: time.Now().Add(-1 * time.Hour),
		UpdatedAt: time.Now().Add(-1 * time.Hour),
	}
	orderRepo.orders[order3.ID] = order3

	// Step 5: Close shift and calculate costs
	t.Log("Closing shift and calculating costs...")
	summary, err := costCalculator.CalculateShiftOrderCosts(ctx, testShift.ID)
	require.NoError(t, err)
	require.NotNil(t, summary)
	
	t.Logf("Cost calculation summary: %d orders, %d items, total cost: %.2f",
		summary.TotalOrders, summary.TotalItems, summary.TotalAccountingCost)

	// Step 6: View category profit report
	t.Log("Viewing category profit report...")
	
	// For the test, we'll use a date range that should capture all items
	// Note: The mock's FindByDateRange implementation checks CreatedAt,
	// but order items created by CalculateShiftOrderCosts have zero CreatedAt
	// So we need to manually set CreatedAt on the order items
	now := time.Now()
	for _, item := range orderItemRepo.orderItems {
		item.CreatedAt = now.Add(-2 * time.Hour) // Set to a time within our range
	}
	
	dateRange := DateRange{
		Start: now.Add(-5 * time.Hour),
		End:   now,
	}
	
	categoryProfits, err := profitAnalyzer.GetCategoryProfits(ctx, dateRange)
	require.NoError(t, err)
	require.NotEmpty(t, categoryProfits, "Should have category profit data")

	// Verify category profit data
	var coffeeProfitData, teaProfitData *CategoryProfit
	for i := range categoryProfits {
		if categoryProfits[i].Category == "Coffee" {
			coffeeProfitData = &categoryProfits[i]
		} else if categoryProfits[i].Category == "Tea" {
			teaProfitData = &categoryProfits[i]
		}
	}

	require.NotNil(t, coffeeProfitData, "Should have Coffee category profit data")
	require.NotNil(t, teaProfitData, "Should have Tea category profit data")

	// Verify Coffee category
	// 3 Cappuccino + 2 Latte + 2 Cappuccino = 5 Cappuccino + 2 Latte
	// Revenue: 5*45000 + 2*50000 = 225000 + 100000 = 325000
	expectedCoffeeRevenue := 325000.0
	assert.InDelta(t, expectedCoffeeRevenue, coffeeProfitData.TotalRevenue, 1.0,
		"Coffee revenue should be %.2f", expectedCoffeeRevenue)
	
	// Cost: Cappuccino = 14550, Latte = 17300
	// 5*14550 + 2*17300 = 72750 + 34600 = 107350
	expectedCoffeeCost := 107350.0
	assert.InDelta(t, expectedCoffeeCost, coffeeProfitData.TotalCost, 100.0,
		"Coffee cost should be approximately %.2f", expectedCoffeeCost)
	
	expectedCoffeeProfit := expectedCoffeeRevenue - expectedCoffeeCost
	assert.InDelta(t, expectedCoffeeProfit, coffeeProfitData.TotalProfit, 100.0,
		"Coffee profit should be approximately %.2f", expectedCoffeeProfit)
	
	assert.Greater(t, coffeeProfitData.AverageProfitMargin, 0.0,
		"Coffee profit margin should be positive")
	
	t.Logf("✓ Coffee category: Revenue=%.2f, Cost=%.2f, Profit=%.2f, Margin=%.2f%%",
		coffeeProfitData.TotalRevenue, coffeeProfitData.TotalCost,
		coffeeProfitData.TotalProfit, coffeeProfitData.AverageProfitMargin)

	// Verify Tea category
	// 4 Green Tea + 3 Milk Tea + 1 Green Tea = 5 Green Tea + 3 Milk Tea
	// Revenue: 5*30000 + 3*35000 = 150000 + 105000 = 255000
	expectedTeaRevenue := 255000.0
	assert.InDelta(t, expectedTeaRevenue, teaProfitData.TotalRevenue, 1.0,
		"Tea revenue should be %.2f", expectedTeaRevenue)
	
	// Cost: Green Tea = 1000, Milk Tea = 6500
	// 5*1000 + 3*6500 = 5000 + 19500 = 24500
	expectedTeaCost := 24500.0
	assert.InDelta(t, expectedTeaCost, teaProfitData.TotalCost, 100.0,
		"Tea cost should be approximately %.2f", expectedTeaCost)
	
	expectedTeaProfit := expectedTeaRevenue - expectedTeaCost
	assert.InDelta(t, expectedTeaProfit, teaProfitData.TotalProfit, 100.0,
		"Tea profit should be approximately %.2f", expectedTeaProfit)
	
	t.Logf("✓ Tea category: Revenue=%.2f, Cost=%.2f, Profit=%.2f, Margin=%.2f%%",
		teaProfitData.TotalRevenue, teaProfitData.TotalCost,
		teaProfitData.TotalProfit, teaProfitData.AverageProfitMargin)

	// Step 7: Add operating expenses
	t.Log("Adding operating expenses...")
	
	periodStart := time.Now().Add(-24 * time.Hour)
	periodEnd := time.Now()
	
	operatingExpense := &expense.OperatingExpense{
		ID:             primitive.NewObjectID(),
		PeriodStart:    periodStart,
		PeriodEnd:      periodEnd,
		StaffSalary:    100000,
		Rent:           50000,
		Utilities:      20000,
		MarketingCosts: 10000,
		OtherExpenses:  5000,
		TotalExpenses:  185000,
		CreatedAt:      time.Now(),
		UpdatedAt:      time.Now(),
	}
	operatingExpenseRepo.expenses[operatingExpense.ID] = operatingExpense

	// Step 8: View operating profit report
	t.Log("Viewing operating profit report...")
	
	operatingProfitReport, err := profitAnalyzer.GetOperatingProfit(ctx, dateRange)
	require.NoError(t, err)
	require.NotNil(t, operatingProfitReport, "Should have operating profit report")

	// Verify operating profit report
	// Total revenue = Coffee + Tea = 325000 + 255000 = 580000
	expectedTotalRevenue := 580000.0
	assert.InDelta(t, expectedTotalRevenue, operatingProfitReport.TotalRevenue, 1.0,
		"Total revenue should be %.2f", expectedTotalRevenue)
	
	// Total COGS = Coffee cost + Tea cost = 107350 + 24500 = 131850
	expectedTotalCOGS := 131850.0
	assert.InDelta(t, expectedTotalCOGS, operatingProfitReport.TotalCOGS, 100.0,
		"Total COGS should be approximately %.2f", expectedTotalCOGS)
	
	// Gross profit = Revenue - COGS = 580000 - 131850 = 448150
	expectedGrossProfit := expectedTotalRevenue - expectedTotalCOGS
	assert.InDelta(t, expectedGrossProfit, operatingProfitReport.GrossProfit, 100.0,
		"Gross profit should be approximately %.2f", expectedGrossProfit)
	
	// Verify operating expenses are allocated proportionally
	// Period is 24 hours, date range is 5 hours, so expenses are prorated: 185000 * (5/24) = 38541.67
	expectedAllocatedExpenses := 185000.0 * (5.0 / 24.0)
	assert.InDelta(t, expectedAllocatedExpenses, operatingProfitReport.TotalExpenses, 1.0,
		"Total expenses should be allocated proportionally to %.2f", expectedAllocatedExpenses)
	
	// Verify expense_allocated flag is set
	assert.True(t, operatingProfitReport.ExpenseAllocated, "Expense allocated flag should be true")
	
	// Operating profit = Gross profit - Allocated expenses
	expectedOperatingProfit := expectedGrossProfit - expectedAllocatedExpenses
	assert.InDelta(t, expectedOperatingProfit, operatingProfitReport.OperatingProfit, 100.0,
		"Operating profit should be approximately %.2f", expectedOperatingProfit)
	
	// Operating profit margin = (Operating profit / Revenue) * 100
	expectedOperatingProfitMargin := (expectedOperatingProfit / expectedTotalRevenue) * 100
	assert.InDelta(t, expectedOperatingProfitMargin, operatingProfitReport.OperatingProfitMargin, 1.0,
		"Operating profit margin should be approximately %.2f%%", expectedOperatingProfitMargin)
	
	t.Logf("✓ Operating Profit Report:")
	t.Logf("  - Total Revenue: %.2f", operatingProfitReport.TotalRevenue)
	t.Logf("  - Total COGS: %.2f", operatingProfitReport.TotalCOGS)
	t.Logf("  - Gross Profit: %.2f (%.2f%%)", operatingProfitReport.GrossProfit, operatingProfitReport.GrossProfitMargin)
	t.Logf("  - Total Expenses: %.2f", operatingProfitReport.TotalExpenses)
	t.Logf("  - Operating Profit: %.2f (%.2f%%)", operatingProfitReport.OperatingProfit, operatingProfitReport.OperatingProfitMargin)

	t.Log("✅ Complete profit analysis workflow test passed!")
	t.Log("   - Orders created with various menu items")
	t.Log("   - Shift closed and costs calculated")
	t.Log("   - Category profit report generated correctly")
	t.Log("   - Operating expenses added")
	t.Log("   - Operating profit report generated correctly")
}
