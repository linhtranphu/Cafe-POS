# Tài Liệu Tính Năng Toàn Diện - Hệ Thống Café POS

## Tổng Quan Hệ Thống

Hệ thống Café POS là một ứng dụng quản lý quán cà phê toàn diện được xây dựng với kiến trúc hiện đại:

### Công Nghệ Sử Dụng
- **Backend**: Go 1.21+ với Gin Web Framework, MongoDB
- **Frontend**: Vue.js 3 với Vite, Tailwind CSS, Pinia state management
- **Cơ sở hạ tầng**: Docker & Docker Compose, Nginx, MongoDB 7.0
- **Xác thực**: JWT với phân quyền theo vai trò (RBAC)

---

## Phân Quyền Người Dùng

Hệ thống hỗ trợ 4 vai trò người dùng chính:

### 1. **Manager (Quản lý)** - Quyền truy cập toàn hệ thống
- Quản lý người dùng (tạo, sửa, xóa)
- Quản lý menu và thực đơn
- Theo dõi kho nguyên liệu
- Quản lý cơ sở vật chất & thiết bị
- Theo dõi chi phí và báo cáo
- Xem tất cả đơn hàng và ca làm việc
- Phân tích dữ liệu toàn hệ thống

### 2. **Cashier (Thu ngân)** - Xử lý thanh toán & quản lý ca
- Quản lý ca thu ngân (mở/đóng ca)
- Đối soát tiền mặt
- Xử lý thanh toán và thu tiền
- Xử lý sai lệch thanh toán
- Khóa và hoàn thiện đơn hàng
- Báo cáo hàng ngày và bàn giao ca
- Kiểm toán thanh toán

### 3. **Waiter (Phục vụ)** - Quản lý đơn hàng & phục vụ bàn
- Tạo và quản lý đơn hàng
- Quản lý bàn
- Theo dõi trạng thái đơn hàng
- Quản lý ca làm việc (mở/đóng ca)
- Thu tiền khách hàng
- Chỉnh sửa đơn hàng trước khi thanh toán

### 4. **Barista (Pha chế)** - Chuẩn bị đồ uống & hoàn thành đơn
- Quản lý hàng đợi đơn hàng
- Quy trình pha chế đồ uống
- Cập nhật trạng thái đơn hàng (nhận, hoàn thành)
- Theo dõi ca làm việc
- Hoàn thành đơn hàng

---

## Tính Năng Cốt Lõi

### 1. QUẢN LÝ ĐƠN HÀNG

#### Vòng Đời Đơn Hàng (State Machine):
```
CREATED → PAID → QUEUED → IN_PROGRESS → READY → SERVED → LOCKED
```
**Đường dẫn thay thế**: CANCELLED, REFUNDED

#### Các Thao Tác Đơn Hàng:
- **Tạo đơn hàng** với nhiều món
- **Chỉnh sửa đơn hàng** (trước khi thanh toán)
- **Thu tiền** (một phần hoặc toàn bộ)
- **Hoàn tiền** (trước khi vào hàng đợi)
- **Gửi đến quầy bar/pha chế**
- **Hủy đơn hàng** với lý do theo dõi
- **Khóa đơn hàng** sau khi hoàn thành phục vụ

#### Chi Tiết Đơn Hàng Được Theo Dõi:
- **Số đơn hàng** (tự động: YYYYMMDD-HHMMSS-XXX)
- **Tên khách hàng** (tùy chọn)
- **Món và số lượng, giá**
- **Tổng phụ, giảm giá, tổng cộng**
- **Số tiền đã trả và còn nợ**
- **Phương thức thanh toán** (Tiền mặt, QR, Chuyển khoản)
- **Phân công phục vụ và pha chế**
- **Thời gian cho mỗi chuyển đổi trạng thái**
- **Ghi chú và lý do hủy/hoàn tiền**

#### Tính Năng Đơn Hàng:
- Theo dõi trạng thái đơn hàng thời gian thực
- Tính toán phần trăm tiến độ đơn hàng
- Tự động tính hoàn tiền khi chỉnh sửa
- Phát hiện sai lệch thanh toán
- Lịch sử và kiểm toán đơn hàng

