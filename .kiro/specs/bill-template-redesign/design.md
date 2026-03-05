# Tài liệu Thiết kế: Thiết kế lại Template In Bill

## Tổng quan

Feature này thiết kế lại template in hóa đơn (bill) để cải thiện bố cục, tính nhất quán và tính chuyên nghiệp. Template mới sẽ:
- Hiển thị logo ở góc trên bên trái
- Tổ chức các món trong bảng có cấu trúc rõ ràng
- Sử dụng font size đồng đều trong toàn bộ hóa đơn
- Tích hợp mượt mà với hệ thống in hiện có

Thiết kế này tận dụng các module hiện có (TextRenderer, FormatParser, ImageConverter) và mở rộng chúng để hỗ trợ các tính năng mới.

## Kiến trúc

### Luồng xử lý chính


```
Order Payment (PAID)
    ↓
Print Service tạo Print Job
    ↓
Template Renderer nhận template + order data
    ↓
[NEW] Logo Renderer: Load và render logo (nếu có)
    ↓
[NEW] Table Formatter: Format items thành bảng
    ↓
[ENHANCED] Format Parser: Parse với font size rules
    ↓
Text Renderer: Render text với font size đồng đều
    ↓
[NEW] Image Compositor: Kết hợp logo + text content
    ↓
Image Converter: Convert sang ESC/POS
    ↓
Gửi đến máy in
```

### Các thành phần chính

1. **Logo Renderer** (MỚI)
   - Load logo từ file system
   - Resize logo để phù hợp với paper width
   - Convert sang grayscale
   - Xử lý lỗi gracefully

2. **Table Formatter** (MỚI)
   - Format danh sách items thành bảng
   - Tính toán độ rộng cột
   - Căn chỉnh text trong các cột
   - Xử lý text wrapping cho tên món dài

3. **Font Size Manager** (MỚI)
   - Quản lý font size cho các loại nội dung
   - Đảm bảo tính nhất quán
   - Hỗ trợ 3 mức: normal (18pt), header (22pt), total (20pt)

4. **Image Compositor** (MỚI)
   - Kết hợp logo image và text content image
   - Xử lý layout và spacing
   - Tạo ra image cuối cùng để in

5. **Enhanced Format Parser** (CẢI TIẾN)
   - Thêm logic nhận diện table rows
   - Thêm font size detection
   - Giữ nguyên logic hiện có

6. **Enhanced Text Renderer** (CẢI TIẾN)
   - Hỗ trợ nhiều font sizes
   - Render table với alignment chính xác
   - Giữ nguyên logic hiện có

## Các thành phần và Giao diện

### 1. Logo Renderer

```go
// LogoRenderer xử lý việc load và render logo
type LogoRenderer struct {
    maxWidthPercent float64  // % của paper width (default: 25%)
    margin          int      // Margin xung quanh logo
}

// RenderLogo loads logo từ path và render thành grayscale image
// Input: logoPath (đường dẫn file), paperWidth (pixel width của giấy)
// Output: *image.Gray (logo đã resize), error
func (r *LogoRenderer) RenderLogo(logoPath string, paperWidth int) (*image.Gray, error)

// resizeLogo resize logo để fit trong maxWidth
func (r *LogoRenderer) resizeLogo(img image.Image, maxWidth int) *image.Gray

// convertToGrayscale convert color image sang grayscale
func (r *LogoRenderer) convertToGrayscale(img image.Image) *image.Gray
```

**Thuật toán resize logo:**
```
maxLogoWidth = paperWidth * maxWidthPercent
if logoWidth > maxLogoWidth:
    scale = maxLogoWidth / logoWidth
    newWidth = maxLogoWidth
    newHeight = logoHeight * scale
    resize image using bilinear interpolation
```

### 2. Table Formatter

```go
// TableFormatter format items thành bảng
type TableFormatter struct {
    paperWidth  int
    margin      int
    columnGap   int  // Khoảng cách giữa các cột
}

// TableColumn định nghĩa một cột trong bảng
type TableColumn struct {
    Header    string
    Width     int        // Width in characters
    Alignment Alignment  // Left, Center, Right
}

// TableRow định nghĩa một dòng trong bảng
type TableRow struct {
    Cells []string
}

// FormatItemsTable format order items thành table lines
// Input: items (danh sách món), paperWidth
// Output: []string (các dòng text đã format)
func (f *TableFormatter) FormatItemsTable(items []OrderItem, paperWidth int) []string

// calculateColumnWidths tính toán độ rộng tối ưu cho các cột
func (f *TableFormatter) calculateColumnWidths(paperWidth int) []TableColumn

// formatRow format một row với column widths
func (f *TableFormatter) formatRow(row TableRow, columns []TableColumn) string

// wrapCellText wrap text trong cell nếu quá dài
func (f *TableFormatter) wrapCellText(text string, maxWidth int) []string
```

