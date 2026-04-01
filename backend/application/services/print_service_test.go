package services

import (
	"context"
	"errors"
	"testing"
	"time"

	"cafe-pos/backend/domain/order"
	"cafe-pos/backend/domain/printing"
	"cafe-pos/backend/domain/settings"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

// Mock repositories
type MockPrintJobRepository struct {
	mock.Mock
}

func (m *MockPrintJobRepository) Create(ctx context.Context, job *printing.PrintJob) error {
	args := m.Called(ctx, job)
	return args.Error(0)
}

func (m *MockPrintJobRepository) FindByID(ctx context.Context, id primitive.ObjectID) (*printing.PrintJob, error) {
	args := m.Called(ctx, id)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*printing.PrintJob), args.Error(1)
}

func (m *MockPrintJobRepository) FindByOrderID(ctx context.Context, orderID primitive.ObjectID) ([]*printing.PrintJob, error) {
	args := m.Called(ctx, orderID)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]*printing.PrintJob), args.Error(1)
}

func (m *MockPrintJobRepository) FindPending(ctx context.Context, limit int) ([]*printing.PrintJob, error) {
	args := m.Called(ctx, limit)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]*printing.PrintJob), args.Error(1)
}

func (m *MockPrintJobRepository) FindFailed(ctx context.Context) ([]*printing.PrintJob, error) {
	args := m.Called(ctx)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]*printing.PrintJob), args.Error(1)
}

func (m *MockPrintJobRepository) UpdateStatus(ctx context.Context, id primitive.ObjectID, status printing.PrintJobStatus, errorMsg string) error {
	args := m.Called(ctx, id, status, errorMsg)
	return args.Error(0)
}

func (m *MockPrintJobRepository) IncrementRetry(ctx context.Context, id primitive.ObjectID) error {
	args := m.Called(ctx, id)
	return args.Error(0)
}

func (m *MockPrintJobRepository) Delete(ctx context.Context, id primitive.ObjectID) error {
	args := m.Called(ctx, id)
	return args.Error(0)
}

func (m *MockPrintJobRepository) DeleteOldCompleted(ctx context.Context, olderThan time.Time) error {
	args := m.Called(ctx, olderThan)
	return args.Error(0)
}

type MockPrinterConfigRepository struct {
	mock.Mock
}

func (m *MockPrinterConfigRepository) Create(ctx context.Context, config *printing.PrinterConfig) error {
	args := m.Called(ctx, config)
	return args.Error(0)
}

func (m *MockPrinterConfigRepository) FindByID(ctx context.Context, id primitive.ObjectID) (*printing.PrinterConfig, error) {
	args := m.Called(ctx, id)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*printing.PrinterConfig), args.Error(1)
}

func (m *MockPrinterConfigRepository) FindAll(ctx context.Context) ([]*printing.PrinterConfig, error) {
	args := m.Called(ctx)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]*printing.PrinterConfig), args.Error(1)
}

func (m *MockPrinterConfigRepository) FindByType(ctx context.Context, printerType printing.PrinterType) ([]*printing.PrinterConfig, error) {
	args := m.Called(ctx, printerType)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]*printing.PrinterConfig), args.Error(1)
}

func (m *MockPrinterConfigRepository) FindDefault(ctx context.Context, printerType printing.PrinterType) (*printing.PrinterConfig, error) {
	args := m.Called(ctx, printerType)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*printing.PrinterConfig), args.Error(1)
}

func (m *MockPrinterConfigRepository) Update(ctx context.Context, config *printing.PrinterConfig) error {
	args := m.Called(ctx, config)
	return args.Error(0)
}

func (m *MockPrinterConfigRepository) Delete(ctx context.Context, id primitive.ObjectID) error {
	args := m.Called(ctx, id)
	return args.Error(0)
}

type MockPrintTemplateRepository struct {
	mock.Mock
}

func (m *MockPrintTemplateRepository) Create(ctx context.Context, template *printing.PrintTemplate) error {
	args := m.Called(ctx, template)
	return args.Error(0)
}

