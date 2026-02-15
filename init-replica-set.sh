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

echo "🔧 Khởi tạo MongoDB Replica Set"
echo "━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━"
echo ""

# Credentials
MONGO_USER="${MONGO_INITDB_ROOT_USERNAME}"
MONGO_PASS="${MONGO_INITDB_ROOT_PASSWORD}"

echo "📊 Kiểm tra Replica Set status..."
docker exec cafe-pos-mongodb mongosh \
  --username "$MONGO_USER" \
  --password "$MONGO_PASS" \
  --authenticationDatabase admin \
  --eval "
  try {
    var status = rs.status();
    print('✅ Replica Set đã được khởi tạo: ' + status.set);
    print('   State: ' + status.members[0].stateStr);
    quit(0);
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
    quit(0);
  }
  "

echo ""
echo "⏳ Đợi Replica Set ổn định (10 giây)..."
sleep 10

echo ""
echo "📊 Verify Replica Set status..."
docker exec cafe-pos-mongodb mongosh \
  --username "$MONGO_USER" \
  --password "$MONGO_PASS" \
  --authenticationDatabase admin \
  --eval "
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
    print('');
  });
  "

echo ""
echo "✅ Replica Set đã sẵn sàng!"
