# 🧪 Cash Handover E2E Test Cases

## 📋 Test Overview

This document outlines end-to-end test scenarios for the cash handover feature.

---

## 🎯 Test Scenarios

### Scenario 1: Partial Handover - Happy Path

**Actors:** Waiter, Cashier

**Preconditions:**
- Waiter has an open shift with cash > 0
- Cashier has an open shift

**Steps:**
1. **Waiter:** Login and navigate to Shifts page
2. **Waiter:** Verify cash status cards are visible
3. **Waiter:** Click "💰 Bàn giao một phần" button
4. **Waiter:** Enter amount (e.g., 200,000 VND) ≤ remaining cash
5. **Waiter:** Add optional note
6. **Waiter:** Click "Bàn giao" button
7. **Waiter:** Verify pending handover banner appears
8. **Waiter:** Verify handover buttons are disabled
9. **Cashier:** Login and navigate to Dashboard
10. **Cashier:** Verify notification banner shows "1 yêu cầu bàn giao"
11. **Cashier:** Click "Xem ngay" or navigate to Handovers page
12. **Cashier:** Verify handover appears in pending list
13. **Cashier:** Click "✅ Xác nhận" button
14. **Cashier:** Enter actual amount (same as declared)
15. **Cashier:** Add optional note
16. **Cashier:** Click "Xác nhận" button
17. **Cashier:** Verify success message
18. **Cashier:** Verify handover moves to "Today's" section
19. **Waiter:** Refresh Shifts page
20. **Waiter:** Verify pending banner disappears
21. **Waiter:** Verify cash amounts updated correctly
22. **Waiter:** Verify handover appears in history with "Đã xác nhận" status

**Expected Results:**
- ✅ Handover created successfully
- ✅ Pending state visible to both parties
- ✅ Confirmation updates cash amounts
- ✅ History shows completed handover
- ✅ No errors or warnings

---

### Scenario 2: Handover with Discrepancy

**Actors:** Waiter, Cashier

**Preconditions:**
- Waiter has an open shift with cash > 0
- Cashier has an open shift

**Steps:**
1. **Waiter:** Create handover for 150,000 VND
2. **Cashier:** Navigate to pending handovers
3. **Cashier:** Click "✅ Xác nhận" on the handover
4. **Cashier:** Enter actual amount: 145,000 VND (5k short)
5. **Cashier:** Select discrepancy reason: "COUNTING_ERROR"
6. **Cashier:** Select responsibility: "WAITER"
7. **Cashier:** Add note explaining discrepancy
8. **Cashier:** Click "Xác nhận"
9. **Cashier:** Verify success message
10. **Waiter:** Check handover history
11. **Waiter:** Verify discrepancy is shown: "-5,000₫"
12. **Waiter:** Verify cashier note is visible

**Expected Results:**
- ✅ Discrepancy recorded correctly
- ✅ Actual amount used for cash updates
- ✅ Discrepancy visible in history
- ✅ Reason and responsibility tracked

---

### Scenario 3: Large Discrepancy Requiring Approval

**Actors:** Waiter, Cashier, Manager

**Preconditions:**
- Waiter has an open shift
- Cashier has an open shift
- Manager is logged in
- Discrepancy threshold is 100,000 VND

**Steps:**
1. **Waiter:** Create handover for 500,000 VND
2. **Cashier:** Confirm with actual amount: 380,000 VND (120k short)
3. **Cashier:** Fill in discrepancy details
4. **Cashier:** Click "Xác nhận"
5. **Cashier:** Verify warning about manager approval
6. **Cashier:** Verify handover status is "DISCREPANCY"
7. **Manager:** Navigate to pending approvals
8. **Manager:** Review handover details and discrepancy
9. **Manager:** Click "Approve" or "Reject"
10. **Manager:** Add approval note
11. **Manager:** Submit decision
12. **Waiter:** Check handover status
13. **Waiter:** Verify manager decision is reflected

**Expected Results:**
- ✅ Large discrepancy triggers approval workflow
- ✅ Cash amounts NOT updated until approval
- ✅ Manager can review and approve/reject
- ✅ Final status reflects manager decision

---

### Scenario 4: Handover and End Shift

**Actors:** Waiter, Cashier

**Preconditions:**
- Waiter has an open shift with remaining cash
- Cashier has an open shift

