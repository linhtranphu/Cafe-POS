#!/bin/bash

# Backup MongoDB from EC2 to Local
# This script connects to EC2, dumps MongoDB, and downloads to local machine

set -e

echo "🔄 MongoDB Backup from EC2 to Local"
echo "===================================="
echo ""

# Configuration
EC2_USER="${EC2_USER:-ubuntu}"
EC2_HOST="${EC2_HOST:-13.212.27.222}"
EC2_KEY="${EC2_KEY:-/Volumes/MacOS/users/tranphulinh/EC2PEM/OngTaPOS.pem}"
MONGO_CONTAINER="mongodb"
DB_NAME="cafe_pos"
BACKUP_DIR="./backups"
TIMESTAMP=$(date +%Y%m%d_%H%M%S)
BACKUP_NAME="cafe_pos_backup_${TIMESTAMP}"

# Colors
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
RED='\033[0;31m'
NC='\033[0m' # No Color

# Check if EC2_HOST is set
if [ "$EC2_HOST" = "your-ec2-ip" ]; then
  echo -e "${RED}❌ Error: EC2_HOST not set${NC}"
  echo ""
  echo "Usage:"
  echo "  EC2_HOST=13.212.27.222 EC2_USER=ubuntu EC2_KEY=/path/to/key.pem ./scripts/backup-db-from-ec2.sh"
  echo ""
  echo "Or use defaults:"
  echo "  ./scripts/backup-db-from-ec2.sh"
  echo ""
  echo "Default values:"
  echo "  EC2_HOST=13.212.27.222"
  echo "  EC2_USER=ubuntu"
  echo "  EC2_KEY=/Volumes/MacOS/users/tranphulinh/EC2PEM/OngTaPOS.pem"
  exit 1
fi

# Create backup directory
mkdir -p "$BACKUP_DIR"

echo -e "${YELLOW}📋 Configuration:${NC}"
echo "  EC2 Host: $EC2_HOST"
echo "  EC2 User: $EC2_USER"
echo "  EC2 Key: $EC2_KEY"
echo "  Database: $DB_NAME"
echo "  Backup Name: $BACKUP_NAME"
echo ""

# Step 1: Create backup on EC2
echo -e "${YELLOW}1️⃣ Creating backup on EC2...${NC}"
ssh -i "$EC2_KEY" "$EC2_USER@$EC2_HOST" << 'ENDSSH'
  set -e
  
  # Create backup directory on EC2
  mkdir -p ~/backups
  
  # Dump MongoDB from container
  echo "  📦 Dumping MongoDB..."
  docker exec mongodb mongodump \
    --db cafe_pos \
    --out /tmp/backup
  
  # Copy backup from container to EC2 host
  echo "  📋 Copying from container..."
  docker cp mongodb:/tmp/backup ~/backups/
  
  # Create tarball
  echo "  🗜️ Creating tarball..."
  cd ~/backups
  tar -czf cafe_pos_backup.tar.gz backup/
  
  # Cleanup
  rm -rf backup/
  
  echo "  ✅ Backup created on EC2"
ENDSSH

if [ $? -ne 0 ]; then
  echo -e "${RED}❌ Failed to create backup on EC2${NC}"
  exit 1
fi

echo -e "${GREEN}✅ Backup created on EC2${NC}"
echo ""

# Step 2: Download backup from EC2
echo -e "${YELLOW}2️⃣ Downloading backup from EC2...${NC}"
scp -i "$EC2_KEY" \
  "$EC2_USER@$EC2_HOST:~/backups/cafe_pos_backup.tar.gz" \
  "$BACKUP_DIR/${BACKUP_NAME}.tar.gz"

if [ $? -ne 0 ]; then
  echo -e "${RED}❌ Failed to download backup${NC}"
  exit 1
fi

echo -e "${GREEN}✅ Backup downloaded${NC}"
echo ""

# Step 3: Extract backup
echo -e "${YELLOW}3️⃣ Extracting backup...${NC}"
cd "$BACKUP_DIR"
tar -xzf "${BACKUP_NAME}.tar.gz"
mv backup "$BACKUP_NAME"

echo -e "${GREEN}✅ Backup extracted${NC}"
echo ""

# Step 4: Cleanup EC2
echo -e "${YELLOW}4️⃣ Cleaning up EC2...${NC}"
ssh -i "$EC2_KEY" "$EC2_USER@$EC2_HOST" << 'ENDSSH'
  rm -f ~/backups/cafe_pos_backup.tar.gz
  echo "  ✅ Cleanup complete"
ENDSSH

echo -e "${GREEN}✅ EC2 cleanup complete${NC}"
echo ""

# Summary
echo -e "${GREEN}🎉 Backup Complete!${NC}"
echo ""
echo "📁 Backup Location:"
echo "  Directory: $BACKUP_DIR/$BACKUP_NAME/"
echo "  Archive: $BACKUP_DIR/${BACKUP_NAME}.tar.gz"
echo ""
echo "📊 Backup Info:"
du -sh "$BACKUP_DIR/$BACKUP_NAME"
echo ""
echo "🔄 To restore this backup locally:"
echo "  mongorestore --db cafe_pos $BACKUP_DIR/$BACKUP_NAME/cafe_pos"
echo ""
echo "🔄 To restore to EC2:"
echo "  ./scripts/restore-db-to-ec2.sh $BACKUP_DIR/$BACKUP_NAME"
