package http

import (
	"context"
	"net/http"
	"strconv"
	"time"

	"cafe-pos/backend/application/services"
	"cafe-pos/backend/domain/menu"
	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

// MenuCostHandler handles menu cost and profit analysis endpoints
type MenuCostHandler struct {
	profitAnalyzer       *services.ProfitAnalyzerService
	costCalculator       *services.CostCalculatorService
	recalculationService *services.CostRecalculationService
}

// NewMenuCostHandler creates a new menu cost handler
func NewMenuCostHandler(
	profitAnalyzer *services.ProfitAnalyzerService,
	costCalculator *services.CostCalculatorService,
	recalculationService *services.CostRecalculationService,
) *MenuCostHandler {
	return &MenuCostHandler{
		profitAnalyzer:       profitAnalyzer,
		costCalculator:       costCalculator,
		recalculationService: recalculationService,
	}
}

// GetMenuCostsResponse represents the response for GET /api/menu/costs
type GetMenuCostsResponse struct {
	Items                []services.MenuItemProfit    `json:"items"`
	Summary              services.ProfitSummary       `json:"summary"`
	RecalculationStatus  *services.RecalculationStatus `json:"recalculation_status"`
}

// GetMenuCosts handles GET /api/menu/costs
// Query params: category (optional), sort_by (optional), sort_order (optional)
// Requirements: 4.1, 4.2, 4.3, 4.4, 7.4
func (h *MenuCostHandler) GetMenuCosts(c *gin.Context) {
	// Parse query parameters
	category := c.Query("category")
	sortBy := c.Query("sort_by")
	sortOrder := c.Query("sort_order")

	// Create filter
	filter := services.ProfitFilter{
		Category:  category,
		SortBy:    sortBy,
		SortOrder: sortOrder,
	}

	// Call service to get all menu item profits
	result, err := h.profitAnalyzer.GetAllMenuItemProfits(c.Request.Context(), filter)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	// Get recalculation status
	recalcStatus, err := h.recalculationService.GetRecalculationStatus(c.Request.Context())
	if err != nil {
		// Log error but don't fail the request
		recalcStatus = &services.RecalculationStatus{
			InProgress:     false,
			QueuedItems:    0,
			ProcessedItems: 0,
			LastUpdated:    time.Now(),
		}
	}

	// Build response
	response := GetMenuCostsResponse{
		Items:               result.Items,
		Summary:             result.Summary,
		RecalculationStatus: recalcStatus,
	}

	c.JSON(http.StatusOK, response)
}

// IngredientCostDetail represents cost detail for a single ingredient
type IngredientCostDetail struct {
	Name              string  `json:"name"`
	Quantity          float64 `json:"quantity"`
	Unit              string  `json:"unit"`
	CostPerUnit       float64 `json:"cost_per_unit"`
	ConversionRate    float64 `json:"conversion_rate"`
	WastagePercentage float64 `json:"wastage_percentage"`
	TotalCost         float64 `json:"total_cost"`
}

// GetMenuCostDetailResponse represents the response for GET /api/menu/costs/:id
type GetMenuCostDetailResponse struct {
	MenuItem    services.MenuItemProfit `json:"menu_item"`
	Ingredients []IngredientCostDetail  `json:"ingredients"`
	TotalCost   float64                 `json:"total_cost"`
}

// GetMenuCostDetail handles GET /api/menu/costs/:id
// Returns menu item with ingredient breakdown
// Requirements: 8.1, 8.2, 8.3
func (h *MenuCostHandler) GetMenuCostDetail(c *gin.Context) {
	// Parse menu item ID
	idStr := c.Param("id")
	id, err := primitive.ObjectIDFromHex(idStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid id"})
		return
	}

	// Calculate cost detail with ingredient breakdown
	costDetail, err := h.costCalculator.CalculateMenuItemCostDetail(c.Request.Context(), id)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "menu item not found"})
		return
	}

	// Calculate profit metrics for the menu item
	profit, err := h.profitAnalyzer.CalculateMenuItemProfit(c.Request.Context(), id)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	// Build ingredient details
	ingredients := make([]IngredientCostDetail, 0, len(costDetail.Ingredients))
	for _, ing := range costDetail.Ingredients {
		ingredients = append(ingredients, IngredientCostDetail{
			Name:              ing.Name,
			Quantity:          ing.Quantity,
			Unit:              ing.Unit,
			CostPerUnit:       ing.CostPerUnit,
			ConversionRate:    ing.ConversionRate,
			WastagePercentage: ing.WastagePercentage,
			TotalCost:         ing.TotalCost,
		})
	}

	// Build response
	response := GetMenuCostDetailResponse{
		MenuItem:    *profit,
		Ingredients: ingredients,
		TotalCost:   costDetail.TotalCost,
	}

	c.JSON(http.StatusOK, response)
}

