package services

import (
	"context"
	"errors"
	"testing"
	"time"

	"cafe-pos/backend/domain"
	"cafe-pos/backend/domain/menu"
	"cafe-pos/backend/domain/order"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

// Mock repositories for OrderService testing

type mockOrderRepositoryForVariants struct {
	orders map[primitive.ObjectID]*order.Order
}

func newMockOrderRepositoryForVariants() *mockOrderRepositoryForVariants {
	return &mockOrderRepositoryForVariants{
		orders: make(map[primitive.ObjectID]*order.Order),
	}
}

func (m *mockOrderRepositoryForVariants) Create(ctx context.Context, o *order.Order) error {
	if o.ID.IsZero() {
		o.ID = primitive.NewObjectID()
	}
	o.CreatedAt = time.Now()
	o.UpdatedAt = time.Now()
	m.orders[o.ID] = o
	return nil
}

func (m *mockOrderRepositoryForVariants) FindByID(ctx context.Context, id primitive.ObjectID) (*order.Order, error) {
	o, exists := m.orders[id]
	if !exists {
		return nil, errors.New("order not found")
	}
	return o, nil
}

func (m *mockOrderRepositoryForVariants) Update(ctx context.Context, id primitive.ObjectID, o *order.Order) error {
	if _, exists := m.orders[id]; !exists {
		return errors.New("order not found")
	}
	o.UpdatedAt = time.Now()
	m.orders[id] = o
	return nil
}

func (m *mockOrderRepositoryForVariants) FindByShiftID(ctx context.Context, shiftID primitive.ObjectID) ([]*order.Order, error) {
	var orders []*order.Order
	for _, o := range m.orders {
		if o.ShiftID == shiftID {
			orders = append(orders, o)
		}
	}
	return orders, nil
}

func (m *mockOrderRepositoryForVariants) FindByWaiterID(ctx context.Context, waiterID primitive.ObjectID) ([]*order.Order, error) {
	return nil, nil
}

func (m *mockOrderRepositoryForVariants) FindByStatus(ctx context.Context, status order.OrderStatus) ([]*order.Order, error) {
	return nil, nil
}

func (m *mockOrderRepositoryForVariants) FindByOrderNumber(ctx context.Context, orderNumber string) (*order.Order, error) {
	return nil, nil
}

func (m *mockOrderRepositoryForVariants) FindAll(ctx context.Context) ([]*order.Order, error) {
	return nil, nil
}

type mockShiftRepositoryForVariants struct {
	shifts map[primitive.ObjectID]*order.Shift
}

func newMockShiftRepositoryForVariants() *mockShiftRepositoryForVariants {
	return &mockShiftRepositoryForVariants{
		shifts: make(map[primitive.ObjectID]*order.Shift),
	}
}

func (m *mockShiftRepositoryForVariants) Create(ctx context.Context, shift *order.Shift) error {
	if shift.ID.IsZero() {
		shift.ID = primitive.NewObjectID()
	}
	m.shifts[shift.ID] = shift
	return nil
}

func (m *mockShiftRepositoryForVariants) FindByID(ctx context.Context, id primitive.ObjectID) (*order.Shift, error) {
	shift, exists := m.shifts[id]
	if !exists {
		return nil, errors.New("shift not found")
	}
	return shift, nil
}

func (m *mockShiftRepositoryForVariants) Update(ctx context.Context, id primitive.ObjectID, shift *order.Shift) error {
	if _, exists := m.shifts[id]; !exists {
		return errors.New("shift not found")
	}
	m.shifts[id] = shift
	return nil
}

func (m *mockShiftRepositoryForVariants) FindOpenShiftByUser(ctx context.Context, userID primitive.ObjectID, role order.RoleType) (*order.Shift, error) {
	return nil, nil
}

func (m *mockShiftRepositoryForVariants) FindOpenShiftByWaiter(ctx context.Context, waiterID primitive.ObjectID) (*order.Shift, error) {
	return nil, nil
}

func (m *mockShiftRepositoryForVariants) FindOpenShifts(ctx context.Context) ([]*order.Shift, error) {
	return nil, nil
}

func (m *mockShiftRepositoryForVariants) FindByWaiterID(ctx context.Context, waiterID primitive.ObjectID) ([]*order.Shift, error) {
	return nil, nil
}

func (m *mockShiftRepositoryForVariants) FindByUserID(ctx context.Context, userID primitive.ObjectID, roleType order.RoleType) ([]*order.Shift, error) {
	return nil, nil
}

func (m *mockShiftRepositoryForVariants) FindByUserAndDateRange(ctx context.Context, userID primitive.ObjectID, startDate, endDate time.Time) ([]*order.Shift, error) {
	return nil, nil
}

func (m *mockShiftRepositoryForVariants) FindByDateRange(ctx context.Context, startDate, endDate time.Time) ([]*order.Shift, error) {
	return nil, nil
}

func (m *mockShiftRepositoryForVariants) FindByRoleType(ctx context.Context, roleType order.RoleType) ([]*order.Shift, error) {
	return nil, nil
}

func (m *mockShiftRepositoryForVariants) FindAll(ctx context.Context) ([]*order.Shift, error) {
	return nil, nil
}

type mockMenuRepositoryForVariants struct {
	items map[primitive.ObjectID]*menu.MenuItem
}

func newMockMenuRepositoryForVariants() *mockMenuRepositoryForVariants {
	return &mockMenuRepositoryForVariants{
		items: make(map[primitive.ObjectID]*menu.MenuItem),
	}
}

func (m *mockMenuRepositoryForVariants) Create(ctx context.Context, item *menu.MenuItem) error {
	if item.ID.IsZero() {
		item.ID = primitive.NewObjectID()
	}
	m.items[item.ID] = item
	return nil
}

func (m *mockMenuRepositoryForVariants) FindAll(ctx context.Context) ([]*menu.MenuItem, error) {
	var items []*menu.MenuItem
	for _, item := range m.items {
		items = append(items, item)
	}
	return items, nil
}

func (m *mockMenuRepositoryForVariants) FindByID(ctx context.Context, id primitive.ObjectID) (*menu.MenuItem, error) {
	item, exists := m.items[id]
	if !exists {
		return nil, errors.New("menu item not found")
	}
	return item, nil
}

func (m *mockMenuRepositoryForVariants) FindByCategory(ctx context.Context, category string) ([]*menu.MenuItem, error) {
	return nil, nil
}

func (m *mockMenuRepositoryForVariants) FindByIngredientName(ctx context.Context, ingredientName string) ([]*menu.MenuItem, error) {
	return nil, nil
}

func (m *mockMenuRepositoryForVariants) Update(ctx context.Context, id primitive.ObjectID, item *menu.MenuItem) error {
	if _, exists := m.items[id]; !exists {
		return errors.New("menu item not found")
	}
	m.items[id] = item
	return nil
}

func (m *mockMenuRepositoryForVariants) Delete(ctx context.Context, id primitive.ObjectID) error {
	return nil
}

// Test backward compatibility - single-size items

func TestCreateOrder_SingleSize(t *testing.T) {
	orderRepo := newMockOrderRepositoryForVariants()
	shiftRepo := newMockShiftRepositoryForVariants()
	menuRepo := newMockMenuRepositoryForVariants()
	smManager := domain.NewStateMachineManager()
	
	service := NewOrderService(orderRepo, shiftRepo, menuRepo, smManager)
	ctx := context.Background()

	// Create open shift
	shift := &order.Shift{
		ID:     primitive.NewObjectID(),
		Status: order.ShiftOpen,
	}
	shiftRepo.Create(ctx, shift)

	// Create single-size menu item
	menuItem := &menu.MenuItem{
		ID:          primitive.NewObjectID(),
		Name:        "Bánh mì",
		Category:    "Món ăn",
		HasVariants: false,
		Price:       20000,
		Available:   true,
	}
	menuRepo.Create(ctx, menuItem)

	// Create order request
	req := &order.CreateOrderRequest{
		CustomerName: "Khách 1",
		ShiftID:      shift.ID.Hex(),
		Items: []order.OrderItem{
			{
				MenuItemID: menuItem.ID,
				Quantity:   2,
			},
		},
	}

	waiterID := primitive.NewObjectID().Hex()
	o, err := service.CreateOrder(ctx, req, waiterID, "Waiter 1")
	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}

	if len(o.Items) != 1 {
		t.Fatalf("Expected 1 item, got %d", len(o.Items))
	}

	item := o.Items[0]
	if item.Name != "Bánh mì" {
		t.Errorf("Expected item name 'Bánh mì', got '%s'", item.Name)
	}
	if item.Price != 20000 {
		t.Errorf("Expected price 20000, got %f", item.Price)
	}
	if item.VariantID != "" {
		t.Errorf("Expected empty variant_id for single-size, got '%s'", item.VariantID)
	}
	if item.VariantName != "" {
		t.Errorf("Expected empty variant_name for single-size, got '%s'", item.VariantName)
	}
	if item.Subtotal != 40000 {
		t.Errorf("Expected subtotal 40000, got %f", item.Subtotal)
	}
	if o.Total != 40000 {
		t.Errorf("Expected total 40000, got %f", o.Total)
	}
}