**Cấu trúc bảng:**
```
Cột 1: Tên món (50% width, align left)
Cột 2: SL (15% width, align right)
Cột 3: Đơn giá (17.5% width, align right)
Cột 4: Thành tiền (17.5% width, align right)
```

**Ví dụ output:**
```
Tên món              SL  Đơn giá    Thành tiền
================================================
Cafe Latte            2   45,000       90,000
Banh Mi Thit          1   35,000       35,000
  (Variant: Đặc biệt)
```

### 3. Font Size Manager

```go
// FontSizeManager quản lý font sizes cho các loại content
type FontSizeManager struct {
    normalSize  float64  // 18pt
    headerSize  float64  // 22pt
    totalSize   float64  // 20pt
}

// FontSizeConfig định nghĩa font size cho một line
type FontSizeConfig struct {
    Size   float64
    Bold   bool
}

// GetFontSizeForLine xác định font size cho một line dựa trên content
func (m *FontSizeManager) GetFontSizeForLine(line string) FontSizeConfig
```

**Quy tắc font size:**
- Header (tên shop, "HÓA ĐƠN BÁN HÀNG"): 22pt, bold
- Total line ("TỔNG CỘNG"): 20pt, bold
- Table header: 18pt, bold
- Table content: 18pt, normal
- Thông tin đơn hàng: 18pt, normal
- Footer: 18pt, normal

### 4. Image Compositor

```go
// ImageCompositor kết hợp logo và text content
type ImageCompositor struct {
    paperWidth int
    margin     int
}

// Compose kết hợp logo (optional) và text content thành một image
// Input: logo (*image.Gray, có thể nil), textContent (*image.Gray)
// Output: *image.Gray (combined image)
func (c *ImageCompositor) Compose(logo *image.Gray, textContent *image.Gray) (*image.Gray, error)

// calculateTotalHeight tính tổng height cần thiết
func (c *ImageCompositor) calculateTotalHeight(logo *image.Gray, textContent *image.Gray) int

// drawLogo vẽ logo lên combined image
func (c *ImageCompositor) drawLogo(dst *image.Gray, logo *image.Gray, x, y int)

// drawTextContent vẽ text content lên combined image
func (c *ImageCompositor) drawTextContent(dst *image.Gray, textContent *image.Gray, x, y int)
```

**Layout logic:**
```
if logo exists:
    [margin]
    [logo - aligned left with margin]
    [spacing: 20px]
    [text content]
    [margin]
else:
    [text content as-is]
```

### 5. Enhanced Format Parser

Mở rộng `FormatParser` hiện có:

```go
// Thêm vào LineFormat struct
type LineFormat struct {
    Text        string
    Bold        bool
    Alignment   Alignment
    IsSeparator bool
    FontSize    float64    // NEW: Font size for this line
    IsTableRow  bool       // NEW: Đánh dấu dòng thuộc table
}

// Thêm method mới
func (p *FormatParser) detectFontSize(line string) float64
func (p *FormatParser) isTableRow(line string) bool
```

**Logic detect font size:**
```
if line contains "HÓA ĐƠN" or is shop name:
    return 22pt
else if line contains "TỔNG CỘNG":
    return 20pt
else:
    return 18pt
```

### 6. Enhanced Text Renderer

Mở rộng `TextRenderer` hiện có:

```go
// Thêm vào RendererConfig
type RendererConfig struct {
    PixelWidth  int
    FontPath    string
    FontSize    float64
    LineSpacing int
    Margin      int
    FontSizes   map[string]float64  // NEW: Map of font sizes
}

// Thêm vào TextRenderer
type TextRenderer struct {
    pixelWidth  int
    fonts       map[float64]FontPair  // NEW: Multiple font sizes
    fontSize    float64
    lineSpacing int
    margin      int
}

type FontPair struct {
    Normal font.Face
    Bold   font.Face
}

// Modify existing method
func (r *TextRenderer) drawLine(img *image.Gray, line LineFormat, y int) int {
    // Select font based on line.FontSize instead of just bold flag
    fontPair := r.fonts[line.FontSize]
    fontFace := fontPair.Normal
    if line.Bold {
        fontFace = fontPair.Bold
    }
    // ... rest of logic
}
```

## Mô hình dữ liệu

### Template Data Structure

Template mới sẽ sử dụng cấu trúc dữ liệu hiện có với một số thêm mới:

