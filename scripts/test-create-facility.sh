#!/bin/bash

TOKEN=$(curl -s -X POST http://localhost:8080/api/login \
  -H "Content-Type: application/json" \
  -d '{"username":"admin","password":"admin123"}' | python3 -c "import sys, json; print(json.load(sys.stdin)['token'])")

echo "=== Testing Create Facility ==="
echo ""

echo "Creating a test facility..."
curl -v -X POST "http://localhost:8080/api/manager/facilities" \
  -H "Authorization: Bearer $TOKEN" \
  -H "Content-Type: application/json" \
  -d '{
    "name": "Máy pha cà phê",
    "type": "Thiết bị",
    "area": "Quầy bar",
    "quantity": 1,
    "status": "in_use",
    "purchase_date": "2024-01-15T00:00:00Z",
    "cost": 50000000,
    "supplier": "La Marzocco",
    "notes": "Máy pha cà phê chuyên nghiệp"
  }'
echo ""
