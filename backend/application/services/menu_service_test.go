package services

import (
	"context"
	"errors"
	"testing"
	"time"

	"cafe-pos/backend/domain/ingredient"
	"cafe-pos/backend/domain/menu"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

// Mock MenuRepository for testing
type mockMenuRepositoryForMenuService struct {
	items map[primitive.ObjectID]*menu.MenuItem
}

func newMockMenuRepositoryForMenuService() *mockMenuRepositoryForMenuService {
	return &mockMenuRepositoryForMenuService{
		items: make(map[primitive.ObjectID]*menu.MenuItem),
	}
}

func (m *mockMenuRepositoryForMenuService) Create(ctx context.Context, item *menu.MenuItem) error {
	if item.ID.IsZero() {
		item.ID = primitive.NewObjectID()
	}
	item.CreatedAt = time.Now()
	item.UpdatedAt = time.Now()
	m.items[item.ID] = item
	return nil
}

func (m *mockMenuRepositoryForMenuService) FindAll(ctx context.Context) ([]*menu.MenuItem, error) {
	var items []*menu.MenuItem
	for _, item := range m.items {
		items = append(items, item)
	}
	return items, nil
}

func (m *mockMenuRepositoryForMenuService) FindByID(ctx context.Context, id primitive.ObjectID) (*menu.MenuItem, error) {
	item, exists := m.items[id]
	if !exists {
		return nil, errors.New("menu item not found")
	}
	return item, nil
}

func (m *mockMenuRepositoryForMenuService) FindByCategory(ctx context.Context, category string) ([]*menu.MenuItem, error) {
	var items []*menu.MenuItem
	for _, item := range m.items {
		if item.Category == category {
			items = append(items, item)
		}
	}
	return items, nil
}

func (m *mockMenuRepositoryForMenuService) FindByIngredientName(ctx context.Context, ingredientName string) ([]*menu.MenuItem, error) {
	var items []*menu.MenuItem
	for _, item := range m.items {
		for _, ing := range item.Ingredients {
			if ing.Name == ingredientName {
				items = append(items, item)
				break
			}
		}
	}
	return items, nil
}

func (m *mockMenuRepositoryForMenuService) Update(ctx context.Context, id primitive.ObjectID, item *menu.MenuItem) error {
	if _, exists := m.items[id]; !exists {
		return errors.New("menu item not found")
	}
	item.UpdatedAt = time.Now()
	m.items[id] = item
	return nil
}

func (m *mockMenuRepositoryForMenuService) Delete(ctx context.Context, id primitive.ObjectID) error {
	if _, exists := m.items[id]; !exists {
		return errors.New("menu item not found")
	}
	delete(m.items, id)
	return nil
}

// Test backward compatibility - single-size items

func TestCreateMenuItem_SingleSize(t *testing.T) {
	repo := newMockMenuRepositoryForMenuService()
	service := NewMenuService(repo)
	ctx := context.Background()

	req := &menu.CreateMenuItemRequest{
		Name:        "Bánh mì",
		Category:    "Món ăn",
		Description: "Bánh mì Việt Nam",
		HasVariants: false,
		Price:       20000,
		Ingredients: []menu.Ingredient{
			{Name: "Bánh mì", Quantity: 1, Unit: ingredient.UnitPiece},
		},
	}

	item, err := service.CreateMenuItem(ctx, req)
	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}

	if item.Name != "Bánh mì" {
		t.Errorf("Expected name 'Bánh mì', got '%s'", item.Name)
	}
	if item.HasVariants != false {
		t.Errorf("Expected has_variants false, got %v", item.HasVariants)
	}
	if item.Price != 20000 {
		t.Errorf("Expected price 20000, got %f", item.Price)
	}
	if len(item.Ingredients) != 1 {
		t.Errorf("Expected 1 ingredient, got %d", len(item.Ingredients))
	}
	if len(item.Variants) != 0 {
		t.Errorf("Expected 0 variants, got %d", len(item.Variants))
	}
}

func TestUpdateMenuItem_SingleSize(t *testing.T) {
	repo := newMockMenuRepositoryForMenuService()
	service := NewMenuService(repo)
	ctx := context.Background()

	// Create initial item
	createReq := &menu.CreateMenuItemRequest{
		Name:        "Bánh mì",
		Category:    "Món ăn",
		HasVariants: false,
		Price:       20000,
	}
	item, _ := service.CreateMenuItem(ctx, createReq)

	// Update item
	updateReq := &menu.UpdateMenuItemRequest{
		Name:  "Bánh mì updated",
		Price: 22000,
	}
	updated, err := service.UpdateMenuItem(ctx, item.ID, updateReq)
	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}

	if updated.Name != "Bánh mì updated" {
		t.Errorf("Expected name 'Bánh mì updated', got '%s'", updated.Name)
	}
	if updated.Price != 22000 {
		t.Errorf("Expected price 22000, got %f", updated.Price)
	}
}

