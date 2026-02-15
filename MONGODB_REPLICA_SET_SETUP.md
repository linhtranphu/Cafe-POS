# MongoDB Replica Set Setup Guide

## Tại Sao Cần Replica Set?

MongoDB transactions (cần cho batch management system) chỉ hoạt động trên:
- Replica Set (khuyến nghị)
- Sharded Cluster

Hiện tại MongoDB đang chạy ở standalone mode, không hỗ trợ transactions.

## Lỗi Hiện Tại

```
Transaction numbers are only allowed on a replica set member or mongos
```

## Giải Pháp

### Option 1: Cấu Hình Replica Set Đơn Giản (Development/Testing)

Đây là cách đơn giản nhất cho môi trường development và testing.

#### Bước 1: Dừng MongoDB Hiện Tại

```bash
# Nếu đang chạy qua Docker Compose
docker-compose down

# Hoặc nếu chạy standalone
brew services stop mongodb-community
# hoặc
sudo systemctl stop mongod
```

#### Bước 2: Cập Nhật docker-compose.yml

Thêm cấu hình replica set vào MongoDB service:

```yaml
version: '3.8'

services:
  mongodb:
    image: mongo:6.0
    container_name: cafe-pos-mongodb
    ports:
      - "27017:27017"
    environment:
      MONGO_INITDB_ROOT_USERNAME: admin
      MONGO_INITDB_ROOT_PASSWORD: password
      MONGO_INITDB_DATABASE: cafe_pos
    volumes:
      - mongodb_data:/data/db
      - ./scripts/init-replica-set.sh:/docker-entrypoint-initdb.d/init-replica-set.sh
    command: ["--replSet", "rs0", "--bind_ip_all"]
    healthcheck:
      test: echo "try { rs.status() } catch (err) { rs.initiate({_id:'rs0',members:[{_id:0,host:'localhost:27017'}]}) }" | mongosh --port 27017 --quiet
      interval: 5s
      timeout: 30s
      start_period: 0s
      start_interval: 1s
      retries: 30

  backend:
    # ... existing backend config
    depends_on:
      mongodb:
        condition: service_healthy

volumes:
  mongodb_data:
```

#### Bước 3: Tạo Script Khởi Tạo Replica Set

Tạo file `scripts/init-replica-set.sh`:

```bash
#!/bin/bash

echo "Waiting for MongoDB to start..."
sleep 10

echo "Initiating replica set..."
mongosh --host localhost:27017 <<EOF
var config = {
    "_id": "rs0",
    "version": 1,
    "members": [
        {
            "_id": 0,
            "host": "localhost:27017",
            "priority": 1
        }
    ]
};
rs.initiate(config, { force: true });
EOF

echo "Replica set initiated!"
```

#### Bước 4: Khởi Động Lại

```bash
# Tạo script executable
chmod +x scripts/init-replica-set.sh

# Khởi động
docker-compose up -d mongodb

# Đợi MongoDB khởi động và replica set được khởi tạo
sleep 15

# Kiểm tra status
docker exec cafe-pos-mongodb mongosh --eval "rs.status()"
```

#### Bước 5: Verify Replica Set

```bash
# Kết nối và kiểm tra
docker exec -it cafe-pos-mongodb mongosh

# Trong mongosh shell:
rs.status()
# Nên thấy: "ok": 1 và state: 1 (PRIMARY)

# Test transaction
use cafe_pos
db.test.insertOne({test: 1})
session = db.getMongo().startSession()
session.startTransaction()
session.getDatabase("cafe_pos").test.insertOne({test: 2})
session.commitTransaction()
# Nếu không có lỗi, transactions đang hoạt động!
```

### Option 2: Replica Set Với Docker Compose (Recommended for CI/CD)

File `docker-compose.test.yml` cho testing:

```yaml
version: '3.8'

services:
  mongodb-primary:
    image: mongo:6.0
    container_name: mongodb-primary
    command: ["--replSet", "rs0", "--bind_ip_all", "--port", "27017"]
    ports:
      - "27017:27017"
    environment:
      MONGO_INITDB_ROOT_USERNAME: admin
      MONGO_INITDB_ROOT_PASSWORD: password
    volumes:
      - mongodb_primary_data:/data/db
    healthcheck:
      test: echo "try { rs.status() } catch (err) { rs.initiate({_id:'rs0',members:[{_id:0,host:'mongodb-primary:27017'}]}) }" | mongosh --port 27017 --quiet
      interval: 5s
      timeout: 30s
      retries: 30

  mongodb-init:
    image: mongo:6.0
    depends_on:
      mongodb-primary:
        condition: service_healthy
    command: >
      mongosh --host mongodb-primary:27017 --eval "
        var config = {
          _id: 'rs0',
          members: [
            {_id: 0, host: 'mongodb-primary:27017'}
          ]
        };
        try {
          rs.status();
          print('Replica set already initialized');
        } catch(e) {
          rs.initiate(config);
          print('Replica set initialized');
        }
      "
    restart: "no"

volumes:
  mongodb_primary_data:
```

Chạy:

```bash
docker-compose -f docker-compose.test.yml up -d
```

### Option 3: Standalone MongoDB với Replica Set (Local Development)

Nếu đang chạy MongoDB local (không dùng Docker):

#### macOS (Homebrew)

```bash
# 1. Dừng MongoDB
brew services stop mongodb-community

# 2. Tạo thư mục data cho replica set
mkdir -p ~/mongodb-replica-set/data

# 3. Khởi động MongoDB với replica set
mongod --replSet rs0 --port 27017 --dbpath ~/mongodb-replica-set/data --bind_ip localhost

# 4. Trong terminal khác, khởi tạo replica set
mongosh --eval "rs.initiate({_id: 'rs0', members: [{_id: 0, host: 'localhost:27017'}]})"

# 5. Verify
mongosh --eval "rs.status()"
```

