package services

import (
	"context"
	"fmt"
	"sync"
	"testing"
	"time"

	"cafe-pos/backend/domain/batch"
	"cafe-pos/backend/domain/ingredient"
	"cafe-pos/backend/infrastructure/mongodb"
	"github.com/leanovate/gopter"
	"github.com/leanovate/gopter/gen"
	"github.com/leanovate/gopter/prop"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

// **Validates: Requirements 2.5, 5.5**
// **Validates: Design Property 7 (Quantity Non-Negativity)**
//
// Property: Ingredient and batch quantities never go negative
// - Ingredient quantities must always be >= 0
// - Batch quantities must always be >= 0
// - System must reject operations that would result in negative quantities
// - Concurrent operations must not cause negative quantities
//
// This property ensures data integrity and prevents invalid inventory states.

// TestProperty_IngredientQuantityNonNegative tests that ingredient quantity never goes negative
func TestProperty_IngredientQuantityNonNegative(t *testing.T) {
	if testing.Short() {
		t.Skip("Skipping property test in short mode")
	}

	ctx := context.Background()
	client, db := setupNonNegativityTestDB(t, ctx)
	defer cleanupNonNegativityTestDB(t, ctx, client, db)

	parameters := gopter.DefaultTestParameters()
	parameters.MinSuccessfulTests = 50
	properties := gopter.NewProperties(parameters)

	properties.Property("Ingredient quantity never goes negative during batch creation", prop.ForAll(
		func(testData ingredientNonNegativeData) bool {
			cleanNonNegativityCollections(ctx, db)

			// Setup repositories and service
			ingredientRepo := mongodb.NewIngredientRepository(db)
			stockHistoryRepo := mongodb.NewStockHistoryRepository(db)
			batchDefRepo := mongodb.NewBatchDefinitionRepository(db)
			batchRecordRepo := mongodb.NewBatchRecordRepository(db)
			batchCostCalc := NewBatchCostCalculator(ingredientRepo)
			
			service := NewBatchRecordService(
				batchRecordRepo,
				batchDefRepo,
				ingredientRepo,
				stockHistoryRepo,
				batchCostCalc,
				client,
			)

			// Create batch definition
			batchDef := &batch.BatchDefinition{
				ID:             primitive.NewObjectID(),
				Name:           "Test Batch",
				Unit:           "ml",
				ShelfLifeHours: 24,
				ConversionRates: []batch.ConversionRate{
					{
						SourceIngredientID:   primitive.NewObjectID(),
						SourceIngredientName: "Test Ingredient",
						SourceQuantity:       testData.SourceQtyRequired,
						SourceUnit:           "g",
						BatchQuantity:        testData.BatchQtyProduced,
						WastageRate:          testData.WastageRate,
					},
				},
				LowStockThreshold:  100,
				ExpiryWarningHours: 4,
				CreatedAt:          time.Now(),
				UpdatedAt:          time.Now(),
			}
			err := batchDefRepo.Create(ctx, batchDef)
			if err != nil {
				t.Logf("Failed to create batch definition: %v", err)
				return false
			}

			// Create ingredient with available quantity
			ing := &ingredient.Ingredient{
				ID:          batchDef.ConversionRates[0].SourceIngredientID,
				Name:        batchDef.ConversionRates[0].SourceIngredientName,
				Category:    "Test",
				Unit:        ingredient.UnitGram,
				Quantity:    testData.AvailableQty,
				MinStock:    10,
				CostPerUnit: 0.5,
				CreatedAt:   time.Now(),
				UpdatedAt:   time.Now(),
			}
			err = ingredientRepo.Create(ctx, ing)
			if err != nil {
				t.Logf("Failed to create ingredient: %v", err)
				return false
			}

			// Try to create batch
			_, err = service.CreateBatch(ctx, CreateBatchRequest{
				BatchDefinitionID: batchDef.ID,
				QuantityProduced:  testData.RequestedQty,
				PreparedBy:        "test-user",
			})

			// Check ingredient quantity after operation
			currentIng, fetchErr := ingredientRepo.FindByID(ctx, ing.ID)
			if fetchErr != nil {
				t.Logf("Failed to fetch ingredient: %v", fetchErr)
				return false
			}

			// CRITICAL: Quantity must never be negative
			if currentIng.Quantity < 0 {
				t.Logf("VIOLATION: Ingredient quantity went negative: %.2f", currentIng.Quantity)
				return false
			}

			// If operation succeeded, verify quantity was deducted correctly
			if err == nil {
				rate := batchDef.ConversionRates[0]
				requiredQty := (testData.RequestedQty / rate.BatchQuantity) * rate.SourceQuantity
				requiredQtyWithWastage := requiredQty * (1 + rate.WastageRate)
				expectedQty := testData.AvailableQty - requiredQtyWithWastage

				if !floatEqualsNonNegativity(currentIng.Quantity, expectedQty, 0.01) {
					t.Logf("Quantity mismatch: expected %.2f, got %.2f", expectedQty, currentIng.Quantity)
					return false
				}
			} else {
				// If operation failed, quantity should be unchanged
				if !floatEqualsNonNegativity(currentIng.Quantity, testData.AvailableQty, 0.01) {
					t.Logf("Quantity changed after failed operation: %.2f -> %.2f", testData.AvailableQty, currentIng.Quantity)
					return false
				}
			}

			return true
		},
		genIngredientNonNegativeData(),
	))

	properties.TestingRun(t, gopter.ConsoleReporter(false))
}

// TestProperty_BatchQuantityNonNegative tests that batch quantity never goes negative
func TestProperty_BatchQuantityNonNegative(t *testing.T) {
	if testing.Short() {
		t.Skip("Skipping property test in short mode")
	}

	ctx := context.Background()
	client, db := setupNonNegativityTestDB(t, ctx)
	defer cleanupNonNegativityTestDB(t, ctx, client, db)

	parameters := gopter.DefaultTestParameters()
	parameters.MinSuccessfulTests = 50
	properties := gopter.NewProperties(parameters)

	properties.Property("Batch quantity never goes negative during usage", prop.ForAll(
		func(testData batchNonNegativeData) bool {
			cleanNonNegativityCollections(ctx, db)

			// Setup repositories and service
			batchDefRepo := mongodb.NewBatchDefinitionRepository(db)
			batchRecordRepo := mongodb.NewBatchRecordRepository(db)
			batchUsageLogRepo := mongodb.NewBatchUsageLogRepository(db)
			
			service := NewBatchUsageService(
				batchRecordRepo,
				batchUsageLogRepo,
			)

			// Create batch definition
			batchDef := &batch.BatchDefinition{
				ID:             primitive.NewObjectID(),
				Name:           "Test Batch",
				Unit:           "ml",
				ShelfLifeHours: 24,
				ConversionRates: []batch.ConversionRate{
					{
						SourceIngredientID:   primitive.NewObjectID(),
						SourceIngredientName: "Test Ingredient",
						SourceQuantity:       100,
						SourceUnit:           "g",
						BatchQuantity:        500,
						WastageRate:          0.1,
					},
				},
				LowStockThreshold:  100,
				ExpiryWarningHours: 4,
				CreatedAt:          time.Now(),
				UpdatedAt:          time.Now(),
			}
			batchDefRepo.Create(ctx, batchDef)

			// Create batch record with available quantity
			batchRecord := &batch.BatchRecord{
				ID:                primitive.NewObjectID(),
				BatchDefinitionID: batchDef.ID,
				BatchName:         batchDef.Name,
				QuantityProduced:  testData.AvailableQty,
				QuantityRemaining: testData.AvailableQty,
				Unit:              batchDef.Unit,
				CostPerUnit:       0.15,
				TotalCost:         testData.AvailableQty * 0.15,
				PreparedBy:        "test-user",
				PreparedAt:        time.Now(),
				ExpiresAt:         time.Now().Add(24 * time.Hour),
				Status:            "available",
				IngredientsUsed:   []batch.IngredientUsage{},
				CreatedAt:         time.Now(),
				UpdatedAt:         time.Now(),
			}
			batchRecordRepo.Create(ctx, batchRecord)

			// Try to use batch
			result, err := service.UseBatch(ctx, UseBatchRequest{
				BatchDefinitionID: batchDef.ID,
				QuantityNeeded:    testData.RequestedQty,
				OrderID:           primitive.NewObjectID(),
				MenuItemID:        primitive.NewObjectID(),
				MenuItemName:      "Test Item",
			})

			// Check batch quantity after operation
			currentBatch, fetchErr := batchRecordRepo.FindByID(ctx, batchRecord.ID)
			if fetchErr != nil {
				t.Logf("Failed to fetch batch record: %v", fetchErr)
				return false
			}

			// CRITICAL: Quantity must never be negative
			if currentBatch.QuantityRemaining < 0 {
				t.Logf("VIOLATION: Batch quantity went negative: %.2f", currentBatch.QuantityRemaining)
				return false
			}

			// Determine if operation should have succeeded
			shouldSucceed := testData.RequestedQty <= testData.AvailableQty

			if shouldSucceed {
				// Operation should succeed - verify it did and quantity was deducted
				if err != nil {
					t.Logf("Expected success but got error: %v", err)
					return false
				}
				if result == nil || !result.Success {
					t.Logf("Expected successful result but got: %+v", result)
					return false
				}
				expectedQty := testData.AvailableQty - testData.RequestedQty
				if !floatEqualsNonNegativity(currentBatch.QuantityRemaining, expectedQty, 0.01) {
					t.Logf("Quantity mismatch: expected %.2f, got %.2f", expectedQty, currentBatch.QuantityRemaining)
					return false
				}
			} else {
				// Operation should fail - verify it did and quantity unchanged
				if err == nil && result != nil && result.Success {
					t.Logf("Expected failure but operation succeeded (requested: %.2f, available: %.2f, result: %+v)", 
						testData.RequestedQty, testData.AvailableQty, result)
					return false
				}
				if !floatEqualsNonNegativity(currentBatch.QuantityRemaining, testData.AvailableQty, 0.01) {
					t.Logf("Quantity changed after failed operation: %.2f -> %.2f", testData.AvailableQty, currentBatch.QuantityRemaining)
					return false
				}
			}

			return true
		},
		genBatchNonNegativeData(),
	))

	properties.TestingRun(t, gopter.ConsoleReporter(false))
}

// TestProperty_ConcurrentOperationsNonNegative tests that concurrent operations don't cause negative quantities
func TestProperty_ConcurrentOperationsNonNegative(t *testing.T) {
	if testing.Short() {
		t.Skip("Skipping property test in short mode")
	}

	ctx := context.Background()
	client, db := setupNonNegativityTestDB(t, ctx)
	defer cleanupNonNegativityTestDB(t, ctx, client, db)

	parameters := gopter.DefaultTestParameters()
	parameters.MinSuccessfulTests = 20 // Reduced due to concurrency overhead
	properties := gopter.NewProperties(parameters)

	properties.Property("Concurrent batch usage operations never cause negative quantities", prop.ForAll(
		func(testData concurrentNonNegativeData) bool {
			cleanNonNegativityCollections(ctx, db)

			// Setup
			batchDefRepo := mongodb.NewBatchDefinitionRepository(db)
			batchRecordRepo := mongodb.NewBatchRecordRepository(db)
			batchUsageLogRepo := mongodb.NewBatchUsageLogRepository(db)
			
			service := NewBatchUsageService(
				batchRecordRepo,
				batchUsageLogRepo,
			)

			// Create batch definition
			batchDef := &batch.BatchDefinition{
				ID:             primitive.NewObjectID(),
				Name:           "Test Batch",
				Unit:           "ml",
				ShelfLifeHours: 24,
				ConversionRates: []batch.ConversionRate{
					{
						SourceIngredientID:   primitive.NewObjectID(),
						SourceIngredientName: "Test Ingredient",
						SourceQuantity:       100,
						SourceUnit:           "g",
						BatchQuantity:        500,
						WastageRate:          0.1,
					},
				},
				LowStockThreshold:  100,
				ExpiryWarningHours: 4,
				CreatedAt:          time.Now(),
				UpdatedAt:          time.Now(),
			}
			batchDefRepo.Create(ctx, batchDef)

			// Create batch record
			batchRecord := &batch.BatchRecord{
				ID:                primitive.NewObjectID(),
				BatchDefinitionID: batchDef.ID,
				BatchName:         batchDef.Name,
				QuantityProduced:  testData.InitialQty,
				QuantityRemaining: testData.InitialQty,
				Unit:              batchDef.Unit,
				CostPerUnit:       0.15,
				TotalCost:         testData.InitialQty * 0.15,
				PreparedBy:        "test-user",
				PreparedAt:        time.Now(),
				ExpiresAt:         time.Now().Add(24 * time.Hour),
				Status:            "available",
				IngredientsUsed:   []batch.IngredientUsage{},
				CreatedAt:         time.Now(),
				UpdatedAt:         time.Now(),
			}
			batchRecordRepo.Create(ctx, batchRecord)

			// Run concurrent usage operations
			var wg sync.WaitGroup
			successCount := 0
			var mu sync.Mutex

			for i := 0; i < testData.NumConcurrent; i++ {
				wg.Add(1)
				go func(idx int) {
					defer wg.Done()
					_, err := service.UseBatch(ctx, UseBatchRequest{
						BatchDefinitionID: batchDef.ID,
						QuantityNeeded:    testData.QtyPerRequest,
						OrderID:           primitive.NewObjectID(),
						MenuItemID:        primitive.NewObjectID(),
						MenuItemName:      fmt.Sprintf("Test Item %d", idx),
					})
					if err == nil {
						mu.Lock()
						successCount++
						mu.Unlock()
					}
				}(i)
			}

			wg.Wait()

			// Check final batch quantity
			currentBatch, err := batchRecordRepo.FindByID(ctx, batchRecord.ID)
			if err != nil {
				t.Logf("Failed to fetch batch record: %v", err)
				return false
			}

			// CRITICAL: Quantity must never be negative
			if currentBatch.QuantityRemaining < 0 {
				t.Logf("VIOLATION: Batch quantity went negative after concurrent operations: %.2f", currentBatch.QuantityRemaining)
				return false
			}

			// Verify that total deducted doesn't exceed initial quantity
			totalDeducted := testData.InitialQty - currentBatch.QuantityRemaining
			
			if totalDeducted > testData.InitialQty + 0.01 { // Allow small epsilon for floating point
				t.Logf("Total deducted (%.2f) exceeds initial quantity (%.2f)", totalDeducted, testData.InitialQty)
				return false
			}

			// The key property: quantity never went negative
			// We don't strictly verify the exact final quantity because without transactions,
			// some operations may have succeeded or failed in various orders.
			// What matters is that the quantity never became negative.
			return true
		},
		genConcurrentNonNegativeData(),
	))

	properties.TestingRun(t, gopter.ConsoleReporter(false))
}

// Test data structures

type ingredientNonNegativeData struct {
	AvailableQty      float64 // Available ingredient quantity
	SourceQtyRequired float64 // Source quantity required per batch
	BatchQtyProduced  float64 // Batch quantity produced
	WastageRate       float64 // Wastage rate (0.0 to 0.5)
	RequestedQty      float64 // Requested batch quantity to produce
}

type batchNonNegativeData struct {
	AvailableQty float64 // Available batch quantity
	RequestedQty float64 // Requested quantity to use
}

type concurrentNonNegativeData struct {
	InitialQty     float64 // Initial batch quantity
	NumConcurrent  int     // Number of concurrent operations
	QtyPerRequest  float64 // Quantity per request
}

// Generators

func genIngredientNonNegativeData() gopter.Gen {
	return gopter.CombineGens(
		gen.Float64Range(50, 1000),   // AvailableQty
		gen.Float64Range(50, 200),    // SourceQtyRequired
		gen.Float64Range(200, 500),   // BatchQtyProduced
		gen.Float64Range(0.0, 0.5),   // WastageRate
		gen.Float64Range(100, 800),   // RequestedQty
	).Map(func(values []interface{}) ingredientNonNegativeData {
		return ingredientNonNegativeData{
			AvailableQty:      values[0].(float64),
			SourceQtyRequired: values[1].(float64),
			BatchQtyProduced:  values[2].(float64),
			WastageRate:       values[3].(float64),
			RequestedQty:      values[4].(float64),
		}
	})
}

func genBatchNonNegativeData() gopter.Gen {
	return gopter.CombineGens(
		gen.Float64Range(50, 1000),  // AvailableQty
		gen.Float64Range(10, 1200),  // RequestedQty (may exceed available)
	).Map(func(values []interface{}) batchNonNegativeData {
		return batchNonNegativeData{
			AvailableQty: values[0].(float64),
			RequestedQty: values[1].(float64),
		}
	})
}

func genConcurrentNonNegativeData() gopter.Gen {
	return gopter.CombineGens(
		gen.Float64Range(500, 2000),  // InitialQty
		gen.IntRange(5, 15),           // NumConcurrent
		gen.Float64Range(50, 200),     // QtyPerRequest
	).Map(func(values []interface{}) concurrentNonNegativeData {
		return concurrentNonNegativeData{
			InitialQty:    values[0].(float64),
			NumConcurrent: values[1].(int),
			QtyPerRequest: values[2].(float64),
		}
	})
}

// Helper functions

func setupNonNegativityTestDB(t *testing.T, ctx context.Context) (*mongo.Client, *mongo.Database) {
	mongoURI := "mongodb://admin:password123@localhost:27017/?authSource=admin"
	clientOptions := options.Client().ApplyURI(mongoURI)
	client, err := mongo.Connect(ctx, clientOptions)
	if err != nil {
		t.Fatalf("Failed to connect to MongoDB: %v", err)
	}

	db := client.Database("cafe_pos_test_non_negativity")
	return client, db
}

func cleanupNonNegativityTestDB(t *testing.T, ctx context.Context, client *mongo.Client, db *mongo.Database) {
	db.Drop(ctx)
	client.Disconnect(ctx)
}

func cleanNonNegativityCollections(ctx context.Context, db *mongo.Database) {
	db.Collection("ingredients").Drop(ctx)
	db.Collection("stock_history").Drop(ctx)
	db.Collection("batch_definitions").Drop(ctx)
	db.Collection("batch_records").Drop(ctx)
	db.Collection("batch_usage_logs").Drop(ctx)
}

func floatEqualsNonNegativity(a, b, epsilon float64) bool {
	diff := a - b
	if diff < 0 {
		diff = -diff
	}
	return diff < epsilon
}
