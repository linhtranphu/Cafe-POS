package services

import (
	"context"
	"fmt"

	"cafe-pos/backend/domain/fund"
	"cafe-pos/backend/infrastructure/mongodb"

	"go.mongodb.org/mongo-driver/bson/primitive"
	mongoDriver "go.mongodb.org/mongo-driver/mongo"
)

// JournalService handles all double-entry bookkeeping for internal fund transfers.
// Every method creates a balanced JournalEntry (Σ DEBIT == Σ CREDIT).
// External payments (expenses) stay in fund_transactions (single-entry).
type JournalService struct {
	repo        *mongodb.JournalRepository
	mongoClient *mongoDriver.Client
}

func NewJournalService(repo *mongodb.JournalRepository, mongoClient *mongoDriver.Client) *JournalService {
	return &JournalService{repo: repo, mongoClient: mongoClient}
}

// ─── Balance Queries ──────────────────────────────────────────────────────────

func (s *JournalService) GetFundBalance(ctx context.Context, fundType fund.FundType) (*fund.FundBalance, error) {
	return s.repo.GetFundBalance(ctx, fundType)
}

func (s *JournalService) GetAllFundBalances(ctx context.Context) (map[fund.FundType]*fund.FundBalance, error) {
	return s.repo.GetAllFundBalances(ctx)
}

func (s *JournalService) ListEntries(ctx context.Context, filter mongodb.JournalFilter) ([]*fund.JournalEntry, int64, error) {
	return s.repo.List(ctx, filter)
}

func (s *JournalService) GetEntry(ctx context.Context, id primitive.ObjectID) (*fund.JournalEntry, error) {
	return s.repo.FindByID(ctx, id)
}

// ─── Business Events ──────────────────────────────────────────────────────────

// RecordCashierShiftStart creates:
//   DEBIT  cash_drawer  X  (ngăn kéo nhận tiền)
//   CREDIT operating    X  (quỹ vận hành xuất tiền)
func (s *JournalService) RecordCashierShiftStart(
	ctx context.Context,
	startingFloat float64,
	cashierShiftID primitive.ObjectID,
	cashierID primitive.ObjectID,
	cashierName string,
) (*fund.JournalEntry, error) {
	if startingFloat <= 0 {
		return nil, fmt.Errorf("starting float must be greater than 0")
	}

	session, err := s.mongoClient.StartSession()
	if err != nil {
		return nil, fmt.Errorf("failed to start session: %w", err)
	}
	defer session.EndSession(ctx)

	var entry *fund.JournalEntry
	err = mongoDriver.WithSession(ctx, session, func(sc mongoDriver.SessionContext) error {
		if err := session.StartTransaction(); err != nil {
			return err
		}
		e, err := s.recordCashierShiftStart(sc, startingFloat, cashierShiftID, cashierID, cashierName)
		if err != nil {
			session.AbortTransaction(sc)
			return err
		}
		entry = e
		return session.CommitTransaction(sc)
	})
	return entry, err
}

// recordCashierShiftStart is the session-aware inner implementation
func (s *JournalService) recordCashierShiftStart(
	ctx context.Context,
	startingFloat float64,
	cashierShiftID primitive.ObjectID,
	cashierID primitive.ObjectID,
	cashierName string,
) (*fund.JournalEntry, error) {
	// Get current balances for both funds
	operatingBefore, err := s.repo.GetFundBalance(ctx, fund.FundTypeOperating)
	if err != nil {
		return nil, fmt.Errorf("failed to get operating balance: %w", err)
	}
	drawerBefore, err := s.repo.GetFundBalance(ctx, fund.FundTypeCashDrawer)
	if err != nil {
		return nil, fmt.Errorf("failed to get cash_drawer balance: %w", err)
	}

	desc := fmt.Sprintf("Tiền đầu ca thu ngân %s - %.0fđ", cashierName, startingFloat)
	entry := fund.NewJournalEntry(
		fund.EventCashierShiftStart,
		cashierShiftID,
		desc,
		cashierID,
		cashierName,
		"cashier",
	)

	// DEBIT cash_drawer (receives money)
	drawerAfter := fund.FundBalance{
		Cash:     drawerBefore.Cash + startingFloat,
		Transfer: drawerBefore.Transfer,
		Total:    drawerBefore.Total + startingFloat,
	}
	entry.AddLine(fund.FundTypeCashDrawer, fund.DirectionDebit, startingFloat, 0, *drawerBefore, drawerAfter)

	// CREDIT operating (sends money)
	operatingAfter := fund.FundBalance{
		Cash:     operatingBefore.Cash - startingFloat,
		Transfer: operatingBefore.Transfer,
		Total:    operatingBefore.Total - startingFloat,
	}
	entry.AddLine(fund.FundTypeOperating, fund.DirectionCredit, startingFloat, 0, *operatingBefore, operatingAfter)

	if err := entry.Validate(); err != nil {
		return nil, err
	}

	if err := s.repo.Create(ctx, entry); err != nil {
		return nil, fmt.Errorf("failed to persist journal entry: %w", err)
	}
	return entry, nil
}

