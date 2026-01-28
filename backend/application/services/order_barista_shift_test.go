package services

import (
	"context"
	"errors"
	"testing"
	"time"

	"cafe-pos/backend/domain/order"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

// TestAcceptOrder_BaristaWithoutShift tests BR-13: Barista must have open shift
func TestAcceptOrder_BaristaWithoutShift(t *testing.T) {
	mockOrderRepo := &MockOrderRepositoryForBarista{
		orders: make(map[string]*order.Order),
	}
	mockShiftRepo := NewMockShiftRepository()
	service := NewOrderService(mockOrderRepo, mockShiftRepo)

	// Create a QUEUED order
	orderID := primitive.NewObjectID()
	testOrder := &order.Order{
		ID:          orderID,
		OrderNumber: "ORD-001",
		Status:      order.StatusQueued,
		Items:       []order.OrderItem{{Name: "Coffee", Quantity: 1, Price: 50000}},
		Total:       50000,
		CreatedAt:   time.Now(),
	}
	mockOrderRepo.orders[orderID.Hex()] = testOrder

	baristaID := primitive.NewObjectID()

	// Try to accept order without opening shift
	_, err := service.AcceptOrder(context.Background(), orderID, baristaID.Hex(), "Barista 1")

	if err == nil {
		t.Error("Expected error when barista accepts order without open shift, got nil")
	}

	if err.Error() != "barista must open a shift before accepting orders" {
		t.Errorf("Expected specific error message, got: %v", err)
	}

	// Verify order status unchanged
	unchangedOrder, _ := mockOrderRepo.FindByID(context.Background(), orderID)
	if unchangedOrder.Status != order.StatusQueued {
		t.Errorf("Expected order to remain QUEUED, got %s", unchangedOrder.Status)
	}

	if unchangedOrder.BaristaID != primitive.NilObjectID {
		t.Error("Expected barista_id to remain unset")
	}
}

// TestAcceptOrder_BaristaWithOpenShift tests successful order acceptance with open shift
func TestAcceptOrder_BaristaWithOpenShift(t *testing.T) {
	mockOrderRepo := &MockOrderRepositoryForBarista{
		orders: make(map[string]*order.Order),
	}
	mockShiftRepo := NewMockShiftRepository()
	service := NewOrderService(mockOrderRepo, mockShiftRepo)

	// Create a QUEUED order
	orderID := primitive.NewObjectID()
	testOrder := &order.Order{
		ID:          orderID,
		OrderNumber: "ORD-001",
		Status:      order.StatusQueued,
		Items:       []order.OrderItem{{Name: "Coffee", Quantity: 1, Price: 50000}},
		Total:       50000,
		CreatedAt:   time.Now(),
	}
	mockOrderRepo.orders[orderID.Hex()] = testOrder

	baristaID := primitive.NewObjectID()

	// Open barista shift
	shiftReq := &order.StartShiftRequest{
		Type:      order.ShiftMorning,
		StartCash: 0,
	}
	shiftService := NewShiftService(mockShiftRepo, mockOrderRepo)
	_, err := shiftService.StartShift(context.Background(), shiftReq, baristaID.Hex(), "Barista 1", order.RoleBarista)
	if err != nil {
		t.Fatalf("Failed to start shift: %v", err)
	}

	// Now accept order (should succeed)
	acceptedOrder, err := service.AcceptOrder(context.Background(), orderID, baristaID.Hex(), "Barista 1")

	if err != nil {
		t.Fatalf("Expected no error when barista with open shift accepts order, got: %v", err)
	}

	if acceptedOrder.Status != order.StatusInProgress {
		t.Errorf("Expected order status to be IN_PROGRESS, got %s", acceptedOrder.Status)
	}

	if acceptedOrder.BaristaID != baristaID {
		t.Error("Expected barista_id to be set")
	}

	if acceptedOrder.BaristaName != "Barista 1" {
		t.Errorf("Expected barista_name to be 'Barista 1', got %s", acceptedOrder.BaristaName)
	}

	if acceptedOrder.AcceptedAt == nil {
		t.Error("Expected accepted_at to be set")
	}
}

// TestAcceptOrder_BaristaWithClosedShift tests rejection when shift is closed
func TestAcceptOrder_BaristaWithClosedShift(t *testing.T) {
	mockOrderRepo := &MockOrderRepositoryForBarista{
		orders: make(map[string]*order.Order),
	}
	mockShiftRepo := NewMockShiftRepository()
	service := NewOrderService(mockOrderRepo, mockShiftRepo)

	// Create a QUEUED order
	orderID := primitive.NewObjectID()
	testOrder := &order.Order{
		ID:          orderID,
		OrderNumber: "ORD-001",
		Status:      order.StatusQueued,
		Items:       []order.OrderItem{{Name: "Coffee", Quantity: 1, Price: 50000}},
		Total:       50000,
		CreatedAt:   time.Now(),
	}
	mockOrderRepo.orders[orderID.Hex()] = testOrder

	baristaID := primitive.NewObjectID()

	// Open and close shift
	shiftService := NewShiftService(mockShiftRepo, mockOrderRepo)
	shift, _ := shiftService.StartShift(context.Background(), &order.StartShiftRequest{
		Type:      order.ShiftMorning,
		StartCash: 0,
	}, baristaID.Hex(), "Barista 1", order.RoleBarista)

	// Close the shift
	endReq := &order.EndShiftRequest{EndCash: 0}
	_, err := shiftService.EndShift(context.Background(), shift.ID, endReq)
	if err != nil {
		t.Fatalf("Failed to end shift: %v", err)
	}

	// Try to accept order with closed shift
	_, err = service.AcceptOrder(context.Background(), orderID, baristaID.Hex(), "Barista 1")

	if err == nil {
		t.Error("Expected error when barista accepts order with closed shift, got nil")
	}

	if err.Error() != "barista must open a shift before accepting orders" {
		t.Errorf("Expected specific error message, got: %v", err)
	}
}

// TestAcceptOrder_MultipleBaristasDifferentShifts tests multiple baristas with their own shifts
func TestAcceptOrder_MultipleBaristasDifferentShifts(t *testing.T) {
	mockOrderRepo := &MockOrderRepositoryForBarista{
		orders: make(map[string]*order.Order),
	}
	mockShiftRepo := NewMockShiftRepository()
	orderService := NewOrderService(mockOrderRepo, mockShiftRepo)
	shiftService := NewShiftService(mockShiftRepo, mockOrderRepo)

	// Create two orders
	order1ID := primitive.NewObjectID()
	order2ID := primitive.NewObjectID()
	
	mockOrderRepo.orders[order1ID.Hex()] = &order.Order{
		ID:          order1ID,
		OrderNumber: "ORD-001",
		Status:      order.StatusQueued,
		Items:       []order.OrderItem{{Name: "Coffee", Quantity: 1, Price: 50000}},
		Total:       50000,
		CreatedAt:   time.Now(),
	}
	
	mockOrderRepo.orders[order2ID.Hex()] = &order.Order{
		ID:          order2ID,
		OrderNumber: "ORD-002",
		Status:      order.StatusQueued,
		Items:       []order.OrderItem{{Name: "Tea", Quantity: 1, Price: 40000}},
		Total:       40000,
		CreatedAt:   time.Now(),
	}

	barista1ID := primitive.NewObjectID()
	barista2ID := primitive.NewObjectID()

	// Barista 1 opens shift
	_, err := shiftService.StartShift(context.Background(), &order.StartShiftRequest{
		Type:      order.ShiftMorning,
		StartCash: 0,
	}, barista1ID.Hex(), "Barista 1", order.RoleBarista)
	if err != nil {
		t.Fatalf("Failed to start shift for barista 1: %v", err)
	}

	// Barista 1 accepts order 1 (should succeed)
	acceptedOrder1, err := orderService.AcceptOrder(context.Background(), order1ID, barista1ID.Hex(), "Barista 1")
	if err != nil {
		t.Fatalf("Expected barista 1 to accept order, got error: %v", err)
	}
	if acceptedOrder1.BaristaID != barista1ID {
		t.Error("Expected order 1 to be assigned to barista 1")
	}

	// Barista 2 tries to accept order 2 without shift (should fail)
	_, err = orderService.AcceptOrder(context.Background(), order2ID, barista2ID.Hex(), "Barista 2")
	if err == nil {
		t.Error("Expected barista 2 to fail without open shift")
	}

	// Barista 2 opens shift
	_, err = shiftService.StartShift(context.Background(), &order.StartShiftRequest{
		Type:      order.ShiftMorning,
		StartCash: 0,
	}, barista2ID.Hex(), "Barista 2", order.RoleBarista)
	if err != nil {
		t.Fatalf("Failed to start shift for barista 2: %v", err)
	}

	// Barista 2 accepts order 2 (should succeed)
	acceptedOrder2, err := orderService.AcceptOrder(context.Background(), order2ID, barista2ID.Hex(), "Barista 2")
	if err != nil {
		t.Fatalf("Expected barista 2 to accept order, got error: %v", err)
	}
	if acceptedOrder2.BaristaID != barista2ID {
		t.Error("Expected order 2 to be assigned to barista 2")
	}
}

// MockOrderRepositoryForBarista - Extended mock for barista tests
type MockOrderRepositoryForBarista struct {
	orders map[string]*order.Order
}

func (m *MockOrderRepositoryForBarista) Create(ctx context.Context, o *order.Order) error {
	o.ID = primitive.NewObjectID()
	m.orders[o.ID.Hex()] = o
	return nil
}

func (m *MockOrderRepositoryForBarista) FindByID(ctx context.Context, id primitive.ObjectID) (*order.Order, error) {
	if o, exists := m.orders[id.Hex()]; exists {
		return o, nil
	}
	return nil, errors.New("order not found")
}

func (m *MockOrderRepositoryForBarista) Update(ctx context.Context, id primitive.ObjectID, o *order.Order) error {
	if _, exists := m.orders[id.Hex()]; exists {
		m.orders[id.Hex()] = o
		return nil
	}
	return errors.New("order not found")
}

func (m *MockOrderRepositoryForBarista) FindByShiftID(ctx context.Context, shiftID primitive.ObjectID) ([]*order.Order, error) {
	return []*order.Order{}, nil
}

func (m *MockOrderRepositoryForBarista) FindByWaiterID(ctx context.Context, waiterID primitive.ObjectID) ([]*order.Order, error) {
	return []*order.Order{}, nil
}

func (m *MockOrderRepositoryForBarista) FindByStatus(ctx context.Context, status order.OrderStatus) ([]*order.Order, error) {
	var orders []*order.Order
	for _, o := range m.orders {
		if o.Status == status {
			orders = append(orders, o)
		}
	}
	return orders, nil
}

func (m *MockOrderRepositoryForBarista) FindByOrderNumber(ctx context.Context, orderNumber string) (*order.Order, error) {
	for _, o := range m.orders {
		if o.OrderNumber == orderNumber {
			return o, nil
		}
	}
	return nil, errors.New("order not found")
}

func (m *MockOrderRepositoryForBarista) FindAll(ctx context.Context) ([]*order.Order, error) {
	var orders []*order.Order
	for _, o := range m.orders {
		orders = append(orders, o)
	}
	return orders, nil
}
