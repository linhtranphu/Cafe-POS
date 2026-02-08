# 🧪 Hướng Dẫn Test Pull-to-Refresh

## ✅ Đã Fix

Pull-to-refresh giờ đây chỉ hoạt động khi **ở đầu trang** (scrollTop = 0).

Khi đang scroll giữa trang → Pull down = Tiếp tục scroll bình thường (KHÔNG refresh).

## 📱 Cách Test Trên iPhone 14

### Test 1: Pull-to-Refresh Ở Đầu Trang ✅

**Các bước:**
1. Mở bất kỳ view nào (Dashboard, Barista, Orders, v.v.)
2. Scroll lên đầu trang (top)
3. Kéo xuống từ từ
4. **Kỳ vọng:**
   - ✅ Thấy icon ⬇️ và text "Kéo xuống để làm mới"
   - ✅ Kéo qua 80px → Icon đổi thành 🎯 và text "Thả để làm mới"
   - ✅ Thả ra → Icon 🔄 quay và text "Đang tải..."
   - ✅ Data reload thành công

### Test 2: Scroll Bình Thường Ở Giữa Trang ✅

**Các bước:**
1. Mở view có nhiều content (Dashboard, Orders)
2. Scroll xuống giữa trang
3. Thử kéo xuống để scroll tiếp
4. **Kỳ vọng:**
   - ✅ Scroll tiếp bình thường
   - ❌ KHÔNG thấy pull-to-refresh indicator
   - ❌ KHÔNG trigger refresh
   - ✅ Scroll mượt mà, không giật lag

### Test 3: Pull Rồi Scroll ✅

**Các bước:**
1. Scroll lên đầu trang
2. Bắt đầu kéo xuống (thấy indicator)
3. Trước khi thả, scroll xuống
4. **Kỳ vọng:**
   - ✅ Indicator biến mất
   - ✅ Scroll tiếp bình thường
   - ❌ KHÔNG trigger refresh

### Test 4: Scroll Nhanh (Flick) ✅

**Các bước:**
1. Scroll lên đầu trang
2. Vuốt nhanh xuống (flick gesture)
3. **Kỳ vọng:**
   - ✅ Scroll nhanh xuống
   - ❌ KHÔNG thấy pull-to-refresh
   - ❌ KHÔNG trigger refresh

## 📋 Danh Sách Views Cần Test

### Container Scroll Views (16 views):

- [ ] **Dashboard** (`/dashboard`)
  - Test: Scroll stats, quick actions
  
- [ ] **Barista** (`/barista`)
  - Test: Queue, working, ready tabs
  
- [ ] **Cashier Dashboard** (`/cashier`)
  - Test: Payment list, stats
  
- [ ] **Cashier Handover** (`/cashier/handover`)
  - Test: Handover list
  
- [ ] **Cashier Reports** (`/cashier/reports`)
  - Test: Report list
  
- [ ] **Cashier Shift Closure** (`/cashier/shift-closure`)
  - Test: Closure form
  
- [ ] **Expense Management** (`/expenses`)
  - Test: Expense list
  
- [ ] **Facility Management** (`/facilities`)
  - Test: Facility list
  
- [ ] **Facility Add/Edit** (`/facilities/add`, `/facilities/edit/:id`)
  - Test: Form fields
  
- [ ] **Ingredient Management** (`/ingredients`)
  - Test: Ingredient list
  
- [ ] **Manager Shifts** (`/manager/shifts`)
  - Test: Shift list, tabs
  
- [ ] **Menu** (`/menu`)
  - Test: Menu items, categories
  
- [ ] **Orders** (`/orders`)
  - Test: Order list, tabs
  
- [ ] **Profile** (`/profile`)
  - Test: Profile info
  
- [ ] **Shifts** (`/shifts`)
  - Test: Shift list, current shift
  
- [ ] **User Management** (`/users`)
  - Test: User list

## ✅ Kết Quả Mong Đợi

### PASS ✅

Tất cả views phải:
- ✅ Pull-to-refresh chỉ hoạt động ở đầu trang
- ✅ Scroll bình thường ở mọi vị trí khác
- ✅ Không có giật lag khi scroll
- ✅ Indicator hiện/ẩn đúng lúc
- ✅ Refresh trigger đúng khi thả ở đầu trang

### FAIL ❌

Nếu thấy:
- ❌ Pull-to-refresh trigger khi đang scroll giữa trang
- ❌ Không scroll được bình thường
- ❌ Indicator hiện khi đang scroll
- ❌ Scroll bị giật, lag
- ❌ Refresh trigger khi không mong muốn

→ **Báo lại để fix!**

## 🔍 Debug

Nếu có vấn đề, check:

1. **Console log:**
   ```javascript
   // Thêm vào usePullToRefresh.js để debug
   console.log('scrollTop:', getScrollTop())
   console.log('isPulling:', isPulling.value)
   console.log('pullDistance:', pullDistance.value)
   ```

2. **Scroll container:**
   ```javascript
   // Check xem có tìm được scroll container không
   console.log('scrollContainer:', scrollContainer)
   ```

3. **Touch events:**
   ```javascript
   // Check touch events có fire không
   console.log('touchStart:', startY)
   console.log('touchMove:', currentY)
   ```

## 💡 Tips

### Để Test Tốt:

1. **Test trên real device** (iPhone 14) - Không dùng simulator
2. **Test ở webapp mode** - Add to Home Screen
3. **Test nhiều lần** - Đảm bảo consistent
4. **Test các edge cases** - Pull nhanh, pull chậm, pull rồi scroll
5. **Test tất cả views** - Đảm bảo không có regression

### Nếu Thấy Bug:

1. Note lại view nào bị bug
2. Note lại hành vi cụ thể
3. Note lại steps to reproduce
4. Screenshot/video nếu có thể
5. Báo lại để fix

## 🎯 Mục Tiêu

**Trải nghiệm người dùng:**
- Scroll mượt mà, tự nhiên
- Pull-to-refresh chỉ khi cần
- Không có hành vi bất ngờ
- Không có giật lag

**Kỹ thuật:**
- Container scroll detection hoạt động đúng
- Scroll position check chính xác
- Touch events handle đúng
- Performance tốt (60fps)

---

**Ngày:** 6 tháng 2, 2026  
**Tính năng:** Pull-to-Refresh Container Scroll  
**Trạng thái:** ✅ Đã implement, chờ test  
**Người test:** User trên iPhone 14
