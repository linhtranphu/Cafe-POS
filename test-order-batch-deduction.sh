#!/bin/bash

echo "🧪 Testing Order Batch Deduction"
echo "━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━"
echo ""

# Colors
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
RED='\033[0;31m'
NC='\033[0m'

# Get auth token
echo "📝 Step 1: Getting auth token..."
TOKEN_RESPONSE=$(curl -s -X POST http://localhost:3000/api/login \
  -H "Content-Type: application/json" \
  -d '{
    "username": "admin",
    "password": "admin123"
  }')

TOKEN=$(echo $TOKEN_RESPONSE | grep -o '"token":"[^"]*"' | cut -d'"' -f4)

if [ -z "$TOKEN" ]; then
    echo -e "${RED}❌ Failed to get auth token${NC}"
    echo "Response: $TOKEN_RESPONSE"
    exit 1
fi

echo -e "${GREEN}✅ Got auth token${NC}"
echo ""

# Create an ingredient first
echo "📝 Step 2: Creating ingredient..."
INGREDIENT_RESPONSE=$(curl -s -X POST http://localhost:3000/api/manager/ingredients \
  -H "Content-Type: application/json" \
  -H "Authorization: Bearer $TOKEN" \
  -d '{
    "name": "Coffee Beans Test",
    "category": "Raw Materials",
    "unit": "g",
    "quantity": 1000,
    "cost_per_unit": 0.5
  }')

INGREDIENT_ID=$(echo $INGREDIENT_RESPONSE | grep -o '"_id":"[^"]*"' | head -1 | cut -d'"' -f4)
if [ -z "$INGREDIENT_ID" ]; then
    INGREDIENT_ID=$(echo $INGREDIENT_RESPONSE | grep -o '"id":"[^"]*"' | head -1 | cut -d'"' -f4)
fi

if [ -z "$INGREDIENT_ID" ]; then
    echo -e "${RED}❌ Failed to create ingredient${NC}"
    echo "Response: $INGREDIENT_RESPONSE"
    exit 1
fi

echo -e "${GREEN}✅ Created ingredient: $INGREDIENT_ID${NC}"
echo ""

# Create a batch definition
echo "📝 Step 3: Creating batch definition..."
BATCH_DEF_RESPONSE=$(curl -s -X POST http://localhost:3000/api/manager/batch-definitions \
  -H "Content-Type: application/json" \
  -H "Authorization: Bearer $TOKEN" \
  -d "{
    \"name\": \"Test Coffee Batch\",
    \"unit\": \"ml\",
    \"shelf_life_hours\": 168,
    \"low_stock_threshold\": 100,
    \"expiry_warning_hours\": 24,
    \"conversion_rates\": [
      {
        \"source_ingredient_id\": \"$INGREDIENT_ID\",
        \"source_ingredient_name\": \"Coffee Beans Test\",
        \"source_quantity\": 100,
        \"source_unit\": \"g\",
        \"batch_quantity\": 500,
        \"wastage_rate\": 0.1
      }
    ]
  }")

BATCH_DEF_ID=$(echo $BATCH_DEF_RESPONSE | grep -o '"_id":"[^"]*"' | head -1 | cut -d'"' -f4)
if [ -z "$BATCH_DEF_ID" ]; then
    BATCH_DEF_ID=$(echo $BATCH_DEF_RESPONSE | grep -o '"id":"[^"]*"' | head -1 | cut -d'"' -f4)
fi

if [ -z "$BATCH_DEF_ID" ]; then
    echo -e "${RED}❌ Failed to create batch definition${NC}"
    echo "Response: $BATCH_DEF_RESPONSE"
    exit 1
fi

echo -e "${GREEN}✅ Created batch definition: $BATCH_DEF_ID${NC}"
echo ""

