package services

import (
	"context"
	"testing"
	"time"

	"cafe-pos/backend/domain/ingredient"
	"cafe-pos/backend/domain/menu"
	"cafe-pos/backend/domain/order"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

// Mock repositories for testing
type mockMenuRepository struct {
	menuItems map[primitive.ObjectID]*menu.MenuItem
}

func (m *mockMenuRepository) Create(ctx context.Context, item *menu.MenuItem) error {
	if item.ID.IsZero() {
		item.ID = primitive.NewObjectID()
	}
	m.menuItems[item.ID] = item
	return nil
}

func (m *mockMenuRepository) FindAll(ctx context.Context) ([]*menu.MenuItem, error) {
	items := make([]*menu.MenuItem, 0, len(m.menuItems))
	for _, item := range m.menuItems {
		items = append(items, item)
	}
	return items, nil
}

func (m *mockMenuRepository) FindByID(ctx context.Context, id primitive.ObjectID) (*menu.MenuItem, error) {
	item, exists := m.menuItems[id]
	if !exists {
		return nil, primitive.ErrInvalidHex
	}
	return item, nil
}

func (m *mockMenuRepository) Update(ctx context.Context, id primitive.ObjectID, item *menu.MenuItem) error {
	m.menuItems[id] = item
	return nil
}

func (m *mockMenuRepository) Delete(ctx context.Context, id primitive.ObjectID) error {
	delete(m.menuItems, id)
	return nil
}

func (m *mockMenuRepository) FindByIngredientName(ctx context.Context, ingredientName string) ([]*menu.MenuItem, error) {
	items := make([]*menu.MenuItem, 0)
	for _, item := range m.menuItems {
		for _, ing := range item.Ingredients {
			if ing.Name == ingredientName {
				items = append(items, item)
				break
			}
		}
	}
	return items, nil
}

func (m *mockMenuRepository) FindByCategory(ctx context.Context, category string) ([]*menu.MenuItem, error) {
	items := make([]*menu.MenuItem, 0)
	for _, item := range m.menuItems {
		if item.Category == category {
			items = append(items, item)
		}
	}
	return items, nil
}

func (m *mockMenuRepository) FindByBatchDefinitionID(ctx context.Context, batchDefID primitive.ObjectID) ([]*menu.MenuItem, error) {
	return nil, nil
}

type mockIngredientRepository struct {
	ingredients map[primitive.ObjectID]*ingredient.Ingredient
}

func (m *mockIngredientRepository) Create(ctx context.Context, item *ingredient.Ingredient) error {
	if item.ID.IsZero() {
		item.ID = primitive.NewObjectID()
	}
	m.ingredients[item.ID] = item
	return nil
}

func (m *mockIngredientRepository) FindAll(ctx context.Context) ([]*ingredient.Ingredient, error) {
	items := make([]*ingredient.Ingredient, 0, len(m.ingredients))
	for _, item := range m.ingredients {
		items = append(items, item)
	}
	return items, nil
}

func (m *mockIngredientRepository) FindByID(ctx context.Context, id primitive.ObjectID) (*ingredient.Ingredient, error) {
	item, exists := m.ingredients[id]
	if !exists {
		return nil, primitive.ErrInvalidHex
	}
	return item, nil
}

func (m *mockIngredientRepository) Update(ctx context.Context, id primitive.ObjectID, item *ingredient.Ingredient) error {
	m.ingredients[id] = item
	return nil
}

func (m *mockIngredientRepository) Delete(ctx context.Context, id primitive.ObjectID) error {
	delete(m.ingredients, id)
	return nil
}

func (m *mockIngredientRepository) FindLowStock(ctx context.Context) ([]*ingredient.Ingredient, error) {
	return []*ingredient.Ingredient{}, nil
}

func (m *mockIngredientRepository) CreateCategory(ctx context.Context, cat *ingredient.IngredientCategory) error {
	return nil
}

func (m *mockIngredientRepository) GetCategories(ctx context.Context) ([]ingredient.IngredientCategory, error) {
	return []ingredient.IngredientCategory{}, nil
}

func (m *mockIngredientRepository) DeleteCategory(ctx context.Context, id primitive.ObjectID) error {
	return nil
}

// Test helper to create test data
func setupTestData() (*mockMenuRepository, *mockIngredientRepository) {
	menuRepo := &mockMenuRepository{
		menuItems: make(map[primitive.ObjectID]*menu.MenuItem),
	}

	espressoID := primitive.NewObjectID()
	milkID := primitive.NewObjectID()
	sugarID := primitive.NewObjectID()
	chocolateID := primitive.NewObjectID()

	ingredientRepo := &mockIngredientRepository{
		ingredients: map[primitive.ObjectID]*ingredient.Ingredient{
			espressoID: {
				ID:                espressoID,
				Name:              "Espresso",
				CostPerUnit:       200.0,
				ConversionRate:    1.0,
				WastagePercentage: 5.0,
			},
			milkID: {
				ID:                milkID,
				Name:              "Milk",
				CostPerUnit:       50.0,
				ConversionRate:    1.0,
				WastagePercentage: 10.0,
			},
			sugarID: {
				ID:                sugarID,
				Name:              "Sugar",
				CostPerUnit:       30.0,
				ConversionRate:    1.0,
				WastagePercentage: 0.0,
			},
			chocolateID: {
				ID:                chocolateID,
				Name:              "Chocolate",
				CostPerUnit:       0.0, // Missing cost
				ConversionRate:    1.0,
				WastagePercentage: 0.0,
			},
		},
	}

	return menuRepo, ingredientRepo
}

func TestCalculateMenuItemCost_BasicCalculation(t *testing.T) {
	menuRepo, ingredientRepo := setupTestData()
	orderRepo := &mockOrderRepository{orders: make(map[primitive.ObjectID]*order.Order)}
	orderItemRepo := &mockOrderItemRepository{orderItems: make([]*order.OrderItemWithCost, 0)}
	service := NewCostCalculatorService(menuRepo, ingredientRepo, orderRepo, orderItemRepo)

	// Create a menu item with ingredients
	menuItem := &menu.MenuItem{
		ID:   primitive.NewObjectID(),
		Name: "Cappuccino",
		Ingredients: []menu.Ingredient{
			{Name: "Espresso", Quantity: 30, Unit: ingredient.UnitMilliliter},  // 30 * 200 * 1.0 * 1.05 = 6300
			{Name: "Milk", Quantity: 150, Unit: ingredient.UnitMilliliter},     // 150 * 50 * 1.0 * 1.10 = 8250
		},
	}
	menuRepo.menuItems[menuItem.ID] = menuItem

	// Calculate cost
	result, err := service.CalculateMenuItemCost(context.Background(), menuItem.ID)
	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}

	// Expected: 6300 + 8250 = 14550
	expectedCost := 14550.0
	if result.CurrentCost != expectedCost {
		t.Errorf("Expected cost %v, got %v", expectedCost, result.CurrentCost)
	}

	if result.CostStatus != menu.CostStatusFinal {
		t.Errorf("Expected status FINAL, got %v", result.CostStatus)
	}

	if len(result.MissingIngredients) != 0 {
		t.Errorf("Expected no missing ingredients, got %v", result.MissingIngredients)
	}
}

