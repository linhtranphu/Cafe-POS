# Formatters Utility Implementation

## 📋 Tổng Quan

Đã tạo utility tập trung `frontend/src/utils/formatters.js` để xử lý tất cả các transformations và formatting, tránh lỗi không nhất quán giữa frontend và backend.

## 🎯 Vấn Đề Đã Giải Quyết

### Vấn đề trước đây:
- **Date Format Mismatch**: Frontend gửi `"2026-01-31"` nhưng backend expect `"2026-01-31T00:00:00Z"`
- **Code Duplication**: Mỗi view có hàm `formatDate()` và `formatPrice()` riêng
- **Inconsistency**: Các views format dữ liệu khác nhau
- **Maintenance**: Khó maintain khi cần thay đổi format

### Giải pháp:
✅ Centralized utility với tất cả formatting functions
✅ Consistent date/time handling
✅ Type-safe data sanitization
✅ Reusable across all views

## 📦 Các Functions Có Sẵn

### 1. Date/Time Functions

#### `toISODate(date, includeTime = false)`
Convert date sang ISO format cho backend
```javascript
// Input: "2026-01-31" (from date input)
// Output: "2026-01-31T00:00:00Z"

toISODate("2026-01-31") // "2026-01-31T00:00:00Z"
toISODate(new Date(), true) // "2026-01-31T10:30:45.123Z"
toISODate(null) // null
```

#### `fromISODate(isoDate)`
Convert ISO date sang local format cho date input
```javascript
// Input: "2026-01-31T00:00:00Z"
// Output: "2026-01-31"

fromISODate("2026-01-31T00:00:00Z") // "2026-01-31"
fromISODate(null) // ""
```

#### `formatDate(date, options = {})`
Format date cho display (Vietnamese locale)
```javascript
formatDate("2026-01-31T00:00:00Z") // "31/01/2026"
formatDate(new Date()) // "31/01/2026"
formatDate(null) // "N/A"

// Custom options
formatDate(date, { 
  year: 'numeric', 
  month: 'long', 
  day: 'numeric' 
}) // "31 Tháng 1, 2026"
```

#### `formatDateTime(date)`
Format date và time cho display
```javascript
formatDateTime("2026-01-31T10:30:45Z") // "31/01/2026, 10:30:45"
```

### 2. Number/Currency Functions

#### `formatPrice(price, showSymbol = true)`
Format giá tiền (Vietnamese currency)
```javascript
formatPrice(50000) // "50.000 ₫"
formatPrice(50000, false) // "50.000"
formatPrice(null) // "0 ₫"
```

#### `formatNumber(num)`
Format số với thousand separators
```javascript
formatNumber(1234567) // "1.234.567"
formatNumber(null) // "0"
```

### 3. Data Transformation Functions

#### `sanitizeFormData(data, schema)`
Sanitize form data trước khi gửi backend
```javascript
const formData = {
  name: 'Test',
  purchase_date: '2026-01-31', // date input value
  cost: '50000', // string from input
  supplier: '', // empty string
  notes: null
}

const sanitized = sanitizeFormData(formData, {
  name: { type: 'string' },
  purchase_date: { type: 'date', default: new Date().toISOString() },
  cost: { type: 'number', default: 0 },
  supplier: { type: 'string', default: '' },
  notes: { type: 'string', default: '' }
})

// Result:
// {
//   name: 'Test',
//   purchase_date: '2026-01-31T00:00:00Z', // ✅ ISO format
//   cost: 50000, // ✅ number
//   supplier: '', // ✅ empty string
//   notes: '' // ✅ converted from null
// }
```

#### `parseBackendData(data, schema)`
Parse backend data cho form display
```javascript
const backendData = {
  name: 'Test',
  purchase_date: '2026-01-31T00:00:00Z', // ISO format
  cost: 50000
}

const parsed = parseBackendData(backendData, {
  purchase_date: { type: 'date' }
})

// Result:
// {
//   name: 'Test',
//   purchase_date: '2026-01-31', // ✅ local format for date input
//   cost: 50000
// }
```

### 4. Validation Functions

#### `validateRequired(data, requiredFields)`
Validate required fields
```javascript
const result = validateRequired(formData, ['name', 'type', 'area'])

// Result:
// {
//   valid: false,
//   errors: ['Trường "name" là bắt buộc']
// }
```

#### `deepClone(obj)`
Deep clone object
```javascript
const cloned = deepClone(originalObject)
```

## 🔧 Cách Sử Dụng

### Trong Vue Component

```vue
<script setup>
import { 
  toISODate, 
  fromISODate, 
  formatDate, 
  formatPrice,
  sanitizeFormData,
  parseBackendData
} from '../utils/formatters'

// 1. Khi save form data
const saveFacility = async () => {
  const dataToSend = sanitizeFormData(formData.value, {
    name: { type: 'string' },
    purchase_date: { type: 'date', default: new Date().toISOString() },
    cost: { type: 'number', default: 0 }
  })
  
  await facilityStore.createFacility(dataToSend)
}

// 2. Khi load data vào form
const openEditModal = (facility) => {
  formData.value = parseBackendData({ ...facility }, {
    purchase_date: { type: 'date' }
  })
}

// 3. Trong template
const displayDate = formatDate(facility.purchase_date)
const displayPrice = formatPrice(facility.cost)
</script>

<template>
  <div>{{ formatDate(item.date) }}</div>
  <div>{{ formatPrice(item.price) }}</div>
</template>
```

