# 🔍 Phát hiện từ chẩn đoán ban đầu

## Thông tin hệ thống

### Local Machine (Mac)
- CPU: Intel i7-7820HQ @ 2.90GHz (8 cores)
- RAM: 7.656 GiB (Docker limit)
- Disk: 458 GB available
- Docker: 26.1.1

### Containers hiện tại
- ✅ MongoDB: Running (4 days, 381 MB RAM, 4.86% CPU)
- ✅ Print Bridge: Running (2 hours, 26 MB RAM, 0% CPU)
- ❌ Backend: NOT running
- ❌ Frontend: NOT running

## Phát hiện quan trọng

### 1. Backend đang chạy NGOÀI Docker
```
Port 3000 (Backend):
backend 38440 tranphulinh   11u  IPv6 ... TCP *:hbci (LISTEN)
```
- Backend process ID: 38440
- Đang chạy trực tiếp trên host, KHÔNG phải trong container

### 2. MongoDB URI không khớp
```env
# Trong .env
MONGODB_URI=mongodb://admin:password123@localhost:27017/cafe_pos?replicaSet=rs0&authSource=admin
```
- URI trỏ đến `localhost:27017` (OK cho host process)
- Nhưng KHÔNG OK cho container (phải dùng `mongodb:27017`)

### 3. Backend image size tăng đột ngột
```
linhtranphu/cafe-pos-backend:latest    948MB   (mới nhất)
linhtranphu/cafe-pos-backend:2.0       39.8MB  (cũ)
cafe-pos-backend:local                 26.9MB  (local)
```
- Image size tăng từ 27MB → 948MB (x35 lần!)
- Có thể do thêm Chromium và fonts

### 4. Không có resource limits
- Không có memory limits trong docker-compose
- Không có CPU limits
- Container có thể dùng hết RAM

## Giả thuyết về nguyên nhân server EC2 chết

### Giả thuyết 1: Out of Memory (OOM)
- Backend image 948MB + MongoDB + Frontend
- EC2 instance có thể có RAM thấp (t2.micro = 1GB)
- Chromium trong backend tốn nhiều RAM khi render PDF
- **Khả năng cao: 80%**

### Giả thuyết 2: Chromium process không cleanup
- Backend dùng chromedp để render bill
- Nếu không cleanup đúng cách, process tích lũy
- **Khả năng: 60%**

### Giả thuyết 3: MongoDB replica set không init
- .env có `replicaSet=rs0` nhưng MongoDB chưa init
- Backend crash khi cố dùng transactions
- **Khả năng: 40%**

### Giả thuyết 4: Health check fail liên tục
- Health check endpoint `/api/state-machines`
- Nếu fail → Docker restart → fail → restart (loop)
- **Khả năng: 30%**

## Bước tiếp theo

1. ✅ Stop backend đang chạy ngoài Docker
2. ✅ Deploy local với monitoring
3. ✅ Quan sát resource usage
4. ✅ Test các chức năng tốn RAM (print, render PDF)
5. ✅ Xác định nguyên nhân chính xác

## Dự đoán

Nguyên nhân chính: **OOM (Out of Memory)** do:
- Backend image quá lớn (948MB)
- Chromium tốn RAM khi render
- Không có memory limits
- EC2 instance RAM thấp

Giải pháp dự kiến:
- Thêm memory limits
- Optimize Chromium usage
- Tăng RAM cho EC2 hoặc dùng instance lớn hơn
- Cleanup Chromium processes sau mỗi render
