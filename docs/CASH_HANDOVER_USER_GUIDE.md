# 📖 Hướng Dẫn Sử Dụng Tính Năng Bàn Giao Tiền

## 🎯 Tổng Quan

Tính năng bàn giao tiền cho phép nhân viên phục vụ (waiter) bàn giao tiền mặt thu được từ khách hàng cho thu ngân (cashier) một cách an toàn và có kiểm soát.

**Lợi ích:**
- ✅ Theo dõi chính xác luồng tiền mặt
- ✅ Giảm rủi ro mất mát tiền
- ✅ Kiểm soát chênh lệch
- ✅ Audit trail đầy đủ
- ✅ Quy trình rõ ràng, minh bạch

---

## 👨‍💼 Hướng Dẫn Cho Nhân Viên Phục Vụ (Waiter)

### 1. Xem Trạng Thái Tiền Mặt

Khi bạn đang trong ca làm việc, bạn sẽ thấy 3 thẻ hiển thị:

- **Tiền hiện có:** Số tiền bạn đang giữ
- **Đã bàn giao:** Tổng tiền đã bàn giao cho thu ngân
- **Tổng thu:** Tổng tiền đã thu từ khách hàng

### 2. Bàn Giao Một Phần Tiền

**Khi nào cần bàn giao một phần?**
- Khi bạn giữ quá nhiều tiền mặt (> 500k)
- Khi muốn giảm rủi ro
- Theo quy định của quán

**Các bước thực hiện:**

1. Nhấn nút **"💰 Bàn giao một phần"**
2. Nhập số tiền muốn bàn giao (≤ tiền hiện có)
3. Thêm ghi chú (tùy chọn)
4. Nhấn **"Bàn giao"**
5. Chờ thu ngân xác nhận

**Lưu ý:**
- Bạn chỉ có thể bàn giao số tiền ≤ tiền hiện có
- Sau khi gửi yêu cầu, bạn không thể tạo yêu cầu mới cho đến khi thu ngân xác nhận
- Bạn có thể hủy yêu cầu nếu chưa được xác nhận

### 3. Hủy Yêu Cầu Bàn Giao

Nếu bạn muốn hủy yêu cầu đang chờ:

1. Nhìn thấy banner màu vàng "🕐 Đang chờ xác nhận bàn giao"
2. Nhấn nút **"Hủy"**
3. Xác nhận hủy
4. Yêu cầu sẽ bị hủy và bạn có thể tạo yêu cầu mới

### 4. Bàn Giao Toàn Bộ và Đóng Ca

**Khi nào sử dụng?**
- Khi kết thúc ca làm việc
- Khi cần bàn giao toàn bộ tiền còn lại

**Các bước thực hiện:**

1. Nhấn nút **"🏁 Bàn giao và đóng ca"**
2. Đọc cảnh báo (ca sẽ tự động đóng sau khi thu ngân xác nhận)
3. Nhập tiền cuối ca (thường là 0)
4. Thêm ghi chú (tùy chọn)
5. Nhấn **"Bàn giao và đóng ca"**
6. Chờ thu ngân xác nhận
7. Ca sẽ tự động đóng sau khi xác nhận

**Lưu ý:**
- Thao tác này không thể hoàn tác
- Ca sẽ tự động đóng sau khi thu ngân xác nhận
- Tất cả orders trong ca sẽ bị khóa

### 5. Xem Lịch Sử Bàn Giao

Bạn có thể xem tất cả các lần bàn giao trong ca:

- Số tiền đã bàn giao
- Thời gian bàn giao
- Trạng thái (Chờ xác nhận / Đã xác nhận / Đã từ chối)
- Ghi chú của bạn
- Phản hồi từ thu ngân
- Chênh lệch (nếu có)

**Các trạng thái:**
- 🟡 **Chờ xác nhận:** Thu ngân chưa xử lý
- 🟢 **Đã xác nhận:** Thu ngân đã nhận tiền
- 🔴 **Đã từ chối:** Thu ngân từ chối (xem lý do trong ghi chú)
- 🟠 **Có chênh lệch:** Có sự khác biệt giữa số tiền khai báo và thực nhận

### 6. Xử Lý Chênh Lệch

Nếu có chênh lệch giữa số tiền bạn khai báo và số tiền thu ngân đếm được:

**Chênh lệch nhỏ (< 100k):**
- Thu ngân sẽ ghi nhận và xác nhận
- Bạn sẽ thấy thông tin chênh lệch trong lịch sử
- Số tiền thực tế sẽ được cập nhật

**Chênh lệch lớn (≥ 100k):**
- Cần sự phê duyệt từ quản lý
- Trạng thái sẽ là "Có chênh lệch"
- Chờ quản lý xem xét và phê duyệt

