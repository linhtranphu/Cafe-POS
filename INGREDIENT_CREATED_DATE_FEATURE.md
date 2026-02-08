# Feature: Ngày Tạo Nguyên Liệu & Hiển Thị Trong Lịch Sử

## 🎯 Mục Tiêu

Cho phép user chọn ngày tạo khi tạo nguyên liệu mới và hiển thị ngày giờ đầy đủ trong lịch sử.

---

## ✅ Thay Đổi

### 1. Frontend - Form Tạo Nguyên Liệu

#### A. Thêm Field "Ngày tạo"

**File**: `frontend/src/views/IngredientManagementView.vue`

**Location**: Sau field "Ghi chú" trong form tạo

```vue
<!-- Ngày tạo (only when creating) -->
<div v-if="!isEditing">
  <label class="block text-sm font-medium text-gray-700 mb-3">
    Ngày tạo
    <span class="text-xs text-gray-500 font-normal">(mặc định: hôm nay)</span>
  </label>
  <input v-model="formData.created_date" type="datetime-local"
    class="w-full px-4 py-4 text-base border border-gray-300 rounded-lg focus:ring-2 focus:ring-blue-500 focus:border-transparent" />
  <p class="text-xs text-gray-500 mt-2">
    💡 Để trống để sử dụng thời gian hiện tại
  </p>
</div>
```

**Features**:
- Chỉ hiển thị khi **tạo mới** (không hiển thị khi edit)
- Input type: `datetime-local` (cho phép chọn cả ngày và giờ)
- Default value: Thời gian hiện tại
- Optional: User có thể để trống

#### B. Update FormData

```javascript
const formData = ref({
  name: '',
  category: '',
  unit: '',
  quantity: 0,
  min_stock: 0,
  cost_per_unit: 0,
  supplier: '',
  notes: '',
  created_date: '' // Will be set to current datetime when opening create modal
})
```

#### C. Set Default DateTime

```javascript
const openCreateModal = () => {
  isEditing.value = false
  currentIngredient.value = null
  
  // Set default created_date to current datetime in local timezone
  const now = new Date()
  const year = now.getFullYear()
  const month = String(now.getMonth() + 1).padStart(2, '0')
  const day = String(now.getDate()).padStart(2, '0')
  const hours = String(now.getHours()).padStart(2, '0')
  const minutes = String(now.getMinutes()).padStart(2, '0')
  const defaultDateTime = `${year}-${month}-${day}T${hours}:${minutes}`
  
  formData.value = {
    // ... other fields
    created_date: defaultDateTime
  }
  // ...
}
```

### 2. Frontend - Hiển Thị Lịch Sử

#### A. Cải Thiện formatCompactDate

**File**: `frontend/src/views/IngredientManagementView.vue`

**Before**:
```javascript
const formatCompactDate = (date) => {
  // ... relative time logic ...
  // More than 7 days: show date only
  return d.toLocaleDateString('vi-VN', { day: '2-digit', month: '2-digit' })
}
```

**After**:
```javascript
const formatCompactDate = (date) => {
  if (!date) return 'N/A'
  const d = new Date(date)
  const now = new Date()
  const diffMs = now - d
  const diffMins = Math.floor(diffMs / 60000)
  const diffHours = Math.floor(diffMs / 3600000)
  const diffDays = Math.floor(diffMs / 86400000)
  
  // Less than 1 hour: show minutes
  if (diffMins < 60) {
    return `${diffMins} phút trước`
  }
  // Less than 24 hours: show hours
  if (diffHours < 24) {
    return `${diffHours} giờ trước`
  }
  // Less than 7 days: show days
  if (diffDays < 7) {
    return `${diffDays} ngày trước`
  }
  // More than 7 days: show full date and time
  return d.toLocaleString('vi-VN', {
    day: '2-digit',
    month: '2-digit',
    year: 'numeric',
    hour: '2-digit',
    minute: '2-digit'
  })
}
```

**Improvements**:
- Hiển thị **ngày giờ đầy đủ** cho records cũ hơn 7 ngày
- Format: `DD/MM/YYYY, HH:MM` (VD: `05/02/2026, 14:30`)

#### B. Lịch Sử Hiển Thị

Lịch sử đã hiển thị ngày tạo ở góc phải:

```vue
<span class="text-[10px] text-gray-500">
  {{ formatCompactDate(record.created_at) }}
</span>
```

