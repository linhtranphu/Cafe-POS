# Chiến lược Kế toán Chi phí Món

## 🎯 CÂU HỎI QUAN TRỌNG

1. **Logic chi phí theo món như thế nào?**
2. **Có nên lưu lịch sử chi phí theo ngày?**
3. **Nếu 1 ngày có 2 batch với chi phí khác nhau thì lấy trung bình?**

---

## 📊 LOGIC HIỆN TẠI (Current Implementation)

### 1. **Menu Item Cost = Current Cost (Chi phí hiện tại)**

```
menu_items collection:
{
  _id: "...",
  name: "áddd",
  price: 50000,
  current_cost: 7200,              ← CHỈ LƯU 1 GIÁ TRỊ
  cost_last_calculated_at: "...",
  cost_status: "FINAL"
}
```

**Đặc điểm**:
- ✅ Chỉ lưu chi phí hiện tại
- ✅ Được cập nhật khi tính lại
- ❌ KHÔNG lưu lịch sử
- ❌ KHÔNG theo dõi theo ngày

### 2. **Order Item Cost = Snapshot Cost (Chi phí tại thời điểm bán)**

```
order_items collection:
{
  _id: "...",
  order_id: "...",
  menu_item_id: "...",
  name: "áddd",
  price: 50000,
  cost: 7200,                      ← LƯU CHI PHÍ LÚC BÁN
  created_at: "2024-02-20T10:00:00Z"
}
```

**Đặc điểm**:
- ✅ Lưu chi phí tại thời điểm order
- ✅ Immutable (không thay đổi)
- ✅ Dùng để tính profit thực tế
- ✅ Có timestamp

### 3. **Batch Record Cost = Historical Cost (Chi phí lịch sử)**

```
batch_records collection:
{
  _id: "...",
  batch_definition_id: "...",
  cost_per_unit: 36,               ← CHI PHÍ CỐ ĐỊNH
  total_cost: 36000,
  prepared_at: "2024-02-20T10:00:00Z",
  expires_at: "2024-02-22T10:00:00Z"
}
```

**Đặc điểm**:
- ✅ Mỗi batch có chi phí riêng
- ✅ Immutable (không thay đổi)
- ✅ Phản ánh giá nguyên liệu lúc sản xuất
- ✅ FIFO inventory

---

## 🔄 LUỒNG CHI PHÍ HIỆN TẠI

### Scenario: Giá nguyên liệu thay đổi

```
Timeline:

Ngày 20/02 - Sáng:
  Cà phê hạt: 200 đ/g
  Tạo batch 1: cost_per_unit = 36 đ/ml
  Món "áddd": current_cost = 7,200 đ (200ml × 36)

Ngày 20/02 - Chiều:
  Cà phê hạt tăng giá: 250 đ/g
  Tạo batch 2: cost_per_unit = 45 đ/ml
  Món "áddd": current_cost = 7,200 đ (vẫn dùng batch 1)

Ngày 21/02:
  Batch 1 hết
  Món "áddd": current_cost = 9,000 đ (200ml × 45, dùng batch 2)
  
Orders:
  Order 1 (20/02 sáng): cost = 7,200 đ ✅
  Order 2 (20/02 chiều): cost = 7,200 đ ✅
  Order 3 (21/02): cost = 9,000 đ ✅
```

**Kết quả**:
- ✅ Mỗi order lưu chi phí chính xác tại thời điểm bán
- ✅ Profit analysis chính xác
- ❌ Menu cost không có lịch sử

---

## 💡 CÁC PHƯƠNG ÁN KẾ TOÁN

### Phương án 1: CURRENT COST ONLY (Hiện tại)

```javascript
menu_items: {
  current_cost: 7200,
  cost_last_calculated_at: "2024-02-20T10:00:00Z"
}

order_items: {
  cost: 7200,  // Snapshot tại thời điểm bán
  created_at: "2024-02-20T10:00:00Z"
}
```

**Ưu điểm**:
- ✅ Đơn giản, dễ implement
- ✅ Performance tốt
- ✅ Order items có lịch sử chính xác

**Nhược điểm**:
- ❌ Không biết chi phí món thay đổi như thế nào
- ❌ Không phân tích được trend chi phí
- ❌ Khó so sánh chi phí theo thời gian

**Phù hợp khi**:
- Giá nguyên liệu ổn định
- Chỉ quan tâm profit hiện tại
- Không cần phân tích lịch sử chi phí món

