# BR-13: Barista Must Have Open Shift Before Accepting Orders

## Business Rule

**BR-13**: Barista must open a shift before accepting orders from the queue.

### Rationale
- **Accountability**: Track which barista worked on which orders
- **Time Tracking**: Measure barista working hours
- **Shift Management**: Ensure proper shift opening/closing procedures
- **Audit Trail**: Link orders to specific shifts for reporting

## Implementation

### Backend Changes

#### 1. Order Service (`order_service.go`)

Added shift validation in `AcceptOrder()` method:

```go
// AcceptOrder - BR-06: Only Barista can move order to IN_PROGRESS
// BR-13: Barista must have an open shift to accept orders
func (s *OrderService) AcceptOrder(ctx context.Context, id primitive.ObjectID, baristaID, baristaName string) (*order.Order, error) {
	o, err := s.orderRepo.FindByID(ctx, id)
	if err != nil {
		return nil, err
	}

	if !o.CanTransitionTo(order.StatusInProgress) {
		return nil, errors.New("cannot accept order in current state")
	}

	// BR-13: Check if barista has an open shift
	baristaOID, _ := primitive.ObjectIDFromHex(baristaID)
	shift, err := s.shiftRepo.FindOpenShiftByUser(ctx, baristaOID, order.RoleBarista)
	if err != nil || shift == nil {
		return nil, errors.New("barista must open a shift before accepting orders")
	}

	now := time.Now()
	o.Status = order.StatusInProgress
	o.BaristaID = baristaOID
	o.BaristaName = baristaName
	o.AcceptedAt = &now

	if err := s.orderRepo.Update(ctx, id, o); err != nil {
		return nil, err
	}
	return o, nil
}
```

**Key Points**:
- Checks for open shift using `FindOpenShiftByUser()`
- Validates shift is for `RoleBarista`
- Returns clear error message if no shift found
- Only proceeds if shift exists and is OPEN

### API Behavior

#### Endpoint
```
POST /api/barista/orders/:id/accept
Authorization: Bearer <token>
```

#### Success Response (200 OK)
```json
{
  "id": "697a1f99a928f9d7ff2311df",
  "order_number": "ORD-001",
  "status": "IN_PROGRESS",
  "barista_id": "697a1234a928f9d7ff231111",
  "barista_name": "Barista 1",
  "accepted_at": "2026-01-28T21:39:34Z",
  ...
}
```

#### Error Response - No Shift (400 Bad Request)
```json
{
  "error": "barista must open a shift before accepting orders"
}
```

#### Error Response - Invalid Order State (400 Bad Request)
```json
{
  "error": "cannot accept order in current state"
}
```

## Unit Tests

### Test Coverage

Created `order_barista_shift_test.go` with 4 test cases:

#### 1. TestAcceptOrder_BaristaWithoutShift
**Purpose**: Verify barista cannot accept orders without shift

**Steps**:
1. Create QUEUED order
2. Try to accept without opening shift
3. Verify error: "barista must open a shift before accepting orders"
4. Verify order remains QUEUED
5. Verify barista_id not set

**Result**: ✅ PASS

---

#### 2. TestAcceptOrder_BaristaWithOpenShift
**Purpose**: Verify barista can accept orders with open shift

**Steps**:
1. Create QUEUED order
2. Open barista shift
3. Accept order
4. Verify status = IN_PROGRESS
5. Verify barista_id and barista_name set
6. Verify accepted_at timestamp set

**Result**: ✅ PASS

---

#### 3. TestAcceptOrder_BaristaWithClosedShift
**Purpose**: Verify barista cannot accept orders with closed shift

**Steps**:
1. Create QUEUED order
2. Open barista shift
3. Close the shift
4. Try to accept order
5. Verify error: "barista must open a shift before accepting orders"

**Result**: ✅ PASS

---

#### 4. TestAcceptOrder_MultipleBaristasDifferentShifts
**Purpose**: Verify multiple baristas with independent shifts

**Steps**:
1. Create 2 orders
2. Barista 1 opens shift
3. Barista 1 accepts order 1 (success)
4. Barista 2 tries to accept order 2 without shift (fail)
5. Barista 2 opens shift
6. Barista 2 accepts order 2 (success)
7. Verify each order assigned to correct barista

**Result**: ✅ PASS

---

### Running Tests

```bash
# Run all barista shift tests
go test -v ./application/services -run TestAcceptOrder

# Run specific test
go test -v ./application/services -run TestAcceptOrder_BaristaWithoutShift

# Run with coverage
go test -cover ./application/services -run TestAcceptOrder
```

### Test Output
```
=== RUN   TestAcceptOrder_BaristaWithoutShift
--- PASS: TestAcceptOrder_BaristaWithoutShift (0.00s)
=== RUN   TestAcceptOrder_BaristaWithOpenShift
--- PASS: TestAcceptOrder_BaristaWithOpenShift (0.00s)
=== RUN   TestAcceptOrder_BaristaWithClosedShift
--- PASS: TestAcceptOrder_BaristaWithClosedShift (0.00s)
=== RUN   TestAcceptOrder_MultipleBaristasDifferentShifts
--- PASS: TestAcceptOrder_MultipleBaristasDifferentShifts (0.00s)
PASS
ok      cafe-pos/backend/application/services   0.012s
```

## Frontend Integration

### User Flow

1. **Barista logs in**
   - Redirected to dashboard
   - Dashboard shows "Chưa mở ca" warning

