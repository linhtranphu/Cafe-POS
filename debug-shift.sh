#!/bin/bash

echo "=== Debug Shift Data ==="
echo ""

# Login as waiter
TOKEN=$(curl -s -X POST http://localhost:3000/api/login \
  -H "Content-Type: application/json" \
  -d '{"username":"waiter","password":"waiter123"}' | grep -o '"token":"[^"]*"' | cut -d'"' -f4)

if [ -z "$TOKEN" ]; then
  echo "❌ Login failed"
  exit 1
fi

echo "✅ Logged in"
echo ""

# Get current shift
echo "=== Current Shift ==="
curl -s -X GET http://localhost:3000/api/waiter/shifts/current \
  -H "Authorization: Bearer $TOKEN" | python3 -m json.tool 2>/dev/null || echo "No shift or invalid JSON"

echo ""
echo ""

# Get orders for the shift (if shift exists)
SHIFT_ID=$(curl -s -X GET http://localhost:3000/api/waiter/shifts/current \
  -H "Authorization: Bearer $TOKEN" | grep -o '"id":"[^"]*"' | head -1 | cut -d'"' -f4)

if [ ! -z "$SHIFT_ID" ]; then
  echo "=== Orders for Shift $SHIFT_ID ==="
  curl -s -X GET "http://localhost:3000/api/waiter/orders?shift_id=$SHIFT_ID" \
    -H "Authorization: Bearer $TOKEN" | python3 -m json.tool 2>/dev/null || echo "No orders or invalid JSON"
fi
