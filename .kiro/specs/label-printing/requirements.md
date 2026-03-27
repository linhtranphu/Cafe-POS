# Tài Liệu Yêu Cầu: In Tem Đơn Hàng (Label Printing)

## Giới Thiệu

Chức năng in tem thay thế chức năng in bill tạm hiện tại. Khi người dùng bấm nút "In bill tạm" (📄), hệ thống sẽ in tem nhãn cho từng món trong đơn hàng thay vì in bill HTML. Tem sử dụng giao thức TSPL (TSC Printer Language) và hiển thị thông tin: tên món, thời gian, tên khách hàng, số lượng, và ghi chú.

## Thuật Ngữ (Glossary)

- **System**: Hệ thống quản lý quán cafe (Cafe POS)
- **TSPL_Generator**: Component chịu trách nhiệm tạo lệnh TSPL từ template và dữ liệu
- **Print_Bridge**: Service trung gian giữa backend và máy in tem
- **Label_Printer**: Máy in tem vật lý hỗ trợ giao thức TSPL
- **Label_Template**: Template chứa các lệnh TSPL với placeholders
- **Order**: Đơn hàng chứa thông tin khách hàng và danh sách món
- **Order_Item**: Một món trong đơn hàng
- **Shop_Settings**: Cấu hình của cửa hàng bao gồm thông tin máy in
- **Template_Editor**: Giao diện cho phép chỉnh sửa template TSPL

## Yêu Cầu

### Requirement 1: Tạo và Gửi Lệnh TSPL

**User Story:** Là một nhân viên quán, tôi muốn in tem cho từng món trong đơn hàng, để bếp và bar có thể xác định món cần làm cho khách nào.

#### Acceptance Criteria

1. WHEN người dùng bấm nút "In bill tạm" (📄) trên giao diện đơn hàng, THE System SHALL gọi API endpoint `/api/orders/:id/print-temp-bill`
2. WHEN API endpoint `/api/orders/:id/print-temp-bill` được gọi, THE System SHALL tạo lệnh TSPL cho từng món trong đơn hàng thay vì tạo bill HTML
3. WHEN tạo lệnh TSPL cho một món, THE TSPL_Generator SHALL load template từ file `label_template.tspl`
4. WHEN tạo lệnh TSPL, THE TSPL_Generator SHALL thay thế các placeholders trong template bằng dữ liệu thực tế: mã order, tên món, note, thời gian, tên khách hàng
5. WHEN lệnh TSPL được tạo xong, THE System SHALL gửi lệnh đến Print_Bridge qua HTTP POST request đến endpoint `/print-tspl`

### Requirement 2: Xử Lý Dữ Liệu Tem

**User Story:** Là một nhân viên quán, tôi muốn tem hiển thị đầy đủ thông tin món ăn theo thứ tự: mã order, tên món + note, thời gian, tên khách hàng, để dễ dàng xác định và chế biến.

#### Acceptance Criteria

1. THE System SHALL hiển thị thông tin trên tem theo thứ tự: (1) Mã order, (2) Tên món, (3) Note (nếu có), (4) Thời gian, (5) Tên khách hàng
2. WHEN một món có variant (ví dụ: "Size L"), THE System SHALL hiển thị tên món kèm variant theo format "Tên món (Variant)"
3. WHEN một món có note, THE System SHALL hiển thị note trong ngoặc đơn ngay dưới tên món: "(Note)"
4. WHEN tên món vượt quá 20 ký tự, THE System SHALL cắt ngắn text và thêm "..." ở cuối
5. WHEN tên khách hàng vượt quá 15 ký tự, THE System SHALL cắt ngắn text và thêm "..." ở cuối
6. WHEN ghi chú vượt quá 25 ký tự, THE System SHALL cắt ngắn text và thêm "..." ở cuối
7. WHEN cắt ngắn text, THE System SHALL đếm theo số ký tự Unicode (runes) chứ không phải bytes, để hỗ trợ tiếng Việt đúng cách

### Requirement 3: Quản Lý Template TSPL

**User Story:** Là một quản lý quán, tôi muốn tùy chỉnh template tem, để phù hợp với kích thước tem và nhu cầu hiển thị của quán.

#### Acceptance Criteria

1. THE System SHALL cung cấp API endpoint `GET /api/label-templates/order-item` để lấy nội dung template hiện tại
2. THE System SHALL cung cấp API endpoint `PUT /api/label-templates/order-item` để cập nhật template
3. WHEN người dùng cập nhật template, THE System SHALL tạo file backup với suffix `.backup` trước khi lưu template mới
4. WHEN template file không tồn tại, THE System SHALL tạo template mặc định với kích thước 50mm x 30mm
5. THE System SHALL cung cấp giao diện Template_Editor cho phép chỉnh sửa template TSPL trực tiếp

### Requirement 4: Cấu Hình Máy In Tem

**User Story:** Là một quản lý quán, tôi muốn cấu hình thông tin máy in tem, để hệ thống có thể kết nối và in tem đúng cách.

#### Acceptance Criteria

1. THE System SHALL lưu trữ cấu hình máy in tem trong Shop_Settings bao gồm: `label_printer_enabled`, `label_printer_ip`, `label_printer_port`, `label_width`, `label_height`
2. THE System SHALL cung cấp giao diện trong Settings để bật/tắt máy in tem
3. THE System SHALL cung cấp giao diện trong Settings để nhập IP address và port của máy in tem
4. THE System SHALL cung cấp dropdown để chọn kích thước tem: 40x30mm, 50x30mm, 60x40mm
5. WHEN `label_printer_enabled` là false, THE System SHALL không thực hiện in tem và trả về lỗi thông báo "label printer not enabled"

