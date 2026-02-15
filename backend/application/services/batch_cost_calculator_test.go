package services

import (
	"context"
	"testing"
	"time"

	"cafe-pos/backend/domain/batch"
	"cafe-pos/backend/domain/ingredient"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

// Mock ingredient repository for testing
type mockIngredientRepo struct {
	ingredients map[primitive.ObjectID]*ingredient.Ingredient
}

func newMockIngredientRepo() *mockIngredientRepo {
	return &mockIngredientRepo{
		ingredients: make(map[primitive.ObjectID]*ingredient.Ingredient),
	}
}

func (m *mockIngredientRepo) Create(ctx context.Context, item *ingredient.Ingredient) error {
	m.ingredients[item.ID] = item
	return nil
}

func (m *mockIngredientRepo) FindAll(ctx context.Context) ([]*ingredient.Ingredient, error) {
	result := make([]*ingredient.Ingredient, 0, len(m.ingredients))
	for _, ing := range m.ingredients {
		result = append(result, ing)
	}
	return result, nil
}

func (m *mockIngredientRepo) FindByID(ctx context.Context, id primitive.ObjectID) (*ingredient.Ingredient, error) {
	ing, exists := m.ingredients[id]
	if !exists {
		return nil, context.DeadlineExceeded // Use a standard error
	}
	return ing, nil
}

func (m *mockIngredientRepo) Update(ctx context.Context, id primitive.ObjectID, item *ingredient.Ingredient) error {
	m.ingredients[id] = item
	return nil
}

func (m *mockIngredientRepo) Delete(ctx context.Context, id primitive.ObjectID) error {
	delete(m.ingredients, id)
	return nil
}

func (m *mockIngredientRepo) FindLowStock(ctx context.Context) ([]*ingredient.Ingredient, error) {
	return nil, nil
}

func (m *mockIngredientRepo) CreateCategory(ctx context.Context, cat *ingredient.IngredientCategory) error {
	return nil
}

func (m *mockIngredientRepo) GetCategories(ctx context.Context) ([]ingredient.IngredientCategory, error) {
	return nil, nil
}

func (m *mockIngredientRepo) DeleteCategory(ctx context.Context, id primitive.ObjectID) error {
	return nil
}

func TestBatchCostCalculator_CalculateBatchCost_SingleIngredient(t *testing.T) {
	ctx := context.Background()
	repo := newMockIngredientRepo()
	calculator := NewBatchCostCalculator(repo)

	// Create test ingredient
	coffeeID := primitive.NewObjectID()
	coffee := &ingredient.Ingredient{
		ID:          coffeeID,
		Name:        "Hạt Cà Phê",
		Unit:        "g",
		CostPerUnit: 0.50, // 0.50 per gram
	}
	repo.ingredients[coffeeID] = coffee

	// Create batch definition: 100g coffee -> 500ml concentrate
	// Wastage rate: 10%
	batchDef := &batch.BatchDefinition{
		ID:             primitive.NewObjectID(),
		Name:           "Cà Phê Concentrate",
		Unit:           "ml",
		ShelfLifeHours: 24,
		ConversionRates: []batch.ConversionRate{
			{
				SourceIngredientID:   coffeeID,
				SourceIngredientName: "Hạt Cà Phê",
				SourceQuantity:       100,
				SourceUnit:           "g",
				BatchQuantity:        500,
				WastageRate:          0.10,
			},
		},
	}

	// Calculate cost for 500ml batch
	breakdown, err := calculator.CalculateBatchCost(ctx, batchDef, 500)
	require.NoError(t, err)
	require.NotNil(t, breakdown)

	// Verify calculations
	// Quantity needed: (500 / 500) * 100 = 100g
	// With wastage: 100 * (1 + 0.10) = 110g
	// Cost: 110 * 0.50 = 55.00
	assert.Equal(t, 1, len(breakdown.IngredientCosts))
	assert.InDelta(t, 110.0, breakdown.IngredientCosts[0].Quantity, 0.01)
	assert.InDelta(t, 55.0, breakdown.IngredientCosts[0].TotalCost, 0.01)
	assert.InDelta(t, 55.0, breakdown.TotalCost, 0.01)
	assert.InDelta(t, 0.11, breakdown.CostPerUnit, 0.01) // 55 / 500
}

func TestBatchCostCalculator_CalculateBatchCost_MultipleIngredients(t *testing.T) {
	ctx := context.Background()
	repo := newMockIngredientRepo()
	calculator := NewBatchCostCalculator(repo)

	// Create test ingredients
	coffeeID := primitive.NewObjectID()
	coffee := &ingredient.Ingredient{
		ID:          coffeeID,
		Name:        "Hạt Cà Phê",
		Unit:        "g",
		CostPerUnit: 0.50,
	}
	repo.ingredients[coffeeID] = coffee

	sugarID := primitive.NewObjectID()
	sugar := &ingredient.Ingredient{
		ID:          sugarID,
		Name:        "Đường",
		Unit:        "g",
		CostPerUnit: 0.02,
	}
	repo.ingredients[sugarID] = sugar

	// Create batch definition with multiple ingredients
	batchDef := &batch.BatchDefinition{
		ID:             primitive.NewObjectID(),
		Name:           "Cà Phê Đường",
		Unit:           "ml",
		ShelfLifeHours: 24,
		ConversionRates: []batch.ConversionRate{
			{
				SourceIngredientID:   coffeeID,
				SourceIngredientName: "Hạt Cà Phê",
				SourceQuantity:       100,
				SourceUnit:           "g",
				BatchQuantity:        500,
				WastageRate:          0.10,
			},
			{
				SourceIngredientID:   sugarID,
				SourceIngredientName: "Đường",
				SourceQuantity:       50,
				SourceUnit:           "g",
				BatchQuantity:        500,
				WastageRate:          0.05,
			},
		},
	}

	// Calculate cost for 500ml batch
	breakdown, err := calculator.CalculateBatchCost(ctx, batchDef, 500)
	require.NoError(t, err)
	require.NotNil(t, breakdown)

	// Verify calculations
	assert.Equal(t, 2, len(breakdown.IngredientCosts))
	
	// Coffee: (500/500) * 100 * (1 + 0.10) = 110g * 0.50 = 55.00
	assert.Equal(t, "Hạt Cà Phê", breakdown.IngredientCosts[0].IngredientName)
	assert.InDelta(t, 110.0, breakdown.IngredientCosts[0].Quantity, 0.01)
	assert.InDelta(t, 55.0, breakdown.IngredientCosts[0].TotalCost, 0.01)
	
	// Sugar: (500/500) * 50 * (1 + 0.05) = 52.5g * 0.02 = 1.05
	assert.Equal(t, "Đường", breakdown.IngredientCosts[1].IngredientName)
	assert.InDelta(t, 52.5, breakdown.IngredientCosts[1].Quantity, 0.01)
	assert.InDelta(t, 1.05, breakdown.IngredientCosts[1].TotalCost, 0.01)
	
	// Total: 55.00 + 1.05 = 56.05
	assert.InDelta(t, 56.05, breakdown.TotalCost, 0.01)
	assert.InDelta(t, 0.11, breakdown.CostPerUnit, 0.01) // 56.05 / 500 = 0.1121, rounded to 0.11
}

func TestBatchCostCalculator_CalculateBatchCost_DifferentQuantity(t *testing.T) {
	ctx := context.Background()
	repo := newMockIngredientRepo()
	calculator := NewBatchCostCalculator(repo)

	// Create test ingredient
	coffeeID := primitive.NewObjectID()
	coffee := &ingredient.Ingredient{
		ID:          coffeeID,
		Name:        "Hạt Cà Phê",
		Unit:        "g",
		CostPerUnit: 0.50,
	}
	repo.ingredients[coffeeID] = coffee

	// Create batch definition: 100g coffee -> 500ml concentrate
	batchDef := &batch.BatchDefinition{
		ID:             primitive.NewObjectID(),
		Name:           "Cà Phê Concentrate",
		Unit:           "ml",
		ShelfLifeHours: 24,
		ConversionRates: []batch.ConversionRate{
			{
				SourceIngredientID:   coffeeID,
				SourceIngredientName: "Hạt Cà Phê",
				SourceQuantity:       100,
				SourceUnit:           "g",
				BatchQuantity:        500,
				WastageRate:          0.10,
			},
		},
	}

	// Calculate cost for 1000ml batch (double the base quantity)
	breakdown, err := calculator.CalculateBatchCost(ctx, batchDef, 1000)
	require.NoError(t, err)
	require.NotNil(t, breakdown)

	// Verify calculations
	// Quantity needed: (1000 / 500) * 100 = 200g
	// With wastage: 200 * (1 + 0.10) = 220g
	// Cost: 220 * 0.50 = 110.00
	assert.InDelta(t, 220.0, breakdown.IngredientCosts[0].Quantity, 0.01)
	assert.InDelta(t, 110.0, breakdown.IngredientCosts[0].TotalCost, 0.01)
	assert.InDelta(t, 110.0, breakdown.TotalCost, 0.01)
	assert.InDelta(t, 0.11, breakdown.CostPerUnit, 0.01) // 110 / 1000
}

func TestBatchCostCalculator_CalculateBatchCost_NoWastage(t *testing.T) {
	ctx := context.Background()
	repo := newMockIngredientRepo()
	calculator := NewBatchCostCalculator(repo)

	// Create test ingredient
	coffeeID := primitive.NewObjectID()
	coffee := &ingredient.Ingredient{
		ID:          coffeeID,
		Name:        "Hạt Cà Phê",
		Unit:        "g",
		CostPerUnit: 0.50,
	}
	repo.ingredients[coffeeID] = coffee

	// Create batch definition with no wastage
	batchDef := &batch.BatchDefinition{
		ID:             primitive.NewObjectID(),
		Name:           "Cà Phê Concentrate",
		Unit:           "ml",
		ShelfLifeHours: 24,
		ConversionRates: []batch.ConversionRate{
			{
				SourceIngredientID:   coffeeID,
				SourceIngredientName: "Hạt Cà Phê",
				SourceQuantity:       100,
				SourceUnit:           "g",
				BatchQuantity:        500,
				WastageRate:          0.0, // No wastage
			},
		},
	}

	// Calculate cost for 500ml batch
	breakdown, err := calculator.CalculateBatchCost(ctx, batchDef, 500)
	require.NoError(t, err)
	require.NotNil(t, breakdown)

	// Verify calculations
	// Quantity needed: (500 / 500) * 100 = 100g
	// With wastage: 100 * (1 + 0.0) = 100g
	// Cost: 100 * 0.50 = 50.00
	assert.Equal(t, 100.0, breakdown.IngredientCosts[0].Quantity)
	assert.Equal(t, 50.0, breakdown.IngredientCosts[0].TotalCost)
	assert.Equal(t, 50.0, breakdown.TotalCost)
	assert.Equal(t, 0.10, breakdown.CostPerUnit) // 50 / 500
}

func TestBatchCostCalculator_CalculateBatchCost_InvalidQuantity(t *testing.T) {
	ctx := context.Background()
	repo := newMockIngredientRepo()
	calculator := NewBatchCostCalculator(repo)

	batchDef := &batch.BatchDefinition{
		ID:              primitive.NewObjectID(),
		Name:            "Test Batch",
		Unit:            "ml",
		ConversionRates: []batch.ConversionRate{},
	}

	// Test with zero quantity
	_, err := calculator.CalculateBatchCost(ctx, batchDef, 0)
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "must be greater than 0")

	// Test with negative quantity
	_, err = calculator.CalculateBatchCost(ctx, batchDef, -100)
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "must be greater than 0")
}

