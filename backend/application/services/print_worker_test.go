package services

import (
	"context"
	"errors"
	"testing"
	"time"

	"cafe-pos/backend/domain/printing"
	"cafe-pos/backend/infrastructure/printbridge"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

// MockPrinter is a mock implementation of Printer
type MockPrinter struct {
	mock.Mock
}

func (m *MockPrinter) Connect() error {
	args := m.Called()
	return args.Error(0)
}

func (m *MockPrinter) Disconnect() error {
	args := m.Called()
	return args.Error(0)
}

func (m *MockPrinter) Print(content string) error {
	args := m.Called(content)
	return args.Error(0)
}

func (m *MockPrinter) GetStatus() (PrinterStatus, error) {
	args := m.Called()
	return args.Get(0).(PrinterStatus), args.Error(1)
}

// MockPrinterManager is a mock implementation of PrinterManager
type MockPrinterManager struct {
	mock.Mock
}

func (m *MockPrinterManager) GetPrinter(config *printing.PrinterConfig) (Printer, error) {
	args := m.Called(config)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(Printer), args.Error(1)
}

func (m *MockPrinterManager) TestConnection(config *printing.PrinterConfig) error {
	args := m.Called(config)
	return args.Error(0)
}

func (m *MockPrinterManager) SetPrintBridgeClient(client *printbridge.Client) {}

// Helper function to create a test print job for worker tests
func createTestPrintJobForWorker(status printing.PrintJobStatus, retryCount int) *printing.PrintJob {
	return &printing.PrintJob{
		ID:          primitive.NewObjectID(),
		Type:        printing.PrintJobTypeBill,
		OrderID:     primitive.NewObjectID(),
		OrderNumber: "ORD-001",
		PrinterID:   primitive.NewObjectID(),
		Content:     "Test content",
		Status:      status,
		RetryCount:  retryCount,
		MaxRetries:  3,
		CreatedAt:   time.Now(),
		UpdatedAt:   time.Now(),
	}
}

func TestNewPrintWorker(t *testing.T) {
	mockJobRepo := new(MockPrintJobRepository)
	mockConfigRepo := new(MockPrinterConfigRepository)
	mockManager := new(MockPrinterManager)

	worker := NewPrintWorker(PrintWorkerConfig{
		PrintJobRepo:      mockJobRepo,
		PrinterConfigRepo: mockConfigRepo,
		PrinterManager:    mockManager,
		PollInterval:      5 * time.Second,
	})

	assert.NotNil(t, worker)
}

func TestNewPrintWorker_DefaultPollInterval(t *testing.T) {
	mockJobRepo := new(MockPrintJobRepository)
	mockConfigRepo := new(MockPrinterConfigRepository)
	mockManager := new(MockPrinterManager)

	worker := NewPrintWorker(PrintWorkerConfig{
		PrintJobRepo:      mockJobRepo,
		PrinterConfigRepo: mockConfigRepo,
		PrinterManager:    mockManager,
		// PollInterval not set - should default to 10 seconds
	})

	assert.NotNil(t, worker)
	// Check that default poll interval is set
	w := worker.(*printWorker)
	assert.Equal(t, 10*time.Second, w.pollInterval)
}

func TestProcessJob_Success(t *testing.T) {
	// Setup
	mockJobRepo := new(MockPrintJobRepository)
	mockConfigRepo := new(MockPrinterConfigRepository)
	mockManager := new(MockPrinterManager)
	mockPrinter := new(MockPrinter)

	job := createTestPrintJobForWorker(printing.PrintJobStatusPending, 0)
	printerConfig := createTestPrinterConfig(printing.PrinterTypeBill)

	// Mock expectations
	mockJobRepo.On("UpdateStatus", mock.Anything, job.ID, printing.PrintJobStatusPrinting, "").Return(nil)
	mockConfigRepo.On("FindByID", mock.Anything, job.PrinterID).Return(printerConfig, nil)
	mockManager.On("GetPrinter", printerConfig).Return(mockPrinter, nil)
	mockPrinter.On("Connect").Return(nil)
	mockPrinter.On("Print", job.Content).Return(nil)
	mockPrinter.On("Disconnect").Return(nil)
	mockJobRepo.On("UpdateStatus", mock.Anything, job.ID, printing.PrintJobStatusCompleted, "").Return(nil)

	// Create worker
	worker := NewPrintWorker(PrintWorkerConfig{
		PrintJobRepo:      mockJobRepo,
		PrinterConfigRepo: mockConfigRepo,
		PrinterManager:    mockManager,
	})

	// Execute
	err := worker.ProcessJob(context.Background(), job)

	// Assert
	assert.NoError(t, err)
	mockJobRepo.AssertExpectations(t)
	mockConfigRepo.AssertExpectations(t)
	mockManager.AssertExpectations(t)
	mockPrinter.AssertExpectations(t)
}

func TestProcessJob_PrinterNotFound(t *testing.T) {
	// Setup
	mockJobRepo := new(MockPrintJobRepository)
	mockConfigRepo := new(MockPrinterConfigRepository)
	mockManager := new(MockPrinterManager)

	job := createTestPrintJobForWorker(printing.PrintJobStatusPending, 0)

	// Mock expectations
	mockJobRepo.On("UpdateStatus", mock.Anything, job.ID, printing.PrintJobStatusPrinting, "").Return(nil)
	mockConfigRepo.On("FindByID", mock.Anything, job.PrinterID).Return(nil, nil)
	mockJobRepo.On("IncrementRetry", mock.Anything, job.ID).Return(nil)
	mockJobRepo.On("UpdateStatus", mock.Anything, job.ID, printing.PrintJobStatusPending, mock.Anything).Return(nil)

	// Create worker
	worker := NewPrintWorker(PrintWorkerConfig{
		PrintJobRepo:      mockJobRepo,
		PrinterConfigRepo: mockConfigRepo,
		PrinterManager:    mockManager,
	})

	// Execute
	err := worker.ProcessJob(context.Background(), job)

	// Assert
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "printer not found")
	mockJobRepo.AssertExpectations(t)
	mockConfigRepo.AssertExpectations(t)
}

