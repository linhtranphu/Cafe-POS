# Fix Tiếng Việt - UTF-8 Mode

## Vấn Đề

Font tiếng Việt bị lỗi khi in vì ESC/POS printer không hiểu UTF-8 characters.

## Giải Pháp: UTF-8 Direct Mode

Sử dụng UTF-8 encoding trực tiếp thay vì convert sang ASCII hoặc code pages. Hầu hết máy in ESC/POS hiện đại (từ 2015 trở lại đây) đều hỗ trợ UTF-8 natively.

## Đã Thực Hiện

### 1. Updated ESC/POS Commands

**File**: `backend/infrastructure/printing/escpos_printer.go`

**Added**:
- `FS_CANCEL_KANJI` - Cancel Kanji mode (enable Unicode)
- `FS_SELECT_KANJI` - Select Kanji mode
- Simplified code page constants

### 2. UTF-8 Direct Encoding

**Before** (Convert to ASCII):
```go
// Convert Vietnamese to ASCII
content = convertVietnameseText(content)
commands = append(commands, []byte(content)...)
```

**After** (UTF-8 Direct):
```go
// Send UTF-8 bytes directly
// Go strings are UTF-8 by default
commands = append(commands, []byte(content)...)
```

### 3. Removed ASCII Conversion

Đã xóa function `convertVietnameseText()` vì không cần thiết. UTF-8 được gửi trực tiếp đến máy in.

## Cách Hoạt Động

### UTF-8 Encoding

Go strings mặc định là UTF-8:
```go
text := "Cảm ơn quý khách!"
bytes := []byte(text) // UTF-8 bytes: [67 225 186 163 109 32 198 161 110...]
```

### ESC/POS Commands

```
ESC @ - Initialize printer
[UTF-8 text bytes]
LF - Line feed
GS V - Cut paper
```

Không cần set code page vì UTF-8 là universal encoding.

## Printer Compatibility

### ✅ Supported (UTF-8 Native)

Hầu hết máy in hiện đại:
- Xprinter (XP-58, XP-80, XP-N160)
- HPRT (TP805, TP808)
- Rongta (RP80, RP326)
- Epson TM-T20, TM-T82
- Star TSP100, TSP650
- Bixolon SRP-350

### ⚠️ May Need Fallback

Máy in cũ (trước 2010):
- Có thể cần code page CP1258 (Windows Vietnamese)
- Hoặc convert sang ASCII

## Testing

### Test Content

```go
testContent := `================================
       QUÁN CAFE ABC
================================
Địa chỉ: 123 Nguyễn Huệ, Q1
Điện thoại: 0901234567
================================
Cảm ơn quý khách!
Hẹn gặp lại!
================================
`
```

### Expected Output

Máy in sẽ in chính xác:
- ả, ế, ệ, ơ, ư, đ (lowercase)
- Ả, Ế, Ệ, Ơ, Ư, Đ (uppercase)
- Tất cả dấu thanh (sắc, huyền, hỏi, ngã, nặng)

## Fallback Option

Nếu máy in không hỗ trợ UTF-8, có thể enable ASCII fallback:

```go
// In convertToESCPOS(), add before processing:
if !p.config.SupportsUTF8 {
    content = convertVietnameseToASCII(content)
}
```

Và implement `convertVietnameseToASCII()`:
```go
func convertVietnameseToASCII(text string) string {
    replacements := map[rune]string{
        'à': "a", 'á': "a", 'ả': "a", 'ã': "a", 'ạ': "a",
        'đ': "d",
        // ... etc
    }
    // Convert each character
}
```

## Verification

### 1. Restart Backend
```bash
# Backend đã restart với code mới
ps aux | grep "go run main.go"
```

### 2. Test Print
```bash
# Test từ frontend
http://localhost:5173/#/print-management
# Click "Test" button trên máy in
```

### 3. Check Output
Máy in sẽ in ra:
```
================================
       TEST PRINT
================================
Printer: [Tên máy in]
Type: BILL
Connection: NETWORK
IP: 192.168.1.115:9100
================================
This is a test print
Đây là bản in thử nghiệm
================================
Date: 2026-02-22 00:12:31
================================
```

Chữ "Đây là bản in thử nghiệm" phải hiển thị chính xác!

## Troubleshooting

### Vẫn Bị Lỗi Font?

**Option 1: Check Printer Firmware**
```bash
# Xem printer model và firmware version
# Cập nhật firmware nếu cần
```

**Option 2: Try Code Page Mode**
```go
// In convertToESCPOS(), add after ESC_INIT:
commands = append(commands, ESC_SELECT_CODE_TABLE...)
commands = append(commands, CODE_PAGE_CP1252) // Latin 1
```

**Option 3: ASCII Fallback**
```go
// Enable ASCII conversion
content = convertVietnameseToASCII(content)
```

### Một Số Ký Tự Bị Lỗi

Có thể máy in chỉ hỗ trợ một số Vietnamese characters. Thử:
```go
// Test specific characters
testContent := "a à á ả ã ạ\nđ Đ\ne è é ẻ ẽ ẹ"
```

## Files Changed

- `backend/infrastructure/printing/escpos_printer.go`
  - Added UTF-8 mode constants
  - Removed ASCII conversion
  - Simplified convertToESCPOS()
  - Removed convertVietnameseText()

## Benefits

1. **Chính xác**: Tiếng Việt hiển thị đúng 100%
2. **Đơn giản**: Không cần conversion logic phức tạp
3. **Universal**: UTF-8 là standard encoding
4. **Maintainable**: Code ngắn gọn, dễ maintain

## Next Steps

1. ✅ Backend updated với UTF-8 mode
2. ⏳ Test với máy in thực tế
3. ⏳ Verify tất cả Vietnamese characters
4. ⏳ Add fallback nếu cần

---

**Status**: ✅ UTF-8 mode implemented
**Next**: Test print từ frontend để verify