// RecordCashierShiftEnd creates:
//   DEBIT  operating    X  (quỹ vận hành nhận tiền về)
//   CREDIT cash_drawer  X  (ngăn kéo xuất tiền)
func (s *JournalService) RecordCashierShiftEnd(
	ctx context.Context,
	cashAmount, transferAmount float64,
	cashierShiftID primitive.ObjectID,
	cashierID primitive.ObjectID,
	cashierName string,
) (*fund.JournalEntry, error) {
	total := cashAmount + transferAmount
	if total <= 0 {
		return nil, fmt.Errorf("handover amount must be greater than 0")
	}

	session, err := s.mongoClient.StartSession()
	if err != nil {
		return nil, fmt.Errorf("failed to start session: %w", err)
	}
	defer session.EndSession(ctx)

	var entry *fund.JournalEntry
	err = mongoDriver.WithSession(ctx, session, func(sc mongoDriver.SessionContext) error {
		if err := session.StartTransaction(); err != nil {
			return err
		}
		e, err := s.recordCashierShiftEnd(sc, cashAmount, transferAmount, cashierShiftID, cashierID, cashierName)
		if err != nil {
			session.AbortTransaction(sc)
			return err
		}
		entry = e
		return session.CommitTransaction(sc)
	})
	return entry, err
}

func (s *JournalService) recordCashierShiftEnd(
	ctx context.Context,
	cashAmount, transferAmount float64,
	cashierShiftID primitive.ObjectID,
	cashierID primitive.ObjectID,
	cashierName string,
) (*fund.JournalEntry, error) {
	total := cashAmount + transferAmount

	operatingBefore, err := s.repo.GetFundBalance(ctx, fund.FundTypeOperating)
	if err != nil {
		return nil, fmt.Errorf("failed to get operating balance: %w", err)
	}
	drawerBefore, err := s.repo.GetFundBalance(ctx, fund.FundTypeCashDrawer)
	if err != nil {
		return nil, fmt.Errorf("failed to get cash_drawer balance: %w", err)
	}

	desc := fmt.Sprintf("Bàn giao quỹ cuối ca thu ngân %s - tiền mặt: %.0fđ, chuyển khoản: %.0fđ", cashierName, cashAmount, transferAmount)
	entry := fund.NewJournalEntry(
		fund.EventCashierShiftEnd,
		cashierShiftID,
		desc,
		cashierID,
		cashierName,
		"cashier",
	)

	// DEBIT operating (receives money back)
	operatingAfter := fund.FundBalance{
		Cash:     operatingBefore.Cash + cashAmount,
		Transfer: operatingBefore.Transfer + transferAmount,
		Total:    operatingBefore.Total + total,
	}
	entry.AddLine(fund.FundTypeOperating, fund.DirectionDebit, cashAmount, transferAmount, *operatingBefore, operatingAfter)

	// CREDIT cash_drawer (sends money out)
	drawerAfter := fund.FundBalance{
		Cash:     drawerBefore.Cash - cashAmount,
		Transfer: drawerBefore.Transfer - transferAmount,
		Total:    drawerBefore.Total - total,
	}
	entry.AddLine(fund.FundTypeCashDrawer, fund.DirectionCredit, cashAmount, transferAmount, *drawerBefore, drawerAfter)

	if err := entry.Validate(); err != nil {
		return nil, err
	}

	if err := s.repo.Create(ctx, entry); err != nil {
		return nil, fmt.Errorf("failed to persist journal entry: %w", err)
	}
	return entry, nil
}

