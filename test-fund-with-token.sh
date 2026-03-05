#!/bin/bash

# Test fund deposit with your actual token
# Get token from browser: localStorage.getItem('token')

if [ -z "$1" ]; then
  echo "Usage: ./test-fund-with-token.sh YOUR_TOKEN"
  echo ""
  echo "Get your token from browser console:"
  echo "  localStorage.getItem('token')"
  exit 1
fi

TOKEN=$1

echo "Testing Fund Deposit API with provided token..."
echo ""

# Test balance first
echo "1. Testing GET balance..."
curl -s -w "\nHTTP_CODE:%{http_code}\n" http://localhost:3000/api/manager/fund/balance \
  -H "Authorization: Bearer $TOKEN" | tee /tmp/balance_response.txt
echo ""

# Test deposit
echo "2. Testing POST deposit..."
curl -s -w "\nHTTP_CODE:%{http_code}\n" -X POST http://localhost:3000/api/manager/fund/deposit \
  -H "Content-Type: application/json" \
  -H "Authorization: Bearer $TOKEN" \
  -d '{
    "cash_amount": 50000,
    "transfer_amount": 0,
    "reason": "Test deposit from debug script"
  }' | tee /tmp/deposit_response.txt

echo ""
echo "Responses saved to /tmp/balance_response.txt and /tmp/deposit_response.txt"
