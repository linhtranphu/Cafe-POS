# Manager Views Fix - Facility & Ingredient

## Vấn Đề

Manager không thấy views cho Facility và Ingredient management.

## Nguyên Nhân

Views đã được tạo nhưng có một số vấn đề nhỏ:
1. ✅ Views sử dụng computed properties không đúng với store structure
2. ✅ Store sử dụng `items` thay vì `facilities`/`ingredients`

## Giải Pháp

### 1. Fixed FacilityManagementView.vue ✅

**Before**:
```javascript
const facilities = computed(() => facilityStore.facilities)
const maintenanceSchedule = computed(() => facilityStore.maintenanceSchedule)
const issueReports = computed(() => facilityStore.issueReports)
```

**After**:
```javascript
const facilities = computed(() => facilityStore.items || [])
const maintenanceSchedule = ref([])
const issueReports = ref([])
```

**Reason**: Store sử dụng `items` property, không phải `facilities`

### 2. Fixed IngredientManagementView.vue ✅

**Before**:
```javascript
const ingredients = computed(() => ingredientStore.ingredients)
const stockHistory = computed(() => ingredientStore.stockHistory)
```

**After**:
```javascript
const ingredients = computed(() => ingredientStore.items || [])
const stockHistory = ref([])
```

**Reason**: Store sử dụng `items` property, không phải `ingredients`

### 3. Fixed Data Loading ✅

**FacilityManagementView.vue**:
```javascript
onMounted(async () => {
  await facilityStore.fetchFacilities()
  maintenanceSchedule.value = await facilityStore.fetchScheduledMaintenance()
  issueReports.value = await facilityStore.fetchIssueReports()
})
```

**IngredientManagementView.vue**:
```javascript
const viewHistory = async (ingredient) => {
  currentIngredient.value = ingredient
  stockHistory.value = await ingredientStore.fetchStockHistory(ingredient.id)
  showHistoryModal.value = true
}
```

## Files Modified

1. ✅ `frontend/src/views/FacilityManagementView.vue`
   - Fixed computed properties
   - Fixed data loading

2. ✅ `frontend/src/views/IngredientManagementView.vue`
   - Fixed computed properties
   - Fixed data loading

## How to Test

### 1. Start Backend
```bash
cd backend
go run main.go
# Backend should run on http://localhost:8080
```

### 2. Start Frontend
```bash
cd frontend
npm run dev
# Frontend should run on http://localhost:5173
```

### 3. Login as Manager
- Username: `admin` (or any manager account)
- Password: your password

### 4. Navigate to Views
- Click "Cơ sở vật chất" (Facilities) in navigation
- Click "Nguyên liệu" (Ingredients) in navigation

### 5. Verify Features Work
- ✅ Dashboard stats display
- ✅ List displays
- ✅ Search works
- ✅ Create/Edit/Delete work
- ✅ Modals open correctly

## Expected Behavior

### Facility Management
- **URL**: `/facilities`
- **Stats**: 4 cards (Total, Operational, Maintenance, Broken)
- **Table**: List of facilities with actions
- **Modals**: Create/Edit, Maintenance Schedule, Issue Reports

### Ingredient Management
- **URL**: `/ingredients`
- **Stats**: 4 cards (Total, In Stock, Low Stock, Out of Stock)
- **Table**: List of ingredients with actions
- **Modals**: Create/Edit, Adjust Stock, Stock History

## API Endpoints Used

### Facility
```
GET    /api/manager/facilities
POST   /api/manager/facilities
PUT    /api/manager/facilities/:id
DELETE /api/manager/facilities/:id
GET    /api/manager/maintenance/scheduled
GET    /api/manager/issues
```

### Ingredient
```
GET    /api/manager/ingredients
POST   /api/manager/ingredients
PUT    /api/manager/ingredients/:id
DELETE /api/manager/ingredients/:id
POST   /api/manager/ingredients/:id/adjust
GET    /api/manager/ingredients/:id/history
GET    /api/manager/ingredients/low-stock
```

## Store Structure

### Facility Store (`stores/facility.js`)
```javascript
state: {
  items: [],        // ← Array of facilities
  loading: false,
  error: null
}
```

### Ingredient Store (`stores/ingredient.js`)
```javascript
state: {
  items: [],        // ← Array of ingredients
  lowStockItems: [],
  loading: false,
  error: null
}
```

## Common Issues & Solutions

### Issue 1: "Cannot read property 'length' of undefined"
**Solution**: Added `|| []` fallback in computed properties

### Issue 2: Maintenance schedule not loading
**Solution**: Changed from computed to ref and load data in onMounted

### Issue 3: Stock history not showing
**Solution**: Store result in ref variable instead of relying on store state

## Status

✅ **FIXED and READY TO USE**

Manager có thể:
- ✅ Xem danh sách facilities và ingredients
- ✅ Thêm mới
- ✅ Chỉnh sửa
- ✅ Xóa
- ✅ Điều chỉnh tồn kho (ingredients)
- ✅ Xem lịch sử
- ✅ Xem lịch bảo trì (facilities)
- ✅ Xem báo cáo sự cố (facilities)

