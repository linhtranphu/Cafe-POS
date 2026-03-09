package services

import (
	"context"
	"errors"
	"fmt"
	"time"

	"cafe-pos/backend/domain/fund"
	"cafe-pos/backend/infrastructure/mongodb"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	mongoDriver "go.mongodb.org/mongo-driver/mongo"
)

type FundService struct {
	fundTransactionRepo *mongodb.FundTransactionRepository
	fundHandoverRepo    *mongodb.FundHandoverRepository
	cashHandoverRepo    *mongodb.CashHandoverRepository
	cashierShiftRepo    *mongodb.CashierShiftRepository
	mongoClient         *mongoDriver.Client
}

func NewFundService(
	fundTransactionRepo *mongodb.FundTransactionRepository,
	fundHandoverRepo *mongodb.FundHandoverRepository,
	cashHandoverRepo *mongodb.CashHandoverRepository,
	cashierShiftRepo *mongodb.CashierShiftRepository,
	mongoClient *mongoDriver.Client,
) *FundService {
	return &FundService{
		fundTransactionRepo: fundTransactionRepo,
		fundHandoverRepo:    fundHandoverRepo,
		cashHandoverRepo:    cashHandoverRepo,
		cashierShiftRepo:    cashierShiftRepo,
		mongoClient:         mongoClient,
	}
}

// ─── Balance Types ────────────────────────────────────────────────────────────

// TodayBalanceSummary represents today's balance summary
type TodayBalanceSummary struct {
	OpeningBalance *fund.FundBalance `json:"opening_balance"`
	TotalInflow    *fund.FundBalance `json:"total_inflow"`
	TotalOutflow   *fund.FundBalance `json:"total_outflow"`
}

// AllFundBalances holds the balance of all 4 fund types + total
type AllFundBalances struct {
	Operating  fund.FundBalance `json:"operating"`
	Inventory  fund.FundBalance `json:"inventory"`
	Profit     fund.FundBalance `json:"profit"`
	CashDrawer fund.FundBalance `json:"cash_drawer"`
	Total      fund.FundBalance `json:"total"`
}

// ─── Transfer Types ───────────────────────────────────────────────────────────

// FundTransferRequest represents a request to transfer between funds
type FundTransferRequest struct {
	FromFundType    fund.FundType
	ToFundType      fund.FundType
	CashAmount      float64
	TransferAmount  float64
	Reason          string
	PerformedBy     primitive.ObjectID
	PerformedByName string
	PerformedByRole string
}

// FundTransferResult contains the two transactions created by a transfer
type FundTransferResult struct {
	WithdrawalTx *fund.FundTransaction `json:"withdrawal_tx"`
	DepositTx    *fund.FundTransaction `json:"deposit_tx"`
	FromBalance  *fund.FundBalance     `json:"from_balance"`
	ToBalance    *fund.FundBalance     `json:"to_balance"`
}

// ─── Balance Calculation ──────────────────────────────────────────────────────

// CalculateCurrentBalance calculates the total fund balance from all sources (backward compat)
func (s *FundService) CalculateCurrentBalance(ctx context.Context) (*fund.FundBalance, error) {
	balance := &fund.FundBalance{}

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
		ctx, time.Time{}, time.Now(),
		ptrTransactionType(fund.TransactionTypeDeposit), 10000, 0,
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
		ctx, time.Time{}, time.Now(),
		ptrTransactionType(fund.TransactionTypeWithdrawal), 10000, 0,
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

// CalculateBalanceByFundType calculates balance for a specific fund type
// Backward compat: "operating" includes transactions with empty/missing fund_type
// and all fund_handovers from the fund handover repo
func (s *FundService) CalculateBalanceByFundType(ctx context.Context, fundType fund.FundType) (*fund.FundBalance, error) {
	balance := &fund.FundBalance{}

	// Build filter — operating includes legacy transactions (no fund_type set)
	var filter bson.M
	if fundType == fund.FundTypeOperating {
		filter = bson.M{
			"$or": bson.A{
				bson.M{"fund_type": fund.FundTypeOperating},
				bson.M{"fund_type": ""},
				bson.M{"fund_type": bson.M{"$exists": false}},
			},
		}
	} else {
		filter = bson.M{"fund_type": fundType}
	}

	txs, err := s.fundTransactionRepo.FindByFilter(ctx, filter, 10000, 0)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch transactions for %s: %w", fundType, err)
	}

	for _, tx := range txs {
		switch tx.Type {
		case fund.TransactionTypeDeposit, fund.TransactionTypeFundHandover:
			balance.Cash += tx.CashAmount
			balance.Transfer += tx.TransferAmount
		case fund.TransactionTypeWithdrawal, fund.TransactionTypeStartingFloat:
			balance.Cash -= tx.CashAmount
			balance.Transfer -= tx.TransferAmount
		}
	}

	// For operating: also include historical fund_handovers from FundHandoverRepo
	if fundType == fund.FundTypeOperating {
		fundHandovers, _, err := s.fundHandoverRepo.FindByDateRange(ctx, time.Time{}, time.Now(), 1, 10000)
		if err != nil {
			return nil, fmt.Errorf("failed to fetch fund handovers: %w", err)
		}
		for _, fh := range fundHandovers {
			balance.Cash += fh.CashAmount
			balance.Transfer += fh.TransferAmount
		}
	}

	balance.Total = balance.Cash + balance.Transfer
	return balance, nil
}

