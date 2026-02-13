#!/bin/bash

# Test script for variant cost analysis endpoints
# Task 6.3 - Implement cost analysis API endpoints

BASE_URL="http://localhost:8080/api"
TOKEN=""

echo "=== Testing Variant Cost Analysis Endpoints ==="
echo ""

# Login as manager to get token
echo "1. Logging in as manager..."
LOGIN_RESPONSE=$(curl -s -X POST "$BASE_URL/auth/login" \
  -H "Content-Type: application/json" \
  -d '{
    "username": "admin",
    "password": "admin123"
  }')

TOKEN=$(echo $LOGIN_RESPONSE | grep -o '"token":"[^"]*' | cut -d'"' -f4)

if [ -z "$TOKEN" ]; then
  echo "❌ Failed to login"
  exit 1
fi

echo "✅ Login successful"
echo ""

# Get a menu item ID (assuming there's at least one menu item)
echo "2. Fetching menu items..."
MENU_RESPONSE=$(curl -s -X GET "$BASE_URL/menu" \
  -H "Authorization: Bearer $TOKEN")

MENU_ITEM_ID=$(echo $MENU_RESPONSE | grep -o '"id":"[^"]*' | head -1 | cut -d'"' -f4)

if [ -z "$MENU_ITEM_ID" ]; then
  echo "❌ No menu items found"
  exit 1
fi

echo "✅ Found menu item: $MENU_ITEM_ID"
echo ""

# Test GET /api/menu/:id/cost-breakdown
echo "3. Testing GET /api/menu/$MENU_ITEM_ID/cost-breakdown"
COST_BREAKDOWN=$(curl -s -X GET "$BASE_URL/menu/$MENU_ITEM_ID/cost-breakdown" \
  -H "Authorization: Bearer $TOKEN")

echo "Response:"
echo $COST_BREAKDOWN | jq '.' 2>/dev/null || echo $COST_BREAKDOWN
echo ""

# Test GET /api/menu/:id/profit-analysis
echo "4. Testing GET /api/menu/$MENU_ITEM_ID/profit-analysis"
PROFIT_ANALYSIS=$(curl -s -X GET "$BASE_URL/menu/$MENU_ITEM_ID/profit-analysis" \
  -H "Authorization: Bearer $TOKEN")

echo "Response:"
echo $PROFIT_ANALYSIS | jq '.' 2>/dev/null || echo $PROFIT_ANALYSIS
echo ""

# Test POST /api/menu/:id/calculate-cost
echo "5. Testing POST /api/menu/$MENU_ITEM_ID/calculate-cost"
CALCULATE_COST=$(curl -s -X POST "$BASE_URL/menu/$MENU_ITEM_ID/calculate-cost" \
  -H "Authorization: Bearer $TOKEN")

echo "Response:"
echo $CALCULATE_COST | jq '.' 2>/dev/null || echo $CALCULATE_COST
echo ""

echo "=== Test Complete ==="
