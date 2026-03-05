#!/bin/bash

# Interactive MongoDB Restore to EC2
# Usage: ./restore-to-ec2-interactive.sh

set -e

# Colors for output
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
BLUE='\033[0;34m'
CYAN='\033[0;36m'
NC='\033[0m' # No Color

clear
echo -e "${BLUE}╔════════════════════════════════════════════════════════╗${NC}"
echo -e "${BLUE}║                                                        ║${NC}"
echo -e "${BLUE}║        MongoDB Database Restore to EC2                 ║${NC}"
echo -e "${BLUE}║                                                        ║${NC}"
echo -e "${BLUE}╚════════════════════════════════════════════════════════╝${NC}"
echo ""

# Step 1: Get EC2 Host
echo -e "${CYAN}📍 Step 1: EC2 Configuration${NC}"
echo -e "${YELLOW}────────────────────────────────────────────────────────${NC}"
echo ""
read -p "Enter EC2 Host (IP or domain): " EC2_HOST

if [ -z "$EC2_HOST" ]; then
    echo -e "${RED}❌ Error: EC2 Host is required${NC}"
    exit 1
fi

echo -e "${GREEN}✓ EC2 Host: $EC2_HOST${NC}"
echo ""

# Step 2: Get EC2 User
read -p "Enter EC2 User [default: ubuntu]: " EC2_USER
EC2_USER=${EC2_USER:-ubuntu}
echo -e "${GREEN}✓ EC2 User: $EC2_USER${NC}"
echo ""

# Step 3: Check PEM file
echo -e "${CYAN}🔑 Step 2: PEM Key File${NC}"
echo -e "${YELLOW}────────────────────────────────────────────────────────${NC}"
echo ""

# Get script directory
SCRIPT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
EC2_KEY="$SCRIPT_DIR/TaCafePOS.pem"

if [ ! -f "$EC2_KEY" ]; then
    echo -e "${RED}❌ Error: PEM file not found: $EC2_KEY${NC}"
    echo -e "${YELLOW}   Expected location: $SCRIPT_DIR/TaCafePOS.pem${NC}"
    exit 1
fi

echo -e "${GREEN}✓ PEM file found: $EC2_KEY${NC}"
echo ""

# Step 4: Get backup file
echo -e "${CYAN}📦 Step 3: Select Backup File${NC}"
echo -e "${YELLOW}────────────────────────────────────────────────────────${NC}"
echo ""

# Use script directory for backup location
BACKUP_DIR="$SCRIPT_DIR/backups_from_ec2"

if [ ! -d "$BACKUP_DIR" ]; then
    echo -e "${RED}❌ Error: Backup directory not found: $BACKUP_DIR${NC}"
    exit 1
fi

echo -e "${BLUE}Available backup files in: ${CYAN}$BACKUP_DIR${NC}"
echo ""

