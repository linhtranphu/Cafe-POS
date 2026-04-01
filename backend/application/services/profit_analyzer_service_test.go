package services

import (
	"context"
	"testing"
	"time"

	"cafe-pos/backend/domain/expense"
	"cafe-pos/backend/domain/menu"
	"cafe-pos/backend/domain/order"
	"cafe-pos/backend/domain/settings"
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

func (m *mockMenuRepo) FindByBatchDefinitionID(ctx context.Context, batchDefID primitive.ObjectID) ([]*menu.MenuItem, error) {
	return nil, nil
}

type mockOrderItemRepo struct {
	items []*order.OrderItemWithCost
}

func (m *mockOrderItemRepo) CreateMany(ctx context.Context, items []*order.OrderItemWithCost) error {
	m.items = append(m.items, items...)
	return nil
}

func (m *mockOrderItemRepo) FindByOrderIDs(ctx context.Context, orderIDs []primitive.ObjectID) ([]*order.OrderItemWithCost, error) {
	return m.items, nil
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

type mockSettingsRepo struct {
	settings *settings.ShopSettings
}

func (m *mockSettingsRepo) FindFirst(ctx context.Context) (*settings.ShopSettings, error) {
	return m.settings, nil
}

func (m *mockSettingsRepo) Update(ctx context.Context, id primitive.ObjectID, settings *settings.ShopSettings) error {
	m.settings = settings
	return nil
}

type mockOperatingExpenseRepo struct {
	expenses []*expense.OperatingExpense
}

func (m *mockOperatingExpenseRepo) FindByDateRange(ctx context.Context, start, end time.Time) ([]*expense.OperatingExpense, error) {
	return m.FindByPeriod(ctx, start, end)
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

func (m *mockOperatingExpenseRepo) FindByDate(ctx context.Context, date time.Time) (*expense.OperatingExpense, error) {
	return m.FindForDate(ctx, date)
}

func (m *mockOperatingExpenseRepo) FindForDate(ctx context.Context, date time.Time) (*expense.OperatingExpense, error) {
	for _, exp := range m.expenses {
		if date.After(exp.PeriodStart) && date.Before(exp.PeriodEnd) {
			return exp, nil
		}
	}
	return nil, nil
}

func (m *mockOperatingExpenseRepo) FindByID(ctx context.Context, id primitive.ObjectID) (*expense.OperatingExpense, error) {
	for _, exp := range m.expenses {
		if exp.ID == id {
			return exp, nil
		}
	}
	return nil, nil
}

func (m *mockOperatingExpenseRepo) FindAll(ctx context.Context, startDate, endDate *time.Time) ([]*expense.OperatingExpense, error) {
	if startDate == nil || endDate == nil {
		return m.expenses, nil
	}
	return m.FindByPeriod(ctx, *startDate, *endDate)
}

func (m *mockOperatingExpenseRepo) Create(ctx context.Context, expense *expense.OperatingExpense) error {
	m.expenses = append(m.expenses, expense)
	return nil
}

func (m *mockOperatingExpenseRepo) Update(ctx context.Context, id primitive.ObjectID, expense *expense.OperatingExpense) error {
	return nil
}

func (m *mockOperatingExpenseRepo) Upsert(ctx context.Context, operatingExpense *expense.OperatingExpense) (*expense.OperatingExpense, error) {
	// Check if expense exists for this period
	for i, exp := range m.expenses {
		if exp.PeriodStart.Equal(operatingExpense.PeriodStart) && exp.PeriodEnd.Equal(operatingExpense.PeriodEnd) {
			m.expenses[i] = operatingExpense
			return operatingExpense, nil
		}
	}
	// Create new
	if operatingExpense.ID.IsZero() {
		operatingExpense.ID = primitive.NewObjectID()
	}
	m.expenses = append(m.expenses, operatingExpense)
	return operatingExpense, nil
}

// Test CalculateMenuItemProfit
func TestCalculateMenuItemProfit(t *testing.T) {
	// Setup
	menuRepo := &mockMenuRepo{
		items: []*menu.MenuItem{
			{
				ID:          primitive.NewObjectID(),
				Name:        "Cappuccino",
				Category:    "Coffee",
				Price:       45000,
				CurrentCost: 15000,
				CostStatus:  menu.CostStatusFinal,
			},
		},
	}

	service := NewProfitAnalyzerService(menuRepo, nil, nil, nil)

	// Test
	profit, err := service.CalculateMenuItemProfit(context.Background(), menuRepo.items[0].ID)

	// Verify
	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}

	if profit.ProfitMargin != 66.67 {
		t.Errorf("Expected profit margin 66.67, got %.2f", profit.ProfitMargin)
	}

	if profit.AbsoluteProfit != 30000 {
		t.Errorf("Expected absolute profit 30000, got %.2f", profit.AbsoluteProfit)
	}

	if profit.WarningStatus != WarningNone {
		t.Errorf("Expected warning status none, got %s", profit.WarningStatus)
	}
}

// Test DetectWarningStatus - Loss case
func TestDetectWarningStatus_Loss(t *testing.T) {
	// Setup - menu item with cost > price
	menuRepo := &mockMenuRepo{
		items: []*menu.MenuItem{
			{
				ID:          primitive.NewObjectID(),
				Name:        "Loss Item",
				Category:    "Coffee",
				Price:       20000,
				CurrentCost: 25000,
				CostStatus:  menu.CostStatusFinal,
			},
		},
	}

	settingsRepo := &mockSettingsRepo{
		settings: &settings.ShopSettings{
			LowMarginThreshold: 20.0,
		},
	}

	service := NewProfitAnalyzerService(menuRepo, nil, settingsRepo, nil)

	// Test
	warnings, err := service.DetectWarningStatus(context.Background(), 20.0)

	// Verify
	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}

	if warnings.LossCount != 1 {
		t.Errorf("Expected 1 loss item, got %d", warnings.LossCount)
	}

	if len(warnings.LossItems) != 1 {
		t.Errorf("Expected 1 loss item in array, got %d", len(warnings.LossItems))
	}

	if warnings.LossItems[0].WarningStatus != WarningLoss {
		t.Errorf("Expected warning status loss, got %s", warnings.LossItems[0].WarningStatus)
	}
}

// Test DetectWarningStatus - Low margin case
func TestDetectWarningStatus_LowMargin(t *testing.T) {
	// Setup - menu item with profit margin < threshold
	menuRepo := &mockMenuRepo{
		items: []*menu.MenuItem{
			{
				ID:          primitive.NewObjectID(),
				Name:        "Low Margin Item",
				Category:    "Coffee",
				Price:       30000,
				CurrentCost: 25000,
				CostStatus:  menu.CostStatusFinal,
			},
		},
	}

	settingsRepo := &mockSettingsRepo{
		settings: &settings.ShopSettings{
			LowMarginThreshold: 20.0,
		},
	}

	service := NewProfitAnalyzerService(menuRepo, nil, settingsRepo, nil)

	// Test
	warnings, err := service.DetectWarningStatus(context.Background(), 20.0)

	// Verify
	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}

	if warnings.LowMarginCount != 1 {
		t.Errorf("Expected 1 low margin item, got %d", warnings.LowMarginCount)
	}

	if len(warnings.LowMarginItems) != 1 {
		t.Errorf("Expected 1 low margin item in array, got %d", len(warnings.LowMarginItems))
	}

	if warnings.LowMarginItems[0].WarningStatus != WarningLowMargin {
		t.Errorf("Expected warning status low_margin, got %s", warnings.LowMarginItems[0].WarningStatus)
	}
}

