# Tính năng Ghi chú món (Item Notes)

## Tổng quan

Đã thêm tính năng cho phép waiter ghi chú yêu cầu đặc biệt cho từng món khi tạo order. Ví dụ: "Ít đường", "Nhiều đá", "Nóng", v.v.

## Thay đổi

### Frontend (`frontend/src/components/CreateOrderModal.vue`)

**1. Thêm nút "Ghi chú" cho mỗi món:**
- Hiển thị khi món đã được thêm vào cart (qty > 0)
- Màu vàng nhạt để dễ nhận biết
- Icon 📝 để trực quan

**2. Thêm Note Modal:**
- Hiển thị tên món đang ghi chú
- Gợi ý nhanh (quick suggestions):
  - Ít đường / Nhiều đường / Không đường
  - Ít đá / Nhiều đá / Không đá
  - Nóng / Ấm / Lạnh
- Ô nhập text tự do
- Nút Lưu/Hủy

**3. State mới:**
```javascript
const itemNotes = ref({}) // { itemId: note, itemId_variantId: note }
const showNoteModal = ref(false)
const currentNote = ref('')
const noteModalItemId = ref(null)
const noteModalVariantId = ref(null)
const noteModalItemName = ref('')
```

**4. Methods mới:**
- `getItemNote(itemId, variantId)` - Lấy note của món
- `openNoteModal(itemId, variantId)` - Mở modal ghi chú
- `closeNoteModal()` - Đóng modal
- `addNoteSuggestion(suggestion)` - Thêm gợi ý vào note
- `saveNote()` - Lưu note

**5. Cập nhật handleConfirm:**
- Include note khi tạo order
- Clear notes sau khi confirm

### Backend

Backend đã hỗ trợ sẵn field `note` trong `OrderItem`:
```go
type OrderItem struct {
    MenuItemID  primitive.ObjectID `bson:"menu_item_id" json:"menu_item_id"`
    VariantID   string             `bson:"variant_id,omitempty" json:"variant_id,omitempty"`
    Name        string             `bson:"name" json:"name"`
    VariantName string             `bson:"variant_name,omitempty" json:"variant_name,omitempty"`
    Price       float64            `bson:"price" json:"price"`
    Quantity    int                `bson:"quantity" json:"quantity"`
    Note        string             `bson:"note,omitempty" json:"note,omitempty"` // ✅ Đã có
    Subtotal    float64            `bson:"subtotal" json:"subtotal"`
}
```

Không cần thay đổi backend!

## Luồng sử dụng

### 1. Thêm món vào cart
- Waiter tap "+" để thêm món
- Nút +/- hiện ra

### 2. Thêm ghi chú
- Tap nút "📝 Thêm ghi chú" (màu vàng)
- Modal ghi chú hiện ra

### 3. Nhập ghi chú
**Option A: Dùng gợi ý nhanh**
- Tap các nút gợi ý: "Ít đường", "Nhiều đá", v.v.
- Có thể tap nhiều gợi ý, sẽ tự động nối bằng dấu phẩy

**Option B: Nhập tự do**
- Gõ trực tiếp vào ô text
- Ví dụ: "Ít đường, nhiều đá, nóng"

### 4. Lưu ghi chú
- Tap "Lưu"
- Note hiện ra dưới nút ghi chú
- Nút đổi thành "📝 Sửa ghi chú"

### 5. Xem/Sửa ghi chú
- Note hiển thị với icon 💬
- Tap "Sửa ghi chú" để chỉnh sửa
- Có thể xóa note bằng cách xóa hết text và tap Lưu

### 6. Xác nhận order
- Tap "Xác nhận" ở floating button
- Note được gửi kèm order đến backend

## UI/UX

### Màu sắc
- Nút ghi chú: Vàng nhạt (`bg-amber-50`, `border-amber-300`)
- Text ghi chú: Vàng đậm (`text-amber-700`)
- Modal: Trắng với accent xanh

### Icon
- Nút: 📝
- Note hiển thị: 💬

### Vị trí
- **Single-size items**: Nút ghi chú nằm dưới nút +/-
- **Multi-size items**: Nút ghi chú nằm dưới mỗi variant có qty > 0

### Responsive
- Modal full-screen trên mobile
- Slide-up animation
- Safe area support

## Gợi ý nhanh (Quick Suggestions)

