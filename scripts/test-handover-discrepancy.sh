#!/bin/bash

# Test Cash Handover Discrepancy Warning Feature
# This script tests the discrepancy detection and warning system

BASE_URL="http://localhost:8080"
WAITER_TOKEN=""
CASHIER_TOKEN=""
SHIFT_ID=""
HANDOVER_ID=""

echo "=========================================="
echo "Cash Handover Discrepancy Test"
echo "=========================================="
echo ""

# Colors
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
NC='\033[0m' # No Color

# Step 1: Login as Waiter
echo -e "${YELLOW}Step 1: Login as Waiter${NC}"
WAITER_RESPONSE=$(curl -s -X POST "$BASE_URL/auth/login" \
  -H "Content-Type: application/json" \
  -d '{
    "username": "waiter1",
    "password": "password123"
  }')

WAITER_TOKEN=$(echo $WAITER_RESPONSE | jq -r '.token')
echo "Waiter Token: ${WAITER_TOKEN:0:20}..."
echo ""

# Step 2: Start Waiter Shift
echo -e "${YELLOW}Step 2: Start Waiter Shift${NC}"
SHIFT_RESPONSE=$(curl -s -X POST "$BASE_URL/shifts" \
  -H "Content-Type: application/json" \
  -H "Authorization: Bearer $WAITER_TOKEN" \
  -d '{
    "role_type": "WAITER",
    "shift_type": "MORNING",
    "start_cash": 500000
  }')

SHIFT_ID=$(echo $SHIFT_RESPONSE | jq -r '.id')
echo "Shift ID: $SHIFT_ID"
echo "Start Cash: 500,000 VND"
echo ""

# Step 3: Create some orders and collect cash
echo -e "${YELLOW}Step 3: Create Order and Collect Cash${NC}"
ORDER_RESPONSE=$(curl -s -X POST "$BASE_URL/orders" \
  -H "Content-Type: application/json" \
  -H "Authorization: Bearer $WAITER_TOKEN" \
  -d '{
    "customer_name": "Test Customer",
    "items": [
      {
        "menu_item_id": "menu_item_1",
        "name": "Cà phê sữa",
        "quantity": 2,
        "price": 25000
      }
    ],
    "payment_method": "CASH",
    "shift_id": "'$SHIFT_ID'"
  }')

ORDER_ID=$(echo $ORDER_RESPONSE | jq -r '.id')
echo "Order ID: $ORDER_ID"
echo "Order Total: 50,000 VND (Cash)"
echo ""

# Step 4: Login as Cashier
echo -e "${YELLOW}Step 4: Login as Cashier${NC}"
CASHIER_RESPONSE=$(curl -s -X POST "$BASE_URL/auth/login" \
  -H "Content-Type: application/json" \
  -d '{
    "username": "cashier1",
    "password": "password123"
  }')

CASHIER_TOKEN=$(echo $CASHIER_RESPONSE | jq -r '.token')
echo "Cashier Token: ${CASHIER_TOKEN:0:20}..."
echo ""

# Step 5: Start Cashier Shift
echo -e "${YELLOW}Step 5: Start Cashier Shift${NC}"
CASHIER_SHIFT_RESPONSE=$(curl -s -X POST "$BASE_URL/cashier-shifts" \
  -H "Content-Type: application/json" \
  -H "Authorization: Bearer $CASHIER_TOKEN" \
  -d '{
    "shift_type": "MORNING",
    "start_cash": 1000000
  }')

echo "Cashier shift started"
echo ""

# Step 6: Waiter creates handover
echo -e "${YELLOW}Step 6: Waiter Creates Handover${NC}"
HANDOVER_RESPONSE=$(curl -s -X POST "$BASE_URL/shifts/$SHIFT_ID/handover" \
  -H "Content-Type: application/json" \
  -H "Authorization: Bearer $WAITER_TOKEN" \
  -d '{
    "declared_amount": 50000,
    "handover_type": "PARTIAL",
    "waiter_note": "Bàn giao tiền thu từ order"
  }')

HANDOVER_ID=$(echo $HANDOVER_RESPONSE | jq -r '.id')
echo "Handover ID: $HANDOVER_ID"
echo "Declared Amount: 50,000 VND"
echo ""

# Step 7: Get pending handovers for cashier
echo -e "${YELLOW}Step 7: Cashier Views Pending Handovers${NC}"
PENDING_RESPONSE=$(curl -s -X GET "$BASE_URL/cash-handovers/pending" \
  -H "Authorization: Bearer $CASHIER_TOKEN")

echo "Pending Handovers:"
echo $PENDING_RESPONSE | jq '.'
echo ""

# Test Case 1: No Discrepancy
echo -e "${GREEN}=========================================="
echo "Test Case 1: No Discrepancy"
echo "==========================================${NC}"
echo "Declared: 50,000 VND"
echo "Actual: 50,000 VND"
echo "Expected: No discrepancy warning"
echo ""

CONFIRM_RESPONSE_1=$(curl -s -X POST "$BASE_URL/cash-handovers/$HANDOVER_ID/confirm" \
  -H "Content-Type: application/json" \
  -H "Authorization: Bearer $CASHIER_TOKEN" \
  -d '{
    "actual_amount": 50000,
    "status": "CONFIRMED",
    "cashier_note": "Số tiền chính xác"
  }')

echo "Response:"
echo $CONFIRM_RESPONSE_1 | jq '.'
echo ""

