# Hướng Dẫn Migration MongoDB Replica Set - KHÔNG MẤT DỮ LIỆU

## 🎯 Mục Tiêu
Chuyển từ MongoDB standalone sang Replica Set để hỗ trợ transactions (batch system) mà KHÔNG MẤT dữ liệu production.

## ⚠️ LƯU Ý QUAN TRỌNG
- **BACKUP trước khi làm bất cứ điều gì**
- Replica Set cần thiết cho MongoDB transactions (batch system yêu cầu)
- Downtime: ~5-10 phút trong quá trình migration

---

## 📋 Bước 1: Backup Dữ Liệu Hiện Tại

### 1.1. Backup toàn bộ database
```bash
# Trên server production
docker exec cafe-pos-mongodb mongodump \
  --username admin \
  --password <YOUR_PASSWORD> \
  --authenticationDatabase admin \
  --out /data/backup-$(date +%Y%m%d-%H%M%S)

# Copy backup ra ngoài container (an toàn hơn)
docker cp cafe-pos-mongodb:/data/backup-$(date +%Y%m%d-%H%M%S) ./mongodb-backup
```

### 1.2. Verify backup
```bash
ls -lh ./mongodb-backup
# Phải thấy folder cafe_pos với các collection
```

---

## 📋 Bước 2: Tạo Keyfile cho Replica Set

Replica Set cần keyfile để authentication giữa các nodes:

```bash
# Tạo keyfile với quyền đúng
openssl rand -base64 756 > mongodb-keyfile
chmod 400 mongodb-keyfile
sudo chown 999:999 mongodb-keyfile  # MongoDB user trong container
```

---

## 📋 Bước 3: Tạo Docker Compose Production với Replica Set

Tạo file `docker-compose.prod.yml`:

```yaml
version: '3.8'

services:
  mongodb:
    image: mongo:7.0
    container_name: cafe-pos-mongodb
    restart: always
    environment:
      MONGO_INITDB_ROOT_USERNAME: ${MONGO_INITDB_ROOT_USERNAME}
      MONGO_INITDB_ROOT_PASSWORD: ${MONGO_INITDB_ROOT_PASSWORD}
      MONGO_INITDB_DATABASE: ${MONGODB_DATABASE:-cafe_pos}
    ports:
      - "27017:27017"
    volumes:
      - mongodb_data:/data/db
      - mongodb_config:/data/configdb
      - ./mongodb-keyfile:/data/keyfile:ro
    # QUAN TRỌNG: Thêm replica set config
    command: >
      mongod 
      --replSet rs0 
      --bind_ip_all 
      --keyFile /data/keyfile
      --auth
    networks:
      - cafe-pos-network
    healthcheck:
      test: |
        mongosh --username admin --password ${MONGO_INITDB_ROOT_PASSWORD} --authenticationDatabase admin --eval "
        try {
          var status = rs.status();
          if (status.ok === 1) {
            quit(0);
          }
        } catch(e) {
          quit(1);
        }
        " || exit 1
      interval: 10s
      timeout: 5s
      retries: 10
      start_period: 40s

  backend:
    image: linhtranphu/cafe-pos-backend:latest
    container_name: cafe-pos-backend
    restart: always
    environment:
      # QUAN TRỌNG: Connection string với replicaSet
      - MONGODB_URI=mongodb://${MONGO_INITDB_ROOT_USERNAME}:${MONGO_INITDB_ROOT_PASSWORD}@mongodb:27017/${MONGODB_DATABASE}?replicaSet=rs0&authSource=admin
      - MONGODB_DATABASE=${MONGODB_DATABASE:-cafe_pos}
      - JWT_SECRET=${JWT_SECRET}
      - PORT=3000
    ports:
      - "3000:3000"
    depends_on:
      mongodb:
        condition: service_healthy
    networks:
      - cafe-pos-network
    healthcheck:
      test: ["CMD", "wget", "--quiet", "--tries=1", "--spider", "http://localhost:3000/api/state-machines"]
      interval: 30s
      timeout: 10s
      retries: 3

  frontend:
    image: linhtranphu/cafe-pos-frontend:latest
    container_name: cafe-pos-frontend
    restart: always
    ports:
      - "80:80"
    depends_on:
      - backend
    networks:
      - cafe-pos-network
    healthcheck:
      test: ["CMD", "wget", "--quiet", "--tries=1", "--spider", "http://localhost"]
      interval: 30s
      timeout: 10s
      retries: 3

volumes:
  mongodb_data:
    driver: local
  mongodb_config:
    driver: local

networks:
  cafe-pos-network:
    driver: bridge
```

