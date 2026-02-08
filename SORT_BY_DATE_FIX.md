# ✅ Fix: Sort Items by Date (Newest First)

**Date:** February 6, 2026  
**Task:** Sort Facility, Ingredient, and Expense lists by creation date (newest to oldest)

---

## 🎯 User Request

> "review lại màn hình facility, ingredient, expense. hãy sort theo thứ tự thời gian từ mới đến cũ"

**Translation:** Sort the facility, ingredient, and expense screens by time from newest to oldest.

---

## ✅ What Was Fixed

### 3 Views Updated:

1. **ExpenseManagementView** - Sort expenses by date (newest first)
2. **IngredientManagementView** - Sort ingredients by created_at (newest first)
3. **FacilityManagementView** - Sort facilities by created_at (newest first)

---

## 🔧 Implementation

### 1. ExpenseManagementView.vue

**Before:**
```javascript
const filteredExpenses = computed(() => {
  let filtered = expenses.value
  
  // Filter by source type
  if (sourceFilter.value) {
    filtered = filtered.filter(e => {
      if (sourceFilter.value === 'manual') {
        return !e.source_type || e.source_type === 'manual'
      }
      return e.source_type === sourceFilter.value
    })
  }
  
  // Filter by search query
  if (searchQuery.value) {
    const query = searchQuery.value.toLowerCase()
    filtered = filtered.filter(e => 
      e.description?.toLowerCase().includes(query) ||
      e.vendor?.toLowerCase().includes(query)
    )
  }
  
  return filtered  // ❌ No sorting
})
```

**After:**
```javascript
const filteredExpenses = computed(() => {
  let filtered = expenses.value
  
  // Filter by source type
  if (sourceFilter.value) {
    filtered = filtered.filter(e => {
      if (sourceFilter.value === 'manual') {
        return !e.source_type || e.source_type === 'manual'
      }
      return e.source_type === sourceFilter.value
    })
  }
  
  // Filter by search query
  if (searchQuery.value) {
    const query = searchQuery.value.toLowerCase()
    filtered = filtered.filter(e => 
      e.description?.toLowerCase().includes(query) ||
      e.vendor?.toLowerCase().includes(query)
    )
  }
  
  // Sort by date (newest first) ✅
  return [...filtered].sort((a, b) => {
    const dateA = new Date(a.date || a.created_at || 0)
    const dateB = new Date(b.date || b.created_at || 0)
    return dateB - dateA // Newest first
  })
})
```

**Key changes:**
- Added sorting by `date` field (primary) or `created_at` (fallback)
- Newest items appear first
- Uses spread operator `[...filtered]` to avoid mutating original array

---

### 2. IngredientManagementView.vue

**Before:**
```javascript
const filteredIngredients = computed(() => {
  if (!searchQuery.value) return ingredients.value
  const query = searchQuery.value.toLowerCase()
  return ingredients.value.filter(i => 
    i.name.toLowerCase().includes(query) ||
    i.category.toLowerCase().includes(query) ||
    i.supplier?.toLowerCase().includes(query)
  )  // ❌ No sorting
})
```

**After:**
```javascript
const filteredIngredients = computed(() => {
  let filtered = ingredients.value
  
  // Filter by search query
  if (searchQuery.value) {
    const query = searchQuery.value.toLowerCase()
    filtered = filtered.filter(i => 
      i.name.toLowerCase().includes(query) ||
      i.category.toLowerCase().includes(query) ||
      i.supplier?.toLowerCase().includes(query)
    )
  }
  
  // Sort by created_at (newest first) ✅
  return [...filtered].sort((a, b) => {
    const dateA = new Date(a.created_at || 0)
    const dateB = new Date(b.created_at || 0)
    return dateB - dateA // Newest first
  })
})
```

**Key changes:**
- Added sorting by `created_at` field
- Newest ingredients appear first
- Consistent filtering logic (always filter, even if no search query)

---

### 3. FacilityManagementView.vue

**Before:**
```javascript
const filteredFacilities = computed(() => {
  if (!searchQuery.value) return facilities.value
  const query = searchQuery.value.toLowerCase()
  return facilities.value.filter(f => 
    f.name?.toLowerCase().includes(query) ||
    f.type?.toLowerCase().includes(query) ||
    f.area?.toLowerCase().includes(query)
  )  // ❌ No sorting
})
```

**After:**
```javascript
const filteredFacilities = computed(() => {
  let filtered = facilities.value
  
  // Filter by search query
  if (searchQuery.value) {
    const query = searchQuery.value.toLowerCase()
    filtered = filtered.filter(f => 
      f.name?.toLowerCase().includes(query) ||
      f.type?.toLowerCase().includes(query) ||
      f.area?.toLowerCase().includes(query)
    )
  }
  
  // Sort by created_at (newest first) ✅
  return [...filtered].sort((a, b) => {
    const dateA = new Date(a.created_at || 0)
    const dateB = new Date(b.created_at || 0)
    return dateB - dateA // Newest first
  })
})
```

