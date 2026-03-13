package services

import (
	"context"
	"time"
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
	CreateCategory(ctx context.Context, cat *ingredient.IngredientCategory) error
	GetCategories(ctx context.Context) ([]ingredient.IngredientCategory, error)
	DeleteCategory(ctx context.Context, id primitive.ObjectID) error
}

type StockHistoryRepository interface {
	Create(ctx context.Context, history *ingredient.StockHistory) error
	FindByIngredientID(ctx context.Context, ingredientID primitive.ObjectID) ([]*ingredient.StockHistory, error)
}

type IngredientService struct {
	ingredientRepo      IngredientRepository
	stockHistoryRepo    StockHistoryRepository
	autoExpenseService  *AutoExpenseService
	costCalculatorService *CostCalculatorService
}

func NewIngredientService(ingredientRepo IngredientRepository, stockHistoryRepo StockHistoryRepository) *IngredientService {
	return &IngredientService{
		ingredientRepo:   ingredientRepo,
		stockHistoryRepo: stockHistoryRepo,
	}
}

// SetAutoExpenseService sets the AutoExpenseService for automatic expense tracking
// This is called after service initialization to avoid circular dependencies
func (s *IngredientService) SetAutoExpenseService(autoExpenseService *AutoExpenseService) {
	s.autoExpenseService = autoExpenseService
}

// SetCostCalculatorService sets the CostCalculatorService for triggering cost recalculation
// This is called after service initialization to avoid circular dependencies
// Requirements: 1.3, 9.1
func (s *IngredientService) SetCostCalculatorService(costCalculatorService *CostCalculatorService) {
	s.costCalculatorService = costCalculatorService
}

// triggerCostRecalculation queues cost recalculation for menu items using this ingredient
// This is called when ingredient cost_per_unit changes
// Requirements: 1.3, 9.1
func (s *IngredientService) triggerCostRecalculation(ctx context.Context, ingredientID primitive.ObjectID) {
	if s.costCalculatorService == nil {
		// Service not configured, skip recalculation
		return
	}
	
	// Queue recalculation asynchronously to avoid blocking the ingredient update
	go func() {
		// Create a background context with timeout to avoid hanging
		bgCtx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()
		
		// Queue cost recalculation for all menu items using this ingredient
		err := s.costCalculatorService.QueueCostRecalculation(bgCtx, ingredientID)
		if err != nil {
			// Log error but don't fail the ingredient update
			// In production, this should use proper logging
		}
	}()
}