---

## 📋 Bước 4: Migration Script (KHÔNG MẤT DỮ LIỆU)

Tạo file `migrate-to-replica-set.sh`:

```bash
#!/bin/bash
set -e

echo "🚀 Bắt đầu migration sang MongoDB Replica Set"
echo "⚠️  Đảm bảo bạn đã backup dữ liệu!"
read -p "Đã backup chưa? (yes/no): " confirm

if [ "$confirm" != "yes" ]; then
    echo "❌ Hãy backup trước! Chạy: docker exec cafe-pos-mongodb mongodump ..."
    exit 1
fi

# Load environment variables
if [ -f .env ]; then
    export $(cat .env | grep -v '^#' | xargs)
fi

echo "📊 Kiểm tra dữ liệu hiện tại..."
docker exec cafe-pos-mongodb mongosh \
  --username ${MONGO_INITDB_ROOT_USERNAME} \
  --password ${MONGO_INITDB_ROOT_PASSWORD} \
  --authenticationDatabase admin \
  --eval "db.getSiblingDB('${MONGODB_DATABASE}').getCollectionNames()"

echo ""
echo "🛑 Dừng services hiện tại..."
docker-compose down

echo ""
echo "✅ Dữ liệu vẫn an toàn trong volume: mongodb_data"
docker volume ls | grep mongodb_data

echo ""
echo "🔧 Tạo keyfile..."
if [ ! -f mongodb-keyfile ]; then
    openssl rand -base64 756 > mongodb-keyfile
    chmod 400 mongodb-keyfile
    sudo chown 999:999 mongodb-keyfile
    echo "✅ Keyfile created"
else
    echo "✅ Keyfile already exists"
fi

echo ""
echo "🚀 Khởi động MongoDB với Replica Set config..."
docker-compose -f docker-compose.prod.yml up -d mongodb

echo ""
echo "⏳ Đợi MongoDB khởi động (30 giây)..."
sleep 30

echo ""
echo "🔧 Khởi tạo Replica Set..."
docker exec cafe-pos-mongodb mongosh \
  --username ${MONGO_INITDB_ROOT_USERNAME} \
  --password ${MONGO_INITDB_ROOT_PASSWORD} \
  --authenticationDatabase admin \
  --eval "
  try {
    var status = rs.status();
    print('✅ Replica Set đã được khởi tạo');
  } catch(e) {
    print('🔧 Khởi tạo Replica Set mới...');
    rs.initiate({
      _id: 'rs0',
      members: [
        { _id: 0, host: 'mongodb:27017', priority: 1 }
      ]
    });
    print('✅ Replica Set initialized');
  }
  "

echo ""
echo "⏳ Đợi Replica Set ổn định (10 giây)..."
sleep 10

echo ""
echo "📊 Kiểm tra Replica Set status..."
docker exec cafe-pos-mongodb mongosh \
  --username ${MONGO_INITDB_ROOT_USERNAME} \
  --password ${MONGO_INITDB_ROOT_PASSWORD} \
  --authenticationDatabase admin \
  --eval "rs.status()"

echo ""
echo "📊 Verify dữ liệu vẫn còn..."
docker exec cafe-pos-mongodb mongosh \
  --username ${MONGO_INITDB_ROOT_USERNAME} \
  --password ${MONGO_INITDB_ROOT_PASSWORD} \
  --authenticationDatabase admin \
  --eval "
  db.getSiblingDB('${MONGODB_DATABASE}').getCollectionNames().forEach(function(col) {
    var count = db.getSiblingDB('${MONGODB_DATABASE}').getCollection(col).countDocuments();
    print(col + ': ' + count + ' documents');
  });
  "

echo ""
echo "🚀 Khởi động Backend và Frontend..."
docker-compose -f docker-compose.prod.yml up -d backend frontend

echo ""
echo "⏳ Đợi services khởi động (20 giây)..."
sleep 20

echo ""
echo "✅ Migration hoàn tất!"
echo ""
echo "📊 Kiểm tra services:"
docker-compose -f docker-compose.prod.yml ps

echo ""
echo "🧪 Test transaction (batch system):"
echo "Truy cập: http://localhost:3000/api/batch/definitions"
echo ""
echo "📝 Logs:"
echo "  - MongoDB: docker logs cafe-pos-mongodb"
echo "  - Backend: docker logs cafe-pos-backend"
echo "  - Frontend: docker logs cafe-pos-frontend"
```

