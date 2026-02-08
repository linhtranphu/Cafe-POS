# Lịch Sử Nguyên Liệu - Hiển Thị Giá Cải Tiến

## Tóm Tắt
Cải thiện modal lịch sử tồn kho để hiển thị rõ ràng giá nhập và thông tin chi phí, giúp dễ dàng theo dõi biến động giá nguyên liệu theo thời gian.

## Cải Tiến Chính

### 1. Header Nâng Cấp
**Trước:**
- Tiêu đề đơn giản: "Lịch sử tồn kho"
- Chỉ có tên nguyên liệu

**Sau:**
- Tiêu đề rõ ràng: "Lịch sử nhập hàng"
- Tên nguyên liệu được nhấn mạnh
- Hiển thị giá hiện tại để tham khảo
```vue
<p class="text-xs text-gray-500 text-center">
  Giá hiện tại: <span class="font-semibold text-green-600">
    {{ formatCurrency(currentIngredient?.cost_per_unit) }}/{{ currentIngredient?.unit }}
  </span>
</p>
```

### 2. Phân Biệt Trực Quan Theo Loại Giao Dịch
**Card có màu sắc:**
- Viền/nền xanh lá: Nhập kho (mua hàng)
- Viền/nền đỏ: Xuất kho (sử dụng/hao hụt)

```vue
:class="record.quantity > 0 ? 'border-green-200 bg-green-50' : 'border-red-200 bg-red-50'"
```

### 3. Phần Thông Tin Giá Nổi Bật
Với giao dịch mua hàng (quantity > 0), hiển thị card giá được highlight:

**Tính năng:**
- Nền gradient xanh lá (from-green-100 to-emerald-100)
- Viền xanh để nhấn mạnh
- Icon tiền (💰) ở header
- Label "THÔNG TIN GIÁ" in hoa

**Thông tin hiển thị:**
1. **Đơn giá**: Giá mỗi đơn vị cho lần nhập này
2. **Tổng chi phí**: Tổng số tiền đã trả
3. **Công thức tính**: Hiển thị phép tính (số lượng × đơn giá)
4. **So sánh giá**: So với giá hiện tại (↑ hoặc ↓)

```vue
<div class="bg-gradient-to-br from-green-100 to-emerald-100 border-2 border-green-300 rounded-xl p-3">
  <!-- Đơn giá -->
  <div class="flex items-center justify-between bg-white rounded-lg px-3 py-2">
    <span class="text-xs text-gray-600">Đơn giá lần này:</span>
    <span class="font-bold text-green-700">
      {{ formatCurrency(record.cost_per_unit) }}/{{ unit }}
    </span>
  </div>
  
  <!-- Tổng chi phí -->
  <div class="flex items-center justify-between bg-white rounded-lg px-3 py-2">
    <span class="text-xs text-gray-600">Tổng chi phí:</span>
    <span class="font-bold text-green-800 text-base">
      {{ formatCurrency(record.total_cost) }}
    </span>
  </div>
  
  <!-- Công thức -->
  <div class="text-xs text-green-700 text-center bg-green-50 rounded-lg py-1">
    = {{ quantity }} {{ unit }} × {{ formatCurrency(cost_per_unit) }}
  </div>
  
  <!-- So sánh giá -->
  <div class="flex items-center justify-center gap-2 text-xs pt-1">
    <span class="text-gray-600">So với giá hiện tại:</span>
    <span class="text-red-600 font-semibold">↑ 5.000 ₫</span>
  </div>
</div>
```

### 4. Header Giao Dịch Rõ Ràng
**Layout cải tiến:**
- Badge loại giao dịch bên trái
- Số lượng lớn, đậm với dấu +/-
- Màu sắc (xanh cho +, đỏ cho -)
- Số lượng Trước → Sau bên phải

```vue
<div class="flex items-center justify-between mb-3">
  <div class="flex items-center gap-2">
    <span class="badge">Nhập thêm</span>
    <span class="text-lg font-bold text-green-700">+5 kg</span>
  </div>
  <div class="text-right">
    <p class="text-xs text-gray-500">Trước → Sau</p>
    <p class="text-sm font-bold">10 → 15</p>
  </div>
</div>
```

### 5. Hiển Thị Lý Do
Chuyển sang card trắng riêng để dễ đọc hơn:
```vue
<div class="mb-3 bg-white rounded-lg p-3">
  <p class="text-sm text-gray-700">
    <span class="font-semibold">Lý do:</span> {{ record.reason }}
  </p>
</div>
```

### 6. Chỉ Báo Không Có Giá Cho Xuất Kho
Với giao dịch xuất kho, hiển thị thông báo rõ ràng:
```vue
<div class="bg-red-100 border border-red-200 rounded-lg p-2">
  <p class="text-xs text-red-700 text-center">
    ⚠️ Xuất kho - Không có thông tin giá
  </p>
</div>
```

### 7. Footer Metadata Nâng Cấp
```vue
<div class="flex items-center justify-between text-xs text-gray-500 pt-2 border-t">
  <div class="flex items-center gap-1">
    <span>👤</span>
    <span class="font-medium">{{ username }}</span>
  </div>
  <div class="flex items-center gap-1">
    <span>🕐</span>
    <span>{{ formatDateTime(created_at) }}</span>
  </div>
</div>
```

