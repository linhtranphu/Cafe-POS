# 🏗️ Café POS System - Requirements Document

## 📋 Core Requirements

### Order Management
- ✅ Tạo order mới, gắn order với bàn, CRUD món trong order
- ✅ Tính tổng tiền realtime
- ✅ State machine: CREATED → CONFIRMED → PAID → BILLED

### Menu & Pricing
- ✅ Danh sách món uống, giá bán
- 🔄 Size/option (future)

### Table Management
- ✅ Danh sách bàn, trạng thái: Empty → Occupied → Paid
- ✅ 1 order = 1 bàn tại 1 thời điểm

### Payment & Billing
- ✅ Cash/Transfer/QR, ghi nhận thời gian + phương thức
- ✅ PAID chỉ 1 lần, không rollback (trừ admin)
- ✅ In bill sau thanh toán, in lại bill (audit)

### User & Role
- ✅ JWT Authentication, Role-based Authorization

## 🏗️ Architecture

```
cafe-pos/
├── backend/                 # Go + Gin + MongoDB
│   ├── domain/             # order, menu, table, payment, ingredient, user, facility, expense
│   ├── application/services/
│   ├── infrastructure/mongodb/
│   └── interfaces/http/
├── frontend/               # Vue.js 3 + Pinia
│   ├── views/
│   ├── services/
│   ├── stores/
│   └── components/
```

## 📊 Database Collections
**Core:** orders, menu_items, tables, payments, users, bills, shifts, cash_handovers
**Inventory:** ingredients, stock_history
**Facility:** facilities, maintenance_records
**Expense:** expenses, expense_categories, recurring_expenses, prepaid_expenses

## 📝 Functional Requirements

### 🍽️ WAITER

**FR-W-01: Mở ca làm việc**
- Waiter đăng nhập
- Chọn ca (sáng / chiều / tối)
- Gắn với: Ca, Quầy, Cashier phụ trách ca

**FR-W-02: Kết thúc ca**
- Gửi yêu cầu chốt ca
- Không được tự đóng ca nếu chưa đối soát tiền

**FR-W-03: Tạo order cho bàn**
- Chọn bàn
- Thêm món
- Sửa số lượng
- Ghi chú món

**FR-W-04: Gửi order xuống quầy pha chế**
- Trạng thái: Draft → Sent → Preparing → Served

**FR-W-05: Gộp / tách bàn (nếu cho phép)**
- Có thể cấu hình: chỉ cho phép trước khi tính tiền

**FR-W-06: Tính tiền cho order**
- Xem tổng tiền
- Áp dụng: Khuyến mãi, Giảm giá (nếu được quyền)

**FR-W-07: Thu tiền từ khách**
- Chọn phương thức: Tiền mặt, Chuyển khoản, Ví điện tử
- Ghi nhận: Số tiền thu, Phương thức, Thời điểm, Người thu = waiter
- ⚠️ Waiter chỉ "ghi nhận thu tiền", không quản lý quỹ

**FR-W-08: In / gửi bill cho khách**

**FR-W-09: KHÔNG được**
- Sửa bill đã thanh toán
- Hủy bill đã thu tiền
- Xem báo cáo doanh thu tổng
- Chỉnh sửa giá bán / cost / nguyên liệu

**FR-W-10: Xem báo cáo theo ca của bản thân**
- Số bill
- Tổng tiền đã thu
- Theo phương thức thanh toán

---

### 💵 CASHIER

#### A. Quản lý ca (Shift Management)

**FR-CASH-01 – Mở ca**
- Cashier có thể mở ca làm việc
- Nhập:
  - Thời gian bắt đầu ca
  - Số tiền đầu ca (cash float)
- Hệ thống ghi nhận:
  - Cashier mở ca
  - Thời điểm mở ca

**FR-CASH-02 – Theo dõi trạng thái ca**
- Cashier có thể xem:
  - Tổng số order trong ca
  - Tổng tiền đã thu (theo từng phương thức)
  - Order chưa thanh toán
- Dữ liệu cập nhật real-time

**FR-CASH-03 – Chốt ca**
- Cashier có thể thực hiện chốt ca
- Khi chốt ca, hệ thống:
  - Tính tổng doanh thu theo ca
  - Phân loại theo:
    - Tiền mặt
    - Chuyển khoản
  - So sánh:
    - Tiền thực tế nhập vào
    - Tiền hệ thống ghi nhận