```go
type BillTemplateData struct {
    // Shop info (existing)
    ShopName           string
    ShopAddress        string
    ShopPhone          string
    ShowAddress        bool
    ShowPhone          bool
    ShowLogo           bool    // NEW
    LogoPath           string  // NEW
    CustomMessage      string
    ShowCustomMessage  bool
    
    // Order info (existing)
    Order OrderData
}

type OrderData struct {
    OrderNumber    string
    CreatedAt      time.Time
    TableNumber    string
    CustomerName   string
    Items          []OrderItem
    Subtotal       float64
    DiscountAmount float64
    TaxAmount      float64
    Total          float64
}

type OrderItem struct {
    Name         string
    VariantName  string
    Quantity     int
    UnitPrice    float64
    TotalPrice   float64
}
```

### Shop Settings Extension

ShopSettings đã có sẵn các fields cần thiết:
- `LogoURL string` - URL/path của logo
- `ShowLogo bool` - Flag để show/hide logo

Không cần thêm fields mới.

### Template Content Format

Template mới sẽ có format như sau:

```
{{if .ShowLogo}}[LOGO]{{end}}

{{.ShopName}}
{{if .ShowAddress}}{{.ShopAddress}}{{end}}
{{if .ShowPhone}}Tel: {{.ShopPhone}}{{end}}
================================
HÓA ĐƠN BÁN HÀNG
Order: {{.Order.OrderNumber}}
Ngày: {{.Order.CreatedAt.Format "02/01/2006 15:04"}}
{{if .Order.TableNumber}}Bàn: {{.Order.TableNumber}}{{end}}
{{if .Order.CustomerName}}Khách: {{.Order.CustomerName}}{{end}}
================================
[TABLE_START]
Tên món              SL  Đơn giá    Thành tiền
------------------------------------------------
{{range .Order.Items}}
{{.Name}}{{if .VariantName}}
  ({{.VariantName}}){{end}}  {{.Quantity}}  {{formatPrice .UnitPrice}}  {{formatPrice .TotalPrice}}
{{end}}
[TABLE_END]
================================
Tổng tiền: {{formatPrice .Order.Subtotal}} VND
{{if gt .Order.DiscountAmount 0.0}}Giảm giá: -{{formatPrice .Order.DiscountAmount}} VND{{end}}
{{if gt .Order.TaxAmount 0.0}}Thuế: {{formatPrice .Order.TaxAmount}} VND{{end}}
TỔNG CỘNG: {{formatPrice .Order.Total}} VND
================================
{{if .ShowCustomMessage}}{{.CustomMessage}}{{end}}
Cảm ơn quý khách!
Hẹn gặp lại!
```

**Markers đặc biệt:**
- `[LOGO]` - Placeholder cho logo, sẽ được replace bởi LogoRenderer
- `[TABLE_START]` và `[TABLE_END]` - Đánh dấu vùng table, sẽ được xử lý bởi TableFormatter

## Correctness Properties

*A property is a characteristic or behavior that should hold true across all valid executions of a system—essentially, a formal statement about what the system should do. Properties serve as the bridge between human-readable specifications and machine-verifiable correctness guarantees.*


### Property 1: Logo hiển thị ở góc trên bên trái

*For any* hóa đơn có logo được cấu hình, khi render hóa đơn, logo phải xuất hiện ở góc trên bên trái của image với x-coordinate bằng margin và y-coordinate ở vị trí đầu tiên.

**Validates: Requirements 1.1, 1.5**

### Property 2: Logo được load từ đường dẫn đã cấu hình

*For any* đường dẫn logo hợp lệ được cấu hình trong shop settings, hệ thống phải load logo từ đúng đường dẫn đó.

**Validates: Requirements 1.2**

### Property 3: Xử lý graceful khi không có logo

*For any* hóa đơn không có logo được cấu hình, hệ thống phải render hóa đơn thành công mà không có logo và không có lỗi.

**Validates: Requirements 1.3**

### Property 4: Logo sizing constraint

*For any* logo và paper width, chiều rộng của logo sau khi resize phải <= 25% của paper width.

**Validates: Requirements 1.4, 1.7**

### Property 5: Logo format support

*For any* file logo có định dạng PNG, JPG, hoặc JPEG, hệ thống phải load và render thành công.

**Validates: Requirements 1.6**

### Property 6: Logo grayscale conversion

*For any* logo (color hoặc grayscale), output image của logo phải là grayscale (tất cả pixels có R=G=B).

**Validates: Requirements 1.8**

### Property 7: Table structure completeness

*For any* order có items, bảng món phải chứa 4 cột: Tên món, Số lượng, Đơn giá, và Thành tiền.

