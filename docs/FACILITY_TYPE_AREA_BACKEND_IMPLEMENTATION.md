# Facility Type & Area Backend Implementation

**Date**: January 31, 2026  
**Status**: ✅ COMPLETE

## Overview

Implemented backend API for managing facility types and areas, replacing localStorage with proper database storage.

---

## Backend Implementation

### 1. Domain Models ✅

**File**: `backend/domain/facility/facility.go`

**Added**:
```go
// FacilityType represents a facility type/category
type FacilityType struct {
    ID        primitive.ObjectID `json:"id" bson:"_id,omitempty"`
    Name      string             `json:"name" bson:"name"`
    CreatedAt time.Time          `json:"created_at" bson:"created_at"`
}

// FacilityArea represents a facility area/location
type FacilityArea struct {
    ID        primitive.ObjectID `json:"id" bson:"_id,omitempty"`
    Name      string             `json:"name" bson:"name"`
    CreatedAt time.Time          `json:"created_at" bson:"created_at"`
}

// Helper functions
func GetDefaultFacilityTypes() []string
func GetDefaultFacilityAreas() []string
```

**Default Types**:
- Bàn ghế (Furniture)
- Máy móc (Machine)
- Dụng cụ (Utensil)
- Điện tử (Electric)
- Khác (Other)

**Default Areas**:
- Phòng khách (Dining Room)
- Bếp (Kitchen)
- Quầy bar (Counter)
- Kho (Storage)
- Văn phòng (Office)
- Khác (Other)

---

### 2. Repository Layer ✅

**File**: `backend/infrastructure/mongodb/facility_repository.go`

**Added Methods**:
```go
// FacilityType management
CreateFacilityType(ctx, *facility.FacilityType) error
GetFacilityTypes(ctx) ([]facility.FacilityType, error)
DeleteFacilityType(ctx, primitive.ObjectID) error

// FacilityArea management
CreateFacilityArea(ctx, *facility.FacilityArea) error
GetFacilityAreas(ctx) ([]facility.FacilityArea, error)
DeleteFacilityArea(ctx, primitive.ObjectID) error
```

**Collections**:
- `facility_types` - Stores facility types
- `facility_areas` - Stores facility areas

---

### 3. Service Layer ✅

**File**: `backend/application/services/facility_service.go`

**Added Methods**:
```go
// FacilityType management
CreateFacilityType(ctx, name string) (*facility.FacilityType, error)
GetFacilityTypes(ctx) ([]facility.FacilityType, error)
DeleteFacilityType(ctx, id primitive.ObjectID) error

// FacilityArea management
CreateFacilityArea(ctx, name string) (*facility.FacilityArea, error)
GetFacilityAreas(ctx) ([]facility.FacilityArea, error)
DeleteFacilityArea(ctx, id primitive.ObjectID) error
```

---

### 4. HTTP Handlers ✅

**File**: `backend/interfaces/http/facility_handler.go`

**Added Endpoints**:

**FacilityType**:
- `POST /api/manager/facility-types` - Create type
- `GET /api/manager/facility-types` - Get all types
- `DELETE /api/manager/facility-types/:id` - Delete type

**FacilityArea**:
- `POST /api/manager/facility-areas` - Create area
- `GET /api/manager/facility-areas` - Get all areas
- `DELETE /api/manager/facility-areas/:id` - Delete area

**Request/Response**:
```json
// Create Type/Area
POST /api/manager/facility-types
{
  "name": "Thiết bị điện tử"
}

// Response
{
  "id": "507f1f77bcf86cd799439011",
  "name": "Thiết bị điện tử",
  "created_at": "2026-01-31T10:00:00Z"
}

// Get All Types
GET /api/manager/facility-types
[
  {
    "id": "507f1f77bcf86cd799439011",
    "name": "Bàn ghế",
    "created_at": "2026-01-31T10:00:00Z"
  },
  ...
]

// Delete Type
DELETE /api/manager/facility-types/:id
{
  "message": "Xóa thành công"
}
```

---

### 5. Routes ✅

**File**: `backend/main.go`

**Added Routes**:
```go
// Facility type and area routes
manager.POST("/facility-types", facilityHandler.CreateFacilityType)
manager.GET("/facility-types", facilityHandler.GetFacilityTypes)
manager.DELETE("/facility-types/:id", facilityHandler.DeleteFacilityType)
manager.POST("/facility-areas", facilityHandler.CreateFacilityArea)
manager.GET("/facility-areas", facilityHandler.GetFacilityAreas)
manager.DELETE("/facility-areas/:id", facilityHandler.DeleteFacilityArea)
```

---

## Frontend Implementation

### 1. Service Layer ✅

**File**: `frontend/src/services/facility.js`

**Added Methods**:
```javascript
// Facility Type management
async getFacilityTypes()
async createFacilityType(type)
async deleteFacilityType(id)

// Facility Area management
async getFacilityAreas()
async createFacilityArea(area)
async deleteFacilityArea(id)
```

---

### 2. Store Layer ✅

**File**: `frontend/src/stores/facility.js`

**Added State**:
```javascript
state: () => ({
  items: [],
  types: [],  // NEW
  areas: [],  // NEW
  loading: false,
  error: null
})
```