// Test GetAllMenuItemProfits with filtering
func TestGetAllMenuItemProfits_CategoryFilter(t *testing.T) {
	// Setup
	menuRepo := &mockMenuRepo{
		items: []*menu.MenuItem{
			{
				ID:          primitive.NewObjectID(),
				Name:        "Cappuccino",
				Category:    "Coffee",
				Price:       45000,
				CurrentCost: 15000,
				CostStatus:  menu.CostStatusFinal,
			},
			{
				ID:          primitive.NewObjectID(),
				Name:        "Green Tea",
				Category:    "Tea",
				Price:       30000,
				CurrentCost: 10000,
				CostStatus:  menu.CostStatusFinal,
			},
		},
	}

	settingsRepo := &mockSettingsRepo{
		settings: &settings.ShopSettings{
			LowMarginThreshold: 20.0,
		},
	}

	service := NewProfitAnalyzerService(menuRepo, nil, settingsRepo, nil)

	// Test - filter by Coffee category
	filter := ProfitFilter{
		Category: "Coffee",
	}
	response, err := service.GetAllMenuItemProfits(context.Background(), filter)

	// Verify
	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}

	if len(response.Items) != 1 {
		t.Errorf("Expected 1 item, got %d", len(response.Items))
	}

	if response.Items[0].Category != "Coffee" {
		t.Errorf("Expected category Coffee, got %s", response.Items[0].Category)
	}
}

// Test GetAllMenuItemProfits with sorting
func TestGetAllMenuItemProfits_Sorting(t *testing.T) {
	// Setup
	menuRepo := &mockMenuRepo{
		items: []*menu.MenuItem{
			{
				ID:          primitive.NewObjectID(),
				Name:        "Low Profit",
				Category:    "Coffee",
				Price:       30000,
				CurrentCost: 25000,
				CostStatus:  menu.CostStatusFinal,
			},
			{
				ID:          primitive.NewObjectID(),
				Name:        "High Profit",
				Category:    "Coffee",
				Price:       45000,
				CurrentCost: 15000,
				CostStatus:  menu.CostStatusFinal,
			},
		},
	}

	settingsRepo := &mockSettingsRepo{
		settings: &settings.ShopSettings{
			LowMarginThreshold: 20.0,
		},
	}

	service := NewProfitAnalyzerService(menuRepo, nil, settingsRepo, nil)

	// Test - sort by profit_margin descending (default)
	filter := ProfitFilter{
		SortBy:    "profit_margin",
		SortOrder: "desc",
	}
	response, err := service.GetAllMenuItemProfits(context.Background(), filter)

	// Verify
	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}

	if len(response.Items) != 2 {
		t.Errorf("Expected 2 items, got %d", len(response.Items))
	}

	// First item should have higher profit margin
	if response.Items[0].ProfitMargin < response.Items[1].ProfitMargin {
		t.Errorf("Expected items sorted by profit margin descending")
	}
}

// Test edge case: price = 0 (promotional item)
func TestCalculateMenuItemProfit_ZeroPrice(t *testing.T) {
	// Setup
	menuRepo := &mockMenuRepo{
		items: []*menu.MenuItem{
			{
				ID:          primitive.NewObjectID(),
				Name:        "Free Item",
				Category:    "Promo",
				Price:       0,
				CurrentCost: 10000,
				CostStatus:  menu.CostStatusFinal,
			},
		},
	}

	service := NewProfitAnalyzerService(menuRepo, nil, nil, nil)

	// Test
	profit, err := service.CalculateMenuItemProfit(context.Background(), menuRepo.items[0].ID)

	// Verify
	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}

	if profit.ProfitMargin != 0 {
		t.Errorf("Expected profit margin 0, got %.2f", profit.ProfitMargin)
	}

	if profit.AbsoluteProfit != -10000 {
		t.Errorf("Expected absolute profit -10000, got %.2f", profit.AbsoluteProfit)
	}

	if profit.WarningStatus != WarningNone {
		t.Errorf("Expected warning status none for promotional item, got %s", profit.WarningStatus)
	}
}

// Test edge case: cost = price (break-even)
func TestCalculateMenuItemProfit_BreakEven(t *testing.T) {
	// Setup
	menuRepo := &mockMenuRepo{
		items: []*menu.MenuItem{
			{
				ID:          primitive.NewObjectID(),
				Name:        "Break Even Item",
				Category:    "Coffee",
				Price:       20000,
				CurrentCost: 20000,
				CostStatus:  menu.CostStatusFinal,
			},
		},
	}

	service := NewProfitAnalyzerService(menuRepo, nil, nil, nil)

	// Test
	profit, err := service.CalculateMenuItemProfit(context.Background(), menuRepo.items[0].ID)

	// Verify
	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}

	if profit.ProfitMargin != 0 {
		t.Errorf("Expected profit margin 0, got %.2f", profit.ProfitMargin)
	}

	if profit.AbsoluteProfit != 0 {
		t.Errorf("Expected absolute profit 0, got %.2f", profit.AbsoluteProfit)
	}
}

// Test GetAllMenuItemProfits with summary statistics
func TestGetAllMenuItemProfits_SummaryStatistics(t *testing.T) {
	// Setup - mix of profitable, low margin, and loss items
	menuRepo := &mockMenuRepo{
		items: []*menu.MenuItem{
			{
				ID:          primitive.NewObjectID(),
				Name:        "High Profit Item",
				Category:    "Coffee",
				Price:       45000,
				CurrentCost: 15000,
				CostStatus:  menu.CostStatusFinal,
			},
			{
				ID:          primitive.NewObjectID(),
				Name:        "Low Margin Item",
				Category:    "Coffee",
				Price:       30000,
				CurrentCost: 25000,
				CostStatus:  menu.CostStatusFinal,
			},
			{
				ID:          primitive.NewObjectID(),
				Name:        "Loss Item",
				Category:    "Coffee",
				Price:       20000,
				CurrentCost: 25000,
				CostStatus:  menu.CostStatusFinal,
			},
			{
				ID:          primitive.NewObjectID(),
				Name:        "Incomplete Item",
				Category:    "Coffee",
				Price:       30000,
				CurrentCost: 0,
				CostStatus:  menu.CostStatusIncomplete,
			},
		},
	}

	settingsRepo := &mockSettingsRepo{
		settings: &settings.ShopSettings{
			LowMarginThreshold: 20.0,
		},
	}

	service := NewProfitAnalyzerService(menuRepo, nil, settingsRepo, nil)

	// Test
	filter := ProfitFilter{}
	response, err := service.GetAllMenuItemProfits(context.Background(), filter)

	// Verify
	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}

	// Check total items
	if response.Summary.TotalItems != 4 {
		t.Errorf("Expected 4 total items, got %d", response.Summary.TotalItems)
	}

	// Check loss count
	if response.Summary.LossCount != 1 {
		t.Errorf("Expected 1 loss item, got %d", response.Summary.LossCount)
	}

	// Check low margin count
	if response.Summary.LowMarginCount != 1 {
		t.Errorf("Expected 1 low margin item, got %d", response.Summary.LowMarginCount)
	}

	// Check average profit margin
	// Average should be calculated from items with price > 0 and not incomplete
	// High Profit: 66.67%, Low Margin: 16.67%, Loss: -25%
	// Average = (66.67 + 16.67 + (-25)) / 3 = 19.45
	expectedAvg := 19.45
	if response.Summary.AverageProfitMargin < expectedAvg-0.1 || response.Summary.AverageProfitMargin > expectedAvg+0.1 {
		t.Errorf("Expected average profit margin around %.2f, got %.2f", expectedAvg, response.Summary.AverageProfitMargin)
	}
}

