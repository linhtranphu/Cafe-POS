# 🔧 Cash Handover Troubleshooting - Waiter Role

## ❌ Vấn Đề: Không Thấy Nút Bàn Giao

**URL:** `http://localhost:5173/#/shifts`  
**Role:** Waiter  
**Vấn đề:** Không có nút "💰 Bàn giao một phần" hoặc "🏁 Bàn giao và đóng ca"

---

## 🔍 Checklist Debug

### ✅ Bước 1: Kiểm Tra Role

**Mở Console (F12) và chạy:**
```javascript
console.log('User role:', localStorage.getItem('user'))
```

**Expected:**
```json
{
  "id": "...",
  "username": "waiter1",
  "role": "waiter",  // ✅ Phải là "waiter"
  "name": "..."
}
```

**Nếu không phải waiter:**
- Logout
- Login lại với tài khoản waiter

---

### ✅ Bước 2: Kiểm Tra Ca Đang Mở

**Trong trang `/shifts`, bạn phải thấy:**

```
┌─────────────────────────────────────┐
│  Ca đang mở                         │
│  ☀️ Ca sáng                         │
│  Phục vụ                            │
│                    [ĐANG MỞ]        │
└─────────────────────────────────────┘
```

**Nếu KHÔNG thấy card "Ca đang mở":**

1. Scroll xuống tìm form "Mở ca làm việc"
2. Chọn ca (☀️ Ca sáng / 🌤️ Ca chiều / 🌙 Ca tối)
3. Nhập tiền đầu ca: `0` (hoặc số tiền bạn muốn)
4. Click "Mở ca"

---

### ✅ Bước 3: Kiểm Tra Tiền Hiện Có

**Trong card "Ca đang mở", kiểm tra:**

```
┌─────────────────────────────────────┐
│  Tiền hiện có: 0đ          ❌       │
│  Đã bàn giao: 0đ                    │
│  Tổng thu: 0đ                       │
└─────────────────────────────────────┘
```

**Nếu "Tiền hiện có" = 0:**
- ❌ Nút bàn giao KHÔNG hiển thị
- ✅ Chỉ có nút "Kết thúc ca"

**Để có tiền:**

#### Option 1: Tạo Order với Cash Payment
```
1. Go to /orders
2. Click "Tạo đơn mới"
3. Chọn bàn
4. Thêm món (từ menu đã seed)
5. Click "Hoàn thành"
6. Chọn payment: "Tiền mặt"
7. Nhập số tiền
8. Click "Thanh toán"
9. Quay lại /shifts
10. Tiền hiện có > 0 ✅
```

#### Option 2: Mở Ca với Tiền Đầu Ca
```
1. Khi mở ca mới
2. Nhập tiền đầu ca: 500000
3. Click "Mở ca"
4. Tiền hiện có = 500,000đ ✅
```

---

### ✅ Bước 4: Kiểm Tra Pending Handover

**Nếu có bàn giao đang chờ:**

```
┌─────────────────────────────────────┐
│  🕐 Đang chờ xác nhận bàn giao      │
│  500,000đ                   [Hủy]  │
└─────────────────────────────────────┘
```

**Khi có pending:**
- ❌ Nút bàn giao bị ẩn
- ✅ Chỉ thấy "Chờ cashier xác nhận..."

**Giải pháp:**
- Chờ cashier xác nhận
- Hoặc click "Hủy" để hủy bàn giao

---

### ✅ Bước 5: Kiểm Tra Frontend Build

**Nếu code đã update nhưng UI không thay đổi:**

```bash
# Stop frontend
Ctrl+C

# Rebuild
cd frontend
npm run build

# Restart
npm run dev
```

**Hoặc hard refresh browser:**
- Chrome/Edge: `Ctrl+Shift+R` (Windows) / `Cmd+Shift+R` (Mac)
- Firefox: `Ctrl+F5` (Windows) / `Cmd+Shift+R` (Mac)

---

## 🎯 Kịch Bản Đầy Đủ

### Scenario: Từ Đầu Đến Khi Thấy Nút Bàn Giao

```
1. Login as waiter
   Username: waiter1
   Password: password123
   ↓
2. Navigate to /shifts
   ↓
3. Nếu chưa có ca → Mở ca
   - Chọn: ☀️ Ca sáng
   - Tiền đầu ca: 0
   - Click "Mở ca"
   ↓
4. Thấy card "Ca đang mở"
   - Tiền hiện có: 0đ ❌
   - Chỉ có nút "Kết thúc ca"
   ↓
5. Tạo order để có tiền
   - Go to /orders
   - Tạo order mới
   - Thêm món: Cà phê đen (25,000đ)
   - Hoàn thành order
   - Payment: Tiền mặt
   - Nhập: 30,000đ
   - Thanh toán
   ↓
6. Quay lại /shifts
   - Tiền hiện có: 25,000đ ✅
   - Thấy nút "💰 Bàn giao một phần" ✅
   - Thấy nút "🏁 Bàn giao và đóng ca" ✅
```

