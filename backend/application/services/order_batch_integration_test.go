package services

import (
	"context"
	"testing"
	"time"

	"cafe-pos/backend/domain/batch"
	"cafe-pos/backend/domain/menu"
	"cafe-pos/backend/domain/order"
	"github.com/stretchr/testify/assert"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

// Test: Order creation with batch ingredients should deduct batch quantities
func TestOrderService_CreateOrder_WithBatchIngredients_DeductsBatches(t *testing.T) {
	// Setup
	ctx := context.Background()
	
	// Create mock repositories
	orderRepo := newMockOrderRepoForBatchIntegration()
	shiftRepo := newMockShiftRepoForBatchIntegration()
	menuRepo := newMockMenuRepoForBatchIntegration()
	batchRecordRepo := newMockBatchRecordRepoForBatchIntegration()
	batchUsageLogRepo := newMockBatchUsageLogRepoForBatchIntegration()
	
	// Create services
	batchUsageService := NewBatchUsageService(batchRecordRepo, batchUsageLogRepo)
	orderService := NewOrderService(orderRepo, shiftRepo, menuRepo, nil, batchUsageService)
	
	// Create test data
	// 1. Create a batch record with 1000ml available
	batchID := primitive.NewObjectID()
	batchDefID := primitive.NewObjectID()
	batchRecord := &batch.BatchRecord{
		ID:                batchID,
		BatchDefinitionID: batchDefID,
		BatchName:         "Coffee Concentrate",
		QuantityProduced:  1000,
		QuantityRemaining: 1000,
		Unit:              "ml",
		CostPerUnit:       0.15,
		TotalCost:         150.0,
		Status:            batch.BatchStatusAvailable,
		ExpiresAt:         time.Now().Add(24 * time.Hour),
	}
	batchRecordRepo.batches[batchID] = batchRecord
	batchRecordRepo.batchesByDef[batchDefID] = []*batch.BatchRecord{batchRecord}
	
	// 2. Create a menu item that uses batch ingredient
	menuItemID := primitive.NewObjectID()
	menuItem := &menu.MenuItem{
		ID:          menuItemID,
		Name:        "Iced Coffee",
		Category:    "Coffee",
		Price:       50000,
		Available:   true,
		HasVariants: false,
		Ingredients: []menu.Ingredient{
			{
				Name:           "Coffee Concentrate",
				Quantity:       30,
				Unit:           "ml",
				IngredientType: string(menu.IngredientTypeBatch),
				BatchID:        &batchDefID,
			},
		},
	}
	menuRepo.items[menuItemID] = menuItem
	
	// 3. Create an open shift
	shiftID := primitive.NewObjectID()
	shift := &order.Shift{
		ID:     shiftID,
		Status: order.ShiftOpen,
	}
	shiftRepo.shifts[shiftID] = shift
	
	// Create order request
	req := &order.CreateOrderRequest{
		ShiftID:      shiftID.Hex(),
		CustomerName: "Test Customer",
		Items: []order.OrderItem{
			{
				MenuItemID: menuItemID,
				Quantity:   2, // 2 drinks = 60ml total
			},
		},
	}
	
	// Execute
	createdOrder, err := orderService.CreateOrder(ctx, req, "waiter123", "Test Waiter")
	
	// Assert
	assert.NoError(t, err)
	assert.NotNil(t, createdOrder)
	assert.Equal(t, 1, len(createdOrder.Items)) // 1 item with quantity 2
	assert.Equal(t, 2, createdOrder.Items[0].Quantity)
	
	// Verify batch was deducted
	assert.Equal(t, 940.0, batchRecord.QuantityRemaining) // 1000 - 60 = 940
	
	// Verify usage log was created
	assert.Equal(t, 1, len(batchUsageLogRepo.logs))
	usageLog := batchUsageLogRepo.logs[0]
	assert.Equal(t, batchID, usageLog.BatchRecordID)
	assert.Equal(t, 60.0, usageLog.QuantityUsed)
	assert.Equal(t, 0.15, usageLog.CostPerUnit)
	assert.Equal(t, 9.0, usageLog.TotalCost) // 60 * 0.15 = 9.0
	
	// Verify batch cost is tracked in order note
	assert.Contains(t, createdOrder.Note, "Batch Cost: 9.00 VND")
}

// Test: Order creation with insufficient batch should fail and rollback
func TestOrderService_CreateOrder_InsufficientBatch_RollsBack(t *testing.T) {
	// Setup
	ctx := context.Background()
	
	orderRepo := newMockOrderRepoForBatchIntegration()
	shiftRepo := newMockShiftRepoForBatchIntegration()
	menuRepo := newMockMenuRepoForBatchIntegration()
	batchRecordRepo := newMockBatchRecordRepoForBatchIntegration()
	batchUsageLogRepo := newMockBatchUsageLogRepoForBatchIntegration()
	
	batchUsageService := NewBatchUsageService(batchRecordRepo, batchUsageLogRepo)
	orderService := NewOrderService(orderRepo, shiftRepo, menuRepo, nil, batchUsageService)
	
	// Create batch with only 20ml available
	batchID := primitive.NewObjectID()
	batchDefID := primitive.NewObjectID()
	batchRecord := &batch.BatchRecord{
		ID:                batchID,
		BatchDefinitionID: batchDefID,
		BatchName:         "Coffee Concentrate",
		QuantityProduced:  1000,
		QuantityRemaining: 20, // Only 20ml available
		Unit:              "ml",
		CostPerUnit:       0.15,
		Status:            batch.BatchStatusAvailable,
		ExpiresAt:         time.Now().Add(24 * time.Hour),
	}
	batchRecordRepo.batches[batchID] = batchRecord
	batchRecordRepo.batchesByDef[batchDefID] = []*batch.BatchRecord{batchRecord}
	
	// Create menu item that needs 30ml per drink
	menuItemID := primitive.NewObjectID()
	menuItem := &menu.MenuItem{
		ID:          menuItemID,
		Name:        "Iced Coffee",
		Category:    "Coffee",
		Price:       50000,
		Available:   true,
		HasVariants: false,
		Ingredients: []menu.Ingredient{
			{
				Name:           "Coffee Concentrate",
				Quantity:       30,
				Unit:           "ml",
				IngredientType: string(menu.IngredientTypeBatch),
				BatchID:        &batchDefID,
			},
		},
	}
	menuRepo.items[menuItemID] = menuItem
	
	// Create open shift
	shiftID := primitive.NewObjectID()
	shift := &order.Shift{
		ID:     shiftID,
		Status: order.ShiftOpen,
	}
	shiftRepo.shifts[shiftID] = shift
	
	// Try to create order for 2 drinks (needs 60ml, but only 20ml available)
	req := &order.CreateOrderRequest{
		ShiftID:      shiftID.Hex(),
		CustomerName: "Test Customer",
		Items: []order.OrderItem{
			{
				MenuItemID: menuItemID,
				Quantity:   2,
			},
		},
	}
	
	// Execute
	createdOrder, err := orderService.CreateOrder(ctx, req, "waiter123", "Test Waiter")
	
	// Assert
	assert.Error(t, err)
	assert.Nil(t, createdOrder)
	assert.Contains(t, err.Error(), "insufficient batch")
	
	// Verify batch quantity was not changed (rollback worked)
	assert.Equal(t, 20.0, batchRecord.QuantityRemaining)
	
	// Verify no usage logs were created
	assert.Equal(t, 0, len(batchUsageLogRepo.logs))
	
	// Verify no order was created
	assert.Equal(t, 0, len(orderRepo.orders))
}

// Test: Order with multiple batch ingredients deducts all correctly
func TestOrderService_CreateOrder_MultipleBatchIngredients_DeductsAll(t *testing.T) {
	// Setup
	ctx := context.Background()
	
	orderRepo := newMockOrderRepoForBatchIntegration()
	shiftRepo := newMockShiftRepoForBatchIntegration()
	menuRepo := newMockMenuRepoForBatchIntegration()
	batchRecordRepo := newMockBatchRecordRepoForBatchIntegration()
	batchUsageLogRepo := newMockBatchUsageLogRepoForBatchIntegration()
	
	batchUsageService := NewBatchUsageService(batchRecordRepo, batchUsageLogRepo)
	orderService := NewOrderService(orderRepo, shiftRepo, menuRepo, nil, batchUsageService)
	
	// Create two different batches
	coffeeBatchID := primitive.NewObjectID()
	coffeeDefID := primitive.NewObjectID()
	coffeeBatch := &batch.BatchRecord{
		ID:                coffeeBatchID,
		BatchDefinitionID: coffeeDefID,
		BatchName:         "Coffee Concentrate",
		QuantityProduced:  1000,
		QuantityRemaining: 1000,
		Unit:              "ml",
		CostPerUnit:       0.15,
		Status:            batch.BatchStatusAvailable,
		ExpiresAt:         time.Now().Add(24 * time.Hour),
	}
	
	milkBatchID := primitive.NewObjectID()
	milkDefID := primitive.NewObjectID()
	milkBatch := &batch.BatchRecord{
		ID:                milkBatchID,
		BatchDefinitionID: milkDefID,
		BatchName:         "Steamed Milk",
		QuantityProduced:  2000,
		QuantityRemaining: 2000,
		Unit:              "ml",
		CostPerUnit:       0.05,
		Status:            batch.BatchStatusAvailable,
		ExpiresAt:         time.Now().Add(12 * time.Hour),
	}
	
	batchRecordRepo.batches[coffeeBatchID] = coffeeBatch
	batchRecordRepo.batches[milkBatchID] = milkBatch
	batchRecordRepo.batchesByDef[coffeeDefID] = []*batch.BatchRecord{coffeeBatch}
	batchRecordRepo.batchesByDef[milkDefID] = []*batch.BatchRecord{milkBatch}
	
	// Create menu item with two batch ingredients
	menuItemID := primitive.NewObjectID()
	menuItem := &menu.MenuItem{
		ID:          menuItemID,
		Name:        "Latte",
		Category:    "Coffee",
		Price:       60000,
		Available:   true,
		HasVariants: false,
		Ingredients: []menu.Ingredient{
			{
				Name:           "Coffee Concentrate",
				Quantity:       30,
				Unit:           "ml",
				IngredientType: string(menu.IngredientTypeBatch),
				BatchID:        &coffeeDefID,
			},
			{
				Name:           "Steamed Milk",
				Quantity:       150,
				Unit:           "ml",
				IngredientType: string(menu.IngredientTypeBatch),
				BatchID:        &milkDefID,
			},
		},
	}
	menuRepo.items[menuItemID] = menuItem
	
	// Create open shift
	shiftID := primitive.NewObjectID()
	shift := &order.Shift{
		ID:     shiftID,
		Status: order.ShiftOpen,
	}
	shiftRepo.shifts[shiftID] = shift
	
	// Create order
	req := &order.CreateOrderRequest{
		ShiftID:      shiftID.Hex(),
		CustomerName: "Test Customer",
		Items: []order.OrderItem{
			{
				MenuItemID: menuItemID,
				Quantity:   1,
			},
		},
	}
	
	// Execute
	createdOrder, err := orderService.CreateOrder(ctx, req, "waiter123", "Test Waiter")
	
	// Assert
	assert.NoError(t, err)
	assert.NotNil(t, createdOrder)
	
	// Verify both batches were deducted
	assert.Equal(t, 970.0, coffeeBatch.QuantityRemaining) // 1000 - 30 = 970
	assert.Equal(t, 1850.0, milkBatch.QuantityRemaining)  // 2000 - 150 = 1850
	
	// Verify two usage logs were created
	assert.Equal(t, 2, len(batchUsageLogRepo.logs))
	
	// Verify total batch cost
	_ = 30*0.15 + 150*0.05 // 4.5 + 7.5 = 12.0
	assert.Contains(t, createdOrder.Note, "Batch Cost: 12.00 VND")
}

// Mock repositories for batch integration tests

type mockOrderRepoForBatchIntegration struct {
	orders map[primitive.ObjectID]*order.Order
}

func newMockOrderRepoForBatchIntegration() *mockOrderRepoForBatchIntegration {
	return &mockOrderRepoForBatchIntegration{
		orders: make(map[primitive.ObjectID]*order.Order),
	}
}

func (m *mockOrderRepoForBatchIntegration) Create(ctx context.Context, o *order.Order) error {
	if o.ID.IsZero() {
		o.ID = primitive.NewObjectID()
	}
	m.orders[o.ID] = o
	return nil
}

func (m *mockOrderRepoForBatchIntegration) FindByID(ctx context.Context, id primitive.ObjectID) (*order.Order, error) {
	if o, ok := m.orders[id]; ok {
		return o, nil
	}
	return nil, nil
}

func (m *mockOrderRepoForBatchIntegration) Update(ctx context.Context, id primitive.ObjectID, o *order.Order) error {
	m.orders[id] = o
	return nil
}

func (m *mockOrderRepoForBatchIntegration) Delete(ctx context.Context, id primitive.ObjectID) error {
	delete(m.orders, id)
	return nil
}

func (m *mockOrderRepoForBatchIntegration) FindByShiftID(ctx context.Context, shiftID primitive.ObjectID) ([]*order.Order, error) {
	return nil, nil
}

func (m *mockOrderRepoForBatchIntegration) FindByWaiterID(ctx context.Context, waiterID primitive.ObjectID) ([]*order.Order, error) {
	return nil, nil
}

func (m *mockOrderRepoForBatchIntegration) FindByStatus(ctx context.Context, status order.OrderStatus) ([]*order.Order, error) {
	return nil, nil
}

func (m *mockOrderRepoForBatchIntegration) FindByOrderNumber(ctx context.Context, orderNumber string) (*order.Order, error) {
	return nil, nil
}

func (m *mockOrderRepoForBatchIntegration) FindAll(ctx context.Context) ([]*order.Order, error) {
	return nil, nil
}

type mockShiftRepoForBatchIntegration struct {
	shifts map[primitive.ObjectID]*order.Shift
}

func newMockShiftRepoForBatchIntegration() *mockShiftRepoForBatchIntegration {
	return &mockShiftRepoForBatchIntegration{
		shifts: make(map[primitive.ObjectID]*order.Shift),
	}
}

func (m *mockShiftRepoForBatchIntegration) FindByID(ctx context.Context, id primitive.ObjectID) (*order.Shift, error) {
	if s, ok := m.shifts[id]; ok {
		return s, nil
	}
	return nil, nil
}

func (m *mockShiftRepoForBatchIntegration) Create(ctx context.Context, s *order.Shift) error {
	return nil
}

func (m *mockShiftRepoForBatchIntegration) Update(ctx context.Context, id primitive.ObjectID, s *order.Shift) error {
	return nil
}

func (m *mockShiftRepoForBatchIntegration) FindOpenShiftByUser(ctx context.Context, userID primitive.ObjectID, role order.RoleType) (*order.Shift, error) {
	return nil, nil
}

func (m *mockShiftRepoForBatchIntegration) FindByStatus(ctx context.Context, status order.ShiftStatus) ([]*order.Shift, error) {
	return nil, nil
}

func (m *mockShiftRepoForBatchIntegration) FindAll(ctx context.Context) ([]*order.Shift, error) {
	return nil, nil
}

func (m *mockShiftRepoForBatchIntegration) FindByDateRange(ctx context.Context, startDate, endDate time.Time) ([]*order.Shift, error) {
	return nil, nil
}

func (m *mockShiftRepoForBatchIntegration) FindByRoleType(ctx context.Context, roleType order.RoleType) ([]*order.Shift, error) {
	return nil, nil
}

func (m *mockShiftRepoForBatchIntegration) FindByUserID(ctx context.Context, userID primitive.ObjectID, roleType order.RoleType) ([]*order.Shift, error) {
	return nil, nil
}

func (m *mockShiftRepoForBatchIntegration) FindOpenShiftByWaiter(ctx context.Context, waiterID primitive.ObjectID) (*order.Shift, error) {
	return nil, nil
}

func (m *mockShiftRepoForBatchIntegration) FindOpenShifts(ctx context.Context) ([]*order.Shift, error) {
	return nil, nil
}

func (m *mockShiftRepoForBatchIntegration) FindByWaiterID(ctx context.Context, waiterID primitive.ObjectID) ([]*order.Shift, error) {
	return nil, nil
}

type mockMenuRepoForBatchIntegration struct {
	items map[primitive.ObjectID]*menu.MenuItem
}

func newMockMenuRepoForBatchIntegration() *mockMenuRepoForBatchIntegration {
	return &mockMenuRepoForBatchIntegration{
		items: make(map[primitive.ObjectID]*menu.MenuItem),
	}
}

func (m *mockMenuRepoForBatchIntegration) FindByID(ctx context.Context, id primitive.ObjectID) (*menu.MenuItem, error) {
	if item, ok := m.items[id]; ok {
		return item, nil
	}
	return nil, nil
}

func (m *mockMenuRepoForBatchIntegration) Create(ctx context.Context, item *menu.MenuItem) error {
	return nil
}

func (m *mockMenuRepoForBatchIntegration) FindAll(ctx context.Context) ([]*menu.MenuItem, error) {
	return nil, nil
}

func (m *mockMenuRepoForBatchIntegration) FindByCategory(ctx context.Context, category string) ([]*menu.MenuItem, error) {
	return nil, nil
}

func (m *mockMenuRepoForBatchIntegration) FindByIngredientName(ctx context.Context, ingredientName string) ([]*menu.MenuItem, error) {
	return nil, nil
}

func (m *mockMenuRepoForBatchIntegration) Update(ctx context.Context, id primitive.ObjectID, item *menu.MenuItem) error {
	return nil
}

func (m *mockMenuRepoForBatchIntegration) Delete(ctx context.Context, id primitive.ObjectID) error {
	return nil
}

type mockBatchRecordRepoForBatchIntegration struct {
	batches      map[primitive.ObjectID]*batch.BatchRecord
	batchesByDef map[primitive.ObjectID][]*batch.BatchRecord
}

func newMockBatchRecordRepoForBatchIntegration() *mockBatchRecordRepoForBatchIntegration {
	return &mockBatchRecordRepoForBatchIntegration{
		batches:      make(map[primitive.ObjectID]*batch.BatchRecord),
		batchesByDef: make(map[primitive.ObjectID][]*batch.BatchRecord),
	}
}

func (m *mockBatchRecordRepoForBatchIntegration) Create(ctx context.Context, record *batch.BatchRecord) error {
	return nil
}

func (m *mockBatchRecordRepoForBatchIntegration) Update(ctx context.Context, record *batch.BatchRecord) error {
	m.batches[record.ID] = record
	return nil
}

func (m *mockBatchRecordRepoForBatchIntegration) Delete(ctx context.Context, id primitive.ObjectID) error {
	return nil
}

func (m *mockBatchRecordRepoForBatchIntegration) FindByID(ctx context.Context, id primitive.ObjectID) (*batch.BatchRecord, error) {
	if b, ok := m.batches[id]; ok {
		return b, nil
	}
	return nil, nil
}

func (m *mockBatchRecordRepoForBatchIntegration) FindAll(ctx context.Context, filter batch.BatchRecordFilter) ([]*batch.BatchRecord, int64, error) {
	return nil, 0, nil
}

func (m *mockBatchRecordRepoForBatchIntegration) FindAvailableByDefinition(ctx context.Context, defID primitive.ObjectID) ([]*batch.BatchRecord, error) {
	if batches, ok := m.batchesByDef[defID]; ok {
		return batches, nil
	}
	return nil, nil
}

func (m *mockBatchRecordRepoForBatchIntegration) UpdateQuantity(ctx context.Context, id primitive.ObjectID, newQuantity float64) error {
	return nil
}

func (m *mockBatchRecordRepoForBatchIntegration) GetTotalAvailableQuantity(ctx context.Context, defID primitive.ObjectID) (float64, error) {
	return 0, nil
}

type mockBatchUsageLogRepoForBatchIntegration struct {
	logs []*batch.BatchUsageLog
}

func newMockBatchUsageLogRepoForBatchIntegration() *mockBatchUsageLogRepoForBatchIntegration {
	return &mockBatchUsageLogRepoForBatchIntegration{
		logs: make([]*batch.BatchUsageLog, 0),
	}
}

func (m *mockBatchUsageLogRepoForBatchIntegration) Create(ctx context.Context, log *batch.BatchUsageLog) error {
	if log.ID.IsZero() {
		log.ID = primitive.NewObjectID()
	}
	m.logs = append(m.logs, log)
	return nil
}

func (m *mockBatchUsageLogRepoForBatchIntegration) FindAll(ctx context.Context, filter batch.BatchUsageLogFilter) ([]*batch.BatchUsageLog, int64, error) {
	// Simple filter implementation for rollback testing
	if filter.OrderID != nil {
		var filtered []*batch.BatchUsageLog
		for _, log := range m.logs {
			if log.OrderID == *filter.OrderID {
				filtered = append(filtered, log)
			}
		}
		return filtered, int64(len(filtered)), nil
	}
	return m.logs, int64(len(m.logs)), nil
}

func (m *mockBatchUsageLogRepoForBatchIntegration) Delete(ctx context.Context, id primitive.ObjectID) error {
	for i, log := range m.logs {
		if log.ID == id {
			m.logs = append(m.logs[:i], m.logs[i+1:]...)
			return nil
		}
	}
	return nil
}
