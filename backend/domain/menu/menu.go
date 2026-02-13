package menu

import (
	"fmt"
	"time"
	"cafe-pos/backend/domain/ingredient"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

// CostStatus represents the status of cost calculation for a menu item
type CostStatus string

const (
	CostStatusFinal      CostStatus = "FINAL"      // Cost has been calculated and finalized
	CostStatusEstimated  CostStatus = "ESTIMATED"  // Cost is estimated (shift not closed)
	CostStatusIncomplete CostStatus = "INCOMPLETE" // Missing ingredient cost data
)

type Ingredient struct {
	Name     string               `bson:"name" json:"name"`
	Quantity float64              `bson:"quantity" json:"quantity"`
	Unit     ingredient.UnitType  `bson:"unit" json:"unit"`
}

// MenuItemVariant represents a size variant of a menu item
type MenuItemVariant struct {
	ID          string       `bson:"id" json:"id"`                     // e.g., "M", "L", "XL"
	Name        string       `bson:"name" json:"name"`                 // e.g., "Size M"
	Price       float64      `bson:"price" json:"price"`
	Ingredients []Ingredient `bson:"ingredients" json:"ingredients"`
	Available   bool         `bson:"available" json:"available"`
	IsDefault   bool         `bson:"is_default" json:"is_default"`
	
	// Cost tracking per variant
	CurrentCost          float64    `bson:"current_cost" json:"current_cost"`
	CostLastCalculatedAt time.Time  `bson:"cost_last_calculated_at" json:"cost_last_calculated_at"`
	CostStatus           CostStatus `bson:"cost_status" json:"cost_status"`
}

type MenuItem struct {
	ID          primitive.ObjectID `bson:"_id,omitempty" json:"id"`
	Name        string             `bson:"name" json:"name"`
	Category    string             `bson:"category" json:"category"`
	Description string             `bson:"description" json:"description"`
	Available   bool               `bson:"available" json:"available"`
	
	// Variants support (NEW)
	HasVariants bool               `bson:"has_variants" json:"has_variants"` // Default: false
	Variants    []MenuItemVariant  `bson:"variants,omitempty" json:"variants,omitempty"`
	
	// Backward compatibility (single-size items) - KEEP these fields
	Price       float64      `bson:"price,omitempty" json:"price,omitempty"`
	Ingredients []Ingredient `bson:"ingredients,omitempty" json:"ingredients,omitempty"`
	
	// Cost tracking fields (for single-size items)
	CurrentCost          float64    `bson:"current_cost,omitempty" json:"current_cost,omitempty"`
	CostLastCalculatedAt time.Time  `bson:"cost_last_calculated_at,omitempty" json:"cost_last_calculated_at,omitempty"`
	CostStatus           CostStatus `bson:"cost_status,omitempty" json:"cost_status,omitempty"`
	
	CreatedAt   time.Time `bson:"created_at" json:"created_at"`
	UpdatedAt   time.Time `bson:"updated_at" json:"updated_at"`
}

type CreateMenuItemRequest struct {
	Name        string             `json:"name" binding:"required"`
	Category    string             `json:"category" binding:"required"`
	Description string             `json:"description"`
	
	// Variants support (NEW)
	HasVariants bool               `json:"has_variants"`
	Variants    []MenuItemVariant  `json:"variants,omitempty"`
	
	// Backward compatibility (single-size items)
	Price       float64      `json:"price,omitempty" binding:"omitempty,min=0"`
	Ingredients []Ingredient `json:"ingredients,omitempty"`
}

type UpdateMenuItemRequest struct {
	Name        string       `json:"name"`
	Category    string       `json:"category"`
	Description string       `json:"description"`
	Available   *bool        `json:"available"`
	
	// Variants support (NEW)
	HasVariants *bool              `json:"has_variants"`
	Variants    []MenuItemVariant  `json:"variants,omitempty"`
	
	// Backward compatibility (single-size items)
	Price       float64      `json:"price,omitempty" binding:"omitempty,min=0"`
	Ingredients []Ingredient `json:"ingredients,omitempty"`
}

// MenuItem Helper Methods

// GetDefaultVariant returns the default variant or nil if no variants
func (m *MenuItem) GetDefaultVariant() *MenuItemVariant {
	if !m.HasVariants {
		return nil
	}
	for i := range m.Variants {
		if m.Variants[i].IsDefault {
			return &m.Variants[i]
		}
	}
	// Fallback to first variant if no default set
	if len(m.Variants) > 0 {
		return &m.Variants[0]
	}
	return nil
}

// GetVariantByID returns the variant with the given ID or nil if not found
func (m *MenuItem) GetVariantByID(variantID string) *MenuItemVariant {
	if !m.HasVariants {
		return nil
	}
	for i := range m.Variants {
		if m.Variants[i].ID == variantID {
			return &m.Variants[i]
		}
	}
	return nil
}

// GetPrice returns the price for the item, optionally for a specific variant
func (m *MenuItem) GetPrice(variantID string) float64 {
	if m.HasVariants {
		if variantID != "" {
			variant := m.GetVariantByID(variantID)
			if variant != nil {
				return variant.Price
			}
		}
		// Fallback to default variant
		defaultVariant := m.GetDefaultVariant()
		if defaultVariant != nil {
			return defaultVariant.Price
		}
		return 0
	}
	return m.Price
}

// GetIngredients returns the ingredients for the item, optionally for a specific variant
func (m *MenuItem) GetIngredients(variantID string) []Ingredient {
	if m.HasVariants {
		variant := m.GetVariantByID(variantID)
		if variant != nil {
			return variant.Ingredients
		}
		return nil
	}
	return m.Ingredients
}

// Validate validates the MenuItem data integrity
func (m *MenuItem) Validate() error {
	// Basic validation
	if m.Name == "" {
		return fmt.Errorf("name is required")
	}
	if m.Category == "" {
		return fmt.Errorf("category is required")
	}
	
	if m.HasVariants {
		// Multi-size validation
		if len(m.Variants) == 0 {
			return fmt.Errorf("variants required when has_variants=true")
		}
		
		// Must have exactly 1 default
		defaultCount := 0
		for _, v := range m.Variants {
			if v.IsDefault {
				defaultCount++
			}
		}
		if defaultCount == 0 {
			return fmt.Errorf("must have at least one default variant")
		}
		if defaultCount > 1 {
			return fmt.Errorf("must have exactly one default variant")
		}
		
		// Variant IDs must be unique
		ids := make(map[string]bool)
		for _, v := range m.Variants {
			if v.ID == "" {
				return fmt.Errorf("variant ID is required")
			}
			if ids[v.ID] {
				return fmt.Errorf("duplicate variant ID: %s", v.ID)
			}
			ids[v.ID] = true
			
			// Variant must have valid price
			if v.Price <= 0 {
				return fmt.Errorf("variant %s price must be > 0", v.ID)
			}
		}
		
		// IMPORTANT: When has_variants=true, old fields should be empty
		// This prevents data inconsistency
		if m.Price > 0 {
			return fmt.Errorf("price should not be set when has_variants=true (use variants instead)")
		}
		if len(m.Ingredients) > 0 {
			return fmt.Errorf("ingredients should not be set when has_variants=true (use variants instead)")
		}
	} else {
		// Single-size validation (backward compatible)
		if m.Price <= 0 {
			return fmt.Errorf("price must be > 0 for single-size item")
		}
		
		// IMPORTANT: When has_variants=false, variants should be empty
		if len(m.Variants) > 0 {
			return fmt.Errorf("variants should not be set when has_variants=false")
		}
	}
	
	return nil
}