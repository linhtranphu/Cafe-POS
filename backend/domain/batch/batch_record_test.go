package batch

import (
	"testing"
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

func TestBatchRecord_Creation(t *testing.T) {
	now := time.Now()
	expiresAt := now.Add(24 * time.Hour)
	batchDefID := primitive.NewObjectID()

	batchRecord := BatchRecord{
		ID:                primitive.NewObjectID(),
		BatchDefinitionID: batchDefID,
		BatchName:         "Cà Phê Concentrate",
		QuantityProduced:  500,
		QuantityRemaining: 500,
		Unit:              "ml",
		CostPerUnit:       0.15,
		TotalCost:         75.0,
		PreparedBy:        "user_id_123",
		PreparedAt:        now,
		ExpiresAt:         expiresAt,
		Status:            BatchStatusAvailable,
		IngredientsUsed:   []IngredientUsage{},
		CreatedAt:         now,
		UpdatedAt:         now,
	}

	if batchRecord.BatchName != "Cà Phê Concentrate" {
		t.Errorf("Expected batch name 'Cà Phê Concentrate', got '%s'", batchRecord.BatchName)
	}

	if batchRecord.QuantityProduced != 500 {
		t.Errorf("Expected quantity produced 500, got %f", batchRecord.QuantityProduced)
	}

	if batchRecord.QuantityRemaining != 500 {
		t.Errorf("Expected quantity remaining 500, got %f", batchRecord.QuantityRemaining)
	}

	if batchRecord.Status != BatchStatusAvailable {
		t.Errorf("Expected status 'available', got '%s'", batchRecord.Status)
	}
}

func TestIngredientUsage_Creation(t *testing.T) {
	ingredientID := primitive.NewObjectID()

	ingredientUsage := IngredientUsage{
		IngredientID:   ingredientID,
		IngredientName: "Hạt Cà Phê",
		Quantity:       110,
		Unit:           "g",
		CostPerUnit:    0.68,
		TotalCost:      75.0,
	}

	if ingredientUsage.IngredientID != ingredientID {
		t.Error("Expected ingredient ID to match")
	}

	if ingredientUsage.IngredientName != "Hạt Cà Phê" {
		t.Errorf("Expected ingredient name 'Hạt Cà Phê', got '%s'", ingredientUsage.IngredientName)
	}

	if ingredientUsage.Quantity != 110 {
		t.Errorf("Expected quantity 110, got %f", ingredientUsage.Quantity)
	}

	if ingredientUsage.CostPerUnit != 0.68 {
		t.Errorf("Expected cost per unit 0.68, got %f", ingredientUsage.CostPerUnit)
	}

	if ingredientUsage.TotalCost != 75.0 {
		t.Errorf("Expected total cost 75.0, got %f", ingredientUsage.TotalCost)
	}
}

func TestBatchRecord_IsExpired(t *testing.T) {
	tests := []struct {
		name      string
		expiresAt time.Time
		expected  bool
	}{
		{
			name:      "Not expired - future date",
			expiresAt: time.Now().Add(24 * time.Hour),
			expected:  false,
		},
		{
			name:      "Expired - past date",
			expiresAt: time.Now().Add(-1 * time.Hour),
			expected:  true,
		},
		{
			name:      "Just expired - 1 second ago",
			expiresAt: time.Now().Add(-1 * time.Second),
			expected:  true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			batchRecord := BatchRecord{
				ExpiresAt: tt.expiresAt,
			}

			result := batchRecord.IsExpired()
			if result != tt.expected {
				t.Errorf("Expected IsExpired() to be %v, got %v", tt.expected, result)
			}
		})
	}
}

func TestBatchRecord_IsDepleted(t *testing.T) {
	tests := []struct {
		name              string
		quantityRemaining float64
		expected          bool
	}{
		{
			name:              "Not depleted - has quantity",
			quantityRemaining: 100,
			expected:          false,
		},
		{
			name:              "Depleted - zero quantity",
			quantityRemaining: 0,
			expected:          true,
		},
		{
			name:              "Depleted - negative quantity",
			quantityRemaining: -1,
			expected:          true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			batchRecord := BatchRecord{
				QuantityRemaining: tt.quantityRemaining,
			}

			result := batchRecord.IsDepleted()
			if result != tt.expected {
				t.Errorf("Expected IsDepleted() to be %v, got %v", tt.expected, result)
			}
		})
	}
}

