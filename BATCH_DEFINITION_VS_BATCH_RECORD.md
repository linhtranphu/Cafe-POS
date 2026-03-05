# Phân biệt: Batch Definition vs Batch Record

## 🎯 ĐIỂM QUAN TRỌNG

**Khi tạo menu item với batch ingredient, `batch_id` phải trỏ đến `batch_definition_id`, KHÔNG phải `batch_record_id`.**

## 📊 SO SÁNH 2 KHÁI NIỆM

### 1. **Batch Definition** (Công thức batch)

```
Collection: batch_definitions

Ví dụ: "Cà phê cốt"
{
  _id: "507f1f77bcf86cd799439011",
  name: "Cà phê cốt",
  unit: "L",
  shelf_life_hours: 48,
  conversion_rates: [
    {
      source_ingredient_id: "...",
      source_ingredient_name: "Cà phê hạt",
      source_quantity: 100,
      source_unit: "g",
      batch_quantity: 1000,
      wastage_rate: 0.05
    }
  ]
}
```

**Đặc điểm**:
- ✅ Tạo 1 LẦN duy nhất
- ✅ Là "công thức" để sản xuất batch
- ✅ Không có cost (chỉ có công thức)
- ✅ Không có expiry date
- ✅ Dùng để reference trong menu items

### 2. **Batch Record** (Lô sản xuất thực tế)

```
Collection: batch_records

Ví dụ: Lô cà phê cốt ngày 20/02
{
  _id: "65abc123...",
  batch_definition_id: "507f1f77bcf86cd799439011",  ← Trỏ đến definition
  batch_name: "Cà phê cốt",
  quantity_produced: 1000,
  quantity_remaining: 400,
  unit: "ml",
  cost_per_unit: 36,              ← Có cost cụ thể
  total_cost: 36000,
  prepared_at: "2024-02-20T10:00:00Z",
  expires_at: "2024-02-22T10:00:00Z",  ← Có hạn sử dụng
  status: "available"
}
```

**Đặc điểm**:
- ✅ Tạo NHIỀU LẦN (mỗi lần sản xuất)
- ✅ Là "sản phẩm thực tế" đã được làm
- ✅ Có cost cụ thể (tính từ giá nguyên liệu lúc đó)
- ✅ Có expiry date
- ✅ KHÔNG dùng để reference trong menu items

---

## 🔄 WORKFLOW ĐÚNG

### Bước 1: Tạo Batch Definition (1 lần)

```
Manager tạo công thức "Cà phê cốt"
  ↓
POST /api/manager/batch-definitions
  {
    name: "Cà phê cốt",
    unit: "L",
    shelf_life_hours: 48,
    conversion_rates: [...]
  }
  ↓
Lưu vào batch_definitions
  _id: "507f1f77bcf86cd799439011"
```

### Bước 2: Tạo Menu Item với Batch Ingredient

```
Manager tạo món "áddd"
  ↓
POST /api/manager/menu
  {
    name: "áddd",
    price: 50000,
    ingredients: [
      {
        name: "cfe cot",
        quantity: 200,
        unit: "ml",
        ingredient_type: "batch",
        batch_id: "507f1f77bcf86cd799439011"  ← batch_definition_id
      }
    ]
  }
  ↓
Lưu vào menu_items
```

**Lưu ý**: Lúc này CHƯA CẦN có batch record nào!

### Bước 3: Sản xuất Batch (nhiều lần)

```
Barista chế biến cà phê cốt
  ↓
POST /api/barista/batch-records
  {
    batch_definition_id: "507f1f77bcf86cd799439011",
    quantity_produced: 1000,
    prepared_by: "admin"
  }
  ↓
BatchCostCalculator tính cost
  → cost_per_unit = 36 đ/ml
  ↓
Lưu vào batch_records
  _id: "65abc123..."
```

### Bước 4: Tính Chi Phí Món

```
User xem menu costs
  ↓
GET /api/manager/menu/costs
  ↓
CostCalculatorService.CalculateMenuItemCost()
  ↓
For ingredient "cfe cot":
  → IsBatchIngredient() → true
  → batch_id = "507f1f77bcf86cd799439011"
  ↓
  getBatchCostPerUnit(batch_id)
    ↓
    batchRecRepo.FindAvailableByDefinition("507f1f77bcf86cd799439011")
      ↓
      Tìm tất cả batch_records có:
        - batch_definition_id = "507f1f77bcf86cd799439011"
        - status = "available"
        - quantity_remaining > 0
        - expires_at > now
      ↓
      Sắp xếp theo expires_at (FIFO)
      ↓
      Lấy batch cũ nhất
      ↓
      Return: cost_per_unit = 36 đ/ml
  ↓
  Calculate: 200ml × 36 = 7,200 đ
```

---

## 💡 TẠI SAO THIẾT KẾ NHƯ VẬY?

### 1. **Flexibility (Linh hoạt)**

```
Tạo menu trước, sản xuất batch sau:

Ngày 1: Tạo batch definition "Cà phê cốt"
Ngày 2: Tạo món "áddd" dùng "Cà phê cốt"
        → Món đã có công thức
        → Nhưng chưa tính được cost (chưa có batch record)
Ngày 3: Sản xuất batch "Cà phê cốt"
        → Batch có cost = 36 đ/ml
Ngày 4: Tính cost món "áddd"
        → Cost = 7,200 đ
```

