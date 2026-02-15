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

# Configuration
CONTAINER_NAME="cafe-pos-mongodb"
MONGO_USER="${MONGO_INITDB_ROOT_USERNAME}"
MONGO_PASS="${MONGO_INITDB_ROOT_PASSWORD}"
BACKUP_DIR="/home/ubuntu/mongodb-backups"
TIMESTAMP=$(date +%Y%m%d-%H%M%S)
BACKUP_NAME="cafe_pos_backup_${TIMESTAMP}.gz"
BACKUP_PATH="${BACKUP_DIR}/${BACKUP_NAME}"

# Retention (số ngày giữ backup)
RETENTION_DAYS=7

echo "🗄️  MongoDB Replica Set Backup"
echo "━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━"
echo ""
echo "📅 Timestamp: $(date '+%Y-%m-%d %H:%M:%S')"
echo "📁 Backup file: ${BACKUP_NAME}"
echo ""

# Create backup directory if not exists
mkdir -p ${BACKUP_DIR}

# Check if container is running
if ! docker ps | grep -q ${CONTAINER_NAME}; then
    echo "❌ Container ${CONTAINER_NAME} không chạy!"
    exit 1
fi

# Check replica set status
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
      print('OK:' + primary.name + ':' + primary.stateStr);
    } else {
      print('ERROR:Replica set not healthy');
    }
  } catch(e) {
    print('ERROR:' + e);
  }
  " 2>&1)

if [[ $RS_STATUS == OK:* ]]; then
    echo "✅ Replica Set healthy: ${RS_STATUS#OK:}"
else
    echo "⚠️  Warning: ${RS_STATUS}"
    echo "   Backup sẽ tiếp tục nhưng có thể không consistent..."
fi

echo ""
echo "🔄 Đang backup database..."

# Perform backup with readPreference=primary to ensure consistency
# --oplog: Include oplog for point-in-time consistency
docker exec ${CONTAINER_NAME} mongodump \
  --username ${MONGO_USER} \
  --password ${MONGO_PASS} \
  --authenticationDatabase admin \
  --readPreference=primary \
  --oplog \
  --archive --gzip > ${BACKUP_PATH}

if [ $? -eq 0 ]; then
    echo "✅ Backup thành công!"
else
    echo "❌ Backup thất bại!"
    exit 1
fi

echo ""
echo "━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━"
echo "📊 Thông tin backup:"
echo "━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━"
echo ""
echo "Location: ${BACKUP_PATH}"
echo "Size: $(du -sh ${BACKUP_PATH} | cut -f1)"
echo "MD5: $(md5sum ${BACKUP_PATH} | cut -d' ' -f1)"
echo ""

# Cleanup old backups
echo "🧹 Dọn dẹp backup cũ (giữ ${RETENTION_DAYS} ngày)..."
find ${BACKUP_DIR} -name "cafe_pos_backup_*.gz" -type f -mtime +${RETENTION_DAYS} -delete
echo "✅ Đã xóa backup cũ hơn ${RETENTION_DAYS} ngày"

echo ""
echo "━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━"
echo "📝 Danh sách backup hiện có:"
echo "━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━"
ls -lh ${BACKUP_DIR}/cafe_pos_backup_*.gz 2>/dev/null | tail -5 || echo "  (Không có backup)"

echo ""
echo "━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━"
echo "💡 Để restore backup này:"
echo "━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━"
echo ""
echo "docker exec -i ${CONTAINER_NAME} mongorestore \\"
echo "  --username ${MONGO_USER} \\"
echo "  --password ${MONGO_PASS} \\"
echo "  --authenticationDatabase admin \\"
echo "  --oplogReplay \\"
echo "  --drop \\"
echo "  --archive --gzip < ${BACKUP_PATH}"
echo ""
echo "✅ Backup hoàn tất!"
