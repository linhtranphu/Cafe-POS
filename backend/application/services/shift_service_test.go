package services

import (
	"context"
	"errors"
	"testing"
	"time"

	"cafe-pos/backend/domain/order"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

// Mock ShiftRepository
type MockShiftRepository struct {
	shifts       map[string]*order.Shift
	createError  error
	findError    error
	updateError  error
}

func NewMockShiftRepository() *MockShiftRepository {
	return &MockShiftRepository{
		shifts: make(map[string]*order.Shift),
	}
}

func (m *MockShiftRepository) Create(ctx context.Context, s *order.Shift) error {
	if m.createError != nil {
		return m.createError
	}
	s.ID = primitive.NewObjectID()
	key := s.UserID.Hex() + "_" + string(s.RoleType)
	m.shifts[key] = s
	return nil
}

func (m *MockShiftRepository) FindByID(ctx context.Context, id primitive.ObjectID) (*order.Shift, error) {
	if m.findError != nil {
		return nil, m.findError
	}
	for _, shift := range m.shifts {
		if shift.ID == id {
			return shift, nil
		}
	}
	return nil, errors.New("shift not found")
}

func (m *MockShiftRepository) Update(ctx context.Context, id primitive.ObjectID, s *order.Shift) error {
	if m.updateError != nil {
		return m.updateError
	}
	for key, shift := range m.shifts {
		if shift.ID == id {
			m.shifts[key] = s
			return nil
		}
	}
	return errors.New("shift not found")
}

func (m *MockShiftRepository) FindOpenShiftByWaiter(ctx context.Context, waiterID primitive.ObjectID) (*order.Shift, error) {
	if m.findError != nil {
		return nil, m.findError
	}
	for _, shift := range m.shifts {
		if shift.WaiterID == waiterID && shift.Status == order.ShiftOpen {
			return shift, nil
		}
	}
	return nil, errors.New("no open shift found")
}

func (m *MockShiftRepository) FindOpenShiftByUser(ctx context.Context, userID primitive.ObjectID, roleType order.RoleType) (*order.Shift, error) {
	if m.findError != nil {
		return nil, m.findError
	}
	key := userID.Hex() + "_" + string(roleType)
	if shift, exists := m.shifts[key]; exists && shift.Status == order.ShiftOpen {
		return shift, nil
	}
	return nil, errors.New("no open shift found")
}

func (m *MockShiftRepository) FindOpenShifts(ctx context.Context) ([]*order.Shift, error) {
	if m.findError != nil {
		return nil, m.findError
	}
	var openShifts []*order.Shift
	for _, shift := range m.shifts {
		if shift.Status == order.ShiftOpen {
			openShifts = append(openShifts, shift)
		}
	}
	return openShifts, nil
}

func (m *MockShiftRepository) FindByWaiterID(ctx context.Context, waiterID primitive.ObjectID) ([]*order.Shift, error) {
	if m.findError != nil {
		return nil, m.findError
	}
	var shifts []*order.Shift
	for _, shift := range m.shifts {
		if shift.WaiterID == waiterID {
			shifts = append(shifts, shift)
		}
	}
	return shifts, nil
}

func (m *MockShiftRepository) FindByUserID(ctx context.Context, userID primitive.ObjectID, roleType order.RoleType) ([]*order.Shift, error) {
	if m.findError != nil {
		return nil, m.findError
	}
	var shifts []*order.Shift
	for _, shift := range m.shifts {
		if shift.UserID == userID && shift.RoleType == roleType {
			shifts = append(shifts, shift)
		}
	}
	return shifts, nil
}

func (m *MockShiftRepository) FindByDateRange(ctx context.Context, startDate, endDate time.Time) ([]*order.Shift, error) {
	if m.findError != nil {
		return nil, m.findError
	}
	var shifts []*order.Shift
	for _, shift := range m.shifts {
		if shift.StartedAt.After(startDate) && shift.StartedAt.Before(endDate) {
			shifts = append(shifts, shift)
		}
	}
	return shifts, nil
}

func (m *MockShiftRepository) FindByRoleType(ctx context.Context, roleType order.RoleType) ([]*order.Shift, error) {
	if m.findError != nil {
		return nil, m.findError
	}
	var shifts []*order.Shift
	for _, shift := range m.shifts {
		if shift.RoleType == roleType {
			shifts = append(shifts, shift)
		}
	}
	return shifts, nil
}

func (m *MockShiftRepository) FindAll(ctx context.Context) ([]*order.Shift, error) {
	if m.findError != nil {
		return nil, m.findError
	}
	var shifts []*order.Shift
	for _, shift := range m.shifts {
		shifts = append(shifts, shift)
	}
	return shifts, nil
}

// Mock OrderRepository
type MockOrderRepository struct{}

func (m *MockOrderRepository) Create(ctx context.Context, o *order.Order) error {
	return nil
}

func (m *MockOrderRepository) FindByID(ctx context.Context, id primitive.ObjectID) (*order.Order, error) {
	return nil, nil
}

func (m *MockOrderRepository) FindByOrderNumber(ctx context.Context, orderNumber string) (*order.Order, error) {
	return nil, errors.New("order not found")
}

func (m *MockOrderRepository) Update(ctx context.Context, id primitive.ObjectID, o *order.Order) error {
	return nil
}

func (m *MockOrderRepository) FindAll(ctx context.Context) ([]*order.Order, error) {
	return []*order.Order{}, nil
}

func (m *MockOrderRepository) FindByShiftID(ctx context.Context, shiftID primitive.ObjectID) ([]*order.Order, error) {
	return []*order.Order{}, nil
}

func (m *MockOrderRepository) FindByWaiterID(ctx context.Context, waiterID primitive.ObjectID) ([]*order.Order, error) {
	return []*order.Order{}, nil
}

func (m *MockOrderRepository) FindByStatus(ctx context.Context, status order.OrderStatus) ([]*order.Order, error) {
	return []*order.Order{}, nil
}

func (m *MockOrderRepository) FindByDateRange(ctx context.Context, startDate, endDate time.Time) ([]*order.Order, error) {
	return []*order.Order{}, nil
}

// Tests
func TestStartShift_WaiterRole(t *testing.T) {
	mockShiftRepo := NewMockShiftRepository()
	mockOrderRepo := &MockOrderRepository{}
	service := NewShiftService(mockShiftRepo, mockOrderRepo)

	userID := primitive.NewObjectID()
	req := &order.StartShiftRequest{
		Type:      order.ShiftMorning,
		StartCash: 1000000,
	}

	shift, err := service.StartShift(context.Background(), req, userID.Hex(), "Waiter 1", order.RoleWaiter)

	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}

	if shift.RoleType != order.RoleWaiter {
		t.Errorf("Expected role_type to be 'waiter', got %s", shift.RoleType)
	}

	if shift.UserID != userID {
		t.Errorf("Expected user_id to be %s, got %s", userID.Hex(), shift.UserID.Hex())
	}

	if shift.UserName != "Waiter 1" {
		t.Errorf("Expected user_name to be 'Waiter 1', got %s", shift.UserName)
	}

	if shift.Status != order.ShiftOpen {
		t.Errorf("Expected status to be OPEN, got %s", shift.Status)
	}

	if shift.StartCash != 1000000 {
		t.Errorf("Expected start_cash to be 1000000, got %f", shift.StartCash)
	}

	// Check legacy fields
	if shift.WaiterID != userID {
		t.Errorf("Expected waiter_id to be set for backward compatibility")
	}

	if shift.WaiterName != "Waiter 1" {
		t.Errorf("Expected waiter_name to be set for backward compatibility")
	}
}

