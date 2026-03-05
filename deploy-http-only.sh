#!/bin/bash

echo "=========================================="
echo "Deploy HTTP-Only Architecture"
echo "=========================================="
echo ""

# Colors
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
NC='\033[0m' # No Color

# Get local machine IP for print bridge
echo "Step 1: Detect local machine IP"
echo "--------------------------------"
LOCAL_IP=$(ifconfig | grep "inet " | grep -v 127.0.0.1 | awk '{print $2}' | head -1)
echo -e "${GREEN}Local IP detected: $LOCAL_IP${NC}"
echo ""

read -p "Is this correct? (y/n): " confirm
if [ "$confirm" != "y" ]; then
    read -p "Enter your local machine IP: " LOCAL_IP
fi

PRINT_BRIDGE_URL="http://${LOCAL_IP}:3001"
echo -e "${GREEN}Print Bridge URL: $PRINT_BRIDGE_URL${NC}"
echo ""

# Step 2: Build new backend image
echo "Step 2: Build backend image (without Chromium)"
echo "-----------------------------------------------"
cd backend
docker build -t linhtranphu/cafe-pos-backend:http-only .
docker tag linhtranphu/cafe-pos-backend:http-only linhtranphu/cafe-pos-backend:latest
echo -e "${GREEN}✅ Backend image built${NC}"
echo ""

# Step 3: Push to Docker Hub
echo "Step 3: Push to Docker Hub"
echo "--------------------------"
read -p "Push to Docker Hub? (y/n): " push_confirm
if [ "$push_confirm" = "y" ]; then
    docker push linhtranphu/cafe-pos-backend:http-only
    docker push linhtranphu/cafe-pos-backend:latest
    echo -e "${GREEN}✅ Images pushed${NC}"
else
    echo -e "${YELLOW}⚠️  Skipped push${NC}"
fi
echo ""

# Step 4: Update EC2 .env
echo "Step 4: Update EC2 environment"
echo "-------------------------------"
echo "Run this command on EC2:"
echo ""
echo -e "${YELLOW}ssh ubuntu@tacafe.store${NC}"
echo -e "${YELLOW}cd ~/cafe-pos${NC}"
echo -e "${YELLOW}echo 'PRINT_BRIDGE_URL=$PRINT_BRIDGE_URL' >> .env${NC}"
echo -e "${YELLOW}docker-compose pull backend${NC}"
echo -e "${YELLOW}docker-compose up -d backend${NC}"
echo -e "${YELLOW}docker logs -f cafe-pos-backend${NC}"
echo ""
read -p "Press Enter when EC2 is updated..."
echo ""

# Step 5: Setup Print Bridge locally
echo "Step 5: Setup Print Bridge (Local)"
echo "-----------------------------------"
cd ../local-print-bridge

if [ ! -f .env ]; then
    echo "Creating .env file..."
    cp .env.example .env
    echo -e "${YELLOW}⚠️  Please edit .env and set your printer IPs${NC}"
    read -p "Press Enter after editing .env..."
fi

echo "Installing dependencies..."
npm install

echo ""
echo "Starting Print Bridge..."
echo -e "${YELLOW}Run: npm start${NC}"
echo -e "${YELLOW}Or with PM2: pm2 start src/index.js --name print-bridge${NC}"
echo ""

# Step 6: Build and deploy frontend
echo "Step 6: Build frontend (optional)"
echo "----------------------------------"
read -p "Rebuild frontend? (y/n): " frontend_confirm
if [ "$frontend_confirm" = "y" ]; then
    cd ../frontend
    docker build -t linhtranphu/cafe-pos-frontend:latest .
    
    read -p "Push frontend to Docker Hub? (y/n): " push_frontend
    if [ "$push_frontend" = "y" ]; then
        docker push linhtranphu/cafe-pos-frontend:latest
        echo -e "${GREEN}✅ Frontend pushed${NC}"
    fi
    
    echo ""
    echo "Deploy frontend on EC2:"
    echo -e "${YELLOW}ssh ubuntu@tacafe.store${NC}"
    echo -e "${YELLOW}cd ~/cafe-pos${NC}"
    echo -e "${YELLOW}docker-compose pull frontend${NC}"
    echo -e "${YELLOW}docker-compose up -d frontend${NC}"
fi
echo ""

# Summary
echo "=========================================="
echo "Deployment Summary"
echo "=========================================="
echo ""
echo "✅ Backend image built (39.9MB, no Chromium)"
echo "✅ Print Bridge URL: $PRINT_BRIDGE_URL"
echo ""
echo "Next steps:"
echo "1. Start Print Bridge locally: npm start"
echo "2. Test: curl http://localhost:3001/health"
echo "3. Login to https://tacafe.store"
echo "4. Go to Settings → Test Print"
echo ""
echo "Troubleshooting:"
echo "- Check backend logs: docker logs cafe-pos-backend"
echo "- Check print bridge logs: npm start output"
echo "- Test connection: curl $PRINT_BRIDGE_URL/health"
echo ""
echo "Documentation: WEBSOCKET_REMOVED_HTTP_ONLY.md"
echo "=========================================="
