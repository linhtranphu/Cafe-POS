package services

import (
	"context"
	"fmt"
	"cafe-pos/backend/domain/menu"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type MenuRepository interface {
	Create(ctx context.Context, item *menu.MenuItem) error
	FindAll(ctx context.Context) ([]*menu.MenuItem, error)
	FindByID(ctx context.Context, id primitive.ObjectID) (*menu.MenuItem, error)
	FindByCategory(ctx context.Context, category string) ([]*menu.MenuItem, error)
	FindByIngredientName(ctx context.Context, ingredientName string) ([]*menu.MenuItem, error)
	FindByBatchDefinitionID(ctx context.Context, batchDefID primitive.ObjectID) ([]*menu.MenuItem, error)
	Update(ctx context.Context, id primitive.ObjectID, item *menu.MenuItem) error
	Delete(ctx context.Context, id primitive.ObjectID) error
}

// MenuBatchDefinitionRepository interface for checking batch availability
type MenuBatchDefinitionRepository interface {
	FindByID(ctx context.Context, id primitive.ObjectID) (interface{}, error)
}

// MenuBatchRecordRepository interface for checking batch availability
type MenuBatchRecordRepository interface {
	GetTotalAvailableQuantity(ctx context.Context, defID primitive.ObjectID) (float64, error)
}

type MenuService struct {
	menuRepo      MenuRepository
	batchDefRepo  MenuBatchDefinitionRepository
	batchRecRepo  MenuBatchRecordRepository
}

func NewMenuService(menuRepo MenuRepository) *MenuService {
	return &MenuService{
		menuRepo:     menuRepo,
		batchDefRepo: nil,
		batchRecRepo: nil,
	}
}

// SetBatchRepositories sets the batch repositories for validation (optional)
// This allows the service to validate batch availability when creating/updating menu items
func (s *MenuService) SetBatchRepositories(batchDefRepo MenuBatchDefinitionRepository, batchRecRepo MenuBatchRecordRepository) {
	s.batchDefRepo = batchDefRepo
	s.batchRecRepo = batchRecRepo
}

func (s *MenuService) CreateMenuItem(ctx context.Context, req *menu.CreateMenuItemRequest) (*menu.MenuItem, error) {
	item := &menu.MenuItem{
		Name:        req.Name,
		Category:    req.Category,
		Description: req.Description,
		Available:   true,
		HasVariants: req.HasVariants,
	}

	// Handle single-size vs multi-size items
	if req.HasVariants {
		// Multi-size item - use variants only
		item.Variants = req.Variants
		// Do NOT copy price and ingredients - they should be empty for multi-size
		item.Price = 0
		item.Ingredients = nil
		
		// Validate batch ingredients in variants
		for _, variant := range req.Variants {
			if err := s.validateBatchIngredients(ctx, variant.Ingredients); err != nil {
				return nil, fmt.Errorf("variant %s: %w", variant.ID, err)
			}
		}
	} else {
		// Single-size item - use price and ingredients (backward compatible)
		item.Price = req.Price
		item.Ingredients = req.Ingredients
		// Do NOT copy variants - they should be empty for single-size
		item.Variants = nil
		
		// Validate batch ingredients
		if err := s.validateBatchIngredients(ctx, req.Ingredients); err != nil {
			return nil, err
		}
	}

	// Validate the item before saving
	// This will catch ambiguous states (e.g., has_variants=true but price > 0)
	if err := item.Validate(); err != nil {
		return nil, fmt.Errorf("validation failed: %w", err)
	}

	err := s.menuRepo.Create(ctx, item)
	if err != nil {
		return nil, err
	}

	return item, nil
}

func (s *MenuService) GetAllMenuItems(ctx context.Context) ([]*menu.MenuItem, error) {
	return s.menuRepo.FindAll(ctx)
}

func (s *MenuService) GetMenuItem(ctx context.Context, id primitive.ObjectID) (*menu.MenuItem, error) {
	return s.menuRepo.FindByID(ctx, id)
}

func (s *MenuService) UpdateMenuItem(ctx context.Context, id primitive.ObjectID, req *menu.UpdateMenuItemRequest) (*menu.MenuItem, error) {
	item, err := s.menuRepo.FindByID(ctx, id)
	if err != nil {
		return nil, err
	}

	// Update basic fields
	if req.Name != "" {
		item.Name = req.Name
	}
	if req.Category != "" {
		item.Category = req.Category
	}
	if req.Description != "" {
		item.Description = req.Description
	}
	if req.Available != nil {
		item.Available = *req.Available
	}

	// Handle has_variants toggle
	if req.HasVariants != nil {
		oldHasVariants := item.HasVariants
		newHasVariants := *req.HasVariants

		if oldHasVariants != newHasVariants {
			// Toggling between single-size and multi-size
			if newHasVariants {
				// Changing from single-size to multi-size
				// Clear old single-size fields
				item.Price = 0
				item.Ingredients = nil
				item.CurrentCost = 0
				item.CostStatus = ""
				// Set variants from request
				item.Variants = req.Variants
				
				// Validate batch ingredients in variants
				for _, variant := range req.Variants {
					if err := s.validateBatchIngredients(ctx, variant.Ingredients); err != nil {
						return nil, fmt.Errorf("variant %s: %w", variant.ID, err)
					}
				}
			} else {
				// Changing from multi-size to single-size
				// Clear variants
				item.Variants = nil
				// Set single-size fields from request
				if req.Price > 0 {
					item.Price = req.Price
				}
				// Always update ingredients if provided (even if empty)
				if req.Ingredients != nil {
					item.Ingredients = req.Ingredients
					
					// Validate batch ingredients (only if not empty)
					if len(req.Ingredients) > 0 {
						if err := s.validateBatchIngredients(ctx, req.Ingredients); err != nil {
							return nil, err
						}
					}
				}
			}
			item.HasVariants = newHasVariants
		} else {
			// Not toggling, just updating within same mode
			if item.HasVariants {
				// Multi-size: update variants if provided
				if len(req.Variants) > 0 {
					item.Variants = req.Variants
					// Ensure root-level ingredients are cleared for multi-size items
					item.Ingredients = nil
					item.Price = 0

					// Validate batch ingredients in variants
					for _, variant := range req.Variants {
						if err := s.validateBatchIngredients(ctx, variant.Ingredients); err != nil {
							return nil, fmt.Errorf("variant %s: %w", variant.ID, err)
						}
					}
				}
			} else {
				// Single-size: update price/ingredients if provided
				if req.Price > 0 {
					item.Price = req.Price
				}
				// Always update ingredients if provided in request (even if empty array)
				if req.Ingredients != nil {
					item.Ingredients = req.Ingredients
					
					// Validate batch ingredients (only if not empty)
					if len(req.Ingredients) > 0 {
						if err := s.validateBatchIngredients(ctx, req.Ingredients); err != nil {
							return nil, err
						}
					}
				}
			}
		}
	} else {
		// has_variants not specified in request, update fields based on current mode
		if item.HasVariants {
			// Multi-size: update variants if provided
			if len(req.Variants) > 0 {
				item.Variants = req.Variants
				// Ensure root-level ingredients are cleared for multi-size items
				item.Ingredients = nil
				item.Price = 0

				// Validate batch ingredients in variants
				for _, variant := range req.Variants {
					if err := s.validateBatchIngredients(ctx, variant.Ingredients); err != nil {
						return nil, fmt.Errorf("variant %s: %w", variant.ID, err)
					}
				}
			}
		} else {
			// Single-size: update price/ingredients if provided
			if req.Price > 0 {
				item.Price = req.Price
			}
			// Always update ingredients if provided in request (even if empty array)
			if req.Ingredients != nil {
				item.Ingredients = req.Ingredients
				
				// Validate batch ingredients (only if not empty)
				if len(req.Ingredients) > 0 {
					if err := s.validateBatchIngredients(ctx, req.Ingredients); err != nil {
						return nil, err
					}
				}
			}
		}
	}

	// Validate before saving
	if err := item.Validate(); err != nil {
		return nil, fmt.Errorf("validation failed: %w", err)
	}

	err = s.menuRepo.Update(ctx, id, item)
	if err != nil {
		return nil, err
	}

	return item, nil
}

func (s *MenuService) DeleteMenuItem(ctx context.Context, id primitive.ObjectID) error {
	return s.menuRepo.Delete(ctx, id)
}

// validateBatchIngredients validates that batch ingredients exist and have sufficient quantity
func (s *MenuService) validateBatchIngredients(ctx context.Context, ingredients []menu.Ingredient) error {
	// If batch repositories are not set, skip validation
	if s.batchDefRepo == nil || s.batchRecRepo == nil {
		return nil
	}
	
	for _, ing := range ingredients {
		// Only validate batch ingredients
		if !ing.IsBatchIngredient() {
			continue
		}
		
		// Check if BatchID is set
		if ing.BatchID == nil {
			return fmt.Errorf("batch ingredient '%s' missing batch_id", ing.Name)
		}
		
		// Check if batch definition exists
		_, err := s.batchDefRepo.FindByID(ctx, *ing.BatchID)
		if err != nil {
			return fmt.Errorf("batch definition not found for ingredient '%s': %w", ing.Name, err)
		}
		
		// Check available quantity (warning only, not blocking)
		availableQty, err := s.batchRecRepo.GetTotalAvailableQuantity(ctx, *ing.BatchID)
		if err != nil {
			// Log warning but don't block creation
			fmt.Printf("Warning: failed to check batch availability for '%s': %v\n", ing.Name, err)
			continue
		}
		
		// If no batches available, return warning (not error, as batches can be prepared later)
		if availableQty == 0 {
			fmt.Printf("Warning: no available batches for ingredient '%s'\n", ing.Name)
		}
	}
	
	return nil
}
