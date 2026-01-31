package http

import (
	"net/http"

	"cafe-pos/backend/application/services"
	"cafe-pos/backend/domain/user"

	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

// CashierShiftHandler handles HTTP requests for cashier shift management.
// This is separate from the regular shift handler to maintain clear separation
// between cashier shifts and waiter/barista shifts.
type CashierShiftHandler struct {
	cashierShiftService *services.CashierShiftService
}

// NewCashierShiftHandler creates a new handler for cashier shifts.
func NewCashierShiftHandler(cashierShiftService *services.CashierShiftService) *CashierShiftHandler {
	return &CashierShiftHandler{
		cashierShiftService: cashierShiftService,
	}
}

// StartCashierShiftRequest represents the request body for starting a cashier shift.
type StartCashierShiftRequest struct {
	StartingFloat float64 `json:"starting_float" binding:"required,min=0"`
}

// StartCashierShift handles POST /api/v1/cashier-shifts
// Creates a new cashier shift for the authenticated cashier.
func (h *CashierShiftHandler) StartCashierShift(c *gin.Context) {
	// Get authenticated user info
	userID, exists := c.Get("user_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "user not authenticated"})
		return
	}

	username, exists := c.Get("username")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "username not found"})
		return
	}

	role, exists := c.Get("role")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "user role not found"})
		return
	}

	// Validate user is cashier or manager
	userRole := role.(user.Role)
	if userRole != user.RoleCashier && userRole != user.RoleManager {
		c.JSON(http.StatusForbidden, gin.H{
			"error": "only cashiers and managers can start cashier shifts",
		})
		return
	}

	// Parse request body
	var req StartCashierShiftRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Convert user ID to ObjectID
	cashierOID, err := primitive.ObjectIDFromHex(userID.(string))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid user ID"})
		return
	}

	// Start cashier shift
	shift, err := h.cashierShiftService.StartCashierShift(
		c.Request.Context(),
		cashierOID,
		username.(string),
		req.StartingFloat,
	)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, shift)
}

// GetCurrentCashierShift handles GET /api/v1/cashier-shifts/current
// Retrieves the current open cashier shift for the authenticated cashier.
func (h *CashierShiftHandler) GetCurrentCashierShift(c *gin.Context) {
	// Get authenticated user info
	userID, exists := c.Get("user_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "user not authenticated"})
		return
	}

	// Convert user ID to ObjectID
	cashierOID, err := primitive.ObjectIDFromHex(userID.(string))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid user ID"})
		return
	}

	// Get current cashier shift
	shift, err := h.cashierShiftService.GetCurrentCashierShift(c.Request.Context(), cashierOID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	if shift == nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "no open cashier shift found"})
		return
	}

	c.JSON(http.StatusOK, shift)
}

// GetAllCashierShifts handles GET /api/v1/cashier-shifts
// Retrieves all cashier shifts (for managers/admins).
func (h *CashierShiftHandler) GetAllCashierShifts(c *gin.Context) {
	// Get user role
	role, exists := c.Get("role")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "user role not found"})
		return
	}

	// Only managers can view all cashier shifts
	userRole := role.(user.Role)
	if userRole != user.RoleManager {
		c.JSON(http.StatusForbidden, gin.H{
			"error": "only managers can view all cashier shifts",
		})
		return
	}

	// Get all cashier shifts
	shifts, err := h.cashierShiftService.GetAllCashierShifts(c.Request.Context())
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, shifts)
}

// GetCashierShift handles GET /api/v1/cashier-shifts/:id
// Retrieves a specific cashier shift by ID.
func (h *CashierShiftHandler) GetCashierShift(c *gin.Context) {
	// Parse shift ID from URL
	shiftIDStr := c.Param("id")
	shiftID, err := primitive.ObjectIDFromHex(shiftIDStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid shift ID"})
		return
	}

	// Get cashier shift
	shift, err := h.cashierShiftService.GetCashierShift(c.Request.Context(), shiftID)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "cashier shift not found"})
		return
	}

	// Check authorization - user must be the cashier or a manager
	userID, _ := c.Get("user_id")
	role, _ := c.Get("role")
	userRole := role.(user.Role)
	userOID, _ := primitive.ObjectIDFromHex(userID.(string))

	if userRole != user.RoleManager && shift.CashierID != userOID {
		c.JSON(http.StatusForbidden, gin.H{
			"error": "you can only view your own cashier shifts",
		})
		return
	}

	c.JSON(http.StatusOK, shift)
}

// GetMyCashierShifts handles GET /api/v1/cashier-shifts/my-shifts
// Retrieves all cashier shifts for the authenticated cashier.
func (h *CashierShiftHandler) GetMyCashierShifts(c *gin.Context) {
	// Get authenticated user info
	userID, exists := c.Get("user_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "user not authenticated"})
		return
	}

	// Convert user ID to ObjectID
	cashierOID, err := primitive.ObjectIDFromHex(userID.(string))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid user ID"})
		return
	}

	// Get cashier shifts for this user
	shifts, err := h.cashierShiftService.GetCashierShiftsByUser(c.Request.Context(), cashierOID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, shifts)
}