#### B. Quản lý thanh toán (Payment Control)

**FR-CASH-04 – Giám sát thanh toán**
- Cashier có thể xem danh sách order:
  - Paid
  - Unpaid
- Thấy rõ:
  - Ai thu tiền (Waiter nào)
  - Thời điểm thu
  - Phương thức thanh toán

**FR-CASH-05 – Xử lý sai lệch thanh toán**
- Cashier có thể:
  - Đánh dấu order có sai lệch
  - Ghi chú lý do (thiếu tiền, nhầm tiền, khách thiếu…)
- Order bị đánh dấu sẽ:
  - Không cho khóa ca nếu chưa xử lý

#### C. Đối soát & Audit

**FR-CASH-06 – Đối soát tiền mặt**
- Khi chốt ca, Cashier nhập:
  - Số tiền mặt thực tế
- Hệ thống tự động:
  - Tính chênh lệch
  - Ghi nhận trạng thái:
    - Khớp
    - Dư
    - Thiếu

**FR-CASH-07 – Đối soát chuyển khoản**
- Cashier có thể:
  - Xác nhận các giao dịch chuyển khoản
  - Đánh dấu giao dịch treo / nghi ngờ
- Các giao dịch chưa xác nhận:
  - Không được tính là hoàn tất ca

#### D. Can thiệp nghiệp vụ (Controlled Override)

**FR-CASH-08 – Hủy/điều chỉnh thanh toán**
- Cashier có quyền:
  - Hủy trạng thái paid trong trường hợp đặc biệt
- Bắt buộc:
  - Nhập lý do
  - Ghi log audit (ai – khi nào – lý do)

**FR-CASH-09 – Khóa order**
- Cashier có thể:
  - Khóa order đã hoàn tất
- Order bị khóa:
  - Không cho chỉnh sửa
  - Không cho hủy thanh toán

#### E. Báo cáo & bàn giao

**FR-CASH-10 – Báo cáo ca**
- Cashier có thể xuất báo cáo ca:
  - Tổng order
  - Tổng tiền
  - Chênh lệch
  - Danh sách order bất thường

**FR-CASH-11 – Bàn giao ca**
- Cashier có thể:
  - Bàn giao ca cho cashier khác
- Hệ thống ghi nhận:
  - Người bàn giao
  - Người nhận ca
  - Thời điểm

---

### 🥬 INGREDIENT MANAGEMENT

**FR-IM-01: Xem danh sách nguyên liệu**
- Tên nguyên liệu
- Loại nguyên liệu (category)
- Đơn vị chuẩn ISO (kg, g, L, ml, piece, box, pack)
- Số lượng tồn
- Trạng thái (Còn hàng / Sắp hết / Hết hàng)

**FR-IM-02: Thêm mới nguyên liệu**
- Tên nguyên liệu (duy nhất)
- Loại nguyên liệu (Cà phê, Trà, Sữa, Đường, Trái cây, Bánh, Khác)
- Đơn vị chuẩn
- Số lượng ban đầu (không âm)
- Giá mỗi đơn vị
- Nhà cung cấp

**FR-IM-03: Cập nhật thông tin nguyên liệu**
- Chỉnh sửa: Tên, Loại, Ngưỡng cảnh báo, Trạng thái
- Business rule: Không thay đổi đơn vị nếu đã có giao dịch

**FR-IM-04: Cập nhật tồn kho**
- Nhập thêm (mua hàng)
- Xuất hủy (hư hỏng)
- Điều chỉnh kiểm kê
- Ghi nhận: Loại điều chỉnh, Số lượng, Lý do, Người thực hiện, Thời gian

**FR-IM-05: Xem lịch sử biến động**
- Trước/sau điều chỉnh
- Loại giao dịch (adjustment, order, purchase, waste)
- Order liên quan (nếu có)
- Người thao tác
- UI hiển thị 50 records gần nhất

**FR-IM-06: Cảnh báo nguyên liệu sắp hết**
- Highlight trên danh sách
- Thông báo trong dashboard
- Low stock alert panel

**FR-IM-07: Tìm kiếm & lọc**
- Tìm theo tên
- Lọc theo loại
- Lọc theo trạng thái tồn kho

