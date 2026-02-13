# Tasks: Quản Lý Nguyên Liệu Batch

## Phase 1: Backend Foundation

### 1. Domain Layer Implementation

- [ ] 1.1 Tạo domain entities
  - [ ] 1.1.1 Tạo BatchDefinition entity với ConversionRate value object
  - [ ] 1.1.2 Tạo BatchRecord entity với IngredientUsage value object
  - [ ] 1.1.3 Tạo BatchUsageLog entity
  - [ ] 1.1.4 Tạo Alert entities (LowStockAlert, ExpiringAlert, ExpiredAlert)
  - [ ] 1.1.5 Viết unit tests cho domain entities

### 2. Repository Layer Implementation

- [ ] 2.1 Implement BatchDefinitionRepository
  - [ ] 2.1.1 Tạo interface và MongoDB implementation
  - [ ] 2.1.2 Implement CRUD operations
  - [ ] 2.1.3 Implement FindAll với filters
  - [ ] 2.1.4 Viết unit tests với mock MongoDB

- [ ] 2.2 Implement BatchRecordRepository
  - [ ] 2.2.1 Tạo interface và MongoDB implementation
  - [ ] 2.2.2 Implement CRUD operations
  - [ ] 2.2.3 Implement FindAvailableByDefinition (FIFO query)
  - [ ] 2.2.4 Implement GetTotalAvailableQuantity
  - [ ] 2.2.5 Implement UpdateQuantity với concurrency control
  - [ ] 2.2.6 Viết unit tests với mock MongoDB

- [ ] 2.3 Implement BatchUsageLogRepository
  - [ ] 2.3.1 Tạo interface và MongoDB implementation
  - [ ] 2.3.2 Implement Create và FindAll với filters
  - [ ] 2.3.3 Viết unit tests

- [ ] 2.4 Tạo MongoDB indexes
  - [ ] 2.4.1 Tạo indexes cho batch_definitions collection
  - [ ] 2.4.2 Tạo indexes cho batch_records collection
  - [ ] 2.4.3 Tạo indexes cho batch_usage_logs collection
  - [ ] 2.4.4 Viết migration script

### 3. Service Layer Implementation

- [ ] 3.1 Implement BatchDefinitionService
  - [ ] 3.1.1 Implement Create với validation
  - [ ] 3.1.2 Implement Update, Delete, GetByID, List
  - [ ] 3.1.3 Implement ValidateConversionRates (check ingredients exist)
  - [ ] 3.1.4 Viết unit tests với mocked repositories

- [ ] 3.2 Implement BatchCostCalculator
  - [ ] 3.2.1 Implement CalculateBatchCost với wastage calculation
  - [ ] 3.2.2 Integrate với Cost_Calculator service hiện có
  - [ ] 3.2.3 Implement cost caching
  - [ ] 3.2.4 Viết unit tests
  - [ ] 3.2.5 **Property Test: Cost Accuracy** - Verify chi phí được tính chính xác

- [ ] 3.3 Implement BatchRecordService
  - [ ] 3.3.1 Implement CreateBatch với transaction
  - [ ] 3.3.2 Implement ingredient deduction logic
  - [ ] 3.3.3 Implement expiry calculation
  - [ ] 3.3.4 Implement GetByID, List với filters
  - [ ] 3.3.5 Implement UpdateQuantity, MarkAsExpired, Delete
  - [ ] 3.3.6 Implement GetAvailableBatches
  - [ ] 3.3.7 Viết unit tests
  - [ ] 3.3.8 **Property Test: Inventory Conservation** - Verify tồn kho được bảo toàn


- [ ] 3.4 Implement BatchUsageService
  - [ ] 3.4.1 Implement UseBatch với FIFO logic
  - [ ] 3.4.2 Implement multi-batch usage (khi 1 batch không đủ)
  - [ ] 3.4.3 Implement expiry checking trước khi use
  - [ ] 3.4.4 Implement transaction cho batch deduction
  - [ ] 3.4.5 Implement GetUsageHistory
  - [ ] 3.4.6 Viết unit tests
  - [ ] 3.4.7 **Property Test: FIFO Ordering** - Verify batch cũ nhất được dùng trước
  - [ ] 3.4.8 **Property Test: Expiry Enforcement** - Verify batch hết hạn không được dùng

