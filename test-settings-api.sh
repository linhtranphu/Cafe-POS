#!/bin/bash

# Login to get token
echo "=== Login ==="
TOKEN=$(curl -s -X POST http://localhost:3000/api/auth/login \
  -H "Content-Type: application/json" \
  -d '{"username":"admin","password":"admin123"}' | jq -r '.token')

echo "Token: ${TOKEN:0:20}..."

# Test /api/settings
echo ""
echo "=== GET /api/settings ==="
curl -s http://localhost:3000/api/settings \
  -H "Authorization: Bearer $TOKEN" | jq '.'

# Test /api/manager/shop-settings
echo ""
echo "=== GET /api/manager/shop-settings ==="
curl -s http://localhost:3000/api/manager/shop-settings \
  -H "Authorization: Bearer $TOKEN" | jq '.'
