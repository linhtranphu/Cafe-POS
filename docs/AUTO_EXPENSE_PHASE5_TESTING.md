# Auto Expense Tracking - Phase 5: Testing & Validation

**Date**: January 31, 2026  
**Status**: ✅ READY FOR TESTING

## Testing Checklist

### 5.1: Manual Testing

#### Test Case 1: Create Ingredient with Initial Stock
**Steps**:
1. Navigate to Ingredient Management
2. Click "Tạo nguyên liệu"
3. Fill in form:
   - Name: "Test Coffee Beans"
   - Category: "Beverage"
   - Unit: "kg"
   - Quantity: 10
   - Cost per unit: 200,000
   - Supplier: "Test Supplier"
4. Verify auto-expense indicator shows: 2,000,000 ₫
5. Click "Thêm mới"
6. Navigate to Expense Management
7. Filter by "🥬 Nguyên liệu"
8. Verify expense created:
   - Description: "Nhập nguyên liệu: Test Coffee Beans"
   - Amount: 2,000,000 ₫
   - Category: "Nguyên liệu"
   - Source badge: "🥬 Tự động"

**Expected Result**: ✅ Expense automatically created and visible

---

#### Test Case 2: Adjust Stock IN
**Steps**:
1. Navigate to Ingredient Management
2. Select existing ingredient
3. Click "📦 Điều chỉnh"
4. Select type: "Thêm vào kho"
5. Enter quantity: 5
6. Enter reason: "Restocking"
7. Verify auto-expense indicator shows calculated amount
8. Click "Xác nhận"
9. Navigate to Expense Management
10. Filter by "🥬 Nguyên liệu"
11. Verify new expense created for stock adjustment

**Expected Result**: ✅ Expense automatically created for stock IN

---

#### Test Case 3: Adjust Stock OUT (No Expense)
**Steps**:
1. Navigate to Ingredient Management
2. Select existing ingredient
3. Click "📦 Điều chỉnh"
4. Select type: "Lấy ra khỏi kho"
5. Enter quantity: -3
6. Enter reason: "Used for production"
7. Verify NO auto-expense indicator shows
8. Click "Xác nhận"
9. Navigate to Expense Management
10. Verify NO new expense created

**Expected Result**: ✅ No expense created for stock OUT

---

#### Test Case 4: Create Ingredient with Zero Cost
**Steps**:
1. Navigate to Ingredient Management
2. Click "Tạo nguyên liệu"
3. Fill in form with cost_per_unit = 0
4. Verify NO auto-expense indicator shows
5. Click "Thêm mới"
6. Navigate to Expense Management
7. Verify NO expense created

**Expected Result**: ✅ No expense created for zero cost

---

#### Test Case 5: Create Facility
**Steps**:
1. Navigate to Facility Management
2. Click "➕" to create facility
3. Fill in form:
   - Name: "Test Espresso Machine"
   - Type: "Equipment"
   - Area: "Kitchen"
   - Quantity: 1
   - Cost: 15,000,000
   - Supplier: "Equipment Ltd."
4. Verify auto-expense indicator shows: 15,000,000 ₫
5. Click "Thêm mới"
6. Navigate to Expense Management
7. Filter by "🏢 Cơ sở vật chất"
8. Verify expense created:
   - Description: "Mua thiết bị: Test Espresso Machine"
   - Amount: 15,000,000 ₫
   - Category: "Cơ sở vật chất"
   - Source badge: "🏢 Tự động"

**Expected Result**: ✅ Expense automatically created

---

#### Test Case 6: Create Maintenance Record
**Steps**:
1. Navigate to Facility Management
2. Select existing facility
3. Create maintenance record:
   - Description: "Replace parts"
   - Cost: 500,000
   - Type: "repair"
4. Click "Thêm mới"
5. Navigate to Expense Management
6. Filter by "🔧 Bảo trì"
7. Verify expense created:
   - Description: "Bảo trì: [Facility Name]"
   - Amount: 500,000 ₫
   - Category: "Bảo trì"
   - Source badge: "🔧 Tự động"

**Expected Result**: ✅ Expense automatically created

---

#### Test Case 7: Expense Source Filtering
**Steps**:
1. Navigate to Expense Management
2. Create mix of expenses (manual + auto)
3. Test each filter:
   - Click "Tất cả" → Shows all expenses
   - Click "✍️ Thủ công" → Shows only manual expenses
   - Click "🥬 Nguyên liệu" → Shows only ingredient expenses
   - Click "🏢 Cơ sở vật chất" → Shows only facility expenses
   - Click "🔧 Bảo trì" → Shows only maintenance expenses
4. Verify filter counts match displayed items
5. Combine with search query
6. Verify both filters work together

**Expected Result**: ✅ Filtering works correctly

---

