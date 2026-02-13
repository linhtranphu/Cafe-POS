#!/bin/bash

# Script to view cost analysis for menu items
# Shows cost breakdown and profit analysis

# Get backend URL from environment or use default
BACKEND_URL="${BACKEND_URL:-http://localhost:8080}"

# Colors for output
GREEN='\033[0;32m'
BLUE='\033[0;34m'
YELLOW='\033[1;33m'
RED='\033[0;31m'
NC='\033[0m' # No Color

echo "━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━"
echo "📊 MENU COST & PROFIT ANALYSIS"
echo "━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━"
echo ""

# Get all menu items with costs
echo "📋 Fetching menu items..."
MENU_DATA=$(curl -s "${BACKEND_URL}/api/menu")

if [ -z "$MENU_DATA" ]; then
    echo "❌ Failed to fetch menu items"
    exit 1
fi

# Parse menu items
MENU_ITEMS=$(echo "$MENU_DATA" | jq -r '.data[] | @base64')

if [ -z "$MENU_ITEMS" ]; then
    echo "❌ No menu items found"
    exit 1
fi

echo "✅ Found menu items"
echo ""

# Function to format currency
format_currency() {
    printf "%'d VND" "$1"
}

# Function to display single-size item
display_single_size() {
    local ITEM=$1
    local NAME=$(echo "$ITEM" | jq -r '.name')
    local CATEGORY=$(echo "$ITEM" | jq -r '.category')
    local PRICE=$(echo "$ITEM" | jq -r '.price')
    local COST=$(echo "$ITEM" | jq -r '.current_cost // 0')
    local STATUS=$(echo "$ITEM" | jq -r '.cost_status // "NOT_CALCULATED"')
    
    echo "━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━"
    echo -e "${BLUE}📦 $NAME${NC} ($CATEGORY)"
    echo "━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━"
    echo ""
    
    if [ "$STATUS" = "FINAL" ] || [ "$STATUS" = "INCOMPLETE" ]; then
        local PROFIT=$(echo "$PRICE - $COST" | bc)
        local MARGIN=$(echo "scale=1; ($PROFIT / $PRICE) * 100" | bc)
        
        echo "💰 Price:        $(format_currency $PRICE)"
        echo "💵 Cost:         $(format_currency $COST)"
        echo "📈 Profit:       $(format_currency $PROFIT)"
        echo "📊 Margin:       ${MARGIN}%"
        
        if [ "$STATUS" = "INCOMPLETE" ]; then
            echo -e "${YELLOW}⚠️  Status:      $STATUS (missing ingredient data)${NC}"
        else
            echo -e "${GREEN}✅ Status:      $STATUS${NC}"
        fi
    else
        echo "💰 Price:        $(format_currency $PRICE)"
        echo -e "${RED}❌ Cost:        Not calculated${NC}"
        echo "📊 Status:       $STATUS"
    fi
    
    echo ""
}

# Function to display multi-size item
display_multi_size() {
    local ITEM=$1
    local NAME=$(echo "$ITEM" | jq -r '.name')
    local CATEGORY=$(echo "$ITEM" | jq -r '.category')
    local VARIANTS=$(echo "$ITEM" | jq -c '.variants[]')
    
    echo "━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━"
    echo -e "${BLUE}🔄 $NAME${NC} ($CATEGORY) - MULTI-SIZE"
    echo "━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━"
    echo ""
    
    # Display table header
    printf "%-10s | %12s | %12s | %12s | %10s\n" "Variant" "Price" "Cost" "Profit" "Margin"
    echo "-----------|--------------|--------------|--------------|------------"
    
    local BEST_MARGIN=0
    local BEST_VARIANT=""
    
    while IFS= read -r VARIANT; do
        local VAR_NAME=$(echo "$VARIANT" | jq -r '.name')
        local VAR_PRICE=$(echo "$VARIANT" | jq -r '.price')
        local VAR_COST=$(echo "$VARIANT" | jq -r '.current_cost // 0')
        local VAR_STATUS=$(echo "$VARIANT" | jq -r '.cost_status // "NOT_CALCULATED"')
        
        if [ "$VAR_STATUS" = "FINAL" ] || [ "$VAR_STATUS" = "INCOMPLETE" ]; then
            local VAR_PROFIT=$(echo "$VAR_PRICE - $VAR_COST" | bc)
            local VAR_MARGIN=$(echo "scale=1; ($VAR_PROFIT / $VAR_PRICE) * 100" | bc)
            
            printf "%-10s | %10.0f đ | %10.0f đ | %10.0f đ | %9.1f%%\n" \
                "$VAR_NAME" "$VAR_PRICE" "$VAR_COST" "$VAR_PROFIT" "$VAR_MARGIN"
            
            # Track best margin
            if (( $(echo "$VAR_MARGIN > $BEST_MARGIN" | bc -l) )); then
                BEST_MARGIN=$VAR_MARGIN
                BEST_VARIANT=$VAR_NAME
            fi
        else
            printf "%-10s | %10.0f đ | %12s | %12s | %10s\n" \
                "$VAR_NAME" "$VAR_PRICE" "N/A" "N/A" "N/A"
        fi
    done <<< "$VARIANTS"
    
    echo ""
    
    if [ -n "$BEST_VARIANT" ]; then
        echo -e "${GREEN}⭐ Most profitable: $BEST_VARIANT (${BEST_MARGIN}% margin)${NC}"
    fi
    
    echo ""
}

# Process each menu item
for ITEM_BASE64 in $MENU_ITEMS; do
    ITEM=$(echo "$ITEM_BASE64" | base64 --decode)
    HAS_VARIANTS=$(echo "$ITEM" | jq -r '.has_variants')
    
    if [ "$HAS_VARIANTS" = "true" ]; then
        display_multi_size "$ITEM"
    else
        display_single_size "$ITEM"
    fi
done

echo "━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━"
echo "🎉 Analysis complete!"
echo ""
echo "💡 Tips:"
echo "  • Items with higher profit margins are more profitable"
echo "  • INCOMPLETE status means some ingredients are missing cost data"
echo "  • For multi-size items, compare margins across sizes"
echo ""
echo "📈 API Endpoints:"
echo "  • Cost breakdown: GET ${BACKEND_URL}/api/menu/:id/cost-breakdown"
echo "  • Profit analysis: GET ${BACKEND_URL}/api/menu/:id/profit-analysis"
echo "  • All menu costs: GET ${BACKEND_URL}/api/menu/costs"
echo "━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━"
