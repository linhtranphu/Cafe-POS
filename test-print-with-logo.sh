#!/bin/bash

echo "Testing Print with Logo..."
echo ""

# Get token
TOKEN=$(curl -s -X POST http://localhost:3000/api/login \
  -H "Content-Type: application/json" \
  -d '{"username":"admin","password":"admin123"}' | jq -r '.token')

if [ "$TOKEN" = "null" ] || [ -z "$TOKEN" ]; then
  echo "❌ Failed to get token"
  exit 1
fi
echo "✅ Token obtained"

# Get first order
ORDER_ID=$(curl -s -H "Authorization: Bearer $TOKEN" \
  http://localhost:3000/api/manager/orders | jq -r '.[0].id')

if [ "$ORDER_ID" = "null" ] || [ -z "$ORDER_ID" ]; then
  echo "❌ No orders found"
  exit 1
fi
echo "✅ Order ID: $ORDER_ID"

# Test print
echo ""
echo "Sending print job to 192.168.1.115..."
RESULT=$(curl -s -X POST \
  -H "Authorization: Bearer $TOKEN" \
  -H "Content-Type: application/json" \
  -d "{\"order_id\":\"$ORDER_ID\",\"printer_ip\":\"192.168.1.115\"}" \
  http://localhost:3000/api/manager/html-templates/test-print)

echo "$RESULT" | jq '.'

echo ""
echo "Check backend logs for details:"
tail -30 backend.log | grep -i "logo\|error\|success\|capture"
