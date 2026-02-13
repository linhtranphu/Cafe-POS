package services

import (
	"context"
	"testing"

	"cafe-pos/backend/domain/ingredient"
	"cafe-pos/backend/domain/menu"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

// Test 15.3: Cost calculation edge cases
// Requirements: FR-6.4-FR-6.9, NFR-5.4-NFR-5.7

// TestCostCalculation_MissingIngredientData tests cost calculation with missing ingredient data
// Should set INCOMPLETE status when ingredient data is missing
func TestCostCalculation_MissingIngredientData(t *testing.T) {
	tests := []struct {
		name           string
		menuItem       *menu.MenuItem
		ingredients    []*ingredient.Ingredient
		expectedStatus menu.CostStatus
		expectedCost   float64
		description    string
	}{
		{
			name: "Missing ingredient in database",
			menuItem: &menu.MenuItem{
				ID:          primitive.NewObjectID(),
				Name:        "Coffee with missing ingredient",
				HasVariants: false,
				Ingredients: []menu.Ingredient{
					{Name: "Coffee Beans", Quantity: 20, Unit: ingredient.UnitGram},
					{Name: "Missing Ingredient", Quantity: 10, Unit: ingredient.UnitGram},
				},
			},
			ingredients: []*ingredient.Ingredient{
				{
					ID:          primitive.NewObjectID(),
					Name:        "Coffee Beans",
					Unit:        ingredient.UnitGram,
					CostPerUnit: 500,
				},
				// Missing Ingredient not in database
			},
			expectedStatus: menu.CostStatusIncomplete,
			expectedCost:   10000, // Only coffee beans: 20 * 500
			description:    "Should mark as INCOMPLETE when ingredient not found in database",
		},
		{
			name: "Ingredient with zero cost_per_unit",
			menuItem: &menu.MenuItem{
				ID:          primitive.NewObjectID(),
				Name:        "Latte with zero cost ingredient",
				HasVariants: false,
				Ingredients: []menu.Ingredient{
					{Name: "Espresso", Quantity: 30, Unit: ingredient.UnitMilliliter},
					{Name: "Milk", Quantity: 150, Unit: ingredient.UnitMilliliter},
				},
			},
			ingredients: []*ingredient.Ingredient{
				{
					ID:          primitive.NewObjectID(),
					Name:        "Espresso",
					Unit:        ingredient.UnitMilliliter,
					CostPerUnit: 0, // Zero cost
				},
				{
					ID:          primitive.NewObjectID(),
					Name:        "Milk",
					Unit:        ingredient.UnitMilliliter,
					CostPerUnit: 50,
				},
			},
			expectedStatus: menu.CostStatusIncomplete,
			expectedCost:   7500, // Only milk: 150 * 50
			description:    "Should mark as INCOMPLETE when ingredient has zero cost_per_unit",
		},
		{
			name: "All ingredients missing cost data",
			menuItem: &menu.MenuItem{
				ID:          primitive.NewObjectID(),
				Name:        "Item with all missing costs",
				HasVariants: false,
				Ingredients: []menu.Ingredient{
					{Name: "Ingredient1", Quantity: 10, Unit: ingredient.UnitGram},
					{Name: "Ingredient2", Quantity: 20, Unit: ingredient.UnitGram},
				},
			},
			ingredients: []*ingredient.Ingredient{
				{
					ID:          primitive.NewObjectID(),
					Name:        "Ingredient1",
					Unit:        ingredient.UnitGram,
					CostPerUnit: 0,
				},
				{
					ID:          primitive.NewObjectID(),
					Name:        "Ingredient2",
					Unit:        ingredient.UnitGram,
					CostPerUnit: 0,
				},
			},
			expectedStatus: menu.CostStatusIncomplete,
			expectedCost:   0,
			description:    "Should mark as INCOMPLETE with zero cost when all ingredients missing cost",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Setup
			menuRepo := &mockMenuRepository{
				menuItems: map[primitive.ObjectID]*menu.MenuItem{
					tt.menuItem.ID: tt.menuItem,
				},
			}
			ingredientRepo := &mockIngredientRepository{
				ingredients: make(map[primitive.ObjectID]*ingredient.Ingredient),
			}
			for _, ing := range tt.ingredients {
				ingredientRepo.ingredients[ing.ID] = ing
			}

			service := NewCostCalculatorService(menuRepo, ingredientRepo, nil, nil)

			// Execute
			result, err := service.CalculateMenuItemCost(context.Background(), tt.menuItem.ID)

			// Verify
			require.NoError(t, err, tt.description)
			assert.Equal(t, tt.expectedStatus, result.CostStatus, tt.description)
			assert.Equal(t, tt.expectedCost, result.CurrentCost, tt.description)
			assert.NotEmpty(t, result.MissingIngredients, "Should list missing ingredients")
		})
	}
}

