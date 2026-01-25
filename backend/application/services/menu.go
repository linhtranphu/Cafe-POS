package services

import (
	"context"
	"cafe-pos/backend/domain/menu"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type MenuRepository interface {
	Create(ctx context.Context, item *menu.MenuItem) error
	FindAll(ctx context.Context) ([]*menu.MenuItem, error)
	FindByID(ctx context.Context, id primitive.ObjectID) (*menu.MenuItem, error)
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
		Price:       req.Price,
		Category:    req.Category,
		Description: req.Description,
		Ingredients: req.Ingredients,
		Available:   true,
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

	if req.Name != "" {
		item.Name = req.Name
	}
	if req.Price > 0 {
		item.Price = req.Price
	}
	if req.Category != "" {
		item.Category = req.Category
	}
	if req.Description != "" {
		item.Description = req.Description
	}
	if len(req.Ingredients) > 0 {
		item.Ingredients = req.Ingredients
	}
	if req.Available != nil {
		item.Available = *req.Available
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