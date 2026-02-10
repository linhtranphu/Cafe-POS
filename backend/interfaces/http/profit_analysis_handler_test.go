package http

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"cafe-pos/backend/application/services"
	"cafe-pos/backend/domain/expense"
	"cafe-pos/backend/domain/menu"
	"cafe-pos/backend/domain/order"
	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

// Mock repositories for testing

type mockMenuRepo struct {
	items []*menu.MenuItem
}

func (m *mockMenuRepo) Create(ctx context.Context, item *menu.MenuItem) error {
	m.items = append(m.items, item)
	return nil
}

func (m *mockMenuRepo) FindAll(ctx context.Context) ([]*menu.MenuItem, error) {
	return m.items, nil
}

func (m *mockMenuRepo) FindByID(ctx context.Context, id primitive.ObjectID) (*menu.MenuItem, error) {
	for _, item := range m.items {
		if item.ID == id {
			return item, nil
		}
	}
	return nil, nil
}

func (m *mockMenuRepo) FindByCategory(ctx context.Context, category string) ([]*menu.MenuItem, error) {
	var result []*menu.MenuItem
	for _, item := range m.items {
		if item.Category == category {
			result = append(result, item)
		}
	}
	return result, nil
}

func (m *mockMenuRepo) FindByIngredientName(ctx context.Context, ingredientName string) ([]*menu.MenuItem, error) {
	return nil, nil
}

func (m *mockMenuRepo) Update(ctx context.Context, id primitive.ObjectID, item *menu.MenuItem) error {
	return nil
}

func (m *mockMenuRepo) Delete(ctx context.Context, id primitive.ObjectID) error {
	return nil
}

type mockOrderItemRepo struct {
	items []*order.OrderItemWithCost
}

func (m *mockOrderItemRepo) CreateMany(ctx context.Context, items []*order.OrderItemWithCost) error {
	m.items = append(m.items, items...)
	return nil
}

func (m *mockOrderItemRepo) FindByDateRange(ctx context.Context, startDate, endDate time.Time) ([]*order.OrderItemWithCost, error) {
	var result []*order.OrderItemWithCost
	for _, item := range m.items {
		if item.CreatedAt.After(startDate) && item.CreatedAt.Before(endDate) {
			result = append(result, item)
		}
	}
	return result, nil
}

func (m *mockOrderItemRepo) FindByOrderIDs(ctx context.Context, orderIDs []primitive.ObjectID) ([]*order.OrderItemWithCost, error) {
	return m.items, nil
}

type mockOperatingExpenseRepo struct {
	expenses []*expense.OperatingExpense
}

func (m *mockOperatingExpenseRepo) Create(ctx context.Context, operatingExpense *expense.OperatingExpense) error {
	m.expenses = append(m.expenses, operatingExpense)
	return nil
}

func (m *mockOperatingExpenseRepo) Update(ctx context.Context, id primitive.ObjectID, operatingExpense *expense.OperatingExpense) error {
	return nil
}

func (m *mockOperatingExpenseRepo) FindByID(ctx context.Context, id primitive.ObjectID) (*expense.OperatingExpense, error) {
	return nil, nil
}

func (m *mockOperatingExpenseRepo) FindByPeriod(ctx context.Context, start, end time.Time) ([]*expense.OperatingExpense, error) {
	var result []*expense.OperatingExpense
	for _, exp := range m.expenses {
		// Check if expense period overlaps with date range
		if exp.PeriodEnd.After(start) && exp.PeriodStart.Before(end) {
			result = append(result, exp)
		}
	}
	return result, nil
}

func (m *mockOperatingExpenseRepo) FindForDate(ctx context.Context, date time.Time) (*expense.OperatingExpense, error) {
	return nil, nil
}

func (m *mockOperatingExpenseRepo) FindAll(ctx context.Context, startDate, endDate *time.Time) ([]*expense.OperatingExpense, error) {
	return m.expenses, nil
}

func (m *mockOperatingExpenseRepo) Upsert(ctx context.Context, operatingExpense *expense.OperatingExpense) (*expense.OperatingExpense, error) {
	return operatingExpense, nil
}