func TestCalculateMenuItemCost_WithMissingCost(t *testing.T) {
	menuRepo, ingredientRepo := setupTestData()
	orderRepo := &mockOrderRepository{orders: make(map[primitive.ObjectID]*order.Order)}
	orderItemRepo := &mockOrderItemRepository{orderItems: make([]*order.OrderItemWithCost, 0)}
	service := NewCostCalculatorService(menuRepo, ingredientRepo, orderRepo, orderItemRepo)

	// Create a menu item with an ingredient that has missing cost
	menuItem := &menu.MenuItem{
		ID:   primitive.NewObjectID(),
		Name: "Mocha",
		Ingredients: []menu.Ingredient{
			{Name: "Espresso", Quantity: 30, Unit: ingredient.UnitMilliliter},
			{Name: "Chocolate", Quantity: 20, Unit: ingredient.UnitGram}, // Missing cost
		},
	}
	menuRepo.menuItems[menuItem.ID] = menuItem

	// Calculate cost
	result, err := service.CalculateMenuItemCost(context.Background(), menuItem.ID)
	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}

	// Expected: only Espresso cost = 30 * 200 * 1.0 * 1.05 = 6300
	expectedCost := 6300.0
	if result.CurrentCost != expectedCost {
		t.Errorf("Expected cost %v, got %v", expectedCost, result.CurrentCost)
	}

	if result.CostStatus != menu.CostStatusIncomplete {
		t.Errorf("Expected status INCOMPLETE, got %v", result.CostStatus)
	}

	if len(result.MissingIngredients) != 1 || result.MissingIngredients[0] != "Chocolate" {
		t.Errorf("Expected missing ingredient 'Chocolate', got %v", result.MissingIngredients)
	}
}

func TestCalculateMenuItemCost_NoIngredients(t *testing.T) {
	menuRepo, ingredientRepo := setupTestData()
	orderRepo := &mockOrderRepository{orders: make(map[primitive.ObjectID]*order.Order)}
	orderItemRepo := &mockOrderItemRepository{orderItems: make([]*order.OrderItemWithCost, 0)}
	service := NewCostCalculatorService(menuRepo, ingredientRepo, orderRepo, orderItemRepo)

	// Create a menu item with no ingredients
	menuItem := &menu.MenuItem{
		ID:          primitive.NewObjectID(),
		Name:        "Service Item",
		Ingredients: []menu.Ingredient{},
	}
	menuRepo.menuItems[menuItem.ID] = menuItem

	// Calculate cost
	result, err := service.CalculateMenuItemCost(context.Background(), menuItem.ID)
	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}

	if result.CurrentCost != 0.0 {
		t.Errorf("Expected cost 0, got %v", result.CurrentCost)
	}

	if result.CostStatus != menu.CostStatusFinal {
		t.Errorf("Expected status FINAL, got %v", result.CostStatus)
	}
}

func TestCalculateMenuItemCost_RoundingTo2Decimals(t *testing.T) {
	menuRepo, ingredientRepo := setupTestData()
	orderRepo := &mockOrderRepository{orders: make(map[primitive.ObjectID]*order.Order)}
	orderItemRepo := &mockOrderItemRepository{orderItems: make([]*order.OrderItemWithCost, 0)}
	service := NewCostCalculatorService(menuRepo, ingredientRepo, orderRepo, orderItemRepo)

	// Add ingredient with price that will create decimal places
	coffeeBeanID := primitive.NewObjectID()
	ingredientRepo.ingredients[coffeeBeanID] = &ingredient.Ingredient{
		ID:                coffeeBeanID,
		Name:              "Coffee Bean",
		CostPerUnit:       33.333,
		ConversionRate:    1.0,
		WastagePercentage: 0.0,
	}

	menuItem := &menu.MenuItem{
		ID:   primitive.NewObjectID(),
		Name: "Americano",
		Ingredients: []menu.Ingredient{
			{Name: "Coffee Bean", Quantity: 10, Unit: ingredient.UnitGram}, // 10 * 33.333 = 333.33
		},
	}
	menuRepo.menuItems[menuItem.ID] = menuItem

	result, err := service.CalculateMenuItemCost(context.Background(), menuItem.ID)
	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}

	// Should be rounded to 2 decimal places
	expectedCost := 333.33
	if result.CurrentCost != expectedCost {
		t.Errorf("Expected cost %v, got %v", expectedCost, result.CurrentCost)
	}
}

