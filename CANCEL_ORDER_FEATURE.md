# Tính năng Hủy Order cho Waiter

## Tổng quan

Đã thêm tính năng cho phép waiter hủy order khi chưa thanh toán (status = CREATED).

## Thay đổi

### 1. Backend (`backend/main.go`)

**Thêm endpoint cancel order cho waiter:**
```go
waiter.POST("/orders/:id/cancel", orderHandler.CancelOrder)
```

**Quyền truy cập:**
- Waiter ✅
- Cashier ✅
- Manager ✅

**Điều kiện hủy order:**
- Order phải ở trạng thái: CREATED, PAID, QUEUED, hoặc IN_PROGRESS
- Không thể hủy order đã SERVED hoặc LOCKED
- Phải cung cấp lý do hủy (bắt buộc)

### 2. Frontend Service (`frontend/src/services/order.js`)

**Cập nhật endpoint:**
```javascript
async cancelOrder(id, reason) {
  const response = await api.post(`/waiter/orders/${id}/cancel`, { reason })
  return response.data
}
```

Thay đổi từ `/cashier/orders/...` sang `/waiter/orders/...` để waiter có thể sử dụng.

### 3. Frontend UI (`frontend/src/views/OrderView.vue`)

**Thêm nút "Hủy order":**
- Hiển thị trong order detail modal
- Chỉ hiện khi order ở trạng thái CREATED (chưa thanh toán)
- Màu đỏ để cảnh báo hành động quan trọng

**Thêm Cancel Order Modal:**
- Hiển thị thông tin order cần hủy
- Yêu cầu nhập lý do hủy (bắt buộc)
- Cảnh báo hành động không thể hoàn tác
- Nút xác nhận chỉ active khi đã nhập lý do

**State mới:**
```javascript
const showCancelModal = ref(false)
const cancelReason = ref('')
const cancelingOrder = ref(null)
```

**Methods mới:**
- `confirmCancelOrder(order)` - Mở modal xác nhận
- `closeCancelModal()` - Đóng modal
- `processCancelOrder()` - Xử lý hủy order

## Luồng sử dụng

1. **Waiter mở order detail** (tap vào order trong danh sách)
2. **Thấy nút "❌ Hủy order"** (chỉ khi order chưa thanh toán)
3. **Tap nút "Hủy order"**
4. **Modal xác nhận hiện ra** với:
   - Thông tin order (số order, khách hàng, tổng tiền)
   - Cảnh báo hành động không thể hoàn tác
   - Ô nhập lý do hủy (bắt buộc)
5. **Nhập lý do hủy** (ví dụ: "Khách đổi ý", "Order nhầm", etc.)
6. **Tap "Xác nhận hủy"**
7. **Order chuyển sang trạng thái CANCELLED**
8. **Thông báo thành công**

## Validation

### Backend
- Kiểm tra order tồn tại
- Kiểm tra quyền truy cập
- Kiểm tra trạng thái order (state machine)
- Kiểm tra lý do hủy không rỗng

### Frontend
- Nút "Xác nhận hủy" disabled khi chưa nhập lý do
- Hiển thị cảnh báo rõ ràng
- Xác nhận trước khi thực hiện

## Trạng thái Order có thể hủy

✅ **CREATED** - Chưa thanh toán (waiter có thể hủy)
✅ **PAID** - Đã thanh toán nhưng chưa gửi bar (waiter có thể hủy)
✅ **QUEUED** - Đang chờ pha (waiter có thể hủy)
✅ **IN_PROGRESS** - Đang pha chế (waiter có thể hủy)
❌ **READY** - Đã pha xong (không thể hủy)
❌ **SERVED** - Đã giao khách (không thể hủy)
❌ **LOCKED** - Đã khóa ca (không thể hủy)

## UI/UX

### Màu sắc
- Nút hủy: Đỏ (`bg-red-500`)
- Modal header: Đỏ (`text-red-600`)
- Cảnh báo: Nền đỏ nhạt (`bg-red-50`)

### Icon
- Nút: ❌
- Modal: ❌

### Text
- Nút: "❌ Hủy order"
- Modal title: "❌ Hủy Order"
- Cảnh báo: "Cảnh báo: Bạn đang hủy order này. Hành động này không thể hoàn tác."

### Vị trí
- Nút "Hủy order" nằm sau nút "Chỉnh sửa"
- Chỉ hiện khi order ở trạng thái CREATED

## Testing Checklist

- [ ] Waiter có thể thấy nút "Hủy order" khi order CREATED
- [ ] Tap nút "Hủy order" mở modal xác nhận
- [ ] Modal hiển thị đúng thông tin order
- [ ] Không thể xác nhận khi chưa nhập lý do
- [ ] Có thể nhập lý do hủy
- [ ] Tap "Quay lại" đóng modal
- [ ] Tap "Xác nhận hủy" với lý do hợp lệ → hủy thành công
- [ ] Order chuyển sang trạng thái CANCELLED
- [ ] Order biến mất khỏi danh sách CREATED
- [ ] Hiển thị thông báo thành công
- [ ] Không thể hủy order đã thanh toán và gửi bar
- [ ] Không thể hủy order đã served
- [ ] Backend validate đúng quyền truy cập

## Lưu ý

1. **Không thể hoàn tác**: Khi đã hủy order, không thể khôi phục lại
2. **Lý do bắt buộc**: Phải nhập lý do để tracking và báo cáo
3. **Quyền hạn**: Waiter chỉ có thể hủy order chưa thanh toán
4. **State machine**: Backend sử dụng state machine để validate transition

## Files thay đổi

1. `backend/main.go` - Thêm endpoint cancel cho waiter
2. `frontend/src/services/order.js` - Cập nhật endpoint
3. `frontend/src/views/OrderView.vue` - Thêm UI và logic

## API Endpoint

```
POST /waiter/orders/:id/cancel
```

**Request Body:**
```json
{
  "reason": "Khách đổi ý"
}
```

**Response:**
```json
{
  "id": "...",
  "order_number": "ORD-001",
  "status": "CANCELLED",
  "cancelled_reason": "Khách đổi ý",
  ...
}
```

**Error Response:**
```json
{
  "error": "cannot cancel order in SERVED status",
  "status": "SERVED",
  "can_cancel": false
}
```

---

**Ngày thực hiện:** 4 tháng 3, 2026
**Trạng thái:** ✅ Hoàn thành
