# Transfer Revenue Display Fix

## Problem
Transfer revenue (tiền CK) was not updating in the ShiftView for waiter role. The backend had the `CalculateTransferRevenue()` method but it was never being called when fetching shift data.

## Root Cause
When the frontend fetched shift data via `GET /api/waiter/shifts/current` or `GET /api/waiter/shifts/:id`, the backend returned the shift data directly from the database without calculating the transfer revenue from orders.

## Solution
Modified the shift handler to automatically calculate transfer revenue when returning shift data for waiter shifts:

### Changes Made

**File: `backend/interfaces/http/shift_handler.go`**

1. **GetCurrentShift()** - Added automatic transfer revenue calculation:
```go
// Calculate transfer revenue for waiter shifts
if shift.RoleType == order.RoleWaiter {
    _ = h.shiftService.CalculateTransferRevenue(c.Request.Context(), shift.ID)
    // Refetch shift to get updated values
    shift, _ = h.shiftService.GetShift(c.Request.Context(), shift.ID)
}
```

2. **GetShift()** - Added automatic transfer revenue calculation:
```go
// Calculate transfer revenue for waiter shifts
if shift.RoleType == order.RoleWaiter {
    _ = h.shiftService.CalculateTransferRevenue(c.Request.Context(), shift.ID)
    // Refetch shift to get updated values
    shift, _ = h.shiftService.GetShift(c.Request.Context(), shift.ID)
}
```

## How It Works

1. When a waiter fetches their current shift, the handler:
   - Retrieves the shift from database
   - Checks if it's a waiter shift
   - Calls `CalculateTransferRevenue()` to recalculate from orders
   - Refetches the shift to get updated values
   - Returns the updated shift with correct transfer revenue

2. The `CalculateTransferRevenue()` method:
   - Fetches all orders for the shift
   - Separates cash and transfer payments
   - Updates `transfer_revenue`, `remaining_transfer`, and related fields
   - Saves the updated shift to database

## Benefits

- Transfer revenue is always up-to-date when viewing shift data
- No frontend changes required
- Automatic calculation on every shift fetch
- Works for both current shift and historical shifts

## Testing

1. Start a waiter shift
2. Create orders with transfer payment method (CK or QR)
3. View the shift - transfer revenue should display correctly
4. Create more orders - refresh to see updated transfer revenue

## Related Files

- `backend/interfaces/http/shift_handler.go` - Handler with automatic calculation
- `backend/application/services/shift_service.go` - Contains CalculateTransferRevenue method
- `frontend/src/views/ShiftView.vue` - Displays transfer revenue (no changes needed)
- `frontend/src/stores/shift.js` - Fetches shift data (no changes needed)