// Test new functionality - multi-size items

func TestCreateMenuItem_MultiSize(t *testing.T) {
	repo := newMockMenuRepositoryForMenuService()
	service := NewMenuService(repo)
	ctx := context.Background()

	req := &menu.CreateMenuItemRequest{
		Name:        "Cà phê sữa đá",
		Category:    "Cà phê",
		Description: "Cà phê phin truyền thống",
		HasVariants: true,
		Variants: []menu.MenuItemVariant{
			{
				ID:        "M",
				Name:      "Size M",
				Price:     25000,
				Available: true,
				IsDefault: true,
				Ingredients: []menu.Ingredient{
					{Name: "Cà phê", Quantity: 20, Unit: ingredient.UnitGram},
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
				},
			},
		},
	}

	item, err := service.CreateMenuItem(ctx, req)
	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}

	if item.Name != "Cà phê sữa đá" {
		t.Errorf("Expected name 'Cà phê sữa đá', got '%s'", item.Name)
	}
	if item.HasVariants != true {
		t.Errorf("Expected has_variants true, got %v", item.HasVariants)
	}
	if len(item.Variants) != 2 {
		t.Errorf("Expected 2 variants, got %d", len(item.Variants))
	}
	if item.Price != 0 {
		t.Errorf("Expected price 0 (should be empty for multi-size), got %f", item.Price)
	}
	if len(item.Ingredients) != 0 {
		t.Errorf("Expected 0 ingredients (should be empty for multi-size), got %d", len(item.Ingredients))
	}
}

func TestUpdateMenuItem_MultiSize(t *testing.T) {
	repo := newMockMenuRepositoryForMenuService()
	service := NewMenuService(repo)
	ctx := context.Background()

	// Create initial multi-size item
	createReq := &menu.CreateMenuItemRequest{
		Name:        "Cà phê",
		Category:    "Cà phê",
		HasVariants: true,
		Variants: []menu.MenuItemVariant{
			{ID: "M", Name: "Size M", Price: 25000, IsDefault: true, Available: true},
		},
	}
	item, _ := service.CreateMenuItem(ctx, createReq)

	// Update variants
	hasVariants := true
	updateReq := &menu.UpdateMenuItemRequest{
		HasVariants: &hasVariants,
		Variants: []menu.MenuItemVariant{
			{ID: "M", Name: "Size M", Price: 26000, IsDefault: true, Available: true},
			{ID: "L", Name: "Size L", Price: 31000, IsDefault: false, Available: true},
		},
	}
	updated, err := service.UpdateMenuItem(ctx, item.ID, updateReq)
	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}

	if len(updated.Variants) != 2 {
		t.Errorf("Expected 2 variants, got %d", len(updated.Variants))
	}
	if updated.Variants[0].Price != 26000 {
		t.Errorf("Expected variant M price 26000, got %f", updated.Variants[0].Price)
	}
}

// Test toggling between modes

func TestUpdateMenuItem_SingleToMulti(t *testing.T) {
	repo := newMockMenuRepositoryForMenuService()
	service := NewMenuService(repo)
	ctx := context.Background()

	// Create single-size item
	createReq := &menu.CreateMenuItemRequest{
		Name:        "Bánh mì",
		Category:    "Món ăn",
		HasVariants: false,
		Price:       20000,
		Ingredients: []menu.Ingredient{
			{Name: "Bánh mì", Quantity: 1, Unit: ingredient.UnitPiece},
		},
	}
	item, _ := service.CreateMenuItem(ctx, createReq)

	// Toggle to multi-size
	hasVariants := true
	updateReq := &menu.UpdateMenuItemRequest{
		HasVariants: &hasVariants,
		Variants: []menu.MenuItemVariant{
			{ID: "M", Name: "Size M", Price: 20000, IsDefault: true, Available: true},
			{ID: "L", Name: "Size L", Price: 25000, IsDefault: false, Available: true},
		},
	}
	updated, err := service.UpdateMenuItem(ctx, item.ID, updateReq)
	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}

	if updated.HasVariants != true {
		t.Errorf("Expected has_variants true, got %v", updated.HasVariants)
	}
	if len(updated.Variants) != 2 {
		t.Errorf("Expected 2 variants, got %d", len(updated.Variants))
	}
	if updated.Price != 0 {
		t.Errorf("Expected price 0 (cleared), got %f", updated.Price)
	}
	if len(updated.Ingredients) != 0 {
		t.Errorf("Expected 0 ingredients (cleared), got %d", len(updated.Ingredients))
	}
}