**FR-IM-08: Phân quyền truy cập**
- Manager: toàn quyền (CRUD + stock adjustment)
- Staff: chỉ xem (GET /waiter/ingredients)

**FR-IM-09: Quản lý danh mục nguyên liệu**
- Thêm danh mục mới
- Xóa danh mục (nếu chưa có nguyên liệu)
- Hiển thị số lượng nguyên liệu theo danh mục
- Business rule: Không cho xóa danh mục đã có nguyên liệu

---

### 🏢 FACILITY MANAGEMENT

**FR-FM-01: Xem danh sách cơ sở vật chất**
- Tên tài sản, Loại tài sản, Khu vực sử dụng
- Số lượng, Tình trạng, Ngày mua / đưa vào sử dụng

**FR-FM-02: Thêm mới cơ sở vật chất**
- Thông tin bắt buộc: Tên, Loại (bàn, máy móc, dụng cụ…), Số lượng, Khu vực, Ngày mua, Tình trạng ban đầu
- Business Rules: Tài sản đơn chiếc (máy) quản lý theo từng item, Dụng cụ tiêu hao (ly, thìa) quản lý theo số lượng

**FR-FM-03: Cập nhật thông tin cơ sở vật chất**
- Cho phép chỉnh sửa: Khu vực, Tình trạng, Số lượng (với tài sản nhóm), Ghi chú
- Không cho phép: Xóa tài sản đã có lịch sử bảo trì (chỉ được inactive)

**FR-FM-04: Quản lý tình trạng tài sản**
- Trạng thái: Đang sử dụng, Hỏng, Đang sửa, Ngừng sử dụng, Thanh lý

**FR-FM-05: Báo hư hỏng (Staff)**
- Thông tin: Tài sản, Mô tả sự cố, Mức độ ảnh hưởng, Hình ảnh (optional)

**FR-FM-06: Quản lý bảo trì / sửa chữa**
- Thông tin bảo trì: Tài sản liên quan, Loại bảo trì (định kỳ / phát sinh), Nội dung, Ngày thực hiện, Chi phí, Đơn vị sửa chữa

**FR-FM-07: Lịch sử tài sản**
- Bao gồm: Thay đổi trạng thái, Bảo trì, Di chuyển khu vực, Thanh lý

**FR-FM-08: Tìm kiếm & lọc**
- Theo tên, Theo loại, Theo khu vực, Theo tình trạng

**FR-FM-09: Phân quyền truy cập**
- Manager: toàn quyền
- Staff: Xem danh sách, Báo hư hỏng, Không chỉnh sửa tài sản

---

### 💰 EXPENSE MANAGEMENT

**FR-EX-01: Xem danh sách chi phí**
- Hiển thị: Tên chi phí, Nhóm chi phí, Số tiền, Tháng áp dụng

**FR-EX-02: Khai báo loại chi phí**
- Tạo danh mục: Thuê nhà, Điện, Nước, Internet
- Gắn: Cố định / biến đổi / một lần, Có định kỳ hay không

**FR-EX-03: Ghi nhận chi phí phát sinh**
- Nhập chi phí thủ công
- Đính kèm hóa đơn

**FR-EX-04: Chi phí định kỳ**
- Thiết lập: Chu kỳ (tháng), Số tiền dự kiến
- Hệ thống nhắc nhập chi phí hàng tháng

**FR-EX-05: Quản lý chi phí trả trước / phân bổ**
- Nhập: Tổng số tiền, Thời gian áp dụng
- Hệ thống tự phân bổ chi phí theo tháng

**FR-EX-06: Báo cáo chi phí**
- Tổng chi phí theo: Tháng, Loại
- So sánh: Thực tế vs dự kiến

**FR-EX-07: Phân quyền**
- Manager: toàn quyền
- Accountant: xem & xuất báo cáo

**Business Rules:**
- Một chi phí phải thuộc 1 loại chi phí
- Chi phí phân bổ: Không được chỉnh sửa kỳ đã chốt, Không cho xóa chi phí đã dùng cho báo cáo
- Chi phí định kỳ: Không auto ghi nhận nếu chưa xác nhận

---

### 📋 ORDER MANAGEMENT (IMPLEMENTED)

**FR-OM-01: Tạo order mới**
- Chọn bàn (table_id)
- Chọn món từ menu (items[])
- Ghi chú cho order
- Gắn tự động: waiter_id, shift_id
- Trạng thái ban đầu: CREATED

