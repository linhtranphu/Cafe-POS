package http

import (
	"bytes"
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"cafe-pos/backend/application/services"
	"cafe-pos/backend/domain/batch"
	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

// Mock BatchUsageService for testing
type mockBatchUsageService struct {
	useBatchFunc        func(ctx context.Context, req services.UseBatchRequest) (*services.BatchUsageResult, error)
	getUsageHistoryFunc func(ctx context.Context, filter batch.BatchUsageLogFilter) ([]*batch.BatchUsageLog, error)
}

func (m *mockBatchUsageService) UseBatch(ctx context.Context, req services.UseBatchRequest) (*services.BatchUsageResult, error) {
	if m.useBatchFunc != nil {
		return m.useBatchFunc(ctx, req)
	}
	return nil, nil
}

func (m *mockBatchUsageService) GetUsageHistory(ctx context.Context, filter batch.BatchUsageLogFilter) ([]*batch.BatchUsageLog, error) {
	if m.getUsageHistoryFunc != nil {
		return m.getUsageHistoryFunc(ctx, filter)
	}
	return nil, nil
}

func TestUseBatch_Success(t *testing.T) {
	gin.SetMode(gin.TestMode)

	batchDefID := primitive.NewObjectID()
	orderID := primitive.NewObjectID()
	menuItemID := primitive.NewObjectID()
	batchRecordID := primitive.NewObjectID()

	mockService := &mockBatchUsageService{
		useBatchFunc: func(ctx context.Context, req services.UseBatchRequest) (*services.BatchUsageResult, error) {
			assert.Equal(t, batchDefID, req.BatchDefinitionID)
			assert.Equal(t, 30.0, req.QuantityNeeded)
			assert.Equal(t, orderID, req.OrderID)
			assert.Equal(t, menuItemID, req.MenuItemID)
			assert.Equal(t, "Cà Phê Đen", req.MenuItemName)

			return &services.BatchUsageResult{
				Success: true,
				BatchesUsed: []services.BatchUsageDetail{
					{
						BatchRecordID: batchRecordID,
						QuantityUsed:  30.0,
						CostPerUnit:   0.15,
					},
				},
				TotalCost: 4.5,
				Message:   "Batch used successfully",
			}, nil
		},
	}

	handler := NewBatchUsageHandler(mockService)

	// Create request
	reqBody := UseBatchRequest{
		BatchDefinitionID: batchDefID.Hex(),
		QuantityNeeded:    30.0,
		OrderID:           orderID.Hex(),
		MenuItemID:        menuItemID.Hex(),
		MenuItemName:      "Cà Phê Đen",
	}
	body, _ := json.Marshal(reqBody)

	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Request = httptest.NewRequest("POST", "/api/batch-usage", bytes.NewBuffer(body))
	c.Request.Header.Set("Content-Type", "application/json")

	// Execute
	handler.UseBatch(c)

	// Assert
	assert.Equal(t, http.StatusOK, w.Code)

	var response UseBatchResponse
	err := json.Unmarshal(w.Body.Bytes(), &response)
	assert.NoError(t, err)
	assert.True(t, response.Success)
	assert.Equal(t, 1, len(response.BatchesUsed))
	assert.Equal(t, batchRecordID.Hex(), response.BatchesUsed[0].BatchRecordID)
	assert.Equal(t, 30.0, response.BatchesUsed[0].QuantityUsed)
	assert.Equal(t, 0.15, response.BatchesUsed[0].CostPerUnit)
	assert.Equal(t, 4.5, response.TotalCost)
	assert.Equal(t, "Batch used successfully", response.Message)
}

func TestUseBatch_InsufficientBatch(t *testing.T) {
	gin.SetMode(gin.TestMode)

	batchDefID := primitive.NewObjectID()
	orderID := primitive.NewObjectID()
	menuItemID := primitive.NewObjectID()

	mockService := &mockBatchUsageService{
		useBatchFunc: func(ctx context.Context, req services.UseBatchRequest) (*services.BatchUsageResult, error) {
			return &services.BatchUsageResult{
				Success: false,
				Message: "Insufficient batch quantity. Need: 30.00, Available: 20.00",
			}, nil
		},
	}

	handler := NewBatchUsageHandler(mockService)

	// Create request
	reqBody := UseBatchRequest{
		BatchDefinitionID: batchDefID.Hex(),
		QuantityNeeded:    30.0,
		OrderID:           orderID.Hex(),
		MenuItemID:        menuItemID.Hex(),
		MenuItemName:      "Cà Phê Đen",
	}
	body, _ := json.Marshal(reqBody)

	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Request = httptest.NewRequest("POST", "/api/batch-usage", bytes.NewBuffer(body))
	c.Request.Header.Set("Content-Type", "application/json")

	// Execute
	handler.UseBatch(c)

	// Assert
	assert.Equal(t, http.StatusOK, w.Code)

	var response UseBatchResponse
	err := json.Unmarshal(w.Body.Bytes(), &response)
	assert.NoError(t, err)
	assert.False(t, response.Success)
	assert.Contains(t, response.Message, "Insufficient batch quantity")
}

func TestUseBatch_InvalidBatchDefinitionID(t *testing.T) {
	gin.SetMode(gin.TestMode)

	mockService := &mockBatchUsageService{}
	handler := NewBatchUsageHandler(mockService)

	// Create request with invalid ID
	reqBody := UseBatchRequest{
		BatchDefinitionID: "invalid-id",
		QuantityNeeded:    30.0,
		OrderID:           primitive.NewObjectID().Hex(),
		MenuItemID:        primitive.NewObjectID().Hex(),
		MenuItemName:      "Cà Phê Đen",
	}
	body, _ := json.Marshal(reqBody)

	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Request = httptest.NewRequest("POST", "/api/batch-usage", bytes.NewBuffer(body))
	c.Request.Header.Set("Content-Type", "application/json")

	// Execute
	handler.UseBatch(c)

	// Assert
	assert.Equal(t, http.StatusBadRequest, w.Code)
	assert.Contains(t, w.Body.String(), "invalid batch_definition_id")
}

func TestUseBatch_InvalidOrderID(t *testing.T) {
	gin.SetMode(gin.TestMode)

	mockService := &mockBatchUsageService{}
	handler := NewBatchUsageHandler(mockService)

	// Create request with invalid order ID
	reqBody := UseBatchRequest{
		BatchDefinitionID: primitive.NewObjectID().Hex(),
		QuantityNeeded:    30.0,
		OrderID:           "invalid-id",
		MenuItemID:        primitive.NewObjectID().Hex(),
		MenuItemName:      "Cà Phê Đen",
	}
	body, _ := json.Marshal(reqBody)

	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Request = httptest.NewRequest("POST", "/api/batch-usage", bytes.NewBuffer(body))
	c.Request.Header.Set("Content-Type", "application/json")

	// Execute
	handler.UseBatch(c)

	// Assert
	assert.Equal(t, http.StatusBadRequest, w.Code)
	assert.Contains(t, w.Body.String(), "invalid order_id")
}

func TestUseBatch_InvalidMenuItemID(t *testing.T) {
	gin.SetMode(gin.TestMode)

	mockService := &mockBatchUsageService{}
	handler := NewBatchUsageHandler(mockService)

	// Create request with invalid menu item ID
	reqBody := UseBatchRequest{
		BatchDefinitionID: primitive.NewObjectID().Hex(),
		QuantityNeeded:    30.0,
		OrderID:           primitive.NewObjectID().Hex(),
		MenuItemID:        "invalid-id",
		MenuItemName:      "Cà Phê Đen",
	}
	body, _ := json.Marshal(reqBody)

	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Request = httptest.NewRequest("POST", "/api/batch-usage", bytes.NewBuffer(body))
	c.Request.Header.Set("Content-Type", "application/json")

	// Execute
	handler.UseBatch(c)

	// Assert
	assert.Equal(t, http.StatusBadRequest, w.Code)
	assert.Contains(t, w.Body.String(), "invalid menu_item_id")
}

func TestUseBatch_MissingRequiredFields(t *testing.T) {
	gin.SetMode(gin.TestMode)

	mockService := &mockBatchUsageService{}
	handler := NewBatchUsageHandler(mockService)

	// Create request with missing fields
	reqBody := UseBatchRequest{
		BatchDefinitionID: primitive.NewObjectID().Hex(),
		// Missing other required fields
	}
	body, _ := json.Marshal(reqBody)

	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Request = httptest.NewRequest("POST", "/api/batch-usage", bytes.NewBuffer(body))
	c.Request.Header.Set("Content-Type", "application/json")

	// Execute
	handler.UseBatch(c)

	// Assert
	assert.Equal(t, http.StatusBadRequest, w.Code)
}

func TestUseBatch_ZeroQuantity(t *testing.T) {
	gin.SetMode(gin.TestMode)

	mockService := &mockBatchUsageService{}
	handler := NewBatchUsageHandler(mockService)

	// Create request with zero quantity
	reqBody := UseBatchRequest{
		BatchDefinitionID: primitive.NewObjectID().Hex(),
		QuantityNeeded:    0,
		OrderID:           primitive.NewObjectID().Hex(),
		MenuItemID:        primitive.NewObjectID().Hex(),
		MenuItemName:      "Cà Phê Đen",
	}
	body, _ := json.Marshal(reqBody)

	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Request = httptest.NewRequest("POST", "/api/batch-usage", bytes.NewBuffer(body))
	c.Request.Header.Set("Content-Type", "application/json")

	// Execute
	handler.UseBatch(c)

	// Assert
	assert.Equal(t, http.StatusBadRequest, w.Code)
}

func TestGetUsageHistory_Success(t *testing.T) {
	gin.SetMode(gin.TestMode)

	batchRecordID := primitive.NewObjectID()
	orderID := primitive.NewObjectID()
	menuItemID := primitive.NewObjectID()

	mockService := &mockBatchUsageService{
		getUsageHistoryFunc: func(ctx context.Context, filter batch.BatchUsageLogFilter) ([]*batch.BatchUsageLog, error) {
			assert.NotNil(t, filter.BatchRecordID)
			assert.Equal(t, batchRecordID, *filter.BatchRecordID)
			assert.NotNil(t, filter.OrderID)
			assert.Equal(t, orderID, *filter.OrderID)
			assert.NotNil(t, filter.MenuItemID)
			assert.Equal(t, menuItemID, *filter.MenuItemID)

			return []*batch.BatchUsageLog{
				{
					ID:            primitive.NewObjectID(),
					BatchRecordID: batchRecordID,
					BatchName:     "Cà Phê Concentrate",
					OrderID:       orderID,
					MenuItemID:    menuItemID,
					MenuItemName:  "Cà Phê Đen",
					QuantityUsed:  30.0,
					Unit:          "ml",
					CostPerUnit:   0.15,
					TotalCost:     4.5,
					UsedAt:        time.Now(),
				},
			}, nil
		},
	}

	handler := NewBatchUsageHandler(mockService)

	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Request = httptest.NewRequest("GET", "/api/batch-usage/history?batch_record_id="+batchRecordID.Hex()+"&order_id="+orderID.Hex()+"&menu_item_id="+menuItemID.Hex(), nil)

	// Execute
	handler.GetUsageHistory(c)

	// Assert
	assert.Equal(t, http.StatusOK, w.Code)

	var response map[string]interface{}
	err := json.Unmarshal(w.Body.Bytes(), &response)
	assert.NoError(t, err)
	assert.NotNil(t, response["data"])
	assert.Equal(t, float64(1), response["total"])
}

func TestGetUsageHistory_WithDateFilters(t *testing.T) {
	gin.SetMode(gin.TestMode)

	fromDate := time.Now().Add(-24 * time.Hour)
	toDate := time.Now()

	mockService := &mockBatchUsageService{
		getUsageHistoryFunc: func(ctx context.Context, filter batch.BatchUsageLogFilter) ([]*batch.BatchUsageLog, error) {
			assert.NotNil(t, filter.FromDate)
			assert.NotNil(t, filter.ToDate)
			assert.True(t, filter.FromDate.Equal(fromDate) || filter.FromDate.After(fromDate.Add(-time.Second)))
			assert.True(t, filter.ToDate.Equal(toDate) || filter.ToDate.Before(toDate.Add(time.Second)))

			return []*batch.BatchUsageLog{}, nil
		},
	}

	handler := NewBatchUsageHandler(mockService)

	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Request = httptest.NewRequest("GET", "/api/batch-usage/history?from_date="+fromDate.Format(time.RFC3339)+"&to_date="+toDate.Format(time.RFC3339), nil)

	// Execute
	handler.GetUsageHistory(c)

	// Assert
	assert.Equal(t, http.StatusOK, w.Code)
}

func TestGetUsageHistory_InvalidBatchRecordID(t *testing.T) {
	gin.SetMode(gin.TestMode)

	mockService := &mockBatchUsageService{}
	handler := NewBatchUsageHandler(mockService)

	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Request = httptest.NewRequest("GET", "/api/batch-usage/history?batch_record_id=invalid-id", nil)

	// Execute
	handler.GetUsageHistory(c)

	// Assert
	assert.Equal(t, http.StatusBadRequest, w.Code)
	assert.Contains(t, w.Body.String(), "invalid batch_record_id")
}

func TestGetUsageHistory_InvalidOrderID(t *testing.T) {
	gin.SetMode(gin.TestMode)

	mockService := &mockBatchUsageService{}
	handler := NewBatchUsageHandler(mockService)

	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Request = httptest.NewRequest("GET", "/api/batch-usage/history?order_id=invalid-id", nil)

	// Execute
	handler.GetUsageHistory(c)

	// Assert
	assert.Equal(t, http.StatusBadRequest, w.Code)
	assert.Contains(t, w.Body.String(), "invalid order_id")
}

func TestGetUsageHistory_InvalidMenuItemID(t *testing.T) {
	gin.SetMode(gin.TestMode)

	mockService := &mockBatchUsageService{}
	handler := NewBatchUsageHandler(mockService)

	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Request = httptest.NewRequest("GET", "/api/batch-usage/history?menu_item_id=invalid-id", nil)

	// Execute
	handler.GetUsageHistory(c)

	// Assert
	assert.Equal(t, http.StatusBadRequest, w.Code)
	assert.Contains(t, w.Body.String(), "invalid menu_item_id")
}

func TestGetUsageHistory_InvalidFromDate(t *testing.T) {
	gin.SetMode(gin.TestMode)

	mockService := &mockBatchUsageService{}
	handler := NewBatchUsageHandler(mockService)

	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Request = httptest.NewRequest("GET", "/api/batch-usage/history?from_date=invalid-date", nil)

	// Execute
	handler.GetUsageHistory(c)

	// Assert
	assert.Equal(t, http.StatusBadRequest, w.Code)
	assert.Contains(t, w.Body.String(), "invalid from_date format")
}

func TestGetUsageHistory_InvalidToDate(t *testing.T) {
	gin.SetMode(gin.TestMode)

	mockService := &mockBatchUsageService{}
	handler := NewBatchUsageHandler(mockService)

	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Request = httptest.NewRequest("GET", "/api/batch-usage/history?to_date=invalid-date", nil)

	// Execute
	handler.GetUsageHistory(c)

	// Assert
	assert.Equal(t, http.StatusBadRequest, w.Code)
	assert.Contains(t, w.Body.String(), "invalid to_date format")
}

func TestGetUsageHistory_NoFilters(t *testing.T) {
	gin.SetMode(gin.TestMode)

	mockService := &mockBatchUsageService{
		getUsageHistoryFunc: func(ctx context.Context, filter batch.BatchUsageLogFilter) ([]*batch.BatchUsageLog, error) {
			// Should be called with empty filter (all pointers nil)
			assert.Nil(t, filter.BatchRecordID)
			assert.Nil(t, filter.OrderID)
			assert.Nil(t, filter.MenuItemID)
			assert.Nil(t, filter.FromDate)
			assert.Nil(t, filter.ToDate)
			return []*batch.BatchUsageLog{}, nil
		},
	}

	handler := NewBatchUsageHandler(mockService)

	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Request = httptest.NewRequest("GET", "/api/batch-usage/history", nil)

	// Execute
	handler.GetUsageHistory(c)

	// Assert
	assert.Equal(t, http.StatusOK, w.Code)
}
