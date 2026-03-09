package services

import (
	"bytes"
	"context"
	"encoding/base64"
	"fmt"
	"html/template"
	"log"
	"os"
	"path/filepath"
	"time"

	"cafe-pos/backend/domain/order"
	"cafe-pos/backend/domain/printing"
	"cafe-pos/backend/domain/settings"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

// PrintService defines the interface for print job management
type PrintService interface {
	// CreatePrintJobsForOrder creates print jobs (1 bill + N labels) for an order
	CreatePrintJobsForOrder(ctx context.Context, ord *order.Order) error

	// CreateTempBillJob creates a temporary bill print job for an order
	CreateTempBillJob(ctx context.Context, ord *order.Order) error

	// ReprintBill creates a new bill print job for an existing order
	ReprintBill(ctx context.Context, orderID primitive.ObjectID) error

	// ReprintLabel creates a new label print job for a specific order item
	ReprintLabel(ctx context.Context, orderID primitive.ObjectID, itemIndex int) error

	// GetPendingJobs returns all pending print jobs
	GetPendingJobs(ctx context.Context) ([]*printing.PrintJob, error)

	// GetFailedJobs returns all failed print jobs
	GetFailedJobs(ctx context.Context) ([]*printing.PrintJob, error)

	// RetryJob retries a failed print job
	RetryJob(ctx context.Context, jobID primitive.ObjectID) error

	// CancelJob cancels a pending print job
	CancelJob(ctx context.Context, jobID primitive.ObjectID) error
}

// WebSocketBroadcaster defines the interface for broadcasting WebSocket events
type WebSocketBroadcaster interface {
	BroadcastPrintJobCreated(job *printing.PrintJob)
	BroadcastPrintJobStatusChanged(jobID string, status printing.PrintJobStatus, errorMsg string)
	BroadcastPrintJobFailed(jobID string, errorMsg string)
}

// printService implements the PrintService interface
type printService struct {
	printJobRepo       printing.PrintJobRepository
	printerConfigRepo  printing.PrinterConfigRepository
	templateRepo       printing.PrintTemplateRepository
	templateRenderer   TemplateRenderer
	chromedpRenderer   *ChromedpBillRendererOptimized
	orderRepo          OrderRepository // To fetch order data for reprints
	shopSettingsRepo   settings.ShopSettingsRepository
	wsBroadcaster      WebSocketBroadcaster
}

// PrintServiceConfig contains configuration for the print service
type PrintServiceConfig struct {
	PrintJobRepo      printing.PrintJobRepository
	PrinterConfigRepo printing.PrinterConfigRepository
	TemplateRepo      printing.PrintTemplateRepository
	TemplateRenderer  TemplateRenderer
	OrderRepo         OrderRepository
	ShopSettingsRepo  settings.ShopSettingsRepository
	WSBroadcaster     WebSocketBroadcaster
}

// NewPrintService creates a new print service
func NewPrintService(config PrintServiceConfig) PrintService {
	// Initialize chromedp renderer
	chromedpRenderer, err := NewChromedpBillRendererOptimized()
	if err != nil {
		log.Printf("Warning: Failed to initialize Chromedp renderer: %v", err)
	}
	
	return &printService{
		printJobRepo:      config.PrintJobRepo,
		printerConfigRepo: config.PrinterConfigRepo,
		templateRepo:      config.TemplateRepo,
		templateRenderer:  config.TemplateRenderer,
		chromedpRenderer:  chromedpRenderer,
		orderRepo:         config.OrderRepo,
		shopSettingsRepo:  config.ShopSettingsRepo,
		wsBroadcaster:     config.WSBroadcaster,
	}
}

// CreatePrintJobsForOrder creates print jobs for an order (1 bill + N labels)
func (s *printService) CreatePrintJobsForOrder(ctx context.Context, ord *order.Order) error {
	if ord == nil {
		return fmt.Errorf("order cannot be nil")
	}

	log.Printf("[PRINT] Creating print jobs for order %s (ID: %s)", ord.OrderNumber, ord.ID.Hex())

	// Get default printers
	billPrinter, err := s.printerConfigRepo.FindDefault(ctx, printing.PrinterTypeBill)
	if err != nil {
		log.Printf("[PRINT ERROR] Failed to get default bill printer for order %s: %v", ord.OrderNumber, err)
		return fmt.Errorf("failed to get default bill printer: %w", err)
	}
	if billPrinter == nil {
		log.Printf("[PRINT ERROR] No default bill printer configured for order %s", ord.OrderNumber)
		return fmt.Errorf("no default bill printer configured")
	}

	labelPrinter, err := s.printerConfigRepo.FindDefault(ctx, printing.PrinterTypeLabel)
	if err != nil {
		log.Printf("[PRINT ERROR] Failed to get default label printer for order %s: %v", ord.OrderNumber, err)
		return fmt.Errorf("failed to get default label printer: %w", err)
	}
	if labelPrinter == nil {
		log.Printf("[PRINT ERROR] No default label printer configured for order %s", ord.OrderNumber)
		return fmt.Errorf("no default label printer configured")
	}

	// Get default templates (only needed if not using print bridge)
	// When using print bridge, we render HTML directly from file
	var billTemplate *printing.PrintTemplate
	var labelTemplate *printing.PrintTemplate
	
	// Check if using print bridge
	shopSettings, err := s.shopSettingsRepo.GetSettings(ctx)
	if err == nil && shopSettings != nil && shopSettings.PrintBridgeURL != "" {
		// Using print bridge - templates not needed
		log.Printf("[PRINT] Using print bridge, skipping template lookup")
	} else {
		// Not using print bridge - need templates from database
		billTemplate, err = s.templateRepo.FindDefault(ctx, printing.TemplateTypeBill)
		if err != nil {
			log.Printf("[PRINT ERROR] Failed to get default bill template for order %s: %v", ord.OrderNumber, err)
			return fmt.Errorf("failed to get default bill template: %w", err)
		}

		labelTemplate, err = s.templateRepo.FindDefault(ctx, printing.TemplateTypeLabel)
		if err != nil {
			log.Printf("[PRINT ERROR] Failed to get default label template for order %s: %v", ord.OrderNumber, err)
			return fmt.Errorf("failed to get default label template: %w", err)
		}
	}

	// Create bill print job
	if err := s.createBillJob(ctx, ord, billPrinter, billTemplate); err != nil {
		log.Printf("[PRINT ERROR] Failed to create bill job for order %s: %v", ord.OrderNumber, err)
		return fmt.Errorf("failed to create bill job: %w", err)
	}

	// Create label print jobs for each item
	for i := range ord.Items {
		if err := s.createLabelJob(ctx, ord, i, labelPrinter, labelTemplate); err != nil {
			log.Printf("[PRINT ERROR] Failed to create label job for item %d in order %s: %v", i, ord.OrderNumber, err)
			return fmt.Errorf("failed to create label job for item %d: %w", i, err)
		}
	}

	log.Printf("[PRINT] Successfully created %d print jobs for order %s (1 bill + %d labels)", 1+len(ord.Items), ord.OrderNumber, len(ord.Items))
	return nil
}

// ReprintBill creates a new bill print job for an existing order
func (s *printService) ReprintBill(ctx context.Context, orderID primitive.ObjectID) error {
	// Fetch order data
	ord, err := s.orderRepo.FindByID(ctx, orderID)
	if err != nil {
		return fmt.Errorf("failed to fetch order: %w", err)
	}
	if ord == nil {
		return fmt.Errorf("order not found: %s", orderID.Hex())
	}

	// Get default printer and template
	billPrinter, err := s.printerConfigRepo.FindDefault(ctx, printing.PrinterTypeBill)
	if err != nil {
		return fmt.Errorf("failed to get default bill printer: %w", err)
	}
	if billPrinter == nil {
		return fmt.Errorf("no default bill printer configured")
	}

	billTemplate, err := s.templateRepo.FindDefault(ctx, printing.TemplateTypeBill)
	if err != nil {
		return fmt.Errorf("failed to get default bill template: %w", err)
	}

	// Create new bill job
	if err := s.createBillJob(ctx, ord, billPrinter, billTemplate); err != nil {
		return fmt.Errorf("failed to create bill job: %w", err)
	}

	return nil
}

// ReprintLabel creates a new label print job for a specific order item
func (s *printService) ReprintLabel(ctx context.Context, orderID primitive.ObjectID, itemIndex int) error {
	// Fetch order data
	ord, err := s.orderRepo.FindByID(ctx, orderID)
	if err != nil {
		return fmt.Errorf("failed to fetch order: %w", err)
	}
	if ord == nil {
		return fmt.Errorf("order not found: %s", orderID.Hex())
	}

	// Validate item index
	if itemIndex < 0 || itemIndex >= len(ord.Items) {
		return fmt.Errorf("invalid item index: %d (order has %d items)", itemIndex, len(ord.Items))
	}

	// Get default printer and template
	labelPrinter, err := s.printerConfigRepo.FindDefault(ctx, printing.PrinterTypeLabel)
	if err != nil {
		return fmt.Errorf("failed to get default label printer: %w", err)
	}
	if labelPrinter == nil {
		return fmt.Errorf("no default label printer configured")
	}

	labelTemplate, err := s.templateRepo.FindDefault(ctx, printing.TemplateTypeLabel)
	if err != nil {
		return fmt.Errorf("failed to get default label template: %w", err)
	}

	// Create new label job
	if err := s.createLabelJob(ctx, ord, itemIndex, labelPrinter, labelTemplate); err != nil {
		return fmt.Errorf("failed to create label job: %w", err)
	}

	return nil
}

// GetPendingJobs returns all pending print jobs
func (s *printService) GetPendingJobs(ctx context.Context) ([]*printing.PrintJob, error) {
	jobs, err := s.printJobRepo.FindPending(ctx, 100) // Limit to 100 jobs
	if err != nil {
		return nil, fmt.Errorf("failed to fetch pending jobs: %w", err)
	}
	return jobs, nil
}

// GetFailedJobs returns all failed print jobs
func (s *printService) GetFailedJobs(ctx context.Context) ([]*printing.PrintJob, error) {
	jobs, err := s.printJobRepo.FindFailed(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch failed jobs: %w", err)
	}
	return jobs, nil
}

// RetryJob retries a failed print job
func (s *printService) RetryJob(ctx context.Context, jobID primitive.ObjectID) error {
	// Fetch the job
	job, err := s.printJobRepo.FindByID(ctx, jobID)
	if err != nil {
		return fmt.Errorf("failed to fetch job: %w", err)
	}
	if job == nil {
		return fmt.Errorf("job not found: %s", jobID.Hex())
	}

	// Only retry failed jobs
	if job.Status != printing.PrintJobStatusFailed {
		return fmt.Errorf("can only retry failed jobs (current status: %s)", job.Status)
	}

	// Reset retry count and status to pending
	job.RetryCount = 0
	job.Status = printing.PrintJobStatusPending
	job.ErrorMsg = ""
	job.UpdatedAt = time.Now()

	// Update the job
	if err := s.printJobRepo.UpdateStatus(ctx, jobID, printing.PrintJobStatusPending, ""); err != nil {
		return fmt.Errorf("failed to update job status: %w", err)
	}

	return nil
}

// CancelJob cancels a pending print job
func (s *printService) CancelJob(ctx context.Context, jobID primitive.ObjectID) error {
	// Fetch the job
	job, err := s.printJobRepo.FindByID(ctx, jobID)
	if err != nil {
		return fmt.Errorf("failed to fetch job: %w", err)
	}
	if job == nil {
		return fmt.Errorf("job not found: %s", jobID.Hex())
	}

	// Only cancel pending jobs
	if job.Status != printing.PrintJobStatusPending {
		return fmt.Errorf("can only cancel pending jobs (current status: %s)", job.Status)
	}

	// Delete the job
	if err := s.printJobRepo.Delete(ctx, jobID); err != nil {
		return fmt.Errorf("failed to delete job: %w", err)
	}

	return nil
}

// createBillJob creates a bill print job
func (s *printService) createBillJob(ctx context.Context, ord *order.Order, printer *printing.PrinterConfig, template *printing.PrintTemplate) error {
	// Fetch shop settings
	shopSettings, err := s.shopSettingsRepo.GetSettings(ctx)
	if err != nil {
		log.Printf("[PRINT ERROR] Failed to fetch shop settings for bill - order_id=%s, order_number=%s, error=%v, timestamp=%s",
			ord.ID.Hex(), ord.OrderNumber, err, time.Now().Format(time.RFC3339))
		return fmt.Errorf("failed to fetch shop settings: %w", err)
	}

	var content string
	var contentType string

	// When using print bridge, render HTML and let print bridge handle conversion
	// Print bridge will convert HTML → image → ESC/POS
	if shopSettings.PrintBridgeURL != "" {
		// Render HTML bill (print bridge will handle the rest)
		htmlContent, err := s.renderBillHTML(ord, shopSettings)
		if err != nil {
			log.Printf("[PRINT ERROR] Failed to render HTML for bill - order_id=%s, order_number=%s, error=%v, timestamp=%s",
				ord.ID.Hex(), ord.OrderNumber, err, time.Now().Format(time.RFC3339))
			return fmt.Errorf("failed to render HTML bill: %w", err)
		}
		content = htmlContent
		contentType = "html"
		log.Printf("[PRINT] Using HTML for print bridge - order_id=%s, order_number=%s, bridge_url=%s", 
			ord.ID.Hex(), ord.OrderNumber, shopSettings.PrintBridgeURL)
	} else {
		// Fallback: No print bridge, try local renderers
		if s.chromedpRenderer != nil {
			escposData, err := s.chromedpRenderer.RenderBillToESCPOS(ord, shopSettings)
			if err != nil {
				log.Printf("[PRINT ERROR] Chromedp rendering failed - order_id=%s, error=%v", ord.OrderNumber, err)
				return fmt.Errorf("failed to render bill with chromedp: %w", err)
			}
			content = base64.StdEncoding.EncodeToString(escposData)
			contentType = "binary"
			log.Printf("[PRINT] Using Chromedp renderer - order_id=%s", ord.OrderNumber)
		} else {
			// Last resort: text template
			content, err = s.templateRenderer.RenderBill(ord, template, shopSettings)
			if err != nil {
				log.Printf("[PRINT ERROR] Text template rendering failed - order_id=%s, error=%v", ord.OrderNumber, err)
				return fmt.Errorf("failed to render bill: %w", err)
			}
			contentType = "text"
			log.Printf("[PRINT] Using text template - order_id=%s", ord.OrderNumber)
		}
	}

	// Create print job
	job := &printing.PrintJob{
		Type:        printing.PrintJobTypeBill,
		OrderID:     ord.ID,
		OrderNumber: ord.OrderNumber,
		PrinterID:   printer.ID,
		Content:     content,
		ContentType: contentType,
		Status:      printing.PrintJobStatusPending,
		RetryCount:  0,
		MaxRetries:  3,
		CreatedAt:   time.Now(),
		UpdatedAt:   time.Now(),
	}

	if err := s.printJobRepo.Create(ctx, job); err != nil {
		log.Printf("[PRINT ERROR] Failed to save print job - job_type=BILL, order_id=%s, order_number=%s, error=%v, timestamp=%s",
			ord.ID.Hex(), ord.OrderNumber, err, time.Now().Format(time.RFC3339))
		return fmt.Errorf("failed to create print job: %w", err)
	}

	log.Printf("[PRINT] Bill job created - job_id=%s, order_id=%s, order_number=%s, printer=%s, content_type=%s, timestamp=%s",
		job.ID.Hex(), ord.ID.Hex(), ord.OrderNumber, printer.Name, contentType, time.Now().Format(time.RFC3339))
	
	// Broadcast WebSocket event
	if s.wsBroadcaster != nil {
		s.wsBroadcaster.BroadcastPrintJobCreated(job)
	}
	
	return nil
}

// createLabelJob creates a label print job for a specific item
func (s *printService) createLabelJob(ctx context.Context, ord *order.Order, itemIndex int, printer *printing.PrinterConfig, template *printing.PrintTemplate) error {
	// Fetch shop settings
	shopSettings, err := s.shopSettingsRepo.GetSettings(ctx)
	if err != nil {
		log.Printf("[PRINT ERROR] Failed to fetch shop settings for label - order_id=%s, order_number=%s, item_index=%d, error=%v, timestamp=%s",
			ord.ID.Hex(), ord.OrderNumber, itemIndex, err, time.Now().Format(time.RFC3339))
		return fmt.Errorf("failed to fetch shop settings: %w", err)
	}

	// Render label content
	content, err := s.templateRenderer.RenderLabel(ord, itemIndex, template, shopSettings)
	if err != nil {
		log.Printf("[PRINT ERROR] Template rendering failed for label - order_id=%s, order_number=%s, item_index=%d, error=%v, timestamp=%s",
			ord.ID.Hex(), ord.OrderNumber, itemIndex, err, time.Now().Format(time.RFC3339))
		return fmt.Errorf("failed to render label: %w", err)
	}

	// Create print job
	job := &printing.PrintJob{
		Type:        printing.PrintJobTypeLabel,
		OrderID:     ord.ID,
		OrderNumber: ord.OrderNumber,
		PrinterID:   printer.ID,
		Content:     content,
		Status:      printing.PrintJobStatusPending,
		RetryCount:  0,
		MaxRetries:  3,
		CreatedAt:   time.Now(),
		UpdatedAt:   time.Now(),
	}

	if err := s.printJobRepo.Create(ctx, job); err != nil {
		log.Printf("[PRINT ERROR] Failed to save print job - job_type=LABEL, order_id=%s, order_number=%s, item_index=%d, error=%v, timestamp=%s",
			ord.ID.Hex(), ord.OrderNumber, itemIndex, err, time.Now().Format(time.RFC3339))
		return fmt.Errorf("failed to create print job: %w", err)
	}

	log.Printf("[PRINT] Label job created - job_id=%s, order_id=%s, order_number=%s, item_index=%d, printer=%s, timestamp=%s",
		job.ID.Hex(), ord.ID.Hex(), ord.OrderNumber, itemIndex, printer.Name, time.Now().Format(time.RFC3339))
	
	// Broadcast WebSocket event
	if s.wsBroadcaster != nil {
		s.wsBroadcaster.BroadcastPrintJobCreated(job)
	}
	
	return nil
}

// renderBillHTML renders bill HTML from template file
func (s *printService) renderBillHTML(ord *order.Order, shopSettings *settings.ShopSettings) (string, error) {
	// Read HTML template file
	templatePath := "./application/services/templates/bill_template_optimized.html"
	templateContent, err := os.ReadFile(templatePath)
	if err != nil {
		return "", fmt.Errorf("failed to read HTML template: %w", err)
	}

	// Parse template
	tmpl, err := template.New("bill").Parse(string(templateContent))
	if err != nil {
		return "", fmt.Errorf("failed to parse HTML template: %w", err)
	}

	// Prepare data
	data := s.prepareBillData(ord, shopSettings)

	// Execute template
	var buf bytes.Buffer
	if err := tmpl.Execute(&buf, data); err != nil {
		return "", fmt.Errorf("failed to execute HTML template: %w", err)
	}

	return buf.String(), nil
}

// prepareBillData prepares data for bill template rendering
func (s *printService) prepareBillData(ord *order.Order, shopSettings *settings.ShopSettings) map[string]interface{} {
	// Format items
	items := make([]map[string]interface{}, 0, len(ord.Items))
	for i, item := range ord.Items {
		itemName := item.Name
		if item.VariantName != "" {
			itemName = fmt.Sprintf("%s (%s)", item.Name, item.VariantName)
		}

		itemData := map[string]interface{}{
			"STT":       i + 1,
			"Name":      itemName,
			"Quantity":  item.Quantity,
			"UnitPrice": formatMoneyVN(item.Price),
			"Total":     formatMoneyVN(item.Subtotal),
		}
		
		// Add note if present (e.g., "50%", "25%")
		if item.Note != "" {
			itemData["Note"] = item.Note
		}
		
		items = append(items, itemData)
	}

	data := map[string]interface{}{
		"ShopName":      shopSettings.ShopName,
		"ShopAddress":   shopSettings.ShopAddress,
		"ShopPhone":     shopSettings.ShopPhone,
		"ShowLogo":      shopSettings.ShowLogo,
		"ShowAddress":   shopSettings.ShowAddress,
		"ShowPhone":     shopSettings.ShowPhone,
		"OrderNumber":   ord.OrderNumber,
		"CustomerName":  ord.CustomerName,
		"WaiterName":    ord.WaiterName,
		"PaymentMethod": string(ord.PaymentMethod),
		"CreatedDate":   ord.CreatedAt.Format("02/01/2006 03:04 PM"),
		"Items":         items,
		"Total":         formatMoneyVN(ord.Total),
		"CustomMessage": shopSettings.CustomMessage,
		"ShowCustomMsg": shopSettings.ShowCustomMessage,
		"LogoURL":       shopSettings.LogoURL,
	}
	
	return data
}

// CreateTempBillJob creates a temporary bill print job for an order
func (s *printService) CreateTempBillJob(ctx context.Context, ord *order.Order) error {
	if ord == nil {
		return fmt.Errorf("order cannot be nil")
	}

	log.Printf("[PRINT] Creating temporary bill job for order %s (ID: %s)", ord.OrderNumber, ord.ID.Hex())

	// Get default bill printer
	billPrinter, err := s.printerConfigRepo.FindDefault(ctx, printing.PrinterTypeBill)
	if err != nil {
		log.Printf("[PRINT ERROR] Failed to get default bill printer for order %s: %v", ord.OrderNumber, err)
		return fmt.Errorf("failed to get default bill printer: %w", err)
	}
	if billPrinter == nil {
		log.Printf("[PRINT ERROR] No default bill printer configured for order %s", ord.OrderNumber)
		return fmt.Errorf("no default bill printer configured")
	}

	// Fetch shop settings
	shopSettings, err := s.shopSettingsRepo.GetSettings(ctx)
	if err != nil {
		log.Printf("[PRINT ERROR] Failed to fetch shop settings for temp bill - order_id=%s, error=%v", ord.ID.Hex(), err)
		return fmt.Errorf("failed to fetch shop settings: %w", err)
	}

	var content string
	var contentType string

	// Render temp bill HTML
	htmlContent, err := s.renderTempBillHTML(ord, shopSettings)
	if err != nil {
		log.Printf("[PRINT ERROR] Failed to render temp bill HTML - order_id=%s, error=%v", ord.ID.Hex(), err)
		return fmt.Errorf("failed to render temp bill HTML: %w", err)
	}
	content = htmlContent
	contentType = "html"
	log.Printf("[PRINT] Using temp bill HTML template - order_id=%s", ord.OrderNumber)

	// Create print job
	job := &printing.PrintJob{
		Type:        printing.PrintJobTypeBill,
		OrderID:     ord.ID,
		OrderNumber: ord.OrderNumber,
		PrinterID:   billPrinter.ID,
		Content:     content,
		ContentType: contentType,
		Status:      printing.PrintJobStatusPending,
		RetryCount:  0,
		MaxRetries:  3,
		CreatedAt:   time.Now(),
		UpdatedAt:   time.Now(),
	}

	if err := s.printJobRepo.Create(ctx, job); err != nil {
		log.Printf("[PRINT ERROR] Failed to save temp bill print job - order_id=%s, error=%v", ord.ID.Hex(), err)
		return fmt.Errorf("failed to create print job: %w", err)
	}

	log.Printf("[PRINT] Temporary bill job created - job_id=%s, order_id=%s", job.ID.Hex(), ord.ID.Hex())
	
	// Broadcast WebSocket event
	if s.wsBroadcaster != nil {
		s.wsBroadcaster.BroadcastPrintJobCreated(job)
	}

	return nil
}

// renderTempBillHTML renders temporary bill HTML from template file
func (s *printService) renderTempBillHTML(ord *order.Order, shopSettings *settings.ShopSettings) (string, error) {
	// Read temp bill HTML template file
	templatePath := "./application/services/templates/temp_bill_template.html"
	
	// Log the absolute path for debugging
	absPath, _ := filepath.Abs(templatePath)
	log.Printf("[PRINT] Reading temp bill template from: %s", absPath)
	
	templateContent, err := os.ReadFile(templatePath)
	if err != nil {
		log.Printf("[PRINT ERROR] Failed to read temp bill template from %s: %v", absPath, err)
		return "", fmt.Errorf("failed to read temp bill HTML template: %w", err)
	}
	
	log.Printf("[PRINT] Temp bill template loaded successfully, size: %d bytes", len(templateContent))

	// Parse template
	tmpl, err := template.New("temp_bill").Parse(string(templateContent))
	if err != nil {
		return "", fmt.Errorf("failed to parse temp bill HTML template: %w", err)
	}

	// Prepare data (same as regular bill)
	data := s.prepareBillData(ord, shopSettings)

	// Execute template
	var buf bytes.Buffer
	if err := tmpl.Execute(&buf, data); err != nil {
		return "", fmt.Errorf("failed to execute temp bill HTML template: %w", err)
	}
	
	log.Printf("[PRINT] Temp bill HTML rendered successfully, size: %d bytes", buf.Len())

	return buf.String(), nil
}
