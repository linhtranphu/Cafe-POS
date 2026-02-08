# ✅ Facility Constants Cleanup

**Date:** February 6, 2026  
**Task:** Simplify FACILITY_TYPE_OPTIONS to match backend and reduce redundancy

---

## 🎯 Issue Identified

User asked: "còn FACILITY_TYPE_OPTIONS thì sao? đang sử dụng FACILITY_TYPES hay FACILITY_TYPE_OPTIONS"

This revealed a redundancy issue in the constants.

---

## 🔍 Problem Analysis

### Backend (Go):

```go
// backend/domain/facility/facility.go
const (
    TypeFurniture = "Bàn ghế"
    TypeMachine   = "Máy móc"
    TypeUtensil   = "Dụng cụ"
    TypeElectric  = "Điện tử"
    TypeOther     = "Khác"
)
```

**Backend stores:** Vietnamese strings directly (`"Bàn ghế"`, `"Máy móc"`, etc.)

### Frontend (Before):

```javascript
export const FACILITY_TYPES = {
  FURNITURE: 'Bàn ghế',
  MACHINE: 'Máy móc',
  UTENSIL: 'Dụng cụ',
  ELECTRIC: 'Điện tử',
  OTHER: 'Khác'
}

export const FACILITY_TYPE_OPTIONS = [
  { value: FACILITY_TYPES.FURNITURE, label: 'Bàn ghế' },  // value='Bàn ghế', label='Bàn ghế'
  { value: FACILITY_TYPES.MACHINE, label: 'Máy móc' },    // value='Máy móc', label='Máy móc'
  { value: FACILITY_TYPES.UTENSIL, label: 'Dụng cụ' },    // value='Dụng cụ', label='Dụng cụ'
  { value: FACILITY_TYPES.ELECTRIC, label: 'Điện tử' },   // value='Điện tử', label='Điện tử'
  { value: FACILITY_TYPES.OTHER, label: 'Khác' }          // value='Khác', label='Khác'
]
```

**Problem:**
- `value` and `label` are identical
- Hardcoded strings in `label` duplicate `FACILITY_TYPES` values
- Redundant and error-prone (if FACILITY_TYPES changes, must update OPTIONS too)

### Frontend Usage:

```javascript
const facilityCategories = computed(() => {
  const defaultCategories = FACILITY_TYPE_OPTIONS.map(opt => opt.label)
  // Result: ['Bàn ghế', 'Máy móc', 'Dụng cụ', 'Điện tử', 'Khác']
  
  const backendTypes = facilityStore.types.map(t => t.name)
  return [...new Set([...defaultCategories, ...backendTypes])]
})
```

**Issue:** Using `opt.label` when `opt.value` would work the same.

---

## ✅ Solution

### Frontend (After):

```javascript
export const FACILITY_TYPES = {
  FURNITURE: 'Bàn ghế',
  MACHINE: 'Máy móc',
  UTENSIL: 'Dụng cụ',
  ELECTRIC: 'Điện tử',
  OTHER: 'Khác'
}

// For backward compatibility and convenience
export const FACILITY_TYPE_OPTIONS = [
  { value: FACILITY_TYPES.FURNITURE, label: FACILITY_TYPES.FURNITURE },
  { value: FACILITY_TYPES.MACHINE, label: FACILITY_TYPES.MACHINE },
  { value: FACILITY_TYPES.UTENSIL, label: FACILITY_TYPES.UTENSIL },
  { value: FACILITY_TYPES.ELECTRIC, label: FACILITY_TYPES.ELECTRIC },
  { value: FACILITY_TYPES.OTHER, label: FACILITY_TYPES.OTHER }
]
```

**Benefits:**
- ✅ Single source of truth (`FACILITY_TYPES`)
- ✅ No hardcoded duplicates
- ✅ If `FACILITY_TYPES` changes, `OPTIONS` updates automatically
- ✅ Clearer that `value` and `label` are the same
- ✅ Easier to maintain

---

## 🔧 Changes Made

### File Modified:

**frontend/src/constants/facility.js**

### Before:

```javascript
export const FACILITY_TYPE_OPTIONS = [
  { value: FACILITY_TYPES.FURNITURE, label: 'Bàn ghế' },
  { value: FACILITY_TYPES.MACHINE, label: 'Máy móc' },
  { value: FACILITY_TYPES.UTENSIL, label: 'Dụng cụ' },
  { value: FACILITY_TYPES.ELECTRIC, label: 'Điện tử' },
  { value: FACILITY_TYPES.OTHER, label: 'Khác' }
]
```

### After:

```javascript
export const FACILITY_TYPE_OPTIONS = [
  { value: FACILITY_TYPES.FURNITURE, label: FACILITY_TYPES.FURNITURE },
  { value: FACILITY_TYPES.MACHINE, label: FACILITY_TYPES.MACHINE },
  { value: FACILITY_TYPES.UTENSIL, label: FACILITY_TYPES.UTENSIL },
  { value: FACILITY_TYPES.ELECTRIC, label: FACILITY_TYPES.ELECTRIC },
  { value: FACILITY_TYPES.OTHER, label: FACILITY_TYPES.OTHER }
]
```

---

## 📊 Impact Analysis

### No Breaking Changes:

