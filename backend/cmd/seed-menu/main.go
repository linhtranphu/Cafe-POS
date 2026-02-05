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

type MenuItem struct {
	Name        string       `bson:"name"`
	Price       float64      `bson:"price"`
	Category    string       `bson:"category"`
	Description string       `bson:"description"`
	Ingredients []Ingredient `bson:"ingredients"`
	Available   bool         `bson:"available"`
	CreatedAt   time.Time    `bson:"created_at"`
	UpdatedAt   time.Time    `bson:"updated_at"`
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

	// Menu items to seed
	menuItems := []interface{}{
		// ========== CÀ PHÊ ==========
		MenuItem{
			Name:        "Cà phê đen",
			Price:       25000,
			Category:    "Cà phê",
			Description: "Cà phê phin truyền thống, đậm đà",
			Ingredients: []Ingredient{
				{Name: "Cà phê", Quantity: 20, Unit: "g"},
				{Name: "Nước nóng", Quantity: 100, Unit: "ml"},
			},
			Available: true,
			CreatedAt: time.Now(),
			UpdatedAt: time.Now(),
		},
		MenuItem{
			Name:        "Cà phê sữa",
			Price:       30000,
			Category:    "Cà phê",
			Description: "Cà phê phin với sữa đặc ngọt ngào",
			Ingredients: []Ingredient{
				{Name: "Cà phê", Quantity: 20, Unit: "g"},
				{Name: "Sữa đặc", Quantity: 30, Unit: "ml"},
				{Name: "Nước nóng", Quantity: 100, Unit: "ml"},
			},
			Available: true,
			CreatedAt: time.Now(),
			UpdatedAt: time.Now(),
		},
		MenuItem{
			Name:        "Bạc xỉu",
			Price:       32000,
			Category:    "Cà phê",
			Description: "Cà phê sữa nhiều sữa, ít cà phê",
			Ingredients: []Ingredient{
				{Name: "Cà phê", Quantity: 15, Unit: "g"},
				{Name: "Sữa đặc", Quantity: 50, Unit: "ml"},
				{Name: "Nước nóng", Quantity: 100, Unit: "ml"},
			},
			Available: true,
			CreatedAt: time.Now(),
			UpdatedAt: time.Now(),
		},
		MenuItem{
			Name:        "Cà phê đá",
			Price:       28000,
			Category:    "Cà phê",
			Description: "Cà phê đen mát lạnh",
			Ingredients: []Ingredient{
				{Name: "Cà phê", Quantity: 20, Unit: "g"},
				{Name: "Nước nóng", Quantity: 80, Unit: "ml"},
				{Name: "Đá", Quantity: 100, Unit: "g"},
			},
			Available: true,
			CreatedAt: time.Now(),
			UpdatedAt: time.Now(),
		},
		MenuItem{
			Name:        "Cà phê sữa đá",
			Price:       32000,
			Category:    "Cà phê",
			Description: "Cà phê sữa mát lạnh",
			Ingredients: []Ingredient{
				{Name: "Cà phê", Quantity: 20, Unit: "g"},
				{Name: "Sữa đặc", Quantity: 30, Unit: "ml"},
				{Name: "Nước nóng", Quantity: 80, Unit: "ml"},
				{Name: "Đá", Quantity: 100, Unit: "g"},
			},
			Available: true,
			CreatedAt: time.Now(),
			UpdatedAt: time.Now(),
		},

		// ========== TRÀ SỮA ==========
		MenuItem{
			Name:        "Trà sữa truyền thống",
			Price:       35000,
			Category:    "Trà sữa",
			Description: "Trà sữa đài loan cổ điển",
			Ingredients: []Ingredient{
				{Name: "Trà đen", Quantity: 10, Unit: "g"},
				{Name: "Sữa tươi", Quantity: 100, Unit: "ml"},
				{Name: "Đường", Quantity: 20, Unit: "g"},
				{Name: "Trân châu", Quantity: 50, Unit: "g"},
			},
			Available: true,
			CreatedAt: time.Now(),
			UpdatedAt: time.Now(),
		},
		MenuItem{
			Name:        "Trà sữa matcha",
			Price:       40000,
			Category:    "Trà sữa",
			Description: "Trà xanh matcha Nhật Bản với sữa",
			Ingredients: []Ingredient{
				{Name: "Matcha", Quantity: 5, Unit: "g"},
				{Name: "Sữa tươi", Quantity: 100, Unit: "ml"},
				{Name: "Đường", Quantity: 20, Unit: "g"},
				{Name: "Trân châu", Quantity: 50, Unit: "g"},
			},
			Available: true,
			CreatedAt: time.Now(),
			UpdatedAt: time.Now(),
		},
		MenuItem{
			Name:        "Trà sữa socola",
			Price:       38000,
			Category:    "Trà sữa",
			Description: "Trà sữa vị socola đậm đà",
			Ingredients: []Ingredient{
				{Name: "Trà đen", Quantity: 10, Unit: "g"},
				{Name: "Sữa tươi", Quantity: 100, Unit: "ml"},
				{Name: "Bột socola", Quantity: 15, Unit: "g"},
				{Name: "Đường", Quantity: 20, Unit: "g"},
				{Name: "Trân châu", Quantity: 50, Unit: "g"},
			},
			Available: true,
			CreatedAt: time.Now(),
			UpdatedAt: time.Now(),
		},

		// ========== TRÀ TRÁI CÂY ==========
		MenuItem{
			Name:        "Trà đào cam sả",
			Price:       35000,
			Category:    "Trà trái cây",
			Description: "Trà đào tươi mát với cam và sả thơm",
			Ingredients: []Ingredient{
				{Name: "Trà xanh", Quantity: 10, Unit: "g"},
				{Name: "Đào", Quantity: 50, Unit: "g"},
				{Name: "Cam", Quantity: 30, Unit: "g"},
				{Name: "Sả", Quantity: 5, Unit: "g"},
				{Name: "Đường", Quantity: 20, Unit: "g"},
			},
			Available: true,
			CreatedAt: time.Now(),
			UpdatedAt: time.Now(),
		},
		MenuItem{
			Name:        "Trà chanh leo",
			Price:       32000,
			Category:    "Trà trái cây",
			Description: "Trà xanh với chanh leo chua ngọt",
			Ingredients: []Ingredient{
				{Name: "Trà xanh", Quantity: 10, Unit: "g"},
				{Name: "Chanh leo", Quantity: 50, Unit: "g"},
				{Name: "Đường", Quantity: 25, Unit: "g"},
			},
			Available: true,
			CreatedAt: time.Now(),
			UpdatedAt: time.Now(),
		},
		MenuItem{
			Name:        "Trà vải",
			Price:       33000,
			Category:    "Trà trái cây",
			Description: "Trà xanh với vải tươi ngọt mát",
			Ingredients: []Ingredient{
				{Name: "Trà xanh", Quantity: 10, Unit: "g"},
				{Name: "Vải", Quantity: 60, Unit: "g"},
				{Name: "Đường", Quantity: 20, Unit: "g"},
			},
			Available: true,
			CreatedAt: time.Now(),
			UpdatedAt: time.Now(),
		},

		// ========== SINH TỐ ==========
		MenuItem{
			Name:        "Sinh tố bơ",
			Price:       40000,
			Category:    "Sinh tố",
			Description: "Sinh tố bơ béo ngậy",
			Ingredients: []Ingredient{
				{Name: "Bơ", Quantity: 150, Unit: "g"},
				{Name: "Sữa tươi", Quantity: 100, Unit: "ml"},
				{Name: "Đường", Quantity: 20, Unit: "g"},
				{Name: "Đá", Quantity: 50, Unit: "g"},
			},
			Available: true,
			CreatedAt: time.Now(),
			UpdatedAt: time.Now(),
		},
		MenuItem{
			Name:        "Sinh tố dâu",
			Price:       38000,
			Category:    "Sinh tố",
			Description: "Sinh tố dâu tây tươi mát",
			Ingredients: []Ingredient{
				{Name: "Dâu tây", Quantity: 100, Unit: "g"},
				{Name: "Sữa tươi", Quantity: 100, Unit: "ml"},
				{Name: "Đường", Quantity: 20, Unit: "g"},
				{Name: "Đá", Quantity: 50, Unit: "g"},
			},
			Available: true,
			CreatedAt: time.Now(),
			UpdatedAt: time.Now(),
		},
		MenuItem{
			Name:        "Sinh tố xoài",
			Price:       38000,
			Category:    "Sinh tố",
			Description: "Sinh tố xoài ngọt thơm",
			Ingredients: []Ingredient{
				{Name: "Xoài", Quantity: 150, Unit: "g"},
				{Name: "Sữa tươi", Quantity: 100, Unit: "ml"},
				{Name: "Đường", Quantity: 15, Unit: "g"},
				{Name: "Đá", Quantity: 50, Unit: "g"},
			},
			Available: true,
			CreatedAt: time.Now(),
			UpdatedAt: time.Now(),
		},

		// ========== NƯỚC ÉP ==========
		MenuItem{
			Name:        "Nước ép cam",
			Price:       35000,
			Category:    "Nước ép",
			Description: "Nước cam tươi 100%",
			Ingredients: []Ingredient{
				{Name: "Cam", Quantity: 200, Unit: "g"},
				{Name: "Đường", Quantity: 10, Unit: "g"},
			},
			Available: true,
			CreatedAt: time.Now(),
			UpdatedAt: time.Now(),
		},
		MenuItem{
			Name:        "Nước ép dưa hấu",
			Price:       30000,
			Category:    "Nước ép",
			Description: "Nước dưa hấu mát lạnh",
			Ingredients: []Ingredient{
				{Name: "Dưa hấu", Quantity: 250, Unit: "g"},
				{Name: "Đường", Quantity: 5, Unit: "g"},
			},
			Available: true,
			CreatedAt: time.Now(),
			UpdatedAt: time.Now(),
		},

		// ========== BÁNH NGỌT ==========
		MenuItem{
			Name:        "Bánh tiramisu",
			Price:       45000,
			Category:    "Bánh ngọt",
			Description: "Bánh tiramisu Ý truyền thống",
			Ingredients: []Ingredient{
				{Name: "Bánh tiramisu", Quantity: 1, Unit: "miếng"},
			},
			Available: true,
			CreatedAt: time.Now(),
			UpdatedAt: time.Now(),
		},
		MenuItem{
			Name:        "Bánh cheesecake",
			Price:       42000,
			Category:    "Bánh ngọt",
			Description: "Bánh phô mai mềm mịn",
			Ingredients: []Ingredient{
				{Name: "Bánh cheesecake", Quantity: 1, Unit: "miếng"},
			},
			Available: true,
			CreatedAt: time.Now(),
			UpdatedAt: time.Now(),
		},
		MenuItem{
			Name:        "Bánh croissant",
			Price:       35000,
			Category:    "Bánh ngọt",
			Description: "Bánh sừng bò Pháp giòn tan",
			Ingredients: []Ingredient{
				{Name: "Bánh croissant", Quantity: 1, Unit: "cái"},
			},
			Available: true,
			CreatedAt: time.Now(),
			UpdatedAt: time.Now(),
		},
		MenuItem{
			Name:        "Bánh muffin",
			Price:       30000,
			Category:    "Bánh ngọt",
			Description: "Bánh muffin chocolate chip",
			Ingredients: []Ingredient{
				{Name: "Bánh muffin", Quantity: 1, Unit: "cái"},
			},
			Available: true,
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
	fmt.Println("📋 Menu Categories:")
	fmt.Println("- Cà phê: 5 items")
	fmt.Println("- Trà sữa: 3 items")
	fmt.Println("- Trà trái cây: 3 items")
	fmt.Println("- Sinh tố: 3 items")
	fmt.Println("- Nước ép: 2 items")
	fmt.Println("- Bánh ngọt: 4 items")
	fmt.Println("\n💰 Price range: 25,000đ - 45,000đ")
	fmt.Println("\n🎉 Ready to create orders!")
}
