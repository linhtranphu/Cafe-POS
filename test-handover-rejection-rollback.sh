#!/bin/bash

# Test script to verify handover rejection rollback with MongoDB transaction
# This tests that when cashier rejects handover, remaining amounts are restored

echo "🧪 Testing Handover Rejection Rollback"
echo "======================================"
echo ""

# Colors for output
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
NC='\033[0m' # No Color

# Get waiter auth token
echo "1️⃣ Logging in as waiter..."
WAITER_LOGIN=$(curl -s -X POST http://localhost:3000/api/auth/login \
  -H "Content-Type: application/json" \
  -d '{
    "username": "waiter1",
    "password": "password123"
  }')

WAITER_TOKEN=$(echo $WAITER_LOGIN | grep -o '"token":"[^"]*' | cut -d'"' -f4)
WAITER_ID=$(echo $WAITER_LOGIN | grep -o '"id":"[^"]*' | cut -d'"' -f4)

if [ -z "$WAITER_TOKEN" ]; then
  echo -e "${RED}❌ Waiter login failed${NC}"
  exit 1
fi

echo -e "${GREEN}✅ Waiter logged in${NC}"
echo ""

# Get cashier auth token
echo "2️⃣ Logging in as cashier..."
CASHIER_LOGIN=$(curl -s -X POST http://localhost:3000/api/auth/login \
  -H "Content-Type: application/json" \
  -d '{
    "username": "cashier1",
    "password": "password123"
  }')

CASHIER_TOKEN=$(echo $CASHIER_LOGIN | grep -o '"token":"[^"]*' | cut -d'"' -f4)
CASHIER_ID=$(echo $CASHIER_LOGIN | grep -o '"id":"[^"]*' | cut -d'"' -f4)

if [ -z "$CASHIER_TOKEN" ]; then
  echo -e "${RED}❌ Cashier login failed${NC}"
  exit 1
fi

echo -e "${GREEN}✅ Cashier logged in${NC}"
echo ""

# Start waiter shift
echo "3️⃣ Starting waiter shift with 100k start cash..."
START_SHIFT=$(curl -s -X POST http://localhost:3000/api/shifts \
  -H "Content-Type: application/json" \
  -H "Authorization: Bearer $WAITER_TOKEN" \
  -d '{
    "type": "morning",
    "start_cash": 100000
  }')

SHIFT_ID=$(echo $START_SHIFT | grep -o '"id":"[^"]*' | cut -d'"' -f4)

if [ -z "$SHIFT_ID" ]; then
  echo -e "${RED}❌ Failed to start shift${NC}"
  echo "Response: $START_SHIFT"
  exit 1
fi

echo -e "${GREEN}✅ Shift started: $SHIFT_ID${NC}"
echo "   Start Cash: 100,000 VND"
echo ""

# Create and collect payment to add more cash
echo "4️⃣ Creating order and collecting payment (50k)..."
CREATE_ORDER=$(curl -s -X POST http://localhost:3000/api/orders \
  -H "Content-Type: application/json" \
  -H "Authorization: Bearer $WAITER_TOKEN" \
  -d '{
    "items": [
      {
        "product_id": "product1",
        "product_name": "Cà phê sữa",
        "quantity": 1,
        "price": 50000
      }
    ],
    "table_number": "A1"
  }')

ORDER_ID=$(echo $CREATE_ORDER | grep -o '"id":"[^"]*' | cut -d'"' -f4)

PAYMENT=$(curl -s -X POST http://localhost:3000/api/orders/$ORDER_ID/collect-payment \
  -H "Content-Type: application/json" \
  -H "Authorization: Bearer $WAITER_TOKEN" \
  -d '{
    "payment_method": "CASH",
    "amount": 50000
  }')

echo -e "${GREEN}✅ Payment collected: 50,000 VND${NC}"
echo ""

# Get shift to check remaining_cash before handover
echo "5️⃣ Checking shift before handover..."
SHIFT_BEFORE=$(curl -s -X GET http://localhost:3000/api/shifts/current \
  -H "Authorization: Bearer $WAITER_TOKEN")

REMAINING_CASH_BEFORE=$(echo $SHIFT_BEFORE | grep -o '"remaining_cash":[0-9]*' | cut -d':' -f2)
REMAINING_TRANSFER_BEFORE=$(echo $SHIFT_BEFORE | grep -o '"remaining_transfer":[0-9]*' | cut -d':' -f2)

echo "   Remaining Cash: ${REMAINING_CASH_BEFORE} VND"
echo "   Remaining Transfer: ${REMAINING_TRANSFER_BEFORE} VND"
echo ""

# Create handover
echo "6️⃣ Creating handover (30k cash)..."
CREATE_HANDOVER=$(curl -s -X POST http://localhost:3000/api/handovers \
  -H "Content-Type: application/json" \
  -H "Authorization: Bearer $WAITER_TOKEN" \
  -d '{
    "cash_amount": 30000,
    "transfer_amount": 0,
    "handover_type": "PARTIAL",
    "waiter_note": "Test handover rejection"
  }')

HANDOVER_ID=$(echo $CREATE_HANDOVER | grep -o '"id":"[^"]*' | cut -d'"' -f4)