// TestCostCalculation_ZeroWastagePercentage tests cost calculation with zero wastage
// Should calculate cost correctly without wastage multiplier
func TestCostCalculation_ZeroWastagePercentage(t *testing.T) {
	tests := []struct {
		name         string
		menuItem     *menu.MenuItem
		ingredients  []*ingredient.Ingredient
		expectedCost float64
		description  string
	}{
		{
			name: "Zero wastage percentage",
			menuItem: &menu.MenuItem{
				ID:          primitive.NewObjectID(),
				Name:        "Coffee with zero wastage",
				HasVariants: false,
				Ingredients: []menu.Ingredient{
					{Name: "Coffee Beans", Quantity: 20, Unit: ingredient.UnitGram},
				},
			},
			ingredients: []*ingredient.Ingredient{
				{
					ID:                primitive.NewObjectID(),
					Name:              "Coffee Beans",
					Unit:              ingredient.UnitGram,
					CostPerUnit:       500,
					WastagePercentage: 0, // Zero wastage
				},
			},
			expectedCost: 10000, // 20 * 500 * 1.0 * (1 + 0/100) = 10000
			description:  "Should calculate cost without wastage multiplier when wastage is 0",
		},
		{
			name: "Negative wastage percentage (treated as zero)",
			menuItem: &menu.MenuItem{
				ID:          primitive.NewObjectID(),
				Name:        "Item with negative wastage",
				HasVariants: false,
				Ingredients: []menu.Ingredient{
					{Name: "Ingredient", Quantity: 10, Unit: ingredient.UnitGram},
				},
			},
			ingredients: []*ingredient.Ingredient{
				{
					ID:                primitive.NewObjectID(),
					Name:              "Ingredient",
					Unit:              ingredient.UnitGram,
					CostPerUnit:       100,
					WastagePercentage: -5, // Negative wastage (invalid, should be treated as 0)
				},
			},
			expectedCost: 1000, // 10 * 100 * 1.0 * (1 + 0/100) = 1000
			description:  "Should treat negative wastage as zero",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Setup
			menuRepo := &mockMenuRepository{
				menuItems: map[primitive.ObjectID]*menu.MenuItem{
					tt.menuItem.ID: tt.menuItem,
				},
			}
			ingredientRepo := &mockIngredientRepository{
				ingredients: make(map[primitive.ObjectID]*ingredient.Ingredient),
			}
			for _, ing := range tt.ingredients {
				ingredientRepo.ingredients[ing.ID] = ing
			}

			service := NewCostCalculatorService(menuRepo, ingredientRepo, nil, nil)

			// Execute
			result, err := service.CalculateMenuItemCost(context.Background(), tt.menuItem.ID)

			// Verify
			require.NoError(t, err, tt.description)
			assert.Equal(t, menu.CostStatusFinal, result.CostStatus, tt.description)
			assert.Equal(t, tt.expectedCost, result.CurrentCost, tt.description)
		})
	}
}

