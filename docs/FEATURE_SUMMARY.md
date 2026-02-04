# Tóm Tắt Tính Năng - Hệ Thống Café POS

## 🎯 Tổng Quan
Hệ thống quản lý quán cà phê toàn diện với 4 vai trò người dùng và quy trình làm việc được tự động hóa.

## 👥 Vai Trò Người Dùng

| Vai Trò | Quyền Hạn Chính |
|---------|-----------------|
| **Manager** | Quản lý toàn hệ thống, người dùng, menu, kho, cơ sở vật chất, chi phí |
| **Cashier** | Quản lý ca thu ngân, đối soát tiền, xử lý thanh toán, khóa đơn hàng |
| **Waiter** | Tạo đơn hàng, thu tiền, quản lý ca phục vụ, phục vụ khách hàng |
| **Barista** | Nhận đơn hàng, pha chế, cập nhật trạng thái, quản lý hàng đợi |

## 📋 Tính Năng Chính

### 1. Quản Lý Đơn Hàng
- **Vòng đời**: CREATED → PAID → QUEUED → IN_PROGRESS → READY → SERVED → LOCKED
- **Tính năng**: Tạo, chỉnh sửa, thanh toán, hoàn tiền, hủy, theo dõi thời gian thực
- **Thanh toán**: Tiền mặt, QR, Chuyển khoản (hỗ trợ thanh toán từng phần)

### 2. Quản Lý Ca Làm Việc
#### Ca Phục Vụ/Pha Chế:
- Mở/đóng ca theo vai trò
- Theo dõi thời gian và doanh thu
- Ngăn nhiều ca mở cùng lúc

#### Ca Thu Ngân (Phức Tạp):
- **Quy trình**: OPEN → CLOSURE_INITIATED → CLOSED
- **Bước đóng ca**: Khởi tạo → Đếm tiền → Tính chênh lệch → Ghi chép → Xác nhận → Đóng

### 3. Quản Lý Thực Đơn
- Tạo/sửa/xóa món ăn
- Quản lý giá và danh mục
- Theo dõi thành phần nguyên liệu
- Trạng thái có sẵn

### 4. Quản Lý Kho Nguyên Liệu
- Theo dõi số lượng, đơn vị (kg, L, cái)
- Cảnh báo hết hàng
- Lịch sử điều chỉnh kho
- Tự động tạo chi phí khi mua

### 5. Quản Lý Cơ Sở Vật Chất
- **Loại**: Nội thất, Máy móc, Dụng cụ, Điện tử
- **Khu vực**: Phòng khách, Bếp, Quầy bar, Kho, Văn phòng
- **Trạng thái**: Đang dùng, Hỏng, Đang sửa, Không hoạt động, Đã thanh lý
- **Bảo trì**: Định kỳ, khẩn cấp, theo dõi chi phí

### 6. Quản Lý Chi Phí
- **Loại**: Thủ công, Tự động (từ nguyên liệu, cơ sở vật chất, bảo trì)
- **Tính năng**: Chi phí định kỳ, trả trước, phân loại
- Tự động tạo chi phí khi mua hàng/bảo trì

### 7. Báo Cáo & Phân Tích
- **Thu ngân**: Báo cáo ca, đối soát, sai lệch, doanh thu
- **Đơn hàng**: Theo trạng thái, thời gian, nhân viên
- **Kho**: Cảnh báo hết hàng, lịch sử sử dụng
- **Cơ sở vật chất**: Lịch sử bảo trì, chi phí
- **Chi phí**: Theo danh mục, thời gian, định kỳ

## 🔧 Tính Năng Kỹ Thuật

### State Machines
- **Đơn hàng**: Xác thực chuyển đổi trạng thái, tính tiến độ
- **Ca làm việc**: Quản lý vòng đời ca, ngăn ca đồng thời
- **Ca thu ngân**: Quy trình đóng ca nhiều bước

### Bảo Mật
- JWT authentication
- Phân quyền theo vai trò (RBAC)
- Mã hóa mật khẩu (bcrypt)
- HTTPS support

### Kiểm Toán
- Lịch sử đơn hàng với timestamp
- Nhật ký ca thu ngân
- Theo dõi thay đổi cơ sở vật chất
- Lịch sử điều chỉnh kho

## 🖥️ Giao Diện

### Views Chính:
- **Dashboard**: Theo vai trò với thống kê nhanh
- **OrderView**: Quản lý đơn hàng với filter trạng thái
- **BaristaView**: 3 tab (Queue, Working, Ready)
- **CashierDashboard**: Thông tin ca, thanh toán
- **CashierShiftClosure**: Quy trình đóng ca từng bước
- **Management Views**: Menu, Kho, Cơ sở vật chất, Chi phí, Người dùng

### Tính Năng UI:
- Responsive mobile-first design
- Real-time updates
- Status filtering
- Quick actions
- Progress indicators

## 🚀 Triển Khai

### Tech Stack:
- **Backend**: Go + Gin + MongoDB
- **Frontend**: Vue.js 3 + Vite + Tailwind CSS
- **Infrastructure**: Docker + Nginx

### API Endpoints:
- `/api/auth/*` - Xác thực
- `/api/orders/*` - Đơn hàng
- `/api/shifts/*` - Ca làm việc
- `/api/cashier-shifts/*` - Ca thu ngân
- `/api/cashier/*` - Thao tác thu ngân
- `/api/menu/*` - Thực đơn
- `/api/ingredients/*` - Nguyên liệu
- `/api/facilities/*` - Cơ sở vật chất
- `/api/expenses/*` - Chi phí
- `/api/users/*` - Người dùng

## 📊 Quy Trình Kinh Doanh

### Xử Lý Đơn Hàng:
1. Waiter tạo đơn → 2. Thu tiền → 3. Gửi bar → 4. Barista nhận → 5. Pha chế → 6. Sẵn sàng → 7. Phục vụ → 8. Cashier khóa

### Đóng Ca Thu Ngân:
1. Khởi tạo → 2. Đếm tiền → 3. Tính chênh lệch → 4. Ghi chép (nếu có) → 5. Xác nhận → 6. Đóng ca

### Quản Lý Kho:
1. Tạo nguyên liệu → 2. Sử dụng → 3. Điều chỉnh → 4. Theo dõi → 5. Cảnh báo → 6. Mua bổ sung

## 🎯 Điểm Nổi Bật

✅ **State Machine** - Đảm bảo quy trình chính xác  
✅ **Multi-role** - 4 vai trò với quyền hạn riêng biệt  
✅ **Real-time** - Cập nhật trạng thái tức thời  
✅ **Audit Trail** - Theo dõi tất cả thay đổi  
✅ **Auto Expense** - Tự động tạo chi phí  
✅ **Mobile First** - Thiết kế tối ưu cho mobile  
✅ **Comprehensive** - Quản lý toàn diện quán cà phê  

---

*Hệ thống này cung cấp giải pháp hoàn chỉnh cho việc quản lý quán cà phê từ đơn hàng đến báo cáo, với quy trình được tự động hóa và kiểm soát chặt chẽ.*