// GetAllFundBalances returns balances for all 4 fund types plus total
func (s *FundService) GetAllFundBalances(ctx context.Context) (*AllFundBalances, error) {
	result := &AllFundBalances{}

	fundTypes := []fund.FundType{
		fund.FundTypeOperating,
		fund.FundTypeInventory,
		fund.FundTypeProfit,
		fund.FundTypeCashDrawer,
	}

	for _, ft := range fundTypes {
		bal, err := s.CalculateBalanceByFundType(ctx, ft)
		if err != nil {
			return nil, fmt.Errorf("balance calculation failed for %s: %w", ft, err)
		}
		switch ft {
		case fund.FundTypeOperating:
			result.Operating = *bal
		case fund.FundTypeInventory:
			result.Inventory = *bal
		case fund.FundTypeProfit:
			result.Profit = *bal
		case fund.FundTypeCashDrawer:
			result.CashDrawer = *bal
		}
		result.Total.Cash += bal.Cash
		result.Total.Transfer += bal.Transfer
		result.Total.Total += bal.Total
	}

	return result, nil
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

	// Fund handovers today
	fundHandovers, _, err := s.fundHandoverRepo.FindByDateRange(ctx, startOfDay, now, 1, 10000)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch today's fund handovers: %w", err)
	}
	for _, fh := range fundHandovers {
		summary.TotalInflow.Cash += fh.CashAmount
		summary.TotalInflow.Transfer += fh.TransferAmount
	}

	// Deposits today
	deposits, err := s.fundTransactionRepo.FindByDateRange(
		ctx, startOfDay, now,
		ptrTransactionType(fund.TransactionTypeDeposit), 10000, 0,
	)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch today's deposits: %w", err)
	}
	for _, d := range deposits {
		summary.TotalInflow.Cash += d.CashAmount
		summary.TotalInflow.Transfer += d.TransferAmount
	}
	summary.TotalInflow.Total = summary.TotalInflow.Cash + summary.TotalInflow.Transfer

	// Withdrawals today
	withdrawals, err := s.fundTransactionRepo.FindByDateRange(
		ctx, startOfDay, now,
		ptrTransactionType(fund.TransactionTypeWithdrawal), 10000, 0,
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

// ─── Deposit / Withdrawal ─────────────────────────────────────────────────────

// CreateDeposit creates a new deposit transaction
// Optional fundType defaults to "operating" for backward compat
func (s *FundService) CreateDeposit(
	ctx context.Context,
	cashAmount, transferAmount float64,
	reason string,
	performedBy primitive.ObjectID,
	performedByName, performedByRole string,
	fundType ...fund.FundType,
) (*fund.FundTransaction, *fund.FundBalance, error) {
	ft := fund.FundTypeOperating
	if len(fundType) > 0 && fundType[0] != "" {
		ft = fundType[0]
	}

	totalAmount := cashAmount + transferAmount
	if err := fund.ValidateAmount(totalAmount); err != nil {
		return nil, nil, err
	}
	if err := fund.ValidateReason(reason); err != nil {
		return nil, nil, err
	}

	balanceBefore, err := s.CalculateBalanceByFundType(ctx, ft)
	if err != nil {
		return nil, nil, fmt.Errorf("failed to calculate balance: %w", err)
	}

	transaction, err := fund.NewFundTransaction(
		fund.TransactionTypeDeposit,
		cashAmount, transferAmount,
		reason,
		performedBy, performedByName, performedByRole,
	)
	if err != nil {
		return nil, nil, err
	}
	transaction.SetFundType(ft)

	balanceAfter := &fund.FundBalance{
		Cash:     balanceBefore.Cash + cashAmount,
		Transfer: balanceBefore.Transfer + transferAmount,
		Total:    balanceBefore.Total + totalAmount,
	}
	transaction.SetBalances(balanceBefore, balanceAfter)

	session, err := s.mongoClient.StartSession()
	if err != nil {
		return nil, nil, fmt.Errorf("failed to start session: %w", err)
	}
	defer session.EndSession(ctx)

	err = mongoDriver.WithSession(ctx, session, func(sc mongoDriver.SessionContext) error {
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
// Optional fundType defaults to "operating" for backward compat
func (s *FundService) CreateWithdrawal(
	ctx context.Context,
	cashAmount, transferAmount float64,
	reason string,
	performedBy primitive.ObjectID,
	performedByName, performedByRole string,
	fundType ...fund.FundType,
) (*fund.FundTransaction, *fund.FundBalance, error) {
	ft := fund.FundTypeOperating
	if len(fundType) > 0 && fundType[0] != "" {
		ft = fundType[0]
	}

	totalAmount := cashAmount + transferAmount
	if err := fund.ValidateAmount(totalAmount); err != nil {
		return nil, nil, err
	}
	if err := fund.ValidateReason(reason); err != nil {
		return nil, nil, err
	}

	balanceBefore, err := s.CalculateBalanceByFundType(ctx, ft)
	if err != nil {
		return nil, nil, fmt.Errorf("failed to calculate balance: %w", err)
	}

	if cashAmount > balanceBefore.Cash {
		return nil, nil, fmt.Errorf("insufficient cash balance in %s: have %.0f, need %.0f", ft, balanceBefore.Cash, cashAmount)
	}
	if transferAmount > balanceBefore.Transfer {
		return nil, nil, fmt.Errorf("insufficient transfer balance in %s: have %.0f, need %.0f", ft, balanceBefore.Transfer, transferAmount)
	}

	transaction, err := fund.NewFundTransaction(
		fund.TransactionTypeWithdrawal,
		cashAmount, transferAmount,
		reason,
		performedBy, performedByName, performedByRole,
	)
	if err != nil {
		return nil, nil, err
	}
	transaction.SetFundType(ft)

	balanceAfter := &fund.FundBalance{
		Cash:     balanceBefore.Cash - cashAmount,
		Transfer: balanceBefore.Transfer - transferAmount,
		Total:    balanceBefore.Total - totalAmount,
	}
	transaction.SetBalances(balanceBefore, balanceAfter)

	session, err := s.mongoClient.StartSession()
	if err != nil {
		return nil, nil, fmt.Errorf("failed to start session: %w", err)
	}
	defer session.EndSession(ctx)

	err = mongoDriver.WithSession(ctx, session, func(sc mongoDriver.SessionContext) error {
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

// ─── Inter-Fund Transfer ──────────────────────────────────────────────────────

// TransferBetweenFunds atomically moves money from one fund to another
func (s *FundService) TransferBetweenFunds(ctx context.Context, req FundTransferRequest) (*FundTransferResult, error) {
	if req.FromFundType == req.ToFundType {
		return nil, errors.New("source and destination funds cannot be the same")
	}

	totalAmount := req.CashAmount + req.TransferAmount
	if totalAmount <= 0 {
		return nil, errors.New("transfer amount must be greater than 0")
	}
	if len(req.Reason) < 10 {
		return nil, errors.New("reason must be at least 10 characters")
	}

	// Validate source balance
	fromBalance, err := s.CalculateBalanceByFundType(ctx, req.FromFundType)
	if err != nil {
		return nil, fmt.Errorf("failed to get source balance: %w", err)
	}
	if req.CashAmount > fromBalance.Cash {
		return nil, fmt.Errorf("insufficient cash in source fund: have %.0f, need %.0f", fromBalance.Cash, req.CashAmount)
	}
	if req.TransferAmount > fromBalance.Transfer {
		return nil, fmt.Errorf("insufficient transfer in source fund: have %.0f, need %.0f", fromBalance.Transfer, req.TransferAmount)
	}

	toBalance, err := s.CalculateBalanceByFundType(ctx, req.ToFundType)
	if err != nil {
		return nil, fmt.Errorf("failed to get destination balance: %w", err)
	}

	// Pre-generate IDs so they can cross-reference each other
	withdrawalID := primitive.NewObjectID()
	depositID := primitive.NewObjectID()

	// Withdrawal side (from source fund)
	withdrawalTx, err := fund.NewFundTransaction(
		fund.TransactionTypeWithdrawal,
		req.CashAmount, req.TransferAmount,
		req.Reason,
		req.PerformedBy, req.PerformedByName, req.PerformedByRole,
	)
	if err != nil {
		return nil, err
	}
	withdrawalTx.ID = withdrawalID
	withdrawalTx.SetFundType(req.FromFundType)
	withdrawalTx.SetSource(fund.SourceTypeFundTransfer, depositID) //nolint
	withdrawalTx.LinkTo(depositID)

	fromBalanceAfter := &fund.FundBalance{
		Cash:     fromBalance.Cash - req.CashAmount,
		Transfer: fromBalance.Transfer - req.TransferAmount,
		Total:    fromBalance.Total - totalAmount,
	}
	withdrawalTx.SetBalances(fromBalance, fromBalanceAfter)

	// Deposit side (into destination fund)
	depositTx, err := fund.NewFundTransaction(
		fund.TransactionTypeDeposit,
		req.CashAmount, req.TransferAmount,
		req.Reason,
		req.PerformedBy, req.PerformedByName, req.PerformedByRole,
	)
	if err != nil {
		return nil, err
	}
	depositTx.ID = depositID
	depositTx.SetFundType(req.ToFundType)
	depositTx.SetSource(fund.SourceTypeFundTransfer, withdrawalID) //nolint
	depositTx.LinkTo(withdrawalID)

	toBalanceAfter := &fund.FundBalance{
		Cash:     toBalance.Cash + req.CashAmount,
		Transfer: toBalance.Transfer + req.TransferAmount,
		Total:    toBalance.Total + totalAmount,
	}
	depositTx.SetBalances(toBalance, toBalanceAfter)

	// Atomic insert of both transactions
	session, err := s.mongoClient.StartSession()
	if err != nil {
		return nil, fmt.Errorf("failed to start session: %w", err)
	}
	defer session.EndSession(ctx)

	err = mongoDriver.WithSession(ctx, session, func(sc mongoDriver.SessionContext) error {
		if err := session.StartTransaction(); err != nil {
			return err
		}
		if err := s.fundTransactionRepo.Create(sc, withdrawalTx); err != nil {
			session.AbortTransaction(sc)
			return fmt.Errorf("failed to create withdrawal: %w", err)
		}
		if err := s.fundTransactionRepo.Create(sc, depositTx); err != nil {
			session.AbortTransaction(sc)
			return fmt.Errorf("failed to create deposit: %w", err)
		}
		return session.CommitTransaction(sc)
	})
	if err != nil {
		return nil, err
	}

	return &FundTransferResult{
		WithdrawalTx: withdrawalTx,
		DepositTx:    depositTx,
		FromBalance:  fromBalanceAfter,
		ToBalance:    toBalanceAfter,
	}, nil
}

// ─── Transaction History ──────────────────────────────────────────────────────

// GetTransactionDetail retrieves a transaction by ID
// RecordHandover creates a FundTransaction for a cashier shift handover event.
// Called automatically when CloseShiftWithFundHandover succeeds.
func (s *FundService) RecordHandover(
	ctx context.Context,
	cashAmount, transferAmount float64,
	shiftID primitive.ObjectID,
	performedBy primitive.ObjectID,
	performedByName, performedByRole string,
) error {
	reason := fmt.Sprintf("Bàn giao quỹ cuối ca (shift %s)", shiftID.Hex())

	transaction, err := fund.NewFundTransaction(
		fund.TransactionTypeFundHandover,
		cashAmount, transferAmount,
		reason,
		performedBy, performedByName, performedByRole,
	)
	if err != nil {
		return fmt.Errorf("failed to build handover transaction: %w", err)
	}
	transaction.SetFundType(fund.FundTypeCashDrawer)
	if err := transaction.SetSource(fund.SourceTypeHandover, shiftID); err != nil {
		return err
	}

	balanceBefore, err := s.CalculateBalanceByFundType(ctx, fund.FundTypeCashDrawer)
	if err != nil {
		return fmt.Errorf("failed to calculate balance before handover: %w", err)
	}
	balanceAfter := &fund.FundBalance{
		Cash:     balanceBefore.Cash + cashAmount,
		Transfer: balanceBefore.Transfer + transferAmount,
		Total:    balanceBefore.Total + cashAmount + transferAmount,
	}
	transaction.SetBalances(balanceBefore, balanceAfter)

	return s.fundTransactionRepo.Create(ctx, transaction)
}

func (s *FundService) GetTransactionDetail(ctx context.Context, id primitive.ObjectID) (*fund.FundTransaction, error) {
	return s.fundTransactionRepo.FindByID(ctx, id)
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
	Type            string                 `json:"type"`
	FundType        string                 `json:"fund_type,omitempty"`
	SourceType      string                 `json:"source_type,omitempty"`
	SourceID        string                 `json:"source_id,omitempty"`
	LinkedTxID      string                 `json:"linked_tx_id,omitempty"`
	CashAmount      float64                `json:"cash_amount"`
	TransferAmount  float64                `json:"transfer_amount"`
	TotalAmount     float64                `json:"total_amount"`
	Description     string                 `json:"description"`
	PerformedBy     string                 `json:"performed_by"`
	PerformedByRole string                 `json:"performed_by_role"`
	Timestamp       time.Time              `json:"timestamp"`
	BalanceBefore   *fund.FundBalance      `json:"balance_before,omitempty"`
	BalanceAfter    *fund.FundBalance      `json:"balance_after,omitempty"`
	Metadata        map[string]interface{} `json:"metadata,omitempty"`
}

// GetAggregatedTransactionHistory retrieves transaction history from all sources with filters
func (s *FundService) GetAggregatedTransactionHistory(
	ctx context.Context,
	fromDate, toDate time.Time,
	transactionType *string, // "deposit","withdrawal","handover","fund_transfer","all"
	moneyType *string,       // "cash","transfer","all"
	fundTypeFilter *string,  // "operating","inventory","profit","cash_drawer","all"
	sourceTypeFilter *string, // "expense","ingredient","handover","fund_transfer","manual","all"
	limit, offset int,
) ([]*TransactionHistoryItem, int64, error) {
	var allItems []*TransactionHistoryItem

	// Determine if we should include handovers and fund_transactions
	includeHandovers := transactionType == nil || *transactionType == "all" || *transactionType == "handover"
	includeFundTxs := transactionType == nil || *transactionType == "all" ||
		*transactionType == "deposit" || *transactionType == "withdrawal" ||
		*transactionType == "fund_transfer"

	// When sourceType filter is "handover", only show handovers
	if sourceTypeFilter != nil && *sourceTypeFilter == "handover" {
		includeFundTxs = false
		includeHandovers = true
	}
	// When sourceType filter is specific (not handover/all), only show fund_txs
	if sourceTypeFilter != nil && *sourceTypeFilter != "all" && *sourceTypeFilter != "handover" {
		includeHandovers = false
		includeFundTxs = true
	}

	// 1. Fund transactions
	if includeFundTxs {
		filter := bson.M{
			"timestamp": bson.M{"$gte": fromDate, "$lte": toDate},
		}

		// Transaction type filter
		if transactionType != nil && *transactionType != "all" && *transactionType != "handover" {
			if *transactionType == "fund_transfer" {
				filter["source_type"] = fund.SourceTypeFundTransfer
			} else {
				filter["type"] = *transactionType
			}
		}

		// Fund type filter
		if fundTypeFilter != nil && *fundTypeFilter != "all" && *fundTypeFilter != "" {
			if *fundTypeFilter == "operating" {
				filter["$or"] = bson.A{
					bson.M{"fund_type": fund.FundTypeOperating},
					bson.M{"fund_type": ""},
					bson.M{"fund_type": bson.M{"$exists": false}},
				}
			} else {
				filter["fund_type"] = *fundTypeFilter
			}
		}

		// Source type filter
		if sourceTypeFilter != nil && *sourceTypeFilter != "all" && *sourceTypeFilter != "handover" {
			if *sourceTypeFilter == "manual" {
				// Manual = no source_type set
				filter["source_type"] = bson.M{"$in": bson.A{"", nil}}
			} else {
				filter["source_type"] = *sourceTypeFilter
			}
		}

		txs, err := s.fundTransactionRepo.FindByFilter(ctx, filter, 10000, 0)
		if err != nil {
			return nil, 0, fmt.Errorf("failed to fetch fund transactions: %w", err)
		}

		for _, tx := range txs {
			if !matchesMoneyType(tx.CashAmount, tx.TransferAmount, moneyType) {
				continue
			}
			desc := tx.Description
			if desc == "" {
				desc = tx.Reason
			}
			item := &TransactionHistoryItem{
				ID:              tx.ID.Hex(),
				Type:            string(tx.Type),
				FundType:        string(tx.FundType),
				SourceType:      string(tx.SourceType),
				SourceID:        tx.SourceID.Hex(),
				LinkedTxID:      tx.LinkedTxID.Hex(),
				CashAmount:      tx.CashAmount,
				TransferAmount:  tx.TransferAmount,
				TotalAmount:     tx.TotalAmount,
				Description:     desc,
				PerformedBy:     tx.PerformedByName,
				PerformedByRole: tx.PerformedByRole,
				Timestamp:       tx.Timestamp,
				BalanceBefore:   tx.BalanceBefore,
				BalanceAfter:    tx.BalanceAfter,
				Metadata:        map[string]interface{}{},
			}
			// Clean zero-value ObjectID hex strings
			if tx.SourceID.IsZero() {
				item.SourceID = ""
			}
			if tx.LinkedTxID.IsZero() {
				item.LinkedTxID = ""
			}
			allItems = append(allItems, item)
		}
	}

	// 2. Fund handovers from FundHandoverRepo (only for operating / all funds)
	showHandoversForFund := fundTypeFilter == nil || *fundTypeFilter == "all" ||
		*fundTypeFilter == "operating" || *fundTypeFilter == ""
	if includeHandovers && showHandoversForFund {
		fundHandovers, _, err := s.fundHandoverRepo.FindByDateRange(ctx, fromDate, toDate, 1, 10000)
		if err != nil {
			return nil, 0, fmt.Errorf("failed to fetch fund handovers: %w", err)
		}

		for _, fh := range fundHandovers {
			if !matchesMoneyType(fh.CashAmount, fh.TransferAmount, moneyType) {
				continue
			}
			desc := "Bàn giao từ ca thu ngân"
			if fh.VarianceAmount != 0 {
				desc += fmt.Sprintf(" (Chênh lệch: %.0f)", fh.VarianceAmount)
			}
			item := &TransactionHistoryItem{
				ID:              fh.ID.Hex(),
				Type:            "fund_handover",
				FundType:        "operating",
				SourceType:      "handover",
				SourceID:        fh.CashierShiftID.Hex(),
				CashAmount:      fh.CashAmount,
				TransferAmount:  fh.TransferAmount,
				TotalAmount:     fh.CashAmount + fh.TransferAmount,
				Description:     desc,
				PerformedBy:     fh.CashierName,
				PerformedByRole: "cashier",
				Timestamp:       fh.HandoverAt,
				Metadata: map[string]interface{}{
					"cashier_shift_id": fh.CashierShiftID.Hex(),
					"variance_amount":  fh.VarianceAmount,
					"has_variance":     fh.VarianceAmount != 0,
				},
			}
			allItems = append(allItems, item)
		}
	}

	// Sort by timestamp descending
	sortTransactionsByTimestamp(allItems)

	total := int64(len(allItems))

	// Paginate
	start := offset
	if start > len(allItems) {
		start = len(allItems)
	}
	end := start + limit
	if end > len(allItems) {
		end = len(allItems)
	}

	return allItems[start:end], total, nil
}

// ─── Helpers ─────────────────────────────────────────────────────────────────

func ptrTransactionType(t fund.TransactionType) *fund.TransactionType {
	return &t
}

func matchesMoneyType(cashAmount, transferAmount float64, moneyType *string) bool {
	if moneyType == nil || *moneyType == "all" || *moneyType == "" {
		return true
	}
	if *moneyType == "cash" && cashAmount == 0 {
		return false
	}
	if *moneyType == "transfer" && transferAmount == 0 {
		return false
	}
	return true
}

func sortTransactionsByTimestamp(transactions []*TransactionHistoryItem) {
	for i := 0; i < len(transactions); i++ {
		for j := i + 1; j < len(transactions); j++ {
			if transactions[i].Timestamp.Before(transactions[j].Timestamp) {
				transactions[i], transactions[j] = transactions[j], transactions[i]
			}
		}
	}
}
