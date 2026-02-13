#!/bin/bash

# Test script to verify variant ingredient units are saved correctly

BASE_URL="http://localhost:3000"
TOKEN=""

echo "=== Test Variant Ingredient Unit Saving ==="
echo ""

# Step 1: Login
echo "1. Logging in..."
LOGIN_RESPONSE=$(curl -s -X POST "$BASE_URL/api/auth/login" \
  -H "Content-Type: application/json" \
  -d '{
    "username": "admin",
    "password": "admin123"
  }')

TOKEN=$(echo $LOGIN_RESPONSE | grep -o '"token":"[^"]*' | cut -d'"' -f4)

if [ -z "$TOKEN" ]; then
  echo "❌ Login failed"
  echo "Response: $LOGIN_RESPONSE"
  exit 1
fi

echo "✅ Login successful"
echo ""

# Step 2: Get an ingredient to use
echo "2. Fetching ingredients..."
INGREDIENTS_RESPONSE=$(curl -s -X GET "$BASE_URL/api/manager/ingredients" \
  -H "Authorization: Bearer $TOKEN")

INGREDIENT_NAME=$(echo $INGREDIENTS_RESPONSE | jq -r '.[0].name')
INGREDIENT_UNIT=$(echo $INGREDIENTS_RESPONSE | jq -r '.[0].unit')

if [ -z "$INGREDIENT_NAME" ]; then
  echo "❌ No ingredients found. Please create some ingredients first."
  exit 1
fi

echo "✅ Found ingredient: $INGREDIENT_NAME (unit: $INGREDIENT_UNIT)"
echo ""

# Step 3: Create multi-size menu item with ingredients
echo "3. Creating multi-size menu item with ingredients..."

CREATE_RESPONSE=$(curl -s -X POST "$BASE_URL/api/manager/menu" \
  -H "Content-Type: application/json" \
  -H "Authorization: Bearer $TOKEN" \
  -d "{
    \"name\": \"Test Unit Coffee\",
    \"category\": \"Cà phê\",
    \"description\": \"Testing ingredient unit saving\",
    \"has_variants\": true,
    \"variants\": [
      {
        \"id\": \"test-m-$(date +%s)\",
        \"name\": \"Size M\",
        \"price\": 25000,
        \"ingredients\": [
          {
            \"name\": \"$INGREDIENT_NAME\",
            \"quantity\": 10,
            \"unit\": \"$INGREDIENT_UNIT\"
          }
        ],
        \"available\": true,
        \"is_default\": true
      },
      {
        \"id\": \"test-l-$(date +%s)\",
        \"name\": \"Size L\",
        \"price\": 30000,
        \"ingredients\": [
          {
            \"name\": \"$INGREDIENT_NAME\",
            \"quantity\": 15,
            \"unit\": \"$INGREDIENT_UNIT\"
          }
        ],
        \"available\": true,
        \"is_default\": false
      }
    ]
  }")

MENU_ID=$(echo $CREATE_RESPONSE | grep -o '"id":"[^"]*' | head -1 | cut -d'"' -f4)

if [ -z "$MENU_ID" ]; then
  echo "❌ Failed to create menu item"
  echo "Response: $CREATE_RESPONSE"
  exit 1
fi

echo "✅ Menu item created with ID: $MENU_ID"
echo ""

# Step 4: Verify the menu item
echo "4. Verifying menu item and ingredient units..."

GET_RESPONSE=$(curl -s -X GET "$BASE_URL/api/manager/menu/$MENU_ID" \
  -H "Authorization: Bearer $TOKEN")

echo "Menu item details:"
echo "$GET_RESPONSE" | jq '.'
echo ""

# Check variant ingredients
echo "Checking variant ingredients:"
VARIANT_1_ING_UNIT=$(echo $GET_RESPONSE | jq -r '.variants[0].ingredients[0].unit')
VARIANT_1_ING_QTY=$(echo $GET_RESPONSE | jq -r '.variants[0].ingredients[0].quantity')
VARIANT_2_ING_UNIT=$(echo $GET_RESPONSE | jq -r '.variants[1].ingredients[0].unit')
VARIANT_2_ING_QTY=$(echo $GET_RESPONSE | jq -r '.variants[1].ingredients[0].quantity')

echo "  Size M: $VARIANT_1_ING_QTY $VARIANT_1_ING_UNIT of $INGREDIENT_NAME"
echo "  Size L: $VARIANT_2_ING_QTY $VARIANT_2_ING_UNIT of $INGREDIENT_NAME"
echo ""

# Verify units are saved correctly
if [ "$VARIANT_1_ING_UNIT" = "$INGREDIENT_UNIT" ] && [ "$VARIANT_2_ING_UNIT" = "$INGREDIENT_UNIT" ]; then
  echo "✅ Ingredient units saved correctly!"
else
  echo "❌ Ingredient units NOT saved correctly"
  echo "   Expected: $INGREDIENT_UNIT"
  echo "   Got Size M: $VARIANT_1_ING_UNIT"
  echo "   Got Size L: $VARIANT_2_ING_UNIT"
  exit 1
fi

# Verify quantities
if [ "$VARIANT_1_ING_QTY" = "10" ] && [ "$VARIANT_2_ING_QTY" = "15" ]; then
  echo "✅ Ingredient quantities saved correctly!"
else
  echo "❌ Ingredient quantities NOT saved correctly"
  echo "   Expected Size M: 10, got: $VARIANT_1_ING_QTY"
  echo "   Expected Size L: 15, got: $VARIANT_2_ING_QTY"
  exit 1
fi

echo ""
echo "=== Test Complete - All checks passed! ==="
