# Manual Testing Checklist - Cashier Fund Handover

## Pre-Testing Setup

- [ ] Backend server running
- [ ] Frontend dev server running
- [ ] MongoDB connected
- [ ] Test cashier account created
- [ ] Browser DevTools open (Console + Network tabs)

## Test Session Information

**Date**: _______________
**Tester**: _______________
**Browser**: _______________
**Device**: _______________
**Environment**: Development / Staging / Production

---

## Part 1: Dashboard Testing

### 1.1 Initial Load
- [ ] Dashboard loads without errors
- [ ] No console errors
- [ ] All API calls successful (check Network tab)
- [ ] Loading states display correctly

### 1.2 Managed Funds Section
- [ ] Section title "💰 Tiền đang quản lý" displays
- [ ] Cash amount displays with 💵 icon
- [ ] Cash has green background (bg-green-50)
- [ ] Transfer amount displays with 💳 icon
- [ ] Transfer has blue background (bg-blue-50)
- [ ] Total displays with orange gradient
- [ ] Warning message displays: "⚠️ Bạn chịu trách nhiệm..."
- [ ] Handover count displays correctly
- [ ] All amounts formatted with commas and ₫ symbol

**Test Data Recorded**:
- Starting Float: _______________
- Received Cash: _______________
- Received Transfer: _______________
- Total: _______________
- Handover Count: _______________

### 1.3 Pull to Refresh
- [ ] Pull down gesture works
- [ ] Refresh indicator appears
- [ ] Data reloads
- [ ] Managed funds update
- [ ] Indicator disappears after load

### 1.4 Refresh Button
- [ ] Refresh button (🔄) visible
- [ ] Button disabled during loading
- [ ] Spin animation during loading
- [ ] Data refreshes on click

---

## Part 2: Closure Flow Testing

### 2.1 Navigation to Closure
- [ ] "Đóng ca" button visible on dashboard
- [ ] Button disabled if waiter shifts open
- [ ] Clicking button navigates to closure view
- [ ] URL changes to /cashier-shift-closure

### 2.2 Closure View - Initial Load
- [ ] Header displays "🔒 Đóng ca thu ngân"
- [ ] Back button (← Quay lại) visible
- [ ] Shift info card displays
- [ ] Managed funds summary card displays
- [ ] No console errors

### 2.3 Managed Funds Summary Card
- [ ] Card title "💰 Tiền đang quản lý" displays
- [ ] Starting float row displays
- [ ] "Nhận từ waiter" section displays
- [ ] Cash row with green background
- [ ] Transfer row with blue background
- [ ] Expected cash calculation displays
- [ ] All amounts match dashboard values

**Verify Calculations**:
- Expected Cash = Starting Float + Received Cash
- Calculated: _______________ = _______________ + _______________
- [ ] Calculation correct

### 2.4 Cash Counting - No Variance
- [ ] Cash input field visible
- [ ] Can enter amount
- [ ] Enter exact expected cash amount
- [ ] Variance shows 0₫
- [ ] No variance documentation form appears
- [ ] Can proceed to confirmation

**Test Data**:
- Expected Cash: _______________
- Actual Cash Entered: _______________
- Variance: 0₫

### 2.5 Cash Counting - Shortage
- [ ] Enter amount less than expected
- [ ] Variance calculates automatically
- [ ] Variance shows negative amount
- [ ] Red styling for shortage
- [ ] Variance documentation form appears
- [ ] Reason dropdown displays
- [ ] Notes textarea displays
- [ ] Cannot proceed without reason
- [ ] Cannot proceed with notes < 10 chars
- [ ] Error messages display

**Test Data**:
- Expected Cash: _______________
- Actual Cash Entered: _______________
- Variance: _______________ (negative)

### 2.6 Cash Counting - Overage
- [ ] Enter amount more than expected
- [ ] Variance calculates automatically
- [ ] Variance shows positive amount
- [ ] Green styling for overage
- [ ] Variance documentation form appears
- [ ] Reason and notes required

**Test Data**:
- Expected Cash: _______________
- Actual Cash Entered: _______________
- Variance: _______________ (positive)

### 2.7 Variance Documentation
- [ ] Reason dropdown has options:
  - [ ] COUNTING_ERROR
  - [ ] CUSTOMER_DISPUTE
  - [ ] SYSTEM_ERROR
  - [ ] OTHER
- [ ] Can select reason
- [ ] Notes field accepts text
- [ ] Notes field shows character count
- [ ] Validation: notes must be ≥ 10 characters
- [ ] Error message if too short
- [ ] Can proceed with valid data

**Test Data**:
- Reason Selected: _______________
- Notes Entered: _______________
- Notes Length: _______________

### 2.8 Confirmation Summary
- [ ] Summary card displays "📋 Tóm tắt bàn giao"
- [ ] Cash handover amount displays
- [ ] Transfer recorded amount displays
- [ ] Total amount displays
- [ ] Variance displays (if exists)
- [ ] Variance reason displays (if exists)
- [ ] All amounts match previous steps
- [ ] "Xác nhận và đóng ca" button enabled

**Verify Summary**:
- Cash Handover: _______________
- Transfer Recorded: _______________
- Total: _______________
- Variance: _______________

### 2.9 Submit Closure
- [ ] Click "Xác nhận và đóng ca"
- [ ] Loading indicator appears
- [ ] Button disabled during submit
- [ ] Check Network tab for API call:
  - [ ] POST to /close-with-fund-handover
  - [ ] Request body correct
  - [ ] Status 200 OK
  - [ ] Response contains shift and fund_handover
