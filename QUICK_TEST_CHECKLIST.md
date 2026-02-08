# ✅ Quick Test Checklist - Điều Chỉnh Nguyên Liệu

## 🎯 URL Test
```
http://localhost:5173/#/ingredients
```

---

## ✅ Test 1: Adjust Giảm - Giá KHÔNG Đổi (2 phút)

1. Chọn nguyên liệu bất kỳ
2. Ghi nhớ: **Đơn giá hiện tại** = ?
3. Click "📦 Điều chỉnh" → Chọn "Điều chỉnh"
4. Nhập số lượng **GIẢM** (VD: 10 → 8)
5. **KHÔNG nhập giá** (để 0)
6. Nhập lý do: "Test giảm"
7. Click "Xác nhận"

**✅ Kết quả mong đợi**:
- Số lượng giảm
- **Đơn giá KHÔNG ĐỔI**

---

## ✅ Test 2: Adjust Tăng Không Nhập Giá - KHÔNG Tạo Expense + Lịch Sử Ghi 0đ (3 phút)

1. Chọn nguyên liệu
2. Ghi nhớ: **Đơn giá hiện tại** = ?
3. Đếm số expense hiện tại trong menu "Chi phí"
4. Click "📦 Điều chỉnh" → Chọn "Điều chỉnh"
5. Nhập số lượng **TĂNG** (VD: 8 → 12)
6. **KHÔNG nhập giá** (để 0)
7. Nhập lý do: "Được tặng"
8. Click "Xác nhận"
9. Click "📊 Lịch sử" để xem lịch sử mới nhất
10. Vào menu "Chi phí" kiểm tra

**✅ Kết quả mong đợi**:
- Số lượng tăng
- **Đơn giá KHÔNG ĐỔI**
- **KHÔNG có expense mới** (số expense không tăng)
- **Lịch sử ghi: Đơn giá 0đ, Tổng chi phí 0đ** (không phải giá hiện tại)

---

## ✅ Test 3: Adjust Tăng Có Nhập Giá - TẠO Expense + Lịch Sử Ghi Đúng (3 phút)

1. Chọn nguyên liệu
2. Ghi nhớ: **Đơn giá hiện tại** = ?
3. Đếm số expense hiện tại
4. Click "📦 Điều chỉnh" → Chọn "Điều chỉnh"
5. Nhập số lượng **TĂNG** (VD: 12 → 15, tăng 3)
6. **Nhập giá mới** (VD: 30,000đ)
7. Nhập lý do: "Mua thêm"
8. Click "Xác nhận"
9. Click "📊 Lịch sử" để xem lịch sử mới nhất
10. Vào menu "Chi phí" kiểm tra

**✅ Kết quả mong đợi**:
- Số lượng tăng
- **Đơn giá thay đổi** (weighted average)
- **CÓ expense mới** (số expense tăng 1)
- Expense amount = số lượng tăng × giá nhập (VD: 3 × 30,000 = 90,000đ)
- **Lịch sử ghi: Đơn giá 30,000đ, Tổng chi phí 90,000đ**

---

## ✅ Test 4: Form "Sửa" - KHÔNG Thể Thay Đổi Tồn Kho (1 phút)

1. Click nút "✏️ Sửa"
2. Thử click vào field "Số lượng nhập"

**✅ Kết quả mong đợi**:
- Field số lượng **màu xám** (disabled)
- **KHÔNG thể nhập** hoặc thay đổi
- Có warning: "Dùng 'Điều chỉnh' để thay đổi tồn kho"
- **KHÔNG có section nhập giá**
- Thay vào đó có "📊 Thông tin hiện tại (chỉ xem)"

---

## ✅ Test 5: UI Chỉ Có 4 Nút (30 giây)

1. Xem danh sách nguyên liệu
2. Đếm số nút action

**✅ Kết quả mong đợi**:
- Chỉ có **4 nút**: 📦 Điều chỉnh, 📊 Lịch sử, ✏️ Sửa, 🗑️ Xóa
- **KHÔNG có** nút "Nhập nhanh" và "Xuất nhanh"

---

## 📊 Bảng Tổng Hợp

| Test | Thời gian | Kết quả |
|------|-----------|---------|
| 1. Adjust giảm | 2 phút | ⬜ PASS / FAIL |
| 2. Adjust tăng không giá | 3 phút | ⬜ PASS / FAIL |
| 3. Adjust tăng có giá | 3 phút | ⬜ PASS / FAIL |
| 4. Form "Sửa" disabled | 1 phút | ⬜ PASS / FAIL |
| 5. UI 4 nút | 30 giây | ⬜ PASS / FAIL |

**Tổng thời gian**: ~10 phút

---

## 🐛 Nếu Có Lỗi

### Lỗi 1: Đơn giá vẫn thay đổi khi giảm
→ Báo lại, cần check frontend logic

### Lỗi 2: Vẫn tạo expense khi không nhập giá
→ Báo lại, cần check backend logic

### Lỗi 3: Form "Sửa" vẫn cho phép thay đổi số lượng
→ Báo lại, cần check frontend UI

### Lỗi 4: Vẫn thấy 6 nút
→ Báo lại, cần check frontend UI

---

## ✅ Nếu Tất Cả PASS

1. Báo "All tests PASS"
2. Sẽ xóa console.log debugging
3. Build và deploy

---

**Ngày**: 2026-02-07  
**Thời gian test**: ~10 phút  
**Quan trọng**: Test 2 và Test 3 (auto-expense logic)