func (m *MockPrintTemplateRepository) FindByID(ctx context.Context, id primitive.ObjectID) (*printing.PrintTemplate, error) {
	args := m.Called(ctx, id)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*printing.PrintTemplate), args.Error(1)
}

func (m *MockPrintTemplateRepository) FindByType(ctx context.Context, templateType printing.TemplateType) ([]*printing.PrintTemplate, error) {
	args := m.Called(ctx, templateType)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]*printing.PrintTemplate), args.Error(1)
}

func (m *MockPrintTemplateRepository) FindDefault(ctx context.Context, templateType printing.TemplateType) (*printing.PrintTemplate, error) {
	args := m.Called(ctx, templateType)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*printing.PrintTemplate), args.Error(1)
}

func (m *MockPrintTemplateRepository) Update(ctx context.Context, template *printing.PrintTemplate) error {
	args := m.Called(ctx, template)
	return args.Error(0)
}

func (m *MockPrintTemplateRepository) Delete(ctx context.Context, id primitive.ObjectID) error {
	args := m.Called(ctx, id)
	return args.Error(0)
}

type MockTemplateRenderer struct {
	mock.Mock
}

func (m *MockTemplateRenderer) RenderBill(ord *order.Order, tmpl *printing.PrintTemplate, shopSettings *settings.ShopSettings) (string, error) {
	args := m.Called(ord, tmpl, shopSettings)
	return args.String(0), args.Error(1)
}

func (m *MockTemplateRenderer) RenderLabel(ord *order.Order, itemIndex int, tmpl *printing.PrintTemplate, shopSettings *settings.ShopSettings) (string, error) {
	args := m.Called(ord, itemIndex, tmpl, shopSettings)
	return args.String(0), args.Error(1)
}

type MockOrderRepository struct {
	mock.Mock
}

func (m *MockOrderRepository) Create(ctx context.Context, o *order.Order) error {
	args := m.Called(ctx, o)
	return args.Error(0)
}

func (m *MockOrderRepository) FindByID(ctx context.Context, id primitive.ObjectID) (*order.Order, error) {
	args := m.Called(ctx, id)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*order.Order), args.Error(1)
}

func (m *MockOrderRepository) Update(ctx context.Context, id primitive.ObjectID, o *order.Order) error {
	args := m.Called(ctx, id, o)
	return args.Error(0)
}

func (m *MockOrderRepository) Delete(ctx context.Context, id primitive.ObjectID) error {
	args := m.Called(ctx, id)
	return args.Error(0)
}

func (m *MockOrderRepository) FindByShiftID(ctx context.Context, shiftID primitive.ObjectID) ([]*order.Order, error) {
	args := m.Called(ctx, shiftID)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]*order.Order), args.Error(1)
}

func (m *MockOrderRepository) FindByWaiterID(ctx context.Context, waiterID primitive.ObjectID) ([]*order.Order, error) {
	args := m.Called(ctx, waiterID)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]*order.Order), args.Error(1)
}

func (m *MockOrderRepository) FindByStatus(ctx context.Context, status order.OrderStatus) ([]*order.Order, error) {
	args := m.Called(ctx, status)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]*order.Order), args.Error(1)
}

func (m *MockOrderRepository) FindByOrderNumber(ctx context.Context, orderNumber string) (*order.Order, error) {
	args := m.Called(ctx, orderNumber)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*order.Order), args.Error(1)
}

func (m *MockOrderRepository) FindAll(ctx context.Context) ([]*order.Order, error) {
	args := m.Called(ctx)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]*order.Order), args.Error(1)
}

func (m *MockOrderRepository) FindByIDs(ctx context.Context, ids []primitive.ObjectID) ([]*order.Order, error) {
	return nil, nil
}

