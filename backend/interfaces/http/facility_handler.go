package http

import (
	"net/http"
	"strconv"
	"time"
	"cafe-pos/backend/application/services"
	"cafe-pos/backend/domain/facility"
	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type FacilityHandler struct {
	service *services.FacilityService
}

func NewFacilityHandler(service *services.FacilityService) *FacilityHandler {
	return &FacilityHandler{service: service}
}

func (h *FacilityHandler) GetAllFacilities(c *gin.Context) {
	facilities, err := h.service.GetAllFacilities(c.Request.Context())
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, facilities)
}

func (h *FacilityHandler) GetFacility(c *gin.Context) {
	id, err := primitive.ObjectIDFromHex(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "ID không hợp lệ"})
		return
	}

	f, err := h.service.GetFacilityByID(c.Request.Context(), id)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Không tìm thấy tài sản"})
		return
	}
	c.JSON(http.StatusOK, f)
}

func (h *FacilityHandler) CreateFacility(c *gin.Context) {
	var f facility.Facility
	if err := c.ShouldBindJSON(&f); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	userID, _ := primitive.ObjectIDFromHex(c.GetString("user_id"))
	username := c.GetString("username")

	if err := h.service.CreateFacility(c.Request.Context(), &f, userID, username); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusCreated, f)
}

func (h *FacilityHandler) UpdateFacility(c *gin.Context) {
	id, err := primitive.ObjectIDFromHex(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "ID không hợp lệ"})
		return
	}

	var f facility.Facility
	if err := c.ShouldBindJSON(&f); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	userID, _ := primitive.ObjectIDFromHex(c.GetString("user_id"))
	username := c.GetString("username")

	if err := h.service.UpdateFacility(c.Request.Context(), id, &f, userID, username); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, f)
}

func (h *FacilityHandler) DeleteFacility(c *gin.Context) {
	id, err := primitive.ObjectIDFromHex(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "ID không hợp lệ"})
		return
	}

	userID, _ := primitive.ObjectIDFromHex(c.GetString("user_id"))
	username := c.GetString("username")

	if err := h.service.DeleteFacility(c.Request.Context(), id, userID, username); err != nil {
		// Check if it's a business rule validation error
		if err.Error() == "không thể xóa tài sản đã có lịch sử bảo trì" {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "Xóa thành công"})
}

func (h *FacilityHandler) CreateMaintenanceRecord(c *gin.Context) {
	var record facility.MaintenanceRecord
	if err := c.ShouldBindJSON(&record); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	record.UserID, _ = primitive.ObjectIDFromHex(c.GetString("user_id"))
	record.Username = c.GetString("username")

	if err := h.service.CreateMaintenanceRecord(c.Request.Context(), &record); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusCreated, record)
}

func (h *FacilityHandler) GetMaintenanceHistory(c *gin.Context) {
	id, err := primitive.ObjectIDFromHex(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "ID không hợp lệ"})
		return
	}

	records, err := h.service.GetMaintenanceHistory(c.Request.Context(), id)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, records)
}

func (h *FacilityHandler) CreateIssueReport(c *gin.Context) {
	var report facility.IssueReport
	if err := c.ShouldBindJSON(&report); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	report.ReportedBy, _ = primitive.ObjectIDFromHex(c.GetString("user_id"))
	report.Username = c.GetString("username")

	if err := h.service.CreateIssueReport(c.Request.Context(), &report); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusCreated, report)
}

func (h *FacilityHandler) GetIssueReports(c *gin.Context) {
	reports, err := h.service.GetIssueReports(c.Request.Context())
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, reports)
}

func (h *FacilityHandler) GetFacilityHistory(c *gin.Context) {
	id, err := primitive.ObjectIDFromHex(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "ID không hợp lệ"})
		return
	}

	history, err := h.service.GetFacilityHistory(c.Request.Context(), id)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, history)
}

func (h *FacilityHandler) GetHistoryWithFilter(c *gin.Context) {
	var filter facility.FacilityHistoryFilter
	
	// Parse query parameters
	if facilityIDStr := c.Query("facility_id"); facilityIDStr != "" {
		if facilityID, err := primitive.ObjectIDFromHex(facilityIDStr); err == nil {
			filter.FacilityID = &facilityID
		}
	}
	
	if action := c.Query("action"); action != "" {
		filter.Action = &action
	}
	
	if dateFromStr := c.Query("date_from"); dateFromStr != "" {
		if dateFrom, err := time.Parse("2006-01-02", dateFromStr); err == nil {
			filter.DateFrom = &dateFrom
		}
	}
	
	if dateToStr := c.Query("date_to"); dateToStr != "" {
		if dateTo, err := time.Parse("2006-01-02", dateToStr); err == nil {
			filter.DateTo = &dateTo
		}
	}
	
	if limit := c.Query("limit"); limit != "" {
		if l, err := strconv.Atoi(limit); err == nil {
			filter.Limit = l
		}
	} else {
		filter.Limit = 50 // Default limit
	}
	
	if offset := c.Query("offset"); offset != "" {
		if o, err := strconv.Atoi(offset); err == nil {
			filter.Offset = o
		}
	}

	history, err := h.service.GetHistoryWithFilter(c.Request.Context(), filter)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, history)
}

