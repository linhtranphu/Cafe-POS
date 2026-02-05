# Cash Handover Discrepancy - Quick Guide 🚀

## 🎯 Tóm tắt nhanh

Khi cashier xác nhận bàn giao và nhập số tiền thực nhận **khác** với số tiền khai báo, hệ thống sẽ:

1. ⚠️ **Hiển thị cảnh báo** màu đỏ (thiếu) hoặc xanh (thừa)
2. 📝 **Yêu cầu nhập lý do** chênh lệch
3. 👤 **Yêu cầu chọn trách nhiệm** (Waiter/Cashier/Customer/System/Unknown)
4. 🔔 **Thông báo nếu cần manager approval** (chênh lệch > 100k)

## 📱 Hướng dẫn sử dụng

### Bước 1: Vào màn hình bàn giao
```
http://localhost:5173/#/cashier/handovers
```

### Bước 2: Click "Xác nhận" trên handover request

### Bước 3: Nhập số tiền thực nhận

**Ví dụ:**
- Waiter khai báo: **100,000₫**
- Bạn đếm thực tế: **95,000₫**
- Nhập: `95000`

### Bước 4: Cảnh báo tự động hiển thị

```
⚠️ Thiếu tiền
Chênh lệch: 5,000₫
```

### Bước 5: Nhập thông tin bắt buộc

**Lý do chênh lệch:** (textarea)
```
Khách trả thiếu 5k, waiter không để ý
```

**Trách nhiệm:** (dropdown)
```
[x] Waiter
[ ] Cashier (Tôi)
[ ] Khách hàng
[ ] Hệ thống
[ ] Chưa rõ
```

### Bước 6: Click "Xác nhận"

## 🎨 Màu sắc cảnh báo

| Tình huống | Màu | Icon | Ý nghĩa |
|------------|-----|------|---------|
| Thiếu tiền | 🔴 Đỏ | 📉 | SHORTAGE - Actual < Declared |
| Thừa tiền | 🟢 Xanh | 📈 | OVERAGE - Actual > Declared |
| Cần approval | 🟠 Cam | 🔔 | Chênh lệch > 100k |

## 📊 Các trường hợp

### ✅ Không chênh lệch
```
Khai báo: 100,000₫
Thực nhận: 100,000₫
→ Không có cảnh báo
→ Xác nhận bình thường
```

### 🔴 Thiếu tiền nhỏ
```
Khai báo: 100,000₫
Thực nhận: 95,000₫
Chênh lệch: -5,000₫
→ Cảnh báo đỏ
→ Nhập lý do + trách nhiệm
→ Xác nhận ngay (không cần manager)
```

### 🔴 Thiếu tiền lớn
```
Khai báo: 200,000₫
Thực nhận: 50,000₫
Chênh lệch: -150,000₫
→ Cảnh báo đỏ + badge cam
→ Nhập lý do + trách nhiệm
→ Chờ manager phê duyệt
```

### 🟢 Thừa tiền nhỏ
```
Khai báo: 100,000₫
Thực nhận: 110,000₫
Chênh lệch: +10,000₫
→ Cảnh báo xanh
→ Nhập lý do + trách nhiệm
→ Xác nhận ngay (không cần manager)
```

### 🟢 Thừa tiền lớn
```
Khai báo: 100,000₫
Thực nhận: 250,000₫
Chênh lệch: +150,000₫
→ Cảnh báo xanh + badge cam
→ Nhập lý do + trách nhiệm
→ Chờ manager phê duyệt
```

## 💡 Tips

### Lý do thường gặp

**Thiếu tiền:**
- "Khách trả thiếu, waiter không để ý"
- "Mất tiền trên đường bàn giao"
- "Tính sai tiền thối"
- "Khách chưa trả đủ"

**Thừa tiền:**
- "Khách tip thêm"
- "Tính nhầm, thu thừa"
- "Khách trả nhầm mệnh giá"
- "Tìm thấy tiền thừa trong túi"

### Chọn trách nhiệm

| Trách nhiệm | Khi nào chọn |
|-------------|--------------|
| **WAITER** | Waiter sai sót (đếm sai, mất tiền, thu thiếu) |
| **CASHIER** | Bạn đếm sai, nhập sai số liệu |
| **CUSTOMER** | Khách trả thiếu/thừa, tip |
| **SYSTEM** | Lỗi hệ thống, tính toán sai |
| **UNKNOWN** | Chưa rõ nguyên nhân |

## ⚠️ Lưu ý quan trọng

1. **Bắt buộc nhập đầy đủ** - Không thể submit nếu thiếu lý do hoặc trách nhiệm
2. **Threshold 100k** - Chênh lệch > 100k cần manager phê duyệt
3. **Không thể sửa** - Sau khi xác nhận không thể thay đổi
4. **Audit trail** - Mọi thông tin được lưu lại để kiểm tra

## 🔍 Kiểm tra lại

Trước khi xác nhận, hãy:
- ✅ Đếm lại tiền thật cẩn thận
- ✅ Kiểm tra có tiền giả không
- ✅ Xác nhận số tiền với waiter
- ✅ Ghi chú rõ ràng lý do chênh lệch

## 📞 Cần hỗ trợ?

Nếu gặp vấn đề:
1. Gọi manager ngay lập tức
2. Không xác nhận nếu chưa rõ
3. Giữ nguyên hiện trường
4. Ghi chép chi tiết

## 🎓 Training Checklist

- [ ] Hiểu cách tính discrepancy
- [ ] Biết khi nào cần manager approval
- [ ] Biết cách nhập lý do rõ ràng
- [ ] Biết cách chọn trách nhiệm đúng
- [ ] Thực hành với các scenarios khác nhau
- [ ] Biết cách xử lý khi có vấn đề

## 📚 Tài liệu liên quan

- [CASH_HANDOVER_DISCREPANCY_FIX.md](./CASH_HANDOVER_DISCREPANCY_FIX.md) - Chi tiết kỹ thuật
- [CASH_HANDOVER_UI_GUIDE.md](./CASH_HANDOVER_UI_GUIDE.md) - Hướng dẫn UI đầy đủ
- [CASH_HANDOVER_TROUBLESHOOTING.md](./CASH_HANDOVER_TROUBLESHOOTING.md) - Xử lý sự cố
