# Quản Lý Nguyên Liệu - Ngăn Chặn Submit Trùng

## Vấn Đề
User có thể vô tình tạo bản ghi trùng do:
- Chạm 2 lần vào nút submit trên mobile
- Click nhiều lần khi đợi phản hồi
- Mạng chậm gây nhầm lẫn về trạng thái submit

## Giải Pháp Đã Triển Khai

### 1. Quản Lý Trạng Thái Loading
Thêm 3 cờ trạng thái loading:
```javascript
const isSubmitting = ref(false)  // Cho thao tác tạo/sửa
const isDeleting = ref(false)    // Cho thao tác xóa
const isAdjusting = ref(false)   // Cho điều chỉnh tồn kho
```

### 2. Vô Hiệu Hóa Nút Khi Đang Xử Lý
Tất cả nút action bị vô hiệu hóa khi đang xử lý:
```vue
<button 
  :disabled="isSubmitting"
  class="disabled:opacity-50 disabled:cursor-not-allowed">
```

### 3. Phản Hồi Trực Quan
**Spinner Loading:**
```vue
<span v-if="isSubmitting" 
  class="inline-block w-4 h-4 border-2 border-white border-t-transparent rounded-full animate-spin">
</span>
```

**Text Nút Động:**
- Bình thường: "Thêm mới" / "Cập nhật" / "Xác nhận"
- Đang xử lý: "Đang lưu..." / "Đang xử lý..."

### 4. Kiểm Tra Dữ Liệu
Thêm validation toàn diện trước khi submit:

**Validation Tạo/Sửa:**
```javascript
- Trường bắt buộc: tên, danh mục, đơn vị
- Giá trị không âm: số lượng, tối thiểu, đơn giá
- Thông báo lỗi cụ thể cho user
```

**Validation Điều Chỉnh:**
```javascript
- Số lượng phải > 0
- Lý do không được để trống
- Số lượng mới không được âm
- Thông báo lỗi cụ thể cho user
```

### 5. Pattern Early Return
```javascript
if (isSubmitting.value) return  // Ngăn submit trùng
```

### 6. Reset Trạng Thái
Trạng thái loading được reset khi:
- Modal đóng
- Thao tác hoàn thành (thành công hoặc lỗi)
- Modal mở lại

## Thay Đổi UI

### Footer Modal Tạo/Sửa
```vue
<!-- Trước -->
<button @click="saveIngredient">Thêm mới</button>

<!-- Sau -->
<button 
  @click="saveIngredient" 
  :disabled="isSubmitting"
  class="disabled:opacity-50 disabled:cursor-not-allowed flex items-center justify-center gap-2">
  <span v-if="isSubmitting" class="spinner"></span>
  <span>{{ isSubmitting ? 'Đang lưu...' : 'Thêm mới' }}</span>
</button>
```

### Footer Modal Điều Chỉnh
```vue
<button 
  @click="adjustStock" 
  :disabled="isAdjusting"
  class="disabled:opacity-50 disabled:cursor-not-allowed flex items-center justify-center gap-2">
  <span v-if="isAdjusting" class="spinner"></span>
  <span>{{ isAdjusting ? 'Đang xử lý...' : 'Xác nhận' }}</span>
</button>
```

### Nút Action Trên Card
Tất cả nút action bị vô hiệu hóa khi đang xử lý:
```vue
<button @click="openAdjustModal(ingredient)" :disabled="isAdjusting">
  📦 Điều chỉnh
</button>
<button @click="openEditModal(ingredient)" :disabled="isSubmitting">
  ✏️ Sửa
</button>
<button @click="deleteIngredient(ingredient)" :disabled="isDeleting">
  🗑️ Xóa
</button>
```

## Luồng Trải Nghiệm User

### Tình Huống 1: Tạo Nguyên Liệu
1. User điền form và click "Thêm mới"
2. Nút ngay lập tức hiển thị spinner + "Đang lưu..."
3. Nút bị vô hiệu hóa (mờ đi, không click được)
4. Nút Hủy cũng bị vô hiệu hóa
5. Tất cả nút trên card nguyên liệu bị vô hiệu hóa
6. Sau khi thành công: Modal đóng, trạng thái reset
7. Sau khi lỗi: Hiển thị thông báo, nút được kích hoạt lại

