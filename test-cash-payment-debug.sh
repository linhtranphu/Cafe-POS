#!/bin/bash

# Debug script to test cash payment update

echo "🔍 Debug: Testing Cash Payment Update"
echo "======================================"
echo ""

API_URL="http://localhost:8080/api"

# Login as waiter
echo "1️⃣  Login..."
LOGIN=$(curl -s -X POST "$API_URL/auth/login" \
  -H "Content-Type: application/json" \
  -d '{"username": "waiter1", "password": "password123"}')

TOKEN=$(echo $LOGIN | jq -r '.token')
USER_ID=$(echo $LOGIN | jq -r '.user.id')

if [ "$TOKEN" == "null" ]; then
  echo "❌ Login failed"
  exit 1
fi

echo "✅ Token: ${TOKEN:0:20}..."
echo ""

# Start shift
echo "2️⃣  Start shift with 100,000 VND..."
SHIFT=$(curl -s -X POST "$API_URL/waiter/shifts/start" \
  -H "Authorization: Bearer $TOKEN" \
  -H "Content-Type: application/json" \
  -d '{
    "type": "MORNING",
    "start_cash": 100000,
    "role_type": "waiter"
  }')

SHIFT_ID=$(echo $SHIFT | jq -r '.id')
echo "✅ Shift ID: $SHIFT_ID"
echo ""

# Get shift before
echo "3️⃣  Shift BEFORE payment:"
BEFORE=$(curl -s -X GET "$API_URL/waiter/shifts/$SHIFT_ID" \
  -H "Authorization: Bearer $TOKEN")
echo $BEFORE | jq '{
  id,
  current_cash,
  remaining_cash,
  transfer_revenue,
  remaining_transfer,
  total_revenue,
  start_cash
}'
echo ""

# Get first menu item
echo "4️⃣  Get menu items..."
MENU=$(curl -s -X GET "$API_URL/menu" \
  -H "Authorization: Bearer $TOKEN")
MENU_ITEM_ID=$(echo $MENU | jq -r '.[0].id')
MENU_ITEM_NAME=$(echo $MENU | jq -r '.[0].name')
HAS_VARIANTS=$(echo $MENU | jq -r '.[0].has_variants')

echo "Menu Item: $MENU_ITEM_NAME (ID: $MENU_ITEM_ID)"
echo "Has Variants: $HAS_VARIANTS"
echo ""

# Create order
echo "5️⃣  Create order..."
if [ "$HAS_VARIANTS" == "true" ]; then
  VARIANT_ID=$(echo $MENU | jq -r '.[0].variants[0].id')
  echo "Using variant: $VARIANT_ID"
  ORDER=$(curl -s -X POST "$API_URL/waiter/orders" \
    -H "Authorization: Bearer $TOKEN" \
    -H "Content-Type: application/json" \
    -d "{
      \"shift_id\": \"$SHIFT_ID\",
      \"customer_name\": \"Test Customer\",
      \"items\": [{
        \"menu_item_id\": \"$MENU_ITEM_ID\",
        \"variant_id\": \"$VARIANT_ID\",
        \"quantity\": 1
      }]
    }")
else
  ORDER=$(curl -s -X POST "$API_URL/waiter/orders" \
    -H "Authorization: Bearer $TOKEN" \
    -H "Content-Type: application/json" \
    -d "{
      \"shift_id\": \"$SHIFT_ID\",
      \"customer_name\": \"Test Customer\",
      \"items\": [{
        \"menu_item_id\": \"$MENU_ITEM_ID\",
        \"quantity\": 1
      }]
    }")
fi

ORDER_ID=$(echo $ORDER | jq -r '.id')
ORDER_TOTAL=$(echo $ORDER | jq -r '.total')

if [ "$ORDER_ID" == "null" ]; then
  echo "❌ Failed to create order"
  echo $ORDER | jq .
  exit 1
fi

echo "✅ Order ID: $ORDER_ID"
echo "✅ Order Total: $ORDER_TOTAL VND"
echo ""

# Pay with CASH
echo "6️⃣  Pay with CASH ($ORDER_TOTAL VND)..."
PAYMENT=$(curl -s -X POST "$API_URL/waiter/orders/$ORDER_ID/payment" \
  -H "Authorization: Bearer $TOKEN" \
  -H "Content-Type: application/json" \
  -d "{
    \"payment_method\": \"CASH\",
    \"amount\": $ORDER_TOTAL,
    \"collector_id\": \"$USER_ID\",
    \"collector_name\": \"Waiter 1\"
  }")

echo "Payment response:"
echo $PAYMENT | jq '{
  id,
  status,
  payment_method,
  amount_paid,
  total
}'
echo ""

# Wait a bit for update
sleep 1

# Get shift after
echo "7️⃣  Shift AFTER payment:"
AFTER=$(curl -s -X GET "$API_URL/waiter/shifts/$SHIFT_ID" \
  -H "Authorization: Bearer $TOKEN")
echo $AFTER | jq '{
  id,
  current_cash,
  remaining_cash,
  transfer_revenue,
  remaining_transfer,
  total_revenue,
  start_cash
}'
echo ""

# Compare
CURRENT_CASH=$(echo $AFTER | jq -r '.current_cash')
REMAINING_CASH=$(echo $AFTER | jq -r '.remaining_cash')
TOTAL_REVENUE=$(echo $AFTER | jq -r '.total_revenue')

EXPECTED_CASH=$(echo "100000 + $ORDER_TOTAL" | bc)

echo "8️⃣  Verification:"
echo "================"
echo "Expected current_cash: $EXPECTED_CASH"
echo "Actual current_cash: $CURRENT_CASH"
echo ""
echo "Expected remaining_cash: $EXPECTED_CASH"
echo "Actual remaining_cash: $REMAINING_CASH"
echo ""
echo "Expected total_revenue: $ORDER_TOTAL"
echo "Actual total_revenue: $TOTAL_REVENUE"
echo ""

if [ "$CURRENT_CASH" == "$EXPECTED_CASH" ]; then
  echo "✅ TEST PASSED: current_cash updated correctly!"
else
  echo "❌ TEST FAILED: current_cash NOT updated"
  echo ""
  echo "🔍 Check backend logs for DEBUG messages"
fi
