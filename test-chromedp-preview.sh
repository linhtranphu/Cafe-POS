#!/bin/bash

echo "Testing Chromedp HTML Preview..."
echo ""

# Get token
echo "1. Getting auth token..."
TOKEN=$(curl -s -X POST http://localhost:3000/api/login \
  -H "Content-Type: application/json" \
  -d '{"username":"admin","password":"admin123"}' | jq -r '.token')

if [ "$TOKEN" = "null" ] || [ -z "$TOKEN" ]; then
  echo "❌ Failed to get token"
  exit 1
fi
echo "✅ Token obtained"

# Get first order
echo ""
echo "2. Getting first order..."
ORDER_ID=$(curl -s -H "Authorization: Bearer $TOKEN" \
  http://localhost:3000/api/manager/orders | jq -r '.[0].id')

if [ "$ORDER_ID" = "null" ] || [ -z "$ORDER_ID" ]; then
  echo "❌ No orders found"
  exit 1
fi
echo "✅ Order ID: $ORDER_ID"

# Get shop settings to check logo
echo ""
echo "3. Checking shop settings..."
SETTINGS=$(curl -s -H "Authorization: Bearer $TOKEN" \
  http://localhost:3000/api/manager/settings)
echo "$SETTINGS" | jq '{show_logo, logo_url}'

# Create preview
echo ""
echo "4. Creating preview..."
RESULT=$(curl -s -X POST \
  -H "Authorization: Bearer $TOKEN" \
  -H "Content-Type: application/json" \
  -d "{\"order_id\":\"$ORDER_ID\"}" \
  http://localhost:3000/api/manager/html-templates/preview)

echo "$RESULT" | jq '.'

FILENAME=$(echo "$RESULT" | jq -r '.filename')

if [ "$FILENAME" != "null" ] && [ -n "$FILENAME" ]; then
  echo ""
  echo "✅ Preview created: $FILENAME"
  echo ""
  echo "5. Checking files..."
  ls -lh backend/raw_preview_*.png backend/preview_*.png 2>/dev/null | tail -2
  echo ""
  echo "6. To view the preview:"
  echo "   open backend/$FILENAME"
  echo "   open backend/raw_preview_*.png"
else
  echo "❌ Failed to create preview"
  echo "$RESULT"
fi

echo ""
echo "7. Checking backend logs for logo loading..."
tail -20 backend.log | grep -i "logo\|error\|success"
