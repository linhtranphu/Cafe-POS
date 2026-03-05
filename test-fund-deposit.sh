#!/bin/bash

# Test fund deposit API
# Make sure backend is running on port 3000

echo "Testing Fund Deposit API..."
echo ""

# Get token first (replace with your actual credentials)
echo "1. Login to get token..."
LOGIN_RESPONSE=$(curl -s -X POST http://localhost:3000/api/auth/login \
  -H "Content-Type: application/json" \
  -d '{
    "username": "manager",
    "password": "manager123"
  }')

TOKEN=$(echo $LOGIN_RESPONSE | grep -o '"token":"[^"]*' | cut -d'"' -f4)

if [ -z "$TOKEN" ]; then
  echo "❌ Failed to get token. Response:"
  echo $LOGIN_RESPONSE
  exit 1
fi

echo "✅ Got token: ${TOKEN:0:20}..."
echo ""

# Test deposit
echo "2. Testing deposit..."
DEPOSIT_RESPONSE=$(curl -s -w "\nHTTP_CODE:%{http_code}" -X POST http://localhost:3000/api/manager/fund/deposit \
  -H "Content-Type: application/json" \
  -H "Authorization: Bearer $TOKEN" \
  -d '{
    "cash_amount": 100000,
    "transfer_amount": 0,
    "reason": "Test deposit from script"
  }')

HTTP_CODE=$(echo "$DEPOSIT_RESPONSE" | grep "HTTP_CODE:" | cut -d: -f2)
BODY=$(echo "$DEPOSIT_RESPONSE" | sed '/HTTP_CODE:/d')

echo "HTTP Status: $HTTP_CODE"
echo "Response Body:"
echo "$BODY" | jq '.' 2>/dev/null || echo "$BODY"

if [ "$HTTP_CODE" = "201" ]; then
  echo ""
  echo "✅ Deposit successful!"
else
  echo ""
  echo "❌ Deposit failed with status $HTTP_CODE"
fi
