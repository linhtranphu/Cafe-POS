#!/bin/bash

# Cash Handover API Integration Test Script
# Tests the complete handover workflow from waiter to cashier

set -e

# Colors for output
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
NC='\033[0m' # No Color

# Configuration
API_URL="${API_URL:-http://localhost:8080/api}"
WAITER_TOKEN=""
CASHIER_TOKEN=""
MANAGER_TOKEN=""
WAITER_SHIFT_ID=""
CASHIER_SHIFT_ID=""
HANDOVER_ID=""

echo -e "${YELLOW}=== Cash Handover API Integration Tests ===${NC}\n"

# Helper functions
print_success() {
    echo -e "${GREEN}✓ $1${NC}"
}

print_error() {
    echo -e "${RED}✗ $1${NC}"
}

print_info() {
    echo -e "${YELLOW}ℹ $1${NC}"
}

# Test 1: Login as waiter
test_waiter_login() {
    print_info "Test 1: Login as waiter..."
    
    RESPONSE=$(curl -s -X POST "$API_URL/login" \
        -H "Content-Type: application/json" \
        -d '{
            "username": "waiter1",
            "password": "password123"
        }')
    
    WAITER_TOKEN=$(echo $RESPONSE | jq -r '.token')
    
    if [ "$WAITER_TOKEN" != "null" ] && [ -n "$WAITER_TOKEN" ]; then
        print_success "Waiter login successful"
    else
        print_error "Waiter login failed"
        exit 1
    fi
}

# Test 2: Login as cashier
test_cashier_login() {
    print_info "Test 2: Login as cashier..."
    
    RESPONSE=$(curl -s -X POST "$API_URL/login" \
        -H "Content-Type: application/json" \
        -d '{
            "username": "cashier1",
            "password": "password123"
        }')
    
    CASHIER_TOKEN=$(echo $RESPONSE | jq -r '.token')
    
    if [ "$CASHIER_TOKEN" != "null" ] && [ -n "$CASHIER_TOKEN" ]; then
        print_success "Cashier login successful"
    else
        print_error "Cashier login failed"
        exit 1
    fi
}

# Test 3: Start waiter shift
test_start_waiter_shift() {
    print_info "Test 3: Start waiter shift..."
    
    RESPONSE=$(curl -s -X POST "$API_URL/shifts" \
        -H "Content-Type: application/json" \
        -H "Authorization: Bearer $WAITER_TOKEN" \
        -d '{
            "type": "MORNING",
            "start_cash": 100000
        }')
    
    WAITER_SHIFT_ID=$(echo $RESPONSE | jq -r '.id')
    
    if [ "$WAITER_SHIFT_ID" != "null" ] && [ -n "$WAITER_SHIFT_ID" ]; then
        print_success "Waiter shift started: $WAITER_SHIFT_ID"
    else
        print_error "Failed to start waiter shift"
        exit 1
    fi
}

# Test 4: Start cashier shift
test_start_cashier_shift() {
    print_info "Test 4: Start cashier shift..."
    
    RESPONSE=$(curl -s -X POST "$API_URL/cashier-shifts" \
        -H "Content-Type: application/json" \
        -H "Authorization: Bearer $CASHIER_TOKEN" \
        -d '{
            "type": "MORNING",
            "start_cash": 500000
        }')
    
    CASHIER_SHIFT_ID=$(echo $RESPONSE | jq -r '.id')
    
    if [ "$CASHIER_SHIFT_ID" != "null" ] && [ -n "$CASHIER_SHIFT_ID" ]; then
        print_success "Cashier shift started: $CASHIER_SHIFT_ID"
    else
        print_error "Failed to start cashier shift"
        exit 1
    fi
}