// Test new functionality - multi-size items

func TestCreateOrder_MultiSize_WithVariantID(t *testing.T) {
	orderRepo := newMockOrderRepositoryForVariants()
	shiftRepo := newMockShiftRepositoryForVariants()
	menuRepo := newMockMenuRepositoryForVariants()
	smManager := domain.NewStateMachineManager()
	
	service := NewOrderService(orderRepo, shiftRepo, menuRepo, smManager)
	ctx := context.Background()

	// Create open shift
	shift := &order.Shift{
		ID:     primitive.NewObjectID(),
		Status: order.ShiftOpen,
	}
	shiftRepo.Create(ctx, shift)

	// Create multi-size menu item
	menuItem := &menu.MenuItem{
		ID:          primitive.NewObjectID(),
		Name:        "Cà phê sữa đá",
		Category:    "Cà phê",
		HasVariants: true,
		Available:   true,
		Variants: []menu.MenuItemVariant{
			{
				ID:        "M",
				Name:      "Size M",
				Price:     25000,
				Available: true,
				IsDefault: true,
			},
			{
				ID:        "L",
				Name:      "Size L",
				Price:     30000,
				Available: true,
				IsDefault: false,
			},
		},
	}
	menuRepo.Create(ctx, menuItem)

	// Create order request with variant_id
	req := &order.CreateOrderRequest{
		CustomerName: "Khách 1",
		ShiftID:      shift.ID.Hex(),
		Items: []order.OrderItem{
			{
				MenuItemID: menuItem.ID,
				VariantID:  "L", // Order Size L
				Quantity:   2,
			},
		},
	}

	waiterID := primitive.NewObjectID().Hex()
	o, err := service.CreateOrder(ctx, req, waiterID, "Waiter 1")
	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}

	if len(o.Items) != 1 {
		t.Fatalf("Expected 1 item, got %d", len(o.Items))
	}

	item := o.Items[0]
	if item.Name != "Cà phê sữa đá" {
		t.Errorf("Expected item name 'Cà phê sữa đá', got '%s'", item.Name)
	}
	if item.VariantID != "L" {
		t.Errorf("Expected variant_id 'L', got '%s'", item.VariantID)
	}
	if item.VariantName != "Size L" {
		t.Errorf("Expected variant_name 'Size L', got '%s'", item.VariantName)
	}
	if item.Price != 30000 {
		t.Errorf("Expected price 30000 (Size L), got %f", item.Price)
	}
	if item.Subtotal != 60000 {
		t.Errorf("Expected subtotal 60000, got %f", item.Subtotal)
	}
	if o.Total != 60000 {
		t.Errorf("Expected total 60000, got %f", o.Total)
	}
}

