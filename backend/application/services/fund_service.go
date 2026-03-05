package services

import (
	"context"
	"fmt"
	"time"

	"cafe-pos/backend/domain/fund"
	"cafe-pos/backend/infrastructure/mongodb"

	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

type FundService struct {
	fundTransactionRepo *mongodb.FundTransactionRepository
	fundHandoverRepo    *mongodb.FundHandoverRepository
	cashHandoverRepo    *mongodb.CashHandoverRepository
	cashierShiftRepo    *mongodb.CashierShiftRepository
	mongoClient         *mongo.Client
}

func NewFundService(
	fundTransactionRepo *mongodb.FundTransactionRepository,
	fundHandoverRepo *mongodb.FundHandoverRepository,
	cashHandoverRepo *mongodb.CashHandoverRepository,
	cashierShiftRepo *mongodb.CashierShiftRepository,
	mongoClient *mongo.Client,
) *FundService {
	return &FundService{
		fundTransactionRepo: fundTransactionRepo,
		fundHandoverRepo:    fundHandoverRepo,
		cashHandoverRepo:    cashHandoverRepo,
		cashierShiftRepo:    cashierShiftRepo,
		mongoClient:         mongoClient,
	}
}

// CalculateCurrentBalance calculates the current fund balance from all sources
func (s *FundService) CalculateCurrentBalance(ctx context.Context) (*fund.FundBalance, error) {
	balance := &fund.FundBalance{
		Cash:     0,
		Transfer: 0,
		Total:    0,
	}

	// 1. Add fund handovers (cashier → fund)
	fundHandovers, _, err := s.fundHandoverRepo.FindByDateRange(ctx, time.Time{}, time.Now(), 1, 10000)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch fund handovers: %w", err)
	}
	for _, fh := range fundHandovers {
		balance.Cash += fh.CashAmount
		balance.Transfer += fh.TransferAmount
	}

	// 2. Add deposits
	deposits, err := s.fundTransactionRepo.FindByDateRange(
		ctx,
		time.Time{}, // from beginning
		time.Now(),
		ptrTransactionType(fund.TransactionTypeDeposit),
		10000, // large limit
		0,
	)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch deposits: %w", err)
	}
	for _, d := range deposits {
		balance.Cash += d.CashAmount
		balance.Transfer += d.TransferAmount
	}

	// 3. Subtract withdrawals
	withdrawals, err := s.fundTransactionRepo.FindByDateRange(
		ctx,
		time.Time{}, // from beginning
		time.Now(),
		ptrTransactionType(fund.TransactionTypeWithdrawal),
		10000, // large limit
		0,
	)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch withdrawals: %w", err)
	}
	for _, w := range withdrawals {
		balance.Cash -= w.CashAmount
		balance.Transfer -= w.TransferAmount
	}

	balance.Total = balance.Cash + balance.Transfer

	return balance, nil
}

// CalculateTodayBalance calculates today's opening balance, inflow, and outflow
func (s *FundService) CalculateTodayBalance(ctx context.Context) (*TodayBalanceSummary, error) {
	now := time.Now()
	startOfDay := time.Date(now.Year(), now.Month(), now.Day(), 0, 0, 0, 0, now.Location())
	
	summary := &TodayBalanceSummary{
		OpeningBalance: &fund.FundBalance{},
		TotalInflow:    &fund.FundBalance{},
		TotalOutflow:   &fund.FundBalance{},
	}

	// Calculate opening balance (balance at start of day)
	// This would require historical tracking - for now, we'll calculate from all-time
	// TODO: Implement daily balance snapshots for accurate opening balance
	
	// Calculate today's inflow
	// 1. Fund handovers today
	fundHandovers, _, err := s.fundHandoverRepo.FindByDateRange(ctx, startOfDay, now, 1, 10000)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch today's fund handovers: %w", err)
	}
	for _, fh := range fundHandovers {
		summary.TotalInflow.Cash += fh.CashAmount
		summary.TotalInflow.Transfer += fh.TransferAmount
	}

	// 2. Deposits today
	deposits, err := s.fundTransactionRepo.FindByDateRange(
		ctx,
		startOfDay,
		now,
		ptrTransactionType(fund.TransactionTypeDeposit),
		10000,
		0,
	)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch today's deposits: %w", err)
	}
	for _, d := range deposits {
		summary.TotalInflow.Cash += d.CashAmount
		summary.TotalInflow.Transfer += d.TransferAmount
	}

	summary.TotalInflow.Total = summary.TotalInflow.Cash + summary.TotalInflow.Transfer

	// Calculate today's outflow
	// 1. Withdrawals today
	withdrawals, err := s.fundTransactionRepo.FindByDateRange(
		ctx,
		startOfDay,
		now,
		ptrTransactionType(fund.TransactionTypeWithdrawal),
		10000,
		0,
	)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch today's withdrawals: %w", err)
	}
	for _, w := range withdrawals {
		summary.TotalOutflow.Cash += w.CashAmount
		summary.TotalOutflow.Transfer += w.TransferAmount
	}

	summary.TotalOutflow.Total = summary.TotalOutflow.Cash + summary.TotalOutflow.Transfer

	return summary, nil
}

