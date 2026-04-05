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
	var body struct {
		schedule.SchedulePeriod
		WeeklyTemplateID string `json:"weekly_template_id"`
	}
	if err := c.ShouldBindJSON(&body); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	userIDStr, _ := c.Get("user_id")
	userName, _ := c.Get("username")
	if oid, err := primitive.ObjectIDFromHex(userIDStr.(string)); err == nil {
		body.CreatedBy = oid
	}
	body.CreatedByName, _ = userName.(string)
	var weeklyTemplateID *primitive.ObjectID
	if body.WeeklyTemplateID != "" {
		if oid, err := primitive.ObjectIDFromHex(body.WeeklyTemplateID); err == nil {
			weeklyTemplateID = &oid
		}
	}
	if err := h.svc.CreatePeriodWithWeeklyTemplate(c.Request.Context(), &body.SchedulePeriod, weeklyTemplateID); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusCreated, body.SchedulePeriod)
}

func (h *ScheduleHandler) GetPeriods(c *gin.Context) {
	list, err := h.svc.GetPeriods(c.Request.Context())
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, list)
}

func (h *ScheduleHandler) UpdatePeriodMaxHoursPercent(c *gin.Context) {
	id, err := primitive.ObjectIDFromHex(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid id"})
		return
	}
	var body struct {
		MaxHoursPercent float64 `json:"max_hours_percent"`
	}
	if err := c.ShouldBindJSON(&body); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	if err := h.svc.UpdatePeriodMaxHoursPercent(c.Request.Context(), id, body.MaxHoursPercent); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "ok"})
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

// GetRegistrationsByDate returns scheduled staff for a given date (manager)
func (h *ScheduleHandler) GetRegistrationsByDate(c *gin.Context) {
	dateStr := c.Query("date")
	if dateStr == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "date is required"})
		return
	}
	date, err := time.Parse("2006-01-02", dateStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid date format, use YYYY-MM-DD"})
		return
	}
	result, err := h.svc.GetRegistrationsByDate(c.Request.Context(), date)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, result)
}

// GetRegistrationsByPeriod returns all scheduled staff for a given period (manager)
func (h *ScheduleHandler) GetRegistrationsByPeriod(c *gin.Context) {
	var body struct {
		PeriodID string `json:"period_id"`
	}
	if err := c.ShouldBindJSON(&body); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	periodID, err := primitive.ObjectIDFromHex(body.PeriodID)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid period_id"})
		return
	}
	result, err := h.svc.GetRegistrationsByPeriod(c.Request.Context(), periodID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, result)
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

// ─── Weekly template endpoints (manager) ─────────────────────────────────────

func (h *ScheduleHandler) CreateWeeklyTemplate(c *gin.Context) {
	var wt schedule.WeeklyShiftTemplate
	if err := c.ShouldBindJSON(&wt); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	if err := h.svc.CreateWeeklyTemplate(c.Request.Context(), &wt); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusCreated, wt)
}

func (h *ScheduleHandler) GetWeeklyTemplates(c *gin.Context) {
	list, err := h.svc.GetWeeklyTemplates(c.Request.Context())
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, list)
}

func (h *ScheduleHandler) UpdateWeeklyTemplate(c *gin.Context) {
	id, err := primitive.ObjectIDFromHex(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid id"})
		return
	}
	var wt schedule.WeeklyShiftTemplate
	if err := c.ShouldBindJSON(&wt); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	if err := h.svc.UpdateWeeklyTemplate(c.Request.Context(), id, &wt); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, wt)
}

func (h *ScheduleHandler) DeleteWeeklyTemplate(c *gin.Context) {
	id, err := primitive.ObjectIDFromHex(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid id"})
		return
	}
	if err := h.svc.DeleteWeeklyTemplate(c.Request.Context(), id); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "ok"})
}