func TestCreateOrder_MultiSize_MissingVariantID(t *testing.T) {
	orderRepo := newMockOrderRepositoryForVariants()
	shiftRepo := newMockShiftRepositoryForVariants()
	menuRepo := newMockMenuRepositoryForVariants()
	smManager := domain.NewStateMachineManager()
	
	service := NewOrderService(orderRepo, shiftRepo, menuRepo, smManager)
	ctx := context.Background()

	// Create open shift
	shift := &order.Shift{
		ID:     primitive.NewObjectID(),
		Status: order.ShiftOpen,
	}
	shiftRepo.Create(ctx, shift)

	// Create multi-size menu item
	menuItem := &menu.MenuItem{
		ID:          primitive.NewObjectID(),
		Name:        "Cà phê sữa đá",
		Category:    "Cà phê",
		HasVariants: true,
		Available:   true,
		Variants: []menu.MenuItemVariant{
			{ID: "M", Name: "Size M", Price: 25000, Available: true, IsDefault: true},
		},
	}
	menuRepo.Create(ctx, menuItem)

	// Create order request WITHOUT variant_id
	req := &order.CreateOrderRequest{
		CustomerName: "Khách 1",
		ShiftID:      shift.ID.Hex(),
		Items: []order.OrderItem{
			{
				MenuItemID: menuItem.ID,
				// Missing VariantID
				Quantity: 1,
			},
		},
	}

	waiterID := primitive.NewObjectID().Hex()
	_, err := service.CreateOrder(ctx, req, waiterID, "Waiter 1")
	if err == nil {
		t.Error("Expected error for multi-size item without variant_id")
	}
}