// RecordWaiterShiftStart creates:
//   DEBIT  waiter_float X  (waiter nhận tiền)
//   CREDIT cash_drawer  X  (ngăn kéo xuất tiền)
func (s *JournalService) RecordWaiterShiftStart(
	ctx context.Context,
	startCash float64,
	waiterShiftID primitive.ObjectID,
	waiterID primitive.ObjectID,
	waiterName string,
) (*fund.JournalEntry, error) {
	if startCash <= 0 {
		return nil, fmt.Errorf("start cash must be greater than 0")
	}

	session, err := s.mongoClient.StartSession()
	if err != nil {
		return nil, fmt.Errorf("failed to start session: %w", err)
	}
	defer session.EndSession(ctx)

	var entry *fund.JournalEntry
	err = mongoDriver.WithSession(ctx, session, func(sc mongoDriver.SessionContext) error {
		if err := session.StartTransaction(); err != nil {
			return err
		}
		e, err := s.recordWaiterShiftStart(sc, startCash, waiterShiftID, waiterID, waiterName)
		if err != nil {
			session.AbortTransaction(sc)
			return err
		}
		entry = e
		return session.CommitTransaction(sc)
	})
	return entry, err
}

func (s *JournalService) recordWaiterShiftStart(
	ctx context.Context,
	startCash float64,
	waiterShiftID primitive.ObjectID,
	waiterID primitive.ObjectID,
	waiterName string,
) (*fund.JournalEntry, error) {
	drawerBefore, err := s.repo.GetFundBalance(ctx, fund.FundTypeCashDrawer)
	if err != nil {
		return nil, fmt.Errorf("failed to get cash_drawer balance: %w", err)
	}
	if drawerBefore.Cash < startCash {
		return nil, fmt.Errorf("ngăn kéo không đủ tiền mặt: có %.0f, cần %.0f", drawerBefore.Cash, startCash)
	}

	waiterBefore, err := s.repo.GetFundBalance(ctx, fund.FundTypeWaiterFloat)
	if err != nil {
		return nil, fmt.Errorf("failed to get waiter_float balance: %w", err)
	}

	desc := fmt.Sprintf("Tiền đầu ca phục vụ %s - %.0fđ", waiterName, startCash)
	entry := fund.NewJournalEntry(
		fund.EventWaiterShiftStart,
		waiterShiftID,
		desc,
		waiterID,
		waiterName,
		"waiter",
	)

	// DEBIT waiter_float (waiter receives money)
	waiterAfter := fund.FundBalance{
		Cash:     waiterBefore.Cash + startCash,
		Transfer: waiterBefore.Transfer,
		Total:    waiterBefore.Total + startCash,
	}
	entry.AddLine(fund.FundTypeWaiterFloat, fund.DirectionDebit, startCash, 0, *waiterBefore, waiterAfter)

	// CREDIT cash_drawer (drawer gives money)
	drawerAfter := fund.FundBalance{
		Cash:     drawerBefore.Cash - startCash,
		Transfer: drawerBefore.Transfer,
		Total:    drawerBefore.Total - startCash,
	}
	entry.AddLine(fund.FundTypeCashDrawer, fund.DirectionCredit, startCash, 0, *drawerBefore, drawerAfter)

	if err := entry.Validate(); err != nil {
		return nil, err
	}

	if err := s.repo.Create(ctx, entry); err != nil {
		return nil, fmt.Errorf("failed to persist journal entry: %w", err)
	}
	return entry, nil
}