# Create a batch record
echo "📝 Step 4: Creating batch record..."
BATCH_RECORD_RESPONSE=$(curl -s -X POST http://localhost:3000/api/batch-records \
  -H "Content-Type: application/json" \
  -H "Authorization: Bearer $TOKEN" \
  -d "{
    \"batch_definition_id\": \"$BATCH_DEF_ID\",
    \"quantity_produced\": 500,
    \"prepared_by\": \"Test User\",
    \"production_date\": \"$(date -u +%Y-%m-%dT%H:%M:%SZ)\",
    \"expiry_date\": \"$(date -u -v+7d +%Y-%m-%dT%H:%M:%SZ 2>/dev/null || date -u -d '+7 days' +%Y-%m-%dT%H:%M:%SZ)\",
    \"notes\": \"Test batch for order deduction\"
  }")

BATCH_RECORD_ID=$(echo $BATCH_RECORD_RESPONSE | grep -o '"_id":"[^"]*"' | head -1 | cut -d'"' -f4)
if [ -z "$BATCH_RECORD_ID" ]; then
    BATCH_RECORD_ID=$(echo $BATCH_RECORD_RESPONSE | grep -o '"id":"[^"]*"' | head -1 | cut -d'"' -f4)
fi

if [ -z "$BATCH_RECORD_ID" ]; then
    echo -e "${RED}❌ Failed to create batch record${NC}"
    echo "Response: $BATCH_RECORD_RESPONSE"
    exit 1
fi

echo -e "${GREEN}✅ Created batch record: $BATCH_RECORD_ID${NC}"
echo ""

# Check initial batch quantity
echo "📝 Step 5: Checking initial batch quantity..."
INITIAL_BATCH=$(curl -s -X GET "http://localhost:3000/api/batch-records/$BATCH_RECORD_ID" \
  -H "Authorization: Bearer $TOKEN")

INITIAL_QTY=$(echo $INITIAL_BATCH | grep -o '"quantity_remaining":[0-9.]*' | cut -d':' -f2)
echo -e "${YELLOW}Initial quantity: $INITIAL_QTY ml${NC}"
echo ""

# Create a menu item with batch ingredient
echo "📝 Step 6: Creating menu item with batch ingredient..."
MENU_ITEM_RESPONSE=$(curl -s -X POST http://localhost:3000/api/manager/menu \
  -H "Content-Type: application/json" \
  -H "Authorization: Bearer $TOKEN" \
  -d "{
    \"name\": \"Test Coffee Drink\",
    \"category\": \"Coffee\",
    \"price\": 50000,
    \"available\": true,
    \"has_variants\": false,
    \"ingredients\": [
      {
        \"name\": \"Coffee Concentrate\",
        \"quantity\": 30,
        \"unit\": \"ml\",
        \"ingredient_type\": \"batch\",
        \"batch_id\": \"$BATCH_DEF_ID\"
      }
    ]
  }")

MENU_ITEM_ID=$(echo $MENU_ITEM_RESPONSE | grep -o '"_id":"[^"]*"' | head -1 | cut -d'"' -f4)
if [ -z "$MENU_ITEM_ID" ]; then
    MENU_ITEM_ID=$(echo $MENU_ITEM_RESPONSE | grep -o '"id":"[^"]*"' | head -1 | cut -d'"' -f4)
fi

if [ -z "$MENU_ITEM_ID" ]; then
    echo -e "${RED}❌ Failed to create menu item${NC}"
    echo "Response: $MENU_ITEM_RESPONSE"
    exit 1
fi

echo -e "${GREEN}✅ Created menu item: $MENU_ITEM_ID${NC}"
echo ""

# Start a shift or get current shift
echo "📝 Step 7: Getting or starting a shift..."
SHIFT_RESPONSE=$(curl -s -X GET http://localhost:3000/api/shifts/current \
  -H "Authorization: Bearer $TOKEN")

SHIFT_ID=$(echo $SHIFT_RESPONSE | grep -o '"_id":"[^"]*"' | head -1 | cut -d'"' -f4)
if [ -z "$SHIFT_ID" ]; then
    SHIFT_ID=$(echo $SHIFT_RESPONSE | grep -o '"id":"[^"]*"' | head -1 | cut -d'"' -f4)
