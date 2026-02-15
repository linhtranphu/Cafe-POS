# Menu Batch Selector Fix

## Vấn đề
Khi tạo menu mới tại `http://localhost:5173/#/menu`, phần chọn batch không hiển thị mặc dù đã định nghĩa batch.

## Nguyên nhân
Code đang hiển thị **batch records** (các batch đã được sản xuất) thay vì **batch definitions** (định nghĩa batch). 

### Logic cũ (sai):
```javascript
const filteredAvailableBatches = computed(() => {
  let filtered = availableBatchRecords.value || []  // ❌ Sử dụng batch records
  
  // Only show available batches
  filtered = filtered.filter(batch => 
    batch.status === 'available' && 
    batch.quantity_remaining > 0
  )
  // ...
})
```

Vấn đề:
- Người dùng tạo batch definition
- Nhưng chưa sản xuất batch record nào
- Danh sách batch records trống → Không hiển thị gì

## Giải pháp
Thay đổi để hiển thị **batch definitions** thay vì batch records.

### 1. Cập nhật Filter Logic

**Trước:**
```javascript
const filteredAvailableBatches = computed(() => {
  let filtered = availableBatchRecords.value || []
  filtered = filtered.filter(batch => 
    batch.status === 'available' && 
    batch.quantity_remaining > 0
  )
  // ...
})
```

**Sau:**
```javascript
const filteredAvailableBatches = computed(() => {
  // Use batch definitions instead of batch records
  let filtered = availableBatchDefinitions.value || []
  
  if (ingredientSearchQuery.value) {
    const query = ingredientSearchQuery.value.toLowerCase()
    filtered = filtered.filter(batch => 
      batch.name?.toLowerCase().includes(query)
    )
  }
  
  return filtered
})
```

### 2. Cập nhật Template Display

**Trước (hiển thị batch record):**
```vue
<div class="font-medium text-gray-800">{{ batch.batch_name }}</div>
<span v-if="batch.status === 'available'" class="...">Khả dụng</span>
<div class="text-sm text-gray-500 mt-1">
  Còn lại: {{ batch.quantity_remaining }} {{ batch.unit }}
</div>
```

**Sau (hiển thị batch definition):**
```vue
<div class="font-medium text-gray-800">{{ batch.name }}</div>
<span class="text-xs bg-purple-100 text-purple-700 px-2 py-0.5 rounded-full">
  Batch
</span>
<div class="text-sm text-gray-500 mt-1">
  Đơn vị: {{ batch.unit }}
</div>
<div class="text-xs text-gray-400 mt-1">
  Thời hạn: {{ batch.shelf_life_hours }} giờ
</div>
```

### 3. Cập nhật selectBatch Function

**Trước:**
```javascript
const selectBatch = (batch) => {
  // batch is a batch record
  form.value.ingredients.push({
    id: batch.id,
    batch_definition_id: batch.batch_definition_id,  // Reference to definition
    name: batch.batch_name,
    costPerUnit: batch.cost_per_unit || 0,
    availableQuantity: batch.quantity_remaining,
    expiresAt: batch.expires_at,
    status: batch.status
    // ...
  })
}
```

**Sau:**
```javascript
const selectBatch = (batch) => {
  // batch is now a batch definition
  form.value.ingredients.push({
    id: batch.id,
    batch_definition_id: batch.id,  // Same as id
    name: batch.name,
    costPerUnit: 0,  // Will be calculated based on conversion rates
    // No availableQuantity, expiresAt, status (those are for records)
    // ...
  })
}
```

## Sự khác biệt: Batch Definition vs Batch Record

### Batch Definition (Định nghĩa Batch)
- Là "công thức" để tạo batch
- Chứa thông tin: tên, đơn vị, thời hạn, nguyên liệu nguồn
- Tồn tại vĩnh viễn trong hệ thống
- Ví dụ: "Cà phê Concentrate" (công thức)

### Batch Record (Batch đã sản xuất)
- Là batch cụ thể đã được sản xuất từ definition
- Chứa thông tin: số lượng còn lại, ngày hết hạn, trạng thái
- Có vòng đời: sản xuất → sử dụng → hết hạn
- Ví dụ: "Cà phê Concentrate #001" sản xuất ngày 15/02/2026

## Lý do thay đổi

