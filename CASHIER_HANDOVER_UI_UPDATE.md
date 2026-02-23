# Cashier Handover View - UI Updates

## Thay Đổi

Frontend view của cashier giờ đây phân biệt rõ ràng giữa tiền mặt và tiền chuyển khoản.

## Các Cải Tiến

### 1. Badge Loại Tiền

Mỗi handover giờ có badge màu sắc để phân biệt:

- 💵 **Tiền mặt** - Badge xanh lá (green)
- 💳 **Chuyển khoản** - Badge xanh dương (blue)  
- 💰 **Cả hai** - Badge tím (purple)

```vue
<!-- Payment type badge -->
<span v-if="handover.cash_declared_amount > 0 && handover.transfer_declared_amount > 0"
  class="bg-purple-100 text-purple-800">
  💰 Cả hai
</span>
<span v-else-if="handover.transfer_declared_amount > 0"
  class="bg-blue-100 text-blue-800">
  💳 Chuyển khoản
</span>
<span v-else-if="handover.cash_declared_amount > 0"
  class="bg-green-100 text-green-800">
  💵 Tiền mặt
</span>
```

### 2. Hiển Thị Số Tiền Riêng Biệt

#### Trong Danh Sách Pending

```
┌─────────────────────────────────────┐
│ waiter                              │
│ 23/02/2026 11:29                    │
│ [Một phần] [💳 Chuyển khoản]        │
│                                     │
│                     💳 30.000 ₫     │
└─────────────────────────────────────┘
```

Thay vì chỉ hiển thị tổng:
```
┌─────────────────────────────────────┐
│ waiter                              │
│ 23/02/2026 11:29                    │
│ [Một phần]                          │
│                                     │
│                     30.000 ₫        │
└─────────────────────────────────────┘
```

#### Khi Có Cả Hai Loại

```
┌─────────────────────────────────────┐
│ waiter                              │
│ 23/02/2026 11:29                    │
│ [Một phần] [💰 Cả hai]              │
│                                     │
│                     💵 50.000 ₫     │
│                     💳 30.000 ₫     │
│                     Tổng: 80.000 ₫  │
└─────────────────────────────────────┘
```

### 3. Modal Xác Nhận

Khi cashier xác nhận handover, modal hiển thị rõ:

```
┌─────────────────────────────────────┐
│ ✅ Xác nhận bàn giao                │
├─────────────────────────────────────┤
│ Waiter: waiter                      │
│                                     │
│ 💳 Tiền CK khai báo:    30.000 ₫   │
│                                     │
│ Loại: Một phần                      │
├─────────────────────────────────────┤
│ Số tiền thực nhận (VNĐ) *          │
│ [30000                    ]         │
│                                     │
│ Ghi chú của bạn                     │
│ [                         ]         │
│                                     │
│ [Hủy]              [Xác nhận]       │
└─────────────────────────────────────┘
```

Hoặc khi có cả hai:

```
┌─────────────────────────────────────┐
│ ✅ Xác nhận bàn giao                │
├─────────────────────────────────────┤
│ Waiter: waiter                      │
│                                     │
│ 💵 Tiền mặt khai báo:   50.000 ₫   │
│ 💳 Tiền CK khai báo:    30.000 ₫   │
│ ─────────────────────────────────   │
│ Tổng:                   80.000 ₫   │
│                                     │
│ Loại: Một phần                      │
└─────────────────────────────────────┘
```

### 4. Today's Handovers

Danh sách handovers hôm nay cũng hiển thị tương tự:

```
📋 Bàn giao hôm nay

┌─────────────────────────────────────┐
│ waiter                              │
│ 11:29                               │
│ [Một phần] [💳 CK]                  │
│                                     │
│                     💳 30.000 ₫     │
│                     [Đã xác nhận]   │
└─────────────────────────────────────┘

┌─────────────────────────────────────┐
│ waiter2                             │
│ 10:15                               │
│ [Một phần] [💵 Mặt]                 │
│                                     │
│                     💵 50.000 ₫     │
│                     [Đã xác nhận]   │
└─────────────────────────────────────┘
```

## Màu Sắc

### Badge Colors

- **Tiền mặt**: `bg-green-100 text-green-800`
- **Chuyển khoản**: `bg-blue-100 text-blue-800`
- **Cả hai**: `bg-purple-100 text-purple-800`

### Amount Colors

- **Tiền mặt**: `text-green-600`
- **Chuyển khoản**: `text-blue-600`
- **Tổng**: `text-gray-900`

## Logic Hiển Thị

### Kiểm Tra Format

```javascript
// Check if using new format (separate amounts)
if (handover.cash_declared_amount > 0 || handover.transfer_declared_amount > 0) {
  // Display separate amounts
} else {
  // Fallback to old format (declared_amount)
}
```

### Badge Logic

```javascript
// Determine payment type
if (cash > 0 && transfer > 0) {
  // Show "Cả hai" badge
} else if (transfer > 0) {
  // Show "Chuyển khoản" badge
} else if (cash > 0) {
  // Show "Tiền mặt" badge
}
```

## Backward Compatibility

✅ Handovers cũ (chỉ có `declared_amount`) vẫn hiển thị đúng
✅ Không có badge loại tiền cho handovers cũ (để tránh nhầm lẫn)
✅ Hiển thị số tiền như cũ cho old format

## Files Changed

- `frontend/src/views/CashierHandoverView.vue`
  - Thêm badge loại tiền
  - Hiển thị separate amounts
  - Update modal confirm
  - Update today's handovers list

## Screenshots (Mô Tả)

### Before
```
waiter
23/02/2026 11:29
[Một phần]
                    30.000 ₫
```

### After
```
waiter
23/02/2026 11:29
[Một phần] [💳 Chuyển khoản]
                    💳 30.000 ₫
```

## Benefits

1. ✅ **Rõ ràng hơn**: Cashier biết ngay đây là tiền mặt hay chuyển khoản
2. ✅ **Tránh nhầm lẫn**: Không còn tình trạng tưởng là tiền mặt nhưng lại là CK
3. ✅ **Dễ đối soát**: Phân biệt rõ từng loại tiền khi xác nhận
4. ✅ **Trực quan**: Màu sắc và icon giúp nhận diện nhanh
5. ✅ **Backward compatible**: Vẫn hoạt động với handovers cũ

## Testing

### Test Case 1: Handover chỉ tiền mặt
- Badge: 💵 Tiền mặt (green)
- Amount: 💵 50.000 ₫ (green)

### Test Case 2: Handover chỉ tiền CK
- Badge: 💳 Chuyển khoản (blue)
- Amount: 💳 30.000 ₫ (blue)

### Test Case 3: Handover cả hai
- Badge: 💰 Cả hai (purple)
- Amounts: 
  - 💵 50.000 ₫ (green)
  - 💳 30.000 ₫ (blue)
  - Tổng: 80.000 ₫ (gray)

### Test Case 4: Old handover
- No payment type badge
- Single amount display
- Works as before

## Next Steps

1. ✅ UI phân biệt loại tiền (DONE)
2. ⏳ Backend auto-convert (DONE)
3. ⏳ Test với real data
4. ⏳ User feedback
5. ⏳ Reports phân biệt cash/transfer handovers