func TestCalculateMenuItemCost_WithConversionAndWastage(t *testing.T) {
	menuRepo, ingredientRepo := setupTestData()
	orderRepo := &mockOrderRepository{orders: make(map[primitive.ObjectID]*order.Order)}
	orderItemRepo := &mockOrderItemRepository{orderItems: make([]*order.OrderItemWithCost, 0)}
	service := NewCostCalculatorService(menuRepo, ingredientRepo, orderRepo, orderItemRepo)

	// Add ingredient with stock unit in kg
	flourID := primitive.NewObjectID()
	ingredientRepo.ingredients[flourID] = &ingredient.Ingredient{
		ID:                flourID,
		Name:              "Flour",
		Unit:              ingredient.UnitKilogram,  // Stock unit: kg
		CostPerUnit:       100000.0,                 // 100,000 VND per kg
		WastagePercentage: 15.0,                     // 15% wastage
	}

	// Menu uses grams (recipe unit different from stock unit)
	menuItem := &menu.MenuItem{
		ID:   primitive.NewObjectID(),
		Name: "Cake",
		Ingredients: []menu.Ingredient{
			// 50g flour
			// Conversion: g → kg = 0.001
			// Cost = 50 * 100,000 * 0.001 * 1.15 = 5,750
			{Name: "Flour", Quantity: 50, Unit: ingredient.UnitGram},
		},
	}
	menuRepo.menuItems[menuItem.ID] = menuItem

	result, err := service.CalculateMenuItemCost(context.Background(), menuItem.ID)
	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}

	// Expected: 50g * 100,000 VND/kg * 0.001 (g→kg) * 1.15 (wastage) = 5,750
	expectedCost := 5750.0
	if result.CurrentCost != expectedCost {
		t.Errorf("Expected cost %v, got %v", expectedCost, result.CurrentCost)
	}
}

func TestCalculateMenuItemCost_DefaultConversionAndWastage(t *testing.T) {
	menuRepo, ingredientRepo := setupTestData()
	orderRepo := &mockOrderRepository{orders: make(map[primitive.ObjectID]*order.Order)}
	orderItemRepo := &mockOrderItemRepository{orderItems: make([]*order.OrderItemWithCost, 0)}
	service := NewCostCalculatorService(menuRepo, ingredientRepo, orderRepo, orderItemRepo)

	// Add ingredient with zero/negative conversion and wastage (should default to 1.0 and 0.0)
	waterID := primitive.NewObjectID()
	ingredientRepo.ingredients[waterID] = &ingredient.Ingredient{
		ID:                waterID,
		Name:              "Water",
		CostPerUnit:       10.0,
		ConversionRate:    0.0,  // Should default to 1.0
		WastagePercentage: -5.0, // Should default to 0.0
	}

	menuItem := &menu.MenuItem{
		ID:   primitive.NewObjectID(),
		Name: "Tea",
		Ingredients: []menu.Ingredient{
			{Name: "Water", Quantity: 200, Unit: ingredient.UnitMilliliter}, // 200 * 10 * 1.0 * 1.0 = 2000
		},
	}
	menuRepo.menuItems[menuItem.ID] = menuItem

	result, err := service.CalculateMenuItemCost(context.Background(), menuItem.ID)
	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}

	expectedCost := 2000.0
	if result.CurrentCost != expectedCost {
		t.Errorf("Expected cost %v, got %v", expectedCost, result.CurrentCost)
	}
}

func TestCalculateMenuItemCost_IngredientNotInDatabase(t *testing.T) {
	menuRepo, ingredientRepo := setupTestData()
	orderRepo := &mockOrderRepository{orders: make(map[primitive.ObjectID]*order.Order)}
	orderItemRepo := &mockOrderItemRepository{orderItems: make([]*order.OrderItemWithCost, 0)}
	service := NewCostCalculatorService(menuRepo, ingredientRepo, orderRepo, orderItemRepo)

	menuItem := &menu.MenuItem{
		ID:   primitive.NewObjectID(),
		Name: "Mystery Drink",
		Ingredients: []menu.Ingredient{
			{Name: "Unknown Ingredient", Quantity: 10, Unit: ingredient.UnitGram},
		},
	}
	menuRepo.menuItems[menuItem.ID] = menuItem

	result, err := service.CalculateMenuItemCost(context.Background(), menuItem.ID)
	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}

	if result.CostStatus != menu.CostStatusIncomplete {
		t.Errorf("Expected status INCOMPLETE, got %v", result.CostStatus)
	}

	if len(result.MissingIngredients) != 1 || result.MissingIngredients[0] != "Unknown Ingredient" {
		t.Errorf("Expected missing ingredient 'Unknown Ingredient', got %v", result.MissingIngredients)
	}
}

type mockOrderRepository struct {
	orders map[primitive.ObjectID]*order.Order
}

func (m *mockOrderRepository) Create(ctx context.Context, o *order.Order) error {
	if o.ID.IsZero() {
		o.ID = primitive.NewObjectID()
	}
	m.orders[o.ID] = o
	return nil
}

func (m *mockOrderRepository) FindByID(ctx context.Context, id primitive.ObjectID) (*order.Order, error) {
	o, exists := m.orders[id]
	if !exists {
		return nil, primitive.ErrInvalidHex
	}
	return o, nil
}

func (m *mockOrderRepository) Update(ctx context.Context, id primitive.ObjectID, o *order.Order) error {
	m.orders[id] = o
	return nil
}

func (m *mockOrderRepository) FindByShiftID(ctx context.Context, shiftID primitive.ObjectID) ([]*order.Order, error) {
	orders := make([]*order.Order, 0)
	for _, o := range m.orders {
		if o.ShiftID == shiftID {
			orders = append(orders, o)
		}
	}
	return orders, nil
}

func (m *mockOrderRepository) FindByWaiterID(ctx context.Context, waiterID primitive.ObjectID) ([]*order.Order, error) {
	return []*order.Order{}, nil
}

func (m *mockOrderRepository) FindByStatus(ctx context.Context, status order.OrderStatus) ([]*order.Order, error) {
	return []*order.Order{}, nil
}

