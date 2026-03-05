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
	log.Println("=== Test Chromedp Optimized Capture ===")
	
	// Tạo renderer
	renderer, err := services.NewChromedpBillRendererOptimized()
	if err != nil {
		log.Fatalf("Failed to create renderer: %v", err)
	}
	defer renderer.Close()
	
	// Tạo test order
	testOrder := &order.Order{
		OrderNumber:   "TEST-001",
		WaiterName:    "Nguyễn Văn A",
		PaymentMethod: order.PaymentCash,
		CreatedAt:     time.Now(),
		Total:         150000,
		Items: []order.OrderItem{
			{
				Name:        "Cà phê sữa đá",
				VariantName: "Size L",
				Quantity:    2,
				Price:       35000,
				Subtotal:    70000,
			},
			{
				Name:        "Trà sữa trân châu",
				VariantName: "",
				Quantity:    1,
				Price:       45000,
				Subtotal:    45000,
			},
			{
				Name:        "Bánh mì thịt",
				VariantName: "",
				Quantity:    1,
				Price:       35000,
				Subtotal:    35000,
			},
		},
	}
	
	// Tạo shop settings
	shopSettings := &settings.ShopSettings{
		ShopName:          "CAFE KIRO",
		ShopAddress:       "123 Nguyễn Huệ, Q.1, TP.HCM",
		ShopPhone:         "0901234567",
		ShowLogo:          true,
		LogoURL:           "/uploads/logo.png",
		ShowAddress:       true,
		ShowPhone:         true,
		CustomMessage:     "Cảm ơn quý khách! Hẹn gặp lại!",
		ShowCustomMessage: true,
	}
	
	log.Println("Test 1: Render bill to ESC/POS")
	escposData, err := renderer.RenderBillToESCPOS(testOrder, shopSettings)
	if err != nil {
		log.Printf("ERROR: Failed to render: %v", err)
	} else {
		log.Printf("SUCCESS: ESC/POS data generated: %d bytes", len(escposData))
		
		// Save to file
		filename := "test_bill_optimized.bin"
		if err := os.WriteFile(filename, escposData, 0644); err != nil {
			log.Printf("ERROR: Failed to save file: %v", err)
		} else {
			log.Printf("SUCCESS: Saved to %s", filename)
		}
	}
	
	log.Println("\nTest 2: Save preview image")
	previewFile := "test_bill_optimized_preview.png"
	if err := renderer.SavePreviewImage(testOrder, shopSettings, previewFile); err != nil {
		log.Printf("ERROR: Failed to save preview: %v", err)
	} else {
		log.Printf("SUCCESS: Preview saved to %s", previewFile)
	}
	
	fmt.Println("\n=== Test completed ===")
	fmt.Println("Kiểm tra các file:")
	fmt.Println("- test_bill_optimized.bin (ESC/POS data)")
	fmt.Println("- test_bill_optimized_preview.png (Preview image)")
}