**Hiển thị**:
- `5 phút trước` (< 1 giờ)
- `2 giờ trước` (< 24 giờ)
- `3 ngày trước` (< 7 ngày)
- `05/02/2026, 14:30` (≥ 7 ngày)

### 3. Backend Changes

#### A. Update CreateIngredientRequest

**File**: `backend/domain/ingredient/ingredient.go`

```go
type CreateIngredientRequest struct {
	Name        string    `json:"name" binding:"required"`
	Category    string    `json:"category" binding:"required"`
	Unit        UnitType  `json:"unit" binding:"required"`
	Quantity    float64   `json:"quantity" binding:"required,min=0"`
	MinStock    float64   `json:"min_stock" binding:"min=0"`
	CostPerUnit float64   `json:"cost_per_unit" binding:"min=0"`
	Supplier    string    `json:"supplier"`
	CreatedDate *string   `json:"created_date"` // Optional: custom creation date (ISO 8601 format)
}
```

**Features**:
- Field `CreatedDate` là optional (`*string`)
- Nhận format: `YYYY-MM-DDTHH:MM` hoặc ISO 8601 full

#### B. Update CreateIngredient Service

**File**: `backend/application/services/ingredient.go`

```go
func (s *IngredientService) CreateIngredient(ctx context.Context, req *ingredient.CreateIngredientRequest, userIDStr string, username string) (*ingredient.Ingredient, error) {
	item := &ingredient.Ingredient{
		Name:        req.Name,
		Category:    req.Category,
		Unit:        req.Unit,
		Quantity:    req.Quantity,
		MinStock:    req.MinStock,
		CostPerUnit: req.CostPerUnit,
		Supplier:    req.Supplier,
	}
	
	// Set created_at from request if provided, otherwise use current time
	if req.CreatedDate != nil && *req.CreatedDate != "" {
		if parsedTime, err := time.Parse(time.RFC3339, *req.CreatedDate); err == nil {
			item.CreatedAt = parsedTime
		} else {
			// Try parsing as datetime-local format (YYYY-MM-DDTHH:MM)
			if parsedTime, err := time.Parse("2006-01-02T15:04", *req.CreatedDate); err == nil {
				item.CreatedAt = parsedTime
			} else {
				item.CreatedAt = time.Now() // Fallback to current time
			}
		}
	} else {
		item.CreatedAt = time.Now()
	}
	item.UpdatedAt = item.CreatedAt

	// ... rest of code
	
	// Create stock history with same created_at
	history := &ingredient.StockHistory{
		// ... other fields
		CreatedAt: item.CreatedAt, // Use the same created_at as ingredient
	}
	// ...
}
```

**Logic**:
1. Nếu `req.CreatedDate` có giá trị → Parse và sử dụng
2. Thử parse ISO 8601 format trước (`2006-01-02T15:04:05Z07:00`)
3. Nếu fail, thử parse datetime-local format (`2006-01-02T15:04`)
4. Nếu vẫn fail hoặc không có giá trị → Dùng `time.Now()`
5. Stock history cũng sử dụng cùng `created_at`

#### C. Update Repository

**File**: `backend/infrastructure/mongodb/ingredient_repository.go`

```go
func (r *IngredientRepository) Create(ctx context.Context, item *ingredient.Ingredient) error {
	// Only set CreatedAt if not already set (allow custom creation date)
	if item.CreatedAt.IsZero() {
		item.CreatedAt = time.Now()
	}
	// Always update UpdatedAt
	item.UpdatedAt = time.Now()
	result, err := r.collection.InsertOne(ctx, item)
	// ...
}
```

**Logic**:
- Chỉ set `CreatedAt = time.Now()` nếu chưa được set
- Cho phép service set custom `CreatedAt` trước khi gọi repository

---

## 📊 Use Cases

### Use Case 1: Tạo Nguyên Liệu Với Ngày Hiện Tại

**Bước thực hiện**:
1. Click "➕ Tạo nguyên liệu"
2. Điền thông tin
3. **Không thay đổi** field "Ngày tạo" (để mặc định)
4. Click "Thêm mới"

**Kết quả**:
- Nguyên liệu được tạo với `created_at = thời gian hiện tại`
- Lịch sử hiển thị: "Vài phút trước"

### Use Case 2: Tạo Nguyên Liệu Với Ngày Tùy Chỉnh

**Bước thực hiện**:
1. Click "➕ Tạo nguyên liệu"
2. Điền thông tin
3. **Thay đổi** field "Ngày tạo" → Chọn: `01/02/2026 10:00`
4. Click "Thêm mới"

