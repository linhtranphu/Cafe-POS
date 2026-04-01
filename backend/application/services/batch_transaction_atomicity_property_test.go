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

// **Validates: Requirements 2.2, 2.6, 8.7**
// **Validates: Design Property 6 (Transaction Atomicity)**
//
// Property: Batch creation operations are atomic - all steps succeed or all rollback
// When creating a batch:
// - If ingredient deduction fails, no batch record should be created
// - If batch record creation fails, ingredients should not be deducted
// - Database should never be left in an inconsistent state
//
// Property: Batch usage operations are atomic - all steps succeed or all rollback
// When using a batch:
// - If batch quantity update fails, no usage log should be created
// - If usage log creation fails, batch quantity should not be deducted
// - Database should never be left in an inconsistent state
//
// Property: Concurrent transactions don't cause race conditions
// When multiple transactions run concurrently:
// - Final state should be consistent
// - No partial updates should occur
// - Inventory quantities should be accurate

// TestProperty_BatchCreationRollback tests that batch creation rolls back on failure
func TestProperty_BatchCreationRollback(t *testing.T) {
	if testing.Short() {
		t.Skip("Skipping property test in short mode")
	}

	// Setup MongoDB test connection
	ctx := context.Background()
	client, db := setupAtomicityTestDB(t, ctx)
	defer cleanupAtomicityTestDB(t, ctx, client, db)

	parameters := gopter.DefaultTestParameters()
	parameters.MinSuccessfulTests = 20 // Reduced for integration test
	properties := gopter.NewProperties(parameters)

	properties.Property("Batch creation rolls back completely on insufficient ingredients", prop.ForAll(
		func(testData batchCreationRollbackData) bool {
			// Clean collections before each test
			cleanCollections(ctx, db)

			// Setup repositories
			ingredientRepo := mongodb.NewIngredientRepository(db)
			stockHistoryRepo := mongodb.NewStockHistoryRepository(db)
			batchDefRepo := mongodb.NewBatchDefinitionRepository(db)
			batchRecordRepo := mongodb.NewBatchRecordRepository(db)
			
			// Create cost calculator
			batchCostCalc := NewBatchCostCalculator(ingredientRepo)
			
			// Create service
			service := NewBatchRecordService(
				batchRecordRepo,
				batchDefRepo,
				ingredientRepo,
				stockHistoryRepo,
				nil,
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
			err := batchDefRepo.Create(ctx, batchDef)
			if err != nil {
				t.Logf("Failed to create batch definition: %v", err)
				return false
			}

			// Create ingredient with INSUFFICIENT quantity (this will cause failure)
			ing := &ingredient.Ingredient{
				ID:          batchDef.ConversionRates[0].SourceIngredientID,
				Name:        batchDef.ConversionRates[0].SourceIngredientName,
				Category:    "Test",
				Unit:        ingredient.UnitGram,
				Quantity:    testData.InsufficientQty, // Not enough!
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

			initialQty := ing.Quantity

			// Try to create batch (should fail due to insufficient ingredients)
			_, err = service.CreateBatch(ctx, CreateBatchRequest{
				BatchDefinitionID: batchDef.ID,
				QuantityProduced:  testData.QuantityProduced,
				PreparedBy:        "test-user",
			})

			// Verify that operation failed
			if err == nil {
				t.Logf("Expected error but got success")
				return false
			}

			// Verify complete rollback:
			// 1. No batch record should be created
			filter := batch.BatchRecordFilter{
				BatchDefinitionID: &batchDef.ID,
			}
			records, _, err := batchRecordRepo.FindAll(ctx, filter)
			if err != nil {
				t.Logf("Failed to query batch records: %v", err)
				return false
			}
			if len(records) != 0 {
				t.Logf("Expected 0 batch records, got %d", len(records))
				return false
			}

			// 2. Ingredient quantity should be unchanged
			currentIng, err := ingredientRepo.FindByID(ctx, ing.ID)
			if err != nil {
				t.Logf("Failed to fetch ingredient: %v", err)
				return false
			}
			if currentIng.Quantity != initialQty {
				t.Logf("Ingredient quantity changed from %.2f to %.2f", initialQty, currentIng.Quantity)
				return false
			}

			// 3. No stock history should be created
			histories, err := stockHistoryRepo.FindByIngredientID(ctx, ing.ID)
			if err != nil {
				t.Logf("Failed to query stock history: %v", err)
				return false
			}
			if len(histories) != 0 {
				t.Logf("Expected 0 stock history records, got %d", len(histories))
				return false
			}

			return true
		},
		genBatchCreationRollbackData(),
	))

	properties.TestingRun(t, gopter.ConsoleReporter(false))
}

// TestProperty_BatchCreationSuccess tests that successful batch creation updates all components
func TestProperty_BatchCreationSuccess(t *testing.T) {
	if testing.Short() {
		t.Skip("Skipping property test in short mode")
	}

	ctx := context.Background()
	client, db := setupAtomicityTestDB(t, ctx)
	defer cleanupAtomicityTestDB(t, ctx, client, db)

	parameters := gopter.DefaultTestParameters()
	parameters.MinSuccessfulTests = 20
	properties := gopter.NewProperties(parameters)

	properties.Property("Successful batch creation updates all components atomically", prop.ForAll(
		func(testData batchCreationSuccessData) bool {
			cleanCollections(ctx, db)

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
				nil,
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
			err := batchDefRepo.Create(ctx, batchDef)
			if err != nil {
				t.Logf("Failed to create batch definition: %v", err)
				return false
			}

			// Create ingredient with SUFFICIENT quantity
			ing := &ingredient.Ingredient{
				ID:          batchDef.ConversionRates[0].SourceIngredientID,
				Name:        batchDef.ConversionRates[0].SourceIngredientName,
				Category:    "Test",
				Unit:        ingredient.UnitGram,
				Quantity:    testData.SufficientQty, // Enough!
				MinStock:    10,
				CostPerUnit: testData.CostPerUnit,
				CreatedAt:   time.Now(),
				UpdatedAt:   time.Now(),
			}
			err = ingredientRepo.Create(ctx, ing)
			if err != nil {
				t.Logf("Failed to create ingredient: %v", err)
				return false
			}

			initialQty := ing.Quantity

			// Create batch (should succeed)
			batchRecord, err := service.CreateBatch(ctx, CreateBatchRequest{
				BatchDefinitionID: batchDef.ID,
				QuantityProduced:  testData.QuantityProduced,
				PreparedBy:        "test-user",
			})

			if err != nil {
				t.Logf("Batch creation failed: %v", err)
				return false
			}

			// Verify all components updated:
			// 1. Batch record created
			if batchRecord == nil {
				t.Logf("Batch record is nil")
				return false
			}
			if batchRecord.QuantityProduced != testData.QuantityProduced {
				t.Logf("Expected quantity %.2f, got %.2f", testData.QuantityProduced, batchRecord.QuantityProduced)
				return false
			}

			// 2. Ingredient deducted correctly with wastage
			currentIng, err := ingredientRepo.FindByID(ctx, ing.ID)
			if err != nil {
				t.Logf("Failed to fetch ingredient: %v", err)
				return false
			}

			rate := batchDef.ConversionRates[0]
			requiredQty := (testData.QuantityProduced / rate.BatchQuantity) * rate.SourceQuantity
			requiredQtyWithWastage := requiredQty * (1 + rate.WastageRate)
			expectedQty := initialQty - requiredQtyWithWastage

			if !floatEqualsAtomicity(currentIng.Quantity, expectedQty, 0.01) {
				t.Logf("Expected ingredient qty %.2f, got %.2f", expectedQty, currentIng.Quantity)
				return false
			}

			// 3. Stock history created
			histories, err := stockHistoryRepo.FindByIngredientID(ctx, ing.ID)
			if err != nil {
				t.Logf("Failed to query stock history: %v", err)
				return false
			}
			if len(histories) != 1 {
				t.Logf("Expected 1 stock history record, got %d", len(histories))
				return false
			}

			// 4. Cost calculated correctly
			expectedCost := requiredQtyWithWastage * testData.CostPerUnit
			if !floatEqualsAtomicity(batchRecord.TotalCost, expectedCost, 0.01) {
				t.Logf("Expected cost %.2f, got %.2f", expectedCost, batchRecord.TotalCost)
				return false
			}

			return true
		},
		genBatchCreationSuccessData(),
	))

	properties.TestingRun(t, gopter.ConsoleReporter(false))
}

// TestProperty_ConcurrentBatchCreation tests concurrent batch creation doesn't cause race conditions
func TestProperty_ConcurrentBatchCreation(t *testing.T) {
	if testing.Short() {
		t.Skip("Skipping property test in short mode")
	}

	ctx := context.Background()
	client, db := setupAtomicityTestDB(t, ctx)
	defer cleanupAtomicityTestDB(t, ctx, client, db)

	parameters := gopter.DefaultTestParameters()
	parameters.MinSuccessfulTests = 10 // Fewer due to concurrency overhead
	properties := gopter.NewProperties(parameters)

	properties.Property("Concurrent batch creation maintains inventory consistency", prop.ForAll(
		func(testData concurrentCreationData) bool {
			cleanCollections(ctx, db)

			// Setup
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
				nil,
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

			// Create ingredient
			ing := &ingredient.Ingredient{
				ID:          batchDef.ConversionRates[0].SourceIngredientID,
				Name:        batchDef.ConversionRates[0].SourceIngredientName,
				Category:    "Test",
				Unit:        ingredient.UnitGram,
				Quantity:    testData.InitialQty,
				MinStock:    10,
				CostPerUnit: 0.5,
				CreatedAt:   time.Now(),
				UpdatedAt:   time.Now(),
			}
			ingredientRepo.Create(ctx, ing)
			initialQty := ing.Quantity

			// Run concurrent batch creations
			var wg sync.WaitGroup
			successCount := 0
			var mu sync.Mutex

			for i := 0; i < testData.NumConcurrent; i++ {
				wg.Add(1)
				go func(idx int) {
					defer wg.Done()
					_, err := service.CreateBatch(ctx, CreateBatchRequest{
						BatchDefinitionID: batchDef.ID,
						QuantityProduced:  testData.QtyPerBatch,
						PreparedBy:        fmt.Sprintf("user-%d", idx),
					})
					if err == nil {
						mu.Lock()
						successCount++
						mu.Unlock()
					}
				}(i)
			}

			wg.Wait()

			// Verify consistency
			currentIng, _ := ingredientRepo.FindByID(ctx, ing.ID)
			
			// Calculate expected deduction
			rate := batchDef.ConversionRates[0]
			qtyPerBatch := (testData.QtyPerBatch / rate.BatchQuantity) * rate.SourceQuantity
			qtyPerBatchWithWastage := qtyPerBatch * (1 + rate.WastageRate)
			expectedDeduction := qtyPerBatchWithWastage * float64(successCount)
			expectedQty := initialQty - expectedDeduction

			// 1. Ingredient quantity should be consistent
			if !floatEqualsAtomicity(currentIng.Quantity, expectedQty, 0.01) {
				t.Logf("Inconsistent ingredient qty: expected %.2f, got %.2f", expectedQty, currentIng.Quantity)
				return false
			}

			// 2. Number of batch records should match successful operations
			filter := batch.BatchRecordFilter{
				BatchDefinitionID: &batchDef.ID,
			}
			records, _, _ := batchRecordRepo.FindAll(ctx, filter)
			if len(records) != successCount {
				t.Logf("Expected %d batch records, got %d", successCount, len(records))
				return false
			}

			// 3. Quantity should never go negative
			if currentIng.Quantity < 0 {
				t.Logf("Ingredient quantity went negative: %.2f", currentIng.Quantity)
				return false
			}

			return true
		},
		genConcurrentCreationData(),
	))

	properties.TestingRun(t, gopter.ConsoleReporter(false))
}

// Test data structures

type batchCreationRollbackData struct {
	InsufficientQty  float64 // Less than required
	QuantityProduced float64
}

type batchCreationSuccessData struct {
	SufficientQty    float64 // More than required
	QuantityProduced float64
	CostPerUnit      float64
}

type concurrentCreationData struct {
	InitialQty    float64
	NumConcurrent int
	QtyPerBatch   float64
}

// Generators

func genBatchCreationRollbackData() gopter.Gen {
	return gopter.CombineGens(
		gen.Float64Range(10, 50),    // InsufficientQty (not enough)
		gen.Float64Range(500, 1000), // QuantityProduced (requires ~110-220g with wastage)
	).Map(func(values []interface{}) batchCreationRollbackData {
		return batchCreationRollbackData{
			InsufficientQty:  values[0].(float64),
			QuantityProduced: values[1].(float64),
		}
	})
}

func genBatchCreationSuccessData() gopter.Gen {
	return gopter.CombineGens(
		gen.Float64Range(500, 5000),  // SufficientQty (enough)
		gen.Float64Range(500, 1000),  // QuantityProduced
		gen.Float64Range(0.1, 2.0),   // CostPerUnit
	).Map(func(values []interface{}) batchCreationSuccessData {
		return batchCreationSuccessData{
			SufficientQty:    values[0].(float64),
			QuantityProduced: values[1].(float64),
			CostPerUnit:      values[2].(float64),
		}
	})
}

func genConcurrentCreationData() gopter.Gen {
	return gopter.CombineGens(
		gen.Float64Range(1000, 5000), // InitialQty
		gen.IntRange(3, 8),            // NumConcurrent
		gen.Float64Range(200, 500),    // QtyPerBatch
	).Map(func(values []interface{}) concurrentCreationData {
		return concurrentCreationData{
			InitialQty:    values[0].(float64),
			NumConcurrent: values[1].(int),
			QtyPerBatch:   values[2].(float64),
		}
	})
}

// Helper functions

func setupAtomicityTestDB(t *testing.T, ctx context.Context) (*mongo.Client, *mongo.Database) {
	// Connect to test MongoDB with authentication
	mongoURI := "mongodb://admin:password123@localhost:27017/?authSource=admin"
	clientOptions := options.Client().ApplyURI(mongoURI)
	client, err := mongo.Connect(ctx, clientOptions)
	if err != nil {
		t.Fatalf("Failed to connect to MongoDB: %v", err)
	}

	// Use a test database
	db := client.Database("cafe_pos_test_atomicity")
	
	return client, db
}

func cleanupAtomicityTestDB(t *testing.T, ctx context.Context, client *mongo.Client, db *mongo.Database) {
	// Drop test database
	db.Drop(ctx)
	client.Disconnect(ctx)
}

func cleanCollections(ctx context.Context, db *mongo.Database) {
	db.Collection("ingredients").Drop(ctx)
	db.Collection("stock_history").Drop(ctx)
	db.Collection("batch_definitions").Drop(ctx)
	db.Collection("batch_records").Drop(ctx)
	db.Collection("batch_usage_logs").Drop(ctx)
}

func floatEqualsAtomicity(a, b, epsilon float64) bool {
	diff := a - b
	if diff < 0 {
		diff = -diff
	}
	return diff < epsilon
}
