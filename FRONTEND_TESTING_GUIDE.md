# Frontend Testing Guide - Cashier Fund Handover

## Overview

This guide provides comprehensive testing procedures for the Cashier Fund Handover feature frontend implementation.

## Test Environment Setup

### Prerequisites
1. Backend server running on `http://localhost:8080`
2. Frontend dev server running on `http://localhost:5173`
3. MongoDB running with test data
4. Valid JWT token for cashier role

### Test Data Setup

```bash
# 1. Start backend
cd backend
go run main.go

# 2. Start frontend (in another terminal)
cd frontend
npm run dev

# 3. Login as cashier to get JWT token
# Use browser DevTools > Application > Local Storage to get token
```

## Test Cases

### TC1: Dashboard - Managed Funds Display

**Objective**: Verify managed funds section displays correctly on dashboard

**Preconditions**:
- Cashier has an open shift
- Cashier has received handovers from waiters

**Steps**:
1. Login as cashier
2. Navigate to Cashier Dashboard
3. Observe "Tiền đang quản lý" section

**Expected Results**:
- ✅ Section displays with title "💰 Tiền đang quản lý"
- ✅ Shows "Tiền mặt" (cash) with green styling
- ✅ Shows "Tiền CK" (transfer) with blue styling
- ✅ Shows "Tổng cộng" (total) with orange gradient
- ✅ Shows warning message: "⚠️ Bạn chịu trách nhiệm trên số tiền này"
- ✅ Shows handover count
- ✅ All amounts formatted correctly (e.g., "1,500,000₫")

**Test Data**:
```javascript
{
  starting_float: 500000,
  received_cash: 1500000,
  received_transfer: 800000,
  total_managed_funds: 2300000,
  expected_cash: 2000000,
  handover_count: 5
}
```

---

### TC2: Dashboard - Pull to Refresh

**Objective**: Verify pull-to-refresh updates managed funds

**Steps**:
1. On Cashier Dashboard
2. Pull down from top of screen
3. Release to trigger refresh

**Expected Results**:
- ✅ Refresh indicator appears
- ✅ Managed funds data reloads
- ✅ Updated amounts display
- ✅ Refresh indicator disappears

---

### TC3: Closure Flow - Managed Funds Summary

**Objective**: Verify managed funds summary displays in closure flow

**Preconditions**:
- Cashier has open shift with received handovers
- All waiter shifts are closed

**Steps**:
1. Navigate to Cashier Dashboard
2. Click "Đóng ca" button
3. Observe managed funds summary card

**Expected Results**:
- ✅ Card displays with title "💰 Tiền đang quản lý"
- ✅ Shows starting float
- ✅ Shows received cash (green background)
- ✅ Shows received transfer (blue background)
- ✅ Shows expected cash calculation
- ✅ All amounts match dashboard values

---

### TC4: Closure Flow - Cash Counting

**Objective**: Verify cash counting and variance calculation

**Steps**:
1. In closure flow, view managed funds summary
2. Note the "expected_cash" value
3. Enter actual cash amount (try different scenarios):
   - Scenario A: Exact match (variance = 0)
   - Scenario B: Shortage (actual < expected)
   - Scenario C: Overage (actual > expected)

**Expected Results**:

**Scenario A (No Variance)**:
- ✅ Variance shows 0₫
- ✅ No variance documentation required
- ✅ Can proceed to confirmation

**Scenario B (Shortage)**:
- ✅ Variance shows negative amount (e.g., -5,000₫)
- ✅ Red styling for shortage
- ✅ Variance documentation form appears
- ✅ Reason dropdown required
- ✅ Notes field required (min 10 chars)

**Scenario C (Overage)**:
- ✅ Variance shows positive amount (e.g., +10,000₫)
- ✅ Green styling for overage
- ✅ Variance documentation form appears
- ✅ Reason and notes required

---

### TC5: Closure Flow - Variance Documentation

**Objective**: Verify variance documentation validation

**Preconditions**:
- Variance exists (actual ≠ expected)

**Steps**:
1. Try to proceed without selecting reason
2. Try to proceed with notes < 10 characters
3. Select valid reason and enter valid notes (≥10 chars)
4. Proceed to confirmation

**Expected Results**:
- ✅ Cannot proceed without reason selected
- ✅ Cannot proceed with short notes
- ✅ Error messages display clearly
- ✅ Can proceed with valid data
- ✅ Reason and notes preserved in confirmation

---

### TC6: Closure Flow - Confirmation Summary

**Objective**: Verify final confirmation displays all data correctly

**Steps**:
1. Complete all closure steps
2. View confirmation summary (Step 4)

**Expected Results**:
- ✅ Shows "📋 Tóm tắt bàn giao"
- ✅ Shows cash handover amount
- ✅ Shows transfer recorded amount
- ✅ Shows total amount
- ✅ Shows variance (if exists)
- ✅ Shows variance reason (if exists)
- ✅ "Xác nhận và đóng ca" button enabled
- ✅ All amounts match previous steps

---

### TC7: Closure Flow - Submit and Success

**Objective**: Verify successful closure with fund handover

**Steps**:
1. Complete all steps
2. Click "Xác nhận và đóng ca"
3. Wait for API response

**Expected Results**:
- ✅ Loading indicator appears
- ✅ API call to `/close-with-fund-handover` endpoint
- ✅ Success message displays
- ✅ Redirects to dashboard
- ✅ Shift status updated to CLOSED
- ✅ Fund handover record created in database

