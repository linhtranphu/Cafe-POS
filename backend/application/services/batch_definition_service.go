package services

import (
	"context"
	"errors"
	"fmt"
	"time"

	"cafe-pos/backend/domain/batch"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

var (
	ErrBatchDefinitionNotFound = errors.New("batch definition not found")
	ErrInvalidConversionRates  = errors.New("invalid conversion rates")
	ErrIngredientNotFound      = errors.New("source ingredient not found")
)

// BatchDefinitionService handles business logic for batch definitions
type BatchDefinitionService struct {
	batchDefRepo   batch.BatchDefinitionRepository
	ingredientRepo IngredientRepository
}

// NewBatchDefinitionService creates a new BatchDefinitionService
func NewBatchDefinitionService(
	batchDefRepo batch.BatchDefinitionRepository,
	ingredientRepo IngredientRepository,
) *BatchDefinitionService {
	return &BatchDefinitionService{
		batchDefRepo:   batchDefRepo,
		ingredientRepo: ingredientRepo,
	}
}

// Create creates a new batch definition with validation
// Validates: Requirements 1.1, 1.2, 1.3, 1.6
func (s *BatchDefinitionService) Create(ctx context.Context, req *batch.CreateBatchDefinitionRequest) (*batch.BatchDefinition, error) {
	// Validate conversion rates
	if err := s.ValidateConversionRates(ctx, req.ConversionRates); err != nil {
		return nil, err
	}

	// Create batch definition entity
	now := time.Now()
	def := &batch.BatchDefinition{
		Name:               req.Name,
		Unit:               req.Unit,
		ShelfLifeHours:     req.ShelfLifeHours,
		ConversionRates:    req.ConversionRates,
		LowStockThreshold:  req.LowStockThreshold,
		ExpiryWarningHours: req.ExpiryWarningHours,
		CreatedAt:          now,
		UpdatedAt:          now,
	}

	// Persist to repository
	if err := s.batchDefRepo.Create(ctx, def); err != nil {
		return nil, fmt.Errorf("failed to create batch definition: %w", err)
	}

	return def, nil
}

// Update updates an existing batch definition
// Validates: Requirements 1.1, 1.2, 1.3, 1.6
func (s *BatchDefinitionService) Update(ctx context.Context, id primitive.ObjectID, req *batch.UpdateBatchDefinitionRequest) (*batch.BatchDefinition, error) {
	// Find existing definition
	def, err := s.batchDefRepo.FindByID(ctx, id)
	if err != nil {
		return nil, ErrBatchDefinitionNotFound
	}

	// Update fields if provided
	if req.Name != "" {
		def.Name = req.Name
	}
	if req.Unit != "" {
		def.Unit = req.Unit
	}
	if req.ShelfLifeHours != nil {
		def.ShelfLifeHours = *req.ShelfLifeHours
	}
	if req.ConversionRates != nil {
		// Validate new conversion rates
		if err := s.ValidateConversionRates(ctx, req.ConversionRates); err != nil {
			return nil, err
		}
		def.ConversionRates = req.ConversionRates
	}
	if req.LowStockThreshold != nil {
		def.LowStockThreshold = *req.LowStockThreshold
	}
	if req.ExpiryWarningHours != nil {
		def.ExpiryWarningHours = *req.ExpiryWarningHours
	}

	def.UpdatedAt = time.Now()

	// Persist changes
	if err := s.batchDefRepo.Update(ctx, def); err != nil {
		return nil, fmt.Errorf("failed to update batch definition: %w", err)
	}

	return def, nil
}

// Delete deletes a batch definition
// Validates: Requirements 7.2
func (s *BatchDefinitionService) Delete(ctx context.Context, id primitive.ObjectID) error {
	// Check if definition exists
	_, err := s.batchDefRepo.FindByID(ctx, id)
	if err != nil {
		return ErrBatchDefinitionNotFound
	}

	// Delete the definition
	if err := s.batchDefRepo.Delete(ctx, id); err != nil {
		return fmt.Errorf("failed to delete batch definition: %w", err)
	}

	return nil
}

// GetByID retrieves a batch definition by ID
// Validates: Requirements 7.1
func (s *BatchDefinitionService) GetByID(ctx context.Context, id primitive.ObjectID) (*batch.BatchDefinition, error) {
	def, err := s.batchDefRepo.FindByID(ctx, id)
	if err != nil {
		return nil, ErrBatchDefinitionNotFound
	}
	return def, nil
}

// List retrieves batch definitions with optional filtering
// Validates: Requirements 7.1
func (s *BatchDefinitionService) List(ctx context.Context, filter batch.BatchDefinitionFilter) ([]*batch.BatchDefinition, int64, error) {
	defs, total, err := s.batchDefRepo.FindAll(ctx, filter)
	if err != nil {
		return nil, 0, fmt.Errorf("failed to list batch definitions: %w", err)
	}
	return defs, total, nil
}

// ValidateConversionRates validates that all source ingredients exist in the inventory system
// Validates: Requirements 1.6
func (s *BatchDefinitionService) ValidateConversionRates(ctx context.Context, rates []batch.ConversionRate) error {
	if len(rates) == 0 {
		return fmt.Errorf("%w: at least one conversion rate is required", ErrInvalidConversionRates)
	}

	for i, rate := range rates {
		// Validate quantities are positive
		if rate.SourceQuantity <= 0 {
			return fmt.Errorf("%w: source quantity must be positive at index %d", ErrInvalidConversionRates, i)
		}
		if rate.BatchQuantity <= 0 {
			return fmt.Errorf("%w: batch quantity must be positive at index %d", ErrInvalidConversionRates, i)
		}

		// Validate wastage rate is between 0 and 1
		if rate.WastageRate < 0 || rate.WastageRate > 1 {
			return fmt.Errorf("%w: wastage rate must be between 0 and 1 at index %d", ErrInvalidConversionRates, i)
		}

		// Check if source ingredient exists
		ingredient, err := s.ingredientRepo.FindByID(ctx, rate.SourceIngredientID)
		if err != nil {
			return fmt.Errorf("%w: ingredient with ID %s not found", ErrIngredientNotFound, rate.SourceIngredientID.Hex())
		}

		// Validate that the ingredient name matches (optional but good for consistency)
		if rate.SourceIngredientName != "" && rate.SourceIngredientName != ingredient.Name {
			return fmt.Errorf("%w: ingredient name mismatch at index %d (expected: %s, got: %s)", 
				ErrInvalidConversionRates, i, ingredient.Name, rate.SourceIngredientName)
		}
	}

	return nil
}
