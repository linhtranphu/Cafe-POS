package http

import (
	"net/http"
	"cafe-pos/backend/application/services"
	"cafe-pos/backend/domain"
	"cafe-pos/backend/domain/order"
	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type OrderHandler struct {
	orderService        *services.OrderService
	stateMachineManager *domain.StateMachineManager
}

func NewOrderHandler(
	orderService *services.OrderService,
	stateMachineManager *domain.StateMachineManager,
) *OrderHandler {
	return &OrderHandler{
		orderService:        orderService,
		stateMachineManager: stateMachineManager,
	}
}

func (h *OrderHandler) CreateOrder(c *gin.Context) {
	var req order.CreateOrderRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	userID, _ := c.Get("user_id")
	username, _ := c.Get("username")
	
	o, err := h.orderService.CreateOrder(c.Request.Context(), &req, userID.(string), username.(string))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, o)
}

func (h *OrderHandler) CollectPayment(c *gin.Context) {
	idStr := c.Param("id")
	id, err := primitive.ObjectIDFromHex(idStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid id"})
		return
	}

	var req order.PaymentRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	userID, _ := c.Get("user_id")
	username, _ := c.Get("username")
	req.CollectorID = userID.(string)
	req.CollectorName = username.(string)

	// Service layer handles validation
	o, err := h.orderService.CollectPayment(c.Request.Context(), id, &req)
	if err != nil {
		// Get order for error context
		ord, _ := h.orderService.GetOrder(c.Request.Context(), id)
		if ord != nil {
			c.JSON(http.StatusBadRequest, gin.H{
				"error":       err.Error(),
				"next_action": h.stateMachineManager.GetOrderNextAction(ord),
				"can_cancel":  h.stateMachineManager.CanCancelOrder(ord),
			})
		} else {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		}
		return
	}

	c.JSON(http.StatusOK, o)
}

func (h *OrderHandler) EditOrder(c *gin.Context) {
	idStr := c.Param("id")
	id, err := primitive.ObjectIDFromHex(idStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid id"})
		return
	}

	var req order.EditOrderRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Get the order first to validate state
	o, err := h.orderService.GetOrder(c.Request.Context(), id)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "order not found"})
		return
	}

	// Check if order can be modified using state machine
	if !h.stateMachineManager.CanModifyOrder(o) {
		c.JSON(http.StatusBadRequest, gin.H{
			"error":       "cannot modify order in current state",
			"status":      o.Status,
			"next_action": h.stateMachineManager.GetOrderNextAction(o),
		})
		return
	}

	response, err := h.orderService.EditOrder(c.Request.Context(), id, &req)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, response)
}

func (h *OrderHandler) RefundPartial(c *gin.Context) {
	idStr := c.Param("id")
	id, err := primitive.ObjectIDFromHex(idStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid id"})
		return
	}

	var req order.RefundRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Get the order first to validate state
	o, err := h.orderService.GetOrder(c.Request.Context(), id)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "order not found"})
		return
	}

	// Validate state transition using state machine
	err = h.stateMachineManager.ValidateOrderTransition(o, order.EventRefundOrder)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error":       err.Error(),
			"status":      o.Status,
			"can_refund":  h.stateMachineManager.CanRefundOrder(o),
			"next_action": h.stateMachineManager.GetOrderNextAction(o),
		})
		return
	}

	o, err = h.orderService.RefundPartial(c.Request.Context(), id, &req)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, o)
}

func (h *OrderHandler) SendToBar(c *gin.Context) {
	idStr := c.Param("id")
	id, err := primitive.ObjectIDFromHex(idStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid id"})
		return
	}

	// Service layer handles validation
	o, err := h.orderService.SendToBar(c.Request.Context(), id)
	if err != nil {
		// Get order for error context
		ord, _ := h.orderService.GetOrder(c.Request.Context(), id)
		if ord != nil {
			c.JSON(http.StatusBadRequest, gin.H{
				"error":       err.Error(),
				"status":      ord.Status,
				"next_action": h.stateMachineManager.GetOrderNextAction(ord),
			})
		} else {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		}
		return
	}

	c.JSON(http.StatusOK, o)
}

// AcceptOrder - Barista accepts order from queue
func (h *OrderHandler) AcceptOrder(c *gin.Context) {
	idStr := c.Param("id")
	id, err := primitive.ObjectIDFromHex(idStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid id"})
		return
	}

	userID, _ := c.Get("user_id")
	username, _ := c.Get("username")

	// Get the order first to validate state
	o, err := h.orderService.GetOrder(c.Request.Context(), id)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "order not found"})
		return
	}

	// Validate state transition using state machine
	err = h.stateMachineManager.ValidateOrderTransition(o, order.EventStartPreparing)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error":       err.Error(),
			"status":      o.Status,
			"next_action": h.stateMachineManager.GetOrderNextAction(o),
		})
		return
	}

	o, err = h.orderService.AcceptOrder(c.Request.Context(), id, userID.(string), username.(string))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, o)
}

