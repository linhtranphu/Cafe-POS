package http

import (
	"context"
	"net/http"
	"time"

	"cafe-pos/backend/application/services"
	"cafe-pos/backend/domain/order"
	"cafe-pos/backend/domain/printing"
	"cafe-pos/backend/domain/settings"

	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

// PrintTemplateHandler handles HTTP requests for print template management
type PrintTemplateHandler struct {
	templateRepo     printing.PrintTemplateRepository
	templateRenderer services.TemplateRenderer
	shopSettingsRepo settings.ShopSettingsRepository
}

// NewPrintTemplateHandler creates a new PrintTemplateHandler
func NewPrintTemplateHandler(templateRepo printing.PrintTemplateRepository, templateRenderer services.TemplateRenderer, shopSettingsRepo settings.ShopSettingsRepository) *PrintTemplateHandler {
	return &PrintTemplateHandler{
		templateRepo:     templateRepo,
		templateRenderer: templateRenderer,
		shopSettingsRepo: shopSettingsRepo,
	}
}

// CreateTemplateRequest represents the request body for creating a template
type CreateTemplateRequest struct {
	Type      printing.TemplateType `json:"type" binding:"required"`
	Name      string                `json:"name" binding:"required"`
	Content   string                `json:"content" binding:"required"`
	IsDefault bool                  `json:"is_default"`
}

// UpdateTemplateRequest represents the request body for updating a template
type UpdateTemplateRequest struct {
	Type      printing.TemplateType `json:"type"`
	Name      string                `json:"name"`
	Content   string                `json:"content"`
	IsDefault bool                  `json:"is_default"`
}

// PreviewTemplateRequest represents the request body for previewing a template
type PreviewTemplateRequest struct {
	Content   string `json:"content" binding:"required"`
	Type      printing.TemplateType `json:"type" binding:"required"`
	ItemIndex int    `json:"item_index"` // For label preview
}

// ListTemplates handles GET /api/print-templates
func (h *PrintTemplateHandler) ListTemplates(c *gin.Context) {
	ctx := c.Request.Context()

	// Get query parameter for filtering by type
	templateType := c.Query("type")

	var templates []*printing.PrintTemplate
	var err error

	if templateType != "" {
		templates, err = h.templateRepo.FindByType(ctx, printing.TemplateType(templateType))
	} else {
		// Get all templates by fetching both types
		billTemplates, err1 := h.templateRepo.FindByType(ctx, printing.TemplateTypeBill)
		labelTemplates, err2 := h.templateRepo.FindByType(ctx, printing.TemplateTypeLabel)
		
		if err1 != nil || err2 != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch templates"})
			return
		}
		
		templates = append(billTemplates, labelTemplates...)
	}

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch templates"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"templates": templates})
}

// GetTemplate handles GET /api/print-templates/:id
func (h *PrintTemplateHandler) GetTemplate(c *gin.Context) {
	ctx := c.Request.Context()

	id, err := primitive.ObjectIDFromHex(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid template ID format"})
		return
	}

	template, err := h.templateRepo.FindByID(ctx, id)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			c.JSON(http.StatusNotFound, gin.H{"error": "Template not found"})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch template"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"template": template})
}

// CreateTemplate handles POST /api/print-templates
func (h *PrintTemplateHandler) CreateTemplate(c *gin.Context) {
	ctx := c.Request.Context()

	var req CreateTemplateRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// If setting as default, unset other defaults of the same type
	if req.IsDefault {
		if err := h.unsetDefaultTemplates(ctx, req.Type); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update default templates"})
			return
		}
	}

	// Create template
	template := &printing.PrintTemplate{
		Type:      req.Type,
		Name:      req.Name,
		Content:   req.Content,
		IsDefault: req.IsDefault,
	}

	if err := h.templateRepo.Create(ctx, template); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create template"})
		return
	}

	c.JSON(http.StatusCreated, gin.H{"template": template})
}

