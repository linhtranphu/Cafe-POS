package expense

import (
	"testing"
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

func TestOperatingExpense_CalculateTotalExpenses(t *testing.T) {
	tests := []struct {
		name     string
		expense  OperatingExpense
		expected float64
	}{
		{
			name: "Calculate total with all expense types",
			expense: OperatingExpense{
				StaffSalary:    2000000,
				Rent:           1000000,
				Utilities:      500000,
				MarketingCosts: 300000,
				OtherExpenses:  200000,
			},
			expected: 4000000,
		},
		{
			name: "Calculate total with zero values",
			expense: OperatingExpense{
				StaffSalary:    0,
				Rent:           0,
				Utilities:      0,
				MarketingCosts: 0,
				OtherExpenses:  0,
			},
			expected: 0,
		},
		{
			name: "Calculate total with partial expenses",
			expense: OperatingExpense{
				StaffSalary:    1500000,
				Rent:           800000,
				Utilities:      0,
				MarketingCosts: 0,
				OtherExpenses:  100000,
			},
			expected: 2400000,
		},
		{
			name: "Calculate total with decimal values",
			expense: OperatingExpense{
				StaffSalary:    1234567.89,
				Rent:           987654.32,
				Utilities:      456789.01,
				MarketingCosts: 234567.89,
				OtherExpenses:  123456.78,
			},
			expected: 3037035.89,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.expense.CalculateTotalExpenses()
			// Use a small tolerance for floating point comparison
			tolerance := 0.01
			diff := tt.expense.TotalExpenses - tt.expected
			if diff < 0 {
				diff = -diff
			}
			if diff > tolerance {
				t.Errorf("CalculateTotalExpenses() = %v, want %v (diff: %v)", tt.expense.TotalExpenses, tt.expected, diff)
			}
		})
	}
}

func TestAllocatedExpense_Structure(t *testing.T) {
	// Test that AllocatedExpense can be created with all fields
	allocated := AllocatedExpense{
		Date:              time.Now(),
		StaffSalary:       66666.67,
		Rent:              33333.33,
		Utilities:         16666.67,
		MarketingCosts:    10000.00,
		OtherExpenses:     6666.67,
		TotalExpenses:     133333.34,
		ExpenseAllocated:  true,
		AllocationNote:    "Chi phí được phân bổ từ tháng",
		SourceExpenseID:   primitive.NewObjectID(),
	}

	if !allocated.ExpenseAllocated {
		t.Error("ExpenseAllocated should be true")
	}

	if allocated.AllocationNote == "" {
		t.Error("AllocationNote should not be empty")
	}

	if allocated.SourceExpenseID.IsZero() {
		t.Error("SourceExpenseID should not be zero")
	}
}

func TestOperatingExpenseRequest_Validation(t *testing.T) {
	// Test that OperatingExpenseRequest structure is correct
	req := OperatingExpenseRequest{
		PeriodStart:    "2024-01-01",
		PeriodEnd:      "2024-01-31",
		StaffSalary:    2000000,
		Rent:           1000000,
		Utilities:      500000,
		MarketingCosts: 300000,
		OtherExpenses:  200000,
	}

	if req.PeriodStart == "" {
		t.Error("PeriodStart should not be empty")
	}

	if req.PeriodEnd == "" {
		t.Error("PeriodEnd should not be empty")
	}

	if req.StaffSalary < 0 {
		t.Error("StaffSalary should not be negative")
	}

	if req.Rent < 0 {
		t.Error("Rent should not be negative")
	}

	if req.Utilities < 0 {
		t.Error("Utilities should not be negative")
	}

	if req.MarketingCosts < 0 {
		t.Error("MarketingCosts should not be negative")
	}

	if req.OtherExpenses < 0 {
		t.Error("OtherExpenses should not be negative")
	}
}
