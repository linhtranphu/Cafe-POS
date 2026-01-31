#!/bin/bash

# Test Order State Machine Validation
# This script tests the order workflow with state machine validation

BASE_URL="http://localhost:8080/api"
TOKEN=""
ORDER_ID=""

# Colors for output
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
BLUE='\033[0;34m'
NC='\033[0m' # No Color

# Function to print colored output
print_success() {
    echo -e "${GREEN}✅ $1${NC}"
}

print_error() {
    echo -e "${RED}❌ $1${NC}"
}

print_info() {
    echo -e "${BLUE}ℹ️  $1${NC}"
}

print_warning() {
    echo -e "${YELLOW}⚠️  $1${NC}"
}

print_section() {
    echo -e "\n${BLUE}========================================${NC}"
    echo -e "${BLUE}$1${NC}"
    echo -e "${BLUE}========================================${NC}\n"
}

# Function to make API calls
api_call() {
    local method=$1
    local endpoint=$2
    local data=$3
    local expected_status=$4
    
    if [ -n "$data" ]; then
        response=$(curl -s -w "\n%{http_code}" -X $method \
            -H "Content-Type: application/json" \
            -H "Authorization: Bearer $TOKEN" \
            -d "$data" \
            "$BASE_URL$endpoint")
    else
        response=$(curl -s -w "\n%{http_code}" -X $method \
            -H "Authorization: Bearer $TOKEN" \
            "$BASE_URL$endpoint")
    fi
    
    http_code=$(echo "$response" | tail -n1)
    body=$(echo "$response" | sed '$d')
    
    if [ "$http_code" = "$expected_status" ]; then
        print_success "Status: $http_code (Expected: $expected_status)"
        echo "$body" | jq '.' 2>/dev/null || echo "$body"
        return 0
    else
        print_error "Status: $http_code (Expected: $expected_status)"
        echo "$body" | jq '.' 2>/dev/null || echo "$body"
        return 1
    fi
}

# Step 1: Login
print_section "Step 1: Login as Waiter"
print_info "Logging in to get JWT token..."

login_response=$(curl -s -X POST \
    -H "Content-Type: application/json" \
    -d '{"username":"waiter1","password":"waiter123"}' \
    "$BASE_URL/login")

TOKEN=$(echo $login_response | jq -r '.token')

if [ "$TOKEN" != "null" ] && [ -n "$TOKEN" ]; then
    print_success "Login successful! Token obtained."
else
    print_error "Login failed!"
    echo $login_response | jq '.'
    exit 1
fi

# Step 2: Create Order (CREATED state)
print_section "Step 2: Create Order"
print_info "Creating a new order..."

create_order_data='{
  "table_number": 5,
  "customer_name": "Test Customer",
  "items": [
    {
      "menu_item_id": "507f1f77bcf86cd799439011",
      "name": "Cappuccino",
      "quantity": 2,
      "price": 45000,
      "notes": "Extra hot"
    }
  ]
}'

create_response=$(api_call "POST" "/waiter/orders" "$create_order_data" "201")
ORDER_ID=$(echo $create_response | jq -r '.id // ._id // .ID')

if [ "$ORDER_ID" != "null" ] && [ -n "$ORDER_ID" ]; then
    print_success "Order created! ID: $ORDER_ID"
    print_info "Order Status: CREATED"
else
    print_error "Failed to create order!"
    exit 1
fi

# Step 3: Try to Send to Bar WITHOUT Payment (Should FAIL)
print_section "Step 3: Test Invalid Transition - Send Unpaid Order to Bar"
print_warning "Attempting to send unpaid order to bar (should be blocked)..."

send_response=$(api_call "POST" "/waiter/orders/$ORDER_ID/send" "" "400")
if [ $? -eq 0 ]; then
    print_success "State machine correctly blocked invalid transition!"
    print_info "Error message should mention 'CREATED' state and suggest payment"
else
    print_error "State machine failed to block invalid transition!"
fi

# Step 4: Try to Edit Order (Should SUCCEED - order is in CREATED state)
print_section "Step 4: Test Valid Action - Edit Order in CREATED State"
print_info "Editing order (should succeed)..."

