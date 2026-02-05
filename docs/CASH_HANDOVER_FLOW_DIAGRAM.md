# 🔄 Cash Handover Flow Diagram

## 📊 System Architecture

```
┌─────────────────────────────────────────────────────────────────┐
│                         FRONTEND                                 │
├─────────────────────────────────────────────────────────────────┤
│                                                                   │
│  ┌──────────────┐    ┌──────────────────┐    ┌──────────────┐  │
│  │ ShiftView    │    │ CashierHandover  │    │ Cashier      │  │
│  │ (Waiter)     │    │ View (Cashier)   │    │ Dashboard    │  │
│  └──────┬───────┘    └────────┬─────────┘    └──────┬───────┘  │
│         │                     │                      │           │
│         └─────────────────────┼──────────────────────┘           │
│                               │                                  │
│  ┌────────────────────────────┴─────────────────────────────┐   │
│  │              Stores (shift.js, cashier.js)               │   │
│  └────────────────────────────┬─────────────────────────────┘   │
│                               │                                  │
│  ┌────────────────────────────┴─────────────────────────────┐   │
│  │              Services (handover.js)                       │   │
│  └────────────────────────────┬─────────────────────────────┘   │
│                               │                                  │
└───────────────────────────────┼──────────────────────────────────┘
                                │ HTTP/REST API
┌───────────────────────────────┼──────────────────────────────────┐
│                         BACKEND                                  │
├───────────────────────────────┼──────────────────────────────────┤
│                               │                                  │
│  ┌────────────────────────────┴─────────────────────────────┐   │
│  │         HTTP Handlers (cash_handover_handler.go)         │   │
│  └────────────────────────────┬─────────────────────────────┘   │
│                               │                                  │
│  ┌────────────────────────────┴─────────────────────────────┐   │
│  │      Services (cash_handover_service.go)                 │   │
│  │  - CreateHandover()                                      │   │
│  │  - ConfirmHandoverWithReconciliation()                   │   │
│  │  - ApproveDiscrepancy()                                  │   │
│  └────────────────────────────┬─────────────────────────────┘   │
│                               │                                  │
│  ┌────────────────────────────┴─────────────────────────────┐   │
│  │      Repositories (MongoDB)                              │   │
│  │  - CashHandoverRepository                                │   │
│  │  - CashDiscrepancyRepository                             │   │
│  └────────────────────────────┬─────────────────────────────┘   │
│                               │                                  │
└───────────────────────────────┼──────────────────────────────────┘
                                │
┌───────────────────────────────┼──────────────────────────────────┐
│                         DATABASE                                 │
├───────────────────────────────┼──────────────────────────────────┤
│                               │                                  │
│  ┌────────────────┐  ┌────────────────┐  ┌─────────────────┐   │
│  │ cash_handovers │  │ cash_          │  │ shifts          │   │
│  │                │  │ discrepancies  │  │ cashier_shifts  │   │
│  └────────────────┘  └────────────────┘  └─────────────────┘   │
│                                                                  │
└──────────────────────────────────────────────────────────────────┘
```

---

## 🔄 Workflow: Partial Handover (No Discrepancy)

```
WAITER                    SYSTEM                    CASHIER
  │                         │                         │
  │ 1. Click "Bàn giao      │                         │
  │    một phần"            │                         │
  │─────────────────────────>                         │
  │                         │                         │
  │ 2. Enter amount         │                         │
  │    (500,000 VND)        │                         │
  │─────────────────────────>                         │
  │                         │                         │
  │                         │ 3. Create handover      │
  │                         │    status: PENDING      │
  │                         │    declared: 500k       │
  │                         │                         │
  │ 4. Show pending banner  │                         │
  │<─────────────────────────                         │
  │                         │                         │
  │                         │ 5. Send notification    │
  │                         │─────────────────────────>
  │                         │                         │
  │                         │                         │ 6. View pending
  │                         │                         │    handover
  │                         │                         │
  │                         │                         │ 7. Count cash
  │                         │                         │    (500,000 VND)
  │                         │                         │
  │                         │                         │ 8. Enter actual
  │                         │                         │    amount: 500k
  │                         │<─────────────────────────
  │                         │                         │
  │                         │ 9. Calculate            │
  │                         │    discrepancy: 0       │
  │                         │    status: CONFIRMED    │
  │                         │                         │
  │                         │ 10. Update shifts:      │
  │                         │     Waiter:             │
  │                         │     - handed_over += 500k
  │                         │     - remaining -= 500k │
  │                         │     Cashier:            │
  │                         │     - received += 500k  │
  │                         │                         │
  │ 11. Update UI           │                         │
  │     Remove pending      │                         │
  │     Update cash amounts │                         │
  │<─────────────────────────                         │
  │                         │                         │
  │                         │ 12. Confirm success     │
  │                         │─────────────────────────>
  │                         │                         │
```

---

## ⚠️ Workflow: Handover with Small Discrepancy