**Steps:**
1. **Waiter:** Click "🏁 Bàn giao và đóng ca" button
2. **Waiter:** Read warning notice
3. **Waiter:** Verify "Tiền sẽ bàn giao" shows all remaining cash
4. **Waiter:** Enter end cash (usually 0)
5. **Waiter:** Add optional note
6. **Waiter:** Click "Bàn giao và đóng ca"
7. **Waiter:** Verify pending status message
8. **Waiter:** Verify cannot perform other actions
9. **Cashier:** Navigate to pending handovers
10. **Cashier:** Verify handover type shows "Toàn bộ + Đóng ca"
11. **Cashier:** Confirm handover
12. **Waiter:** Verify shift automatically closes
13. **Waiter:** Verify shift appears in history as "Đã đóng"
14. **Waiter:** Verify cannot reopen or modify closed shift

**Expected Results:**
- ✅ All remaining cash handed over
- ✅ Shift closes automatically after confirmation
- ✅ Orders locked after shift closure
- ✅ Shift appears in history
- ✅ Cannot modify closed shift

---

### Scenario 5: Reject Handover

**Actors:** Waiter, Cashier

**Preconditions:**
- Waiter has created a handover request

**Steps:**
1. **Cashier:** Navigate to pending handovers
2. **Cashier:** Click "❌ Từ chối" on a handover
3. **Cashier:** Enter rejection reason (required)
4. **Cashier:** Click "Từ chối"
5. **Cashier:** Verify success message
6. **Waiter:** Check handover history
7. **Waiter:** Verify status shows "Đã từ chối"
8. **Waiter:** Read cashier's rejection reason
9. **Waiter:** Verify cash amounts unchanged
10. **Waiter:** Verify can create new handover

**Expected Results:**
- ✅ Handover rejected with reason
- ✅ Cash amounts not updated
- ✅ Waiter can see rejection reason
- ✅ Waiter can create new handover

---

### Scenario 6: Cancel Handover (Waiter)

**Actors:** Waiter

**Preconditions:**
- Waiter has a pending handover

**Steps:**
1. **Waiter:** Navigate to Shifts page
2. **Waiter:** Verify pending handover banner visible
3. **Waiter:** Click "Hủy" button on banner
4. **Waiter:** Confirm cancellation in dialog
5. **Waiter:** Verify banner disappears
6. **Waiter:** Verify handover buttons re-enabled
7. **Waiter:** Check handover history
8. **Waiter:** Verify cancelled handover not in history (or marked as cancelled)

**Expected Results:**
- ✅ Handover cancelled successfully
- ✅ Pending state cleared
- ✅ Can create new handover
- ✅ Cashier no longer sees cancelled handover

---

### Scenario 7: Quick Confirm from Dashboard

**Actors:** Cashier

**Preconditions:**
- Multiple pending handovers exist

**Steps:**
1. **Cashier:** Login and view Dashboard
2. **Cashier:** Verify "⚡ Bàn giao nhanh" section visible
3. **Cashier:** See up to 3 pending handovers
4. **Cashier:** Click "✅" button on first handover
5. **Cashier:** Verify success message
6. **Cashier:** Verify handover removed from quick actions
7. **Cashier:** Verify count updated
8. **Cashier:** Click "❌" button on another handover
9. **Cashier:** Verify rejection processed

**Expected Results:**
- ✅ Quick actions work without opening modal
- ✅ Assumes declared amount = actual amount
- ✅ Fast processing for exact matches
- ✅ Dashboard updates immediately

---

### Scenario 8: Multiple Handovers in One Shift

**Actors:** Waiter, Cashier

**Preconditions:**
- Waiter has open shift with sufficient cash

**Steps:**
1. **Waiter:** Create 1st partial handover (100k)
2. **Cashier:** Confirm 1st handover
3. **Waiter:** Verify cash updated
4. **Waiter:** Create 2nd partial handover (150k)
5. **Cashier:** Confirm 2nd handover
6. **Waiter:** Verify cash updated again
7. **Waiter:** Create 3rd partial handover (200k)
8. **Cashier:** Confirm 3rd handover
9. **Waiter:** Check handover history
10. **Waiter:** Verify all 3 handovers listed
11. **Waiter:** Verify total handed over = 450k

**Expected Results:**
- ✅ Multiple handovers allowed per shift
- ✅ Cash amounts cumulative
- ✅ History shows all handovers
- ✅ Remaining cash decreases correctly

---

### Scenario 9: No Cashier Shift Open

**Actors:** Waiter

**Preconditions:**
- Waiter has open shift
- NO cashier shift is open

**Steps:**
1. **Waiter:** Click "💰 Bàn giao một phần"
2. **Waiter:** Enter amount and note
3. **Waiter:** Click "Bàn giao"
4. **Waiter:** Verify error message: "no active cashier shift found"
5. **Waiter:** Verify handover not created
6. **Waiter:** Verify cash amounts unchanged

**Expected Results:**
- ✅ Error message displayed
- ✅ Handover not created
- ✅ User informed to wait for cashier

