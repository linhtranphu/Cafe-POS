# 🧪 Hướng Dẫn Test Chức Năng Điều Chỉnh Nguyên Liệu

## 📍 URL Test
```
http://localhost:5173/#/ingredients
```

---

## ✅ Test 1: Điều Chỉnh GIẢM - Giá KHÔNG ĐỔI

### Bước 1: Chuẩn bị
- Mở trang Nguyên liệu
- Tìm nguyên liệu có tồn kho (VD: Đường 10kg @ 25,000đ/kg)
- **Ghi nhớ đơn giá hiện tại**: 25,000đ/kg

### Bước 2: Thực hiện điều chỉnh
1. Click nút **"📦 Điều chỉnh"**
2. Chọn loại: **"Điều chỉnh"**
3. Nhập số lượng mới: **8** (giảm từ 10 xuống 8)
4. **KHÔNG nhập giá** (để 0 hoặc trống)
5. Nhập lý do: "Test giảm số lượng"
6. Click **"Xác nhận"**

### Bước 3: Kiểm tra kết quả
✅ **Tồn kho**: 8kg (đã giảm)  
✅ **Đơn giá**: 25,000đ/kg (KHÔNG ĐỔI)  
✅ **Console log** (F12 → Console):
```
=== ADJUST DEBUG ===
Current quantity: 10
New quantity: 8
Is increase? false
cost_per_unit before: 0
Data to send: { new_quantity: 8, cost_per_unit: 0, reason: "..." }
===================
```

### ❌ Nếu sai:
- Đơn giá thay đổi → BUG chưa fix
- Console log không hiển thị → Kiểm tra lại code

---

## ✅ Test 2: Điều Chỉnh TĂNG Không Nhập Giá - Giá KHÔNG ĐỔI

### Bước 1: Chuẩn bị
- Nguyên liệu: Đường 8kg @ 25,000đ/kg
- **Ghi nhớ đơn giá**: 25,000đ/kg

### Bước 2: Thực hiện
1. Click **"📦 Điều chỉnh"**
2. Chọn: **"Điều chỉnh"**
3. Nhập số lượng: **12** (tăng từ 8 lên 12)
4. **KHÔNG nhập giá** (để 0)
5. Lý do: "Tìm thấy thêm hàng"
6. Click **"Xác nhận"**

### Bước 3: Kiểm tra
✅ **Tồn kho**: 12kg (đã tăng)  
✅ **Đơn giá**: 25,000đ/kg (KHÔNG ĐỔI)  
✅ **Console log**:
```
Is increase? true
cost_per_unit before: 0
Data to send: { new_quantity: 12, cost_per_unit: 0, ... }
```

---

## ✅ Test 3: Điều Chỉnh TĂNG Có Nhập Giá - Tính Weighted Average

### Bước 1: Chuẩn bị
- Nguyên liệu: Đường 8kg @ 25,000đ/kg
- **Tính toán trước**:
  - Giá trị cũ: 8kg × 25,000 = 200,000đ
  - Mua thêm: 4kg × 30,000 = 120,000đ
  - Tổng: 12kg = 320,000đ
  - **Đơn giá mới**: 320,000 ÷ 12 = 26,667đ/kg

### Bước 2: Thực hiện
1. Click **"📦 Điều chỉnh"**
2. Chọn: **"Điều chỉnh"**
3. Nhập số lượng: **12**
4. **Nhập giá mới**: **30,000**
5. Lý do: "Mua thêm với giá mới"
6. Click **"Xác nhận"**

### Bước 3: Kiểm tra
✅ **Tồn kho**: 12kg  
✅ **Đơn giá**: ~26,667đ/kg (weighted average)  
✅ **Console log**:
```
Is increase? true
cost_per_unit before: 30000
Data to send: { new_quantity: 12, cost_per_unit: 30000, ... }
```

---

## ✅ Test 4: Form "Sửa" - KHÔNG Thể Thay Đổi Tồn Kho

### Bước 1: Mở form sửa
1. Click nút **"✏️ Sửa"** (nút màu xanh)
2. Form "Cập nhật nguyên liệu" mở ra

