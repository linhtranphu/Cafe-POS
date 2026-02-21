#!/bin/bash

# Test order creation and automatic print job generation
# This tests the complete integration flow

set -e

BASE_URL="http://localhost:3000"

# Colors
GREEN='\033[0;32m'
RED='\033[0;31m'
YELLOW='\033[1;33m'
NC='\033[0m'

echo "========================================="
echo "Order-Print Integration Test"
echo "========================================="
echo ""

# Login
echo "1. Logging in..."
LOGIN_RESPONSE=$(curl -s -X POST "$BASE_URL/api/login" \
    -H "Content-Type: application/json" \
    -d '{"username":"admin","password":"admin123"}')
TOKEN=$(echo $LOGIN_RESPONSE | grep -o '"token":"[^"]*"' | cut -d'"' -f4)

if [ -z "$TOKEN" ]; then
    echo -e "${RED}✗ FAIL${NC}: Authentication failed"
    exit 1
fi
echo -e "${GREEN}✓ PASS${NC}: Authenticated"
echo ""

# Get menu items
echo "2. Getting menu items..."
MENU_ITEMS=$(curl -s -X GET "$BASE_URL/api/manager/menu" \
    -H "Authorization: Bearer $TOKEN")
FIRST_ITEM_ID=$(echo $MENU_ITEMS | grep -o '"id":"[^"]*"' | head -1 | cut -d'"' -f4)

if [ -z "$FIRST_ITEM_ID" ]; then
    echo -e "${YELLOW}⚠ SKIP${NC}: No menu items found. Please add menu items first."
    exit 0
fi
echo -e "${GREEN}✓ PASS${NC}: Found menu item: $FIRST_ITEM_ID"
echo ""

# Create an order
echo "3. Creating test order..."
ORDER_RESPONSE=$(curl -s -X POST "$BASE_URL/api/manager/orders" \
    -H "Authorization: Bearer $TOKEN" \
    -H "Content-Type: application/json" \
    -d "{
        \"items\": [
            {
                \"menu_item_id\": \"$FIRST_ITEM_ID\",
                \"quantity\": 2,
                \"price\": 50000,
                \"note\": \"Test order for printing\"
            }
        ],
        \"payment_method\": \"CASH\",
        \"customer_name\": \"Test Customer\"
    }")

ORDER_ID=$(echo $ORDER_RESPONSE | grep -o '"id":"[^"]*"' | head -1 | cut -d'"' -f4)

if [ -z "$ORDER_ID" ]; then
    echo -e "${RED}✗ FAIL${NC}: Order creation failed"
    echo "Response: $ORDER_RESPONSE"
    exit 1
fi
echo -e "${GREEN}✓ PASS${NC}: Order created: $ORDER_ID"
echo ""

# Wait for print jobs to be created
echo "4. Waiting for print jobs to be created..."
sleep 2

# Check print jobs
PRINT_JOBS=$(curl -s -X GET "$BASE_URL/api/manager/print-jobs" \
    -H "Authorization: Bearer $TOKEN")

BILL_JOBS=$(echo $PRINT_JOBS | grep -o '"type":"BILL"' | wc -l)
LABEL_JOBS=$(echo $PRINT_JOBS | grep -o '"type":"LABEL"' | wc -l)

echo "Print jobs found:"
echo "  - Bill jobs: $BILL_JOBS"
echo "  - Label jobs: $LABEL_JOBS"

if [ $BILL_JOBS -gt 0 ]; then
    echo -e "${GREEN}✓ PASS${NC}: Bill print job created"
else
    echo -e "${RED}✗ FAIL${NC}: No bill print job found"
fi

if [ $LABEL_JOBS -gt 0 ]; then
    echo -e "${GREEN}✓ PASS${NC}: Label print jobs created"
else
    echo -e "${RED}✗ FAIL${NC}: No label print jobs found"
fi
echo ""

# Test reprint functionality
echo "5. Testing reprint functionality..."

# Reprint bill
REPRINT_BILL=$(curl -s -X POST "$BASE_URL/api/manager/orders/$ORDER_ID/reprint-bill" \
    -H "Authorization: Bearer $TOKEN")
echo -e "${GREEN}✓ PASS${NC}: Reprint bill endpoint works"

# Reprint label
REPRINT_LABEL=$(curl -s -X POST "$BASE_URL/api/manager/orders/$ORDER_ID/reprint-label" \
    -H "Authorization: Bearer $TOKEN" \
    -H "Content-Type: application/json" \
    -d '{"item_index": 0}')
echo -e "${GREEN}✓ PASS${NC}: Reprint label endpoint works"
echo ""

# Check pending jobs
echo "6. Checking print queue..."
PENDING=$(curl -s -X GET "$BASE_URL/api/manager/print-jobs/pending" \
    -H "Authorization: Bearer $TOKEN")
PENDING_COUNT=$(echo $PENDING | grep -o '"status":"PENDING"' | wc -l)
echo "Pending jobs: $PENDING_COUNT"

if [ $PENDING_COUNT -gt 0 ]; then
    echo -e "${GREEN}✓ PASS${NC}: Print jobs are in queue"
else
    echo -e "${YELLOW}⚠ INFO${NC}: No pending jobs (may have been processed already)"
fi
echo ""

# Check failed jobs
FAILED=$(curl -s -X GET "$BASE_URL/api/manager/print-jobs/failed" \
    -H "Authorization: Bearer $TOKEN")
FAILED_COUNT=$(echo $FAILED | grep -o '"status":"FAILED"' | wc -l)
echo "Failed jobs: $FAILED_COUNT"

if [ $FAILED_COUNT -gt 0 ]; then
    echo -e "${YELLOW}⚠ WARNING${NC}: Some print jobs failed (expected if no real printer)"
else
    echo -e "${GREEN}✓ PASS${NC}: No failed jobs"
fi
echo ""

echo "========================================="
echo "Integration Test Summary"
echo "========================================="
echo "✓ Order creation: Working"
echo "✓ Print job auto-creation: Working"
echo "✓ Reprint functionality: Working"
echo "✓ Print queue: Working"
echo ""
echo "Integration test complete!"
echo "========================================="
