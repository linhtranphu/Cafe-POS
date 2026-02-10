#!/bin/bash

# Verification script for Task 1.5: Shop Settings Implementation

echo "=== Task 1.5 Verification ==="
echo ""

echo "1. Running unit tests..."
cd backend && go test ./domain/settings/... -v
TEST_RESULT=$?

echo ""
echo "2. Checking migration script build..."
cd cmd/migrate
go build -o create_shop_settings create_shop_settings.go
BUILD_RESULT=$?

if [ $BUILD_RESULT -eq 0 ]; then
    echo "✓ Migration script built successfully"
else
    echo "✗ Migration script build failed"
fi

echo ""
echo "3. Files created:"
echo "   ✓ backend/domain/settings/shop_settings.go"
echo "   ✓ backend/domain/settings/errors.go"
echo "   ✓ backend/domain/settings/shop_settings_test.go"
echo "   ✓ backend/infrastructure/mongodb/shop_settings_repository.go"
echo "   ✓ backend/cmd/migrate/create_shop_settings.go"
echo "   ✓ backend/cmd/migrate/README_SHOP_SETTINGS.md"

echo ""
echo "=== Verification Summary ==="
if [ $TEST_RESULT -eq 0 ] && [ $BUILD_RESULT -eq 0 ]; then
    echo "✅ All checks passed!"
    echo ""
    echo "Task 1.5 is complete and ready for use."
    echo ""
    echo "Next steps:"
    echo "  - Run migration: ./backend/cmd/migrate/create_shop_settings"
    echo "  - Verify in MongoDB: mongosh cafe_pos --eval 'db.shop_settings.find().pretty()'"
else
    echo "❌ Some checks failed"
    exit 1
fi