```
WAITER                    SYSTEM                    CASHIER
  │                         │                         │
  │ 1. Declare 500k         │                         │
  │─────────────────────────>                         │
  │                         │                         │
  │                         │ 2. Create handover      │
  │                         │    declared: 500k       │
  │                         │    status: PENDING      │
  │                         │                         │
  │                         │ 3. Notify cashier       │
  │                         │─────────────────────────>
  │                         │                         │
  │                         │                         │ 4. Count cash
  │                         │                         │    (480,000 VND)
  │                         │                         │    Short 20k!
  │                         │                         │
  │                         │                         │ 5. Enter actual:
  │                         │                         │    480k
  │                         │                         │    Reason: COUNTING_ERROR
  │                         │                         │    Responsibility: WAITER
  │                         │<─────────────────────────
  │                         │                         │
  │                         │ 6. Calculate:           │
  │                         │    discrepancy: -20k    │
  │                         │    type: SHORTAGE       │
  │                         │                         │
  │                         │ 7. Check threshold:     │
  │                         │    20k < 100k           │
  │                         │    No approval needed   │
  │                         │                         │
  │                         │ 8. Create discrepancy   │
  │                         │    record               │
  │                         │                         │
  │                         │ 9. Update shifts:       │
  │                         │    Waiter:              │
  │                         │    - handed_over += 480k│
  │                         │    - remaining -= 500k  │
  │                         │    - discrepancy += -20k│
  │                         │    Cashier:             │
  │                         │    - received += 480k   │
  │                         │    - discrepancy += -20k│
  │                         │    - discrepancy_count++│
  │                         │                         │
  │ 10. Notify waiter       │                         │
  │     "Confirmed with     │                         │
  │      20k shortage"      │                         │
  │<─────────────────────────                         │
  │                         │                         │
```

---

## 🚨 Workflow: Large Discrepancy (Requires Manager Approval)

```
WAITER          SYSTEM          CASHIER         MANAGER
  │               │               │               │
  │ 1. Declare    │               │               │
  │    1,000,000  │               │               │
  │───────────────>               │               │
  │               │               │               │
  │               │ 2. Notify     │               │
  │               │───────────────>               │
  │               │               │               │
  │               │               │ 3. Count:     │
  │               │               │    800,000    │
  │               │               │    Short 200k!│
  │               │               │               │
  │               │               │ 4. Enter:     │
  │               │               │    actual: 800k
  │               │               │    reason     │
  │               │<───────────────               │
  │               │               │               │
  │               │ 5. Calculate: │               │
  │               │    -200k      │               │
  │               │               │               │
  │               │ 6. Check:     │               │
  │               │    200k > 100k│               │
  │               │    REQUIRES   │               │
  │               │    APPROVAL!  │               │
  │               │               │               │
  │               │ 7. Set status:│               │
  │               │    DISCREPANCY│               │
  │               │    requires_  │               │
  │               │    approval:  │               │
  │               │    true       │               │
  │               │               │               │
  │ 8. Notify:    │               │               │
  │    "Waiting   │               │               │
  │     manager   │               │               │
  │     approval" │               │               │
  │<───────────────               │               │
  │               │               │               │
  │               │ 9. Notify manager             │
  │               │───────────────────────────────>
  │               │               │               │
  │               │               │               │ 10. Review
  │               │               │               │     details
  │               │               │               │
  │               │               │               │ 11. Approve
  │               │<───────────────────────────────
  │               │               │               │
  │               │ 12. Update:   │               │
  │               │     status:   │               │
  │               │     CONFIRMED │               │
  │               │     Update    │               │
  │               │     shifts    │               │
  │               │               │               │
  │ 13. Notify:   │               │               │
  │     "Approved"│               │               │
  │<───────────────               │               │
  │               │               │               │
```

---

## 🏁 Workflow: Handover and End Shift

```
WAITER                    SYSTEM                    CASHIER
  │                         │                         │
  │ 1. Click "Bàn giao      │                         │
  │    và đóng ca"          │                         │
  │─────────────────────────>                         │
  │                         │                         │
  │ 2. Enter end_cash: 0    │                         │
  │─────────────────────────>                         │
  │                         │                         │
  │                         │ 3. Create handover:     │
  │                         │    type: END_SHIFT      │
  │                         │    amount: remaining_cash
  │                         │    status: PENDING      │
  │                         │                         │
  │ 4. Show "Waiting for    │                         │
  │    cashier to close     │                         │
  │    shift"               │                         │
  │<─────────────────────────                         │
  │                         │                         │
  │                         │ 5. Notify cashier       │
  │                         │─────────────────────────>
  │                         │                         │
  │                         │                         │ 6. Confirm
  │                         │<─────────────────────────
  │                         │                         │
  │                         │ 7. Update waiter shift: │
  │                         │    - handed_over += amt │
  │                         │    - remaining = 0      │
  │                         │    - status = CLOSED    │
  │                         │    - ended_at = now     │
  │                         │                         │
  │                         │ 8. Lock orders:         │
  │                         │    SERVED → LOCKED      │
  │                         │                         │
  │                         │ 9. Update cashier:      │
  │                         │    - received += amt    │
  │                         │                         │
  │ 10. Redirect to         │                         │
  │     shift history       │                         │
  │<─────────────────────────                         │
  │                         │                         │
  │ 11. Shift CLOSED        │                         │
  │     Cannot reopen       │                         │
  │                         │                         │
```