- [ ] 3.5 Implement BatchAlertService
  - [ ] 3.5.1 Implement CheckLowStock
  - [ ] 3.5.2 Implement CheckExpiring
  - [ ] 3.5.3 Implement CheckExpired
  - [ ] 3.5.4 Implement GetAlerts (aggregate all alerts)
  - [ ] 3.5.5 Implement alert caching
  - [ ] 3.5.6 Viết unit tests
  - [ ] 3.5.7 **Property Test: Alert Correctness** - Verify alerts hiển thị đúng

- [ ] 3.6 Implement Background Jobs
  - [ ] 3.6.1 Tạo job để mark expired batches (chạy mỗi giờ)
  - [ ] 3.6.2 Tạo job để check alerts (chạy mỗi 5 phút)
  - [ ] 3.6.3 Implement job scheduler
  - [ ] 3.6.4 Viết tests cho background jobs

### 4. HTTP Handler Layer Implementation

- [ ] 4.1 Implement BatchDefinitionHandler
  - [ ] 4.1.1 Implement POST /api/batch-definitions
  - [ ] 4.1.2 Implement GET /api/batch-definitions
  - [ ] 4.1.3 Implement GET /api/batch-definitions/:id
  - [ ] 4.1.4 Implement PUT /api/batch-definitions/:id
  - [ ] 4.1.5 Implement DELETE /api/batch-definitions/:id
  - [ ] 4.1.6 Add authentication & authorization middleware
  - [ ] 4.1.7 Add input validation
  - [ ] 4.1.8 Viết integration tests

- [ ] 4.2 Implement BatchRecordHandler
  - [ ] 4.2.1 Implement POST /api/batch-records
  - [ ] 4.2.2 Implement GET /api/batch-records
  - [ ] 4.2.3 Implement GET /api/batch-records/:id
  - [ ] 4.2.4 Implement PATCH /api/batch-records/:id/quantity
  - [ ] 4.2.5 Implement PATCH /api/batch-records/:id/expire
  - [ ] 4.2.6 Implement DELETE /api/batch-records/:id
  - [ ] 4.2.7 Add authentication & authorization
  - [ ] 4.2.8 Viết integration tests

- [ ] 4.3 Implement BatchUsageHandler
  - [ ] 4.3.1 Implement POST /api/batch-usage
  - [ ] 4.3.2 Implement GET /api/batch-usage/history
  - [ ] 4.3.3 Add authentication & authorization
  - [ ] 4.3.4 Viết integration tests

- [ ] 4.4 Implement BatchAlertHandler
  - [ ] 4.4.1 Implement GET /api/batch-alerts
  - [ ] 4.4.2 Add authentication & authorization
  - [ ] 4.4.3 Viết integration tests

- [ ] 4.5 Implement BatchReportHandler
  - [ ] 4.5.1 Implement GET /api/batch-reports/production
  - [ ] 4.5.2 Implement GET /api/batch-reports/wastage
  - [ ] 4.5.3 Implement GET /api/batch-reports/usage
  - [ ] 4.5.4 Add authentication & authorization (manager only)
  - [ ] 4.5.5 Viết integration tests

### 5. Integration với Existing Systems

- [ ] 5.1 Integrate với Inventory System
  - [ ] 5.1.1 Update ingredient deduction logic để support batch creation
  - [ ] 5.1.2 Add rollback logic khi batch creation fails
  - [ ] 5.1.3 Viết integration tests

- [ ] 5.2 Integrate với Menu System
  - [ ] 5.2.1 Extend Recipe schema để support batch ingredients
  - [ ] 5.2.2 Update recipe validation để check batch availability
  - [ ] 5.2.3 Update cost calculation để support batch costs
  - [ ] 5.2.4 Viết integration tests

- [ ] 5.3 Integrate với Order System
  - [ ] 5.3.1 Update order processing để deduct batches
  - [ ] 5.3.2 Update order cost calculation để use batch costs
  - [ ] 5.3.3 Add rollback logic khi order fails
  - [ ] 5.3.4 Viết integration tests

### 6. Property-Based Testing

