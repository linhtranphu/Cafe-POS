#!/bin/bash

# Script để fix WebSocket connection trên EC2
# Vấn đề: Nginx thiếu proxy config cho Socket.IO

set -e

GREEN='\033[0;32m'
YELLOW='\033[1;33m'
RED='\033[0;31m'
NC='\033[0m'

echo -e "${YELLOW}🔧 Fixing WebSocket Connection on EC2${NC}"
echo ""

# Kiểm tra nginx.conf đã được update chưa
if grep -q "location /socket.io/" frontend/nginx.conf; then
    echo -e "${GREEN}✅ nginx.conf đã có Socket.IO proxy config${NC}"
else
    echo -e "${RED}❌ nginx.conf chưa có Socket.IO proxy config${NC}"
    echo "Vui lòng chạy lại script này sau khi update nginx.conf"
    exit 1
fi

echo ""
echo -e "${YELLOW}📋 Các bước thực hiện:${NC}"
echo "1. Rebuild frontend image với nginx config mới"
echo "2. Push image lên Docker Hub"
echo "3. Pull và restart trên EC2"
echo ""

read -p "Tiếp tục? (y/n) " -n 1 -r
echo
if [[ ! $REPLY =~ ^[Yy]$ ]]; then
    exit 1
fi

# Step 1: Rebuild frontend
echo ""
echo -e "${YELLOW}1️⃣  Rebuilding frontend image...${NC}"
cd frontend
docker build --no-cache -t linhtranphu/cafe-pos-frontend:latest .
if [ $? -ne 0 ]; then
    echo -e "${RED}❌ Build failed${NC}"
    exit 1
fi
cd ..
echo -e "${GREEN}✅ Frontend rebuilt${NC}"

# Step 2: Push to Docker Hub
echo ""
echo -e "${YELLOW}2️⃣  Pushing to Docker Hub...${NC}"
docker push linhtranphu/cafe-pos-frontend:latest
if [ $? -ne 0 ]; then
    echo -e "${RED}❌ Push failed${NC}"
    exit 1
fi
echo -e "${GREEN}✅ Pushed to Docker Hub${NC}"

# Step 3: Instructions for EC2
echo ""
echo -e "${GREEN}✅ Build và push thành công!${NC}"
echo ""
echo -e "${YELLOW}📝 Trên EC2, chạy các lệnh sau:${NC}"
echo ""
echo "# Pull image mới"
echo "docker pull linhtranphu/cafe-pos-frontend:latest"
echo ""
echo "# Restart frontend container"
echo "docker-compose -f docker-compose.prod.yml up -d --force-recreate frontend"
echo ""
echo "# Hoặc restart toàn bộ"
echo "docker-compose -f docker-compose.prod.yml restart"
echo ""
echo -e "${YELLOW}🧪 Test WebSocket sau khi restart:${NC}"
echo "1. Mở browser console tại https://tacafe.store"
echo "2. Kiểm tra log: [WebSocket] Connected"
echo "3. Không còn lỗi 'WebSocket is closed before the connection is established'"
echo ""