func TestProcessJob_PrinterDisabled(t *testing.T) {
	// Setup
	mockJobRepo := new(MockPrintJobRepository)
	mockConfigRepo := new(MockPrinterConfigRepository)
	mockManager := new(MockPrinterManager)

	job := createTestPrintJobForWorker(printing.PrintJobStatusPending, 0)
	printerConfig := createTestPrinterConfig(printing.PrinterTypeBill)
	printerConfig.IsEnabled = false

	// Mock expectations
	mockJobRepo.On("UpdateStatus", mock.Anything, job.ID, printing.PrintJobStatusPrinting, "").Return(nil)
	mockConfigRepo.On("FindByID", mock.Anything, job.PrinterID).Return(printerConfig, nil)
	mockJobRepo.On("IncrementRetry", mock.Anything, job.ID).Return(nil)
	mockJobRepo.On("UpdateStatus", mock.Anything, job.ID, printing.PrintJobStatusPending, mock.Anything).Return(nil)

	// Create worker
	worker := NewPrintWorker(PrintWorkerConfig{
		PrintJobRepo:      mockJobRepo,
		PrinterConfigRepo: mockConfigRepo,
		PrinterManager:    mockManager,
	})

	// Execute
	err := worker.ProcessJob(context.Background(), job)

	// Assert
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "printer is disabled")
	mockJobRepo.AssertExpectations(t)
	mockConfigRepo.AssertExpectations(t)
}

func TestProcessJob_ConnectionFailed(t *testing.T) {
	// Setup
	mockJobRepo := new(MockPrintJobRepository)
	mockConfigRepo := new(MockPrinterConfigRepository)
	mockManager := new(MockPrinterManager)
	mockPrinter := new(MockPrinter)

	job := createTestPrintJobForWorker(printing.PrintJobStatusPending, 0)
	printerConfig := createTestPrinterConfig(printing.PrinterTypeBill)

	// Mock expectations
	mockJobRepo.On("UpdateStatus", mock.Anything, job.ID, printing.PrintJobStatusPrinting, "").Return(nil)
	mockConfigRepo.On("FindByID", mock.Anything, job.PrinterID).Return(printerConfig, nil)
	mockManager.On("GetPrinter", printerConfig).Return(mockPrinter, nil)
	mockPrinter.On("Connect").Return(errors.New("connection refused"))
	mockJobRepo.On("IncrementRetry", mock.Anything, job.ID).Return(nil)
	mockJobRepo.On("UpdateStatus", mock.Anything, job.ID, printing.PrintJobStatusPending, mock.Anything).Return(nil)

	// Create worker
	worker := NewPrintWorker(PrintWorkerConfig{
		PrintJobRepo:      mockJobRepo,
		PrinterConfigRepo: mockConfigRepo,
		PrinterManager:    mockManager,
	})

	// Execute
	err := worker.ProcessJob(context.Background(), job)

	// Assert
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "failed to connect to printer")
	mockJobRepo.AssertExpectations(t)
	mockConfigRepo.AssertExpectations(t)
	mockManager.AssertExpectations(t)
	mockPrinter.AssertExpectations(t)
}