**Validates: Requirements 2.1**

### Property 8: Table column alignment consistency

*For any* bảng món, tên món phải căn trái, và số lượng, đơn giá, thành tiền phải căn phải, với các cột được căn chỉnh đồng đều trong toàn bộ bảng.

**Validates: Requirements 2.2, 2.3, 2.8**

### Property 9: Table separator lines

*For any* bảng món, phải có đường kẻ ngang phân tách giữa header và các dòng món.

**Validates: Requirements 2.4**

### Property 10: Long item name wrapping

*For any* item có tên dài vượt quá chiều rộng cột, tên món phải được wrap sang dòng mới trong cùng cột.

**Validates: Requirements 2.5**

### Property 11: Variant display on sub-line

*For any* item có variant, variant phải được hiển thị trên dòng phụ bên dưới tên món.

**Validates: Requirements 2.6**

### Property 12: Column width calculation

*For any* paper width, tổng chiều rộng các cột (bao gồm gaps) phải <= paper width - 2*margin.

**Validates: Requirements 2.7**

### Property 13: Font size consistency for regular content

*For any* nội dung thông thường (table content, order info, footer), font size phải là 18pt.

**Validates: Requirements 3.1, 3.4, 3.5**

### Property 14: Header font size

*For any* header line (tên shop, "HÓA ĐƠN BÁN HÀNG"), font size phải là 22pt.

**Validates: Requirements 3.2**

### Property 15: Total line font size

*For any* dòng tổng tiền ("TỔNG CỘNG"), font size phải là 20pt.

**Validates: Requirements 3.3**

### Property 16: Font weight for headers and totals

*For any* header line hoặc total line, font weight phải là bold.

**Validates: Requirements 3.6**

### Property 17: Font weight for regular content

*For any* nội dung thông thường (không phải header hoặc total), font weight phải là normal.

**Validates: Requirements 3.7**

### Property 18: Default template usage

*For any* template được đặt làm default, tất cả print jobs mới phải sử dụng template đó.

**Validates: Requirements 4.2, 4.7**

### Property 19: Template variable support

*For any* template data hợp lệ, tất cả các biến template (ShopName, Order.Items, Order.Total, v.v.) phải được render đúng giá trị.

**Validates: Requirements 4.5**

### Property 20: Paper width compatibility

*For any* paper width (58mm hoặc 80mm), template phải render thành công và fit trong paper width.

**Validates: Requirements 4.6**

### Property 21: Logo path persistence

*For any* logo được upload, đường dẫn logo phải được lưu vào shop_settings collection.

**Validates: Requirements 5.2**

### Property 22: Logo format validation

*For any* file upload, nếu định dạng không phải PNG, JPG, hoặc JPEG, upload phải bị reject.

**Validates: Requirements 5.5**

### Property 23: Logo file size validation

*For any* file upload, nếu kích thước > 2MB, upload phải bị reject.

**Validates: Requirements 5.6**

### Property 24: Logo storage location

*For any* logo được upload thành công, file phải được lưu trong thư mục uploads/logos/.

**Validates: Requirements 5.7**

### Property 25: Template preservation

*For any* template mới được tạo, tất cả templates cũ phải vẫn tồn tại trong database.

**Validates: Requirements 6.1, 6.4**

### Property 26: Template switching

*For any* template được chọn (mới hoặc cũ), print jobs phải sử dụng đúng template đó.

**Validates: Requirements 6.2, 6.3**

### Property 27: Missing logo graceful handling

*For any* logo path được cấu hình nhưng file không tồn tại, hệ thống phải render hóa đơn thành công mà không có logo.

**Validates: Requirements 7.1**

### Property 28: Corrupt logo graceful handling

*For any* logo file bị corrupt, hệ thống phải log lỗi và render hóa đơn thành công mà không có logo.

**Validates: Requirements 7.2**

### Property 29: Template rendering fallback

*For any* template rendering failure, hệ thống phải fallback về template mặc định cũ và render thành công.

**Validates: Requirements 7.3**

### Property 30: Error logging completeness

*For any* lỗi liên quan đến logo hoặc template rendering, lỗi phải được log với đầy đủ thông tin (timestamp, error message, context).

**Validates: Requirements 7.4**

### Property 31: Extreme logo size handling

*For any* logo quá lớn không thể resize, hệ thống phải bỏ qua logo và render phần còn lại thành công.

**Validates: Requirements 7.5**

## Xử lý lỗi

### Lỗi liên quan đến Logo

1. **Logo file không tồn tại**
   - Log: "Logo file not found: {path}"
   - Action: Tiếp tục render mà không có logo
   - User notification: Không cần (silent fallback)

