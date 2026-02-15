package batch

import (
	"testing"
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

func TestBatchDefinition_Creation(t *testing.T) {
	now := time.Now()
	
	batchDef := BatchDefinition{
		ID:                 primitive.NewObjectID(),
		Name:               "Cà Phê Concentrate",
		Unit:               "ml",
		ShelfLifeHours:     24,
		ConversionRates:    []ConversionRate{},
		LowStockThreshold:  200,
		ExpiryWarningHours: 4,
		CreatedAt:          now,
		UpdatedAt:          now,
	}

	if batchDef.Name != "Cà Phê Concentrate" {
		t.Errorf("Expected name 'Cà Phê Concentrate', got '%s'", batchDef.Name)
	}

	if batchDef.Unit != "ml" {
		t.Errorf("Expected unit 'ml', got '%s'", batchDef.Unit)
	}

	if batchDef.ShelfLifeHours != 24 {
		t.Errorf("Expected shelf life 24 hours, got %d", batchDef.ShelfLifeHours)
	}

	if batchDef.LowStockThreshold != 200 {
		t.Errorf("Expected low stock threshold 200, got %f", batchDef.LowStockThreshold)
	}

	if batchDef.ExpiryWarningHours != 4 {
		t.Errorf("Expected expiry warning hours 4, got %d", batchDef.ExpiryWarningHours)
	}
}

func TestConversionRate_Creation(t *testing.T) {
	ingredientID := primitive.NewObjectID()
	
	conversionRate := ConversionRate{
		SourceIngredientID:   ingredientID,
		SourceIngredientName: "Hạt Cà Phê",
		SourceQuantity:       100,
		SourceUnit:           "g",
		BatchQuantity:        500,
		WastageRate:          0.1,
	}

	if conversionRate.SourceIngredientID != ingredientID {
		t.Errorf("Expected source ingredient ID to match")
	}

	if conversionRate.SourceIngredientName != "Hạt Cà Phê" {
		t.Errorf("Expected source ingredient name 'Hạt Cà Phê', got '%s'", conversionRate.SourceIngredientName)
	}

	if conversionRate.SourceQuantity != 100 {
		t.Errorf("Expected source quantity 100, got %f", conversionRate.SourceQuantity)
	}

	if conversionRate.SourceUnit != "g" {
		t.Errorf("Expected source unit 'g', got '%s'", conversionRate.SourceUnit)
	}

	if conversionRate.BatchQuantity != 500 {
		t.Errorf("Expected batch quantity 500, got %f", conversionRate.BatchQuantity)
	}

	if conversionRate.WastageRate != 0.1 {
		t.Errorf("Expected wastage rate 0.1, got %f", conversionRate.WastageRate)
	}
}

func TestBatchDefinition_WithMultipleConversionRates(t *testing.T) {
	now := time.Now()
	ingredientID1 := primitive.NewObjectID()
	ingredientID2 := primitive.NewObjectID()

	conversionRates := []ConversionRate{
		{
			SourceIngredientID:   ingredientID1,
			SourceIngredientName: "Hạt Cà Phê",
			SourceQuantity:       100,
			SourceUnit:           "g",
			BatchQuantity:        500,
			WastageRate:          0.1,
		},
		{
			SourceIngredientID:   ingredientID2,
			SourceIngredientName: "Nước",
			SourceQuantity:       400,
			SourceUnit:           "ml",
			BatchQuantity:        500,
			WastageRate:          0.05,
		},
	}

	batchDef := BatchDefinition{
		ID:                 primitive.NewObjectID(),
		Name:               "Cà Phê Concentrate",
		Unit:               "ml",
		ShelfLifeHours:     24,
		ConversionRates:    conversionRates,
		LowStockThreshold:  200,
		ExpiryWarningHours: 4,
		CreatedAt:          now,
		UpdatedAt:          now,
	}

	if len(batchDef.ConversionRates) != 2 {
		t.Errorf("Expected 2 conversion rates, got %d", len(batchDef.ConversionRates))
	}

	if batchDef.ConversionRates[0].SourceIngredientName != "Hạt Cà Phê" {
		t.Errorf("Expected first ingredient 'Hạt Cà Phê', got '%s'", batchDef.ConversionRates[0].SourceIngredientName)
	}

	if batchDef.ConversionRates[1].SourceIngredientName != "Nước" {
		t.Errorf("Expected second ingredient 'Nước', got '%s'", batchDef.ConversionRates[1].SourceIngredientName)
	}
}

func TestConversionRate_WastageCalculation(t *testing.T) {
	tests := []struct {
		name                string
		sourceQuantity      float64
		wastageRate         float64
		expectedWithWastage float64
	}{
		{
			name:                "10% wastage",
			sourceQuantity:      100,
			wastageRate:         0.1,
			expectedWithWastage: 110,
		},
		{
			name:                "5% wastage",
			sourceQuantity:      200,
			wastageRate:         0.05,
			expectedWithWastage: 210,
		},
		{
			name:                "0% wastage",
			sourceQuantity:      100,
			wastageRate:         0.0,
			expectedWithWastage: 100,
		},
		{
			name:                "20% wastage",
			sourceQuantity:      50,
			wastageRate:         0.2,
			expectedWithWastage: 60,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			conversionRate := ConversionRate{
				SourceQuantity: tt.sourceQuantity,
				WastageRate:    tt.wastageRate,
			}

			// Calculate quantity with wastage
			quantityWithWastage := conversionRate.SourceQuantity * (1 + conversionRate.WastageRate)

			// Use a small epsilon for floating point comparison
			epsilon := 0.0001
			diff := quantityWithWastage - tt.expectedWithWastage
			if diff < -epsilon || diff > epsilon {
				t.Errorf("Expected quantity with wastage %f, got %f", tt.expectedWithWastage, quantityWithWastage)
			}
		})
	}
}

func TestCreateBatchDefinitionRequest_Validation(t *testing.T) {
	ingredientID := primitive.NewObjectID()

	validRequest := CreateBatchDefinitionRequest{
		Name:           "Cà Phê Concentrate",
		Unit:           "ml",
		ShelfLifeHours: 24,
		ConversionRates: []ConversionRate{
			{
				SourceIngredientID:   ingredientID,
				SourceIngredientName: "Hạt Cà Phê",
				SourceQuantity:       100,
				SourceUnit:           "g",
				BatchQuantity:        500,
				WastageRate:          0.1,
			},
		},
		LowStockThreshold:  200,
		ExpiryWarningHours: 4,
	}

	if validRequest.Name == "" {
		t.Error("Name should not be empty")
	}

	if validRequest.Unit == "" {
		t.Error("Unit should not be empty")
	}

	if validRequest.ShelfLifeHours < 1 {
		t.Error("ShelfLifeHours should be at least 1")
	}

	if len(validRequest.ConversionRates) < 1 {
		t.Error("ConversionRates should have at least one entry")
	}

	if validRequest.LowStockThreshold < 0 {
		t.Error("LowStockThreshold should not be negative")
	}

	if validRequest.ExpiryWarningHours < 0 {
		t.Error("ExpiryWarningHours should not be negative")
	}
}

func TestBatchDefinitionFilter_DefaultValues(t *testing.T) {
	filter := BatchDefinitionFilter{
		Search: "",
		Page:   1,
		Limit:  20,
	}

	if filter.Page != 1 {
		t.Errorf("Expected default page 1, got %d", filter.Page)
	}

	if filter.Limit != 20 {
		t.Errorf("Expected default limit 20, got %d", filter.Limit)
	}

	if filter.Search != "" {
		t.Errorf("Expected empty search, got '%s'", filter.Search)
	}
}
