package http

import (
	"bytes"
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

// TestBatchUsageHandler_Integration tests the complete batch usage flow
// Task 4.3.4: Integration tests for BatchUsageHandler
// Requirements: 5.2, 5.3, 5.4, 5.5, 5.6
//
// This integration test verifies:
// 1. Create batch records via repository
// 2. Use batch via POST /api/batch-usage
// 3. Verify batch quantity is deducted (FIFO)
// 4. Verify usage is logged
// 5. Get usage history via GET /api/batch-usage/history
func TestBatchUsageHandler_Integration(t *testing.T) {
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
	db := client.Database("cafe_pos_test_batch_usage_" + primitive.NewObjectID().Hex())
	defer db.Drop(ctx)

	// Setup repositories
	batchRecordRepo := mongodb.NewBatchRecordRepository(db)
	batchUsageLogRepo := mongodb.NewBatchUsageLogRepository(db)

	// Setup service and handler
	batchUsageService := services.NewBatchUsageService(batchRecordRepo, batchUsageLogRepo)
	handler := NewBatchUsageHandler(batchUsageService)

	// Setup Gin router
	gin.SetMode(gin.TestMode)
	router := gin.New()
	router.POST("/api/batch-usage", handler.UseBatch)
	router.GET("/api/batch-usage/history", handler.GetUsageHistory)

	// ========================================
	// STEP 1: Create test batch records
	// ========================================
	batchDefID := primitive.NewObjectID()
	orderID := primitive.NewObjectID()
	menuItemID := primitive.NewObjectID()

	// Create first batch (older, should be used first - FIFO)
	batch1 := &batch.BatchRecord{
		ID:                primitive.NewObjectID(),
		BatchDefinitionID: batchDefID,
		BatchName:         "Test Batch",
		QuantityProduced:  100.0,
		QuantityRemaining: 100.0,
		Unit:              "ml",
		CostPerUnit:       0.15,
		TotalCost:         15.0,
		PreparedBy:        "test_user",
		PreparedAt:        time.Now().Add(-2 * time.Hour),
		ExpiresAt:         time.Now().Add(22 * time.Hour), // Expires in 22 hours
		Status:            batch.BatchStatusAvailable,
		IngredientsUsed:   []batch.IngredientUsage{},
		CreatedAt:         time.Now(),
		UpdatedAt:         time.Now(),
	}
	err = batchRecordRepo.Create(ctx, batch1)
	require.NoError(t, err)

	// Create second batch (newer, should be used second - FIFO)
	batch2 := &batch.BatchRecord{
		ID:                primitive.NewObjectID(),
		BatchDefinitionID: batchDefID,
		BatchName:         "Test Batch",
		QuantityProduced:  100.0,
		QuantityRemaining: 100.0,
		Unit:              "ml",
		CostPerUnit:       0.15,
		TotalCost:         15.0,
		PreparedBy:        "test_user",
		PreparedAt:        time.Now().Add(-1 * time.Hour),
		ExpiresAt:         time.Now().Add(23 * time.Hour), // Expires in 23 hours
		Status:            batch.BatchStatusAvailable,
		IngredientsUsed:   []batch.IngredientUsage{},
		CreatedAt:         time.Now(),
		UpdatedAt:         time.Now(),
	}
	err = batchRecordRepo.Create(ctx, batch2)
	require.NoError(t, err)

	t.Log("✓ Step 1: Created 2 batch records")

	// ========================================
	// STEP 2: Use batch via API (should use batch1 first - FIFO)
	// ========================================
	usageReq := UseBatchRequest{
		BatchDefinitionID: batchDefID.Hex(),
		QuantityNeeded:    30.0,
		OrderID:           orderID.Hex(),
		MenuItemID:        menuItemID.Hex(),
		MenuItemName:      "Test Menu Item",
	}
	reqBody, _ := json.Marshal(usageReq)

	w := httptest.NewRecorder()
	req := httptest.NewRequest("POST", "/api/batch-usage", bytes.NewBuffer(reqBody))
	req.Header.Set("Content-Type", "application/json")
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)

	var usageResp UseBatchResponse
	err = json.Unmarshal(w.Body.Bytes(), &usageResp)
	require.NoError(t, err)
	assert.True(t, usageResp.Success)
	assert.Equal(t, 1, len(usageResp.BatchesUsed))
	assert.Equal(t, batch1.ID.Hex(), usageResp.BatchesUsed[0].BatchRecordID)
	assert.Equal(t, 30.0, usageResp.BatchesUsed[0].QuantityUsed)
	assert.Equal(t, 4.5, usageResp.TotalCost) // 30 * 0.15

	t.Log("✓ Step 2: Used 30ml from batch1 (FIFO)")

	// ========================================
	// STEP 3: Verify batch1 quantity was deducted
	// ========================================
	updatedBatch1, err := batchRecordRepo.FindByID(ctx, batch1.ID)
	require.NoError(t, err)
	assert.Equal(t, 70.0, updatedBatch1.QuantityRemaining)
	assert.Equal(t, batch.BatchStatusAvailable, updatedBatch1.Status)

	t.Log("✓ Step 3: Batch1 quantity deducted correctly (100 -> 70)")

	// ========================================
	// STEP 4: Verify usage was logged
	// ========================================
	usageLogs, err := batchUsageLogRepo.FindAll(ctx, batch.BatchUsageLogFilter{
		BatchRecordID: &batch1.ID,
	})
	require.NoError(t, err)
	assert.Equal(t, 1, len(usageLogs))
	assert.Equal(t, batch1.ID, usageLogs[0].BatchRecordID)
	assert.Equal(t, orderID, usageLogs[0].OrderID)
	assert.Equal(t, menuItemID, usageLogs[0].MenuItemID)
	assert.Equal(t, 30.0, usageLogs[0].QuantityUsed)
	assert.Equal(t, 0.15, usageLogs[0].CostPerUnit)
	assert.Equal(t, 4.5, usageLogs[0].TotalCost)

	t.Log("✓ Step 4: Usage logged correctly")

	// ========================================
	// STEP 5: Get usage history via API
	// ========================================
	w = httptest.NewRecorder()
	req = httptest.NewRequest("GET", "/api/batch-usage/history?batch_record_id="+batch1.ID.Hex(), nil)
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)

	var historyResp map[string]interface{}
	err = json.Unmarshal(w.Body.Bytes(), &historyResp)
	require.NoError(t, err)
	assert.Equal(t, float64(1), historyResp["total"])

	data := historyResp["data"].([]interface{})
	assert.Equal(t, 1, len(data))

	t.Log("✓ Step 5: Retrieved usage history via API")

	// ========================================
	// STEP 6: Use more batch (should still use batch1 until depleted)
	// ========================================
	usageReq.QuantityNeeded = 80.0 // More than batch1 remaining (70), should use both batches
	reqBody, _ = json.Marshal(usageReq)

	w = httptest.NewRecorder()
	req = httptest.NewRequest("POST", "/api/batch-usage", bytes.NewBuffer(reqBody))
	req.Header.Set("Content-Type", "application/json")
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)

	err = json.Unmarshal(w.Body.Bytes(), &usageResp)
	require.NoError(t, err)
	assert.True(t, usageResp.Success)
	assert.Equal(t, 2, len(usageResp.BatchesUsed)) // Should use both batches
	
	// First batch should be fully depleted (70ml)
	assert.Equal(t, batch1.ID.Hex(), usageResp.BatchesUsed[0].BatchRecordID)
	assert.Equal(t, 70.0, usageResp.BatchesUsed[0].QuantityUsed)
	
	// Second batch should provide remaining (10ml)
	assert.Equal(t, batch2.ID.Hex(), usageResp.BatchesUsed[1].BatchRecordID)
	assert.Equal(t, 10.0, usageResp.BatchesUsed[1].QuantityUsed)
	
	// Total cost: (70 * 0.15) + (10 * 0.15) = 12.0
	assert.Equal(t, 12.0, usageResp.TotalCost)

	t.Log("✓ Step 6: Used 80ml across both batches (FIFO: 70ml from batch1, 10ml from batch2)")

	// ========================================
	// STEP 7: Verify batch1 is depleted
	// ========================================
	updatedBatch1, err = batchRecordRepo.FindByID(ctx, batch1.ID)
	require.NoError(t, err)
	assert.Equal(t, 0.0, updatedBatch1.QuantityRemaining)
	assert.Equal(t, batch.BatchStatusDepleted, updatedBatch1.Status)

	t.Log("✓ Step 7: Batch1 is now depleted")

	// ========================================
	// STEP 8: Verify batch2 quantity was deducted
	// ========================================
	updatedBatch2, err := batchRecordRepo.FindByID(ctx, batch2.ID)
	require.NoError(t, err)
	assert.Equal(t, 90.0, updatedBatch2.QuantityRemaining)
	assert.Equal(t, batch.BatchStatusAvailable, updatedBatch2.Status)

	t.Log("✓ Step 8: Batch2 quantity deducted correctly (100 -> 90)")

	// ========================================
	// STEP 9: Test insufficient batch scenario
	// ========================================
	usageReq.QuantityNeeded = 100.0 // More than available (90ml in batch2)
	reqBody, _ = json.Marshal(usageReq)

	w = httptest.NewRecorder()
	req = httptest.NewRequest("POST", "/api/batch-usage", bytes.NewBuffer(reqBody))
	req.Header.Set("Content-Type", "application/json")
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)

	err = json.Unmarshal(w.Body.Bytes(), &usageResp)
	require.NoError(t, err)
	assert.False(t, usageResp.Success)
	assert.Contains(t, usageResp.Message, "Insufficient")

	t.Log("✓ Step 9: Correctly rejected request for insufficient batch")

	// ========================================
	// FINAL VERIFICATION
	// ========================================
	t.Log("\n=== Integration Test Summary ===")
	t.Log("✓ Step 1: Created 2 batch records")
	t.Log("✓ Step 2: Used batch via POST /api/batch-usage")
	t.Log("✓ Step 3: Verified FIFO ordering (older batch used first)")
	t.Log("✓ Step 4: Verified usage logging")
	t.Log("✓ Step 5: Retrieved usage history via GET /api/batch-usage/history")
	t.Log("✓ Step 6: Verified multi-batch usage when single batch insufficient")
	t.Log("✓ Step 7: Verified batch status updates (depleted)")
	t.Log("✓ Step 8: Verified quantity deduction across batches")
	t.Log("✓ Step 9: Verified insufficient batch error handling")
	t.Log("\n✅ All integration tests passed!")
}
