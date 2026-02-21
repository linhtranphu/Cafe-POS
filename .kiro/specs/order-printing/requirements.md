# Tài Liệu Yêu Cầu: In Bill và Tem Đơn Hàng

## Giới Thiệu

Chức năng in bill và tem cho phép quán in hóa đơn thanh toán (bill) để đưa cho khách hàng và in nhãn dán (tem) để dán lên sản phẩm sau khi đơn hàng được tạo. Điều này giúp tối ưu hóa quy trình phục vụ, đảm bảo thông tin chính xác và tăng tốc độ xử lý đơn hàng.

## Thuật Ngữ

- **System**: Hệ thống quản lý quán
- **Bill**: Hóa đơn thanh toán được in ra để đưa cho khách hàng
- **Tem**: Nhãn dán được in ra để dán lên ly/sản phẩm
- **Order**: Đơn hàng được tạo trong hệ thống
- **Order_Item**: Món trong đơn hàng (bao gồm tên món, số lượng, giá, variant)
- **Printer**: Thiết bị in (máy in bill hoặc máy in tem)
- **Print_Job**: Công việc in được gửi đến máy in
- **Variant**: Biến thể của món (ví dụ: size S/M/L, đá/ít đá)

## Yêu Cầu

### Yêu Cầu 1: In Bill Thanh Toán

**User Story:** Là nhân viên quán, tôi muốn in bill thanh toán sau khi tạo đơn hàng, để khách hàng có thể biết tổng số tiền cần thanh toán và giữ làm chứng từ.

#### Tiêu Chí Chấp Nhận

1. WHEN một Order được tạo thành công và được xác nhận, THE System SHALL tự động tạo Print_Job cho Bill
1a. THE System SHALL in Bill ngay sau khi Order được tạo và trước khi bắt đầu pha chế
1b. THE System SHALL cho phép in lại Bill bất kỳ lúc nào từ chi tiết đơn hàng
2. THE Bill SHALL hiển thị tên quán và thông tin liên hệ
3. THE Bill SHALL hiển thị số thứ tự Order
4. THE Bill SHALL hiển thị thời gian tạo Order
5. FOR EACH Order_Item trong Order, THE Bill SHALL hiển thị tên món, variant (nếu có), số lượng, đơn giá và thành tiền
6. THE Bill SHALL hiển thị tổng số tiền cần thanh toán
7. THE Bill SHALL có định dạng phù hợp với máy in nhiệt (thermal printer) khổ 58mm hoặc 80mm
8. WHEN Print_Job được tạo, THE System SHALL gửi lệnh in đến Printer được cấu hình ngay lập tức
9. IF Printer không khả dụng, THEN THE System SHALL lưu Print_Job vào hàng đợi và thử lại sau

### Yêu Cầu 2: In Tem Sản Phẩm

**User Story:** Là nhân viên pha chế, tôi muốn in tem dán lên sản phẩm, để đảm bảo món được làm đúng yêu cầu và giao đúng khách hàng.

#### Tiêu Chí Chấp Nhận

1. WHEN một Order được tạo thành công, THE System SHALL tự động tạo Print_Job cho Tem của mỗi Order_Item
2. FOR EACH Order_Item, THE System SHALL tạo một Tem riêng biệt
3. THE Tem SHALL hiển thị số thứ tự Order
4. THE Tem SHALL hiển thị tên món
5. THE Tem SHALL hiển thị variant của món (ví dụ: Size M, Đá, Ít đường)
6. THE Tem SHALL hiển thị số thứ tự món trong Order (ví dụ: 1/3, 2/3, 3/3)
7. THE Tem SHALL hiển thị thời gian tạo Order
8. THE Tem SHALL có định dạng phù hợp với máy in tem nhãn (label printer)
9. WHEN Print_Job được tạo, THE System SHALL gửi lệnh in đến Printer được cấu hình
10. IF Printer không khả dụng, THEN THE System SHALL lưu Print_Job vào hàng đợi và thử lại sau

### Yêu Cầu 3: Quản Lý Cấu Hình Máy In

**User Story:** Là quản lý quán, tôi muốn cấu hình các máy in trong hệ thống, để hệ thống biết in bill và tem ra máy nào.

#### Tiêu Chí Chấp Nhận

