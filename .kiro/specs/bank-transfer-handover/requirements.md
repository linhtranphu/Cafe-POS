# Requirements Document - Bank Transfer Handover

## Introduction

Tính năng này mở rộng quy trình bàn giao tiền hiện tại để hỗ trợ cả tiền chuyển khoản (bank transfer) bên cạnh tiền mặt. Waiter có thể bàn giao cả hai loại tiền trong ca làm việc, và Cashier sẽ đối soát với tài khoản ngân hàng bên ngoài hệ thống trước khi xác nhận.

## Glossary

- **Waiter**: Nhân viên phục vụ, người thu tiền từ khách hàng và bàn giao cho Cashier
- **Cashier**: Thu ngân, người nhận và xác nhận tiền bàn giao từ Waiter
- **Bank_Transfer**: Phương thức thanh toán chuyển khoản ngân hàng (bao gồm cả QR code)
- **Cash**: Tiền mặt
- **Handover**: Quy trình bàn giao tiền từ Waiter sang Cashier
- **Declared_Amount**: Số tiền Waiter khai báo khi bàn giao
- **Actual_Amount**: Số tiền Cashier xác nhận sau khi đối soát
- **Discrepancy**: Chênh lệch giữa số tiền khai báo và số tiền thực tế
- **Shift**: Ca làm việc
- **External_Bank_Account**: Tài khoản ngân hàng bên ngoài hệ thống dùng để đối soát

## Requirements

### Requirement 1: Waiter Bank Transfer Handover

**User Story:** Là một Waiter, tôi muốn bàn giao tiền chuyển khoản tương tự như tiền mặt, để tôi có thể hoàn tất quy trình bàn giao cho tất cả các loại thanh toán trong ca.

#### Acceptance Criteria

1. WHEN a Waiter creates a handover request, THE System SHALL allow selection of handover type (cash, bank transfer, or both)
2. WHEN a Waiter selects bank transfer handover, THE System SHALL display the total bank transfer amount collected in the shift
3. WHEN a Waiter declares bank transfer amount, THE System SHALL validate that the declared amount does not exceed the total bank transfer revenue in the shift
4. WHEN a Waiter submits a handover request with bank transfer, THE System SHALL create a handover record with separate cash and bank transfer amounts
5. WHEN a Waiter creates an end-shift handover, THE System SHALL automatically include both remaining cash and total bank transfer amounts

### Requirement 2: Separate Tracking of Payment Methods

**User Story:** Là hệ thống, tôi cần phân biệt và theo dõi riêng biệt tiền mặt và tiền chuyển khoản, để đảm bảo tính chính xác trong quản lý tài chính.

#### Acceptance Criteria

1. THE System SHALL track cash revenue separately from bank transfer revenue for each shift
2. THE System SHALL calculate total bank transfer revenue by summing all TRANSFER and QR payment methods
3. WHEN displaying shift summary, THE System SHALL show cash amount and bank transfer amount as separate line items
4. WHEN calculating handover amounts, THE System SHALL maintain separate totals for cash and bank transfer
5. THE System SHALL store cash_declared_amount and transfer_declared_amount as separate fields in handover records

### Requirement 3: Cashier Bank Transfer Reconciliation

**User Story:** Là một Cashier, tôi muốn đối soát tiền chuyển khoản với tài khoản ngân hàng bên ngoài, để tôi có thể xác nhận chính xác số tiền thực tế nhận được.

#### Acceptance Criteria

1. WHEN a Cashier views a pending handover with bank transfer, THE System SHALL display both declared cash amount and declared bank transfer amount
2. WHEN a Cashier confirms a handover, THE System SHALL require input of actual amounts for both cash and bank transfer separately
3. WHEN a Cashier enters actual bank transfer amount, THE System SHALL calculate discrepancy between declared and actual bank transfer amounts
4. THE System SHALL allow Cashier to add notes explaining the bank transfer reconciliation process
5. WHEN bank transfer discrepancy exceeds threshold, THE System SHALL flag the handover for manager approval

