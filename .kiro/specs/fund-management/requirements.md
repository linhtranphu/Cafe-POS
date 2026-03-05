# Fund Management - Requirements

## Tổng quan
Manager cần một view để quản lý quỹ tiền của quán, bao gồm:
- Xem tổng quan số dư quỹ (tiền mặt và chuyển khoản)
- Xem lịch sử giao dịch ra/vào quỹ
- Thêm/rút tiền từ quỹ (deposit/withdrawal)
- Xem chi tiết các giao dịch

## User Stories

### US1: Xem tổng quan quỹ
**Là** manager  
**Tôi muốn** xem tổng quan số dư quỹ hiện tại  
**Để** biết quán đang có bao nhiêu tiền

**Acceptance Criteria:**
- Hiển thị số dư tiền mặt hiện tại
- Hiển thị số dư chuyển khoản hiện tại
- Hiển thị tổng số dư
- Hiển thị số dư đầu ngày (opening balance)
- Hiển thị tổng thu trong ngày (từ orders + handovers)
- Hiển thị tổng chi trong ngày (withdrawals)

### US2: Xem lịch sử giao dịch
**Là** manager  
**Tôi muốn** xem lịch sử tất cả giao dịch ra/vào quỹ  
**Để** theo dõi dòng tiền

**Acceptance Criteria:**
- Hiển thị danh sách giao dịch theo thời gian (mới nhất trước)
- Mỗi giao dịch hiển thị:
  - Loại giao dịch (deposit, withdrawal, handover, fund_handover)
  - Số tiền (cash/transfer)
  - Người thực hiện
  - Thời gian
  - Ghi chú/lý do
- Có thể filter theo:
  - Loại giao dịch
  - Ngày (hôm nay, hôm qua, tuần này, tháng này, custom range)
  - Loại tiền (cash, transfer, all)

### US3: Thêm tiền vào quỹ (Deposit)
**Là** manager  
**Tôi muốn** thêm tiền vào quỹ  
**Để** bổ sung vốn hoặc ghi nhận tiền từ nguồn khác

**Acceptance Criteria:**
- Form nhập:
  - Loại tiền (cash/transfer/both)
  - Số tiền cash (nếu có)
  - Số tiền transfer (nếu có)
  - Lý do/ghi chú (required, min 10 chars)
- Validate số tiền > 0
- Ghi nhận người thực hiện và thời gian
- Cập nhật số dư quỹ ngay lập tức

### US4: Rút tiền từ quỹ (Withdrawal)
**Là** manager  
**Tôi muốn** rút tiền từ quỹ  
**Để** chi tiêu hoặc chuyển tiền ra ngoài

**Acceptance Criteria:**
- Form nhập:
  - Loại tiền (cash/transfer/both)
  - Số tiền cash (nếu có)
  - Số tiền transfer (nếu có)
  - Lý do/ghi chú (required, min 10 chars)
- Validate số tiền > 0
- Validate số tiền rút <= số dư hiện tại
- Ghi nhận người thực hiện và thời gian
- Cập nhật số dư quỹ ngay lập tức

### US5: Xem chi tiết giao dịch
**Là** manager  
**Tôi muốn** xem chi tiết một giao dịch  
**Để** hiểu rõ hơn về giao dịch đó

**Acceptance Criteria:**
- Click vào giao dịch để xem chi tiết
- Hiển thị đầy đủ thông tin:
  - ID giao dịch
  - Loại giao dịch
  - Số tiền (cash/transfer/total)
  - Người thực hiện (tên + role)
  - Thời gian
  - Ghi chú/lý do
  - Số dư trước/sau giao dịch (nếu có)
  - Metadata khác (shift_id, order_id, etc.)

## Data Sources

### Giao dịch vào quỹ (Inflow):
1. **Cash Handovers** (từ waiter → cashier)
   - Collection: `cash_handovers`
   - Fields: cash_amount, transfer_amount, from_shift_id, to_shift_id
   
2. **Fund Handovers** (từ cashier → fund)
   - Collection: `fund_handovers`
   - Fields: cash_amount, transfer_amount, cashier_shift_id
   
3. **Deposits** (manager thêm tiền)
   - Collection: `fund_transactions` (new)
   - Type: "deposit"

### Giao dịch ra quỹ (Outflow):
1. **Withdrawals** (manager rút tiền)
   - Collection: `fund_transactions` (new)
   - Type: "withdrawal"

2. **Starting Float** (tiền đầu ca cho cashier)
   - Collection: `cashier_shifts`
   - Field: starting_float

## Technical Requirements

### Backend:
1. Tạo domain model `FundTransaction`
2. Tạo repository `FundTransactionRepository`
3. Tạo service `FundService` với methods:
   - `GetCurrentBalance()` - tính số dư hiện tại
   - `GetTransactionHistory(filters)` - lấy lịch sử
   - `Deposit(amount, type, reason)` - thêm tiền
   - `Withdraw(amount, type, reason)` - rút tiền
   - `GetTransactionDetail(id)` - chi tiết giao dịch
4. Tạo HTTP handlers cho manager role
5. API endpoints:
   - `GET /api/manager/fund/balance` - số dư hiện tại
   - `GET /api/manager/fund/transactions` - lịch sử (with filters)
   - `POST /api/manager/fund/deposit` - thêm tiền
   - `POST /api/manager/fund/withdraw` - rút tiền
   - `GET /api/manager/fund/transactions/:id` - chi tiết

### Frontend:
1. Tạo view `FundManagementView.vue`
2. Tạo service `fundService.js`
3. Tạo store `fundStore.js` (optional)
4. Components:
   - `FundBalanceCard.vue` - hiển thị số dư
   - `FundTransactionList.vue` - danh sách giao dịch
   - `DepositModal.vue` - form thêm tiền
   - `WithdrawalModal.vue` - form rút tiền
   - `TransactionDetailModal.vue` - chi tiết giao dịch

## UI/UX Requirements

### Mobile-First Design:
- Card-based layout
- Pull-to-refresh
- Touch-friendly buttons (min 44px height)
- Bottom navigation
- Responsive grid

### Color Coding:
- 💵 Green: Tiền mặt (cash)
- 💳 Blue: Chuyển khoản (transfer)
- 📥 Green: Giao dịch vào (inflow/deposit)
- 📤 Red: Giao dịch ra (outflow/withdrawal)
- 🔄 Orange: Handover transactions

### Navigation:
- Thêm menu item "💰 Quỹ tiền" trong manager dashboard
- Route: `/manager/fund`

## Security & Permissions
- Chỉ manager role mới được truy cập
- Tất cả operations phải ghi audit log
- Validate permissions ở cả frontend và backend

## Future Enhancements (Out of Scope)
- Export báo cáo Excel/PDF
- Biểu đồ thống kê dòng tiền
- Alerts khi số dư thấp
- Multi-currency support
- Reconciliation với bank statements
