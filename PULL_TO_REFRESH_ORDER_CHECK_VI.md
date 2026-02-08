# Kiểm tra thứ tự khởi tạo Pull-to-Refresh

## Lỗi tìm thấy
**CashierDashboard.vue** có lỗi khi `usePullToRefresh(refreshData)` được gọi TRƯỚC khi `refreshData` được định nghĩa, gây ra:
```
ReferenceError: Cannot access 'refreshData' before initialization
```

## Nguyên nhân
Trong JavaScript, arrow functions (`const refreshData = async () => {}`) KHÔNG được hoisted như function declarations. Chúng phải được định nghĩa trước khi sử dụng.

## Cách sửa
Di chuyển định nghĩa `refreshData` lên TRƯỚC khi gọi `usePullToRefresh()` trong CashierDashboard.vue.

## Tất cả Views đã kiểm tra ✅

### Views có thứ tự đúng (refreshData định nghĩa TRƯỚC usePullToRefresh)

1. **BaristaView.vue** - ✅ Không có pull-to-refresh (barista không cần)
2. **CashierDashboard.vue** - ✅ ĐÃ SỬA (trước đó bị lỗi, giờ đã đúng)
3. **CashierHandoverView.vue** - ✅ Thứ tự đúng
4. **CashierReports.vue** - ✅ Thứ tự đúng
5. **CashierShiftClosure.vue** - ✅ Thứ tự đúng
6. **DashboardView.vue** - ✅ Thứ tự đúng
7. **ExpenseManagementView.vue** - ✅ Thứ tự đúng
8. **FacilityManagementView.vue** - ✅ Thứ tự đúng
9. **IngredientManagementView.vue** - ✅ Thứ tự đúng
10. **LoginView.vue** - ✅ Không có pull-to-refresh (trang login không cần)
11. **ManagerShiftView.vue** - ✅ Thứ tự đúng
12. **MenuView.vue** - ✅ Thứ tự đúng
13. **OrderView.vue** - ✅ Không có pull-to-refresh (nhận order không cần)
14. **ProfileView.vue** - ✅ Thứ tự đúng
15. **ShiftView.vue** - ✅ Thứ tự đúng
16. **UserManagementView.vue** - ✅ Thứ tự đúng

## Pattern cần tuân theo

### ✅ THỨ TỰ ĐÚNG
```javascript
// 1. Định nghĩa refreshData TRƯỚC
const refreshData = async () => {
  await someStore.fetchData()
}

// 2. Sử dụng nó trong usePullToRefresh SAU
const { pullDistance, isRefreshing } = usePullToRefresh(refreshData)
```

### ❌ THỨ TỰ SAI (gây lỗi)
```javascript
// 1. Sử dụng refreshData TRƯỚC khi định nghĩa
const { pullDistance, isRefreshing } = usePullToRefresh(refreshData)

// 2. Định nghĩa refreshData SAU (QUÁ MUỘN!)
const refreshData = async () => {
  await someStore.fetchData()
}
```

## Tại sao điều này quan trọng

**Arrow functions KHÔNG được hoisted:**
```javascript
// Cái này BỊ LỖI
console.log(myFunc()) // ReferenceError
const myFunc = () => 'hello'
```

**Function declarations ĐƯỢC hoisted:**
```javascript
// Cái này HOẠT ĐỘNG
console.log(myFunc()) // 'hello'
function myFunc() { return 'hello' }
```

## Checklist kiểm tra
- [x] CashierDashboard.vue - Đã sửa và test
- [x] Tất cả views khác - Đã xác nhận thứ tự đúng
- [x] Không tìm thấy lỗi khởi tạo nào khác

## Tóm tắt
✅ **Tất cả views giờ có thứ tự khởi tạo đúng**
✅ **Lỗi CashierDashboard.vue đã được sửa**
✅ **Không có view nào khác bị lỗi này**

---
**Trạng thái**: ✅ Hoàn thành
**Ngày**: 2026-02-07
**File đã sửa**: 1 (CashierDashboard.vue)
**File đã kiểm tra**: 16 views