func TestProcessJob_PrintFailed(t *testing.T) {
	// Setup
	mockJobRepo := new(MockPrintJobRepository)
	mockConfigRepo := new(MockPrinterConfigRepository)
	mockManager := new(MockPrinterManager)
	mockPrinter := new(MockPrinter)

	job := createTestPrintJobForWorker(printing.PrintJobStatusPending, 0)
	printerConfig := createTestPrinterConfig(printing.PrinterTypeBill)

	// Mock expectations
	mockJobRepo.On("UpdateStatus", mock.Anything, job.ID, printing.PrintJobStatusPrinting, "").Return(nil)
	mockConfigRepo.On("FindByID", mock.Anything, job.PrinterID).Return(printerConfig, nil)
	mockManager.On("GetPrinter", printerConfig).Return(mockPrinter, nil)
	mockPrinter.On("Connect").Return(nil)
	mockPrinter.On("Print", job.Content).Return(errors.New("paper jam"))
	mockPrinter.On("Disconnect").Return(nil)
	mockJobRepo.On("IncrementRetry", mock.Anything, job.ID).Return(nil)
	mockJobRepo.On("UpdateStatus", mock.Anything, job.ID, printing.PrintJobStatusPending, mock.Anything).Return(nil)

	// Create worker
	worker := NewPrintWorker(PrintWorkerConfig{
		PrintJobRepo:      mockJobRepo,
		PrinterConfigRepo: mockConfigRepo,
		PrinterManager:    mockManager,
	})

	// Execute
	err := worker.ProcessJob(context.Background(), job)

	// Assert
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "failed to print")
	mockJobRepo.AssertExpectations(t)
	mockConfigRepo.AssertExpectations(t)
	mockManager.AssertExpectations(t)
	mockPrinter.AssertExpectations(t)
}

func TestProcessJob_MaxRetriesExceeded(t *testing.T) {
	// Setup
	mockJobRepo := new(MockPrintJobRepository)
	mockConfigRepo := new(MockPrinterConfigRepository)
	mockManager := new(MockPrinterManager)
	mockPrinter := new(MockPrinter)

	// Job that has already failed 3 times (max retries)
	job := createTestPrintJobForWorker(printing.PrintJobStatusPending, 3)
	printerConfig := createTestPrinterConfig(printing.PrinterTypeBill)

	// Mock expectations
	mockJobRepo.On("UpdateStatus", mock.Anything, job.ID, printing.PrintJobStatusPrinting, "").Return(nil)
	mockConfigRepo.On("FindByID", mock.Anything, job.PrinterID).Return(printerConfig, nil)
	mockManager.On("GetPrinter", printerConfig).Return(mockPrinter, nil)
	mockPrinter.On("Connect").Return(errors.New("connection refused"))
	mockJobRepo.On("IncrementRetry", mock.Anything, job.ID).Return(nil)
	mockJobRepo.On("UpdateStatus", mock.Anything, job.ID, printing.PrintJobStatusFailed, mock.Anything).Return(nil)

	// Create worker
	worker := NewPrintWorker(PrintWorkerConfig{
		PrintJobRepo:      mockJobRepo,
		PrinterConfigRepo: mockConfigRepo,
		PrinterManager:    mockManager,
	})

	// Execute
	err := worker.ProcessJob(context.Background(), job)

	// Assert
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "job failed after")
	mockJobRepo.AssertExpectations(t)
	mockConfigRepo.AssertExpectations(t)
	mockManager.AssertExpectations(t)
	mockPrinter.AssertExpectations(t)
}

func TestProcessJob_NotPendingStatus(t *testing.T) {
	// Setup
	mockJobRepo := new(MockPrintJobRepository)
	mockConfigRepo := new(MockPrinterConfigRepository)
	mockManager := new(MockPrinterManager)

	// Job with COMPLETED status
	job := createTestPrintJobForWorker(printing.PrintJobStatusCompleted, 0)

	// Create worker
	worker := NewPrintWorker(PrintWorkerConfig{
		PrintJobRepo:      mockJobRepo,
		PrinterConfigRepo: mockConfigRepo,
		PrinterManager:    mockManager,
	})

	// Execute
	err := worker.ProcessJob(context.Background(), job)

	// Assert
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "is not pending")
}

func TestProcessJob_NilJob(t *testing.T) {
	// Setup
	mockJobRepo := new(MockPrintJobRepository)
	mockConfigRepo := new(MockPrinterConfigRepository)
	mockManager := new(MockPrinterManager)

	// Create worker
	worker := NewPrintWorker(PrintWorkerConfig{
		PrintJobRepo:      mockJobRepo,
		PrinterConfigRepo: mockConfigRepo,
		PrinterManager:    mockManager,
	})

	// Execute
	err := worker.ProcessJob(context.Background(), nil)

	// Assert
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "job cannot be nil")
}

