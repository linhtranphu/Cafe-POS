package batch

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

func TestNewBatchAlerts(t *testing.T) {
	alerts := NewBatchAlerts()

	assert.NotNil(t, alerts)
	assert.NotNil(t, alerts.LowStock)
	assert.NotNil(t, alerts.Expiring)
	assert.NotNil(t, alerts.Expired)
	assert.Equal(t, 0, len(alerts.LowStock))
	assert.Equal(t, 0, len(alerts.Expiring))
	assert.Equal(t, 0, len(alerts.Expired))
	assert.False(t, alerts.LastChecked.IsZero())
}

func TestBatchAlerts_AddLowStockAlert(t *testing.T) {
	alerts := NewBatchAlerts()
	lowStockAlert := NewLowStockAlert(
		primitive.NewObjectID(),
		"Coffee Concentrate",
		150.0,
		200.0,
		"ml",
	)

	alerts.AddLowStockAlert(lowStockAlert)

	assert.Equal(t, 1, len(alerts.LowStock))
	assert.Equal(t, "Coffee Concentrate", alerts.LowStock[0].BatchName)
	assert.Equal(t, 150.0, alerts.LowStock[0].CurrentStock)
}

func TestBatchAlerts_AddExpiringAlert(t *testing.T) {
	alerts := NewBatchAlerts()
	expiresAt := time.Now().Add(2 * time.Hour)
	expiringAlert := NewExpiringAlert(
		primitive.NewObjectID(),
		"Coffee Concentrate",
		200.0,
		"ml",
		expiresAt,
	)

	alerts.AddExpiringAlert(expiringAlert)

	assert.Equal(t, 1, len(alerts.Expiring))
	assert.Equal(t, "Coffee Concentrate", alerts.Expiring[0].BatchName)
	assert.Equal(t, 200.0, alerts.Expiring[0].QuantityRemaining)
}

func TestBatchAlerts_AddExpiredAlert(t *testing.T) {
	alerts := NewBatchAlerts()
	expiredAt := time.Now().Add(-1 * time.Hour)
	expiredAlert := NewExpiredAlert(
		primitive.NewObjectID(),
		"Coffee Concentrate",
		100.0,
		"ml",
		15.0,
		expiredAt,
	)

	alerts.AddExpiredAlert(expiredAlert)

	assert.Equal(t, 1, len(alerts.Expired))
	assert.Equal(t, "Coffee Concentrate", alerts.Expired[0].BatchName)
	assert.Equal(t, 100.0, alerts.Expired[0].QuantityWasted)
	assert.Equal(t, 15.0, alerts.Expired[0].CostWasted)
}

func TestBatchAlerts_HasAlerts(t *testing.T) {
	tests := []struct {
		name     string
		setup    func(*BatchAlerts)
		expected bool
	}{
		{
			name: "no alerts",
			setup: func(ba *BatchAlerts) {
				// No alerts added
			},
			expected: false,
		},
		{
			name: "has low stock alert",
			setup: func(ba *BatchAlerts) {
				ba.AddLowStockAlert(NewLowStockAlert(
					primitive.NewObjectID(),
					"Coffee",
					50.0,
					100.0,
					"ml",
				))
			},
			expected: true,
		},
		{
			name: "has expiring alert",
			setup: func(ba *BatchAlerts) {
				ba.AddExpiringAlert(NewExpiringAlert(
					primitive.NewObjectID(),
					"Coffee",
					100.0,
					"ml",
					time.Now().Add(2*time.Hour),
				))
			},
			expected: true,
		},
		{
			name: "has expired alert",
			setup: func(ba *BatchAlerts) {
				ba.AddExpiredAlert(NewExpiredAlert(
					primitive.NewObjectID(),
					"Coffee",
					50.0,
					"ml",
					10.0,
					time.Now().Add(-1*time.Hour),
				))
			},
			expected: true,
		},
		{
			name: "has multiple alerts",
			setup: func(ba *BatchAlerts) {
				ba.AddLowStockAlert(NewLowStockAlert(
					primitive.NewObjectID(),
					"Coffee",
					50.0,
					100.0,
					"ml",
				))
				ba.AddExpiringAlert(NewExpiringAlert(
					primitive.NewObjectID(),
					"Tea",
					100.0,
					"ml",
					time.Now().Add(2*time.Hour),
				))
			},
			expected: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			alerts := NewBatchAlerts()
			tt.setup(alerts)
			assert.Equal(t, tt.expected, alerts.HasAlerts())
		})
	}
}

