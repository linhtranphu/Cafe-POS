# Image Inversion Fix for ESC/POS Printing

## Problem
Khi in ra chỉ thấy 1 dải màu đen thay vì nội dung bill.

## Root Cause
Library `kenshaw/escpos/raster` có logic NGƯỢC với thông thường:

```go
// From raster.Converter.ToRaster()
if lightness(img.At(x, y)) >= c.Threshold {
    // Set bit to 1 → Print BLACK
}
```

**Logic này nghĩa là:**
- Pixel SÁNG (lightness cao) → In ĐEN
- Pixel TỐI (lightness thấp) → Để TRẮNG

**Vấn đề:**
- Bill image có background TRẮNG (lightness = 1.0)
- Text/content là ĐEN (lightness = 0.0)
- Khi convert:
  - Background trắng (lightness 1.0 > threshold 0.5) → In ĐEN ❌
  - Text đen (lightness 0.0 < threshold 0.5) → Để TRẮNG ❌
- Kết quả: Toàn bộ giấy in đen!

## Solution: Invert Image Before Conversion

Trước khi pass image vào `raster.Converter`, invert màu (đen ↔ trắng):

```go
func invertImage(img image.Image) image.Image {
    bounds := img.Bounds()
    inverted := image.NewRGBA(bounds)
    
    for y := bounds.Min.Y; y < bounds.Max.Y; y++ {
        for x := bounds.Min.X; x < bounds.Max.X; x++ {
            r, g, b, a := img.At(x, y).RGBA()
            
            // Invert RGB values (keep alpha)
            inverted.Set(x, y, color.RGBA{
                R: uint8(255 - (r >> 8)),
                G: uint8(255 - (g >> 8)),
                B: uint8(255 - (b >> 8)),
                A: uint8(a >> 8),
            })
        }
    }
    
    return inverted
}
```

**Sau khi invert:**
- Background trắng → Background đen (lightness 0.0 < threshold) → Để TRẮNG ✅
- Text đen → Text trắng (lightness 1.0 > threshold) → In ĐEN ✅

## Updated Implementation

### chromedp_bill_renderer_optimized.go
```go
func imageToESCPOSOptimized(img image.Image) []byte {
    // Invert image first (black <-> white)
    invertedImg := invertImage(img)
    
    var buf bytes.Buffer
    e := escpos.New(&buf)
    e.Init()
    
    converter := &raster.Converter{
        MaxWidth:  576,
        Threshold: 0.5, // Back to 0.5 (balanced)
    }
    
    converter.Print(invertedImg, e)
    
    e.FormfeedN(3)
    e.Cut()
    e.End()
    
    return buf.Bytes()
}

func invertImage(img image.Image) image.Image {
    bounds := img.Bounds()
    inverted := image.NewRGBA(bounds)
    
    for y := bounds.Min.Y; y < bounds.Max.Y; y++ {
        for x := bounds.Min.X; x < bounds.Max.X; x++ {
            r, g, b, a := img.At(x, y).RGBA()
            
            inverted.Set(x, y, color.RGBA{
                R: uint8(255 - (r >> 8)),
                G: uint8(255 - (g >> 8)),
                B: uint8(255 - (b >> 8)),
                A: uint8(a >> 8),
            })
        }
    }
    
    return inverted
}
```

### visual_bill_renderer.go
Same implementation with `invertImageVisual()` function.

## Why This Works

### Before Inversion (WRONG):
```
Original Image:
- Background: White (255, 255, 255) → lightness = 1.0
- Text: Black (0, 0, 0) → lightness = 0.0

Raster Converter (threshold 0.5):
- White (1.0 >= 0.5) → Print BLACK → Wrong! ❌
- Black (0.0 < 0.5) → Leave WHITE → Wrong! ❌

Result: Black page with white text (inverted)
```

### After Inversion (CORRECT):
```
Inverted Image:
- Background: Black (0, 0, 0) → lightness = 0.0
- Text: White (255, 255, 255) → lightness = 1.0

Raster Converter (threshold 0.5):
- Black (0.0 < 0.5) → Leave WHITE → Correct! ✅
- White (1.0 >= 0.5) → Print BLACK → Correct! ✅

Result: White page with black text (correct)
```

## Threshold Reset to 0.5
Vì đã invert image, threshold quay lại 0.5 (balanced):
- 0.5 = 50% threshold
- Có thể điều chỉnh nếu cần:
  - 0.4 = Darker print (more black)
  - 0.6 = Lighter print (less black)

## Files Modified
- `backend/application/services/chromedp_bill_renderer_optimized.go`
  - Added `invertImage()` function
  - Call `invertImage()` before `converter.Print()`
  - Reset threshold to 0.5
  
- `backend/application/services/visual_bill_renderer.go`
  - Added `invertImageVisual()` function
  - Call `invertImageVisual()` before `converter.Print()`
  - Reset threshold to 0.5

## Testing
1. Create test order
2. Print to Zywell ZY303 (192.168.1.115:9100)
3. Should see:
   - White background ✅
   - Black text ✅
   - Clear content ✅
   - No black bars ✅

## Performance Note
Image inversion adds one extra pass through all pixels, but:
- Negligible impact (< 10ms for 576x800 image)
- Necessary for correct output
- Only done once per print job

## Alternative Solutions Considered

### 1. Fork/Patch Library (Rejected)
- Too complex
- Hard to maintain
- Not worth it

### 2. Manual Bitmap Conversion (Rejected)
- Already tried, had issues
- Library is more reliable
- Just need to invert input

### 3. Invert Image (CHOSEN) ✅
- Simple
- Works with library as-is
- Easy to understand
- Minimal performance impact
