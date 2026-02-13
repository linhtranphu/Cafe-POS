package services

import (
	"context"
	"testing"

	"cafe-pos/backend/domain/ingredient"
	"cafe-pos/backend/domain/menu"
)

// Task 14.2: Test edge cases
// - Item with 1 variant
// - Item with 10 variants (max)
// - Toggle has_variants multiple times
// - Delete default variant (should reassign or fail)

// Test: Item with 1 variant (minimum valid multi-size item)
func TestEdgeCase_ItemWith1Variant(t *testing.T) {
	repo := newMockMenuRepositoryForMenuService()
	service := NewMenuService(repo)
	ctx := context.Background()

	req := &menu.CreateMenuItemRequest{
		Name:        "Cà phê đen",
		Category:    "Cà phê",
		Description: "Cà phê đen đá",
		HasVariants: true,
		Variants: []menu.MenuItemVariant{
			{
				ID:        "ONLY",
				Name:      "Size duy nhất",
				Price:     20000,
				Available: true,
				IsDefault: true,
				Ingredients: []menu.Ingredient{
					{Name: "Cà phê", Quantity: 20, Unit: ingredient.UnitGram},
				},
			},
		},
	}

	item, err := service.CreateMenuItem(ctx, req)
	if err != nil {
		t.Fatalf("Expected no error for item with 1 variant, got %v", err)
	}

	// Verify item was created successfully
	if item.Name != "Cà phê đen" {
		t.Errorf("Expected name 'Cà phê đen', got '%s'", item.Name)
	}
	if !item.HasVariants {
		t.Error("Expected has_variants to be true")
	}
	if len(item.Variants) != 1 {
		t.Errorf("Expected 1 variant, got %d", len(item.Variants))
	}
	if !item.Variants[0].IsDefault {
		t.Error("Expected the only variant to be default")
	}

	// Verify we can get the default variant
	defaultVariant := item.GetDefaultVariant()
	if defaultVariant == nil {
		t.Fatal("Expected to get default variant, got nil")
	}
	if defaultVariant.ID != "ONLY" {
		t.Errorf("Expected default variant ID 'ONLY', got '%s'", defaultVariant.ID)
	}

	// Verify we can get price
	price := item.GetPrice("")
	if price != 20000 {
		t.Errorf("Expected price 20000, got %f", price)
	}
}

// Test: Item with 10 variants (maximum recommended)
func TestEdgeCase_ItemWith10Variants(t *testing.T) {
	repo := newMockMenuRepositoryForMenuService()
	service := NewMenuService(repo)
	ctx := context.Background()

	// Create 10 variants
	variants := make([]menu.MenuItemVariant, 10)
	for i := 0; i < 10; i++ {
		variants[i] = menu.MenuItemVariant{
			ID:        string(rune('A' + i)), // A, B, C, ..., J
			Name:      "Size " + string(rune('A'+i)),
			Price:     float64(20000 + i*5000),
			Available: true,
			IsDefault: i == 0, // First one is default
			Ingredients: []menu.Ingredient{
				{Name: "Cà phê", Quantity: float64(10 + i*5), Unit: ingredient.UnitGram},
			},
		}
	}

	req := &menu.CreateMenuItemRequest{
		Name:        "Cà phê đặc biệt",
		Category:    "Cà phê",
		Description: "Cà phê với nhiều size",
		HasVariants: true,
		Variants:    variants,
	}

	item, err := service.CreateMenuItem(ctx, req)
	if err != nil {
		t.Fatalf("Expected no error for item with 10 variants, got %v", err)
	}

	// Verify all variants were created
	if len(item.Variants) != 10 {
		t.Errorf("Expected 10 variants, got %d", len(item.Variants))
	}

	// Verify default variant
	defaultVariant := item.GetDefaultVariant()
	if defaultVariant == nil {
		t.Fatal("Expected to get default variant, got nil")
	}
	if defaultVariant.ID != "A" {
		t.Errorf("Expected default variant ID 'A', got '%s'", defaultVariant.ID)
	}

	// Verify we can get each variant by ID
	for i := 0; i < 10; i++ {
		variantID := string(rune('A' + i))
		variant := item.GetVariantByID(variantID)
		if variant == nil {
			t.Errorf("Expected to find variant '%s', got nil", variantID)
			continue
		}
		if variant.ID != variantID {
			t.Errorf("Expected variant ID '%s', got '%s'", variantID, variant.ID)
		}
		expectedPrice := float64(20000 + i*5000)
		if variant.Price != expectedPrice {
			t.Errorf("Expected variant %s price %f, got %f", variantID, expectedPrice, variant.Price)
		}
	}

	// Verify we can get price for each variant
	for i := 0; i < 10; i++ {
		variantID := string(rune('A' + i))
		price := item.GetPrice(variantID)
		expectedPrice := float64(20000 + i*5000)
		if price != expectedPrice {
			t.Errorf("Expected price %f for variant %s, got %f", expectedPrice, variantID, price)
		}
	}
}

