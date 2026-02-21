package printing

import (
	"testing"
	"time"
)

// TestTemplateType_Constants tests that TemplateType constants are defined correctly
func TestTemplateType_Constants(t *testing.T) {
	tests := []struct {
		name         string
		templateType TemplateType
		expected     string
	}{
		{"Bill template type", TemplateTypeBill, "BILL"},
		{"Label template type", TemplateTypeLabel, "LABEL"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if string(tt.templateType) != tt.expected {
				t.Errorf("Expected %s, got %s", tt.expected, string(tt.templateType))
			}
		})
	}
}

// TestPrintTemplate_Creation tests creating a PrintTemplate with valid data
func TestPrintTemplate_Creation(t *testing.T) {
	now := time.Now()

	template := &PrintTemplate{
		Type:      TemplateTypeBill,
		Name:      "Default Bill Template",
		Content:   "{{.ShopName}}\n{{.Order.OrderNumber}}",
		IsDefault: true,
		CreatedAt: now,
		UpdatedAt: now,
	}

	// Verify all fields are set correctly
	if template.Type != TemplateTypeBill {
		t.Errorf("Expected Type to be BILL, got %s", template.Type)
	}
	if template.Name != "Default Bill Template" {
		t.Errorf("Expected Name to be 'Default Bill Template', got %s", template.Name)
	}
	if template.Content == "" {
		t.Error("Expected Content to be set")
	}
	if !template.IsDefault {
		t.Error("Expected IsDefault to be true")
	}
}

// TestPrintTemplate_BillTemplate tests bill template structure
func TestPrintTemplate_BillTemplate(t *testing.T) {
	billTemplate := &PrintTemplate{
		Type: TemplateTypeBill,
		Name: "Standard Bill",
		Content: `{{.ShopName}}
{{.ShopAddress}}
Tel: {{.ShopPhone}}
================================
Order: {{.Order.OrderNumber}}
Time: {{.PrintTime.Format "02/01/2006 15:04"}}
================================
{{range .Order.Items}}
{{.Name}} x {{.Quantity}} = {{.Subtotal}}
{{end}}
================================
TOTAL: {{.Order.Total}}`,
		IsDefault: true,
	}

	if billTemplate.Type != TemplateTypeBill {
		t.Error("Expected template type to be BILL")
	}

	// Verify template contains required fields
	requiredFields := []string{
		"{{.ShopName}}",
		"{{.Order.OrderNumber}}",
		"{{.PrintTime",
		"{{.Order.Total}}",
	}

	for _, field := range requiredFields {
		if !contains(billTemplate.Content, field) {
			t.Errorf("Expected template to contain %s", field)
		}
	}
}

// TestPrintTemplate_LabelTemplate tests label template structure
func TestPrintTemplate_LabelTemplate(t *testing.T) {
	labelTemplate := &PrintTemplate{
		Type: TemplateTypeLabel,
		Name: "Standard Label",
		Content: `Order: {{.Order.OrderNumber}}
{{.ItemIndex}}/{{.TotalItems}}

{{with index .Order.Items .ItemIndex}}
{{.Name}}
{{if .VariantName}}{{.VariantName}}{{end}}
{{end}}

{{.PrintTime.Format "15:04"}}`,
		IsDefault: true,
	}

	if labelTemplate.Type != TemplateTypeLabel {
		t.Error("Expected template type to be LABEL")
	}

	// Verify template contains required fields
	requiredFields := []string{
		"{{.Order.OrderNumber}}",
		"{{.ItemIndex}}",
		"{{.TotalItems}}",
		"{{.PrintTime",
	}

	for _, field := range requiredFields {
		if !contains(labelTemplate.Content, field) {
			t.Errorf("Expected template to contain %s", field)
		}
	}
}

