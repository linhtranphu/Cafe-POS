# Phân tích: Batch Cost Calculation Flow

## 🎯 CÂU HỎI: CostCalculatorService tính cost cho batch thời điểm nào?

**TRẢ LỜI NGẮN**: `CostCalculatorService` KHÔNG tính cost cho batch. Nó chỉ **đọc cost đã được tính sẵn** từ batch record.

---

## 📊 PHÂN BIỆT 2 LOẠI COST CALCULATION

### 1. **Batch Cost** (Chi phí sản xuất batch)
- **Ai tính**: `BatchCostCalculator`
- **Khi nào**: Khi TẠO batch record (CreateBatch)
- **Tính từ đâu**: Nguyên liệu thô (raw ingredients)
- **Lưu ở đâu**: `batch_records.cost_per_unit`, `batch_records.total_cost`

### 2. **Menu Item Cost** (Chi phí món ăn)
- **Ai tính**: `CostCalculatorService`
- **Khi nào**: Khi tính chi phí món (CalculateMenuItemCost)
- **Tính từ đâu**: 
  - Raw ingredients → Lấy từ `ingredients.cost_per_unit`
  - Batch ingredients → Lấy từ `batch_records.cost_per_unit` (đã tính sẵn)
- **Lưu ở đâu**: `menu_items.current_cost`

---

## 🔄 FLOW 1: TẠO BATCH (Tính cost cho batch)

```
User tạo batch "cfe cot" (1L)
    ↓
BatchRecordService.CreateBatch()
    ↓
Step 1: Fetch batch definition
    ↓
Step 2: Look up username
    ↓
Step 3: BatchCostCalculator.CalculateBatchCost() ← TÍNH COST Ở ĐÂY
    ↓
    Công thức:
    For each ingredient in batch definition:
        quantity_needed = (quantityProduced / batch_quantity) × source_quantity
        actual_quantity = quantity_needed × (1 + wastage_rate)
        quantity_in_stock_unit = convertUnit(actual_quantity, source_unit, stock_unit)
        ingredient_cost = quantity_in_stock_unit × cost_per_unit
    
    total_cost = sum(all ingredient_costs)
    cost_per_unit = total_cost / quantityProduced
    ↓
Step 4: Check ingredient availability
    ↓
Step 5: Start MongoDB transaction
    ↓
Step 6: Deduct ingredients from inventory
    ↓
Step 7: Calculate expiry time
    ↓
Step 8: Create batch record
    BatchRecord {
        cost_per_unit: 36 đ/ml    ← LƯU COST ĐÃ TÍNH
        total_cost: 36,000 đ
        quantity_produced: 1000 ml
        ...
    }
    ↓
Step 9: Commit transaction
    ↓
Batch "cfe cot" được tạo với cost = 36 đ/ml
```

### Ví dụ tính cost batch "cfe cot":

```
Batch Definition: "cfe cot" (1L)
Ingredients:
  - Cà phê hạt: 100g @ 200 đ/g = 20,000 đ
  - Nước: 900ml @ 0.01 đ/ml = 9 đ
  - Đường: 50g @ 20 đ/g = 1,000 đ
  - Wastage: 5%

Calculation:
  Raw cost = 20,000 + 9 + 1,000 = 21,009 đ
  With wastage = 21,009 × 1.05 = 22,059 đ
  Cost per unit = 22,059 / 1000ml = 22.06 đ/ml

Result:
  batch_records.cost_per_unit = 22.06 đ/ml
  batch_records.total_cost = 22,059 đ
```

---

## 🔄 FLOW 2: TÍNH CHI PHÍ MÓN (Đọc cost từ batch)

```
User xem menu costs hoặc tạo order
    ↓
CostCalculatorService.CalculateMenuItemCost()
    ↓
Fetch menu item "áddd"
    ↓
For each ingredient in menu item:
    ↓
    Is batch ingredient? (IsBatchIngredient())
    ↓
    YES → getBatchCostPerUnit() ← ĐỌC COST TỪ BATCH
        ↓
        batchRecRepo.FindAvailableByDefinition(batchDefID)
        ↓
        Get oldest batch (FIFO)
        ↓
        Return: batch.CostPerUnit (36 đ/ml) ← ĐÃ TÍNH SẴN
        ↓
        Calculate: 200ml × 36 đ/ml = 7,200 đ
    ↓
    NO → Get from ingredients.cost_per_unit
        ↓
        Calculate with conversion & wastage
    ↓
Sum all ingredient costs
    ↓
Update menu_items.current_cost = 7,200 đ
```

---

## 🔑 ĐIỂM QUAN TRỌNG

### 1. **Batch cost được tính 1 LẦN duy nhất**
- Khi tạo batch record
- Dựa trên giá nguyên liệu thô tại thời điểm đó
- Lưu vào `batch_records.cost_per_unit`
- **KHÔNG thay đổi** khi giá nguyên liệu thô thay đổi

### 2. **Menu cost được tính NHIỀU LẦN**
- Mỗi khi gọi `CalculateMenuItemCost()`
- Khi xem menu costs page
- Khi tạo order (nếu có auto-calculate)
- Khi nguyên liệu thay đổi giá

