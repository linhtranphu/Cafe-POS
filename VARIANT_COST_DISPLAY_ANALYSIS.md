# Phân tích: Hiển thị Chi phí ước tính khi Edit Variant

## Vấn đề

Khi thêm ingredient vào variant, chi phí ước tính được tính toán và hiển thị. Nhưng khi load lại (edit menu item), chi phí không hiển thị.

## Phân tích

### Tình huống hiện tại:

**Khi thêm ingredient mới:**
```javascript
// Add ingredient to variant
variant.ingredients.push({
  // ... other fields
  estimatedCost: 0
})

// Calculate initial cost
updateVariantIngredientCost(currentVariantIndex.value, ingIndex)
// → estimatedCost được tính toán ngay
```

**Khi load lại (edit):**
```javascript
// OLD CODE - estimatedCost luôn = 0
return {
  // ... other fields
  estimatedCost: 0  // ❌ Không tính toán
}
```

### Câu hỏi: Có hợp lý khi hiển thị chi phí ước tính khi edit?

**Trả lời: CÓ, rất hợp lý!**

## Lý do nên hiển thị:

### 1. Tính nhất quán UX
- Nếu hiển thị khi tạo mới → cũng nên hiển thị khi edit
- Người dùng mong đợi thấy cùng thông tin trong cả 2 trường hợp

### 2. Giúp người dùng ra quyết định
- Khi edit, người dùng có thể:
  - Thấy ngay chi phí thay đổi khi điều chỉnh số lượng
  - Thấy ngay chi phí thay đổi khi đổi đơn vị
  - So sánh chi phí giữa các size
  - Đánh giá giá bán có hợp lý không

### 3. Thông tin có sẵn
- Tất cả dữ liệu cần thiết đều có:
  - `costPerUnit` từ ingredient store
  - `quantity` từ menu item
  - `unit` từ menu item
  - `wastage` từ ingredient store
- Không cần gọi API thêm

### 4. Tính năng hữu ích
```
Ví dụ thực tế:
┌─────────────────────────────────────────┐
│ Cà phê sữa đá - Size M                  │
│                                         │
│ Nguyên liệu:                            │
│ • Cà phê: 20g                           │
│   Chi phí ước tính: 15,000 VNĐ         │ ← Giúp người dùng
│                                         │
│ • Sữa: 50ml                             │
│   Chi phí ước tính: 8,000 VNĐ          │ ← Đánh giá nhanh
│                                         │
│ 💰 Tổng chi phí: 23,000 VNĐ            │
│ Giá bán: 30,000 VNĐ                    │
│ → Lợi nhuận: 7,000 VNĐ (23%)           │ ← Quyết định có hợp lý
└─────────────────────────────────────────┘
```

## Giải pháp đã implement

### Code mới:
```javascript
const editItem = (item) => {
  // ...
  const enrichedIngredients = variant.ingredients.map(ing => {
    const ingredientData = availableIngredients.value.find(i => i.name === ing.name)
    if (ingredientData) {
      // ✅ Calculate estimated cost when loading
      const breakdown = calculateCostBreakdown(
        ing.quantity,
        ing.unit,
        ingredientData.cost_per_unit || 0,
        ingredientData.unit,
        ingredientData.wastage_percentage || 0
      )
      
      return {
        // ... other fields
        estimatedCost: breakdown.totalCost  // ✅ Tính toán ngay
      }
    }
    return ing
  })
  // ...
}
```

## Lợi ích

### 1. UX nhất quán
- ✅ Chi phí hiển thị cả khi tạo mới và edit
- ✅ Người dùng không bị bối rối

### 2. Thông tin real-time
- ✅ Chi phí được tính dựa trên giá nguyên liệu hiện tại
- ✅ Nếu giá nguyên liệu thay đổi, chi phí ước tính cũng thay đổi khi edit

### 3. Hỗ trợ ra quyết định
- ✅ Người dùng thấy ngay tổng chi phí của mỗi size
- ✅ Dễ dàng so sánh và điều chỉnh giá bán

### 4. Không có nhược điểm
- ✅ Không tốn thêm API call
- ✅ Không làm chậm UI
- ✅ Tính toán nhanh (client-side)

## Kết luận

**Quyết định: Hiển thị chi phí ước tính khi edit là HỢP LÝ và CẦN THIẾT**

Lý do:
1. Tính nhất quán UX
2. Thông tin hữu ích cho người dùng
3. Dữ liệu có sẵn, không tốn chi phí
4. Giúp người dùng ra quyết định tốt hơn về giá bán

## Testing

### Manual Test:
1. Tạo menu item multi-size với ingredients có cost
2. Lưu menu item
3. Edit lại menu item
4. Kiểm tra:
   - ✅ Chi phí ước tính hiển thị cho mỗi ingredient
   - ✅ Tổng chi phí hiển thị cho mỗi variant
   - ✅ Thay đổi số lượng → chi phí cập nhật ngay
   - ✅ Thay đổi đơn vị → chi phí cập nhật ngay

### Expected Result:
```
Size M:
  Cà phê: 20g @ 750 VNĐ/g
  Chi phí ước tính: 15,000 VNĐ ✅
  
  Sữa: 50ml @ 160 VNĐ/ml
  Chi phí ước tính: 8,000 VNĐ ✅
  
  💰 Tổng chi phí: 23,000 VNĐ ✅
```

## Files thay đổi

- `frontend/src/views/MenuView.vue`: Tính toán estimatedCost trong editItem()
- `VARIANT_COST_DISPLAY_ANALYSIS.md`: Tài liệu phân tích này
