package http

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"cafe-pos/backend/application/services"
	"cafe-pos/backend/domain/ingredient"
	"cafe-pos/backend/domain/menu"
	"cafe-pos/backend/domain/order"
	"cafe-pos/backend/domain/settings"
	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

// Mock repositories for testing
type mockMenuRepository struct {
	menuItems []*menu.MenuItem
}

func (m *mockMenuRepository) Create(ctx context.Context, item *menu.MenuItem) error {
	m.menuItems = append(m.menuItems, item)
	return nil
}

func (m *mockMenuRepository) FindAll(ctx context.Context) ([]*menu.MenuItem, error) {
	return m.menuItems, nil
}

func (m *mockMenuRepository) FindByID(ctx context.Context, id primitive.ObjectID) (*menu.MenuItem, error) {
	for _, item := range m.menuItems {
		if item.ID == id {
			return item, nil
		}
	}
	return nil, fmt.Errorf("menu item not found")
}

func (m *mockMenuRepository) FindByCategory(ctx context.Context, category string) ([]*menu.MenuItem, error) {
	var filtered []*menu.MenuItem
	for _, item := range m.menuItems {
		if item.Category == category {
			filtered = append(filtered, item)
		}
	}
	return filtered, nil
}

func (m *mockMenuRepository) FindByIngredientName(ctx context.Context, ingredientName string) ([]*menu.MenuItem, error) {
	var filtered []*menu.MenuItem
	for _, item := range m.menuItems {
		for _, ing := range item.Ingredients {
			if ing.Name == ingredientName {
				filtered = append(filtered, item)
				break
			}
		}
	}
	return filtered, nil
}

func (m *mockMenuRepository) Update(ctx context.Context, id primitive.ObjectID, item *menu.MenuItem) error {
	for i, existing := range m.menuItems {
		if existing.ID == id {
			m.menuItems[i] = item
			return nil
		}
	}
	return nil
}

func (m *mockMenuRepository) Delete(ctx context.Context, id primitive.ObjectID) error {
	return nil
}

type mockIngredientRepository struct {
	ingredients []*ingredient.Ingredient
}

func (m *mockIngredientRepository) Create(ctx context.Context, ing *ingredient.Ingredient) error {
	m.ingredients = append(m.ingredients, ing)
	return nil
}

func (m *mockIngredientRepository) FindAll(ctx context.Context) ([]*ingredient.Ingredient, error) {
	return m.ingredients, nil
}

func (m *mockIngredientRepository) FindByID(ctx context.Context, id primitive.ObjectID) (*ingredient.Ingredient, error) {
	return nil, nil
}

func (m *mockIngredientRepository) Update(ctx context.Context, id primitive.ObjectID, ing *ingredient.Ingredient) error {
	return nil
}

func (m *mockIngredientRepository) Delete(ctx context.Context, id primitive.ObjectID) error {
	return nil
}

func (m *mockIngredientRepository) FindLowStock(ctx context.Context) ([]*ingredient.Ingredient, error) {
	return nil, nil
}

func (m *mockIngredientRepository) CreateCategory(ctx context.Context, category *ingredient.IngredientCategory) error {
	return nil
}

func (m *mockIngredientRepository) FindAllCategories(ctx context.Context) ([]*ingredient.IngredientCategory, error) {
	return nil, nil
}

func (m *mockIngredientRepository) GetCategories(ctx context.Context) ([]ingredient.IngredientCategory, error) {
	return nil, nil
}

func (m *mockIngredientRepository) DeleteCategory(ctx context.Context, id primitive.ObjectID) error {
	return nil
}

type mockOrderItemRepository struct {
	orderItems []*order.OrderItemWithCost
}

func (m *mockOrderItemRepository) CreateMany(ctx context.Context, items []*order.OrderItemWithCost) error {
	m.orderItems = append(m.orderItems, items...)
	return nil
}

func (m *mockOrderItemRepository) FindByOrderIDs(ctx context.Context, orderIDs []primitive.ObjectID) ([]*order.OrderItemWithCost, error) {
	return m.orderItems, nil
}

func (m *mockOrderItemRepository) FindByDateRange(ctx context.Context, startDate, endDate time.Time) ([]*order.OrderItemWithCost, error) {
	return m.orderItems, nil
}

type mockSettingsRepository struct {
	settings *settings.ShopSettings
}

func (m *mockSettingsRepository) FindFirst(ctx context.Context) (*settings.ShopSettings, error) {
	if m.settings == nil {
		m.settings = &settings.ShopSettings{
			LowMarginThreshold: 20.0,
		}
	}
	return m.settings, nil
}

