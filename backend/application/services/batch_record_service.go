package services

import (
	"context"
	"fmt"
	"time"

	"cafe-pos/backend/domain/batch"
	"cafe-pos/backend/domain/ingredient"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

// BatchRecordService handles batch record operations
type BatchRecordService struct {
	batchRecordRepo      batch.BatchRecordRepository
	batchDefinitionRepo  batch.BatchDefinitionRepository
	ingredientRepo       IngredientRepository
	stockHistoryRepo     StockHistoryRepository
	userRepo             UserRepository
	batchCostCalculator  *BatchCostCalculator
	mongoClient          *mongo.Client
}

// NewBatchRecordService creates a new batch record service
func NewBatchRecordService(
	batchRecordRepo batch.BatchRecordRepository,
	batchDefinitionRepo batch.BatchDefinitionRepository,
	ingredientRepo IngredientRepository,
	stockHistoryRepo StockHistoryRepository,
	userRepo UserRepository,
	batchCostCalculator *BatchCostCalculator,
	mongoClient *mongo.Client,
) *BatchRecordService {
	return &BatchRecordService{
		batchRecordRepo:     batchRecordRepo,
		batchDefinitionRepo: batchDefinitionRepo,
		ingredientRepo:      ingredientRepo,
		stockHistoryRepo:    stockHistoryRepo,
		userRepo:            userRepo,
		batchCostCalculator: batchCostCalculator,
		mongoClient:         mongoClient,
	}
}

// CreateBatchRequest represents a request to create a batch record
type CreateBatchRequest struct {
	BatchDefinitionID primitive.ObjectID
	QuantityProduced  float64
	PreparedBy        string
}

// CreateBatch creates a new batch record with transaction support
// Steps:
// 1. Fetch batch definition
// 2. Calculate required ingredients with wastage
// 3. Check ingredient availability
// 4. Calculate cost
// 5. Start MongoDB transaction
// 6. Deduct ingredients from inventory
// 7. Create stock history records
// 8. Create batch record
// 9. Commit transaction (or rollback on error)
func (s *BatchRecordService) CreateBatch(ctx context.Context, req CreateBatchRequest) (*batch.BatchRecord, error) {
	if req.QuantityProduced <= 0 {
		return nil, fmt.Errorf("quantity produced must be greater than 0")
	}

	// 1. Fetch batch definition
	batchDef, err := s.batchDefinitionRepo.FindByID(ctx, req.BatchDefinitionID)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch batch definition: %w", err)
	}

	// 2. Look up username from user ID
	preparedByName := req.PreparedBy // Default to ID if lookup fails
	if userID, err := primitive.ObjectIDFromHex(req.PreparedBy); err == nil {
		if user, err := s.userRepo.FindByID(ctx, userID); err == nil && user != nil {
			preparedByName = user.Username
		}
	}

	// 3. Calculate cost breakdown (includes wastage calculation)
	costBreakdown, err := s.batchCostCalculator.CalculateBatchCost(ctx, batchDef, req.QuantityProduced)
	if err != nil {
		return nil, fmt.Errorf("failed to calculate batch cost: %w", err)
	}

	// 4. Check ingredient availability before starting transaction
	for _, ingredientCost := range costBreakdown.IngredientCosts {
		ing, err := s.ingredientRepo.FindByID(ctx, ingredientCost.IngredientID)
		if err != nil {
			return nil, fmt.Errorf("ingredient %s not found: %w", ingredientCost.IngredientName, err)
		}

		// Convert ingredientCost.Quantity from source unit to stock unit for comparison
		conversionRate := ingredient.GetConversionRate(ingredient.UnitType(ing.Unit), ingredient.UnitType(ingredientCost.Unit))
		quantityInStockUnit := ingredientCost.Quantity * conversionRate

		if ing.Quantity < quantityInStockUnit {
			return nil, fmt.Errorf("insufficient %s: need %.2f%s (%.2f%s in stock unit), have %.2f%s",
				ingredientCost.IngredientName,
				ingredientCost.Quantity,
				ingredientCost.Unit,
				quantityInStockUnit,
				ing.Unit,
				ing.Quantity,
				ing.Unit,
			)
		}
	}

	// 5. Start MongoDB transaction for atomic operations
	session, err := s.mongoClient.StartSession()
	if err != nil {
		return nil, fmt.Errorf("failed to start transaction session: %w", err)
	}
	defer session.EndSession(ctx)

	var batchRecord *batch.BatchRecord

	// Execute transaction with automatic rollback on error
	err = mongo.WithSession(ctx, session, func(sessCtx mongo.SessionContext) error {
		if err := session.StartTransaction(); err != nil {
			return fmt.Errorf("failed to start transaction: %w", err)
		}

		// 6. Deduct ingredients from inventory and create stock history
		ingredientsUsed := make([]batch.IngredientUsage, 0, len(costBreakdown.IngredientCosts))
		
		for _, ingredientCost := range costBreakdown.IngredientCosts {
			ing, err := s.ingredientRepo.FindByID(sessCtx, ingredientCost.IngredientID)
			if err != nil {
				session.AbortTransaction(sessCtx)
				return fmt.Errorf("failed to fetch ingredient %s: %w", ingredientCost.IngredientName, err)
			}

			beforeQty := ing.Quantity
			
			// Convert ingredientCost.Quantity from source unit to stock unit for deduction
			conversionRate := ingredient.GetConversionRate(ingredient.UnitType(ing.Unit), ingredient.UnitType(ingredientCost.Unit))
			quantityToDeduct := ingredientCost.Quantity * conversionRate
			
			// Deduct quantity in stock unit
			ing.Quantity -= quantityToDeduct
			if ing.Quantity < 0 {
				ing.Quantity = 0
			}
			afterQty := ing.Quantity

			err = s.ingredientRepo.Update(sessCtx, ingredientCost.IngredientID, ing)
			if err != nil {
				session.AbortTransaction(sessCtx)
				return fmt.Errorf("failed to deduct ingredient %s: %w", ingredientCost.IngredientName, err)
			}

			// Create stock history record for batch creation
			userID := primitive.NilObjectID
			if req.PreparedBy != "" {
				if oid, err := primitive.ObjectIDFromHex(req.PreparedBy); err == nil {
					userID = oid
				}
			}

			history := &ingredient.StockHistory{
				IngredientID: ingredientCost.IngredientID,
				Type:         ingredient.TransactionOrder, // Using "order" type for batch creation
				Quantity:     -quantityToDeduct,           // Negative for deduction, in stock unit
				BeforeQty:    beforeQty,
				AfterQty:     afterQty,
				Reason:       fmt.Sprintf("Chế biến batch: %s (%.2f%s)", batchDef.Name, req.QuantityProduced, batchDef.Unit),
				UserID:       userID,
				Username:     preparedByName, // Use the looked-up username
				CostPerUnit:  ingredientCost.CostPerUnit,
				TotalCost:    -ingredientCost.TotalCost, // Negative for cost deduction
				CreatedAt:    time.Now(),
			}

			err = s.stockHistoryRepo.Create(sessCtx, history)
			if err != nil {
				session.AbortTransaction(sessCtx)
				return fmt.Errorf("failed to create stock history for %s: %w", ingredientCost.IngredientName, err)
			}

			// Record ingredient usage (keep in source unit for display)
			ingredientsUsed = append(ingredientsUsed, batch.IngredientUsage{
				IngredientID:   ingredientCost.IngredientID,
				IngredientName: ingredientCost.IngredientName,
				Quantity:       ingredientCost.Quantity,
				Unit:           ingredientCost.Unit,
				CostPerUnit:    ingredientCost.CostPerUnit,
				TotalCost:      ingredientCost.TotalCost,
			})
		}

		// 7. Calculate expiry time
		preparedAt := time.Now()
		expiresAt := preparedAt.Add(time.Duration(batchDef.ShelfLifeHours) * time.Hour)

		// 8. Create batch record
		batchRecord = &batch.BatchRecord{
			BatchDefinitionID: req.BatchDefinitionID,
			BatchName:         batchDef.Name,
			QuantityProduced:  req.QuantityProduced,
			QuantityRemaining: req.QuantityProduced,
			Unit:              batchDef.Unit,
			CostPerUnit:       costBreakdown.CostPerUnit,
			TotalCost:         costBreakdown.TotalCost,
			PreparedBy:        req.PreparedBy,
			PreparedByName:    preparedByName, // Add the username
			PreparedAt:        preparedAt,
			ExpiresAt:         expiresAt,
			Status:            batch.BatchStatusAvailable,
			IngredientsUsed:   ingredientsUsed,
			CreatedAt:         preparedAt,
			UpdatedAt:         preparedAt,
		}

		err = s.batchRecordRepo.Create(sessCtx, batchRecord)
		if err != nil {
			session.AbortTransaction(sessCtx)
			return fmt.Errorf("failed to create batch record: %w", err)
		}

		// 9. Commit transaction
		if err := session.CommitTransaction(sessCtx); err != nil {
			return fmt.Errorf("failed to commit transaction: %w", err)
		}

		return nil
	})

	if err != nil {
		return nil, err
	}

	return batchRecord, nil
}