2. **Barista opens shift**
   - Navigate to "Ca làm việc"
   - Select shift type (Morning/Afternoon/Evening)
   - Enter start_cash: 0 (barista doesn't handle cash)
   - Click "Mở ca"

3. **Barista views queue**
   - Navigate to "Barista" tab
   - See list of QUEUED orders
   - Each order shows: order number, items, time

4. **Barista accepts order**
   - Click "Nhận pha" button
   - If no shift: Show error "Bạn phải mở ca trước khi nhận order"
   - If shift open: Order moves to "Đang pha" tab

### Error Handling

```javascript
// frontend/src/stores/barista.js
async acceptOrder(orderId) {
  try {
    const response = await baristaService.acceptOrder(orderId)
    // Update local state
    this.removeFromQueue(orderId)
    this.addToWorking(response.data)
  } catch (error) {
    if (error.response?.data?.error === 'barista must open a shift before accepting orders') {
      // Show user-friendly message
      alert('Bạn phải mở ca trước khi nhận order')
      // Optionally redirect to shift page
      router.push('/shifts')
    } else {
      alert('Lỗi: ' + error.message)
    }
  }
}
```

## Testing Checklist

### Manual Testing

- [ ] Login as barista without shift
- [ ] Try to accept order → Should show error
- [ ] Open shift
- [ ] Try to accept order → Should succeed
- [ ] Close shift
- [ ] Try to accept another order → Should show error
- [ ] Reopen shift
- [ ] Accept order → Should succeed

### API Testing

Use `test-barista-shift-validation.sh`:

```bash
./test-barista-shift-validation.sh
```

Expected output:
```
=== Testing BR-13: Barista Shift Validation ===

1. Login as barista1...
✅ Logged in successfully

2. Check current shift (should not exist)...
   Response: {"error":"no open shift found"}

3. Get queued orders...
   Found 1 orders in queue

4. Try to accept order WITHOUT opening shift...
   ✅ PASS: Correctly rejected with error: barista must open a shift before accepting orders

5. Open barista shift...
   ✅ Shift opened successfully

6. Try to accept order WITH open shift...
   ✅ PASS: Order accepted successfully

=== All Tests Passed! ===
```

## Related Business Rules

- **BR-06**: Only Barista can move order to IN_PROGRESS
- **BR-07**: No modification after barista accepts (IN_PROGRESS)
- **BR-09**: READY indicates drink completed
- **BR-10**: User can have 1 OPEN shift per role
- **BR-13**: Barista must have open shift to accept orders ⭐ (This rule)

## Database Schema

### Shift Collection
```javascript
{
  _id: ObjectId,
  role_type: "barista",  // Important for validation
  user_id: ObjectId,     // Barista ID
  user_name: String,     // Barista name
  status: "OPEN",        // Must be OPEN
  type: "MORNING",
  start_cash: 0,
  started_at: ISODate,
  ...
}
```

### Order Collection
```javascript
{
  _id: ObjectId,
  status: "IN_PROGRESS",
  barista_id: ObjectId,    // Set when accepted
  barista_name: String,    // Set when accepted
  accepted_at: ISODate,    // Set when accepted
  ...
}
```

## Performance Considerations

### Query Optimization
- `FindOpenShiftByUser()` uses compound index: `{user_id: 1, role_type: 1, status: 1}`
- Query is fast: O(log n) with index
- No performance impact on order acceptance

### Caching (Future Enhancement)
- Cache barista's open shift in memory
- Invalidate on shift close
- Reduces DB queries

## Security Considerations

- ✅ JWT token validates barista identity
- ✅ Middleware ensures only barista role can access endpoint
- ✅ Service validates shift belongs to requesting barista
- ✅ Cannot accept orders for other baristas' shifts

## Monitoring & Logging

### Log Messages
```
INFO: Barista {barista_id} accepted order {order_id} in shift {shift_id}
WARN: Barista {barista_id} tried to accept order without shift
ERROR: Failed to find shift for barista {barista_id}
```

### Metrics to Track
- Number of orders accepted per shift
- Average time to accept order after queued
- Number of "no shift" errors per day
- Barista shift duration vs orders completed

## Troubleshooting

### Issue: Barista can't accept orders
**Symptoms**: Always get "must open shift" error

**Checks**:
1. Is shift actually open?
   ```bash
   curl -X GET http://localhost:8080/api/shifts/current \
     -H "Authorization: Bearer $TOKEN"
   ```

2. Is shift for barista role?
   ```javascript
   shift.role_type === "barista"
   ```

3. Is shift status OPEN?
   ```javascript
   shift.status === "OPEN"
   ```

4. Does user_id match?
   ```javascript
   shift.user_id === barista_id
   ```

### Issue: Backend not reflecting changes
**Solution**: Restart backend
```bash
./restart-backend.sh
```

## Future Enhancements

1. **Shift Reminder**: Notify barista to open shift when logging in
2. **Auto-Open Shift**: Option to auto-open shift on first order acceptance
3. **Shift Templates**: Pre-configured shift settings per barista
4. **Shift Handover**: Transfer orders when shift changes
5. **Break Time**: Pause shift without closing (barista on break)

## References

- Implementation: `backend/application/services/order_service.go`
- Tests: `backend/application/services/order_barista_shift_test.go`
- API: `backend/interfaces/http/order_handler.go`
- Frontend: `frontend/src/stores/barista.js`
- Documentation: `documents/SHIFT_BY_ROLE.md`
