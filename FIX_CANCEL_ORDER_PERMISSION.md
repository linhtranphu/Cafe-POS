# Fix: Lỗi "insufficient permissions" khi Waiter hủy order

## Vấn đề

Khi waiter cố gắng hủy order (chưa thanh toán), hệ thống báo lỗi:
```
❌ Lỗi: insufficient permissions
```

## Nguyên nhân

Backend đang chạy code cũ chưa có endpoint `/waiter/orders/:id/cancel`. Code mới đã được thêm vào `backend/main.go` nhưng backend container chưa được rebuild với code mới.

## Giải pháp

### Cách 1: Rebuild Backend (Khuyến nghị)

```bash
# Chạy script tự động
./rebuild-backend-with-cancel.sh
```

Hoặc thủ công:

```bash
# 1. Stop backend
docker-compose stop backend

# 2. Remove container
docker-compose rm -f backend

# 3. Rebuild (no cache)
cd backend
docker build --no-cache -t cafe-pos-backend:latest .
cd ..

# 4. Start backend
docker-compose up -d backend

# 5. Check logs
docker logs -f cafe-pos-backend
```

### Cách 2: Rebuild toàn bộ stack

```bash
# Stop all
docker-compose down

# Rebuild all
docker-compose up -d --build

# Check status
docker-compose ps
```

### Cách 3: Deploy lên production (EC2)

Nếu đang test trên production:

```bash
# 1. SSH vào EC2
ssh your-ec2-instance

# 2. Pull code mới
cd /path/to/cafe-pos
git pull origin main

# 3. Rebuild backend
cd backend
docker build --no-cache -t cafe-pos-backend:latest .
cd ..

# 4. Restart
docker-compose up -d backend

# 5. Check
docker logs -f cafe-pos-backend
```

## Kiểm tra

### 1. Kiểm tra endpoint đã được đăng ký

```bash
# Check if cancel endpoint exists in code
grep -n "waiter.POST.*cancel" backend/main.go

# Should show:
# 523:    waiter.POST("/orders/:id/cancel", orderHandler.CancelOrder)
```

### 2. Kiểm tra backend đang chạy code mới

```bash
# Check backend logs for waiter routes
docker logs cafe-pos-backend 2>&1 | grep -i "waiter"

# Or check the running code
docker exec cafe-pos-backend cat /app/main.go | grep "waiter.POST.*cancel"
```

### 3. Test API trực tiếp

```bash
# Get auth token first (login as waiter)
TOKEN="your-jwt-token"

# Try to cancel an order
curl -X POST \
  http://localhost:8080/api/waiter/orders/ORDER_ID/cancel \
  -H "Authorization: Bearer $TOKEN" \
  -H "Content-Type: application/json" \
  -d '{"reason": "Test cancel"}'

# Should return order with status CANCELLED
# If returns 403 or "insufficient permissions", backend needs rebuild
```

## Xác nhận đã fix

Sau khi rebuild, test lại:

1. Login với tài khoản waiter
2. Tạo một order mới (chưa thanh toán)
3. Mở order detail
4. Tap nút "❌ Hủy order"
5. Nhập lý do hủy
6. Tap "Xác nhận hủy"
7. ✅ Order chuyển sang trạng thái CANCELLED

## Code đã thêm

### Backend (`backend/main.go`)

```go
// Waiter routes
waiter := protected.Group("/waiter")
waiter.Use(http.RequireRole(user.RoleWaiter, user.RoleCashier, user.RoleManager))
{
    // ... other routes ...
    waiter.POST("/orders/:id/cancel", orderHandler.CancelOrder)  // ← NEW
    // ... other routes ...
}
```

### Frontend (`frontend/src/services/order.js`)

```javascript
async cancelOrder(id, reason) {
  const response = await api.post(`/waiter/orders/${id}/cancel`, { reason })
  return response.data
}
```

## Lưu ý

1. **Không cần thay đổi database** - Chỉ cần rebuild backend
2. **Frontend không cần rebuild** - Đã dùng đúng endpoint
3. **Chỉ cần restart backend** - Không ảnh hưởng services khác
4. **Downtime tối thiểu** - Chỉ vài giây khi restart backend

## Troubleshooting

### Vẫn báo lỗi sau khi rebuild?

1. **Kiểm tra container đang chạy:**
   ```bash
   docker ps | grep backend
   ```

2. **Kiểm tra logs có lỗi:**
   ```bash
   docker logs cafe-pos-backend 2>&1 | tail -50
   ```

3. **Kiểm tra code trong container:**
   ```bash
   docker exec cafe-pos-backend cat /app/main.go | grep -A 5 "Waiter routes"
   ```

4. **Clear browser cache:**
   - Ctrl+Shift+R (hard refresh)
   - Hoặc clear cache trong DevTools

5. **Kiểm tra JWT token:**
   - Token có thể đã expire
   - Logout và login lại

### Backend không start được?

```bash
# Check logs
docker logs cafe-pos-backend

# Common issues:
# - Port 8080 already in use
# - MongoDB connection failed
# - Environment variables missing

# Fix: Check .env file
cat backend/.env
```

## Scripts hỗ trợ

1. **check-cancel-endpoint.sh** - Kiểm tra endpoint có được đăng ký không
2. **rebuild-backend-with-cancel.sh** - Tự động rebuild backend

---

**Ngày tạo:** 4 tháng 3, 2026
**Trạng thái:** Đã fix - Cần rebuild backend