// Test GetAllMenuItemProfits with all sorting options
func TestGetAllMenuItemProfits_SortByAbsoluteProfit(t *testing.T) {
	// Setup
	menuRepo := &mockMenuRepo{
		items: []*menu.MenuItem{
			{
				ID:          primitive.NewObjectID(),
				Name:        "Item A",
				Category:    "Coffee",
				Price:       30000,
				CurrentCost: 25000,
				CostStatus:  menu.CostStatusFinal,
			},
			{
				ID:          primitive.NewObjectID(),
				Name:        "Item B",
				Category:    "Coffee",
				Price:       45000,
				CurrentCost: 15000,
				CostStatus:  menu.CostStatusFinal,
			},
		},
	}

	settingsRepo := &mockSettingsRepo{
		settings: &settings.ShopSettings{
			LowMarginThreshold: 20.0,
		},
	}

	service := NewProfitAnalyzerService(menuRepo, nil, settingsRepo, nil)

	// Test - sort by absolute_profit descending
	filter := ProfitFilter{
		SortBy:    "absolute_profit",
		SortOrder: "desc",
	}
	response, err := service.GetAllMenuItemProfits(context.Background(), filter)

	// Verify
	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}

	if len(response.Items) != 2 {
		t.Errorf("Expected 2 items, got %d", len(response.Items))
	}

	// First item should have higher absolute profit (30000 > 5000)
	if response.Items[0].AbsoluteProfit < response.Items[1].AbsoluteProfit {
		t.Errorf("Expected items sorted by absolute profit descending, got %.2f and %.2f", 
			response.Items[0].AbsoluteProfit, response.Items[1].AbsoluteProfit)
	}
}

// Test GetAllMenuItemProfits sort by name
func TestGetAllMenuItemProfits_SortByName(t *testing.T) {
	// Setup
	menuRepo := &mockMenuRepo{
		items: []*menu.MenuItem{
			{
				ID:          primitive.NewObjectID(),
				Name:        "Zebra Coffee",
				Category:    "Coffee",
				Price:       30000,
				CurrentCost: 10000,
				CostStatus:  menu.CostStatusFinal,
			},
			{
				ID:          primitive.NewObjectID(),
				Name:        "Apple Tea",
				Category:    "Tea",
				Price:       25000,
				CurrentCost: 8000,
				CostStatus:  menu.CostStatusFinal,
			},
		},
	}

	settingsRepo := &mockSettingsRepo{
		settings: &settings.ShopSettings{
			LowMarginThreshold: 20.0,
		},
	}

	service := NewProfitAnalyzerService(menuRepo, nil, settingsRepo, nil)

	// Test - sort by name ascending
	filter := ProfitFilter{
		SortBy:    "name",
		SortOrder: "asc",
	}
	response, err := service.GetAllMenuItemProfits(context.Background(), filter)

	// Verify
	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}

	if len(response.Items) != 2 {
		t.Errorf("Expected 2 items, got %d", len(response.Items))
	}

	// First item should be "Apple Tea" (alphabetically first)
	if response.Items[0].Name != "Apple Tea" {
		t.Errorf("Expected first item to be 'Apple Tea', got '%s'", response.Items[0].Name)
	}

	if response.Items[1].Name != "Zebra Coffee" {
		t.Errorf("Expected second item to be 'Zebra Coffee', got '%s'", response.Items[1].Name)
	}
}

// Test GetCategoryProfits - basic functionality
func TestGetCategoryProfits_Basic(t *testing.T) {
	// Setup
	menuItemID1 := primitive.NewObjectID()
	menuItemID2 := primitive.NewObjectID()
	orderID1 := primitive.NewObjectID()
	orderID2 := primitive.NewObjectID()

	menuRepo := &mockMenuRepo{
		items: []*menu.MenuItem{
			{
				ID:       menuItemID1,
				Name:     "Cappuccino",
				Category: "Coffee",
			},
			{
				ID:       menuItemID2,
				Name:     "Green Tea",
				Category: "Tea",
			},
		},
	}

	baseTime := time.Date(2024, 1, 15, 10, 0, 0, 0, time.UTC)

	orderItemRepo := &mockOrderItemRepo{
		items: []*order.OrderItemWithCost{
			{
				ID:               primitive.NewObjectID(),
				OrderID:          orderID1,
				MenuItemID:       menuItemID1,
				Name:             "Cappuccino",
				Price:            45000,
				Quantity:         2,
				Subtotal:         90000,
				AccountingCost:   30000, // 15000 per item * 2
				CostStatus:       order.CostStatusFinal,
				CostCalculatedAt: baseTime,
				CreatedAt:        baseTime,
			},
			{
				ID:               primitive.NewObjectID(),
				OrderID:          orderID2,
				MenuItemID:       menuItemID2,
				Name:             "Green Tea",
				Price:            30000,
				Quantity:         3,
				Subtotal:         90000,
				AccountingCost:   30000, // 10000 per item * 3
				CostStatus:       order.CostStatusFinal,
				CostCalculatedAt: baseTime,
				CreatedAt:        baseTime,
			},
		},
	}

	service := NewProfitAnalyzerService(menuRepo, orderItemRepo, nil, nil)

	// Test
	dateRange := DateRange{
		Start: baseTime.Add(-1 * time.Hour),
		End:   baseTime.Add(1 * time.Hour),
	}
	categories, err := service.GetCategoryProfits(context.Background(), dateRange)

	// Verify
	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}

	if len(categories) != 2 {
		t.Fatalf("Expected 2 categories, got %d", len(categories))
	}

	// Find Coffee category
	var coffeeCategory *CategoryProfit
	var teaCategory *CategoryProfit
	for i := range categories {
		if categories[i].Category == "Coffee" {
			coffeeCategory = &categories[i]
		} else if categories[i].Category == "Tea" {
			teaCategory = &categories[i]
		}
	}

	if coffeeCategory == nil {
		t.Fatal("Coffee category not found")
	}
	if teaCategory == nil {
		t.Fatal("Tea category not found")
	}

	// Verify Coffee category
	if coffeeCategory.TotalRevenue != 90000 {
		t.Errorf("Expected Coffee total revenue 90000, got %.2f", coffeeCategory.TotalRevenue)
	}
	if coffeeCategory.TotalCost != 30000 {
		t.Errorf("Expected Coffee total cost 30000, got %.2f", coffeeCategory.TotalCost)
	}
	if coffeeCategory.TotalProfit != 60000 {
		t.Errorf("Expected Coffee total profit 60000, got %.2f", coffeeCategory.TotalProfit)
	}
	expectedMargin := 66.67
	if coffeeCategory.AverageProfitMargin < expectedMargin-0.1 || coffeeCategory.AverageProfitMargin > expectedMargin+0.1 {
		t.Errorf("Expected Coffee average profit margin %.2f, got %.2f", expectedMargin, coffeeCategory.AverageProfitMargin)
	}
	if coffeeCategory.ItemCount != 2 {
		t.Errorf("Expected Coffee item count 2, got %d", coffeeCategory.ItemCount)
	}
	if coffeeCategory.OrderCount != 1 {
		t.Errorf("Expected Coffee order count 1, got %d", coffeeCategory.OrderCount)
	}

	// Verify Tea category
	if teaCategory.TotalRevenue != 90000 {
		t.Errorf("Expected Tea total revenue 90000, got %.2f", teaCategory.TotalRevenue)
	}
	if teaCategory.TotalCost != 30000 {
		t.Errorf("Expected Tea total cost 30000, got %.2f", teaCategory.TotalCost)
	}
	if teaCategory.TotalProfit != 60000 {
		t.Errorf("Expected Tea total profit 60000, got %.2f", teaCategory.TotalProfit)
	}
	if teaCategory.ItemCount != 3 {
		t.Errorf("Expected Tea item count 3, got %d", teaCategory.ItemCount)
	}
	if teaCategory.OrderCount != 1 {
		t.Errorf("Expected Tea order count 1, got %d", teaCategory.OrderCount)
	}
}

