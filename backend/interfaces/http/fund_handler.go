package http

import (
	"fmt"
	"net/http"
	"strconv"
	"time"

	"cafe-pos/backend/application/services"
	"cafe-pos/backend/domain/fund"
	"cafe-pos/backend/domain/user"
	"cafe-pos/backend/infrastructure/mongodb"

	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type FundHandler struct {
	journalService *services.JournalService
}

func NewFundHandler(journalService *services.JournalService) *FundHandler {
	return &FundHandler{journalService: journalService}
}

// getUserObjectID extracts user ObjectID from Gin context
func getUserObjectID(c *gin.Context) (primitive.ObjectID, error) {
	userID, ok := c.Get("user_id")
	if !ok {
		return primitive.NilObjectID, fmt.Errorf("user_id not found in context")
	}
	switch v := userID.(type) {
	case primitive.ObjectID:
		return v, nil
	case string:
		return primitive.ObjectIDFromHex(v)
	default:
		return primitive.NilObjectID, fmt.Errorf("invalid user_id type: %T", v)
	}
}

// getUserInfo extracts user id, name and role from Gin context
func getUserInfo(c *gin.Context) (userObjID primitive.ObjectID, userName, userRole string, err error) {
	userObjID, err = getUserObjectID(c)
	if err != nil {
		return primitive.NilObjectID, "", "", err
	}

	if name, ok := c.Get("username"); ok {
		userName, _ = name.(string)
	}

	if role, ok := c.Get("role"); ok {
		if roleType, ok := role.(user.Role); ok {
			userRole = string(roleType)
		} else if roleStr, ok := role.(string); ok {
			userRole = roleStr
		}
	}

	if userRole == "" {
		userRole = "manager"
	}

	return userObjID, userName, userRole, nil
}

// ─── Request types ────────────────────────────────────────────────────────────

type DepositRequest struct {
	CashAmount     float64 `json:"cash_amount"`
	TransferAmount float64 `json:"transfer_amount"`
	Reason         string  `json:"reason"`
	FundType       string  `json:"fund_type"` // optional, default "operating"
}

type WithdrawRequest struct {
	CashAmount     float64 `json:"cash_amount"`
	TransferAmount float64 `json:"transfer_amount"`
	Reason         string  `json:"reason"`
	FundType       string  `json:"fund_type"` // optional, default "operating"
}

type FundTransferRequest struct {
	FromFundType   string  `json:"from_fund_type"`
	ToFundType     string  `json:"to_fund_type"`
	CashAmount     float64 `json:"cash_amount"`
	TransferAmount float64 `json:"transfer_amount"`
	Reason         string  `json:"reason"`
}

// ─── Handlers ─────────────────────────────────────────────────────────────────

// Deposit handles POST /api/manager/fund/deposit
func (h *FundHandler) Deposit(c *gin.Context) {
	ctx := c.Request.Context()

	userObjID, userName, userRole, err := getUserInfo(c)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "unauthorized: " + err.Error()})
		return
	}

	var req DepositRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid request body"})
		return
	}

	ft := fund.FundType(req.FundType)
	if ft == "" {
		ft = fund.FundTypeOperating
	}

	entry, err := h.journalService.RecordManagerDeposit(
		ctx,
		ft,
		req.CashAmount, req.TransferAmount,
		req.Reason,
		userObjID, userName, userRole,
	)
	if err != nil {
		fmt.Printf("ERROR in Deposit handler: %v\n", err)
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, entry)
}

// Withdraw handles POST /api/manager/fund/withdraw
func (h *FundHandler) Withdraw(c *gin.Context) {
	ctx := c.Request.Context()

	userObjID, userName, userRole, err := getUserInfo(c)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "unauthorized: " + err.Error()})
		return
	}

	var req WithdrawRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid request body"})
		return
	}

	ft := fund.FundType(req.FundType)
	if ft == "" {
		ft = fund.FundTypeOperating
	}

	entry, err := h.journalService.RecordManagerWithdrawal(
		ctx,
		ft,
		req.CashAmount, req.TransferAmount,
		req.Reason,
		userObjID, userName, userRole,
	)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, entry)
}

// Transfer handles POST /api/manager/fund/transfer
func (h *FundHandler) Transfer(c *gin.Context) {
	ctx := c.Request.Context()

	userObjID, userName, userRole, err := getUserInfo(c)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "unauthorized: " + err.Error()})
		return
	}

	var req FundTransferRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid request body"})
		return
	}

	if req.FromFundType == "" || req.ToFundType == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "from_fund_type and to_fund_type are required"})
		return
	}

	entry, err := h.journalService.RecordFundTransfer(
		ctx,
		fund.FundType(req.FromFundType),
		fund.FundType(req.ToFundType),
		req.CashAmount,
		req.TransferAmount,
		req.Reason,
		userObjID, userName, userRole,
	)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, entry)
}

// ─── Journal / Double-Entry Endpoints ─────────────────────────────────────────

// GetJournalFundBalances handles GET /api/manager/fund/journal-balances
func (h *FundHandler) GetJournalFundBalances(c *gin.Context) {
	ctx := c.Request.Context()
	balances, err := h.journalService.GetAllFundBalances(ctx)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, balances)
}

// GetJournalEntries handles GET /api/manager/fund/journal
// Query params: fund_type, event_type, from_date, to_date, limit, offset
func (h *FundHandler) GetJournalEntries(c *gin.Context) {
	ctx := c.Request.Context()

	filter := mongodb.JournalFilter{
		Limit:  20,
		Offset: 0,
	}

	if v := c.Query("fund_type"); v != "" {
		ft := fund.FundType(v)
		filter.FundType = &ft
	}
	if v := c.Query("event_type"); v != "" {
		et := fund.JournalEventType(v)
		filter.EventType = &et
	}
	if v := c.Query("from_date"); v != "" {
		if t, err := time.Parse(time.RFC3339, v); err == nil {
			filter.FromDate = &t
		}
	}
	if v := c.Query("to_date"); v != "" {
		if t, err := time.Parse(time.RFC3339, v); err == nil {
			filter.ToDate = &t
		}
	}
	if v := c.Query("limit"); v != "" {
		if n, err := strconv.Atoi(v); err == nil && n > 0 {
			filter.Limit = n
		}
	}
	if v := c.Query("offset"); v != "" {
		if n, err := strconv.Atoi(v); err == nil && n >= 0 {
			filter.Offset = n
		}
	}

	entries, total, err := h.journalService.ListEntries(ctx, filter)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"entries": entries,
		"total":   total,
		"limit":   filter.Limit,
		"offset":  filter.Offset,
	})
}

// GetJournalEntry handles GET /api/manager/fund/journal/:id
func (h *FundHandler) GetJournalEntry(c *gin.Context) {
	ctx := c.Request.Context()

	id, err := primitive.ObjectIDFromHex(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid journal entry ID"})
		return
	}

	entry, err := h.journalService.GetEntry(ctx, id)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "journal entry not found"})
		return
	}

	c.JSON(http.StatusOK, entry)
}
