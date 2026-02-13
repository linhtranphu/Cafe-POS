package order

import (
	"encoding/json"
	"strings"
	"testing"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

// Test OrderItem with variant (new functionality)

func TestOrderItem_WithVariant(t *testing.T) {
	menuItemID := primitive.NewObjectID()
	
	item := OrderItem{
		MenuItemID:  menuItemID,
		VariantID:   "L",
		Name:        "Cà phê sữa đá",
		VariantName: "Size L",
		Price:       30000,
		Quantity:    2,
		Subtotal:    60000,
	}
	
	if item.VariantID != "L" {
		t.Errorf("Expected variant_id 'L', got '%s'", item.VariantID)
	}
	if item.VariantName != "Size L" {
		t.Errorf("Expected variant_name 'Size L', got '%s'", item.VariantName)
	}
	if item.Price != 30000 {
		t.Errorf("Expected price 30000, got %f", item.Price)
	}
}

// Test OrderItem without variant (backward compatible)

func TestOrderItem_WithoutVariant(t *testing.T) {
	menuItemID := primitive.NewObjectID()
	
	item := OrderItem{
		MenuItemID: menuItemID,
		Name:       "Bánh mì",
		Price:      20000,
		Quantity:   1,
		Subtotal:   20000,
	}
	
	if item.VariantID != "" {
		t.Errorf("Expected empty variant_id, got '%s'", item.VariantID)
	}
	if item.VariantName != "" {
		t.Errorf("Expected empty variant_name, got '%s'", item.VariantName)
	}
	if item.Price != 20000 {
		t.Errorf("Expected price 20000, got %f", item.Price)
	}
}

// Test JSON marshaling/unmarshaling

func TestOrderItem_JSON_WithVariant(t *testing.T) {
	menuItemID := primitive.NewObjectID()
	
	item := OrderItem{
		MenuItemID:  menuItemID,
		VariantID:   "L",
		Name:        "Cà phê sữa đá",
		VariantName: "Size L",
		Price:       30000,
		Quantity:    2,
		Subtotal:    60000,
	}
	
	// Marshal to JSON
	jsonData, err := json.Marshal(item)
	if err != nil {
		t.Fatalf("Failed to marshal JSON: %v", err)
	}
	
	jsonStr := string(jsonData)
	
	// Verify variant fields are present
	if !strings.Contains(jsonStr, "variant_id") {
		t.Error("Expected variant_id in JSON")
	}
	if !strings.Contains(jsonStr, "variant_name") {
		t.Error("Expected variant_name in JSON")
	}
	
	// Unmarshal back
	var unmarshaled OrderItem
	err = json.Unmarshal(jsonData, &unmarshaled)
	if err != nil {
		t.Fatalf("Failed to unmarshal JSON: %v", err)
	}
	
	if unmarshaled.VariantID != "L" {
		t.Errorf("Expected variant_id 'L', got '%s'", unmarshaled.VariantID)
	}
	if unmarshaled.VariantName != "Size L" {
		t.Errorf("Expected variant_name 'Size L', got '%s'", unmarshaled.VariantName)
	}
}

func TestOrderItem_JSON_WithoutVariant(t *testing.T) {
	menuItemID := primitive.NewObjectID()
	
	item := OrderItem{
		MenuItemID: menuItemID,
		Name:       "Bánh mì",
		Price:      20000,
		Quantity:   1,
		Subtotal:   20000,
	}
	
	// Marshal to JSON
	jsonData, err := json.Marshal(item)
	if err != nil {
		t.Fatalf("Failed to marshal JSON: %v", err)
	}
	
	// Unmarshal back
	var unmarshaled OrderItem
	err = json.Unmarshal(jsonData, &unmarshaled)
	if err != nil {
		t.Fatalf("Failed to unmarshal JSON: %v", err)
	}
	
	if unmarshaled.VariantID != "" {
		t.Errorf("Expected empty variant_id, got '%s'", unmarshaled.VariantID)
	}
	if unmarshaled.VariantName != "" {
		t.Errorf("Expected empty variant_name, got '%s'", unmarshaled.VariantName)
	}
	if unmarshaled.Name != "Bánh mì" {
		t.Errorf("Expected name 'Bánh mì', got '%s'", unmarshaled.Name)
	}
}

// Test CalculateTotal() with variants

func TestCalculateTotal_WithVariants(t *testing.T) {
	menuItemID1 := primitive.NewObjectID()
	menuItemID2 := primitive.NewObjectID()
	
	order := &Order{
		Items: []OrderItem{
			{
				MenuItemID:  menuItemID1,
				VariantID:   "M",
				Name:        "Cà phê sữa đá",
				VariantName: "Size M",
				Price:       25000,
				Quantity:    2,
			},
			{
				MenuItemID:  menuItemID2,
				VariantID:   "L",
				Name:        "Cà phê sữa đá",
				VariantName: "Size L",
				Price:       30000,
				Quantity:    1,
			},
		},
		Discount: 5000,
	}
	
	order.CalculateTotal()
	
	// Check item subtotals
	if order.Items[0].Subtotal != 50000 {
		t.Errorf("Expected item 0 subtotal 50000, got %f", order.Items[0].Subtotal)
	}
	if order.Items[1].Subtotal != 30000 {
		t.Errorf("Expected item 1 subtotal 30000, got %f", order.Items[1].Subtotal)
	}
	
	// Check order subtotal
	if order.Subtotal != 80000 {
		t.Errorf("Expected order subtotal 80000, got %f", order.Subtotal)
	}
	
	// Check order total (subtotal - discount)
	if order.Total != 75000 {
		t.Errorf("Expected order total 75000, got %f", order.Total)
	}
}

// Test CalculateTotal() without variants (backward compatible)

func TestCalculateTotal_WithoutVariants(t *testing.T) {
	menuItemID1 := primitive.NewObjectID()
	menuItemID2 := primitive.NewObjectID()
	
	order := &Order{
		Items: []OrderItem{
			{
				MenuItemID: menuItemID1,
				Name:       "Bánh mì",
				Price:      20000,
				Quantity:   2,
			},
			{
				MenuItemID: menuItemID2,
				Name:       "Nước ngọt",
				Price:      10000,
				Quantity:   1,
			},
		},
		Discount: 0,
	}
	
	order.CalculateTotal()
	
	// Check item subtotals
	if order.Items[0].Subtotal != 40000 {
		t.Errorf("Expected item 0 subtotal 40000, got %f", order.Items[0].Subtotal)
	}
	if order.Items[1].Subtotal != 10000 {
		t.Errorf("Expected item 1 subtotal 10000, got %f", order.Items[1].Subtotal)
	}
	
	// Check order subtotal
	if order.Subtotal != 50000 {
		t.Errorf("Expected order subtotal 50000, got %f", order.Subtotal)
	}
	
	// Check order total
	if order.Total != 50000 {
		t.Errorf("Expected order total 50000, got %f", order.Total)
	}
}

// Test mixed order (with and without variants)

func TestCalculateTotal_MixedOrder(t *testing.T) {
	menuItemID1 := primitive.NewObjectID()
	menuItemID2 := primitive.NewObjectID()
	
	order := &Order{
		Items: []OrderItem{
			{
				MenuItemID: menuItemID1,
				Name:       "Bánh mì",
				Price:      20000,
				Quantity:   1,
			},
			{
				MenuItemID:  menuItemID2,
				VariantID:   "L",
				Name:        "Cà phê sữa đá",
				VariantName: "Size L",
				Price:       30000,
				Quantity:    2,
			},
		},
		Discount: 10000,
	}
	
	order.CalculateTotal()
	
	// Check item subtotals
	if order.Items[0].Subtotal != 20000 {
		t.Errorf("Expected item 0 subtotal 20000, got %f", order.Items[0].Subtotal)
	}
	if order.Items[1].Subtotal != 60000 {
		t.Errorf("Expected item 1 subtotal 60000, got %f", order.Items[1].Subtotal)
	}
	
	// Check order subtotal
	if order.Subtotal != 80000 {
		t.Errorf("Expected order subtotal 80000, got %f", order.Subtotal)
	}
	
	// Check order total (subtotal - discount)
	if order.Total != 70000 {
		t.Errorf("Expected order total 70000, got %f", order.Total)
	}
}

// Test CalculateTotal with amount paid

func TestCalculateTotal_WithAmountPaid(t *testing.T) {
	menuItemID := primitive.NewObjectID()
	
	order := &Order{
		Items: []OrderItem{
			{
				MenuItemID:  menuItemID,
				VariantID:   "L",
				Name:        "Cà phê sữa đá",
				VariantName: "Size L",
				Price:       30000,
				Quantity:    2,
			},
		},
		Discount:   0,
		AmountPaid: 50000,
	}
	
	order.CalculateTotal()
	
	// Check order total
	if order.Total != 60000 {
		t.Errorf("Expected order total 60000, got %f", order.Total)
	}
	
	// Check amount due
	if order.AmountDue != 10000 {
		t.Errorf("Expected amount due 10000, got %f", order.AmountDue)
	}
}

// Test CalculateTotal with overpayment

func TestCalculateTotal_WithOverpayment(t *testing.T) {
	menuItemID := primitive.NewObjectID()
	
	order := &Order{
		Items: []OrderItem{
			{
				MenuItemID: menuItemID,
				Name:       "Bánh mì",
				Price:      20000,
				Quantity:   1,
			},
		},
		Discount:   0,
		AmountPaid: 30000,
	}
	
	order.CalculateTotal()
	
	// Check order total
	if order.Total != 20000 {
		t.Errorf("Expected order total 20000, got %f", order.Total)
	}
	
	// Check amount due (should be 0, not negative)
	if order.AmountDue != 0 {
		t.Errorf("Expected amount due 0, got %f", order.AmountDue)
	}
}

// Test CalculateTotal with discount greater than subtotal

func TestCalculateTotal_DiscountGreaterThanSubtotal(t *testing.T) {
	menuItemID := primitive.NewObjectID()
	
	order := &Order{
		Items: []OrderItem{
			{
				MenuItemID: menuItemID,
				Name:       "Bánh mì",
				Price:      20000,
				Quantity:   1,
			},
		},
		Discount: 30000,
	}
	
	order.CalculateTotal()
	
	// Check order total (should be 0, not negative)
	if order.Total != 0 {
		t.Errorf("Expected order total 0, got %f", order.Total)
	}
}
