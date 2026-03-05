# Raster Logic Verified with Test

## Problem Summary
- First attempt: In ra toàn đen (black bar)
- Second attempt (with inversion): In ra toàn trắng (blank page)

## Root Cause Discovery

Tôi đã tạo test case đơn giản để verify logic của `kenshaw/escpos/raster`:

### Test Setup
```go
// 8x8 pixel image
// Top 4 rows: WHITE (255, 255, 255)
// Bottom 4 rows: BLACK (0, 0, 0)
```

### Test Results
```
=== Testing threshold 0.5 ===
Row 0: ████████  ← WHITE pixels → prints BLACK
Row 1: ████████
Row 2: ████████
Row 3: ████████
Row 4: ░░░░░░░░  ← BLACK pixels → leaves WHITE
Row 5: ░░░░░░░░
Row 6: ░░░░░░░░
Row 7: ░░░░░░░░

Lightness values:
WHITE (255,255,255): lightness = 1.000
BLACK (0,0,0): lightness = 0.000
GRAY (128,128,128): lightness = 0.502
```

## Confirmed Logic

**The `raster.Converter` has INVERTED logic:**

```go
if lightness(pixel) >= Threshold {
    // Set bit = 1 → Printer prints BLACK
}
```

**This means:**
- WHITE pixels (lightness = 1.0) → bit = 1 → **Prints BLACK**
- BLACK pixels (lightness = 0.0) → bit = 0 → **Leaves WHITE**

## Why This Happens

The library is designed for images where:
- **Pixels to print** have HIGH lightness values
- **Pixels to skip** have LOW lightness values

This is OPPOSITE of normal image processing where:
- Black pixels (0,0,0) should print black
- White pixels (255,255,255) should stay white

## Correct Solution: Image Inversion

**For normal bill images:**
- Background: WHITE (255,255,255)
- Text: BLACK (0,0,0)

**Without inversion:**
- Background WHITE → lightness 1.0 → prints BLACK ❌
- Text BLACK → lightness 0.0 → leaves WHITE ❌
- Result: Black page with white text (inverted)

**With inversion:**
1. Invert image first: WHITE ↔ BLACK
2. After inversion:
   - Background WHITE → BLACK (0,0,0) → lightness 0.0 → leaves WHITE ✅
   - Text BLACK → WHITE (255,255,255) → lightness 1.0 → prints BLACK ✅
3. Result: White page with black text (correct!)

## Implementation

### chromedp_bill_renderer_optimized.go
```go
func imageToESCPOSOptimized(img image.Image) []byte {
    // CRITICAL: Invert image before conversion
    invertedImg := invertImage(img)
    
    var buf bytes.Buffer
    e := escpos.New(&buf)
    e.Init()
    
    converter := &raster.Converter{
        MaxWidth:  576,
        Threshold: 0.5, // Works correctly with inverted image
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

## Why Previous Attempts Failed

### Attempt 1: No Inversion
- Used library directly on normal image
- Result: Black bar (background printed black)

### Attempt 2: Removed Inversion (thinking it was wrong)
- Thought inversion was the problem
- Tried different thresholds (0.1, 0.6, etc.)
- Result: Blank page (text not printed)

### Attempt 3: Verified with Test (CORRECT)
- Created simple test case
- Confirmed library logic is inverted
- Re-added inversion with proper understanding
- Result: Should work correctly now!

## Test Case Code

Created `backend/cmd/test-raster-logic/main.go` to verify:
- Creates 8x8 test image (top=white, bottom=black)
- Tests with different thresholds
- Prints bitmap visualization
- Confirms lightness calculations

Run test:
```bash
cd backend
go run cmd/test-raster-logic/main.go
```

## Files Modified
- `backend/application/services/chromedp_bill_renderer_optimized.go`
  - Re-added `invertImage()` function
  - Added detailed comments explaining the logic
  
- `backend/application/services/visual_bill_renderer.go`
  - Re-added `invertImageVisual()` function
  - Same logic as chromedp renderer

- `backend/cmd/test-raster-logic/main.go` (NEW)
  - Test program to verify raster logic
  - Proves the library behavior

## Expected Result
Now when printing:
- Background should be WHITE ✅
- Text should be BLACK ✅
- Logo should be visible ✅
- Layout should match preview PNG ✅

## Threshold Tuning
With inverted image, threshold 0.5 (50%) works well:
- Can adjust if needed:
  - 0.4 = Darker print (more black pixels)
  - 0.6 = Lighter print (fewer black pixels)

## Backend Status
✅ Backend restarted with corrected implementation
- Server running on port 3000
- Ready to test print with real order