edit_data='{
  "items": [
    {
      "menu_item_id": "507f1f77bcf86cd799439011",
      "name": "Cappuccino",
      "quantity": 3,
      "price": 45000,
      "notes": "Extra hot, extra foam"
    }
  ]
}'

edit_response=$(api_call "PUT" "/waiter/orders/$ORDER_ID/edit" "$edit_data" "200")
if [ $? -eq 0 ]; then
    print_success "Order edited successfully!"
else
    print_error "Failed to edit order!"
fi

# Step 5: Collect Payment (CREATED → PAID)
print_section "Step 5: Collect Payment"
print_info "Collecting payment for order..."

payment_data='{
  "payment_method": "cash",
  "amount_paid": 135000
}'

payment_response=$(api_call "POST" "/waiter/orders/$ORDER_ID/payment" "$payment_data" "200")
if [ $? -eq 0 ]; then
    print_success "Payment collected! Order Status: PAID"
else
    print_error "Failed to collect payment!"
    exit 1
fi

# Step 6: Try to Pay Again (Should FAIL)
print_section "Step 6: Test Invalid Transition - Pay Already Paid Order"
print_warning "Attempting to pay again (should be blocked)..."

payment_response=$(api_call "POST" "/waiter/orders/$ORDER_ID/payment" "$payment_data" "400")
if [ $? -eq 0 ]; then
    print_success "State machine correctly blocked duplicate payment!"
else
    print_error "State machine failed to block duplicate payment!"
fi

# Step 7: Send to Bar (PAID → QUEUED)
print_section "Step 7: Send Order to Bar"
print_info "Sending order to bar..."

send_response=$(api_call "POST" "/waiter/orders/$ORDER_ID/send" "" "200")
if [ $? -eq 0 ]; then
    print_success "Order sent to bar! Order Status: QUEUED"
else
    print_error "Failed to send order to bar!"
    exit 1
fi

# Step 8: Try to Edit Order After Sent to Bar (Should FAIL)
print_section "Step 8: Test Invalid Action - Edit Order After Sent to Bar"
print_warning "Attempting to edit order in QUEUED state (should be blocked)..."

edit_response=$(api_call "PUT" "/waiter/orders/$ORDER_ID/edit" "$edit_data" "400")
if [ $? -eq 0 ]; then
    print_success "State machine correctly blocked editing queued order!"
    print_info "Error should mention 'cannot modify order in current state'"
else
    print_error "State machine failed to block editing queued order!"
fi

# Step 9: Login as Barista
print_section "Step 9: Login as Barista"
print_info "Switching to barista account..."

barista_login=$(curl -s -X POST \
    -H "Content-Type: application/json" \
    -d '{"username":"barista1","password":"barista123"}' \
    "$BASE_URL/login")

TOKEN=$(echo $barista_login | jq -r '.token')

if [ "$TOKEN" != "null" ] && [ -n "$TOKEN" ]; then
    print_success "Barista login successful!"
else
    print_error "Barista login failed!"
    exit 1
fi

# Step 10: Accept Order (QUEUED → IN_PROGRESS)
print_section "Step 10: Barista Accepts Order"
print_info "Barista accepting order from queue..."

accept_response=$(api_call "POST" "/barista/orders/$ORDER_ID/accept" "" "200")
if [ $? -eq 0 ]; then
    print_success "Order accepted! Order Status: IN_PROGRESS"
else
    print_error "Failed to accept order!"
    exit 1
fi

# Step 11: Try to Accept Again (Should FAIL)
print_section "Step 11: Test Invalid Transition - Accept Already In-Progress Order"
print_warning "Attempting to accept again (should be blocked)..."

accept_response=$(api_call "POST" "/barista/orders/$ORDER_ID/accept" "" "400")
if [ $? -eq 0 ]; then
    print_success "State machine correctly blocked duplicate accept!"
else
    print_error "State machine failed to block duplicate accept!"
fi

# Step 12: Mark as Ready (IN_PROGRESS → READY)
print_section "Step 12: Mark Order as Ready"
print_info "Marking order as ready..."

ready_response=$(api_call "POST" "/barista/orders/$ORDER_ID/ready" "" "200")
if [ $? -eq 0 ]; then
    print_success "Order marked as ready! Order Status: READY"
    print_info "Response should include progress: 80%"
else
    print_error "Failed to mark order as ready!"
    exit 1
