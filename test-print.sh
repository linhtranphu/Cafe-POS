#!/bin/bash

# Test Print Script
# Tests printing to thermal printers via local-print-bridge

set -e

# Colors
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
RED='\033[0;31m'
NC='\033[0m'

BRIDGE_URL="http://localhost:3001"

echo "━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━"
echo "🖨️  Testing Local Print Bridge"
echo "━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━"
echo ""

# Check if bridge is running
echo -e "${YELLOW}📋 Step 1: Checking if print bridge is running...${NC}"
HEALTH=$(curl -s ${BRIDGE_URL}/health 2>/dev/null || echo "")

if [ -z "$HEALTH" ]; then
    echo -e "${RED}❌ Print bridge is not running${NC}"
    echo "Please start it first:"
    echo "  docker start print-bridge"
    echo "  OR"
    echo "  cd local-print-bridge && docker-compose up -d"
    exit 1
fi

echo -e "${GREEN}✅ Print bridge is running${NC}"
echo "$HEALTH" | grep -o '"status":"[^"]*"'
echo ""

# Test bill printer
echo -e "${YELLOW}📋 Step 2: Testing bill printer (80mm)...${NC}"
BILL_RESPONSE=$(curl -s -X POST ${BRIDGE_URL}/print \
  -H "Content-Type: application/json" \
  -d '{
    "jobId": "test-bill-001",
    "printerIP": "192.168.1.115",
    "printerPort": 9100,
    "type": "bill",
    "content": "================================\n       CAFE POS TEST\n================================\n\nTest Print - Bill Printer\nDate: 2026-02-21 21:45:00\n\nThis is a test print from\nLocal Print Bridge\n\n================================\n      Thank you!\n================================\n\n\n\n"
  }')

echo "Response: $BILL_RESPONSE"

if echo "$BILL_RESPONSE" | grep -q '"success":true'; then
    echo -e "${GREEN}✅ Bill printer test successful${NC}"
else
    echo -e "${RED}❌ Bill printer test failed${NC}"
    echo "Response: $BILL_RESPONSE"
fi

echo ""

# Test label printer
echo -e "${YELLOW}📋 Step 3: Testing label printer (58mm)...${NC}"
LABEL_RESPONSE=$(curl -s -X POST ${BRIDGE_URL}/print \
  -H "Content-Type: application/json" \
  -d '{
    "jobId": "test-label-001",
    "printerIP": "192.168.1.101",
    "printerPort": 9100,
    "type": "label",
    "content": "====================\n  CAFE POS TEST\n====================\n\nTest Label Print\n2026-02-21\n\nItem: Test Item\nQty: 1\n\n====================\n\n\n"
  }')

echo "Response: $LABEL_RESPONSE"

if echo "$LABEL_RESPONSE" | grep -q '"success":true'; then
    echo -e "${GREEN}✅ Label printer test successful${NC}"
else
    echo -e "${RED}❌ Label printer test failed${NC}"
    echo "Response: $LABEL_RESPONSE"
fi

echo ""
echo "━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━"
echo "✅ Print Test Complete!"
echo "━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━"
echo ""
echo "Check your printers for output."
echo ""
echo "To view bridge logs:"
echo "  docker logs -f print-bridge"
echo ""
