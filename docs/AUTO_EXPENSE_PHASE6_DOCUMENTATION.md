# Auto Expense Tracking - Phase 6: Documentation & Deployment

**Date**: January 31, 2026  
**Status**: ✅ COMPLETE

## 6.1: User Guide

### For Managers: Auto Expense Tracking

#### What is Auto Expense Tracking?

Auto Expense Tracking automatically creates expense records when you:
- Purchase ingredients (create new or adjust stock IN)
- Purchase facilities/equipment
- Create maintenance records

This eliminates manual expense entry and ensures accurate financial tracking.

---

#### How It Works

##### 1. Creating Ingredients

When you create a new ingredient with initial stock:

**Steps**:
1. Go to **Nguyên liệu** (Ingredients)
2. Click **➕ Tạo nguyên liệu**
3. Fill in the form:
   - Name, category, unit
   - **Quantity**: Initial stock amount
   - **Cost per unit**: Purchase price
   - Supplier (optional)

**What Happens**:
- You'll see a green box showing: "✅ Tự động ghi nhận chi phí"
- The system calculates: `Quantity × Cost per unit`
- When you click "Thêm mới", the ingredient is created
- An expense is automatically created in category "Nguyên liệu"

**Example**:
```
Ingredient: Coffee Beans
Quantity: 10 kg
Cost per unit: 200,000 ₫
→ Auto expense: 2,000,000 ₫
```

---

##### 2. Adjusting Stock (Stock IN)

When you add stock to existing ingredients:

**Steps**:
1. Go to **Nguyên liệu**
2. Select an ingredient
3. Click **📦 Điều chỉnh**
4. Select type: **Thêm vào kho** (Add to stock)
5. Enter quantity and reason

**What Happens**:
- You'll see a green box showing the expense amount
- When you confirm, stock is adjusted
- An expense is automatically created

**Note**: Stock OUT (removing from stock) does NOT create an expense.

---

##### 3. Purchasing Facilities

When you purchase new equipment or facilities:

**Steps**:
1. Go to **Cơ sở vật chất** (Facilities)
2. Click **➕**
3. Fill in the form:
   - Name, type, area
   - **Cost**: Purchase price
   - Supplier (optional)

**What Happens**:
- You'll see a green box showing: "✅ Tự động ghi nhận chi phí"
- When you click "Thêm mới", the facility is created
- An expense is automatically created in category "Cơ sở vật chất"

**Example**:
```
Facility: Espresso Machine
Cost: 15,000,000 ₫
→ Auto expense: 15,000,000 ₫
```

---

##### 4. Recording Maintenance

When you create maintenance records:

**Steps**:
1. Go to **Cơ sở vật chất**
2. Select a facility
3. Create maintenance record with cost

**What Happens**:
- An expense is automatically created in category "Bảo trì"
- The expense links to the facility

---

#### Viewing Auto-Generated Expenses

**Steps**:
1. Go to **Chi phí** (Expenses)
2. Use the filter buttons at the top:
   - **Tất cả**: All expenses
   - **✍️ Thủ công**: Manual expenses only
   - **🥬 Nguyên liệu**: Ingredient purchases
   - **🏢 Cơ sở vật chất**: Facility purchases
   - **🔧 Bảo trì**: Maintenance costs

**Visual Indicators**:
- Auto-generated expenses show a colored badge:
  - 🥬 Tự động (Green) - Ingredient
  - 🏢 Tự động (Blue) - Facility
  - 🔧 Tự động (Orange) - Maintenance
- Manual expenses have no badge

---

#### Benefits

✅ **No Manual Entry**: Expenses created automatically  
✅ **Accurate Tracking**: No forgotten expenses  
✅ **Easy Reporting**: Filter by source type  
✅ **Audit Trail**: Link expenses to source items  
✅ **Time Saving**: Focus on operations, not data entry

---

#### Important Notes

⚠️ **Zero Cost Items**: Items with zero cost don't create expenses  
⚠️ **Editing Items**: Editing existing items doesn't create new expenses  
⚠️ **Stock OUT**: Removing stock doesn't create expenses  
⚠️ **Manual Expenses**: You can still create manual expenses as before