## Thứ Tự Ưu Tiên Trực Quan

### Ưu tiên 1: Loại & Số Lượng Giao Dịch
- Số lượng lớn, đậm có màu
- Badge rõ ràng cho loại giao dịch

### Ưu tiên 2: Thông Tin Giá (cho mua hàng)
- Card highlight với nền gradient
- Đơn giá và tổng chi phí hiển thị nổi bật
- Phân tích công thức để minh bạch

### Ưu tiên 3: Ngữ Cảnh
- Lý do trong card riêng
- Số lượng Trước/Sau
- User và timestamp

## Trường Hợp Sử Dụng

### Trường Hợp 1: Theo Dõi Thay Đổi Giá
**Tình huống:** Quản lý muốn xem giá nguyên liệu có tăng không

**Giải pháp:**
1. Mở lịch sử nguyên liệu
2. Xem giá hiện tại ở đầu
3. Cuộn qua các bản ghi mua hàng
4. Mỗi bản ghi hiển thị:
   - Giá tại thời điểm mua
   - So sánh với giá hiện tại (↑ hoặc ↓)
   - Chỉ báo trực quan (đỏ tăng, xanh giảm)

### Trường Hợp 2: Xác Minh Chi Phí Mua Hàng
**Tình huống:** Quản lý muốn xác minh đơn hàng gần đây

**Giải pháp:**
1. Tìm bản ghi mua hàng (card xanh)
2. Xem phần giá được highlight:
   - Đơn giá đã trả
   - Tổng chi phí
   - Phân tích công thức
3. So sánh với hóa đơn

### Trường Hợp 3: Phân Tích Xu Hướng Chi Phí
**Tình huống:** Quản lý muốn hiểu xu hướng chi phí

**Giải pháp:**
1. Cuộn qua lịch sử theo thời gian
2. Mỗi lần mua hiển thị:
   - Ngày giờ
   - Số lượng mua
   - Giá đã trả
   - So sánh với giá hiện tại
3. Nhận diện patterns (thay đổi theo mùa, đổi nhà cung cấp, v.v.)

## Tối Ưu Mobile

### Thân Thiện Với Touch
- Card lớn với khoảng cách tốt
- Phân tách trực quan rõ ràng
- Dễ cuộn

### Mật Độ Thông Tin
- Cân bằng: Không quá đông, không quá thưa
- Thông tin quan trọng (giá) nổi bật
- Thông tin phụ (metadata) nhỏ hơn nhưng đọc được

### Mã Màu
- Xanh lá: Tích cực (mua hàng, nhập kho)
- Đỏ: Tiêu cực (sử dụng, xuất kho)
- Trực quan và nhất quán

## Lợi Ích

### 1. Minh Bạch Giá
- Hiển thị rõ giá lịch sử
- Dễ theo dõi thay đổi giá
- So sánh với giá hiện tại

### 2. Xác Minh Chi Phí
- Tổng chi phí hiển thị rõ
- Công thức tính được cung cấp
- Dễ xác minh với hóa đơn

### 3. Ra Quyết Định Tốt Hơn
- Dữ liệu giá lịch sử hiển thị
- Xu hướng dễ nhận diện
- Quyết định mua hàng có thông tin

### 4. Giao Diện Chuyên Nghiệp
- Thiết kế sạch, hiện đại
- Mã màu để rõ ràng
- Hiệu ứng gradient để nhấn mạnh

## Chi Tiết Kỹ Thuật

### Render Có Điều Kiện
```javascript
// Hiển thị thông tin giá chỉ cho mua hàng có dữ liệu giá
v-if="record.cost_per_unit > 0 && record.quantity > 0"

// Hiển thị chỉ báo không có giá cho xuất kho
v-else-if="record.quantity < 0"
```

### Logic So Sánh Giá
```javascript
// Tính chênh lệch
const priceDiff = Math.abs(record.cost_per_unit - currentIngredient?.cost_per_unit)

// Xác định hướng
const isIncrease = record.cost_per_unit > currentIngredient?.cost_per_unit

// Hiển thị mũi tên
{{ isIncrease ? '↑' : '↓' }} {{ formatCurrency(priceDiff) }}
```

### Class Màu
```javascript
// Màu động dựa trên thay đổi giá
:class="record.cost_per_unit > currentIngredient?.cost_per_unit 
  ? 'text-red-600'    // Giá tăng
  : 'text-green-600'" // Giá giảm
```

## Cải Tiến Tương Lai

Cân nhắc thêm:
1. **Biểu đồ giá**: Đồ thị trực quan thay đổi giá theo thời gian
2. **Giá trung bình**: Tính và hiển thị giá mua trung bình
3. **Cảnh báo giá**: Thông báo khi giá thay đổi đáng kể
4. **Export**: Xuất lịch sử ra CSV/Excel
5. **Bộ lọc**: Lọc theo khoảng thời gian, khoảng giá, user
6. **Thống kê**: Giá min/max/trung bình, tổng chi tiêu

## File Đã Sửa

1. `frontend/src/views/IngredientManagementView.vue`
   - Nâng cấp UI modal lịch sử tồn kho
   - Thêm phần thông tin giá
   - Thêm mã màu cho loại giao dịch
   - Thêm logic so sánh giá
   - Cải thiện layout và spacing