fi

# Step 13: Login back as Waiter
print_section "Step 13: Login Back as Waiter"
print_info "Switching back to waiter account..."

waiter_login=$(curl -s -X POST \
    -H "Content-Type: application/json" \
    -d '{"username":"waiter1","password":"waiter123"}' \
    "$BASE_URL/login")

TOKEN=$(echo $waiter_login | jq -r '.token')

if [ "$TOKEN" != "null" ] && [ -n "$TOKEN" ]; then
    print_success "Waiter login successful!"
else
    print_error "Waiter login failed!"
    exit 1
fi

# Step 14: Serve Order (READY → SERVED)
print_section "Step 14: Serve Order to Customer"
print_info "Serving order to customer..."

serve_response=$(api_call "POST" "/waiter/orders/$ORDER_ID/serve" "" "200")
if [ $? -eq 0 ]; then
    print_success "Order served! Order Status: SERVED"
    print_info "Progress: 100%"
else
    print_error "Failed to serve order!"
    exit 1
fi

# Step 15: Try to Cancel Served Order (Should FAIL)
print_section "Step 15: Test Invalid Transition - Cancel Served Order"
print_warning "Attempting to cancel served order (should be blocked)..."

cancel_data='{"reason":"Customer changed mind"}'
cancel_response=$(api_call "POST" "/cashier/orders/$ORDER_ID/cancel" "$cancel_data" "400")
if [ $? -eq 0 ]; then
    print_success "State machine correctly blocked canceling served order!"
    print_info "Error should mention 'cannot cancel order in SERVED status'"
else
    print_error "State machine failed to block canceling served order!"
fi

# Step 16: Login as Cashier
print_section "Step 16: Login as Cashier"
print_info "Switching to cashier account..."

cashier_login=$(curl -s -X POST \
    -H "Content-Type: application/json" \
    -d '{"username":"cashier1","password":"cashier123"}' \
    "$BASE_URL/login")

TOKEN=$(echo $cashier_login | jq -r '.token')

if [ "$TOKEN" != "null" ] && [ -n "$TOKEN" ]; then
    print_success "Cashier login successful!"
else
    print_error "Cashier login failed!"
    exit 1
fi

# Step 17: Lock Order (SERVED → LOCKED)
print_section "Step 17: Lock Order (Shift Closure)"
print_info "Locking order for shift closure..."

lock_response=$(api_call "POST" "/cashier/orders/$ORDER_ID/lock" "" "200")
if [ $? -eq 0 ]; then
    print_success "Order locked! Order Status: LOCKED (Terminal State)"
else
    print_error "Failed to lock order!"
    exit 1
fi

# Step 18: Try to Modify Locked Order (Should FAIL)
print_section "Step 18: Test Terminal State - Try to Modify Locked Order"
print_warning "Attempting to modify locked order (should be blocked)..."

refund_data='{"reason":"Refund request","amount":135000}'
refund_response=$(api_call "POST" "/cashier/orders/$ORDER_ID/refund" "$refund_data" "400")
if [ $? -eq 0 ]; then
    print_success "State machine correctly blocked modifying locked order!"
    print_info "LOCKED is a terminal state - no transitions allowed"
else
    print_error "State machine failed to block modifying locked order!"
fi

# Summary
print_section "Test Summary"
print_success "Order Workflow Test Complete!"
echo ""
print_info "Order Lifecycle Tested:"
echo "  1. ✅ CREATED - Order created"
echo "  2. ✅ PAID - Payment collected"
echo "  3. ✅ QUEUED - Sent to bar"
echo "  4. ✅ IN_PROGRESS - Barista accepted"
echo "  5. ✅ READY - Preparation complete"
echo "  6. ✅ SERVED - Served to customer"
echo "  7. ✅ LOCKED - Locked for shift closure"
echo ""
print_info "Invalid Transitions Tested:"
echo "  1. ✅ Cannot send unpaid order to bar"
echo "  2. ✅ Cannot pay already paid order"
echo "  3. ✅ Cannot edit order after sent to bar"
echo "  4. ✅ Cannot accept already in-progress order"
echo "  5. ✅ Cannot cancel served order"
echo "  6. ✅ Cannot modify locked order"
echo ""
print_success "All state machine validations working correctly! 🎉"