// TestCostCalculation_HighWastagePercentage tests cost calculation with high wastage (>50%)
// Should apply wastage multiplier correctly even with high percentages
func TestCostCalculation_HighWastagePercentage(t *testing.T) {
	tests := []struct {
		name         string
		menuItem     *menu.MenuItem
		ingredients  []*ingredient.Ingredient
		expectedCost float64
		description  string
	}{
		{
			name: "50% wastage percentage",
			menuItem: &menu.MenuItem{
				ID:          primitive.NewObjectID(),
				Name:        "Item with 50% wastage",
				HasVariants: false,
				Ingredients: []menu.Ingredient{
					{Name: "Ingredient", Quantity: 10, Unit: ingredient.UnitGram},
				},
			},
			ingredients: []*ingredient.Ingredient{
				{
					ID:                primitive.NewObjectID(),
					Name:              "Ingredient",
					Unit:              ingredient.UnitGram,
					CostPerUnit:       100,
					WastagePercentage: 50, // 50% wastage
				},
			},
			expectedCost: 1500, // 10 * 100 * 1.0 * (1 + 50/100) = 1500
			description:  "Should apply 50% wastage multiplier correctly",
		},
		{
			name: "75% wastage percentage",
			menuItem: &menu.MenuItem{
				ID:          primitive.NewObjectID(),
				Name:        "Item with 75% wastage",
				HasVariants: false,
				Ingredients: []menu.Ingredient{
					{Name: "Ingredient", Quantity: 20, Unit: ingredient.UnitGram},
				},
			},
			ingredients: []*ingredient.Ingredient{
				{
					ID:                primitive.NewObjectID(),
					Name:              "Ingredient",
					Unit:              ingredient.UnitGram,
					CostPerUnit:       200,
					WastagePercentage: 75, // 75% wastage
				},
			},
			expectedCost: 7000, // 20 * 200 * 1.0 * (1 + 75/100) = 7000
			description:  "Should apply 75% wastage multiplier correctly",
		},
		{
			name: "100% wastage percentage (doubles cost)",
			menuItem: &menu.MenuItem{
				ID:          primitive.NewObjectID(),
				Name:        "Item with 100% wastage",
				HasVariants: false,
				Ingredients: []menu.Ingredient{
					{Name: "Ingredient", Quantity: 15, Unit: ingredient.UnitGram},
				},
			},
			ingredients: []*ingredient.Ingredient{
				{
					ID:                primitive.NewObjectID(),
					Name:              "Ingredient",
					Unit:              ingredient.UnitGram,
					CostPerUnit:       300,
					WastagePercentage: 100, // 100% wastage (doubles the cost)
				},
			},
			expectedCost: 9000, // 15 * 300 * 1.0 * (1 + 100/100) = 9000
			description:  "Should double cost with 100% wastage",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Setup
			menuRepo := &mockMenuRepository{
				menuItems: map[primitive.ObjectID]*menu.MenuItem{
					tt.menuItem.ID: tt.menuItem,
				},
			}
			ingredientRepo := &mockIngredientRepository{
				ingredients: make(map[primitive.ObjectID]*ingredient.Ingredient),
			}
			for _, ing := range tt.ingredients {
				ingredientRepo.ingredients[ing.ID] = ing
			}

			service := NewCostCalculatorService(menuRepo, ingredientRepo, nil, nil)

			// Execute
			result, err := service.CalculateMenuItemCost(context.Background(), tt.menuItem.ID)

			// Verify
			require.NoError(t, err, tt.description)
			assert.Equal(t, menu.CostStatusFinal, result.CostStatus, tt.description)
			assert.Equal(t, tt.expectedCost, result.CurrentCost, tt.description)
		})
	}
}

