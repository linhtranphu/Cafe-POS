# Kế Hoạch Triển Khai: In Bill và Tem Đơn Hàng

## Tổng Quan

Triển khai hệ thống in tự động bill và tem cho đơn hàng, sử dụng kiến trúc event-driven với queue-based printing. Hệ thống tích hợp với order service hiện có và hỗ trợ máy in nhiệt ESC/POS.

## Tasks

- [x] 1. Thiết lập cấu trúc domain layer và entities
  - Tạo package `backend/domain/printing/`
  - Định nghĩa entities: `PrintJob`, `PrinterConfig`, `PrintTemplate`
  - Định nghĩa constants cho types và statuses
  - Định nghĩa repository interfaces
  - _Requirements: 1.1, 2.1, 3.1, 3.2, 4.1, 5.1_

- [x] 1.1 Viết unit tests cho domain entities
  - Test validation logic cho PrintJob
  - Test status transitions
  - Test PrinterConfig validation
  - _Requirements: 1.1, 2.1, 3.2, 4.1_

- [ ] 2. Triển khai MongoDB repositories
  - [x] 2.1 Tạo `backend/infrastructure/mongodb/print_job_repository.go`
    - Implement PrintJobRepository interface
    - Implement các methods: Create, FindByID, FindByOrderID, FindPending, FindFailed, UpdateStatus, IncrementRetry, Delete, DeleteOldCompleted
    - _Requirements: 4.1, 4.2, 4.7_

  - [x] 2.2 Tạo `backend/infrastructure/mongodb/printer_config_repository.go`
    - Implement PrinterConfigRepository interface
    - Implement các methods: Create, FindByID, FindAll, FindByType, FindDefault, Update, Delete
    - _Requirements: 3.1, 3.2, 3.3, 3.4_

  - [x] 2.3 Tạo `backend/infrastructure/mongodb/print_template_repository.go`
    - Implement PrintTemplateRepository interface
    - Implement các methods: Create, FindByID, FindByType, FindDefault, Update, Delete
    - _Requirements: 5.1, 5.2, 5.3_

  - [x] 2.4 Tạo migration script cho collections và indexes
    - Tạo `backend/cmd/migrate/create_printing_collections.go`
    - Tạo collections: print_jobs, printer_configs, print_templates
    - Tạo indexes theo design document
    - Tạo TTL index cho print_jobs (7 days)
    - _Requirements: 4.7_

- [x] 2.5 Viết integration tests cho repositories
  - Test CRUD operations với MongoDB
  - Test query methods (FindPending, FindFailed, FindDefault)
  - Test TTL index behavior
  - _Requirements: 4.1, 4.7_

- [ ] 3. Triển khai template rendering service
  - [x] 3.1 Tạo `backend/application/services/template_renderer.go`
    - Implement TemplateRenderer interface
    - Implement RenderBill method với Go template
    - Implement RenderLabel method
    - Xử lý template errors và fallback to default
    - _Requirements: 1.2, 1.3, 1.4, 1.5, 1.6, 2.3, 2.4, 2.5, 2.6, 2.7_

  - [x] 3.2 Tạo default templates
    - Tạo default bill template với ESC/POS formatting
    - Tạo default label template
    - Đảm bảo templates tuân thủ width constraints
    - _Requirements: 1.7, 2.8, 5.4, 5.5_

- [ ]* 3.3 Viết property test cho template rendering
  - **Property 2: Bill content completeness**
  - **Validates: Requirements 1.2, 1.3, 1.4, 1.5, 1.6**

- [ ]* 3.4 Viết property test cho label rendering
  - **Property 3: Label content completeness**
  - **Validates: Requirements 2.3, 2.4, 2.5, 2.6, 2.7**

- [ ]* 3.5 Viết property test cho format constraints
  - **Property 4: Bill format width constraint**
  - **Property 5: Label format size constraint**
  - **Validates: Requirements 1.7, 2.8**

