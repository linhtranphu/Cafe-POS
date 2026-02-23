# Frontend-Driven Cashier Shift Closure

## Overview

This document describes the frontend-driven approach for cashier shift closure, where all data is collected on the client side and submitted in a single transaction.

## Problem Statement

The previous implementation used a step-by-step approach where each step (initiate, record cash, document variance, close) was saved immediately to the database. This meant:

- If user clicked "Back" after Step 2, the data was already saved
- User expected "database transaction" behavior: if not completed, nothing should be saved
- Rollback logic was complex and error-prone

## Solution: Frontend-Driven Approach (Option A)

### How It Works

1. **Frontend collects all data locally** (in Vue component state)
2. **No API calls until user clicks "Hoàn tất đóng ca"**
3. **One API call with all data** → Backend processes in ONE transaction (~100ms)
4. **If user clicks "Back"** → Local data discarded, no API call, nothing saved

### Benefits

✅ True "all-or-nothing" behavior from user perspective
✅ Simple rollback: just discard local state
✅ Short transaction duration (~100ms)
✅ No complex state management on backend
✅ Better UX: user can review all data before committing

## Implementation

### Backend

#### New Service Method: `CompleteClosure`

```go
// backend/application/services/cashier_shift_service.go

func (s *CashierShiftService) CompleteClosure(
	ctx context.Context,
	shiftID primitive.ObjectID,
	actualCash float64,
	varianceReason *cashier.VarianceReason,
	varianceNotes *string,
	userID, deviceID string,
) (*cashier.CashierShift, error)
```

**What it does:**
1. Starts MongoDB transaction
2. Validates shift is in OPEN status
3. Checks all waiter shifts are closed
4. Initiates closure
5. Records actual cash
6. Documents variance (if needed)
7. Closes shift
8. Saves to database

**All in ONE transaction** - if any step fails, everything rolls back.

#### New Handler: `CompleteClosure`

```go
// backend/interfaces/http/cashier_shift_closure_handler.go

func (h *CashierShiftClosureHandler) CompleteClosure(c *gin.Context)
```

**Request payload:**
```json
{
  "actual_cash": 500000,
  "variance_reason": "COUNTING_ERROR",  // optional, required if variance exists
  "variance_notes": "Đếm nhầm tiền lẻ"  // optional, required if variance exists
}
```

#### New Route

```go
// backend/main.go
cashierShifts.POST("/:id/complete-closure", cashierShiftClosureHandler.CompleteClosure)
```

### Frontend

#### New View: `CashierShiftClosureV2.vue`

**Key features:**
- All data stored in local `closureData` ref
- No API calls until "Hoàn tất đóng ca" clicked
- "Quay lại" button discards local data
- Variance calculated locally
- Variance documentation shown conditionally

**Data structure:**
```javascript
const closureData = ref({
  actualCash: null,
  variance: null,
  varianceReason: '',
  varianceNotes: ''
})
```

**Workflow:**
1. Check waiter shifts (API call)
2. User enters actual cash → variance calculated locally
3. If variance exists → show variance documentation form
4. User clicks "Hoàn tất đóng ca" → ONE API call with all data
5. Backend processes in transaction → shift closed

#### Service Method: `completeClosure`

```javascript
// frontend/src/services/cashierShift.js

async completeClosure(shiftId, data) {
  const response = await api.post(`/cashier-shifts/${shiftId}/complete-closure`, data)
  return response.data
}
```

#### Router Update

```javascript
// frontend/src/router/index.js
import CashierShiftClosureV2 from '../views/CashierShiftClosureV2.vue'

{
  path: '/cashier/shift-closure/:id',
  name: 'CashierShiftClosure',
  component: CashierShiftClosureV2,
  meta: { requiresAuth: true, requiresCashier: true }
}
```

## User Experience

### Scenario 1: No Variance

1. User clicks "Đóng ca" from dashboard
2. System checks waiter shifts
3. User enters actual cash (same as system cash)
4. Variance shows: 0₫
5. User clicks "Hoàn tất đóng ca"
6. ✅ Shift closed immediately

### Scenario 2: With Variance

1. User clicks "Đóng ca" from dashboard
2. System checks waiter shifts
3. User enters actual cash (different from system cash)
4. Variance shows: +50,000₫ or -50,000₫
5. System shows variance documentation form
6. User selects reason and enters notes
7. User clicks "Hoàn tất đóng ca"
8. ✅ Shift closed with variance documented

### Scenario 3: User Changes Mind

