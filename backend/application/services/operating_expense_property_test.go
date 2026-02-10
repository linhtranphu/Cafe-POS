package services

import (
	"context"
	"math"
	"testing"
	"time"

	"cafe-pos/backend/domain/expense"
	"github.com/leanovate/gopter"
	"github.com/leanovate/gopter/gen"
	"github.com/leanovate/gopter/prop"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

// Feature: menu-cost-profit-analysis, Property 9: Expense Allocation
// **Validates: Requirements 6.5.8**
//
// Property: For any monthly operating expense and target date within that month,
// the allocated daily expense should equal monthly_expense / days_in_month,
// with expense_allocated flag set to true.
func TestProperty_ExpenseAllocation(t *testing.T) {
	parameters := gopter.DefaultTestParameters()
	parameters.MinSuccessfulTests = 100 // Minimum 100 iterations as per spec
	properties := gopter.NewProperties(parameters)

	properties.Property("Daily expense equals monthly expense divided by days in period", prop.ForAll(
		func(expenseData testExpenseData, dayOffset int) bool {
			// Skip invalid cases
			if expenseData.PeriodDays <= 0 || dayOffset < 0 || dayOffset >= expenseData.PeriodDays {
				return true
			}

			ctx := context.Background()

			// Create operating expense
			periodStart := time.Date(2024, 1, 1, 0, 0, 0, 0, time.UTC)
			periodEnd := periodStart.AddDate(0, 0, expenseData.PeriodDays-1)

			operatingExpense := &expense.OperatingExpense{
				ID:             primitive.NewObjectID(),
				PeriodStart:    periodStart,
				PeriodEnd:      periodEnd,
				StaffSalary:    expenseData.StaffSalary,
				Rent:           expenseData.Rent,
				Utilities:      expenseData.Utilities,
				MarketingCosts: expenseData.MarketingCosts,
				OtherExpenses:  expenseData.OtherExpenses,
			}
			operatingExpense.CalculateTotalExpenses()

			// Create service
			expenseRepo := &mockOperatingExpenseRepository{
				expenses: make(map[primitive.ObjectID]*expense.OperatingExpense),
			}
			expenseRepo.expenses[operatingExpense.ID] = operatingExpense
			service := NewOperatingExpenseService(expenseRepo)

			// Calculate target date within the period
			targetDate := periodStart.AddDate(0, 0, dayOffset)

			// Allocate daily expense
			allocated, err := service.AllocateDailyExpense(ctx, operatingExpense, targetDate)
			if err != nil {
				t.Logf("AllocateDailyExpense failed: %v", err)
				return false
			}

			// Calculate expected daily expenses
			daysInPeriod := float64(expenseData.PeriodDays)
			expectedStaffSalary := math.Round((expenseData.StaffSalary/daysInPeriod)*100) / 100
			expectedRent := math.Round((expenseData.Rent/daysInPeriod)*100) / 100
			expectedUtilities := math.Round((expenseData.Utilities/daysInPeriod)*100) / 100
			expectedMarketingCosts := math.Round((expenseData.MarketingCosts/daysInPeriod)*100) / 100
			expectedOtherExpenses := math.Round((expenseData.OtherExpenses/daysInPeriod)*100) / 100
			expectedTotalExpenses := math.Round((operatingExpense.TotalExpenses/daysInPeriod)*100) / 100

			// Verify allocated expenses match expected values
			tolerance := 0.01
			if math.Abs(allocated.StaffSalary-expectedStaffSalary) > tolerance {
				t.Logf("StaffSalary mismatch: expected %.2f, got %.2f",
					expectedStaffSalary, allocated.StaffSalary)
				return false
			}

			if math.Abs(allocated.Rent-expectedRent) > tolerance {
				t.Logf("Rent mismatch: expected %.2f, got %.2f",
					expectedRent, allocated.Rent)
				return false
			}

			if math.Abs(allocated.Utilities-expectedUtilities) > tolerance {
				t.Logf("Utilities mismatch: expected %.2f, got %.2f",
					expectedUtilities, allocated.Utilities)
				return false
			}

			if math.Abs(allocated.MarketingCosts-expectedMarketingCosts) > tolerance {
				t.Logf("MarketingCosts mismatch: expected %.2f, got %.2f",
					expectedMarketingCosts, allocated.MarketingCosts)
				return false
			}

			if math.Abs(allocated.OtherExpenses-expectedOtherExpenses) > tolerance {
				t.Logf("OtherExpenses mismatch: expected %.2f, got %.2f",
					expectedOtherExpenses, allocated.OtherExpenses)
				return false
			}

			if math.Abs(allocated.TotalExpenses-expectedTotalExpenses) > tolerance {
				t.Logf("TotalExpenses mismatch: expected %.2f, got %.2f",
					expectedTotalExpenses, allocated.TotalExpenses)
				return false
			}

			// Verify expense_allocated flag is true
			if !allocated.ExpenseAllocated {
				t.Logf("ExpenseAllocated flag should be true")
				return false
			}

			// Verify allocation note is set
			if allocated.AllocationNote == "" {
				t.Logf("AllocationNote should not be empty")
				return false
			}

			// Verify target date matches
			if !allocated.Date.Equal(targetDate) {
				t.Logf("Date mismatch: expected %s, got %s",
					targetDate.Format("2006-01-02"), allocated.Date.Format("2006-01-02"))
				return false
			}

			// Verify source expense ID matches
			if allocated.SourceExpenseID != operatingExpense.ID {
				t.Logf("SourceExpenseID mismatch")
				return false
			}

			return true
		},
		genExpenseData(),
		gen.IntRange(0, 30), // Day offset within period
	))

	properties.TestingRun(t, gopter.ConsoleReporter(false))
}

// Feature: menu-cost-profit-analysis, Property 9: Expense Allocation (Sum Consistency)
// **Validates: Requirements 6.5.8**
//
// Property: The sum of all daily allocated expenses over a period should approximately
// equal the total monthly expense (within rounding tolerance).
func TestProperty_ExpenseAllocation_SumConsistency(t *testing.T) {
	parameters := gopter.DefaultTestParameters()
	parameters.MinSuccessfulTests = 100
	properties := gopter.NewProperties(parameters)

	properties.Property("Sum of daily allocations equals monthly expense", prop.ForAll(
		func(expenseData testExpenseData) bool {
			// Skip invalid cases
			if expenseData.PeriodDays <= 0 || expenseData.PeriodDays > 31 {
				return true
			}

			ctx := context.Background()

			// Create operating expense
			periodStart := time.Date(2024, 1, 1, 0, 0, 0, 0, time.UTC)
			periodEnd := periodStart.AddDate(0, 0, expenseData.PeriodDays-1)

			operatingExpense := &expense.OperatingExpense{
				ID:             primitive.NewObjectID(),
				PeriodStart:    periodStart,
				PeriodEnd:      periodEnd,
				StaffSalary:    expenseData.StaffSalary,
				Rent:           expenseData.Rent,
				Utilities:      expenseData.Utilities,
				MarketingCosts: expenseData.MarketingCosts,
				OtherExpenses:  expenseData.OtherExpenses,
			}
			operatingExpense.CalculateTotalExpenses()

			// Create service
			expenseRepo := &mockOperatingExpenseRepository{
				expenses: make(map[primitive.ObjectID]*expense.OperatingExpense),
			}
			expenseRepo.expenses[operatingExpense.ID] = operatingExpense
			service := NewOperatingExpenseService(expenseRepo)

			// Allocate for each day in the period
			var sumStaffSalary, sumRent, sumUtilities, sumMarketingCosts, sumOtherExpenses, sumTotal float64

			for dayOffset := 0; dayOffset < expenseData.PeriodDays; dayOffset++ {
				targetDate := periodStart.AddDate(0, 0, dayOffset)
				allocated, err := service.AllocateDailyExpense(ctx, operatingExpense, targetDate)
				if err != nil {
					t.Logf("AllocateDailyExpense failed for day %d: %v", dayOffset, err)
					return false
				}

				sumStaffSalary += allocated.StaffSalary
				sumRent += allocated.Rent
				sumUtilities += allocated.Utilities
				sumMarketingCosts += allocated.MarketingCosts
				sumOtherExpenses += allocated.OtherExpenses
				sumTotal += allocated.TotalExpenses
			}

			// Round sums to 2 decimal places
			sumStaffSalary = math.Round(sumStaffSalary*100) / 100
			sumRent = math.Round(sumRent*100) / 100
			sumUtilities = math.Round(sumUtilities*100) / 100
			sumMarketingCosts = math.Round(sumMarketingCosts*100) / 100
			sumOtherExpenses = math.Round(sumOtherExpenses*100) / 100
			sumTotal = math.Round(sumTotal*100) / 100

			// Verify sums are close to original expenses (within rounding tolerance)
			// Tolerance is higher because of accumulated rounding errors
			tolerance := float64(expenseData.PeriodDays) * 0.01

			if math.Abs(sumStaffSalary-expenseData.StaffSalary) > tolerance {
				t.Logf("StaffSalary sum mismatch: expected %.2f, got %.2f (diff: %.2f)",
					expenseData.StaffSalary, sumStaffSalary, math.Abs(sumStaffSalary-expenseData.StaffSalary))
				return false
			}

			if math.Abs(sumRent-expenseData.Rent) > tolerance {
				t.Logf("Rent sum mismatch: expected %.2f, got %.2f (diff: %.2f)",
					expenseData.Rent, sumRent, math.Abs(sumRent-expenseData.Rent))
				return false
			}

			if math.Abs(sumUtilities-expenseData.Utilities) > tolerance {
				t.Logf("Utilities sum mismatch: expected %.2f, got %.2f (diff: %.2f)",
					expenseData.Utilities, sumUtilities, math.Abs(sumUtilities-expenseData.Utilities))
				return false
			}

			if math.Abs(sumMarketingCosts-expenseData.MarketingCosts) > tolerance {
				t.Logf("MarketingCosts sum mismatch: expected %.2f, got %.2f (diff: %.2f)",
					expenseData.MarketingCosts, sumMarketingCosts, math.Abs(sumMarketingCosts-expenseData.MarketingCosts))
				return false
			}

			if math.Abs(sumOtherExpenses-expenseData.OtherExpenses) > tolerance {
				t.Logf("OtherExpenses sum mismatch: expected %.2f, got %.2f (diff: %.2f)",
					expenseData.OtherExpenses, sumOtherExpenses, math.Abs(sumOtherExpenses-expenseData.OtherExpenses))
				return false
			}

			if math.Abs(sumTotal-operatingExpense.TotalExpenses) > tolerance {
				t.Logf("TotalExpenses sum mismatch: expected %.2f, got %.2f (diff: %.2f)",
					operatingExpense.TotalExpenses, sumTotal, math.Abs(sumTotal-operatingExpense.TotalExpenses))
				return false
			}

			return true
		},
		genExpenseData(),
	))

	properties.TestingRun(t, gopter.ConsoleReporter(false))
}

// Feature: menu-cost-profit-analysis, Property 9: Expense Allocation (Rounding)
// **Validates: Requirements 6.5.8**
//
// Property: All allocated daily expenses should be rounded to 2 decimal places.
func TestProperty_ExpenseAllocation_Rounding(t *testing.T) {
	parameters := gopter.DefaultTestParameters()
	parameters.MinSuccessfulTests = 100
	properties := gopter.NewProperties(parameters)

	properties.Property("Allocated expenses are rounded to 2 decimal places", prop.ForAll(
		func(monthlyExpense float64, periodDays int, dayOffset int) bool {
			// Skip invalid cases
			if monthlyExpense <= 0 || periodDays <= 0 || periodDays > 31 || dayOffset < 0 || dayOffset >= periodDays {
				return true
			}

			ctx := context.Background()

			// Create operating expense with a single expense category
			periodStart := time.Date(2024, 1, 1, 0, 0, 0, 0, time.UTC)
			periodEnd := periodStart.AddDate(0, 0, periodDays-1)

			operatingExpense := &expense.OperatingExpense{
				ID:             primitive.NewObjectID(),
				PeriodStart:    periodStart,
				PeriodEnd:      periodEnd,
				StaffSalary:    monthlyExpense,
				Rent:           0,
				Utilities:      0,
				MarketingCosts: 0,
				OtherExpenses:  0,
			}
			operatingExpense.CalculateTotalExpenses()

			// Create service
			expenseRepo := &mockOperatingExpenseRepository{
				expenses: make(map[primitive.ObjectID]*expense.OperatingExpense),
			}
			expenseRepo.expenses[operatingExpense.ID] = operatingExpense
			service := NewOperatingExpenseService(expenseRepo)

			// Allocate for a specific day
			targetDate := periodStart.AddDate(0, 0, dayOffset)
			allocated, err := service.AllocateDailyExpense(ctx, operatingExpense, targetDate)
			if err != nil {
				t.Logf("AllocateDailyExpense failed: %v", err)
				return false
			}

			// Verify all expense fields are rounded to 2 decimal places
			checkRounding := func(value float64, name string) bool {
				rounded := math.Round(value*100) / 100
				if math.Abs(value-rounded) > 0.0001 {
					t.Logf("%s not properly rounded to 2 decimals: %.10f (rounded: %.2f)",
						name, value, rounded)
					return false
				}
				return true
			}

			if !checkRounding(allocated.StaffSalary, "StaffSalary") {
				return false
			}
			if !checkRounding(allocated.Rent, "Rent") {
				return false
			}
			if !checkRounding(allocated.Utilities, "Utilities") {
				return false
			}
			if !checkRounding(allocated.MarketingCosts, "MarketingCosts") {
				return false
			}
			if !checkRounding(allocated.OtherExpenses, "OtherExpenses") {
				return false
			}
			if !checkRounding(allocated.TotalExpenses, "TotalExpenses") {
				return false
			}

			return true
		},
		gen.Float64Range(1000.0, 10000000.0), // Monthly expense
		gen.IntRange(1, 31),                   // Period days
		gen.IntRange(0, 30),                   // Day offset
	))

	properties.TestingRun(t, gopter.ConsoleReporter(false))
}

// Test data structures for expense allocation property testing

type testExpenseData struct {
	PeriodDays     int
	StaffSalary    float64
	Rent           float64
	Utilities      float64
	MarketingCosts float64
	OtherExpenses  float64
}

// Generator for expense data
func genExpenseData() gopter.Gen {
	return gopter.CombineGens(
		gen.IntRange(1, 31),                  // Period days (1-31)
		gen.Float64Range(1000000, 10000000),  // Staff salary
		gen.Float64Range(500000, 5000000),    // Rent
		gen.Float64Range(200000, 2000000),    // Utilities
		gen.Float64Range(100000, 1000000),    // Marketing costs
		gen.Float64Range(100000, 1000000),    // Other expenses
	).Map(func(values []interface{}) testExpenseData {
		return testExpenseData{
			PeriodDays:     values[0].(int),
			StaffSalary:    values[1].(float64),
			Rent:           values[2].(float64),
			Utilities:      values[3].(float64),
			MarketingCosts: values[4].(float64),
			OtherExpenses:  values[5].(float64),
		}
	})
}

// Mock repository for operating expenses
type mockOperatingExpenseRepository struct {
	expenses map[primitive.ObjectID]*expense.OperatingExpense
}

func (m *mockOperatingExpenseRepository) Create(ctx context.Context, operatingExpense *expense.OperatingExpense) error {
	if operatingExpense.ID.IsZero() {
		operatingExpense.ID = primitive.NewObjectID()
	}
	m.expenses[operatingExpense.ID] = operatingExpense
	return nil
}

func (m *mockOperatingExpenseRepository) Update(ctx context.Context, id primitive.ObjectID, operatingExpense *expense.OperatingExpense) error {
	m.expenses[id] = operatingExpense
	return nil
}

func (m *mockOperatingExpenseRepository) FindByID(ctx context.Context, id primitive.ObjectID) (*expense.OperatingExpense, error) {
	exp, exists := m.expenses[id]
	if !exists {
		return nil, nil
	}
	return exp, nil
}

func (m *mockOperatingExpenseRepository) FindByPeriod(ctx context.Context, startDate, endDate time.Time) ([]*expense.OperatingExpense, error) {
	var result []*expense.OperatingExpense
	for _, exp := range m.expenses {
		// Check if periods overlap
		if !exp.PeriodStart.After(endDate) && !exp.PeriodEnd.Before(startDate) {
			result = append(result, exp)
		}
	}
	return result, nil
}

func (m *mockOperatingExpenseRepository) FindForDate(ctx context.Context, date time.Time) (*expense.OperatingExpense, error) {
	for _, exp := range m.expenses {
		if !date.Before(exp.PeriodStart) && !date.After(exp.PeriodEnd) {
			return exp, nil
		}
	}
	return nil, nil
}

func (m *mockOperatingExpenseRepository) FindAll(ctx context.Context, startDate, endDate *time.Time) ([]*expense.OperatingExpense, error) {
	if startDate == nil || endDate == nil {
		var result []*expense.OperatingExpense
		for _, exp := range m.expenses {
			result = append(result, exp)
		}
		return result, nil
	}
	return m.FindByPeriod(ctx, *startDate, *endDate)
}

func (m *mockOperatingExpenseRepository) Upsert(ctx context.Context, operatingExpense *expense.OperatingExpense) (*expense.OperatingExpense, error) {
	// Check if expense exists for this period
	for id, exp := range m.expenses {
		if exp.PeriodStart.Equal(operatingExpense.PeriodStart) && exp.PeriodEnd.Equal(operatingExpense.PeriodEnd) {
			operatingExpense.ID = id
			m.expenses[id] = operatingExpense
			return operatingExpense, nil
		}
	}

	// Create new
	if operatingExpense.ID.IsZero() {
		operatingExpense.ID = primitive.NewObjectID()
	}
	m.expenses[operatingExpense.ID] = operatingExpense
	return operatingExpense, nil
}

// Legacy methods for backward compatibility
func (m *mockOperatingExpenseRepository) FindByDateRange(ctx context.Context, start, end time.Time) ([]*expense.OperatingExpense, error) {
	return m.FindByPeriod(ctx, start, end)
}

func (m *mockOperatingExpenseRepository) FindByDate(ctx context.Context, date time.Time) (*expense.OperatingExpense, error) {
	return m.FindForDate(ctx, date)
}
