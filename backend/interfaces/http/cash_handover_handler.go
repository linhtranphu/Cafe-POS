package http

import (
	"net/http"

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

	var reqBody struct {
		CashAmount     float64               `json:"cash_amount"`
		TransferAmount float64               `json:"transfer_amount"`
		HandoverType   handover.HandoverType `json:"handover_type" binding:"required"`
		WaiterNote     string                `json:"waiter_note"`
		EndCash        float64               `json:"end_cash"`     // Tiền mặt cuối ca
		EndTransfer    float64               `json:"end_transfer"` // Tiền CK cuối ca
	}

	if err := c.ShouldBindJSON(&reqBody); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	userID, _ := c.Get("user_id")
	username, _ := c.Get("username")

	handoverResult, err := h.handoverService.CreateHandover(
		c.Request.Context(),
		shiftOID,
		reqBody.CashAmount,
		reqBody.TransferAmount,
		reqBody.HandoverType,
		reqBody.WaiterNote,
		reqBody.EndCash,
		reqBody.EndTransfer,
		userID.(string),
		username.(string),
	)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, handoverResult)
}

// ConfirmHandover confirms or rejects a handover (cashier)
func (h *CashHandoverHandler) ConfirmHandover(c *gin.Context) {
	handoverID := c.Param("id")
	handoverOID, err := primitive.ObjectIDFromHex(handoverID)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid handover ID"})
		return
	}

	var reqBody struct {
		ActualCashAmount          float64                      `json:"actual_cash_amount"`
		ActualTransferAmount      float64                      `json:"actual_transfer_amount"`
		Status                    handover.HandoverStatus      `json:"status" binding:"required"`
		CashierNote               string                       `json:"cashier_note"`
		DiscrepancyReason         string                       `json:"discrepancy_reason"`
		DiscrepancyResponsibility handover.ResponsibilityType  `json:"discrepancy_responsibility"`
	}

	if err := c.ShouldBindJSON(&reqBody); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	userID, _ := c.Get("user_id")

	err = h.handoverService.ConfirmHandover(
		c.Request.Context(),
		handoverOID,
		reqBody.ActualCashAmount,
		reqBody.ActualTransferAmount,
		reqBody.Status,
		reqBody.CashierNote,
		reqBody.DiscrepancyReason,
		reqBody.DiscrepancyResponsibility,
		userID.(string),
	)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "handover confirmed successfully"})
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