// Helper functions
func createTestOrder() *order.Order {
	return &order.Order{
		ID:          primitive.NewObjectID(),
		OrderNumber: "ORD-001",
		Items: []order.OrderItem{
			{
				Name:     "Cà phê sữa",
				Quantity: 1,
				Price:    25000,
				Subtotal: 25000,
			},
			{
				Name:        "Trà sữa",
				VariantName: "Size L",
				Quantity:    2,
				Price:       30000,
				Subtotal:    60000,
			},
		},
		Subtotal:  85000,
		Total:     85000,
		Status:    order.StatusPaid,
		CreatedAt: time.Now(),
	}
}

func createTestPrinterConfig(printerType printing.PrinterType) *printing.PrinterConfig {
	return &printing.PrinterConfig{
		ID:             primitive.NewObjectID(),
		Name:           "Test Printer",
		Type:           printerType,
		ConnectionType: printing.ConnectionTypeNetwork,
		IPAddress:      "192.168.1.100",
		Port:           9100,
		PaperWidth:     80,
		IsDefault:      true,
		IsEnabled:      true,
	}
}

func createTestTemplate(templateType printing.TemplateType) *printing.PrintTemplate {
	return &printing.PrintTemplate{
		ID:        primitive.NewObjectID(),
		Type:      templateType,
		Name:      "Default Template",
		Content:   "Test template content",
		IsDefault: true,
	}
}

// Tests
func TestCreatePrintJobsForOrder_Success(t *testing.T) {
	// Setup
	ctx := context.Background()
	ord := createTestOrder()
	billPrinter := createTestPrinterConfig(printing.PrinterTypeBill)
	labelPrinter := createTestPrinterConfig(printing.PrinterTypeLabel)
	billTemplate := createTestTemplate(printing.TemplateTypeBill)
	labelTemplate := createTestTemplate(printing.TemplateTypeLabel)

	mockPrintJobRepo := new(MockPrintJobRepository)
	mockPrinterConfigRepo := new(MockPrinterConfigRepository)
	mockTemplateRepo := new(MockPrintTemplateRepository)
	mockRenderer := new(MockTemplateRenderer)
	mockOrderRepo := new(MockOrderRepository)

	// Expectations
	mockPrinterConfigRepo.On("FindDefault", ctx, printing.PrinterTypeBill).Return(billPrinter, nil)
	mockPrinterConfigRepo.On("FindDefault", ctx, printing.PrinterTypeLabel).Return(labelPrinter, nil)
	mockTemplateRepo.On("FindDefault", ctx, printing.TemplateTypeBill).Return(billTemplate, nil)
	mockTemplateRepo.On("FindDefault", ctx, printing.TemplateTypeLabel).Return(labelTemplate, nil)
	mockRenderer.On("RenderBill", ord, billTemplate, mock.Anything).Return("Bill content", nil)
	mockRenderer.On("RenderLabel", ord, 0, labelTemplate, mock.Anything).Return("Label 1 content", nil)
	mockRenderer.On("RenderLabel", ord, 1, labelTemplate, mock.Anything).Return("Label 2 content", nil)
	mockPrintJobRepo.On("Create", ctx, mock.MatchedBy(func(job *printing.PrintJob) bool {
		return job.Type == printing.PrintJobTypeBill && job.Status == printing.PrintJobStatusPending
	})).Return(nil)
	mockPrintJobRepo.On("Create", ctx, mock.MatchedBy(func(job *printing.PrintJob) bool {
		return job.Type == printing.PrintJobTypeLabel && job.Status == printing.PrintJobStatusPending
	})).Return(nil)

	service := NewPrintService(PrintServiceConfig{
		PrintJobRepo:      mockPrintJobRepo,
		PrinterConfigRepo: mockPrinterConfigRepo,
		TemplateRepo:      mockTemplateRepo,
		TemplateRenderer:  mockRenderer,
		OrderRepo:         mockOrderRepo,
	})

	// Execute
	err := service.CreatePrintJobsForOrder(ctx, ord)

	// Assert
	assert.NoError(t, err)
	mockPrintJobRepo.AssertNumberOfCalls(t, "Create", 3) // 1 bill + 2 labels
	mockPrinterConfigRepo.AssertExpectations(t)
	mockTemplateRepo.AssertExpectations(t)
	mockRenderer.AssertExpectations(t)
}

