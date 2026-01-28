# 🔄 State Machine Implementation - With Barista Role

## 📋 Overview

Đã implement state machine mới theo đúng nghiệp vụ với vai trò Barista được tách biệt rõ ràng.

## 🎯 New State Machine

```
CREATED → PAID → QUEUED → IN_PROGRESS → READY → SERVED → LOCKED
                    ↓          ↓
                CANCELLED → LOCKED
```

### States Definition

| State | Ý nghĩa | Ai trigger |
|-------|---------|------------|
| **CREATED** | Order chưa thanh toán | System |
| **PAID** | Đã thu tiền, chưa giao cho pha chế | Waiter |
| **QUEUED** | Đã gửi cho barista, chờ nhận | Waiter |
| **IN_PROGRESS** | Barista đã nhận và đang pha | Barista |
| **READY** | Pha xong, chờ giao | Barista |
| **SERVED** | Đã giao cho khách | Waiter |
| **CANCELLED** | Đã hủy | Waiter/Cashier |
| **LOCKED** | Đã chốt ca | Cashier |

## 🔐 Business Rules Implemented

### BR-01: State Transitions
```go
transitions := map[OrderStatus][]OrderStatus{
    StatusCreated:    {StatusPaid, StatusCancelled},
    StatusPaid:       {StatusPaid, StatusQueued, StatusCancelled},
    StatusQueued:     {StatusInProgress, StatusCancelled},
    StatusInProgress: {StatusReady},
    StatusReady:      {StatusServed},
    StatusServed:     {StatusLocked},
    StatusCancelled:  {StatusLocked},
    StatusLocked:     {},
}
```

### BR-06: Only Barista can move order to IN_PROGRESS
```go
func (s *OrderService) AcceptOrder(ctx context.Context, id primitive.ObjectID, baristaID, baristaName string) (*order.Order, error) {
    // Only barista can accept order
    // Moves from QUEUED → IN_PROGRESS
}
```

### BR-07: No modification after barista accepts
```go
func (o *Order) CanModify() bool {
    // Once order enters IN_PROGRESS, no modification or refund is allowed
    return o.Status == StatusCreated || o.Status == StatusPaid || o.Status == StatusQueued
}

func (s *OrderService) CancelOrder(...) {
    // Cannot cancel once barista has accepted
    if o.Status == StatusInProgress || o.Status == StatusReady {
        return errors.New("cannot cancel order after barista has started preparing")
    }
}
```

### BR-08: Payment adjustments only before QUEUED
```go
func (o *Order) IsEditable() bool {
    return o.Status == StatusCreated || o.Status == StatusPaid
}

func (o *Order) CanRefund() bool {
    return o.Status == StatusPaid && o.AmountPaid > 0
}
```

### BR-09: READY indicates drink completed
```go
func (s *OrderService) FinishPreparing(ctx context.Context, id primitive.ObjectID) (*order.Order, error) {
    // Barista marks order as READY
    // Moves from IN_PROGRESS → READY
}
```

## 🎭 Role Responsibilities

### Waiter
**Actions:**
- Create order (→ CREATED)
- Collect payment (CREATED → PAID)
- Edit order (CREATED, PAID only)
- Send to bar (PAID → QUEUED)
- Deliver drink (READY → SERVED)
- Cancel order (before IN_PROGRESS)

**Endpoints:**
```
POST   /api/waiter/orders              # Create order
POST   /api/waiter/orders/:id/payment  # Collect payment
PUT    /api/waiter/orders/:id/edit     # Edit order
POST   /api/waiter/orders/:id/send     # Send to bar
POST   /api/waiter/orders/:id/serve    # Deliver drink
GET    /api/waiter/orders               # Get my orders
```

### Barista
**Actions:**
- View queue (QUEUED orders)
- Accept order (QUEUED → IN_PROGRESS)
- Mark as ready (IN_PROGRESS → READY)
- View my orders (IN_PROGRESS + READY)

**Endpoints:**
```
GET    /api/barista/orders/queue       # View queued orders
POST   /api/barista/orders/:id/accept  # Accept order
POST   /api/barista/orders/:id/ready   # Mark as ready
GET    /api/barista/orders/my          # Get my orders
GET    /api/barista/orders/:id         # View order details
```

### Cashier
**Actions:**
- View all orders
- Cancel order (before IN_PROGRESS)
- Refund (PAID only, before QUEUED)
- Lock order (SERVED/CANCELLED → LOCKED)
- Close shift (locks all orders)

