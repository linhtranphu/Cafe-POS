# 🍽️ Menu Seed Guide

## 📋 Overview

Script để seed menu items vào database, giúp test và demo hệ thống.

---

## 🚀 Quick Start

### Cách 1: Sử dụng Go Binary (Recommended)

```bash
# Build
cd backend/cmd/seed-menu
go build -o seed-menu

# Run
./seed-menu
```

### Cách 2: Chạy trực tiếp

```bash
cd backend/cmd/seed-menu
go run main.go
```

---

## 📦 Menu Items Seeded

### Tổng quan
- **Total items:** 20 món
- **Categories:** 6 loại
- **Price range:** 25,000đ - 45,000đ

---

### 1. Cà Phê (5 items)

| Tên món | Giá | Mô tả |
|---------|-----|-------|
| Cà phê đen | 25,000đ | Cà phê phin truyền thống, đậm đà |
| Cà phê sữa | 30,000đ | Cà phê phin với sữa đặc ngọt ngào |
| Bạc xỉu | 32,000đ | Cà phê sữa nhiều sữa, ít cà phê |
| Cà phê đá | 28,000đ | Cà phê đen mát lạnh |
| Cà phê sữa đá | 32,000đ | Cà phê sữa mát lạnh |

---

### 2. Trà Sữa (3 items)

| Tên món | Giá | Mô tả |
|---------|-----|-------|
| Trà sữa truyền thống | 35,000đ | Trà sữa đài loan cổ điển |
| Trà sữa matcha | 40,000đ | Trà xanh matcha Nhật Bản với sữa |
| Trà sữa socola | 38,000đ | Trà sữa vị socola đậm đà |

---

### 3. Trà Trái Cây (3 items)

| Tên món | Giá | Mô tả |
|---------|-----|-------|
| Trà đào cam sả | 35,000đ | Trà đào tươi mát với cam và sả thơm |
| Trà chanh leo | 32,000đ | Trà xanh với chanh leo chua ngọt |
| Trà vải | 33,000đ | Trà xanh với vải tươi ngọt mát |

---

### 4. Sinh Tố (3 items)

| Tên món | Giá | Mô tả |
|---------|-----|-------|
| Sinh tố bơ | 40,000đ | Sinh tố bơ béo ngậy |
| Sinh tố dâu | 38,000đ | Sinh tố dâu tây tươi mát |
| Sinh tố xoài | 38,000đ | Sinh tố xoài ngọt thơm |

---

### 5. Nước Ép (2 items)

| Tên món | Giá | Mô tả |
|---------|-----|-------|
| Nước ép cam | 35,000đ | Nước cam tươi 100% |
| Nước ép dưa hấu | 30,000đ | Nước dưa hấu mát lạnh |

---

### 6. Bánh Ngọt (4 items)

| Tên món | Giá | Mô tả |
|---------|-----|-------|
| Bánh tiramisu | 45,000đ | Bánh tiramisu Ý truyền thống |
| Bánh cheesecake | 42,000đ | Bánh phô mai mềm mịn |
| Bánh croissant | 35,000đ | Bánh sừng bò Pháp giòn tan |
| Bánh muffin | 30,000đ | Bánh muffin chocolate chip |

---

## 🧪 Testing Workflow

### 1. Seed Menu Items
```bash
cd backend/cmd/seed-menu
./seed-menu
```

### 2. Login as Waiter
```
Username: waiter1
Password: password123
```

### 3. Start Shift
```
1. Go to /shift
2. Select shift type
3. Enter start cash: 0
4. Click "Mở ca"
```

### 4. Create Order
```
1. Go to /orders
2. Click "Tạo đơn mới"
3. Select table
4. Add items from menu
5. Complete order
```

### 5. Process Payment
```
1. Select payment method: Cash
2. Enter amount
3. Complete payment
4. Cash will be added to shift
```

### 6. Test Handover
```
1. Go to /shift
2. See "Tiền hiện có" > 0
3. Click "💰 Bàn giao một phần"
4. Enter amount
5. Submit
```

---

## 🔄 Re-seeding

Nếu muốn seed lại:

```bash
# Script sẽ hỏi confirm nếu đã có data
./seed-menu

# Output:
⚠️  Found 20 existing menu items
Do you want to clear and reseed? (y/N): y
🗑️  Cleared existing menu items
✅ Seeded 20 menu items successfully!
```

---

## 📊 Database Structure

### Collection: `menu_items`

```javascript
{
  _id: ObjectId("..."),
  name: "Cà phê đen",
  price: 25000,
  category: "Cà phê",
  description: "Cà phê phin truyền thống, đậm đà",
  ingredients: [
    { name: "Cà phê", quantity: 20, unit: "g" },
    { name: "Nước nóng", quantity: 100, unit: "ml" }
  ],
  available: true,
  created_at: ISODate("2026-02-04T..."),
  updated_at: ISODate("2026-02-04T...")
}
```

---

## 🎯 Use Cases

### 1. Demo System
- Show full menu to clients
- Test order creation
- Test payment processing

### 2. Development
- Test menu filtering by category
- Test search functionality
- Test price calculations

### 3. Testing Handover
- Create orders with cash payment
- Accumulate cash in shift
- Test handover workflow

---

## 🛠️ Customization

### Add More Items

Edit `backend/cmd/seed-menu/main.go`:

```go
menuItems := []interface{}{
    // ... existing items ...
    
    // Add new item
    MenuItem{
        Name:        "Cà phê espresso",
        Price:       28000,
        Category:    "Cà phê",
        Description: "Cà phê espresso đậm đà",
        Ingredients: []Ingredient{
            {Name: "Cà phê", Quantity: 18, Unit: "g"},
        },
        Available: true,
        CreatedAt: time.Now(),
        UpdatedAt: time.Now(),
    },
}
```

### Change Prices

```go
MenuItem{
    Name:  "Cà phê đen",
    Price: 30000,  // Changed from 25000
    // ...
}
```

### Add New Category

```go
// ========== ĐỒ ĂN ==========
MenuItem{
    Name:        "Sandwich",
    Price:       45000,
    Category:    "Đồ ăn",  // New category
    Description: "Sandwich thịt nguội",
    // ...
}
```

---

## 🐛 Troubleshooting

### Error: Cannot connect to MongoDB

**Solution:**
```bash
# Check MongoDB is running
docker ps | grep mongo

# Check connection string
export MONGODB_URI="mongodb://admin:password123@localhost:27017"
```

### Error: Permission denied

**Solution:**
```bash
chmod +x seed-menu
```

### Items not showing in UI

**Solution:**
1. Check database:
```bash
# Using MongoDB Compass or CLI
db.menu_items.find().count()
```

2. Refresh frontend
3. Check API endpoint: `GET /api/menu`

---

## 📁 Files

### Seed Script
- `backend/cmd/seed-menu/main.go` - Go seed script
- `scripts/seed-menu-items.js` - MongoDB shell script (alternative)
- `scripts/seed-menu.sh` - Bash wrapper (requires mongosh)

### Related Files
- `backend/domain/menu/menu.go` - Menu domain model
- `backend/infrastructure/mongodb/menu_repository.go` - Menu repository
- `backend/application/services/menu.go` - Menu service
- `backend/interfaces/http/menu_handler.go` - Menu HTTP handlers

---

## 🔗 Related Documentation

- [ORDER_IMPLEMENTATION.md](./ORDER_IMPLEMENTATION.md) - Order system
- [CASH_HANDOVER_WAITER_GUIDE.md](./CASH_HANDOVER_WAITER_GUIDE.md) - Handover guide
- [PROJECT_STRUCTURE.md](./PROJECT_STRUCTURE.md) - Project structure

---

## ✅ Verification

After seeding, verify:

```bash
# 1. Check count
curl http://localhost:8080/api/menu | jq 'length'
# Expected: 20

# 2. Check categories
curl http://localhost:8080/api/menu | jq 'group_by(.category) | map({category: .[0].category, count: length})'

# 3. Check price range
curl http://localhost:8080/api/menu | jq '[.[].price] | min, max'
# Expected: 25000, 45000
```

---

**Last Updated:** 2026-02-04  
**Version:** 1.0  
**Status:** ✅ Complete
