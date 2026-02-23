#!/bin/bash

echo "=== Test Transfer Revenue ==="
echo ""

# Test với user waiter / waiter123
echo "Testing with waiter / waiter123..."
echo ""

# Login
echo "1. Login..."
RESPONSE=$(curl -s -X POST http://localhost:3000/api/login \
  -H "Content-Type: application/json" \
  -d '{"username":"waiter","password":"waiter123"}')

echo "Login response: $RESPONSE"
echo ""

# Extract token (simple grep method)
TOKEN=$(echo $RESPONSE | grep -o '"token":"[^"]*"' | cut -d'"' -f4)

if [ -z "$TOKEN" ]; then
  echo "❌ Login failed. Please check username/password."
  echo ""
  echo "Try these credentials:"
  echo "- waiter / waiter123"
  echo "- admin / admin123"
  exit 1
fi

echo "✅ Login successful"
echo "Token: ${TOKEN:0:30}..."
echo ""

# Get current shift
echo "2. Get current shift..."
SHIFT_RESPONSE=$(curl -s -X GET http://localhost:3000/api/waiter/shifts/current \
  -H "Authorization: Bearer $TOKEN")

echo "Shift response:"
echo "$SHIFT_RESPONSE"
echo ""

# Check if transfer_revenue exists in response
if echo "$SHIFT_RESPONSE" | grep -q "transfer_revenue"; then
  echo "✅ transfer_revenue field exists in response"
  
  # Try to extract value (basic method)
  TRANSFER_REV=$(echo "$SHIFT_RESPONSE" | grep -o '"transfer_revenue":[0-9]*' | cut -d':' -f2)
  echo "   Value: $TRANSFER_REV VND"
else
  echo "❌ transfer_revenue field NOT found in response"
fi

if echo "$SHIFT_RESPONSE" | grep -q "remaining_transfer"; then
  echo "✅ remaining_transfer field exists in response"
  
  REMAINING=$(echo "$SHIFT_RESPONSE" | grep -o '"remaining_transfer":[0-9]*' | cut -d':' -f2)
  echo "   Value: $REMAINING VND"
else
  echo "❌ remaining_transfer field NOT found in response"
fi

echo ""
echo "=== Test Complete ==="
