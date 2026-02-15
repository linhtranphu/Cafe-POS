package http

import (
	"context"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"cafe-pos/backend/application/services"
	"cafe-pos/backend/domain/batch"
	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

// Mock repositories for testing
type mockBatchRecordRepo struct {
	mock.Mock
}

func (m *mockBatchRecordRepo) Create(ctx context.Context, record *batch.BatchRecord) error {
	args := m.Called(ctx, record)
	return args.Error(0)
}

func (m *mockBatchRecordRepo) Update(ctx context.Context, record *batch.BatchRecord) error {
	args := m.Called(ctx, record)
	return args.Error(0)
}

func (m *mockBatchRecordRepo) Delete(ctx context.Context, id primitive.ObjectID) error {
	args := m.Called(ctx, id)
	return args.Error(0)
}

func (m *mockBatchRecordRepo) FindByID(ctx context.Context, id primitive.ObjectID) (*batch.BatchRecord, error) {
	args := m.Called(ctx, id)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*batch.BatchRecord), args.Error(1)
}

func (m *mockBatchRecordRepo) FindAll(ctx context.Context, filter batch.BatchRecordFilter) ([]*batch.BatchRecord, int64, error) {
	args := m.Called(ctx, filter)
	if args.Get(0) == nil {
		return nil, 0, args.Error(2)
	}
	return args.Get(0).([]*batch.BatchRecord), int64(args.Int(1)), args.Error(2)
}

func (m *mockBatchRecordRepo) FindAvailableByDefinition(ctx context.Context, defID primitive.ObjectID) ([]*batch.BatchRecord, error) {
	args := m.Called(ctx, defID)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]*batch.BatchRecord), args.Error(1)
}

func (m *mockBatchRecordRepo) UpdateQuantity(ctx context.Context, id primitive.ObjectID, newQuantity float64) error {
	args := m.Called(ctx, id, newQuantity)
	return args.Error(0)
}

func (m *mockBatchRecordRepo) GetTotalAvailableQuantity(ctx context.Context, defID primitive.ObjectID) (float64, error) {
	args := m.Called(ctx, defID)
	return args.Get(0).(float64), args.Error(1)
}

type mockBatchUsageLogRepo struct {
	mock.Mock
}

func (m *mockBatchUsageLogRepo) Create(ctx context.Context, log *batch.BatchUsageLog) error {
	args := m.Called(ctx, log)
	return args.Error(0)
}

func (m *mockBatchUsageLogRepo) FindAll(ctx context.Context, filter batch.BatchUsageLogFilter) ([]*batch.BatchUsageLog, int64, error) {
	args := m.Called(ctx, filter)
	if args.Get(0) == nil {
		return nil, 0, args.Error(2)
	}
	return args.Get(0).([]*batch.BatchUsageLog), int64(args.Int(1)), args.Error(2)
}

type mockBatchDefRepo struct {
	mock.Mock
}

func (m *mockBatchDefRepo) Create(ctx context.Context, def *batch.BatchDefinition) error {
	args := m.Called(ctx, def)
	return args.Error(0)
}

func (m *mockBatchDefRepo) Update(ctx context.Context, def *batch.BatchDefinition) error {
	args := m.Called(ctx, def)
	return args.Error(0)
}

func (m *mockBatchDefRepo) Delete(ctx context.Context, id primitive.ObjectID) error {
	args := m.Called(ctx, id)
	return args.Error(0)
}

func (m *mockBatchDefRepo) FindByID(ctx context.Context, id primitive.ObjectID) (*batch.BatchDefinition, error) {
	args := m.Called(ctx, id)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*batch.BatchDefinition), args.Error(1)
}

func (m *mockBatchDefRepo) FindAll(ctx context.Context, filter batch.BatchDefinitionFilter) ([]*batch.BatchDefinition, int64, error) {
	args := m.Called(ctx, filter)
	if args.Get(0) == nil {
		return nil, 0, args.Error(2)
	}
	return args.Get(0).([]*batch.BatchDefinition), int64(args.Int(1)), args.Error(2)
}

// TestBatchReportHandler_GetProductionReport tests the production report endpoint
func TestBatchReportHandler_GetProductionReport(t *testing.T) {
	gin.SetMode(gin.TestMode)

	// Setup mocks
	mockRecordRepo := new(mockBatchRecordRepo)
	mockUsageLogRepo := new(mockBatchUsageLogRepo)
	mockDefRepo := new(mockBatchDefRepo)

	// Create service and handler
	reportService := services.NewBatchReportService(mockRecordRepo, mockUsageLogRepo, mockDefRepo)
	handler := NewBatchReportHandler(reportService)

	// Setup test data
	now := time.Now()
	records := []*batch.BatchRecord{
		{
			ID:                primitive.NewObjectID(),
			BatchName:         "Coffee Concentrate",
			QuantityProduced:  500,
			TotalCost:         75.0,
			PreparedBy:        "user1",
			PreparedAt:        now,
		},
		{
			ID:                primitive.NewObjectID(),
			BatchName:         "Coffee Concentrate",
			QuantityProduced:  300,
			TotalCost:         45.0,
			PreparedBy:        "user2",
			PreparedAt:        now,
		},
	}

	// Setup expectations
	mockRecordRepo.On("FindAll", mock.Anything, mock.Anything).Return(records, 2, nil)

	// Create request
	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	fromDate := now.Add(-24 * time.Hour).Format(time.RFC3339)
	toDate := now.Add(24 * time.Hour).Format(time.RFC3339)
	c.Request = httptest.NewRequest("GET", "/api/batch-reports/production?from_date="+fromDate+"&to_date="+toDate, nil)

	// Execute
	handler.GetProductionReport(c)

	// Verify
	assert.Equal(t, http.StatusOK, w.Code)
	mockRecordRepo.AssertExpectations(t)
}