fi

# If no current shift, start a new one
if [ -z "$SHIFT_ID" ]; then
    echo "No current shift, starting new one..."
    SHIFT_RESPONSE=$(curl -s -X POST http://localhost:3000/api/shifts/start \
      -H "Content-Type: application/json" \
      -H "Authorization: Bearer $TOKEN" \
      -d '{
        "type": "waiter",
        "starting_cash": 0
      }')
    
    SHIFT_ID=$(echo $SHIFT_RESPONSE | grep -o '"_id":"[^"]*"' | head -1 | cut -d'"' -f4)
    if [ -z "$SHIFT_ID" ]; then
        SHIFT_ID=$(echo $SHIFT_RESPONSE | grep -o '"id":"[^"]*"' | head -1 | cut -d'"' -f4)
    fi
fi

if [ -z "$SHIFT_ID" ]; then
    echo -e "${RED}❌ Failed to get or start shift${NC}"
    echo "Response: $SHIFT_RESPONSE"
    exit 1
fi

echo -e "${GREEN}✅ Started shift: $SHIFT_ID${NC}"
echo ""

# Create an order
echo "📝 Step 8: Creating order (should deduct batch)..."
ORDER_RESPONSE=$(curl -s -X POST http://localhost:3000/api/waiter/orders \
  -H "Content-Type: application/json" \
  -H "Authorization: Bearer $TOKEN" \
  -d "{
    \"shift_id\": \"$SHIFT_ID\",
    \"customer_name\": \"Test Customer\",
    \"items\": [
      {
        \"menu_item_id\": \"$MENU_ITEM_ID\",
        \"quantity\": 2
      }
    ]
  }")

ORDER_ID=$(echo $ORDER_RESPONSE | grep -o '"_id":"[^"]*"' | head -1 | cut -d'"' -f4)
if [ -z "$ORDER_ID" ]; then
    ORDER_ID=$(echo $ORDER_RESPONSE | grep -o '"id":"[^"]*"' | head -1 | cut -d'"' -f4)
fi

if [ -z "$ORDER_ID" ]; then
    echo -e "${RED}❌ Failed to create order${NC}"
    echo "Response: $ORDER_RESPONSE"
    exit 1
fi

echo -e "${GREEN}✅ Created order: $ORDER_ID${NC}"
echo ""

# Check batch quantity after order
echo "📝 Step 9: Checking batch quantity after order..."
sleep 2
FINAL_BATCH=$(curl -s -X GET "http://localhost:3000/api/batch-records/$BATCH_RECORD_ID" \
  -H "Authorization: Bearer $TOKEN")

FINAL_QTY=$(echo $FINAL_BATCH | grep -o '"quantity_remaining":[0-9.]*' | cut -d':' -f2)
echo -e "${YELLOW}Final quantity: $FINAL_QTY ml${NC}"
echo ""

# Calculate expected quantity
EXPECTED_QTY=$(echo "$INITIAL_QTY - 60" | bc) # 2 drinks * 30ml each = 60ml

echo "━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━"
echo "📊 Results:"
echo "━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━"
echo "Initial quantity:  $INITIAL_QTY ml"
echo "Order quantity:    2 drinks × 30ml = 60ml"
echo "Expected final:    $EXPECTED_QTY ml"
echo "Actual final:      $FINAL_QTY ml"
echo ""

if [ "$FINAL_QTY" = "$EXPECTED_QTY" ]; then
    echo -e "${GREEN}✅ SUCCESS: Batch was deducted correctly!${NC}"
else
    echo -e "${RED}❌ FAILED: Batch was NOT deducted correctly!${NC}"
    echo ""
    echo "Possible issues:"
    echo "1. BatchUsageService not properly integrated"
    echo "2. Menu item ingredient_type not set to 'batch'"
    echo "3. batch_id not properly linked"
    echo ""
    echo "Check backend logs for errors:"
    echo "  tail -f backend.log"
fi

echo ""
echo "━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━"