func TestStartShift_BaristaRole(t *testing.T) {
	mockShiftRepo := NewMockShiftRepository()
	mockOrderRepo := &MockOrderRepository{}
	service := NewShiftService(mockShiftRepo, mockOrderRepo)

	userID := primitive.NewObjectID()
	req := &order.StartShiftRequest{
		Type:      order.ShiftMorning,
		StartCash: 0, // Barista doesn't handle cash
	}

	shift, err := service.StartShift(context.Background(), req, userID.Hex(), "Barista 1", order.RoleBarista)

	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}

	if shift.RoleType != order.RoleBarista {
		t.Errorf("Expected role_type to be 'barista', got %s", shift.RoleType)
	}

	if shift.UserID != userID {
		t.Errorf("Expected user_id to be %s, got %s", userID.Hex(), shift.UserID.Hex())
	}

	if shift.UserName != "Barista 1" {
		t.Errorf("Expected user_name to be 'Barista 1', got %s", shift.UserName)
	}

	// Barista should not have waiter_id set
	if shift.WaiterID != primitive.NilObjectID {
		t.Errorf("Expected waiter_id to be nil for barista role")
	}
}

func TestStartShift_CashierRole(t *testing.T) {
	mockShiftRepo := NewMockShiftRepository()
	mockOrderRepo := &MockOrderRepository{}
	service := NewShiftService(mockShiftRepo, mockOrderRepo)

	userID := primitive.NewObjectID()
	req := &order.StartShiftRequest{
		Type:      order.ShiftAfternoon,
		StartCash: 2000000,
	}

	shift, err := service.StartShift(context.Background(), req, userID.Hex(), "Cashier 1", order.RoleCashier)

	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}

	if shift.RoleType != order.RoleCashier {
		t.Errorf("Expected role_type to be 'cashier', got %s", shift.RoleType)
	}

	// Check legacy fields for cashier
	if shift.CashierID != userID {
		t.Errorf("Expected cashier_id to be set for backward compatibility")
	}

	if shift.CashierName != "Cashier 1" {
		t.Errorf("Expected cashier_name to be set for backward compatibility")
	}
}

