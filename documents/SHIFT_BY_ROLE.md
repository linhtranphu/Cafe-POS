# Shift Management by Role

## Tổng quan

Hệ thống shift được tách riêng theo từng role để quản lý ca làm việc độc lập cho mỗi vai trò.

## Kiến trúc

### 1. Shift Model

```go
type Shift struct {
    ID            primitive.ObjectID
    Type          ShiftType          // MORNING, AFTERNOON, EVENING
    Status        ShiftStatus        // OPEN, CLOSED
    RoleType      RoleType           // waiter, cashier, barista
    UserID        primitive.ObjectID // ID của user
    UserName      string             // Tên user
    
    // Legacy fields (backward compatibility)
    WaiterID      primitive.ObjectID
    WaiterName    string
    CashierID     primitive.ObjectID
    CashierName   string
    
    StartCash     float64
    EndCash       float64
    TotalRevenue  float64
    TotalOrders   int
    StartedAt     time.Time
    EndedAt       *time.Time
}
```

### 2. Role Types

```go
type RoleType string

const (
    RoleWaiter  RoleType = "waiter"
    RoleCashier RoleType = "cashier"
    RoleBarista RoleType = "barista"
)
```

## Business Rules

### BR-10: Shift Independence by Role
- Mỗi user có thể có **1 shift OPEN** cho **mỗi role**
- Ví dụ: User A có thể đồng thời:
  - 1 waiter shift OPEN
  - 1 barista shift OPEN (nếu có cả 2 role)

### BR-11: Role-Specific Shift Data
- **Waiter Shift**: Track orders, revenue, cash handling
- **Barista Shift**: Track drinks prepared, working time
- **Cashier Shift**: Track cash reconciliation, payment oversight

### BR-12: Shift Closure Rules
- Waiter/Cashier: Phải nhập `end_cash` để đóng ca
- Barista: Không cần `end_cash`, chỉ cần đóng ca

## API Endpoints

### Start Shift
```
POST /api/shifts/start
Authorization: Bearer <token>

Request:
{
  "type": "MORNING",
  "start_cash": 1000000
}

Response:
{
  "id": "...",
  "type": "MORNING",
  "status": "OPEN",
  "role_type": "barista",
  "user_id": "...",
  "user_name": "Barista 1",
  "start_cash": 1000000,
  "started_at": "2024-01-28T08:00:00Z"
}
```

**Note**: `role_type` được tự động lấy từ JWT token

### Get Current Shift
```
GET /api/shifts/current
Authorization: Bearer <token>

Response:
{
  "id": "...",
  "type": "MORNING",
  "status": "OPEN",
  "role_type": "barista",
  "user_id": "...",
  "user_name": "Barista 1",
  "started_at": "2024-01-28T08:00:00Z"
}
```

**Note**: Trả về shift OPEN của user với role hiện tại

### Get My Shifts
```
GET /api/shifts/my
Authorization: Bearer <token>

Response:
[
  {
    "id": "...",
    "type": "MORNING",
    "status": "CLOSED",
    "role_type": "barista",
    "user_id": "...",
    "user_name": "Barista 1",
    "started_at": "2024-01-27T08:00:00Z",
    "ended_at": "2024-01-27T16:00:00Z"
  },
  ...
]
```

**Note**: Trả về tất cả shifts của user với role hiện tại

### End Shift
```
POST /api/shifts/:id/end
Authorization: Bearer <token>

Request:
{
  "end_cash": 1500000
}

Response:
{
  "id": "...",
  "status": "CLOSED",
  "end_cash": 1500000,
  "total_revenue": 500000,
  "total_orders": 25,
  "ended_at": "2024-01-28T16:00:00Z"
}
```

## Frontend Integration

### Shift Store
```javascript
// frontend/src/stores/shift.js
export const useShiftStore = defineStore('shift', {
  actions: {
    async fetchCurrentShift() {
      // Tự động lấy shift theo role của user hiện tại
      const response = await shiftService.getCurrentShift()
      this.currentShift = response.data
    }
  }
})
```

### Shift View
- Waiter: Hiển thị orders, revenue, cash
- Barista: Hiển thị drinks prepared, working time
- Cashier: Hiển thị cash reconciliation

## Migration Strategy

### Backward Compatibility
- Giữ lại `waiter_id`, `waiter_name`, `cashier_id`, `cashier_name`
- Khi tạo shift mới, set cả `user_id`/`user_name` và legacy fields
- Queries cũ vẫn hoạt động với `waiter_id`
- Queries mới sử dụng `user_id` + `role_type`

### Data Migration (Optional)
```javascript
// Script để migrate data cũ
db.shifts.updateMany(
  { role_type: { $exists: false } },
  [
    {
      $set: {
        role_type: "waiter",
        user_id: "$waiter_id",
        user_name: "$waiter_name"
      }
    }
  ]
)
```

## Benefits

1. **Separation of Concerns**: Mỗi role có shift riêng
2. **Flexible Tracking**: Track metrics khác nhau cho mỗi role
3. **Scalability**: Dễ thêm role mới (e.g., kitchen staff)
4. **Audit Trail**: Rõ ràng ai làm gì, khi nào
5. **Backward Compatible**: Không break existing code

## Example Scenarios

### Scenario 1: Barista Shift
```
08:00 - Barista1 mở ca (MORNING)
08:30 - Accept order #001
09:00 - Finish order #001
...
16:00 - Đóng ca (8 hours, 25 drinks)
```

### Scenario 2: Waiter Shift
```
08:00 - Waiter1 mở ca (MORNING, start_cash: 1M)
08:30 - Create order #001
09:00 - Collect payment #001
...
16:00 - Đóng ca (end_cash: 1.5M, revenue: 500K)
```

### Scenario 3: Multi-Role User
```
User A (có cả waiter và barista role):
- 08:00 - Mở waiter shift
- 08:05 - Mở barista shift
- Có thể switch giữa 2 shifts
- Đóng từng shift độc lập
```

## Testing

### Test Cases
1. ✅ User có thể mở shift cho role của mình
2. ✅ User không thể mở 2 shifts cùng role
3. ✅ User có thể mở shifts cho nhiều roles khác nhau
4. ✅ Get current shift trả về đúng shift theo role
5. ✅ Get my shifts chỉ trả về shifts của role hiện tại
6. ✅ Legacy queries vẫn hoạt động

## Future Enhancements

1. **Shift Handover**: Chuyển giao ca giữa users
2. **Shift Reports**: Báo cáo chi tiết theo role
3. **Shift Scheduling**: Lên lịch ca trước
4. **Multi-Facility**: Shift theo cơ sở