#### Test Case 8: Edit Existing Items (No New Expense)
**Steps**:
1. Edit existing ingredient (change name, quantity, etc.)
2. Verify NO auto-expense indicator shows
3. Save changes
4. Navigate to Expense Management
5. Verify NO new expense created
6. Repeat for facility

**Expected Result**: ✅ No expense created when editing

---

### 5.2: Integration Testing

#### Test Case 9: Concurrent Operations
**Steps**:
1. Open multiple browser tabs
2. Create ingredients simultaneously in different tabs
3. Verify all expenses created correctly
4. Check for duplicate expenses
5. Verify category cache works correctly

**Expected Result**: ✅ No race conditions or duplicates

---

#### Test Case 10: Error Handling
**Steps**:
1. Simulate backend error (disconnect network)
2. Create ingredient
3. Verify ingredient created but expense fails gracefully
4. Check logs for error messages
5. Verify main operation succeeded

**Expected Result**: ✅ Graceful degradation works

---

#### Test Case 11: Large Data Sets
**Steps**:
1. Create 100+ ingredients with costs
2. Verify all expenses created
3. Test filtering performance
4. Check memory usage
5. Verify no performance degradation

**Expected Result**: ✅ Performance acceptable

---

### 5.3: User Acceptance Testing

#### Test Case 12: User Workflow - Ingredient Purchase
**Scenario**: Manager purchases 10kg coffee beans for 2,000,000 ₫

**Steps**:
1. Manager logs in
2. Navigates to Ingredient Management
3. Creates new ingredient with purchase details
4. Sees auto-expense indicator
5. Confirms creation
6. Checks Expense Management
7. Sees expense automatically recorded
8. Verifies expense details match purchase

**Expected Result**: ✅ Seamless workflow, no manual expense entry needed

---

#### Test Case 13: User Workflow - Facility Purchase
**Scenario**: Manager purchases espresso machine for 15,000,000 ₫

**Steps**:
1. Manager logs in
2. Navigates to Facility Management
3. Creates new facility with purchase details
4. Sees auto-expense indicator
5. Confirms creation
6. Checks Expense Management
7. Sees expense automatically recorded
8. Can filter to see only facility expenses

**Expected Result**: ✅ Seamless workflow, clear expense tracking

---

#### Test Case 14: User Workflow - Monthly Report
**Scenario**: Manager reviews monthly expenses by source

**Steps**:
1. Manager navigates to Expense Management
2. Uses source filters to review:
   - All ingredient purchases
   - All facility purchases
   - All maintenance costs
   - All manual expenses
3. Compares totals
4. Exports or analyzes data

**Expected Result**: ✅ Easy to analyze expenses by source

---

## Testing Results Template

```
Test Case: [Number and Name]
Date: [Date]
Tester: [Name]
Status: [ ] PASS [ ] FAIL [ ] BLOCKED

Steps Executed:
1. [Step 1]
2. [Step 2]
...

Actual Result:
[What actually happened]

Expected Result:
[What should have happened]

Issues Found:
- [Issue 1]
- [Issue 2]

Screenshots:
[Attach screenshots if applicable]

Notes:
[Additional observations]
```

## Bug Report Template

```
Bug ID: [AUTO-EXP-XXX]
Severity: [ ] Critical [ ] High [ ] Medium [ ] Low
Priority: [ ] P0 [ ] P1 [ ] P2 [ ] P3

Title: [Short description]

Description:
[Detailed description of the bug]

Steps to Reproduce:
1. [Step 1]
2. [Step 2]
...

Expected Behavior:
[What should happen]

Actual Behavior:
[What actually happens]

Environment:
- Browser: [Chrome/Firefox/Safari]
- OS: [Windows/Mac/Linux]
- Backend Version: [Version]
- Frontend Version: [Version]

Logs:
[Relevant log entries]

Screenshots:
[Attach screenshots]

Workaround:
[If any workaround exists]
```

## Performance Benchmarks

### Target Metrics:
- Ingredient creation with expense: < 500ms
- Facility creation with expense: < 500ms
- Expense filtering: < 100ms
- Category cache hit rate: > 95%
- Concurrent operations: No race conditions

### Monitoring:
- Backend logs for expense tracking
- Frontend console for errors
- Network tab for API calls
- Database queries for performance

## Sign-Off Checklist

- [ ] All unit tests passing
- [ ] All manual test cases executed
- [ ] Integration tests completed
- [ ] User acceptance tests completed
- [ ] Performance benchmarks met
- [ ] No critical or high-priority bugs
- [ ] Documentation updated
- [ ] Code reviewed
- [ ] Ready for production deployment

---

**Phase 5 Status**: ✅ READY FOR TESTING  
**Test Cases**: 14  
**Estimated Testing Time**: 4-6 hours  
**Required Testers**: 2-3 (1 QA, 1-2 end users)
