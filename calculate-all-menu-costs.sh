#!/bin/bash

# Script to calculate costs for all menu items
# This will trigger cost calculation for each menu item

echo "🔄 Calculating costs for all menu items..."
echo ""

# Get backend URL from environment or use default
BACKEND_URL="${BACKEND_URL:-http://localhost:8080}"

# Get all menu items
echo "📋 Fetching menu items..."
MENU_ITEMS=$(curl -s "${BACKEND_URL}/api/menu" | jq -r '.data[].id')

if [ -z "$MENU_ITEMS" ]; then
    echo "❌ No menu items found. Please seed menu items first:"
    echo "   go run backend/cmd/seed-menu-variants/main.go"
    exit 1
fi

# Count items
ITEM_COUNT=$(echo "$MENU_ITEMS" | wc -l | tr -d ' ')
echo "✅ Found $ITEM_COUNT menu items"
echo ""

# Calculate cost for each item
COUNTER=0
SUCCESS=0
FAILED=0

for ITEM_ID in $MENU_ITEMS; do
    COUNTER=$((COUNTER + 1))
    echo "[$COUNTER/$ITEM_COUNT] Calculating cost for item: $ITEM_ID"
    
    RESPONSE=$(curl -s -X POST "${BACKEND_URL}/api/menu/${ITEM_ID}/calculate-cost")
    
    if echo "$RESPONSE" | jq -e '.cost_status' > /dev/null 2>&1; then
        COST_STATUS=$(echo "$RESPONSE" | jq -r '.cost_status')
        CURRENT_COST=$(echo "$RESPONSE" | jq -r '.current_cost')
        
        if [ "$COST_STATUS" = "FINAL" ]; then
            echo "  ✅ Status: $COST_STATUS, Cost: ${CURRENT_COST} VND"
            SUCCESS=$((SUCCESS + 1))
        elif [ "$COST_STATUS" = "INCOMPLETE" ]; then
            echo "  ⚠️  Status: $COST_STATUS (missing ingredient data)"
            SUCCESS=$((SUCCESS + 1))
        else
            echo "  ❓ Status: $COST_STATUS"
            SUCCESS=$((SUCCESS + 1))
        fi
    else
        echo "  ❌ Failed to calculate cost"
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
echo "🎉 Cost calculation complete!"
echo ""
echo "📈 Next steps:"
echo "  1. View cost breakdown: GET ${BACKEND_URL}/api/menu/:id/cost-breakdown"
echo "  2. View profit analysis: GET ${BACKEND_URL}/api/menu/:id/profit-analysis"
echo "  3. View all menu costs: GET ${BACKEND_URL}/api/menu/costs"
echo ""
