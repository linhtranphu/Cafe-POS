#!/bin/bash

echo "🧪 Testing HTML Template API"
echo "=============================="

# Colors
GREEN='\033[0;32m'
RED='\033[0;31m'
YELLOW='\033[1;33m'
NC='\033[0m' # No Color

BASE_URL="http://localhost:8080/api/manager"

# Get auth token (assuming you have a valid token)
# Replace with your actual token or login first
TOKEN="your-token-here"

echo ""
echo "${YELLOW}1. Testing GET /html-templates/bill${NC}"
curl -s -X GET "${BASE_URL}/html-templates/bill" \
  -H "Authorization: Bearer ${TOKEN}" \
  | jq '.' || echo "${RED}❌ Failed${NC}"

echo ""
echo "${YELLOW}2. Testing POST /html-templates/preview${NC}"
echo "   (Need order_id)"
# Replace with actual order ID
ORDER_ID="675e0e8e8e8e8e8e8e8e8e8e"
curl -s -X POST "${BASE_URL}/html-templates/preview" \
  -H "Authorization: Bearer ${TOKEN}" \
  -H "Content-Type: application/json" \
  -d "{\"order_id\": \"${ORDER_ID}\"}" \
  | jq '.' || echo "${RED}❌ Failed${NC}"

echo ""
echo "${YELLOW}3. Testing POST /html-templates/test-print${NC}"
echo "   (Need order_id and printer_ip)"
curl -s -X POST "${BASE_URL}/html-templates/test-print" \
  -H "Authorization: Bearer ${TOKEN}" \
  -H "Content-Type: application/json" \
  -d "{\"order_id\": \"${ORDER_ID}\", \"printer_ip\": \"192.168.1.115\"}" \
  | jq '.' || echo "${RED}❌ Failed${NC}"

echo ""
echo "✅ Tests completed"
