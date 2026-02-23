# Testing Quick Start - Cashier Fund Handover

## 🚀 Quick Setup (5 minutes)

### 1. Start Servers

```bash
# Terminal 1: Backend
cd backend
go run main.go

# Terminal 2: Frontend
cd frontend
npm run dev
```

### 2. Login as Cashier

1. Open browser: `http://localhost:5173`
2. Login with cashier credentials
3. Open DevTools (F12)

### 3. Get JWT Token

```javascript
// In browser console
localStorage.getItem('token')
// Copy the token value
```

---

## 🧪 Run Automated Tests

```bash
# Set your token
export TOKEN="your_jwt_token_here"

# Run API tests
./test-fund-handover-api.sh

# Run frontend integration tests
node test-frontend-fund-handover.js
```

---

## 📋 Manual Testing Flow

### Test Scenario 1: Happy Path (No Variance)

1. **Dashboard**
   - ✅ See managed funds section
   - ✅ Note the amounts

2. **Start Closure**
   - Click "Đóng ca"
   - ✅ See managed funds summary

3. **Count Cash**
   - Enter exact expected cash
   - ✅ Variance = 0₫
   - ✅ No documentation needed

4. **Confirm**
   - Review summary
   - Click "Xác nhận và đóng ca"
   - ✅ Success!

**Time**: ~2 minutes

---

### Test Scenario 2: With Variance

1. **Dashboard**
   - ✅ See managed funds

2. **Start Closure**
   - Click "Đóng ca"

3. **Count Cash**
   - Enter amount different from expected
   - Example: Expected 2,000,000₫ → Enter 1,995,000₫
   - ✅ Variance = -5,000₫ (shortage)

4. **Document Variance**
   - Select reason: "COUNTING_ERROR"
   - Enter notes: "Đếm nhầm tờ 50k thành 100k"
   - ✅ Notes ≥ 10 characters

5. **Confirm**
   - Review summary with variance
   - Click "Xác nhận và đóng ca"
   - ✅ Success!

**Time**: ~3 minutes

---

## 🔍 What to Check

### Dashboard
- [ ] 💰 Managed funds section visible
- [ ] 💵 Cash (green) + 💳 Transfer (blue)
- [ ] 📊 Total (orange gradient)
- [ ] ⚠️ Warning message
- [ ] Pull-to-refresh works

### Closure Flow
- [ ] Managed funds summary matches dashboard
- [ ] Cash input works
- [ ] Variance calculates automatically
- [ ] Variance form appears when needed
- [ ] Confirmation shows all data
- [ ] Submit works

### Mobile
- [ ] Responsive layout
- [ ] Touch interactions smooth
- [ ] No horizontal scroll
- [ ] Keyboard doesn't cover inputs

---

## 🐛 Common Issues to Watch For

### Issue 1: Managed Funds Not Loading
**Symptom**: Empty or loading forever
**Check**: 
- Network tab for API call
- Console for errors
- Backend logs

### Issue 2: Variance Not Calculating
**Symptom**: Shows 0 when should show variance
**Check**:
- Expected cash value
- Actual cash input
- Console for calculation errors

### Issue 3: Cannot Submit
**Symptom**: Button disabled or error
**Check**:
- Variance documentation complete
- Notes ≥ 10 characters
- Network connectivity
- Console errors

### Issue 4: Mobile Layout Broken
**Symptom**: Overflow, cut-off text
**Check**:
- Viewport meta tag
- CSS responsive classes
- Safe area insets

---

## 📊 Quick Verification

### After Successful Closure

**1. Check Frontend**
```
✅ Redirected to dashboard
✅ Shift status = CLOSED
✅ Success message shown
```

**2. Check Database**
```javascript
// MongoDB shell
db.fund_handovers.findOne({ cashier_shift_id: ObjectId("...") })
// Should return the fund handover record

db.cashier_shifts.findOne({ _id: ObjectId("...") })
// Should show status: "CLOSED"
```

**3. Check API Response**
```json
{
  "shift": { "status": "CLOSED", ... },
  "fund_handover": {
    "cash_amount": 1995000,
    "transfer_amount": 800000,
    "variance_amount": -5000,
    ...
  }
}
```

---

## 🎯 Testing Priorities

### Priority 1: Critical Path (Must Work)
1. ✅ Dashboard displays managed funds
2. ✅ Closure flow completes successfully
3. ✅ Fund handover record created
4. ✅ No data loss

### Priority 2: Important Features
1. ✅ Variance calculation correct
2. ✅ Variance documentation works
3. ✅ Mobile responsive
4. ✅ Error handling

### Priority 3: Nice to Have
1. ✅ Pull-to-refresh smooth
2. ✅ Animations smooth
3. ✅ Loading states
4. ✅ Accessibility

---

## 📱 Mobile Testing Shortcuts

### iOS Safari
1. Connect iPhone via USB
2. Safari > Develop > iPhone > localhost
3. Test touch interactions
4. Check safe area insets

### Android Chrome
1. Enable USB debugging
2. Chrome > chrome://inspect
3. Select device
4. Test touch interactions

### Quick Mobile Checks
- [ ] Tap targets ≥ 44x44px
- [ ] Text readable without zoom
- [ ] No horizontal scroll
- [ ] Keyboard behavior correct

---

## 🔧 Troubleshooting

### Backend Not Running
```bash
cd backend
go run main.go
# Should see: Server running on :8080
```

### Frontend Not Running
```bash
cd frontend
npm run dev
# Should see: Local: http://localhost:5173
```

### MongoDB Not Connected
```bash
# Check MongoDB status
mongosh
# Should connect successfully
```

### Token Expired
```javascript
// Get new token
// 1. Logout
// 2. Login again
// 3. Copy new token from localStorage
```

---

## 📞 Need Help?

### Check Logs
- **Frontend**: Browser console (F12)
- **Backend**: Terminal running `go run main.go`
- **MongoDB**: `mongosh` and check collections

### Documentation
- Full testing guide: `FRONTEND_TESTING_GUIDE.md`
- Manual checklist: `MANUAL_TESTING_CHECKLIST.md`
- API docs: `CASHIER_FUND_HANDOVER_PHASE_4_COMPLETE.md`

### Common Commands
```bash
# View backend logs
cd backend && go run main.go

# View MongoDB data
mongosh
use cafe_pos
db.fund_handovers.find().pretty()
db.cashier_shifts.find({ status: "CLOSED" }).pretty()

# Clear test data
db.fund_handovers.deleteMany({})
db.cashier_shifts.updateMany({}, { $set: { status: "OPEN" } })
```

---

## ✅ Testing Complete?

After testing, verify:
- [ ] All critical paths work
- [ ] No console errors
- [ ] Mobile responsive
- [ ] Database records correct
- [ ] Issues documented

Then:
1. Fill out `MANUAL_TESTING_CHECKLIST.md`
2. Document any bugs found
3. Get stakeholder approval
4. Ready for deployment! 🚀