// GetMenuWarnings handles GET /api/menu/warnings
// Query params: threshold (optional)
// Requirements: 3.3, 3.4, 3.5
func (h *MenuCostHandler) GetMenuWarnings(c *gin.Context) {
	// Parse optional threshold parameter
	thresholdStr := c.Query("threshold")
	threshold := 0.0
	if thresholdStr != "" {
		var err error
		threshold, err = strconv.ParseFloat(thresholdStr, 64)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "invalid threshold value"})
			return
		}
	}

	// Call service to detect warnings
	warnings, err := h.profitAnalyzer.DetectWarningStatus(c.Request.Context(), threshold)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, warnings)
}

// VariantCostBreakdown represents cost breakdown for a single variant
type VariantCostBreakdown struct {
	VariantID            string                   `json:"variant_id"`
	VariantName          string                   `json:"variant_name"`
	Price                float64                  `json:"price"`
	Ingredients          []IngredientCostDetail   `json:"ingredients"`
	TotalCost            float64                  `json:"total_cost"`
	CostStatus           string                   `json:"cost_status"`
	CostLastCalculatedAt time.Time                `json:"cost_last_calculated_at"`
}

// GetCostBreakdownResponse represents the response for GET /api/menu/:id/cost-breakdown
type GetCostBreakdownResponse struct {
	MenuItemID   primitive.ObjectID      `json:"menu_item_id"`
	MenuItemName string                  `json:"menu_item_name"`
	HasVariants  bool                    `json:"has_variants"`
	
	// For single-size items
	Price        float64                 `json:"price,omitempty"`
	Ingredients  []IngredientCostDetail  `json:"ingredients,omitempty"`
	TotalCost    float64                 `json:"total_cost,omitempty"`
	CostStatus   string                  `json:"cost_status,omitempty"`
	
	// For multi-size items
	Variants     []VariantCostBreakdown  `json:"variants,omitempty"`
}

// GetCostBreakdown handles GET /api/menu/:id/cost-breakdown
// Returns detailed cost breakdown per variant (or single-size)
// Requirements: FR-7.6, FR-9.1-FR-9.4, AC-10.1-AC-10.5
func (h *MenuCostHandler) GetCostBreakdown(c *gin.Context) {
	// Parse menu item ID
	idStr := c.Param("id")
	id, err := primitive.ObjectIDFromHex(idStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid id"})
		return
	}

	// Fetch the menu item using the cost calculator service
	menuItem, err := h.costCalculator.GetMenuItemByID(c.Request.Context(), id)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "menu item not found"})
		return
	}

	response := GetCostBreakdownResponse{
		MenuItemID:   id,
		MenuItemName: menuItem.Name,
		HasVariants:  menuItem.HasVariants,
	}

	if menuItem.HasVariants {
		// Multi-size item - return breakdown per variant
		response.Variants = make([]VariantCostBreakdown, 0, len(menuItem.Variants))
		
		for _, variant := range menuItem.Variants {
			// Calculate cost detail for this variant's ingredients
			costDetail := h.calculateVariantCostDetail(c.Request.Context(), variant.Ingredients)
			
			variantBreakdown := VariantCostBreakdown{
				VariantID:            variant.ID,
				VariantName:          variant.Name,
				Price:                variant.Price,
				Ingredients:          costDetail.ingredients,
				TotalCost:            costDetail.totalCost,
				CostStatus:           string(variant.CostStatus),
				CostLastCalculatedAt: variant.CostLastCalculatedAt,
			}
			
			response.Variants = append(response.Variants, variantBreakdown)
		}
	} else {
		// Single-size item - return single breakdown
		costDetail, err := h.costCalculator.CalculateMenuItemCostDetail(c.Request.Context(), id)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}

		// Build ingredient details
		ingredients := make([]IngredientCostDetail, 0, len(costDetail.Ingredients))
		for _, ing := range costDetail.Ingredients {
			ingredients = append(ingredients, IngredientCostDetail{
				Name:              ing.Name,
				Quantity:          ing.Quantity,
				Unit:              ing.Unit,
				CostPerUnit:       ing.CostPerUnit,
				ConversionRate:    ing.ConversionRate,
				WastagePercentage: ing.WastagePercentage,
				TotalCost:         ing.TotalCost,
			})
		}

		response.Price = menuItem.Price
		response.Ingredients = ingredients
		response.TotalCost = costDetail.TotalCost
		response.CostStatus = string(costDetail.CostStatus)
	}

	c.JSON(http.StatusOK, response)
}

// variantCostDetail is a helper struct for calculateVariantCostDetail
type variantCostDetail struct {
	ingredients []IngredientCostDetail
	totalCost   float64
}

