# Cashier Shift Data Refresh Fix

## Problem
After a handover is confirmed by the cashier, the fields (✅ Đã BG, 💼 Chưa BG, ⏳ Chờ) in the cashier dashboard were not updating, even with manual refresh or auto-refresh every 10 seconds.

## Root Cause
The `GetAllShifts` API endpoint was returning raw shift data from the database without calculating the latest transfer revenue and handover fields. 

When fetching individual shifts or the current shift, the backend calls `CalculateTransferRevenue()` to update:
- `current_cash`
- `remaining_cash`
- `handed_over_cash`
- `transfer_revenue`
- `remaining_transfer`
- `handed_over_transfer`

However, `GetAllShifts` was missing this calculation step, so the frontend received stale data.

## Solution
Updated `backend/interfaces/http/shift_handler.go` in the `GetAllShifts` handler to:

1. Fetch all shifts from database
2. For each **open waiter shift**, call `CalculateTransferRevenue()` to update handover fields
3. Refetch all shifts to get the updated values
4. Return the fresh data to frontend

```go
func (h *ShiftHandler) GetAllShifts(c *gin.Context) {
	shifts, err := h.shiftService.GetAllShifts(c.Request.Context())
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	// Calculate transfer revenue for all open waiter shifts
	for _, shift := range shifts {
		if shift.RoleType == order.RoleWaiter && shift.Status == order.ShiftOpen {
			_ = h.shiftService.CalculateTransferRevenue(c.Request.Context(), shift.ID)
		}
	}
	
	// Refetch shifts to get updated values
	shifts, err = h.shiftService.GetAllShifts(c.Request.Context())
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, shifts)
}
```

## Impact
- ✅ Cashier dashboard now shows real-time handover data
- ✅ Auto-refresh (every 10 seconds) works correctly
- ✅ Manual refresh updates all fields properly
- ✅ After handover confirmation, waiter shift data reflects immediately

## Files Changed
- `backend/interfaces/http/shift_handler.go` - Added transfer revenue calculation to `GetAllShifts`

## Testing
1. Open cashier shift at http://localhost:5173/#/cashier
2. Open waiter shift and create orders with cash/transfer payments
3. Waiter creates handover request
4. Cashier confirms handover
5. Verify that:
   - ✅ Đã BG increases
   - 💼 Chưa BG decreases
   - Per-waiter cards show updated amounts
   - Auto-refresh continues to show correct data

## Related Files
- `frontend/src/components/CashierShiftManager.vue` - Cashier dashboard component
- `frontend/src/stores/shift.js` - Shift store with `fetchAllShifts()`
- `backend/application/services/shift_service.go` - `CalculateTransferRevenue()` method
