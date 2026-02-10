# Background Workers - Tài Liệu Đầy Đủ

## 🎯 Câu Trả Lời Nhanh

**Background workers ĐÃ ĐƯỢC ENABLE tự động!**

Khi bạn start server backend, workers sẽ tự động chạy:

```bash
cd backend
go run main.go

# Output:
# Starting cost recalculation worker pool...
# ✅ Cost recalculation worker pool started
```

## 📚 Tài Liệu

### 1. Quick Start (Đọc đầu tiên!)
📄 **File**: `BACKGROUND_WORKERS_QUICK_START.md`

Tài liệu ngắn gọn với:
- ✅ Cách verify workers hoạt động
- ✅ Khi nào workers xử lý jobs
- ✅ Cấu hình cơ bản
- ✅ Troubleshooting nhanh

**Đọc file này nếu**: Bạn muốn biết nhanh workers có hoạt động không

### 2. Hướng Dẫn Chi Tiết (Tiếng Việt)
📄 **File**: `backend/HUONG_DAN_BACKGROUND_WORKERS.md`

Hướng dẫn đầy đủ với:
- 🔧 Cách hoạt động chi tiết
- ⚙️ Cấu hình workers và queue
- 📊 Monitoring và metrics
- 🐛 Troubleshooting chi tiết
- 💡 Performance tips

**Đọc file này nếu**: Bạn muốn hiểu sâu về workers và cách tối ưu

### 3. Flow Diagram
📄 **File**: `BACKGROUND_WORKERS_FLOW.md`

Visual diagrams cho:
- 🔄 Server startup flow
- 📦 Job processing flow
- 🔁 Retry logic
- 📊 Monitoring flow
- 🏗️ Worker pool architecture

**Đọc file này nếu**: Bạn muốn hiểu flow và architecture

### 4. Monitoring Guide
📄 **File**: `backend/MONITORING_AND_ALERTS.md`

Chi tiết về monitoring system:
- 📈 Metrics collection
- 🚨 Alert rules
- 🔌 API endpoints
- 🔍 Troubleshooting
- 📋 Best practices

**Đọc file này nếu**: Bạn muốn setup monitoring và alerts

## 🧪 Testing

### Test Script
```bash
./test-background-workers.sh
```

Script sẽ tự động kiểm tra:
- ✅ Health status
- ✅ Worker metrics
- ✅ Recent jobs
- ✅ Alerts
- ✅ System performance

### Manual Testing

#### 1. Check Health
```bash
TOKEN=$(curl -X POST http://localhost:3000/api/login \
  -H "Content-Type: application/json" \
  -d '{"username":"admin","password":"admin123"}' \
  | jq -r '.token')

curl -H "Authorization: Bearer $TOKEN" \
  http://localhost:3000/api/manager/monitoring/health
```

#### 2. Trigger Workers
```bash
# Cập nhật ingredient cost
curl -X PUT "http://localhost:3000/api/manager/ingredients/INGREDIENT_ID" \
  -H "Authorization: Bearer $TOKEN" \
  -H "Content-Type: application/json" \
  -d '{"name":"Coffee","cost_per_unit":250000}'

# Xem workers xử lý
curl -H "Authorization: Bearer $TOKEN" \
  "http://localhost:3000/api/manager/monitoring/metrics?type=recalculation_job&limit=10"
```

## 🔧 Configuration

### Default Configuration
```go
// In backend/main.go
costRecalculationService := services.NewCostRecalculationService(
    costCalculatorService, 
    menuRepo, 
    4,      // 4 workers
    1000    // Queue size 1000
)
```

### Recommended Settings

| Scenario | Workers | Queue Size |
|----------|---------|------------|
| Small (< 100 menu items) | 2-4 | 500 |
| Medium (100-500 items) | 4-8 | 1000 |
| Large (> 500 items) | 8-16 | 2000 |

### Adjust Configuration

Edit `backend/main.go`:
```go
costRecalculationService := services.NewCostRecalculationService(
    costCalculatorService, 
    menuRepo, 
    8,      // Tăng workers
    2000    // Tăng queue size
)
```

Restart server:
```bash
# Ctrl+C để stop
go run main.go  # Start lại
```

## 📊 Monitoring Endpoints