// CreateDeposit creates a new deposit transaction
func (s *FundService) CreateDeposit(
	ctx context.Context,
	cashAmount, transferAmount float64,
	reason string,
	performedBy primitive.ObjectID,
	performedByName, performedByRole string,
) (*fund.FundTransaction, *fund.FundBalance, error) {
	// Validate
	totalAmount := cashAmount + transferAmount
	if err := fund.ValidateAmount(totalAmount); err != nil {
		return nil, nil, err
	}
	if err := fund.ValidateReason(reason); err != nil {
		return nil, nil, err
	}

	// Get current balance for audit
	balanceBefore, err := s.CalculateCurrentBalance(ctx)
	if err != nil {
		return nil, nil, fmt.Errorf("failed to calculate balance: %w", err)
	}

	// Create transaction
	transaction, err := fund.NewFundTransaction(
		fund.TransactionTypeDeposit,
		cashAmount,
		transferAmount,
		reason,
		performedBy,
		performedByName,
		performedByRole,
	)
	if err != nil {
		return nil, nil, err
	}

	// Calculate balance after
	balanceAfter := &fund.FundBalance{
		Cash:     balanceBefore.Cash + cashAmount,
		Transfer: balanceBefore.Transfer + transferAmount,
		Total:    balanceBefore.Total + totalAmount,
	}
	transaction.SetBalances(balanceBefore, balanceAfter)

	// Use MongoDB transaction for atomicity
	session, err := s.mongoClient.StartSession()
	if err != nil {
		return nil, nil, fmt.Errorf("failed to start session: %w", err)
	}
	defer session.EndSession(ctx)

	err = mongo.WithSession(ctx, session, func(sc mongo.SessionContext) error {
		if err := session.StartTransaction(); err != nil {
			return err
		}

		if err := s.fundTransactionRepo.Create(sc, transaction); err != nil {
			session.AbortTransaction(sc)
			return fmt.Errorf("failed to create deposit: %w", err)
		}

		return session.CommitTransaction(sc)
	})

	if err != nil {
		return nil, nil, err
	}

	return transaction, balanceAfter, nil
}