### Bước 2: Kiểm tra UI
✅ **Field "Số lượng nhập"**:
- Màu xám (disabled)
- Không thể click hoặc nhập
- Có text warning: "⚠️ Dùng 'Điều chỉnh' để thay đổi tồn kho"

✅ **Section giá**:
- KHÔNG có section "💰 Thông tin giá"
- Thay vào đó có "📊 Thông tin hiện tại (chỉ xem)"
- Hiển thị tồn kho và giá ở chế độ read-only
- Có warning: "⚠️ Để thay đổi tồn kho hoặc giá, vui lòng sử dụng chức năng 'Điều chỉnh'"

### Bước 3: Test thay đổi
✅ **Có thể sửa**:
- Tên nguyên liệu
- Danh mục
- Đơn vị
- Mức tối thiểu
- Nhà cung cấp
- Ghi chú

❌ **KHÔNG thể sửa**:
- Số lượng tồn kho
- Đơn giá

---

## ✅ Test 5: UI Chỉ Có 4 Nút

### Kiểm tra danh sách nguyên liệu

Mỗi nguyên liệu có **4 nút** trên **1 hàng**:

```
┌─────────────────────────────────────┐
│ Đường                               │
│ Nguyên liệu khô                     │
│ Tồn: 10kg | Giá: 25,000đ/kg        │
├─────────────────────────────────────┤
│ [📦 Điều chỉnh] [📊 Lịch sử]       │
│ [✏️ Sửa]        [🗑️ Xóa]           │
└─────────────────────────────────────┘
```

✅ **Có 4 nút**:
1. 📦 Điều chỉnh (màu vàng)
2. 📊 Lịch sử (màu tím)
3. ✏️ Sửa (màu xanh)
4. 🗑️ Xóa (màu đỏ)

❌ **KHÔNG có**:
- ~~Nhập nhanh~~
- ~~Xuất nhanh~~

---

## 🐛 Các Lỗi Có Thể Gặp

### Lỗi 1: Đơn giá vẫn bị thay đổi khi giảm
**Nguyên nhân**: Logic frontend chưa đúng  
**Kiểm tra**: Console log có hiển thị `cost_per_unit: 0` không?  
**Fix**: Xem lại code line 1489-1507

### Lỗi 2: Console log không hiển thị
**Nguyên nhân**: Chưa mở Developer Tools  
**Fix**: Nhấn F12 → Tab Console

### Lỗi 3: Form "Sửa" vẫn cho phép thay đổi số lượng
**Nguyên nhân**: `:disabled="isEditing"` chưa được apply  
**Fix**: Kiểm tra lại code line ~310

### Lỗi 4: Vẫn thấy 6 nút thay vì 4 nút
**Nguyên nhân**: Code UI chưa được update  
**Fix**: Kiểm tra lại code line ~122-148

---

## 📊 Bảng Tổng Hợp Kết Quả Test

| Test Case | Tồn kho trước | Tồn kho sau | Giá trước | Giá sau | Kết quả |
|-----------|---------------|-------------|-----------|---------|---------|
| Giảm không nhập giá | 10kg | 8kg | 25,000 | 25,000 | ✅ PASS |
| Tăng không nhập giá | 8kg | 12kg | 25,000 | 25,000 | ✅ PASS |
| Tăng có nhập giá | 8kg | 12kg | 25,000 | 26,667 | ✅ PASS |
| Form sửa disabled | - | - | - | - | ✅ PASS |
| UI 4 nút | - | - | - | - | ✅ PASS |

---

## 🎯 Kết Luận

Sau khi test xong **5 test cases** trên:

### ✅ Nếu tất cả PASS:
1. Xóa console.log debugging (line 1493-1503)
2. Build frontend: `cd frontend && npm run build`
3. Deploy lên production

### ❌ Nếu có test FAIL:
1. Ghi lại test case nào fail
2. Chụp screenshot console log
3. Báo lại để fix

---

**Ngày tạo**: 2026-02-07  
**File liên quan**: `frontend/src/views/IngredientManagementView.vue`
