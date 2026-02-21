#!/bin/bash

echo "🔄 Switching MongoDB to Replica Set Mode"
echo "━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━"
echo ""

# Colors
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
RED='\033[0;31m'
NC='\033[0m'

# Check if Docker is running
if ! docker info > /dev/null 2>&1; then
    echo -e "${RED}❌ Docker is not running. Please start Docker first.${NC}"
    exit 1
fi

echo -e "${YELLOW}📋 Step 1: Stopping current containers...${NC}"
docker-compose down

echo ""
echo -e "${YELLOW}📋 Step 2: Starting MongoDB with replica set configuration...${NC}"
docker-compose -f docker-compose.replica-set.yml up -d mongodb

echo ""
echo -e "${YELLOW}⏳ Step 3: Waiting for MongoDB to start (20 seconds)...${NC}"
sleep 20

echo ""
echo -e "${YELLOW}📋 Step 4: Initializing replica set...${NC}"
docker exec cafe-pos-mongodb mongosh \
  --username admin \
  --password password123 \
  --authenticationDatabase admin \
  --eval "
  try {
    var status = rs.status();
    print('✅ Replica Set already initialized: ' + status.set);
  } catch(e) {
    print('🔧 Initializing Replica Set...');
    var result = rs.initiate({
      _id: 'rs0',
      members: [
        { _id: 0, host: 'localhost:27017' }
      ]
    });
    print('✅ Replica Set initialized');
  }
  "

echo ""
echo -e "${YELLOW}⏳ Step 5: Waiting for replica set to stabilize (15 seconds)...${NC}"
sleep 15

echo ""
echo -e "${YELLOW}📋 Step 6: Verifying replica set status...${NC}"
docker exec cafe-pos-mongodb mongosh \
  --username admin \
  --password password123 \
  --authenticationDatabase admin \
  --eval "
  var status = rs.status();
  if (status.ok === 1) {
    print('✅ Replica Set Status: OK');
    print('   Set Name: ' + status.set);
    var primary = status.members.find(m => m.stateStr === 'PRIMARY');
    if (primary) {
      print('   Primary: ' + primary.name + ' (' + primary.stateStr + ')');
    }
  } else {
    print('❌ Replica Set not ready');
    quit(1);
  }
  "

if [ $? -eq 0 ]; then
    echo ""
    echo -e "${GREEN}✅ MongoDB Replica Set is ready!${NC}"
    echo ""
    echo -e "${GREEN}Next steps:${NC}"
    echo -e "  1. Start backend: ${YELLOW}docker-compose -f docker-compose.replica-set.yml up -d backend${NC}"
    echo -e "  2. Start frontend: ${YELLOW}docker-compose -f docker-compose.replica-set.yml up -d frontend${NC}"
    echo -e "  3. Or start all: ${YELLOW}docker-compose -f docker-compose.replica-set.yml up -d${NC}"
    echo ""
    echo -e "${GREEN}Connection string:${NC}"
    echo -e "  ${YELLOW}mongodb://admin:password123@localhost:27017/cafe_pos?replicaSet=rs0&authSource=admin${NC}"
else
    echo ""
    echo -e "${RED}❌ Setup failed. Check logs:${NC}"
    echo -e "  ${YELLOW}docker logs cafe-pos-mongodb${NC}"
    exit 1
fi