// RecordWaiterHandover creates:
//   DEBIT  cash_drawer  X  (ngăn kéo nhận tiền về)
//   CREDIT waiter_float X  (waiter hoàn trả tiền)
func (s *JournalService) RecordWaiterHandover(
	ctx context.Context,
	cashAmount, transferAmount float64,
	handoverID primitive.ObjectID,
	waiterID primitive.ObjectID,
	waiterName string,
) (*fund.JournalEntry, error) {
	total := cashAmount + transferAmount
	if total <= 0 {
		return nil, fmt.Errorf("handover amount must be greater than 0")
	}

	session, err := s.mongoClient.StartSession()
	if err != nil {
		return nil, fmt.Errorf("failed to start session: %w", err)
	}
	defer session.EndSession(ctx)

	var entry *fund.JournalEntry
	err = mongoDriver.WithSession(ctx, session, func(sc mongoDriver.SessionContext) error {
		if err := session.StartTransaction(); err != nil {
			return err
		}
		e, err := s.recordWaiterHandover(sc, cashAmount, transferAmount, handoverID, waiterID, waiterName)
		if err != nil {
			session.AbortTransaction(sc)
			return err
		}
		entry = e
		return session.CommitTransaction(sc)
	})
	return entry, err
}

func (s *JournalService) recordWaiterHandover(
	ctx context.Context,
	cashAmount, transferAmount float64,
	handoverID primitive.ObjectID,
	waiterID primitive.ObjectID,
	waiterName string,
) (*fund.JournalEntry, error) {
	total := cashAmount + transferAmount

	drawerBefore, err := s.repo.GetFundBalance(ctx, fund.FundTypeCashDrawer)
	if err != nil {
		return nil, fmt.Errorf("failed to get cash_drawer balance: %w", err)
	}
	waiterBefore, err := s.repo.GetFundBalance(ctx, fund.FundTypeWaiterFloat)
	if err != nil {
		return nil, fmt.Errorf("failed to get waiter_float balance: %w", err)
	}

	desc := fmt.Sprintf("Bàn giao ca phục vụ %s - tiền mặt: %.0fđ, chuyển khoản: %.0fđ", waiterName, cashAmount, transferAmount)
	entry := fund.NewJournalEntry(
		fund.EventWaiterHandover,
		handoverID,
		desc,
		waiterID,
		waiterName,
		"waiter",
	)

	// DEBIT cash_drawer (drawer receives money back)
	drawerAfter := fund.FundBalance{
		Cash:     drawerBefore.Cash + cashAmount,
		Transfer: drawerBefore.Transfer + transferAmount,
		Total:    drawerBefore.Total + total,
	}
	entry.AddLine(fund.FundTypeCashDrawer, fund.DirectionDebit, cashAmount, transferAmount, *drawerBefore, drawerAfter)

	// CREDIT waiter_float (waiter returns money)
	waiterAfter := fund.FundBalance{
		Cash:     waiterBefore.Cash - cashAmount,
		Transfer: waiterBefore.Transfer - transferAmount,
		Total:    waiterBefore.Total - total,
	}
	entry.AddLine(fund.FundTypeWaiterFloat, fund.DirectionCredit, cashAmount, transferAmount, *waiterBefore, waiterAfter)

	if err := entry.Validate(); err != nil {
		return nil, err
	}

	if err := s.repo.Create(ctx, entry); err != nil {
		return nil, fmt.Errorf("failed to persist journal entry: %w", err)
	}
	return entry, nil
}