func (s *IngredientService) CreateIngredient(ctx context.Context, req *ingredient.CreateIngredientRequest, userIDStr string, username string) (*ingredient.Ingredient, error) {
	// Set default values for conversion_rate and wastage_percentage
	conversionRate := 1.0
	if req.ConversionRate != nil {
		conversionRate = *req.ConversionRate
	}
	
	wastagePercentage := 0.0
	if req.WastagePercentage != nil {
		wastagePercentage = *req.WastagePercentage
	}
	
	item := &ingredient.Ingredient{
		Name:              req.Name,
		Category:          req.Category,
		Unit:              req.Unit,
		Quantity:          req.Quantity,
		MinStock:          req.MinStock,
		CostPerUnit:       req.CostPerUnit,
		Supplier:          req.Supplier,
		ConversionRate:    conversionRate,
		WastagePercentage: wastagePercentage,
	}
	
	// Set created_at from request if provided, otherwise use current time
	if req.CreatedDate != nil && *req.CreatedDate != "" {
		if parsedTime, err := time.Parse(time.RFC3339, *req.CreatedDate); err == nil {
			item.CreatedAt = parsedTime
		} else {
			// Try parsing as datetime-local format (YYYY-MM-DDTHH:MM)
			if parsedTime, err := time.Parse("2006-01-02T15:04", *req.CreatedDate); err == nil {
				item.CreatedAt = parsedTime
			} else {
				item.CreatedAt = time.Now() // Fallback to current time
			}
		}
	} else {
		item.CreatedAt = time.Now()
	}
	item.UpdatedAt = item.CreatedAt

	err := s.ingredientRepo.Create(ctx, item)
	if err != nil {
		return nil, err
	}

	// Create initial stock history record if quantity > 0
	if req.Quantity > 0 {
		userID := primitive.NilObjectID
		if userIDStr != "" {
			if oid, err := primitive.ObjectIDFromHex(userIDStr); err == nil {
				userID = oid
			}
		}
		
		history := &ingredient.StockHistory{
			IngredientID: item.ID,
			Type:         ingredient.TransactionPurchase,
			Quantity:     req.Quantity,
			BeforeQty:    0,                              // Initial creation, before quantity is 0
			AfterQty:     req.Quantity,
			Reason:       "Tạo nguyên liệu mới - Nhập kho đầu tiên",
			UserID:       userID,
			Username:     username,
			CostPerUnit:  req.CostPerUnit,
			TotalCost:    req.Quantity * req.CostPerUnit,
			CreatedAt:    item.CreatedAt, // Use the same created_at as ingredient
		}
		
		// Create stock history (don't fail if this fails)
		if err := s.stockHistoryRepo.Create(ctx, history); err != nil {
			// Log error but don't fail the operation
		}
	}

	// Track expense for initial purchase if AutoExpenseService is configured
	if s.autoExpenseService != nil && req.Quantity > 0 {
		if err := s.autoExpenseService.TrackIngredientPurchase(ctx, item, req.Quantity, username); err != nil {
			// Log error but don't fail the operation
			// The ingredient was created successfully, expense tracking is secondary
		}
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

	// Track if cost_per_unit changed to trigger recalculation
	costChanged := false
	if req.CostPerUnit != nil && *req.CostPerUnit != item.CostPerUnit {
		costChanged = true
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
	if req.ConversionRate != nil {
		item.ConversionRate = *req.ConversionRate
	}
	if req.WastagePercentage != nil {
		item.WastagePercentage = *req.WastagePercentage
	}

	err = s.ingredientRepo.Update(ctx, id, item)
	if err != nil {
		return nil, err
	}

	// Trigger cost recalculation if cost_per_unit changed
	// Requirements: 1.3, 9.1
	if costChanged {
		s.triggerCostRecalculation(ctx, id)
	}

	return item, nil
}

func (s *IngredientService) DeleteIngredient(ctx context.Context, id primitive.ObjectID) error {
	return s.ingredientRepo.Delete(ctx, id)
}

// StockIn - Add stock (purchase/receive)
func (s *IngredientService) StockIn(ctx context.Context, id primitive.ObjectID, req *ingredient.StockInRequest) (*ingredient.Ingredient, error) {
	item, err := s.ingredientRepo.FindByID(ctx, id)
	if err != nil {
		return nil, err
	}

	beforeQty := item.Quantity
	item.Quantity += req.Quantity
	afterQty := item.Quantity

	// Determine cost per unit for this transaction
	costPerUnit := req.CostPerUnit  // Use provided price, or 0 if not provided
	userProvidedNewPrice := req.CostPerUnit > 0 // Track if user provided a new price

	// Track if cost_per_unit changed to trigger recalculation
	costChanged := false

	// Calculate weighted average cost ONLY when new price is provided and different
	if req.CostPerUnit > 0 && req.CostPerUnit != item.CostPerUnit && afterQty > 0 {
		// Weighted average: (old_qty * old_price + new_qty * new_price) / total_qty
		oldValue := beforeQty * item.CostPerUnit
		newValue := req.Quantity * req.CostPerUnit
		item.CostPerUnit = (oldValue + newValue) / afterQty
		costChanged = true
	}

	err = s.ingredientRepo.Update(ctx, id, item)
	if err != nil {
		return nil, err
	}

	// Create stock history record
	// For history: only record cost if user provided a price
	// If gifted/found (no price), record as 0 cost
	userID, _ := primitive.ObjectIDFromHex(req.UserID)
	reason := req.Reason
	if reason == "" {
		reason = "Nhập kho"
	}
	
	history := &ingredient.StockHistory{
		IngredientID: id,
		Type:         ingredient.TransactionPurchase,
		Quantity:     req.Quantity,
		BeforeQty:    beforeQty,
		AfterQty:     afterQty,
		Reason:       reason,
		UserID:       userID,
		Username:     req.Username,
		CostPerUnit:  costPerUnit,  // 0 if gifted, req.CostPerUnit if purchased
		TotalCost:    req.Quantity * costPerUnit,  // 0 if gifted
	}
	s.stockHistoryRepo.Create(ctx, history)

	// Track expense ONLY if user provided a new price (actual purchase, not gifted)
	if s.autoExpenseService != nil && userProvidedNewPrice {
		tempItem := *item
		tempItem.CostPerUnit = req.CostPerUnit // Use the price user provided
		if err := s.autoExpenseService.TrackIngredientPurchase(ctx, &tempItem, req.Quantity, req.Username); err != nil {
			// Log error but don't fail the operation
		}
	}

	// Trigger cost recalculation if cost_per_unit changed
	// Requirements: 1.3, 9.1
	if costChanged {
		s.triggerCostRecalculation(ctx, id)
	}

	return item, nil
}

// StockOut - Remove stock (usage/waste)
func (s *IngredientService) StockOut(ctx context.Context, id primitive.ObjectID, req *ingredient.StockOutRequest) (*ingredient.Ingredient, error) {
	item, err := s.ingredientRepo.FindByID(ctx, id)
	if err != nil {
		return nil, err
	}

	beforeQty := item.Quantity
	item.Quantity -= req.Quantity
	if item.Quantity < 0 {
		item.Quantity = 0
	}
	afterQty := item.Quantity

	// Price never changes on stock out
	err = s.ingredientRepo.Update(ctx, id, item)
	if err != nil {
		return nil, err
	}

	// Create stock history record
	userID, _ := primitive.ObjectIDFromHex(req.UserID)
	history := &ingredient.StockHistory{
		IngredientID: id,
		Type:         ingredient.TransactionWaste,
		Quantity:     -req.Quantity, // Negative for removal
		BeforeQty:    beforeQty,
		AfterQty:     afterQty,
		Reason:       req.Reason,
		UserID:       userID,
		Username:     req.Username,
		CostPerUnit:  item.CostPerUnit, // Current price (unchanged)
		TotalCost:    -req.Quantity * item.CostPerUnit,
	}
	s.stockHistoryRepo.Create(ctx, history)

	return item, nil
}

// StockAdjust - Set stock to specific quantity (inventory correction)
func (s *IngredientService) StockAdjust(ctx context.Context, id primitive.ObjectID, req *ingredient.StockAdjustRequest) (*ingredient.Ingredient, error) {
	item, err := s.ingredientRepo.FindByID(ctx, id)
	if err != nil {
		return nil, err
	}

	beforeQty := item.Quantity
	quantityDiff := req.NewQuantity - beforeQty
	item.Quantity = req.NewQuantity
	afterQty := item.Quantity

	// Determine cost per unit for this transaction
	costPerUnit := float64(0)       // Default to 0 (no cost)
	userProvidedNewPrice := false   // Track if user actually provided a new price
	costChanged := false            // Track if cost_per_unit changed
	
	// Only recalculate price if:
	// 1. Quantity increased (positive diff)
	// 2. New price provided and different from current
	if quantityDiff > 0 && req.CostPerUnit > 0 && req.CostPerUnit != item.CostPerUnit && afterQty > 0 {
		// Weighted average for the increase
		oldValue := beforeQty * item.CostPerUnit
		newValue := quantityDiff * req.CostPerUnit
		item.CostPerUnit = (oldValue + newValue) / afterQty
		costPerUnit = req.CostPerUnit // Use new price for history
		userProvidedNewPrice = true   // User provided a new price
		costChanged = true
	}

	err = s.ingredientRepo.Update(ctx, id, item)
	if err != nil {
		return nil, err
	}

	// Create stock history record
	// For history: only record cost if user provided a new price
	// If gifted/found (no price), record as 0 cost
	userID, _ := primitive.ObjectIDFromHex(req.UserID)
	history := &ingredient.StockHistory{
		IngredientID: id,
		Type:         ingredient.TransactionAdjustment,
		Quantity:     quantityDiff,
		BeforeQty:    beforeQty,
		AfterQty:     afterQty,
		Reason:       req.Reason,
		UserID:       userID,
		Username:     req.Username,
		CostPerUnit:  costPerUnit,  // 0 if gifted, req.CostPerUnit if purchased
		TotalCost:    quantityDiff * costPerUnit,  // 0 if gifted
	}
	s.stockHistoryRepo.Create(ctx, history)

	// Track expense ONLY if:
	// 1. Quantity increased AND
	// 2. User provided a new price (not gifted/found)
	if s.autoExpenseService != nil && quantityDiff > 0 && userProvidedNewPrice {
		tempItem := *item
		tempItem.CostPerUnit = costPerUnit
		if err := s.autoExpenseService.TrackIngredientPurchase(ctx, &tempItem, quantityDiff, req.Username); err != nil {
			// Log error but don't fail the operation
		}
	}

	// Trigger cost recalculation if cost_per_unit changed
	// Requirements: 1.3, 9.1
	if costChanged {
		s.triggerCostRecalculation(ctx, id)
	}

	return item, nil
}

// AdjustStock - Legacy method for backward compatibility
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

	// Determine cost per unit for this transaction
	costPerUnit := req.CostPerUnit
	if costPerUnit <= 0 {
		// If no cost provided, use current cost
		costPerUnit = item.CostPerUnit
	}

	// Track if cost_per_unit changed to trigger recalculation
	costChanged := false

	// Calculate weighted average cost ONLY for stock IN with NEW price
	// Only recalculate if:
	// 1. Quantity is positive (adding stock)
	// 2. A new cost per unit is provided (req.CostPerUnit > 0)
	// 3. The new price is different from current price
	// 4. After quantity is positive
	if req.Quantity > 0 && req.CostPerUnit > 0 && req.CostPerUnit != item.CostPerUnit && afterQty > 0 {
		// Weighted average: (old_qty * old_price + new_qty * new_price) / total_qty
		oldValue := beforeQty * item.CostPerUnit
		newValue := req.Quantity * req.CostPerUnit
		item.CostPerUnit = (oldValue + newValue) / afterQty
		costChanged = true
	}
	// If no new price provided (req.CostPerUnit = 0), keep current price
	// This handles: stock OUT, stock SET without price, stock IN without price

	err = s.ingredientRepo.Update(ctx, id, item)
	if err != nil {
		return nil, err
	}

	// Create stock history record with price information
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
		CostPerUnit:  costPerUnit,                    // Price at time of transaction
		TotalCost:    req.Quantity * costPerUnit,     // Total cost for this transaction
	}
	s.stockHistoryRepo.Create(ctx, history)

	// Track expense for stock IN (positive quantity adjustment)
	// Use the actual cost per unit for this transaction
	if s.autoExpenseService != nil && req.Quantity > 0 && costPerUnit > 0 {
		// Create a temporary ingredient with the transaction price for expense tracking
		tempItem := *item
		tempItem.CostPerUnit = costPerUnit
		if err := s.autoExpenseService.TrackIngredientPurchase(ctx, &tempItem, req.Quantity, req.Username); err != nil {
			// Log error but don't fail the operation
			// The stock adjustment was successful, expense tracking is secondary
		}
	}

	// Trigger cost recalculation if cost_per_unit changed
	// Requirements: 1.3, 9.1
	if costChanged {
		s.triggerCostRecalculation(ctx, id)
	}

	return item, nil
}

