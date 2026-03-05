# Tài liệu Yêu cầu: Thiết kế lại Template In Bill

## Giới thiệu

Hệ thống hiện tại đã có khả năng in hóa đơn (bill) tự động sử dụng máy in nhiệt ESC/POS với template có thể tùy chỉnh. Tuy nhiên, template hiện tại cần được cải thiện về mặt bố cục và tính nhất quán. Feature này sẽ thiết kế lại template in bill để có bố cục gọn gàng hơn với các món được tổ chức trong bảng, font size đồng đều, và thêm logo ở góc trên bên trái.

## Thuật ngữ

- **Bill**: Hóa đơn bán hàng được in ra cho khách hàng
- **Template**: Mẫu định dạng nội dung hóa đơn
- **ESC/POS**: Giao thức lệnh cho máy in nhiệt
- **Text_Renderer**: Module chuyển đổi text thành hình ảnh để in
- **Format_Parser**: Module phân tích và định dạng nội dung template
- **Logo**: Hình ảnh đại diện thương hiệu của cửa hàng
- **Item_Table**: Bảng hiển thị danh sách các món trong đơn hàng
- **Font_Size**: Kích thước chữ được đo bằng điểm (pt)
- **Paper_Width**: Chiều rộng giấy in (58mm hoặc 80mm)

## Yêu cầu

### Yêu cầu 1: Hiển thị Logo ở góc trên bên trái

**User Story:** Là chủ cửa hàng, tôi muốn có logo thương hiệu xuất hiện ở góc trên bên trái của hóa đơn, để tăng tính nhận diện thương hiệu và tính chuyên nghiệp.

#### Tiêu chí chấp nhận

1. WHEN một hóa đơn được in THEN THE System SHALL hiển thị logo ở góc trên bên trái của hóa đơn
2. WHERE logo được cấu hình THEN THE System SHALL tải và hiển thị logo từ đường dẫn đã cấu hình
3. WHEN logo không được cấu hình THEN THE System SHALL bỏ qua phần logo và in phần còn lại bình thường
4. THE System SHALL thay đổi kích thước logo để phù hợp với chiều rộng giấy in (tối đa 25% chiều rộng)
5. THE System SHALL căn chỉnh logo về phía trái với khoảng cách padding phù hợp
6. THE System SHALL hỗ trợ các định dạng hình ảnh phổ biến (PNG, JPG, JPEG)
7. WHEN logo có kích thước quá lớn THEN THE System SHALL tự động thu nhỏ logo để vừa với giấy in
8. THE System SHALL chuyển đổi logo sang ảnh đen trắng (grayscale) để phù hợp với máy in nhiệt

### Yêu cầu 2: Tổ chức các món trong bảng

**User Story:** Là người dùng, tôi muốn các món trong hóa đơn được hiển thị trong bảng có cấu trúc rõ ràng, để dễ đọc và theo dõi thông tin.

#### Tiêu chí chấp nhận

1. THE System SHALL hiển thị các món trong định dạng bảng với các cột: Tên món, Số lượng, Đơn giá, Thành tiền
2. THE System SHALL căn chỉnh tên món về bên trái
3. THE System SHALL căn chỉnh số lượng, đơn giá, và thành tiền về bên phải
4. THE System SHALL sử dụng đường kẻ ngang để phân tách header bảng và các dòng món
5. THE System SHALL tự động xuống dòng cho tên món dài vượt quá chiều rộng cột
6. WHEN một món có variant THEN THE System SHALL hiển thị variant trên dòng phụ bên dưới tên món
7. THE System SHALL tính toán độ rộng các cột dựa trên chiều rộng giấy in
8. THE System SHALL đảm bảo các cột được căn chỉnh đồng đều trong toàn bộ bảng

### Yêu cầu 3: Font size đồng đều

**User Story:** Là người đọc hóa đơn, tôi muốn tất cả nội dung có font size nhất quán, để hóa đơn trông chuyên nghiệp và dễ đọc.

#### Tiêu chí chấp nhận

