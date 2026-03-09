# Phân Tích Chức Năng Gộp Bill (Merge Bills)

## 📋 Tổng Quan

Chức năng gộp bill cho phép nhân viên phục vụ (waiter) hoặc thu ngân (cashier) gộp nhiều bill chưa thanh toán thành một bill duy nhất. Điều này hữu ích khi:
- Khách ngồi nhiều bàn muốn gộp chung thanh toán
- Khách đặt nhiều lần trong cùng một lần đến quán
- Cần hợp nhất các order để dễ quản lý

## 🎯 Yêu Cầu Nghiệp Vụ

### Điều Kiện Gộp Bill
1. **Trạng thái order**: Chỉ gộp được các order có status = `CREATED` hoặc `PAID`
2. **Chưa gửi quầy bar**: Không gộp được order đã ở trạng thái `QUEUED`, `IN_PROGRESS`, `READY`, `SERVED`
3. **Cùng ca làm việc**: Các order phải thuộc cùng một shift_id
4. **Tối thiểu 2 orders**: Phải chọn ít nhất 2 orders để gộp

### Quy Tắc Gộp
1. **Items**: Gộp tất cả items từ các orders được chọn
2. **Tổng tiền**: Tính lại subtotal và total từ tất cả items
3. **Thanh toán**: 
   - Nếu có order đã PAID → order mới có status = PAID, amount_paid = tổng amount_paid của các orders
   - Nếu tất cả CREATED → order mới có status = CREATED, amount_paid = 0
4. **Thông tin khách**: Lấy customer_name từ order đầu tiên (hoặc cho phép nhập lại)
5. **Waiter**: Lấy waiter từ order đầu tiên (hoặc người thực hiện merge)
6. **Order cũ**: Đánh dấu các order cũ là CANCELLED với lý do "Đã gộp vào order #XXX"

## 🏗️ Thiết Kế Kỹ Thuật

### Backend Changes

#### 1. Domain Model Updates (`backend/domain/order/order.go`)

Thêm struct mới:
```go
type MergeOrdersRequest struct {
    OrderIDs      []string `json:"order_ids" binding:"required,min=2"`
    CustomerName  string   `json:"customer_name"`
    Note          string   `json:"note"`
}

type MergeOrdersResponse struct {
    MergedOrder    *Order   `json:"merged_order"`
    CancelledOrders []string `json:"cancelled_orders"`
    Message        string   `json:"message"`
}
```

#### 2. Service Layer (`backend/application/services/order_service.go`)

Thêm method mới:
```go
func (s *OrderService) MergeOrders(ctx context.Context, req *MergeOrdersRequest, userID, userName string) (*MergeOrdersResponse, error)
```

**Logic xử lý**:
1. Validate tất cả orders tồn tại
2. Kiểm tra điều kiện gộp:
   - Tất cả orders phải có status = CREATED hoặc PAID
   - Tất cả orders phải cùng shift_id
   - Không có order nào đã QUEUED trở đi
3. Tạo order mới:
   - Gộp tất cả items
   - Tính tổng amount_paid
   - Xác định status (PAID nếu có bất kỳ order nào PAID, ngược lại CREATED)
   - Tạo order_number mới
4. Hủy các orders cũ với lý do merge
5. Trả về merged order và danh sách cancelled orders

#### 3. Repository Layer (`backend/infrastructure/mongodb/order_repository.go`)

Thêm method:
```go
func (r *OrderRepository) FindByIDs(ctx context.Context, ids []primitive.ObjectID) ([]*order.Order, error)
```

#### 4. HTTP Handler (`backend/interfaces/http/order_handler.go`)

Thêm endpoint mới:
```go
// POST /waiter/orders/merge
// POST /cashier/orders/merge
func (h *OrderHandler) MergeOrders(c *gin.Context)
```

### Frontend Changes

#### 1. Service Layer (`frontend/src/services/order.js`)

Thêm method:
```javascript
async mergeOrders(orderIds, customerName, note) {
  const response = await api.post('/waiter/orders/merge', {
    order_ids: orderIds,
    customer_name: customerName,
    note: note
  })
  return response.data
}
```

#### 2. Store Updates (`frontend/src/stores/order.js`)

Thêm action:
```javascript
async mergeOrders(orderIds, customerName, note) {
  this.error = null
  try {
    const response = await orderService.mergeOrders(orderIds, customerName, note)
    // Remove cancelled orders from list
    this.orders = this.orders.filter(o => !orderIds.includes(o.id))
    // Add merged order to list
    this.orders.unshift(response.merged_order)
    return response
  } catch (error) {
    this.error = error.response?.data?.error || 'Lỗi gộp bill'
    throw error
  }
}
```

#### 3. UI Component (`frontend/src/components/MergeBillsModal.vue`)

Component mới với các tính năng:
- Hiển thị danh sách orders có thể gộp (CREATED hoặc PAID, cùng shift)
- Checkbox để chọn nhiều orders
- Hiển thị preview: tổng items, tổng tiền, số tiền đã thu
- Input để nhập/sửa customer_name
- Nút "Gộp Bill" để thực hiện merge
- Xác nhận trước khi gộp

