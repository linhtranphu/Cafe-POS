package services

import (
	"context"
	"errors"
	"testing"
	"time"

	"cafe-pos/backend/domain/batch"
	"cafe-pos/backend/domain/ingredient"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

// Mock BatchDefinitionRepository for testing
type mockBatchDefinitionRepository struct {
	definitions map[primitive.ObjectID]*batch.BatchDefinition
	createErr   error
	updateErr   error
	deleteErr   error
	findErr     error
}

func newMockBatchDefinitionRepository() *mockBatchDefinitionRepository {
	return &mockBatchDefinitionRepository{
		definitions: make(map[primitive.ObjectID]*batch.BatchDefinition),
	}
}

func (m *mockBatchDefinitionRepository) Create(ctx context.Context, def *batch.BatchDefinition) error {
	if m.createErr != nil {
		return m.createErr
	}
	if def.ID.IsZero() {
		def.ID = primitive.NewObjectID()
	}
	m.definitions[def.ID] = def
	return nil
}

func (m *mockBatchDefinitionRepository) Update(ctx context.Context, def *batch.BatchDefinition) error {
	if m.updateErr != nil {
		return m.updateErr
	}
	if _, exists := m.definitions[def.ID]; !exists {
		return errors.New("definition not found")
	}
	m.definitions[def.ID] = def
	return nil
}

func (m *mockBatchDefinitionRepository) Delete(ctx context.Context, id primitive.ObjectID) error {
	if m.deleteErr != nil {
		return m.deleteErr
	}
	if _, exists := m.definitions[id]; !exists {
		return errors.New("definition not found")
	}
	delete(m.definitions, id)
	return nil
}

func (m *mockBatchDefinitionRepository) FindByID(ctx context.Context, id primitive.ObjectID) (*batch.BatchDefinition, error) {
	if m.findErr != nil {
		return nil, m.findErr
	}
	def, exists := m.definitions[id]
	if !exists {
		return nil, errors.New("definition not found")
	}
	return def, nil
}

func (m *mockBatchDefinitionRepository) FindAll(ctx context.Context, filter batch.BatchDefinitionFilter) ([]*batch.BatchDefinition, int64, error) {
	if m.findErr != nil {
		return nil, 0, m.findErr
	}
	
	var result []*batch.BatchDefinition
	for _, def := range m.definitions {
		// Simple search filter
		if filter.Search == "" || containsString(def.Name, filter.Search) {
			result = append(result, def)
		}
	}
	
	return result, int64(len(result)), nil
}

func containsString(s, substr string) bool {
	return len(s) >= len(substr) && (s == substr || len(substr) == 0)
}

// Mock IngredientRepository for testing
type mockIngredientRepositoryForBatch struct {
	ingredients map[primitive.ObjectID]*ingredient.Ingredient
	findErr     error
}

func newMockIngredientRepositoryForBatch() *mockIngredientRepositoryForBatch {
	return &mockIngredientRepositoryForBatch{
		ingredients: make(map[primitive.ObjectID]*ingredient.Ingredient),
	}
}

func (m *mockIngredientRepositoryForBatch) Create(ctx context.Context, item *ingredient.Ingredient) error {
	return nil
}

func (m *mockIngredientRepositoryForBatch) FindAll(ctx context.Context) ([]*ingredient.Ingredient, error) {
	return nil, nil
}

func (m *mockIngredientRepositoryForBatch) FindByID(ctx context.Context, id primitive.ObjectID) (*ingredient.Ingredient, error) {
	if m.findErr != nil {
		return nil, m.findErr
	}
	ing, exists := m.ingredients[id]
	if !exists {
		return nil, errors.New("ingredient not found")
	}
	return ing, nil
}

func (m *mockIngredientRepositoryForBatch) Update(ctx context.Context, id primitive.ObjectID, item *ingredient.Ingredient) error {
	return nil
}

func (m *mockIngredientRepositoryForBatch) Delete(ctx context.Context, id primitive.ObjectID) error {
	return nil
}

func (m *mockIngredientRepositoryForBatch) FindLowStock(ctx context.Context) ([]*ingredient.Ingredient, error) {
	return nil, nil
}

func (m *mockIngredientRepositoryForBatch) CreateCategory(ctx context.Context, cat *ingredient.IngredientCategory) error {
	return nil
}

func (m *mockIngredientRepositoryForBatch) GetCategories(ctx context.Context) ([]ingredient.IngredientCategory, error) {
	return nil, nil
}

func (m *mockIngredientRepositoryForBatch) DeleteCategory(ctx context.Context, id primitive.ObjectID) error {
	return nil
}

// Test Create method
func TestBatchDefinitionService_Create(t *testing.T) {
	ctx := context.Background()

	t.Run("Success - Create valid batch definition", func(t *testing.T) {
		batchRepo := newMockBatchDefinitionRepository()
		ingredientRepo := newMockIngredientRepositoryForBatch()
		service := NewBatchDefinitionService(batchRepo, ingredientRepo)

		// Add mock ingredient
		ingredientID := primitive.NewObjectID()
		ingredientRepo.ingredients[ingredientID] = &ingredient.Ingredient{
			ID:   ingredientID,
			Name: "Hạt Cà Phê",
			Unit: "g",
		}

		req := &batch.CreateBatchDefinitionRequest{
			Name:           "Cà Phê Concentrate",
			Unit:           "ml",
			ShelfLifeHours: 24,
			ConversionRates: []batch.ConversionRate{
				{
					SourceIngredientID:   ingredientID,
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

		def, err := service.Create(ctx, req)
		if err != nil {
			t.Fatalf("Expected no error, got %v", err)
		}

		if def.ID.IsZero() {
			t.Error("Expected definition ID to be set")
		}
		if def.Name != req.Name {
			t.Errorf("Expected name %s, got %s", req.Name, def.Name)
		}
		if def.Unit != req.Unit {
			t.Errorf("Expected unit %s, got %s", req.Unit, def.Unit)
		}
		if def.ShelfLifeHours != req.ShelfLifeHours {
			t.Errorf("Expected shelf life %d, got %d", req.ShelfLifeHours, def.ShelfLifeHours)
		}
		if len(def.ConversionRates) != 1 {
			t.Errorf("Expected 1 conversion rate, got %d", len(def.ConversionRates))
		}
		if def.CreatedAt.IsZero() {
			t.Error("Expected CreatedAt to be set")
		}
		if def.UpdatedAt.IsZero() {
			t.Error("Expected UpdatedAt to be set")
		}
	})

	t.Run("Error - Invalid conversion rates (empty)", func(t *testing.T) {
		batchRepo := newMockBatchDefinitionRepository()
		ingredientRepo := newMockIngredientRepositoryForBatch()
		service := NewBatchDefinitionService(batchRepo, ingredientRepo)

		req := &batch.CreateBatchDefinitionRequest{
			Name:              "Test Batch",
			Unit:              "ml",
			ShelfLifeHours:    24,
			ConversionRates:   []batch.ConversionRate{}, // Empty
			LowStockThreshold: 100,
		}

		_, err := service.Create(ctx, req)
		if err == nil {
			t.Error("Expected error for empty conversion rates")
		}
		if !errors.Is(err, ErrInvalidConversionRates) {
			t.Errorf("Expected ErrInvalidConversionRates, got %v", err)
		}
	})

	t.Run("Error - Ingredient not found", func(t *testing.T) {
		batchRepo := newMockBatchDefinitionRepository()
		ingredientRepo := newMockIngredientRepositoryForBatch()
		service := NewBatchDefinitionService(batchRepo, ingredientRepo)

		nonExistentID := primitive.NewObjectID()
		req := &batch.CreateBatchDefinitionRequest{
			Name:           "Test Batch",
			Unit:           "ml",
			ShelfLifeHours: 24,
			ConversionRates: []batch.ConversionRate{
				{
					SourceIngredientID: nonExistentID,
					SourceQuantity:     100,
					BatchQuantity:      500,
					WastageRate:        0.1,
				},
			},
		}

		_, err := service.Create(ctx, req)
		if err == nil {
			t.Error("Expected error for non-existent ingredient")
		}
		if !errors.Is(err, ErrIngredientNotFound) {
			t.Errorf("Expected ErrIngredientNotFound, got %v", err)
		}
	})

	t.Run("Error - Invalid source quantity", func(t *testing.T) {
		batchRepo := newMockBatchDefinitionRepository()
		ingredientRepo := newMockIngredientRepositoryForBatch()
		service := NewBatchDefinitionService(batchRepo, ingredientRepo)

		ingredientID := primitive.NewObjectID()
		ingredientRepo.ingredients[ingredientID] = &ingredient.Ingredient{
			ID:   ingredientID,
			Name: "Test Ingredient",
		}

		req := &batch.CreateBatchDefinitionRequest{
			Name:           "Test Batch",
			Unit:           "ml",
			ShelfLifeHours: 24,
			ConversionRates: []batch.ConversionRate{
				{
					SourceIngredientID: ingredientID,
					SourceQuantity:     0, // Invalid
					BatchQuantity:      500,
					WastageRate:        0.1,
				},
			},
		}

		_, err := service.Create(ctx, req)
		if err == nil {
			t.Error("Expected error for invalid source quantity")
		}
	})

	t.Run("Error - Invalid wastage rate", func(t *testing.T) {
		batchRepo := newMockBatchDefinitionRepository()
		ingredientRepo := newMockIngredientRepositoryForBatch()
		service := NewBatchDefinitionService(batchRepo, ingredientRepo)

		ingredientID := primitive.NewObjectID()
		ingredientRepo.ingredients[ingredientID] = &ingredient.Ingredient{
			ID:   ingredientID,
			Name: "Test Ingredient",
		}

		req := &batch.CreateBatchDefinitionRequest{
			Name:           "Test Batch",
			Unit:           "ml",
			ShelfLifeHours: 24,
			ConversionRates: []batch.ConversionRate{
				{
					SourceIngredientID: ingredientID,
					SourceQuantity:     100,
					BatchQuantity:      500,
					WastageRate:        1.5, // Invalid (> 1)
				},
			},
		}

		_, err := service.Create(ctx, req)
		if err == nil {
			t.Error("Expected error for invalid wastage rate")
		}
	})
}

// Test Update method
func TestBatchDefinitionService_Update(t *testing.T) {
	ctx := context.Background()

	t.Run("Success - Update batch definition", func(t *testing.T) {
		batchRepo := newMockBatchDefinitionRepository()
		ingredientRepo := newMockIngredientRepositoryForBatch()
		service := NewBatchDefinitionService(batchRepo, ingredientRepo)

		// Create initial definition
		defID := primitive.NewObjectID()
		initialDef := &batch.BatchDefinition{
			ID:                 defID,
			Name:               "Original Name",
			Unit:               "ml",
			ShelfLifeHours:     24,
			ConversionRates:    []batch.ConversionRate{},
			LowStockThreshold:  100,
			ExpiryWarningHours: 4,
			CreatedAt:          time.Now(),
			UpdatedAt:          time.Now(),
		}
		batchRepo.definitions[defID] = initialDef

		// Update
		newShelfLife := 48
		newThreshold := 200.0
		req := &batch.UpdateBatchDefinitionRequest{
			Name:              "Updated Name",
			ShelfLifeHours:    &newShelfLife,
			LowStockThreshold: &newThreshold,
		}

		updated, err := service.Update(ctx, defID, req)
		if err != nil {
			t.Fatalf("Expected no error, got %v", err)
		}

		if updated.Name != "Updated Name" {
			t.Errorf("Expected name 'Updated Name', got %s", updated.Name)
		}
		if updated.ShelfLifeHours != 48 {
			t.Errorf("Expected shelf life 48, got %d", updated.ShelfLifeHours)
		}
		if updated.LowStockThreshold != 200.0 {
			t.Errorf("Expected threshold 200.0, got %f", updated.LowStockThreshold)
		}
		if updated.Unit != "ml" {
			t.Error("Expected unit to remain unchanged")
		}
	})

	t.Run("Error - Definition not found", func(t *testing.T) {
		batchRepo := newMockBatchDefinitionRepository()
		ingredientRepo := newMockIngredientRepositoryForBatch()
		service := NewBatchDefinitionService(batchRepo, ingredientRepo)

		nonExistentID := primitive.NewObjectID()
		req := &batch.UpdateBatchDefinitionRequest{
			Name: "Updated Name",
		}

		_, err := service.Update(ctx, nonExistentID, req)
		if err == nil {
			t.Error("Expected error for non-existent definition")
		}
		if !errors.Is(err, ErrBatchDefinitionNotFound) {
			t.Errorf("Expected ErrBatchDefinitionNotFound, got %v", err)
		}
	})

	t.Run("Success - Update conversion rates", func(t *testing.T) {
		batchRepo := newMockBatchDefinitionRepository()
		ingredientRepo := newMockIngredientRepositoryForBatch()
		service := NewBatchDefinitionService(batchRepo, ingredientRepo)

		// Add mock ingredient
		ingredientID := primitive.NewObjectID()
		ingredientRepo.ingredients[ingredientID] = &ingredient.Ingredient{
			ID:   ingredientID,
			Name: "New Ingredient",
		}

		// Create initial definition
		defID := primitive.NewObjectID()
		initialDef := &batch.BatchDefinition{
			ID:              defID,
			Name:            "Test Batch",
			Unit:            "ml",
			ShelfLifeHours:  24,
			ConversionRates: []batch.ConversionRate{},
			CreatedAt:       time.Now(),
			UpdatedAt:       time.Now(),
		}
		batchRepo.definitions[defID] = initialDef

		// Update with new conversion rates
		req := &batch.UpdateBatchDefinitionRequest{
			ConversionRates: []batch.ConversionRate{
				{
					SourceIngredientID: ingredientID,
					SourceQuantity:     200,
					BatchQuantity:      1000,
					WastageRate:        0.05,
				},
			},
		}

		updated, err := service.Update(ctx, defID, req)
		if err != nil {
			t.Fatalf("Expected no error, got %v", err)
		}

		if len(updated.ConversionRates) != 1 {
			t.Errorf("Expected 1 conversion rate, got %d", len(updated.ConversionRates))
		}
	})
}

// Test Delete method
func TestBatchDefinitionService_Delete(t *testing.T) {
	ctx := context.Background()

	t.Run("Success - Delete batch definition", func(t *testing.T) {
		batchRepo := newMockBatchDefinitionRepository()
		ingredientRepo := newMockIngredientRepositoryForBatch()
		service := NewBatchDefinitionService(batchRepo, ingredientRepo)

		// Create definition
		defID := primitive.NewObjectID()
		def := &batch.BatchDefinition{
			ID:   defID,
			Name: "Test Batch",
		}
		batchRepo.definitions[defID] = def

		// Delete
		err := service.Delete(ctx, defID)
		if err != nil {
			t.Fatalf("Expected no error, got %v", err)
		}

		// Verify deleted
		_, exists := batchRepo.definitions[defID]
		if exists {
			t.Error("Expected definition to be deleted")
		}
	})

	t.Run("Error - Definition not found", func(t *testing.T) {
		batchRepo := newMockBatchDefinitionRepository()
		ingredientRepo := newMockIngredientRepositoryForBatch()
		service := NewBatchDefinitionService(batchRepo, ingredientRepo)

		nonExistentID := primitive.NewObjectID()
		err := service.Delete(ctx, nonExistentID)
		if err == nil {
			t.Error("Expected error for non-existent definition")
		}
		if !errors.Is(err, ErrBatchDefinitionNotFound) {
			t.Errorf("Expected ErrBatchDefinitionNotFound, got %v", err)
		}
	})
}

// Test GetByID method
func TestBatchDefinitionService_GetByID(t *testing.T) {
	ctx := context.Background()

	t.Run("Success - Get batch definition by ID", func(t *testing.T) {
		batchRepo := newMockBatchDefinitionRepository()
		ingredientRepo := newMockIngredientRepositoryForBatch()
		service := NewBatchDefinitionService(batchRepo, ingredientRepo)

		// Create definition
		defID := primitive.NewObjectID()
		def := &batch.BatchDefinition{
			ID:   defID,
			Name: "Test Batch",
			Unit: "ml",
		}
		batchRepo.definitions[defID] = def

		// Get by ID
		result, err := service.GetByID(ctx, defID)
		if err != nil {
			t.Fatalf("Expected no error, got %v", err)
		}

		if result.ID != defID {
			t.Errorf("Expected ID %s, got %s", defID.Hex(), result.ID.Hex())
		}
		if result.Name != "Test Batch" {
			t.Errorf("Expected name 'Test Batch', got %s", result.Name)
		}
	})

	t.Run("Error - Definition not found", func(t *testing.T) {
		batchRepo := newMockBatchDefinitionRepository()
		ingredientRepo := newMockIngredientRepositoryForBatch()
		service := NewBatchDefinitionService(batchRepo, ingredientRepo)

		nonExistentID := primitive.NewObjectID()
		_, err := service.GetByID(ctx, nonExistentID)
		if err == nil {
			t.Error("Expected error for non-existent definition")
		}
		if !errors.Is(err, ErrBatchDefinitionNotFound) {
			t.Errorf("Expected ErrBatchDefinitionNotFound, got %v", err)
		}
	})
}

// Test List method
func TestBatchDefinitionService_List(t *testing.T) {
	ctx := context.Background()

	t.Run("Success - List all batch definitions", func(t *testing.T) {
		batchRepo := newMockBatchDefinitionRepository()
		ingredientRepo := newMockIngredientRepositoryForBatch()
		service := NewBatchDefinitionService(batchRepo, ingredientRepo)

		// Create multiple definitions
		def1 := &batch.BatchDefinition{
			ID:   primitive.NewObjectID(),
			Name: "Batch 1",
		}
		def2 := &batch.BatchDefinition{
			ID:   primitive.NewObjectID(),
			Name: "Batch 2",
		}
		batchRepo.definitions[def1.ID] = def1
		batchRepo.definitions[def2.ID] = def2

		// List
		filter := batch.BatchDefinitionFilter{}
		results, total, err := service.List(ctx, filter)
		if err != nil {
			t.Fatalf("Expected no error, got %v", err)
		}

		if total != 2 {
			t.Errorf("Expected total 2, got %d", total)
		}
		if len(results) != 2 {
			t.Errorf("Expected 2 results, got %d", len(results))
		}
	})

	t.Run("Success - List with search filter", func(t *testing.T) {
		batchRepo := newMockBatchDefinitionRepository()
		ingredientRepo := newMockIngredientRepositoryForBatch()
		service := NewBatchDefinitionService(batchRepo, ingredientRepo)

		// Create definitions
		def1 := &batch.BatchDefinition{
			ID:   primitive.NewObjectID(),
			Name: "Coffee Batch",
		}
		def2 := &batch.BatchDefinition{
			ID:   primitive.NewObjectID(),
			Name: "Tea Batch",
		}
		batchRepo.definitions[def1.ID] = def1
		batchRepo.definitions[def2.ID] = def2

		// List with filter
		filter := batch.BatchDefinitionFilter{
			Search: "Coffee Batch",
		}
		results, total, err := service.List(ctx, filter)
		if err != nil {
			t.Fatalf("Expected no error, got %v", err)
		}

		if total != 1 {
			t.Errorf("Expected total 1, got %d", total)
		}
		if len(results) != 1 {
			t.Errorf("Expected 1 result, got %d", len(results))
		}
		if len(results) > 0 && results[0].Name != "Coffee Batch" {
			t.Errorf("Expected 'Coffee Batch', got %s", results[0].Name)
		}
	})
}

// Test ValidateConversionRates method
func TestBatchDefinitionService_ValidateConversionRates(t *testing.T) {
	ctx := context.Background()

	t.Run("Success - Valid conversion rates", func(t *testing.T) {
		batchRepo := newMockBatchDefinitionRepository()
		ingredientRepo := newMockIngredientRepositoryForBatch()
		service := NewBatchDefinitionService(batchRepo, ingredientRepo)

		// Add mock ingredient
		ingredientID := primitive.NewObjectID()
		ingredientRepo.ingredients[ingredientID] = &ingredient.Ingredient{
			ID:   ingredientID,
			Name: "Test Ingredient",
		}

		rates := []batch.ConversionRate{
			{
				SourceIngredientID: ingredientID,
				SourceQuantity:     100,
				BatchQuantity:      500,
				WastageRate:        0.1,
			},
		}

		err := service.ValidateConversionRates(ctx, rates)
		if err != nil {
			t.Errorf("Expected no error, got %v", err)
		}
	})

	t.Run("Error - Empty conversion rates", func(t *testing.T) {
		batchRepo := newMockBatchDefinitionRepository()
		ingredientRepo := newMockIngredientRepositoryForBatch()
		service := NewBatchDefinitionService(batchRepo, ingredientRepo)

		rates := []batch.ConversionRate{}
		err := service.ValidateConversionRates(ctx, rates)
		if err == nil {
			t.Error("Expected error for empty conversion rates")
		}
	})

	t.Run("Error - Negative source quantity", func(t *testing.T) {
		batchRepo := newMockBatchDefinitionRepository()
		ingredientRepo := newMockIngredientRepositoryForBatch()
		service := NewBatchDefinitionService(batchRepo, ingredientRepo)

		ingredientID := primitive.NewObjectID()
		ingredientRepo.ingredients[ingredientID] = &ingredient.Ingredient{
			ID: ingredientID,
		}

		rates := []batch.ConversionRate{
			{
				SourceIngredientID: ingredientID,
				SourceQuantity:     -100, // Invalid
				BatchQuantity:      500,
				WastageRate:        0.1,
			},
		}

		err := service.ValidateConversionRates(ctx, rates)
		if err == nil {
			t.Error("Expected error for negative source quantity")
		}
	})

	t.Run("Error - Wastage rate out of range", func(t *testing.T) {
		batchRepo := newMockBatchDefinitionRepository()
		ingredientRepo := newMockIngredientRepositoryForBatch()
		service := NewBatchDefinitionService(batchRepo, ingredientRepo)

		ingredientID := primitive.NewObjectID()
		ingredientRepo.ingredients[ingredientID] = &ingredient.Ingredient{
			ID: ingredientID,
		}

		rates := []batch.ConversionRate{
			{
				SourceIngredientID: ingredientID,
				SourceQuantity:     100,
				BatchQuantity:      500,
				WastageRate:        -0.1, // Invalid
			},
		}

		err := service.ValidateConversionRates(ctx, rates)
		if err == nil {
			t.Error("Expected error for negative wastage rate")
		}
	})
}
