# Hướng Dẫn Background Workers

## Tổng Quan

Background workers (công nhân nền) đã được **tự động kích hoạt** khi server khởi động. Không cần cấu hình thêm!

## Cách Hoạt Động

### 1. Khởi Động Tự Động

Khi bạn chạy server backend, background workers sẽ tự động start:

```bash
cd backend
go run main.go
```

Bạn sẽ thấy log:
```
Starting cost recalculation worker pool...
✅ Cost recalculation worker pool started
```

### 2. Cấu Hình Worker Pool

Trong file `main.go`, worker pool được cấu hình với:

```go
costRecalculationService := services.NewCostRecalculationService(
    costCalculatorService, 
    menuRepo, 
    4,      // 4 workers (số lượng goroutines xử lý đồng thời)
    1000    // Queue size 1000 (số lượng jobs tối đa trong hàng đợi)
)
```

**Giải thích**:
- **4 workers**: Có 4 goroutines chạy song song để xử lý cost recalculation
- **Queue size 1000**: Có thể chứa tối đa 1000 jobs đang chờ xử lý

### 3. Khi Nào Workers Hoạt Động?

Workers tự động xử lý khi:

#### A. Cập Nhật Giá Nguyên Liệu
Khi manager cập nhật `cost_per_unit` của một ingredient:

```bash
# Ví dụ: Cập nhật giá cà phê
PATCH /api/manager/ingredients/:id
{
  "cost_per_unit": 250000
}
```

**Điều gì xảy ra**:
1. Ingredient được cập nhật trong database
2. Hệ thống tìm tất cả menu items sử dụng ingredient này
3. Mỗi menu item được queue vào background job
4. Workers xử lý từng job để tính lại `current_cost`

#### B. Xem Trạng Thái Workers

Kiểm tra workers có đang hoạt động:

```bash
curl -H "Authorization: Bearer $TOKEN" \
  http://localhost:3000/api/manager/monitoring/metrics/aggregated
```

Response:
```json
{
  "total_recalc_jobs": 150,
  "successful_recalc_jobs": 148,
  "failed_recalc_jobs": 2,
  "average_recalc_duration": 200000000
}
```

### 4. Graceful Shutdown

Khi tắt server (Ctrl+C), workers sẽ tự động dừng an toàn:

```
Stopping cost recalculation worker pool...
✅ Cost recalculation worker pool stopped
```

Tất cả jobs đang xử lý sẽ hoàn thành trước khi server tắt.

## Kiểm Tra Workers Hoạt Động

### Test 1: Cập Nhật Ingredient và Xem Queue

```bash
# 1. Lấy token
TOKEN=$(curl -X POST http://localhost:3000/api/login \
  -H "Content-Type: application/json" \
  -d '{"username":"admin","password":"admin123"}' \
  | jq -r '.token')

# 2. Cập nhật giá ingredient
curl -X PATCH http://localhost:3000/api/manager/ingredients/INGREDIENT_ID \
  -H "Authorization: Bearer $TOKEN" \
  -H "Content-Type: application/json" \
  -d '{"cost_per_unit": 300000}'

# 3. Kiểm tra metrics ngay lập tức
curl -H "Authorization: Bearer $TOKEN" \
  http://localhost:3000/api/manager/monitoring/metrics?type=recalculation_job&limit=10
```

### Test 2: Xem Health Status

```bash
curl -H "Authorization: Bearer $TOKEN" \
  http://localhost:3000/api/manager/monitoring/health
```

Response:
```json
{
  "status": "healthy",
  "error_rates": {
    "cost_calculation": 0.5,
    "recalculation_job": 0.2,
    "shift_closure": 0.0
  },
  "recent_alerts": {
    "critical": 0,
    "warning": 0,
    "total": 0
  }
}
```

## Điều Chỉnh Cấu Hình

### Tăng Số Lượng Workers

Nếu hệ thống có nhiều menu items và cần xử lý nhanh hơn:

```go
// Trong main.go, thay đổi từ 4 thành 8 workers
costRecalculationService := services.NewCostRecalculationService(
    costCalculatorService, 
    menuRepo, 
    8,      // Tăng lên 8 workers
    1000
)
```

**Khi nào cần tăng**:
- Có >500 menu items
- Cập nhật ingredient thường xuyên
- Thấy queue đầy (error "queue is full")

### Tăng Queue Size

Nếu thấy lỗi "recalculation queue is full":

```go
costRecalculationService := services.NewCostRecalculationService(
    costCalculatorService, 
    menuRepo, 
    4,
    2000    // Tăng queue size lên 2000
)
```

## Monitoring Workers

### 1. Xem Metrics Chi Tiết

```bash
# Xem tất cả recalculation jobs
curl -H "Authorization: Bearer $TOKEN" \
  "http://localhost:3000/api/manager/monitoring/metrics?type=recalculation_job&limit=50"
```

### 2. Xem Alerts

```bash
# Xem critical alerts
curl -H "Authorization: Bearer $TOKEN" \
  "http://localhost:3000/api/manager/monitoring/alerts?level=critical&limit=20"
```

### 3. Xem Aggregated Stats

```bash
curl -H "Authorization: Bearer $TOKEN" \
  "http://localhost:3000/api/manager/monitoring/metrics/aggregated"
```

## Troubleshooting

### Vấn Đề 1: Workers Không Xử Lý Jobs

**Triệu chứng**: Cập nhật ingredient nhưng menu cost không thay đổi

**Kiểm tra**:
```bash
# Xem log server
tail -f backend/server.log

# Kiểm tra metrics
curl -H "Authorization: Bearer $TOKEN" \
  http://localhost:3000/api/manager/monitoring/metrics/aggregated
```

**Giải pháp**:
- Restart server
- Kiểm tra MongoDB connection
- Xem alerts để biết lỗi cụ thể

### Vấn Đề 2: Queue Full

**Triệu chứng**: Error "recalculation queue is full"

**Giải pháp**:
1. Tăng queue size trong `main.go`
2. Tăng số workers
3. Restart server

### Vấn Đề 3: High Error Rate

**Triệu chứng**: Alert "High Recalculation Job Failure Rate"

**Kiểm tra**:
```bash
# Xem failed jobs
curl -H "Authorization: Bearer $TOKEN" \
  "http://localhost:3000/api/manager/monitoring/metrics?type=recalculation_job&limit=100" \
  | jq '.metrics[] | select(.status == "failure")'
```

**Giải pháp**:
- Kiểm tra ingredient data quality
- Verify menu item recipes
- Check database performance

## Performance Tips

### 1. Optimal Worker Count

```
Số workers = Số CPU cores * 2
```

Ví dụ: Server có 4 cores → 8 workers

### 2. Queue Size

```
Queue size = Số menu items * 2
```

Ví dụ: 500 menu items → queue size 1000

### 3. Monitoring Frequency

Kiểm tra health status:
- **Hàng ngày**: Xem aggregated metrics
- **Hàng tuần**: Review error trends
- **Khi có alert**: Investigate ngay lập tức

## Tóm Tắt

✅ **Background workers tự động start** khi server khởi động
✅ **Không cần cấu hình thêm** - hoạt động ngay out of the box
✅ **Tự động xử lý** khi cập nhật ingredient costs
✅ **Graceful shutdown** khi tắt server
✅ **Full monitoring** với metrics và alerts

**Để verify workers đang chạy**:
```bash
# Xem log khi start server
go run main.go

# Hoặc check metrics
curl -H "Authorization: Bearer $TOKEN" \
  http://localhost:3000/api/manager/monitoring/health
```

Nếu thấy `"status": "healthy"` → Workers đang hoạt động tốt! 🎉