- [ ]* 3.6 Viết unit tests cho template rendering
  - Test rendering với specific order data
  - Test error handling cho invalid templates
  - Test fallback to default template
  - _Requirements: 1.2, 1.3, 1.4, 1.5, 1.6, 2.3, 2.4, 2.5, 2.6, 2.7_

- [ ] 4. Triển khai printer manager và drivers
  - [x] 4.1 Tạo `backend/application/services/printer_manager.go`
    - Implement PrinterManager interface
    - Implement GetPrinter method (factory pattern)
    - Implement TestConnection method
    - _Requirements: 3.5, 3.6_

  - [x] 4.2 Tạo `backend/infrastructure/printing/escpos_printer.go`
    - Implement Printer interface cho ESC/POS thermal printers
    - Implement Connect, Disconnect, Print, GetStatus methods
    - Hỗ trợ network connection (TCP/IP)
    - Convert template output sang ESC/POS commands
    - _Requirements: 1.7, 1.8, 1.9, 3.5_

  - [x] 4.3 Tạo `backend/infrastructure/printing/label_printer.go`
    - Implement Printer interface cho label printers
    - Hỗ trợ các kích thước label khác nhau
    - _Requirements: 2.8, 2.9, 2.10_

- [x] 4.4 Viết unit tests cho printer drivers
  - Test ESC/POS command generation
  - Test connection handling
  - Test error detection
  - Mock network connections
  - _Requirements: 1.8, 1.9, 2.9, 2.10, 3.5_

- [ ]* 4.5 Viết property test cho printer status check
  - **Property 12: Printer status check**
  - **Validates: Requirements 3.5, 7.5**

- [ ] 5. Triển khai print service
  - [x] 5.1 Tạo `backend/application/services/print_service.go`
    - Implement PrintService interface
    - Implement CreatePrintJobsForOrder method
    - Implement ReprintBill và ReprintLabel methods
    - Implement GetPendingJobs, GetFailedJobs methods
    - Implement RetryJob và CancelJob methods
    - _Requirements: 1.1, 1a, 1b, 2.1, 2.2, 4.4, 4.5, 4.6, 6.1_

  - [x] 5.2 Implement auto-create print jobs logic
    - Tạo 1 bill job và N label jobs cho mỗi order
    - Validate order status (PAID)
    - Get default printers cho bill và label
    - Render templates
    - Save jobs với status PENDING
    - _Requirements: 1.1, 1a, 2.1, 2.2, 6.1, 6.2_

- [ ]* 5.3 Viết property test cho auto-create print jobs
  - **Property 1: Tự động tạo print jobs khi order được tạo**
  - **Validates: Requirements 1.1, 1a, 2.1, 2.2, 6.1**

- [ ]* 5.4 Viết property test cho reprint capability
  - **Property 8: Reprint capability**
  - **Validates: Requirements 1b**

- [ ]* 5.5 Viết unit tests cho print service
  - Test CreatePrintJobsForOrder với specific orders
  - Test ReprintBill và ReprintLabel
  - Test error handling khi printer không tồn tại
  - _Requirements: 1.1, 1b, 2.1, 2.2, 6.1_

- [ ] 6. Triển khai print worker và queue processing
  - [x] 6.1 Tạo `backend/application/services/print_worker.go`
    - Implement background worker để process print queue
    - Implement ProcessJob method
    - Implement retry logic với exponential backoff
    - Implement status updates (PENDING → PRINTING → COMPLETED/FAILED)
    - _Requirements: 1.8, 1.9, 2.9, 2.10, 4.2, 4.3_

  - [x] 6.2 Implement retry mechanism
    - Auto retry tối đa 3 lần
    - Exponential backoff (30s, 1m, 2m)
    - Update retry_count
    - Set status FAILED sau 3 lần
    - _Requirements: 4.2, 4.3_

  - [x] 6.3 Implement queue polling
    - Poll pending jobs mỗi 10 giây
    - Process jobs theo thứ tự created_at
    - Handle concurrent processing
    - _Requirements: 1.8, 1.9, 2.9, 2.10_

