#!/bin/bash

# Test script to verify ingredient and expense values are in sync
# This script creates an ingredient and checks if the expense is created correctly

API_URL="http://localhost:8080"
TOKEN=""

echo "=== Test Ingredient-Expense Sync ==="
echo ""

# Step 1: Login to get token
echo "Step 1: Login..."
LOGIN_RESPONSE=$(curl -s -X POST "$API_URL/auth/login" \
  -H "Content-Type: application/json" \
  -d '{
    "username": "admin",
    "password": "admin123"
  }')

TOKEN=$(echo $LOGIN_RESPONSE | jq -r '.token')

if [ "$TOKEN" == "null" ] || [ -z "$TOKEN" ]; then
  echo "❌ Login failed"
  echo "Response: $LOGIN_RESPONSE"
  exit 1
fi

echo "✅ Login successful"
echo ""

# Step 2: Get ingredient categories
echo "Step 2: Get ingredient categories..."
CATEGORIES=$(curl -s -X GET "$API_URL/manager/ingredient-categories" \
  -H "Authorization: Bearer $TOKEN")

CATEGORY_ID=$(echo $CATEGORIES | jq -r '.[0].id // empty')

if [ -z "$CATEGORY_ID" ]; then
  echo "⚠️  No categories found, creating one..."
  CREATE_CAT=$(curl -s -X POST "$API_URL/manager/ingredient-categories" \
    -H "Authorization: Bearer $TOKEN" \
    -H "Content-Type: application/json" \
    -d '{"name": "Test Category"}')
  CATEGORY_ID=$(echo $CREATE_CAT | jq -r '.id')
fi

echo "✅ Category ID: $CATEGORY_ID"
echo ""

# Step 3: Create ingredient with specific values
echo "Step 3: Create ingredient..."
INGREDIENT_DATA='{
  "name": "Test Coffee Beans",
  "category": "Test Category",
  "unit": "kg",
  "quantity": 10,
  "min_stock": 5,
  "cost_per_unit": 200000,
  "supplier": "Test Supplier"
}'

echo "Sending data:"
echo "$INGREDIENT_DATA" | jq '.'

INGREDIENT_RESPONSE=$(curl -s -X POST "$API_URL/manager/ingredients" \
  -H "Authorization: Bearer $TOKEN" \
  -H "Content-Type: application/json" \
  -d "$INGREDIENT_DATA")

INGREDIENT_ID=$(echo $INGREDIENT_RESPONSE | jq -r '.id // empty')

if [ -z "$INGREDIENT_ID" ]; then
  echo "❌ Failed to create ingredient"
  echo "Response: $INGREDIENT_RESPONSE"
  exit 1
fi

echo "✅ Ingredient created"
echo "Ingredient ID: $INGREDIENT_ID"
echo "Response:"
echo "$INGREDIENT_RESPONSE" | jq '.'
echo ""

# Step 4: Wait a moment for expense to be created
echo "Step 4: Waiting for expense creation..."
sleep 2
echo ""

# Step 5: Get expenses and find the one for this ingredient
echo "Step 5: Check expense..."
EXPENSES=$(curl -s -X GET "$API_URL/manager/expenses" \
  -H "Authorization: Bearer $TOKEN")

# Find expense with matching source_id
EXPENSE=$(echo $EXPENSES | jq --arg id "$INGREDIENT_ID" '.[] | select(.source_id == $id)')

if [ -z "$EXPENSE" ]; then
  echo "❌ No expense found for ingredient"
  echo "All expenses:"
  echo "$EXPENSES" | jq '.'
  exit 1
fi

echo "✅ Expense found"
echo "Expense details:"
echo "$EXPENSE" | jq '.'
echo ""

# Step 6: Verify values
echo "Step 6: Verify values..."
INGREDIENT_COST=$(echo $INGREDIENT_RESPONSE | jq -r '.cost_per_unit')
INGREDIENT_QTY=$(echo $INGREDIENT_RESPONSE | jq -r '.quantity')
EXPECTED_AMOUNT=$(echo "$INGREDIENT_COST * $INGREDIENT_QTY" | bc)

EXPENSE_AMOUNT=$(echo $EXPENSE | jq -r '.amount')

echo "Ingredient cost_per_unit: $INGREDIENT_COST VND"
echo "Ingredient quantity: $INGREDIENT_QTY kg"
echo "Expected expense amount: $EXPECTED_AMOUNT VND"
echo "Actual expense amount: $EXPENSE_AMOUNT VND"
echo ""

if [ "$EXPECTED_AMOUNT" == "$EXPENSE_AMOUNT" ]; then
  echo "✅ VALUES MATCH! Expense amount is correct."
else
  echo "❌ VALUES MISMATCH!"
  echo "   Expected: $EXPECTED_AMOUNT VND"
  echo "   Got: $EXPENSE_AMOUNT VND"
  echo "   Difference: $(echo "$EXPENSE_AMOUNT - $EXPECTED_AMOUNT" | bc) VND"
fi

echo ""
echo "=== Test Complete ==="
