package services

import (
	"bytes"
	"fmt"
	"log"
	"strings"
	"text/template"
	"time"

	"cafe-pos/backend/domain/order"
	"cafe-pos/backend/domain/printing"
	"cafe-pos/backend/domain/settings"
	printingInfra "cafe-pos/backend/infrastructure/printing"
)

// TemplateRenderer defines the interface for rendering print templates
type TemplateRenderer interface {
	RenderBill(order *order.Order, tmpl *printing.PrintTemplate, shopSettings *settings.ShopSettings) (string, error)
	RenderLabel(order *order.Order, itemIndex int, tmpl *printing.PrintTemplate, shopSettings *settings.ShopSettings) (string, error)
}

// templateRenderer implements the TemplateRenderer interface
type templateRenderer struct {
	defaultBillTemplate  string
	defaultLabelTemplate string
	tableFormatter       *printingInfra.TableFormatter
	formatParser         *printingInfra.FormatParser
}

// NewTemplateRenderer creates a new template renderer with default templates
func NewTemplateRenderer() TemplateRenderer {
	// Default paper width for 80mm paper (576 pixels)
	paperWidth := 576
	margin := 10
	
	// Initialize table formatter
	tableFormatter := printingInfra.NewTableFormatter(paperWidth, margin, 1)
	
	// Initialize format parser
	formatParser := printingInfra.NewFormatParser(paperWidth)
	
	return &templateRenderer{
		defaultBillTemplate:  getDefaultBillTemplate(),
		defaultLabelTemplate: getDefaultLabelTemplate(),
		tableFormatter:       tableFormatter,
		formatParser:         formatParser,
	}
}

// RenderBill renders a bill template with order data
func (r *templateRenderer) RenderBill(ord *order.Order, tmpl *printing.PrintTemplate, shopSettings *settings.ShopSettings) (string, error) {
	if ord == nil {
		return "", fmt.Errorf("order cannot be nil")
	}
	if shopSettings == nil {
		return "", fmt.Errorf("shop settings cannot be nil")
	}

	// Prepare template data with shop settings
	data := printing.TemplateData{
		ShopName:          shopSettings.ShopName,
		ShopAddress:       shopSettings.ShopAddress,
		ShopPhone:         shopSettings.ShopPhone,
		LogoURL:           shopSettings.LogoURL,
		CustomMessage:     shopSettings.CustomMessage,
		ShowLogo:          shopSettings.ShowLogo,
		ShowAddress:       shopSettings.ShowAddress,
		ShowPhone:         shopSettings.ShowPhone,
		ShowCustomMessage: shopSettings.ShowCustomMessage,
		Order:             ord,
		PrintTime:         time.Now(),
	}

	// Use custom template if provided, otherwise use default
	templateContent := r.defaultBillTemplate
	if tmpl != nil && tmpl.Content != "" {
		templateContent = tmpl.Content
	}

	// Try to render with the provided/default template
	content, err := r.renderTemplate("bill", templateContent, data)
	if err != nil {
		// If custom template fails, fallback to default
		if tmpl != nil && tmpl.Content != "" {
			content, fallbackErr := r.renderTemplate("bill_fallback", r.defaultBillTemplate, data)
			if fallbackErr != nil {
				return "", fmt.Errorf("template rendering failed and fallback failed: %w", err)
			}
			log.Printf("Template rendering failed: %v. Using fallback template.", err)
			// Still process the fallback content
			processedContent, processErr := r.processTemplateContent(content, ord, shopSettings)
			if processErr != nil {
				log.Printf("Template processing failed: %v. Using unprocessed content.", processErr)
				return content, nil
			}
			return processedContent, nil
		}
		return "", fmt.Errorf("template rendering failed: %w", err)
	}

	// Process the rendered content with new modules
	processedContent, err := r.processTemplateContent(content, ord, shopSettings)
	if err != nil {
		// If processing fails, log error and return original content as fallback
		log.Printf("Template processing failed: %v. Using unprocessed content.", err)
		return content, nil
	}

	return processedContent, nil
}

// processTemplateContent processes the template content (currently just removes markers)
func (r *templateRenderer) processTemplateContent(content string, ord *order.Order, shopSettings *settings.ShopSettings) (string, error) {
	// Remove [LOGO] marker if present (logo rendering not implemented)
	content = strings.ReplaceAll(content, "[LOGO]", "")
	
	// Handle [TABLE_START] and [TABLE_END] markers
	content = r.processTableMarkers(content, ord)

	return content, nil
}

// processTableMarkers processes [TABLE_START] and [TABLE_END] markers and formats the table
func (r *templateRenderer) processTableMarkers(content string, ord *order.Order) string {
	// Find table markers
	tableStartIdx := strings.Index(content, "[TABLE_START]")
	tableEndIdx := strings.Index(content, "[TABLE_END]")

	// If no table markers found, return content as-is
	if tableStartIdx == -1 || tableEndIdx == -1 || tableEndIdx <= tableStartIdx {
		return content
	}

	// Extract content before, inside, and after table markers
	beforeTable := content[:tableStartIdx]
	afterTable := content[tableEndIdx+len("[TABLE_END]"):]

	// Format items table using table formatter
	items := make([]printingInfra.OrderItem, len(ord.Items))
	for i, item := range ord.Items {
		items[i] = printingInfra.OrderItem{
			Name:        item.Name,
			VariantName: item.VariantName,
			Quantity:    item.Quantity,
			UnitPrice:   item.Price,
			TotalPrice:  item.Subtotal,
		}
	}

	paperWidth := 576 // Default 80mm paper width
	tableLines := r.tableFormatter.FormatItemsTable(items, paperWidth)

	// Join table lines
	tableContent := strings.Join(tableLines, "\n")

	// Reconstruct content with formatted table
	return beforeTable + "\n" + tableContent + "\n" + afterTable
}