- [ ]* 6.4 Viết property test cho print command sent immediately
  - **Property 6: Print command sent immediately**
  - **Validates: Requirements 1.8, 2.9**

- [ ]* 6.5 Viết property test cho queue and retry
  - **Property 7: Queue and retry on printer unavailable**
  - **Validates: Requirements 1.9, 2.10**

- [ ]* 6.6 Viết property test cho automatic retry with limit
  - **Property 15: Automatic retry with limit**
  - **Validates: Requirements 4.2, 4.3**

- [ ]* 6.7 Viết unit tests cho print worker
  - Test job processing với mock printer
  - Test retry logic với different error types
  - Test status transitions
  - _Requirements: 1.8, 1.9, 2.9, 2.10, 4.2, 4.3_

- [ ] 7. Tích hợp với order service
  - [x] 7.1 Thêm event emission vào order service
    - Emit OrderCreated event khi order chuyển sang PAID
    - Đảm bảo order được commit trước khi emit event
    - _Requirements: 6.1, 6.2_

  - [x] 7.2 Implement event handler trong print service
    - Listen to OrderCreated event
    - Call CreatePrintJobsForOrder
    - Handle errors mà không ảnh hưởng order creation
    - Log print job creation
    - _Requirements: 6.1, 6.3, 6.4_

  - [x] 7.3 Thêm auto-print setting
    - Tạo shop_settings field cho auto_print_enabled
    - Check setting trước khi tạo print jobs
    - _Requirements: 6.5, 6.6_

- [ ]* 7.4 Viết property test cho order persistence before print jobs
  - **Property 20: Order persistence before print jobs**
  - **Validates: Requirements 6.2**

- [ ]* 7.5 Viết property test cho print failure isolation
  - **Property 21: Print failure isolation**
  - **Validates: Requirements 6.3**

- [ ]* 7.6 Viết property test cho auto-print toggle
  - **Property 23: Auto-print toggle**
  - **Validates: Requirements 6.5, 6.6**

- [ ]* 7.7 Viết integration tests cho order-print flow
  - Test end-to-end: create order → print jobs created
  - Test với auto-print enabled/disabled
  - Test error handling
  - _Requirements: 6.1, 6.2, 6.3, 6.5, 6.6_

- [x] 8. Checkpoint - Đảm bảo backend core functionality hoạt động
  - Chạy tất cả tests
  - Test thủ công: tạo order → verify print jobs created
  - Test retry mechanism
  - Hỏi user nếu có vấn đề

- [ ] 9. Triển khai HTTP handlers
  - [x] 9.1 Tạo `backend/interfaces/http/print_job_handler.go`
    - GET /api/print-jobs - List print jobs
    - GET /api/print-jobs/:id - Get job detail
    - GET /api/print-jobs/pending - Get pending jobs
    - GET /api/print-jobs/failed - Get failed jobs
    - POST /api/print-jobs/:id/retry - Retry job
    - DELETE /api/print-jobs/:id - Cancel job
    - _Requirements: 4.4, 4.5, 4.6_

  - [x] 9.2 Tạo `backend/interfaces/http/printer_config_handler.go`
    - GET /api/printers - List printers
    - GET /api/printers/:id - Get printer detail
    - POST /api/printers - Create printer config
    - PUT /api/printers/:id - Update printer config
    - DELETE /api/printers/:id - Delete printer config
    - POST /api/printers/:id/test - Test connection
    - _Requirements: 3.1, 3.2, 3.3, 3.4, 3.5, 3.6, 3.7_

  - [x] 9.3 Tạo `backend/interfaces/http/print_template_handler.go`
    - GET /api/print-templates - List templates
    - GET /api/print-templates/:id - Get template detail
    - POST /api/print-templates - Create template
    - PUT /api/print-templates/:id - Update template
    - DELETE /api/print-templates/:id - Delete template
    - POST /api/print-templates/:id/preview - Preview template
    - _Requirements: 5.1, 5.2, 5.3, 5.6, 5.7_

  - [x] 9.4 Thêm reprint endpoints vào order handler
    - POST /api/orders/:id/reprint-bill - Reprint bill
    - POST /api/orders/:id/reprint-label - Reprint label
    - _Requirements: 1b_

  - [x] 9.5 Register routes trong main.go
    - Wire up handlers với services
    - Add authentication middleware
    - _Requirements: All API endpoints_