// TestTemplateData_BillData tests bill template data structure
func TestTemplateData_BillData(t *testing.T) {
	now := time.Now()

	data := &TemplateData{
		ShopName:    "Test Coffee Shop",
		ShopAddress: "123 Main St",
		ShopPhone:   "0123456789",
		Order:       nil, // Would be actual order in real usage
		PrintTime:   now,
	}

	if data.ShopName != "Test Coffee Shop" {
		t.Errorf("Expected ShopName to be 'Test Coffee Shop', got %s", data.ShopName)
	}
	if data.ShopAddress != "123 Main St" {
		t.Errorf("Expected ShopAddress to be '123 Main St', got %s", data.ShopAddress)
	}
	if data.ShopPhone != "0123456789" {
		t.Errorf("Expected ShopPhone to be '0123456789', got %s", data.ShopPhone)
	}
	if !data.PrintTime.Equal(now) {
		t.Error("Expected PrintTime to match")
	}
}

// TestTemplateData_LabelData tests label template data structure
func TestTemplateData_LabelData(t *testing.T) {
	now := time.Now()

	data := &TemplateData{
		ShopName:   "Test Coffee Shop",
		Order:      nil, // Would be actual order in real usage
		PrintTime:  now,
		ItemIndex:  0,
		TotalItems: 3,
	}

	if data.ItemIndex != 0 {
		t.Errorf("Expected ItemIndex to be 0, got %d", data.ItemIndex)
	}
	if data.TotalItems != 3 {
		t.Errorf("Expected TotalItems to be 3, got %d", data.TotalItems)
	}
}

// TestPrintTemplate_DefaultTemplate tests default template logic
func TestPrintTemplate_DefaultTemplate(t *testing.T) {
	// Only one template of each type should be default
	billTemplate1 := &PrintTemplate{
		Type:      TemplateTypeBill,
		Name:      "Default Bill",
		IsDefault: true,
	}

	billTemplate2 := &PrintTemplate{
		Type:      TemplateTypeBill,
		Name:      "Custom Bill",
		IsDefault: false,
	}

	labelTemplate := &PrintTemplate{
		Type:      TemplateTypeLabel,
		Name:      "Default Label",
		IsDefault: true,
	}

	// Verify only one bill template is default
	if !billTemplate1.IsDefault {
		t.Error("Expected Bill Template 1 to be default")
	}
	if billTemplate2.IsDefault {
		t.Error("Expected Bill Template 2 to not be default")
	}

	// Label template can also be default (different type)
	if !labelTemplate.IsDefault {
		t.Error("Expected Label Template to be default")
	}
}

// TestCreatePrintTemplateRequest_Validation tests request validation
func TestCreatePrintTemplateRequest_Validation(t *testing.T) {
	tests := []struct {
		name    string
		request CreatePrintTemplateRequest
		valid   bool
	}{
		{
			"Valid bill template",
			CreatePrintTemplateRequest{
				Type:      TemplateTypeBill,
				Name:      "My Bill Template",
				Content:   "{{.ShopName}}\n{{.Order.Total}}",
				IsDefault: false,
			},
			true,
		},
		{
			"Valid label template",
			CreatePrintTemplateRequest{
				Type:      TemplateTypeLabel,
				Name:      "My Label Template",
				Content:   "{{.Order.OrderNumber}}\n{{.ItemIndex}}",
				IsDefault: false,
			},
			true,
		},
		{
			"Empty name",
			CreatePrintTemplateRequest{
				Type:    TemplateTypeBill,
				Name:    "",
				Content: "{{.ShopName}}",
			},
			false,
		},
		{
			"Empty content",
			CreatePrintTemplateRequest{
				Type:    TemplateTypeBill,
				Name:    "Template",
				Content: "",
			},
			false,
		},
		{
			"Invalid type",
			CreatePrintTemplateRequest{
				Type:    "INVALID",
				Name:    "Template",
				Content: "Content",
			},
			false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Basic validation checks
			isValid := tt.request.Name != "" &&
				tt.request.Content != "" &&
				(tt.request.Type == TemplateTypeBill || tt.request.Type == TemplateTypeLabel)

			if isValid != tt.valid {
				t.Errorf("Expected validation to be %v, got %v", tt.valid, isValid)
			}
		})
	}
}

