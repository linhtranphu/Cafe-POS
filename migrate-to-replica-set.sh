#!/bin/bash
set -e

echo "🚀 Bắt đầu migration sang MongoDB Replica Set"
echo "⚠️  Đảm bảo bạn đã backup dữ liệu!"
echo ""
read -p "Đã backup chưa? (yes/no): " confirm

if [ "$confirm" != "yes" ]; then
    echo "❌ Hãy backup trước!"
    echo ""
    echo "Chạy lệnh sau để backup:"
    echo "docker exec cafe-pos-mongodb mongodump --username admin --password <PASSWORD> --authenticationDatabase admin --out /data/backup-\$(date +%Y%m%d-%H%M%S)"
    echo "docker cp cafe-pos-mongodb:/data/backup-<TIMESTAMP> ./mongodb-backup"
    exit 1
fi

# Load environment variables
if [ -f .env.production ]; then
    export $(cat .env.production | grep -v '^#' | xargs)
elif [ -f .env ]; then
    export $(cat .env | grep -v '^#' | xargs)
else
    echo "❌ File .env.production hoặc .env không tồn tại!"
    exit 1
fi

# Validate required env vars
if [ -z "$MONGO_INITDB_ROOT_USERNAME" ] || [ -z "$MONGO_INITDB_ROOT_PASSWORD" ]; then
    echo "❌ Thiếu MONGO_INITDB_ROOT_USERNAME hoặc MONGO_INITDB_ROOT_PASSWORD trong .env"
    exit 1
fi

MONGODB_DATABASE=${MONGODB_DATABASE:-cafe_pos}

echo ""
echo "📊 Kiểm tra dữ liệu hiện tại..."
docker exec cafe-pos-mongodb mongosh \
  --username ${MONGO_INITDB_ROOT_USERNAME} \
  --password ${MONGO_INITDB_ROOT_PASSWORD} \
  --authenticationDatabase admin \
  --eval "db.getSiblingDB('${MONGODB_DATABASE}').getCollectionNames()" || {
    echo "⚠️  Không thể kết nối MongoDB. Container có đang chạy không?"
    docker ps | grep mongodb
    exit 1
  }

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
    
    # Try to change owner (may need sudo)
    if command -v sudo &> /dev/null; then
        sudo chown 999:999 mongodb-keyfile 2>/dev/null || {
            echo "⚠️  Không thể chown keyfile. Sẽ thử tiếp..."
        }
    fi
    echo "✅ Keyfile created"
else
    echo "✅ Keyfile already exists"
    chmod 400 mongodb-keyfile
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
    print('✅ Replica Set đã được khởi tạo: ' + status.set);
    print('   State: ' + status.members[0].stateStr);
  } catch(e) {
    print('🔧 Khởi tạo Replica Set mới...');
    var result = rs.initiate({
      _id: 'rs0',
      members: [
        { _id: 0, host: 'mongodb:27017', priority: 1 }
      ]
    });
    print('✅ Replica Set initialized');
    print(JSON.stringify(result));
  }
  "

echo ""
echo "⏳ Đợi Replica Set ổn định (15 giây)..."
sleep 15

echo ""
echo "📊 Kiểm tra Replica Set status..."
docker exec cafe-pos-mongodb mongosh \
  --username ${MONGO_INITDB_ROOT_USERNAME} \
  --password ${MONGO_INITDB_ROOT_PASSWORD} \
  --authenticationDatabase admin \
  --eval "
  var status = rs.status();
  print('Replica Set: ' + status.set);
  print('Members:');
  status.members.forEach(function(m) {
    print('  - ' + m.name + ': ' + m.stateStr + ' (health: ' + m.health + ')');
  });
  "

echo ""
echo "📊 Verify dữ liệu vẫn còn..."
docker exec cafe-pos-mongodb mongosh \
  --username ${MONGO_INITDB_ROOT_USERNAME} \
  --password ${MONGO_INITDB_ROOT_PASSWORD} \
  --authenticationDatabase admin \
  --eval "
  var db = db.getSiblingDB('${MONGODB_DATABASE}');
  var collections = db.getCollectionNames();
  print('Database: ${MONGODB_DATABASE}');
  print('Collections: ' + collections.length);
  collections.forEach(function(col) {
    var count = db.getCollection(col).countDocuments();
    print('  - ' + col + ': ' + count + ' documents');
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
echo "━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━"
echo "📊 Kiểm tra services:"
echo "━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━"
docker-compose -f docker-compose.prod.yml ps

echo ""
echo "━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━"
echo "🧪 Các bước tiếp theo:"
echo "━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━"
echo ""
echo "1. Kiểm tra logs:"
echo "   docker logs cafe-pos-mongodb"
echo "   docker logs cafe-pos-backend"
echo "   docker logs cafe-pos-frontend"
echo ""
echo "2. Test batch system (transactions):"
echo "   curl http://localhost:3000/api/batch/definitions"
echo ""
echo "3. Truy cập ứng dụng:"
echo "   http://localhost"
echo ""
echo "4. Nếu có vấn đề, xem MONGODB_REPLICA_MIGRATION_GUIDE.md"
echo "   phần Rollback Plan"
echo ""
echo "━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━"
