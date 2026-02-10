# Menu Categories API - Bug Fix

## Vấn đề
```
Failed to fetch categories: TypeError: Cannot read properties of null (reading 'data')
at fetchCategories (MenuView.vue:107:37)
```

## Nguyên nhân
Backend API responses không consistent:
- Một số endpoints trả về trực tiếp data: `categories`
- Frontend expect response có structure: `{ data: [...] }`

## Giải pháp

### Backend Changes
Đã cập nhật tất cả endpoints trong `backend/interfaces/http/menu_category_handler.go` để wrap response trong object `{ data: ... }`:

#### Trước:
```go
c.JSON(http.StatusOK, categories)
c.JSON(http.StatusCreated, category)
```

#### Sau:
```go
c.JSON(http.StatusOK, gin.H{"data": categories})
c.JSON(http.StatusCreated, gin.H{"data": category})
```

### Frontend Changes
Đã cập nhật `frontend/src/views/MenuView.vue` để xử lý response đúng:

#### Trước:
```javascript
const response = await menuCategoryService.getCategories()
menuCategories.value = response.data || []
```

#### Sau:
```javascript
const data = await menuCategoryService.getCategories()
// API returns { data: [...] }
menuCategories.value = data?.data || []
```

## Response Structure (Consistent)

### GET /api/manager/menu-categories
```json
{
  "data": [
    { "id": "...", "name": "Cà phê" }
  ]
}
```

### POST /api/manager/menu-categories
```json
{
  "data": {
    "id": "...",
    "name": "Sinh tố"
  }
}
```

### PUT /api/manager/menu-categories/:id
```json
{
  "data": {
    "id": "...",
    "name": "Trà sữa"
  }
}
```

### DELETE /api/manager/menu-categories/:id
```json
{
  "message": "category deleted successfully"
}
```

## Files Changed
- `backend/interfaces/http/menu_category_handler.go` - Wrapped all responses
- `frontend/src/views/MenuView.vue` - Updated fetchCategories to handle response correctly

## Testing
- [x] Backend compiles successfully
- [x] Frontend has no diagnostic errors
- [ ] Test API endpoints with actual server
- [ ] Verify categories load correctly in UI

## Status
✅ **Fixed**: Response structure now consistent across all endpoints
✅ **Compiled**: Both backend and frontend compile successfully
⏳ **Next**: Test with running server
