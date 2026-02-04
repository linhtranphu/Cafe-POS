package http

import (
	"net/http"
	"strconv"
	"time"

	"cafe-pos/backend/application/services"
	"cafe-pos/backend/domain/handover"
	"cafe-pos/backend/domain/user"
	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type CashHandoverHandler struct {
	service *services.CashHandoverService
}

func NewCashHandoverHandler(service *services.CashHandoverService) *CashHandoverHandler {
	return &CashHandoverHandler{service: service}
}

// CreateHandover creates a new cash handover request
// POST /api/shifts/:id/handover
func (h *CashHandoverHandler) CreateHandover(c *gin.Context) {
	// Get shift ID from URL
	shiftIDStr := c.Param("id")
	shiftID, err := primitive.ObjectIDFromHex(shiftIDStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid shift ID"})
		return
	}

	// Get user info from context
	userID, exists := c.Get("user_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "User not authenticated"})
		return
	}

	username, exists := c.Get("username")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Username not found"})
		return
	}

	// Parse request
	var req handover.CreateHandoverRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Convert user ID
	waiterID, err := primitive.ObjectIDFromHex(userID.(string))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid user ID"})
		return
	}

	// Create handover
	handoverResult, err := h.service.CreateHandover(
		c.Request.Context(),
		shiftID,
		&req,
		waiterID,
		username.(string),
	)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, handover.NewHandoverResponse(handoverResult))
}

// CreateHandoverAndEndShift creates a handover and ends the shift
// POST /api/shifts/:id/handover-and-end
func (h *CashHandoverHandler) CreateHandoverAndEndShift(c *gin.Context) {
	// Get shift ID from URL
	shiftIDStr := c.Param("id")
	shiftID, err := primitive.ObjectIDFromHex(shiftIDStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid shift ID"})
		return
	}

	// Get user info from context
	userID, exists := c.Get("user_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "User not authenticated"})
		return
	}

	username, exists := c.Get("username")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Username not found"})
		return
	}

	// Parse request
	var req handover.CreateHandoverRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Convert user ID
	waiterID, err := primitive.ObjectIDFromHex(userID.(string))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid user ID"})
		return
	}

	// Create handover and end shift
	handoverResult, err := h.service.CreateHandoverAndEndShift(
		c.Request.Context(),
		shiftID,
		&req,
		waiterID,
		username.(string),
	)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, handover.NewHandoverResponse(handoverResult))
}