- [ ] 6.1 **Property Test: Transaction Atomicity**
  - [ ] 6.1.1 Test batch creation transaction rollback
  - [ ] 6.1.2 Test batch usage transaction rollback
  - [ ] 6.1.3 Test concurrent transactions

- [ ] 6.2 **Property Test: Quantity Non-Negativity**
  - [ ] 6.2.1 Test ingredient quantity never goes negative
  - [ ] 6.2.2 Test batch quantity never goes negative
  - [ ] 6.2.3 Test concurrent operations

### 7. Backend Testing & Documentation

- [ ] 7.1 Viết comprehensive integration tests
  - [ ] 7.1.1 Test full batch creation flow
  - [ ] 7.1.2 Test full batch usage flow
  - [ ] 7.1.3 Test alert generation flow
  - [ ] 7.1.4 Test report generation flow

- [ ] 7.2 Viết API documentation
  - [ ] 7.2.1 Document tất cả endpoints với examples
  - [ ] 7.2.2 Document error responses
  - [ ] 7.2.3 Document authentication requirements

- [ ] 7.3 Performance testing
  - [ ] 7.3.1 Test với 1000+ batch records
  - [ ] 7.3.2 Test concurrent batch creation
  - [ ] 7.3.3 Test concurrent batch usage
  - [ ] 7.3.4 Optimize queries nếu cần

## Phase 2: Frontend Implementation

### 8. State Management (Pinia Stores)

- [ ] 8.1 Implement useBatchDefinitionStore
  - [ ] 8.1.1 Define state, getters, actions
  - [ ] 8.1.2 Implement fetchDefinitions
  - [ ] 8.1.3 Implement createDefinition, updateDefinition, deleteDefinition
  - [ ] 8.1.4 Add error handling
  - [ ] 8.1.5 Viết unit tests

- [ ] 8.2 Implement useBatchRecordStore
  - [ ] 8.2.1 Define state với filters và pagination
  - [ ] 8.2.2 Implement fetchRecords với filters
  - [ ] 8.2.3 Implement createRecord, updateQuantity, markAsExpired, deleteRecord
  - [ ] 8.2.4 Add error handling
  - [ ] 8.2.5 Viết unit tests

- [ ] 8.3 Implement useBatchAlertStore
  - [ ] 8.3.1 Define state
  - [ ] 8.3.2 Implement fetchAlerts
  - [ ] 8.3.3 Implement auto-refresh logic
  - [ ] 8.3.4 Add error handling
  - [ ] 8.3.5 Viết unit tests

- [ ] 8.4 Implement useBatchReportStore
  - [ ] 8.4.1 Define state
  - [ ] 8.4.2 Implement fetchProductionReport, fetchWastageReport, fetchUsageReport
  - [ ] 8.4.3 Implement exportReport
  - [ ] 8.4.4 Add error handling
  - [ ] 8.4.5 Viết unit tests

### 9. API Service Layer

- [ ] 9.1 Tạo batchDefinitionService.js
  - [ ] 9.1.1 Implement API calls cho CRUD operations
  - [ ] 9.1.2 Add error handling và retry logic

- [ ] 9.2 Tạo batchRecordService.js
  - [ ] 9.2.1 Implement API calls cho batch record operations
  - [ ] 9.2.2 Add error handling

- [ ] 9.3 Tạo batchAlertService.js
  - [ ] 9.3.1 Implement API call cho alerts
  - [ ] 9.3.2 Add polling logic

- [ ] 9.4 Tạo batchReportService.js
  - [ ] 9.4.1 Implement API calls cho reports
  - [ ] 9.4.2 Implement export functionality

### 10. Batch Definition Components

- [ ] 10.1 Implement BatchDefinitionList.vue
  - [ ] 10.1.1 Tạo table view với columns
  - [ ] 10.1.2 Add search functionality
  - [ ] 10.1.3 Add create button
  - [ ] 10.1.4 Add edit/delete actions
  - [ ] 10.1.5 Add responsive design
  - [ ] 10.1.6 Viết component tests

- [ ] 10.2 Implement BatchDefinitionForm.vue
  - [ ] 10.2.1 Tạo form với validation
  - [ ] 10.2.2 Implement ingredient selector (autocomplete)
  - [ ] 10.2.3 Implement dynamic conversion rates list
  - [ ] 10.2.4 Add cost preview
  - [ ] 10.2.5 Add responsive design
  - [ ] 10.2.6 Viết component tests

