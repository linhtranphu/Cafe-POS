package ingredient

import (
	"testing"
)

// TestGetConversionRate_SameUnit tests conversion when stock and recipe units are the same
func TestGetConversionRate_SameUnit(t *testing.T) {
	tests := []struct {
		name     string
		unit     UnitType
		expected float64
	}{
		{"kg to kg", UnitKilogram, 1.0},
		{"g to g", UnitGram, 1.0},
		{"L to L", UnitLiter, 1.0},
		{"ml to ml", UnitMilliliter, 1.0},
		{"piece to piece", UnitPiece, 1.0},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := GetConversionRate(tt.unit, tt.unit)
			if result != tt.expected {
				t.Errorf("GetConversionRate(%v, %v) = %v, want %v", 
					tt.unit, tt.unit, result, tt.expected)
			}
		})
	}
}

// TestGetConversionRate_MassConversions tests mass unit conversions
func TestGetConversionRate_MassConversions(t *testing.T) {
	tests := []struct {
		name       string
		stockUnit  UnitType
		recipeUnit UnitType
		expected   float64
	}{
		{"kg to g", UnitKilogram, UnitGram, 0.001},
		{"g to kg", UnitGram, UnitKilogram, 1000.0},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := GetConversionRate(tt.stockUnit, tt.recipeUnit)
			if result != tt.expected {
				t.Errorf("GetConversionRate(%v, %v) = %v, want %v", 
					tt.stockUnit, tt.recipeUnit, result, tt.expected)
			}
		})
	}
}

// TestGetConversionRate_VolumeConversions tests volume unit conversions
func TestGetConversionRate_VolumeConversions(t *testing.T) {
	tests := []struct {
		name       string
		stockUnit  UnitType
		recipeUnit UnitType
		expected   float64
	}{
		{"L to ml", UnitLiter, UnitMilliliter, 0.001},
		{"ml to L", UnitMilliliter, UnitLiter, 1000.0},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := GetConversionRate(tt.stockUnit, tt.recipeUnit)
			if result != tt.expected {
				t.Errorf("GetConversionRate(%v, %v) = %v, want %v", 
					tt.stockUnit, tt.recipeUnit, result, tt.expected)
			}
		})
	}
}

// TestGetConversionRate_InvalidConversions tests invalid conversions (different categories)
func TestGetConversionRate_InvalidConversions(t *testing.T) {
	tests := []struct {
		name       string
		stockUnit  UnitType
		recipeUnit UnitType
		expected   float64
	}{
		{"kg to L (mass to volume)", UnitKilogram, UnitLiter, 1.0},
		{"L to kg (volume to mass)", UnitLiter, UnitKilogram, 1.0},
		{"g to ml (mass to volume)", UnitGram, UnitMilliliter, 1.0},
		{"piece to kg (count to mass)", UnitPiece, UnitKilogram, 1.0},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := GetConversionRate(tt.stockUnit, tt.recipeUnit)
			if result != tt.expected {
				t.Errorf("GetConversionRate(%v, %v) = %v, want %v (fallback)", 
					tt.stockUnit, tt.recipeUnit, result, tt.expected)
			}
		})
	}
}

// TestGetConversionRate_RealWorldScenarios tests real-world use cases
func TestGetConversionRate_RealWorldScenarios(t *testing.T) {
	tests := []struct {
		name         string
		stockUnit    UnitType
		recipeUnit   UnitType
		quantity     float64
		costPerUnit  float64
		expectedCost float64
		description  string
	}{
		{
			name:         "Milk: L to ml",
			stockUnit:    UnitLiter,
			recipeUnit:   UnitMilliliter,
			quantity:     150,           // 150ml
			costPerUnit:  50000,         // 50,000 VND/L
			expectedCost: 7500,          // 150 * 50,000 * 0.001 = 7,500
			description:  "Coffee shop uses 150ml milk, stock in liters",
		},
		{
			name:         "Coffee: kg to g",
			stockUnit:    UnitKilogram,
			recipeUnit:   UnitGram,
			quantity:     20,            // 20g
			costPerUnit:  200000,        // 200,000 VND/kg
			expectedCost: 4000,          // 20 * 200,000 * 0.001 = 4,000
			description:  "Coffee shop uses 20g coffee beans, stock in kg",
		},
		{
			name:         "Sugar: kg to g",
			stockUnit:    UnitKilogram,
			recipeUnit:   UnitGram,
			quantity:     10,            // 10g
			costPerUnit:  25000,         // 25,000 VND/kg
			expectedCost: 250,           // 10 * 25,000 * 0.001 = 250
			description:  "Coffee shop uses 10g sugar, stock in kg",
		},
		{
			name:         "Water: L to L (no conversion)",
			stockUnit:    UnitLiter,
			recipeUnit:   UnitLiter,
			quantity:     0.5,           // 0.5L
			costPerUnit:  10000,         // 10,000 VND/L
			expectedCost: 5000,          // 0.5 * 10,000 * 1.0 = 5,000
			description:  "Same unit, no conversion needed",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			conversionRate := GetConversionRate(tt.stockUnit, tt.recipeUnit)
			actualCost := tt.quantity * tt.costPerUnit * conversionRate
			
			if actualCost != tt.expectedCost {
				t.Errorf("%s: Expected cost %v, got %v (conversion rate: %v)", 
					tt.description, tt.expectedCost, actualCost, conversionRate)
			}
		})
	}
}

