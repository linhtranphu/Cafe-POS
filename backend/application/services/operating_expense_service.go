package services

import (
	"context"
	"fmt"
	"math"
	"time"

	"cafe-pos/backend/domain/expense"
)

// Note: OperatingExpenseRepository interface is defined in profit_analyzer_service.go
// to avoid duplication. It includes all methods needed by both services.

// OperatingExpenseService handles operating expense management and allocation
type OperatingExpenseService struct {
	expenseRepo OperatingExpenseRepository
}

// NewOperatingExpenseService creates a new operating expense service
func NewOperatingExpenseService(expenseRepo OperatingExpenseRepository) *OperatingExpenseService {
	return &OperatingExpenseService{
		expenseRepo: expenseRepo,
	}
}

// UpsertOperatingExpense creates or updates an operating expense for a period
// Validates that period_start <= period_end and all amounts >= 0
// Requirements: 6.5.2
func (s *OperatingExpenseService) UpsertOperatingExpense(ctx context.Context, req *expense.OperatingExpenseRequest) (*expense.OperatingExpense, error) {
	// Parse dates
	periodStart, err := time.Parse("2006-01-02", req.PeriodStart)
	if err != nil {
		return nil, fmt.Errorf("invalid period_start date format: %w", err)
	}
	
	periodEnd, err := time.Parse("2006-01-02", req.PeriodEnd)
	if err != nil {
		return nil, fmt.Errorf("invalid period_end date format: %w", err)
	}
	
	// Validate period_start <= period_end
	if periodStart.After(periodEnd) {
		return nil, fmt.Errorf("period_start must be before or equal to period_end")
	}
	
	// Validate all amounts >= 0 (already validated by binding tags, but double-check)
	if req.StaffSalary < 0 || req.Rent < 0 || req.Utilities < 0 || req.MarketingCosts < 0 || req.OtherExpenses < 0 {
		return nil, fmt.Errorf("all expense amounts must be >= 0")
	}
	
	// Create operating expense entity
	operatingExpense := &expense.OperatingExpense{
		PeriodStart:    periodStart,
		PeriodEnd:      periodEnd,
		StaffSalary:    req.StaffSalary,
		Rent:           req.Rent,
		Utilities:      req.Utilities,
		MarketingCosts: req.MarketingCosts,
		OtherExpenses:  req.OtherExpenses,
	}
	
	// Calculate total expenses
	operatingExpense.CalculateTotalExpenses()
	
	// Upsert (create or update)
	result, err := s.expenseRepo.Upsert(ctx, operatingExpense)
	if err != nil {
		return nil, fmt.Errorf("failed to upsert operating expense: %w", err)
	}
	
	return result, nil
}

// GetOperatingExpenseForDate retrieves the operating expense for a specific date
// Returns nil if no expense is found for that date
// Requirements: 6.5.7
func (s *OperatingExpenseService) GetOperatingExpenseForDate(ctx context.Context, date time.Time) (*expense.OperatingExpense, error) {
	operatingExpense, err := s.expenseRepo.FindForDate(ctx, date)
	if err != nil {
		// Return nil if not found (not an error condition)
		return nil, nil
	}
	
	return operatingExpense, nil
}

// GetOperatingExpenses retrieves operating expenses with optional date range filter
// If startDate and endDate are nil, returns all expenses
// If provided, returns expenses that overlap with the date range
// Requirements: 6.5.7
func (s *OperatingExpenseService) GetOperatingExpenses(ctx context.Context, startDate, endDate *time.Time) ([]*expense.OperatingExpense, error) {
	// Validate date range if provided
	if startDate != nil && endDate != nil && startDate.After(*endDate) {
		return nil, fmt.Errorf("start_date must be before or equal to end_date")
	}
	
	expenses, err := s.expenseRepo.FindAll(ctx, startDate, endDate)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch operating expenses: %w", err)
	}
	
	return expenses, nil
}