func TestBatchCostCalculator_CalculateBatchCost_IngredientNotFound(t *testing.T) {
	ctx := context.Background()
	repo := newMockIngredientRepo()
	calculator := NewBatchCostCalculator(repo)

	// Create batch definition with non-existent ingredient
	nonExistentID := primitive.NewObjectID()
	batchDef := &batch.BatchDefinition{
		ID:             primitive.NewObjectID(),
		Name:           "Test Batch",
		Unit:           "ml",
		ShelfLifeHours: 24,
		ConversionRates: []batch.ConversionRate{
			{
				SourceIngredientID:   nonExistentID,
				SourceIngredientName: "Non-existent",
				SourceQuantity:       100,
				SourceUnit:           "g",
				BatchQuantity:        500,
				WastageRate:          0.0,
			},
		},
	}

	// Calculate cost should fail
	_, err := calculator.CalculateBatchCost(ctx, batchDef, 500)
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "failed to get cost")
}

func TestBatchCostCalculator_CalculateBatchCost_IngredientNoCost(t *testing.T) {
	ctx := context.Background()
	repo := newMockIngredientRepo()
	calculator := NewBatchCostCalculator(repo)

	// Create ingredient with no cost
	coffeeID := primitive.NewObjectID()
	coffee := &ingredient.Ingredient{
		ID:          coffeeID,
		Name:        "Hạt Cà Phê",
		Unit:        "g",
		CostPerUnit: 0, // No cost set
	}
	repo.ingredients[coffeeID] = coffee

	batchDef := &batch.BatchDefinition{
		ID:             primitive.NewObjectID(),
		Name:           "Test Batch",
		Unit:           "ml",
		ShelfLifeHours: 24,
		ConversionRates: []batch.ConversionRate{
			{
				SourceIngredientID:   coffeeID,
				SourceIngredientName: "Hạt Cà Phê",
				SourceQuantity:       100,
				SourceUnit:           "g",
				BatchQuantity:        500,
				WastageRate:          0.0,
			},
		},
	}

	// Calculate cost should fail
	_, err := calculator.CalculateBatchCost(ctx, batchDef, 500)
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "has no cost set")
}

