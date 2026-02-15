#!/bin/bash
set -e

# Load environment variables
if [ -f .env ]; then
    export $(cat .env | grep -v '^#' | xargs)
else
    echo "❌ File .env không tồn tại!"
    exit 1
fi

BACKUP_DIR="./mongodb-backup"
TIMESTAMP=$(date +%Y%m%d-%H%M%S)
BACKUP_PATH="${BACKUP_DIR}/backup-${TIMESTAMP}"

echo "🗄️  MongoDB Backup Script"
echo "━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━"
echo ""
echo "📁 Backup location: ${BACKUP_PATH}"
echo ""

# Create backup directory
mkdir -p ${BACKUP_DIR}

echo "🔄 Đang backup MongoDB..."
docker exec cafe-pos-mongodb mongodump \
  --username ${MONGO_INITDB_ROOT_USERNAME} \
  --password ${MONGO_INITDB_ROOT_PASSWORD} \
  --authenticationDatabase admin \
  --out /data/backup-${TIMESTAMP}

echo ""
echo "📦 Copy backup ra ngoài container..."
docker cp cafe-pos-mongodb:/data/backup-${TIMESTAMP} ${BACKUP_PATH}

echo ""
echo "🧹 Dọn dẹp backup trong container..."
docker exec cafe-pos-mongodb rm -rf /data/backup-${TIMESTAMP}

echo ""
echo "✅ Backup hoàn tất!"
echo ""
echo "━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━"
echo "📊 Thông tin backup:"
echo "━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━"
echo ""
echo "Location: ${BACKUP_PATH}"
echo "Size: $(du -sh ${BACKUP_PATH} | cut -f1)"
echo ""
echo "Collections:"
ls -lh ${BACKUP_PATH}/${MONGODB_DATABASE:-cafe_pos}/ 2>/dev/null || echo "  (Không tìm thấy collections)"
echo ""
echo "━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━"
echo "💡 Để restore backup này:"
echo "━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━"
echo ""
echo "docker cp ${BACKUP_PATH} cafe-pos-mongodb:/data/"
echo "docker exec cafe-pos-mongodb mongorestore \\"
echo "  --username ${MONGO_INITDB_ROOT_USERNAME} \\"
echo "  --password ${MONGO_INITDB_ROOT_PASSWORD} \\"
echo "  --authenticationDatabase admin \\"
echo "  --drop \\"
echo "  /data/backup-${TIMESTAMP}"
echo ""
