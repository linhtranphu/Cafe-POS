#!/bin/bash

# Test Ingredient Creation with Stock History
# This script tests that stock history is created when creating a new ingredient

BASE_URL="http://localhost:8080/api"
TOKEN=""

echo "=== Test Ingredient Creation with Stock History ==="
echo ""

# Step 1: Login
echo "Step 1: Login as admin..."
LOGIN_RESPONSE=$(curl -s -X POST "$BASE_URL/auth/login" \
  -H "Content-Type: application/json" \
  -d '{
    "username": "admin",
    "password": "admin123"
  }')

TOKEN=$(echo $LOGIN_RESPONSE | jq -r '.token')
if [ "$TOKEN" = "null" ] || [ -z "$TOKEN" ]; then
  echo "❌ Login failed"
  echo "Response: $LOGIN_RESPONSE"
  exit 1
fi
echo "✅ Login successful"
echo ""

# Step 2: Create ingredient with initial quantity
echo "Step 2: Create ingredient with initial quantity..."
INGREDIENT_NAME="Test Ingredient $(date +%s)"
CREATE_RESPONSE=$(curl -s -X POST "$BASE_URL/manager/ingredients" \
  -H "Content-Type: application/json" \
  -H "Authorization: Bearer $TOKEN" \
  -d "{
    \"name\": \"$INGREDIENT_NAME\",
    \"category\": \"Test Category\",
    \"unit\": \"kg\",
    \"quantity\": 10,
    \"min_stock\": 2,
    \"cost_per_unit\": 50000,
    \"supplier\": \"Test Supplier\"
  }")

INGREDIENT_ID=$(echo $CREATE_RESPONSE | jq -r '.id')
if [ "$INGREDIENT_ID" = "null" ] || [ -z "$INGREDIENT_ID" ]; then
  echo "❌ Failed to create ingredient"
  echo "Response: $CREATE_RESPONSE"
  exit 1
fi
echo "✅ Ingredient created with ID: $INGREDIENT_ID"
echo ""

# Step 3: Fetch stock history
echo "Step 3: Fetch stock history..."
sleep 1  # Give backend time to process
HISTORY_RESPONSE=$(curl -s -X GET "$BASE_URL/manager/ingredients/$INGREDIENT_ID/history" \
  -H "Authorization: Bearer $TOKEN")

HISTORY_COUNT=$(echo $HISTORY_RESPONSE | jq '. | length')
echo "History records found: $HISTORY_COUNT"
echo ""

if [ "$HISTORY_COUNT" -eq 0 ]; then
  echo "❌ No stock history created!"
  echo "Response: $HISTORY_RESPONSE"
  exit 1
fi

# Step 4: Verify history details
echo "Step 4: Verify history details..."
FIRST_RECORD=$(echo $HISTORY_RESPONSE | jq '.[0]')
echo "First history record:"
echo "$FIRST_RECORD" | jq '.'
echo ""

# Check fields
QUANTITY=$(echo $FIRST_RECORD | jq -r '.quantity')
BEFORE_QTY=$(echo $FIRST_RECORD | jq -r '.before_qty')
AFTER_QTY=$(echo $FIRST_RECORD | jq -r '.after_qty')
COST_PER_UNIT=$(echo $FIRST_RECORD | jq -r '.cost_per_unit')
TOTAL_COST=$(echo $FIRST_RECORD | jq -r '.total_cost')
REASON=$(echo $FIRST_RECORD | jq -r '.reason')

echo "Verification:"
echo "  Quantity: $QUANTITY (expected: 10)"
echo "  Before Qty: $BEFORE_QTY (expected: 0)"
echo "  After Qty: $AFTER_QTY (expected: 10)"
echo "  Cost Per Unit: $COST_PER_UNIT (expected: 50000)"
echo "  Total Cost: $TOTAL_COST (expected: 500000)"
echo "  Reason: $REASON"
echo ""

# Validate
ERRORS=0
if [ "$QUANTITY" != "10" ]; then
  echo "❌ Quantity mismatch"
  ERRORS=$((ERRORS + 1))
fi
if [ "$BEFORE_QTY" != "0" ]; then
  echo "❌ Before Qty should be 0"
  ERRORS=$((ERRORS + 1))
fi
if [ "$AFTER_QTY" != "10" ]; then
  echo "❌ After Qty mismatch"
  ERRORS=$((ERRORS + 1))
fi
if [ "$COST_PER_UNIT" != "50000" ]; then
  echo "❌ Cost Per Unit mismatch"
  ERRORS=$((ERRORS + 1))
fi
if [ "$TOTAL_COST" != "500000" ]; then
  echo "❌ Total Cost mismatch"
  ERRORS=$((ERRORS + 1))
fi

if [ $ERRORS -eq 0 ]; then
  echo "✅ All validations passed!"
else
  echo "❌ $ERRORS validation(s) failed"
  exit 1
fi
echo ""

# Step 5: Test with zero quantity
echo "Step 5: Test creating ingredient with zero quantity..."
INGREDIENT_NAME_ZERO="Test Zero Ingredient $(date +%s)"
CREATE_ZERO_RESPONSE=$(curl -s -X POST "$BASE_URL/manager/ingredients" \
  -H "Content-Type: application/json" \
  -H "Authorization: Bearer $TOKEN" \
  -d "{
    \"name\": \"$INGREDIENT_NAME_ZERO\",
    \"category\": \"Test Category\",
    \"unit\": \"kg\",
    \"quantity\": 0,
    \"min_stock\": 2,
    \"cost_per_unit\": 50000,
    \"supplier\": \"Test Supplier\"
  }")

INGREDIENT_ID_ZERO=$(echo $CREATE_ZERO_RESPONSE | jq -r '.id')
if [ "$INGREDIENT_ID_ZERO" = "null" ] || [ -z "$INGREDIENT_ID_ZERO" ]; then
  echo "❌ Failed to create ingredient with zero quantity"
  exit 1
fi
echo "✅ Ingredient with zero quantity created: $INGREDIENT_ID_ZERO"

sleep 1
HISTORY_ZERO_RESPONSE=$(curl -s -X GET "$BASE_URL/manager/ingredients/$INGREDIENT_ID_ZERO/history" \
  -H "Authorization: Bearer $TOKEN")

HISTORY_ZERO_COUNT=$(echo $HISTORY_ZERO_RESPONSE | jq '. | length')
echo "History records for zero quantity: $HISTORY_ZERO_COUNT (expected: 0)"

if [ "$HISTORY_ZERO_COUNT" -eq 0 ]; then
  echo "✅ Correctly no history for zero quantity ingredient"
else
  echo "❌ Should not create history for zero quantity"
  exit 1
fi
echo ""

# Cleanup
echo "Cleanup: Deleting test ingredients..."
curl -s -X DELETE "$BASE_URL/manager/ingredients/$INGREDIENT_ID" \
  -H "Authorization: Bearer $TOKEN" > /dev/null
curl -s -X DELETE "$BASE_URL/manager/ingredients/$INGREDIENT_ID_ZERO" \
  -H "Authorization: Bearer $TOKEN" > /dev/null
echo "✅ Cleanup complete"
echo ""

echo "=== All Tests Passed! ==="