// CreateWithdrawal creates a new withdrawal transaction
func (s *FundService) CreateWithdrawal(
	ctx context.Context,
	cashAmount, transferAmount float64,
	reason string,
	performedBy primitive.ObjectID,
	performedByName, performedByRole string,
) (*fund.FundTransaction, *fund.FundBalance, error) {
	// Validate
	totalAmount := cashAmount + transferAmount
	if err := fund.ValidateAmount(totalAmount); err != nil {
		return nil, nil, err
	}
	if err := fund.ValidateReason(reason); err != nil {
		return nil, nil, err
	}

	// Get current balance
	balanceBefore, err := s.CalculateCurrentBalance(ctx)
	if err != nil {
		return nil, nil, fmt.Errorf("failed to calculate balance: %w", err)
	}

	// Validate sufficient balance
	if cashAmount > balanceBefore.Cash {
		return nil, nil, fmt.Errorf("insufficient cash balance: have %.2f, need %.2f", balanceBefore.Cash, cashAmount)
	}
	if transferAmount > balanceBefore.Transfer {
		return nil, nil, fmt.Errorf("insufficient transfer balance: have %.2f, need %.2f", balanceBefore.Transfer, transferAmount)
	}

	// Create transaction
	transaction, err := fund.NewFundTransaction(
		fund.TransactionTypeWithdrawal,
		cashAmount,
		transferAmount,
		reason,
		performedBy,
		performedByName,
		performedByRole,
	)
	if err != nil {
		return nil, nil, err
	}

	// Calculate balance after
	balanceAfter := &fund.FundBalance{
		Cash:     balanceBefore.Cash - cashAmount,
		Transfer: balanceBefore.Transfer - transferAmount,
		Total:    balanceBefore.Total - totalAmount,
	}
	transaction.SetBalances(balanceBefore, balanceAfter)

	// Use MongoDB transaction for atomicity
	session, err := s.mongoClient.StartSession()
	if err != nil {
		return nil, nil, fmt.Errorf("failed to start session: %w", err)
	}
	defer session.EndSession(ctx)

	err = mongo.WithSession(ctx, session, func(sc mongo.SessionContext) error {
		if err := session.StartTransaction(); err != nil {
			return err
		}

		if err := s.fundTransactionRepo.Create(sc, transaction); err != nil {
			session.AbortTransaction(sc)
			return fmt.Errorf("failed to create withdrawal: %w", err)
		}

		return session.CommitTransaction(sc)
	})

	if err != nil {
		return nil, nil, err
	}

	return transaction, balanceAfter, nil
}

// GetTransactionDetail retrieves a transaction by ID
func (s *FundService) GetTransactionDetail(ctx context.Context, id primitive.ObjectID) (*fund.FundTransaction, error) {
	return s.fundTransactionRepo.FindByID(ctx, id)
}

// TodayBalanceSummary represents today's balance summary
type TodayBalanceSummary struct {
	OpeningBalance *fund.FundBalance `json:"opening_balance"`
	TotalInflow    *fund.FundBalance `json:"total_inflow"`
	TotalOutflow   *fund.FundBalance `json:"total_outflow"`
}

// Helper function to create pointer to TransactionType
func ptrTransactionType(t fund.TransactionType) *fund.TransactionType {
	return &t
}

// GetTransactionHistory retrieves transaction history with filters
func (s *FundService) GetTransactionHistory(
	ctx context.Context,
	fromDate, toDate time.Time,
	transactionType *fund.TransactionType,
	limit, offset int,
) ([]*fund.FundTransaction, error) {
	return s.fundTransactionRepo.FindByDateRange(ctx, fromDate, toDate, transactionType, limit, offset)
}

// CountTransactions counts transactions matching the filter
func (s *FundService) CountTransactions(
	ctx context.Context,
	fromDate, toDate time.Time,
	transactionType *fund.TransactionType,
) (int64, error) {
	return s.fundTransactionRepo.Count(ctx, fromDate, toDate, transactionType)
}

// TransactionHistoryItem represents a unified transaction item from any source
type TransactionHistoryItem struct {
	ID              string                 `json:"id"`
	Type            string                 `json:"type"` // "deposit", "withdrawal", "fund_handover", "cash_handover", "starting_float"
	CashAmount      float64                `json:"cash_amount"`
	TransferAmount  float64                `json:"transfer_amount"`
	TotalAmount     float64                `json:"total_amount"`
	Description     string                 `json:"description"`
	PerformedBy     string                 `json:"performed_by"`
	PerformedByRole string                 `json:"performed_by_role"`
	Timestamp       time.Time              `json:"timestamp"`
	Metadata        map[string]interface{} `json:"metadata,omitempty"`
}