---

### 2. QUẢN LÝ CA LÀM VIỆC

#### Hai Loại Ca Làm Việc:

#### A. Ca Phục Vụ/Pha Chế:
- **Mở/đóng ca** theo vai trò (phục vụ hoặc pha chế)
- **Theo dõi thời gian ca**
- **Tính doanh thu ca**
- **Đếm đơn hàng mỗi ca**
- **Ngăn nhiều ca mở cùng lúc**
- **Loại ca**: SÁNG, CHIỀU, TỐI
- **Khóa đơn hàng** khi đóng ca

#### B. Ca Thu Ngân (Quy Trình Phức Tạp):
**Luồng Trạng Thái**: `OPEN → CLOSURE_INITIATED → CLOSED`

**Các Bước Quy Trình Đóng Ca**:
1. **Khởi tạo đóng ca**
2. **Ghi nhận số tiền thực tế**
3. **Tính toán chênh lệch** (tiền hệ thống vs tiền thực tế)
4. **Ghi chép chênh lệch** với lý do và ghi chú (nếu khác 0)
5. **Xác nhận trách nhiệm**
6. **Đóng ca**

#### Tính Năng Ca Thu Ngân:
- **Theo dõi tiền đầu ca**
- **Tính toán tiền hệ thống**
- **Ghi nhận tiền thực tế**
- **Tính toán và ghi chép chênh lệch**
- **Xác nhận trách nhiệm**
- **Nhật ký kiểm toán** cho tất cả hành động
- **Bất biến** sau khi đóng
- **Ngăn đóng ca** nếu ca phục vụ còn mở

---

### 3. XỬ LÝ THANH TOÁN

#### Phương Thức Thanh Toán:
- **Tiền mặt** (💵)
- **Mã QR** (📱)
- **Chuyển khoản ngân hàng** (🏦)

#### Tính Năng Thanh Toán:
- **Hỗ trợ thanh toán từng phần**
- **Nhiều lần thu tiền** cho một đơn hàng
- **Theo dõi phương thức thanh toán**
- **Xác định người thu tiền**
- **Kiểm toán thanh toán**
- **Phát hiện và báo cáo sai lệch**
- **Khả năng ghi đè thanh toán** (quản lý/thu ngân)
- **Đối soát tiền mặt**

#### Giám Sát Thanh Toán:
- **Theo dõi tất cả thanh toán** theo ca
- **Phát hiện sai lệch thanh toán**
- **Báo cáo và giải quyết sai lệch**
- **Lịch sử kiểm toán thanh toán**
- **Báo cáo thanh toán hàng ngày**

---

### 4. QUẢN LÝ THỰC ĐƠN

#### Quản Lý Món Ăn:
- **Tạo/sửa/xóa** món trong thực đơn
- **Quản lý giá**
- **Tổ chức danh mục**
- **Mô tả món**
- **Theo dõi thành phần nguyên liệu**
- **Trạng thái có sẵn** (có/không có)
- **Tìm kiếm và lọc** món

#### Danh Mục:
- **Cà phê** (☕)
- **Trà** (🍵)
- **Nước ép** (🧃)
- **Đồ ăn** (🍰)
- **Danh mục tùy chỉnh**

---

### 5. QUẢN LÝ NGUYÊN LIỆU

#### Theo Dõi Kho:
- **Tên nguyên liệu, danh mục, loại đơn vị**
- **Số lượng hiện tại**
- **Mức tồn kho tối thiểu**
- **Giá mỗi đơn vị**
- **Thông tin nhà cung cấp**
- **Lịch sử kho** với thời gian

#### Loại Đơn Vị Hỗ Trợ:
- **Khối lượng**: kg, g
- **Thể tích**: L, ml
- **Số lượng**: cái, hộp, gói

#### Thao Tác Kho:
- **Điều chỉnh kho** với theo dõi lý do
- **Cảnh báo hết hàng**
- **Kiểm toán lịch sử kho**
- **Theo dõi người dùng** cho điều chỉnh
- **Tự động theo dõi chi phí** cho mua hàng

