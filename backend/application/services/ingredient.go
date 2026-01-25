package services

import (
	"context"
	"cafe-pos/backend/domain/ingredient"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type IngredientRepository interface {
	Create(ctx context.Context, item *ingredient.Ingredient) error
	FindAll(ctx context.Context) ([]*ingredient.Ingredient, error)
	FindByID(ctx context.Context, id primitive.ObjectID) (*ingredient.Ingredient, error)
	Update(ctx context.Context, id primitive.ObjectID, item *ingredient.Ingredient) error
	Delete(ctx context.Context, id primitive.ObjectID) error
	FindLowStock(ctx context.Context) ([]*ingredient.Ingredient, error)
}

type StockHistoryRepository interface {
	Create(ctx context.Context, history *ingredient.StockHistory) error
	FindByIngredientID(ctx context.Context, ingredientID primitive.ObjectID) ([]*ingredient.StockHistory, error)
}

type IngredientService struct {
	ingredientRepo IngredientRepository
	stockHistoryRepo StockHistoryRepository
}

func NewIngredientService(ingredientRepo IngredientRepository, stockHistoryRepo StockHistoryRepository) *IngredientService {
	return &IngredientService{
		ingredientRepo: ingredientRepo,
		stockHistoryRepo: stockHistoryRepo,
	}
}

func (s *IngredientService) CreateIngredient(ctx context.Context, req *ingredient.CreateIngredientRequest) (*ingredient.Ingredient, error) {
	item := &ingredient.Ingredient{
		Name:        req.Name,
		Category:    req.Category,
		Unit:        req.Unit,
		Quantity:    req.Quantity,
		MinStock:    req.MinStock,
		CostPerUnit: req.CostPerUnit,
		Supplier:    req.Supplier,
	}

	err := s.ingredientRepo.Create(ctx, item)
	if err != nil {
		return nil, err
	}

	return item, nil
}

func (s *IngredientService) GetAllIngredients(ctx context.Context) ([]*ingredient.Ingredient, error) {
	return s.ingredientRepo.FindAll(ctx)
}

func (s *IngredientService) GetIngredient(ctx context.Context, id primitive.ObjectID) (*ingredient.Ingredient, error) {
	return s.ingredientRepo.FindByID(ctx, id)
}

func (s *IngredientService) UpdateIngredient(ctx context.Context, id primitive.ObjectID, req *ingredient.UpdateIngredientRequest) (*ingredient.Ingredient, error) {
	item, err := s.ingredientRepo.FindByID(ctx, id)
	if err != nil {
		return nil, err
	}

	if req.Name != "" {
		item.Name = req.Name
	}
	if req.Category != "" {
		item.Category = req.Category
	}
	if req.Unit != "" {
		item.Unit = req.Unit
	}
	if req.Quantity != nil {
		item.Quantity = *req.Quantity
	}
	if req.MinStock != nil {
		item.MinStock = *req.MinStock
	}
	if req.CostPerUnit != nil {
		item.CostPerUnit = *req.CostPerUnit
	}
	if req.Supplier != "" {
		item.Supplier = req.Supplier
	}

	err = s.ingredientRepo.Update(ctx, id, item)
	if err != nil {
		return nil, err
	}

	return item, nil
}

func (s *IngredientService) DeleteIngredient(ctx context.Context, id primitive.ObjectID) error {
	return s.ingredientRepo.Delete(ctx, id)
}

func (s *IngredientService) AdjustStock(ctx context.Context, id primitive.ObjectID, req *ingredient.StockAdjustmentRequest) (*ingredient.Ingredient, error) {
	item, err := s.ingredientRepo.FindByID(ctx, id)
	if err != nil {
		return nil, err
	}

	beforeQty := item.Quantity
	item.Quantity += req.Quantity
	if item.Quantity < 0 {
		item.Quantity = 0
	}
	afterQty := item.Quantity

	err = s.ingredientRepo.Update(ctx, id, item)
	if err != nil {
		return nil, err
	}

	// Create stock history record
	userID, _ := primitive.ObjectIDFromHex(req.UserID)
	history := &ingredient.StockHistory{
		IngredientID: id,
		Type:         ingredient.TransactionAdjustment,
		Quantity:     req.Quantity,
		BeforeQty:    beforeQty,
		AfterQty:     afterQty,
		Reason:       req.Reason,
		UserID:       userID,
		Username:     req.Username,
	}
	s.stockHistoryRepo.Create(ctx, history)

	return item, nil
}

func (s *IngredientService) GetLowStockIngredients(ctx context.Context) ([]*ingredient.Ingredient, error) {
	return s.ingredientRepo.FindLowStock(ctx)
}

func (s *IngredientService) GetStockHistory(ctx context.Context, id primitive.ObjectID) ([]*ingredient.StockHistory, error) {
	return s.stockHistoryRepo.FindByIngredientID(ctx, id)
}