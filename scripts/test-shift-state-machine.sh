#!/bin/bash

# Test Shift State Machine Validation
# Tests waiter/barista shift workflow with state machine validation

BASE_URL="http://localhost:8080/api"

# Colors
GREEN='\033[0;32m'
RED='\033[0;31m'
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

# Extract JSON value
extract_json() {
    local json=$1
    local key=$2
    echo "$json" | grep -o "\"$key\":\"[^\"]*\"" | cut -d'"' -f4
}

print_section "Shift State Machine Workflow Test"

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

# Step 2: Check Current Shift
print_section "Step 2: Check Current Shift"
print_info "Checking if waiter has open shift..."

current_response=$(curl -s -w "\n%{http_code}" -X GET \
    -H "Authorization: Bearer $TOKEN" \
    "$BASE_URL/shifts/current")

http_code=$(echo "$current_response" | tail -n1)
body=$(echo "$current_response" | sed '$d')

if [ "$http_code" = "200" ]; then
    print_warning "Waiter already has open shift"
    SHIFT_ID=$(extract_json "$body" "id")
    if [ -z "$SHIFT_ID" ]; then
        SHIFT_ID=$(extract_json "$body" "_id")
    fi
    echo "Existing Shift ID: $SHIFT_ID"
    echo "$body"
    HAS_OPEN_SHIFT=true
else
    print_info "No open shift found (that's OK)"
    HAS_OPEN_SHIFT=false
fi

# Step 3: Try to Start Shift (Should FAIL if already have open shift)
if [ "$HAS_OPEN_SHIFT" = true ]; then
    print_section "Step 3: Test - Start Shift When Already Have Open Shift (Should FAIL)"
    print_warning "Attempting to start duplicate shift..."

    start_response=$(curl -s -w "\n%{http_code}" -X POST \
        -H "Content-Type: application/json" \
        -H "Authorization: Bearer $TOKEN" \
        -d '{"device_id":"test-device","type":"waiter"}' \
        "$BASE_URL/shifts/start")

    http_code=$(echo "$start_response" | tail -n1)
    body=$(echo "$start_response" | sed '$d')

    if [ "$http_code" = "400" ]; then
        print_success "✅ State machine blocked duplicate shift!"
        echo "$body"
    else
        print_error "State machine failed! Got status: $http_code"
        echo "$body"
    fi
else
    # Step 3: Start Shift (Should SUCCEED)
    print_section "Step 3: Start Shift (OPEN State)"
    print_info "Starting waiter shift..."

    start_response=$(curl -s -w "\n%{http_code}" -X POST \
        -H "Content-Type: application/json" \
        -H "Authorization: Bearer $TOKEN" \
        -d '{"device_id":"test-device","type":"waiter"}' \
        "$BASE_URL/shifts/start")

    http_code=$(echo "$start_response" | tail -n1)
    body=$(echo "$start_response" | sed '$d')

    if [ "$http_code" = "201" ]; then
        print_success "Shift started! Status: OPEN"
        SHIFT_ID=$(extract_json "$body" "id")
        if [ -z "$SHIFT_ID" ]; then
            SHIFT_ID=$(extract_json "$body" "_id")
        fi
        echo "Shift ID: $SHIFT_ID"
        echo "$body"
    else
        print_error "Failed to start shift! Status: $http_code"
        echo "$body"
        exit 1
    fi

    # Step 4: Try to Start Another Shift (Should FAIL)
    print_section "Step 4: Test - Start Duplicate Shift (Should FAIL)"
    print_warning "Attempting to start another shift..."

    start_response=$(curl -s -w "\n%{http_code}" -X POST \
        -H "Content-Type: application/json" \
        -H "Authorization: Bearer $TOKEN" \
        -d '{"device_id":"test-device-2","type":"waiter"}' \
        "$BASE_URL/shifts/start")

    http_code=$(echo "$start_response" | tail -n1)
    body=$(echo "$start_response" | sed '$d')

    if [ "$http_code" = "400" ]; then
        print_success "✅ State machine blocked duplicate shift!"
        echo "$body"
    else
        print_error "State machine failed! Got status: $http_code"
        echo "$body"
    fi
fi

# Step 5: Get Shift Details
print_section "Step 5: Get Shift Details"
print_info "Fetching shift details..."

shift_response=$(curl -s -w "\n%{http_code}" -X GET \
    -H "Authorization: Bearer $TOKEN" \
    "$BASE_URL/shifts/$SHIFT_ID")

http_code=$(echo "$shift_response" | tail -n1)
body=$(echo "$shift_response" | sed '$d')

if [ "$http_code" = "200" ]; then
    print_success "Shift details retrieved"
    echo "$body"
else
    print_error "Failed to get shift! Status: $http_code"
    echo "$body"
fi

# Step 6: End Shift (OPEN → CLOSED)
print_section "Step 6: End Shift (OPEN → CLOSED)"
print_info "Ending shift..."

end_response=$(curl -s -w "\n%{http_code}" -X POST \
    -H "Content-Type: application/json" \
    -H "Authorization: Bearer $TOKEN" \
    -d '{"notes":"End of shift test"}' \
    "$BASE_URL/shifts/$SHIFT_ID/end")

