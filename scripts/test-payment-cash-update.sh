#!/bin/bash

# Test script to verify cash payment updates shift

echo "🧪 Testing Cash Payment Updates Shift"
echo "======================================"
echo ""

# Configuration
API_URL="http://localhost:3000/api"

# Step 1: Login as waiter
echo "1️⃣ Login as waiter..."
LOGIN_RESPONSE=$(curl -s -X POST "$API_URL/login" \
  -H "Content-Type: application/json" \
  -d '{"username":"waiter1","password":"password123"}')

TOKEN=$(echo $LOGIN_RESPONSE | jq -r '.token')
USER_ID=$(echo $LOGIN_RESPONSE | jq -r '.user.id')

if [ "$TOKEN" == "null" ] || [ -z "$TOKEN" ]; then
  echo "❌ Login failed!"
  echo "Response: $LOGIN_RESPONSE"
  exit 1
fi

echo "✅ Logged in successfully"
echo "Token: ${TOKEN:0:20}..."
echo ""

# Step 2: Start shift
echo "2️⃣ Starting shift..."
SHIFT_RESPONSE=$(curl -s -X POST "$API_URL/shifts/start" \
  -H "Authorization: Bearer $TOKEN" \
  -H "Content-Type: application/json" \
  -d '{"type":"MORNING","start_cash":0}')

SHIFT_ID=$(echo $SHIFT_RESPONSE | jq -r '.id')

if [ "$SHIFT_ID" == "null" ] || [ -z "$SHIFT_ID" ]; then
  echo "❌ Failed to start shift!"
  echo "Response: $SHIFT_RESPONSE"
  exit 1
fi

echo "✅ Shift started"
echo "Shift ID: $SHIFT_ID"
echo ""

# Step 3: Get menu items
echo "3️⃣ Getting menu items..."
MENU_RESPONSE=$(curl -s -X GET "$API_URL/menu" \
  -H "Authorization: Bearer $TOKEN")

FIRST_ITEM=$(echo $MENU_RESPONSE | jq -r '.[0]')
ITEM_ID=$(echo $FIRST_ITEM | jq -r '.id')
ITEM_NAME=$(echo $FIRST_ITEM | jq -r '.name')
ITEM_PRICE=$(echo $FIRST_ITEM | jq -r '.price')

echo "✅ Got menu items"
echo "First item: $ITEM_NAME ($ITEM_PRICE VND)"
echo ""

# Step 4: Create order
echo "4️⃣ Creating order..."
ORDER_RESPONSE=$(curl -s -X POST "$API_URL/orders" \
  -H "Authorization: Bearer $TOKEN" \
  -H "Content-Type: application/json" \
  -d "{
    \"table_number\": \"1\",
    \"items\": [{
      \"menu_item_id\": \"$ITEM_ID\",
      \"name\": \"$ITEM_NAME\",
      \"quantity\": 1,
      \"price\": $ITEM_PRICE
    }]
  }")

ORDER_ID=$(echo $ORDER_RESPONSE | jq -r '.id')

if [ "$ORDER_ID" == "null" ] || [ -z "$ORDER_ID" ]; then
  echo "❌ Failed to create order!"
  echo "Response: $ORDER_RESPONSE"
  exit 1
fi

echo "✅ Order created"
echo "Order ID: $ORDER_ID"
echo ""

# Step 5: Complete order
echo "5️⃣ Completing order..."
COMPLETE_RESPONSE=$(curl -s -X POST "$API_URL/orders/$ORDER_ID/complete" \
  -H "Authorization: Bearer $TOKEN")

echo "✅ Order completed"
echo ""

# Step 6: Pay with cash
echo "6️⃣ Paying with cash..."
PAYMENT_RESPONSE=$(curl -s -X POST "$API_URL/orders/$ORDER_ID/payment" \
  -H "Authorization: Bearer $TOKEN" \
  -H "Content-Type: application/json" \
  -d "{
    \"amount\": $ITEM_PRICE,
    \"payment_method\": \"cash\",
    \"collector_id\": \"$USER_ID\"
  }")

echo "Payment response: $PAYMENT_RESPONSE"
echo ""

# Step 7: Check shift cash
echo "7️⃣ Checking shift cash..."
sleep 1
SHIFT_CHECK=$(curl -s -X GET "$API_URL/shifts/current" \
  -H "Authorization: Bearer $TOKEN")

REMAINING_CASH=$(echo $SHIFT_CHECK | jq -r '.remaining_cash')
CURRENT_CASH=$(echo $SHIFT_CHECK | jq -r '.current_cash')
TOTAL_REVENUE=$(echo $SHIFT_CHECK | jq -r '.total_revenue')

echo "Shift cash status:"
echo "  Remaining cash: $REMAINING_CASH VND"
echo "  Current cash: $CURRENT_CASH VND"
echo "  Total revenue: $TOTAL_REVENUE VND"
echo ""

# Verify
if [ "$REMAINING_CASH" == "$ITEM_PRICE" ]; then
  echo "✅ SUCCESS! Cash was updated correctly!"
else
  echo "❌ FAILED! Cash was NOT updated!"
  echo "Expected: $ITEM_PRICE VND"
  echo "Got: $REMAINING_CASH VND"
  exit 1
fi

echo ""
echo "🎉 Test completed successfully!"