**Added Actions**:
```javascript
// Facility Type management
async fetchFacilityTypes()
async createFacilityType(type)
async deleteFacilityType(id)

// Facility Area management
async fetchFacilityAreas()
async createFacilityArea(area)
async deleteFacilityArea(id)
```

---

### 3. View Layer ✅

**File**: `frontend/src/views/FacilityManagementView.vue`

**Changes**:
- Replaced localStorage with backend API calls
- `facilityCategories` computed property now uses `facilityStore.types`
- `addCategory()` calls `facilityStore.createFacilityType()`
- `deleteCategory()` calls `facilityStore.deleteFacilityType()`
- `onMounted()` calls `facilityStore.fetchFacilityTypes()`

**Before** (localStorage):
```javascript
const customCategories = JSON.parse(localStorage.getItem('facilityCategories') || '[]')
```

**After** (Backend API):
```javascript
const backendTypes = facilityStore.types.map(t => t.name)
```

---

## Data Migration

### From localStorage to Backend

**Migration Steps**:
1. User opens Facility Management view
2. Frontend fetches types from backend
3. If backend is empty, default types are available from constants
4. User can add custom types via UI
5. Custom types are saved to backend database

**No Data Loss**:
- Default types always available from constants
- Custom types can be re-added via UI
- No automatic migration needed

---

## Testing

### Backend Tests

**Compilation**: ✅ PASS
```bash
go build ./...
# Exit Code: 0
```

**Manual Testing Needed**:
- [ ] Create facility type
- [ ] Get all facility types
- [ ] Delete facility type
- [ ] Create facility area
- [ ] Get all facility areas
- [ ] Delete facility area
- [ ] Cannot delete type in use
- [ ] Cannot delete area in use

---

### Frontend Tests

**Build**: ✅ PASS
```bash
npm run build
# ✓ built in 3.50s
```

**Manual Testing Needed**:
- [ ] Fetch types on mount
- [ ] Display types in dropdown
- [ ] Add new type via modal
- [ ] Delete custom type
- [ ] Cannot delete default type
- [ ] Cannot delete type in use
- [ ] Types persist after page reload

---

## API Documentation

### Facility Types

#### Create Facility Type
```http
POST /api/manager/facility-types
Authorization: Bearer <token>
Content-Type: application/json

{
  "name": "Thiết bị điện tử"
}

Response: 201 Created
{
  "id": "507f1f77bcf86cd799439011",
  "name": "Thiết bị điện tử",
  "created_at": "2026-01-31T10:00:00Z"
}
```

#### Get All Facility Types
```http
GET /api/manager/facility-types
Authorization: Bearer <token>

Response: 200 OK
[
  {
    "id": "507f1f77bcf86cd799439011",
    "name": "Bàn ghế",
    "created_at": "2026-01-31T10:00:00Z"
  },
  {
    "id": "507f1f77bcf86cd799439012",
    "name": "Máy móc",
    "created_at": "2026-01-31T10:00:00Z"
  }
]
```

#### Delete Facility Type
```http
DELETE /api/manager/facility-types/:id
Authorization: Bearer <token>

Response: 200 OK
{
  "message": "Xóa thành công"
}

Error: 400 Bad Request
{
  "error": "Không thể xóa loại thiết bị đang được sử dụng"
}
```

### Facility Areas

Same structure as Facility Types, replace `/facility-types` with `/facility-areas`.

---

## Database Schema

### facility_types Collection
```javascript
{
  _id: ObjectId,
  name: String,
  created_at: ISODate
}
```

**Indexes**:
- `_id` (primary key)
- `name` (unique, for preventing duplicates)

### facility_areas Collection
```javascript
{
  _id: ObjectId,
  name: String,
  created_at: ISODate
}
```

**Indexes**:
- `_id` (primary key)
- `name` (unique, for preventing duplicates)

---

## Benefits

✅ **Centralized Data**: Types and areas stored in database  
✅ **Multi-User**: Changes visible to all users  
✅ **Persistent**: Data survives browser cache clear  
✅ **Scalable**: Easy to add more fields (description, icon, etc.)  
✅ **Consistent**: Same pattern as ingredient categories  
✅ **API-First**: Ready for mobile apps or external integrations

---

## Future Enhancements (Optional)

- Add type/area descriptions
- Add icons/colors for types
- Add type/area ordering
- Add type/area usage statistics
- Add bulk import/export
- Add type/area permissions
- Add type/area archiving (soft delete)

---

## Files Modified

### Backend (5 files)
1. `backend/domain/facility/facility.go` - Added models
2. `backend/infrastructure/mongodb/facility_repository.go` - Added methods
3. `backend/application/services/facility_service.go` - Added methods
4. `backend/interfaces/http/facility_handler.go` - Added handlers
5. `backend/main.go` - Added routes

### Frontend (3 files)
1. `frontend/src/services/facility.js` - Added API calls
2. `frontend/src/stores/facility.js` - Added state & actions
3. `frontend/src/views/FacilityManagementView.vue` - Updated to use backend

---

**Status**: ✅ COMPLETE  
**Build**: ✅ SUCCESS (Backend & Frontend)  
**Ready for**: Testing and deployment