// TestValidateUnitConversion_ValidConversions tests valid unit conversions
func TestValidateUnitConversion_ValidConversions(t *testing.T) {
	tests := []struct {
		name       string
		stockUnit  UnitType
		recipeUnit UnitType
		expected   bool
	}{
		{"kg to g (mass)", UnitKilogram, UnitGram, true},
		{"g to kg (mass)", UnitGram, UnitKilogram, true},
		{"L to ml (volume)", UnitLiter, UnitMilliliter, true},
		{"ml to L (volume)", UnitMilliliter, UnitLiter, true},
		{"same unit kg", UnitKilogram, UnitKilogram, true},
		{"same unit L", UnitLiter, UnitLiter, true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := ValidateUnitConversion(tt.stockUnit, tt.recipeUnit)
			if result != tt.expected {
				t.Errorf("ValidateUnitConversion(%v, %v) = %v, want %v", 
					tt.stockUnit, tt.recipeUnit, result, tt.expected)
			}
		})
	}
}

// TestValidateUnitConversion_InvalidConversions tests invalid unit conversions
func TestValidateUnitConversion_InvalidConversions(t *testing.T) {
	tests := []struct {
		name       string
		stockUnit  UnitType
		recipeUnit UnitType
		expected   bool
	}{
		{"kg to L (mass to volume)", UnitKilogram, UnitLiter, false},
		{"L to kg (volume to mass)", UnitLiter, UnitKilogram, false},
		{"g to ml (mass to volume)", UnitGram, UnitMilliliter, false},
		{"ml to g (volume to mass)", UnitMilliliter, UnitGram, false},
		{"piece to kg (count to mass)", UnitPiece, UnitKilogram, false},
		{"kg to piece (mass to count)", UnitKilogram, UnitPiece, false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := ValidateUnitConversion(tt.stockUnit, tt.recipeUnit)
			if result != tt.expected {
				t.Errorf("ValidateUnitConversion(%v, %v) = %v, want %v", 
					tt.stockUnit, tt.recipeUnit, result, tt.expected)
			}
		})
	}
}

// TestGetCompatibleUnits tests getting compatible units for a given unit
func TestGetCompatibleUnits(t *testing.T) {
	tests := []struct {
		name     string
		unit     UnitType
		expected []UnitType
	}{
		{
			name:     "kg returns mass units",
			unit:     UnitKilogram,
			expected: []UnitType{UnitKilogram, UnitGram},
		},
		{
			name:     "g returns mass units",
			unit:     UnitGram,
			expected: []UnitType{UnitKilogram, UnitGram},
		},
		{
			name:     "L returns volume units",
			unit:     UnitLiter,
			expected: []UnitType{UnitLiter, UnitMilliliter},
		},
		{
			name:     "ml returns volume units",
			unit:     UnitMilliliter,
			expected: []UnitType{UnitLiter, UnitMilliliter},
		},
		{
			name:     "piece returns count units",
			unit:     UnitPiece,
			expected: []UnitType{UnitPiece, UnitBox, UnitPack},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := GetCompatibleUnits(tt.unit)
			
			// Check length
			if len(result) != len(tt.expected) {
				t.Errorf("GetCompatibleUnits(%v) returned %d units, want %d", 
					tt.unit, len(result), len(tt.expected))
				return
			}
			
			// Check each unit
			for i, expectedUnit := range tt.expected {
				if result[i] != expectedUnit {
					t.Errorf("GetCompatibleUnits(%v)[%d] = %v, want %v", 
						tt.unit, i, result[i], expectedUnit)
				}
			}
		})
	}
}

// TestGetConversionRate_Precision tests floating point precision
func TestGetConversionRate_Precision(t *testing.T) {
	// Test that conversion rates are precise enough for cost calculations
	tests := []struct {
		name         string
		stockUnit    UnitType
		recipeUnit   UnitType
		quantity     float64
		costPerUnit  float64
		expectedCost float64
		tolerance    float64
	}{
		{
			name:         "Small quantity precision",
			stockUnit:    UnitKilogram,
			recipeUnit:   UnitGram,
			quantity:     0.5,           // 0.5g
			costPerUnit:  1000000,       // 1,000,000 VND/kg
			expectedCost: 500,           // 0.5 * 1,000,000 * 0.001 = 500
			tolerance:    0.01,
		},
		{
			name:         "Large quantity precision",
			stockUnit:    UnitLiter,
			recipeUnit:   UnitMilliliter,
			quantity:     5000,          // 5000ml = 5L
			costPerUnit:  50000,         // 50,000 VND/L
			expectedCost: 250000,        // 5000 * 50,000 * 0.001 = 250,000
			tolerance:    0.01,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			conversionRate := GetConversionRate(tt.stockUnit, tt.recipeUnit)
			actualCost := tt.quantity * tt.costPerUnit * conversionRate
			
			diff := actualCost - tt.expectedCost
			if diff < 0 {
				diff = -diff
			}
			
			if diff > tt.tolerance {
				t.Errorf("Precision error: expected %v, got %v (diff: %v, tolerance: %v)", 
					tt.expectedCost, actualCost, diff, tt.tolerance)
			}
		})
	}
}

// BenchmarkGetConversionRate benchmarks the conversion rate calculation
func BenchmarkGetConversionRate(b *testing.B) {
	for i := 0; i < b.N; i++ {
		GetConversionRate(UnitLiter, UnitMilliliter)
	}
}

// BenchmarkValidateUnitConversion benchmarks the validation function
func BenchmarkValidateUnitConversion(b *testing.B) {
	for i := 0; i < b.N; i++ {
		ValidateUnitConversion(UnitLiter, UnitMilliliter)
	}
}

// BenchmarkGetCompatibleUnits benchmarks getting compatible units
func BenchmarkGetCompatibleUnits(b *testing.B) {
	for i := 0; i < b.N; i++ {
		GetCompatibleUnits(UnitLiter)
	}
}