---

## 6.2: Admin Documentation

### System Architecture

```
User Action (Create/Adjust)
        ↓
Frontend Form
        ↓
HTTP Handler
        ↓
Service Layer
        ├─ Main Operation (Create/Update)
        └─ AutoExpenseService.Track*()
                ↓
        ExpenseService.CreateExpense()
                ↓
        Database (Expense Record)
```

### Configuration

#### Enable/Disable Auto-Expense

**Location**: `backend/main.go`

**To Disable**:
```go
// Comment out these lines:
// ingredientService.SetAutoExpenseService(autoExpenseService)
// facilityService.SetAutoExpenseService(autoExpenseService)
```

**To Enable**:
```go
// Ensure these lines are present:
autoExpenseService := services.NewAutoExpenseService(expenseService)
ingredientService.SetAutoExpenseService(autoExpenseService)
facilityService.SetAutoExpenseService(autoExpenseService)
```

---

### Category Management

**Default Categories** (auto-created):
- Nguyên liệu (Ingredient)
- Cơ sở vật chất (Facility)
- Bảo trì (Maintenance)
- Tiện ích (Utility)
- Nhân sự (Salary)
- Marketing
- Khác (Other)

**Location**: `backend/domain/expense/category.go`

**To Add Categories**:
Edit `GetDefaultCategories()` function.

---

### Database Schema

**Expense Collection**:
```json
{
  "_id": ObjectId,
  "date": ISODate,
  "category_id": ObjectId,
  "amount": Number,
  "description": String,
  "payment_method": String,
  "vendor": String,
  "notes": String,
  "source_type": String,      // "ingredient", "facility", "maintenance", "manual"
  "source_id": ObjectId,       // Links to source entity
  "created_at": ISODate,
  "updated_at": ISODate
}
```

**Indexes**:
- `category_id` (for filtering)
- `source_type` (for filtering)
- `source_id` (for linking)
- `date` (for date range queries)

---

### Monitoring

**Logs to Monitor**:
```
[AutoExpense] Created new category: Nguyên liệu (ID: ...)
[AutoExpense] Tracked ingredient purchase: Coffee Beans (5.00 kg) - Amount: 1000000.00 VND
[AutoExpense] Tracked facility purchase: Espresso Machine - Amount: 15000000.00 VND
[AutoExpense] Tracked maintenance: Coffee Grinder - Amount: 500000.00 VND
[AutoExpense] Skipping ingredient purchase tracking: zero cost or quantity
[AutoExpense] Category cache cleared
```

**Error Logs**:
```
[AutoExpense] Failed to get/create category for ingredient: [error]
[AutoExpense] Failed to create expense for ingredient purchase: [error]
```

**Note**: Errors are logged but don't fail main operations.

---

### Troubleshooting

#### Issue: Expenses Not Created

**Check**:
1. AutoExpenseService is wired up in `main.go`
2. Backend logs for errors
3. Category exists or can be created
4. Cost/quantity values are > 0

**Solution**:
- Check backend logs
- Verify database connection
- Ensure ExpenseService is working

---

#### Issue: Duplicate Expenses

**Possible Causes**:
- Race condition in concurrent requests
- Category cache not working

**Solution**:
- Check for duplicate API calls
- Review category caching logic
- Check database for duplicate records

---

#### Issue: Wrong Expense Amount

**Check**:
1. Frontend form values
2. Backend calculation logic
3. Currency conversion (if applicable)

**Solution**:
- Verify form data sent to backend
- Check AutoExpenseService calculation
- Review ingredient/facility cost values

---

## 6.3: Deployment Checklist

### Pre-Deployment

- [ ] All tests passing (Phase 5)
- [ ] Code reviewed and approved
- [ ] Documentation complete
- [ ] Database backup created
- [ ] Rollback plan prepared

### Backend Deployment

- [ ] Build backend: `go build`
- [ ] Run tests: `go test ./...`
- [ ] Deploy to staging
- [ ] Smoke test on staging
- [ ] Deploy to production
- [ ] Verify logs for errors
- [ ] Monitor for 1 hour

### Frontend Deployment

