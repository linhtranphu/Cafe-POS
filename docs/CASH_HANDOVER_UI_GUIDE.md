# 🎨 Hướng Dẫn Giao Diện Bàn Giao Tiền

## 📱 Tổng Quan UI

Hệ thống bàn giao tiền có **3 giao diện chính**:

1. **Waiter Interface** - Giao diện phục vụ (ShiftView.vue)
2. **Cashier Interface** - Giao diện thu ngân (CashierHandoverView.vue)
3. **Manager Interface** - Giao diện quản lý (CashierHandoverView.vue)

---

## 1️⃣ Giao Diện Phục Vụ (Waiter)

### 📍 Vị trí: `/shift` hoặc `/waiter`

### 🎯 Chức năng chính:

#### A. Bàn Giao Một Phần (Partial Handover)
**Khi nào dùng:** Bàn giao tiền trong ca, ca vẫn tiếp tục

**UI Components:**
```
┌─────────────────────────────────────┐
│  💰 Bàn Giao Tiền                   │
├─────────────────────────────────────┤
│  Số tiền bàn giao: [_________] VND  │
│  Ghi chú: [___________________]     │
│                                     │
│  [Hủy]  [Xác nhận bàn giao] ✅      │
└─────────────────────────────────────┘
```

**Các trường:**
- ✅ Số tiền bàn giao (required)
- ✅ Ghi chú (optional)
- ✅ Tự động gửi đến thu ngân đang mở ca

**Flow:**
1. Click nút "💰 Bàn giao tiền"
2. Nhập số tiền muốn bàn giao
3. Thêm ghi chú (nếu cần)
4. Click "Xác nhận bàn giao"
5. ✅ Thành công → Hiển thị thông báo
6. ⏳ Chờ thu ngân xác nhận

---

#### B. Bàn Giao & Kết Thúc Ca (Handover & End Shift)
**Khi nào dùng:** Bàn giao tiền và kết thúc ca làm việc

**UI Components:**
```
┌─────────────────────────────────────┐
│  🏁 Bàn Giao & Kết Thúc Ca          │
├─────────────────────────────────────┤
│  Tổng tiền thu: 1,500,000 VND       │
│  Số tiền bàn giao: [_________] VND  │
│  Ghi chú: [___________________]     │
│                                     │
│  ⚠️ Lưu ý: Ca sẽ kết thúc sau khi   │
│  thu ngân xác nhận bàn giao         │
│                                     │
│  [Hủy]  [Bàn giao & Kết thúc] 🏁   │
└─────────────────────────────────────┘
```

**Các trường:**
- ✅ Tổng tiền thu (hiển thị tự động)
- ✅ Số tiền bàn giao (required)
- ✅ Ghi chú (optional)

**Flow:**
1. Click nút "🏁 Kết thúc ca"
2. Hệ thống hiển thị tổng tiền thu trong ca
3. Nhập số tiền bàn giao
4. Thêm ghi chú (nếu cần)
5. Click "Bàn giao & Kết thúc"
6. ⏳ Ca chuyển sang trạng thái "ENDING"
7. ⏳ Chờ thu ngân xác nhận
8. ✅ Thu ngân xác nhận → Ca kết thúc

---

#### C. Trạng Thái Bàn Giao

**Pending (Đang chờ):**
```
┌─────────────────────────────────────┐
│  ⏳ Bàn giao đang chờ xác nhận       │
├─────────────────────────────────────┤
│  Số tiền: 500,000 VND               │
│  Thu ngân: Nguyễn Văn A             │
│  Thời gian: 14:30 - 04/02/2026      │
│  Ghi chú: Bàn giao tiền ca chiều    │
│                                     │
│  [Hủy bàn giao] ❌                  │
└─────────────────────────────────────┘
```

**Confirmed (Đã xác nhận):**
```
┌─────────────────────────────────────┐
│  ✅ Bàn giao đã xác nhận             │
├─────────────────────────────────────┤
│  Số tiền: 500,000 VND               │
│  Thu ngân: Nguyễn Văn A             │
│  Xác nhận lúc: 14:35 - 04/02/2026   │
│  Trạng thái: ✅ Chính xác           │
└─────────────────────────────────────┘
```

**Rejected (Bị từ chối):**
```
┌─────────────────────────────────────┐
│  ❌ Bàn giao bị từ chối              │
├─────────────────────────────────────┤
│  Số tiền: 500,000 VND               │
│  Thu ngân: Nguyễn Văn A             │
│  Lý do: Thiếu 50,000 VND            │
│  Thời gian: 14:35 - 04/02/2026      │
└─────────────────────────────────────┘
```

