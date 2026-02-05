# 👨‍💼 Hướng Dẫn Bàn Giao Tiền - Vai Trò Waiter

## 📍 Vị Trí Giao Diện

**Giao diện bàn giao cho waiter được tích hợp trong trang quản lý ca:**

- **Route:** `/shift` hoặc `/waiter`
- **Component:** `ShiftView.vue`
- **Không có trang riêng** - Tất cả chức năng bàn giao nằm trong trang ca làm việc

---

## 🚀 Cách Truy Cập

### Bước 1: Login as Waiter
```
1. Mở app
2. Login với tài khoản waiter
   - Username: waiter1 (hoặc tài khoản waiter khác)
   - Password: password123
```

### Bước 2: Vào Trang Ca Làm Việc
```
1. Click vào bottom navigation: "⏰ Ca làm"
   HOẶC
2. Navigate đến: /shift
```

### Bước 3: Mở Ca (Nếu Chưa Có Ca Mở)
```
1. Chọn ca: ☀️ Ca sáng / 🌤️ Ca chiều / 🌙 Ca tối
2. Nhập tiền đầu ca: 0 (hoặc số tiền bạn có)
3. Click "Mở ca"
```

### Bước 4: Có Tiền Để Bàn Giao
```
Để thấy nút bàn giao, cần có:
- Tiền hiện có > 0
- Hoặc đã thu tiền từ orders
```

---

## 🎨 Giao Diện Waiter

### 1. Card Ca Đang Mở

```
┌─────────────────────────────────────────┐
│  Ca đang mở                             │
│  ☀️ Ca sáng                             │
│  Phục vụ                                │
│                          [ĐANG MỞ]      │
├─────────────────────────────────────────┤
│  Bắt đầu: 09:00        Tiền đầu ca: 0đ  │
├─────────────────────────────────────────┤
│  Tiền hiện có  │ Đã bàn giao │ Tổng thu │
│  500,000đ      │ 0đ          │ 500,000đ │
├─────────────────────────────────────────┤
│  [💰 Bàn giao một phần]                 │
│  [🏁 Bàn giao và đóng ca]               │
└─────────────────────────────────────────┘
```

---

### 2. Nút Bàn Giao

#### A. Bàn Giao Một Phần
**Hiển thị khi:**
- ✅ Có tiền hiện có > 0
- ✅ Không có bàn giao đang chờ
- ✅ Ca đang mở

**Button:**
```
┌─────────────────────────────────────┐
│  💰 Bàn giao một phần               │
└─────────────────────────────────────┘
```

**Chức năng:**
- Bàn giao một phần tiền
- Ca vẫn tiếp tục
- Có thể bàn giao nhiều lần

---

#### B. Bàn Giao và Đóng Ca
**Hiển thị khi:**
- ✅ Có tiền hiện có > 0
- ✅ Không có bàn giao đang chờ
- ✅ Ca đang mở

**Button:**
```
┌─────────────────────────────────────┐
│  🏁 Bàn giao và đóng ca             │
└─────────────────────────────────────┘
```

**Chức năng:**
- Bàn giao toàn bộ tiền
- Kết thúc ca làm việc
- Chờ cashier xác nhận

---

### 3. Trạng Thái Đang Chờ

**Khi có bàn giao đang chờ:**
```
┌─────────────────────────────────────┐
│  🕐 Đang chờ xác nhận bàn giao      │
│  500,000đ                           │
│  Bàn giao một phần          [Hủy]  │
└─────────────────────────────────────┘

┌─────────────────────────────────────┐
│  Chờ cashier xác nhận...            │
└─────────────────────────────────────┘
```

**Các nút bị disable:**
- ❌ Không thể bàn giao thêm
- ❌ Không thể đóng ca
- ✅ Có thể hủy bàn giao

---

### 4. Lịch Sử Bàn Giao