func TestCreatePrintJobsForOrder_NilOrder(t *testing.T) {
	// Setup
	ctx := context.Background()
	service := NewPrintService(PrintServiceConfig{
		PrintJobRepo:      new(MockPrintJobRepository),
		PrinterConfigRepo: new(MockPrinterConfigRepository),
		TemplateRepo:      new(MockPrintTemplateRepository),
		TemplateRenderer:  new(MockTemplateRenderer),
		OrderRepo:         new(MockOrderRepository),
	})

	// Execute
	err := service.CreatePrintJobsForOrder(ctx, nil)

	// Assert
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "order cannot be nil")
}

func TestCreatePrintJobsForOrder_NoDefaultBillPrinter(t *testing.T) {
	// Setup
	ctx := context.Background()
	ord := createTestOrder()
	mockPrinterConfigRepo := new(MockPrinterConfigRepository)

	mockPrinterConfigRepo.On("FindDefault", ctx, printing.PrinterTypeBill).Return(nil, nil)

	service := NewPrintService(PrintServiceConfig{
		PrintJobRepo:      new(MockPrintJobRepository),
		PrinterConfigRepo: mockPrinterConfigRepo,
		TemplateRepo:      new(MockPrintTemplateRepository),
		TemplateRenderer:  new(MockTemplateRenderer),
		OrderRepo:         new(MockOrderRepository),
	})

	// Execute
	err := service.CreatePrintJobsForOrder(ctx, ord)

	// Assert
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "no default bill printer configured")
}

func TestReprintBill_Success(t *testing.T) {
	// Setup
	ctx := context.Background()
	ord := createTestOrder()
	billPrinter := createTestPrinterConfig(printing.PrinterTypeBill)
	billTemplate := createTestTemplate(printing.TemplateTypeBill)

	mockPrintJobRepo := new(MockPrintJobRepository)
	mockPrinterConfigRepo := new(MockPrinterConfigRepository)
	mockTemplateRepo := new(MockPrintTemplateRepository)
	mockRenderer := new(MockTemplateRenderer)
	mockOrderRepo := new(MockOrderRepository)

	// Expectations
	mockOrderRepo.On("FindByID", ctx, ord.ID).Return(ord, nil)
	mockPrinterConfigRepo.On("FindDefault", ctx, printing.PrinterTypeBill).Return(billPrinter, nil)
	mockTemplateRepo.On("FindDefault", ctx, printing.TemplateTypeBill).Return(billTemplate, nil)
	mockRenderer.On("RenderBill", ord, billTemplate, mock.Anything).Return("Bill content", nil)
	mockPrintJobRepo.On("Create", ctx, mock.MatchedBy(func(job *printing.PrintJob) bool {
		return job.Type == printing.PrintJobTypeBill && job.Status == printing.PrintJobStatusPending
	})).Return(nil)

	service := NewPrintService(PrintServiceConfig{
		PrintJobRepo:      mockPrintJobRepo,
		PrinterConfigRepo: mockPrinterConfigRepo,
		TemplateRepo:      mockTemplateRepo,
		TemplateRenderer:  mockRenderer,
		OrderRepo:         mockOrderRepo,
	})

	// Execute
	err := service.ReprintBill(ctx, ord.ID)

	// Assert
	assert.NoError(t, err)
	mockOrderRepo.AssertExpectations(t)
	mockPrintJobRepo.AssertExpectations(t)
}

func TestReprintBill_OrderNotFound(t *testing.T) {
	// Setup
	ctx := context.Background()
	orderID := primitive.NewObjectID()
	mockOrderRepo := new(MockOrderRepository)

	mockOrderRepo.On("FindByID", ctx, orderID).Return(nil, nil)

	service := NewPrintService(PrintServiceConfig{
		PrintJobRepo:      new(MockPrintJobRepository),
		PrinterConfigRepo: new(MockPrinterConfigRepository),
		TemplateRepo:      new(MockPrintTemplateRepository),
		TemplateRenderer:  new(MockTemplateRenderer),
		OrderRepo:         mockOrderRepo,
	})

	// Execute
	err := service.ReprintBill(ctx, orderID)

	// Assert
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "order not found")
}