# Test 5: Create partial handover
test_create_partial_handover() {
    print_info "Test 5: Create partial handover..."
    
    RESPONSE=$(curl -s -X POST "$API_URL/shifts/$WAITER_SHIFT_ID/handover" \
        -H "Content-Type: application/json" \
        -H "Authorization: Bearer $WAITER_TOKEN" \
        -d '{
            "declared_amount": 200000,
            "handover_type": "PARTIAL",
            "waiter_note": "Test partial handover"
        }')
    
    HANDOVER_ID=$(echo $RESPONSE | jq -r '.id')
    
    if [ "$HANDOVER_ID" != "null" ] && [ -n "$HANDOVER_ID" ]; then
        print_success "Partial handover created: $HANDOVER_ID"
    else
        print_error "Failed to create partial handover"
        echo "Response: $RESPONSE"
        exit 1
    fi
}

# Test 6: Get pending handover (waiter)
test_get_pending_handover() {
    print_info "Test 6: Get pending handover (waiter)..."
    
    RESPONSE=$(curl -s -X GET "$API_URL/shifts/$WAITER_SHIFT_ID/pending-handover" \
        -H "Authorization: Bearer $WAITER_TOKEN")
    
    PENDING_ID=$(echo $RESPONSE | jq -r '.id')
    
    if [ "$PENDING_ID" == "$HANDOVER_ID" ]; then
        print_success "Pending handover retrieved correctly"
    else
        print_error "Failed to get pending handover"
        exit 1
    fi
}

# Test 7: Get pending handovers (cashier)
test_get_pending_handovers_cashier() {
    print_info "Test 7: Get pending handovers (cashier)..."
    
    RESPONSE=$(curl -s -X GET "$API_URL/cash-handovers/pending" \
        -H "Authorization: Bearer $CASHIER_TOKEN")
    
    COUNT=$(echo $RESPONSE | jq '. | length')
    
    if [ "$COUNT" -gt 0 ]; then
        print_success "Cashier can see $COUNT pending handover(s)"
    else
        print_error "Cashier cannot see pending handovers"
        exit 1
    fi
}

# Test 8: Confirm handover without discrepancy
test_confirm_handover_no_discrepancy() {
    print_info "Test 8: Confirm handover without discrepancy..."
    
    RESPONSE=$(curl -s -X POST "$API_URL/cash-handovers/$HANDOVER_ID/confirm" \
        -H "Content-Type: application/json" \
        -H "Authorization: Bearer $CASHIER_TOKEN" \
        -d '{
            "actual_amount": 200000,
            "status": "CONFIRMED",
            "cashier_note": "Confirmed - exact amount"
        }')
    
    if echo "$RESPONSE" | jq -e '.message' > /dev/null 2>&1; then
        print_success "Handover confirmed without discrepancy"
    else
        print_error "Failed to confirm handover"
        echo "Response: $RESPONSE"
        exit 1
    fi
}

# Test 9: Create handover with discrepancy
test_create_handover_with_discrepancy() {
    print_info "Test 9: Create handover with discrepancy..."
    
    # Create another handover
    RESPONSE=$(curl -s -X POST "$API_URL/shifts/$WAITER_SHIFT_ID/handover" \
        -H "Content-Type: application/json" \
        -H "Authorization: Bearer $WAITER_TOKEN" \
        -d '{
            "declared_amount": 150000,
            "handover_type": "PARTIAL",
            "waiter_note": "Test handover with discrepancy"
        }')
    
    HANDOVER_ID_2=$(echo $RESPONSE | jq -r '.id')
    
    if [ "$HANDOVER_ID_2" != "null" ] && [ -n "$HANDOVER_ID_2" ]; then
        print_success "Second handover created: $HANDOVER_ID_2"
        
        # Confirm with discrepancy
        RESPONSE=$(curl -s -X POST "$API_URL/cash-handovers/$HANDOVER_ID_2/confirm" \
            -H "Content-Type: application/json" \
            -H "Authorization: Bearer $CASHIER_TOKEN" \
            -d '{
                "actual_amount": 145000,
                "status": "CONFIRMED",
                "cashier_note": "Short 5k",
                "discrepancy_reason": "COUNTING_ERROR",
                "discrepancy_responsibility": "WAITER"
            }')
        
        if echo "$RESPONSE" | jq -e '.message' > /dev/null 2>&1; then
            print_success "Handover confirmed with discrepancy"
        else
            print_error "Failed to confirm handover with discrepancy"
            exit 1
        fi
    else
        print_error "Failed to create second handover"
        exit 1
    fi
}

