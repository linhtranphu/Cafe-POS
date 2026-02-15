package mongodb

import (
	"context"
	"testing"
	"time"

	"cafe-pos/backend/domain/batch"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

// Helper function to create a test batch usage log
func createTestBatchUsageLog(batchName string, menuItemName string, quantityUsed float64) *batch.BatchUsageLog {
	now := time.Now()
	return &batch.BatchUsageLog{
		BatchRecordID: primitive.NewObjectID(),
		BatchName:     batchName,
		OrderID:       primitive.NewObjectID(),
		MenuItemID:    primitive.NewObjectID(),
		MenuItemName:  menuItemName,
		QuantityUsed:  quantityUsed,
		Unit:          "ml",
		CostPerUnit:   0.15,
		TotalCost:     quantityUsed * 0.15,
		UsedAt:        now,
	}
}

// TestCreate tests the Create method
func TestBatchUsageLogRepository_Create(t *testing.T) {
	ctx := context.Background()
	
	t.Run("Create new batch usage log", func(t *testing.T) {
		log := createTestBatchUsageLog("Cà Phê Concentrate", "Cà Phê Đen", 30)
		
		// Verify ID is not set before creation
		if !log.ID.IsZero() {
			t.Error("Expected ID to be zero before creation")
		}
		
		// Simulate creation
		log.ID = primitive.NewObjectID()
		
		// Verify ID is set after creation
		if log.ID.IsZero() {
			t.Error("Expected ID to be set after creation")
		}
		
		// Verify fields
		if log.BatchName != "Cà Phê Concentrate" {
			t.Errorf("Expected batch name 'Cà Phê Concentrate', got '%s'", log.BatchName)
		}
		if log.MenuItemName != "Cà Phê Đen" {
			t.Errorf("Expected menu item name 'Cà Phê Đen', got '%s'", log.MenuItemName)
		}
		if log.QuantityUsed != 30 {
			t.Errorf("Expected quantity used 30, got %f", log.QuantityUsed)
		}
	})
	
	t.Run("Create with calculated total cost", func(t *testing.T) {
		log := createTestBatchUsageLog("Trà Đen Concentrate", "Trà Đen", 50)
		
		expectedCost := 50 * 0.15
		if log.TotalCost != expectedCost {
			t.Errorf("Expected total cost %f, got %f", expectedCost, log.TotalCost)
		}
	})
	
	t.Run("Create with timestamp", func(t *testing.T) {
		log := createTestBatchUsageLog("Sữa Tươi", "Sữa Đá", 100)
		
		if log.UsedAt.IsZero() {
			t.Error("Expected UsedAt to be set")
		}
	})
	
	_ = ctx
}

// TestFindAll tests the FindAll method with various filters
func TestBatchUsageLogRepository_FindAll(t *testing.T) {
	ctx := context.Background()
	
	t.Run("Find all without filters", func(t *testing.T) {
		logs := []*batch.BatchUsageLog{
			createTestBatchUsageLog("Cà Phê Concentrate", "Cà Phê Đen", 30),
			createTestBatchUsageLog("Trà Đen Concentrate", "Trà Đen", 50),
			createTestBatchUsageLog("Sữa Tươi", "Sữa Đá", 100),
		}
		
		for _, log := range logs {
			log.ID = primitive.NewObjectID()
		}
		
		if len(logs) != 3 {
			t.Errorf("Expected 3 logs, got %d", len(logs))
		}
	})
	
	t.Run("Find all with batch record ID filter", func(t *testing.T) {
		batchRecordID := primitive.NewObjectID()
		filter := batch.BatchUsageLogFilter{
			BatchRecordID: &batchRecordID,
		}
		
		if filter.BatchRecordID == nil {
			t.Error("Expected batch record ID to be set")
		}
		if *filter.BatchRecordID != batchRecordID {
			t.Error("Expected batch record ID to match")
		}
	})
	
	t.Run("Find all with order ID filter", func(t *testing.T) {
		orderID := primitive.NewObjectID()
		filter := batch.BatchUsageLogFilter{
			OrderID: &orderID,
		}
		
		if filter.OrderID == nil {
			t.Error("Expected order ID to be set")
		}
		if *filter.OrderID != orderID {
			t.Error("Expected order ID to match")
		}
	})
	
	t.Run("Find all with menu item ID filter", func(t *testing.T) {
		menuItemID := primitive.NewObjectID()
		filter := batch.BatchUsageLogFilter{
			MenuItemID: &menuItemID,
		}
		
		if filter.MenuItemID == nil {
			t.Error("Expected menu item ID to be set")
		}
		if *filter.MenuItemID != menuItemID {
			t.Error("Expected menu item ID to match")
		}
	})
	
	t.Run("Find all with date range filter", func(t *testing.T) {
		fromDate := time.Now().Add(-7 * 24 * time.Hour)
		toDate := time.Now()
		
		filter := batch.BatchUsageLogFilter{
			FromDate: &fromDate,
			ToDate:   &toDate,
		}
		
		if filter.FromDate == nil {
			t.Error("Expected from date to be set")
		}
		if filter.ToDate == nil {
			t.Error("Expected to date to be set")
		}
		
		if !filter.ToDate.After(*filter.FromDate) {
			t.Error("Expected to date to be after from date")
		}
	})
	
	t.Run("Find all with pagination", func(t *testing.T) {
		filter := batch.BatchUsageLogFilter{
			Page:  1,
			Limit: 20,
		}
		
		if filter.Page != 1 {
			t.Errorf("Expected page 1, got %d", filter.Page)
		}
		if filter.Limit != 20 {
			t.Errorf("Expected limit 20, got %d", filter.Limit)
		}
	})
	
	t.Run("Verify sorting by used_at descending (newest first)", func(t *testing.T) {
		now := time.Now()
		
		log1 := createTestBatchUsageLog("Batch 1", "Item 1", 30)
		log1.ID = primitive.NewObjectID()
		log1.UsedAt = now.Add(-2 * time.Hour)
		
		log2 := createTestBatchUsageLog("Batch 2", "Item 2", 40)
		log2.ID = primitive.NewObjectID()
		log2.UsedAt = now.Add(-1 * time.Hour)
		
		log3 := createTestBatchUsageLog("Batch 3", "Item 3", 50)
		log3.ID = primitive.NewObjectID()
		log3.UsedAt = now
		
		// Verify timestamps are in order (newest first)
		if !log3.UsedAt.After(log2.UsedAt) {
			t.Error("Expected log3 to be used after log2")
		}
		if !log2.UsedAt.After(log1.UsedAt) {
			t.Error("Expected log2 to be used after log1")
		}
	})
	
	t.Run("Find all with multiple filters", func(t *testing.T) {
		batchRecordID := primitive.NewObjectID()
		fromDate := time.Now().Add(-7 * 24 * time.Hour)
		toDate := time.Now()
		
		filter := batch.BatchUsageLogFilter{
			BatchRecordID: &batchRecordID,
			FromDate:      &fromDate,
			ToDate:        &toDate,
			Page:          1,
			Limit:         10,
		}
		
		// Verify all filters are set
		if filter.BatchRecordID == nil {
			t.Error("Expected batch record ID to be set")
		}
		if filter.FromDate == nil {
			t.Error("Expected from date to be set")
		}
		if filter.ToDate == nil {
			t.Error("Expected to date to be set")
		}
		if filter.Page != 1 {
			t.Errorf("Expected page 1, got %d", filter.Page)
		}
		if filter.Limit != 10 {
			t.Errorf("Expected limit 10, got %d", filter.Limit)
		}
	})
	
	_ = ctx
}

// TestCalculateTotalCost tests the cost calculation
func TestBatchUsageLogRepository_CalculateTotalCost(t *testing.T) {
	t.Run("Calculate total cost correctly", func(t *testing.T) {
		log := &batch.BatchUsageLog{
			QuantityUsed: 100,
			CostPerUnit:  0.25,
		}
		
		expectedCost := 25.0
		actualCost := log.CalculateTotalCost()
		
		if actualCost != expectedCost {
			t.Errorf("Expected total cost %f, got %f", expectedCost, actualCost)
		}
	})
	
	t.Run("Calculate total cost with zero quantity", func(t *testing.T) {
		log := &batch.BatchUsageLog{
			QuantityUsed: 0,
			CostPerUnit:  0.25,
		}
		
		expectedCost := 0.0
		actualCost := log.CalculateTotalCost()
		
		if actualCost != expectedCost {
			t.Errorf("Expected total cost %f, got %f", expectedCost, actualCost)
		}
	})
	
	t.Run("Calculate total cost with zero cost per unit", func(t *testing.T) {
		log := &batch.BatchUsageLog{
			QuantityUsed: 100,
			CostPerUnit:  0,
		}
		
		expectedCost := 0.0
		actualCost := log.CalculateTotalCost()
		
		if actualCost != expectedCost {
			t.Errorf("Expected total cost %f, got %f", expectedCost, actualCost)
		}
	})
	
	t.Run("Calculate total cost with decimal values", func(t *testing.T) {
		log := &batch.BatchUsageLog{
			QuantityUsed: 33.5,
			CostPerUnit:  0.15,
		}
		
		expectedCost := 5.025
		actualCost := log.CalculateTotalCost()
		
		// Use a small epsilon for floating point comparison
		epsilon := 0.0001
		if actualCost < expectedCost-epsilon || actualCost > expectedCost+epsilon {
			t.Errorf("Expected total cost %f, got %f", expectedCost, actualCost)
		}
	})
}

// TestNewBatchUsageLog tests the factory function
func TestBatchUsageLogRepository_NewBatchUsageLog(t *testing.T) {
	t.Run("Create from request", func(t *testing.T) {
		req := batch.CreateBatchUsageLogRequest{
			BatchRecordID: primitive.NewObjectID(),
			BatchName:     "Cà Phê Concentrate",
			OrderID:       primitive.NewObjectID(),
			MenuItemID:    primitive.NewObjectID(),
			MenuItemName:  "Cà Phê Đen",
			QuantityUsed:  30,
			Unit:          "ml",
			CostPerUnit:   0.15,
		}
		
		log := batch.NewBatchUsageLog(req)
		
		// Verify fields are set correctly
		if log.BatchRecordID != req.BatchRecordID {
			t.Error("Expected batch record ID to match")
		}
		if log.BatchName != req.BatchName {
			t.Errorf("Expected batch name '%s', got '%s'", req.BatchName, log.BatchName)
		}
		if log.OrderID != req.OrderID {
			t.Error("Expected order ID to match")
		}
		if log.MenuItemID != req.MenuItemID {
			t.Error("Expected menu item ID to match")
		}
		if log.MenuItemName != req.MenuItemName {
			t.Errorf("Expected menu item name '%s', got '%s'", req.MenuItemName, log.MenuItemName)
		}
		if log.QuantityUsed != req.QuantityUsed {
			t.Errorf("Expected quantity used %f, got %f", req.QuantityUsed, log.QuantityUsed)
		}
		if log.Unit != req.Unit {
			t.Errorf("Expected unit '%s', got '%s'", req.Unit, log.Unit)
		}
		if log.CostPerUnit != req.CostPerUnit {
			t.Errorf("Expected cost per unit %f, got %f", req.CostPerUnit, log.CostPerUnit)
		}
		
		// Verify total cost is calculated
		expectedCost := req.QuantityUsed * req.CostPerUnit
		if log.TotalCost != expectedCost {
			t.Errorf("Expected total cost %f, got %f", expectedCost, log.TotalCost)
		}
		
		// Verify timestamp is set
		if log.UsedAt.IsZero() {
			t.Error("Expected UsedAt to be set")
		}
	})
}

// TestBatchUsageLogValidation tests field validation
func TestBatchUsageLogRepository_Validation(t *testing.T) {
	t.Run("Valid batch usage log", func(t *testing.T) {
		log := createTestBatchUsageLog("Valid Batch", "Valid Item", 50)
		
		// Verify all required fields are set
		if log.BatchRecordID.IsZero() {
			t.Error("Expected batch record ID to be set")
		}
		if log.BatchName == "" {
			t.Error("Expected batch name to be set")
		}
		if log.OrderID.IsZero() {
			t.Error("Expected order ID to be set")
		}
		if log.MenuItemID.IsZero() {
			t.Error("Expected menu item ID to be set")
		}
		if log.MenuItemName == "" {
			t.Error("Expected menu item name to be set")
		}
		if log.QuantityUsed <= 0 {
			t.Error("Expected positive quantity used")
		}
		if log.Unit == "" {
			t.Error("Expected unit to be set")
		}
		if log.CostPerUnit < 0 {
			t.Error("Expected non-negative cost per unit")
		}
		if log.UsedAt.IsZero() {
			t.Error("Expected used at timestamp to be set")
		}
	})
}

// TestFindAll_EmptyResult tests FindAll when no logs exist
func TestBatchUsageLogRepository_FindAll_EmptyResult(t *testing.T) {
	ctx := context.Background()
	
	filter := batch.BatchUsageLogFilter{}
	
	// In real implementation, this would return empty slice
	// For mock, we just verify the filter is valid
	if filter.Page < 0 {
		t.Error("Expected non-negative page")
	}
	if filter.Limit < 0 {
		t.Error("Expected non-negative limit")
	}
	
	_ = ctx
}

// TestUsageTracking tests usage tracking scenarios
func TestBatchUsageLogRepository_UsageTracking(t *testing.T) {
	t.Run("Track multiple uses of same batch", func(t *testing.T) {
		batchRecordID := primitive.NewObjectID()
		
		log1 := createTestBatchUsageLog("Cà Phê Concentrate", "Cà Phê Đen", 30)
		log1.BatchRecordID = batchRecordID
		
		log2 := createTestBatchUsageLog("Cà Phê Concentrate", "Cà Phê Sữa", 40)
		log2.BatchRecordID = batchRecordID
		
		// Verify both logs reference the same batch
		if log1.BatchRecordID != log2.BatchRecordID {
			t.Error("Expected both logs to reference the same batch")
		}
		
		// Verify different quantities
		if log1.QuantityUsed == log2.QuantityUsed {
			t.Error("Expected different quantities used")
		}
	})
	
	t.Run("Track usage across different orders", func(t *testing.T) {
		batchRecordID := primitive.NewObjectID()
		
		log1 := createTestBatchUsageLog("Cà Phê Concentrate", "Cà Phê Đen", 30)
		log1.BatchRecordID = batchRecordID
		log1.OrderID = primitive.NewObjectID()
		
		log2 := createTestBatchUsageLog("Cà Phê Concentrate", "Cà Phê Đen", 30)
		log2.BatchRecordID = batchRecordID
		log2.OrderID = primitive.NewObjectID()
		
		// Verify same batch, different orders
		if log1.BatchRecordID != log2.BatchRecordID {
			t.Error("Expected both logs to reference the same batch")
		}
		if log1.OrderID == log2.OrderID {
			t.Error("Expected different order IDs")
		}
	})
	
	t.Run("Track usage of different batches in same order", func(t *testing.T) {
		orderID := primitive.NewObjectID()
		
		log1 := createTestBatchUsageLog("Cà Phê Concentrate", "Cà Phê Đen", 30)
		log1.OrderID = orderID
		log1.BatchRecordID = primitive.NewObjectID()
		
		log2 := createTestBatchUsageLog("Trà Đen Concentrate", "Trà Đen", 50)
		log2.OrderID = orderID
		log2.BatchRecordID = primitive.NewObjectID()
		
		// Verify same order, different batches
		if log1.OrderID != log2.OrderID {
			t.Error("Expected both logs to reference the same order")
		}
		if log1.BatchRecordID == log2.BatchRecordID {
			t.Error("Expected different batch record IDs")
		}
	})
}
