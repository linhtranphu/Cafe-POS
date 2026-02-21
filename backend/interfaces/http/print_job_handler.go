package http

import (
	"net/http"

	"cafe-pos/backend/application/services"
	"cafe-pos/backend/domain/printing"

	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

// PrintJobHandler handles HTTP requests for print job management
type PrintJobHandler struct {
	printService  services.PrintService
	printJobRepo  printing.PrintJobRepository
	wsBroadcaster services.WebSocketBroadcaster
}

// NewPrintJobHandler creates a new PrintJobHandler
func NewPrintJobHandler(printService services.PrintService, printJobRepo printing.PrintJobRepository, wsBroadcaster services.WebSocketBroadcaster) *PrintJobHandler {
	return &PrintJobHandler{
		printService:  printService,
		printJobRepo:  printJobRepo,
		wsBroadcaster: wsBroadcaster,
	}
}

// ListPrintJobs handles GET /api/print-jobs
func (h *PrintJobHandler) ListPrintJobs(c *gin.Context) {
	ctx := c.Request.Context()

	// Get query parameters for filtering
	status := c.Query("status")
	orderID := c.Query("order_id")

	var jobs []*printing.PrintJob
	var err error

	// Filter by order_id if provided
	if orderID != "" {
		oid, err := primitive.ObjectIDFromHex(orderID)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid order_id format"})
			return
		}
		jobs, err = h.printJobRepo.FindByOrderID(ctx, oid)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch print jobs"})
			return
		}
	} else if status == "pending" {
		jobs, err = h.printService.GetPendingJobs(ctx)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch pending jobs"})
			return
		}
	} else if status == "failed" {
		jobs, err = h.printService.GetFailedJobs(ctx)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch failed jobs"})
			return
		}
	} else {
		// Return all pending and failed jobs by default
		pendingJobs, err := h.printService.GetPendingJobs(ctx)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch pending jobs"})
			return
		}
		failedJobs, err := h.printService.GetFailedJobs(ctx)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch failed jobs"})
			return
		}
		jobs = append(pendingJobs, failedJobs...)
	}

	c.JSON(http.StatusOK, gin.H{"jobs": jobs})
}

// GetPrintJob handles GET /api/print-jobs/:id
func (h *PrintJobHandler) GetPrintJob(c *gin.Context) {
	ctx := c.Request.Context()

	id, err := primitive.ObjectIDFromHex(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid job ID format"})
		return
	}

	job, err := h.printJobRepo.FindByID(ctx, id)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Print job not found"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"job": job})
}

// GetPendingJobs handles GET /api/print-jobs/pending
func (h *PrintJobHandler) GetPendingJobs(c *gin.Context) {
	ctx := c.Request.Context()

	jobs, err := h.printService.GetPendingJobs(ctx)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch pending jobs"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"jobs": jobs})
}

// GetFailedJobs handles GET /api/print-jobs/failed
func (h *PrintJobHandler) GetFailedJobs(c *gin.Context) {
	ctx := c.Request.Context()

	jobs, err := h.printService.GetFailedJobs(ctx)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch failed jobs"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"jobs": jobs})
}

// RetryPrintJob handles POST /api/print-jobs/:id/retry
func (h *PrintJobHandler) RetryPrintJob(c *gin.Context) {
	ctx := c.Request.Context()

	id, err := primitive.ObjectIDFromHex(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid job ID format"})
		return
	}

	if err := h.printService.RetryJob(ctx, id); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Print job queued for retry"})
}

// CancelPrintJob handles DELETE /api/print-jobs/:id
func (h *PrintJobHandler) CancelPrintJob(c *gin.Context) {
	ctx := c.Request.Context()

	id, err := primitive.ObjectIDFromHex(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid job ID format"})
		return
	}

	if err := h.printService.CancelJob(ctx, id); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Print job cancelled"})
}

// UpdatePrintJobStatus handles PUT /api/print-jobs/:id/status
// This endpoint is called by the local print bridge to update job status
func (h *PrintJobHandler) UpdatePrintJobStatus(c *gin.Context) {
	ctx := c.Request.Context()

	id, err := primitive.ObjectIDFromHex(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid job ID format"})
		return
	}

	// Parse request body
	var req struct {
		Status   string `json:"status" binding:"required"`
		ErrorMsg string `json:"error_msg"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request body"})
		return
	}

	// Validate status
	validStatuses := map[string]bool{
		"PENDING":   true,
		"PRINTING":  true,
		"COMPLETED": true,
		"FAILED":    true,
	}

	if !validStatuses[req.Status] {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid status value"})
		return
	}

	// Get the job
	job, err := h.printJobRepo.FindByID(ctx, id)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Print job not found"})
		return
	}

	// Update status
	if err := h.printJobRepo.UpdateStatus(ctx, id, printing.PrintJobStatus(req.Status), req.ErrorMsg); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update job status"})
		return
	}

	// Broadcast WebSocket event for status change
	if h.wsBroadcaster != nil {
		h.wsBroadcaster.BroadcastPrintJobStatusChanged(job.ID.Hex(), printing.PrintJobStatus(req.Status), req.ErrorMsg)
		
		// Also broadcast failed event if status is FAILED
		if req.Status == "FAILED" {
			h.wsBroadcaster.BroadcastPrintJobFailed(job.ID.Hex(), req.ErrorMsg)
		}
	}
	
	c.JSON(http.StatusOK, gin.H{
		"message": "Print job status updated",
		"job_id":  job.ID.Hex(),
		"status":  req.Status,
	})
}