func (h *FacilityHandler) SearchFacilities(c *gin.Context) {
	var filter facility.FacilityFilter
	
	// Parse query parameters
	if name := c.Query("name"); name != "" {
		filter.Name = &name
	}
	
	if facilityType := c.Query("type"); facilityType != "" {
		filter.Type = &facilityType
	}
	
	if area := c.Query("area"); area != "" {
		filter.Area = &area
	}
	
	if status := c.Query("status"); status != "" {
		filter.Status = &status
	}
	
	if limit := c.Query("limit"); limit != "" {
		if l, err := strconv.Atoi(limit); err == nil {
			filter.Limit = l
		}
	} else {
		filter.Limit = 20 // Default limit
	}
	
	if offset := c.Query("offset"); offset != "" {
		if o, err := strconv.Atoi(offset); err == nil {
			filter.Offset = o
		}
	}

	facilities, total, err := h.service.SearchFacilities(c.Request.Context(), filter)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	
	c.JSON(http.StatusOK, gin.H{
		"data":   facilities,
		"total":  total,
		"limit":  filter.Limit,
		"offset": filter.Offset,
	})
}

func (h *FacilityHandler) GetScheduledMaintenance(c *gin.Context) {
	tasks, err := h.service.GetScheduledMaintenance(c.Request.Context())
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, tasks)
}

func (h *FacilityHandler) GetNextMaintenanceDate(c *gin.Context) {
	id, err := primitive.ObjectIDFromHex(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "ID không hợp lệ"})
		return
	}

	result, err := h.service.GetNextMaintenanceDate(c.Request.Context(), id)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, result)
}

func (h *FacilityHandler) GetMaintenanceDue(c *gin.Context) {
	tasks, err := h.service.GetMaintenanceDue(c.Request.Context())
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, tasks)
}

func (h *FacilityHandler) GetMaintenanceStats(c *gin.Context) {
	id, err := primitive.ObjectIDFromHex(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "ID không hợp lệ"})
		return
	}

	stats, err := h.service.GetMaintenanceStats(c.Request.Context(), id)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, stats)
}

func (h *FacilityHandler) GetStatusHistory(c *gin.Context) {
	id, err := primitive.ObjectIDFromHex(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "ID không hợp lệ"})
		return
	}

	history, err := h.service.GetStatusHistory(c.Request.Context(), id)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, history)
}

// FacilityType handlers
func (h *FacilityHandler) CreateFacilityType(c *gin.Context) {
	var req struct {
		Name string `json:"name" binding:"required"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	ft, err := h.service.CreateFacilityType(c.Request.Context(), req.Name)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusCreated, ft)
}

func (h *FacilityHandler) GetFacilityTypes(c *gin.Context) {
	types, err := h.service.GetFacilityTypes(c.Request.Context())
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, types)
}

func (h *FacilityHandler) DeleteFacilityType(c *gin.Context) {
	id, err := primitive.ObjectIDFromHex(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "ID không hợp lệ"})
		return
	}

	if err := h.service.DeleteFacilityType(c.Request.Context(), id); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Không thể xóa loại thiết bị đang được sử dụng"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "Xóa thành công"})
}

// FacilityArea handlers
func (h *FacilityHandler) CreateFacilityArea(c *gin.Context) {
	var req struct {
		Name string `json:"name" binding:"required"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	fa, err := h.service.CreateFacilityArea(c.Request.Context(), req.Name)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusCreated, fa)
}

func (h *FacilityHandler) GetFacilityAreas(c *gin.Context) {
	areas, err := h.service.GetFacilityAreas(c.Request.Context())
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, areas)
}

func (h *FacilityHandler) DeleteFacilityArea(c *gin.Context) {
	id, err := primitive.ObjectIDFromHex(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "ID không hợp lệ"})
		return
	}

	if err := h.service.DeleteFacilityArea(c.Request.Context(), id); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Không thể xóa khu vực đang được sử dụng"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "Xóa thành công"})
}
