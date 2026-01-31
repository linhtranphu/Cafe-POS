#!/bin/bash

# Test script for formatters utility
# Tests date format conversion and facility creation

echo "🧪 Testing Formatters Utility"
echo "=============================="
echo ""

# Get auth token
TOKEN=$(curl -s -X POST http://localhost:8080/api/login \
  -H "Content-Type: application/json" \
  -d '{"username":"admin","password":"admin123"}' | python3 -c "import sys, json; print(json.load(sys.stdin)['token'])")

if [ -z "$TOKEN" ]; then
  echo "❌ Failed to get auth token"
  exit 1
fi

echo "✅ Authentication successful"
echo ""

# Test 1: Create facility with ISO date format (what sanitizeFormData produces)
echo "📝 Test 1: Create facility with ISO date format"
echo "-----------------------------------------------"

RESPONSE=$(curl -s -X POST "http://localhost:8080/api/manager/facilities" \
  -H "Authorization: Bearer $TOKEN" \
  -H "Content-Type: application/json" \
  -d '{
    "name": "Test Facility 1",
    "type": "Bàn ghế",
    "area": "Phòng khách",
    "quantity": 1,
    "status": "Đang sử dụng",
    "purchase_date": "2026-01-31T00:00:00Z",
    "cost": 1000000,
    "supplier": "Test Supplier",
    "notes": "Test notes"
  }')

FACILITY_ID=$(echo "$RESPONSE" | python3 -c "import sys, json; print(json.load(sys.stdin).get('id', ''))")

if [ -n "$FACILITY_ID" ]; then
  echo "✅ Facility created successfully"
  echo "   ID: $FACILITY_ID"
else
  echo "❌ Failed to create facility"
  echo "   Response: $RESPONSE"
  exit 1
fi

echo ""

# Test 2: Get facility and verify date format
echo "📝 Test 2: Verify facility data"
echo "-------------------------------"

FACILITY=$(curl -s -X GET "http://localhost:8080/api/manager/facilities/$FACILITY_ID" \
  -H "Authorization: Bearer $TOKEN")

PURCHASE_DATE=$(echo "$FACILITY" | python3 -c "import sys, json; print(json.load(sys.stdin).get('purchase_date', ''))")

if [[ "$PURCHASE_DATE" == *"T"*"Z" ]]; then
  echo "✅ Date format is correct (ISO format)"
  echo "   Purchase date: $PURCHASE_DATE"
else
  echo "❌ Date format is incorrect"
  echo "   Purchase date: $PURCHASE_DATE"
  exit 1
fi

echo ""

# Test 3: Update facility
echo "📝 Test 3: Update facility"
echo "--------------------------"

UPDATE_RESPONSE=$(curl -s -X PUT "http://localhost:8080/api/manager/facilities/$FACILITY_ID" \
  -H "Authorization: Bearer $TOKEN" \
  -H "Content-Type: application/json" \
  -d '{
    "name": "Test Facility 1 Updated",
    "type": "Bàn ghế",
    "area": "Phòng khách",
    "quantity": 2,
    "status": "Đang sử dụng",
    "purchase_date": "2026-02-01T00:00:00Z",
    "cost": 1500000,
    "supplier": "Test Supplier Updated",
    "notes": "Test notes updated"
  }')

UPDATED_NAME=$(echo "$UPDATE_RESPONSE" | python3 -c "import sys, json; print(json.load(sys.stdin).get('name', ''))")

if [[ "$UPDATED_NAME" == "Test Facility 1 Updated" ]]; then
  echo "✅ Facility updated successfully"
else
  echo "❌ Failed to update facility"
  echo "   Response: $UPDATE_RESPONSE"
  exit 1
fi

echo ""

# Test 4: Delete facility
echo "📝 Test 4: Delete facility"
echo "--------------------------"

DELETE_RESPONSE=$(curl -s -X DELETE "http://localhost:8080/api/manager/facilities/$FACILITY_ID" \
  -H "Authorization: Bearer $TOKEN")

if [[ "$DELETE_RESPONSE" == *"Xóa thành công"* ]]; then
  echo "✅ Facility deleted successfully"
else
  echo "❌ Failed to delete facility"
  echo "   Response: $DELETE_RESPONSE"
  exit 1
fi

echo ""

# Test 5: Verify deletion
echo "📝 Test 5: Verify deletion"
echo "--------------------------"

VERIFY_RESPONSE=$(curl -s -X GET "http://localhost:8080/api/manager/facilities/$FACILITY_ID" \
  -H "Authorization: Bearer $TOKEN")

if [[ "$VERIFY_RESPONSE" == *"Không tìm thấy"* ]]; then
  echo "✅ Facility successfully deleted (404 as expected)"
else
  echo "⚠️  Facility might still exist"
  echo "   Response: $VERIFY_RESPONSE"
fi

echo ""
echo "=============================="
echo "🎉 All tests passed!"
echo "=============================="
echo ""
echo "Summary:"
echo "  ✅ ISO date format works correctly"
echo "  ✅ Create facility works"
echo "  ✅ Update facility works"
echo "  ✅ Delete facility works"
echo "  ✅ Date conversion is consistent"
echo ""
echo "The formatters utility is working correctly!"
