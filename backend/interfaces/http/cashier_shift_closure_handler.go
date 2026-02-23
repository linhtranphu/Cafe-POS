package http

import (
	"net/http"
	"time"

	"cafe-pos/backend/application/services"
	"cafe-pos/backend/domain"
	"cafe-pos/backend/domain/cashier"

	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

// CashierShiftClosureHandler handles HTTP requests for cashier shift closure workflow
type CashierShiftClosureHandler struct {
	cashierShiftService *services.CashierShiftService
	stateMachineManager *domain.StateMachineManager
}

// NewCashierShiftClosureHandler creates a new handler
func NewCashierShiftClosureHandler(
	cashierShiftService *services.CashierShiftService,
	stateMachineManager *domain.StateMachineManager,
) *CashierShiftClosureHandler {
	return &CashierShiftClosureHandler{
		cashierShiftService: cashierShiftService,
		stateMachineManager: stateMachineManager,
	}
}

// InitiateClosureRequest represents the request to initiate shift closure
type InitiateClosureRequest struct {
	UserID   string `json:"user_id"`
	DeviceID string `json:"device_id"`
}

// RecordActualCashRequest represents the request to record actual cash
type RecordActualCashRequest struct {
	ActualCash float64 `json:"actual_cash" binding:"required,min=0"`
	UserID     string  `json:"user_id"`
	DeviceID   string  `json:"device_id"`
}

// DocumentVarianceRequest represents the request to document variance
type DocumentVarianceRequest struct {
	Reason   string `json:"reason" binding:"required"`
	Notes    string `json:"notes" binding:"required,min=10"`
	UserID   string `json:"user_id"`
	DeviceID string `json:"device_id"`
}

// ConfirmResponsibilityRequest represents the request to confirm responsibility
type ConfirmResponsibilityRequest struct {
	UserID   string `json:"user_id"`
	DeviceID string `json:"device_id"`
}

// CloseShiftRequest represents the request to close the shift
type CloseShiftRequest struct {
	UserID   string `json:"user_id"`
	DeviceID string `json:"device_id"`
}

// CompleteClosureRequest represents the request to complete entire closure workflow
type CompleteClosureRequest struct {
	ActualCash     float64  `json:"actual_cash" binding:"required,min=0"`
	VarianceReason *string  `json:"variance_reason"`
	VarianceNotes  *string  `json:"variance_notes"`
	UserID         string   `json:"user_id"`
	DeviceID       string   `json:"device_id"`
}

// InitiateClosure starts the shift closure process
func (h *CashierShiftClosureHandler) InitiateClosure(c *gin.Context) {
	shiftID := c.Param("id")
	
	shiftObjID, err := primitive.ObjectIDFromHex(shiftID)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid shift ID"})
		return
	}

	// Get user_id from JWT token
	userID := c.GetString("user_id")
	if userID == "" {
		userID = "system"
	}
	deviceID := "web"

	// Use transaction-wrapped method
	shift, err := h.cashierShiftService.InitiateClosureWithTransaction(
		c.Request.Context(),
		shiftObjID,
		userID,
		deviceID,
	)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, shift)
}

// RecordActualCash records the actual cash counted
func (h *CashierShiftClosureHandler) RecordActualCash(c *gin.Context) {
	shiftID := c.Param("id")
	
	shiftObjID, err := primitive.ObjectIDFromHex(shiftID)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid shift ID"})
		return
	}

	var req RecordActualCashRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Get user_id from JWT token
	userID := c.GetString("user_id")
	if userID == "" {
		userID = "system"
	}
	deviceID := "web"

	// Use transaction-wrapped method
	shift, variance, err := h.cashierShiftService.RecordActualCashWithTransaction(
		c.Request.Context(),
		shiftObjID,
		req.ActualCash,
		userID,
		deviceID,
	)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"shift":    shift,
		"variance": variance,
	})
}

