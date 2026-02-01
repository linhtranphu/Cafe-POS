#!/bin/bash

# Colors for output
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
NC='\033[0m' # No Color

echo -e "${YELLOW}🔄 Stopping existing processes on port 3000...${NC}"
lsof -ti:3000 | xargs -r kill -9 2>/dev/null
pkill -f "cafe-pos-server" 2>/dev/null
pkill -f "go run main.go" 2>/dev/null
sleep 2

echo -e "${GREEN}✅ Ports cleared${NC}"
echo ""

# Check if MongoDB is running
echo -e "${YELLOW}🔍 Checking MongoDB...${NC}"
if ! lsof -i :27017 > /dev/null 2>&1; then
    echo -e "${RED}❌ MongoDB is not running on port 27017${NC}"
    echo -e "${YELLOW}Starting MongoDB with Docker...${NC}"
    docker run -d -p 27017:27017 \
        -e MONGO_INITDB_ROOT_USERNAME=admin \
        -e MONGO_INITDB_ROOT_PASSWORD=password123 \
        --name cafe-pos-mongodb \
        mongo:7.0 2>/dev/null || true
    echo -e "${YELLOW}Waiting for MongoDB to be ready (10 seconds)...${NC}"
    sleep 10
fi
echo -e "${GREEN}✅ MongoDB is running${NC}"
echo ""

# Start backend
echo -e "${YELLOW}🚀 Starting backend on port 3000...${NC}"
cd backend
MONGODB_URI="mongodb://admin:password123@localhost:27017" go run main.go