func TestReprintLabel_Success(t *testing.T) {
	// Setup
	ctx := context.Background()
	ord := createTestOrder()
	labelPrinter := createTestPrinterConfig(printing.PrinterTypeLabel)
	labelTemplate := createTestTemplate(printing.TemplateTypeLabel)

	mockPrintJobRepo := new(MockPrintJobRepository)
	mockPrinterConfigRepo := new(MockPrinterConfigRepository)
	mockTemplateRepo := new(MockPrintTemplateRepository)
	mockRenderer := new(MockTemplateRenderer)
	mockOrderRepo := new(MockOrderRepository)

	// Expectations
	mockOrderRepo.On("FindByID", ctx, ord.ID).Return(ord, nil)
	mockPrinterConfigRepo.On("FindDefault", ctx, printing.PrinterTypeLabel).Return(labelPrinter, nil)
	mockTemplateRepo.On("FindDefault", ctx, printing.TemplateTypeLabel).Return(labelTemplate, nil)
	mockRenderer.On("RenderLabel", ord, 0, labelTemplate, mock.Anything).Return("Label content", nil)
	mockPrintJobRepo.On("Create", ctx, mock.MatchedBy(func(job *printing.PrintJob) bool {
		return job.Type == printing.PrintJobTypeLabel && job.Status == printing.PrintJobStatusPending
	})).Return(nil)

	service := NewPrintService(PrintServiceConfig{
		PrintJobRepo:      mockPrintJobRepo,
		PrinterConfigRepo: mockPrinterConfigRepo,
		TemplateRepo:      mockTemplateRepo,
		TemplateRenderer:  mockRenderer,
		OrderRepo:         mockOrderRepo,
	})

	// Execute
	err := service.ReprintLabel(ctx, ord.ID, 0)

	// Assert
	assert.NoError(t, err)
	mockOrderRepo.AssertExpectations(t)
	mockPrintJobRepo.AssertExpectations(t)
}

func TestReprintLabel_InvalidItemIndex(t *testing.T) {
	// Setup
	ctx := context.Background()
	ord := createTestOrder()
	mockOrderRepo := new(MockOrderRepository)

	mockOrderRepo.On("FindByID", ctx, ord.ID).Return(ord, nil)

	service := NewPrintService(PrintServiceConfig{
		PrintJobRepo:      new(MockPrintJobRepository),
		PrinterConfigRepo: new(MockPrinterConfigRepository),
		TemplateRepo:      new(MockPrintTemplateRepository),
		TemplateRenderer:  new(MockTemplateRenderer),
		OrderRepo:         mockOrderRepo,
	})

	// Execute
	err := service.ReprintLabel(ctx, ord.ID, 10) // Invalid index

	// Assert
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "invalid item index")
}

func TestGetPendingJobs_Success(t *testing.T) {
	// Setup
	ctx := context.Background()
	expectedJobs := []*printing.PrintJob{
		{ID: primitive.NewObjectID(), Status: printing.PrintJobStatusPending},
		{ID: primitive.NewObjectID(), Status: printing.PrintJobStatusPending},
	}

	mockPrintJobRepo := new(MockPrintJobRepository)
	mockPrintJobRepo.On("FindPending", ctx, 100).Return(expectedJobs, nil)

	service := NewPrintService(PrintServiceConfig{
		PrintJobRepo:      mockPrintJobRepo,
		PrinterConfigRepo: new(MockPrinterConfigRepository),
		TemplateRepo:      new(MockPrintTemplateRepository),
		TemplateRenderer:  new(MockTemplateRenderer),
		OrderRepo:         new(MockOrderRepository),
	})

	// Execute
	jobs, err := service.GetPendingJobs(ctx)

	// Assert
	assert.NoError(t, err)
	assert.Equal(t, expectedJobs, jobs)
	mockPrintJobRepo.AssertExpectations(t)
}

