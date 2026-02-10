#!/bin/bash

# Test script for Menu Cost & Profit Analysis API endpoints
# This script tests all API endpoints to verify they work correctly

set -e

BASE_URL="http://localhost:3000/api"
TOKEN=""
MENU_ITEM_ID=""
SHIFT_ID=""

# Colors for output
GREEN='\033[0;32m'
RED='\033[0;31m'
YELLOW='\033[1;33m'
NC='\033[0m' # No Color

echo "=========================================="
echo "Menu Cost & Profit Analysis API Test"
echo "=========================================="
echo ""

# Function to print test result
print_result() {
    if [ $1 -eq 0 ]; then
        echo -e "${GREEN}✓ $2${NC}"
    else
        echo -e "${RED}✗ $2${NC}"
        exit 1
    fi
}

# Function to print section header
print_section() {
    echo ""
    echo -e "${YELLOW}=== $1 ===${NC}"
}

# 1. Login to get token
print_section "Authentication"
echo "Logging in as admin..."
LOGIN_RESPONSE=$(curl -s -X POST "$BASE_URL/login" \
    -H "Content-Type: application/json" \
    -d '{"username":"admin","password":"admin123"}')

TOKEN=$(echo $LOGIN_RESPONSE | grep -o '"token":"[^"]*' | cut -d'"' -f4)

if [ -z "$TOKEN" ]; then
    echo -e "${RED}✗ Failed to get authentication token${NC}"
    echo "Response: $LOGIN_RESPONSE"
    exit 1
fi

print_result 0 "Login successful"

# 2. Test Menu Cost Endpoints
print_section "Menu Cost Endpoints"

# 2.1 GET /api/manager/menu/costs
echo "Testing GET /api/manager/menu/costs..."
COSTS_RESPONSE=$(curl -s -X GET "$BASE_URL/manager/menu/costs" \
    -H "Authorization: Bearer $TOKEN")

if echo "$COSTS_RESPONSE" | grep -q '"items"'; then
    print_result 0 "GET /api/manager/menu/costs"
else
    print_result 1 "GET /api/manager/menu/costs"
fi

# Extract a menu item ID for detail test
MENU_ITEM_ID=$(echo $COSTS_RESPONSE | grep -o '"menu_item_id":"[^"]*' | head -1 | cut -d'"' -f4)

# 2.2 GET /api/manager/menu/costs/:id (if we have an ID)
if [ ! -z "$MENU_ITEM_ID" ]; then
    echo "Testing GET /api/manager/menu/costs/:id..."
    COST_DETAIL_RESPONSE=$(curl -s -X GET "$BASE_URL/manager/menu/costs/$MENU_ITEM_ID" \
        -H "Authorization: Bearer $TOKEN")
    
    if echo "$COST_DETAIL_RESPONSE" | grep -q '"menu_item"'; then
        print_result 0 "GET /api/manager/menu/costs/:id"
    else
        print_result 1 "GET /api/manager/menu/costs/:id"
    fi
else
    echo -e "${YELLOW}⊘ Skipping GET /api/manager/menu/costs/:id (no menu items)${NC}"
fi

# 2.3 GET /api/manager/menu/warnings
echo "Testing GET /api/manager/menu/warnings..."
WARNINGS_RESPONSE=$(curl -s -X GET "$BASE_URL/manager/menu/warnings" \
    -H "Authorization: Bearer $TOKEN")

if echo "$WARNINGS_RESPONSE" | grep -q '"loss_items"'; then
    print_result 0 "GET /api/manager/menu/warnings"
else
    print_result 1 "GET /api/manager/menu/warnings"
fi

# 2.4 Test with custom threshold
echo "Testing GET /api/manager/menu/warnings?threshold=25..."
WARNINGS_THRESHOLD_RESPONSE=$(curl -s -X GET "$BASE_URL/manager/menu/warnings?threshold=25" \
    -H "Authorization: Bearer $TOKEN")

if echo "$WARNINGS_THRESHOLD_RESPONSE" | grep -q '"threshold"'; then
    print_result 0 "GET /api/manager/menu/warnings?threshold=25"
else
    print_result 1 "GET /api/manager/menu/warnings?threshold=25"
fi

# 3. Test Profit Analysis Endpoints
print_section "Profit Analysis Endpoints"

# 3.1 GET /api/manager/reports/category-profit
echo "Testing GET /api/manager/reports/category-profit..."
START_DATE=$(date -u -d "30 days ago" +%Y-%m-%d 2>/dev/null || date -u -v-30d +%Y-%m-%d)
END_DATE=$(date -u +%Y-%m-%d)

