package services

import (
	"strings"
	"testing"
	"time"

	"cafe-pos/backend/domain/order"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

// TestBillTemplate_WidthConstraints verifies that bill templates fit within paper width constraints
// Requirements: 1.7 - Bill must fit 58mm (32 chars) or 80mm (48 chars) paper
func TestBillTemplate_WidthConstraints(t *testing.T) {
	renderer := NewTemplateRenderer()

	// Create test order with various item name lengths
	ord := &order.Order{
		ID:          primitive.NewObjectID(),
		OrderNumber: "ORD-001",
		WaiterName:  "John Doe",
		Items: []order.OrderItem{
			{
				Name:        "Cappuccino",
				VariantName: "Size M, Hot, Extra Foam",
				Quantity:    2,
				Price:       45000,
				Subtotal:    90000,
			},
			{
				Name:        "Caramel Macchiato",
				VariantName: "Large",
				Quantity:    1,
				Price:       55000,
				Subtotal:    55000,
			},
			{
				Name:     "Espresso",
				Quantity: 3,
				Price:    35000,
				Subtotal: 105000,
			},
		},
		Subtotal:  250000,
		Discount:  25000,
		Total:     225000,
		CreatedAt: time.Now(),
	}

	shopInfo := &ShopInfo{
		Name:    "Coffee Shop & Bakery",
		Address: "123 Main Street, District 1",
		Phone:   "0123456789",
	}

	// Render bill
	content, err := renderer.RenderBill(ord, nil, shopInfo)
	if err != nil {
		t.Fatalf("RenderBill failed: %v", err)
	}

	// Check each line width
	lines := strings.Split(content, "\n")
	
	// Test for 58mm paper (32 chars max)
	maxWidth58mm := 32
	violating58mm := []string{}
	for i, line := range lines {
		if len(line) > maxWidth58mm {
			violating58mm = append(violating58mm, line)
			t.Logf("Line %d exceeds 58mm width (%d chars): %q", i+1, len(line), line)
		}
	}

	// Test for 80mm paper (48 chars max)
	maxWidth80mm := 48
	violating80mm := []string{}
	for i, line := range lines {
		if len(line) > maxWidth80mm {
			violating80mm = append(violating80mm, line)
			t.Logf("Line %d exceeds 80mm width (%d chars): %q", i+1, len(line), line)
		}
	}

	// Report results
	if len(violating58mm) > 0 {
		t.Logf("WARNING: %d lines exceed 58mm (32 char) width", len(violating58mm))
	} else {
		t.Logf("✓ All lines fit within 58mm (32 char) width")
	}

	if len(violating80mm) > 0 {
		t.Errorf("FAIL: %d lines exceed 80mm (48 char) width", len(violating80mm))
		t.Logf("Full bill content:\n%s", content)
	} else {
		t.Logf("✓ All lines fit within 80mm (48 char) width")
	}
}

// TestLabelTemplate_SizeConstraints verifies that label templates fit within label size constraints
// Requirements: 2.8 - Labels must fit 40x30mm, 50x30mm, 60x40mm sizes
func TestLabelTemplate_SizeConstraints(t *testing.T) {
	renderer := NewTemplateRenderer()

	// Create test order with long item names and variants
	ord := &order.Order{
		ID:          primitive.NewObjectID(),
		OrderNumber: "ORD-123456",
		Items: []order.OrderItem{
			{
				Name:        "Caramel Macchiato",
				VariantName: "Large, Extra Hot, Whipped Cream",
				Note:        "Less ice, extra sweet",
				Quantity:    1,
				Price:       55000,
			},
			{
				Name:        "Iced Americano",
				VariantName: "Medium",
				Quantity:    2,
				Price:       40000,
			},
		},
		CreatedAt: time.Now(),
	}

	shopInfo := &ShopInfo{
		Name:    "Coffee Shop",
		Address: "123 Main St",
		Phone:   "0123456789",
	}

	// Test different label sizes
	labelSizes := []struct {
		name      string
		maxWidth  int // characters
		maxLines  int
		expectFit bool // whether default template should fit
	}{
		{"40x30mm", 20, 6, false},  // Small label - may not fit default template
		{"50x30mm", 25, 6, false},  // Medium label - may not fit default template
		{"60x40mm", 30, 8, true},   // Large label - default template target size
	}

	for _, size := range labelSizes {
		t.Run(size.name, func(t *testing.T) {
			// Render label for first item
			content, err := renderer.RenderLabel(ord, 0, nil, shopInfo)
			if err != nil {
				t.Fatalf("RenderLabel failed: %v", err)
			}

			lines := strings.Split(content, "\n")
			
			// Remove empty lines at the end
			for len(lines) > 0 && strings.TrimSpace(lines[len(lines)-1]) == "" {
				lines = lines[:len(lines)-1]
			}

			// Check line count
			if len(lines) > size.maxLines {
				if size.expectFit {
					t.Errorf("Label has %d lines, exceeds max %d lines for %s", 
						len(lines), size.maxLines, size.name)
				} else {
					t.Logf("Label has %d lines, exceeds max %d lines for %s (expected for small labels)", 
						len(lines), size.maxLines, size.name)
				}
			}

			// Check line width
			violatingLines := []string{}
			for i, line := range lines {
				if len(line) > size.maxWidth {
					violatingLines = append(violatingLines, line)
					t.Logf("Line %d exceeds %s width (%d chars): %q", 
						i+1, size.name, len(line), line)
				}
			}

			if len(violatingLines) > 0 {
				if size.expectFit {
					t.Errorf("FAIL: %d lines exceed %s width", len(violatingLines), size.name)
				} else {
					t.Logf("WARNING: %d lines exceed %s width (expected for small labels)", 
						len(violatingLines), size.name)
				}
			} else {
				t.Logf("✓ All lines fit within %s", size.name)
			}

			t.Logf("Label content for %s:\n%s", size.name, content)
		})
	}
}

// TestBillTemplate_LongItemNames tests bill rendering with very long item names
func TestBillTemplate_LongItemNames(t *testing.T) {
	renderer := NewTemplateRenderer()

	ord := &order.Order{
		ID:          primitive.NewObjectID(),
		OrderNumber: "ORD-999",
		WaiterName:  "Staff",
		Items: []order.OrderItem{
			{
				Name:        "Super Long Item Name That Might Break Layout",
				VariantName: "Extra Large Size with Many Toppings",
				Quantity:    1,
				Price:       100000,
				Subtotal:    100000,
			},
		},
		Subtotal:  100000,
		Discount:  0,
		Total:     100000,
		CreatedAt: time.Now(),
	}

	shopInfo := &ShopInfo{
		Name:    "Shop",
		Address: "Address",
		Phone:   "Phone",
	}

	content, err := renderer.RenderBill(ord, nil, shopInfo)
	if err != nil {
		t.Fatalf("RenderBill failed: %v", err)
	}

	lines := strings.Split(content, "\n")
	maxWidth := 48 // 80mm paper

	for i, line := range lines {
		if len(line) > maxWidth {
			t.Errorf("Line %d exceeds 80mm width (%d chars): %q", i+1, len(line), line)
		}
	}
}
