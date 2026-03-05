#!/bin/bash

# Fix print jobs timeout issue on EC2
# This script creates missing collections and indexes

set -e

echo "=========================================="
echo "🔧 Fix Print Jobs Timeout Issue"
echo "=========================================="
echo ""

# Configuration
EC2_USER="${EC2_USER:-ubuntu}"
EC2_HOST="${EC2_HOST:-47.128.65.142}"
EC2_KEY="${EC2_KEY:-/Volumes/Linh-DAT/TaCafePOS/TaCafePOS.pem}"

# Colors
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
BLUE='\033[0;34m'
NC='\033[0m' # No Color

echo -e "${BLUE}📋 Configuration:${NC}"
echo "  EC2 Host: $EC2_HOST"
echo "  EC2 User: $EC2_USER"
echo ""

# Test SSH connection
echo -e "${YELLOW}🔍 Testing SSH connection...${NC}"
if ! ssh -i "$EC2_KEY" -o ConnectTimeout=10 -o StrictHostKeyChecking=no "$EC2_USER@$EC2_HOST" "echo 'Connection successful'" > /dev/null 2>&1; then
    echo -e "${RED}❌ Failed to connect to EC2${NC}"
    exit 1
fi
echo -e "${GREEN}✅ SSH connection successful${NC}"
echo ""

# Run fix on EC2
echo -e "${YELLOW}🔧 Running fix on EC2...${NC}"
ssh -i "$EC2_KEY" "$EC2_USER@$EC2_HOST" bash << 'ENDSSH'
  set -e
  
  MONGO_CONTAINER="cafe-pos-mongodb"
  DATABASE_NAME="cafe_pos"
  
  echo "  📦 Checking MongoDB container..."
  if ! docker ps | grep -q "$MONGO_CONTAINER"; then
    echo "  ❌ MongoDB container '$MONGO_CONTAINER' is not running"
    exit 1
  fi
  echo "  ✅ MongoDB container is running"
  echo ""
  
  # Get MongoDB password
  MONGO_PASSWORD=$(docker exec $MONGO_CONTAINER printenv MONGO_INITDB_ROOT_PASSWORD)
  
  echo "  🔧 Creating indexes for print_jobs..."
  docker exec $MONGO_CONTAINER mongosh $DATABASE_NAME \
    --username admin \
    --password "$MONGO_PASSWORD" \
    --authenticationDatabase admin \
    --quiet --eval '
    // Create indexes for print_jobs
    db.print_jobs.createIndex({ "status": 1, "created_at": -1 });
    db.print_jobs.createIndex({ "order_id": 1 });
    db.print_jobs.createIndex({ "printer_id": 1 });
    print("  ✅ Created print_jobs indexes");
  '
  
  echo ""
  echo "  🔧 Creating indexes for printer_configs..."
  docker exec $MONGO_CONTAINER mongosh $DATABASE_NAME \
    --username admin \
    --password "$MONGO_PASSWORD" \
    --authenticationDatabase admin \
    --quiet --eval '
    db.printer_configs.createIndex({ "type": 1, "is_default": 1 });
    print("  ✅ Created printer_configs indexes");
  '
  
  echo ""
  echo "  🔧 Creating indexes for print_templates..."
  docker exec $MONGO_CONTAINER mongosh $DATABASE_NAME \
    --username admin \
    --password "$MONGO_PASSWORD" \
    --authenticationDatabase admin \
    --quiet --eval '
    db.print_templates.createIndex({ "type": 1, "is_default": 1 });
    print("  ✅ Created print_templates indexes");
  '
  
  echo ""
  echo "  📊 Checking collections..."
  docker exec $MONGO_CONTAINER mongosh $DATABASE_NAME \
    --username admin \
    --password "$MONGO_PASSWORD" \
    --authenticationDatabase admin \
    --quiet --eval '
    var printJobs = db.print_jobs.countDocuments();
    var printerConfigs = db.printer_configs.countDocuments();
    var printTemplates = db.print_templates.countDocuments();
    print("  - print_jobs: " + printJobs + " documents");
    print("  - printer_configs: " + printerConfigs + " documents");
    print("  - print_templates: " + printTemplates + " documents");
  '
  
  echo ""
  echo "  🔄 Restarting backend..."
  if docker ps | grep -q "backend"; then
    docker restart backend > /dev/null 2>&1
    sleep 3
    echo "  ✅ Backend restarted"
  else
    echo "  ⚠️  Backend container not found"
  fi
  
  echo ""
  echo "  ✅ Fix completed successfully"
ENDSSH

if [ $? -ne 0 ]; then
    echo -e "${RED}❌ Failed to run fix${NC}"
    exit 1
fi

echo ""
echo -e "${GREEN}=========================================="
echo "✅ Fix Completed Successfully!"
echo "==========================================${NC}"
echo ""
echo -e "${BLUE}📝 What was fixed:${NC}"
echo "  ✅ Created indexes for print_jobs"
echo "  ✅ Created indexes for printer_configs"
echo "  ✅ Created indexes for print_templates"
echo "  ✅ Restarted backend"
echo ""
echo -e "${BLUE}🌐 Please refresh your browser:${NC}"
echo "   https://tacafe.store/#/print-management"
echo ""
