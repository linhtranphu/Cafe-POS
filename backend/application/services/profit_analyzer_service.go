package services

import (
	"context"
	"fmt"
	"math"
	"time"

	"cafe-pos/backend/domain/expense"
	"cafe-pos/backend/domain/menu"
	"cafe-pos/backend/domain/order"
	"cafe-pos/backend/domain/settings"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

// WarningStatus represents the warning status for a menu item
type WarningStatus string

const (
	WarningNone      WarningStatus = "none"
	WarningLowMargin WarningStatus = "low_margin"
	WarningLoss      WarningStatus = "loss"
)

// VariantProfit represents profit metrics for a single variant
type VariantProfit struct {
	VariantID      string          `json:"variant_id"`
	VariantName    string          `json:"variant_name"`
	Price          float64         `json:"price"`
	CurrentCost    float64         `json:"current_cost"`
	ProfitMargin   float64         `json:"profit_margin"`
	AbsoluteProfit float64         `json:"absolute_profit"`
	CostStatus     menu.CostStatus `json:"cost_status"`
	IsDefault      bool            `json:"is_default"`
	WarningStatus  WarningStatus   `json:"warning_status"`
}

// MenuItemProfit represents profit metrics for a menu item
type MenuItemProfit struct {
	MenuItemID           primitive.ObjectID `json:"menu_item_id"`
	Name                 string             `json:"name"`
	Category             string             `json:"category"`
	Price                float64            `json:"price"`
	CurrentCost          float64            `json:"current_cost"`
	ProfitMargin         float64            `json:"profit_margin"`
	AbsoluteProfit       float64            `json:"absolute_profit"`
	CostStatus           menu.CostStatus    `json:"cost_status"`
	CostLastCalculatedAt time.Time          `json:"cost_last_calculated_at"`
	WarningStatus        WarningStatus      `json:"warning_status"`
	HasVariants          bool               `json:"has_variants"`
	Variants             []VariantProfit    `json:"variants,omitempty"`
}

// CategoryProfit represents profit analysis for a menu category
type CategoryProfit struct {
	Category            string  `json:"category"`
	TotalRevenue        float64 `json:"total_revenue"`
	TotalCost           float64 `json:"total_cost"`
	TotalProfit         float64 `json:"total_profit"`
	AverageProfitMargin float64 `json:"average_profit_margin"`
	OrderCount          int     `json:"order_count"`
	ItemCount           int     `json:"item_count"`
}

// OperatingProfitReport represents operating profit analysis
type OperatingProfitReport struct {
	DateRange             DateRange `json:"date_range"`
	TotalRevenue          float64   `json:"total_revenue"`
	TotalCOGS             float64   `json:"total_cogs"`
	GrossProfit           float64   `json:"gross_profit"`
	GrossProfitMargin     float64   `json:"gross_profit_margin"`
	StaffSalary           float64   `json:"staff_salary"`
	Rent                  float64   `json:"rent"`
	Utilities             float64   `json:"utilities"`
	MarketingCosts        float64   `json:"marketing_costs"`
	OtherExpenses         float64   `json:"other_expenses"`
	TotalExpenses         float64   `json:"total_expenses"`
	OperatingProfit       float64   `json:"operating_profit"`
	OperatingProfitMargin float64   `json:"operating_profit_margin"`
	ExpenseAllocated      bool      `json:"expense_allocated"`
	AllocationNote        string    `json:"allocation_note,omitempty"`
}

// DateRange represents a date range filter
type DateRange struct {
	Start time.Time `json:"start"`
	End   time.Time `json:"end"`
}

// ProfitFilter represents filter options for profit analysis
type ProfitFilter struct {
	Category  string
	SortBy    string // "profit_margin", "absolute_profit", "name"
	SortOrder string // "asc", "desc"
}

// ProfitWarnings represents loss and low margin warnings
type ProfitWarnings struct {
	LossItems      []MenuItemProfit `json:"loss_items"`
	LowMarginItems []MenuItemProfit `json:"low_margin_items"`
	LossCount      int              `json:"loss_count"`
	LowMarginCount int              `json:"low_margin_count"`
	Threshold      float64          `json:"threshold"`
}

// SettingsRepository interface for shop settings
type SettingsRepository interface {
	FindFirst(ctx context.Context) (*settings.ShopSettings, error)
	Update(ctx context.Context, id primitive.ObjectID, settings *settings.ShopSettings) error
}

// OperatingExpenseRepository interface for operating expenses
// This interface is shared by both ProfitAnalyzerService and OperatingExpenseService
type OperatingExpenseRepository interface {
	Create(ctx context.Context, operatingExpense *expense.OperatingExpense) error
	Update(ctx context.Context, id primitive.ObjectID, operatingExpense *expense.OperatingExpense) error
	FindByID(ctx context.Context, id primitive.ObjectID) (*expense.OperatingExpense, error)
	FindByPeriod(ctx context.Context, startDate, endDate time.Time) ([]*expense.OperatingExpense, error)
	FindForDate(ctx context.Context, date time.Time) (*expense.OperatingExpense, error)
	FindAll(ctx context.Context, startDate, endDate *time.Time) ([]*expense.OperatingExpense, error)
	Upsert(ctx context.Context, operatingExpense *expense.OperatingExpense) (*expense.OperatingExpense, error)
	// Legacy methods for backward compatibility
	FindByDateRange(ctx context.Context, start, end time.Time) ([]*expense.OperatingExpense, error)
	FindByDate(ctx context.Context, date time.Time) (*expense.OperatingExpense, error)
}

// ProfitAnalyzerService handles profit analysis for menu items
type ProfitAnalyzerService struct {
	menuRepo            MenuRepository
	orderItemRepo       OrderItemRepository
	settingsRepo        SettingsRepository
	operatingExpenseRepo OperatingExpenseRepository
}

// NewProfitAnalyzerService creates a new profit analyzer service
func NewProfitAnalyzerService(
	menuRepo MenuRepository,
	orderItemRepo OrderItemRepository,
	settingsRepo SettingsRepository,
	operatingExpenseRepo OperatingExpenseRepository,
) *ProfitAnalyzerService {
	return &ProfitAnalyzerService{
		menuRepo:            menuRepo,
		orderItemRepo:       orderItemRepo,
		settingsRepo:        settingsRepo,
		operatingExpenseRepo: operatingExpenseRepo,
	}
}

// CalculateMenuItemProfit calculates profit metrics for a menu item
// Requirements: 2.1, 2.5, 2.6
func (s *ProfitAnalyzerService) CalculateMenuItemProfit(ctx context.Context, menuItemID primitive.ObjectID) (*MenuItemProfit, error) {
	// Fetch the menu item
	menuItem, err := s.menuRepo.FindByID(ctx, menuItemID)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch menu item: %w", err)
	}

	// Calculate profit metrics
	profit := s.calculateProfitMetrics(menuItem, 0) // threshold will be fetched separately if needed

	return profit, nil
}