```javascript
const noteSuggestions = [
  'Ít đường',
  'Nhiều đường',
  'Không đường',
  'Ít đá',
  'Nhiều đá',
  'Không đá',
  'Nóng',
  'Ấm',
  'Lạnh'
]
```

Có thể mở rộng thêm:
- Độ ngọt: Ít ngọt, Vừa ngọt, Rất ngọt
- Nhiệt độ: Nóng sốt, Ấm vừa, Lạnh vừa, Đá xay
- Đặc biệt: Không sữa, Thêm shot, Ít béo, v.v.

## Data Flow

### 1. Khi thêm note:
```javascript
// User taps "Thêm ghi chú" cho Cà phê sữa (itemId: "abc123")
openNoteModal("abc123", null)

// User nhập "Ít đường, nhiều đá"
currentNote.value = "Ít đường, nhiều đá"

// User taps "Lưu"
saveNote()
// → itemNotes.value["abc123"] = "Ít đường, nhiều đá"
```

### 2. Khi confirm order:
```javascript
handleConfirm()
// → Tạo cartArray với note:
[
  {
    menu_item_id: "abc123",
    name: "Cà phê sữa",
    price: 25000,
    quantity: 2,
    note: "Ít đường, nhiều đá" // ← Note được include
  }
]
```

### 3. Backend nhận data:
```json
{
  "customer_name": "Khách A",
  "items": [
    {
      "menu_item_id": "abc123",
      "name": "Cà phê sữa",
      "price": 25000,
      "quantity": 2,
      "note": "Ít đường, nhiều đá"
    }
  ]
}
```

## Hiển thị Note

### Trong CreateOrderModal
- Note hiển thị dưới nút ghi chú
- Background vàng nhạt
- Border vàng
- Icon 💬

### Trong Order Detail (OrderView)
- Note hiển thị trong danh sách items
- Có thể thêm styling đặc biệt để highlight

### Trong Bill (Print)
- Note được in ra cùng với tên món
- Giúp barista biết yêu cầu đặc biệt

## Testing Checklist

- [ ] Thêm món vào cart → Nút ghi chú hiện ra
- [ ] Tap nút "Thêm ghi chú" → Modal mở
- [ ] Tap gợi ý nhanh → Text được thêm vào
- [ ] Tap nhiều gợi ý → Text nối bằng dấu phẩy
- [ ] Nhập text tự do → Text hiển thị đúng
- [ ] Tap "Lưu" → Note được lưu
- [ ] Note hiển thị dưới nút ghi chú
- [ ] Tap "Sửa ghi chú" → Modal mở với note cũ
- [ ] Xóa hết text và lưu → Note bị xóa
- [ ] Xóa món khỏi cart → Note cũng bị xóa
- [ ] Confirm order → Note được gửi đến backend
- [ ] Note hiển thị trong order detail
- [ ] Note được in ra trong bill

## Variants Support

Note hoạt động độc lập cho mỗi variant:

```javascript
// Cà phê sữa - Size M
itemNotes.value["abc123_variant-m"] = "Ít đường"

// Cà phê sữa - Size L
itemNotes.value["abc123_variant-l"] = "Nhiều đá"
```

Mỗi size có thể có note riêng!

## Lưu ý

1. **Note không bắt buộc**: Món có thể không có note
2. **Note được lưu theo item**: Mỗi món/variant có note riêng
3. **Note bị xóa khi xóa món**: Khi qty = 0, note cũng bị xóa
4. **Note được clear sau confirm**: Sau khi tạo order, notes được reset
5. **Backend đã hỗ trợ**: Không cần thay đổi backend

## Mở rộng trong tương lai

1. **Lưu gợi ý thường dùng**: Học từ lịch sử note của user
2. **Note templates**: Tạo template note cho từng loại món
3. **Voice input**: Nhập note bằng giọng nói
4. **Note cho toàn order**: Thêm note chung cho cả order (đã có field `note` trong Order)

## Files thay đổi

1. `frontend/src/components/CreateOrderModal.vue` - Thêm UI và logic note

## API

Không cần thay đổi API. Backend đã hỗ trợ field `note` trong `OrderItem`.

---

**Ngày thực hiện:** 4 tháng 3, 2026
**Trạng thái:** ✅ Hoàn thành