**FR-OM-02: Xác nhận order**
- Chuyển trạng thái: CREATED → UNPAID
- Validate: Order phải có items
- Chỉ waiter tạo order mới được xác nhận

**FR-OM-03: Thu tiền**
- Chọn phương thức: CASH, TRANSFER, QR
- Nhập số tiền, áp dụng discount (nếu có)
- Chuyển trạng thái: UNPAID → PAID
- Ghi nhận: payment_method, paid_at, collected_by

**FR-OM-04: Gửi pha chế**
- Chuyển trạng thái: PAID → IN_PROGRESS
- Business rule: Phải PAID trước khi gửi kitchen
- Ghi nhận: sent_to_kitchen_at

**FR-OM-05: Phục vụ order**
- Chuyển trạng thái: IN_PROGRESS → SERVED
- Ghi nhận: served_at

**FR-OM-06: Hủy order (Cashier/Manager)**
- Chuyển trạng thái: UNPAID → CANCELLED
- Ghi nhận: cancelled_at, cancelled_by, cancel_reason
- Business rule: Chỉ hủy được order UNPAID

**FR-OM-07: Hoàn tiền (Cashier/Manager)**
- Chuyển trạng thái: PAID/IN_PROGRESS → REFUNDED
- Ghi nhận: refunded_at, refunded_by, refund_reason
- Business rule: Không hoàn tiền order đã SERVED

**FR-OM-08: Khóa order (Cashier)**
- Chuyển trạng thái: SERVED/CANCELLED/REFUNDED → LOCKED
- Business rule: Order LOCKED không thể sửa/xóa
- Auto lock khi chốt ca

**FR-OM-09: Xem danh sách orders**
- Waiter: Xem orders của mình trong ca hiện tại
- Cashier/Manager: Xem tất cả orders
- Lọc theo: Status, Shift, Waiter, Table, Date range

**FR-OM-10: Xem chi tiết order**
- Thông tin order: Items, Total, Discount, Payment
- Lịch sử state transitions
- Thông tin waiter, shift, table

---

### 🪑 TABLE MANAGEMENT (IMPLEMENTED)

**FR-TM-01: Xem danh sách bàn**
- Hiển thị: Tên bàn, Sức chứa, Khu vực, Trạng thái
- Trạng thái: EMPTY, OCCUPIED
- Lọc theo: Status, Area

**FR-TM-02: Tạo bàn mới (Manager)**
- Thông tin: Tên bàn (duy nhất), Sức chứa, Khu vực
- Trạng thái mặc định: EMPTY

**FR-TM-03: Cập nhật thông tin bàn (Manager)**
- Chỉnh sửa: Tên, Sức chứa, Khu vực
- Business rule: Không sửa bàn đang OCCUPIED

**FR-TM-04: Xóa bàn (Manager)**
- Business rule: Chỉ xóa bàn EMPTY, không có order liên quan

**FR-TM-05: Cập nhật trạng thái bàn**
- Auto update khi tạo/thanh toán order
- EMPTY → OCCUPIED (khi tạo order)
- OCCUPIED → EMPTY (khi order PAID)

---

### ⏰ SHIFT MANAGEMENT (IMPLEMENTED)

**FR-SM-01: Mở ca (Waiter)**
- Chọn loại ca: MORNING, AFTERNOON, EVENING
- Nhập tiền đầu ca (start_cash)
- Ghi nhận: waiter_id, started_at
- Business rule: Waiter không thể mở 2 ca cùng lúc

**FR-SM-02: Xem ca hiện tại**
- Hiển thị: Loại ca, Thời gian bắt đầu, Tiền đầu ca
- Số orders trong ca, Tổng doanh thu tạm tính

**FR-SM-03: Kết ca (Waiter)**
- Nhập tiền cuối ca (end_cash)
- Hệ thống tính: Total revenue, Total orders
- Chuyển trạng thái: OPEN → ENDED
- Ghi nhận: ended_at

**FR-SM-04: Chốt ca (Cashier)**
- Đối soát tiền với waiter
- Auto lock tất cả orders trong ca
- Chuyển trạng thái: ENDED → CLOSED
- Ghi nhận: closed_at, closed_by
- Business rule: Chỉ Cashier mới được chốt ca

