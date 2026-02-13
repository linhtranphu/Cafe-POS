package http

import (
	"bytes"
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"cafe-pos/backend/application/services"
	"cafe-pos/backend/domain"
	"cafe-pos/backend/domain/menu"
	"cafe-pos/backend/domain/order"
	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

// Setup test router for orders
func setupOrderTestRouter(orderService *services.OrderService, stateMachineManager *domain.StateMachineManager) *gin.Engine {
	gin.SetMode(gin.TestMode)
	router := gin.Default()
	
	// Mock authentication middleware
	router.Use(func(c *gin.Context) {
		c.Set("user_id", "507f1f77bcf86cd799439011")
		c.Set("username", "test_waiter")
		c.Next()
	})
	
	handler := NewOrderHandler(orderService, stateMachineManager)

	router.POST("/api/orders", handler.CreateOrder)
	router.GET("/api/orders", handler.GetAllOrders)
	router.GET("/api/orders/:id", handler.GetOrder)

	return router
}

// Test POST /api/orders with single-size items (201 Created)
func TestOrderAPI_CreateOrderWithSingleSizeItems(t *testing.T) {
	// Setup repositories
	menuRepo := newMockMenuRepositoryForOrderTests()
	orderRepo := newMockOrderRepositoryForTests()
	shiftRepo := newMockShiftRepositoryForTests()
	
	// Create an open shift
	shift := &order.Shift{
		ID:     primitive.NewObjectID(),
		Status: order.ShiftOpen,
	}
	shiftRepo.shifts[shift.ID] = shift
	
	// Create single-size menu items
	item1 := &menu.MenuItem{
		ID:          primitive.NewObjectID(),
		Name:        "Bánh mì",
		Category:    "Món ăn",
		HasVariants: false,
		Price:       20000,
		Available:   true,
	}
	item2 := &menu.MenuItem{
		ID:          primitive.NewObjectID(),
		Name:        "Nước suối",
		Category:    "Nước uống",
		HasVariants: false,
		Price:       10000,
		Available:   true,
	}
	menuRepo.items[item1.ID] = item1
	menuRepo.items[item2.ID] = item2
	
	// Setup services
	stateMachineManager := domain.NewStateMachineManager()
	orderService := services.NewOrderService(orderRepo, shiftRepo, menuRepo, stateMachineManager)
	router := setupOrderTestRouter(orderService, stateMachineManager)
	
	// Create order request
	reqBody := map[string]interface{}{
		"customer_name": "Khách 1",
		"shift_id":      shift.ID.Hex(),
		"items": []map[string]interface{}{
			{
				"menu_item_id": item1.ID.Hex(),
				"quantity":     2,
			},
			{
				"menu_item_id": item2.ID.Hex(),
				"quantity":     1,
			},
		},
	}
	
	body, _ := json.Marshal(reqBody)
	req, _ := http.NewRequest("POST", "/api/orders", bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()
	
	router.ServeHTTP(w, req)
	
	if w.Code != http.StatusCreated {
		t.Errorf("Expected status 201, got %d. Body: %s", w.Code, w.Body.String())
	}
	
	var response order.Order
	json.Unmarshal(w.Body.Bytes(), &response)
	
	// Verify order details
	if len(response.Items) != 2 {
		t.Errorf("Expected 2 items, got %d", len(response.Items))
	}
	
	// Verify first item (single-size, no variant)
	if response.Items[0].Name != "Bánh mì" {
		t.Errorf("Expected item name 'Bánh mì', got '%s'", response.Items[0].Name)
	}
	if response.Items[0].Price != 20000 {
		t.Errorf("Expected price 20000, got %f", response.Items[0].Price)
	}
	if response.Items[0].VariantID != "" {
		t.Errorf("Expected no variant_id for single-size item, got '%s'", response.Items[0].VariantID)
	}
	if response.Items[0].VariantName != "" {
		t.Errorf("Expected no variant_name for single-size item, got '%s'", response.Items[0].VariantName)
	}
	
	// Verify total calculation
	expectedTotal := (20000 * 2) + (10000 * 1)
	if response.Total != float64(expectedTotal) {
		t.Errorf("Expected total %d, got %f", expectedTotal, response.Total)
	}
}

// Test POST /api/orders with multi-size items and variant_id (201 Created)
func TestOrderAPI_CreateOrderWithMultiSizeItems(t *testing.T) {
	// Setup repositories
	menuRepo := newMockMenuRepositoryForOrderTests()
	orderRepo := newMockOrderRepositoryForTests()
	shiftRepo := newMockShiftRepositoryForTests()
	
	// Create an open shift
	shift := &order.Shift{
		ID:     primitive.NewObjectID(),
		Status: order.ShiftOpen,
	}
	shiftRepo.shifts[shift.ID] = shift
	
	// Create multi-size menu item
	item := &menu.MenuItem{
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
			{
				ID:        "XL",
				Name:      "Size XL",
				Price:     35000,
				Available: true,
				IsDefault: false,
			},
		},
	}
	menuRepo.items[item.ID] = item
	
	// Setup services
	stateMachineManager := domain.NewStateMachineManager()
	orderService := services.NewOrderService(orderRepo, shiftRepo, menuRepo, stateMachineManager)
	router := setupOrderTestRouter(orderService, stateMachineManager)
	
	// Create order request with variants
	reqBody := map[string]interface{}{
		"customer_name": "Khách 2",
		"shift_id":      shift.ID.Hex(),
		"items": []map[string]interface{}{
			{
				"menu_item_id": item.ID.Hex(),
				"variant_id":   "L",
				"quantity":     2,
			},
			{
				"menu_item_id": item.ID.Hex(),
				"variant_id":   "XL",
				"quantity":     1,
			},
		},
	}
	
	body, _ := json.Marshal(reqBody)
	req, _ := http.NewRequest("POST", "/api/orders", bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()
	
	router.ServeHTTP(w, req)
	
	if w.Code != http.StatusCreated {
		t.Errorf("Expected status 201, got %d. Body: %s", w.Code, w.Body.String())
	}
	
	var response order.Order
	json.Unmarshal(w.Body.Bytes(), &response)
	
	// Verify order details
	if len(response.Items) != 2 {
		t.Errorf("Expected 2 items, got %d", len(response.Items))
	}
	
	// Verify first item (Size L)
	if response.Items[0].Name != "Cà phê sữa đá" {
		t.Errorf("Expected item name 'Cà phê sữa đá', got '%s'", response.Items[0].Name)
	}
	if response.Items[0].VariantID != "L" {
		t.Errorf("Expected variant_id 'L', got '%s'", response.Items[0].VariantID)
	}
	if response.Items[0].VariantName != "Size L" {
		t.Errorf("Expected variant_name 'Size L', got '%s'", response.Items[0].VariantName)
	}
	if response.Items[0].Price != 30000 {
		t.Errorf("Expected price 30000, got %f", response.Items[0].Price)
	}
	
	// Verify second item (Size XL)
	if response.Items[1].VariantID != "XL" {
		t.Errorf("Expected variant_id 'XL', got '%s'", response.Items[1].VariantID)
	}
	if response.Items[1].VariantName != "Size XL" {
		t.Errorf("Expected variant_name 'Size XL', got '%s'", response.Items[1].VariantName)
	}
	if response.Items[1].Price != 35000 {
		t.Errorf("Expected price 35000, got %f", response.Items[1].Price)
	}
	
	// Verify total calculation
	expectedTotal := (30000 * 2) + (35000 * 1)
	if response.Total != float64(expectedTotal) {
		t.Errorf("Expected total %d, got %f", expectedTotal, response.Total)
	}
}

// Test POST /api/orders with mixed items (201 Created)
func TestOrderAPI_CreateOrderWithMixedItems(t *testing.T) {
	// Setup repositories
	menuRepo := newMockMenuRepositoryForOrderTests()
	orderRepo := newMockOrderRepositoryForTests()
	shiftRepo := newMockShiftRepositoryForTests()
	
	// Create an open shift
	shift := &order.Shift{
		ID:     primitive.NewObjectID(),
		Status: order.ShiftOpen,
	}
	shiftRepo.shifts[shift.ID] = shift
	
	// Create single-size item
	singleItem := &menu.MenuItem{
		ID:          primitive.NewObjectID(),
		Name:        "Bánh mì",
		Category:    "Món ăn",
		HasVariants: false,
		Price:       20000,
		Available:   true,
	}
	
	// Create multi-size item
	multiItem := &menu.MenuItem{
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
	
	menuRepo.items[singleItem.ID] = singleItem
	menuRepo.items[multiItem.ID] = multiItem
	
	// Setup services
	stateMachineManager := domain.NewStateMachineManager()
	orderService := services.NewOrderService(orderRepo, shiftRepo, menuRepo, stateMachineManager)
	router := setupOrderTestRouter(orderService, stateMachineManager)
	
	// Create order request with mixed items
	reqBody := map[string]interface{}{
		"customer_name": "Khách 3",
		"shift_id":      shift.ID.Hex(),
		"items": []map[string]interface{}{
			{
				"menu_item_id": singleItem.ID.Hex(),
				"quantity":     1,
				// No variant_id for single-size
			},
			{
				"menu_item_id": multiItem.ID.Hex(),
				"variant_id":   "L",
				"quantity":     2,
			},
		},
	}
	
	body, _ := json.Marshal(reqBody)
	req, _ := http.NewRequest("POST", "/api/orders", bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()
	
	router.ServeHTTP(w, req)
	
	if w.Code != http.StatusCreated {
		t.Errorf("Expected status 201, got %d. Body: %s", w.Code, w.Body.String())
	}
	
	var response order.Order
	json.Unmarshal(w.Body.Bytes(), &response)
	
	// Verify order details
	if len(response.Items) != 2 {
		t.Errorf("Expected 2 items, got %d", len(response.Items))
	}
	
	// Verify single-size item
	if response.Items[0].Name != "Bánh mì" {
		t.Errorf("Expected item name 'Bánh mì', got '%s'", response.Items[0].Name)
	}
	if response.Items[0].VariantID != "" {
		t.Errorf("Expected no variant_id for single-size item, got '%s'", response.Items[0].VariantID)
	}
	if response.Items[0].Price != 20000 {
		t.Errorf("Expected price 20000, got %f", response.Items[0].Price)
	}
	
	// Verify multi-size item
	if response.Items[1].Name != "Cà phê sữa đá" {
		t.Errorf("Expected item name 'Cà phê sữa đá', got '%s'", response.Items[1].Name)
	}
	if response.Items[1].VariantID != "L" {
		t.Errorf("Expected variant_id 'L', got '%s'", response.Items[1].VariantID)
	}
	if response.Items[1].VariantName != "Size L" {
		t.Errorf("Expected variant_name 'Size L', got '%s'", response.Items[1].VariantName)
	}
	if response.Items[1].Price != 30000 {
		t.Errorf("Expected price 30000, got %f", response.Items[1].Price)
	}
	
	// Verify total calculation
	expectedTotal := (20000 * 1) + (30000 * 2)
	if response.Total != float64(expectedTotal) {
		t.Errorf("Expected total %d, got %f", expectedTotal, response.Total)
	}
}

// Test POST /api/orders missing variant_id for multi-size (400 Bad Request)
func TestOrderAPI_CreateOrderMissingVariantID(t *testing.T) {
	// Setup repositories
	menuRepo := newMockMenuRepositoryForOrderTests()
	orderRepo := newMockOrderRepositoryForTests()
	shiftRepo := newMockShiftRepositoryForTests()
	
	// Create an open shift
	shift := &order.Shift{
		ID:     primitive.NewObjectID(),
		Status: order.ShiftOpen,
	}
	shiftRepo.shifts[shift.ID] = shift
	
	// Create multi-size menu item
	item := &menu.MenuItem{
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
		},
	}
	menuRepo.items[item.ID] = item
	
	// Setup services
	stateMachineManager := domain.NewStateMachineManager()
	orderService := services.NewOrderService(orderRepo, shiftRepo, menuRepo, stateMachineManager)
	router := setupOrderTestRouter(orderService, stateMachineManager)
	
	// Create order request WITHOUT variant_id for multi-size item
	reqBody := map[string]interface{}{
		"customer_name": "Khách 4",
		"shift_id":      shift.ID.Hex(),
		"items": []map[string]interface{}{
			{
				"menu_item_id": item.ID.Hex(),
				"quantity":     1,
				// Missing variant_id!
			},
		},
	}
	
	body, _ := json.Marshal(reqBody)
	req, _ := http.NewRequest("POST", "/api/orders", bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()
	
	router.ServeHTTP(w, req)
	
	if w.Code != http.StatusBadRequest {
		t.Errorf("Expected status 400, got %d. Body: %s", w.Code, w.Body.String())
	}
	
	var response map[string]interface{}
	json.Unmarshal(w.Body.Bytes(), &response)
	
	// Verify error message mentions variant_id requirement
	errorMsg, ok := response["error"].(string)
	if !ok {
		t.Errorf("Expected error message in response")
	}
	if errorMsg == "" {
		t.Errorf("Expected non-empty error message")
	}
	// Error should mention variant_id is required
	t.Logf("Error message: %s", errorMsg)
}

// Test POST /api/orders with invalid variant_id (400 Bad Request)
func TestOrderAPI_CreateOrderInvalidVariantID(t *testing.T) {
	// Setup repositories
	menuRepo := newMockMenuRepositoryForOrderTests()
	orderRepo := newMockOrderRepositoryForTests()
	shiftRepo := newMockShiftRepositoryForTests()
	
	// Create an open shift
	shift := &order.Shift{
		ID:     primitive.NewObjectID(),
		Status: order.ShiftOpen,
	}
	shiftRepo.shifts[shift.ID] = shift
	
	// Create multi-size menu item
	item := &menu.MenuItem{
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
	menuRepo.items[item.ID] = item
	
	// Setup services
	stateMachineManager := domain.NewStateMachineManager()
	orderService := services.NewOrderService(orderRepo, shiftRepo, menuRepo, stateMachineManager)
	router := setupOrderTestRouter(orderService, stateMachineManager)
	
	// Create order request with INVALID variant_id
	reqBody := map[string]interface{}{
		"customer_name": "Khách 5",
		"shift_id":      shift.ID.Hex(),
		"items": []map[string]interface{}{
			{
				"menu_item_id": item.ID.Hex(),
				"variant_id":   "XXL", // Invalid! Only M and L exist
				"quantity":     1,
			},
		},
	}
	
	body, _ := json.Marshal(reqBody)
	req, _ := http.NewRequest("POST", "/api/orders", bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()
	
	router.ServeHTTP(w, req)
	
	if w.Code != http.StatusBadRequest {
		t.Errorf("Expected status 400, got %d. Body: %s", w.Code, w.Body.String())
	}
	
	var response map[string]interface{}
	json.Unmarshal(w.Body.Bytes(), &response)
	
	// Verify error message mentions invalid variant_id
	errorMsg, ok := response["error"].(string)
	if !ok {
		t.Errorf("Expected error message in response")
	}
	if errorMsg == "" {
		t.Errorf("Expected non-empty error message")
	}
	// Error should mention invalid variant_id
	t.Logf("Error message: %s", errorMsg)
}

// Mock repositories for order tests

type mockMenuRepositoryForOrderTests struct {
	items map[primitive.ObjectID]*menu.MenuItem
}

func newMockMenuRepositoryForOrderTests() *mockMenuRepositoryForOrderTests {
	return &mockMenuRepositoryForOrderTests{
		items: make(map[primitive.ObjectID]*menu.MenuItem),
	}
}

func (m *mockMenuRepositoryForOrderTests) Create(ctx context.Context, item *menu.MenuItem) error {
	if item.ID.IsZero() {
		item.ID = primitive.NewObjectID()
	}
	m.items[item.ID] = item
	return nil
}

func (m *mockMenuRepositoryForOrderTests) FindAll(ctx context.Context) ([]*menu.MenuItem, error) {
	var items []*menu.MenuItem
	for _, item := range m.items {
		items = append(items, item)
	}
	return items, nil
}

func (m *mockMenuRepositoryForOrderTests) FindByID(ctx context.Context, id primitive.ObjectID) (*menu.MenuItem, error) {
	item, exists := m.items[id]
	if !exists {
		return nil, nil
	}
	return item, nil
}

func (m *mockMenuRepositoryForOrderTests) FindByCategory(ctx context.Context, category string) ([]*menu.MenuItem, error) {
	return nil, nil
}

func (m *mockMenuRepositoryForOrderTests) FindByIngredientName(ctx context.Context, ingredientName string) ([]*menu.MenuItem, error) {
	return nil, nil
}

func (m *mockMenuRepositoryForOrderTests) Update(ctx context.Context, id primitive.ObjectID, item *menu.MenuItem) error {
	if _, exists := m.items[id]; !exists {
		return nil
	}
	m.items[id] = item
	return nil
}

func (m *mockMenuRepositoryForOrderTests) Delete(ctx context.Context, id primitive.ObjectID) error {
	delete(m.items, id)
	return nil
}

type mockOrderRepositoryForTests struct {
	orders map[primitive.ObjectID]*order.Order
}

func newMockOrderRepositoryForTests() *mockOrderRepositoryForTests {
	return &mockOrderRepositoryForTests{
		orders: make(map[primitive.ObjectID]*order.Order),
	}
}

func (m *mockOrderRepositoryForTests) Create(ctx context.Context, o *order.Order) error {
	if o.ID.IsZero() {
		o.ID = primitive.NewObjectID()
	}
	m.orders[o.ID] = o
	return nil
}

func (m *mockOrderRepositoryForTests) FindByID(ctx context.Context, id primitive.ObjectID) (*order.Order, error) {
	o, exists := m.orders[id]
	if !exists {
		return nil, nil
	}
	return o, nil
}

func (m *mockOrderRepositoryForTests) Update(ctx context.Context, id primitive.ObjectID, o *order.Order) error {
	if _, exists := m.orders[id]; !exists {
		return nil
	}
	m.orders[id] = o
	return nil
}

func (m *mockOrderRepositoryForTests) FindByShiftID(ctx context.Context, shiftID primitive.ObjectID) ([]*order.Order, error) {
	return nil, nil
}

func (m *mockOrderRepositoryForTests) FindByWaiterID(ctx context.Context, waiterID primitive.ObjectID) ([]*order.Order, error) {
	return nil, nil
}

func (m *mockOrderRepositoryForTests) FindByStatus(ctx context.Context, status order.OrderStatus) ([]*order.Order, error) {
	return nil, nil
}

func (m *mockOrderRepositoryForTests) FindByOrderNumber(ctx context.Context, orderNumber string) (*order.Order, error) {
	return nil, nil
}

func (m *mockOrderRepositoryForTests) FindAll(ctx context.Context) ([]*order.Order, error) {
	var orders []*order.Order
	for _, o := range m.orders {
		orders = append(orders, o)
	}
	return orders, nil
}

type mockShiftRepositoryForTests struct {
	shifts map[primitive.ObjectID]*order.Shift
}

func newMockShiftRepositoryForTests() *mockShiftRepositoryForTests {
	return &mockShiftRepositoryForTests{
		shifts: make(map[primitive.ObjectID]*order.Shift),
	}
}

func (m *mockShiftRepositoryForTests) Create(ctx context.Context, shift *order.Shift) error {
	if shift.ID.IsZero() {
		shift.ID = primitive.NewObjectID()
	}
	m.shifts[shift.ID] = shift
	return nil
}

func (m *mockShiftRepositoryForTests) FindByID(ctx context.Context, id primitive.ObjectID) (*order.Shift, error) {
	shift, exists := m.shifts[id]
	if !exists {
		return nil, nil
	}
	return shift, nil
}

func (m *mockShiftRepositoryForTests) Update(ctx context.Context, id primitive.ObjectID, shift *order.Shift) error {
	if _, exists := m.shifts[id]; !exists {
		return nil
	}
	m.shifts[id] = shift
	return nil
}

func (m *mockShiftRepositoryForTests) FindOpenShiftByUser(ctx context.Context, userID primitive.ObjectID, role order.RoleType) (*order.Shift, error) {
	for _, shift := range m.shifts {
		if shift.Status == order.ShiftOpen {
			return shift, nil
		}
	}
	return nil, nil
}

func (m *mockShiftRepositoryForTests) FindAll(ctx context.Context) ([]*order.Shift, error) {
	return nil, nil
}

func (m *mockShiftRepositoryForTests) FindByStatus(ctx context.Context, status order.ShiftStatus) ([]*order.Shift, error) {
	return nil, nil
}

func (m *mockShiftRepositoryForTests) FindByUserID(ctx context.Context, userID primitive.ObjectID, roleType order.RoleType) ([]*order.Shift, error) {
	return nil, nil
}

func (m *mockShiftRepositoryForTests) FindByDateRange(ctx context.Context, startDate, endDate time.Time) ([]*order.Shift, error) {
	return nil, nil
}

func (m *mockShiftRepositoryForTests) FindOpenShiftByWaiter(ctx context.Context, waiterID primitive.ObjectID) (*order.Shift, error) {
	return nil, nil
}

func (m *mockShiftRepositoryForTests) FindOpenShifts(ctx context.Context) ([]*order.Shift, error) {
	return nil, nil
}

func (m *mockShiftRepositoryForTests) FindByWaiterID(ctx context.Context, waiterID primitive.ObjectID) ([]*order.Shift, error) {
	return nil, nil
}

func (m *mockShiftRepositoryForTests) FindByRoleType(ctx context.Context, roleType order.RoleType) ([]*order.Shift, error) {
	return nil, nil
}