// Test GetCategoryProfits - empty date range
func TestGetCategoryProfits_EmptyDateRange(t *testing.T) {
	// Setup
	menuRepo := &mockMenuRepo{
		items: []*menu.MenuItem{
			{
				ID:       primitive.NewObjectID(),
				Name:     "Cappuccino",
				Category: "Coffee",
			},
		},
	}

	orderItemRepo := &mockOrderItemRepo{
		items: []*order.OrderItemWithCost{},
	}

	service := NewProfitAnalyzerService(menuRepo, orderItemRepo, nil, nil)

	// Test
	dateRange := DateRange{
		Start: time.Date(2024, 1, 1, 0, 0, 0, 0, time.UTC),
		End:   time.Date(2024, 1, 31, 23, 59, 59, 0, time.UTC),
	}
	categories, err := service.GetCategoryProfits(context.Background(), dateRange)

	// Verify
	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}

	if len(categories) != 0 {
		t.Errorf("Expected 0 categories for empty date range, got %d", len(categories))
	}
}

// Test GetCategoryProfits - skip incomplete cost status
func TestGetCategoryProfits_SkipIncomplete(t *testing.T) {
	// Setup
	menuItemID1 := primitive.NewObjectID()
	menuItemID2 := primitive.NewObjectID()
	orderID1 := primitive.NewObjectID()

	menuRepo := &mockMenuRepo{
		items: []*menu.MenuItem{
			{
				ID:       menuItemID1,
				Name:     "Cappuccino",
				Category: "Coffee",
			},
			{
				ID:       menuItemID2,
				Name:     "Incomplete Item",
				Category: "Coffee",
			},
		},
	}

	baseTime := time.Date(2024, 1, 15, 10, 0, 0, 0, time.UTC)

	orderItemRepo := &mockOrderItemRepo{
		items: []*order.OrderItemWithCost{
			{
				ID:               primitive.NewObjectID(),
				OrderID:          orderID1,
				MenuItemID:       menuItemID1,
				Name:             "Cappuccino",
				Price:            45000,
				Quantity:         1,
				Subtotal:         45000,
				AccountingCost:   15000,
				CostStatus:       order.CostStatusFinal,
				CostCalculatedAt: baseTime,
				CreatedAt:        baseTime,
			},
			{
				ID:               primitive.NewObjectID(),
				OrderID:          orderID1,
				MenuItemID:       menuItemID2,
				Name:             "Incomplete Item",
				Price:            30000,
				Quantity:         1,
				Subtotal:         30000,
				AccountingCost:   0,
				CostStatus:       order.CostStatusIncomplete,
				CostCalculatedAt: baseTime,
				CreatedAt:        baseTime,
			},
		},
	}

	service := NewProfitAnalyzerService(menuRepo, orderItemRepo, nil, nil)

	// Test
	dateRange := DateRange{
		Start: baseTime.Add(-1 * time.Hour),
		End:   baseTime.Add(1 * time.Hour),
	}
	categories, err := service.GetCategoryProfits(context.Background(), dateRange)

	// Verify
	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}

	if len(categories) != 1 {
		t.Fatalf("Expected 1 category (incomplete item should be skipped), got %d", len(categories))
	}

	// Verify only the complete item is included
	if categories[0].TotalRevenue != 45000 {
		t.Errorf("Expected total revenue 45000 (incomplete item excluded), got %.2f", categories[0].TotalRevenue)
	}
	if categories[0].ItemCount != 1 {
		t.Errorf("Expected item count 1 (incomplete item excluded), got %d", categories[0].ItemCount)
	}
}

// Test GetCategoryProfits - multiple orders in same category
func TestGetCategoryProfits_MultipleOrders(t *testing.T) {
	// Setup
	menuItemID := primitive.NewObjectID()
	orderID1 := primitive.NewObjectID()
	orderID2 := primitive.NewObjectID()
	orderID3 := primitive.NewObjectID()

	menuRepo := &mockMenuRepo{
		items: []*menu.MenuItem{
			{
				ID:       menuItemID,
				Name:     "Cappuccino",
				Category: "Coffee",
			},
		},
	}

	baseTime := time.Date(2024, 1, 15, 10, 0, 0, 0, time.UTC)

	orderItemRepo := &mockOrderItemRepo{
		items: []*order.OrderItemWithCost{
			{
				ID:               primitive.NewObjectID(),
				OrderID:          orderID1,
				MenuItemID:       menuItemID,
				Name:             "Cappuccino",
				Price:            45000,
				Quantity:         1,
				Subtotal:         45000,
				AccountingCost:   15000,
				CostStatus:       order.CostStatusFinal,
				CostCalculatedAt: baseTime,
				CreatedAt:        baseTime,
			},
			{
				ID:               primitive.NewObjectID(),
				OrderID:          orderID2,
				MenuItemID:       menuItemID,
				Name:             "Cappuccino",
				Price:            45000,
				Quantity:         2,
				Subtotal:         90000,
				AccountingCost:   30000,
				CostStatus:       order.CostStatusFinal,
				CostCalculatedAt: baseTime,
				CreatedAt:        baseTime,
			},
			{
				ID:               primitive.NewObjectID(),
				OrderID:          orderID3,
				MenuItemID:       menuItemID,
				Name:             "Cappuccino",
				Price:            45000,
				Quantity:         1,
				Subtotal:         45000,
				AccountingCost:   15000,
				CostStatus:       order.CostStatusFinal,
				CostCalculatedAt: baseTime,
				CreatedAt:        baseTime,
			},
		},
	}

	service := NewProfitAnalyzerService(menuRepo, orderItemRepo, nil, nil)

	// Test
	dateRange := DateRange{
		Start: baseTime.Add(-1 * time.Hour),
		End:   baseTime.Add(1 * time.Hour),
	}
	categories, err := service.GetCategoryProfits(context.Background(), dateRange)

	// Verify
	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}

	if len(categories) != 1 {
		t.Fatalf("Expected 1 category, got %d", len(categories))
	}

	// Verify aggregated values
	expectedRevenue := 45000.0 + 90000.0 + 45000.0 // 180000
	if categories[0].TotalRevenue != expectedRevenue {
		t.Errorf("Expected total revenue %.2f, got %.2f", expectedRevenue, categories[0].TotalRevenue)
	}

	expectedCost := 15000.0 + 30000.0 + 15000.0 // 60000
	if categories[0].TotalCost != expectedCost {
		t.Errorf("Expected total cost %.2f, got %.2f", expectedCost, categories[0].TotalCost)
	}

	expectedProfit := expectedRevenue - expectedCost // 120000
	if categories[0].TotalProfit != expectedProfit {
		t.Errorf("Expected total profit %.2f, got %.2f", expectedProfit, categories[0].TotalProfit)
	}

	expectedItemCount := 1 + 2 + 1 // 4 items total
	if categories[0].ItemCount != expectedItemCount {
		t.Errorf("Expected item count %d, got %d", expectedItemCount, categories[0].ItemCount)
	}

	// Verify unique order count
	if categories[0].OrderCount != 3 {
		t.Errorf("Expected order count 3 (unique orders), got %d", categories[0].OrderCount)
	}
}

