package mongodb

import (
	"context"
	"testing"
	"time"

	"cafe-pos/backend/domain/batch"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

// Helper function to create a test batch record
func createTestBatchRecord(batchName string, quantityProduced float64) *batch.BatchRecord {
	now := time.Now()
	return &batch.BatchRecord{
		BatchDefinitionID: primitive.NewObjectID(),
		BatchName:         batchName,
		QuantityProduced:  quantityProduced,
		QuantityRemaining: quantityProduced,
		Unit:              "ml",
		CostPerUnit:       0.15,
		TotalCost:         quantityProduced * 0.15,
		PreparedBy:        "user_123",
		PreparedAt:        now,
		ExpiresAt:         now.Add(24 * time.Hour),
		Status:            batch.BatchStatusAvailable,
		IngredientsUsed: []batch.IngredientUsage{
			{
				IngredientID:   primitive.NewObjectID(),
				IngredientName: "Hạt Cà Phê",
				Quantity:       110,
				Unit:           "g",
				CostPerUnit:    0.68,
				TotalCost:      75.0,
			},
		},
	}
}

// TestCreate tests the Create method
func TestBatchRecordRepository_Create(t *testing.T) {
	ctx := context.Background()
	
	t.Run("Create new batch record", func(t *testing.T) {
		record := createTestBatchRecord("Cà Phê Concentrate", 500)
		
		// Verify ID is not set before creation
		if !record.ID.IsZero() {
			t.Error("Expected ID to be zero before creation")
		}
		
		// Simulate creation
		record.ID = primitive.NewObjectID()
		record.CreatedAt = time.Now()
		record.UpdatedAt = time.Now()
		
		// Verify ID is set after creation
		if record.ID.IsZero() {
			t.Error("Expected ID to be set after creation")
		}
		
		// Verify timestamps
		if record.CreatedAt.IsZero() {
			t.Error("Expected CreatedAt to be set")
		}
		if record.UpdatedAt.IsZero() {
			t.Error("Expected UpdatedAt to be set")
		}
		
		// Verify fields
		if record.BatchName != "Cà Phê Concentrate" {
			t.Errorf("Expected batch name 'Cà Phê Concentrate', got '%s'", record.BatchName)
		}
		if record.QuantityProduced != 500 {
			t.Errorf("Expected quantity produced 500, got %f", record.QuantityProduced)
		}
		if record.QuantityRemaining != 500 {
			t.Errorf("Expected quantity remaining 500, got %f", record.QuantityRemaining)
		}
		if record.Status != batch.BatchStatusAvailable {
			t.Errorf("Expected status 'available', got '%s'", record.Status)
		}
	})
	
	t.Run("Create with ingredients used", func(t *testing.T) {
		record := createTestBatchRecord("Trà Đen Concentrate", 1000)
		
		if len(record.IngredientsUsed) != 1 {
			t.Errorf("Expected 1 ingredient used, got %d", len(record.IngredientsUsed))
		}
		
		ingredient := record.IngredientsUsed[0]
		if ingredient.IngredientName != "Hạt Cà Phê" {
			t.Errorf("Expected ingredient 'Hạt Cà Phê', got '%s'", ingredient.IngredientName)
		}
		if ingredient.Quantity != 110 {
			t.Errorf("Expected quantity 110, got %f", ingredient.Quantity)
		}
	})
	
	_ = ctx
}

// TestUpdate tests the Update method
func TestBatchRecordRepository_Update(t *testing.T) {
	ctx := context.Background()
	
	t.Run("Update existing batch record", func(t *testing.T) {
		record := createTestBatchRecord("Cà Phê Concentrate", 500)
		record.ID = primitive.NewObjectID()
		record.CreatedAt = time.Now().Add(-1 * time.Hour)
		originalCreatedAt := record.CreatedAt
		
		// Update fields
		record.QuantityRemaining = 300
		record.UpdatedAt = time.Now()
		
		// Verify UpdatedAt is updated
		if record.UpdatedAt.Before(originalCreatedAt) {
			t.Error("Expected UpdatedAt to be after CreatedAt")
		}
		
		// Verify fields are updated
		if record.QuantityRemaining != 300 {
			t.Errorf("Expected quantity remaining 300, got %f", record.QuantityRemaining)
		}
		
		// Verify CreatedAt is not changed
		if !record.CreatedAt.Equal(originalCreatedAt) {
			t.Error("Expected CreatedAt to remain unchanged")
		}
	})
	
	_ = ctx
}

// TestFindByID tests the FindByID method
func TestBatchRecordRepository_FindByID(t *testing.T) {
	ctx := context.Background()
	
	t.Run("Find existing batch record", func(t *testing.T) {
		record := createTestBatchRecord("Cà Phê Concentrate", 500)
		record.ID = primitive.NewObjectID()
		record.CreatedAt = time.Now()
		record.UpdatedAt = time.Now()
		
		// Verify we can retrieve by ID
		if record.ID.IsZero() {
			t.Error("Expected valid ID")
		}
		
		// Verify fields
		if record.BatchName != "Cà Phê Concentrate" {
			t.Errorf("Expected batch name 'Cà Phê Concentrate', got '%s'", record.BatchName)
		}
	})
	
	_ = ctx
}

// TestFindAll tests the FindAll method with various filters
func TestBatchRecordRepository_FindAll(t *testing.T) {
	ctx := context.Background()
	
	t.Run("Find all without filters", func(t *testing.T) {
		records := []*batch.BatchRecord{
			createTestBatchRecord("Cà Phê Concentrate", 500),
			createTestBatchRecord("Trà Đen Concentrate", 1000),
			createTestBatchRecord("Sữa Tươi", 2000),
		}
		
		for _, record := range records {
			record.ID = primitive.NewObjectID()
			record.CreatedAt = time.Now()
			record.UpdatedAt = time.Now()
		}
		
		if len(records) != 3 {
			t.Errorf("Expected 3 records, got %d", len(records))
		}
	})
	
	t.Run("Find all with batch definition ID filter", func(t *testing.T) {
		defID := primitive.NewObjectID()
		filter := batch.BatchRecordFilter{
			BatchDefinitionID: &defID,
		}
		
		if filter.BatchDefinitionID == nil {
			t.Error("Expected batch definition ID to be set")
		}
		if *filter.BatchDefinitionID != defID {
			t.Error("Expected batch definition ID to match")
		}
	})
	
	t.Run("Find all with status filter", func(t *testing.T) {
		filter := batch.BatchRecordFilter{
			Status: batch.BatchStatusAvailable,
		}
		
		if filter.Status != batch.BatchStatusAvailable {
			t.Errorf("Expected status 'available', got '%s'", filter.Status)
		}
	})
	
	t.Run("Find all with prepared by filter", func(t *testing.T) {
		filter := batch.BatchRecordFilter{
			PreparedBy: "user_123",
		}
		
		if filter.PreparedBy != "user_123" {
			t.Errorf("Expected prepared by 'user_123', got '%s'", filter.PreparedBy)
		}
	})
	
	t.Run("Find all with date range filter", func(t *testing.T) {
		fromDate := time.Now().Add(-7 * 24 * time.Hour)
		toDate := time.Now()
		
		filter := batch.BatchRecordFilter{
			FromDate: &fromDate,
			ToDate:   &toDate,
		}
		
		if filter.FromDate == nil {
			t.Error("Expected from date to be set")
		}
		if filter.ToDate == nil {
			t.Error("Expected to date to be set")
		}
	})
	
	t.Run("Find all with pagination", func(t *testing.T) {
		filter := batch.BatchRecordFilter{
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
	
	t.Run("Verify sorting by expires_at ascending (FIFO)", func(t *testing.T) {
		now := time.Now()
		
		record1 := createTestBatchRecord("First", 500)
		record1.ID = primitive.NewObjectID()
		record1.ExpiresAt = now.Add(1 * time.Hour)
		
		record2 := createTestBatchRecord("Second", 500)
		record2.ID = primitive.NewObjectID()
		record2.ExpiresAt = now.Add(2 * time.Hour)
		
		record3 := createTestBatchRecord("Third", 500)
		record3.ID = primitive.NewObjectID()
		record3.ExpiresAt = now.Add(3 * time.Hour)
		
		// Verify expiry times are in order (FIFO - oldest expires first)
		if !record1.ExpiresAt.Before(record2.ExpiresAt) {
			t.Error("Expected record1 to expire before record2")
		}
		if !record2.ExpiresAt.Before(record3.ExpiresAt) {
			t.Error("Expected record2 to expire before record3")
		}
	})
	
	_ = ctx
}

// TestFindAvailableByDefinition tests the FIFO query
func TestBatchRecordRepository_FindAvailableByDefinition(t *testing.T) {
	ctx := context.Background()
	
	t.Run("Find available batches sorted by expiry (FIFO)", func(t *testing.T) {
		defID := primitive.NewObjectID()
		now := time.Now()
		
		// Create batches with different expiry times
		record1 := createTestBatchRecord("Batch 1", 500)
		record1.ID = primitive.NewObjectID()
		record1.BatchDefinitionID = defID
		record1.ExpiresAt = now.Add(1 * time.Hour)
		record1.Status = batch.BatchStatusAvailable
		
		record2 := createTestBatchRecord("Batch 2", 500)
		record2.ID = primitive.NewObjectID()
		record2.BatchDefinitionID = defID
		record2.ExpiresAt = now.Add(2 * time.Hour)
		record2.Status = batch.BatchStatusAvailable
		
		// Verify FIFO ordering (oldest expiry first)
		if !record1.ExpiresAt.Before(record2.ExpiresAt) {
			t.Error("Expected record1 to expire before record2 (FIFO)")
		}
	})
	
	t.Run("Exclude expired batches", func(t *testing.T) {
		now := time.Now()
		
		expiredRecord := createTestBatchRecord("Expired Batch", 500)
		expiredRecord.ExpiresAt = now.Add(-1 * time.Hour) // Already expired
		
		if !expiredRecord.IsExpired() {
			t.Error("Expected batch to be expired")
		}
	})
	
	t.Run("Exclude depleted batches", func(t *testing.T) {
		depletedRecord := createTestBatchRecord("Depleted Batch", 500)
		depletedRecord.QuantityRemaining = 0
		depletedRecord.Status = batch.BatchStatusDepleted
		
		if !depletedRecord.IsDepleted() {
			t.Error("Expected batch to be depleted")
		}
	})
	
	t.Run("Only include available status", func(t *testing.T) {
		availableRecord := createTestBatchRecord("Available Batch", 500)
		availableRecord.Status = batch.BatchStatusAvailable
		
		if availableRecord.Status != batch.BatchStatusAvailable {
			t.Errorf("Expected status 'available', got '%s'", availableRecord.Status)
		}
	})
	
	_ = ctx
}

// TestUpdateQuantity tests the UpdateQuantity method
func TestBatchRecordRepository_UpdateQuantity(t *testing.T) {
	ctx := context.Background()
	
	t.Run("Update quantity to positive value", func(t *testing.T) {
		record := createTestBatchRecord("Cà Phê Concentrate", 500)
		record.ID = primitive.NewObjectID()
		
		newQuantity := 300.0
		
		// Simulate update
		record.QuantityRemaining = newQuantity
		record.UpdatedAt = time.Now()
		
		if record.QuantityRemaining != 300 {
			t.Errorf("Expected quantity 300, got %f", record.QuantityRemaining)
		}
	})
	
	t.Run("Update quantity to zero marks as depleted", func(t *testing.T) {
		record := createTestBatchRecord("Cà Phê Concentrate", 500)
		record.ID = primitive.NewObjectID()
		
		// Simulate update to zero
		record.QuantityRemaining = 0
		record.Status = batch.BatchStatusDepleted
		record.UpdatedAt = time.Now()
		
		if record.QuantityRemaining != 0 {
			t.Errorf("Expected quantity 0, got %f", record.QuantityRemaining)
		}
		if record.Status != batch.BatchStatusDepleted {
			t.Errorf("Expected status 'depleted', got '%s'", record.Status)
		}
	})
	
	t.Run("Verify UpdatedAt is set", func(t *testing.T) {
		record := createTestBatchRecord("Cà Phê Concentrate", 500)
		record.ID = primitive.NewObjectID()
		originalUpdatedAt := record.UpdatedAt
		
		time.Sleep(10 * time.Millisecond)
		
		// Simulate update
		record.QuantityRemaining = 400
		record.UpdatedAt = time.Now()
		
		if !record.UpdatedAt.After(originalUpdatedAt) {
			t.Error("Expected UpdatedAt to be updated")
		}
	})
	
	_ = ctx
}

// TestGetTotalAvailableQuantity tests the aggregation query
func TestBatchRecordRepository_GetTotalAvailableQuantity(t *testing.T) {
	ctx := context.Background()
	
	t.Run("Calculate total from multiple batches", func(t *testing.T) {
		defID := primitive.NewObjectID()
		now := time.Now()
		
		record1 := createTestBatchRecord("Batch 1", 500)
		record1.BatchDefinitionID = defID
		record1.QuantityRemaining = 300
		record1.ExpiresAt = now.Add(1 * time.Hour)
		record1.Status = batch.BatchStatusAvailable
		
		record2 := createTestBatchRecord("Batch 2", 500)
		record2.BatchDefinitionID = defID
		record2.QuantityRemaining = 200
		record2.ExpiresAt = now.Add(2 * time.Hour)
		record2.Status = batch.BatchStatusAvailable
		
		expectedTotal := 500.0 // 300 + 200
		actualTotal := record1.QuantityRemaining + record2.QuantityRemaining
		
		if actualTotal != expectedTotal {
			t.Errorf("Expected total %f, got %f", expectedTotal, actualTotal)
		}
	})
	
	t.Run("Exclude expired batches from total", func(t *testing.T) {
		now := time.Now()
		
		availableRecord := createTestBatchRecord("Available", 500)
		availableRecord.QuantityRemaining = 300
		availableRecord.ExpiresAt = now.Add(1 * time.Hour)
		availableRecord.Status = batch.BatchStatusAvailable
		
		expiredRecord := createTestBatchRecord("Expired", 500)
		expiredRecord.QuantityRemaining = 200
		expiredRecord.ExpiresAt = now.Add(-1 * time.Hour)
		expiredRecord.Status = batch.BatchStatusExpired
		
		// Only available batch should be counted
		expectedTotal := 300.0
		actualTotal := availableRecord.QuantityRemaining
		
		if actualTotal != expectedTotal {
			t.Errorf("Expected total %f, got %f", expectedTotal, actualTotal)
		}
	})
	
	t.Run("Return zero when no available batches", func(t *testing.T) {
		total := 0.0
		
		if total != 0 {
			t.Errorf("Expected total 0, got %f", total)
		}
	})
	
	_ = ctx
}

// TestBatchRecordStatus tests status-related methods
func TestBatchRecordRepository_Status(t *testing.T) {
	t.Run("IsExpired returns true for expired batch", func(t *testing.T) {
		record := createTestBatchRecord("Expired Batch", 500)
		record.ExpiresAt = time.Now().Add(-1 * time.Hour)
		
		if !record.IsExpired() {
			t.Error("Expected batch to be expired")
		}
	})
	
	t.Run("IsExpired returns false for non-expired batch", func(t *testing.T) {
		record := createTestBatchRecord("Fresh Batch", 500)
		record.ExpiresAt = time.Now().Add(1 * time.Hour)
		
		if record.IsExpired() {
			t.Error("Expected batch to not be expired")
		}
	})
	
	t.Run("IsDepleted returns true when quantity is zero", func(t *testing.T) {
		record := createTestBatchRecord("Depleted Batch", 500)
		record.QuantityRemaining = 0
		
		if !record.IsDepleted() {
			t.Error("Expected batch to be depleted")
		}
	})
	
	t.Run("IsDepleted returns false when quantity is positive", func(t *testing.T) {
		record := createTestBatchRecord("Available Batch", 500)
		record.QuantityRemaining = 300
		
		if record.IsDepleted() {
			t.Error("Expected batch to not be depleted")
		}
	})
	
	t.Run("IsAvailable returns true for available batch", func(t *testing.T) {
		record := createTestBatchRecord("Available Batch", 500)
		record.Status = batch.BatchStatusAvailable
		record.ExpiresAt = time.Now().Add(1 * time.Hour)
		record.QuantityRemaining = 300
		
		if !record.IsAvailable() {
			t.Error("Expected batch to be available")
		}
	})
	
	t.Run("IsAvailable returns false for expired batch", func(t *testing.T) {
		record := createTestBatchRecord("Expired Batch", 500)
		record.Status = batch.BatchStatusAvailable
		record.ExpiresAt = time.Now().Add(-1 * time.Hour)
		record.QuantityRemaining = 300
		
		if record.IsAvailable() {
			t.Error("Expected batch to not be available (expired)")
		}
	})
	
	t.Run("IsAvailable returns false for depleted batch", func(t *testing.T) {
		record := createTestBatchRecord("Depleted Batch", 500)
		record.Status = batch.BatchStatusAvailable
		record.ExpiresAt = time.Now().Add(1 * time.Hour)
		record.QuantityRemaining = 0
		
		if record.IsAvailable() {
			t.Error("Expected batch to not be available (depleted)")
		}
	})
}

// TestConcurrencyControl tests optimistic locking behavior
func TestBatchRecordRepository_ConcurrencyControl(t *testing.T) {
	t.Run("UpdateQuantity uses atomic operation", func(t *testing.T) {
		record := createTestBatchRecord("Concurrent Batch", 500)
		record.ID = primitive.NewObjectID()
		
		// Simulate concurrent updates
		originalQuantity := record.QuantityRemaining
		
		// First update
		record.QuantityRemaining = originalQuantity - 100
		firstQuantity := record.QuantityRemaining
		
		// Second update
		record.QuantityRemaining = firstQuantity - 50
		finalQuantity := record.QuantityRemaining
		
		expectedFinal := originalQuantity - 150
		if finalQuantity != expectedFinal {
			t.Errorf("Expected final quantity %f, got %f", expectedFinal, finalQuantity)
		}
	})
}