func (m *mockOrderRepository) FindAll(ctx context.Context) ([]*order.Order, error) {
	orders := make([]*order.Order, 0, len(m.orders))
	for _, o := range m.orders {
		orders = append(orders, o)
	}
	return orders, nil
}

func (m *mockOrderRepository) FindByOrderNumber(ctx context.Context, orderNumber string) (*order.Order, error) {
	return nil, primitive.ErrInvalidHex
}

func (m *mockOrderRepository) Delete(ctx context.Context, id primitive.ObjectID) error {
	delete(m.orders, id)
	return nil
}

func (m *mockOrderRepository) FindByIDs(ctx context.Context, ids []primitive.ObjectID) ([]*order.Order, error) {
	var result []*order.Order
	for _, id := range ids {
		if o, ok := m.orders[id]; ok {
			result = append(result, o)
		}
	}
	return result, nil
}

type mockOrderItemRepository struct {
	orderItems []*order.OrderItemWithCost
}

func (m *mockOrderItemRepository) CreateMany(ctx context.Context, items []*order.OrderItemWithCost) error {
	for _, item := range items {
		if item.ID.IsZero() {
			item.ID = primitive.NewObjectID()
		}
		m.orderItems = append(m.orderItems, item)
	}
	return nil
}

func (m *mockOrderItemRepository) FindByOrderID(ctx context.Context, orderID primitive.ObjectID) ([]*order.OrderItemWithCost, error) {
	items := make([]*order.OrderItemWithCost, 0)
	for _, item := range m.orderItems {
		if item.OrderID == orderID {
			items = append(items, item)
		}
	}
	return items, nil
}

func (m *mockOrderItemRepository) FindByOrderIDs(ctx context.Context, orderIDs []primitive.ObjectID) ([]*order.OrderItemWithCost, error) {
	items := make([]*order.OrderItemWithCost, 0)
	orderIDMap := make(map[primitive.ObjectID]bool)
	for _, id := range orderIDs {
		orderIDMap[id] = true
	}
	for _, item := range m.orderItems {
		if orderIDMap[item.OrderID] {
			items = append(items, item)
		}
	}
	return items, nil
}

func (m *mockOrderItemRepository) FindByDateRange(ctx context.Context, startDate, endDate time.Time) ([]*order.OrderItemWithCost, error) {
	items := make([]*order.OrderItemWithCost, 0)
	for _, item := range m.orderItems {
		if item.CreatedAt.After(startDate) && item.CreatedAt.Before(endDate) {
			items = append(items, item)
		}
	}
	return items, nil
}

// TestCalculateShiftOrderCosts tests the shift order cost calculation
func TestCalculateShiftOrderCosts(t *testing.T) {
	ctx := context.Background()

	// Setup test data
	shiftID := primitive.NewObjectID()
	menuItemID1 := primitive.NewObjectID()
	menuItemID2 := primitive.NewObjectID()

	// Create mock repositories
	menuRepo := &mockMenuRepository{
		menuItems: make(map[primitive.ObjectID]*menu.MenuItem),
	}
	ingredientRepo := &mockIngredientRepository{
		ingredients: make(map[primitive.ObjectID]*ingredient.Ingredient),
	}
	orderRepo := &mockOrderRepository{
		orders: make(map[primitive.ObjectID]*order.Order),
	}
	orderItemRepo := &mockOrderItemRepository{
		orderItems: make([]*order.OrderItemWithCost, 0),
	}

	// Create test ingredients
	espresso := &ingredient.Ingredient{
		ID:                primitive.NewObjectID(),
		Name:              "Espresso",
		CostPerUnit:       200.0,
		ConversionRate:    1.0,
		WastagePercentage: 5.0,
	}
	milk := &ingredient.Ingredient{
		ID:                primitive.NewObjectID(),
		Name:              "Milk",
		CostPerUnit:       50.0,
		ConversionRate:    1.0,
		WastagePercentage: 10.0,
	}
	ingredientRepo.ingredients[espresso.ID] = espresso
	ingredientRepo.ingredients[milk.ID] = milk

	// Create test menu items
	cappuccino := &menu.MenuItem{
		ID:    menuItemID1,
		Name:  "Cappuccino",
		Price: 45000,
		Ingredients: []menu.Ingredient{
			{Name: "Espresso", Quantity: 30},
			{Name: "Milk", Quantity: 150},
		},
	}
	latte := &menu.MenuItem{
		ID:    menuItemID2,
		Name:  "Latte",
		Price: 50000,
		Ingredients: []menu.Ingredient{
			{Name: "Espresso", Quantity: 30},
			{Name: "Milk", Quantity: 200},
		},
	}
	menuRepo.menuItems[menuItemID1] = cappuccino
	menuRepo.menuItems[menuItemID2] = latte

	// Create test orders
	order1 := &order.Order{
		ID:      primitive.NewObjectID(),
		ShiftID: shiftID,
		Items: []order.OrderItem{
			{
				MenuItemID: menuItemID1,
				Name:       "Cappuccino",
				Price:      45000,
				Quantity:   2,
				Subtotal:   90000,
			},
		},
	}
	order2 := &order.Order{
		ID:      primitive.NewObjectID(),
		ShiftID: shiftID,
		Items: []order.OrderItem{
			{
				MenuItemID: menuItemID2,
				Name:       "Latte",
				Price:      50000,
				Quantity:   1,
				Subtotal:   50000,
			},
			{
				MenuItemID: menuItemID1,
				Name:       "Cappuccino",
				Price:      45000,
				Quantity:   1,
				Subtotal:   45000,
			},
		},
	}
	orderRepo.orders[order1.ID] = order1
	orderRepo.orders[order2.ID] = order2

	// Create service
	service := NewCostCalculatorService(menuRepo, ingredientRepo, orderRepo, orderItemRepo)

	// Test: Calculate shift order costs
	result, err := service.CalculateShiftOrderCosts(ctx, shiftID)
	if err != nil {
		t.Fatalf("CalculateShiftOrderCosts failed: %v", err)
	}

	// Verify result
	if result.TotalOrders != 2 {
		t.Errorf("Expected 2 orders, got %d", result.TotalOrders)
	}
	if result.TotalItems != 3 {
		t.Errorf("Expected 3 items, got %d", result.TotalItems)
	}
	if result.ItemsWithFinalCost != 3 {
		t.Errorf("Expected 3 items with final cost, got %d", result.ItemsWithFinalCost)
	}
	if result.ItemsWithIncompleteCost != 0 {
		t.Errorf("Expected 0 items with incomplete cost, got %d", result.ItemsWithIncompleteCost)
	}

	// Verify order items were created
	if len(orderItemRepo.orderItems) != 3 {
		t.Errorf("Expected 3 order items created, got %d", len(orderItemRepo.orderItems))
	}

	// Verify cost calculation for first order item (Cappuccino x2)
	// Espresso: 30 * 200 * 1.0 * 1.05 = 6300
	// Milk: 150 * 50 * 1.0 * 1.10 = 8250
	// Total per item: 14550
	// Total for 2 items: 29100
	expectedCostPerCappuccino := 14550.0
	expectedTotalForOrder1 := 29100.0

	orderItem1 := orderItemRepo.orderItems[0]
	if orderItem1.AccountingCost != expectedTotalForOrder1 {
		t.Errorf("Expected accounting cost %.2f for order 1, got %.2f", expectedTotalForOrder1, orderItem1.AccountingCost)
	}
	if orderItem1.CostStatus != order.CostStatusFinal {
		t.Errorf("Expected cost status FINAL, got %s", orderItem1.CostStatus)
	}

	// Verify cost calculation for second order item (Latte x1)
	// Espresso: 30 * 200 * 1.0 * 1.05 = 6300
	// Milk: 200 * 50 * 1.0 * 1.10 = 11000
	// Total: 17300
	expectedCostForLatte := 17300.0

	orderItem2 := orderItemRepo.orderItems[1]
	if orderItem2.AccountingCost != expectedCostForLatte {
		t.Errorf("Expected accounting cost %.2f for latte, got %.2f", expectedCostForLatte, orderItem2.AccountingCost)
	}

	// Verify total accounting cost
	expectedTotalCost := expectedTotalForOrder1 + expectedCostForLatte + expectedCostPerCappuccino
	if result.TotalAccountingCost != expectedTotalCost {
		t.Errorf("Expected total accounting cost %.2f, got %.2f", expectedTotalCost, result.TotalAccountingCost)
	}
}

