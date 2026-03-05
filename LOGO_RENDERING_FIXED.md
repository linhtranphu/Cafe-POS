# Logo Rendering Fix - COMPLETED ✅

## Vấn đề
Logo không hiển thị trong bill khi render với chromedp, mặc dù logo base64 đã được load thành công.

## Nguyên nhân
Go's `html/template` package tự động escape HTML attributes để bảo mật. Khi sử dụng `string` type cho LogoBase64, template engine escape data URL thành `#ZgotmplZ` (error marker).

## Giải pháp

### 1. Sử dụng `template.URL` type
Thay đổi type của `LogoBase64` từ `string` sang `template.URL`:

```go
type BillTemplateDataOptimized struct {
    // ...
    LogoBase64    template.URL // Changed from string to template.URL
    // ...
}
```

### 2. Convert khi assign value
```go
data.LogoBase64 = template.URL(logoBase64) // Convert to template.URL for src attribute
```

### 3. Template HTML không cần thay đổi
```html
<img src="{{.LogoBase64}}" alt="Logo">
```

## Kết quả

✅ Logo hiển thị đúng trong rendered HTML
✅ Logo render trong chromedp screenshot
✅ Logo in ra bill ESC/POS
✅ Hoạt động với logo thực tế từ uploads/

## Test Results

### Test với logo uploaded (702x374 JPEG)
```
Logo resized: 702x374 → 200x106
Logo base64: 51822 bytes (original: 107481 bytes)
✅ Logo base64 found in rendered HTML
Screenshot captured: 96348 bytes
✅ Preview saved successfully
```

### Files đã sửa
1. `backend/application/services/chromedp_bill_renderer_optimized.go`
   - Changed `LogoBase64` type to `template.URL`
   - Convert string to `template.URL` when assigning
   - Added debug logging

2. `backend/interfaces/http/order_handler.go`
   - Fixed unused import issue

## Chromedp Optimization (Bonus)

Đã áp dụng các fix để capture "nét căng":

1. **Data URL approach** thay vì temp file
2. **4 chốt chặn** để đảm bảo render hoàn chỉnh:
   - WaitReady("body")
   - Promise đợi images load
   - Sleep 200ms cho font/CSS
   - FullScreenshot quality 100

## Testing

Chạy test với logo thực tế:
```bash
cd backend
go run cmd/test-uploaded-logo/main.go
```

Output files:
- `test_uploaded_logo_preview.png` - Preview image
- `test_uploaded_logo.bin` - ESC/POS data
- `debug_rendered.html` - HTML for inspection

## Notes

- `template.URL` type cho phép data URL trong `src` attribute
- `template.HTML` type cho phép raw HTML trong content
- Logo được auto-resize về max width 200px để giảm kích thước base64
- Data URL approach hoạt động tốt trong Docker environment