**Lưu ý:**
- Luôn đếm kỹ tiền trước khi bàn giao
- Ghi chú rõ ràng nếu có vấn đề
- Liên hệ quản lý nếu chênh lệch lớn

---

## 💰 Hướng Dẫn Cho Thu Ngân (Cashier)

### 1. Nhận Thông Báo

Khi có yêu cầu bàn giao mới:

- Dashboard sẽ hiển thị banner màu vàng
- Số lượng yêu cầu đang chờ
- Nút **"Xem ngay"** để xem chi tiết

### 2. Xem Danh Sách Chờ Xác Nhận

Trên trang **"Quản lý bàn giao"**, bạn sẽ thấy:

- Tên nhân viên phục vụ
- Số tiền khai báo
- Loại bàn giao (Một phần / Toàn bộ + Đóng ca)
- Thời gian tạo yêu cầu
- Ghi chú từ nhân viên

### 3. Xác Nhận Bàn Giao

**Các bước thực hiện:**

1. Nhấn nút **"✅ Xác nhận"** trên yêu cầu
2. Đếm tiền thực tế nhận được
3. Nhập số tiền thực nhận vào form
4. Thêm ghi chú (tùy chọn)
5. Nhấn **"Xác nhận"**

**Trường hợp không có chênh lệch:**
- Nhập số tiền = số tiền khai báo
- Thêm ghi chú xác nhận
- Xác nhận

**Trường hợp có chênh lệch:**
- Nhập số tiền thực nhận (khác số tiền khai báo)
- Hệ thống tự động tính chênh lệch
- Chọn lý do chênh lệch:
  - Lỗi đếm tiền
  - Lỗi giao dịch
  - Vấn đề khách hàng
  - Lỗi hệ thống
  - Khác
- Chọn trách nhiệm:
  - Nhân viên phục vụ
  - Thu ngân
  - Khách hàng
  - Hệ thống
  - Không rõ
- Thêm ghi chú giải thích
- Xác nhận

**Chênh lệch lớn (≥ 100k):**
- Hệ thống sẽ cảnh báo
- Yêu cầu cần phê duyệt từ quản lý
- Tiền chưa được cập nhật cho đến khi quản lý phê duyệt

### 4. Từ Chối Bàn Giao

Nếu có vấn đề với yêu cầu:

1. Nhấn nút **"❌ Từ chối"**
2. Nhập lý do từ chối (bắt buộc)
3. Nhấn **"Từ chối"**
4. Nhân viên sẽ nhận được thông báo

**Khi nào nên từ chối?**
- Số tiền không khớp với hệ thống
- Có dấu hiệu bất thường
- Cần xác minh thêm thông tin

### 5. Xác Nhận Nhanh (Quick Confirm)

Trên Dashboard, bạn có thể xác nhận nhanh:

1. Xem danh sách "⚡ Bàn giao nhanh"
2. Nhấn **"✅"** để xác nhận (giả định số tiền đúng)
3. Nhấn **"❌"** để từ chối
4. Không cần nhập chi tiết

**Khi nào sử dụng?**
- Khi số tiền chính xác
- Khi cần xử lý nhanh
- Khi tin tưởng nhân viên

### 6. Xem Lịch Sử Hôm Nay

Phần **"Bàn giao hôm nay"** hiển thị:

- Tất cả bàn giao đã xử lý hôm nay
- Trạng thái của từng bàn giao
- Ghi chú của bạn
- Thông tin chênh lệch (nếu có)

---

## 👨‍💼 Hướng Dẫn Cho Quản Lý (Manager)

### 1. Xem Yêu Cầu Phê Duyệt

Khi có chênh lệch lớn (≥ 100k):

1. Truy cập trang **"Pending Approvals"**
2. Xem danh sách yêu cầu cần phê duyệt
3. Xem chi tiết:
   - Nhân viên phục vụ
   - Thu ngân
   - Số tiền khai báo
   - Số tiền thực nhận
   - Chênh lệch
   - Lý do chênh lệch
   - Trách nhiệm

### 2. Phê Duyệt Chênh Lệch

**Các bước thực hiện:**

1. Xem xét thông tin chi tiết
2. Xác minh với nhân viên liên quan
3. Quyết định phê duyệt hoặc từ chối
4. Nhập ghi chú giải thích quyết định
5. Nhấn **"Approve"** hoặc **"Reject"**

**Nếu phê duyệt:**
- Số tiền sẽ được cập nhật
- Trạng thái chuyển thành "Đã xác nhận"
- Chênh lệch được ghi nhận

**Nếu từ chối:**
- Trạng thái chuyển thành "Đã từ chối"
- Số tiền không được cập nhật
- Cần xử lý lại với nhân viên