```
┌─────────────────────────────────────┐
│  📋 Lịch sử bàn giao                │
├─────────────────────────────────────┤
│  500,000đ              [✅ Đã xác nhận]│
│  14:30 - 04/02/2026                 │
│  Bàn giao một phần                  │
│  Ghi chú: Bàn giao tiền ca chiều    │
│  Phản hồi cashier: OK               │
├─────────────────────────────────────┤
│  300,000đ              [⏳ Đang chờ] │
│  10:15 - 04/02/2026                 │
│  Bàn giao một phần                  │
└─────────────────────────────────────┘
```

---

## 🔄 Workflow Bàn Giao

### Flow 1: Bàn Giao Một Phần

```
1. Waiter mở ca
   ↓
2. Thu tiền từ khách (qua orders)
   ↓
3. Tiền hiện có: 500,000đ
   ↓
4. Click "💰 Bàn giao một phần"
   ↓
5. Modal hiện ra:
   ┌─────────────────────────────────┐
   │  💰 Bàn giao tiền               │
   ├─────────────────────────────────┤
   │  Số tiền: [_________] VND       │
   │  Ghi chú: [_______________]     │
   │                                 │
   │  [Hủy]  [Xác nhận bàn giao] ✅  │
   └─────────────────────────────────┘
   ↓
6. Nhập số tiền: 300,000
   ↓
7. Nhập ghi chú: "Bàn giao tiền ca chiều"
   ↓
8. Click "Xác nhận bàn giao"
   ↓
9. Thành công! Hiển thị trạng thái "Đang chờ"
   ↓
10. Tiền hiện có còn: 200,000đ
    ↓
11. Ca vẫn tiếp tục
```

---

### Flow 2: Bàn Giao và Đóng Ca

```
1. Waiter muốn kết thúc ca
   ↓
2. Tiền hiện có: 500,000đ
   ↓
3. Click "🏁 Bàn giao và đóng ca"
   ↓
4. Modal hiện ra:
   ┌─────────────────────────────────┐
   │  🏁 Bàn giao & Kết thúc ca      │
   ├─────────────────────────────────┤
   │  Tổng tiền thu: 1,500,000 VND   │
   │  Số tiền bàn giao: [_____] VND  │
   │  Ghi chú: [_______________]     │
   │                                 │
   │  ⚠️ Ca sẽ kết thúc sau khi      │
   │  cashier xác nhận               │
   │                                 │
   │  [Hủy]  [Bàn giao & Kết thúc]  │
   └─────────────────────────────────┘
   ↓
5. Nhập số tiền: 500,000
   ↓
6. Click "Bàn giao & Kết thúc"
   ↓
7. Ca chuyển sang "ENDING"
   ↓
8. Chờ cashier xác nhận
   ↓
9. Cashier xác nhận → Ca kết thúc
```

---

## 🎯 Các Trường Hợp Hiển Thị

### Case 1: Không Có Tiền
**Điều kiện:**
- Tiền hiện có = 0
- Chưa thu tiền từ orders

**Hiển thị:**
```
┌─────────────────────────────────────┐
│  [Kết thúc ca]                      │
└─────────────────────────────────────┘
```

**Không hiển thị:**
- ❌ Bàn giao một phần
- ❌ Bàn giao và đóng ca

---

### Case 2: Có Tiền, Không Có Bàn Giao Đang Chờ
**Điều kiện:**
- Tiền hiện có > 0
- Không có pending handover

**Hiển thị:**
```
┌─────────────────────────────────────┐
│  [💰 Bàn giao một phần]             │
│  [🏁 Bàn giao và đóng ca]           │
└─────────────────────────────────────┘
```

---

### Case 3: Có Bàn Giao Đang Chờ
**Điều kiện:**
- Có pending handover

**Hiển thị:**
```
┌─────────────────────────────────────┐
│  🕐 Đang chờ xác nhận bàn giao      │
│  500,000đ                   [Hủy]  │
├─────────────────────────────────────┤
│  [Chờ cashier xác nhận...]          │
└─────────────────────────────────────┘
```

**Không hiển thị:**
- ❌ Bàn giao một phần
- ❌ Bàn giao và đóng ca
- ❌ Kết thúc ca

---

## 🐛 Troubleshooting

### Vấn đề 1: Không Thấy Nút Bàn Giao

**Nguyên nhân:**
- Chưa login as waiter
- Chưa mở ca
- Tiền hiện có = 0
- Có bàn giao đang chờ