#### Danh Mục Nguyên Liệu:
- **Tạo/quản lý** danh mục nguyên liệu
- **Tổ chức theo loại**

---

### 6. QUẢN LÝ CƠ SỞ VẬT CHẤT & THIẾT BỊ

#### Theo Dõi Cơ Sở Vật Chất:
- **Tên, loại, khu vực/vị trí**
- **Số lượng**
- **Trạng thái** (Đang dùng, Hỏng, Đang sửa, Không hoạt động, Đã thanh lý)
- **Ngày mua và chi phí**
- **Thông tin nhà cung cấp**
- **Ghi chú**

#### Loại Cơ Sở Vật Chất:
- **Nội thất** (Bàn ghế)
- **Máy móc** (Máy móc)
- **Dụng cụ** (Dụng cụ)
- **Điện tử** (Điện tử)
- **Khác** (Khác)

#### Khu Vực Cơ Sở Vật Chất:
- **Phòng khách** (Phòng khách)
- **Bếp** (Bếp)
- **Quầy/Bar** (Quầy bar)
- **Kho** (Kho)
- **Văn phòng** (Văn phòng)
- **Khác** (Khác)

#### Quản Lý Bảo Trì:
- **Tạo hồ sơ bảo trì**
- **Loại bảo trì**: định kỳ, khẩn cấp
- **Theo dõi chi phí bảo trì**
- **Thông tin nhà cung cấp**
- **Lịch sử bảo trì** mỗi cơ sở vật chất
- **Lập kế hoạch bảo trì định kỳ**
- **Báo cáo sự cố** với mức độ nghiêm trọng
- **Thống kê bảo trì**

#### Lịch Sử Cơ Sở Vật Chất:
- **Theo dõi tất cả thay đổi** (tạo, cập nhật, di chuyển, thay đổi trạng thái)
- **Theo dõi người dùng và thời gian**
- **So sánh giá trị cũ và mới**
- **Tự động theo dõi chi phí** cho bảo trì

---

### 7. QUẢN LÝ CHI PHÍ

#### Loại Chi Phí:
- **Chi phí thủ công**
- **Chi phí tự động theo dõi** (từ nguyên liệu, cơ sở vật chất, bảo trì)

#### Danh Mục Chi Phí:
- **Tạo/quản lý** danh mục chi phí
- **Phân loại tất cả chi phí**

#### Tính Năng Chi Phí:
- **Theo dõi ngày**
- **Số tiền và phương thức thanh toán**
- **Thông tin nhà cung cấp**
- **Mô tả và ghi chú**
- **Theo dõi nguồn** (nguyên liệu, cơ sở vật chất, bảo trì, thủ công)
- **Theo dõi người dùng** (ai tạo chi phí)

#### Chi Phí Định Kỳ:
- **Thiết lập chi phí định kỳ**
- **Tùy chọn tần suất**: hàng ngày, hàng tuần, hàng tháng, hàng quý, hàng năm
- **Theo dõi ngày đến hạn tiếp theo**
- **Trạng thái hoạt động/không hoạt động**

#### Chi Phí Trả Trước:
- **Theo dõi số tiền trả trước**
- **Ngày bắt đầu và kết thúc**
- **Theo dõi khấu hao**

#### Theo Dõi Chi Phí Tự Động:
- **Tự động tạo chi phí** cho mua nguyên liệu
- **Tự động tạo chi phí** cho mua cơ sở vật chất
- **Tự động tạo chi phí** cho chi phí bảo trì
- **Liên kết với hồ sơ nguồn** để kiểm toán

---

### 8. BÁO CÁO & PHÂN TÍCH

#### Báo Cáo Thu Ngân:
- **Báo cáo ca hàng ngày**
- **Báo cáo đóng ca**
- **Báo cáo đối soát thanh toán**
- **Báo cáo sai lệch**
- **Doanh thu theo phương thức thanh toán**
- **Phân tích dòng tiền**

#### Báo Cáo Đơn Hàng:
- **Đơn hàng theo trạng thái**
- **Đơn hàng theo khoảng thời gian**
- **Đơn hàng theo phục vụ**
- **Đơn hàng theo pha chế**
- **Phân tích doanh thu**
- **Tỷ lệ hoàn thành đơn hàng**

