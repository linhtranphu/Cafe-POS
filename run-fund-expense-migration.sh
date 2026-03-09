#!/bin/bash

# Script to run Fund-Expense Integration migration
# This script connects to MongoDB and executes the migration

set -e

echo "🚀 Running Fund-Expense Integration Migration..."
echo ""

# Check if MongoDB is running
if ! docker ps | grep -q mongo; then
    echo "❌ Error: MongoDB container is not running"
    echo "Please start MongoDB first with: docker-compose up -d"
    exit 1
fi

# Get the actual MongoDB container name
MONGO_CONTAINER=$(docker ps --format "{{.Names}}" | grep -i mongo)

# Get MongoDB credentials from .env or use defaults
if [ -f .env ]; then
    source .env
fi
MONGO_USER=${MONGO_INITDB_ROOT_USERNAME:-admin}
MONGO_PASS=${MONGO_INITDB_ROOT_PASSWORD:-password}

echo "📡 Connecting to MongoDB container: $MONGO_CONTAINER"
echo "   Using credentials: $MONGO_USER"
echo ""

# Run the migration script with authentication
docker exec -i $MONGO_CONTAINER mongosh \
  --username "$MONGO_USER" \
  --password "$MONGO_PASS" \
  --authenticationDatabase admin \
  cafe_pos < mongodb-init/migrate-fund-expense-integration.js

echo ""
echo "✅ Migration script executed"
echo ""
echo "📋 Next steps:"
echo "   1. Verify the migration results above"
echo "   2. Check that all expenses have 'paid_from_fund' and 'fund_transaction_id' fields"
echo "   3. Check that all fund_transactions have 'source_type' and 'source_id' fields"
echo "   4. Proceed with backend code implementation (Task 2)"
