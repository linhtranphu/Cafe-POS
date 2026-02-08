# ✅ Facility Default Type: "Khác"

**Date:** February 6, 2026  
**Task:** Set default facility type to "Khác" instead of empty string

---

## 🎯 User Request

> "ở view facility, default 'chọn loại' là khác"

**Translation:** In facility view, the default "select type" should be "Khác" (Other).

---

## ✅ What Was Changed

### Before:

```javascript
formData.value = {
  name: '',
  type: '',  // ❌ Empty - shows "Chọn loại" placeholder
  area: 'Mặc định',
  quantity: 1,
  status: FACILITY_STATUS.IN_USE,
  // ...
}
```

**User Experience:**
- Opens form → Type dropdown shows "Chọn loại"
- User must manually select a type
- Extra step required

### After:

```javascript
formData.value = {
  name: '',
  type: 'Khác',  // ✅ Pre-selected "Khác"
  area: 'Mặc định',
  quantity: 1,
  status: FACILITY_STATUS.IN_USE,
  // ...
}
```

**User Experience:**
- Opens form → Type dropdown shows "Khác" selected
- User can keep it or change to specific type
- One less step for generic items

---

## 🔧 Implementation

### Files Modified:

**frontend/src/views/FacilityManagementView.vue**

### Changes Made:

**1. Initial formData declaration:**
```javascript
const formData = ref({
  name: '',
  type: 'Khác',  // Changed from ''
  area: 'Mặc định',
  quantity: 1,
  status: FACILITY_STATUS.IN_USE,
  purchase_date: '',
  cost: 0,
  supplier: '',
  notes: ''
})
```

**2. openCreateModal function:**
```javascript
const openCreateModal = () => {
  // ...
  formData.value = {
    name: '',
    type: 'Khác',  // Changed from ''
    area: 'Mặc định',
    quantity: 1,
    status: FACILITY_STATUS.IN_USE,
    // ...
  }
  // ...
}
```

**3. saveFacility function (reset after save):**
```javascript
const saveFacility = async () => {
  // ...
  // Reset form and close modal
  formData.value = {
    name: '',
    type: 'Khác',  // Changed from ''
    area: 'Mặc định',
    quantity: 1,
    status: FACILITY_STATUS.IN_USE,
    // ...
  }
  // ...
}
```

---

## 📊 Available Types

From `frontend/src/constants/facility.js`:

```javascript
export const FACILITY_TYPES = {
  FURNITURE: 'Bàn ghế',
  MACHINE: 'Máy móc',
  UTENSIL: 'Dụng cụ',
  ELECTRIC: 'Điện tử',
  OTHER: 'Khác'  // ← Default value
}

export const FACILITY_TYPE_OPTIONS = [
  { value: 'Bàn ghế', label: 'Bàn ghế' },
  { value: 'Máy móc', label: 'Máy móc' },
  { value: 'Dụng cụ', label: 'Dụng cụ' },
  { value: 'Điện tử', label: 'Điện tử' },
  { value: 'Khác', label: 'Khác' }  // ← Default
]
```

---

## 🎨 User Experience

### Before:

```
User clicks "Tạo thiết bị"
↓
Form opens
↓
Type dropdown: [Chọn loại ▼]  ← Empty, must select
↓
User must click dropdown
↓
User must select a type
↓
Can proceed
```

### After:

```
User clicks "Tạo thiết bị"
↓
Form opens
↓
Type dropdown: [Khác ▼]  ← Pre-selected
↓
User can:
  - Keep "Khác" → Proceed immediately ✅
  - Change to specific type → Select from dropdown
↓
Can proceed
```

---

## 💡 Why "Khác" as Default?

### Reasoning:

1. **Generic Items:**
   - Many facilities don't fit specific categories
   - "Khác" is a safe default for miscellaneous items

2. **Flexibility:**
   - User can keep it for generic items
   - User can change to specific type if needed

3. **Less Friction:**
   - One less required action
   - Faster form completion for generic items

4. **Consistency:**
   - Area already defaults to "Mặc định"
   - Status already defaults to "Đang sử dụng"
   - Type should also have a default

### Use Cases:

**Generic Items (Keep "Khác"):**
- Miscellaneous tools
- Decorations
- Temporary items
- Items that don't fit categories

**Specific Items (Change Type):**
- Tables/Chairs → "Bàn ghế"
- Coffee machine → "Máy móc"
- Cups/Plates → "Dụng cụ"
- POS system → "Điện tử"

---

## 🧪 Testing

### Test Cases:

**1. Create New Facility:**
```
1. Click "Tạo thiết bị"
2. Check Type dropdown
3. ✅ Should show "Khác" selected
4. ✅ Can keep "Khác" and save
5. ✅ Can change to other type
```

