#!/bin/bash

# Script để chẩn đoán tại sao không in bill khi thu tiền

set -e

RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
BLUE='\033[0;34m'
NC='\033[0m'

echo -e "${BLUE}=== Diagnose Auto-Print On Payment ===${NC}"
echo ""

# 1. Check auto-print setting
echo "1. Checking auto-print setting..."
SETTINGS=$(docker exec backend curl -s http://localhost:3000/api/settings 2>/dev/null || echo "")

if [ -z "$SETTINGS" ]; then
    echo -e "${RED}✗ Cannot connect to backend${NC}"
    exit 1
fi

AUTO_PRINT=$(echo "$SETTINGS" | grep -o '"auto_print_enabled":[^,}]*' | cut -d':' -f2 | tr -d ' ')

if [ "$AUTO_PRINT" = "true" ]; then
    echo -e "${GREEN}✓ Auto-print is ENABLED${NC}"
elif [ "$AUTO_PRINT" = "false" ]; then
    echo -e "${RED}✗ Auto-print is DISABLED${NC}"
    echo ""
    echo "To enable auto-print:"
    echo "  1. Go to Settings page"
    echo "  2. Enable 'Tự động in khi thu tiền'"
    echo "  3. Save settings"
    echo ""
    echo "Or enable via API:"
    echo '  curl -X PUT http://localhost:3000/api/settings \'
    echo '    -H "Content-Type: application/json" \'
    echo '    -d '"'"'{"auto_print_enabled": true}'"'"
    exit 0
else
    echo -e "${YELLOW}⚠ Cannot determine auto-print status${NC}"
    echo "Settings response: $SETTINGS"
fi

echo ""

# 2. Check Print Bridge connection
echo "2. Checking Print Bridge connection..."

PRINT_BRIDGE_URL=$(grep "^PRINT_BRIDGE_URL=" .env | cut -d'=' -f2 | tr -d '"' | tr -d "'")

if [ -z "$PRINT_BRIDGE_URL" ]; then
    echo -e "${RED}✗ PRINT_BRIDGE_URL not set in .env${NC}"
    exit 1
fi

echo "Print Bridge URL: $PRINT_BRIDGE_URL"

if curl -s "$PRINT_BRIDGE_URL/health" > /dev/null 2>&1; then
    echo -e "${GREEN}✓ Print Bridge is accessible${NC}"
else
    echo -e "${RED}✗ Cannot connect to Print Bridge${NC}"
    echo "Check:"
    echo "  1. Print Bridge is running"
    echo "  2. PRINT_BRIDGE_URL is correct in .env"
    echo "  3. Network connectivity"
    exit 1
fi

echo ""

# 3. Check recent orders
echo "3. Checking recent orders..."

# Get recent orders (last 5)
ORDERS=$(docker exec backend curl -s "http://localhost:3000/api/waiter/orders?limit=5" 2>/dev/null || echo "")

if [ -z "$ORDERS" ]; then
    echo -e "${RED}✗ Cannot fetch orders${NC}"
    exit 1
fi

echo "Recent orders:"
echo "$ORDERS" | grep -o '"order_number":"[^"]*"' | head -5 || echo "No orders found"

echo ""

# 4. Check print jobs
echo "4. Checking print jobs..."

PRINT_JOBS=$(docker exec backend curl -s "http://localhost:3000/api/print-jobs?limit=10" 2>/dev/null || echo "")

if [ -z "$PRINT_JOBS" ]; then
    echo -e "${YELLOW}⚠ Cannot fetch print jobs${NC}"
else
    PENDING_COUNT=$(echo "$PRINT_JOBS" | grep -o '"status":"pending"' | wc -l)
    COMPLETED_COUNT=$(echo "$PRINT_JOBS" | grep -o '"status":"completed"' | wc -l)
    FAILED_COUNT=$(echo "$PRINT_JOBS" | grep -o '"status":"failed"' | wc -l)
    
    echo "Print jobs status:"
    echo "  Pending: $PENDING_COUNT"
    echo "  Completed: $COMPLETED_COUNT"
    echo "  Failed: $FAILED_COUNT"
    
    if [ "$PENDING_COUNT" -gt 0 ]; then
        echo ""
        echo -e "${YELLOW}⚠ There are pending print jobs${NC}"
        echo "Check print worker logs:"
        echo "  docker logs backend | grep -i print"
    fi
    
    if [ "$FAILED_COUNT" -gt 0 ]; then
        echo ""
        echo -e "${RED}✗ There are failed print jobs${NC}"
        echo "Recent failed jobs:"
        echo "$PRINT_JOBS" | grep -A 5 '"status":"failed"' | head -20
    fi
fi

echo ""

# 5. Check backend logs for print-related errors
echo "5. Checking backend logs for print errors..."

PRINT_ERRORS=$(docker logs --tail=100 backend 2>&1 | grep -i "print\|ERROR" | tail -10)

if [ -z "$PRINT_ERRORS" ]; then
    echo -e "${GREEN}✓ No recent print errors in logs${NC}"
else
    echo -e "${YELLOW}Recent print-related logs:${NC}"
    echo "$PRINT_ERRORS"
fi

echo ""

# 6. Summary
echo -e "${BLUE}=== Summary ===${NC}"
echo ""

if [ "$AUTO_PRINT" = "true" ]; then
    echo -e "${GREEN}✓ Auto-print is enabled${NC}"
    echo ""
    echo "If bills are still not printing when you collect payment:"
    echo ""
    echo "1. Check if order is fully paid:"
    echo "   - Order must have amount_due = 0"
    echo "   - Order status must change to 'paid'"
    echo ""
    echo "2. Check print jobs are created:"
    echo "   - After payment, check: curl http://localhost:3000/api/print-jobs"
    echo "   - Should see new print jobs with order_number"
    echo ""
    echo "3. Check print worker is processing:"
    echo "   - docker logs -f backend | grep -i print"
    echo "   - Should see 'Processing print job' messages"
    echo ""
    echo "4. Check Print Bridge is working:"
    echo "   - curl $PRINT_BRIDGE_URL/health"
    echo "   - Test print from UI: Settings > Print Management"
    echo ""
    echo "5. Test full flow:"
    echo "   - Create order"
    echo "   - Collect payment (full amount)"
    echo "   - Watch logs: docker logs -f backend"
    echo "   - Check print jobs: curl http://localhost:3000/api/print-jobs"
else
    echo -e "${RED}✗ Auto-print is DISABLED${NC}"
    echo ""
    echo "Enable auto-print in Settings to automatically print bills on payment."
fi

echo ""
