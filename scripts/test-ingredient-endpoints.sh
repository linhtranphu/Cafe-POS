#!/bin/bash

# Test Ingredient Endpoints
# Make sure backend is running on localhost:8080

echo "🧪 Testing Ingredient Endpoints"
echo "================================"

# Get token (replace with actual login)
TOKEN="your-jwt-token-here"

echo ""
echo "1️⃣ Testing GET /api/manager/ingredients"
curl -X GET http://localhost:8080/api/manager/ingredients \
  -H "Authorization: Bearer $TOKEN" \
  -w "\nStatus: %{http_code}\n"

echo ""
echo "2️⃣ Testing GET /api/manager/ingredients/low-stock"
curl -X GET http://localhost:8080/api/manager/ingredients/low-stock \
  -H "Authorization: Bearer $TOKEN" \
  -w "\nStatus: %{http_code}\n"

echo ""
echo "3️⃣ Testing GET /api/manager/ingredients/:id/history"
curl -X GET http://localhost:8080/api/manager/ingredients/INGREDIENT_ID/history \
  -H "Authorization: Bearer $TOKEN" \
  -w "\nStatus: %{http_code}\n"

echo ""
echo "4️⃣ Testing POST /api/manager/ingredients/:id/adjust"
curl -X POST http://localhost:8080/api/manager/ingredients/INGREDIENT_ID/adjust \
  -H "Authorization: Bearer $TOKEN" \
  -H "Content-Type: application/json" \
  -d '{"quantity": 10, "reason": "Nhập hàng"}' \
  -w "\nStatus: %{http_code}\n"

echo ""
echo "✅ Test completed"
