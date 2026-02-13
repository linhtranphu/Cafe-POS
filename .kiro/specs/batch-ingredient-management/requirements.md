# Tài Liệu Yêu Cầu: Quản Lý Nguyên Liệu Batch

## Giới Thiệu

Hệ thống quản lý nguyên liệu batch cho phép theo dõi các nguyên liệu trung gian được chế biến từ nguyên liệu thô và sau đó được sử dụng trong các sản phẩm cuối cùng. Ví dụ: cà phê concentrate được chế biến từ hạt cà phê, sau đó được sử dụng để pha các loại đồ uống cà phê.

Hệ thống này tích hợp với hệ thống quản lý nguyên liệu hiện có, hệ thống menu với công thức, và hệ thống tính toán chi phí để cung cấp khả năng theo dõi toàn diện từ nguyên liệu thô đến sản phẩm cuối cùng.

## Bảng Thuật Ngữ

- **Source_Ingredient**: Nguyên liệu thô được sử dụng để chế biến batch (ví dụ: hạt cà phê, sữa tươi)
- **Batch**: Nguyên liệu trung gian được chế biến từ Source_Ingredient (ví dụ: cà phê concentrate, sữa đã tiệt trùng)
- **Batch_Definition**: Định nghĩa cách chế biến một loại batch, bao gồm công thức và tỷ lệ chuyển đổi
- **Conversion_Rate**: Tỷ lệ chuyển đổi từ Source_Ingredient sang Batch (ví dụ: 100g hạt → 500ml concentrate)
- **Wastage_Rate**: Tỷ lệ hao hụt trong quá trình chế biến (phần trăm)
- **Shelf_Life**: Thời gian batch có thể sử dụng được sau khi chế biến (tính bằng giờ)
- **Batch_Record**: Bản ghi cụ thể của một lần chế biến batch
- **Low_Stock_Threshold**: Ngưỡng cảnh báo khi số lượng batch còn lại thấp
- **Expiry_Warning_Hours**: Số giờ trước khi hết hạn để hiển thị cảnh báo
- **Menu_Item**: Món ăn/đồ uống trong menu sử dụng batch
- **Recipe**: Công thức của Menu_Item, có thể sử dụng batch hoặc nguyên liệu thô
- **Cost_Calculator**: Hệ thống tính toán chi phí hiện có
- **Inventory_System**: Hệ thống quản lý tồn kho hiện có
- **Preparer**: Người chế biến batch

## Yêu Cầu

### Yêu Cầu 1: Định Nghĩa Batch

**User Story:** Là một quản lý, tôi muốn định nghĩa các loại batch có thể được chế biến, để hệ thống biết cách tính toán chi phí và theo dõi tồn kho.

#### Tiêu Chí Chấp Nhận

1. THE Batch_Definition SHALL bao gồm tên batch, đơn vị đo lường, và Shelf_Life
2. THE Batch_Definition SHALL chỉ định một hoặc nhiều Source_Ingredient với Conversion_Rate tương ứng
3. WHEN định nghĩa Conversion_Rate, THE System SHALL yêu cầu số lượng Source_Ingredient và số lượng Batch được tạo ra
4. THE Batch_Definition SHALL cho phép chỉ định Wastage_Rate cho mỗi Source_Ingredient
5. THE Batch_Definition SHALL cho phép thiết lập Low_Stock_Threshold và Expiry_Warning_Hours
6. WHEN tạo Batch_Definition, THE System SHALL xác thực rằng tất cả Source_Ingredient tồn tại trong Inventory_System

### Yêu Cầu 2: Ghi Nhận Chế Biến Batch

**User Story:** Là một barista, tôi muốn ghi nhận khi tôi chế biến một batch cà phê concentrate, để hệ thống theo dõi tồn kho và chi phí.

#### Tiêu Chí Chấp Nhận

1. WHEN ghi nhận chế biến batch, THE System SHALL tạo Batch_Record với số lượng, thời gian chế biến, và Preparer
2. WHEN tạo Batch_Record, THE System SHALL tự động trừ Source_Ingredient từ Inventory_System dựa trên Conversion_Rate
3. WHEN tạo Batch_Record, THE System SHALL áp dụng Wastage_Rate để tính toán số lượng Source_Ingredient thực tế cần sử dụng
4. WHEN tạo Batch_Record, THE System SHALL tính toán thời gian hết hạn dựa trên Shelf_Life
5. IF Source_Ingredient không đủ số lượng, THEN THE System SHALL từ chối tạo Batch_Record và hiển thị thông báo lỗi
6. WHEN tạo Batch_Record thành công, THE System SHALL cập nhật tổng số lượng batch khả dụng

