package services

import (
	"context"
	
	"cafe-pos/backend/domain/batch"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

// BatchRecordRepository interface for adapter (to avoid import cycle)
type BatchRecordRepositoryForAdapter interface {
	FindAvailableByDefinition(ctx context.Context, defID primitive.ObjectID) ([]*batch.BatchRecord, error)
}

// BatchRecordCostAdapter adapts BatchRecordRepository to CostCalculatorBatchRecordRepository
type BatchRecordCostAdapter struct {
	repo BatchRecordRepositoryForAdapter
}

// NewBatchRecordCostAdapter creates a new adapter
func NewBatchRecordCostAdapter(repo BatchRecordRepositoryForAdapter) *BatchRecordCostAdapter {
	return &BatchRecordCostAdapter{
		repo: repo,
	}
}

// FindAvailableByDefinition implements CostCalculatorBatchRecordRepository interface
func (a *BatchRecordCostAdapter) FindAvailableByDefinition(ctx context.Context, defID primitive.ObjectID) ([]*CostCalculatorBatchRecord, error) {
	records, err := a.repo.FindAvailableByDefinition(ctx, defID)
	if err != nil {
		return nil, err
	}
	
	// Convert to CostCalculatorBatchRecord
	result := make([]*CostCalculatorBatchRecord, len(records))
	for i, r := range records {
		result[i] = &CostCalculatorBatchRecord{
			CostPerUnit: r.CostPerUnit,
			Unit:        r.Unit,
		}
	}
	return result, nil
}