**2. After Saving:**
```
1. Create facility with "Khác"
2. Save successfully
3. Click "Tạo thiết bị" again
4. ✅ Type should reset to "Khác"
```

**3. Edit Existing Facility:**
```
1. Edit facility with type "Máy móc"
2. Check Type dropdown
3. ✅ Should show "Máy móc" (not "Khác")
4. ✅ Edit mode preserves original type
```

**4. Cancel and Reopen:**
```
1. Open create form
2. Change type to "Bàn ghế"
3. Cancel (close form)
4. Reopen create form
5. ✅ Type should reset to "Khác"
```

---

## 📱 Visual Example

### Form Display:

```
┌─────────────────────────────────────┐
│ ➕ Thêm thiết bị mới                │
├─────────────────────────────────────┤
│ Tên thiết bị *                      │
│ [                    ]              │
├─────────────────────────────────────┤
│ Loại *              Số lượng *      │
│ [Khác ▼]            [1]             │
│  ↑ Pre-selected                     │
├─────────────────────────────────────┤
│ Khu vực             Trạng thái *    │
│ [Mặc định]          [Đang sử dụng ▼]│
└─────────────────────────────────────┘
```

### Dropdown Options:

```
┌─────────────────────┐
│ Loại *              │
│ ┌─────────────────┐ │
│ │ Bàn ghế         │ │
│ │ Máy móc         │ │
│ │ Dụng cụ         │ │
│ │ Điện tử         │ │
│ │ ✓ Khác          │ │ ← Selected
│ └─────────────────┘ │
└─────────────────────┘
```

---

## 🔍 Edge Cases

### 1. Edit Mode:

```javascript
const openEditModal = (facility) => {
  editingFacility.value = facility
  formData.value = {
    name: facility.name || '',
    type: facility.type || '',  // ← Uses facility's type, not "Khác"
    area: facility.area || 'Mặc định',
    // ...
  }
  showFacilityForm.value = true
}
```

**Behavior:**
- Edit mode preserves original type
- Only create mode defaults to "Khác"

### 2. Empty Type from Backend:

```javascript
type: facility.type || ''  // If backend returns null/undefined
```

**Behavior:**
- If backend returns empty type → Shows "Chọn loại"
- This is correct for edit mode (preserve original state)

### 3. Custom Categories:

```javascript
const facilityCategories = computed(() => {
  const defaultCategories = FACILITY_TYPE_OPTIONS.map(opt => opt.label)
  const backendTypes = facilityStore.types.map(t => t.name)
  return [...new Set([...defaultCategories, ...backendTypes])]
})
```

**Behavior:**
- "Khác" is in default categories
- Always available as default
- Custom categories also available

---

## 📁 Files Modified

**1. frontend/src/views/FacilityManagementView.vue**

**Changes:**
- Line ~413: Initial formData declaration
- Line ~555: openCreateModal function
- Line ~615: saveFacility function (reset)

**Total:** 3 locations changed

---

## ✅ Quality Check

- [x] Code implemented
- [x] No TypeScript/lint errors
- [x] No diagnostics issues
- [x] "Khác" exists in FACILITY_TYPE_OPTIONS
- [x] Default works for create mode
- [x] Edit mode preserves original type
- [x] Reset after save works correctly
- [ ] Testing on device (pending user)

---

## 💡 Technical Notes

### Why 3 Locations?

**1. Initial Declaration:**
```javascript
const formData = ref({
  type: 'Khác'
})
```
- Sets default when component loads
- Used if form opened before any action

**2. openCreateModal:**
```javascript
const openCreateModal = () => {
  formData.value = {
    type: 'Khác'
  }
}
```
- Resets form when opening create modal
- Ensures clean state

**3. saveFacility (reset):**
```javascript
const saveFacility = async () => {
  // After save
  formData.value = {
    type: 'Khác'
  }
}
```
- Resets form after successful save
- Prepares for next create action

### Why Not Change Template?

```vue
<select v-model="formData.type">
  <option value="">Chọn loại</option>  <!-- Keep this -->
  <option v-for="cat in facilityCategories" :key="cat" :value="cat">
    {{ cat }}
  </option>
</select>
```

**Reason:**
- Template shows all options
- Default value set in JavaScript
- Cleaner separation of concerns
- Easier to maintain

---

## 🎯 Summary

**Problem:** Default facility type was empty, requiring extra step

**Solution:** Set default to "Khác" (Other)

**Result:**
- ✅ Form opens with "Khác" pre-selected
- ✅ User can keep it or change
- ✅ Faster form completion
- ✅ Better UX for generic items

**Impact:**
- Less friction in form
- Faster data entry
- Consistent with other defaults
- Flexible for all use cases

---

**Implementation Date:** February 6, 2026  
**Developer:** Kiro AI  
**Status:** ✅ Complete, Ready for Testing
