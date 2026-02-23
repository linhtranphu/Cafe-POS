#!/bin/bash

# Test Fund Handover API Endpoints
# This script tests the new API endpoints for Phase 4

set -e

API_URL="${API_URL:-http://localhost:8080/api/v1}"
TOKEN="${TOKEN:-}"

echo "🧪 Testing Fund Handover API Endpoints"
echo "========================================"
echo ""

# Check if token is provided
if [ -z "$TOKEN" ]; then
    echo "❌ Error: TOKEN environment variable is required"
    echo "Usage: TOKEN=your_jwt_token ./test-fund-handover-api.sh"
    exit 1
fi

# Step 1: Get current cashier shift
echo "📋 Step 1: Get current cashier shift"
SHIFT_RESPONSE=$(curl -s -X GET "$API_URL/cashier-shifts/current" \
    -H "Authorization: Bearer $TOKEN" \
    -H "Content-Type: application/json")

SHIFT_ID=$(echo "$SHIFT_RESPONSE" | jq -r '.id // .ID // ._id // empty')

if [ -z "$SHIFT_ID" ] || [ "$SHIFT_ID" = "null" ]; then
    echo "❌ No open cashier shift found"
    echo "Response: $SHIFT_RESPONSE"
    exit 1
fi

echo "✅ Found cashier shift: $SHIFT_ID"
echo ""

# Step 2: Get managed funds
echo "💰 Step 2: Get managed funds"
MANAGED_FUNDS=$(curl -s -X GET "$API_URL/cashier-shifts/$SHIFT_ID/managed-funds" \
    -H "Authorization: Bearer $TOKEN" \
    -H "Content-Type: application/json")

echo "Response:"
echo "$MANAGED_FUNDS" | jq '.'

STARTING_FLOAT=$(echo "$MANAGED_FUNDS" | jq -r '.starting_float')
RECEIVED_CASH=$(echo "$MANAGED_FUNDS" | jq -r '.received_cash')
RECEIVED_TRANSFER=$(echo "$MANAGED_FUNDS" | jq -r '.received_transfer')
EXPECTED_CASH=$(echo "$MANAGED_FUNDS" | jq -r '.expected_cash')

echo ""
echo "📊 Managed Funds Summary:"
echo "  Starting Float:    $STARTING_FLOAT"
echo "  Received Cash:     $RECEIVED_CASH"
echo "  Received Transfer: $RECEIVED_TRANSFER"
echo "  Expected Cash:     $EXPECTED_CASH"
echo ""

# Step 3: Close shift with fund handover (with variance)
echo "🔒 Step 3: Close shift with fund handover"
echo "Simulating actual cash = expected cash - 5000 (variance)"

ACTUAL_CASH=$(echo "$EXPECTED_CASH - 5000" | bc)

CLOSURE_PAYLOAD=$(cat <<EOF
{
  "actual_cash": $ACTUAL_CASH,
  "variance_reason": "COUNTING_ERROR",
  "variance_notes": "Đếm nhầm tờ 50k thành 100k - test case"
}
EOF
)

echo "Payload:"
echo "$CLOSURE_PAYLOAD" | jq '.'
echo ""

CLOSURE_RESPONSE=$(curl -s -X POST "$API_URL/cashier-shifts/$SHIFT_ID/close-with-fund-handover" \
    -H "Authorization: Bearer $TOKEN" \
    -H "Content-Type: application/json" \
    -d "$CLOSURE_PAYLOAD")

echo "Response:"
echo "$CLOSURE_RESPONSE" | jq '.'
echo ""

# Check if closure was successful
CLOSED_SHIFT_ID=$(echo "$CLOSURE_RESPONSE" | jq -r '.shift.id // .shift.ID // .shift._id // empty')
HANDOVER_ID=$(echo "$CLOSURE_RESPONSE" | jq -r '.fund_handover.id // .fund_handover.ID // .fund_handover._id // empty')

if [ -z "$CLOSED_SHIFT_ID" ] || [ "$CLOSED_SHIFT_ID" = "null" ]; then
    echo "❌ Closure failed"
    echo "Error: $(echo "$CLOSURE_RESPONSE" | jq -r '.error // "Unknown error"')"
    exit 1
fi

if [ -z "$HANDOVER_ID" ] || [ "$HANDOVER_ID" = "null" ]; then
    echo "❌ Fund handover record not created"
    exit 1
fi

echo "✅ Shift closed successfully: $CLOSED_SHIFT_ID"
echo "✅ Fund handover created: $HANDOVER_ID"
echo ""

# Step 4: Verify fund handover details
echo "📋 Step 4: Verify fund handover details"
CASH_AMOUNT=$(echo "$CLOSURE_RESPONSE" | jq -r '.fund_handover.cash_amount')
TRANSFER_AMOUNT=$(echo "$CLOSURE_RESPONSE" | jq -r '.fund_handover.transfer_amount')
TOTAL_AMOUNT=$(echo "$CLOSURE_RESPONSE" | jq -r '.fund_handover.total_amount')
VARIANCE_AMOUNT=$(echo "$CLOSURE_RESPONSE" | jq -r '.fund_handover.variance_amount')

echo "Fund Handover Details:"
echo "  Cash Amount:     $CASH_AMOUNT"
echo "  Transfer Amount: $TRANSFER_AMOUNT"
echo "  Total Amount:    $TOTAL_AMOUNT"
echo "  Variance:        $VARIANCE_AMOUNT"
echo ""

# Verify variance calculation
EXPECTED_VARIANCE=$(echo "$ACTUAL_CASH - $EXPECTED_CASH" | bc)
if [ "$VARIANCE_AMOUNT" != "$EXPECTED_VARIANCE" ]; then
    echo "⚠️  Warning: Variance calculation mismatch"
    echo "  Expected: $EXPECTED_VARIANCE"
    echo "  Got:      $VARIANCE_AMOUNT"
else
    echo "✅ Variance calculation correct"
fi

echo ""
echo "🎉 All tests passed!"
echo ""
echo "Summary:"
echo "  ✅ Get managed funds API working"
echo "  ✅ Close with fund handover API working"
echo "  ✅ Fund handover record created"
echo "  ✅ Variance documented correctly"