func TestCreateOrder_MultiSize_InvalidVariantID(t *testing.T) {
	orderRepo := newMockOrderRepositoryForVariants()
	shiftRepo := newMockShiftRepositoryForVariants()
	menuRepo := newMockMenuRepositoryForVariants()
	smManager := domain.NewStateMachineManager()
	
	service := NewOrderService(orderRepo, shiftRepo, menuRepo, smManager)
	ctx := context.Background()

	// Create open shift
	shift := &order.Shift{
		ID:     primitive.NewObjectID(),
		Status: order.ShiftOpen,
	}
	shiftRepo.Create(ctx, shift)

	// Create multi-size menu item
	menuItem := &menu.MenuItem{
		ID:          primitive.NewObjectID(),
		Name:        "Cà phê sữa đá",
		Category:    "Cà phê",
		HasVariants: true,
		Available:   true,
		Variants: []menu.MenuItemVariant{
			{ID: "M", Name: "Size M", Price: 25000, Available: true, IsDefault: true},
			{ID: "L", Name: "Size L", Price: 30000, Available: true, IsDefault: false},
		},
	}
	menuRepo.Create(ctx, menuItem)

	// Create order request with INVALID variant_id
	req := &order.CreateOrderRequest{
		CustomerName: "Khách 1",
		ShiftID:      shift.ID.Hex(),
		Items: []order.OrderItem{
			{
				MenuItemID: menuItem.ID,
				VariantID:  "XL", // Invalid - doesn't exist
				Quantity:   1,
			},
		},
	}

	waiterID := primitive.NewObjectID().Hex()
	_, err := service.CreateOrder(ctx, req, waiterID, "Waiter 1")
	if err == nil {
		t.Error("Expected error for invalid variant_id")
	}
}

// Test mixed orders