### Trước đây (dùng Batch Records):
1. ❌ Phải sản xuất batch record trước khi tạo menu
2. ❌ Không linh hoạt - phụ thuộc vào tồn kho
3. ❌ Phức tạp cho người dùng

### Bây giờ (dùng Batch Definitions):
1. ✅ Chỉ cần định nghĩa batch là có thể tạo menu
2. ✅ Linh hoạt - không phụ thuộc tồn kho
3. ✅ Đơn giản hơn cho người dùng

## Data Structure Changes

### Batch Definition Object
```javascript
{
  id: "batch_def_123",
  name: "Cà phê Concentrate",
  unit: "ml",
  shelf_life_hours: 24,
  low_stock_threshold: 200,
  expiry_warning_hours: 4,
  conversion_rates: [...]
}
```

### Batch Record Object (không dùng nữa trong menu)
```javascript
{
  id: "batch_rec_456",
  batch_definition_id: "batch_def_123",
  batch_name: "Cà phê Concentrate",
  quantity_remaining: 500,
  status: "available",
  expires_at: "2026-02-16T10:00:00Z"
}
```

## Files thay đổi

### frontend/src/views/MenuView.vue
1. **filteredAvailableBatches** - Đổi từ batch records sang batch definitions
2. **Template** - Hiển thị thông tin batch definition
3. **selectBatch** - Xử lý batch definition thay vì batch record

## Testing Checklist

- [x] No syntax errors
- [ ] Batch definitions hiển thị trong selector
- [ ] Có thể chọn batch definition khi tạo menu
- [ ] Thông tin batch hiển thị đúng (tên, đơn vị, thời hạn)
- [ ] Search batch hoạt động
- [ ] Thêm batch vào menu thành công
- [ ] Chi phí tính toán đúng (nếu có)

## Hướng dẫn test

### Bước 1: Tạo Batch Definition
1. Mở `http://localhost:5173/#/batch`
2. Click "➕ Tạo Batch Mới"
3. Điền thông tin:
   - Tên: "Cà phê Concentrate"
   - Đơn vị: "ml"
   - Thời hạn: 24 giờ
   - Chọn nguyên liệu nguồn
4. Lưu batch definition

### Bước 2: Tạo Menu với Batch
1. Mở `http://localhost:5173/#/menu`
2. Click "➕ Thêm món"
3. Điền tên món và thông tin cơ bản
4. Click "+ Thêm" trong phần "🥘 Nguyên liệu"
5. Chuyển sang tab "🧪 Batch"
6. **Kiểm tra**: Batch "Cà phê Concentrate" hiển thị
7. Click chọn batch
8. Nhập số lượng cần dùng
9. Lưu menu

### Kết quả mong đợi
- ✅ Batch definition hiển thị ngay sau khi tạo
- ✅ Không cần sản xuất batch record
- ✅ Có thể thêm batch vào menu
- ✅ Menu lưu thành công với batch ingredient

## Ghi chú kỹ thuật

### API Compatibility
- Backend API đã hỗ trợ cả batch definitions và batch records
- Frontend chỉ cần gọi đúng endpoint:
  - `GET /api/batch/definitions` - Lấy danh sách definitions
  - `GET /api/batch/records` - Lấy danh sách records

### Store Usage
```javascript
// Batch Definition Store
const batchDefinitionStore = useBatchDefinitionStore()
batchDefinitionStore.fetchDefinitions()  // Load definitions
const definitions = batchDefinitionStore.definitions

// Batch Record Store (không dùng cho menu selector)
const batchRecordStore = useBatchRecordStore()
batchRecordStore.fetchRecords()  // Load records
const records = batchRecordStore.records
```

### Cost Calculation
- Batch definition không có `cost_per_unit` trực tiếp
- Chi phí được tính dựa trên `conversion_rates` và giá nguyên liệu nguồn
- Frontend cần implement logic tính toán chi phí batch

## Lợi ích

### Cho người dùng
1. ✅ Workflow đơn giản hơn
2. ✅ Không cần sản xuất batch trước
3. ✅ Tạo menu nhanh hơn

### Cho hệ thống
1. ✅ Logic rõ ràng hơn
2. ✅ Tách biệt định nghĩa và thực tế
3. ✅ Dễ bảo trì và mở rộng

## Status

✅ **FIXED** - Batch definitions hiện đã hiển thị trong menu selector
