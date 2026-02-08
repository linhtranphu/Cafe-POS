# ✅ Expense Filter: Source Type → Creator

**Date:** February 6, 2026  
**Task:** Replace source type filter with creator filter in ExpenseManagementView

---

## 🎯 User Request

> "edit màn hình expense, bên dưới phần 'Search Bar', thay vì hiển thị source type filter. hiển thị danh sách người tạo expense"

**Translation:** Replace the source type filter below the search bar with a list of expense creators.

---

## ✅ What Was Changed

### ExpenseManagementView.vue

**Before:**
```vue
<!-- Source Type Filter -->
<div class="mt-2 flex gap-2 overflow-x-auto pb-2">
  <button @click="sourceFilter = ''">Tất cả</button>
  <button @click="sourceFilter = 'manual'">✍️ Thủ công</button>
  <button @click="sourceFilter = 'ingredient'">🥬 Nguyên liệu</button>
  <button @click="sourceFilter = 'facility'">🏢 Cơ sở vật chất</button>
  <button @click="sourceFilter = 'maintenance'">🔧 Bảo trì</button>
</div>
```

**After:**
```vue
<!-- Creator Filter -->
<div class="mt-2 flex gap-2 overflow-x-auto pb-2">
  <button @click="creatorFilter = ''">👥 Tất cả</button>
  <button v-for="creator in uniqueCreators" :key="creator"
    @click="creatorFilter = creator">
    👤 {{ creator }}
  </button>
</div>
```

---

## 🔧 Implementation Details

### 1. Template Changes

**Replaced:**
- Source type filter buttons (manual, ingredient, facility, maintenance)

**With:**
- Dynamic creator filter buttons (generated from expense data)

### 2. Script Changes

**Added:**
```javascript
// Get unique creators from expenses
const uniqueCreators = computed(() => {
  const creators = expenses.value
    .map(e => e.created_by || 'Hệ thống')
    .filter((value, index, self) => self.indexOf(value) === index)
  return creators.sort()
})
```

**Changed:**
```javascript
// Before
const sourceFilter = ref('')

// After
const creatorFilter = ref('')
```

**Updated Filter Logic:**
```javascript
// Before - Filter by source type
if (sourceFilter.value) {
  filtered = filtered.filter(e => {
    if (sourceFilter.value === 'manual') {
      return !e.source_type || e.source_type === 'manual'
    }
    return e.source_type === sourceFilter.value
  })
}

// After - Filter by creator
if (creatorFilter.value) {
  filtered = filtered.filter(e => {
    const creator = e.created_by || 'Hệ thống'
    return creator === creatorFilter.value
  })
}
```

---

## 🎨 User Experience

### Before:

```
┌─────────────────────────────────┐
│ 💰 Chi phí                      │
├─────────────────────────────────┤
│ [Search Bar]                    │
├─────────────────────────────────┤
│ Filter by Source Type:          │
│ [Tất cả] [✍️ Thủ công]          │
│ [🥬 Nguyên liệu] [🏢 CSVC]      │
│ [🔧 Bảo trì]                    │
└─────────────────────────────────┘
```

### After:

```
┌─────────────────────────────────┐
│ 💰 Chi phí                      │
├─────────────────────────────────┤
│ [Search Bar]                    │
├─────────────────────────────────┤
│ Filter by Creator:              │
│ [👥 Tất cả] [👤 Admin]          │
│ [👤 Hệ thống] [👤 Manager]      │
│ [👤 Cashier]                    │
└─────────────────────────────────┘
```

---

## 💡 How It Works

### Dynamic Creator List:

```javascript
const uniqueCreators = computed(() => {
  // 1. Get all creators from expenses
  const creators = expenses.value.map(e => e.created_by || 'Hệ thống')
  
  // 2. Remove duplicates
  const unique = creators.filter((value, index, self) => 
    self.indexOf(value) === index
  )
  
  // 3. Sort alphabetically
  return unique.sort()
})
```

**Example:**
```
Expenses:
- Expense 1: created_by = "Admin"
- Expense 2: created_by = "Manager"
- Expense 3: created_by = null (→ "Hệ thống")
- Expense 4: created_by = "Admin"
- Expense 5: created_by = "Cashier"

uniqueCreators = ["Admin", "Cashier", "Hệ thống", "Manager"]
                  ↑ Sorted alphabetically, no duplicates
```

### Filter Logic:

```javascript
if (creatorFilter.value) {
  filtered = filtered.filter(e => {
    const creator = e.created_by || 'Hệ thống'
    return creator === creatorFilter.value
  })
}
```

**Example:**
```
User clicks: [👤 Admin]
↓
creatorFilter = "Admin"
↓
Show only expenses where created_by = "Admin"
```

---

## 📱 Visual Example

### Filter Buttons:

```
┌─────────────────────────────────────────────────┐
│ [👥 Tất cả] [👤 Admin] [👤 Cashier]             │
│ [👤 Hệ thống] [👤 Manager]                      │
└─────────────────────────────────────────────────┘
     ↑ Active (purple)    ↑ Inactive (white)
```

### Filtered Results:

