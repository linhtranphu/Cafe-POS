# Điều Chỉnh Nguyên Liệu - Chế Độ Nhập Tổng Giá

## Tóm Tắt
Đã thêm chế độ nhập tổng giá vào modal điều chỉnh tồn kho, giúp nhất quán với flow tạo nguyên liệu mới và dễ sử dụng hơn.

## Thay Đổi

### 1. Thêm Chế Độ Nhập Giá
- **Chế độ mặc định**: Nhập tổng giá (trực quan hơn cho user)
- **Chế độ thay thế**: Nhập đơn giá (cho user nâng cao)
- Nút chuyển đổi giữa 2 chế độ

### 2. Tính Năng Chế Độ Tổng Giá
```javascript
// Khi user nhập tổng giá:
- Nhập: Tổng tiền đã trả (VD: 50.000đ cho 2 kg)
- Hệ thống tính: Đơn giá = Tổng tiền / Số lượng
- Hiển thị: Đơn giá được tính kèm công thức
- Tham khảo: Giá hiện tại để so sánh
```

### 3. Tính Năng Chế Độ Đơn Giá
```javascript
// Khi user nhập đơn giá:
- Nhập: Giá cho 1 đơn vị (VD: 25.000đ/kg)
- Hệ thống dùng: Đơn giá đã nhập
- Hiển thị: Giá hiện tại làm placeholder
- Lưu ý: Để 0 nếu dùng giá hiện tại
```

### 4. Cải Tiến UI
- **Phân biệt rõ ràng**: Viền xanh lá cho tổng giá, viền xám cho đơn giá
- **Tính toán realtime**: Hiển thị đơn giá được tính ngay lập tức
- **Hướng dẫn**: Hiển thị giá hiện tại để tham khảo
- **Mặc định thông minh**: Xóa trường đối diện khi chuyển chế độ

## Luồng Sử Dụng

### Tình Huống 1: User biết tổng tiền (phổ biến nhất)
1. Click "Điều chỉnh" trên nguyên liệu
2. Chọn "Nhập thêm"
3. Nhập số lượng: `2 kg`
4. Nhập tổng giá: `50.000đ`
5. Hệ thống hiển thị: "Đơn giá được tính: 25.000 ₫/kg"
6. Tự động tạo chi phí: 50.000đ

### Tình Huống 2: User biết đơn giá
1. Click chuyển đổi: "Nhập đơn giá"
2. Nhập số lượng: `2 kg`
3. Nhập đơn giá: `25.000đ/kg`
4. Hệ thống tính: Tổng = 50.000đ
5. Tự động tạo chi phí: 50.000đ

### Tình Huống 3: Giá không đổi
1. Chỉ nhập số lượng
2. Để giá = 0 (hoặc không nhập)
3. Hệ thống dùng giá hiện tại của nguyên liệu
4. Chi phí tự động dùng giá hiện tại

## Triển Khai Kỹ Thuật

### Biến State Mới
```javascript
const adjustPriceMode = ref('total') // 'total' hoặc 'unit'
const adjustTotalPrice = ref(0)
```

### Hàm Mới
```javascript
toggleAdjustPriceMode() // Chuyển đổi chế độ
calculateAdjustUnitPrice() // Tổng giá → Đơn giá
```

### Tính Toán Reactive
```javascript
// Khi số lượng hoặc tổng giá thay đổi:
@input="adjustPriceMode === 'total' && calculateAdjustUnitPrice()"

// Công thức:
đơn_giá = tổng_giá / số_lượng
```

## Lợi Ích

### 1. Nhất Quán
- UX giống modal tạo nguyên liệu mới
- User học 1 lần, dùng mọi nơi

### 2. Dễ Sử Dụng
- Hầu hết user biết tổng tiền đã trả, không phải đơn giá
- Không cần tính toán thủ công
- Giảm lỗi nhập liệu

### 3. Linh Hoạt
- Power user vẫn có thể nhập đơn giá trực tiếp
- Chuyển đổi dễ dàng giữa các chế độ
- Giữ nguyên chức năng cũ

### 4. Minh Bạch
- Hiển thị công thức tính toán
- Hiển thị giá hiện tại để tham khảo
- Phản hồi trực quan rõ ràng

## Tích Hợp Tự Động Ghi Chi Phí

Tính toán chi phí tự động hoạt động với cả 2 chế độ:

```javascript
giá_hiệu_lực = adjustData.cost_per_unit > 0 
  ? adjustData.cost_per_unit 
  : currentIngredient.cost_per_unit

số_tiền_chi_phí = số_lượng × giá_hiệu_lực
```

**Kết quả**: Dù user nhập tổng giá hay đơn giá, chi phí đều được ghi đúng.

## Checklist Kiểm Tra

- [ ] Chế độ tổng giá: Nhập tổng, kiểm tra đơn giá được tính
- [ ] Chế độ đơn giá: Nhập đơn giá, kiểm tra chi phí được tính
- [ ] Chuyển đổi chế độ: Kiểm tra các trường xóa/tính lại
- [ ] Để giá = 0: Kiểm tra dùng giá hiện tại
- [ ] Thay đổi số lượng: Kiểm tra đơn giá tính lại ở chế độ tổng giá
- [ ] Chi phí tự động: Kiểm tra số tiền ghi đúng
- [ ] Cập nhật giá: Kiểm tra cost_per_unit của nguyên liệu được cập nhật nếu nhập giá mới

## File Đã Sửa

1. `frontend/src/views/IngredientManagementView.vue`
   - Thêm state `adjustPriceMode` và `adjustTotalPrice`
   - Thêm hàm `toggleAdjustPriceMode()`
   - Thêm hàm `calculateAdjustUnitPrice()`
   - Cập nhật UI modal điều chỉnh với nút chuyển đổi chế độ
   - Cập nhật input số lượng để trigger tính toán lại

## Bước Tiếp Theo

Cân nhắc áp dụng pattern tương tự cho:
- Điều chỉnh tồn kho cơ sở vật chất (nếu có)
- Các tính năng quản lý kho khác
- Form nhập chi phí (ngược lại: nhập tổng, tính chi tiết)