**Verify in Database**:
```javascript
// Check fund handover created
db.fund_handovers.findOne({ cashier_shift_id: ObjectId("...") })

// Check shift closed
db.cashier_shifts.findOne({ _id: ObjectId("..."), status: "CLOSED" })
```

---

### TC8: Closure Flow - Error Handling

**Objective**: Verify error handling for various failure scenarios

**Test Scenarios**:

**A. Network Error**:
1. Disconnect network
2. Try to close shift
3. Expected: Error message displays, can retry

**B. Waiter Shifts Still Open**:
1. Leave waiter shift open
2. Try to close cashier shift
3. Expected: Error message about open waiter shifts

**C. Invalid Data**:
1. Enter negative actual cash
2. Try to proceed
3. Expected: Validation error

**D. Session Expired**:
1. Clear JWT token
2. Try to close shift
3. Expected: Redirect to login

---

### TC9: Mobile Responsiveness

**Objective**: Verify UI works on mobile devices

**Test Devices**:
- iPhone (Safari)
- Android (Chrome)
- iPad (Safari)

**Steps**:
1. Open on mobile device
2. Test all closure flow steps
3. Test pull-to-refresh
4. Test button interactions

**Expected Results**:
- ✅ All text readable (no overflow)
- ✅ Buttons easily tappable (min 44x44px)
- ✅ Cards display properly
- ✅ No horizontal scrolling
- ✅ Safe area insets respected
- ✅ Keyboard doesn't cover inputs
- ✅ Pull-to-refresh works smoothly

---

### TC10: Back Navigation

**Objective**: Verify back button behavior

**Steps**:
1. Start closure flow
2. Click "Quay lại" at various steps
3. Observe behavior

**Expected Results**:
- ✅ Returns to dashboard
- ✅ No data submitted to API
- ✅ Shift remains in OPEN status
- ✅ Can restart closure flow

---

## API Integration Tests

### Test API Calls

Use browser DevTools > Network tab to verify:

**1. Get Managed Funds**
```
GET /api/v1/cashier-shifts/{id}/managed-funds
Status: 200 OK
Response: {
  cashier_shift_id: "...",
  starting_float: 500000,
  received_cash: 1500000,
  received_transfer: 800000,
  total_managed_funds: 2300000,
  expected_cash: 2000000,
  handover_count: 5
}
```

**2. Close with Fund Handover**
```
POST /api/v1/cashier-shifts/{id}/close-with-fund-handover
Request: {
  actual_cash: 1995000,
  variance_reason: "COUNTING_ERROR",
  variance_notes: "Đếm nhầm tờ 50k thành 100k"
}
Status: 200 OK
Response: {
  shift: { ... },
  fund_handover: { ... }
}
```

---

## Performance Tests

### Load Time
- ✅ Dashboard loads < 2 seconds
- ✅ Managed funds API call < 500ms
- ✅ Closure submission < 1 second

### Smooth Scrolling
- ✅ No lag when scrolling
- ✅ Pull-to-refresh smooth
- ✅ Animations smooth (60fps)

---

## Accessibility Tests

### Keyboard Navigation
- ✅ Can tab through all inputs
- ✅ Enter key submits forms
- ✅ Escape key closes modals

### Screen Reader
- ✅ All amounts announced correctly
- ✅ Error messages announced
- ✅ Button labels clear

### Color Contrast
- ✅ Text readable on all backgrounds
- ✅ WCAG AA compliance

---

## Browser Compatibility

Test on:
- ✅ Chrome (latest)
- ✅ Safari (latest)
- ✅ Firefox (latest)
- ✅ Edge (latest)
- ✅ Mobile Safari (iOS)
- ✅ Mobile Chrome (Android)

---

## Regression Tests

After any changes, verify:
- ✅ Existing closure flow still works
- ✅ Dashboard displays correctly
- ✅ No console errors
- ✅ No broken styles
- ✅ API calls successful

---

## Test Checklist

Use this checklist for each test session:

```
Dashboard Tests:
[ ] TC1: Managed funds display
[ ] TC2: Pull to refresh

Closure Flow Tests:
[ ] TC3: Managed funds summary
[ ] TC4: Cash counting (3 scenarios)
[ ] TC5: Variance documentation
[ ] TC6: Confirmation summary
[ ] TC7: Submit and success
[ ] TC8: Error handling (4 scenarios)

UI/UX Tests:
[ ] TC9: Mobile responsiveness
[ ] TC10: Back navigation

Integration Tests:
[ ] API calls verified
[ ] Database records verified

Performance Tests:
[ ] Load times acceptable
[ ] Smooth scrolling

Accessibility Tests:
[ ] Keyboard navigation
[ ] Screen reader
[ ] Color contrast

Browser Tests:
[ ] Chrome
[ ] Safari
[ ] Firefox
[ ] Mobile browsers
```

---

## Bug Reporting Template

When reporting bugs, use this format:

```markdown
**Bug ID**: FH-XXX
**Severity**: Critical / High / Medium / Low
**Test Case**: TCXX
**Environment**: Browser, OS, Device

**Steps to Reproduce**:
1. 
2. 
3. 

**Expected Result**:


**Actual Result**:


**Screenshots**:
[Attach screenshots]

**Console Errors**:
[Paste console errors]

**Network Logs**:
[Paste relevant API calls]
```

---

## Next Steps

After completing frontend testing:
1. Document all bugs found
2. Fix critical and high priority bugs
3. Retest fixed bugs
4. Perform final regression test
5. Get stakeholder approval
6. Proceed to deployment