// Test: Toggle has_variants multiple times
func TestEdgeCase_ToggleHasVariantsMultipleTimes(t *testing.T) {
	repo := newMockMenuRepositoryForMenuService()
	service := NewMenuService(repo)
	ctx := context.Background()

	// 1. Create as single-size
	createReq := &menu.CreateMenuItemRequest{
		Name:        "Bánh mì",
		Category:    "Món ăn",
		HasVariants: false,
		Price:       20000,
		Ingredients: []menu.Ingredient{
			{Name: "Bánh mì", Quantity: 1, Unit: ingredient.UnitPiece},
		},
	}
	item, err := service.CreateMenuItem(ctx, createReq)
	if err != nil {
		t.Fatalf("Failed to create initial single-size item: %v", err)
	}

	// Verify initial state
	if item.HasVariants {
		t.Error("Expected has_variants to be false initially")
	}
	if item.Price != 20000 {
		t.Errorf("Expected price 20000, got %f", item.Price)
	}

	// 2. Toggle to multi-size (first toggle)
	hasVariants := true
	updateReq1 := &menu.UpdateMenuItemRequest{
		HasVariants: &hasVariants,
		Variants: []menu.MenuItemVariant{
			{ID: "M", Name: "Size M", Price: 20000, IsDefault: true, Available: true},
			{ID: "L", Name: "Size L", Price: 25000, IsDefault: false, Available: true},
		},
	}
	item, err = service.UpdateMenuItem(ctx, item.ID, updateReq1)
	if err != nil {
		t.Fatalf("Failed to toggle to multi-size (1st toggle): %v", err)
	}

	// Verify after first toggle
	if !item.HasVariants {
		t.Error("Expected has_variants to be true after 1st toggle")
	}
	if len(item.Variants) != 2 {
		t.Errorf("Expected 2 variants after 1st toggle, got %d", len(item.Variants))
	}
	if item.Price != 0 {
		t.Errorf("Expected price to be cleared (0) after 1st toggle, got %f", item.Price)
	}

	// 3. Toggle back to single-size (second toggle)
	hasVariants = false
	updateReq2 := &menu.UpdateMenuItemRequest{
		HasVariants: &hasVariants,
		Price:       22000,
		Ingredients: []menu.Ingredient{
			{Name: "Bánh mì", Quantity: 1, Unit: ingredient.UnitPiece},
		},
	}
	item, err = service.UpdateMenuItem(ctx, item.ID, updateReq2)
	if err != nil {
		t.Fatalf("Failed to toggle back to single-size (2nd toggle): %v", err)
	}

	// Verify after second toggle
	if item.HasVariants {
		t.Error("Expected has_variants to be false after 2nd toggle")
	}
	if item.Price != 22000 {
		t.Errorf("Expected price 22000 after 2nd toggle, got %f", item.Price)
	}
	if len(item.Variants) != 0 {
		t.Errorf("Expected variants to be cleared after 2nd toggle, got %d", len(item.Variants))
	}

	// 4. Toggle to multi-size again (third toggle)
	hasVariants = true
	updateReq3 := &menu.UpdateMenuItemRequest{
		HasVariants: &hasVariants,
		Variants: []menu.MenuItemVariant{
			{ID: "S", Name: "Size S", Price: 18000, IsDefault: true, Available: true},
			{ID: "M", Name: "Size M", Price: 22000, IsDefault: false, Available: true},
			{ID: "L", Name: "Size L", Price: 26000, IsDefault: false, Available: true},
		},
	}
	item, err = service.UpdateMenuItem(ctx, item.ID, updateReq3)
	if err != nil {
		t.Fatalf("Failed to toggle to multi-size again (3rd toggle): %v", err)
	}

	// Verify after third toggle
	if !item.HasVariants {
		t.Error("Expected has_variants to be true after 3rd toggle")
	}
	if len(item.Variants) != 3 {
		t.Errorf("Expected 3 variants after 3rd toggle, got %d", len(item.Variants))
	}
	if item.Price != 0 {
		t.Errorf("Expected price to be cleared (0) after 3rd toggle, got %f", item.Price)
	}

	// 5. Toggle back to single-size one more time (fourth toggle)
	hasVariants = false
	updateReq4 := &menu.UpdateMenuItemRequest{
		HasVariants: &hasVariants,
		Price:       24000,
	}
	item, err = service.UpdateMenuItem(ctx, item.ID, updateReq4)
	if err != nil {
		t.Fatalf("Failed to toggle back to single-size (4th toggle): %v", err)
	}

	// Verify after fourth toggle
	if item.HasVariants {
		t.Error("Expected has_variants to be false after 4th toggle")
	}
	if item.Price != 24000 {
		t.Errorf("Expected price 24000 after 4th toggle, got %f", item.Price)
	}
	if len(item.Variants) != 0 {
		t.Errorf("Expected variants to be cleared after 4th toggle, got %d", len(item.Variants))
	}

	// Verify item is still functional
	finalPrice := item.GetPrice("")
	if finalPrice != 24000 {
		t.Errorf("Expected final price 24000, got %f", finalPrice)
	}
}

