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
	log.Println("=== Test with Uploaded Logo ===")
	
	// Check which logo to use
	logoPath := "uploads/logos/logo_24094.jpeg"
	if _, err := os.Stat(logoPath); os.IsNotExist(err) {
		logoPath = "uploads/logo.png"
	}
	
	log.Printf("Using logo: %s", logoPath)
	
	// Create shop settings with uploaded logo
	shopSettings := &settings.ShopSettings{
		ShopName:          "CAFE KIRO TEST",
		ShopAddress:       "123 Nguyễn Huệ, Quận 1, TP.HCM",
		ShopPhone:         "0901234567",
		ShowLogo:          true,
		LogoURL:           logoPath,
		ShowAddress:       true,
		ShowPhone:         true,
		CustomMessage:     "Cảm ơn quý khách! Hẹn gặp lại!",
		ShowCustomMessage: true,
	}
	
	// Create test order
	testOrder := &order.Order{
		OrderNumber:   "UPLOAD-TEST-001",
		WaiterName:    "Nguyễn Văn A",
		PaymentMethod: order.PaymentCash,
		CreatedAt:     time.Now(),
		Total:         250000,
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
				VariantName: "Size M",
				Quantity:    2,
				Price:       45000,
				Subtotal:    90000,
			},
			{
				Name:        "Bánh mì thịt",
				VariantName: "",
				Quantity:    2,
				Price:       35000,
				Subtotal:    70000,
			},
			{
				Name:        "Nước cam ép",
				VariantName: "",
				Quantity:    1,
				Price:       20000,
				Subtotal:    20000,
			},
		},
	}
	
	// Create renderer
	log.Println("\nCreating chromedp renderer...")
	renderer, err := services.NewChromedpBillRendererOptimized()
	if err != nil {
		log.Fatalf("Failed to create renderer: %v", err)
	}
	defer renderer.Close()
	
	// Test render to ESC/POS
	log.Println("\n=== Test 1: Render to ESC/POS ===")
	escposData, err := renderer.RenderBillToESCPOS(testOrder, shopSettings)
	if err != nil {
		log.Printf("❌ ERROR: %v", err)
	} else {
		log.Printf("✅ SUCCESS: ESC/POS data generated: %d bytes", len(escposData))
		
		// Save to file
		filename := "test_uploaded_logo.bin"
		if err := os.WriteFile(filename, escposData, 0644); err != nil {
			log.Printf("ERROR: Failed to save file: %v", err)
		} else {
			log.Printf("✅ Saved to %s", filename)
		}
	}
	
	// Test save preview
	log.Println("\n=== Test 2: Save Preview Image ===")
	previewFile := "test_uploaded_logo_preview.png"
	if err := renderer.SavePreviewImage(testOrder, shopSettings, previewFile); err != nil {
		log.Printf("❌ ERROR: %v", err)
	} else {
		log.Printf("✅ SUCCESS: Preview saved to %s", previewFile)
	}
	
	fmt.Println("\n=== Test completed ===")
	fmt.Println("Kiểm tra các file:")
	fmt.Println("- test_uploaded_logo.bin (ESC/POS data)")
	fmt.Println("- test_uploaded_logo_preview.png (Preview image)")
	fmt.Println("- debug_rendered.html (Rendered HTML for inspection)")
	fmt.Println("\nMở preview image để xem logo có hiển thị đúng không:")
	fmt.Println("  open test_uploaded_logo_preview.png")
}
