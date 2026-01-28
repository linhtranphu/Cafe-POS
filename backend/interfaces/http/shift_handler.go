package http

import (
	"net/http"
	"cafe-pos/backend/application/services"
	"cafe-pos/backend/domain/order"
	"cafe-pos/backend/domain/user"
	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type ShiftHandler struct {
	shiftService *services.ShiftService
}

func NewShiftHandler(shiftService *services.ShiftService) *ShiftHandler {
	return &ShiftHandler{shiftService: shiftService}
}

func (h *ShiftHandler) StartShift(c *gin.Context) {
	var req order.StartShiftRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	userID, _ := c.Get("user_id")
	username, _ := c.Get("username")
	role, _ := c.Get("role")
	
	// Convert user.Role to order.RoleType
	roleType := order.ParseRoleType(string(role.(user.Role)))

	shift, err := h.shiftService.StartShift(c.Request.Context(), &req, userID.(string), username.(string), roleType)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, shift)
}

func (h *ShiftHandler) EndShift(c *gin.Context) {
	idStr := c.Param("id")
	id, err := primitive.ObjectIDFromHex(idStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid id"})
		return
	}

	var req order.EndShiftRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	shift, err := h.shiftService.EndShift(c.Request.Context(), id, &req)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, shift)
}

func (h *ShiftHandler) CloseShift(c *gin.Context) {
	idStr := c.Param("id")
	id, err := primitive.ObjectIDFromHex(idStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid id"})
		return
	}

	var req order.EndShiftRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	shift, err := h.shiftService.CloseShiftAndLockOrders(c.Request.Context(), id, &req)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, shift)
}

func (h *ShiftHandler) GetCurrentShift(c *gin.Context) {
	userID, _ := c.Get("user_id")
	role, _ := c.Get("role")
	userOID, _ := primitive.ObjectIDFromHex(userID.(string))
	
	// Convert user.Role to order.RoleType
	roleType := order.ParseRoleType(string(role.(user.Role)))

	shift, err := h.shiftService.GetCurrentShift(c.Request.Context(), userOID, roleType)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "no open shift found"})
		return
	}

	c.JSON(http.StatusOK, shift)
}

func (h *ShiftHandler) GetMyShifts(c *gin.Context) {
	userID, _ := c.Get("user_id")
	role, _ := c.Get("role")
	userOID, _ := primitive.ObjectIDFromHex(userID.(string))
	
	// Convert user.Role to order.RoleType
	roleType := order.ParseRoleType(string(role.(user.Role)))

	shifts, err := h.shiftService.GetShiftsByUser(c.Request.Context(), userOID, roleType)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, shifts)
}

func (h *ShiftHandler) GetAllShifts(c *gin.Context) {
	shifts, err := h.shiftService.GetAllShifts(c.Request.Context())
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, shifts)
}

func (h *ShiftHandler) GetShift(c *gin.Context) {
	idStr := c.Param("id")
	id, err := primitive.ObjectIDFromHex(idStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid id"})
		return
	}

	shift, err := h.shiftService.GetShift(c.Request.Context(), id)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "shift not found"})
		return
	}

	c.JSON(http.StatusOK, shift)
}