- [ ]* 9.6 Viết HTTP handler tests
  - Test request validation
  - Test response formatting
  - Test error handling
  - _Requirements: 4.4, 4.5, 4.6, 3.1-3.7, 5.1-5.7, 1b_

- [ ]* 9.7 Viết integration tests cho HTTP APIs
  - Test end-to-end API calls
  - Test authentication
  - Test error responses
  - _Requirements: All API endpoints_

- [ ]* 9.8 Viết property tests cho printer configuration
  - **Property 9: Multiple printer configuration**
  - **Property 10: Default printer uniqueness**
  - **Property 11: Printer connection validation**
  - **Property 13: Printer disable without delete**
  - **Validates: Requirements 3.1, 3.2, 3.3, 3.4, 3.6, 3.7**

- [ ]* 9.9 Viết property tests cho print job management
  - **Property 14: Print job status tracking**
  - **Property 16: Manual retry capability**
  - **Property 17: Cancel pending jobs**
  - **Validates: Requirements 4.1, 4.4, 4.5, 4.6**

- [ ] 10. Triển khai error handling và logging
  - [x] 10.1 Implement error logging
    - Log mỗi print job creation
    - Log mỗi print attempt (success/failure)
    - Log printer connection errors
    - Log template rendering errors
    - Include timestamp, job_id, order_id, error details
    - _Requirements: 6.4, 7.2_

  - [x] 10.2 Implement error notifications
    - Tạo notification system cho print failures
    - Hiển thị printer offline warnings
    - Hiển thị hardware errors (paper out, jam)
    - _Requirements: 7.1, 7.3, 7.5_

  - [x] 10.3 Implement error history
    - Store error details trong print_jobs
    - Provide API để query error history
    - _Requirements: 7.4_

- [ ]* 10.4 Viết property test cho print activity logging
  - **Property 22: Print activity logging**
  - **Validates: Requirements 6.4, 7.2**

- [ ]* 10.5 Viết property test cho error history
  - **Property 24: Error history queryable**
  - **Validates: Requirements 7.4**

- [ ]* 10.6 Viết unit tests cho error handling
  - Test different error types
  - Test error message formatting
  - Test notification triggering
  - _Requirements: 7.1, 7.2, 7.3, 7.4, 7.5_

- [ ] 11. Triển khai background jobs
  - [x] 11.1 Tạo cleanup job cho old print jobs
    - Chạy daily
    - Xóa completed jobs cũ hơn 7 ngày
    - Log cleanup results
    - _Requirements: 4.7_

  - [x] 11.2 Start print worker trong main.go
    - Start worker goroutine khi app starts
    - Graceful shutdown handling
    - _Requirements: 1.8, 1.9, 2.9, 2.10_

- [ ]* 11.3 Viết property test cho automatic cleanup
  - **Property 18: Automatic cleanup of old jobs**
  - **Validates: Requirements 4.7**

- [ ]* 11.4 Viết unit tests cho background jobs
  - Test cleanup job logic
  - Test worker lifecycle
  - _Requirements: 4.7_

- [x] 12. Checkpoint - Đảm bảo backend hoàn chỉnh
  - Chạy tất cả tests (unit + property + integration)
  - Test thủ công tất cả API endpoints
  - Test error scenarios
  - Hỏi user nếu có vấn đề

