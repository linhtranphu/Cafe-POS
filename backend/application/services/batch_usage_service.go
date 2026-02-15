package services

import (
	"context"
	"fmt"
	"math"
	"time"

	"cafe-pos/backend/domain/batch"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

// BatchUsageService handles batch usage operations
type BatchUsageService struct {
	batchRecordRepo    batch.BatchRecordRepository
	batchUsageLogRepo  batch.BatchUsageLogRepository
}

// NewBatchUsageService creates a new batch usage service
func NewBatchUsageService(
	batchRecordRepo batch.BatchRecordRepository,
	batchUsageLogRepo batch.BatchUsageLogRepository,
) *BatchUsageService {
	return &BatchUsageService{
		batchRecordRepo:   batchRecordRepo,
		batchUsageLogRepo: batchUsageLogRepo,
	}
}

// UseBatchRequest represents a request to use batch
type UseBatchRequest struct {
	BatchDefinitionID primitive.ObjectID
	QuantityNeeded    float64
	OrderID           primitive.ObjectID
	MenuItemID        primitive.ObjectID
	MenuItemName      string
}

// BatchUsageResult represents the result of batch usage
type BatchUsageResult struct {
	Success     bool
	BatchesUsed []BatchUsageDetail
	TotalCost   float64
	Message     string
}

// BatchUsageDetail represents details of a single batch usage
type BatchUsageDetail struct {
	BatchRecordID primitive.ObjectID
	QuantityUsed  float64
	CostPerUnit   float64
}

// UseBatch uses batch with FIFO logic
// Steps:
// 1. Get available batches sorted by expiry (FIFO)
// 2. Check expiry before using
// 3. Use batches in order until quantity is satisfied
// 4. Update batch quantities
// 5. Log usage
func (s *BatchUsageService) UseBatch(ctx context.Context, req UseBatchRequest) (*BatchUsageResult, error) {
	if req.QuantityNeeded <= 0 {
		return nil, fmt.Errorf("quantity needed must be greater than 0")
	}

	// 1. Get available batches sorted by expiry date (FIFO)
	batches, err := s.batchRecordRepo.FindAvailableByDefinition(ctx, req.BatchDefinitionID)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch available batches: %w", err)
	}

	if len(batches) == 0 {
		return &BatchUsageResult{
			Success: false,
			Message: "No available batches found",
		}, nil
	}

	// 2. Filter out expired batches
	now := time.Now()
	availableBatches := make([]*batch.BatchRecord, 0)
	for _, b := range batches {
		if now.Before(b.ExpiresAt) && b.QuantityRemaining > 0 {
			availableBatches = append(availableBatches, b)
		}
	}

	if len(availableBatches) == 0 {
		return &BatchUsageResult{
			Success: false,
			Message: "No non-expired batches available",
		}, nil
	}

	// 3. Calculate total available quantity
	totalAvailable := 0.0
	for _, b := range availableBatches {
		totalAvailable += b.QuantityRemaining
	}

	if totalAvailable < req.QuantityNeeded {
		return &BatchUsageResult{
			Success: false,
			Message: fmt.Sprintf("Insufficient batch quantity. Need: %.2f, Available: %.2f",
				req.QuantityNeeded, totalAvailable),
		}, nil
	}

	// 4. Use batches in FIFO order
	remainingNeeded := req.QuantityNeeded
	batchesUsed := make([]BatchUsageDetail, 0)
	totalCost := 0.0

	for _, batchRecord := range availableBatches {
		if remainingNeeded <= 0 {
			break
		}

		// Calculate how much to use from this batch
		quantityToUse := math.Min(batchRecord.QuantityRemaining, remainingNeeded)

		// Update batch quantity
		batchRecord.QuantityRemaining -= quantityToUse
		batchRecord.UpdatedAt = now

		// Update status if depleted
		if batchRecord.QuantityRemaining == 0 {
			batchRecord.Status = batch.BatchStatusDepleted
		}

		err = s.batchRecordRepo.Update(ctx, batchRecord)
		if err != nil {
			// TODO: Implement rollback mechanism
			return nil, fmt.Errorf("failed to update batch record: %w", err)
		}

		// Calculate cost for this usage
		usageCost := quantityToUse * batchRecord.CostPerUnit
		totalCost += usageCost

		// Log usage
		usageLog := &batch.BatchUsageLog{
			BatchRecordID:  batchRecord.ID,
			BatchName:      batchRecord.BatchName,
			OrderID:        req.OrderID,
			MenuItemID:     req.MenuItemID,
			MenuItemName:   req.MenuItemName,
			QuantityUsed:   quantityToUse,
			Unit:           batchRecord.Unit,
			CostPerUnit:    batchRecord.CostPerUnit,
			TotalCost:      usageCost,
			UsedAt:         now,
		}

		err = s.batchUsageLogRepo.Create(ctx, usageLog)
		if err != nil {
			// Log error but don't fail the operation
			// The batch quantity has already been deducted
		}

		// Record batch usage detail
		batchesUsed = append(batchesUsed, BatchUsageDetail{
			BatchRecordID: batchRecord.ID,
			QuantityUsed:  quantityToUse,
			CostPerUnit:   batchRecord.CostPerUnit,
		})

		remainingNeeded -= quantityToUse
	}

	// Round total cost to 2 decimal places
	totalCost = math.Round(totalCost*100) / 100

	return &BatchUsageResult{
		Success:     true,
		BatchesUsed: batchesUsed,
		TotalCost:   totalCost,
		Message:     "Batch used successfully",
	}, nil
}

// GetUsageHistory retrieves usage history with filters
func (s *BatchUsageService) GetUsageHistory(ctx context.Context, filter batch.BatchUsageLogFilter) ([]*batch.BatchUsageLog, error) {
	logs, _, err := s.batchUsageLogRepo.FindAll(ctx, filter)
	return logs, err
}

// RollbackBatchUsage rolls back batch usage for an order
// This is used when an order creation fails after batches have been deducted
func (s *BatchUsageService) RollbackBatchUsage(ctx context.Context, orderID primitive.ObjectID) error {
	// Find all usage logs for this order
	filter := batch.BatchUsageLogFilter{
		OrderID: &orderID,
	}
	
	logs, err := s.GetUsageHistory(ctx, filter)
	if err != nil {
		return fmt.Errorf("failed to fetch usage logs for rollback: %w", err)
	}

	// Restore batch quantities
	for _, log := range logs {
		// Get the batch record
		batchRecord, err := s.batchRecordRepo.FindByID(ctx, log.BatchRecordID)
		if err != nil {
			// Log error but continue with other rollbacks
			continue
		}

		// Restore quantity
		batchRecord.QuantityRemaining += log.QuantityUsed
		batchRecord.UpdatedAt = time.Now()

		// Update status if it was depleted
		if batchRecord.Status == batch.BatchStatusDepleted && batchRecord.QuantityRemaining > 0 {
			batchRecord.Status = batch.BatchStatusAvailable
		}

		err = s.batchRecordRepo.Update(ctx, batchRecord)
		if err != nil {
			// Log error but continue with other rollbacks
			continue
		}
	}

	// Delete the usage logs
	// Note: In a production system, you might want to mark them as "rolled back" instead of deleting
	for _, log := range logs {
		_ = s.batchUsageLogRepo.Delete(ctx, log.ID)
	}

	return nil
}
