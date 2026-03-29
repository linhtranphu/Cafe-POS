package http

import (
	"net/http"
	"time"

	"cafe-pos/backend/application/services"
	"cafe-pos/backend/domain/attendance"
	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type AttendanceHandler struct {
	svc *services.AttendanceService
}

func NewAttendanceHandler(svc *services.AttendanceService) *AttendanceHandler {
	return &AttendanceHandler{svc: svc}
}

func (h *AttendanceHandler) AddEntry(c *gin.Context) {
	var e attendance.HoursEntry
	if err := c.ShouldBindJSON(&e); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	userIDStr, _ := c.Get("user_id")
	userName, _ := c.Get("username")
	if oid, err := primitive.ObjectIDFromHex(userIDStr.(string)); err == nil {
		e.CreatedBy = oid
	}
	e.CreatedByName, _ = userName.(string)
	if err := h.svc.AddEntry(c.Request.Context(), &e); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusCreated, e)
}

func (h *AttendanceHandler) GetEntries(c *gin.Context) {
	fromStr := c.Query("from")
	toStr := c.Query("to")

	var from, to time.Time
	var err error

	if fromStr != "" {
		from, err = time.Parse("2006-01-02", fromStr)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "invalid from date"})
			return
		}
	} else {
		now := time.Now().UTC()
		day := now.Weekday()
		if day == time.Sunday {
			day = 7
		}
		from = now.AddDate(0, 0, -int(day-1))
		from = time.Date(from.Year(), from.Month(), from.Day(), 0, 0, 0, 0, time.UTC)
	}

	if toStr != "" {
		to, err = time.Parse("2006-01-02", toStr)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "invalid to date"})
			return
		}
		to = time.Date(to.Year(), to.Month(), to.Day(), 23, 59, 59, 999999999, time.UTC)
	} else {
		to = from.AddDate(0, 0, 6)
		to = time.Date(to.Year(), to.Month(), to.Day(), 23, 59, 59, 999999999, time.UTC)
	}

	entries, err := h.svc.GetEntries(c.Request.Context(), from, to)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	if entries == nil {
		entries = []*attendance.HoursEntry{}
	}
	c.JSON(http.StatusOK, entries)
}

func (h *AttendanceHandler) DeleteEntry(c *gin.Context) {
	id, err := primitive.ObjectIDFromHex(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid id"})
		return
	}
	if err := h.svc.DeleteEntry(c.Request.Context(), id); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "ok"})
}