func TestBatchCostCalculator_CostCaching(t *testing.T) {
	ctx := context.Background()
	repo := newMockIngredientRepo()
	calculator := NewBatchCostCalculator(repo)

	// Create test ingredient
	coffeeID := primitive.NewObjectID()
	coffee := &ingredient.Ingredient{
		ID:          coffeeID,
		Name:        "Hạt Cà Phê",
		Unit:        "g",
		CostPerUnit: 0.50,
	}
	repo.ingredients[coffeeID] = coffee

	batchDef := &batch.BatchDefinition{
		ID:             primitive.NewObjectID(),
		Name:           "Cà Phê Concentrate",
		Unit:           "ml",
		ShelfLifeHours: 24,
		ConversionRates: []batch.ConversionRate{
			{
				SourceIngredientID:   coffeeID,
				SourceIngredientName: "Hạt Cà Phê",
				SourceQuantity:       100,
				SourceUnit:           "g",
				BatchQuantity:        500,
				WastageRate:          0.10,
			},
		},
	}

	// First calculation - should hit database
	breakdown1, err := calculator.CalculateBatchCost(ctx, batchDef, 500)
	require.NoError(t, err)
	assert.Equal(t, 55.0, breakdown1.TotalCost)

	// Change ingredient cost in repo
	coffee.CostPerUnit = 1.00
	repo.ingredients[coffeeID] = coffee

	// Second calculation - should use cache (old cost)
	breakdown2, err := calculator.CalculateBatchCost(ctx, batchDef, 500)
	require.NoError(t, err)
	assert.Equal(t, 55.0, breakdown2.TotalCost) // Still old cost

	// Invalidate cache
	calculator.InvalidateCache()

	// Third calculation - should hit database (new cost)
	breakdown3, err := calculator.CalculateBatchCost(ctx, batchDef, 500)
	require.NoError(t, err)
	assert.Equal(t, 110.0, breakdown3.TotalCost) // New cost: 110 * 1.00
}

