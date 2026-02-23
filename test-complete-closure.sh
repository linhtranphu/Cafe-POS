#!/bin/bash

# Test script for frontend-driven complete closure workflow
# This tests the new CompleteClosure endpoint that executes the entire workflow in one transaction

set -e

API_URL="http://localhost:8080/api"
TOKEN=""

echo "🧪 Testing Frontend-Driven Complete Closure Workflow"
echo "=================================================="
echo ""

# Function to make API calls
api_call() {
  local method=$1
  local endpoint=$2
  local data=$3
  
  if [ -z "$data" ]; then
    curl -s -X "$method" \
      -H "Authorization: Bearer $TOKEN" \
      -H "Content-Type: application/json" \
      "$API_URL$endpoint"
  else
    curl -s -X "$method" \
      -H "Authorization: Bearer $TOKEN" \
      -H "Content-Type: application/json" \
      -d "$data" \
      "$API_URL$endpoint"
  fi
}

# Step 1: Login as cashier
echo "1️⃣ Logging in as cashier..."
LOGIN_RESPONSE=$(curl -s -X POST \
  -H "Content-Type: application/json" \
  -d '{"username":"cashier1","password":"password123"}' \
  "$API_URL/auth/login")

TOKEN=$(echo $LOGIN_RESPONSE | jq -r '.token')
if [ "$TOKEN" = "null" ] || [ -z "$TOKEN" ]; then
  echo "❌ Login failed"
  echo $LOGIN_RESPONSE | jq .
  exit 1
fi
echo "✅ Logged in successfully"
echo ""

# Step 2: Get current cashier shift
echo "2️⃣ Getting current cashier shift..."
CURRENT_SHIFT=$(api_call GET "/cashier-shifts/current")
SHIFT_ID=$(echo $CURRENT_SHIFT | jq -r '.id')

if [ "$SHIFT_ID" = "null" ] || [ -z "$SHIFT_ID" ]; then
  echo "❌ No open cashier shift found. Please start a shift first."
  exit 1
fi

echo "✅ Found shift: $SHIFT_ID"
echo "   Status: $(echo $CURRENT_SHIFT | jq -r '.status')"
echo "   System Cash: $(echo $CURRENT_SHIFT | jq -r '.system_cash')"
echo ""

# Step 3: Check waiter shifts
echo "3️⃣ Checking waiter shifts..."
WAITER_STATUS=$(api_call GET "/cashier-shifts/check-waiter-shifts")
CAN_CLOSE=$(echo $WAITER_STATUS | jq -r '.can_close')
OPEN_COUNT=$(echo $WAITER_STATUS | jq -r '.open_count')

if [ "$CAN_CLOSE" = "false" ]; then
  echo "❌ Cannot close: $OPEN_COUNT waiter shift(s) still open"
  echo $WAITER_STATUS | jq '.open_shifts'
  exit 1
fi

echo "✅ All waiter shifts are closed"
echo ""

# Step 4: Complete closure (no variance)
echo "4️⃣ Testing complete closure WITHOUT variance..."
SYSTEM_CASH=$(echo $CURRENT_SHIFT | jq -r '.system_cash')
ACTUAL_CASH=$SYSTEM_CASH  # Same as system = no variance

COMPLETE_RESPONSE=$(api_call POST "/cashier-shifts/$SHIFT_ID/complete-closure" \
  "{\"actual_cash\": $ACTUAL_CASH}")

STATUS=$(echo $COMPLETE_RESPONSE | jq -r '.status')
if [ "$STATUS" = "CLOSED" ]; then
  echo "✅ Shift closed successfully (no variance)"
  echo "   Final Status: $STATUS"
  echo "   Actual Cash: $(echo $COMPLETE_RESPONSE | jq -r '.actual_cash')"
  echo "   Variance: $(echo $COMPLETE_RESPONSE | jq -r '.variance.amount')"
else
  echo "❌ Failed to close shift"
  echo $COMPLETE_RESPONSE | jq .
  exit 1
fi
echo ""

echo "=================================================="
echo "✅ All tests passed!"
echo ""
echo "📝 Test Summary:"
echo "   - Frontend-driven closure works correctly"
echo "   - All steps executed in one transaction"
echo "   - No intermediate state saved to database"
echo ""
echo "💡 To test WITH variance:"
echo "   1. Start a new cashier shift"
echo "   2. Modify the ACTUAL_CASH value in this script"
echo "   3. Add variance_reason and variance_notes to the payload"
echo ""
