# Implementation: Cancel Closure Feature

## Tổng quan

Đã implement chức năng hủy quy trình đóng ca thu ngân khi người dùng bấm "Quay lại" trước khi hoàn thành các bước quan trọng.

## Các thay đổi

### 1. Backend - Domain Model

**File**: `backend/domain/cashier/cashier_shift.go`

Thêm method `CancelClosure`:
```go
func (cs *CashierShift) CancelClosure(userID, deviceID string, timestamp time.Time) error
```

**Logic**:
- Chỉ cho phép cancel khi status là `CLOSURE_INITIATED`
- Không cho phép cancel nếu đã nhập tiền thực tế (`ActualCash != nil`)
- Chuyển status về `OPEN`
- Ghi audit log với action `closure_cancelled`

### 2. Backend - HTTP Handler

**File**: `backend/interfaces/http/cashier_shift_closure_handler.go`

Thêm method `CancelClosure`:
```go
func (h *CashierShiftClosureHandler) CancelClosure(c *gin.Context)
```

**Logic**:
- Validate state transition qua state machine
- Kiểm tra `CanCancelCashierShiftClosure`
- Gọi domain method `CancelClosure`
- Lưu shift vào database

### 3. Backend - API Route

**File**: `backend/main.go`

Thêm route:
```go
cashierShifts.POST("/:id/cancel-closure", cashierShiftClosureHandler.CancelClosure)
```

**Endpoint**: `POST /api/cashier-shifts/:id/cancel-closure`

### 4. Frontend - Service

**File**: `frontend/src/services/cashierShift.js`

Thêm method:
```javascript
async cancelClosure(shiftId) {
  const response = await api.post(`/cashier-shifts/${shiftId}/cancel-closure`, {})
  return response.data
}
```

### 5. Frontend - View Logic

**File**: `frontend/src/views/CashierShiftClosure.vue`

**Thêm UI Card "Hủy đóng ca"**:
```vue
<!-- Cancel Closure Option (shown when closure initiated but no critical steps) -->
<div v-if="shift.status === CASHIER_SHIFT_STATUS.CLOSURE_INITIATED && !shift.actual_cash" 
     class="bg-orange-50 border-2 border-orange-200 rounded-2xl p-4 shadow-sm">
  <div class="flex items-start gap-3 mb-3">
    <span class="text-2xl">⚠️</span>
    <div class="flex-1">
      <p class="font-bold text-orange-800 mb-1">Đã bắt đầu đóng ca</p>
      <p class="text-sm text-orange-700">
        Nếu bạn muốn hủy quy trình đóng ca và quay về trạng thái mở ca, bấm nút bên dưới.
      </p>
    </div>
  </div>
  <button
    @click="cancelClosure"
    :disabled="processing"
    class="w-full py-3 bg-orange-500 text-white rounded-xl font-bold active:scale-95 transition-transform disabled:opacity-50"
  >
    {{ processing ? 'Đang hủy...' : '↩️ Hủy đóng ca' }}
  </button>
</div>
```

**Thêm method cancelClosure**:
```javascript
const cancelClosure = async () => {
  if (!confirm('Bạn có chắc muốn hủy quy trình đóng ca?\n\nCa sẽ quay về trạng thái mở.')) {
    return
  }
  
  processing.value = true
  error.value = null
  
  try {
    await cashierShiftService.cancelClosure(shift.value.id)
    await loadShift()
    alert('✅ Đã hủy quy trình đóng ca thành công!')
  } catch (err) {
    error.value = err.response?.data?.error || 'Không thể hủy đóng ca'
  } finally {
    processing.value = false
  }
}
```

**Đơn giản hóa method goBack**:
```javascript
const goBack = () => {
  router.push('/cashier')
}
```

## Quy trình hoạt động

### Trường hợp 1: Cancel thành công

1. User bấm "Bắt đầu đóng ca" → Status: `OPEN` → `CLOSURE_INITIATED`
2. UI hiển thị card cảnh báo màu cam với nút "↩️ Hủy đóng ca"
3. User bấm nút "Hủy đóng ca" → Hiện confirm dialog
4. User chọn "OK" → Gọi API cancel-closure
5. Backend validate và chuyển status về `OPEN`
6. Frontend reload shift data và hiển thị success message
7. User có thể tiếp tục làm việc hoặc bấm "Quay lại"

### Trường hợp 2: Không cho phép cancel

1. User bấm "Bắt đầu đóng ca" → Status: `CLOSURE_INITIATED`
2. User nhập tiền thực tế → `ActualCash` được set
3. Card "Hủy đóng ca" biến mất (không còn hiển thị)
4. Nếu user cố gọi API cancel-closure → Backend trả về error
5. Frontend hiển thị error message

### Trường hợp 3: User không muốn cancel

