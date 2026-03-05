# Kế hoạch Triển khai: Thiết kế lại Template In Bill

## Tổng quan

Triển khai template in hóa đơn mới với logo ở góc trên bên trái, bảng món có cấu trúc rõ ràng, và font size đồng đều. Feature này mở rộng hệ thống in hiện có với các module mới: LogoRenderer, TableFormatter, FontSizeManager, và ImageCompositor.

## Tasks

- [x] 1. Triển khai Logo Renderer
  - [x] 1.1 Tạo `backend/infrastructure/printing/logo_renderer.go`
    - Implement LogoRenderer struct với maxWidthPercent và margin
    - Implement RenderLogo method: load logo từ path, resize, convert to grayscale
    - Implement resizeLogo method với bilinear interpolation
    - Implement convertToGrayscale method
    - Handle errors gracefully (missing file, corrupt file, resize failure)
    - _Requirements: 1.1, 1.2, 1.3, 1.4, 1.5, 1.6, 1.7, 1.8, 7.1, 7.2, 7.5_

  - [ ]* 1.2 Viết property test cho logo positioning
    - **Property 1: Logo hiển thị ở góc trên bên trái**
    - **Validates: Requirements 1.1, 1.5**

  - [ ]* 1.3 Viết property test cho logo loading
    - **Property 2: Logo được load từ đường dẫn đã cấu hình**
    - **Validates: Requirements 1.2**

  - [ ]* 1.4 Viết property test cho missing logo handling
    - **Property 3: Xử lý graceful khi không có logo**
    - **Validates: Requirements 1.3**

  - [ ]* 1.5 Viết property test cho logo sizing
    - **Property 4: Logo sizing constraint**
    - **Validates: Requirements 1.4, 1.7**

  - [ ]* 1.6 Viết property test cho logo format support
    - **Property 5: Logo format support**
    - **Validates: Requirements 1.6**

  - [ ]* 1.7 Viết property test cho grayscale conversion
    - **Property 6: Logo grayscale conversion**
    - **Validates: Requirements 1.8**

  - [ ]* 1.8 Viết unit tests cho logo renderer
    - Test load PNG, JPG, JPEG files
    - Test resize với different sizes
    - Test error handling (missing file, corrupt file)
    - _Requirements: 1.1, 1.2, 1.3, 1.4, 1.6, 1.8_

- [x] 2. Triển khai Table Formatter
  - [x] 2.1 Tạo `backend/infrastructure/printing/table_formatter.go`
    - Implement TableFormatter struct với paperWidth, margin, columnGap
    - Define TableColumn và TableRow structs
    - Implement FormatItemsTable method: format order items thành table lines
    - Implement calculateColumnWidths: tính column widths (50%, 15%, 17.5%, 17.5%)
    - Implement formatRow: format một row với column alignment
    - Implement wrapCellText: wrap text trong cell nếu quá dài
    - _Requirements: 2.1, 2.2, 2.3, 2.4, 2.5, 2.6, 2.7, 2.8_

  - [ ]* 2.2 Viết property test cho table structure
    - **Property 7: Table structure completeness**
    - **Validates: Requirements 2.1**

  - [ ]* 2.3 Viết property test cho column alignment
    - **Property 8: Table column alignment consistency**
    - **Validates: Requirements 2.2, 2.3, 2.8**

  - [ ]* 2.4 Viết property test cho separator lines
    - **Property 9: Table separator lines**
    - **Validates: Requirements 2.4**

  - [ ]* 2.5 Viết property test cho text wrapping
    - **Property 10: Long item name wrapping**
    - **Validates: Requirements 2.5**

  - [ ]* 2.6 Viết property test cho variant display
    - **Property 11: Variant display on sub-line**
    - **Validates: Requirements 2.6**

  - [ ]* 2.7 Viết property test cho column width calculation
    - **Property 12: Column width calculation**
    - **Validates: Requirements 2.7**

  - [ ]* 2.8 Viết unit tests cho table formatter
    - Test table với different số lượng items
    - Test với long item names
    - Test với items có variants
    - Test column width calculation cho 58mm và 80mm
    - _Requirements: 2.1, 2.2, 2.3, 2.5, 2.6, 2.7_

