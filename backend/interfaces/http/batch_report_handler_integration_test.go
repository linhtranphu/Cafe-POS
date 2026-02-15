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

// TestBatchReportHandler_ProductionReport_Integration tests the production report endpoint
func TestBatchReportHandler_ProductionReport_Integration(t *testing.T) {
	// Skip if not running integration tests
	if testing.Short() {
		t.Skip("Skipping integration test")
	}

	// Setup MongoDB connection
	ctx := context.Background()
	client, err := mongo.Connect(ctx, options.Client().ApplyURI("mongodb://localhost:27017"))
	require.NoError(t, err)
	defer client.Disconnect(ctx)

	// Use test database
	db := client.Database("cafe_pos_test_batch_report_production")
	defer db.Drop(ctx)

	// Create repositories
	batchRecordRepo := mongodb.NewBatchRecordRepository(db)
	batchUsageLogRepo := mongodb.NewBatchUsageLogRepository(db)
	batchDefRepo := mongodb.NewBatchDefinitionRepository(db)

	// Create service and handler
	reportService := services.NewBatchReportService(batchRecordRepo, batchUsageLogRepo, batchDefRepo)
	handler := NewBatchReportHandler(reportService)

	// Setup Gin router
	gin.SetMode(gin.TestMode)
	router := gin.New()
	router.GET("/api/batch-reports/production", handler.GetProductionReport)

	// ========================================
	// STEP 1: Create test data
	// ========================================

	// Create batch definition
	batchDef := &batch.BatchDefinition{
		ID:   primitive.NewObjectID(),
		Name: "Coffee Concentrate",
	}
	err = batchDefRepo.Create(ctx, batchDef)
	require.NoError(t, err)

	// Create batch records
	now := time.Now()
	fromDate := now.Add(-24 * time.Hour)
	toDate := now.Add(24 * time.Hour)

	record1 := &batch.BatchRecord{
		ID:                primitive.NewObjectID(),
		BatchDefinitionID: batchDef.ID,
		BatchName:         "Coffee Concentrate",
		QuantityProduced:  500,
		QuantityRemaining: 500,
		Unit:              "ml",
		CostPerUnit:       0.15,
		TotalCost:         75.0,
		PreparedBy:        "user1",
		PreparedAt:        now.Add(-12 * time.Hour),
		ExpiresAt:         now.Add(12 * time.Hour),
		Status:            batch.BatchStatusAvailable,
	}
	err = batchRecordRepo.Create(ctx, record1)
	require.NoError(t, err)

	record2 := &batch.BatchRecord{
		ID:                primitive.NewObjectID(),
		BatchDefinitionID: batchDef.ID,
		BatchName:         "Coffee Concentrate",
		QuantityProduced:  300,
		QuantityRemaining: 300,
		Unit:              "ml",
		CostPerUnit:       0.15,
		TotalCost:         45.0,
		PreparedBy:        "user2",
		PreparedAt:        now.Add(-6 * time.Hour),
		ExpiresAt:         now.Add(18 * time.Hour),
		Status:            batch.BatchStatusAvailable,
	}
	err = batchRecordRepo.Create(ctx, record2)
	require.NoError(t, err)

	// ========================================
	// STEP 2: Test GET /api/batch-reports/production
	// ========================================

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/api/batch-reports/production?from_date="+fromDate.Format(time.RFC3339)+"&to_date="+toDate.Format(time.RFC3339), nil)
	router.ServeHTTP(w, req)

	// Verify response
	assert.Equal(t, http.StatusOK, w.Code)

	var report services.ProductionReport
	err = json.Unmarshal(w.Body.Bytes(), &report)
	require.NoError(t, err)

	// Verify report data
	assert.Equal(t, 2, report.TotalBatchesProduced)
	assert.Equal(t, 800.0, report.TotalQuantityProduced)
	assert.Equal(t, 120.0, report.TotalCost)
	assert.Len(t, report.ByBatchType, 1)
	assert.Equal(t, "Coffee Concentrate", report.ByBatchType[0].BatchName)
	assert.Equal(t, 2, report.ByBatchType[0].Count)
	assert.Len(t, report.ByPreparer, 2)

	// ========================================
	// STEP 3: Test with batch_definition_id filter
	// ========================================

	w = httptest.NewRecorder()
	req, _ = http.NewRequest("GET", "/api/batch-reports/production?from_date="+fromDate.Format(time.RFC3339)+"&to_date="+toDate.Format(time.RFC3339)+"&batch_definition_id="+batchDef.ID.Hex(), nil)
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)

	err = json.Unmarshal(w.Body.Bytes(), &report)
	require.NoError(t, err)
	assert.Equal(t, 2, report.TotalBatchesProduced)

	// ========================================
	// STEP 4: Test with prepared_by filter
	// ========================================

	w = httptest.NewRecorder()
	req, _ = http.NewRequest("GET", "/api/batch-reports/production?from_date="+fromDate.Format(time.RFC3339)+"&to_date="+toDate.Format(time.RFC3339)+"&prepared_by=user1", nil)
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)

	err = json.Unmarshal(w.Body.Bytes(), &report)
	require.NoError(t, err)
	assert.Equal(t, 1, report.TotalBatchesProduced)
	assert.Equal(t, 500.0, report.TotalQuantityProduced)

	t.Log("\n=== Integration Test Summary ===")
	t.Log("✓ Step 1: Created test batch records")
	t.Log("✓ Step 2: Retrieved production report")
	t.Log("✓ Step 3: Tested batch_definition_id filter")
	t.Log("✓ Step 4: Tested prepared_by filter")
}