**Giải pháp:**
1. Kiểm tra role: Phải là "waiter"
2. Kiểm tra ca: Phải có ca đang mở
3. Kiểm tra tiền: Phải có tiền > 0
4. Kiểm tra pending: Không có bàn giao đang chờ

---

### Vấn đề 2: Nút Bị Disable

**Nguyên nhân:**
- Có bàn giao đang chờ xác nhận

**Giải pháp:**
- Chờ cashier xác nhận
- Hoặc hủy bàn giao hiện tại

---

### Vấn đề 3: Không Thấy Lịch Sử

**Nguyên nhân:**
- Chưa có bàn giao nào
- Chưa load data

**Giải pháp:**
- Tạo bàn giao mới
- Refresh trang

---

## 📊 Data Requirements

### Để Thấy UI Bàn Giao

**Minimum Requirements:**
```javascript
{
  user: {
    role: "waiter"  // ✅ Required
  },
  currentShift: {
    id: "...",
    status: "OPEN",  // ✅ Required
    role_type: "waiter",
    remaining_cash: 500000,  // ✅ Must be > 0
    current_cash: 500000
  },
  pendingHandover: null  // ✅ Must be null to show buttons
}
```

---

## 🎨 UI States

### State 1: Ready to Handover
```
✅ Role: waiter
✅ Shift: OPEN
✅ Cash: > 0
✅ Pending: null

→ Show: Bàn giao buttons
```

### State 2: Waiting for Confirmation
```
✅ Role: waiter
✅ Shift: OPEN or ENDING
✅ Cash: any
✅ Pending: exists

→ Show: Pending status + Cancel button
→ Hide: Bàn giao buttons
```

### State 3: No Cash
```
✅ Role: waiter
✅ Shift: OPEN
❌ Cash: 0
✅ Pending: null

→ Show: End shift button only
→ Hide: Bàn giao buttons
```

---

## 🔍 Debug Checklist

Nếu không thấy UI bàn giao, kiểm tra console:

```javascript
// 1. Check user role
console.log('User role:', authStore.user?.role)
// Expected: "waiter"

// 2. Check current shift
console.log('Current shift:', currentShift.value)
// Expected: { id: "...", status: "OPEN", ... }

// 3. Check cash amount
console.log('Remaining cash:', currentShift.value?.remaining_cash)
// Expected: > 0

// 4. Check pending handover
console.log('Pending handover:', pendingHandover.value)
// Expected: null (to show buttons) or object (to show pending status)

// 5. Check isWaiter computed
console.log('Is waiter:', isWaiter.value)
// Expected: true
```

---

## 📸 Screenshots Reference

### Desktop View
```
/shift page → Scroll down → See handover buttons
```

### Mobile View
```
Bottom nav → "⏰ Ca làm" → See handover buttons in shift card
```

---

## 🔗 Related Files

### Frontend
- `frontend/src/views/ShiftView.vue` - Main UI
- `frontend/src/stores/shift.js` - State management
- `frontend/src/services/handover.js` - API calls

### Backend
- `backend/interfaces/http/cash_handover_handler.go` - Handlers
- `backend/application/services/cash_handover_service.go` - Business logic

### Documentation
- [CASH_HANDOVER_UI_GUIDE.md](./CASH_HANDOVER_UI_GUIDE.md) - Complete UI guide
- [CASH_HANDOVER_USER_GUIDE.md](./CASH_HANDOVER_USER_GUIDE.md) - User guide (Vietnamese)

---

## ✅ Quick Test

### Test Steps:
1. Login as waiter (username: waiter1, password: password123)
2. Go to `/shift`
3. If no shift → Start shift (any type, start_cash: 0)
4. Create an order with cash payment (to have cash > 0)
5. Go back to `/shift`
6. Should see: "💰 Bàn giao một phần" button
7. Click button → Modal appears
8. Enter amount → Submit
9. Should see: "🕐 Đang chờ xác nhận bàn giao"

---

**Last Updated:** 2026-02-04  
**Version:** 1.0  
**Status:** ✅ Complete
