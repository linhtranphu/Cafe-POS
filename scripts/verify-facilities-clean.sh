#!/bin/bash

TOKEN=$(curl -s -X POST http://localhost:8080/api/login \
  -H "Content-Type: application/json" \
  -d '{"username":"admin","password":"admin123"}' | python3 -c "import sys, json; print(json.load(sys.stdin)['token'])")

echo "=== Verifying Facilities Cleanup ==="
echo ""

echo "1. Get All Facilities:"
curl -s -X GET "http://localhost:8080/api/manager/facilities" \
  -H "Authorization: Bearer $TOKEN"
echo ""
echo ""

echo "2. Get Maintenance Schedule:"
curl -s -X GET "http://localhost:8080/api/manager/maintenance/scheduled" \
  -H "Authorization: Bearer $TOKEN"
echo ""
echo ""

echo "3. Get Issue Reports:"
curl -s -X GET "http://localhost:8080/api/manager/issues" \
  -H "Authorization: Bearer $TOKEN"
echo ""