// TestCostCalculation_UnitConversion tests cost calculation with unit conversions
// Should apply conversion rates correctly (g to kg, ml to L, etc.)
func TestCostCalculation_UnitConversion(t *testing.T) {
	tests := []struct {
		name         string
		menuItem     *menu.MenuItem
		ingredients  []*ingredient.Ingredient
		expectedCost float64
		description  string
	}{
		{
			name: "Gram to Kilogram conversion",
			menuItem: &menu.MenuItem{
				ID:          primitive.NewObjectID(),
				Name:        "Item with g to kg conversion",
				HasVariants: false,
				Ingredients: []menu.Ingredient{
					{Name: "Flour", Quantity: 500, Unit: ingredient.UnitGram}, // Recipe uses grams
				},
			},
			ingredients: []*ingredient.Ingredient{
				{
					ID:          primitive.NewObjectID(),
					Name:        "Flour",
					Unit:        ingredient.UnitKilogram, // Stock in kilograms
					CostPerUnit: 20000,                   // 20,000 per kg
				},
			},
			expectedCost: 10000, // 500g = 0.5kg, 0.5 * 20000 = 10000
			description:  "Should convert grams to kilograms correctly",
		},
		{
			name: "Milliliter to Liter conversion",
			menuItem: &menu.MenuItem{
				ID:          primitive.NewObjectID(),
				Name:        "Item with ml to L conversion",
				HasVariants: false,
				Ingredients: []menu.Ingredient{
					{Name: "Milk", Quantity: 250, Unit: ingredient.UnitMilliliter}, // Recipe uses ml
				},
			},
			ingredients: []*ingredient.Ingredient{
				{
					ID:          primitive.NewObjectID(),
					Name:        "Milk",
					Unit:        ingredient.UnitLiter, // Stock in liters
					CostPerUnit: 30000,                // 30,000 per liter
				},
			},
			expectedCost: 7500, // 250ml = 0.25L, 0.25 * 30000 = 7500
			description:  "Should convert milliliters to liters correctly",
		},
		{
			name: "Kilogram to Gram conversion",
			menuItem: &menu.MenuItem{
				ID:          primitive.NewObjectID(),
				Name:        "Item with kg to g conversion",
				HasVariants: false,
				Ingredients: []menu.Ingredient{
					{Name: "Sugar", Quantity: 2, Unit: ingredient.UnitKilogram}, // Recipe uses kg
				},
			},
			ingredients: []*ingredient.Ingredient{
				{
					ID:          primitive.NewObjectID(),
					Name:        "Sugar",
					Unit:        ingredient.UnitGram, // Stock in grams
					CostPerUnit: 10,                  // 10 per gram
				},
			},
			expectedCost: 20000, // 2kg = 2000g, 2000 * 10 = 20000
			description:  "Should convert kilograms to grams correctly",
		},
		{
			name: "Liter to Milliliter conversion",
			menuItem: &menu.MenuItem{
				ID:          primitive.NewObjectID(),
				Name:        "Item with L to ml conversion",
				HasVariants: false,
				Ingredients: []menu.Ingredient{
					{Name: "Water", Quantity: 1.5, Unit: ingredient.UnitLiter}, // Recipe uses liters
				},
			},
			ingredients: []*ingredient.Ingredient{
				{
					ID:          primitive.NewObjectID(),
					Name:        "Water",
					Unit:        ingredient.UnitMilliliter, // Stock in ml
					CostPerUnit: 5,                         // 5 per ml
				},
			},
			expectedCost: 7500, // 1.5L = 1500ml, 1500 * 5 = 7500
			description:  "Should convert liters to milliliters correctly",
		},
		{
			name: "Same unit (no conversion)",
			menuItem: &menu.MenuItem{
				ID:          primitive.NewObjectID(),
				Name:        "Item with same unit",
				HasVariants: false,
				Ingredients: []menu.Ingredient{
					{Name: "Coffee", Quantity: 30, Unit: ingredient.UnitGram},
				},
			},
			ingredients: []*ingredient.Ingredient{
				{
					ID:          primitive.NewObjectID(),
					Name:        "Coffee",
					Unit:        ingredient.UnitGram, // Same unit
					CostPerUnit: 500,
				},
			},
			expectedCost: 15000, // 30 * 500 * 1.0 = 15000
			description:  "Should handle same unit without conversion",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Setup
			menuRepo := &mockMenuRepository{
				menuItems: map[primitive.ObjectID]*menu.MenuItem{
					tt.menuItem.ID: tt.menuItem,
				},
			}
			ingredientRepo := &mockIngredientRepository{
				ingredients: make(map[primitive.ObjectID]*ingredient.Ingredient),
			}
			for _, ing := range tt.ingredients {
				ingredientRepo.ingredients[ing.ID] = ing
			}

			service := NewCostCalculatorService(menuRepo, ingredientRepo, nil, nil)

			// Execute
			result, err := service.CalculateMenuItemCost(context.Background(), tt.menuItem.ID)

			// Verify
			require.NoError(t, err, tt.description)
			assert.Equal(t, menu.CostStatusFinal, result.CostStatus, tt.description)
			assert.Equal(t, tt.expectedCost, result.CurrentCost, tt.description)
		})
	}
}

