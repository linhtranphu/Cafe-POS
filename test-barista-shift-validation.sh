#!/bin/bash

# Test script for BR-13: Barista must have open shift before accepting orders

BASE_URL="http://localhost:8080"

echo "=== Testing BR-13: Barista Shift Validation ==="
echo ""

# Step 1: Login as barista
echo "1. Login as barista1..."
LOGIN_RESPONSE=$(curl -s -X POST "$BASE_URL/api/auth/login" \
  -H "Content-Type: application/json" \
  -d '{
    "username": "barista1",
    "password": "barista123"
  }')

TOKEN=$(echo $LOGIN_RESPONSE | jq -r '.token')
BARISTA_ID=$(echo $LOGIN_RESPONSE | jq -r '.user.id')

if [ "$TOKEN" == "null" ] || [ -z "$TOKEN" ]; then
  echo "❌ Failed to login"
  echo $LOGIN_RESPONSE | jq .
  exit 1
fi

echo "✅ Logged in successfully"
echo "   Token: ${TOKEN:0:20}..."
echo "   Barista ID: $BARISTA_ID"
echo ""

# Step 2: Check current shift (should not exist)
echo "2. Check current shift (should not exist)..."
SHIFT_RESPONSE=$(curl -s -X GET "$BASE_URL/api/shifts/current" \
  -H "Authorization: Bearer $TOKEN")

echo "   Response: $SHIFT_RESPONSE"
echo ""

# Step 3: Get queued orders
echo "3. Get queued orders..."
QUEUE_RESPONSE=$(curl -s -X GET "$BASE_URL/api/barista/orders/queue" \
  -H "Authorization: Bearer $TOKEN")

ORDER_COUNT=$(echo $QUEUE_RESPONSE | jq '. | length')
echo "   Found $ORDER_COUNT orders in queue"

if [ "$ORDER_COUNT" -gt 0 ]; then
  ORDER_ID=$(echo $QUEUE_RESPONSE | jq -r '.[0].id')
  ORDER_NUMBER=$(echo $QUEUE_RESPONSE | jq -r '.[0].order_number')
  echo "   First order: $ORDER_NUMBER (ID: $ORDER_ID)"
else
  echo "   No orders in queue. Creating one..."
  
  # Login as waiter to create order
  WAITER_LOGIN=$(curl -s -X POST "$BASE_URL/api/auth/login" \
    -H "Content-Type: application/json" \
    -d '{"username": "waiter1", "password": "waiter123"}')
  
  WAITER_TOKEN=$(echo $WAITER_LOGIN | jq -r '.token')
  
  # Start waiter shift
  curl -s -X POST "$BASE_URL/api/shifts/start" \
    -H "Authorization: Bearer $WAITER_TOKEN" \
    -H "Content-Type: application/json" \
    -d '{"type": "MORNING", "start_cash": 1000000}' > /dev/null
  
  # Create and pay order
  CREATE_ORDER=$(curl -s -X POST "$BASE_URL/api/orders" \
    -H "Authorization: Bearer $WAITER_TOKEN" \
    -H "Content-Type: application/json" \
    -d '{
      "customer_name": "Test Customer",
      "items": [{"menu_item_id": "000000000000000000000001", "name": "Coffee", "quantity": 1, "price": 50000}],
      "total": 50000
    }')
  
  ORDER_ID=$(echo $CREATE_ORDER | jq -r '.id')
  
  # Pay order
  curl -s -X POST "$BASE_URL/api/orders/$ORDER_ID/pay" \
    -H "Authorization: Bearer $WAITER_TOKEN" \
    -H "Content-Type: application/json" \
    -d '{"payment_method": "CASH", "amount_paid": 50000}' > /dev/null
  
  # Send to bar
  curl -s -X POST "$BASE_URL/api/orders/$ORDER_ID/send-to-bar" \
    -H "Authorization: Bearer $WAITER_TOKEN" > /dev/null
  
  echo "   Created order: $ORDER_ID"
fi
echo ""

# Step 4: Try to accept order WITHOUT shift (should FAIL)
echo "4. Try to accept order WITHOUT opening shift..."
ACCEPT_RESPONSE=$(curl -s -X POST "$BASE_URL/api/barista/orders/$ORDER_ID/accept" \
  -H "Authorization: Bearer $TOKEN")

ERROR=$(echo $ACCEPT_RESPONSE | jq -r '.error')

if [ "$ERROR" == "barista must open a shift before accepting orders" ]; then
  echo "   ✅ PASS: Correctly rejected with error: $ERROR"
else
  echo "   ❌ FAIL: Expected rejection but got: $ACCEPT_RESPONSE"
  exit 1
fi
echo ""

# Step 5: Open shift
echo "5. Open barista shift..."
SHIFT_START=$(curl -s -X POST "$BASE_URL/api/shifts/start" \
  -H "Authorization: Bearer $TOKEN" \
  -H "Content-Type: application/json" \
  -d '{
    "type": "MORNING",
    "start_cash": 0
  }')

SHIFT_ID=$(echo $SHIFT_START | jq -r '.id')

if [ "$SHIFT_ID" == "null" ] || [ -z "$SHIFT_ID" ]; then
  echo "   ❌ Failed to start shift"
  echo $SHIFT_START | jq .
  exit 1
fi

echo "   ✅ Shift opened successfully"
echo "   Shift ID: $SHIFT_ID"
echo ""

# Step 6: Try to accept order WITH shift (should SUCCEED)
echo "6. Try to accept order WITH open shift..."
ACCEPT_RESPONSE2=$(curl -s -X POST "$BASE_URL/api/barista/orders/$ORDER_ID/accept" \
  -H "Authorization: Bearer $TOKEN")

ACCEPTED_STATUS=$(echo $ACCEPT_RESPONSE2 | jq -r '.status')

if [ "$ACCEPTED_STATUS" == "IN_PROGRESS" ]; then
  echo "   ✅ PASS: Order accepted successfully"
  echo "   Order status: $ACCEPTED_STATUS"
  echo "   Barista: $(echo $ACCEPT_RESPONSE2 | jq -r '.barista_name')"
else
  ERROR2=$(echo $ACCEPT_RESPONSE2 | jq -r '.error')
  echo "   ❌ FAIL: Expected success but got error: $ERROR2"
  echo "   Full response: $ACCEPT_RESPONSE2"
  exit 1
fi
echo ""

echo "=== All Tests Passed! ==="
echo ""
echo "Summary:"
echo "  ✅ Barista cannot accept orders without shift"
echo "  ✅ Barista can accept orders with open shift"
echo "  ✅ BR-13 validation is working correctly"
