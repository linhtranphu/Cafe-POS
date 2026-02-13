package main

import (
	"context"
	"fmt"
	"log"
	"os"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type Ingredient struct {
	Name     string  `bson:"name"`
	Quantity float64 `bson:"quantity"`
	Unit     string  `bson:"unit"`
}

type MenuItemVariant struct {
	ID                   string       `bson:"id"`
	Name                 string       `bson:"name"`
	Price                float64      `bson:"price"`
	Ingredients          []Ingredient `bson:"ingredients"`
	Available            bool         `bson:"available"`
	IsDefault            bool         `bson:"is_default"`
	CurrentCost          float64      `bson:"current_cost"`
	CostStatus           string       `bson:"cost_status"`
	CostLastCalculatedAt time.Time    `bson:"cost_last_calculated_at"`
}

type MenuItem struct {
	Name        string            `bson:"name"`
	Category    string            `bson:"category"`
	Description string            `bson:"description"`
	Available   bool              `bson:"available"`
	HasVariants bool              `bson:"has_variants"`
	Variants    []MenuItemVariant `bson:"variants,omitempty"`
	Price       float64           `bson:"price,omitempty"`
	Ingredients []Ingredient      `bson:"ingredients,omitempty"`
	CurrentCost float64           `bson:"current_cost,omitempty"`
	CostStatus  string            `bson:"cost_status,omitempty"`
	CreatedAt   time.Time         `bson:"created_at"`
	UpdatedAt   time.Time         `bson:"updated_at"`
}

func main() {
	// Get MongoDB URI from environment or use default
	mongoURI := os.Getenv("MONGODB_URI")
	if mongoURI == "" {
		mongoURI = "mongodb://admin:password123@localhost:27017"
	}

	// Connect to MongoDB
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	client, err := mongo.Connect(ctx, options.Client().ApplyURI(mongoURI))
	if err != nil {
		log.Fatal(err)
	}
	defer client.Disconnect(ctx)

	// Get collection
	collection := client.Database("cafe_pos").Collection("menu_items")

	// Menu items with variants
	menuItems := []interface{}{
		// ========== CÀ PHÊ SỮA ĐÁ (MULTI-SIZE) ==========
		MenuItem{
			Name:        "Cà phê sữa đá",
			Category:    "Cà phê",
			Description: "Cà phê phin truyền thống với sữa đá",
			Available:   true,
			HasVariants: true,
			Variants: []MenuItemVariant{
				{
					ID:        "M",
					Name:      "Size M",
					Price:     25000,
					Available: true,
					IsDefault: true,
					Ingredients: []Ingredient{
						{Name: "Cà phê", Quantity: 20, Unit: "g"},
						{Name: "Sữa đặc", Quantity: 30, Unit: "ml"},
						{Name: "Đá", Quantity: 100, Unit: "g"},
					},
				},
				{
					ID:        "L",
					Name:      "Size L",
					Price:     30000,
					Available: true,
					IsDefault: false,
					Ingredients: []Ingredient{
						{Name: "Cà phê", Quantity: 30, Unit: "g"},
						{Name: "Sữa đặc", Quantity: 45, Unit: "ml"},
						{Name: "Đá", Quantity: 150, Unit: "g"},
					},
				},
				{
					ID:        "XL",
					Name:      "Size XL",
					Price:     35000,
					Available: true,
					IsDefault: false,
					Ingredients: []Ingredient{
						{Name: "Cà phê", Quantity: 40, Unit: "g"},
						{Name: "Sữa đặc", Quantity: 60, Unit: "ml"},
						{Name: "Đá", Quantity: 200, Unit: "g"},
					},
				},
			},
			CreatedAt: time.Now(),
			UpdatedAt: time.Now(),
		},

		// ========== CÀ PHÊ ĐEN ĐÁ (MULTI-SIZE) ==========
		MenuItem{
			Name:        "Cà phê đen đá",
			Category:    "Cà phê",
			Description: "Cà phê phin đen đậm đà",
			Available:   true,
			HasVariants: true,
			Variants: []MenuItemVariant{
				{
					ID:        "M",
					Name:      "Size M",
					Price:     20000,
					Available: true,
					IsDefault: true,
					Ingredients: []Ingredient{
						{Name: "Cà phê", Quantity: 20, Unit: "g"},
						{Name: "Đá", Quantity: 100, Unit: "g"},
					},
				},
				{
					ID:        "L",
					Name:      "Size L",
					Price:     25000,
					Available: true,
					IsDefault: false,
					Ingredients: []Ingredient{
						{Name: "Cà phê", Quantity: 30, Unit: "g"},
						{Name: "Đá", Quantity: 150, Unit: "g"},
					},
				},
			},
			CreatedAt: time.Now(),
			UpdatedAt: time.Now(),
		},

		// ========== TRÀ SỮA TRUYỀN THỐNG (MULTI-SIZE) ==========
		MenuItem{
			Name:        "Trà sữa truyền thống",
			Category:    "Trà sữa",
			Description: "Trà sữa đài loan cổ điển với trân châu",
			Available:   true,
			HasVariants: true,
			Variants: []MenuItemVariant{
				{
					ID:        "M",
					Name:      "Size M",
					Price:     35000,
					Available: true,
					IsDefault: true,
					Ingredients: []Ingredient{
						{Name: "Trà đen", Quantity: 10, Unit: "g"},
						{Name: "Sữa tươi", Quantity: 100, Unit: "ml"},
						{Name: "Đường", Quantity: 20, Unit: "g"},
						{Name: "Trân châu", Quantity: 50, Unit: "g"},
					},
				},
				{
					ID:        "L",
					Name:      "Size L",
					Price:     42000,
					Available: true,
					IsDefault: false,
					Ingredients: []Ingredient{
						{Name: "Trà đen", Quantity: 15, Unit: "g"},
						{Name: "Sữa tươi", Quantity: 150, Unit: "ml"},
						{Name: "Đường", Quantity: 30, Unit: "g"},
						{Name: "Trân châu", Quantity: 75, Unit: "g"},
					},
				},
				{
					ID:        "XL",
					Name:      "Size XL",
					Price:     48000,
					Available: true,
					IsDefault: false,
					Ingredients: []Ingredient{
						{Name: "Trà đen", Quantity: 20, Unit: "g"},
						{Name: "Sữa tươi", Quantity: 200, Unit: "ml"},
						{Name: "Đường", Quantity: 40, Unit: "g"},
						{Name: "Trân châu", Quantity: 100, Unit: "g"},
					},
				},
			},
			CreatedAt: time.Now(),
			UpdatedAt: time.Now(),
		},

		// ========== TRÀ SỮA MATCHA (MULTI-SIZE) ==========
		MenuItem{
			Name:        "Trà sữa matcha",
			Category:    "Trà sữa",
			Description: "Trà xanh matcha Nhật Bản với sữa",
			Available:   true,
			HasVariants: true,
			Variants: []MenuItemVariant{
				{
					ID:        "M",
					Name:      "Size M",
					Price:     40000,
					Available: true,
					IsDefault: true,
					Ingredients: []Ingredient{
						{Name: "Matcha", Quantity: 5, Unit: "g"},
						{Name: "Sữa tươi", Quantity: 100, Unit: "ml"},
						{Name: "Đường", Quantity: 20, Unit: "g"},
					},
				},
				{
					ID:        "L",
					Name:      "Size L",
					Price:     48000,
					Available: true,
					IsDefault: false,
					Ingredients: []Ingredient{
						{Name: "Matcha", Quantity: 8, Unit: "g"},
						{Name: "Sữa tươi", Quantity: 150, Unit: "ml"},
						{Name: "Đường", Quantity: 30, Unit: "g"},
					},
				},
			},
			CreatedAt: time.Now(),
			UpdatedAt: time.Now(),
		},

		// ========== TRÀ ĐÀO CAM SẢ (MULTI-SIZE) ==========
		MenuItem{
			Name:        "Trà đào cam sả",
			Category:    "Trà trái cây",
			Description: "Trà đào tươi mát với cam và sả thơm",
			Available:   true,
			HasVariants: true,
			Variants: []MenuItemVariant{
				{
					ID:        "M",
					Name:      "Size M",
					Price:     35000,
					Available: true,
					IsDefault: true,
					Ingredients: []Ingredient{
						{Name: "Trà xanh", Quantity: 10, Unit: "g"},
						{Name: "Đào", Quantity: 50, Unit: "g"},
						{Name: "Cam", Quantity: 30, Unit: "g"},
						{Name: "Sả", Quantity: 5, Unit: "g"},
						{Name: "Đường", Quantity: 20, Unit: "g"},
					},
				},
				{
					ID:        "L",
					Name:      "Size L",
					Price:     42000,
					Available: true,
					IsDefault: false,
					Ingredients: []Ingredient{
						{Name: "Trà xanh", Quantity: 15, Unit: "g"},
						{Name: "Đào", Quantity: 75, Unit: "g"},
						{Name: "Cam", Quantity: 45, Unit: "g"},
						{Name: "Sả", Quantity: 8, Unit: "g"},
						{Name: "Đường", Quantity: 30, Unit: "g"},
					},
				},
			},
			CreatedAt: time.Now(),
			UpdatedAt: time.Now(),
		},

		// ========== SINH TỐ BƠ (MULTI-SIZE) ==========
		MenuItem{
			Name:        "Sinh tố bơ",
			Category:    "Sinh tố",
			Description: "Sinh tố bơ béo ngậy",
			Available:   true,
			HasVariants: true,
			Variants: []MenuItemVariant{
				{
					ID:        "M",
					Name:      "Size M",
					Price:     40000,
					Available: true,
					IsDefault: true,
					Ingredients: []Ingredient{
						{Name: "Bơ", Quantity: 150, Unit: "g"},
						{Name: "Sữa tươi", Quantity: 100, Unit: "ml"},
						{Name: "Đường", Quantity: 20, Unit: "g"},
						{Name: "Đá", Quantity: 50, Unit: "g"},
					},
				},
				{
					ID:        "L",
					Name:      "Size L",
					Price:     50000,
					Available: true,
					IsDefault: false,
					Ingredients: []Ingredient{
						{Name: "Bơ", Quantity: 250, Unit: "g"},
						{Name: "Sữa tươi", Quantity: 150, Unit: "ml"},
						{Name: "Đường", Quantity: 30, Unit: "g"},
						{Name: "Đá", Quantity: 75, Unit: "g"},
					},
				},
			},
			CreatedAt: time.Now(),
			UpdatedAt: time.Now(),
		},

		// ========== SINGLE-SIZE ITEMS ==========
		MenuItem{
			Name:        "Bánh mì thịt",
			Category:    "Món ăn",
			Description: "Bánh mì Việt Nam truyền thống",
			Available:   true,
			HasVariants: false,
			Price:       20000,
			Ingredients: []Ingredient{
				{Name: "Bánh mì", Quantity: 1, Unit: "cái"},
				{Name: "Thịt", Quantity: 50, Unit: "g"},
			},
			CreatedAt: time.Now(),
			UpdatedAt: time.Now(),
		},

		MenuItem{
			Name:        "Bánh tiramisu",
			Category:    "Bánh ngọt",
			Description: "Bánh tiramisu Ý truyền thống",
			Available:   true,
			HasVariants: false,
			Price:       45000,
			Ingredients: []Ingredient{
				{Name: "Bánh tiramisu", Quantity: 1, Unit: "miếng"},
			},
			CreatedAt: time.Now(),
			UpdatedAt: time.Now(),
		},

		MenuItem{
			Name:        "Bánh croissant",
			Category:    "Bánh ngọt",
			Description: "Bánh sừng bò Pháp giòn tan",
			Available:   true,
			HasVariants: false,
			Price:       35000,
			Ingredients: []Ingredient{
				{Name: "Bánh croissant", Quantity: 1, Unit: "cái"},
			},
			CreatedAt: time.Now(),
			UpdatedAt: time.Now(),
		},

		MenuItem{
			Name:        "Nước ép cam",
			Category:    "Nước ép",
			Description: "Nước cam tươi 100%",
			Available:   true,
			HasVariants: false,
			Price:       35000,
			Ingredients: []Ingredient{
				{Name: "Cam", Quantity: 200, Unit: "g"},
				{Name: "Đường", Quantity: 10, Unit: "g"},
			},
			CreatedAt: time.Now(),
			UpdatedAt: time.Now(),
		},
	}

	// Check if menu items already exist
	count, err := collection.CountDocuments(ctx, bson.M{})
	if err != nil {
		log.Fatal(err)
	}

	if count > 0 {
		fmt.Printf("⚠️  Found %d existing menu items\n", count)
		fmt.Print("Do you want to clear and reseed? (y/N): ")
		var response string
		fmt.Scanln(&response)
		if response != "y" && response != "Y" {
			fmt.Println("❌ Cancelled")
			return
		}
		// Clear existing items
		_, err = collection.DeleteMany(ctx, bson.M{})
		if err != nil {
			log.Fatal(err)
		}
		fmt.Println("🗑️  Cleared existing menu items")
	}

	// Insert menu items
	result, err := collection.InsertMany(ctx, menuItems)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Printf("\n✅ Seeded %d menu items successfully!\n\n", len(result.InsertedIDs))
	fmt.Println("📋 Menu Items:")
	fmt.Println("\n🔄 Multi-size items (with variants):")
	fmt.Println("  1. Cà phê sữa đá (M: 25k, L: 30k, XL: 35k)")
	fmt.Println("  2. Cà phê đen đá (M: 20k, L: 25k)")
	fmt.Println("  3. Trà sữa truyền thống (M: 35k, L: 42k, XL: 48k)")
	fmt.Println("  4. Trà sữa matcha (M: 40k, L: 48k)")
	fmt.Println("  5. Trà đào cam sả (M: 35k, L: 42k)")
	fmt.Println("  6. Sinh tố bơ (M: 40k, L: 50k)")
	fmt.Println("\n📦 Single-size items:")
	fmt.Println("  7. Bánh mì thịt (20k)")
	fmt.Println("  8. Bánh tiramisu (45k)")
	fmt.Println("  9. Bánh croissant (35k)")
	fmt.Println("  10. Nước ép cam (35k)")
	fmt.Println("\n💰 Price range: 20,000đ - 50,000đ")
	fmt.Println("\n📊 Next steps:")
	fmt.Println("  1. Run: go run backend/cmd/seed-menu-variants/main.go")
	fmt.Println("  2. Calculate costs: POST /api/menu/:id/calculate-cost")
	fmt.Println("  3. View cost breakdown: GET /api/menu/:id/cost-breakdown")
	fmt.Println("  4. View profit analysis: GET /api/menu/:id/profit-analysis")
	fmt.Println("\n🎉 Ready to test cost analysis!")

	// Print inserted IDs for reference
	fmt.Println("\n📝 Inserted Menu Item IDs:")
	for i, id := range result.InsertedIDs {
		fmt.Printf("  %d. %v\n", i+1, id)
	}
}