2. **Logo file corrupt**
   - Log: "Failed to load logo: {path}, error: {error}"
   - Action: Tiếp tục render mà không có logo
   - User notification: Không cần (silent fallback)

3. **Logo quá lớn**
   - Log: "Logo resize failed: {path}, size: {size}"
   - Action: Tiếp tục render mà không có logo
   - User notification: Không cần (silent fallback)

4. **Logo format không hợp lệ (upload)**
   - Log: "Invalid logo format: {filename}, format: {format}"
   - Action: Reject upload
   - User notification: "Chỉ hỗ trợ định dạng PNG, JPG, JPEG"

5. **Logo file quá lớn (upload)**
   - Log: "Logo file too large: {filename}, size: {size}MB"
   - Action: Reject upload
   - User notification: "Kích thước file tối đa 2MB"

### Lỗi liên quan đến Template

1. **Template rendering failure**
   - Log: "Template rendering failed: {template_id}, error: {error}"
   - Action: Fallback to default template
   - User notification: "Lỗi render template, sử dụng template mặc định"

2. **Template not found**
   - Log: "Template not found: {template_id}"
   - Action: Use default template
   - User notification: Không cần

3. **Invalid template data**
   - Log: "Invalid template data: {error}"
   - Action: Use default values
   - User notification: Không cần

### Lỗi liên quan đến Table Formatting

1. **Item name quá dài**
   - Action: Wrap text trong cột
   - No error, handled gracefully

2. **Column width calculation overflow**
   - Log: "Column width overflow: paper_width={width}"
   - Action: Adjust column widths proportionally
   - No user notification

## Chiến lược Testing

### Dual Testing Approach

Feature này sẽ sử dụng cả unit tests và property-based tests:

**Unit Tests:**
- Test specific examples của logo rendering
- Test specific table layouts
- Test error handling với specific error cases
- Test integration giữa các components

**Property-Based Tests:**
- Test universal properties với random inputs
- Minimum 100 iterations per property test
- Sử dụng thư viện: `gopter` (Go property testing)
- Mỗi property test phải tag với comment:
  ```go
  // Feature: bill-template-redesign, Property 1: Logo hiển thị ở góc trên bên trái
  ```

### Test Coverage

**Logo Renderer:**
- Unit tests: Load PNG/JPG/JPEG, resize logic, grayscale conversion
- Property tests: Property 1-6 (logo positioning, sizing, format support)

**Table Formatter:**
- Unit tests: Specific table layouts, column calculations
- Property tests: Property 7-12 (table structure, alignment, wrapping)

**Font Size Manager:**
- Unit tests: Specific font size assignments
- Property tests: Property 13-17 (font size consistency)

**Image Compositor:**
- Unit tests: Compose with/without logo
- Property tests: Combined with logo and text rendering properties

**Template Integration:**
- Unit tests: Template rendering với specific data
- Property tests: Property 18-20, 25-26 (template usage, switching)

**Logo Upload:**
- Unit tests: Upload flow, validation
- Property tests: Property 21-24 (persistence, validation, storage)

**Error Handling:**
- Unit tests: Specific error scenarios
- Property tests: Property 27-31 (graceful degradation)

### Test Data Generation

Sử dụng generators cho property tests:

```go
// Logo generators
GenLogoPath() // Random valid/invalid paths
GenLogoImage() // Random images với different sizes/formats
GenPaperWidth() // 58mm (448px) hoặc 80mm (576px)

// Order generators
GenOrder() // Random orders với items
GenOrderItem() // Random items với long/short names, variants
GenOrderWithManyItems() // Orders với nhiều items để test table

// Template generators
GenTemplateData() // Random template data
GenShopSettings() // Random shop settings với/không có logo
```

### Integration Testing

Test end-to-end flow:
1. Upload logo → verify saved to correct location
2. Configure shop settings → verify logo path saved
3. Create order → trigger print job
4. Render template → verify logo + table + font sizes
5. Convert to ESC/POS → verify output format
6. Switch templates → verify correct template used

### Manual Testing Checklist

- [ ] Upload logo PNG, JPG, JPEG - verify preview
- [ ] Upload logo > 2MB - verify rejection
- [ ] Upload invalid format - verify rejection
- [ ] Create order với items dài - verify wrapping
- [ ] Create order với variants - verify sub-line display
- [ ] Test với 58mm paper - verify layout
- [ ] Test với 80mm paper - verify layout
- [ ] Delete logo - verify bills render without logo
- [ ] Corrupt logo file - verify graceful fallback
- [ ] Switch giữa template mới và cũ - verify correct rendering

