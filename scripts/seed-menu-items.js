// MongoDB Script to Seed Menu Items
// Run: mongosh mongodb://admin:password123@localhost:27017/cafe_pos --authenticationDatabase admin < scripts/seed-menu-items.js

// Switch to database
db = db.getSiblingDB('cafe_pos');

// Clear existing menu items (optional - comment out if you want to keep existing items)
// db.menu_items.deleteMany({});

// Seed menu items
const menuItems = [
  // ========== CÀ PHÊ ==========
  {
    name: "Cà phê đen",
    price: 25000,
    category: "Cà phê",
    description: "Cà phê phin truyền thống, đậm đà",
    ingredients: [
      { name: "Cà phê", quantity: 20, unit: "g" },
      { name: "Nước nóng", quantity: 100, unit: "ml" }
    ],
    available: true,
    created_at: new Date(),
    updated_at: new Date()
  },
  {
    name: "Cà phê sữa",
    price: 30000,
    category: "Cà phê",
    description: "Cà phê phin với sữa đặc ngọt ngào",
    ingredients: [
      { name: "Cà phê", quantity: 20, unit: "g" },
      { name: "Sữa đặc", quantity: 30, unit: "ml" },
      { name: "Nước nóng", quantity: 100, unit: "ml" }
    ],
    available: true,
    created_at: new Date(),
    updated_at: new Date()
  },
  {
    name: "Bạc xỉu",
    price: 32000,
    category: "Cà phê",
    description: "Cà phê sữa nhiều sữa, ít cà phê",
    ingredients: [
      { name: "Cà phê", quantity: 15, unit: "g" },
      { name: "Sữa đặc", quantity: 50, unit: "ml" },
      { name: "Nước nóng", quantity: 100, unit: "ml" }
    ],
    available: true,
    created_at: new Date(),
    updated_at: new Date()
  },
  {
    name: "Cà phê đá",
    price: 28000,
    category: "Cà phê",
    description: "Cà phê đen mát lạnh",
    ingredients: [
      { name: "Cà phê", quantity: 20, unit: "g" },
      { name: "Nước nóng", quantity: 80, unit: "ml" },
      { name: "Đá", quantity: 100, unit: "g" }
    ],
    available: true,
    created_at: new Date(),
    updated_at: new Date()
  },
  {
    name: "Cà phê sữa đá",
    price: 32000,
    category: "Cà phê",
    description: "Cà phê sữa mát lạnh",
    ingredients: [
      { name: "Cà phê", quantity: 20, unit: "g" },
      { name: "Sữa đặc", quantity: 30, unit: "ml" },
      { name: "Nước nóng", quantity: 80, unit: "ml" },
      { name: "Đá", quantity: 100, unit: "g" }
    ],
    available: true,
    created_at: new Date(),
    updated_at: new Date()
  },

  // ========== TRÀ SỮA ==========
  {
    name: "Trà sữa truyền thống",
    price: 35000,
    category: "Trà sữa",
    description: "Trà sữa đài loan cổ điển",
    ingredients: [
      { name: "Trà đen", quantity: 10, unit: "g" },
      { name: "Sữa tươi", quantity: 100, unit: "ml" },
      { name: "Đường", quantity: 20, unit: "g" },
      { name: "Trân châu", quantity: 50, unit: "g" }
    ],
    available: true,
    created_at: new Date(),
    updated_at: new Date()
  },
  {
    name: "Trà sữa matcha",
    price: 40000,
    category: "Trà sữa",
    description: "Trà xanh matcha Nhật Bản với sữa",
    ingredients: [
      { name: "Matcha", quantity: 5, unit: "g" },
      { name: "Sữa tươi", quantity: 100, unit: "ml" },
      { name: "Đường", quantity: 20, unit: "g" },
      { name: "Trân châu", quantity: 50, unit: "g" }
    ],
    available: true,
    created_at: new Date(),
    updated_at: new Date()
  },
  {
    name: "Trà sữa socola",
    price: 38000,
    category: "Trà sữa",
    description: "Trà sữa vị socola đậm đà",
    ingredients: [
      { name: "Trà đen", quantity: 10, unit: "g" },
      { name: "Sữa tươi", quantity: 100, unit: "ml" },
      { name: "Bột socola", quantity: 15, unit: "g" },
      { name: "Đường", quantity: 20, unit: "g" },
      { name: "Trân châu", quantity: 50, unit: "g" }
    ],
    available: true,
    created_at: new Date(),
    updated_at: new Date()
  },

  // ========== TRÀ TRÁI CÂY ==========
  {
    name: "Trà đào cam sả",
    price: 35000,
    category: "Trà trái cây",
    description: "Trà đào tươi mát với cam và sả thơm",
    ingredients: [
      { name: "Trà xanh", quantity: 10, unit: "g" },
      { name: "Đào", quantity: 50, unit: "g" },
      { name: "Cam", quantity: 30, unit: "g" },
      { name: "Sả", quantity: 5, unit: "g" },
      { name: "Đường", quantity: 20, unit: "g" }
    ],
    available: true,
    created_at: new Date(),
    updated_at: new Date()
  },
  {
    name: "Trà chanh leo",
    price: 32000,
    category: "Trà trái cây",
    description: "Trà xanh với chanh leo chua ngọt",
    ingredients: [
      { name: "Trà xanh", quantity: 10, unit: "g" },
      { name: "Chanh leo", quantity: 50, unit: "g" },
      { name: "Đường", quantity: 25, unit: "g" }
    ],
    available: true,
    created_at: new Date(),
    updated_at: new Date()
  },
  {
    name: "Trà vải",
    price: 33000,
    category: "Trà trái cây",
    description: "Trà xanh với vải tươi ngọt mát",
    ingredients: [
      { name: "Trà xanh", quantity: 10, unit: "g" },
      { name: "Vải", quantity: 60, unit: "g" },
      { name: "Đường", quantity: 20, unit: "g" }
    ],
    available: true,
    created_at: new Date(),
    updated_at: new Date()
  },

  // ========== SINH TỐ ==========
  {
    name: "Sinh tố bơ",
    price: 40000,
    category: "Sinh tố",
    description: "Sinh tố bơ béo ngậy",
    ingredients: [
      { name: "Bơ", quantity: 150, unit: "g" },
      { name: "Sữa tươi", quantity: 100, unit: "ml" },
      { name: "Đường", quantity: 20, unit: "g" },
      { name: "Đá", quantity: 50, unit: "g" }
    ],
    available: true,
    created_at: new Date(),
    updated_at: new Date()
  },
  {
    name: "Sinh tố dâu",
    price: 38000,
    category: "Sinh tố",
    description: "Sinh tố dâu tây tươi mát",
    ingredients: [
      { name: "Dâu tây", quantity: 100, unit: "g" },
      { name: "Sữa tươi", quantity: 100, unit: "ml" },
      { name: "Đường", quantity: 20, unit: "g" },
      { name: "Đá", quantity: 50, unit: "g" }
    ],
    available: true,
    created_at: new Date(),
    updated_at: new Date()
  },
  {
    name: "Sinh tố xoài",
    price: 38000,
    category: "Sinh tố",
    description: "Sinh tố xoài ngọt thơm",
    ingredients: [
      { name: "Xoài", quantity: 150, unit: "g" },
      { name: "Sữa tươi", quantity: 100, unit: "ml" },
      { name: "Đường", quantity: 15, unit: "g" },
      { name: "Đá", quantity: 50, unit: "g" }
    ],
    available: true,
    created_at: new Date(),
    updated_at: new Date()
  },

  // ========== NƯỚC ÉP ==========
  {
    name: "Nước ép cam",
    price: 35000,
    category: "Nước ép",
    description: "Nước cam tươi 100%",
    ingredients: [
      { name: "Cam", quantity: 200, unit: "g" },
      { name: "Đường", quantity: 10, unit: "g" }
    ],
    available: true,
    created_at: new Date(),
    updated_at: new Date()
  },
  {
    name: "Nước ép dưa hấu",
    price: 30000,
    category: "Nước ép",
    description: "Nước dưa hấu mát lạnh",
    ingredients: [
      { name: "Dưa hấu", quantity: 250, unit: "g" },
      { name: "Đường", quantity: 5, unit: "g" }
    ],
    available: true,
    created_at: new Date(),
    updated_at: new Date()
  },

  // ========== BÁNH NGỌT ==========
  {
    name: "Bánh tiramisu",
    price: 45000,
    category: "Bánh ngọt",
    description: "Bánh tiramisu Ý truyền thống",
    ingredients: [
      { name: "Bánh tiramisu", quantity: 1, unit: "miếng" }
    ],
    available: true,
    created_at: new Date(),
    updated_at: new Date()
  },
  {
    name: "Bánh cheesecake",
    price: 42000,
    category: "Bánh ngọt",
    description: "Bánh phô mai mềm mịn",
    ingredients: [
      { name: "Bánh cheesecake", quantity: 1, unit: "miếng" }
    ],
    available: true,
    created_at: new Date(),
    updated_at: new Date()
  },
  {
    name: "Bánh croissant",
    price: 35000,
    category: "Bánh ngọt",
    description: "Bánh sừng bò Pháp giòn tan",
    ingredients: [
      { name: "Bánh croissant", quantity: 1, unit: "cái" }
    ],
    available: true,
    created_at: new Date(),
    updated_at: new Date()
  },
  {
    name: "Bánh muffin",
    price: 30000,
    category: "Bánh ngọt",
    description: "Bánh muffin chocolate chip",
    ingredients: [
      { name: "Bánh muffin", quantity: 1, unit: "cái" }
    ],
    available: true,
    created_at: new Date(),
    updated_at: new Date()
  }
];

// Insert menu items
const result = db.menu_items.insertMany(menuItems);

print("✅ Seeded " + result.insertedIds.length + " menu items successfully!");
print("\n📋 Menu Categories:");
print("- Cà phê: 5 items");
print("- Trà sữa: 3 items");
print("- Trà trái cây: 3 items");
print("- Sinh tố: 3 items");
print("- Nước ép: 2 items");
print("- Bánh ngọt: 4 items");
print("\n💰 Price range: 25,000đ - 45,000đ");
print("\n🎉 Ready to create orders!");
