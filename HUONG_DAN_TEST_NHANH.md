# Hướng Dẫn Test Nhanh - Bàn Giao Quỹ Thu Ngân

## 🚀 Chuẩn Bị (5 phút)

### 1. Khởi động Server

```bash
# Terminal 1: Backend
cd backend
go run main.go

# Terminal 2: Frontend
cd frontend
npm run dev
```

### 2. Đăng nhập với vai trò Cashier

1. Mở trình duyệt: `http://localhost:5173`
2. Đăng nhập với tài khoản cashier
3. Mở DevTools (F12)

### 3. Lấy JWT Token

```javascript
// Trong console của trình duyệt
localStorage.getItem('token')
// Copy giá trị token
```

---

## 🧪 Chạy Test Tự Động

```bash
# Thiết lập token
export TOKEN="token_của_bạn"

# Chạy test API
./test-fund-handover-api.sh

# Chạy test frontend
node test-frontend-fund-handover.js
```

---

## 📋 Quy Trình Test Thủ Công

### Kịch Bản 1: Không Có Chênh Lệch

1. **Dashboard**
   - ✅ Xem phần "Tiền đang quản lý"
   - ✅ Ghi nhận các số tiền

2. **Bắt Đầu Đóng Ca**
   - Nhấn "Đóng ca"
   - ✅ Xem tổng quan tiền quản lý

3. **Đếm Tiền**
   - Nhập đúng số tiền mặt lý thuyết
   - ✅ Chênh lệch = 0₫
   - ✅ Không cần ghi nhận lý do

4. **Xác Nhận**
   - Xem lại tóm tắt
   - Nhấn "Xác nhận và đóng ca"
   - ✅ Thành công!

**Thời gian**: ~2 phút

---

### Kịch Bản 2: Có Chênh Lệch

1. **Dashboard**
   - ✅ Xem tiền đang quản lý

2. **Bắt Đầu Đóng Ca**
   - Nhấn "Đóng ca"

3. **Đếm Tiền**
   - Nhập số tiền khác với lý thuyết
   - Ví dụ: Lý thuyết 2,000,000₫ → Nhập 1,995,000₫
   - ✅ Chênh lệch = -5,000₫ (thiếu)

4. **Ghi Nhận Chênh Lệch**
   - Chọn lý do: "COUNTING_ERROR"
   - Nhập ghi chú: "Đếm nhầm tờ 50k thành 100k"
   - ✅ Ghi chú ≥ 10 ký tự

5. **Xác Nhận**
   - Xem lại tóm tắt có chênh lệch
   - Nhấn "Xác nhận và đóng ca"
   - ✅ Thành công!

**Thời gian**: ~3 phút

---

## 🔍 Những Gì Cần Kiểm Tra

### Dashboard
- [ ] 💰 Phần "Tiền đang quản lý" hiển thị
- [ ] 💵 Tiền mặt (xanh lá) + 💳 Tiền CK (xanh dương)
- [ ] 📊 Tổng cộng (gradient cam)
- [ ] ⚠️ Thông báo cảnh báo
- [ ] Kéo để làm mới hoạt động

### Quy Trình Đóng Ca
- [ ] Tổng quan tiền quản lý khớp với dashboard
- [ ] Nhập tiền mặt hoạt động
- [ ] Chênh lệch tự động tính
- [ ] Form ghi nhận chênh lệch xuất hiện khi cần
- [ ] Xác nhận hiển thị đầy đủ dữ liệu
- [ ] Gửi thành công

### Mobile
- [ ] Giao diện responsive
- [ ] Tương tác cảm ứng mượt mà
- [ ] Không cuộn ngang
- [ ] Bàn phím không che input

---

## 🐛 Các Vấn Đề Thường Gặp

### Vấn Đề 1: Tiền Quản Lý Không Tải
**Triệu chứng**: Trống hoặc loading mãi
**Kiểm tra**: 
- Tab Network để xem API call
- Console để xem lỗi
- Log backend

### Vấn Đề 2: Chênh Lệch Không Tính
**Triệu chứng**: Hiển thị 0 khi phải có chênh lệch
**Kiểm tra**:
- Giá trị tiền mặt lý thuyết
- Số tiền nhập vào
- Console để xem lỗi tính toán

### Vấn Đề 3: Không Gửi Được
**Triệu chứng**: Nút bị disable hoặc lỗi
**Kiểm tra**:
- Ghi nhận chênh lệch đầy đủ
- Ghi chú ≥ 10 ký tự
- Kết nối mạng
- Lỗi console

### Vấn Đề 4: Giao Diện Mobile Bị Vỡ
**Triệu chứng**: Tràn, chữ bị cắt
**Kiểm tra**:
- Viewport meta tag
- CSS responsive classes
- Safe area insets

---

## 📊 Kiểm Tra Nhanh

### Sau Khi Đóng Ca Thành Công

**1. Kiểm Tra Frontend**
```
✅ Chuyển về dashboard
✅ Trạng thái ca = CLOSED
✅ Thông báo thành công
```

