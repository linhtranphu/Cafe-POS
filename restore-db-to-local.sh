#!/bin/bash

# Restore MongoDB backup to local database
# Usage: ./scripts/restore-db-to-local.sh <backup-directory>

set -e

echo "🔄 Restore MongoDB Backup to Local"
echo "==================================="
echo ""

# Check arguments
if [ $# -eq 0 ]; then
  echo "❌ Error: Backup directory not specified"
  echo ""
  echo "Usage:"
  echo "  ./scripts/restore-db-to-local.sh <backup-directory>"
  echo ""
  echo "Example:"
  echo "  ./scripts/restore-db-to-local.sh ./backups/cafe_pos_backup_20260205_150000"
  exit 1
fi

BACKUP_DIR="$1"
DB_NAME="cafe_pos"

# Check if backup directory exists
if [ ! -d "$BACKUP_DIR" ]; then
  echo "❌ Error: Backup directory not found: $BACKUP_DIR"
  exit 1
fi

# Check if cafe_pos directory exists in backup
if [ ! -d "$BACKUP_DIR/cafe_pos" ]; then
  echo "❌ Error: cafe_pos database not found in backup"
  echo "Expected: $BACKUP_DIR/cafe_pos"
  exit 1
fi

echo "📋 Configuration:"
echo "  Backup: $BACKUP_DIR"
echo "  Database: $DB_NAME"
echo ""

# Warning
echo "⚠️  WARNING: This will replace your local database!"
echo ""
read -p "Continue? (y/N): " -n 1 -r
echo ""

if [[ ! $REPLY =~ ^[Yy]$ ]]; then
  echo "❌ Restore cancelled"
  exit 1
fi

echo ""
echo "🔄 Restoring database..."

# Check if MongoDB is running in Docker
if docker ps | grep -q mongodb; then
  echo "  📦 Using Docker MongoDB..."
  
  # Copy backup to container
  docker cp "$BACKUP_DIR/cafe_pos" mongodb:/tmp/
  
  # Restore in container
  docker exec mongodb mongorestore \
    --db cafe_pos \
    --drop \
    /tmp/cafe_pos
  
  # Cleanup
  docker exec mongodb rm -rf /tmp/cafe_pos
  
elif command -v mongorestore &> /dev/null; then
  echo "  💻 Using local MongoDB..."
  
  # Restore to local MongoDB
  mongorestore \
    --db cafe_pos \
    --drop \
    "$BACKUP_DIR/cafe_pos"
else
  echo "❌ Error: MongoDB not found"
  echo "Please ensure MongoDB is running (Docker or local)"
  exit 1
fi

echo ""
echo "✅ Restore complete!"
echo ""
echo "🔍 Verify restore:"
echo "  docker exec -it mongodb mongosh cafe_pos --eval 'db.stats()'"
