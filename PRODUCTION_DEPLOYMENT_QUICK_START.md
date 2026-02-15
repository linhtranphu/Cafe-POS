# 🚀 Production Deployment - Quick Start Guide

## Tóm Tắt Nhanh

Hướng dẫn deploy production với MongoDB Replica Set (cần thiết cho batch system) mà **KHÔNG MẤT DỮ LIỆU**.

---

## ⚡ Quick Commands

### 1. Backup Dữ Liệu (BẮT BUỘC)
```bash
chmod +x backup-mongodb.sh
./backup-mongodb.sh
```

### 2. Migration sang Replica Set
```bash
chmod +x migrate-to-replica-set.sh
./migrate-to-replica-set.sh
```

### 3. Verify Migration
```bash
chmod +x verify-replica-set.sh
./verify-replica-set.sh
```

---

## 📋 Chi Tiết Từng Bước

### Bước 1: Chuẩn Bị

```bash
# Đảm bảo có file .env với credentials
cat .env

# Cần có:
# MONGO_INITDB_ROOT_USERNAME=admin
# MONGO_INITDB_ROOT_PASSWORD=<your-password>
# MONGODB_DATABASE=cafe_pos
# JWT_SECRET=<your-jwt-secret>
```

### Bước 2: Backup

```bash
# Backup tự động
./backup-mongodb.sh

# Hoặc manual
docker exec cafe-pos-mongodb mongodump \
  --username admin \
  --password <PASSWORD> \
  --authenticationDatabase admin \
  --out /data/backup-$(date +%Y%m%d-%H%M%S)

docker cp cafe-pos-mongodb:/data/backup-<TIMESTAMP> ./mongodb-backup
```

### Bước 3: Migration

```bash
# Chạy migration script
./migrate-to-replica-set.sh

# Script sẽ:
# 1. Kiểm tra backup
# 2. Dừng services
# 3. Tạo keyfile
# 4. Khởi động MongoDB với replica set
# 5. Khởi tạo replica set
# 6. Verify dữ liệu
# 7. Khởi động backend & frontend
```

### Bước 4: Verify

```bash
# Kiểm tra toàn bộ hệ thống
./verify-replica-set.sh

# Kiểm tra logs
docker logs cafe-pos-mongodb
docker logs cafe-pos-backend
docker logs cafe-pos-frontend
```

---

## 🔍 Kiểm Tra Nhanh

### MongoDB Replica Set Status
```bash
docker exec cafe-pos-mongodb mongosh \
  --username admin \
  --password <PASSWORD> \
  --authenticationDatabase admin \
  --eval "rs.status()"
```

Kết quả phải có:
- `"ok": 1`
- `"stateStr": "PRIMARY"`

### Kiểm Tra Dữ Liệu
```bash
docker exec cafe-pos-mongodb mongosh \
  --username admin \
  --password <PASSWORD> \
  --authenticationDatabase admin \
  --eval "db.getSiblingDB('cafe_pos').getCollectionNames()"
```

### Test Backend
```bash
curl http://localhost:3000/api/state-machines
curl http://localhost:3000/api/batch/definitions
```

### Test Frontend
```bash
curl http://localhost
```

---

## 🆘 Troubleshooting

### Lỗi: "not master and slaveOk=false"
```bash
docker exec cafe-pos-mongodb mongosh \
  --username admin \
  --password <PASSWORD> \
  --authenticationDatabase admin \
  --eval "
  cfg = rs.conf();
  cfg.members[0].priority = 2;
  rs.reconfig(cfg, {force: true});
  "
```

### Lỗi: Backend không connect
```bash
# Kiểm tra connection string
docker exec cafe-pos-backend env | grep MONGODB_URI

# Phải có: ?replicaSet=rs0&authSource=admin
```

### Lỗi: Keyfile permissions
```bash
chmod 400 mongodb-keyfile
sudo chown 999:999 mongodb-keyfile
```

---

## 🔄 Rollback (Nếu Cần)

```bash
# 1. Dừng services
docker-compose -f docker-compose.prod.yml down

# 2. Khởi động MongoDB standalone
docker-compose up -d mongodb
sleep 10

# 3. Restore backup
docker cp ./mongodb-backup/backup-<TIMESTAMP> cafe-pos-mongodb:/data/
docker exec cafe-pos-mongodb mongorestore \
  --username admin \
  --password <PASSWORD> \
  --authenticationDatabase admin \
  --drop \
  /data/backup-<TIMESTAMP>

# 4. Khởi động với config cũ
docker-compose -f docker-compose.hub.yml up -d
```

