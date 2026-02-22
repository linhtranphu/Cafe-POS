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
		mongoURI = "mongodb://admin:password123@localhost:27017/cafe_pos?replicaSet=rs0&authSource=admin"
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

	// Create repositories
	printerRepo := mongodb.NewPrinterConfigRepository(db)
	templateRepo := mongodb.NewPrintTemplateRepository(db)
	settingsRepo := mongodb.NewShopSettingsRepository(db)

	fmt.Println("=== Checking Printer Configuration ===")
	fmt.Println()

	ctx = context.Background()

	// Check all printers
	fmt.Println("1. All Printers:")
	allPrinters, err := printerRepo.FindAll(ctx)
	if err != nil {
		log.Printf("Error fetching printers: %v", err)
	} else if len(allPrinters) == 0 {
		fmt.Println("   ❌ No printers configured")
	} else {
		for _, p := range allPrinters {
			fmt.Printf("   - %s (Type: %s, Default: %v, IP: %s:%d)\n", 
				p.Name, p.Type, p.IsDefault, p.IPAddress, p.Port)
		}
	}
	fmt.Println()

	// Check default BILL printer
	fmt.Println("2. Default BILL Printer:")
	billPrinter, err := printerRepo.FindDefault(ctx, printing.PrinterTypeBill)
	if err != nil {
		fmt.Printf("   ❌ Error: %v\n", err)
	} else if billPrinter == nil {
		fmt.Println("   ❌ No default BILL printer found")
	} else {
		fmt.Printf("   ✅ Found: %s (IP: %s:%d)\n", billPrinter.Name, billPrinter.IPAddress, billPrinter.Port)
	}
	fmt.Println()

	// Check default LABEL printer
	fmt.Println("3. Default LABEL Printer:")
	labelPrinter, err := printerRepo.FindDefault(ctx, printing.PrinterTypeLabel)
	if err != nil {
		fmt.Printf("   ❌ Error: %v\n", err)
	} else if labelPrinter == nil {
		fmt.Println("   ❌ No default LABEL printer found")
	} else {
		fmt.Printf("   ✅ Found: %s (IP: %s:%d)\n", labelPrinter.Name, labelPrinter.IPAddress, labelPrinter.Port)
	}
	fmt.Println()

	// Check templates
	fmt.Println("4. Templates:")
	billTemplate, err := templateRepo.FindDefault(ctx, printing.TemplateTypeBill)
	if err != nil {
		fmt.Printf("   ❌ Bill template error: %v\n", err)
	} else if billTemplate == nil {
		fmt.Println("   ❌ No default BILL template")
	} else {
		fmt.Printf("   ✅ Bill template: %s\n", billTemplate.Name)
	}

	labelTemplate, err := templateRepo.FindDefault(ctx, printing.TemplateTypeLabel)
	if err != nil {
		fmt.Printf("   ❌ Label template error: %v\n", err)
	} else if labelTemplate == nil {
		fmt.Println("   ❌ No default LABEL template")
	} else {
		fmt.Printf("   ✅ Label template: %s\n", labelTemplate.Name)
	}
	fmt.Println()

	// Check auto-print setting
	fmt.Println("5. Auto-Print Setting:")
	settings, err := settingsRepo.GetSettings(ctx)
	if err != nil {
		fmt.Printf("   ⚠️  Error fetching settings: %v\n", err)
	} else if settings == nil {
		fmt.Println("   ⚠️  No shop settings found")
	} else {
		if settings.AutoPrintEnabled {
			fmt.Println("   ✅ Auto-print ENABLED")
		} else {
			fmt.Println("   ❌ Auto-print DISABLED")
		}
	}
	fmt.Println()

	fmt.Println("=== Check Complete ===")
	fmt.Println()

	// Summary
	fmt.Println("Summary:")
	canAutoPrint := true
	if billPrinter == nil {
		fmt.Println("   ❌ Missing: Default BILL printer")
		canAutoPrint = false
	}
	if labelPrinter == nil {
		fmt.Println("   ❌ Missing: Default LABEL printer")
		canAutoPrint = false
	}
	if billTemplate == nil {
		fmt.Println("   ❌ Missing: Default BILL template")
		canAutoPrint = false
	}
	if labelTemplate == nil {
		fmt.Println("   ❌ Missing: Default LABEL template")
		canAutoPrint = false
	}
	if settings != nil && !settings.AutoPrintEnabled {
		fmt.Println("   ❌ Auto-print is disabled")
		canAutoPrint = false
	}

	if canAutoPrint {
		fmt.Println("   ✅ Auto-print should work!")
	} else {
		fmt.Println("   ❌ Auto-print will NOT work - fix issues above")
	}
}