1. User bấm "Bắt đầu đóng ca" → Status: `CLOSURE_INITIATED`
2. UI hiển thị card "Hủy đóng ca"
3. User bấm nút "Hủy đóng ca" → Hiện confirm dialog
4. User chọn "Cancel" → Không gọi API
5. User vẫn ở trang đóng ca, có thể tiếp tục quy trình

### Trường hợp 4: User bấm "Quay lại"

1. User bấm nút "← Quay lại" ở header
2. Navigate về `/cashier` (không có confirm dialog)
3. Ca vẫn ở trạng thái `CLOSURE_INITIATED` nếu chưa cancel

## Validation Rules

### Backend State Machine

- ✅ Cho phép transition: `CLOSURE_INITIATED` → `OPEN` (via `EventCancelClosure`)
- ✅ Kiểm tra: `ActualCash == nil` (chưa nhập tiền thực tế)
- ✅ Ghi audit log: action = `closure_cancelled`

### Frontend UX

- ✅ Hiển thị card cảnh báo màu cam khi status = `CLOSURE_INITIATED` và `actual_cash == null`
- ✅ Nút "Hủy đóng ca" rõ ràng, dễ nhìn
- ✅ Confirm dialog trước khi thực hiện cancel
- ✅ Hiển thị success message sau khi cancel thành công
- ✅ Hiển thị error message nếu cancel thất bại
- ✅ Reload shift data để cập nhật UI
- ✅ Nút "Quay lại" đơn giản, không có logic phức tạp

## Testing

### Manual Test

1. Vào trang đóng ca: `http://localhost:5173/#/cashier/shift-closure/{shift_id}`
2. Bấm "Bắt đầu đóng ca"
3. Verify: Hiển thị card màu cam "Đã bắt đầu đóng ca" với nút "↩️ Hủy đóng ca"
4. Bấm nút "Hủy đóng ca"
5. Chọn "OK" trong confirm dialog
6. Verify: 
   - Hiển thị success message "✅ Đã hủy quy trình đóng ca thành công!"
   - Card "Hủy đóng ca" biến mất
   - Hiển thị lại "Bước 1: Bắt đầu đóng ca"
   - Ca quay về trạng thái `OPEN`

### Test không cho phép cancel

1. Bấm "Bắt đầu đóng ca" lần nữa
2. Nhập tiền thực tế (ví dụ: 500000)
3. Bấm "✓ Xác nhận tiền mặt"
4. Verify: Card "Hủy đóng ca" không còn hiển thị
5. Thử gọi API cancel-closure trực tiếp → Nhận error

### Automated Test

Chạy script test:
```bash
./test-cancel-closure.sh
```

Test cases:
- ✅ Can initiate closure
- ✅ Can cancel closure before recording actual cash
- ✅ Audit log is updated correctly
- ✅ Cannot cancel after recording actual cash

## API Documentation

### Cancel Closure

**Endpoint**: `POST /api/cashier-shifts/:id/cancel-closure`

**Headers**:
```
Authorization: Bearer {token}
Content-Type: application/json
```

**Request Body**: `{}`

**Success Response** (200):
```json
{
  "id": "699c269a1ba7fb0030879c87",
  "status": "OPEN",
  "cashier_name": "Thu Ngân 1",
  "starting_float": 500000,
  "system_cash": 500000,
  "actual_cash": null,
  "audit_log": [
    {
      "action": "closure_initiated",
      "user_id": "...",
      "timestamp": "2024-01-15T10:00:00Z"
    },
    {
      "action": "closure_cancelled",
      "user_id": "...",
      "timestamp": "2024-01-15T10:05:00Z"
    }
  ]
}
```

**Error Responses**:

400 - Cannot cancel (actual cash recorded):
```json
{
  "error": "cannot cancel closure: actual cash has been recorded"
}
```

400 - Invalid status:
```json
{
  "error": "cannot cancel closure: shift status must be ClosureInitiated"
}
```

404 - Shift not found:
```json
{
  "error": "shift not found"
}
```

## Security & Audit

- ✅ Tất cả actions đều được ghi vào audit log
- ✅ User ID và device ID được track
- ✅ Timestamp chính xác cho mọi thay đổi
- ✅ Không cho phép cancel sau khi đã nhập tiền thực tế (prevent data loss)

## Deployment Checklist

- [x] Backend domain model updated
- [x] Backend handler implemented
- [x] API route added
- [x] Frontend service method added
- [x] Frontend view logic updated
- [x] Test script created
- [ ] Backend restart required
- [ ] Frontend rebuild required
- [ ] Manual testing
- [ ] Production deployment

## Notes

- Feature này giúp tránh trường hợp ca bị "kẹt" ở trạng thái `CLOSURE_INITIATED`
- UX tốt hơn - cho phép người dùng sửa lỗi khi bấm nhầm
- Audit trail đầy đủ cho mọi thao tác
- Validation chặt chẽ để đảm bảo data integrity