### 3. Xem Thống Kê Chênh Lệch

Truy cập **"Discrepancy Stats"** để xem:

- Tổng số lần bàn giao
- Tổng chênh lệch
- Số lần thiếu tiền
- Số lần thừa tiền
- Số tiền thiếu
- Số tiền thừa
- Số lần cần phê duyệt

**Sử dụng để:**
- Đánh giá hiệu suất nhân viên
- Phát hiện vấn đề hệ thống
- Cải thiện quy trình
- Đào tạo nhân viên

---

## 💡 Mẹo và Thực Hành Tốt

### Cho Nhân Viên Phục Vụ

1. **Đếm kỹ tiền trước khi bàn giao**
   - Đếm ít nhất 2 lần
   - Tách riêng các mệnh giá
   - Ghi chú nếu có tiền lẻ

2. **Bàn giao thường xuyên**
   - Không giữ quá nhiều tiền (> 500k)
   - Bàn giao sau mỗi 3-4 orders lớn
   - Giảm rủi ro mất mát

3. **Ghi chú rõ ràng**
   - Ghi rõ thời gian thu tiền
   - Ghi chú nếu có vấn đề
   - Dễ tra cứu sau này

4. **Kiểm tra lịch sử**
   - Xem lại các lần bàn giao
   - Học từ chênh lệch
   - Cải thiện độ chính xác

### Cho Thu Ngân

1. **Đếm tiền cẩn thận**
   - Đếm trước mặt nhân viên (nếu có thể)
   - Sử dụng máy đếm tiền nếu có
   - Xác nhận từng mệnh giá

2. **Ghi chú chi tiết**
   - Ghi rõ tình trạng tiền
   - Ghi chú nếu có vấn đề
   - Giúp audit sau này

3. **Xử lý nhanh**
   - Không để nhân viên chờ lâu
   - Sử dụng quick confirm khi phù hợp
   - Ưu tiên yêu cầu cũ trước

4. **Báo cáo chênh lệch lớn**
   - Thông báo quản lý ngay
   - Ghi chú chi tiết nguyên nhân
   - Đề xuất giải pháp

### Cho Quản Lý

1. **Theo dõi thống kê**
   - Xem báo cáo hàng ngày
   - Phát hiện xu hướng
   - Hành động kịp thời

2. **Đào tạo nhân viên**
   - Hướng dẫn cách đếm tiền
   - Giải thích tầm quan trọng
   - Chia sẻ best practices

3. **Cải thiện quy trình**
   - Lắng nghe phản hồi
   - Điều chỉnh threshold nếu cần
   - Tối ưu hóa workflow

4. **Xử lý công bằng**
   - Xem xét kỹ từng trường hợp
   - Không vội kết luận
   - Tạo môi trường tin tưởng

---

## ❓ Câu Hỏi Thường Gặp (FAQ)

### Q1: Tôi có thể bàn giao bao nhiêu lần trong một ca?
**A:** Không giới hạn. Bạn có thể bàn giao nhiều lần tùy nhu cầu.

### Q2: Nếu thu ngân không xác nhận thì sao?
**A:** Bạn có thể hủy yêu cầu và tạo yêu cầu mới, hoặc liên hệ quản lý.

### Q3: Chênh lệch bao nhiêu thì cần phê duyệt quản lý?
**A:** Chênh lệch ≥ 100,000 VND cần phê duyệt quản lý.

### Q4: Tôi có thể sửa số tiền sau khi gửi yêu cầu không?
**A:** Không. Bạn cần hủy yêu cầu cũ và tạo yêu cầu mới.

### Q5: Nếu tôi quên bàn giao tiền trước khi đóng ca?
**A:** Sử dụng chức năng "Bàn giao và đóng ca" để bàn giao toàn bộ tiền còn lại.

### Q6: Lịch sử bàn giao được lưu bao lâu?
**A:** Vĩnh viễn. Tất cả bàn giao đều được lưu trữ để audit.

### Q7: Tôi có thể xem bàn giao của người khác không?
**A:** Không. Bạn chỉ xem được bàn giao của mình (waiter) hoặc được phân công (cashier).

### Q8: Nếu có tranh chấp về số tiền thì sao?
**A:** Liên hệ quản lý ngay. Quản lý sẽ xem xét và quyết định.

---

## 🆘 Hỗ Trợ

Nếu bạn gặp vấn đề hoặc có câu hỏi:

1. **Xem hướng dẫn này trước**
2. **Liên hệ quản lý ca**
3. **Gọi hotline hỗ trợ:** 1900-xxxx
4. **Email:** support@cafepos.com

---

**Phiên bản:** 1.0  
**Cập nhật:** 04/02/2026  
**Người soạn:** Đội ngũ phát triển Café POS
