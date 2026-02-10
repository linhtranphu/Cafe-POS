package http

import (
	"net/http"
	"strconv"
	"time"

	"cafe-pos/backend/application/services"
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
