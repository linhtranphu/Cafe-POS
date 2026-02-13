#!/bin/bash

# Complete setup script for cost analysis testing

echo "━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━"
echo "🚀 SETUP & CALCULATE COSTS - COMPLETE WORKFLOW"
echo "━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━"
echo ""

# Colors
GREEN='\033[0;32m'
BLUE='\033[0;34m'
YELLOW='\033[1;33m'
RED='\033[0;31m'
NC='\033[0m' # No Color

# Backend URL (port 3000 for direct run)
BACKEND_URL="${BACKEND_URL:-http://localhost:3000}"

# Step 1: Check if backend is running
echo "━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━"
echo "📡 Step 1: Checking Backend Status"
echo "━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━"
echo ""

# Try to connect to backend (any response means it's running)
BACKEND_CHECK=$(curl -s -o /dev/null -w "%{http_code}" ${BACKEND_URL}/api/auth/login 2>/dev/null)

if [ "$BACKEND_CHECK" = "200" ] || [ "$BACKEND_CHECK" = "400" ] || [ "$BACKEND_CHECK" = "401" ] || [ "$BACKEND_CHECK" = "405" ]; then
    echo -e "${GREEN}✅ Backend is running at ${BACKEND_URL}${NC}"
else
    echo -e "${RED}❌ Backend is NOT running (HTTP code: $BACKEND_CHECK)${NC}"
    echo ""
    echo "Please start the backend first:"
    echo ""
    echo -e "${YELLOW}Start backend (direct run):${NC}"
    echo "  cd backend"
    echo "  go run main.go"
    echo ""
    echo "Backend should run on port 3000"
    echo ""
    echo "Then run this script again."
    exit 1
fi

echo ""

# Step 2: Check ingredients
echo "━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━"
echo "🥬 Step 2: Getting Auth Token & Checking Ingredients"
echo "━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━"
echo ""

# Get auth token
if [ ! -f ".auth_token" ]; then
    echo "🔐 Getting authentication token..."
    ./get-auth-token.sh
    echo ""
fi

AUTH_TOKEN=$(cat .auth_token 2>/dev/null)

if [ -z "$AUTH_TOKEN" ]; then
    echo -e "${RED}❌ Failed to get auth token${NC}"
    exit 1
fi

INGREDIENTS_RESPONSE=$(curl -s ${BACKEND_URL}/api/manager/ingredients \
  -H "Authorization: Bearer ${AUTH_TOKEN}")
INGREDIENTS_COUNT=$(echo "$INGREDIENTS_RESPONSE" | grep -o '"id":"[^"]*"' | wc -l | tr -d ' ')

if [ "$INGREDIENTS_COUNT" -gt "0" ]; then
    echo -e "${GREEN}✅ Found $INGREDIENTS_COUNT ingredients${NC}"
else
    echo -e "${YELLOW}⚠️  No ingredients found. Seeding ingredients...${NC}"
    echo ""
    cd backend
    go run cmd/seed/main.go
    cd ..
    echo ""
    echo -e "${GREEN}✅ Ingredients seeded${NC}"
fi

echo ""

# Step 3: Seed menu items
echo "━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━"
echo "🍽️  Step 3: Seeding Menu Items"
echo "━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━"
echo ""

./seed-menu-variants-auto.sh

echo ""

# Step 4: Wait a bit for data to be ready
echo "⏳ Waiting for data to be ready..."
sleep 2
echo ""

# Step 5: Calculate costs
echo "━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━"
echo "💰 Step 4: Calculating Costs"
echo "━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━"
echo ""

./calculate-costs-simple.sh

echo ""

# Step 6: Summary
echo "━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━"
echo "🎉 SETUP COMPLETE!"
echo "━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━"
echo ""
echo -e "${GREEN}✅ Backend is running${NC}"
echo -e "${GREEN}✅ Ingredients seeded${NC}"
echo -e "${GREEN}✅ Menu items seeded (10 items)${NC}"
echo -e "${GREEN}✅ Costs calculated${NC}"
echo ""
echo "📊 View Results:"
echo ""
echo "1. Via API (with auth token):"
echo "   export AUTH_TOKEN=\$(cat .auth_token)"
echo "   curl -H \"Authorization: Bearer \$AUTH_TOKEN\" ${BACKEND_URL}/api/manager/menu"
echo "   curl -H \"Authorization: Bearer \$AUTH_TOKEN\" ${BACKEND_URL}/api/manager/menu/{ID}/cost-breakdown"
echo "   curl -H \"Authorization: Bearer \$AUTH_TOKEN\" ${BACKEND_URL}/api/manager/menu/{ID}/profit-analysis"
echo ""
echo "2. Via Frontend:"
echo "   http://localhost:5173/cost-analysis"
echo ""
echo "3. Get sample menu ID:"
echo "   curl -s -H \"Authorization: Bearer \$AUTH_TOKEN\" ${BACKEND_URL}/api/manager/menu | grep -o '\"id\":\"[^\"]*\"' | head -1"
echo ""
echo "━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━"