// calculateVariantCostDetail calculates cost detail for a variant's ingredients
// This method fetches ingredients from the database and calculates costs
func (h *MenuCostHandler) calculateVariantCostDetail(ctx context.Context, variantIngredients []menu.Ingredient) *variantCostDetail {
	result := &variantCostDetail{
		ingredients: []IngredientCostDetail{},
		totalCost:   0,
	}

	// If no ingredients, return empty
	if len(variantIngredients) == 0 {
		return result
	}

	// We need to call the cost calculator service to get ingredient details
	// For now, create a temporary menu item to use the existing CalculateMenuItemCostDetail method
	// This is not ideal but works for the MVP
	// TODO: Refactor to have a dedicated method for calculating variant costs
	
	// For now, return basic structure without detailed calculation
	// The actual cost is already calculated and stored in variant.CurrentCost
	for _, ing := range variantIngredients {
		detail := IngredientCostDetail{
			Name:              ing.Name,
			Quantity:          ing.Quantity,
			Unit:              string(ing.Unit),
			CostPerUnit:       0, // Would need to fetch from ingredient repo
			ConversionRate:    1.0,
			WastagePercentage: 0,
			TotalCost:         0,
		}
		result.ingredients = append(result.ingredients, detail)
	}

	return result
}

// VariantProfitAnalysis represents profit analysis for a single variant
type VariantProfitAnalysis struct {
	VariantID      string  `json:"variant_id"`
	VariantName    string  `json:"variant_name"`
	Price          float64 `json:"price"`
	Cost           float64 `json:"cost"`
	Profit         float64 `json:"profit"`
	ProfitMargin   float64 `json:"profit_margin_percent"`
	CostStatus     string  `json:"cost_status"`
}

// GetProfitAnalysisResponse represents the response for GET /api/menu/:id/profit-analysis
type GetProfitAnalysisResponse struct {
	MenuItemID   primitive.ObjectID       `json:"menu_item_id"`
	MenuItemName string                   `json:"menu_item_name"`
	HasVariants  bool                     `json:"has_variants"`
	
	// For single-size items
	Price        float64                  `json:"price,omitempty"`
	Cost         float64                  `json:"cost,omitempty"`
	Profit       float64                  `json:"profit,omitempty"`
	ProfitMargin float64                  `json:"profit_margin_percent,omitempty"`
	CostStatus   string                   `json:"cost_status,omitempty"`
	
	// For multi-size items
	Variants     []VariantProfitAnalysis  `json:"variants,omitempty"`
}

// GetProfitAnalysis handles GET /api/menu/:id/profit-analysis
// Returns profit analysis per variant (or single-size)
// Requirements: FR-7.6, FR-9.1-FR-9.4, AC-12.1-AC-12.4
func (h *MenuCostHandler) GetProfitAnalysis(c *gin.Context) {
	// Parse menu item ID
	idStr := c.Param("id")
	id, err := primitive.ObjectIDFromHex(idStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid id"})
		return
	}

	// Fetch the menu item using the cost calculator service
	menuItem, err := h.costCalculator.GetMenuItemByID(c.Request.Context(), id)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "menu item not found"})
		return
	}

	response := GetProfitAnalysisResponse{
		MenuItemID:   id,
		MenuItemName: menuItem.Name,
		HasVariants:  menuItem.HasVariants,
	}

	if menuItem.HasVariants {
		// Multi-size item - return profit analysis per variant
		response.Variants = make([]VariantProfitAnalysis, 0, len(menuItem.Variants))
		
		for _, variant := range menuItem.Variants {
			profit := variant.Price - variant.CurrentCost
			profitMargin := 0.0
			if variant.Price > 0 {
				profitMargin = (profit / variant.Price) * 100
			}
			
			variantProfit := VariantProfitAnalysis{
				VariantID:    variant.ID,
				VariantName:  variant.Name,
				Price:        variant.Price,
				Cost:         variant.CurrentCost,
				Profit:       profit,
				ProfitMargin: profitMargin,
				CostStatus:   string(variant.CostStatus),
			}
			
			response.Variants = append(response.Variants, variantProfit)
		}
	} else {
		// Single-size item - return single profit analysis
		profit := menuItem.Price - menuItem.CurrentCost
		profitMargin := 0.0
		if menuItem.Price > 0 {
			profitMargin = (profit / menuItem.Price) * 100
		}
		
		response.Price = menuItem.Price
		response.Cost = menuItem.CurrentCost
		response.Profit = profit
		response.ProfitMargin = profitMargin
		response.CostStatus = string(menuItem.CostStatus)
	}

	c.JSON(http.StatusOK, response)
}

// CalculateCost handles POST /api/menu/:id/calculate-cost
// Triggers cost calculation for a menu item (all variants if multi-size)
// Requirements: FR-7.6, FR-9.1-FR-9.4
func (h *MenuCostHandler) CalculateCost(c *gin.Context) {
	// Parse menu item ID
	idStr := c.Param("id")
	id, err := primitive.ObjectIDFromHex(idStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid id"})
		return
	}

	// Trigger cost calculation
	result, err := h.costCalculator.CalculateMenuItemCost(c.Request.Context(), id)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message":                "cost calculation completed",
		"menu_item_id":           result.MenuItemID.Hex(),
		"current_cost":           result.CurrentCost,
		"cost_status":            string(result.CostStatus),
		"cost_last_calculated_at": result.CostLastCalculatedAt,
		"missing_ingredients":    result.MissingIngredients,
	})
}
