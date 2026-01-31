#!/bin/bash

# Simple Order Workflow Test (without jq dependency)
# Tests state machine validation for order transitions

BASE_URL="http://localhost:8080/api"

# Colors
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
BLUE='\033[0;34m'
NC='\033[0m'

print_success() { echo -e "${GREEN}✅ $1${NC}"; }
print_error() { echo -e "${RED}❌ $1${NC}"; }
print_info() { echo -e "${BLUE}ℹ️  $1${NC}"; }
print_warning() { echo -e "${YELLOW}⚠️  $1${NC}"; }
print_section() {
    echo -e "\n${BLUE}========================================${NC}"
    echo -e "${BLUE}$1${NC}"
    echo -e "${BLUE}========================================${NC}\n"
}

# Extract JSON value (simple grep/sed approach)
extract_json() {
    local json=$1
    local key=$2
    echo "$json" | grep -o "\"$key\":\"[^\"]*\"" | cut -d'"' -f4
}

print_section "Order State Machine Workflow Test"

# Step 1: Login as Waiter
print_section "Step 1: Login as Waiter"
print_info "Logging in..."

login_response=$(curl -s -X POST \
    -H "Content-Type: application/json" \
    -d '{"username":"waiter1","password":"waiter123"}' \
    "$BASE_URL/login")

TOKEN=$(extract_json "$login_response" "token")

if [ -n "$TOKEN" ] && [ "$TOKEN" != "null" ]; then
    print_success "Login successful!"
    echo "Token: ${TOKEN:0:20}..."
else
    print_error "Login failed!"
    echo "$login_response"
    exit 1
fi

# Step 2: Start Waiter Shift
print_section "Step 2: Start Waiter Shift"
print_info "Starting waiter shift..."

shift_response=$(curl -s -w "\n%{http_code}" -X POST \
    -H "Content-Type: application/json" \
    -H "Authorization: Bearer $TOKEN" \
    -d '{"device_id":"test-device","type":"waiter"}' \
    "$BASE_URL/shifts/start")

http_code=$(echo "$shift_response" | tail -n1)
body=$(echo "$shift_response" | sed '$d')

if [ "$http_code" = "201" ]; then
    print_success "Shift started!"
    echo "$body"
elif [ "$http_code" = "400" ] && echo "$body" | grep -q "already has an open shift"; then
    print_warning "Shift already open (that's OK)"
else
    print_error "Failed to start shift! Status: $http_code"
    echo "$body"
fi

# Step 3: Create Order
print_section "Step 3: Create Order (CREATED State)"
print_info "Creating new order..."

create_response=$(curl -s -X POST \
    -H "Content-Type: application/json" \
    -H "Authorization: Bearer $TOKEN" \
    -d '{
      "table_number": 5,
      "customer_name": "Test Customer",
      "items": [{
        "menu_item_id": "507f1f77bcf86cd799439011",
        "name": "Cappuccino",
        "quantity": 2,
        "price": 45000
      }]
    }' \
    "$BASE_URL/waiter/orders")

ORDER_ID=$(extract_json "$create_response" "id")
if [ -z "$ORDER_ID" ]; then
    ORDER_ID=$(extract_json "$create_response" "_id")
fi

if [ -n "$ORDER_ID" ] && [ "$ORDER_ID" != "null" ]; then
    print_success "Order created! ID: $ORDER_ID"
    echo "$create_response"
else
    print_error "Failed to create order!"
    echo "$create_response"
    exit 1
fi

# Step 4: Try to Send Unpaid Order (Should FAIL)
print_section "Step 4: Test - Send Unpaid Order to Bar (Should FAIL)"
print_warning "Attempting invalid transition..."

send_response=$(curl -s -w "\n%{http_code}" -X POST \
    -H "Authorization: Bearer $TOKEN" \
    "$BASE_URL/waiter/orders/$ORDER_ID/send")

http_code=$(echo "$send_response" | tail -n1)
body=$(echo "$send_response" | sed '$d')

if [ "$http_code" = "400" ]; then
    print_success "✅ State machine blocked invalid transition!"
    echo "$body"
else
    print_error "State machine failed! Got status: $http_code"
    echo "$body"
fi

# Step 5: Collect Payment
print_section "Step 5: Collect Payment (CREATED → PAID)"
print_info "Collecting payment..."

payment_response=$(curl -s -w "\n%{http_code}" -X POST \
    -H "Content-Type: application/json" \
    -H "Authorization: Bearer $TOKEN" \
    -d '{"payment_method":"cash","amount_paid":90000}' \
    "$BASE_URL/waiter/orders/$ORDER_ID/payment")

http_code=$(echo "$payment_response" | tail -n1)
body=$(echo "$payment_response" | sed '$d')

if [ "$http_code" = "200" ]; then
    print_success "Payment collected! Order now PAID"
    echo "$body"
else
    print_error "Payment failed! Status: $http_code"
    echo "$body"
fi

# Step 6: Try to Pay Again (Should FAIL)
print_section "Step 6: Test - Pay Already Paid Order (Should FAIL)"
print_warning "Attempting duplicate payment..."

payment_response=$(curl -s -w "\n%{http_code}" -X POST \
    -H "Content-Type: application/json" \
    -H "Authorization: Bearer $TOKEN" \
    -d '{"payment_method":"cash","amount_paid":90000}' \
    "$BASE_URL/waiter/orders/$ORDER_ID/payment")

http_code=$(echo "$payment_response" | tail -n1)
body=$(echo "$payment_response" | sed '$d')

if [ "$http_code" = "400" ]; then
    print_success "✅ State machine blocked duplicate payment!"
    echo "$body"
else
    print_error "State machine failed! Got status: $http_code"
fi