// GetByID retrieves a batch record by ID
func (s *BatchRecordService) GetByID(ctx context.Context, id primitive.ObjectID) (*batch.BatchRecord, error) {
	return s.batchRecordRepo.FindByID(ctx, id)
}

// List retrieves batch records with filters
func (s *BatchRecordService) List(ctx context.Context, filter batch.BatchRecordFilter) ([]*batch.BatchRecord, error) {
	records, _, err := s.batchRecordRepo.FindAll(ctx, filter)
	return records, err
}

// UpdateQuantity updates the remaining quantity of a batch record
func (s *BatchRecordService) UpdateQuantity(ctx context.Context, id primitive.ObjectID, newQuantity float64) error {
	if newQuantity < 0 {
		return fmt.Errorf("quantity cannot be negative")
	}

	record, err := s.batchRecordRepo.FindByID(ctx, id)
	if err != nil {
		return fmt.Errorf("batch record not found: %w", err)
	}

	record.QuantityRemaining = newQuantity
	record.UpdatedAt = time.Now()

	// Update status if depleted
	if newQuantity == 0 {
		record.Status = batch.BatchStatusDepleted
	}

	return s.batchRecordRepo.Update(ctx, record)
}

// MarkAsExpired marks a batch record as expired
func (s *BatchRecordService) MarkAsExpired(ctx context.Context, id primitive.ObjectID) error {
	record, err := s.batchRecordRepo.FindByID(ctx, id)
	if err != nil {
		return fmt.Errorf("batch record not found: %w", err)
	}

	record.Status = batch.BatchStatusExpired
	record.UpdatedAt = time.Now()

	return s.batchRecordRepo.Update(ctx, record)
}

