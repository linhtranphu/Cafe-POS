#!/bin/bash

# Test script for order-printing integration
# This script verifies that the backend core functionality works

set -e

echo "========================================="
echo "Testing Order-Printing Backend Core"
echo "========================================="
echo ""

# 1. Check if backend compiles
echo "✓ Step 1: Checking if backend compiles..."
cd backend
go build -o /tmp/backend-test ./main.go
if [ $? -eq 0 ]; then
    echo "  ✅ Backend compiles successfully"
else
    echo "  ❌ Backend compilation failed"
    exit 1
fi
echo ""

# 2. Run label printer tests
echo "✓ Step 2: Running label printer tests..."
go test -v ./infrastructure/printing/... > /tmp/label-printer-tests.log 2>&1
if [ $? -eq 0 ]; then
    echo "  ✅ Label printer tests pass"
    grep "PASS" /tmp/label-printer-tests.log | head -5
else
    echo "  ❌ Label printer tests failed"
    cat /tmp/label-printer-tests.log
    exit 1
fi
echo ""

# 3. Check domain entities exist
echo "✓ Step 3: Checking domain entities..."
if [ -f "domain/printing/print_job.go" ] && \
   [ -f "domain/printing/printer_config.go" ] && \
   [ -f "domain/printing/print_template.go" ]; then
    echo "  ✅ All domain entities exist"
else
    echo "  ❌ Some domain entities are missing"
    exit 1
fi
echo ""

# 4. Check repositories exist
echo "✓ Step 4: Checking repositories..."
if [ -f "infrastructure/mongodb/print_job_repository.go" ] && \
   [ -f "infrastructure/mongodb/printer_config_repository.go" ] && \
   [ -f "infrastructure/mongodb/print_template_repository.go" ]; then
    echo "  ✅ All repositories exist"
else
    echo "  ❌ Some repositories are missing"
    exit 1
fi
echo ""

# 5. Check services exist
echo "✓ Step 5: Checking services..."
if [ -f "application/services/print_service.go" ] && \
   [ -f "application/services/print_worker.go" ] && \
   [ -f "application/services/printer_manager.go" ] && \
   [ -f "application/services/template_renderer.go" ]; then
    echo "  ✅ All services exist"
else
    echo "  ❌ Some services are missing"
    exit 1
fi
echo ""

# 6. Check integration with order service
echo "✓ Step 6: Checking order service integration..."
if grep -q "CreatePrintJobsForOrder" application/services/order_service.go && \
   grep -q "printService" application/services/order_service.go && \
   grep -q "SetPrintService" application/services/order_service.go; then
    echo "  ✅ Order service is integrated with print service"
else
    echo "  ❌ Order service integration is incomplete"
    exit 1
fi
echo ""

# 7. Check if print service is wired in main.go
echo "✓ Step 7: Checking main.go wiring..."
if grep -q "NewPrintService" main.go && \
   grep -q "SetPrintService" main.go; then
    echo "  ✅ Print service is wired in main.go"
else
    echo "  ❌ Print service is not properly wired in main.go"
    exit 1
fi
echo ""

# 8. Check migration script exists
echo "✓ Step 8: Checking migration script..."
if [ -f "cmd/migrate/create_printing_collections.go" ]; then
    echo "  ✅ Migration script exists"
else
    echo "  ❌ Migration script is missing"
    exit 1
fi
echo ""

echo "========================================="
echo "✅ All Backend Core Checks Passed!"
echo "========================================="
echo ""
echo "Summary:"
echo "  ✓ Backend compiles without errors"
echo "  ✓ Label printer tests pass"
echo "  ✓ Domain entities exist"
echo "  ✓ Repositories exist"
echo "  ✓ Services exist"
echo "  ✓ Order service integration complete"
echo "  ✓ Main.go wiring complete"
echo "  ✓ Migration script exists"
echo ""
echo "Next Steps:"
echo "  1. Start MongoDB (if not running)"
echo "  2. Run migration: go run cmd/migrate/main.go"
echo "  3. Start backend: go run main.go"
echo "  4. Create an order and mark it as PAID"
echo "  5. Verify print jobs are created in the database"
echo ""