// FinishPreparing - Barista marks order as ready
func (h *OrderHandler) FinishPreparing(c *gin.Context) {
	idStr := c.Param("id")
	id, err := primitive.ObjectIDFromHex(idStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid id"})
		return
	}

	// Get the order first to validate state
	o, err := h.orderService.GetOrder(c.Request.Context(), id)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "order not found"})
		return
	}

	// Validate state transition using state machine
	err = h.stateMachineManager.ValidateOrderTransition(o, order.EventMarkReady)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error":       err.Error(),
			"status":      o.Status,
			"progress":    h.stateMachineManager.GetOrderProgress(o),
			"next_action": h.stateMachineManager.GetOrderNextAction(o),
		})
		return
	}

	o, err = h.orderService.FinishPreparing(c.Request.Context(), id)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, o)
}

func (h *OrderHandler) ServeOrder(c *gin.Context) {
	idStr := c.Param("id")
	id, err := primitive.ObjectIDFromHex(idStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid id"})
		return
	}

	// Get the order first to validate state
	o, err := h.orderService.GetOrder(c.Request.Context(), id)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "order not found"})
		return
	}

	// Validate state transition using state machine
	err = h.stateMachineManager.ValidateOrderTransition(o, order.EventServeOrder)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error":       err.Error(),
			"status":      o.Status,
			"next_action": h.stateMachineManager.GetOrderNextAction(o),
		})
		return
	}

	o, err = h.orderService.ServeOrder(c.Request.Context(), id)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, o)
}

func (h *OrderHandler) CancelOrder(c *gin.Context) {
	idStr := c.Param("id")
	id, err := primitive.ObjectIDFromHex(idStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid id"})
		return
	}

	var req order.CancelOrderRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Get the order first to validate state
	o, err := h.orderService.GetOrder(c.Request.Context(), id)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "order not found"})
		return
	}

	// Validate state transition using state machine
	err = h.stateMachineManager.ValidateOrderTransition(o, order.EventCancelOrder)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error":       err.Error(),
			"status":      o.Status,
			"can_cancel":  h.stateMachineManager.CanCancelOrder(o),
			"next_action": h.stateMachineManager.GetOrderNextAction(o),
		})
		return
	}

	o, err = h.orderService.CancelOrder(c.Request.Context(), id, &req)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, o)
}

func (h *OrderHandler) LockOrder(c *gin.Context) {
	idStr := c.Param("id")
	id, err := primitive.ObjectIDFromHex(idStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid id"})
		return
	}

	// Get the order first to validate state
	o, err := h.orderService.GetOrder(c.Request.Context(), id)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "order not found"})
		return
	}

	// Validate state transition using state machine
	err = h.stateMachineManager.ValidateOrderTransition(o, order.EventLockOrder)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error":       err.Error(),
			"status":      o.Status,
			"can_lock":    h.stateMachineManager.CanLockOrder(o),
			"next_action": h.stateMachineManager.GetOrderNextAction(o),
		})
		return
	}

	o, err = h.orderService.LockOrder(c.Request.Context(), id)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, o)
}

func (h *OrderHandler) GetMyOrders(c *gin.Context) {
	userID, _ := c.Get("user_id")
	waiterID, _ := primitive.ObjectIDFromHex(userID.(string))

	orders, err := h.orderService.GetOrdersByWaiter(c.Request.Context(), waiterID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, orders)
}

func (h *OrderHandler) GetAllOrders(c *gin.Context) {
	orders, err := h.orderService.GetAllOrders(c.Request.Context())
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, orders)
}

// GetQueuedOrders - Get orders waiting for barista
func (h *OrderHandler) GetQueuedOrders(c *gin.Context) {
	orders, err := h.orderService.GetQueuedOrders(c.Request.Context())
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, orders)
}

// GetMyBaristaOrders - Get orders assigned to current barista
func (h *OrderHandler) GetMyBaristaOrders(c *gin.Context) {
	userID, _ := c.Get("user_id")
	baristaID, _ := primitive.ObjectIDFromHex(userID.(string))

	orders, err := h.orderService.GetBaristaOrders(c.Request.Context(), baristaID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, orders)
}

func (h *OrderHandler) GetOrder(c *gin.Context) {
	idStr := c.Param("id")
	id, err := primitive.ObjectIDFromHex(idStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid id"})
		return
	}

	o, err := h.orderService.GetOrder(c.Request.Context(), id)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "order not found"})
		return
	}

	c.JSON(http.StatusOK, o)
}
