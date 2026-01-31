# Facility Constants Synchronization ✅

## Problem
Frontend and backend were using different status values, causing facilities to not be created properly:
- Frontend was using: `in_use`, `maintenance`, `broken`, `retired`
- Backend was expecting: `Đang sử dụng`, `Đang sửa`, `Hỏng`, `Ngừng sử dụng`

This mismatch caused validation errors and data inconsistency.

## Solution
Created a centralized constants file to ensure frontend and backend are synchronized.

## Files Created

### 1. `frontend/src/constants/facility.js`
Centralized constants file that matches backend definitions in `backend/domain/facility/facility.go`

#### Constants Defined:

**Facility Status:**
- `IN_USE`: 'Đang sử dụng'
- `BROKEN`: 'Hỏng'
- `REPAIRING`: 'Đang sửa'
- `INACTIVE`: 'Ngừng sử dụng'
- `DISPOSED`: 'Thanh lý'

**Facility Types:**
- `FURNITURE`: 'Bàn ghế'
- `MACHINE`: 'Máy móc'
- `UTENSIL`: 'Dụng cụ'
- `ELECTRIC`: 'Điện tử'
- `OTHER`: 'Khác'

**Facility Areas:**
- `DINING_ROOM`: 'Phòng khách'
- `KITCHEN`: 'Bếp'
- `COUNTER`: 'Quầy bar'
- `STORAGE`: 'Kho'
- `OFFICE`: 'Văn phòng'
- `OTHER`: 'Khác'

**Maintenance Types:**
- `SCHEDULED`: 'scheduled'
- `EMERGENCY`: 'emergency'
- `PREVENTIVE`: 'preventive'
- `CORRECTIVE`: 'corrective'

**Issue Severity:**
- `LOW`: 'low'
- `MEDIUM`: 'medium'
- `HIGH`: 'high'
- `CRITICAL`: 'critical'

**Issue Status:**
- `OPEN`: 'open'
- `IN_PROGRESS`: 'in_progress'
- `RESOLVED`: 'resolved'

#### Helper Functions:
```javascript
getFacilityStatusClass(status)  // Returns Tailwind CSS classes
getIssueStatusClass(status)     // Returns Tailwind CSS classes
getIssueSeverityClass(severity) // Returns Tailwind CSS classes
```

## Files Modified

### 1. `frontend/src/views/FacilityManagementView.vue`
Updated to use constants:
- Import constants from `../constants/facility`
- Use `FACILITY_STATUS.IN_USE` instead of hardcoded strings
- Use `FACILITY_STATUS_OPTIONS` for select dropdown
- Use helper functions for CSS classes
- Updated computed properties to use constants

### Changes Made:
```javascript
// Before
status: 'in_use'

// After
import { FACILITY_STATUS, FACILITY_STATUS_OPTIONS } from '../constants/facility'
status: FACILITY_STATUS.IN_USE
```

## Benefits

### 1. Single Source of Truth
- All status values defined in one place
- Easy to update across entire application
- Reduces bugs from typos or mismatches

### 2. Type Safety
- Constants prevent typos
- IDE autocomplete support
- Easier refactoring

### 3. Maintainability
- Clear documentation of all possible values
- Easy to add new statuses
- Consistent naming conventions

### 4. Frontend-Backend Sync
- Frontend constants match backend exactly
- No more validation errors from mismatched values
- Data consistency guaranteed

## Backend Constants Reference
From `backend/domain/facility/facility.go`:
```go
const (
	StatusInUse     = "Đang sử dụng"
	StatusBroken    = "Hỏng"
	StatusRepairing = "Đang sửa"
	StatusInactive  = "Ngừng sử dụng"
	StatusDisposed  = "Thanh lý"
)
```

## Usage Example

### In Vue Component:
```vue
<script setup>
import { FACILITY_STATUS, FACILITY_STATUS_OPTIONS, getFacilityStatusClass } from '../constants/facility'

const formData = ref({
  status: FACILITY_STATUS.IN_USE
})
</script>

<template>
  <select v-model="formData.status">
    <option v-for="option in FACILITY_STATUS_OPTIONS" :key="option.value" :value="option.value">
      {{ option.label }}
    </option>
  </select>
  
  <span :class="getFacilityStatusClass(facility.status)">
    {{ facility.status }}
  </span>
</template>
```

## Testing
After this fix:
- ✅ Facilities can be created from frontend
- ✅ Status values match backend expectations
- ✅ No validation errors
- ✅ Data consistency maintained
- ✅ CSS classes applied correctly

## Future Improvements
1. Consider creating TypeScript types for better type safety
2. Add validation functions in constants file
3. Create similar constants files for other domains (orders, shifts, etc.)
4. Generate constants from backend OpenAPI spec

## Status: COMPLETE ✅
Frontend and backend are now synchronized using centralized constants.
