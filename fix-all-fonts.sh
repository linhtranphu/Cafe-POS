#!/bin/bash

# Script tổng hợp để cài fonts cho Backend và Print Bridge
# Fix lỗi: no Vietnamese-compatible system fonts found

set -e

RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
BLUE='\033[0;34m'
NC='\033[0m'

echo -e "${BLUE}========================================${NC}"
echo -e "${BLUE}   Fix Fonts for Vietnamese Printing${NC}"
echo -e "${BLUE}========================================${NC}"
echo ""

# 1. Fix Backend Fonts
echo -e "${BLUE}=== Part 1: Backend Fonts ===${NC}"
echo ""

if docker ps --format '{{.Names}}' | grep -q "backend"; then
    echo "Fixing backend fonts..."
    ./fix-backend-fonts.sh
    echo ""
else
    echo -e "${YELLOW}⚠ Backend container not running, skipping...${NC}"
    echo ""
fi

# 2. Fix Print Bridge Fonts
echo -e "${BLUE}=== Part 2: Print Bridge Fonts ===${NC}"
echo ""

if docker ps --format '{{.Names}}' | grep -q "local-print-bridge"; then
    echo "Fixing Print Bridge fonts..."
    ./fix-print-bridge-fonts.sh
    echo ""
else
    echo -e "${YELLOW}⚠ Print Bridge container not running, skipping...${NC}"
    echo ""
fi

# 3. Final Summary
echo ""
echo -e "${BLUE}========================================${NC}"
echo -e "${BLUE}           Summary${NC}"
echo -e "${BLUE}========================================${NC}"
echo ""

BACKEND_RUNNING=$(docker ps --format '{{.Names}}' | grep -c "backend" || echo "0")
PRINT_BRIDGE_RUNNING=$(docker ps --format '{{.Names}}' | grep -c "local-print-bridge" || echo "0")

if [ "$BACKEND_RUNNING" -gt 0 ]; then
    echo -e "${GREEN}✓ Backend fonts fixed${NC}"
else
    echo -e "${YELLOW}⚠ Backend not running${NC}"
fi

if [ "$PRINT_BRIDGE_RUNNING" -gt 0 ]; then
    echo -e "${GREEN}✓ Print Bridge fonts fixed${NC}"
else
    echo -e "${YELLOW}⚠ Print Bridge not running${NC}"
fi

echo ""
echo "Fonts installed:"
echo "  • Noto Sans (basic Latin)"
echo "  • Noto Sans CJK (Vietnamese, Chinese, Japanese, Korean)"
echo "  • DejaVu Sans"
echo "  • Liberation Sans"
echo "  • Roboto"
echo ""
echo "Next steps:"
echo "  1. Test print from UI"
echo "  2. Check backend logs: docker logs -f backend"
echo "  3. Check Print Bridge logs: docker logs -f local-print-bridge"
echo ""
echo "For permanent fix (recommended):"
echo "  1. Rebuild images: ./build_docker_hub.sh"
echo "  2. Push to DockerHub"
echo "  3. Pull on production: docker compose pull"
echo ""
echo -e "${GREEN}✓ All done!${NC}"
