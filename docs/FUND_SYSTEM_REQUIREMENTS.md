# Cafe POS — Fund & Accounting System: Requirements

**Version:** 2.0 (Double-Entry Rewrite)  
**Last Updated:** 2026-03-11  
**Status:** Core implemented. See FUND_SYSTEM_TASKS.md for remaining work.

---

## 1. Overview

Mọi luồng tiền trong quán đều được ghi nhận bằng **double-entry bookkeeping**. Mỗi sự kiện tạo một `JournalEntry` cân bằng (Σ DEBIT = Σ CREDIT). Journal là nguồn sự thật duy nhất để tính số dư quỹ.

### 1.1 Accounting Convention

| Side   | Ý nghĩa với real fund accounts |
|--------|-------------------------------|
| DEBIT  | Số dư quỹ **tăng** (tiền vào) |
| CREDIT | Số dư quỹ **giảm** (tiền ra) |

External counterpart accounts không có số dư thực — chỉ để cân bằng double-entry.

---

## 2. Fund Accounts

### 2.1 Real Fund Accounts (có số dư thực)

| Account        | Key           | Mô tả |
|----------------|---------------|-------|
| Quỹ vận hành  | `operating`   | Chi phí vận hành: điện, nước, thuê mặt bằng, linh tinh |
| Quỹ hàng hóa  | `inventory`   | Mua nguyên liệu, cơ sở vật chất |
| Quỹ lợi nhuận | `profit`      | Lợi nhuận tích lũy; chủ quán có thể rút |
| Ngăn kéo tiền | `cash_drawer` | Tiền mặt vật lý tại quầy thu ngân |
| Tiền phục vụ  | `waiter_float`| Tiền waiter đang giữ trong ca |

### 2.2 External Counterpart Accounts (virtual — chỉ để cân bằng)

| Account         | Key              | Dùng khi nào |
|-----------------|------------------|--------------|
| Chủ quán        | `owner`          | Nạp tiền / rút tiền của manager |
| Nhà cung cấp    | `supplier`       | Chi tiêu / mua nguyên liệu / mua tài sản |
| Khách hàng      | `customer`       | Doanh thu từ đơn hàng (tương lai) |
| Thiếu tiền      | `cash_shortage`  | Waiter bàn giao thiếu |
| Thừa tiền       | `cash_overage`   | Waiter bàn giao thừa |

---

## 3. Journal Event Types & Double-Entry Entries

### 3.1 Manager Operations

#### FR-FUND-01: Manager Deposit — Nạp tiền vào quỹ
```
DEBIT  {fund_type}   +amount   ← quỹ nhận tiền
CREDIT owner         +amount   ← chủ quán bỏ tiền vào
```
Rules: amount > 0; reason ≥ 10 chars; fund_type = any real fund

#### FR-FUND-02: Manager Withdrawal — Rút tiền từ quỹ
```
DEBIT  owner         +amount   ← chủ quán nhận lại tiền
CREDIT {fund_type}   +amount   ← quỹ xuất tiền
```
Rules: amount > 0; reason ≥ 10 chars; sufficient balance

#### FR-FUND-03: Fund Transfer — Chuyển giữa các quỹ
```
DEBIT  to_fund       +amount   ← quỹ đích nhận tiền
CREDIT from_fund     +amount   ← quỹ nguồn xuất tiền
```
Rules: from ≠ to; amount > 0; sufficient balance in from_fund

---

### 3.2 Cashier Shift Operations

#### FR-SHIFT-01: Cashier Shift Start — Đầu ca thu ngân
```
DEBIT  cash_drawer   +float    ← ngăn kéo nhận tiền lẻ đầu ca
CREDIT operating     +float    ← quỹ vận hành xuất tiền lẻ
```
Rules: float > 0; no existing open cashier shift

#### FR-SHIFT-02: Cashier Shift End — Đóng ca thu ngân
```
DEBIT  operating     +actual   ← quỹ vận hành nhận tiền về
CREDIT cash_drawer   +actual   ← ngăn kéo bàn giao tiền
```
Rules:
- Tất cả waiter shifts phải đóng trước
- expectedCash = startingFloat + receivedCash - distributedCash
- Nếu |actualCash - expectedCash| ≥ 0.01 → bắt buộc ghi lý do + ghi chú ≥ 10 ký tự
- Variance reasons: COUNTING_ERROR, UNRECORDED_SALE, THEFT, CHANGE_ERROR, SYSTEM_ERROR, OTHER

---

### 3.3 Waiter Shift Operations

#### FR-WAITER-01: Waiter Shift Start — Phát tiền đầu ca phục vụ
```
DEBIT  waiter_float  +amount   ← waiter nhận tiền lẻ
CREDIT cash_drawer   +amount   ← ngăn kéo xuất tiền cho waiter
```
Rules: amount > 0; cashier shift phải đang mở