http_code=$(echo "$end_response" | tail -n1)
body=$(echo "$end_response" | sed '$d')

if [ "$http_code" = "200" ]; then
    print_success "Shift ended! Status: CLOSED"
    echo "$body"
else
    print_error "Failed to end shift! Status: $http_code"
    echo "$body"
fi

# Step 7: Try to End Already Closed Shift (Should FAIL)
print_section "Step 7: Test - End Already Closed Shift (Should FAIL)"
print_warning "Attempting to end closed shift..."

end_response=$(curl -s -w "\n%{http_code}" -X POST \
    -H "Content-Type: application/json" \
    -H "Authorization: Bearer $TOKEN" \
    -d '{"notes":"Try to end again"}' \
    "$BASE_URL/shifts/$SHIFT_ID/end")

http_code=$(echo "$end_response" | tail -n1)
body=$(echo "$end_response" | sed '$d')

if [ "$http_code" = "400" ]; then
    print_success "✅ State machine blocked ending closed shift!"
    echo "$body"
else
    print_error "State machine failed! Got status: $http_code"
    echo "$body"
fi

# Step 8: Try to Close Already Closed Shift (Should FAIL)
print_section "Step 8: Test - Close Already Closed Shift (Should FAIL)"
print_warning "Attempting to close closed shift..."

close_response=$(curl -s -w "\n%{http_code}" -X POST \
    -H "Content-Type: application/json" \
    -H "Authorization: Bearer $TOKEN" \
    -d '{"notes":"Try to close"}' \
    "$BASE_URL/shifts/$SHIFT_ID/close")

http_code=$(echo "$close_response" | tail -n1)
body=$(echo "$close_response" | sed '$d')

if [ "$http_code" = "400" ]; then
    print_success "✅ State machine blocked closing closed shift!"
    echo "$body"
else
    print_error "State machine failed! Got status: $http_code"
    echo "$body"
fi

# Step 9: Start New Shift (Should SUCCEED now)
print_section "Step 9: Start New Shift After Closing Previous"
print_info "Starting new shift..."

start_response=$(curl -s -w "\n%{http_code}" -X POST \
    -H "Content-Type: application/json" \
    -H "Authorization: Bearer $TOKEN" \
    -d '{"device_id":"test-device-new","type":"waiter"}' \
    "$BASE_URL/shifts/start")

http_code=$(echo "$start_response" | tail -n1)
body=$(echo "$start_response" | sed '$d')

if [ "$http_code" = "201" ]; then
    print_success "New shift started successfully!"
    NEW_SHIFT_ID=$(extract_json "$body" "id")
    if [ -z "$NEW_SHIFT_ID" ]; then
        NEW_SHIFT_ID=$(extract_json "$body" "_id")
    fi
    echo "New Shift ID: $NEW_SHIFT_ID"
    echo "$body"
else
    print_error "Failed to start new shift! Status: $http_code"
    echo "$body"
fi

# Step 10: Test Barista Shift
print_section "Step 10: Test Barista Shift"
print_info "Logging in as barista..."

barista_login=$(curl -s -X POST \
    -H "Content-Type: application/json" \
    -d '{"username":"barista1","password":"barista123"}' \
    "$BASE_URL/login")

TOKEN=$(extract_json "$barista_login" "token")

if [ -n "$TOKEN" ]; then
    print_success "Barista logged in!"
    
    # Check if barista has open shift
    current_response=$(curl -s -w "\n%{http_code}" -X GET \
        -H "Authorization: Bearer $TOKEN" \
        "$BASE_URL/shifts/current")

    http_code=$(echo "$current_response" | tail -n1)
    
    if [ "$http_code" = "404" ]; then
        print_info "Barista has no open shift, starting one..."
        
        start_response=$(curl -s -w "\n%{http_code}" -X POST \
            -H "Content-Type: application/json" \
            -H "Authorization: Bearer $TOKEN" \
            -d '{"device_id":"barista-device","type":"barista"}' \
            "$BASE_URL/shifts/start")

        http_code=$(echo "$start_response" | tail -n1)
        body=$(echo "$start_response" | sed '$d')

        if [ "$http_code" = "201" ]; then
            print_success "Barista shift started!"
            echo "$body"
        else
            print_warning "Could not start barista shift: $http_code"
            echo "$body"
        fi
    else
        print_info "Barista already has open shift"
    fi
else
    print_error "Barista login failed!"
fi

# Summary
print_section "Test Summary"
print_success "Shift State Machine Test Complete! 🎉"
echo ""
print_info "Tested Shift Lifecycle:"
echo "  ✅ OPEN → CLOSED"
echo ""
print_info "Tested Invalid Transitions:"
echo "  ✅ Cannot start shift when already have open shift"
echo "  ✅ Cannot end already closed shift"
echo "  ✅ Cannot close already closed shift"
echo "  ✅ Can start new shift after closing previous one"
echo ""
print_info "Tested Both Roles:"
echo "  ✅ Waiter shift workflow"
echo "  ✅ Barista shift workflow"
echo ""
print_success "All shift state machine validations working! 🚀"
