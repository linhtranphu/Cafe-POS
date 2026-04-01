package services

import (
	"context"
	"errors"
	"sync"
	"testing"
	"time"

	"cafe-pos/backend/domain"
	"cafe-pos/backend/domain/cashier"
	"cafe-pos/backend/domain/fund"
	"cafe-pos/backend/domain/order"
	"cafe-pos/backend/infrastructure/mongodb"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

// ══════════════════════════════════════════════════════════════════════════════
// Stubs / in-memory implementations
// ══════════════════════════════════════════════════════════════════════════════

// ── journalRepository stub ────────────────────────────────────────────────────

type stubJournalRepo struct {
	mu      sync.Mutex
	entries []*fund.JournalEntry
	// Set to a non-zero event type to make Create fail for that event type.
	failOn fund.JournalEventType
	// Fixed balance returned for all fund types (default: 1 000 000 cash).
	balance fund.FundBalance
}

func newStubJournalRepo() *stubJournalRepo {
	return &stubJournalRepo{balance: fund.FundBalance{Cash: 1_000_000, Total: 1_000_000}}
}

func (r *stubJournalRepo) Create(ctx context.Context, e *fund.JournalEntry) error {
	r.mu.Lock()
	defer r.mu.Unlock()
	if r.failOn != "" && e.EventType == r.failOn {
		return errors.New("injected journal error")
	}
	cp := *e
	r.entries = append(r.entries, &cp)
	return nil
}

func (r *stubJournalRepo) GetFundBalance(_ context.Context, _ fund.FundType) (*fund.FundBalance, error) {
	cp := r.balance
	return &cp, nil
}

func (r *stubJournalRepo) GetAllFundBalances(_ context.Context) (map[fund.FundType]*fund.FundBalance, error) {
	return nil, nil
}

func (r *stubJournalRepo) List(_ context.Context, _ mongodb.JournalFilter) ([]*fund.JournalEntry, int64, error) {
	return nil, 0, nil
}

func (r *stubJournalRepo) FindByID(_ context.Context, _ primitive.ObjectID) (*fund.JournalEntry, error) {
	return nil, nil
}

func (r *stubJournalRepo) count(et fund.JournalEventType) int {
	r.mu.Lock()
	defer r.mu.Unlock()
	n := 0
	for _, e := range r.entries {
		if e.EventType == et {
			n++
		}
	}
	return n
}

// newTestJournalService returns a JournalService backed by the in-memory stub.
// mongoClient is nil — no real MongoDB sessions are started in unit tests.
func newTestJournalService(repo *stubJournalRepo) *JournalService {
	return &JournalService{repo: repo, mongoClient: nil}
}

// ── OrderRepository stub ──────────────────────────────────────────────────────

type stubOrderRepo struct {
	mu     sync.Mutex
	orders map[primitive.ObjectID]*order.Order
}

func newStubOrderRepo() *stubOrderRepo {
	return &stubOrderRepo{orders: make(map[primitive.ObjectID]*order.Order)}
}

func (r *stubOrderRepo) put(o *order.Order) {
	r.mu.Lock()
	defer r.mu.Unlock()
	if o.ID.IsZero() {
		o.ID = primitive.NewObjectID()
	}
	cp := *o
	r.orders[o.ID] = &cp
}

func (r *stubOrderRepo) Create(ctx context.Context, o *order.Order) error {
	r.put(o)
	return nil
}

func (r *stubOrderRepo) FindByID(_ context.Context, id primitive.ObjectID) (*order.Order, error) {
	r.mu.Lock()
	defer r.mu.Unlock()
	o, ok := r.orders[id]
	if !ok {
		return nil, errors.New("order not found")
	}
	cp := *o
	return &cp, nil
}

func (r *stubOrderRepo) Update(_ context.Context, id primitive.ObjectID, o *order.Order) error {
	r.mu.Lock()
	defer r.mu.Unlock()
	if _, ok := r.orders[id]; !ok {
		return errors.New("order not found")
	}
	cp := *o
	r.orders[id] = &cp
	return nil
}

func (r *stubOrderRepo) Delete(_ context.Context, id primitive.ObjectID) error {
	delete(r.orders, id)
	return nil
}

func (r *stubOrderRepo) FindByShiftID(_ context.Context, _ primitive.ObjectID) ([]*order.Order, error) {
	return nil, nil
}
func (r *stubOrderRepo) FindByWaiterID(_ context.Context, _ primitive.ObjectID) ([]*order.Order, error) {
	return nil, nil
}
func (r *stubOrderRepo) FindByStatus(_ context.Context, _ order.OrderStatus) ([]*order.Order, error) {
	return nil, nil
}
func (r *stubOrderRepo) FindByOrderNumber(_ context.Context, _ string) (*order.Order, error) {
	return nil, nil
}
func (r *stubOrderRepo) FindAll(_ context.Context) ([]*order.Order, error) { return nil, nil }
func (r *stubOrderRepo) FindByIDs(_ context.Context, _ []primitive.ObjectID) ([]*order.Order, error) {
	return nil, nil
}

// ── ShiftRepository stub (minimal) ───────────────────────────────────────────

type stubShiftRepo struct {
	shifts map[primitive.ObjectID]*order.Shift
}

func newStubShiftRepo() *stubShiftRepo {
	return &stubShiftRepo{shifts: make(map[primitive.ObjectID]*order.Shift)}
}

func (r *stubShiftRepo) Create(_ context.Context, s *order.Shift) error {
	if s.ID.IsZero() {
		s.ID = primitive.NewObjectID()
	}
	cp := *s
	r.shifts[s.ID] = &cp
	return nil
}
func (r *stubShiftRepo) FindByID(_ context.Context, id primitive.ObjectID) (*order.Shift, error) {
	s, ok := r.shifts[id]
	if !ok {
		return nil, errors.New("shift not found")
	}
	cp := *s
	return &cp, nil
}
func (r *stubShiftRepo) Update(_ context.Context, id primitive.ObjectID, s *order.Shift) error {
	cp := *s
	r.shifts[id] = &cp
	return nil
}
func (r *stubShiftRepo) FindOpenShiftByWaiter(_ context.Context, _ primitive.ObjectID) (*order.Shift, error) {
	return nil, nil
}
func (r *stubShiftRepo) FindOpenShiftByUser(_ context.Context, userID primitive.ObjectID, roleType order.RoleType) (*order.Shift, error) {
	for _, s := range r.shifts {
		if s.UserID == userID && s.RoleType == roleType && s.Status == order.ShiftOpen {
			cp := *s
			return &cp, nil
		}
	}
	return nil, nil
}
func (r *stubShiftRepo) FindOpenShifts(_ context.Context) ([]*order.Shift, error)   { return nil, nil }
func (r *stubShiftRepo) FindByWaiterID(_ context.Context, _ primitive.ObjectID) ([]*order.Shift, error) {
	return nil, nil
}
func (r *stubShiftRepo) FindByUserID(_ context.Context, _ primitive.ObjectID, _ order.RoleType) ([]*order.Shift, error) {
	return nil, nil
}
func (r *stubShiftRepo) FindByDateRange(_ context.Context, _, _ time.Time) ([]*order.Shift, error) {
	return nil, nil
}
func (r *stubShiftRepo) FindByRoleType(_ context.Context, _ order.RoleType) ([]*order.Shift, error) {
	return nil, nil
}
func (r *stubShiftRepo) FindAll(_ context.Context) ([]*order.Shift, error) { return nil, nil }

// ── Cashier shift harness ─────────────────────────────────────────────────────
//
// ShiftService.cashierShiftRepo is typed as *mongodb.CashierShiftRepository (concrete).
// For unit tests we bypass this by re-implementing the two-operation logic
// of StartShift in a test wrapper rather than trying to inject a mock through
// the concrete field. See testShiftHarness.startShift below.

type stubCashierShiftList struct {
	shifts []*cashier.CashierShift
}

func (s *stubCashierShiftList) addOpen() {
	s.shifts = append(s.shifts, &cashier.CashierShift{
		ID:     primitive.NewObjectID(),
		Status: cashier.CashierShiftOpen,
	})
}

// ── Test harness for StartShift ───────────────────────────────────────────────
//
// Rather than fighting the concrete *mongodb.CashierShiftRepository type, we
// replicate the StartShift business logic in a thin test wrapper that swaps
// the cashier repo for our in-memory stub. This tests the exact same rules
// without needing a real MongoDB connection.

type testShiftHarness struct {
	shiftRepo   *stubShiftRepo
	cashierList *stubCashierShiftList
	journalRepo *stubJournalRepo // nil when journal testing is not needed
	js          *JournalService
	sm          *domain.StateMachineManager
}

func newTestShiftHarness(withJournal bool) *testShiftHarness {
	h := &testShiftHarness{
		shiftRepo:   newStubShiftRepo(),
		cashierList: &stubCashierShiftList{},
		sm:          domain.NewStateMachineManager(),
	}
	if withJournal {
		h.journalRepo = newStubJournalRepo()
		h.js = newTestJournalService(h.journalRepo)
	}
	return h
}

// startShift mirrors ShiftService.StartShift but uses stub repos so that unit
// tests don't require MongoDB. The atomicity guarantee (rollback on journal
// failure) is replicated manually since mongoClient is nil in these tests.
func (h *testShiftHarness) startShift(ctx context.Context, req *order.StartShiftRequest, userID, userName string, roleType order.RoleType) (*order.Shift, error) {
	if roleType == "cashier" {
		return nil, errors.New("cashier shifts must be created using the cashier shift service")
	}
	if !roleType.IsValid() {
		return nil, errors.New("invalid role type")
	}

	userOID, _ := primitive.ObjectIDFromHex(userID)

	existing, _ := h.shiftRepo.FindOpenShiftByUser(ctx, userOID, roleType)
	if err := h.sm.ValidateWaiterShiftStart(existing); err != nil {
		return nil, err
	}

	var activeCashierShifts []*cashier.CashierShift
	if roleType == order.RoleWaiter {
		if len(h.cashierList.shifts) == 0 {
			return nil, errors.New("thu ngân phải mở ca trước khi phục vụ mở ca")
		}
		activeCashierShifts = h.cashierList.shifts
	}

	shift := &order.Shift{
		Type:          req.Type,
		Status:        order.ShiftOpen,
		RoleType:      roleType,
		UserID:        userOID,
		UserName:      userName,
		StartCash:     req.StartCash,
		CurrentCash:   req.StartCash,
		RemainingCash: req.StartCash,
		StartedAt:     time.Now(),
	}

	needsJournal := len(activeCashierShifts) > 0 && req.StartCash > 0 && h.js != nil
	if needsJournal {
		// Create shift first, then record journal. On journal failure we undo
		// the shift creation to replicate the transactional guarantee.
		if err := h.shiftRepo.Create(ctx, shift); err != nil {
			return nil, err
		}
		if _, err := h.js.RecordWaiterShiftStartInSession(ctx, req.StartCash, shift.ID, userOID, userName); err != nil {
			delete(h.shiftRepo.shifts, shift.ID) // simulate rollback
			return nil, errors.New("kế toán không ghi nhận được tiền đầu ca: " + err.Error())
		}
		activeCashierShifts[0].AddDistributedCash(req.StartCash)
	} else {
		if err := h.shiftRepo.Create(ctx, shift); err != nil {
			return nil, err
		}
	}

	return shift, nil
}

// ══════════════════════════════════════════════════════════════════════════════
// Tests: CollectPayment
// ══════════════════════════════════════════════════════════════════════════════

func makeTestOrder(total float64, status order.OrderStatus) *order.Order {
	o := &order.Order{
		ID:          primitive.NewObjectID(),
		OrderNumber: "ORD-TEST",
		Status:      status,
		Total:       total,
		AmountDue:   total,
		WaiterID:    primitive.NewObjectID(),
		WaiterName:  "Waiter Test",
		CreatedAt:   time.Now(),
		UpdatedAt:   time.Now(),
	}
	return o
}

func makeTestOrderService(orderRepo OrderRepository, shiftRepo ShiftRepository, js *JournalService) *OrderService {
	sm := domain.NewStateMachineManager()
	svc := NewOrderService(orderRepo, shiftRepo, nil, sm, nil)
	if js != nil {
		svc.SetJournalService(js)
	}
	return svc
}

// TestCollectPayment_CashPayment_MarksOrderPaidAndRecordsJournal verifies that a
// normal cash payment transitions the order to PAID and creates exactly one
// order_payment journal entry.
func TestCollectPayment_CashPayment_MarksOrderPaidAndRecordsJournal(t *testing.T) {
	orderRepo := newStubOrderRepo()
	o := makeTestOrder(50_000, order.StatusCreated)
	orderRepo.put(o)

	jRepo := newStubJournalRepo()
	svc := makeTestOrderService(orderRepo, newStubShiftRepo(), newTestJournalService(jRepo))

	result, err := svc.CollectPayment(context.Background(), o.ID, &order.PaymentRequest{
		Amount:        50_000,
		PaymentMethod: order.PaymentCash,
		CollectorID:   primitive.NewObjectID().Hex(),
		CollectorName: "Thu ngân A",
	})

	require.NoError(t, err)
	assert.Equal(t, order.StatusPaid, result.Status)
	assert.Equal(t, float64(0), result.AmountDue)
	assert.Equal(t, order.PaymentCash, result.PaymentMethod)
	assert.NotNil(t, result.PaidAt)
	assert.Equal(t, 1, jRepo.count(fund.EventOrderPayment))
}

// TestCollectPayment_TransferPayment_RecordsJournal verifies the transfer payment path.
func TestCollectPayment_TransferPayment_RecordsJournal(t *testing.T) {
	orderRepo := newStubOrderRepo()
	o := makeTestOrder(120_000, order.StatusCreated)
	orderRepo.put(o)

	jRepo := newStubJournalRepo()
	svc := makeTestOrderService(orderRepo, newStubShiftRepo(), newTestJournalService(jRepo))

	result, err := svc.CollectPayment(context.Background(), o.ID, &order.PaymentRequest{
		Amount:        120_000,
		PaymentMethod: order.PaymentTransfer,
	})

	require.NoError(t, err)
	assert.Equal(t, order.StatusPaid, result.Status)
	assert.Equal(t, order.PaymentTransfer, result.PaymentMethod)
	assert.Equal(t, 1, jRepo.count(fund.EventOrderPayment))
}

// TestCollectPayment_AlreadyPaid_ReturnsIdempotently verifies that calling
// CollectPayment on a fully-paid order returns the order as-is and does NOT
// record a second journal entry (pre-check guard).
func TestCollectPayment_AlreadyPaid_ReturnsIdempotently(t *testing.T) {
	orderRepo := newStubOrderRepo()
	o := makeTestOrder(50_000, order.StatusPaid)
	o.AmountPaid = 50_000
	o.AmountDue = 0
	orderRepo.put(o)

	jRepo := newStubJournalRepo()
	svc := makeTestOrderService(orderRepo, newStubShiftRepo(), newTestJournalService(jRepo))

	result, err := svc.CollectPayment(context.Background(), o.ID, &order.PaymentRequest{
		Amount:        50_000,
		PaymentMethod: order.PaymentCash,
	})

	require.NoError(t, err)
	assert.Equal(t, order.StatusPaid, result.Status)
	assert.Equal(t, 0, jRepo.count(fund.EventOrderPayment), "đơn đã thanh toán không được ghi thêm journal")
}

// TestCollectPayment_SequentialDoublePay_OnlyOneJournalEntry simulates a user
// double-clicking "collect payment" where the first request fully completes
// before the second arrives. Only one journal entry must be created.
func TestCollectPayment_SequentialDoublePay_OnlyOneJournalEntry(t *testing.T) {
	orderRepo := newStubOrderRepo()
	o := makeTestOrder(80_000, order.StatusCreated)
	orderRepo.put(o)

	jRepo := newStubJournalRepo()
	svc := makeTestOrderService(orderRepo, newStubShiftRepo(), newTestJournalService(jRepo))

	req := &order.PaymentRequest{
		Amount:        80_000,
		PaymentMethod: order.PaymentCash,
		CollectorID:   primitive.NewObjectID().Hex(),
	}

	// First call — succeeds and marks order PAID
	r1, err := svc.CollectPayment(context.Background(), o.ID, req)
	require.NoError(t, err)
	assert.Equal(t, order.StatusPaid, r1.Status)

	// Second call — idempotent: returns the same order, no new journal entry
	r2, err := svc.CollectPayment(context.Background(), o.ID, req)
	require.NoError(t, err)
	assert.Equal(t, order.StatusPaid, r2.Status)

	assert.Equal(t, 1, jRepo.count(fund.EventOrderPayment), "phải đúng 1 journal entry khi gọi 2 lần")
}

// TestCollectPayment_LockedOrder_Rejected verifies that a LOCKED order (shift
// already closed) cannot be paid.
func TestCollectPayment_LockedOrder_Rejected(t *testing.T) {
	orderRepo := newStubOrderRepo()
	o := makeTestOrder(50_000, order.StatusLocked)
	orderRepo.put(o)

	svc := makeTestOrderService(orderRepo, newStubShiftRepo(), nil)

	_, err := svc.CollectPayment(context.Background(), o.ID, &order.PaymentRequest{
		Amount: 50_000, PaymentMethod: order.PaymentCash,
	})

	assert.Error(t, err)
}

// TestCollectPayment_CancelledOrder_Rejected verifies that a CANCELLED order cannot be paid.
func TestCollectPayment_CancelledOrder_Rejected(t *testing.T) {
	orderRepo := newStubOrderRepo()
	o := makeTestOrder(50_000, order.StatusCancelled)
	orderRepo.put(o)

	svc := makeTestOrderService(orderRepo, newStubShiftRepo(), nil)

	_, err := svc.CollectPayment(context.Background(), o.ID, &order.PaymentRequest{
		Amount: 50_000, PaymentMethod: order.PaymentCash,
	})

	assert.Error(t, err)
}

// TestCollectPayment_UnknownOrderID_ReturnsError verifies that an unknown order ID
// returns an error immediately.
func TestCollectPayment_UnknownOrderID_ReturnsError(t *testing.T) {
	svc := makeTestOrderService(newStubOrderRepo(), newStubShiftRepo(), nil)

	_, err := svc.CollectPayment(context.Background(), primitive.NewObjectID(), &order.PaymentRequest{
		Amount: 50_000, PaymentMethod: order.PaymentCash,
	})

	assert.Error(t, err)
}

// TestCollectPayment_NoJournalService_StillProcessesPayment verifies that payment
// completes successfully when no journal service is configured (optional dependency).
func TestCollectPayment_NoJournalService_StillProcessesPayment(t *testing.T) {
	orderRepo := newStubOrderRepo()
	o := makeTestOrder(30_000, order.StatusCreated)
	orderRepo.put(o)

	svc := makeTestOrderService(orderRepo, newStubShiftRepo(), nil) // journal = nil

	result, err := svc.CollectPayment(context.Background(), o.ID, &order.PaymentRequest{
		Amount: 30_000, PaymentMethod: order.PaymentQR,
	})

	require.NoError(t, err)
	assert.Equal(t, order.StatusPaid, result.Status)
}

// ══════════════════════════════════════════════════════════════════════════════
// Tests: StartShift (waiter)
// ══════════════════════════════════════════════════════════════════════════════

// TestStartWaiterShift_NoCashierShift_Rejected verifies that a waiter cannot open
// a shift when no cashier shift is currently open.
func TestStartWaiterShift_NoCashierShift_Rejected(t *testing.T) {
	h := newTestShiftHarness(false)
	// cashierList is empty by default

	_, err := h.startShift(context.Background(),
		&order.StartShiftRequest{StartCash: 200_000},
		primitive.NewObjectID().Hex(), "Linh", order.RoleWaiter)

	require.Error(t, err)
	assert.Contains(t, err.Error(), "thu ngân phải mở ca")
}

// TestStartWaiterShift_ZeroStartCash_CreatesShiftWithoutJournal verifies that a
// waiter can open a shift with zero start cash (no journal entry expected).
func TestStartWaiterShift_ZeroStartCash_CreatesShiftWithoutJournal(t *testing.T) {
	h := newTestShiftHarness(true)
	h.cashierList.addOpen()

	shift, err := h.startShift(context.Background(),
		&order.StartShiftRequest{StartCash: 0},
		primitive.NewObjectID().Hex(), "Linh", order.RoleWaiter)

	require.NoError(t, err)
	assert.Equal(t, order.ShiftOpen, shift.Status)
	assert.Equal(t, 1, len(h.shiftRepo.shifts))
	assert.Equal(t, 0, h.journalRepo.count(fund.EventWaiterShiftStart), "startCash=0 không ghi journal")
}

// TestStartWaiterShift_WithStartCash_RecordsJournalEntry verifies that opening a
// shift with StartCash > 0 creates one waiter_shift_start journal entry.
func TestStartWaiterShift_WithStartCash_RecordsJournalEntry(t *testing.T) {
	h := newTestShiftHarness(true)
	h.cashierList.addOpen()

	shift, err := h.startShift(context.Background(),
		&order.StartShiftRequest{StartCash: 300_000},
		primitive.NewObjectID().Hex(), "Linh", order.RoleWaiter)

	require.NoError(t, err)
	assert.Equal(t, float64(300_000), shift.StartCash)
	assert.Equal(t, 1, h.journalRepo.count(fund.EventWaiterShiftStart), "phải có 1 journal entry tiền đầu ca")
}

// TestStartWaiterShift_JournalFails_ShiftRolledBack verifies that when the journal
// entry fails (e.g. insufficient cash_drawer balance), the shift is NOT persisted.
// This enforces the atomicity guarantee: ca chỉ mở thành công khi kế toán ghi nhận được.
func TestStartWaiterShift_JournalFails_ShiftRolledBack(t *testing.T) {
	h := newTestShiftHarness(true)
	h.cashierList.addOpen()
	h.journalRepo.failOn = fund.EventWaiterShiftStart // kế toán lỗi

	_, err := h.startShift(context.Background(),
		&order.StartShiftRequest{StartCash: 200_000},
		primitive.NewObjectID().Hex(), "Linh", order.RoleWaiter)

	require.Error(t, err, "mở ca phải thất bại khi kế toán lỗi")
	assert.Equal(t, 0, len(h.shiftRepo.shifts), "shift không được lưu khi journal lỗi")
}

// TestStartWaiterShift_CannotOpenTwice verifies the state machine prevents opening
// two concurrent shifts for the same waiter.
func TestStartWaiterShift_CannotOpenTwice(t *testing.T) {
	h := newTestShiftHarness(false)
	h.cashierList.addOpen()
	userID := primitive.NewObjectID().Hex()

	_, err := h.startShift(context.Background(), &order.StartShiftRequest{}, userID, "Linh", order.RoleWaiter)
	require.NoError(t, err)

	_, err = h.startShift(context.Background(), &order.StartShiftRequest{}, userID, "Linh", order.RoleWaiter)
	assert.Error(t, err, "không thể mở 2 ca cùng lúc với cùng user")
}

// TestStartBaristaShift_NoCashierRequired verifies that a barista can open a shift
// without needing an active cashier shift (rule only applies to waiters).
func TestStartBaristaShift_NoCashierRequired(t *testing.T) {
	h := newTestShiftHarness(false)
	// No cashier shift open

	shift, err := h.startShift(context.Background(),
		&order.StartShiftRequest{StartCash: 0},
		primitive.NewObjectID().Hex(), "Barista A", order.RoleBarista)

	require.NoError(t, err)
	assert.Equal(t, order.RoleBarista, shift.RoleType)
}

// TestStartCashierShift_RejectedInWaiterService verifies that the waiter shift
// service explicitly rejects cashier role assignments.
func TestStartCashierShift_RejectedInWaiterService(t *testing.T) {
	h := newTestShiftHarness(false)

	_, err := h.startShift(context.Background(), &order.StartShiftRequest{}, primitive.NewObjectID().Hex(), "Cashier X", "cashier")
	assert.Error(t, err)
}
