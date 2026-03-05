# Requirements Document

## Introduction

Tính năng Fund-Expense Integration tích hợp hệ thống quản lý quỹ (fund) với các module chi tiêu hiện có (expense, ingredient restock, facility maintenance) trong hệ thống Cafe POS. Tính năng này cho phép Manager chi tiền trực tiếp từ quỹ khi tạo các giao dịch chi tiêu, đồng thời tự động ghi nhận fund transaction và cập nhật số dư quỹ.

## Implementation Scope

### Phase 1 (REQUIRED)
- Requirement 1: Tạo Expense từ Quỹ
- Requirement 2: Mua Ingredient từ Quỹ
- Requirement 4: Xem Lịch sử Chi tiêu từ Quỹ
- Requirement 6: Kiểm tra Số dư Quỹ
- Requirement 7: Audit Trail cho Fund Transactions
- Requirement 8: Database Schema Changes
- Requirement 9: Transaction Atomicity
- Requirement 10: Payment Method Consistency
- Requirement 11: Prevent Double Spending

### Future Phases (OPTIONAL)
- Requirement 3: Chi phí Sửa chữa Facility từ Quỹ (optional feature)
- Requirement 5: Báo cáo Sử dụng Quỹ theo Mục đích (optional feature)

## Glossary

- **Fund_Manager**: Module quản lý quỹ tiền mặt của cửa hàng
- **Expense_Module**: Module ghi nhận các khoản chi phí vận hành
- **Ingredient_Module**: Module quản lý nguyên liệu và restock
- **Facility_Module**: Module quản lý thiết bị và bảo trì
- **Fund_Transaction**: Giao dịch nạp/rút tiền từ quỹ
- **Withdrawal**: Giao dịch rút tiền từ quỹ
- **Source_Record**: Record gốc tạo ra fund transaction (expense/ingredient/facility)
- **Manager**: Người dùng có quyền quản lý quỹ và chi tiêu
- **Fund_Balance**: Số dư hiện tại của quỹ
- **Audit_Trail**: Lịch sử đầy đủ các giao dịch để kiểm soát

## Requirements

### Requirement 1: Tạo Expense từ Quỹ

**User Story:** Là Manager, tôi muốn đánh dấu expense được chi từ quỹ khi tạo expense mới, để hệ thống tự động trừ tiền từ quỹ và ghi nhận giao dịch.

#### Acceptance Criteria

1. THE Expense_Module SHALL provide a "paid from fund" option in the expense creation form
2. WHEN Manager selects "paid from fund" option, THE Expense_Module SHALL validate that Fund_Balance is sufficient for the expense amount
3. WHEN Manager creates an expense with "paid from fund" selected AND Fund_Balance is sufficient, THE Fund_Manager SHALL create a Withdrawal transaction with amount equal to expense amount
4. WHEN a Withdrawal transaction is created for an expense, THE Fund_Manager SHALL update Fund_Balance by subtracting the withdrawal amount
5. WHEN an expense is created from fund, THE Expense_Module SHALL store the fund_transaction_id reference in the expense record
6. WHEN an expense is created from fund, THE Fund_Manager SHALL store source_type as "expense" and source_id as expense record identifier in the Fund_Transaction

### Requirement 2: Mua Ingredient từ Quỹ

**User Story:** Là Manager, tôi muốn nhập thêm nguyên liệu và chi tiền từ quỹ, để theo dõi chi phí mua nguyên liệu từ quỹ.

#### Acceptance Criteria

1. THE Ingredient_Module SHALL provide a "paid from fund" option in the ingredient restock form
2. WHEN Manager selects "paid from fund" option for restock, THE Ingredient_Module SHALL validate that Fund_Balance is sufficient for the restock cost
3. WHEN Manager restocks ingredient with "paid from fund" selected AND Fund_Balance is sufficient, THE Expense_Module SHALL create an expense record with category "ingredient purchase"
4. WHEN an ingredient restock is paid from fund, THE Fund_Manager SHALL create a Withdrawal transaction with amount equal to restock cost
5. WHEN an ingredient restock is paid from fund, THE Ingredient_Module SHALL update ingredient stock quantity
6. WHEN an ingredient restock creates fund withdrawal, THE Fund_Manager SHALL store source_type as "ingredient" and source_id as restock record identifier in the Fund_Transaction
7. WHEN an ingredient restock is paid from fund, THE Expense_Module SHALL link the expense record to the Fund_Transaction via fund_transaction_id

### Requirement 3: Chi phí Sửa chữa Facility từ Quỹ (OPTIONAL)