---

#### D. Lịch Sử Bàn Giao

**UI Components:**
```
┌─────────────────────────────────────┐
│  📜 Lịch sử bàn giao                │
├─────────────────────────────────────┤
│  ✅ 500,000 VND - 14:30             │
│     Thu ngân: Nguyễn Văn A          │
│     Trạng thái: Đã xác nhận         │
├─────────────────────────────────────┤
│  ⏳ 300,000 VND - 10:15             │
│     Thu ngân: Trần Thị B            │
│     Trạng thái: Đang chờ            │
├─────────────────────────────────────┤
│  ❌ 200,000 VND - 09:00             │
│     Thu ngân: Nguyễn Văn A          │
│     Trạng thái: Bị từ chối          │
│     Lý do: Thiếu tiền               │
└─────────────────────────────────────┘
```

---

## 2️⃣ Giao Diện Thu Ngân (Cashier)

### 📍 Vị trí: `/cashier/handovers`

### 🎯 Chức năng chính:

#### A. Dashboard Overview (Trong CashierDashboard.vue)

**Notification Banner:**
```
┌─────────────────────────────────────┐
│  ⚠️ 3 yêu cầu bàn giao đang chờ     │
│                                     │
│  [Xem ngay] →                       │
└─────────────────────────────────────┘
```

**Quick Actions:**
```
┌─────────────────────────────────────┐
│  ⚡ Bàn giao nhanh                   │
├─────────────────────────────────────┤
│  Nguyễn Văn A - 500,000 VND         │
│  [✅ Xác nhận]  [❌ Từ chối]        │
├─────────────────────────────────────┤
│  Trần Thị B - 300,000 VND           │
│  [✅ Xác nhận]  [❌ Từ chối]        │
├─────────────────────────────────────┤
│  [Xem tất cả 3 yêu cầu →]           │
└─────────────────────────────────────┘
```

---

#### B. Trang Quản Lý Bàn Giao Chi Tiết

**Tab Navigation:**
```
┌─────────────────────────────────────┐
│  [Đang chờ (3)] [Hôm nay (12)]      │
└─────────────────────────────────────┘
```

**Tab 1: Đang Chờ Xác Nhận**
```
┌─────────────────────────────────────┐
│  ⏳ Yêu cầu bàn giao #1              │
├─────────────────────────────────────┤
│  Phục vụ: Nguyễn Văn A              │
│  Số tiền khai báo: 500,000 VND      │
│  Thời gian: 14:30 - 04/02/2026      │
│  Ghi chú: Bàn giao tiền ca chiều    │
│                                     │
│  [Xác nhận chi tiết] [Từ chối] ❌   │
└─────────────────────────────────────┘
```

**Tab 2: Bàn Giao Hôm Nay**
```
┌─────────────────────────────────────┐
│  📊 Tổng quan hôm nay               │
├─────────────────────────────────────┤
│  Tổng số bàn giao: 12               │
│  Đã xác nhận: 10 ✅                 │
│  Bị từ chối: 2 ❌                   │
│  Tổng tiền: 5,500,000 VND           │
└─────────────────────────────────────┘

┌─────────────────────────────────────┐
│  ✅ Bàn giao #12 - 16:45            │
│  Phục vụ: Nguyễn Văn A              │
│  Số tiền: 500,000 VND               │
│  Trạng thái: Đã xác nhận            │
├─────────────────────────────────────┤
│  ❌ Bàn giao #11 - 16:30            │
│  Phục vụ: Trần Thị B                │
│  Số tiền: 300,000 VND               │
│  Trạng thái: Bị từ chối             │
│  Lý do: Thiếu 50,000 VND            │
└─────────────────────────────────────┘
```

---

#### C. Modal Xác Nhận Chi Tiết

**Reconciliation Modal:**
```
┌─────────────────────────────────────┐
│  💰 Xác nhận bàn giao               │
├─────────────────────────────────────┤
│  Phục vụ: Nguyễn Văn A              │
│  Số tiền khai báo: 500,000 VND      │
│                                     │
│  Số tiền thực tế: [_________] VND   │
│                                     │
│  Chênh lệch: 0 VND ✅               │
│                                     │
│  Ghi chú thu ngân:                  │
│  [_________________________]        │
│                                     │
│  [Hủy]  [Xác nhận] ✅               │
└─────────────────────────────────────┘
```

