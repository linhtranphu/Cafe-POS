package batch

import (
	"context"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

// BatchDefinitionRepository defines the interface for batch definition persistence
type BatchDefinitionRepository interface {
	Create(ctx context.Context, def *BatchDefinition) error
	Update(ctx context.Context, def *BatchDefinition) error
	Delete(ctx context.Context, id primitive.ObjectID) error
	FindByID(ctx context.Context, id primitive.ObjectID) (*BatchDefinition, error)
	FindAll(ctx context.Context, filter BatchDefinitionFilter) ([]*BatchDefinition, int64, error)
}

// BatchRecordRepository defines the interface for batch record persistence
type BatchRecordRepository interface {
	Create(ctx context.Context, record *BatchRecord) error
	Update(ctx context.Context, record *BatchRecord) error
	Delete(ctx context.Context, id primitive.ObjectID) error
	FindByID(ctx context.Context, id primitive.ObjectID) (*BatchRecord, error)
	FindAll(ctx context.Context, filter BatchRecordFilter) ([]*BatchRecord, int64, error)
	FindAvailableByDefinition(ctx context.Context, defID primitive.ObjectID) ([]*BatchRecord, error)
	UpdateQuantity(ctx context.Context, id primitive.ObjectID, newQuantity float64) error
	GetTotalAvailableQuantity(ctx context.Context, defID primitive.ObjectID) (float64, error)
}

// BatchUsageLogRepository defines the interface for batch usage log persistence
type BatchUsageLogRepository interface {
	Create(ctx context.Context, log *BatchUsageLog) error
	FindAll(ctx context.Context, filter BatchUsageLogFilter) ([]*BatchUsageLog, int64, error)
	Delete(ctx context.Context, id primitive.ObjectID) error
}