// RecordManagerDeposit creates:
//
//	DEBIT  fundType X  (quỹ nhận tiền từ chủ quán)
//	CREDIT owner    X  (chủ quán bỏ tiền vào)
func (s *JournalService) RecordManagerDeposit(
	ctx context.Context,
	fundType fund.FundType,
	cashAmount, transferAmount float64,
	reason string,
	performedBy primitive.ObjectID,
	performedByName, performedByRole string,
) (*fund.JournalEntry, error) {
	total := cashAmount + transferAmount
	if total <= 0 {
		return nil, fmt.Errorf("deposit amount must be greater than 0")
	}
	if len(reason) < 10 {
		return nil, fmt.Errorf("reason must be at least 10 characters")
	}

	session, err := s.mongoClient.StartSession()
	if err != nil {
		return nil, fmt.Errorf("failed to start session: %w", err)
	}
	defer session.EndSession(ctx)

	var entry *fund.JournalEntry
	err = mongoDriver.WithSession(ctx, session, func(sc mongoDriver.SessionContext) error {
		if err := session.StartTransaction(); err != nil {
			return err
		}

		fundBefore, err := s.repo.GetFundBalance(sc, fundType)
		if err != nil {
			session.AbortTransaction(sc)
			return fmt.Errorf("failed to get %s balance: %w", fundType, err)
		}
		ownerBefore, err := s.repo.GetFundBalance(sc, fund.FundTypeOwner)
		if err != nil {
			session.AbortTransaction(sc)
			return fmt.Errorf("failed to get owner balance: %w", err)
		}

		desc := fmt.Sprintf("Nạp tiền vào %s: %.0fđ - %s", fundType, total, reason)
		e := fund.NewJournalEntry(
			fund.EventManagerDeposit,
			primitive.NewObjectID(),
			desc,
			performedBy,
			performedByName,
			performedByRole,
		)

		// DEBIT fund (quỹ nhận tiền từ chủ)
		fundAfter := fund.FundBalance{
			Cash:     fundBefore.Cash + cashAmount,
			Transfer: fundBefore.Transfer + transferAmount,
			Total:    fundBefore.Total + total,
		}
		e.AddLine(fundType, fund.DirectionDebit, cashAmount, transferAmount, *fundBefore, fundAfter)

		// CREDIT owner (chủ bỏ tiền vào hệ thống)
		ownerAfter := fund.FundBalance{
			Cash:     ownerBefore.Cash - cashAmount,
			Transfer: ownerBefore.Transfer - transferAmount,
			Total:    ownerBefore.Total - total,
		}
		e.AddLine(fund.FundTypeOwner, fund.DirectionCredit, cashAmount, transferAmount, *ownerBefore, ownerAfter)

		if err := e.Validate(); err != nil {
			session.AbortTransaction(sc)
			return err
		}
		if err := s.repo.Create(sc, e); err != nil {
			session.AbortTransaction(sc)
			return fmt.Errorf("failed to persist journal entry: %w", err)
		}
		entry = e
		return session.CommitTransaction(sc)
	})
	return entry, err
}

