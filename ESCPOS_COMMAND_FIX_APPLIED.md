# ESC/POS Command Fix Applied - COMPLETED ✅

## Problem
Khi in bill, máy in chỉ ra "dãy ký tự vô nghĩa" thay vì in hình ảnh bill.

## Root Cause
ESC/POS library (kenshaw/escpos) sử dụng command **GS 8 L** (Graphics data) để in ảnh. Command này không được support bởi nhiều máy in ESC/POS, đặc biệt là các máy in giá rẻ hoặc cũ.

## Solution
Thay đổi từ **GS 8 L** sang **GS v 0** (Raster bit image) - command được support rộng rãi hơn bởi hầu hết các máy in ESC/POS.

## Implementation

### Before (Using escpos library):
```go
func imageToESCPOSOptimized(img image.Image) []byte {
    invertedImg := invertImage(img)
    
    var buf bytes.Buffer
    e := escpos.New(&buf)
    e.Init()
    
    converter := &raster.Converter{
        MaxWidth:  576,
        Threshold: 0.5,
    }
    
    converter.Print(invertedImg, e)  // Uses GS 8 L command
    
    e.FormfeedN(3)
    e.Cut()
    e.End()
    
    return buf.Bytes()
}
```

**Output**: `1b 40 1d 38 4c ...` (GS 8 L command)

### After (Manual GS v 0 implementation):
```go
func imageToESCPOSOptimized(img image.Image) []byte {
    bounds := img.Bounds()
    width := bounds.Dx()
    height := bounds.Dy()
    
    var buf bytes.Buffer
    
    // ESC @ - Initialize printer
    buf.Write([]byte{0x1B, 0x40})
    
    // Process image line by line
    bytesPerLine := (width + 7) / 8
    
    for y := 0; y < height; y++ {
        // GS v 0 - Print raster bit image
        buf.Write([]byte{0x1D, 0x76, 0x30, 0x00})
        
        // Width in bytes (little endian)
        buf.WriteByte(byte(bytesPerLine & 0xFF))
        buf.WriteByte(byte((bytesPerLine >> 8) & 0xFF))
        
        // Height (1 line at a time)
        buf.WriteByte(0x01)
        buf.WriteByte(0x00)
        
        // Convert line to bitmap
        lineData := make([]byte, bytesPerLine)
        for x := 0; x < width; x++ {
            r, g, b, _ := img.At(x, y).RGBA()
            gray := (r*299 + g*587 + b*114) / 1000
            
            if gray < 32768 {  // Dark pixel
                byteIndex := x / 8
                bitIndex := 7 - (x % 8)
                lineData[byteIndex] |= (1 << bitIndex)
            }
        }
        buf.Write(lineData)
    }
    
    // Feed and cut
    buf.Write([]byte{0x1B, 0x64, 0x03})  // Feed 3 lines
    buf.Write([]byte{0x1D, 0x56, 0x00})  // Cut
    
    return buf.Bytes()
}
```

**Output**: `1b 40 1d 76 30 00 ...` (GS v 0 command)

## Changes Made

### 1. chromedp_bill_renderer_optimized.go
**Modified:**
- Replaced `imageToESCPOSOptimized()` function
- Removed escpos library dependency
- Removed raster library dependency
- Implemented manual GS v 0 raster bit image conversion

**Removed imports:**
```go
"github.com/kenshaw/escpos"
"github.com/kenshaw/escpos/raster"
```

## Technical Details

### GS v 0 Command Format:
```
GS v 0 m xL xH yL yH d1...dk

Where:
- GS = 0x1D
- v = 0x76
- 0 = 0x30 (normal mode)
- m = 0x00 (mode)
- xL xH = width in bytes (little endian)
- yL yH = height in dots (little endian)
- d1...dk = bitmap data
```

### Bitmap Conversion:
- Process image line by line
- Convert each pixel to grayscale using weighted average: `(R*299 + G*587 + B*114) / 1000`
- Threshold at 50%: pixels darker than 50% gray → print black (bit=1)
- Pack 8 pixels into 1 byte (MSB first)

### Data Size Comparison:
- **GS 8 L**: 76,426 bytes (compressed format)
- **GS v 0**: 84,888 bytes (uncompressed, line-by-line)

GS v 0 uses more data but has better compatibility.

## Testing

### Test with uploaded logo:
```bash
cd backend
go run cmd/test-uploaded-logo/main.go
```

**Results:**
- ✅ ESC/POS data generated: 84,888 bytes
- ✅ Format verified: GS v 0 command (1d 76 30 00)
- ✅ Logo rendered correctly in preview
- ✅ Ready for printing

### Verify ESC/POS format:
```bash
hexdump -C test_uploaded_logo.bin | head -10
```

**Expected output:**
```
00000000  1b 40 1d 76 30 00 48 00  01 00 ...
          ^ESC@ ^GS v 0   ^width  ^height
```

## Benefits

✅ **Better compatibility** - GS v 0 supported by most ESC/POS printers
✅ **No external library** - Direct control over ESC/POS commands
✅ **Predictable output** - Line-by-line processing
✅ **Logo rendering works** - With template.URL fix applied

## Testing with Real Printer

1. **Create order** from frontend
2. **Auto print** will trigger
3. **Check printer output** - should print bill with logo correctly

If printer still shows garbage:
- Check printer supports ESC/POS (not just text)
- Verify printer IP and port (usually 9100)
- Check print bridge logs: `docker logs local-print-bridge`

## Files Modified

1. `backend/application/services/chromedp_bill_renderer_optimized.go`
   - Replaced imageToESCPOSOptimized() function
   - Removed escpos/raster imports
   - Implemented manual GS v 0 conversion

## Backend Status

✅ Backend restarted successfully
✅ All services running
✅ Ready for order creation and printing

Access: http://localhost:5173
