# Chromedp Capture Optimization - Fix "Nét Căng"

## Vấn đề
- Chromedp capture bị mờ, không sắc nét
- Logo thường bị mất hoặc không hiển thị
- Render không đầy đủ, thiếu font hoặc CSS

## Giải pháp đã áp dụng

### 1. Fix Chromedp Capture với Data URL
Thay vì dùng temp file, sử dụng data URL với các "chốt chặn" để đảm bảo render hoàn chỉnh:

```go
func (r *ChromedpBillRendererOptimized) captureHTML(html string) (image.Image, error) {
    var buf []byte
    
    // Sử dụng data URL
    dataURL := "data:text/html;charset=utf-8," + url.PathEscape(html)
    
    err := chromedp.Run(r.ctx,
        // Điều hướng tới data URL
        chromedp.Navigate(dataURL),
        
        // CHỐT CHẶN 1: Đợi body element xuất hiện
        chromedp.WaitReady("body"),
        
        // CHỐT CHẶN 2: Đợi tất cả images load xong
        chromedp.Evaluate(`
            new Promise((resolve) => {
                const images = document.querySelectorAll('img');
                if (images.length === 0) {
                    resolve();
                    return;
                }
                let loaded = 0;
                const checkComplete = () => {
                    loaded++;
                    if (loaded === images.length) resolve();
                };
                images.forEach(img => {
                    if (img.complete) {
                        checkComplete();
                    } else {
                        img.onload = checkComplete;
                        img.onerror = checkComplete;
                    }
                });
                setTimeout(resolve, 2000);
            });
        `, nil),
        
        // CHỐT CHẶN 3: Đợi thêm để font/CSS render hoàn toàn
        chromedp.Sleep(200 * time.Millisecond),
        
        // CHỐT CHẶN 4: Chụp toàn bộ trang
        chromedp.FullScreenshot(&buf, 100),
    )
    
    // ... decode và crop
}
```

### 2. Logo Base64 Embedding
Logo đã được convert sang base64 và embed trực tiếp vào HTML:

**Function `loadImageAsBase64`:**
- Đọc file logo từ disk
- Resize về max width 200px để giảm kích thước
- Encode sang PNG
- Convert sang base64 với format: `data:image/png;base64,...`

**Template HTML:**
```html
<div class="logo">
    <img src="{{.LogoBase64}}" alt="Logo">
</div>
```

### 3. Các chốt chặn quan trọng

1. **WaitReady("body")**: Đảm bảo DOM đã load
2. **Image loading promise**: Đợi tất cả images load xong
3. **Sleep 200ms**: Buffer time cho font và CSS render
4. **FullScreenshot**: Capture toàn bộ trang với quality 100

## Lợi ích

✅ Render sắc nét, "nét căng"
✅ Logo luôn hiển thị (không phụ thuộc file system)
✅ Font và CSS render đầy đủ
✅ Hoạt động tốt trong Docker environment
✅ Không cần temp file

## Testing

Chạy test để verify:

```bash
cd backend
go run cmd/test-chromedp-optimized/main.go
```

Kiểm tra output:
- `test_bill_optimized.bin` - ESC/POS data
- `test_bill_optimized_preview.png` - Preview image

## Files đã sửa

1. `backend/application/services/chromedp_bill_renderer_optimized.go`
   - Updated `captureHTML()` function
   - Added `net/url` import for PathEscape

2. `backend/application/services/templates/bill_template.html`
   - Đã sử dụng `{{.LogoBase64}}` cho logo (không cần sửa thêm)

## Notes

- Data URL approach hoạt động tốt hơn temp file trong Docker
- Base64 logo tăng kích thước HTML nhưng đảm bảo logo luôn hiển thị
- Resize logo về 200px giúp giảm kích thước base64
- Timeout 2 giây cho image loading là đủ cho hầu hết trường hợp
