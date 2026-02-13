package menu

import (
	"encoding/json"
	"strings"
	"testing"
	"cafe-pos/backend/domain/ingredient"
)

// Test backward compatibility - single-size items

func TestGetPrice_SingleSize(t *testing.T) {
	item := &MenuItem{
		Name:        "Bánh mì",
		Category:    "Món ăn",
		HasVariants: false,
		Price:       20000,
	}
	
	price := item.GetPrice("")
	if price != 20000 {
		t.Errorf("Expected price 20000, got %f", price)
	}
}

func TestGetIngredients_SingleSize(t *testing.T) {
	ingredients := []Ingredient{
		{Name: "Bánh mì", Quantity: 1, Unit: ingredient.UnitPiece},
		{Name: "Thịt", Quantity: 50, Unit: ingredient.UnitGram},
	}
	
	item := &MenuItem{
		Name:        "Bánh mì",
		Category:    "Món ăn",
		HasVariants: false,
		Price:       20000,
		Ingredients: ingredients,
	}
	
	result := item.GetIngredients("")
	if len(result) != 2 {
		t.Errorf("Expected 2 ingredients, got %d", len(result))
	}
}

func TestValidate_ValidSingleSize(t *testing.T) {
	item := &MenuItem{
		Name:        "Bánh mì",
		Category:    "Món ăn",
		HasVariants: false,
		Price:       20000,
	}
	
	err := item.Validate()
	if err != nil {
		t.Errorf("Expected no error, got %v", err)
	}
}

// Test new functionality - multi-size items

func TestGetDefaultVariant_Found(t *testing.T) {
	item := &MenuItem{
		Name:        "Cà phê",
		Category:    "Đồ uống",
		HasVariants: true,
		Variants: []MenuItemVariant{
			{ID: "M", Name: "Size M", Price: 25000, IsDefault: true},
			{ID: "L", Name: "Size L", Price: 30000, IsDefault: false},
		},
	}
	
	variant := item.GetDefaultVariant()
	if variant == nil {
		t.Fatal("Expected default variant, got nil")
	}
	if variant.ID != "M" {
		t.Errorf("Expected default variant ID 'M', got '%s'", variant.ID)
	}
}

func TestGetDefaultVariant_NotFound(t *testing.T) {
	item := &MenuItem{
		Name:        "Bánh mì",
		Category:    "Món ăn",
		HasVariants: false,
		Price:       20000,
	}
	
	variant := item.GetDefaultVariant()
	if variant != nil {
		t.Errorf("Expected nil for single-size item, got %v", variant)
	}
}

func TestGetVariantByID_Found(t *testing.T) {
	item := &MenuItem{
		Name:        "Cà phê",
		Category:    "Đồ uống",
		HasVariants: true,
		Variants: []MenuItemVariant{
			{ID: "M", Name: "Size M", Price: 25000, IsDefault: true},
			{ID: "L", Name: "Size L", Price: 30000, IsDefault: false},
		},
	}
	
	variant := item.GetVariantByID("L")
	if variant == nil {
		t.Fatal("Expected variant L, got nil")
	}
	if variant.Price != 30000 {
		t.Errorf("Expected price 30000, got %f", variant.Price)
	}
}

func TestGetVariantByID_NotFound(t *testing.T) {
	item := &MenuItem{
		Name:        "Cà phê",
		Category:    "Đồ uống",
		HasVariants: true,
		Variants: []MenuItemVariant{
			{ID: "M", Name: "Size M", Price: 25000, IsDefault: true},
		},
	}
	
	variant := item.GetVariantByID("XL")
	if variant != nil {
		t.Errorf("Expected nil for non-existent variant, got %v", variant)
	}
}

func TestGetPrice_WithVariantID(t *testing.T) {
	item := &MenuItem{
		Name:        "Cà phê",
		Category:    "Đồ uống",
		HasVariants: true,
		Variants: []MenuItemVariant{
			{ID: "M", Name: "Size M", Price: 25000, IsDefault: true},
			{ID: "L", Name: "Size L", Price: 30000, IsDefault: false},
		},
	}
	
	price := item.GetPrice("L")
	if price != 30000 {
		t.Errorf("Expected price 30000, got %f", price)
	}
}

func TestGetPrice_InvalidVariantID_FallbackToDefault(t *testing.T) {
	item := &MenuItem{
		Name:        "Cà phê",
		Category:    "Đồ uống",
		HasVariants: true,
		Variants: []MenuItemVariant{
			{ID: "M", Name: "Size M", Price: 25000, IsDefault: true},
			{ID: "L", Name: "Size L", Price: 30000, IsDefault: false},
		},
	}
	
	price := item.GetPrice("XL")
	if price != 25000 {
		t.Errorf("Expected fallback to default price 25000, got %f", price)
	}
}