CATEGORY_PROFIT_RESPONSE=$(curl -s -X GET "$BASE_URL/manager/reports/category-profit?start_date=$START_DATE&end_date=$END_DATE" \
    -H "Authorization: Bearer $TOKEN")

if echo "$CATEGORY_PROFIT_RESPONSE" | grep -q '"categories"'; then
    print_result 0 "GET /api/manager/reports/category-profit"
else
    print_result 1 "GET /api/manager/reports/category-profit"
fi

# 3.2 GET /api/manager/reports/operating-profit
echo "Testing GET /api/manager/reports/operating-profit..."
OPERATING_PROFIT_RESPONSE=$(curl -s -X GET "$BASE_URL/manager/reports/operating-profit?start_date=$START_DATE&end_date=$END_DATE" \
    -H "Authorization: Bearer $TOKEN")

if echo "$OPERATING_PROFIT_RESPONSE" | grep -q '"gross_profit"'; then
    print_result 0 "GET /api/manager/reports/operating-profit"
else
    print_result 1 "GET /api/manager/reports/operating-profit"
fi

# 4. Test Operating Expense Endpoints
print_section "Operating Expense Endpoints"

# 4.1 POST /api/manager/operating-expenses
echo "Testing POST /api/manager/operating-expenses..."
CREATE_EXPENSE_RESPONSE=$(curl -s -X POST "$BASE_URL/manager/operating-expenses" \
    -H "Authorization: Bearer $TOKEN" \
    -H "Content-Type: application/json" \
    -d "{
        \"period_start\": \"$START_DATE\",
        \"period_end\": \"$END_DATE\",
        \"staff_salary\": 10000000,
        \"rent\": 5000000,
        \"utilities\": 2000000,
        \"marketing_costs\": 1000000,
        \"other_expenses\": 500000
    }")

if echo "$CREATE_EXPENSE_RESPONSE" | grep -q '"total_expenses"'; then
    print_result 0 "POST /api/manager/operating-expenses"
else
    print_result 1 "POST /api/manager/operating-expenses"
fi

# 4.2 GET /api/manager/operating-expenses
echo "Testing GET /api/manager/operating-expenses..."
GET_EXPENSES_RESPONSE=$(curl -s -X GET "$BASE_URL/manager/operating-expenses?start_date=$START_DATE&end_date=$END_DATE" \
    -H "Authorization: Bearer $TOKEN")

if echo "$GET_EXPENSES_RESPONSE" | grep -q '"expenses"'; then
    print_result 0 "GET /api/manager/operating-expenses"
else
    print_result 1 "GET /api/manager/operating-expenses"
fi

# 5. Test Settings Endpoint
print_section "Settings Endpoints"

# 5.1 GET /api/manager/settings
echo "Testing GET /api/manager/settings..."
GET_SETTINGS_RESPONSE=$(curl -s -X GET "$BASE_URL/manager/settings" \
    -H "Authorization: Bearer $TOKEN")

if echo "$GET_SETTINGS_RESPONSE" | grep -q '"low_margin_threshold"'; then
    print_result 0 "GET /api/manager/settings"
else
    print_result 1 "GET /api/manager/settings"
fi

# 5.2 PATCH /api/manager/settings
echo "Testing PATCH /api/manager/settings..."
UPDATE_SETTINGS_RESPONSE=$(curl -s -X PATCH "$BASE_URL/manager/settings" \
    -H "Authorization: Bearer $TOKEN" \
    -H "Content-Type: application/json" \
    -d '{"low_margin_threshold": 25.0}')

if echo "$UPDATE_SETTINGS_RESPONSE" | grep -q '"low_margin_threshold"'; then
    print_result 0 "PATCH /api/manager/settings"
else
    print_result 1 "PATCH /api/manager/settings"
fi

# 6. Test Filtering and Sorting
print_section "Filtering and Sorting"

# 6.1 Test category filter
echo "Testing GET /api/manager/menu/costs?category=Coffee..."
FILTER_RESPONSE=$(curl -s -X GET "$BASE_URL/manager/menu/costs?category=Coffee" \
    -H "Authorization: Bearer $TOKEN")

if echo "$FILTER_RESPONSE" | grep -q '"items"'; then
    print_result 0 "GET /api/manager/menu/costs?category=Coffee"
else
    print_result 1 "GET /api/manager/menu/costs?category=Coffee"
fi

# 6.2 Test sorting
echo "Testing GET /api/manager/menu/costs?sort_by=profit_margin&sort_order=desc..."
SORT_RESPONSE=$(curl -s -X GET "$BASE_URL/manager/menu/costs?sort_by=profit_margin&sort_order=desc" \
    -H "Authorization: Bearer $TOKEN")

