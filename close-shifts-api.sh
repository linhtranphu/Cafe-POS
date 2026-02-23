#!/bin/bash

echo "🔄 Closing All Open Waiter Shifts via API"
echo "=========================================="
echo ""

API_URL="http://localhost:8080/api"

# Login as manager
echo "1️⃣  Logging in as manager..."
LOGIN=$(curl -s -X POST "$API_URL/auth/login" \
  -H "Content-Type: application/json" \
  -d '{
    "username": "manager1",
    "password": "password123"
  }')

TOKEN=$(echo $LOGIN | jq -r '.token')

if [ "$TOKEN" == "null" ] || [ -z "$TOKEN" ]; then
  echo "❌ Login failed. Trying admin..."
  
  LOGIN=$(curl -s -X POST "$API_URL/auth/login" \
    -H "Content-Type: application/json" \
    -d '{
      "username": "admin",
      "password": "password123"
    }')
  
  TOKEN=$(echo $LOGIN | jq -r '.token')
  
  if [ "$TOKEN" == "null" ] || [ -z "$TOKEN" ]; then
    echo "❌ Login failed"
    exit 1
  fi
fi

echo "✅ Logged in successfully"
echo ""

# Get all shifts
echo "2️⃣  Getting all shifts..."
SHIFTS=$(curl -s -X GET "$API_URL/manager/shifts" \
  -H "Authorization: Bearer $TOKEN")

# Count open shifts
OPEN_COUNT=$(echo $SHIFTS | jq '[.[] | select(.status == "OPEN")] | length')

if [ "$OPEN_COUNT" == "0" ]; then
  echo "✅ No open shifts found"
  exit 0
fi

echo "Found $OPEN_COUNT open shift(s)"
echo ""

# Get open shift IDs
SHIFT_IDS=$(echo $SHIFTS | jq -r '.[] | select(.status == "OPEN") | .id')

# Close each shift
COUNT=0
for SHIFT_ID in $SHIFT_IDS; do
  COUNT=$((COUNT + 1))
  echo "[$COUNT/$OPEN_COUNT] Closing shift: $SHIFT_ID"
  
  # Get shift details
  SHIFT=$(echo $SHIFTS | jq ".[] | select(.id == \"$SHIFT_ID\")")
  USER_NAME=$(echo $SHIFT | jq -r '.user_name')
  ROLE=$(echo $SHIFT | jq -r '.role_type')
  
  echo "   User: $USER_NAME ($ROLE)"
  
  # Close shift via API (if endpoint exists)
  # For now, we'll use direct MongoDB update
  echo "   ⚠️  API endpoint for closing shift not available"
  echo "   Please use MongoDB script or manual close"
  echo ""
done

echo "ℹ️  To close shifts, please:"
echo "   1. Install pymongo: pip3 install pymongo"
echo "   2. Run: python3 close_all_shifts.py"
echo ""
echo "Or manually in MongoDB:"
echo "   db.shifts.updateMany({status: 'OPEN'}, {\$set: {status: 'CLOSED', ended_at: new Date()}})"
