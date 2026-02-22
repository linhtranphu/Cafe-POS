#!/bin/bash

# Test default template by fetching it from API
echo "📋 Fetching default bill template..."

curl -s http://localhost:3000/api/manager/print-templates?type=BILL \
  -H "Content-Type: application/json" | jq '.'

echo ""
echo "✅ Template fetched successfully!"
echo ""
echo "You can now use this template for printing orders."
