package main

import (
	"context"
	"fmt"
	"log"
	"os"
	"time"

	"cafe-pos/backend/domain/printing"
	"cafe-pos/backend/infrastructure/mongodb"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

func main() {
	// Get MongoDB URI from environment
	mongoURI := os.Getenv("MONGODB_URI")
	if mongoURI == "" {
		log.Fatal("MONGODB_URI environment variable is required")
	}

	// Connect to MongoDB
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	client, err := mongo.Connect(ctx, options.Client().ApplyURI(mongoURI))
	if err != nil {
		log.Fatalf("Failed to connect to MongoDB: %v", err)
	}
	defer client.Disconnect(ctx)

	// Get database
	dbName := os.Getenv("MONGODB_DATABASE")
	if dbName == "" {
		dbName = "cafe_pos"
	}
	db := client.Database(dbName)

	// Create template repository
	templateRepo := mongodb.NewPrintTemplateRepository(db)

	// Create default bill template
	defaultBillTemplate := &printing.PrintTemplate{
		Type:      printing.TemplateTypeBill,
		Name:      "Hóa đơn mặc định",
		IsDefault: true,
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
		Content: `{{.ShopName}}
{{if .ShowAddress}}{{.ShopAddress}}{{end}}
{{if .ShowPhone}}Tel: {{.ShopPhone}}{{end}}
================================
HOA DON BAN HANG
Order: {{.Order.OrderNumber}}
Ngay: {{.Order.CreatedAt.Format "02/01/2006 15:04"}}
{{if .Order.TableNumber}}Ban: {{.Order.TableNumber}}{{end}}
{{if .Order.CustomerName}}Khach: {{.Order.CustomerName}}{{end}}
================================
{{range .Order.Items}}{{.Name}}{{if .VariantName}} - {{.VariantName}}{{end}}
  {{.Quantity}} x {{formatPrice .UnitPrice}} = {{formatPrice .TotalPrice}}
{{end}}================================
Tong tien: {{formatPrice .Order.Subtotal}} VND
{{if gt .Order.DiscountAmount 0.0}}Giam gia: -{{formatPrice .Order.DiscountAmount}} VND
{{end}}{{if gt .Order.TaxAmount 0.0}}Thue: {{formatPrice .Order.TaxAmount}} VND
{{end}}TONG CONG: {{formatPrice .Order.Total}} VND
================================
{{if .ShowCustomMessage}}{{.CustomMessage}}
{{end}}Cam on quy khach!
Hen gap lai!`,
	}

	// Create default label template
	defaultLabelTemplate := &printing.PrintTemplate{
		Type:      printing.TemplateTypeLabel,
		Name:      "Nhãn món mặc định",
		IsDefault: true,
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
		Content: `{{.ShopName}}
================================
Order: {{.Order.OrderNumber}}
{{if .Order.TableNumber}}Ban: {{.Order.TableNumber}}{{end}}
================================
{{with index .Order.Items .ItemIndex}}{{.Name}}
{{if .VariantName}}Size: {{.VariantName}}{{end}}
So luong: {{.Quantity}}
{{if .Note}}Ghi chu: {{.Note}}{{end}}{{end}}
================================
Thoi gian: {{.PrintTime.Format "15:04"}}`,
	}

	// Check if default bill template already exists
	ctx = context.Background()
	existingBillTemplates, err := templateRepo.FindByType(ctx, printing.TemplateTypeBill)
	if err != nil {
		log.Fatalf("Failed to check existing bill templates: %v", err)
	}

	hasBillDefault := false
	for _, tmpl := range existingBillTemplates {
		if tmpl.IsDefault {
			hasBillDefault = true
			log.Printf("Default bill template already exists: %s (ID: %s)", tmpl.Name, tmpl.ID.Hex())
			break
		}
	}

	if !hasBillDefault {
		if err := templateRepo.Create(ctx, defaultBillTemplate); err != nil {
			log.Fatalf("Failed to create default bill template: %v", err)
		}
		log.Printf("✅ Created default bill template: %s", defaultBillTemplate.Name)
	} else {
		log.Println("ℹ️  Default bill template already exists, skipping creation")
	}

	// Check if default label template already exists
	existingLabelTemplates, err := templateRepo.FindByType(ctx, printing.TemplateTypeLabel)
	if err != nil {
		log.Fatalf("Failed to check existing label templates: %v", err)
	}

	hasLabelDefault := false
	for _, tmpl := range existingLabelTemplates {
		if tmpl.IsDefault {
			hasLabelDefault = true
			log.Printf("Default label template already exists: %s (ID: %s)", tmpl.Name, tmpl.ID.Hex())
			break
		}
	}

	if !hasLabelDefault {
		if err := templateRepo.Create(ctx, defaultLabelTemplate); err != nil {
			log.Fatalf("Failed to create default label template: %v", err)
		}
		log.Printf("✅ Created default label template: %s", defaultLabelTemplate.Name)
	} else {
		log.Println("ℹ️  Default label template already exists, skipping creation")
	}

	fmt.Println("\n✅ Default template setup complete!")
}
