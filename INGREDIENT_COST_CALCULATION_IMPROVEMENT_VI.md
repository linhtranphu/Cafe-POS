# Cải thiện UI tính toán chi phí Nguyên liệu

## Tổng quan
Cải thiện giao diện quản lý nguyên liệu để hiển thị rõ ràng cách tính chi phí và cách expense được tự động ghi nhận.

## Vấn đề
UI trước đây không rõ ràng về:
- "cost_per_unit" nghĩa là gì (đơn giá)
- Tổng chi phí được tính như thế nào (số lượng × đơn giá)
- Khi nào expense được tự động tạo

## Giải pháp
Cải thiện UI với nhãn rõ ràng, công thức tính toán, và chỉ báo trực quan.

## Các thay đổi

### 1. Form Tạo/Sửa Nguyên liệu

#### Cải thiện nhãn trường
**Trước:**
- "Số lượng" (mơ hồ)
- "Giá/Đơn vị" (không rõ)

**Sau:**
- "Số lượng nhập" với gợi ý đơn vị: `(kg)`
- "Đơn giá" với gợi ý giá: `(VND/kg)`
- Thêm placeholder ví dụ: "VD: 10", "VD: 200000"

#### Thêm hiển thị Tổng chi phí
Thẻ màu xanh dương mới hiển thị:
```
💰 Tổng chi phí: 2,000,000₫
= 10 kg × 200,000₫/kg
```

Xuất hiện khi đã nhập cả số lượng và đơn giá.

#### Cải thiện chỉ báo Auto-Expense
**Trước:**
- Hộp xanh lá đơn giản với tổng số tiền

**Sau:**
- Thẻ xanh lá nổi bật với:
  - Tiêu đề rõ ràng: "Tự động ghi nhận chi phí"
  - Hộp trắng hiển thị số tiền expense
  - Công thức tính: `(10 kg × 200,000₫)`
  - Thông tin danh mục và phương thức thanh toán

### 2. Hiển thị Danh sách Nguyên liệu

#### Cải thiện layout thông tin
**Trước:**
```
Tồn kho: 10 kg
Tối thiểu: 5 kg
Đơn giá: 200,000₫
```

**Sau:**
```
📦 Tồn kho: 10 kg (đậm, lớn hơn)
⚠️ Tối thiểu: 5 kg
💵 Đơn giá: 200,000₫/kg (nổi bật trong hộp xanh)
```

### 3. Modal Điều chỉnh Tồn kho

#### Cải thiện chỉ báo Auto-Expense
Khi điều chỉnh nhập thêm (thêm số lượng):
- Hiển thị chi phí được tính cho số lượng thêm vào
- Công thức: `số_lượng_thêm × đơn_giá`
- Ví dụ: `5 kg × 200,000₫/kg = 1,000,000₫`

### 4. Thêm Computed Property

```javascript
const totalCost = computed(() => {
  const quantity = formData.value.quantity || 0
  const costPerUnit = formData.value.cost_per_unit || 0
  return quantity * costPerUnit
})
```

## Cải thiện Trực quan

### Mã màu
- **Thẻ xanh dương**: Hiển thị tính toán (thông tin)
- **Thẻ xanh lá**: Chỉ báo auto-expense (hành động sẽ được thực hiện)
- **Icons**: 💰 cho chi phí, 📝 cho danh mục, 💵 cho phương thức thanh toán

### Typography
- **Đậm** cho số quan trọng
- **Font lớn hơn** cho tổng số tiền
- **Chữ nhỏ màu xám** cho giải thích

### Layout
- Grid responsive cho các trường form
- Phân cấp trực quan rõ ràng
- Khoảng cách phù hợp giữa các phần

## Ví dụ Luồng Người dùng

### Tạo Nguyên liệu Mới

1. **Nhập thông tin cơ bản:**
   - Tên: "Cà phê hạt"
   - Danh mục: "Nguyên liệu"
   - Đơn vị: "kg"

2. **Nhập số lượng và giá:**
   - Số lượng nhập: `10` kg
   - Đơn giá: `200,000` VND/kg

3. **Xem tính toán:**
   - Thẻ xanh dương hiển thị: "💰 Tổng chi phí: 2,000,000₫"
   - Công thức: "= 10 kg × 200,000₫/kg"

4. **Xem chỉ báo auto-expense:**
   - Thẻ xanh lá xác nhận expense sẽ được tạo
   - Hiển thị số tiền chính xác: 2,000,000₫
   - Hiển thị danh mục và phương thức thanh toán

5. **Gửi:**
   - Nguyên liệu được tạo ✓
   - Expense tự động được tạo ✓

### Điều chỉnh Tồn kho

1. **Chọn nguyên liệu** từ danh sách
2. **Click "📦 Điều chỉnh"**
3. **Chọn loại:** "Nhập thêm"
4. **Nhập số lượng:** `5` kg
5. **Xem auto-expense:**
   - Thẻ xanh lá hiển thị: "Chi phí nhập thêm: 1,000,000₫"
   - Công thức: "= 5 kg × 200,000₫/kg"
6. **Gửi:**
   - Tồn kho được cập nhật ✓
   - Expense được tạo ✓

## Logic Backend (Không thay đổi)

Logic backend vẫn giữ nguyên:
```go
amount := ing.CostPerUnit * quantity
```

Cải thiện UI này chỉ làm cho tính toán hiển thị rõ ràng và dễ hiểu hơn cho người dùng.

## Lợi ích

1. **Rõ ràng**: Người dùng hiểu họ đang nhập gì
2. **Minh bạch**: Tính toán hiển thị trước khi gửi
3. **Tự tin**: Người dùng biết expense sẽ được theo dõi đúng
4. **Giáo dục**: Người dùng mới học hệ thống nhanh chóng
5. **Ngăn lỗi**: Nhãn rõ ràng giảm sai sót khi nhập

## File đã sửa

- `frontend/src/views/IngredientManagementView.vue`

## Checklist kiểm tra

- [ ] Tạo nguyên liệu với số lượng và đơn giá
- [ ] Xác minh tổng chi phí hiển thị đúng
- [ ] Xác minh chỉ báo auto-expense hiển thị số tiền đúng
- [ ] Điều chỉnh nhập thêm và xác minh chỉ báo expense
- [ ] Kiểm tra layout responsive trên mobile
- [ ] Xác minh nhãn đơn vị cập nhật động
- [ ] Test với các đơn vị khác nhau (kg, L, piece, v.v.)
- [ ] Xác minh expense được tạo với số tiền đúng

## Công thức Tham khảo

```
Tổng Chi phí = Số lượng × Đơn giá

Ví dụ:
- 10 kg × 200,000₫/kg = 2,000,000₫
- 5 L × 50,000₫/L = 250,000₫
- 20 cái × 5,000₫/cái = 100,000₫
```

---
**Trạng thái**: ✅ Hoàn thành
**Ngày**: 2026-02-07
**Tác động**: Cải thiện độ rõ ràng UI, không thay đổi backend
