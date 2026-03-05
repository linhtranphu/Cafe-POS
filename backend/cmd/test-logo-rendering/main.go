package main

import (
	"context"
	"fmt"
	"log"
	"time"

	"cafe-pos/backend/application/services"
	"cafe-pos/backend/domain/order"
	"cafe-pos/backend/domain/printing"
	"cafe-pos/backend/infrastructure/mongodb"

	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

func main() {
	// Connect to MongoDB
	ctx := context.Background()
	mongoURI := "mongodb://admin:password123@localhost:27017/cafe_pos?replicaSet=rs0&authSource=admin"
	
	client, err := mongo.Connect(ctx, options.Client().ApplyURI(mongoURI))
	if err != nil {
		log.Fatal("Failed to connect to MongoDB:", err)
	}
	defer client.Disconnect(ctx)
	
	db := client.Database("cafe_pos")
	
	// Create repositories
	shopSettingsRepo := mongodb.NewShopSettingsRepository(db)
	templateRepo := mongodb.NewPrintTemplateRepository(db)
	
	// Create template renderer
	templateRenderer := services.NewTemplateRenderer()
	
	// Get shop settings
	shopSettings, err := shopSettingsRepo.GetSettings(ctx)
	if err != nil {
		log.Fatal("Failed to get shop settings:", err)
	}
	
	fmt.Println("=== Shop Settings ===")
	fmt.Printf("ShowLogo: %v\n", shopSettings.ShowLogo)
	fmt.Printf("LogoURL: %s\n", shopSettings.LogoURL)
	
	// Get default template
	templates, err := templateRepo.FindByType(ctx, printing.TemplateTypeBill)
	if err != nil {
		log.Fatal("Failed to get templates:", err)
	}
	
	var defaultTemplate *printing.PrintTemplate
	for _, t := range templates {
		if t.IsDefault {
			defaultTemplate = t
			break
		}
	}
	
	if defaultTemplate == nil {
		log.Fatal("No default template found")
	}
	
	fmt.Println("\n=== Template ===")
	fmt.Printf("Name: %s\n", defaultTemplate.Name)
	hasLogoMarker := false
	for i := 0; i < len(defaultTemplate.Content)-6; i++ {
		if defaultTemplate.Content[i:i+6] == "[LOGO]" {
			hasLogoMarker = true
			break
		}
	}
	fmt.Printf("Has [LOGO] marker: %v\n", hasLogoMarker)
	
	// Create test order
	testOrder := &order.Order{
		ID:          primitive.NewObjectID(),
		OrderNumber: "TEST-LOGO-001",
		Items: []order.OrderItem{
			{
				Name:     "Test Coffee",
				Quantity: 1,
				Price:    30000,
				Subtotal: 30000,
			},
		},
		Subtotal:   30000,
		Discount:   0,
		Total:      30000,
		WaiterName: "Test",
		CreatedAt:  time.Now(),
		UpdatedAt:  time.Now(),
	}
	
	fmt.Println("\n=== Rendering Bill ===")
	content, err := templateRenderer.RenderBill(testOrder, defaultTemplate, shopSettings)
	if err != nil {
		log.Fatal("Failed to render bill:", err)
	}
	
	fmt.Println("\n=== Rendered Content (first 500 chars) ===")
	if len(content) > 500 {
		fmt.Println(content[:500])
	} else {
		fmt.Println(content)
	}
	
	fmt.Println("\n=== Logo Check ===")
	hasEscPosImage := false
	for i := 0; i < len(content)-2; i++ {
		if content[i] == '\x1D' && content[i+1] == '\x76' && content[i+2] == '\x30' {
			hasEscPosImage = true
			break
		}
	}
	fmt.Printf("Has ESC/POS image command (GS v 0): %v\n", hasEscPosImage)
	fmt.Printf("Content length: %d bytes\n", len(content))
	
	if hasEscPosImage {
		fmt.Println("\n✅ SUCCESS: Logo ESC/POS commands found in content!")
	} else {
		fmt.Println("\n❌ FAILED: No logo ESC/POS commands in content")
	}
}