#### Báo Cáo Kho:
- **Cảnh báo hết hàng**
- **Theo dõi sử dụng nguyên liệu**
- **Lịch sử kho**

#### Báo Cáo Cơ Sở Vật Chất:
- **Lịch sử bảo trì**
- **Trạng thái thiết bị**
- **Chi phí bảo trì**
- **Theo dõi sự cố**

#### Báo Cáo Chi Phí:
- **Chi phí theo danh mục**
- **Chi phí theo khoảng thời gian**
- **Theo dõi chi phí định kỳ**
- **Khấu hao chi phí trả trước**

---

### 9. STATE MACHINES & QUY TRÌNH

#### State Machine Đơn Hàng:
- **Xác thực tất cả chuyển đổi** đơn hàng
- **Thực thi quy tắc kinh doanh**
- **Ngăn thay đổi trạng thái không hợp lệ**
- **Tính toán tiến độ đơn hàng**
- **Xác định hành động tiếp theo dự kiến**

#### State Machine Ca Phục Vụ/Pha Chế:
- **Quản lý vòng đời ca**
- **Xác thực bắt đầu/kết thúc ca**
- **Tính toán thời gian ca**
- **Ngăn ca đồng thời**

#### State Machine Ca Thu Ngân:
- **Quy trình đóng ca nhiều bước phức tạp**
- **Xác thực từng bước**
- **Thực thi hoàn thành điều kiện tiên quyết**
- **Ngăn bỏ qua bước**
- **Cho phép hủy đóng ca** (với hạn chế)

#### Quản Lý State Machine Tập Trung:
- **Cung cấp truy cập thống nhất** cho tất cả state machines
- **Xác thực chuyển đổi** qua các domain
- **Thực thi quy tắc kinh doanh**
- **Tính toán tiến độ và hành động tiếp theo**

---

### 10. QUẢN LÝ NGƯỜI DÙNG

#### Thao Tác Người Dùng:
- **Tạo người dùng** với phân quyền vai trò
- **Chỉnh sửa thông tin người dùng**
- **Đặt lại mật khẩu người dùng** (bởi quản lý)
- **Thay đổi mật khẩu riêng**
- **Bật/tắt trạng thái người dùng**
- **Xóa người dùng** (với hạn chế)
- **Xem tất cả người dùng** hoặc lọc theo vai trò
- **Theo dõi lần đăng nhập cuối**

#### Xác Thực Người Dùng:
- **Tính duy nhất tên người dùng**
- **Mã hóa mật khẩu** (bcrypt)
- **Xác thực vai trò**
- **Ngăn xóa quản lý cuối cùng**

---

### 11. XÁC THỰC & BẢO MẬT

#### Xác Thực:
- **Xác thực token JWT**
- **Đăng nhập với tên người dùng và mật khẩu**
- **Hết hạn token**
- **Theo dõi lần đăng nhập cuối**

#### Tính Năng Bảo Mật:
- **Mã hóa mật khẩu** với bcrypt
- **Kiểm soát truy cập dựa trên vai trò** (RBAC)
- **Middleware cho xác thực/ủy quyền**
- **Cấu hình dựa trên môi trường**
- **Xác thực MongoDB**
- **Hỗ trợ HTTPS** trong sản xuất

---

## Giao Diện Frontend

### Các View Chính:

1. **LoginView** - Xác thực người dùng
2. **DashboardView** - Dashboard theo vai trò
   - **Manager**: Tổng quan hệ thống, thống kê nhanh
   - **Cashier**: Thông tin ca, theo dõi doanh thu, ca mở
   - **Waiter**: Thống kê đơn hàng, đơn hàng chờ
   - **Barista**: Thống kê hàng đợi, đơn hàng đang làm

3. **OrderView** - Quản lý đơn hàng
4. **BaristaView** - Quy trình pha chế
5. **ShiftView** - Quản lý ca làm việc
6. **CashierDashboard** - Thao tác thu ngân
7. **CashierShiftClosure** - Quy trình đóng ca
8. **MenuView** - Quản lý thực đơn
9. **IngredientManagementView** - Quản lý kho
10. **FacilityManagementView** - Quản lý thiết bị
11. **ExpenseManagementView** - Theo dõi chi phí
12. **UserManagementView** - Quản trị người dùng
13. **ProfileView** - Hồ sơ người dùng
14. **CashierReports** - Báo cáo