func TestUpdateMenuItem_MultiToSingle(t *testing.T) {
	repo := newMockMenuRepositoryForMenuService()
	service := NewMenuService(repo)
	ctx := context.Background()

	// Create multi-size item
	createReq := &menu.CreateMenuItemRequest{
		Name:        "Cà phê",
		Category:    "Cà phê",
		HasVariants: true,
		Variants: []menu.MenuItemVariant{
			{ID: "M", Name: "Size M", Price: 25000, IsDefault: true, Available: true},
			{ID: "L", Name: "Size L", Price: 30000, IsDefault: false, Available: true},
		},
	}
	item, _ := service.CreateMenuItem(ctx, createReq)

	// Toggle to single-size
	hasVariants := false
	updateReq := &menu.UpdateMenuItemRequest{
		HasVariants: &hasVariants,
		Price:       28000,
		Ingredients: []menu.Ingredient{
			{Name: "Cà phê", Quantity: 25, Unit: ingredient.UnitGram},
		},
	}
	updated, err := service.UpdateMenuItem(ctx, item.ID, updateReq)
	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}

	if updated.HasVariants != false {
		t.Errorf("Expected has_variants false, got %v", updated.HasVariants)
	}
	if updated.Price != 28000 {
		t.Errorf("Expected price 28000, got %f", updated.Price)
	}
	if len(updated.Ingredients) != 1 {
		t.Errorf("Expected 1 ingredient, got %d", len(updated.Ingredients))
	}
	if len(updated.Variants) != 0 {
		t.Errorf("Expected 0 variants (cleared), got %d", len(updated.Variants))
	}
}

// Test error cases

func TestCreateMenuItem_MultiSize_NoVariants(t *testing.T) {
	repo := newMockMenuRepositoryForMenuService()
	service := NewMenuService(repo)
	ctx := context.Background()

	req := &menu.CreateMenuItemRequest{
		Name:        "Cà phê",
		Category:    "Cà phê",
		HasVariants: true,
		Variants:    []menu.MenuItemVariant{}, // Empty variants
	}

	_, err := service.CreateMenuItem(ctx, req)
	if err == nil {
		t.Error("Expected error for multi-size with no variants")
	}
}

func TestCreateMenuItem_MultiSize_NoDefault(t *testing.T) {
	repo := newMockMenuRepositoryForMenuService()
	service := NewMenuService(repo)
	ctx := context.Background()

	req := &menu.CreateMenuItemRequest{
		Name:        "Cà phê",
		Category:    "Cà phê",
		HasVariants: true,
		Variants: []menu.MenuItemVariant{
			{ID: "M", Name: "Size M", Price: 25000, IsDefault: false, Available: true},
			{ID: "L", Name: "Size L", Price: 30000, IsDefault: false, Available: true},
		},
	}

	_, err := service.CreateMenuItem(ctx, req)
	if err == nil {
		t.Error("Expected error for multi-size with no default variant")
	}
}

func TestCreateMenuItem_MultiSize_MultipleDefaults(t *testing.T) {
	repo := newMockMenuRepositoryForMenuService()
	service := NewMenuService(repo)
	ctx := context.Background()

	req := &menu.CreateMenuItemRequest{
		Name:        "Cà phê",
		Category:    "Cà phê",
		HasVariants: true,
		Variants: []menu.MenuItemVariant{
			{ID: "M", Name: "Size M", Price: 25000, IsDefault: true, Available: true},
			{ID: "L", Name: "Size L", Price: 30000, IsDefault: true, Available: true},
		},
	}

	_, err := service.CreateMenuItem(ctx, req)
	if err == nil {
		t.Error("Expected error for multi-size with multiple default variants")
	}
}

func TestCreateMenuItem_AmbiguousState(t *testing.T) {
	repo := newMockMenuRepositoryForMenuService()
	service := NewMenuService(repo)
	ctx := context.Background()

	// Note: With the current design, the service layer handles ambiguous requests
	// by only using the appropriate fields based on has_variants flag.
	// This test verifies that even if client sends both price and variants,
	// the service correctly ignores the inappropriate field.
	
	req := &menu.CreateMenuItemRequest{
		Name:        "Cà phê",
		Category:    "Cà phê",
		HasVariants: true,
		Price:       20000, // This will be ignored since has_variants=true
		Variants: []menu.MenuItemVariant{
			{ID: "M", Name: "Size M", Price: 25000, IsDefault: true, Available: true},
		},
	}

	item, err := service.CreateMenuItem(ctx, req)
	if err != nil {
		t.Errorf("Should handle ambiguous request by ignoring price field: %v", err)
	}
	
	// Verify that price was not set (should be 0)
	if item.Price != 0 {
		t.Errorf("Expected price to be 0 for multi-size item, got %f", item.Price)
	}
	
	// Verify variants were set correctly
	if len(item.Variants) != 1 {
		t.Errorf("Expected 1 variant, got %d", len(item.Variants))
	}
}

func TestCreateMenuItem_SingleSize_ZeroPrice(t *testing.T) {
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
		t.Error("Expected error for single-size with zero price")
	}
}
