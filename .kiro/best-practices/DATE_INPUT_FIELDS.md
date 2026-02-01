# Date Input Fields - Best Practices

## Problem
HTML date input fields have inconsistent rendering across browsers and devices. On mobile, they often:
- Display with larger default width than other input fields
- Break layout alignment
- Show browser-specific styling that doesn't match the design system

## Solution

### CSS Classes for Date Inputs
Always apply these classes to `<input type="date">` fields:

```html
<input 
  v-model="formData.date" 
  type="date" 
  class="w-full px-3 py-3 text-sm border border-gray-300 rounded-lg focus:ring-2 focus:ring-blue-500 appearance-none"
/>
```

### Key CSS Properties

| Property | Value | Purpose |
|----------|-------|---------|
| `appearance: none` | **CRITICAL** | Removes browser default styling, normalizes width |
| `px-3` | Padding | Reduced horizontal padding (not `px-4`) |
| `py-3` | Padding | Reduced vertical padding (not `py-4`) |
| `text-sm` | Font size | Smaller text size (not `text-base`) |
| `w-full` | Width | Full container width |
| `border` | Border | Standard border styling |
| `rounded-lg` | Border radius | Consistent with other inputs |
| `focus:ring-2` | Focus state | Accessibility and UX |

### Why These Values?

1. **`appearance: none`** - Strips browser defaults that cause width inconsistencies
2. **`px-3 py-3`** - Tighter spacing prevents overflow on mobile screens
3. **`text-sm`** - Smaller font ensures content fits within constrained width
4. **`w-full`** - Ensures responsive behavior in grid layouts

## Implementation Examples

### Single Date Field
```vue
<div>
  <label class="block text-sm font-medium text-gray-700 mb-3">Ngày mua</label>
  <input 
    v-model="formData.purchase_date" 
    type="date" 
    class="w-full px-3 py-3 text-sm border border-gray-300 rounded-lg focus:ring-2 focus:ring-blue-500 appearance-none"
  />
</div>
```

### Date Field in Grid Layout
```vue
<div class="grid grid-cols-1 sm:grid-cols-2 gap-4">
  <div>
    <label class="block text-sm font-medium text-gray-700 mb-3">Ngày mua</label>
    <input 
      v-model="formData.purchase_date" 
      type="date" 
      class="w-full px-3 py-3 text-sm border border-gray-300 rounded-lg focus:ring-2 focus:ring-blue-500 appearance-none"
    />
  </div>
  <div>
    <!-- Other field -->
  </div>
</div>
```

### Date Field in Form
```vue
<div class="grid grid-cols-2 gap-3">
  <div>
    <label class="block text-sm font-medium text-gray-700 mb-1">Ngày *</label>
    <input 
      v-model="formData.date" 
      type="date" 
      class="w-full px-3 py-3 text-sm border border-gray-300 rounded-lg focus:ring-2 focus:ring-blue-500 appearance-none"
    />
  </div>
  <div>
    <!-- Other field -->
  </div>
</div>
```

## Files Using This Pattern

- ✅ `frontend/src/views/FacilityAddEditView.vue`
- ✅ `frontend/src/views/ExpenseManagementView.vue`
- ✅ `frontend/src/views/FacilityManagementView.vue`
- ✅ `frontend/src/views/CashierReports.vue`

## Testing Checklist

- [ ] Date field displays correctly on mobile (iPhone, Android)
- [ ] Date field aligns with other input fields in the same row
- [ ] Date picker opens when clicked
- [ ] Date value is properly formatted (YYYY-MM-DD)
- [ ] Focus state shows blue ring
- [ ] No horizontal scrolling caused by date field

## Related Best Practices

- See `FORM_INPUTS.md` for general input field styling
- See `RESPONSIVE_LAYOUTS.md` for grid layout patterns
- See `MOBILE_OPTIMIZATION.md` for mobile-specific considerations
