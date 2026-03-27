# Requirements Document

## Introduction

Tính năng Tách Bill (Split Bill) cho phép waiter tách một order đang ở trạng thái CREATED thành hai order riêng biệt. Waiter chọn những items (và số lượng) muốn tách ra, phần đó sẽ thành order mới, phần còn lại ở order gốc. Cả hai order đều thuộc cùng shift và waiter. Tính năng này là nghịch đảo của Gộp Bill (Merge Orders) và không ảnh hưởng đến batch ingredients đã được deduct khi tạo order ban đầu.

## Glossary

- **Order**: Đơn hàng trong hệ thống POS, chứa danh sách items, thông tin thanh toán và trạng thái.
- **OrderItem**: Một dòng sản phẩm trong order, bao gồm tên, variant, giá, số lượng và ghi chú.
- **Split_Service**: Service xử lý logic tách bill.
- **Order_Repository**: Repository lưu trữ và truy vấn orders.
- **Shift**: Ca làm việc hiện tại, xác định phạm vi hoạt động của waiter.
- **Waiter**: Nhân viên phục vụ thực hiện thao tác tách bill.
- **SplitItem**: Thông tin một item được chọn để tách, gồm MenuItemID, VariantID và số lượng muốn tách.
- **Original_Order**: Order gốc được tách, sau khi tách sẽ chứa phần items còn lại.
- **New_Order**: Order mới được tạo ra chứa các items đã tách.

## Requirements

### Requirement 1: Điều kiện tách bill

**User Story:** As a waiter, I want to only split orders that are in CREATED status, so that I don't accidentally split orders that are already being processed or paid.

#### Acceptance Criteria

1. WHEN a split request is received, THE Split_Service SHALL verify the order exists before processing.
2. WHEN a split request is received for an order with status other than CREATED, THE Split_Service SHALL reject the request with an error message indicating the order cannot be split.
3. WHEN a split request is received for an order belonging to a closed Shift, THE Split_Service SHALL reject the request with an error message indicating the shift is closed.
4. WHEN a split request is received for an order with fewer than 2 total item units across all OrderItems, THE Split_Service SHALL reject the request with an error message indicating there are not enough items to split.

---

### Requirement 2: Chọn items để tách

**User Story:** As a waiter, I want to select specific items and quantities to move to a new order, so that I can flexibly split a bill according to customer requests.

#### Acceptance Criteria

1. WHEN submitting a split request, THE Waiter SHALL provide a list of SplitItems, each specifying a MenuItemID, VariantID, and quantity to split.
2. WHEN a SplitItem references a MenuItemID and VariantID combination not present in the Original_Order, THE Split_Service SHALL reject the request with an error identifying the invalid item.
3. WHEN a SplitItem specifies a quantity greater than the available quantity of that item in the Original_Order, THE Split_Service SHALL reject the request with an error indicating the quantity exceeds available stock.
4. WHEN a SplitItem specifies a quantity of zero or a negative number, THE Split_Service SHALL reject the request with an error indicating the quantity is invalid.
5. WHEN the selected SplitItems cover all item units in the Original_Order (nothing remains), THE Split_Service SHALL reject the request with an error indicating at least one item must remain in the original order.
6. THE Split_Service SHALL support partial quantity splits, where a SplitItem quantity is less than the full quantity of that item in the Original_Order.

---

### Requirement 3: Tạo order mới từ items được tách

**User Story:** As a waiter, I want the split items to form a new valid order, so that each group of customers can be billed separately.

#### Acceptance Criteria

1. WHEN a split is executed, THE Split_Service SHALL create a New_Order containing exactly the items and quantities specified in the SplitItems list.
2. THE Split_Service SHALL assign the New_Order the same ShiftID, WaiterID, and WaiterName as the Original_Order.
3. THE Split_Service SHALL set the New_Order status to CREATED.
4. THE Split_Service SHALL generate a unique OrderNumber for the New_Order using the same format as existing order creation.
5. WHERE a CustomerName is provided in the split request, THE Split_Service SHALL assign that CustomerName to the New_Order; otherwise THE Split_Service SHALL leave CustomerName empty.
6. THE Split_Service SHALL calculate Subtotal, Discount (0), and Total for the New_Order based on the items and their prices.
7. THE Split_Service SHALL set BillPrinted to false on the New_Order.

---

### Requirement 4: Cập nhật order gốc sau khi tách

**User Story:** As a waiter, I want the original order to be updated to reflect only the remaining items, so that the original bill is accurate after splitting.

#### Acceptance Criteria

1. WHEN a split is executed, THE Split_Service SHALL remove the split quantities from the corresponding items in the Original_Order.
2. WHEN a split reduces an OrderItem's quantity to zero, THE Split_Service SHALL remove that OrderItem entirely from the Original_Order.
3. THE Split_Service SHALL recalculate Subtotal and Total of the Original_Order after removing the split items.
4. THE Split_Service SHALL preserve the Original_Order's OrderNumber, WaiterID, WaiterName, ShiftID, CustomerName, Note, and Discount unchanged.
5. THE Split_Service SHALL set BillPrinted to false on the Original_Order after modification.

---

### Requirement 5: Tính toàn vẹn dữ liệu sau khi tách

**User Story:** As a system, I want the split operation to maintain data integrity, so that no items are lost or duplicated.

#### Acceptance Criteria

1. WHEN a split is executed, THE Split_Service SHALL ensure the total quantity of each item across Original_Order and New_Order equals the quantity in the Original_Order before the split.
2. WHEN a split is executed, THE Split_Service SHALL ensure the sum of Subtotals of Original_Order and New_Order equals the Subtotal of the Original_Order before the split.
3. IF the Order_Repository fails to save either the updated Original_Order or the New_Order, THEN THE Split_Service SHALL return an error and leave the Original_Order unchanged.
4. THE Split_Service SHALL NOT deduct batch ingredients again during a split, as ingredients were already deducted when the Original_Order was created.

---

### Requirement 6: API endpoint tách bill

**User Story:** As a waiter using the POS app, I want a dedicated API endpoint to split a bill, so that the frontend can trigger the operation reliably.

#### Acceptance Criteria

1. THE Split_Service SHALL expose the split operation via POST /waiter/orders/{id}/split.
2. WHEN a valid split request is submitted, THE Split_Service SHALL return a response containing both the updated Original_Order and the New_Order.
3. WHEN an invalid split request is submitted, THE Split_Service SHALL return an appropriate HTTP error code and a descriptive error message in Vietnamese.
4. THE Split_Service SHALL require waiter authentication to access the split endpoint.

---

### Requirement 7: Ghi nhận lịch sử tách bill

**User Story:** As a manager, I want to know when and how a bill was split, so that I can audit order history.

#### Acceptance Criteria

1. WHEN a split is executed, THE Split_Service SHALL set the Note field of the New_Order to include a reference to the Original_Order number in the format "Tách từ order #[OrderNumber]".
2. WHEN a split is executed, THE Split_Service SHALL append a note to the Original_Order indicating a split occurred, in the format "Đã tách ra order #[NewOrderNumber]".
