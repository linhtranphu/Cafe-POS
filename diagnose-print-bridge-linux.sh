#!/bin/bash

# Script để chẩn đoán Print Bridge trên Linux
# Kiểm tra Chromium và dependencies

set -e

RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
BLUE='\033[0;34m'
NC='\033[0m'

echo -e "${BLUE}=== Print Bridge Linux Diagnostic ===${NC}"
echo ""

# 1. Check if running on Linux
echo "1. Checking OS..."
if [[ "$OSTYPE" == "linux-gnu"* ]]; then
    echo -e "${GREEN}✓ Running on Linux${NC}"
else
    echo -e "${YELLOW}⚠ Not running on Linux (OS: $OSTYPE)${NC}"
fi
echo ""

# 2. Check Docker container
echo "2. Checking Docker container..."
if docker ps --format '{{.Names}}' | grep -q "local-print-bridge"; then
    echo -e "${GREEN}✓ Container is running${NC}"
    
    # Get container ID
    CONTAINER_ID=$(docker ps --filter "name=local-print-bridge" --format "{{.ID}}")
    echo "   Container ID: $CONTAINER_ID"
else
    echo -e "${RED}✗ Container is not running${NC}"
    exit 1
fi
echo ""

# 3. Check Chromium inside container
echo "3. Checking Chromium in container..."

# Check if chromium-browser exists
if docker exec local-print-bridge which chromium-browser > /dev/null 2>&1; then
    echo -e "${GREEN}✓ chromium-browser found${NC}"
    CHROME_PATH=$(docker exec local-print-bridge which chromium-browser)
    echo "   Path: $CHROME_PATH"
else
    echo -e "${RED}✗ chromium-browser not found${NC}"
fi

# Check Chromium version
echo ""
echo "Chromium version:"
docker exec local-print-bridge chromium-browser --version 2>&1 || echo -e "${RED}Failed to get version${NC}"

echo ""

# 4. Test Chromium execution
echo "4. Testing Chromium execution..."
echo "Running headless test..."

TEST_OUTPUT=$(docker exec local-print-bridge chromium-browser \
    --headless \
    --disable-gpu \
    --no-sandbox \
    --disable-dev-shm-usage \
    --dump-dom \
    about:blank 2>&1 || echo "FAILED")

if echo "$TEST_OUTPUT" | grep -q "html"; then
    echo -e "${GREEN}✓ Chromium can render HTML${NC}"
else
    echo -e "${RED}✗ Chromium failed to render${NC}"
    echo "Error output:"
    echo "$TEST_OUTPUT" | head -20
fi
echo ""

# 5. Check shared memory
echo "5. Checking shared memory..."
SHM_SIZE=$(docker inspect local-print-bridge | grep -i "ShmSize" | awk '{print $2}' | tr -d ',')
if [ ! -z "$SHM_SIZE" ] && [ "$SHM_SIZE" != "0" ]; then
    SHM_MB=$((SHM_SIZE / 1024 / 1024))
    echo -e "${GREEN}✓ Shared memory: ${SHM_MB}MB${NC}"
    if [ "$SHM_MB" -lt 128 ]; then
        echo -e "${YELLOW}⚠ Shared memory is low, recommend 256MB+${NC}"
    fi
else
    echo -e "${YELLOW}⚠ Shared memory not configured${NC}"
fi
echo ""

# 6. Check environment variables
echo "6. Checking environment variables..."
docker exec local-print-bridge env | grep -E "CHROME|DISPLAY" || echo "No Chrome env vars set"
echo ""

# 7. Check recent logs for errors
echo "7. Checking recent logs for errors..."
echo "Last 20 lines:"
docker logs --tail=20 local-print-bridge 2>&1
echo ""

# 8. Test print-html endpoint
echo "8. Testing /print-html endpoint..."

# Get Print Bridge URL from cloudflare tunnel or localhost
PRINT_BRIDGE_URL="http://localhost:3001"

echo "Testing with simple HTML..."
TEST_RESPONSE=$(curl -s -w "\nHTTP_CODE:%{http_code}" -X POST "$PRINT_BRIDGE_URL/print-html" \
    -H "Content-Type: application/json" \
    -d '{
        "html": "<html><body><h1>Test Print</h1></body></html>",
        "printerIP": "192.168.1.100",
        "printerPort": 9100,
        "paperWidth": 80
    }' 2>&1)

HTTP_CODE=$(echo "$TEST_RESPONSE" | grep "HTTP_CODE:" | cut -d: -f2)
RESPONSE_BODY=$(echo "$TEST_RESPONSE" | grep -v "HTTP_CODE:")

echo "HTTP Code: $HTTP_CODE"
echo "Response: $RESPONSE_BODY"

if [ "$HTTP_CODE" = "200" ]; then
    echo -e "${GREEN}✓ Print endpoint responded successfully${NC}"
else
    echo -e "${RED}✗ Print endpoint failed${NC}"
fi
echo ""

# 9. Check for common issues
echo "9. Checking for common issues..."

# Check if running in privileged mode or with proper security opts
SECURITY_OPT=$(docker inspect local-print-bridge | grep -i "seccomp")
if [ ! -z "$SECURITY_OPT" ]; then
    echo -e "${GREEN}✓ Security options configured${NC}"
else
    echo -e "${YELLOW}⚠ No security options found${NC}"
    echo "   Chromium may need --no-sandbox flag"
fi
echo ""

# 10. Recommendations
echo -e "${BLUE}=== Recommendations ===${NC}"
echo ""

if echo "$TEST_OUTPUT" | grep -q "FAILED\|error\|Error"; then
    echo "Issues detected with Chromium. Try these fixes:"
    echo ""
    echo "1. Rebuild with updated Dockerfile:"
    echo "   cd local-print-bridge"
    echo "   docker compose down"
    echo "   docker compose build --no-cache"
    echo "   docker compose up -d"
    echo ""
    echo "2. Increase shared memory in docker-compose.yml:"
    echo "   shm_size: '512mb'"
    echo ""
    echo "3. Add security options in docker-compose.yml:"
    echo "   security_opt:"
    echo "     - seccomp:unconfined"
    echo ""
    echo "4. Check logs after restart:"
    echo "   docker logs -f local-print-bridge"
else
    echo -e "${GREEN}✓ Print Bridge appears to be working correctly${NC}"
    echo ""
    echo "If you still have issues:"
    echo "1. Check printer IP and port are correct"
    echo "2. Verify printer is accessible from Linux machine"
    echo "3. Check firewall settings"
fi

echo ""
echo -e "${BLUE}=== Diagnostic Complete ===${NC}"
