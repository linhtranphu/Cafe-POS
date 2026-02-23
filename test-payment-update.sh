#!/bin/bash

# Test script to verify current_cash and transfer_revenue updates on payment

echo "🧪 Testing Payment Updates to Shift"
echo "===================================="
echo ""

# Colors
GREEN='\033[0;32m'
RED='\033[0;31m'
YELLOW='\033[1;33m'
NC='\033[0m' # No Color

# API Base URL
API_URL="http://localhost:8080/api"

# Login as waiter
echo "1️⃣  Logging in as waiter..."
LOGIN_RESPONSE=$(curl -s -X POST "$API_URL/auth/login" \
  -H "Content-Type: application/json" \
  -d '{
    "username": "waiter1",
    "password": "password123"
  }')

TOKEN=$(echo $LOGIN_RESPONSE | jq -r '.token')
USER_ID=$(echo $LOGIN_RESPONSE | jq -r '.user.id')

if [ "$TOKEN" == "null" ] || [ -z "$TOKEN" ]; then
  echo -e "${RED}❌ Login failed${NC}"
  exit 1
fi

echo -e "${GREEN}✅ Logged in successfully${NC}"
echo ""

# Start a shift
echo "2️⃣  Starting a new shift..."
SHIFT_RESPONSE=$(curl -s -X POST "$API_URL/waiter/shifts/start" \
  -H "Authorization: Bearer $TOKEN" \
  -H "Content-Type: application/json" \
  -d '{
    "type": "MORNING",
    "start_cash": 100000,
    "role_type": "waiter"
  }')

SHIFT_ID=$(echo $SHIFT_RESPONSE | jq -r '.id')

if [ "$SHIFT_ID" == "null" ] || [ -z "$SHIFT_ID" ]; then
  echo -e "${RED}❌ Failed to start shift${NC}"
  echo $SHIFT_RESPONSE | jq .
  exit 1
fi

echo -e "${GREEN}✅ Shift started: $SHIFT_ID${NC}"
echo ""

# Get shift details before orders
echo "3️⃣  Getting shift details before orders..."
SHIFT_BEFORE=$(curl -s -X GET "$API_URL/waiter/shifts/$SHIFT_ID" \
  -H "Authorization: Bearer $TOKEN")

echo "Shift before orders:"
echo $SHIFT_BEFORE | jq '{
  current_cash,
  remaining_cash,
  transfer_revenue,
  remaining_transfer,
  total_revenue
}'
echo ""

# Create order 1 (will pay with CASH)
echo "4️⃣  Creating order 1..."
ORDER1_RESPONSE=$(curl -s -X POST "$API_URL/waiter/orders" \
  -H "Authorization: Bearer $TOKEN" \
  -H "Content-Type: application/json" \
  -d "{
    \"shift_id\": \"$SHIFT_ID\",
    \"customer_name\": \"Customer 1\",
    \"items\": [
      {
        \"menu_item_id\": \"menu_item_id_here\",
        \"quantity\": 2,
        \"note\": \"Test order 1\"
      }
    ]
  }")

ORDER1_ID=$(echo $ORDER1_RESPONSE | jq -r '.id')

if [ "$ORDER1_ID" == "null" ] || [ -z "$ORDER1_ID" ]; then
  echo -e "${YELLOW}⚠️  Failed to create order 1 (might need valid menu_item_id)${NC}"
  echo $ORDER1_RESPONSE | jq .
else
  echo -e "${GREEN}✅ Order 1 created: $ORDER1_ID${NC}"
  
  # Pay order 1 with CASH
  echo "5️⃣  Paying order 1 with CASH (50,000 VND)..."
  PAYMENT1_RESPONSE=$(curl -s -X POST "$API_URL/waiter/orders/$ORDER1_ID/payment" \
    -H "Authorization: Bearer $TOKEN" \
    -H "Content-Type: application/json" \
    -d "{
      \"payment_method\": \"CASH\",
      \"amount\": 50000,
      \"collector_id\": \"$USER_ID\",
      \"collector_name\": \"Waiter 1\"
    }")
  
  echo $PAYMENT1_RESPONSE | jq .
  echo ""
