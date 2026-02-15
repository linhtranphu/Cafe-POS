package http

import (
	"context"
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"cafe-pos/backend/domain/batch"
	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

// mockBatchAlertService is a mock implementation of BatchAlertService for testing
type mockBatchAlertService struct {
	alerts *batch.BatchAlerts
	err    error
}

func (m *mockBatchAlertService) GetAlerts(ctx context.Context) (*batch.BatchAlerts, error) {
	if m.err != nil {
		return nil, m.err
	}
	return m.alerts, nil
}

func (m *mockBatchAlertService) CheckLowStock(ctx context.Context) ([]*batch.LowStockAlert, error) {
	return nil, nil
}

func (m *mockBatchAlertService) CheckExpiring(ctx context.Context) ([]*batch.ExpiringAlert, error) {
	return nil, nil
}

func (m *mockBatchAlertService) CheckExpired(ctx context.Context) ([]*batch.ExpiredAlert, error) {
	return nil, nil
}

func (m *mockBatchAlertService) InvalidateCache() {}

// TestBatchAlertHandler_GetAlerts_Success tests successful alert retrieval
func TestBatchAlertHandler_GetAlerts_Success(t *testing.T) {
	// Setup mock service
	mockAlerts := &batch.BatchAlerts{
		LowStock: []*batch.LowStockAlert{
			{
				BatchDefinitionID: primitive.NewObjectID(),
				BatchName:         "Coffee Concentrate",
				CurrentStock:      30.0,
				Threshold:         50.0,
				Unit:              "ml",
			},
		},
		Expiring: []*batch.ExpiringAlert{
			{
				BatchRecordID:     primitive.NewObjectID(),
				BatchName:         "Tea Concentrate",
				QuantityRemaining: 80.0,
				Unit:              "ml",
				ExpiresAt:         time.Now().Add(2 * time.Hour),
				HoursUntilExpiry:  2,
			},
		},
		Expired: []*batch.ExpiredAlert{
			{
				BatchRecordID:  primitive.NewObjectID(),
				BatchName:      "Milk",
				QuantityWasted: 50.0,
				Unit:           "ml",
				CostWasted:     10.0,
				ExpiredAt:      time.Now().Add(-1 * time.Hour),
			},
		},
		LastChecked: time.Now(),
	}

	mockService := &mockBatchAlertService{
		alerts: mockAlerts,
	}

	handler := NewBatchAlertHandler(mockService)

	// Setup Gin
	gin.SetMode(gin.TestMode)
	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)

	// Create request
	req, _ := http.NewRequest("GET", "/api/batch-alerts", nil)
	c.Request = req

	// Call handler
	handler.GetAlerts(c)

	// Verify response
	assert.Equal(t, http.StatusOK, w.Code)

	var response batch.BatchAlerts
	err := json.Unmarshal(w.Body.Bytes(), &response)
	assert.NoError(t, err)

	// Verify alerts
	assert.Len(t, response.LowStock, 1)
	assert.Equal(t, "Coffee Concentrate", response.LowStock[0].BatchName)
	assert.Equal(t, 30.0, response.LowStock[0].CurrentStock)

	assert.Len(t, response.Expiring, 1)
	assert.Equal(t, "Tea Concentrate", response.Expiring[0].BatchName)
	assert.Equal(t, 2, response.Expiring[0].HoursUntilExpiry)

	assert.Len(t, response.Expired, 1)
	assert.Equal(t, "Milk", response.Expired[0].BatchName)
	assert.Equal(t, 10.0, response.Expired[0].CostWasted)

	assert.False(t, response.LastChecked.IsZero())
}

