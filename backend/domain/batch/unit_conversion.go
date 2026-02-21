package batch

import "fmt"

// ConvertQuantity converts a quantity from one unit to another
// Supports volume units (ml, l) and weight units (g, kg)
func ConvertQuantity(quantity float64, fromUnit, toUnit string) (float64, error) {
	if fromUnit == toUnit {
		return quantity, nil
	}

	// Normalize units to lowercase
	fromUnit = normalizeUnit(fromUnit)
	toUnit = normalizeUnit(toUnit)

	if fromUnit == toUnit {
		return quantity, nil
	}

	// Volume conversions (base unit: ml)
	volumeUnits := map[string]float64{
		"ml": 1,
		"l":  1000,
	}

	// Weight conversions (base unit: g)
	weightUnits := map[string]float64{
		"g":  1,
		"kg": 1000,
	}

	// Try volume conversion
	if fromBase, ok := volumeUnits[fromUnit]; ok {
		if toBase, ok := volumeUnits[toUnit]; ok {
			// Convert to base unit (ml), then to target unit
			return quantity * fromBase / toBase, nil
		}
	}

	// Try weight conversion
	if fromBase, ok := weightUnits[fromUnit]; ok {
		if toBase, ok := weightUnits[toUnit]; ok {
			// Convert to base unit (g), then to target unit
			return quantity * fromBase / toBase, nil
		}
	}

	return 0, fmt.Errorf("cannot convert from %s to %s: incompatible or unsupported units", fromUnit, toUnit)
}

// normalizeUnit normalizes unit strings to lowercase for comparison
func normalizeUnit(unit string) string {
	// Simple lowercase conversion
	// In production, might want more sophisticated normalization
	switch unit {
	case "ML", "Ml", "mL":
		return "ml"
	case "L":
		return "l"
	case "G":
		return "g"
	case "KG", "Kg":
		return "kg"
	default:
		return unit
	}
}