func TestBatchCostCalculator_CacheExpiry(t *testing.T) {
	ctx := context.Background()
	repo := newMockIngredientRepo()
	calculator := NewBatchCostCalculator(repo)
	
	// Set very short TTL for testing
	calculator.costCache.ttl = 100 * time.Millisecond

	// Create test ingredient
	coffeeID := primitive.NewObjectID()
	coffee := &ingredient.Ingredient{
		ID:          coffeeID,
		Name:        "Hạt Cà Phê",
		Unit:        "g",
		CostPerUnit: 0.50,
	}
	repo.ingredients[coffeeID] = coffee

	batchDef := &batch.BatchDefinition{
		ID:             primitive.NewObjectID(),
		Name:           "Cà Phê Concentrate",
		Unit:           "ml",
		ShelfLifeHours: 24,
		ConversionRates: []batch.ConversionRate{
			{
				SourceIngredientID:   coffeeID,
				SourceIngredientName: "Hạt Cà Phê",
				SourceQuantity:       100,
				SourceUnit:           "g",
				BatchQuantity:        500,
				WastageRate:          0.10,
			},
		},
	}

	// First calculation
	breakdown1, err := calculator.CalculateBatchCost(ctx, batchDef, 500)
	require.NoError(t, err)
	assert.Equal(t, 55.0, breakdown1.TotalCost)

	// Change ingredient cost
	coffee.CostPerUnit = 1.00
	repo.ingredients[coffeeID] = coffee

	// Wait for cache to expire
	time.Sleep(150 * time.Millisecond)

	// Second calculation - cache expired, should use new cost
	breakdown2, err := calculator.CalculateBatchCost(ctx, batchDef, 500)
	require.NoError(t, err)
	assert.Equal(t, 110.0, breakdown2.TotalCost) // New cost
}

