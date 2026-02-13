#!/bin/bash

# Simple script to calculate costs without jq dependency
# Supports authentication

echo "🔄 Calculating costs for all menu items..."
echo ""

# Get backend URL from environment or use default (port 3000)
BACKEND_URL="${BACKEND_URL:-http://localhost:3000}"

# Get auth token
if [ -f ".auth_token" ]; then
    AUTH_TOKEN=$(cat .auth_token)
    echo "🔐 Using saved auth token"
else
    echo "🔐 Getting auth token..."
    ./get-auth-token.sh > /dev/null 2>&1
    if [ -f ".auth_token" ]; then
        AUTH_TOKEN=$(cat .auth_token)
    else
        echo "❌ Failed to get auth token"
        echo "   Run: ./get-auth-token.sh"
        exit 1
    fi
fi

echo ""

# Get all menu items
echo "📋 Fetching menu items..."
RESPONSE=$(curl -s "${BACKEND_URL}/api/manager/menu" \
  -H "Authorization: Bearer ${AUTH_TOKEN}")

# Check if we got a response
if [ -z "$RESPONSE" ]; then
    echo "❌ Failed to connect to backend at ${BACKEND_URL}"
    echo "   Make sure the backend is running!"
    exit 1
fi

# Check for auth error
if echo "$RESPONSE" | grep -q "authorization"; then
    echo "❌ Authentication failed"
    echo "   Getting new token..."
    rm -f .auth_token
    ./get-auth-token.sh
    if [ -f ".auth_token" ]; then
        AUTH_TOKEN=$(cat .auth_token)
        RESPONSE=$(curl -s "${BACKEND_URL}/api/manager/menu" \
          -H "Authorization: Bearer ${AUTH_TOKEN}")
    else
        exit 1
    fi
fi

# Extract menu item IDs using grep and sed (no jq needed)
MENU_IDS=$(echo "$RESPONSE" | grep -o '"id":"[^"]*"' | sed 's/"id":"//g' | sed 's/"//g')

if [ -z "$MENU_IDS" ]; then
    echo "❌ No menu items found"
    echo "   Response: $RESPONSE"
    exit 1
fi

# Count items
ITEM_COUNT=$(echo "$MENU_IDS" | wc -l | tr -d ' ')
echo "✅ Found $ITEM_COUNT menu items"
echo ""

# Calculate cost for each item
COUNTER=0
SUCCESS=0
FAILED=0

for ITEM_ID in $MENU_IDS; do
    COUNTER=$((COUNTER + 1))
    echo "[$COUNTER/$ITEM_COUNT] Calculating cost for: $ITEM_ID"
    
    CALC_RESPONSE=$(curl -s -X POST "${BACKEND_URL}/api/manager/menu/${ITEM_ID}/calculate-cost" \
      -H "Authorization: Bearer ${AUTH_TOKEN}")
    
    # Check if response contains cost_status
    if echo "$CALC_RESPONSE" | grep -q "cost_status"; then
        # Extract status and cost
        STATUS=$(echo "$CALC_RESPONSE" | grep -o '"cost_status":"[^"]*"' | sed 's/"cost_status":"//g' | sed 's/"//g')
        COST=$(echo "$CALC_RESPONSE" | grep -o '"current_cost":[0-9.]*' | sed 's/"current_cost"://g')
        
        if [ "$STATUS" = "FINAL" ]; then
            echo "  ✅ Status: $STATUS, Cost: ${COST} VND"
            SUCCESS=$((SUCCESS + 1))
        elif [ "$STATUS" = "INCOMPLETE" ]; then
            echo "  ⚠️  Status: $STATUS (missing ingredient data)"
            SUCCESS=$((SUCCESS + 1))
        else
            echo "  ❓ Status: $STATUS"
            SUCCESS=$((SUCCESS + 1))
        fi
    else
        echo "  ❌ Failed to calculate cost"
        echo "     Response: $CALC_RESPONSE"
        FAILED=$((FAILED + 1))
    fi
    echo ""
done

echo "━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━"
echo "📊 Summary:"
echo "  Total items: $ITEM_COUNT"
echo "  ✅ Success: $SUCCESS"
echo "  ❌ Failed: $FAILED"
echo "━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━"
echo ""

if [ $SUCCESS -gt 0 ]; then
    echo "🎉 Cost calculation complete!"
    echo ""
    echo "📈 View results:"
    echo "  • GET ${BACKEND_URL}/api/manager/menu"
    echo "  • GET ${BACKEND_URL}/api/manager/menu/:id/cost-breakdown"
    echo "  • GET ${BACKEND_URL}/api/manager/menu/:id/profit-analysis"
else
    echo "⚠️  No costs calculated successfully"
    echo "   Check if ingredients are seeded: go run backend/cmd/seed/main.go"
fi
echo ""
