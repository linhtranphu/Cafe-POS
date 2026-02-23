#!/bin/bash

# Test Transfer Handover Flow
# This script tests the complete flow:
# 1. Create order with transfer payment
# 2. Handover transfer amount
# 3. Verify waiter's cash is not affected

set -e

API_URL="http://localhost:3000/api"
FRONTEND_URL="http://localhost:5173"

echo "==================================="
echo "🧪 TEST: Transfer Handover Flow"
echo "==================================="
echo ""

# Colors
GREEN='\033[0;32m'
RED='\033[0;31m'
YELLOW='\033[1;33m'
BLUE='\033[0;34m'
NC='\033[0m' # No Color

# Step 1: Login as waiter
echo -e "${BLUE}Step 1: Login as waiter${NC}"
WAITER_LOGIN=$(curl -s -X POST "$API_URL/login" \
  -H "Content-Type: application/json" \
  -d '{
    "username": "waiter",
    "password": "waiter123"
  }')

WAITER_TOKEN=$(echo $WAITER_LOGIN | grep -o '"token":"[^"]*"' | cut -d'"' -f4)
WAITER_ID=$(echo $WAITER_LOGIN | grep -o '"user_id":"[^"]*"' | cut -d'"' -f4)

if [ -z "$WAITER_TOKEN" ]; then
  echo -e "${RED}❌ Failed to login as waiter${NC}"
  exit 1
fi

echo -e "${GREEN}✅ Logged in as waiter${NC}"
echo "   Token: ${WAITER_TOKEN:0:20}..."
echo ""

# Step 2: Get current shift
echo -e "${BLUE}Step 2: Get waiter's current shift${NC}"
SHIFT=$(curl -s -X GET "$API_URL/shifts/current" \
  -H "Authorization: Bearer $WAITER_TOKEN")

SHIFT_ID=$(echo $SHIFT | grep -o '"id":"[^"]*"' | head -1 | cut -d'"' -f4)

if [ -z "$SHIFT_ID" ]; then
  echo -e "${YELLOW}⚠️  No open shift found, creating new shift...${NC}"
  
  # Create new shift
  NEW_SHIFT=$(curl -s -X POST "$API_URL/shifts" \
    -H "Authorization: Bearer $WAITER_TOKEN" \
    -H "Content-Type: application/json" \
    -d '{
      "shift_type": "MORNING",
      "start_cash": 50000
    }')
  
  SHIFT_ID=$(echo $NEW_SHIFT | grep -o '"id":"[^"]*"' | head -1 | cut -d'"' -f4)
  
  if [ -z "$SHIFT_ID" ]; then
    echo -e "${RED}❌ Failed to create shift${NC}"
    exit 1
  fi
  
  echo -e "${GREEN}✅ Created new shift${NC}"
fi

echo -e "${GREEN}✅ Found shift: $SHIFT_ID${NC}"

# Get shift details BEFORE
SHIFT_BEFORE=$(curl -s -X GET "$API_URL/shifts/$SHIFT_ID" \
  -H "Authorization: Bearer $WAITER_TOKEN")

CASH_BEFORE=$(echo $SHIFT_BEFORE | grep -o '"remaining_cash":[0-9]*' | cut -d':' -f2)
TRANSFER_BEFORE=$(echo $SHIFT_BEFORE | grep -o '"remaining_transfer":[0-9]*' | cut -d':' -f2)
CURRENT_CASH_BEFORE=$(echo $SHIFT_BEFORE | grep -o '"current_cash":[0-9]*' | cut -d':' -f2)

echo ""
echo -e "${YELLOW}📊 Shift BEFORE:${NC}"
echo "   Current Cash: ${CURRENT_CASH_BEFORE} VND"
echo "   Remaining Cash: ${CASH_BEFORE} VND"
echo "   Remaining Transfer: ${TRANSFER_BEFORE} VND"
echo ""

# Step 3: Create order with transfer payment
echo -e "${BLUE}Step 3: Create order with TRANSFER payment${NC}"

# Get menu items
MENU=$(curl -s -X GET "$API_URL/menu" \
  -H "Authorization: Bearer $WAITER_TOKEN")

MENU_ITEM_ID=$(echo $MENU | grep -o '"id":"[^"]*"' | head -1 | cut -d'"' -f4)

if [ -z "$MENU_ITEM_ID" ]; then
  echo -e "${RED}❌ No menu items found${NC}"
  exit 1
fi

# Create order
ORDER=$(curl -s -X POST "$API_URL/orders" \
  -H "Authorization: Bearer $WAITER_TOKEN" \
  -H "Content-Type: application/json" \
  -d "{
    \"shift_id\": \"$SHIFT_ID\",
    \"customer_name\": \"Test Customer\",
    \"items\": [
      {
        \"menu_item_id\": \"$MENU_ITEM_ID\",
        \"quantity\": 1
      }
    ]
  }")

ORDER_ID=$(echo $ORDER | grep -o '"id":"[^"]*"' | head -1 | cut -d'"' -f4)
ORDER_TOTAL=$(echo $ORDER | grep -o '"total":[0-9]*' | cut -d':' -f2)