**Endpoints:**
```
GET    /api/cashier/orders              # View all orders
POST   /api/cashier/orders/:id/cancel   # Cancel order
POST   /api/cashier/orders/:id/refund   # Refund
POST   /api/cashier/orders/:id/lock     # Lock order
POST   /api/cashier/shifts/:id/close    # Close shift
```

### Manager
**Actions:**
- All of the above
- Full access to all operations

## 📊 Order Model Changes

### New Fields
```go
type Order struct {
    // ... existing fields ...
    
    // Barista tracking
    BaristaID       primitive.ObjectID `bson:"barista_id,omitempty" json:"barista_id,omitempty"`
    BaristaName     string             `bson:"barista_name,omitempty" json:"barista_name,omitempty"`
    
    // Timestamps
    QueuedAt        *time.Time         `bson:"queued_at,omitempty" json:"queued_at,omitempty"`
    AcceptedAt      *time.Time         `bson:"accepted_at,omitempty" json:"accepted_at,omitempty"`
    ReadyAt         *time.Time         `bson:"ready_at,omitempty" json:"ready_at,omitempty"`
    
    // Removed: SentToBarAt (replaced by QueuedAt)
}
```

### New Methods
```go
func (o *Order) CanModify() bool
func (o *Order) CanRefund() bool
```

## 🔄 Workflow Examples

### Happy Path
```
1. Waiter creates order → CREATED
2. Waiter collects payment → PAID
3. Waiter sends to bar → QUEUED
4. Barista accepts order → IN_PROGRESS
5. Barista finishes preparing → READY
6. Waiter delivers to customer → SERVED
7. Cashier closes shift → LOCKED
```

### Edit Before Barista
```
1. Order is PAID
2. Waiter edits order ✅ (allowed)
3. Waiter sends to bar → QUEUED
4. Waiter tries to edit ✅ (still allowed)
5. Barista accepts → IN_PROGRESS
6. Waiter tries to edit ❌ (not allowed)
```

### Cancel Scenarios
```
# Can cancel:
- CREATED → CANCELLED ✅
- PAID → CANCELLED ✅
- QUEUED → CANCELLED ✅

# Cannot cancel:
- IN_PROGRESS → CANCELLED ❌ (barista already working)
- READY → CANCELLED ❌ (drink already made)
- SERVED → CANCELLED ❌ (already delivered)
```

### Refund Scenarios
```
# Can refund:
- PAID (before QUEUED) ✅

# Cannot refund:
- QUEUED ❌ (already sent to barista)
- IN_PROGRESS ❌ (barista working)
- READY ❌ (drink made)
- SERVED ❌ (already delivered)
```

## 📈 Metrics & Tracking

### Queue Metrics
```go
// Time in queue
queueTime = AcceptedAt - QueuedAt

// Preparation time
prepTime = ReadyAt - AcceptedAt

// Delivery time
deliveryTime = ServedAt - ReadyAt

// Total time
totalTime = ServedAt - CreatedAt
```

### Barista Performance
```go
// Orders accepted by barista
SELECT COUNT(*) WHERE barista_id = X AND status >= IN_PROGRESS

// Average prep time
SELECT AVG(ready_at - accepted_at) WHERE barista_id = X

// Orders in progress
SELECT COUNT(*) WHERE barista_id = X AND status = IN_PROGRESS
```

## 🎨 Frontend Updates Needed

### Order Status Display
```javascript
const statuses = [
  { value: 'CREATED', label: 'Mới tạo', icon: '🆕', color: 'gray' },
  { value: 'PAID', label: 'Đã thu', icon: '💰', color: 'green' },
  { value: 'QUEUED', label: 'Chờ pha', icon: '⏳', color: 'yellow' },
  { value: 'IN_PROGRESS', label: 'Đang pha', icon: '🍹', color: 'blue' },
  { value: 'READY', label: 'Sẵn sàng', icon: '✅', color: 'purple' },
  { value: 'SERVED', label: 'Đã giao', icon: '🎉', color: 'green' },
  { value: 'CANCELLED', label: 'Đã hủy', icon: '❌', color: 'red' },
  { value: 'LOCKED', label: 'Đã khóa', icon: '🔒', color: 'gray' }
]
```

