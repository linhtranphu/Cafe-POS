package services

import (
	"context"
	"strings"
	"testing"

	"cafe-pos/backend/domain"
	"cafe-pos/backend/domain/ingredient"
	"cafe-pos/backend/domain/menu"
	"cafe-pos/backend/domain/order"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

// Task 14.1: Test validation errors
// This file consolidates all validation error tests for menu size variants

// Test 1: Create multi-size without variants (should fail)
func TestValidationError_CreateMultiSize_NoVariants(t *testing.T) {
	repo := newMockMenuRepositoryForMenuService()
	service := NewMenuService(repo)
	ctx := context.Background()

	req := &menu.CreateMenuItemRequest{
		Name:        "Cà phê",
		Category:    "Cà phê",
		Description: "Cà phê phin",
		HasVariants: true,
		Variants:    []menu.MenuItemVariant{}, // Empty - should fail
	}

	_, err := service.CreateMenuItem(ctx, req)
	if err == nil {
		t.Error("Expected error for multi-size item without variants, got nil")
	}

	// Verify error message is clear
	if !strings.Contains(err.Error(), "variants required") {
		t.Errorf("Expected error message to mention 'variants required', got: %v", err)
	}
}

// Test 2: Create multi-size without default variant (should fail)
func TestValidationError_CreateMultiSize_NoDefaultVariant(t *testing.T) {
	repo := newMockMenuRepositoryForMenuService()
	service := NewMenuService(repo)
	ctx := context.Background()

	req := &menu.CreateMenuItemRequest{
		Name:        "Cà phê",
		Category:    "Cà phê",
		Description: "Cà phê phin",
		HasVariants: true,
		Variants: []menu.MenuItemVariant{
			{
				ID:        "M",
				Name:      "Size M",
				Price:     25000,
				Available: true,
				IsDefault: false, // No default
			},
			{
				ID:        "L",
				Name:      "Size L",
				Price:     30000,
				Available: true,
				IsDefault: false, // No default
			},
		},
	}

	_, err := service.CreateMenuItem(ctx, req)
	if err == nil {
		t.Error("Expected error for multi-size item without default variant, got nil")
	}

	// Verify error message is clear
	if !strings.Contains(err.Error(), "default variant") {
		t.Errorf("Expected error message to mention 'default variant', got: %v", err)
	}
}

// Test 3: Create multi-size with duplicate variant IDs (should fail)
func TestValidationError_CreateMultiSize_DuplicateVariantIDs(t *testing.T) {
	repo := newMockMenuRepositoryForMenuService()
	service := NewMenuService(repo)
	ctx := context.Background()

	req := &menu.CreateMenuItemRequest{
		Name:        "Cà phê",
		Category:    "Cà phê",
		Description: "Cà phê phin",
		HasVariants: true,
		Variants: []menu.MenuItemVariant{
			{
				ID:        "M",
				Name:      "Size M",
				Price:     25000,
				Available: true,
				IsDefault: true,
			},
			{
				ID:        "M", // Duplicate ID
				Name:      "Size M Duplicate",
				Price:     30000,
				Available: true,
				IsDefault: false,
			},
		},
	}

	_, err := service.CreateMenuItem(ctx, req)
	if err == nil {
		t.Error("Expected error for duplicate variant IDs, got nil")
	}

	// Verify error message is clear
	if !strings.Contains(err.Error(), "duplicate") {
		t.Errorf("Expected error message to mention 'duplicate', got: %v", err)
	}
}

// Test 4: Create item with both price and variants (should fail)
func TestValidationError_CreateItem_BothPriceAndVariants(t *testing.T) {
	repo := newMockMenuRepositoryForMenuService()
	service := NewMenuService(repo)
	ctx := context.Background()

	req := &menu.CreateMenuItemRequest{
		Name:        "Cà phê",
		Category:    "Cà phê",
		Description: "Cà phê phin",
		HasVariants: true,
		Price:       20000, // Should NOT be set when has_variants=true
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

	_, err := service.CreateMenuItem(ctx, req)
	if err == nil {
		t.Error("Expected error for item with both price and variants, got nil")
	}

	// Verify error message is clear
	if !strings.Contains(err.Error(), "price should not be set") {
		t.Errorf("Expected error message to mention 'price should not be set', got: %v", err)
	}
}

// Test 5: Order multi-size without variant_id (should fail with clear error)
func TestValidationError_OrderMultiSize_MissingVariantID(t *testing.T) {
	orderRepo := newMockOrderRepositoryForVariants()
	shiftRepo := newMockShiftRepositoryForVariants()
	menuRepo := newMockMenuRepositoryForVariants()
	smManager := domain.NewStateMachineManager()

	service := NewOrderService(orderRepo, shiftRepo, menuRepo, smManager, nil)
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

	// Create order WITHOUT variant_id
	req := &order.CreateOrderRequest{
		CustomerName: "Khách 1",
		ShiftID:      shift.ID.Hex(),
		Items: []order.OrderItem{
			{
				MenuItemID: menuItem.ID,
				// Missing VariantID - should fail
				Quantity: 1,
			},
		},
	}

	waiterID := primitive.NewObjectID().Hex()
	_, err := service.CreateOrder(ctx, req, waiterID, "Waiter 1")
	if err == nil {
		t.Error("Expected error for ordering multi-size item without variant_id, got nil")
	}

	// Verify error message is clear and mentions variant_id
	errMsg := err.Error()
	if !strings.Contains(errMsg, "variant_id") && !strings.Contains(errMsg, "variant") {
		t.Errorf("Expected error message to mention 'variant_id' or 'variant', got: %v", err)
	}
}

// Test 6: Order with invalid variant_id (should fail with clear error)
func TestValidationError_OrderMultiSize_InvalidVariantID(t *testing.T) {
	orderRepo := newMockOrderRepositoryForVariants()
	shiftRepo := newMockShiftRepositoryForVariants()
	menuRepo := newMockMenuRepositoryForVariants()
	smManager := domain.NewStateMachineManager()

	service := NewOrderService(orderRepo, shiftRepo, menuRepo, smManager, nil)
	ctx := context.Background()

	// Create open shift
	shift := &order.Shift{
		ID:     primitive.NewObjectID(),
		Status: order.ShiftOpen,
	}
	shiftRepo.Create(ctx, shift)

	// Create multi-size menu item with only M and L sizes
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

	// Create order with INVALID variant_id (XL doesn't exist)
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
		t.Error("Expected error for ordering with invalid variant_id, got nil")
	}

	// Verify error message is clear and mentions the invalid variant_id
	errMsg := err.Error()
	if !strings.Contains(errMsg, "invalid") && !strings.Contains(errMsg, "variant") {
		t.Errorf("Expected error message to mention 'invalid variant', got: %v", err)
	}
}

// Additional validation tests for completeness

// Test: Create single-size with variants set (should fail)
func TestValidationError_CreateSingleSize_WithVariants(t *testing.T) {
	repo := newMockMenuRepositoryForMenuService()
	service := NewMenuService(repo)
	ctx := context.Background()

	req := &menu.CreateMenuItemRequest{
		Name:        "Bánh mì",
		Category:    "Món ăn",
		HasVariants: false,
		Price:       20000,
		Variants: []menu.MenuItemVariant{ // Should NOT be set when has_variants=false
			{
				ID:        "M",
				Name:      "Size M",
				Price:     25000,
				Available: true,
				IsDefault: true,
			},
		},
	}

	_, err := service.CreateMenuItem(ctx, req)
	if err == nil {
		t.Error("Expected error for single-size item with variants set, got nil")
		return
	}

	// Verify error message is clear
	if !strings.Contains(err.Error(), "variants should not be set") {
		t.Errorf("Expected error message to mention 'variants should not be set', got: %v", err)
	}
}

// Test: Create multi-size with ingredients set at item level (should fail)
func TestValidationError_CreateMultiSize_WithIngredientsAtItemLevel(t *testing.T) {
	repo := newMockMenuRepositoryForMenuService()
	service := NewMenuService(repo)
	ctx := context.Background()

	req := &menu.CreateMenuItemRequest{
		Name:        "Cà phê",
		Category:    "Cà phê",
		HasVariants: true,
		Ingredients: []menu.Ingredient{ // Should NOT be set when has_variants=true
			{Name: "Cà phê", Quantity: 20, Unit: ingredient.UnitGram},
		},
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

	_, err := service.CreateMenuItem(ctx, req)
	if err == nil {
		t.Error("Expected error for multi-size item with ingredients at item level, got nil")
	}

	// Verify error message is clear
	if !strings.Contains(err.Error(), "ingredients should not be set") {
		t.Errorf("Expected error message to mention 'ingredients should not be set', got: %v", err)
	}
}

// Test: Create single-size with zero price (should fail)
func TestValidationError_CreateSingleSize_ZeroPrice(t *testing.T) {
	repo := newMockMenuRepositoryForMenuService()
	service := NewMenuService(repo)
	ctx := context.Background()

	req := &menu.CreateMenuItemRequest{
		Name:        "Bánh mì",
		Category:    "Món ăn",
		HasVariants: false,
		Price:       0, // Invalid
	}

	_, err := service.CreateMenuItem(ctx, req)
	if err == nil {
		t.Error("Expected error for single-size item with zero price, got nil")
	}

	// Verify error message is clear
	if !strings.Contains(err.Error(), "price must be > 0") {
		t.Errorf("Expected error message to mention 'price must be > 0', got: %v", err)
	}
}

// Test: Create multi-size with zero variant price (should fail)
func TestValidationError_CreateMultiSize_ZeroVariantPrice(t *testing.T) {
	repo := newMockMenuRepositoryForMenuService()
	service := NewMenuService(repo)
	ctx := context.Background()

	req := &menu.CreateMenuItemRequest{
		Name:        "Cà phê",
		Category:    "Cà phê",
		HasVariants: true,
		Variants: []menu.MenuItemVariant{
			{
				ID:        "M",
				Name:      "Size M",
				Price:     0, // Invalid
				Available: true,
				IsDefault: true,
			},
		},
	}

	_, err := service.CreateMenuItem(ctx, req)
	if err == nil {
		t.Error("Expected error for variant with zero price, got nil")
	}

	// Verify error message is clear
	if !strings.Contains(err.Error(), "price must be > 0") {
		t.Errorf("Expected error message to mention 'price must be > 0', got: %v", err)
	}
}

// Test: Create multi-size with empty variant ID (should fail)
func TestValidationError_CreateMultiSize_EmptyVariantID(t *testing.T) {
	repo := newMockMenuRepositoryForMenuService()
	service := NewMenuService(repo)
	ctx := context.Background()

	req := &menu.CreateMenuItemRequest{
		Name:        "Cà phê",
		Category:    "Cà phê",
		HasVariants: true,
		Variants: []menu.MenuItemVariant{
			{
				ID:        "", // Empty - invalid
				Name:      "Size M",
				Price:     25000,
				Available: true,
				IsDefault: true,
			},
		},
	}

	_, err := service.CreateMenuItem(ctx, req)
	if err == nil {
		t.Error("Expected error for variant with empty ID, got nil")
	}

	// Verify error message is clear
	if !strings.Contains(err.Error(), "variant ID is required") {
		t.Errorf("Expected error message to mention 'variant ID is required', got: %v", err)
	}
}

// Test: Create multi-size with multiple default variants (should fail)
func TestValidationError_CreateMultiSize_MultipleDefaults(t *testing.T) {
	repo := newMockMenuRepositoryForMenuService()
	service := NewMenuService(repo)
	ctx := context.Background()

	req := &menu.CreateMenuItemRequest{
		Name:        "Cà phê",
		Category:    "Cà phê",
		HasVariants: true,
		Variants: []menu.MenuItemVariant{
			{
				ID:        "M",
				Name:      "Size M",
				Price:     25000,
				Available: true,
				IsDefault: true, // Default
			},
			{
				ID:        "L",
				Name:      "Size L",
				Price:     30000,
				Available: true,
				IsDefault: true, // Also default - invalid
			},
		},
	}

	_, err := service.CreateMenuItem(ctx, req)
	if err == nil {
		t.Error("Expected error for multiple default variants, got nil")
	}

	// Verify error message is clear
	if !strings.Contains(err.Error(), "exactly one default") {
		t.Errorf("Expected error message to mention 'exactly one default', got: %v", err)
	}
}

// Test: Order multi-size with unavailable variant (should fail)
func TestValidationError_OrderMultiSize_UnavailableVariant(t *testing.T) {
	orderRepo := newMockOrderRepositoryForVariants()
	shiftRepo := newMockShiftRepositoryForVariants()
	menuRepo := newMockMenuRepositoryForVariants()
	smManager := domain.NewStateMachineManager()

	service := NewOrderService(orderRepo, shiftRepo, menuRepo, smManager, nil)
	ctx := context.Background()

	// Create open shift
	shift := &order.Shift{
		ID:     primitive.NewObjectID(),
		Status: order.ShiftOpen,
	}
	shiftRepo.Create(ctx, shift)

	// Create multi-size menu item with unavailable variant
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
				Available: false, // Unavailable
				IsDefault: false,
			},
		},
	}
	menuRepo.Create(ctx, menuItem)

	// Try to order unavailable variant
	req := &order.CreateOrderRequest{
		CustomerName: "Khách 1",
		ShiftID:      shift.ID.Hex(),
		Items: []order.OrderItem{
			{
				MenuItemID: menuItem.ID,
				VariantID:  "L", // Unavailable
				Quantity:   1,
			},
		},
	}

	waiterID := primitive.NewObjectID().Hex()
	_, err := service.CreateOrder(ctx, req, waiterID, "Waiter 1")
	if err == nil {
		t.Error("Expected error for ordering unavailable variant, got nil")
	}

	// Verify error message is clear
	errMsg := err.Error()
	if !strings.Contains(errMsg, "not available") {
		t.Errorf("Expected error message to mention 'not available', got: %v", err)
	}
}