**User Story:** Là Manager, tôi muốn ghi nhận chi phí sửa chữa thiết bị từ quỹ, để theo dõi chi phí bảo trì từ quỹ.

**Note:** This is an optional feature for future implementation. All acceptance criteria use WHERE clause to indicate conditional implementation.

#### Acceptance Criteria

1. WHERE facility maintenance feature is enabled, THE Facility_Module SHALL provide a "paid from fund" option in the facility issue report form
2. WHERE facility maintenance is paid from fund, WHEN Manager reports a facility issue with repair cost, THE Facility_Module SHALL validate that Fund_Balance is sufficient for the repair cost
3. WHERE facility maintenance is paid from fund, WHEN Manager reports facility issue with cost AND Fund_Balance is sufficient, THE Facility_Module SHALL create a facility history record
4. WHERE facility maintenance is paid from fund, WHEN a facility issue with cost is reported, THE Expense_Module SHALL create an expense record with category "facility maintenance"
5. WHERE facility maintenance is paid from fund, WHEN a facility issue with cost is reported, THE Fund_Manager SHALL create a Withdrawal transaction with amount equal to repair cost
6. WHERE facility maintenance is paid from fund, WHEN facility issue creates fund withdrawal, THE Fund_Manager SHALL store source_type as "facility" and source_id as facility history record identifier in the Fund_Transaction
7. WHERE facility maintenance is paid from fund, WHEN facility issue is reported, THE Expense_Module SHALL link the expense record to the Fund_Transaction via fund_transaction_id

### Requirement 4: Xem Lịch sử Chi tiêu từ Quỹ

**User Story:** Là Manager, tôi muốn xem expense nào được chi từ quỹ, để audit và kiểm soát chi tiêu.

#### Acceptance Criteria

1. WHEN Manager views expense list, THE Expense_Module SHALL display a visual indicator for expenses that have paid_from_fund flag set to true
2. WHEN Manager views an expense detail that was paid from fund, THE Expense_Module SHALL display a link to the associated Fund_Transaction
3. THE Expense_Module SHALL provide a filter option to show only expenses where paid_from_fund is true
4. WHEN Manager applies "paid from fund" filter, THE Expense_Module SHALL display only expenses with fund_transaction_id populated

### Requirement 5: Báo cáo Sử dụng Quỹ theo Mục đích (OPTIONAL)

**User Story:** Là Manager, tôi muốn xem báo cáo quỹ được dùng cho mục đích gì, để phân tích và lập kế hoạch ngân sách.

**Note:** This is an optional feature for future implementation. All acceptance criteria use WHERE clause to indicate conditional implementation.

#### Acceptance Criteria

1. WHERE fund usage reporting is enabled, THE Fund_Manager SHALL provide a report showing fund withdrawals grouped by source_type
2. WHERE fund usage reporting is enabled, WHEN Manager views fund usage report, THE Fund_Manager SHALL display total withdrawal amount for each source_type
3. WHERE fund usage reporting is enabled, WHEN Manager views fund usage report, THE Fund_Manager SHALL display count of transactions for each source_type
4. WHERE fund usage reporting is enabled, THE Fund_Manager SHALL provide date range filter for fund usage report
5. WHERE fund usage reporting is enabled, WHEN Manager applies date range filter, THE Fund_Manager SHALL display fund usage data only for Fund_Transactions within the specified date range

### Requirement 6: Kiểm tra Số dư Quỹ

**User Story:** Là hệ thống, tôi cần kiểm tra số dư quỹ đủ trước khi tạo withdrawal, để tránh chi vượt quỹ.

#### Acceptance Criteria

1. WHEN a withdrawal request is initiated, THE Fund_Manager SHALL retrieve current Fund_Balance before processing
2. WHEN a withdrawal request amount exceeds Fund_Balance, THE Fund_Manager SHALL reject the withdrawal request
3. WHEN a withdrawal request is rejected due to insufficient balance, THE Fund_Manager SHALL return an error message indicating current Fund_Balance and requested amount
4. WHEN Manager attempts to create expense with "paid from fund" AND Fund_Balance is insufficient, THE Expense_Module SHALL display the error message and prevent expense creation
5. WHEN Manager attempts to restock ingredient with "paid from fund" AND Fund_Balance is insufficient, THE Ingredient_Module SHALL display the error message and prevent restock operation
6. WHERE facility maintenance is enabled, WHEN Manager attempts to report facility issue with "paid from fund" AND Fund_Balance is insufficient, THE Facility_Module SHALL display the error message and prevent issue creation

### Requirement 7: Audit Trail cho Fund Transactions

**User Story:** Là Manager, tôi muốn xem ai đã chi tiền từ quỹ, khi nào, và cho mục đích gì, để kiểm soát và audit.