func TestStartShift_DuplicateShiftSameRole(t *testing.T) {
	mockShiftRepo := NewMockShiftRepository()
	mockOrderRepo := &MockOrderRepository{}
	service := NewShiftService(mockShiftRepo, mockOrderRepo)

	userID := primitive.NewObjectID()
	req := &order.StartShiftRequest{
		Type:      order.ShiftMorning,
		StartCash: 1000000,
	}

	// Start first shift
	_, err := service.StartShift(context.Background(), req, userID.Hex(), "Waiter 1", order.RoleWaiter)
	if err != nil {
		t.Fatalf("Expected no error on first shift, got %v", err)
	}

	// Try to start second shift with same role
	_, err = service.StartShift(context.Background(), req, userID.Hex(), "Waiter 1", order.RoleWaiter)
	if err == nil {
		t.Error("Expected error when starting duplicate shift for same role, got nil")
	}

	if err.Error() != "user already has an open shift for this role" {
		t.Errorf("Expected specific error message, got %v", err)
	}
}

func TestStartShift_MultipleRolesSameUser(t *testing.T) {
	mockShiftRepo := NewMockShiftRepository()
	mockOrderRepo := &MockOrderRepository{}
	service := NewShiftService(mockShiftRepo, mockOrderRepo)

	userID := primitive.NewObjectID()

	// Start waiter shift
	waiterReq := &order.StartShiftRequest{
		Type:      order.ShiftMorning,
		StartCash: 1000000,
	}
	waiterShift, err := service.StartShift(context.Background(), waiterReq, userID.Hex(), "User 1", order.RoleWaiter)
	if err != nil {
		t.Fatalf("Expected no error on waiter shift, got %v", err)
	}

	// Start barista shift for same user
	baristaReq := &order.StartShiftRequest{
		Type:      order.ShiftMorning,
		StartCash: 0,
	}
	baristaShift, err := service.StartShift(context.Background(), baristaReq, userID.Hex(), "User 1", order.RoleBarista)
	if err != nil {
		t.Fatalf("Expected no error on barista shift, got %v", err)
	}

	// Both shifts should exist
	if waiterShift.RoleType != order.RoleWaiter {
		t.Error("Waiter shift should have waiter role")
	}

	if baristaShift.RoleType != order.RoleBarista {
		t.Error("Barista shift should have barista role")
	}

	// Both should be OPEN
	if waiterShift.Status != order.ShiftOpen {
		t.Error("Waiter shift should be OPEN")
	}

	if baristaShift.Status != order.ShiftOpen {
		t.Error("Barista shift should be OPEN")
	}
}

