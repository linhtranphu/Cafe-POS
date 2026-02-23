# Testing Guide: Cashier Shift View Redesign

## Overview
The Cashier Shift View has been completely redesigned with a new component `CashierShiftView.vue`. This guide helps you test all the new features.

## Prerequisites
- Login as a user with Cashier role
- Have some shifts created (both waiter and barista shifts)
- Have some payments in waiter shifts

## Test Scenarios

### 1. Tab Switching
- [ ] Open http://localhost:5173/#/shifts as Cashier
- [ ] Verify you see two tabs: "Ca phục vụ" and "Ca pha chế"
- [ ] Click on "Ca pha chế" tab - should turn blue
- [ ] Click back on "Ca phục vụ" tab - should turn blue
- [ ] Verify selected shift is cleared when switching tabs

### 2. Date Filter
- [ ] Select today's date in the date picker
- [ ] Verify shifts from today appear in the dropdown
- [ ] Select a different date
- [ ] Verify only shifts from that date appear
- [ ] Verify selected shift is cleared when date changes

### 3. Shift Selection - Waiter
- [ ] Select "Ca phục vụ" tab
- [ ] Choose a waiter shift from the dropdown
- [ ] Verify shift summary card appears with:
  - Yellow-orange gradient background
  - User name, shift type, time, status
  - 4 stat cards: Total orders, Revenue, Cash, Transfer
  - "Chốt ca" button (if shift is open)

### 4. Shift Selection - Barista
- [ ] Select "Ca pha chế" tab
- [ ] Choose a barista shift from the dropdown
- [ ] Verify shift summary card appears with:
  - Purple-indigo gradient background
  - User name, shift type, time, status
  - 2 stat cards: Total orders, In progress
  - "Chốt ca" button (if shift is open)

### 5. Payment List (Waiter Shifts Only)
- [ ] Select a waiter shift with payments
- [ ] Verify payment list appears below shift summary
- [ ] Each payment card should show:
  - Customer name, time, amount
  - Payment method badge (Cash/Transfer/QR)
  - Status badge
  - Three action buttons: Điều chỉnh, Báo lỗi, Khóa

### 6. Payment Actions
- [ ] Click "Điều chỉnh" on a payment
- [ ] Verify OverridePaymentModal opens
- [ ] Test override functionality
- [ ] Click "Báo lỗi" on a payment
- [ ] Verify DiscrepancyModal opens
- [ ] Test discrepancy reporting
- [ ] Click "🔒" button on a payment
- [ ] Verify confirmation dialog appears
- [ ] Test lock order functionality

### 7. Close Shift
- [ ] Select an open shift
- [ ] Click "Chốt ca" button
- [ ] Verify close shift modal appears
- [ ] Enter end cash amount
- [ ] Submit form
- [ ] Verify shift is closed and list refreshes

### 8. Loading States
- [ ] Verify loading spinner appears when:
  - Fetching shifts on mount
  - Loading shift details after selection
  - Loading payments for waiter shifts

### 9. Empty States
- [ ] Select a date with no shifts
- [ ] Verify "Chọn ca để xem chi tiết" message appears
- [ ] Select a waiter shift with no payments
- [ ] Verify "Chưa có thanh toán nào" message appears

### 10. Error Handling
- [ ] Test with network disconnected
- [ ] Verify error messages appear in alerts
- [ ] Test with invalid shift ID
- [ ] Verify graceful error handling

### 11. Mobile Responsiveness
- [ ] Test on mobile device or browser dev tools
- [ ] Verify all buttons are touch-friendly (min 44px height)
- [ ] Verify payment action buttons stack vertically on small screens
- [ ] Verify text wraps properly
- [ ] Verify gradients and colors look good

### 12. Role-Based Access
- [ ] Login as Waiter - verify original ShiftView appears
- [ ] Login as Barista - verify original ShiftView appears
- [ ] Login as Cashier - verify new CashierShiftView appears
- [ ] Login as Manager - verify new CashierShiftView appears

## Known Limitations
- Payment actions require OverridePaymentModal and DiscrepancyModal components to be implemented
- Close shift requires shiftStore.closeShift API to be working
- All API endpoints must be available and working

## Success Criteria
✅ All tab switching works smoothly
✅ Date filter correctly filters shifts
✅ Shift details load without errors
✅ Payment list displays correctly for waiter shifts
✅ All payment actions work as expected
✅ Close shift functionality works
✅ Loading states appear appropriately
✅ Error messages are user-friendly
✅ Mobile responsive design works well
✅ Only Cashier/Manager roles see the new view