### Requirement 5: Kiểm Tra Kết Nối và Xử Lý Lỗi

**User Story:** Là một nhân viên quán, tôi muốn biết ngay khi có lỗi in tem, để có thể xử lý kịp thời.

#### Acceptance Criteria

1. WHEN thực hiện in tem, THE System SHALL kiểm tra Print_Bridge có đang chạy và accessible hay không trước khi gửi lệnh
2. IF Print_Bridge không available, THEN THE System SHALL trả về lỗi với message "print bridge not available"
3. IF máy in tem không kết nối được (offline hoặc sai IP), THEN THE System SHALL trả về lỗi với message "label printer offline or unreachable"
4. IF template file không tồn tại hoặc không đọc được, THEN THE System SHALL trả về lỗi với message "failed to read template"
5. IF template có syntax lỗi, THEN THE System SHALL trả về lỗi với message "invalid template syntax"
6. WHEN có lỗi xảy ra trong quá trình in tem, THE System SHALL log chi tiết lỗi bao gồm order ID, item name, và error message

### Requirement 6: In Nhiều Tem Cho Một Đơn Hàng

**User Story:** Là một nhân viên quán, tôi muốn mỗi món trong đơn hàng có một tem riêng, để dễ dàng phân phối cho các bộ phận khác nhau (bếp, bar).

#### Acceptance Criteria

1. WHEN một đơn hàng có N món, THE System SHALL tạo và in đúng N tem (một tem cho mỗi món)
2. WHEN in tem cho một món thất bại, THE System SHALL tiếp tục in tem cho các món còn lại thay vì dừng toàn bộ quá trình
3. WHEN in tem cho một món thất bại, THE System SHALL log lỗi cho món đó nhưng không trả về lỗi cho toàn bộ request
4. THE System SHALL in tem tuần tự (sequential) cho từng món, không in song song (parallel)

### Requirement 7: Test Print

**User Story:** Là một quản lý quán, tôi muốn test in tem với dữ liệu mẫu, để kiểm tra máy in và template hoạt động đúng trước khi sử dụng thực tế.

#### Acceptance Criteria

1. THE System SHALL cung cấp API endpoint `POST /api/label-templates/test-print` để test in tem
2. WHEN gọi test print API, THE System SHALL chấp nhận dữ liệu mẫu bao gồm: `item_name`, `note`, `customer_name`, `printer_ip`, `port`
3. WHEN test print, THE System SHALL tạo lệnh TSPL với dữ liệu mẫu và gửi đến máy in được chỉ định
4. WHEN test print thành công, THE System SHALL trả về response với message "Test print successful"
5. WHEN test print thất bại, THE System SHALL trả về lỗi với message chi tiết về nguyên nhân

### Requirement 8: Thay Thế Chức Năng In Bill Tạm

**User Story:** Là một quản lý quán, tôi muốn chức năng in tem thay thế hoàn toàn chức năng in bill tạm cũ, để đơn giản hóa quy trình và tránh nhầm lẫn.

#### Acceptance Criteria

1. THE System SHALL thay đổi logic của method `CreateTempBillJob()` từ render HTML sang generate TSPL
2. THE System SHALL xóa hoặc không sử dụng file template `temp_bill_template.html`
3. THE System SHALL thay thế tab "Temp Bill Template" trong Print Management bằng tab "Label Template"
4. THE System SHALL xóa component `TempBillTemplateEditor.vue` và thay bằng `LabelTemplateEditor.vue`
5. THE System SHALL giữ nguyên button "📄 In bill tạm" trên giao diện nhưng thay đổi message thành công thành "✅ Đã in tem thành công"
6. THE System SHALL giữ nguyên API endpoint `/api/orders/:id/print-temp-bill` để không break frontend

### Requirement 9: Validation và Security

**User Story:** Là một quản lý hệ thống, tôi muốn đảm bảo template và cấu hình máy in được validate, để tránh lỗi bảo mật và hoạt động không mong muốn.

#### Acceptance Criteria

1. WHEN người dùng cập nhật template, THE System SHALL kiểm tra template không chứa các lệnh nguy hiểm như "DOWNLOAD", "ERASE", "KILL"
2. WHEN người dùng cập nhật template, THE System SHALL kiểm tra kích thước template không vượt quá 10000 bytes
3. WHEN người dùng nhập IP máy in, THE System SHALL validate IP address có format hợp lệ
4. WHEN người dùng nhập port máy in, THE System SHALL validate port number là số dương hợp lệ
5. THE System SHALL yêu cầu authentication cho tất cả API endpoints liên quan đến label printing
6. THE System SHALL chỉ cho phép user có role "admin" hoặc "manager" chỉnh sửa template và cấu hình máy in

### Requirement 10: Tích Hợp với Print Bridge

**User Story:** Là một developer, tôi muốn hệ thống tích hợp với Print Bridge một cách đáng tin cậy, để đảm bảo lệnh TSPL được gửi đến máy in đúng cách.

#### Acceptance Criteria

1. THE System SHALL gửi lệnh TSPL đến Print_Bridge qua HTTP POST request với Content-Type là "application/json"
2. THE System SHALL gửi request body chứa: `commands` (TSPL string), `printer_ip`, và `port`
3. WHEN Print_Bridge trả về HTTP status 200, THE System SHALL coi như in tem thành công
4. WHEN Print_Bridge trả về HTTP status khác 200, THE System SHALL coi như in tem thất bại và trả về lỗi
5. THE System SHALL set timeout 30 giây cho mỗi request đến Print_Bridge