func TestBatchRecord_IsAvailable(t *testing.T) {
	now := time.Now()

	tests := []struct {
		name              string
		status            string
		expiresAt         time.Time
		quantityRemaining float64
		expected          bool
	}{
		{
			name:              "Available - all conditions met",
			status:            BatchStatusAvailable,
			expiresAt:         now.Add(24 * time.Hour),
			quantityRemaining: 100,
			expected:          true,
		},
		{
			name:              "Not available - expired status",
			status:            BatchStatusExpired,
			expiresAt:         now.Add(24 * time.Hour),
			quantityRemaining: 100,
			expected:          false,
		},
		{
			name:              "Not available - depleted status",
			status:            BatchStatusDepleted,
			expiresAt:         now.Add(24 * time.Hour),
			quantityRemaining: 0,
			expected:          false,
		},
		{
			name:              "Not available - expired time",
			status:            BatchStatusAvailable,
			expiresAt:         now.Add(-1 * time.Hour),
			quantityRemaining: 100,
			expected:          false,
		},
		{
			name:              "Not available - zero quantity",
			status:            BatchStatusAvailable,
			expiresAt:         now.Add(24 * time.Hour),
			quantityRemaining: 0,
			expected:          false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			batchRecord := BatchRecord{
				Status:            tt.status,
				ExpiresAt:         tt.expiresAt,
				QuantityRemaining: tt.quantityRemaining,
			}

			result := batchRecord.IsAvailable()
			if result != tt.expected {
				t.Errorf("Expected IsAvailable() to be %v, got %v", tt.expected, result)
			}
		})
	}
}

func TestBatchRecord_CalculateExpiryTime(t *testing.T) {
	preparedAt := time.Date(2026, 2, 13, 10, 0, 0, 0, time.UTC)

	tests := []struct {
		name           string
		shelfLifeHours int
		expectedExpiry time.Time
	}{
		{
			name:           "24 hours shelf life",
			shelfLifeHours: 24,
			expectedExpiry: time.Date(2026, 2, 14, 10, 0, 0, 0, time.UTC),
		},
		{
			name:           "48 hours shelf life",
			shelfLifeHours: 48,
			expectedExpiry: time.Date(2026, 2, 15, 10, 0, 0, 0, time.UTC),
		},
		{
			name:           "12 hours shelf life",
			shelfLifeHours: 12,
			expectedExpiry: time.Date(2026, 2, 13, 22, 0, 0, 0, time.UTC),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			batchRecord := BatchRecord{
				PreparedAt: preparedAt,
			}

			result := batchRecord.CalculateExpiryTime(tt.shelfLifeHours)
			if !result.Equal(tt.expectedExpiry) {
				t.Errorf("Expected expiry time %v, got %v", tt.expectedExpiry, result)
			}
		})
	}
}

func TestBatchRecord_UpdateStatus(t *testing.T) {
	now := time.Now()

	tests := []struct {
		name              string
		quantityRemaining float64
		expiresAt         time.Time
		expectedStatus    string
	}{
		{
			name:              "Should be depleted",
			quantityRemaining: 0,
			expiresAt:         now.Add(24 * time.Hour),
			expectedStatus:    BatchStatusDepleted,
		},
		{
			name:              "Should be expired",
			quantityRemaining: 100,
			expiresAt:         now.Add(-1 * time.Hour),
			expectedStatus:    BatchStatusExpired,
		},
		{
			name:              "Should be available",
			quantityRemaining: 100,
			expiresAt:         now.Add(24 * time.Hour),
			expectedStatus:    BatchStatusAvailable,
		},
		{
			name:              "Depleted takes precedence over expired",
			quantityRemaining: 0,
			expiresAt:         now.Add(-1 * time.Hour),
			expectedStatus:    BatchStatusDepleted,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			batchRecord := BatchRecord{
				QuantityRemaining: tt.quantityRemaining,
				ExpiresAt:         tt.expiresAt,
			}

			batchRecord.UpdateStatus()
			if batchRecord.Status != tt.expectedStatus {
				t.Errorf("Expected status '%s', got '%s'", tt.expectedStatus, batchRecord.Status)
			}
		})
	}
}

