# Hướng Dẫn Test Chi Phí và Lợi Nhuận Món Ăn

## Tổng Quan

Hướng dẫn này giúp bạn test tính năng phân tích chi phí và lợi nhuận cho món ăn có nhiều size (variants).

## Bước 1: Seed Dữ Liệu Menu

### 1.1. Chạy Seed Command

```bash
go run backend/cmd/seed-menu-variants/main.go
```

### 1.2. Dữ Liệu Được Tạo

**Món có nhiều size (Multi-size):**
1. ☕ **Cà phê sữa đá** (M: 25k, L: 30k, XL: 35k)
2. ☕ **Cà phê đen đá** (M: 20k, L: 25k)
3. 🧋 **Trà sữa truyền thống** (M: 35k, L: 42k, XL: 48k)
4. 🍵 **Trà sữa matcha** (M: 40k, L: 48k)
5. 🍑 **Trà đào cam sả** (M: 35k, L: 42k)
6. 🥑 **Sinh tố bơ** (M: 40k, L: 50k)

**Món size đơn (Single-size):**
7. 🥖 **Bánh mì thịt** (20k)
8. 🍰 **Bánh tiramisu** (45k)
9. 🥐 **Bánh croissant** (35k)
10. 🍊 **Nước ép cam** (35k)

## Bước 2: Tính Toán Chi Phí

### 2.1. Tính Chi Phí Tất Cả Món

```bash
./calculate-all-menu-costs.sh
```

Script này sẽ:
- Lấy danh sách tất cả món ăn
- Gọi API tính chi phí cho từng món
- Hiển thị kết quả (FINAL/INCOMPLETE status)

### 2.2. Tính Chi Phí Một Món Cụ Thể

```bash
curl -X POST http://localhost:8080/api/menu/{MENU_ITEM_ID}/calculate-cost
```

**Response:**
```json
{
  "menu_item_id": "...",
  "current_cost": 13800,
  "cost_status": "FINAL",
  "cost_last_calculated_at": "2026-02-13T10:00:00Z"
}
```

## Bước 3: Xem Phân Tích Chi Phí

### 3.1. Xem Tổng Quan Tất Cả Món

```bash
./view-menu-cost-analysis.sh
```

Script này hiển thị:
- Giá bán (Price)
- Chi phí (Cost)
- Lợi nhuận (Profit)
- Tỷ lệ lợi nhuận (Margin %)
- So sánh giữa các size

**Ví dụ Output:**

```
━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━
🔄 Cà phê sữa đá (Cà phê) - MULTI-SIZE
━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━

Variant    | Price        | Cost         | Profit       | Margin
-----------|--------------|--------------|--------------|------------
Size M     |     25000 đ |     13800 đ |     11200 đ |      44.8%
Size L     |     30000 đ |     20700 đ |      9300 đ |      31.0%
Size XL    |     35000 đ |     27600 đ |      7400 đ |      21.1%

⭐ Most profitable: Size M (44.8% margin)
```

### 3.2. Xem Chi Tiết Chi Phí Một Món

```bash
curl http://localhost:8080/api/menu/{MENU_ITEM_ID}/cost-breakdown
```

**Response cho món có nhiều size:**
```json
{
  "menu_item_id": "...",
  "menu_item_name": "Cà phê sữa đá",
  "has_variants": true,
  "variants": [
    {
      "variant_id": "M",
      "variant_name": "Size M",
      "price": 25000,
      "current_cost": 13800,
      "cost_status": "FINAL",
      "ingredients": [
        {
          "name": "Cà phê",
          "quantity": 20,
          "unit": "g",
          "cost_per_unit": 500,
          "conversion_rate": 1.0,
          "wastage_percentage": 5.0,
          "total_cost": 10500
        },
        {
          "name": "Sữa đặc",
          "quantity": 30,
          "unit": "ml",
          "cost_per_unit": 100,
          "conversion_rate": 1.0,
          "wastage_percentage": 10.0,
          "total_cost": 3300
        }
      ]
    }
  ]
}
```

### 3.3. Xem Phân Tích Lợi Nhuận

```bash
curl http://localhost:8080/api/menu/{MENU_ITEM_ID}/profit-analysis
```

**Response:**
```json
{
  "menu_item_id": "...",
  "menu_item_name": "Cà phê sữa đá",
  "has_variants": true,
  "variants": [
    {
      "variant_id": "M",
      "variant_name": "Size M",
      "price": 25000,
      "cost": 13800,
      "profit": 11200,
      "profit_margin": 44.8
    },
    {
      "variant_id": "L",
      "variant_name": "Size L",
      "price": 30000,
      "cost": 20700,
      "profit": 9300,
      "profit_margin": 31.0
    },
    {
      "variant_id": "XL",
      "variant_name": "Size XL",
      "price": 35000,
      "cost": 27600,
      "profit": 7400,
      "profit_margin": 21.1
    }
  ]
}
```

## Bước 4: Phân Tích Kết Quả

### 4.1. Hiểu Các Chỉ Số