// DeductStock deducts ingredient quantity when an order is PAID, and records StockHistory
func (s *IngredientService) DeductStock(
	ctx context.Context,
	ingredientID primitive.ObjectID,
	quantity float64,
	unit ingredient.UnitType,
	orderID primitive.ObjectID,
	orderNumber string,
) error {
	ing, err := s.ingredientRepo.FindByID(ctx, ingredientID)
	if err != nil {
		return err
	}

	convRate := ingredient.GetConversionRate(ingredient.UnitType(ing.Unit), unit)
	qtyInStock := quantity * convRate

	beforeQty := ing.Quantity
	ing.Quantity -= qtyInStock
	if ing.Quantity < 0 {
		ing.Quantity = 0
	}
	afterQty := ing.Quantity

	if err := s.ingredientRepo.Update(ctx, ingredientID, ing); err != nil {
		return err
	}

	history := &ingredient.StockHistory{
		IngredientID: ingredientID,
		Type:         ingredient.TransactionOrder,
		Quantity:     -qtyInStock,
		BeforeQty:    beforeQty,
		AfterQty:     afterQty,
		Reason:       "Order #" + orderNumber,
		CostPerUnit:  ing.CostPerUnit,
		TotalCost:    -qtyInStock * ing.CostPerUnit,
	}
	s.stockHistoryRepo.Create(ctx, history)

	return nil
}

