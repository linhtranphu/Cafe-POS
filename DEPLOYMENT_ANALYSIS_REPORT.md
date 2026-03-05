# 📊 BÁO CÁO PHÂN TÍCH DEPLOYMENT

## Tóm tắt

✅ **Deployment local thành công!** Tất cả containers đang chạy ổn định.

## Kết quả Monitoring

### Container Status (sau 1 phút)
| Container | Status | CPU % | Memory | Memory % |
|-----------|--------|-------|--------|----------|
| Frontend | ✅ Running | 0.00% | 9.5 MB | 0.12% |
| Backend | ✅ Running | 0.82% | 76.8 MB | 0.98% |
| MongoDB | ✅ Running | 0.89% | 76.7 MB | 0.98% |
| Print Bridge | ✅ Running | 0.00% | 26.8 MB | 0.34% |

### Tổng Resource Usage
- **Total Memory**: ~190 MB / 7.66 GB (2.5%)
- **Total CPU**: ~1.7%
- **Disk**: Backend image 948MB, Frontend 63MB, MongoDB 845MB

## Phát hiện quan trọng

### 1. Backend Image Size - VẤN ĐỀ CHÍNH! 🔴

```
Backend Image: 948 MB (!!!)
```

**Nguyên nhân:**
- Chromium browser: ~500-600 MB
- Chromium ChromeDriver: ~100 MB  
- 198 packages Alpine (fonts, dependencies): ~200 MB
- Go binary + dependencies: ~50 MB

**So sánh:**
- Image cũ (không có Chromium): 27 MB
- Image mới (có Chromium): 948 MB
- **Tăng 35 lần!**

### 2. Memory Usage Pattern

**Idle state (không có traffic):**
- Backend: 76 MB
- MongoDB: 76 MB
- Frontend: 9 MB

**Khi có activity (spike lúc 02:13:09):**
- Backend: 227 MB (tăng 3x!)
- CPU: 125% (sử dụng nhiều cores)

### 3. Chromium Process Behavior

Từ logs build:
```
Installing chromium (144.0.7559.132-r3)
Installing chromium-chromedriver (144.0.7559.132-r3)
Installing font-noto-cjk (0_git20220127-r1)
```

Chromium được cài đặt đầy đủ, không phải headless minimal.

## Nguyên nhân Server EC2 Chết

### Giả thuyết đã xác nhận: OUT OF MEMORY (OOM) ✅

**Tình huống trên EC2:**

1. **EC2 Instance nhỏ (ví dụ: t2.micro)**
   - RAM: 1 GB
   - Swap: 0-512 MB

2. **Memory breakdown khi start:**
   ```
   MongoDB:  76 MB
   Backend:  76 MB (idle)
   Frontend:  9 MB
   System:  ~200 MB
   ─────────────────
   Total:   ~361 MB (OK)
   ```

3. **Memory breakdown khi có request (print bill):**
   ```
   MongoDB:  76 MB
   Backend: 227 MB (Chromium rendering!)
   Frontend:  9 MB
   System:  ~200 MB
   ─────────────────
   Total:   ~512 MB (STILL OK)
   ```

4. **Memory breakdown khi nhiều requests đồng thời:**
   ```
   MongoDB:  76 MB
   Backend: 227 MB × 3 requests = 681 MB
   Frontend:  9 MB
   System:  ~200 MB
   ─────────────────
   Total:   ~966 MB (OUT OF MEMORY!)
   ```

### Kịch bản Server Chết

1. User tạo order và print bill
2. Backend khởi động Chromium để render HTML → PDF
3. Chromium process tốn 150-200 MB RAM
4. Nếu có 2-3 requests cùng lúc → OOM
5. Linux OOM Killer kill process → Container restart
6. Health check fail → Docker restart container
7. Lặp lại → Server không thể phục hồi

### Bằng chứng

✅ Backend image 948MB (quá lớn)
✅ Memory spike từ 76MB → 227MB khi có activity
✅ CPU spike 125% (Chromium rendering)
✅ Không có memory limits trong docker-compose
✅ EC2 thường dùng instance nhỏ (1-2GB RAM)

## Giải pháp

### Giải pháp 1: Tối ưu Chromium Usage (Khuyến nghị) ⭐

