package services

import (
	"context"
	"testing"
	"time"

	"cafe-pos/backend/domain/batch"
	"cafe-pos/backend/domain/ingredient"
	"cafe-pos/backend/domain/menu"
	"cafe-pos/backend/infrastructure/mongodb"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

// Adapter for MenuBatchDefinitionRepository
type menuBatchDefRepoAdapter struct {
	repo *mongodb.BatchDefinitionRepository
}

func (a *menuBatchDefRepoAdapter) FindByID(ctx context.Context, id primitive.ObjectID) (interface{}, error) {
	return a.repo.FindByID(ctx, id)
}

// Adapter for MenuBatchRecordRepository
type menuBatchRecRepoAdapter struct {
	repo *mongodb.BatchRecordRepository
}

func (a *menuBatchRecRepoAdapter) GetTotalAvailableQuantity(ctx context.Context, defID primitive.ObjectID) (float64, error) {
	return a.repo.GetTotalAvailableQuantity(ctx, defID)
}

// Adapter for CostCalculatorBatchRecordRepository
type costCalcBatchRecRepoAdapter struct {
	repo *mongodb.BatchRecordRepository
}

func (a *costCalcBatchRecRepoAdapter) FindAvailableByDefinition(ctx context.Context, defID primitive.ObjectID) ([]*CostCalculatorBatchRecord, error) {
	records, err := a.repo.FindAvailableByDefinition(ctx, defID)
	if err != nil {
		return nil, err
	}
	
	// Convert to CostCalculatorBatchRecord
	result := make([]*CostCalculatorBatchRecord, len(records))
	for i, r := range records {
		result[i] = &CostCalculatorBatchRecord{
			CostPerUnit: r.CostPerUnit,
			Unit:        r.Unit,
		}
	}
	return result, nil
}

// isReplicaSet checks if MongoDB is running as a replica set
func isReplicaSet(t *testing.T, client *mongo.Client) bool {
	var result struct {
		Hosts []string `bson:"hosts"`
	}
	err := client.Database("admin").RunCommand(context.Background(), map[string]interface{}{"isMaster": 1}).Decode(&result)
	if err != nil {
		return false
	}
	return len(result.Hosts) > 0
}

// TestMenuBatchIntegration tests the integration between menu system and batch ingredients
// Requirements: 5.1, 5.6, 3.3
func TestMenuBatchIntegration(t *testing.T) {
	// Setup MongoDB test database
	ctx := context.Background()
	
	// Get MongoDB URI from environment or use default with auth
	mongoURI := "mongodb://admin:password123@localhost:27017/?authSource=admin"
	client, err := mongo.Connect(ctx, options.Client().ApplyURI(mongoURI))
	require.NoError(t, err)
	defer client.Disconnect(ctx)

	// Use a test database
	testDB := client.Database("cafe_pos_test_menu_batch_" + primitive.NewObjectID().Hex())
	defer testDB.Drop(ctx)

	// Setup repositories
	menuRepo := mongodb.NewMenuRepository(testDB)
	ingredientRepo := mongodb.NewIngredientRepository(testDB)
	batchDefinitionRepo := mongodb.NewBatchDefinitionRepository(testDB)
	batchRecordRepo := mongodb.NewBatchRecordRepository(testDB)
	stockHistoryRepo := mongodb.NewStockHistoryRepository(testDB)

	// Setup services
	menuService := NewMenuService(menuRepo)
	menuService.SetBatchRepositories(
		&menuBatchDefRepoAdapter{repo: batchDefinitionRepo},
		&menuBatchRecRepoAdapter{repo: batchRecordRepo},
	)
	
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
	
	costCalculatorService := NewCostCalculatorService(menuRepo, ingredientRepo, nil, nil)
	costCalculatorService.SetBatchRecordRepository(&costCalcBatchRecRepoAdapter{repo: batchRecordRepo})

	t.Run("CreateMenuItem_WithBatchIngredient_ValidatesSchema", func(t *testing.T) {
		// Create a batch definition
		batchDef := &batch.BatchDefinition{
			Name:           "Coffee Concentrate",
			Unit:           "ml",
			ShelfLifeHours: 24,
			ConversionRates: []batch.ConversionRate{
				{
					SourceIngredientID:   primitive.NewObjectID(),
					SourceIngredientName: "Coffee Beans",
					SourceQuantity:       100.0,
					SourceUnit:           "g",
					BatchQuantity:        500.0,
					WastageRate:          0.1,
				},
			},
			LowStockThreshold:  200.0,
			ExpiryWarningHours: 4,
			CreatedAt:          time.Now(),
			UpdatedAt:          time.Now(),
		}
		err := batchDefinitionRepo.Create(ctx, batchDef)
		require.NoError(t, err)

		// Create menu item with batch ingredient
		req := &menu.CreateMenuItemRequest{
			Name:        "Iced Coffee",
			Category:    "Beverages",
			Description: "Cold coffee drink",
			Price:       45000,
			Ingredients: []menu.Ingredient{
				{
					Name:           "Coffee Concentrate",
					Quantity:       30.0,
					Unit:           ingredient.UnitMilliliter,
					IngredientType: string(menu.IngredientTypeBatch),
					BatchID:        &batchDef.ID,
				},
			},
		}

		menuItem, err := menuService.CreateMenuItem(ctx, req)
		require.NoError(t, err)
		assert.NotNil(t, menuItem)
		assert.Equal(t, "Iced Coffee", menuItem.Name)
		assert.Len(t, menuItem.Ingredients, 1)
		assert.Equal(t, string(menu.IngredientTypeBatch), menuItem.Ingredients[0].IngredientType)
		assert.NotNil(t, menuItem.Ingredients[0].BatchID)
		assert.Equal(t, batchDef.ID, *menuItem.Ingredients[0].BatchID)
	})

	t.Run("CreateMenuItem_WithBatchIngredient_MissingBatchID_Fails", func(t *testing.T) {
		// Attempt to create menu item with batch ingredient but no BatchID
		req := &menu.CreateMenuItemRequest{
			Name:        "Invalid Coffee",
			Category:    "Beverages",
			Description: "Should fail",
			Price:       45000,
			Ingredients: []menu.Ingredient{
				{
					Name:           "Coffee Concentrate",
					Quantity:       30.0,
					Unit:           ingredient.UnitMilliliter,
					IngredientType: string(menu.IngredientTypeBatch),
					BatchID:        nil, // Missing BatchID
				},
			},
		}

		_, err := menuService.CreateMenuItem(ctx, req)
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "missing batch_id")
	})

	t.Run("CreateMenuItem_WithBatchIngredient_NonExistentBatch_Fails", func(t *testing.T) {
		// Attempt to create menu item with non-existent batch
		nonExistentBatchID := primitive.NewObjectID()
		req := &menu.CreateMenuItemRequest{
			Name:        "Invalid Coffee",
			Category:    "Beverages",
			Description: "Should fail",
			Price:       45000,
			Ingredients: []menu.Ingredient{
				{
					Name:           "Coffee Concentrate",
					Quantity:       30.0,
					Unit:           ingredient.UnitMilliliter,
					IngredientType: string(menu.IngredientTypeBatch),
					BatchID:        &nonExistentBatchID,
				},
			},
		}

		_, err := menuService.CreateMenuItem(ctx, req)
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "batch definition not found")
	})

	t.Run("CreateMenuItem_WithMixedIngredients_RawAndBatch", func(t *testing.T) {
		// Create raw ingredient
		rawIngredient := &ingredient.Ingredient{
			Name:        "Milk",
			Category:    "Dairy",
			Unit:        ingredient.UnitMilliliter,
			Quantity:    5000.0,
			MinStock:    1000.0,
			CostPerUnit: 0.02,
			Supplier:    "Test Supplier",
			CreatedAt:   time.Now(),
			UpdatedAt:   time.Now(),
		}
		err := ingredientRepo.Create(ctx, rawIngredient)
		require.NoError(t, err)

		// Create batch definition
		batchDef := &batch.BatchDefinition{
			Name:           "Espresso Concentrate",
			Unit:           "ml",
			ShelfLifeHours: 12,
			ConversionRates: []batch.ConversionRate{
				{
					SourceIngredientID:   primitive.NewObjectID(),
					SourceIngredientName: "Espresso Beans",
					SourceQuantity:       50.0,
					SourceUnit:           "g",
					BatchQuantity:        200.0,
					WastageRate:          0.05,
				},
			},
			LowStockThreshold:  100.0,
			ExpiryWarningHours: 2,
			CreatedAt:          time.Now(),
			UpdatedAt:          time.Now(),
		}
		err = batchDefinitionRepo.Create(ctx, batchDef)
		require.NoError(t, err)

		// Create menu item with both raw and batch ingredients
		req := &menu.CreateMenuItemRequest{
			Name:        "Latte",
			Category:    "Beverages",
			Description: "Espresso with milk",
			Price:       55000,
			Ingredients: []menu.Ingredient{
				{
					Name:           "Espresso Concentrate",
					Quantity:       30.0,
					Unit:           ingredient.UnitMilliliter,
					IngredientType: string(menu.IngredientTypeBatch),
					BatchID:        &batchDef.ID,
				},
				{
					Name:           "Milk",
					Quantity:       200.0,
					Unit:           ingredient.UnitMilliliter,
					IngredientType: string(menu.IngredientTypeRaw),
					BatchID:        nil,
				},
			},
		}

		menuItem, err := menuService.CreateMenuItem(ctx, req)
		require.NoError(t, err)
		assert.NotNil(t, menuItem)
		assert.Len(t, menuItem.Ingredients, 2)
		
		// Verify batch ingredient
		assert.Equal(t, string(menu.IngredientTypeBatch), menuItem.Ingredients[0].IngredientType)
		assert.NotNil(t, menuItem.Ingredients[0].BatchID)
		
		// Verify raw ingredient
		assert.Equal(t, string(menu.IngredientTypeRaw), menuItem.Ingredients[1].IngredientType)
		assert.Nil(t, menuItem.Ingredients[1].BatchID)
	})

	t.Run("CalculateMenuItemCost_WithBatchIngredient", func(t *testing.T) {
		// Skip if not running on replica set (transactions required)
		if !isReplicaSet(t, client) {
			t.Skip("Skipping test that requires MongoDB replica set for transactions")
		}
		
		// Create raw ingredient for batch
		rawIngredient := &ingredient.Ingredient{
			Name:        "Tea Leaves",
			Category:    "Raw Materials",
			Unit:        ingredient.UnitGram,
			Quantity:    2000.0,
			MinStock:    500.0,
			CostPerUnit: 0.10,
			Supplier:    "Test Supplier",
			CreatedAt:   time.Now(),
			UpdatedAt:   time.Now(),
		}
		err := ingredientRepo.Create(ctx, rawIngredient)
		require.NoError(t, err)

		// Create batch definition
		batchDef := &batch.BatchDefinition{
			Name:           "Tea Concentrate",
			Unit:           "ml",
			ShelfLifeHours: 24,
			ConversionRates: []batch.ConversionRate{
				{
					SourceIngredientID:   rawIngredient.ID,
					SourceIngredientName: rawIngredient.Name,
					SourceQuantity:       100.0,
					SourceUnit:           string(rawIngredient.Unit),
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

		// Create batch record
		batchReq := CreateBatchRequest{
			BatchDefinitionID: batchDef.ID,
			QuantityProduced:  500.0,
			PreparedBy:        "test-user",
		}
		batchRecord, err := batchRecordService.CreateBatch(ctx, batchReq)
		require.NoError(t, err)
		
		// Expected batch cost: 100g * 1.1 (wastage) * 0.10 = 11.0 total
		// Cost per unit: 11.0 / 500ml = 0.022 per ml
		expectedBatchCostPerUnit := 0.022

		// Create menu item with batch ingredient
		menuReq := &menu.CreateMenuItemRequest{
			Name:        "Iced Tea",
			Category:    "Beverages",
			Description: "Cold tea drink",
			Price:       35000,
			Ingredients: []menu.Ingredient{
				{
					Name:           "Tea Concentrate",
					Quantity:       50.0,
					Unit:           ingredient.UnitMilliliter,
					IngredientType: string(menu.IngredientTypeBatch),
					BatchID:        &batchDef.ID,
				},
			},
		}

		menuItem, err := menuService.CreateMenuItem(ctx, menuReq)
		require.NoError(t, err)

		// Calculate cost
		costResult, err := costCalculatorService.CalculateMenuItemCost(ctx, menuItem.ID)
		require.NoError(t, err)
		
		// Expected menu item cost: 50ml * 0.022 = 1.1
		expectedMenuItemCost := 50.0 * expectedBatchCostPerUnit
		
		assert.InDelta(t, expectedMenuItemCost, costResult.CurrentCost, 0.01)
		assert.Equal(t, menu.CostStatusFinal, costResult.CostStatus)
		
		// Verify batch record cost
		assert.InDelta(t, expectedBatchCostPerUnit, batchRecord.CostPerUnit, 0.001)
	})

	t.Run("CalculateMenuItemCost_WithMixedIngredients", func(t *testing.T) {
		// Skip if not running on replica set (transactions required)
		if !isReplicaSet(t, client) {
			t.Skip("Skipping test that requires MongoDB replica set for transactions")
		}
		
		// Create raw ingredients
		sugar := &ingredient.Ingredient{
			Name:        "Sugar",
			Category:    "Sweeteners",
			Unit:        ingredient.UnitGram,
			Quantity:    5000.0,
			MinStock:    1000.0,
			CostPerUnit: 0.01,
			Supplier:    "Test Supplier",
			CreatedAt:   time.Now(),
			UpdatedAt:   time.Now(),
		}
		err := ingredientRepo.Create(ctx, sugar)
		require.NoError(t, err)

		coffeeBeans := &ingredient.Ingredient{
			Name:        "Premium Coffee Beans",
			Category:    "Raw Materials",
			Unit:        ingredient.UnitGram,
			Quantity:    3000.0,
			MinStock:    500.0,
			CostPerUnit: 0.15,
			Supplier:    "Test Supplier",
			CreatedAt:   time.Now(),
			UpdatedAt:   time.Now(),
		}
		err = ingredientRepo.Create(ctx, coffeeBeans)
		require.NoError(t, err)

		// Create batch definition
		batchDef := &batch.BatchDefinition{
			Name:           "Premium Coffee Concentrate",
			Unit:           "ml",
			ShelfLifeHours: 24,
			ConversionRates: []batch.ConversionRate{
				{
					SourceIngredientID:   coffeeBeans.ID,
					SourceIngredientName: coffeeBeans.Name,
					SourceQuantity:       100.0,
					SourceUnit:           string(coffeeBeans.Unit),
					BatchQuantity:        400.0,
					WastageRate:          0.15, // 15% wastage
				},
			},
			LowStockThreshold:  200.0,
			ExpiryWarningHours: 4,
			CreatedAt:          time.Now(),
			UpdatedAt:          time.Now(),
		}
		err = batchDefinitionRepo.Create(ctx, batchDef)
		require.NoError(t, err)

		// Create batch record
		batchReq := CreateBatchRequest{
			BatchDefinitionID: batchDef.ID,
			QuantityProduced:  400.0,
			PreparedBy:        "test-user",
		}
		batchRecord, err := batchRecordService.CreateBatch(ctx, batchReq)
		require.NoError(t, err)
		
		// Expected batch cost: 100g * 1.15 (wastage) * 0.15 = 17.25 total
		// Cost per unit: 17.25 / 400ml = 0.043125 per ml

		// Create menu item with mixed ingredients
		menuReq := &menu.CreateMenuItemRequest{
			Name:        "Sweet Coffee",
			Category:    "Beverages",
			Description: "Coffee with sugar",
			Price:       50000,
			Ingredients: []menu.Ingredient{
				{
					Name:           "Premium Coffee Concentrate",
					Quantity:       40.0,
					Unit:           ingredient.UnitMilliliter,
					IngredientType: string(menu.IngredientTypeBatch),
					BatchID:        &batchDef.ID,
				},
				{
					Name:           "Sugar",
					Quantity:       10.0,
					Unit:           ingredient.UnitGram,
					IngredientType: string(menu.IngredientTypeRaw),
					BatchID:        nil,
				},
			},
		}

		menuItem, err := menuService.CreateMenuItem(ctx, menuReq)
		require.NoError(t, err)

		// Calculate cost
		costResult, err := costCalculatorService.CalculateMenuItemCost(ctx, menuItem.ID)
		require.NoError(t, err)
		
		// Expected costs:
		// Batch: 40ml * 0.043125 = 1.725
		// Sugar: 10g * 0.01 = 0.1
		// Total: 1.825
		expectedBatchCost := 40.0 * batchRecord.CostPerUnit
		expectedSugarCost := 10.0 * 0.01
		expectedTotalCost := expectedBatchCost + expectedSugarCost
		
		assert.InDelta(t, expectedTotalCost, costResult.CurrentCost, 0.01)
		assert.Equal(t, menu.CostStatusFinal, costResult.CostStatus)
	})

	t.Run("CalculateMenuItemCost_WithBatchIngredient_NoBatchesAvailable", func(t *testing.T) {
		// Create batch definition without any batch records
		batchDef := &batch.BatchDefinition{
			Name:           "Unavailable Batch",
			Unit:           "ml",
			ShelfLifeHours: 24,
			ConversionRates: []batch.ConversionRate{
				{
					SourceIngredientID:   primitive.NewObjectID(),
					SourceIngredientName: "Some Ingredient",
					SourceQuantity:       100.0,
					SourceUnit:           "g",
					BatchQuantity:        500.0,
					WastageRate:          0.0,
				},
			},
			LowStockThreshold:  200.0,
			ExpiryWarningHours: 4,
			CreatedAt:          time.Now(),
			UpdatedAt:          time.Now(),
		}
		err := batchDefinitionRepo.Create(ctx, batchDef)
		require.NoError(t, err)

		// Create menu item with batch ingredient (no batches available)
		menuReq := &menu.CreateMenuItemRequest{
			Name:        "Unavailable Drink",
			Category:    "Beverages",
			Description: "Should have incomplete cost",
			Price:       40000,
			Ingredients: []menu.Ingredient{
				{
					Name:           "Unavailable Batch",
					Quantity:       30.0,
					Unit:           ingredient.UnitMilliliter,
					IngredientType: string(menu.IngredientTypeBatch),
					BatchID:        &batchDef.ID,
				},
			},
		}

		menuItem, err := menuService.CreateMenuItem(ctx, menuReq)
		require.NoError(t, err)

		// Calculate cost (should return incomplete status)
		costResult, err := costCalculatorService.CalculateMenuItemCost(ctx, menuItem.ID)
		require.NoError(t, err)
		
		assert.Equal(t, 0.0, costResult.CurrentCost)
		assert.Equal(t, menu.CostStatusIncomplete, costResult.CostStatus)
		assert.Contains(t, costResult.MissingIngredients, "Unavailable Batch")
	})

	t.Run("UpdateMenuItem_ChangeBatchIngredient", func(t *testing.T) {
		// Create two batch definitions
		batchDef1 := &batch.BatchDefinition{
			Name:           "Batch A",
			Unit:           "ml",
			ShelfLifeHours: 24,
			ConversionRates: []batch.ConversionRate{
				{
					SourceIngredientID:   primitive.NewObjectID(),
					SourceIngredientName: "Ingredient A",
					SourceQuantity:       100.0,
					SourceUnit:           "g",
					BatchQuantity:        500.0,
					WastageRate:          0.0,
				},
			},
			LowStockThreshold:  200.0,
			ExpiryWarningHours: 4,
			CreatedAt:          time.Now(),
			UpdatedAt:          time.Now(),
		}
		err := batchDefinitionRepo.Create(ctx, batchDef1)
		require.NoError(t, err)

		batchDef2 := &batch.BatchDefinition{
			Name:           "Batch B",
			Unit:           "ml",
			ShelfLifeHours: 24,
			ConversionRates: []batch.ConversionRate{
				{
					SourceIngredientID:   primitive.NewObjectID(),
					SourceIngredientName: "Ingredient B",
					SourceQuantity:       100.0,
					SourceUnit:           "g",
					BatchQuantity:        500.0,
					WastageRate:          0.0,
				},
			},
			LowStockThreshold:  200.0,
			ExpiryWarningHours: 4,
			CreatedAt:          time.Now(),
			UpdatedAt:          time.Now(),
		}
		err = batchDefinitionRepo.Create(ctx, batchDef2)
		require.NoError(t, err)

		// Create menu item with first batch
		createReq := &menu.CreateMenuItemRequest{
			Name:        "Test Drink",
			Category:    "Beverages",
			Description: "Test",
			Price:       40000,
			Ingredients: []menu.Ingredient{
				{
					Name:           "Batch A",
					Quantity:       30.0,
					Unit:           ingredient.UnitMilliliter,
					IngredientType: string(menu.IngredientTypeBatch),
					BatchID:        &batchDef1.ID,
				},
			},
		}

		menuItem, err := menuService.CreateMenuItem(ctx, createReq)
		require.NoError(t, err)

		// Update to use second batch
		updateReq := &menu.UpdateMenuItemRequest{
			Ingredients: []menu.Ingredient{
				{
					Name:           "Batch B",
					Quantity:       40.0,
					Unit:           ingredient.UnitMilliliter,
					IngredientType: string(menu.IngredientTypeBatch),
					BatchID:        &batchDef2.ID,
				},
			},
		}

		updatedItem, err := menuService.UpdateMenuItem(ctx, menuItem.ID, updateReq)
		require.NoError(t, err)
		assert.Len(t, updatedItem.Ingredients, 1)
		assert.Equal(t, "Batch B", updatedItem.Ingredients[0].Name)
		assert.Equal(t, batchDef2.ID, *updatedItem.Ingredients[0].BatchID)
		assert.Equal(t, 40.0, updatedItem.Ingredients[0].Quantity)
	})
}