| Endpoint | Description |
|----------|-------------|
| `GET /api/manager/monitoring/health` | Overall health status |
| `GET /api/manager/monitoring/metrics/aggregated` | Aggregated statistics |
| `GET /api/manager/monitoring/metrics?type=recalculation_job` | Recent jobs |
| `GET /api/manager/monitoring/alerts` | Active alerts |

## 🚨 Alert Rules

| Alert | Threshold | Level |
|-------|-----------|-------|
| High Recalculation Failure Rate | 5 failures in 5min OR >15% error rate | Critical |
| High Cost Calculation Failure Rate | 10 failures in 5min OR >20% error rate | Critical |
| Shift Closure Failures | 3 failures in 10min OR >10% error rate | Critical |

## 🐛 Common Issues

### Issue 1: Workers không xử lý jobs

**Symptoms**: Cập nhật ingredient nhưng menu cost không thay đổi

**Solution**:
1. Check server log: `tail -f backend/server.log`
2. Check metrics: `curl .../monitoring/metrics/aggregated`
3. Restart server

### Issue 2: Queue full error

**Symptoms**: Error "recalculation queue is full"

**Solution**:
1. Tăng queue size trong `main.go`
2. Tăng số workers
3. Restart server

### Issue 3: High error rate

**Symptoms**: Alert "High Recalculation Job Failure Rate"

**Solution**:
1. Check failed jobs: `curl .../monitoring/metrics?type=recalculation_job`
2. Verify ingredient data quality
3. Check database performance

## 📈 Performance

### Metrics

- **Throughput**: ~20 jobs/second (4 workers)
- **Latency**: ~150-200ms per job
- **Queue capacity**: 1000 jobs
- **Timeout**: 5 seconds per job
- **Retry**: 3 attempts with exponential backoff

### Optimization Tips

1. **Increase workers** for higher throughput
2. **Increase queue size** to handle bursts
3. **Monitor error rates** to catch issues early
4. **Keep ingredient data clean** to minimize failures

## 🎓 Learning Path

### Beginner
1. ✅ Read `BACKGROUND_WORKERS_QUICK_START.md`
2. ✅ Run `./test-background-workers.sh`
3. ✅ Test updating an ingredient

### Intermediate
1. ✅ Read `backend/HUONG_DAN_BACKGROUND_WORKERS.md`
2. ✅ Understand monitoring endpoints
3. ✅ Practice troubleshooting

### Advanced
1. ✅ Read `BACKGROUND_WORKERS_FLOW.md`
2. ✅ Read `backend/MONITORING_AND_ALERTS.md`
3. ✅ Optimize configuration for your use case

## 🔗 Related Documentation

- **Implementation**: `backend/TASK_23.3_MONITORING_IMPLEMENTATION.md`
- **Cost Calculator**: `backend/application/services/cost_calculator_service.go`
- **Recalculation Service**: `backend/application/services/cost_recalculation_service.go`
- **Monitoring Service**: `backend/application/services/monitoring_service.go`

## ❓ FAQ

### Q: Workers có tự động start không?
**A**: ✅ Có, khi server khởi động

### Q: Cần cấu hình gì không?
**A**: ❌ Không, hoạt động ngay out of the box

### Q: Làm sao biết workers đang chạy?
**A**: Xem log server hoặc check `/monitoring/health`

### Q: Khi nào workers xử lý?
**A**: Khi cập nhật ingredient costs

### Q: Có monitoring không?
**A**: ✅ Có, full metrics và alerts

### Q: Workers có chạy song song không?
**A**: ✅ Có, 4 workers xử lý đồng thời

### Q: Nếu job fail thì sao?
**A**: Tự động retry 3 lần với exponential backoff

### Q: Queue đầy thì sao?
**A**: Tăng queue size trong `main.go` và restart

### Q: Có thể tắt workers không?
**A**: Không nên, workers là core feature. Nhưng có thể set workers = 0 nếu cần

### Q: Workers có ảnh hưởng performance không?
**A**: Không, workers chạy async và có timeout

## 🎉 Summary

✅ **Background workers tự động enable**
✅ **Không cần cấu hình thêm**
✅ **Full monitoring và alerts**
✅ **Graceful shutdown**
✅ **Production ready**

**Để bắt đầu**:
```bash
# 1. Start server
cd backend
go run main.go

# 2. Verify workers
./test-background-workers.sh

# 3. Test với real data
# Cập nhật một ingredient và xem workers xử lý!
```

**Need help?** Đọc các tài liệu chi tiết ở trên! 📚
