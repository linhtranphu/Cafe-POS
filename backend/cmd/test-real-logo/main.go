package main

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"time"

	"cafe-pos/backend/application/services"
	"cafe-pos/backend/domain/order"
	"cafe-pos/backend/domain/settings"
)

func main() {
	log.Println("=== Test with Real Logo from Settings API ===")
	
	// 1. Fetch settings from API
	log.Println("Fetching settings from http://localhost:3000/api/settings...")
	resp, err := http.Get("http://localhost:3000/api/settings")
	if err != nil {
		log.Fatalf("Failed to fetch settings: %v", err)
	}
	defer resp.Body.Close()
	
	if resp.StatusCode != 200 {
		body, _ := io.ReadAll(resp.Body)
		log.Fatalf("API returned status %d: %s", resp.StatusCode, string(body))
	}
	
	var shopSettings settings.ShopSettings
	if err := json.NewDecoder(resp.Body).Decode(&shopSettings); err != nil {
		log.Fatalf("Failed to decode settings: %v", err)
	}
	
	log.Printf("✅ Settings loaded:")
	log.Printf("  - Shop Name: %s", shopSettings.ShopName)
	log.Printf("  - Show Logo: %v", shopSettings.ShowLogo)
	log.Printf("  - Logo URL: %s", shopSettings.LogoURL)
	log.Printf("  - Show Address: %v", shopSettings.ShowAddress)
	log.Printf("  - Show Phone: %v", shopSettings.ShowPhone)
	
	// 2. Create test order
	testOrder := &order.Order{
		OrderNumber:   "REAL-TEST-001",
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
	
	// 3. Create renderer
	log.Println("\nCreating chromedp renderer...")
	renderer, err := services.NewChromedpBillRendererOptimized()
	if err != nil {
		log.Fatalf("Failed to create renderer: %v", err)
	}
	defer renderer.Close()
	
	// 4. Test render to ESC/POS
	log.Println("\n=== Test 1: Render to ESC/POS ===")
	escposData, err := renderer.RenderBillToESCPOS(testOrder, &shopSettings)
	if err != nil {
		log.Printf("❌ ERROR: %v", err)
	} else {
		log.Printf("✅ SUCCESS: ESC/POS data generated: %d bytes", len(escposData))
		
		// Save to file
		filename := "test_real_logo.bin"
		if err := os.WriteFile(filename, escposData, 0644); err != nil {
			log.Printf("ERROR: Failed to save file: %v", err)
		} else {
			log.Printf("✅ Saved to %s", filename)
		}
	}
	
	// 5. Test save preview
	log.Println("\n=== Test 2: Save Preview Image ===")
	previewFile := "test_real_logo_preview.png"
	if err := renderer.SavePreviewImage(testOrder, &shopSettings, previewFile); err != nil {
		log.Printf("❌ ERROR: %v", err)
	} else {
		log.Printf("✅ SUCCESS: Preview saved to %s", previewFile)
	}
	
	fmt.Println("\n=== Test completed ===")
	fmt.Println("Kiểm tra các file:")
	fmt.Println("- test_real_logo.bin (ESC/POS data)")
	fmt.Println("- test_real_logo_preview.png (Preview image)")
	fmt.Println("- debug_rendered.html (Rendered HTML for inspection)")
	fmt.Println("\nMở preview image để xem logo có hiển thị đúng không:")
	fmt.Println("  open test_real_logo_preview.png")
}