// calculateProfitMetrics is a helper method that calculates profit metrics for a menu item
func (s *ProfitAnalyzerService) calculateProfitMetrics(menuItem *menu.MenuItem, lowMarginThreshold float64) *MenuItemProfit {
	profit := &MenuItemProfit{
		MenuItemID:           menuItem.ID,
		Name:                 menuItem.Name,
		Category:             menuItem.Category,
		CostLastCalculatedAt: menuItem.CostLastCalculatedAt,
		HasVariants:          menuItem.HasVariants,
	}

	if menuItem.HasVariants && len(menuItem.Variants) > 0 {
		// Multi-variant: calculate per-variant profit metrics
		variants := make([]VariantProfit, 0, len(menuItem.Variants))
		for _, v := range menuItem.Variants {
			vp := VariantProfit{
				VariantID:   v.ID,
				VariantName: v.Name,
				Price:       v.Price,
				CurrentCost: v.CurrentCost,
				CostStatus:  v.CostStatus,
				IsDefault:   v.IsDefault,
			}
			if v.Price > 0 {
				margin := ((v.Price - v.CurrentCost) / v.Price) * 100
				vp.ProfitMargin = math.Round(margin*100) / 100
				vp.AbsoluteProfit = v.Price - v.CurrentCost
			} else {
				vp.AbsoluteProfit = -v.CurrentCost
			}
			vp.WarningStatus = s.determineWarningStatus(v.Price, v.CurrentCost, vp.ProfitMargin, lowMarginThreshold)
			variants = append(variants, vp)
		}
		profit.Variants = variants

		// Use default variant for top-level sort/display fields
		defaultVariant := menuItem.GetDefaultVariant()
		if defaultVariant == nil {
			defaultVariant = &menuItem.Variants[0]
		}
		profit.Price = defaultVariant.Price
		profit.CurrentCost = defaultVariant.CurrentCost
		profit.CostStatus = defaultVariant.CostStatus
		if defaultVariant.Price > 0 {
			margin := ((defaultVariant.Price - defaultVariant.CurrentCost) / defaultVariant.Price) * 100
			profit.ProfitMargin = math.Round(margin*100) / 100
			profit.AbsoluteProfit = defaultVariant.Price - defaultVariant.CurrentCost
		} else {
			profit.AbsoluteProfit = -defaultVariant.CurrentCost
		}
		profit.WarningStatus = s.determineWarningStatus(profit.Price, profit.CurrentCost, profit.ProfitMargin, lowMarginThreshold)
		return profit
	}

	// Single-size item
	profit.Price = menuItem.Price
	profit.CurrentCost = menuItem.CurrentCost
	profit.CostStatus = menuItem.CostStatus

	if menuItem.Price <= 0 {
		profit.ProfitMargin = 0
		profit.AbsoluteProfit = -menuItem.CurrentCost
		profit.WarningStatus = WarningNone
		return profit
	}

	profitMargin := ((menuItem.Price - menuItem.CurrentCost) / menuItem.Price) * 100
	profit.ProfitMargin = math.Round(profitMargin*100) / 100
	profit.AbsoluteProfit = menuItem.Price - menuItem.CurrentCost
	profit.WarningStatus = s.determineWarningStatus(menuItem.Price, menuItem.CurrentCost, profit.ProfitMargin, lowMarginThreshold)

	return profit
}

