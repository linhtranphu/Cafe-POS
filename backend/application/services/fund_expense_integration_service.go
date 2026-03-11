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
	ErrInsufficientFundBalance = errors.New("insufficient fund balance")
	ErrInvalidSourceType       = errors.New("invalid source type")
	ErrPaymentMethodMismatch   = errors.New("payment method does not match fund payment")
	ErrCannotModifyFundExpense = errors.New("cannot modify expense paid from fund")
)

// FundExpenseIntegrationService orchestrates fund-paid operations
type FundExpenseIntegrationService struct {
	expenseRepo    *mongodb.ExpenseRepository
	ingredientRepo *mongodb.IngredientRepository
	restockRepo    *mongodb.IngredientRestockRepository
	facilityRepo   *mongodb.FacilityRepository
	journalService *JournalService
	mongoClient    *mongo.Client
}

// NewFundExpenseIntegrationService creates a new service instance
func NewFundExpenseIntegrationService(
	expenseRepo *mongodb.ExpenseRepository,
	ingredientRepo *mongodb.IngredientRepository,
	restockRepo *mongodb.IngredientRestockRepository,
	facilityRepo *mongodb.FacilityRepository,
	journalService *JournalService,
	mongoClient *mongo.Client,
) *FundExpenseIntegrationService {
	return &FundExpenseIntegrationService{
		expenseRepo:    expenseRepo,
		ingredientRepo: ingredientRepo,
		restockRepo:    restockRepo,
		facilityRepo:   facilityRepo,
		journalService: journalService,
		mongoClient:    mongoClient,
	}
}

// CreateExpenseFromFundRequest represents the request to create an expense from fund
type CreateExpenseFromFundRequest struct {
	Date        time.Time
	CategoryID  primitive.ObjectID
	Amount      float64
	Description string
	Vendor      string
	Notes       string
	MoneyType   string // "cash" or "transfer"
	UserID      primitive.ObjectID
	UserName    string
	UserRole    string
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
	Expense      *expense.Expense
	JournalEntry *fund.JournalEntry
}

// RestockFromFundResult represents the result of restocking ingredient from fund
type RestockFromFundResult struct {
	RestockRecord *ingredient.IngredientRestockRecord
	Expense       *expense.Expense
	JournalEntry  *fund.JournalEntry
	UpdatedStock  float64
}

// ValidateFundBalanceForType validates that a specific fund has sufficient balance
func (s *FundExpenseIntegrationService) ValidateFundBalanceForType(ctx context.Context, requiredAmount float64, moneyType string, fundType fund.FundType) error {
	b, err := s.journalService.GetFundBalance(ctx, fundType)
	if err != nil {
		return fmt.Errorf("failed to calculate fund balance: %w", err)
	}
	if moneyType == "transfer" {
		if requiredAmount > b.Transfer {
			return fmt.Errorf("%w (%s - chuyển khoản): cần=%.2f, có=%.2f", ErrInsufficientFundBalance, fundType, requiredAmount, b.Transfer)
		}
	} else {
		if requiredAmount > b.Cash {
			return fmt.Errorf("%w (%s - tiền mặt): cần=%.2f, có=%.2f", ErrInsufficientFundBalance, fundType, requiredAmount, b.Cash)
		}
	}
	return nil
}

