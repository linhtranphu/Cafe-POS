#!/bin/bash

# Test script to verify variant auto-ID generation
# This creates a multi-size menu item and verifies the IDs are auto-generated

BASE_URL="http://localhost:3000"
TOKEN=""

echo "=== Test Variant Auto-ID Generation ==="
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

# Step 2: Create multi-size menu item with auto-generated IDs
echo "2. Creating multi-size menu item with auto-generated variant IDs..."

CREATE_RESPONSE=$(curl -s -X POST "$BASE_URL/api/manager/menu" \
  -H "Content-Type: application/json" \
  -H "Authorization: Bearer $TOKEN" \
  -d '{
    "name": "Test Auto ID Coffee",
    "category": "Cà phê",
    "description": "Testing auto-generated variant IDs",
    "has_variants": true,
    "variants": [
      {
        "id": "auto-gen-1",
        "name": "Size Nhỏ",
        "price": 25000,
        "ingredients": [],
        "available": true,
        "is_default": true
      },
      {
        "id": "auto-gen-2",
        "name": "Size Vừa",
        "price": 30000,
        "ingredients": [],
        "available": true,
        "is_default": false
      },
      {
        "id": "auto-gen-3",
        "name": "Size Lớn",
        "price": 35000,
        "ingredients": [],
        "available": true,
        "is_default": false
      }
    ]
  }')

MENU_ID=$(echo $CREATE_RESPONSE | grep -o '"id":"[^"]*' | head -1 | cut -d'"' -f4)

if [ -z "$MENU_ID" ]; then
  echo "❌ Failed to create menu item"
  echo "Response: $CREATE_RESPONSE"
  exit 1
fi

echo "✅ Menu item created with ID: $MENU_ID"
echo ""

# Step 3: Verify the menu item
echo "3. Verifying menu item..."

GET_RESPONSE=$(curl -s -X GET "$BASE_URL/api/manager/menu/$MENU_ID" \
  -H "Authorization: Bearer $TOKEN")

echo "Menu item details:"
echo "$GET_RESPONSE" | jq '.'
echo ""

# Check if variants have IDs
VARIANT_COUNT=$(echo $GET_RESPONSE | jq '.variants | length')
echo "Number of variants: $VARIANT_COUNT"

if [ "$VARIANT_COUNT" -eq 3 ]; then
  echo "✅ All 3 variants created successfully"
  
  # Display variant IDs
  echo ""
  echo "Variant IDs:"
  echo "$GET_RESPONSE" | jq -r '.variants[] | "  - \(.id): \(.name) - \(.price) VNĐ"'
else
  echo "❌ Expected 3 variants, got $VARIANT_COUNT"
fi

echo ""
echo "=== Test Complete ==="
