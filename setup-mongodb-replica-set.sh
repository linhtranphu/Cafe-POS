#!/bin/bash

echo "🚀 Setting up MongoDB Replica Set for Cafe POS"
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

echo -e "${YELLOW}📋 Step 1: Stopping existing containers...${NC}"
docker-compose down

echo ""
echo -e "${YELLOW}📋 Step 2: Starting MongoDB with replica set...${NC}"
docker-compose -f docker-compose.replica-set.yml up -d mongodb

echo ""
echo -e "${YELLOW}⏳ Step 3: Waiting for MongoDB to initialize (30 seconds)...${NC}"
sleep 30

echo ""
echo -e "${YELLOW}📋 Step 4: Checking replica set status...${NC}"
docker exec cafe-pos-mongodb mongosh --quiet --eval "
var status = rs.status();
if (status.ok === 1) {
    print('✓ Replica set initialized successfully');
    print('  Set name: ' + status.set);
    var primary = status.members.find(m => m.stateStr === 'PRIMARY');
    if (primary) {
        print('  Primary: ' + primary.name);
    }
} else {
    print('✗ Replica set not ready yet');
    print('Attempting manual initialization...');
    rs.initiate({
        _id: 'rs0',
        members: [{_id: 0, host: 'localhost:27017'}]
    });
}
"

echo ""
echo -e "${YELLOW}⏳ Waiting for replica set to stabilize (10 seconds)...${NC}"
sleep 10

echo ""
echo -e "${YELLOW}📋 Step 5: Testing transactions...${NC}"
./test-mongodb-transactions.sh

if [ $? -eq 0 ]; then
    echo ""
    echo -e "${GREEN}✅ MongoDB Replica Set setup complete!${NC}"
    echo ""
    echo -e "${GREEN}Next steps:${NC}"
    echo -e "  1. Start backend: ${YELLOW}docker-compose -f docker-compose.replica-set.yml up -d backend${NC}"
    echo -e "  2. Run tests: ${YELLOW}cd backend && go test -v -run=\"Batch\" ./application/services/...${NC}"
    echo ""
    echo -e "${GREEN}Connection string:${NC}"
    echo -e "  ${YELLOW}mongodb://admin:password@localhost:27017/cafe_pos?replicaSet=rs0&authSource=admin${NC}"
else
    echo ""
    echo -e "${RED}❌ Setup failed. Please check the logs:${NC}"
    echo -e "  ${YELLOW}docker logs cafe-pos-mongodb${NC}"
    echo ""
    echo -e "${RED}For manual setup, see: MONGODB_REPLICA_SET_SETUP.md${NC}"
    exit 1
fi