### 11. Batch Record Components

- [ ] 11.1 Implement BatchRecordList.vue
  - [ ] 11.1.1 Tạo table view với color coding
  - [ ] 11.1.2 Add filters (batch type, status, date range, preparer)
  - [ ] 11.1.3 Add sorting (expiry date, prepared date, quantity)
  - [ ] 11.1.4 Add pagination
  - [ ] 11.1.5 Add actions (view, mark expired, delete)
  - [ ] 11.1.6 Add responsive design
  - [ ] 11.1.7 Viết component tests

- [ ] 11.2 Implement BatchRecordForm.vue
  - [ ] 11.2.1 Tạo form với batch definition selector
  - [ ] 11.2.2 Add quantity input
  - [ ] 11.2.3 Display required ingredients và expected cost
  - [ ] 11.2.4 Add confirmation dialog
  - [ ] 11.2.5 Add error handling (insufficient ingredients)
  - [ ] 11.2.6 Add responsive design
  - [ ] 11.2.7 Viết component tests

- [ ] 11.3 Implement BatchRecordDetail.vue
  - [ ] 11.3.1 Display batch record information
  - [ ] 11.3.2 Display ingredients used breakdown
  - [ ] 11.3.3 Display cost breakdown
  - [ ] 11.3.4 Display usage history
  - [ ] 11.3.5 Add timeline visualization
  - [ ] 11.3.6 Add actions (mark expired, delete)
  - [ ] 11.3.7 Viết component tests

### 12. Alert Components

- [ ] 12.1 Implement BatchAlertPanel.vue
  - [ ] 12.1.1 Tạo three sections (low stock, expiring, expired)
  - [ ] 12.1.2 Add badge với số lượng alerts
  - [ ] 12.1.3 Implement expandable sections
  - [ ] 12.1.4 Add auto-refresh (every 5 minutes)
  - [ ] 12.1.5 Add responsive design
  - [ ] 12.1.6 Viết component tests

- [ ] 12.2 Implement BatchAlertCard.vue
  - [ ] 12.2.1 Tạo card với icon và color coding
  - [ ] 12.2.2 Display alert information
  - [ ] 12.2.3 Add action buttons
  - [ ] 12.2.4 Add dismiss functionality
  - [ ] 12.2.5 Viết component tests

### 13. Report Components

- [ ] 13.1 Implement BatchProductionReport.vue
  - [ ] 13.1.1 Add date range picker
  - [ ] 13.1.2 Add filters (batch type, preparer)
  - [ ] 13.1.3 Display summary cards
  - [ ] 13.1.4 Add production trend chart
  - [ ] 13.1.5 Add breakdown table
  - [ ] 13.1.6 Add export to CSV
  - [ ] 13.1.7 Viết component tests

- [ ] 13.2 Implement BatchWastageReport.vue
  - [ ] 13.2.1 Add date range picker và filters
  - [ ] 13.2.2 Display summary cards
  - [ ] 13.2.3 Add wastage trend chart
  - [ ] 13.2.4 Add breakdown table
  - [ ] 13.2.5 Add recommendations section
  - [ ] 13.2.6 Viết component tests

- [ ] 13.3 Implement BatchUsageReport.vue
  - [ ] 13.3.1 Add date range picker và filters
  - [ ] 13.3.2 Display summary cards
  - [ ] 13.3.3 Add usage trend chart
  - [ ] 13.3.4 Add breakdown by menu item table
  - [ ] 13.3.5 Add most used batches ranking
  - [ ] 13.3.6 Viết component tests

### 14. Integration Components

- [ ] 14.1 Enhance MenuRecipeEditor.vue
  - [ ] 14.1.1 Add toggle để chọn "Nguyên Liệu Thô" vs "Batch"
  - [ ] 14.1.2 Add batch selector
  - [ ] 14.1.3 Display available batch quantity
  - [ ] 14.1.4 Add warning nếu batch không đủ
  - [ ] 14.1.5 Update cost calculation
  - [ ] 14.1.6 Viết component tests