---

### Scenario 10: Exceed Remaining Cash

**Actors:** Waiter

**Preconditions:**
- Waiter has open shift with 100k remaining cash

**Steps:**
1. **Waiter:** Click "💰 Bàn giao một phần"
2. **Waiter:** Try to enter 150k (> remaining cash)
3. **Waiter:** Verify input validation prevents submission
4. **Waiter:** Or verify error message after submission
5. **Waiter:** Enter valid amount (≤ 100k)
6. **Waiter:** Verify submission succeeds

**Expected Results:**
- ✅ Cannot handover more than remaining cash
- ✅ Validation prevents invalid amounts
- ✅ Clear error message if attempted

---

## 📱 Mobile Testing Checklist

### Waiter Mobile Tests
- [ ] Cash status cards display correctly
- [ ] Buttons are touch-friendly (min 44x44px)
- [ ] Modals slide up from bottom
- [ ] Forms are easy to fill on mobile keyboard
- [ ] Pending banner is visible and clear
- [ ] History list scrolls smoothly
- [ ] All text is readable without zooming

### Cashier Mobile Tests
- [ ] Notification banner visible at top
- [ ] Quick actions section usable with thumbs
- [ ] Pending handovers list scrolls well
- [ ] Confirm modal fits on screen
- [ ] Input fields work with mobile keyboard
- [ ] Status badges are clear and readable
- [ ] Navigation between sections smooth

---

## 🔒 Security Testing

### Authorization Tests
- [ ] Waiter cannot access cashier endpoints
- [ ] Cashier cannot access other cashier's handovers
- [ ] Manager can access all handovers
- [ ] Unauthorized requests return 401
- [ ] Wrong role returns 403

### Data Validation Tests
- [ ] Negative amounts rejected
- [ ] Zero amounts rejected
- [ ] Non-numeric input rejected
- [ ] SQL injection attempts blocked
- [ ] XSS attempts sanitized

---

## ⚡ Performance Testing

### Load Tests
- [ ] 100 concurrent handover requests
- [ ] 1000 handovers in database
- [ ] Response time < 500ms
- [ ] No memory leaks
- [ ] Database queries optimized

### Stress Tests
- [ ] Multiple waiters creating handovers simultaneously
- [ ] Cashier processing many handovers quickly
- [ ] Large discrepancy calculations
- [ ] History with 100+ handovers

---

## 🐛 Error Handling Tests

### Network Errors
- [ ] Offline mode handling
- [ ] Timeout handling
- [ ] Retry mechanism
- [ ] Error messages clear

### Edge Cases
- [ ] Shift closed during handover
- [ ] Cashier shift closed during confirmation
- [ ] Concurrent handover attempts
- [ ] Database connection lost
- [ ] Invalid IDs in URLs

---

## ✅ Test Execution Checklist

### Before Testing
- [ ] Backend server running
- [ ] Frontend server running
- [ ] Database seeded with test data
- [ ] Test users created (waiter, cashier, manager)
- [ ] Test environment configured

### During Testing
- [ ] Record all test results
- [ ] Screenshot any errors
- [ ] Note performance issues
- [ ] Document unexpected behavior
- [ ] Test on multiple devices

### After Testing
- [ ] All scenarios passed
- [ ] No critical bugs found
- [ ] Performance acceptable
- [ ] Security verified
- [ ] Documentation updated

---

## 📊 Test Results Template

```
Test Date: YYYY-MM-DD
Tester: [Name]
Environment: [Dev/Staging/Production]

Scenario 1: Partial Handover - Happy Path
Status: ✅ PASS / ❌ FAIL
Notes: [Any observations]

Scenario 2: Handover with Discrepancy
Status: ✅ PASS / ❌ FAIL
Notes: [Any observations]

[... continue for all scenarios ...]

Overall Result: ✅ PASS / ❌ FAIL
Critical Issues: [List any critical issues]
Minor Issues: [List any minor issues]
Recommendations: [Any recommendations]
```

---

## 🚀 Automated Testing

### Unit Tests
```bash
cd backend
go test ./application/services/cash_handover_service_test.go -v
```

### Integration Tests
```bash
cd scripts
./test-handover-workflow.sh
```

### E2E Tests (Manual)
Follow scenarios 1-10 above

---

## 📝 Notes

- Test with real-world amounts (Vietnamese currency)
- Test during peak hours to verify performance
- Test on actual mobile devices, not just emulators
- Test with slow network connections
- Test with multiple users simultaneously
- Document all edge cases discovered
- Update test cases as features evolve

---

**Last Updated:** 2026-02-04  
**Status:** Ready for Testing  
**Coverage:** Backend + Frontend + E2E
