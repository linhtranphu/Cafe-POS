# Batch Ingredient Management - User Acceptance Testing Guide

## Mục Đích

Tài liệu này hướng dẫn kiểm thử chấp nhận người dùng (UAT) cho hệ thống Quản Lý Nguyên Liệu Batch. UAT đảm bảo hệ thống đáp ứng nhu cầu thực tế của người dùng cuối.

## Thông Tin Hệ Thống

- **Backend API**: http://localhost:8080
- **Frontend**: http://localhost:5173
- **Test Accounts**:
  - Manager: `manager@test.com` / `password123`
  - Barista: `barista@test.com` / `password123`

## Chuẩn Bị Trước Khi Test

### 1. Khởi Động Hệ Thống

```bash
# Terminal 1: Start MongoDB
docker-compose up -d mongodb

# Terminal 2: Start Backend
cd backend
go run main.go

# Terminal 3: Start Frontend
cd frontend
npm run dev
```

### 2. Seed Dữ Liệu Test

```bash
# Tạo admin user
cd backend
go run cmd/seed-admin/main.go

# Tạo nguyên liệu mẫu
go run cmd/seed/main.go

# Tạo menu items
go run cmd/seed-menu/main.go
```

### 3. Xác Nhận Hệ Thống Hoạt Động

- [ ] Backend API responding: `curl http://localhost:8080/health`
- [ ] Frontend loading: Mở `http://localhost:5173`
- [ ] Login thành công với tài khoản test
- [ ] Dashboard hiển thị đúng

---

## PHẦN 1: MANAGER TESTING

### Test Case 1.1: Quản Lý Batch Definition

**Mục tiêu**: Kiểm tra khả năng tạo, xem, sửa, xóa batch definition

#### 1.1.1 Tạo Batch Definition Mới

**Bước thực hiện**:
1. Đăng nhập với tài khoản Manager
2. Vào menu "Batch" → "Định Nghĩa Batch"
3. Click nút "Tạo Định Nghĩa Mới"
4. Nhập thông tin:
   - Tên: "Cà Phê Concentrate"
   - Đơn vị: "ml"
   - Thời gian bảo quản: 24 giờ
   - Ngưỡng tồn kho thấp: 500 ml
   - Cảnh báo hết hạn: 4 giờ
5. Thêm nguyên liệu nguồn:
   - Chọn "Hạt Cà Phê"
   - Số lượng nguồn: 100g
   - Số lượng batch tạo ra: 500ml
   - Tỷ lệ hao hụt: 10%
6. Click "Lưu"

**Kết quả mong đợi**:
- [ ] Form validation hoạt động (không cho submit nếu thiếu thông tin)
- [ ] Hiển thị preview chi phí dự kiến
- [ ] Batch definition được tạo thành công
- [ ] Hiển thị thông báo thành công
- [ ] Batch definition xuất hiện trong danh sách
- [ ] Chi phí được tính đúng (bao gồm wastage)

**Ghi chú lỗi** (nếu có):
```
[Ghi chú các vấn đề phát hiện ở đây]
```

#### 1.1.2 Xem Danh Sách Batch Definitions

**Bước thực hiện**:
1. Vào "Batch" → "Định Nghĩa Batch"
2. Kiểm tra danh sách hiển thị

**Kết quả mong đợi**:
- [ ] Hiển thị tất cả batch definitions
- [ ] Thông tin đầy đủ: Tên, Đơn vị, Thời gian bảo quản
- [ ] Search box hoạt động
- [ ] Có thể tìm kiếm theo tên
- [ ] Responsive trên mobile

**Ghi chú lỗi**:
```
[Ghi chú các vấn đề phát hiện ở đây]
```

#### 1.1.3 Chỉnh Sửa Batch Definition

**Bước thực hiện**:
1. Click nút "Sửa" trên một batch definition
2. Thay đổi thông tin (ví dụ: tăng ngưỡng tồn kho thấp lên 600ml)
3. Click "Lưu"

**Kết quả mong đợi**:
- [ ] Form hiển thị đúng thông tin hiện tại
- [ ] Có thể chỉnh sửa tất cả các trường
- [ ] Lưu thành công
- [ ] Thông tin cập nhật trong danh sách

**Ghi chú lỗi**:
```
[Ghi chú các vấn đề phát hiện ở đây]
```

