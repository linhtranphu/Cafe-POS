# HTML Template Management Implementation

## Tổng quan

Đã implement tính năng quản lý HTML templates trong tab Templates của Print Management. Tính năng này cho phép:

1. Load/Save HTML template từ/đến backend
2. Live preview template với dữ liệu mẫu
3. Test print với order thật
4. Preview PNG với order thật

## Kiến trúc

### Frontend

**Component:** `frontend/src/components/printing/HTMLTemplateEditor.vue`

**Tính năng:**
- HTML editor với syntax highlighting
- Live preview trong iframe
- Sample data panel (show/hide)
- Order search và selection
- Test print với order thật
- Preview PNG với order thật
- Auto-save với debounce

**Layout:**
```
┌─────────────────────────────────────────────────────┐
│ Header: Reload | Save                               │
├──────────────────────┬──────────────────────────────┤
│ HTML Editor          │ Preview & Test               │
│                      │ ┌──────────────────────────┐ │
│ <textarea>           │ │ Preview Frame (576px)    │ │
│   HTML template      │ │                          │ │
│   with Go syntax     │ │ [Rendered HTML]          │ │
│ </textarea>          │ │                          │ │
│                      │ └──────────────────────────┘ │
│                      │                              │
│                      │ Test Print Section:          │
│                      │ - Printer IP                 │
│                      │ - Order search               │
│                      │ - Test Print | Preview PNG   │
└──────────────────────┴──────────────────────────────┘
```

### Backend

**Handler:** `backend/interfaces/http/html_template_handler.go`

**Endpoints:**
- `GET /api/manager/html-templates/bill` - Load template
- `PUT /api/manager/html-templates/bill` - Save template
- `POST /api/manager/html-templates/test-print` - Test print với order
- `POST /api/manager/html-templates/preview` - Preview PNG với order

**Template Location:**
`backend/application/services/templates/bill_template_optimized.html`

## Workflow

### 1. Load Template

```
User opens Templates tab
→ Frontend: GET /api/manager/html-templates/bill
→ Backend: Read template file
→ Response: { content: "HTML...", path: "..." }
→ Frontend: Display in editor + preview
```

### 2. Edit & Preview

```
User edits HTML
→ Debounced update (500ms)
→ Process Go template syntax
→ Replace with sample data
→ Render in iframe
→ Auto-adjust iframe height
```

### 3. Save Template

```
User clicks "Lưu Template"
→ Frontend: PUT /api/manager/html-templates/bill
→ Backend:
  1. Backup current template (.backup)
  2. Write new template
→ Response: { success: true, backup: "..." }
```

### 4. Test Print

```
User:
  1. Enters printer IP
  2. Searches for order
  3. Selects order
  4. Clicks "Test Print"
→ Frontend: POST /api/manager/html-templates/test-print
→ Backend:
  1. Fetch order from database
  2. Fetch shop settings
  3. Render template with real data (Chromedp)
  4. Convert to ESC/POS
  5. Send to printer
→ Response: { success: true, order_number: "..." }
```

### 5. Preview PNG

```
User clicks "Preview PNG"
→ Frontend: POST /api/manager/html-templates/preview
→ Backend:
  1. Fetch order from database
  2. Fetch shop settings
  3. Render template with real data (Chromedp)
  4. Save as PNG file
→ Response: { success: true, filename: "preview_html_template_[order].png" }
```

## Go Template Syntax Support

Template sử dụng Go template syntax:

### Variables
```html
{{.ShopName}}
{{.OrderNumber}}
{{.Total}}
```

### Conditionals
```html
{{if .ShowLogo}}
  <img src="{{.LogoPath}}" />
{{end}}
```

### Loops
```html
{{range $index, $item := .Items}}
  <div>{{add $index 1}}. {{$item.Name}}</div>
{{end}}
```

### Functions
```html
{{formatMoney .Total}}
{{add $index 1}}
```

## Template Data Structure

```go
type BillTemplateDataOptimized struct {
    ShopName      string
    ShopAddress   string
    ShopPhone     string
    LogoBase64    string
    ShowLogo      bool
    ShowAddress   bool
    ShowPhone     bool
    OrderNumber   string
    WaiterName    string
    PaymentMethod string
    CreatedDate   string
    Items         []BillItemDataOptimized
    Total         string
    CustomMessage string
    ShowCustomMsg bool
}

type BillItemDataOptimized struct {
    STT       int
    Name      string
    Quantity  int
    UnitPrice string
    Total     string
}
```

## Frontend Template Processing

Frontend xử lý Go template syntax để preview:

```javascript
const processTemplate = (html) => {
  // Replace variables
  html = html.replace(/\{\{\.ShopName\}\}/g, sampleData.ShopName)
  
  // Handle conditionals
  html = html.replace(/\{\{if \.ShowLogo\}\}[\s\S]*?\{\{end\}\}/g, '')
  
  // Handle loops
  const itemTemplate = extractItemTemplate(html)
  let itemsHtml = ''
  sampleData.Items.forEach((item, index) => {
    itemsHtml += renderItem(itemTemplate, item, index)
  })
  html = html.replace(/\{\{range.*?\}\}[\s\S]*?\{\{end\}\}/, itemsHtml)
  
  return html
}
```

## Integration với Chromedp

Template được render bởi `ChromedpBillRendererOptimized`:

1. Parse Go template
2. Execute với order data
3. Chromedp capture HTML → PNG
4. Binarization (đen/trắng)
5. Convert PNG → ESC/POS
6. Send to printer

## Files Changed/Created

### Frontend
- ✅ `frontend/src/components/printing/HTMLTemplateEditor.vue` (UPDATED)
  - Added test print section
  - Added order search/selection
  - Integrated with backend APIs

### Backend
- ✅ `backend/interfaces/http/html_template_handler.go` (NEW)
  - Template CRUD operations
  - Test print with real orders
  - Preview generation
- ✅ `backend/interfaces/http/chromedp_print_handler.go` (UPDATED)
  - Added GetRenderer() method
- ✅ `backend/main.go` (UPDATED)
  - Initialize htmlTemplateHandler
  - Register routes

## API Reference

### GET /api/manager/html-templates/bill

Load current template.

**Response:**
```json
{
  "success": true,
  "content": "<!DOCTYPE html>...",
  "path": "./backend/application/services/templates/bill_template_optimized.html",
  "filename": "bill_template_optimized.html"
}
```

### PUT /api/manager/html-templates/bill

Save template.

**Request:**
```json
{
  "content": "<!DOCTYPE html>..."
}
```

**Response:**
```json
{
  "success": true,
  "message": "Template saved successfully",
  "backup": "./backend/application/services/templates/bill_template_optimized.html.backup"
}
```

### POST /api/manager/html-templates/test-print

Test print with real order.

**Request:**
```json
{
  "order_id": "507f1f77bcf86cd799439011",
  "printer_ip": "192.168.1.115"
}
```

**Response:**
```json
{
  "success": true,
  "message": "Test print successful",
  "order_number": "20260222-095703-168"
}
```

### POST /api/manager/html-templates/preview

Generate preview PNG with real order.

**Request:**
```json
{
  "order_id": "507f1f77bcf86cd799439011"
}
```

**Response:**
```json
{
  "success": true,
  "message": "Preview created successfully",
  "filename": "preview_html_template_20260222-095703-168.png",
  "order_number": "20260222-095703-168"
}
```

## Usage Guide

### 1. Edit Template

1. Mở Print Management → Tab "Templates"
2. Click "🌐 HTML Template"
3. Edit HTML/CSS trong editor
4. Preview tự động update (debounced 500ms)
5. Click "💾 Lưu Template" để save

### 2. Test Print

1. Nhập IP máy in (ví dụ: 192.168.1.115)
2. Tìm order bằng search box
3. Click vào order để chọn
4. Click "🖨️ Test Print"
5. Kiểm tra output trên máy in

### 3. Preview PNG

1. Chọn order (như trên)
2. Click "👁️ Preview PNG"
3. File PNG được tạo trong thư mục backend
4. Mở file để xem kết quả

## Best Practices

### Template Design

1. **Fixed Width**: Luôn dùng 576px width
2. **Black & White**: Tránh màu sắc, chỉ dùng đen/trắng
3. **Simple Fonts**: Dùng Arial, sans-serif
4. **No External Resources**: Embed images as base64
5. **Test Frequently**: Test print thường xuyên

### CSS Guidelines

```css
/* Good */
body {
  width: 576px;
  background: white;
  color: black;
  font-family: Arial, sans-serif;
}

/* Avoid */
body {
  width: 100%; /* Don't use percentage */
  background: #f0f0f0; /* Don't use gray */
  color: #333; /* Use pure black */
  font-family: 'Custom Font'; /* May not render */
}
```

### Performance Tips

1. **Minimize HTML**: Giảm kích thước HTML
2. **Inline CSS**: Đừng dùng external stylesheets
3. **Optimize Images**: Compress logo trước khi embed
4. **Cache Template**: Backend cache parsed template

## Troubleshooting

### Template không load
- Kiểm tra file path trong backend
- Kiểm tra permissions (chmod 644)
- Xem backend logs

### Preview không hiển thị đúng
- Kiểm tra Go template syntax
- Xem browser console
- Refresh preview manually

### Test print failed
- Kiểm tra printer IP và port
- Ping printer: `ping 192.168.1.115`
- Kiểm tra order ID có tồn tại
- Xem backend logs

### Template save failed
- Kiểm tra file permissions
- Kiểm tra disk space
- Backup file có được tạo không

## Security Considerations

1. **Input Validation**: Backend validate HTML content
2. **File Permissions**: Template file chỉ writable bởi backend
3. **Backup**: Tự động backup trước khi save
4. **Access Control**: Chỉ manager role có thể edit
5. **XSS Prevention**: Sanitize user input trong template

## Future Enhancements

1. **Multiple Templates**: Hỗ trợ nhiều templates (bill, receipt, invoice)
2. **Template Library**: Pre-built templates để chọn
3. **Version Control**: Git-like versioning cho templates
4. **Template Validation**: Validate syntax trước khi save
5. **Live Collaboration**: Multiple users edit cùng lúc
6. **Template Marketplace**: Share/download templates
7. **A/B Testing**: Test multiple templates
8. **Analytics**: Track template performance

## Kết luận

Đã implement thành công tính năng quản lý HTML templates với:
- ✅ Load/Save templates
- ✅ Live preview
- ✅ Test print với order thật
- ✅ Preview PNG generation
- ✅ Integration với Chromedp renderer
- ✅ Full CRUD operations

Tính năng này cho phép user dễ dàng customize bill layout mà không cần code Go, chỉ cần edit HTML/CSS.