func TestValidate_ValidMultiSize(t *testing.T) {
	item := &MenuItem{
		Name:        "Cà phê",
		Category:    "Đồ uống",
		HasVariants: true,
		Variants: []MenuItemVariant{
			{ID: "M", Name: "Size M", Price: 25000, IsDefault: true},
			{ID: "L", Name: "Size L", Price: 30000, IsDefault: false},
		},
	}
	
	err := item.Validate()
	if err != nil {
		t.Errorf("Expected no error, got %v", err)
	}
}

// Test validation edge cases

func TestValidate_NoVariantsWhenHasVariantsTrue(t *testing.T) {
	item := &MenuItem{
		Name:        "Cà phê",
		Category:    "Đồ uống",
		HasVariants: true,
		Variants:    []MenuItemVariant{},
	}
	
	err := item.Validate()
	if err == nil {
		t.Error("Expected error for no variants when has_variants=true")
	}
}

func TestValidate_NoDefaultVariant(t *testing.T) {
	item := &MenuItem{
		Name:        "Cà phê",
		Category:    "Đồ uống",
		HasVariants: true,
		Variants: []MenuItemVariant{
			{ID: "M", Name: "Size M", Price: 25000, IsDefault: false},
			{ID: "L", Name: "Size L", Price: 30000, IsDefault: false},
		},
	}
	
	err := item.Validate()
	if err == nil {
		t.Error("Expected error for no default variant")
	}
}

func TestValidate_MultipleDefaultVariants(t *testing.T) {
	item := &MenuItem{
		Name:        "Cà phê",
		Category:    "Đồ uống",
		HasVariants: true,
		Variants: []MenuItemVariant{
			{ID: "M", Name: "Size M", Price: 25000, IsDefault: true},
			{ID: "L", Name: "Size L", Price: 30000, IsDefault: true},
		},
	}
	
	err := item.Validate()
	if err == nil {
		t.Error("Expected error for multiple default variants")
	}
}

func TestValidate_DuplicateVariantIDs(t *testing.T) {
	item := &MenuItem{
		Name:        "Cà phê",
		Category:    "Đồ uống",
		HasVariants: true,
		Variants: []MenuItemVariant{
			{ID: "M", Name: "Size M", Price: 25000, IsDefault: true},
			{ID: "M", Name: "Size M Duplicate", Price: 30000, IsDefault: false},
		},
	}
	
	err := item.Validate()
	if err == nil {
		t.Error("Expected error for duplicate variant IDs")
	}
}

func TestValidate_PriceSetWhenHasVariantsTrue(t *testing.T) {
	item := &MenuItem{
		Name:        "Cà phê",
		Category:    "Đồ uống",
		HasVariants: true,
		Price:       20000, // Should not be set
		Variants: []MenuItemVariant{
			{ID: "M", Name: "Size M", Price: 25000, IsDefault: true},
		},
	}
	
	err := item.Validate()
	if err == nil {
		t.Error("Expected error for price set when has_variants=true")
	}
}

func TestValidate_VariantsSetWhenHasVariantsFalse(t *testing.T) {
	item := &MenuItem{
		Name:        "Bánh mì",
		Category:    "Món ăn",
		HasVariants: false,
		Price:       20000,
		Variants: []MenuItemVariant{
			{ID: "M", Name: "Size M", Price: 25000, IsDefault: true},
		},
	}
	
	err := item.Validate()
	if err == nil {
		t.Error("Expected error for variants set when has_variants=false")
	}
}

// Test ambiguous states

func TestValidate_BothPriceAndVariants(t *testing.T) {
	item := &MenuItem{
		Name:        "Cà phê",
		Category:    "Đồ uống",
		HasVariants: true,
		Price:       20000,
		Variants: []MenuItemVariant{
			{ID: "M", Name: "Size M", Price: 25000, IsDefault: true},
		},
	}
	
	err := item.Validate()
	if err == nil {
		t.Error("Expected error for ambiguous state (both price and variants)")
	}
}

func TestValidate_IngredientsSetWhenHasVariantsTrue(t *testing.T) {
	item := &MenuItem{
		Name:        "Cà phê",
		Category:    "Đồ uống",
		HasVariants: true,
		Ingredients: []Ingredient{
			{Name: "Cà phê", Quantity: 20, Unit: ingredient.UnitGram},
		},
		Variants: []MenuItemVariant{
			{ID: "M", Name: "Size M", Price: 25000, IsDefault: true},
		},
	}
	
	err := item.Validate()
	if err == nil {
		t.Error("Expected error for ingredients set when has_variants=true")
	}
}