### Yêu Cầu 3: Tính Toán Chi Phí Batch Tự Động

**User Story:** Là một quản lý, tôi muốn xem chi phí của mỗi batch được tính tự động từ nguyên liệu nguồn, để tôi có thể định giá chính xác.

#### Tiêu Chí Chấp Nhận

1. WHEN tạo Batch_Record, THE System SHALL tự động tính toán chi phí batch từ Cost_Calculator
2. THE System SHALL tính chi phí batch bằng cách cộng chi phí của tất cả Source_Ingredient đã sử dụng (bao gồm wastage)
3. WHEN hiển thị thông tin batch, THE System SHALL hiển thị chi phí trên mỗi đơn vị
4. WHEN chi phí Source_Ingredient thay đổi, THE System SHALL tính lại chi phí cho các Batch_Record mới
5. THE System SHALL lưu trữ chi phí tại thời điểm chế biến cho mỗi Batch_Record để theo dõi lịch sử

### Yêu Cầu 4: Cảnh Báo Tồn Kho Thấp và Hết Hạn

**User Story:** Là một barista, tôi muốn được cảnh báo khi batch sắp hết hoặc sắp hết hạn, để tôi có thể chế biến thêm hoặc sử dụng kịp thời.

#### Tiêu Chí Chấp Nhận

1. WHEN tổng số lượng batch khả dụng nhỏ hơn hoặc bằng Low_Stock_Threshold, THE System SHALL hiển thị cảnh báo tồn kho thấp
2. WHEN thời gian còn lại của Batch_Record nhỏ hơn hoặc bằng Expiry_Warning_Hours, THE System SHALL hiển thị cảnh báo sắp hết hạn
3. WHEN Batch_Record đã hết hạn, THE System SHALL đánh dấu batch là không khả dụng và không cho phép sử dụng
4. THE System SHALL hiển thị danh sách tất cả các cảnh báo trên dashboard
5. WHEN hiển thị danh sách batch, THE System SHALL sắp xếp theo thời gian hết hạn (batch sắp hết hạn nhất ở đầu)
6. THE System SHALL hiển thị số lượng batch khả dụng (chưa hết hạn) và số lượng batch đã hết hạn riêng biệt

### Yêu Cầu 5: Tích Hợp Menu và Trừ Tự Động

**User Story:** Là một quản lý, tôi muốn batch tự động được trừ khỏi tồn kho khi được sử dụng trong món ăn, để tôi có thể theo dõi chính xác.

#### Tiêu Chí Chấp Nhận

1. THE Recipe SHALL cho phép sử dụng batch như một thành phần thay vì Source_Ingredient
2. WHEN tạo đơn hàng với Menu_Item sử dụng batch, THE System SHALL tự động trừ số lượng batch từ Batch_Record
3. THE System SHALL trừ batch theo thứ tự FIFO (First In First Out) để sử dụng batch cũ nhất trước
4. WHEN trừ batch, THE System SHALL ưu tiên batch sắp hết hạn nhất trong số các batch khả dụng
5. IF không đủ batch khả dụng, THEN THE System SHALL từ chối đơn hàng và hiển thị thông báo lỗi
6. WHEN tính chi phí Menu_Item, THE System SHALL sử dụng chi phí batch thực tế đã được trừ

### Yêu Cầu 6: Báo Cáo và Theo Dõi

**User Story:** Là một quản lý, tôi muốn xem báo cáo về việc sử dụng batch, để tôi có thể tối ưu hóa quy trình chế biến.

#### Tiêu Chí Chấp Nhận

1. THE System SHALL cung cấp báo cáo hiển thị tổng số lượng batch đã chế biến theo khoảng thời gian
2. THE System SHALL hiển thị tỷ lệ wastage thực tế so với Wastage_Rate dự kiến
3. THE System SHALL hiển thị số lượng batch đã hết hạn và giá trị chi phí bị lãng phí
4. THE System SHALL hiển thị batch nào được sử dụng nhiều nhất trong Menu_Item
5. THE System SHALL cho phép lọc báo cáo theo loại batch, Preparer, và khoảng thời gian
6. THE System SHALL hiển thị xu hướng sử dụng batch theo thời gian để hỗ trợ dự đoán nhu cầu