---

### Phương án 2: DAILY COST HISTORY (Lịch sử theo ngày)

```javascript
menu_items: {
  current_cost: 9000,
  cost_last_calculated_at: "2024-02-21T10:00:00Z"
}

menu_cost_history: {
  menu_item_id: "...",
  date: "2024-02-20",
  cost: 7200,
  cost_status: "FINAL",
  calculated_at: "2024-02-20T10:00:00Z"
}
```

**Ưu điểm**:
- ✅ Theo dõi được trend chi phí
- ✅ Phân tích được biến động
- ✅ So sánh chi phí theo thời gian
- ✅ Audit trail đầy đủ

**Nhược điểm**:
- ❌ Phức tạp hơn
- ❌ Cần storage nhiều hơn
- ❌ Query phức tạp hơn

**Phù hợp khi**:
- Giá nguyên liệu thay đổi thường xuyên
- Cần phân tích trend
- Cần báo cáo chi phí theo thời gian

---

### Phương án 3: SNAPSHOT ON CHANGE (Snapshot khi thay đổi)

```javascript
menu_items: {
  current_cost: 9000,
  cost_last_calculated_at: "2024-02-21T10:00:00Z"
}

menu_cost_snapshots: {
  menu_item_id: "...",
  cost: 7200,
  valid_from: "2024-02-20T10:00:00Z",
  valid_to: "2024-02-21T10:00:00Z",
  reason: "Batch cost changed"
}
```

**Ưu điểm**:
- ✅ Chỉ lưu khi có thay đổi
- ✅ Tiết kiệm storage
- ✅ Có lịch sử đầy đủ
- ✅ Query hiệu quả

**Nhược điểm**:
- ❌ Logic phức tạp hơn
- ❌ Cần detect changes

**Phù hợp khi**:
- Giá thay đổi không thường xuyên
- Cần balance giữa history và performance
- Cần biết tại sao chi phí thay đổi

---

## 🎲 XỬ LÝ NHIỀU BATCH CÙNG NGÀY

### Scenario: 1 ngày có 2 batch với chi phí khác nhau

```
Ngày 20/02:

Batch 1 (Sáng):
  - cost_per_unit: 36 đ/ml
  - quantity: 1000ml
  - prepared_at: 08:00

Batch 2 (Chiều):
  - cost_per_unit: 45 đ/ml
  - quantity: 1000ml
  - prepared_at: 14:00

Orders:
  Order 1 (10:00): Dùng batch 1 → cost = 7,200 đ
  Order 2 (16:00): Dùng batch 1 → cost = 7,200 đ (FIFO)
  Order 3 (18:00): Batch 1 hết, dùng batch 2 → cost = 9,000 đ
```

### Option A: FIFO (First In First Out) - HIỆN TẠI

```go
// Logic hiện tại
getBatchCostPerUnit(batchDefID) {
  batches := FindAvailableByDefinition(batchDefID)
  // Sắp xếp theo expires_at (batch cũ nhất)
  sort(batches, by: expires_at)
  return batches[0].cost_per_unit  // 36 đ/ml
}
```

**Kết quả**:
- Order 1: 7,200 đ (batch 1)
- Order 2: 7,200 đ (batch 1)
- Order 3: 9,000 đ (batch 2)

**Ưu điểm**:
- ✅ Chính xác với inventory thực tế
- ✅ Đúng nguyên tắc kế toán FIFO
- ✅ Mỗi order có chi phí đúng

**Nhược điểm**:
- ❌ Chi phí món thay đổi trong ngày
- ❌ Khó dự đoán chi phí

---

### Option B: WEIGHTED AVERAGE (Trung bình gia quyền)

```go
getBatchCostPerUnit(batchDefID) {
  batches := FindAvailableByDefinition(batchDefID)
  
  totalCost := 0
  totalQuantity := 0
  
  for batch in batches {
    totalCost += batch.cost_per_unit * batch.quantity_remaining
    totalQuantity += batch.quantity_remaining
  }
  
  return totalCost / totalQuantity
}
```

**Tính toán**:
```
Batch 1: 36 đ/ml × 1000ml = 36,000 đ
Batch 2: 45 đ/ml × 1000ml = 45,000 đ
Total: 81,000 đ / 2000ml = 40.5 đ/ml
```

**Kết quả**:
- Order 1: 8,100 đ (40.5 × 200ml)
- Order 2: 8,100 đ (40.5 × 200ml)
- Order 3: 8,100 đ (40.5 × 200ml)