if [ -z "$HANDOVER_ID" ]; then
  echo -e "${RED}❌ Failed to create handover${NC}"
  echo "Response: $CREATE_HANDOVER"
  exit 1
fi

echo -e "${GREEN}✅ Handover created: $HANDOVER_ID${NC}"
echo ""

# Get shift after handover creation
echo "7️⃣ Checking shift after handover creation..."
SHIFT_AFTER_CREATE=$(curl -s -X GET http://localhost:3000/api/shifts/current \
  -H "Authorization: Bearer $WAITER_TOKEN")

REMAINING_CASH_AFTER_CREATE=$(echo $SHIFT_AFTER_CREATE | grep -o '"remaining_cash":[0-9]*' | cut -d':' -f2)
REMAINING_TRANSFER_AFTER_CREATE=$(echo $SHIFT_AFTER_CREATE | grep -o '"remaining_transfer":[0-9]*' | cut -d':' -f2)

echo "   Remaining Cash: ${REMAINING_CASH_AFTER_CREATE} VND"
echo "   Remaining Transfer: ${REMAINING_TRANSFER_AFTER_CREATE} VND"

EXPECTED_AFTER_CREATE=$((REMAINING_CASH_BEFORE - 30000))
if [ "$REMAINING_CASH_AFTER_CREATE" = "$EXPECTED_AFTER_CREATE" ]; then
  echo -e "${GREEN}✅ Cash deducted correctly: ${REMAINING_CASH_BEFORE} - 30000 = ${REMAINING_CASH_AFTER_CREATE}${NC}"
else
  echo -e "${RED}❌ Cash deduction wrong: expected ${EXPECTED_AFTER_CREATE}, got ${REMAINING_CASH_AFTER_CREATE}${NC}"
fi
echo ""

# Cashier REJECTS handover
echo "8️⃣ Cashier rejecting handover..."
REJECT_HANDOVER=$(curl -s -X POST http://localhost:3000/api/handovers/$HANDOVER_ID/confirm \
  -H "Content-Type: application/json" \
  -H "Authorization: Bearer $CASHIER_TOKEN" \
  -d '{
    "cash_actual_amount": 0,
    "transfer_actual_amount": 0,
    "status": "REJECTED",
    "cashier_note": "Testing rejection rollback"
  }')

echo -e "${YELLOW}⚠️  Handover rejected by cashier${NC}"
echo ""

# Get shift after rejection - THIS IS THE KEY TEST
echo "9️⃣ Checking shift after rejection (should restore amounts)..."
sleep 1  # Give DB time to update
SHIFT_AFTER_REJECT=$(curl -s -X GET http://localhost:3000/api/shifts/current \
  -H "Authorization: Bearer $WAITER_TOKEN")

REMAINING_CASH_AFTER_REJECT=$(echo $SHIFT_AFTER_REJECT | grep -o '"remaining_cash":[0-9]*' | cut -d':' -f2)
REMAINING_TRANSFER_AFTER_REJECT=$(echo $SHIFT_AFTER_REJECT | grep -o '"remaining_transfer":[0-9]*' | cut -d':' -f2)

echo "   Remaining Cash: ${REMAINING_CASH_AFTER_REJECT} VND"
echo "   Remaining Transfer: ${REMAINING_TRANSFER_AFTER_REJECT} VND"
echo ""

# Verify rollback
echo "🔍 Verifying rollback..."
echo "======================================"

if [ "$REMAINING_CASH_AFTER_REJECT" = "$REMAINING_CASH_BEFORE" ]; then
  echo -e "${GREEN}✅ ROLLBACK SUCCESS: Cash restored to ${REMAINING_CASH_BEFORE} VND${NC}"
  ROLLBACK_SUCCESS=true
else
  echo -e "${RED}❌ ROLLBACK FAILED: Expected ${REMAINING_CASH_BEFORE}, got ${REMAINING_CASH_AFTER_REJECT}${NC}"
  ROLLBACK_SUCCESS=false
fi

if [ "$REMAINING_TRANSFER_AFTER_REJECT" = "$REMAINING_TRANSFER_BEFORE" ]; then
  echo -e "${GREEN}✅ Transfer unchanged: ${REMAINING_TRANSFER_BEFORE} VND${NC}"
else
  echo -e "${RED}❌ Transfer changed unexpectedly${NC}"
  ROLLBACK_SUCCESS=false
fi

echo ""
echo "======================================"
echo "Summary:"
echo "  Before handover: ${REMAINING_CASH_BEFORE} VND"
echo "  After create:    ${REMAINING_CASH_AFTER_CREATE} VND (deducted 30k)"
echo "  After reject:    ${REMAINING_CASH_AFTER_REJECT} VND (should be ${REMAINING_CASH_BEFORE})"
echo "======================================"

if [ "$ROLLBACK_SUCCESS" = true ]; then
  echo -e "${GREEN}✅ TEST PASSED: Transaction rollback works correctly!${NC}"
  exit 0
else
  echo -e "${RED}❌ TEST FAILED: Transaction rollback not working${NC}"
  exit 1
fi