## ✅ Views Đã Cập Nhật

### 1. FacilityManagementView.vue ✅
- ✅ Import formatters utility
- ✅ Sử dụng `sanitizeFormData()` trong `saveFacility()`
- ✅ Sử dụng `parseBackendData()` trong `openEditModal()`
- ✅ Sử dụng `formatDate()` và `formatPrice()` trong template
- ✅ Removed duplicate formatDate/formatPrice functions

## 📝 Views Cần Cập Nhật

Các views sau vẫn có duplicate formatDate/formatPrice functions:

1. **UserManagementView.vue**
   - Has: `formatDate()`
   - Should use: `formatDate` from formatters

2. **ShiftView.vue**
   - Has: `formatPrice()`, `formatDate()`
   - Should use: `formatPrice`, `formatDate` from formatters

3. **IngredientView.vue**
   - Has: `formatPrice()`, `formatDate()`
   - Should use: `formatPrice`, `formatDate` from formatters

4. **ExpenseView.vue**
   - Has: `formatPrice()`, `formatDate()`
   - Should use: `formatPrice`, `formatDate` from formatters

5. **MenuView.vue**
   - Has: `formatPrice()`
   - Should use: `formatPrice` from formatters

6. **CashierDashboard.vue**
   - Has: `formatPrice()`, `formatDate()`
   - Should use: `formatPrice`, `formatDate` from formatters

7. **CashierShiftClosure.vue**
   - Has: `formatPrice()`
   - Should use: `formatPrice` from formatters

8. **OrderView.vue**
   - Has: `formatPrice()`, `formatDate()`
   - Should use: `formatPrice`, `formatDate` from formatters

## 🎯 Migration Plan

### Phase 1: Critical Views (Date Handling) ✅
- [x] FacilityManagementView.vue - DONE

### Phase 2: High Priority Views
- [ ] IngredientManagementView.vue - Has date inputs
- [ ] ExpenseView.vue - Has date inputs
- [ ] CashierShiftClosure.vue - Has date/time display

### Phase 3: Medium Priority Views
- [ ] UserManagementView.vue
- [ ] ShiftView.vue
- [ ] OrderView.vue
- [ ] CashierDashboard.vue

### Phase 4: Low Priority Views
- [ ] MenuView.vue
- [ ] BaristaView.vue
- [ ] DashboardView.vue

## 🔍 Testing Checklist

Sau khi migrate mỗi view, test:

- [ ] Create new record với date field
- [ ] Edit existing record với date field
- [ ] Display dates correctly in list view
- [ ] Display prices correctly with Vietnamese format
- [ ] Form validation works
- [ ] No console errors

## 📚 Best Practices

### DO ✅
- Always use `sanitizeFormData()` trước khi gửi data lên backend
- Always use `parseBackendData()` khi load data vào form
- Always use `formatDate()` và `formatPrice()` cho display
- Define schema rõ ràng cho mỗi form

### DON'T ❌
- Không tự convert date format manually
- Không tạo duplicate formatDate/formatPrice functions
- Không gửi date string trực tiếp từ date input lên backend
- Không hardcode date format strings

## 🐛 Common Issues & Solutions

### Issue 1: Date không được accept bởi backend
**Symptom**: 400 Bad Request khi create/update
**Solution**: Sử dụng `sanitizeFormData()` với schema type 'date'

### Issue 2: Date input không hiển thị giá trị khi edit
**Symptom**: Date input trống khi mở edit modal
**Solution**: Sử dụng `parseBackendData()` với schema type 'date'

### Issue 3: Price hiển thị không đúng format
**Symptom**: "50000" thay vì "50.000 ₫"
**Solution**: Sử dụng `formatPrice()` trong template

## 📖 Related Files

- **Utility**: `frontend/src/utils/formatters.js`
- **Constants**: `frontend/src/constants/facility.js`
- **Example Usage**: `frontend/src/views/FacilityManagementView.vue`
- **Test Script**: `test-create-facility.sh`

## 🎉 Benefits

1. **Consistency**: Tất cả views format data giống nhau
2. **Maintainability**: Chỉ cần update 1 file khi cần thay đổi format
3. **Type Safety**: Schema-based validation và transformation
4. **Error Prevention**: Tránh lỗi date format mismatch
5. **Code Reuse**: Không duplicate code
6. **Testing**: Dễ test vì logic tập trung

## 🔄 Future Enhancements

- [ ] Add unit tests cho formatters.js
- [ ] Add TypeScript types
- [ ] Add more validation functions
- [ ] Add currency conversion support
- [ ] Add timezone handling
- [ ] Add locale switching support
