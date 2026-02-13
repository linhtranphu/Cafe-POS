package http

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"cafe-pos/backend/application/services"
)

// Task 14.1: Test validation errors at API level
// This file tests that validation errors are properly returned via HTTP with clear error messages

// Test 1: API - Create multi-size without variants (should return 400/500 with clear error)
func TestAPIValidationError_CreateMultiSize_NoVariants(t *testing.T) {
	menuRepo := newMockMenuRepositoryForAPI()
	menuService := services.NewMenuService(menuRepo)
	router := setupMenuTestRouter(menuService)

	reqBody := map[string]interface{}{
		"name":         "Cà phê",
		"category":     "Cà phê",
		"description":  "Cà phê phin",
		"has_variants": true,
		"variants":     []map[string]interface{}{}, // Empty - should fail
	}

	body, _ := json.Marshal(reqBody)
	req, _ := http.NewRequest("POST", "/api/menu", bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()

	router.ServeHTTP(w, req)

	// Should return error status (400 or 500)
	if w.Code != http.StatusBadRequest && w.Code != http.StatusInternalServerError {
		t.Errorf("Expected status 400 or 500, got %d", w.Code)
	}

	// Check error message
	responseBody := w.Body.String()
	if !strings.Contains(responseBody, "variants") {
		t.Errorf("Expected error message to mention 'variants', got: %s", responseBody)
	}
}

// Test 2: API - Create multi-size without default variant (should return 400/500 with clear error)
func TestAPIValidationError_CreateMultiSize_NoDefaultVariant(t *testing.T) {
	menuRepo := newMockMenuRepositoryForAPI()
	menuService := services.NewMenuService(menuRepo)
	router := setupMenuTestRouter(menuService)

	reqBody := map[string]interface{}{
		"name":         "Cà phê",
		"category":     "Cà phê",
		"description":  "Cà phê phin",
		"has_variants": true,
		"variants": []map[string]interface{}{
			{
				"id":         "M",
				"name":       "Size M",
				"price":      25000,
				"available":  true,
				"is_default": false, // No default
			},
			{
				"id":         "L",
				"name":       "Size L",
				"price":      30000,
				"available":  true,
				"is_default": false, // No default
			},
		},
	}

	body, _ := json.Marshal(reqBody)
	req, _ := http.NewRequest("POST", "/api/menu", bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()

	router.ServeHTTP(w, req)

	// Should return error status
	if w.Code != http.StatusBadRequest && w.Code != http.StatusInternalServerError {
		t.Errorf("Expected status 400 or 500, got %d", w.Code)
	}

	// Check error message
	responseBody := w.Body.String()
	if !strings.Contains(responseBody, "default") {
		t.Errorf("Expected error message to mention 'default', got: %s", responseBody)
	}
}

// Test 3: API - Create multi-size with duplicate variant IDs (should return 400/500 with clear error)
func TestAPIValidationError_CreateMultiSize_DuplicateVariantIDs(t *testing.T) {
	menuRepo := newMockMenuRepositoryForAPI()
	menuService := services.NewMenuService(menuRepo)
	router := setupMenuTestRouter(menuService)

	reqBody := map[string]interface{}{
		"name":         "Cà phê",
		"category":     "Cà phê",
		"description":  "Cà phê phin",
		"has_variants": true,
		"variants": []map[string]interface{}{
			{
				"id":         "M",
				"name":       "Size M",
				"price":      25000,
				"available":  true,
				"is_default": true,
			},
			{
				"id":         "M", // Duplicate
				"name":       "Size M Duplicate",
				"price":      30000,
				"available":  true,
				"is_default": false,
			},
		},
	}

	body, _ := json.Marshal(reqBody)
	req, _ := http.NewRequest("POST", "/api/menu", bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()

	router.ServeHTTP(w, req)

	// Should return error status
	if w.Code != http.StatusBadRequest && w.Code != http.StatusInternalServerError {
		t.Errorf("Expected status 400 or 500, got %d", w.Code)
	}

	// Check error message
	responseBody := w.Body.String()
	if !strings.Contains(responseBody, "duplicate") {
		t.Errorf("Expected error message to mention 'duplicate', got: %s", responseBody)
	}
}

// Test 4: API - Create item with both price and variants (should return 400/500 with clear error)
func TestAPIValidationError_CreateItem_BothPriceAndVariants(t *testing.T) {
	menuRepo := newMockMenuRepositoryForAPI()
	menuService := services.NewMenuService(menuRepo)
	router := setupMenuTestRouter(menuService)

	reqBody := map[string]interface{}{
		"name":         "Cà phê",
		"category":     "Cà phê",
		"description":  "Cà phê phin",
		"has_variants": true,
		"price":        20000, // Should NOT be set
		"variants": []map[string]interface{}{
			{
				"id":         "M",
				"name":       "Size M",
				"price":      25000,
				"available":  true,
				"is_default": true,
			},
		},
	}

	body, _ := json.Marshal(reqBody)
	req, _ := http.NewRequest("POST", "/api/menu", bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()

	router.ServeHTTP(w, req)

	// Should return error status
	if w.Code != http.StatusBadRequest && w.Code != http.StatusInternalServerError {
		t.Errorf("Expected status 400 or 500, got %d", w.Code)
	}

	// Check error message
	responseBody := w.Body.String()
	if !strings.Contains(responseBody, "price") {
		t.Errorf("Expected error message to mention 'price', got: %s", responseBody)
	}
}

// Test 5: API - Create single-size with variants (should return 400/500 with clear error)
func TestAPIValidationError_CreateSingleSize_WithVariants(t *testing.T) {
	menuRepo := newMockMenuRepositoryForAPI()
	menuService := services.NewMenuService(menuRepo)
	router := setupMenuTestRouter(menuService)

	reqBody := map[string]interface{}{
		"name":         "Bánh mì",
		"category":     "Món ăn",
		"has_variants": false,
		"price":        20000,
		"variants": []map[string]interface{}{ // Should NOT be set
			{
				"id":         "M",
				"name":       "Size M",
				"price":      25000,
				"available":  true,
				"is_default": true,
			},
		},
	}

	body, _ := json.Marshal(reqBody)
	req, _ := http.NewRequest("POST", "/api/menu", bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()

	router.ServeHTTP(w, req)

	// Should return error status
	if w.Code != http.StatusBadRequest && w.Code != http.StatusInternalServerError {
		t.Errorf("Expected status 400 or 500, got %d", w.Code)
	}

	// Check error message
	responseBody := w.Body.String()
	if !strings.Contains(responseBody, "variants") {
		t.Errorf("Expected error message to mention 'variants', got: %s", responseBody)
	}
}

// Test 6: API - Create single-size with zero price (should return 400/500 with clear error)
func TestAPIValidationError_CreateSingleSize_ZeroPrice(t *testing.T) {
	menuRepo := newMockMenuRepositoryForAPI()
	menuService := services.NewMenuService(menuRepo)
	router := setupMenuTestRouter(menuService)

	reqBody := map[string]interface{}{
		"name":         "Bánh mì",
		"category":     "Món ăn",
		"has_variants": false,
		"price":        0, // Invalid
	}

	body, _ := json.Marshal(reqBody)
	req, _ := http.NewRequest("POST", "/api/menu", bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()

	router.ServeHTTP(w, req)

	// Should return error status
	if w.Code != http.StatusBadRequest && w.Code != http.StatusInternalServerError {
		t.Errorf("Expected status 400 or 500, got %d", w.Code)
	}

	// Check error message
	responseBody := w.Body.String()
	if !strings.Contains(responseBody, "price") {
		t.Errorf("Expected error message to mention 'price', got: %s", responseBody)
	}
}

// Test 7: API - Create multi-size with zero variant price (should return 400/500 with clear error)
func TestAPIValidationError_CreateMultiSize_ZeroVariantPrice(t *testing.T) {
	menuRepo := newMockMenuRepositoryForAPI()
	menuService := services.NewMenuService(menuRepo)
	router := setupMenuTestRouter(menuService)

	reqBody := map[string]interface{}{
		"name":         "Cà phê",
		"category":     "Cà phê",
		"has_variants": true,
		"variants": []map[string]interface{}{
			{
				"id":         "M",
				"name":       "Size M",
				"price":      0, // Invalid
				"available":  true,
				"is_default": true,
			},
		},
	}

	body, _ := json.Marshal(reqBody)
	req, _ := http.NewRequest("POST", "/api/menu", bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()

	router.ServeHTTP(w, req)

	// Should return error status
	if w.Code != http.StatusBadRequest && w.Code != http.StatusInternalServerError {
		t.Errorf("Expected status 400 or 500, got %d", w.Code)
	}

	// Check error message
	responseBody := w.Body.String()
	if !strings.Contains(responseBody, "price") {
		t.Errorf("Expected error message to mention 'price', got: %s", responseBody)
	}
}

// Test 8: API - Create multi-size with empty variant ID (should return 400/500 with clear error)
func TestAPIValidationError_CreateMultiSize_EmptyVariantID(t *testing.T) {
	menuRepo := newMockMenuRepositoryForAPI()
	menuService := services.NewMenuService(menuRepo)
	router := setupMenuTestRouter(menuService)

	reqBody := map[string]interface{}{
		"name":         "Cà phê",
		"category":     "Cà phê",
		"has_variants": true,
		"variants": []map[string]interface{}{
			{
				"id":         "", // Empty - invalid
				"name":       "Size M",
				"price":      25000,
				"available":  true,
				"is_default": true,
			},
		},
	}

	body, _ := json.Marshal(reqBody)
	req, _ := http.NewRequest("POST", "/api/menu", bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()

	router.ServeHTTP(w, req)

	// Should return error status
	if w.Code != http.StatusBadRequest && w.Code != http.StatusInternalServerError {
		t.Errorf("Expected status 400 or 500, got %d", w.Code)
	}

	// Check error message
	responseBody := w.Body.String()
	if !strings.Contains(responseBody, "variant ID") {
		t.Errorf("Expected error message to mention 'variant ID', got: %s", responseBody)
	}
}

// Test 9: API - Create multi-size with multiple defaults (should return 400/500 with clear error)
func TestAPIValidationError_CreateMultiSize_MultipleDefaults(t *testing.T) {
	menuRepo := newMockMenuRepositoryForAPI()
	menuService := services.NewMenuService(menuRepo)
	router := setupMenuTestRouter(menuService)

	reqBody := map[string]interface{}{
		"name":         "Cà phê",
		"category":     "Cà phê",
		"has_variants": true,
		"variants": []map[string]interface{}{
			{
				"id":         "M",
				"name":       "Size M",
				"price":      25000,
				"available":  true,
				"is_default": true, // Default
			},
			{
				"id":         "L",
				"name":       "Size L",
				"price":      30000,
				"available":  true,
				"is_default": true, // Also default - invalid
			},
		},
	}

	body, _ := json.Marshal(reqBody)
	req, _ := http.NewRequest("POST", "/api/menu", bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()

	router.ServeHTTP(w, req)

	// Should return error status
	if w.Code != http.StatusBadRequest && w.Code != http.StatusInternalServerError {
		t.Errorf("Expected status 400 or 500, got %d", w.Code)
	}

	// Check error message
	responseBody := w.Body.String()
	if !strings.Contains(responseBody, "default") {
		t.Errorf("Expected error message to mention 'default', got: %s", responseBody)
	}
}

// Test 10: API - Create multi-size with ingredients at item level (should return 400/500 with clear error)
func TestAPIValidationError_CreateMultiSize_WithIngredientsAtItemLevel(t *testing.T) {
	menuRepo := newMockMenuRepositoryForAPI()
	menuService := services.NewMenuService(menuRepo)
	router := setupMenuTestRouter(menuService)

	reqBody := map[string]interface{}{
		"name":         "Cà phê",
		"category":     "Cà phê",
		"has_variants": true,
		"ingredients": []map[string]interface{}{ // Should NOT be set
			{"name": "Cà phê", "quantity": 20, "unit": "g"},
		},
		"variants": []map[string]interface{}{
			{
				"id":         "M",
				"name":       "Size M",
				"price":      25000,
				"available":  true,
				"is_default": true,
			},
		},
	}

	body, _ := json.Marshal(reqBody)
	req, _ := http.NewRequest("POST", "/api/menu", bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()

	router.ServeHTTP(w, req)

	// Should return error status
	if w.Code != http.StatusBadRequest && w.Code != http.StatusInternalServerError {
		t.Errorf("Expected status 400 or 500, got %d", w.Code)
	}

	// Check error message
	responseBody := w.Body.String()
	if !strings.Contains(responseBody, "ingredients") {
		t.Errorf("Expected error message to mention 'ingredients', got: %s", responseBody)
	}
}

// Note: Order validation errors (missing variant_id, invalid variant_id) are tested
// in the order handler tests and integration tests. They require more complex setup
// with shifts and menu items, so they're better tested in the order service tests.