- [ ] 14.2 Implement BatchStatusWidget.vue
  - [ ] 14.2.1 Display summary (total batches, available quantity)
  - [ ] 14.2.2 Display alert count badges
  - [ ] 14.2.3 Add quick links
  - [ ] 14.2.4 Add mini usage trend chart
  - [ ] 14.2.5 Add compact mode
  - [ ] 14.2.6 Viết component tests

### 15. Routing & Navigation

- [ ] 15.1 Add batch routes
  - [ ] 15.1.1 Define routes cho batch management
  - [ ] 15.1.2 Add route guards (authentication, authorization)
  - [ ] 15.1.3 Add navigation menu items

- [ ] 15.2 Update dashboard
  - [ ] 15.2.1 Add BatchStatusWidget to dashboard
  - [ ] 15.2.2 Add quick access links

### 16. Styling & UX

- [ ] 16.1 Implement color coding system
  - [ ] 16.1.1 Define colors cho batch status
  - [ ] 16.1.2 Apply colors consistently across components

- [ ] 16.2 Implement responsive design
  - [ ] 16.2.1 Test trên mobile devices
  - [ ] 16.2.2 Optimize layouts cho small screens
  - [ ] 16.2.3 Test trên tablets

- [ ] 16.3 Add loading states
  - [ ] 16.3.1 Add spinners cho async operations
  - [ ] 16.3.2 Add skeleton screens cho lists

- [ ] 16.4 Add error states
  - [ ] 16.4.1 Design error messages
  - [ ] 16.4.2 Add retry buttons
  - [ ] 16.4.3 Add fallback UI

### 17. Frontend Testing

- [ ] 17.1 Viết unit tests cho stores
  - [ ] 17.1.1 Test tất cả actions
  - [ ] 17.1.2 Test error handling
  - [ ] 17.1.3 Test state mutations

- [ ] 17.2 Viết component tests
  - [ ] 17.2.1 Test user interactions
  - [ ] 17.2.2 Test form validation
  - [ ] 17.2.3 Test error handling

- [ ] 17.3 Viết E2E tests
  - [ ] 17.3.1 Test batch definition creation flow
  - [ ] 17.3.2 Test batch record creation flow
  - [ ] 17.3.3 Test alert viewing flow
  - [ ] 17.3.4 Test report generation flow

## Phase 3: Integration & Testing

### 18. End-to-End Integration

- [ ] 18.1 Test complete batch lifecycle
  - [ ] 18.1.1 Create batch definition
  - [ ] 18.1.2 Create batch record
  - [ ] 18.1.3 Use batch in order
  - [ ] 18.1.4 Verify inventory deduction
  - [ ] 18.1.5 Verify cost calculation

- [ ] 18.2 Test alert system
  - [ ] 18.2.1 Create batch với low stock
  - [ ] 18.2.2 Verify low stock alert appears
  - [ ] 18.2.3 Create batch với short shelf life
  - [ ] 18.2.4 Wait for expiry warning
  - [ ] 18.2.5 Verify expiring alert appears
  - [ ] 18.2.6 Wait for expiry
  - [ ] 18.2.7 Verify expired alert appears

- [ ] 18.3 Test report generation
  - [ ] 18.3.1 Create multiple batches
  - [ ] 18.3.2 Use batches in orders
  - [ ] 18.3.3 Generate production report
  - [ ] 18.3.4 Generate wastage report
  - [ ] 18.3.5 Generate usage report
  - [ ] 18.3.6 Verify data accuracy

### 19. Performance Testing

- [ ] 19.1 Backend performance
  - [ ] 19.1.1 Test với 1000+ batch records
  - [ ] 19.1.2 Test concurrent batch creation (100 requests)
  - [ ] 19.1.3 Test concurrent batch usage
  - [ ] 19.1.4 Measure API response times
  - [ ] 19.1.5 Optimize nếu cần

- [ ] 19.2 Frontend performance
  - [ ] 19.2.1 Test rendering với large lists
  - [ ] 19.2.2 Test pagination performance
  - [ ] 19.2.3 Test chart rendering
  - [ ] 19.2.4 Optimize nếu cần

### 20. Security Testing

- [ ] 20.1 Test authentication
  - [ ] 20.1.1 Verify tất cả endpoints require auth
  - [ ] 20.1.2 Test với invalid tokens
  - [ ] 20.1.3 Test với expired tokens