func TestBatchRecord_DeductQuantity(t *testing.T) {
	now := time.Now()

	tests := []struct {
		name                  string
		initialQuantity       float64
		deductAmount          float64
		expectedRemaining     float64
		expectedError         error
		expectedStatus        string
	}{
		{
			name:              "Valid deduction",
			initialQuantity:   500,
			deductAmount:      100,
			expectedRemaining: 400,
			expectedError:     nil,
			expectedStatus:    BatchStatusAvailable,
		},
		{
			name:              "Deduct all quantity",
			initialQuantity:   100,
			deductAmount:      100,
			expectedRemaining: 0,
			expectedError:     nil,
			expectedStatus:    BatchStatusDepleted,
		},
		{
			name:              "Insufficient quantity",
			initialQuantity:   50,
			deductAmount:      100,
			expectedRemaining: 50,
			expectedError:     ErrInsufficientQuantity,
			expectedStatus:    BatchStatusAvailable,
		},
		{
			name:              "Negative deduction",
			initialQuantity:   100,
			deductAmount:      -10,
			expectedRemaining: 100,
			expectedError:     ErrInvalidQuantity,
			expectedStatus:    BatchStatusAvailable,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			batchRecord := BatchRecord{
				QuantityRemaining: tt.initialQuantity,
				ExpiresAt:         now.Add(24 * time.Hour),
				Status:            BatchStatusAvailable,
			}

			err := batchRecord.DeductQuantity(tt.deductAmount)

			if tt.expectedError != nil {
				if err == nil {
					t.Errorf("Expected error %v, got nil", tt.expectedError)
				} else if err.Error() != tt.expectedError.Error() {
					t.Errorf("Expected error '%v', got '%v'", tt.expectedError, err)
				}
			} else if err != nil {
				t.Errorf("Expected no error, got %v", err)
			}

			if batchRecord.QuantityRemaining != tt.expectedRemaining {
				t.Errorf("Expected remaining quantity %f, got %f", tt.expectedRemaining, batchRecord.QuantityRemaining)
			}

			if batchRecord.Status != tt.expectedStatus {
				t.Errorf("Expected status '%s', got '%s'", tt.expectedStatus, batchRecord.Status)
			}
		})
	}
}

func TestBatchRecord_CalculateTotalCost(t *testing.T) {
	ingredientID1 := primitive.NewObjectID()
	ingredientID2 := primitive.NewObjectID()

	tests := []struct {
		name          string
		ingredients   []IngredientUsage
		expectedTotal float64
	}{
		{
			name: "Single ingredient",
			ingredients: []IngredientUsage{
				{
					IngredientID:   ingredientID1,
					IngredientName: "Hạt Cà Phê",
					Quantity:       110,
					Unit:           "g",
					CostPerUnit:    0.68,
					TotalCost:      75.0,
				},
			},
			expectedTotal: 75.0,
		},
		{
			name: "Multiple ingredients",
			ingredients: []IngredientUsage{
				{
					IngredientID:   ingredientID1,
					IngredientName: "Hạt Cà Phê",
					Quantity:       110,
					Unit:           "g",
					CostPerUnit:    0.68,
					TotalCost:      75.0,
				},
				{
					IngredientID:   ingredientID2,
					IngredientName: "Nước",
					Quantity:       420,
					Unit:           "ml",
					CostPerUnit:    0.01,
					TotalCost:      4.2,
				},
			},
			expectedTotal: 79.2,
		},
		{
			name:          "No ingredients",
			ingredients:   []IngredientUsage{},
			expectedTotal: 0,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			batchRecord := BatchRecord{
				IngredientsUsed: tt.ingredients,
			}

			result := batchRecord.CalculateTotalCost()

			// Use epsilon for floating point comparison
			epsilon := 0.0001
			diff := result - tt.expectedTotal
			if diff < -epsilon || diff > epsilon {
				t.Errorf("Expected total cost %f, got %f", tt.expectedTotal, result)
			}
		})
	}
}

func TestBatchRecord_CalculateCostPerUnit(t *testing.T) {
	tests := []struct {
		name             string
		quantityProduced float64
		totalCost        float64
		expectedCostPer  float64
	}{
		{
			name:             "Normal calculation",
			quantityProduced: 500,
			totalCost:        75.0,
			expectedCostPer:  0.15,
		},
		{
			name:             "Different values",
			quantityProduced: 1000,
			totalCost:        100.0,
			expectedCostPer:  0.1,
		},
		{
			name:             "Zero quantity",
			quantityProduced: 0,
			totalCost:        75.0,
			expectedCostPer:  0,
		},
		{
			name:             "Zero cost",
			quantityProduced: 500,
			totalCost:        0,
			expectedCostPer:  0,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			batchRecord := BatchRecord{
				QuantityProduced: tt.quantityProduced,
				TotalCost:        tt.totalCost,
			}

			result := batchRecord.CalculateCostPerUnit()

			// Use epsilon for floating point comparison
			epsilon := 0.0001
			diff := result - tt.expectedCostPer
			if diff < -epsilon || diff > epsilon {
				t.Errorf("Expected cost per unit %f, got %f", tt.expectedCostPer, result)
			}
		})
	}
}