---

## API Endpoints

### Xác Thực:
- `POST /api/auth/login` - Đăng nhập
- `POST /api/auth/logout` - Đăng xuất

### Đơn Hàng:
- `POST /api/orders` - Tạo đơn hàng
- `GET /api/orders` - Lấy tất cả đơn hàng
- `GET /api/orders/:id` - Lấy chi tiết đơn hàng
- `POST /api/orders/:id/payment` - Thu tiền
- `POST /api/orders/:id/edit` - Chỉnh sửa đơn hàng
- `POST /api/orders/:id/refund` - Hoàn tiền
- `POST /api/orders/:id/send-to-bar` - Gửi đến quầy bar
- `POST /api/orders/:id/accept` - Nhận đơn hàng (pha chế)
- `POST /api/orders/:id/finish` - Đánh dấu sẵn sàng (pha chế)
- `POST /api/orders/:id/serve` - Phục vụ đơn hàng
- `POST /api/orders/:id/cancel` - Hủy đơn hàng
- `POST /api/orders/:id/lock` - Khóa đơn hàng

### Ca Làm Việc:
- `POST /api/shifts` - Bắt đầu ca
- `GET /api/shifts/current` - Lấy ca hiện tại
- `POST /api/shifts/:id/end` - Kết thúc ca
- `GET /api/shifts` - Lấy tất cả ca
- `GET /api/shifts/open` - Lấy ca đang mở

### Ca Thu Ngân:
- `POST /api/cashier-shifts` - Bắt đầu ca thu ngân
- `GET /api/cashier-shifts/current` - Lấy ca thu ngân hiện tại
- `GET /api/cashier-shifts/:id` - Lấy chi tiết ca thu ngân

### Thao Tác Thu Ngân:
- `POST /api/cashier/reconcile` - Đối soát tiền mặt
- `POST /api/cashier/discrepancy` - Báo cáo sai lệch
- `POST /api/cashier/override-payment/:id` - Ghi đè thanh toán
- `POST /api/cashier/lock-order/:id` - Khóa đơn hàng

### Thực Đơn:
- `POST /api/menu` - Tạo món
- `GET /api/menu` - Lấy tất cả món
- `PUT /api/menu/:id` - Cập nhật món
- `DELETE /api/menu/:id` - Xóa món

### Nguyên Liệu:
- `POST /api/ingredients` - Tạo nguyên liệu
- `GET /api/ingredients` - Lấy tất cả nguyên liệu
- `PUT /api/ingredients/:id` - Cập nhật nguyên liệu
- `POST /api/ingredients/:id/adjust` - Điều chỉnh kho

### Cơ Sở Vật Chất:
- `POST /api/facilities` - Tạo cơ sở vật chất
- `GET /api/facilities` - Lấy tất cả cơ sở vật chất
- `PUT /api/facilities/:id` - Cập nhật cơ sở vật chất
- `POST /api/facilities/:id/maintenance` - Tạo hồ sơ bảo trì

### Chi Phí:
- `POST /api/expenses` - Tạo chi phí
- `GET /api/expenses` - Lấy chi phí
- `PUT /api/expenses/:id` - Cập nhật chi phí

### Người Dùng:
- `POST /api/users` - Tạo người dùng
- `GET /api/users` - Lấy tất cả người dùng
- `PUT /api/users/:id` - Cập nhật người dùng
- `POST /api/users/:id/reset-password` - Đặt lại mật khẩu

---

## Mô Hình Cơ Sở Dữ Liệu

### Collections:
- `users` - Người dùng
- `orders` - Đơn hàng
- `shifts` - Ca làm việc
- `cashier_shifts` - Ca thu ngân
- `menu_items` - Món trong thực đơn
- `ingredients` - Nguyên liệu
- `ingredient_categories` - Danh mục nguyên liệu
- `stock_history` - Lịch sử kho
- `facilities` - Cơ sở vật chất
- `facility_types` - Loại cơ sở vật chất
- `facility_areas` - Khu vực cơ sở vật chất
- `maintenance_records` - Hồ sơ bảo trì
- `expenses` - Chi phí
- `expense_categories` - Danh mục chi phí