**Ưu điểm**:
- ✅ Chi phí ổn định trong ngày
- ✅ Dễ dự đoán
- ✅ Công bằng

**Nhược điểm**:
- ❌ Không phản ánh inventory thực tế
- ❌ Không đúng FIFO
- ❌ Profit analysis không chính xác

---

### Option C: DAILY AVERAGE (Trung bình theo ngày)

```go
getDailyCost(menuItemID, date) {
  orders := FindOrdersByDate(menuItemID, date)
  
  totalCost := 0
  totalOrders := 0
  
  for order in orders {
    totalCost += order.cost
    totalOrders++
  }
  
  return totalCost / totalOrders
}
```

**Kết quả**:
```
Order 1: 7,200 đ (FIFO - batch 1)
Order 2: 7,200 đ (FIFO - batch 1)
Order 3: 9,000 đ (FIFO - batch 2)

Daily average: (7,200 + 7,200 + 9,000) / 3 = 7,800 đ
```

**Ưu điểm**:
- ✅ Phản ánh chi phí thực tế trong ngày
- ✅ Dùng cho báo cáo
- ✅ Smooth out variations

**Nhược điểm**:
- ❌ Chỉ tính được sau khi có orders
- ❌ Không dùng được cho pricing

---

## 🎯 KHUYẾN NGHỊ

### Cho hệ thống hiện tại (Cafe POS):

**Giữ nguyên FIFO + Thêm Cost History**

```javascript
// 1. Menu Items - Current cost (như hiện tại)
menu_items: {
  current_cost: 9000,
  cost_last_calculated_at: "2024-02-21T10:00:00Z"
}

// 2. Order Items - Snapshot cost (như hiện tại)
order_items: {
  cost: 7200,  // Chi phí thực tế lúc bán
  created_at: "2024-02-20T10:00:00Z"
}

// 3. NEW: Menu Cost History (thêm mới)
menu_cost_history: {
  menu_item_id: "...",
  date: "2024-02-20",
  min_cost: 7200,
  max_cost: 7200,
  avg_cost: 7200,
  sample_count: 10,  // Số lần tính trong ngày
  calculated_at: "2024-02-20T23:59:59Z"
}
```

### Lý do:

1. **FIFO cho orders**:
   - ✅ Chính xác với inventory
   - ✅ Đúng nguyên tắc kế toán
   - ✅ Profit analysis chính xác

2. **Daily history cho analysis**:
   - ✅ Theo dõi trend
   - ✅ Phân tích biến động
   - ✅ Báo cáo theo thời gian

3. **Min/Max/Avg cho insights**:
   - ✅ Biết range chi phí trong ngày
   - ✅ Detect anomalies
   - ✅ Planning & forecasting

---

## 📊 IMPLEMENTATION STRATEGY

### Phase 1: Giữ nguyên (Đã có)
```
✅ Menu current_cost
✅ Order snapshot cost
✅ Batch FIFO logic
```

### Phase 2: Thêm History (Tùy chọn)
```
□ Tạo menu_cost_history collection
□ Background job tính daily summary
□ API endpoints cho history
□ UI charts cho trend analysis
```

### Phase 3: Advanced Analytics (Tương lai)
```
□ Cost forecasting
□ Price optimization
□ Margin analysis by time
□ Ingredient cost impact analysis
```

---

## 🎓 KẾT LUẬN

### Trả lời câu hỏi:

**1. Logic chi phí theo món như thế nào?**
- Hiện tại: FIFO - Dùng batch cũ nhất
- Mỗi order lưu chi phí tại thời điểm bán
- Menu chỉ lưu current cost

**2. Có nên lưu lịch sử chi phí theo ngày?**
- ✅ NÊN - Nếu cần phân tích trend
- ✅ NÊN - Nếu giá thay đổi thường xuyên
- ❌ KHÔNG CẦN - Nếu chỉ quan tâm profit hiện tại

**3. Nếu 1 ngày có 2 batch khác chi phí thì lấy trung bình?**
- ❌ KHÔNG - Giữ FIFO cho orders (chính xác)
- ✅ CÓ - Lưu average trong daily history (analysis)
- ✅ CÓ - Lưu cả min/max để biết range

### Khuyến nghị cuối cùng:

**Giữ FIFO + Thêm Daily Summary**
- Orders: FIFO (chính xác)
- History: Daily avg/min/max (analysis)
- Best of both worlds!