func (m *mockSettingsRepository) Update(ctx context.Context, id primitive.ObjectID, s *settings.ShopSettings) error {
	m.settings = s
	return nil
}

type mockOrderRepository struct{}

func (m *mockOrderRepository) Create(ctx context.Context, o *order.Order) error {
	return nil
}

func (m *mockOrderRepository) FindAll(ctx context.Context) ([]*order.Order, error) {
	return nil, nil
}

func (m *mockOrderRepository) FindByID(ctx context.Context, id primitive.ObjectID) (*order.Order, error) {
	return nil, nil
}

func (m *mockOrderRepository) Update(ctx context.Context, id primitive.ObjectID, o *order.Order) error {
	return nil
}

func (m *mockOrderRepository) FindByShiftID(ctx context.Context, shiftID primitive.ObjectID) ([]*order.Order, error) {
	return nil, nil
}

func (m *mockOrderRepository) FindByWaiterID(ctx context.Context, waiterID primitive.ObjectID) ([]*order.Order, error) {
	return nil, nil
}

func (m *mockOrderRepository) FindByBaristaID(ctx context.Context, baristaID primitive.ObjectID) ([]*order.Order, error) {
	return nil, nil
}

func (m *mockOrderRepository) FindByStatus(ctx context.Context, status order.OrderStatus) ([]*order.Order, error) {
	return nil, nil
}

func (m *mockOrderRepository) FindByOrderNumber(ctx context.Context, orderNumber string) (*order.Order, error) {
	return nil, nil
}

// Setup test handler
func setupTestHandler() (*MenuCostHandler, *mockMenuRepository, *mockIngredientRepository) {
	gin.SetMode(gin.TestMode)

	menuRepo := &mockMenuRepository{menuItems: []*menu.MenuItem{}}
	ingredientRepo := &mockIngredientRepository{ingredients: []*ingredient.Ingredient{}}
	orderItemRepo := &mockOrderItemRepository{orderItems: []*order.OrderItemWithCost{}}
	settingsRepo := &mockSettingsRepository{}
	orderRepo := &mockOrderRepository{}

	// Create services
	costCalculator := services.NewCostCalculatorService(menuRepo, ingredientRepo, orderRepo, orderItemRepo)
	profitAnalyzer := services.NewProfitAnalyzerService(menuRepo, orderItemRepo, settingsRepo, nil)
	recalcService := services.NewCostRecalculationService(costCalculator, menuRepo, 2, 100)

	// Create handler
	handler := NewMenuCostHandler(profitAnalyzer, costCalculator, recalcService)

	return handler, menuRepo, ingredientRepo
}

