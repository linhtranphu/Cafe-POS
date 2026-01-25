# 🏗️ Café POS System - Project Structure

## 📋 Core Requirements Analysis

### 5.1 Order Management
- ✅ Tạo order mới
- ✅ Gắn order với bàn  
- ✅ CRUD món trong order
- ✅ Tính tổng tiền realtime
- ✅ State machine: CREATED → CONFIRMED → PAID → BILLED

### 5.2 Menu & Pricing
- ✅ Danh sách món uống
- ✅ Giá bán
- 🔄 Size/option (future)

### 5.3 Table Management
- ✅ Danh sách bàn
- ✅ Trạng thái: Empty → Occupied → Paid
- ✅ 1 order = 1 bàn tại 1 thời điểm

### 5.4 Payment
- ✅ Cash/Transfer
- ✅ Ghi nhận thời gian + phương thức
- ✅ PAID chỉ 1 lần, không rollback (trừ admin)

### 5.5 Billing/Printing
- ✅ In bill sau thanh toán
- ✅ In lại bill (audit)
- ✅ Thông tin đầy đủ

### 5.6 User & Role
- ✅ Nhân viên vs Quản lý
- ✅ Phân quyền rõ ràng
- ✅ JWT Authentication
- ✅ Role-based Authorization

### 5.7 Ingredient Management
- ✅ Quản lý nguyên liệu với đơn vị chuẩn ISO
- ✅ Theo dõi tồn kho realtime
- ✅ Cảnh báo sắp hết hàng
- ✅ Lịch sử biến động chi tiết
- ✅ Tìm kiếm và lọc nguyên liệu
- ✅ Phân quyền Manager/Staff

## 🏗️ Architecture Design

```
cafe-pos/
├── backend/                 # Go + Gin + MongoDB
│   ├── domain/
│   │   ├── order/          # Order entity & business logic
│   │   ├── menu/           # Menu items & pricing
│   │   ├── table/          # Table management
│   │   ├── payment/        # Payment processing
│   │   ├── ingredient/     # Ingredient inventory & stock
│   │   └── user/           # User & roles
│   ├── application/
│   │   └── services/       # Business services
│   ├── infrastructure/
│   │   ├── mongodb/        # Database layer
│   │   └── printer/        # Bill printing
│   └── interfaces/
│       └── http/           # REST API
├── frontend/               # Vue.js 3 POS Interface
│   ├── views/
│   │   ├── OrderView.vue   # Main POS screen
│   │   ├── TableView.vue   # Table management
│   │   ├── MenuView.vue    # Menu management (Manager)
│   │   ├── IngredientView.vue # Ingredient management (Manager)
│   │   └── ReportView.vue  # Reports (manager)
│   └── components/
└── docker-compose.yml      # Development setup
```

## 📊 Database Schema

### Collections:
1. **orders** - Order management
2. **menu_items** - Menu & pricing
3. **tables** - Table status
4. **payments** - Payment records
5. **users** - Staff & managers
6. **bills** - Bill history
7. **ingredients** - Ingredient inventory
8. **stock_history** - Stock movement tracking

## 🔐 User Roles & Permissions

### 1. Waiter/Staff (Nhân viên order)
**Quyền hạn:**
- Tạo order mới
- Nhập món, số lượng
- Gắn bàn
- Xem & thông báo tổng tiền
- Chọn phương thức thanh toán
- Xác nhận đã thu tiền
- In bill

**Hạn chế:**
- ❌ Không sửa order sau khi đã thanh toán
- ❌ Không xem báo cáo doanh thu

### 2. Cashier (Thu ngân)
**Quyền hạn:**
- Xem order đã tạo
- Thu tiền
- Xác nhận thanh toán
- In / in lại bill

### 3. Manager/Store Owner (Quản lý)
**Quyền hạn:**
- Xem tất cả order
- Xem báo cáo doanh thu
- In lại bill
- Chỉnh sửa / hủy order đã thanh toán (có log)
- Quản lý menu & giá
- Quản lý bàn
- Quản lý user

### Default Users:
- `admin/admin123` (Manager)
- `waiter1/waiter123` (Waiter)
- `cashier1/cashier123` (Cashier)

## 🥬 Ingredient Management System

### ✅ Functional Requirements Implemented:

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

### 📊 Sample Data:
- **18 ingredients** với đầy đủ thông tin
- **ISO standard units**: kg, g, L, ml, piece, box, pack
- **Categories**: Cà phê, Trà, Sữa, Đường, Trái cây, Bánh, Khác
- **Cost tracking**: Giá vốn mỗi đơn vị
- **Stock history**: Tự động log mọi thay đổi

## 🏢 Facility Management System

### ✅ Functional Requirements:

**FR-FM-01: Xem danh sách cơ sở vật chất**

Mô tả: Hệ thống cho phép người dùng xem danh sách tất cả tài sản/cơ sở vật chất.

Thông tin hiển thị:
- Tên tài sản
- Loại tài sản
- Khu vực sử dụng
- Số lượng
- Tình trạng
- Ngày mua / đưa vào sử dụng

**FR-FM-02: Thêm mới cơ sở vật chất**

Mô tả: Hệ thống cho phép quản lý thêm mới một tài sản hoặc nhóm tài sản.

Thông tin bắt buộc:
- Tên
- Loại (bàn, máy móc, dụng cụ…)
- Số lượng
- Khu vực
- Ngày mua
- Tình trạng ban đầu

