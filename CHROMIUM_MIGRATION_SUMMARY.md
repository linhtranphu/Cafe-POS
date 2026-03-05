# ✅ Tóm tắt: Chuyển Chromium sang Print Bridge

## Đã hoàn thành

### 1. Phân tích nguyên nhân ✅
- Xác định: Backend EC2 bị OOM do Chromium (948MB image)
- Root cause: HTML template rendering tốn nhiều RAM
- Giải pháp: Chuyển rendering sang local print bridge

### 2. Tạo HTML Renderer cho Print Bridge ✅

**File mới:**
- `local-print-bridge/src/services/htmlRenderer.js`

**Chức năng:**
- Sử dụng Puppeteer (Node.js Chromium automation)
- Render HTML → PNG → Grayscale → ESC/POS
- Singleton pattern để reuse browser instance
- Memory optimizations

**Dependencies đã cài:**
```json
{
  "puppeteer": "^latest",  // Chromium automation
  "sharp": "^latest"        // Image processing
}
```

### 3. Thêm API endpoints ✅

**Endpoints mới:**
```
POST /render-html
- Input: { html, width, shopSettings }
- Output: { escposData (base64), previewImage (base64), stats }

POST /test-render
- Test endpoint với sample HTML
- Verify rendering works
```

### 4. Update Print Bridge ✅

**Changes:**
- Import htmlRenderer service
- Add render endpoints
- Graceful shutdown cleanup

## Bước tiếp theo

### Phase 2: Update Backend (Chưa làm)

**Cần làm:**

1. **Tạo Print Bridge Client** trong backend
   ```go
   // backend/infrastructure/printbridge/client.go
   type PrintBridgeClient struct {
       baseURL string
       timeout time.Duration
   }
   
   func (c *PrintBridgeClient) RenderHTMLToESCPOS(html string) ([]byte, error)
   ```

2. **Update HTML Template Handler**
   ```go
   // backend/interfaces/http/html_template_handler.go
   // OLD: h.chromedpRenderer.RenderBillToESCPOS(...)
   // NEW: h.printBridgeClient.RenderHTMLToESCPOS(...)
   ```

3. **Remove Chromedp initialization**
   ```go
   // backend/main.go
   // Comment out chromedp handler creation
   ```

### Phase 3: Remove Chromium from Backend (Chưa làm)

**Cần làm:**

1. **Update Dockerfile**
   ```dockerfile
   # backend/Dockerfile
   # XÓA các dòng:
   # RUN apk --no-cache add chromium chromium-chromedriver font-noto-cjk
   ```

2. **Remove chromedp code**
   - Xóa hoặc comment out chromedp files
   - Keep for reference nếu cần rollback

3. **Rebuild & test**
   ```bash
   cd backend
   docker build -t cafe-pos-backend:no-chromium .
   docker images | grep cafe-pos-backend
   # Verify size: ~27MB instead of 948MB
   ```

### Phase 4: Deploy & Verify (Chưa làm)

**Deployment steps:**

1. **Deploy new print bridge** (có Chromium)
   ```bash
   cd local-print-bridge
   npm install
   npm start
   # Verify: curl http://localhost:3001/test-render
   ```

2. **Deploy new backend** (không có Chromium)
   ```bash
   # Build & push to Docker Hub
   docker build -t linhtranphu/cafe-pos-backend:no-chromium .
   docker push linhtranphu/cafe-pos-backend:no-chromium
   
   # Deploy to EC2
   ssh ec2-user@your-ec2-ip
   docker pull linhtranphu/cafe-pos-backend:no-chromium
   docker-compose up -d
   ```

3. **Test end-to-end**
   - Test HTML template printing
   - Verify print quality
   - Check memory usage
   - Monitor for 24-48h

## Expected Results

### Before (Current):
```
Backend EC2:
- Image: 948 MB
- Memory idle: 76 MB
- Memory active: 227 MB
- OOM risk: HIGH

Print Bridge:
- No Chromium
- Memory: 26 MB
```

### After (Target):
```
Backend EC2:
- Image: 27 MB (↓ 97%)
- Memory idle: 76 MB
- Memory active: 76 MB (↓ 66%)
- OOM risk: LOW

Print Bridge:
- With Chromium
- Memory idle: 100 MB
- Memory active: 250 MB (when rendering)
- Runs on local machine (more resources)
```

## Testing Checklist

### Print Bridge Testing
- [ ] Puppeteer installs correctly
- [ ] /test-render works
- [ ] /render-html accepts HTML
- [ ] ESC/POS output is valid
- [ ] Image quality is good
- [ ] Vietnamese fonts work
- [ ] Logo rendering works
- [ ] Memory doesn't leak

### Backend Testing
- [ ] Print bridge client works
- [ ] HTML template handler updated
- [ ] Error handling works
- [ ] Fallback to text template if bridge down
- [ ] Backend image size reduced
- [ ] No chromedp dependencies

### Integration Testing
- [ ] Frontend → Backend → Print Bridge → Printer
- [ ] End-to-end HTML template printing
- [ ] Multiple concurrent requests
- [ ] Error scenarios handled

### EC2 Testing
- [ ] Backend deploys successfully
- [ ] Memory usage reduced
- [ ] No OOM errors for 48h
- [ ] Print functionality works
- [ ] Performance acceptable

## Rollback Plan

Nếu có vấn đề:

1. **Keep old backend image** (948MB) as backup tag
2. **Revert to old image:**
   ```bash
   docker pull linhtranphu/cafe-pos-backend:latest-with-chromium
   docker-compose up -d
   ```
3. **Feature flag** trong code để switch rendering method

## Configuration

### Print Bridge .env
```env
PORT=3001
BACKEND_URL=https://tacafe.store
DEFAULT_BILL_PRINTER_IP=192.168.1.115
DEFAULT_LABEL_PRINTER_IP=192.168.1.116
```

### Backend .env (NEW)
```env
PRINT_BRIDGE_URL=http://192.168.1.X:3001
PRINT_BRIDGE_TIMEOUT=30000
PRINT_BRIDGE_ENABLED=true
```

## Timeline Estimate

- ✅ Phase 1 (Print Bridge): 2 hours - DONE
- ⏳ Phase 2 (Backend Update): 2 hours - TODO
- ⏳ Phase 3 (Remove Chromium): 1 hour - TODO
- ⏳ Phase 4 (Deploy & Test): 2 hours - TODO
- **Total remaining: ~5 hours**

## Next Action

**Immediate:** Implement Phase 2 - Update Backend to use Print Bridge

```bash
# 1. Create print bridge client
# 2. Update HTML template handler
# 3. Test integration locally
# 4. Proceed to Phase 3
```

---

**Status:** Phase 1 Complete ✅  
**Next:** Phase 2 - Backend Integration  
**Priority:** HIGH  
**Risk:** LOW (có rollback plan)