// TestCalculateShiftOrderCosts_EmptyShift tests calculation with no orders
func TestCalculateShiftOrderCosts_EmptyShift(t *testing.T) {
	ctx := context.Background()
	shiftID := primitive.NewObjectID()

	// Create mock repositories
	menuRepo := &mockMenuRepository{menuItems: make(map[primitive.ObjectID]*menu.MenuItem)}
	ingredientRepo := &mockIngredientRepository{ingredients: make(map[primitive.ObjectID]*ingredient.Ingredient)}
	orderRepo := &mockOrderRepository{orders: make(map[primitive.ObjectID]*order.Order)}
	orderItemRepo := &mockOrderItemRepository{orderItems: make([]*order.OrderItemWithCost, 0)}

	service := NewCostCalculatorService(menuRepo, ingredientRepo, orderRepo, orderItemRepo)

	result, err := service.CalculateShiftOrderCosts(ctx, shiftID)
	if err != nil {
		t.Fatalf("CalculateShiftOrderCosts failed: %v", err)
	}

	if result.TotalOrders != 0 {
		t.Errorf("Expected 0 orders, got %d", result.TotalOrders)
	}
	if result.TotalItems != 0 {
		t.Errorf("Expected 0 items, got %d", result.TotalItems)
	}
}

// TestCalculateShiftOrderCosts_IncompleteCost tests handling of missing ingredient costs
func TestCalculateShiftOrderCosts_IncompleteCost(t *testing.T) {
	ctx := context.Background()
	shiftID := primitive.NewObjectID()
	menuItemID := primitive.NewObjectID()

	// Create mock repositories
	menuRepo := &mockMenuRepository{menuItems: make(map[primitive.ObjectID]*menu.MenuItem)}
	ingredientRepo := &mockIngredientRepository{ingredients: make(map[primitive.ObjectID]*ingredient.Ingredient)}
	orderRepo := &mockOrderRepository{orders: make(map[primitive.ObjectID]*order.Order)}
	orderItemRepo := &mockOrderItemRepository{orderItems: make([]*order.OrderItemWithCost, 0)}

	// Create ingredient with missing cost
	espresso := &ingredient.Ingredient{
		ID:          primitive.NewObjectID(),
		Name:        "Espresso",
		CostPerUnit: 0.0, // Missing cost
	}
	ingredientRepo.ingredients[espresso.ID] = espresso

	// Create menu item
	cappuccino := &menu.MenuItem{
		ID:    menuItemID,
		Name:  "Cappuccino",
		Price: 45000,
		Ingredients: []menu.Ingredient{
			{Name: "Espresso", Quantity: 30},
		},
	}
	menuRepo.menuItems[menuItemID] = cappuccino

	// Create order
	testOrder := &order.Order{
		ID:      primitive.NewObjectID(),
		ShiftID: shiftID,
		Items: []order.OrderItem{
			{
				MenuItemID: menuItemID,
				Name:       "Cappuccino",
				Price:      45000,
				Quantity:   1,
				Subtotal:   45000,
			},
		},
	}
	orderRepo.orders[testOrder.ID] = testOrder

	service := NewCostCalculatorService(menuRepo, ingredientRepo, orderRepo, orderItemRepo)

	result, err := service.CalculateShiftOrderCosts(ctx, shiftID)
	if err != nil {
		t.Fatalf("CalculateShiftOrderCosts failed: %v", err)
	}

	if result.ItemsWithIncompleteCost != 1 {
		t.Errorf("Expected 1 item with incomplete cost, got %d", result.ItemsWithIncompleteCost)
	}
	if result.ItemsWithFinalCost != 0 {
		t.Errorf("Expected 0 items with final cost, got %d", result.ItemsWithFinalCost)
	}

	// Verify order item was created with INCOMPLETE status
	if len(orderItemRepo.orderItems) != 1 {
		t.Fatalf("Expected 1 order item created, got %d", len(orderItemRepo.orderItems))
	}
	if orderItemRepo.orderItems[0].CostStatus != order.CostStatusIncomplete {
		t.Errorf("Expected cost status INCOMPLETE, got %s", orderItemRepo.orderItems[0].CostStatus)
	}
}