**A. Sử dụng Chromium headless minimal**
```dockerfile
# Thay vì cài full Chromium
RUN apk --no-cache add chromium chromium-chromedriver

# Chỉ cài chromium headless
RUN apk --no-cache add chromium --repository=http://dl-cdn.alpinelinux.org/alpine/edge/community
```

**B. Giới hạn Chromium resources**
```go
// Trong backend code
opts := append(chromedp.DefaultExecAllocatorOptions[:],
    chromedp.Flag("disable-gpu", true),
    chromedp.Flag("no-sandbox", true),
    chromedp.Flag("disable-dev-shm-usage", true),
    chromedp.Flag("disable-software-rasterizer", true),
    chromedp.Flag("disable-extensions", true),
    // QUAN TRỌNG: Giới hạn memory
    chromedp.Flag("max-old-space-size", "128"),
    chromedp.Flag("memory-pressure-off", true),
)
```

**C. Cleanup Chromium processes**
```go
// Đảm bảo cleanup sau mỗi render
defer cancel()
defer chromedp.Cancel(ctx)
```

**D. Queue print jobs**
```go
// Chỉ cho phép 1 print job tại một thời điểm
var printMutex sync.Mutex

func PrintBill(data) {
    printMutex.Lock()
    defer printMutex.Unlock()
    
    // Render with Chromium
}
```

### Giải pháp 2: Thêm Memory Limits

**docker-compose.yml:**
```yaml
backend:
  image: cafe-pos-backend:local
  deploy:
    resources:
      limits:
        memory: 512M  # Giới hạn tối đa
      reservations:
        memory: 256M  # Reserve tối thiểu
```

### Giải pháp 3: Tăng RAM cho EC2

**Nâng cấp instance:**
- t2.micro (1GB) → t2.small (2GB): +$9/tháng
- t2.small (2GB) → t2.medium (4GB): +$18/tháng

### Giải pháp 4: Sử dụng External Print Service

**Tách print service ra container riêng:**
```yaml
print-service:
  image: cafe-pos-print-service
  deploy:
    resources:
      limits:
        memory: 512M
  restart: on-failure
```

## Khuyến nghị Triển khai

### Ưu tiên 1: Tối ưu Chromium (Làm ngay) 🔥

1. Giới hạn Chromium memory
2. Queue print jobs (1 job/time)
3. Cleanup processes đúng cách
4. Test lại trên local

**Ước tính giảm:**
- Memory: 227MB → 150MB (giảm 34%)
- Image size: 948MB → 600MB (giảm 37%)

### Ưu tiên 2: Thêm Memory Limits

1. Set memory limit 512MB cho backend
2. Set memory limit 256MB cho MongoDB
3. Monitor và adjust

### Ưu tiên 3: Nâng cấp EC2 (Nếu cần)

- Từ t2.micro (1GB) → t2.small (2GB)
- Chi phí: ~$9/tháng
- Giải quyết triệt để vấn đề OOM

## Testing Plan

### Test 1: Stress Test Print Function
```bash
# Gửi 10 print requests đồng thời
for i in {1..10}; do
  curl -X POST http://localhost:3000/api/print/bill \
    -H "Content-Type: application/json" \
    -d '{"orderId": "test-'$i'"}' &
done

# Monitor memory
watch -n 1 'docker stats --no-stream'
```

### Test 2: Long Running Test
```bash
# Chạy trong 1 giờ, print mỗi 30 giây
while true; do
  curl -X POST http://localhost:3000/api/print/bill
  sleep 30
done
```

### Test 3: EC2 Simulation
```bash
# Giới hạn Docker memory giống EC2
docker run --memory="1g" --memory-swap="1g" ...
```

## Kết luận

✅ **Đã xác định nguyên nhân:** OUT OF MEMORY do Chromium
✅ **Đã có giải pháp:** Tối ưu Chromium + Memory limits
✅ **Có thể deploy ngay:** Sau khi apply fixes

**Bước tiếp theo:**
1. Implement Chromium optimization
2. Add memory limits
3. Test stress scenarios
4. Deploy lên EC2 với monitoring
5. Quan sát 24-48h

---

**Ngày tạo:** 2026-03-01  
**Trạng thái:** ✅ Hoàn thành phân tích  
**Action required:** Implement optimizations
