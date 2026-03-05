package main

import (
	"context"
	"encoding/base64"
	"fmt"
	"log"
	"os"
	"time"

	"cafe-pos/backend/infrastructure/mongodb"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

func main() {
	// Connect to MongoDB
	mongoURI := "mongodb://admin:password123@localhost:27017/cafe_pos?replicaSet=rs0&authSource=admin"
	
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	
	client, err := mongo.Connect(ctx, options.Client().ApplyURI(mongoURI))
	if err != nil {
		log.Fatal("Failed to connect to MongoDB:", err)
	}
	defer client.Disconnect(context.Background())
	
	// Get shop settings
	settingsRepo := mongodb.NewShopSettingsRepository(client.Database("cafe_pos"))
	shopSettings, err := settingsRepo.GetSettings(context.Background())
	if err != nil {
		log.Fatal("Failed to get shop settings:", err)
	}
	
	fmt.Println("=== SHOP SETTINGS ===")
	fmt.Printf("Shop Name: %s\n", shopSettings.ShopName)
	fmt.Printf("Show Logo: %v\n", shopSettings.ShowLogo)
	fmt.Printf("Logo URL: %s\n", shopSettings.LogoURL)
	fmt.Println()
	
	// Check if logo file exists
	if shopSettings.LogoURL != "" {
		// Try original path first
		logoPath := shopSettings.LogoURL
		_, err := os.Stat(logoPath)
		
		// If failed and path starts with /, try prepending "."
		if err != nil && len(logoPath) > 0 && logoPath[0] == '/' {
			logoPath = "." + logoPath
			_, err = os.Stat(logoPath)
		}
		
		if err == nil {
			fmt.Printf("✅ Logo file exists: %s\n", logoPath)
			
			// Get file info
			fileInfo, _ := os.Stat(logoPath)
			fmt.Printf("   File size: %d bytes\n", fileInfo.Size())
			
			// Try to load as base64
			logoBase64, err := loadImageAsBase64(shopSettings.LogoURL)
			if err != nil {
				fmt.Printf("❌ Failed to load logo as base64: %v\n", err)
			} else {
				fmt.Printf("✅ Successfully loaded logo as base64\n")
				fmt.Printf("   Base64 length: %d characters\n", len(logoBase64))
				fmt.Printf("   Base64 prefix: %s...\n", logoBase64[:min(100, len(logoBase64))])
			}
		} else {
			fmt.Printf("❌ Logo file does not exist: %s\n", shopSettings.LogoURL)
			fmt.Printf("   Error: %v\n", err)
		}
	} else {
		fmt.Println("⚠️  No logo URL configured")
	}
	
	fmt.Println()
	
	// List all files in uploads directory
	fmt.Println("=== FILES IN UPLOADS DIRECTORY ===")
	uploadsDir := "./uploads"
	if entries, err := os.ReadDir(uploadsDir); err == nil {
		for _, entry := range entries {
			if !entry.IsDir() {
				info, _ := entry.Info()
				fmt.Printf("  - %s (%d bytes)\n", entry.Name(), info.Size())
			}
		}
	} else {
		fmt.Printf("❌ Cannot read uploads directory: %v\n", err)
	}
}

func loadImageAsBase64(path string) (string, error) {
	// Try original path first
	data, err := os.ReadFile(path)
	
	// If failed and path starts with /, try prepending "."
	if err != nil && len(path) > 0 && path[0] == '/' {
		data, err = os.ReadFile("." + path)
	}
	
	if err != nil {
		return "", fmt.Errorf("failed to read image from %s: %w", path, err)
	}

	// Detect MIME type
	mimeType := "image/jpeg"
	if len(data) > 4 && data[0] == 0x89 && data[1] == 0x50 && data[2] == 0x4E && data[3] == 0x47 {
		mimeType = "image/png"
	}

	encoded := base64.StdEncoding.EncodeToString(data)
	return fmt.Sprintf("data:%s;base64,%s", mimeType, encoded), nil
}

func min(a, b int) int {
	if a < b {
		return a
	}
	return b
}

// ShopSettings represents shop configuration
type ShopSettings struct {
	ID                bson.M `bson:"_id,omitempty"`
	ShopName          string `bson:"shop_name"`
	ShopAddress       string `bson:"shop_address"`
	ShopPhone         string `bson:"shop_phone"`
	LogoURL           string `bson:"logo_url"`
	ShowLogo          bool   `bson:"show_logo"`
	ShowAddress       bool   `bson:"show_address"`
	ShowPhone         bool   `bson:"show_phone"`
	CustomMessage     string `bson:"custom_message"`
	ShowCustomMessage bool   `bson:"show_custom_message"`
}
