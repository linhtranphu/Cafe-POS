package http

import (
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson/primitive"

	"cafe-pos/backend/application/services"
	"cafe-pos/backend/domain/handover"
)

// CashHandoverHandler handles HTTP requests for cash handovers
type CashHandoverHandler struct {
	handoverService *services.CashHandoverService
}

// NewCashHandoverHandler creates a new cash handover handler
func NewCashHandoverHandler(handoverService *services.CashHandoverService) *CashHandoverHandler {
	return &CashHandoverHandler{
		handoverService: handoverService,
	}
}

// CreateHandover creates a new handover request (waiter)
func (h *CashHandoverHandler) CreateHandover(c *gin.Context) {
	shiftID := c.Param("id")
	shiftOID, err := primitive.ObjectIDFromHex(shiftID)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid shift ID"})
		return
	}

	var req handover.CreateHandoverRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	userID, _ := c.Get("user_id")
	username, _ := c.Get("username")

	handoverResult, err := h.handoverService.CreateHandover(
		c.Request.Context(),
		shiftOID,
		&req,
		userID.(string),
		username.(string),
	)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, handoverResult)
}

// CreateHandoverAndEndShift creates handover and ends shift (waiter)
func (h *CashHandoverHandler) CreateHandoverAndEndShift(c *gin.Context) {
	shiftID := c.Param("id")
	shiftOID, err := primitive.ObjectIDFromHex(shiftID)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid shift ID"})
		return
	}

	var req handover.CreateHandoverAndEndShiftRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	userID, _ := c.Get("user_id")
	username, _ := c.Get("username")

	handoverResult, err := h.handoverService.CreateHandoverAndEndShift(
		c.Request.Context(),
		shiftOID,
		&req,
		userID.(string),
		username.(string),
	)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, handoverResult)
}

// GetPendingHandover gets pending handover for a shift (waiter)
func (h *CashHandoverHandler) GetPendingHandover(c *gin.Context) {
	shiftID := c.Param("id")
	shiftOID, err := primitive.ObjectIDFromHex(shiftID)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid shift ID"})
		return
	}

	handoverResult, err := h.handoverService.GetPendingHandover(c.Request.Context(), shiftOID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	if handoverResult == nil {
		c.JSON(http.StatusOK, gin.H{"handover": nil})
		return
	}

	c.JSON(http.StatusOK, handoverResult)
}

// GetHandoverHistory gets handover history for a shift (waiter)
func (h *CashHandoverHandler) GetHandoverHistory(c *gin.Context) {
	shiftID := c.Param("id")
	shiftOID, err := primitive.ObjectIDFromHex(shiftID)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid shift ID"})
		return
	}

	handovers, err := h.handoverService.GetHandoverHistory(c.Request.Context(), shiftOID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, handovers)
}

// CancelHandover cancels a pending handover (waiter)
func (h *CashHandoverHandler) CancelHandover(c *gin.Context) {
	handoverID := c.Param("id")
	handoverOID, err := primitive.ObjectIDFromHex(handoverID)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid handover ID"})
		return
	}

	userID, _ := c.Get("user_id")

	err = h.handoverService.CancelHandover(c.Request.Context(), handoverOID, userID.(string))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "handover cancelled successfully"})
}

// GetPendingHandovers gets all pending handovers for current cashier
func (h *CashHandoverHandler) GetPendingHandovers(c *gin.Context) {
	userID, _ := c.Get("user_id")
	userOID, err := primitive.ObjectIDFromHex(userID.(string))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid user ID"})
		return
	}

	handovers, err := h.handoverService.GetPendingByCashier(c.Request.Context(), userOID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, handovers)
}

// GetAllPendingHandovers gets all pending handovers (any cashier)
func (h *CashHandoverHandler) GetAllPendingHandovers(c *gin.Context) {
	handovers, err := h.handoverService.GetAllPending(c.Request.Context())
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, handovers)
}

