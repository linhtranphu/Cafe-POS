#!/bin/bash

# Quick test for cancel-closure API

API_URL="http://localhost:3000/api"

echo "=== Test Cancel Closure API ==="
echo ""

# Login
echo "1. Login as cashier..."
LOGIN_RESPONSE=$(curl -s -X POST "$API_URL/auth/login" \
  -H "Content-Type: application/json" \
  -d '{
    "username": "cashier1",
    "password": "password123"
  }')

TOKEN=$(echo $LOGIN_RESPONSE | jq -r '.token')

if [ "$TOKEN" == "null" ] || [ -z "$TOKEN" ]; then
  echo "❌ Login failed"
  exit 1
fi

echo "✅ Login successful"
echo ""

# Get current shift
echo "2. Get current shift..."
SHIFT=$(curl -s -X GET "$API_URL/cashier-shifts/current" \
  -H "Authorization: Bearer $TOKEN")

SHIFT_ID=$(echo $SHIFT | jq -r '.id')
STATUS=$(echo $SHIFT | jq -r '.status')

echo "   Shift ID: $SHIFT_ID"
echo "   Status: $STATUS"
echo ""

# If not in CLOSURE_INITIATED, initiate it
if [ "$STATUS" != "CLOSURE_INITIATED" ]; then
  echo "3. Initiate closure..."
  curl -s -X POST "$API_URL/cashier-shifts/$SHIFT_ID/initiate-closure" \
    -H "Authorization: Bearer $TOKEN" \
    -H "Content-Type: application/json" \
    -d '{}' | jq .
  echo ""
fi

# Test cancel-closure
echo "4. Test cancel-closure endpoint..."
CANCEL_RESPONSE=$(curl -s -X POST "$API_URL/cashier-shifts/$SHIFT_ID/cancel-closure" \
  -H "Authorization: Bearer $TOKEN" \
  -H "Content-Type: application/json" \
  -d '{}')

echo $CANCEL_RESPONSE | jq .

NEW_STATUS=$(echo $CANCEL_RESPONSE | jq -r '.status')

if [ "$NEW_STATUS" == "OPEN" ]; then
  echo ""
  echo "✅ Cancel closure successful!"
else
  echo ""
  echo "❌ Cancel closure failed"
fi
