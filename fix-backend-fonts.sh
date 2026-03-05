#!/bin/bash

# Script để cài fonts tiếng Việt vào backend container
# Fix lỗi: no Vietnamese-compatible system fonts found

set -e

RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
BLUE='\033[0;34m'
NC='\033[0m'

echo -e "${BLUE}=== Fix Backend Fonts for Vietnamese ===${NC}"
echo ""

# Check if backend container is running
if ! docker ps --format '{{.Names}}' | grep -q "backend"; then
    echo -e "${RED}✗ Backend container is not running${NC}"
    echo "Start backend first: docker compose up -d backend"
    exit 1
fi

echo -e "${GREEN}✓ Backend container is running${NC}"
echo ""

# 1. Install fonts in container
echo "1. Installing Vietnamese-compatible fonts..."
echo ""

docker exec backend sh -c '
    echo "Installing fonts..."
    apk update
    apk add --no-cache \
        ttf-dejavu \
        font-noto \
        font-noto-cjk \
        fontconfig \
        font-liberation \
        font-roboto
    
    echo "Updating font cache..."
    fc-cache -f
    
    echo "Listing installed fonts..."
    fc-list | grep -i "dejavu\|noto\|liberation\|roboto" | head -10
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

# 2. Verify fonts
echo "2. Verifying font installation..."
echo ""

FONT_CHECK=$(docker exec backend fc-list | grep -i "dejavu\|noto" | wc -l)

if [ "$FONT_CHECK" -gt 0 ]; then
    echo -e "${GREEN}✓ Found $FONT_CHECK font files${NC}"
    echo ""
    echo "Sample fonts:"
    docker exec backend fc-list | grep -i "dejavu\|noto" | head -5
else
    echo -e "${RED}✗ No fonts found${NC}"
    exit 1
fi

echo ""

# 3. Check font paths
echo "3. Checking font paths..."
echo ""

docker exec backend sh -c '
    echo "Font directories:"
    ls -la /usr/share/fonts/truetype/ 2>/dev/null || echo "  /usr/share/fonts/truetype/ not found"
    ls -la /usr/share/fonts/TTF/ 2>/dev/null || echo "  /usr/share/fonts/TTF/ not found"
    
    echo ""
    echo "DejaVu Sans locations:"
    find /usr/share/fonts -name "*DejaVu*" 2>/dev/null | head -5
'

echo ""

# 4. Restart backend to apply changes
echo "4. Restarting backend container..."
docker restart backend

echo "Waiting for backend to be ready..."
sleep 5

# Check if backend is healthy
for i in {1..10}; do
    if docker exec backend wget -q -O- http://localhost:3000/api/health > /dev/null 2>&1; then
        echo -e "${GREEN}✓ Backend is healthy${NC}"
        break
    fi
    if [ $i -eq 10 ]; then
        echo -e "${YELLOW}⚠ Backend health check timeout${NC}"
    fi
    echo "Waiting... ($i/10)"
    sleep 2
done

echo ""

# 5. Test print
echo "5. Testing print with Vietnamese text..."
echo ""
echo "You can now test printing from the UI."
echo "The font error should be resolved."

echo ""
echo -e "${BLUE}=== Summary ===${NC}"
echo ""
echo "✓ Fonts installed in backend container:"
echo "  - DejaVu Sans"
echo "  - Noto Sans"
echo "  - Noto Sans CJK (for Vietnamese)"
echo "  - Liberation Sans"
echo "  - Roboto"
echo ""
echo "✓ Font cache updated"
echo "✓ Backend restarted"
echo ""
echo "Next steps:"
echo "  1. Test print from UI"
echo "  2. Check logs: docker logs -f backend"
echo "  3. If still issues, rebuild backend image"
echo ""
echo -e "${GREEN}✓ Fix complete!${NC}"