// RecordManagerWithdrawal creates:
//
//	DEBIT  owner    X  (chủ quán rút tiền ra)
//	CREDIT fundType X  (quỹ xuất tiền)
func (s *JournalService) RecordManagerWithdrawal(
	ctx context.Context,
	fundType fund.FundType,
	cashAmount, transferAmount float64,
	reason string,
	performedBy primitive.ObjectID,
	performedByName, performedByRole string,
) (*fund.JournalEntry, error) {
	total := cashAmount + transferAmount
	if total <= 0 {
		return nil, fmt.Errorf("withdrawal amount must be greater than 0")
	}
	if len(reason) < 10 {
		return nil, fmt.Errorf("reason must be at least 10 characters")
	}

	session, err := s.mongoClient.StartSession()
	if err != nil {
		return nil, fmt.Errorf("failed to start session: %w", err)
	}
	defer session.EndSession(ctx)

	var entry *fund.JournalEntry
	err = mongoDriver.WithSession(ctx, session, func(sc mongoDriver.SessionContext) error {
		if err := session.StartTransaction(); err != nil {
			return err
		}

		fundBefore, err := s.repo.GetFundBalance(sc, fundType)
		if err != nil {
			session.AbortTransaction(sc)
			return fmt.Errorf("failed to get %s balance: %w", fundType, err)
		}
		// Validate sufficient balance
		if cashAmount > 0 && cashAmount > fundBefore.Cash {
			session.AbortTransaction(sc)
			return fmt.Errorf("insufficient cash in %s: have %.0f, need %.0f", fundType, fundBefore.Cash, cashAmount)
		}
		if transferAmount > 0 && transferAmount > fundBefore.Transfer {
			session.AbortTransaction(sc)
			return fmt.Errorf("insufficient transfer in %s: have %.0f, need %.0f", fundType, fundBefore.Transfer, transferAmount)
		}

		ownerBefore, err := s.repo.GetFundBalance(sc, fund.FundTypeOwner)
		if err != nil {
			session.AbortTransaction(sc)
			return fmt.Errorf("failed to get owner balance: %w", err)
		}

		desc := fmt.Sprintf("Rút tiền từ %s: %.0fđ - %s", fundType, total, reason)
		e := fund.NewJournalEntry(
			fund.EventManagerWithdrawal,
			primitive.NewObjectID(),
			desc,
			performedBy,
			performedByName,
			performedByRole,
		)

		// DEBIT owner (chủ nhận tiền ra ngoài)
		ownerAfter := fund.FundBalance{
			Cash:     ownerBefore.Cash + cashAmount,
			Transfer: ownerBefore.Transfer + transferAmount,
			Total:    ownerBefore.Total + total,
		}
		e.AddLine(fund.FundTypeOwner, fund.DirectionDebit, cashAmount, transferAmount, *ownerBefore, ownerAfter)

		// CREDIT fund (quỹ xuất tiền)
		fundAfter := fund.FundBalance{
			Cash:     fundBefore.Cash - cashAmount,
			Transfer: fundBefore.Transfer - transferAmount,
			Total:    fundBefore.Total - total,
		}
		e.AddLine(fundType, fund.DirectionCredit, cashAmount, transferAmount, *fundBefore, fundAfter)

		if err := e.Validate(); err != nil {
			session.AbortTransaction(sc)
			return err
		}
		if err := s.repo.Create(sc, e); err != nil {
			session.AbortTransaction(sc)
			return fmt.Errorf("failed to persist journal entry: %w", err)
		}
		entry = e
		return session.CommitTransaction(sc)
	})
	return entry, err
}

// RecordFundWithdrawalInSession records a fund withdrawal within an existing MongoDB session.
// Use this inside another service's transaction (no new session is created).
// Creates: DEBIT counterpart / CREDIT fundType
// counterpart: typically FundTypeSupplier for purchases, FundTypeOwner for manager withdrawal
func (s *JournalService) RecordFundWithdrawalInSession(
	ctx context.Context, // must already be a mongo.SessionContext
	eventType fund.JournalEventType,
	fundType fund.FundType,
	counterpart fund.FundType, // external counterpart account (owner/supplier/customer)
	cashAmount, transferAmount float64,
	description string,
	eventID primitive.ObjectID,
	performedBy primitive.ObjectID,
	performedByName, performedByRole string,
) (*fund.JournalEntry, error) {
	total := cashAmount + transferAmount
	if total <= 0 {
		return nil, fmt.Errorf("withdrawal amount must be greater than 0")
	}

	fundBefore, err := s.repo.GetFundBalance(ctx, fundType)
	if err != nil {
		return nil, fmt.Errorf("failed to get %s balance: %w", fundType, err)
	}
	if cashAmount > 0 && cashAmount > fundBefore.Cash {
		return nil, fmt.Errorf("insufficient cash in %s: have %.0f, need %.0f", fundType, fundBefore.Cash, cashAmount)
	}
	if transferAmount > 0 && transferAmount > fundBefore.Transfer {
		return nil, fmt.Errorf("insufficient transfer in %s: have %.0f, need %.0f", fundType, fundBefore.Transfer, transferAmount)
	}

	counterpartBefore, err := s.repo.GetFundBalance(ctx, counterpart)
	if err != nil {
		return nil, fmt.Errorf("failed to get %s balance: %w", counterpart, err)
	}

	e := fund.NewJournalEntry(eventType, eventID, description, performedBy, performedByName, performedByRole)

	// DEBIT counterpart (nhà cung cấp / chủ quán nhận tiền ra)
	counterpartAfter := fund.FundBalance{
		Cash:     counterpartBefore.Cash + cashAmount,
		Transfer: counterpartBefore.Transfer + transferAmount,
		Total:    counterpartBefore.Total + total,
	}
	e.AddLine(counterpart, fund.DirectionDebit, cashAmount, transferAmount, *counterpartBefore, counterpartAfter)

	// CREDIT fund (quỹ xuất tiền)
	fundAfter := fund.FundBalance{
		Cash:     fundBefore.Cash - cashAmount,
		Transfer: fundBefore.Transfer - transferAmount,
		Total:    fundBefore.Total - total,
	}
	e.AddLine(fundType, fund.DirectionCredit, cashAmount, transferAmount, *fundBefore, fundAfter)

	if err := e.Validate(); err != nil {
		return nil, err
	}
	if err := s.repo.Create(ctx, e); err != nil {
		return nil, fmt.Errorf("failed to persist journal entry: %w", err)
	}
	return e, nil
}