**Khi có chênh lệch > 100,000 VND:**
```
┌─────────────────────────────────────┐
│  ⚠️ Phát hiện chênh lệch lớn        │
├─────────────────────────────────────┤
│  Số tiền khai báo: 500,000 VND      │
│  Số tiền thực tế: 350,000 VND       │
│  Chênh lệch: -150,000 VND ⚠️        │
│                                     │
│  ⚠️ Chênh lệch vượt ngưỡng cho phép │
│  (100,000 VND). Cần quản lý phê     │
│  duyệt.                             │
│                                     │
│  Ghi chú thu ngân:                  │
│  [Thiếu 150k, cần kiểm tra]         │
│                                     │
│  [Hủy]  [Gửi quản lý phê duyệt] 📤  │
└─────────────────────────────────────┘
```

---

## 3️⃣ Giao Diện Quản Lý (Manager)

### 📍 Vị trí: `/cashier/handovers` (với quyền manager)

### 🎯 Chức năng chính:

#### A. Tab Cần Phê Duyệt

**Tab Navigation:**
```
┌─────────────────────────────────────┐
│  [Đang chờ] [Hôm nay] [Cần duyệt (2)]│
└─────────────────────────────────────┘
```

**Danh sách cần phê duyệt:**
```
┌─────────────────────────────────────┐
│  ⚠️ Chênh lệch cần phê duyệt #1     │
├─────────────────────────────────────┤
│  Phục vụ: Nguyễn Văn A              │
│  Thu ngân: Trần Thị B               │
│  Số tiền khai báo: 500,000 VND      │
│  Số tiền thực tế: 350,000 VND       │
│  Chênh lệch: -150,000 VND ⚠️        │
│  Ghi chú TN: Thiếu 150k             │
│                                     │
│  [Chi tiết] [Phê duyệt] [Từ chối]  │
└─────────────────────────────────────┘
```

---

#### B. Modal Phê Duyệt Chênh Lệch

```
┌─────────────────────────────────────┐
│  👨‍💼 Phê duyệt chênh lệch            │
├─────────────────────────────────────┤
│  📋 Thông tin bàn giao              │
│  Phục vụ: Nguyễn Văn A              │
│  Thu ngân: Trần Thị B               │
│  Thời gian: 14:30 - 04/02/2026      │
│                                     │
│  💰 Chi tiết tiền                   │
│  Khai báo: 500,000 VND              │
│  Thực tế: 350,000 VND               │
│  Chênh lệch: -150,000 VND ⚠️        │
│                                     │
│  📝 Ghi chú                         │
│  Phục vụ: Bàn giao tiền ca chiều    │
│  Thu ngân: Thiếu 150k, cần kiểm tra │
│                                     │
│  Quyết định:                        │
│  ○ Phê duyệt (chấp nhận chênh lệch) │
│  ○ Từ chối (yêu cầu kiểm tra lại)   │
│                                     │
│  Ghi chú quản lý:                   │
│  [_________________________]        │
│                                     │
│  [Hủy]  [Xác nhận quyết định] ✅    │
└─────────────────────────────────────┘
```

---

#### C. Thống Kê Chênh Lệch

```
┌─────────────────────────────────────┐
│  📊 Thống kê chênh lệch             │
├─────────────────────────────────────┤
│  Khoảng thời gian: [____] - [____]  │
│  [Xem thống kê]                     │
└─────────────────────────────────────┘

┌─────────────────────────────────────┐
│  📈 Kết quả thống kê                │
├─────────────────────────────────────┤
│  Tổng số bàn giao: 150              │
│  Có chênh lệch: 12 (8%)             │
│  Đã phê duyệt: 10                   │
│  Bị từ chối: 2                      │
│                                     │
│  Tổng chênh lệch: -500,000 VND      │
│  Chênh lệch TB: -41,667 VND         │
│  Chênh lệch lớn nhất: -200,000 VND  │
└─────────────────────────────────────┘
```

---

## 🎨 Design System

### Colors

**Status Colors:**
- 🟢 Green (`bg-green-500`) - Confirmed, Success
- 🔴 Red (`bg-red-500`) - Rejected, Error
- 🟡 Yellow (`bg-yellow-500`) - Pending, Warning
- 🔵 Blue (`bg-blue-500`) - Info, Primary action
- ⚪ Gray (`bg-gray-500`) - Neutral, Disabled

**Background Colors:**
- `bg-white` - Card background
- `bg-gray-50` - Secondary background
- `bg-yellow-50` - Warning background
- `bg-green-50` - Success background
- `bg-red-50` - Error background

### Typography

