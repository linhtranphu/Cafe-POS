#!/bin/bash
set -e

# Load environment variables from .env.production
if [ -f .env.production ]; then
    export $(cat .env.production | grep -v '^#' | xargs)
elif [ -f .env ]; then
    export $(cat .env | grep -v '^#' | xargs)
else
    echo "❌ File .env.production hoặc .env không tồn tại!"
    echo "   Tạo file .env.production với nội dung:"
    echo "   MONGO_INITDB_ROOT_USERNAME=admin"
    echo "   MONGO_INITDB_ROOT_PASSWORD=your_password"
    exit 1
fi

# Configuration
CONTAINER_NAME="cafe-pos-mongodb"
MONGO_USER="${MONGO_INITDB_ROOT_USERNAME}"
MONGO_PASS="${MONGO_INITDB_ROOT_PASSWORD}"
BACKUP_DIR="/home/ubuntu/mongodb-backups"

echo "🔄 MongoDB Replica Set Restore"
echo "━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━"
echo ""

# List available backups
echo "📋 Danh sách backup có sẵn:"
echo ""
ls -lht ${BACKUP_DIR}/cafe_pos_backup_*.gz 2>/dev/null | head -10 || {
    echo "❌ Không tìm thấy backup nào trong ${BACKUP_DIR}"
    exit 1
}

echo ""
echo "━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━"
read -p "Nhập tên file backup (hoặc 'latest' cho backup mới nhất): " BACKUP_INPUT

if [ "$BACKUP_INPUT" == "latest" ]; then
    BACKUP_FILE=$(ls -t ${BACKUP_DIR}/cafe_pos_backup_*.gz 2>/dev/null | head -1)
    if [ -z "$BACKUP_FILE" ]; then
        echo "❌ Không tìm thấy backup nào!"
        exit 1
    fi
    echo "📁 Sử dụng backup mới nhất: $(basename $BACKUP_FILE)"
else
    BACKUP_FILE="${BACKUP_DIR}/${BACKUP_INPUT}"
    if [ ! -f "$BACKUP_FILE" ]; then
        echo "❌ File không tồn tại: ${BACKUP_FILE}"
        exit 1
    fi
fi

echo ""
echo "⚠️  CẢNH BÁO: Restore sẽ XÓA toàn bộ dữ liệu hiện tại!"
echo "   Backup file: $(basename $BACKUP_FILE)"
echo "   Size: $(du -sh $BACKUP_FILE | cut -f1)"
echo ""
read -p "Bạn có chắc chắn muốn restore? (yes/no): " CONFIRM

if [ "$CONFIRM" != "yes" ]; then
    echo "❌ Hủy restore"
    exit 0
fi

# Check if container is running
if ! docker ps | grep -q ${CONTAINER_NAME}; then
    echo "❌ Container ${CONTAINER_NAME} không chạy!"
    exit 1
fi

# Check replica set status
echo ""
echo "🔍 Kiểm tra Replica Set status..."
RS_STATUS=$(docker exec ${CONTAINER_NAME} mongosh \
  --username ${MONGO_USER} \
  --password ${MONGO_PASS} \
  --authenticationDatabase admin \
  --quiet \
  --eval "
  try {
    var status = rs.status();
    if (status.ok === 1) {
      var primary = status.members.find(m => m.stateStr === 'PRIMARY');
      print('OK:' + primary.stateStr);
    } else {
      print('ERROR:Replica set not healthy');
    }
  } catch(e) {
    print('ERROR:' + e);
  }
  " 2>&1)

if [[ $RS_STATUS == OK:* ]]; then
    echo "✅ Replica Set healthy"
else
    echo "❌ Replica Set không healthy: ${RS_STATUS}"
    echo "   Không thể restore khi replica set không ổn định!"
    exit 1
fi

echo ""
echo "🔄 Đang restore database..."
echo "   (Quá trình này có thể mất vài phút...)"
echo ""

# Perform restore
cat ${BACKUP_FILE} | docker exec -i ${CONTAINER_NAME} mongorestore \
  --username ${MONGO_USER} \
  --password ${MONGO_PASS} \
  --authenticationDatabase admin \
  --oplogReplay \
  --drop \
  --archive --gzip

if [ $? -eq 0 ]; then
    echo ""
    echo "✅ Restore thành công!"
else
    echo ""
    echo "❌ Restore thất bại!"
    exit 1
fi

echo ""
echo "━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━"
echo "📊 Verify dữ liệu sau restore:"
echo "━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━"
echo ""

docker exec ${CONTAINER_NAME} mongosh \
  --username ${MONGO_USER} \
  --password ${MONGO_PASS} \
  --authenticationDatabase admin \
  --eval "
  var db = db.getSiblingDB('cafe_pos');
  var collections = db.getCollectionNames();
  print('Database: cafe_pos');
  print('Total Collections: ' + collections.length);
  print('');
  collections.forEach(function(col) {
    var count = db.getCollection(col).countDocuments();
    print('  - ' + col + ': ' + count + ' documents');
  });
  "

echo ""
echo "✅ Restore hoàn tất!"
echo ""
echo "💡 Khởi động lại backend để apply changes:"
echo "   docker-compose -f docker-compose.prod.yml restart backend"