// CreateExpenseFromFund creates an expense paid from fund with atomicity guarantee
func (s *FundExpenseIntegrationService) CreateExpenseFromFund(
	ctx context.Context,
	req CreateExpenseFromFundRequest,
) (*ExpenseFromFundResult, error) {
	if req.Amount <= 0 {
		return nil, errors.New("amount must be greater than 0")
	}
	if req.MoneyType == "" {
		req.MoneyType = "cash"
	}

	expenseFundType := fund.FundTypeForSource(fund.SourceTypeExpense)
	if err := s.ValidateFundBalanceForType(ctx, req.Amount, req.MoneyType, expenseFundType); err != nil {
		return nil, err
	}

	session, err := s.mongoClient.StartSession()
	if err != nil {
		return nil, fmt.Errorf("failed to start session: %w", err)
	}
	defer session.EndSession(ctx)

	var result *ExpenseFromFundResult

	err = mongo.WithSession(ctx, session, func(sc mongo.SessionContext) error {
		if err := session.StartTransaction(); err != nil {
			return err
		}

		// Create expense record
		exp := &expense.Expense{
			Date:          req.Date,
			CategoryID:    req.CategoryID,
			Amount:        req.Amount,
			Description:   req.Description,
			PaymentMethod: expense.PaymentMethodFund,
			Vendor:        req.Vendor,
			Notes:         req.Notes,
			SourceType:    expense.SourceTypeManual,
			PaidFromFund:  true,
			CreatedBy:     req.UserName,
			CreatedAt:     time.Now(),
			UpdatedAt:     time.Now(),
		}
		if err := s.expenseRepo.CreateExpense(sc, exp); err != nil {
			session.AbortTransaction(sc)
			return fmt.Errorf("failed to create expense: %w", err)
		}

		// Calculate cash/transfer split
		var cashAmount, transferAmount float64
		if req.MoneyType == "transfer" {
			transferAmount = req.Amount
		} else {
			cashAmount = req.Amount
		}

		// Record fund withdrawal in journal (within same session/transaction)
		desc := fmt.Sprintf("Chi tiêu: %s", req.Description)
		je, err := s.journalService.RecordFundWithdrawalInSession(
			sc,
			fund.EventExpense,
			expenseFundType,
			fund.FundTypeSupplier, // chi tiêu → trả nhà cung cấp
			cashAmount, transferAmount,
			desc,
			exp.ID,
			req.UserID, req.UserName, req.UserRole,
		)
		if err != nil {
			session.AbortTransaction(sc)
			return fmt.Errorf("failed to record journal entry: %w", err)
		}

		// Store journal entry ID in expense
		exp.FundTransactionID = je.ID
		if err := s.expenseRepo.UpdateExpense(sc, exp.ID, exp); err != nil {
			session.AbortTransaction(sc)
			return fmt.Errorf("failed to update expense with journal entry ID: %w", err)
		}

		if err := session.CommitTransaction(sc); err != nil {
			return fmt.Errorf("failed to commit transaction: %w", err)
		}

		result = &ExpenseFromFundResult{Expense: exp, JournalEntry: je}
		return nil
	})

	if err != nil {
		return nil, err
	}
	return result, nil
}

// RestockIngredientFromFund restocks ingredient paid from fund with atomicity guarantee
func (s *FundExpenseIntegrationService) RestockIngredientFromFund(
	ctx context.Context,
	req RestockFromFundRequest,
) (*RestockFromFundResult, error) {
	if req.Quantity <= 0 {
		return nil, errors.New("quantity must be greater than 0")
	}
	if req.CostPerUnit < 0 {
		return nil, errors.New("cost per unit cannot be negative")
	}

	totalCost := req.Quantity * req.CostPerUnit

	ingredientFundType := fund.FundTypeForSource(fund.SourceTypeIngredient)
	if err := s.ValidateFundBalanceForType(ctx, totalCost, req.MoneyType, ingredientFundType); err != nil {
		return nil, err
	}

	ing, err := s.ingredientRepo.FindByID(ctx, req.IngredientID)
	if err != nil {
		return nil, fmt.Errorf("failed to find ingredient: %w", err)
	}

	session, err := s.mongoClient.StartSession()
	if err != nil {
		return nil, fmt.Errorf("failed to start session: %w", err)
	}
	defer session.EndSession(ctx)

	var result *RestockFromFundResult

	err = mongo.WithSession(ctx, session, func(sc mongo.SessionContext) error {
		if err := session.StartTransaction(); err != nil {
			return err
		}

		// Create restock record
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
		if err := s.restockRepo.Create(sc, restockRecord); err != nil {
			session.AbortTransaction(sc)
			return fmt.Errorf("failed to create restock record: %w", err)
		}

		// Find or create "ingredient purchase" expense category
		categories, err := s.expenseRepo.GetCategories(sc, "")
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
		if ingredientCategoryID.IsZero() {
			cat := &expense.Category{Name: "Mua nguyên liệu"}
			if err := s.expenseRepo.CreateCategory(sc, cat); err != nil {
				session.AbortTransaction(sc)
				return fmt.Errorf("failed to create ingredient category: %w", err)
			}
			ingredientCategoryID = cat.ID
		}

		// Create expense record
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

		// Calculate cash/transfer split
		var cashAmount, transferAmount float64
		if req.MoneyType == "transfer" {
			transferAmount = totalCost
		} else {
			cashAmount = totalCost
		}

		// Record fund withdrawal in journal
		desc := fmt.Sprintf("Mua nguyên liệu: %s", ing.Name)
		je, err := s.journalService.RecordFundWithdrawalInSession(
			sc,
			fund.EventIngredientRestock,
			ingredientFundType,
			fund.FundTypeSupplier, // mua nguyên liệu → trả nhà cung cấp
			cashAmount, transferAmount,
			desc,
			restockRecord.ID,
			req.UserID, req.UserName, req.UserRole,
		)
		if err != nil {
			session.AbortTransaction(sc)
			return fmt.Errorf("failed to record journal entry: %w", err)
		}

		// Link restock record to expense and journal entry
		if err := restockRecord.SetFundPayment(exp.ID, je.ID); err != nil {
			session.AbortTransaction(sc)
			return fmt.Errorf("failed to set fund payment: %w", err)
		}
		exp.FundTransactionID = je.ID
		if err := s.expenseRepo.UpdateExpense(sc, exp.ID, exp); err != nil {
			session.AbortTransaction(sc)
			return fmt.Errorf("failed to update expense: %w", err)
		}

		// Update ingredient stock quantity
		beforeQty := ing.Quantity
		ing.Quantity += req.Quantity
		if req.CostPerUnit > 0 && req.CostPerUnit != ing.CostPerUnit && ing.Quantity > 0 {
			oldValue := beforeQty * ing.CostPerUnit
			newValue := req.Quantity * req.CostPerUnit
			ing.CostPerUnit = (oldValue + newValue) / ing.Quantity
		}
		if err := s.ingredientRepo.Update(sc, req.IngredientID, ing); err != nil {
			session.AbortTransaction(sc)
			return fmt.Errorf("failed to update ingredient stock: %w", err)
		}

		if err := session.CommitTransaction(sc); err != nil {
			return fmt.Errorf("failed to commit transaction: %w", err)
		}

		result = &RestockFromFundResult{
			RestockRecord: restockRecord,
			Expense:       exp,
			JournalEntry:  je,
			UpdatedStock:  ing.Quantity,
		}
		return nil
	})

	if err != nil {
		return nil, err
	}
	return result, nil
}

