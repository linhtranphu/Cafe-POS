#!/bin/bash

# Test Print Bridge URL Fix

GREEN='\033[0;32m'
YELLOW='\033[1;33m'
RED='\033[0;31m'
BLUE='\033[0;34m'
NC='\033[0m'

echo -e "${BLUE}╔════════════════════════════════════════════════════════╗${NC}"
echo -e "${BLUE}║   Testing Print Bridge URL Save Fix                   ║${NC}"
echo -e "${BLUE}╚════════════════════════════════════════════════════════╝${NC}"
echo ""

# Configuration
API_URL="${API_URL:-https://tacafe.store}"
TOKEN="${TOKEN:-}"

if [ -z "$TOKEN" ]; then
    echo -e "${YELLOW}Please provide your auth token:${NC}"
    echo "export TOKEN='your_token_here'"
    echo ""
    echo "Or run with:"
    echo "TOKEN='your_token' ./test-print-bridge-url-fix.sh"
    exit 1
fi

echo -e "${YELLOW}Step 1: Get current settings${NC}"
echo "GET ${API_URL}/api/manager/shop-settings"
echo ""

response=$(curl -s -w "\n%{http_code}" \
    -H "Authorization: Bearer ${TOKEN}" \
    "${API_URL}/api/manager/shop-settings")

http_code=$(echo "$response" | tail -n1)
body=$(echo "$response" | head -n-1)

if [ "$http_code" != "200" ]; then
    echo -e "${RED}❌ Failed to get settings${NC}"
    echo "HTTP Code: $http_code"
    echo "Response: $body"
    exit 1
fi

echo -e "${GREEN}✅ Got settings${NC}"
echo ""

# Parse settings
settings_id=$(echo "$body" | jq -r '.id')
current_url=$(echo "$body" | jq -r '.print_bridge_url // "not set"')

echo "Settings ID: $settings_id"
echo "Current Print Bridge URL: $current_url"
echo ""

# Update with new URL
TEST_URL="http://192.168.1.100:3001"

echo -e "${YELLOW}Step 2: Update Print Bridge URL${NC}"
echo "PUT ${API_URL}/api/manager/shop-settings/${settings_id}"
echo "New URL: ${TEST_URL}"
echo ""

# Get all current values
shop_name=$(echo "$body" | jq -r '.shop_name')
shop_address=$(echo "$body" | jq -r '.shop_address // ""')
shop_phone=$(echo "$body" | jq -r '.shop_phone // ""')
logo_url=$(echo "$body" | jq -r '.logo_url // ""')
custom_message=$(echo "$body" | jq -r '.custom_message // ""')
show_logo=$(echo "$body" | jq -r '.show_logo // true')
show_address=$(echo "$body" | jq -r '.show_address // true')
show_phone=$(echo "$body" | jq -r '.show_phone // true')
show_custom_message=$(echo "$body" | jq -r '.show_custom_message // true')
auto_print_enabled=$(echo "$body" | jq -r '.auto_print_enabled // true')

update_response=$(curl -s -w "\n%{http_code}" \
    -X PUT \
    -H "Authorization: Bearer ${TOKEN}" \
    -H "Content-Type: application/json" \
    -d "{
        \"shop_name\": \"${shop_name}\",
        \"shop_address\": \"${shop_address}\",
        \"shop_phone\": \"${shop_phone}\",
        \"logo_url\": \"${logo_url}\",
        \"custom_message\": \"${custom_message}\",
        \"print_bridge_url\": \"${TEST_URL}\",
        \"show_logo\": ${show_logo},
        \"show_address\": ${show_address},
        \"show_phone\": ${show_phone},
        \"show_custom_message\": ${show_custom_message},
        \"auto_print_enabled\": ${auto_print_enabled}
    }" \
    "${API_URL}/api/manager/shop-settings/${settings_id}")

update_http_code=$(echo "$update_response" | tail -n1)
update_body=$(echo "$update_response" | head -n-1)

if [ "$update_http_code" != "200" ]; then
    echo -e "${RED}❌ Failed to update settings${NC}"
    echo "HTTP Code: $update_http_code"
    echo "Response: $update_body"
    exit 1
fi

echo -e "${GREEN}✅ Settings updated${NC}"
echo ""

# Verify the update
echo -e "${YELLOW}Step 3: Verify Print Bridge URL was saved${NC}"
echo "GET ${API_URL}/api/manager/shop-settings"
echo ""

verify_response=$(curl -s -w "\n%{http_code}" \
    -H "Authorization: Bearer ${TOKEN}" \
    "${API_URL}/api/manager/shop-settings")

verify_http_code=$(echo "$verify_response" | tail -n1)
verify_body=$(echo "$verify_response" | head -n-1)

if [ "$verify_http_code" != "200" ]; then
    echo -e "${RED}❌ Failed to verify settings${NC}"
    echo "HTTP Code: $verify_http_code"
    exit 1
fi

saved_url=$(echo "$verify_body" | jq -r '.print_bridge_url // "not set"')

echo "Saved Print Bridge URL: $saved_url"
echo ""

if [ "$saved_url" = "$TEST_URL" ]; then
    echo -e "${GREEN}╔════════════════════════════════════════════════════════╗${NC}"
    echo -e "${GREEN}║   ✅ SUCCESS! Print Bridge URL is now saved!          ║${NC}"
    echo -e "${GREEN}╚════════════════════════════════════════════════════════╝${NC}"
    echo ""
    echo "Before: $current_url"
    echo "After:  $saved_url"
    echo ""
    echo -e "${YELLOW}Next steps:${NC}"
    echo "1. Deploy updated backend to EC2"
    echo "2. Test from UI: https://tacafe.store/#/print-management"
    echo "3. Update Print Bridge URL and click 'Lưu cài đặt'"
    echo "4. Verify it's saved by refreshing the page"
else
    echo -e "${RED}╔════════════════════════════════════════════════════════╗${NC}"
    echo -e "${RED}║   ❌ FAILED! Print Bridge URL was not saved!          ║${NC}"
    echo -e "${RED}╚════════════════════════════════════════════════════════╝${NC}"
    echo ""
    echo "Expected: $TEST_URL"
    echo "Got:      $saved_url"
    echo ""
    echo "The fix may not be deployed yet."
    exit 1
fi
