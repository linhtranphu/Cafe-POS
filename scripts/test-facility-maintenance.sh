#!/bin/bash

# Test facility maintenance history
FACILITY_ID="697629c9d0ac2facdbb23baa"
TOKEN=$(curl -s -X POST http://localhost:8080/api/login \
  -H "Content-Type: application/json" \
  -d '{"username":"admin","password":"admin123"}' | python3 -c "import sys, json; print(json.load(sys.stdin)['token'])")

echo "=== Testing Facility Maintenance History ==="
echo "Facility ID: $FACILITY_ID"
echo ""

echo "1. Get Facility Details:"
curl -s -X GET "http://localhost:8080/api/manager/facilities/$FACILITY_ID" \
  -H "Authorization: Bearer $TOKEN"
echo ""
echo ""

echo "2. Get Maintenance History:"
curl -s -X GET "http://localhost:8080/api/manager/facilities/$FACILITY_ID/maintenance" \
  -H "Authorization: Bearer $TOKEN"
echo ""
echo ""

echo "3. Try to Delete:"
curl -s -X DELETE "http://localhost:8080/api/manager/facilities/$FACILITY_ID" \
  -H "Authorization: Bearer $TOKEN"
echo ""