- [ ] 20.2 Test authorization
  - [ ] 20.2.1 Test manager-only endpoints với barista role
  - [ ] 20.2.2 Test barista endpoints với waiter role
  - [ ] 20.2.3 Verify proper error messages

- [ ] 20.3 Test input validation
  - [ ] 20.3.1 Test với invalid data types
  - [ ] 20.3.2 Test với negative quantities
  - [ ] 20.3.3 Test với SQL injection attempts
  - [ ] 20.3.4 Test với XSS attempts

### 21. User Acceptance Testing

- [ ] 21.1 Manager testing
  - [ ] 21.1.1 Test batch definition management
  - [ ] 21.1.2 Test report viewing
  - [ ] 21.1.3 Test batch record management
  - [ ] 21.1.4 Gather feedback

- [ ] 21.2 Barista testing
  - [ ] 21.2.1 Test batch record creation
  - [ ] 21.2.2 Test alert viewing
  - [ ] 21.2.3 Test mobile usability
  - [ ] 21.2.4 Gather feedback

- [ ] 21.3 Bug fixes và improvements
  - [ ] 21.3.1 Fix reported bugs
  - [ ] 21.3.2 Implement requested improvements
  - [ ] 21.3.3 Re-test

## Phase 4: Deployment & Monitoring

### 22. Deployment Preparation

- [ ] 22.1 Database migration
  - [ ] 22.1.1 Tạo migration scripts
  - [ ] 22.1.2 Test migration trên staging
  - [ ] 22.1.3 Prepare rollback scripts

- [ ] 22.2 Environment configuration
  - [ ] 22.2.1 Configure staging environment
  - [ ] 22.2.2 Configure production environment
  - [ ] 22.2.3 Set up environment variables

- [ ] 22.3 Documentation
  - [ ] 22.3.1 Write user guide cho managers
  - [ ] 22.3.2 Write user guide cho baristas
  - [ ] 22.3.3 Write admin guide
  - [ ] 22.3.4 Write API documentation

### 23. Deployment

- [ ] 23.1 Deploy to staging
  - [ ] 23.1.1 Deploy backend
  - [ ] 23.1.2 Deploy frontend
  - [ ] 23.1.3 Run smoke tests
  - [ ] 23.1.4 Fix issues

- [ ] 23.2 Deploy to production
  - [ ] 23.2.1 Schedule maintenance window
  - [ ] 23.2.2 Backup database
  - [ ] 23.2.3 Deploy backend
  - [ ] 23.2.4 Deploy frontend
  - [ ] 23.2.5 Run smoke tests
  - [ ] 23.2.6 Monitor for issues

### 24. Monitoring & Maintenance

- [ ] 24.1 Set up monitoring
  - [ ] 24.1.1 Monitor API response times
  - [ ] 24.1.2 Monitor error rates
  - [ ] 24.1.3 Monitor database performance
  - [ ] 24.1.4 Set up alerts

- [ ] 24.2 Set up logging
  - [ ] 24.2.1 Configure structured logging
  - [ ] 24.2.2 Set up log aggregation
  - [ ] 24.2.3 Create dashboards

- [ ] 24.3 Post-deployment monitoring
  - [ ] 24.3.1 Monitor for 24 hours
  - [ ] 24.3.2 Check metrics
  - [ ] 24.3.3 Gather user feedback
  - [ ] 24.3.4 Address issues

## Notes

### Testing Framework
- Backend: Go testing package + testify
- Property-Based Testing: gopter
- Frontend: Vitest + Vue Test Utils
- E2E: Playwright

### Priority Levels
- High: Core functionality (batch creation, usage, FIFO)
- Medium: Alerts, reports
- Low: Advanced features, optimizations

### Dependencies
- Tasks 1-7 phải hoàn thành trước khi bắt đầu Phase 2
- Tasks 8-17 có thể parallel với nhau
- Tasks 18-21 phải đợi Phase 1 và 2 hoàn thành

### Estimated Timeline
- Phase 1: 2 weeks
- Phase 2: 2 weeks
- Phase 3: 1 week
- Phase 4: 1 week
- Total: 6 weeks