- [ ] Success message displays
- [ ] Redirects to dashboard
- [ ] Shift status updated to CLOSED

**API Response Verification**:
- Shift ID: _______________
- Fund Handover ID: _______________
- Status Code: _______________

### 2.10 Database Verification
```javascript
// Run in MongoDB shell
db.fund_handovers.findOne({ cashier_shift_id: ObjectId("_______________") })
db.cashier_shifts.findOne({ _id: ObjectId("_______________") })
```

- [ ] Fund handover record exists
- [ ] Cash amount correct
- [ ] Transfer amount correct
- [ ] Variance amount correct
- [ ] Variance reason correct (if exists)
- [ ] Variance notes correct (if exists)
- [ ] Cashier shift status = CLOSED
- [ ] Timestamps correct

---

## Part 3: Error Handling

### 3.1 Network Errors
- [ ] Disconnect network
- [ ] Try to load managed funds
- [ ] Error message displays
- [ ] Can retry after reconnecting

### 3.2 Waiter Shifts Open
- [ ] Leave waiter shift open
- [ ] Try to close cashier shift
- [ ] Error message about open waiter shifts
- [ ] Cannot proceed

### 3.3 Invalid Input
- [ ] Try negative actual cash
- [ ] Validation error displays
- [ ] Try empty variance notes
- [ ] Validation error displays

### 3.4 Session Expired
- [ ] Clear JWT token from localStorage
- [ ] Try to access dashboard
- [ ] Redirects to login
- [ ] No console errors

---

## Part 4: Mobile Testing

### 4.1 iPhone Testing
**Device**: _______________
**iOS Version**: _______________

- [ ] Dashboard loads correctly
- [ ] Managed funds section displays
- [ ] Text readable (no overflow)
- [ ] Buttons tappable (min 44x44px)
- [ ] Pull-to-refresh works
- [ ] Closure flow works
- [ ] Keyboard doesn't cover inputs
- [ ] Safe area insets respected
- [ ] No horizontal scrolling

### 4.2 Android Testing
**Device**: _______________
**Android Version**: _______________

- [ ] Dashboard loads correctly
- [ ] Managed funds section displays
- [ ] Text readable
- [ ] Buttons tappable
- [ ] Pull-to-refresh works
- [ ] Closure flow works
- [ ] Keyboard behavior correct
- [ ] No horizontal scrolling

### 4.3 Tablet Testing
**Device**: _______________

- [ ] Layout adapts to larger screen
- [ ] All features work
- [ ] No UI issues

---

## Part 5: Browser Compatibility

### 5.1 Chrome
**Version**: _______________
- [ ] All features work
- [ ] No console errors
- [ ] Styling correct

### 5.2 Safari
**Version**: _______________
- [ ] All features work
- [ ] No console errors
- [ ] Styling correct

### 5.3 Firefox
**Version**: _______________
- [ ] All features work
- [ ] No console errors
- [ ] Styling correct

### 5.4 Edge
**Version**: _______________
- [ ] All features work
- [ ] No console errors
- [ ] Styling correct

---

## Part 6: Performance

### 6.1 Load Times
- Dashboard initial load: _______________ ms
- Managed funds API call: _______________ ms
- Closure submission: _______________ ms

**Targets**:
- [ ] Dashboard < 2 seconds
- [ ] API calls < 500ms
- [ ] Submission < 1 second

### 6.2 Smooth Scrolling
- [ ] No lag when scrolling
- [ ] Pull-to-refresh smooth
- [ ] Animations smooth (60fps)

---

## Part 7: Accessibility

### 7.1 Keyboard Navigation
- [ ] Can tab through all inputs
- [ ] Enter key submits forms
- [ ] Escape key closes modals
- [ ] Focus indicators visible

### 7.2 Screen Reader
- [ ] Test with VoiceOver (iOS) or TalkBack (Android)
- [ ] All amounts announced correctly
- [ ] Error messages announced
- [ ] Button labels clear

### 7.3 Color Contrast
- [ ] Text readable on all backgrounds
- [ ] Meets WCAG AA standards
- [ ] Color not sole indicator of meaning

---

## Part 8: Edge Cases

### 8.1 Zero Handovers
- [ ] Dashboard displays with 0 handovers
- [ ] Managed funds shows 0 for received amounts
- [ ] Closure flow works

### 8.2 Large Amounts
- [ ] Test with amounts > 100,000,000₫
- [ ] Formatting correct
- [ ] No overflow issues

### 8.3 Multiple Rapid Clicks
- [ ] Click submit button rapidly
- [ ] Only one API call made
- [ ] Button disabled after first click

---

## Issues Found

| # | Severity | Description | Steps to Reproduce | Screenshot |
|---|----------|-------------|-------------------|------------|
| 1 |          |             |                   |            |
| 2 |          |             |                   |            |
| 3 |          |             |                   |            |

**Severity Levels**: Critical / High / Medium / Low

---

## Sign-Off

**Tester Signature**: _______________
**Date**: _______________

**Status**: ☐ Passed  ☐ Failed  ☐ Passed with Issues

**Notes**:
_______________________________________________
_______________________________________________
_______________________________________________

---

## Next Steps

After completing this checklist:
1. [ ] Document all issues found
2. [ ] Prioritize issues by severity
3. [ ] Create bug tickets
4. [ ] Fix critical and high priority issues
5. [ ] Retest fixed issues
6. [ ] Get stakeholder approval
7. [ ] Proceed to deployment