1. User clicks "Đóng ca" from dashboard
2. System checks waiter shifts
3. User enters actual cash
4. User realizes they made a mistake
5. User clicks "Quay lại"
6. Confirmation: "Tất cả dữ liệu đã nhập sẽ bị mất"
7. User confirms
8. ✅ Returns to dashboard, NO data saved

## Transaction Behavior

### What Happens in the Transaction

```
START TRANSACTION
  1. Get shift from database
  2. Validate shift status = OPEN
  3. Check waiter shifts are closed
  4. shift.InitiateClosure()
  5. shift.RecordActualCash()
  6. IF variance exists:
       shift.DocumentVariance()
  7. shift.Close()
  8. Save shift to database
COMMIT TRANSACTION
```

**Duration:** ~100ms (very short, safe for transactions)

### Rollback Scenarios

**Automatic rollback if:**
- Shift not found
- Shift not in OPEN status
- Waiter shifts still open
- Validation fails at any step
- Database save fails

**Manual rollback:**
- User clicks "Quay lại" → Local data discarded, no API call

## Testing

### Test Script: `test-complete-closure.sh`

```bash
./test-complete-closure.sh
```

**What it tests:**
1. Login as cashier
2. Get current shift
3. Check waiter shifts
4. Complete closure (no variance)
5. Verify shift is closed

### Manual Testing

1. **Test no variance:**
   - Start cashier shift
   - Enter actual cash = system cash
   - Click "Hoàn tất đóng ca"
   - Verify shift closed

2. **Test with variance:**
   - Start cashier shift
   - Enter actual cash ≠ system cash
   - Fill variance documentation
   - Click "Hoàn tất đóng ca"
   - Verify shift closed with variance

3. **Test cancel:**
   - Start cashier shift
   - Enter actual cash
   - Click "Quay lại"
   - Confirm discard
   - Verify no data saved

4. **Test validation:**
   - Try to close with waiter shifts open
   - Verify error message
   - Close waiter shifts
   - Try again → should work

## Comparison with Previous Approach

### Previous: Step-by-Step (Per-Step Transactions)

```
User Action          API Call                    DB State
-----------          --------                    --------
Click "Đóng ca"   → POST /initiate-closure   → status = CLOSURE_INITIATED ✅ SAVED
Enter cash        → POST /record-actual-cash → actual_cash = 500000 ✅ SAVED
Click "Back"      → POST /cancel-closure     → Rollback to OPEN ✅ SAVED
```

**Problem:** Data saved at each step, complex rollback logic

### New: Frontend-Driven (Single Transaction)

```
User Action          API Call                    DB State
-----------          --------                    --------
Click "Đóng ca"   → (none)                    → (no change)
Enter cash        → (none)                    → (no change)
Click "Back"      → (none)                    → (no change)
Click "Hoàn tất"  → POST /complete-closure   → status = CLOSED ✅ SAVED
```

**Benefit:** No data saved until completion, simple discard logic

## Migration Notes

### Old Endpoints (Still Available)

These endpoints are still available for backward compatibility:
- `POST /:id/initiate-closure`
- `POST /:id/record-actual-cash`
- `POST /:id/document-variance`
- `POST /:id/close`
- `POST /:id/cancel-closure`

### New Endpoint

- `POST /:id/complete-closure` (recommended)

### Deprecation Plan

1. ✅ Implement new frontend-driven approach
2. ✅ Update router to use V2 view
3. Test thoroughly
4. Monitor usage
5. Consider deprecating old endpoints after 1-2 months

## Files Changed

### Backend
- `backend/application/services/cashier_shift_service.go` (added `CompleteClosure`)
- `backend/interfaces/http/cashier_shift_closure_handler.go` (added `CompleteClosure` handler)
- `backend/main.go` (added route)

### Frontend
- `frontend/src/views/CashierShiftClosureV2.vue` (NEW FILE)
- `frontend/src/services/cashierShift.js` (added `completeClosure` method)
- `frontend/src/router/index.js` (updated to use V2 view)

### Documentation
- `TRANSACTION_PATTERNS_EXPLAINED.md` (explains 3 patterns)
- `FRONTEND_DRIVEN_CLOSURE.md` (this file)

### Testing
- `test-complete-closure.sh` (test script)

## Conclusion

The frontend-driven approach provides true "all-or-nothing" behavior from the user's perspective while maintaining short transaction durations and simple rollback logic. This is the recommended approach for workflows where users expect to review all data before committing.

## Next Steps

1. ✅ Implementation complete
2. Test the new workflow end-to-end
3. Verify transaction behavior
4. Monitor for any issues
5. Consider deprecating old step-by-step endpoints
