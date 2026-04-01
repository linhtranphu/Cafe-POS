package services

import (
	"context"
	"errors"
	"fmt"
	"time"

	"cafe-pos/backend/domain"
	"cafe-pos/backend/domain/cashier"
	"cafe-pos/backend/domain/order"
	"cafe-pos/backend/infrastructure/mongodb"
	"go.mongodb.org/mongo-driver/bson/primitive"
	mongoDriver "go.mongodb.org/mongo-driver/mongo"
)

type ShiftRepository interface {
	Create(ctx context.Context, s *order.Shift) error
	FindByID(ctx context.Context, id primitive.ObjectID) (*order.Shift, error)
	Update(ctx context.Context, id primitive.ObjectID, s *order.Shift) error
	FindOpenShiftByWaiter(ctx context.Context, waiterID primitive.ObjectID) (*order.Shift, error)
	FindOpenShiftByUser(ctx context.Context, userID primitive.ObjectID, roleType order.RoleType) (*order.Shift, error)
	FindOpenShifts(ctx context.Context) ([]*order.Shift, error)
	FindByWaiterID(ctx context.Context, waiterID primitive.ObjectID) ([]*order.Shift, error)
	FindByUserID(ctx context.Context, userID primitive.ObjectID, roleType order.RoleType) ([]*order.Shift, error)
	FindByDateRange(ctx context.Context, startDate, endDate time.Time) ([]*order.Shift, error)
	FindByRoleType(ctx context.Context, roleType order.RoleType) ([]*order.Shift, error)
	FindAll(ctx context.Context) ([]*order.Shift, error)
}

type ShiftService struct {
	shiftRepo           ShiftRepository
	orderRepo           OrderRepository
	stateMachineManager *domain.StateMachineManager
	journalService      *JournalService
	cashierShiftRepo    *mongodb.CashierShiftRepository
	mongoClient         *mongoDriver.Client
}

func NewShiftService(
	shiftRepo ShiftRepository,
	orderRepo OrderRepository,
	stateMachineManager *domain.StateMachineManager,
	journalService *JournalService,
	cashierShiftRepo *mongodb.CashierShiftRepository,
	mongoClient *mongoDriver.Client,
) *ShiftService {
	return &ShiftService{
		shiftRepo:           shiftRepo,
		orderRepo:           orderRepo,
		stateMachineManager: stateMachineManager,
		journalService:      journalService,
		cashierShiftRepo:    cashierShiftRepo,
		mongoClient:         mongoClient,
	}
}