// TestQueueCostRecalculation tests the QueueCostRecalculation method
func TestQueueCostRecalculation(t *testing.T) {
	menuRepo, ingredientRepo := setupTestData()
	orderRepo := &mockOrderRepository{orders: make(map[primitive.ObjectID]*order.Order)}
	orderItemRepo := &mockOrderItemRepository{orderItems: make([]*order.OrderItemWithCost, 0)}
	service := NewCostCalculatorService(menuRepo, ingredientRepo, orderRepo, orderItemRepo)
	
	// Create recalculation service and wire it up
	recalcService := NewCostRecalculationService(service, menuRepo, 2, 100)
	service.SetCostRecalculationService(recalcService)
	recalcService.Start()
	defer recalcService.Stop()

	// Create menu items that use Espresso
	cappuccino := &menu.MenuItem{
		ID:   primitive.NewObjectID(),
		Name: "Cappuccino",
		Ingredients: []menu.Ingredient{
			{Name: "Espresso", Quantity: 30, Unit: ingredient.UnitMilliliter},
			{Name: "Milk", Quantity: 150, Unit: ingredient.UnitMilliliter},
		},
	}
	menuRepo.menuItems[cappuccino.ID] = cappuccino

	latte := &menu.MenuItem{
		ID:   primitive.NewObjectID(),
		Name: "Latte",
		Ingredients: []menu.Ingredient{
			{Name: "Espresso", Quantity: 30, Unit: ingredient.UnitMilliliter},
			{Name: "Milk", Quantity: 200, Unit: ingredient.UnitMilliliter},
		},
	}
	menuRepo.menuItems[latte.ID] = latte

	americano := &menu.MenuItem{
		ID:   primitive.NewObjectID(),
		Name: "Americano",
		Ingredients: []menu.Ingredient{
			{Name: "Espresso", Quantity: 60, Unit: ingredient.UnitMilliliter},
		},
	}
	menuRepo.menuItems[americano.ID] = americano

	// Create a menu item that doesn't use Espresso
	milkshake := &menu.MenuItem{
		ID:   primitive.NewObjectID(),
		Name: "Milkshake",
		Ingredients: []menu.Ingredient{
			{Name: "Milk", Quantity: 300, Unit: ingredient.UnitMilliliter},
			{Name: "Sugar", Quantity: 20, Unit: ingredient.UnitGram},
		},
	}
	menuRepo.menuItems[milkshake.ID] = milkshake

	// Get the Espresso ingredient ID
	var espressoID primitive.ObjectID
	for id, ing := range ingredientRepo.ingredients {
		if ing.Name == "Espresso" {
			espressoID = id
			break
		}
	}

	// Queue cost recalculation for Espresso
	err := service.QueueCostRecalculation(context.Background(), espressoID)
	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}

	// Wait for recalculations to complete
	time.Sleep(500 * time.Millisecond)

	// Verify that the menu items using Espresso were recalculated
	// Check that their costs were updated
	cappuccinoUpdated := menuRepo.menuItems[cappuccino.ID]
	if cappuccinoUpdated.CurrentCost == 0 {
		t.Error("Expected Cappuccino cost to be updated")
	}

	latteUpdated := menuRepo.menuItems[latte.ID]
	if latteUpdated.CurrentCost == 0 {
		t.Error("Expected Latte cost to be updated")
	}

	americanoUpdated := menuRepo.menuItems[americano.ID]
	if americanoUpdated.CurrentCost == 0 {
		t.Error("Expected Americano cost to be updated")
	}

	// Milkshake should not be updated (doesn't use Espresso)
	milkshakeUpdated := menuRepo.menuItems[milkshake.ID]
	if milkshakeUpdated.CurrentCost != 0 {
		t.Error("Expected Milkshake cost to remain 0 (doesn't use Espresso)")
	}
}