**Click [👤 Admin]:**
```
┌─────────────────────────────────┐
│ Mua cà phê                      │
│ 👤 Admin                        │
│ -500,000 ₫                      │
├─────────────────────────────────┤
│ Sửa máy pha                     │
│ 👤 Admin                        │
│ -1,200,000 ₫                    │
└─────────────────────────────────┘
Only expenses created by Admin
```

**Click [👥 Tất cả]:**
```
┌─────────────────────────────────┐
│ Mua cà phê                      │
│ 👤 Admin                        │
├─────────────────────────────────┤
│ Nhập nguyên liệu                │
│ 👤 Hệ thống                     │
├─────────────────────────────────┤
│ Tiền điện                       │
│ 👤 Manager                      │
└─────────────────────────────────┘
All expenses
```

---

## 🧪 Testing

### Test Cases:

**1. Display Unique Creators:**
```
1. Open expense view
2. Check filter buttons
3. ✅ Should show unique creators only (no duplicates)
4. ✅ Should be sorted alphabetically
5. ✅ Should include "Hệ thống" for auto-created expenses
```

**2. Filter by Creator:**
```
1. Click a creator button (e.g., "Admin")
2. Check expense list
3. ✅ Should show only expenses by that creator
4. ✅ Other expenses should be hidden
```

**3. Show All:**
```
1. Click "Tất cả" button
2. Check expense list
3. ✅ Should show all expenses
4. ✅ No filtering applied
```

**4. Empty State:**
```
1. Filter by creator with no expenses
2. Check display
3. ✅ Should show "Không có chi phí nào"
```

**5. Dynamic Update:**
```
1. Create new expense with new creator
2. Check filter buttons
3. ✅ New creator should appear in filter list
```

---

## 📊 Benefits

### Why This Change?

**Before (Source Type Filter):**
- ❌ Fixed categories (manual, ingredient, facility, maintenance)
- ❌ Not flexible
- ❌ Doesn't show who created expenses

**After (Creator Filter):**
- ✅ Dynamic list based on actual data
- ✅ Shows all creators automatically
- ✅ Easy to track expenses by person
- ✅ Better for accountability

### Use Cases:

1. **Manager wants to see expenses by specific staff:**
   - Click staff name → See their expenses

2. **Check auto-created expenses:**
   - Click "Hệ thống" → See system-generated expenses

3. **Audit expenses by creator:**
   - Filter by each person → Review their spending

---

## 🔍 Edge Cases

### 1. Null/Undefined created_by:

```javascript
const creator = e.created_by || 'Hệ thống'
```

**Handling:**
- If `created_by` is null/undefined → Show as "Hệ thống"
- Consistent display

### 2. Empty Expense List:

```javascript
const uniqueCreators = computed(() => {
  const creators = expenses.value.map(...)
  // If expenses.value is empty → creators = []
  return creators.sort() // Returns []
})
```

**Result:**
- Only "Tất cả" button shows
- No creator buttons (no data)

### 3. Single Creator:

```
uniqueCreators = ["Admin"]
```

**Display:**
```
[👥 Tất cả] [👤 Admin]
```

---

## 📁 Files Modified

**1. frontend/src/views/ExpenseManagementView.vue**

**Changes:**
- Template: Replaced source type filter with creator filter
- Script: Added `uniqueCreators` computed property
- Script: Changed `sourceFilter` to `creatorFilter`
- Script: Updated filter logic

**Lines changed:** ~30 lines

---

## ✅ Quality Check

- [x] Code implemented
- [x] No TypeScript/lint errors
- [x] No diagnostics issues
- [x] Dynamic creator list
- [x] Handles null/undefined creators
- [x] Sorted alphabetically
- [x] No duplicates
- [ ] Testing on device (pending user)

---

## 💡 Technical Notes

### Array.filter() for Unique Values:

```javascript
const unique = array.filter((value, index, self) => 
  self.indexOf(value) === index
)
```

**How it works:**
- `self.indexOf(value)` → First occurrence index
- `index` → Current index
- If they match → First occurrence → Keep it
- If they don't match → Duplicate → Remove it

**Example:**
```javascript
["Admin", "Manager", "Admin", "Cashier"]
         ↓
["Admin", "Manager", "Cashier"]
```

### Computed Property Reactivity:

```javascript
const uniqueCreators = computed(() => {
  const creators = expenses.value.map(...)
  // Automatically updates when expenses.value changes
})
```

**Benefits:**
- Auto-updates when expenses change
- No manual refresh needed
- Efficient (only recalculates when dependencies change)

---

## 🎯 Summary

**Problem:** Source type filter not useful for tracking who created expenses

**Solution:** Replace with dynamic creator filter

**Result:**
- ✅ Shows all unique creators
- ✅ Filter expenses by creator
- ✅ Dynamic (updates automatically)
- ✅ Sorted alphabetically
- ✅ Handles edge cases

**Impact:**
- Better expense tracking
- Easy to see who created what
- Better accountability
- More flexible than fixed categories

---

**Implementation Date:** February 6, 2026  
**Developer:** Kiro AI  
**Status:** ✅ Complete, Ready for Testing
