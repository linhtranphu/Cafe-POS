# Hướng dẫn sửa Timezone cho EC2

## Vấn đề
Khi chạy trên EC2, thời gian trên hóa đơn có thể không đúng với giờ Việt Nam do:
- EC2 server mặc định dùng UTC timezone
- Docker container cũng dùng UTC timezone

## Giải pháp đã áp dụng

### 1. Code Level (backend/main.go)
```go
// Set timezone to Vietnam
loc, err := time.LoadLocation("Asia/Ho_Chi_Minh")
if err != nil {
    log.Printf("⚠️  Warning: Could not load Asia/Ho_Chi_Minh timezone: %v", err)
} else {
    time.Local = loc
    log.Println("✅ Timezone set to Asia/Ho_Chi_Minh")
}
```

### 2. Dockerfile Level (backend/Dockerfile)
```dockerfile
# Set timezone to Vietnam
ENV TZ=Asia/Ho_Chi_Minh
```

### 3. Docker Compose Level (docker-compose.yml)
```yaml
backend:
  environment:
    - TZ=Asia/Ho_Chi_Minh
```

### 4. Environment File (.env.ec2)
```bash
TZ=Asia/Ho_Chi_Minh
```

## Cách kiểm tra trên EC2

### Bước 1: Kiểm tra timezone hiện tại
```bash
./check-timezone.sh
```

Script này sẽ hiển thị:
- System timezone
- Backend container timezone
- MongoDB container timezone
- So sánh với giờ Việt Nam

### Bước 2: Sửa timezone nếu sai
```bash
sudo ./fix-ec2-timezone.sh
```

Script này sẽ:
1. Set system timezone thành Asia/Ho_Chi_Minh
2. Restart backend container
3. Verify timezone đã đúng

### Bước 3: Deploy lại backend với timezone mới

```bash
# Build và push image mới
./build_docker_hub.sh

# Trên EC2, pull và restart
docker-compose pull backend
docker-compose up -d backend
```

## Kiểm tra kết quả

1. Tạo một order mới
2. In hóa đơn
3. Kiểm tra thời gian trên hóa đơn:
   - Format: `02/01/2006 15:04` (24 giờ)
   - Timezone: UTC+7 (giờ Việt Nam)

## Lưu ý

- Thời gian được lưu trong MongoDB vẫn là UTC (chuẩn quốc tế)
- Chỉ khi hiển thị mới convert sang giờ Việt Nam
- Điều này đảm bảo tính nhất quán khi có nhiều múi giờ

## Troubleshooting

### Vấn đề: Thời gian vẫn sai sau khi deploy
**Giải pháp:**
1. Kiểm tra backend logs: `docker logs cafe-pos-backend | grep timezone`
2. Verify TZ environment: `docker exec cafe-pos-backend printenv TZ`
3. Restart container: `docker restart cafe-pos-backend`

### Vấn đề: System timezone không thay đổi được
**Giải pháp:**
```bash
# Trên EC2 Ubuntu/Amazon Linux
sudo timedatectl set-timezone Asia/Ho_Chi_Minh

# Verify
timedatectl
```

### Vấn đề: Container timezone không đúng
**Giải pháp:**
1. Rebuild image với Dockerfile mới
2. Hoặc mount timezone file:
```yaml
volumes:
  - /etc/localtime:/etc/localtime:ro
  - /etc/timezone:/etc/timezone:ro
```