**Kết quả**:
- Nguyên liệu được tạo với `created_at = 01/02/2026 10:00`
- Lịch sử hiển thị: "01/02/2026, 10:00" (nếu > 7 ngày)

### Use Case 3: Xem Lịch Sử

**Bước thực hiện**:
1. Click "📊 Lịch sử" của nguyên liệu
2. Xem các records

**Hiển thị**:
- Record mới (< 1h): "15 phút trước"
- Record hôm nay (< 24h): "5 giờ trước"
- Record tuần này (< 7 ngày): "3 ngày trước"
- Record cũ (≥ 7 ngày): "01/02/2026, 10:00"

---

## 🧪 Testing

### Test Case 1: Default DateTime

1. Mở form tạo nguyên liệu
2. Kiểm tra field "Ngày tạo"
3. **Kết quả mong đợi**: Hiển thị thời gian hiện tại (VD: `2026-02-07T23:00`)

### Test Case 2: Custom DateTime

1. Mở form tạo nguyên liệu
2. Thay đổi "Ngày tạo" → `2026-02-01T10:00`
3. Điền thông tin khác và tạo
4. Xem lịch sử
5. **Kết quả mong đợi**: Lịch sử hiển thị "01/02/2026, 10:00"

### Test Case 3: Empty DateTime

1. Mở form tạo nguyên liệu
2. Xóa hết field "Ngày tạo" (để trống)
3. Tạo nguyên liệu
4. **Kết quả mong đợi**: Sử dụng thời gian hiện tại

### Test Case 4: History Display

1. Tạo nguyên liệu với ngày cũ (VD: 1 tháng trước)
2. Tạo nguyên liệu với ngày hiện tại
3. Xem lịch sử
4. **Kết quả mong đợi**:
   - Record cũ: Hiển thị ngày giờ đầy đủ
   - Record mới: Hiển thị "vài phút trước"

---

## 📁 Files Changed

### Frontend
- `frontend/src/views/IngredientManagementView.vue`
  - Added `created_date` field to form
  - Updated `formData` initialization
  - Updated `openCreateModal()` to set default datetime
  - Improved `formatCompactDate()` to show full datetime for old records

### Backend
- `backend/domain/ingredient/ingredient.go`
  - Added `CreatedDate *string` to `CreateIngredientRequest`
- `backend/application/services/ingredient.go`
  - Added `time` import
  - Updated `CreateIngredient()` to parse and use custom created_date
  - Set stock history `CreatedAt` to same value
- `backend/infrastructure/mongodb/ingredient_repository.go`
  - Updated `Create()` to only set `CreatedAt` if not already set

---

## 🎯 Benefits

1. **Flexibility**: User có thể chọn ngày tạo tùy ý (VD: nhập liệu lịch sử)
2. **Default Value**: Mặc định là thời gian hiện tại (UX tốt)
3. **Full DateTime**: Lịch sử hiển thị cả ngày và giờ
4. **Smart Display**: Hiển thị relative time cho records gần, full datetime cho records cũ
5. **Backward Compatible**: Không ảnh hưởng đến data cũ

---

## 📝 Notes

### Date Format

**Frontend → Backend**:
- Format: `YYYY-MM-DDTHH:MM` (từ `<input type="datetime-local">`)
- Example: `2026-02-07T23:00`

**Backend Parsing**:
1. Try ISO 8601: `2006-01-02T15:04:05Z07:00`
2. Try datetime-local: `2006-01-02T15:04`
3. Fallback: `time.Now()`

**Display**:
- Recent (< 7 days): Relative time
- Old (≥ 7 days): Full datetime `DD/MM/YYYY, HH:MM`

### Timezone

- Frontend: Local timezone của user
- Backend: Parse as-is (không convert timezone)
- Database: Lưu as-is

---

## ✅ Checklist

- ✅ Frontend: Added created_date field to form
- ✅ Frontend: Set default datetime
- ✅ Frontend: Improved formatCompactDate
- ✅ Backend: Added CreatedDate to request
- ✅ Backend: Parse custom created_date
- ✅ Backend: Update repository to allow custom CreatedAt
- ✅ Backend: Set stock history CreatedAt
- ✅ Backend: Restart server
- ⏳ Testing: User test

---

**Ngày**: 2026-02-07  
**Status**: ✅ COMPLETE - Ready for testing  
**Backend**: ✅ Running on :3000  
**Frontend**: ✅ Running on :5173