// TestQueueCostRecalculation_NoMenuItems tests queuing when no menu items use the ingredient
func TestQueueCostRecalculation_NoMenuItems(t *testing.T) {
	menuRepo, ingredientRepo := setupTestData()
	orderRepo := &mockOrderRepository{orders: make(map[primitive.ObjectID]*order.Order)}
	orderItemRepo := &mockOrderItemRepository{orderItems: make([]*order.OrderItemWithCost, 0)}
	service := NewCostCalculatorService(menuRepo, ingredientRepo, orderRepo, orderItemRepo)
	
	// Create recalculation service and wire it up
	recalcService := NewCostRecalculationService(service, menuRepo, 2, 100)
	service.SetCostRecalculationService(recalcService)
	recalcService.Start()
	defer recalcService.Stop()

	// Create a menu item that doesn't use Chocolate
	cappuccino := &menu.MenuItem{
		ID:   primitive.NewObjectID(),
		Name: "Cappuccino",
		Ingredients: []menu.Ingredient{
			{Name: "Espresso", Quantity: 30, Unit: ingredient.UnitMilliliter},
			{Name: "Milk", Quantity: 150, Unit: ingredient.UnitMilliliter},
		},
	}
	menuRepo.menuItems[cappuccino.ID] = cappuccino

	// Get the Chocolate ingredient ID
	var chocolateID primitive.ObjectID
	for id, ing := range ingredientRepo.ingredients {
		if ing.Name == "Chocolate" {
			chocolateID = id
			break
		}
	}

	// Queue cost recalculation for Chocolate
	err := service.QueueCostRecalculation(context.Background(), chocolateID)
	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}

	// Wait a bit for any potential recalculations
	time.Sleep(200 * time.Millisecond)

	// Check that cappuccino was not updated (doesn't use Chocolate)
	cappuccinoUpdated := menuRepo.menuItems[cappuccino.ID]
	if cappuccinoUpdated.CurrentCost != 0 {
		t.Error("Expected Cappuccino cost to remain 0 (doesn't use Chocolate)")
	}
}
// TestQueueCostRecalculation_InvalidIngredient tests queuing with an invalid ingredient ID
func TestQueueCostRecalculation_InvalidIngredient(t *testing.T) {
	menuRepo, ingredientRepo := setupTestData()
	orderRepo := &mockOrderRepository{orders: make(map[primitive.ObjectID]*order.Order)}
	orderItemRepo := &mockOrderItemRepository{orderItems: make([]*order.OrderItemWithCost, 0)}
	service := NewCostCalculatorService(menuRepo, ingredientRepo, orderRepo, orderItemRepo)

	// Try to queue with an invalid ingredient ID
	invalidID := primitive.NewObjectID()
	err := service.QueueCostRecalculation(context.Background(), invalidID)
	if err == nil {
		t.Error("Expected error for invalid ingredient ID, got nil")
	}
}

// Test variant cost calculation

func TestCalculateMenuItemCost_MultiSize_WithVariants(t *testing.T) {
	menuRepo, ingredientRepo := setupTestData()
	orderRepo := &mockOrderRepository{orders: make(map[primitive.ObjectID]*order.Order)}
	orderItemRepo := &mockOrderItemRepository{orderItems: make([]*order.OrderItemWithCost, 0)}
	service := NewCostCalculatorService(menuRepo, ingredientRepo, orderRepo, orderItemRepo)

	// Add ingredients needed for test
	coffeeID := primitive.NewObjectID()
	milkID := primitive.NewObjectID()
	ingredientRepo.ingredients[coffeeID] = &ingredient.Ingredient{
		ID:                coffeeID,
		Name:              "Cà phê",
		Unit:              ingredient.UnitGram,
		CostPerUnit:       500.0, // 500đ per gram
		WastagePercentage: 0.0,
	}
	ingredientRepo.ingredients[milkID] = &ingredient.Ingredient{
		ID:                milkID,
		Name:              "Sữa đặc",
		Unit:              ingredient.UnitMilliliter,
		CostPerUnit:       200.0, // 200đ per ml
		WastagePercentage: 0.0,
	}

	// Create multi-size menu item with variants
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
				Ingredients: []menu.Ingredient{
					{Name: "Cà phê", Quantity: 20, Unit: ingredient.UnitGram},
					{Name: "Sữa đặc", Quantity: 30, Unit: ingredient.UnitMilliliter},
				},
			},
			{
				ID:        "L",
				Name:      "Size L",
				Price:     30000,
				Available: true,
				IsDefault: false,
				Ingredients: []menu.Ingredient{
					{Name: "Cà phê", Quantity: 30, Unit: ingredient.UnitGram},
					{Name: "Sữa đặc", Quantity: 45, Unit: ingredient.UnitMilliliter},
				},
			},
		},
	}
	menuRepo.Create(context.Background(), menuItem)

	// Calculate cost
	result, err := service.CalculateMenuItemCost(context.Background(), menuItem.ID)
	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}

	// Verify result is for default variant
	if result.CostStatus != menu.CostStatusFinal {
		t.Errorf("Expected cost status FINAL, got %s", result.CostStatus)
	}

	// Fetch updated menu item to verify variant costs
	updated, _ := menuRepo.FindByID(context.Background(), menuItem.ID)
	
	// Verify Size M cost
	variantM := updated.GetVariantByID("M")
	if variantM == nil {
		t.Fatal("Expected variant M to exist")
	}
	// Cost = (20g * 500đ/g) + (30ml * 200đ/ml) = 10,000 + 6,000 = 16,000
	expectedCostM := 16000.0
	if variantM.CurrentCost != expectedCostM {
		t.Errorf("Expected variant M cost %f, got %f", expectedCostM, variantM.CurrentCost)
	}
	if variantM.CostStatus != menu.CostStatusFinal {
		t.Errorf("Expected variant M status FINAL, got %s", variantM.CostStatus)
	}

	// Verify Size L cost
	variantL := updated.GetVariantByID("L")
	if variantL == nil {
		t.Fatal("Expected variant L to exist")
	}
	// Cost = (30g * 500đ/g) + (45ml * 200đ/ml) = 15,000 + 9,000 = 24,000
	expectedCostL := 24000.0
	if variantL.CurrentCost != expectedCostL {
		t.Errorf("Expected variant L cost %f, got %f", expectedCostL, variantL.CurrentCost)
	}
	if variantL.CostStatus != menu.CostStatusFinal {
		t.Errorf("Expected variant L status FINAL, got %s", variantL.CostStatus)
	}

	// Verify old cost fields are cleared
	if updated.CurrentCost != 0 {
		t.Errorf("Expected menu item CurrentCost to be 0 (cleared), got %f", updated.CurrentCost)
	}
	if updated.CostStatus != "" {
		t.Errorf("Expected menu item CostStatus to be empty (cleared), got %s", updated.CostStatus)
	}
}

