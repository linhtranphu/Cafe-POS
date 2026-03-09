# MongoDB Backup & Restore Scripts

Scripts để backup và restore MongoDB database trong Docker container.

## Prerequisites

- Docker đang chạy
- MongoDB container tên `cafe-pos-mongodb` (hoặc sửa tên trong script)
- `mongodump` và `mongorestore` có sẵn trong MongoDB container

## Backup Database

### Cách 1: Backup với timestamp tự động

```bash
./backup-mongodb.sh
```

Sẽ tạo folder `mongodb_backup_YYYYMMDD_HHMMSS` với backup data.

### Cách 2: Backup vào folder chỉ định

```bash
./backup-mongodb.sh ./my_backup_folder
```

## Restore Database

### Restore từ backup directory

```bash
./restore-mongodb.sh ./mongodb_backup_20260307_123456/cafe_pos
```

hoặc nếu backup folder có cấu trúc khác:

```bash
./restore-mongodb.sh ./path/to/backup/cafe_pos
```

### Script sẽ tự động:

1. ✅ Kiểm tra container có đang chạy không
2. ✅ Start container nếu đang stop
3. ✅ Đợi MongoDB sẵn sàng
4. ✅ Copy backup vào container
5. ✅ Restore database (với `--drop` để xóa data cũ)
6. ✅ Initialize replica set nếu chưa có
7. ✅ Verify restore thành công
8. ✅ Clean up temporary files

## Troubleshooting

### Container không start được

```bash
# Xem log
docker logs cafe-pos-mongodb

# Restart container
docker restart cafe-pos-mongodb

# Hoặc tạo container mới
docker run -d --name cafe-pos-mongodb \
  -p 27017:27017 \
  mongo:7.0 --replSet rs0 --bind_ip_all
```

### Disk space đầy

```bash
# Dọn dẹp Docker
docker system prune -a

# Xóa volumes không dùng
docker volume prune
```

### Restore bị lỗi

```bash
# Kiểm tra cấu trúc backup folder
ls -la ./mongodb_backup_folder/

# Phải có cấu trúc:
# mongodb_backup_folder/
#   cafe_pos/
#     collection1.bson
#     collection1.metadata.json
#     collection2.bson
#     ...
```

## Manual Commands

### Backup thủ công

```bash
# Vào container
docker exec -it cafe-pos-mongodb bash

# Backup
mongodump --db cafe_pos --out /tmp/backup

# Exit và copy ra
docker cp cafe-pos-mongodb:/tmp/backup ./my_backup
```

### Restore thủ công

```bash
# Copy backup vào container
docker cp ./my_backup cafe-pos-mongodb:/tmp/backup

# Vào container
docker exec -it cafe-pos-mongodb bash

# Restore
mongorestore --db cafe_pos --drop /tmp/backup/cafe_pos
```

### Initialize Replica Set

```bash
docker exec -it cafe-pos-mongodb mongosh

# Trong mongosh:
rs.initiate({
  _id: 'rs0',
  members: [{ _id: 0, host: 'localhost:27017' }]
})

rs.status()
```

## Notes

- Script sử dụng `--drop` khi restore, nghĩa là sẽ xóa toàn bộ data cũ
- Backup không bao gồm system collections
- Replica set cần được initialize để backend hoạt động đúng
- Container name mặc định là `cafe-pos-mongodb`, có thể sửa trong script nếu khác
