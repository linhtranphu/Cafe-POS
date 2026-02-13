# Seed Data & Cost Analysis Tools - Summary

## Tổng Quan

Đã tạo đầy đủ seed data và tools để test tính năng phân tích chi phí và lợi nhuận cho món ăn có nhiều size (variants).

## Files Đã Tạo

### 1. Seed Command
📁 `backend/cmd/seed-menu-variants/main.go`
- Tạo 10 món ăn (6 multi-size, 4 single-size)
- Mỗi món có đầy đủ thông tin: name, category, price, ingredients
- Multi-size items có 2-3 variants với giá và nguyên liệu khác nhau

### 2. Scripts

#### 📜 `calculate-all-menu-costs.sh`
- Tự động tính chi phí cho tất cả món
- Hiển thị progress và summary
- Báo cáo success/failed items

#### 📜 `view-menu-cost-analysis.sh`
- Hiển thị phân tích chi phí đẹp mắt
- So sánh giữa các size
- Highlight size có lợi nhuận cao nhất
- Sử dụng colors để dễ đọc

### 3. Documentation

#### 📖 `HUONG_DAN_TEST_CHI_PHI_LOI_NHUAN.md`
- Hướng dẫn chi tiết từng bước
- Giải thích các chỉ số (cost, profit, margin)
- Test scenarios
- Troubleshooting guide

#### 📖 `QUICK_START_COST_ANALYSIS.md`
- Quick start guide 3 bước
- Kết quả mẫu
- Checklist

## Dữ Liệu Seed

### Multi-Size Items (6 món)

1. **Cà phê sữa đá**
   - Size M: 25,000đ (Cà phê 20g, Sữa đặc 30ml)
   - Size L: 30,000đ (Cà phê 30g, Sữa đặc 45ml)
   - Size XL: 35,000đ (Cà phê 40g, Sữa đặc 60ml)

2. **Cà phê đen đá**
   - Size M: 20,000đ (Cà phê 20g)
   - Size L: 25,000đ (Cà phê 30g)

3. **Trà sữa truyền thống**
   - Size M: 35,000đ
   - Size L: 42,000đ
   - Size XL: 48,000đ

4. **Trà sữa matcha**
   - Size M: 40,000đ
   - Size L: 48,000đ

5. **Trà đào cam sả**
   - Size M: 35,000đ
   - Size L: 42,000đ

6. **Sinh tố bơ**
   - Size M: 40,000đ
   - Size L: 50,000đ

### Single-Size Items (4 món)

7. Bánh mì thịt - 20,000đ
8. Bánh tiramisu - 45,000đ
9. Bánh croissant - 35,000đ
10. Nước ép cam - 35,000đ

## Workflow Test

```
┌─────────────────────┐
│  1. Seed Menu Data  │
│  (10 items)         │
└──────────┬──────────┘
           │
           ▼
┌─────────────────────┐
│  2. Calculate Costs │
│  (all items)        │
└──────────┬──────────┘
           │
           ▼
┌─────────────────────┐
│  3. View Analysis   │
│  (terminal/API/UI)  │
└─────────────────────┘
```

## Kết Quả Mong Đợi

### Cost Calculation
- ✅ Tất cả món có cost_status = FINAL hoặc INCOMPLETE
- ✅ Chi phí được tính theo công thức: quantity × cost_per_unit × conversion_rate × (1 + wastage/100)
- ✅ Mỗi variant có chi phí riêng

### Profit Analysis
- ✅ Profit = Price - Cost
- ✅ Margin = (Profit / Price) × 100%
- ✅ Size M thường có margin cao nhất (ít nguyên liệu hơn)
- ✅ Size XL có margin thấp nhất (nhiều nguyên liệu)

### Example: Cà phê sữa đá

| Size | Price | Cost   | Profit | Margin |
|------|-------|--------|--------|--------|
| M    | 25k   | 13.8k  | 11.2k  | 44.8%  |
| L    | 30k   | 20.7k  | 9.3k   | 31.0%  |
| XL   | 35k   | 27.6k  | 7.4k   | 21.1%  |

## Commands Cheat Sheet

```bash
# 1. Seed menu data
go run backend/cmd/seed-menu-variants/main.go

# 2. Calculate all costs
./calculate-all-menu-costs.sh

# 3. View analysis
./view-menu-cost-analysis.sh

# 4. Calculate single item
curl -X POST http://localhost:8080/api/menu/{ID}/calculate-cost

# 5. View cost breakdown
curl http://localhost:8080/api/menu/{ID}/cost-breakdown

# 6. View profit analysis
curl http://localhost:8080/api/menu/{ID}/profit-analysis
```

## API Endpoints

### Cost Calculation
- `POST /api/menu/:id/calculate-cost` - Tính chi phí
- `GET /api/menu/:id/cost-breakdown` - Chi tiết chi phí
- `GET /api/menu/:id/profit-analysis` - Phân tích lợi nhuận

### Menu Management
- `GET /api/menu` - Danh sách món (có cost data)
- `GET /api/menu/:id` - Chi tiết món
- `GET /api/menu/costs` - Tổng hợp chi phí
- `GET /api/menu/warnings` - Món có cost INCOMPLETE

## Frontend Views

### Cost Analysis View
- URL: `http://localhost:5173/cost-analysis`
- Features:
  - Danh sách món với chi phí
  - Filter theo cost_status
  - Cost breakdown modal
  - Profit comparison modal

## Test Coverage

### Integration Test
✅ `TestCostAnalysisFlow_Integration` - 470 lines
- Create multi-size item
- Calculate costs
- View cost breakdown
- View profit analysis
- Verify cost_status updates
- Update ingredient price & recalculate

### Requirements Verified
- ✅ AC-10.1 to AC-10.5 (Cost display)
- ✅ AC-11.1 to AC-11.7 (Cost calculation)
- ✅ AC-12.1 to AC-12.4 (Profit analysis)
- ✅ FR-6.4, FR-6.6 (Cost recalculation)

## Business Insights

### Pricing Strategy
1. **Size M** thường có margin cao nhất
   - Khuyến khích khách mua Size M
   - Upsell từ M → L cần cẩn thận

2. **Size XL** có margin thấp
   - Cân nhắc tăng giá
   - Hoặc giảm lượng nguyên liệu

3. **Cost vs Price**
   - Chi phí tăng nhanh hơn giá bán
   - Cần review pricing strategy

### Cost Management
1. Monitor ingredient costs
2. Track wastage percentage
3. Optimize recipes for better margins
4. Identify high-cost items

## Next Steps

1. ✅ Seed data created
2. ✅ Scripts created
3. ✅ Documentation complete
4. ⏭️ Run seed command
5. ⏭️ Calculate costs
6. ⏭️ Test on frontend
7. ⏭️ Analyze results

## Support

Nếu gặp vấn đề:
1. Xem `HUONG_DAN_TEST_CHI_PHI_LOI_NHUAN.md` (Troubleshooting section)
2. Check backend logs
3. Verify ingredients seeded
4. Check API responses

---

**Tạo bởi:** Task 13.5 Implementation
**Ngày:** 2026-02-13
**Status:** ✅ Ready to use
