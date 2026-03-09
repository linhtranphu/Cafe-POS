package http

import (
	"fmt"
	"net/http"
	"strconv"
	"time"

	"cafe-pos/backend/application/services"
	"cafe-pos/backend/domain/fund"
	"cafe-pos/backend/domain/user"

	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type FundHandler struct {
	fundService *services.FundService
}

func NewFundHandler(fundService *services.FundService) *FundHandler {
	return &FundHandler{fundService: fundService}
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

// ─── Response types ───────────────────────────────────────────────────────────

type BalanceResponse struct {
	CurrentBalance *fund.FundBalance             `json:"current_balance"`
	TodaySummary   *services.TodayBalanceSummary `json:"today_summary"`
	LastUpdated    time.Time                     `json:"last_updated"`
}

type TransactionHistoryResponse struct {
	Transactions []*services.TransactionHistoryItem `json:"transactions"`
	Total        int64                              `json:"total"`
	Limit        int                                `json:"limit"`
	Offset       int                                `json:"offset"`
}

type TransactionResponse struct {
	Transaction *fund.FundTransaction `json:"transaction"`
	NewBalance  *fund.FundBalance     `json:"new_balance"`
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

// GetBalance handles GET /api/manager/fund/balance
func (h *FundHandler) GetBalance(c *gin.Context) {
	ctx := c.Request.Context()

	currentBalance, err := h.fundService.CalculateCurrentBalance(ctx)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	todaySummary, err := h.fundService.CalculateTodayBalance(ctx)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, BalanceResponse{
		CurrentBalance: currentBalance,
		TodaySummary:   todaySummary,
		LastUpdated:    time.Now(),
	})
}

// GetAllFundBalances handles GET /api/manager/fund/balances
func (h *FundHandler) GetAllFundBalances(c *gin.Context) {
	ctx := c.Request.Context()

	balances, err := h.fundService.GetAllFundBalances(ctx)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"balances":     balances,
		"last_updated": time.Now(),
	})
}

// GetTransactions handles GET /api/manager/fund/transactions
func (h *FundHandler) GetTransactions(c *gin.Context) {
	ctx := c.Request.Context()

	var transactionType *string
	if typeStr := c.Query("type"); typeStr != "" && typeStr != "all" {
		transactionType = &typeStr
	}

	var moneyType *string
	if mtStr := c.Query("money_type"); mtStr != "" && mtStr != "all" {
		moneyType = &mtStr
	}

	var fundTypeFilter *string
	if ftStr := c.Query("fund_type"); ftStr != "" && ftStr != "all" {
		fundTypeFilter = &ftStr
	}

	var sourceTypeFilter *string
	if stStr := c.Query("source_type"); stStr != "" && stStr != "all" {
		sourceTypeFilter = &stStr
	}

	fromDate := time.Now().Truncate(24 * time.Hour)
	if fromStr := c.Query("from_date"); fromStr != "" {
		if parsed, err := time.Parse(time.RFC3339, fromStr); err == nil {
			fromDate = parsed
		}
	}

	toDate := time.Now()
	if toStr := c.Query("to_date"); toStr != "" {
		if parsed, err := time.Parse(time.RFC3339, toStr); err == nil {
			toDate = parsed
		}
	}

	limit := 50
	if limitStr := c.Query("limit"); limitStr != "" {
		if parsed, err := strconv.Atoi(limitStr); err == nil && parsed > 0 && parsed <= 200 {
			limit = parsed
		}
	}

	offset := 0
	if offsetStr := c.Query("offset"); offsetStr != "" {
		if parsed, err := strconv.Atoi(offsetStr); err == nil && parsed >= 0 {
			offset = parsed
		}
	}

	transactions, total, err := h.fundService.GetAggregatedTransactionHistory(
		ctx, fromDate, toDate,
		transactionType, moneyType,
		fundTypeFilter, sourceTypeFilter,
		limit, offset,
	)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, TransactionHistoryResponse{
		Transactions: transactions,
		Total:        total,
		Limit:        limit,
		Offset:       offset,
	})
}

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
	transaction, newBalance, err := h.fundService.CreateDeposit(
		ctx,
		req.CashAmount, req.TransferAmount,
		req.Reason,
		userObjID, userName, userRole,
		ft,
	)
	if err != nil {
		fmt.Printf("ERROR in Deposit handler: %v\n", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, TransactionResponse{
		Transaction: transaction,
		NewBalance:  newBalance,
	})
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
	transaction, newBalance, err := h.fundService.CreateWithdrawal(
		ctx,
		req.CashAmount, req.TransferAmount,
		req.Reason,
		userObjID, userName, userRole,
		ft,
	)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, TransactionResponse{
		Transaction: transaction,
		NewBalance:  newBalance,
	})
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

	result, err := h.fundService.TransferBetweenFunds(ctx, services.FundTransferRequest{
		FromFundType:    fund.FundType(req.FromFundType),
		ToFundType:      fund.FundType(req.ToFundType),
		CashAmount:      req.CashAmount,
		TransferAmount:  req.TransferAmount,
		Reason:          req.Reason,
		PerformedBy:     userObjID,
		PerformedByName: userName,
		PerformedByRole: userRole,
	})
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, result)
}

// GetTransactionDetail handles GET /api/manager/fund/transactions/:id
func (h *FundHandler) GetTransactionDetail(c *gin.Context) {
	ctx := c.Request.Context()

	id, err := primitive.ObjectIDFromHex(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid transaction ID"})
		return
	}

	transaction, err := h.fundService.GetTransactionDetail(ctx, id)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	if transaction == nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "transaction not found"})
		return
	}

	c.JSON(http.StatusOK, transaction)
}