if [ -z "$ORDER_ID" ]; then
  echo -e "${RED}❌ Failed to create order${NC}"
  exit 1
fi

echo -e "${GREEN}✅ Created order: $ORDER_ID${NC}"
echo "   Total: ${ORDER_TOTAL} VND"
echo ""

# Step 4: Pay with TRANSFER
echo -e "${BLUE}Step 4: Pay order with TRANSFER${NC}"

PAYMENT=$(curl -s -X POST "$API_URL/orders/$ORDER_ID/payment" \
  -H "Authorization: Bearer $WAITER_TOKEN" \
  -H "Content-Type: application/json" \
  -d "{
    \"amount\": $ORDER_TOTAL,
    \"payment_method\": \"TRANSFER\",
    \"collector_id\": \"$WAITER_ID\",
    \"collector_name\": \"waiter\"
  }")

PAYMENT_STATUS=$(echo $PAYMENT | grep -o '"status":"[^"]*"' | head -1 | cut -d'"' -f4)

if [ "$PAYMENT_STATUS" != "PAID" ]; then
  echo -e "${RED}❌ Payment failed${NC}"
  exit 1
fi

echo -e "${GREEN}✅ Payment successful (TRANSFER)${NC}"
echo ""

# Wait a bit for DB to update
sleep 1

# Get shift details AFTER payment
SHIFT_AFTER_PAYMENT=$(curl -s -X GET "$API_URL/shifts/$SHIFT_ID" \
  -H "Authorization: Bearer $WAITER_TOKEN")

CASH_AFTER_PAYMENT=$(echo $SHIFT_AFTER_PAYMENT | grep -o '"remaining_cash":[0-9]*' | cut -d':' -f2)
TRANSFER_AFTER_PAYMENT=$(echo $SHIFT_AFTER_PAYMENT | grep -o '"remaining_transfer":[0-9]*' | cut -d':' -f2)
CURRENT_CASH_AFTER_PAYMENT=$(echo $SHIFT_AFTER_PAYMENT | grep -o '"current_cash":[0-9]*' | cut -d':' -f2)

echo -e "${YELLOW}📊 Shift AFTER PAYMENT:${NC}"
echo "   Current Cash: ${CURRENT_CASH_AFTER_PAYMENT} VND"
echo "   Remaining Cash: ${CASH_AFTER_PAYMENT} VND"
echo "   Remaining Transfer: ${TRANSFER_AFTER_PAYMENT} VND"
echo ""

# Verify cash not changed, transfer increased
if [ "$CASH_AFTER_PAYMENT" -ne "$CASH_BEFORE" ]; then
  echo -e "${RED}❌ FAIL: Cash changed after transfer payment!${NC}"
  echo "   Expected: ${CASH_BEFORE} VND"
  echo "   Got: ${CASH_AFTER_PAYMENT} VND"
  exit 1
fi

if [ "$TRANSFER_AFTER_PAYMENT" -le "$TRANSFER_BEFORE" ]; then
  echo -e "${RED}❌ FAIL: Transfer did not increase!${NC}"
  echo "   Before: ${TRANSFER_BEFORE} VND"
  echo "   After: ${TRANSFER_AFTER_PAYMENT} VND"
  exit 1
fi

echo -e "${GREEN}✅ Payment correctly updated transfer, not cash${NC}"
echo ""

# Step 5: Create handover for transfer
echo -e "${BLUE}Step 5: Create handover for TRANSFER${NC}"

HANDOVER=$(curl -s -X POST "$API_URL/handovers" \
  -H "Authorization: Bearer $WAITER_TOKEN" \
  -H "Content-Type: application/json" \
  -d "{
    \"shift_id\": \"$SHIFT_ID\",
    \"cash_amount\": 0,
    \"transfer_amount\": $ORDER_TOTAL,
    \"handover_type\": \"PARTIAL\",
    \"waiter_note\": \"Test transfer handover\"
  }")

HANDOVER_ID=$(echo $HANDOVER | grep -o '"id":"[^"]*"' | head -1 | cut -d'"' -f4)

if [ -z "$HANDOVER_ID" ]; then
  echo -e "${RED}❌ Failed to create handover${NC}"
  echo "$HANDOVER"
  exit 1
fi

echo -e "${GREEN}✅ Created handover: $HANDOVER_ID${NC}"
echo "   Transfer Amount: ${ORDER_TOTAL} VND"
echo ""

# Step 6: Login as cashier and confirm handover
echo -e "${BLUE}Step 6: Login as cashier${NC}"

CASHIER_LOGIN=$(curl -s -X POST "$API_URL/login" \
  -H "Content-Type: application/json" \
  -d '{
    "username": "cashier",
    "password": "cashier123"
  }')

CASHIER_TOKEN=$(echo $CASHIER_LOGIN | grep -o '"token":"[^"]*"' | cut -d'"' -f4)

if [ -z "$CASHIER_TOKEN" ]; then
  echo -e "${RED}❌ Failed to login as cashier${NC}"
  exit 1
fi

echo -e "${GREEN}✅ Logged in as cashier${NC}"
echo ""

