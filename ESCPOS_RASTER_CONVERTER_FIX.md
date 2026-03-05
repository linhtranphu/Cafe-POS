# ESC/POS Raster Converter Fix

## Problem
Mẫu in ra không giống hình render mặc dù đã sử dụng `kenshaw/escpos` library. Implementation thủ công của bitmap conversion chưa đúng.

## Root Cause
Code trước đó tự implement bitmap conversion:
```go
// Manual bitmap conversion (WRONG)
for y := 0; y < height; y++ {
    for x := 0; x < width; x++ {
        r, g, b, _ := img.At(x, y).RGBA()
        gray := (r + g + b) / 3  // Simple average
        
        if gray < 32768 {  // 50% threshold
            byteIndex := y*bytesPerLine + x/8
            bitIndex := 7 - (x % 8)
            imgBW[byteIndex] |= (1 << bitIndex)
        }
    }
}
```

**Vấn đề:**
1. Cách tính grayscale đơn giản (average RGB) không chính xác
2. Threshold cứng (32768) không tối ưu
3. Bit packing có thể sai
4. Không sử dụng proper lightness calculation

## Solution: Use `raster.Converter`

Library `kenshaw/escpos` có package `raster` với `Converter` đã implement đúng cách:

### Lightness Calculation (Proper)
```go
const (
    lumR, lumG, lumB = 55, 182, 18  // Weighted for human perception
)

func lightness(c color.Color) float64 {
    r, g, b, _ := c.RGBA()
    return float64(lumR*r+lumG*g+lumB*b) / float64(0xffff*(lumR+lumG+lumB))
}
```

### Bitmap Conversion (Proper)
```go
func (c *Converter) ToRaster(img image.Image) (data []byte, imageWidth, bytesWidth int) {
    sz := img.Bounds().Size()
    
    imageWidth = sz.X
    if imageWidth > c.MaxWidth {
        imageWidth = c.MaxWidth
    }
    
    bytesWidth = imageWidth / 8
    if imageWidth%8 != 0 {
        bytesWidth += 1
    }
    
    data = make([]byte, bytesWidth*sz.Y)
    
    for y := 0; y < sz.Y; y++ {
        for x := 0; x < imageWidth; x++ {
            if lightness(img.At(x, y)) >= c.Threshold {
                data[y*bytesWidth+x/8] |= 0x80 >> uint(x%8)
            }
        }
    }
    
    return
}
```

## Updated Implementation

### chromedp_bill_renderer_optimized.go
```go
import (
    // ... other imports
    "github.com/kenshaw/escpos"
    "github.com/kenshaw/escpos/raster"
)

func imageToESCPOSOptimized(img image.Image) []byte {
    var buf bytes.Buffer
    
    // Create escpos writer
    e := escpos.New(&buf)
    
    // Initialize printer
    e.Init()
    
    // Use raster converter from escpos library
    converter := &raster.Converter{
        MaxWidth:  576, // K80 printer width
        Threshold: 0.5, // 50% threshold (0.0 = all black, 1.0 = all white)
    }
    
    // Convert and print image
    converter.Print(img, e)
    
    // Feed paper
    e.FormfeedN(3)
    
    // Cut paper
    e.Cut()
    
    // End
    e.End()
    
    return buf.Bytes()
}
```

### visual_bill_renderer.go
```go
import (
    // ... other imports
    "github.com/kenshaw/escpos"
    "github.com/kenshaw/escpos/raster"
)

func (r *VisualBillRenderer) imageToESCPOS(img image.Image) ([]byte, error) {
    var buf bytes.Buffer
    
    // Create escpos writer
    e := escpos.New(&buf)
    
    // Initialize printer
    e.Init()
    
    // Use raster converter from escpos library
    converter := &raster.Converter{
        MaxWidth:  576, // K80 printer width
        Threshold: 0.5, // 50% threshold (0.0 = all black, 1.0 = all white)
    }
    
    // Convert and print image
    converter.Print(img, e)
    
    // Feed paper
    e.FormfeedN(3)
    
    // Cut paper
    e.Cut()
    
    // End
    e.End()
    
    return buf.Bytes(), nil
}
```

## Key Improvements

### 1. Proper Lightness Calculation
- Uses weighted RGB values (55:182:18) based on human perception
- Green contributes most to perceived brightness
- More accurate than simple average

### 2. Correct Bit Packing
- Uses `0x80 >> uint(x%8)` instead of `1 << bitIndex`
- Matches ESC/POS raster format specification
- MSB first, left to right

### 3. Configurable Threshold
- `Threshold: 0.5` = 50% (default)
- Can adjust: 0.0 (all black) to 1.0 (all white)
- Lower value = darker print
- Higher value = lighter print

### 4. Automatic Width Handling
- Truncates if image > MaxWidth
- Proper bytesWidth calculation
- Handles non-multiple-of-8 widths

## Threshold Tuning

If print quality is not good, adjust threshold:

```go
converter := &raster.Converter{
    MaxWidth:  576,
    Threshold: 0.4,  // Darker (more black pixels)
    // or
    Threshold: 0.6,  // Lighter (fewer black pixels)
}
```

**Recommended values:**
- 0.3-0.4: For light/faded prints (need more black)
- 0.5: Default (balanced)
- 0.6-0.7: For dark/heavy prints (need less black)

## Files Modified
- `backend/application/services/chromedp_bill_renderer_optimized.go`
  - Replaced manual bitmap conversion with `raster.Converter`
  - Added import for `github.com/kenshaw/escpos/raster`
  
- `backend/application/services/visual_bill_renderer.go`
  - Replaced manual bitmap conversion with `raster.Converter`
  - Added import for `github.com/kenshaw/escpos/raster`

## Testing
1. Create a test order
2. Print bill to Zywell ZY303 (192.168.1.115:9100)
3. Compare print output with preview PNG
4. Adjust threshold if needed

## Expected Results
- Print output should match preview PNG closely
- Text should be clear and readable
- Logo should be recognizable
- Lines and borders should be sharp
- No distortion or artifacts

## Troubleshooting

### Print too light
```go
Threshold: 0.4  // Lower threshold = more black pixels
```

### Print too dark
```go
Threshold: 0.6  // Higher threshold = fewer black pixels
```

### Print distorted
- Check image width = 576px
- Verify printer DPI = 203
- Check paper width = 80mm (K80)

## References
- Library: https://github.com/kenshaw/escpos
- Raster package: https://github.com/kenshaw/escpos/blob/master/raster/raster.go
- ESC/POS spec: GS 8 L (Store raster format graphics data)
