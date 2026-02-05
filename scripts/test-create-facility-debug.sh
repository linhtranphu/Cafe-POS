#!/bin/bash

# Test Create Facility with Debug
# This script tests facility creation with various scenarios

echo "🧪 Testing Facility Creation..."
echo ""

# Get auth token
echo "1️⃣ Logging in as manager..."
LOGIN_RESPONSE=$(curl -s -X POST http://localhost:3000/api/login \
  -H "Content-Type: application/json" \
  -d '{
    "username": "admin",
    "password": "admin123"
  }')

TOKEN=$(echo $LOGIN_RESPONSE | jq -r '.token')

if [ "$TOKEN" = "null" ] || [ -z "$TOKEN" ]; then
  echo "❌ Login failed"
  echo "Response: $LOGIN_RESPONSE"
  exit 1
fi

echo "✅ Login successful"
echo "Token: ${TOKEN:0:20}..."
echo ""

# Test 1: Create facility WITHOUT purchase_date
echo "2️⃣ Test 1: Create facility WITHOUT purchase_date"
RESPONSE1=$(curl -s -X POST http://localhost:3000/api/manager/facilities \
  -H "Content-Type: application/json" \
  -H "Authorization: Bearer $TOKEN" \
  -d '{
    "name": "Test Bàn gỗ 1",
    "type": "Bàn ghế",
    "area": "Phòng khách",
    "quantity": 5,
    "status": "Đang sử dụng",
    "cost": 5000000,
    "supplier": "Nhà cung cấp A",
    "notes": "Test without purchase_date"
  }')

echo "Response:"
echo $RESPONSE1 | jq '.'
echo ""

# Test 2: Create facility WITH purchase_date
echo "3️⃣ Test 2: Create facility WITH purchase_date"
RESPONSE2=$(curl -s -X POST http://localhost:3000/api/manager/facilities \
  -H "Content-Type: application/json" \
  -H "Authorization: Bearer $TOKEN" \
  -d '{
    "name": "Test Bàn gỗ 2",
    "type": "Bàn ghế",
    "area": "Phòng khách",
    "quantity": 3,
    "status": "Đang sử dụng",
    "purchase_date": "2026-02-05T00:00:00Z",
    "cost": 3000000,
    "supplier": "Nhà cung cấp B",
    "notes": "Test with purchase_date"
  }')

echo "Response:"
echo $RESPONSE2 | jq '.'
echo ""

# Test 3: Create facility with EMPTY purchase_date (should fail)
echo "4️⃣ Test 3: Create facility with EMPTY purchase_date (should fail)"
RESPONSE3=$(curl -s -X POST http://localhost:3000/api/manager/facilities \
  -H "Content-Type: application/json" \
  -H "Authorization: Bearer $TOKEN" \
  -d '{
    "name": "Test Bàn gỗ 3",
    "type": "Bàn ghế",
    "area": "Phòng khách",
    "quantity": 2,
    "status": "Đang sử dụng",
    "purchase_date": "",
    "cost": 2000000,
    "supplier": "Nhà cung cấp C",
    "notes": "Test with empty purchase_date"
  }')

echo "Response:"
echo $RESPONSE3 | jq '.'
echo ""

# Test 4: Create facility with minimal data
echo "5️⃣ Test 4: Create facility with minimal data"
RESPONSE4=$(curl -s -X POST http://localhost:3000/api/manager/facilities \
  -H "Content-Type: application/json" \
  -H "Authorization: Bearer $TOKEN" \
  -d '{
    "name": "Test Ghế đơn",
    "type": "Bàn ghế",
    "area": "Phòng khách",
    "quantity": 10,
    "status": "Đang sử dụng"
  }')

echo "Response:"
echo $RESPONSE4 | jq '.'
echo ""

echo "✅ All tests completed!"
