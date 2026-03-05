# Debug White Preview Issue

## Problem
Preview PNG file ra trắng hoàn toàn.

## Possible Causes

### 1. Chromedp không capture được HTML
**Symptoms**: Cả raw và binarized đều trắng

**Causes**:
- Chrome không chạy
- HTML template có lỗi syntax
- Timeout quá ngắn (300ms)
- Data URL quá dài
- Font không load được

### 2. Binarization threshold quá cao
**Symptoms**: Raw có nội dung, binarized trắng

**Causes**:
- Threshold 128 quá cao
- Tất cả pixels > 128 → trở thành trắng (255)
- Chỉ pixels rất tối (< 128) mới thành đen

## Debug Steps

### Step 1: Check Raw Capture

Code đã được update để save 2 files:
- `raw_preview_html_template_*.png` - Raw capture (chưa binarize)
- `preview_html_template_*.png` - Sau binarize

**Test**:
```bash
# Create preview
# (use frontend or curl)

# Check files
ls -lh backend/raw_preview_*.png
ls -lh backend/preview_*.png

# View raw file
open backend/raw_preview_*.png  # macOS
```

**Expected**:
- Raw file có nội dung bill đầy đủ
- Binarized file có thể trắng nếu threshold sai

### Step 2: If Raw is White

Chromedp không capture được. Check:

1. **Chrome process**:
```bash
ps aux | grep chrome
```

2. **Backend logs**:
```bash
tail -f backend.log | grep -i "chromedp\|capture\|error"
```

3. **Template syntax**:
```bash
# Check template file
cat backend/application/services/templates/bill_template_optimized.html
```

4. **Test simple HTML**:
Tạo test với HTML đơn giản:
```html
<!DOCTYPE html>
<html>
<body style="width:576px;background:white;">
  <h1>TEST</h1>
  <p>Hello World</p>
</body>
</html>
```

### Step 3: If Raw Has Content but Binarized is White

Binarization threshold quá cao.

**Current threshold**: 128 (0-255 scale)

**Logic**:
```go
if gray > 128 {
    pixel = WHITE (255)
} else {
    pixel = BLACK (0)
}
```

**Problem**: 
- Background trắng (255) > 128 → WHITE ✓
- Text đen (0) < 128 → BLACK ✓
- BUT: Nếu text có anti-aliasing (gray 100-200) → WHITE ❌

**Solution**: Lower threshold

```go
// Try threshold 64 or 96
const ThresholdValue = 64  // More aggressive
```

### Step 4: Check Chromedp Settings

Current settings:
```go
chromedp.Flag("headless", true),
chromedp.Flag("disable-gpu", true),
chromedp.Flag("no-sandbox", true),
```

**Add more flags**:
```go
chromedp.Flag("disable-dev-shm-usage", true),
chromedp.Flag("disable-software-rasterizer", true),
chromedp.WindowSize(576, 1200),
```

### Step 5: Increase Wait Time

Current: 300ms

```go
chromedp.Sleep(300*time.Millisecond)
```

**Try**: 1000ms or 2000ms

```go
chromedp.Sleep(2*time.Second)
```

## Quick Fixes

### Fix 1: Lower Binarization Threshold

```go
// In chromedp_bill_renderer_optimized.go
const ThresholdValue = 64  // Changed from 128
```

### Fix 2: Disable Binarization (Test Only)

```go
// In SavePreviewImage
// Comment out binarization
// bwImg := binarizeImageOptimized(img, ThresholdValue)

// Save raw image directly
return png.Encode(f, img)
```

### Fix 3: Increase Chromedp Timeout

```go
// In captureHTML
chromedp.Sleep(2*time.Second),  // Changed from 300ms
```

### Fix 4: Add Viewport Size

```go
// In NewChromedpBillRendererOptimized
opts := append(chromedp.DefaultExecAllocatorOptions[:],
    chromedp.Flag("headless", true),
    chromedp.Flag("disable-gpu", true),
    chromedp.Flag("no-sandbox", true),
    chromedp.WindowSize(576, 1200),  // ADD THIS
)
```

## Testing

### Test 1: Check Raw Capture
```bash
# After creating preview
ls -lh backend/raw_preview_*.png
file backend/raw_preview_*.png
```

**Expected**: PNG file with actual content

### Test 2: Check Binarized
```bash
ls -lh backend/preview_*.png
file backend/preview_*.png
```

**Expected**: PNG file with black & white content

### Test 3: Compare with preview.go
```bash
cd backend/cmd/test-vietnamese-print
go run . preview.go main.go
open preview_bill.png
```

**Expected**: This should work (uses gg library directly)

## Recommended Solution

Based on the issue, most likely:

1. **If raw is white**: Chromedp issue
   - Increase timeout to 2s
   - Add WindowSize
   - Check Chrome process

2. **If raw has content**: Threshold issue
   - Lower threshold to 64 or 96
   - Or remove binarization entirely
   - Use grayscale instead

## Code Changes

### Change 1: Lower Threshold
```go
const ThresholdValue = 64  // More aggressive binarization
```

### Change 2: Save Raw for Debug
```go
// Already done - saves raw_preview_*.png
```

### Change 3: Increase Timeout
```go
chromedp.Sleep(2*time.Second),
```

## Next Steps

1. ✅ Create preview again
2. ✅ Check `raw_preview_*.png` file
3. ⚠️ If raw is white → Fix chromedp
4. ⚠️ If raw has content → Lower threshold
5. ✅ Test print after fix
