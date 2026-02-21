package printing

import (
	"context"
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

// PrintJobRepository defines the interface for print job persistence
type PrintJobRepository interface {
	Create(ctx context.Context, job *PrintJob) error
	FindByID(ctx context.Context, id primitive.ObjectID) (*PrintJob, error)
	FindByOrderID(ctx context.Context, orderID primitive.ObjectID) ([]*PrintJob, error)
	FindPending(ctx context.Context, limit int) ([]*PrintJob, error)
	FindFailed(ctx context.Context) ([]*PrintJob, error)
	UpdateStatus(ctx context.Context, id primitive.ObjectID, status PrintJobStatus, errorMsg string) error
	IncrementRetry(ctx context.Context, id primitive.ObjectID) error
	Delete(ctx context.Context, id primitive.ObjectID) error
	DeleteOldCompleted(ctx context.Context, olderThan time.Time) error
}

// PrinterConfigRepository defines the interface for printer configuration persistence
type PrinterConfigRepository interface {
	Create(ctx context.Context, config *PrinterConfig) error
	FindByID(ctx context.Context, id primitive.ObjectID) (*PrinterConfig, error)
	FindAll(ctx context.Context) ([]*PrinterConfig, error)
	FindByType(ctx context.Context, printerType PrinterType) ([]*PrinterConfig, error)
	FindDefault(ctx context.Context, printerType PrinterType) (*PrinterConfig, error)
	Update(ctx context.Context, config *PrinterConfig) error
	Delete(ctx context.Context, id primitive.ObjectID) error
}

// PrintTemplateRepository defines the interface for print template persistence
type PrintTemplateRepository interface {
	Create(ctx context.Context, template *PrintTemplate) error
	FindByID(ctx context.Context, id primitive.ObjectID) (*PrintTemplate, error)
	FindByType(ctx context.Context, templateType TemplateType) ([]*PrintTemplate, error)
	FindDefault(ctx context.Context, templateType TemplateType) (*PrintTemplate, error)
	Update(ctx context.Context, template *PrintTemplate) error
	Delete(ctx context.Context, id primitive.ObjectID) error
}