func TestBatchRecord_WithMultipleIngredients(t *testing.T) {
	now := time.Now()
	batchDefID := primitive.NewObjectID()
	ingredientID1 := primitive.NewObjectID()
	ingredientID2 := primitive.NewObjectID()

	ingredientsUsed := []IngredientUsage{
		{
			IngredientID:   ingredientID1,
			IngredientName: "Hạt Cà Phê",
			Quantity:       110,
			Unit:           "g",
			CostPerUnit:    0.68,
			TotalCost:      75.0,
		},
		{
			IngredientID:   ingredientID2,
			IngredientName: "Nước",
			Quantity:       420,
			Unit:           "ml",
			CostPerUnit:    0.01,
			TotalCost:      4.2,
		},
	}

	batchRecord := BatchRecord{
		ID:                primitive.NewObjectID(),
		BatchDefinitionID: batchDefID,
		BatchName:         "Cà Phê Concentrate",
		QuantityProduced:  500,
		QuantityRemaining: 500,
		Unit:              "ml",
		PreparedBy:        "user_id_123",
		PreparedAt:        now,
		ExpiresAt:         now.Add(24 * time.Hour),
		Status:            BatchStatusAvailable,
		IngredientsUsed:   ingredientsUsed,
		CreatedAt:         now,
		UpdatedAt:         now,
	}

	// Calculate costs
	batchRecord.TotalCost = batchRecord.CalculateTotalCost()
	batchRecord.CostPerUnit = batchRecord.CalculateCostPerUnit()

	if len(batchRecord.IngredientsUsed) != 2 {
		t.Errorf("Expected 2 ingredients used, got %d", len(batchRecord.IngredientsUsed))
	}

	expectedTotal := 79.2
	epsilon := 0.0001
	diff := batchRecord.TotalCost - expectedTotal
	if diff < -epsilon || diff > epsilon {
		t.Errorf("Expected total cost %f, got %f", expectedTotal, batchRecord.TotalCost)
	}

	expectedCostPerUnit := 0.1584
	diff = batchRecord.CostPerUnit - expectedCostPerUnit
	if diff < -epsilon || diff > epsilon {
		t.Errorf("Expected cost per unit %f, got %f", expectedCostPerUnit, batchRecord.CostPerUnit)
	}
}

func TestCreateBatchRecordRequest_Validation(t *testing.T) {
	batchDefID := primitive.NewObjectID()

	validRequest := CreateBatchRecordRequest{
		BatchDefinitionID: batchDefID,
		QuantityProduced:  500,
		PreparedBy:        "user_id_123",
	}

	if validRequest.BatchDefinitionID.IsZero() {
		t.Error("BatchDefinitionID should not be zero")
	}

	if validRequest.QuantityProduced <= 0 {
		t.Error("QuantityProduced should be greater than zero")
	}

	if validRequest.PreparedBy == "" {
		t.Error("PreparedBy should not be empty")
	}
}

func TestBatchRecordFilter_DefaultValues(t *testing.T) {
	filter := BatchRecordFilter{
		Page:  1,
		Limit: 20,
	}

	if filter.Page != 1 {
		t.Errorf("Expected default page 1, got %d", filter.Page)
	}

	if filter.Limit != 20 {
		t.Errorf("Expected default limit 20, got %d", filter.Limit)
	}

	if filter.BatchDefinitionID != nil {
		t.Error("Expected BatchDefinitionID to be nil")
	}

	if filter.Status != "" {
		t.Errorf("Expected empty status, got '%s'", filter.Status)
	}
}

func TestBatchError_Error(t *testing.T) {
	err := &BatchError{
		Code:    "TEST_ERROR",
		Message: "This is a test error",
	}

	if err.Error() != "This is a test error" {
		t.Errorf("Expected error message 'This is a test error', got '%s'", err.Error())
	}
}
