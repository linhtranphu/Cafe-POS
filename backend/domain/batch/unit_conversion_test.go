package batch

import (
	"testing"
)

func TestConvertQuantity(t *testing.T) {
	tests := []struct {
		name        string
		quantity    float64
		fromUnit    string
		toUnit      string
		expected    float64
		expectError bool
	}{
		// Same unit - no conversion
		{
			name:        "Same unit ml",
			quantity:    100,
			fromUnit:    "ml",
			toUnit:      "ml",
			expected:    100,
			expectError: false,
		},
		{
			name:        "Same unit l",
			quantity:    1,
			fromUnit:    "l",
			toUnit:      "l",
			expected:    1,
			expectError: false,
		},
		// Volume conversions
		{
			name:        "ml to l",
			quantity:    1000,
			fromUnit:    "ml",
			toUnit:      "l",
			expected:    1,
			expectError: false,
		},
		{
			name:        "l to ml",
			quantity:    1,
			fromUnit:    "l",
			toUnit:      "ml",
			expected:    1000,
			expectError: false,
		},
		{
			name:        "l to ml (0.5l)",
			quantity:    0.5,
			fromUnit:    "l",
			toUnit:      "ml",
			expected:    500,
			expectError: false,
		},
		{
			name:        "ml to l (200ml)",
			quantity:    200,
			fromUnit:    "ml",
			toUnit:      "l",
			expected:    0.2,
			expectError: false,
		},
		// Weight conversions
		{
			name:        "g to kg",
			quantity:    1000,
			fromUnit:    "g",
			toUnit:      "kg",
			expected:    1,
			expectError: false,
		},
		{
			name:        "kg to g",
			quantity:    1,
			fromUnit:    "kg",
			toUnit:      "g",
			expected:    1000,
			expectError: false,
		},
		{
			name:        "kg to g (0.5kg)",
			quantity:    0.5,
			fromUnit:    "kg",
			toUnit:      "g",
			expected:    500,
			expectError: false,
		},
		{
			name:        "g to kg (250g)",
			quantity:    250,
			fromUnit:    "g",
			toUnit:      "kg",
			expected:    0.25,
			expectError: false,
		},
		// Case insensitive
		{
			name:        "ML to l (uppercase)",
			quantity:    1000,
			fromUnit:    "ML",
			toUnit:      "l",
			expected:    1,
			expectError: false,
		},
		{
			name:        "L to ml (uppercase)",
			quantity:    1,
			fromUnit:    "L",
			toUnit:      "ml",
			expected:    1000,
			expectError: false,
		},
		// Error cases
		{
			name:        "Incompatible units (volume to weight)",
			quantity:    100,
			fromUnit:    "ml",
			toUnit:      "g",
			expected:    0,
			expectError: true,
		},
		{
			name:        "Incompatible units (weight to volume)",
			quantity:    100,
			fromUnit:    "g",
			toUnit:      "ml",
			expected:    0,
			expectError: true,
		},
		{
			name:        "Unsupported unit",
			quantity:    100,
			fromUnit:    "oz",
			toUnit:      "ml",
			expected:    0,
			expectError: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result, err := ConvertQuantity(tt.quantity, tt.fromUnit, tt.toUnit)

			if tt.expectError {
				if err == nil {
					t.Errorf("Expected error but got none")
				}
				return
			}

			if err != nil {
				t.Errorf("Unexpected error: %v", err)
				return
			}

			if result != tt.expected {
				t.Errorf("Expected %.2f, got %.2f", tt.expected, result)
			}
		})
	}
}

func TestNormalizeUnit(t *testing.T) {
	tests := []struct {
		input    string
		expected string
	}{
		{"ml", "ml"},
		{"ML", "ml"},
		{"Ml", "ml"},
		{"mL", "ml"},
		{"l", "l"},
		{"L", "l"},
		{"g", "g"},
		{"G", "g"},
		{"kg", "kg"},
		{"KG", "kg"},
		{"Kg", "kg"},
		{"unknown", "unknown"},
	}

	for _, tt := range tests {
		t.Run(tt.input, func(t *testing.T) {
			result := normalizeUnit(tt.input)
			if result != tt.expected {
				t.Errorf("normalizeUnit(%s) = %s, expected %s", tt.input, result, tt.expected)
			}
		})
	}
}