// UpdateTemplate handles PUT /api/print-templates/:id
func (h *PrintTemplateHandler) UpdateTemplate(c *gin.Context) {
	ctx := c.Request.Context()

	id, err := primitive.ObjectIDFromHex(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid template ID format"})
		return
	}

	// Fetch existing template
	existing, err := h.templateRepo.FindByID(ctx, id)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			c.JSON(http.StatusNotFound, gin.H{"error": "Template not found"})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch template"})
		return
	}

	var req UpdateTemplateRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Update fields if provided
	if req.Type != "" {
		existing.Type = req.Type
	}
	if req.Name != "" {
		existing.Name = req.Name
	}
	if req.Content != "" {
		existing.Content = req.Content
	}
	existing.IsDefault = req.IsDefault

	// If setting as default, unset other defaults of the same type
	if req.IsDefault {
		if err := h.unsetDefaultTemplates(ctx, existing.Type); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update default templates"})
			return
		}
	}

	if err := h.templateRepo.Update(ctx, existing); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update template"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"template": existing})
}

// DeleteTemplate handles DELETE /api/print-templates/:id
func (h *PrintTemplateHandler) DeleteTemplate(c *gin.Context) {
	ctx := c.Request.Context()

	id, err := primitive.ObjectIDFromHex(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid template ID format"})
		return
	}

	// Check if it's a default template
	template, err := h.templateRepo.FindByID(ctx, id)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			c.JSON(http.StatusNotFound, gin.H{"error": "Template not found"})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch template"})
		return
	}

	if template.IsDefault {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Cannot delete default template"})
		return
	}

	if err := h.templateRepo.Delete(ctx, id); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to delete template"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Template deleted successfully"})
}

// PreviewTemplate handles POST /api/print-templates/:id/preview
func (h *PrintTemplateHandler) PreviewTemplate(c *gin.Context) {
	var req PreviewTemplateRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	ctx := c.Request.Context()

	// Fetch shop settings
	shopSettings, err := h.shopSettingsRepo.GetSettings(ctx)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch shop settings"})
		return
	}

	// Create sample order for preview
	sampleOrder := h.createSampleOrder()

	// Create temporary template
	template := &printing.PrintTemplate{
		Type:    req.Type,
		Content: req.Content,
	}

	var content string

	// Render based on type
	if req.Type == printing.TemplateTypeBill {
		content, err = h.templateRenderer.RenderBill(sampleOrder, template, shopSettings)
	} else if req.Type == printing.TemplateTypeLabel {
		itemIndex := req.ItemIndex
		if itemIndex < 0 || itemIndex >= len(sampleOrder.Items) {
			itemIndex = 0
		}
		content, err = h.templateRenderer.RenderLabel(sampleOrder, itemIndex, template, shopSettings)
	} else {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid template type"})
		return
	}

	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error":   "Template rendering failed",
			"details": err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"preview": content,
		"success": true,
	})
}

// unsetDefaultTemplates unsets the is_default flag for all templates of a specific type
func (h *PrintTemplateHandler) unsetDefaultTemplates(ctx context.Context, templateType printing.TemplateType) error {
	templates, err := h.templateRepo.FindByType(ctx, templateType)
	if err != nil {
		return err
	}

	for _, template := range templates {
		if template.IsDefault {
			template.IsDefault = false
			if err := h.templateRepo.Update(ctx, template); err != nil {
				return err
			}
		}
	}

	return nil
}

// createSampleOrder creates a sample order for template preview
func (h *PrintTemplateHandler) createSampleOrder() *order.Order {
	now := time.Now()
	return &order.Order{
		ID:          primitive.NewObjectID(),
		OrderNumber: "ORD-001",
		Items: []order.OrderItem{
			{
				Name:        "Cà phê sữa đá",
				VariantName: "Size M",
				Quantity:    2,
				Price:       25000,
				Subtotal:    50000,
				Note:        "Ít đường",
			},
			{
				Name:        "Trà sữa trân châu",
				VariantName: "Size L",
				Quantity:    1,
				Price:       35000,
				Subtotal:    35000,
			},
		},
		Subtotal:    85000,
		Discount:    5000,
		Total:       80000,
		WaiterName:  "Nguyễn Văn A",
		CreatedAt:   now,
		UpdatedAt:   now,
	}
}
