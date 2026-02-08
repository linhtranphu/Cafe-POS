# Checklist Kiểm Tra Constants Refactor

## Mục Đích
Kiểm tra xem việc refactor từ hardcoded strings sang constants có hoạt động đúng không.

---

## 1. ShiftView.vue ✅

### Kiểm tra User Role
- [ ] Login với tài khoản **Manager** → Xem được tất cả ca
- [ ] Login với tài khoản **Cashier** → Xem được tất cả ca
- [ ] Login với tài khoản **Waiter** → Chỉ xem được ca của mình
- [ ] Login với tài khoản **Barista** → Chỉ xem được ca của mình

### Kiểm tra Shift Status
- [ ] Mở ca mới → Status hiển thị "Đang mở"
- [ ] Đóng ca → Status chuyển sang "Đã đóng"
- [ ] Filter theo status hoạt động đúng

### Kiểm tra Handover (Waiter)
- [ ] Waiter có thể tạo bàn giao một phần
- [ ] Waiter có thể tạo bàn giao toàn bộ và đóng ca
- [ ] Pending handover hiển thị đúng
- [ ] Hủy handover hoạt động

---

## 2. ManagerShiftView.vue ✅

### Kiểm tra Filter Tabs
- [ ] Click "Tất cả" → Hiển thị tất cả ca
- [ ] Click "Đang mở" → Chỉ hiển thị ca OPEN
- [ ] Click "Đã đóng" → Chỉ hiển thị ca CLOSED

### Kiểm tra Waiter Shifts
- [ ] Hiển thị đúng số lượng ca waiter đang mở
- [ ] Status badge hiển thị đúng màu
- [ ] Click vào ca → Modal chi tiết hiển thị đúng

### Kiểm tra Barista Shifts
- [ ] Hiển thị đúng số lượng ca barista đang mở
- [ ] Status badge hiển thị đúng màu
- [ ] Click vào ca → Modal chi tiết hiển thị đúng

### Kiểm tra Cashier Shifts
- [ ] Hiển thị đúng số lượng ca cashier đang mở
- [ ] Status badge hiển thị đúng (OPEN, CLOSURE_INITIATED, CLOSED)
- [ ] Click vào ca → Modal chi tiết hiển thị đúng

---

## 3. CashierShiftClosure.vue ✅

### Kiểm tra Step 1: Initiate Closure
- [ ] Chỉ hiển thị khi status = OPEN
- [ ] Click "Bắt đầu đóng ca" → Chuyển sang CLOSURE_INITIATED

### Kiểm tra Step 2: Record Actual Cash
- [ ] Chỉ hiển thị khi status = CLOSURE_INITIATED và chưa có actual_cash
- [ ] Nhập tiền thực tế → Lưu thành công

### Kiểm tra Step 3: Document Variance
- [ ] Chỉ hiển thị khi có chênh lệch và chưa có reason
- [ ] Chọn lý do và nhập ghi chú → Lưu thành công

### Kiểm tra Step 4: Confirm Responsibility
- [ ] Chỉ hiển thị khi đã hoàn thành các bước trước
- [ ] Click "Tôi xác nhận" → Lưu thành công

### Kiểm tra Step 5: Close Shift
- [ ] Chỉ hiển thị khi đã confirm responsibility
- [ ] Kiểm tra waiter shifts → Cảnh báo nếu còn ca waiter mở
- [ ] Click "Đóng ca" → Status chuyển sang CLOSED

### Kiểm tra Completed State
- [ ] Hiển thị khi status = CLOSED
- [ ] Hiển thị thông báo hoàn thành
- [ ] Button "Quay lại Dashboard" hoạt động

---

## 4. DashboardView.vue ✅

### Kiểm tra Role-based Display
- [ ] Manager → Hiển thị tất cả thống kê
- [ ] Cashier → Hiển thị thống kê cashier
- [ ] Waiter → Hiển thị thống kê waiter
- [ ] Barista → Hiển thị thống kê barista