### Waiter Actions
```javascript
// Show "Gửi bar" button only for PAID orders
if (order.status === 'PAID' && order.amount_due <= 0) {
  <button onClick={() => sendToBar(order.id)}>🍹 Gửi bar</button>
}

// Show "Giao khách" button only for READY orders
if (order.status === 'READY') {
  <button onClick={() => serveOrder(order.id)}>🎉 Giao khách</button>
}

// Can edit only before IN_PROGRESS
if (['CREATED', 'PAID', 'QUEUED'].includes(order.status)) {
  <button onClick={() => editOrder(order)}>✏️ Sửa</button>
}
```

### Barista View (New)
```javascript
// Queue view
<div className="queue">
  {queuedOrders.map(order => (
    <OrderCard 
      order={order}
      action={() => acceptOrder(order.id)}
      actionLabel="Nhận order"
    />
  ))}
</div>

// My orders view
<div className="my-orders">
  {myOrders.map(order => (
    <OrderCard 
      order={order}
      action={order.status === 'IN_PROGRESS' 
        ? () => markReady(order.id)
        : null
      }
      actionLabel={order.status === 'IN_PROGRESS' ? 'Hoàn tất' : 'Chờ giao'}
    />
  ))}
</div>
```

## 🔧 Migration Notes

### Database Migration
```javascript
// Update existing orders
db.orders.updateMany(
  { status: 'IN_PROGRESS', sent_to_bar_at: { $exists: true } },
  {
    $set: { 
      status: 'QUEUED',
      queued_at: '$sent_to_bar_at'
    },
    $unset: { sent_to_bar_at: '' }
  }
)

// Orders that were IN_PROGRESS without barista
db.orders.updateMany(
  { status: 'IN_PROGRESS', barista_id: { $exists: false } },
  {
    $set: { status: 'QUEUED' }
  }
)
```

### User Migration
```javascript
// Create barista users
db.users.insertMany([
  {
    username: 'barista1',
    password: hashedPassword,
    role: 'barista',
    name: 'Barista 1',
    active: true,
    created_at: new Date(),
    updated_at: new Date()
  }
])
```

## ✅ Testing Checklist

### State Transitions
- [ ] CREATED → PAID (waiter collects payment)
- [ ] PAID → QUEUED (waiter sends to bar)
- [ ] QUEUED → IN_PROGRESS (barista accepts)
- [ ] IN_PROGRESS → READY (barista finishes)
- [ ] READY → SERVED (waiter delivers)
- [ ] SERVED → LOCKED (cashier closes shift)
- [ ] PAID → CANCELLED (before barista accepts)
- [ ] QUEUED → CANCELLED (before barista accepts)
- [ ] IN_PROGRESS → CANCELLED (should fail)

### Business Rules
- [ ] BR-06: Only barista can accept order
- [ ] BR-07: Cannot modify after IN_PROGRESS
- [ ] BR-08: Cannot refund after QUEUED
- [ ] BR-09: READY state works correctly

### Role Permissions
- [ ] Waiter can create, pay, send, serve
- [ ] Barista can view queue, accept, mark ready
- [ ] Cashier can cancel (before IN_PROGRESS), refund (before QUEUED)
- [ ] Manager has full access

### Edge Cases
- [ ] Edit order in QUEUED state
- [ ] Cancel order in QUEUED state
- [ ] Refund attempt in QUEUED state (should fail)
- [ ] Multiple baristas accepting same order
- [ ] Order stuck in IN_PROGRESS

## 🎉 Benefits

### For Business
1. **Clear Accountability**: Biết ai làm gì, khi nào
2. **Quality Control**: Barista chịu trách nhiệm về chất lượng
3. **Performance Tracking**: Đo được thời gian pha chế
4. **Queue Management**: Quản lý hàng đợi hiệu quả

### For Operations
1. **No Confusion**: Không còn "waiter tự bấm đã pha"
2. **Audit Trail**: Đầy đủ lịch sử thay đổi
3. **Error Detection**: Phát hiện lỗi pha chế
4. **Resource Planning**: Biết barista nào đang bận

### For Users
1. **Waiter**: Biết order nào sẵn sàng để giao
2. **Barista**: Có queue riêng, không bị nhầm lẫn
3. **Cashier**: Kiểm soát tốt hơn
4. **Manager**: Có metrics để đánh giá

## 🚀 Next Steps

1. ✅ Backend implementation (Done)
2. ⏳ Frontend updates (Pending)
   - Update OrderView với QUEUED, READY states
   - Tạo BaristaView mới
   - Update status badges và colors
3. ⏳ Database migration (Pending)
4. ⏳ Testing (Pending)
5. ⏳ Documentation update (Pending)
6. ⏳ Training (Pending)