// TestBatchReportHandler_WastageReport_Integration tests the wastage report endpoint
func TestBatchReportHandler_WastageReport_Integration(t *testing.T) {
	// Skip if not running integration tests
	if testing.Short() {
		t.Skip("Skipping integration test")
	}

	// Setup MongoDB connection
	ctx := context.Background()
	client, err := mongo.Connect(ctx, options.Client().ApplyURI("mongodb://localhost:27017"))
	require.NoError(t, err)
	defer client.Disconnect(ctx)

	// Use test database
	db := client.Database("cafe_pos_test_batch_report_wastage")
	defer db.Drop(ctx)

	// Create repositories
	batchRecordRepo := mongodb.NewBatchRecordRepository(db)
	batchUsageLogRepo := mongodb.NewBatchUsageLogRepository(db)
	batchDefRepo := mongodb.NewBatchDefinitionRepository(db)

	// Create service and handler
	reportService := services.NewBatchReportService(batchRecordRepo, batchUsageLogRepo, batchDefRepo)
	handler := NewBatchReportHandler(reportService)

	// Setup Gin router
	gin.SetMode(gin.TestMode)
	router := gin.New()
	router.GET("/api/batch-reports/wastage", handler.GetWastageReport)

	// ========================================
	// STEP 1: Create test data with expired batches
	// ========================================

	now := time.Now()
	fromDate := now.Add(-24 * time.Hour)
	toDate := now.Add(24 * time.Hour)

	// Create expired batch records
	expiredRecord1 := &batch.BatchRecord{
		ID:                primitive.NewObjectID(),
		BatchDefinitionID: primitive.NewObjectID(),
		BatchName:         "Coffee Concentrate",
		QuantityProduced:  500,
		QuantityRemaining: 200, // 200ml wasted
		Unit:              "ml",
		CostPerUnit:       0.15,
		TotalCost:         75.0,
		PreparedBy:        "user1",
		PreparedAt:        now.Add(-48 * time.Hour),
		ExpiresAt:         now.Add(-12 * time.Hour),
		Status:            batch.BatchStatusExpired,
	}
	err = batchRecordRepo.Create(ctx, expiredRecord1)
	require.NoError(t, err)

	expiredRecord2 := &batch.BatchRecord{
		ID:                primitive.NewObjectID(),
		BatchDefinitionID: primitive.NewObjectID(),
		BatchName:         "Tea Concentrate",
		QuantityProduced:  300,
		QuantityRemaining: 100, // 100ml wasted
		Unit:              "ml",
		CostPerUnit:       0.10,
		TotalCost:         30.0,
		PreparedBy:        "user2",
		PreparedAt:        now.Add(-36 * time.Hour),
		ExpiresAt:         now.Add(-6 * time.Hour),
		Status:            batch.BatchStatusExpired,
	}
	err = batchRecordRepo.Create(ctx, expiredRecord2)
	require.NoError(t, err)

	// ========================================
	// STEP 2: Test GET /api/batch-reports/wastage
	// ========================================

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/api/batch-reports/wastage?from_date="+fromDate.Format(time.RFC3339)+"&to_date="+toDate.Format(time.RFC3339), nil)
	router.ServeHTTP(w, req)

	// Verify response
	assert.Equal(t, http.StatusOK, w.Code)

	var report services.WastageReport
	err = json.Unmarshal(w.Body.Bytes(), &report)
	require.NoError(t, err)

	// Verify report data
	assert.Equal(t, 2, report.TotalExpiredBatches)
	assert.Equal(t, 300.0, report.TotalQuantityWasted) // 200 + 100
	assert.Equal(t, 40.0, report.TotalCostWasted)      // (200 * 0.15) + (100 * 0.10)
	assert.Len(t, report.WastageByType, 2)

	t.Log("\n=== Integration Test Summary ===")
	t.Log("✓ Step 1: Created expired batch records")
	t.Log("✓ Step 2: Retrieved wastage report")
	t.Logf("✓ Total expired batches: %d", report.TotalExpiredBatches)
	t.Logf("✓ Total quantity wasted: %.2f", report.TotalQuantityWasted)
	t.Logf("✓ Total cost wasted: %.2f", report.TotalCostWasted)
}

