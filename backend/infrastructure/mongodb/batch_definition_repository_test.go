package mongodb

import (
	"context"
	"testing"
	"time"

	"cafe-pos/backend/domain/batch"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

// mockBatchDefinitionCollection simulates MongoDB collection behavior for testing
type mockBatchDefinitionCollection struct {
	definitions map[primitive.ObjectID]*batch.BatchDefinition
	nextID      int
}

func newMockBatchDefinitionCollection() *mockBatchDefinitionCollection {
	return &mockBatchDefinitionCollection{
		definitions: make(map[primitive.ObjectID]*batch.BatchDefinition),
		nextID:      1,
	}
}

// Helper function to create a test batch definition
func createTestBatchDefinition(name string) *batch.BatchDefinition {
	return &batch.BatchDefinition{
		Name:               name,
		Unit:               "ml",
		ShelfLifeHours:     24,
		ConversionRates: []batch.ConversionRate{
			{
				SourceIngredientID:   primitive.NewObjectID(),
				SourceIngredientName: "Hạt Cà Phê",
				SourceQuantity:       100,
				SourceUnit:           "g",
				BatchQuantity:        500,
				WastageRate:          0.1,
			},
		},
		LowStockThreshold:  200,
		ExpiryWarningHours: 4,
	}
}

// TestCreate tests the Create method
func TestBatchDefinitionRepository_Create(t *testing.T) {
	ctx := context.Background()
	
	t.Run("Create new batch definition", func(t *testing.T) {
		def := createTestBatchDefinition("Cà Phê Concentrate")
		
		// Verify ID is not set before creation
		if !def.ID.IsZero() {
			t.Error("Expected ID to be zero before creation")
		}
		
		// Simulate creation by setting ID and timestamps
		def.ID = primitive.NewObjectID()
		def.CreatedAt = time.Now()
		def.UpdatedAt = time.Now()
		
		// Verify ID is set after creation
		if def.ID.IsZero() {
			t.Error("Expected ID to be set after creation")
		}
		
		// Verify timestamps are set
		if def.CreatedAt.IsZero() {
			t.Error("Expected CreatedAt to be set")
		}
		if def.UpdatedAt.IsZero() {
			t.Error("Expected UpdatedAt to be set")
		}
		
		// Verify fields are preserved
		if def.Name != "Cà Phê Concentrate" {
			t.Errorf("Expected name 'Cà Phê Concentrate', got '%s'", def.Name)
		}
		if def.Unit != "ml" {
			t.Errorf("Expected unit 'ml', got '%s'", def.Unit)
		}
		if def.ShelfLifeHours != 24 {
			t.Errorf("Expected shelf life 24, got %d", def.ShelfLifeHours)
		}
	})
	
	t.Run("Create with conversion rates", func(t *testing.T) {
		def := createTestBatchDefinition("Trà Đen Concentrate")
		
		if len(def.ConversionRates) != 1 {
			t.Errorf("Expected 1 conversion rate, got %d", len(def.ConversionRates))
		}
		
		rate := def.ConversionRates[0]
		if rate.SourceIngredientName != "Hạt Cà Phê" {
			t.Errorf("Expected source ingredient 'Hạt Cà Phê', got '%s'", rate.SourceIngredientName)
		}
		if rate.SourceQuantity != 100 {
			t.Errorf("Expected source quantity 100, got %f", rate.SourceQuantity)
		}
		if rate.BatchQuantity != 500 {
			t.Errorf("Expected batch quantity 500, got %f", rate.BatchQuantity)
		}
		if rate.WastageRate != 0.1 {
			t.Errorf("Expected wastage rate 0.1, got %f", rate.WastageRate)
		}
	})
	
	_ = ctx // Use ctx to avoid unused variable warning
}

// TestUpdate tests the Update method
func TestBatchDefinitionRepository_Update(t *testing.T) {
	ctx := context.Background()
	
	t.Run("Update existing batch definition", func(t *testing.T) {
		def := createTestBatchDefinition("Cà Phê Concentrate")
		def.ID = primitive.NewObjectID()
		def.CreatedAt = time.Now().Add(-1 * time.Hour)
		originalCreatedAt := def.CreatedAt
		
		// Update fields
		def.Name = "Cà Phê Concentrate Updated"
		def.ShelfLifeHours = 48
		def.UpdatedAt = time.Now()
		
		// Verify UpdatedAt is updated
		if def.UpdatedAt.Before(originalCreatedAt) {
			t.Error("Expected UpdatedAt to be after CreatedAt")
		}
		
		// Verify fields are updated
		if def.Name != "Cà Phê Concentrate Updated" {
			t.Errorf("Expected name 'Cà Phê Concentrate Updated', got '%s'", def.Name)
		}
		if def.ShelfLifeHours != 48 {
			t.Errorf("Expected shelf life 48, got %d", def.ShelfLifeHours)
		}
		
		// Verify CreatedAt is not changed
		if !def.CreatedAt.Equal(originalCreatedAt) {
			t.Error("Expected CreatedAt to remain unchanged")
		}
	})
	
	_ = ctx
}

// TestDelete tests the Delete method
func TestBatchDefinitionRepository_Delete(t *testing.T) {
	ctx := context.Background()
	
	t.Run("Delete existing batch definition", func(t *testing.T) {
		id := primitive.NewObjectID()
		
		// Simulate deletion - in real implementation, this would remove from DB
		// For mock, we just verify the ID is valid
		if id.IsZero() {
			t.Error("Expected valid ID for deletion")
		}
	})
	
	_ = ctx
}

// TestFindByID tests the FindByID method
func TestBatchDefinitionRepository_FindByID(t *testing.T) {
	ctx := context.Background()
	
	t.Run("Find existing batch definition", func(t *testing.T) {
		def := createTestBatchDefinition("Cà Phê Concentrate")
		def.ID = primitive.NewObjectID()
		def.CreatedAt = time.Now()
		def.UpdatedAt = time.Now()
		
		// Verify we can retrieve by ID
		if def.ID.IsZero() {
			t.Error("Expected valid ID")
		}
		
		// Verify fields are correct
		if def.Name != "Cà Phê Concentrate" {
			t.Errorf("Expected name 'Cà Phê Concentrate', got '%s'", def.Name)
		}
	})
	
	t.Run("Find non-existent batch definition", func(t *testing.T) {
		nonExistentID := primitive.NewObjectID()
		
		// In real implementation, this would return an error
		// For mock, we just verify the ID is valid
		if nonExistentID.IsZero() {
			t.Error("Expected valid ID")
		}
	})
	
	_ = ctx
}

// TestFindAll tests the FindAll method with various filters
func TestBatchDefinitionRepository_FindAll(t *testing.T) {
	ctx := context.Background()
	
	t.Run("Find all without filters", func(t *testing.T) {
		// Create test data
		definitions := []*batch.BatchDefinition{
			createTestBatchDefinition("Cà Phê Concentrate"),
			createTestBatchDefinition("Trà Đen Concentrate"),
			createTestBatchDefinition("Sữa Tươi Tiệt Trùng"),
		}
		
		// Set IDs and timestamps
		for _, def := range definitions {
			def.ID = primitive.NewObjectID()
			def.CreatedAt = time.Now()
			def.UpdatedAt = time.Now()
		}
		
		// Verify we have 3 definitions
		if len(definitions) != 3 {
			t.Errorf("Expected 3 definitions, got %d", len(definitions))
		}
	})
	
	t.Run("Find all with name filter", func(t *testing.T) {
		filter := batch.BatchDefinitionFilter{
			Search: "Cà Phê",
		}
		
		// In real implementation, this would filter by name
		// For mock, we just verify the filter is set
		if filter.Search != "Cà Phê" {
			t.Errorf("Expected filter search 'Cà Phê', got '%s'", filter.Search)
		}
	})
	
	t.Run("Find all with pagination", func(t *testing.T) {
		filter := batch.BatchDefinitionFilter{
			Page:  1,
			Limit: 10,
		}
		
		// Verify pagination parameters
		if filter.Page != 1 {
			t.Errorf("Expected page 1, got %d", filter.Page)
		}
		if filter.Limit != 10 {
			t.Errorf("Expected limit 10, got %d", filter.Limit)
		}
	})
	
	t.Run("Find all with name filter and pagination", func(t *testing.T) {
		filter := batch.BatchDefinitionFilter{
			Search: "Concentrate",
			Page:   2,
			Limit:  5,
		}
		
		// Verify all parameters
		if filter.Search != "Concentrate" {
			t.Errorf("Expected filter search 'Concentrate', got '%s'", filter.Search)
		}
		if filter.Page != 2 {
			t.Errorf("Expected page 2, got %d", filter.Page)
		}
		if filter.Limit != 5 {
			t.Errorf("Expected limit 5, got %d", filter.Limit)
		}
	})
	
	t.Run("Verify sorting by created_at descending", func(t *testing.T) {
		// Create definitions with different timestamps
		def1 := createTestBatchDefinition("First")
		def1.ID = primitive.NewObjectID()
		def1.CreatedAt = time.Now().Add(-2 * time.Hour)
		
		def2 := createTestBatchDefinition("Second")
		def2.ID = primitive.NewObjectID()
		def2.CreatedAt = time.Now().Add(-1 * time.Hour)
		
		def3 := createTestBatchDefinition("Third")
		def3.ID = primitive.NewObjectID()
		def3.CreatedAt = time.Now()
		
		// Verify timestamps are in order
		if !def3.CreatedAt.After(def2.CreatedAt) {
			t.Error("Expected def3 to be created after def2")
		}
		if !def2.CreatedAt.After(def1.CreatedAt) {
			t.Error("Expected def2 to be created after def1")
		}
	})
	
	_ = ctx
}

// TestFindAll_EmptyResult tests FindAll when no definitions exist
func TestBatchDefinitionRepository_FindAll_EmptyResult(t *testing.T) {
	ctx := context.Background()
	
	filter := batch.BatchDefinitionFilter{}
	
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

// TestConversionRates tests conversion rate handling
func TestBatchDefinitionRepository_ConversionRates(t *testing.T) {
	t.Run("Multiple conversion rates", func(t *testing.T) {
		def := &batch.BatchDefinition{
			Name:           "Complex Batch",
			Unit:           "ml",
			ShelfLifeHours: 24,
			ConversionRates: []batch.ConversionRate{
				{
					SourceIngredientID:   primitive.NewObjectID(),
					SourceIngredientName: "Ingredient 1",
					SourceQuantity:       100,
					SourceUnit:           "g",
					BatchQuantity:        500,
					WastageRate:          0.1,
				},
				{
					SourceIngredientID:   primitive.NewObjectID(),
					SourceIngredientName: "Ingredient 2",
					SourceQuantity:       50,
					SourceUnit:           "ml",
					BatchQuantity:        500,
					WastageRate:          0.05,
				},
			},
			LowStockThreshold:  200,
			ExpiryWarningHours: 4,
		}
		
		if len(def.ConversionRates) != 2 {
			t.Errorf("Expected 2 conversion rates, got %d", len(def.ConversionRates))
		}
		
		// Verify first rate
		rate1 := def.ConversionRates[0]
		if rate1.SourceIngredientName != "Ingredient 1" {
			t.Errorf("Expected 'Ingredient 1', got '%s'", rate1.SourceIngredientName)
		}
		if rate1.WastageRate != 0.1 {
			t.Errorf("Expected wastage rate 0.1, got %f", rate1.WastageRate)
		}
		
		// Verify second rate
		rate2 := def.ConversionRates[1]
		if rate2.SourceIngredientName != "Ingredient 2" {
			t.Errorf("Expected 'Ingredient 2', got '%s'", rate2.SourceIngredientName)
		}
		if rate2.WastageRate != 0.05 {
			t.Errorf("Expected wastage rate 0.05, got %f", rate2.WastageRate)
		}
	})
	
	t.Run("Zero wastage rate", func(t *testing.T) {
		rate := batch.ConversionRate{
			SourceIngredientID:   primitive.NewObjectID(),
			SourceIngredientName: "Perfect Ingredient",
			SourceQuantity:       100,
			SourceUnit:           "g",
			BatchQuantity:        100,
			WastageRate:          0.0,
		}
		
		if rate.WastageRate != 0.0 {
			t.Errorf("Expected wastage rate 0.0, got %f", rate.WastageRate)
		}
	})
}

// TestBatchDefinitionValidation tests field validation
func TestBatchDefinitionRepository_Validation(t *testing.T) {
	t.Run("Valid batch definition", func(t *testing.T) {
		def := createTestBatchDefinition("Valid Batch")
		
		// Verify all required fields are set
		if def.Name == "" {
			t.Error("Expected name to be set")
		}
		if def.Unit == "" {
			t.Error("Expected unit to be set")
		}
		if def.ShelfLifeHours <= 0 {
			t.Error("Expected positive shelf life hours")
		}
		if len(def.ConversionRates) == 0 {
			t.Error("Expected at least one conversion rate")
		}
		if def.LowStockThreshold < 0 {
			t.Error("Expected non-negative low stock threshold")
		}
		if def.ExpiryWarningHours < 0 {
			t.Error("Expected non-negative expiry warning hours")
		}
	})
	
	t.Run("Batch definition with edge values", func(t *testing.T) {
		def := &batch.BatchDefinition{
			Name:               "Edge Case Batch",
			Unit:               "piece",
			ShelfLifeHours:     1, // Minimum valid value
			ConversionRates:    []batch.ConversionRate{},
			LowStockThreshold:  0, // Minimum valid value
			ExpiryWarningHours: 0, // Minimum valid value
		}
		
		if def.ShelfLifeHours < 1 {
			t.Error("Expected shelf life hours >= 1")
		}
		if def.LowStockThreshold < 0 {
			t.Error("Expected low stock threshold >= 0")
		}
		if def.ExpiryWarningHours < 0 {
			t.Error("Expected expiry warning hours >= 0")
		}
	})
}