---

## 📊 Monitoring

### Xem Logs Real-time
```bash
# MongoDB
docker logs -f cafe-pos-mongodb

# Backend
docker logs -f cafe-pos-backend

# Frontend
docker logs -f cafe-pos-frontend

# Tất cả
docker-compose -f docker-compose.prod.yml logs -f
```

### Kiểm Tra Resource Usage
```bash
docker stats
```

### Kiểm Tra Replica Set Health
```bash
watch -n 5 'docker exec cafe-pos-mongodb mongosh \
  --username admin \
  --password <PASSWORD> \
  --authenticationDatabase admin \
  --eval "rs.status().members.forEach(m => print(m.name + \": \" + m.stateStr))"'
```

---

## 🔐 Security Checklist

- [ ] Đổi `MONGO_INITDB_ROOT_PASSWORD` trong .env
- [ ] Đổi `JWT_SECRET` trong .env
- [ ] Keyfile có permissions 400
- [ ] Backup được lưu ở nơi an toàn
- [ ] Firewall chỉ cho phép port 80, 443
- [ ] MongoDB port 27017 không expose ra internet (chỉ internal)

---

## 📁 File Structure

```
.
├── docker-compose.prod.yml          # Production config với replica set
├── docker-compose.hub.yml           # Fallback config (standalone)
├── .env                             # Environment variables
├── mongodb-keyfile                  # Replica set keyfile (auto-generated)
├── mongodb-backup/                  # Backup directory
│   └── backup-YYYYMMDD-HHMMSS/
├── migrate-to-replica-set.sh       # Migration script
├── backup-mongodb.sh                # Backup script
├── verify-replica-set.sh            # Verification script
└── MONGODB_REPLICA_MIGRATION_GUIDE.md  # Chi tiết đầy đủ
```

---

## ✅ Checklist Hoàn Thành

### Pre-Migration
- [ ] Đọc MONGODB_REPLICA_MIGRATION_GUIDE.md
- [ ] Có file .env với credentials đúng
- [ ] Backup dữ liệu thành công
- [ ] Backup được copy ra ngoài container

### Migration
- [ ] Chạy migrate-to-replica-set.sh thành công
- [ ] Replica set status = PRIMARY
- [ ] Dữ liệu vẫn còn nguyên (verify bằng count)
- [ ] Backend khởi động không lỗi
- [ ] Frontend truy cập được

### Post-Migration
- [ ] Test batch system (tạo batch definition)
- [ ] Test transactions (verify-replica-set.sh)
- [ ] Logs không có error
- [ ] Setup backup tự động (cron job)

---

## 🎯 Next Steps

### 1. Setup Backup Tự Động
```bash
# Thêm vào crontab
crontab -e

# Backup mỗi ngày lúc 2 giờ sáng
0 2 * * * cd /path/to/project && ./backup-mongodb.sh >> /var/log/mongodb-backup.log 2>&1
```

### 2. Setup Monitoring
- Cài đặt monitoring tool (Prometheus, Grafana)
- Alert khi replica set không healthy
- Alert khi disk space thấp

### 3. Setup Log Rotation
```bash
# Giới hạn log size
docker-compose -f docker-compose.prod.yml up -d --log-opt max-size=10m --log-opt max-file=3
```

---

## 📚 Tài Liệu

- [MONGODB_REPLICA_MIGRATION_GUIDE.md](./MONGODB_REPLICA_MIGRATION_GUIDE.md) - Chi tiết đầy đủ
- [MongoDB Replica Set Docs](https://www.mongodb.com/docs/manual/replication/)
- [MongoDB Transactions](https://www.mongodb.com/docs/manual/core/transactions/)

---

## 💡 Tips

1. **Luôn backup trước khi thay đổi gì**
2. **Test trên staging trước khi production**
3. **Monitor logs sau migration ít nhất 24h**
4. **Giữ backup ít nhất 7 ngày**
5. **Document mọi thay đổi**

---

## 🆘 Support

Nếu gặp vấn đề:
1. Kiểm tra logs: `docker logs cafe-pos-mongodb`
2. Chạy verify: `./verify-replica-set.sh`
3. Xem troubleshooting trong MONGODB_REPLICA_MIGRATION_GUIDE.md
4. Rollback nếu cần thiết