### 2. **Multiple Batches (Nhiều lô)**

```
Batch Definition: "Cà phê cốt" (ID: 507f...)

Batch Record 1 (20/02):
  - cost_per_unit: 36 đ/ml
  - quantity_remaining: 400ml
  - expires_at: 22/02

Batch Record 2 (21/02):
  - cost_per_unit: 38 đ/ml  ← Giá nguyên liệu tăng
  - quantity_remaining: 1000ml
  - expires_at: 23/02

Menu Item "áddd":
  - batch_id: 507f...  ← Trỏ đến definition
  - Khi tính cost:
    → Dùng batch cũ nhất (FIFO)
    → Cost = 200ml × 36 = 7,200 đ
```

### 3. **Cost Tracking (Theo dõi chi phí)**

```
Mỗi batch record có cost riêng:
  - Batch 1: 36 đ/ml (giá cà phê 200đ/g)
  - Batch 2: 38 đ/ml (giá cà phê 220đ/g)
  - Batch 3: 35 đ/ml (giá cà phê 190đ/g)

Menu item luôn trỏ đến definition:
  → Tự động dùng batch available
  → Cost thay đổi theo batch đang dùng
  → Accurate cost tracking
```

---

## 🔴 VẤN ĐỀ VỚI MÓN "áddd"

### Scenario 1: Ingredient type không đúng

```javascript
// SAI:
{
  name: "cfe cot",
  ingredient_type: "",        // Mặc định = "raw"
  batch_id: null
}

// Logic:
IsBatchIngredient() → false
→ Tìm trong ingredients collection
→ Không tìm thấy "cfe cot"
→ cost = 0 ❌
```

### Scenario 2: Batch ID không được set

```javascript
// SAI:
{
  name: "cfe cot",
  ingredient_type: "batch",   // ĐÚNG
  batch_id: null              // SAI - Thiếu definition ID
}

// Logic:
IsBatchIngredient() → true
BatchID == nil → true
→ missingIngredients.append("cfe cot")
→ cost = 0 ❌
```

### Scenario 3: Batch ID trỏ sai (trỏ đến record thay vì definition)

```javascript
// SAI:
{
  name: "cfe cot",
  ingredient_type: "batch",
  batch_id: "65abc123..."     // SAI - Đây là batch_record_id
}

// Logic:
IsBatchIngredient() → true
getBatchCostPerUnit("65abc123...")
→ FindAvailableByDefinition("65abc123...")
→ Không tìm thấy batch nào có batch_definition_id = "65abc123..."
→ Error: "no available batches"
→ cost = 0 ❌
```

### Scenario 4: ĐÚNG

```javascript
// ĐÚNG:
{
  name: "cfe cot",
  ingredient_type: "batch",
  batch_id: "507f1f77bcf86cd799439011"  // batch_definition_id
}

// Logic:
IsBatchIngredient() → true
getBatchCostPerUnit("507f1f77bcf86cd799439011")
→ FindAvailableByDefinition("507f1f77bcf86cd799439011")
→ Tìm thấy batch records:
    [
      { cost_per_unit: 36, expires_at: "22/02" },
      { cost_per_unit: 38, expires_at: "23/02" }
    ]
→ Lấy batch cũ nhất (FIFO): cost_per_unit = 36
→ Calculate: 200ml × 36 = 7,200 đ ✅
```

---

## 📝 CHECKLIST FIX MÓN "áddd"

### 1. Kiểm tra Batch Definition tồn tại

```bash
# Tìm batch definition "cfe cot"
db.batch_definitions.findOne({ name: "cfe cot" })

# Lấy _id
# Ví dụ: "507f1f77bcf86cd799439011"
```

### 2. Kiểm tra Batch Records available

```bash
# Tìm batch records của definition này
db.batch_records.find({
  batch_definition_id: ObjectId("507f1f77bcf86cd799439011"),
  status: "available",
  quantity_remaining: { $gt: 0 }
})

# Phải có ít nhất 1 batch available
```

### 3. Update Menu Item

```bash
# Update ingredient trong món "áddd"
db.menu_items.updateOne(
  { name: "áddd" },
  {
    $set: {
      "ingredients.$[elem].ingredient_type": "batch",
      "ingredients.$[elem].batch_id": ObjectId("507f1f77bcf86cd799439011")
    }
  },
  {
    arrayFilters: [{ "elem.name": "cfe cot" }]
  }
)
```

### 4. Trigger Cost Calculation

```bash
# Gọi API tính lại chi phí
POST /api/manager/menu/:menu_item_id/calculate-cost
```

---

## 🎯 KẾT LUẬN

**Batch ID trong menu item = Batch Definition ID**

- ✅ Cho phép tạo menu trước khi có batch thực tế
- ✅ Tự động dùng batch available (FIFO)
- ✅ Cost thay đổi theo batch đang dùng
- ✅ Flexible và accurate cost tracking

**Món "áddd" không có cost vì**:
1. `ingredient_type` không phải "batch"
2. `batch_id` không được set hoặc set sai
3. Cần update data và trigger recalculation
