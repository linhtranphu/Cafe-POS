package http

import (
	"bytes"
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"cafe-pos/backend/application/services"
	"cafe-pos/backend/domain/menu"
	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

// Setup test router
func setupMenuTestRouter(menuService *services.MenuService) *gin.Engine {
	gin.SetMode(gin.TestMode)
	router := gin.Default()
	handler := NewMenuHandler(menuService)

	router.POST("/api/menu", handler.CreateMenuItem)
	router.GET("/api/menu", handler.GetAllMenuItems)
	router.GET("/api/menu/:id", handler.GetMenuItem)
	router.PUT("/api/menu/:id", handler.UpdateMenuItem)
	router.DELETE("/api/menu/:id", handler.DeleteMenuItem)

	return router
}

// Test POST /api/menu - single-size item

func TestMenuAPI_CreateSingleSize(t *testing.T) {
	menuRepo := newMockMenuRepositoryForAPI()
	menuService := services.NewMenuService(menuRepo)
	router := setupMenuTestRouter(menuService)

	reqBody := map[string]interface{}{
		"name":         "Bánh mì",
		"category":     "Món ăn",
		"description":  "Bánh mì Việt Nam",
		"has_variants": false,
		"price":        20000,
		"ingredients": []map[string]interface{}{
			{"name": "Bánh mì", "quantity": 1, "unit": "piece"},
		},
	}

	body, _ := json.Marshal(reqBody)
	req, _ := http.NewRequest("POST", "/api/menu", bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()

	router.ServeHTTP(w, req)

	if w.Code != http.StatusCreated {
		t.Errorf("Expected status 201, got %d", w.Code)
	}

	var response menu.MenuItem
	json.Unmarshal(w.Body.Bytes(), &response)

	if response.Name != "Bánh mì" {
		t.Errorf("Expected name 'Bánh mì', got '%s'", response.Name)
	}
	if response.HasVariants != false {
		t.Errorf("Expected has_variants false, got %v", response.HasVariants)
	}
	if response.Price != 20000 {
		t.Errorf("Expected price 20000, got %f", response.Price)
	}
}

// Test POST /api/menu - multi-size item

func TestMenuAPI_CreateMultiSize(t *testing.T) {
	menuRepo := newMockMenuRepositoryForAPI()
	menuService := services.NewMenuService(menuRepo)
	router := setupMenuTestRouter(menuService)

	reqBody := map[string]interface{}{
		"name":         "Cà phê sữa đá",
		"category":     "Cà phê",
		"description":  "Cà phê phin truyền thống",
		"has_variants": true,
		"variants": []map[string]interface{}{
			{
				"id":         "M",
				"name":       "Size M",
				"price":      25000,
				"available":  true,
				"is_default": true,
				"ingredients": []map[string]interface{}{
					{"name": "Cà phê", "quantity": 20, "unit": "g"},
				},
			},
			{
				"id":         "L",
				"name":       "Size L",
				"price":      30000,
				"available":  true,
				"is_default": false,
				"ingredients": []map[string]interface{}{
					{"name": "Cà phê", "quantity": 30, "unit": "g"},
				},
			},
		},
	}

	body, _ := json.Marshal(reqBody)
	req, _ := http.NewRequest("POST", "/api/menu", bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()

	router.ServeHTTP(w, req)

	if w.Code != http.StatusCreated {
		t.Errorf("Expected status 201, got %d. Body: %s", w.Code, w.Body.String())
	}

	var response menu.MenuItem
	json.Unmarshal(w.Body.Bytes(), &response)

	if response.Name != "Cà phê sữa đá" {
		t.Errorf("Expected name 'Cà phê sữa đá', got '%s'", response.Name)
	}
	if response.HasVariants != true {
		t.Errorf("Expected has_variants true, got %v", response.HasVariants)
	}
	if len(response.Variants) != 2 {
		t.Errorf("Expected 2 variants, got %d", len(response.Variants))
	}
}

// Test POST /api/menu - validation error (no variants when has_variants=true)

func TestMenuAPI_CreateMultiSize_NoVariants(t *testing.T) {
	menuRepo := newMockMenuRepositoryForAPI()
	menuService := services.NewMenuService(menuRepo)
	router := setupMenuTestRouter(menuService)

	reqBody := map[string]interface{}{
		"name":         "Cà phê",
		"category":     "Cà phê",
		"has_variants": true,
		"variants":     []map[string]interface{}{}, // Empty
	}

	body, _ := json.Marshal(reqBody)
	req, _ := http.NewRequest("POST", "/api/menu", bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()

	router.ServeHTTP(w, req)

	if w.Code != http.StatusInternalServerError {
		t.Errorf("Expected status 500 (validation error), got %d", w.Code)
	}
}

// Test GET /api/menu - list items with variants

func TestMenuAPI_GetAllItems(t *testing.T) {
	menuRepo := newMockMenuRepositoryForAPI()
	menuService := services.NewMenuService(menuRepo)
	router := setupMenuTestRouter(menuService)

	// Create single-size item
	singleItem := &menu.MenuItem{
		ID:          primitive.NewObjectID(),
		Name:        "Bánh mì",
		Category:    "Món ăn",
		HasVariants: false,
		Price:       20000,
		Available:   true,
	}
	menuRepo.items[singleItem.ID] = singleItem

	// Create multi-size item
	multiItem := &menu.MenuItem{
		ID:          primitive.NewObjectID(),
		Name:        "Cà phê",
		Category:    "Cà phê",
		HasVariants: true,
		Available:   true,
		Variants: []menu.MenuItemVariant{
			{ID: "M", Name: "Size M", Price: 25000, IsDefault: true, Available: true},
			{ID: "L", Name: "Size L", Price: 30000, IsDefault: false, Available: true},
		},
	}
	menuRepo.items[multiItem.ID] = multiItem

	req, _ := http.NewRequest("GET", "/api/menu", nil)
	w := httptest.NewRecorder()

	router.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Errorf("Expected status 200, got %d", w.Code)
	}

	var response []*menu.MenuItem
	json.Unmarshal(w.Body.Bytes(), &response)

	if len(response) != 2 {
		t.Errorf("Expected 2 items, got %d", len(response))
	}
}

// Test PUT /api/menu/:id - toggle single to multi

func TestMenuAPI_UpdateSingleToMulti(t *testing.T) {
	menuRepo := newMockMenuRepositoryForAPI()
	menuService := services.NewMenuService(menuRepo)
	router := setupMenuTestRouter(menuService)

	// Create initial single-size item
	item := &menu.MenuItem{
		ID:          primitive.NewObjectID(),
		Name:        "Bánh mì",
		Category:    "Món ăn",
		HasVariants: false,
		Price:       20000,
		Available:   true,
	}
	menuRepo.items[item.ID] = item

	// Update to multi-size
	reqBody := map[string]interface{}{
		"has_variants": true,
		"variants": []map[string]interface{}{
			{
				"id":         "M",
				"name":       "Size M",
				"price":      20000,
				"available":  true,
				"is_default": true,
			},
			{
				"id":         "L",
				"name":       "Size L",
				"price":      25000,
				"available":  true,
				"is_default": false,
			},
		},
	}

	body, _ := json.Marshal(reqBody)
	req, _ := http.NewRequest("PUT", "/api/menu/"+item.ID.Hex(), bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()

	router.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Errorf("Expected status 200, got %d. Body: %s", w.Code, w.Body.String())
	}

	var response menu.MenuItem
	json.Unmarshal(w.Body.Bytes(), &response)

	if response.HasVariants != true {
		t.Errorf("Expected has_variants true, got %v", response.HasVariants)
	}
	if len(response.Variants) != 2 {
		t.Errorf("Expected 2 variants, got %d", len(response.Variants))
	}
	if response.Price != 0 {
		t.Errorf("Expected price 0 (cleared), got %f", response.Price)
	}
}

// Mock repository for API tests
type mockMenuRepositoryForAPI struct {
	items map[primitive.ObjectID]*menu.MenuItem
}

func newMockMenuRepositoryForAPI() *mockMenuRepositoryForAPI {
	return &mockMenuRepositoryForAPI{
		items: make(map[primitive.ObjectID]*menu.MenuItem),
	}
}

func (m *mockMenuRepositoryForAPI) Create(ctx context.Context, item *menu.MenuItem) error {
	if item.ID.IsZero() {
		item.ID = primitive.NewObjectID()
	}
	m.items[item.ID] = item
	return nil
}

func (m *mockMenuRepositoryForAPI) FindAll(ctx context.Context) ([]*menu.MenuItem, error) {
	var items []*menu.MenuItem
	for _, item := range m.items {
		items = append(items, item)
	}
	return items, nil
}

func (m *mockMenuRepositoryForAPI) FindByID(ctx context.Context, id primitive.ObjectID) (*menu.MenuItem, error) {
	item, exists := m.items[id]
	if !exists {
		return nil, nil
	}
	return item, nil
}

func (m *mockMenuRepositoryForAPI) FindByCategory(ctx context.Context, category string) ([]*menu.MenuItem, error) {
	return nil, nil
}

func (m *mockMenuRepositoryForAPI) FindByIngredientName(ctx context.Context, ingredientName string) ([]*menu.MenuItem, error) {
	return nil, nil
}

func (m *mockMenuRepositoryForAPI) Update(ctx context.Context, id primitive.ObjectID, item *menu.MenuItem) error {
	if _, exists := m.items[id]; !exists {
		return nil
	}
	m.items[id] = item
	return nil
}

func (m *mockMenuRepositoryForAPI) Delete(ctx context.Context, id primitive.ObjectID) error {
	delete(m.items, id)
	return nil
}
