package http

import (
	"net/http"
	"time"

	"cafe-pos/backend/application/services"
	"cafe-pos/backend/domain/schedule"
	"cafe-pos/backend/domain/user"
	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type ScheduleHandler struct {
	svc *services.ScheduleService
}

func NewScheduleHandler(svc *services.ScheduleService) *ScheduleHandler {
	return &ScheduleHandler{svc: svc}
}

// ─── Template endpoints (manager) ────────────────────────────────────────────

func (h *ScheduleHandler) CreateTemplate(c *gin.Context) {
	var t schedule.ShiftTemplate
	if err := c.ShouldBindJSON(&t); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	if err := h.svc.CreateTemplate(c.Request.Context(), &t); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusCreated, t)
}

func (h *ScheduleHandler) GetTemplates(c *gin.Context) {
	list, err := h.svc.GetTemplates(c.Request.Context())
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, list)
}

func (h *ScheduleHandler) UpdateTemplate(c *gin.Context) {
	id, err := primitive.ObjectIDFromHex(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid id"})
		return
	}
	var t schedule.ShiftTemplate
	if err := c.ShouldBindJSON(&t); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	if err := h.svc.UpdateTemplate(c.Request.Context(), id, &t); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, t)
}

func (h *ScheduleHandler) DeleteTemplate(c *gin.Context) {
	id, err := primitive.ObjectIDFromHex(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid id"})
		return
	}
	if err := h.svc.DeleteTemplate(c.Request.Context(), id); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "ok"})
}

// ─── Period endpoints (manager) ───────────────────────────────────────────────

func (h *ScheduleHandler) CreatePeriod(c *gin.Context) {
	var p schedule.SchedulePeriod
	if err := c.ShouldBindJSON(&p); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	userIDStr, _ := c.Get("user_id")
	userName, _ := c.Get("username")
	if oid, err := primitive.ObjectIDFromHex(userIDStr.(string)); err == nil {
		p.CreatedBy = oid
	}
	p.CreatedByName, _ = userName.(string)
	if err := h.svc.CreatePeriod(c.Request.Context(), &p); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusCreated, p)
}

func (h *ScheduleHandler) GetPeriods(c *gin.Context) {
	list, err := h.svc.GetPeriods(c.Request.Context())
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, list)
}

func (h *ScheduleHandler) SetPeriodStatus(c *gin.Context) {
	id, err := primitive.ObjectIDFromHex(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid id"})
		return
	}
	var body struct {
		Status schedule.PeriodStatus `json:"status"`
	}
	if err := c.ShouldBindJSON(&body); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	if err := h.svc.SetPeriodStatus(c.Request.Context(), id, body.Status); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "ok"})
}

func (h *ScheduleHandler) GetPeriodDetail(c *gin.Context) {
	id, err := primitive.ObjectIDFromHex(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid id"})
		return
	}
	// Optionally pass current user so MyRegistration is populated
	var userID *primitive.ObjectID
	if userIDStr, ok := c.Get("user_id"); ok {
		if oid, err := primitive.ObjectIDFromHex(userIDStr.(string)); err == nil {
			userID = &oid
		}
	}
	detail, err := h.svc.GetPeriodDetail(c.Request.Context(), id, userID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, detail)
}

// ─── Slot endpoints (manager) ─────────────────────────────────────────────────

func (h *ScheduleHandler) AddSlot(c *gin.Context) {
	periodID, err := primitive.ObjectIDFromHex(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid period id"})
		return
	}
	var body struct {
		TemplateID string    `json:"template_id"`
		Date       time.Time `json:"date"`
	}
	if err := c.ShouldBindJSON(&body); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	templateID, err := primitive.ObjectIDFromHex(body.TemplateID)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid template_id"})
		return
	}
	slot, err := h.svc.AddSlot(c.Request.Context(), periodID, templateID, body.Date)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusCreated, slot)
}

func (h *ScheduleHandler) RemoveSlot(c *gin.Context) {
	slotID, err := primitive.ObjectIDFromHex(c.Param("slotId"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid slot id"})
		return
	}
	if err := h.svc.RemoveSlot(c.Request.Context(), slotID); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "ok"})
}

// CancelRegistration allows a manager to cancel any staff registration
func (h *ScheduleHandler) CancelRegistration(c *gin.Context) {
	regID, err := primitive.ObjectIDFromHex(c.Param("regId"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid registration id"})
		return
	}
	if err := h.svc.CancelRegistration(c.Request.Context(), regID); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "ok"})
}

// ─── Staff endpoints ──────────────────────────────────────────────────────────

// GetCurrentPeriod returns the active open/published period with my registrations
func (h *ScheduleHandler) GetCurrentPeriod(c *gin.Context) {
	var userID *primitive.ObjectID
	if userIDStr, ok := c.Get("user_id"); ok {
		if oid, err := primitive.ObjectIDFromHex(userIDStr.(string)); err == nil {
			userID = &oid
		}
	}
	p, err := h.svc.GetCurrentPeriod(c.Request.Context())
	if err != nil {
		c.JSON(http.StatusOK, nil)
		return
	}
	detail, err := h.svc.GetPeriodDetail(c.Request.Context(), p.ID, userID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, detail)
}

func (h *ScheduleHandler) Register(c *gin.Context) {
	slotID, err := primitive.ObjectIDFromHex(c.Param("slotId"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid slot id"})
		return
	}
	userIDStr, _ := c.Get("user_id")
	userName, _ := c.Get("username")
	roleVal, _ := c.Get("role")
	oid, err := primitive.ObjectIDFromHex(userIDStr.(string))
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "invalid user"})
		return
	}
	reg, err := h.svc.Register(
		c.Request.Context(),
		slotID,
		oid,
		userName.(string),
		string(roleVal.(user.Role)),
	)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusCreated, reg)
}

func (h *ScheduleHandler) Unregister(c *gin.Context) {
	slotID, err := primitive.ObjectIDFromHex(c.Param("slotId"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid slot id"})
		return
	}
	userIDStr, _ := c.Get("user_id")
	oid, err := primitive.ObjectIDFromHex(userIDStr.(string))
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "invalid user"})
		return
	}
	if err := h.svc.Unregister(c.Request.Context(), slotID, oid); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "ok"})
}

func (h *ScheduleHandler) GetMySchedule(c *gin.Context) {
	userIDStr, _ := c.Get("user_id")
	oid, err := primitive.ObjectIDFromHex(userIDStr.(string))
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "invalid user"})
		return
	}
	detail, err := h.svc.GetMySchedule(c.Request.Context(), oid)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, detail)
}