// TestBatchReportHandler_GetProductionReport_MissingDates tests validation
func TestBatchReportHandler_GetProductionReport_MissingDates(t *testing.T) {
	gin.SetMode(gin.TestMode)

	// Setup mocks
	mockRecordRepo := new(mockBatchRecordRepo)
	mockUsageLogRepo := new(mockBatchUsageLogRepo)
	mockDefRepo := new(mockBatchDefRepo)

	// Create service and handler
	reportService := services.NewBatchReportService(mockRecordRepo, mockUsageLogRepo, mockDefRepo)
	handler := NewBatchReportHandler(reportService)

	// Test missing from_date
	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Request = httptest.NewRequest("GET", "/api/batch-reports/production?to_date=2026-02-13T00:00:00Z", nil)
	handler.GetProductionReport(c)
	assert.Equal(t, http.StatusBadRequest, w.Code)

	// Test missing to_date
	w = httptest.NewRecorder()
	c, _ = gin.CreateTestContext(w)
	c.Request = httptest.NewRequest("GET", "/api/batch-reports/production?from_date=2026-02-13T00:00:00Z", nil)
	handler.GetProductionReport(c)
	assert.Equal(t, http.StatusBadRequest, w.Code)
}

// TestBatchReportHandler_GetWastageReport tests the wastage report endpoint
func TestBatchReportHandler_GetWastageReport(t *testing.T) {
	gin.SetMode(gin.TestMode)

	// Setup mocks
	mockRecordRepo := new(mockBatchRecordRepo)
	mockUsageLogRepo := new(mockBatchUsageLogRepo)
	mockDefRepo := new(mockBatchDefRepo)

	// Create service and handler
	reportService := services.NewBatchReportService(mockRecordRepo, mockUsageLogRepo, mockDefRepo)
	handler := NewBatchReportHandler(reportService)

	// Setup test data
	now := time.Now()
	records := []*batch.BatchRecord{
		{
			ID:                primitive.NewObjectID(),
			BatchName:         "Coffee Concentrate",
			QuantityRemaining: 200,
			CostPerUnit:       0.15,
			Status:            batch.BatchStatusExpired,
		},
	}

	// Setup expectations
	mockRecordRepo.On("FindAll", mock.Anything, mock.Anything).Return(records, 1, nil)

	// Create request
	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	fromDate := now.Add(-24 * time.Hour).Format(time.RFC3339)
	toDate := now.Add(24 * time.Hour).Format(time.RFC3339)
	c.Request = httptest.NewRequest("GET", "/api/batch-reports/wastage?from_date="+fromDate+"&to_date="+toDate, nil)

	// Execute
	handler.GetWastageReport(c)

	// Verify
	assert.Equal(t, http.StatusOK, w.Code)
	mockRecordRepo.AssertExpectations(t)
}

// TestBatchReportHandler_GetUsageReport tests the usage report endpoint
func TestBatchReportHandler_GetUsageReport(t *testing.T) {
	gin.SetMode(gin.TestMode)

	// Setup mocks
	mockRecordRepo := new(mockBatchRecordRepo)
	mockUsageLogRepo := new(mockBatchUsageLogRepo)
	mockDefRepo := new(mockBatchDefRepo)

	// Create service and handler
	reportService := services.NewBatchReportService(mockRecordRepo, mockUsageLogRepo, mockDefRepo)
	handler := NewBatchReportHandler(reportService)

	// Setup test data
	now := time.Now()
	logs := []*batch.BatchUsageLog{
		{
			ID:            primitive.NewObjectID(),
			BatchName:     "Coffee Concentrate",
			MenuItemName:  "Black Coffee",
			QuantityUsed:  30,
			TotalCost:     4.5,
			UsedAt:        now,
		},
	}

	// Setup expectations
	mockUsageLogRepo.On("FindAll", mock.Anything, mock.Anything).Return(logs, 1, nil)

	// Create request
	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	fromDate := now.Add(-24 * time.Hour).Format(time.RFC3339)
	toDate := now.Add(24 * time.Hour).Format(time.RFC3339)
	c.Request = httptest.NewRequest("GET", "/api/batch-reports/usage?from_date="+fromDate+"&to_date="+toDate, nil)

	// Execute
	handler.GetUsageReport(c)

	// Verify
	assert.Equal(t, http.StatusOK, w.Code)
	mockUsageLogRepo.AssertExpectations(t)
}