**1. Usage in FacilityManagementView:**
```javascript
const facilityCategories = computed(() => {
  const defaultCategories = FACILITY_TYPE_OPTIONS.map(opt => opt.label)
  // Before: ['Bàn ghế', 'Máy móc', 'Dụng cụ', 'Điện tử', 'Khác']
  // After:  ['Bàn ghế', 'Máy móc', 'Dụng cụ', 'Điện tử', 'Khác']
  // ✅ Same result
})
```

**2. Select Options:**
```vue
<select v-model="formData.type">
  <option v-for="cat in facilityCategories" :key="cat" :value="cat">
    {{ cat }}
  </option>
</select>
```
- `facilityCategories` still contains same values
- ✅ No change in behavior

**3. Backend Communication:**
```javascript
await facilityStore.createFacility({
  type: 'Khác'  // Still sends Vietnamese string
})
```
- Backend expects Vietnamese strings
- ✅ No change needed

### Why This Works:

**Before:**
```javascript
FACILITY_TYPES.OTHER = 'Khác'
OPTIONS[4] = { value: 'Khác', label: 'Khác' }
```

**After:**
```javascript
FACILITY_TYPES.OTHER = 'Khác'
OPTIONS[4] = { value: FACILITY_TYPES.OTHER, label: FACILITY_TYPES.OTHER }
// Evaluates to: { value: 'Khác', label: 'Khác' }
```

**Result:** Identical runtime values, but cleaner code.

---

## 💡 Why Keep OPTIONS?

### Question: Why not just use FACILITY_TYPES directly?

**Answer:** Backward compatibility and flexibility.

### Current Usage Pattern:

```javascript
// Pattern 1: Get labels
const labels = FACILITY_TYPE_OPTIONS.map(opt => opt.label)

// Pattern 2: Use in select
<option v-for="opt in FACILITY_TYPE_OPTIONS" :value="opt.value">
  {{ opt.label }}
</option>
```

### If We Removed OPTIONS:

```javascript
// Would need to change to:
const labels = Object.values(FACILITY_TYPES)

// Or:
<option v-for="type in Object.values(FACILITY_TYPES)" :value="type">
  {{ type }}
</option>
```

**Decision:** Keep OPTIONS for:
- Backward compatibility
- Consistent API with other constants (STATUS_OPTIONS, etc.)
- Easier to extend in future (e.g., add icons, colors)

---

## 🎯 Best Practices

### Single Source of Truth:

```javascript
// ✅ GOOD - Reference constant
{ value: FACILITY_TYPES.OTHER, label: FACILITY_TYPES.OTHER }

// ❌ BAD - Hardcoded duplicate
{ value: FACILITY_TYPES.OTHER, label: 'Khác' }
```

### Why?

**If we need to change "Khác" to "Khác (Other)":**

**Before (BAD):**
```javascript
export const FACILITY_TYPES = {
  OTHER: 'Khác (Other)'  // Changed here
}

export const FACILITY_TYPE_OPTIONS = [
  { value: FACILITY_TYPES.OTHER, label: 'Khác' }  // ❌ Forgot to update!
]
// Result: value='Khác (Other)', label='Khác' → Inconsistent!
```

**After (GOOD):**
```javascript
export const FACILITY_TYPES = {
  OTHER: 'Khác (Other)'  // Changed here
}

export const FACILITY_TYPE_OPTIONS = [
  { value: FACILITY_TYPES.OTHER, label: FACILITY_TYPES.OTHER }  // ✅ Auto-updates!
]
// Result: value='Khác (Other)', label='Khác (Other)' → Consistent!
```

---

## 🧪 Testing

### Test Cases:

**1. Display in Dropdown:**
```
1. Open facility create form
2. Check "Loại" dropdown
3. ✅ Should show: Bàn ghế, Máy móc, Dụng cụ, Điện tử, Khác
4. ✅ Default should be "Khác"
```

**2. Save Facility:**
```
1. Create facility with type "Máy móc"
2. Save
3. Check backend data
4. ✅ Should store: type: "Máy móc"
```

**3. Edit Facility:**
```
1. Edit facility with type "Bàn ghế"
2. Check dropdown
3. ✅ Should show "Bàn ghế" selected
```

**4. Custom Categories:**
```
1. Add custom category "Trang trí"
2. Check dropdown
3. ✅ Should show: Bàn ghế, Máy móc, ..., Khác, Trang trí
```

---

## 📁 Files Modified

**1. frontend/src/constants/facility.js**

**Changes:**
- Updated `FACILITY_TYPE_OPTIONS` to reference `FACILITY_TYPES` instead of hardcoded strings
- Added comment explaining purpose

**Lines changed:** ~6 lines

---

## ✅ Quality Check

- [x] Code implemented
- [x] No TypeScript/lint errors
- [x] No diagnostics issues
- [x] Backward compatible (no breaking changes)
- [x] Single source of truth
- [x] Matches backend constants
- [x] Easier to maintain
- [ ] Testing on device (pending user)

---

## 🎯 Summary

**Problem:** Redundant hardcoded strings in FACILITY_TYPE_OPTIONS

**Solution:** Reference FACILITY_TYPES constant instead

**Result:**
- ✅ Single source of truth
- ✅ No duplicates
- ✅ Easier to maintain
- ✅ Backward compatible
- ✅ Matches backend

**Impact:**
- No breaking changes
- Same runtime behavior
- Cleaner code
- Future-proof

---

**Implementation Date:** February 6, 2026  
**Developer:** Kiro AI  
**Status:** ✅ Complete, Ready for Testing