# Step 7: Confirm handover
echo -e "${BLUE}Step 7: Confirm handover${NC}"

CONFIRM=$(curl -s -X POST "$API_URL/handovers/$HANDOVER_ID/confirm" \
  -H "Authorization: Bearer $CASHIER_TOKEN" \
  -H "Content-Type: application/json" \
  -d "{
    \"status\": \"CONFIRMED\",
    \"actual_amount\": $ORDER_TOTAL,
    \"cashier_note\": \"Test confirmation\"
  }")

CONFIRM_STATUS=$(echo $CONFIRM | grep -o '"status":"[^"]*"' | head -1 | cut -d'"' -f4)

if [ "$CONFIRM_STATUS" != "CONFIRMED" ]; then
  echo -e "${RED}❌ Failed to confirm handover${NC}"
  echo "$CONFIRM"
  exit 1
fi

echo -e "${GREEN}✅ Handover confirmed${NC}"
echo ""

# Wait for DB update
sleep 1

# Step 8: Verify final shift state
echo -e "${BLUE}Step 8: Verify final shift state${NC}"

SHIFT_FINAL=$(curl -s -X GET "$API_URL/shifts/$SHIFT_ID" \
  -H "Authorization: Bearer $WAITER_TOKEN")

CASH_FINAL=$(echo $SHIFT_FINAL | grep -o '"remaining_cash":[0-9]*' | cut -d':' -f2)
TRANSFER_FINAL=$(echo $SHIFT_FINAL | grep -o '"remaining_transfer":[0-9]*' | cut -d':' -f2)
CURRENT_CASH_FINAL=$(echo $SHIFT_FINAL | grep -o '"current_cash":[0-9]*' | cut -d':' -f2)
HANDED_OVER_CASH=$(echo $SHIFT_FINAL | grep -o '"handed_over_cash":[0-9]*' | cut -d':' -f2)
HANDED_OVER_TRANSFER=$(echo $SHIFT_FINAL | grep -o '"handed_over_transfer":[0-9]*' | cut -d':' -f2)

echo ""
echo -e "${YELLOW}📊 Shift FINAL STATE:${NC}"
echo "   Current Cash: ${CURRENT_CASH_FINAL} VND"
echo "   Remaining Cash: ${CASH_FINAL} VND"
echo "   Remaining Transfer: ${TRANSFER_FINAL} VND"
echo "   Handed Over Cash: ${HANDED_OVER_CASH} VND"
echo "   Handed Over Transfer: ${HANDED_OVER_TRANSFER} VND"
echo ""

# Verify results
echo -e "${BLUE}Step 9: Verify results${NC}"
echo ""

PASS=true

# Check 1: Cash should not change from initial
if [ "$CASH_FINAL" -ne "$CASH_BEFORE" ]; then
  echo -e "${RED}❌ FAIL: Cash changed!${NC}"
  echo "   Initial: ${CASH_BEFORE} VND"
  echo "   Final: ${CASH_FINAL} VND"
  PASS=false
else
  echo -e "${GREEN}✅ PASS: Cash unchanged (${CASH_FINAL} VND)${NC}"
fi

# Check 2: Transfer should be 0 (all handed over)
if [ "$TRANSFER_FINAL" -ne 0 ]; then
  echo -e "${RED}❌ FAIL: Transfer not fully handed over!${NC}"
  echo "   Expected: 0 VND"
  echo "   Got: ${TRANSFER_FINAL} VND"
  PASS=false
else
  echo -e "${GREEN}✅ PASS: Transfer fully handed over (0 VND)${NC}"
fi

# Check 3: Handed over cash should be 0
if [ "$HANDED_OVER_CASH" -ne 0 ]; then
  echo -e "${RED}❌ FAIL: Handed over cash should be 0!${NC}"
  echo "   Expected: 0 VND"
  echo "   Got: ${HANDED_OVER_CASH} VND"
  PASS=false
else
  echo -e "${GREEN}✅ PASS: No cash handed over (0 VND)${NC}"
fi

# Check 4: Handed over transfer should equal order total
if [ "$HANDED_OVER_TRANSFER" -ne "$ORDER_TOTAL" ]; then
  echo -e "${RED}❌ FAIL: Handed over transfer incorrect!${NC}"
  echo "   Expected: ${ORDER_TOTAL} VND"
  echo "   Got: ${HANDED_OVER_TRANSFER} VND"
  PASS=false
else
  echo -e "${GREEN}✅ PASS: Transfer handed over correctly (${HANDED_OVER_TRANSFER} VND)${NC}"
fi

echo ""
echo "==================================="
if [ "$PASS" = true ]; then
  echo -e "${GREEN}✅ ALL TESTS PASSED!${NC}"
  echo "==================================="
  echo ""
  echo "Frontend URLs to verify:"
  echo "  Waiter Shift: ${FRONTEND_URL}/#/shifts"
  echo "  Cashier View: ${FRONTEND_URL}/#/cashier/handovers"
  exit 0
else
  echo -e "${RED}❌ SOME TESTS FAILED!${NC}"
  echo "==================================="
  exit 1
fi