// Test: Delete default variant (should fail - must always have exactly one default)
func TestEdgeCase_DeleteDefaultVariant_ShouldFail(t *testing.T) {
	repo := newMockMenuRepositoryForMenuService()
	service := NewMenuService(repo)
	ctx := context.Background()

	// Create multi-size item with 3 variants
	createReq := &menu.CreateMenuItemRequest{
		Name:        "Cà phê",
		Category:    "Cà phê",
		HasVariants: true,
		Variants: []menu.MenuItemVariant{
			{ID: "M", Name: "Size M", Price: 25000, IsDefault: true, Available: true},
			{ID: "L", Name: "Size L", Price: 30000, IsDefault: false, Available: true},
			{ID: "XL", Name: "Size XL", Price: 35000, IsDefault: false, Available: true},
		},
	}
	item, err := service.CreateMenuItem(ctx, createReq)
	if err != nil {
		t.Fatalf("Failed to create item: %v", err)
	}

	// Try to update by removing the default variant (M)
	// This should fail because we'd have no default variant
	hasVariants := true
	updateReq := &menu.UpdateMenuItemRequest{
		HasVariants: &hasVariants,
		Variants: []menu.MenuItemVariant{
			// Removed M (the default)
			{ID: "L", Name: "Size L", Price: 30000, IsDefault: false, Available: true},
			{ID: "XL", Name: "Size XL", Price: 35000, IsDefault: false, Available: true},
		},
	}
	_, err = service.UpdateMenuItem(ctx, item.ID, updateReq)
	if err == nil {
		t.Error("Expected error when removing default variant without reassigning, got nil")
	}

	// Verify error message mentions default variant
	if err != nil && !contains(err.Error(), "default") {
		t.Errorf("Expected error message to mention 'default', got: %v", err)
	}
}

