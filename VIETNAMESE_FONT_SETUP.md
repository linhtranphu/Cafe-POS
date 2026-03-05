# Vietnamese Font Setup for EC2/Linux

## Problem
Khi deploy trên EC2 Linux, chromedp có thể không render đúng tiếng Việt vì thiếu font hỗ trợ Unicode Vietnamese.

## Solution

### For Ubuntu/Debian (EC2):
```bash
sudo apt update
sudo apt install -y fonts-noto fonts-noto-cjk fonts-noto-color-emoji
```

### For Amazon Linux 2:
```bash
sudo yum install -y google-noto-sans-fonts google-noto-serif-fonts
```

### For Alpine Linux (Docker):
```dockerfile
RUN apk add --no-cache \
    chromium \
    nss \
    freetype \
    harfbuzz \
    ca-certificates \
    ttf-freefont \
    font-noto
```

## Verify Font Installation

```bash
# List installed fonts
fc-list | grep -i noto

# Check Vietnamese support
fc-list :lang=vi
```

## Chromedp Configuration Applied

### 1. Window Size for 80mm Bill
```go
chromedp.WindowSize(576, 2000)  // 576px = 80mm at 72 DPI
```

### 2. Viewport Emulation
```go
chromedp.EmulateViewport(576, 2000)  // Exact bill width
```

### 3. Font Rendering Flags
```go
chromedp.Flag("font-render-hinting", "none")
chromedp.Flag("disable-font-subpixel-positioning", false)
```

### 4. Font Loading Wait
```javascript
document.fonts.ready.then(() => {
    console.log('Fonts loaded');
});
```

## Testing

### Test Vietnamese Characters:
```bash
cd backend
go run cmd/test-uploaded-logo/main.go
```

Check preview image for:
- ✅ Dấu sắc, huyền, hỏi, ngã, nặng
- ✅ Chữ đ, ă, â, ê, ô, ơ, ư
- ✅ No boxes or question marks

## HTML Template Font Stack

Template uses comprehensive font stack with Vietnamese support:
```css
font-family: -apple-system, BlinkMacSystemFont, "Segoe UI", Arial, "Noto Sans", sans-serif;
```

**Font priority:**
1. **-apple-system** (macOS) - Excellent Vietnamese support
2. **BlinkMacSystemFont** (Chrome on macOS)
3. **Segoe UI** (Windows) - Good Vietnamese support
4. **Arial** - Fallback for older systems
5. **Noto Sans** (Linux) - Best Vietnamese support on Linux
6. **sans-serif** - Generic fallback

**Why this works:**
- ✅ macOS: Uses San Francisco font (perfect Vietnamese)
- ✅ Windows: Uses Segoe UI (good Vietnamese)
- ✅ Linux: Falls back to Noto Sans (if installed)
- ✅ No custom font loading needed
- ✅ Fast rendering

## Common Issues

### Issue 1: Boxes instead of Vietnamese characters
**Cause**: No Vietnamese font installed
**Fix**: Install fonts-noto

### Issue 2: Wrong diacritics position
**Cause**: Font doesn't support Vietnamese combining marks
**Fix**: Use Noto Sans or Roboto

### Issue 3: Bill width too wide/narrow
**Cause**: Viewport not set correctly
**Fix**: Use EmulateViewport(576, 2000)

## Production Deployment Checklist

- [ ] Install Vietnamese fonts on server
- [ ] Verify font installation with `fc-list :lang=vi`
- [ ] Test bill printing with Vietnamese text
- [ ] Check preview images for correct rendering
- [ ] Monitor chromedp logs for font warnings

## Docker Deployment

If using Docker, add to Dockerfile:
```dockerfile
FROM golang:1.21-alpine

# Install chromium and fonts
RUN apk add --no-cache \
    chromium \
    nss \
    freetype \
    harfbuzz \
    font-noto \
    font-noto-cjk

# Set chromium path
ENV CHROME_BIN=/usr/bin/chromium-browser
ENV CHROME_PATH=/usr/bin/chromium-browser

# ... rest of Dockerfile
```

## Verification Script

```bash
#!/bin/bash
echo "Checking Vietnamese font support..."

# Check if fonts are installed
if fc-list :lang=vi | grep -q "Noto"; then
    echo "✅ Vietnamese fonts installed"
else
    echo "❌ Vietnamese fonts missing"
    echo "Run: sudo apt install fonts-noto"
fi

# Check chromium
if command -v chromium &> /dev/null; then
    echo "✅ Chromium installed"
else
    echo "❌ Chromium missing"
fi
```

Save as `check-fonts.sh` and run:
```bash
chmod +x check-fonts.sh
./check-fonts.sh
```
