
## 📝 Checklist Phát Triển Page Mới

### Frontend Development:

**1. View Layer** (`/frontend/src/views/`)
- [ ] Tạo file `[Feature]View.vue`
- [ ] Import stores và services cần thiết
- [ ] Implement UI components (table, form, modal)
- [ ] Xử lý state management (loading, error)
- [ ] Implement CRUD operations
- [ ] Thêm validation cho forms
- [ ] Responsive design

**2. Service Layer** (`/frontend/src/services/`)
- [ ] Tạo file `[feature].js`
- [ ] Import api instance
- [ ] Implement GET methods (list, detail)
- [ ] Implement POST methods (create)
- [ ] Implement PUT methods (update)
- [ ] Implement DELETE methods (delete)
- [ ] Xử lý query parameters cho filtering
- [ ] Export service object

**3. Store Layer** (`/frontend/src/stores/`)
- [ ] Tạo file `[feature].js`
- [ ] Import service tương ứng
- [ ] Define state (data, loading, error)
- [ ] Implement fetch actions
- [ ] Implement create actions
- [ ] Implement update actions
- [ ] Implement delete actions
- [ ] Xử lý error handling
- [ ] Export store với defineStore

**4. Router** (`/frontend/src/router/index.js`)
- [ ] Thêm route mới
- [ ] Cấu hình meta (requiresAuth, role)
- [ ] Import view component

**5. Navigation** (`/frontend/src/components/Navigation.vue`)
- [ ] Thêm menu item mới
- [ ] Kiểm tra role-based visibility

### Backend Development:

**1. Domain Layer** (`/backend/domain/[feature]/`)
- [ ] Tạo thư mục feature
- [ ] Tạo file `[feature].go`
- [ ] Define structs với bson và json tags
- [ ] Thêm ObjectID, timestamps
- [ ] Define business entities

**2. Repository Layer** (`/backend/infrastructure/mongodb/`)
- [ ] Tạo file `[feature]_repository.go`
- [ ] Define repository struct với mongo.Collection
- [ ] Implement Create method
- [ ] Implement Get/Find methods
- [ ] Implement Update method
- [ ] Implement Delete method
- [ ] Xử lý context và errors
- [ ] Export NewRepository constructor

**3. Service Layer** (`/backend/application/services/`)
- [ ] Tạo file `[feature]_service.go`
- [ ] Define service struct với repository
- [ ] Implement business logic methods
- [ ] Validate input data
- [ ] Call repository methods
- [ ] Export NewService constructor

**4. Handler Layer** (`/backend/interfaces/http/`)
- [ ] Tạo file `[feature]_handler.go`
- [ ] Define handler struct với service
- [ ] Implement Create handler
- [ ] Implement Get/List handlers
- [ ] Implement Update handler
- [ ] Implement Delete handler
- [ ] Parse request body/params
- [ ] Return JSON responses
- [ ] Xử lý HTTP status codes
- [ ] Export NewHandler constructor

**5. Routes** (`/backend/main.go`)
- [ ] Import handler package
- [ ] Khởi tạo repository
- [ ] Khởi tạo service với repository
- [ ] Khởi tạo handler với service
- [ ] Thêm routes vào manager group
- [ ] Áp dụng middleware (auth, role)
- [ ] Test endpoints

### Testing & Validation:
- [ ] Test tất cả API endpoints với Postman/curl
- [ ] Kiểm tra validation rules
- [ ] Test error handling
- [ ] Kiểm tra role-based access
- [ ] Test UI interactions
- [ ] Verify data persistence
- [ ] Check responsive design

### Example Implementation:

```
Expense Management:
✅ Frontend:
  ✅ View: ExpenseView.vue
  ✅ Service: expense.js
  ✅ Store: expense.js
  ✅ Router: Added /expenses route
  ✅ Navigation: Added menu item

✅ Backend:
  ✅ Domain: domain/expense/expense.go
  ✅ Repository: infrastructure/mongodb/expense_repository.go
  ✅ Service: application/services/expense_service.go
  ✅ Handler: interfaces/http/expense_handler.go
  ✅ Routes: main.go (manager group)
```