// GetPendingHandover gets pending handover for a shift
// GET /api/shifts/:id/pending-handover
func (h *CashHandoverHandler) GetPendingHandover(c *gin.Context) {
	// Get shift ID from URL
	shiftIDStr := c.Param("id")
	shiftID, err := primitive.ObjectIDFromHex(shiftIDStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid shift ID"})
		return
	}

	// Get handover history (we'll filter for pending in the response)
	handovers, err := h.service.GetHandoverHistory(c.Request.Context(), shiftID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	// Find pending handover
	for _, h := range handovers {
		if h.Status == handover.HandoverStatusPending {
			c.JSON(http.StatusOK, handover.NewHandoverResponse(h))
			return
		}
	}

	c.JSON(http.StatusNotFound, gin.H{"message": "No pending handover found"})
}

// GetHandoverHistory gets handover history for a shift
// GET /api/shifts/:id/handovers
func (h *CashHandoverHandler) GetHandoverHistory(c *gin.Context) {
	// Get shift ID from URL
	shiftIDStr := c.Param("id")
	shiftID, err := primitive.ObjectIDFromHex(shiftIDStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid shift ID"})
		return
	}

	// Get handover history
	handovers, err := h.service.GetHandoverHistory(c.Request.Context(), shiftID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	// Convert to response format
	responses := make([]*handover.HandoverResponse, len(handovers))
	for i, h := range handovers {
		responses[i] = handover.NewHandoverResponse(h)
	}

	c.JSON(http.StatusOK, gin.H{
		"handovers": responses,
		"count":     len(responses),
	})
}

// CancelHandover cancels a pending handover
// DELETE /api/cash-handovers/:id
func (h *CashHandoverHandler) CancelHandover(c *gin.Context) {
	// Get handover ID from URL
	handoverIDStr := c.Param("id")
	handoverID, err := primitive.ObjectIDFromHex(handoverIDStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid handover ID"})
		return
	}

	// Get user info from context
	userID, exists := c.Get("user_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "User not authenticated"})
		return
	}

	// Convert user ID
	waiterID, err := primitive.ObjectIDFromHex(userID.(string))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid user ID"})
		return
	}

	// Cancel handover
	if err := h.service.CancelHandover(c.Request.Context(), handoverID, waiterID); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Handover cancelled successfully"})
}

// GetPendingHandovers gets pending handovers for cashier
// GET /api/cash-handovers/pending
func (h *CashHandoverHandler) GetPendingHandovers(c *gin.Context) {
	// Get user info from context
	userID, exists := c.Get("user_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "User not authenticated"})
		return
	}

	// Convert user ID
	cashierID, err := primitive.ObjectIDFromHex(userID.(string))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid user ID"})
		return
	}

	// Get pending handovers
	handovers, err := h.service.GetPendingHandovers(c.Request.Context(), cashierID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	// Convert to response format
	responses := make([]*handover.HandoverResponse, len(handovers))
	for i, h := range handovers {
		responses[i] = handover.NewHandoverResponse(h)
	}

	c.JSON(http.StatusOK, gin.H{
		"handovers": responses,
		"count":     len(responses),
	})
}

// GetTodayHandovers gets today's handovers
// GET /api/cash-handovers/today
func (h *CashHandoverHandler) GetTodayHandovers(c *gin.Context) {
	// Get today's handovers
	handovers, err := h.service.GetTodayHandovers(c.Request.Context())
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	// Convert to response format
	responses := make([]*handover.HandoverResponse, len(handovers))
	for i, h := range handovers {
		responses[i] = handover.NewHandoverResponse(h)
	}

	c.JSON(http.StatusOK, gin.H{
		"handovers": responses,
		"count":     len(responses),
	})
}

// ReconcileHandover reconciles a handover with discrepancy details
// POST /api/cash-handovers/:id/reconcile
func (h *CashHandoverHandler) ReconcileHandover(c *gin.Context) {
	// Get handover ID from URL
	handoverIDStr := c.Param("id")
	handoverID, err := primitive.ObjectIDFromHex(handoverIDStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid handover ID"})
		return
	}

	// Get user info from context
	userID, exists := c.Get("user_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "User not authenticated"})
		return
	}

	username, exists := c.Get("username")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Username not found"})
		return
	}

	// Parse request
	var req handover.ReconcileHandoverRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Convert user ID
	cashierID, err := primitive.ObjectIDFromHex(userID.(string))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid user ID"})
		return
	}

	// For reconciliation, we need to get the cashier's current shift
	// This is a simplified approach - in a real system, you might want to pass the shift ID
	cashierShiftID := primitive.NewObjectID() // This should be retrieved from current open shift

	// Convert reconcile request to confirm request
	confirmReq := &handover.ConfirmHandoverRequest{
		ActualAmount: req.ActualAmount,
		CashierNotes: req.CashierNotes,
	}

	// Confirm handover with reconciliation
	if err := h.service.ConfirmHandoverWithReconciliation(
		c.Request.Context(),
		handoverID,
		confirmReq,
		cashierShiftID,
		cashierID,
		username.(string),
	); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Handover reconciled successfully"})
}

// QuickConfirm quickly confirms a handover with exact amount
// POST /api/cash-handovers/:id/quick-confirm
func (h *CashHandoverHandler) QuickConfirm(c *gin.Context) {
	// Get handover ID from URL
	handoverIDStr := c.Param("id")
	handoverID, err := primitive.ObjectIDFromHex(handoverIDStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid handover ID"})
		return
	}

	// Get user info from context
	userID, exists := c.Get("user_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "User not authenticated"})
		return
	}

	username, exists := c.Get("username")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Username not found"})
		return
	}

	// Parse request (optional notes)
	var req struct {
		CashierNotes string `json:"cashier_notes,omitempty"`
	}
	c.ShouldBindJSON(&req)

	// Convert user ID
	cashierID, err := primitive.ObjectIDFromHex(userID.(string))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid user ID"})
		return
	}

	// For quick confirm, we need to get the handover first to use its requested amount
	// This is a simplified approach - in a real system, you might want to pass the shift ID
	cashierShiftID := primitive.NewObjectID() // This should be retrieved from current open shift

	// Get the handover to use its requested amount as actual amount
	handoverService := h.service
	handovers, err := handoverService.GetPendingHandovers(c.Request.Context(), cashierID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	var targetHandover *handover.CashHandover
	for _, h := range handovers {
		if h.ID == handoverID {
			targetHandover = h
			break
		}
	}

	if targetHandover == nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Handover not found"})
		return
	}

	// Create confirm request with exact amount
	confirmReq := &handover.ConfirmHandoverRequest{
		ActualAmount: targetHandover.RequestedAmount, // Exact amount
		CashierNotes: req.CashierNotes,
		QuickConfirm: true,
	}

	// Confirm handover
	if err := h.service.ConfirmHandoverWithReconciliation(
		c.Request.Context(),
		handoverID,
		confirmReq,
		cashierShiftID,
		cashierID,
		username.(string),
	); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Handover confirmed successfully"})
}