// Test GetCategoryProfits - uncategorized items
func TestGetCategoryProfits_Uncategorized(t *testing.T) {
	// Setup
	menuItemID := primitive.NewObjectID()
	orderID := primitive.NewObjectID()

	menuRepo := &mockMenuRepo{
		items: []*menu.MenuItem{
			{
				ID:       menuItemID,
				Name:     "Mystery Item",
				Category: "", // Empty category
			},
		},
	}

	baseTime := time.Date(2024, 1, 15, 10, 0, 0, 0, time.UTC)

	orderItemRepo := &mockOrderItemRepo{
		items: []*order.OrderItemWithCost{
			{
				ID:               primitive.NewObjectID(),
				OrderID:          orderID,
				MenuItemID:       menuItemID,
				Name:             "Mystery Item",
				Price:            25000,
				Quantity:         1,
				Subtotal:         25000,
				AccountingCost:   10000,
				CostStatus:       order.CostStatusFinal,
				CostCalculatedAt: baseTime,
				CreatedAt:        baseTime,
			},
		},
	}

	service := NewProfitAnalyzerService(menuRepo, orderItemRepo, nil, nil)

	// Test
	dateRange := DateRange{
		Start: baseTime.Add(-1 * time.Hour),
		End:   baseTime.Add(1 * time.Hour),
	}
	categories, err := service.GetCategoryProfits(context.Background(), dateRange)

	// Verify
	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}

	if len(categories) != 1 {
		t.Fatalf("Expected 1 category, got %d", len(categories))
	}

	// Verify category is "Uncategorized"
	if categories[0].Category != "Uncategorized" {
		t.Errorf("Expected category 'Uncategorized', got '%s'", categories[0].Category)
	}
}

// Test GetCategoryProfits - menu item not found (orphaned order item)
func TestGetCategoryProfits_MenuItemNotFound(t *testing.T) {
	// Setup
	menuItemID := primitive.NewObjectID()
	orphanedMenuItemID := primitive.NewObjectID()
	orderID := primitive.NewObjectID()

	menuRepo := &mockMenuRepo{
		items: []*menu.MenuItem{
			{
				ID:       menuItemID,
				Name:     "Cappuccino",
				Category: "Coffee",
			},
			// orphanedMenuItemID is not in the menu
		},
	}

	baseTime := time.Date(2024, 1, 15, 10, 0, 0, 0, time.UTC)

	orderItemRepo := &mockOrderItemRepo{
		items: []*order.OrderItemWithCost{
			{
				ID:               primitive.NewObjectID(),
				OrderID:          orderID,
				MenuItemID:       menuItemID,
				Name:             "Cappuccino",
				Price:            45000,
				Quantity:         1,
				Subtotal:         45000,
				AccountingCost:   15000,
				CostStatus:       order.CostStatusFinal,
				CostCalculatedAt: baseTime,
				CreatedAt:        baseTime,
			},
			{
				ID:               primitive.NewObjectID(),
				OrderID:          orderID,
				MenuItemID:       orphanedMenuItemID,
				Name:             "Deleted Item",
				Price:            30000,
				Quantity:         1,
				Subtotal:         30000,
				AccountingCost:   10000,
				CostStatus:       order.CostStatusFinal,
				CostCalculatedAt: baseTime,
				CreatedAt:        baseTime,
			},
		},
	}

	service := NewProfitAnalyzerService(menuRepo, orderItemRepo, nil, nil)

	// Test
	dateRange := DateRange{
		Start: baseTime.Add(-1 * time.Hour),
		End:   baseTime.Add(1 * time.Hour),
	}
	categories, err := service.GetCategoryProfits(context.Background(), dateRange)

	// Verify
	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}

	// Only the valid menu item should be included
	if len(categories) != 1 {
		t.Fatalf("Expected 1 category (orphaned item skipped), got %d", len(categories))
	}

	if categories[0].TotalRevenue != 45000 {
		t.Errorf("Expected total revenue 45000 (orphaned item excluded), got %.2f", categories[0].TotalRevenue)
	}
}

// Test GetCategoryProfits - date range filtering
func TestGetCategoryProfits_DateRangeFiltering(t *testing.T) {
	// Setup
	menuItemID := primitive.NewObjectID()
	orderID1 := primitive.NewObjectID()
	orderID2 := primitive.NewObjectID()

	menuRepo := &mockMenuRepo{
		items: []*menu.MenuItem{
			{
				ID:       menuItemID,
				Name:     "Cappuccino",
				Category: "Coffee",
			},
		},
	}

	baseTime := time.Date(2024, 1, 15, 10, 0, 0, 0, time.UTC)
	outsideTime := time.Date(2024, 2, 1, 10, 0, 0, 0, time.UTC)

	orderItemRepo := &mockOrderItemRepo{
		items: []*order.OrderItemWithCost{
			{
				ID:               primitive.NewObjectID(),
				OrderID:          orderID1,
				MenuItemID:       menuItemID,
				Name:             "Cappuccino",
				Price:            45000,
				Quantity:         1,
				Subtotal:         45000,
				AccountingCost:   15000,
				CostStatus:       order.CostStatusFinal,
				CostCalculatedAt: baseTime,
				CreatedAt:        baseTime, // Inside date range
			},
			{
				ID:               primitive.NewObjectID(),
				OrderID:          orderID2,
				MenuItemID:       menuItemID,
				Name:             "Cappuccino",
				Price:            45000,
				Quantity:         1,
				Subtotal:         45000,
				AccountingCost:   15000,
				CostStatus:       order.CostStatusFinal,
				CostCalculatedAt: outsideTime,
				CreatedAt:        outsideTime, // Outside date range
			},
		},
	}

	service := NewProfitAnalyzerService(menuRepo, orderItemRepo, nil, nil)

	// Test - date range that includes only the first order
	dateRange := DateRange{
		Start: time.Date(2024, 1, 1, 0, 0, 0, 0, time.UTC),
		End:   time.Date(2024, 1, 31, 23, 59, 59, 0, time.UTC),
	}
	categories, err := service.GetCategoryProfits(context.Background(), dateRange)

	// Verify
	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}

	if len(categories) != 1 {
		t.Fatalf("Expected 1 category, got %d", len(categories))
	}

	// Only the first order should be included
	if categories[0].TotalRevenue != 45000 {
		t.Errorf("Expected total revenue 45000 (only first order), got %.2f", categories[0].TotalRevenue)
	}
	if categories[0].OrderCount != 1 {
		t.Errorf("Expected order count 1, got %d", categories[0].OrderCount)
	}
}