// Test GET /api/menu/costs with various filters
func TestGetMenuCosts_WithFilters(t *testing.T) {
	handler, menuRepo, ingredientRepo := setupTestHandler()

	// Setup test data
	coffeeID := primitive.NewObjectID()
	teaID := primitive.NewObjectID()
	espressoID := primitive.NewObjectID()
	milkID := primitive.NewObjectID()

	// Create ingredients
	ingredientRepo.ingredients = []*ingredient.Ingredient{
		{
			ID:          espressoID,
			Name:        "Espresso",
			CostPerUnit: 200,
			Unit:        ingredient.UnitMilliliter,
		},
		{
			ID:          milkID,
			Name:        "Milk",
			CostPerUnit: 50,
			Unit:        ingredient.UnitMilliliter,
		},
	}

	// Create menu items
	menuRepo.menuItems = []*menu.MenuItem{
		{
			ID:       coffeeID,
			Name:     "Cappuccino",
			Category: "Coffee",
			Price:    45000,
			Ingredients: []menu.Ingredient{
				{Name: "Espresso", Quantity: 30, Unit: string(ingredient.UnitMilliliter)},
				{Name: "Milk", Quantity: 150, Unit: string(ingredient.UnitMilliliter)},
			},
			CurrentCost:          13500,
			CostStatus:           menu.CostStatusFinal,
			CostLastCalculatedAt: time.Now(),
		},
		{
			ID:       teaID,
			Name:     "Green Tea",
			Category: "Tea",
			Price:    30000,
			Ingredients: []menu.Ingredient{
				{Name: "Tea Leaves", Quantity: 5, Unit: string(ingredient.UnitGram)},
			},
			CurrentCost:          0,
			CostStatus:           menu.CostStatusIncomplete,
			CostLastCalculatedAt: time.Now(),
		},
	}

	// Test 1: Get all menu costs without filter
	t.Run("GetAllMenuCosts", func(t *testing.T) {
		w := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(w)
		c.Request, _ = http.NewRequest("GET", "/api/menu/costs", nil)

		handler.GetMenuCosts(c)

		assert.Equal(t, http.StatusOK, w.Code)

		var response GetMenuCostsResponse
		err := json.Unmarshal(w.Body.Bytes(), &response)
		assert.NoError(t, err)
		assert.Equal(t, 2, len(response.Items))
		assert.Equal(t, 2, response.Summary.TotalItems)
	})

	// Test 2: Filter by category
	t.Run("FilterByCategory", func(t *testing.T) {
		w := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(w)
		c.Request, _ = http.NewRequest("GET", "/api/menu/costs?category=Coffee", nil)

		handler.GetMenuCosts(c)

		assert.Equal(t, http.StatusOK, w.Code)

		var response GetMenuCostsResponse
		err := json.Unmarshal(w.Body.Bytes(), &response)
		assert.NoError(t, err)
		assert.Equal(t, 1, len(response.Items))
		assert.Equal(t, "Cappuccino", response.Items[0].Name)
	})

	// Test 3: Sort by profit_margin descending
	t.Run("SortByProfitMargin", func(t *testing.T) {
		w := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(w)
		c.Request, _ = http.NewRequest("GET", "/api/menu/costs?sort_by=profit_margin&sort_order=desc", nil)

		handler.GetMenuCosts(c)

		assert.Equal(t, http.StatusOK, w.Code)

		var response GetMenuCostsResponse
		err := json.Unmarshal(w.Body.Bytes(), &response)
		assert.NoError(t, err)
		assert.Equal(t, 2, len(response.Items))
		// First item should have higher profit margin
		if len(response.Items) >= 2 {
			assert.GreaterOrEqual(t, response.Items[0].ProfitMargin, response.Items[1].ProfitMargin)
		}
	})
}

// Test GET /api/menu/costs/:id with valid and invalid IDs
func TestGetMenuCostDetail(t *testing.T) {
	handler, menuRepo, ingredientRepo := setupTestHandler()

	// Setup test data
	coffeeID := primitive.NewObjectID()
	espressoID := primitive.NewObjectID()
	milkID := primitive.NewObjectID()

	// Create ingredients
	ingredientRepo.ingredients = []*ingredient.Ingredient{
		{
			ID:                espressoID,
			Name:              "Espresso",
			CostPerUnit:       200,
			Unit:              ingredient.UnitMilliliter,
			ConversionRate:    1.0,
			WastagePercentage: 5.0,
		},
		{
			ID:                milkID,
			Name:              "Milk",
			CostPerUnit:       50,
			Unit:              ingredient.UnitMilliliter,
			ConversionRate:    1.0,
			WastagePercentage: 10.0,
		},
	}

	// Create menu item
	menuRepo.menuItems = []*menu.MenuItem{
		{
			ID:       coffeeID,
			Name:     "Cappuccino",
			Category: "Coffee",
			Price:    45000,
			Ingredients: []menu.Ingredient{
				{Name: "Espresso", Quantity: 30, Unit: string(ingredient.UnitMilliliter)},
				{Name: "Milk", Quantity: 150, Unit: string(ingredient.UnitMilliliter)},
			},
			CurrentCost:          13500,
			CostStatus:           menu.CostStatusFinal,
			CostLastCalculatedAt: time.Now(),
		},
	}

	// Test 1: Valid ID
	t.Run("ValidID", func(t *testing.T) {
		w := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(w)
		c.Request, _ = http.NewRequest("GET", "/api/menu/costs/"+coffeeID.Hex(), nil)
		c.Params = gin.Params{{Key: "id", Value: coffeeID.Hex()}}

		handler.GetMenuCostDetail(c)

		assert.Equal(t, http.StatusOK, w.Code)

		var response GetMenuCostDetailResponse
		err := json.Unmarshal(w.Body.Bytes(), &response)
		assert.NoError(t, err)
		assert.Equal(t, "Cappuccino", response.MenuItem.Name)
		assert.Equal(t, 2, len(response.Ingredients))
		assert.Greater(t, response.TotalCost, 0.0)
	})

	// Test 2: Invalid ID format
	t.Run("InvalidIDFormat", func(t *testing.T) {
		w := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(w)
		c.Request, _ = http.NewRequest("GET", "/api/menu/costs/invalid", nil)
		c.Params = gin.Params{{Key: "id", Value: "invalid"}}

		handler.GetMenuCostDetail(c)

		assert.Equal(t, http.StatusBadRequest, w.Code)
	})

	// Test 3: Non-existent ID
	t.Run("NonExistentID", func(t *testing.T) {
		nonExistentID := primitive.NewObjectID()
		w := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(w)
		c.Request, _ = http.NewRequest("GET", "/api/menu/costs/"+nonExistentID.Hex(), nil)
		c.Params = gin.Params{{Key: "id", Value: nonExistentID.Hex()}}

		handler.GetMenuCostDetail(c)

		assert.Equal(t, http.StatusNotFound, w.Code)
	})
}

