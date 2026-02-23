#!/bin/bash

# Test script to verify shift revenue realtime fix
# This tests that shift revenue fields are NOT recalculated from orders

echo "🧪 Testing Shift Revenue Realtime Fix"
echo "======================================"
echo ""

# Get auth token
echo "1️⃣ Logging in as waiter..."
LOGIN_RESPONSE=$(curl -s -X POST http://localhost:3000/api/auth/login \
  -H "Content-Type: application/json" \
  -d '{
    "username": "waiter1",
    "password": "password123"
  }')

TOKEN=$(echo $LOGIN_RESPONSE | grep -o '"token":"[^"]*' | cut -d'"' -f4)

if [ -z "$TOKEN" ]; then
  echo "❌ Login failed"
  exit 1
fi

echo "✅ Logged in successfully"
echo ""

# Start a new shift
echo "2️⃣ Starting new shift with 100k start cash..."
START_SHIFT_RESPONSE=$(curl -s -X POST http://localhost:3000/api/shifts \
  -H "Content-Type: application/json" \
  -H "Authorization: Bearer $TOKEN" \
  -d '{
    "type": "morning",
    "start_cash": 100000
  }')

SHIFT_ID=$(echo $START_SHIFT_RESPONSE | grep -o '"id":"[^"]*' | cut -d'"' -f4)

if [ -z "$SHIFT_ID" ]; then
  echo "❌ Failed to start shift"
  echo "Response: $START_SHIFT_RESPONSE"
  exit 1
fi

echo "✅ Shift started: $SHIFT_ID"
echo "   Start Cash: 100,000 VND"
echo ""

# Get current shift to verify initial values
echo "3️⃣ Getting current shift (initial state)..."
CURRENT_SHIFT=$(curl -s -X GET http://localhost:3000/api/shifts/current \
  -H "Authorization: Bearer $TOKEN")

CURRENT_CASH=$(echo $CURRENT_SHIFT | grep -o '"current_cash":[0-9]*' | cut -d':' -f2)
REMAINING_CASH=$(echo $CURRENT_SHIFT | grep -o '"remaining_cash":[0-9]*' | cut -d':' -f2)

echo "   Current Cash: $CURRENT_CASH VND"
echo "   Remaining Cash: $REMAINING_CASH VND"
echo ""

# Create an order
echo "4️⃣ Creating order..."
CREATE_ORDER_RESPONSE=$(curl -s -X POST http://localhost:3000/api/orders \
  -H "Content-Type: application/json" \
  -H "Authorization: Bearer $TOKEN" \
  -d '{
    "items": [
      {
        "product_id": "product1",
        "product_name": "Cà phê sữa",
        "quantity": 1,
        "price": 30000
      }
    ],
    "table_number": "A1"
  }')

ORDER_ID=$(echo $CREATE_ORDER_RESPONSE | grep -o '"id":"[^"]*' | cut -d'"' -f4)

if [ -z "$ORDER_ID" ]; then
  echo "❌ Failed to create order"
  exit 1
fi

echo "✅ Order created: $ORDER_ID"
echo "   Total: 30,000 VND"
echo ""

# Collect payment
echo "5️⃣ Collecting CASH payment (30k)..."
PAYMENT_RESPONSE=$(curl -s -X POST http://localhost:3000/api/orders/$ORDER_ID/collect-payment \
  -H "Content-Type: application/json" \
  -H "Authorization: Bearer $TOKEN" \
  -d '{
    "payment_method": "CASH",
    "amount": 30000
  }')

echo "✅ Payment collected"
echo ""

# Get shift again - THIS IS THE KEY TEST
echo "6️⃣ Getting current shift (after payment)..."
AFTER_PAYMENT_SHIFT=$(curl -s -X GET http://localhost:3000/api/shifts/current \
  -H "Authorization: Bearer $TOKEN")

CURRENT_CASH_AFTER=$(echo $AFTER_PAYMENT_SHIFT | grep -o '"current_cash":[0-9]*' | cut -d':' -f2)
REMAINING_CASH_AFTER=$(echo $AFTER_PAYMENT_SHIFT | grep -o '"remaining_cash":[0-9]*' | cut -d':' -f2)
TOTAL_REVENUE=$(echo $AFTER_PAYMENT_SHIFT | grep -o '"total_revenue":[0-9]*' | cut -d':' -f2)

echo "   Current Cash: $CURRENT_CASH_AFTER VND"
echo "   Remaining Cash: $REMAINING_CASH_AFTER VND"
echo "   Total Revenue: $TOTAL_REVENUE VND"
echo ""

# Verify the fix
echo "7️⃣ Verifying fix..."
EXPECTED_CURRENT_CASH=130000
EXPECTED_REMAINING_CASH=130000

if [ "$CURRENT_CASH_AFTER" = "$EXPECTED_CURRENT_CASH" ]; then
  echo "✅ Current Cash is CORRECT: $CURRENT_CASH_AFTER = $EXPECTED_CURRENT_CASH"
else
  echo "❌ Current Cash is WRONG: $CURRENT_CASH_AFTER (expected $EXPECTED_CURRENT_CASH)"
  echo "   This means CalculateTransferRevenue is still being called!"
fi

if [ "$REMAINING_CASH_AFTER" = "$EXPECTED_REMAINING_CASH" ]; then
  echo "✅ Remaining Cash is CORRECT: $REMAINING_CASH_AFTER = $EXPECTED_REMAINING_CASH"
else
  echo "❌ Remaining Cash is WRONG: $REMAINING_CASH_AFTER (expected $EXPECTED_REMAINING_CASH)"
fi

echo ""
echo "======================================"
if [ "$CURRENT_CASH_AFTER" = "$EXPECTED_CURRENT_CASH" ] && [ "$REMAINING_CASH_AFTER" = "$EXPECTED_REMAINING_CASH" ]; then
  echo "✅ TEST PASSED: Realtime values are preserved!"
else
  echo "❌ TEST FAILED: Values are being recalculated"
fi
echo "======================================"