// determineWarningStatus determines the warning status based on cost, price, and margin
func (s *ProfitAnalyzerService) determineWarningStatus(price, cost, profitMargin, threshold float64) WarningStatus {
	// Check if cost > price → mark as "loss"
	if cost > price {
		return WarningLoss
	}

	// Check if profit_margin < threshold → mark as "low_margin"
	if threshold > 0 && profitMargin < threshold {
		return WarningLowMargin
	}

	// Otherwise mark as "none"
	return WarningNone
}

// DetectWarningStatus detects loss and low margin items
// Requirements: 3.1, 3.2, 3.6
func (s *ProfitAnalyzerService) DetectWarningStatus(ctx context.Context, lowMarginThreshold float64) (*ProfitWarnings, error) {
	// If threshold not provided, fetch from settings
	if lowMarginThreshold <= 0 {
		settings, err := s.settingsRepo.FindFirst(ctx)
		if err != nil {
			return nil, fmt.Errorf("failed to fetch shop settings: %w", err)
		}
		lowMarginThreshold = settings.LowMarginThreshold
	}

	// Fetch all menu items
	menuItems, err := s.menuRepo.FindAll(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch menu items: %w", err)
	}

	warnings := &ProfitWarnings{
		LossItems:      []MenuItemProfit{},
		LowMarginItems: []MenuItemProfit{},
		Threshold:      lowMarginThreshold,
	}

	// Check each menu item for warnings
	for _, menuItem := range menuItems {
		// Skip items with incomplete cost status
		if menuItem.CostStatus == menu.CostStatusIncomplete {
			continue
		}

		// Calculate profit metrics
		profit := s.calculateProfitMetrics(menuItem, lowMarginThreshold)

		// Categorize by warning status
		if profit.WarningStatus == WarningLoss {
			warnings.LossItems = append(warnings.LossItems, *profit)
			warnings.LossCount++
		} else if profit.WarningStatus == WarningLowMargin {
			warnings.LowMarginItems = append(warnings.LowMarginItems, *profit)
			warnings.LowMarginCount++
		}
	}

	return warnings, nil
}

// ProfitSummary represents summary statistics for profit analysis
type ProfitSummary struct {
	TotalItems          int     `json:"total_items"`
	LossCount           int     `json:"loss_count"`
	LowMarginCount      int     `json:"low_margin_count"`
	AverageProfitMargin float64 `json:"average_profit_margin"`
}

// GetAllMenuItemProfitsResponse represents the response for GetAllMenuItemProfits
type GetAllMenuItemProfitsResponse struct {
	Items   []MenuItemProfit `json:"items"`
	Summary ProfitSummary    `json:"summary"`
}