// GetExpensesPaidFromFund retrieves expenses paid from fund with filtering
func (s *FundExpenseIntegrationService) GetExpensesPaidFromFund(
	ctx context.Context,
	limit, offset int,
) ([]expense.Expense, error) {
	filter := bson.M{
		"paid_from_fund":      true,
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
	Facility     *facility.Facility
	Expense      *expense.Expense
	JournalEntry *fund.JournalEntry
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
	facilityFundType := fund.FundTypeForSource(fund.SourceTypeFacility)
	if err := s.ValidateFundBalanceForType(ctx, totalCost, req.MoneyType, facilityFundType); err != nil {
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

		// Create facility record
		now := time.Now()
		req.Facility.PaidFromFund = true
		req.Facility.CreatedAt = now
		req.Facility.UpdatedAt = now
		if err := s.facilityRepo.Create(sc, req.Facility); err != nil {
			session.AbortTransaction(sc)
			return fmt.Errorf("failed to create facility: %w", err)
		}

		// Find or create facility expense category
		categories, err := s.expenseRepo.GetCategories(sc, "")
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

		// Calculate cash/transfer split
		var cashAmount, transferAmount float64
		if req.MoneyType == "transfer" {
			transferAmount = totalCost
		} else {
			cashAmount = totalCost
		}

		// Record fund withdrawal in journal
		desc := fmt.Sprintf("Mua tài sản: %s", req.Facility.Name)
		je, err := s.journalService.RecordFundWithdrawalInSession(
			sc,
			fund.EventFacilityPurchase,
			facilityFundType,
			fund.FundTypeSupplier, // mua tài sản → trả nhà cung cấp
			cashAmount, transferAmount,
			desc,
			req.Facility.ID,
			req.UserID, req.UserName, req.UserRole,
		)
		if err != nil {
			session.AbortTransaction(sc)
			return fmt.Errorf("failed to record journal entry: %w", err)
		}

		// Link facility and expense to journal entry
		req.Facility.ExpenseID = exp.ID
		req.Facility.FundTransactionID = je.ID
		exp.FundTransactionID = je.ID

		if err := s.facilityRepo.Update(sc, req.Facility.ID, req.Facility); err != nil {
			session.AbortTransaction(sc)
			return fmt.Errorf("failed to update facility with journal links: %w", err)
		}
		if err := s.expenseRepo.UpdateExpense(sc, exp.ID, exp); err != nil {
			session.AbortTransaction(sc)
			return fmt.Errorf("failed to update expense with journal entry ID: %w", err)
		}

		if err := session.CommitTransaction(sc); err != nil {
			return fmt.Errorf("failed to commit transaction: %w", err)
		}

		result = &PurchaseFacilityFromFundResult{
			Facility:     req.Facility,
			Expense:      exp,
			JournalEntry: je,
		}
		return nil
	})

	if err != nil {
		return nil, err
	}
	return result, nil
}

// GetRestockHistory retrieves restock history for an ingredient with pagination
func (s *FundExpenseIntegrationService) GetRestockHistory(
	ctx context.Context,
	ingredientID primitive.ObjectID,
	limit, offset int,
) ([]*ingredient.IngredientRestockRecord, error) {
	records, err := s.restockRepo.FindByIngredientID(ctx, ingredientID, limit, offset)
	if err != nil {
		return nil, fmt.Errorf("failed to get restock history: %w", err)
	}
	return records, nil
}