#### Linux

```bash
# 1. Dừng MongoDB
sudo systemctl stop mongod

# 2. Sửa file config
sudo nano /etc/mongod.conf

# Thêm/sửa:
replication:
  replSetName: "rs0"

# 3. Khởi động lại
sudo systemctl start mongod

# 4. Khởi tạo replica set
mongosh --eval "rs.initiate({_id: 'rs0', members: [{_id: 0, host: 'localhost:27017'}]})"

# 5. Verify
mongosh --eval "rs.status()"
```

## Cập Nhật Connection String

Sau khi cấu hình replica set, cập nhật connection string trong code:

### Trước (Standalone):
```
mongodb://localhost:27017/cafe_pos
```

### Sau (Replica Set):
```
mongodb://localhost:27017/cafe_pos?replicaSet=rs0
```

Hoặc với authentication:
```
mongodb://admin:password@localhost:27017/cafe_pos?replicaSet=rs0&authSource=admin
```

## Cập Nhật Backend Code

File `.env` hoặc config:

```bash
# Development
MONGODB_URI=mongodb://localhost:27017/cafe_pos?replicaSet=rs0

# Docker
MONGODB_URI=mongodb://mongodb:27017/cafe_pos?replicaSet=rs0

# With auth
MONGODB_URI=mongodb://admin:password@mongodb:27017/cafe_pos?replicaSet=rs0&authSource=admin
```

## Testing Script

Tạo file `test-mongodb-transactions.sh`:

```bash
#!/bin/bash

echo "Testing MongoDB Transactions..."

mongosh --eval "
use cafe_pos;

// Test basic connection
print('✓ Connected to MongoDB');

// Test replica set status
var status = rs.status();
if (status.ok === 1) {
    print('✓ Replica set is active');
    print('  - Set name: ' + status.set);
    print('  - Primary: ' + status.members.find(m => m.stateStr === 'PRIMARY').name);
} else {
    print('✗ Replica set not configured');
    quit(1);
}

// Test transaction
try {
    var session = db.getMongo().startSession();
    session.startTransaction();
    
    session.getDatabase('cafe_pos').test_transactions.insertOne({
        test: 'transaction_test',
        timestamp: new Date()
    });
    
    session.commitTransaction();
    session.endSession();
    
    print('✓ Transactions are working!');
    
    // Cleanup
    db.test_transactions.drop();
} catch(e) {
    print('✗ Transaction failed: ' + e);
    quit(1);
}

print('\\n✅ All tests passed! MongoDB is ready for batch management.');
"
```

Chạy test:

```bash
chmod +x test-mongodb-transactions.sh
./test-mongodb-transactions.sh
```

## Chạy Lại Backend Tests

Sau khi cấu hình replica set:

```bash
cd backend

# Chạy tất cả batch tests
go test -v -run="Batch" ./application/services/...

# Chạy test transaction cụ thể
go test -v -run="TestProperty_BatchCreationSuccess" ./application/services/...
```

## Troubleshooting

### Lỗi: "not master and slaveOk=false"

```bash
# Trong mongosh
rs.status()
# Kiểm tra xem node có phải PRIMARY không

# Nếu không, force reconfigure
rs.reconfig(rs.conf(), {force: true})
```

### Lỗi: "no reachable servers"

```bash
# Kiểm tra MongoDB đang chạy
docker ps | grep mongodb
# hoặc
ps aux | grep mongod

# Kiểm tra port
netstat -an | grep 27017
```

### Lỗi: "connection refused"

```bash
# Kiểm tra bind_ip
# MongoDB phải bind_ip_all hoặc bind đúng interface

# Trong docker-compose
command: ["--replSet", "rs0", "--bind_ip_all"]
```

### Reset Replica Set

Nếu cần reset hoàn toàn:

```bash
# Dừng MongoDB
docker-compose down -v  # -v để xóa volumes

# Hoặc local
rm -rf ~/mongodb-replica-set/data/*

# Khởi động lại và init lại
```

## Production Considerations

Cho production, nên dùng:

1. **3-node Replica Set** (minimum recommended):
   - 1 Primary
   - 2 Secondary (cho high availability)

2. **Separate Servers**: Mỗi node trên server riêng

3. **Monitoring**: Sử dụng MongoDB Atlas hoặc monitoring tools

4. **Backups**: Automated backups từ secondary nodes

5. **Security**: 
   - Enable authentication
   - Use TLS/SSL
   - Network isolation

## Quick Start Commands

```bash
# 1. Cập nhật docker-compose.yml (xem Option 1 ở trên)

# 2. Khởi động
docker-compose down
docker-compose up -d mongodb

# 3. Đợi khởi động
sleep 15

# 4. Verify
docker exec cafe-pos-mongodb mongosh --eval "rs.status()"

# 5. Test transactions
./test-mongodb-transactions.sh

# 6. Chạy backend tests
cd backend && go test -v -run="Batch" ./application/services/...
```

## Kết Luận

Sau khi cấu hình replica set:
- ✅ Transactions sẽ hoạt động
- ✅ Tất cả batch tests sẽ pass
- ✅ Production-ready với high availability
- ✅ Hỗ trợ ACID transactions cho batch operations

**Khuyến nghị:** Sử dụng Option 1 (Docker Compose với healthcheck) cho development và testing - đơn giản và reliable nhất.
