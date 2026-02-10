package services

import (
	"context"
	"fmt"

	"cafe-pos/backend/domain/menu"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

// MenuCategoryRepository interface
type MenuCategoryRepository interface {
	Create(ctx context.Context, category *menu.MenuCategory) error
	FindAll(ctx context.Context) ([]*menu.MenuCategory, error)
	FindByID(ctx context.Context, id primitive.ObjectID) (*menu.MenuCategory, error)
	FindByName(ctx context.Context, name string) (*menu.MenuCategory, error)
	Update(ctx context.Context, id primitive.ObjectID, category *menu.MenuCategory) error
	Delete(ctx context.Context, id primitive.ObjectID) error
}

// MenuCategoryService handles menu category business logic
type MenuCategoryService struct {
	categoryRepo MenuCategoryRepository
	menuRepo     MenuRepository
}

// NewMenuCategoryService creates a new menu category service
func NewMenuCategoryService(categoryRepo MenuCategoryRepository, menuRepo MenuRepository) *MenuCategoryService {
	return &MenuCategoryService{
		categoryRepo: categoryRepo,
		menuRepo:     menuRepo,
	}
}

// CreateCategory creates a new menu category
func (s *MenuCategoryService) CreateCategory(ctx context.Context, req *menu.CreateMenuCategoryRequest) (*menu.MenuCategory, error) {
	// Check if category already exists
	existing, err := s.categoryRepo.FindByName(ctx, req.Name)
	if err == nil && existing != nil {
		return nil, fmt.Errorf("category with name '%s' already exists", req.Name)
	}

	category := &menu.MenuCategory{
		Name: req.Name,
	}

	err = s.categoryRepo.Create(ctx, category)
	if err != nil {
		return nil, fmt.Errorf("failed to create category: %w", err)
	}

	return category, nil
}

// GetAllCategories returns all menu categories
func (s *MenuCategoryService) GetAllCategories(ctx context.Context) ([]*menu.MenuCategory, error) {
	categories, err := s.categoryRepo.FindAll(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch categories: %w", err)
	}
	return categories, nil
}

// GetCategory returns a single category by ID
func (s *MenuCategoryService) GetCategory(ctx context.Context, id primitive.ObjectID) (*menu.MenuCategory, error) {
	category, err := s.categoryRepo.FindByID(ctx, id)
	if err != nil {
		return nil, fmt.Errorf("category not found: %w", err)
	}
	return category, nil
}

// UpdateCategory updates a menu category
func (s *MenuCategoryService) UpdateCategory(ctx context.Context, id primitive.ObjectID, req *menu.UpdateMenuCategoryRequest) (*menu.MenuCategory, error) {
	// Check if category exists
	category, err := s.categoryRepo.FindByID(ctx, id)
	if err != nil {
		return nil, fmt.Errorf("category not found: %w", err)
	}

	// Check if new name conflicts with existing category
	if req.Name != category.Name {
		existing, err := s.categoryRepo.FindByName(ctx, req.Name)
		if err == nil && existing != nil && existing.ID != id {
			return nil, fmt.Errorf("category with name '%s' already exists", req.Name)
		}
	}

	category.Name = req.Name

	err = s.categoryRepo.Update(ctx, id, category)
	if err != nil {
		return nil, fmt.Errorf("failed to update category: %w", err)
	}

	return category, nil
}

// DeleteCategory deletes a menu category
func (s *MenuCategoryService) DeleteCategory(ctx context.Context, id primitive.ObjectID) error {
	// Check if category exists
	category, err := s.categoryRepo.FindByID(ctx, id)
	if err != nil {
		return fmt.Errorf("category not found: %w", err)
	}

	// Check if any menu items use this category
	menuItems, err := s.menuRepo.FindByCategory(ctx, category.Name)
	if err != nil {
		return fmt.Errorf("failed to check menu items: %w", err)
	}

	if len(menuItems) > 0 {
		return fmt.Errorf("cannot delete category '%s' because it has %d menu items", category.Name, len(menuItems))
	}

	err = s.categoryRepo.Delete(ctx, id)
	if err != nil {
		return fmt.Errorf("failed to delete category: %w", err)
	}

	return nil
}
