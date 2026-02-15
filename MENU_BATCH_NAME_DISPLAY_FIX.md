# Menu Batch Name Display Fix

## Vấn đề
Khi thêm batch vào menu mới tại `http://localhost:5173/#/menu`, tên batch không hiển thị trong danh sách nguyên liệu.

## Nguyên nhân có thể
1. Field `name` của batch definition không được set đúng
2. Template không có fallback khi `name` undefined
3. Khó phân biệt batch vs raw ingredient

## Giải pháp

### 1. Thêm Fallback cho Name
Thay vì chỉ hiển thị `{{ ing.name }}`, thêm fallback:
```vue
{{ ing.name || 'Không có tên' }}
```

### 2. Thêm Badge phân biệt Batch
Thêm badge "🧪 Batch" để dễ nhận biết:
```vue
<div class="flex items-center gap-2">
  <div class="font-bold text-sm text-gray-800 truncate">{{ ing.name || 'Không có tên' }}</div>
  <span v-if="ing.type === 'batch'" class="text-xs bg-purple-100 text-purple-700 px-2 py-0.5 rounded-full flex-shrink-0">
    🧪 Batch
  </span>
</div>
```

## Code Changes

### Single-Size Ingredients (Trước)
```vue
<div class="font-bold text-sm text-gray-800 truncate">{{ ing.name }}</div>
```

### Single-Size Ingredients (Sau)
```vue
<div class="flex items-center gap-2">
  <div class="font-bold text-sm text-gray-800 truncate">{{ ing.name || 'Không có tên' }}</div>
  <span v-if="ing.type === 'batch'" class="text-xs bg-purple-100 text-purple-700 px-2 py-0.5 rounded-full flex-shrink-0">
    🧪 Batch
  </span>
</div>
```

### Variant Ingredients
Áp dụng tương tự cho variant ingredients.

## Lợi ích

### 1. Fallback Name
- ✅ Hiển thị "Không có tên" nếu name undefined
- ✅ Dễ debug khi có vấn đề
- ✅ Không bị blank space

### 2. Badge Phân biệt
- ✅ Dễ nhận biết batch vs raw ingredient
- ✅ UI rõ ràng hơn
- ✅ Màu tím (purple) nhất quán với batch theme

## Testing

### Test Case 1: Thêm Raw Ingredient
1. Mở menu form
2. Thêm nguyên liệu thô (ví dụ: Cà phê hạt)
3. **Kết quả mong đợi**: Hiển thị tên, không có badge

### Test Case 2: Thêm Batch
1. Mở menu form
2. Chuyển sang tab "🧪 Batch"
3. Chọn một batch definition
4. **Kết quả mong đợi**: 
   - Hiển thị tên batch
   - Có badge "🧪 Batch" màu tím

### Test Case 3: Batch không có tên
1. Nếu batch.name undefined
2. **Kết quả mong đợi**: Hiển thị "Không có tên"

## Debug Steps

Nếu vẫn không hiển thị tên:

### 1. Kiểm tra Console
```javascript
// Trong selectBatch function, thêm:
console.log('Selected batch:', batch)
console.log('Batch name:', batch.name)
```

### 2. Kiểm tra Batch Definition API
```bash
curl http://localhost:3000/api/batch/definitions
```

Xem response có field `name` không:
```json
{
  "id": "...",
  "name": "Cà phê Concentrate",  // ← Phải có field này
  "unit": "ml",
  ...
}
```

### 3. Kiểm tra Store
```javascript
// Trong browser console:
const batchStore = useBatchDefinitionStore()
console.log('Definitions:', batchStore.definitions)
```

## Files Modified

### frontend/src/views/MenuView.vue
1. Single-size ingredient header (line ~303)
2. Variant ingredient header (line ~467)

## Related Issues

### Nếu batch.name vẫn undefined
Có thể do:
1. Backend không trả về field `name`
2. Store không map đúng
3. Service không parse response đúng

**Giải pháp**: Kiểm tra chain:
```
Backend API → Service → Store → Component
```

## Visual Changes

### Before
```
┌─────────────────────────────┐
│                             │  ← Blank (no name)
│ Kho: ml @ 0đ/ml            │
└─────────────────────────────┘
```

### After (Raw Ingredient)
```
┌─────────────────────────────┐
│ Cà phê hạt                  │
│ Kho: kg @ 200,000đ/kg      │
└─────────────────────────────┘
```

### After (Batch)
```
┌─────────────────────────────┐
│ Cà phê Concentrate 🧪 Batch │
│ Kho: ml @ 0đ/ml            │
└─────────────────────────────┘
```

## Status

✅ **FIXED** - Thêm fallback và badge để hiển thị tên batch rõ ràng hơn

## Next Steps

1. Test thêm batch vào menu
2. Verify tên hiển thị đúng
3. Verify badge "🧪 Batch" xuất hiện
4. Nếu vẫn không hiển thị → debug theo steps trên
