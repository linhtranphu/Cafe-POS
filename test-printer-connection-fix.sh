#!/bin/bash

# Script để test fix printer connection qua print bridge

set -e

GREEN='\033[0;32m'
YELLOW='\033[1;33m'
RED='\033[0;31m'
NC='\033[0m'

echo "=========================================="
echo "TEST PRINTER CONNECTION FIX"
echo "=========================================="
echo ""

echo -e "${YELLOW}1. Rebuild backend...${NC}"
docker-compose build backend
echo -e "${GREEN}✓ Backend rebuilt${NC}"
echo ""

echo -e "${YELLOW}2. Restart backend...${NC}"
docker-compose up -d backend
echo "Waiting for backend to be ready..."
sleep 5

for i in {1..10}; do
    if curl -s http://localhost:3000/health > /dev/null 2>&1; then
        echo -e "${GREEN}✓ Backend is ready${NC}"
        break
    fi
    if [ $i -eq 10 ]; then
        echo -e "${RED}❌ Backend health check timeout${NC}"
        exit 1
    fi
    echo "Waiting... ($i/10)"
    sleep 2
done
echo ""

echo -e "${YELLOW}3. Check print bridge...${NC}"
BRIDGE_URL=$(curl -s http://localhost:3000/api/settings | jq -r '.print_bridge_url // "NOT SET"')

if [ "$BRIDGE_URL" == "NOT SET" ] || [ "$BRIDGE_URL" == "null" ]; then
    echo -e "${RED}❌ Print Bridge URL chưa được cấu hình!${NC}"
    echo ""
    echo "Vui lòng cấu hình Print Bridge URL:"
    echo "  1. Vào /#/print-management"
    echo "  2. Tab 'Cài Đặt'"
    echo "  3. Điền Print Bridge URL (ví dụ: http://192.168.1.100:3001)"
    echo "  4. Bật 'Tự động in bill và tem khi tạo đơn hàng'"
    echo "  5. Bấm 'Lưu cài đặt'"
    echo ""
    echo "Hoặc chạy script setup:"
    echo "  ./setup-print-bridge-settings.sh"
    exit 1
fi

echo "Print Bridge URL: $BRIDGE_URL"

# Test print bridge
if curl -s -f "$BRIDGE_URL/health" > /dev/null 2>&1; then
    echo -e "${GREEN}✓ Print Bridge is available${NC}"
else
    echo -e "${RED}❌ Print Bridge is not available${NC}"
    exit 1
fi
echo ""

echo -e "${YELLOW}4. Check printers...${NC}"
PRINTERS=$(curl -s http://localhost:3000/api/manager/printers)
PRINTER_COUNT=$(echo $PRINTERS | jq '.printers | length')

echo "Found $PRINTER_COUNT printer(s)"

if [ "$PRINTER_COUNT" -eq 0 ]; then
    echo -e "${YELLOW}⚠️  No printers configured${NC}"
    echo "Please add a printer first at /#/print-management"
    exit 0
fi

# Get first printer
FIRST_PRINTER=$(echo $PRINTERS | jq -r '.printers[0]')
PRINTER_ID=$(echo $FIRST_PRINTER | jq -r '.id')
PRINTER_NAME=$(echo $FIRST_PRINTER | jq -r '.name')
PRINTER_IP=$(echo $FIRST_PRINTER | jq -r '.ip_address')

echo "Testing printer: $PRINTER_NAME ($PRINTER_IP)"
echo ""

echo -e "${YELLOW}5. Test printer connection...${NC}"

# Get token
TOKEN=$(cat ~/.cafe-pos-token 2>/dev/null || echo "")

if [ -z "$TOKEN" ]; then
    echo -e "${YELLOW}⚠️  No token found, please login first${NC}"
    exit 1
fi

# Test connection
RESULT=$(curl -s -X POST http://localhost:3000/api/manager/printers/$PRINTER_ID/test \
    -H "Authorization: Bearer $TOKEN")

SUCCESS=$(echo $RESULT | jq -r '.success // false')
MESSAGE=$(echo $RESULT | jq -r '.message // .error')

if [ "$SUCCESS" == "true" ]; then
    echo -e "${GREEN}✓ Test successful!${NC}"
    echo "Message: $MESSAGE"
else
    echo -e "${RED}❌ Test failed!${NC}"
    echo "Error: $MESSAGE"
    exit 1
fi

echo ""
echo "=========================================="
echo "HOÀN TẤT!"
echo "=========================================="
echo ""
echo -e "${GREEN}✓ Printer test connection đã hoạt động qua Print Bridge${NC}"
echo ""
echo "Bây giờ bạn có thể:"
echo "  1. Test connection ở tab Printers (sẽ qua Print Bridge)"
echo "  2. Không còn lỗi font nữa"
echo "  3. Collect payment sẽ tự động in bill"
echo ""