// GetAllMenuItemProfits fetches all menu items with profit metrics
// Applies filters (category, sort) and returns with summary statistics
// Requirements: 2.1, 2.9, 4.3, 4.4
func (s *ProfitAnalyzerService) GetAllMenuItemProfits(ctx context.Context, filter ProfitFilter) (*GetAllMenuItemProfitsResponse, error) {
	// Fetch shop settings for low margin threshold
	settings, err := s.settingsRepo.FindFirst(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch shop settings: %w", err)
	}
	lowMarginThreshold := settings.LowMarginThreshold

	// Fetch all menu items
	var menuItems []*menu.MenuItem
	if filter.Category != "" {
		// Filter by category
		menuItems, err = s.menuRepo.FindByCategory(ctx, filter.Category)
		if err != nil {
			return nil, fmt.Errorf("failed to fetch menu items by category: %w", err)
		}
	} else {
		// Fetch all menu items
		menuItems, err = s.menuRepo.FindAll(ctx)
		if err != nil {
			return nil, fmt.Errorf("failed to fetch menu items: %w", err)
		}
	}

	// Calculate profit metrics for each item
	profits := make([]MenuItemProfit, 0, len(menuItems))
	var totalProfitMargin float64
	var profitMarginCount int
	lossCount := 0
	lowMarginCount := 0

	for _, menuItem := range menuItems {
		profit := s.calculateProfitMetrics(menuItem, lowMarginThreshold)
		profits = append(profits, *profit)

		// Update summary statistics
		if profit.WarningStatus == WarningLoss {
			lossCount++
		} else if profit.WarningStatus == WarningLowMargin {
			lowMarginCount++
		}

		// Calculate average profit margin (excluding items with price <= 0)
		if profit.Price > 0 && profit.CostStatus != menu.CostStatusIncomplete {
			totalProfitMargin += profit.ProfitMargin
			profitMarginCount++
		}
	}

	// Sort items based on filter
	s.sortProfits(profits, filter.SortBy, filter.SortOrder)

	// Calculate average profit margin
	averageProfitMargin := 0.0
	if profitMarginCount > 0 {
		averageProfitMargin = math.Round((totalProfitMargin/float64(profitMarginCount))*100) / 100
	}

	response := &GetAllMenuItemProfitsResponse{
		Items: profits,
		Summary: ProfitSummary{
			TotalItems:          len(profits),
			LossCount:           lossCount,
			LowMarginCount:      lowMarginCount,
			AverageProfitMargin: averageProfitMargin,
		},
	}

	return response, nil
}

// sortProfits sorts the profit items based on the specified field and order
func (s *ProfitAnalyzerService) sortProfits(profits []MenuItemProfit, sortBy, sortOrder string) {
	// Default sort by profit_margin descending
	if sortBy == "" {
		sortBy = "profit_margin"
	}
	if sortOrder == "" {
		sortOrder = "desc"
	}

	// Simple bubble sort for simplicity (can be optimized with sort.Slice)
	n := len(profits)
	for i := 0; i < n-1; i++ {
		for j := 0; j < n-i-1; j++ {
			shouldSwap := false

			switch sortBy {
			case "profit_margin":
				if sortOrder == "asc" {
					shouldSwap = profits[j].ProfitMargin > profits[j+1].ProfitMargin
				} else {
					shouldSwap = profits[j].ProfitMargin < profits[j+1].ProfitMargin
				}
			case "absolute_profit":
				if sortOrder == "asc" {
					shouldSwap = profits[j].AbsoluteProfit > profits[j+1].AbsoluteProfit
				} else {
					shouldSwap = profits[j].AbsoluteProfit < profits[j+1].AbsoluteProfit
				}
			case "name":
				if sortOrder == "asc" {
					shouldSwap = profits[j].Name > profits[j+1].Name
				} else {
					shouldSwap = profits[j].Name < profits[j+1].Name
				}
			}

			if shouldSwap {
				profits[j], profits[j+1] = profits[j+1], profits[j]
			}
		}
	}
}