# Test 10: Get handover history
test_get_handover_history() {
    print_info "Test 10: Get handover history..."
    
    RESPONSE=$(curl -s -X GET "$API_URL/shifts/$WAITER_SHIFT_ID/handovers" \
        -H "Authorization: Bearer $WAITER_TOKEN")
    
    COUNT=$(echo $RESPONSE | jq '. | length')
    
    if [ "$COUNT" -ge 2 ]; then
        print_success "Handover history retrieved: $COUNT handovers"
    else
        print_error "Failed to get handover history"
        exit 1
    fi
}

# Test 11: Reject handover
test_reject_handover() {
    print_info "Test 11: Reject handover..."
    
    # Create another handover
    RESPONSE=$(curl -s -X POST "$API_URL/shifts/$WAITER_SHIFT_ID/handover" \
        -H "Content-Type: application/json" \
        -H "Authorization: Bearer $WAITER_TOKEN" \
        -d '{
            "declared_amount": 50000,
            "handover_type": "PARTIAL",
            "waiter_note": "Test rejection"
        }')
    
    HANDOVER_ID_3=$(echo $RESPONSE | jq -r '.id')
    
    if [ "$HANDOVER_ID_3" != "null" ] && [ -n "$HANDOVER_ID_3" ]; then
        # Reject it
        RESPONSE=$(curl -s -X POST "$API_URL/cash-handovers/$HANDOVER_ID_3/confirm" \
            -H "Content-Type: application/json" \
            -H "Authorization: Bearer $CASHIER_TOKEN" \
            -d '{
                "status": "REJECTED",
                "cashier_note": "Amount does not match records"
            }')
        
        if echo "$RESPONSE" | jq -e '.message' > /dev/null 2>&1; then
            print_success "Handover rejected successfully"
        else
            print_error "Failed to reject handover"
            exit 1
        fi
    else
        print_error "Failed to create handover for rejection test"
        exit 1
    fi
}

# Test 12: Cancel handover (waiter)
test_cancel_handover() {
    print_info "Test 12: Cancel handover (waiter)..."
    
    # Create another handover
    RESPONSE=$(curl -s -X POST "$API_URL/shifts/$WAITER_SHIFT_ID/handover" \
        -H "Content-Type: application/json" \
        -H "Authorization: Bearer $WAITER_TOKEN" \
        -d '{
            "declared_amount": 30000,
            "handover_type": "PARTIAL",
            "waiter_note": "Test cancellation"
        }')
    
    HANDOVER_ID_4=$(echo $RESPONSE | jq -r '.id')
    
    if [ "$HANDOVER_ID_4" != "null" ] && [ -n "$HANDOVER_ID_4" ]; then
        # Cancel it
        RESPONSE=$(curl -s -X DELETE "$API_URL/cash-handovers/$HANDOVER_ID_4" \
            -H "Authorization: Bearer $WAITER_TOKEN")
        
        if echo "$RESPONSE" | jq -e '.message' > /dev/null 2>&1; then
            print_success "Handover cancelled successfully"
        else
            print_error "Failed to cancel handover"
            exit 1
        fi
    else
        print_error "Failed to create handover for cancellation test"
        exit 1
    fi
}

# Run all tests
main() {
    echo -e "${YELLOW}Starting API integration tests...${NC}\n"
    
    test_waiter_login
    test_cashier_login
    test_start_waiter_shift
    test_start_cashier_shift
    test_create_partial_handover
    test_get_pending_handover
    test_get_pending_handovers_cashier
    test_confirm_handover_no_discrepancy
    test_create_handover_with_discrepancy
    test_get_handover_history
    test_reject_handover
    test_cancel_handover
    
    echo -e "\n${GREEN}=== All tests passed! ===${NC}"
}

# Check if jq is installed
if ! command -v jq &> /dev/null; then
    print_error "jq is not installed. Please install jq to run this script."
    exit 1
fi

# Run tests
main