func TestGetFailedJobs_Success(t *testing.T) {
	// Setup
	ctx := context.Background()
	expectedJobs := []*printing.PrintJob{
		{ID: primitive.NewObjectID(), Status: printing.PrintJobStatusFailed},
	}

	mockPrintJobRepo := new(MockPrintJobRepository)
	mockPrintJobRepo.On("FindFailed", ctx).Return(expectedJobs, nil)

	service := NewPrintService(PrintServiceConfig{
		PrintJobRepo:      mockPrintJobRepo,
		PrinterConfigRepo: new(MockPrinterConfigRepository),
		TemplateRepo:      new(MockPrintTemplateRepository),
		TemplateRenderer:  new(MockTemplateRenderer),
		OrderRepo:         new(MockOrderRepository),
	})

	// Execute
	jobs, err := service.GetFailedJobs(ctx)

	// Assert
	assert.NoError(t, err)
	assert.Equal(t, expectedJobs, jobs)
	mockPrintJobRepo.AssertExpectations(t)
}

func TestRetryJob_Success(t *testing.T) {
	// Setup
	ctx := context.Background()
	jobID := primitive.NewObjectID()
	failedJob := &printing.PrintJob{
		ID:         jobID,
		Status:     printing.PrintJobStatusFailed,
		RetryCount: 3,
		ErrorMsg:   "Connection failed",
	}

	mockPrintJobRepo := new(MockPrintJobRepository)
	mockPrintJobRepo.On("FindByID", ctx, jobID).Return(failedJob, nil)
	mockPrintJobRepo.On("UpdateStatus", ctx, jobID, printing.PrintJobStatusPending, "").Return(nil)

	service := NewPrintService(PrintServiceConfig{
		PrintJobRepo:      mockPrintJobRepo,
		PrinterConfigRepo: new(MockPrinterConfigRepository),
		TemplateRepo:      new(MockPrintTemplateRepository),
		TemplateRenderer:  new(MockTemplateRenderer),
		OrderRepo:         new(MockOrderRepository),
	})

	// Execute
	err := service.RetryJob(ctx, jobID)

	// Assert
	assert.NoError(t, err)
	mockPrintJobRepo.AssertExpectations(t)
}

func TestRetryJob_NotFailedStatus(t *testing.T) {
	// Setup
	ctx := context.Background()
	jobID := primitive.NewObjectID()
	pendingJob := &printing.PrintJob{
		ID:     jobID,
		Status: printing.PrintJobStatusPending,
	}

	mockPrintJobRepo := new(MockPrintJobRepository)
	mockPrintJobRepo.On("FindByID", ctx, jobID).Return(pendingJob, nil)

	service := NewPrintService(PrintServiceConfig{
		PrintJobRepo:      mockPrintJobRepo,
		PrinterConfigRepo: new(MockPrinterConfigRepository),
		TemplateRepo:      new(MockPrintTemplateRepository),
		TemplateRenderer:  new(MockTemplateRenderer),
		OrderRepo:         new(MockOrderRepository),
	})

	// Execute
	err := service.RetryJob(ctx, jobID)

	// Assert
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "can only retry failed jobs")
}

func TestCancelJob_Success(t *testing.T) {
	// Setup
	ctx := context.Background()
	jobID := primitive.NewObjectID()
	pendingJob := &printing.PrintJob{
		ID:     jobID,
		Status: printing.PrintJobStatusPending,
	}

	mockPrintJobRepo := new(MockPrintJobRepository)
	mockPrintJobRepo.On("FindByID", ctx, jobID).Return(pendingJob, nil)
	mockPrintJobRepo.On("Delete", ctx, jobID).Return(nil)

	service := NewPrintService(PrintServiceConfig{
		PrintJobRepo:      mockPrintJobRepo,
		PrinterConfigRepo: new(MockPrinterConfigRepository),
		TemplateRepo:      new(MockPrintTemplateRepository),
		TemplateRenderer:  new(MockTemplateRenderer),
		OrderRepo:         new(MockOrderRepository),
	})

	// Execute
	err := service.CancelJob(ctx, jobID)

	// Assert
	assert.NoError(t, err)
	mockPrintJobRepo.AssertExpectations(t)
}