func TestCalculateBackoffDelay(t *testing.T) {
	// Setup
	mockJobRepo := new(MockPrintJobRepository)
	mockConfigRepo := new(MockPrinterConfigRepository)
	mockManager := new(MockPrinterManager)

	worker := NewPrintWorker(PrintWorkerConfig{
		PrintJobRepo:      mockJobRepo,
		PrinterConfigRepo: mockConfigRepo,
		PrinterManager:    mockManager,
	})

	w := worker.(*printWorker)

	// Test cases
	tests := []struct {
		name        string
		retryCount  int
		expected    time.Duration
		description string
	}{
		{
			name:        "No retry",
			retryCount:  0,
			expected:    0,
			description: "First attempt should have no delay",
		},
		{
			name:        "First retry",
			retryCount:  1,
			expected:    30 * time.Second,
			description: "First retry should wait 30 seconds",
		},
		{
			name:        "Second retry",
			retryCount:  2,
			expected:    60 * time.Second,
			description: "Second retry should wait 1 minute",
		},
		{
			name:        "Third retry",
			retryCount:  3,
			expected:    120 * time.Second,
			description: "Third retry should wait 2 minutes",
		},
		{
			name:        "Fourth retry (capped)",
			retryCount:  4,
			expected:    240 * time.Second,
			description: "Fourth retry should wait 4 minutes",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			delay := w.calculateBackoffDelay(tt.retryCount)
			assert.Equal(t, tt.expected, delay, tt.description)
		})
	}
}

func TestProcessJob_RetryBackoff(t *testing.T) {
	// Setup
	mockJobRepo := new(MockPrintJobRepository)
	mockConfigRepo := new(MockPrinterConfigRepository)
	mockManager := new(MockPrinterManager)

	// Job that was just updated (should not retry yet)
	job := createTestPrintJobForWorker(printing.PrintJobStatusPending, 1)
	job.UpdatedAt = time.Now() // Just updated

	// Create worker
	worker := NewPrintWorker(PrintWorkerConfig{
		PrintJobRepo:      mockJobRepo,
		PrinterConfigRepo: mockConfigRepo,
		PrinterManager:    mockManager,
	})

	// Execute
	err := worker.ProcessJob(context.Background(), job)

	// Assert - should return nil (not ready to retry yet)
	assert.NoError(t, err)
}

func TestProcessJob_RetryAfterBackoff(t *testing.T) {
	// Setup
	mockJobRepo := new(MockPrintJobRepository)
	mockConfigRepo := new(MockPrinterConfigRepository)
	mockManager := new(MockPrinterManager)
	mockPrinter := new(MockPrinter)

	// Job that was updated 31 seconds ago (ready for retry)
	job := createTestPrintJobForWorker(printing.PrintJobStatusPending, 1)
	job.UpdatedAt = time.Now().Add(-31 * time.Second)
	printerConfig := createTestPrinterConfig(printing.PrinterTypeBill)

	// Mock expectations
	mockJobRepo.On("UpdateStatus", mock.Anything, job.ID, printing.PrintJobStatusPrinting, "").Return(nil)
	mockConfigRepo.On("FindByID", mock.Anything, job.PrinterID).Return(printerConfig, nil)
	mockManager.On("GetPrinter", printerConfig).Return(mockPrinter, nil)
	mockPrinter.On("Connect").Return(nil)
	mockPrinter.On("Print", job.Content).Return(nil)
	mockPrinter.On("Disconnect").Return(nil)
	mockJobRepo.On("UpdateStatus", mock.Anything, job.ID, printing.PrintJobStatusCompleted, "").Return(nil)

	// Create worker
	worker := NewPrintWorker(PrintWorkerConfig{
		PrintJobRepo:      mockJobRepo,
		PrinterConfigRepo: mockConfigRepo,
		PrinterManager:    mockManager,
	})

	// Execute
	err := worker.ProcessJob(context.Background(), job)

	// Assert
	assert.NoError(t, err)
	mockJobRepo.AssertExpectations(t)
	mockConfigRepo.AssertExpectations(t)
	mockManager.AssertExpectations(t)
	mockPrinter.AssertExpectations(t)
}

func TestStop(t *testing.T) {
	// Setup
	mockJobRepo := new(MockPrintJobRepository)
	mockConfigRepo := new(MockPrinterConfigRepository)
	mockManager := new(MockPrinterManager)

	worker := NewPrintWorker(PrintWorkerConfig{
		PrintJobRepo:      mockJobRepo,
		PrinterConfigRepo: mockConfigRepo,
		PrinterManager:    mockManager,
	})

	// Execute
	worker.Stop()

	// Assert - should not panic
	// The stop channel should be closed
	w := worker.(*printWorker)
	select {
	case <-w.stopChan:
		// Channel is closed, test passes
	default:
		t.Error("Stop channel should be closed")
	}
}
