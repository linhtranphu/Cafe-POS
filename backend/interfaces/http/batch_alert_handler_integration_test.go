package http

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"cafe-pos/backend/application/services"
	"cafe-pos/backend/domain/batch"
	"cafe-pos/backend/infrastructure/mongodb"
	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

// TestBatchAlertHandler_Integration tests the complete batch alert flow
// Task 4.4.3: Integration tests for BatchAlertHandler
// Requirements: 4.1, 4.2, 4.4, 4.5, 4.6
//
// This integration test verifies:
// 1. Create batch definitions and records via repository
// 2. Get alerts via GET /api/batch-alerts
// 3. Verify low stock alerts are returned correctly
// 4. Verify expiring alerts are returned correctly
// 5. Verify expired alerts are returned correctly
func TestBatchAlertHandler_Integration(t *testing.T) {
	// Skip if not running integration tests
	if testing.Short() {
		t.Skip("Skipping integration test")
	}

	ctx := context.Background()

	// Setup MongoDB connection
	client, err := mongo.Connect(ctx, options.Client().ApplyURI("mongodb://localhost:27017"))
	require.NoError(t, err)
	defer client.Disconnect(ctx)

	// Use test database
	db := client.Database("cafe_pos_test_batch_alerts_" + primitive.NewObjectID().Hex())
	defer db.Drop(ctx)

	// Setup repositories
	batchDefinitionRepo := mongodb.NewBatchDefinitionRepository(db)
	batchRecordRepo := mongodb.NewBatchRecordRepository(db)

	// Setup service and handler
	batchAlertService := services.NewBatchAlertService(batchDefinitionRepo, batchRecordRepo)
	handler := NewBatchAlertHandler(batchAlertService)

	// Setup Gin router
	gin.SetMode(gin.TestMode)
	router := gin.New()
	router.GET("/api/batch-alerts", handler.GetAlerts)

	// ========================================
	// STEP 1: Create test batch definitions
	// ========================================

	// Definition 1: Low stock threshold = 50ml
	def1 := &batch.BatchDefinition{
		ID:                  primitive.NewObjectID(),
		Name:                "Coffee Concentrate",
		Unit:                "ml",
		ShelfLifeHours:      24,
		ConversionRates:     []batch.ConversionRate{},
		LowStockThreshold:   50.0,
		ExpiryWarningHours:  4,
		CreatedAt:           time.Now(),
		UpdatedAt:           time.Now(),
	}
	err = batchDefinitionRepo.Create(ctx, def1)
	require.NoError(t, err)

	// Definition 2: Low stock threshold = 100ml
	def2 := &batch.BatchDefinition{
		ID:                  primitive.NewObjectID(),
		Name:                "Tea Concentrate",
		Unit:                "ml",
		ShelfLifeHours:      12,
		ConversionRates:     []batch.ConversionRate{},
		LowStockThreshold:   100.0,
		ExpiryWarningHours:  2,
		CreatedAt:           time.Now(),
		UpdatedAt:           time.Now(),
	}
	err = batchDefinitionRepo.Create(ctx, def2)
	require.NoError(t, err)

	// ========================================
	// STEP 2: Create test batch records
	// ========================================

	// Batch 1: Low stock (30ml < 50ml threshold)
	batch1 := &batch.BatchRecord{
		ID:                primitive.NewObjectID(),
		BatchDefinitionID: def1.ID,
		BatchName:         def1.Name,
		QuantityProduced:  100.0,
		QuantityRemaining: 30.0, // Below threshold
		Unit:              def1.Unit,
		CostPerUnit:       0.15,
		TotalCost:         15.0,
		PreparedBy:        "test_user",
		PreparedAt:        time.Now().Add(-20 * time.Hour),
		ExpiresAt:         time.Now().Add(4 * time.Hour), // Not expiring yet
		Status:            batch.BatchStatusAvailable,
		IngredientsUsed:   []batch.IngredientUsage{},
		CreatedAt:         time.Now(),
		UpdatedAt:         time.Now(),
	}
	err = batchRecordRepo.Create(ctx, batch1)
	require.NoError(t, err)

	// Batch 2: Expiring soon (2 hours until expiry, warning threshold = 4 hours)
	batch2 := &batch.BatchRecord{
		ID:                primitive.NewObjectID(),
		BatchDefinitionID: def1.ID,
		BatchName:         def1.Name,
		QuantityProduced:  100.0,
		QuantityRemaining: 80.0,
		Unit:              def1.Unit,
		CostPerUnit:       0.15,
		TotalCost:         15.0,
		PreparedBy:        "test_user",
		PreparedAt:        time.Now().Add(-22 * time.Hour),
		ExpiresAt:         time.Now().Add(2 * time.Hour), // Expiring in 2 hours
		Status:            batch.BatchStatusAvailable,
		IngredientsUsed:   []batch.IngredientUsage{},
		CreatedAt:         time.Now(),
		UpdatedAt:         time.Now(),
	}
	err = batchRecordRepo.Create(ctx, batch2)
	require.NoError(t, err)

	// Batch 3: Already expired with remaining quantity
	batch3 := &batch.BatchRecord{
		ID:                primitive.NewObjectID(),
		BatchDefinitionID: def2.ID,
		BatchName:         def2.Name,
		QuantityProduced:  100.0,
		QuantityRemaining: 50.0, // Wasted quantity
		Unit:              def2.Unit,
		CostPerUnit:       0.20,
		TotalCost:         20.0,
		PreparedBy:        "test_user",
		PreparedAt:        time.Now().Add(-15 * time.Hour),
		ExpiresAt:         time.Now().Add(-1 * time.Hour), // Expired 1 hour ago
		Status:            batch.BatchStatusExpired,
		IngredientsUsed:   []batch.IngredientUsage{},
		CreatedAt:         time.Now(),
		UpdatedAt:         time.Now(),
	}
	err = batchRecordRepo.Create(ctx, batch3)
	require.NoError(t, err)

	// Batch 4: Normal batch (not low stock, not expiring)
	batch4 := &batch.BatchRecord{
		ID:                primitive.NewObjectID(),
		BatchDefinitionID: def2.ID,
		BatchName:         def2.Name,
		QuantityProduced:  200.0,
		QuantityRemaining: 150.0, // Above threshold
		Unit:              def2.Unit,
		CostPerUnit:       0.20,
		TotalCost:         40.0,
		PreparedBy:        "test_user",
		PreparedAt:        time.Now().Add(-2 * time.Hour),
		ExpiresAt:         time.Now().Add(10 * time.Hour), // Not expiring soon
		Status:            batch.BatchStatusAvailable,
		IngredientsUsed:   []batch.IngredientUsage{},
		CreatedAt:         time.Now(),
		UpdatedAt:         time.Now(),
	}
	err = batchRecordRepo.Create(ctx, batch4)
	require.NoError(t, err)

	// ========================================
	// STEP 3: Test GET /api/batch-alerts
	// ========================================

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/api/batch-alerts", nil)
	router.ServeHTTP(w, req)

	// Verify response
	assert.Equal(t, http.StatusOK, w.Code)

	var alerts batch.BatchAlerts
	err = json.Unmarshal(w.Body.Bytes(), &alerts)
	require.NoError(t, err)

	// ========================================
	// STEP 4: Verify low stock alerts
	// ========================================

	// Should have 1 low stock alert (Coffee Concentrate: 30ml < 50ml)
	assert.Len(t, alerts.LowStock, 1, "Should have 1 low stock alert")
	
	if len(alerts.LowStock) > 0 {
		lowStockAlert := alerts.LowStock[0]
		assert.Equal(t, def1.ID, lowStockAlert.BatchDefinitionID)
		assert.Equal(t, "Coffee Concentrate", lowStockAlert.BatchName)
		assert.Equal(t, 30.0, lowStockAlert.CurrentStock)
		assert.Equal(t, 50.0, lowStockAlert.Threshold)
		assert.Equal(t, "ml", lowStockAlert.Unit)
	}

	// ========================================
	// STEP 5: Verify expiring alerts
	// ========================================

	// Should have 1 expiring alert (batch2: expires in 2 hours, warning = 4 hours)
	assert.Len(t, alerts.Expiring, 1, "Should have 1 expiring alert")
	
	if len(alerts.Expiring) > 0 {
		expiringAlert := alerts.Expiring[0]
		assert.Equal(t, batch2.ID, expiringAlert.BatchRecordID)
		assert.Equal(t, "Coffee Concentrate", expiringAlert.BatchName)
		assert.Equal(t, 80.0, expiringAlert.QuantityRemaining)
		assert.Equal(t, "ml", expiringAlert.Unit)
		// Hours until expiry should be around 2 (allow some tolerance)
		assert.GreaterOrEqual(t, expiringAlert.HoursUntilExpiry, 1)
		assert.LessOrEqual(t, expiringAlert.HoursUntilExpiry, 3)
	}

	// ========================================
	// STEP 6: Verify expired alerts
	// ========================================

	// Should have 1 expired alert (batch3: expired with 50ml wasted)
	assert.Len(t, alerts.Expired, 1, "Should have 1 expired alert")
	
	if len(alerts.Expired) > 0 {
		expiredAlert := alerts.Expired[0]
		assert.Equal(t, batch3.ID, expiredAlert.BatchRecordID)
		assert.Equal(t, "Tea Concentrate", expiredAlert.BatchName)
		assert.Equal(t, 50.0, expiredAlert.QuantityWasted)
		assert.Equal(t, "ml", expiredAlert.Unit)
		// Cost wasted = 50ml * 0.20 = 10.0
		assert.Equal(t, 10.0, expiredAlert.CostWasted)
	}

	// ========================================
	// STEP 7: Verify last checked timestamp
	// ========================================

	assert.False(t, alerts.LastChecked.IsZero(), "LastChecked should be set")
	assert.WithinDuration(t, time.Now(), alerts.LastChecked, 5*time.Second)
}

