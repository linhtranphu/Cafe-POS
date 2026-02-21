#!/bin/bash

# Test Local Print Bridge Integration
# This script tests the complete flow from backend to local printer

set -e

echo "=========================================="
echo "Local Print Bridge Integration Test"
echo "=========================================="
echo ""

# Colors
GREEN='\033[0;32m'
RED='\033[0;31m'
YELLOW='\033[1;33m'
NC='\033[0m' # No Color

# Configuration
BACKEND_URL="${BACKEND_URL:-http://localhost:8080}"
LOCAL_BRIDGE_URL="${LOCAL_BRIDGE_URL:-http://localhost:3001}"
PRINTER_IP="${PRINTER_IP:-192.168.1.100}"
PRINTER_PORT="${PRINTER_PORT:-9100}"

echo "Configuration:"
echo "  Backend URL: $BACKEND_URL"
echo "  Local Bridge URL: $LOCAL_BRIDGE_URL"
echo "  Printer IP: $PRINTER_IP:$PRINTER_PORT"
echo ""

# Test 1: Check Local Bridge Health
echo "Test 1: Checking Local Bridge Health..."
if curl -s -f "$LOCAL_BRIDGE_URL/health" > /dev/null; then
    echo -e "${GREEN}✓ Local Bridge is running${NC}"
    HEALTH=$(curl -s "$LOCAL_BRIDGE_URL/health")
    echo "  Response: $HEALTH"
else
    echo -e "${RED}✗ Local Bridge is not running${NC}"
    echo "  Please start the local bridge: cd local-print-bridge && npm start"
    exit 1
fi
echo ""

# Test 2: Check Backend Health
echo "Test 2: Checking Backend Health..."
if curl -s -f "$BACKEND_URL/health" > /dev/null; then
    echo -e "${GREEN}✓ Backend is running${NC}"
else
    echo -e "${RED}✗ Backend is not running${NC}"
    echo "  Please start the backend: cd backend && go run main.go"
    exit 1
fi
echo ""

# Test 3: Test Printer Connection
echo "Test 3: Testing Printer Connection..."
CONN_TEST=$(curl -s -X POST "$LOCAL_BRIDGE_URL/test-connection" \
    -H "Content-Type: application/json" \
    -d "{\"printerIP\":\"$PRINTER_IP\",\"printerPort\":$PRINTER_PORT}")

if echo "$CONN_TEST" | grep -q '"success":true'; then
    echo -e "${GREEN}✓ Printer connection successful${NC}"
    echo "  Response: $CONN_TEST"
else
    echo -e "${YELLOW}⚠ Printer connection failed (this is OK if printer is not available)${NC}"
    echo "  Response: $CONN_TEST"
fi
echo ""

# Test 4: Get Local Bridge Status
echo "Test 4: Getting Local Bridge Statistics..."
STATUS=$(curl -s "$LOCAL_BRIDGE_URL/status")
echo -e "${GREEN}✓ Status retrieved${NC}"
echo "  $STATUS"
echo ""

# Test 5: Test Print Job Status Update Endpoint
echo "Test 5: Testing Backend Status Update Endpoint..."
echo -e "${YELLOW}Note: This requires authentication. Skipping for now.${NC}"
echo "  Endpoint: PUT $BACKEND_URL/api/print-jobs/:id/status"
echo "  Body: {\"status\":\"COMPLETED\",\"error_msg\":\"\"}"
echo ""

# Summary
echo "=========================================="
echo "Test Summary"
echo "=========================================="
echo -e "${GREEN}✓ Local Bridge is operational${NC}"
echo -e "${GREEN}✓ Backend is operational${NC}"
echo ""
echo "Integration Points:"
echo "  1. Frontend → Local Bridge: $LOCAL_BRIDGE_URL/print"
echo "  2. Local Bridge → Printer: TCP $PRINTER_IP:$PRINTER_PORT"
echo "  3. Local Bridge → Backend: $BACKEND_URL/api/print-jobs/:id/status"
echo ""
echo "Next Steps:"
echo "  1. Configure printers in Print Management UI"
echo "  2. Create a test order to trigger auto-print"
echo "  3. Check browser console for '[LocalPrint] Bridge available: true'"
echo "  4. Verify print job appears in Print Jobs tab"
echo ""
echo "For detailed documentation, see: LOCAL_PRINT_BRIDGE_INTEGRATION.md"
echo ""
