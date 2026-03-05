# Test Chromedp HTML Capture

## Mục đích
Verify xem chromedp có capture HTML thành image đúng không trước khi convert sang ESC/POS.

## Cách test

### 1. Sử dụng Preview Endpoint (Recommended)

**Endpoint**: `POST /api/manager/html-templates/preview`

**Request**:
```json
{
  "order_id": "YOUR_ORDER_ID_HERE"
}
```

**Response**:
```json
{
  "success": true,
  "message": "Preview created successfully",
  "filename": "preview_html_template_20260222-095703-168.png",
  "order_number": "20260222-095703-168"
}
```

**File được tạo**: `preview_html_template_{order_number}.png` trong thư mục backend root

### 2. Sử dụng curl

```bash
# Get an order ID first
curl -H "Authorization: Bearer YOUR_TOKEN" \
  http://localhost:3000/api/manager/orders | jq '.[0].id'

# Create preview
curl -X POST \
  -H "Authorization: Bearer YOUR_TOKEN" \
  -H "Content-Type: application/json" \
  -d '{"order_id":"YOUR_ORDER_ID"}' \
  http://localhost:3000/api/manager/html-templates/preview
```

### 3. Sử dụng Frontend

1. Vào `http://localhost:5173/#/print-management`
2. Tab "Templates"
3. Chọn một order
4. Click "Preview PNG"
5. Check console log để xem filename
6. File sẽ được tạo trong thư mục backend

## Kiểm tra kết quả

### Stage 1: Raw Capture (FullScreenshot)
File preview sẽ cho thấy image sau khi:
1. ✅ Chromedp capture HTML
2. ✅ Binarization (black & white)
3. ❌ CHƯA invert (vì SavePreviewImage không invert)

**Expected**: 
- Background: WHITE
- Text: BLACK
- Logo: Visible
- Layout: Match HTML template

### Stage 2: After Inversion (sent to printer)
Khi in thật, image sẽ được invert trước khi convert sang ESC/POS:
- Background WHITE → BLACK (lightness 0.0) → không in → giấy trắng ✅
- Text BLACK → WHITE (lightness 1.0) → in đen → text đen ✅

## Các vấn đề có thể gặp

### 1. Preview PNG bị trắng hoàn toàn
**Nguyên nhân**: Chromedp không capture được HTML
**Kiểm tra**:
- Chrome có chạy không? (check process)
- HTML template có lỗi syntax không?
- Logo có load được không?

### 2. Preview PNG có nội dung nhưng layout sai
**Nguyên nhân**: HTML/CSS chưa đúng
**Kiểm tra**:
- Width có đúng 576px không?
- Font có load được không?
- CSS có lỗi không?

### 3. Preview PNG đúng nhưng in ra sai
**Nguyên nhân**: ESC/POS conversion có vấn đề
**Kiểm tra**:
- Inversion có được apply không?
- Raster converter có hoạt động đúng không?
- Threshold có phù hợp không?

## Debug Steps

### Step 1: Verify Chromedp Capture
```bash
# Create preview
curl -X POST \
  -H "Authorization: Bearer YOUR_TOKEN" \
  -H "Content-Type: application/json" \
  -d '{"order_id":"YOUR_ORDER_ID"}' \
  http://localhost:3000/api/manager/html-templates/preview

# Check file
ls -lh preview_html_template_*.png
open preview_html_template_*.png  # macOS
```

**Expected**: PNG file với nội dung bill đúng

### Step 2: Verify Binarization
Preview PNG đã được binarize (black & white only).

**Check**: 
- Không có màu xám (grayscale)
- Chỉ có đen và trắng
- Text rõ ràng, không bị mờ

### Step 3: Verify Inversion Logic
Code hiện tại:
```go
// In imageToESCPOSOptimized()
invertedImg := invertImage(img)  // ← This inverts before sending to printer
```

**Verify**: Check code có dòng này không

### Step 4: Test Print
```bash
# Test print với order
curl -X POST \
  -H "Authorization: Bearer YOUR_TOKEN" \
  -H "Content-Type: application/json" \
  -d '{"order_id":"YOUR_ORDER_ID","printer_ip":"192.168.1.115"}' \
  http://localhost:3000/api/manager/html-templates/test-print
```

## Expected Flow

```
HTML Template
    ↓
Chromedp FullScreenshot
    ↓
PNG Image (color)
    ↓
Binarization (threshold 128)
    ↓
Black & White Image
    ↓
Inversion (black ↔ white)
    ↓
Inverted B&W Image
    ↓
Raster Converter (threshold 0.5)
    ↓
ESC/POS Commands
    ↓
Printer
```

## Files to Check

1. **Template**: `backend/application/services/templates/bill_template_optimized.html`
2. **Renderer**: `backend/application/services/chromedp_bill_renderer_optimized.go`
3. **Handler**: `backend/interfaces/http/html_template_handler.go`
4. **Preview Output**: `preview_html_template_*.png` (in backend root)

## Common Issues

### Issue: "Failed to capture"
- Chrome not running
- Template syntax error
- Timeout too short

### Issue: "Preview blank"
- HTML not rendering
- CSS not loading
- Font not found

### Issue: "Preview correct but print wrong"
- Inversion not applied
- Raster logic incorrect
- Threshold wrong

## Next Steps

1. ✅ Create preview PNG
2. ✅ Verify preview looks correct
3. ✅ Verify code has inversion
4. ✅ Test actual print
5. ⚠️ Adjust threshold if needed