Business Rules:
- Tài sản đơn chiếc (máy) quản lý theo từng item
- Dụng cụ tiêu hao (ly, thìa) quản lý theo số lượng

**FR-FM-03: Cập nhật thông tin cơ sở vật chất**

Mô tả: Cho phép chỉnh sửa thông tin tài sản.

Cho phép chỉnh sửa:
- Khu vực
- Tình trạng
- Số lượng (với tài sản nhóm)
- Ghi chú

Không cho phép:
- Xóa tài sản đã có lịch sử bảo trì (chỉ được inactive)

**FR-FM-04: Quản lý tình trạng tài sản**

Mô tả: Hệ thống cho phép cập nhật trạng thái tài sản.

Trạng thái:
- Đang sử dụng
- Hỏng
- Đang sửa
- Ngừng sử dụng
- Thanh lý

**FR-FM-05: Báo hư hỏng (Staff)**

Mô tả: Nhân viên có thể tạo yêu cầu báo hư hỏng.

Thông tin:
- Tài sản
- Mô tả sự cố
- Mức độ ảnh hưởng
- Hình ảnh (optional)

**FR-FM-06: Quản lý bảo trì / sửa chữa**

Mô tả: Quản lý tạo và theo dõi các hoạt động bảo trì.

Thông tin bảo trì:
- Tài sản liên quan
- Loại bảo trì (định kỳ / phát sinh)
- Nội dung
- Ngày thực hiện
- Chi phí
- Đơn vị sửa chữa

**FR-FM-07: Lịch sử tài sản**

Mô tả: Hệ thống lưu toàn bộ lịch sử của tài sản.

Bao gồm:
- Thay đổi trạng thái
- Bảo trì
- Di chuyển khu vực
- Thanh lý

**FR-FM-08: Tìm kiếm & lọc**

Mô tả:
- Theo tên
- Theo loại
- Theo khu vực
- Theo tình trạng

**FR-FM-09: Phân quyền truy cập**

Mô tả:
- Manager: toàn quyền
- Staff:
  - Xem danh sách
  - Báo hư hỏng
  - Không chỉnh sửa tài sản

## 💰 Expense Management System

### 4. Phân tích nghiệp vụ cốt lõi

**4.1 Nghiệp vụ ghi nhận chi phí**

Use case: Ghi nhận chi phí

Trigger: Có chi phí phát sinh

Thông tin cần nhập:
- Loại chi phí (điện, nước, thuê nhà…)
- Số tiền
- Kỳ áp dụng (tháng/năm)
- Ngày phát sinh
- Hình thức thanh toán
- Ghi chú / hóa đơn (ảnh)

**4.2 Chi phí định kỳ (Recurring Cost)**

Ví dụ:
- Tiền thuê nhà hàng tháng
- Tiền điện, nước

Đặc điểm nghiệp vụ:
- Lặp theo chu kỳ (tháng)
- Có thể:
  - Ghi nhận tự động
  - Hoặc nhắc nhập

**4.3 Chi phí trả trước / phân bổ**

Ví dụ: Tiền cọc nhà
- Trả 1 lần
- Áp dụng cho nhiều tháng

Business rule:
- Chi phí gốc: 1 lần
- Có thể:
  - Không phân bổ (chỉ theo dõi)
  - Phân bổ theo thời gian (12/24 tháng)

### ✅ Functional Requirements:

**FR-EX-01: Xem danh sách chi phí**

Mô tả: Danh sách chi phí theo thời gian

Hiển thị:
- Tên chi phí
- Nhóm chi phí
- Số tiền
- Tháng áp dụng

**FR-EX-02: Khai báo loại chi phí**

Mô tả: Tạo danh mục:
- Thuê nhà
- Điện
- Nước
- Internet

Gắn:
- Cố định / biến đổi / một lần
- Có định kỳ hay không

**FR-EX-03: Ghi nhận chi phí phát sinh**

Mô tả:
- Nhập chi phí thủ công
- Đính kèm hóa đơn

**FR-EX-04: Chi phí định kỳ**

Mô tả: Thiết lập:
- Chu kỳ (tháng)
- Số tiền dự kiến
- Hệ thống nhắc nhập chi phí hàng tháng

**FR-EX-05: Quản lý chi phí trả trước / phân bổ**

Ví dụ: Tiền cọc nhà

Nhập:
- Tổng số tiền
- Thời gian áp dụng
- Hệ thống tự phân bổ chi phí theo tháng

**FR-EX-06: Báo cáo chi phí**

Mô tả: Tổng chi phí theo:
- Tháng
- Loại

So sánh:
- Thực tế vs dự kiến

**FR-EX-07: Phân quyền**

Mô tả:
- Manager: toàn quyền
- Accountant: xem & xuất báo cáo

### 6. Business Rules quan trọng

- Một chi phí phải thuộc 1 loại chi phí
- Chi phí phân bổ:
  - Không được chỉnh sửa kỳ đã chốt
  - Không cho xóa chi phí đã dùng cho báo cáo
- Chi phí định kỳ:
  - Không auto ghi nhận nếu chưa xác nhận

## 🎯 MVP Features Priority

### Phase 1 (Core POS):
1. ✅ Menu management
2. ✅ Table management  
3. ✅ Order creation & management
4. ✅ Payment processing
5. ✅ Basic billing
6. ✅ Ingredient management

### Phase 2 (Enhanced):
1. 🔄 Bill printing
2. 🔄 Reports & analytics
3. 🔄 Advanced user management

Bạn có muốn tôi bắt đầu implement Café POS System này không?