#!/bin/bash

# Script để setup Print Bridge URL trong settings
# Giải quyết vấn đề test connection không qua print bridge

set -e

GREEN='\033[0;32m'
YELLOW='\033[1;33m'
RED='\033[0;31m'
NC='\033[0m'

echo "=========================================="
echo "SETUP PRINT BRIDGE SETTINGS"
echo "=========================================="
echo ""

# 1. Nhập thông tin
echo -e "${YELLOW}Nhập thông tin cấu hình:${NC}"
echo ""

read -p "Print Bridge URL (ví dụ: http://192.168.1.100:3001): " BRIDGE_URL
read -p "Tên quán: " SHOP_NAME
read -p "Địa chỉ quán (optional): " SHOP_ADDRESS
read -p "Số điện thoại (optional): " SHOP_PHONE

# Validate Bridge URL
if [ -z "$BRIDGE_URL" ]; then
    echo -e "${RED}❌ Print Bridge URL không được để trống!${NC}"
    exit 1
fi

if [ -z "$SHOP_NAME" ]; then
    echo -e "${RED}❌ Tên quán không được để trống!${NC}"
    exit 1
fi

echo ""
echo -e "${YELLOW}2. Kiểm tra Print Bridge...${NC}"

# Test Print Bridge connection
if curl -s -f "$BRIDGE_URL/health" > /dev/null 2>&1; then
    echo -e "${GREEN}✅ Print Bridge đang hoạt động: $BRIDGE_URL${NC}"
else
    echo -e "${RED}❌ Không thể kết nối đến Print Bridge: $BRIDGE_URL${NC}"
    echo "Vui lòng kiểm tra:"
    echo "  1. Print Bridge đang chạy: cd local-print-bridge && docker-compose ps"
    echo "  2. URL đúng (bao gồm http:// và port)"
    echo "  3. Firewall không chặn"
    exit 1
fi

echo ""
echo -e "${YELLOW}3. Kiểm tra settings hiện tại...${NC}"

# Check if settings exist
SETTINGS_CHECK=$(curl -s http://localhost:3000/api/settings 2>&1)

if echo "$SETTINGS_CHECK" | grep -q "shop_name"; then
    echo -e "${YELLOW}⚠️  Settings đã tồn tại, sẽ cập nhật...${NC}"
    METHOD="PUT"
else
    echo -e "${GREEN}✓ Chưa có settings, sẽ tạo mới...${NC}"
    METHOD="POST"
fi

echo ""
echo -e "${YELLOW}4. Lưu settings...${NC}"

# Get token (assuming you're logged in)
TOKEN=$(cat ~/.cafe-pos-token 2>/dev/null || echo "")

if [ -z "$TOKEN" ]; then
    echo -e "${YELLOW}⚠️  Không tìm thấy token, cần đăng nhập...${NC}"
    read -p "Nhập username: " USERNAME
    read -sp "Nhập password: " PASSWORD
    echo ""
    
    # Login to get token
    LOGIN_RESPONSE=$(curl -s -X POST http://localhost:3000/api/auth/login \
        -H "Content-Type: application/json" \
        -d "{\"username\":\"$USERNAME\",\"password\":\"$PASSWORD\"}")
    
    TOKEN=$(echo $LOGIN_RESPONSE | jq -r '.token // empty')
    
    if [ -z "$TOKEN" ]; then
        echo -e "${RED}❌ Đăng nhập thất bại!${NC}"
        exit 1
    fi
    
    # Save token for future use
    echo "$TOKEN" > ~/.cafe-pos-token
    echo -e "${GREEN}✓ Đăng nhập thành công${NC}"
fi

# Create/Update settings
RESPONSE=$(curl -s -X $METHOD http://localhost:3000/api/settings \
    -H "Content-Type: application/json" \
    -H "Authorization: Bearer $TOKEN" \
    -d "{
        \"shop_name\": \"$SHOP_NAME\",
        \"shop_address\": \"$SHOP_ADDRESS\",
        \"shop_phone\": \"$SHOP_PHONE\",
        \"print_bridge_url\": \"$BRIDGE_URL\",
        \"auto_print_enabled\": true,
        \"show_logo\": true,
        \"show_address\": true,
        \"show_phone\": true,
        \"show_custom_message\": true,
        \"custom_message\": \"Cảm ơn quý khách! Hẹn gặp lại...\"
    }")

if echo "$RESPONSE" | grep -q "shop_name"; then
    echo -e "${GREEN}✅ Đã lưu settings thành công!${NC}"
else
    echo -e "${RED}❌ Lỗi khi lưu settings:${NC}"
    echo "$RESPONSE"
    exit 1
fi

echo ""
echo -e "${YELLOW}5. Xác nhận settings...${NC}"

# Verify settings
VERIFY=$(curl -s http://localhost:3000/api/settings)
SAVED_BRIDGE_URL=$(echo $VERIFY | jq -r '.print_bridge_url // "NOT SET"')
SAVED_AUTO_PRINT=$(echo $VERIFY | jq -r '.auto_print_enabled // false')

echo "Print Bridge URL: $SAVED_BRIDGE_URL"
echo "Auto Print Enabled: $SAVED_AUTO_PRINT"

if [ "$SAVED_BRIDGE_URL" == "$BRIDGE_URL" ]; then
    echo -e "${GREEN}✅ Settings đã được lưu đúng!${NC}"
else
    echo -e "${RED}❌ Settings không khớp!${NC}"
    exit 1
fi

echo ""
echo -e "${YELLOW}6. Khởi động lại backend để áp dụng...${NC}"

# Restart backend to reload settings
if docker ps --format '{{.Names}}' | grep -q "backend"; then
    echo "Restarting backend container..."
    docker restart backend
    echo "Waiting for backend to be ready..."
    sleep 5
    
    # Check backend health
    for i in {1..10}; do
        if curl -s http://localhost:3000/health > /dev/null 2>&1; then
            echo -e "${GREEN}✅ Backend đã sẵn sàng${NC}"
            break
        fi
        if [ $i -eq 10 ]; then
            echo -e "${YELLOW}⚠️  Backend health check timeout${NC}"
        fi
        echo "Waiting... ($i/10)"
        sleep 2
    done
else
    echo -e "${YELLOW}⚠️  Backend container không chạy, vui lòng restart thủ công${NC}"
fi

echo ""
echo "=========================================="
echo "HOÀN TẤT!"
echo "=========================================="
echo ""
echo -e "${GREEN}✅ Print Bridge đã được cấu hình${NC}"
echo ""
echo "Bây giờ bạn có thể:"
echo "  1. Test connection ở tab Printers (sẽ qua Print Bridge)"
echo "  2. Test print ở tab Templates"
echo "  3. Collect payment sẽ tự động in bill"
echo ""
echo "Nếu vẫn gặp lỗi, check logs:"
echo "  - Backend: docker logs -f backend"
echo "  - Print Bridge: docker logs -f local-print-bridge"
echo ""