func TestGetCurrentShift_ByRole(t *testing.T) {
	mockShiftRepo := NewMockShiftRepository()
	mockOrderRepo := &MockOrderRepository{}
	service := NewShiftService(mockShiftRepo, mockOrderRepo)

	userID := primitive.NewObjectID()

	// Start waiter shift
	waiterReq := &order.StartShiftRequest{
		Type:      order.ShiftMorning,
		StartCash: 1000000,
	}
	_, err := service.StartShift(context.Background(), waiterReq, userID.Hex(), "User 1", order.RoleWaiter)
	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}

	// Start barista shift
	baristaReq := &order.StartShiftRequest{
		Type:      order.ShiftMorning,
		StartCash: 0,
	}
	_, err = service.StartShift(context.Background(), baristaReq, userID.Hex(), "User 1", order.RoleBarista)
	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}

	// Get waiter shift
	waiterShift, err := service.GetCurrentShift(context.Background(), userID, order.RoleWaiter)
	if err != nil {
		t.Fatalf("Expected to find waiter shift, got error: %v", err)
	}

	if waiterShift.RoleType != order.RoleWaiter {
		t.Errorf("Expected waiter role, got %s", waiterShift.RoleType)
	}

	// Get barista shift
	baristaShift, err := service.GetCurrentShift(context.Background(), userID, order.RoleBarista)
	if err != nil {
		t.Fatalf("Expected to find barista shift, got error: %v", err)
	}

	if baristaShift.RoleType != order.RoleBarista {
		t.Errorf("Expected barista role, got %s", baristaShift.RoleType)
	}

	// Try to get cashier shift (should not exist)
	_, err = service.GetCurrentShift(context.Background(), userID, order.RoleCashier)
	if err == nil {
		t.Error("Expected error when getting non-existent cashier shift")
	}
}

func TestGetShiftsByUser_FilteredByRole(t *testing.T) {
	mockShiftRepo := NewMockShiftRepository()
	mockOrderRepo := &MockOrderRepository{}
	service := NewShiftService(mockShiftRepo, mockOrderRepo)

	userID := primitive.NewObjectID()

	// Create multiple shifts for different roles
	waiterReq := &order.StartShiftRequest{Type: order.ShiftMorning, StartCash: 1000000}
	service.StartShift(context.Background(), waiterReq, userID.Hex(), "User 1", order.RoleWaiter)

	baristaReq := &order.StartShiftRequest{Type: order.ShiftAfternoon, StartCash: 0}
	service.StartShift(context.Background(), baristaReq, userID.Hex(), "User 1", order.RoleBarista)

	// Get waiter shifts only
	waiterShifts, err := service.GetShiftsByUser(context.Background(), userID, order.RoleWaiter)
	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}

	if len(waiterShifts) != 1 {
		t.Errorf("Expected 1 waiter shift, got %d", len(waiterShifts))
	}

	if waiterShifts[0].RoleType != order.RoleWaiter {
		t.Error("Expected waiter role in results")
	}

	// Get barista shifts only
	baristaShifts, err := service.GetShiftsByUser(context.Background(), userID, order.RoleBarista)
	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}

	if len(baristaShifts) != 1 {
		t.Errorf("Expected 1 barista shift, got %d", len(baristaShifts))
	}

	if baristaShifts[0].RoleType != order.RoleBarista {
		t.Error("Expected barista role in results")
	}
}

func TestGetShiftsByRole(t *testing.T) {
	mockShiftRepo := NewMockShiftRepository()
	mockOrderRepo := &MockOrderRepository{}
	service := NewShiftService(mockShiftRepo, mockOrderRepo)

	user1ID := primitive.NewObjectID()
	user2ID := primitive.NewObjectID()

	// Create shifts for different users and roles
	service.StartShift(context.Background(), &order.StartShiftRequest{Type: order.ShiftMorning, StartCash: 1000000}, user1ID.Hex(), "User 1", order.RoleWaiter)
	service.StartShift(context.Background(), &order.StartShiftRequest{Type: order.ShiftMorning, StartCash: 0}, user1ID.Hex(), "User 1", order.RoleBarista)
	service.StartShift(context.Background(), &order.StartShiftRequest{Type: order.ShiftAfternoon, StartCash: 0}, user2ID.Hex(), "User 2", order.RoleBarista)

	// Get all barista shifts
	baristaShifts, err := service.GetShiftsByRole(context.Background(), order.RoleBarista)
	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}

	if len(baristaShifts) != 2 {
		t.Errorf("Expected 2 barista shifts, got %d", len(baristaShifts))
	}

	for _, shift := range baristaShifts {
		if shift.RoleType != order.RoleBarista {
			t.Error("Expected all shifts to be barista role")
		}
	}

	// Get all waiter shifts
	waiterShifts, err := service.GetShiftsByRole(context.Background(), order.RoleWaiter)
	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}

	if len(waiterShifts) != 1 {
		t.Errorf("Expected 1 waiter shift, got %d", len(waiterShifts))
	}
}