// TestCostCalculation_DecimalPrecision tests cost accuracy to 2 decimal places
// Should round costs correctly to 2 decimal places
func TestCostCalculation_DecimalPrecision(t *testing.T) {
	tests := []struct {
		name         string
		menuItem     *menu.MenuItem
		ingredients  []*ingredient.Ingredient
		expectedCost float64
		description  string
	}{
		{
			name: "Cost with many decimal places",
			menuItem: &menu.MenuItem{
				ID:          primitive.NewObjectID(),
				Name:        "Item with precise cost",
				HasVariants: false,
				Ingredients: []menu.Ingredient{
					{Name: "Ingredient", Quantity: 7.333, Unit: ingredient.UnitGram},
				},
			},
			ingredients: []*ingredient.Ingredient{
				{
					ID:          primitive.NewObjectID(),
					Name:        "Ingredient",
					Unit:        ingredient.UnitGram,
					CostPerUnit: 123.456,
				},
			},
			expectedCost: 905.30, // 7.333 * 123.456 = 905.277648, rounded to 905.30 (rounds up)
			description:  "Should round to 2 decimal places",
		},
		{
			name: "Cost with conversion and wastage (complex rounding)",
			menuItem: &menu.MenuItem{
				ID:          primitive.NewObjectID(),
				Name:        "Complex rounding item",
				HasVariants: false,
				Ingredients: []menu.Ingredient{
					{Name: "Ingredient", Quantity: 333, Unit: ingredient.UnitGram},
				},
			},
			ingredients: []*ingredient.Ingredient{
				{
					ID:                primitive.NewObjectID(),
					Name:              "Ingredient",
					Unit:              ingredient.UnitKilogram,
					CostPerUnit:       15789.50,
					WastagePercentage: 12.5,
				},
			},
			expectedCost: 5915.14, // 333g = 0.333kg, 0.333 * 15789.50 * 1.125 = 5915.136375, rounded to 5915.14
			description:  "Should round complex calculation to 2 decimal places",
		},
		{
			name: "Multiple ingredients with rounding",
			menuItem: &menu.MenuItem{
				ID:          primitive.NewObjectID(),
				Name:        "Multiple ingredients rounding",
				HasVariants: false,
				Ingredients: []menu.Ingredient{
					{Name: "Ingredient1", Quantity: 10.5, Unit: ingredient.UnitGram},
					{Name: "Ingredient2", Quantity: 20.3, Unit: ingredient.UnitGram},
				},
			},
			ingredients: []*ingredient.Ingredient{
				{
					ID:          primitive.NewObjectID(),
					Name:        "Ingredient1",
					Unit:        ingredient.UnitGram,
					CostPerUnit: 99.99,
				},
				{
					ID:          primitive.NewObjectID(),
					Name:        "Ingredient2",
					Unit:        ingredient.UnitGram,
					CostPerUnit: 77.77,
				},
			},
			expectedCost: 2628.63, // (10.5 * 99.99) + (20.3 * 77.77) = 1049.895 + 1578.731 = 2628.626, rounded to 2628.63
			description:  "Should round sum of multiple ingredients to 2 decimal places",
		},
		{
			name: "Cost rounds down",
			menuItem: &menu.MenuItem{
				ID:          primitive.NewObjectID(),
				Name:        "Round down item",
				HasVariants: false,
				Ingredients: []menu.Ingredient{
					{Name: "Ingredient", Quantity: 10, Unit: ingredient.UnitGram},
				},
			},
			ingredients: []*ingredient.Ingredient{
				{
					ID:          primitive.NewObjectID(),
					Name:        "Ingredient",
					Unit:        ingredient.UnitGram,
					CostPerUnit: 100.004, // Should round down
				},
			},
			expectedCost: 1000.04, // 10 * 100.004 = 1000.04
			description:  "Should round down when third decimal < 5",
		},
		{
			name: "Cost rounds up",
			menuItem: &menu.MenuItem{
				ID:          primitive.NewObjectID(),
				Name:        "Round up item",
				HasVariants: false,
				Ingredients: []menu.Ingredient{
					{Name: "Ingredient", Quantity: 10, Unit: ingredient.UnitGram},
				},
			},
			ingredients: []*ingredient.Ingredient{
				{
					ID:          primitive.NewObjectID(),
					Name:        "Ingredient",
					Unit:        ingredient.UnitGram,
					CostPerUnit: 100.006, // Should round up
				},
			},
			expectedCost: 1000.06, // 10 * 100.006 = 1000.06
			description:  "Should round up when third decimal >= 5",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Setup
			menuRepo := &mockMenuRepository{
				menuItems: map[primitive.ObjectID]*menu.MenuItem{
					tt.menuItem.ID: tt.menuItem,
				},
			}
			ingredientRepo := &mockIngredientRepository{
				ingredients: make(map[primitive.ObjectID]*ingredient.Ingredient),
			}
			for _, ing := range tt.ingredients {
				ingredientRepo.ingredients[ing.ID] = ing
			}

			service := NewCostCalculatorService(menuRepo, ingredientRepo, nil, nil)

			// Execute
			result, err := service.CalculateMenuItemCost(context.Background(), tt.menuItem.ID)

			// Verify
			require.NoError(t, err, tt.description)
			assert.Equal(t, menu.CostStatusFinal, result.CostStatus, tt.description)
			assert.InDelta(t, tt.expectedCost, result.CurrentCost, 0.01, tt.description)
		})
	}
}