---

## Quy Trình Kinh Doanh Chính

### 1. Quy Trình Xử Lý Đơn Hàng:
1. Phục vụ tạo đơn hàng
2. Phục vụ thu tiền
3. Phục vụ gửi đến quầy bar
4. Pha chế nhận đơn hàng
5. Pha chế chuẩn bị đồ uống
6. Pha chế đánh dấu sẵn sàng
7. Phục vụ phục vụ khách hàng
8. Thu ngân khóa đơn hàng

### 2. Đóng Ca Thu Ngân:
1. Thu ngân khởi tạo đóng ca
2. Thu ngân đếm tiền thực tế
3. Hệ thống tính chênh lệch
4. Thu ngân ghi chép chênh lệch (nếu cần)
5. Thu ngân xác nhận trách nhiệm
6. Hệ thống đóng ca

### 3. Quản Lý Kho:
1. Quản lý tạo nguyên liệu
2. Pha chế sử dụng nguyên liệu trong đơn hàng
3. Quản lý điều chỉnh kho
4. Hệ thống theo dõi lịch sử kho
5. Hệ thống cảnh báo hết hàng
6. Quản lý tạo đơn đặt hàng

### 4. Bảo Trì Cơ Sở Vật Chất:
1. Quản lý tạo cơ sở vật chất
2. Thực hiện bảo trì
3. Quản lý ghi nhận bảo trì
4. Hệ thống theo dõi lịch sử bảo trì
5. Hệ thống tính chi phí bảo trì
6. Quản lý xem báo cáo bảo trì

---

## Tính Năng Nâng Cao

### Xác Thực State Machine:
- Ngăn chuyển đổi đơn hàng không hợp lệ
- Thực thi quy tắc kinh doanh
- Tính toán tiến độ
- Xác định hành động tiếp theo

### Kiểm Toán:
- Lịch sử đơn hàng với thời gian
- Nhật ký kiểm toán ca thu ngân
- Lịch sử thay đổi cơ sở vật chất
- Lịch sử điều chỉnh kho
- Theo dõi hành động người dùng

### Theo Dõi Chi Phí Tự Động:
- Tự động tạo chi phí cho mua nguyên liệu
- Tự động tạo chi phí cho mua cơ sở vật chất
- Tự động tạo chi phí cho bảo trì
- Liên kết với hồ sơ nguồn

### Cập Nhật Thời Gian Thực:
- Cập nhật trạng thái đơn hàng
- Theo dõi ca
- Quản lý hàng đợi
- Theo dõi thanh toán

### Dashboard Theo Vai Trò:
- **Manager**: Tổng quan hệ thống
- **Cashier**: Theo dõi ca và thanh toán
- **Waiter**: Quản lý đơn hàng
- **Barista**: Theo dõi hàng đợi và chuẩn bị

---

## Triển Khai & Cấu Hình

### Biến Môi Trường:
- Chuỗi kết nối cơ sở dữ liệu
- JWT secret (tối thiểu 32 ký tự)
- Cổng API
- URL Frontend
- Cấu hình HTTPS

### Triển Khai Docker:
- Container Backend
- Container Frontend
- Container MongoDB
- Nginx reverse proxy

### Checklist Sản Xuất:
- Thay đổi mật khẩu mặc định
- Đặt JWT secret mạnh
- Bật HTTPS
- Cấu hình sao lưu
- Thiết lập giám sát
- Cấu hình logging

---

Hệ thống Café POS toàn diện này cung cấp khả năng quản lý hoàn chỉnh cho hoạt động quán cà phê bao gồm quản lý đơn hàng, xử lý thanh toán, theo dõi kho, quản lý cơ sở vật chất, theo dõi chi phí, và báo cáo chi tiết với kiểm soát truy cập dựa trên vai trò và quy trình điều khiển bằng state machine.