- [x] 3. Triển khai Font Size Manager
  - [x] 3.1 Tạo `backend/infrastructure/printing/font_size_manager.go`
    - Implement FontSizeManager struct với normalSize (18pt), headerSize (22pt), totalSize (20pt)
    - Define FontSizeConfig struct với Size và Bold
    - Implement GetFontSizeForLine method: detect font size dựa trên line content
    - Implement detection logic cho headers, totals, regular content
    - _Requirements: 3.1, 3.2, 3.3, 3.4, 3.5, 3.6, 3.7_

  - [ ]* 3.2 Viết property test cho regular content font size
    - **Property 13: Font size consistency for regular content**
    - **Validates: Requirements 3.1, 3.4, 3.5**

  - [ ]* 3.3 Viết property test cho header font size
    - **Property 14: Header font size**
    - **Validates: Requirements 3.2**

  - [ ]* 3.4 Viết property test cho total font size
    - **Property 15: Total line font size**
    - **Validates: Requirements 3.3**

  - [ ]* 3.5 Viết property test cho font weight
    - **Property 16: Font weight for headers and totals**
    - **Property 17: Font weight for regular content**
    - **Validates: Requirements 3.6, 3.7**

  - [ ]* 3.6 Viết unit tests cho font size manager
    - Test detection cho different line types
    - Test font size assignment
    - Test font weight assignment
    - _Requirements: 3.1, 3.2, 3.3, 3.6, 3.7_

- [x] 4. Triển khai Image Compositor
  - [x] 4.1 Tạo `backend/infrastructure/printing/image_compositor.go`
    - Implement ImageCompositor struct với paperWidth và margin
    - Implement Compose method: kết hợp logo và text content
    - Implement calculateTotalHeight: tính total height cần thiết
    - Implement drawLogo: vẽ logo lên combined image
    - Implement drawTextContent: vẽ text content lên combined image
    - Handle case khi logo là nil (không có logo)
    - _Requirements: 1.1, 1.3, 1.5_

  - [ ]* 4.2 Viết unit tests cho image compositor
    - Test compose với logo
    - Test compose không có logo
    - Test height calculation
    - Test image positioning
    - _Requirements: 1.1, 1.3, 1.5_

- [x] 5. Checkpoint - Đảm bảo core modules hoạt động
  - Chạy tất cả unit tests
  - Test thủ công từng module riêng lẻ
  - Hỏi user nếu có vấn đề

- [x] 6. Mở rộng Format Parser
  - [x] 6.1 Cập nhật `backend/infrastructure/printing/format_parser.go`
    - Thêm FontSize field vào LineFormat struct
    - Thêm IsTableRow field vào LineFormat struct
    - Implement detectFontSize method
    - Implement isTableRow method
    - Update Parse method để set FontSize và IsTableRow
    - _Requirements: 3.1, 3.2, 3.3_

  - [ ]* 6.2 Viết unit tests cho enhanced format parser
    - Test font size detection
    - Test table row detection
    - Test backward compatibility với existing logic
    - _Requirements: 3.1, 3.2, 3.3_

- [x] 7. Mở rộng Text Renderer
  - [x] 7.1 Cập nhật `backend/infrastructure/printing/text_renderer.go`
    - Thêm FontSizes map vào RendererConfig
    - Thêm fonts map[float64]FontPair vào TextRenderer
    - Update NewTextRenderer để load multiple font sizes (18pt, 20pt, 22pt)
    - Update drawLine để select font dựa trên line.FontSize
    - Update calculateImageHeight để account for different font sizes
    - _Requirements: 3.1, 3.2, 3.3_

  - [ ]* 7.2 Viết unit tests cho enhanced text renderer
    - Test rendering với different font sizes
    - Test height calculation với mixed font sizes
    - Test backward compatibility
    - _Requirements: 3.1, 3.2, 3.3_

