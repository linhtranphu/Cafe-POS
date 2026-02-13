#!/bin/bash

# Test creating a multi-size menu item with variants

curl -X POST http://localhost:8080/api/manager/menu \
  -H "Content-Type: application/json" \
  -d '{
    "name": "Cà phê sữa đá",
    "category": "Cà phê",
    "description": "Cà phê phin truyền thống",
    "has_variants": true,
    "variants": [
      {
        "id": "M",
        "name": "Size M",
        "price": 20000,
        "ingredients": [
          {"name": "Cà phê", "quantity": 20, "unit": "g"},
          {"name": "Sữa đặc", "quantity": 30, "unit": "ml"}
        ],
        "available": true,
        "is_default": true
      },
      {
        "id": "L",
        "name": "Size L",
        "price": 25000,
        "ingredients": [
          {"name": "Cà phê", "quantity": 30, "unit": "g"},
          {"name": "Sữa đặc", "quantity": 45, "unit": "ml"}
        ],
        "available": true,
        "is_default": false
      }
    ]
  }' | jq .