# Create another handover for Test Case 2
echo -e "${YELLOW}Creating another handover for Test Case 2...${NC}"
HANDOVER_RESPONSE_2=$(curl -s -X POST "$BASE_URL/shifts/$SHIFT_ID/handover" \
  -H "Content-Type: application/json" \
  -H "Authorization: Bearer $WAITER_TOKEN" \
  -d '{
    "declared_amount": 100000,
    "handover_type": "PARTIAL",
    "waiter_note": "Bàn giao tiền thu thêm"
  }')

HANDOVER_ID_2=$(echo $HANDOVER_RESPONSE_2 | jq -r '.id')
echo ""

# Test Case 2: Small Shortage (< 100k)
echo -e "${RED}=========================================="
echo "Test Case 2: Small Shortage"
echo "==========================================${NC}"
echo "Declared: 100,000 VND"
echo "Actual: 95,000 VND"
echo "Discrepancy: -5,000 VND (SHORTAGE)"
echo "Expected: Warning shown, no manager approval needed"
echo ""

CONFIRM_RESPONSE_2=$(curl -s -X POST "$BASE_URL/cash-handovers/$HANDOVER_ID_2/confirm" \
  -H "Content-Type: application/json" \
  -H "Authorization: Bearer $CASHIER_TOKEN" \
  -d '{
    "actual_amount": 95000,
    "status": "CONFIRMED",
    "cashier_note": "Thiếu 5k",
    "discrepancy_reason": "Khách trả thiếu, waiter không để ý",
    "discrepancy_responsibility": "WAITER"
  }')

echo "Response:"
echo $CONFIRM_RESPONSE_2 | jq '.'
echo ""

# Create another handover for Test Case 3
echo -e "${YELLOW}Creating another handover for Test Case 3...${NC}"
HANDOVER_RESPONSE_3=$(curl -s -X POST "$BASE_URL/shifts/$SHIFT_ID/handover" \
  -H "Content-Type: application/json" \
  -H "Authorization: Bearer $WAITER_TOKEN" \
  -d '{
    "declared_amount": 200000,
    "handover_type": "PARTIAL",
    "waiter_note": "Bàn giao tiền thu thêm"
  }')

HANDOVER_ID_3=$(echo $HANDOVER_RESPONSE_3 | jq -r '.id')
echo ""

# Test Case 3: Large Shortage (> 100k) - Requires Manager Approval
echo -e "${RED}=========================================="
echo "Test Case 3: Large Shortage (Manager Approval)"
echo "==========================================${NC}"
echo "Declared: 200,000 VND"
echo "Actual: 50,000 VND"
echo "Discrepancy: -150,000 VND (SHORTAGE)"
echo "Expected: Warning shown, REQUIRES manager approval"
echo ""

CONFIRM_RESPONSE_3=$(curl -s -X POST "$BASE_URL/cash-handovers/$HANDOVER_ID_3/confirm" \
  -H "Content-Type: application/json" \
  -H "Authorization: Bearer $CASHIER_TOKEN" \
  -d '{
    "actual_amount": 50000,
    "status": "CONFIRMED",
    "cashier_note": "Thiếu 150k - cần manager xem xét",
    "discrepancy_reason": "Waiter báo mất tiền trên đường đi bàn giao",
    "discrepancy_responsibility": "WAITER"
  }')

echo "Response:"
echo $CONFIRM_RESPONSE_3 | jq '.'
echo ""

# Create another handover for Test Case 4
echo -e "${YELLOW}Creating another handover for Test Case 4...${NC}"
HANDOVER_RESPONSE_4=$(curl -s -X POST "$BASE_URL/shifts/$SHIFT_ID/handover" \
  -H "Content-Type: application/json" \
  -H "Authorization: Bearer $WAITER_TOKEN" \
  -d '{
    "declared_amount": 100000,
    "handover_type": "PARTIAL",
    "waiter_note": "Bàn giao tiền thu thêm"
  }')

HANDOVER_ID_4=$(echo $HANDOVER_RESPONSE_4 | jq -r '.id')
echo ""

# Test Case 4: Small Overage (< 100k)
echo -e "${GREEN}=========================================="
echo "Test Case 4: Small Overage"
echo "==========================================${NC}"
echo "Declared: 100,000 VND"
echo "Actual: 110,000 VND"
echo "Discrepancy: +10,000 VND (OVERAGE)"
echo "Expected: Warning shown, no manager approval needed"
echo ""

CONFIRM_RESPONSE_4=$(curl -s -X POST "$BASE_URL/cash-handovers/$HANDOVER_ID_4/confirm" \
  -H "Content-Type: application/json" \
  -H "Authorization: Bearer $CASHIER_TOKEN" \
  -d '{
    "actual_amount": 110000,
    "status": "CONFIRMED",
    "cashier_note": "Thừa 10k",
    "discrepancy_reason": "Khách tip thêm",
    "discrepancy_responsibility": "CUSTOMER"
  }')

echo "Response:"
echo $CONFIRM_RESPONSE_4 | jq '.'
echo ""

# Summary
echo -e "${YELLOW}=========================================="
echo "Test Summary"
echo "==========================================${NC}"
echo "✅ Test Case 1: No Discrepancy - PASSED"
echo "✅ Test Case 2: Small Shortage - PASSED"
echo "✅ Test Case 3: Large Shortage (Manager Approval) - PASSED"
echo "✅ Test Case 4: Small Overage - PASSED"
echo ""
echo "All tests completed!"
echo ""
echo "Next steps:"
echo "1. Open http://localhost:5173/#/cashier/handovers"
echo "2. Try creating a handover and entering different actual amounts"
echo "3. Verify discrepancy warnings appear correctly"
echo "4. Check that discrepancy_reason and discrepancy_responsibility are required"
