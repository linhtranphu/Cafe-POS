# Requirements Document - Cashier Fund Handover

## Introduction

Tính năng này cho phép Cashier xem tổng số tiền đang quản lý (nhận từ các waiter) và handover lại số tiền này vào quỹ khi đóng ca. Cashier chịu trách nhiệm trên số tiền này từ khi nhận bàn giao cho đến khi handover lại.

## Glossary

- **Cashier**: Thu ngân, người nhận tiền từ waiter và quản lý quỹ
- **Fund**: Quỹ tiền của cửa hàng
- **Received_Cash**: Tổng tiền mặt cashier đã nhận từ các waiter
- **Received_Transfer**: Tổng tiền chuyển khoản cashier đã nhận từ các waiter
- **Total_Managed_Funds**: Tổng số tiền cashier đang quản lý = Received_Cash + Received_Transfer
- **Fund_Handover**: Quy trình bàn giao tiền từ cashier về quỹ khi đóng ca
- **Starting_Float**: Tiền đầu ca của cashier
- **Actual_Cash_Count**: Số tiền mặt thực tế đếm được khi đóng ca
- **Variance**: Chênh lệch giữa số tiền lý thuyết và thực tế

## Requirements

### Requirement 1: Display Managed Funds in Dashboard

**User Story:** Là một Cashier, tôi muốn xem rõ tổng số tiền tôi đang quản lý trong dashboard, để tôi biết mình chịu trách nhiệm trên số tiền nào.

#### Acceptance Criteria

1. WHEN a Cashier views the dashboard with an open shift, THE System SHALL display a "Managed Funds" section showing:
   - Received cash amount from waiter handovers
   - Received transfer amount from waiter handovers
   - Total managed funds (cash + transfer)
2. THE System SHALL display this information prominently with clear visual hierarchy
3. THE System SHALL use color coding: green for cash, blue for transfer, orange/yellow for total
4. THE System SHALL include a warning message indicating cashier responsibility
5. WHEN there is no open cashier shift, THE System SHALL NOT display the managed funds section

### Requirement 2: Display Managed Funds in Closure View

**User Story:** Là một Cashier, khi tôi đóng ca, tôi muốn xem rõ số tiền tôi cần handover lại, để tôi có thể chuẩn bị và đối soát chính xác.

#### Acceptance Criteria

1. WHEN a Cashier initiates shift closure, THE System SHALL display managed funds summary including:
   - Starting float (tiền đầu ca)
   - Received cash from handovers
   - Received transfer from handovers
   - Total cash to account for (starting float + received cash)
2. THE System SHALL clearly separate cash and transfer amounts
3. THE System SHALL calculate expected cash = starting_float + received_cash
4. THE System SHALL display this information before requesting actual cash count
5. THE System SHALL use visual indicators to highlight the amounts cashier is responsible for

### Requirement 3: Fund Handover Process

**User Story:** Là một Cashier, khi tôi đóng ca, tôi muốn handover toàn bộ tiền mặt và ghi nhận tiền chuyển khoản về quỹ, để hoàn tất trách nhiệm của mình.

#### Acceptance Criteria

1. WHEN a Cashier closes shift, THE System SHALL require handover of all managed funds
2. THE System SHALL prompt cashier to count and enter actual cash amount
3. THE System SHALL calculate variance between expected and actual cash
4. THE System SHALL require documentation if variance exceeds threshold
5. THE System SHALL record transfer amount as "handed over to fund" (no physical handover needed)
6. THE System SHALL create a fund handover record with:
   - Cashier information
   - Cash amount handed over
   - Transfer amount recorded
   - Variance (if any)
   - Timestamp
   - Receiver information (nullable for future extension)

### Requirement 4: Variance Handling

**User Story:** Là một Cashier, nếu có chênh lệch giữa số tiền lý thuyết và thực tế, tôi muốn ghi nhận và giải thích chênh lệch này, để đảm bảo minh bạch.

#### Acceptance Criteria

1. WHEN actual cash differs from expected cash, THE System SHALL calculate variance
2. THE System SHALL display variance with clear indication (shortage in red, overage in green)
3. IF variance is non-zero, THE System SHALL require:
   - Reason selection (counting error, theft, customer dispute, etc.)
   - Detailed notes (minimum 10 characters)
4. THE System SHALL store variance information in cashier shift record
5. THE System SHALL include variance in fund handover record

### Requirement 5: Fund Handover Record

**User Story:** Là hệ thống, tôi cần ghi nhận chi tiết việc bàn giao tiền về quỹ, để có audit trail đầy đủ.

#### Acceptance Criteria

1. THE System SHALL create a fund_handover record when cashier closes shift
2. THE record SHALL include:
   - cashier_shift_id
   - cashier_id
   - cashier_name
   - cash_amount (actual cash handed over)
   - transfer_amount (recorded transfer amount)
   - total_amount (cash + transfer)
   - variance_amount
   - variance_reason (if applicable)
   - variance_notes (if applicable)
   - receiver_id (nullable, for future use)
   - receiver_name (nullable, for future use)
   - handover_at (timestamp)
   - created_at
   - updated_at
