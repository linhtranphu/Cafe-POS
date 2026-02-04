#!/bin/bash

# Test script for Cash Handover API endpoints
# This script tests the basic functionality of the cash handover API

BASE_URL="http://localhost:3000/api"
TOKEN=""

echo "🧪 Testing Cash Handover API Endpoints"
echo "======================================="

# Function to make authenticated requests
make_request() {
    local method=$1
    local endpoint=$2
    local data=$3
    
    if [ -n "$data" ]; then
        curl -s -X $method \
            -H "Content-Type: application/json" \
            -H "Authorization: Bearer $TOKEN" \
            -d "$data" \
            "$BASE_URL$endpoint"
    else
        curl -s -X $method \
            -H "Authorization: Bearer $TOKEN" \
            "$BASE_URL$endpoint"
    fi
}

# Step 1: Login to get token
echo "1. 🔐 Logging in as admin..."
LOGIN_RESPONSE=$(curl -s -X POST \
    -H "Content-Type: application/json" \
    -d '{"username":"admin","password":"admin123"}' \
    "$BASE_URL/login")

TOKEN=$(echo $LOGIN_RESPONSE | grep -o '"token":"[^"]*"' | cut -d'"' -f4)

if [ -z "$TOKEN" ]; then
    echo "❌ Login failed. Please ensure the server is running and admin user exists."
    echo "Response: $LOGIN_RESPONSE"
    exit 1
fi

echo "✅ Login successful"

# Step 2: Test getting pending handovers (should be empty initially)
echo ""
echo "2. 📋 Getting pending handovers..."
PENDING_RESPONSE=$(make_request GET "/cash-handovers/pending")
echo "Response: $PENDING_RESPONSE"

# Step 3: Test getting today's handovers (should be empty initially)
echo ""
echo "3. 📅 Getting today's handovers..."
TODAY_RESPONSE=$(make_request GET "/cash-handovers/today")
echo "Response: $TODAY_RESPONSE"

# Step 4: Test getting discrepancy stats
echo ""
echo "4. 📊 Getting discrepancy statistics..."
STATS_RESPONSE=$(make_request GET "/cash-handovers/discrepancy-stats")
echo "Response: $STATS_RESPONSE"

# Step 5: Test manager endpoints
echo ""
echo "5. 👨‍💼 Testing manager endpoints..."
echo "5a. Getting pending approvals..."
APPROVALS_RESPONSE=$(make_request GET "/manager/cash-handovers/pending-approval")
echo "Response: $APPROVALS_RESPONSE"

echo ""
echo "5b. Getting manager discrepancy stats..."
MANAGER_STATS_RESPONSE=$(make_request GET "/manager/discrepancies/stats")
echo "Response: $MANAGER_STATS_RESPONSE"

echo ""
echo "🎉 API endpoint tests completed!"
echo "Note: To test handover creation, you need an active waiter shift."
echo "Note: To test handover confirmation, you need an active cashier shift."