### Kiểm tra Order Status
- [ ] Đếm đúng số order theo từng status
- [ ] Filter theo status hoạt động đúng

---

## 5. UserManagementView.vue ✅

### Kiểm tra Role Badge
- [ ] Manager → Badge màu đỏ
- [ ] Cashier → Badge màu vàng
- [ ] Waiter → Badge màu xanh dương
- [ ] Barista → Badge màu tím

### Kiểm tra Role Filter
- [ ] Filter theo role hoạt động đúng
- [ ] Tạo user mới với role → Lưu đúng
- [ ] Cập nhật role → Lưu đúng

---

## 6. IngredientManagementView.vue ✅

### Kiểm tra Adjustment Types
- [ ] Chọn "Nhập kho" → Type = ADD
- [ ] Chọn "Xuất kho" → Type = REMOVE
- [ ] Chọn "Điều chỉnh" → Type = ADJUST

### Kiểm tra Price Logic
- [ ] ADD với giá mới → Tính weighted average
- [ ] REMOVE → Không thay đổi giá
- [ ] ADJUST giảm → Không thay đổi giá
- [ ] ADJUST tăng với giá mới → Tính weighted average

---

## 7. OrderView.vue ✅

### Kiểm tra Status Filter
- [ ] Filter theo từng status hoạt động đúng
- [ ] Badge màu hiển thị đúng cho từng status

### Kiểm tra Payment Method
- [ ] Chọn "Tiền mặt" → CASH
- [ ] Chọn "Chuyển khoản" → TRANSFER
- [ ] Chọn "QR" → QR

---

## 8. BaristaView.vue ✅

### Kiểm tra Order Status
- [ ] Queue tab → Hiển thị orders QUEUED
- [ ] Working tab → Hiển thị orders IN_PROGRESS
- [ ] Ready tab → Hiển thị orders READY

---

## 9. CashierDashboard.vue ✅

### Kiểm tra Payment Method Display
- [ ] Tiền mặt → Badge xanh lá
- [ ] Chuyển khoản → Badge xanh dương
- [ ] QR → Badge tím

### Kiểm tra Shift Status
- [ ] Hiển thị đúng status của waiter shifts
- [ ] Hiển thị đúng status của cashier shifts

---

## 10. ExpenseManagementView.vue ✅

### Kiểm tra Payment Method
- [ ] Chọn "Tiền mặt" → CASH
- [ ] Chọn "Chuyển khoản" → BANK
- [ ] Chọn "Thẻ" → CARD

---

## 11. FacilityManagementView.vue ✅

### Kiểm tra Facility Status
- [ ] Status "Đang sử dụng" → IN_USE
- [ ] Status "Đang sửa chữa" → REPAIRING
- [ ] Status "Hỏng" → BROKEN

---

## Kiểm Tra Chung

### Console Errors
- [ ] Không có lỗi trong console
- [ ] Không có warning về undefined constants

### Backend Communication
- [ ] API calls vẫn hoạt động bình thường
- [ ] Data được gửi đúng format
- [ ] Response được xử lý đúng

### Performance
- [ ] Không có lag khi chuyển trang
- [ ] Filter hoạt động nhanh
- [ ] Load data không chậm hơn trước

---

## Kết Quả

**Tổng số test cases**: ~80+  
**Passed**: ___  
**Failed**: ___  
**Skipped**: ___  

---

## Ghi Chú

- Nếu có lỗi, ghi rõ:
  - View nào
  - Test case nào
  - Lỗi gì
  - Steps để reproduce

- Nếu cần fix:
  - Kiểm tra import constants
  - Kiểm tra tên constant có đúng không
  - Kiểm tra backend có trả đúng giá trị không

---

**Ngày tạo**: 2026-02-07  
**Người kiểm tra**: ___________  
**Ngày kiểm tra**: ___________