// TestBatchAlertHandler_NoAlerts tests when there are no alerts
func TestBatchAlertHandler_NoAlerts(t *testing.T) {
	// Skip if not running integration tests
	if testing.Short() {
		t.Skip("Skipping integration test")
	}

	ctx := context.Background()

	// Setup MongoDB connection
	client, err := mongo.Connect(ctx, options.Client().ApplyURI("mongodb://localhost:27017"))
	require.NoError(t, err)
	defer client.Disconnect(ctx)

	// Use test database
	db := client.Database("cafe_pos_test_batch_alerts_no_alerts_" + primitive.NewObjectID().Hex())
	defer db.Drop(ctx)

	// Setup repositories
	batchDefinitionRepo := mongodb.NewBatchDefinitionRepository(db)
	batchRecordRepo := mongodb.NewBatchRecordRepository(db)

	// Setup service and handler
	batchAlertService := services.NewBatchAlertService(batchDefinitionRepo, batchRecordRepo)
	handler := NewBatchAlertHandler(batchAlertService)

	// Setup Gin router
	gin.SetMode(gin.TestMode)
	router := gin.New()
	router.GET("/api/batch-alerts", handler.GetAlerts)

	// Create a batch definition with high threshold
	def := &batch.BatchDefinition{
		ID:                  primitive.NewObjectID(),
		Name:                "Test Batch",
		Unit:                "ml",
		ShelfLifeHours:      24,
		ConversionRates:     []batch.ConversionRate{},
		LowStockThreshold:   10.0, // Very low threshold
		ExpiryWarningHours:  1,    // Very short warning
		CreatedAt:           time.Now(),
		UpdatedAt:           time.Now(),
	}
	err = batchDefinitionRepo.Create(ctx, def)
	require.NoError(t, err)

	// Create a normal batch (not low stock, not expiring)
	batchRecord := &batch.BatchRecord{
		ID:                primitive.NewObjectID(),
		BatchDefinitionID: def.ID,
		BatchName:         def.Name,
		QuantityProduced:  100.0,
		QuantityRemaining: 100.0, // Well above threshold
		Unit:              def.Unit,
		CostPerUnit:       0.15,
		TotalCost:         15.0,
		PreparedBy:        "test_user",
		PreparedAt:        time.Now(),
		ExpiresAt:         time.Now().Add(24 * time.Hour), // Not expiring soon
		Status:            batch.BatchStatusAvailable,
		IngredientsUsed:   []batch.IngredientUsage{},
		CreatedAt:         time.Now(),
		UpdatedAt:         time.Now(),
	}
	err = batchRecordRepo.Create(ctx, batchRecord)
	require.NoError(t, err)

	// Test GET /api/batch-alerts
	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/api/batch-alerts", nil)
	router.ServeHTTP(w, req)

	// Verify response
	assert.Equal(t, http.StatusOK, w.Code)

	var alerts batch.BatchAlerts
	err = json.Unmarshal(w.Body.Bytes(), &alerts)
	require.NoError(t, err)

	// Should have no alerts
	assert.Empty(t, alerts.LowStock, "Should have no low stock alerts")
	assert.Empty(t, alerts.Expiring, "Should have no expiring alerts")
	assert.Empty(t, alerts.Expired, "Should have no expired alerts")
	assert.False(t, alerts.LastChecked.IsZero(), "LastChecked should be set")
}