// GetCategoryProfits calculates profit analysis aggregated by menu category
// Uses accounting_cost from order items (not current_cost)
// Requirements: 6.1, 6.2, 6.3, 6.4
func (s *ProfitAnalyzerService) GetCategoryProfits(ctx context.Context, dateRange DateRange) ([]CategoryProfit, error) {
	// Fetch all order items within the date range
	orderItems, err := s.orderItemRepo.FindByDateRange(ctx, dateRange.Start, dateRange.End)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch order items: %w", err)
	}

	// If no order items, return empty result
	if len(orderItems) == 0 {
		return []CategoryProfit{}, nil
	}

	// Fetch all menu items to get category information
	menuItems, err := s.menuRepo.FindAll(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch menu items: %w", err)
	}

	// Build menu item lookup map by ID
	menuItemMap := make(map[primitive.ObjectID]*menu.MenuItem)
	for _, item := range menuItems {
		menuItemMap[item.ID] = item
	}

	// Aggregate by category
	categoryMap := make(map[string]*CategoryProfit)

	for _, orderItem := range orderItems {
		// Skip items with incomplete cost status
		if orderItem.CostStatus == order.CostStatusIncomplete {
			continue
		}

		// Find the menu item to get category
		menuItem, exists := menuItemMap[orderItem.MenuItemID]
		if !exists {
			// Menu item not found, skip
			continue
		}

		category := menuItem.Category
		if category == "" {
			category = "Uncategorized"
		}

		// Initialize category if not exists
		if _, exists := categoryMap[category]; !exists {
			categoryMap[category] = &CategoryProfit{
				Category: category,
			}
		}

		// Calculate revenue and cost for this order item
		revenue := orderItem.Price * float64(orderItem.Quantity)
		cost := orderItem.AccountingCost // Already includes quantity

		// Update category totals
		categoryMap[category].TotalRevenue += revenue
		categoryMap[category].TotalCost += cost
		categoryMap[category].ItemCount += orderItem.Quantity
		
		// Count unique orders (we'll need to track order IDs)
		// For now, we'll increment order count per item (not accurate but simple)
		// TODO: Track unique order IDs for accurate order count
	}

	// Calculate profit and margin for each category
	categories := make([]CategoryProfit, 0, len(categoryMap))
	for _, cat := range categoryMap {
		cat.TotalProfit = cat.TotalRevenue - cat.TotalCost
		
		// Calculate average profit margin
		if cat.TotalRevenue > 0 {
			cat.AverageProfitMargin = (cat.TotalProfit / cat.TotalRevenue) * 100
			cat.AverageProfitMargin = math.Round(cat.AverageProfitMargin*100) / 100
		}

		// Round totals to 2 decimal places
		cat.TotalRevenue = math.Round(cat.TotalRevenue*100) / 100
		cat.TotalCost = math.Round(cat.TotalCost*100) / 100
		cat.TotalProfit = math.Round(cat.TotalProfit*100) / 100

		categories = append(categories, *cat)
	}

	// Count unique orders per category
	// We need to track order IDs to get accurate order count
	orderIDsByCategory := make(map[string]map[primitive.ObjectID]bool)
	for _, orderItem := range orderItems {
		if orderItem.CostStatus == order.CostStatusIncomplete {
			continue
		}
		menuItem, exists := menuItemMap[orderItem.MenuItemID]
		if !exists {
			continue
		}
		category := menuItem.Category
		if category == "" {
			category = "Uncategorized"
		}
		if orderIDsByCategory[category] == nil {
			orderIDsByCategory[category] = make(map[primitive.ObjectID]bool)
		}
		orderIDsByCategory[category][orderItem.OrderID] = true
	}

	// Update order counts
	for i := range categories {
		if orderIDs, exists := orderIDsByCategory[categories[i].Category]; exists {
			categories[i].OrderCount = len(orderIDs)
		}
	}

	return categories, nil
}