func (m *mockOperatingExpenseRepo) FindByDateRange(ctx context.Context, start, end time.Time) ([]*expense.OperatingExpense, error) {
	return m.FindByPeriod(ctx, start, end)
}

func (m *mockOperatingExpenseRepo) FindByDate(ctx context.Context, date time.Time) (*expense.OperatingExpense, error) {
	return m.FindForDate(ctx, date)
}

// Test GET /api/reports/category-profit with various date ranges
// Requirements: 6.1
func TestGetCategoryProfit_WithDateRanges(t *testing.T) {
	gin.SetMode(gin.TestMode)

	// Setup test data
	menuItem1 := &menu.MenuItem{
		ID:          primitive.NewObjectID(),
		Name:        "Cappuccino",
		Category:    "Coffee",
		Price:       50000,
		CurrentCost: 15000,
		CostStatus:  menu.CostStatusFinal,
	}

	menuItem2 := &menu.MenuItem{
		ID:          primitive.NewObjectID(),
		Name:        "Green Tea",
		Category:    "Tea",
		Price:       30000,
		CurrentCost: 10000,
		CostStatus:  menu.CostStatusFinal,
	}

	// Create order items with different dates
	baseDate := time.Date(2024, 1, 15, 12, 0, 0, 0, time.UTC)
	
	orderItem1 := &order.OrderItemWithCost{
		ID:             primitive.NewObjectID(),
		OrderID:        primitive.NewObjectID(),
		MenuItemID:     menuItem1.ID,
		Name:           menuItem1.Name,
		Price:          menuItem1.Price,
		Quantity:       2,
		AccountingCost: 30000,
		CostStatus:     order.CostStatusFinal,
		CreatedAt:      baseDate.AddDate(0, 0, 1), // Jan 16
	}

	orderItem2 := &order.OrderItemWithCost{
		ID:             primitive.NewObjectID(),
		OrderID:        primitive.NewObjectID(),
		MenuItemID:     menuItem2.ID,
		Name:           menuItem2.Name,
		Price:          menuItem2.Price,
		Quantity:       1,
		AccountingCost: 10000,
		CostStatus:     order.CostStatusFinal,
		CreatedAt:      baseDate.AddDate(0, 0, 5), // Jan 20
	}

	orderItem3 := &order.OrderItemWithCost{
		ID:             primitive.NewObjectID(),
		OrderID:        primitive.NewObjectID(),
		MenuItemID:     menuItem1.ID,
		Name:           menuItem1.Name,
		Price:          menuItem1.Price,
		Quantity:       1,
		AccountingCost: 15000,
		CostStatus:     order.CostStatusFinal,
		CreatedAt:      baseDate.AddDate(0, 0, 10), // Jan 25
	}

	// Setup mock repositories
	menuRepo := &mockMenuRepo{
		items: []*menu.MenuItem{menuItem1, menuItem2},
	}

	orderItemRepo := &mockOrderItemRepo{
		items: []*order.OrderItemWithCost{orderItem1, orderItem2, orderItem3},
	}

	// Create service and handler
	profitAnalyzerService := services.NewProfitAnalyzerService(menuRepo, orderItemRepo, nil, nil)
	handler := NewProfitAnalysisHandler(profitAnalyzerService)

	// Test case 1: Date range includes all items (Jan 15 - Jan 31)
	t.Run("Date range includes all items", func(t *testing.T) {
		w := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(w)
		c.Request = httptest.NewRequest("GET", "/api/reports/category-profit?start_date=2024-01-15&end_date=2024-01-31", nil)

		handler.GetCategoryProfit(c)

		assert.Equal(t, http.StatusOK, w.Code)

		var response map[string]interface{}
		err := json.Unmarshal(w.Body.Bytes(), &response)
		assert.NoError(t, err)

		categories := response["categories"].([]interface{})
		assert.Equal(t, 2, len(categories)) // Coffee and Tea

		// Verify Coffee category
		var coffeeCategory map[string]interface{}
		for _, cat := range categories {
			catMap := cat.(map[string]interface{})
			if catMap["category"] == "Coffee" {
				coffeeCategory = catMap
				break
			}
		}
		assert.NotNil(t, coffeeCategory)
		assert.Equal(t, 150000.0, coffeeCategory["total_revenue"]) // 2*50000 + 1*50000
		assert.Equal(t, 45000.0, coffeeCategory["total_cost"])     // 30000 + 15000
	})

	// Test case 2: Date range includes only first two items (Jan 15 - Jan 22)
	t.Run("Date range includes only first two items", func(t *testing.T) {
		w := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(w)
		c.Request = httptest.NewRequest("GET", "/api/reports/category-profit?start_date=2024-01-15&end_date=2024-01-22", nil)

		handler.GetCategoryProfit(c)

		assert.Equal(t, http.StatusOK, w.Code)

		var response map[string]interface{}
		err := json.Unmarshal(w.Body.Bytes(), &response)
		assert.NoError(t, err)

		categories := response["categories"].([]interface{})
		assert.Equal(t, 2, len(categories)) // Coffee and Tea

		// Verify Coffee category
		var coffeeCategory map[string]interface{}
		for _, cat := range categories {
			catMap := cat.(map[string]interface{})
			if catMap["category"] == "Coffee" {
				coffeeCategory = catMap
				break
			}
		}
		assert.NotNil(t, coffeeCategory)
		assert.Equal(t, 100000.0, coffeeCategory["total_revenue"]) // 2*50000
		assert.Equal(t, 30000.0, coffeeCategory["total_cost"])     // 30000
	})

	// Test case 3: Date range includes no items (Jan 1 - Jan 10)
	t.Run("Date range includes no items", func(t *testing.T) {
		w := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(w)
		c.Request = httptest.NewRequest("GET", "/api/reports/category-profit?start_date=2024-01-01&end_date=2024-01-10", nil)

		handler.GetCategoryProfit(c)

		assert.Equal(t, http.StatusOK, w.Code)

		var response map[string]interface{}
		err := json.Unmarshal(w.Body.Bytes(), &response)
		assert.NoError(t, err)

		categories := response["categories"].([]interface{})
		assert.Equal(t, 0, len(categories)) // No categories
	})
}