func TestCreateOrder_MixedSingleAndMultiSize(t *testing.T) {
	orderRepo := newMockOrderRepositoryForVariants()
	shiftRepo := newMockShiftRepositoryForVariants()
	menuRepo := newMockMenuRepositoryForVariants()
	smManager := domain.NewStateMachineManager()
	
	service := NewOrderService(orderRepo, shiftRepo, menuRepo, smManager)
	ctx := context.Background()

	// Create open shift
	shift := &order.Shift{
		ID:     primitive.NewObjectID(),
		Status: order.ShiftOpen,
	}
	shiftRepo.Create(ctx, shift)

	// Create single-size menu item
	singleItem := &menu.MenuItem{
		ID:          primitive.NewObjectID(),
		Name:        "Bánh mì",
		Category:    "Món ăn",
		HasVariants: false,
		Price:       20000,
		Available:   true,
	}
	menuRepo.Create(ctx, singleItem)

	// Create multi-size menu item
	multiItem := &menu.MenuItem{
		ID:          primitive.NewObjectID(),
		Name:        "Cà phê sữa đá",
		Category:    "Cà phê",
		HasVariants: true,
		Available:   true,
		Variants: []menu.MenuItemVariant{
			{ID: "M", Name: "Size M", Price: 25000, Available: true, IsDefault: true},
			{ID: "L", Name: "Size L", Price: 30000, Available: true, IsDefault: false},
		},
	}
	menuRepo.Create(ctx, multiItem)

	// Create order with both types
	req := &order.CreateOrderRequest{
		CustomerName: "Khách 1",
		ShiftID:      shift.ID.Hex(),
		Items: []order.OrderItem{
			{
				MenuItemID: singleItem.ID,
				Quantity:   1,
			},
			{
				MenuItemID: multiItem.ID,
				VariantID:  "L",
				Quantity:   2,
			},
		},
	}

	waiterID := primitive.NewObjectID().Hex()
	o, err := service.CreateOrder(ctx, req, waiterID, "Waiter 1")
	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}

	if len(o.Items) != 2 {
		t.Fatalf("Expected 2 items, got %d", len(o.Items))
	}

	// Check single-size item
	item1 := o.Items[0]
	if item1.Name != "Bánh mì" {
		t.Errorf("Expected item 1 name 'Bánh mì', got '%s'", item1.Name)
	}
	if item1.Price != 20000 {
		t.Errorf("Expected item 1 price 20000, got %f", item1.Price)
	}
	if item1.VariantID != "" {
		t.Errorf("Expected item 1 empty variant_id, got '%s'", item1.VariantID)
	}

	// Check multi-size item
	item2 := o.Items[1]
	if item2.Name != "Cà phê sữa đá" {
		t.Errorf("Expected item 2 name 'Cà phê sữa đá', got '%s'", item2.Name)
	}
	if item2.VariantID != "L" {
		t.Errorf("Expected item 2 variant_id 'L', got '%s'", item2.VariantID)
	}
	if item2.VariantName != "Size L" {
		t.Errorf("Expected item 2 variant_name 'Size L', got '%s'", item2.VariantName)
	}
	if item2.Price != 30000 {
		t.Errorf("Expected item 2 price 30000, got %f", item2.Price)
	}

	// Check total
	expectedTotal := 20000 + (30000 * 2)
	if o.Total != float64(expectedTotal) {
		t.Errorf("Expected total %d, got %f", expectedTotal, o.Total)
	}
}

// Test EditOrder with variants

func TestEditOrder_WithVariants(t *testing.T) {
	orderRepo := newMockOrderRepositoryForVariants()
	shiftRepo := newMockShiftRepositoryForVariants()
	menuRepo := newMockMenuRepositoryForVariants()
	smManager := domain.NewStateMachineManager()
	
	service := NewOrderService(orderRepo, shiftRepo, menuRepo, smManager)
	ctx := context.Background()

	// Create open shift
	shift := &order.Shift{
		ID:     primitive.NewObjectID(),
		Status: order.ShiftOpen,
	}
	shiftRepo.Create(ctx, shift)

	// Create multi-size menu item
	menuItem := &menu.MenuItem{
		ID:          primitive.NewObjectID(),
		Name:        "Cà phê sữa đá",
		Category:    "Cà phê",
		HasVariants: true,
		Available:   true,
		Variants: []menu.MenuItemVariant{
			{ID: "M", Name: "Size M", Price: 25000, Available: true, IsDefault: true},
			{ID: "L", Name: "Size L", Price: 30000, Available: true, IsDefault: false},
		},
	}
	menuRepo.Create(ctx, menuItem)

	// Create initial order with Size M
	createReq := &order.CreateOrderRequest{
		CustomerName: "Khách 1",
		ShiftID:      shift.ID.Hex(),
		Items: []order.OrderItem{
			{
				MenuItemID: menuItem.ID,
				VariantID:  "M",
				Quantity:   1,
			},
		},
	}
	waiterID := primitive.NewObjectID().Hex()
	o, _ := service.CreateOrder(ctx, createReq, waiterID, "Waiter 1")

	// Edit order to change to Size L
	editReq := &order.EditOrderRequest{
		Items: []order.OrderItem{
			{
				MenuItemID: menuItem.ID,
				VariantID:  "L", // Change to Size L
				Quantity:   1,
			},
		},
	}

	response, err := service.EditOrder(ctx, o.ID, editReq)
	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}

	edited := response.Order
	if len(edited.Items) != 1 {
		t.Fatalf("Expected 1 item, got %d", len(edited.Items))
	}

	item := edited.Items[0]
	if item.VariantID != "L" {
		t.Errorf("Expected variant_id 'L', got '%s'", item.VariantID)
	}
	if item.VariantName != "Size L" {
		t.Errorf("Expected variant_name 'Size L', got '%s'", item.VariantName)
	}
	if item.Price != 30000 {
		t.Errorf("Expected price 30000, got %f", item.Price)
	}
}