// DocumentVariance documents the reason for variance
func (h *CashierShiftClosureHandler) DocumentVariance(c *gin.Context) {
	shiftID := c.Param("id")
	
	shiftObjID, err := primitive.ObjectIDFromHex(shiftID)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid shift ID"})
		return
	}

	var req DocumentVarianceRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Get user_id from JWT token
	userID := c.GetString("user_id")
	if userID == "" {
		userID = "system"
	}
	deviceID := "web"

	// Use transaction-wrapped method
	reason := cashier.VarianceReason(req.Reason)
	shift, err := h.cashierShiftService.DocumentVarianceWithTransaction(
		c.Request.Context(),
		shiftObjID,
		reason,
		req.Notes,
		userID,
		deviceID,
	)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, shift)
}

// ConfirmResponsibility confirms cashier's responsibility
func (h *CashierShiftClosureHandler) ConfirmResponsibility(c *gin.Context) {
	shiftID := c.Param("id")
	
	shiftObjID, err := primitive.ObjectIDFromHex(shiftID)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid shift ID"})
		return
	}

	// Get user_id from JWT token
	userID := c.GetString("user_id")
	if userID == "" {
		userID = "system"
	}
	deviceID := "web"

	// Get the shift
	shift, err := h.cashierShiftService.GetCashierShift(c.Request.Context(), shiftObjID)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "shift not found"})
		return
	}

	// Validate step using state machine
	err = h.stateMachineManager.ValidateCashierShiftStep(shift, "confirm_responsibility")
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error":     err.Error(),
			"next_step": h.stateMachineManager.GetCashierShiftNextStep(shift),
		})
		return
	}

	// Confirm responsibility
	err = shift.ConfirmResponsibility(userID, deviceID, time.Now())
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Save the shift
	err = h.cashierShiftService.SaveCashierShift(c.Request.Context(), shift)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to save shift"})
		return
	}

	c.JSON(http.StatusOK, shift)
}

// CloseShift finalizes the shift closure
func (h *CashierShiftClosureHandler) CloseShift(c *gin.Context) {
	shiftID := c.Param("id")
	
	shiftObjID, err := primitive.ObjectIDFromHex(shiftID)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid shift ID"})
		return
	}

	// Get user_id from JWT token
	userID := c.GetString("user_id")
	if userID == "" {
		userID = "system"
	}
	deviceID := "web"

	// Use transaction-wrapped method
	shift, err := h.cashierShiftService.CloseShiftWithTransaction(
		c.Request.Context(),
		shiftObjID,
		userID,
		deviceID,
	)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, shift)
}


// CancelClosure cancels the shift closure process and returns to OPEN status
func (h *CashierShiftClosureHandler) CancelClosure(c *gin.Context) {
	shiftID := c.Param("id")
	
	shiftObjID, err := primitive.ObjectIDFromHex(shiftID)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid shift ID"})
		return
	}

	// Get user_id from JWT token
	userID := c.GetString("user_id")
	if userID == "" {
		userID = "system"
	}
	deviceID := "web"

	// Use transaction-wrapped method
	shift, err := h.cashierShiftService.CancelClosureWithTransaction(
		c.Request.Context(),
		shiftObjID,
		userID,
		deviceID,
	)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, shift)
}

// CheckWaiterShifts checks if there are any open waiter shifts
func (h *CashierShiftClosureHandler) CheckWaiterShifts(c *gin.Context) {
	allClosed, openShifts, err := h.cashierShiftService.CheckWaiterShiftsClosed(c.Request.Context())
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	// Return detailed information about open shifts
	var openShiftDetails []map[string]interface{}
	for _, shift := range openShifts {
		openShiftDetails = append(openShiftDetails, map[string]interface{}{
			"id":         shift.ID.Hex(),
			"user_name":  shift.UserName,
			"role_type":  shift.RoleType,
			"started_at": shift.StartedAt,
		})
	}

	c.JSON(http.StatusOK, gin.H{
		"all_closed":   allClosed,
		"open_shifts":  openShiftDetails,
		"open_count":   len(openShifts),
		"can_close":    allClosed,
	})
}


