#!/bin/bash

BASE_URL="http://localhost:8080"

echo "=== Quick Test: Barista Shift Validation ==="
echo ""

# Login as barista
echo "1. Login as barista1..."
LOGIN=$(curl -s -X POST "$BASE_URL/api/auth/login" \
  -H "Content-Type: application/json" \
  -d '{"username": "barista1", "password": "barista123"}')

TOKEN=$(echo $LOGIN | jq -r '.token')
echo "   Token: ${TOKEN:0:30}..."
echo ""

# Check current shift
echo "2. Check current shift..."
SHIFT=$(curl -s -X GET "$BASE_URL/api/shifts/current" \
  -H "Authorization: Bearer $TOKEN")
echo "   Response: $SHIFT"
echo ""

# Try to accept an order (use the ID from your log)
ORDER_ID="697a1f99a928f9d7ff2311df"
echo "3. Try to accept order $ORDER_ID without shift..."
ACCEPT=$(curl -s -w "\nHTTP_CODE:%{http_code}" -X POST "$BASE_URL/api/barista/orders/$ORDER_ID/accept" \
  -H "Authorization: Bearer $TOKEN")

echo "$ACCEPT"
echo ""

# Extract HTTP code
HTTP_CODE=$(echo "$ACCEPT" | grep "HTTP_CODE:" | cut -d: -f2)
RESPONSE=$(echo "$ACCEPT" | grep -v "HTTP_CODE:")

echo "   HTTP Status: $HTTP_CODE"
echo "   Response Body: $RESPONSE"
echo ""

if [ "$HTTP_CODE" == "400" ]; then
  ERROR=$(echo "$RESPONSE" | jq -r '.error')
  echo "   ✅ Got 400 error as expected"
  echo "   Error message: $ERROR"
  
  if [[ "$ERROR" == *"shift"* ]]; then
    echo "   ✅ Error mentions shift - validation is working!"
  else
    echo "   ⚠️  Error doesn't mention shift: $ERROR"
  fi
else
  echo "   ❌ Expected 400, got $HTTP_CODE"
fi