#### Acceptance Criteria

1. WHEN a Fund_Transaction is created, THE Fund_Manager SHALL record the user identifier of the Manager who initiated the transaction
2. WHEN a Fund_Transaction is created, THE Fund_Manager SHALL record the user role of the Manager who initiated the transaction
3. WHEN a Fund_Transaction is created, THE Fund_Manager SHALL record the timestamp with timezone information
4. WHEN a Fund_Transaction is created, THE Fund_Manager SHALL record the reason or description for the transaction
5. WHEN a Fund_Transaction is created, THE Fund_Manager SHALL record Fund_Balance before the transaction
6. WHEN a Fund_Transaction is created, THE Fund_Manager SHALL record Fund_Balance after the transaction
7. WHEN Manager views Fund_Transaction detail, THE Fund_Manager SHALL display a link to the Source_Record based on source_type and source_id
8. WHEN Manager clicks on Source_Record link, THE Fund_Manager SHALL navigate to the detail view of the corresponding expense, ingredient restock, or facility issue record

### Requirement 8: Database Schema Changes

**User Story:** Là hệ thống, tôi cần lưu trữ thông tin liên kết giữa fund transactions và source records, để duy trì tính toàn vẹn dữ liệu.

#### Acceptance Criteria

1. THE Expense_Module SHALL store fund_transaction_id field in expense records to reference associated Fund_Transaction
2. THE Expense_Module SHALL store paid_from_fund boolean field in expense records to indicate payment source
3. THE Fund_Manager SHALL store source_type field in Fund_Transaction records with allowed values "expense", "ingredient", or "facility"
4. THE Fund_Manager SHALL store source_id field in Fund_Transaction records to reference the originating record identifier
5. WHEN a Fund_Transaction has source_type populated, THE Fund_Manager SHALL ensure source_id is also populated
6. WHEN an expense has paid_from_fund set to true, THE Expense_Module SHALL ensure fund_transaction_id is populated

### Requirement 9: Transaction Atomicity

**User Story:** Là hệ thống, tôi cần đảm bảo tính nguyên tử của các thao tác liên quan đến fund và expense, để tránh trạng thái dữ liệu không nhất quán.

#### Acceptance Criteria

1. WHEN creating expense from fund, IF Fund_Transaction creation fails, THEN THE Expense_Module SHALL rollback the expense creation
2. WHEN creating expense from fund, IF expense creation fails after Fund_Transaction is created, THEN THE Fund_Manager SHALL rollback the Fund_Transaction and restore Fund_Balance
3. WHEN restocking ingredient from fund, IF any operation fails in the sequence, THEN THE Ingredient_Module SHALL rollback all changes including expense, Fund_Transaction, and stock update
4. WHERE facility maintenance is enabled, WHEN reporting facility issue from fund, IF any operation fails in the sequence, THEN THE Facility_Module SHALL rollback all changes including facility history, expense, and Fund_Transaction
5. WHEN a rollback occurs, THE Fund_Manager SHALL restore Fund_Balance to the value before the failed transaction

### Requirement 10: Payment Method Consistency

**User Story:** Là hệ thống, tôi cần đảm bảo tính nhất quán của phương thức thanh toán, để tránh xung đột dữ liệu.

#### Acceptance Criteria

1. WHEN Manager selects "paid from fund" for expense, THE Expense_Module SHALL set payment method to "fund" or equivalent fund payment indicator
2. WHEN an expense is paid from fund, THE Expense_Module SHALL prevent Manager from selecting other payment methods simultaneously
3. WHEN Manager creates Fund_Transaction with source_type "expense", THE Fund_Manager SHALL verify that the referenced expense has payment method set to fund payment indicator
4. IF payment method inconsistency is detected, THEN THE Fund_Manager SHALL reject the Fund_Transaction creation and return validation error

### Requirement 11: Prevent Double Spending

**User Story:** Là hệ thống, tôi cần ngăn chặn việc chi tiêu trùng lặp từ quỹ, để đảm bảo tính chính xác của số dư quỹ.

#### Acceptance Criteria

1. WHEN an expense already has fund_transaction_id populated, THE Expense_Module SHALL prevent Manager from modifying the "paid from fund" flag
2. WHEN a Fund_Transaction already exists with specific source_type and source_id, THE Fund_Manager SHALL reject creation of duplicate Fund_Transaction with same source_type and source_id
3. WHEN Manager attempts to edit an expense that was paid from fund, THE Expense_Module SHALL prevent changes to the expense amount
4. IF Manager needs to correct an expense paid from fund, THE Expense_Module SHALL require Manager to void the original expense and create a new expense record