- [x] 13. Triển khai frontend services
  - [x] 13.1 Tạo `frontend/src/services/printJob.js`
    - API client methods cho print job endpoints
    - fetchPrintJobs, fetchPendingJobs, fetchFailedJobs
    - retryJob, cancelJob
    - _Requirements: 4.4, 4.5, 4.6_

  - [x] 13.2 Tạo `frontend/src/services/printerConfig.js`
    - API client methods cho printer config endpoints
    - fetchPrinters, createPrinter, updatePrinter, deletePrinter
    - testConnection
    - _Requirements: 3.1, 3.2, 3.3, 3.4, 3.5, 3.6, 3.7_

  - [x] 13.3 Tạo `frontend/src/services/printTemplate.js`
    - API client methods cho template endpoints
    - fetchTemplates, createTemplate, updateTemplate, deleteTemplate
    - previewTemplate
    - _Requirements: 5.1, 5.2, 5.3, 5.6, 5.7_

- [x] 14. Triển khai Pinia stores
  - [x] 14.1 Tạo `frontend/src/stores/printJob.js`
    - State: printJobs, pendingJobs, failedJobs
    - Actions: fetchJobs, retryJob, cancelJob
    - Getters: jobsByOrderId, jobsByStatus
    - _Requirements: 4.4, 4.5, 4.6_

  - [x] 14.2 Tạo `frontend/src/stores/printerConfig.js`
    - State: printers, defaultBillPrinter, defaultLabelPrinter
    - Actions: fetchPrinters, savePrinter, deletePrinter, testConnection
    - Getters: printersByType, enabledPrinters
    - _Requirements: 3.1, 3.2, 3.3, 3.4, 3.5, 3.6, 3.7_

  - [x] 14.3 Tạo `frontend/src/stores/printTemplate.js`
    - State: templates, defaultTemplates
    - Actions: fetchTemplates, saveTemplate, deleteTemplate, previewTemplate
    - Getters: templatesByType
    - _Requirements: 5.1, 5.2, 5.3, 5.6, 5.7_

- [x] 15. Triển khai Vue components
  - [x] 15.1 Tạo `frontend/src/components/printing/PrintJobList.vue`
    - Hiển thị danh sách print jobs
    - Filter theo status (pending, failed, completed)
    - Actions: retry, cancel
    - Real-time updates
    - _Requirements: 4.4, 4.5, 4.6_

  - [x] 15.2 Tạo `frontend/src/components/printing/PrinterConfigForm.vue`
    - Form để tạo/edit printer config
    - Fields: name, type, connection_type, connection details, paper_width
    - Test connection button
    - Set default checkbox
    - Enable/disable toggle
    - _Requirements: 3.1, 3.2, 3.3, 3.4, 3.5, 3.6, 3.7_

  - [x] 15.3 Tạo `frontend/src/components/printing/PrinterConfigList.vue`
    - Hiển thị danh sách printers
    - Show status (online/offline)
    - Actions: edit, delete, test connection
    - Highlight default printers
    - _Requirements: 3.1, 3.2, 3.3, 3.4, 3.5_

  - [x] 15.4 Tạo `frontend/src/components/printing/PrintTemplateEditor.vue`
    - Code editor cho template content
    - Preview pane với sample data
    - Save/load templates
    - Set default checkbox
    - _Requirements: 5.1, 5.2, 5.3, 5.6, 5.7_

  - [x] 15.5 Thêm reprint buttons vào OrderDetailView
    - Button "In lại Bill"
    - Button "In lại Tem" cho mỗi item
    - Show loading state
    - Show success/error notifications
    - _Requirements: 1b_

  - [x] 15.6 Tạo `frontend/src/components/printing/PrintErrorNotification.vue`
    - Component để hiển thị print errors
    - Show error message
    - Action buttons (retry, dismiss)
    - Auto-dismiss after timeout
    - _Requirements: 7.1, 7.3_

- [ ]* 15.7 Viết component tests
  - Test PrintJobList rendering và actions
  - Test PrinterConfigForm validation và submission
  - Test PrintTemplateEditor editing và preview
  - _Requirements: 4.4, 4.5, 4.6, 3.1-3.7, 5.1-5.7_