// TestCostCalculation_MultiSizeWithEdgeCases tests multi-size items with edge cases
// Should handle edge cases correctly for variant cost calculation
func TestCostCalculation_MultiSizeWithEdgeCases(t *testing.T) {
	tests := []struct {
		name                string
		menuItem            *menu.MenuItem
		ingredients         []*ingredient.Ingredient
		expectedVariantCost map[string]float64 // variant ID -> expected cost
		expectedStatus      map[string]menu.CostStatus
		description         string
	}{
		{
			name: "Multi-size with missing ingredient data",
			menuItem: &menu.MenuItem{
				ID:          primitive.NewObjectID(),
				Name:        "Coffee with variants and missing data",
				HasVariants: true,
				Variants: []menu.MenuItemVariant{
					{
						ID:        "M",
						Name:      "Size M",
						Price:     25000,
						IsDefault: true,
						Ingredients: []menu.Ingredient{
							{Name: "Coffee", Quantity: 20, Unit: ingredient.UnitGram},
							{Name: "Missing", Quantity: 10, Unit: ingredient.UnitGram},
						},
					},
					{
						ID:        "L",
						Name:      "Size L",
						Price:     30000,
						IsDefault: false,
						Ingredients: []menu.Ingredient{
							{Name: "Coffee", Quantity: 30, Unit: ingredient.UnitGram},
						},
					},
				},
			},
			ingredients: []*ingredient.Ingredient{
				{
					ID:          primitive.NewObjectID(),
					Name:        "Coffee",
					Unit:        ingredient.UnitGram,
					CostPerUnit: 500,
				},
			},
			expectedVariantCost: map[string]float64{
				"M": 10000, // Only coffee: 20 * 500
				"L": 15000, // 30 * 500
			},
			expectedStatus: map[string]menu.CostStatus{
				"M": menu.CostStatusIncomplete, // Missing ingredient
				"L": menu.CostStatusFinal,      // All ingredients present
			},
			description: "Should handle missing ingredient in one variant",
		},
		{
			name: "Multi-size with high wastage",
			menuItem: &menu.MenuItem{
				ID:          primitive.NewObjectID(),
				Name:        "Item with high wastage variants",
				HasVariants: true,
				Variants: []menu.MenuItemVariant{
					{
						ID:        "M",
						Name:      "Size M",
						Price:     20000,
						IsDefault: true,
						Ingredients: []menu.Ingredient{
							{Name: "Ingredient", Quantity: 10, Unit: ingredient.UnitGram},
						},
					},
					{
						ID:        "L",
						Name:      "Size L",
						Price:     30000,
						IsDefault: false,
						Ingredients: []menu.Ingredient{
							{Name: "Ingredient", Quantity: 20, Unit: ingredient.UnitGram},
						},
					},
				},
			},
			ingredients: []*ingredient.Ingredient{
				{
					ID:                primitive.NewObjectID(),
					Name:              "Ingredient",
					Unit:              ingredient.UnitGram,
					CostPerUnit:       100,
					WastagePercentage: 80, // 80% wastage
				},
			},
			expectedVariantCost: map[string]float64{
				"M": 1800, // 10 * 100 * (1 + 80/100) = 1800
				"L": 3600, // 20 * 100 * (1 + 80/100) = 3600
			},
			expectedStatus: map[string]menu.CostStatus{
				"M": menu.CostStatusFinal,
				"L": menu.CostStatusFinal,
			},
			description: "Should apply high wastage to all variants",
		},
		{
			name: "Multi-size with unit conversion",
			menuItem: &menu.MenuItem{
				ID:          primitive.NewObjectID(),
				Name:        "Item with conversion in variants",
				HasVariants: true,
				Variants: []menu.MenuItemVariant{
					{
						ID:        "M",
						Name:      "Size M",
						Price:     25000,
						IsDefault: true,
						Ingredients: []menu.Ingredient{
							{Name: "Milk", Quantity: 200, Unit: ingredient.UnitMilliliter},
						},
					},
					{
						ID:        "L",
						Name:      "Size L",
						Price:     35000,
						IsDefault: false,
						Ingredients: []menu.Ingredient{
							{Name: "Milk", Quantity: 350, Unit: ingredient.UnitMilliliter},
						},
					},
				},
			},
			ingredients: []*ingredient.Ingredient{
				{
					ID:          primitive.NewObjectID(),
					Name:        "Milk",
					Unit:        ingredient.UnitLiter, // Stock in liters
					CostPerUnit: 30000,                // 30,000 per liter
				},
			},
			expectedVariantCost: map[string]float64{
				"M": 6000,  // 200ml = 0.2L, 0.2 * 30000 = 6000
				"L": 10500, // 350ml = 0.35L, 0.35 * 30000 = 10500
			},
			expectedStatus: map[string]menu.CostStatus{
				"M": menu.CostStatusFinal,
				"L": menu.CostStatusFinal,
			},
			description: "Should apply unit conversion to all variants",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Setup
			menuRepo := &mockMenuRepository{
				menuItems: map[primitive.ObjectID]*menu.MenuItem{
					tt.menuItem.ID: tt.menuItem,
				},
			}
			ingredientRepo := &mockIngredientRepository{
				ingredients: make(map[primitive.ObjectID]*ingredient.Ingredient),
			}
			for _, ing := range tt.ingredients {
				ingredientRepo.ingredients[ing.ID] = ing
			}

			service := NewCostCalculatorService(menuRepo, ingredientRepo, nil, nil)

			// Execute
			_, err := service.CalculateMenuItemCost(context.Background(), tt.menuItem.ID)

			// Verify
			require.NoError(t, err, tt.description)

			// Fetch updated menu item
			updatedItem, err := menuRepo.FindByID(context.Background(), tt.menuItem.ID)
			require.NoError(t, err)

			// Verify each variant
			for _, variant := range updatedItem.Variants {
				expectedCost, ok := tt.expectedVariantCost[variant.ID]
				require.True(t, ok, "Expected cost for variant %s", variant.ID)
				assert.Equal(t, expectedCost, variant.CurrentCost, "Cost for variant %s: %s", variant.ID, tt.description)

				expectedStatus, ok := tt.expectedStatus[variant.ID]
				require.True(t, ok, "Expected status for variant %s", variant.ID)
				assert.Equal(t, expectedStatus, variant.CostStatus, "Status for variant %s: %s", variant.ID, tt.description)
			}
		})
	}
}