// TestBatchAlertHandler_GetAlerts_NoAlerts tests when there are no alerts
func TestBatchAlertHandler_GetAlerts_NoAlerts(t *testing.T) {
	// Setup mock service with empty alerts
	mockAlerts := &batch.BatchAlerts{
		LowStock:    []*batch.LowStockAlert{},
		Expiring:    []*batch.ExpiringAlert{},
		Expired:     []*batch.ExpiredAlert{},
		LastChecked: time.Now(),
	}

	mockService := &mockBatchAlertService{
		alerts: mockAlerts,
	}

	handler := NewBatchAlertHandler(mockService)

	// Setup Gin
	gin.SetMode(gin.TestMode)
	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)

	// Create request
	req, _ := http.NewRequest("GET", "/api/batch-alerts", nil)
	c.Request = req

	// Call handler
	handler.GetAlerts(c)

	// Verify response
	assert.Equal(t, http.StatusOK, w.Code)

	var response batch.BatchAlerts
	err := json.Unmarshal(w.Body.Bytes(), &response)
	assert.NoError(t, err)

	// Verify no alerts
	assert.Empty(t, response.LowStock)
	assert.Empty(t, response.Expiring)
	assert.Empty(t, response.Expired)
	assert.False(t, response.LastChecked.IsZero())
}

// TestBatchAlertHandler_GetAlerts_ServiceError tests error handling
func TestBatchAlertHandler_GetAlerts_ServiceError(t *testing.T) {
	// Setup mock service with error
	mockService := &mockBatchAlertService{
		err: errors.New("database connection failed"),
	}

	handler := NewBatchAlertHandler(mockService)

	// Setup Gin
	gin.SetMode(gin.TestMode)
	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)

	// Create request
	req, _ := http.NewRequest("GET", "/api/batch-alerts", nil)
	c.Request = req

	// Call handler
	handler.GetAlerts(c)

	// Verify error response
	assert.Equal(t, http.StatusInternalServerError, w.Code)

	var response map[string]string
	err := json.Unmarshal(w.Body.Bytes(), &response)
	assert.NoError(t, err)
	assert.Contains(t, response["error"], "database connection failed")
}

// TestBatchAlertHandler_GetAlerts_MultipleAlerts tests multiple alerts of each type
func TestBatchAlertHandler_GetAlerts_MultipleAlerts(t *testing.T) {
	// Setup mock service with multiple alerts
	mockAlerts := &batch.BatchAlerts{
		LowStock: []*batch.LowStockAlert{
			{
				BatchDefinitionID: primitive.NewObjectID(),
				BatchName:         "Coffee Concentrate",
				CurrentStock:      30.0,
				Threshold:         50.0,
				Unit:              "ml",
			},
			{
				BatchDefinitionID: primitive.NewObjectID(),
				BatchName:         "Tea Concentrate",
				CurrentStock:      20.0,
				Threshold:         100.0,
				Unit:              "ml",
			},
		},
		Expiring: []*batch.ExpiringAlert{
			{
				BatchRecordID:     primitive.NewObjectID(),
				BatchName:         "Milk",
				QuantityRemaining: 50.0,
				Unit:              "ml",
				ExpiresAt:         time.Now().Add(1 * time.Hour),
				HoursUntilExpiry:  1,
			},
			{
				BatchRecordID:     primitive.NewObjectID(),
				BatchName:         "Cream",
				QuantityRemaining: 30.0,
				Unit:              "ml",
				ExpiresAt:         time.Now().Add(3 * time.Hour),
				HoursUntilExpiry:  3,
			},
		},
		Expired: []*batch.ExpiredAlert{
			{
				BatchRecordID:  primitive.NewObjectID(),
				BatchName:      "Old Coffee",
				QuantityWasted: 100.0,
				Unit:           "ml",
				CostWasted:     15.0,
				ExpiredAt:      time.Now().Add(-2 * time.Hour),
			},
		},
		LastChecked: time.Now(),
	}

	mockService := &mockBatchAlertService{
		alerts: mockAlerts,
	}

	handler := NewBatchAlertHandler(mockService)

	// Setup Gin
	gin.SetMode(gin.TestMode)
	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)

	// Create request
	req, _ := http.NewRequest("GET", "/api/batch-alerts", nil)
	c.Request = req

	// Call handler
	handler.GetAlerts(c)

	// Verify response
	assert.Equal(t, http.StatusOK, w.Code)

	var response batch.BatchAlerts
	err := json.Unmarshal(w.Body.Bytes(), &response)
	assert.NoError(t, err)

	// Verify multiple alerts
	assert.Len(t, response.LowStock, 2)
	assert.Len(t, response.Expiring, 2)
	assert.Len(t, response.Expired, 1)
}