---

## 📸 Screenshots Mẫu

### 1. Không Có Tiền (Tiền hiện có = 0)

```
┌─────────────────────────────────────┐
│  Ca đang mở                         │
│  ☀️ Ca sáng                         │
├─────────────────────────────────────┤
│  Tiền hiện có: 0đ                   │
│  Đã bàn giao: 0đ                    │
│  Tổng thu: 0đ                       │
├─────────────────────────────────────┤
│  [Kết thúc ca]                      │
└─────────────────────────────────────┘
```

**❌ Không có nút bàn giao**

---

### 2. Có Tiền (Tiền hiện có > 0)

```
┌─────────────────────────────────────┐
│  Ca đang mở                         │
│  ☀️ Ca sáng                         │
├─────────────────────────────────────┤
│  Tiền hiện có: 500,000đ             │
│  Đã bàn giao: 0đ                    │
│  Tổng thu: 500,000đ                 │
├─────────────────────────────────────┤
│  [💰 Bàn giao một phần]             │
│  [🏁 Bàn giao và đóng ca]           │
└─────────────────────────────────────┘
```

**✅ Có nút bàn giao**

---

### 3. Có Bàn Giao Đang Chờ

```
┌─────────────────────────────────────┐
│  Ca đang mở                         │
│  ☀️ Ca sáng                         │
├─────────────────────────────────────┤
│  🕐 Đang chờ xác nhận bàn giao      │
│  500,000đ                   [Hủy]  │
├─────────────────────────────────────┤
│  [Chờ cashier xác nhận...]          │
└─────────────────────────────────────┘
```

**❌ Nút bàn giao bị ẩn**

---

## 🐛 Common Issues

### Issue 1: "Không thấy card Ca đang mở"

**Nguyên nhân:** Chưa mở ca

**Giải pháp:**
1. Scroll xuống
2. Tìm form "Mở ca làm việc"
3. Điền thông tin và mở ca

---

### Issue 2: "Chỉ thấy nút Kết thúc ca"

**Nguyên nhân:** Tiền hiện có = 0

**Giải pháp:**
1. Tạo order với cash payment
2. Hoặc mở ca mới với tiền đầu ca > 0

---

### Issue 3: "Nút bị disable/gray"

**Nguyên nhân:** Có bàn giao đang chờ

**Giải pháp:**
1. Chờ cashier xác nhận
2. Hoặc hủy bàn giao hiện tại

---

### Issue 4: "Code đã update nhưng UI không đổi"

**Nguyên nhân:** Browser cache

**Giải pháp:**
```bash
# 1. Rebuild frontend
cd frontend
npm run build

# 2. Hard refresh browser
Ctrl+Shift+R (Windows)
Cmd+Shift+R (Mac)

# 3. Clear browser cache
Settings → Clear browsing data → Cached images and files
```

---

## 🔍 Debug Commands

### Check User Role
```javascript
// In browser console
const user = JSON.parse(localStorage.getItem('user'))
console.log('Role:', user?.role)
// Expected: "waiter"
```

### Check Current Shift
```javascript
// In browser console
fetch('/api/shifts/current', {
  headers: {
    'Authorization': 'Bearer ' + localStorage.getItem('token')
  }
})
.then(r => r.json())
.then(data => console.log('Current shift:', data))
```

### Check Remaining Cash
```javascript
// In browser console
fetch('/api/shifts/current', {
  headers: {
    'Authorization': 'Bearer ' + localStorage.getItem('token')
  }
})
.then(r => r.json())
.then(data => console.log('Remaining cash:', data.remaining_cash))
// Must be > 0 to see handover buttons
```

---

## ✅ Success Criteria

Bạn sẽ thấy nút bàn giao khi:

1. ✅ Role = "waiter"
2. ✅ Có ca đang mở (status = "OPEN")
3. ✅ Tiền hiện có > 0
4. ✅ Không có bàn giao đang chờ
5. ✅ Frontend đã build mới nhất

---

## 📞 Still Not Working?

Nếu vẫn không thấy sau khi check tất cả:

1. **Check browser console for errors:**
   - F12 → Console tab
   - Look for red errors

2. **Check network requests:**
   - F12 → Network tab
   - Reload page
   - Check if `/api/shifts/current` returns data

3. **Check Vue DevTools:**
   - Install Vue DevTools extension
   - Check `isWaiter` computed property
   - Check `currentShift` data
   - Check `pendingHandover` data

4. **Provide debug info:**
   ```javascript
   // Copy this output
   console.log({
     role: JSON.parse(localStorage.getItem('user'))?.role,
     hasShift: !!currentShift.value,
     remainingCash: currentShift.value?.remaining_cash,
     hasPending: !!pendingHandover.value
   })
   ```

---

**Last Updated:** 2026-02-04  
**Version:** 1.0  
**Status:** ✅ Complete