**Font Sizes:**
- `text-xs` (12px) - Small labels
- `text-sm` (14px) - Body text
- `text-base` (16px) - Default text
- `text-lg` (18px) - Section headers
- `text-xl` (20px) - Page headers
- `text-2xl` (24px) - Main headers

**Font Weights:**
- `font-normal` (400) - Body text
- `font-medium` (500) - Emphasized text
- `font-bold` (700) - Headers

### Spacing

**Padding:**
- `p-2` (8px) - Tight spacing
- `p-4` (16px) - Default spacing
- `p-6` (24px) - Comfortable spacing

**Margin:**
- `mb-2` (8px) - Small gap
- `mb-4` (16px) - Default gap
- `mb-6` (24px) - Large gap

### Borders

**Border Radius:**
- `rounded-lg` (8px) - Buttons, inputs
- `rounded-xl` (12px) - Cards
- `rounded-2xl` (16px) - Large cards

**Border Width:**
- `border` (1px) - Default border
- `border-2` (2px) - Emphasized border
- `border-l-4` (4px left) - Notification banner

---

## 📱 Responsive Design

### Mobile (< 768px)
- Single column layout
- Full width cards
- Stacked buttons
- Simplified navigation

### Tablet (768px - 1024px)
- Two column layout
- Compact cards
- Side-by-side buttons

### Desktop (> 1024px)
- Three column layout
- Detailed cards
- Inline actions
- Full navigation

---

## 🔔 Notifications & Feedback

### Success Messages
```
✅ Bàn giao thành công!
✅ Đã xác nhận bàn giao
✅ Đã phê duyệt chênh lệch
```

### Error Messages
```
❌ Không thể tạo bàn giao
❌ Số tiền không hợp lệ
❌ Không tìm thấy thu ngân
```

### Warning Messages
```
⚠️ Chênh lệch vượt ngưỡng cho phép
⚠️ Cần quản lý phê duyệt
⚠️ Ca sẽ kết thúc sau khi xác nhận
```

### Info Messages
```
ℹ️ Đang chờ thu ngân xác nhận
ℹ️ Bàn giao đã được gửi
ℹ️ Đang tải dữ liệu...
```

---

## 🚀 User Flows

### Flow 1: Bàn Giao Một Phần (Waiter)
1. Waiter mở trang `/shift`
2. Click "💰 Bàn giao tiền"
3. Nhập số tiền và ghi chú
4. Click "Xác nhận bàn giao"
5. Thấy notification "✅ Bàn giao thành công"
6. Thấy trạng thái "⏳ Đang chờ xác nhận"

### Flow 2: Xác Nhận Bàn Giao (Cashier)
1. Cashier thấy notification "⚠️ 1 yêu cầu bàn giao"
2. Click "Xem ngay" → Đến `/cashier/handovers`
3. Click "Xác nhận chi tiết"
4. Nhập số tiền thực tế
5. Kiểm tra chênh lệch
6. Click "Xác nhận"
7. Thấy notification "✅ Đã xác nhận"

### Flow 3: Phê Duyệt Chênh Lệch (Manager)
1. Manager vào tab "Cần duyệt"
2. Thấy danh sách chênh lệch
3. Click "Chi tiết"
4. Xem thông tin đầy đủ
5. Chọn "Phê duyệt" hoặc "Từ chối"
6. Nhập ghi chú
7. Click "Xác nhận quyết định"
8. Thấy notification "✅ Đã phê duyệt"

---

## 📸 Screenshots Locations

**Files to check:**
- `frontend/src/views/ShiftView.vue` - Waiter interface
- `frontend/src/views/CashierHandoverView.vue` - Cashier/Manager interface
- `frontend/src/views/CashierDashboard.vue` - Dashboard notifications

**To see the UI:**
1. Start the application
2. Login as waiter → Go to `/shift`
3. Login as cashier → Go to `/cashier/handovers`
4. Login as manager → Go to `/cashier/handovers`

---

## 🔗 Related Documentation

- [CASH_HANDOVER_README.md](./CASH_HANDOVER_README.md) - Tổng quan tính năng
- [CASH_HANDOVER_API_DOCUMENTATION.md](./CASH_HANDOVER_API_DOCUMENTATION.md) - API endpoints
- [CASH_HANDOVER_USER_GUIDE.md](./CASH_HANDOVER_USER_GUIDE.md) - Hướng dẫn sử dụng
- [CASH_HANDOVER_PHASE6_COMPLETE.md](./CASH_HANDOVER_PHASE6_COMPLETE.md) - Implementation details

---

**Last Updated:** 2026-02-04  
**Version:** 1.0  
**Status:** ✅ Complete
