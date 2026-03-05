package main

import (
	"fmt"
	"log"
	"os"
	"time"

	"cafe-pos/backend/application/services"
	"cafe-pos/backend/domain/order"
	"cafe-pos/backend/domain/settings"
)

func main() {
	fmt.Println("🚀 Test Chromedp Bill Renderer")
	fmt.Println("================================")

	// Create renderer
	renderer, err := services.NewChromedpBillRenderer()
	if err != nil {
		log.Fatalf("Failed to create renderer: %v", err)
	}
	defer renderer.Close()

	// Mock order data
	ord := &order.Order{
		OrderNumber:   "20260226-123456-001",
		WaiterName:    "Waiter1",
		PaymentMethod: "Tiền mặt",
		CreatedAt:     time.Now(),
		Total:         105000,
		Items: []order.OrderItem{
			{
				Name:      "Cà phê sữa đá",
				Quantity:  2,
				Price:     25000,
				Subtotal:  50000,
			},
			{
				Name:      "Trà đào cam sả",
				Quantity:  1,
				Price:     35000,
				Subtotal:  35000,
			},
			{
				Name:      "Bánh mì",
				Quantity:  1,
				Price:     20000,
				Subtotal:  20000,
			},
		},
	}

	// Mock shop settings
	shopSettings := &settings.ShopSettings{
		ShopName:          "Tiệm cà phê Ông Tạ",
		ShopAddress:       "Đ/c: 10/8 Trần Nhật Duật, P. Tân Định, Quận 1, HCM",
		ShopPhone:         "0906990602",
		ShowLogo:          true,
		LogoURL:           "backend/uploads/logos/logo_24094.jpeg",
		ShowAddress:       true,
		ShowPhone:         true,
		ShowCustomMessage: true,
		CustomMessage:     "Cảm ơn quý khách!",
	}

	fmt.Println("📝 Rendering bill to ESC/POS...")
	escposData, err := renderer.RenderBillToESCPOS(ord, shopSettings)
	if err != nil {
		log.Fatalf("Failed to render: %v", err)
	}

	fmt.Printf("✓ Generated %d bytes of ESC/POS data\n", len(escposData))

	// Optional: Save to file for inspection
	filename := "test_bill_escpos.bin"
	if err := os.WriteFile(filename, escposData, 0644); err != nil {
		log.Printf("Warning: Failed to save file: %v", err)
	} else {
		fmt.Printf("✓ Saved to %s\n", filename)
	}

	// Test sending to printer (uncomment if you have a printer)
	// printerIP := "192.168.1.100:9100"
	// fmt.Printf("📤 Sending to printer %s...\n", printerIP)
	// if err := services.SendToPrinter(printerIP, escposData); err != nil {
	// 	log.Printf("Failed to send to printer: %v", err)
	// } else {
	// 	fmt.Println("✓ Sent to printer successfully")
	// }

	fmt.Println("\n✅ Test completed!")
}