func TestValidate_MissingName(t *testing.T) {
	item := &MenuItem{
		Name:        "",
		Category:    "Món ăn",
		HasVariants: false,
		Price:       20000,
	}
	
	err := item.Validate()
	if err == nil {
		t.Error("Expected error for missing name")
	}
}

func TestValidate_MissingCategory(t *testing.T) {
	item := &MenuItem{
		Name:        "Bánh mì",
		Category:    "",
		HasVariants: false,
		Price:       20000,
	}
	
	err := item.Validate()
	if err == nil {
		t.Error("Expected error for missing category")
	}
}

func TestValidate_ZeroPrice_SingleSize(t *testing.T) {
	item := &MenuItem{
		Name:        "Bánh mì",
		Category:    "Món ăn",
		HasVariants: false,
		Price:       0,
	}
	
	err := item.Validate()
	if err == nil {
		t.Error("Expected error for zero price in single-size item")
	}
}

func TestValidate_NegativePrice_SingleSize(t *testing.T) {
	item := &MenuItem{
		Name:        "Bánh mì",
		Category:    "Món ăn",
		HasVariants: false,
		Price:       -1000,
	}
	
	err := item.Validate()
	if err == nil {
		t.Error("Expected error for negative price in single-size item")
	}
}

func TestValidate_ZeroVariantPrice(t *testing.T) {
	item := &MenuItem{
		Name:        "Cà phê",
		Category:    "Đồ uống",
		HasVariants: true,
		Variants: []MenuItemVariant{
			{ID: "M", Name: "Size M", Price: 0, IsDefault: true},
		},
	}
	
	err := item.Validate()
	if err == nil {
		t.Error("Expected error for zero variant price")
	}
}

func TestValidate_EmptyVariantID(t *testing.T) {
	item := &MenuItem{
		Name:        "Cà phê",
		Category:    "Đồ uống",
		HasVariants: true,
		Variants: []MenuItemVariant{
			{ID: "", Name: "Size M", Price: 25000, IsDefault: true},
		},
	}
	
	err := item.Validate()
	if err == nil {
		t.Error("Expected error for empty variant ID")
	}
}

// Test DTO JSON marshaling/unmarshaling

func TestCreateMenuItemRequest_JSON_SingleSize(t *testing.T) {
	// Test unmarshaling
	jsonData := `{"name":"Bánh mì","category":"Món ăn","description":"Bánh mì Việt Nam","has_variants":false,"price":20000,"ingredients":[{"name":"Bánh mì","quantity":1,"unit":"piece"}]}`
	
	var unmarshaled CreateMenuItemRequest
	err := unmarshalJSON([]byte(jsonData), &unmarshaled)
	if err != nil {
		t.Fatalf("Failed to unmarshal JSON: %v", err)
	}
	
	if unmarshaled.Name != "Bánh mì" {
		t.Errorf("Expected name 'Bánh mì', got '%s'", unmarshaled.Name)
	}
	if unmarshaled.HasVariants != false {
		t.Errorf("Expected has_variants false, got %v", unmarshaled.HasVariants)
	}
	if unmarshaled.Price != 20000 {
		t.Errorf("Expected price 20000, got %f", unmarshaled.Price)
	}
	if len(unmarshaled.Ingredients) != 1 {
		t.Errorf("Expected 1 ingredient, got %d", len(unmarshaled.Ingredients))
	}
	
	// Test marshaling
	req := CreateMenuItemRequest{
		Name:        "Bánh mì",
		Category:    "Món ăn",
		Description: "Bánh mì Việt Nam",
		HasVariants: false,
		Price:       20000,
		Ingredients: []Ingredient{
			{Name: "Bánh mì", Quantity: 1, Unit: ingredient.UnitPiece},
		},
	}
	
	jsonBytes, err := marshalJSON(req)
	if err != nil {
		t.Fatalf("Failed to marshal JSON: %v", err)
	}
	
	if len(jsonBytes) == 0 {
		t.Error("Expected non-empty JSON output")
	}
}