// GetOperatingProfit calculates operating profit after deducting operating expenses
// Requirements: 6.5.1, 6.5.3, 6.5.4, 6.5.6
func (s *ProfitAnalyzerService) GetOperatingProfit(ctx context.Context, dateRange DateRange) (*OperatingProfitReport, error) {
	// Initialize report
	report := &OperatingProfitReport{
		DateRange: dateRange,
	}

	// Fetch all order items within the date range
	orderItems, err := s.orderItemRepo.FindByDateRange(ctx, dateRange.Start, dateRange.End)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch order items: %w", err)
	}

	// Calculate gross profit from orders
	var totalRevenue float64
	var totalCOGS float64

	for _, orderItem := range orderItems {
		// Skip items with incomplete cost status
		if orderItem.CostStatus == order.CostStatusIncomplete {
			continue
		}

		// Calculate revenue and cost for this order item
		revenue := orderItem.Price * float64(orderItem.Quantity)
		cost := orderItem.AccountingCost // Already includes quantity

		totalRevenue += revenue
		totalCOGS += cost
	}

	// Round to 2 decimal places
	report.TotalRevenue = math.Round(totalRevenue*100) / 100
	report.TotalCOGS = math.Round(totalCOGS*100) / 100

	// Calculate gross profit
	report.GrossProfit = report.TotalRevenue - report.TotalCOGS
	report.GrossProfit = math.Round(report.GrossProfit*100) / 100

	// Calculate gross profit margin
	if report.TotalRevenue > 0 {
		report.GrossProfitMargin = (report.GrossProfit / report.TotalRevenue) * 100
		report.GrossProfitMargin = math.Round(report.GrossProfitMargin*100) / 100
	}

	// Fetch operating expenses for the date range
	expenses, err := s.operatingExpenseRepo.FindByPeriod(ctx, dateRange.Start, dateRange.End)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch operating expenses: %w", err)
	}

	// If no expenses found, return gross profit only with note
	if len(expenses) == 0 {
		// Operating profit equals gross profit when there are no expenses
		report.OperatingProfit = report.GrossProfit
		report.OperatingProfitMargin = report.GrossProfitMargin
		report.AllocationNote = "Chưa nhập chi phí vận hành"
		return report, nil
	}

	// Aggregate expenses
	// If multiple expense periods overlap with the date range, we need to allocate them
	var totalStaffSalary float64
	var totalRent float64
	var totalUtilities float64
	var totalMarketingCosts float64
	var totalOtherExpenses float64
	expenseAllocated := false

	for _, exp := range expenses {
		// Check if expense period exactly matches the date range
		if exp.PeriodStart.Equal(dateRange.Start) && exp.PeriodEnd.Equal(dateRange.End) {
			// Exact match, use full expense
			totalStaffSalary += exp.StaffSalary
			totalRent += exp.Rent
			totalUtilities += exp.Utilities
			totalMarketingCosts += exp.MarketingCosts
			totalOtherExpenses += exp.OtherExpenses
		} else {
			// Need to allocate expenses proportionally
			expenseAllocated = true
			
			// Calculate the number of days in the expense period
			expenseDays := exp.PeriodEnd.Sub(exp.PeriodStart).Hours() / 24
			if expenseDays <= 0 {
				expenseDays = 1
			}

			// Calculate the number of days in the date range that overlap with expense period
			overlapStart := dateRange.Start
			if exp.PeriodStart.After(overlapStart) {
				overlapStart = exp.PeriodStart
			}
			overlapEnd := dateRange.End
			if exp.PeriodEnd.Before(overlapEnd) {
				overlapEnd = exp.PeriodEnd
			}
			overlapDays := overlapEnd.Sub(overlapStart).Hours() / 24
			if overlapDays <= 0 {
				continue
			}

			// Calculate allocation ratio
			allocationRatio := overlapDays / expenseDays

			// Allocate expenses
			totalStaffSalary += exp.StaffSalary * allocationRatio
			totalRent += exp.Rent * allocationRatio
			totalUtilities += exp.Utilities * allocationRatio
			totalMarketingCosts += exp.MarketingCosts * allocationRatio
			totalOtherExpenses += exp.OtherExpenses * allocationRatio
		}
	}

	// Round expenses to 2 decimal places
	report.StaffSalary = math.Round(totalStaffSalary*100) / 100
	report.Rent = math.Round(totalRent*100) / 100
	report.Utilities = math.Round(totalUtilities*100) / 100
	report.MarketingCosts = math.Round(totalMarketingCosts*100) / 100
	report.OtherExpenses = math.Round(totalOtherExpenses*100) / 100

	// Calculate total expenses
	report.TotalExpenses = report.StaffSalary + report.Rent + report.Utilities + report.MarketingCosts + report.OtherExpenses
	report.TotalExpenses = math.Round(report.TotalExpenses*100) / 100

	// Calculate operating profit
	report.OperatingProfit = report.GrossProfit - report.TotalExpenses
	report.OperatingProfit = math.Round(report.OperatingProfit*100) / 100

	// Calculate operating profit margin
	if report.TotalRevenue > 0 {
		report.OperatingProfitMargin = (report.OperatingProfit / report.TotalRevenue) * 100
		report.OperatingProfitMargin = math.Round(report.OperatingProfitMargin*100) / 100
	}

	// Set expense allocated flag and note
	report.ExpenseAllocated = expenseAllocated
	if expenseAllocated {
		report.AllocationNote = "Chi phí được phân bổ từ tháng"
	}

	return report, nil
}