func (s *ShiftService) StartShift(ctx context.Context, req *order.StartShiftRequest, userID, userName string, roleType order.RoleType) (*order.Shift, error) {
	// Reject cashier role - cashier shifts are handled separately
	if roleType == "cashier" {
		return nil, errors.New("cashier shifts must be created using the cashier shift service")
	}
	
	// Validate role type is waiter or barista
	if !roleType.IsValid() {
		return nil, errors.New("invalid role type: must be waiter or barista")
	}
	
	userOID, _ := primitive.ObjectIDFromHex(userID)

	// Check if user already has an open shift for this role
	existingShift, _ := s.shiftRepo.FindOpenShiftByUser(ctx, userOID, roleType)

	// Validate using state machine
	if err := s.stateMachineManager.ValidateWaiterShiftStart(existingShift); err != nil {
		return nil, err
	}

	// Waiter phải có cashier ca mở trước khi mở ca — validate BEFORE creating the shift
	var activeCashierShifts []*cashier.CashierShift
	if roleType == order.RoleWaiter && s.cashierShiftRepo != nil {
		openCashierShifts, err := s.cashierShiftRepo.FindOpen(ctx)
		if err != nil {
			return nil, errors.New("không thể kiểm tra ca thu ngân")
		}
		if len(openCashierShifts) == 0 {
			return nil, errors.New("thu ngân phải mở ca trước khi phục vụ mở ca")
		}
		activeCashierShifts = openCashierShifts
	}

	shift := &order.Shift{
		Type:               req.Type,
		Status:             order.ShiftOpen,
		RoleType:           roleType,
		UserID:             userOID,
		UserName:           userName,
		StartCash:          req.StartCash,
		CurrentCash:        req.StartCash,
		RemainingCash:      req.StartCash,
		TransferRevenue:    0,
		RemainingTransfer:  0,
		HandedOverCash:     0,
		HandedOverTransfer: 0,
		StartedAt:          time.Now(),
	}

	// Wrap shift creation + journal entry + cashier update in a single transaction
	// so that a failed journal recording also rolls back the shift creation.
	needsJournal := len(activeCashierShifts) > 0 && req.StartCash > 0 && s.journalService != nil && s.mongoClient != nil
	if needsJournal {
		activeCashierShift := activeCashierShifts[0]
		session, err := s.mongoClient.StartSession()
		if err != nil {
			return nil, fmt.Errorf("không thể bắt đầu transaction: %w", err)
		}
		defer session.EndSession(ctx)

		err = mongoDriver.WithSession(ctx, session, func(sc mongoDriver.SessionContext) error {
			if err := session.StartTransaction(); err != nil {
				return err
			}
			if err := s.shiftRepo.Create(sc, shift); err != nil {
				session.AbortTransaction(sc)
				return err
			}
			if _, err := s.journalService.RecordWaiterShiftStartInSession(sc, req.StartCash, shift.ID, userOID, userName); err != nil {
				session.AbortTransaction(sc)
				return fmt.Errorf("kế toán không ghi nhận được tiền đầu ca: %w", err)
			}
			activeCashierShift.AddDistributedCash(req.StartCash)
			if err := s.cashierShiftRepo.Save(sc, activeCashierShift); err != nil {
				session.AbortTransaction(sc)
				return fmt.Errorf("không thể cập nhật ca thu ngân: %w", err)
			}
			return session.CommitTransaction(sc)
		})
		if err != nil {
			return nil, err
		}
	} else {
		if err := s.shiftRepo.Create(ctx, shift); err != nil {
			return nil, err
		}
	}

	return shift, nil
}

func (s *ShiftService) EndShift(ctx context.Context, shiftID primitive.ObjectID, req *order.EndShiftRequest) (*order.Shift, error) {
	shift, err := s.shiftRepo.FindByID(ctx, shiftID)
	if err != nil {
		return nil, err
	}

	// Validate state transition using state machine
	if err := s.stateMachineManager.ValidateWaiterShiftTransition(shift, order.EventEndShift); err != nil {
		return nil, err
	}

	orders, err := s.orderRepo.FindByShiftID(ctx, shiftID)
	if err != nil {
		return nil, err
	}

	totalRevenue := 0.0
	for _, o := range orders {
		if o.Status == order.StatusPaid || o.Status == order.StatusInProgress || o.Status == order.StatusServed {
			totalRevenue += o.Total
		}
	}

	now := time.Now()
	shift.Status = order.ShiftClosed
	shift.EndCash = req.EndCash
	shift.TotalRevenue = totalRevenue
	shift.TotalOrders = len(orders)
	shift.EndedAt = &now

	if err := s.shiftRepo.Update(ctx, shiftID, shift); err != nil {
		return nil, err
	}
	return shift, nil
}

func (s *ShiftService) GetCurrentShift(ctx context.Context, userID primitive.ObjectID, roleType order.RoleType) (*order.Shift, error) {
	return s.shiftRepo.FindOpenShiftByUser(ctx, userID, roleType)
}

func (s *ShiftService) GetOpenShifts(ctx context.Context) ([]*order.Shift, error) {
	return s.shiftRepo.FindOpenShifts(ctx)
}

func (s *ShiftService) GetShiftsByUser(ctx context.Context, userID primitive.ObjectID, roleType order.RoleType) ([]*order.Shift, error) {
	return s.shiftRepo.FindByUserID(ctx, userID, roleType)
}

func (s *ShiftService) GetShiftsByWaiter(ctx context.Context, waiterID primitive.ObjectID) ([]*order.Shift, error) {
	return s.shiftRepo.FindByWaiterID(ctx, waiterID)
}

func (s *ShiftService) GetShiftsByRole(ctx context.Context, roleType order.RoleType) ([]*order.Shift, error) {
	return s.shiftRepo.FindByRoleType(ctx, roleType)
}