- [ ] Build frontend: `npm run build`
- [ ] Test build locally
- [ ] Deploy to staging
- [ ] Test all flows on staging
- [ ] Deploy to production
- [ ] Clear browser cache
- [ ] Test on production

### Post-Deployment

- [ ] Verify auto-expense creation
- [ ] Check expense filtering
- [ ] Monitor error logs
- [ ] Test with real data
- [ ] User acceptance sign-off
- [ ] Update documentation
- [ ] Announce to users

### Rollback Plan

**If Issues Found**:

1. **Backend Rollback**:
   ```bash
   # Revert to previous version
   git checkout <previous-commit>
   go build
   # Restart server
   ```

2. **Frontend Rollback**:
   ```bash
   # Revert to previous version
   git checkout <previous-commit>
   npm run build
   # Deploy previous build
   ```

3. **Database Rollback**:
   - Auto-generated expenses can be deleted by source_type
   - No schema changes, so no migration needed

4. **Disable Feature**:
   ```go
   // In main.go, comment out:
   // ingredientService.SetAutoExpenseService(autoExpenseService)
   // facilityService.SetAutoExpenseService(autoExpenseService)
   ```

---

## 6.4: Training Materials

### Quick Start Guide (1 Page)

**For Managers**:

**Auto Expense Tracking - Quick Start**

1. **Creating Ingredients**:
   - Fill in quantity and cost
   - See green box with expense amount
   - Click "Thêm mới"
   - Expense created automatically ✅

2. **Adjusting Stock**:
   - Select "Thêm vào kho" (Add)
   - Enter quantity
   - See green box with expense amount
   - Click "Xác nhận"
   - Expense created automatically ✅

3. **Purchasing Facilities**:
   - Fill in cost
   - See green box with expense amount
   - Click "Thêm mới"
   - Expense created automatically ✅

4. **Viewing Expenses**:
   - Go to "Chi phí"
   - Use filter buttons:
     - 🥬 Nguyên liệu
     - 🏢 Cơ sở vật chất
     - 🔧 Bảo trì
   - See auto-generated expenses with badges

**Benefits**: No manual expense entry needed!

---

### Video Tutorial Script

**Title**: "Auto Expense Tracking - Complete Guide"

**Duration**: 5 minutes

**Script**:

1. **Introduction** (30s)
   - What is auto expense tracking?
   - Why it's useful
   - What you'll learn

2. **Creating Ingredients** (1m 30s)
   - Navigate to Ingredients
   - Create new ingredient
   - Show auto-expense indicator
   - Confirm creation
   - Show expense in Expense Management

3. **Adjusting Stock** (1m 30s)
   - Select ingredient
   - Adjust stock IN
   - Show auto-expense indicator
   - Confirm adjustment
   - Show expense created

4. **Purchasing Facilities** (1m)
   - Navigate to Facilities
   - Create new facility
   - Show auto-expense indicator
   - Confirm creation
   - Show expense created

5. **Filtering Expenses** (1m)
   - Navigate to Expenses
   - Show filter buttons
   - Filter by source type
   - Show badges on expenses

6. **Conclusion** (30s)
   - Recap benefits
   - Where to get help
   - Thank you

---

## 6.5: Support Resources

### FAQ

**Q: Can I disable auto expense tracking?**  
A: Yes, contact your system administrator.

**Q: What if I don't want an expense created?**  
A: Set cost to 0, or edit the item after creation.

**Q: Can I edit auto-generated expenses?**  
A: Yes, they're regular expenses and can be edited or deleted.

**Q: What happens if expense creation fails?**  
A: The main operation (create ingredient/facility) still succeeds. Check with admin.

**Q: Can I see which ingredient/facility created an expense?**  
A: Yes, expenses have source badges and link to source items.

---

### Contact Information

**Technical Support**:
- Email: support@example.com
- Phone: +84 xxx xxx xxx
- Hours: Mon-Fri 9AM-6PM

**System Administrator**:
- Email: admin@example.com
- For configuration and troubleshooting

---

**Phase 6 Status**: ✅ COMPLETE  
**Documentation Pages**: 6  
**Training Materials**: 3  
**Support Resources**: 2