#### 4. OrderView Updates (`frontend/src/views/OrderView.vue`)

Thêm:
- Nút "Gộp Bill" ở header (chỉ hiện khi có ít nhất 2 unpaid orders)
- Chế độ chọn nhiều orders (selection mode)
- Hiển thị số lượng orders đã chọn
- Mở MergeBillsModal khi click "Gộp Bill"

## 📱 User Flow

### Waiter/Cashier Workflow

1. **Vào màn hình Orders**
   - Xem danh sách orders chưa thanh toán
   - Nhận thấy có nhiều orders của cùng một khách

2. **Bật chế độ gộp bill**
   - Click nút "Gộp Bill" ở header
   - Giao diện chuyển sang selection mode

3. **Chọn orders cần gộp**
   - Checkbox xuất hiện trên mỗi order card
   - Chọn 2 hoặc nhiều orders (chỉ hiện orders CREATED/PAID cùng shift)
   - Hiển thị preview tổng tiền ở bottom bar

4. **Xác nhận gộp**
   - Click "Tiếp tục" → Mở modal xác nhận
   - Modal hiển thị:
     - Danh sách orders sẽ gộp (order number, items, total)
     - Tổng items sau khi gộp
     - Tổng tiền sau khi gộp
     - Số tiền đã thu (nếu có)
     - Input customer_name (mặc định từ order đầu tiên)
   - Click "Xác nhận gộp"

5. **Kết quả**
   - Hiển thị thông báo thành công
   - Orders cũ biến mất khỏi danh sách
   - Order mới xuất hiện với tất cả items đã gộp
   - Có thể tiếp tục thu tiền hoặc chỉnh sửa order mới

## 🔒 Validation Rules

### Backend Validation
```go
// Kiểm tra số lượng orders
if len(orderIDs) < 2 {
    return nil, errors.New("phải chọn ít nhất 2 orders để gộp")
}

// Kiểm tra tất cả orders tồn tại
orders, err := s.orderRepo.FindByIDs(ctx, orderIDs)
if len(orders) != len(orderIDs) {
    return nil, errors.New("một số orders không tồn tại")
}

// Kiểm tra status
for _, o := range orders {
    if o.Status != order.StatusCreated && o.Status != order.StatusPaid {
        return nil, fmt.Errorf("order %s không thể gộp (status: %s)", o.OrderNumber, o.Status)
    }
}

// Kiểm tra cùng shift
shiftID := orders[0].ShiftID
for _, o := range orders {
    if o.ShiftID != shiftID {
        return nil, errors.New("các orders phải thuộc cùng một ca làm việc")
    }
}

// Kiểm tra shift còn mở
shift, err := s.shiftRepo.FindByID(ctx, shiftID)
if err != nil || shift.Status != order.ShiftOpen {
    return nil, errors.New("ca làm việc đã đóng, không thể gộp orders")
}
```

### Frontend Validation
```javascript
// Chỉ hiển thị orders có thể gộp
const mergeableOrders = computed(() => {
  return orders.value.filter(o => 
    (o.status === 'CREATED' || o.status === 'PAID') &&
    o.shift_id === currentShiftId.value
  )
})

// Disable nút gộp nếu chọn < 2 orders
const canMerge = computed(() => selectedOrders.value.length >= 2)
```

## 📊 Database Impact

### Orders Collection
- Tạo 1 order mới (merged order)
- Update N orders cũ: status = CANCELLED, cancel_reason = "Đã gộp vào order #XXX"

### Shift Revenue
- Không ảnh hưởng vì chỉ gộp orders chưa hoặc đã thanh toán trong cùng shift
- Nếu có orders đã PAID → amount_paid được chuyển sang order mới

### Batch Ingredients
- Không cần xử lý lại vì ingredients đã được deduct khi tạo orders ban đầu
- Order mới không deduct thêm ingredients

## 🎨 UI/UX Design

### Selection Mode
```
┌─────────────────────────────────────┐
│ 📋 Orders        [Hủy] [Gộp Bill]  │
├─────────────────────────────────────┤
│ ☑️ #001 - Nguyễn Văn A              │
│    2 items • 50,000đ                │
├─────────────────────────────────────┤
│ ☑️ #002 - Nguyễn Văn A              │
│    1 item • 25,000đ                 │
├─────────────────────────────────────┤
│ ☐ #003 - Trần Thị B                 │
│    3 items • 75,000đ                │
└─────────────────────────────────────┘
┌─────────────────────────────────────┐
│ 2 orders đã chọn • Tổng: 75,000đ   │
│              [Tiếp tục gộp] ────────┤
└─────────────────────────────────────┘
```