#### FR-WAITER-02: Waiter Cash Handover — Waiter bàn giao tiền (không có chênh lệch)
```
DEBIT  cash_drawer   +amount   ← ngăn kéo nhận tiền
CREDIT waiter_float  +amount   ← waiter xuất tiền
```

#### FR-WAITER-03: Waiter Handover with Shortage — Waiter thiếu tiền
```
DEBIT  cash_drawer   +actual            ← ngăn kéo nhận số thực tế
DEBIT  cash_shortage +shortage_amount   ← ghi nhận khoản thiếu (audit)
CREDIT waiter_float  +expected          ← waiter xuất số kỳ vọng
```
where: shortage = expected - actual

#### FR-WAITER-04: Waiter Handover with Overage — Waiter thừa tiền
```
DEBIT  cash_drawer   +actual            ← ngăn kéo nhận số thực tế
CREDIT cash_overage  +overage_amount    ← ghi nhận khoản thừa (audit)
CREDIT waiter_float  +expected          ← waiter xuất số kỳ vọng
```
where: overage = actual - expected

---

### 3.4 Fund-Paid Expense Operations

#### FR-EX-01: Operating Expense — Chi tiêu vận hành từ quỹ
```
DEBIT  supplier    +amount   ← tiền ra ngoài (trả nhà cung cấp)
CREDIT operating   +amount   ← quỹ vận hành xuất tiền
```
Atomic: JournalEntry + Expense record trong cùng MongoDB transaction

#### FR-EX-02: Ingredient Restock — Mua nguyên liệu từ quỹ
```
DEBIT  supplier    +cost     ← trả nhà cung cấp
CREDIT inventory   +cost     ← quỹ hàng hóa xuất tiền
```
Atomic: JournalEntry + RestockRecord + Expense + ingredient quantity update

#### FR-EX-03: Facility Purchase — Mua tài sản từ quỹ
```
DEBIT  supplier    +cost     ← trả nhà cung cấp
CREDIT operating   +cost     ← quỹ vận hành xuất tiền (hoặc inventory)
```
Atomic: JournalEntry + Facility + Expense record

---

### 3.5 Order Revenue (Planned — chưa implement)

#### FR-REV-01: Order Payment — Thu tiền từ khách
```
Cash:     DEBIT cash_drawer  +cash_amount   ; CREDIT customer  +cash_amount
Transfer: DEBIT cash_drawer  +transfer_amt  ; CREDIT customer  +transfer_amt
```

---

## 4. Balance Calculation

```
fund_balance = Σ(DEBIT lines where fund_type = X) − Σ(CREDIT lines where fund_type = X)
```

- Tách biệt cash / transfer / total
- 5 real fund types xuất hiện trên dashboard
- External accounts (owner, supplier, customer, cash_shortage, cash_overage) bị loại khỏi "all balances"

---

## 5. Business Rules

| # | Rule |
|---|------|
| BR-1 | Fund-paid operations phải atomic (fail cả 2 nếu bất kỳ bước nào lỗi) |
| BR-2 | Withdrawal phải kiểm tra sufficient balance (cash và transfer riêng biệt) |
| BR-3 | Cashier không thể đóng ca nếu còn waiter shift đang mở |
| BR-4 | Variance ≥ 0.01 bắt buộc có reason + notes (≥ 10 chars) |
| BR-5 | Journal entries bất biến — không sửa, không xóa |
| BR-6 | Starting float mặc định xuất từ `operating` fund |
| BR-7 | expectedCash = startingFloat + receivedCash - distributedCash |

---

## 6. API Endpoints

| Method | Path | Description |
|--------|------|-------------|
| GET | `/api/manager/fund/journal-balances` | 5 real fund balances |
| GET | `/api/manager/fund/journal` | List entries (filters: fund_type, event_type, from_date, to_date, limit, offset) |
| GET | `/api/manager/fund/journal/:id` | Single entry detail |
| POST | `/api/manager/fund/deposit` | Manager deposit |
| POST | `/api/manager/fund/withdraw` | Manager withdrawal |
| POST | `/api/manager/fund/transfer` | Transfer between funds |

---

## 7. Frontend Views Liên Quan

| Route | View | Mục đích |
|-------|------|---------|
| `/manager/fund` | FundManagementView | Số dư + lịch sử journal |
| `/operating-expenses` | OperatingExpenseView | Chi tiêu vận hành (journal-based) |
| `/manager/ingredients` | IngredientManagementView | Nhập hàng từ quỹ |
| `/manager/facilities` | FacilityManagementView | Mua tài sản từ quỹ |
| `/cashier/shift-closure` | CashierShiftClosureV2 | Đóng ca cashier |
| `/cashier/handover` | CashierHandoverView | Waiter bàn giao tiền |
