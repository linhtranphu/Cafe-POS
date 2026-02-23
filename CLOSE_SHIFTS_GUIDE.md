# Hướng dẫn đóng ca làm việc

## Scripts có sẵn

### 1. Xem danh sách ca đang mở

```bash
./list-open-shifts.sh
```

Script này sẽ hiển thị:
- Danh sách tất cả ca đang mở
- Thông tin chi tiết về mỗi ca:
  - User name và role
  - Thời gian bắt đầu
  - Tiền mặt (start, current, remaining, handed over)
  - Tiền chuyển khoản (revenue, remaining, handed over)
  - Số lượng orders và doanh thu

**Khi nào dùng:**
- Trước khi đóng ca để xem có bao nhiêu ca đang mở
- Kiểm tra trạng thái tiền mặt và chuyển khoản
- Debug vấn đề về ca làm việc

### 2. Đóng tất cả ca đang mở

```bash
./close-all-waiter-shifts.sh
```

Script này sẽ:
1. Yêu cầu xác nhận (phải gõ "yes")
2. Tìm tất cả ca đang mở
3. Đóng từng ca một:
   - Tính tổng doanh thu từ orders
   - Cập nhật trạng thái ca thành CLOSED
   - Lock tất cả orders đã hoàn thành
   - Hiển thị summary cho mỗi ca

**⚠️ Cảnh báo:**
- Script này sẽ đóng TẤT CẢ ca đang mở
- Không thể hoàn tác sau khi đóng
- Nên chạy `list-open-shifts.sh` trước để kiểm tra

**Khi nào dùng:**
- Cuối ngày để đóng tất cả ca
- Khi cần reset hệ thống
- Khi có lỗi và cần đóng ca thủ công

## Quy trình khuyến nghị

### Đóng ca hàng ngày

```bash
# Bước 1: Xem danh sách ca đang mở
./list-open-shifts.sh

# Bước 2: Kiểm tra thông tin
# - Xem có ca nào cần xử lý đặc biệt không
# - Kiểm tra số tiền có hợp lý không

# Bước 3: Đóng tất cả ca
./close-all-waiter-shifts.sh
# Gõ "yes" để xác nhận

# Bước 4: Xác nhận lại
./list-open-shifts.sh
# Nên thấy "No open shifts found"
```

### Debug vấn đề ca làm việc

```bash
# Xem chi tiết các ca đang mở
./list-open-shifts.sh

# Nếu thấy ca có vấn đề (tiền âm, doanh thu sai, etc.)
# Có thể:
# 1. Sửa trực tiếp trong MongoDB
# 2. Hoặc đóng ca và tạo lại
```

## Chi tiết kỹ thuật

### Script đóng ca làm gì?

1. **Tìm tất cả ca OPEN**
   ```go
   openShifts, err := shiftRepo.FindOpenShifts(ctx)
   ```

2. **Với mỗi ca:**
   - Lấy tất cả orders của ca đó
   - Tính tổng doanh thu (chỉ orders PAID/IN_PROGRESS/SERVED)
   - Phân loại doanh thu (cash/transfer)
   - Cập nhật shift:
     ```go
     shift.Status = order.ShiftClosed
     shift.EndedAt = &now
     shift.TotalRevenue = totalRevenue
     shift.TotalOrders = totalOrders
     shift.EndCash = shift.CurrentCash
     ```
   - Lock các orders đã hoàn thành (SERVED/CANCELLED)

3. **Hiển thị summary**
   - Số orders
   - Doanh thu (total, cash, transfer)
   - Tiền còn lại (cash, transfer)
   - Số orders đã lock

### Dữ liệu được cập nhật

**Shift:**
- `status`: OPEN → CLOSED
- `ended_at`: timestamp hiện tại
- `total_revenue`: tổng doanh thu từ orders
- `total_orders`: số lượng orders
- `end_cash`: = current_cash
- `updated_at`: timestamp hiện tại

**Orders:**
- `status`: SERVED/CANCELLED → LOCKED
- `locked_at`: timestamp hiện tại

### Dữ liệu KHÔNG thay đổi

- `current_cash`, `remaining_cash` - giữ nguyên
- `transfer_revenue`, `remaining_transfer` - giữ nguyên
- `handed_over_cash`, `handed_over_transfer` - giữ nguyên
- Orders với status khác (CREATED, PAID, IN_PROGRESS) - không lock

## Troubleshooting

### Lỗi: "Failed to find open shifts"

**Nguyên nhân:** Không kết nối được MongoDB

**Giải pháp:**
1. Kiểm tra MongoDB đang chạy: `mongosh`
2. Kiểm tra `.env` file có đúng MONGO_URI không
3. Kiểm tra network/firewall

### Lỗi: "Failed to close shift"

**Nguyên nhân:** Lỗi khi update shift trong DB

**Giải pháp:**
1. Kiểm tra MongoDB logs
2. Kiểm tra permissions
3. Thử update thủ công trong MongoDB:
   ```javascript
   db.shifts.updateOne(
     { _id: ObjectId("...") },
     { $set: { status: "CLOSED", ended_at: new Date() } }
   )
   ```

### Ca đã đóng nhưng vẫn hiển thị OPEN

**Nguyên nhân:** Cache hoặc DB không sync

**Giải pháp:**
1. Refresh browser
2. Kiểm tra trực tiếp trong MongoDB:
   ```javascript
   db.shifts.find({ status: "OPEN" })
   ```
3. Nếu DB đúng nhưng UI sai, clear cache browser

### Tiền bị âm sau khi đóng ca

**Nguyên nhân:** 
- Bàn giao sai (trừ nhầm loại tiền)
- Orders có payment method sai

**Giải pháp:**
1. Kiểm tra lại logic bàn giao (xem TRANSFER_HANDOVER_FIX.md)
2. Kiểm tra orders:
   ```javascript
   db.orders.find({ 
     shift_id: ObjectId("..."),
     payment_method: { $exists: true }
   })
   ```
3. Sửa lại remaining_cash/remaining_transfer nếu cần

## Files liên quan

- `backend/cmd/close-all-waiter-shifts/main.go` - Script đóng ca
- `backend/cmd/list-open-shifts/main.go` - Script xem ca
- `close-all-waiter-shifts.sh` - Bash wrapper
- `list-open-shifts.sh` - Bash wrapper
- `backend/application/services/shift_service.go` - Logic shift service
- `backend/infrastructure/mongodb/shift_repository.go` - Shift repository

## Lưu ý quan trọng

1. **Backup trước khi đóng ca hàng loạt**
   ```bash
   mongodump --db cafe_pos --out backup_$(date +%Y%m%d)
   ```

2. **Không đóng ca khi:**
   - Có orders đang pending payment
   - Có handover đang pending
   - Đang trong giờ cao điểm

3. **Sau khi đóng ca:**
   - Kiểm tra reports có đúng không
   - Xác nhận không còn ca OPEN
   - Backup database

4. **Khôi phục nếu cần:**
   ```bash
   mongorestore --db cafe_pos backup_20240223/cafe_pos
   ```