#### 1.1.4 Xóa Batch Definition

**Bước thực hiện**:
1. Click nút "Xóa" trên một batch definition chưa có batch record
2. Xác nhận xóa

**Kết quả mong đợi**:
- [ ] Hiển thị dialog xác nhận
- [ ] Xóa thành công
- [ ] Batch definition biến mất khỏi danh sách
- [ ] Không cho xóa nếu đã có batch record sử dụng

**Ghi chú lỗi**:
```
[Ghi chú các vấn đề phát hiện ở đây]
```

---

### Test Case 1.2: Xem Báo Cáo

**Mục tiêu**: Kiểm tra các báo cáo sản xuất, lãng phí, sử dụng

#### 1.2.1 Báo Cáo Sản Xuất

**Bước thực hiện**:
1. Vào "Batch" → "Báo Cáo" → "Sản Xuất"
2. Chọn khoảng thời gian (ví dụ: 7 ngày qua)
3. Xem báo cáo

**Kết quả mong đợi**:
- [ ] Hiển thị tổng số batch đã sản xuất
- [ ] Hiển thị tổng số lượng và chi phí
- [ ] Chart hiển thị xu hướng sản xuất theo thời gian
- [ ] Bảng phân tích theo loại batch
- [ ] Bảng phân tích theo người chế biến
- [ ] Có thể filter theo batch type và preparer
- [ ] Có thể export CSV

**Ghi chú lỗi**:
```
[Ghi chú các vấn đề phát hiện ở đây]
```

#### 1.2.2 Báo Cáo Lãng Phí

**Bước thực hiện**:
1. Vào "Batch" → "Báo Cáo" → "Lãng Phí"
2. Chọn khoảng thời gian
3. Xem báo cáo

**Kết quả mong đợi**:
- [ ] Hiển thị tổng số batch hết hạn
- [ ] Hiển thị tổng số lượng và chi phí lãng phí
- [ ] Chart hiển thị xu hướng lãng phí
- [ ] Bảng phân tích theo loại batch
- [ ] Hiển thị recommendations để giảm lãng phí
- [ ] Có thể export CSV

**Ghi chú lỗi**:
```
[Ghi chú các vấn đề phát hiện ở đây]
```

#### 1.2.3 Báo Cáo Sử Dụng

**Bước thực hiện**:
1. Vào "Batch" → "Báo Cáo" → "Sử Dụng"
2. Chọn khoảng thời gian
3. Xem báo cáo

**Kết quả mong đợi**:
- [ ] Hiển thị tổng số lần sử dụng
- [ ] Hiển thị tổng số lượng và chi phí
- [ ] Chart hiển thị xu hướng sử dụng
- [ ] Bảng phân tích theo menu item
- [ ] Ranking batch được sử dụng nhiều nhất
- [ ] Có thể filter theo batch type và menu item
- [ ] Có thể export CSV

**Ghi chú lỗi**:
```
[Ghi chú các vấn đề phát hiện ở đây]
```

---

### Test Case 1.3: Quản Lý Batch Records

**Mục tiêu**: Kiểm tra khả năng xem, điều chỉnh, xóa batch records

#### 1.3.1 Xem Danh Sách Batch Records

**Bước thực hiện**:
1. Vào "Batch" → "Batch Records"
2. Xem danh sách

**Kết quả mong đợi**:
- [ ] Hiển thị tất cả batch records
- [ ] Color coding đúng (xanh: available, vàng: low stock, đỏ: expiring, xám: expired)
- [ ] Filters hoạt động (batch type, status, date range, preparer)
- [ ] Sorting hoạt động (expiry date, prepared date, quantity)
- [ ] Pagination hoạt động
- [ ] Responsive trên mobile

**Ghi chú lỗi**:
```
[Ghi chú các vấn đề phát hiện ở đây]
```

#### 1.3.2 Xem Chi Tiết Batch Record

**Bước thực hiện**:
1. Click vào một batch record
2. Xem chi tiết

**Kết quả mong đợi**:
- [ ] Hiển thị đầy đủ thông tin batch
- [ ] Hiển thị breakdown nguyên liệu đã sử dụng
- [ ] Hiển thị breakdown chi phí
- [ ] Hiển thị lịch sử sử dụng (orders đã dùng batch này)
- [ ] Timeline visualization (Prepared → Used → Expired)
- [ ] Actions: Mark as Expired, Delete (nếu chưa dùng)