// Delete deletes a batch record and restores ingredients to inventory (only if not used)
func (s *BatchRecordService) Delete(ctx context.Context, id primitive.ObjectID) error {
	record, err := s.batchRecordRepo.FindByID(ctx, id)
	if err != nil {
		return fmt.Errorf("batch record not found: %w", err)
	}

	// Check if batch has been used
	if record.QuantityRemaining < record.QuantityProduced {
		return fmt.Errorf("cannot delete batch that has been partially used")
	}

	// Start transaction to restore ingredients and delete batch atomically
	session, err := s.mongoClient.StartSession()
	if err != nil {
		return fmt.Errorf("failed to start transaction session: %w", err)
	}
	defer session.EndSession(ctx)

	err = mongo.WithSession(ctx, session, func(sessCtx mongo.SessionContext) error {
		if err := session.StartTransaction(); err != nil {
			return fmt.Errorf("failed to start transaction: %w", err)
		}

		// Restore ingredients to inventory
		for _, ingredientUsage := range record.IngredientsUsed {
			ing, err := s.ingredientRepo.FindByID(sessCtx, ingredientUsage.IngredientID)
			if err != nil {
				session.AbortTransaction(sessCtx)
				return fmt.Errorf("failed to fetch ingredient %s: %w", ingredientUsage.IngredientName, err)
			}

			beforeQty := ing.Quantity
			ing.Quantity += ingredientUsage.Quantity
			afterQty := ing.Quantity

			err = s.ingredientRepo.Update(sessCtx, ingredientUsage.IngredientID, ing)
			if err != nil {
				session.AbortTransaction(sessCtx)
				return fmt.Errorf("failed to restore ingredient %s: %w", ingredientUsage.IngredientName, err)
			}

			// Create stock history record for ingredient restoration
			userID := primitive.NilObjectID
			if record.PreparedBy != "" {
				if oid, err := primitive.ObjectIDFromHex(record.PreparedBy); err == nil {
					userID = oid
				}
			}

			history := &ingredient.StockHistory{
				IngredientID: ingredientUsage.IngredientID,
				Type:         ingredient.TransactionAdjustment,
				Quantity:     ingredientUsage.Quantity, // Positive for restoration
				BeforeQty:    beforeQty,
				AfterQty:     afterQty,
				Reason:       fmt.Sprintf("Hoàn trả từ batch đã xóa: %s", record.BatchName),
				UserID:       userID,
				Username:     record.PreparedBy,
				CostPerUnit:  ingredientUsage.CostPerUnit,
				TotalCost:    ingredientUsage.TotalCost,
				CreatedAt:    time.Now(),
			}

			err = s.stockHistoryRepo.Create(sessCtx, history)
			if err != nil {
				session.AbortTransaction(sessCtx)
				return fmt.Errorf("failed to create stock history for %s: %w", ingredientUsage.IngredientName, err)
			}
		}

		// Delete batch record
		err = s.batchRecordRepo.Delete(sessCtx, id)
		if err != nil {
			session.AbortTransaction(sessCtx)
			return fmt.Errorf("failed to delete batch record: %w", err)
		}

		// Commit transaction
		if err := session.CommitTransaction(sessCtx); err != nil {
			return fmt.Errorf("failed to commit transaction: %w", err)
		}

		return nil
	})

	return err
}

// GetAvailableBatches retrieves available batches for a batch definition
// Sorted by expiry date (FIFO)
func (s *BatchRecordService) GetAvailableBatches(ctx context.Context, batchDefID primitive.ObjectID) ([]*batch.BatchRecord, error) {
	return s.batchRecordRepo.FindAvailableByDefinition(ctx, batchDefID)
}

// MarkExpiredBatches marks all expired batches as expired (background job)
func (s *BatchRecordService) MarkExpiredBatches(ctx context.Context) error {
	// Find all available batches that have expired
	filter := batch.BatchRecordFilter{
		Status: batch.BatchStatusAvailable,
	}
	
	records, _, err := s.batchRecordRepo.FindAll(ctx, filter)
	if err != nil {
		return fmt.Errorf("failed to fetch batch records: %w", err)
	}

	now := time.Now()
	expiredCount := 0

	for _, record := range records {
		if now.After(record.ExpiresAt) {
			record.Status = batch.BatchStatusExpired
			record.UpdatedAt = now
			
			err = s.batchRecordRepo.Update(ctx, record)
			if err != nil {
				// Log error but continue with other records
				continue
			}
			expiredCount++
		}
	}

	return nil
}