func (s *ShiftService) GetAllShifts(ctx context.Context) ([]*order.Shift, error) {
	return s.shiftRepo.FindAll(ctx)
}

func (s *ShiftService) GetShift(ctx context.Context, id primitive.ObjectID) (*order.Shift, error) {
	return s.shiftRepo.FindByID(ctx, id)
}

func (s *ShiftService) GetPendingOrdersByShift(ctx context.Context, shiftID primitive.ObjectID) ([]*order.Order, error) {
	orders, err := s.orderRepo.FindByShiftID(ctx, shiftID)
	if err != nil {
		return nil, err
	}

	pendingStatuses := map[order.OrderStatus]bool{
		order.StatusCreated:    true,
		order.StatusPaid:       true,
		order.StatusQueued:     true,
		order.StatusInProgress: true,
		order.StatusReady:      true,
	}

	var pending []*order.Order
	for _, o := range orders {
		if pendingStatuses[o.Status] {
			pending = append(pending, o)
		}
	}
	return pending, nil
}

func (s *ShiftService) CloseShiftAndLockOrders(ctx context.Context, shiftID primitive.ObjectID, req *order.EndShiftRequest) (*order.Shift, error) {
	// Get shift first to validate
	shift, err := s.shiftRepo.FindByID(ctx, shiftID)
	if err != nil {
		return nil, err
	}

	// Validate state transition using state machine
	if err := s.stateMachineManager.ValidateWaiterShiftTransition(shift, order.EventEndShift); err != nil {
		return nil, err
	}

	// Kiểm tra không còn order chưa thanh toán trước khi đóng ca
	orders, _ := s.orderRepo.FindByShiftID(ctx, shiftID)
	unpaidCount := 0
	for _, o := range orders {
		if o.Status == order.StatusCreated {
			unpaidCount++
		}
	}
	if unpaidCount > 0 {
		return nil, fmt.Errorf("ca còn %d order chưa thanh toán, vui lòng xử lý trước khi đóng ca", unpaidCount)
	}

	// End the shift
	shift, err = s.EndShift(ctx, shiftID, req)
	if err != nil {
		return nil, err
	}

	now := time.Now()
	for _, o := range orders {
		switch o.Status {
		case order.StatusCreated:
			// Không thể xảy ra (đã kiểm tra ở trên), bỏ qua
			continue
		case order.StatusPaid, order.StatusQueued, order.StatusInProgress, order.StatusReady, order.StatusServed, order.StatusCancelled:
			// Khoá đơn đã thanh toán và đơn hoàn thành
			o.Status = order.StatusLocked
			o.LockedAt = &now
			o.UpdatedAt = now
			s.orderRepo.Update(ctx, o.ID, o)
		}
	}

	return shift, nil
}

// CalculateTransferRevenue calculates transfer revenue from orders and updates shift
func (s *ShiftService) CalculateTransferRevenue(ctx context.Context, shiftID primitive.ObjectID) error {
	shift, err := s.shiftRepo.FindByID(ctx, shiftID)
	if err != nil {
		return err
	}

	orders, err := s.orderRepo.FindByShiftID(ctx, shiftID)
	if err != nil {
		return err
	}

	// Calculate cash and transfer revenue separately
	cashRevenue := 0.0
	transferRevenue := 0.0
	for _, o := range orders {
		if o.Status == order.StatusPaid || o.Status == order.StatusInProgress || o.Status == order.StatusServed {
			if o.PaymentMethod == order.PaymentCash {
				cashRevenue += o.Total
			} else if o.PaymentMethod == order.PaymentTransfer || o.PaymentMethod == order.PaymentQR {
				transferRevenue += o.Total
			}
		}
	}

	// Update shift with transfer revenue
	shift.TransferRevenue = transferRevenue
	shift.RemainingTransfer = transferRevenue - shift.HandedOverTransfer
	shift.CurrentCash = shift.StartCash + cashRevenue
	shift.RemainingCash = shift.CurrentCash - shift.HandedOverCash
	shift.TotalRevenue = cashRevenue + transferRevenue
	shift.UpdatedAt = time.Now()

	return s.shiftRepo.Update(ctx, shiftID, shift)
}