// AllocateDailyExpense allocates a monthly/period expense to a specific day
// Calculates proportional daily expense: monthly_expense / days_in_period
// Requirements: 6.5.8
func (s *OperatingExpenseService) AllocateDailyExpense(ctx context.Context, monthlyExpense *expense.OperatingExpense, targetDate time.Time) (*expense.AllocatedExpense, error) {
	// Validate that targetDate is within the expense period
	if targetDate.Before(monthlyExpense.PeriodStart) || targetDate.After(monthlyExpense.PeriodEnd) {
		return nil, fmt.Errorf("target date %s is outside expense period %s to %s",
			targetDate.Format("2006-01-02"),
			monthlyExpense.PeriodStart.Format("2006-01-02"),
			monthlyExpense.PeriodEnd.Format("2006-01-02"))
	}
	
	// Calculate number of days in the period
	daysInPeriod := monthlyExpense.PeriodEnd.Sub(monthlyExpense.PeriodStart).Hours()/24 + 1 // +1 to include both start and end dates
	
	if daysInPeriod <= 0 {
		return nil, fmt.Errorf("invalid period: period_end must be after period_start")
	}
	
	// Calculate daily allocation for each expense category
	dailyStaffSalary := monthlyExpense.StaffSalary / daysInPeriod
	dailyRent := monthlyExpense.Rent / daysInPeriod
	dailyUtilities := monthlyExpense.Utilities / daysInPeriod
	dailyMarketingCosts := monthlyExpense.MarketingCosts / daysInPeriod
	dailyOtherExpenses := monthlyExpense.OtherExpenses / daysInPeriod
	dailyTotalExpenses := monthlyExpense.TotalExpenses / daysInPeriod
	
	// Round to 2 decimal places
	dailyStaffSalary = math.Round(dailyStaffSalary*100) / 100
	dailyRent = math.Round(dailyRent*100) / 100
	dailyUtilities = math.Round(dailyUtilities*100) / 100
	dailyMarketingCosts = math.Round(dailyMarketingCosts*100) / 100
	dailyOtherExpenses = math.Round(dailyOtherExpenses*100) / 100
	dailyTotalExpenses = math.Round(dailyTotalExpenses*100) / 100
	
	// Create allocation note
	allocationNote := fmt.Sprintf("Chi phí được phân bổ từ kỳ %s đến %s (%d ngày)",
		monthlyExpense.PeriodStart.Format("02/01/2006"),
		monthlyExpense.PeriodEnd.Format("02/01/2006"),
		int(daysInPeriod))
	
	// Create allocated expense
	allocatedExpense := &expense.AllocatedExpense{
		Date:              targetDate,
		StaffSalary:       dailyStaffSalary,
		Rent:              dailyRent,
		Utilities:         dailyUtilities,
		MarketingCosts:    dailyMarketingCosts,
		OtherExpenses:     dailyOtherExpenses,
		TotalExpenses:     dailyTotalExpenses,
		ExpenseAllocated:  true,
		AllocationNote:    allocationNote,
		SourceExpenseID:   monthlyExpense.ID,
	}
	
	return allocatedExpense, nil
}

// GetOrAllocateExpenseForDate retrieves the operating expense for a date
// If an exact expense exists for that date, returns it
// If a period expense exists that contains the date, allocates daily expense
// If no expense exists, returns nil
func (s *OperatingExpenseService) GetOrAllocateExpenseForDate(ctx context.Context, date time.Time) (*expense.AllocatedExpense, error) {
	// Try to find an expense that contains this date
	operatingExpense, err := s.GetOperatingExpenseForDate(ctx, date)
	if err != nil {
		return nil, err
	}
	
	if operatingExpense == nil {
		// No expense found for this date
		return nil, nil
	}
	
	// Check if this is a single-day expense or a period expense
	periodStart := operatingExpense.PeriodStart
	periodEnd := operatingExpense.PeriodEnd
	
	// Normalize dates to compare only date parts (ignore time)
	periodStartDate := time.Date(periodStart.Year(), periodStart.Month(), periodStart.Day(), 0, 0, 0, 0, time.UTC)
	periodEndDate := time.Date(periodEnd.Year(), periodEnd.Month(), periodEnd.Day(), 0, 0, 0, 0, time.UTC)
	targetDate := time.Date(date.Year(), date.Month(), date.Day(), 0, 0, 0, 0, time.UTC)
	
	if periodStartDate.Equal(periodEndDate) && periodStartDate.Equal(targetDate) {
		// Single-day expense, return as-is (not allocated)
		return &expense.AllocatedExpense{
			Date:              date,
			StaffSalary:       operatingExpense.StaffSalary,
			Rent:              operatingExpense.Rent,
			Utilities:         operatingExpense.Utilities,
			MarketingCosts:    operatingExpense.MarketingCosts,
			OtherExpenses:     operatingExpense.OtherExpenses,
			TotalExpenses:     operatingExpense.TotalExpenses,
			ExpenseAllocated:  false,
			AllocationNote:    "",
			SourceExpenseID:   operatingExpense.ID,
		}, nil
	}
	
	// Period expense, allocate daily
	return s.AllocateDailyExpense(ctx, operatingExpense, date)
}

// GetExpensesForDateRange retrieves or allocates expenses for a date range
// For each day in the range, gets or allocates the expense
// Returns a map of date -> allocated expense
func (s *OperatingExpenseService) GetExpensesForDateRange(ctx context.Context, startDate, endDate time.Time) (map[string]*expense.AllocatedExpense, error) {
	// Validate date range
	if startDate.After(endDate) {
		return nil, fmt.Errorf("start_date must be before or equal to end_date")
	}
	
	// Find all expenses that overlap with the date range
	expenses, err := s.expenseRepo.FindByPeriod(ctx, startDate, endDate)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch expenses: %w", err)
	}
	
	// Create result map
	result := make(map[string]*expense.AllocatedExpense)
	
	// Iterate through each day in the range
	currentDate := startDate
	for !currentDate.After(endDate) {
		dateKey := currentDate.Format("2006-01-02")
		
		// Find the expense that contains this date
		var matchingExpense *expense.OperatingExpense
		for _, exp := range expenses {
			if !currentDate.Before(exp.PeriodStart) && !currentDate.After(exp.PeriodEnd) {
				matchingExpense = exp
				break
			}
		}
		
		if matchingExpense != nil {
			// Allocate expense for this date
			allocated, err := s.AllocateDailyExpense(ctx, matchingExpense, currentDate)
			if err != nil {
				return nil, fmt.Errorf("failed to allocate expense for date %s: %w", dateKey, err)
			}
			result[dateKey] = allocated
		}
		
		// Move to next day
		currentDate = currentDate.AddDate(0, 0, 1)
	}
	
	return result, nil
}