---

## 📊 Data Flow

```
┌─────────────────────────────────────────────────────────────┐
│                    CREATE HANDOVER                          │
└─────────────────────────────────────────────────────────────┘
                              │
                              ▼
                    ┌──────────────────┐
                    │ Validate Request │
                    │ - Shift open?    │
                    │ - Amount valid?  │
                    │ - Cashier open?  │
                    └────────┬─────────┘
                             │
                             ▼
                    ┌──────────────────┐
                    │ Create Handover  │
                    │ Record           │
                    │ status: PENDING  │
                    └────────┬─────────┘
                             │
                             ▼
                    ┌──────────────────┐
                    │ Send Notification│
                    │ to Cashier       │
                    └──────────────────┘

┌─────────────────────────────────────────────────────────────┐
│                    CONFIRM HANDOVER                         │
└─────────────────────────────────────────────────────────────┘
                              │
                              ▼
                    ┌──────────────────┐
                    │ Get Handover     │
                    │ Record           │
                    └────────┬─────────┘
                             │
                             ▼
                    ┌──────────────────┐
                    │ Calculate        │
                    │ Discrepancy      │
                    │ actual - declared│
                    └────────┬─────────┘
                             │
                    ┌────────┴─────────┐
                    │                  │
         discrepancy = 0    discrepancy ≠ 0
                    │                  │
                    ▼                  ▼
          ┌──────────────┐   ┌──────────────────┐
          │ Status:      │   │ Check Threshold  │
          │ CONFIRMED    │   │                  │
          └──────┬───────┘   └────────┬─────────┘
                 │                    │
                 │           ┌────────┴─────────┐
                 │           │                  │
                 │    < threshold      > threshold
                 │           │                  │
                 │           ▼                  ▼
                 │  ┌──────────────┐  ┌──────────────┐
                 │  │ Status:      │  │ Status:      │
                 │  │ CONFIRMED    │  │ DISCREPANCY  │
                 │  │              │  │ requires_    │
                 │  │              │  │ approval:true│
                 │  └──────┬───────┘  └──────┬───────┘
                 │         │                 │
                 │         │                 │
                 └─────────┴─────────────────┘
                           │
                           ▼
                 ┌──────────────────┐
                 │ Create           │
                 │ Discrepancy      │
                 │ Record (if any)  │
                 └────────┬─────────┘
                          │
                          ▼
                 ┌──────────────────┐
                 │ Update Cash      │
                 │ Amounts          │
                 │ (if confirmed)   │
                 └────────┬─────────┘
                          │
                          ▼
                 ┌──────────────────┐
                 │ Send             │
                 │ Notifications    │
                 └──────────────────┘
```

---

## 🔐 Security Flow

```
┌─────────────────────────────────────────────────────────────┐
│                    REQUEST VALIDATION                        │
└─────────────────────────────────────────────────────────────┘

HTTP Request
     │
     ▼
┌──────────────────┐
│ Authentication   │ ──> JWT Token Valid?
│ Middleware       │     │
└────────┬─────────┘     │ No
         │ Yes           ▼
         │          401 Unauthorized
         ▼
┌──────────────────┐
│ Authorization    │ ──> Role Allowed?
│ Middleware       │     │
└────────┬─────────┘     │ No
         │ Yes           ▼
         │          403 Forbidden
         ▼
┌──────────────────┐
│ Handler          │
│ Validation       │ ──> Input Valid?
└────────┬─────────┘     │
         │ Yes           │ No
         │               ▼
         │          400 Bad Request
         ▼
┌──────────────────┐
│ Business Logic   │ ──> Business Rules OK?
│ Validation       │     │
└────────┬─────────┘     │ No
         │ Yes           ▼
         │          422 Unprocessable
         ▼
┌──────────────────┐
│ Execute          │
│ Operation        │
└────────┬─────────┘
         │
         ▼
    200 Success
```

---

## 📈 Performance Considerations

### Database Queries Optimization:

```
1. Indexes:
   - cash_handovers: (cashier_id, status)
   - cash_handovers: (waiter_shift_id)
   - cash_handovers: (handover_at DESC)

2. Query Patterns:
   - Get pending by cashier: Use compound index
   - Get history by shift: Use shift_id index
   - Get today's handovers: Use date range with index

3. Caching Strategy:
   - Cache pending count for cashier dashboard
   - Invalidate on handover create/confirm
   - TTL: 30 seconds
```

### Frontend Optimization:

```
1. Lazy Loading:
   - CashierHandoverView loaded on demand
   - Handover history paginated

2. Real-time Updates:
   - Poll pending handovers every 30s
   - Show loading states
   - Optimistic UI updates

3. Mobile Performance:
   - Minimize re-renders
   - Use computed properties
   - Debounce input handlers
```

---

**This diagram should be referenced during implementation to understand the complete flow!**