**Ghi chú lỗi**:
```
[Ghi chú các vấn đề phát hiện ở đây]
```

#### 1.3.3 Đánh Dấu Batch Hết Hạn

**Bước thực hiện**:
1. Vào chi tiết một batch record available
2. Click "Đánh Dấu Hết Hạn"
3. Xác nhận

**Kết quả mong đợi**:
- [ ] Hiển thị dialog xác nhận
- [ ] Status chuyển sang "expired"
- [ ] Batch không còn khả dụng để sử dụng
- [ ] Xuất hiện trong expired alerts

**Ghi chú lỗi**:
```
[Ghi chú các vấn đề phát hiện ở đây]
```

#### 1.3.4 Xóa Batch Record

**Bước thực hiện**:
1. Vào chi tiết một batch record chưa được sử dụng
2. Click "Xóa"
3. Xác nhận

**Kết quả mong đợi**:
- [ ] Hiển thị dialog xác nhận
- [ ] Xóa thành công
- [ ] Nguyên liệu nguồn được hoàn trả vào inventory
- [ ] Không cho xóa nếu batch đã được sử dụng một phần

**Ghi chú lỗi**:
```
[Ghi chú các vấn đề phát hiện ở đây]
```

---

## PHẦN 2: BARISTA TESTING

### Test Case 2.1: Tạo Batch Record

**Mục tiêu**: Kiểm tra quy trình chế biến batch từ góc nhìn barista

#### 2.1.1 Tạo Batch Thành Công

**Bước thực hiện**:
1. Đăng nhập với tài khoản Barista
2. Vào "Batch" → "Tạo Batch"
3. Chọn batch definition: "Cà Phê Concentrate"
4. Nhập số lượng: 500ml
5. Xem preview nguyên liệu cần và chi phí
6. Click "Xác Nhận"

**Kết quả mong đợi**:
- [ ] Form đơn giản, dễ sử dụng
- [ ] Hiển thị rõ ràng nguyên liệu cần: "110g Hạt Cà Phê" (100g + 10% wastage)
- [ ] Hiển thị chi phí dự kiến
- [ ] Confirmation dialog với breakdown chi tiết
- [ ] Tạo batch thành công
- [ ] Nguyên liệu nguồn bị trừ đúng số lượng
- [ ] Batch record xuất hiện trong danh sách
- [ ] Thời gian hết hạn được tính đúng (24 giờ sau)

**Ghi chú lỗi**:
```
[Ghi chú các vấn đề phát hiện ở đây]
```

#### 2.1.2 Tạo Batch Khi Không Đủ Nguyên Liệu

**Bước thực hiện**:
1. Chọn batch definition cần nhiều nguyên liệu
2. Nhập số lượng lớn (đảm bảo không đủ nguyên liệu)
3. Click "Xác Nhận"

**Kết quả mong đợi**:
- [ ] Hiển thị error message rõ ràng
- [ ] Thông báo chính xác nguyên liệu nào thiếu và thiếu bao nhiêu
- [ ] Ví dụ: "Không đủ Hạt Cà Phê. Cần: 1100g, Có: 500g"
- [ ] Không tạo batch record
- [ ] Không trừ nguyên liệu

**Ghi chú lỗi**:
```
[Ghi chú các vấn đề phát hiện ở đây]
```

---

### Test Case 2.2: Xem Alerts

**Mục tiêu**: Kiểm tra hệ thống cảnh báo

#### 2.2.1 Xem Tất Cả Alerts

**Bước thực hiện**:
1. Vào "Batch" → "Cảnh Báo"
2. Xem panel alerts

**Kết quả mong đợi**:
- [ ] Hiển thị 3 sections: Low Stock, Expiring Soon, Expired
- [ ] Badge hiển thị số lượng alerts mỗi loại
- [ ] Sections có thể expand/collapse
- [ ] Auto-refresh mỗi 5 phút
- [ ] Responsive trên mobile

**Ghi chú lỗi**:
```
[Ghi chú các vấn đề phát hiện ở đây]
```

#### 2.2.2 Low Stock Alert

**Bước thực hiện**:
1. Tạo batch cho đến khi tổng số lượng available < threshold
2. Kiểm tra Low Stock section