// RejectHandover rejects a handover request
// POST /api/cash-handovers/:id/reject
func (h *CashHandoverHandler) RejectHandover(c *gin.Context) {
	// Get handover ID from URL
	handoverIDStr := c.Param("id")
	handoverID, err := primitive.ObjectIDFromHex(handoverIDStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid handover ID"})
		return
	}

	// Get user info from context
	userID, exists := c.Get("user_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "User not authenticated"})
		return
	}

	username, exists := c.Get("username")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Username not found"})
		return
	}

	// Parse request
	var req struct {
		Reason string `json:"reason" binding:"required"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Convert user ID
	cashierID, err := primitive.ObjectIDFromHex(userID.(string))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid user ID"})
		return
	}

	// For rejection, we need the cashier's current shift
	cashierShiftID := primitive.NewObjectID() // This should be retrieved from current open shift

	// Reject handover
	if err := h.service.RejectHandover(
		c.Request.Context(),
		handoverID,
		req.Reason,
		cashierShiftID,
		cashierID,
		username.(string),
	); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Handover rejected successfully"})
}

// GetDiscrepancyStats gets discrepancy statistics
// GET /api/cash-handovers/discrepancy-stats
func (h *CashHandoverHandler) GetDiscrepancyStats(c *gin.Context) {
	// Parse query parameters for date range
	startDateStr := c.Query("start_date")
	endDateStr := c.Query("end_date")

	var startDate, endDate time.Time
	var err error

	if startDateStr != "" {
		startDate, err = time.Parse("2006-01-02", startDateStr)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid start_date format. Use YYYY-MM-DD"})
			return
		}
	} else {
		// Default to start of current month
		now := time.Now()
		startDate = time.Date(now.Year(), now.Month(), 1, 0, 0, 0, 0, now.Location())
	}

	if endDateStr != "" {
		endDate, err = time.Parse("2006-01-02", endDateStr)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid end_date format. Use YYYY-MM-DD"})
			return
		}
	} else {
		// Default to end of today
		now := time.Now()
		endDate = time.Date(now.Year(), now.Month(), now.Day(), 23, 59, 59, 0, now.Location())
	}

	// Get discrepancy stats
	stats, err := h.service.GetDiscrepancyStats(c.Request.Context(), startDate, endDate)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"stats":      stats,
		"start_date": startDate.Format("2006-01-02"),
		"end_date":   endDate.Format("2006-01-02"),
	})
}

// Manager-specific endpoints

// GetPendingApprovals gets handovers requiring manager approval
// GET /api/cash-handovers/pending-approval
func (h *CashHandoverHandler) GetPendingApprovals(c *gin.Context) {
	// Get handovers requiring manager approval
	handovers, err := h.service.GetPendingApprovals(c.Request.Context())
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	// Convert to response format
	responses := make([]*handover.HandoverResponse, len(handovers))
	for i, h := range handovers {
		responses[i] = handover.NewHandoverResponse(h)
	}

	c.JSON(http.StatusOK, gin.H{
		"handovers": responses,
		"count":     len(responses),
	})
}

// ApproveDiscrepancy approves or rejects a discrepancy
// POST /api/cash-handovers/:id/approve
func (h *CashHandoverHandler) ApproveDiscrepancy(c *gin.Context) {
	// Get handover ID from URL
	handoverIDStr := c.Param("id")
	handoverID, err := primitive.ObjectIDFromHex(handoverIDStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid handover ID"})
		return
	}

	// Get user info from context
	userID, exists := c.Get("user_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "User not authenticated"})
		return
	}

	// Parse request
	var req handover.ApproveDiscrepancyRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Convert user ID
	managerID, err := primitive.ObjectIDFromHex(userID.(string))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid user ID"})
		return
	}

	// Approve discrepancy
	if err := h.service.ApproveDiscrepancy(
		c.Request.Context(),
		handoverID,
		managerID,
		req.Approved,
		req.ManagerNotes,
	); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	action := "rejected"
	if req.Approved {
		action = "approved"
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "Discrepancy " + action + " successfully",
	})
}