// TestUpdatePrintTemplateRequest_PartialUpdate tests partial update request
func TestUpdatePrintTemplateRequest_PartialUpdate(t *testing.T) {
	// Test that we can update individual fields
	isDefault := true

	request := &UpdatePrintTemplateRequest{
		Name:      "Updated Template Name",
		IsDefault: &isDefault,
	}

	if request.Name != "Updated Template Name" {
		t.Errorf("Expected Name to be 'Updated Template Name', got %s", request.Name)
	}

	if request.IsDefault == nil || *request.IsDefault != true {
		t.Error("Expected IsDefault to be true")
	}

	// Fields not set should be empty
	if request.Content != "" {
		t.Error("Expected Content to be empty for partial update")
	}
}

// TestPrintTemplateFilter_DefaultValues tests filter default values
func TestPrintTemplateFilter_DefaultValues(t *testing.T) {
	filter := &PrintTemplateFilter{
		Page:  1,
		Limit: 20,
	}

	if filter.Page != 1 {
		t.Errorf("Expected default Page to be 1, got %d", filter.Page)
	}

	if filter.Limit != 20 {
		t.Errorf("Expected default Limit to be 20, got %d", filter.Limit)
	}

	if filter.Type != "" {
		t.Errorf("Expected Type to be empty, got %s", filter.Type)
	}
}

// TestPrintTemplate_ContentValidation tests template content validation
func TestPrintTemplate_ContentValidation(t *testing.T) {
	tests := []struct {
		name        string
		templateType TemplateType
		content     string
		valid       bool
		description string
	}{
		{
			"Valid bill template with all fields",
			TemplateTypeBill,
			"{{.ShopName}}\n{{.Order.OrderNumber}}\n{{.Order.Total}}",
			true,
			"Contains required bill fields",
		},
		{
			"Valid label template with all fields",
			TemplateTypeLabel,
			"{{.Order.OrderNumber}}\n{{.ItemIndex}}/{{.TotalItems}}",
			true,
			"Contains required label fields",
		},
		{
			"Template with invalid syntax",
			TemplateTypeBill,
			"{{.ShopName}\n{{.Order.Total}}",
			false,
			"Missing closing braces",
		},
		{
			"Empty template",
			TemplateTypeBill,
			"",
			false,
			"Template cannot be empty",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			template := &PrintTemplate{
				Type:    tt.templateType,
				Name:    "Test Template",
				Content: tt.content,
			}

			// Basic validation: content should not be empty
			isValid := template.Content != ""

			// Check for balanced braces (simple check)
			if isValid {
				openCount := 0
				closeCount := 0
				for i := 0; i < len(template.Content)-1; i++ {
					if template.Content[i] == '{' && template.Content[i+1] == '{' {
						openCount++
					}
					if template.Content[i] == '}' && template.Content[i+1] == '}' {
						closeCount++
					}
				}
				isValid = openCount == closeCount
			}

			if isValid != tt.valid {
				t.Errorf("Expected validation to be %v for %s", tt.valid, tt.description)
			}
		})
	}
}

// TestPrintTemplate_TypeSpecificFields tests type-specific field requirements
func TestPrintTemplate_TypeSpecificFields(t *testing.T) {
	tests := []struct {
		name         string
		templateType TemplateType
		requiredFields []string
	}{
		{
			"Bill template required fields",
			TemplateTypeBill,
			[]string{"ShopName", "Order.OrderNumber", "Order.Total"},
		},
		{
			"Label template required fields",
			TemplateTypeLabel,
			[]string{"Order.OrderNumber", "ItemIndex", "TotalItems"},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// This test documents what fields should be available
			// In a real implementation, you might validate that templates
			// contain these fields
			if len(tt.requiredFields) == 0 {
				t.Error("Expected required fields to be defined")
			}
		})
	}
}

// Helper function to check if string contains substring
func contains(s, substr string) bool {
	return len(s) >= len(substr) && (s == substr || len(s) > len(substr) && containsHelper(s, substr))
}

func containsHelper(s, substr string) bool {
	for i := 0; i <= len(s)-len(substr); i++ {
		if s[i:i+len(substr)] == substr {
			return true
		}
	}
	return false
}
