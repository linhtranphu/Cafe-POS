package main

import (
	"context"
	"log"
	"time"
	"cafe-pos/backend/application/services"
	"cafe-pos/backend/domain/facility"
	"cafe-pos/backend/domain/ingredient"
	"cafe-pos/backend/domain/menu"
	"cafe-pos/backend/domain/user"
	"cafe-pos/backend/infrastructure/mongodb"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

func seedData() {
	// MongoDB connection
	client, err := mongo.Connect(context.TODO(), options.Client().ApplyURI("mongodb://localhost:27017"))
	if err != nil {
		log.Fatal(err)
	}
	defer client.Disconnect(context.TODO())

	db := client.Database("cafe_pos")
	ctx := context.Background()

	// Repositories
	userRepo := mongodb.NewUserRepository(db)
	menuRepo := mongodb.NewMenuRepository(db)

	// Services
	authService := services.NewAuthService(userRepo, services.NewJWTService("secret"))

	// Seed Users
	seedUsers(ctx, authService, userRepo)
	
	// Seed Menu Items
	seedMenuItems(ctx, menuRepo)
	
	// Seed Ingredients
	ingredientRepo := mongodb.NewIngredientRepository(db)
	seedIngredients(ctx, ingredientRepo)
	
	// Seed Facilities
	facilityRepo := mongodb.NewFacilityRepository(db)
	seedFacilities(ctx, facilityRepo)

	log.Println("✅ Seed data completed!")
}

func seedUsers(ctx context.Context, authService *services.AuthService, userRepo *mongodb.UserRepository) {
	// Only seed manager account for production
	users := []struct {
		username string
		password string
		role     user.Role
		name     string
	}{
		{"admin", "admin123", user.RoleManager, "Quản lý"},
	}

	for _, u := range users {
		// Check if user exists
		if _, err := userRepo.FindByUsername(ctx, u.username); err == nil {
			continue
		}

		hashedPassword, _ := authService.HashPassword(u.password)
		newUser := &user.User{
			Username:  u.username,
			Password:  hashedPassword,
			Role:      u.role,
			Name:      u.name,
			Active:    true,
			CreatedAt: time.Now(),
			UpdatedAt: time.Now(),
		}
		userRepo.Create(ctx, newUser)
		log.Printf("Created user: %s (%s)", u.username, u.role)
	}
}

func seedIngredients(ctx context.Context, ingredientRepo *mongodb.IngredientRepository) {
	ingredients := []ingredient.Ingredient{
		// Cà phê
		{Name: "Cà phê bột Robusta", Category: "Cà phê", Unit: "kg", Quantity: 5.0, MinStock: 2.0, CostPerUnit: 150000, Supplier: "Trung Nguyên"},
		{Name: "Cà phê bột Arabica", Category: "Cà phê", Unit: "kg", Quantity: 3.0, MinStock: 1.5, CostPerUnit: 200000, Supplier: "Highland Coffee"},
		{Name: "Espresso blend", Category: "Cà phê", Unit: "kg", Quantity: 2.5, MinStock: 1.0, CostPerUnit: 180000, Supplier: "Lavazza"},
		
		// Sữa
		{Name: "Sữa tươi", Category: "Sữa", Unit: "L", Quantity: 20.0, MinStock: 10.0, CostPerUnit: 25000, Supplier: "Vinamilk"},
		{Name: "Sữa đặc", Category: "Sữa", Unit: "ml", Quantity: 2000.0, MinStock: 500.0, CostPerUnit: 80, Supplier: "Ông Thọ"},
		{Name: "Sữa hạnh nhân", Category: "Sữa", Unit: "L", Quantity: 5.0, MinStock: 2.0, CostPerUnit: 45000, Supplier: "Alpro"},
		
		// Trà
		{Name: "Trà đen Ceylon", Category: "Trà", Unit: "g", Quantity: 500.0, MinStock: 200.0, CostPerUnit: 200, Supplier: "Lipton"},
		{Name: "Trà xanh Thái Nguyên", Category: "Trà", Unit: "g", Quantity: 300.0, MinStock: 100.0, CostPerUnit: 300, Supplier: "Phúc Long"},
		{Name: "Trà Earl Grey", Category: "Trà", Unit: "g", Quantity: 200.0, MinStock: 50.0, CostPerUnit: 400, Supplier: "Twinings"},
		
		// Đường
		{Name: "Đường trắng", Category: "Đường", Unit: "kg", Quantity: 10.0, MinStock: 5.0, CostPerUnit: 18000, Supplier: "Biên Hòa"},
		{Name: "Đường nâu", Category: "Đường", Unit: "kg", Quantity: 3.0, MinStock: 1.0, CostPerUnit: 25000, Supplier: "Organic"},
		{Name: "Mật ong", Category: "Đường", Unit: "ml", Quantity: 1000.0, MinStock: 300.0, CostPerUnit: 150, Supplier: "Honeywell"},
		
		// Trái cây
		{Name: "Chanh tươi", Category: "Trái cây", Unit: "piece", Quantity: 50.0, MinStock: 20.0, CostPerUnit: 3000, Supplier: "Đà Lạt Farm"},
		{Name: "Cam tươi", Category: "Trái cây", Unit: "piece", Quantity: 30.0, MinStock: 15.0, CostPerUnit: 5000, Supplier: "Đà Lạt Farm"},
		{Name: "Đào đóng hộp", Category: "Trái cây", Unit: "g", Quantity: 800.0, MinStock: 200.0, CostPerUnit: 50, Supplier: "Del Monte"},
		{Name: "Dâu tây đông lạnh", Category: "Trái cây", Unit: "g", Quantity: 1000.0, MinStock: 300.0, CostPerUnit: 80, Supplier: "Đà Lạt"},
		
		// Bánh
		{Name: "Bánh croissant đông lạnh", Category: "Bánh", Unit: "piece", Quantity: 24.0, MinStock: 10.0, CostPerUnit: 8000, Supplier: "Tous Les Jours"},
		{Name: "Bánh mì sandwich", Category: "Bánh", Unit: "piece", Quantity: 20.0, MinStock: 8.0, CostPerUnit: 5000, Supplier: "ABC Bakery"},
		{Name: "Bánh quy digestive", Category: "Bánh", Unit: "pack", Quantity: 15.0, MinStock: 5.0, CostPerUnit: 25000, Supplier: "McVitie's"},
		
		// Khác
		{Name: "Ly giấy 8oz", Category: "Khác", Unit: "piece", Quantity: 500.0, MinStock: 100.0, CostPerUnit: 800, Supplier: "Bao bì Sài Gòn"},
		{Name: "Ly giấy 12oz", Category: "Khác", Unit: "piece", Quantity: 300.0, MinStock: 80.0, CostPerUnit: 1000, Supplier: "Bao bì Sài Gòn"},
		{Name: "Nắp ly", Category: "Khác", Unit: "piece", Quantity: 800.0, MinStock: 200.0, CostPerUnit: 300, Supplier: "Bao bì Sài Gòn"},
		{Name: "Ống hút giấy", Category: "Khác", Unit: "piece", Quantity: 1000.0, MinStock: 300.0, CostPerUnit: 200, Supplier: "Eco Straw"},
		{Name: "Khăn giấy", Category: "Khác", Unit: "pack", Quantity: 20.0, MinStock: 5.0, CostPerUnit: 15000, Supplier: "Saigon Paper"},
	}

	for _, ing := range ingredients {
		ing.CreatedAt = time.Now()
		ing.UpdatedAt = time.Now()
		ingredientRepo.Create(ctx, &ing)
		log.Printf("Created ingredient: %s - %s (%v %s)", ing.Name, ing.Category, ing.Quantity, ing.Unit)
	}
}
	func seedMenuItems(ctx context.Context, menuRepo *mongodb.MenuRepository) {
	menuItems := []menu.MenuItem{
		// Cà phê
		{Name: "Cà phê đen", Price: 25000, Category: "Cà phê", Description: "Cà phê đen truyền thống", Available: true,
			Ingredients: []menu.Ingredient{
				{Name: "Cà phê bột", Quantity: 20, Unit: "gram"},
				{Name: "Nước", Quantity: 150, Unit: "ml"},
			}},
		{Name: "Cà phê sữa", Price: 30000, Category: "Cà phê", Description: "Cà phê sữa đá", Available: true,
			Ingredients: []menu.Ingredient{
				{Name: "Cà phê bột", Quantity: 20, Unit: "gram"},
				{Name: "Sữa đặc", Quantity: 30, Unit: "ml"},
				{Name: "Nước", Quantity: 120, Unit: "ml"},
				{Name: "Đá", Quantity: 100, Unit: "gram"},
			}},
		{Name: "Cappuccino", Price: 45000, Category: "Cà phê", Description: "Cappuccino Ý", Available: true,
			Ingredients: []menu.Ingredient{
				{Name: "Espresso", Quantity: 30, Unit: "ml"},
				{Name: "Sữa tươi", Quantity: 150, Unit: "ml"},
			}},
		{Name: "Latte", Price: 50000, Category: "Cà phê", Description: "Latte với foam art", Available: true,
			Ingredients: []menu.Ingredient{
				{Name: "Espresso", Quantity: 30, Unit: "ml"},
				{Name: "Sữa tươi", Quantity: 200, Unit: "ml"},
			}},
		{Name: "Americano", Price: 35000, Category: "Cà phê", Description: "Americano đậm đà", Available: true,
			Ingredients: []menu.Ingredient{
				{Name: "Espresso", Quantity: 60, Unit: "ml"},
				{Name: "Nước nóng", Quantity: 120, Unit: "ml"},
			}},
		
		// Trà
		{Name: "Trà đào", Price: 35000, Category: "Trà", Description: "Trà đào cam sả", Available: true,
			Ingredients: []menu.Ingredient{
				{Name: "Trà xanh", Quantity: 5, Unit: "gram"},
				{Name: "Đào", Quantity: 50, Unit: "gram"},
				{Name: "Cam", Quantity: 30, Unit: "gram"},
				{Name: "Sả", Quantity: 10, Unit: "gram"},
				{Name: "Nước", Quantity: 300, Unit: "ml"},
			}},
		{Name: "Trà sữa", Price: 40000, Category: "Trà", Description: "Trà sữa trân châu", Available: true,
			Ingredients: []menu.Ingredient{
				{Name: "Trà đen", Quantity: 8, Unit: "gram"},
				{Name: "Sữa tươi", Quantity: 100, Unit: "ml"},
				{Name: "Trân châu", Quantity: 30, Unit: "gram"},
				{Name: "Đường", Quantity: 20, Unit: "gram"},
			}},
		{Name: "Trà xanh", Price: 25000, Category: "Trà", Description: "Trà xanh Thái Nguyên", Available: true,
			Ingredients: []menu.Ingredient{
				{Name: "Trà xanh", Quantity: 5, Unit: "gram"},
				{Name: "Nước nóng", Quantity: 200, Unit: "ml"},
			}},
		{Name: "Trà chanh", Price: 30000, Category: "Trà", Description: "Trà chanh mật ong", Available: true,
			Ingredients: []menu.Ingredient{
				{Name: "Trà đen", Quantity: 5, Unit: "gram"},
				{Name: "Chanh", Quantity: 40, Unit: "gram"},
				{Name: "Mật ong", Quantity: 15, Unit: "ml"},
				{Name: "Nước", Quantity: 250, Unit: "ml"},
			}},
		
		// Nước ép
		{Name: "Nước ép cam", Price: 35000, Category: "Nước ép", Description: "Cam tươi vắt", Available: true,
			Ingredients: []menu.Ingredient{
				{Name: "Cam", Quantity: 200, Unit: "gram"},
				{Name: "Đường", Quantity: 10, Unit: "gram"},
			}},
		{Name: "Sinh tố bơ", Price: 45000, Category: "Nước ép", Description: "Sinh tố bơ sữa", Available: true,
			Ingredients: []menu.Ingredient{
				{Name: "Bơ", Quantity: 150, Unit: "gram"},
				{Name: "Sữa tươi", Quantity: 200, Unit: "ml"},
				{Name: "Đường", Quantity: 15, Unit: "gram"},
				{Name: "Đá", Quantity: 50, Unit: "gram"},
			}},
		{Name: "Nước ép dưa hấu", Price: 30000, Category: "Nước ép", Description: "Dưa hấu tươi", Available: true,
			Ingredients: []menu.Ingredient{
				{Name: "Dưa hấu", Quantity: 300, Unit: "gram"},
				{Name: "Đá", Quantity: 50, Unit: "gram"},
			}},
		
		// Bánh ngọt
		{Name: "Bánh tiramisu", Price: 55000, Category: "Bánh ngọt", Description: "Tiramisu Ý", Available: true,
			Ingredients: []menu.Ingredient{
				{Name: "Bánh ladyfinger", Quantity: 6, Unit: "cái"},
				{Name: "Mascarpone", Quantity: 100, Unit: "gram"},
				{Name: "Cà phê espresso", Quantity: 50, Unit: "ml"},
				{Name: "Bột cacao", Quantity: 5, Unit: "gram"},
			}},
		{Name: "Bánh cheesecake", Price: 50000, Category: "Bánh ngọt", Description: "Cheesecake dâu", Available: true,
			Ingredients: []menu.Ingredient{
				{Name: "Cream cheese", Quantity: 120, Unit: "gram"},
				{Name: "Bánh quy", Quantity: 50, Unit: "gram"},
				{Name: "Dâu tây", Quantity: 30, Unit: "gram"},
				{Name: "Đường", Quantity: 20, Unit: "gram"},
			}},
		{Name: "Bánh croissant", Price: 35000, Category: "Bánh ngọt", Description: "Croissant bơ", Available: true,
			Ingredients: []menu.Ingredient{
				{Name: "Bột mì", Quantity: 80, Unit: "gram"},
				{Name: "Bơ", Quantity: 40, Unit: "gram"},
				{Name: "Trứng", Quantity: 1, Unit: "quả"},
			}},
		
		// Món nhẹ
		{Name: "Sandwich", Price: 45000, Category: "Món nhẹ", Description: "Sandwich gà", Available: true,
			Ingredients: []menu.Ingredient{
				{Name: "Bánh mì", Quantity: 2, Unit: "lát"},
				{Name: "Thịt gà", Quantity: 80, Unit: "gram"},
				{Name: "Rau xà lách", Quantity: 20, Unit: "gram"},
				{Name: "Cà chua", Quantity: 30, Unit: "gram"},
				{Name: "Mayonnaise", Quantity: 10, Unit: "gram"},
			}},
		{Name: "Salad", Price: 40000, Category: "Món nhẹ", Description: "Salad rau củ", Available: true,
			Ingredients: []menu.Ingredient{
				{Name: "Rau xà lách", Quantity: 50, Unit: "gram"},
				{Name: "Cà chua", Quantity: 40, Unit: "gram"},
				{Name: "Dưa chuột", Quantity: 30, Unit: "gram"},
				{Name: "Dầu olive", Quantity: 10, Unit: "ml"},
			}},
		{Name: "Khoai tây chiên", Price: 30000, Category: "Món nhẹ", Description: "Khoai tây chiên giòn", Available: true,
			Ingredients: []menu.Ingredient{
				{Name: "Khoai tây", Quantity: 200, Unit: "gram"},
				{Name: "Dầu ăn", Quantity: 50, Unit: "ml"},
				{Name: "Muối", Quantity: 2, Unit: "gram"},
			}},
	}

	for _, item := range menuItems {
		item.CreatedAt = time.Now()
		item.UpdatedAt = time.Now()
		menuRepo.Create(ctx, &item)
		log.Printf("Created menu item: %s - %s", item.Name, item.Category)
	}
}

func main() {
	seedData()
}

func seedFacilities(ctx context.Context, facilityRepo *mongodb.FacilityRepository) {
	facilities := []facility.Facility{
		// Bàn ghế
		{Name: "Bàn khách 2 chỗ", Type: facility.TypeFurniture, Area: facility.AreaDiningRoom, Quantity: 8, Status: facility.StatusInUse, PurchaseDate: time.Date(2023, 1, 15, 0, 0, 0, 0, time.UTC), Cost: 1200000, Supplier: "Nội thất Hòa Phát", Notes: "Bàn gỗ tự nhiên, màu nâu"},
		{Name: "Bàn khách 4 chỗ", Type: facility.TypeFurniture, Area: facility.AreaDiningRoom, Quantity: 6, Status: facility.StatusInUse, PurchaseDate: time.Date(2023, 1, 15, 0, 0, 0, 0, time.UTC), Cost: 1800000, Supplier: "Nội thất Hòa Phát", Notes: "Bàn gỗ tự nhiên, màu nâu"},
		{Name: "Ghế cafe", Type: facility.TypeFurniture, Area: facility.AreaDiningRoom, Quantity: 32, Status: facility.StatusInUse, PurchaseDate: time.Date(2023, 1, 15, 0, 0, 0, 0, time.UTC), Cost: 450000, Supplier: "Nội thất Hòa Phát", Notes: "Ghế gỗ có đệm"},
		{Name: "Sofa góc", Type: facility.TypeFurniture, Area: facility.AreaDiningRoom, Quantity: 2, Status: facility.StatusInUse, PurchaseDate: time.Date(2023, 2, 10, 0, 0, 0, 0, time.UTC), Cost: 3500000, Supplier: "Sofa House", Notes: "Sofa da màu nâu"},
		{Name: "Bàn bar cao", Type: facility.TypeFurniture, Area: facility.AreaCounter, Quantity: 4, Status: facility.StatusInUse, PurchaseDate: time.Date(2023, 3, 5, 0, 0, 0, 0, time.UTC), Cost: 800000, Supplier: "Bar Furniture", Notes: "Bàn cao 110cm"},
		{Name: "Ghế bar cao", Type: facility.TypeFurniture, Area: facility.AreaCounter, Quantity: 12, Status: facility.StatusInUse, PurchaseDate: time.Date(2023, 3, 5, 0, 0, 0, 0, time.UTC), Cost: 350000, Supplier: "Bar Furniture", Notes: "Ghế xoay có tựa lưng"},
		
		// Máy móc
		{Name: "Máy pha cà phê Espresso", Type: facility.TypeMachine, Area: facility.AreaCounter, Quantity: 1, Status: facility.StatusInUse, PurchaseDate: time.Date(2023, 1, 20, 0, 0, 0, 0, time.UTC), Cost: 25000000, Supplier: "La Marzocco Vietnam", Notes: "Máy pha chuyên nghiệp 2 group"},
		{Name: "Máy xay cà phê", Type: facility.TypeMachine, Area: facility.AreaCounter, Quantity: 2, Status: facility.StatusInUse, PurchaseDate: time.Date(2023, 1, 20, 0, 0, 0, 0, time.UTC), Cost: 8000000, Supplier: "Mazzer Vietnam", Notes: "Máy xay burr chuyên nghiệp"},
		{Name: "Tủ lạnh 2 cánh", Type: facility.TypeMachine, Area: facility.AreaKitchen, Quantity: 1, Status: facility.StatusInUse, PurchaseDate: time.Date(2023, 1, 10, 0, 0, 0, 0, time.UTC), Cost: 12000000, Supplier: "Samsung", Notes: "Tủ lạnh 500L"},
		{Name: "Tủ đông", Type: facility.TypeMachine, Area: facility.AreaKitchen, Quantity: 1, Status: facility.StatusInUse, PurchaseDate: time.Date(2023, 1, 10, 0, 0, 0, 0, time.UTC), Cost: 8000000, Supplier: "LG", Notes: "Tủ đông 300L"},
		{Name: "Máy làm đá", Type: facility.TypeMachine, Area: facility.AreaCounter, Quantity: 1, Status: facility.StatusInUse, PurchaseDate: time.Date(2023, 2, 1, 0, 0, 0, 0, time.UTC), Cost: 15000000, Supplier: "Hoshizaki", Notes: "Máy làm đá viên 50kg/ngày"},
		{Name: "Máy ép trái cây", Type: facility.TypeMachine, Area: facility.AreaCounter, Quantity: 1, Status: facility.StatusInUse, PurchaseDate: time.Date(2023, 2, 15, 0, 0, 0, 0, time.UTC), Cost: 3500000, Supplier: "Zumex", Notes: "Máy ép cam tự động"},
		{Name: "Lò vi sóng", Type: facility.TypeMachine, Area: facility.AreaKitchen, Quantity: 1, Status: facility.StatusBroken, PurchaseDate: time.Date(2022, 8, 10, 0, 0, 0, 0, time.UTC), Cost: 2500000, Supplier: "Panasonic", Notes: "Cần sửa chữa bảng điều khiển"},
		
		// Dụng cụ
		{Name: "Ly thủy tinh 250ml", Type: facility.TypeUtensil, Area: facility.AreaCounter, Quantity: 50, Status: facility.StatusInUse, PurchaseDate: time.Date(2023, 1, 25, 0, 0, 0, 0, time.UTC), Cost: 25000, Supplier: "Ocean Glass", Notes: "Ly thủy tinh trong suốt"},
		{Name: "Ly thủy tinh 350ml", Type: facility.TypeUtensil, Area: facility.AreaCounter, Quantity: 40, Status: facility.StatusInUse, PurchaseDate: time.Date(2023, 1, 25, 0, 0, 0, 0, time.UTC), Cost: 30000, Supplier: "Ocean Glass", Notes: "Ly thủy tinh trong suốt"},
		{Name: "Tách cà phê ceramic", Type: facility.TypeUtensil, Area: facility.AreaCounter, Quantity: 30, Status: facility.StatusInUse, PurchaseDate: time.Date(2023, 1, 25, 0, 0, 0, 0, time.UTC), Cost: 45000, Supplier: "Minh Long", Notes: "Tách sứ trắng có đĩa lót"},
		{Name: "Thìa inox", Type: facility.TypeUtensil, Area: facility.AreaCounter, Quantity: 60, Status: facility.StatusInUse, PurchaseDate: time.Date(2023, 1, 25, 0, 0, 0, 0, time.UTC), Cost: 8000, Supplier: "Inox Việt", Notes: "Thìa cà phê inox 304"},
		{Name: "Bình đựng sữa inox", Type: facility.TypeUtensil, Area: facility.AreaCounter, Quantity: 4, Status: facility.StatusInUse, PurchaseDate: time.Date(2023, 2, 5, 0, 0, 0, 0, time.UTC), Cost: 150000, Supplier: "Barista Tools", Notes: "Bình đựng sữa 600ml"},
		{Name: "Tamper cà phê", Type: facility.TypeUtensil, Area: facility.AreaCounter, Quantity: 2, Status: facility.StatusInUse, PurchaseDate: time.Date(2023, 2, 5, 0, 0, 0, 0, time.UTC), Cost: 350000, Supplier: "Barista Tools", Notes: "Tamper 58mm"},
		{Name: "Khăn lau microfiber", Type: facility.TypeUtensil, Area: facility.AreaCounter, Quantity: 20, Status: facility.StatusInUse, PurchaseDate: time.Date(2023, 3, 1, 0, 0, 0, 0, time.UTC), Cost: 15000, Supplier: "Clean Pro", Notes: "Khăn lau chuyên dụng"},
		
		// Điện tử
		{Name: "Máy tính tiền", Type: facility.TypeElectric, Area: facility.AreaCounter, Quantity: 1, Status: facility.StatusInUse, PurchaseDate: time.Date(2023, 1, 5, 0, 0, 0, 0, time.UTC), Cost: 8000000, Supplier: "Casio", Notes: "Máy tính tiền điện tử có màn hình"},
		{Name: "Máy in hóa đơn", Type: facility.TypeElectric, Area: facility.AreaCounter, Quantity: 1, Status: facility.StatusInUse, PurchaseDate: time.Date(2023, 1, 5, 0, 0, 0, 0, time.UTC), Cost: 2500000, Supplier: "Epson", Notes: "Máy in nhiệt 80mm"},
		{Name: "Cân điện tử", Type: facility.TypeElectric, Area: facility.AreaKitchen, Quantity: 2, Status: facility.StatusInUse, PurchaseDate: time.Date(2023, 2, 10, 0, 0, 0, 0, time.UTC), Cost: 800000, Supplier: "Tanita", Notes: "Cân chính xác 0.1g"},
		{Name: "Loa bluetooth", Type: facility.TypeElectric, Area: facility.AreaDiningRoom, Quantity: 4, Status: facility.StatusInUse, PurchaseDate: time.Date(2023, 3, 15, 0, 0, 0, 0, time.UTC), Cost: 1200000, Supplier: "JBL", Notes: "Loa trần âm thanh nền"},
		{Name: "Đèn LED trang trí", Type: facility.TypeElectric, Area: facility.AreaDiningRoom, Quantity: 12, Status: facility.StatusInUse, PurchaseDate: time.Date(2023, 1, 30, 0, 0, 0, 0, time.UTC), Cost: 250000, Supplier: "Philips", Notes: "Đèn LED dây 5m"},
		{Name: "Quạt trần", Type: facility.TypeElectric, Area: facility.AreaDiningRoom, Quantity: 3, Status: facility.StatusRepairing, PurchaseDate: time.Date(2022, 12, 1, 0, 0, 0, 0, time.UTC), Cost: 1500000, Supplier: "Panasonic", Notes: "1 chiếc cần thay motor"},
		
		// Khác
		{Name: "Kệ để đồ inox", Type: facility.TypeOther, Area: facility.AreaKitchen, Quantity: 3, Status: facility.StatusInUse, PurchaseDate: time.Date(2023, 1, 12, 0, 0, 0, 0, time.UTC), Cost: 800000, Supplier: "Inox Kitchen", Notes: "Kệ 4 tầng inox 304"},
		{Name: "Thùng rác inox", Type: facility.TypeOther, Area: facility.AreaDiningRoom, Quantity: 6, Status: facility.StatusInUse, PurchaseDate: time.Date(2023, 1, 12, 0, 0, 0, 0, time.UTC), Cost: 350000, Supplier: "Inox Kitchen", Notes: "Thùng rác có nắp đạp"},
		{Name: "Bảng menu treo tường", Type: facility.TypeOther, Area: facility.AreaDiningRoom, Quantity: 2, Status: facility.StatusInUse, PurchaseDate: time.Date(2023, 2, 20, 0, 0, 0, 0, time.UTC), Cost: 500000, Supplier: "Quảng cáo Minh Tâm", Notes: "Bảng mica có đèn LED"},
		{Name: "Két sắt mini", Type: facility.TypeOther, Area: facility.AreaOffice, Quantity: 1, Status: facility.StatusInUse, PurchaseDate: time.Date(2023, 1, 8, 0, 0, 0, 0, time.UTC), Cost: 2200000, Supplier: "Safeguard", Notes: "Két sắt điện tử"},
		{Name: "Bàn làm việc", Type: facility.TypeFurniture, Area: facility.AreaOffice, Quantity: 1, Status: facility.StatusInUse, PurchaseDate: time.Date(2023, 1, 8, 0, 0, 0, 0, time.UTC), Cost: 1500000, Supplier: "Nội thất Hòa Phát", Notes: "Bàn gỗ có ngăn kéo"},
		{Name: "Ghế xoay văn phòng", Type: facility.TypeFurniture, Area: facility.AreaOffice, Quantity: 1, Status: facility.StatusInUse, PurchaseDate: time.Date(2023, 1, 8, 0, 0, 0, 0, time.UTC), Cost: 800000, Supplier: "Nội thất Hòa Phát", Notes: "Ghế xoay có tựa lưng"},
	}

	for _, fac := range facilities {
		fac.CreatedAt = time.Now()
		fac.UpdatedAt = time.Now()
		facilityRepo.Create(ctx, &fac)
		log.Printf("Created facility: %s - %s (%s)", fac.Name, fac.Type, fac.Status)
	}

	// Seed some maintenance records and issue reports
	seedMaintenanceRecords(ctx, facilityRepo)
	seedIssueReports(ctx, facilityRepo)
}

func seedMaintenanceRecords(ctx context.Context, facilityRepo *mongodb.FacilityRepository) {
	// Get some facilities for maintenance records
	facilities, _ := facilityRepo.GetAll(ctx)
	if len(facilities) == 0 {
		return
	}

	maintenanceRecords := []facility.MaintenanceRecord{
		{
			FacilityID:  facilities[0].ID,
			Type:        "scheduled",
			Description: "Bảo dưỡng định kỳ máy pha cà phê - thay filter và vệ sinh",
			Cost:        500000,
			Vendor:      "La Marzocco Service",
			Date:        time.Date(2024, 1, 15, 0, 0, 0, 0, time.UTC),
			Username:    "admin",
			CreatedAt:   time.Now(),
		},
		{
			FacilityID:  facilities[1].ID,
			Type:        "emergency",
			Description: "Sửa chữa máy xay cà phê - thay burr mới",
			Cost:        1200000,
			Vendor:      "Mazzer Service",
			Date:        time.Date(2024, 2, 10, 0, 0, 0, 0, time.UTC),
			Username:    "admin",
			CreatedAt:   time.Now(),
		},
		{
			FacilityID:  facilities[2].ID,
			Type:        "scheduled",
			Description: "Vệ sinh và kiểm tra tủ lạnh",
			Cost:        200000,
			Vendor:      "Điện lạnh Minh Phát",
			Date:        time.Date(2024, 3, 1, 0, 0, 0, 0, time.UTC),
			Username:    "admin",
			CreatedAt:   time.Now(),
		},
	}

	for _, record := range maintenanceRecords {
		facilityRepo.CreateMaintenanceRecord(ctx, &record)
		log.Printf("Created maintenance record for facility: %s", record.Description)
	}
}

func seedIssueReports(ctx context.Context, facilityRepo *mongodb.FacilityRepository) {
	// Get some facilities for issue reports
	facilities, _ := facilityRepo.GetAll(ctx)
	if len(facilities) == 0 {
		return
	}

	issueReports := []facility.IssueReport{
		{
			FacilityID:  facilities[0].ID,
			Description: "Máy pha cà phê có tiếng kêu lạ khi hoạt động",
			Severity:    "medium",
			Status:      "open",
			Username:    "waiter1",
			CreatedAt:   time.Now().AddDate(0, 0, -2),
		},
		{
			FacilityID:  facilities[1].ID,
			Description: "Ghế cafe số 5 bị lỏng ốc vít",
			Severity:    "low",
			Status:      "resolved",
			Username:    "waiter2",
			CreatedAt:   time.Now().AddDate(0, 0, -5),
			ResolvedAt:  &[]time.Time{time.Now().AddDate(0, 0, -3)}[0],
		},
		{
			FacilityID:  facilities[2].ID,
			Description: "Tủ lạnh không đủ lạnh, cần kiểm tra",
			Severity:    "high",
			Status:      "in_progress",
			Username:    "cashier1",
			CreatedAt:   time.Now().AddDate(0, 0, -1),
		},
	}

	for _, report := range issueReports {
		facilityRepo.CreateIssueReport(ctx, &report)
		log.Printf("Created issue report: %s (%s)", report.Description, report.Severity)
	}
	
	// Seed status history
	seedStatusHistory(ctx, facilityRepo)
}

func seedStatusHistory(ctx context.Context, facilityRepo *mongodb.FacilityRepository) {
	facilities, _ := facilityRepo.GetAll(ctx)
	if len(facilities) == 0 {
		return
	}

	statusHistories := []facility.FacilityHistory{
		{
			FacilityID:  facilities[0].ID,
			Action:      "status_change",
			Description: "Thay đổi trạng thái từ 'Đang sử dụng' sang 'Đang sửa'",
			OldValue:    "Đang sử dụng",
			NewValue:    "Đang sửa",
			UserID:      primitive.NewObjectID(),
			Username:    "admin",
			CreatedAt:   time.Now().AddDate(0, 0, -10),
		},
		{
			FacilityID:  facilities[0].ID,
			Action:      "status_change",
			Description: "Thay đổi trạng thái từ 'Đang sửa' sang 'Đang sử dụng'",
			OldValue:    "Đang sửa",
			NewValue:    "Đang sử dụng",
			UserID:      primitive.NewObjectID(),
			Username:    "admin",
			CreatedAt:   time.Now().AddDate(0, 0, -5),
		},
		{
			FacilityID:  facilities[1].ID,
			Action:      "status_change",
			Description: "Thay đổi trạng thái từ 'Đang sử dụng' sang 'Hỏng'",
			OldValue:    "Đang sử dụng",
			NewValue:    "Hỏng",
			UserID:      primitive.NewObjectID(),
			Username:    "waiter1",
			CreatedAt:   time.Now().AddDate(0, 0, -7),
		},
		{
			FacilityID:  facilities[2].ID,
			Action:      "status_change",
			Description: "Thay đổi trạng thái từ 'Hỏng' sang 'Đang sửa'",
			OldValue:    "Hỏng",
			NewValue:    "Đang sửa",
			UserID:      primitive.NewObjectID(),
			Username:    "admin",
			CreatedAt:   time.Now().AddDate(0, 0, -3),
		},
		{
			FacilityID:  facilities[3].ID,
			Action:      "status_change",
			Description: "Thay đổi trạng thái từ 'Đang sử dụng' sang 'Ngừng sử dụng'",
			OldValue:    "Đang sử dụng",
			NewValue:    "Ngừng sử dụng",
			UserID:      primitive.NewObjectID(),
			Username:    "admin",
			CreatedAt:   time.Now().AddDate(0, 0, -15),
		},
	}

	for _, history := range statusHistories {
		facilityRepo.CreateHistory(ctx, &history)
		log.Printf("Created status history: %s", history.Description)
	}
}