// TestBatchReportHandler_UsageReport_Integration tests the usage report endpoint
func TestBatchReportHandler_UsageReport_Integration(t *testing.T) {
	// Skip if not running integration tests
	if testing.Short() {
		t.Skip("Skipping integration test")
	}

	// Setup MongoDB connection
	ctx := context.Background()
	client, err := mongo.Connect(ctx, options.Client().ApplyURI("mongodb://localhost:27017"))
	require.NoError(t, err)
	defer client.Disconnect(ctx)

	// Use test database
	db := client.Database("cafe_pos_test_batch_report_usage")
	defer db.Drop(ctx)

	// Create repositories
	batchRecordRepo := mongodb.NewBatchRecordRepository(db)
	batchUsageLogRepo := mongodb.NewBatchUsageLogRepository(db)
	batchDefRepo := mongodb.NewBatchDefinitionRepository(db)

	// Create service and handler
	reportService := services.NewBatchReportService(batchRecordRepo, batchUsageLogRepo, batchDefRepo)
	handler := NewBatchReportHandler(reportService)

	// Setup Gin router
	gin.SetMode(gin.TestMode)
	router := gin.New()
	router.GET("/api/batch-reports/usage", handler.GetUsageReport)

	// ========================================
	// STEP 1: Create test data with usage logs
	// ========================================

	now := time.Now()
	fromDate := now.Add(-24 * time.Hour)
	toDate := now.Add(24 * time.Hour)

	batchRecordID := primitive.NewObjectID()
	orderID := primitive.NewObjectID()
	menuItemID := primitive.NewObjectID()

	// Create usage logs
	usageLog1 := &batch.BatchUsageLog{
		ID:            primitive.NewObjectID(),
		BatchRecordID: batchRecordID,
		BatchName:     "Coffee Concentrate",
		OrderID:       orderID,
		MenuItemID:    menuItemID,
		MenuItemName:  "Black Coffee",
		QuantityUsed:  30,
		Unit:          "ml",
		CostPerUnit:   0.15,
		TotalCost:     4.5,
		UsedAt:        now.Add(-12 * time.Hour),
	}
	err = batchUsageLogRepo.Create(ctx, usageLog1)
	require.NoError(t, err)

	usageLog2 := &batch.BatchUsageLog{
		ID:            primitive.NewObjectID(),
		BatchRecordID: batchRecordID,
		BatchName:     "Coffee Concentrate",
		OrderID:       primitive.NewObjectID(),
		MenuItemID:    menuItemID,
		MenuItemName:  "Black Coffee",
		QuantityUsed:  30,
		Unit:          "ml",
		CostPerUnit:   0.15,
		TotalCost:     4.5,
		UsedAt:        now.Add(-6 * time.Hour),
	}
	err = batchUsageLogRepo.Create(ctx, usageLog2)
	require.NoError(t, err)

	usageLog3 := &batch.BatchUsageLog{
		ID:            primitive.NewObjectID(),
		BatchRecordID: batchRecordID,
		BatchName:     "Coffee Concentrate",
		OrderID:       primitive.NewObjectID(),
		MenuItemID:    primitive.NewObjectID(),
		MenuItemName:  "Latte",
		QuantityUsed:  50,
		Unit:          "ml",
		CostPerUnit:   0.15,
		TotalCost:     7.5,
		UsedAt:        now.Add(-3 * time.Hour),
	}
	err = batchUsageLogRepo.Create(ctx, usageLog3)
	require.NoError(t, err)

	// ========================================
	// STEP 2: Test GET /api/batch-reports/usage
	// ========================================

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/api/batch-reports/usage?from_date="+fromDate.Format(time.RFC3339)+"&to_date="+toDate.Format(time.RFC3339), nil)
	router.ServeHTTP(w, req)

	// Verify response
	assert.Equal(t, http.StatusOK, w.Code)

	var report services.UsageReport
	err = json.Unmarshal(w.Body.Bytes(), &report)
	require.NoError(t, err)

	// Verify report data
	assert.Equal(t, 3, report.TotalUsageCount)
	assert.Equal(t, 110.0, report.TotalQuantityUsed) // 30 + 30 + 50
	assert.Equal(t, 16.5, report.TotalCost)          // 4.5 + 4.5 + 7.5
	assert.Len(t, report.ByMenuItem, 2)
	assert.Len(t, report.UsageTrend, 3) // 3 different dates

	// ========================================
	// STEP 3: Test with menu_item_id filter
	// ========================================

	w = httptest.NewRecorder()
	req, _ = http.NewRequest("GET", "/api/batch-reports/usage?from_date="+fromDate.Format(time.RFC3339)+"&to_date="+toDate.Format(time.RFC3339)+"&menu_item_id="+menuItemID.Hex(), nil)
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)

	err = json.Unmarshal(w.Body.Bytes(), &report)
	require.NoError(t, err)
	assert.Equal(t, 2, report.TotalUsageCount)
	assert.Equal(t, 60.0, report.TotalQuantityUsed) // 30 + 30

	t.Log("\n=== Integration Test Summary ===")
	t.Log("✓ Step 1: Created usage logs")
	t.Log("✓ Step 2: Retrieved usage report")
	t.Log("✓ Step 3: Tested menu_item_id filter")
	t.Logf("✓ Total usage count: %d", report.TotalUsageCount)
	t.Logf("✓ Total quantity used: %.2f", report.TotalQuantityUsed)
}

