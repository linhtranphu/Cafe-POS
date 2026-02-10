#!/bin/bash

# Test script to verify CalculateShiftOrderCosts implementation
# This script runs only the new tests for shift order cost calculation

cd "$(dirname "$0")"

echo "Running CalculateShiftOrderCosts tests..."
go test -v -run "^TestCalculateShiftOrderCosts" ./application/services/ 2>&1 | grep -E "(PASS|FAIL|RUN|===|---)"

exit_code=${PIPESTATUS[0]}

if [ $exit_code -eq 0 ]; then
    echo ""
    echo "✓ All CalculateShiftOrderCosts tests passed!"
else
    echo ""
    echo "✗ Some tests failed or build errors occurred"
fi

exit $exit_code
