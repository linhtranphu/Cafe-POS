#!/bin/bash

set -e

# Colors
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
BLUE='\033[0;34m'
NC='\033[0m'

echo -e "${BLUE}╔════════════════════════════════════════════════════════════╗${NC}"
echo -e "${BLUE}║     🚀 EC2 Deployment - Setup & Deploy                   ║${NC}"
echo -e "${BLUE}╚════════════════════════════════════════════════════════════╝${NC}"
echo ""

# Check if .env exists
if [ ! -f .env ]; then
    echo -e "${RED}❌ Error: .env file not found!${NC}"
    echo "Please create .env file with MongoDB password"
    exit 1
fi

source .env

# Step 1: Install Docker
echo -e "${BLUE}📦 Step 1: Installing Docker...${NC}"
if ! command -v docker &> /dev/null; then
    curl -fsSL https://get.docker.com -o get-docker.sh
    sudo sh get-docker.sh
    sudo usermod -aG docker $USER
    echo -e "${GREEN}✅ Docker installed${NC}"
else
    echo -e "${GREEN}✅ Docker already installed${NC}"
fi
echo ""

# Step 2: Install Docker Compose
echo -e "${BLUE}� Step 2: Installing Docker Compose...${NC}"
if ! command -v docker-compose &> /dev/null; then
    sudo curl -L "https://github.com/docker/compose/releases/latest/download/docker-compose-$(uname -s)-$(uname -m)" -o /usr/local/bin/docker-compose
    sudo chmod +x /usr/local/bin/docker-compose
    echo -e "${GREEN}✅ Docker Compose installed${NC}"
else
    echo -e "${GREEN}✅ Docker Compose already installed${NC}"
fi
echo ""

# Step 3: Pull images
echo -e "${BLUE}� Step 3: Pulling Docker images...${NC}"
docker-compose -f docker-compose.hub.yml pull
echo -e "${GREEN}✅ Images pulled${NC}"
echo ""

# Step 4: Start services
echo -e "${BLUE}🚀 Step 4: Starting services...${NC}"
docker-compose -f docker-compose.hub.yml up -d
echo -e "${GREEN}✅ Services started${NC}"
echo ""

# Step 5: Wait for services
echo -e "${BLUE}⏳ Step 5: Waiting for services to be ready...${NC}"
sleep 30
echo -e "${GREEN}✅ Services ready${NC}"
echo ""

# Step 6: Verify MongoDB
echo -e "${BLUE}🔐 Step 6: Verifying MongoDB authentication...${NC}"
if docker exec cafe-pos-mongodb mongosh \
    -u "$MONGO_INITDB_ROOT_USERNAME" \
    -p "$MONGO_INITDB_ROOT_PASSWORD" \
    --authenticationDatabase admin \
    --eval "db.adminCommand('ping')" > /dev/null 2>&1; then
    echo -e "${GREEN}✅ MongoDB authentication successful${NC}"
else
    echo -e "${RED}❌ MongoDB authentication failed${NC}"
    exit 1
fi
echo ""

# Step 7: Seed data
echo -e "${BLUE}🌱 Step 7: Seeding initial data...${NC}"
docker exec -it cafe-pos-backend ./cafe-pos-server seed
echo -e "${GREEN}✅ Data seeded${NC}"
echo ""

# Summary
echo -e "${GREEN}╔════════════════════════════════════════════════════════════╗${NC}"
echo -e "${GREEN}║          ✅ EC2 Deployment Completed Successfully!        ║${NC}"
echo -e "${GREEN}╚════════════════════════════════════════════════════════════╝${NC}"
echo ""

echo -e "${BLUE}🌐 Access Application:${NC}"
echo "   Frontend: http://$(curl -s http://169.254.169.254/latest/meta-data/public-ipv4)"
echo "   Backend:  http://$(curl -s http://169.254.169.254/latest/meta-data/public-ipv4):3000"
echo ""

echo -e "${BLUE}🔑 Credentials:${NC}"
echo "   App Username: admin"
echo "   App Password: admin123"
echo "   ⚠️  Change password after first login!"
echo ""

echo -e "${BLUE}� Useful Commands:${NC}"
echo "   View logs:    docker-compose -f docker-compose.hub.yml logs -f"
echo "   Restart:      docker-compose -f docker-compose.hub.yml restart"
echo "   Stop:         docker-compose -f docker-compose.hub.yml down"
echo ""