- [x] 8. Tích hợp các modules vào Template Renderer
  - [x] 8.1 Cập nhật `backend/application/services/template_renderer.go`
    - Inject LogoRenderer, TableFormatter, FontSizeManager, ImageCompositor
    - Update RenderBill method:
      - Check if logo is configured
      - If yes, use LogoRenderer to render logo
      - Use TableFormatter to format items table
      - Use FontSizeManager to assign font sizes
      - Use TextRenderer to render text content
      - Use ImageCompositor to combine logo + text
    - Handle [LOGO], [TABLE_START], [TABLE_END] markers
    - Implement error handling và fallback logic
    - _Requirements: 1.1, 1.2, 1.3, 2.1, 3.1, 4.5, 7.1, 7.2, 7.3_

  - [ ]* 8.2 Viết property test cho template variable support
    - **Property 19: Template variable support**
    - **Validates: Requirements 4.5**

  - [ ]* 8.3 Viết property test cho paper width compatibility
    - **Property 20: Paper width compatibility**
    - **Validates: Requirements 4.6**

  - [ ]* 8.4 Viết property test cho error handling
    - **Property 27: Missing logo graceful handling**
    - **Property 28: Corrupt logo graceful handling**
    - **Property 29: Template rendering fallback**
    - **Property 31: Extreme logo size handling**
    - **Validates: Requirements 7.1, 7.2, 7.3, 7.5**

  - [ ]* 8.5 Viết integration tests cho template rendering
    - Test end-to-end rendering với logo + table + font sizes
    - Test với 58mm và 80mm paper
    - Test error scenarios
    - _Requirements: 1.1, 2.1, 3.1, 4.5, 4.6_

- [x] 9. Checkpoint - Đảm bảo backend integration hoạt động
  - Chạy tất cả tests (unit + property + integration)
  - Test thủ công: render bill với logo và table
  - Test error scenarios
  - Hỏi user nếu có vấn đề

- [x] 10. Tạo template mới trong database
  - [x] 10.1 Tạo `backend/cmd/create-new-bill-template/main.go`
    - Script để tạo template mới với logo và table format
    - Template content theo design document
    - Set is_default = false (không override template cũ)
    - _Requirements: 4.1, 6.1, 6.4_

  - [ ]* 10.2 Viết unit test cho template creation
    - Test template được tạo với correct type (BILL)
    - Test template không override existing default
    - _Requirements: 4.1, 6.1_

- [x] 11. Triển khai logo upload backend
  - [x] 11.1 Tạo `backend/interfaces/http/logo_upload_handler.go`
    - POST /api/settings/logo - Upload logo
    - DELETE /api/settings/logo - Delete logo
    - Validate file format (PNG, JPG, JPEG)
    - Validate file size (<= 2MB)
    - Save file to uploads/logos/ directory
    - Update shop_settings.logo_url
    - Return logo URL
    - _Requirements: 5.1, 5.2, 5.4, 5.5, 5.6, 5.7_

  - [ ]* 11.2 Viết property test cho logo upload validation
    - **Property 22: Logo format validation**
    - **Property 23: Logo file size validation**
    - **Validates: Requirements 5.5, 5.6**

  - [ ]* 11.3 Viết property test cho logo persistence
    - **Property 21: Logo path persistence**
    - **Property 24: Logo storage location**
    - **Validates: Requirements 5.2, 5.7**

  - [ ]* 11.4 Viết unit tests cho logo upload handler
    - Test upload với valid files
    - Test validation errors
    - Test file storage
    - Test database update
    - _Requirements: 5.2, 5.5, 5.6, 5.7_

- [x] 12. Triển khai template management backend
  - [x] 12.1 Cập nhật `backend/interfaces/http/print_template_handler.go`
    - Thêm logic để list tất cả templates (không chỉ default)
    - Thêm endpoint để set template as default
    - Ensure old templates không bị xóa
    - _Requirements: 4.2, 6.1, 6.2, 6.5_

  - [ ]* 12.2 Viết property test cho template management
    - **Property 18: Default template usage**
    - **Property 25: Template preservation**
    - **Property 26: Template switching**
    - **Validates: Requirements 4.2, 4.7, 6.1, 6.2, 6.3, 6.4**

  - [ ]* 12.3 Viết unit tests cho template management
    - Test set default template
    - Test template switching
    - Test template preservation
    - _Requirements: 4.2, 6.1, 6.2_