// Test: Delete default variant but reassign to another (should succeed)
func TestEdgeCase_DeleteDefaultVariant_WithReassignment(t *testing.T) {
	repo := newMockMenuRepositoryForMenuService()
	service := NewMenuService(repo)
	ctx := context.Background()

	// Create multi-size item with 3 variants
	createReq := &menu.CreateMenuItemRequest{
		Name:        "Cà phê",
		Category:    "Cà phê",
		HasVariants: true,
		Variants: []menu.MenuItemVariant{
			{ID: "M", Name: "Size M", Price: 25000, IsDefault: true, Available: true},
			{ID: "L", Name: "Size L", Price: 30000, IsDefault: false, Available: true},
			{ID: "XL", Name: "Size XL", Price: 35000, IsDefault: false, Available: true},
		},
	}
	item, err := service.CreateMenuItem(ctx, createReq)
	if err != nil {
		t.Fatalf("Failed to create item: %v", err)
	}

	// Update by removing M (old default) and making L the new default
	hasVariants := true
	updateReq := &menu.UpdateMenuItemRequest{
		HasVariants: &hasVariants,
		Variants: []menu.MenuItemVariant{
			{ID: "L", Name: "Size L", Price: 30000, IsDefault: true, Available: true}, // Now default
			{ID: "XL", Name: "Size XL", Price: 35000, IsDefault: false, Available: true},
		},
	}
	updatedItem, err := service.UpdateMenuItem(ctx, item.ID, updateReq)
	if err != nil {
		t.Fatalf("Expected no error when reassigning default, got %v", err)
	}

	// Verify update succeeded
	if len(updatedItem.Variants) != 2 {
		t.Errorf("Expected 2 variants after update, got %d", len(updatedItem.Variants))
	}

	// Verify L is now the default
	defaultVariant := updatedItem.GetDefaultVariant()
	if defaultVariant == nil {
		t.Fatal("Expected to get default variant, got nil")
	}
	if defaultVariant.ID != "L" {
		t.Errorf("Expected default variant to be 'L', got '%s'", defaultVariant.ID)
	}

	// Verify M no longer exists
	mVariant := updatedItem.GetVariantByID("M")
	if mVariant != nil {
		t.Error("Expected M variant to be deleted, but it still exists")
	}

	// Verify we can still get price (should use new default)
	price := updatedItem.GetPrice("")
	if price != 30000 {
		t.Errorf("Expected price 30000 (new default), got %f", price)
	}
}

// Test: Reduce variants from many to 1 (edge case)
func TestEdgeCase_ReduceVariantsToOne(t *testing.T) {
	repo := newMockMenuRepositoryForMenuService()
	service := NewMenuService(repo)
	ctx := context.Background()

	// Create item with 5 variants
	createReq := &menu.CreateMenuItemRequest{
		Name:        "Cà phê",
		Category:    "Cà phê",
		HasVariants: true,
		Variants: []menu.MenuItemVariant{
			{ID: "XS", Name: "Size XS", Price: 15000, IsDefault: false, Available: true},
			{ID: "S", Name: "Size S", Price: 20000, IsDefault: false, Available: true},
			{ID: "M", Name: "Size M", Price: 25000, IsDefault: true, Available: true},
			{ID: "L", Name: "Size L", Price: 30000, IsDefault: false, Available: true},
			{ID: "XL", Name: "Size XL", Price: 35000, IsDefault: false, Available: true},
		},
	}
	item, err := service.CreateMenuItem(ctx, createReq)
	if err != nil {
		t.Fatalf("Failed to create item: %v", err)
	}

	// Verify initial state
	if len(item.Variants) != 5 {
		t.Errorf("Expected 5 variants initially, got %d", len(item.Variants))
	}

	// Update to keep only 1 variant
	hasVariants := true
	updateReq := &menu.UpdateMenuItemRequest{
		HasVariants: &hasVariants,
		Variants: []menu.MenuItemVariant{
			{ID: "M", Name: "Size M", Price: 25000, IsDefault: true, Available: true},
		},
	}
	updatedItem, err := service.UpdateMenuItem(ctx, item.ID, updateReq)
	if err != nil {
		t.Fatalf("Expected no error when reducing to 1 variant, got %v", err)
	}

	// Verify only 1 variant remains
	if len(updatedItem.Variants) != 1 {
		t.Errorf("Expected 1 variant after update, got %d", len(updatedItem.Variants))
	}

	// Verify it's the correct variant
	if updatedItem.Variants[0].ID != "M" {
		t.Errorf("Expected remaining variant to be 'M', got '%s'", updatedItem.Variants[0].ID)
	}

	// Verify it's still the default
	if !updatedItem.Variants[0].IsDefault {
		t.Error("Expected remaining variant to be default")
	}

	// Verify item is still functional
	defaultVariant := updatedItem.GetDefaultVariant()
	if defaultVariant == nil {
		t.Fatal("Expected to get default variant, got nil")
	}
	if defaultVariant.ID != "M" {
		t.Errorf("Expected default variant 'M', got '%s'", defaultVariant.ID)
	}
}

// Helper function for string contains check
func contains(s, substr string) bool {
	return len(s) >= len(substr) && (s == substr || len(substr) == 0 || 
		(len(s) > 0 && len(substr) > 0 && findSubstring(s, substr)))
}

func findSubstring(s, substr string) bool {
	for i := 0; i <= len(s)-len(substr); i++ {
		if s[i:i+len(substr)] == substr {
			return true
		}
	}
	return false
}
