package http

import (
	"net/http"
	"time"

	"cafe-pos/backend/application/services"
	"github.com/gin-gonic/gin"
)

// ProfitAnalysisHandler handles profit analysis API endpoints
type ProfitAnalysisHandler struct {
	profitAnalyzerService *services.ProfitAnalyzerService
}

// NewProfitAnalysisHandler creates a new profit analysis handler
func NewProfitAnalysisHandler(profitAnalyzerService *services.ProfitAnalyzerService) *ProfitAnalysisHandler {
	return &ProfitAnalysisHandler{
		profitAnalyzerService: profitAnalyzerService,
	}
}

// GetCategoryProfit handles GET /api/reports/category-profit
// Requirements: 6.1, 6.4, 7.1
func (h *ProfitAnalysisHandler) GetCategoryProfit(c *gin.Context) {
	// Get query parameters
	startDateStr := c.Query("start_date")
	endDateStr := c.Query("end_date")

	// Validate required parameters
	if startDateStr == "" || endDateStr == "" {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "start_date and end_date are required",
		})
		return
	}

	// Parse dates (ISO 8601 format)
	startDate, err := time.Parse("2006-01-02", startDateStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "invalid start_date format, expected YYYY-MM-DD",
		})
		return
	}

	endDate, err := time.Parse("2006-01-02", endDateStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "invalid end_date format, expected YYYY-MM-DD",
		})
		return
	}

	// Validate date range
	if startDate.After(endDate) {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "start_date must be before or equal to end_date",
		})
		return
	}

	// Set time to start and end of day
	startDate = time.Date(startDate.Year(), startDate.Month(), startDate.Day(), 0, 0, 0, 0, startDate.Location())
	endDate = time.Date(endDate.Year(), endDate.Month(), endDate.Day(), 23, 59, 59, 999999999, endDate.Location())

	// Create date range
	dateRange := services.DateRange{
		Start: startDate,
		End:   endDate,
	}

	// Call service to get category profits
	categories, err := h.profitAnalyzerService.GetCategoryProfits(c.Request.Context(), dateRange)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "failed to calculate category profits",
		})
		return
	}

	// Return response
	c.JSON(http.StatusOK, gin.H{
		"date_range": gin.H{
			"start": startDateStr,
			"end":   endDateStr,
		},
		"categories": categories,
	})
}


// GetOperatingProfit handles GET /api/reports/operating-profit
// Requirements: 6.5.1, 6.5.6, 6.5.9
func (h *ProfitAnalysisHandler) GetOperatingProfit(c *gin.Context) {
	// Get query parameters
	startDateStr := c.Query("start_date")
	endDateStr := c.Query("end_date")

	// Validate required parameters
	if startDateStr == "" || endDateStr == "" {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "start_date and end_date are required",
		})
		return
	}

	// Parse dates (ISO 8601 format)
	startDate, err := time.Parse("2006-01-02", startDateStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "invalid start_date format, expected YYYY-MM-DD",
		})
		return
	}

	endDate, err := time.Parse("2006-01-02", endDateStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "invalid end_date format, expected YYYY-MM-DD",
		})
		return
	}

	// Validate date range
	if startDate.After(endDate) {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "start_date must be before or equal to end_date",
		})
		return
	}

	// Set time to start and end of day
	startDate = time.Date(startDate.Year(), startDate.Month(), startDate.Day(), 0, 0, 0, 0, startDate.Location())
	endDate = time.Date(endDate.Year(), endDate.Month(), endDate.Day(), 23, 59, 59, 999999999, endDate.Location())

	// Create date range
	dateRange := services.DateRange{
		Start: startDate,
		End:   endDate,
	}

	// Call service to get operating profit
	report, err := h.profitAnalyzerService.GetOperatingProfit(c.Request.Context(), dateRange)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "failed to calculate operating profit",
		})
		return
	}

	// Return response
	c.JSON(http.StatusOK, report)
}
