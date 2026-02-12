package ingredient

// GetConversionRate calculates the conversion rate from stock unit to recipe unit
// Returns the multiplier to convert recipe quantity to stock quantity
// Example: stockUnit="L", recipeUnit="ml" => returns 0.001 (1ml = 0.001L)
func GetConversionRate(stockUnit, recipeUnit UnitType) float64 {
	// If same unit, no conversion needed
	if stockUnit == recipeUnit {
		return 1.0
	}
	
	// Mass conversions (base: kg)
	massToKg := map[UnitType]float64{
		UnitKilogram: 1.0,
		UnitGram:     0.001, // 1g = 0.001kg
	}
	
	// Volume conversions (base: L)
	volumeToL := map[UnitType]float64{
		UnitLiter:      1.0,
		UnitMilliliter: 0.001, // 1ml = 0.001L
	}
	
	// Try mass conversion
	if stockFactor, stockOk := massToKg[stockUnit]; stockOk {
		if recipeFactor, recipeOk := massToKg[recipeUnit]; recipeOk {
			// conversion_rate = recipe_factor / stock_factor
			return recipeFactor / stockFactor
		}
	}
	
	// Try volume conversion
	if stockFactor, stockOk := volumeToL[stockUnit]; stockOk {
		if recipeFactor, recipeOk := volumeToL[recipeUnit]; recipeOk {
			return recipeFactor / stockFactor
		}
	}
	
	// No conversion found, return 1.0 (assume same unit)
	return 1.0
}

// ValidateUnitConversion checks if conversion between units is valid
// Returns true if both units are in the same category (mass or volume)
func ValidateUnitConversion(stockUnit, recipeUnit UnitType) bool {
	// Same unit is always valid
	if stockUnit == recipeUnit {
		return true
	}
	
	massUnits := map[UnitType]bool{
		UnitKilogram: true,
		UnitGram:     true,
	}
	
	volumeUnits := map[UnitType]bool{
		UnitLiter:      true,
		UnitMilliliter: true,
	}
	
	// Both must be mass or both must be volume
	stockIsMass := massUnits[stockUnit]
	recipeIsMass := massUnits[recipeUnit]
	stockIsVolume := volumeUnits[stockUnit]
	recipeIsVolume := volumeUnits[recipeUnit]
	
	// Valid if both are mass or both are volume
	return (stockIsMass && recipeIsMass) || (stockIsVolume && recipeIsVolume)
}

// GetCompatibleUnits returns all units that can be converted to/from the given unit
func GetCompatibleUnits(unit UnitType) []UnitType {
	massUnits := []UnitType{UnitKilogram, UnitGram}
	volumeUnits := []UnitType{UnitLiter, UnitMilliliter}
	countUnits := []UnitType{UnitPiece, UnitBox, UnitPack}
	
	// Check which category the unit belongs to
	for _, u := range massUnits {
		if u == unit {
			return massUnits
		}
	}
	
	for _, u := range volumeUnits {
		if u == unit {
			return volumeUnits
		}
	}
	
	for _, u := range countUnits {
		if u == unit {
			return countUnits
		}
	}
	
	// Unknown unit, return only itself
	return []UnitType{unit}
}