// GetTodayHandovers gets today's handovers for current cashier
func (h *CashHandoverHandler) GetTodayHandovers(c *gin.Context) {
	userID, _ := c.Get("user_id")
	userOID, err := primitive.ObjectIDFromHex(userID.(string))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid user ID"})
		return
	}

	handovers, err := h.handoverService.GetTodayByCashier(c.Request.Context(), userOID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, handovers)
}

// ConfirmHandover confirms or rejects a handover with reconciliation (cashier)
func (h *CashHandoverHandler) ConfirmHandover(c *gin.Context) {
	handoverID := c.Param("id")
	handoverOID, err := primitive.ObjectIDFromHex(handoverID)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid handover ID"})
		return
	}

	var req handover.ConfirmHandoverRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	userID, _ := c.Get("user_id")

	err = h.handoverService.ConfirmHandoverWithReconciliation(
		c.Request.Context(),
		handoverOID,
		&req,
		userID.(string),
	)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "handover confirmed successfully"})
}

// QuickConfirm quickly confirms or rejects without detailed reconciliation (cashier)
func (h *CashHandoverHandler) QuickConfirm(c *gin.Context) {
	handoverID := c.Param("id")
	handoverOID, err := primitive.ObjectIDFromHex(handoverID)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid handover ID"})
		return
	}

	var req struct {
		Status handover.HandoverStatus `json:"status" binding:"required"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	userID, _ := c.Get("user_id")

	// Get handover by ID to use declared amount as actual amount
	h_obj, err := h.handoverService.GetHandoverByID(c.Request.Context(), handoverOID)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// For quick confirm, assume declared = actual (no discrepancy)
	confirmReq := &handover.ConfirmHandoverRequest{
		ActualAmount: h_obj.DeclaredAmount,
		Status:       req.Status,
		CashierNote:  "Quick confirm",
	}

	err = h.handoverService.ConfirmHandoverWithReconciliation(
		c.Request.Context(),
		handoverOID,
		confirmReq,
		userID.(string),
	)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "handover confirmed successfully"})
}

// GetPendingApprovals gets handovers requiring manager approval (manager)
func (h *CashHandoverHandler) GetPendingApprovals(c *gin.Context) {
	handovers, err := h.handoverService.GetRequiringApproval(c.Request.Context())
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, handovers)
}

// ApproveDiscrepancy approves or rejects a discrepancy (manager)
func (h *CashHandoverHandler) ApproveDiscrepancy(c *gin.Context) {
	handoverID := c.Param("id")
	handoverOID, err := primitive.ObjectIDFromHex(handoverID)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid handover ID"})
		return
	}

	var req struct {
		Approved bool   `json:"approved" binding:"required"`
		Note     string `json:"note"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	userID, _ := c.Get("user_id")

	err = h.handoverService.ApproveDiscrepancy(
		c.Request.Context(),
		handoverOID,
		userID.(string),
		req.Approved,
		req.Note,
	)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "discrepancy processed successfully"})
}

// GetDiscrepancyStats gets discrepancy statistics (manager)
func (h *CashHandoverHandler) GetDiscrepancyStats(c *gin.Context) {
	startDateStr := c.Query("start_date")
	endDateStr := c.Query("end_date")

	var startDate, endDate time.Time
	var err error

	if startDateStr != "" {
		startDate, err = time.Parse("2006-01-02", startDateStr)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "invalid start_date format"})
			return
		}
	} else {
		// Default to start of today
		now := time.Now()
		startDate = time.Date(now.Year(), now.Month(), now.Day(), 0, 0, 0, 0, now.Location())
	}

	if endDateStr != "" {
		endDate, err = time.Parse("2006-01-02", endDateStr)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "invalid end_date format"})
			return
		}
		// Set to end of day
		endDate = time.Date(endDate.Year(), endDate.Month(), endDate.Day(), 23, 59, 59, 0, endDate.Location())
	} else {
		// Default to end of today
		endDate = time.Now()
	}

	stats, err := h.handoverService.GetDiscrepancyStats(c.Request.Context(), startDate, endDate)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, stats)
}