// Test GET /api/reports/category-profit with invalid date ranges
func TestGetCategoryProfit_InvalidDateRanges(t *testing.T) {
	gin.SetMode(gin.TestMode)

	// Create service and handler with empty repos
	profitAnalyzerService := services.NewProfitAnalyzerService(&mockMenuRepo{}, &mockOrderItemRepo{}, nil, nil)
	handler := NewProfitAnalysisHandler(profitAnalyzerService)

	// Test case 1: Missing start_date
	t.Run("Missing start_date", func(t *testing.T) {
		w := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(w)
		c.Request = httptest.NewRequest("GET", "/api/reports/category-profit?end_date=2024-01-31", nil)

		handler.GetCategoryProfit(c)

		assert.Equal(t, http.StatusBadRequest, w.Code)
	})

	// Test case 2: Missing end_date
	t.Run("Missing end_date", func(t *testing.T) {
		w := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(w)
		c.Request = httptest.NewRequest("GET", "/api/reports/category-profit?start_date=2024-01-01", nil)

		handler.GetCategoryProfit(c)

		assert.Equal(t, http.StatusBadRequest, w.Code)
	})

	// Test case 3: Invalid date format
	t.Run("Invalid date format", func(t *testing.T) {
		w := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(w)
		c.Request = httptest.NewRequest("GET", "/api/reports/category-profit?start_date=2024/01/01&end_date=2024-01-31", nil)

		handler.GetCategoryProfit(c)

		assert.Equal(t, http.StatusBadRequest, w.Code)
	})

	// Test case 4: start_date after end_date
	t.Run("start_date after end_date", func(t *testing.T) {
		w := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(w)
		c.Request = httptest.NewRequest("GET", "/api/reports/category-profit?start_date=2024-01-31&end_date=2024-01-01", nil)

		handler.GetCategoryProfit(c)

		assert.Equal(t, http.StatusBadRequest, w.Code)
	})
}