Chmod script:
```bash
chmod +x migrate-to-replica-set.sh
```

---

## 📋 Bước 5: Chạy Migration

```bash
# 1. Đảm bảo có file .env với credentials
cat .env

# 2. Chạy migration
./migrate-to-replica-set.sh
```

---

## 📋 Bước 6: Verify Migration Thành Công

### 6.1. Kiểm tra Replica Set
```bash
docker exec cafe-pos-mongodb mongosh \
  --username admin \
  --password <YOUR_PASSWORD> \
  --authenticationDatabase admin \
  --eval "rs.status()"
```

Kết quả phải có:
- `"ok": 1`
- `"stateStr": "PRIMARY"`

### 6.2. Kiểm tra dữ liệu
```bash
docker exec cafe-pos-mongodb mongosh \
  --username admin \
  --password <YOUR_PASSWORD> \
  --authenticationDatabase admin \
  --eval "
  db.getSiblingDB('cafe_pos').getCollectionNames().forEach(function(col) {
    var count = db.getSiblingDB('cafe_pos').getCollection(col).countDocuments();
    print(col + ': ' + count + ' documents');
  });
  "
```

### 6.3. Test transaction (batch system)
```bash
# Test tạo batch definition
curl -X POST http://localhost:3000/api/batch/definitions \
  -H "Content-Type: application/json" \
  -H "Authorization: Bearer <YOUR_TOKEN>" \
  -d '{
    "name": "Test Batch",
    "unit": "lít",
    "conversion_rates": [...]
  }'
```

---

## 🔄 Rollback Plan (Nếu Có Vấn Đề)

### Nếu migration thất bại:

```bash
# 1. Dừng services
docker-compose -f docker-compose.prod.yml down

# 2. Restore từ backup
docker-compose up -d mongodb
sleep 10

docker exec cafe-pos-mongodb mongorestore \
  --username admin \
  --password <YOUR_PASSWORD> \
  --authenticationDatabase admin \
  --drop \
  /data/backup-<TIMESTAMP>

# 3. Khởi động lại với config cũ
docker-compose -f docker-compose.hub.yml up -d
```

---

## 📊 Monitoring sau Migration

### Kiểm tra logs
```bash
# MongoDB logs
docker logs -f cafe-pos-mongodb

# Backend logs (xem transaction errors)
docker logs -f cafe-pos-backend
```

### Kiểm tra replica set health
```bash
# Chạy định kỳ
docker exec cafe-pos-mongodb mongosh \
  --username admin \
  --password <YOUR_PASSWORD> \
  --authenticationDatabase admin \
  --eval "rs.status().members.forEach(m => print(m.name + ': ' + m.stateStr))"
```

---

## ✅ Checklist Hoàn Thành

- [ ] Backup dữ liệu thành công
- [ ] Tạo mongodb-keyfile
- [ ] Tạo docker-compose.prod.yml
- [ ] Chạy migration script
- [ ] Verify replica set status = PRIMARY
- [ ] Verify dữ liệu còn nguyên
- [ ] Test batch system (transactions)
- [ ] Backend logs không có lỗi transaction
- [ ] Frontend hoạt động bình thường

---

## 🆘 Troubleshooting

### Lỗi: "not master and slaveOk=false"
```bash
# Kiểm tra replica set status
docker exec cafe-pos-mongodb mongosh --eval "rs.status()"

# Nếu không phải PRIMARY, force reconfigure
docker exec cafe-pos-mongodb mongosh --eval "
cfg = rs.conf();
cfg.members[0].priority = 2;
rs.reconfig(cfg, {force: true});
"
```

### Lỗi: "keyFile permissions too open"
```bash
chmod 400 mongodb-keyfile
sudo chown 999:999 mongodb-keyfile
```

### Lỗi: Backend không connect được
```bash
# Kiểm tra connection string trong .env
# Phải có: ?replicaSet=rs0&authSource=admin

# Test connection
docker exec cafe-pos-backend wget -O- http://localhost:3000/api/health
```

---

## 📚 Tài Liệu Tham Khảo

- [MongoDB Replica Set Deployment](https://www.mongodb.com/docs/manual/tutorial/deploy-replica-set/)
- [MongoDB Transactions](https://www.mongodb.com/docs/manual/core/transactions/)
- [Docker MongoDB Replica Set](https://www.mongodb.com/compatibility/deploying-a-mongodb-cluster-with-docker)