### Confirmation Modal
```
┌─────────────────────────────────────┐
│         Xác nhận gộp bill           │
├─────────────────────────────────────┤
│ Gộp các orders sau:                 │
│                                     │
│ • #001: 2 items - 50,000đ          │
│ • #002: 1 item - 25,000đ           │
│                                     │
│ ─────────────────────────────────   │
│ Tổng: 3 items - 75,000đ            │
│ Đã thu: 0đ                          │
│ Còn lại: 75,000đ                    │
│                                     │
│ Tên khách: [Nguyễn Văn A        ]  │
│ Ghi chú:   [                    ]  │
│                                     │
│     [Hủy]        [Xác nhận gộp]    │
└─────────────────────────────────────┘
```

## 🧪 Test Cases

### Backend Tests
1. ✅ Gộp 2 orders CREATED thành công
2. ✅ Gộp 2 orders PAID thành công
3. ✅ Gộp 1 CREATED + 1 PAID thành công (kết quả: PAID)
4. ❌ Gộp < 2 orders → Error
5. ❌ Gộp orders khác shift → Error
6. ❌ Gộp order đã QUEUED → Error
7. ❌ Gộp order không tồn tại → Error
8. ✅ Items được gộp đúng
9. ✅ Tổng tiền tính đúng
10. ✅ Amount_paid tính đúng
11. ✅ Orders cũ bị cancel với lý do đúng

### Frontend Tests
1. ✅ Hiển thị nút "Gộp Bill" khi có >= 2 unpaid orders
2. ✅ Ẩn nút khi không đủ điều kiện
3. ✅ Selection mode hoạt động đúng
4. ✅ Chỉ hiển thị orders có thể gộp
5. ✅ Preview tổng tiền đúng
6. ✅ Modal confirmation hiển thị đúng thông tin
7. ✅ Sau merge, orders cũ biến mất, order mới xuất hiện
8. ✅ Error handling hiển thị đúng

## 🚀 Implementation Plan

### Phase 1: Backend (2-3 giờ)
1. ✅ Thêm MergeOrdersRequest/Response structs
2. ✅ Implement OrderService.MergeOrders()
3. ✅ Thêm OrderRepository.FindByIDs()
4. ✅ Thêm HTTP handler và routes
5. ✅ Viết unit tests

### Phase 2: Frontend (3-4 giờ)
1. ✅ Thêm orderService.mergeOrders()
2. ✅ Thêm store action
3. ✅ Tạo MergeBillsModal component
4. ✅ Update OrderView với selection mode
5. ✅ Styling và responsive

### Phase 3: Testing & Polish (1-2 giờ)
1. ✅ Integration testing
2. ✅ UI/UX refinement
3. ✅ Error handling
4. ✅ Documentation

**Tổng thời gian ước tính: 6-9 giờ**

## 🔄 Alternative Approaches

### Approach 1: Soft Merge (Không hủy orders cũ)
- Tạo order mới nhưng giữ orders cũ với flag `merged_into`
- Ưu điểm: Có thể trace lại lịch sử
- Nhược điểm: Phức tạp hơn, dễ nhầm lẫn

### Approach 2: In-place Merge (Giữ 1 order, xóa các order còn lại)
- Chọn 1 order làm "master", gộp items vào đó
- Ưu điểm: Giữ nguyên order_number
- Nhược điểm: Khó xác định order nào là master

### Approach 3: Current Design (Tạo mới + Cancel cũ) ✅ RECOMMENDED
- Tạo order hoàn toàn mới, cancel tất cả orders cũ
- Ưu điểm: Rõ ràng, dễ hiểu, dễ implement
- Nhược điểm: Mất order_number gốc (nhưng có thể trace qua cancel_reason)

## 📝 Notes

1. **Permissions**: Chỉ waiter và cashier có quyền merge orders
2. **Audit Trail**: Cancel reason ghi rõ "Đã gộp vào order #XXX" để trace
3. **Batch Ingredients**: Không cần xử lý vì đã deduct khi tạo orders ban đầu
4. **Print**: Order mới chưa in bill, có thể in lại sau khi merge
5. **Shift Revenue**: Không ảnh hưởng vì chỉ gộp trong cùng shift

## ❓ Questions to Clarify

1. Có cho phép gộp orders của nhiều waiters khác nhau không?
   - **Đề xuất**: Có, vì khách có thể được phục vụ bởi nhiều người
   
2. Có cho phép gộp orders có discount khác nhau không?
   - **Đề xuất**: Có, tính tổng discount = sum(discounts)
   
3. Có giới hạn số lượng orders tối đa có thể gộp không?
   - **Đề xuất**: Không giới hạn, nhưng UI nên cảnh báo nếu > 10 orders
   
4. Có cho phép un-merge (tách lại) không?
   - **Đề xuất**: Không, quá phức tạp. Nếu cần, tạo orders mới

## 🎯 Success Criteria

1. ✅ Waiter có thể gộp nhiều bills chưa thanh toán thành 1 bill
2. ✅ Tất cả items được gộp đúng, không mất dữ liệu
3. ✅ Tổng tiền tính chính xác
4. ✅ Orders cũ được đánh dấu cancelled với lý do rõ ràng
5. ✅ UI/UX trực quan, dễ sử dụng
6. ✅ Không ảnh hưởng đến shift revenue tracking
7. ✅ Không ảnh hưởng đến batch ingredient tracking