**Kết quả mong đợi**:
- [ ] Alert xuất hiện khi tổng số lượng <= threshold
- [ ] Hiển thị: Batch name, Current stock, Threshold
- [ ] Icon và màu vàng
- [ ] Action button: "Chế Biến Thêm" (link đến create batch form)

**Ghi chú lỗi**:
```
[Ghi chú các vấn đề phát hiện ở đây]
```

#### 2.2.3 Expiring Soon Alert

**Bước thực hiện**:
1. Tạo batch với shelf life ngắn (hoặc đợi batch sắp hết hạn)
2. Kiểm tra Expiring Soon section

**Kết quả mong đợi**:
- [ ] Alert xuất hiện khi thời gian còn lại <= expiry warning hours
- [ ] Hiển thị: Batch name, Quantity remaining, Expires at, Hours until expiry
- [ ] Icon và màu cam/đỏ
- [ ] Action button: "Sử Dụng Ngay"

**Ghi chú lỗi**:
```
[Ghi chú các vấn đề phát hiện ở đây]
```

#### 2.2.4 Expired Alert

**Bước thực hiện**:
1. Đợi batch hết hạn (hoặc manually mark as expired)
2. Kiểm tra Expired section

**Kết quả mong đợi**:
- [ ] Alert xuất hiện khi batch hết hạn
- [ ] Hiển thị: Batch name, Quantity wasted, Cost wasted, Expired at
- [ ] Icon và màu xám
- [ ] Dismiss button để ẩn alert

**Ghi chú lỗi**:
```
[Ghi chú các vấn đề phát hiện ở đây]
```

---

### Test Case 2.3: Mobile Usability

**Mục tiêu**: Kiểm tra trải nghiệm trên mobile device

#### 2.3.1 Test Trên Mobile

**Bước thực hiện**:
1. Mở app trên mobile device hoặc responsive mode
2. Thực hiện các thao tác chính:
   - Tạo batch record
   - Xem danh sách batch records
   - Xem alerts
   - Xem chi tiết batch

**Kết quả mong đợi**:
- [ ] Layout responsive, không bị vỡ
- [ ] Text đọc được, không quá nhỏ
- [ ] Buttons đủ lớn để tap
- [ ] Forms dễ nhập liệu
- [ ] Tables/lists scroll được
- [ ] Navigation menu hoạt động tốt
- [ ] Loading states hiển thị rõ ràng

**Ghi chú lỗi**:
```
[Ghi chú các vấn đề phát hiện ở đây]
```

---

## PHẦN 3: INTEGRATION TESTING

### Test Case 3.1: Batch Trong Menu

**Mục tiêu**: Kiểm tra tích hợp batch với menu system

#### 3.1.1 Thêm Batch Vào Recipe

**Bước thực hiện**:
1. Vào "Menu" → Chọn một menu item
2. Click "Sửa"
3. Trong recipe editor, thêm ingredient mới
4. Toggle sang "Batch"
5. Chọn "Cà Phê Concentrate"
6. Nhập số lượng: 30ml
7. Lưu

**Kết quả mong đợi**:
- [ ] Toggle "Nguyên Liệu Thô" / "Batch" hoạt động
- [ ] Batch selector hiển thị tất cả batch definitions
- [ ] Hiển thị số lượng batch available
- [ ] Warning nếu batch không đủ
- [ ] Chi phí menu item được tính từ batch cost
- [ ] Lưu thành công

**Ghi chú lỗi**:
```
[Ghi chú các vấn đề phát hiện ở đây]
```

#### 3.1.2 Tạo Order Sử Dụng Batch

**Bước thực hiện**:
1. Tạo order với menu item có sử dụng batch
2. Xác nhận order

**Kết quả mong đợi**:
- [ ] Order tạo thành công
- [ ] Batch quantity bị trừ đúng số lượng
- [ ] FIFO: Batch cũ nhất (sắp hết hạn nhất) được dùng trước
- [ ] Usage log được tạo
- [ ] Chi phí order sử dụng batch cost thực tế
- [ ] Nếu không đủ batch, order bị reject với error message rõ ràng

**Ghi chú lỗi**:
```
[Ghi chú các vấn đề phát hiện ở đây]
```

---

### Test Case 3.2: Dashboard Widget

