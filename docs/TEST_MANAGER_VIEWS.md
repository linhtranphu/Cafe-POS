# Test Manager Views

## Checklist để test Facility và Ingredient Management

### 1. Login as Manager
- [ ] Login với tài khoản manager
- [ ] Verify navigation hiển thị đúng

### 2. Test Facility Management

#### Access
- [ ] Click vào "Cơ sở vật chất" trong navigation
- [ ] URL: `/facilities`
- [ ] View load thành công

#### Dashboard Stats
- [ ] Hiển thị 4 stats cards:
  - Tổng Thiết Bị
  - Hoạt Động (màu xanh)
  - Bảo Trì (màu vàng)
  - Hỏng Hóc (màu đỏ)

#### Facility List
- [ ] Bảng hiển thị danh sách thiết bị
- [ ] Columns: Tên, Loại, Vị Trí, Trạng Thái, Bảo Trì Tiếp, Thao Tác
- [ ] Search box hoạt động

#### Create Facility
- [ ] Click "➕ Thêm Thiết Bị"
- [ ] Modal mở ra
- [ ] Fill form:
  - Tên Thiết Bị (required)
  - Loại (required)
  - Model
  - Vị Trí (required)
  - Trạng Thái (required)
  - Ngày Mua
  - Chu Kỳ Bảo Trì
  - Ghi Chú
- [ ] Click "Thêm Mới"
- [ ] Thiết bị mới xuất hiện trong danh sách

#### Edit Facility
- [ ] Click "✏️ Sửa" trên một thiết bị
- [ ] Modal mở với data đã điền
- [ ] Sửa thông tin
- [ ] Click "Cập Nhật"
- [ ] Thông tin được cập nhật

#### Delete Facility
- [ ] Click "🗑️ Xóa"
- [ ] Confirmation dialog xuất hiện
- [ ] Confirm xóa
- [ ] Thiết bị bị xóa khỏi danh sách

#### Maintenance Schedule
- [ ] Click "📅 Lịch Bảo Trì"
- [ ] Modal hiển thị lịch bảo trì
- [ ] Hiển thị thiết bị cần bảo trì
- [ ] Cảnh báo quá hạn (màu đỏ)

#### Issue Reports
- [ ] Click "⚠️ Báo Cáo Sự Cố"
- [ ] Modal hiển thị danh sách sự cố
- [ ] Hiển thị trạng thái xử lý

### 3. Test Ingredient Management

#### Access
- [ ] Click vào "Nguyên liệu" trong navigation
- [ ] URL: `/ingredients`
- [ ] View load thành công

#### Dashboard Stats
- [ ] Hiển thị 4 stats cards:
  - Tổng Nguyên Liệu
  - Đủ Hàng (màu xanh)
  - Sắp Hết (màu vàng)
  - Hết Hàng (màu đỏ)

#### Ingredient List
- [ ] Bảng hiển thị danh sách nguyên liệu
- [ ] Columns: Tên, Danh Mục, Số Lượng, Đơn Vị, Trạng Thái, Giá, Thao Tác
- [ ] Search box hoạt động

#### Create Ingredient
- [ ] Click "➕ Thêm Nguyên Liệu"
- [ ] Modal mở ra
- [ ] Fill form:
  - Tên Nguyên Liệu (required)
  - Danh Mục (required) - dropdown
  - Đơn Vị (required) - dropdown
  - Số Lượng (required)
  - Số Lượng Tối Thiểu (required)
  - Giá/Đơn Vị (required)
  - Nhà Cung Cấp
  - Ghi Chú
- [ ] Click "Thêm Mới"
- [ ] Nguyên liệu mới xuất hiện trong danh sách

#### Edit Ingredient
- [ ] Click "✏️ Sửa" trên một nguyên liệu
- [ ] Modal mở với data đã điền
- [ ] Sửa thông tin
- [ ] Click "Cập Nhật"
- [ ] Thông tin được cập nhật

#### Adjust Stock
- [ ] Click "📦 Điều Chỉnh" trên một nguyên liệu
- [ ] Modal mở ra
- [ ] Hiển thị tồn kho hiện tại
- [ ] Select loại điều chỉnh:
  - Nhập Hàng (Add)
  - Xuất Hàng (Remove)
  - Điều Chỉnh (Adjust)
- [ ] Nhập số lượng
- [ ] Nhập lý do
- [ ] Hiển thị số lượng sau điều chỉnh
- [ ] Click "Xác Nhận"
- [ ] Tồn kho được cập nhật

#### View Stock History
- [ ] Click "📊 Lịch Sử" trên một nguyên liệu
- [ ] Modal hiển thị lịch sử
- [ ] Hiển thị:
  - Loại điều chỉnh (Nhập/Xuất/Điều chỉnh)
  - Số lượng
  - Lý do
  - Người thực hiện
  - Thời gian
  - Số lượng trước/sau

#### Delete Ingredient
- [ ] Click "🗑️ Xóa"
- [ ] Confirmation dialog xuất hiện
- [ ] Confirm xóa
- [ ] Nguyên liệu bị xóa khỏi danh sách

#### Low Stock Alert
- [ ] Click "⚠️ Sắp Hết Hàng"
- [ ] Danh sách lọc chỉ hiển thị nguyên liệu sắp hết

### 4. Common Tests

#### Responsive Design
- [ ] Test trên desktop
- [ ] Test trên tablet
- [ ] Test trên mobile
- [ ] Layout điều chỉnh đúng

#### Error Handling
- [ ] Test với network error
- [ ] Error message hiển thị đúng
- [ ] Không crash app

#### Loading States
- [ ] Loading indicator hiển thị khi fetch data
- [ ] Không block UI

## Troubleshooting

### Nếu views không load:

1. **Check console errors**
   ```bash
   # Open browser DevTools (F12)
   # Check Console tab for errors
   ```

2. **Check network requests**
   ```bash
   # Open Network tab
   # Check if API calls are made
   # Check response status codes
   ```

3. **Check authentication**
   ```bash
   # Verify logged in as manager
   # Check localStorage for token
   ```

4. **Check backend**
   ```bash
   # Verify backend is running on port 8080
   curl http://localhost:8080/api/manager/facilities
   ```

5. **Check stores**
   ```javascript
   // In browser console
   import { useFacilityStore } from './stores/facility'
   const store = useFacilityStore()
   console.log(store.items)
   ```

## Expected API Responses

### Facilities
```json
[
  {
    "id": "...",
    "name": "Máy pha cà phê",
    "type": "equipment",
    "model": "Breville BES870XL",
    "location": "Quầy bar",
    "status": "operational",
    "purchase_date": "2024-01-01",
    "maintenance_interval_days": 30,
    "next_maintenance_date": "2024-02-01",
    "notes": "..."
  }
]
```

### Ingredients
```json
[
  {
    "id": "...",
    "name": "Cà phê Arabica",
    "category": "coffee",
    "unit": "kg",
    "quantity": 50,
    "min_quantity": 10,
    "unit_price": 200000,
    "supplier": "Trung Nguyên",
    "notes": "..."
  }
]
```

## Notes

- Views sử dụng stores có sẵn (`facility.js`, `ingredient.js`)
- Stores sử dụng property `items` thay vì `facilities`/`ingredients`
- Backend API endpoints đã có sẵn
- Chỉ cần manager role để access

