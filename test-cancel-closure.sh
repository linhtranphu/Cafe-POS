#!/bin/bash

# Test Cancel Closure Flow
# This script tests the cancel closure functionality

set -e

API_URL="http://localhost:8080/api"
TOKEN=""

echo "=== Test Cancel Closure Flow ==="
echo ""

# Step 1: Login as cashier
echo "1. Login as cashier..."
LOGIN_RESPONSE=$(curl -s -X POST "$API_URL/auth/login" \
  -H "Content-Type: application/json" \
  -d '{
    "username": "cashier1",
    "password": "password123"
  }')

TOKEN=$(echo $LOGIN_RESPONSE | jq -r '.token')
USER_ID=$(echo $LOGIN_RESPONSE | jq -r '.user.id')

if [ "$TOKEN" == "null" ] || [ -z "$TOKEN" ]; then
  echo "❌ Login failed"
  echo $LOGIN_RESPONSE | jq .
  exit 1
fi

echo "✅ Login successful"
echo "   User ID: $USER_ID"
echo ""

# Step 2: Get current cashier shift
echo "2. Get current cashier shift..."
CURRENT_SHIFT=$(curl -s -X GET "$API_URL/cashier-shifts/current" \
  -H "Authorization: Bearer $TOKEN")

SHIFT_ID=$(echo $CURRENT_SHIFT | jq -r '.id')

if [ "$SHIFT_ID" == "null" ] || [ -z "$SHIFT_ID" ]; then
  echo "❌ No open cashier shift found"
  echo "   Creating a new shift..."
  
  # Create new shift
  NEW_SHIFT=$(curl -s -X POST "$API_URL/cashier-shifts" \
    -H "Authorization: Bearer $TOKEN" \
    -H "Content-Type: application/json" \
    -d '{
      "starting_float": 500000
    }')
  
  SHIFT_ID=$(echo $NEW_SHIFT | jq -r '.id')
  
  if [ "$SHIFT_ID" == "null" ] || [ -z "$SHIFT_ID" ]; then
    echo "❌ Failed to create shift"
    echo $NEW_SHIFT | jq .
    exit 1
  fi
  
  echo "✅ New shift created: $SHIFT_ID"
else
  echo "✅ Found open shift: $SHIFT_ID"
fi

SHIFT_STATUS=$(echo $CURRENT_SHIFT | jq -r '.status')
echo "   Status: $SHIFT_STATUS"
echo ""

# Step 3: Initiate closure
echo "3. Initiate closure..."
INITIATE_RESPONSE=$(curl -s -X POST "$API_URL/cashier-shifts/$SHIFT_ID/initiate-closure" \
  -H "Authorization: Bearer $TOKEN" \
  -H "Content-Type: application/json" \
  -d '{}')

NEW_STATUS=$(echo $INITIATE_RESPONSE | jq -r '.status')

if [ "$NEW_STATUS" != "CLOSURE_INITIATED" ]; then
  echo "❌ Failed to initiate closure"
  echo $INITIATE_RESPONSE | jq .
  exit 1
fi

echo "✅ Closure initiated"
echo "   Status: $NEW_STATUS"
echo ""

# Step 4: Cancel closure
echo "4. Cancel closure..."
CANCEL_RESPONSE=$(curl -s -X POST "$API_URL/cashier-shifts/$SHIFT_ID/cancel-closure" \
  -H "Authorization: Bearer $TOKEN" \
  -H "Content-Type: application/json" \
  -d '{}')

FINAL_STATUS=$(echo $CANCEL_RESPONSE | jq -r '.status')

if [ "$FINAL_STATUS" != "OPEN" ]; then
  echo "❌ Failed to cancel closure"
  echo $CANCEL_RESPONSE | jq .
  exit 1
fi

echo "✅ Closure cancelled successfully"
echo "   Status: $FINAL_STATUS"
echo ""

# Step 5: Verify audit log
echo "5. Verify audit log..."
AUDIT_LOG=$(echo $CANCEL_RESPONSE | jq '.audit_log')
LAST_ACTION=$(echo $AUDIT_LOG | jq -r '.[-1].action')

if [ "$LAST_ACTION" != "closure_cancelled" ]; then
  echo "❌ Audit log not updated correctly"
  echo "   Last action: $LAST_ACTION"
  exit 1
fi

echo "✅ Audit log updated correctly"
echo "   Last action: $LAST_ACTION"
echo ""

# Step 6: Test cannot cancel after recording actual cash
echo "6. Test cannot cancel after recording actual cash..."

# Initiate closure again
echo "   6.1. Initiate closure again..."
curl -s -X POST "$API_URL/cashier-shifts/$SHIFT_ID/initiate-closure" \
  -H "Authorization: Bearer $TOKEN" \
  -H "Content-Type: application/json" \
  -d '{}' > /dev/null

echo "   ✅ Closure initiated"

# Record actual cash
echo "   6.2. Record actual cash..."
curl -s -X POST "$API_URL/cashier-shifts/$SHIFT_ID/record-actual-cash" \
  -H "Authorization: Bearer $TOKEN" \
  -H "Content-Type: application/json" \
  -d '{
    "actual_cash": 500000
  }' > /dev/null

echo "   ✅ Actual cash recorded"

# Try to cancel (should fail)
echo "   6.3. Try to cancel (should fail)..."
CANCEL_FAIL_RESPONSE=$(curl -s -X POST "$API_URL/cashier-shifts/$SHIFT_ID/cancel-closure" \
  -H "Authorization: Bearer $TOKEN" \
  -H "Content-Type: application/json" \
  -d '{}')

ERROR_MSG=$(echo $CANCEL_FAIL_RESPONSE | jq -r '.error')

if [[ "$ERROR_MSG" == *"actual cash has been recorded"* ]]; then
  echo "   ✅ Cancel correctly rejected: $ERROR_MSG"
else
  echo "   ❌ Cancel should have been rejected"
  echo $CANCEL_FAIL_RESPONSE | jq .
  exit 1
fi

echo ""
echo "=== All Tests Passed! ==="
echo ""
echo "Summary:"
echo "✅ Can initiate closure"
echo "✅ Can cancel closure before recording actual cash"
echo "✅ Audit log is updated correctly"
echo "✅ Cannot cancel after recording actual cash"
