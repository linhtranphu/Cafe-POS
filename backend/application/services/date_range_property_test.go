package services

import (
	"context"
	"testing"
	"time"

	"cafe-pos/backend/domain/expense"
	"cafe-pos/backend/domain/menu"
	"cafe-pos/backend/domain/order"
	"github.com/leanovate/gopter"
	"github.com/leanovate/gopter/gen"
	"github.com/leanovate/gopter/prop"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

// Feature: menu-cost-profit-analysis, Property 16: Date Range Filtering
// **Validates: Requirements 6.4, 6.5.6**
//
// Property: For any date range, the category profit analysis and operating profit analysis
// should include only orders where the order date falls within the specified range.
func TestProperty_DateRangeFiltering(t *testing.T) {
	parameters := gopter.DefaultTestParameters()
	parameters.MinSuccessfulTests = 100 // Minimum 100 iterations as per spec
	properties := gopter.NewProperties(parameters)

	properties.Property("Only orders within date range are included in profit analysis", prop.ForAll(
		func(daysOffset1, daysOffset2, daysOffset3 int) bool {
			// Create a base date
			baseDate := time.Date(2024, 1, 15, 12, 0, 0, 0, time.UTC)

			// Create three dates from offsets (can be negative or positive)
			date1 := baseDate.AddDate(0, 0, daysOffset1)
			date2 := baseDate.AddDate(0, 0, daysOffset2)
			date3 := baseDate.AddDate(0, 0, daysOffset3)

			// Determine the date range (start and end)
			var startDate, endDate time.Time
			if date1.Before(date2) {
				startDate = date1
				endDate = date2
			} else {
				startDate = date2
				endDate = date1
			}

			// Ensure we have a valid range (at least 1 day)
			if startDate.Equal(endDate) {
				endDate = startDate.AddDate(0, 0, 1)
			}

			// Create menu items
			menuItem1 := &menu.MenuItem{
				ID:          primitive.NewObjectID(),
				Name:        "Test Item 1",
				Category:    "Coffee",
				Price:       50000,
				CurrentCost: 15000,
				CostStatus:  menu.CostStatusFinal,
			}

			// Create order items with different dates
			// Item 1: Within range
			orderItem1 := &order.OrderItemWithCost{
				ID:             primitive.NewObjectID(),
				OrderID:        primitive.NewObjectID(),
				MenuItemID:     menuItem1.ID,
				Name:           menuItem1.Name,
				Price:          menuItem1.Price,
				Quantity:       2,
				AccountingCost: 30000, // 15000 * 2
				CostStatus:     order.CostStatusFinal,
				CreatedAt:      startDate.AddDate(0, 0, 1), // 1 day after start
			}

			// Item 2: Before range
			orderItem2 := &order.OrderItemWithCost{
				ID:             primitive.NewObjectID(),
				OrderID:        primitive.NewObjectID(),
				MenuItemID:     menuItem1.ID,
				Name:           menuItem1.Name,
				Price:          menuItem1.Price,
				Quantity:       1,
				AccountingCost: 15000,
				CostStatus:     order.CostStatusFinal,
				CreatedAt:      startDate.AddDate(0, 0, -5), // 5 days before start
			}

			// Item 3: After range
			orderItem3 := &order.OrderItemWithCost{
				ID:             primitive.NewObjectID(),
				OrderID:        primitive.NewObjectID(),
				MenuItemID:     menuItem1.ID,
				Name:           menuItem1.Name,
				Price:          menuItem1.Price,
				Quantity:       1,
				AccountingCost: 15000,
				CostStatus:     order.CostStatusFinal,
				CreatedAt:      endDate.AddDate(0, 0, 5), // 5 days after end
			}

			// Item 4: At date3 (may or may not be in range)
			orderItem4 := &order.OrderItemWithCost{
				ID:             primitive.NewObjectID(),
				OrderID:        primitive.NewObjectID(),
				MenuItemID:     menuItem1.ID,
				Name:           menuItem1.Name,
				Price:          menuItem1.Price,
				Quantity:       3,
				AccountingCost: 45000, // 15000 * 3
				CostStatus:     order.CostStatusFinal,
				CreatedAt:      date3,
			}

			// Setup mock repositories
			menuRepo := &mockMenuRepo{
				items: []*menu.MenuItem{menuItem1},
			}

			orderItemRepo := &mockOrderItemRepo{
				items: []*order.OrderItemWithCost{orderItem1, orderItem2, orderItem3, orderItem4},
			}

			// Create service
			service := NewProfitAnalyzerService(menuRepo, orderItemRepo, nil, nil)

			// Set time to start and end of day for the range
			rangeStart := time.Date(startDate.Year(), startDate.Month(), startDate.Day(), 0, 0, 0, 0, startDate.Location())
			rangeEnd := time.Date(endDate.Year(), endDate.Month(), endDate.Day(), 23, 59, 59, 999999999, endDate.Location())

			// Call GetCategoryProfits with date range
			dateRange := DateRange{
				Start: rangeStart,
				End:   rangeEnd,
			}
			categories, err := service.GetCategoryProfits(context.Background(), dateRange)
			if err != nil {
				t.Logf("Unexpected error: %v", err)
				return false
			}

			// Calculate expected revenue and cost
			// Only items within the date range should be included
			expectedRevenue := 0.0
			expectedCost := 0.0

			// Check each order item
			if orderItem1.CreatedAt.After(rangeStart) && orderItem1.CreatedAt.Before(rangeEnd) {
				expectedRevenue += orderItem1.Price * float64(orderItem1.Quantity)
				expectedCost += orderItem1.AccountingCost
			}
			if orderItem2.CreatedAt.After(rangeStart) && orderItem2.CreatedAt.Before(rangeEnd) {
				expectedRevenue += orderItem2.Price * float64(orderItem2.Quantity)
				expectedCost += orderItem2.AccountingCost
			}
			if orderItem3.CreatedAt.After(rangeStart) && orderItem3.CreatedAt.Before(rangeEnd) {
				expectedRevenue += orderItem3.Price * float64(orderItem3.Quantity)
				expectedCost += orderItem3.AccountingCost
			}
			if orderItem4.CreatedAt.After(rangeStart) && orderItem4.CreatedAt.Before(rangeEnd) {
				expectedRevenue += orderItem4.Price * float64(orderItem4.Quantity)
				expectedCost += orderItem4.AccountingCost
			}

			// If no items in range, categories should be empty
			if expectedRevenue == 0 {
				return len(categories) == 0
			}

			// Find the Coffee category
			var coffeeCategory *CategoryProfit
			for i := range categories {
				if categories[i].Category == "Coffee" {
					coffeeCategory = &categories[i]
					break
				}
			}

			// Should have found the category
			if coffeeCategory == nil {
				t.Logf("Coffee category not found in results")
				return false
			}

			// Verify revenue and cost match expected
			tolerance := 0.01
			if coffeeCategory.TotalRevenue < expectedRevenue-tolerance || coffeeCategory.TotalRevenue > expectedRevenue+tolerance {
				t.Logf("Revenue mismatch: expected %v, got %v", expectedRevenue, coffeeCategory.TotalRevenue)
				return false
			}

			if coffeeCategory.TotalCost < expectedCost-tolerance || coffeeCategory.TotalCost > expectedCost+tolerance {
				t.Logf("Cost mismatch: expected %v, got %v", expectedCost, coffeeCategory.TotalCost)
				return false
			}

			return true
		},
		gen.IntRange(-30, 30), // daysOffset1: -30 to +30 days from base
		gen.IntRange(-30, 30), // daysOffset2: -30 to +30 days from base
		gen.IntRange(-30, 30), // daysOffset3: -30 to +30 days from base
	))

	properties.TestingRun(t, gopter.ConsoleReporter(false))
}

// Feature: menu-cost-profit-analysis, Property 16: Date Range Filtering (Operating Profit)
// **Validates: Requirements 6.5.6**
//
// Property: For any date range, the operating profit analysis should include only orders
// where the order date falls within the specified range.
func TestProperty_DateRangeFiltering_OperatingProfit(t *testing.T) {
	parameters := gopter.DefaultTestParameters()
	parameters.MinSuccessfulTests = 100 // Minimum 100 iterations as per spec
	properties := gopter.NewProperties(parameters)

	properties.Property("Only orders within date range are included in operating profit", prop.ForAll(
		func(daysOffset1, daysOffset2 int) bool {
			// Create a base date
			baseDate := time.Date(2024, 1, 15, 12, 0, 0, 0, time.UTC)

			// Create two dates from offsets
			date1 := baseDate.AddDate(0, 0, daysOffset1)
			date2 := baseDate.AddDate(0, 0, daysOffset2)

			// Determine the date range (start and end)
			var startDate, endDate time.Time
			if date1.Before(date2) {
				startDate = date1
				endDate = date2
			} else {
				startDate = date2
				endDate = date1
			}

			// Ensure we have a valid range (at least 1 day)
			if startDate.Equal(endDate) {
				endDate = startDate.AddDate(0, 0, 1)
			}

			// Create menu items
			menuItem1 := &menu.MenuItem{
				ID:          primitive.NewObjectID(),
				Name:        "Test Item 1",
				Category:    "Coffee",
				Price:       50000,
				CurrentCost: 15000,
				CostStatus:  menu.CostStatusFinal,
			}

			// Create order items with different dates
			// Item 1: Within range
			orderItem1 := &order.OrderItemWithCost{
				ID:             primitive.NewObjectID(),
				OrderID:        primitive.NewObjectID(),
				MenuItemID:     menuItem1.ID,
				Name:           menuItem1.Name,
				Price:          menuItem1.Price,
				Quantity:       2,
				AccountingCost: 30000,
				CostStatus:     order.CostStatusFinal,
				CreatedAt:      startDate.AddDate(0, 0, 1), // 1 day after start
			}

			// Item 2: Before range
			orderItem2 := &order.OrderItemWithCost{
				ID:             primitive.NewObjectID(),
				OrderID:        primitive.NewObjectID(),
				MenuItemID:     menuItem1.ID,
				Name:           menuItem1.Name,
				Price:          menuItem1.Price,
				Quantity:       1,
				AccountingCost: 15000,
				CostStatus:     order.CostStatusFinal,
				CreatedAt:      startDate.AddDate(0, 0, -5), // 5 days before start
			}

			// Setup mock repositories
			menuRepo := &mockMenuRepo{
				items: []*menu.MenuItem{menuItem1},
			}

			orderItemRepo := &mockOrderItemRepo{
				items: []*order.OrderItemWithCost{orderItem1, orderItem2},
			}

			// Create mock operating expense repo (empty, no expenses)
			operatingExpenseRepo := &mockOperatingExpenseRepo{
				expenses: []*expense.OperatingExpense{},
			}

			// Create service
			service := NewProfitAnalyzerService(menuRepo, orderItemRepo, nil, operatingExpenseRepo)

			// Set time to start and end of day for the range
			rangeStart := time.Date(startDate.Year(), startDate.Month(), startDate.Day(), 0, 0, 0, 0, startDate.Location())
			rangeEnd := time.Date(endDate.Year(), endDate.Month(), endDate.Day(), 23, 59, 59, 999999999, endDate.Location())

			// Call GetOperatingProfit with date range
			dateRange := DateRange{
				Start: rangeStart,
				End:   rangeEnd,
			}
			report, err := service.GetOperatingProfit(context.Background(), dateRange)
			if err != nil {
				t.Logf("Unexpected error: %v", err)
				return false
			}

			// Calculate expected revenue and cost
			// Only items within the date range should be included
			expectedRevenue := 0.0
			expectedCOGS := 0.0

			// Check each order item
			if orderItem1.CreatedAt.After(rangeStart) && orderItem1.CreatedAt.Before(rangeEnd) {
				expectedRevenue += orderItem1.Price * float64(orderItem1.Quantity)
				expectedCOGS += orderItem1.AccountingCost
			}
			if orderItem2.CreatedAt.After(rangeStart) && orderItem2.CreatedAt.Before(rangeEnd) {
				expectedRevenue += orderItem2.Price * float64(orderItem2.Quantity)
				expectedCOGS += orderItem2.AccountingCost
			}

			// Verify revenue and COGS match expected
			tolerance := 0.01
			if report.TotalRevenue < expectedRevenue-tolerance || report.TotalRevenue > expectedRevenue+tolerance {
				t.Logf("Revenue mismatch: expected %v, got %v", expectedRevenue, report.TotalRevenue)
				return false
			}

			if report.TotalCOGS < expectedCOGS-tolerance || report.TotalCOGS > expectedCOGS+tolerance {
				t.Logf("COGS mismatch: expected %v, got %v", expectedCOGS, report.TotalCOGS)
				return false
			}

			return true
		},
		gen.IntRange(-30, 30), // daysOffset1: -30 to +30 days from base
		gen.IntRange(-30, 30), // daysOffset2: -30 to +30 days from base
	))

	properties.TestingRun(t, gopter.ConsoleReporter(false))
}