### Tình Huống 2: Điều Chỉnh Tồn Kho
1. User nhập điều chỉnh và click "Xác nhận"
2. Nút hiển thị spinner + "Đang xử lý..."
3. Tất cả nút bị vô hiệu hóa
4. Sau khi thành công: Modal đóng, danh sách refresh
5. Sau khi lỗi: Hiển thị thông báo, có thể thử lại

### Tình Huống 3: Xóa Nguyên Liệu
1. User click "Xóa", xác nhận dialog
2. Nút xóa bị vô hiệu hóa ngay lập tức
3. Các nút action khác vẫn hoạt động
4. Sau khi thành công: Item bị xóa khỏi danh sách
5. Sau khi lỗi: Hiển thị thông báo, nút được kích hoạt lại

## Triển Khai Kỹ Thuật

### Guard Trạng Thái Loading
```javascript
const saveIngredient = async () => {
  if (isSubmitting.value) return  // Guard clause
  
  // Validation
  if (!formData.value.name) {
    alert('Vui lòng điền đầy đủ thông tin')
    return
  }
  
  isSubmitting.value = true
  try {
    await ingredientStore.createIngredient(formData.value)
    closeModal()
  } catch (error) {
    alert('Có lỗi xảy ra')
  } finally {
    isSubmitting.value = false  // Luôn reset
  }
}
```

### CSS Cho Trạng Thái Disabled
```css
.disabled\:opacity-50:disabled {
  opacity: 0.5;
}

.disabled\:cursor-not-allowed:disabled {
  cursor: not-allowed;
}
```

### Animation Spinner
```css
@keyframes spin {
  to { transform: rotate(360deg); }
}

.animate-spin {
  animation: spin 1s linear infinite;
}
```

## Lợi Ích

### 1. Ngăn Chặn Bản Ghi Trùng
- Không còn submit trùng vô tình
- Mạng chậm không gây nhầm lẫn
- Xử lý vấn đề tap trên mobile

### 2. Phản Hồi Rõ Ràng Cho User
- Spinner hiển thị thao tác đang xử lý
- Text nút thay đổi để chỉ ra hành động
- Trạng thái disabled ngăn tương tác

### 3. Xử Lý Lỗi Tốt Hơn
- Validation trước khi submit
- Thông báo lỗi cụ thể
- Trạng thái luôn reset đúng cách

### 4. UX Chuyên Nghiệp
- Phù hợp với tiêu chuẩn app hiện đại
- Giảm frustration cho user
- Xây dựng niềm tin vào hệ thống

## Checklist Kiểm Tra

- [ ] Tạo nguyên liệu: Click submit nhiều lần nhanh
- [ ] Sửa nguyên liệu: Chạm 2 lần nút lưu
- [ ] Điều chỉnh tồn kho: Click xác nhận 2 lần nhanh
- [ ] Xóa nguyên liệu: Click xóa nhiều lần
- [ ] Mạng chậm: Kiểm tra spinner hiển thị khi delay
- [ ] Trường hợp lỗi: Kiểm tra nút kích hoạt lại sau lỗi
- [ ] Đóng modal: Kiểm tra trạng thái reset khi đóng
- [ ] Nhiều nguyên liệu: Kiểm tra chỉ 1 thao tác tại 1 thời điểm

## Xử Lý Edge Cases

1. **Network timeout**: Finally block đảm bảo reset trạng thái
2. **Đóng modal khi đang xử lý**: Trạng thái reset khi đóng
3. **Nhiều modal**: Mỗi modal có trạng thái loading độc lập
4. **Click nhanh**: Guard clause ngăn thực thi
5. **Lỗi validation**: Trạng thái không set nếu validation fail

## Cải Tiến Tương Lai

Cân nhắc thêm:
- Toast notifications thay vì alerts
- Animation thành công
- Optimistic UI updates
- Chức năng Undo
- Batch operations với progress bar
- Chỉ báo trạng thái mạng

## File Đã Sửa

1. `frontend/src/views/IngredientManagementView.vue`
   - Thêm loading state refs
   - Cập nhật tất cả action functions với guards
   - Thêm logic validation
   - Cập nhật button templates với disabled states
   - Thêm spinner component
   - Thêm CSS animations