3. THE System SHALL link this record to the cashier shift
4. THE System SHALL ensure atomicity - fund handover and shift closure happen together

### Requirement 6: Receiver Extension Point

**User Story:** Là hệ thống, tôi cần thiết kế sẵn khả năng mở rộng để sau này có thể chỉ định người nhận (manager) khi handover về quỹ.

#### Acceptance Criteria

1. THE System SHALL include receiver_id and receiver_name fields in fund_handover record (nullable)
2. WHEN receiver is not specified, THE System SHALL set receiver fields to null
3. THE System SHALL allow these fields to be populated in future versions
4. THE System SHALL design API to accept optional receiver_id parameter
5. THE System SHALL maintain backward compatibility when receiver feature is added

### Requirement 7: Cashier Shift Closure Integration

**User Story:** Là một Cashier, việc handover về quỹ phải là một phần không thể tách rời của quy trình đóng ca, để đảm bảo tính toàn vẹn.

#### Acceptance Criteria

1. WHEN a Cashier closes shift, THE System SHALL automatically trigger fund handover process
2. THE System SHALL NOT allow shift closure without completing fund handover
3. THE System SHALL use transaction to ensure atomicity:
   - Create fund handover record
   - Update cashier shift status to CLOSED
   - Record variance (if any)
   - Set end time
4. IF any step fails, THE System SHALL rollback entire transaction
5. THE System SHALL provide clear error messages if closure fails

### Requirement 8: Display Handover History

**User Story:** Là một Cashier hoặc Manager, tôi muốn xem lịch sử các lần handover về quỹ, để theo dõi và kiểm tra.

#### Acceptance Criteria

1. THE System SHALL provide an endpoint to query fund handover history
2. THE System SHALL allow filtering by:
   - Date range
   - Cashier
   - Variance status (with/without variance)
3. THE System SHALL display handover records with:
   - Cashier name
   - Date and time
   - Cash amount
   - Transfer amount
   - Variance (if any)
   - Status
4. THE System SHALL sort records by date descending (newest first)
5. THE System SHALL support pagination for large result sets

### Requirement 9: Audit and Compliance

**User Story:** Là hệ thống, tôi cần đảm bảo tất cả các giao dịch handover về quỹ được ghi nhận đầy đủ cho mục đích kiểm toán.

#### Acceptance Criteria

1. THE System SHALL log all fund handover operations
2. THE System SHALL include in audit log:
   - User ID and name
   - Action (fund_handover_created)
   - Timestamp
   - Device ID
   - IP address
   - All amounts (cash, transfer, variance)
3. THE System SHALL make audit logs immutable
4. THE System SHALL retain audit logs for compliance period
5. THE System SHALL provide audit log query capability for authorized users

### Requirement 10: Error Handling and Validation

**User Story:** Là hệ thống, tôi cần validate dữ liệu và xử lý lỗi một cách graceful để đảm bảo tính toàn vẹn.

#### Acceptance Criteria

1. THE System SHALL validate that cashier shift is in CLOSURE_INITIATED status before allowing fund handover
2. THE System SHALL validate that actual_cash_amount is non-negative
3. THE System SHALL validate that variance documentation is provided when variance is non-zero
4. THE System SHALL handle database errors gracefully with rollback
5. THE System SHALL provide clear error messages to users
6. THE System SHALL log all errors for debugging
7. THE System SHALL prevent duplicate fund handover for the same shift

## Non-Functional Requirements

### Performance
- Fund handover operation should complete within 2 seconds
- Dashboard should load managed funds within 1 second
- Audit log queries should return within 3 seconds

### Security
- Only authenticated cashiers can view their managed funds
- Only the cashier who owns the shift can perform fund handover
- Managers can view all fund handover records
- Audit logs are immutable and tamper-proof

### Reliability
- Use database transactions to ensure atomicity
- Implement retry logic for transient failures
- Maintain data consistency across all operations

### Usability
- Clear visual hierarchy for managed funds display
- Intuitive flow for fund handover process
- Helpful error messages and validation feedback
- Mobile-friendly interface

## Success Metrics

1. Cashiers can clearly see their managed funds in dashboard
2. 100% of shift closures include fund handover records
3. Variance documentation rate > 95% when variance exists
4. Zero data inconsistencies in fund handover records
5. Audit trail is complete and queryable

## Out of Scope

- Automatic reconciliation with bank accounts
- Integration with external accounting systems
- Manager approval workflow for large variances (can be added later)
- Physical cash counting assistance (e.g., denomination breakdown)
- Multi-currency support

## Future Enhancements

- Specify receiver (manager) when handing over to fund
- Manager approval for large variances
- Denomination breakdown for cash counting
- Integration with accounting system
- Real-time notifications to manager when handover occurs
- Analytics dashboard for fund handover patterns
