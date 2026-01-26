package http

import (
	"net/http"
	"time"
	"cafe-pos/backend/application/services"
	"cafe-pos/backend/domain/expense"
	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type ExpenseHandler struct {
	service *services.ExpenseService
}

func NewExpenseHandler(service *services.ExpenseService) *ExpenseHandler {
	return &ExpenseHandler{service: service}
}

func (h *ExpenseHandler) CreateExpense(c *gin.Context) {
	var e expense.Expense
	if err := c.ShouldBindJSON(&e); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	if err := h.service.CreateExpense(c.Request.Context(), &e); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusCreated, e)
}

func (h *ExpenseHandler) GetExpenses(c *gin.Context) {
	filter := bson.M{}
	if startDate := c.Query("start_date"); startDate != "" {
		if t, err := time.Parse("2006-01-02", startDate); err == nil {
			filter["date"] = bson.M{"$gte": t}
		}
	}
	if endDate := c.Query("end_date"); endDate != "" {
		if t, err := time.Parse("2006-01-02", endDate); err == nil {
			if filter["date"] != nil {
				filter["date"].(bson.M)["$lte"] = t
			} else {
				filter["date"] = bson.M{"$lte": t}
			}
		}
	}
	if categoryID := c.Query("category_id"); categoryID != "" {
		if id, err := primitive.ObjectIDFromHex(categoryID); err == nil {
			filter["category_id"] = id
		}
	}
	expenses, err := h.service.GetExpenses(c.Request.Context(), filter)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, expenses)
}

func (h *ExpenseHandler) UpdateExpense(c *gin.Context) {
	id, _ := primitive.ObjectIDFromHex(c.Param("id"))
	var e expense.Expense
	if err := c.ShouldBindJSON(&e); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	if err := h.service.UpdateExpense(c.Request.Context(), id, &e); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, e)
}

func (h *ExpenseHandler) DeleteExpense(c *gin.Context) {
	id, _ := primitive.ObjectIDFromHex(c.Param("id"))
	if err := h.service.DeleteExpense(c.Request.Context(), id); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.Status(http.StatusNoContent)
}

func (h *ExpenseHandler) CreateCategory(c *gin.Context) {
	var cat expense.Category
	if err := c.ShouldBindJSON(&cat); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	if err := h.service.CreateCategory(c.Request.Context(), &cat); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusCreated, cat)
}

func (h *ExpenseHandler) GetCategories(c *gin.Context) {
	categories, err := h.service.GetCategories(c.Request.Context())
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, categories)
}

func (h *ExpenseHandler) DeleteCategory(c *gin.Context) {
	id, _ := primitive.ObjectIDFromHex(c.Param("id"))
	if err := h.service.DeleteCategory(c.Request.Context(), id); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.Status(http.StatusNoContent)
}

func (h *ExpenseHandler) CreateRecurring(c *gin.Context) {
	var re expense.RecurringExpense
	if err := c.ShouldBindJSON(&re); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	if err := h.service.CreateRecurring(c.Request.Context(), &re); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusCreated, re)
}

func (h *ExpenseHandler) GetRecurring(c *gin.Context) {
	recurring, err := h.service.GetRecurring(c.Request.Context())
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, recurring)
}

func (h *ExpenseHandler) DeleteRecurring(c *gin.Context) {
	id, _ := primitive.ObjectIDFromHex(c.Param("id"))
	if err := h.service.DeleteRecurring(c.Request.Context(), id); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.Status(http.StatusNoContent)
}

func (h *ExpenseHandler) CreatePrepaid(c *gin.Context) {
	var pe expense.PrepaidExpense
	if err := c.ShouldBindJSON(&pe); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	if err := h.service.CreatePrepaid(c.Request.Context(), &pe); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusCreated, pe)
}

func (h *ExpenseHandler) GetPrepaid(c *gin.Context) {
	prepaid, err := h.service.GetPrepaid(c.Request.Context())
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, prepaid)
}

func (h *ExpenseHandler) DeletePrepaid(c *gin.Context) {
	id, _ := primitive.ObjectIDFromHex(c.Param("id"))
	if err := h.service.DeletePrepaid(c.Request.Context(), id); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.Status(http.StatusNoContent)
}
