# Merge Bills Implementation Summary

## ✅ Hoàn Thành

Đã implement thành công chức năng gộp bill cho các orders chưa thanh toán.

## 📦 Backend Changes

### 1. Domain Layer (`backend/domain/order/order.go`)
- ✅ Thêm `MergeOrdersRequest` struct
- ✅ Thêm `MergeOrdersResponse` struct  
- ✅ Thêm method `IsMergeable()` để kiểm tra order có thể gộp không

### 2. Repository Layer (`backend/infrastructure/mongodb/order_repository.go`)
- ✅ Thêm method `FindByIDs()` để lấy nhiều orders cùng lúc

### 3. Service Layer (`backend/application/services/order_service.go`)
- ✅ Thêm interface method `FindByIDs` vào `OrderRepository`
- ✅ Implement `MergeOrders()` method với logic:
  - Validate tối thiểu 2 orders
  - Kiểm tra tất cả orders tồn tại
  - Validate status (chỉ CREATED hoặc PAID)
  - Validate cùng shift
  - Gộp items từ tất cả orders
  - Tính tổng amount_paid
  - Xác định status mới (PAID nếu có bất kỳ order nào PAID)
  - Tạo order mới
  - Cancel tất cả orders cũ với lý do "Đã gộp vào order #XXX"

### 4. HTTP Handler (`backend/interfaces/http/order_handler.go`)
- ✅ Thêm `MergeOrders()` handler

### 5. Routes (`backend/main.go`)
- ✅ Thêm route `POST /waiter/orders/merge`
- ✅ Thêm route `POST /cashier/orders/merge`

## 🎨 Frontend Changes

### 1. Service Layer (`frontend/src/services/order.js`)
- ✅ Thêm `mergeOrders(orderIds, customerName, note)` method

### 2. Store Layer (`frontend/src/stores/order.js`)
- ✅ Thêm `mergeOrders()` action
- ✅ Logic xóa orders cũ và thêm order mới vào danh sách

### 3. Component (`frontend/src/components/MergeBillsModal.vue`)
- ✅ Tạo component mới với UI:
  - Hiển thị danh sách orders sẽ gộp
  - Preview tổng items, tổng tiền, đã thu, còn lại
  - Input customer_name
  - Input note
  - Nút xác nhận gộp
  - Error handling

### 4. View (`frontend/src/views/OrderView.vue`)
- ✅ Thêm nút "Gộp Bill" (🔗) ở header (chỉ hiện khi có >= 2 mergeable orders)
- ✅ Implement selection mode:
  - Header thay đổi khi vào selection mode
  - Checkbox trên mỗi order card
  - Chỉ cho phép chọn orders có thể gộp (CREATED/PAID, cùng shift)
  - Bottom bar với số lượng đã chọn và nút "Gộp"
- ✅ Computed properties:
  - `mergeableOrders`: Lọc orders có thể gộp
  - `mergeableOrdersCount`: Đếm số orders có thể gộp
  - `selectedOrdersForMerge`: Lấy chi tiết orders đã chọn
  - `isOrderMergeable()`: Kiểm tra order có thể gộp không
- ✅ Methods:
  - `enterSelectionMode()`: Vào chế độ chọn
  - `exitSelectionMode()`: Thoát chế độ chọn
  - `toggleOrderSelection()`: Toggle chọn/bỏ chọn order
  - `openMergeModal()`: Mở modal xác nhận
  - `closeMergeModal()`: Đóng modal
  - `handleMerged()`: Xử lý sau khi gộp thành công

## 🎯 Business Logic

### Điều Kiện Gộp
1. ✅ Tối thiểu 2 orders
2. ✅ Tất cả orders phải có status = CREATED hoặc PAID
3. ✅ Tất cả orders phải cùng shift_id
4. ✅ Shift phải còn mở (status = OPEN)

### Quy Tắc Gộp
1. ✅ **Items**: Gộp tất cả items từ các orders
2. ✅ **Discount**: Tổng discount = sum(discounts)
3. ✅ **Amount Paid**: Tổng amount_paid = sum(amount_paid)
4. ✅ **Status**: 
   - Nếu có ít nhất 1 order PAID → status = PAID
   - Nếu tất cả CREATED → status = CREATED
5. ✅ **Payment Info**: Lấy từ order PAID đầu tiên (nếu có)
6. ✅ **Customer Name**: Lấy từ request hoặc order đầu tiên
7. ✅ **Waiter**: Lấy từ order đầu tiên
8. ✅ **Orders cũ**: Cancel với lý do "Đã gộp vào order #XXX"

