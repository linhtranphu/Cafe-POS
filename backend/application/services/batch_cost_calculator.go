package services

import (
	"context"
	"fmt"
	"math"
	"sync"
	"time"

	"cafe-pos/backend/domain/batch"
	"cafe-pos/backend/domain/ingredient"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

// BatchCostCalculator handles cost calculation for batch production
type BatchCostCalculator struct {
	ingredientRepo IngredientRepository
	costCache      *CostCache
	mu             sync.RWMutex
}

// CostCache stores ingredient costs with TTL
type CostCache struct {
	costs map[string]CachedCost
	mu    sync.RWMutex
	ttl   time.Duration
}

// CachedCost represents a cached ingredient cost with expiry
type CachedCost struct {
	Cost      float64
	ExpiresAt time.Time
}

// NewBatchCostCalculator creates a new batch cost calculator
func NewBatchCostCalculator(ingredientRepo IngredientRepository) *BatchCostCalculator {
	return &BatchCostCalculator{
		ingredientRepo: ingredientRepo,
		costCache: &CostCache{
			costs: make(map[string]CachedCost),
			ttl:   5 * time.Minute,
		},
	}
}

// CostBreakdown represents the detailed cost breakdown for a batch
type CostBreakdown struct {
	TotalCost       float64
	CostPerUnit     float64
	IngredientCosts []IngredientCost
}

// IngredientCost represents the cost contribution of a single ingredient
type IngredientCost struct {
	IngredientID   primitive.ObjectID
	IngredientName string
	Quantity       float64
	Unit           string
	CostPerUnit    float64
	TotalCost      float64
}

// CalculateBatchCost calculates the total cost to produce a batch
// Formula: For each ingredient: quantity_needed * cost_per_unit * (1 + wastage_rate)
// where quantity_needed = (batch_quantity / conversion_rate.batch_quantity) * conversion_rate.source_quantity
// IMPORTANT: Handles unit conversion between source_unit and ingredient stock unit
func (c *BatchCostCalculator) CalculateBatchCost(ctx context.Context, def *batch.BatchDefinition, quantityProduced float64) (*CostBreakdown, error) {
	if quantityProduced <= 0 {
		return nil, fmt.Errorf("quantity produced must be greater than 0")
	}

	breakdown := &CostBreakdown{
		IngredientCosts: make([]IngredientCost, 0, len(def.ConversionRates)),
	}

	var totalCost float64

	// Calculate cost for each source ingredient
	for _, convRate := range def.ConversionRates {
		// Calculate how much of this ingredient is needed
		// quantity_needed = (quantityProduced / batch_quantity) * source_quantity
		quantityNeeded := (quantityProduced / convRate.BatchQuantity) * convRate.SourceQuantity

		// Apply wastage rate
		// actual_quantity = quantity_needed * (1 + wastage_rate)
		actualQuantity := quantityNeeded * (1 + convRate.WastageRate)

		// Get ingredient details (with caching)
		ing, costPerUnit, err := c.getIngredientWithCost(ctx, convRate.SourceIngredientID)
		if err != nil {
			return nil, fmt.Errorf("failed to get cost for ingredient %s: %w", convRate.SourceIngredientName, err)
		}

		// Convert quantity from source_unit to ingredient stock unit
		// Example: if source_unit is "g" and stock unit is "kg", convert 500g to 0.5kg
		quantityInStockUnit := c.convertUnit(actualQuantity, convRate.SourceUnit, string(ing.Unit))

		// Calculate total cost for this ingredient using stock unit quantity
		ingredientTotalCost := quantityInStockUnit * costPerUnit
		ingredientTotalCost = math.Round(ingredientTotalCost*100) / 100

		// Add to breakdown (display in source_unit for clarity)
		breakdown.IngredientCosts = append(breakdown.IngredientCosts, IngredientCost{
			IngredientID:   convRate.SourceIngredientID,
			IngredientName: convRate.SourceIngredientName,
			Quantity:       actualQuantity,
			Unit:           convRate.SourceUnit,
			CostPerUnit:    costPerUnit,
			TotalCost:      ingredientTotalCost,
		})

		totalCost += ingredientTotalCost
	}

	// Round total cost to 2 decimal places
	breakdown.TotalCost = math.Round(totalCost*100) / 100

	// Calculate cost per unit
	breakdown.CostPerUnit = math.Round((totalCost/quantityProduced)*100) / 100

	return breakdown, nil
}

// convertUnit converts quantity from one unit to another
// Supports: kg <-> g, L <-> ml
func (c *BatchCostCalculator) convertUnit(quantity float64, fromUnit, toUnit string) float64 {
	// Normalize units to lowercase
	from := normalizeUnit(fromUnit)
	to := normalizeUnit(toUnit)

	// If units are the same, no conversion needed
	if from == to {
		return quantity
	}

	// Define conversion rates (to base unit)
	conversions := map[string]float64{
		"kg": 1.0,
		"g":  0.001,
		"l":  1.0,
		"ml": 0.001,
	}

	// Get conversion rates
	fromRate, fromExists := conversions[from]
	toRate, toExists := conversions[to]

	if !fromExists || !toExists {
		// Unknown units, return original quantity (no conversion)
		return quantity
	}

	// Convert: quantity * (fromRate / toRate)
	// Example: 500g to kg = 500 * (0.001 / 1.0) = 0.5kg
	return quantity * (fromRate / toRate)
}

// normalizeUnit normalizes unit strings to lowercase
func normalizeUnit(unit string) string {
	switch unit {
	case "kg", "Kg", "KG", "kilogram", "kilograms":
		return "kg"
	case "g", "G", "gram", "grams":
		return "g"
	case "l", "L", "liter", "liters", "litre", "litres":
		return "l"
	case "ml", "ML", "milliliter", "milliliters", "millilitre", "millilitres":
		return "ml"
	default:
		return unit
	}
}

// getIngredientCost retrieves ingredient cost with caching
func (c *BatchCostCalculator) getIngredientCost(ctx context.Context, ingredientID primitive.ObjectID) (float64, error) {
	// Try cache first
	if cost, found := c.costCache.Get(ingredientID.Hex()); found {
		return cost, nil
	}

	// Cache miss, query database
	ingredient, err := c.ingredientRepo.FindByID(ctx, ingredientID)
	if err != nil {
		return 0, fmt.Errorf("ingredient not found: %w", err)
	}

	if ingredient.CostPerUnit <= 0 {
		return 0, fmt.Errorf("ingredient %s has no cost set", ingredient.Name)
	}

	// Store in cache
	c.costCache.Set(ingredientID.Hex(), ingredient.CostPerUnit)

	return ingredient.CostPerUnit, nil
}

// getIngredientWithCost retrieves ingredient details and cost with caching
func (c *BatchCostCalculator) getIngredientWithCost(ctx context.Context, ingredientID primitive.ObjectID) (*ingredient.Ingredient, float64, error) {
	// Query database for ingredient details
	ing, err := c.ingredientRepo.FindByID(ctx, ingredientID)
	if err != nil {
		return nil, 0, fmt.Errorf("ingredient not found: %w", err)
	}

	if ing.CostPerUnit <= 0 {
		return nil, 0, fmt.Errorf("ingredient %s has no cost set", ing.Name)
	}

	// Try cache for cost
	if cost, found := c.costCache.Get(ingredientID.Hex()); found {
		return ing, cost, nil
	}

	// Store in cache
	c.costCache.Set(ingredientID.Hex(), ing.CostPerUnit)

	return ing, ing.CostPerUnit, nil
}

// InvalidateCache clears the cost cache (called when ingredient costs change)
func (c *BatchCostCalculator) InvalidateCache() {
	c.costCache.Clear()
}

// InvalidateCacheForIngredient clears cache for a specific ingredient
func (c *BatchCostCalculator) InvalidateCacheForIngredient(ingredientID primitive.ObjectID) {
	c.costCache.Delete(ingredientID.Hex())
}

// CostCache methods

// Get retrieves a cost from cache if not expired
func (cc *CostCache) Get(key string) (float64, bool) {
	cc.mu.RLock()
	defer cc.mu.RUnlock()

	cached, exists := cc.costs[key]
	if !exists {
		return 0, false
	}

	// Check if expired
	if time.Now().After(cached.ExpiresAt) {
		return 0, false
	}

	return cached.Cost, true
}

// Set stores a cost in cache with TTL
func (cc *CostCache) Set(key string, cost float64) {
	cc.mu.Lock()
	defer cc.mu.Unlock()

	cc.costs[key] = CachedCost{
		Cost:      cost,
		ExpiresAt: time.Now().Add(cc.ttl),
	}
}

// Delete removes a specific entry from cache
func (cc *CostCache) Delete(key string) {
	cc.mu.Lock()
	defer cc.mu.Unlock()

	delete(cc.costs, key)
}

// Clear removes all entries from cache
func (cc *CostCache) Clear() {
	cc.mu.Lock()
	defer cc.mu.Unlock()

	cc.costs = make(map[string]CachedCost)
}
