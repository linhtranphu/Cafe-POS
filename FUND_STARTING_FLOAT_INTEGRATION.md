# Tích hợp Starting Float vào Quản lý Quỹ

## Tổng quan

Khi cashier mở ca với starting_float, hệ thống sẽ TỰ ĐỘNG tạo một giao dịch `withdrawal` trong `fund_transactions` để trừ tiền từ quỹ. Khi đóng ca, tiền được trả lại qua `fund_handover`.

## Luồng tiền

### 1. Khi Cashier mở ca (Starting Float)
- **Hành động**: Tạo `fund_transaction` với type = `withdrawal`
- **Reason**: "Tiền đầu ca cho thu ngân [Tên]"
- **Hiển thị**: Giao dịch RÚT tiền từ quỹ
- **Màu sắc**: Đỏ (red)
- **Ảnh hưởng**: Giảm số dư quỹ NGAY LẬP TỨC

### 2. Khi Cashier đóng ca (Handover)
- **Nguồn dữ liệu**: `fund_handovers`
- **Hiển thị**: Giao dịch NỘP tiền vào quỹ
- **Loại**: `fund_handover`
- **Màu sắc**: Xanh dương (blue)
- **Ảnh hưởng**: Tăng số dư quỹ

## Thay đổi Backend

### 1. CashierShiftService - StartCashierShift()

Khi mở ca với starting_float > 0, tự động tạo fund_transaction:

```go
// If starting float > 0, create fund withdrawal transaction
// Use MongoDB transaction for atomicity
if startingFloat > 0 {
    session, err := s.mongoClient.StartSession()
    if err != nil {
        return nil, fmt.Errorf("failed to start session: %w", err)
    }
    defer session.EndSession(ctx)

    err = mongo.WithSession(ctx, session, func(sc mongo.SessionContext) error {
        if err := session.StartTransaction(); err != nil {
            return err
        }

        // Create cashier shift first to get the ID
        if err := s.cashierShiftRepo.Create(sc, shift); err != nil {
            session.AbortTransaction(sc)
            return fmt.Errorf("failed to create cashier shift: %w", err)
        }

        // Create fund withdrawal for starting float
        reason := fmt.Sprintf("Tiền đầu ca cho thu ngân %s", cashierName)
        _, _, err := s.fundService.CreateWithdrawal(
            sc,
            startingFloat, // cash amount
            0,             // transfer amount
            reason,
            cashierID,
            cashierName,
            "cashier",
        )
        if err != nil {
            session.AbortTransaction(sc)
            return fmt.Errorf("failed to create fund withdrawal for starting float: %w", err)
        }

        return session.CommitTransaction(sc)
    })

    if err != nil {
        return nil, err
    }
}
```

**Quan trọng**: 
- Sử dụng MongoDB transaction để đảm bảo atomicity
- Nếu tạo fund_transaction thất bại, cashier shift cũng sẽ bị rollback
- Reason rõ ràng để dễ audit

### 2. FundService - CalculateCurrentBalance()

Logic đơn giản, KHÔNG cần trừ starting_float vì đã có fund_transaction:

```go
func (s *FundService) CalculateCurrentBalance(ctx context.Context) (*fund.FundBalance, error) {
    balance := &fund.FundBalance{
        Cash:     0,
        Transfer: 0,
        Total:    0,
    }

    // 1. Add fund handovers (cashier → fund)
    // 2. Add deposits
    // 3. Subtract withdrawals (bao gồm cả starting_float)
    
    balance.Total = balance.Cash + balance.Transfer
    return balance, nil
}
```

### 3. FundService - GetAggregatedTransactionHistory()

KHÔNG cần đọc từ cashier_shifts vì starting_float đã là fund_transaction:

```go
// Chỉ cần đọc từ:
// 1. fund_transactions (deposits và withdrawals, bao gồm starting_float)
// 2. fund_handovers (cashier → fund)
```

## Thay đổi Frontend

### 1. Constants (fund.js)

Đã có sẵn:
```javascript
export const TRANSACTION_TYPES = {
  STARTING_FLOAT: 'starting_float',
  // ...
}

export const TRANSACTION_TYPE_LABELS = {
  [TRANSACTION_TYPES.STARTING_FLOAT]: 'Tiền đầu ca',
  // ...
}

export const TRANSACTION_TYPE_ICONS = {
  [TRANSACTION_TYPES.STARTING_FLOAT]: '💰',
  // ...
}
```

### 2. FundManagementView.vue

#### Helper Functions
```javascript
const getSourceLabel = (type) => {
  const labels = {
    starting_float: 'Tiền đầu ca',
    // ...
  }
  return labels[type] || 'Khác'
}

const getSourceBadgeClass = (type) => {
  const classes = {
    starting_float: 'bg-orange-100 text-orange-700',
    // ...
  }
  return classes[type] || 'bg-gray-100 text-gray-700'
}
```