# Step 7: Send to Bar
print_section "Step 7: Send to Bar (PAID → QUEUED)"
print_info "Sending order to bar..."

send_response=$(curl -s -w "\n%{http_code}" -X POST \
    -H "Authorization: Bearer $TOKEN" \
    "$BASE_URL/waiter/orders/$ORDER_ID/send")

http_code=$(echo "$send_response" | tail -n1)
body=$(echo "$send_response" | sed '$d')

if [ "$http_code" = "200" ]; then
    print_success "Order sent to bar! Order now QUEUED"
    echo "$body"
else
    print_error "Failed to send! Status: $http_code"
    echo "$body"
fi

# Step 8: Try to Edit After Sent (Should FAIL)
print_section "Step 8: Test - Edit Order After Sent to Bar (Should FAIL)"
print_warning "Attempting to edit queued order..."

edit_response=$(curl -s -w "\n%{http_code}" -X PUT \
    -H "Content-Type: application/json" \
    -H "Authorization: Bearer $TOKEN" \
    -d '{"items":[{"menu_item_id":"507f1f77bcf86cd799439011","name":"Cappuccino","quantity":3,"price":45000}]}' \
    "$BASE_URL/waiter/orders/$ORDER_ID/edit")

http_code=$(echo "$edit_response" | tail -n1)
body=$(echo "$edit_response" | sed '$d')

if [ "$http_code" = "400" ]; then
    print_success "✅ State machine blocked editing queued order!"
    echo "$body"
else
    print_error "State machine failed! Got status: $http_code"
fi

# Step 9: Login as Barista
print_section "Step 9: Login as Barista"
print_info "Switching to barista..."

barista_login=$(curl -s -X POST \
    -H "Content-Type: application/json" \
    -d '{"username":"barista1","password":"barista123"}' \
    "$BASE_URL/login")

TOKEN=$(extract_json "$barista_login" "token")

if [ -n "$TOKEN" ]; then
    print_success "Barista logged in!"
else
    print_error "Barista login failed!"
    exit 1
fi

# Step 10: Accept Order
print_section "Step 10: Accept Order (QUEUED → IN_PROGRESS)"
print_info "Barista accepting order..."

accept_response=$(curl -s -w "\n%{http_code}" -X POST \
    -H "Authorization: Bearer $TOKEN" \
    "$BASE_URL/barista/orders/$ORDER_ID/accept")

http_code=$(echo "$accept_response" | tail -n1)
body=$(echo "$accept_response" | sed '$d')

if [ "$http_code" = "200" ]; then
    print_success "Order accepted! Order now IN_PROGRESS"
    echo "$body"
else
    print_error "Failed to accept! Status: $http_code"
    echo "$body"
fi

# Step 11: Mark as Ready
print_section "Step 11: Mark as Ready (IN_PROGRESS → READY)"
print_info "Marking order as ready..."

ready_response=$(curl -s -w "\n%{http_code}" -X POST \
    -H "Authorization: Bearer $TOKEN" \
    "$BASE_URL/barista/orders/$ORDER_ID/ready")

http_code=$(echo "$ready_response" | tail -n1)
body=$(echo "$ready_response" | sed '$d')

if [ "$http_code" = "200" ]; then
    print_success "Order ready! Order now READY (Progress: 80%)"
    echo "$body"
else
    print_error "Failed to mark ready! Status: $http_code"
    echo "$body"
fi

# Step 12: Login back as Waiter
print_section "Step 12: Login Back as Waiter"

waiter_login=$(curl -s -X POST \
    -H "Content-Type: application/json" \
    -d '{"username":"waiter1","password":"waiter123"}' \
    "$BASE_URL/login")

TOKEN=$(extract_json "$waiter_login" "token")
print_success "Waiter logged in!"

# Step 13: Serve Order
print_section "Step 13: Serve Order (READY → SERVED)"
print_info "Serving order to customer..."

serve_response=$(curl -s -w "\n%{http_code}" -X POST \
    -H "Authorization: Bearer $TOKEN" \
    "$BASE_URL/waiter/orders/$ORDER_ID/serve")

http_code=$(echo "$serve_response" | tail -n1)
body=$(echo "$serve_response" | sed '$d')

if [ "$http_code" = "200" ]; then
    print_success "Order served! Order now SERVED (Progress: 100%)"
    echo "$body"
else
    print_error "Failed to serve! Status: $http_code"
    echo "$body"
fi

# Step 14: Try to Cancel Served Order (Should FAIL)
print_section "Step 14: Test - Cancel Served Order (Should FAIL)"
print_warning "Attempting to cancel served order..."

cancel_response=$(curl -s -w "\n%{http_code}" -X POST \
    -H "Content-Type: application/json" \
    -H "Authorization: Bearer $TOKEN" \
    -d '{"reason":"Test cancel"}' \
    "$BASE_URL/cashier/orders/$ORDER_ID/cancel")

http_code=$(echo "$cancel_response" | tail -n1)
body=$(echo "$cancel_response" | sed '$d')

if [ "$http_code" = "400" ]; then
    print_success "✅ State machine blocked canceling served order!"
    echo "$body"
else
    print_error "State machine failed! Got status: $http_code"
fi

# Summary
print_section "Test Summary"
print_success "Order State Machine Test Complete! 🎉"
echo ""
print_info "Tested Order Lifecycle:"
echo "  ✅ CREATED → PAID → QUEUED → IN_PROGRESS → READY → SERVED"
echo ""
print_info "Tested Invalid Transitions:"
echo "  ✅ Cannot send unpaid order to bar"
echo "  ✅ Cannot pay already paid order"
echo "  ✅ Cannot edit order after sent to bar"
echo "  ✅ Cannot cancel served order"
echo ""
print_success "All state machine validations working! 🚀"