fi

# Create order 2 (will pay with TRANSFER)
echo "6️⃣  Creating order 2..."
ORDER2_RESPONSE=$(curl -s -X POST "$API_URL/waiter/orders" \
  -H "Authorization: Bearer $TOKEN" \
  -H "Content-Type: application/json" \
  -d "{
    \"shift_id\": \"$SHIFT_ID\",
    \"customer_name\": \"Customer 2\",
    \"items\": [
      {
        \"menu_item_id\": \"menu_item_id_here\",
        \"quantity\": 1,
        \"note\": \"Test order 2\"
      }
    ]
  }")

ORDER2_ID=$(echo $ORDER2_RESPONSE | jq -r '.id')

if [ "$ORDER2_ID" == "null" ] || [ -z "$ORDER2_ID" ]; then
  echo -e "${YELLOW}⚠️  Failed to create order 2 (might need valid menu_item_id)${NC}"
  echo $ORDER2_RESPONSE | jq .
else
  echo -e "${GREEN}✅ Order 2 created: $ORDER2_ID${NC}"
  
  # Pay order 2 with TRANSFER
  echo "7️⃣  Paying order 2 with TRANSFER (75,000 VND)..."
  PAYMENT2_RESPONSE=$(curl -s -X POST "$API_URL/waiter/orders/$ORDER2_ID/payment" \
    -H "Authorization: Bearer $TOKEN" \
    -H "Content-Type: application/json" \
    -d "{
      \"payment_method\": \"TRANSFER\",
      \"amount\": 75000,
      \"collector_id\": \"$USER_ID\",
      \"collector_name\": \"Waiter 1\"
    }")
  
  echo $PAYMENT2_RESPONSE | jq .
  echo ""
fi

# Get shift details after payments
echo "8️⃣  Getting shift details after payments..."
SHIFT_AFTER=$(curl -s -X GET "$API_URL/waiter/shifts/$SHIFT_ID" \
  -H "Authorization: Bearer $TOKEN")

echo "Shift after payments:"
echo $SHIFT_AFTER | jq '{
  current_cash,
  remaining_cash,
  transfer_revenue,
  remaining_transfer,
  total_revenue
}'
echo ""

# Verify results
echo "9️⃣  Verification:"
echo "================"

CURRENT_CASH=$(echo $SHIFT_AFTER | jq -r '.current_cash')
REMAINING_CASH=$(echo $SHIFT_AFTER | jq -r '.remaining_cash')
TRANSFER_REVENUE=$(echo $SHIFT_AFTER | jq -r '.transfer_revenue')
REMAINING_TRANSFER=$(echo $SHIFT_AFTER | jq -r '.remaining_transfer')
TOTAL_REVENUE=$(echo $SHIFT_AFTER | jq -r '.total_revenue')

echo "Expected:"
echo "  - current_cash: 150,000 (100,000 start + 50,000 cash payment)"
echo "  - remaining_cash: 150,000"
echo "  - transfer_revenue: 75,000"
echo "  - remaining_transfer: 75,000"
echo "  - total_revenue: 125,000 (50,000 + 75,000)"
echo ""
echo "Actual:"
echo "  - current_cash: $CURRENT_CASH"
echo "  - remaining_cash: $REMAINING_CASH"
echo "  - transfer_revenue: $TRANSFER_REVENUE"
echo "  - remaining_transfer: $REMAINING_TRANSFER"
echo "  - total_revenue: $TOTAL_REVENUE"
echo ""

# Check if values are correct
if [ "$CURRENT_CASH" == "150000" ] && [ "$TRANSFER_REVENUE" == "75000" ]; then
  echo -e "${GREEN}✅ TEST PASSED: Both cash and transfer are updated correctly!${NC}"
else
  echo -e "${RED}❌ TEST FAILED: Values don't match expected results${NC}"
fi