func TestCancelJob_NotPendingStatus(t *testing.T) {
	// Setup
	ctx := context.Background()
	jobID := primitive.NewObjectID()
	completedJob := &printing.PrintJob{
		ID:     jobID,
		Status: printing.PrintJobStatusCompleted,
	}

	mockPrintJobRepo := new(MockPrintJobRepository)
	mockPrintJobRepo.On("FindByID", ctx, jobID).Return(completedJob, nil)

	service := NewPrintService(PrintServiceConfig{
		PrintJobRepo:      mockPrintJobRepo,
		PrinterConfigRepo: new(MockPrinterConfigRepository),
		TemplateRepo:      new(MockPrintTemplateRepository),
		TemplateRenderer:  new(MockTemplateRenderer),
		OrderRepo:         new(MockOrderRepository),
	})

	// Execute
	err := service.CancelJob(ctx, jobID)

	// Assert
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "can only cancel pending jobs")
}

func TestCancelJob_JobNotFound(t *testing.T) {
	// Setup
	ctx := context.Background()
	jobID := primitive.NewObjectID()

	mockPrintJobRepo := new(MockPrintJobRepository)
	mockPrintJobRepo.On("FindByID", ctx, jobID).Return(nil, nil)

	service := NewPrintService(PrintServiceConfig{
		PrintJobRepo:      mockPrintJobRepo,
		PrinterConfigRepo: new(MockPrinterConfigRepository),
		TemplateRepo:      new(MockPrintTemplateRepository),
		TemplateRenderer:  new(MockTemplateRenderer),
		OrderRepo:         new(MockOrderRepository),
	})

	// Execute
	err := service.CancelJob(ctx, jobID)

	// Assert
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "job not found")
}

func TestCreatePrintJobsForOrder_TemplateRenderingError(t *testing.T) {
	// Setup
	ctx := context.Background()
	ord := createTestOrder()
	billPrinter := createTestPrinterConfig(printing.PrinterTypeBill)
	labelPrinter := createTestPrinterConfig(printing.PrinterTypeLabel)
	billTemplate := createTestTemplate(printing.TemplateTypeBill)
	labelTemplate := createTestTemplate(printing.TemplateTypeLabel)

	mockPrintJobRepo := new(MockPrintJobRepository)
	mockPrinterConfigRepo := new(MockPrinterConfigRepository)
	mockTemplateRepo := new(MockPrintTemplateRepository)
	mockRenderer := new(MockTemplateRenderer)
	mockOrderRepo := new(MockOrderRepository)

	// Expectations - bill rendering fails
	mockPrinterConfigRepo.On("FindDefault", ctx, printing.PrinterTypeBill).Return(billPrinter, nil)
	mockPrinterConfigRepo.On("FindDefault", ctx, printing.PrinterTypeLabel).Return(labelPrinter, nil)
	mockTemplateRepo.On("FindDefault", ctx, printing.TemplateTypeBill).Return(billTemplate, nil)
	mockTemplateRepo.On("FindDefault", ctx, printing.TemplateTypeLabel).Return(labelTemplate, nil)
	mockRenderer.On("RenderBill", ord, billTemplate, mock.Anything).Return("", errors.New("template error"))

	service := NewPrintService(PrintServiceConfig{
		PrintJobRepo:      mockPrintJobRepo,
		PrinterConfigRepo: mockPrinterConfigRepo,
		TemplateRepo:      mockTemplateRepo,
		TemplateRenderer:  mockRenderer,
		OrderRepo:         mockOrderRepo,
	})

	// Execute
	err := service.CreatePrintJobsForOrder(ctx, ord)

	// Assert
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "failed to render bill")
	mockPrintJobRepo.AssertNotCalled(t, "Create")
}
