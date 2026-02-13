#!/bin/bash

# Script to get authentication token for API calls

BACKEND_URL="${BACKEND_URL:-http://localhost:3000}"

# Default admin credentials (adjust if needed)
USERNAME="${API_USERNAME:-admin}"
PASSWORD="${API_PASSWORD:-admin123}"

echo "🔐 Getting authentication token..."
echo ""

# Login and get token
RESPONSE=$(curl -s -X POST "${BACKEND_URL}/api/auth/login" \
  -H "Content-Type: application/json" \
  -d "{\"username\":\"${USERNAME}\",\"password\":\"${PASSWORD}\"}")

# Extract token
TOKEN=$(echo "$RESPONSE" | grep -o '"token":"[^"]*"' | sed 's/"token":"//g' | sed 's/"//g')

if [ -z "$TOKEN" ]; then
    echo "❌ Failed to get token"
    echo "Response: $RESPONSE"
    echo ""
    echo "Please check:"
    echo "  1. Backend is running on ${BACKEND_URL}"
    echo "  2. Username: ${USERNAME}"
    echo "  3. Password: ${PASSWORD}"
    exit 1
fi

echo "✅ Token obtained successfully!"
echo ""
echo "Token: $TOKEN"
echo ""
echo "Export to use in other scripts:"
echo "export AUTH_TOKEN=\"$TOKEN\""
echo ""

# Save to file for other scripts to use
echo "$TOKEN" > .auth_token
echo "💾 Token saved to .auth_token file"