func TestBatchAlerts_TotalAlertCount(t *testing.T) {
	alerts := NewBatchAlerts()

	// Initially zero
	assert.Equal(t, 0, alerts.TotalAlertCount())

	// Add one of each type
	alerts.AddLowStockAlert(NewLowStockAlert(
		primitive.NewObjectID(),
		"Coffee",
		50.0,
		100.0,
		"ml",
	))
	assert.Equal(t, 1, alerts.TotalAlertCount())

	alerts.AddExpiringAlert(NewExpiringAlert(
		primitive.NewObjectID(),
		"Tea",
		100.0,
		"ml",
		time.Now().Add(2*time.Hour),
	))
	assert.Equal(t, 2, alerts.TotalAlertCount())

	alerts.AddExpiredAlert(NewExpiredAlert(
		primitive.NewObjectID(),
		"Milk",
		50.0,
		"ml",
		10.0,
		time.Now().Add(-1*time.Hour),
	))
	assert.Equal(t, 3, alerts.TotalAlertCount())
}

func TestNewLowStockAlert(t *testing.T) {
	batchDefID := primitive.NewObjectID()
	alert := NewLowStockAlert(batchDefID, "Coffee Concentrate", 150.0, 200.0, "ml")

	assert.NotNil(t, alert)
	assert.Equal(t, batchDefID, alert.BatchDefinitionID)
	assert.Equal(t, "Coffee Concentrate", alert.BatchName)
	assert.Equal(t, 150.0, alert.CurrentStock)
	assert.Equal(t, 200.0, alert.Threshold)
	assert.Equal(t, "ml", alert.Unit)
}

func TestLowStockAlert_IsLowStock(t *testing.T) {
	tests := []struct {
		name         string
		currentStock float64
		threshold    float64
		expected     bool
	}{
		{
			name:         "stock below threshold",
			currentStock: 150.0,
			threshold:    200.0,
			expected:     true,
		},
		{
			name:         "stock equal to threshold",
			currentStock: 200.0,
			threshold:    200.0,
			expected:     true,
		},
		{
			name:         "stock above threshold",
			currentStock: 250.0,
			threshold:    200.0,
			expected:     false,
		},
		{
			name:         "zero stock",
			currentStock: 0.0,
			threshold:    100.0,
			expected:     true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			alert := NewLowStockAlert(
				primitive.NewObjectID(),
				"Test Batch",
				tt.currentStock,
				tt.threshold,
				"ml",
			)
			assert.Equal(t, tt.expected, alert.IsLowStock())
		})
	}
}

func TestNewExpiringAlert(t *testing.T) {
	batchRecordID := primitive.NewObjectID()
	expiresAt := time.Now().Add(3 * time.Hour)
	alert := NewExpiringAlert(batchRecordID, "Coffee Concentrate", 200.0, "ml", expiresAt)

	assert.NotNil(t, alert)
	assert.Equal(t, batchRecordID, alert.BatchRecordID)
	assert.Equal(t, "Coffee Concentrate", alert.BatchName)
	assert.Equal(t, 200.0, alert.QuantityRemaining)
	assert.Equal(t, "ml", alert.Unit)
	assert.Equal(t, expiresAt, alert.ExpiresAt)
	// HoursUntilExpiry should be approximately 3 (may be 2 due to rounding)
	assert.True(t, alert.HoursUntilExpiry >= 2 && alert.HoursUntilExpiry <= 3)
}

