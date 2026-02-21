#!/bin/bash
# Test script for print service only

cd "$(dirname "$0")"

echo "Running print service tests..."
go test -v -run "TestCreatePrintJobsForOrder|TestReprintBill|TestReprintLabel|TestGetPendingJobs|TestGetFailedJobs|TestRetryJob|TestCancelJob" \
  -tags=unit \
  ./application/services/print_service.go \
  ./application/services/print_service_test.go \
  ./application/services/template_renderer.go \
  ./application/services/order_service.go \
  2>&1 | grep -E "(PASS|FAIL|RUN|===|---)"

echo ""
echo "Test run complete"
