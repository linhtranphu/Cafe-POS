package ingredient

import (
	"errors"
	"time"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

// IngredientRestockRecord represents a record of ingredient restock operation
// This is specifically for tracking restocks that may be paid from fund
type IngredientRestockRecord struct {
	ID           primitive.ObjectID `bson:"_id,omitempty" json:"id"`
	IngredientID primitive.ObjectID `bson:"ingredient_id" json:"ingredient_id"`
	Quantity     float64            `bson:"quantity" json:"quantity"`
	CostPerUnit  float64            `bson:"cost_per_unit" json:"cost_per_unit"`
	TotalCost    float64            `bson:"total_cost" json:"total_cost"`
	
	// Fund integration fields
	PaidFromFund      bool               `bson:"paid_from_fund" json:"paid_from_fund"`
	ExpenseID         primitive.ObjectID `bson:"expense_id,omitempty" json:"expense_id,omitempty"`
	FundTransactionID primitive.ObjectID `bson:"fund_transaction_id,omitempty" json:"fund_transaction_id,omitempty"`
	
	// Audit fields
	PerformedBy     string    `bson:"performed_by" json:"performed_by"`
	PerformedByName string    `bson:"performed_by_name" json:"performed_by_name"`
	Reason          string    `bson:"reason" json:"reason"`
	CreatedAt       time.Time `bson:"created_at" json:"created_at"`
}

// NewIngredientRestockRecord creates a new ingredient restock record
func NewIngredientRestockRecord(
	ingredientID primitive.ObjectID,
	quantity, costPerUnit float64,
	performedBy, performedByName, reason string,
) (*IngredientRestockRecord, error) {
	if ingredientID.IsZero() {
		return nil, errors.New("ingredient_id is required")
	}
	if quantity <= 0 {
		return nil, errors.New("quantity must be greater than 0")
	}
	if costPerUnit < 0 {
		return nil, errors.New("cost_per_unit cannot be negative")
	}
	if performedBy == "" {
		return nil, errors.New("performed_by is required")
	}
	if performedByName == "" {
		return nil, errors.New("performed_by_name is required")
	}
	totalCost := quantity * costPerUnit
	
	return &IngredientRestockRecord{
		IngredientID:    ingredientID,
		Quantity:        quantity,
		CostPerUnit:     costPerUnit,
		TotalCost:       totalCost,
		PerformedBy:     performedBy,
		PerformedByName: performedByName,
		Reason:          reason,
		CreatedAt:       time.Now(),
	}, nil
}

// Validate validates the restock record
// Requirement 2.1, 2.3, 2.4, 2.5: Validate all required fields
func (r *IngredientRestockRecord) Validate() error {
	if r.IngredientID.IsZero() {
		return errors.New("ingredient_id is required")
	}
	if r.Quantity <= 0 {
		return errors.New("quantity must be greater than 0")
	}
	if r.CostPerUnit < 0 {
		return errors.New("cost_per_unit cannot be negative")
	}
	if r.TotalCost < 0 {
		return errors.New("total_cost cannot be negative")
	}
	if r.PerformedBy == "" {
		return errors.New("performed_by is required")
	}
	if r.PerformedByName == "" {
		return errors.New("performed_by_name is required")
	}

	// Validate fund payment consistency
	// Requirement 2.3, 2.4: When paid from fund, expense_id and fund_transaction_id must be set
	if r.PaidFromFund {
		if r.ExpenseID.IsZero() {
			return errors.New("expense_id is required when paid_from_fund is true")
		}
		if r.FundTransactionID.IsZero() {
			return errors.New("fund_transaction_id is required when paid_from_fund is true")
		}
	}
	
	return nil
}

// SetFundPayment sets the fund payment information for the restock record
// Requirement 2.3, 2.4, 2.6, 2.7: Link to expense and fund transaction
func (r *IngredientRestockRecord) SetFundPayment(expenseID, fundTransactionID primitive.ObjectID) error {
	if expenseID.IsZero() {
		return errors.New("expense_id cannot be zero")
	}
	if fundTransactionID.IsZero() {
		return errors.New("fund_transaction_id cannot be zero")
	}
	
	r.PaidFromFund = true
	r.ExpenseID = expenseID
	r.FundTransactionID = fundTransactionID
	
	return nil
}

// IsPaidFromFund checks if the restock was paid from fund
func (r *IngredientRestockRecord) IsPaidFromFund() bool {
	return r.PaidFromFund && !r.ExpenseID.IsZero() && !r.FundTransactionID.IsZero()
}

// CalculateTotalCost calculates and updates the total cost
func (r *IngredientRestockRecord) CalculateTotalCost() {
	r.TotalCost = r.Quantity * r.CostPerUnit
}