func TestBatchCostCalculator_InvalidateCacheForIngredient(t *testing.T) {
	ctx := context.Background()
	repo := newMockIngredientRepo()
	calculator := NewBatchCostCalculator(repo)

	// Create two ingredients
	coffeeID := primitive.NewObjectID()
	coffee := &ingredient.Ingredient{
		ID:          coffeeID,
		Name:        "Hạt Cà Phê",
		Unit:        "g",
		CostPerUnit: 0.50,
	}
	repo.ingredients[coffeeID] = coffee

	sugarID := primitive.NewObjectID()
	sugar := &ingredient.Ingredient{
		ID:          sugarID,
		Name:        "Đường",
		Unit:        "g",
		CostPerUnit: 0.02,
	}
	repo.ingredients[sugarID] = sugar

	batchDef := &batch.BatchDefinition{
		ID:             primitive.NewObjectID(),
		Name:           "Cà Phê Đường",
		Unit:           "ml",
		ShelfLifeHours: 24,
		ConversionRates: []batch.ConversionRate{
			{
				SourceIngredientID:   coffeeID,
				SourceIngredientName: "Hạt Cà Phê",
				SourceQuantity:       100,
				SourceUnit:           "g",
				BatchQuantity:        500,
				WastageRate:          0.10,
			},
			{
				SourceIngredientID:   sugarID,
				SourceIngredientName: "Đường",
				SourceQuantity:       50,
				SourceUnit:           "g",
				BatchQuantity:        500,
				WastageRate:          0.05,
			},
		},
	}

	// First calculation - cache both ingredients
	breakdown1, err := calculator.CalculateBatchCost(ctx, batchDef, 500)
	require.NoError(t, err)
	assert.Equal(t, 56.05, breakdown1.TotalCost)

	// Change both ingredient costs
	coffee.CostPerUnit = 1.00
	sugar.CostPerUnit = 0.04
	repo.ingredients[coffeeID] = coffee
	repo.ingredients[sugarID] = sugar

	// Invalidate only coffee cache
	calculator.InvalidateCacheForIngredient(coffeeID)

	// Second calculation - coffee uses new cost, sugar uses cached cost
	breakdown2, err := calculator.CalculateBatchCost(ctx, batchDef, 500)
	require.NoError(t, err)
	// Coffee: 110 * 1.00 = 110.00
	// Sugar: 52.5 * 0.02 = 1.05 (cached)
	// Total: 111.05
	assert.Equal(t, 111.05, breakdown2.TotalCost)
}
