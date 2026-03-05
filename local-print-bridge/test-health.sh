#!/bin/bash

# Test Local Print Bridge Health Check

GREEN='\033[0;32m'
YELLOW='\033[1;33m'
RED='\033[0;31m'
NC='\033[0m'

PRINT_BRIDGE_URL="${1:-http://localhost:3001}"

echo "Testing Print Bridge Health Check..."
echo "URL: ${PRINT_BRIDGE_URL}/health"
echo ""

# Test health endpoint
response=$(curl -s -w "\n%{http_code}" "${PRINT_BRIDGE_URL}/health" 2>&1)
http_code=$(echo "$response" | tail -n1)
body=$(echo "$response" | head -n-1)

if [ "$http_code" = "200" ]; then
    echo -e "${GREEN}✅ Health check passed!${NC}"
    echo ""
    echo "Response:"
    echo "$body" | jq . 2>/dev/null || echo "$body"
    echo ""
    
    # Test status endpoint
    echo "Testing /status endpoint..."
    status_response=$(curl -s "${PRINT_BRIDGE_URL}/status")
    echo "$status_response" | jq . 2>/dev/null || echo "$status_response"
    echo ""
    
    echo -e "${GREEN}✅ Print Bridge is running and healthy!${NC}"
else
    echo -e "${RED}❌ Health check failed!${NC}"
    echo "HTTP Code: $http_code"
    echo "Response: $body"
    echo ""
    echo -e "${YELLOW}Troubleshooting:${NC}"
    echo "1. Check if print bridge is running:"
    echo "   ps aux | grep print-bridge"
    echo ""
    echo "2. Check if port 3001 is in use:"
    echo "   lsof -i :3001"
    echo ""
    echo "3. Start print bridge:"
    echo "   ./print-bridge"
    echo "   # or"
    echo "   go run main.go"
    echo "   # or"
    echo "   docker-compose up -d"
    exit 1
fi
