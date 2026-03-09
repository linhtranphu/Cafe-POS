package services

import (
	"context"
	"errors"
	"fmt"
	"time"

	"cafe-pos/backend/domain/expense"
	"cafe-pos/backend/domain/facility"
	"cafe-pos/backend/domain/fund"
	"cafe-pos/backend/domain/ingredient"
	"cafe-pos/backend/infrastructure/mongodb"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

// Custom errors for fund expense integration
var (
	ErrInsufficientFundBalance  = errors.New("insufficient fund balance")
	ErrInvalidSourceType        = errors.New("invalid source type")
	ErrDuplicateFundTransaction = errors.New("duplicate fund transaction for source")
	ErrPaymentMethodMismatch    = errors.New("payment method does not match fund payment")
	ErrCannotModifyFundExpense  = errors.New("cannot modify expense paid from fund")
	ErrTransactionRollbackFailed = errors.New("transaction rollback failed")
)

// FundExpenseIntegrationService orchestrates fund-paid operations
type FundExpenseIntegrationService struct {
	expenseRepo         *mongodb.ExpenseRepository
	fundTxRepo          *mongodb.FundTransactionRepository
	ingredientRepo      *mongodb.IngredientRepository
	restockRepo         *mongodb.IngredientRestockRepository
	facilityRepo        *mongodb.FacilityRepository
	fundService         *FundService
	mongoClient         *mongo.Client
}

// NewFundExpenseIntegrationService creates a new service instance
func NewFundExpenseIntegrationService(
	expenseRepo *mongodb.ExpenseRepository,
	fundTxRepo *mongodb.FundTransactionRepository,
	ingredientRepo *mongodb.IngredientRepository,
	restockRepo *mongodb.IngredientRestockRepository,
	facilityRepo *mongodb.FacilityRepository,
	fundService *FundService,
	mongoClient *mongo.Client,
) *FundExpenseIntegrationService {
	return &FundExpenseIntegrationService{
		expenseRepo:    expenseRepo,
		fundTxRepo:     fundTxRepo,
		ingredientRepo: ingredientRepo,
		restockRepo:    restockRepo,
		facilityRepo:   facilityRepo,
		fundService:    fundService,
		mongoClient:    mongoClient,
	}
}

// CreateExpenseFromFundRequest represents the request to create an expense from fund
type CreateExpenseFromFundRequest struct {
	Date          time.Time
	CategoryID    primitive.ObjectID
	Amount        float64
	Description   string
	Vendor        string
	Notes         string
	MoneyType     string // "cash" or "transfer"
	UserID        primitive.ObjectID
	UserName      string
	UserRole      string
}

// RestockFromFundRequest represents the request to restock ingredient from fund
type RestockFromFundRequest struct {
	IngredientID primitive.ObjectID
	Quantity     float64
	CostPerUnit  float64
	Reason       string
	MoneyType    string // "cash" or "transfer"
	UserID       primitive.ObjectID
	UserName     string
	UserRole     string
}

// ExpenseFromFundResult represents the result of creating an expense from fund
type ExpenseFromFundResult struct {
	Expense         *expense.Expense
	FundTransaction *fund.FundTransaction
}

// RestockFromFundResult represents the result of restocking ingredient from fund
type RestockFromFundResult struct {
	RestockRecord   *ingredient.IngredientRestockRecord
	Expense         *expense.Expense
	FundTransaction *fund.FundTransaction
	UpdatedStock    float64
}

// ValidateFundBalance validates that the fund has sufficient balance for the requested amount
// Requirements: 6.1, 6.2, 6.3
func (s *FundExpenseIntegrationService) ValidateFundBalance(ctx context.Context, requiredAmount float64) error {
	// Get current fund balance
	currentBalance, err := s.fundService.CalculateCurrentBalance(ctx)
	if err != nil {
		return fmt.Errorf("failed to calculate fund balance: %w", err)
	}

	// Check if balance is sufficient
	if requiredAmount > currentBalance.Total {
		return fmt.Errorf("%w: required=%.2f, available=%.2f", 
			ErrInsufficientFundBalance, requiredAmount, currentBalance.Total)
	}

	return nil
}

// CreateExpenseFromFund creates an expense paid from fund with atomicity guarantee
// Requirements: 1.1, 1.2, 1.3, 1.4, 1.5, 1.6, 9.1, 9.2, 10.1
func (s *FundExpenseIntegrationService) CreateExpenseFromFund(
	ctx context.Context,
	req CreateExpenseFromFundRequest,
) (*ExpenseFromFundResult, error) {
	// Validate amount
	if req.Amount <= 0 {
		return nil, errors.New("amount must be greater than 0")
	}

	// Default money_type to cash
	if req.MoneyType == "" {
		req.MoneyType = "cash"
	}

	// Validate fund balance before starting transaction
	// Requirement 6.1, 6.2: Check balance before processing
	if err := s.ValidateFundBalance(ctx, req.Amount); err != nil {
		return nil, err
	}

	// Start MongoDB session for transaction
	session, err := s.mongoClient.StartSession()
	if err != nil {
		return nil, fmt.Errorf("failed to start session: %w", err)
	}
	defer session.EndSession(ctx)

	var result *ExpenseFromFundResult

	// Execute transaction
	err = mongo.WithSession(ctx, session, func(sc mongo.SessionContext) error {
		if err := session.StartTransaction(); err != nil {
			return err
		}

		// Get current balance for audit
		balanceBefore, err := s.fundService.CalculateCurrentBalance(sc)
		if err != nil {
			session.AbortTransaction(sc)
			return fmt.Errorf("failed to calculate balance: %w", err)
		}

		// Double-check balance within transaction and validate per money type
		// Requirement 6.2: Validate sufficient balance
		if req.MoneyType == "transfer" {
			if req.Amount > balanceBefore.Transfer {
				session.AbortTransaction(sc)
				return fmt.Errorf("insufficient fund balance (chuyển khoản): cần=%.2f, có=%.2f",
					req.Amount, balanceBefore.Transfer)
			}
		} else {
			if req.Amount > balanceBefore.Cash {
				session.AbortTransaction(sc)
				return fmt.Errorf("insufficient fund balance (tiền mặt): cần=%.2f, có=%.2f",
					req.Amount, balanceBefore.Cash)
			}
		}

		// Create expense record
		// Requirements 1.1, 1.5, 10.1: Create expense with fund fields
		exp := &expense.Expense{
			Date:          req.Date,
			CategoryID:    req.CategoryID,
			Amount:        req.Amount,
			Description:   req.Description,
			PaymentMethod: expense.PaymentMethodFund, // Requirement 10.1
			Vendor:        req.Vendor,
			Notes:         req.Notes,
			SourceType:    expense.SourceTypeManual,
			PaidFromFund:  true, // Requirement 1.1
			CreatedBy:     req.UserName,
			CreatedAt:     time.Now(),
			UpdatedAt:     time.Now(),
		}

		// Create expense
		if err := s.expenseRepo.CreateExpense(sc, exp); err != nil {
			session.AbortTransaction(sc)
			return fmt.Errorf("failed to create expense: %w", err) // Requirement 9.1
		}

		// Calculate balance after withdrawal based on money type
		var cashAmount, transferAmount float64
		if req.MoneyType == "transfer" {
			transferAmount = req.Amount
		} else {
			cashAmount = req.Amount
		}
		balanceAfter := &fund.FundBalance{
			Cash:     balanceBefore.Cash - cashAmount,
			Transfer: balanceBefore.Transfer - transferAmount,
			Total:    balanceBefore.Total - req.Amount,
		}

		// Create fund withdrawal transaction with source linking
		// Requirements 1.3, 1.4, 1.6: Create withdrawal with source reference
		fundTx, err := fund.NewFundTransaction(
			fund.TransactionTypeWithdrawal,
			cashAmount,
			transferAmount,
			fmt.Sprintf("Chi tiêu: %s", req.Description),
			req.UserID,
			req.UserName,
			req.UserRole,
		)
		if err != nil {
			session.AbortTransaction(sc)
			return fmt.Errorf("failed to create fund transaction: %w", err)
		}

		// Set source linking - Requirement 1.6
		if err := fundTx.SetSource(fund.SourceTypeExpense, exp.ID); err != nil {
			session.AbortTransaction(sc)
			return fmt.Errorf("failed to set fund transaction source: %w", err)
		}

		// Set balance information for audit - Requirement 1.4
		fundTx.SetBalances(balanceBefore, balanceAfter)

		// Create fund transaction
		if err := s.fundTxRepo.Create(sc, fundTx); err != nil {
			session.AbortTransaction(sc)
			return fmt.Errorf("failed to create fund transaction: %w", err) // Requirement 9.2
		}

		// Update expense with fund transaction ID - Requirement 1.5
		exp.FundTransactionID = fundTx.ID
		if err := s.expenseRepo.UpdateExpense(sc, exp.ID, exp); err != nil {
			session.AbortTransaction(sc)
			return fmt.Errorf("failed to update expense with fund transaction ID: %w", err)
		}

		// Commit transaction
		if err := session.CommitTransaction(sc); err != nil {
			return fmt.Errorf("failed to commit transaction: %w", err)
		}

		result = &ExpenseFromFundResult{
			Expense:         exp,
			FundTransaction: fundTx,
		}

		return nil
	})

	if err != nil {
		return nil, err
	}

	return result, nil
}

// RestockIngredientFromFund restocks ingredient paid from fund with atomicity guarantee
// Requirements: 2.1, 2.2, 2.3, 2.4, 2.5, 2.6, 2.7, 9.3
func (s *FundExpenseIntegrationService) RestockIngredientFromFund(
	ctx context.Context,
	req RestockFromFundRequest,
) (*RestockFromFundResult, error) {
	// Validate inputs
	if req.Quantity <= 0 {
		return nil, errors.New("quantity must be greater than 0")
	}
	if req.CostPerUnit < 0 {
		return nil, errors.New("cost per unit cannot be negative")
	}

	totalCost := req.Quantity * req.CostPerUnit

	// Validate fund balance before starting transaction
	// Requirement 2.2: Check balance before processing
	if err := s.ValidateFundBalance(ctx, totalCost); err != nil {
		return nil, err
	}

	// Get ingredient to verify it exists
	ing, err := s.ingredientRepo.FindByID(ctx, req.IngredientID)
	if err != nil {
		return nil, fmt.Errorf("failed to find ingredient: %w", err)
	}

	// Start MongoDB session for transaction
	session, err := s.mongoClient.StartSession()
	if err != nil {
		return nil, fmt.Errorf("failed to start session: %w", err)
	}
	defer session.EndSession(ctx)

	var result *RestockFromFundResult

	// Execute transaction
	// Requirement 9.3: Use transaction for atomicity
	err = mongo.WithSession(ctx, session, func(sc mongo.SessionContext) error {
		if err := session.StartTransaction(); err != nil {
			return err
		}

		// Get current balance for audit
		balanceBefore, err := s.fundService.CalculateCurrentBalance(sc)
		if err != nil {
			session.AbortTransaction(sc)
			return fmt.Errorf("failed to calculate balance: %w", err)
		}

		// Double-check balance within transaction based on money_type
		if req.MoneyType == "transfer" {
			if totalCost > balanceBefore.Transfer {
				session.AbortTransaction(sc)
				return fmt.Errorf("%w (transfer): required=%.2f, available=%.2f", 
					ErrInsufficientFundBalance, totalCost, balanceBefore.Transfer)
			}
		} else {
			if totalCost > balanceBefore.Cash {
				session.AbortTransaction(sc)
				return fmt.Errorf("%w (cash): required=%.2f, available=%.2f", 
					ErrInsufficientFundBalance, totalCost, balanceBefore.Cash)
			}
		}

		// Create ingredient restock record
		// Requirement 2.1: Create restock record
		restockRecord, err := ingredient.NewIngredientRestockRecord(
			req.IngredientID,
			req.Quantity,
			req.CostPerUnit,
			req.UserID.Hex(),
			req.UserName,
			req.Reason,
		)
		if err != nil {
			session.AbortTransaction(sc)
			return fmt.Errorf("failed to create restock record: %w", err)
		}
		
		// Create restock record in DB first to get ID
		if err := s.restockRepo.Create(sc, restockRecord); err != nil {
			session.AbortTransaction(sc)
			return fmt.Errorf("failed to create restock record: %w", err)
		}

		// Find or create "ingredient purchase" expense category
		categories, err := s.expenseRepo.GetCategories(sc)
		if err != nil {
			session.AbortTransaction(sc)
			return fmt.Errorf("failed to get expense categories: %w", err)
		}

		var ingredientCategoryID primitive.ObjectID
		for _, cat := range categories {
			if cat.Name == "Mua nguyên liệu" || cat.Name == "ingredient purchase" {
				ingredientCategoryID = cat.ID
				break
			}
		}

		// If category doesn't exist, create it
		if ingredientCategoryID.IsZero() {
			cat := &expense.Category{
				Name: "Mua nguyên liệu",
			}
			if err := s.expenseRepo.CreateCategory(sc, cat); err != nil {
				session.AbortTransaction(sc)
				return fmt.Errorf("failed to create ingredient category: %w", err)
			}
			ingredientCategoryID = cat.ID
		}

		// Create expense record for ingredient purchase
		// Requirement 2.3: Create expense with category "ingredient purchase"
		exp := &expense.Expense{
			Date:          time.Now(),
			CategoryID:    ingredientCategoryID,
			Amount:        totalCost,
			Description:   fmt.Sprintf("Mua nguyên liệu: %s (%.2f %s)", ing.Name, req.Quantity, ing.Unit),
			PaymentMethod: expense.PaymentMethodFund,
			Vendor:        ing.Supplier,
			Notes:         req.Reason,
			SourceType:    expense.SourceTypeIngredient,
			SourceID:      restockRecord.ID,
			PaidFromFund:  true,
			CreatedBy:     req.UserName,
			CreatedAt:     time.Now(),
			UpdatedAt:     time.Now(),
		}

		if err := s.expenseRepo.CreateExpense(sc, exp); err != nil {
			session.AbortTransaction(sc)
			return fmt.Errorf("failed to create expense: %w", err)
		}

		// Calculate balance after withdrawal based on money_type
		var cashAmount, transferAmount float64
		if req.MoneyType == "transfer" {
			transferAmount = totalCost
			cashAmount = 0
		} else {
			cashAmount = totalCost
			transferAmount = 0
		}
		
		balanceAfter := &fund.FundBalance{
			Cash:     balanceBefore.Cash - cashAmount,
			Transfer: balanceBefore.Transfer - transferAmount,
			Total:    balanceBefore.Total - totalCost,
		}

		// Create fund withdrawal transaction
		// Requirement 2.4: Create withdrawal transaction
		fundTx, err := fund.NewFundTransaction(
			fund.TransactionTypeWithdrawal,
			cashAmount,
			transferAmount,
			fmt.Sprintf("Mua nguyên liệu: %s", ing.Name),
			req.UserID,
			req.UserName,
			req.UserRole,
		)
		if err != nil {
			session.AbortTransaction(sc)
			return fmt.Errorf("failed to create fund transaction: %w", err)
		}

		// Set source linking - Requirement 2.6
		if err := fundTx.SetSource(fund.SourceTypeIngredient, restockRecord.ID); err != nil {
			session.AbortTransaction(sc)
			return fmt.Errorf("failed to set fund transaction source: %w", err)
		}

		fundTx.SetBalances(balanceBefore, balanceAfter)

		if err := s.fundTxRepo.Create(sc, fundTx); err != nil {
			session.AbortTransaction(sc)
			return fmt.Errorf("failed to create fund transaction: %w", err)
		}

		// Link restock record to expense and fund transaction
		// Requirement 2.7: Link expense to fund transaction
		if err := restockRecord.SetFundPayment(exp.ID, fundTx.ID); err != nil {
			session.AbortTransaction(sc)
			return fmt.Errorf("failed to set fund payment: %w", err)
		}

		// Update expense with fund transaction ID
		exp.FundTransactionID = fundTx.ID
		if err := s.expenseRepo.UpdateExpense(sc, exp.ID, exp); err != nil {
			session.AbortTransaction(sc)
			return fmt.Errorf("failed to update expense: %w", err)
		}

		// Update ingredient stock quantity
		// Requirement 2.5: Update stock quantity
		beforeQty := ing.Quantity
		ing.Quantity += req.Quantity
		
		// Update cost per unit using weighted average if new price is different
		if req.CostPerUnit > 0 && req.CostPerUnit != ing.CostPerUnit && ing.Quantity > 0 {
			oldValue := beforeQty * ing.CostPerUnit
			newValue := req.Quantity * req.CostPerUnit
			ing.CostPerUnit = (oldValue + newValue) / ing.Quantity
		}

		if err := s.ingredientRepo.Update(sc, req.IngredientID, ing); err != nil {
			session.AbortTransaction(sc)
			return fmt.Errorf("failed to update ingredient stock: %w", err)
		}

		// Commit transaction
		if err := session.CommitTransaction(sc); err != nil {
			return fmt.Errorf("failed to commit transaction: %w", err)
		}

		result = &RestockFromFundResult{
			RestockRecord:   restockRecord,
			Expense:         exp,
			FundTransaction: fundTx,
			UpdatedStock:    ing.Quantity,
		}

		return nil
	})

	if err != nil {
		return nil, err
	}

	return result, nil
}

// GetExpensesPaidFromFund retrieves expenses paid from fund with filtering
// Requirements: 4.3, 4.4
func (s *FundExpenseIntegrationService) GetExpensesPaidFromFund(
	ctx context.Context,
	limit, offset int,
) ([]expense.Expense, error) {
	filter := bson.M{
		"paid_from_fund": true,
		"fund_transaction_id": bson.M{"$exists": true, "$ne": primitive.NilObjectID},
	}

	expenses, err := s.expenseRepo.GetExpenses(ctx, filter)
	if err != nil {
		return nil, fmt.Errorf("failed to get expenses paid from fund: %w", err)
	}

	return expenses, nil
}

// PurchaseFacilityFromFundRequest represents the request to purchase a facility from fund
type PurchaseFacilityFromFundRequest struct {
	Facility  *facility.Facility
	MoneyType string // "cash" or "transfer"
	UserID    primitive.ObjectID
	UserName  string
	UserRole  string
}

// PurchaseFacilityFromFundResult represents the result of purchasing a facility from fund
type PurchaseFacilityFromFundResult struct {
	Facility        *facility.Facility
	Expense         *expense.Expense
	FundTransaction *fund.FundTransaction
}

// PurchaseFacilityFromFund creates a facility paid from fund with atomicity guarantee
func (s *FundExpenseIntegrationService) PurchaseFacilityFromFund(
	ctx context.Context,
	req PurchaseFacilityFromFundRequest,
) (*PurchaseFacilityFromFundResult, error) {
	if req.Facility.Cost <= 0 {
		return nil, errors.New("giá trị tài sản phải lớn hơn 0 để trừ quỹ")
	}

	totalCost := req.Facility.Cost

	// Validate fund balance before starting transaction
	if err := s.ValidateFundBalance(ctx, totalCost); err != nil {
		return nil, err
	}

	session, err := s.mongoClient.StartSession()
	if err != nil {
		return nil, fmt.Errorf("failed to start session: %w", err)
	}
	defer session.EndSession(ctx)

	var result *PurchaseFacilityFromFundResult

	err = mongo.WithSession(ctx, session, func(sc mongo.SessionContext) error {
		if err := session.StartTransaction(); err != nil {
			return err
		}

		// Get current balance for audit
		balanceBefore, err := s.fundService.CalculateCurrentBalance(sc)
		if err != nil {
			session.AbortTransaction(sc)
			return fmt.Errorf("failed to calculate balance: %w", err)
		}

		// Double-check balance within transaction based on money_type
		if req.MoneyType == "transfer" {
			if totalCost > balanceBefore.Transfer {
				session.AbortTransaction(sc)
				return fmt.Errorf("%w (chuyển khoản): cần=%.2f, có=%.2f",
					ErrInsufficientFundBalance, totalCost, balanceBefore.Transfer)
			}
		} else {
			if totalCost > balanceBefore.Cash {
				session.AbortTransaction(sc)
				return fmt.Errorf("%w (tiền mặt): cần=%.2f, có=%.2f",
					ErrInsufficientFundBalance, totalCost, balanceBefore.Cash)
			}
		}

		// Create facility record
		now := time.Now()
		req.Facility.PaidFromFund = true
		req.Facility.CreatedAt = now
		req.Facility.UpdatedAt = now
		if err := s.facilityRepo.Create(sc, req.Facility); err != nil {
			session.AbortTransaction(sc)
			return fmt.Errorf("failed to create facility: %w", err)
		}

		// Find or create "facility purchase" expense category
		categories, err := s.expenseRepo.GetCategories(sc)
		if err != nil {
			session.AbortTransaction(sc)
			return fmt.Errorf("failed to get expense categories: %w", err)
		}

		var facilityCategoryID primitive.ObjectID
		for _, cat := range categories {
			if cat.Name == "Cơ sở vật chất" || cat.Name == "facility purchase" {
				facilityCategoryID = cat.ID
				break
			}
		}
		if facilityCategoryID.IsZero() {
			cat := &expense.Category{Name: "Cơ sở vật chất"}
			if err := s.expenseRepo.CreateCategory(sc, cat); err != nil {
				session.AbortTransaction(sc)
				return fmt.Errorf("failed to create facility category: %w", err)
			}
			facilityCategoryID = cat.ID
		}

		// Create expense record
		exp := &expense.Expense{
			Date:          now,
			CategoryID:    facilityCategoryID,
			Amount:        totalCost,
			Description:   fmt.Sprintf("Mua tài sản: %s (SL: %d)", req.Facility.Name, req.Facility.Quantity),
			PaymentMethod: expense.PaymentMethodFund,
			Vendor:        req.Facility.Supplier,
			Notes:         req.Facility.Notes,
			SourceType:    expense.SourceTypeFacility,
			SourceID:      req.Facility.ID,
			PaidFromFund:  true,
			CreatedBy:     req.UserName,
			CreatedAt:     now,
			UpdatedAt:     now,
		}
		if err := s.expenseRepo.CreateExpense(sc, exp); err != nil {
			session.AbortTransaction(sc)
			return fmt.Errorf("failed to create expense: %w", err)
		}

		// Calculate balance after withdrawal
		var cashAmount, transferAmount float64
		if req.MoneyType == "transfer" {
			transferAmount = totalCost
		} else {
			cashAmount = totalCost
		}
		balanceAfter := &fund.FundBalance{
			Cash:     balanceBefore.Cash - cashAmount,
			Transfer: balanceBefore.Transfer - transferAmount,
			Total:    balanceBefore.Total - totalCost,
		}

		// Create fund withdrawal transaction
		fundTx, err := fund.NewFundTransaction(
			fund.TransactionTypeWithdrawal,
			cashAmount,
			transferAmount,
			fmt.Sprintf("Mua tài sản: %s", req.Facility.Name),
			req.UserID,
			req.UserName,
			req.UserRole,
		)
		if err != nil {
			session.AbortTransaction(sc)
			return fmt.Errorf("failed to create fund transaction: %w", err)
		}

		if err := fundTx.SetSource(fund.SourceTypeFacility, req.Facility.ID); err != nil {
			session.AbortTransaction(sc)
			return fmt.Errorf("failed to set fund transaction source: %w", err)
		}
		fundTx.SetBalances(balanceBefore, balanceAfter)

		if err := s.fundTxRepo.Create(sc, fundTx); err != nil {
			session.AbortTransaction(sc)
			return fmt.Errorf("failed to create fund transaction: %w", err)
		}

		// Link facility to expense and fund transaction, then update
		req.Facility.ExpenseID = exp.ID
		req.Facility.FundTransactionID = fundTx.ID
		exp.FundTransactionID = fundTx.ID

		if err := s.facilityRepo.Update(sc, req.Facility.ID, req.Facility); err != nil {
			session.AbortTransaction(sc)
			return fmt.Errorf("failed to update facility with fund links: %w", err)
		}
		if err := s.expenseRepo.UpdateExpense(sc, exp.ID, exp); err != nil {
			session.AbortTransaction(sc)
			return fmt.Errorf("failed to update expense with fund transaction ID: %w", err)
		}

		if err := session.CommitTransaction(sc); err != nil {
			return fmt.Errorf("failed to commit transaction: %w", err)
		}

		result = &PurchaseFacilityFromFundResult{
			Facility:        req.Facility,
			Expense:         exp,
			FundTransaction: fundTx,
		}
		return nil
	})

	if err != nil {
		return nil, err
	}

	return result, nil
}

// GetRestockHistory retrieves restock history for an ingredient with pagination
// Requirements: 2.5, 2.6, 2.7
func (s *FundExpenseIntegrationService) GetRestockHistory(
	ctx context.Context,
	ingredientID primitive.ObjectID,
	limit, offset int,
) ([]*ingredient.IngredientRestockRecord, error) {
	// Query ingredient_restock_history collection
	records, err := s.restockRepo.FindByIngredientID(ctx, ingredientID, limit, offset)
	if err != nil {
		return nil, fmt.Errorf("failed to get restock history: %w", err)
	}

	return records, nil
}