// Test GET /api/reports/operating-profit with and without expenses
// Requirements: 6.5.1
func TestGetOperatingProfit_WithAndWithoutExpenses(t *testing.T) {
	gin.SetMode(gin.TestMode)

	// Setup test data
	menuItem1 := &menu.MenuItem{
		ID:          primitive.NewObjectID(),
		Name:        "Cappuccino",
		Category:    "Coffee",
		Price:       50000,
		CurrentCost: 15000,
		CostStatus:  menu.CostStatusFinal,
	}

	// Create order items
	baseDate := time.Date(2024, 1, 15, 12, 0, 0, 0, time.UTC)
	
	orderItem1 := &order.OrderItemWithCost{
		ID:             primitive.NewObjectID(),
		OrderID:        primitive.NewObjectID(),
		MenuItemID:     menuItem1.ID,
		Name:           menuItem1.Name,
		Price:          menuItem1.Price,
		Quantity:       2,
		AccountingCost: 30000,
		CostStatus:     order.CostStatusFinal,
		CreatedAt:      baseDate.AddDate(0, 0, 1),
	}

	// Setup mock repositories
	menuRepo := &mockMenuRepo{
		items: []*menu.MenuItem{menuItem1},
	}

	orderItemRepo := &mockOrderItemRepo{
		items: []*order.OrderItemWithCost{orderItem1},
	}

	// Test case 1: Without expenses
	t.Run("Without expenses", func(t *testing.T) {
		operatingExpenseRepo := &mockOperatingExpenseRepo{
			expenses: []*expense.OperatingExpense{},
		}

		profitAnalyzerService := services.NewProfitAnalyzerService(menuRepo, orderItemRepo, nil, operatingExpenseRepo)
		handler := NewProfitAnalysisHandler(profitAnalyzerService)

		w := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(w)
		c.Request = httptest.NewRequest("GET", "/api/reports/operating-profit?start_date=2024-01-15&end_date=2024-01-31", nil)

		handler.GetOperatingProfit(c)

		assert.Equal(t, http.StatusOK, w.Code)

		var response services.OperatingProfitReport
		err := json.Unmarshal(w.Body.Bytes(), &response)
		assert.NoError(t, err)

		// Verify gross profit
		assert.Equal(t, 100000.0, response.TotalRevenue) // 2*50000
		assert.Equal(t, 30000.0, response.TotalCOGS)
		assert.Equal(t, 70000.0, response.GrossProfit)

		// Verify operating profit equals gross profit when no expenses
		assert.Equal(t, 70000.0, response.OperatingProfit)
		assert.Equal(t, 0.0, response.TotalExpenses)
		assert.Contains(t, response.AllocationNote, "Chưa nhập chi phí vận hành")
	})

	// Test case 2: With expenses
	t.Run("With expenses", func(t *testing.T) {
		operatingExpense := &expense.OperatingExpense{
			ID:             primitive.NewObjectID(),
			PeriodStart:    time.Date(2024, 1, 1, 0, 0, 0, 0, time.UTC),
			PeriodEnd:      time.Date(2024, 1, 31, 23, 59, 59, 0, time.UTC),
			StaffSalary:    2000000,
			Rent:           1000000,
			Utilities:      500000,
			MarketingCosts: 300000,
			OtherExpenses:  200000,
			TotalExpenses:  4000000,
		}

		operatingExpenseRepo := &mockOperatingExpenseRepo{
			expenses: []*expense.OperatingExpense{operatingExpense},
		}

		profitAnalyzerService := services.NewProfitAnalyzerService(menuRepo, orderItemRepo, nil, operatingExpenseRepo)
		handler := NewProfitAnalysisHandler(profitAnalyzerService)

		w := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(w)
		c.Request = httptest.NewRequest("GET", "/api/reports/operating-profit?start_date=2024-01-15&end_date=2024-01-31", nil)

		handler.GetOperatingProfit(c)

		assert.Equal(t, http.StatusOK, w.Code)

		var response services.OperatingProfitReport
		err := json.Unmarshal(w.Body.Bytes(), &response)
		assert.NoError(t, err)

		// Verify gross profit
		assert.Equal(t, 100000.0, response.TotalRevenue)
		assert.Equal(t, 30000.0, response.TotalCOGS)
		assert.Equal(t, 70000.0, response.GrossProfit)

		// Verify expenses are included
		assert.Greater(t, response.TotalExpenses, 0.0)
		assert.Equal(t, response.GrossProfit-response.TotalExpenses, response.OperatingProfit)
	})
}