func TestCalculateMenuItemCost_MultiSize_WithMissingCost(t *testing.T) {
	menuRepo, ingredientRepo := setupTestData()
	orderRepo := &mockOrderRepository{orders: make(map[primitive.ObjectID]*order.Order)}
	orderItemRepo := &mockOrderItemRepository{orderItems: make([]*order.OrderItemWithCost, 0)}
	service := NewCostCalculatorService(menuRepo, ingredientRepo, orderRepo, orderItemRepo)

	// Add coffee ingredient
	coffeeID := primitive.NewObjectID()
	ingredientRepo.ingredients[coffeeID] = &ingredient.Ingredient{
		ID:                coffeeID,
		Name:              "Cà phê",
		Unit:              ingredient.UnitGram,
		CostPerUnit:       500.0,
		WastagePercentage: 0.0,
	}

	// Create multi-size menu item with variants
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
				Ingredients: []menu.Ingredient{
					{Name: "Cà phê", Quantity: 20, Unit: ingredient.UnitGram},
					{Name: "Unknown Ingredient", Quantity: 10, Unit: ingredient.UnitGram}, // Missing cost
				},
			},
		},
	}
	menuRepo.Create(context.Background(), menuItem)

	// Calculate cost
	result, err := service.CalculateMenuItemCost(context.Background(), menuItem.ID)
	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}

	// Verify status is INCOMPLETE
	if result.CostStatus != menu.CostStatusIncomplete {
		t.Errorf("Expected cost status INCOMPLETE, got %s", result.CostStatus)
	}

	// Fetch updated menu item to verify variant status
	updated, _ := menuRepo.FindByID(context.Background(), menuItem.ID)
	variantM := updated.GetVariantByID("M")
	if variantM.CostStatus != menu.CostStatusIncomplete {
		t.Errorf("Expected variant M status INCOMPLETE, got %s", variantM.CostStatus)
	}
}

func TestCalculateMenuItemCost_SingleSize_BackwardCompatible(t *testing.T) {
	menuRepo, ingredientRepo := setupTestData()
	orderRepo := &mockOrderRepository{orders: make(map[primitive.ObjectID]*order.Order)}
	orderItemRepo := &mockOrderItemRepository{orderItems: make([]*order.OrderItemWithCost, 0)}
	service := NewCostCalculatorService(menuRepo, ingredientRepo, orderRepo, orderItemRepo)

	// Add coffee ingredient
	coffeeID := primitive.NewObjectID()
	ingredientRepo.ingredients[coffeeID] = &ingredient.Ingredient{
		ID:                coffeeID,
		Name:              "Cà phê",
		Unit:              ingredient.UnitGram,
		CostPerUnit:       500.0,
		WastagePercentage: 0.0,
	}

	// Create single-size menu item (backward compatible)
	menuItem := &menu.MenuItem{
		ID:          primitive.NewObjectID(),
		Name:        "Bánh mì",
		Category:    "Món ăn",
		HasVariants: false,
		Price:       20000,
		Available:   true,
		Ingredients: []menu.Ingredient{
			{Name: "Cà phê", Quantity: 10, Unit: ingredient.UnitGram},
		},
	}
	menuRepo.Create(context.Background(), menuItem)

	// Calculate cost
	result, err := service.CalculateMenuItemCost(context.Background(), menuItem.ID)
	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}

	// Verify cost is calculated for single-size item
	expectedCost := 5000.0 // 10g * 500đ/g
	if result.CurrentCost != expectedCost {
		t.Errorf("Expected cost %f, got %f", expectedCost, result.CurrentCost)
	}
	if result.CostStatus != menu.CostStatusFinal {
		t.Errorf("Expected cost status FINAL, got %s", result.CostStatus)
	}

	// Fetch updated menu item to verify
	updated, _ := menuRepo.FindByID(context.Background(), menuItem.ID)
	if updated.CurrentCost != expectedCost {
		t.Errorf("Expected menu item cost %f, got %f", expectedCost, updated.CurrentCost)
	}
	if len(updated.Variants) != 0 {
		t.Errorf("Expected 0 variants for single-size item, got %d", len(updated.Variants))
	}
}

func TestCalculateMenuItemCost_MultiSize_WithConversionAndWastage(t *testing.T) {
	menuRepo := &mockMenuRepository{menuItems: make(map[primitive.ObjectID]*menu.MenuItem)}
	ingredientRepo := &mockIngredientRepository{ingredients: make(map[primitive.ObjectID]*ingredient.Ingredient)}
	orderRepo := &mockOrderRepository{orders: make(map[primitive.ObjectID]*order.Order)}
	orderItemRepo := &mockOrderItemRepository{orderItems: make([]*order.OrderItemWithCost, 0)}
	service := NewCostCalculatorService(menuRepo, ingredientRepo, orderRepo, orderItemRepo)

	// Create ingredient with conversion and wastage
	ing := &ingredient.Ingredient{
		ID:                primitive.NewObjectID(),
		Name:              "Cà phê",
		Unit:              ingredient.UnitKilogram, // Stock in kg
		CostPerUnit:       500000,                  // 500,000đ per kg
		WastagePercentage: 10.0,                    // 10% wastage
	}
	ingredientRepo.Create(context.Background(), ing)

	// Create multi-size menu item with variants using grams
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
				Ingredients: []menu.Ingredient{
					{Name: "Cà phê", Quantity: 20, Unit: ingredient.UnitGram}, // Recipe in grams
				},
			},
		},
	}
	menuRepo.Create(context.Background(), menuItem)

	// Calculate cost
	_, err := service.CalculateMenuItemCost(context.Background(), menuItem.ID)
	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}

	// Verify cost calculation with conversion and wastage
	// Cost = 20g * (500,000đ/kg) * (1kg/1000g) * (1 + 10/100)
	// Cost = 20 * 500,000 * 0.001 * 1.1 = 11,000đ
	expectedCost := 11000.0
	
	updated, _ := menuRepo.FindByID(context.Background(), menuItem.ID)
	variantM := updated.GetVariantByID("M")
	if variantM.CurrentCost != expectedCost {
		t.Errorf("Expected variant M cost %f (with conversion and wastage), got %f", expectedCost, variantM.CurrentCost)
	}
}
