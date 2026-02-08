#!/bin/bash

# Test Auto-Expense Logic for Gifted Ingredients
# This script tests that expenses are NOT created when no price is provided

BASE_URL="http://localhost:3000"
TOKEN=""

echo "🧪 Testing Auto-Expense Logic for Gifted Ingredients"
echo "=================================================="
echo ""

# Colors
GREEN='\033[0;32m'
RED='\033[0;31m'
YELLOW='\033[1;33m'
NC='\033[0m' # No Color

# Login as admin
echo "1️⃣  Logging in as admin..."
LOGIN_RESPONSE=$(curl -s -X POST "$BASE_URL/api/auth/login" \
  -H "Content-Type: application/json" \
  -d '{
    "username": "admin",
    "password": "admin123"
  }')

TOKEN=$(echo $LOGIN_RESPONSE | jq -r '.token')

if [ "$TOKEN" == "null" ] || [ -z "$TOKEN" ]; then
  echo -e "${RED}❌ Login failed${NC}"
  exit 1
fi

echo -e "${GREEN}✅ Login successful${NC}"
echo ""

# Get ingredient list
echo "2️⃣  Getting ingredient list..."
INGREDIENTS=$(curl -s -X GET "$BASE_URL/api/manager/ingredients" \
  -H "Authorization: Bearer $TOKEN")

INGREDIENT_ID=$(echo $INGREDIENTS | jq -r '.[0].id')
INGREDIENT_NAME=$(echo $INGREDIENTS | jq -r '.[0].name')
CURRENT_QTY=$(echo $INGREDIENTS | jq -r '.[0].quantity')
CURRENT_PRICE=$(echo $INGREDIENTS | jq -r '.[0].cost_per_unit')

echo -e "Ingredient: ${YELLOW}$INGREDIENT_NAME${NC}"
echo -e "Current Qty: ${YELLOW}$CURRENT_QTY${NC}"
echo -e "Current Price: ${YELLOW}$CURRENT_PRICE${NC}"
echo ""

# Get current expense count
echo "3️⃣  Getting current expense count..."
EXPENSES_BEFORE=$(curl -s -X GET "$BASE_URL/api/manager/expenses" \
  -H "Authorization: Bearer $TOKEN")
EXPENSE_COUNT_BEFORE=$(echo $EXPENSES_BEFORE | jq '. | length')
echo -e "Expenses before: ${YELLOW}$EXPENSE_COUNT_BEFORE${NC}"
echo ""

# Test Case 1: Adjust increase WITHOUT price (gifted) - Should NOT create expense
echo "4️⃣  Test Case 1: Adjust increase WITHOUT price (gifted)"
echo "   Expected: NO expense created"
NEW_QTY=$(echo "$CURRENT_QTY + 5" | bc)

ADJUST_RESPONSE=$(curl -s -X POST "$BASE_URL/api/manager/ingredients/$INGREDIENT_ID/adjust" \
  -H "Authorization: Bearer $TOKEN" \
  -H "Content-Type: application/json" \
  -d "{
    \"new_quantity\": $NEW_QTY,
    \"cost_per_unit\": 0,
    \"reason\": \"Test: Được tặng từ nhà cung cấp\"
  }")

sleep 1

# Check expense count
EXPENSES_AFTER=$(curl -s -X GET "$BASE_URL/api/manager/expenses" \
  -H "Authorization: Bearer $TOKEN")
EXPENSE_COUNT_AFTER=$(echo $EXPENSES_AFTER | jq '. | length')

echo -e "   Expenses after: ${YELLOW}$EXPENSE_COUNT_AFTER${NC}"

if [ "$EXPENSE_COUNT_AFTER" -eq "$EXPENSE_COUNT_BEFORE" ]; then
  echo -e "   ${GREEN}✅ PASS: No expense created (correct!)${NC}"
else
  echo -e "   ${RED}❌ FAIL: Expense was created (wrong!)${NC}"
  echo "   Latest expense:"
  echo $EXPENSES_AFTER | jq '.[0]'
fi
echo ""

# Update counts for next test
EXPENSE_COUNT_BEFORE=$EXPENSE_COUNT_AFTER
CURRENT_QTY=$NEW_QTY

# Test Case 2: Adjust increase WITH price - Should create expense
echo "5️⃣  Test Case 2: Adjust increase WITH price"
echo "   Expected: Expense created"
NEW_QTY=$(echo "$CURRENT_QTY + 3" | bc)
NEW_PRICE=50000

ADJUST_RESPONSE=$(curl -s -X POST "$BASE_URL/api/manager/ingredients/$INGREDIENT_ID/adjust" \
  -H "Authorization: Bearer $TOKEN" \
  -H "Content-Type: application/json" \
  -d "{
    \"new_quantity\": $NEW_QTY,
    \"cost_per_unit\": $NEW_PRICE,
    \"reason\": \"Test: Mua thêm với giá mới\"
  }")

sleep 1

# Check expense count
EXPENSES_AFTER=$(curl -s -X GET "$BASE_URL/api/manager/expenses" \
  -H "Authorization: Bearer $TOKEN")
EXPENSE_COUNT_AFTER=$(echo $EXPENSES_AFTER | jq '. | length')

echo -e "   Expenses after: ${YELLOW}$EXPENSE_COUNT_AFTER${NC}"

if [ "$EXPENSE_COUNT_AFTER" -gt "$EXPENSE_COUNT_BEFORE" ]; then
  echo -e "   ${GREEN}✅ PASS: Expense created (correct!)${NC}"
  echo "   Latest expense:"
  LATEST_EXPENSE=$(echo $EXPENSES_AFTER | jq '.[0]')
  EXPENSE_AMOUNT=$(echo $LATEST_EXPENSE | jq -r '.amount')
  EXPECTED_AMOUNT=$(echo "3 * $NEW_PRICE" | bc)
  echo -e "   Amount: ${YELLOW}$EXPENSE_AMOUNT${NC} (expected: $EXPECTED_AMOUNT)"
  
  if [ "$EXPENSE_AMOUNT" -eq "$EXPECTED_AMOUNT" ]; then
    echo -e "   ${GREEN}✅ Amount is correct${NC}"
  else
    echo -e "   ${RED}❌ Amount is wrong${NC}"
  fi
else
  echo -e "   ${RED}❌ FAIL: No expense created (wrong!)${NC}"
fi
echo ""

# Summary
echo "=================================================="
echo "🎯 Test Summary"
echo "=================================================="
echo ""
echo "Test Case 1: Adjust increase WITHOUT price"
echo "  Expected: NO expense"
echo "  Result: Check above"
echo ""
echo "Test Case 2: Adjust increase WITH price"
echo "  Expected: Expense created"
echo "  Result: Check above"
echo ""
echo "✅ Testing complete!"