1. THE System SHALL cho phép cấu hình nhiều Printer
2. FOR EACH Printer, THE System SHALL lưu trữ tên, loại (bill hoặc tem), địa chỉ IP hoặc cổng kết nối
3. THE System SHALL cho phép chọn Printer mặc định cho bill
4. THE System SHALL cho phép chọn Printer mặc định cho tem
5. THE System SHALL kiểm tra trạng thái kết nối của Printer
6. WHEN Printer được thêm hoặc cập nhật, THE System SHALL xác thực thông tin kết nối
7. THE System SHALL cho phép vô hiệu hóa Printer mà không xóa cấu hình

### Yêu Cầu 4: Quản Lý Hàng Đợi In

**User Story:** Là nhân viên quán, tôi muốn xem và quản lý các công việc in đang chờ, để xử lý khi máy in gặp sự cố.

#### Tiêu Chí Chấp Nhận

1. THE System SHALL lưu trữ tất cả Print_Job với trạng thái (pending, printing, completed, failed)
2. WHEN Print_Job thất bại, THE System SHALL tự động thử lại tối đa 3 lần
3. IF Print_Job thất bại sau 3 lần thử, THEN THE System SHALL đánh dấu trạng thái là failed
4. THE System SHALL hiển thị danh sách Print_Job đang chờ và thất bại
5. THE System SHALL cho phép in lại Print_Job thất bại thủ công
6. THE System SHALL cho phép hủy Print_Job đang chờ
7. THE System SHALL tự động xóa Print_Job đã hoàn thành sau 7 ngày

### Yêu Cầu 5: Tùy Chỉnh Mẫu In

**User Story:** Là quản lý quán, tôi muốn tùy chỉnh nội dung và định dạng của bill và tem, để phù hợp với thương hiệu và nhu cầu của quán.

#### Tiêu Chí Chấp Nhận

1. THE System SHALL cho phép cấu hình thông tin quán (tên, địa chỉ, số điện thoại) hiển thị trên Bill
2. THE System SHALL cho phép bật/tắt hiển thị logo trên Bill
3. THE System SHALL cho phép thêm lời cảm ơn hoặc thông điệp tùy chỉnh ở cuối Bill
4. THE System SHALL cho phép chọn kích thước giấy in (58mm hoặc 80mm) cho Bill
5. THE System SHALL cho phép chọn kích thước tem (40x30mm, 50x30mm, 60x40mm)
6. THE System SHALL cho phép bật/tắt hiển thị các trường thông tin trên Tem
7. WHEN cấu hình mẫu in được thay đổi, THE System SHALL áp dụng cho các đơn hàng mới

### Yêu Cầu 6: Tích Hợp Với Hệ Thống Order

**User Story:** Là nhân viên quán, tôi muốn hệ thống tự động in bill và tem khi tạo đơn hàng, để tiết kiệm thời gian và giảm sai sót.

#### Tiêu Chí Chấp Nhận

1. WHEN một Order được tạo thành công, THE System SHALL tự động kích hoạt in bill và tem
2. THE System SHALL đảm bảo Order được lưu vào database trước khi tạo Print_Job
3. IF việc tạo Print_Job thất bại, THEN THE System SHALL không ảnh hưởng đến việc tạo Order
4. THE System SHALL ghi log mỗi lần in bill và tem
5. THE System SHALL cho phép bật/tắt tính năng tự động in trong cài đặt
6. WHERE tính năng tự động in bị tắt, THE System SHALL vẫn cho phép in thủ công từ chi tiết đơn hàng

### Yêu Cầu 7: Xử Lý Lỗi và Thông Báo

**User Story:** Là nhân viên quán, tôi muốn được thông báo khi có lỗi in ấn, để kịp thời xử lý và không làm gián đoạn phục vụ.

#### Tiêu Chí Chấp Nhận

1. WHEN Print_Job thất bại, THE System SHALL hiển thị thông báo lỗi cho người dùng
2. THE System SHALL ghi log chi tiết lỗi bao gồm thời gian, loại lỗi, và thông tin Printer
3. IF Printer mất kết nối, THEN THE System SHALL hiển thị cảnh báo trên giao diện
4. THE System SHALL cho phép xem lịch sử lỗi in ấn
5. WHEN Printer hết giấy hoặc kẹt giấy, THE System SHALL phát hiện và thông báo (nếu Printer hỗ trợ)
6. THE System SHALL cung cấp hướng dẫn xử lý lỗi phổ biến trong giao diện