**Key changes:**
- Added sorting by `created_at` field
- Newest facilities appear first
- Consistent filtering logic

---

## 📊 Sorting Logic

### Sort Order:

```
Newest (Most Recent)
    ↓
    ↓
    ↓
Oldest
```

### Implementation:

```javascript
return [...filtered].sort((a, b) => {
  const dateA = new Date(a.created_at || 0)
  const dateB = new Date(b.created_at || 0)
  return dateB - dateA // Newest first (descending)
})
```

**Explanation:**
- `dateB - dateA` = Descending order (newest first)
- `dateA - dateB` = Ascending order (oldest first)
- Fallback to `0` if date is missing (appears at bottom)

---

## 🎨 User Experience

### Before Fix:

```
Items displayed in random/database order:
- Item from 2 weeks ago
- Item from yesterday
- Item from 1 month ago
- Item from today
❌ Hard to find recent items
```

### After Fix:

```
Items displayed newest first:
- Item from today ← Most recent
- Item from yesterday
- Item from 2 weeks ago
- Item from 1 month ago
✅ Easy to find recent items
```

---

## 📱 Visual Example

### Expense List (After Fix):

```
┌─────────────────────────────────┐
│ 💰 Chi phí                      │
├─────────────────────────────────┤
│ 📋 Danh sách chi phí            │
├─────────────────────────────────┤
│ ┌─────────────────────────────┐ │
│ │ Mua cà phê                  │ │ ← Today
│ │ 📅 06/02/2026               │ │
│ │ -500,000 ₫                  │ │
│ └─────────────────────────────┘ │
│ ┌─────────────────────────────┐ │
│ │ Sửa máy pha                 │ │ ← Yesterday
│ │ 📅 05/02/2026               │ │
│ │ -1,200,000 ₫                │ │
│ └─────────────────────────────┘ │
│ ┌─────────────────────────────┐ │
│ │ Mua sữa tươi                │ │ ← 3 days ago
│ │ 📅 03/02/2026               │ │
│ │ -300,000 ₫                  │ │
│ └─────────────────────────────┘ │
└─────────────────────────────────┘
```

---

## 🧪 Testing

### Test Cases:

**1. Create New Item:**
```
1. Create new expense/ingredient/facility
2. Check list
3. ✅ New item should appear at TOP
```

**2. Search Filter:**
```
1. Search for items
2. Check filtered results
3. ✅ Results should still be sorted (newest first)
```

**3. Multiple Items Same Day:**
```
1. Create multiple items on same day
2. Check list
3. ✅ Should be sorted by creation time (newest first)
```

**4. Empty List:**
```
1. View empty list
2. ✅ Should show "Không có ... nào" message
3. ✅ No errors
```

---

## 📁 Files Modified

### Code (3 files):
1. `frontend/src/views/ExpenseManagementView.vue`
2. `frontend/src/views/IngredientManagementView.vue`
3. `frontend/src/views/FacilityManagementView.vue`

### Changes per file:
- Updated `filteredExpenses` computed property
- Updated `filteredIngredients` computed property
- Updated `filteredFacilities` computed property

---

## ✅ Quality Check

- [x] Code implemented
- [x] No TypeScript/lint errors
- [x] No diagnostics issues
- [x] Consistent sorting logic across all 3 views
- [x] Handles missing dates gracefully (fallback to 0)
- [x] Doesn't mutate original arrays (uses spread operator)
- [ ] Testing on device (pending user)

---

## 💡 Technical Details

### Why Spread Operator?

```javascript
// ❌ BAD - Mutates original array
return filtered.sort((a, b) => ...)

// ✅ GOOD - Creates new array
return [...filtered].sort((a, b) => ...)
```

**Reason:** Vue's reactivity system tracks array mutations. Sorting in-place can cause unexpected behavior.

### Date Handling:

```javascript
const dateA = new Date(a.created_at || 0)
```

**Fallback to 0:**
- If `created_at` is missing/null/undefined
- `new Date(0)` = January 1, 1970 (Unix epoch)
- These items appear at bottom of list

### Sort Comparison:

```javascript
return dateB - dateA
```

**Result:**
- Positive number → b comes before a
- Negative number → a comes before b
- Zero → keep original order

**Example:**
```
dateB = 2026-02-06 (newer)
dateA = 2026-02-05 (older)
dateB - dateA = positive → b comes first ✅
```

---

## 🎯 Summary

**Problem:** Items displayed in random order, hard to find recent items

**Solution:** Sort by creation date (newest first) in all 3 views

**Result:**
- ✅ Expenses sorted by `date` or `created_at` (newest first)
- ✅ Ingredients sorted by `created_at` (newest first)
- ✅ Facilities sorted by `created_at` (newest first)
- ✅ Consistent sorting across all views
- ✅ Easy to find recent items

**Impact:**
- Better UX - Recent items at top
- Consistent behavior across views
- No breaking changes
- Backward compatible

---

**Implementation Date:** February 6, 2026  
**Developer:** Kiro AI  
**Status:** ✅ Complete, Ready for Testing