// Test GetOperatingProfit - basic functionality with exact date match
func TestGetOperatingProfit_Basic(t *testing.T) {
	// Setup
	menuItemID := primitive.NewObjectID()
	orderID := primitive.NewObjectID()

	menuRepo := &mockMenuRepo{
		items: []*menu.MenuItem{
			{
				ID:       menuItemID,
				Name:     "Cappuccino",
				Category: "Coffee",
			},
		},
	}

	baseTime := time.Date(2024, 1, 15, 10, 0, 0, 0, time.UTC)
	startDate := time.Date(2024, 1, 1, 0, 0, 0, 0, time.UTC)
	endDate := time.Date(2024, 1, 31, 23, 59, 59, 0, time.UTC)

	orderItemRepo := &mockOrderItemRepo{
		items: []*order.OrderItemWithCost{
			{
				ID:               primitive.NewObjectID(),
				OrderID:          orderID,
				MenuItemID:       menuItemID,
				Name:             "Cappuccino",
				Price:            45000,
				Quantity:         10,
				Subtotal:         450000,
				AccountingCost:   150000, // 15000 per item * 10
				CostStatus:       order.CostStatusFinal,
				CostCalculatedAt: baseTime,
				CreatedAt:        baseTime,
			},
		},
	}

	operatingExpenseRepo := &mockOperatingExpenseRepo{
		expenses: []*expense.OperatingExpense{
			{
				ID:             primitive.NewObjectID(),
				PeriodStart:    startDate,
				PeriodEnd:      endDate,
				StaffSalary:    2000000,
				Rent:           1000000,
				Utilities:      500000,
				MarketingCosts: 300000,
				OtherExpenses:  200000,
				TotalExpenses:  4000000,
			},
		},
	}

	service := NewProfitAnalyzerService(menuRepo, orderItemRepo, nil, operatingExpenseRepo)

	// Test
	dateRange := DateRange{
		Start: startDate,
		End:   endDate,
	}
	report, err := service.GetOperatingProfit(context.Background(), dateRange)

	// Verify
	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}

	// Verify gross profit calculations
	expectedRevenue := 450000.0
	if report.TotalRevenue != expectedRevenue {
		t.Errorf("Expected total revenue %.2f, got %.2f", expectedRevenue, report.TotalRevenue)
	}

	expectedCOGS := 150000.0
	if report.TotalCOGS != expectedCOGS {
		t.Errorf("Expected total COGS %.2f, got %.2f", expectedCOGS, report.TotalCOGS)
	}

	expectedGrossProfit := 300000.0
	if report.GrossProfit != expectedGrossProfit {
		t.Errorf("Expected gross profit %.2f, got %.2f", expectedGrossProfit, report.GrossProfit)
	}

	expectedGrossProfitMargin := 66.67
	if report.GrossProfitMargin < expectedGrossProfitMargin-0.1 || report.GrossProfitMargin > expectedGrossProfitMargin+0.1 {
		t.Errorf("Expected gross profit margin %.2f, got %.2f", expectedGrossProfitMargin, report.GrossProfitMargin)
	}

	// Verify operating expenses
	if report.StaffSalary != 2000000 {
		t.Errorf("Expected staff salary 2000000, got %.2f", report.StaffSalary)
	}
	if report.Rent != 1000000 {
		t.Errorf("Expected rent 1000000, got %.2f", report.Rent)
	}
	if report.Utilities != 500000 {
		t.Errorf("Expected utilities 500000, got %.2f", report.Utilities)
	}
	if report.MarketingCosts != 300000 {
		t.Errorf("Expected marketing costs 300000, got %.2f", report.MarketingCosts)
	}
	if report.OtherExpenses != 200000 {
		t.Errorf("Expected other expenses 200000, got %.2f", report.OtherExpenses)
	}

	expectedTotalExpenses := 4000000.0
	if report.TotalExpenses != expectedTotalExpenses {
		t.Errorf("Expected total expenses %.2f, got %.2f", expectedTotalExpenses, report.TotalExpenses)
	}

	// Verify operating profit
	expectedOperatingProfit := -3700000.0 // 300000 - 4000000
	if report.OperatingProfit != expectedOperatingProfit {
		t.Errorf("Expected operating profit %.2f, got %.2f", expectedOperatingProfit, report.OperatingProfit)
	}

	expectedOperatingProfitMargin := -822.22
	if report.OperatingProfitMargin < expectedOperatingProfitMargin-0.1 || report.OperatingProfitMargin > expectedOperatingProfitMargin+0.1 {
		t.Errorf("Expected operating profit margin %.2f, got %.2f", expectedOperatingProfitMargin, report.OperatingProfitMargin)
	}

	// Verify expense allocated flag (should be false for exact match)
	if report.ExpenseAllocated {
		t.Errorf("Expected expense allocated to be false for exact date match")
	}

	if report.AllocationNote != "" {
		t.Errorf("Expected no allocation note for exact date match, got '%s'", report.AllocationNote)
	}
}

// Test GetOperatingProfit - no expenses (should return gross profit only)
func TestGetOperatingProfit_NoExpenses(t *testing.T) {
	// Setup
	menuItemID := primitive.NewObjectID()
	orderID := primitive.NewObjectID()

	menuRepo := &mockMenuRepo{
		items: []*menu.MenuItem{
			{
				ID:       menuItemID,
				Name:     "Cappuccino",
				Category: "Coffee",
			},
		},
	}

	baseTime := time.Date(2024, 1, 15, 10, 0, 0, 0, time.UTC)
	startDate := time.Date(2024, 1, 1, 0, 0, 0, 0, time.UTC)
	endDate := time.Date(2024, 1, 31, 23, 59, 59, 0, time.UTC)

	orderItemRepo := &mockOrderItemRepo{
		items: []*order.OrderItemWithCost{
			{
				ID:               primitive.NewObjectID(),
				OrderID:          orderID,
				MenuItemID:       menuItemID,
				Name:             "Cappuccino",
				Price:            45000,
				Quantity:         5,
				Subtotal:         225000,
				AccountingCost:   75000,
				CostStatus:       order.CostStatusFinal,
				CostCalculatedAt: baseTime,
				CreatedAt:        baseTime,
			},
		},
	}

	operatingExpenseRepo := &mockOperatingExpenseRepo{
		expenses: []*expense.OperatingExpense{}, // No expenses
	}

	service := NewProfitAnalyzerService(menuRepo, orderItemRepo, nil, operatingExpenseRepo)

	// Test
	dateRange := DateRange{
		Start: startDate,
		End:   endDate,
	}
	report, err := service.GetOperatingProfit(context.Background(), dateRange)

	// Verify
	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}

	// Verify gross profit is calculated
	expectedRevenue := 225000.0
	if report.TotalRevenue != expectedRevenue {
		t.Errorf("Expected total revenue %.2f, got %.2f", expectedRevenue, report.TotalRevenue)
	}

	expectedGrossProfit := 150000.0
	if report.GrossProfit != expectedGrossProfit {
		t.Errorf("Expected gross profit %.2f, got %.2f", expectedGrossProfit, report.GrossProfit)
	}

	// Verify all expenses are zero
	if report.StaffSalary != 0 {
		t.Errorf("Expected staff salary 0, got %.2f", report.StaffSalary)
	}
	if report.TotalExpenses != 0 {
		t.Errorf("Expected total expenses 0, got %.2f", report.TotalExpenses)
	}

	// Verify operating profit equals gross profit (no expenses)
	if report.OperatingProfit != expectedGrossProfit {
		t.Errorf("Expected operating profit to equal gross profit %.2f, got %.2f", expectedGrossProfit, report.OperatingProfit)
	}

	// Verify allocation note
	expectedNote := "Chưa nhập chi phí vận hành"
	if report.AllocationNote != expectedNote {
		t.Errorf("Expected allocation note '%s', got '%s'", expectedNote, report.AllocationNote)
	}
}