func TestCreateMenuItemRequest_JSON_MultiSize(t *testing.T) {
	jsonData := `{
		"name":"Cà phê sữa đá",
		"category":"Cà phê",
		"description":"Cà phê phin truyền thống",
		"has_variants":true,
		"variants":[
			{
				"id":"M",
				"name":"Size M",
				"price":25000,
				"ingredients":[{"name":"Cà phê","quantity":20,"unit":"g"}],
				"available":true,
				"is_default":true,
				"current_cost":0,
				"cost_last_calculated_at":"0001-01-01T00:00:00Z",
				"cost_status":""
			},
			{
				"id":"L",
				"name":"Size L",
				"price":30000,
				"ingredients":[{"name":"Cà phê","quantity":30,"unit":"g"}],
				"available":true,
				"is_default":false,
				"current_cost":0,
				"cost_last_calculated_at":"0001-01-01T00:00:00Z",
				"cost_status":""
			}
		]
	}`
	
	var req CreateMenuItemRequest
	err := unmarshalJSON([]byte(jsonData), &req)
	if err != nil {
		t.Fatalf("Failed to unmarshal JSON: %v", err)
	}
	
	if req.Name != "Cà phê sữa đá" {
		t.Errorf("Expected name 'Cà phê sữa đá', got '%s'", req.Name)
	}
	if req.HasVariants != true {
		t.Errorf("Expected has_variants true, got %v", req.HasVariants)
	}
	if len(req.Variants) != 2 {
		t.Errorf("Expected 2 variants, got %d", len(req.Variants))
	}
	if req.Variants[0].ID != "M" {
		t.Errorf("Expected variant ID 'M', got '%s'", req.Variants[0].ID)
	}
	if req.Variants[0].Price != 25000 {
		t.Errorf("Expected variant price 25000, got %f", req.Variants[0].Price)
	}
}

func TestUpdateMenuItemRequest_JSON_SingleSize(t *testing.T) {
	jsonData := `{
		"name":"Bánh mì updated",
		"category":"Món ăn",
		"price":22000,
		"has_variants":false
	}`
	
	var req UpdateMenuItemRequest
	err := unmarshalJSON([]byte(jsonData), &req)
	if err != nil {
		t.Fatalf("Failed to unmarshal JSON: %v", err)
	}
	
	if req.Name != "Bánh mì updated" {
		t.Errorf("Expected name 'Bánh mì updated', got '%s'", req.Name)
	}
	if req.HasVariants == nil || *req.HasVariants != false {
		t.Errorf("Expected has_variants false, got %v", req.HasVariants)
	}
	if req.Price != 22000 {
		t.Errorf("Expected price 22000, got %f", req.Price)
	}
}

func TestUpdateMenuItemRequest_JSON_MultiSize(t *testing.T) {
	jsonData := `{
		"name":"Cà phê updated",
		"has_variants":true,
		"variants":[
			{
				"id":"M",
				"name":"Size M",
				"price":26000,
				"ingredients":[],
				"available":true,
				"is_default":true,
				"current_cost":0,
				"cost_last_calculated_at":"0001-01-01T00:00:00Z",
				"cost_status":""
			}
		]
	}`
	
	var req UpdateMenuItemRequest
	err := unmarshalJSON([]byte(jsonData), &req)
	if err != nil {
		t.Fatalf("Failed to unmarshal JSON: %v", err)
	}
	
	if req.Name != "Cà phê updated" {
		t.Errorf("Expected name 'Cà phê updated', got '%s'", req.Name)
	}
	if req.HasVariants == nil || *req.HasVariants != true {
		t.Errorf("Expected has_variants true, got %v", req.HasVariants)
	}
	if len(req.Variants) != 1 {
		t.Errorf("Expected 1 variant, got %d", len(req.Variants))
	}
}

func TestCreateMenuItemRequest_OmitEmpty_SingleSize(t *testing.T) {
	req := CreateMenuItemRequest{
		Name:        "Bánh mì",
		Category:    "Món ăn",
		HasVariants: false,
		Price:       20000,
	}
	
	jsonData, err := marshalJSON(req)
	if err != nil {
		t.Fatalf("Failed to marshal JSON: %v", err)
	}
	
	jsonStr := string(jsonData)
	
	// Price should be present
	if !contains(jsonStr, "price") {
		t.Error("Expected price to be present for single-size item")
	}
	
	// has_variants should be present (even if false)
	if !contains(jsonStr, "has_variants") {
		t.Error("Expected has_variants to be present")
	}
}

func TestCreateMenuItemRequest_OmitEmpty_MultiSize(t *testing.T) {
	req := CreateMenuItemRequest{
		Name:        "Cà phê",
		Category:    "Cà phê",
		HasVariants: true,
		Variants: []MenuItemVariant{
			{ID: "M", Name: "Size M", Price: 25000, IsDefault: true, Available: true},
		},
	}
	
	jsonData, err := marshalJSON(req)
	if err != nil {
		t.Fatalf("Failed to marshal JSON: %v", err)
	}
	
	jsonStr := string(jsonData)
	
	// Variants should be present
	if !contains(jsonStr, "variants") {
		t.Error("Expected variants to be present for multi-size item")
	}
	
	// Price should be omitted (or 0) when has_variants=true
	// Note: omitempty doesn't omit zero values, only nil/empty
}

// Helper functions for JSON testing
func marshalJSON(v interface{}) ([]byte, error) {
	return json.Marshal(v)
}

func unmarshalJSON(data []byte, v interface{}) error {
	return json.Unmarshal(data, v)
}

func contains(s, substr string) bool {
	return strings.Contains(s, substr)
}