1. THE System SHALL sử dụng font size 18pt cho tất cả nội dung thông thường trong hóa đơn
2. THE System SHALL sử dụng font size 22pt cho tiêu đề chính (tên cửa hàng, "HÓA ĐƠN BÁN HÀNG")
3. THE System SHALL sử dụng font size 20pt cho tổng tiền cuối cùng (TỔNG CỘNG)
4. THE System SHALL đảm bảo font size nhất quán cho tất cả các dòng trong bảng món
5. THE System SHALL đảm bảo font size nhất quán cho thông tin đơn hàng (số order, ngày giờ, bàn, khách)
6. THE System SHALL sử dụng font weight bold cho header bảng và tổng tiền
7. THE System SHALL sử dụng font weight normal cho nội dung thông thường

### Yêu cầu 4: Tích hợp với hệ thống template hiện có

**User Story:** Là quản trị viên, tôi muốn template mới tích hợp mượt mà với hệ thống in hiện có, để không làm gián đoạn quy trình làm việc.

#### Tiêu chí chấp nhận

1. THE System SHALL lưu template mới vào collection print_templates với type BILL
2. THE System SHALL cho phép đặt template mới làm template mặc định
3. THE System SHALL sử dụng Text_Renderer hiện có để render template
4. THE System SHALL sử dụng Format_Parser hiện có để phân tích định dạng
5. THE System SHALL hỗ trợ tất cả các biến template hiện có (ShopName, Order.Items, Order.Total, v.v.)
6. THE System SHALL tương thích với cả giấy in 58mm và 80mm
7. WHEN template mới được chọn làm mặc định THEN THE System SHALL sử dụng template này cho tất cả các lần in bill mới

### Yêu cầu 5: Cấu hình logo trong shop settings

**User Story:** Là quản trị viên, tôi muốn có thể cấu hình logo từ giao diện quản lý, để dễ dàng thay đổi logo khi cần.

#### Tiêu chí chấp nhận

1. THE System SHALL cung cấp trường upload logo trong shop settings
2. THE System SHALL lưu đường dẫn logo vào shop_settings collection
3. THE System SHALL hiển thị preview logo sau khi upload
4. THE System SHALL cho phép xóa logo đã upload
5. WHEN logo được upload THEN THE System SHALL validate định dạng file (PNG, JPG, JPEG)
6. WHEN logo được upload THEN THE System SHALL validate kích thước file (tối đa 2MB)
7. THE System SHALL lưu logo vào thư mục uploads/logos/ trên server

### Yêu cầu 6: Tương thích ngược với template cũ

**User Story:** Là quản trị viên, tôi muốn vẫn có thể sử dụng template cũ nếu cần, để có sự linh hoạt trong việc chọn lựa.

#### Tiêu chí chấp nhận

1. THE System SHALL giữ nguyên template cũ trong database
2. THE System SHALL cho phép chuyển đổi giữa template mới và cũ
3. WHEN chuyển về template cũ THEN THE System SHALL in hóa đơn theo định dạng cũ
4. THE System SHALL không xóa hoặc ghi đè template cũ khi tạo template mới
5. THE System SHALL hiển thị danh sách tất cả templates có sẵn trong Print Management

### Yêu cầu 7: Xử lý lỗi và fallback

**User Story:** Là hệ thống, tôi cần xử lý các lỗi liên quan đến logo và template một cách graceful, để đảm bảo việc in không bị gián đoạn.

#### Tiêu chí chấp nhận

1. WHEN logo file không tồn tại THEN THE System SHALL in hóa đơn mà không có logo
2. WHEN logo file bị lỗi (corrupt) THEN THE System SHALL log lỗi và in hóa đơn mà không có logo
3. WHEN template rendering thất bại THEN THE System SHALL fallback về template mặc định cũ
4. THE System SHALL log tất cả lỗi liên quan đến logo và template rendering
5. WHEN logo quá lớn và không thể resize THEN THE System SHALL bỏ qua logo và in phần còn lại