### Yêu Cầu 7: Quản Lý Batch Record

**User Story:** Là một quản lý, tôi muốn có thể điều chỉnh hoặc xóa batch record khi có sai sót, để dữ liệu luôn chính xác.

#### Tiêu Chí Chấp Nhận

1. THE System SHALL cho phép xem danh sách tất cả Batch_Record với thông tin chi tiết
2. WHEN xóa Batch_Record chưa được sử dụng, THE System SHALL hoàn trả Source_Ingredient vào Inventory_System
3. IF Batch_Record đã được sử dụng một phần, THEN THE System SHALL chỉ cho phép đánh dấu là hết hạn, không cho phép xóa
4. THE System SHALL ghi lại lịch sử thay đổi cho mỗi Batch_Record (audit log)
5. WHEN điều chỉnh số lượng Batch_Record, THE System SHALL tính lại chi phí và cập nhật Inventory_System tương ứng
6. THE System SHALL yêu cầu xác nhận trước khi xóa hoặc điều chỉnh Batch_Record

### Yêu Cầu 8: API và Tích Hợp Backend

**User Story:** Là một developer, tôi muốn có API rõ ràng để quản lý batch, để frontend có thể tích hợp dễ dàng.

#### Tiêu Chí Chấp Nhận

1. THE System SHALL cung cấp RESTful API cho tất cả các thao tác batch (CRUD)
2. THE API SHALL sử dụng JSON cho request và response
3. THE API SHALL trả về mã lỗi HTTP phù hợp (200, 201, 400, 404, 500)
4. THE API SHALL xác thực quyền truy cập cho mỗi endpoint
5. THE API SHALL validate dữ liệu đầu vào và trả về thông báo lỗi rõ ràng
6. THE System SHALL lưu trữ dữ liệu batch trong MongoDB với schema được định nghĩa rõ ràng
7. THE System SHALL đảm bảo tính toàn vẹn dữ liệu khi cập nhật đồng thời (concurrency control)

### Yêu Cầu 9: Giao Diện Người Dùng

**User Story:** Là một barista, tôi muốn giao diện đơn giản để ghi nhận batch nhanh chóng, để không làm gián đoạn công việc.

#### Tiêu Chí Chấp Nhận

1. THE Frontend SHALL cung cấp form đơn giản để ghi nhận chế biến batch với các trường bắt buộc
2. THE Frontend SHALL hiển thị danh sách batch khả dụng với trạng thái tồn kho và thời gian hết hạn
3. THE Frontend SHALL sử dụng màu sắc để phân biệt trạng thái (xanh: đủ, vàng: thấp, đỏ: sắp hết hạn)
4. THE Frontend SHALL hiển thị cảnh báo nổi bật trên dashboard khi có batch cần chú ý
5. THE Frontend SHALL responsive và hoạt động tốt trên mobile device
6. WHEN ghi nhận batch, THE Frontend SHALL hiển thị preview chi phí dự kiến trước khi xác nhận

### Yêu Cầu 10: Hiệu Suất và Khả Năng Mở Rộng

**User Story:** Là một system administrator, tôi muốn hệ thống hoạt động nhanh và ổn định ngay cả khi có nhiều batch record, để đảm bảo trải nghiệm người dùng tốt.

#### Tiêu Chí Chấp Nhận

1. WHEN truy vấn danh sách batch, THE System SHALL trả về kết quả trong vòng 500ms với 1000 batch records
2. THE System SHALL sử dụng index trên MongoDB cho các trường thường xuyên truy vấn (batch_definition_id, expiry_time, created_at)
3. WHEN tính toán chi phí batch, THE System SHALL cache chi phí Source_Ingredient để tránh truy vấn lặp lại
4. THE System SHALL xử lý đồng thời nhiều request tạo batch mà không gây ra race condition
5. THE System SHALL có khả năng xử lý ít nhất 100 batch records được tạo đồng thời
6. THE System SHALL log các thao tác quan trọng để hỗ trợ debugging và monitoring