### Requirement 4: Dual Confirmation Process

**User Story:** Là một Cashier, tôi muốn xác nhận cả tiền mặt và tiền chuyển khoản trong một giao dịch bàn giao, để quy trình xác nhận được đơn giản và nhất quán.

#### Acceptance Criteria

1. WHEN a Cashier confirms a handover, THE System SHALL require confirmation of both cash and bank transfer amounts in a single transaction
2. WHEN a Cashier rejects a handover, THE System SHALL reject both cash and bank transfer components together
3. THE System SHALL calculate total discrepancy as the sum of cash discrepancy and bank transfer discrepancy
4. WHEN updating shift balances, THE System SHALL update both cash and bank transfer balances atomically
5. IF either cash or bank transfer confirmation fails, THEN THE System SHALL rollback the entire handover transaction

### Requirement 5: Display and Reporting

**User Story:** Là một Cashier, tôi muốn xem thông tin chi tiết về tiền mặt và tiền chuyển khoản trong giao diện, để tôi có thể dễ dàng theo dõi và quản lý.

#### Acceptance Criteria

1. WHEN displaying shift summary, THE System SHALL show cash revenue and bank transfer revenue as separate colored badges
2. WHEN displaying handover history, THE System SHALL show both cash and bank transfer amounts for each handover
3. WHEN displaying pending handovers, THE System SHALL clearly distinguish between cash-only, transfer-only, and combined handovers
4. THE System SHALL use consistent color coding: green for cash, blue for bank transfer
5. WHEN displaying discrepancies, THE System SHALL show separate discrepancy amounts for cash and bank transfer

### Requirement 6: Data Integrity and Validation

**User Story:** Là hệ thống, tôi cần đảm bảo tính toàn vẹn của dữ liệu khi xử lý bàn giao tiền chuyển khoản, để tránh mất mát hoặc sai sót trong quản lý tài chính.

#### Acceptance Criteria

1. THE System SHALL validate that total declared amount (cash + bank transfer) does not exceed total shift revenue
2. WHEN a Waiter has pending handover, THE System SHALL prevent creation of new handover for the same shift
3. THE System SHALL ensure that bank transfer amounts are non-negative
4. WHEN calculating remaining balances, THE System SHALL maintain separate remaining_cash and remaining_transfer fields
5. THE System SHALL log all handover transactions with timestamps and user information for audit purposes

### Requirement 7: Backward Compatibility

**User Story:** Là hệ thống, tôi cần duy trì khả năng tương thích ngược với quy trình bàn giao tiền mặt hiện tại, để không ảnh hưởng đến các ca làm việc đang diễn ra.

#### Acceptance Criteria

1. WHEN a Waiter creates a cash-only handover, THE System SHALL process it using the existing cash handover logic
2. THE System SHALL treat handovers without transfer_declared_amount as cash-only handovers
3. WHEN displaying old handover records, THE System SHALL show them correctly even if they lack bank transfer fields
4. THE System SHALL allow Cashier to confirm cash-only handovers without requiring bank transfer input
5. WHEN migrating existing data, THE System SHALL set transfer amounts to zero for old handover records

### Requirement 8: Manager Approval for Large Discrepancies

**User Story:** Là một Manager, tôi muốn được thông báo và phê duyệt các bàn giao có chênh lệch lớn về tiền chuyển khoản, để đảm bảo kiểm soát tài chính chặt chẽ.

#### Acceptance Criteria

1. WHEN total discrepancy (cash + bank transfer) exceeds 100,000 VND, THE System SHALL require manager approval
2. WHEN bank transfer discrepancy alone exceeds 50,000 VND, THE System SHALL flag for manager review
3. THE System SHALL send notification to manager when handover requires approval
4. WHEN a Manager approves a handover, THE System SHALL update shift balances and complete the handover
5. WHEN a Manager rejects a handover, THE System SHALL return it to pending status and notify the Cashier