// Test GetOperatingProfit - expense allocation (monthly to daily)
func TestGetOperatingProfit_ExpenseAllocation(t *testing.T) {
	// Setup
	menuItemID := primitive.NewObjectID()
	orderID := primitive.NewObjectID()

	menuRepo := &mockMenuRepo{
		items: []*menu.MenuItem{
			{
				ID:       menuItemID,
				Name:     "Cappuccino",
				Category: "Coffee",
			},
		},
	}

	// Date range: Jan 1-7 (7 days)
	startDate := time.Date(2024, 1, 1, 0, 0, 0, 0, time.UTC)
	endDate := time.Date(2024, 1, 7, 23, 59, 59, 0, time.UTC)
	baseTime := time.Date(2024, 1, 3, 10, 0, 0, 0, time.UTC)

	orderItemRepo := &mockOrderItemRepo{
		items: []*order.OrderItemWithCost{
			{
				ID:               primitive.NewObjectID(),
				OrderID:          orderID,
				MenuItemID:       menuItemID,
				Name:             "Cappuccino",
				Price:            45000,
				Quantity:         10,
				Subtotal:         450000,
				AccountingCost:   150000,
				CostStatus:       order.CostStatusFinal,
				CostCalculatedAt: baseTime,
				CreatedAt:        baseTime,
			},
		},
	}

	// Monthly expense: Jan 1-31 (31 days)
	monthStart := time.Date(2024, 1, 1, 0, 0, 0, 0, time.UTC)
	monthEnd := time.Date(2024, 1, 31, 23, 59, 59, 0, time.UTC)

	operatingExpenseRepo := &mockOperatingExpenseRepo{
		expenses: []*expense.OperatingExpense{
			{
				ID:             primitive.NewObjectID(),
				PeriodStart:    monthStart,
				PeriodEnd:      monthEnd,
				StaffSalary:    3100000, // 100000 per day
				Rent:           3100000,
				Utilities:      1550000,
				MarketingCosts: 930000,
				OtherExpenses:  620000,
				TotalExpenses:  9300000,
			},
		},
	}

	service := NewProfitAnalyzerService(menuRepo, orderItemRepo, nil, operatingExpenseRepo)

	// Test
	dateRange := DateRange{
		Start: startDate,
		End:   endDate,
	}
	report, err := service.GetOperatingProfit(context.Background(), dateRange)

	// Verify
	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}

	// Verify gross profit
	expectedGrossProfit := 300000.0
	if report.GrossProfit != expectedGrossProfit {
		t.Errorf("Expected gross profit %.2f, got %.2f", expectedGrossProfit, report.GrossProfit)
	}

	// Verify expense allocation
	// 7 days out of 31 days = 7/31 ratio
	// Note: The actual calculation uses hours, so 7 days = 168 hours, 31 days = 744 hours
	// Allocation ratio = 168/744 = 0.2258...
	expectedRatio := 7.0 / 31.0

	expectedStaffSalary := 3100000 * expectedRatio
	if report.StaffSalary < expectedStaffSalary-1000 || report.StaffSalary > expectedStaffSalary+1000 {
		t.Errorf("Expected staff salary around %.2f, got %.2f", expectedStaffSalary, report.StaffSalary)
	}

	expectedRent := 3100000 * expectedRatio
	if report.Rent < expectedRent-1000 || report.Rent > expectedRent+1000 {
		t.Errorf("Expected rent around %.2f, got %.2f", expectedRent, report.Rent)
	}

	// Verify expense allocated flag is true
	if !report.ExpenseAllocated {
		t.Errorf("Expected expense allocated to be true for date range mismatch")
	}

	expectedNote := "Chi phí được phân bổ từ tháng"
	if report.AllocationNote != expectedNote {
		t.Errorf("Expected allocation note '%s', got '%s'", expectedNote, report.AllocationNote)
	}
}

// Test GetOperatingProfit - skip incomplete cost status
func TestGetOperatingProfit_SkipIncomplete(t *testing.T) {
	// Setup
	menuItemID1 := primitive.NewObjectID()
	menuItemID2 := primitive.NewObjectID()
	orderID := primitive.NewObjectID()

	menuRepo := &mockMenuRepo{
		items: []*menu.MenuItem{
			{
				ID:       menuItemID1,
				Name:     "Cappuccino",
				Category: "Coffee",
			},
			{
				ID:       menuItemID2,
				Name:     "Incomplete Item",
				Category: "Coffee",
			},
		},
	}

	baseTime := time.Date(2024, 1, 15, 10, 0, 0, 0, time.UTC)
	startDate := time.Date(2024, 1, 1, 0, 0, 0, 0, time.UTC)
	endDate := time.Date(2024, 1, 31, 23, 59, 59, 0, time.UTC)

	orderItemRepo := &mockOrderItemRepo{
		items: []*order.OrderItemWithCost{
			{
				ID:               primitive.NewObjectID(),
				OrderID:          orderID,
				MenuItemID:       menuItemID1,
				Name:             "Cappuccino",
				Price:            45000,
				Quantity:         1,
				Subtotal:         45000,
				AccountingCost:   15000,
				CostStatus:       order.CostStatusFinal,
				CostCalculatedAt: baseTime,
				CreatedAt:        baseTime,
			},
			{
				ID:               primitive.NewObjectID(),
				OrderID:          orderID,
				MenuItemID:       menuItemID2,
				Name:             "Incomplete Item",
				Price:            30000,
				Quantity:         1,
				Subtotal:         30000,
				AccountingCost:   0,
				CostStatus:       order.CostStatusIncomplete,
				CostCalculatedAt: baseTime,
				CreatedAt:        baseTime,
			},
		},
	}

	operatingExpenseRepo := &mockOperatingExpenseRepo{
		expenses: []*expense.OperatingExpense{},
	}

	service := NewProfitAnalyzerService(menuRepo, orderItemRepo, nil, operatingExpenseRepo)

	// Test
	dateRange := DateRange{
		Start: startDate,
		End:   endDate,
	}
	report, err := service.GetOperatingProfit(context.Background(), dateRange)

	// Verify
	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}

	// Verify only the complete item is included
	expectedRevenue := 45000.0
	if report.TotalRevenue != expectedRevenue {
		t.Errorf("Expected total revenue %.2f (incomplete item excluded), got %.2f", expectedRevenue, report.TotalRevenue)
	}

	expectedCOGS := 15000.0
	if report.TotalCOGS != expectedCOGS {
		t.Errorf("Expected total COGS %.2f (incomplete item excluded), got %.2f", expectedCOGS, report.TotalCOGS)
	}

	expectedGrossProfit := 30000.0
	if report.GrossProfit != expectedGrossProfit {
		t.Errorf("Expected gross profit %.2f, got %.2f", expectedGrossProfit, report.GrossProfit)
	}
}

// Test GetOperatingProfit - multiple overlapping expense periods
func TestGetOperatingProfit_MultipleExpensePeriods(t *testing.T) {
	// Setup
	menuItemID := primitive.NewObjectID()
	orderID := primitive.NewObjectID()

	menuRepo := &mockMenuRepo{
		items: []*menu.MenuItem{
			{
				ID:       menuItemID,
				Name:     "Cappuccino",
				Category: "Coffee",
			},
		},
	}

	baseTime := time.Date(2024, 1, 15, 10, 0, 0, 0, time.UTC)
	startDate := time.Date(2024, 1, 1, 0, 0, 0, 0, time.UTC)
	endDate := time.Date(2024, 1, 31, 23, 59, 59, 0, time.UTC)

	orderItemRepo := &mockOrderItemRepo{
		items: []*order.OrderItemWithCost{
			{
				ID:               primitive.NewObjectID(),
				OrderID:          orderID,
				MenuItemID:       menuItemID,
				Name:             "Cappuccino",
				Price:            45000,
				Quantity:         10,
				Subtotal:         450000,
				AccountingCost:   150000,
				CostStatus:       order.CostStatusFinal,
				CostCalculatedAt: baseTime,
				CreatedAt:        baseTime,
			},
		},
	}

	// Two expense periods that both overlap with the date range
	operatingExpenseRepo := &mockOperatingExpenseRepo{
		expenses: []*expense.OperatingExpense{
			{
				ID:             primitive.NewObjectID(),
				PeriodStart:    time.Date(2024, 1, 1, 0, 0, 0, 0, time.UTC),
				PeriodEnd:      time.Date(2024, 1, 15, 23, 59, 59, 0, time.UTC),
				StaffSalary:    1000000,
				Rent:           500000,
				Utilities:      250000,
				MarketingCosts: 150000,
				OtherExpenses:  100000,
				TotalExpenses:  2000000,
			},
			{
				ID:             primitive.NewObjectID(),
				PeriodStart:    time.Date(2024, 1, 16, 0, 0, 0, 0, time.UTC),
				PeriodEnd:      time.Date(2024, 1, 31, 23, 59, 59, 0, time.UTC),
				StaffSalary:    1000000,
				Rent:           500000,
				Utilities:      250000,
				MarketingCosts: 150000,
				OtherExpenses:  100000,
				TotalExpenses:  2000000,
			},
		},
	}

	service := NewProfitAnalyzerService(menuRepo, orderItemRepo, nil, operatingExpenseRepo)

	// Test
	dateRange := DateRange{
		Start: startDate,
		End:   endDate,
	}
	report, err := service.GetOperatingProfit(context.Background(), dateRange)

	// Verify
	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}

	// Verify gross profit
	expectedGrossProfit := 300000.0
	if report.GrossProfit != expectedGrossProfit {
		t.Errorf("Expected gross profit %.2f, got %.2f", expectedGrossProfit, report.GrossProfit)
	}

	// Both expense periods should be fully included (they cover the entire date range)
	expectedTotalExpenses := 4000000.0
	if report.TotalExpenses != expectedTotalExpenses {
		t.Errorf("Expected total expenses %.2f, got %.2f", expectedTotalExpenses, report.TotalExpenses)
	}

	// Verify operating profit
	expectedOperatingProfit := -3700000.0
	if report.OperatingProfit != expectedOperatingProfit {
		t.Errorf("Expected operating profit %.2f, got %.2f", expectedOperatingProfit, report.OperatingProfit)
	}
}

