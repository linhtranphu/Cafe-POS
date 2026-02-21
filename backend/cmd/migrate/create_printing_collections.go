package main

import (
	"context"
	"fmt"
	"log"
	"os"
	"time"

	"cafe-pos/backend/infrastructure/mongodb"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

func main() {
	// Load environment variables
	mongoURI := os.Getenv("MONGODB_URI")
	if mongoURI == "" {
		mongoURI = "mongodb://localhost:27017"
	}

	dbName := os.Getenv("MONGODB_DATABASE")
	if dbName == "" {
		dbName = "cafe_pos"
	}

	log.Printf("Connecting to MongoDB: %s", mongoURI)

	// Connect to MongoDB
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	client, err := mongo.Connect(ctx, options.Client().ApplyURI(mongoURI))
	if err != nil {
		log.Fatalf("Failed to connect to MongoDB: %v", err)
	}
	defer client.Disconnect(ctx)

	// Ping to verify connection
	if err := client.Ping(ctx, nil); err != nil {
		log.Fatalf("Failed to ping MongoDB: %v", err)
	}

	db := client.Database(dbName)
	fmt.Printf("Connected to database: %s\n", dbName)

	// List of collections to create
	collections := []string{
		"print_jobs",
		"printer_configs",
		"print_templates",
	}

	// Create collections
	fmt.Println("\n=== Creating printing collections ===")
	
	existingCollections, err := db.ListCollectionNames(ctx, map[string]interface{}{})
	if err != nil {
		log.Fatalf("Failed to list collections: %v", err)
	}

	existingMap := make(map[string]bool)
	for _, name := range existingCollections {
		existingMap[name] = true
	}

	for _, collName := range collections {
		if existingMap[collName] {
			fmt.Printf("✓ Collection '%s' already exists\n", collName)
		} else {
			if err := db.CreateCollection(ctx, collName); err != nil {
				log.Fatalf("Failed to create collection '%s': %v", collName, err)
			}
			fmt.Printf("✓ Created collection '%s'\n", collName)
		}
	}

	// Create indexes for print_jobs
	fmt.Println("\n=== Creating indexes for print_jobs ===")
	printJobRepo := mongodb.NewPrintJobRepository(db)
	if err := printJobRepo.CreateIndexes(ctx); err != nil {
		log.Fatalf("Failed to create indexes for print_jobs: %v", err)
	}
	fmt.Println("✓ Created indexes:")
	fmt.Println("  - order_id_idx (order_id)")
	fmt.Println("  - status_created_idx (status, created_at)")
	fmt.Println("  - created_at_ttl_idx (created_at) with TTL 7 days")

	// Create indexes for printer_configs
	fmt.Println("\n=== Creating indexes for printer_configs ===")
	printerConfigRepo := mongodb.NewPrinterConfigRepository(db)
	if err := printerConfigRepo.CreateIndexes(ctx); err != nil {
		log.Fatalf("Failed to create indexes for printer_configs: %v", err)
	}
	fmt.Println("✓ Created indexes:")
	fmt.Println("  - type_default_idx (type, is_default)")
	fmt.Println("  - enabled_idx (is_enabled)")

	// Create indexes for print_templates
	fmt.Println("\n=== Creating indexes for print_templates ===")
	printTemplateRepo := mongodb.NewPrintTemplateRepository(db)
	if err := printTemplateRepo.CreateIndexes(ctx); err != nil {
		log.Fatalf("Failed to create indexes for print_templates: %v", err)
	}
	fmt.Println("✓ Created indexes:")
	fmt.Println("  - type_default_idx (type, is_default)")

	fmt.Println("\n=== Migration completed successfully ===")
	
	// Print schema information
	fmt.Println("\n=== Collection Schemas ===")
	
	fmt.Println("\nprint_jobs:")
	fmt.Println("  - _id: ObjectId")
	fmt.Println("  - type: String (BILL|LABEL)")
	fmt.Println("  - order_id: ObjectId")
	fmt.Println("  - order_number: String")
	fmt.Println("  - printer_id: ObjectId")
	fmt.Println("  - content: String (rendered content)")
	fmt.Println("  - status: String (PENDING|PRINTING|COMPLETED|FAILED)")
	fmt.Println("  - retry_count: Number")
	fmt.Println("  - max_retries: Number")
	fmt.Println("  - error_msg: String")
	fmt.Println("  - created_at: Date")
	fmt.Println("  - updated_at: Date")
	fmt.Println("  - printed_at: Date (optional)")

	fmt.Println("\nprinter_configs:")
	fmt.Println("  - _id: ObjectId")
	fmt.Println("  - name: String")
	fmt.Println("  - type: String (BILL|LABEL)")
	fmt.Println("  - connection_type: String (NETWORK|USB)")
	fmt.Println("  - ip_address: String (optional)")
	fmt.Println("  - port: Number (optional)")
	fmt.Println("  - usb_path: String (optional)")
	fmt.Println("  - paper_width: Number (mm: 58 or 80)")
	fmt.Println("  - is_default: Boolean")
	fmt.Println("  - is_enabled: Boolean")
	fmt.Println("  - created_at: Date")
	fmt.Println("  - updated_at: Date")

	fmt.Println("\nprint_templates:")
	fmt.Println("  - _id: ObjectId")
	fmt.Println("  - type: String (BILL|LABEL)")
	fmt.Println("  - name: String")
	fmt.Println("  - content: String (template string)")
	fmt.Println("  - is_default: Boolean")
	fmt.Println("  - created_at: Date")
	fmt.Println("  - updated_at: Date")

	fmt.Println("\n=== Usage ===")
	fmt.Println("To run this migration:")
	fmt.Println("  go run backend/cmd/migrate/create_printing_collections.go")
	fmt.Println("\nOr with custom MongoDB URI:")
	fmt.Println("  MONGODB_URI=mongodb://localhost:27017 go run backend/cmd/migrate/create_printing_collections.go")
}
