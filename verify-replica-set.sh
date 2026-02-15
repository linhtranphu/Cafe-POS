#!/bin/bash
set -e

# Load environment variables
if [ -f .env.production ]; then
    export $(cat .env.production | grep -v '^#' | xargs)
elif [ -f .env ]; then
    export $(cat .env | grep -v '^#' | xargs)
else
    echo "❌ File .env.production hoặc .env không tồn tại!"
    exit 1
fi

MONGODB_DATABASE=${MONGODB_DATABASE:-cafe_pos}

echo "🔍 Kiểm tra MongoDB Replica Set"
echo "━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━"
echo ""

# Check if container is running
if ! docker ps | grep -q cafe-pos-mongodb; then
    echo "❌ MongoDB container không chạy!"
    echo ""
    echo "Khởi động bằng: docker-compose -f docker-compose.prod.yml up -d mongodb"
    exit 1
fi

echo "✅ MongoDB container đang chạy"
echo ""

# Check replica set status
echo "━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━"
echo "📊 Replica Set Status:"
echo "━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━"
docker exec cafe-pos-mongodb mongosh \
  --username ${MONGO_INITDB_ROOT_USERNAME} \
  --password ${MONGO_INITDB_ROOT_PASSWORD} \
  --authenticationDatabase admin \
  --eval "
  try {
    var status = rs.status();
    print('');
    print('Replica Set Name: ' + status.set);
    print('');
    print('Members:');
    status.members.forEach(function(m) {
      var state = m.stateStr;
      var icon = state === 'PRIMARY' ? '👑' : (state === 'SECONDARY' ? '📋' : '⚠️');
      print('  ' + icon + ' ' + m.name);
      print('     State: ' + state);
      print('     Health: ' + (m.health === 1 ? '✅ Healthy' : '❌ Unhealthy'));
      print('     Uptime: ' + Math.floor(m.uptime / 60) + ' minutes');
      print('');
    });
    
    if (status.ok === 1) {
      print('✅ Replica Set hoạt động bình thường');
    } else {
      print('⚠️  Replica Set có vấn đề');
    }
  } catch(e) {
    print('❌ Lỗi: ' + e);
    print('');
    print('Replica Set chưa được khởi tạo!');
    quit(1);
  }
  " || {
    echo ""
    echo "❌ Không thể kiểm tra Replica Set"
    exit 1
  }

echo ""
echo "━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━"
echo "📊 Database & Collections:"
echo "━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━"
docker exec cafe-pos-mongodb mongosh \
  --username ${MONGO_INITDB_ROOT_USERNAME} \
  --password ${MONGO_INITDB_ROOT_PASSWORD} \
  --authenticationDatabase admin \
  --eval "
  var db = db.getSiblingDB('${MONGODB_DATABASE}');
  var collections = db.getCollectionNames();
  
  print('');
  print('Database: ${MONGODB_DATABASE}');
  print('Total Collections: ' + collections.length);
  print('');
  
  var totalDocs = 0;
  collections.forEach(function(col) {
    var count = db.getCollection(col).countDocuments();
    totalDocs += count;
    var icon = count > 0 ? '📄' : '📭';
    print('  ' + icon + ' ' + col + ': ' + count.toLocaleString() + ' documents');
  });
  
  print('');
  print('Total Documents: ' + totalDocs.toLocaleString());
  "

echo ""
echo "━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━"
echo "🧪 Test Transaction Support:"
echo "━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━"
docker exec cafe-pos-mongodb mongosh \
  --username ${MONGO_INITDB_ROOT_USERNAME} \
  --password ${MONGO_INITDB_ROOT_PASSWORD} \
  --authenticationDatabase admin \
  --eval "
  var db = db.getSiblingDB('${MONGODB_DATABASE}');
  
  print('');
  print('Testing transaction support...');
  print('');
  
  try {
    var session = db.getMongo().startSession();
    session.startTransaction();
    
    // Test transaction
    var testCol = session.getDatabase('${MONGODB_DATABASE}').getCollection('_transaction_test');
    testCol.insertOne({test: true, timestamp: new Date()});
    
    session.commitTransaction();
    session.endSession();
    
    // Cleanup
    db.getCollection('_transaction_test').drop();
    
    print('✅ Transactions hoạt động bình thường!');
    print('   Batch system có thể sử dụng transactions.');
  } catch(e) {
    print('❌ Lỗi transaction: ' + e);
    print('');
    print('⚠️  Batch system có thể không hoạt động đúng!');
    quit(1);
  }
  "

echo ""
echo "━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━"
echo "🌐 Services Status:"
echo "━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━"
echo ""
docker-compose -f docker-compose.prod.yml ps

echo ""
echo "━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━"
echo "✅ Tất cả kiểm tra hoàn tất!"
echo "━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━"
echo ""
echo "💡 Các lệnh hữu ích:"
echo ""
echo "  Xem logs MongoDB:"
echo "    docker logs -f cafe-pos-mongodb"
echo ""
echo "  Xem logs Backend:"
echo "    docker logs -f cafe-pos-backend"
echo ""
echo "  Backup database:"
echo "    ./backup-mongodb.sh"
echo ""
echo "  Truy cập MongoDB shell:"
echo "    docker exec -it cafe-pos-mongodb mongosh -u ${MONGO_INITDB_ROOT_USERNAME} -p"
echo ""