**Mục tiêu**: Kiểm tra batch status widget trên dashboard

#### 3.2.1 Xem Dashboard Widget

**Bước thực hiện**:
1. Vào Dashboard
2. Tìm Batch Status Widget

**Kết quả mong đợi**:
- [ ] Widget hiển thị trên dashboard
- [ ] Hiển thị tổng số batches
- [ ] Hiển thị available quantity
- [ ] Hiển thị alert count badges (low stock, expiring, expired)
- [ ] Quick links: Create Batch, View Alerts, View Reports
- [ ] Mini chart: Usage trend (7 ngày qua)
- [ ] Compact design, không chiếm quá nhiều không gian

**Ghi chú lỗi**:
```
[Ghi chú các vấn đề phát hiện ở đây]
```

---

## PHẦN 4: ERROR HANDLING & EDGE CASES

### Test Case 4.1: Network Errors

**Bước thực hiện**:
1. Tắt backend server
2. Thử tạo batch record
3. Bật lại backend
4. Retry

**Kết quả mong đợi**:
- [ ] Hiển thị error message rõ ràng
- [ ] Có retry button
- [ ] Không crash app
- [ ] Khi retry thành công, hoạt động bình thường

**Ghi chú lỗi**:
```
[Ghi chú các vấn đề phát hiện ở đây]
```

---

### Test Case 4.2: Concurrent Operations

**Bước thực hiện**:
1. Mở 2 browser tabs
2. Tab 1: Tạo batch record sử dụng nguyên liệu X
3. Tab 2: Đồng thời tạo batch record khác cũng sử dụng nguyên liệu X
4. Submit cả 2 gần như cùng lúc

**Kết quả mong đợi**:
- [ ] Một trong hai thành công
- [ ] Cái còn lại fail với error "không đủ nguyên liệu"
- [ ] Không có race condition
- [ ] Inventory quantity chính xác

**Ghi chú lỗi**:
```
[Ghi chú các vấn đề phát hiện ở đây]
```

---

## FEEDBACK COLLECTION

### Manager Feedback

**Câu hỏi**:
1. Giao diện có dễ sử dụng không?
2. Báo cáo có cung cấp đủ thông tin cần thiết không?
3. Có tính năng nào thiếu không?
4. Có tính năng nào không cần thiết không?
5. Performance có chấp nhận được không?

**Phản hồi**:
```
[Ghi phản hồi của manager ở đây]
```

---

### Barista Feedback

**Câu hỏi**:
1. Quy trình tạo batch có nhanh và đơn giản không?
2. Alerts có hữu ích không?
3. Mobile experience có tốt không?
4. Có gặp khó khăn gì không?
5. Đề xuất cải thiện?

**Phản hồi**:
```
[Ghi phản hồi của barista ở đây]
```

---

## BUG TRACKING

### Bugs Phát Hiện

| ID | Severity | Mô tả | Steps to Reproduce | Status |
|----|----------|-------|-------------------|--------|
| 1  |          |       |                   |        |
| 2  |          |       |                   |        |
| 3  |          |       |                   |        |

**Severity Levels**:
- **Critical**: Hệ thống không sử dụng được
- **High**: Tính năng chính không hoạt động
- **Medium**: Tính năng phụ có vấn đề
- **Low**: UI/UX issues, không ảnh hưởng chức năng

---

## ACCEPTANCE CRITERIA

### Criteria for Sign-off

- [ ] Tất cả test cases PASS hoặc có workaround
- [ ] Không có Critical/High bugs
- [ ] Manager và Barista đồng ý sign-off
- [ ] Performance chấp nhận được (< 2s cho mọi thao tác)
- [ ] Mobile usability tốt
- [ ] Documentation đầy đủ

### Sign-off

**Manager**: _________________ Date: _______

**Barista**: _________________ Date: _______

**QA Lead**: _________________ Date: _______

---

## NEXT STEPS

Sau khi UAT hoàn thành:

1. **Fix Bugs**: Sửa tất cả bugs phát hiện (ưu tiên Critical/High)
2. **Implement Improvements**: Thực hiện các cải thiện được đề xuất
3. **Re-test**: Test lại các bugs đã fix
4. **Documentation**: Hoàn thiện user guides
5. **Deployment**: Chuẩn bị deploy lên production

