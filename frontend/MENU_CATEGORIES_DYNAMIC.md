# Menu Categories - Dynamic API Implementation

## Summary
Removed hardcoded menu categories and implemented full API integration for dynamic category management.

## Changes Made

### Backend (Already Completed)
- ✅ Domain model: `backend/domain/menu/category.go`
- ✅ Repository: `backend/infrastructure/mongodb/menu_category_repository.go`
- ✅ Service: `backend/application/services/menu_category_service.go`
- ✅ HTTP Handler: `backend/interfaces/http/menu_category_handler.go`
- ✅ Routes wired in `backend/main.go` under `/manager/menu-categories`

### Frontend (Completed)
- ✅ Service: `frontend/src/services/menuCategory.js`
- ✅ Updated: `frontend/src/views/MenuView.vue`
  - Removed hardcoded category arrays
  - Added `menuCategories` ref for API data
  - Added `categoriesLoading` ref for loading state
  - Added `fetchCategories()` function to load from API
  - Updated `addCategory()` to call API with error handling
  - Updated `deleteCategory()` to call API with error handling
  - Integrated category fetching in `refreshData()` and `onMounted()`
  - Pull-to-refresh now refreshes both menu items and categories

## API Endpoints

### GET /manager/menu-categories
Returns all menu categories
```json
{
  "data": [
    {
      "id": "507f1f77bcf86cd799439011",
      "name": "Cà phê"
    }
  ]
}
```

### POST /manager/menu-categories
Create a new category
```json
{
  "name": "Trà"
}
```

### PUT /manager/menu-categories/:id
Update a category
```json
{
  "name": "Trà sữa"
}
```

### DELETE /manager/menu-categories/:id
Delete a category (only if no menu items use it)

## Features

### Category Management
- ✅ Load categories from API on mount
- ✅ Create new categories via API
- ✅ Delete categories via API (with validation)
- ✅ Prevent deletion if category has menu items
- ✅ Duplicate name validation
- ✅ Pull-to-refresh support
- ✅ Error handling with user-friendly messages

### User Experience
- Categories are fetched from backend on page load
- Pull-to-refresh updates both menu items and categories
- Creating a category immediately adds it to the list
- Deleting a category checks for menu items first
- All operations show success/error messages

## Testing Checklist

### Backend
- [x] Backend compiles successfully
- [ ] Test GET /manager/menu-categories
- [ ] Test POST /manager/menu-categories
- [ ] Test PUT /manager/menu-categories/:id
- [ ] Test DELETE /manager/menu-categories/:id
- [ ] Test delete validation (category with menu items)

### Frontend
- [x] Frontend builds successfully
- [ ] Categories load on page mount
- [ ] Pull-to-refresh updates categories
- [ ] Create category works
- [ ] Delete category works
- [ ] Delete validation prevents removing categories with items
- [ ] Error messages display correctly
- [ ] Category dropdown shows API categories

## Next Steps
1. Start backend server
2. Test category CRUD operations
3. Verify menu items can use new categories
4. Test edge cases (duplicate names, deletion with items)