// Test expense allocation scenarios
// Requirements: 6.5.8
func TestGetOperatingProfit_ExpenseAllocation(t *testing.T) {
	gin.SetMode(gin.TestMode)

	// Setup test data
	menuItem1 := &menu.MenuItem{
		ID:          primitive.NewObjectID(),
		Name:        "Cappuccino",
		Category:    "Coffee",
		Price:       50000,
		CurrentCost: 15000,
		CostStatus:  menu.CostStatusFinal,
	}

	// Create order items
	baseDate := time.Date(2024, 1, 15, 12, 0, 0, 0, time.UTC)
	
	orderItem1 := &order.OrderItemWithCost{
		ID:             primitive.NewObjectID(),
		OrderID:        primitive.NewObjectID(),
		MenuItemID:     menuItem1.ID,
		Name:           menuItem1.Name,
		Price:          menuItem1.Price,
		Quantity:       2,
		AccountingCost: 30000,
		CostStatus:     order.CostStatusFinal,
		CreatedAt:      baseDate.AddDate(0, 0, 1),
	}

	// Setup mock repositories
	menuRepo := &mockMenuRepo{
		items: []*menu.MenuItem{menuItem1},
	}

	orderItemRepo := &mockOrderItemRepo{
		items: []*order.OrderItemWithCost{orderItem1},
	}

	// Test case: Monthly expense allocated to partial month
	t.Run("Monthly expense allocated to partial month", func(t *testing.T) {
		// Monthly expense for entire January
		operatingExpense := &expense.OperatingExpense{
			ID:             primitive.NewObjectID(),
			PeriodStart:    time.Date(2024, 1, 1, 0, 0, 0, 0, time.UTC),
			PeriodEnd:      time.Date(2024, 1, 31, 23, 59, 59, 0, time.UTC),
			StaffSalary:    3100000, // 100k per day for 31 days
			Rent:           1000000,
			Utilities:      500000,
			MarketingCosts: 300000,
			OtherExpenses:  200000,
			TotalExpenses:  5100000,
		}

		operatingExpenseRepo := &mockOperatingExpenseRepo{
			expenses: []*expense.OperatingExpense{operatingExpense},
		}

		profitAnalyzerService := services.NewProfitAnalyzerService(menuRepo, orderItemRepo, nil, operatingExpenseRepo)
		handler := NewProfitAnalysisHandler(profitAnalyzerService)

		// Query for only Jan 15-20 (6 days)
		w := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(w)
		c.Request = httptest.NewRequest("GET", "/api/reports/operating-profit?start_date=2024-01-15&end_date=2024-01-20", nil)

		handler.GetOperatingProfit(c)

		assert.Equal(t, http.StatusOK, w.Code)

		var response services.OperatingProfitReport
		err := json.Unmarshal(w.Body.Bytes(), &response)
		assert.NoError(t, err)

		// Verify expense allocation flag is set
		assert.True(t, response.ExpenseAllocated)
		assert.Contains(t, response.AllocationNote, "Chi phí được phân bổ từ tháng")

		// Verify expenses are allocated proportionally
		// Should be approximately 6/31 of the monthly expense
		assert.Greater(t, response.TotalExpenses, 0.0)
		assert.Less(t, response.TotalExpenses, operatingExpense.TotalExpenses)
	})
}