**2. Kiểm Tra Database**
```javascript
// MongoDB shell
db.fund_handovers.findOne({ cashier_shift_id: ObjectId("...") })
// Phải trả về bản ghi fund handover

db.cashier_shifts.findOne({ _id: ObjectId("...") })
// Phải có status: "CLOSED"
```

**3. Kiểm Tra API Response**
```json
{
  "shift": { "status": "CLOSED", ... },
  "fund_handover": {
    "cash_amount": 1995000,
    "transfer_amount": 800000,
    "variance_amount": -5000,
    ...
  }
}
```

---

## 🎯 Ưu Tiên Test

### Ưu Tiên 1: Chức Năng Chính (Bắt Buộc)
1. ✅ Dashboard hiển thị tiền quản lý
2. ✅ Quy trình đóng ca hoàn tất
3. ✅ Bản ghi fund handover được tạo
4. ✅ Không mất dữ liệu

### Ưu Tiên 2: Tính Năng Quan Trọng
1. ✅ Tính chênh lệch đúng
2. ✅ Ghi nhận chênh lệch hoạt động
3. ✅ Mobile responsive
4. ✅ Xử lý lỗi

### Ưu Tiên 3: Tốt Nếu Có
1. ✅ Kéo để làm mới mượt
2. ✅ Animation mượt
3. ✅ Trạng thái loading
4. ✅ Accessibility

---

## 📱 Test Mobile Nhanh

### iOS Safari
1. Kết nối iPhone qua USB
2. Safari > Develop > iPhone > localhost
3. Test tương tác cảm ứng
4. Kiểm tra safe area insets

### Android Chrome
1. Bật USB debugging
2. Chrome > chrome://inspect
3. Chọn thiết bị
4. Test tương tác cảm ứng

### Kiểm Tra Nhanh Mobile
- [ ] Vùng chạm ≥ 44x44px
- [ ] Chữ đọc được không cần zoom
- [ ] Không cuộn ngang
- [ ] Bàn phím hoạt động đúng

---

## 🔧 Xử Lý Sự Cố

### Backend Không Chạy
```bash
cd backend
go run main.go
# Phải thấy: Server running on :8080
```

### Frontend Không Chạy
```bash
cd frontend
npm run dev
# Phải thấy: Local: http://localhost:5173
```

### MongoDB Không Kết Nối
```bash
# Kiểm tra MongoDB
mongosh
# Phải kết nối thành công
```

### Token Hết Hạn
```javascript
// Lấy token mới
// 1. Đăng xuất
// 2. Đăng nhập lại
// 3. Copy token mới từ localStorage
```

---

## 📞 Cần Trợ Giúp?

### Kiểm Tra Logs
- **Frontend**: Console trình duyệt (F12)
- **Backend**: Terminal chạy `go run main.go`
- **MongoDB**: `mongosh` và kiểm tra collections

### Tài Liệu
- Hướng dẫn đầy đủ: `FRONTEND_TESTING_GUIDE.md`
- Checklist thủ công: `MANUAL_TESTING_CHECKLIST.md`
- Tài liệu API: `CASHIER_FUND_HANDOVER_PHASE_4_COMPLETE.md`

### Lệnh Thường Dùng
```bash
# Xem log backend
cd backend && go run main.go

# Xem dữ liệu MongoDB
mongosh
use cafe_pos
db.fund_handovers.find().pretty()
db.cashier_shifts.find({ status: "CLOSED" }).pretty()

# Xóa dữ liệu test
db.fund_handovers.deleteMany({})
db.cashier_shifts.updateMany({}, { $set: { status: "OPEN" } })
```

---

## ✅ Hoàn Thành Test?

Sau khi test, xác nhận:
- [ ] Tất cả chức năng chính hoạt động
- [ ] Không có lỗi console
- [ ] Mobile responsive
- [ ] Bản ghi database đúng
- [ ] Đã ghi nhận các vấn đề

Sau đó:
1. Điền `MANUAL_TESTING_CHECKLIST.md`
2. Ghi nhận các bug tìm thấy
3. Xin phê duyệt từ stakeholder
4. Sẵn sàng deploy! 🚀

---

## 📚 Tài Liệu Liên Quan

### Tiếng Việt
- `HOAN_THANH_FUND_HANDOVER.md` - Tổng quan tính năng
- `HUONG_DAN_TEST_NHANH.md` - File này

### Tiếng Anh
- `FRONTEND_TESTING_GUIDE.md` - Hướng dẫn test chi tiết
- `MANUAL_TESTING_CHECKLIST.md` - Checklist test thủ công
- `TESTING_QUICK_START.md` - Bắt đầu nhanh
- `FRONTEND_TESTING_COMPLETE.md` - Tổng kết task test

### Scripts
- `test-frontend-fund-handover.js` - Test tự động
- `test-fund-handover-api.sh` - Test API

---

## 🎉 Chúc Bạn Test Thành Công!

Nếu có vấn đề, tham khảo tài liệu chi tiết hoặc kiểm tra logs để debug.