// Test GetOperatingProfit - zero revenue edge case
func TestGetOperatingProfit_ZeroRevenue(t *testing.T) {
	// Setup
	menuRepo := &mockMenuRepo{
		items: []*menu.MenuItem{},
	}

	orderItemRepo := &mockOrderItemRepo{
		items: []*order.OrderItemWithCost{}, // No orders
	}

	startDate := time.Date(2024, 1, 1, 0, 0, 0, 0, time.UTC)
	endDate := time.Date(2024, 1, 31, 23, 59, 59, 0, time.UTC)

	operatingExpenseRepo := &mockOperatingExpenseRepo{
		expenses: []*expense.OperatingExpense{
			{
				ID:             primitive.NewObjectID(),
				PeriodStart:    startDate,
				PeriodEnd:      endDate,
				StaffSalary:    2000000,
				Rent:           1000000,
				Utilities:      500000,
				MarketingCosts: 300000,
				OtherExpenses:  200000,
				TotalExpenses:  4000000,
			},
		},
	}

	service := NewProfitAnalyzerService(menuRepo, orderItemRepo, nil, operatingExpenseRepo)

	// Test
	dateRange := DateRange{
		Start: startDate,
		End:   endDate,
	}
	report, err := service.GetOperatingProfit(context.Background(), dateRange)

	// Verify
	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}

	// Verify all revenue and profit values are zero
	if report.TotalRevenue != 0 {
		t.Errorf("Expected total revenue 0, got %.2f", report.TotalRevenue)
	}
	if report.TotalCOGS != 0 {
		t.Errorf("Expected total COGS 0, got %.2f", report.TotalCOGS)
	}
	if report.GrossProfit != 0 {
		t.Errorf("Expected gross profit 0, got %.2f", report.GrossProfit)
	}
	if report.GrossProfitMargin != 0 {
		t.Errorf("Expected gross profit margin 0, got %.2f", report.GrossProfitMargin)
	}

	// Verify expenses are still included
	expectedTotalExpenses := 4000000.0
	if report.TotalExpenses != expectedTotalExpenses {
		t.Errorf("Expected total expenses %.2f, got %.2f", expectedTotalExpenses, report.TotalExpenses)
	}

	// Verify operating profit is negative (all expenses, no revenue)
	expectedOperatingProfit := -4000000.0
	if report.OperatingProfit != expectedOperatingProfit {
		t.Errorf("Expected operating profit %.2f, got %.2f", expectedOperatingProfit, report.OperatingProfit)
	}

	// Operating profit margin should be 0 when revenue is 0
	if report.OperatingProfitMargin != 0 {
		t.Errorf("Expected operating profit margin 0 when revenue is 0, got %.2f", report.OperatingProfitMargin)
	}
}

// Test GetOperatingProfit - positive operating profit
func TestGetOperatingProfit_PositiveProfit(t *testing.T) {
	// Setup
	menuItemID := primitive.NewObjectID()
	orderID := primitive.NewObjectID()

	menuRepo := &mockMenuRepo{
		items: []*menu.MenuItem{
			{
				ID:       menuItemID,
				Name:     "Cappuccino",
				Category: "Coffee",
			},
		},
	}

	baseTime := time.Date(2024, 1, 15, 10, 0, 0, 0, time.UTC)
	startDate := time.Date(2024, 1, 1, 0, 0, 0, 0, time.UTC)
	endDate := time.Date(2024, 1, 31, 23, 59, 59, 0, time.UTC)

	orderItemRepo := &mockOrderItemRepo{
		items: []*order.OrderItemWithCost{
			{
				ID:               primitive.NewObjectID(),
				OrderID:          orderID,
				MenuItemID:       menuItemID,
				Name:             "Cappuccino",
				Price:            45000,
				Quantity:         200, // High volume
				Subtotal:         9000000,
				AccountingCost:   3000000, // 15000 per item * 200
				CostStatus:       order.CostStatusFinal,
				CostCalculatedAt: baseTime,
				CreatedAt:        baseTime,
			},
		},
	}

	operatingExpenseRepo := &mockOperatingExpenseRepo{
		expenses: []*expense.OperatingExpense{
			{
				ID:             primitive.NewObjectID(),
				PeriodStart:    startDate,
				PeriodEnd:      endDate,
				StaffSalary:    2000000,
				Rent:           1000000,
				Utilities:      500000,
				MarketingCosts: 300000,
				OtherExpenses:  200000,
				TotalExpenses:  4000000,
			},
		},
	}

	service := NewProfitAnalyzerService(menuRepo, orderItemRepo, nil, operatingExpenseRepo)

	// Test
	dateRange := DateRange{
		Start: startDate,
		End:   endDate,
	}
	report, err := service.GetOperatingProfit(context.Background(), dateRange)

	// Verify
	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}

	// Verify gross profit
	expectedRevenue := 9000000.0
	expectedCOGS := 3000000.0
	expectedGrossProfit := 6000000.0

	if report.TotalRevenue != expectedRevenue {
		t.Errorf("Expected total revenue %.2f, got %.2f", expectedRevenue, report.TotalRevenue)
	}
	if report.TotalCOGS != expectedCOGS {
		t.Errorf("Expected total COGS %.2f, got %.2f", expectedCOGS, report.TotalCOGS)
	}
	if report.GrossProfit != expectedGrossProfit {
		t.Errorf("Expected gross profit %.2f, got %.2f", expectedGrossProfit, report.GrossProfit)
	}

	// Verify operating profit is positive
	expectedOperatingProfit := 2000000.0 // 6000000 - 4000000
	if report.OperatingProfit != expectedOperatingProfit {
		t.Errorf("Expected operating profit %.2f, got %.2f", expectedOperatingProfit, report.OperatingProfit)
	}

	// Verify operating profit margin
	expectedOperatingProfitMargin := 22.22 // (2000000 / 9000000) * 100
	if report.OperatingProfitMargin < expectedOperatingProfitMargin-0.1 || report.OperatingProfitMargin > expectedOperatingProfitMargin+0.1 {
		t.Errorf("Expected operating profit margin %.2f, got %.2f", expectedOperatingProfitMargin, report.OperatingProfitMargin)
	}
}