// TestBatchReportHandler_ValidationErrors tests error handling
func TestBatchReportHandler_ValidationErrors(t *testing.T) {
	// Skip if not running integration tests
	if testing.Short() {
		t.Skip("Skipping integration test")
	}

	// Setup MongoDB connection
	ctx := context.Background()
	client, err := mongo.Connect(ctx, options.Client().ApplyURI("mongodb://localhost:27017"))
	require.NoError(t, err)
	defer client.Disconnect(ctx)

	// Use test database
	db := client.Database("cafe_pos_test_batch_report_validation")
	defer db.Drop(ctx)

	// Create repositories
	batchRecordRepo := mongodb.NewBatchRecordRepository(db)
	batchUsageLogRepo := mongodb.NewBatchUsageLogRepository(db)
	batchDefRepo := mongodb.NewBatchDefinitionRepository(db)

	// Create service and handler
	reportService := services.NewBatchReportService(batchRecordRepo, batchUsageLogRepo, batchDefRepo)
	handler := NewBatchReportHandler(reportService)

	// Setup Gin router
	gin.SetMode(gin.TestMode)
	router := gin.New()
	router.GET("/api/batch-reports/production", handler.GetProductionReport)
	router.GET("/api/batch-reports/wastage", handler.GetWastageReport)
	router.GET("/api/batch-reports/usage", handler.GetUsageReport)

	// Test missing from_date
	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/api/batch-reports/production?to_date=2026-02-13T00:00:00Z", nil)
	router.ServeHTTP(w, req)
	assert.Equal(t, http.StatusBadRequest, w.Code)

	// Test missing to_date
	w = httptest.NewRecorder()
	req, _ = http.NewRequest("GET", "/api/batch-reports/production?from_date=2026-02-13T00:00:00Z", nil)
	router.ServeHTTP(w, req)
	assert.Equal(t, http.StatusBadRequest, w.Code)

	// Test invalid date format
	w = httptest.NewRecorder()
	req, _ = http.NewRequest("GET", "/api/batch-reports/production?from_date=invalid&to_date=2026-02-13T00:00:00Z", nil)
	router.ServeHTTP(w, req)
	assert.Equal(t, http.StatusBadRequest, w.Code)

	// Test invalid batch_definition_id
	w = httptest.NewRecorder()
	req, _ = http.NewRequest("GET", "/api/batch-reports/production?from_date=2026-02-13T00:00:00Z&to_date=2026-02-14T00:00:00Z&batch_definition_id=invalid", nil)
	router.ServeHTTP(w, req)
	assert.Equal(t, http.StatusBadRequest, w.Code)

	// Test invalid menu_item_id
	w = httptest.NewRecorder()
	req, _ = http.NewRequest("GET", "/api/batch-reports/usage?from_date=2026-02-13T00:00:00Z&to_date=2026-02-14T00:00:00Z&menu_item_id=invalid", nil)
	router.ServeHTTP(w, req)
	assert.Equal(t, http.StatusBadRequest, w.Code)

	t.Log("\n=== Validation Test Summary ===")
	t.Log("✓ Tested missing from_date")
	t.Log("✓ Tested missing to_date")
	t.Log("✓ Tested invalid date format")
	t.Log("✓ Tested invalid batch_definition_id")
	t.Log("✓ Tested invalid menu_item_id")
}
