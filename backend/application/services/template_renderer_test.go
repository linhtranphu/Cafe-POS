package services

import (
	"strings"
	"testing"
	"time"

	"cafe-pos/backend/domain/order"
	"cafe-pos/backend/domain/printing"
	"cafe-pos/backend/domain/settings"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

func TestRenderBill_Success(t *testing.T) {
	renderer := NewTemplateRenderer()
	
	// Create test order
	ord := &order.Order{
		ID:          primitive.NewObjectID(),
		OrderNumber: "ORD-001",
		WaiterName:  "John Doe",
		Items: []order.OrderItem{
			{
				Name:        "Cappuccino",
				VariantName: "Size M",
				Quantity:    2,
				Price:       45000,
				Subtotal:    90000,
			},
			{
				Name:     "Latte",
				Quantity: 1,
				Price:    50000,
				Subtotal: 50000,
			},
		},
		Subtotal:  140000,
		Discount:  10000,
		Total:     130000,
		CreatedAt: time.Now(),
	}
	
	shopInfo := &settings.ShopSettings{
		ShopName:    "Coffee Shop",
		ShopAddress: "123 Main St",
		ShopPhone:   "0123456789",
	}
	
	// Render bill with default template
	content, err := renderer.RenderBill(ord, nil, shopInfo)
	if err != nil {
		t.Fatalf("RenderBill failed: %v", err)
	}
	
	// Verify all required fields are present
	requiredFields := []string{
		"Coffee Shop",
		"123 Main St",
		"0123456789",
		"ORD-001",
		"John Doe",
		"Cappuccino",
		"Size M",
		"Latte",
		"2 x 45000",
		"1 x 50000",
		"140000", // Subtotal
		"10000",  // Discount
		"130000", // Total
	}
	
	for _, field := range requiredFields {
		if !strings.Contains(content, field) {
			t.Errorf("Bill content missing required field: %s", field)
		}
	}
}

func TestRenderBill_NilOrder(t *testing.T) {
	renderer := NewTemplateRenderer()
	shopInfo := &settings.ShopSettings{Name: "Test", Address: "Test", Phone: "Test"}
	
	_, err := renderer.RenderBill(nil, nil, shopInfo)
	if err == nil {
		t.Error("Expected error for nil order, got nil")
	}
}

func TestRenderBill_NilShopInfo(t *testing.T) {
	renderer := NewTemplateRenderer()
	ord := &order.Order{OrderNumber: "ORD-001"}
	
	_, err := renderer.RenderBill(ord, nil, nil)
	if err == nil {
		t.Error("Expected error for nil shop info, got nil")
	}
}

func TestRenderBill_CustomTemplate(t *testing.T) {
	renderer := NewTemplateRenderer()
	
	ord := &order.Order{
		OrderNumber: "ORD-002",
		Total:       100000,
		CreatedAt:   time.Now(),
	}
	
	shopInfo := &settings.ShopSettings{
		ShopName:    "Custom Shop",
		ShopAddress: "456 Custom St",
		ShopPhone:   "9876543210",
	}
	
	customTemplate := &printing.PrintTemplate{
		Content: `{{.ShopName}} - Order {{.Order.OrderNumber}}
Total: {{formatPrice .Order.Total}} VND`,
	}
	
	content, err := renderer.RenderBill(ord, customTemplate, shopInfo)
	if err != nil {
		t.Fatalf("RenderBill with custom template failed: %v", err)
	}
	
	if !strings.Contains(content, "Custom Shop - Order ORD-002") {
		t.Error("Custom template not applied correctly")
	}
	if !strings.Contains(content, "Total: 100000 VND") {
		t.Error("Custom template missing total")
	}
}

func TestRenderBill_InvalidCustomTemplate_FallbackToDefault(t *testing.T) {
	renderer := NewTemplateRenderer()
	
	ord := &order.Order{
		OrderNumber: "ORD-003",
		WaiterName:  "Jane",
		Total:       50000,
		CreatedAt:   time.Now(),
	}
	
	shopInfo := &settings.ShopSettings{
		ShopName:    "Test Shop",
		ShopAddress: "Test Address",
		ShopPhone:   "1234567890",
	}
	
	// Invalid template with syntax error
	invalidTemplate := &printing.PrintTemplate{
		Content: `{{.ShopName}} {{.InvalidField`,
	}
	
	// Should fallback to default template
	content, err := renderer.RenderBill(ord, invalidTemplate, shopInfo)
	if err != nil {
		t.Fatalf("Expected fallback to succeed, got error: %v", err)
	}
	
	// Should contain fields from default template
	if !strings.Contains(content, "Test Shop") {
		t.Error("Fallback template missing shop name")
	}
	if !strings.Contains(content, "ORD-003") {
		t.Error("Fallback template missing order number")
	}
}

func TestRenderLabel_Success(t *testing.T) {
	renderer := NewTemplateRenderer()
	
	ord := &order.Order{
		ID:          primitive.NewObjectID(),
		OrderNumber: "ORD-004",
		Items: []order.OrderItem{
			{
				Name:        "Espresso",
				VariantName: "Double Shot",
				Note:        "Extra hot",
				Quantity:    1,
				Price:       40000,
			},
			{
				Name:     "Americano",
				Quantity: 1,
				Price:    35000,
			},
		},
		CreatedAt: time.Now(),
	}
	
	shopInfo := &settings.ShopSettings{
		ShopName:    "Coffee Shop",
		ShopAddress: "123 Main St",
		ShopPhone:   "0123456789",
	}
	
	// Render label for first item
	content, err := renderer.RenderLabel(ord, 0, nil, shopInfo)
	if err != nil {
		t.Fatalf("RenderLabel failed: %v", err)
	}
	
	// Verify required fields
	requiredFields := []string{
		"ORD-004",
		"1/2", // Item 1 of 2
		"Espresso",
		"Double Shot",
		"Extra hot", // Note without "Note:" prefix in new template
	}
	
	for _, field := range requiredFields {
		if !strings.Contains(content, field) {
			t.Errorf("Label content missing required field: %s", field)
		}
	}
}

func TestRenderLabel_SecondItem(t *testing.T) {
	renderer := NewTemplateRenderer()
	
	ord := &order.Order{
		OrderNumber: "ORD-005",
		Items: []order.OrderItem{
			{Name: "Item 1"},
			{Name: "Item 2", VariantName: "Large"},
			{Name: "Item 3"},
		},
		CreatedAt: time.Now(),
	}
	
	shopInfo := &settings.ShopSettings{Name: "Shop", Address: "Addr", Phone: "Phone"}
	
	// Render label for second item (index 1)
	content, err := renderer.RenderLabel(ord, 1, nil, shopInfo)
	if err != nil {
		t.Fatalf("RenderLabel failed: %v", err)
	}
	
	if !strings.Contains(content, "2/3") {
		t.Error("Label missing correct item position")
	}
	if !strings.Contains(content, "Item 2") {
		t.Error("Label missing item name")
	}
	if !strings.Contains(content, "Large") {
		t.Error("Label missing variant name")
	}
}

func TestRenderLabel_InvalidItemIndex(t *testing.T) {
	renderer := NewTemplateRenderer()
	
	ord := &order.Order{
		OrderNumber: "ORD-006",
		Items:       []order.OrderItem{{Name: "Item 1"}},
	}
	
	shopInfo := &settings.ShopSettings{Name: "Shop", Address: "Addr", Phone: "Phone"}
	
	// Test negative index
	_, err := renderer.RenderLabel(ord, -1, nil, shopInfo)
	if err == nil {
		t.Error("Expected error for negative index, got nil")
	}
	
	// Test index out of bounds
	_, err = renderer.RenderLabel(ord, 5, nil, shopInfo)
	if err == nil {
		t.Error("Expected error for out of bounds index, got nil")
	}
}

func TestRenderLabel_NilOrder(t *testing.T) {
	renderer := NewTemplateRenderer()
	shopInfo := &settings.ShopSettings{Name: "Test", Address: "Test", Phone: "Test"}
	
	_, err := renderer.RenderLabel(nil, 0, nil, shopInfo)
	if err == nil {
		t.Error("Expected error for nil order, got nil")
	}
}

func TestRenderLabel_NilShopInfo(t *testing.T) {
	renderer := NewTemplateRenderer()
	ord := &order.Order{
		OrderNumber: "ORD-007",
		Items:       []order.OrderItem{{Name: "Item"}},
	}
	
	_, err := renderer.RenderLabel(ord, 0, nil, nil)
	if err == nil {
		t.Error("Expected error for nil shop info, got nil")
	}
}

func TestRenderLabel_CustomTemplate(t *testing.T) {
	renderer := NewTemplateRenderer()
	
	ord := &order.Order{
		OrderNumber: "ORD-008",
		Items: []order.OrderItem{
			{Name: "Custom Item", VariantName: "Special"},
		},
		CreatedAt: time.Now(),
	}
	
	shopInfo := &settings.ShopSettings{Name: "Shop", Address: "Addr", Phone: "Phone"}
	
	customTemplate := &printing.PrintTemplate{
		Content: `Order: {{.Order.OrderNumber}}
{{with index .Order.Items .ItemIndex}}{{.Name}} - {{.VariantName}}{{end}}`,
	}
	
	content, err := renderer.RenderLabel(ord, 0, customTemplate, shopInfo)
	if err != nil {
		t.Fatalf("RenderLabel with custom template failed: %v", err)
	}
	
	if !strings.Contains(content, "Custom Item - Special") {
		t.Error("Custom template not applied correctly")
	}
}

func TestRenderLabel_InvalidCustomTemplate_FallbackToDefault(t *testing.T) {
	renderer := NewTemplateRenderer()
	
	ord := &order.Order{
		OrderNumber: "ORD-009",
		Items: []order.OrderItem{
			{Name: "Test Item"},
		},
		CreatedAt: time.Now(),
	}
	
	shopInfo := &settings.ShopSettings{Name: "Shop", Address: "Addr", Phone: "Phone"}
	
	// Invalid template
	invalidTemplate := &printing.PrintTemplate{
		Content: `{{.InvalidSyntax`,
	}
	
	// Should fallback to default
	content, err := renderer.RenderLabel(ord, 0, invalidTemplate, shopInfo)
	if err != nil {
		t.Fatalf("Expected fallback to succeed, got error: %v", err)
	}
	
	if !strings.Contains(content, "ORD-009") {
		t.Error("Fallback template missing order number")
	}
	if !strings.Contains(content, "Test Item") {
		t.Error("Fallback template missing item name")
	}
}