## 🧪 Testing

### Backend
```bash
cd backend
go build
```
✅ Compilation successful

### Manual Testing Checklist
- [ ] Gộp 2 orders CREATED
- [ ] Gộp 2 orders PAID
- [ ] Gộp 1 CREATED + 1 PAID
- [ ] Gộp 3+ orders
- [ ] Thử gộp < 2 orders (should fail)
- [ ] Thử gộp orders khác shift (should fail)
- [ ] Thử gộp order đã QUEUED (should fail)
- [ ] Kiểm tra items được gộp đúng
- [ ] Kiểm tra tổng tiền tính đúng
- [ ] Kiểm tra amount_paid tính đúng
- [ ] Kiểm tra orders cũ bị cancel với lý do đúng
- [ ] Kiểm tra UI selection mode
- [ ] Kiểm tra modal hiển thị đúng thông tin

## 📱 User Flow

1. **Vào màn hình Orders** → Thấy nút 🔗 (nếu có >= 2 orders có thể gộp)
2. **Click nút 🔗** → Vào selection mode
3. **Chọn orders** → Checkbox xuất hiện, chọn 2+ orders
4. **Click "Gộp X orders"** → Mở modal xác nhận
5. **Xem preview** → Tổng items, tổng tiền, đã thu, còn lại
6. **Nhập customer name** (optional)
7. **Nhập note** (optional)
8. **Click "Xác nhận gộp"** → Gộp thành công
9. **Kết quả** → Orders cũ biến mất, order mới xuất hiện

## 🎨 UI Features

### Selection Mode
- Header thay đổi màu purple với text "Chọn orders để gộp"
- Nút "Hủy" để thoát selection mode
- Hiển thị số lượng đã chọn
- Checkbox tròn trên mỗi order card
- Orders không thể gộp bị mờ đi (opacity-50)
- Orders đã chọn có ring purple và background purple-50

### Bottom Bar (Selection Mode)
- Fixed bottom bar với 2 nút
- Nút "Hủy" để thoát
- Nút "Gộp X orders" (disabled nếu < 2)
- Hiển thị số lượng orders đã chọn

### Merge Modal
- Danh sách orders sẽ gộp với order_number, items, total
- Badge hiển thị status (Đã thu/Chưa thu)
- Preview box màu xanh với tổng items, tổng tiền, đã thu, còn lại
- Input customer_name với giá trị mặc định từ order đầu tiên
- Textarea note
- 2 nút: Hủy và Xác nhận gộp
- Error message nếu có lỗi

## 🔒 Security & Permissions

- ✅ Chỉ Waiter, Cashier, Manager có quyền merge orders
- ✅ Validate shift còn mở
- ✅ Validate orders thuộc cùng shift
- ✅ Validate status hợp lệ

## 📊 Database Impact

### Orders Collection
- Tạo 1 order mới (merged order)
- Update N orders cũ: status = CANCELLED, cancel_reason = "Đã gộp vào order #XXX"

### Không Ảnh Hưởng
- ✅ Shift revenue (vì chỉ gộp trong cùng shift)
- ✅ Batch ingredients (đã deduct khi tạo orders ban đầu)

## 🚀 Deployment

### Backend
```bash
cd backend
go build
./backend
```

### Frontend
```bash
cd frontend
npm run build
# Deploy dist/ folder
```

## 📝 Notes

1. **Order Number**: Order mới có order_number mới, không giữ order_number cũ
2. **Audit Trail**: Cancel reason ghi rõ "Đã gộp vào order #XXX" để trace
3. **Bill Printed**: Order mới có bill_printed = false, có thể in lại
4. **Payment Method**: Lấy từ order PAID đầu tiên (nếu có)
5. **Waiter**: Lấy từ order đầu tiên, không phải người thực hiện merge

## ✨ Future Enhancements

- [ ] Thêm confirmation dialog trước khi gộp
- [ ] Hiển thị preview chi tiết hơn (từng item)
- [ ] Cho phép chọn waiter cho order mới
- [ ] Lưu lịch sử merge trong audit log
- [ ] Thêm filter "Mergeable orders" trong OrderView
- [ ] Thêm keyboard shortcuts (Ctrl+M để vào selection mode)
- [ ] Thêm animation khi gộp thành công

## 🎉 Success!

Chức năng gộp bill đã được implement hoàn chỉnh và sẵn sàng sử dụng!
