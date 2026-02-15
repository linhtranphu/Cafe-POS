/**
 * Unit Conversion Composable
 * Provides functions to calculate conversion rates between different units
 * Used for ingredient cost calculations in menu items
 */
export function useUnitConversion() {
  // Unit categories
  const massUnits = ['kg', 'g']
  const volumeUnits = ['L', 'ml']
  const countUnits = ['piece', 'box', 'pack']
  
  // Conversion factors to base unit
  const massToKg = {
    'kg': 1.0,
    'g': 0.001
  }
  
  const volumeToL = {
    'L': 1.0,
    'ml': 0.001
  }
  
  /**
   * Calculate conversion rate from stock unit to recipe unit
   * @param {string} stockUnit - Unit used in stock/inventory (e.g., "L")
   * @param {string} recipeUnit - Unit used in recipe (e.g., "ml")
   * @returns {number} - Conversion rate (e.g., 0.001 for ml→L)
   * 
   * Example:
   *   getConversionRate("L", "ml") => 0.001
   *   getConversionRate("kg", "g") => 0.001
   *   getConversionRate("g", "kg") => 0.001
   *   getConversionRate("L", "L") => 1.0
   */
  const getConversionRate = (stockUnit, recipeUnit) => {
    // Same unit, no conversion
    if (stockUnit === recipeUnit) {
      return 1.0
    }
    
    // Try mass conversion
    // To convert from stockUnit to recipeUnit, we need: stockToBase / recipeToBase
    // Example: g→kg = 0.001 / 1.0 = 0.001
    if (massToKg[stockUnit] && massToKg[recipeUnit]) {
      return massToKg[stockUnit] / massToKg[recipeUnit]
    }
    
    // Try volume conversion
    if (volumeToL[stockUnit] && volumeToL[recipeUnit]) {
      return volumeToL[stockUnit] / volumeToL[recipeUnit]
    }
    
    // No conversion found, assume same unit
    return 1.0
  }
  
  /**
   * Validate if conversion between units is possible
   * @param {string} stockUnit - Stock unit
   * @param {string} recipeUnit - Recipe unit
   * @returns {boolean} - True if conversion is valid
   */
  const isValidConversion = (stockUnit, recipeUnit) => {
    if (stockUnit === recipeUnit) return true
    
    const stockIsMass = massUnits.includes(stockUnit)
    const recipeIsMass = massUnits.includes(recipeUnit)
    const stockIsVolume = volumeUnits.includes(stockUnit)
    const recipeIsVolume = volumeUnits.includes(recipeUnit)
    
    // Valid if both are same category
    return (stockIsMass && recipeIsMass) || (stockIsVolume && recipeIsVolume)
  }
  
  /**
   * Get compatible units for a given stock unit
   * @param {string} stockUnit - Stock unit
   * @returns {Array<string>} - Array of compatible units
   */
  const getCompatibleUnits = (stockUnit) => {
    if (massUnits.includes(stockUnit)) return massUnits
    if (volumeUnits.includes(stockUnit)) return volumeUnits
    if (countUnits.includes(stockUnit)) return countUnits
    return [stockUnit]
  }
  
  /**
   * Calculate estimated cost for an ingredient
   * @param {number} quantity - Quantity in recipe unit
   * @param {string} recipeUnit - Recipe unit
   * @param {number} costPerUnit - Cost per stock unit
   * @param {string} stockUnit - Stock unit
   * @param {number} wastage - Wastage percentage (default 0)
   * @returns {number} - Total cost rounded to 2 decimal places
   */
  const calculateCost = (quantity, recipeUnit, costPerUnit, stockUnit, wastage = 0) => {
    // Convert quantity from recipe unit to stock unit
    const conversionRate = getConversionRate(recipeUnit, stockUnit)
    const quantityInStockUnit = quantity * conversionRate
    const cost = quantityInStockUnit * costPerUnit * (1 + wastage / 100)
    return Math.round(cost * 100) / 100
  }
  
  /**
   * Calculate base cost and wastage cost separately
   * @param {number} quantity - Quantity in recipe unit
   * @param {string} recipeUnit - Recipe unit
   * @param {number} costPerUnit - Cost per stock unit
   * @param {string} stockUnit - Stock unit
   * @param {number} wastage - Wastage percentage (default 0)
   * @returns {Object} - { baseCost, wastageCost, totalCost }
   */
  const calculateCostBreakdown = (quantity, recipeUnit, costPerUnit, stockUnit, wastage = 0) => {
    // Convert quantity from recipe unit to stock unit
    // Example: 200g → 0.2kg, so we need rate from g→kg
    const conversionRate = getConversionRate(recipeUnit, stockUnit)
    const quantityInStockUnit = quantity * conversionRate
    const baseCost = quantityInStockUnit * costPerUnit
    const wastageCost = baseCost * (wastage / 100)
    const totalCost = baseCost + wastageCost
    
    return {
      baseCost: Math.round(baseCost * 100) / 100,
      wastageCost: Math.round(wastageCost * 100) / 100,
      totalCost: Math.round(totalCost * 100) / 100
    }
  }
  
  /**
   * Format conversion explanation for display
   * @param {string} stockUnit - Stock unit (unit in inventory)
   * @param {string} recipeUnit - Recipe unit (unit used in recipe)
   * @returns {string} - Human-readable explanation
   */
  const getConversionExplanation = (stockUnit, recipeUnit) => {
    if (stockUnit === recipeUnit) {
      return 'Không cần quy đổi'
    }
    
    const rate = getConversionRate(stockUnit, recipeUnit)
    
    // We always want to show the conversion in a natural way
    // Example: "1kg = 1000g" (not "1g = 0.001kg")
    
    // Determine which unit is larger
    const stockToBase = massToKg[stockUnit] || volumeToL[stockUnit] || 1
    const recipeToBase = massToKg[recipeUnit] || volumeToL[recipeUnit] || 1
    
    if (stockToBase > recipeToBase) {
      // Stock unit is larger (e.g., kg > g)
      // Show: 1kg = Xg
      const ratio = stockToBase / recipeToBase
      return `1${stockUnit} = ${ratio}${recipeUnit}`
    } else if (stockToBase < recipeToBase) {
      // Recipe unit is larger (e.g., g < kg)
      // Show: 1kg = Xg
      const ratio = recipeToBase / stockToBase
      return `1${recipeUnit} = ${ratio}${stockUnit}`
    }
    
    return 'Không cần quy đổi'
  }
  
  /**
   * Get unit display name in Vietnamese
   * @param {string} unit - Unit code
   * @returns {string} - Display name
   */
  const getUnitDisplayName = (unit) => {
    const displayNames = {
      'kg': 'Kilogram',
      'g': 'Gram',
      'L': 'Lít',
      'ml': 'Milliliter',
      'piece': 'Cái',
      'box': 'Hộp',
      'pack': 'Gói'
    }
    return displayNames[unit] || unit
  }
  
  return {
    getConversionRate,
    isValidConversion,
    getCompatibleUnits,
    calculateCost,
    calculateCostBreakdown,
    getConversionExplanation,
    getUnitDisplayName
  }
}
