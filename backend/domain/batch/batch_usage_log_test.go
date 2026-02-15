package batch

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

func TestBatchUsageLog_CalculateTotalCost(t *testing.T) {
	tests := []struct {
		name         string
		quantityUsed float64
		costPerUnit  float64
		expected     float64
	}{
		{
			name:         "Basic calculation",
			quantityUsed: 100.0,
			costPerUnit:  0.5,
			expected:     50.0,
		},
		{
			name:         "Decimal quantities",
			quantityUsed: 25.5,
			costPerUnit:  1.2,
			expected:     30.6,
		},
		{
			name:         "Zero quantity",
			quantityUsed: 0,
			costPerUnit:  1.0,
			expected:     0,
		},
		{
			name:         "Zero cost",
			quantityUsed: 100,
			costPerUnit:  0,
			expected:     0,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			log := &BatchUsageLog{
				QuantityUsed: tt.quantityUsed,
				CostPerUnit:  tt.costPerUnit,
			}
			result := log.CalculateTotalCost()
			assert.InDelta(t, tt.expected, result, 0.001)
		})
	}
}

func TestNewBatchUsageLog(t *testing.T) {
	batchRecordID := primitive.NewObjectID()
	orderID := primitive.NewObjectID()
	menuItemID := primitive.NewObjectID()

	req := CreateBatchUsageLogRequest{
		BatchRecordID: batchRecordID,
		BatchName:     "Cà Phê Concentrate",
		OrderID:       orderID,
		MenuItemID:    menuItemID,
		MenuItemName:  "Cà Phê Đen",
		QuantityUsed:  30.0,
		Unit:          "ml",
		CostPerUnit:   0.15,
	}

	log := NewBatchUsageLog(req)

	assert.Equal(t, batchRecordID, log.BatchRecordID)
	assert.Equal(t, "Cà Phê Concentrate", log.BatchName)
	assert.Equal(t, orderID, log.OrderID)
	assert.Equal(t, menuItemID, log.MenuItemID)
	assert.Equal(t, "Cà Phê Đen", log.MenuItemName)
	assert.Equal(t, 30.0, log.QuantityUsed)
	assert.Equal(t, "ml", log.Unit)
	assert.Equal(t, 0.15, log.CostPerUnit)
	assert.InDelta(t, 4.5, log.TotalCost, 0.001)
	assert.WithinDuration(t, time.Now(), log.UsedAt, time.Second)
}

func TestNewBatchUsageLog_TotalCostCalculation(t *testing.T) {
	tests := []struct {
		name         string
		quantityUsed float64
		costPerUnit  float64
		expectedCost float64
	}{
		{
			name:         "Standard usage",
			quantityUsed: 50.0,
			costPerUnit:  0.2,
			expectedCost: 10.0,
		},
		{
			name:         "Large quantity",
			quantityUsed: 1000.0,
			costPerUnit:  0.05,
			expectedCost: 50.0,
		},
		{
			name:         "Small quantity",
			quantityUsed: 5.5,
			costPerUnit:  2.0,
			expectedCost: 11.0,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			req := CreateBatchUsageLogRequest{
				BatchRecordID: primitive.NewObjectID(),
				BatchName:     "Test Batch",
				OrderID:       primitive.NewObjectID(),
				MenuItemID:    primitive.NewObjectID(),
				MenuItemName:  "Test Item",
				QuantityUsed:  tt.quantityUsed,
				Unit:          "ml",
				CostPerUnit:   tt.costPerUnit,
			}

			log := NewBatchUsageLog(req)
			assert.InDelta(t, tt.expectedCost, log.TotalCost, 0.001)
		})
	}
}

func TestBatchUsageLogFilter(t *testing.T) {
	// Test that filter struct can be created with various combinations
	batchRecordID := primitive.NewObjectID()
	orderID := primitive.NewObjectID()
	menuItemID := primitive.NewObjectID()
	fromDate := time.Now().Add(-24 * time.Hour)
	toDate := time.Now()

	filter := BatchUsageLogFilter{
		BatchRecordID: &batchRecordID,
		OrderID:       &orderID,
		MenuItemID:    &menuItemID,
		FromDate:      &fromDate,
		ToDate:        &toDate,
		Page:          1,
		Limit:         20,
	}

	assert.Equal(t, batchRecordID, *filter.BatchRecordID)
	assert.Equal(t, orderID, *filter.OrderID)
	assert.Equal(t, menuItemID, *filter.MenuItemID)
	assert.Equal(t, fromDate, *filter.FromDate)
	assert.Equal(t, toDate, *filter.ToDate)
	assert.Equal(t, 1, filter.Page)
	assert.Equal(t, 20, filter.Limit)
}

func TestBatchUsageLogFilter_EmptyFilter(t *testing.T) {
	// Test that filter can be created with no filters
	filter := BatchUsageLogFilter{
		Page:  1,
		Limit: 50,
	}

	assert.Nil(t, filter.BatchRecordID)
	assert.Nil(t, filter.OrderID)
	assert.Nil(t, filter.MenuItemID)
	assert.Nil(t, filter.FromDate)
	assert.Nil(t, filter.ToDate)
	assert.Equal(t, 1, filter.Page)
	assert.Equal(t, 50, filter.Limit)
}