// RenderLabel renders a label template for a specific order item
func (r *templateRenderer) RenderLabel(ord *order.Order, itemIndex int, tmpl *printing.PrintTemplate, shopSettings *settings.ShopSettings) (string, error) {
	if ord == nil {
		return "", fmt.Errorf("order cannot be nil")
	}
	if shopSettings == nil {
		return "", fmt.Errorf("shop settings cannot be nil")
	}
	if itemIndex < 0 || itemIndex >= len(ord.Items) {
		return "", fmt.Errorf("invalid item index: %d (order has %d items)", itemIndex, len(ord.Items))
	}

	// Prepare template data with shop settings
	data := printing.TemplateData{
		ShopName:          shopSettings.ShopName,
		ShopAddress:       shopSettings.ShopAddress,
		ShopPhone:         shopSettings.ShopPhone,
		LogoURL:           shopSettings.LogoURL,
		CustomMessage:     shopSettings.CustomMessage,
		ShowLogo:          shopSettings.ShowLogo,
		ShowAddress:       shopSettings.ShowAddress,
		ShowPhone:         shopSettings.ShowPhone,
		ShowCustomMessage: shopSettings.ShowCustomMessage,
		Order:             ord,
		PrintTime:         time.Now(),
		ItemIndex:         itemIndex,
		TotalItems:        len(ord.Items),
	}

	// Use custom template if provided, otherwise use default
	templateContent := r.defaultLabelTemplate
	if tmpl != nil && tmpl.Content != "" {
		templateContent = tmpl.Content
	}

	// Try to render with the provided/default template
	content, err := r.renderTemplate("label", templateContent, data)
	if err != nil {
		// If custom template fails, fallback to default
		if tmpl != nil && tmpl.Content != "" {
			content, fallbackErr := r.renderTemplate("label_fallback", r.defaultLabelTemplate, data)
			if fallbackErr != nil {
				return "", fmt.Errorf("template rendering failed and fallback failed: %w", err)
			}
			return content, nil
		}
		return "", fmt.Errorf("template rendering failed: %w", err)
	}

	return content, nil
}

// renderTemplate renders a template with the given data
func (r *templateRenderer) renderTemplate(name string, templateContent string, data printing.TemplateData) (string, error) {
	// Create template with custom functions
	tmpl, err := template.New(name).Funcs(template.FuncMap{
		"formatPrice": func(price float64) string {
			return fmt.Sprintf("%.0f", price)
		},
		"formatTime": func(t time.Time, layout string) string {
			return t.Format(layout)
		},
		"add": func(a, b int) int {
			return a + b
		},
		"truncate": func(s string, maxLen int) string {
			if len(s) <= maxLen {
				return s
			}
			if maxLen <= 3 {
				return s[:maxLen]
			}
			return s[:maxLen-3] + "..."
		},
		"padRight": func(s string, width int) string {
			if len(s) >= width {
				return s
			}
			return s + strings.Repeat(" ", width-len(s))
		},
	}).Parse(templateContent)
	if err != nil {
		return "", fmt.Errorf("failed to parse template: %w", err)
	}

	// Execute template
	var buf bytes.Buffer
	if err := tmpl.Execute(&buf, data); err != nil {
		return "", fmt.Errorf("failed to execute template: %w", err)
	}

	return buf.String(), nil
}

// getDefaultBillTemplate returns the default bill template
// Designed to fit within 80mm (48 chars) paper width
func getDefaultBillTemplate() string {
	return `{{truncate .ShopName 48}}
{{if .ShowAddress}}{{if .ShopAddress}}{{truncate .ShopAddress 48}}
{{end}}{{end}}{{if .ShowPhone}}{{if .ShopPhone}}Tel: {{.ShopPhone}}
{{end}}{{end}}================================
Order: {{.Order.OrderNumber}}
Time: {{formatTime .Order.CreatedAt "02/01 15:04"}}
{{if .Order.WaiterName}}Waiter: {{truncate .Order.WaiterName 40}}
{{end}}================================
{{range .Order.Items}}{{truncate .Name 48}}
{{if .VariantName}}  {{truncate .VariantName 46}}
{{end}}  {{.Quantity}} x {{formatPrice .Price}} = {{formatPrice .Subtotal}}
{{end}}================================
Subtotal: {{padRight (formatPrice .Order.Subtotal) 20}}
Discount: {{padRight (formatPrice .Order.Discount) 20}}
--------------------------------
TOTAL: {{formatPrice .Order.Total}} VND
================================
{{if .ShowCustomMessage}}{{if .CustomMessage}}{{truncate .CustomMessage 48}}
{{end}}{{end}}Thank you!
`
}

// getDefaultLabelTemplate returns the default label template
// Designed to fit within 60x40mm label (30 chars width, 8 lines)
// For smaller labels (40x30mm, 50x30mm), customize template or truncate content
func getDefaultLabelTemplate() string {
	return `{{truncate .Order.OrderNumber 30}} {{add .ItemIndex 1}}/{{.TotalItems}}
{{with index .Order.Items .ItemIndex}}{{truncate .Name 30}}
{{if .VariantName}}{{truncate .VariantName 30}}
{{end}}{{if .Note}}{{truncate .Note 30}}
{{end}}{{end}}{{formatTime .PrintTime "15:04"}}
`
}