// RecordFundTransfer creates:
//   DEBIT  toFund   X  (destination receives money)
//   CREDIT fromFund X  (source sends money)
func (s *JournalService) RecordFundTransfer(
	ctx context.Context,
	fromFund, toFund fund.FundType,
	cashAmount, transferAmount float64,
	reason string,
	performedBy primitive.ObjectID,
	performedByName, performedByRole string,
) (*fund.JournalEntry, error) {
	if fromFund == toFund {
		return nil, fmt.Errorf("source and destination funds cannot be the same")
	}
	total := cashAmount + transferAmount
	if total <= 0 {
		return nil, fmt.Errorf("transfer amount must be greater than 0")
	}
	if len(reason) < 10 {
		return nil, fmt.Errorf("reason must be at least 10 characters")
	}

	session, err := s.mongoClient.StartSession()
	if err != nil {
		return nil, fmt.Errorf("failed to start session: %w", err)
	}
	defer session.EndSession(ctx)

	var entry *fund.JournalEntry
	err = mongoDriver.WithSession(ctx, session, func(sc mongoDriver.SessionContext) error {
		if err := session.StartTransaction(); err != nil {
			return err
		}

		fromBefore, err := s.repo.GetFundBalance(sc, fromFund)
		if err != nil {
			session.AbortTransaction(sc)
			return fmt.Errorf("failed to get source balance: %w", err)
		}
		if cashAmount > fromBefore.Cash {
			session.AbortTransaction(sc)
			return fmt.Errorf("insufficient cash in %s: have %.0f, need %.0f", fromFund, fromBefore.Cash, cashAmount)
		}
		if transferAmount > fromBefore.Transfer {
			session.AbortTransaction(sc)
			return fmt.Errorf("insufficient transfer in %s: have %.0f, need %.0f", fromFund, fromBefore.Transfer, transferAmount)
		}

		toBefore, err := s.repo.GetFundBalance(sc, toFund)
		if err != nil {
			session.AbortTransaction(sc)
			return fmt.Errorf("failed to get destination balance: %w", err)
		}

		e := fund.NewJournalEntry(
			fund.EventFundTransfer,
			primitive.NilObjectID,
			reason,
			performedBy,
			performedByName,
			performedByRole,
		)

		// DEBIT destination
		toAfter := fund.FundBalance{
			Cash:     toBefore.Cash + cashAmount,
			Transfer: toBefore.Transfer + transferAmount,
			Total:    toBefore.Total + total,
		}
		e.AddLine(toFund, fund.DirectionDebit, cashAmount, transferAmount, *toBefore, toAfter)

		// CREDIT source
		fromAfter := fund.FundBalance{
			Cash:     fromBefore.Cash - cashAmount,
			Transfer: fromBefore.Transfer - transferAmount,
			Total:    fromBefore.Total - total,
		}
		e.AddLine(fromFund, fund.DirectionCredit, cashAmount, transferAmount, *fromBefore, fromAfter)

		if err := e.Validate(); err != nil {
			session.AbortTransaction(sc)
			return err
		}

		if err := s.repo.Create(sc, e); err != nil {
			session.AbortTransaction(sc)
			return fmt.Errorf("failed to persist journal entry: %w", err)
		}
		entry = e
		return session.CommitTransaction(sc)
	})
	return entry, err
}