// GetAggregatedTransactionHistory retrieves transaction history from all sources
func (s *FundService) GetAggregatedTransactionHistory(
	ctx context.Context,
	fromDate, toDate time.Time,
	transactionType *string, // "deposit", "withdrawal", "handover", "starting_float", "all"
	moneyType *string,       // "cash", "transfer", "all"
	limit, offset int,
) ([]*TransactionHistoryItem, int64, error) {
	var allTransactions []*TransactionHistoryItem

	// 1. Get fund transactions (deposits and withdrawals)
	if transactionType == nil || *transactionType == "all" || *transactionType == "deposit" || *transactionType == "withdrawal" {
		var txType *fund.TransactionType
		if transactionType != nil && *transactionType != "all" && *transactionType != "handover" && *transactionType != "starting_float" {
			t := fund.TransactionType(*transactionType)
			txType = &t
		}

		fundTxs, err := s.fundTransactionRepo.FindByDateRange(ctx, fromDate, toDate, txType, 10000, 0)
		if err != nil {
			return nil, 0, fmt.Errorf("failed to fetch fund transactions: %w", err)
		}

		for _, tx := range fundTxs {
			// Apply money type filter
			if moneyType != nil && *moneyType != "all" {
				if *moneyType == "cash" && tx.CashAmount == 0 {
					continue
				}
				if *moneyType == "transfer" && tx.TransferAmount == 0 {
					continue
				}
			}

			item := &TransactionHistoryItem{
				ID:              tx.ID.Hex(),
				Type:            string(tx.Type),
				CashAmount:      tx.CashAmount,
				TransferAmount:  tx.TransferAmount,
				TotalAmount:     tx.TotalAmount,
				Description:     tx.Reason,
				PerformedBy:     tx.PerformedByName,
				PerformedByRole: tx.PerformedByRole,
				Timestamp:       tx.Timestamp,
				Metadata:        map[string]interface{}{},
			}
			allTransactions = append(allTransactions, item)
		}
	}

	// 2. Get fund handovers (cashier → fund)
	if transactionType == nil || *transactionType == "all" || *transactionType == "handover" {
		fundHandovers, _, err := s.fundHandoverRepo.FindByDateRange(ctx, fromDate, toDate, 1, 10000)
		if err != nil {
			return nil, 0, fmt.Errorf("failed to fetch fund handovers: %w", err)
		}

		for _, fh := range fundHandovers {
			// Apply money type filter
			if moneyType != nil && *moneyType != "all" {
				if *moneyType == "cash" && fh.CashAmount == 0 {
					continue
				}
				if *moneyType == "transfer" && fh.TransferAmount == 0 {
					continue
				}
			}

			description := fmt.Sprintf("Bàn giao từ ca thu ngân")
			if fh.VarianceAmount != 0 {
				description += fmt.Sprintf(" (Chênh lệch: %.0f)", fh.VarianceAmount)
			}

			item := &TransactionHistoryItem{
				ID:              fh.ID.Hex(),
				Type:            "fund_handover",
				CashAmount:      fh.CashAmount,
				TransferAmount:  fh.TransferAmount,
				TotalAmount:     fh.CashAmount + fh.TransferAmount,
				Description:     description,
				PerformedBy:     fh.CashierName,
				PerformedByRole: "cashier",
				Timestamp:       fh.HandoverAt,
				Metadata: map[string]interface{}{
					"cashier_shift_id": fh.CashierShiftID.Hex(),
					"variance_amount":  fh.VarianceAmount,
					"has_variance":     fh.VarianceAmount != 0,
				},
			}
			allTransactions = append(allTransactions, item)
		}
	}

	// Sort all transactions by timestamp (newest first)
	sortTransactionsByTimestamp(allTransactions)

	// Calculate total count
	total := int64(len(allTransactions))

	// Apply pagination
	start := offset
	end := offset + limit
	if start > len(allTransactions) {
		start = len(allTransactions)
	}
	if end > len(allTransactions) {
		end = len(allTransactions)
	}

	paginatedTransactions := allTransactions[start:end]

	return paginatedTransactions, total, nil
}

// Helper function to sort transactions by timestamp
func sortTransactionsByTimestamp(transactions []*TransactionHistoryItem) {
	// Simple bubble sort for now (can be optimized with sort.Slice)
	for i := 0; i < len(transactions); i++ {
		for j := i + 1; j < len(transactions); j++ {
			if transactions[i].Timestamp.Before(transactions[j].Timestamp) {
				transactions[i], transactions[j] = transactions[j], transactions[i]
			}
		}
	}
}