### 3. **Batch cost = Historical cost (Chi phí lịch sử)**
```
Batch "cfe cot" tạo ngày 20/02:
  - Cà phê hạt: 200 đ/g
  - Cost per unit: 22 đ/ml

Ngày 25/02, cà phê hạt tăng giá lên 250 đ/g:
  - Batch cũ vẫn giữ cost: 22 đ/ml ← KHÔNG ĐỔI
  - Batch mới sẽ có cost: 27.5 đ/ml ← TÍNH LẠI
```

### 4. **Menu cost = Current cost (Chi phí hiện tại)**
```
Món "áddd" dùng batch "cfe cot":
  - Nếu còn batch cũ (22 đ/ml) → Cost = 200ml × 22 = 4,400 đ
  - Nếu hết batch cũ, dùng batch mới (27.5 đ/ml) → Cost = 200ml × 27.5 = 5,500 đ
```

---

## 🎨 DIAGRAM: COST CALCULATION ARCHITECTURE

```
┌─────────────────────────────────────────────────────────────┐
│                    BATCH CREATION                            │
│                                                              │
│  User creates batch                                          │
│         ↓                                                    │
│  BatchRecordService.CreateBatch()                           │
│         ↓                                                    │
│  BatchCostCalculator.CalculateBatchCost()                   │
│         ↓                                                    │
│  Calculate from RAW INGREDIENTS                             │
│    - Coffee: 100g @ 200đ/g = 20,000đ                       │
│    - Water: 900ml @ 0.01đ/ml = 9đ                          │
│    - Sugar: 50g @ 20đ/g = 1,000đ                           │
│    - Wastage: 5%                                            │
│         ↓                                                    │
│  Total: 22,059đ / 1000ml = 22.06đ/ml                       │
│         ↓                                                    │
│  SAVE to batch_records                                       │
│    ┌──────────────────────────────┐                        │
│    │ cost_per_unit: 22.06 đ/ml   │ ← STORED ONCE          │
│    │ total_cost: 22,059 đ        │                        │
│    └──────────────────────────────┘                        │
└─────────────────────────────────────────────────────────────┘

┌─────────────────────────────────────────────────────────────┐
│                  MENU COST CALCULATION                       │
│                                                              │
│  User views menu costs / creates order                       │
│         ↓                                                    │
│  CostCalculatorService.CalculateMenuItemCost()              │
│         ↓                                                    │
│  For each ingredient in menu:                               │
│         ↓                                                    │
│  ┌─────────────────────────────────────┐                   │
│  │ Is Batch Ingredient?                │                   │
│  └─────────────────────────────────────┘                   │
│         ↓                    ↓                              │
│       YES                   NO                              │
│         ↓                    ↓                              │
│  getBatchCostPerUnit()   Get from ingredients              │
│         ↓                    ↓                              │
│  READ from batch_records  cost_per_unit                    │
│    ┌──────────────────┐                                    │
│    │ 22.06 đ/ml       │ ← READ STORED VALUE               │
│    └──────────────────┘                                    │
│         ↓                                                    │
│  Calculate: 200ml × 22.06 = 4,412đ                         │
│         ↓                                                    │
│  Sum all ingredients                                         │
│         ↓                                                    │
│  SAVE to menu_items.current_cost                            │
└─────────────────────────────────────────────────────────────┘
```

---

## 💡 TẠI SAO THIẾT KẾ NHƯ VẬY?

### 1. **Accounting Accuracy (Độ chính xác kế toán)**
- Batch cost phản ánh chi phí thực tế khi sản xuất
- Không thay đổi theo giá nguyên liệu sau này
- Đúng với nguyên tắc kế toán FIFO

### 2. **Performance (Hiệu suất)**
- Không cần tính lại cost mỗi lần dùng batch
- Chỉ cần đọc `cost_per_unit` đã lưu
- Nhanh hơn nhiều so với tính từ nguyên liệu thô

### 3. **Inventory Tracking (Theo dõi tồn kho)**
- Mỗi batch có cost riêng
- Biết chính xác lãi/lỗ từng batch
- Dễ dàng phân tích hiệu quả sản xuất

---

## 🔧 TROUBLESHOOTING

### Vấn đề: Menu item dùng batch không tính được cost

**Nguyên nhân có thể**:
1. ❌ `batchRecRepo` không được wire vào `CostCalculatorService`
2. ❌ Batch record không tồn tại hoặc đã hết
3. ❌ `ingredient_type` không được set = "batch"
4. ❌ `batch_id` không được set trong menu ingredient

**Giải pháp**:
1. ✅ Wire `batchRecRepo` trong `main.go` (đã fix)
2. ✅ Kiểm tra batch còn available không
3. ✅ Set `ingredient_type = "batch"` khi thêm vào menu
4. ✅ Set `batch_id` = batch definition ID

---

## 📚 SUMMARY

| Aspect | Batch Cost | Menu Cost |
|--------|-----------|-----------|
| **Tính bởi** | BatchCostCalculator | CostCalculatorService |
| **Khi nào** | Khi tạo batch | Khi tính chi phí món |
| **Từ đâu** | Raw ingredients | Batch records (đã tính) |
| **Lưu ở đâu** | batch_records | menu_items |
| **Thay đổi** | Không (historical) | Có (current) |
| **Mục đích** | Tracking cost sản xuất | Pricing & profit analysis |
