package http

import (
	"net/http"
	"time"

	"cafe-pos/backend/application/services"

	"github.com/gin-gonic/gin"
)

type CashierHandler struct {
	reconciliationService *services.CashReconciliationService
	oversightService      *services.PaymentOversightService
	reportService         *services.CashierReportService
}

func NewCashierHandler(
	reconciliationService *services.CashReconciliationService,
	oversightService *services.PaymentOversightService,
	reportService *services.CashierReportService,
) *CashierHandler {
	return &CashierHandler{
		reconciliationService: reconciliationService,
		oversightService:      oversightService,
		reportService:         reportService,
	}
}

// FR-CASH-02: Theo dõi trạng thái ca
func (h *CashierHandler) GetShiftStatus(c *gin.Context) {
	shiftID := c.Param("id")

	status, err := h.reconciliationService.GetShiftStatus(shiftID)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, status)
}

// FR-CASH-04: Giám sát thanh toán
func (h *CashierHandler) GetPaymentsByShift(c *gin.Context) {
	shiftID := c.Param("id")

	payments, err := h.oversightService.GetPaymentsByShift(shiftID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"payments": payments})
}

// FR-CASH-05: Xử lý sai lệch thanh toán
func (h *CashierHandler) ReportDiscrepancy(c *gin.Context) {
	var req struct {
		OrderID string  `json:"order_id" binding:"required"`
		Reason  string  `json:"reason" binding:"required"`
		Amount  float64 `json:"amount" binding:"required"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	cashierID := c.GetString("user_id")

	err := h.oversightService.ReportDiscrepancy(req.OrderID, req.Reason, req.Amount, cashierID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, gin.H{"message": "Discrepancy reported successfully"})
}

// FR-CASH-06: Đối soát tiền mặt
func (h *CashierHandler) ReconcileCash(c *gin.Context) {
	var req struct {
		ShiftID    string  `json:"shift_id" binding:"required"`
		ActualCash float64 `json:"actual_cash" binding:"required"`
		Notes      string  `json:"notes"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	cashierID := c.GetString("user_id")

	reconciliation, err := h.reconciliationService.ReconcileCash(req.ShiftID, req.ActualCash, req.Notes, cashierID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, reconciliation)
}

// FR-CASH-08: Hủy/điều chỉnh thanh toán
func (h *CashierHandler) OverridePayment(c *gin.Context) {
	orderID := c.Param("id")
	
	var req struct {
		Reason string `json:"reason" binding:"required"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	cashierID := c.GetString("user_id")

	err := h.oversightService.OverridePayment(orderID, req.Reason, cashierID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Payment overridden successfully"})
}

// FR-CASH-09: Khóa order
func (h *CashierHandler) LockOrder(c *gin.Context) {
	orderID := c.Param("id")
	cashierID := c.GetString("user_id")

	err := h.oversightService.LockOrder(orderID, cashierID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Order locked successfully"})
}

// FR-CASH-10: Báo cáo ca
func (h *CashierHandler) GenerateShiftReport(c *gin.Context) {
	shiftID := c.Param("id")

	report, err := h.reportService.GenerateShiftReport(shiftID)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, report)
}

// FR-CASH-11: Bàn giao ca
func (h *CashierHandler) HandoverShift(c *gin.Context) {
	var req services.HandoverData

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	err := h.reportService.HandoverShift(&req)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Shift handover completed successfully"})
}

// Get pending discrepancies
func (h *CashierHandler) GetPendingDiscrepancies(c *gin.Context) {
	discrepancies, err := h.oversightService.GetPendingDiscrepancies()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"discrepancies": discrepancies})
}

// Resolve discrepancy
func (h *CashierHandler) ResolveDiscrepancy(c *gin.Context) {
	discrepancyID := c.Param("id")

	err := h.oversightService.ResolveDiscrepancy(discrepancyID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Discrepancy resolved successfully"})
}

// Get daily report
func (h *CashierHandler) GetDailyReport(c *gin.Context) {
	dateStr := c.Query("date")
	if dateStr == "" {
		dateStr = time.Now().Format("2006-01-02")
	}

	date, err := time.Parse("2006-01-02", dateStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid date format. Use YYYY-MM-DD"})
		return
	}

	report, err := h.reportService.GetDailyReport(date)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, report)
}

// Get audit trail for order
func (h *CashierHandler) GetOrderAudits(c *gin.Context) {
	orderID := c.Param("id")

	audits, err := h.oversightService.GetAuditsByOrder(orderID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"audits": audits})
}