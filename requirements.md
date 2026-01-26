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

**FR-C-01: Mở quỹ đầu ca**
- Nhập tiền đầu ca
- Gắn với: Ca, Quầy, Cashier

**FR-C-02: Theo dõi quỹ trong ca**
- Tổng tiền mặt thực tế
- Tổng tiền hệ thống ghi nhận

**FR-C-03: Xem danh sách bill theo waiter**
- Theo ca
- Theo phương thức thanh toán

**FR-C-04: Đối soát tiền waiter nộp**
- Mỗi waiter:
  - Tổng tiền phải nộp
  - Tiền thực nhận
  - Chênh lệch (+ / -)

**FR-C-05: Xác nhận chốt ca cho waiter**
- Sau khi đối soát xong
- Khóa dữ liệu ca của waiter

**FR-C-06: Chốt ca**
- Tổng hợp: Doanh thu, Tiền mặt, Không tiền mặt, Chênh lệch

**FR-C-07: Chốt ngày (nếu cashier có quyền)**
- Tổng hợp nhiều ca
- Snapshot dữ liệu

**FR-C-08: Xử lý ngoại lệ**
- Hủy bill (có lý do)
- Điều chỉnh sai sót (ghi log)

**FR-C-09: Xem báo cáo**
- Doanh thu
- Hiệu suất waiter
- Chênh lệch tiền

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