- [x] 16. Triển khai views
  - [x] 16.1 Tạo `frontend/src/views/PrintManagementView.vue`
    - Tabs: Print Jobs, Printers, Templates
    - Integrate PrintJobList, PrinterConfigList, PrintTemplateEditor
    - _Requirements: 4.4, 3.1, 5.1_

  - [x] 16.2 Thêm route cho print management
    - Add route /print-management
    - Add menu item trong navigation
    - Restrict to admin users
    - _Requirements: 4.4, 3.1, 5.1_

- [x] 17. Triển khai real-time updates
  - [x] 17.1 Setup WebSocket connection cho print job updates
    - Listen to print-job-created events
    - Listen to print-job-status-changed events
    - Listen to print-job-failed events
    - Update store state
    - _Requirements: 7.1_

  - [x] 17.2 Implement notification system
    - Show toast notification khi print job fails
    - Show printer offline warnings
    - Action buttons trong notifications
    - _Requirements: 7.1, 7.3_

- [x] 18. Triển khai template customization settings
  - [x] 18.1 Tạo shop settings form cho print configuration
    - Shop info fields (name, address, phone)
    - Logo upload
    - Custom message textarea
    - Paper width selector (58mm/80mm)
    - Label size selector
    - Field visibility togglesb
    - Auto-print enable/disable
    - _Requirements: 5.1, 5.2, 5.3, 5.4, 5.5, 5.6, 6.5_

  - [x] 18.2 Integrate settings vào template rendering
    - Pass shop settings to template renderer
    - Apply settings khi render bill/label
    - _Requirements: 5.1, 5.2, 5.3, 5.7_

- [ ]* 18.3 Viết property test cho template configuration
  - **Property 19: Template configuration affects rendering**
  - **Validates: Requirements 5.1, 5.2, 5.3, 5.4, 5.5, 5.6, 5.7**

- [~] 19. Checkpoint - Đảm bảo frontend hoàn chỉnh
  - Test tất cả components
  - Test integration với backend APIs
  - Test real-time updates
  - Test error handling
  - Hỏi user nếu có vấn đề

- [ ] 20. Testing và bug fixes
  - [~] 20.1 Chạy tất cả property tests (100 iterations)
    - Verify tất cả 24 properties pass
    - Fix any failing tests
    - _Requirements: All_

  - [~] 20.2 Chạy tất cả unit tests
    - Verify coverage >80%
    - Fix any failing tests
    - _Requirements: All_

  - [~] 20.3 Chạy integration tests
    - Test end-to-end flows
    - Test error scenarios
    - _Requirements: All_

  - [~] 20.4 Manual testing
    - Test với real printer (nếu có)
    - Test tất cả user flows
    - Test error handling
    - _Requirements: All_

- [ ] 21. Documentation và deployment
  - [~] 21.1 Viết API documentation
    - Document tất cả endpoints
    - Include request/response examples
    - Document error codes

  - [~] 21.2 Viết user guide
    - Hướng dẫn cấu hình máy in
    - Hướng dẫn tùy chỉnh templates
    - Hướng dẫn xử lý lỗi thường gặp
    - _Requirements: 7.6_

  - [~] 21.3 Setup monitoring và alerts
    - Monitor print job failure rate
    - Monitor printer status
    - Alert khi failure rate cao

  - [~] 21.4 Deploy to production
    - Run migrations
    - Deploy backend
    - Deploy frontend
    - Verify functionality

- [~] 22. Final checkpoint - Hoàn thành feature
  - Verify tất cả requirements được implement
  - Verify tất cả tests pass
  - Verify documentation đầy đủ
  - User acceptance testing
  - Hỏi user nếu cần thêm gì

## Ghi Chú

- Tasks đánh dấu `*` là optional và có thể skip để có MVP nhanh hơn
- Mỗi task tham chiếu đến requirements cụ thể để dễ traceability
- Checkpoints đảm bảo validation từng bước
- Property tests validate universal correctness properties
- Unit tests validate specific examples và edge cases
- Integration tests validate end-to-end flows
