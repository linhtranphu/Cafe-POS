# Manager Shift Management Implementation

## 🎯 Mục Tiêu

Tạo menu "Quản lý ca" cho manager để xem và giám sát tất cả các ca làm việc của waiter/barista và cashier trong hệ thống.

## ✅ Đã Hoàn Thành

### 1. Navigation Update (`frontend/src/components/Navigation.vue`)

**Manager Navigation - 5 menus:**
1. 🏠 Dashboard
2. ⏰ **Quản lý ca** (NEW)
3. 👥 Nhân viên
4. 📊 Báo cáo
5. 👤 Cá nhân

### 2. New View: ManagerShiftView (`frontend/src/views/ManagerShiftView.vue`)

#### Features:

**📊 Stats Cards:**
- Tổng ca (tất cả shifts)
- Waiter đang làm (open waiter shifts)
- Cashier đang làm (open cashier shifts)
- Ca hôm nay (today's shifts)

**🔍 Filter Tabs:**
- Tất cả
- Đang mở (OPEN)
- Đã đóng (CLOSED)

**🍽️ Waiter/Barista Shifts Section:**
- List tất cả waiter/barista shifts
- Hiển thị:
  - Tên nhân viên
  - Role type (Phục vụ/Pha chế)
  - Shift type (Ca sáng/chiều/tối)
  - Status (Đang mở/Đã đóng)
  - Thời gian bắt đầu/kết thúc
  - Duration (nếu đang mở)
  - Stats: Tiền đầu ca, Tiền cuối ca, Doanh thu (nếu đã đóng)

**💵 Cashier Shifts Section:**
- List tất cả cashier shifts
- Hiển thị:
  - Tên thu ngân
  - Status (Đang mở/Đang đóng/Đã đóng)
  - Thời gian mở/đóng ca
  - Duration (nếu đang mở)
  - Stats: Tiền mặt thực tế, Tiền dự kiến, Chênh lệch (nếu đã đóng)

**📱 Shift Detail Modal:**
- Click vào shift để xem chi tiết
- Waiter shift details:
  - Thông tin nhân viên
  - Thời gian
  - Tài chính (nếu đã đóng)
- Cashier shift details:
  - Thông tin thu ngân
  - Thời gian
  - Tài chính và chênh lệch (nếu đã đóng)
  - Lý do chênh lệch (nếu có)

**🔄 Refresh Button:**
- Reload data từ server

### 3. Router Update (`frontend/src/router/index.js`)

Added new route:
```javascript
{
  path: '/manager/shifts',
  name: 'ManagerShifts',
  component: ManagerShiftView,
  meta: { requiresAuth: true, requiresManager: true }
}
```

### 4. Store Update (`frontend/src/stores/cashierShift.js`)

Added:
- `shifts` property (alias for `cashierShifts`)
- `fetchAllShifts()` method (alias for `fetchAllCashierShifts()`)

## 📱 UI Design

### Mobile-First Layout

```
┌─────────────────────────────────┐
│ ⏰ Quản lý ca làm việc      🔄  │
│                                 │
│ [Tất cả] [Đang mở] [Đã đóng]   │
└─────────────────────────────────┘

┌──────────┬──────────┐
│ ⏰ Tổng ca│ ✅ Waiter│
├──────────┼──────────┤
│ 💵 Cashier│ 📅 Hôm nay│
└──────────┴──────────┘

🍽️ Ca Waiter/Barista (X ca)
┌─────────────────────────────────┐
│ Nguyễn Văn A                    │
│ 🍽️ Phục vụ                      │
│ ☀️ Ca sáng                       │
│                                 │
│ Bắt đầu: 31/01/2026, 08:00     │
│ Thời gian: 2h 30m               │
│                                 │
│ [✅ Đang mở]                    │
└─────────────────────────────────┘

💵 Ca Thu ngân (Y ca)
┌─────────────────────────────────┐
│ Trần Thị B                      │
│ 💵 Thu ngân                     │
│                                 │
│ Bắt đầu: 31/01/2026, 08:00     │
│ Kết thúc: 31/01/2026, 17:00    │
│                                 │
│ Tiền mặt: 5.000.000 ₫          │
│ Dự kiến: 4.950.000 ₫           │
│ Chênh lệch: +50.000 ₫          │
│                                 │
│ [Đã đóng]                       │
└─────────────────────────────────┘
```

### Detail Modal

```
┌─────────────────────────────────┐
│ Chi tiết ca làm việc        ×   │
├─────────────────────────────────┤
│                                 │
│ ┌─────────────────────────────┐ │
│ │ Nguyễn Văn A                │ │
│ │ 🍽️ Phục vụ                  │ │
│ │ ☀️ Ca sáng                   │ │
│ └─────────────────────────────┘ │
│                                 │
│ Trạng thái: [Đã đóng]          │
│                                 │
│ Thời gian bắt đầu:             │
│ 31/01/2026, 08:00:00           │
│                                 │
│ Thời gian kết thúc:            │
│ 31/01/2026, 17:00:00           │
│                                 │
│ ┌──────────┬──────────┐        │
│ │ Tiền đầu │ Tiền cuối│        │
│ │ 500.000₫ │ 5.500.000₫│       │
│ └──────────┴──────────┘        │
│                                 │
│ ┌─────────────────────┐        │
│ │ Tổng doanh thu      │        │
│ │ 5.000.000 ₫         │        │
│ └─────────────────────┘        │
└─────────────────────────────────┘
```

## 🔐 Access Control

### Manager có thể:
- ✅ Xem tất cả waiter/barista shifts
- ✅ Xem tất cả cashier shifts
- ✅ Xem chi tiết từng shift
- ✅ Filter shifts theo status
- ✅ Refresh data

### Manager KHÔNG thể:
- ❌ Mở/đóng shift cho nhân viên (nhân viên tự quản lý)
- ❌ Chỉnh sửa shift data
- ❌ Xóa shifts

## 📊 Data Flow

```
ManagerShiftView
    ↓
    ├─→ shiftStore.fetchAllShifts()
    │   └─→ GET /api/manager/shifts
    │       └─→ Returns: waiter/barista shifts
    │
    └─→ cashierShiftStore.fetchAllShifts()
        └─→ GET /api/cashier-shifts
            └─→ Returns: cashier shifts
```

## 🎨 Color Coding

### Shift Status:
- **Đang mở** (OPEN): Green (`bg-green-100 text-green-800`)
- **Đã đóng** (CLOSED): Gray (`bg-gray-100 text-gray-800`)

### Cashier Shift Status:
- **Đang mở** (OPEN): Green (`bg-green-100 text-green-800`)
- **Đang đóng** (CLOSURE_INITIATED): Yellow (`bg-yellow-100 text-yellow-800`)
- **Đã đóng** (CLOSED): Gray (`bg-gray-100 text-gray-800`)

### Role Types:
- **Waiter**: 🍽️ Phục vụ
- **Barista**: 🍹 Pha chế
- **Cashier**: 💵 Thu ngân

### Shift Types:
- **Morning**: ☀️ Ca sáng
- **Afternoon**: 🌤️ Ca chiều
- **Evening**: 🌙 Ca tối

## 🧪 Testing Checklist

### Navigation:
- [ ] Manager navigation shows 5 menus
- [ ] "Quản lý ca" menu is visible
- [ ] Clicking "Quản lý ca" navigates to `/manager/shifts`
- [ ] Non-manager users cannot access this route

### Data Loading:
- [ ] Stats cards show correct counts
- [ ] Waiter shifts load correctly
- [ ] Cashier shifts load correctly
- [ ] Refresh button reloads data

### Filtering:
- [ ] "Tất cả" shows all shifts
- [ ] "Đang mở" shows only open shifts
- [ ] "Đã đóng" shows only closed shifts
- [ ] Filter applies to both waiter and cashier shifts

### Shift Display:
- [ ] Waiter shifts show correct info
- [ ] Cashier shifts show correct info
- [ ] Duration calculates correctly for open shifts
- [ ] Stats show correctly for closed shifts

### Detail Modal:
- [ ] Clicking shift opens detail modal
- [ ] Waiter shift details display correctly
- [ ] Cashier shift details display correctly
- [ ] Close button works
- [ ] Modal scrolls if content is long

### Mobile Responsiveness:
- [ ] Layout works on mobile
- [ ] Cards are touch-friendly
- [ ] Modal slides up from bottom
- [ ] Filter tabs scroll horizontally if needed

## 📝 Files Created/Modified

### Created:
1. `frontend/src/views/ManagerShiftView.vue` - Main view component

### Modified:
1. `frontend/src/components/Navigation.vue` - Added "Quản lý ca" menu
2. `frontend/src/router/index.js` - Added `/manager/shifts` route
3. `frontend/src/stores/cashierShift.js` - Added `shifts` property and `fetchAllShifts()` method

## 🚀 Benefits

1. **Centralized Monitoring**: Manager có thể xem tất cả shifts ở một nơi
2. **Real-time Status**: Biết được nhân viên nào đang làm việc
3. **Financial Oversight**: Xem doanh thu và chênh lệch của từng ca
4. **Easy Filtering**: Nhanh chóng filter theo status
5. **Detailed View**: Xem chi tiết từng shift khi cần

## 🔄 Future Enhancements

- [ ] Add date range filter
- [ ] Add search by employee name
- [ ] Add export to Excel/PDF
- [ ] Add shift statistics/charts
- [ ] Add ability to add notes to shifts
- [ ] Add shift comparison
- [ ] Add notifications for long shifts
- [ ] Add shift scheduling (future shifts)

## 📚 Related Documents

- `MANAGER_NAVIGATION_REDESIGN.md` - Manager navigation design
- `CASHIER_WAITER_SHIFT_SEPARATION_PLAN.md` - Shift separation architecture
- `STATE_MACHINE_DOCUMENTATION.md` - Shift state machine
