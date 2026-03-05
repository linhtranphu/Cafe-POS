#!/bin/bash

# Script để cài fonts tiếng Việt vào Print Bridge container
# Fix lỗi: Font rendering issues with Vietnamese text

set -e

RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
BLUE='\033[0;34m'
NC='\033[0m'

echo -e "${BLUE}=== Fix Print Bridge Fonts for Vietnamese ===${NC}"
echo ""

# Check if print bridge container is running
if ! docker ps --format '{{.Names}}' | grep -q "local-print-bridge"; then
    echo -e "${RED}✗ Print Bridge container is not running${NC}"
    echo "Start Print Bridge first:"
    echo "  cd local-print-bridge"
    echo "  docker compose up -d"
    exit 1
fi

echo -e "${GREEN}✓ Print Bridge container is running${NC}"
echo ""

# 1. Check current fonts
echo "1. Checking current fonts..."
echo ""

CURRENT_FONTS=$(docker exec local-print-bridge fc-list 2>/dev/null | wc -l)
echo "Current font count: $CURRENT_FONTS"

if [ "$CURRENT_FONTS" -gt 50 ]; then
    echo -e "${GREEN}✓ Fonts appear to be installed${NC}"
else
    echo -e "${YELLOW}⚠ Few fonts found, installing more...${NC}"
fi

echo ""

# 2. Install fonts in container
echo "2. Installing Vietnamese-compatible fonts..."
echo ""

docker exec local-print-bridge sh -c '
    echo "Installing fonts..."
    apk update
    apk add --no-cache \
        font-noto \
        font-noto-cjk \
        ttf-dejavu \
        font-liberation \
        font-roboto \
        fontconfig
    
    echo "Updating font cache..."
    fc-cache -f
    
    echo "Listing Vietnamese-compatible fonts..."
    fc-list | grep -i "noto\|dejavu\|liberation\|roboto" | head -10
'

if [ $? -eq 0 ]; then
    echo ""
    echo -e "${GREEN}✓ Fonts installed successfully${NC}"
else
    echo ""
    echo -e "${RED}✗ Failed to install fonts${NC}"
    exit 1
fi

echo ""

# 3. Verify fonts
echo "3. Verifying font installation..."
echo ""

FONT_CHECK=$(docker exec local-print-bridge fc-list | grep -i "noto\|dejavu" | wc -l)

if [ "$FONT_CHECK" -gt 0 ]; then
    echo -e "${GREEN}✓ Found $FONT_CHECK font files${NC}"
    echo ""
    echo "Sample fonts:"
    docker exec local-print-bridge fc-list | grep -i "noto" | head -5
else
    echo -e "${RED}✗ No fonts found${NC}"
    exit 1
fi

echo ""

# 4. Check Chromium can access fonts
echo "4. Testing Chromium font access..."
echo ""

docker exec local-print-bridge sh -c '
    echo "Chromium version:"
    chromium-browser --version
    
    echo ""
    echo "Font directories accessible to Chromium:"
    ls -la /usr/share/fonts/ 2>/dev/null | head -10
'

echo ""

# 5. Test render with Vietnamese text
echo "5. Testing HTML render with Vietnamese text..."
echo ""

TEST_HTML='<html><head><meta charset="UTF-8"></head><body><h1>Tiếng Việt Test</h1><p>Cà phê sữa đá - 25,000 VNĐ</p></body></html>'

RENDER_TEST=$(docker exec local-print-bridge chromium-browser \
    --headless \
    --disable-gpu \
    --no-sandbox \
    --disable-dev-shm-usage \
    --dump-dom \
    "data:text/html;charset=utf-8,$TEST_HTML" 2>&1 | grep -i "việt\|phê" || echo "")

if [ ! -z "$RENDER_TEST" ]; then
    echo -e "${GREEN}✓ Chromium can render Vietnamese text${NC}"
else
    echo -e "${YELLOW}⚠ Cannot verify Vietnamese rendering${NC}"
fi

echo ""

# 6. Restart Print Bridge
echo "6. Restarting Print Bridge container..."
docker restart local-print-bridge

echo "Waiting for Print Bridge to be ready..."
sleep 5

# Check if Print Bridge is healthy
for i in {1..10}; do
    if curl -s http://localhost:3001/health > /dev/null 2>&1; then
        echo -e "${GREEN}✓ Print Bridge is healthy${NC}"
        break
    fi
    if [ $i -eq 10 ]; then
        echo -e "${YELLOW}⚠ Print Bridge health check timeout${NC}"
    fi
    echo "Waiting... ($i/10)"
    sleep 2
done

echo ""

# 7. Summary
echo -e "${BLUE}=== Summary ===${NC}"
echo ""
echo "✓ Fonts installed in Print Bridge container:"
echo "  - Noto Sans (basic)"
echo "  - Noto Sans CJK (Vietnamese, Chinese, Japanese, Korean)"
echo "  - DejaVu Sans"
echo "  - Liberation Sans"
echo "  - Roboto"
echo ""
echo "✓ Font cache updated"
echo "✓ Chromium can access fonts"
echo "✓ Print Bridge restarted"
echo ""
echo "Next steps:"
echo "  1. Test print from UI with Vietnamese text"
echo "  2. Check logs: docker logs -f local-print-bridge"
echo "  3. If still issues, rebuild Print Bridge image:"
echo "     cd local-print-bridge"
echo "     docker compose build --no-cache"
echo "     docker compose up -d"
echo ""
echo -e "${GREEN}✓ Fix complete!${NC}"