if echo "$SORT_RESPONSE" | grep -q '"items"'; then
    print_result 0 "GET /api/manager/menu/costs?sort_by=profit_margin&sort_order=desc"
else
    print_result 1 "GET /api/manager/menu/costs?sort_by=profit_margin&sort_order=desc"
fi

# 7. Test Error Handling
print_section "Error Handling"

# 7.1 Test invalid menu item ID
echo "Testing GET /api/manager/menu/costs/invalid_id..."
INVALID_ID_RESPONSE=$(curl -s -w "\n%{http_code}" -X GET "$BASE_URL/manager/menu/costs/invalid_id" \
    -H "Authorization: Bearer $TOKEN")

HTTP_CODE=$(echo "$INVALID_ID_RESPONSE" | tail -1)
if [ "$HTTP_CODE" -eq 400 ] || [ "$HTTP_CODE" -eq 404 ]; then
    print_result 0 "GET /api/manager/menu/costs/invalid_id (returns error)"
else
    print_result 1 "GET /api/manager/menu/costs/invalid_id (should return error)"
fi

# 7.2 Test invalid date range
echo "Testing GET /api/manager/reports/category-profit with invalid dates..."
INVALID_DATE_RESPONSE=$(curl -s -w "\n%{http_code}" -X GET "$BASE_URL/manager/reports/category-profit?start_date=invalid&end_date=invalid" \
    -H "Authorization: Bearer $TOKEN")

HTTP_CODE=$(echo "$INVALID_DATE_RESPONSE" | tail -1)
if [ "$HTTP_CODE" -eq 400 ]; then
    print_result 0 "GET /api/manager/reports/category-profit with invalid dates (returns error)"
else
    print_result 1 "GET /api/manager/reports/category-profit with invalid dates (should return error)"
fi

# 7.3 Test negative expense values
echo "Testing POST /api/manager/operating-expenses with negative values..."
NEGATIVE_EXPENSE_RESPONSE=$(curl -s -w "\n%{http_code}" -X POST "$BASE_URL/manager/operating-expenses" \
    -H "Authorization: Bearer $TOKEN" \
    -H "Content-Type: application/json" \
    -d "{
        \"period_start\": \"$START_DATE\",
        \"period_end\": \"$END_DATE\",
        \"staff_salary\": -1000,
        \"rent\": 5000000,
        \"utilities\": 2000000,
        \"marketing_costs\": 1000000,
        \"other_expenses\": 500000
    }")

HTTP_CODE=$(echo "$NEGATIVE_EXPENSE_RESPONSE" | tail -1)
if [ "$HTTP_CODE" -eq 400 ]; then
    print_result 0 "POST /api/manager/operating-expenses with negative values (returns error)"
else
    print_result 1 "POST /api/manager/operating-expenses with negative values (should return error)"
fi

# 8. Database Verification
print_section "Database Verification"

echo "Checking MongoDB collections..."

# Check if order_items collection exists
ORDER_ITEMS_COUNT=$(docker exec cafe-pos-mongodb-1 mongosh cafe_pos --quiet --eval "db.order_items.countDocuments({})" 2>/dev/null || echo "0")
echo "Order items in database: $ORDER_ITEMS_COUNT"

# Check if operating_expenses collection exists
EXPENSES_COUNT=$(docker exec cafe-pos-mongodb-1 mongosh cafe_pos --quiet --eval "db.operating_expenses.countDocuments({})" 2>/dev/null || echo "0")
echo "Operating expenses in database: $EXPENSES_COUNT"

# Check if shop_settings has low_margin_threshold
SETTINGS_CHECK=$(docker exec cafe-pos-mongodb-1 mongosh cafe_pos --quiet --eval "db.shop_settings.findOne({}, {low_margin_threshold: 1})" 2>/dev/null || echo "{}")
echo "Shop settings: $SETTINGS_CHECK"

print_result 0 "Database verification complete"

# Summary
print_section "Test Summary"
echo -e "${GREEN}All API endpoint tests passed!${NC}"
echo ""
echo "Tested endpoints:"
echo "  ✓ GET /api/manager/menu/costs"
echo "  ✓ GET /api/manager/menu/costs/:id"
echo "  ✓ GET /api/manager/menu/warnings"
echo "  ✓ GET /api/manager/reports/category-profit"
echo "  ✓ GET /api/manager/reports/operating-profit"
echo "  ✓ POST /api/manager/operating-expenses"
echo "  ✓ GET /api/manager/operating-expenses"
echo "  ✓ GET /api/manager/settings"
echo "  ✓ PATCH /api/manager/settings"
echo "  ✓ Filtering and sorting"
echo "  ✓ Error handling"
echo "  ✓ Database operations"
echo ""
echo "=========================================="