- [x] 13. Triển khai error logging
  - [x] 13.1 Cập nhật error logging trong tất cả modules
    - Logo renderer: log missing file, corrupt file, resize errors
    - Template renderer: log rendering failures
    - Include timestamp, error message, context trong logs
    - _Requirements: 7.1, 7.2, 7.3, 7.4, 7.5_

  - [ ]* 13.2 Viết property test cho error logging
    - **Property 30: Error logging completeness**
    - **Validates: Requirements 7.4**

  - [ ]* 13.3 Viết unit tests cho error logging
    - Test log format
    - Test log content completeness
    - _Requirements: 7.4_

- [x] 14. Checkpoint - Đảm bảo backend hoàn chỉnh
  - Chạy tất cả tests (unit + property + integration)
  - Test thủ công tất cả API endpoints
  - Test error scenarios và logging
  - Hỏi user nếu có vấn đề

- [x] 15. Triển khai logo upload frontend
  - [x] 15.1 Cập nhật `frontend/src/views/SettingsView.vue` (hoặc tương tự)
    - Thêm logo upload field với file input
    - Implement upload handler: call POST /api/settings/logo
    - Hiển thị preview sau khi upload
    - Thêm delete button: call DELETE /api/settings/logo
    - Show validation errors (format, size)
    - _Requirements: 5.1, 5.3, 5.4, 5.5, 5.6_

  - [ ]* 15.2 Viết component tests cho logo upload
    - Test upload flow
    - Test preview display
    - Test delete flow
    - Test validation error display
    - _Requirements: 5.1, 5.3, 5.4_

- [x] 16. Triển khai template management frontend
  - [x] 16.1 Cập nhật `frontend/src/views/PrintManagementView.vue`
    - Hiển thị danh sách tất cả templates (không chỉ default)
    - Thêm button "Set as Default" cho mỗi template
    - Highlight template hiện tại đang là default
    - Thêm preview cho mỗi template
    - _Requirements: 4.2, 6.2, 6.5_

  - [ ]* 16.2 Viết component tests cho template management
    - Test template list display
    - Test set default action
    - Test template switching
    - _Requirements: 4.2, 6.2, 6.5_

- [ ] 17. Testing end-to-end
  - [x] 17.1 Manual testing checklist
    - Upload logo PNG - verify preview và storage
    - Upload logo JPG - verify preview và storage
    - Upload logo > 2MB - verify rejection
    - Upload invalid format - verify rejection
    - Create order với nhiều items - verify table format
    - Create order với long item names - verify wrapping
    - Create order với variants - verify sub-line display
    - Test với 58mm paper - verify layout
    - Test với 80mm paper - verify layout
    - Delete logo - verify bills render without logo
    - Corrupt logo file - verify graceful fallback
    - Switch template mới/cũ - verify correct rendering
    - _Requirements: All_

  - [ ] 17.2 Chạy tất cả property tests (100 iterations)
    - Verify tất cả 31 properties pass
    - Fix any failing tests
    - _Requirements: All_

  - [ ] 17.3 Chạy tất cả unit tests
    - Verify coverage >= 80%
    - Fix any failing tests
    - _Requirements: All_

- [x] 18. Documentation
  - [x] 18.1 Viết user guide
    - Hướng dẫn upload logo
    - Hướng dẫn switch templates
    - Hướng dẫn xử lý lỗi thường gặp
    - _Requirements: All_

  - [x] 18.2 Update API documentation
    - Document logo upload endpoints
    - Document template management endpoints
    - Include request/response examples

- [x] 19. Final checkpoint - Hoàn thành feature
  - Verify tất cả requirements được implement
  - Verify tất cả tests pass
  - Verify documentation đầy đủ
  - User acceptance testing
  - Hỏi user nếu cần thêm gì

## Ghi chú

- Tasks đánh dấu `*` là optional và có thể skip để có MVP nhanh hơn
- Mỗi task tham chiếu đến requirements cụ thể để dễ traceability
- Checkpoints đảm bảo validation từng bước
- Property tests validate universal correctness properties (minimum 100 iterations)
- Unit tests validate specific examples và edge cases
- Integration tests validate end-to-end flows
- Sử dụng thư viện `gopter` cho property-based testing trong Go
