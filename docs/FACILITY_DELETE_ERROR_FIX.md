# Facility Delete Error Fix ✅

## Problem
When trying to delete a facility that has maintenance history, the backend was returning a 500 Internal Server Error. The frontend was showing a generic error message without explaining why the deletion failed.

## Root Cause
1. **Backend**: The `DeleteFacility` handler was returning all errors as 500 (Internal Server Error), even business rule validation errors
2. **Frontend**: The error handling was not extracting the actual error message from the backend response

## Business Rule
The system has a business rule that prevents deletion of facilities that have maintenance history. This is to preserve data integrity and audit trails.

**Rule**: "không thể xóa tài sản đã có lịch sử bảo trì" (Cannot delete facility with maintenance history)

## Solution

### Backend Fix (`backend/interfaces/http/facility_handler.go`)
Modified the `DeleteFacility` handler to return proper HTTP status codes:
- **400 Bad Request**: For business rule validation errors (facility has maintenance history)
- **500 Internal Server Error**: For actual server errors

```go
func (h *FacilityHandler) DeleteFacility(c *gin.Context) {
	// ... validation code ...
	
	if err := h.service.DeleteFacility(c.Request.Context(), id, userID, username); err != nil {
		// Check if it's a business rule validation error
		if err.Error() == "không thể xóa tài sản đã có lịch sử bảo trì" {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "Xóa thành công"})
}
```

### Frontend Fix (`frontend/src/views/FacilityManagementView.vue`)
Improved error handling to extract and display the actual error message from the backend:

```javascript
const deleteFacility = async (facility) => {
  if (confirm(`Bạn có chắc muốn xóa thiết bị "${facility.name}"?`)) {
    try {
      await facilityStore.deleteFacility(facility.id)
      alert('Xóa thiết bị thành công')
    } catch (error) {
      console.error('Error deleting facility:', error)
      const errorMessage = error.response?.data?.error || 'Có lỗi xảy ra khi xóa thiết bị'
      alert(errorMessage)
    }
  }
}
```

## User Experience After Fix

### Scenario 1: Delete facility WITHOUT maintenance history
- ✅ Deletion succeeds
- ✅ Shows success message: "Xóa thiết bị thành công"
- ✅ Facility is removed from the list

### Scenario 2: Delete facility WITH maintenance history
- ❌ Deletion fails (as expected by business rule)
- ✅ Shows clear error message: "không thể xóa tài sản đã có lịch sử bảo trì"
- ✅ User understands why deletion failed
- ✅ Returns 400 Bad Request (not 500 Internal Server Error)

## Files Modified
1. `backend/interfaces/http/facility_handler.go` - Improved error handling with proper HTTP status codes
2. `frontend/src/views/FacilityManagementView.vue` - Extract and display backend error messages

## Testing
After this fix:
- ✅ Backend returns 400 (not 500) for business rule violations
- ✅ Frontend displays the actual error message from backend
- ✅ User understands why deletion failed
- ✅ Success message shown when deletion succeeds

## Future Improvements (Optional)
Consider implementing one of these approaches:
1. **Soft Delete**: Mark facility as deleted instead of removing it
2. **Archive**: Move facility to archived state while preserving history
3. **Force Delete**: Add a "force delete" option for managers with confirmation
4. **Better UI**: Show maintenance history count before attempting delete

## Status: COMPLETE ✅
The error is now properly handled with appropriate HTTP status codes and clear error messages.