**FR-SM-05: Xem lịch sử ca**
- Waiter: Xem shifts của mình
- Cashier/Manager: Xem tất cả shifts
- Lọc theo: Waiter, Date range, Status

**FR-SM-06: Xem báo cáo ca**
- Tổng doanh thu theo ca
- Số orders, Trung bình bill
- Phân bổ theo payment method
- So sánh giữa các ca

---

### 🍽️ MENU MANAGEMENT

**FR-MM-01: Xem danh sách menu**
- Hiển thị: Tên món, Giá, Danh mục, Trạng thái
- Lọc theo: Category, Status (Available/Unavailable)

**FR-MM-02: Thêm món mới (Manager)**
- Thông tin: Tên món, Giá, Danh mục, Mô tả, Hình ảnh
- Công thức: Danh sách nguyên liệu + số lượng
- Trạng thái mặc định: Available

**FR-MM-03: Cập nhật món (Manager)**
- Chỉnh sửa: Tên, Giá, Danh mục, Mô tả, Công thức
- Business rule: Không xóa món đã có trong orders

**FR-MM-04: Quản lý danh mục món**
- Tạo/sửa/xóa danh mục
- Business rule: Không xóa danh mục đã có món

**FR-MM-05: Đánh dấu hết hàng**
- Chuyển trạng thái: Available → Unavailable
- Món unavailable không hiển thị khi tạo order

---

### 📊 REPORTING & ANALYTICS

**FR-RA-01: Báo cáo doanh thu**
- Tổng doanh thu theo: Ngày, Tuần, Tháng
- Phân tích theo: Payment method, Shift, Waiter
- Biểu đồ xu hướng

**FR-RA-02: Báo cáo bán hàng**
- Top món bán chạy
- Doanh thu theo danh mục
- Trung bình giá trị bill

**FR-RA-03: Báo cáo hiệu suất**
- Số orders theo waiter
- Doanh thu theo waiter
- Thời gian phục vụ trung bình

**FR-RA-04: Báo cáo tồn kho**
- Nguyên liệu sắp hết
- Lịch sử nhập/xuất
- Giá trị tồn kho

**FR-RA-05: Báo cáo chi phí**
- Tổng chi phí theo loại
- So sánh dự kiến vs thực tế
- Tỷ lệ chi phí/doanh thu

**FR-RA-06: Dashboard tổng quan**
- Doanh thu hôm nay
- Số orders đang xử lý
- Nguyên liệu cần nhập
- Cảnh báo hệ thống

---

## 🔐 User Roles Summary

| Role | Order | Payment | Reconciliation | Management | Reports |
|------|-------|---------|----------------|------------|---------|
| **Waiter** | Tạo/sửa (trước thanh toán) | Thu tiền từng bill | Xem bill của mình | ❌ | Ca của mình |
| **Cashier** | ❌ | ❌ | Đối soát, Chốt ca/ngày | ❌ | Toàn bộ |
| **Manager** | Toàn quyền | Toàn quyền | Toàn quyền | Toàn quyền | Toàn quyền |

## 🔄 Workflow (5 Bước)

1. **MỞ QUẦY** (7:00 AM - Manager/Cashier): Mở cửa quán, kiểm tra thiết bị, chuẩn bị nguyên liệu
2. **MỞ CA** (7:30 AM - Waiter): Waiter login → Hệ thống tạo Shift → Ghi nhận tiền đầu ca
3. **BÁN HÀNG + THU TIỀN** (7:30-12:00 - Waiter): Tạo order → Thu tiền → In bill (Check cash limit: ≥5M → Bàn giao)
4. **KẾT CA** (12:00 PM - Waiter): Waiter nhấn "Kết Ca" → Hệ thống tính tổng → Waiter đếm tiền thực tế → Xác nhận
5. **CHỐT CA** (12:10 PM - Cashier): Đối soát từng waiter → Nhận tiền → Chốt ca → Tạo snapshot → Khóa dữ liệu

## 🔒 4 BA Rules BẮT BUỘC

**Rule 1: Bill gắn với Shift + User**
- Mỗi bill phải ghi nhận: CollectedBy (user_id), ShiftID (shift_id)

**Rule 2: Không sửa bill sau thanh toán**
- Bill đã paid là immutable, chỉ được refund với log đầy đủ

**Rule 3: Cash Limit**
- Nếu tiền mặt ≥ 5M → Waiter phải bàn giao cho Cashier