# List available backup files
ls -lh "$BACKUP_DIR"/*.gz 2>/dev/null | awk '{printf "  %s  %s\n", $5, $9}' | while read size file; do
    filename=$(basename "$file")
    echo -e "  ${GREEN}•${NC} $filename ${CYAN}($size)${NC}"
done

echo ""
read -p "Enter backup filename (e.g., mongodb-backup-2026-02-28_12-36-31.gz): " BACKUP_FILENAME

if [ -z "$BACKUP_FILENAME" ]; then
    echo -e "${RED}❌ Error: Backup filename is required${NC}"
    exit 1
fi

BACKUP_FILE="$BACKUP_DIR/$BACKUP_FILENAME"

if [ ! -f "$BACKUP_FILE" ]; then
    echo -e "${RED}❌ Error: Backup file not found: $BACKUP_FILE${NC}"
    exit 1
fi

echo -e "${GREEN}✓ Selected: $BACKUP_FILENAME${NC}"
echo ""

# Step 5: Get database name
echo -e "${CYAN}🗄️  Step 4: Database Configuration${NC}"
echo -e "${YELLOW}────────────────────────────────────────────────────────${NC}"
echo ""
read -p "Enter database name [default: cafe_pos]: " DATABASE_NAME
DATABASE_NAME=${DATABASE_NAME:-cafe_pos}
echo -e "${GREEN}✓ Database: $DATABASE_NAME${NC}"
echo ""

# Step 6: Get MongoDB container name
read -p "Enter MongoDB container name [default: cafe-pos-mongodb]: " MONGO_CONTAINER
MONGO_CONTAINER=${MONGO_CONTAINER:-cafe-pos-mongodb}
echo -e "${GREEN}✓ Container: $MONGO_CONTAINER${NC}"
echo ""

# Display summary
echo ""
echo -e "${BLUE}╔════════════════════════════════════════════════════════╗${NC}"
echo -e "${BLUE}║                    Configuration Summary               ║${NC}"
echo -e "${BLUE}╚════════════════════════════════════════════════════════╝${NC}"
echo ""
echo -e "${CYAN}EC2 Host:${NC}       $EC2_HOST"
echo -e "${CYAN}EC2 User:${NC}       $EC2_USER"
echo -e "${CYAN}PEM Key:${NC}        $EC2_KEY"
echo -e "${CYAN}Database:${NC}       $DATABASE_NAME"
echo -e "${CYAN}Container:${NC}      $MONGO_CONTAINER"
echo -e "${CYAN}Backup File:${NC}    $(basename "$BACKUP_FILE")"
echo -e "${CYAN}Backup Size:${NC}    $(ls -lh "$BACKUP_FILE" | awk '{print $5}')"
echo ""
echo -e "${RED}⚠️  WARNING: This will REPLACE the database on EC2!${NC}"
echo -e "${RED}   All current data will be LOST!${NC}"
echo ""
read -p "Are you sure you want to continue? (yes/no): " -r
echo ""

if [[ ! $REPLY =~ ^[Yy][Ee][Ss]$ ]]; then
    echo -e "${YELLOW}Restore cancelled.${NC}"
    exit 0
fi

# Test SSH connection first
echo -e "${YELLOW}🔍 Testing SSH connection...${NC}"
if ! ssh -i "$EC2_KEY" -o ConnectTimeout=10 -o StrictHostKeyChecking=no "$EC2_USER@$EC2_HOST" "echo 'Connection successful'" > /dev/null 2>&1; then
    echo -e "${RED}❌ Failed to connect to EC2${NC}"
    echo -e "${RED}   Please check your EC2 host, user, and PEM key${NC}"
    exit 1
fi
echo -e "${GREEN}✅ SSH connection successful${NC}"
echo ""

# Start restore process
echo -e "${BLUE}╔════════════════════════════════════════════════════════╗${NC}"
echo -e "${BLUE}║                  Starting Restore Process              ║${NC}"
echo -e "${BLUE}╚════════════════════════════════════════════════════════╝${NC}"
echo ""

# Step 1: Upload backup to EC2
echo -e "${YELLOW}1️⃣  Uploading backup to EC2...${NC}"
ssh -i "$EC2_KEY" "$EC2_USER@$EC2_HOST" "mkdir -p /home/$EC2_USER/backups"

BACKUP_FILENAME=$(basename "$BACKUP_FILE")
scp -i "$EC2_KEY" -o StrictHostKeyChecking=no \
    "$BACKUP_FILE" \
    "$EC2_USER@$EC2_HOST:/home/$EC2_USER/backups/$BACKUP_FILENAME"

if [ $? -ne 0 ]; then
    echo -e "${RED}❌ Failed to upload backup${NC}"
    exit 1
fi

echo -e "${GREEN}✅ Backup uploaded successfully${NC}"
echo ""

# Step 2: Restore database on EC2
echo -e "${YELLOW}2️⃣  Restoring database on EC2...${NC}"
ssh -i "$EC2_KEY" "$EC2_USER@$EC2_HOST" bash << ENDSSH
  set -e
  
  MONGO_CONTAINER="$MONGO_CONTAINER"
  DATABASE_NAME="$DATABASE_NAME"
  BACKUP_FILE="/home/$EC2_USER/backups/$BACKUP_FILENAME"
  
  echo "  📦 Checking MongoDB container..."
  if ! docker ps | grep -q "\$MONGO_CONTAINER"; then
    echo "  ❌ MongoDB container '\$MONGO_CONTAINER' is not running"
    exit 1
  fi
  echo "  ✅ MongoDB container is running"
  
  echo "  📦 Copying backup into MongoDB container..."
  docker cp "\$BACKUP_FILE" \$MONGO_CONTAINER:/tmp/restore-backup.gz
  
  echo "  🗑️  Dropping existing database..."
  MONGO_PASSWORD=\$(docker exec \$MONGO_CONTAINER printenv MONGO_INITDB_ROOT_PASSWORD)
  
  docker exec \$MONGO_CONTAINER mongosh \
    --username admin \
    --password "\$MONGO_PASSWORD" \
    --authenticationDatabase admin \
    \$DATABASE_NAME \
    --eval "db.dropDatabase()" 2>/dev/null || true
  
  echo "  📥 Restoring database from backup..."
  docker exec \$MONGO_CONTAINER mongorestore \
    --username admin \
    --password "\$MONGO_PASSWORD" \
    --authenticationDatabase admin \
    --db \$DATABASE_NAME \
    --archive=/tmp/restore-backup.gz \
    --gzip
  
  echo "  🧹 Cleaning up temporary files in container..."
  docker exec \$MONGO_CONTAINER rm -f /tmp/restore-backup.gz
  
  echo ""
  echo "  📊 Database statistics:"
  docker exec \$MONGO_CONTAINER mongosh \
    --username admin \
    --password "\$MONGO_PASSWORD" \
    --authenticationDatabase admin \
    \$DATABASE_NAME \
    --quiet \
    --eval "
      print('  Collections:');
      db.getCollectionNames().forEach(function(col) {
        var count = db[col].countDocuments();
        print('    - ' + col + ': ' + count + ' documents');
      });
    "
  
  echo ""
  echo "  ✅ Database restored successfully"
ENDSSH

if [ $? -ne 0 ]; then
    echo -e "${RED}❌ Failed to restore database${NC}"
    exit 1
fi

echo -e "${GREEN}✅ Database restored on EC2${NC}"
echo ""

# Step 3: Cleanup
echo -e "${YELLOW}3️⃣  Cleaning up...${NC}"
ssh -i "$EC2_KEY" "$EC2_USER@$EC2_HOST" "rm -f /home/$EC2_USER/backups/$BACKUP_FILENAME"
echo -e "${GREEN}✅ Cleanup complete${NC}"
echo ""

# Step 4: Restart backend (optional)
echo -e "${YELLOW}4️⃣  Restarting backend service...${NC}"
ssh -i "$EC2_KEY" "$EC2_USER@$EC2_HOST" bash << 'ENDSSH'
  if docker ps | grep -q "backend"; then
    echo "  🔄 Restarting backend container..."
    docker restart backend > /dev/null 2>&1
    sleep 3
    
    # Check if backend is running
    if docker ps | grep -q "backend"; then
      echo "  ✅ Backend restarted successfully"
    else
      echo "  ⚠️  Backend may have issues, please check logs"
    fi
  else
    echo "  ℹ️  Backend container not found, skipping restart"
  fi
ENDSSH

echo ""
echo -e "${GREEN}╔════════════════════════════════════════════════════════╗${NC}"
echo -e "${GREEN}║                                                        ║${NC}"
echo -e "${GREEN}║              ✅ Restore Completed Successfully!         ║${NC}"
echo -e "${GREEN}║                                                        ║${NC}"
echo -e "${GREEN}╚════════════════════════════════════════════════════════╝${NC}"
echo ""
echo -e "${BLUE}📝 Summary:${NC}"
echo -e "  ✅ Backup uploaded to EC2"
echo -e "  ✅ Database restored: ${CYAN}$DATABASE_NAME${NC}"
echo -e "  ✅ Backend restarted"
echo ""
echo -e "${BLUE}🌐 Your application is now running with the restored data${NC}"
echo -e "   URL: ${CYAN}http://$EC2_HOST${NC}"
echo ""
