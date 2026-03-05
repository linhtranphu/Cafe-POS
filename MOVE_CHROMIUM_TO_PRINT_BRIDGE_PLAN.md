# 📋 Plan: Chuyển Chromium từ Backend EC2 sang Local Print Bridge

## Mục tiêu

Chuyển HTML rendering (Chromium) từ backend EC2 sang local print bridge để:
- ✅ Giảm backend image: 948MB → 27MB (97%)
- ✅ Giảm memory EC2
- ✅ Rendering chạy trên máy local (nhiều resources)
- ✅ Giữ nguyên functionality

## Kiến trúc mới

### Trước (Current):
```
Frontend → Backend EC2 (Chromium) → ESC/POS → Print Bridge → Printer
                ↑
            948MB image
            227MB RAM
```

### Sau (New):
```
Frontend → Backend EC2 → Print Bridge (Chromium) → Printer
                ↑                    ↑
            27MB image          Render HTML
            76MB RAM            Local resources
```

## Thay đổi cần thực hiện

### 1. Backend (EC2) - REMOVE Chromium

**Files cần sửa:**
- `backend/Dockerfile` - Xóa Chromium packages
- `backend/main.go` - Disable chromedp handlers
- `backend/interfaces/http/html_template_handler.go` - Chuyển logic sang API call

**Thay đổi:**
```go
// OLD: Render locally với Chromedp
escposData, err := h.chromedpRenderer.RenderBillToESCPOS(ord, shopSettings)

// NEW: Gửi HTML sang print bridge để render
escposData, err := h.printBridgeClient.RenderHTMLToESCPOS(htmlContent, shopSettings)
```

### 2. Local Print Bridge - ADD Chromium

**Files cần tạo/sửa:**
- `local-print-bridge/package.json` - Add puppeteer
- `local-print-bridge/src/services/htmlRenderer.js` - NEW service
- `local-print-bridge/src/index.js` - Add render endpoint
- `local-print-bridge/Dockerfile` - Add Chromium (optional, for Docker)

**New endpoint:**
```javascript
POST /render-html
{
  "html": "<html>...</html>",
  "width": 576,
  "shopSettings": {...}
}

Response:
{
  "escposData": "base64...",
  "previewImage": "base64..." // optional
}
```

### 3. Frontend - NO CHANGE

Frontend không cần thay đổi gì, vẫn gọi backend API như cũ.

## Implementation Steps

### Phase 1: Add Chromium to Print Bridge ✅

1. Install puppeteer
2. Create HTML renderer service
3. Add /render-html endpoint
4. Test locally

### Phase 2: Update Backend to use Print Bridge ✅

1. Create print bridge client
2. Update HTML template handler
3. Remove chromedp initialization
4. Test integration

### Phase 3: Remove Chromium from Backend ✅

1. Update Dockerfile
2. Remove chromedp code
3. Rebuild image
4. Verify size reduction

### Phase 4: Deploy & Test ✅

1. Deploy new backend to EC2
2. Update print bridge on local machines
3. Test HTML template printing
4. Monitor memory usage

## Technical Details

### Print Bridge HTML Renderer

**Technology:** Puppeteer (Node.js Chromium automation)

**Why Puppeteer vs Chromedp:**
- Node.js native (print bridge is Node.js)
- Easier to install and manage
- Better documentation
- Same Chromium engine

**Process:**
1. Receive HTML + settings
2. Launch headless Chromium
3. Render HTML to image
4. Convert to grayscale
5. Generate ESC/POS commands
6. Return base64 encoded data

### Backend Print Bridge Client

**Responsibilities:**
- Send HTML to print bridge
- Receive ESC/POS data
- Handle errors and retries
- Fallback to text template if bridge unavailable

**Configuration:**
```env
PRINT_BRIDGE_URL=http://192.168.1.X:3001
PRINT_BRIDGE_TIMEOUT=30000
```

## Rollback Plan

Nếu có vấn đề:

1. **Keep old backend image** (948MB) as backup
2. **Feature flag** để switch giữa local render vs bridge render
3. **Fallback** to text template nếu bridge không available

## Testing Checklist

### Local Testing
- [ ] Print bridge can render HTML
- [ ] ESC/POS output is correct
- [ ] Image quality is good
- [ ] Vietnamese fonts work
- [ ] Logo rendering works

### Integration Testing
- [ ] Backend can call print bridge
- [ ] Error handling works
- [ ] Timeout handling works
- [ ] Retry logic works

### EC2 Testing
- [ ] Backend image size reduced
- [ ] Memory usage reduced
- [ ] No OOM errors
- [ ] Print functionality works end-to-end

### Performance Testing
- [ ] Render time acceptable (<5s)
- [ ] Multiple concurrent requests work
- [ ] No memory leaks

## Success Metrics

- ✅ Backend image: 948MB → <50MB
- ✅ Backend memory idle: 227MB → <100MB
- ✅ EC2 stable for 48h without restart
- ✅ HTML template printing works
- ✅ Print quality unchanged

## Timeline

- Phase 1: 2 hours (Add to print bridge)
- Phase 2: 2 hours (Update backend)
- Phase 3: 1 hour (Remove from backend)
- Phase 4: 2 hours (Deploy & test)
- **Total: ~7 hours**

## Next Steps

1. ✅ Create HTML renderer service for print bridge
2. ✅ Add puppeteer to print bridge
3. ✅ Test rendering locally
4. ✅ Update backend to use print bridge
5. ✅ Deploy and verify

---

**Status:** Ready to implement  
**Priority:** HIGH  
**Risk:** LOW (có rollback plan)
