package services

import (
	"context"
	"testing"
	"time"

	"cafe-pos/backend/domain/batch"
	"cafe-pos/backend/domain/ingredient"
	"cafe-pos/backend/infrastructure/mongodb"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

// TestBatchInventoryIntegration tests the integration between batch creation and inventory system
// Note: Full transaction tests require MongoDB replica set. These tests focus on validation logic.
func TestBatchInventoryIntegration(t *testing.T) {
	// Setup MongoDB test database
	ctx := context.Background()
	
	// Get MongoDB URI from environment or use default with auth
	mongoURI := "mongodb://admin:password123@localhost:27017/?authSource=admin"
	client, err := mongo.Connect(ctx, options.Client().ApplyURI(mongoURI))
	require.NoError(t, err)
	defer client.Disconnect(ctx)

	// Use a test database
	testDB := client.Database("cafe_pos_test_batch_inventory_" + primitive.NewObjectID().Hex())
	defer testDB.Drop(ctx)

	// Setup repositories
	ingredientRepo := mongodb.NewIngredientRepository(testDB)
	stockHistoryRepo := mongodb.NewStockHistoryRepository(testDB)
	batchDefinitionRepo := mongodb.NewBatchDefinitionRepository(testDB)
	batchRecordRepo := mongodb.NewBatchRecordRepository(testDB)

	// Setup services
	batchCostCalculator := NewBatchCostCalculator(ingredientRepo)
	batchRecordService := NewBatchRecordService(
		batchRecordRepo,
		batchDefinitionRepo,
		ingredientRepo,
		stockHistoryRepo,
		nil,
		batchCostCalculator,
		client,
	)

	t.Run("CreateBatch_ValidatesInsufficientIngredients", func(t *testing.T) {
		// Create test ingredient with insufficient quantity
		testIngredient := &ingredient.Ingredient{
			Name:        "Milk",
			Category:    "Dairy",
			Unit:        ingredient.UnitMilliliter,
			Quantity:    50.0, // Not enough for batch
			MinStock:    100.0,
			CostPerUnit: 0.02,
			Supplier:    "Test Supplier",
			CreatedAt:   time.Now(),
			UpdatedAt:   time.Now(),
		}
		err := ingredientRepo.Create(ctx, testIngredient)
		require.NoError(t, err)

		// Create batch definition requiring more than available
		batchDef := &batch.BatchDefinition{
			Name:           "Milk Batch",
			Unit:           "ml",
			ShelfLifeHours: 12,
			ConversionRates: []batch.ConversionRate{
				{
					SourceIngredientID:   testIngredient.ID,
					SourceIngredientName: testIngredient.Name,
					SourceQuantity:       100.0, // Need 100ml but only have 50ml
					SourceUnit:           string(testIngredient.Unit),
					BatchQuantity:        100.0,
					WastageRate:          0.0,
				},
			},
			LowStockThreshold:  50.0,
			ExpiryWarningHours: 2,
			CreatedAt:          time.Now(),
			UpdatedAt:          time.Now(),
		}
		err = batchDefinitionRepo.Create(ctx, batchDef)
		require.NoError(t, err)

		// Attempt to create batch (should fail before transaction)
		req := CreateBatchRequest{
			BatchDefinitionID: batchDef.ID,
			QuantityProduced:  100.0,
			PreparedBy:        "test-user",
		}

		_, err = batchRecordService.CreateBatch(ctx, req)
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "insufficient")

		// Verify ingredient quantity was NOT changed
		unchangedIngredient, err := ingredientRepo.FindByID(ctx, testIngredient.ID)
		require.NoError(t, err)
		assert.Equal(t, 50.0, unchangedIngredient.Quantity)

		// Verify no stock history was created
		histories, err := stockHistoryRepo.FindByIngredientID(ctx, testIngredient.ID)
		require.NoError(t, err)
		assert.Len(t, histories, 0)

		// Verify no batch record was created
		filter := batch.BatchRecordFilter{
			BatchDefinitionID: &batchDef.ID,
		}
		records, _, err := batchRecordRepo.FindAll(ctx, filter)
		require.NoError(t, err)
		assert.Len(t, records, 0)
	})

	t.Run("CreateBatch_ValidatesMultipleIngredientsBeforeTransaction", func(t *testing.T) {
		// Create two ingredients, one with insufficient quantity
		ingredient1 := &ingredient.Ingredient{
			Name:        "Chocolate",
			Category:    "Flavors",
			Unit:        ingredient.UnitGram,
			Quantity:    1000.0, // Sufficient
			MinStock:    100.0,
			CostPerUnit: 0.15,
			Supplier:    "Test Supplier",
			CreatedAt:   time.Now(),
			UpdatedAt:   time.Now(),
		}
		err := ingredientRepo.Create(ctx, ingredient1)
		require.NoError(t, err)

		ingredient2 := &ingredient.Ingredient{
			Name:        "Cream",
			Category:    "Dairy",
			Unit:        ingredient.UnitMilliliter,
			Quantity:    50.0, // Insufficient
			MinStock:    100.0,
			CostPerUnit: 0.03,
			Supplier:    "Test Supplier",
			CreatedAt:   time.Now(),
			UpdatedAt:   time.Now(),
		}
		err = ingredientRepo.Create(ctx, ingredient2)
		require.NoError(t, err)

		// Create batch definition
		batchDef := &batch.BatchDefinition{
			Name:           "Chocolate Cream",
			Unit:           "ml",
			ShelfLifeHours: 24,
			ConversionRates: []batch.ConversionRate{
				{
					SourceIngredientID:   ingredient1.ID,
					SourceIngredientName: ingredient1.Name,
					SourceQuantity:       100.0,
					SourceUnit:           string(ingredient1.Unit),
					BatchQuantity:        300.0,
					WastageRate:          0.0,
				},
				{
					SourceIngredientID:   ingredient2.ID,
					SourceIngredientName: ingredient2.Name,
					SourceQuantity:       200.0, // Need 200ml but only have 50ml
					SourceUnit:           string(ingredient2.Unit),
					BatchQuantity:        300.0,
					WastageRate:          0.0,
				},
			},
			LowStockThreshold:  200.0,
			ExpiryWarningHours: 4,
			CreatedAt:          time.Now(),
			UpdatedAt:          time.Now(),
		}
		err = batchDefinitionRepo.Create(ctx, batchDef)
		require.NoError(t, err)

		// Attempt to create batch (should fail before transaction)
		req := CreateBatchRequest{
			BatchDefinitionID: batchDef.ID,
			QuantityProduced:  300.0,
			PreparedBy:        "test-user",
		}

		_, err = batchRecordService.CreateBatch(ctx, req)
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "insufficient")

		// Verify NEITHER ingredient was deducted (validation prevented transaction)
		unchangedIngredient1, err := ingredientRepo.FindByID(ctx, ingredient1.ID)
		require.NoError(t, err)
		assert.Equal(t, 1000.0, unchangedIngredient1.Quantity)

		unchangedIngredient2, err := ingredientRepo.FindByID(ctx, ingredient2.ID)
		require.NoError(t, err)
		assert.Equal(t, 50.0, unchangedIngredient2.Quantity)

		// Verify no stock history was created for either ingredient
		histories1, err := stockHistoryRepo.FindByIngredientID(ctx, ingredient1.ID)
		require.NoError(t, err)
		assert.Len(t, histories1, 0)

		histories2, err := stockHistoryRepo.FindByIngredientID(ctx, ingredient2.ID)
		require.NoError(t, err)
		assert.Len(t, histories2, 0)
	})

	t.Run("CreateBatch_CalculatesCostWithWastage", func(t *testing.T) {
		// Create test ingredient
		testIngredient := &ingredient.Ingredient{
			Name:        "Test Coffee",
			Category:    "Raw Materials",
			Unit:        ingredient.UnitGram,
			Quantity:    1000.0,
			MinStock:    100.0,
			CostPerUnit: 0.05,
			Supplier:    "Test Supplier",
			CreatedAt:   time.Now(),
			UpdatedAt:   time.Now(),
		}
		err := ingredientRepo.Create(ctx, testIngredient)
		require.NoError(t, err)

		// Create batch definition with wastage
		batchDef := &batch.BatchDefinition{
			Name:           "Test Concentrate",
			Unit:           "ml",
			ShelfLifeHours: 24,
			ConversionRates: []batch.ConversionRate{
				{
					SourceIngredientID:   testIngredient.ID,
					SourceIngredientName: testIngredient.Name,
					SourceQuantity:       100.0,
					SourceUnit:           string(testIngredient.Unit),
					BatchQuantity:        500.0,
					WastageRate:          0.1, // 10% wastage
				},
			},
			LowStockThreshold:  200.0,
			ExpiryWarningHours: 4,
			CreatedAt:          time.Now(),
			UpdatedAt:          time.Now(),
		}
		err = batchDefinitionRepo.Create(ctx, batchDef)
		require.NoError(t, err)

		// Calculate expected cost
		// 100g * (1 + 0.1) = 110g needed
		// 110g * 0.05 per gram = 5.5 total cost
		// 5.5 / 500ml = 0.011 per ml
		expectedTotalCost := 110.0 * 0.05
		expectedCostPerUnit := expectedTotalCost / 500.0

		// Test cost calculation (without actually creating batch to avoid transaction)
		costBreakdown, err := batchCostCalculator.CalculateBatchCost(ctx, batchDef, 500.0)
		require.NoError(t, err)
		
		assert.InDelta(t, expectedTotalCost, costBreakdown.TotalCost, 0.01)
		assert.InDelta(t, expectedCostPerUnit, costBreakdown.CostPerUnit, 0.001)
		assert.Len(t, costBreakdown.IngredientCosts, 1)
		assert.InDelta(t, 110.0, costBreakdown.IngredientCosts[0].Quantity, 0.01)
	})

	t.Run("DeleteBatch_ValidatesPartiallyUsedBatch", func(t *testing.T) {
		// Create a batch record that has been partially used
		testIngredient := &ingredient.Ingredient{
			Name:        "Test Ingredient",
			Category:    "Test",
			Unit:        ingredient.UnitGram,
			Quantity:    1000.0,
			MinStock:    100.0,
			CostPerUnit: 0.05,
			Supplier:    "Test Supplier",
			CreatedAt:   time.Now(),
			UpdatedAt:   time.Now(),
		}
		err := ingredientRepo.Create(ctx, testIngredient)
		require.NoError(t, err)

		batchRecord := &batch.BatchRecord{
			BatchDefinitionID: primitive.NewObjectID(),
			BatchName:         "Test Batch",
			QuantityProduced:  100.0,
			QuantityRemaining: 50.0, // Partially used
			Unit:              "ml",
			CostPerUnit:       0.1,
			TotalCost:         10.0,
			PreparedBy:        "test-user",
			PreparedAt:        time.Now(),
			ExpiresAt:         time.Now().Add(24 * time.Hour),
			Status:            batch.BatchStatusAvailable,
			IngredientsUsed: []batch.IngredientUsage{
				{
					IngredientID:   testIngredient.ID,
					IngredientName: testIngredient.Name,
					Quantity:       100.0,
					Unit:           string(testIngredient.Unit),
					CostPerUnit:    0.05,
					TotalCost:      5.0,
				},
			},
			CreatedAt: time.Now(),
			UpdatedAt: time.Now(),
		}
		err = batchRecordRepo.Create(ctx, batchRecord)
		require.NoError(t, err)

		// Attempt to delete partially used batch (should fail)
		err = batchRecordService.Delete(ctx, batchRecord.ID)
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "partially used")

		// Verify batch still exists
		existingBatch, err := batchRecordRepo.FindByID(ctx, batchRecord.ID)
		require.NoError(t, err)
		assert.NotNil(t, existingBatch)
	})
}

// Note: Full transaction tests with rollback require MongoDB replica set
// To test transactions locally, set up a replica set:
// 1. Start MongoDB with replica set: mongod --replSet rs0
// 2. Initialize replica set: rs.initiate()
// 3. Run tests with: go test -v -run TestBatchInventoryIntegration