**Cost Status:**
- `FINAL`: Tất cả nguyên liệu có dữ liệu chi phí đầy đủ
- `INCOMPLETE`: Thiếu dữ liệu chi phí một số nguyên liệu
- `NOT_CALCULATED`: Chưa tính toán chi phí

**Profit Margin (Tỷ lệ lợi nhuận):**
```
Profit Margin = (Price - Cost) / Price × 100%
```

### 4.2. So Sánh Giữa Các Size

**Ví dụ: Cà phê sữa đá**

| Size | Price | Cost  | Profit | Margin |
|------|-------|-------|--------|--------|
| M    | 25k   | 13.8k | 11.2k  | 44.8%  |
| L    | 30k   | 20.7k | 9.3k   | 31.0%  |
| XL   | 35k   | 27.6k | 7.4k   | 21.1%  |

**Phân tích:**
- Size M có tỷ lệ lợi nhuận cao nhất (44.8%)
- Chi phí tăng nhanh hơn giá bán khi tăng size
- Size XL có lợi nhuận thấp nhất (21.1%)

**Khuyến nghị:**
- Nên khuyến khích khách mua Size M (lợi nhuận cao nhất)
- Cân nhắc tăng giá Size L và XL để cải thiện margin
- Hoặc giảm lượng nguyên liệu để giảm chi phí

### 4.3. Công Thức Tính Chi Phí

```
Cost = Σ (quantity × cost_per_unit × conversion_rate × (1 + wastage/100))
```

**Ví dụ: Cà phê Size M**
```
Cà phê:   20g × 500đ/g × 1.0 × 1.05 = 10,500đ
Sữa đặc:  30ml × 100đ/ml × 1.0 × 1.10 = 3,300đ
Total:                                  13,800đ
```

## Bước 5: Test Trên Frontend

### 5.1. Truy Cập Cost Analysis View

```
http://localhost:5173/cost-analysis
```

### 5.2. Tính Năng Có Sẵn

1. **Danh sách món với chi phí**
   - Hiển thị tất cả món
   - Hiển thị chi phí và lợi nhuận mỗi variant
   - Filter theo cost_status

2. **Cost Breakdown Modal**
   - Click vào món để xem chi tiết
   - Hiển thị công thức tính chi phí
   - Hiển thị từng nguyên liệu

3. **Profit Comparison Modal**
   - So sánh lợi nhuận giữa các size
   - Hiển thị size có lợi nhuận cao nhất
   - Chart so sánh trực quan

## Bước 6: Test Scenarios

### Scenario 1: Tính Chi Phí Món Mới

1. Tạo món mới qua API hoặc UI
2. Gọi API calculate-cost
3. Verify cost_status = FINAL
4. Verify chi phí được tính đúng

### Scenario 2: Cập Nhật Giá Nguyên Liệu

1. Cập nhật giá nguyên liệu (ví dụ: cà phê từ 500đ → 600đ/g)
2. Gọi lại API calculate-cost cho các món có cà phê
3. Verify chi phí tăng tương ứng
4. Verify profit margin giảm

### Scenario 3: Thiếu Dữ Liệu Nguyên Liệu

1. Tạo món với nguyên liệu chưa có trong DB
2. Gọi API calculate-cost
3. Verify cost_status = INCOMPLETE
4. Verify hiển thị warning về nguyên liệu thiếu

### Scenario 4: So Sánh Lợi Nhuận

1. Xem profit analysis cho món có nhiều size
2. Verify size nào có margin cao nhất
3. Verify tính toán profit đúng
4. Verify hiển thị comparison chart

## Bước 7: API Endpoints Tổng Hợp

### Menu Cost APIs

```bash
# Get all menu items with costs
GET /api/menu

# Get menu item detail with costs
GET /api/menu/:id

# Calculate cost for menu item
POST /api/menu/:id/calculate-cost

# Get cost breakdown (detailed)
GET /api/menu/:id/cost-breakdown

# Get profit analysis
GET /api/menu/:id/profit-analysis

# Get all menu costs (summary)
GET /api/menu/costs

# Get menu warnings (incomplete costs)
GET /api/menu/warnings
```

## Troubleshooting

### Lỗi: "Ingredient not found"

**Nguyên nhân:** Nguyên liệu chưa có trong database

**Giải pháp:**
```bash
# Seed ingredients first
go run backend/cmd/seed/main.go
```

### Lỗi: "Cost status INCOMPLETE"

**Nguyên nhân:** Một số nguyên liệu thiếu cost_per_unit

**Giải pháp:**
1. Kiểm tra ingredients collection
2. Cập nhật cost_per_unit cho nguyên liệu thiếu
3. Chạy lại calculate-cost

### Chi phí = 0

**Nguyên nhân:** Chưa chạy calculate-cost

**Giải pháp:**
```bash
./calculate-all-menu-costs.sh
```

## Kết Luận

Bạn đã có đầy đủ dữ liệu và công cụ để test tính năng phân tích chi phí và lợi nhuận:

✅ Seed data với 10 món (6 multi-size, 4 single-size)
✅ Scripts tự động tính chi phí
✅ Scripts xem phân tích chi phí
✅ API endpoints đầy đủ
✅ Frontend views để test UI

Chúc bạn test thành công! 🎉