// Test GET /api/menu/warnings with custom threshold
func TestGetMenuWarnings(t *testing.T) {
	handler, menuRepo, ingredientRepo := setupTestHandler()

	// Setup test data
	lossItemID := primitive.NewObjectID()
	lowMarginItemID := primitive.NewObjectID()
	profitableItemID := primitive.NewObjectID()

	// Create ingredients
	ingredientRepo.ingredients = []*ingredient.Ingredient{
		{
			ID:          primitive.NewObjectID(),
			Name:        "Ingredient A",
			CostPerUnit: 100,
			Unit:        ingredient.UnitGram,
		},
	}

	// Create menu items with different profit scenarios
	menuRepo.menuItems = []*menu.MenuItem{
		{
			ID:       lossItemID,
			Name:     "Loss Item",
			Category: "Coffee",
			Price:    20000,
			Ingredients: []menu.Ingredient{
				{Name: "Ingredient A", Quantity: 300, Unit: string(ingredient.UnitGram)},
			},
			CurrentCost:          30000, // Cost > Price (loss)
			CostStatus:           menu.CostStatusFinal,
			CostLastCalculatedAt: time.Now(),
		},
		{
			ID:       lowMarginItemID,
			Name:     "Low Margin Item",
			Category: "Coffee",
			Price:    30000,
			Ingredients: []menu.Ingredient{
				{Name: "Ingredient A", Quantity: 250, Unit: string(ingredient.UnitGram)},
			},
			CurrentCost:          25000, // Margin = 16.67% (low margin)
			CostStatus:           menu.CostStatusFinal,
			CostLastCalculatedAt: time.Now(),
		},
		{
			ID:       profitableItemID,
			Name:     "Profitable Item",
			Category: "Coffee",
			Price:    50000,
			Ingredients: []menu.Ingredient{
				{Name: "Ingredient A", Quantity: 100, Unit: string(ingredient.UnitGram)},
			},
			CurrentCost:          10000, // Margin = 80% (profitable)
			CostStatus:           menu.CostStatusFinal,
			CostLastCalculatedAt: time.Now(),
		},
	}

	// Test 1: Default threshold (from settings)
	t.Run("DefaultThreshold", func(t *testing.T) {
		w := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(w)
		c.Request, _ = http.NewRequest("GET", "/api/menu/warnings", nil)

		handler.GetMenuWarnings(c)

		assert.Equal(t, http.StatusOK, w.Code)

		var response services.ProfitWarnings
		err := json.Unmarshal(w.Body.Bytes(), &response)
		assert.NoError(t, err)
		assert.Equal(t, 1, response.LossCount)
		assert.Equal(t, 1, response.LowMarginCount)
		assert.Equal(t, 20.0, response.Threshold)
	})

	// Test 2: Custom threshold
	t.Run("CustomThreshold", func(t *testing.T) {
		w := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(w)
		c.Request, _ = http.NewRequest("GET", "/api/menu/warnings?threshold=15", nil)

		handler.GetMenuWarnings(c)

		assert.Equal(t, http.StatusOK, w.Code)

		var response services.ProfitWarnings
		err := json.Unmarshal(w.Body.Bytes(), &response)
		assert.NoError(t, err)
		assert.Equal(t, 1, response.LossCount)
		assert.Equal(t, 0, response.LowMarginCount) // Not low margin with 15% threshold (item has 16.67%)
		assert.Equal(t, 15.0, response.Threshold)
	})

	// Test 3: Invalid threshold
	t.Run("InvalidThreshold", func(t *testing.T) {
		w := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(w)
		c.Request, _ = http.NewRequest("GET", "/api/menu/warnings?threshold=invalid", nil)

		handler.GetMenuWarnings(c)

		assert.Equal(t, http.StatusBadRequest, w.Code)
	})
}
