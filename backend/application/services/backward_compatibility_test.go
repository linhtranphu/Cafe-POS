package services

import (
	"context"
	"errors"
	"testing"
	"time"

	"cafe-pos/backend/domain/ingredient"
	"cafe-pos/backend/domain/menu"
	"cafe-pos/backend/domain/order"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

// Mock repositories for backward compatibility testing
type mockIngredientRepositoryForBackwardCompat struct {
	ingredients map[string]*ingredient.Ingredient
}

func newMockIngredientRepositoryForBackwardCompat() *mockIngredientRepositoryForBackwardCompat {
	return &mockIngredientRepositoryForBackwardCompat{
		ingredients: make(map[string]*ingredient.Ingredient),
	}
}

func (m *mockIngredientRepositoryForBackwardCompat) Create(ctx context.Context, item *ingredient.Ingredient) error {
	if item.ID.IsZero() {
		item.ID = primitive.NewObjectID()
	}
	m.ingredients[item.Name] = item
	return nil
}

func (m *mockIngredientRepositoryForBackwardCompat) FindAll(ctx context.Context) ([]*ingredient.Ingredient, error) {
	items := make([]*ingredient.Ingredient, 0, len(m.ingredients))
	for _, item := range m.ingredients {
		items = append(items, item)
	}
	return items, nil
}

func (m *mockIngredientRepositoryForBackwardCompat) FindByID(ctx context.Context, id primitive.ObjectID) (*ingredient.Ingredient, error) {
	for _, item := range m.ingredients {
		if item.ID == id {
			return item, nil
		}
	}
	return nil, errors.New("ingredient not found")
}

func (m *mockIngredientRepositoryForBackwardCompat) Update(ctx context.Context, id primitive.ObjectID, item *ingredient.Ingredient) error {
	m.ingredients[item.Name] = item
	return nil
}

func (m *mockIngredientRepositoryForBackwardCompat) Delete(ctx context.Context, id primitive.ObjectID) error {
	for name, item := range m.ingredients {
		if item.ID == id {
			delete(m.ingredients, name)
			return nil
		}
	}
	return errors.New("ingredient not found")
}

func (m *mockIngredientRepositoryForBackwardCompat) FindLowStock(ctx context.Context) ([]*ingredient.Ingredient, error) {
	return []*ingredient.Ingredient{}, nil
}

func (m *mockIngredientRepositoryForBackwardCompat) CreateCategory(ctx context.Context, cat *ingredient.IngredientCategory) error {
	return nil
}

func (m *mockIngredientRepositoryForBackwardCompat) GetCategories(ctx context.Context) ([]ingredient.IngredientCategory, error) {
	return []ingredient.IngredientCategory{}, nil
}

func (m *mockIngredientRepositoryForBackwardCompat) DeleteCategory(ctx context.Context, id primitive.ObjectID) error {
	return nil
}

type mockOrderItemRepositoryForBackwardCompat struct{}

func newMockOrderItemRepositoryForBackwardCompat() *mockOrderItemRepositoryForBackwardCompat {
	return &mockOrderItemRepositoryForBackwardCompat{}
}

func (m *mockOrderItemRepositoryForBackwardCompat) FindByShiftID(ctx context.Context, shiftID primitive.ObjectID) ([]order.OrderItem, error) {
	return []order.OrderItem{}, nil
}

func (m *mockOrderItemRepositoryForBackwardCompat) CreateMany(ctx context.Context, items []*order.OrderItemWithCost) error {
	return nil
}

func (m *mockOrderItemRepositoryForBackwardCompat) FindByOrderIDs(ctx context.Context, orderIDs []primitive.ObjectID) ([]*order.OrderItemWithCost, error) {
	return []*order.OrderItemWithCost{}, nil
}

func (m *mockOrderItemRepositoryForBackwardCompat) FindByDateRange(ctx context.Context, startDate, endDate time.Time) ([]*order.OrderItemWithCost, error) {
	return []*order.OrderItemWithCost{}, nil
}

// TestBackwardCompatibility_SingleSizeItem verifies that single-size items
// work exactly as before the variants feature was added.
// This test ensures no regressions in existing functionality.
//
// Requirements: FR-1.5, FR-5.1, FR-6.1
func TestBackwardCompatibility_SingleSizeItem(t *testing.T) {
	ctx := context.Background()

	// Initialize mock repositories
	menuRepo := newMockMenuRepositoryForMenuService()
	orderRepo := newMockOrderRepositoryForVariants()
	shiftRepo := newMockShiftRepositoryForVariants()
	ingredientRepo := newMockIngredientRepositoryForBackwardCompat()
	orderItemRepo := newMockOrderItemRepositoryForBackwardCompat()

	// Initialize services
	_ = NewMenuService(menuRepo) // menuService not used in this test
	costCalculator := NewCostCalculatorService(menuRepo, ingredientRepo, orderRepo, orderItemRepo)

	t.Run("Step 1: Create single-size item", func(t *testing.T) {
		// Create a traditional single-size menu item (no variants)
		menuItem := &menu.MenuItem{
			Name:        "Bánh mì thịt",
			Category:    "Món ăn",
			Description: "Bánh mì Việt Nam truyền thống",
			Available:   true,
			HasVariants: false, // Single-size item
			Price:       20000,
			Ingredients: []menu.Ingredient{
				{
					Name:     "Bánh mì",
					Quantity: 1,
					Unit:     "cái",
				},
				{
					Name:     "Thịt",
					Quantity: 50,
					Unit:     "g",
				},
			},
			CreatedAt: time.Now(),
			UpdatedAt: time.Now(),
		}

		err := menuRepo.Create(ctx, menuItem)
		require.NoError(t, err, "Should create single-size item successfully")
		require.NotEmpty(t, menuItem.ID, "Item should have an ID")

		t.Run("Step 2: Verify has_variants=false", func(t *testing.T) {
			// Retrieve the item and verify structure
			retrieved, err := menuRepo.FindByID(ctx, menuItem.ID)
			require.NoError(t, err)
			
			assert.False(t, retrieved.HasVariants, "has_variants should be false")
			assert.Empty(t, retrieved.Variants, "variants array should be empty")
		})

		t.Run("Step 3: Verify price field populated", func(t *testing.T) {
			retrieved, err := menuRepo.FindByID(ctx, menuItem.ID)
			require.NoError(t, err)
			
			assert.Equal(t, 20000.0, retrieved.Price, "price should be set")
			assert.NotEmpty(t, retrieved.Ingredients, "ingredients should be set")
			assert.Len(t, retrieved.Ingredients, 2, "should have 2 ingredients")
		})

		t.Run("Step 4: Order without variant_id", func(t *testing.T) {
			// Create a shift first (required for orders)
			shift := &order.Shift{
				Type:      order.ShiftMorning,
				RoleType:  order.RoleWaiter,
				UserID:    primitive.NewObjectID(),
				StartedAt: time.Now(),
				Status:    order.ShiftOpen,
			}
			err := shiftRepo.Create(ctx, shift)
			require.NoError(t, err)

			// Create an order with the single-size item (no variant_id)
			newOrder := &order.Order{
				CustomerName: "Khách 1",
				Status:       order.StatusCreated,
				ShiftID:      shift.ID,
				Items: []order.OrderItem{
					{
						MenuItemID: menuItem.ID,
						Quantity:   2,
						// No VariantID - this is the key test for backward compatibility
					},
				},
				CreatedAt: time.Now(),
				UpdatedAt: time.Now(),
			}

			// Process the order (this should populate price from menu item)
			for i := range newOrder.Items {
				item := &newOrder.Items[i]
				menuItem, err := menuRepo.FindByID(ctx, item.MenuItemID)
				require.NoError(t, err)

				// Backward compatible logic: use item.Price directly
				assert.False(t, menuItem.HasVariants, "should be single-size")
				item.Name = menuItem.Name
				item.Price = menuItem.Price // Use price field directly
				item.Subtotal = item.Price * float64(item.Quantity)
			}

			newOrder.CalculateTotal()

			err = orderRepo.Create(ctx, newOrder)
			require.NoError(t, err, "Should create order without variant_id")

			// Verify order details
			assert.Equal(t, "Bánh mì thịt", newOrder.Items[0].Name)
			assert.Equal(t, 20000.0, newOrder.Items[0].Price)
			assert.Equal(t, 2, newOrder.Items[0].Quantity)
			assert.Equal(t, 40000.0, newOrder.Items[0].Subtotal)
			assert.Equal(t, 40000.0, newOrder.Total)
			assert.Empty(t, newOrder.Items[0].VariantID, "variant_id should be empty")
			assert.Empty(t, newOrder.Items[0].VariantName, "variant_name should be empty")
		})

		t.Run("Step 5: Calculate cost", func(t *testing.T) {
			// Create ingredient data for cost calculation
			banhMiIngredient := &ingredient.Ingredient{
				Name:              "Bánh mì",
				Unit:              "cái",
				CostPerUnit:       5000,
				WastagePercentage: 5,
				CreatedAt:         time.Now(),
				UpdatedAt:         time.Now(),
			}
			err := ingredientRepo.Create(ctx, banhMiIngredient)
			require.NoError(t, err)

			thitIngredient := &ingredient.Ingredient{
				Name:              "Thịt",
				Unit:              "g",
				CostPerUnit:       100, // 100 VND per gram
				WastagePercentage: 10,
				CreatedAt:         time.Now(),
				UpdatedAt:         time.Now(),
			}
			err = ingredientRepo.Create(ctx, thitIngredient)
			require.NoError(t, err)

			// Calculate cost for single-size item
			_, err = costCalculator.CalculateMenuItemCost(ctx, menuItem.ID)
			require.NoError(t, err, "Should calculate cost successfully")

			// Verify cost calculation
			retrieved, err := menuRepo.FindByID(ctx, menuItem.ID)
			require.NoError(t, err)

			// Expected cost calculation:
			// Bánh mì: 1 × 5000 × 1.0 × 1.05 = 5250
			// Thịt: 50 × 100 × 1.0 × 1.10 = 5500
			// Total: 10750
			expectedCost := 10750.0
			assert.Equal(t, expectedCost, retrieved.CurrentCost, "cost should be calculated correctly")
			assert.Equal(t, menu.CostStatusFinal, retrieved.CostStatus, "cost status should be FINAL")
			assert.False(t, retrieved.CostLastCalculatedAt.IsZero(), "cost_last_calculated_at should be set")

			// Verify variants are NOT used for cost calculation
			assert.Empty(t, retrieved.Variants, "variants should still be empty")
		})

		t.Run("Step 6: Verify no regressions", func(t *testing.T) {
			// Test all basic operations still work
			retrieved, err := menuRepo.FindByID(ctx, menuItem.ID)
			require.NoError(t, err)

			// 1. Can update the item
			retrieved.Price = 22000
			retrieved.Description = "Updated description"
			err = menuRepo.Update(ctx, retrieved.ID, retrieved)
			require.NoError(t, err, "Should update single-size item")

			updated, err := menuRepo.FindByID(ctx, retrieved.ID)
			require.NoError(t, err)
			assert.Equal(t, 22000.0, updated.Price, "price should be updated")
			assert.Equal(t, "Updated description", updated.Description)

			// 2. Can list items (should include single-size items)
			allItems, err := menuRepo.FindAll(ctx)
			require.NoError(t, err)
			assert.NotEmpty(t, allItems, "should find items")
			
			foundItem := false
			for _, item := range allItems {
				if item.ID == menuItem.ID {
					foundItem = true
					assert.False(t, item.HasVariants, "should be single-size")
					assert.Equal(t, 22000.0, item.Price)
				}
			}
			assert.True(t, foundItem, "should find our single-size item in list")

			// 3. Can delete the item
			err = menuRepo.Delete(ctx, menuItem.ID)
			require.NoError(t, err, "Should delete single-size item")

			_, err = menuRepo.FindByID(ctx, menuItem.ID)
			assert.Error(t, err, "Should not find deleted item")
		})
	})
}

// TestBackwardCompatibility_GetPriceMethod verifies the GetPrice method
// works correctly for single-size items (backward compatible behavior)
func TestBackwardCompatibility_GetPriceMethod(t *testing.T) {
	singleSizeItem := &menu.MenuItem{
		Name:        "Cà phê đen",
		HasVariants: false,
		Price:       15000,
	}

	// GetPrice should return the price field directly
	price := singleSizeItem.GetPrice("")
	assert.Equal(t, 15000.0, price, "GetPrice should return price field for single-size items")

	// Even if variant_id is provided, should still return price field
	price = singleSizeItem.GetPrice("M")
	assert.Equal(t, 15000.0, price, "GetPrice should ignore variant_id for single-size items")
}

// TestBackwardCompatibility_GetIngredientsMethod verifies the GetIngredients method
// works correctly for single-size items
func TestBackwardCompatibility_GetIngredientsMethod(t *testing.T) {
	ingredients := []menu.Ingredient{
		{Name: "Cà phê", Quantity: 20, Unit: "g"},
	}

	singleSizeItem := &menu.MenuItem{
		Name:        "Cà phê đen",
		HasVariants: false,
		Ingredients: ingredients,
	}

	// GetIngredients should return the ingredients field directly
	result := singleSizeItem.GetIngredients("")
	assert.Equal(t, ingredients, result, "GetIngredients should return ingredients field for single-size items")

	// Even if variant_id is provided, should still return ingredients field
	result = singleSizeItem.GetIngredients("M")
	assert.Equal(t, ingredients, result, "GetIngredients should ignore variant_id for single-size items")
}

// TestBackwardCompatibility_ValidationStillWorks verifies validation
// rules for single-size items remain unchanged
func TestBackwardCompatibility_ValidationStillWorks(t *testing.T) {
	t.Run("Valid single-size item passes validation", func(t *testing.T) {
		item := &menu.MenuItem{
			Name:        "Test Item",
			Category:    "Test",
			HasVariants: false,
			Price:       10000,
			Ingredients: []menu.Ingredient{
				{Name: "Ingredient 1", Quantity: 1, Unit: "unit"},
			},
		}

		err := item.Validate()
		assert.NoError(t, err, "Valid single-size item should pass validation")
	})

	t.Run("Single-size item without price fails validation", func(t *testing.T) {
		item := &menu.MenuItem{
			Name:        "Test Item",
			Category:    "Test",
			HasVariants: false,
			Price:       0, // Invalid
		}

		err := item.Validate()
		assert.Error(t, err, "Single-size item without price should fail validation")
		assert.Contains(t, err.Error(), "price must be > 0")
	})

	t.Run("Single-size item with variants array fails validation", func(t *testing.T) {
		item := &menu.MenuItem{
			Name:        "Test Item",
			Category:    "Test",
			HasVariants: false,
			Price:       10000,
			Variants: []menu.MenuItemVariant{
				{ID: "M", Name: "Size M", Price: 10000},
			},
		}

		err := item.Validate()
		assert.Error(t, err, "Single-size item with variants should fail validation")
		assert.Contains(t, err.Error(), "variants should not be set when has_variants=false")
	})
}
