# Background Workers - Quick Start

## TL;DR

✅ **Background workers ĐÃ ĐƯỢC ENABLE tự động!**

Không cần làm gì thêm. Workers tự động start khi server khởi động.

## Verify Workers Hoạt Động

### Cách 1: Xem Log Server

```bash
cd backend
go run main.go
```

Bạn sẽ thấy:
```
Starting cost recalculation worker pool...
✅ Cost recalculation worker pool started
```

### Cách 2: Chạy Test Script

```bash
./test-background-workers.sh
```

Script sẽ kiểm tra:
- ✅ Health status
- ✅ Worker metrics
- ✅ Recent jobs
- ✅ Alerts

### Cách 3: Check API

```bash
# Login
TOKEN=$(curl -X POST http://localhost:3000/api/login \
  -H "Content-Type: application/json" \
  -d '{"username":"admin","password":"admin123"}' \
  | jq -r '.token')

# Check health
curl -H "Authorization: Bearer $TOKEN" \
  http://localhost:3000/api/manager/monitoring/health
```

Response:
```json
{
  "status": "healthy",  // ← Workers đang hoạt động tốt!
  "error_rates": {
    "recalculation_job": 0.2
  }
}
```

## Khi Nào Workers Hoạt Động?

Workers tự động xử lý khi:

1. **Cập nhật giá nguyên liệu**:
   ```bash
   PATCH /api/manager/ingredients/:id
   {
     "cost_per_unit": 250000
   }
   ```
   → Workers tự động tính lại cost cho tất cả menu items sử dụng ingredient này

2. **Xem kết quả**:
   ```bash
   GET /api/manager/menu/costs
   ```
   → `current_cost` đã được cập nhật!

## Cấu Hình

Trong `backend/main.go`:

```go
costRecalculationService := services.NewCostRecalculationService(
    costCalculatorService, 
    menuRepo, 
    4,      // 4 workers (có thể tăng nếu cần)
    1000    // Queue size 1000
)
```

**Khi nào cần tăng workers**:
- Có >500 menu items
- Cập nhật ingredient thường xuyên
- Thấy queue full errors

## Monitoring

### Dashboard URLs

- **Health**: `GET /api/manager/monitoring/health`
- **Metrics**: `GET /api/manager/monitoring/metrics/aggregated`
- **Alerts**: `GET /api/manager/monitoring/alerts`

### Key Metrics

```json
{
  "total_recalc_jobs": 150,
  "successful_recalc_jobs": 148,
  "failed_recalc_jobs": 2,
  "average_recalc_duration": 200000000  // 200ms
}
```

## Troubleshooting

### Workers không xử lý jobs?

1. **Check server log**:
   ```bash
   tail -f backend/server.log
   ```

2. **Check metrics**:
   ```bash
   curl -H "Authorization: Bearer $TOKEN" \
     http://localhost:3000/api/manager/monitoring/metrics/aggregated
   ```

3. **Restart server**:
   ```bash
   # Ctrl+C để stop
   go run main.go  # Start lại
   ```

### Queue full error?

Tăng queue size trong `main.go`:
```go
costRecalculationService := services.NewCostRecalculationService(
    costCalculatorService, 
    menuRepo, 
    4,
    2000    // Tăng từ 1000 lên 2000
)
```

## Tài Liệu Chi Tiết

- 📖 **Hướng dẫn đầy đủ**: `backend/HUONG_DAN_BACKGROUND_WORKERS.md`
- 📊 **Monitoring guide**: `backend/MONITORING_AND_ALERTS.md`
- 🧪 **Test script**: `./test-background-workers.sh`

## Tóm Tắt

| Câu Hỏi | Trả Lời |
|---------|---------|
| Workers có tự động start không? | ✅ Có, khi server khởi động |
| Cần cấu hình gì không? | ❌ Không, hoạt động ngay |
| Làm sao biết workers đang chạy? | Xem log hoặc check `/monitoring/health` |
| Khi nào workers xử lý? | Khi cập nhật ingredient costs |
| Có monitoring không? | ✅ Có, full metrics và alerts |

**Bottom line**: Workers đã sẵn sàng! Chỉ cần start server và sử dụng. 🚀