// CompleteClosure completes the entire closure workflow in one transaction
// This is a frontend-driven approach where all data is collected on the client
func (h *CashierShiftClosureHandler) CompleteClosure(c *gin.Context) {
	shiftID := c.Param("id")
	
	shiftObjID, err := primitive.ObjectIDFromHex(shiftID)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid shift ID"})
		return
	}

	var req CompleteClosureRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Get user_id from JWT token
	userID := c.GetString("user_id")
	if userID == "" {
		userID = "system"
	}
	deviceID := "web"

	// Convert variance reason if provided
	var varianceReason *cashier.VarianceReason
	if req.VarianceReason != nil {
		reason := cashier.VarianceReason(*req.VarianceReason)
		varianceReason = &reason
	}

	// Use complete closure method
	shift, err := h.cashierShiftService.CompleteClosure(
		c.Request.Context(),
		shiftObjID,
		req.ActualCash,
		varianceReason,
		req.VarianceNotes,
		userID,
		deviceID,
	)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, shift)
}

// ============================================================================
// FUND HANDOVER - NEW ENDPOINTS
// ============================================================================

// GetManagedFunds returns the current managed funds for a cashier shift
// GET /api/v1/cashier-shifts/:id/managed-funds
func (h *CashierShiftClosureHandler) GetManagedFunds(c *gin.Context) {
	shiftID := c.Param("id")
	
	shiftObjID, err := primitive.ObjectIDFromHex(shiftID)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid shift ID"})
		return
	}

	// Get managed funds summary
	summary, err := h.cashierShiftService.GetManagedFunds(c.Request.Context(), shiftObjID)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, summary)
}

// CloseShiftWithFundHandoverRequest represents the request to close shift with fund handover
type CloseShiftWithFundHandoverRequest struct {
	ActualCash     float64  `json:"actual_cash" binding:"required,min=0"`
	VarianceReason *string  `json:"variance_reason"`
	VarianceNotes  *string  `json:"variance_notes"`
	ReceiverID     *string  `json:"receiver_id"` // Optional, for future use
	UserID         string   `json:"user_id"`
	DeviceID       string   `json:"device_id"`
}

// CloseShiftWithFundHandover closes a cashier shift and creates a fund handover record
// POST /api/v1/cashier-shifts/:id/close-with-fund-handover
func (h *CashierShiftClosureHandler) CloseShiftWithFundHandover(c *gin.Context) {
	shiftID := c.Param("id")
	
	shiftObjID, err := primitive.ObjectIDFromHex(shiftID)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid shift ID"})
		return
	}

	var req CloseShiftWithFundHandoverRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Get user_id from JWT token
	userID := c.GetString("user_id")
	if userID == "" {
		userID = "system"
	}
	deviceID := "web"

	// Convert variance reason if provided
	var varianceReason *cashier.VarianceReason
	if req.VarianceReason != nil {
		reason := cashier.VarianceReason(*req.VarianceReason)
		varianceReason = &reason
	}

	// Convert receiver ID if provided
	var receiverObjID *primitive.ObjectID
	if req.ReceiverID != nil && *req.ReceiverID != "" {
		receiverID, err := primitive.ObjectIDFromHex(*req.ReceiverID)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "invalid receiver ID"})
			return
		}
		receiverObjID = &receiverID
	}

	// Close shift with fund handover
	shift, fundHandover, err := h.cashierShiftService.CloseShiftWithFundHandover(
		c.Request.Context(),
		shiftObjID,
		req.ActualCash,
		varianceReason,
		req.VarianceNotes,
		receiverObjID,
		userID,
		deviceID,
	)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"shift":         shift,
		"fund_handover": fundHandover,
	})
}