// RestoreStock restores ingredient quantity (rollback of a deduction)
func (s *IngredientService) RestoreStock(
	ctx context.Context,
	ingredientID primitive.ObjectID,
	quantity float64,
	unit ingredient.UnitType,
	orderID primitive.ObjectID,
) error {
	ing, err := s.ingredientRepo.FindByID(ctx, ingredientID)
	if err != nil {
		return err
	}

	convRate := ingredient.GetConversionRate(ingredient.UnitType(ing.Unit), unit)
	qtyInStock := quantity * convRate

	beforeQty := ing.Quantity
	ing.Quantity += qtyInStock
	afterQty := ing.Quantity

	if err := s.ingredientRepo.Update(ctx, ingredientID, ing); err != nil {
		return err
	}

	history := &ingredient.StockHistory{
		IngredientID: ingredientID,
		Type:         ingredient.TransactionAdjustment,
		Quantity:     qtyInStock,
		BeforeQty:    beforeQty,
		AfterQty:     afterQty,
		Reason:       "Hoàn trả tồn kho (rollback order)",
		CostPerUnit:  ing.CostPerUnit,
		TotalCost:    qtyInStock * ing.CostPerUnit,
	}
	s.stockHistoryRepo.Create(ctx, history)

	return nil
}

func (s *IngredientService) GetLowStockIngredients(ctx context.Context) ([]*ingredient.Ingredient, error) {
	return s.ingredientRepo.FindLowStock(ctx)
}

func (s *IngredientService) GetStockHistory(ctx context.Context, id primitive.ObjectID) ([]*ingredient.StockHistory, error) {
	return s.stockHistoryRepo.FindByIngredientID(ctx, id)
}

// Category methods
func (s *IngredientService) CreateCategory(ctx context.Context, req *ingredient.CreateCategoryRequest) (*ingredient.IngredientCategory, error) {
	cat := &ingredient.IngredientCategory{
		Name: req.Name,
	}
	
	err := s.ingredientRepo.CreateCategory(ctx, cat)
	if err != nil {
		return nil, err
	}
	
	return cat, nil
}

func (s *IngredientService) GetCategories(ctx context.Context) ([]ingredient.IngredientCategory, error) {
	return s.ingredientRepo.GetCategories(ctx)
}

func (s *IngredientService) DeleteCategory(ctx context.Context, id primitive.ObjectID) error {
	return s.ingredientRepo.DeleteCategory(ctx, id)
}
