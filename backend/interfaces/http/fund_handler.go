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
	return &FundHandler{
		fundService: fundService,
	}
}

// Helper function to get user ObjectID from context
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

// Helper function to get user info from context
func getUserInfo(c *gin.Context) (userObjID primitive.ObjectID, userName, userRole string, err error) {
	userObjID, err = getUserObjectID(c)
	if err != nil {
		return primitive.NilObjectID, "", "", err
	}
	
	if name, ok := c.Get("username"); ok {
		userName, _ = name.(string)
	}
	
	if role, ok := c.Get("role"); ok {
		// Role is stored as user.Role type in context
		if roleType, ok := role.(user.Role); ok {
			userRole = string(roleType)
		} else if roleStr, ok := role.(string); ok {
			userRole = roleStr
		}
	}
	
	// Debug logging
	fmt.Printf("DEBUG getUserInfo - userObjID: %v, userName: %s, userRole: %s\n", userObjID, userName, userRole)
	
	// If role is empty, default to "manager" for backward compatibility with old tokens
	if userRole == "" {
		fmt.Printf("WARN: userRole is empty, defaulting to 'manager' for backward compatibility\n")
		userRole = "manager"
	}
	
	return userObjID, userName, userRole, nil
}

// BalanceResponse represents the fund balance response
type BalanceResponse struct {
	CurrentBalance *fund.FundBalance              `json:"current_balance"`
	TodaySummary   *services.TodayBalanceSummary  `json:"today_summary"`
	LastUpdated    time.Time                      `json:"last_updated"`
}

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

	response := BalanceResponse{
		CurrentBalance: currentBalance,
		TodaySummary:   todaySummary,
		LastUpdated:    time.Now(),
	}

	c.JSON(http.StatusOK, response)
}

// TransactionHistoryResponse represents the transaction history response
type TransactionHistoryResponse struct {
	Transactions []*services.TransactionHistoryItem `json:"transactions"`
	Total        int64                              `json:"total"`
	Limit        int                                `json:"limit"`
	Offset       int                                `json:"offset"`
}

// GetTransactions handles GET /api/manager/fund/transactions
func (h *FundHandler) GetTransactions(c *gin.Context) {
	ctx := c.Request.Context()

	// Parse query parameters
	// Transaction type filter
	var transactionType *string
	if typeStr := c.Query("type"); typeStr != "" {
		transactionType = &typeStr
	}

	// Money type filter
	var moneyType *string
	if mtStr := c.Query("money_type"); mtStr != "" {
		moneyType = &mtStr
	}

	// Date range
	fromDate := time.Now().Truncate(24 * time.Hour) // Default: today 00:00
	if fromStr := c.Query("from_date"); fromStr != "" {
		if parsed, err := time.Parse(time.RFC3339, fromStr); err == nil {
			fromDate = parsed
		}
	}

	toDate := time.Now() // Default: now
	if toStr := c.Query("to_date"); toStr != "" {
		if parsed, err := time.Parse(time.RFC3339, toStr); err == nil {
			toDate = parsed
		}
	}

	// Pagination
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

	// Fetch aggregated transactions from all sources
	transactions, total, err := h.fundService.GetAggregatedTransactionHistory(
		ctx,
		fromDate,
		toDate,
		transactionType,
		moneyType,
		limit,
		offset,
	)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	response := TransactionHistoryResponse{
		Transactions: transactions,
		Total:        total,
		Limit:        limit,
		Offset:       offset,
	}

	c.JSON(http.StatusOK, response)
}

// DepositRequest represents a deposit request
type DepositRequest struct {
	CashAmount     float64 `json:"cash_amount"`
	TransferAmount float64 `json:"transfer_amount"`
	Reason         string  `json:"reason"`
}

// TransactionResponse represents a transaction response
type TransactionResponse struct {
	Transaction *fund.FundTransaction `json:"transaction"`
	NewBalance  *fund.FundBalance     `json:"new_balance"`
}

// Deposit handles POST /api/manager/fund/deposit
func (h *FundHandler) Deposit(c *gin.Context) {
	ctx := c.Request.Context()

	// Get user info from context
	userObjID, userName, userRole, err := getUserInfo(c)
	if err != nil {
		fmt.Printf("ERROR getting user info: %v\n", err)
		c.JSON(http.StatusUnauthorized, gin.H{"error": "unauthorized: " + err.Error()})
		return
	}

	// Parse request
	var req DepositRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid request body"})
		return
	}

	// Create deposit
	transaction, newBalance, err := h.fundService.CreateDeposit(
		ctx,
		req.CashAmount,
		req.TransferAmount,
		req.Reason,
		userObjID,
		userName,
		userRole,
	)
	if err != nil {
		// Log the error for debugging
		fmt.Printf("ERROR in Deposit handler: %v\n", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	response := TransactionResponse{
		Transaction: transaction,
		NewBalance:  newBalance,
	}

	c.JSON(http.StatusCreated, response)
}

// WithdrawRequest represents a withdrawal request
type WithdrawRequest struct {
	CashAmount     float64 `json:"cash_amount"`
	TransferAmount float64 `json:"transfer_amount"`
	Reason         string  `json:"reason"`
}

// Withdraw handles POST /api/manager/fund/withdraw
func (h *FundHandler) Withdraw(c *gin.Context) {
	ctx := c.Request.Context()

	// Get user info from context
	userObjID, userName, userRole, err := getUserInfo(c)
	if err != nil {
		fmt.Printf("ERROR getting user info: %v\n", err)
		c.JSON(http.StatusUnauthorized, gin.H{"error": "unauthorized: " + err.Error()})
		return
	}

	// Parse request
	var req WithdrawRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid request body"})
		return
	}

	// Create withdrawal
	transaction, newBalance, err := h.fundService.CreateWithdrawal(
		ctx,
		req.CashAmount,
		req.TransferAmount,
		req.Reason,
		userObjID,
		userName,
		userRole,
	)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	response := TransactionResponse{
		Transaction: transaction,
		NewBalance:  newBalance,
	}

	c.JSON(http.StatusCreated, response)
}

// GetTransactionDetail handles GET /api/manager/fund/transactions/:id
func (h *FundHandler) GetTransactionDetail(c *gin.Context) {
	ctx := c.Request.Context()

	// Get ID from URL
	idStr := c.Param("id")
	
	id, err := primitive.ObjectIDFromHex(idStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid transaction ID"})
		return
	}

	// Fetch transaction
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