// TestBatchAlertHandler_EmptyDatabase tests when database is empty
func TestBatchAlertHandler_EmptyDatabase(t *testing.T) {
	// Skip if not running integration tests
	if testing.Short() {
		t.Skip("Skipping integration test")
	}

	ctx := context.Background()

	// Setup MongoDB connection
	client, err := mongo.Connect(ctx, options.Client().ApplyURI("mongodb://localhost:27017"))
	require.NoError(t, err)
	defer client.Disconnect(ctx)

	// Use test database
	db := client.Database("cafe_pos_test_batch_alerts_empty_" + primitive.NewObjectID().Hex())
	defer db.Drop(ctx)

	// Setup repositories
	batchDefinitionRepo := mongodb.NewBatchDefinitionRepository(db)
	batchRecordRepo := mongodb.NewBatchRecordRepository(db)

	// Setup service and handler
	batchAlertService := services.NewBatchAlertService(batchDefinitionRepo, batchRecordRepo)
	handler := NewBatchAlertHandler(batchAlertService)

	// Setup Gin router
	gin.SetMode(gin.TestMode)
	router := gin.New()
	router.GET("/api/batch-alerts", handler.GetAlerts)

	// Test GET /api/batch-alerts with empty database
	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/api/batch-alerts", nil)
	router.ServeHTTP(w, req)

	// Verify response
	assert.Equal(t, http.StatusOK, w.Code)

	var alerts batch.BatchAlerts
	err = json.Unmarshal(w.Body.Bytes(), &alerts)
	require.NoError(t, err)

	// Should have no alerts
	assert.Empty(t, alerts.LowStock, "Should have no low stock alerts")
	assert.Empty(t, alerts.Expiring, "Should have no expiring alerts")
	assert.Empty(t, alerts.Expired, "Should have no expired alerts")
	assert.False(t, alerts.LastChecked.IsZero(), "LastChecked should be set")
}
