package expense

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

// OperatingExpense represents operating expenses for a specific period
// Used for calculating operating profit after deducting operating costs from gross profit
type OperatingExpense struct {
	ID             primitive.ObjectID `bson:"_id,omitempty" json:"id"`
	PeriodStart    time.Time          `bson:"period_start" json:"period_start"`
	PeriodEnd      time.Time          `bson:"period_end" json:"period_end"`
	StaffSalary    float64            `bson:"staff_salary" json:"staff_salary"`
	Rent           float64            `bson:"rent" json:"rent"`
	Utilities      float64            `bson:"utilities" json:"utilities"`
	MarketingCosts float64            `bson:"marketing_costs" json:"marketing_costs"`
	OtherExpenses  float64            `bson:"other_expenses" json:"other_expenses"`
	TotalExpenses  float64            `bson:"total_expenses" json:"total_expenses"` // Auto-calculated
	CreatedAt      time.Time          `bson:"created_at" json:"created_at"`
	UpdatedAt      time.Time          `bson:"updated_at" json:"updated_at"`
}

// OperatingExpenseRequest represents the request body for creating/updating operating expenses
type OperatingExpenseRequest struct {
	PeriodStart    string  `json:"period_start" binding:"required"`    // ISO 8601 date format
	PeriodEnd      string  `json:"period_end" binding:"required"`      // ISO 8601 date format
	StaffSalary    float64 `json:"staff_salary" binding:"min=0"`       // Must be >= 0
	Rent           float64 `json:"rent" binding:"min=0"`               // Must be >= 0
	Utilities      float64 `json:"utilities" binding:"min=0"`          // Must be >= 0
	MarketingCosts float64 `json:"marketing_costs" binding:"min=0"`    // Must be >= 0
	OtherExpenses  float64 `json:"other_expenses" binding:"min=0"`     // Must be >= 0
}

// CalculateTotalExpenses calculates the total expenses from all expense categories
func (oe *OperatingExpense) CalculateTotalExpenses() {
	oe.TotalExpenses = oe.StaffSalary + oe.Rent + oe.Utilities + oe.MarketingCosts + oe.OtherExpenses
}

// AllocatedExpense represents daily allocated expenses from a monthly/period expense
type AllocatedExpense struct {
	Date              time.Time `json:"date"`
	StaffSalary       float64   `json:"staff_salary"`
	Rent              float64   `json:"rent"`
	Utilities         float64   `json:"utilities"`
	MarketingCosts    float64   `json:"marketing_costs"`
	OtherExpenses     float64   `json:"other_expenses"`
	TotalExpenses     float64   `json:"total_expenses"`
	ExpenseAllocated  bool      `json:"expense_allocated"`  // true if allocated from period
	AllocationNote    string    `json:"allocation_note"`    // e.g., "Chi phí được phân bổ từ tháng"
	SourceExpenseID   primitive.ObjectID `json:"source_expense_id"` // Reference to original expense
}