func TestExpiringAlert_IsExpiringSoon(t *testing.T) {
	tests := []struct {
		name         string
		hoursUntil   time.Duration
		warningHours int
		expected     bool
	}{
		{
			name:         "expiring within warning period",
			hoursUntil:   2 * time.Hour,
			warningHours: 4,
			expected:     true,
		},
		{
			name:         "expiring exactly at warning hours",
			hoursUntil:   4 * time.Hour,
			warningHours: 4,
			expected:     true,
		},
		{
			name:         "expiring beyond warning period",
			hoursUntil:   6 * time.Hour,
			warningHours: 4,
			expected:     false,
		},
		{
			name:         "already expired",
			hoursUntil:   -1 * time.Hour,
			warningHours: 4,
			expected:     false,
		},
		{
			name:         "expiring in 1 hour with 4 hour warning",
			hoursUntil:   90 * time.Minute, // Use 90 minutes to ensure it rounds to 1 hour
			warningHours: 4,
			expected:     true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			expiresAt := time.Now().Add(tt.hoursUntil)
			alert := NewExpiringAlert(
				primitive.NewObjectID(),
				"Test Batch",
				100.0,
				"ml",
				expiresAt,
			)
			assert.Equal(t, tt.expected, alert.IsExpiringSoon(tt.warningHours))
		})
	}
}

func TestNewExpiredAlert(t *testing.T) {
	batchRecordID := primitive.NewObjectID()
	expiredAt := time.Now().Add(-2 * time.Hour)
	alert := NewExpiredAlert(batchRecordID, "Coffee Concentrate", 100.0, "ml", 15.0, expiredAt)

	assert.NotNil(t, alert)
	assert.Equal(t, batchRecordID, alert.BatchRecordID)
	assert.Equal(t, "Coffee Concentrate", alert.BatchName)
	assert.Equal(t, 100.0, alert.QuantityWasted)
	assert.Equal(t, "ml", alert.Unit)
	assert.Equal(t, 15.0, alert.CostWasted)
	assert.Equal(t, expiredAt, alert.ExpiredAt)
}

func TestExpiredAlert_CalculateCostWasted(t *testing.T) {
	tests := []struct {
		name           string
		quantityWasted float64
		costPerUnit    float64
		expected       float64
	}{
		{
			name:           "normal calculation",
			quantityWasted: 100.0,
			costPerUnit:    0.15,
			expected:       15.0,
		},
		{
			name:           "zero quantity",
			quantityWasted: 0.0,
			costPerUnit:    0.15,
			expected:       0.0,
		},
		{
			name:           "zero cost",
			quantityWasted: 100.0,
			costPerUnit:    0.0,
			expected:       0.0,
		},
		{
			name:           "large quantity",
			quantityWasted: 1000.0,
			costPerUnit:    2.5,
			expected:       2500.0,
		},
		{
			name:           "fractional values",
			quantityWasted: 50.5,
			costPerUnit:    1.25,
			expected:       63.125,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			alert := NewExpiredAlert(
				primitive.NewObjectID(),
				"Test Batch",
				tt.quantityWasted,
				"ml",
				0.0, // Initial cost wasted
				time.Now(),
			)
			result := alert.CalculateCostWasted(tt.costPerUnit)
			assert.Equal(t, tt.expected, result)
		})
	}
}

func TestBatchAlerts_MultipleAlertsScenario(t *testing.T) {
	// Simulate a realistic scenario with multiple alerts
	alerts := NewBatchAlerts()

	// Add 2 low stock alerts
	alerts.AddLowStockAlert(NewLowStockAlert(
		primitive.NewObjectID(),
		"Coffee Concentrate",
		50.0,
		200.0,
		"ml",
	))
	alerts.AddLowStockAlert(NewLowStockAlert(
		primitive.NewObjectID(),
		"Tea Concentrate",
		30.0,
		100.0,
		"ml",
	))

	// Add 3 expiring alerts
	for i := 0; i < 3; i++ {
		alerts.AddExpiringAlert(NewExpiringAlert(
			primitive.NewObjectID(),
			"Milk Batch",
			100.0,
			"ml",
			time.Now().Add(2*time.Hour),
		))
	}

	// Add 1 expired alert
	alerts.AddExpiredAlert(NewExpiredAlert(
		primitive.NewObjectID(),
		"Old Coffee",
		75.0,
		"ml",
		11.25,
		time.Now().Add(-1*time.Hour),
	))

	assert.Equal(t, 2, len(alerts.LowStock))
	assert.Equal(t, 3, len(alerts.Expiring))
	assert.Equal(t, 1, len(alerts.Expired))
	assert.Equal(t, 6, alerts.TotalAlertCount())
	assert.True(t, alerts.HasAlerts())
}