#### Computed Properties
```javascript
const todayWithdrawals = computed(() => {
  const today = new Date().toDateString()
  return transactions.value
    .filter(tx => 
      new Date(tx.timestamp).toDateString() === today && 
      (tx.type === 'withdrawal' || tx.type === 'starting_float')
    )
    .reduce((sum, tx) => sum + tx.total_amount, 0)
})
```

#### UI Display
```vue
<div class="flex justify-between items-center">
  <span class="text-gray-600">📤 Rút ra & Tiền đầu ca</span>
  <span class="font-semibold text-red-600">{{ formatCurrency(todayWithdrawals) }}</span>
</div>
```

## Ưu điểm của thiết kế này

1. **Chủ động**: Hành động mở ca NGAY LẬP TỨC tạo giao dịch trừ tiền
2. **Nhất quán**: Tất cả giao dịch đều được lưu trong `fund_transactions`
3. **Audit trail**: Dễ dàng theo dõi lịch sử giao dịch
4. **Atomicity**: Sử dụng MongoDB transaction đảm bảo tính toàn vẹn dữ liệu
5. **Đơn giản**: Logic tính balance đơn giản, không cần xử lý đặc biệt

## Công thức tính số dư quỹ

```
Số dư quỹ = 
  + Tổng fund_handovers (cashier → fund)
  + Tổng deposits (manager → fund)
  - Tổng withdrawals (fund → manager, bao gồm starting_float)
```

**Đơn giản**: Tất cả giao dịch đều được lưu trong `fund_transactions`, không cần xử lý đặc biệt.

## Phân loại giao dịch

| Loại | Nguồn | Hướng | Màu sắc | Ảnh hưởng | Ghi chú |
|------|-------|-------|---------|-----------|---------|
| `deposit` | Manager bổ sung | → Quỹ | Xanh lá | + | Thủ công |
| `withdrawal` | Manager rút / Starting float | Quỹ → | Đỏ | - | Thủ công hoặc tự động khi mở ca |
| `fund_handover` | Cashier bàn giao | → Quỹ | Xanh dương | + | Tự động khi đóng ca |

## Testing

### Test Backend
```bash
# Compile backend
cd backend && go build -o /dev/null main.go

# Test API
curl http://localhost:3000/api/manager/fund/balance
curl http://localhost:3000/api/manager/fund/transactions
```

### Test Frontend
1. **Kiểm tra số dư ban đầu**: Ghi nhớ số dư quỹ hiện tại (ví dụ: 500,000)
2. **Mở ca cashier**: Mở ca với starting_float = 100,000
3. **Kiểm tra view quỹ NGAY**: http://localhost:5173/#/manager/fund
   - Số dư quỹ giảm NGAY LẬP TỨC xuống 400,000
   - Withdrawal "Tiền đầu ca cho thu ngân [Tên]" hiển thị trong lịch sử
   - Badge màu đỏ với label "Rút ra"
   - "Rút ra" trong phân tích nguồn tăng 100,000
4. **Đóng ca cashier**: Handover với cash_amount = 100,000 (+ doanh thu nếu có)
5. **Kiểm tra lại view quỹ**:
   - Số dư quỹ tăng trở lại 500,000 (hoặc hơn nếu có doanh thu)
   - Fund handover hiển thị trong lịch sử
   - Badge màu xanh dương với label "Doanh thu"

## Files đã thay đổi

### Backend
- `backend/application/services/cashier_shift_service.go`
  - Thêm FundService dependency
  - Cập nhật `StartCashierShift()` để tạo fund_transaction khi starting_float > 0
  - Sử dụng MongoDB transaction cho atomicity
- `backend/main.go`
  - Di chuyển khởi tạo FundService lên trước CashierShiftService
  - Truyền fundService vào CashierShiftService constructor

### Frontend
- `frontend/src/views/FundManagementView.vue`
  - Cập nhật `todayWithdrawals` để chỉ filter `withdrawal` (không cần `starting_float` nữa)
  - Cập nhật `getSourceLabel()` để hiển thị "Rút ra" cho withdrawal
  - Cập nhật label "Rút ra" (không cần "& Tiền đầu ca" nữa)

## Kết luận

Thiết kế này đảm bảo rằng hành động MỞ CA sẽ NGAY LẬP TỨC tạo giao dịch trừ tiền từ quỹ. Điều này giúp:
- Số dư quỹ luôn chính xác theo thời gian thực
- Audit trail đầy đủ cho mọi giao dịch
- Logic đơn giản và dễ hiểu
- Tính toàn vẹn dữ liệu được đảm bảo bằng MongoDB transaction
