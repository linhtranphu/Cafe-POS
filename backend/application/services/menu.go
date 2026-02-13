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
	Update(ctx context.Context, id primitive.ObjectID, item *menu.MenuItem) error
	Delete(ctx context.Context, id primitive.ObjectID) error
}

type MenuService struct {
	menuRepo MenuRepository
}

func NewMenuService(menuRepo MenuRepository) *MenuService {
	return &MenuService{menuRepo: menuRepo}
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
	} else {
		// Single-size item - use price and ingredients (backward compatible)
		item.Price = req.Price
		item.Ingredients = req.Ingredients
		// Do NOT copy variants - they should be empty for single-size
		item.Variants = nil
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
			} else {
				// Changing from multi-size to single-size
				// Clear variants
				item.Variants = nil
				// Set single-size fields from request
				if req.Price > 0 {
					item.Price = req.Price
				}
				if len(req.Ingredients) > 0 {
					item.Ingredients = req.Ingredients
				}
			}
			item.HasVariants = newHasVariants
		} else {
			// Not toggling, just updating within same mode
			if item.HasVariants {
				// Multi-size: update variants if provided
				if len(req.Variants) > 0 {
					item.Variants = req.Variants
				}
			} else {
				// Single-size: update price/ingredients if provided
				if req.Price > 0 {
					item.Price = req.Price
				}
				if len(req.Ingredients) > 0 {
					item.Ingredients = req.Ingredients
				}
			}
		}
	} else {
		// has_variants not specified in request, update fields based on current mode
		if item.HasVariants {
			// Multi-size: update variants if provided
			if len(req.Variants) > 0 {
				item.Variants = req.Variants
			}
		} else {
			// Single-size: update price/ingredients if provided
			if req.Price > 0 {
				item.Price = req.Price
			}
			if len(req.Ingredients) > 0 {
				item.Ingredients = req.Ingredients
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