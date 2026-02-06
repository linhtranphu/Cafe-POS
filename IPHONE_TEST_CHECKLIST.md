# 📱 iPhone Safe Area - Test Checklist

## ✅ Pre-Test Setup

- [ ] Build frontend: `cd frontend && npm run build`
- [ ] Deploy to server
- [ ] Note the app URL: _______________

## 📱 iPhone Testing Steps

### 1. Add to Home Screen
- [ ] Open Safari on iPhone
- [ ] Navigate to app URL
- [ ] Tap Share button (square with arrow up)
- [ ] Select "Add to Home Screen"
- [ ] Name it "Cafe POS" or similar
- [ ] Tap "Add"

### 2. Open as Webapp
- [ ] Close Safari
- [ ] Find app icon on home screen
- [ ] Tap to open
- [ ] Verify it opens in standalone mode (no Safari UI)

## 🧪 Test Each View

### Manager Views
- [ ] **DashboardView** - Header not hidden by notch
- [ ] **OrderView** - Header visible, content scrollable
- [ ] **ShiftView** - Header visible
- [ ] **ManagerShiftView** - Header visible
- [ ] **FacilityManagementView** - Header visible, bottom nav clear
- [ ] **IngredientManagementView** - Header visible, bottom nav clear
- [ ] **ExpenseManagementView** - Header visible, bottom nav clear
- [ ] **UserManagementView** - Header visible
- [ ] **ProfileView** - Header visible, bottom nav clear

### Cashier Views
- [ ] **CashierDashboard** - Header visible, stats cards readable
- [ ] **CashierHandoverView** - Header visible, list scrollable
- [ ] **CashierReports** - Header visible, date picker accessible
- [ ] **CashierShiftClosure** - Header visible, form accessible

### Barista Views
- [ ] **BaristaView** - Header visible, order list scrollable

### Other Views
- [ ] **FacilityAddEditView** - Form header visible, bottom buttons accessible
- [ ] **MenuView** - Modals display correctly

## 🔍 Specific Checks

### Top Area (Notch)
- [ ] Header text fully visible
- [ ] Header icons/buttons not cut off
- [ ] Search bars fully accessible
- [ ] No content hidden behind notch

### Bottom Area (Home Indicator)
- [ ] Bottom navigation fully visible
- [ ] All nav icons accessible
- [ ] Nav labels readable
- [ ] No content hidden behind home indicator

### Scrollable Content
- [ ] Content scrolls smoothly
- [ ] Last item in list visible when scrolled to bottom
- [ ] No content cut off at bottom
- [ ] Pull-to-refresh works correctly

### Modals & Forms
- [ ] Modal headers visible
- [ ] Modal content scrollable
- [ ] Form inputs accessible
- [ ] Bottom buttons in modals accessible
- [ ] Close buttons (×) accessible

### Landscape Mode
- [ ] Rotate to landscape
- [ ] Content still visible
- [ ] No content hidden by rounded corners
- [ ] Navigation still accessible

## 📊 Test Results

### Device Information
- Device Model: _______________
- iOS Version: _______________
- Screen Size: _______________
- Has Notch: Yes / No
- Has Dynamic Island: Yes / No

### Overall Results
- [ ] All headers visible ✅
- [ ] All bottom navigation visible ✅
- [ ] All content scrollable ✅
- [ ] No content hidden ✅
- [ ] Pull-to-refresh works ✅

### Issues Found
List any issues here:
1. _______________
2. _______________
3. _______________

## 🐛 Common Issues to Check

- [ ] Header text touching notch
- [ ] Bottom nav hidden by home indicator
- [ ] Content cut off when scrolling
- [ ] Buttons not clickable near edges
- [ ] Modals not centered properly
- [ ] Form inputs hidden by keyboard

## 📸 Screenshots

Take screenshots of:
- [ ] Home screen with app icon
- [ ] Main dashboard view
- [ ] View with notch visible
- [ ] Bottom navigation
- [ ] Any issues found

## ✅ Sign Off

- Tested By: _______________
- Date: _______________
- Status: Pass / Fail / Needs Review
- Notes: _______________

---

**Testing Complete:** _______________  
**Ready for Production:** Yes / No