**Rule 4: QR Ưu Tiên**
- UI ưu tiên hiển thị QR, dashboard tracking % QR vs Cash, cảnh báo nếu cash > 70%

## 📊 State Machine
```
Store: CLOSED → OPEN
Shift: INACTIVE → ACTIVE → ENDED → CLOSED
Bill: DRAFT → PAID → (REFUNDED)
```

## 📡 API Endpoints

**Waiter:**
```
POST /waiter/shift/start, /shift/end
POST /waiter/orders, /orders/:id/payment, /bills/:id/print
GET  /waiter/shift/summary, /reports/my-shift
```

**Cashier:**
```
POST /cashier/fund/open, /shifts/:id/close, /daily-close
GET  /cashier/reconciliation, /reports/revenue
POST /cashier/bills/:id/cancel
```

**Manager:**
```
POST /manager/store/open, /store/close
GET  /manager/reports/*
```

## 🎯 MVP Features

**Phase 1 (Core POS):** ✅ Menu, Table, Order, Payment, Billing, Ingredient, Facility, Expense
**Phase 2 (Enhanced):** 🔄 Bill printing, Reports & analytics, Advanced user management

## 📝 Development Checklist

**Frontend:** View → Service → Store → Router → Navigation
**Backend:** Domain → Repository → Service → Handler → Routes

**Default Users:**
- `admin/admin123` (Manager)
- `waiter1/waiter123` (Waiter)
- `cashier1/cashier123` (Cashier)

---

## 📋 Order Management Implementation Summary

### ✅ Backend Implementation (13 files)

**Phase 1 - Domain Layer:**
- `order.go` - Order entity với state machine (8 states)
- `table.go` - Table entity
- `shift.go` - Shift entity

**Phase 2 - Repository Layer:**
- `order_repository.go` - CRUD + FindByShiftID, FindByWaiterID, FindByStatus
- `table_repository.go` - CRUD + FindByStatus, UpdateStatus
- `shift_repository.go` - CRUD + FindOpenShiftByWaiter, FindByDateRange

**Phase 3 - Service Layer:**
- `order_service.go` - CreateOrder, ConfirmOrder, PayOrder, SendToKitchen, ServeOrder, CancelOrder, RefundOrder, LockOrder
- `table_service.go` - Full CRUD + status management
- `shift_service.go` - StartShift, EndShift, CloseShiftAndLockOrders

**Phase 4 - Handler Layer:**
- `order_handler.go` - 11 HTTP endpoints
- `table_handler.go` - 5 HTTP endpoints
- `shift_handler.go` - 7 HTTP endpoints

**Phase 5 - Routes Integration:**
- `main.go` - 23 new routes với role-based authorization

### ✅ Frontend Implementation (11 files)

**Phase 6 - Services & Stores:**
- `order.js` (service + store) - Full CRUD + state transitions
- `table.js` (service + store) - CRUD operations
- `shift.js` (service + store) - Start, End, Close shifts

**Phase 7 - Views:**
- `OrderView.vue` - Order management UI với status tabs, payment modal
- `TableView.vue` - Table grid với status filter
- `ShiftView.vue` - Shift management với current shift display
- `router/index.js` - 3 new routes
- `Navigation.vue` - Menu items added

### 🔄 Order State Machine
```
CREATED → UNPAID → PAID → IN_PROGRESS → SERVED → LOCKED
           ↓        ↓         ↓
       CANCELLED  REFUNDED  REFUNDED
           ↓        ↓         ↓
        LOCKED   LOCKED    LOCKED
```

### 🎯 Key Business Rules Implemented
- ✅ Order phải gắn với `waiter_id` và `shift_id`
- ✅ Order chỉ tạo được khi có shift OPEN
- ✅ Order phải PAID trước khi gửi kitchen
- ✅ Order LOCKED không thể sửa/xóa
- ✅ Waiter không thể mở 2 shift cùng lúc
- ✅ Auto calculate revenue khi chốt ca
- ✅ Auto lock orders khi chốt ca
- ✅ Payment methods: CASH, TRANSFER, QR

### 📊 Implementation Stats
- **Total Files:** 24 files created/updated
- **API Endpoints:** 23 new routes
- **State Transitions:** 8 states với validation
- **Roles Supported:** Waiter, Cashier, Manager
- **UI Components:** 3 major views với responsive design
